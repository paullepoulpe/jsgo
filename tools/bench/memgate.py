#!/usr/bin/env python3
"""Memory gate: byte-identity against Java plus heap/GC numbers, one run per task.

  python3 -I tools/bench/memgate.py --jsgo BIN [--tasks a,b] [--fast] [--threads N] [--out DIR]

For each task (tools/bench/tasks/*.json) it runs Java once and keeps its outputs as
the reference (in --golden, reused across calls), then runs jsgo with
JSGO_MEMSTATS and prints one row: gate status, wall, CPU, peak RSS, peak live
heap (largest heap surviving a GC), bytes and objects allocated, GC cycles and
GC CPU. With --fast, jsgo runs with --fast_mode and the gate only checks the exit
status (fast mode may change output; tools/bench/fast_mode_gate.py checks it).
Results are appended to <out>/results.jsonl (default: .memgate/ in this checkout).

Env: JSGO_PEAK_HEAP_PROFILE=<file> makes jsgo write a heap profile at its peak
live heap (cmd/jsgo/memstats.go); inspect it with
`go tool pprof -sample_index=inuse_space -top <binary> <file>`.
"""
import argparse
import json
import os
import shutil
import sys
import time

sys.dont_write_bytecode = True
HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import bench  # noqa: E402

DEFAULT_TASKS = "typed-app,downlevel-simple,checks-only,bundle,shaka-blockly"


def run_jobs(argv0, jobs, task, rundir, cache, env, nice):
    shutil.rmtree(rundir, ignore_errors=True)
    bench.make_workspace(rundir, cache)
    single = len(jobs) == 1
    tot = {"wall_s": 0.0, "cpu_s": 0.0, "rss_mb": 0.0, "exits": []}
    stats = {}
    for j, args in enumerate(jobs):
        arts = bench.artifacts(task, rundir, j, single)
        jenv = dict(env)
        statf = os.path.join(rundir, f"{j}.memstats.json")
        if env is not None and "JSGO_MEMSTATS_ON" in env:
            jenv.pop("JSGO_MEMSTATS_ON")
            jenv["JSGO_MEMSTATS"] = statf
        r = bench.run_proc(argv0 + args, rundir, jenv, arts["stdout"], arts["stderr"], nice)
        with open(arts["exit"], "w") as f:
            f.write(f"{r['exit']}\n")
        tot["wall_s"] += r["wall_s"]
        tot["cpu_s"] += r["cpu_s"]
        tot["rss_mb"] = max(tot["rss_mb"], r["rss_mb"])
        tot["exits"].append(r["exit"])
        if os.path.exists(statf):
            with open(statf) as f:
                s = json.load(f)
            for k, v in s.items():
                if k.startswith("peak"):
                    stats[k] = max(stats.get(k, 0), v)
                else:
                    stats[k] = stats.get(k, 0) + v
    tot.update(stats)
    return tot


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--jsgo", required=True)
    ap.add_argument("--tasks", default=DEFAULT_TASKS)
    # Per checkout (worktrees run side by side); the Java references are shared.
    ap.add_argument("--out", default=os.path.join(bench.REPO, ".memgate"))
    ap.add_argument("--golden", default=os.path.join(bench.default_cache(), "memgate-golden"))
    ap.add_argument("--fast", action="store_true")
    ap.add_argument("--threads", type=int, default=1)
    ap.add_argument("--label", default="")
    ap.add_argument("--nice", type=int, default=0)
    ap.add_argument("--runs", type=int, default=1)
    a = ap.parse_args()
    cache = bench.default_cache()
    jar = os.path.join(cache, bench.read_pin()["CLOSURE_JAR"])
    java_env = {k: v for k, v in os.environ.items() if k not in bench.JAVA_ENV_DROP}
    go_env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP and k != "JSGO_FAST_MODE"}
    for e in (java_env, go_env):
        e["LC_ALL"] = "C"
    go_env["JSGO_MEMSTATS_ON"] = "1"
    jsgo = os.path.abspath(a.jsgo)
    os.makedirs(a.out, exist_ok=True)
    extra = []
    if a.threads != 1:
        extra.append(f"--num_parallel_threads={a.threads}")
    if a.fast:
        extra.append("--fast_mode")
    tasks = bench.load_tasks(a.tasks.split(","))
    hdr = (f"{'task':22} {'gate':6} {'wall':>7} {'cpu':>7} {'rss MB':>7} {'live MB':>8} "
           f"{'scan MB':>8} {'alloc GB':>8} {'Mobj':>7} {'gcs':>5} {'gcCPU':>6}")
    print(hdr, flush=True)
    for task in tasks:
        name = task["name"]
        gdir = os.path.join(a.golden, name)
        ws = os.path.join(a.out, name, "ws")
        bench.make_workspace(ws, cache)
        jobs = bench.task_jobs(task, ws)
        single = len(jobs) == 1
        ref = os.path.join(gdir, "run")
        if not os.path.exists(os.path.join(gdir, "done")):
            print(f"   ({name}: running Java for the reference)", flush=True)
            run_jobs(["java", "-jar", jar], jobs, task, ref, cache, java_env, a.nice)
            with open(os.path.join(gdir, "done"), "w") as f:
                f.write("ok\n")
        for i in range(a.runs):
            got = os.path.join(a.out, name, "jsgo")
            tot = run_jobs([jsgo] + extra, jobs, task, got, cache, go_env, a.nice)
            bad = []
            for j in range(len(jobs)):
                ra = bench.artifacts(task, ref, j, single)
                ga = bench.artifacts(task, got, j, single)
                keys = ["exit"] if a.fast else list(ra)
                for k in keys:
                    if not bench.same_file(ra[k], ga[k]):
                        bad.append((j, k, bench.describe_diff(ra[k], ga[k])))
            status = "MATCH" if not bad else f"DIFF{len(bad)}"
            print(f"{name:22} {status:6} {tot['wall_s']:7.2f} {tot['cpu_s']:7.2f} {tot['rss_mb']:7.0f} "
                  f"{tot.get('peak_live_bytes', 0) / 2**20:8.0f} {tot.get('peak_scan_bytes', 0) / 2**20:8.0f} "
                  f"{tot.get('alloc_bytes', 0) / 2**30:8.2f} "
                  f"{tot.get('alloc_objects', 0) / 1e6:7.1f} {tot.get('gc_cycles', 0):5.0f} "
                  f"{tot.get('gc_cpu_s', 0):6.2f}", flush=True)
            for b in bad[:4]:
                print(f"     job {b[0]} {b[1]}: {b[2]}", flush=True)
            rec = {"time": time.time(), "label": a.label, "jsgo": jsgo, "task": name, "fast": a.fast,
                   "threads": a.threads, "status": status, **{k: v for k, v in tot.items() if k != "exits"}}
            with open(os.path.join(a.out, "results.jsonl"), "a") as f:
                f.write(json.dumps(rec) + "\n")


if __name__ == "__main__":
    main()
