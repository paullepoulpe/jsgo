#!/usr/bin/env python3
"""Prepare/check or explicitly launch ONE frozen paired task in a quiet slot.

No benchmark helper, compiler argv or runtime options are changed. Empirical
shared-host admission is distinct from the actual4GiB cap on each process.
"""
import argparse
import json
import math
import os
from pathlib import Path
import subprocess
import sys

sys.dont_write_bytecode = True
import paired_cli as helper

RESULTS = Path('/workspace/jsgo-results/round4')
PREPARED = RESULTS/'next-fresh-paired-preparation-v1/commands.json'
PRIORITY = ['typed-full-n1', 'bigquery-n1', 'bundle-n2']
GIB = 1 << 30


def save(path, value):
    helper.save(path, value)


def check(plan):
    for path, digest in plan['pins'].items():
        if helper.sha(path) != digest:
            raise ValueError('frozen dependency changed: '+path)
    for task in plan['tasks']:
        spec = json.loads(Path(task['spec']).read_text())
        if helper.input_identity(spec) != json.loads(Path(task['input_manifest']).read_text()):
            raise ValueError('frozen compiler inputs changed')
        argv = task['command']
        if argv[argv.index('--spec')+1] != task['spec'] or argv[argv.index('--go-env')+1] != '{}':
            raise ValueError('prepared command no longer has exact task/default runtime')
    return plan


def prepare(out):
    if out.exists():
        raise ValueError('new preparation directory required')
    original = json.loads(PREPARED.read_text())
    records = {x['name']: x for x in original['records']}
    pins = dict(original['pins']) | {str(PREPARED): helper.sha(PREPARED), str(Path(__file__).resolve()): helper.sha(__file__)}
    tasks = []
    for name in PRIORITY:
        task = dict(records[name])
        diagnostic = Path(task['correctness_reference']['go'])/'result.json'
        observed = json.loads(diagnostic.read_text())
        peak = int(observed['cgroup']['memory.peak'])
        task['shared_host_admission'] = {
            'java_full_cap_bytes': 4*GIB,
            'historical_go_cgroup_peak_bytes': peak,
            'go_empirical_planning_reserve_bytes': max(GIB, math.ceil(1.5*peak)),
            'outer_margin_bytes': GIB//2,
            'warning': 'The Go reserve is empirical1.5x historical cgroup peak, minimum1GiB; not a hard Go cap or a guarantee. Each actual Java/Go cgroup remains4GiB. JVM pause stops CPU, retains resident memory. Any OOM/failure invalidates measurement.'}
        pins[str(diagnostic)] = helper.sha(diagnostic)
        pins[task['spec']] = task['spec_sha256']
        pins[task['input_manifest']] = task['input_manifest_sha256']
        pins[task['argv_record']] = task['argv_sha256']
        target = Path(task['command'][task['command'].index('--out')+1])
        if target.exists():
            raise ValueError('prepared fresh output already exists: '+str(target))
        tasks.append(task)
    plan = {'schema': 'jsgo-single-paired-launch-v1', 'status': 'PREPARED ONLY; no workloads launched',
        'priority': PRIORITY, 'tasks': tasks, 'pins': pins,
        'protocol': original['protocol'],
        'launch_rule': 'Run ONE explicitly selected task only, after root establishes and maintains a dedicated quiet slot. Check fresh output absence, immutable source/input pins, no other compiler/runtime workloads or containers, and shared-host empirical reserve. No helper or command modification.'}
    check(plan)
    out.mkdir(parents=True)
    save(out/'plan.json', plan)
    print(json.dumps({'status': 'PREPARED_NOT_RUN', 'plan': str(out/'plan.json'), 'priority': PRIORITY}))


def preflight(task, path):
    processes = []
    for proc in Path('/proc').iterdir():
        if not proc.name.isdecimal():
            continue
        try:
            executable = Path(os.readlink(proc/'exe').removesuffix(' (deleted)')).name
            if executable in {'go', 'compile', 'link', 'java', 'javac', 'node', 'protoc'} or executable.endswith('.test') or executable.startswith(('jsgo', 'protoc-')):
                processes.append({'pid': int(proc.name), 'executable': executable})
        except (FileNotFoundError, PermissionError, ProcessLookupError):
            pass
    env = {k: v for k, v in os.environ.items() if k not in helper.DOCKER_DROP}
    containers = subprocess.run(helper.DOCKER+['ps', '--format', '{{.ID}}\t{{.Names}}'],
        env=env, capture_output=True, text=True, check=True).stdout.splitlines()
    cg = Path('/sys/fs/cgroup')
    current, limit = [int((cg/name).read_text()) for name in ('memory.current', 'memory.max')]
    stat = dict((k, int(v)) for k, v in (line.split() for line in (cg/'memory.stat').read_text().splitlines()))
    clean = max(0, min(stat.get('file', 0)-stat.get('shmem', 0)-stat.get('file_dirty', 0)-stat.get('file_writeback', 0),
                       stat.get('inactive_file', 0)+stat.get('active_file', 0)))
    policy = task['shared_host_admission']
    required = policy['java_full_cap_bytes']+policy['go_empirical_planning_reserve_bytes']+policy['outer_margin_bytes']
    extra = limit-max(0, current-clean)-required
    proof = {'scope': 'Read-only planning admission; no reclaim/cache-dropping writes or runtime changes',
        'processes': processes, 'containers': containers, 'memory_current': current,
        'memory_max': limit, 'memory_stat': stat, 'clean_reclaimable_file_cache_only': clean,
        'all_other_outer_bytes_reserved': max(0, current-clean), 'shared_host_admission': policy,
        'additional_headroom_bytes': extra, 'admitted': not processes and not containers and extra >= 0}
    save(path, proof)
    if not proof['admitted']:
        raise RuntimeError('quiet/shared-host headroom admission failed; no benchmark launched')


def run(planpath, name, evidence):
    plan = check(json.loads(planpath.read_text()))
    task = next(x for x in plan['tasks'] if x['name'] == name)
    target = Path(task['command'][task['command'].index('--out')+1])
    if evidence.exists() or target.exists():
        raise ValueError('new launcher evidence and exact benchmark output paths required')
    evidence.mkdir(parents=True)
    save(evidence/'exact-command.json', task['command'])
    save(evidence/'selected-task.json', task)
    preflight(task, evidence/'quiet-shared-host-preflight.json')
    check(plan)
    if target.exists():
        raise ValueError('fresh benchmark output appeared before launch')
    before = helper.outer_state(); save(evidence/'outer-before.json', before)
    with (evidence/'stdout').open('wb') as stdout, (evidence/'stderr').open('wb') as stderr:
        done = subprocess.run(task['command'], stdout=stdout, stderr=stderr)
    after = helper.outer_state(); save(evidence/'outer-after.json', after)
    save(evidence/'exit.json', {'exit': done.returncode})
    if done.returncode or any(helper.event_count(before.get('memory.events', ''), k) != helper.event_count(after.get('memory.events', ''), k) for k in ('oom', 'oom_kill')):
        raise RuntimeError('paired task failed or outer OOM increased; evidence preserved, stop')
    summary = json.loads((target/'summary.json').read_text())
    if not summary['complete'] or not summary['warm_target_settled']:
        raise RuntimeError('incomplete or unsettled warm target; result remains inconclusive, stop')
    save(evidence/'result.json', {'status': 'COMPLETE', 'task': name, 'paired_summary': str(target/'summary.json'), 'outer_oom_delta': 0})


def main():
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest='mode', required=True)
    prep = sub.add_parser('prepare'); prep.add_argument('--out', required=True)
    chk = sub.add_parser('check'); chk.add_argument('--plan', required=True)
    launch = sub.add_parser('run'); launch.add_argument('--plan', required=True)
    launch.add_argument('--task', choices=PRIORITY, required=True)
    launch.add_argument('--evidence', required=True)
    a = p.parse_args()
    if a.mode == 'prepare': prepare(Path(a.out).resolve())
    elif a.mode == 'check':
        check(json.loads(Path(a.plan).read_text())); print('PASS: frozen commands/pins/input checks only; no workloads launched')
    else: run(Path(a.plan).resolve(), a.task, Path(a.evidence).resolve())


if __name__ == '__main__':
    main()
