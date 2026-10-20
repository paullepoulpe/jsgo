#!/usr/bin/env python3
"""Multiple Go factors, each paired with stock Java in shared warm-JVM cohorts.

Only name, binary and experiment environment vary. One task/spec/argv, resource
budget and stock Java target apply to every pair. No compiler builds, fixture
generation, downloads, workers or parsed-extern caches occur here.

--variants accepts [{"name":"stock-main","binary":"/abs/jsgo","go_env":{}},
...]. Use --control factor-off to report direct Go control/variant effects too.
Helper sources are hash-pinned below and copied into this run's evidence before
being imported/executed. Mutable source-tree helpers are never run silently.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import random
import re
import statistics
import subprocess
import sys
import types

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
REPO = HERE.parent.parent
PINNED_JAR = 'd6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6'
HELPER_PINS = {
    'paired_cli.py': '914a9c97eb5cb66fa786d01a2f43599b781327c89b51fe0e4fef3ab207fd0854',
    'limited_cli.py': '6b7b4d551cc9521d14b1b0088660cb111ea15aade76ef4f6300f6d11de310cf3',
    'cgroup_exec.py': '3b19b27b04a146b3793b917ab1eef30573ca92bcff06360acfb964ac7b6d025e',
    'java/CliPaired.java': '3e9690c404b3352f5240323a6ab854e4ec53349272b34361db1d069d4c67c1cd',
}


def sha(path):
    h = hashlib.sha256()
    with Path(path).open('rb') as stream:
        for block in iter(lambda: stream.read(1 << 20), b''):
            h.update(block)
    return h.hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, allow_nan=False)+'\n')


def load_variants(path):
    rows = json.loads(Path(path).read_text())
    if isinstance(rows, dict) and set(rows) == {'variants'}:
        rows = rows['variants']
    if not isinstance(rows, list) or not rows:
        raise ValueError('variants must be a nonempty array')
    result = []
    for row in rows:
        if not isinstance(row, dict) or set(row)-{'name', 'binary', 'go_env', 'sha256'} or not {'name', 'binary'} <= row.keys():
            raise ValueError('only name, binary, go_env and optional expected sha256 may vary')
        if not isinstance(row['name'], str) or not re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9_.-]{0,79}', row['name']):
            raise ValueError('unsafe or empty variant name')
        binary = Path(row['binary'])
        if not binary.is_absolute() or not binary.is_file():
            raise ValueError('variant binary must be an existing absolute file')
        env = row.get('go_env', {})
        if not isinstance(env, dict) or any(not isinstance(k, str) or not re.fullmatch(r'JSGO_EXPERIMENT_[A-Z0-9_]+', k)
                                            or not isinstance(v, str) or '\x00' in v for k, v in env.items()):
            raise ValueError('only string-valued JSGO_EXPERIMENT_* flags may vary; no runtime/locale/profile overrides')
        digest = sha(binary)
        if 'sha256' in row and row['sha256'] != digest:
            raise ValueError('variant executable does not match its supplied SHA-256')
        result.append({'name': row['name'], 'binary': str(binary.resolve()), 'go_env': env,
                       'binary_sha256': digest,
                       'configuration_sha256': hashlib.sha256(json.dumps({'binary_sha256': digest, 'go_env': env}, sort_keys=True).encode()).hexdigest()})
    if len({v['name'] for v in result}) != len(result):
        raise ValueError('variant names must be distinct')
    return result


def balanced_schedule(names, cohorts, repetitions, seed):
    """Random Latin rotations: position counts differ by at most one overall.

    A repetition visits every variant exactly once. Complete blocks of V rows
    visit every variant in every position; a partial block uses distinct shifts.
    Each variant's Java/Go pair direction alternates across its repetitions.
    """
    rng = random.Random(seed)
    first_direction = {name: rng.randrange(2) for name in names}
    rows = []
    while len(rows) < cohorts*repetitions:
        base = list(names); rng.shuffle(base)
        shifts = list(range(len(base))); rng.shuffle(shifts)
        for shift in shifts:
            if len(rows) == cohorts*repetitions:
                break
            sequence = len(rows)
            rows.append([{'variant': name, 'position': position+1,
                          'order': ['java', 'go'] if (sequence+first_direction[name]) % 2 == 0 else ['go', 'java']}
                         for position, name in enumerate(base[shift:]+base[:shift])])
    return [{'cohort': cohort+1,
             'repetitions': [{'repetition': repetition+1, 'pairs': rows[cohort*repetitions+repetition]}
                             for repetition in range(repetitions)]}
            for cohort in range(cohorts)]


def freeze_helpers(root, driver):
    frozen = root/'runner-sources'
    records = {}
    for name, expected in HELPER_PINS.items():
        source = HERE/name
        data = source.read_bytes()
        actual = hashlib.sha256(data).hexdigest()
        if actual != expected:
            raise ValueError('helper source changed; audit and repin before running: '+name)
        target = frozen/name
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_bytes(data)
        records[name] = {'source': str(source), 'snapshot': str(target), 'sha256': actual}
    source = Path(__file__).resolve()
    target = frozen/source.name
    target.write_bytes(source.read_bytes())
    records[source.name] = {'source': str(source), 'snapshot': str(target), 'sha256': sha(target)}
    classfile = Path(driver)/'CliPaired.class'
    target = frozen/'java-driver/CliPaired.class'
    target.parent.mkdir()
    target.write_bytes(classfile.read_bytes())
    records['java-driver/CliPaired.class'] = {'source': str(classfile), 'snapshot': str(target), 'sha256': sha(target)}
    # Execute the exact checked snapshot bytes, never a stale timestamp-based
    # .pyc cache or a helper from another checkout on sys.path.
    module = types.ModuleType('_paired_factors_frozen_helpers')
    module.__file__ = str(frozen/'paired_cli.py')
    exec(compile((frozen/'paired_cli.py').read_bytes(), module.__file__, 'exec'), module.__dict__)
    save(root/'source-manifest.json', records)
    return module, records, str(target.parent)


def check_frozen(records):
    for name, row in records.items():
        if sha(row['snapshot']) != row['sha256']:
            raise ValueError('frozen run helper changed: '+name)


def validate_ready(line):
    fields = line.rstrip('\n').split('\t')
    if len(fields) != 5 or fields[:3] != ['READY', '4294967296', '0'] or fields[4] not in ('0-1', '0,1'):
        raise ValueError('Java did not verify the exact kernel resource budget: '+repr(fields))
    quota, period = fields[3].split()
    if quota == 'max' or int(quota) != 2*int(period):
        raise ValueError('Java kernel CPU budget differs from two CPUs')
    return fields


def full_stats(helper, values):
    return helper.stats(values) | {'sample_variance': statistics.variance(values) if len(values) > 1 else None}


def row_stats(helper, rows):
    return {'seconds': full_stats(helper, [r['wall_s'] for r in rows]),
            'cpu_seconds': full_stats(helper, [r['cpu_s'] for r in rows]),
            'peak_rss_mib': full_stats(helper, [r['peak_rss_mib'] for r in rows])}


def summarize_cohort(helper, index, variants, pairs, warm, outer_before, outer_after):
    java = [p['rows']['java'] for p in pairs]
    common_hwm = max(r['peak_rss_mib'] for r in java)
    by_variant = {}
    for variant in variants:
        matched = [p for p in pairs if p['variant'] == variant['name']]
        js, gs = [p['rows']['java'] for p in matched], [p['rows']['go'] for p in matched]
        jstats, gstats = row_stats(helper, js), row_stats(helper, gs)
        by_variant[variant['name']] = {'pairs': len(matched), 'java': jstats, 'go': gstats,
            'java_over_go_time': jstats['seconds']['median']/gstats['seconds']['median'],
            'java_over_go_rss': common_hwm/gstats['peak_rss_mib']['median'],
            'individual_pair_time_ratios': full_stats(helper, [j['wall_s']/g['wall_s'] for j, g in zip(js, gs)]),
            'shared_java_lifetime_peak_rss_mib': common_hwm}
    drift = None
    if len(java) >= 8:
        drift = statistics.median(r['wall_s'] for r in java[-4:])/statistics.median(r['wall_s'] for r in java[:4])-1
    return {'index': index, 'variants': by_variant, 'all_outputs_match': True,
            'warmup_stability': helper.warm_stability(warm), 'measured_java_first4_vs_last4_median_drift': drift,
            'measured_java': row_stats(helper, java), 'outer_before': outer_before, 'outer_after': outer_after,
            'memory_window': 'one common final JVM lifetime VmHWM for every variant; each Go row is a fresh process peak'}


def summarize(helper, cohorts, variants, count, seed, control):
    result = {'complete': len(cohorts) == count, 'cohorts': cohorts,
              'warm_target_settled': all(c['warmup_stability']['stable'] for c in cohorts),
              'status': 'inconclusive warmup' if any(not c['warmup_stability']['stable'] for c in cohorts) else 'measured',
              'ratio_direction': 'Java / Go >1 favors Go; control / variant >1 favors variant', 'variants': {}}
    for offset, variant in enumerate(variants):
        name = variant['name']
        rows = [c['variants'][name] for c in cohorts]
        value = {'time_ratio': helper.ratio_ci([r['java_over_go_time'] for r in rows], seed+offset),
                 'memory_ratio': helper.ratio_ci([r['java_over_go_rss'] for r in rows], seed+1000+offset),
                 'cohort_median_seconds': {'java': full_stats(helper, [r['java']['seconds']['median'] for r in rows]),
                                           'go': full_stats(helper, [r['go']['seconds']['median'] for r in rows])}}
        if control:
            value['go_control'] = {'name': control,
                'time_ratio': helper.ratio_ci([c['variants'][control]['go']['seconds']['median']/c['variants'][name]['go']['seconds']['median'] for c in cohorts], seed+2000+offset),
                'memory_ratio': helper.ratio_ci([c['variants'][control]['go']['peak_rss_mib']['median']/c['variants'][name]['go']['peak_rss_mib']['median'] for c in cohorts], seed+3000+offset)}
        result['variants'][name] = value
    return result


def main():
    p = argparse.ArgumentParser(description=__doc__)
    for name in ('spec', 'variants', 'jar', 'driver', 'java-home', 'image', 'out'):
        p.add_argument('--'+name, required=True)
    p.add_argument('--cohorts', type=int, default=3)
    p.add_argument('--repetitions', type=int, default=5, help='pairs per variant in EACH independent JVM cohort')
    p.add_argument('--warmup', type=int, default=12)
    p.add_argument('--max-warmup', type=int, default=30)
    p.add_argument('--go-warmup', type=int, default=3, help='unmeasured byte-gated fresh processes PER variant/cohort')
    p.add_argument('--seed', type=int, default=20261020)
    p.add_argument('--control', help='variant name used for direct Go control/variant ratios')
    p.add_argument('--readonly-bind', action='append', default=[])
    a = p.parse_args()
    if a.cohorts < 3 or a.repetitions < 2 or a.go_warmup < 0 or not 8 <= a.warmup <= a.max_warmup <= 30:
        p.error('need >=3 cohorts, >=2 repetitions, nonnegative Go warmup, 8<=Java warmup<=max<=30')
    variants = load_variants(a.variants)
    if a.control and a.control not in {v['name'] for v in variants}:
        p.error('control must name a configured variant')
    if '@sha256:' not in a.image or sha(a.jar) != PINNED_JAR:
        p.error('require pinned container image and repository stock Java jar')
    pin_path = REPO/'tools/oracle/pin.env'
    pin = dict(line.strip().split('=', 1) for line in pin_path.read_text().splitlines() if line.strip() and not line.lstrip().startswith('#'))
    if PINNED_JAR not in pin['CLOSURE_JAR_SHA256'].strip('"').split():
        p.error('repository Java comparator pin changed; audit before repinning runner')
    for path in (a.jar, a.driver, a.java_home):
        if not Path(path).is_absolute():
            p.error('jar, driver and Java home must be absolute paths')
    root = Path(a.out).resolve()
    if root.exists():
        p.error('refusing to replace prior evidence')
    root.mkdir(parents=True)
    helper, sources, frozen_driver = freeze_helpers(root, a.driver)
    spec = json.loads(Path(a.spec).read_text())
    if any('\n' in value or '\r' in value for value in spec['argv']):
        raise ValueError('newlines in argv are unsupported')
    identity = helper.input_identity(spec)
    dependencies = helper.java_runtime_dependencies(a.java_home)
    schedule = balanced_schedule([v['name'] for v in variants], a.cohorts, a.repetitions, a.seed)
    variant_map = {v['name']: v for v in variants}
    # Each Engine gets one immutable binary/environment configuration. Never
    # mutate Engine.args under a live Java container or between Go requests.
    def engine_args(variant):
        values = vars(a) | {'binary': variant['binary'], 'go_env': json.dumps(variant['go_env'], sort_keys=True),
                           'driver': frozen_driver, 'readonly_bind': sorted(set(a.readonly_bind+list(spec['links'].values())))}
        return argparse.Namespace(**values)
    engines = {v['name']: helper.Engine(engine_args(v)) for v in variants}
    java_engine = helper.Engine(engine_args(variants[0]))
    binary_manifest = {v['binary']: v['binary_sha256'] for v in variants}
    def unchanged():
        check_frozen(sources)
        if any(sha(binary) != digest for binary, digest in binary_manifest.items()) or helper.input_identity(spec) != identity or helper.java_runtime_dependencies(a.java_home) != dependencies or sha(a.jar) != PINNED_JAR:
            raise ValueError('executable, input, jar or Java runtime configuration changed')
    save(root/'schedule.json', schedule)
    save(root/'metadata.json', vars(a) | {'schema': 'jsgo-paired-factors-v1', 'variants_resolved': variants,
         'spec': spec, 'input_manifest': identity, 'spec_sha256': sha(a.spec), 'variants_file_sha256': sha(a.variants),
         'helper_sources': sources, 'repository_pin_sha256': sha(pin_path), 'java_jar_sha256': PINNED_JAR,
         'java_binary_sha256': sha(Path(a.java_home)/'bin/java'), 'java_runtime_dependencies': dependencies,
         'java_flags': ['-Xmx3g'], 'driver_runtime': frozen_driver,
         'schedule_policy': 'randomized Latin rotations; complete V-row blocks position-balanced, incomplete counts differ <=1; AB/BA alternates per variant',
         'confidence_unit': 'independent restarted JVM cohort median ratios; no per-request pseudo-replication',
         'memory_windows': {'java': 'shared cohort JVM lifetime peak VmHWM including warmup/JIT and ALL variants paired requests',
                            'go': 'individual fresh process lifetime wait4 peak; per-variant median and full spread retained'},
         'paused_java': 'all JVM threads frozen during every Go process; resident pages remain charged, outer OOM invalidates cohort',
         'activation_evidence': 'environment flags alone do not establish factor dispatch; attach independent actual CLI activation probes',
         'same_work': 'one immutable task argv and input/extern identity for ALL variants; fresh stock Compiler and extern parse EVERY Java request; every stream/exit/artifact gated'})
    cohorts, reference, process, stderr = [], None, None, None
    try:
        for plan in schedule:
            unchanged()
            index = plan['cohort']
            directory, warm, pairs, trial = root/f'cohort-{index}', [], [], 0
            before = helper.outer_state()
            jwork = directory/'java'
            helper.workspace(jwork, spec)
            (jwork/'argv.txt').write_text('\n'.join(spec['argv'])+'\n')
            stderr = (jwork/'driver.stderr').open('w')
            command = ['python3', str(helper.HERE/'cgroup_exec.py'), str(Path(a.java_home)/'bin/java'), '-Xmx3g',
                       '-cp', a.jar+':'+frozen_driver, 'CliPaired', str(jwork/'argv.txt')]
            cid = java_engine.create(jwork, command, java=True)
            process = subprocess.Popen(helper.DOCKER+['start', '-a', '-i', cid], stdin=subprocess.PIPE,
                                       stdout=subprocess.PIPE, stderr=stderr, text=True, env=java_engine.env)
            save(directory/'ready.json', validate_ready(process.stdout.readline()))
            for unused in range(a.max_warmup):
                trial += 1
                row = helper.java_request(process, jwork, spec, trial, 'warmup'); warm.append(row)
                if reference is None:
                    reference = row['digests']; save(root/'stock-java-oracle-digests.json', reference)
                if row['digests'] != reference:
                    raise ValueError('Java warmup differs from fixed stock oracle')
                save(directory/'warmup.json', {'rows': warm, 'stability': helper.warm_stability(warm)})
                if len(warm) >= a.warmup and helper.warm_stability(warm)['stable']:
                    break
            # Warm each fresh binary's OS pages equally, without warming/reusing
            # a Go compiler. Every process still performs its complete task.
            java_engine.call('pause', cid, stdout=subprocess.DEVNULL)
            try:
                for variant in variants:
                    engine = engines[variant['name']]
                    for number in range(a.go_warmup):
                        row = helper.go_request(engine, directory/f'go-warmup-{variant["name"]}-{number+1}', spec, engine.args)
                        if row['digests'] != reference:
                            raise ValueError('Go warmup differs from fixed stock oracle: '+variant['name'])
            finally:
                java_engine.call('unpause', cid, stdout=subprocess.DEVNULL)
            # Freeze actual Java request IDs after adaptive warmup and before
            # measurement, including a one-to-one mapping for every Go pair.
            mapping = []
            for repetition in plan['repetitions']:
                for item in repetition['pairs']:
                    trial += 1
                    mapping.append(item | {'cohort': index, 'repetition': repetition['repetition'], 'java_request': trial,
                        'java_protocol': str(jwork/'protocol'/str(trial)),
                        'go_directory': str(directory/f'go-r{repetition["repetition"]}-{item["variant"]}')})
            save(directory/'row-mapping.json', mapping)
            for item in mapping:
                variant = variant_map[item['variant']]
                pair = item | {'binary_sha256': variant['binary_sha256'], 'configuration_sha256': variant['configuration_sha256'], 'rows': {}}
                pairs.append(pair)
                for runtime in item['order']:
                    if runtime == 'java':
                        row = helper.java_request(process, jwork, spec, item['java_request'], 'measure')
                    else:
                        java_engine.call('pause', cid, stdout=subprocess.DEVNULL)
                        try:
                            engine = engines[variant['name']]
                            row = helper.go_request(engine, Path(item['go_directory']), spec, engine.args)
                        finally:
                            java_engine.call('unpause', cid, stdout=subprocess.DEVNULL)
                    pair['rows'][runtime] = row
                    save(directory/'trials.json', pairs)
                    if row['digests'] != reference:
                        raise ValueError(runtime+' output differs from fixed Java oracle: '+variant['name'])
                pair['gate'] = 'MATCH'
                save(directory/'trials.json', pairs)
            process.stdin.write('STOP\n'); process.stdin.flush(); process.stdin.close()
            if process.wait() != 0:
                raise RuntimeError('Java driver exited with a failure')
            stderr.close(); stderr = None
            java_engine.remove(cid)
            after = helper.outer_state()
            if any(helper.event_count(after.get('memory.events', ''), key) != helper.event_count(before.get('memory.events', ''), key) for key in ('oom', 'oom_kill')):
                raise RuntimeError('outer cgroup OOM invalidates entire paired cohort')
            unchanged()
            cohort = summarize_cohort(helper, index, variants, pairs, warm, before, after)
            cohorts.append(cohort)
            save(directory/'summary.json', cohort)
            save(root/'summary.json', summarize(helper, cohorts, variants, a.cohorts, a.seed, a.control))
            print(json.dumps(cohort), flush=True)
    except BaseException as error:
        states = []
        for engine in [java_engine]+list(engines.values()):
            for container in sorted(engine.containers):
                try:
                    state = json.loads(engine.call('inspect', container, capture_output=True, text=True).stdout)[0]['State']
                    states.append({'id': container, 'state': state})
                except Exception as inspection:
                    states.append({'id': container, 'inspection_failure': str(inspection)})
        save(root/'failure.json', {'type': type(error).__name__, 'message': str(error), 'completed_cohorts': len(cohorts),
                                  'containers': states, 'outer_state': helper.outer_state(), 'claim': 'FAILED; never an infinite speedup'})
        raise
    finally:
        for engine in [java_engine]+list(engines.values()):
            engine.close()
        if process is not None and process.poll() is None:
            try:
                process.wait(timeout=10)
            except subprocess.TimeoutExpired:
                process.terminate()
        if stderr is not None:
            stderr.close()


if __name__ == '__main__':
    main()
