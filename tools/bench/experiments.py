#!/usr/bin/env python3
"""Compare prebuilt jsgo variants against a stored Java reference.

Uses bench.py's tasks, process measurements, and compatibility gate. Variant
order alternates across rounds so neither binary always benefits from running
first. Every measured whole-program run must match the Java artifact digests.
No builds take place during timing.

  python3 -I tools/bench/experiments.py --reference /tmp/reference --out /tmp/compare \
      --variant baseline=/path/jsgo-old --variant candidate=/path/jsgo-new --runs 5
"""
import argparse
import json
import os
from pathlib import Path
import random
import statistics
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True
import bench


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument('--reference', required=True)
    ap.add_argument('--out', required=True)
    ap.add_argument('--variant', action='append', required=True, help='NAME=BINARY; no compilation during measurement')
    ap.add_argument('--tasks', default='typed-app,bundle,downlevel-whitespace,downlevel-simple,checks-only')
    ap.add_argument('--runs', type=int, default=5)
    ap.add_argument('--seed', type=int, default=20261018)
    ap.add_argument('--nice', type=int, default=10)
    ap.add_argument('--locale', default='C.UTF-8')
    ap.add_argument('--cache', default=None)
    a = ap.parse_args()
    if a.runs < 1:
        ap.error('--runs must be positive')
    a.cache = a.cache or bench.default_cache()
    variants = {}
    for value in a.variant:
        name, binary = value.split('=', 1)
        if not name or '/' in name or name in ('.', '..') or name in variants:
            ap.error('variant names must be unique simple directory names')
        binary = os.path.abspath(binary)
        if not os.access(binary, os.X_OK):
            ap.error('binary is not executable: ' + binary)
        variants[name] = binary
    root = Path(a.out).resolve()
    root.mkdir(parents=True, exist_ok=True)
    ref = Path(a.reference).resolve()
    runner = bench.Runner(a)
    metadata = bench.machine_info(runner, {'variants': {name: {'path': path, 'sha256': bench.sha256_file(path)}
                                                         for name, path in variants.items()}})
    metadata['cpu_affinity'] = sorted(os.sched_getaffinity(0))
    metadata['reference'] = str(ref)
    metadata['runs'] = a.runs
    metadata['seed'] = a.seed
    bench.save(str(root/'machine.json'), metadata)
    summary = {}
    rng = random.Random(a.seed)
    for task in bench.load_tasks(a.tasks.split(',')):
        name = task['name']
        reference = ref/name
        cold = bench.load(str(reference/'java-cold.json'))
        warm = bench.load(str(reference/'java-warm.json'))
        if not cold or not warm:
            sys.exit('missing Java reference for ' + name)
        expected = cold['runs'][0]['digests']
        jobs = bench.load(str(reference/'jobs.json'))
        warm_values = [i['wall_ms']/1000 for i in warm['iterations'] if i['kind']=='measure']
        if not warm_values or warm.get('mismatch_vs_cold'):
            sys.exit('invalid warm reference for ' + name)
        task_result = {label: [] for label in variants}
        tdirs = {}
        for label in variants:
            tdir = root/name/label
            tdir.mkdir(parents=True, exist_ok=True)
            if (tdir/'jsgo').exists():
                sys.exit('refusing to overwrite previous results: ' + str(tdir))
            (tdir/'java-cold').symlink_to(reference/'java-cold', target_is_directory=True)
            bench.save(str(tdir/'jobs.json'), jobs)
            tdirs[label] = str(tdir)
        order = list(variants)
        rng.shuffle(order)
        for index in range(1, a.runs+1):
            # AB/BA for pairs, rotated/reversed blocks for larger sets.
            shift = (index-1) % len(order)
            block = order[shift:] + order[:shift] if len(order)>2 else order[:]
            if index % 2 == 0:
                block.reverse()
            for label in block:
                binary = variants[label]
                print(f'== {name} / {label} / round {index}', flush=True)
                run = runner.measure_cold(task, jobs, tdirs[label], 'jsgo', index,
                                          lambda argv, binary=binary: [binary]+argv, runner.go_env)
                task_result[label].append(run)
                if index==1:
                    status = bench.gate(task, jobs, tdirs[label])
                    if status != 'MATCH':
                        sys.exit(f'compatibility gate failed: {name} / {label}: {status}')
                if run['digests'] != expected:
                    sys.exit(f'artifact digests changed: {name} / {label} / round {index}')
                bench.save(str(Path(tdirs[label])/'jsgo.json'), {'runs': task_result[label]})
        summary[name] = {'java_warm_s': statistics.median(warm_values), 'variants': {}}
        first_label = next(iter(variants))
        baseline = statistics.median([r['wall_s'] for r in task_result[first_label]])
        for label, runs in task_result.items():
            wall = [r['wall_s'] for r in runs]
            median = statistics.median(wall)
            summary[name]['variants'][label] = {
                'median_s': median, 'min_s': min(wall), 'max_s': max(wall),
                'median_cpu_s': statistics.median([r['cpu_s'] for r in runs]),
                'median_peak_rss_mb': statistics.median([r['rss_mb'] for r in runs]),
                'speed_vs_first_variant': baseline/median,
                'speed_vs_java_warm': summary[name]['java_warm_s']/median,
                'all_outputs_match': True,
            }
        bench.save(str(root/'summary.json'), summary)
        print(json.dumps({name: summary[name]}, indent=2), flush=True)


if __name__=='__main__':
    main()
