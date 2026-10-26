#!/usr/bin/env python3
"""Counterbalanced interleaved suite: cold Java vs one or more jsgo binaries.

  python3 -I tools/bench/interleaved.py --out DIR --tasks a,b [--blocks 5] [--threads 4] \
      --jsgo name=BIN[@THREADS] ... [--java JAVA] [--warm [--warmup 5 --warm-iters 5]] \
      [--no-java --java-ref DIR] [--repo REPO]

Block i runs the contenders in the order Java, jsgo... when i is even and the reverse
when odd (docs/benchmarks.md, "v20261021 on JDK 27"). Every run is gated byte for byte
against the Java reference: Java's first cold run in this suite (DIR/<task>/ref), or
the one an earlier run left in --java-ref. The JVM's own "WARNING: " lines (the
sun.misc.Unsafe notice of JDK 26+) are dropped from Java's stderr first. 'googleapis' is
the three-stage build of run_googleapis_bench.py; its gate is a digest of the chunk
outputs plus each stage's stdout, stderr and exit. With --warm, Java is also run warm
(java/CliBatch.java, one JVM per task); for 'single-files' Java is only run warm and its
last iteration is the reference. One JSON record per run goes to DIR/results.jsonl;
interleaved_stats.py summarizes them.
"""
import argparse
import glob
import hashlib
import json
import math
import os
import shutil
import statistics
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
JAVA = "java"  # set from --java in main
T975_DF = {1: 12.706, 2: 4.303, 3: 3.182, 4: 2.776, 5: 2.571, 6: 2.447, 7: 2.365, 8: 2.306, 9: 2.262}


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--repo", default=os.path.dirname(os.path.dirname(HERE)))
    ap.add_argument("--java", default=os.environ.get("JSGO_BENCH_JAVA", "java"),
                    help="the java launcher (default: $JSGO_BENCH_JAVA or java on PATH)")
    ap.add_argument("--out", required=True)
    ap.add_argument("--tasks", required=True)
    ap.add_argument("--blocks", type=int, default=5)
    ap.add_argument("--threads", type=int, default=4)
    ap.add_argument("--jsgo", action="append", default=[], help="name=path[@threads]")
    ap.add_argument("--java-threads", type=int, default=4)
    ap.add_argument("--warm", action="store_true", help="also run Java warm (CliBatch) per task")
    ap.add_argument("--warmup", type=int, default=5)
    ap.add_argument("--warm-iters", type=int, default=5)
    ap.add_argument("--no-java", action="store_true", help="time only jsgo; gate against --java-ref")
    ap.add_argument("--java-ref", default="", help="dir holding <task>/ref from an earlier suite run")
    ap.add_argument("--jar", default="")
    a = ap.parse_args()
    global JAVA
    JAVA = shutil.which(a.java) or a.java
    sys.path.insert(0, os.path.join(a.repo, "tools", "bench"))
    import bench  # noqa: E402

    cache = bench.default_cache()
    jar = a.jar or os.path.join(cache, bench.read_pin()["CLOSURE_JAR"])
    java_env = {k: v for k, v in os.environ.items() if k not in bench.JAVA_ENV_DROP and k != "JAVA_HOME"}
    go_env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP and not k.startswith("JSGO_")}
    for e in (java_env, go_env):
        e["LC_ALL"] = "C.UTF-8"
    contenders = []
    if not a.no_java:
        contenders.append(("java", [JAVA, "-jar", jar], java_env, True, a.java_threads))
    for spec in a.jsgo:
        name, path = spec.split("=", 1)
        th = a.threads
        if "@" in path:
            path, th = path.rsplit("@", 1)
            th = int(th)
        contenders.append((name, [os.path.abspath(path)], go_env, False, th))
    os.makedirs(a.out, exist_ok=True)
    resf = os.path.join(a.out, "results.jsonl")
    names = [t.strip() for t in a.tasks.split(",") if t.strip()]

    for name in names:
        tdir = os.path.join(a.out, name)
        ref = os.path.join(a.java_ref or a.out, name, "ref")
        cs = contenders
        if name == "single-files":
            # 1,000 JVM launches per cold run: Java is measured warm only, and its last warm
            # iteration is the reference.
            cs = [c for c in contenders if not c[3]]
            if a.warm or not os.path.exists(os.path.join(ref, "done")):
                warm_java(bench, cache, jar, java_env, a, name, tdir, ref, resf)
        for b in range(a.blocks):
            order = cs if b % 2 == 0 else cs[::-1]
            for cname, argv0, env, is_java, th in order:
                tflag = [f"--num_parallel_threads={th}"]
                rundir = os.path.join(tdir, cname)
                if name == "googleapis":
                    rec = run_gapi(bench, argv0, tflag, env, rundir, is_java)
                    if is_java and not os.path.exists(os.path.join(ref, "done")):
                        save_gapi_ref(rundir, ref, rec)
                    rec["gate"] = gate_gapi(rundir, ref, rec)
                else:
                    task = bench.load_tasks([name])[0]
                    bench.make_workspace(rundir, cache)
                    jobs = bench.task_jobs(task, rundir)
                    rec = run_task(bench, argv0, tflag, env, task, jobs, rundir)
                    if is_java and not os.path.exists(os.path.join(ref, "done")):
                        save_task_ref(bench, task, jobs, rundir, ref)
                    rec["gate"] = gate_task(bench, task, jobs, rundir, ref)
                rec.update({"task": name, "who": cname, "block": b, "time": time.time()})
                with open(resf, "a") as f:
                    f.write(json.dumps(rec) + "\n")
                print(f"{name:20} b{b} {cname:8} wall {rec['wall_s']:8.2f}  cpu {rec['cpu_s']:8.2f}  "
                      f"rss {rec['rss_mb']:8.0f} MiB  {rec['gate']}", flush=True)
        if a.warm and name not in ("googleapis", "single-files"):
            warm_java(bench, cache, jar, java_env, a, name, tdir, ref, resf)
    summarize(resf, names)


def warm_java(bench, cache, jar, java_env, a, name, tdir, ref, resf):
    """W warm-up + N measured iterations of the task's invocations in one JVM (tools/bench's
    CliBatch). Its last iteration's outputs are gated against the reference (or become it)."""
    build = os.path.join(cache, "bench-build", "java-warm-" + os.path.basename(jar))
    os.makedirs(build, exist_ok=True)
    src = os.path.join(os.path.dirname(bench.__file__), "java", "CliBatch.java")
    if not os.path.exists(os.path.join(build, "CliBatch.class")):
        subprocess.run([os.path.join(os.path.dirname(JAVA), "javac"), "-nowarn", "-cp", jar, "-d", build, src], env=java_env, check=True)
    task = bench.load_tasks([name])[0]
    rundir = os.path.join(tdir, "java-warm")
    shutil.rmtree(rundir, ignore_errors=True)
    bench.make_workspace(rundir, cache)
    jobs = bench.task_jobs(task, rundir)
    tflag = [f"--num_parallel_threads={a.java_threads}"]
    jobfile = os.path.join(rundir, "jobs.tsv")
    with open(jobfile, "w") as f:
        for args in jobs:
            f.write("\t".join(args + tflag) + "\n")
    os.makedirs(os.path.join(rundir, "batch"), exist_ok=True)
    r = bench.run_proc([JAVA, "-cp", jar + ":" + build, "CliBatch", str(a.warmup), str(a.warm_iters), jobfile, "batch"],
                       rundir, java_env, os.path.join(rundir, "report.tsv"), os.path.join(rundir, "driver.stderr"), 0)
    iters, _ = bench.parse_batch_report(os.path.join(rundir, "report.tsv"))
    single = len(jobs) == 1
    for j in range(len(jobs)):
        strip_jvm_banner(os.path.join(rundir, "batch", f"{j}.stderr"))
    # The warm driver writes job j's last outputs as batch/<j>.{stdout,stderr,exit}; the task's
    # output files (out.js, out.js.map) land in the working directory.
    got = os.path.join(rundir, "gated")
    os.makedirs(got, exist_ok=True)
    for j in range(len(jobs)):
        dst = bench.artifacts(task, got, j, single)
        for k in ("stdout", "stderr", "exit"):
            shutil.copyfile(os.path.join(rundir, "batch", f"{j}.{k}"), dst[k])
        if single:
            for k in task.get("compare", []):
                if os.path.exists(os.path.join(rundir, k)):
                    shutil.copyfile(os.path.join(rundir, k), dst[k])
    if not os.path.exists(os.path.join(ref, "done")):
        save_task_ref(bench, task, jobs, got, ref)
    gate = gate_task(bench, task, jobs, got, ref)
    for it in iters:
        if it["kind"] != "measure":
            continue
        rec = {"task": name, "who": "java-warm", "block": it["iter"], "time": time.time(),
               "wall_s": it["wall_ms"] / 1000, "cpu_s": it["cpu_ms"] / 1000, "rss_mb": r["rss_mb"], "gate": gate}
        with open(resf, "a") as f:
            f.write(json.dumps(rec) + "\n")
    print(f"{name:20} java-warm {len([i for i in iters if i['kind'] == 'measure'])} iters, median "
          f"{statistics.median([i['wall_ms'] for i in iters if i['kind'] == 'measure']) / 1000:.2f} s, "
          f"process rss {r['rss_mb']:.0f} MiB, exit {r['exit']}, {gate}", flush=True)


def run_task(bench, argv0, tflag, env, task, jobs, rundir):
    single = len(jobs) == 1
    tot = {"wall_s": 0.0, "cpu_s": 0.0, "rss_mb": 0.0}
    for j, args in enumerate(jobs):
        arts = bench.artifacts(task, rundir, j, single)
        for k in task.get("compare", []):
            if single and os.path.exists(arts[k]):
                os.remove(arts[k])
        r = bench.run_proc(argv0 + args + tflag, rundir, env, arts["stdout"], arts["stderr"], 0)
        if argv0[0] == JAVA:
            strip_jvm_banner(arts["stderr"])
        with open(arts["exit"], "w") as f:
            f.write(f"{r['exit']}\n")
        tot["wall_s"] += r["wall_s"]
        tot["cpu_s"] += r["cpu_s"]
        tot["rss_mb"] = max(tot["rss_mb"], r["rss_mb"])
    return tot


def strip_jvm_banner(path):
    """Drops the JVM's own 'WARNING: ...' lines (e.g. the sun.misc.Unsafe deprecation banner) from
    Java's stderr; Closure's diagnostics are 'file:line:col: WARNING - ...' and are kept. The raw
    file is kept as <path>.raw."""
    data = open(path, "rb").read()
    lines = data.split(b"\n")
    kept = [l for l in lines if not l.startswith(b"WARNING: ")]
    if len(kept) != len(lines):
        shutil.copyfile(path, path + ".raw")
        open(path, "wb").write(b"\n".join(kept))


def save_task_ref(bench, task, jobs, rundir, ref):
    shutil.rmtree(ref, ignore_errors=True)
    os.makedirs(ref)
    single = len(jobs) == 1
    for j in range(len(jobs)):
        src = bench.artifacts(task, rundir, j, single)
        dst = bench.artifacts(task, ref, j, single)
        for k in src:
            if os.path.exists(src[k]):
                shutil.copyfile(src[k], dst[k])
    open(os.path.join(ref, "done"), "w").write("ok\n")


def gate_task(bench, task, jobs, rundir, ref):
    """MATCH; "MATCH (java-frames: N jobs)" when N invocations differ only in the JVM's stack-frame
    lines on stderr (both compilers threw the same exception; bench.py's java-frames class);
    otherwise "DIFF <count>: " and the first three differences."""
    single = len(jobs) == 1
    bad, frames_only = [], set()
    for j in range(len(jobs)):
        ra = bench.artifacts(task, ref, j, single)
        ga = bench.artifacts(task, rundir, j, single)
        for k in ra:
            if not bench.same_file(ra[k], ga[k]):
                if k == "stderr" and os.path.exists(ga[k]) and bench.without_java_frames(
                        bench.read_text(ra[k])) == bench.without_java_frames(bench.read_text(ga[k])):
                    frames_only.add(j)
                    continue
                bad.append(f"job{j}:{k}: {bench.describe_diff(ra[k], ga[k])}")
    if bad:
        return f"DIFF {len(bad)}: " + " | ".join(bad[:3])
    if frames_only:
        return f"MATCH (java-frames: {len(frames_only)} jobs)"
    return "MATCH"


GAPI_STAGES = ("CHECKS", "OPTIMIZATIONS", "FINALIZATIONS")


def run_gapi(bench, argv0, tflag, env, rundir, is_java):
    gdir = os.path.join(bench.default_cache(), "bench-inputs", "googleapis-bench")
    flagfile = os.path.join(gdir, "flags.txt")
    shutil.rmtree(rundir, ignore_errors=True)
    os.makedirs(os.path.join(rundir, "chunks"))
    base = list(argv0)
    if is_java:
        base = [argv0[0], "-Xmx20g"] + argv0[1:]
    base += [f"--flagfile={flagfile}"] + tflag
    cs, os_ = os.path.join(rundir, "checks.state"), os.path.join(rundir, "optimizations.state")
    extra = {
        "CHECKS": ["--segment_of_compilation_to_run=CHECKS", f"--filename_to_save_to={cs}"],
        "OPTIMIZATIONS": ["--segment_of_compilation_to_run=OPTIMIZATIONS", f"--filename_to_restore_from={cs}",
                          f"--filename_to_save_to={os_}"],
        "FINALIZATIONS": ["--segment_of_compilation_to_run=FINALIZATIONS", f"--filename_to_restore_from={os_}",
                          f"--chunk_output_path_prefix={rundir}/chunks/"],
    }
    tot = {"wall_s": 0.0, "cpu_s": 0.0, "rss_mb": 0.0, "stages": {}}
    for st in GAPI_STAGES:
        r = bench.run_proc(base + extra[st], rundir, env, os.path.join(rundir, f"{st}.stdout"),
                           os.path.join(rundir, f"{st}.stderr"), 0)
        if is_java:
            strip_jvm_banner(os.path.join(rundir, f"{st}.stderr"))
        tot["stages"][st] = {k: r[k] for k in ("wall_s", "cpu_s", "rss_mb", "exit")}
        tot["wall_s"] += r["wall_s"]
        tot["cpu_s"] += r["cpu_s"]
        tot["rss_mb"] = max(tot["rss_mb"], r["rss_mb"])
        if r["exit"] != 0:
            break
    h = hashlib.sha256()
    for p in sorted(glob.glob(os.path.join(rundir, "chunks", "*.js"))):
        h.update(os.path.basename(p).encode())
        h.update(b"\0")
        h.update(hashlib.sha256(open(p, "rb").read()).digest())
    tot["chunks_sha256"] = h.hexdigest()
    for f in (cs, os_):
        if os.path.exists(f):
            os.remove(f)
    shutil.rmtree(os.path.join(rundir, "chunks"))
    return tot


def save_gapi_ref(rundir, ref, rec):
    shutil.rmtree(ref, ignore_errors=True)
    os.makedirs(ref)
    for st in GAPI_STAGES:
        for k in ("stdout", "stderr"):
            p = os.path.join(rundir, f"{st}.{k}")
            if os.path.exists(p):
                shutil.copyfile(p, os.path.join(ref, f"{st}.{k}"))
    json.dump({"chunks_sha256": rec["chunks_sha256"],
               "exits": {s: v["exit"] for s, v in rec["stages"].items()}}, open(os.path.join(ref, "ref.json"), "w"))
    open(os.path.join(ref, "done"), "w").write("ok\n")


def gate_gapi(rundir, ref, rec):
    r = json.load(open(os.path.join(ref, "ref.json")))
    bad = []
    if rec["chunks_sha256"] != r["chunks_sha256"]:
        bad.append("chunks")
    for st in GAPI_STAGES:
        if rec["stages"].get(st, {}).get("exit") != r["exits"].get(st):
            bad.append(f"{st}:exit")
        for k in ("stdout", "stderr"):
            a, b = os.path.join(ref, f"{st}.{k}"), os.path.join(rundir, f"{st}.{k}")
            if (open(a, "rb").read() if os.path.exists(a) else b"") != (open(b, "rb").read() if os.path.exists(b) else b""):
                bad.append(f"{st}:{k}")
    return "MATCH" if not bad else "DIFF " + ",".join(bad)


def ci_ratio(num, den):
    """Geometric mean of num/den with a paired-t 95% CI on log ratios."""
    logs = [math.log(x / y) for x, y in zip(num, den)]
    m = statistics.mean(logs)
    if len(logs) < 2:
        return math.exp(m), math.exp(m), math.exp(m)
    h = T975_DF.get(len(logs) - 1, 2.0) * statistics.stdev(logs) / math.sqrt(len(logs))
    return math.exp(m), math.exp(m - h), math.exp(m + h)


def summarize(resf, names):
    recs = [json.loads(l) for l in open(resf)]
    for name in names:
        rs = [r for r in recs if r["task"] == name]
        who = sorted({r["who"] for r in rs}, key=lambda w: (w != "java", w))
        by = {w: sorted([r for r in rs if r["who"] == w], key=lambda r: r["block"]) for w in who}
        print(f"\n== {name}")
        for w in who:
            g = {r["gate"][:5] for r in by[w]}
            print(f"  {w:8} n={len(by[w])} wall {statistics.mean(r['wall_s'] for r in by[w]):8.2f} s  "
                  f"cpu {statistics.mean(r['cpu_s'] for r in by[w]):8.2f} s  "
                  f"rss {statistics.mean(r['rss_mb'] for r in by[w]):8.0f} MiB  gate {sorted(g)}")
        pairs = [(x, y) for x in who for y in who if x != y and (x == "java" or (x, y) == ("main", "merged"))]
        for x, y in pairs:
            n = min(len(by[x]), len(by[y]))
            for k in ("wall_s", "cpu_s", "rss_mb"):
                g, lo, hi = ci_ratio([r[k] for r in by[x][:n]], [r[k] for r in by[y][:n]])
                print(f"  {x}/{y} {k:7} {g:6.3f}x [{lo:.3f}, {hi:.3f}]")


if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "summary":
        summarize(sys.argv[2], sys.argv[3].split(","))
    else:
        main()
