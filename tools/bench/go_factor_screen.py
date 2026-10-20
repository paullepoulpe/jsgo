#!/usr/bin/env python3
"""Exploratory Go-v-Go factor screening; NOT a Java performance benchmark.

Every row is a fresh CLI under the existing real 4GiB/default-runtime budget,
byte-gated against an explicit pinned-stock-Java oracle for identical work.
Positive candidates require a subsequent live paired-JVM confirmation.
"""
import argparse
import hashlib
import json
from pathlib import Path
import random
import re
import statistics
import sys
import types

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
FACTORS_SHA256 = '6c5f7b60e83c4531be527d5d01a91a1e0bf307455a3ac45d133001a3ae922159'


def sha(path):
    h = hashlib.sha256()
    with Path(path).open('rb') as stream:
        for block in iter(lambda: stream.read(1 << 20), b''):
            h.update(block)
    return h.hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, allow_nan=False)+'\n')


def module_from_bytes(name, path, data):
    module = types.ModuleType(name)
    module.__file__ = str(path)
    exec(compile(data, str(path), 'exec'), module.__dict__)
    return module


def freeze(root):
    frozen, records = root/'runner-sources', {}
    def snapshot(name, expected=None):
        source = HERE/name
        data = source.read_bytes()
        digest = hashlib.sha256(data).hexdigest()
        if expected is not None and digest != expected:
            raise ValueError('helper changed; audit and repin explicitly: '+name)
        target = frozen/name
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_bytes(data)
        records[name] = {'source': str(source), 'snapshot': str(target), 'sha256': digest}
        return data, target
    data, path = snapshot('paired_factors.py', FACTORS_SHA256)
    factors = module_from_bytes('_screen_frozen_factor_model', path, data)
    for name in ('paired_cli.py', 'limited_cli.py', 'cgroup_exec.py'):
        snapshot(name, factors.HELPER_PINS[name])
    snapshot(Path(__file__).name)
    helper_path = frozen/'paired_cli.py'
    helper = module_from_bytes('_screen_frozen_engine', helper_path, helper_path.read_bytes())
    save(root/'source-manifest.json', records)
    return helper, factors, records


def balanced_rounds(names, rounds, seed):
    rng, rows = random.Random(seed), []
    while len(rows) < rounds:
        base = list(names); rng.shuffle(base)
        shifts = list(range(len(base))); rng.shuffle(shifts)
        for shift in shifts:
            if len(rows) == rounds:
                break
            rows.append({'round': len(rows)+1, 'variants': base[shift:]+base[:shift]})
    return rows


def work_identity(spec):
    return {key: spec[key] for key in ('argv', 'links', 'inputs', 'artifacts')} | {
        'absent_artifacts': spec.get('absent_artifacts', []),
        'builtin_externs': spec['same_work']['builtin_externs']}


def validate_oracle(spec, identity, oracle, metadata, pinned_jar):
    if metadata.get('jar_sha256', metadata.get('java_jar_sha256')) != pinned_jar:
        raise ValueError('oracle provenance must identify the pinned stock Java jar')
    if work_identity(metadata['spec']) != work_identity(spec) or metadata['input_manifest'] != identity:
        raise ValueError('stock Java oracle uses different argv, inputs, externs or output requirements')
    expected = {'stdout', 'stderr', 'exit'} | set(spec['artifacts']) | set(spec.get('absent_artifacts', []))
    if set(oracle) != expected:
        raise ValueError('oracle must cover ALL and only the task streams/exit/artifacts/absences')
    absent = set(spec.get('absent_artifacts', []))
    if absent & set(spec['artifacts']):
        raise ValueError('artifact cannot be both required and absent')
    for name, digest in oracle.items():
        if name in absent:
            if digest is not None:
                raise ValueError('oracle expected-absence entry is not null')
        elif not isinstance(digest, str) or not re.fullmatch('[0-9a-f]{64}', digest):
            raise ValueError('oracle missing an actual required artifact digest')
    if oracle['exit'] != hashlib.sha256(b'0\n').hexdigest():
        raise ValueError('a failed compiler task is not a valid timing oracle')


def stats(values):
    return {'median': statistics.median(values), 'min': min(values), 'max': max(values),
            'mean': statistics.mean(values), 'sample_stdev': statistics.stdev(values) if len(values) > 1 else None,
            'sample_variance': statistics.variance(values) if len(values) > 1 else None}


def summarize(rows, variants, control, expected):
    grouped = {v['name']: [r for r in rows if r['variant'] == v['name']] for v in variants}
    controls = {r['round']: r for r in grouped[control]}
    summary = {'schema': 'jsgo-go-factor-screen-v1', 'complete': len(rows) == expected,
               'scope': 'EXPLORATORY Go-v-Go only; positive candidates require live paired stock-JVM confirmation',
               'all_rows_match_stock_artifacts': all(r['gate'] == 'MATCH' for r in rows),
               'control': control, 'ratio_direction': 'Go variant / Go control; <1 means less time or RSS', 'variants': {}}
    for variant in variants:
        samples = grouped[variant['name']]
        if not samples:
            continue
        ratios = [{'round': r['round'], 'variant_row': r['row'], 'control_row': controls[r['round']]['row'],
                   'time_go_over_control': r['metrics']['wall_s']/controls[r['round']]['metrics']['wall_s'],
                   'rss_go_over_control': r['metrics']['peak_rss_mib']/controls[r['round']]['metrics']['peak_rss_mib']}
                  for r in samples if r['round'] in controls]
        summary['variants'][variant['name']] = {'samples': len(samples),
            'seconds': stats([r['metrics']['wall_s'] for r in samples]),
            'cpu_seconds': stats([r['metrics']['cpu_s'] for r in samples]),
            'peak_rss_mib': stats([r['metrics']['peak_rss_mib'] for r in samples]),
            'same_round_ratios': ratios,
            'time_go_over_control': stats([r['time_go_over_control'] for r in ratios]) if ratios else None,
            'rss_go_over_control': stats([r['rss_go_over_control'] for r in ratios]) if ratios else None}
    return summary


def main():
    p = argparse.ArgumentParser(description=__doc__)
    for name in ('spec', 'variants', 'oracle', 'image', 'out', 'control'):
        p.add_argument('--'+name, required=True)
    p.add_argument('--oracle-metadata', help='defaults to metadata.json beside explicit stock Java oracle')
    p.add_argument('--rounds', type=int, default=3)
    p.add_argument('--seed', type=int, default=20261020)
    p.add_argument('--readonly-bind', action='append', default=[])
    p.add_argument('--provenance', action='append', default=[], help='frozen build/source command/manifest/log to copy and hash; repeatable')
    a = p.parse_args()
    if a.rounds < 3 or not re.fullmatch(r'.+@sha256:[0-9a-f]{64}', a.image):
        p.error('require >=3 rounds and exact pinned image digest')
    root = Path(a.out).resolve()
    if root.exists():
        p.error('refusing to replace prior evidence')
    root.mkdir(parents=True)
    helper, factors, sources = freeze(root)
    variants = factors.load_variants(a.variants)
    if a.control not in {v['name'] for v in variants}:
        raise ValueError('control must name a configured Go variant')
    spec = json.loads(Path(a.spec).read_text())
    identity = helper.input_identity(spec)
    oracle_path = Path(a.oracle).resolve()
    oracle_metadata_path = Path(a.oracle_metadata).resolve() if a.oracle_metadata else oracle_path.parent/'metadata.json'
    oracle = json.loads(oracle_path.read_text())
    oracle_metadata = json.loads(oracle_metadata_path.read_text())
    validate_oracle(spec, identity, oracle, oracle_metadata, factors.PINNED_JAR)
    save(root/'stock-java-artifact-oracle.json', oracle)
    save(root/'stock-java-oracle-provenance.json', oracle_metadata)
    provenance = []
    for number, source_name in enumerate(a.provenance):
        source = Path(source_name).resolve()
        target = root/'build-provenance'/f'{number+1}-{source.name}'
        target.parent.mkdir(exist_ok=True)
        target.write_bytes(source.read_bytes())
        provenance.append({'source': str(source), 'snapshot': str(target), 'sha256': sha(target)})
    engines = {}
    for variant in variants:
        values = vars(a) | {'binary': variant['binary'], 'go_env': json.dumps(variant['go_env'], sort_keys=True),
                           'readonly_bind': sorted(set(a.readonly_bind+list(spec['links'].values())))}
        engines[variant['name']] = helper.Engine(argparse.Namespace(**values))
    first_engine = engines[variants[0]['name']]
    # Read-only image inspection fails if absent; do not let docker create pull
    # a missing image during a compiler screen.
    image_info = json.loads(first_engine.call('image', 'inspect', a.image, capture_output=True, text=True).stdout)[0]
    schedule = balanced_rounds([v['name'] for v in variants], a.rounds, a.seed)
    save(root/'schedule.json', schedule)
    row_plan = [{'row': number+1, 'round': block['round'], 'position': position+1, 'variant': name,
                 'work': str(root/f'round-{block["round"]}'/f'{position+1:02}-{name}')}
                for number, (block, position, name) in enumerate((block, position, name)
                    for block in schedule for position, name in enumerate(block['variants']))]
    save(root/'row-mapping.json', row_plan)
    binary_manifest = {v['binary']: v['binary_sha256'] for v in variants}
    def unchanged():
        factors.check_frozen(sources)
        if any(sha(path) != digest for path, digest in binary_manifest.items()) or helper.input_identity(spec) != identity:
            raise ValueError('compiler binary, helper source or task inputs changed')
    save(root/'metadata.json', vars(a) | {'schema': 'jsgo-go-factor-screen-v1', 'spec': spec, 'input_manifest': identity,
         'spec_sha256': sha(a.spec), 'variants_file_sha256': sha(a.variants), 'variants_resolved': variants,
         'source_manifest': sources, 'build_provenance': provenance,
         'build_provenance_status': 'caller-supplied frozen records' if provenance else 'not supplied; attach before any source-attributed claim',
         'image': {'reference': a.image, 'id': image_info['Id'], 'repo_digests': image_info.get('RepoDigests', []), 'os': image_info['Os'], 'architecture': image_info['Architecture']},
         'oracle': {'source': str(oracle_path), 'sha256': sha(oracle_path), 'metadata_source': str(oracle_metadata_path),
                    'metadata_sha256': sha(oracle_metadata_path), 'pinned_stock_jar_sha256': factors.PINNED_JAR},
         'scope': 'EXPLORATORY Go-v-Go screening, no Java timing/RSS comparison; live paired-JVM confirmation required',
         'schedule_policy': 'seeded Latin rotations; full V-round blocks position-balanced; partial counts differ <=1',
         'warmup_policy': 'no compiler warmup rows; all binaries hashed before first round, every row fresh, order balanced',
         'activation_evidence': 'flags are not dispatch proof; attach actual CLI activation probes separately',
         'memory_window': 'each fresh Go process lifetime wait4 peak; container peak additionally includes helper/page cache',
         'budget': 'actual cgroup4GiB, swap0,CPU2,affinity0,1 verified by frozen helpers; runtime defaults unchanged'})
    rows, variant_map = [], {v['name']: v for v in variants}
    before_screen = helper.outer_state()
    active_row = None
    try:
        for block in schedule:
            unchanged()
            before = helper.outer_state()
            for position, name in enumerate(block['variants']):
                variant, engine = variant_map[name], engines[name]
                index = len(rows)+1
                work = root/f'round-{block["round"]}'/f'{position+1:02}-{name}'
                active_row = row_plan[index-1]
                metrics = helper.go_request(engine, work, spec, engine.args)
                row = {'row': index, 'round': block['round'], 'position': position+1, 'variant': name,
                       'configuration_sha256': variant['configuration_sha256'], 'binary_sha256': variant['binary_sha256'],
                       'work': str(work), 'metrics': metrics, 'gate': 'MATCH' if metrics['digests'] == oracle else 'DIFF'}
                rows.append(row); save(root/'rows.json', rows)
                if row['gate'] != 'MATCH':
                    raise ValueError('strict stock Java artifact mismatch: '+name)
                print(json.dumps({'row': index, 'round': block['round'], 'variant': name, 'seconds': metrics['wall_s'],
                                  'peak_rss_mib': metrics['peak_rss_mib'], 'gate': row['gate']}), flush=True)
            after = helper.outer_state()
            save(root/f'round-{block["round"]}'/'outer-state.json', {'before': before, 'after': after})
            if any(helper.event_count(after.get('memory.events', ''), key) != helper.event_count(before_screen.get('memory.events', ''), key) for key in ('oom', 'oom_kill')):
                raise RuntimeError('outer OOM invalidates exploratory screen')
            unchanged()
            save(root/'summary.json', summarize(rows, variants, a.control, len(variants)*a.rounds))
    except BaseException as error:
        states = []
        for engine in engines.values():
            for cid in sorted(engine.containers):
                try:
                    state = json.loads(engine.call('inspect', cid, capture_output=True, text=True).stdout)[0]['State']
                    states.append({'id': cid, 'state': state})
                except Exception as inspection:
                    states.append({'id': cid, 'inspection_failure': str(inspection)})
        save(root/'failure.json', {'type': type(error).__name__, 'message': str(error), 'completed_rows': len(rows),
                                  'outer_state': helper.outer_state(), 'containers': states, 'active_row': active_row,
                                  'scope': 'FAILED exploratory screen; not a performance win'})
        raise
    finally:
        for engine in engines.values():
            engine.close()


if __name__ == '__main__':
    main()
