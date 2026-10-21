"""Audit failed v4 evidence without producing performance statistics or ratios.

No Java, Docker, benchmark imports or large-artifact/JAR hashing. Only small
source/identity files, one stock class entry, and the failed request's three
protocol streams are hashed. Existing metric/JIT rows remain in their raw paths.
"""
import csv
import hashlib
import json
from pathlib import Path
from urllib.parse import unquote, urlparse
import zipfile

OUT = Path(__file__).resolve().parent
ROOT = OUT.parents[1]
RUN = ROOT / 'results' / OUT.name
read = lambda p: json.loads(Path(p).read_text())
sha = lambda p: hashlib.sha256(Path(p).read_bytes()).hexdigest()
counter = lambda text: {k: int(v) for k, v in (line.split() for line in text.replace(';', '\n').splitlines())}
OOM = ('oom', 'oom_kill', 'oom_group_kill')
m, failure, stopped = [read(RUN / name) for name in ('metadata.json', 'failure.json', 'stopped-proof.json')]
manifest, oracle = m['manifest'], read(RUN / 'oracle.json')
assert failure['status'] == 'REJECTED no acceptedperformanceclaim' and failure['type'] == 'AssertionError'
assert not (RUN / 'summary.json').exists()
assert m['baseline'] == 'STOCK' and m['candidate'] == 'accumulator'
assert (m['cohorts'], m['pairs'], m['min_warmup'], m['max_warmup'], m['request_timeout'], m['total_timeout']) == (5, 20, 60, 150, 300, 21600)
assert m['runner_sha256'] == '9adb1c481ea19998d45a7e75a87e7670c1d5efafb381f2138c4d8f64bc18730d' == sha(ROOT / 'bench/run-v4.py')
assert m['manifest_sha256'] == sha(ROOT / 'build/accumulator/manifest.json')
assert manifest == read(ROOT / 'build/accumulator/manifest.json') and manifest['status'] == 'PASS_BUILD_VALIDATION'
for key in ('spec', 'oracle', 'oracle_metadata', 'runtime_grant'):
    assert sha(m[key]) == m[key + '_sha256']
spec, grant = read(m['spec']), read(m['runtime_grant'])
assert read(m['oracle']) == oracle and len(oracle) == 6
assert set(oracle) == {'stdout', 'stderr', 'exit', *spec['artifacts'], *spec.get('absent_artifacts', [])}
for key in ('baseline', 'candidate', 'cohorts', 'pairs', 'min_warmup', 'max_warmup', 'seed', 'request_timeout', 'total_timeout', 'benchmark_protocol', 'manifest_sha256', 'runner_sha256', 'spec_sha256'):
    assert grant[key] == m[key], key
assert Path(grant['out']).resolve() == RUN
for name, expected in m['helpers'].items():
    assert sha(RUN / 'helper-snapshots' / name) == expected
assert sha(ROOT / 'bench/CliFactors.java') == manifest['driver']['source_sha256'] == m['driver_source_sha256']
for name, expected in manifest['driver']['classes'].items():
    assert sha(Path(m['driver']) / name) == expected
families = manifest['compiler_source_families']
assert families == ['com/google/javascript/jscomp/FunctionArgumentInjector']
assert manifest['unmodified_compiler_classes_outside_overlays'] is True
for name, v in manifest['variants'].items():
    if name == 'STOCK':
        assert v['overlay'] is None and v['classes'] == {}
    else:
        assert v['compiled_source_families'] == families and v['javac_flags'] == manifest['javac_flags']
        assert sha(v['build_proof']) == v['build_proof_sha256'] == manifest['compiler_build_receipt_sha256'][name]
        receipt = read(v['build_proof'])
        assert receipt['status'] == 'PASS' and receipt['exit'] == 0 and receipt['classes'] == v['classes']
with zipfile.ZipFile(m['jar']) as jar:
    stock_classes = {family: hashlib.sha256(jar.read(family + '.class')).hexdigest() for family in families}

def pressure_changes(p):
    a, b = p['before'], p['after']
    ea, eb = counter(a['memory.events.local']), counter(b['memory.events.local'])
    sa, sb = counter(a['memory.stat']), counter(b['memory.stat'])
    changes = {'local_high': eb.get('high', 0) - ea.get('high', 0), 'local_max': eb.get('max', 0) - ea.get('max', 0)}
    for key in ('pgscan_direct', 'pgsteal_direct', 'pgmajfault'):
        changes[key] = sb[key] - sa[key] if key in sa and key in sb else None
    assert a['memory.pressure'] is None and b['memory.pressure'] is None
    changes.update(PSI_some_total_us=None, PSI_full_total_us=None)
    assert changes == p['changes'] and all(v is None or v >= 0 for v in changes.values())
    for key in OOM:
        assert counter(b['memory.events']).get(key, 0) - counter(a['memory.events']).get(key, 0) == 0
    return changes

indices, counts, container_ids, anomalies, incomplete = [], [], [], [], []
for number in (1, 2):
    cohort = RUN / f'cohort-{number}'
    warm, pairs = read(cohort / 'warmup.json'), read(cohort / 'pairs.json')
    assert len(pairs) == (20 if number == 1 else 7)
    assert all(len(warm['rows'][name]) == 60 and warm['stability'][name]['stable'] for name in ('STOCK', 'accumulator'))
    for i, pair in enumerate(pairs, 1):
        assert pair['pair'] == i and set(pair['rows']) == {'STOCK', 'accumulator'}
        assert sorted(pair['order']) == ['STOCK', 'accumulator']
        if i > 1:
            assert pair['order'] == pairs[i - 2]['order'][::-1]
    for name in ('STOCK', 'accumulator'):
        work, variant = cohort / name, manifest['variants'][name]
        ready, config, owned = [read(work / file) for file in ('ready.json', 'container-config.json', 'owned-driver.json')]
        caps, host = ready['caps'], config['host_config']
        assert caps[:3] == ['READY', '4294967296', '0'] and caps[4] in ('0-1', '0,1')
        quota, period = map(int, caps[3].split()); assert quota == 2 * period
        assert host['Memory'] == host['MemorySwap'] == 4294967296 and host['NanoCpus'] == 2000000000
        assert host['CpusetCpus'] == '0,1' and host['NetworkMode'] == 'none'
        assert config['id'] == owned['container'] and owned['pid'] == owned['pgid'] and Path(owned['workdir']) == work
        container_ids.append(config['id'])
        assert ready['VM']['runtime_version'].startswith('21.') and ready['VM']['max_heap_bytes'] == 3221225472
        cp = ([variant['overlay']] if variant['overlay'] else []) + [m['jar'], m['driver']]
        assert read(work / 'jvm-command.json')['command'] == ['python3', str(RUN / 'helper-snapshots/cgroup_exec.py'), str(Path(m['java_home']) / 'bin/java'), '-Xmx3g', '-cp', ':'.join(cp), 'CliFactors', str(work / 'argv.txt'), str(work / 'classes.txt')]
        assert (work / 'argv.txt').read_text().splitlines() == spec['argv']
        assert set(ready['loaded_target_classes']) == {f.replace('/', '.') for f in families}
        for family in families:
            origin = ready['loaded_target_classes'][family.replace('/', '.')]
            assert Path(unquote(urlparse(origin['location']).path)).resolve() == Path(variant['overlay'] or m['jar']).resolve()
            assert origin['sha256'] == (variant['classes'][family + '.class'] if variant['overlay'] else stock_classes[family])
        dirs = sorted((work / 'protocol').iterdir(), key=lambda p: int(p.name))
        assert [int(p.name) for p in dirs] == list(range(1, len(dirs) + 1))
        rows = [read(p / 'row.json') for p in dirs if (p / 'row.json').exists()]
        assert len(rows) == (80 if number == 1 else 68 if name == 'STOCK' else 67)
        last_peak = last_hwm = 0
        for p in dirs:
            request = int(p.name)
            pressure = read(p / 'parent-pressure.json')
            changes = pressure_changes(pressure)
            assert pressure['zero_required'] == (request > 60)
            if not (p / 'row.json').exists():
                assert (number, name, request) == (2, 'accumulator', 68)
                assert changes == {'local_high': 0, 'local_max': 387, 'pgscan_direct': 14879, 'pgsteal_direct': 14879, 'pgmajfault': 0, 'PSI_some_total_us': None, 'PSI_full_total_us': None}
                stream_hashes = {stream: sha(p / stream) for stream in ('stdout', 'stderr', 'exit')}
                assert all(stream_hashes[k] == oracle[k] for k in stream_hashes)
                incomplete.append({'cohort': number, 'variant': name, 'request': request, 'pair': 8, 'protocol': str(p),
                                   'pressure': pressure, 'three_protocol_streams_match_oracle': True,
                                   'stream_sha256': stream_hashes, 'six_digest_check': 'NOT_PERFORMED_BEFORE_ABORT',
                                   'metrics_JIT_row': 'NOT_SAVED; no raw RESULT transcript retained',
                                   'final_artifacts_present_unhashed': {key: (work / key).is_file() for key in spec['artifacts']}})
                indices.append({'cohort': number, 'variant': name, 'kind': 'measure', 'request': request, 'pair': 8,
                                'retention': 'FAILED_PRESSURE_NO_METRIC_ROW', 'six_recorded_digests_match': False, 'metric_row': '', 'protocol': str(p)})
                continue
            row = read(p / 'row.json')
            assert row['parent_pressure'] == pressure and row['variant'] == name and row['request'] == request
            assert row['kind'] == ('warmup' if request <= 60 else 'measure')
            assert row['digests'] == oracle and row['exit'] == 0
            assert all(counter(row['cgroup_events']).get(k, 0) == 0 for k in OOM)
            assert row['cgroup_peak_bytes'] <= 4294967296 and row['cgroup_current_bytes'] <= 4294967296
            assert row['cgroup_peak_bytes'] >= last_peak
            if row['peak_rss_mib'] < last_hwm:
                previous = rows[request - 2]
                later = next((r for r in rows[request:] if r['hwm_kib'] > previous['hwm_kib']), None)
                assert request <= 60 and later is not None
                anomalies.append({'cohort': number, 'variant': name, 'request': request, 'prior_request': request - 1,
                                  'before_kib': previous['hwm_kib'], 'after_kib': row['hwm_kib'],
                                  'later_request': later['request'], 'later_hwm_kib': later['hwm_kib'],
                                  'policy': 'EXISTING_WARMUP_ONLY_LATER_STRICTLY_EXCEEDED; raw retained rows unchanged'})
            last_peak, last_hwm = row['cgroup_peak_bytes'], row['peak_rss_mib']
            if request <= 60:
                assert row == warm['rows'][name][request - 1]
                state = 'COMPLETE_WARMUP_ROW'
            else:
                assert all(v is None or v == 0 for v in changes.values())
                pair_index = request - 60
                if pair_index <= len(pairs):
                    assert row == pairs[pair_index - 1]['rows'][name]
                    state = 'COMPLETE_PAIR_ROW'
                else:
                    assert (number, name, request) == (2, 'STOCK', 68)
                    state = 'COMPLETE_UNPAIRED_REQUEST'
            indices.append({'cohort': number, 'variant': name, 'kind': row['kind'], 'request': request,
                            'pair': request - 60 if request > 60 else '', 'retention': state,
                            'six_recorded_digests_match': True, 'metric_row': str(p / 'row.json'), 'protocol': str(p)})
        counts.append({'cohort': number, 'variant': name, 'warm_rows': 60, 'measured_metric_rows': len(rows) - 60,
                       'complete_measured_pairs': len(pairs), 'protocol_directories': len(dirs),
                       'full_recorded_six_digest_rows': len(rows)})
assert len(incomplete) == 1 and len(indices) == 296 and sum(r['six_recorded_digests_match'] for r in indices) == 295
assert read(RUN / 'progress.json')['completed_cohorts'] == [read(RUN / 'cohort-1/cohort-summary.json')]
assert not (RUN / 'cohort-2/cohort-summary.json').exists()
assert all(not (RUN / f'cohort-{n}').exists() for n in (3, 4, 5))
assert len(set(container_ids)) == len(stopped['owned_cleanup']) == 4
assert {r['container'] for r in stopped['owned_cleanup']} == set(container_ids)
assert all(r['removed'] is True and r['driver_exit'] == 0 and not r.get('fallback_removed') for r in stopped['owned_cleanup'])
assert all(stopped[key] == [] for key in ('cleanup_errors', 'remaining_owned_containers', 'active_compilers', 'remaining_owned_driver_processes'))
assert all(stopped['outer_OOM_delta'][k] == 0 for k in OOM)
assert all(counter(stopped['outer_after']['memory.events']).get(k, 0) - counter(stopped['outer_before']['memory.events']).get(k, 0) == 0 for k in OOM)
rejection = read(OUT / 'frozen-reviewer-rejection/review.json')
assert rejection['status'] == 'FAIL independent Java-factor review' and rejection['error'] == 'failure evidence present: no accepted result'
report = {'status': 'AUDIT_COMPLETED_FAILED_RUN_REMAINS_REJECTED', 'results': str(RUN),
          'planned_protocol': m['benchmark_protocol'], 'confirmation_complete': False, 'accepted_performance_claim': False,
          'performance_statistics': 'This failure audit computes none; subsequent partial-descriptive.json was prepared only after root closed NO_RESTART before receiving partial metrics',
          'completed_cohorts': 1, 'partial_cohorts': 1, 'unstarted_cohorts': [3, 4, 5],
          'counts': counts, 'full_metric_digest_rows': 295, 'warm_rows': 240, 'measured_metric_rows': 55,
          'complete_pairs': 27, 'extra_unpaired_STOCK_measurement': {'cohort': 2, 'request': 68, 'pair': 8},
          'failed_request': incomplete[0], 'manifest_runner_grant_spec_driver_class_origin_resource_checks': 'PASS_WITHIN_AUDIT_SCOPE',
          'recorded_six_digest_checks': 'ALL_295_FULL_ROWS_EQUAL_PINNED_ORACLE', 'telemetry_anomalies': anomalies,
          'cleanup': {'owned_JVMs': 4, 'all_driver_exit_codes_zero': True, 'all_owned_containers_removed': True,
                      'fallback_cleanup': False, 'cleanup_errors': [], 'remaining_owned_processes_containers': [],
                      'outer_OOM_delta': stopped['outer_OOM_delta']},
          'PSI': 'UNAVAILABLE_not_zero', 'frozen_reviewer': {'status': rejection['status'], 'reason': rejection['error'], 'exit': 1},
          'small_identity_hashes': {'metadata': sha(RUN / 'metadata.json'), 'failure': sha(RUN / 'failure.json'),
                                    'stopped_proof': sha(RUN / 'stopped-proof.json'), 'runner': m['runner_sha256'],
                                    'reviewer': sha(ROOT / 'bench/review-results-v4.py'), 'manifest': m['manifest_sha256']},
          'limitations': ['Recorded six-digest equality checked for every full row; bulk historical stream/artifact rehashing deferred to avoid interfering with root cache accounting. Only the failed request stdout/stderr/exit were rehashed.',
                         'Failed request completed compiler exit0 and retained streams, but pressure rejection preceded six-digest computation and metric/JIT row persistence; no full-parity or timing reconstruction is claimed.',
                         'The 8GiB parent counters cover the shared ancestor, including resident paused JVMs and filesystem/kernel state; they do not identify an exclusive causal process or establish a candidate slowdown.',
                         'Before/after memory samples do not identify transient peak or reclaimed objects. memory.events.local max counts limit encounters, not bytes; zero new OOM does not satisfy the zero-pressure gate.',
                         'No accepted five-cohort inference, restart, weakening, deletion, smoothing or row replacement. Separate partial descriptive observations are not an accepted benchmark summary. Root README/config/prespec and original results are unchanged.']}
(OUT / 'failure-audit.json').write_text(json.dumps(report, indent=2, allow_nan=False) + '\n')
with (OUT / 'retained-row-index.csv').open('w', newline='') as stream:
    writer = csv.DictWriter(stream, fieldnames=list(indices[0])); writer.writeheader(); writer.writerows(indices)
print(json.dumps({k: report[k] for k in ('status', 'completed_cohorts', 'partial_cohorts', 'full_metric_digest_rows', 'warm_rows', 'measured_metric_rows', 'complete_pairs', 'telemetry_anomalies', 'cleanup', 'performance_statistics')}, indent=2))
