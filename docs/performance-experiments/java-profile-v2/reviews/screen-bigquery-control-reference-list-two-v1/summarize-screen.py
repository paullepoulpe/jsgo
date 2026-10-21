"""Summarize this completed, independently reviewed one-cohort screen; no runtime imports."""
import csv
import json
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
summary['finding'] = ('One-cohort CONTROL screen only. Elapsed, CPU and allocated-byte medians and approximate lifetime HWM increased. '
                      'CompilationMXBean reports approximate accumulated compilation elapsed time, potentially summed across threads, not CPU%. '
                      'Container pauses and later counter updates may affect that proxy; no causal attribution is inferred. '
                      'No independent-cohort interval or speed gain is established; none promoted. '
                      'All rows and measured JIT diagnostics are retained, with no postmeasurement rejection.')
(out / 'compact-summary.json').write_text(json.dumps(summary, indent=2, allow_nan=False) + '\n')
with (out / 'measured-jit.csv').open('w', newline='') as stream:
    writer = csv.DictWriter(stream, fieldnames=list(jit_rows[0])); writer.writeheader(); writer.writerows(jit_rows)
lines = ['PASS independent review; one cohort, five measured pairs; NONE_PROMOTED. No confidence interval or established speed gain.',
         f"All {summary['requests_with_all_six_digests']} requests matched all six oracle digests; all driver exits zero; no new OOM, resource/provenance/pressure-proxy or cleanup violation; PSI unavailable.",
         f"Telemetry anomalies admitted under the existing warmup-only policy: {len(review['telemetry_anomalies'])}. Approximate lifetime HWM includes warmup/JIT; current RSS is a postrequest snapshot; allocation is VM-wide request delta, not retained heap.", '',
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
if review['telemetry_anomalies']:
    lines.extend(['', '| Existing-policy telemetry exception | Prior KiB | Later KiB | Delta KiB | Retained superseding sample |', '| --- | ---: | ---: | ---: | --- |'])
    for item in review['telemetry_anomalies']:
        later = item['later_strictly_larger_sample']
        lines.append(f"| {item['variant']}, warmup {item['prior_request']} to {item['request']} | {item['before']} | {item['after']} | {item['delta_kib']} | warmup {later['request']}: {later['hwm_kib']} KiB |")
    lines.extend(['', 'This is the pre-existing warmup-only VmHWM amendment, originally adopted after an earlier v1 observation. Exact retained rows prove the later strictly larger sample. No new exception, tolerance, smoothing or row/cohort removal was introduced; raw values and maximum-sample RSS ratios remain unchanged. No measured VmHWM or cgroup-peak decrease was admitted.'])
(out / 'compact-summary.md').write_text('\n'.join(lines) + '\n')
print('\n'.join(lines))
