#!/usr/bin/env python3
"""Paired stock-warm-Java/fresh-Go CLI trials under real 4 GiB cgroups.

Requires a pinned container image with Python, the unchanged pinned compiler
jar, and precompiled CliPaired.class. It never builds or pulls an image. Every
compiler sees exactly the spec argv and input bytes. Idle Java is frozen while
Go runs; its resident pages remain charged to its separate 4 GiB cgroup.
"""
import argparse
import hashlib
import json
import math
import os
from pathlib import Path
import random
import shutil
import statistics
import subprocess
import sys
import time

HERE = Path(__file__).resolve().parent
DOCKER = ['docker', '--host=unix:///var/run/docker.sock']
DOCKER_DROP = {'DOCKER_HOST', 'DOCKER_CONTEXT', 'DOCKER_TLS', 'DOCKER_TLS_VERIFY', 'DOCKER_CERT_PATH'}
JAVA_DROP = {'JAVA_TOOL_OPTIONS', '_JAVA_OPTIONS', 'JDK_JAVA_OPTIONS', 'CLASSPATH'}
GO_DROP = {'GOGC', 'GOMEMLIMIT', 'GODEBUG', 'GOMAXPROCS', 'GOFLAGS'}


def save(path, obj):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(obj, indent=2, allow_nan=False)+'\n')


def sha(path):
    h = hashlib.sha256()
    with Path(path).open('rb') as stream:
        for block in iter(lambda: stream.read(1 << 20), b''):
            h.update(block)
    return h.hexdigest()


def outer_state():
    root = Path('/sys/fs/cgroup')
    return {name: (root/name).read_text().strip() for name in
            ('memory.events', 'memory.current', 'memory.max', 'cpu.max') if (root/name).exists()}


def event_count(text, name):
    return int(dict(line.split() for line in text.splitlines() if line.strip()).get(name, 0))


def stats(values):
    return {'median': statistics.median(values), 'min': min(values), 'max': max(values),
            'mean': statistics.mean(values), 'sample_stdev': statistics.stdev(values) if len(values) > 1 else None}


def ratio_ci(ratios, seed):
    if len(ratios) < 3:
        return None
    rng = random.Random(seed)
    logs = [math.log(x) for x in ratios]
    samples = sorted(math.exp(statistics.mean(rng.choices(logs, k=len(logs)))) for _ in range(10000))
    return {'geomean': math.exp(statistics.mean(logs)), 'ci95_low': samples[250],
            'ci95_high': samples[9749], 'independent_cohorts': len(logs)}


def warm_stability(rows):
    if len(rows) < 8:
        return {'stable': False, 'reason': 'need eight diagnostic rows'}
    old, new = rows[-8:-4], rows[-4:]
    drift = statistics.median(r['wall_s'] for r in new)/statistics.median(r['wall_s'] for r in old)-1
    jit = [r['jit_after_ms']-r['jit_before_ms'] for r in new]
    monitored = all(r['jit_before_ms'] >= 0 and r['jit_after_ms'] >= r['jit_before_ms'] for r in new)
    fraction = sum(jit)/(sum(r['wall_s'] for r in new)*1000) if monitored else None
    spread = max(r['wall_s'] for r in new)/min(r['wall_s'] for r in new)
    return {'stable': abs(drift) <= .05 and fraction is not None and fraction <= .02 and spread <= 1.20,
            'last4_vs_prior4_median_drift': drift, 'last4_jit_time_over_wall': fraction,
            'last4_max_over_min': spread, 'policy': 'fixed thresholds; extend warmup12 to at most30, retain all rows'}


class Engine:
    def __init__(self, args):
        self.args = args
        self.env = {k: v for k, v in os.environ.items() if k not in DOCKER_DROP}
        self.containers = set()

    def call(self, *argv, **kwargs):
        return subprocess.run(DOCKER+list(argv), env=self.env, check=True, **kwargs)

    def create(self, work, command, java=False):
        a = self.args
        bindings = set(a.readonly_bind + [str(HERE), str(Path(a.binary).resolve())])
        if java:
            bindings.update((a.java_home, a.jar, a.driver))
        argv = ['create', '-i', '--network=none', '--read-only', '--cpus=2', '--cpuset-cpus=0,1',
                '--memory=4294967296', '--memory-swap=4294967296', '--pids-limit=4096',
                '--tmpfs=/tmp:rw,nosuid,nodev,size=268435456', '--user='+str(os.getuid())+':'+str(os.getgid()),
                '--workdir='+str(work)]
        for path in sorted(bindings):
            argv += ['--mount', 'type=bind,src='+path+',dst='+path+',readonly']
        argv += ['--mount', 'type=bind,src='+str(work)+',dst='+str(work)]
        argv += ['--env', 'LC_ALL=C.UTF-8', a.image] + command
        cid = self.call(*argv, capture_output=True, text=True).stdout.strip()
        self.containers.add(cid)
        info = json.loads(self.call('inspect', cid, capture_output=True, text=True).stdout)[0]
        config = info['HostConfig']
        if config['Memory'] != 4294967296 or config['MemorySwap'] != 4294967296 or config['NanoCpus'] != 2000000000 or config['CpusetCpus'] != '0,1':
            raise ValueError('Docker did not install the exact memory/CPU constraints')
        save(work/'container-config.json', {'id': cid, 'image': info['Image'], 'host_config': {
            k: config[k] for k in ('Memory', 'MemorySwap', 'NanoCpus', 'CpusetCpus', 'PidsLimit', 'NetworkMode')}})
        return cid

    def remove(self, cid):
        self.call('rm', '-f', cid, stdout=subprocess.DEVNULL)
        self.containers.discard(cid)

    def close(self):
        for cid in list(self.containers):
            try:
                self.remove(cid)
            except subprocess.CalledProcessError:
                pass


def workspace(path, spec):
    path.mkdir(parents=True)
    for name, target in spec['links'].items():
        link = path/name
        link.parent.mkdir(parents=True, exist_ok=True)
        link.symlink_to(target, target_is_directory=True)
    for rel in spec['artifacts']+spec.get('absent_artifacts', []):
        target = path/rel
        if Path(rel).is_absolute() or '..' in Path(rel).parts:
            raise ValueError('artifact must be a safe relative path')
        target.parent.mkdir(parents=True, exist_ok=True)


def digests(work, spec, streams):
    paths = {name: streams/name for name in ('stdout', 'stderr', 'exit')}
    paths.update({name: work/name for name in spec['artifacts']})
    result = {name: sha(path) if path.is_file() else None for name, path in paths.items()}
    if any(value is None for value in result.values()):
        raise ValueError('required output missing: '+str(result))
    for name in spec.get('absent_artifacts', []):
        if (work/name).exists():
            raise ValueError('stock no-output artifact unexpectedly emitted: '+name)
        result[name] = None
    return result


def input_identity(spec):
    actual = []
    for item in spec['inputs']:
        path = Path(spec['links'][item['link']])/item['relative']
        digest = sha(path)
        if digest != item['sha256']:
            raise ValueError('input changed: '+str(path))
        actual.append({'link': item['link'], 'relative': item['relative'], 'bytes': path.stat().st_size, 'sha256': digest})
    return actual


def java_request(process, work, spec, number, kind):
    for rel in spec['artifacts']+spec.get('absent_artifacts', []):
        (work/rel).unlink(missing_ok=True)
    process.stdin.write(str(number)+'\t'+kind+'\n'); process.stdin.flush()
    line = process.stdout.readline().rstrip('\n')
    values = line.split('\t')
    if len(values) != 14 or values[:3] != ['RESULT', str(number), kind]:
        raise RuntimeError('Java protocol failed or process exceeded its cap: '+line)
    names = ['wall_ns', 'cpu_ns', 'gc_count', 'gc_ms', 'jit_before_ms', 'jit_after_ms',
             'rss_kib', 'hwm_kib', 'exit', 'cgroup_peak_bytes']
    row = dict(zip(names, map(int, values[3:13])))
    row.update({'wall_s': row['wall_ns']/1e9, 'cpu_s': row['cpu_ns']/1e9,
                'peak_rss_mib': row['hwm_kib']/1024, 'cgroup_events': values[13], 'kind': kind})
    events = values[13].replace(';', '\n')
    if row['exit'] != 0 or event_count(events, 'oom') or event_count(events, 'oom_kill') or row['cgroup_peak_bytes'] > 4294967296:
        raise RuntimeError('Java compilation/cgroup failure: '+str(row))
    row['digests'] = digests(work, spec, work/'protocol'/str(number))
    save(work/'protocol'/str(number)/'row.json', row)
    return row


def go_request(engine, work, spec, args):
    workspace(work, spec)
    env = {'LC_ALL': 'C.UTF-8', 'PATH': '/usr/local/bin:/usr/bin:/bin'} | json.loads(args.go_env)
    if GO_DROP & env.keys():
        raise ValueError('runtime GC/CPU overrides are outside this protocol')
    save(work/'request.json', {'argv': [args.binary]+spec['argv'], 'environment': env})
    cid = engine.create(work, ['python3', str(HERE/'limited_cli.py'), str(work/'request.json')])
    engine.call('start', '-a', cid, stdout=subprocess.DEVNULL)
    row = json.loads((work/'process.json').read_text())
    row['digests'] = digests(work, spec, work)
    save(work/'process.json', row)
    engine.remove(cid)
    if row['exit'] != 0 or event_count(row['cgroup']['memory.events'], 'oom') or event_count(row['cgroup']['memory.events'], 'oom_kill') or int(row['cgroup']['memory.peak']) > 4294967296:
        raise RuntimeError('Go compilation/cgroup failure')
    return row


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--spec', required=True)
    p.add_argument('--binary', required=True)
    p.add_argument('--jar', required=True)
    p.add_argument('--jar-sha256', required=True)
    p.add_argument('--driver', required=True, help='directory containing precompiled CliPaired.class')
    p.add_argument('--java-home', required=True)
    p.add_argument('--image', required=True, help='preinstalled image pinned by @sha256: digest; Python3 required')
    p.add_argument('--readonly-bind', action='append', default=[])
    p.add_argument('--out', required=True)
    p.add_argument('--cohorts', type=int, default=3)
    p.add_argument('--pairs', type=int, default=5)
    p.add_argument('--warmup', type=int, default=12)
    p.add_argument('--max-warmup', type=int, default=30)
    p.add_argument('--go-warmup', type=int, default=3, help='fresh unmeasured Go processes, each byte-gated')
    p.add_argument('--seed', type=int, default=20261020)
    p.add_argument('--go-env', default='{}')
    a = p.parse_args()
    if a.cohorts < 3 or a.pairs < 2 or not 8 <= a.warmup <= a.max_warmup <= 30:
        p.error('need >=3 cohorts, >=2 pairs, 8<=warmup<=max-warmup<=30')
    if '@sha256:' not in a.image or sha(a.jar) != a.jar_sha256:
        p.error('image and stock jar must match exact content pins')
    from bench import read_pin
    if a.jar_sha256 != read_pin()['CLOSURE_JAR_SHA256']:
        p.error('Java comparator must be the repository pinned stock jar')
    root = Path(a.out).resolve()
    if root.exists():
        p.error('refusing to overwrite prior evidence')
    root.mkdir()
    spec = json.loads(Path(a.spec).read_text())
    if any('\n' in value or '\r' in value for value in spec['argv']):
        p.error('newline in argv is unsupported')
    identity = input_identity(spec)
    a.readonly_bind = sorted(set(a.readonly_bind+list(spec['links'].values())))
    for item in [a.binary, a.jar, a.driver, a.java_home]:
        if not Path(item).is_absolute():
            p.error('paths must be absolute')
    binary_hash = sha(a.binary)
    save(root/'metadata.json', vars(a) | {'spec': spec, 'input_manifest': identity,
         'go_binary_sha256': binary_hash, 'driver_class_sha256': sha(Path(a.driver)/'CliPaired.class'),
         'driver_source_sha256': sha(HERE/'java/CliPaired.java'),
         'java_binary_sha256': sha(Path(a.java_home)/'bin/java'), 'java_flags': ['-Xmx3g'],
         'memory_windows': {'java': 'warm JVM lifetime process VmHWM, including warmup/JIT',
                            'go': 'fresh process wait4 lifetime peak; container memory.peak includes helper/page cache'},
         'paused_java': 'all JVM threads frozen during Go trials; resident memory remains charged, outer pressure recorded',
         'confidence_unit': 'independent restarted JVM cohort median ratios; not individual correlated iterations'})
    engine, rng, cohorts, reference = Engine(a), random.Random(a.seed), [], None
    try:
        for index in range(a.cohorts):
            if sha(a.binary) != binary_hash or input_identity(spec) != identity:
                raise ValueError('source inputs or executable changed')
            directory = root/f'cohort-{index+1}'
            jwork = directory/'java'
            workspace(jwork, spec)
            (jwork/'argv.txt').write_text('\n'.join(spec['argv'])+'\n')
            stderr = (jwork/'driver.stderr').open('w')
            cid = engine.create(jwork, ['python3', str(HERE/'cgroup_exec.py'), str(Path(a.java_home)/'bin/java'), '-Xmx3g', '-cp', a.jar+':'+a.driver,
                                        'CliPaired', str(jwork/'argv.txt')], java=True)
            process = subprocess.Popen(DOCKER+['start', '-a', '-i', cid], stdin=subprocess.PIPE,
                                       stdout=subprocess.PIPE, stderr=stderr, text=True, env=engine.env)
            ready = process.stdout.readline().rstrip('\n').split('\t')
            if len(ready) != 5 or ready[:3] != ['READY', '4294967296', '0'] or ready[4] not in ('0-1', '0,1'):
                raise RuntimeError('Java did not verify kernel memory/CPU cap: '+str(ready))
            quota, period = ready[3].split()
            if quota == 'max' or int(quota) != 2*int(period):
                raise RuntimeError('Java kernel CPU budget differs from two CPUs')
            warm, measured, go, orders, trial = [], [], [], [], 0
            before = outer_state()
            for n in range(a.max_warmup):
                trial += 1; row = java_request(process, jwork, spec, trial, 'warmup'); warm.append(row)
                if reference is None:
                    reference = row['digests']; save(root/'stock-java-oracle-digests.json', reference)
                if row['digests'] != reference:
                    raise ValueError('Java warmup output differs from fixed stock oracle')
                save(directory/'warmup.json', {'rows': warm, 'stability': warm_stability(warm)})
                if len(warm) >= a.warmup and warm_stability(warm)['stable']:
                    break
            engine.call('pause', cid, stdout=subprocess.DEVNULL)
            try:
                for n in range(a.go_warmup):
                    row = go_request(engine, directory/f'go-warmup-{n+1}', spec, a)
                    if row['digests'] != reference:
                        raise ValueError('Go unmeasured warmup differs from fixed stock oracle')
            finally:
                engine.call('unpause', cid, stdout=subprocess.DEVNULL)
            first = ['java', 'go']; rng.shuffle(first)
            for pair in range(a.pairs):
                order = first if pair % 2 == 0 else list(reversed(first)); orders.append(order)
                for runtime in order:
                    if runtime == 'java':
                        trial += 1; row = java_request(process, jwork, spec, trial, 'measure'); measured.append(row)
                    else:
                        engine.call('pause', cid, stdout=subprocess.DEVNULL)
                        try:
                            row = go_request(engine, directory/f'go-{pair+1}', spec, a); go.append(row)
                        finally:
                            engine.call('unpause', cid, stdout=subprocess.DEVNULL)
                    if row['digests'] != reference:
                        raise ValueError(runtime+' artifact bytes differ from pinned stock Java')
                    save(directory/'trials.json', {'java': measured, 'go': go, 'orders': orders})
            process.stdin.write('STOP\n'); process.stdin.flush(); process.stdin.close()
            if process.wait() != 0:
                raise RuntimeError('Java driver exit failure')
            stderr.close(); engine.remove(cid)
            after = outer_state()
            if event_count(after.get('memory.events', ''), 'oom_kill') != event_count(before.get('memory.events', ''), 'oom_kill'):
                raise RuntimeError('outer cgroup OOM invalidates the paired cohort')
            stability = warm_stability(warm)
            java_time, go_time = stats([x['wall_s'] for x in measured]), stats([x['wall_s'] for x in go])
            java_rss = measured[-1]['peak_rss_mib']; go_rss = statistics.median(x['peak_rss_mib'] for x in go)
            cohort = {'index': index+1, 'java_seconds': java_time, 'go_seconds': go_time,
                      'java_peak_rss_mib': java_rss, 'go_peak_rss_mib': go_rss,
                      'java_over_go_time': java_time['median']/go_time['median'],
                      'java_over_go_rss': java_rss/go_rss, 'warmup_stability': stability,
                      'outer_before': before, 'outer_after': after, 'all_outputs_match': True}
            cohorts.append(cohort); save(directory/'summary.json', cohort)
            save(root/'summary.json', {'cohorts': cohorts, 'complete': len(cohorts) == a.cohorts,
                 'time_ratio': ratio_ci([x['java_over_go_time'] for x in cohorts], a.seed),
                 'memory_ratio': ratio_ci([x['java_over_go_rss'] for x in cohorts], a.seed+1),
                 'warm_target_settled': all(x['warmup_stability']['stable'] for x in cohorts),
                 'status': 'inconclusive warmup' if not all(x['warmup_stability']['stable'] for x in cohorts) else 'measured',
                 'ratio_direction': 'Java / Go; greater than one favors Go'})
            print(json.dumps(cohort), flush=True)
    except BaseException as error:
        save(root/'failure.json', {'type': type(error).__name__, 'message': str(error),
                                 'completed_cohorts': len(cohorts), 'outer_state': outer_state()})
        raise
    finally:
        engine.close()


if __name__ == '__main__':
    main()
