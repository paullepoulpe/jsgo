"""Summarize this completed, independently reviewed one-cohort screen; no runtime imports."""
import csv
import json
import hashlib
from pathlib import Path

out = Path(__file__).resolve().parent
read = lambda path: json.loads(path.read_text())
review = read(out / 'review.json')
assert review['status'] == 'PASS independent Java-factor review'
assert review['cohorts'] == 1 and review['pairs_per_cohort'] == 5
root = Path(review['results'])
pairs = read(root / 'cohort-1/pairs.json')
cohort = review['per_cohort'][0]
names = [review['baseline'], review['candidate']]
summary = {key: review[key] for key in ('status', 'results', 'baseline', 'candidate', 'cohorts', 'pairs_per_cohort')}
summary.update(review=str(out / 'review.json'), screen_only=True, confidence_intervals=None,
               settled_speed_gain_established=False, promotion_status='NONE_PROMOTED', variants={})
jit_rows = []
for name in names:
    v = cohort['variants'][name]
    rows = [pair['rows'][name] for pair in pairs]
    assert [row['wall_s'] for row in rows] == v['wall_seconds']['raw']
    valid = [row['jit_before_ms'] >= 0 and row['jit_after_ms'] >= row['jit_before_ms'] for row in rows]
    deltas = [row['jit_after_ms'] - row['jit_before_ms'] if ok else None for row, ok in zip(rows, valid)]
    total = sum(deltas) if all(valid) else None
    wall = sum(row['wall_s'] for row in rows)
    fraction = total / (1000 * wall) if total is not None else None
    flag = ('MEASURED_JIT_UNAVAILABLE' if fraction is None else
            'MEASURED_JIT_ABOVE_2_PERCENT' if fraction > .02 else 'MEASURED_JIT_AT_OR_BELOW_2_PERCENT')
    summary['variants'][name] = {
        'warmup_requests': v['warm_requests'], 'measured_requests': len(rows),
        'wall_median_s': v['wall_seconds']['median'], 'wall_mean_s': v['wall_seconds']['mean'],
        'wall_sample_SD_s': v['wall_seconds']['sample_SD'], 'wall_CV': v['wall_seconds']['CV'],
        'CPU_median_s': v['CPU_seconds']['median'], 'CPU_mean_s': v['CPU_seconds']['mean'],
        'allocated_median_bytes': v['allocated_bytes']['median'], 'allocated_mean_bytes': v['allocated_bytes']['mean'],
        'current_RSS_median_MiB': v['current_RSS_MiB']['median'], 'current_RSS_mean_MiB': v['current_RSS_MiB']['mean'],
        'lifetime_HWM_MiB': v['final_lifetime_HWM_MiB'], 'readiness': v['settling'],
        'measured_request_numbers': [row['request'] for row in rows], 'measured_wall_s': [row['wall_s'] for row in rows],
        'measured_JIT_deltas_ms': deltas, 'measured_JIT_sum_ms': total, 'measured_wall_sum_s': wall,
        'measured_batch_JIT_fraction': fraction, 'measured_JIT_flag': flag,
        'measured_JIT_invalid_interval_count': valid.count(False),
        'measured_GC_counts': [row['gc_count'] for row in rows], 'measured_GC_ms': [row['gc_ms'] for row in rows]}
    for pair, row, delta, ok in zip(pairs, rows, deltas, valid):
        jit_rows.append({'cohort': 1, 'variant': name, 'pair': pair['pair'], 'request': row['request'],
                         'wall_s': row['wall_s'], 'JIT_before_ms': row['jit_before_ms'],
                         'JIT_after_ms': row['jit_after_ms'], 'JIT_valid': ok, 'JIT_delta_ms': delta})
b, c = [summary['variants'][name] for name in names]
summary['candidate_over_baseline'] = cohort['candidate_over_baseline']
summary['baseline_over_candidate'] = {k: 1 / v for k, v in cohort['candidate_over_baseline'].items()}
summary['descriptive_candidate_over_baseline'] = {
    'current_RSS_median': c['current_RSS_median_MiB'] / b['current_RSS_median_MiB'],
    'wall_mean': c['wall_mean_s'] / b['wall_mean_s'], 'CPU_mean': c['CPU_mean_s'] / b['CPU_mean_s']}
summary['same_pair_ratios_descriptive_only'] = cohort['same_pair_ratios']
summary.update(all_six_digests_every_request=review['strict_digests_every_request'] == 6,
               requests_with_all_six_digests=review['warmup_requests'] + review['measured_requests'],
               all_graceful_driver_exit_codes_zero=review['all_graceful_driver_exit_codes_zero'],
               telemetry_anomalies=review['telemetry_anomalies'], PSI_telemetry=review['PSI_telemetry'],
               manifest_sha256=review['manifest_sha256'], runner_sha256=review['runner_sha256'],
               stopped_proof_sha256=review['stopped_proof_sha256'],
               limitations=review['limits'], Java_vs_Go_claim=False)
summary['finding'] = ('One-cohort CONTROL screen only. Elapsed and CPU medians increased while allocated bytes fell 1.4506%. '
                      'CONTROL recorded an 8554 ms cumulative compilation-time counter delta on its fifth measured request, '
                      'putting its batch JIT fraction at 42.7717%, versus 1.8556% for the candidate. '
                      'CompilationMXBean measures approximate accumulated compilation elapsed time, potentially summed across threads, not CPU%. '
                      'A compilation spanning a paused container may contribute paused duration when a later counter update occurs; '
                      'that mechanism was not diagnosed for this jump. Treat the large ratio as a potentially contaminated elapsed-time proxy, '
                      'not proof of JIT CPU share or slowdown causality. This single cohort does not establish a settled benefit; '
                      'lower allocation and approximate RSS remain descriptive. '
                      'No independent-cohort interval or speed gain is established; none promoted. '
                      'All rows and measured JIT diagnostics are retained, with no postmeasurement rejection.')
experiment = out.parents[1]
manifest = read(experiment / 'build/source-position-last/manifest.json')
sha = lambda path: hashlib.sha256(path.read_bytes()).hexdigest()
source_identity = {}
for name in names:
    variant = manifest['variants'][name]
    receipt_path = Path(variant['build_proof'])
    receipt = read(receipt_path)
    assert sha(receipt_path) == variant['build_proof_sha256'] == manifest['compiler_build_receipt_sha256'][name]
    assert receipt['status'] == 'PASS' and receipt['exit'] == 0 and receipt['variant'] == name
    assert receipt['javac_flags'] == manifest['javac_flags'] == variant['javac_flags']
    assert receipt['classes'] == variant['classes']
    checks = {}
    for relative, original in receipt['original_sources'].items():
        control = experiment / 'candidates/source-position-last/control-src' / relative
        effective = experiment / 'candidates/source-position-last' / ('control-src' if name == 'CONTROL' else 'src') / relative
        relocated = receipt_path.parent / 'src' / relative
        assert sha(control) == original
        assert sha(effective) == receipt['effective_sources'][relative]
        assert sha(relocated) == receipt['relocated_sources'][relative]
        checks[relative] = {'original_sha256': original, 'effective_sha256': receipt['effective_sources'][relative],
                            'relocated_sha256': receipt['relocated_sources'][relative]}
    source_identity[name] = {'status': 'PASS small source/receipt identity checks', 'build_receipt': str(receipt_path),
                             'build_receipt_sha256': variant['build_proof_sha256'], 'sources': checks, 'classes': receipt['classes']}
summary['source_identity_checks'] = source_identity
(out / 'compact-summary.json').write_text(json.dumps(summary, indent=2, allow_nan=False) + '\n')
with (out / 'measured-jit.csv').open('w', newline='') as stream:
    writer = csv.DictWriter(stream, fieldnames=list(jit_rows[0])); writer.writeheader(); writer.writerows(jit_rows)
lines = ['PASS independent review; one cohort, five measured pairs; NONE_PROMOTED. No confidence interval or established speed gain.',
         f"All {summary['requests_with_all_six_digests']} requests matched all six oracle digests; all driver exits zero; no new OOM, resource/provenance/pressure-proxy or cleanup violation; PSI unavailable.",
         'No telemetry anomalies. Approximate lifetime HWM includes warmup/JIT; current RSS is a postrequest snapshot; allocation is VM-wide request delta, not retained heap.', '',
         f'| Metric | {names[0]} | {names[1]} | Candidate / CONTROL |', '| --- | ---: | ---: | ---: |']
for label, key in [('Elapsed median s', 'wall_median_s'), ('Elapsed mean s', 'wall_mean_s'),
                   ('CPU median s', 'CPU_median_s'), ('CPU mean s', 'CPU_mean_s'),
                   ('Allocated median bytes', 'allocated_median_bytes'), ('Current RSS median MiB', 'current_RSS_median_MiB'),
                   ('Approx. lifetime HWM MiB', 'lifetime_HWM_MiB')]:
    lines.append(f'| {label} | {b[key]:.9f} | {c[key]:.9f} | {c[key] / b[key]:.9f} |')
lines.append('')
for name in names:
    v = summary['variants'][name]; ready = v['readiness']
    lines.extend([f"{name}: {v['warmup_requests']} matched warmups; drift {ready['last4_vs_prior4_median_drift']:.6%}, warm JIT/wall {ready['last4_JIT_time_over_wall']:.6%}, spread {ready['last4_max_over_min']:.9f}.",
                  f"Measured JIT deltas {v['measured_JIT_deltas_ms']} ms; sum {v['measured_JIT_sum_ms']} ms / {v['measured_wall_sum_s']:.9f} s = {v['measured_batch_JIT_fraction']:.6%}; {v['measured_JIT_flag']}; invalid intervals {v['measured_JIT_invalid_interval_count']}.", ''])
lines.extend([summary['finding'], '', 'The independent reviewer checked manifest/runner pins, matched source-family/compiler flags, loaded class locations/digests, the exact compiler command and all six recorded outputs on every request. Owned JVM exits and OOM deltas were zero, and both containers were removed. No full-JAR hashing or JVM/build execution occurred in this review. Source/probe correctness review is separate.', '',
              'Evidence: review.json, rows.csv, cohorts.csv, ratios.csv, measured-jit.csv and compact-summary.json. Original results are unchanged.'])
lines.extend(['', '| Source identity (SHA-256) | CONTROL | source-position-last |', '| --- | --- | --- |'])
relative = 'com/google/javascript/jscomp/parsing/parser/LineNumberScanner.java'
for key in ('original_sha256', 'effective_sha256', 'relocated_sha256'):
    lines.append('| ' + key + ' | ' + source_identity[names[0]]['sources'][relative][key] + ' | ' + source_identity[names[1]]['sources'][relative][key] + ' |')
lines.extend(['', 'Small source/receipt identity checks passed against the matching manifest. The rebuilt family is LineNumberScanner only, with matching CONTROL/candidate compiler flags and independent loaded-origin checks. The change keeps one scanner-local immutable SourcePosition for an identical file/offset/line/column result; no cross-request compiler cache is introduced. Semantic/probe review is separate from this measurement review.'])
(out / 'compact-summary.md').write_text('\n'.join(lines) + '\n')
print('\n'.join(lines))
