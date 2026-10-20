#!/usr/bin/env python3
"""Prepare or explicitly run the five-variant, two-task Go-only screen.

Preparation/check perform source/input validation only: no Docker or compiler.
Execution needs the finished candidate's caller-supplied SHA and a quiet slot.
"""
import argparse
import json
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


def prepare(out):
    if out.exists():
        raise ValueError('new preparation directory required')
    if screen.sha(STOCK) != STOCK_SHA:
        raise ValueError('frozen stock5de binary changed')
    if not CANDIDATE.parent.is_dir():
        raise ValueError('candidate destination directory not prepared')
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
    provenance = [RESULTS/'main-5de/stock-build.json', RESULTS/'sass-slab-5de-v3/source-manifest.json']
    paths += provenance
    plan = {'schema': 'jsgo-sass-slab-screen-plan-v1', 'status': 'PREPARED ONLY; candidate hash pending',
        'stock': str(STOCK), 'stock_sha256': STOCK_SHA, 'candidate': str(CANDIDATE),
        'candidate_sha256': None, 'image': IMAGE, 'rounds': 3, 'seed': 20261020,
        'variants': [{'name': 'stock-5de', 'binary': str(STOCK), 'go_env': {}, 'sha256': STOCK_SHA}]+
            [{'name': name, 'binary': str(CANDIDATE), 'go_env': {SASS: a, SLAB: b}, 'sha256': None}
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
        'candidate_sha256': None, 'tasks': [t['name'] for t in tasks], 'rows': 30}))


def run(planpath, out, digest):
    plan = check(json.loads(planpath.read_text()))
    if not re.fullmatch('[0-9a-f]{64}', digest) or screen.sha(plan['candidate']) != digest:
        raise ValueError('finished candidate must match supplied SHA256')
    if out.exists():
        raise ValueError('new run evidence directory required')
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
    chk = sub.add_parser('check'); chk.add_argument('--plan', required=True)
    launch = sub.add_parser('run'); launch.add_argument('--plan', required=True)
    launch.add_argument('--out', required=True); launch.add_argument('--candidate-sha256', required=True)
    a = p.parse_args()
    if a.mode == 'prepare': prepare(Path(a.out).resolve())
    elif a.mode == 'check':
        check(json.loads(Path(a.plan).read_text())); print('PASS: preparation paths/pins/inputs/oracles; NO workloads launched')
    else: run(Path(a.plan).resolve(), Path(a.out).resolve(), a.candidate_sha256)


if __name__ == '__main__':
    main()
