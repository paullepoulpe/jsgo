#!/usr/bin/env python3
"""Compare equivalent Go/Java JSON workers with stock byte references.

Both workers receive the same warm-up and measured requests, compiler thread
flag, files, JSON/base64 transport, and output checks. Each fresh worker's peak
RSS covers its complete, identical warm-up + measurement request window.
Java should be built with build_java_worker.py, enabling the equivalent bounded
immutable extern-template cache. This does not compare against stock warm Java.
"""
import argparse
import base64
import json
import os
from pathlib import Path
import selectors
import shutil
import statistics
import subprocess
import sys
import time

sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True
import bench


def thread_flags(args, count):
    """Replace both supported spellings without changing the rest of a job."""
    result = []
    skip = False
    for flag in args:
        if skip:
            skip = False
        elif flag == '--num_parallel_threads':
            skip = True
        elif not flag.startswith('--num_parallel_threads='):
            result.append(flag)
    return result + ['--num_parallel_threads=' + str(count)]


def process_status(pid):
    values = {}
    for line in Path(f'/proc/{pid}/status').read_text().splitlines():
        if line.startswith(('VmHWM:', 'Threads:')):
            key, value, *_ = line.split()
            values[key.rstrip(':')] = int(value)
    if 'VmHWM' not in values or 'Threads' not in values:
        raise RuntimeError('worker VmHWM/Threads unavailable')
    return values['VmHWM']/1024, values['Threads']


def percentile(samples, fraction):
    values = sorted(samples)
    index = (len(values)-1)*fraction
    lower = int(index)
    upper = min(lower+1, len(values)-1)
    return values[lower] + (values[upper]-values[lower])*(index-lower)


def sample_diagnostics(samples, drift_limit):
    median = statistics.median(samples)
    half = len(samples)//2
    drift = ((statistics.median(samples[half:]) / statistics.median(samples[:half])-1)*100
             if half else None)
    quarters = [statistics.median(samples[len(samples)*i//4:len(samples)*(i+1)//4])
                for i in range(4)] if len(samples) >= 4 else []
    return {'p05_ms': percentile(samples, 0.05), 'p95_ms': percentile(samples, 0.95),
            'median_absolute_deviation_ms': statistics.median(abs(value-median) for value in samples),
            'measurement_quarter_medians_ms': quarters, 'measurement_half_drift_pct': drift,
            # This is a visible diagnostic, not a proof of JVM convergence.
            'drift_check_limit_pct': drift_limit,
            'passes_drift_check': abs(drift) <= drift_limit if len(samples) >= 4 else None}


def validated_java_build(build_dir, cache):
    build = bench.load(str(build_dir/'build.json'))
    classes = build_dir/'classes'
    actual_names = {str(path.relative_to(classes)) for path in classes.rglob('*.class')}
    if actual_names != set(build['class_sha256']):
        raise RuntimeError('Java build contains missing or unexpected classes')
    for name, digest in build['class_sha256'].items():
        if bench.sha256_file(classes/name) != digest:
            raise RuntimeError('Java build class hash mismatch: '+name)
    pin = bench.read_pin()
    jar = Path(cache).resolve()/pin['CLOSURE_JAR']
    digest = bench.sha256_file(jar)
    if digest != build['jar_sha256'] or digest not in pin['CLOSURE_JAR_SHA256'].split():
        raise RuntimeError('Java build/reference jar hash mismatch')
    command = build['command']
    try:
        cp_index = command.index('-cp')
        classpath = [Path(entry).resolve() for entry in command[cp_index+1].split(os.pathsep)]
    except (ValueError, IndexError, TypeError):
        raise RuntimeError('Java build command lacks its expected classpath')
    if classpath != [classes.resolve(), jar.resolve()]:
        raise RuntimeError('Java build command does not execute the validated classes/jar')
    cache_flags = [flag for flag in command if flag.startswith('-Djsgo.bench.externCache=')]
    if cache_flags != ['-Djsgo.bench.externCache=true'] or command[-1] != 'CliWorker':
        raise RuntimeError('Java build command is not the equivalent cached CLI worker')
    if any(flag.startswith('-Djsgo.bench.verifyCache=') and flag != '-Djsgo.bench.verifyCache=false' for flag in command):
        raise RuntimeError('Java cache counters must be disabled during timing')
    launcher = shutil.which(command[0])
    if not launcher:
        raise RuntimeError('Java executable unavailable: '+command[0])
    launcher = str(Path(launcher).resolve())
    # Pin PATH resolution for the complete sweep, and record the native VM too.
    command = [launcher]+command[1:]
    runtime = {'executable': launcher, 'executable_sha256': bench.sha256_file(launcher)}
    vm = Path(launcher).parent.parent/'lib/server/libjvm.so'
    if vm.exists(): runtime['libjvm_sha256'] = bench.sha256_file(vm)
    release = Path(launcher).parent.parent/'release'
    if release.exists(): runtime['release'] = release.read_text()
    return build, command, runtime


def measure(command, ws, jobs, expected_jobs, warmup, iters, env, timeout, drift_limit):
    rows, peak, warmup_peak, thread_counts = [], 0, 0, []
    with (ws/'protocol.stderr').open('wb') as err:
        process = subprocess.Popen(command, cwd=ws, env=env, stdin=subprocess.PIPE,
                                   stdout=subprocess.PIPE, stderr=err)
        selector = selectors.DefaultSelector()
        selector.register(process.stdout, selectors.EVENT_READ)
        try:
            for iteration in range(warmup + iters):
                elapsed = 0
                for index, (args, expected) in enumerate(zip(jobs, expected_jobs)):
                    # Remove outputs outside the timed request so stale files
                    # cannot make a failed compile appear byte-compatible.
                    for artifact in expected:
                        if artifact not in ('stdout', 'stderr', 'exit'):
                            (ws/artifact).unlink(missing_ok=True)
                    request_id = iteration * len(jobs) + index
                    start = time.perf_counter_ns()
                    request = json.dumps({'id': request_id, 'args': args}).encode() + b'\n'
                    process.stdin.write(request)
                    process.stdin.flush()
                    if not selector.select(timeout):
                        raise RuntimeError(f'worker timeout at request {request_id}')
                    line = process.stdout.readline()
                    if not line:
                        raise RuntimeError(f'worker exited at request {request_id}; see {ws}/protocol.stderr')
                    response = json.loads(line)
                    got = {'stdout': base64.b64decode(response['stdout_base64'] or '', validate=True),
                           'stderr': base64.b64decode(response['stderr_base64'] or '', validate=True),
                           'exit': (str(response['exit'])+'\n').encode()}
                    elapsed += (time.perf_counter_ns() - start) / 1e6
                    if response['id'] != request_id:
                        raise AssertionError(f'response id mismatch: {response}')
                    # Gate every response and file, including every warm-up.
                    for artifact, want in expected.items():
                        actual = got.get(artifact)
                        if artifact not in got:
                            path = ws/artifact
                            actual = path.read_bytes() if path.exists() else None
                        if actual != want:
                            raise AssertionError(f'byte mismatch: request {request_id}, {artifact}, {ws}')
                    rss, threads = process_status(process.pid)
                    thread_counts.append(threads)
                    peak = max(peak, rss)
                    if iteration < warmup: warmup_peak = peak
                rows.append({'iter': iteration, 'kind': 'warmup' if iteration < warmup else 'measure',
                             'wall_ms': elapsed, 'peak_rss_mb': peak})
                bench.save(str(ws/'iterations.json'), rows)
            process.stdin.close()
            if process.wait(timeout=timeout) != 0:
                raise RuntimeError(f'worker failed on shutdown; see {ws}/protocol.stderr')
        finally:
            selector.close()
            if process.poll() is None:
                process.kill()
                process.wait()
            process.stdin.close()
            process.stdout.close()
    samples = [row['wall_ms'] for row in rows if row['kind'] == 'measure']
    return {'median_ms': statistics.median(samples), 'peak_rss_mb': peak,
            'warmup_peak_rss_mb': warmup_peak,
            'os_threads_after_response_min_max': [min(thread_counts), max(thread_counts)],
            'measured_ms': samples, 'iterations': rows, 'requests_checked': len(rows)*len(jobs),
            **sample_diagnostics(samples, drift_limit),
            'all_responses_and_artifacts_match_stock_java': True}


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--reference', action='append', required=True,
                    help='stock reference root; repeat for disjoint task sets')
    ap.add_argument('--cache', default=bench.default_cache())
    ap.add_argument('--out', required=True)
    ap.add_argument('--go-binary', required=True)
    ap.add_argument('--java-build', required=True, help='build_java_worker.py output directory')
    ap.add_argument('--tasks', default='small-whitespace,small-simple,small-advanced,small-checks')
    ap.add_argument('--threads', default='1,2,4')
    ap.add_argument('--warmup', type=int, default=100)
    ap.add_argument('--iters', type=int, default=30)
    ap.add_argument('--timeout', type=float, default=300)
    ap.add_argument('--stability-drift-pct', type=float, default=5,
                    help='flag measured-half drift above this threshold; never changes either worker request count')
    ap.add_argument('--gomaxprocs', type=int, help='explicit independent Go runtime CPU limit')
    args = ap.parse_args()
    if args.stability_drift_pct <= 0:
        ap.error('stability drift threshold must be positive')
    if args.warmup < 0 or args.iters < 1:
        ap.error('warmup must be nonnegative and iters positive')
    counts = [int(value) for value in args.threads.split(',')]
    if not counts or any(count < 1 for count in counts):
        ap.error('thread counts must be positive')
    if args.gomaxprocs is not None and args.gomaxprocs < 1:
        ap.error('gomaxprocs must be positive')
    root = Path(args.out).resolve()
    if root.exists():
        ap.error('refusing to overwrite '+str(root))
    root.mkdir(parents=True)
    binary = Path(args.go_binary).resolve()
    java_build = Path(args.java_build).resolve()
    build, java_command, java_runtime = validated_java_build(java_build, args.cache)
    env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP+bench.JAVA_ENV_DROP}
    env['LC_ALL'] = 'C.UTF-8'
    if args.gomaxprocs:
        env['GOMAXPROCS'] = str(args.gomaxprocs)
    commands = {'go': [str(binary)], 'java': java_command}
    metadata = {'warmup': args.warmup, 'iters': args.iters, 'compiler_threads': counts,
                'gomaxprocs': args.gomaxprocs or 'runtime default', 'commands': commands,
                'affinity': sorted(os.sched_getaffinity(0)), 'go_binary_sha256': bench.sha256_file(binary),
                'java_build': build, 'java_runtime': java_runtime,
                'stability_drift_limit_pct': args.stability_drift_pct, 'references': [str(Path(p).resolve()) for p in args.reference],
                'memory_window': 'process high-water RSS through all warmup and measured requests; MiB',
                'ratio_direction': 'Java / Go; >1 favors Go for both time and memory',
                'timing_window': 'request JSON encoding, pipe write, response read/JSON/base64 decoding; artifact checks excluded'}
    bench.save(str(root/'metadata.json'), metadata)
    summary = {}
    pair_index = 0
    for task in bench.load_tasks(args.tasks.split(',')):
        name = task['name']
        refs = [Path(p).resolve()/name for p in args.reference if (Path(p)/name/'jobs.json').exists()]
        if len(refs) != 1:
            raise RuntimeError(f'expected exactly one stock reference for {name}, got {refs}')
        reference = refs[0]
        original_jobs = bench.load(str(reference/'jobs.json'))
        expected_jobs = []
        for index in range(len(original_jobs)):
            files = bench.artifacts(task, str(reference/'java-cold/run-1'), index, len(original_jobs)==1)
            expected_jobs.append({key: Path(path).read_bytes() if Path(path).exists() else None
                                  for key, path in files.items()})
        summary[name] = {}
        for count in counts:
            jobs = [thread_flags(job, count) for job in original_jobs]
            result = {}
            # Alternate process order so one implementation is not always first.
            order = ['go', 'java'] if pair_index % 2 == 0 else ['java', 'go']
            pair_index += 1
            for runtime in order:
                ws = root/name/f'threads-{count}'/runtime
                if ws.exists():
                    raise RuntimeError('refusing to overwrite '+str(ws))
                bench.make_workspace(str(ws), args.cache)
                bench.save(str(ws/'jobs.json'), jobs)
                print(f'{name}, compiler_threads={count}, {runtime}: {args.warmup}+{args.iters}', flush=True)
                result[runtime] = measure(commands[runtime], ws, jobs, expected_jobs, args.warmup,
                                          args.iters, env, args.timeout, args.stability_drift_pct)
                if result[runtime]['passes_drift_check'] is False:
                    print(f'  WARNING {runtime}: measured-half drift '
                          f'{result[runtime]["measurement_half_drift_pct"]:+.1f}% exceeds '
                          f'{args.stability_drift_pct}%; inspect raw samples before claiming steady state', flush=True)
            result['java_over_go_time'] = result['java']['median_ms']/result['go']['median_ms']
            result['java_over_go_peak_rss'] = result['java']['peak_rss_mb']/result['go']['peak_rss_mb']
            summary[name][str(count)] = result
            bench.save(str(root/'summary.json'), summary)
            print(f'  Java/Go time {result["java_over_go_time"]:.3f}x; '
                  f'peak RSS {result["java_over_go_peak_rss"]:.3f}x; all outputs MATCH', flush=True)
    bench.save(str(root/'complete.json'), {'complete': True})


if __name__ == '__main__':
    main()
