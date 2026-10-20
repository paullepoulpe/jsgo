#!/usr/bin/env python3
"""Prepare or explicitly run the five-variant, two-task Go-only screen.

Preparation/check perform source/input validation only: no Docker or compiler.
Execution needs the finished candidate's caller-supplied SHA and a quiet slot.
"""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys

sys.dont_write_bytecode = True
import go_factor_screen as screen
import paired_cli as helper
import paired_factors as factors

HERE = Path(__file__).resolve().parent
RESULTS = Path('/workspace/jsgo-results/round4')
STOCK = RESULTS/'bin/jsgo-main-5de6e335'
STOCK_SHA = '7ba61a7dc3b8b9880a02f9ef8801c088efe5ea455207186ef6c90d7a4e741989'
CANDIDATE = RESULTS/'sass-slab-5de-v3/bin/jsgo-sass-slab-5de'
IMAGE = 'python@sha256:70729b46c69b4f1e97c4822c1af3df53a1476cf5ddc6c087c0c10bc3a5678c2f'
SASS = 'JSGO_EXPERIMENT_SASS_EXPERT_VALUE'
SLAB = 'JSGO_EXPERIMENT_NODE_SLAB_102'
CONFIGS = [('integrated-off', '0', '0'), ('sass-only', '1', '0'),
           ('slab-only', '0', '1'), ('sass-and-slab', '1', '1')]
TASKS = [
    ('bundle', 'enterprise-specs-5de-bundle-t1-v1/bundle-advanced.json',
     'paired-main5de-bundle-advanced-t1-v1/stock-java-oracle-digests.json',
     'paired-main5de-bundle-advanced-t1-v1/metadata.json',
     'main5de-bundle-t1-bytegate-v1/result.json'),
    ('typed', 'enterprise-specs-5de-next-gates-v1/typed-full.json',
     'latest-main-typed-full-stock-oracle/stock-java-oracle-digests.json',
     'latest-main-typed-full-stock-oracle/metadata.json',
     'main5de-typed-full-bytegate-v3/result.json'),
]


def save(path, value):
    screen.save(path, value)


def check(plan):
    if plan['schema'] != 'jsgo-sass-slab-screen-plan-v1':
        raise ValueError('unsupported plan')
    for path, expected in plan['pins'].items():
        if screen.sha(path) != expected:
            raise ValueError('prepared source/reference changed: '+path)
    for task in plan['tasks']:
        spec = json.loads(Path(task['spec']).read_text())
        identity = helper.input_identity(spec)
        if identity != json.loads(Path(task['input_manifest']).read_text()):
            raise ValueError('prepared task inputs changed')
        screen.validate_oracle(spec, identity,
            json.loads(Path(task['oracle']).read_text()),
            json.loads(Path(task['oracle_metadata']).read_text()), factors.PINNED_JAR)
    return plan


def prepare(out, candidate=CANDIDATE, source_manifest=None, build_result=None, extra_provenance=()):
    if out.exists():
        raise ValueError('new preparation directory required')
    if screen.sha(STOCK) != STOCK_SHA:
        raise ValueError('frozen stock5de binary changed')
    candidate = Path(candidate).resolve()
    if candidate != CANDIDATE and (source_manifest is None or build_result is None or not extra_provenance):
        raise ValueError('custom candidate requires explicit source manifest, build result and build/equivalence provenance')
    source_manifest = Path(source_manifest or RESULTS/'sass-slab-5de-v3/source-manifest.json').resolve()
    build_result = Path(build_result or RESULTS/'sass-slab-5de-v3/steps/cli-build/result.json').resolve()
    build = json.loads(build_result.read_text())
    if build.get('exit') != 0 or Path(build['binary']).resolve() != candidate:
        raise ValueError('successful build result must identify this exact candidate path')
    candidate_sha = build['binary_sha256']
    if screen.sha(candidate) != candidate_sha:
        raise ValueError('candidate does not match frozen build result')
    if json.loads(source_manifest.read_text()).get('base') != '5de6e3359349d70e60c98788e5e62b78cf1cb0e7':
        raise ValueError('candidate source manifest must identify the stock5de base')
    if screen.sha(HERE/'paired_factors.py') != screen.FACTORS_SHA256:
        raise ValueError('screen factor helper pin changed')
    for name, expected in factors.HELPER_PINS.items():
        if screen.sha(HERE/name) != expected:
            raise ValueError('factor helper pin changed: '+name)
    paths = [Path(__file__), STOCK] + [HERE/name for name in
        ['go_factor_screen.py', 'paired_factors.py', 'paired_cli.py', 'limited_cli.py', 'cgroup_exec.py']]
    tasks = []
    out.mkdir(parents=True)
    for name, spec_path, oracle, metadata, diagnostic in TASKS:
        spec_path, oracle, metadata, diagnostic = [RESULTS/x for x in (spec_path, oracle, metadata, diagnostic)]
        spec = json.loads(spec_path.read_text())
        identity = helper.input_identity(spec)
        screen.validate_oracle(spec, identity, json.loads(oracle.read_text()),
                               json.loads(metadata.read_text()), factors.PINNED_JAR)
        if spec['same_work']['compiler_threads'] != 1:
            raise ValueError('both screen tasks must retain their matched N1')
        ip = out/(name+'-input-manifest.json'); save(ip, identity)
        ap = out/(name+'-argv.json'); save(ap, spec['argv'])
        diag = json.loads(diagnostic.read_text())
        tasks.append({'name': name, 'spec': str(spec_path), 'oracle': str(oracle),
            'oracle_metadata': str(metadata), 'input_manifest': str(ip), 'argv_record': str(ap),
            'compiler_threads': 1, 'rows': 15,
            'historical_diagnostic_stock_go_wall_s': diag['wall_s'],
            'diagnostic_only_15row_arithmetic_minutes': 15*diag['wall_s']/60,
            'estimate_warning': 'Busy historical correctness clocks, not a benchmark prediction or Java comparison.'})
        paths += [spec_path, oracle, metadata, diagnostic, ip, ap]
    provenance = [RESULTS/'main-5de/stock-build.json', source_manifest, build_result]+[Path(p).resolve() for p in extra_provenance]
    paths += provenance
    plan = {'schema': 'jsgo-sass-slab-screen-plan-v1', 'status': 'PREPARED ONLY; finished candidate pinned to build result',
        'stock': str(STOCK), 'stock_sha256': STOCK_SHA, 'candidate': str(candidate),
        'candidate_sha256': candidate_sha, 'candidate_source_manifest': str(source_manifest),
        'candidate_build_result': str(build_result), 'image': IMAGE, 'rounds': 3, 'seed': 20261020,
        'variants': [{'name': 'stock-5de', 'binary': str(STOCK), 'go_env': {}, 'sha256': STOCK_SHA}]+
            [{'name': name, 'binary': str(candidate), 'go_env': {SASS: a, SLAB: b}, 'sha256': candidate_sha}
             for name, a, b in CONFIGS],
        'tasks': tasks, 'provenance': list(map(str, provenance)),
        'pins': {str(p): screen.sha(p) for p in paths},
        'scope': 'Exploratory fresh Go-v-Go only; no Java timing/RSS ratios or performance claim.',
        'budget': 'Each fresh process actual4GiB/swap0/CPU2/affinity0,1; default runtime GC and CPU policy.',
        'order': screen.balanced_rounds(['stock-5de']+[x[0] for x in CONFIGS], 3, 20261020),
        'measurement': 'Three fresh rows per variant/task; retain time/CPU/RSS min/median/max/variance; same-round ratios vs stock and integrated-off. No CI from these three exploratory rows.',
        'activation': 'Flags are not proof of activation; attach the candidate actual CLI probes/gates separately.',
        'launch_conditions': 'Root quiet slot after tests/CLI gates; explicit final candidate SHA required. Stop all builds/profilers/generators. Strict fixed Java artifacts and metadata every row.'}
    save(out/'plan.json', plan)
    check(plan)
    print(json.dumps({'status': 'PREPARED_NOT_RUN', 'plan': str(out/'plan.json'),
        'candidate_sha256': candidate_sha, 'tasks': [t['name'] for t in tasks], 'rows': 30}))



def quiet_preflight(path):
    """Runtime-only read-only checks; preparation never contacts Docker."""
    workloads = []
    for proc in Path('/proc').iterdir():
        if not proc.name.isdecimal():
            continue
        try:
            comm = (proc/'comm').read_text().strip()
            executable = Path(os.readlink(proc/'exe').removesuffix(' (deleted)')).name
            if executable in {'go', 'compile', 'link', 'java', 'javac', 'node', 'protoc'} or executable.endswith('.test') or executable.startswith(('jsgo', 'protoc-')):
                workloads.append({'pid': int(proc.name), 'comm': comm, 'executable': executable})
        except (FileNotFoundError, PermissionError, ProcessLookupError):
            pass
    env = {k: v for k, v in os.environ.items() if k not in helper.DOCKER_DROP}
    docker = subprocess.run(helper.DOCKER+['ps', '--format', '{{.ID}}\t{{.Names}}'],
                            env=env, capture_output=True, text=True, check=True)
    containers = docker.stdout.splitlines()
    cg = Path('/sys/fs/cgroup')
    current, limit = [int((cg/name).read_text()) for name in ('memory.current', 'memory.max')]
    stat = dict((k, int(v)) for k, v in (line.split() for line in (cg/'memory.stat').read_text().splitlines()))
    clean = max(0, min(stat.get('file', 0)-stat.get('shmem', 0)-stat.get('file_dirty', 0)-stat.get('file_writeback', 0),
                       stat.get('inactive_file', 0)+stat.get('active_file', 0)))
    extra = limit-max(0, current-clean)-5*(1 << 30)
    proof = {'scope': 'Read-only quiet/headroom admission, no reclaim/cache-dropping writes',
             'workloads': workloads, 'containers': containers, 'memory_current': current,
             'memory_max': limit, 'memory_stat': stat, 'clean_file_cache_only': clean,
             'full_task_cap': 4*(1 << 30), 'outer_margin': 1 << 30, 'additional_headroom': extra}
    proof['admitted'] = not workloads and not containers and extra >= 0
    save(path, proof)
    if not proof['admitted']:
        raise RuntimeError('quiet/headroom preflight failed; no screen launched')


def run(planpath, out, digest):
    plan = check(json.loads(planpath.read_text()))
    if not re.fullmatch('[0-9a-f]{64}', digest) or screen.sha(plan['candidate']) != digest or plan['candidate_sha256'] != digest:
        raise ValueError('finished candidate must match supplied SHA256')
    if out.exists():
        raise ValueError('new run evidence directory required')
    build = json.loads(Path(plan['candidate_build_result']).read_text())
    if build.get('exit') != 0 or build['binary_sha256'] != digest or Path(build['binary']).resolve() != Path(plan['candidate']).resolve():
        raise ValueError('candidate no longer matches its frozen successful build provenance')
    variants = [dict(v, sha256=digest) if v['binary'] == plan['candidate'] else v
                for v in plan['variants']]
    out.mkdir(parents=True)
    variantpath = out/'variants.json'; save(variantpath, variants)
    factors.load_variants(variantpath)
    save(out/'prepared-plan.json', plan)
    save(out/'launch.json', {'candidate_sha256': digest, 'plan_sha256': screen.sha(planpath),
                            'scope': plan['scope']})
    for task in plan['tasks']:
        check(plan)
        quiet_preflight(out/(task['name']+'-quiet-preflight.json'))
        command = [sys.executable, str(HERE/'go_factor_screen.py'), '--spec', task['spec'],
            '--variants', str(variantpath), '--control', 'stock-5de', '--oracle', task['oracle'],
            '--oracle-metadata', task['oracle_metadata'], '--image', plan['image'],
            '--rounds', '3', '--seed', str(plan['seed']), '--out', str(out/task['name'])]
        for path in plan['provenance']+[str(planpath), str(out/'launch.json')]:
            command += ['--provenance', path]
        save(out/(task['name']+'-command.json'), command)
        with (out/(task['name']+'.stdout')).open('wb') as stdout, (out/(task['name']+'.stderr')).open('wb') as stderr:
            done = subprocess.run(command, stdout=stdout, stderr=stderr)
        save(out/(task['name']+'-exit.json'), {'exit': done.returncode})
        if done.returncode:
            raise RuntimeError(task['name']+' screen failed; stop and preserve evidence')
        rows = json.loads((out/task['name']/'rows.json').read_text())
        meta = json.loads((out/task['name']/'metadata.json').read_text())
        save(out/task['name']/'integrated-off-control-summary.json',
             screen.summarize(rows, meta['variants_resolved'], 'integrated-off', 15))
    save(out/'result.json', {'status': 'PASS', 'rows': 30, 'all_rows_match_stock_java_artifacts': True,
                             'scope': plan['scope']})


def main():
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest='mode', required=True)
    prep = sub.add_parser('prepare'); prep.add_argument('--out', required=True)
    prep.add_argument('--candidate', default=str(CANDIDATE))
    prep.add_argument('--candidate-source-manifest')
    prep.add_argument('--candidate-build-result')
    prep.add_argument('--candidate-provenance', action='append', default=[],
                      help='repeat build commands/source-equivalence proof paths; frozen and copied into both runs')
    chk = sub.add_parser('check'); chk.add_argument('--plan', required=True)
    launch = sub.add_parser('run'); launch.add_argument('--plan', required=True)
    launch.add_argument('--out', required=True); launch.add_argument('--candidate-sha256', required=True)
    a = p.parse_args()
    if a.mode == 'prepare': prepare(Path(a.out).resolve(), a.candidate, a.candidate_source_manifest, a.candidate_build_result, a.candidate_provenance)
    elif a.mode == 'check':
        check(json.loads(Path(a.plan).read_text())); print('PASS: preparation paths/pins/inputs/oracles; NO workloads launched')
    else: run(Path(a.plan).resolve(), Path(a.out).resolve(), a.candidate_sha256)


if __name__ == '__main__':
    main()
