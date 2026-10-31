#!/usr/bin/env python3
"""Byte-for-byte oracle for jsgo's hidden --max_optimization_loop_iterations=N flag.

Java's CLI does not expose CompilerOptions.setMaxOptimizationLoopIterations, so
java/MaxLoopRunner.java sets it from the same flag. For every task and every N,
this runs that Java driver once and jsgo in each requested configuration
(compat or --fast_mode, at a thread count) with the same argv, and compares
stdout, stderr, the exit status and the task's output files byte for byte. It
also prints the output size (raw and gzip) against the run without the flag
(N=0, the plain jar), so the cost of stopping the loops early is visible.

    python3 -I tools/bench/maxloop_oracle.py --jsgo ./jsgo --out /tmp/maxloop \\
        --tasks typed-app,typed-full --iterations 0,2,3 --configs compat@1,fast@4
"""
import argparse
import gzip
import hashlib
import json
import os
import subprocess
import sys
import time

sys.dont_write_bytecode = True
HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import bench  # noqa: E402


def digest(ws, task, proc):
    res = {"exit": proc.returncode, "stdout": hashlib.sha256(proc.stdout).hexdigest(),
           "stderr": hashlib.sha256(proc.stderr).hexdigest()}
    for name in task.get("compare", []):
        if name in ("stdout", "stderr", "exit"):
            continue
        p = os.path.join(ws, name)
        if os.path.exists(p):
            with open(p, "rb") as f:
                data = f.read()
            res[name] = hashlib.sha256(data).hexdigest()
            if name.endswith(".js"):
                res["raw"] = len(data)
                res["gzip"] = len(gzip.compress(data, 9))
            os.remove(p)
    return res


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--jsgo", required=True, type=os.path.abspath)
    ap.add_argument("--out", required=True)
    ap.add_argument("--tasks", default="typed-app,bundle-advanced,shaka-blockly,typed-full")
    ap.add_argument("--iterations", default="0,2,3", help="values of N; 0 runs without the flag")
    ap.add_argument("--configs", default="compat@1,fast@4", help="jsgo runs: mode@threads")
    ap.add_argument("--cache", default=bench.default_cache())
    a = ap.parse_args()

    jar = os.path.join(a.cache, bench.read_pin()["CLOSURE_JAR"])
    out = os.path.abspath(a.out)
    build = os.path.join(out, "java")
    os.makedirs(build, exist_ok=True)
    java_env = {k: v for k, v in os.environ.items() if k not in bench.JAVA_ENV_DROP}
    go_env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP + ("JSGO_FAST_MODE",)}
    for env in (java_env, go_env):
        env["LC_ALL"] = "C.UTF-8"
    subprocess.run(["javac", "-nowarn", "-cp", jar, "-d", build, os.path.join(HERE, "java", "MaxLoopRunner.java")],
                   env=java_env, check=True)
    ws = os.path.join(out, "ws")
    bench.make_workspace(ws, a.cache)

    failures = 0
    for task in bench.load_tasks(a.tasks.split(",")):
        jobs = bench.task_jobs(task, ws)
        if len(jobs) != 1:
            sys.exit(f"{task['name']}: only single-invocation tasks are supported")
        argv = jobs[0]
        base = None
        for n in [int(x) for x in a.iterations.split(",")]:
            flag = [f"--max_optimization_loop_iterations={n}"] if n else []
            t0 = time.perf_counter()
            p = subprocess.run(["java", "-cp", jar + ":" + build, "MaxLoopRunner"] + argv + flag,
                               cwd=ws, env=java_env, capture_output=True)
            ref = digest(ws, task, p)
            ref["wall"] = round(time.perf_counter() - t0, 2)
            if base is None:
                base = ref
            rec = {"task": task["name"], "n": n, "java": ref, "jsgo": {}}
            ok = True
            for cfg in a.configs.split(","):
                mode, threads = cfg.split("@")
                extra = [f"--num_parallel_threads={threads}"] + (["--fast_mode"] if mode == "fast" else [])
                t0 = time.perf_counter()
                p = subprocess.run([a.jsgo] + argv + extra + flag, cwd=ws, env=go_env, capture_output=True)
                got = digest(ws, task, p)
                got["wall"] = round(time.perf_counter() - t0, 2)
                same = all(got.get(k) == ref.get(k) for k in ref if k not in ("wall",))
                got["match"] = same
                ok &= same
                rec["jsgo"][cfg] = got
            failures += not ok
            size = ""
            if "raw" in ref and "raw" in base:
                size = " raw %+.3f%% gzip %+.3f%%" % (100 * (ref["raw"] / base["raw"] - 1),
                                                    100 * (ref["gzip"] / base["gzip"] - 1))
            print(f"{task['name']} N={n}: {'MATCH' if ok else 'DIFF'}{size}  java {ref['wall']} s, " +
                  ", ".join(f"{c} {g['wall']} s" for c, g in rec["jsgo"].items()), flush=True)
            with open(os.path.join(out, "results.jsonl"), "a") as f:
                f.write(json.dumps(rec) + "\n")
    sys.exit(1 if failures else 0)


if __name__ == "__main__":
    main()
