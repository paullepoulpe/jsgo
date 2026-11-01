#!/usr/bin/env python3
"""Behavioural oracle for jsgo's fast mode (PLAN.md §8; docs/behavioural-oracle.md).

Compiles each test of a suite with jsgo in compat mode and in fast mode (same
flags otherwise), runs the uncompiled test and both builds under node with a
deterministic prelude (tools/behave/runner.js), and compares the builds'
behaviour, compat being the reference (tools/behave/core.py).

Suites (phase 1):
  test262   the pinned stratified sample of test262's runtime suite
            (tools/behave/t262.py, tools/behave/test262/)
  closure   Closure Compiler's runtime tests (tools/behave/closure_rt.py)

Commands:
  run      run suites; writes <out>/results.jsonl, the compat builds, the
           per-test details of every difference, and <out>/summary.{md,json}.
           With --java, also checks every test that passes uncompiled and fails
           on the compat build against the Java compiler (same command line):
           the same output means upstream behaviour, a different one a port
           bug (<out>/findings/).
  faults   the seeded-fault check: rebuilds the fast builds of a finished run
           with JSGO_FAST_FAULT=<fault> and reports which faults the oracle
           catches (<out>/faults.{md,json}). It needs a jsgo built with the
           faults: `go build -tags jsgo_faults`.
  replay   runs the stored compat builds of a finished run again and checks
           that each behaves as recorded: the runner's determinism (and, after
           a node upgrade, the engine's stability).
  explain  for a finished `faults` run: rebuilds each build that a fault
           changed but no test caught, finds where its output differs from
           the compat build, and runs it with V8 block coverage, to tell
           changes the tests executed from changes they never ran
           (<out>/explain.{md,json}, hunks under <out>/<fault>/explain/).
  report   rewrites the summary of a run directory.
  verify   checks the fetched suites against the digests in pins.json
           (tools/behave/fetch.sh runs it); --update rewrites them.
  probe    test262 only: runs every statically eligible test uncompiled and
           compiles a few tests per feature, the input of `sample`.
  sample   test262 only: derives the feature exclusions and draws the pinned
           sample and per-change subset from a probe.

Usage:
  go build -o /tmp/jsgo ./cmd/jsgo
  bash tools/behave/fetch.sh
  python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-run            # per-change subset
  python3 tools/behave/behave.py run --jsgo /tmp/jsgo --out /tmp/behave-full --subset full --java
  go build -tags jsgo_faults -o /tmp/jsgo-faults ./cmd/jsgo
  python3 tools/behave/behave.py faults --jsgo /tmp/jsgo-faults --baseline /tmp/behave-run --out /tmp/behave-faults

The exit status of `run` is 1 when a test that passes on the compat build
behaves differently on the fast build (a regression or behaviour difference)
and is not listed in tools/behave/expect/compat-fast.txt.
"""
import argparse
import collections
import concurrent.futures
import datetime
import hashlib
import json
import os
import shlex
import shutil
import subprocess
import sys
import threading
import time
import traceback

sys.dont_write_bytecode = True
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import core  # noqa: E402
import rule  # noqa: E402
import t262  # noqa: E402
import closure_rt  # noqa: E402

SUITES = ("test262", "closure")
EXPECT_DIR = os.path.join(core.HERE, "expect")
FAULTS = ("inline-args", "drop-unused-init", "swap-args", "fold-negzero", "hook-negation", "collapse-unsafe",
          "inline-vars")


def log(msg):
    print(f"[{datetime.datetime.now().strftime('%H:%M:%S')}] {msg}", flush=True)


def suite_root(cache, suite):
    pins = core.read_pins()
    return os.path.join(cache, pins[suite]["dir"])


def load_cases(args, suite):
    root = suite_root(args.cache, suite)
    if suite == "test262":
        cs = t262.cases(root, args.subset)
    else:
        cs = closure_rt.cases(root, args.cache)
    if args.filter:
        cs = [c for c in cs if args.filter in c.id]
    if args.limit:
        cs = cs[:args.limit]
    return cs


# ---------------------------------------------------------------------------------------------
# Running cases


class Ctx:
    def __init__(self, args, out, configs, fault=None, baseline=None):
        self.args = args
        self.out = out
        self.configs = configs
        self.fault = fault
        self.baseline = baseline  # {(suite, id, config): results line}
        self.lock = threading.Lock()
        self.work = os.path.join(out, "work")
        os.makedirs(self.work, exist_ok=True)
        self.results = open(os.path.join(out, "results.jsonl"), "a")
        self.done = 0
        self.t0 = time.monotonic()
        self.fast_flags = shlex.split(getattr(args, "fast_flags", "") or "")
        cc = getattr(args, "compat_cache", "none")
        self.cache = core.CompatCache(cc) if cc and cc != "none" else None
        self.release = core.read_pins()["closure_version"]
        self.node_version = subprocess.run([args.node, "--version"], capture_output=True, text=True).stdout.strip()
        self.jsgo_sha = core.sha256_file(args.jsgo)[:16]
        self.astdiff = rule.astdiff_binary(args.cache)
        self.extern_names = rule.extern_names(self.astdiff, args.cache)

    def emit(self, lines):
        with self.lock:
            for ln in lines:
                self.results.write(json.dumps(ln, sort_keys=True) + "\n")
            self.results.flush()


def source_ids(case):
    ids = set()
    for inp in case.inputs:
        with open(os.path.join(case.root, inp), encoding="utf-8", errors="surrogateescape") as f:
            ids |= core.identifiers(f.read())
    if case.uncompiled:
        ids |= core.identifiers(case.uncompiled)
    return frozenset(ids)


def read_text(path):
    with open(path, encoding="utf-8", errors="surrogateescape") as f:
        return f.read()


def builds_path(out, case, config, mode):
    return os.path.join(out, "builds", case.suite, config, case.safe_id() + f".{mode}.js.gz")


def save_diff(ctx, case, config, cmp, co, fo, compat_js, fast_js, fault=None):
    d = os.path.join(ctx.out, "diffs", case.suite, config, case.safe_id())
    os.makedirs(d, exist_ok=True)
    for name, src in (("compat.js", compat_js), ("fast.js", fast_js)):
        if src and os.path.exists(src):
            shutil.copyfile(src, os.path.join(d, name))
    with open(os.path.join(d, "outcomes.json"), "w") as f:
        json.dump({"cmp": cmp, "compat": co, "fast": fo, "fault": fault,
                   "argv_compat": case.compile_argv("jsgo", config, "compat", "compat.js"),
                   "argv_fast": case.compile_argv("jsgo", config, "fast", "fast.js"),
                   "cwd": case.root}, f, indent=1, sort_keys=True)
    return os.path.relpath(d, ctx.out)


def compat_build(ctx, case, config, cjs):
    """Compiles (or takes from the compat cache) the compat build of case into
    cjs. Returns (compile result, text, cached runs, cache key)."""
    a = ctx.args
    key = core.compat_key(case, config, ctx.release) if ctx.cache else None
    if key:
        hit = ctx.cache.get(key)
        if hit is not None:
            comp, text, runs = hit
            if comp["ok"]:
                with open(cjs, "w", encoding="utf-8", errors="surrogateescape") as f:
                    f.write(text)
            with ctx.lock:
                ctx.cache.hits += 1
            return comp, text or "", runs, key
        with ctx.lock:
            ctx.cache.misses += 1
    cc = core.compile_case(a.jsgo, case, config, "compat", cjs)
    text = read_text(cjs) if cc["ok"] else ""
    if key:
        ctx.cache.put(key, cc, text if cc["ok"] else None, {}, ctx.jsgo_sha)
    return cc, text, {}, key


def run_cached(ctx, case, prog, key, runs, coverage_dir=None):
    """Runs prog (a compat build or the uncompiled program), or takes its run
    from the cache when no coverage is asked for."""
    rk = core.run_key(case, ctx.node_version)
    if coverage_dir is None and rk in runs:
        with ctx.lock:
            ctx.cache.run_hits += 1
        return runs[rk]
    r = core.run_program(ctx.args.node, prog, case, prog[:-3], coverage_dir=coverage_dir)
    if key and ctx.cache and not r["timeout"]:
        _cache_run(ctx, key, rk, r)
    return r


def _cache_run(ctx, key, rk, r):
    hit = ctx.cache.get(key)
    if hit is not None:
        ctx.cache.put(key, hit[0], None, {rk: r}, ctx.jsgo_sha)


def coverage_of(cov_dir):
    return rule.Coverage(core.coverage_ranges(cov_dir)) if cov_dir and os.path.isdir(cov_dir) else None


def apply_rule(ctx, case, config, work, co, ctext, cjs, cfile_ok, fo, ftext, fjs, cmp, diff, fast_cov,
               compat_cov_dir):
    """The coverage rule for one differing build. compat_cov_dir: where the
    compat build ran under coverage, or None to run it now if a deletion needs
    it."""
    if "error" in diff or diff.get("isomorphic"):
        return rule.judge(diff, cmp, co, None, None)[0]
    compat_cov = coverage_of(compat_cov_dir)
    if compat_cov is None and any(c["kind"] == "delete" for c in diff["changes"]) and cfile_ok:
        cov = os.path.join(work, config + ".compat.cov")
        shutil.rmtree(cov, ignore_errors=True)
        r = core.run_program(ctx.args.node, cjs, case, cjs[:-3] + ".cov", coverage_dir=cov)
        o = core.outcome(case, co["compile"], r)
        if core.behaviour(o) != core.behaviour(co):
            return {"verdict": "blind", "reason": "compat-nondeterministic"}
        compat_cov = coverage_of(cov)
    summary, details = rule.judge(diff, cmp, co, fast_cov, compat_cov)
    if details:
        d = os.path.join(ctx.out, "rule", case.suite, config)
        os.makedirs(d, exist_ok=True)
        sigs = rule.signatures(diff, ctext, ftext)
        with open(os.path.join(d, case.safe_id() + ".json"), "w") as f:
            json.dump({"summary": summary, "cmp": cmp, "id": case.id,
                       "changes": [dict(c, verdict=v, ran=r, sig=g) for (c, v, r), g in zip(details, sigs)],
                       "units_a": diff["units_a"], "units_b": diff["units_b"]}, f, sort_keys=True)
        core.write_gz(os.path.join(d, case.safe_id() + ".fast.js.gz"), ftext.encode("utf-8", "surrogateescape"))
    return summary


def run_case(ctx, case):
    """Runs one case under every configuration (compat, fast and the uncompiled
    program) and returns its results lines."""
    a = ctx.args
    work = os.path.join(ctx.work, case.suite + "-" + case.safe_id())
    os.makedirs(work, exist_ok=True)
    unc = None
    if case.uncompiled is not None:
        prog = os.path.join(work, "uncompiled.js")
        with open(prog, "w", encoding="utf-8", errors="surrogateescape") as f:
            f.write(case.uncompiled)
        ukey = runs = None
        if ctx.cache:
            ukey = "u" + hashlib.sha256(case.uncompiled.encode("utf-8", "surrogateescape")).hexdigest()
            hit = ctx.cache.get(ukey)
            runs = hit[2] if hit else {}
            if hit is None:
                ctx.cache.put(ukey, {"ok": False, "uncompiled": True}, None, {}, ctx.jsgo_sha)
        unc = core.outcome(case, None, run_cached(ctx, case, prog, ukey, runs or {}) if ctx.cache else
                           core.run_program(a.node, prog, case, os.path.join(work, "uncompiled")))
    ids = None
    lines = []
    for config in case.runs_in(ctx.configs):
        cjs = os.path.join(work, config + ".compat.js")
        fjs = os.path.join(work, config + ".fast.js")
        cc, ctext, cruns, ckey = compat_build(ctx, case, config, cjs)
        if cc["ok"]:
            core.write_gz(builds_path(ctx.out, case, config, "compat"), ctext.encode("utf-8", "surrogateescape"))
        fc = core.compile_case(a.jsgo, case, config, "fast", fjs, flags=ctx.fast_flags)
        ftext = read_text(fjs) if fc["ok"] else ""
        identical = cc["ok"] and fc["ok"] and ctext == ftext
        diff = None
        if cc["ok"] and fc["ok"] and not identical:
            diff = rule.locate(ctx.astdiff, cjs, fjs, rule.renames_properties(config, case),
                               rule.fixed_names_file(case, ctx.extern_names, work))
        need_cov = diff is not None and "error" not in diff and not diff["isomorphic"]
        ccov = fcov = None
        if need_cov:
            ccov, fcov = os.path.join(work, config + ".compat.cov"), os.path.join(work, config + ".fast.cov")
            for c in (ccov, fcov):
                shutil.rmtree(c, ignore_errors=True)
        cr = None
        if cc["ok"]:
            if ctx.cache:
                cr = run_cached(ctx, case, cjs, ckey, cruns, coverage_dir=ccov)
            else:
                cr = core.run_program(a.node, cjs, case, cjs[:-3], coverage_dir=ccov)
        co = core.outcome(case, cc, cr)
        if fc["ok"] and identical and not a.rerun_identical:
            fr = dict(cr, reused=True)
        else:
            fr = core.run_program(a.node, fjs, case, fjs[:-3], coverage_dir=fcov) if fc["ok"] else None
        fo = core.outcome(case, fc, fr)
        if (cc["ok"] or fc["ok"]) and core.behaviour(co) != core.behaviour(fo) and ids is None:
            ids = source_ids(case)
        cmp = core.compare(co, fo, ctext, ftext, ids or frozenset())
        if identical and cmp != "same":
            cmp = "nondeterministic"
        line = {"suite": case.suite, "id": case.id, "stratum": case.stratum, "config": config,
                "uncompiled": unc, "compat": co, "fast": fo, "cmp": cmp, "identical": identical}
        if identical:
            line["rule"] = {"verdict": "identical"}
        elif diff is not None:
            line["rule"] = apply_rule(ctx, case, config, work, co, ctext, cjs, cc["ok"], fo, ftext, fjs, cmp, diff,
                                      coverage_of(fcov), ccov)
        if cmp not in ("same", "same-renamed"):
            line["detail"] = save_diff(ctx, case, config, cmp, co, fo, cjs, fjs)
        lines.append(line)
    if not a.keep_work:
        shutil.rmtree(work, ignore_errors=True)
    return lines


def fault_case(ctx, case):
    """Rebuilds the fast builds of one case with the seeded fault and compares
    them with the baseline's compat outcomes, under the coverage rule."""
    a = ctx.args
    work = os.path.join(ctx.work, case.suite + "-" + case.safe_id())
    os.makedirs(work, exist_ok=True)
    lines = []
    ids = None
    for config in case.runs_in(ctx.configs):
        base = ctx.baseline.get((case.suite, case.id, config))
        if base is None:
            continue
        co = base["compat"]
        cpath = os.path.join(a.baseline, "builds", case.suite, config, case.safe_id() + ".compat.js.gz")
        ctext = core.read_gz(cpath).decode("utf-8", "surrogateescape") if os.path.exists(cpath) else ""
        fjs = os.path.join(work, config + ".fast.js")
        fc = core.compile_case(a.jsgo, case, config, "fast", fjs, fault=ctx.fault, flags=ctx.fast_flags)
        ftext = read_text(fjs) if fc["ok"] else ""
        triggered = not (fc["ok"] and co["compile"]["ok"] and ftext == ctext) and not (
            not fc["ok"] and not co["compile"]["ok"] and fc["errors"] == co["compile"]["errors"])
        if not triggered:
            lines.append({"suite": case.suite, "id": case.id, "config": config, "fault": ctx.fault,
                          "triggered": False, "cmp": "same", "compat_pass": co["pass"],
                          "compat_signal": rule.signal(co), "rule": {"verdict": "identical"}})
            continue
        cjs = os.path.join(work, config + ".compat.js")
        with open(cjs, "w", encoding="utf-8", errors="surrogateescape") as f:
            f.write(ctext)
        diff = None
        if fc["ok"] and co["compile"]["ok"]:
            diff = rule.locate(ctx.astdiff, cjs, fjs, rule.renames_properties(config, case),
                               rule.fixed_names_file(case, ctx.extern_names, work))
        need_cov = diff is not None and "error" not in diff and not diff["isomorphic"]
        fcov = os.path.join(work, config + ".fast.cov") if need_cov else None
        if fcov:
            shutil.rmtree(fcov, ignore_errors=True)
        fr = core.run_program(a.node, fjs, case, fjs[:-3], coverage_dir=fcov) if fc["ok"] else None
        fo = core.outcome(case, fc, fr)
        if ids is None:
            ids = source_ids(case)
        cmp = core.compare(co, fo, ctext, ftext, ids)
        line = {"suite": case.suite, "id": case.id, "config": config, "fault": ctx.fault, "triggered": True,
                "cmp": cmp, "compat_pass": co["pass"], "compat_signal": rule.signal(co), "fast_pass": fo["pass"],
                "fast_why": fo["why"]}
        if diff is not None:
            line["rule"] = apply_rule(ctx, case, config, work, co, ctext, cjs, co["compile"]["ok"], fo, ftext, fjs,
                                      cmp, diff, coverage_of(fcov), None)
        else:  # one side did not compile: a compile difference is checked by the strict layers
            line["rule"] = {"verdict": "fail" if cmp in core.SIGNAL_DIFFS else "blind", "reason": "compile"}
        if cmp in core.SIGNAL_DIFFS and a.keep_diffs:
            line["detail"] = save_diff(ctx, case, config, cmp, co, fo, cjs, fjs, fault=ctx.fault)
        lines.append(line)
    shutil.rmtree(work, ignore_errors=True)
    return lines


def run_parallel(ctx, cases, fn, label):
    total = len(cases)
    if not total:
        return
    log(f"{label}: {total} case(s), {len(ctx.configs)} configuration(s), {ctx.args.jobs} job(s)")
    every = max(1, total // 20)
    with concurrent.futures.ThreadPoolExecutor(max_workers=ctx.args.jobs) as ex:
        futs = {ex.submit(fn, ctx, c): c for c in cases}
        for fut in concurrent.futures.as_completed(futs):
            c = futs[fut]
            try:
                lines = fut.result()
            except Exception as e:  # a harness bug: report it, keep going
                log(f"  harness error on {c.suite}/{c.id}: {e!r}\n{traceback.format_exc()}")
                lines = [{"suite": c.suite, "id": c.id, "config": cfg, "harness_error": repr(e)} for cfg in ctx.configs]
            ctx.emit(lines)
            ctx.done += 1
            if ctx.done % every == 0 or ctx.done == total:
                el = time.monotonic() - ctx.t0
                eta = el / ctx.done * (total - ctx.done)
                log(f"  {ctx.done}/{total} cases, {el / 60:.1f} min elapsed, ~{eta / 60:.1f} min left")


def read_results(out):
    path = os.path.join(out, "results.jsonl")
    lines = []
    if os.path.exists(path):
        with open(path) as f:
            for ln in f:
                ln = ln.strip()
                if ln:
                    try:
                        lines.append(json.loads(ln))
                    except ValueError:
                        pass  # a line cut short by a restart
    return lines


def latest(lines):
    """The last line for each (suite, id, config)."""
    d = {}
    for ln in lines:
        d[(ln["suite"], ln["id"], ln["config"])] = ln
    return d


def cmd_run(args):
    out = os.path.abspath(args.out)
    os.makedirs(out, exist_ok=True)
    configs = args.configs.split(",")
    meta_path = os.path.join(out, "run.json")
    meta = {}
    if args.resume and os.path.exists(meta_path):
        with open(meta_path) as f:
            meta = json.load(f)
    elif os.path.exists(os.path.join(out, "results.jsonl")):
        os.remove(os.path.join(out, "results.jsonl"))
        shutil.rmtree(os.path.join(out, "diffs"), ignore_errors=True)
    meta.update({"jsgo": os.path.abspath(args.jsgo), "jsgo_sha256": core.sha256_file(args.jsgo),
                 "node": args.node, "node_version": subprocess.run([args.node, "--version"], capture_output=True,
                                                                   text=True).stdout.strip(),
                 "subset": args.subset, "configs": configs, "suites": args.suites.split(","), "jobs": args.jobs,
                 "pins": core.read_pins()})
    meta.setdefault("wall", {})
    done = latest(read_results(out)) if args.resume else {}
    ctx = Ctx(args, out, configs)
    for suite in args.suites.split(","):
        cases = load_cases(args, suite)
        todo = [c for c in cases if not all((c.suite, c.id, cfg) in done for cfg in c.runs_in(configs))]
        t0 = time.monotonic()
        ctx.done, ctx.t0 = 0, time.monotonic()
        run_parallel(ctx, todo, run_case, suite)
        meta["wall"][suite] = meta["wall"].get(suite, 0) + round(time.monotonic() - t0, 1)
        meta.setdefault("cases", {})[suite] = len(cases)
        with open(meta_path, "w") as f:
            json.dump(meta, f, indent=1, sort_keys=True)
    ctx.results.close()
    if args.java:
        t0 = time.monotonic()
        java_check(args, out)
        meta["wall"]["java"] = round(time.monotonic() - t0, 1)
        with open(meta_path, "w") as f:
            json.dump(meta, f, indent=1, sort_keys=True)
    if not args.keep_work:
        shutil.rmtree(os.path.join(out, "work"), ignore_errors=True)
    bad = write_report(out)
    return 1 if bad else 0


# ---------------------------------------------------------------------------------------------
# Java check: uncompiled passes, compat fails


def java_classpath(cache):
    pin = {}
    with open(os.path.join(core.REPO, "tools", "oracle", "pin.env")) as f:
        for line in f:
            line = line.strip()
            if line and not line.startswith("#"):
                k, v = line.split("=", 1)
                pin[k] = v.strip('"')
    jar = os.path.join(cache, pin["CLOSURE_JAR"])
    build = os.path.join(cache, "behave", "java-build")
    cls = os.path.join(build, "CliBatch.class")
    src = os.path.join(core.REPO, "tools", "oracle", "java", "CliBatch.java")
    if not os.path.exists(cls) or os.path.getmtime(cls) < os.path.getmtime(src):
        os.makedirs(build, exist_ok=True)
        subprocess.run(["javac", "-nowarn", "-cp", jar, "-d", build, src], check=True,
                       env=core.clean_env(core.JAVA_ENV_DROP))
    return jar + os.pathsep + build


def category(o):
    if o["why"] == "compile-error":
        c = o["compile"]
        if c["crash"]:
            return "compile-crash"
        if "INTERNAL COMPILER ERROR" in c.get("stderr", ""):
            return "compile-internal-error"
        return "compile-error:" + (",".join(c["errors"][:2]) or f"exit-{c['exit']}")
    return "runtime:" + o["why"]


def java_check(args, out):
    """Compiles with Java every (case, config) whose uncompiled program passes
    and whose compat build fails, and compares Java's output with jsgo's."""
    lines = latest(read_results(out))
    todo = [ln for ln in lines.values()
            if ln.get("uncompiled") and ln["uncompiled"]["pass"] and ln.get("compat") and not ln["compat"]["pass"]]
    log(f"java check: {len(todo)} build(s) pass uncompiled and fail on compat")
    if not todo:
        return
    cp = java_classpath(args.cache)
    cases = {}
    for suite in {ln["suite"] for ln in todo}:
        for c in load_cases(args, suite):
            cases[(suite, c.id)] = c
    jdir = os.path.join(out, "java")
    shutil.rmtree(jdir, ignore_errors=True)
    os.makedirs(jdir)
    jobs_by_root = collections.defaultdict(list)
    for i, ln in enumerate(sorted(todo, key=lambda x: (x["suite"], x["id"], x["config"]))):
        c = cases.get((ln["suite"], ln["id"]))
        if c is None:
            continue
        js = os.path.join(jdir, f"{i}.js")
        res = os.path.join(jdir, f"{i}.out")
        argv = c.compile_argv([], ln["config"], "java", js)
        jobs_by_root[c.root].append((i, ln, c, js, res, argv))
    results = []
    for root, jobs in jobs_by_root.items():
        # Split the jobs over JVMs, one per job slot.
        n = max(1, min(args.jobs, len(jobs) // 50 + 1))
        chunks = [jobs[k::n] for k in range(n)]
        procs = []
        for k, chunk in enumerate(chunks):
            jobfile = os.path.join(jdir, f"jobs-{k}.tsv")
            with open(jobfile, "w") as f:
                for (i, ln, c, js, res, argv) in chunk:
                    f.write("\t".join([res] + argv) + "\n")
            procs.append(subprocess.Popen(["java", "-Xss8m", "-Xmx2g", "-cp", cp, "CliBatch"], cwd=root,
                                          stdin=open(jobfile), stdout=subprocess.DEVNULL,
                                          stderr=open(jobfile + ".log", "w"), env=core.clean_env(core.JAVA_ENV_DROP)))
        for p in procs:
            p.wait()
        for (i, ln, c, js, res, argv) in jobs:
            results.append(java_compare(args, out, ln, c, js, res))
    fdir = os.path.join(out, "findings")
    with open(os.path.join(out, "findings.jsonl"), "w") as f:
        for r in results:
            f.write(json.dumps(r, sort_keys=True) + "\n")
    port = [r for r in results if r["kind"] == "port"]
    log(f"java check: {len(results) - len(port)} upstream, {len(port)} port difference(s) (details in {fdir})")


def java_compare(args, out, ln, case, js, res):
    text = ""
    if os.path.exists(res):
        with open(res, encoding="utf-8", errors="replace") as f:
            text = f.read()
    jexit = None
    if "---EXIT---" in text:
        body, _, tail = text.rpartition("---EXIT---")
        try:
            jexit = int(tail.strip())
        except ValueError:
            jexit = tail.strip()
    else:
        body = text
    jerrors = sorted(set(core.ERROR_CODE_RE.findall(body)))
    co = ln["compat"]
    cpath = os.path.join(out, "builds", case.suite, ln["config"], case.safe_id() + ".compat.js.gz")
    gtext = core.read_gz(cpath) if os.path.exists(cpath) else None
    jtext = None
    if os.path.exists(js):
        with open(js, "rb") as f:
            jtext = f.read()
    same_exit = (jexit == 0) == bool(co["compile"]["ok"])
    same = same_exit and gtext == jtext and (co["compile"]["ok"] or jerrors == co["compile"]["errors"])
    r = {"suite": case.suite, "id": case.id, "config": ln["config"], "category": category(co),
         "kind": "upstream" if same else "port", "java_exit": jexit, "java_errors": jerrors}
    if not same:
        d = os.path.join(out, "findings", case.suite, ln["config"], case.safe_id())
        os.makedirs(d, exist_ok=True)
        if gtext is not None:
            with open(os.path.join(d, "jsgo.js"), "wb") as f:
                f.write(gtext)
        if jtext is not None:
            with open(os.path.join(d, "java.js"), "wb") as f:
                f.write(jtext)
        with open(os.path.join(d, "java.stderr"), "w") as f:
            f.write(body)
        with open(os.path.join(d, "info.json"), "w") as f:
            json.dump({"argv": case.compile_argv("jsgo", ln["config"], "compat", "jsgo.js"), "cwd": case.root,
                       "compat": co, "java_exit": jexit, "java_errors": jerrors}, f, indent=1, sort_keys=True)
        # Does Java's own build pass?
        if jtext is not None:
            prefix = os.path.join(d, "java")
            jr = core.run_program(args.node, os.path.join(d, "java.js"), case, prefix)
            r["java_pass"] = core.outcome(case, {"ok": True}, jr)["pass"]
        r["detail"] = os.path.relpath(d, out)
    return r


# ---------------------------------------------------------------------------------------------
# Seeded faults


def check_fault_binary(jsgo, faults):
    """Exits unless jsgo was built with the seeded faults (`go build -tags
    jsgo_faults`) and knows each fault name: a binary without them ignores
    JSGO_FAST_FAULT, so every fault would look harmless."""
    toy = os.path.join(core.HERE, "faults_toy.js")
    env = core.clean_env(extra={"JSGO_FAST_FAULT": "fold-negzero"})
    p = subprocess.run([jsgo, "--fast_mode", "-O", "SIMPLE", "--js", toy], env=env, capture_output=True, text=True)
    if p.returncode != 0 or "-0" in p.stdout:
        sys.exit(f"{jsgo} has no seeded faults (fold-negzero does not fire on {toy}): build it with "
                 "`go build -tags jsgo_faults -o <bin> ./cmd/jsgo`")
    for fault in faults:  # a fault build refuses an unknown fault name at start-up
        p = subprocess.run([jsgo, "--version"], env=core.clean_env(extra={"JSGO_FAST_FAULT": fault}),
                           capture_output=True, text=True)
        if p.returncode != 0:
            sys.exit(f"{jsgo} rejects JSGO_FAST_FAULT={fault}: {p.stderr.strip().splitlines()[:1]}")


def cmd_faults(args):
    out = os.path.abspath(args.out)
    os.makedirs(out, exist_ok=True)
    args.baseline = os.path.abspath(args.baseline)
    base = latest(read_results(args.baseline))
    with open(os.path.join(args.baseline, "run.json")) as f:
        bmeta = json.load(f)
    configs = bmeta["configs"]
    faults = args.faults.split(",")
    check_fault_binary(args.jsgo, faults)
    summary = {"baseline": args.baseline, "jsgo": os.path.abspath(args.jsgo), "faults": {}}
    path = os.path.join(out, "faults.json")
    if args.resume and os.path.exists(path):
        with open(path) as f:
            summary = json.load(f)
    for fault in faults:
        if fault in summary["faults"]:
            continue
        fout = os.path.join(out, fault)
        shutil.rmtree(fout, ignore_errors=True)
        os.makedirs(fout)
        ctx = Ctx(args, fout, configs, fault=fault, baseline=base)
        t0 = time.monotonic()
        suites = sorted({k[0] for k in base})
        for suite in suites:
            ids = {k[1] for k in base if k[0] == suite}
            cases = [c for c in load_cases(args, suite) if c.id in ids]
            ctx.done, ctx.t0 = 0, time.monotonic()
            run_parallel(ctx, cases, fault_case, f"fault {fault}, {suite}")
        ctx.results.close()
        lines = read_results(fout)
        summary["faults"][fault] = fault_summary(lines, round(time.monotonic() - t0, 1))
        shutil.rmtree(os.path.join(fout, "work"), ignore_errors=True)
        with open(path, "w") as f:
            json.dump(summary, f, indent=1, sort_keys=True)
    write_faults_md(out, summary)
    return 0


def fault_summary(lines, wall):
    s = {"wall": wall, "builds": 0, "triggered": 0, "signal_builds": 0, "caught": 0, "weak": 0,
         "caught_cases": 0, "signal_cases": 0, "by_suite": {}, "examples": []}
    caught_cases, signal_cases = set(), set()
    for ln in lines:
        if "harness_error" in ln:
            continue
        key = (ln["suite"], ln["id"])
        bs = s["by_suite"].setdefault(ln["suite"], {"builds": 0, "triggered": 0, "signal_builds": 0, "caught": 0})
        s["builds"] += 1
        bs["builds"] += 1
        if ln["compat_pass"]:
            s["signal_builds"] += 1
            bs["signal_builds"] += 1
            signal_cases.add(key)
        if ln["triggered"]:
            s["triggered"] += 1
            bs["triggered"] += 1
        if ln["cmp"] in core.SIGNAL_DIFFS:
            s["caught"] += 1
            bs["caught"] += 1
            caught_cases.add(key)
            if len(s["examples"]) < 10:
                s["examples"].append(f"{ln['suite']}/{ln['id']} [{ln['config']}] {ln['cmp']}: {ln.get('fast_why')}")
        elif ln["cmp"] not in ("same", "same-renamed"):
            s["weak"] += 1
    s["caught_cases"] = len(caught_cases)
    s["signal_cases"] = len(signal_cases)
    return s


def write_faults_md(out, summary):
    rows = []
    caught = 0
    for fault, s in summary["faults"].items():
        if s["caught"]:
            caught += 1
        pct = 100.0 * s["caught"] / s["signal_builds"] if s["signal_builds"] else 0
        cpct = 100.0 * s["caught_cases"] / s["signal_cases"] if s["signal_cases"] else 0
        rows.append(f"| {fault} | {s['triggered']} | {s['caught']} | {pct:.2f}% | {s['caught_cases']} | {cpct:.2f}% "
                    f"| {s['weak']} | {'yes' if s['caught'] else 'no'} | {s['wall'] / 60:.1f} |")
    n = len(summary["faults"])
    md = [f"# Seeded faults: {caught} of {n} caught (mutation score {100.0 * caught / n if n else 0:.0f}%)", "",
          f"Baseline: `{summary['baseline']}`", "",
          "triggered: builds whose output the fault changed; caught: signal builds (compat passes) whose "
          "behaviour differs; % of the signal builds; cases: tests caught in at least one configuration.", "",
          "| fault | triggered | caught builds | % signal builds | caught tests | % signal tests | weak | caught | min |",
          "|---|---:|---:|---:|---:|---:|---:|---|---:|"] + rows + [""]
    for fault, s in summary["faults"].items():
        if s["examples"]:
            md.append(f"## {fault}")
            md += [f"- {e}" for e in s["examples"]]
            md.append("")
    with open(os.path.join(out, "faults.md"), "w") as f:
        f.write("\n".join(md))
    print("\n".join(md[:len(rows) + 8]))


# ---------------------------------------------------------------------------------------------
# Report


def read_expect(name):
    path = os.path.join(EXPECT_DIR, name)
    if not os.path.exists(path):
        return set()
    with open(path) as f:
        return {ln.split("#", 1)[0].strip() for ln in f if ln.split("#", 1)[0].strip()}


def write_report(out):
    """Writes summary.json and summary.md; returns the unexpected signal
    differences."""
    lines = list(latest(read_results(out)).values())
    meta = {}
    if os.path.exists(os.path.join(out, "run.json")):
        with open(os.path.join(out, "run.json")) as f:
            meta = json.load(f)
    allowed = read_expect("compat-fast.txt")
    table = collections.OrderedDict()
    diffs = collections.defaultdict(list)
    harness_errors = []
    for ln in sorted(lines, key=lambda x: (x["suite"], x["config"], x["id"])):
        if "harness_error" in ln:
            harness_errors.append(ln)
            continue
        t = table.setdefault((ln["suite"], ln["config"]), collections.Counter())
        t["tests"] += 1
        if ln["uncompiled"] is not None:
            t["uncompiled_run"] += 1
            t["uncompiled_pass"] += ln["uncompiled"]["pass"]
        t["compat_compiled"] += ln["compat"]["compile"]["ok"]
        t["fast_compiled"] += ln["fast"]["compile"]["ok"]
        t["compat_pass"] += ln["compat"]["pass"]
        t["fast_pass"] += ln["fast"]["pass"]
        t["identical"] += ln["identical"]
        t["cmp:" + ln["cmp"]] += 1
        if ln["uncompiled"] is not None and ln["uncompiled"]["pass"] and not ln["compat"]["pass"]:
            t["uncompiled_pass_compat_fail"] += 1
        if ln["uncompiled"] is not None and not ln["uncompiled"]["pass"] and ln["compat"]["pass"]:
            t["uncompiled_fail_compat_pass"] += 1
        t["compile_secs_compat"] += ln["compat"]["compile"]["secs"]
        t["compile_secs_fast"] += ln["fast"]["compile"]["secs"]
        t["signal"] += rule.signal(ln["compat"])
        t["tests_compat_pass"] += core.passing_tests(ln["compat"])
        r = ln.get("rule")
        if r:
            t["rule:" + r["verdict"]] += 1
            for k in ("vouched", "blind", "fail"):
                t["changes_" + k] += r.get(k, 0)
                t["units_" + k] += r.get("units_" + k, 0)
        if ln["cmp"] not in ("same", "same-renamed"):
            diffs[ln["cmp"]].append(ln)
    unexpected = [ln for c in list(core.SIGNAL_DIFFS) + ["nondeterministic"] for ln in diffs.get(c, [])
                  if f"{ln['suite']}/{ln['id']} {ln['config']}" not in allowed]
    findings = []
    fpath = os.path.join(out, "findings.jsonl")
    if os.path.exists(fpath):
        with open(fpath) as f:
            findings = [json.loads(x) for x in f if x.strip()]
    summary = {"meta": meta, "table": {f"{s}/{c}": dict(v) for (s, c), v in table.items()},
               "differences": {k: [f"{x['suite']}/{x['id']} {x['config']}" for x in v] for k, v in diffs.items()},
               "unexpected": [f"{x['suite']}/{x['id']} {x['config']}" for x in unexpected],
               "harness_errors": [f"{x['suite']}/{x['id']} {x['config']}: {x['harness_error']}" for x in harness_errors],
               "findings": summarize_findings(findings)}
    with open(os.path.join(out, "summary.json"), "w") as f:
        json.dump(summary, f, indent=1, sort_keys=True)
    md = ["# Behavioural oracle: compat vs fast", ""]
    if meta:
        md += [f"- jsgo: `{meta.get('jsgo')}` (sha256 {meta.get('jsgo_sha256', '')[:12]}), node {meta.get('node_version')}",
               f"- subset: {meta.get('subset')}; configurations: {', '.join(meta.get('configs', []))}; "
               f"jobs: {meta.get('jobs')}",
               "- wall time: " + ", ".join(f"{k} {v / 60:.1f} min" for k, v in sorted(meta.get("wall", {}).items())),
               ""]
    md += ["| suite / config | tests | uncompiled pass | compat compiled | compat pass (signal) | fast pass "
           "| fast bytes = compat | same | same-renamed | regression | behaviour | fast-only-pass | weak "
           "| unc. pass, compat fail |",
           "|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|"]
    for (s, c), t in table.items():
        unc = f"{t['uncompiled_pass']}/{t['uncompiled_run']}" if t["uncompiled_run"] else "n/a"
        md.append(f"| {s} / {c} | {t['tests']} | {unc} | {t['compat_compiled']} | {t['compat_pass']} | {t['fast_pass']} "
                  f"| {t['identical']} | {t['cmp:same']} | {t['cmp:same-renamed']} | {t['cmp:regression']} "
                  f"| {t['cmp:behaviour']} | {t['cmp:fast-only-pass']} | {t['cmp:weak']} "
                  f"| {t['uncompiled_pass_compat_fail']} |")
    md.append("")
    md += ["Coverage rule (PLAN.md §8), per build that is not byte-identical: isomorphic (every function "
           "still alpha-isomorphic), vouched (every change ran and behaved the same in a test that checks "
           "something), blind (some change never ran, or ran unchecked), fail (behaved differently); then the "
           "changes and the changed functions by verdict.", "",
           "| suite / config | signal builds | tests passing (compat) | identical | isomorphic | vouched | blind "
           "| fail | changes vouched / blind / fail | changed functions vouched / blind / fail |",
           "|---|---:|---:|---:|---:|---:|---:|---:|---|---|"]
    for (s_, c), t in table.items():
        md.append(f"| {s_} / {c} | {t['signal']} | {t['tests_compat_pass']} | {t['rule:identical']} "
                  f"| {t['rule:isomorphic']} | {t['rule:vouched']} | {t['rule:blind']} | {t['rule:fail']} "
                  f"| {t['changes_vouched']} / {t['changes_blind']} / {t['changes_fail']} "
                  f"| {t['units_vouched']} / {t['units_blind']} / {t['units_fail']} |")
    md.append("")
    sigs = signature_summary(out)
    if sigs["total"]:
        summary["signatures"] = sigs
        with open(os.path.join(out, "summary.json"), "w") as f:
            json.dump(summary, f, indent=1, sort_keys=True)
        md += [f"Distinct changes (the same function changed the same way in several builds counts once): "
               f"{sigs['total']}; vouched somewhere and failing nowhere {sigs['vouched']}, blind everywhere "
               f"{sigs['blind']}, failing somewhere {sigs['fail']}.", ""]
        if sigs["blind_examples"]:
            md.append("Blind everywhere (examples):")
            md += [f"- {x}" for x in sigs["blind_examples"][:15]]
            md.append("")
    md.append(f"**Unexpected signal differences: {len(unexpected)}**"
              + (" (the run fails)" if unexpected else ""))
    md.append("")
    for cmp_kind in list(core.SIGNAL_DIFFS) + ["nondeterministic", "fast-only-pass", "weak"]:
        if diffs.get(cmp_kind):
            md.append(f"## {cmp_kind} ({len(diffs[cmp_kind])})")
            for ln in diffs[cmp_kind][:200]:
                md.append(f"- {ln['suite']}/{ln['id']} [{ln['config']}]: compat {ln['compat']['why']}, "
                          f"fast {ln['fast']['why']} — `{ln.get('detail', '')}`")
            md.append("")
    if harness_errors:
        md.append(f"## harness errors ({len(harness_errors)})")
        md += [f"- {x}" for x in summary["harness_errors"][:50]]
        md.append("")
    if findings:
        fs = summary["findings"]
        md += [f"## Uncompiled passes, compat fails: {fs['total']} build(s), {fs['upstream']} upstream "
               f"(Java gives the same output), {fs['port']} port difference(s)", ""]
        md += ["| suite / config | category | upstream | port |", "|---|---|---:|---:|"]
        for row in fs["by_category"]:
            md.append(f"| {row[0]} | {row[1]} | {row[2]} | {row[3]} |")
        md.append("")
        if fs["port_list"]:
            md.append("Port differences:")
            md += [f"- {x}" for x in fs["port_list"]]
            md.append("")
    with open(os.path.join(out, "summary.md"), "w") as f:
        f.write("\n".join(md))
    print("\n".join(md[:8 + len(table)]))
    return unexpected


def signature_summary(out):
    """Aggregates the changes of out/rule/ by signature (rule.signature)."""
    import glob
    by = {}
    ex = {}
    for path in glob.glob(os.path.join(out, "rule", "*", "*", "*.json")):
        try:
            with open(path) as f:
                d = json.load(f)
        except ValueError:
            continue
        suite, config = path.split(os.sep)[-3:-1]
        for c in d["changes"]:
            by.setdefault(c["sig"], collections.Counter())[c["verdict"]] += 1
            if c["sig"] not in ex:
                u = d["units_b"][c["unit_b"]] if c["unit_b"] >= 0 else d["units_a"][c["unit_a"]]
                ex[c["sig"]] = f"{suite}/{d.get('id')} [{config}] {c['kind']} {c['what'][:50]} in {u['name']} (line {u['line']})"
    res = {"total": len(by), "vouched": 0, "blind": 0, "fail": 0, "blind_examples": []}
    for sig in sorted(by):
        v = by[sig]
        if v.get("fail"):
            res["fail"] += 1
        elif v.get("vouched"):
            res["vouched"] += 1
        else:
            res["blind"] += 1
            res["blind_examples"].append(ex[sig])
    return res


def summarize_findings(findings):
    by = collections.Counter()
    for r in findings:
        by[(f"{r['suite']}/{r['config']}", r["category"], r["kind"])] += 1
    cats = sorted({(k[0], k[1]) for k in by}, key=lambda k: (k[0], -(by[(k[0], k[1], "upstream")] + by[(k[0], k[1], "port")])))
    return {"total": len(findings), "upstream": sum(1 for r in findings if r["kind"] == "upstream"),
            "port": sum(1 for r in findings if r["kind"] == "port"),
            "by_category": [[s, c, by[(s, c, "upstream")], by[(s, c, "port")]] for s, c in cats],
            "port_list": [f"{r['suite']}/{r['id']} [{r['config']}] {r['category']} java_pass={r.get('java_pass')} "
                          f"`{r.get('detail', '')}`" for r in findings if r["kind"] == "port"]}


def explain_case(args, fault, case, config, base_line, out_dir):
    """Rebuilds one faulted build, diffs it against the compat build and runs
    it under coverage. Returns a dict for explain.json."""
    work = os.path.join(out_dir, "work", f"{case.suite}-{case.safe_id()}-{config}")
    shutil.rmtree(work, ignore_errors=True)
    os.makedirs(work)
    cpath = builds_path(args.baseline, case, config, "compat")
    ctext = core.read_gz(cpath).decode("utf-8", "surrogateescape")
    fjs = os.path.join(work, "fault.js")
    comp = core.compile_case(args.jsgo, case, config, "fast", fjs, fault=fault)
    r = {"suite": case.suite, "id": case.id, "config": config}
    if not comp["ok"]:
        r["class"] = "compile-error"
        shutil.rmtree(work, ignore_errors=True)
        return r
    ftext = read_text(fjs)
    ids = source_ids(case)
    hunks = core.changed_ranges(ctext, ftext, ids)
    cov = os.path.join(work, "cov")
    run = core.run_program(args.node, fjs, case, fjs[:-3], coverage_dir=cov)
    ranges = core.coverage_ranges(cov)
    flags = [core.executed(ranges, a, b) for a, b in hunks]
    r.update({"hunks": len(hunks), "executed_hunks": sum(flags), "renaming_only": not hunks,
              "fast_why": core.outcome(case, comp, run)["why"],
              "class": "renaming-only" if not hunks else ("executed" if any(flags) else "not-executed")})
    d = os.path.join(out_dir, case.suite, config, case.safe_id())
    os.makedirs(d, exist_ok=True)
    with open(os.path.join(d, "hunks.txt"), "w", encoding="utf-8", errors="surrogateescape") as f:
        for (a, b), ex in zip(hunks, flags):
            lo, hi = max(0, a - 160), min(len(ftext), b + 160)
            f.write(f"== fault build [{a},{b}) {'EXECUTED' if ex else 'not executed'}\n"
                    f"{ftext[lo:a]}\u27e6{ftext[a:b]}\u27e7{ftext[b:hi]}\n\n")
    shutil.copyfile(fjs, os.path.join(d, "fault.js"))
    with open(os.path.join(d, "compat.js"), "w", encoding="utf-8", errors="surrogateescape") as f:
        f.write(ctext)
    shutil.rmtree(work, ignore_errors=True)
    return r


def cmd_explain(args):
    out = os.path.abspath(args.out)
    args.baseline = os.path.abspath(args.baseline)
    base = latest(read_results(args.baseline))
    with open(os.path.join(out, "faults.json")) as f:
        summary = json.load(f)
    faults = args.faults.split(",") if args.faults else list(summary["faults"])
    check_fault_binary(args.jsgo, faults)
    cases = {}
    for suite in sorted({k[0] for k in base}):
        ids = {k[1] for k in base if k[0] == suite}
        for c in load_cases(args, suite):
            if c.id in ids:
                cases[(suite, c.id)] = c
    report = {}
    path = os.path.join(out, "explain.json")
    if os.path.exists(path):
        with open(path) as f:
            report = json.load(f)
    for fault in faults:
        lines = read_results(os.path.join(out, fault))
        todo = [ln for ln in lines if ln.get("triggered") and ln.get("compat_pass")
                and ln["cmp"] in ("same", "same-renamed") and (ln["suite"], ln["id"]) in cases]
        changed_signal = sum(1 for ln in lines if ln.get("triggered") and ln.get("compat_pass"))
        caught = sum(1 for ln in lines if ln["cmp"] in core.SIGNAL_DIFFS)
        log(f"explain {fault}: {len(todo)} changed signal build(s) not caught")
        out_dir = os.path.join(out, fault, "explain")
        shutil.rmtree(out_dir, ignore_errors=True)
        with concurrent.futures.ThreadPoolExecutor(max_workers=args.jobs) as ex:
            results = list(ex.map(lambda ln: explain_case(args, fault, cases[(ln["suite"], ln["id"])], ln["config"],
                                                          base[(ln["suite"], ln["id"], ln["config"])], out_dir), todo))
        shutil.rmtree(os.path.join(out_dir, "work"), ignore_errors=True)
        counts = collections.Counter(r["class"] for r in results)
        report[fault] = {"builds": len(lines), "changed": sum(1 for ln in lines if ln.get("triggered")),
                         "changed_signal": changed_signal, "caught": caught, "uncaught": len(todo),
                         "classes": dict(counts), "cases": results}
        with open(path, "w") as f:
            json.dump(report, f, indent=1, sort_keys=True)
    md = ["# Seeded faults explained", "",
          "changed: builds whose output the fault changed (the others are byte-identical to the fault-free "
          "build, so nothing could catch them); signal: of those, builds whose compat build passes; uncaught "
          "changes split by whether the run executed a changed range (V8 block coverage).", "",
          "| fault | builds | changed | identical | changed, signal | caught | caught / changed signal "
          "| uncaught: executed | not executed | renaming only | compile error |",
          "|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|"]
    for fault, v in report.items():
        c = v["classes"]
        pct = f"{100.0 * v['caught'] / v['changed_signal']:.0f}%" if v["changed_signal"] else "–"
        md.append(f"| `{fault}` | {v['builds']} | {v['changed']} | {v['builds'] - v['changed']} | {v['changed_signal']} "
                  f"| {v['caught']} | {pct} | {c.get('executed', 0)} | {c.get('not-executed', 0)} "
                  f"| {c.get('renaming-only', 0)} | {c.get('compile-error', 0)} |")
    with open(os.path.join(out, "explain.md"), "w") as f:
        f.write("\n".join(md) + "\n")
    print("\n".join(md))
    return 0


def cmd_replay(args):
    out = os.path.abspath(args.out)
    meta_path = os.path.join(out, "run.json")
    if os.path.exists(meta_path):
        with open(meta_path) as f:
            args.subset = json.load(f).get("subset", args.subset)
    lines = [ln for ln in latest(read_results(out)).values()
             if "harness_error" not in ln and ln["compat"]["compile"]["ok"]]
    if args.limit:
        lines = lines[:args.limit]
    cases = {}
    for suite in sorted({ln["suite"] for ln in lines}):
        for c in load_cases(args, suite):
            cases[(suite, c.id)] = c
    work = os.path.join(out, "replay-work")
    os.makedirs(work, exist_ok=True)
    lock = threading.Lock()
    bad = []

    def one(ln):
        c = cases.get((ln["suite"], ln["id"]))
        if c is None:
            return
        src = builds_path(out, c, ln["config"], "compat")
        prog = os.path.join(work, f"{c.suite}-{c.safe_id()}-{ln['config']}.js")
        with open(prog, "wb") as f:
            f.write(core.read_gz(src))
        o = core.outcome(c, ln["compat"]["compile"], core.run_program(args.node, prog, c, prog[:-3]))
        os.remove(prog)
        if core.behaviour(o) != core.behaviour(ln["compat"]):
            with lock:
                bad.append(f"{c.suite}/{c.id} [{ln['config']}]: recorded {ln['compat']['why']}, now {o['why']}")

    t0 = time.monotonic()
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.jobs) as ex:
        list(ex.map(one, lines))
    shutil.rmtree(work, ignore_errors=True)
    log(f"replay: {len(lines)} compat build(s) run again in {(time.monotonic() - t0) / 60:.1f} min, "
        f"{len(bad)} behaved differently")
    for b in sorted(bad):
        print("  " + b)
    return 1 if bad else 0


def suite_files(cache, name):
    """The files whose digest pins.json records for suite `name`, relative
    to its directory."""
    pins = core.read_pins()
    root = os.path.join(cache, pins[name]["dir"])
    if name == "test262":
        files = []
        for dirpath, dirnames, filenames in os.walk(os.path.join(root, "harness")):
            files += [os.path.relpath(os.path.join(dirpath, f), root) for f in filenames if f.endswith(".js")]
        files += ["test/" + rel for rel in t262.read_list(t262.SAMPLE)]
        return root, files
    sub = {"closure": pins["closure"]["path"], "closure_library": "closure/goog"}[name]
    files = []
    for dirpath, dirnames, filenames in os.walk(os.path.join(root, sub)):
        files += [os.path.relpath(os.path.join(dirpath, f), root) for f in filenames]
    return root, files


def cmd_verify(args):
    pins = core.read_pins()
    bad = 0
    for name in ("test262", "closure", "closure_library"):
        root, files = suite_files(args.cache, name)
        if not os.path.isdir(root):
            print(f"{name}: {root} is missing (run tools/behave/fetch.sh)")
            bad += 1
            continue
        digest, n = core.tree_digest(root, files)
        want = pins[name].get("sha256")
        if args.update:
            pins[name]["sha256"] = digest
            pins[name]["files"] = n
            print(f"{name}: {n} files, sha256 {digest}")
        elif digest != want:
            print(f"{name}: sha256 {digest} over {n} files, pins.json says {want}")
            bad += 1
        else:
            print(f"{name}: ok ({n} files)")
    if args.update:
        with open(core.PINS, "w") as f:
            json.dump(pins, f, indent=1)
            f.write("\n")
    return 1 if bad else 0


def cmd_report(args):
    return 1 if write_report(os.path.abspath(args.out)) else 0


# ---------------------------------------------------------------------------------------------
# test262 probe and sample


def cmd_probe(args):
    """Runs every statically eligible test262 test uncompiled, and compiles up
    to --per-feature tests per feature (compat, simple-es2015)."""
    out = os.path.abspath(args.out)
    os.makedirs(out, exist_ok=True)
    root = suite_root(args.cache, "test262")
    tests = [t for t in t262.load_tests(root) if t262.static_exclusion(t) is None]
    log(f"probe: {len(tests)} statically eligible tests")
    ctx = Ctx(args, out, [])
    upath = os.path.join(out, "uncompiled.jsonl")
    done = {}
    if os.path.exists(upath):
        with open(upath) as f:
            for ln in f:
                try:
                    r = json.loads(ln)
                    done[r["id"]] = r
                except ValueError:
                    pass
    uf = open(upath, "a")

    def unc(ctx, t):
        c = t262.make_case(root, t)
        work = os.path.join(ctx.work, c.safe_id())
        os.makedirs(work, exist_ok=True)
        prog = os.path.join(work, "u.js")
        with open(prog, "w", encoding="utf-8", errors="surrogateescape") as f:
            f.write(c.uncompiled)
        o = core.outcome(c, None, core.run_program(args.node, prog, c, prog[:-3]))
        shutil.rmtree(work, ignore_errors=True)
        r = {"id": t.rel, "pass": o["pass"], "why": o["why"]}
        with ctx.lock:
            uf.write(json.dumps(r) + "\n")
            uf.flush()
        done[t.rel] = r
        return []

    todo = [t for t in tests if t.rel not in done]

    class Wrap:  # run_parallel wants .suite/.id
        def __init__(self, t):
            self.t, self.suite, self.id = t, "test262", t.rel

    ctx.configs = ["uncompiled"]
    run_parallel(ctx, [Wrap(t) for t in todo], lambda ctx, w: unc(ctx, w.t), "probe uncompiled")
    uf.close()
    # Compile probe: per feature, the first --per-feature passing tests by hash rank.
    by_feature = collections.defaultdict(list)
    for t in tests:
        if done.get(t.rel, {}).get("pass"):
            for ft in t.features:
                by_feature[ft].append(t)
    pick = {}
    for ft, ts in by_feature.items():
        for t in sorted(ts, key=lambda t: t262.rank("probe", t.rel))[:args.per_feature]:
            pick[t.rel] = t
    cpath = os.path.join(out, "compile.jsonl")
    cdone = {}
    if os.path.exists(cpath):
        with open(cpath) as f:
            for ln in f:
                try:
                    r = json.loads(ln)
                    cdone[r["id"]] = r
                except ValueError:
                    pass
    cf = open(cpath, "a")

    def comp(ctx, w):
        c = t262.make_case(root, w.t)
        work = os.path.join(ctx.work, c.safe_id())
        os.makedirs(work, exist_ok=True)
        r = core.compile_case(args.jsgo, c, "simple-es2015", "compat", os.path.join(work, "c.js"))
        shutil.rmtree(work, ignore_errors=True)
        with ctx.lock:
            cf.write(json.dumps({"id": w.t.rel, "ok": r["ok"], "errors": r["errors"]}) + "\n")
            cf.flush()
        return []

    ctx.done, ctx.t0 = 0, time.monotonic()
    ctx.configs = ["simple-es2015"]
    run_parallel(ctx, [Wrap(t) for rel, t in sorted(pick.items()) if rel not in cdone], comp, "probe compile")
    cf.close()
    shutil.rmtree(ctx.work, ignore_errors=True)
    ctx.results.close()
    log(f"probe written to {out}")
    return 0


def cmd_sample(args):
    probe = os.path.abspath(args.probe)
    root = suite_root(args.cache, "test262")
    pins = core.read_pins()
    tests = t262.load_tests(root)
    unc = {}
    with open(os.path.join(probe, "uncompiled.jsonl")) as f:
        for ln in f:
            r = json.loads(ln)
            unc[r["id"]] = r
    comp = {}
    with open(os.path.join(probe, "compile.jsonl")) as f:
        for ln in f:
            r = json.loads(ln)
            comp[r["id"]] = r
    static = collections.Counter()
    eligible = []
    for t in tests:
        e = t262.static_exclusion(t)
        if e:
            static[e] += 1
        else:
            eligible.append(t)
    # Feature exclusions.
    feat_tests = collections.defaultdict(list)
    for t in eligible:
        for ft in t.features:
            feat_tests[ft].append(t)
    node_ex, closure_ex = {}, {}
    for ft, ts in sorted(feat_tests.items()):
        fails = sum(1 for t in ts if not unc.get(t.rel, {}).get("pass"))
        if len(ts) >= 1 and fails / len(ts) >= args.node_threshold:
            node_ex[ft] = {"tests": len(ts), "fail_uncompiled": fails,
                           "reason": f"node {args.node_version} does not implement it: {fails} of {len(ts)} tests "
                                     "fail uncompiled"}
            continue
        probed = [comp[t.rel] for t in ts if t.rel in comp]
        if probed:
            bad = [r for r in probed if not r["ok"]]
            if len(bad) / len(probed) >= args.closure_threshold:
                codes = collections.Counter(c for r in bad for c in r["errors"])
                closure_ex[ft] = {"tests": len(ts), "probed": len(probed), "compile_errors": len(bad),
                                  "codes": dict(codes.most_common(3)),
                                  "reason": f"Closure {pins['closure_version']} rejects it: {len(bad)} of "
                                            f"{len(probed)} probed tests do not compile (simple-es2015)"}
    excluded_features = set(node_ex) | set(closure_ex)
    pool, feat_out, unc_fail = [], collections.Counter(), 0
    for t in eligible:
        hit = sorted(set(t.features) & excluded_features)
        if hit:
            feat_out["node" if any(f in node_ex for f in hit) else "closure"] += 1
            continue
        if not unc.get(t.rel, {}).get("pass"):
            unc_fail += 1
            continue
        pool.append(t)
    sample = t262.stratified_sample(pool, args.size, args.seed, lambda t: t.stratum(), lambda t: t.rel)
    per_change = t262.stratified_sample(sample, args.per_change, args.seed + 1, lambda t: t.stratum(),
                                        lambda t: t.rel)
    header = (f"# test262 {pins['test262']['commit']}, stratified by the first three path components "
              f"under test/ (seed {{seed}}); generated by tools/behave/behave.py sample. {{n}} tests.\n")
    with open(t262.SAMPLE, "w") as f:
        f.write(header.format(seed=args.seed, n=len(sample)))
        f.writelines(t.rel + "\n" for t in sample)
    with open(t262.PER_CHANGE, "w") as f:
        f.write(header.format(seed=args.seed + 1, n=len(per_change)) +
                "# The per-change subset: a stratified sample of sample.txt.\n")
        f.writelines(t.rel + "\n" for t in per_change)
    ex = {"test262_commit": pins["test262"]["commit"], "total_tests": len(tests),
          "static": {k: {"tests": v, "reason": t262.STATIC_REASONS[k]} for k, v in sorted(static.items())},
          "eligible": len(eligible),
          "features_node": node_ex, "features_closure": closure_ex,
          "excluded_by_feature": dict(feat_out),
          "fail_uncompiled": {"tests": unc_fail, "reason": f"the test fails uncompiled on node "
                                                          f"{args.node_version}: no signal"},
          "pool": len(pool), "sample": len(sample), "per_change": len(per_change),
          "seed": args.seed, "node_threshold": args.node_threshold, "closure_threshold": args.closure_threshold}
    with open(t262.EXCLUSIONS, "w") as f:
        json.dump(ex, f, indent=1, sort_keys=True)
        f.write("\n")
    log(f"sample: {len(tests)} tests, {len(eligible)} statically eligible, {len(node_ex)} node and "
        f"{len(closure_ex)} Closure feature exclusions ({sum(feat_out.values())} tests), {unc_fail} fail uncompiled; "
        f"pool {len(pool)}, sample {len(sample)}, per-change {len(per_change)}")
    return 0


# ---------------------------------------------------------------------------------------------


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = ap.add_subparsers(dest="cmd", required=True)

    def common(p, jsgo=True):
        if jsgo:
            p.add_argument("--jsgo", required=True, help="the jsgo binary under test")
        p.add_argument("--node", default=shutil.which("node") or "/opt/node22/bin/node")
        p.add_argument("--cache", default=None, help="the cache directory (default: the main checkout's .cache)")
        p.add_argument("--jobs", type=int, default=os.cpu_count() or 4)
        p.add_argument("--subset", default="per-change",
                       help="test262 subset: per-change (default), full, or a file of test paths")
        p.add_argument("--filter", default="", help="only the cases whose id contains this")
        p.add_argument("--limit", type=int, default=0, help="only the first N cases of each suite")
        p.add_argument("--keep-work", action="store_true", help="keep the work directory")

    p = sub.add_parser("run", help="run suites, compat vs fast")
    common(p)
    p.add_argument("--out", required=True)
    p.add_argument("--suites", default=",".join(SUITES))
    p.add_argument("--configs", default=",".join(core.CONFIGS))
    p.add_argument("--java", action="store_true", help="check uncompiled-pass/compat-fail builds against Java")
    p.add_argument("--compat-cache", default="",
                   help="where compat builds and their runs are cached across runs (default: "
                        "<cache>/behave/compat-cache; 'none': compile and run every compat build)")
    p.add_argument("--fast-flags", default="",
                   help="extra compiler flags for the fast builds only, e.g. "
                        "'--max_optimization_loop_iterations=1' (a test input for the oracle: an output change "
                        "that keeps behaviour)")
    p.add_argument("--rerun-identical", action="store_true",
                   help="run fast builds even when byte-identical to compat (checks the runner's determinism)")
    p.add_argument("--resume", action="store_true", help="keep the results already in --out")
    p.set_defaults(fn=cmd_run)

    p = sub.add_parser("faults", help="seeded-fault check")
    common(p)
    p.add_argument("--baseline", required=True, help="a finished `run` directory")
    p.add_argument("--out", required=True)
    p.add_argument("--faults", default=",".join(FAULTS))
    p.add_argument("--fast-flags", default="", help="extra compiler flags for the faulted fast builds")
    p.add_argument("--keep-diffs", action="store_true")
    p.add_argument("--resume", action="store_true")
    p.set_defaults(fn=cmd_faults)

    p = sub.add_parser("explain", help="where uncaught fault changes are, and whether tests ran them")
    common(p)
    p.add_argument("--baseline", required=True, help="the `run` directory the faults were checked against")
    p.add_argument("--out", required=True, help="the `faults` directory")
    p.add_argument("--faults", default="", help="comma-separated faults (default: all in faults.json)")
    p.set_defaults(fn=cmd_explain)

    p = sub.add_parser("replay", help="run a finished run's compat builds again (determinism)")
    common(p, jsgo=False)
    p.add_argument("--out", required=True, help="a finished `run` directory")
    p.set_defaults(fn=cmd_replay)

    p = sub.add_parser("report", help="rewrite a run's summary")
    p.add_argument("--out", required=True)
    p.set_defaults(fn=cmd_report)

    p = sub.add_parser("verify", help="check the fetched suites against pins.json")
    p.add_argument("--cache", default=None)
    p.add_argument("--update", action="store_true")
    p.set_defaults(fn=cmd_verify)

    p = sub.add_parser("probe", help="test262: uncompiled run and per-feature compile probe")
    common(p)
    p.add_argument("--out", required=True)
    p.add_argument("--per-feature", type=int, default=12)
    p.set_defaults(fn=cmd_probe)

    p = sub.add_parser("sample", help="test262: exclusions, sample and per-change subset from a probe")
    p.add_argument("--probe", required=True)
    p.add_argument("--cache", default=None)
    p.add_argument("--size", type=int, default=5000)
    p.add_argument("--per-change", type=int, default=1500)
    p.add_argument("--seed", type=int, default=20261101)
    p.add_argument("--node-threshold", type=float, default=0.5)
    p.add_argument("--closure-threshold", type=float, default=0.75)
    p.add_argument("--node-version", default="22")
    p.set_defaults(fn=cmd_sample)

    args = ap.parse_args()
    if getattr(args, "cache", "unset") is None:
        args.cache = core.default_cache()
    if getattr(args, "compat_cache", None) == "":
        args.compat_cache = os.path.join(args.cache, "behave", "compat-cache")
    for k in ("jsgo", "node"):
        if getattr(args, k, None) and os.path.sep in getattr(args, k):
            setattr(args, k, os.path.abspath(getattr(args, k)))
    sys.exit(args.fn(args))


if __name__ == "__main__":
    main()
