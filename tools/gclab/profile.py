#!/usr/bin/env python3
"""One profiling run of a task with the gclab hook.

  python3 -I profile.py --bin JSGO --task T --out DIR --mode MODE [--env K=V ...]

Modes (all at the default GOGC unless --env says otherwise):
  live      JSGO_GCLAB_LIVE=1 (a forced GC after every pass): live and scannable heap per pass
  heap      heap profiles at JSGO_GCLAB_HEAPAT and at exit, GODEBUG=memprofilerate=4096
  cpu       a CPU profile of the whole run
  gctrace   GODEBUG=gctrace=1 (stderr not compared)
  locality  JSGO_GCLAB_LOCALITY: AST walk timings at the given passes
Every mode also writes the per-pass log and the GC log.
"""
import argparse, os, shutil, sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from common import *  # noqa

ap = argparse.ArgumentParser()
ap.add_argument("--bin", required=True)
ap.add_argument("--task", required=True)
ap.add_argument("--out", required=True)
ap.add_argument("--mode", required=True)
ap.add_argument("--tag", default="")
ap.add_argument("--env", action="append", default=[])
a = ap.parse_args()
out = os.path.abspath(a.out)
os.makedirs(out, exist_ok=True)
tag = a.mode + (("-" + a.tag) if a.tag else "")
pre = os.path.join(out, tag)
ws = os.path.join(out, "ws-" + tag)
shutil.rmtree(ws, ignore_errors=True)
bench.make_workspace(ws, CACHE)
(_, args), = jobs_of(a.task, ws)
env = {"JSGO_GCLAB_PASSES": pre + "-passes.tsv", "JSGO_GCLAB_GCLOG": pre + "-gclog.tsv",
       "JSGO_GCLAB_STATS": pre + "-stats.json"}
if a.mode == "live":
    env["JSGO_GCLAB_LIVE"] = "1"
elif a.mode == "heap":
    d = pre + "-heap"
    shutil.rmtree(d, ignore_errors=True)
    os.makedirs(d)
    env["JSGO_GCLAB_HEAPDIR"] = d
    env["GODEBUG"] = "memprofilerate=4096"
elif a.mode == "cpu":
    env["JSGO_GCLAB_CPUPROF"] = pre + "-cpu.pprof"
elif a.mode == "gctrace":
    env["GODEBUG"] = "gctrace=1"
elif a.mode == "locality":
    env["JSGO_GCLAB_LOCALITY_OUT"] = pre + "-locality.tsv"
    if os.path.exists(pre + "-locality.tsv"):
        os.remove(pre + "-locality.tsv")
for kv in a.env:
    k, v = kv.split("=", 1)
    env[k] = v
arts = bench.artifacts(TASKS[a.task], ws, 0, True)
p = bench.run_proc([os.path.abspath(a.bin)] + args, ws, go_env(env), arts["stdout"], arts["stderr"], 10)
with open(arts["exit"], "w") as f:
    f.write(f"{p['exit']}\n")
d = bench.digest_artifacts(arts)
shutil.copyfile(arts["stderr"], pre + "-stderr")
ref = ref_digests(a.task, "jsgo")[0]
print(f"{a.task} {tag}: wall {p['wall_s']:.2f}s cpu {p['cpu_s']:.2f}s rss {p['rss_mb']:.0f}MB exit {p['exit']} "
      f"outputs {'identical' if d == ref else 'DIFFER: ' + str([k for k in d if d[k] != ref.get(k)])}", flush=True)
with open(pre + "-run.json", "w") as f:
    json.dump(dict(p, env=env, outputs_identical=(d == ref)), f, indent=1)
shutil.rmtree(ws, ignore_errors=True)
