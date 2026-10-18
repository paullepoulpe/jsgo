#!/usr/bin/env python3
"""Measure a persistent Go compiler against the existing stable warm JVM reference.

The Go driver runs the real CLI entry point and matches the Java batch protocol.
These results are explicitly persistent-process results, not one-shot CLI times.
"""
import argparse
import json
import os
from pathlib import Path
import statistics
import sys
sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True
import bench

ap = argparse.ArgumentParser(description=__doc__)
ap.add_argument('--reference', required=True)
ap.add_argument('--out', required=True)
ap.add_argument('--binary', required=True)
ap.add_argument('--tasks', default='typed-app,bundle,downlevel-whitespace,downlevel-simple,checks-only')
ap.add_argument('--warmup', type=int, default=2)
ap.add_argument('--iters', type=int, default=5)
ap.add_argument('--extern-cache', action='store_true')
ap.add_argument('--gogc', type=str)
a = ap.parse_args()
root = Path(a.out).resolve()
root.mkdir(parents=True, exist_ok=True)
ref = Path(a.reference).resolve()
cache = bench.default_cache()
env = {k:v for k,v in os.environ.items() if k not in bench.GO_ENV_DROP}
env['LC_ALL'] = 'C.UTF-8'
if a.gogc is not None: env['GOGC'] = a.gogc
bench.save(str(root/'metadata.json'), {'binary': os.path.abspath(a.binary), 'sha256': bench.sha256_file(a.binary), 'affinity': sorted(os.sched_getaffinity(0)), 'warmup': a.warmup, 'iters': a.iters, 'extern_cache': a.extern_cache, 'gogc': a.gogc or 'default'})
summary = {}
for task in bench.load_tasks(a.tasks.split(',')):
 name = task['name']
 ws = root/name
 if ws.exists(): sys.exit('refusing to overwrite '+str(ws))
 bench.make_workspace(str(ws), cache)
 jobs = bench.load(str(ref/name/'jobs.json'))
 (ws/'jobs.tsv').write_text(''.join('\t'.join(args)+'\n' for args in jobs))
 argv = [os.path.abspath(a.binary)] + (['-extern-cache'] if a.extern_cache else []) + [str(a.warmup), str(a.iters), 'jobs.tsv', 'batch']
 print('== '+name, flush=True)
 process = bench.run_proc(argv, str(ws), env, str(ws/'report.tsv'), str(ws/'driver.stderr'), 10)
 if process['exit'] != 0: sys.exit('Go driver failed: '+str(ws/'driver.stderr'))
 rows = []
 for line in (ws/'report.tsv').read_text().splitlines():
  if not line or line.startswith(('iter','#')): continue
  c = line.split('\t')
  rows.append({'iter': int(c[0]), 'kind': c[1], 'wall_ms': float(c[2]), 'cpu_ms': float(c[3]), 'gc_count': int(c[4]), 'gc_pause_ms': float(c[5]), 'rss_mb': float(c[6])})
 mismatch = []
 for j in range(len(jobs)):
  expected = bench.artifacts(task, str(ref/name/'java-cold/run-1'), j, len(jobs)==1)
  got = {k: str(ws/'batch'/f'{j}.{k}') for k in ('stdout','stderr','exit')}
  for k in task.get('compare',[]): got[k] = str(ws/k)
  mismatch += [{'job':j,'artifact':k} for k,p in expected.items() if not bench.same_file(p, got[k])]
 if mismatch: sys.exit('Java output mismatch: '+name+': '+repr(mismatch))
 measured = [r for r in rows if r['kind']=='measure']
 if len(measured) != a.iters: sys.exit('missing measured iterations')
 wall = statistics.median(r['wall_ms'] for r in measured)/1000
 java = bench.load(str(ref/name/'java-warm.json'))
 jwall = statistics.median(r['wall_ms'] for r in java['iterations'] if r['kind']=='measure')/1000
 summary[name] = {'median_s': wall, 'java_warm_s': jwall, 'speed_vs_java_warm':jwall/wall, 'median_cpu_s':statistics.median(r['cpu_ms'] for r in measured)/1000, 'peak_rss_mb':process['rss_mb'], 'all_outputs_match': True, 'iterations': rows}
 bench.save(str(root/'summary.json'), summary)
 print(f'   Go persistent {wall:.4f}s; Java warm {jwall:.4f}s; {jwall/wall:.3f}x; outputs MATCH', flush=True)
