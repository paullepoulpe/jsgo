#!/usr/bin/env python3
"""End-to-end benchmark of jsgo against the pinned Java Closure Compiler (PLAN.md §6).

  python3 -I tools/bench/bench.py run    [options]   # measure, into --out
  python3 -I tools/bench/bench.py report [options]   # Markdown report from --out
  python3 -I tools/bench/bench.py sample [--check]   # (re)generate tasks/single-files.txt

See tools/bench/README.md. Tasks are tools/bench/tasks/*.json; the pass-to-phase
roll-up is tools/bench/phases.tsv; the warm Java driver is tools/bench/java/CliBatch.java.

Every compile runs in its own directory holding symlinks to the inputs (inputs/, corpus/,
app/, node_modules/), so both compilers see identical relative paths, and outputs land in
that directory. For each task the runner records:
  - Java cold: `java -jar closure-compiler.jar <args>`, N runs (the jar is the cache's
    CLOSURE_JAR, tools/oracle/pin.env);
  - Java warm: CliBatch, one JVM, W warm-up iterations then N measured ones;
  - jsgo: `jsgo <args>`, N runs (interleaved with the Java cold runs);
  - per pass: the same invocations with --tracer_mode=TIMING_ONLY;
and the correctness gate: jsgo's stdout, stderr, exit status and output files (output,
source map) against Java cold run 1, byte for byte. jsgo is MATCH, DIFF, or PENDING (it
reached a not-yet-ported member: a *PendingPort panic, "pending port: ..."), and only a
MATCH task reports speed ratios. Java runs are also checked against each other.
"""
import argparse
import datetime
import glob
import hashlib
import json
import os
import platform
import random
import re
import shutil
import statistics
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
REPO = os.path.dirname(os.path.dirname(HERE))
TASKS_DIR = os.path.join(HERE, "tasks")
PHASES_TSV = os.path.join(HERE, "phases.tsv")
SINGLE_FILES = os.path.join(TASKS_DIR, "single-files.txt")
# The pinned release and its oracle jar (CLOSURE_JAR, a file name in the cache).
PIN_ENV = os.path.join(REPO, "tools", "oracle", "pin.env")
SAMPLE_SEED = 20260415
SAMPLE_SIZE = 1000
SAMPLE_MAX_BYTES = 256 * 1024

PHASES = ["startup", "parse", "checks", "type checking", "transpilation", "optimization loop",
          "renaming", "other passes", "untraced (inside tracker)", "code printing and source maps"]
PENDING_RE = re.compile(r"pending port: (.*?)(?: \[recovered\])*$", re.M)
# A Go runtime crash (as opposed to Java-style diagnostics or an uncaught Java exception).
JAVA_FRAME_RE = re.compile(r"^\t(?:at .*|\.\.\. \d+ more)$")
CRASH_RE = re.compile(r"^(?:panic: |fatal error: )(.*?)(?: \[recovered\])*$", re.M)
# Variables that would change the JVM's or the Go runtime's defaults.
JAVA_ENV_DROP = ("JAVA_TOOL_OPTIONS", "_JAVA_OPTIONS", "JDK_JAVA_OPTIONS")
GO_ENV_DROP = ("GOGC", "GOMEMLIMIT", "GOMAXPROCS", "GODEBUG")


def read_pin():
    """The KEY=value lines of tools/oracle/pin.env (a value may be double-quoted)."""
    pin = {}
    with open(PIN_ENV) as f:
        for line in f:
            line = line.strip()
            if line and not line.startswith("#"):
                k, v = line.split("=", 1)
                pin[k] = v[1:-1] if len(v) >= 2 and v[0] == v[-1] == '"' else v
    return pin


def default_cache():
    if os.environ.get("JSGO_CACHE"):
        return os.environ["JSGO_CACHE"]
    common = subprocess.run(["git", "-C", REPO, "rev-parse", "--path-format=absolute", "--git-common-dir"],
                            capture_output=True, text=True, check=True).stdout.strip()
    return os.path.join(os.path.dirname(common), ".cache")


# ---------------------------------------------------------------------------------------------
# Tasks


def load_tasks(names=None):
    tasks = []
    for fn in sorted(os.listdir(TASKS_DIR)):
        if fn.endswith(".json"):
            with open(os.path.join(TASKS_DIR, fn)) as f:
                t = json.load(f)
            t["file"] = fn
            tasks.append(t)
    tasks.sort(key=lambda t: t.get("order", 99))
    if names:
        want = set(names)
        unknown = want - {t["name"] for t in tasks}
        if unknown:
            sys.exit(f"unknown task(s): {', '.join(sorted(unknown))}")
        tasks = [t for t in tasks if t["name"] in want]
    return tasks


def make_workspace(d, cache):
    """A directory with the input symlinks every task's relative paths refer to."""
    os.makedirs(d, exist_ok=True)
    links = {
        "inputs": os.path.join(cache, "bench-inputs"),
        "corpus": os.path.join(cache, "corpus"),
        "app": os.path.join(HERE, "app"),
    }
    for name, target in links.items():
        p = os.path.join(d, name)
        if not os.path.lexists(p):
            os.symlink(target, p)
    # A node_modules/ layout for module resolution (--module_resolution=NODE).
    nm = os.path.join(d, "node_modules")
    for pkg, version in (("rxjs", "7.8.2"), ("tslib", "2.8.1"), ("core-js", "3.50.0"),
                         ("@angular/core", "16.2.12")):
        p = os.path.join(nm, pkg)
        os.makedirs(os.path.dirname(p), exist_ok=True)
        if not os.path.lexists(p):
            os.symlink(os.path.join(cache, "bench-inputs", f"{pkg}@{version}"), p)


def expand_inputs(specs, ws):
    out = []
    for spec in specs:
        if "file" in spec:
            files = [spec["file"]]
        else:
            files = sorted(glob.glob(spec["glob"], root_dir=ws, recursive=True))
            excl = spec.get("exclude", [])
            files = [f for f in files if not any(glob_match(f, e) for e in excl)]
            if not files:
                sys.exit(f"input glob matched nothing: {spec['glob']} (run tools/bench/fetch-inputs.sh)")
        for f in files:
            if not os.path.exists(os.path.join(ws, f)):
                sys.exit(f"missing input {f} (run tools/bench/fetch-inputs.sh)")
            out += [spec.get("flag", "--js"), f]
    return out


def glob_match(path, pattern):
    import fnmatch
    return fnmatch.fnmatchcase(path, pattern) or fnmatch.fnmatchcase(os.path.basename(path), pattern)


def task_jobs(task, ws):
    """The task's invocations: a list of argument lists (one per compile)."""
    flags, outs = task.get("flags", []), task.get("output_args", [])
    if "per_file" in task:
        pf = task["per_file"]
        with open(os.path.join(TASKS_DIR, pf["list"])) as f:
            names = [l.strip() for l in f if l.strip() and not l.startswith("#")]
        return [flags + ["--js", pf["prefix"] + n] + outs for n in names]
    return [flags + expand_inputs(task["inputs"], ws) + outs]


# ---------------------------------------------------------------------------------------------
# Running and measuring


def loadavg():
    with open("/proc/loadavg") as f:
        return [float(x) for x in f.read().split()[:3]]


def run_proc(argv, cwd, env, out_path, err_path, nice):
    """Runs argv to completion; returns its wall time, CPU time and peak RSS (wait4 rusage)."""
    load = loadavg()
    launch = time.time()
    t0 = time.monotonic()
    with open(out_path, "wb") as o, open(err_path, "wb") as e:
        p = subprocess.Popen(argv, cwd=cwd, env=env, stdin=subprocess.DEVNULL, stdout=o, stderr=e,
                             preexec_fn=(lambda: os.nice(nice)) if nice else None)
        _, status, ru = os.wait4(p.pid, 0)
    wall = time.monotonic() - t0
    return {
        "wall_s": wall, "user_s": ru.ru_utime, "sys_s": ru.ru_stime, "cpu_s": ru.ru_utime + ru.ru_stime,
        "rss_mb": ru.ru_maxrss / 1024, "exit": os.waitstatus_to_exitcode(status),
        "launch_epoch_ms": launch * 1000, "exit_epoch_ms": time.time() * 1000, "load1": load[0],
    }


def sha256_file(p):
    h = hashlib.sha256()
    with open(p, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def artifacts(task, rundir, job_index, single):
    """{artifact name: path} of one invocation's results."""
    if single:
        a = {"stdout": os.path.join(rundir, "stdout"), "stderr": os.path.join(rundir, "stderr"),
             "exit": os.path.join(rundir, "exit")}
        for name in task.get("compare", []):
            a[name] = os.path.join(rundir, name)
        return a
    return {k: os.path.join(rundir, f"{job_index}.{k}") for k in ("stdout", "stderr", "exit")}


def digest_artifacts(arts):
    return {k: (sha256_file(p) if os.path.exists(p) else None) for k, p in arts.items()}


class Runner:
    def __init__(self, args):
        self.args = args
        self.cache = args.cache
        self.jar = os.path.join(self.cache, read_pin()["CLOSURE_JAR"])
        self.out = os.path.abspath(args.out)
        os.makedirs(self.out, exist_ok=True)
        self.java_env = {k: v for k, v in os.environ.items() if k not in JAVA_ENV_DROP}
        self.go_env = {k: v for k, v in os.environ.items() if k not in GO_ENV_DROP}
        for env in (self.java_env, self.go_env):
            env["LC_ALL"] = args.locale
        self.batch_cp = None
        self.jsgo = None

    # -- tools

    def build_batch(self):
        build = os.path.join(self.cache, "bench-build", "java")
        os.makedirs(build, exist_ok=True)
        src = os.path.join(HERE, "java", "CliBatch.java")
        cls = os.path.join(build, "CliBatch.class")
        if not os.path.exists(cls) or os.path.getmtime(cls) < os.path.getmtime(src):
            subprocess.run(["javac", "-nowarn", "-cp", self.jar, "-d", build, src], env=self.java_env, check=True)
        self.batch_cp = self.jar + ":" + build

    def resolve_jsgo(self):
        a = self.args
        info = {}
        if a.jsgo:
            self.jsgo = os.path.abspath(a.jsgo)
            info["source"] = f"binary {self.jsgo}"
        else:
            if a.jsgo_ref:
                src = os.path.join(self.out, "jsgo-src")
                shutil.rmtree(src, ignore_errors=True)
                os.makedirs(src)
                commit = subprocess.run(["git", "-C", REPO, "rev-parse", a.jsgo_ref + "^{commit}"],
                                        capture_output=True, text=True, check=True).stdout.strip()
                arch = subprocess.Popen(["git", "-C", REPO, "archive", commit], stdout=subprocess.PIPE)
                subprocess.run(["tar", "-x", "-C", src], stdin=arch.stdout, check=True)
                arch.wait()
                info["source"] = f"{a.jsgo_ref} ({commit})"
            else:
                src = REPO
                commit = subprocess.run(["git", "-C", REPO, "rev-parse", "HEAD"], capture_output=True,
                                        text=True).stdout.strip()
                dirty = subprocess.run(["git", "-C", REPO, "status", "--porcelain"], capture_output=True,
                                       text=True).stdout.strip()
                info["source"] = f"this checkout ({commit}{', dirty' if dirty else ''})"
            if not os.path.isdir(os.path.join(src, "cmd", "jsgo")):
                info["status"] = "UNAVAILABLE (no cmd/jsgo in " + info["source"] + ")"
                return info
            binary = os.path.join(self.out, "jsgo")
            r = subprocess.run(["go", "build", "-o", binary, "./cmd/jsgo"], cwd=src, capture_output=True, text=True)
            if r.returncode != 0:
                info["status"] = "BUILD FAILED"
                info["build_log"] = r.stderr[-4000:]
                return info
            self.jsgo = binary
        info["sha256"] = sha256_file(self.jsgo)
        return info

    def java_cold_argv(self, args):
        return ["java", "-jar", self.jar] + args

    # -- one task

    def run_task(self, task):
        a = self.args
        name = task["name"]
        tdir = os.path.join(self.out, name)
        os.makedirs(tdir, exist_ok=True)
        ws = os.path.join(tdir, "ws")
        make_workspace(ws, self.cache)
        jobs = task_jobs(task, ws)
        single = len(jobs) == 1
        with open(os.path.join(tdir, "jobs.json"), "w") as f:
            json.dump(jobs, f)
        runs = a.runs or task.get("runs", 5)
        print(f"== {name}: {len(jobs)} invocation(s), {runs} run(s)", flush=True)

        do_java = not a.skip_java
        do_jsgo = self.jsgo is not None and not a.skip_jsgo

        # Java cold and jsgo runs are interleaved, so that a change in the machine's load
        # affects both alike. jsgo's first run is gated against Java cold run 1; a PENDING
        # jsgo is not timed further.
        cold, jsgo = [], []
        jsgo_status = None
        if do_java:
            shutil.rmtree(os.path.join(tdir, "java-cold"), ignore_errors=True)
        if do_jsgo:
            shutil.rmtree(os.path.join(tdir, "jsgo"), ignore_errors=True)
        for i in range(1, runs + 1):
            if do_java:
                cold.append(self.measure_cold(task, jobs, tdir, "java-cold", i, self.java_cold_argv, self.java_env))
            if do_jsgo and jsgo_status not in ("PENDING", "CRASH"):
                r = self.measure_cold(task, jobs, tdir, "jsgo", i, lambda args: [self.jsgo] + args, self.go_env)
                jsgo.append(r)
                if i == 1:
                    jsgo_status = self.gate(task, jobs, tdir, r)
        if do_java:
            self.check_java_runs(task, jobs, tdir, cold)
            save(os.path.join(tdir, "java-cold.json"), {"runs": cold})
        if do_jsgo:
            save(os.path.join(tdir, "jsgo.json"), {"runs": jsgo})

        if do_java and not a.skip_warm:
            warm = self.measure_warm(task, jobs, tdir)
            save(os.path.join(tdir, "java-warm.json"), warm)

        if not a.skip_trace:
            self.trace(task, jobs, tdir, do_java, do_jsgo and jsgo_status not in ("PENDING", "CRASH"))

    def measure_cold(self, task, jobs, tdir, who, i, argv_of, env):
        """One run of the task: every invocation once, as separate processes."""
        rundir = os.path.join(tdir, who, f"run-{i}")
        shutil.rmtree(rundir, ignore_errors=True)
        make_workspace(rundir, self.cache)
        single = len(jobs) == 1
        total = {"wall_s": 0.0, "user_s": 0.0, "sys_s": 0.0, "cpu_s": 0.0, "rss_mb": 0.0}
        per_job, loads, digests = [], [], []
        for j, args in enumerate(jobs):
            arts = artifacts(task, rundir, j, single)
            for k in task.get("compare", []) if single else []:
                if os.path.exists(arts[k]):
                    os.remove(arts[k])
            r = run_proc(argv_of(args), rundir, env, arts["stdout"], arts["stderr"], self.args.nice)
            with open(arts["exit"], "w") as f:
                f.write(f"{r['exit']}\n")
            for k in ("wall_s", "user_s", "sys_s", "cpu_s"):
                total[k] += r[k]
            total["rss_mb"] = max(total["rss_mb"], r["rss_mb"])
            loads.append(r["load1"])
            per_job.append({k: r[k] for k in ("wall_s", "cpu_s", "rss_mb", "exit")})
            digests.append(digest_artifacts(arts))
        total["load1"] = statistics.median(loads)
        total["exits"] = sorted({p["exit"] for p in per_job})
        if not single:
            total["per_job"] = per_job
        total["digests"] = digests
        print(f"   {who} run {i}: wall {total['wall_s']:.2f}s cpu {total['cpu_s']:.2f}s "
              f"rss {total['rss_mb']:.0f}MB load {total['load1']:.2f} exit {total['exits']}", flush=True)
        # Keep the full outputs of the first run only (Java's is the gate's reference).
        if i > 1:
            for k, p in ((k, p) for j in range(len(jobs)) for k, p in artifacts(task, rundir, j, single).items()):
                if k not in ("exit",) and os.path.exists(p) and os.path.getsize(p) > 4096:
                    os.remove(p)
        return total

    def gate(self, task, jobs, tdir, jsgo_run):
        return gate(task, jobs, tdir)

    def check_java_runs(self, task, jobs, tdir, cold):
        """Java must agree with itself: every cold run with run 1 (by digest)."""
        if not cold:
            return
        ref = cold[0]["digests"]
        bad = [i + 2 for i, r in enumerate(cold[1:]) if r["digests"] != ref]
        save(os.path.join(tdir, "java-consistency.json"), {"cold_runs_differing_from_run_1": bad})
        if bad:
            print(f"   WARNING: Java cold runs {bad} differ from run 1", flush=True)

    def measure_warm(self, task, jobs, tdir):
        a = self.args
        warmup = a.warmup if a.warmup is not None else task.get("warm_warmup", 2)
        iters = a.warm_iters or task.get("warm_iters", 5)
        rundir = os.path.join(tdir, "java-warm")
        shutil.rmtree(rundir, ignore_errors=True)
        make_workspace(rundir, self.cache)
        jobfile = os.path.join(rundir, "jobs.tsv")
        with open(jobfile, "w") as f:
            for args in jobs:
                f.write("\t".join(args) + "\n")
        argv = ["java", "-cp", self.batch_cp, "CliBatch", str(warmup), str(iters), jobfile, "batch"]
        r = run_proc(argv, rundir, self.java_env, os.path.join(rundir, "report.tsv"),
                     os.path.join(rundir, "driver.stderr"), a.nice)
        iters_out, peak = parse_batch_report(os.path.join(rundir, "report.tsv"))
        # The driver's last iteration must reproduce Java cold run 1.
        ref_dir = os.path.join(tdir, "java-cold", "run-1")
        single = len(jobs) == 1
        mismatch = []
        if os.path.isdir(ref_dir):
            for j in range(len(jobs)):
                ref = artifacts(task, ref_dir, j, single)
                got = {k: os.path.join(rundir, "batch", f"{j}.{k}") for k in ("stdout", "stderr", "exit")}
                for k in task.get("compare", []) if single else []:
                    got[k] = os.path.join(rundir, k)
                bad = [k for k in ref if os.path.exists(ref[k]) and not same_file(ref[k], got[k])]
                if bad:
                    mismatch.append({"job": j, "artifacts": bad})
        measured = [i for i in iters_out if i["kind"] == "measure"]
        print(f"   java-warm: {len(measured)} iteration(s), median wall "
              f"{statistics.median([i['wall_ms'] for i in measured]) / 1000:.2f}s, peak rss {peak:.0f}MB, "
              f"exit {r['exit']}, {'reproduces cold run 1' if not mismatch else f'DIFFERS from cold run 1: {mismatch[:3]}'}",
              flush=True)
        return {"warmup": warmup, "iters": iters, "iterations": iters_out, "peak_rss_mb": peak,
                "process": {k: r[k] for k in ("wall_s", "cpu_s", "rss_mb", "exit", "load1")},
                "load1_end": loadavg()[0], "mismatch_vs_cold": mismatch}

    # -- per pass

    def trace(self, task, jobs, tdir, do_java, do_jsgo):
        a = self.args
        k = a.trace_runs or task.get("trace_runs", 3)
        tr_jobs = [args + ["--tracer_mode=TIMING_ONLY"] for args in jobs]
        single = len(jobs) == 1
        result = {}
        if do_java and task.get("trace_cold", True):
            result["java-cold"] = [self.trace_cold(tr_jobs, tdir, "java-cold", i, self.java_cold_argv, self.java_env)
                                   for i in range(1, k + 1)]
        if do_java and not a.skip_warm:
            result["java-warm"] = self.trace_warm(task, tr_jobs, tdir, k)
        if do_jsgo:
            result["jsgo"] = [self.trace_cold(tr_jobs, tdir, "jsgo", i, lambda args: [self.jsgo] + args, self.go_env)
                              for i in range(1, k + 1)]
        path = os.path.join(tdir, "trace.json")
        old = load(path) or {}
        old.update(result)
        save(path, old)

    def trace_cold(self, jobs, tdir, who, i, argv_of, env):
        rundir = os.path.join(tdir, "trace", f"{who}-{i}")
        shutil.rmtree(rundir, ignore_errors=True)
        make_workspace(rundir, self.cache)
        per_job = []
        for j, args in enumerate(jobs):
            out, err = os.path.join(rundir, "stdout"), os.path.join(rundir, f"{j}.stderr")
            r = run_proc(argv_of(args), rundir, env, out, err, self.args.nice)
            rep = parse_tracer(read_text(err))
            per_job.append(phase_times(rep, r["launch_epoch_ms"], r["exit_epoch_ms"]))
            if j > 0 or len(jobs) > 1:
                os.remove(err)
        for f in os.listdir(rundir):
            p = os.path.join(rundir, f)
            if os.path.isfile(p) and not os.path.islink(p) and f != "0.stderr":
                os.remove(p)
        return merge_jobs(per_job)

    def trace_warm(self, task, jobs, tdir, k):
        a = self.args
        warmup = a.warmup if a.warmup is not None else task.get("warm_warmup", 2)
        rundir = os.path.join(tdir, "trace", "java-warm")
        shutil.rmtree(rundir, ignore_errors=True)
        make_workspace(rundir, self.cache)
        jobfile = os.path.join(rundir, "jobs.tsv")
        with open(jobfile, "w") as f:
            for args in jobs:
                f.write("\t".join(args) + "\n")
        argv = ["java", "-cp", self.batch_cp, "CliBatch", "--keep-stderr", str(warmup), str(k), jobfile, "batch"]
        run_proc(argv, rundir, self.java_env, os.path.join(rundir, "report.tsv"),
                 os.path.join(rundir, "driver.stderr"), a.nice)
        iters, _ = parse_batch_report(os.path.join(rundir, "report.tsv"))
        out = []
        for it in [i for i in iters if i["kind"] == "measure"]:
            n = it["iter"]
            reps = []
            for j in range(len(jobs)):
                p = os.path.join(rundir, "batch", f"{j}.{n}.stderr")
                reps.append(parse_tracer(read_text(p)))
                os.remove(p)
            if len(jobs) == 1:
                # The iteration's bounds play the part of process start and exit.
                m = merge_jobs([phase_times(reps[0], it["start_epoch_ms"], it["end_epoch_ms"])])
            else:
                # Inside one JVM, startup and printing cannot be told apart per invocation: the
                # time outside all of the iteration's trackers goes on the printing row.
                m = merge_jobs([phase_times(r, None, None) for r in reps])
                inside = sum(m["phases"].values())
                m["phases"]["code printing and source maps"] = max(0.0, it["wall_ms"] - inside)
            m["wall_ms"] = it["wall_ms"]
            out.append(m)
        return out


def without_java_frames(text):
    """stderr without the stack frames of a Java exception trace ("\tat ...", "\t... N more")."""
    return "\n".join(l for l in text.split("\n") if not JAVA_FRAME_RE.match(l))


def gate(task, jobs, tdir):
    """Compares jsgo's first run with Java cold run 1, invocation by invocation, byte for byte.

    Writes gate.json: the status (MATCH, DIFF or PENDING), counts, the pending signatures, Go
    crashes, and for DIFF invocations their kind and first differences:
      crash        jsgo died with a Go panic or fatal error;
      java-frames  Java threw an internal exception and only its stack frames differ (jsgo prints
                   the same exception and message, exit status and output, but not Java's frames);
      output       anything else.
    Returns the status, or CRASH when jsgo crashes on every invocation (it is then not timed)."""
    single = len(jobs) == 1
    ref_dir = os.path.join(tdir, "java-cold", "run-1")
    got_dir = os.path.join(tdir, "jsgo", "run-1")
    if not os.path.isdir(ref_dir):
        save(os.path.join(tdir, "gate.json"), {"status": "NO-REFERENCE"})
        return None
    counts = {"MATCH": 0, "DIFF": 0, "PENDING": 0}
    kinds = {"crash": 0, "java-frames": 0, "output": 0}
    pending_sigs, diffs, crashes, matched = {}, [], {}, []
    per_kind_kept = {}
    for j in range(len(jobs)):
        ref = artifacts(task, ref_dir, j, single)
        got = artifacts(task, got_dir, j, single)
        err = read_text(got["stderr"])
        m = PENDING_RE.search(err)
        if m and "panic" in err:
            counts["PENDING"] += 1
            pending_sigs[m.group(1)] = pending_sigs.get(m.group(1), 0) + 1
            continue
        bad = [k for k in ref if not same_file(ref[k], got[k])]
        if not bad:
            counts["MATCH"] += 1
            matched.append(j)
            continue
        counts["DIFF"] += 1
        c = CRASH_RE.search(err)
        if c:
            kind = "crash"
            crashes[c.group(1)] = crashes.get(c.group(1), 0) + 1
        elif bad == ["stderr"] and without_java_frames(read_text(ref["stderr"])) == without_java_frames(err):
            kind = "java-frames"
        else:
            kind = "output"
        kinds[kind] += 1
        if per_kind_kept.get(kind, 0) < 5:
            per_kind_kept[kind] = per_kind_kept.get(kind, 0) + 1
            diffs.append({"job": j, "input": jobs[j][-1] if not single else None, "kind": kind, "artifacts": bad,
                          "detail": {k: describe_diff(ref[k], got[k]) for k in bad}})
    status = "MATCH" if counts["MATCH"] == len(jobs) else ("PENDING" if counts["PENDING"] else "DIFF")
    info = {"status": status, "counts": counts, "diff_kinds": kinds, "pending": pending_sigs,
            "crashes": crashes, "diffs": diffs, "matched": matched if not single else None}
    save(os.path.join(tdir, "gate.json"), info)
    print(f"   gate: {status} {counts}" + (f" {next(iter(pending_sigs))}" if pending_sigs else "")
          + (f" diff kinds: {kinds}" if counts["DIFF"] else ""), flush=True)
    if status == "DIFF" and kinds["crash"] == len(jobs):
        return "CRASH"
    return status


def parse_batch_report(path):
    iters, peak = [], 0.0
    with open(path) as f:
        for line in f:
            line = line.rstrip("\n")
            if line.startswith("# peak_rss_mb="):
                peak = float(line.split("=", 1)[1])
            elif line and not line.startswith("iter") and not line.startswith("#"):
                c = line.split("\t")
                iters.append({"iter": int(c[0]), "kind": c[1], "wall_ms": float(c[2]), "cpu_ms": float(c[3]),
                              "gc_count": int(c[4]), "gc_ms": int(c[5]), "rss_mb": float(c[6]),
                              "start_epoch_ms": int(c[7]), "end_epoch_ms": int(c[8])})
    return iters, peak


# ---------------------------------------------------------------------------------------------
# Tracer reports and phases


def load_phase_map():
    m = {}
    with open(PHASES_TSV) as f:
        for line in f:
            if line.startswith("#") or not line.strip():
                continue
            p, ph = line.rstrip("\n").split("\t")
            m[p] = ph
    return m


PHASE_MAP = None


def parse_tracer(text):
    """The TIMING_ONLY report: tracker start/end (epoch ms) and per-pass runtime/runs."""
    rep = {"passes": {}, "start": None, "end": None, "wall": None}
    lines = text.split("\n")
    i = 0
    while i < len(lines):
        l = lines[i]
        if l.startswith("Start time(ms): "):
            rep["start"] = int(l.split(": ")[1])
        elif l.startswith("End time(ms): "):
            rep["end"] = int(l.split(": ")[1])
        elif l.startswith("Wall time(ms): "):
            rep["wall"] = int(l.split(": ")[1])
        elif l == "Summary:" and i + 1 < len(lines) and lines[i + 1].startswith("pass,runtime"):
            i += 2
            while i < len(lines) and lines[i].strip():
                parts = lines[i].rsplit(",", 7)
                if len(parts) == 8:
                    rep["passes"][parts[0]] = {"ms": int(parts[1]), "runs": int(parts[3])}
                i += 1
            continue
        i += 1
    return rep


def phase_times(rep, launch_ms, exit_ms):
    global PHASE_MAP
    if PHASE_MAP is None:
        PHASE_MAP = load_phase_map()
    phases = {ph: 0.0 for ph in PHASES}
    unmapped = []
    for p, s in rep["passes"].items():
        ph = PHASE_MAP.get(p)
        if ph is None:
            unmapped.append(p)
            ph = "other passes"
        phases[ph] += s["ms"]
    passes_sum = sum(s["ms"] for s in rep["passes"].values())
    ok = rep["start"] is not None and rep["end"] is not None
    if ok:
        phases["untraced (inside tracker)"] = max(0.0, (rep["end"] - rep["start"]) - passes_sum)
        if launch_ms is not None:
            phases["startup"] = max(0.0, rep["start"] - launch_ms)
            phases["code printing and source maps"] = max(0.0, exit_ms - rep["end"])
    return {"phases": phases, "passes": rep["passes"], "unmapped": unmapped, "has_report": ok,
            "wall_ms": (exit_ms - launch_ms) if launch_ms is not None else None}


def merge_jobs(per_job):
    """Sums one run's per-invocation phase and pass times."""
    phases = {ph: sum(j["phases"][ph] for j in per_job) for ph in PHASES}
    passes = {}
    for j in per_job:
        for p, s in j["passes"].items():
            t = passes.setdefault(p, {"ms": 0, "runs": 0})
            t["ms"] += s["ms"]
            t["runs"] += s["runs"]
    walls = [j["wall_ms"] for j in per_job if j["wall_ms"] is not None]
    return {"phases": phases, "passes": passes,
            "unmapped": sorted({p for j in per_job for p in j["unmapped"]}),
            "reports": sum(1 for j in per_job if j["has_report"]), "jobs": len(per_job),
            "wall_ms": sum(walls) if walls else None}


# ---------------------------------------------------------------------------------------------
# Helpers


def read_text(p):
    try:
        with open(p, "rb") as f:
            return f.read().decode("utf-8", "replace")
    except FileNotFoundError:
        return ""


def same_file(a, b):
    ea, eb = os.path.exists(a), os.path.exists(b)
    if ea != eb:
        return False
    if not ea:
        return True
    if os.path.getsize(a) != os.path.getsize(b):
        return False
    return sha256_file(a) == sha256_file(b)


def describe_diff(a, b):
    if not os.path.exists(b):
        return "missing in jsgo"
    if not os.path.exists(a):
        return "only in jsgo"
    with open(a, "rb") as f:
        da = f.read()
    with open(b, "rb") as f:
        db = f.read()
    n = next((i for i in range(min(len(da), len(db))) if da[i] != db[i]), min(len(da), len(db)))
    line = da[:n].count(b"\n") + 1
    return (f"java {len(da)} bytes, jsgo {len(db)} bytes; first difference at byte {n} (line {line}): "
            f"java {da[n:n + 60]!r} / jsgo {db[n:n + 60]!r}")


def save(path, obj):
    tmp = path + ".tmp"
    with open(tmp, "w") as f:
        json.dump(obj, f, indent=1)
    os.replace(tmp, path)


def load(path):
    try:
        with open(path) as f:
            return json.load(f)
    except FileNotFoundError:
        return None


def machine_info(runner, jsgo_info):
    def cmd(*argv):
        try:
            r = subprocess.run(argv, capture_output=True, text=True, env=runner.java_env)
            return (r.stdout + r.stderr).strip()
        except OSError as e:
            return str(e)
    cpu = ""
    with open("/proc/cpuinfo") as f:
        for line in f:
            if line.startswith("model name"):
                cpu = line.split(":", 1)[1].strip()
                break
    mem = ""
    with open("/proc/meminfo") as f:
        for line in f:
            if line.startswith("MemTotal"):
                mem = f"{int(line.split()[1]) / 1048576:.1f} GiB"
    osname = ""
    if os.path.exists("/etc/os-release"):
        with open("/etc/os-release") as f:
            for line in f:
                if line.startswith("PRETTY_NAME="):
                    osname = line.split("=", 1)[1].strip().strip('"')
    return {
        "date": datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M UTC"),
        "host": platform.node(), "cpu": cpu, "cpus": os.cpu_count(), "memory": mem,
        "kernel": platform.release(), "os": osname,
        "java": cmd("java", "-version").splitlines()[:3],
        "go": cmd("go", "version"),
        "jar_sha256": sha256_file(runner.jar),
        "jsgo": jsgo_info,
        "nice": runner.args.nice, "locale": runner.args.locale,
        "loadavg_start": loadavg(),
        "bench_commit": subprocess.run(["git", "-C", REPO, "rev-parse", "HEAD"], capture_output=True,
                                       text=True).stdout.strip(),
    }


# ---------------------------------------------------------------------------------------------
# Report


def fmt_s(v):
    if v is None:
        return "–"
    return f"{v:.2f}" if v < 10 else (f"{v:.1f}" if v < 100 else f"{v:.0f}")


def stat(values):
    values = [v for v in values if v is not None]
    if not values:
        return None
    return {"median": statistics.median(values), "min": min(values), "max": max(values), "n": len(values)}


def fmt_stat(s, unit=""):
    if not s:
        return "–"
    if s["n"] == 1:
        return f"{fmt_s(s['median'])}{unit}"
    return f"{fmt_s(s['median'])}{unit} ({fmt_s(s['min'])}–{fmt_s(s['max'])})"


def ratio(a, b):
    if not a or not b or not b["median"]:
        return "–"
    return f"{a['median'] / b['median']:.2f}×"


def report(out_dir, md_path, preliminary=True, notes=None):
    machine = load(os.path.join(out_dir, "machine.json")) or {}
    tasks = load_tasks()
    rows, sections = [], []
    for t in tasks:
        tdir = os.path.join(out_dir, t["name"])
        if not os.path.isdir(tdir):
            continue
        cold = load(os.path.join(tdir, "java-cold.json"))
        warm = load(os.path.join(tdir, "java-warm.json"))
        jsgo = load(os.path.join(tdir, "jsgo.json"))
        gate = load(os.path.join(tdir, "gate.json"))
        consistency = load(os.path.join(tdir, "java-consistency.json"))
        trace = load(os.path.join(tdir, "trace.json")) or {}
        jobs = load(os.path.join(tdir, "jobs.json")) or []

        c_wall = stat([r["wall_s"] for r in cold["runs"]]) if cold else None
        c_cpu = stat([r["cpu_s"] for r in cold["runs"]]) if cold else None
        c_rss = stat([r["rss_mb"] for r in cold["runs"]]) if cold else None
        c_load = stat([r["load1"] for r in cold["runs"]]) if cold else None
        w_meas = [i for i in warm["iterations"] if i["kind"] == "measure"] if warm else []
        w_wall = stat([i["wall_ms"] / 1000 for i in w_meas])
        w_cpu = stat([i["cpu_ms"] / 1000 for i in w_meas])
        w_rss = warm["peak_rss_mb"] if warm else None
        status = (gate or {}).get("status") or ("not run" if not jsgo else "–")
        crashes = gate_crashes(gate, tdir, len(jobs)) if gate else {}
        crashed_all = status == "DIFF" and sum(crashes.values()) == len(jobs)
        g_wall = g_cpu = g_rss = None
        if jsgo and jsgo.get("runs") and status in ("MATCH", "DIFF") and not crashed_all:
            g_wall = stat([r["wall_s"] for r in jsgo["runs"]])
            g_cpu = stat([r["cpu_s"] for r in jsgo["runs"]])
            g_rss = stat([r["rss_mb"] for r in jsgo["runs"]])
        status_cell = status + (" (jsgo crashes)" if crashed_all else "")
        if gate and gate.get("counts") and len(jobs) > 1:
            c = gate["counts"]
            status_cell += f" ({c['MATCH']}/{c['DIFF']}/{c['PENDING']} match/diff/pending)"
        speed_cold = ratio(c_wall, g_wall) if status == "MATCH" else "–"
        speed_warm = ratio(w_wall, g_wall) if status == "MATCH" else "–"
        subset = matching_subset(gate, cold, jsgo) if len(jobs) > 1 and status != "MATCH" else None
        if subset:
            speed_cold = f"{subset['java'] / subset['jsgo']:.2f}× (on the {subset['n']} matching)"
        rows.append(
            f"| {t['title']} | {fmt_stat(c_wall, ' s')} | {fmt_stat(c_cpu, ' s')} | {fmt_s(c_rss and c_rss['median'])} MB "
            f"| {fmt_stat(w_wall, ' s')} | {fmt_stat(w_cpu, ' s')} | {fmt_s(w_rss)} MB "
            f"| {status_cell} | {fmt_stat(g_wall, ' s')} | {fmt_stat(g_cpu, ' s')} "
            f"| {fmt_s(g_rss and g_rss['median'])}{' MB' if g_rss else ''} | {speed_cold} | {speed_warm} "
            f"| {fmt_stat(c_load)} |")

        # Section per task.
        s = [f"### {t['title']} (`{t['name']}`)", "", t.get("description", ""), ""]
        s.append(f"Mode: {t.get('mode', '')}. Invocations per run: {len(jobs)}. "
                 f"Java cold runs: {c_wall['n'] if c_wall else 0}; Java warm: "
                 f"{warm['warmup'] if warm else 0} warm-up + {len(w_meas)} measured iteration(s) in one JVM; "
                 f"task file `tools/bench/tasks/{t['file']}`.")
        s.append("")
        if cold:
            exits = sorted({e for r in cold["runs"] for e in r["exits"]})
            s.append(f"- Java exit status: {exits}.")
            if len(jobs) > 1:
                pj = [p["wall_s"] for p in cold["runs"][0]["per_job"]]
                q = statistics.quantiles(pj, n=10)
                s.append(f"- Java cold, per invocation (run 1): median {statistics.median(pj):.2f} s, "
                         f"p10 {q[0]:.2f} s, p90 {q[-1]:.2f} s.")
        if consistency is not None:
            bad = consistency["cold_runs_differing_from_run_1"]
            s.append("- Java determinism: " + ("all cold runs byte-identical." if not bad
                                               else f"cold runs {bad} differ from run 1."))
        if warm:
            mm = warm.get("mismatch_vs_cold")
            s.append("- Warm driver fidelity: " + ("its last iteration reproduces Java cold run 1 byte for byte."
                                                    if not mm else f"DIFFERS from cold run 1: {mm[:5]}"))
            s.append(f"- Java warm iterations (wall s): "
                     + ", ".join(f"{i['wall_ms'] / 1000:.2f}{'*' if i['kind'] == 'warmup' else ''}"
                                 for i in warm["iterations"]) + " (* = warm-up).")
        if gate:
            if gate.get("status") == "PENDING":
                sigs = sorted(gate["pending"].items(), key=lambda kv: -kv[1])
                s.append("- jsgo: PENDING, it reaches pending ports: "
                         + "; ".join(f"`{k}`" + (f" ({v}×)" if len(jobs) > 1 else "") for k, v in sigs[:5]) + ".")
            elif gate.get("status") == "DIFF":
                s.append("- jsgo: DIFF." + (" It is not timed: it crashes." if crashed_all else ""))
            if crashes:
                s.append("- jsgo crashes (Go runtime): " + "; ".join(
                    f"`{k}`" + (f" ({v}×)" if len(jobs) > 1 else "") for k, v in crashes.items()) + ".")
                if len(jobs) == 1:
                    frames = pending_frames(read_text(os.path.join(tdir, "jsgo", "run-1", "stderr")))
                    if frames:
                        s.append("  Its stacks also hold a pending port being unwound: "
                                 + ", ".join(f"`{f}`" for f in frames) + " (pending_gen.go).")
            if subset:
                s.append(f"- On the {subset['n']} invocations that match byte for byte: Java cold {subset['java']:.0f} s, "
                         f"jsgo {subset['jsgo']:.0f} s in total (per invocation, median over runs), a "
                         f"{subset['java'] / subset['jsgo']:.2f}× speed-up; median per invocation Java "
                         f"{subset['java_med']:.2f} s, jsgo {subset['jsgo_med']:.2f} s.")
            if gate.get("diff_kinds") and len(jobs) > 1:
                s.append("- DIFF kinds: " + ", ".join(f"{v} {KIND_LABEL[k]}" for k, v in gate["diff_kinds"].items() if v)
                         + ". Reproducers: one invocation each, the task's flags plus `--js <input>` (paths relative "
                         "to the run directory, `corpus/` = `.cache/corpus/`).")
            if gate.get("diffs"):
                s.append("- jsgo differs from Java" + (" (up to 5 invocations per kind)" if len(jobs) > 1 else "") + ":")
                for d in gate["diffs"][:15]:
                    where = f"`{jobs[d['job']][-1]}` ({d.get('kind', '?')}), " if len(jobs) > 1 else ""
                    for k, v in d["detail"].items():
                        if len(jobs) > 1 and d.get("kind") == "java-frames":
                            s.append(f"  - {where}{k}: only Java's stack frames differ")
                            break
                        s.append(f"  - {where}{k}: {v}")
                        if len(jobs) > 1 and d.get("kind") == "crash":
                            break
            if gate.get("status") == "MATCH":
                s.append("- jsgo: MATCH (stdout, stderr, exit status and output files byte-identical).")
        s.append("")
        s += phase_table(trace, status == "MATCH")
        sections.append("\n".join(s))

    lines = ["# jsgo vs Java Closure Compiler: end-to-end benchmark", ""]
    if preliminary:
        lines += ["> **Preliminary.** Measured on a shared machine with other jobs running (see the load "
                  "averages), under `nice -n %s`. A quiet, pinned re-run will replace these numbers." % machine.get("nice", 10), ""]
    if notes:
        with open(notes) as f:
            lines += [f.read().rstrip("\n"), ""]
    lines += ["Spec: PLAN.md §6. Harness: `tools/bench/` (README.md there). Values are medians, with the "
              "min–max range over the runs in parentheses. Wall and CPU are seconds; CPU is user+system of the "
              "process (for Java warm: the JVM's process CPU during the iteration, JIT and GC threads included). "
              "RSS is peak resident memory (Java warm: the JVM's peak over its whole life, warm-up included). "
              "Speed-ups are Java time ÷ jsgo time (above 1× means jsgo is faster) and only appear when jsgo's output, "
              "source map and diagnostics match Java's byte for byte. *load* is the 1-minute load average at the start "
              "of each Java cold run.", "",
              "## Overall", "",
              "| Task | Java cold wall | Java cold CPU | Java cold RSS | Java warm wall | Java warm CPU | Java warm RSS "
              "| jsgo | jsgo wall | jsgo CPU | jsgo RSS | jsgo speed-up vs cold | vs warm | load (1 min) |",
              "|---|---|---|---|---|---|---|---|---|---|---|---|---|---|"] + rows + [""]
    lines += ["## Per task and per phase", "",
              "Per-phase times come from separate runs with `--tracer_mode=TIMING_ONLY` (milliseconds, median "
              "over the trace runs; for multi-invocation tasks, summed over the invocations of a run). Passes are "
              "rolled up with `tools/bench/phases.tsv`. *startup* is process start to the tracker's start; "
              "*code printing and source maps* is the tracker's report to process exit (the report is printed "
              "before the output is generated, so this row holds toSource, the source map, writing the output "
              "and exit). For Java warm, the iteration's start and end stand for process start and exit; with several "
              "invocations per iteration (the single-file task) startup and printing cannot be told apart inside the "
              "JVM, so all time outside the trackers is on the printing row and startup is 0. *untraced* is the "
              "tracker's wall time not attributed to any pass. Speed-ups (Java ÷ jsgo) only for MATCH tasks.", ""]
    lines += sections
    lines += ["## Machine", ""]
    for k in ("date", "cpu", "cpus", "memory", "kernel", "os", "go", "jar_sha256", "nice", "locale",
              "loadavg_start", "loadavg_end", "bench_commit", "jsgo_measured"):
        if k in machine:
            lines.append(f"- {k}: {machine[k]}")
    if machine.get("java"):
        lines.append(f"- java: {' / '.join(machine['java'])}")
    if machine.get("jsgo"):
        lines.append("- jsgo: " + ", ".join(f"{k} {v}" for k, v in machine["jsgo"].items() if k != "build_log"))
    lines.append("- environment: JAVA_TOOL_OPTIONS, _JAVA_OPTIONS and JDK_JAVA_OPTIONS removed for Java (default "
                 "JVM settings); GOGC, GOMEMLIMIT, GOMAXPROCS and GODEBUG removed for jsgo (default Go GC).")
    lines.append("")
    text = "\n".join(lines)
    if md_path:
        os.makedirs(os.path.dirname(os.path.abspath(md_path)), exist_ok=True)
        with open(md_path, "w") as f:
            f.write(text)
    return text


def pending_frames(text):
    """Pending-port declarations (pending_gen.go) on the stacks of a crash dump."""
    lines = text.split("\n")
    out = []
    for i in range(1, len(lines)):
        if "/pending_gen.go:" in lines[i]:
            fn = lines[i - 1].rsplit("/", 1)[-1].split("(")[0]
            if fn not in out:
                out.append(fn)
    return out


def matching_subset(gate, cold, jsgo):
    """Java cold and jsgo wall time summed over the invocations that match (multi-invocation tasks)."""
    m = (gate or {}).get("matched")
    if not m or not cold or not jsgo or not jsgo.get("runs"):
        return None
    def per_job(runs, j):
        return statistics.median(r["per_job"][j]["wall_s"] for r in runs if "per_job" in r)
    try:
        java = [per_job(cold["runs"], j) for j in m]
        go = [per_job(jsgo["runs"], j) for j in m]
    except (KeyError, IndexError, statistics.StatisticsError):
        return None
    return {"n": len(m), "java": sum(java), "jsgo": sum(go), "java_med": statistics.median(java),
            "jsgo_med": statistics.median(go)}


KIND_LABEL = {"crash": "Go crashes", "java-frames": "Java internal exceptions whose stack frames jsgo does not print",
              "output": "output/diagnostic differences"}


def gate_crashes(gate, tdir, njobs):
    """The gate's Go runtime crashes; for results gated before they were recorded, from jsgo's stderr."""
    if "crashes" in gate:
        return gate["crashes"]
    out = {}
    if njobs == 1 and gate.get("status") == "DIFF":
        c = CRASH_RE.search(read_text(os.path.join(tdir, "jsgo", "run-1", "stderr")))
        if c:
            out[c.group(1)] = 1
    return out


def phase_table(trace, show_ratio):
    cols = [c for c in ("java-cold", "java-warm", "jsgo") if trace.get(c)]
    missing = [c for c in cols if not any(r.get("reports", 1) for r in trace[c])]
    cols = [c for c in cols if c not in missing]
    if not cols:
        return ["(no per-phase data)", ""]
    med = {}
    for c in cols:
        med[c] = {ph: statistics.median([r["phases"][ph] for r in trace[c]]) for ph in PHASES}
    ratios = [c for c in ("java-cold", "java-warm") if c in cols] if show_ratio and "jsgo" in cols else []
    out = ["| Phase | " + " | ".join(label(c) + " ms" for c in cols)
           + "".join(f" | speed-up vs {label(c)}" for c in ratios) + " |",
           "|---|" + "---|" * (len(cols) + len(ratios))]
    for ph in PHASES:
        cells = [f"{med[c][ph]:.0f}" for c in cols]
        extra = "".join(" | " + (f"{med[c][ph] / med['jsgo'][ph]:.2f}×" if med["jsgo"][ph] and med[c][ph] else "–")
                        for c in ratios)
        out.append(f"| {ph} | " + " | ".join(cells) + extra + " |")
    tot = {c: sum(med[c].values()) for c in cols}
    out.append("| **total** | " + " | ".join(f"**{tot[c]:.0f}**" for c in cols)
               + "".join(f" | **{tot[c] / tot['jsgo']:.2f}×**" for c in ratios) + " |")
    out.append("")
    # Coverage and line-up.
    notes = []
    for c in cols:
        r0 = trace[c][0]
        traced = sum(v for k, v in med[c].items() if k not in ("startup", "code printing and source maps",
                                                               "untraced (inside tracker)"))
        notes.append(f"{label(c)}: {r0.get('reports', 0)}/{r0.get('jobs', 1)} report(s) parsed per run, "
                     f"{len(r0['passes'])} distinct passes, passes cover {100 * traced / tot[c] if tot[c] else 0:.0f}% of the time"
                     + (f"; unmapped: {', '.join(r0['unmapped'])}" if r0["unmapped"] else ""))
    names = {c: {p: s["runs"] for p, s in trace[c][0]["passes"].items()} for c in cols}
    if "java-cold" in names and "jsgo" in names:
        a, b = names["java-cold"], names["jsgo"]
        only_j, only_g = sorted(set(a) - set(b)), sorted(set(b) - set(a))
        runs_diff = sorted(p for p in set(a) & set(b) if a[p] != b[p])
        notes.append("Pass line-up Java cold vs jsgo: " + (
            "identical pass names and run counts." if not (only_j or only_g or runs_diff) else
            f"only in Java: {only_j[:10]}; only in jsgo: {only_g[:10]}; different run counts: {runs_diff[:10]}"))
    if "java-cold" in names and "java-warm" in names:
        same = names["java-cold"] == names["java-warm"]
        notes.append("Pass line-up Java cold vs warm: " + ("identical." if same else "differs."))
    for c in missing:
        notes.append(f"{label(c)}: no tracer report in any trace run (it stopped before printing one).")
    out += ["- " + n for n in notes] + [""]
    # Per-pass details.
    allp = set()
    for c in cols:
        allp |= set(trace[c][0]["passes"])
    pm = load_phase_map()
    def pmed(c, p):
        vals = [r["passes"].get(p, {}).get("ms") for r in trace[c]]
        vals = [v for v in vals if v is not None]
        return statistics.median(vals) if vals else None
    key = cols[0]
    rows = sorted(allp, key=lambda p: -(pmed(key, p) or 0))
    out += ["<details><summary>Per pass (ms, median; runs = times the pass ran in one run of the task)</summary>", "",
            "| Pass | Phase | Runs | " + " | ".join(label(c) for c in cols) + " |",
            "|---|---|---|" + "---|" * len(cols)]
    for p in rows:
        runs = trace[key][0]["passes"].get(p, {}).get("runs", "")
        out.append(f"| {p} | {pm.get(p, 'other passes')} | {runs} | "
                   + " | ".join("–" if pmed(c, p) is None else f"{pmed(c, p):.0f}" for c in cols) + " |")
    out += ["", "</details>", ""]
    return out


def label(c):
    return {"java-cold": "Java cold", "java-warm": "Java warm", "jsgo": "jsgo"}[c]


# ---------------------------------------------------------------------------------------------
# Single-file sample


def sample(cache, check):
    with open(os.path.join(cache, "corpus-files.txt")) as f:
        names = [l.strip() for l in f if l.strip()]
    corpus = os.path.join(cache, "corpus")
    eligible = [n for n in names
                if not re.match(r"test262-parser-tests@[^/]+/package/(fail|early)/", n)
                and os.path.getsize(os.path.join(corpus, n)) <= SAMPLE_MAX_BYTES]
    picked = sorted(random.Random(SAMPLE_SEED).sample(eligible, SAMPLE_SIZE))
    header = [f"# {SAMPLE_SIZE} single-file compiles for tasks/single-files.json, drawn by",
              f"# `bench.py sample` from .cache/corpus-files.txt ({len(names)} files) with",
              f"# random.Random({SAMPLE_SEED}).sample, excluding test262-parser-tests fail/ and early/",
              f"# (invalid on purpose) and files over {SAMPLE_MAX_BYTES // 1024} KiB ({len(eligible)} eligible).",
              "# Paths are relative to .cache/corpus/."]
    text = "\n".join(header + picked) + "\n"
    if check:
        with open(SINGLE_FILES) as f:
            if f.read() != text:
                sys.exit("tasks/single-files.txt is stale: run `bench.py sample`")
        print("tasks/single-files.txt is current")
        return
    with open(SINGLE_FILES, "w") as f:
        f.write(text)
    print(f"wrote {SINGLE_FILES}: {SAMPLE_SIZE} of {len(eligible)} eligible files")


# ---------------------------------------------------------------------------------------------


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = ap.add_subparsers(dest="cmd", required=True)
    r = sub.add_parser("run", help="measure")
    r.add_argument("--out", required=True, help="results directory (re-runs merge into it per task/compiler)")
    r.add_argument("--tasks", help="comma-separated task names (default: all)")
    r.add_argument("--runs", type=int, help="cold runs per compiler (default: the task's, 5)")
    r.add_argument("--warmup", type=int, help="warm-up iterations of the warm driver (default: the task's)")
    r.add_argument("--warm-iters", type=int, help="measured warm iterations (default: the task's, 5)")
    r.add_argument("--trace-runs", type=int, help="--tracer_mode runs per compiler (default: the task's, 3)")
    r.add_argument("--jsgo", help="a jsgo binary (default: build ./cmd/jsgo of this checkout)")
    r.add_argument("--jsgo-ref", help="build jsgo from this git ref instead (git archive + go build)")
    r.add_argument("--skip-java", action="store_true", help="only jsgo (gated against the stored Java run 1)")
    r.add_argument("--skip-warm", action="store_true", help="no Java warm rows")
    r.add_argument("--skip-jsgo", action="store_true")
    r.add_argument("--skip-trace", action="store_true", help="no per-phase runs")
    r.add_argument("--nice", type=int, default=10, help="niceness of the measured processes (default 10)")
    r.add_argument("--locale", default="C.UTF-8", help="LC_ALL for both compilers (default C.UTF-8)")
    r.add_argument("--cache", default=None, help="the .cache directory (default: the main checkout's)")
    r.add_argument("--report", help="also write the Markdown report here")
    p = sub.add_parser("report", help="write the Markdown report")
    p.add_argument("--out", required=True)
    p.add_argument("--md", help="output file (default: stdout)")
    p.add_argument("--final", action="store_true", help="omit the 'preliminary' banner")
    p.add_argument("--notes", help="a Markdown file to include after the banner (findings, context)")
    g = sub.add_parser("regate", help="recompute gate.json from the stored first runs (no compiling)")
    g.add_argument("--out", required=True)
    g.add_argument("--tasks", help="comma-separated task names (default: all present)")
    s = sub.add_parser("sample", help="(re)generate tasks/single-files.txt")
    s.add_argument("--check", action="store_true")
    s.add_argument("--cache", default=None)
    a = ap.parse_args()

    if a.cmd == "report":
        text = report(a.out, a.md, preliminary=not a.final, notes=a.notes)
        if not a.md:
            print(text)
        return
    if a.cmd == "regate":
        for t in load_tasks(a.tasks.split(",") if a.tasks else None):
            tdir = os.path.join(a.out, t["name"])
            if os.path.isdir(os.path.join(tdir, "jsgo", "run-1")):
                print(f"== {t['name']}")
                gate(t, load(os.path.join(tdir, "jobs.json")), tdir)
        return
    a.cache = a.cache or default_cache()
    if a.cmd == "sample":
        sample(a.cache, a.check)
        return

    runner = Runner(a)
    if not a.skip_java:
        runner.build_batch()
    jsgo_info = {"status": "skipped"} if a.skip_jsgo else runner.resolve_jsgo()
    print(f"jsgo: {jsgo_info}", flush=True)
    mpath = os.path.join(runner.out, "machine.json")
    machine = load(mpath) or {}
    if not a.skip_java or not machine:
        first = {k: machine[k] for k in ("date", "loadavg_start") if k in machine}
        machine.update(machine_info(runner, jsgo_info))
        machine.update(first)
    else:
        machine["jsgo"] = jsgo_info
        machine["jsgo_measured"] = (datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M UTC")
                                    + ", separately from the Java runs (--skip-java), not interleaved with them")
    save(mpath, machine)
    for t in load_tasks(a.tasks.split(",") if a.tasks else None):
        runner.run_task(t)
        machine["loadavg_end"] = loadavg()
        save(mpath, machine)
    if a.report:
        report(runner.out, a.report)


if __name__ == "__main__":
    main()
