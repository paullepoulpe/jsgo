#!/usr/bin/env python3
"""GC settings sweep, measured as tools/bench/bench.py measures (run_proc: nice 10, wait4 wall,
CPU and peak RSS; GOGC & co removed from the environment, LC_ALL=C.UTF-8). The settings are
interleaved run by run, in alternating order. Every run also records runtime/metrics totals at
exit (JSGO_GCLAB_STATS, read once at exit) and checks its outputs byte for byte against the
stored GOGC=100 jsgo run 1 and Java cold run 1 of bench2.

  python3 -I sweep.py --bin JSGO --out DIR [--tasks a,b] [--runs N] [--settings a,b]
"""
import argparse, json, os, shutil, statistics, sys, time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from common import *  # noqa

SETTINGS = {
    "gogc100": {},
    "gogc200": {"GOGC": "200"},
    "gogc400": {"GOGC": "400"},
    "off10g": {"GOGC": "off", "GOMEMLIMIT": "10GiB"},
    "offjava": {"GOGC": "off"},  # GOMEMLIMIT = Java cold's peak RSS (per task, per file for single-100)
    "pristine": {},  # GOGC=100 with --pristine, the binary built without the hook (control)
}
STAT_KEYS = {
    "gc_cycles": "/gc/cycles/total:gc-cycles",
    "gc_cpu_s": "/cpu/classes/gc/total:cpu-seconds",
    "gc_assist_s": "/cpu/classes/gc/mark/assist:cpu-seconds",
    "gc_dedicated_s": "/cpu/classes/gc/mark/dedicated:cpu-seconds",
    "gc_idle_s": "/cpu/classes/gc/mark/idle:cpu-seconds",
    "gc_pause_cpu_s": "/cpu/classes/gc/pause:cpu-seconds",
    "gc_pause_wall_s": "/sched/pauses/total/gc:seconds#sum_upper",
    "scavenge_cpu_s": "/cpu/classes/scavenge/total:cpu-seconds",
    "alloc_bytes": "/gc/heap/allocs:bytes",
    "alloc_objs": "/gc/heap/allocs:objects",
    "live_end": "/gc/heap/live:bytes",
    "scan_heap_end": "/gc/scan/heap:bytes",
    "mem_total_end": "/memory/classes/total:bytes",
}


def java_rss(task):
    """Java cold peak RSS in MiB: the task's median; for single-100 per file (median over runs)."""
    runs = json.load(open(os.path.join(REF, ref_task(task), "java-cold.json")))["runs"]
    if task == "single-100":
        n = len(runs[0]["per_job"])
        return [statistics.median(r["per_job"][j]["rss_mb"] for r in runs) for j in range(n)]
    return statistics.median(r["rss_mb"] for r in runs)


def wait_quiet(limit=1.0):
    while bench.loadavg()[0] >= limit:
        time.sleep(10)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--bin", required=True)
    ap.add_argument("--out", required=True)
    ap.add_argument("--tasks", default="typed-app,bundle,downlevel-simple,checks-only,single-100")
    ap.add_argument("--settings", default=",".join(SETTINGS))
    ap.add_argument("--runs", type=int, default=5)
    ap.add_argument("--pristine", help="binary without the hook, for the 'pristine' setting")
    a = ap.parse_args()
    out = os.path.abspath(a.out)
    a.bin = os.path.abspath(a.bin)
    os.makedirs(out, exist_ok=True)
    resf = os.path.join(out, "sweep.json")
    res = bench.load(resf) or {}
    settings = a.settings.split(",")
    for task in a.tasks.split(","):
        tdir = os.path.join(out, task)
        ws = os.path.join(tdir, "ws")
        bench.make_workspace(ws, CACHE)
        jobs = jobs_of(task, ws)
        single = task != "single-100"
        ref_go, ref_java = ref_digests(task, "jsgo"), ref_digests(task, "java-cold")
        jrss = java_rss(task)
        r = {s: [] for s in settings}
        # One untimed warm-up pass (page cache, binary) at the default setting, then a quiet machine.
        wdir = os.path.join(tdir, "warmup")
        bench.make_workspace(wdir, CACHE)
        for j, args in jobs:
            arts = bench.artifacts(TASKS[ref_task(task)], wdir, j, single)
            bench.run_proc([a.bin] + args, wdir, go_env(), arts["stdout"], arts["stderr"], 10)
        shutil.rmtree(wdir, ignore_errors=True)
        wait_quiet()
        print(f"== {task}: {len(jobs)} job(s), load {bench.loadavg()[0]:.2f}", flush=True)
        for i in range(a.runs):
            order = settings if i % 2 == 0 else list(reversed(settings))
            for s in order:
                rundir = os.path.join(tdir, s)
                shutil.rmtree(rundir, ignore_errors=True)
                bench.make_workspace(rundir, CACHE)
                tot = {"wall_s": 0.0, "cpu_s": 0.0, "user_s": 0.0, "sys_s": 0.0, "rss_mb": 0.0}
                st = {k: 0.0 for k in STAT_KEYS}
                st["live_end"] = st["scan_heap_end"] = st["mem_total_end"] = 0
                same_go = same_java = 0
                loads = []
                for n, (j, args) in enumerate(jobs):
                    arts = bench.artifacts(TASKS[ref_task(task)], rundir, j, single)
                    sp = os.path.join(rundir, f"stats-{j}.json")
                    extra = dict(SETTINGS[s], JSGO_GCLAB_STATS=sp)
                    if s == "offjava":
                        mb = jrss[j] if isinstance(jrss, list) else jrss
                        extra["GOMEMLIMIT"] = f"{int(round(mb))}MiB"
                    binary = os.path.abspath(a.pristine) if s == "pristine" else a.bin
                    p = bench.run_proc([binary] + args, rundir, go_env(extra), arts["stdout"], arts["stderr"], 10)
                    with open(arts["exit"], "w") as f:
                        f.write(f"{p['exit']}\n")
                    for k in ("wall_s", "cpu_s", "user_s", "sys_s"):
                        tot[k] += p[k]
                    tot["rss_mb"] = max(tot["rss_mb"], p["rss_mb"])
                    loads.append(p["load1"])
                    d = bench.digest_artifacts(arts)
                    ref_idx = j if not single else 0
                    same_go += d == ref_go[ref_idx]
                    same_java += d == ref_java[ref_idx]
                    js = json.load(open(sp)) if os.path.exists(sp) else {m: 0 for m in STAT_KEYS.values()}
                    for k, m in STAT_KEYS.items():
                        if k in ("live_end", "scan_heap_end", "mem_total_end"):
                            st[k] = max(st[k], js[m])
                        else:
                            st[k] += js[m]
                    for k in arts.values():
                        if os.path.exists(k) and os.path.getsize(k) > 4096:
                            os.remove(k)
                rec = dict(tot, **st, load1=statistics.median(loads), jobs=len(jobs),
                           same_as_jsgo_ref=same_go, same_as_java_ref=same_java)
                r[s].append(rec)
                print(f"   {s:8} run {i + 1}: wall {tot['wall_s']:7.2f}s cpu {tot['cpu_s']:7.2f}s rss {tot['rss_mb']:6.0f}MB "
                      f"gc {st['gc_cycles']:5.0f} gccpu {st['gc_cpu_s']:6.2f}s alloc {st['alloc_bytes'] / 2**30:6.2f}GiB "
                      f"load {rec['load1']:.2f} same {same_go}/{len(jobs)} (java {same_java})", flush=True)
        res[task] = {"java_rss_mb": jrss, "runs": r}
        bench.save(resf, res)
        shutil.rmtree(tdir, ignore_errors=True)
    print("done", flush=True)


if __name__ == "__main__":
    main()
