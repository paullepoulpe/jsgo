#!/usr/bin/env python3
"""Compare stock warm Java and jsgo at requested compiler thread counts.

--num_parallel_threads controls compiler prebuild helpers, not the process's
total thread count or available CPU cores. Run under a fixed CPU affinity and
quota. Historical Go binaries execute prebuild helpers sequentially even when this
flag exceeds one. New binaries may execute bounded
parallel prebuild helpers. This script compares fresh Go CLI processes with
stock warm Java. Both create fresh compilers and parse externs normally. No
builds run during timing; the Java driver is compiled before the sweep. Every result
must reproduce the stored Java reference, including generated artifacts.
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


def with_compiler_threads(argv, count):
    """Replace both thread-flag spellings once; preserve all other argv bytes."""
    result = []
    index = 0
    while index < len(argv):
        value = argv[index]
        if value == '--':
            # Do not reinterpret positional values after the option terminator.
            return result + [f'--num_parallel_threads={count}'] + argv[index:]
        if value == '--num_parallel_threads':
            if index+1 >= len(argv) or argv[index+1].startswith('--'):
                raise ValueError('existing --num_parallel_threads is missing its value')
            index += 2
        elif value.startswith('--num_parallel_threads='):
            index += 1
        else:
            result.append(value)
            index += 1
    return result + [f'--num_parallel_threads={count}']


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--binary', required=True)
    ap.add_argument('--reference', required=True)
    ap.add_argument('--out', required=True)
    ap.add_argument('--tasks', default='typed-app,bundle,downlevel-whitespace')
    ap.add_argument('--threads', default='1,2,4')
    ap.add_argument('--runs', type=int, default=5, help='Go CLI runs per setting')
    ap.add_argument('--warmup', type=int, default=12, help='Java warm-up iterations')
    ap.add_argument('--warm-iters', type=int, default=10)
    ap.add_argument('--cache', default=None)
    ap.add_argument('--nice', type=int, default=10)
    ap.add_argument('--locale', default='C.UTF-8')
    a = ap.parse_args()
    counts = [int(n) for n in a.threads.split(',')]
    if min(counts) < 1 or len(set(counts)) != len(counts) or a.runs < 1:
        ap.error('thread counts must be unique positive integers and runs positive')
    if a.warmup < 0 or a.warm_iters < 1:
        ap.error('warmup must be nonnegative and warm-iters positive')
    root, ref = Path(a.out).resolve(), Path(a.reference).resolve()
    if root.exists():
        ap.error('refusing to overwrite an existing output directory')
    a.out, a.cache = str(root), a.cache or bench.default_cache()
    binary = str(Path(a.binary).resolve())
    runner = bench.Runner(a)
    if bench.sha256_file(runner.jar) not in bench.read_pin()['CLOSURE_JAR_SHA256'].split():
        ap.error('stock Java jar does not match the pinned SHA-256')
    runner.build_batch()
    machine = bench.machine_info(runner, {'binary': binary, 'sha256': bench.sha256_file(binary)})
    machine.update({'cpu_affinity': sorted(os.sched_getaffinity(0)), 'mode': 'cli',
                    'thread_counts': counts, 'reference': str(ref),
                    'thread_control': '--num_parallel_threads; existing flags replaced; runtime CPU detection unchanged',
                    'gomaxprocs': 'runtime default; inherited GOMAXPROCS removed; compiler flag does not change it',
                    'java_compiler': 'pinned stock jar; fresh CommandLineRunner/compiler and normal extern parsing per invocation',
                    'java_classpath': runner.batch_cp,
                    'java_batch_source_sha256': bench.sha256_file(Path(bench.HERE)/'java/CliBatch.java'),
                    'ratio_direction': 'Java / Go; above 1 favors Go for time and resident memory',
                    'go_memory_window': 'median of fresh CLI process peak RSS values',
                    'java_memory_window': 'peak RSS of warm JVM including all warmup and measured iterations',
                    'go_prebuild': 'depends on supplied binary; bounded parallel helpers in updated port',
                    'memory_unit': 'MiB; field names ending _mb are retained for compatibility'})
    for name in ['cpu.max', 'cpuset.cpus.effective']:
        p = Path('/sys/fs/cgroup')/name
        if p.exists():
            machine[name] = p.read_text().strip()
    bench.save(str(root/'machine.json'), machine)
    summary = {}
    rng = random.Random(20261018)
    for task in bench.load_tasks(a.tasks.split(',')):
        name = task['name']
        jobs = bench.load(str(ref/name/'jobs.json'))
        cold = bench.load(str(ref/name/'java-cold.json'))
        if not jobs or not cold:
            sys.exit('missing Java reference for ' + name)
        dirs, go_runs = {}, {n: [] for n in counts}
        order = counts[:]
        rng.shuffle(order)
        for n in counts:
            td = root/f'threads-{n}'/name
            td.mkdir(parents=True)
            (td/'java-cold').symlink_to(ref/name/'java-cold', target_is_directory=True)
            args = [with_compiler_threads(job, n) for job in jobs]
            bench.save(str(td/'jobs.json'), args)
            dirs[n] = (td, args)
        for i in range(1, a.runs+1):
            for n in order if i % 2 else list(reversed(order)):
                td, args = dirs[n]
                print(f'== {name}: Go CLI threads={n}, run={i}', flush=True)
                run = runner.measure_cold(task, args, str(td), 'jsgo', i,
                                         lambda argv: [binary]+argv, runner.go_env)
                if run['digests'] != cold['runs'][0]['digests']:
                    sys.exit(f'Go output mismatch: {name}, threads={n}, run={i}')
                go_runs[n].append(run)
                bench.save(str(td/'jsgo.json'), {'runs': go_runs[n]})
        summary[name] = {}
        for n in order:
            td, args = dirs[n]
            print(f'== {name}: warm Java threads={n}', flush=True)
            warm = runner.measure_warm(task, args, str(td))
            bench.save(str(td/'java-warm.json'), warm)
            if warm['process']['exit'] != 0 or warm['mismatch_vs_cold']:
                sys.exit(f'Java output mismatch: {name}, threads={n}')
            java_rows = [r for r in warm['iterations'] if r['kind'] == 'measure']
            if len(java_rows) != a.warm_iters:
                sys.exit(f'missing Java measured iterations: {name}, threads={n}')
            jtime = statistics.median(r['wall_ms'] for r in java_rows)/1000
            jmem = warm['process']['rss_mb']
            gtime = statistics.median(r['wall_s'] for r in go_runs[n])
            gmem = statistics.median(r['rss_mb'] for r in go_runs[n])
            raw = {'go_cli_runs': go_runs[n]}
            half = len(java_rows)//2
            drift = (statistics.median(r['wall_ms'] for r in java_rows[half:]) /
                     statistics.median(r['wall_ms'] for r in java_rows[:half])-1) if half else None
            summary[name][str(n)] = {
                'go_median_s': gtime, 'java_warm_median_s': jtime,
                'go_peak_rss_mb': gmem, 'java_peak_rss_mb': jmem,
                'time_ratio_java_over_go': jtime/gtime,
                'memory_ratio_java_over_go': jmem/gmem,
                'java_warmup': a.warmup, 'java_measured': a.warm_iters,
                'java_measured_half_drift': drift,
                'java_measured_half_drift_pct': drift*100 if drift is not None else None,
                'java_iterations': warm['iterations'], 'all_outputs_match': True, **raw,
            }
            bench.save(str(root/'summary.json'), summary)
            print(f'   Java/Go time={jtime/gtime:.3f}x, RSS={jmem/gmem:.3f}x', flush=True)


if __name__ == '__main__':
    main()
