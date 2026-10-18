#!/usr/bin/env python3
"""Time a prebuilt jsgo-worker against single-invocation Java references.

Includes JSON transport and client decoding. Checks every response and generated
artifact outside the timed interval. Use the same CPU affinity for the Java
reference and this process; no builds should run during either measurement.
"""
import argparse
import base64
import json
import os
from pathlib import Path
import statistics
import subprocess
import sys
import time

sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True
import bench


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--binary', required=True)
    ap.add_argument('--reference', required=True)
    ap.add_argument('--out', required=True)
    ap.add_argument('--tasks', default='small-advanced,small-simple,small-whitespace,small-checks')
    ap.add_argument('--warmup', type=int, default=10)
    ap.add_argument('--iters', type=int, default=30)
    ap.add_argument('--cache', default=None)
    a = ap.parse_args()
    if a.warmup < 0 or a.iters < 1:
        ap.error('warmup must be nonnegative and iters positive')
    binary = str(Path(a.binary).resolve())
    reference, root = Path(a.reference).resolve(), Path(a.out).resolve()
    if root.exists():
        ap.error('refusing to overwrite existing output directory')
    root.mkdir(parents=True)
    env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP}
    env['LC_ALL'] = 'C.UTF-8'
    summary = {}
    for task in bench.load_tasks(a.tasks.split(',')):
        name = task['name']
        jobs = bench.load(str(reference/name/'jobs.json'))
        warm = bench.load(str(reference/name/'java-warm.json'))
        if not jobs or len(jobs) != 1 or not warm or warm.get('mismatch_vs_cold'):
            sys.exit('need a matching single-invocation Java reference: ' + name)
        java = [r['wall_ms'] for r in warm['iterations'] if r['kind'] == 'measure']
        if not java:
            sys.exit('missing measured Java iterations: ' + name)
        ws = root/name
        bench.make_workspace(str(ws), a.cache or bench.default_cache())
        artifacts = bench.artifacts(task, str(reference/name/'java-cold/run-1'), 0, True)
        expected = {k: Path(p).read_bytes() if Path(p).exists() else None
                    for k, p in artifacts.items()}
        rows, rss = [], []
        with (ws/'protocol.stderr').open('wb') as errors:
            p = subprocess.Popen([binary], cwd=ws, env=env, stdin=subprocess.PIPE,
                                 stdout=subprocess.PIPE, stderr=errors)
            try:
                for i in range(a.warmup+a.iters):
                    request = json.dumps({'id': i, 'args': jobs[0]}).encode()+b'\n'
                    start = time.perf_counter()
                    p.stdin.write(request)
                    p.stdin.flush()
                    response = json.loads(p.stdout.readline())
                    stdout = base64.b64decode(response['stdout_base64'], validate=True)
                    stderr = base64.b64decode(response['stderr_base64'], validate=True)
                    elapsed = (time.perf_counter()-start)*1000
                    if response['id'] != i:
                        sys.exit('worker response ID mismatch')
                    got = {'stdout': stdout, 'stderr': stderr,
                           'exit': (str(response['exit'])+'\n').encode()}
                    for key, want in expected.items():
                        if key not in got:
                            file = ws/key
                            got[key] = file.read_bytes() if file.exists() else None
                        if got[key] != want:
                            sys.exit(f'artifact mismatch: {name}, request {i}, {key}')
                    for line in Path(f'/proc/{p.pid}/status').read_text().splitlines():
                        if line.startswith('VmHWM:'):
                            rss.append(int(line.split()[1])/1024)
                    if i >= a.warmup:
                        rows.append(elapsed)
                p.stdin.close()
                if p.wait(timeout=10) != 0:
                    sys.exit('worker failed on shutdown: ' + name)
            finally:
                if p.poll() is None:
                    p.kill()
                    p.wait()
                p.stdout.close()
        median = statistics.median(rows)
        summary[name] = {
            'worker_protocol_median_ms': median,
            'java_warm_median_ms': statistics.median(java),
            'speed_vs_java_warm': statistics.median(java)/median,
            'time_ratio_java_over_go': statistics.median(java)/median,
            'worker_peak_rss_mb': max(rss),
            'java_peak_rss_mb': warm['process']['rss_mb'],
            'memory_ratio_java_over_go': warm['process']['rss_mb']/max(rss),
            'all_responses_match_java': True,
            'worker_warmup': a.warmup, 'worker_iters': a.iters,
            'measured_ms': rows, 'java_measure_ms': java,
            'java_warmup': warm['warmup'], 'java_iters': warm['iters'],
        }
        bench.save(str(root/'summary.json'), summary)
        print(f'{name}: worker {median:.3f}ms / {max(rss):.1f}MiB peak RSS, '
              f"Java {statistics.median(java):.3f}ms / {warm['process']['rss_mb']:.1f}MiB peak RSS; "
              f"Java/Go time {statistics.median(java)/median:.3f}x, RSS {warm['process']['rss_mb']/max(rss):.3f}x; "
              f'all {a.warmup+a.iters} requests MATCH', flush=True)
    bench.save(str(root/'provenance.json'), {
        'binary_sha256': bench.sha256_file(binary),
        'cpu_affinity': sorted(os.sched_getaffinity(0)),
        'reference': str(reference),
        'memory_unit': 'MiB (legacy _mb field names)',
        'memory_method': 'worker: /proc/PID/status VmHWM after each request; '
                         'Java: wait4 ru_maxrss for the full reference process; '
                         'both include warmup',
    })


if __name__ == '__main__':
    main()
