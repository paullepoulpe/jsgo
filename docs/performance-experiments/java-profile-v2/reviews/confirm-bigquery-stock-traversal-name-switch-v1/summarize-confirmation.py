#!/usr/bin/env python3
"""Small post-run report transform; no runtime or benchmark imports."""
from pathlib import Path
import csv
import json
import math
import statistics

ROOT = Path(__file__).resolve().parent
review = json.loads((ROOT / 'review.json').read_text())
with (ROOT / 'rows.csv').open(newline='') as stream:
    rows = list(csv.DictReader(stream))
assert review['status'] == 'PASS independent Java-factor review'
assert review['cohorts'] == 3 and review['pairs_per_cohort'] == 5
baseline, candidate = review['baseline'], review['candidate']


def interval(values):
    logs = list(map(math.log, values))
    center = statistics.mean(logs)
    half = 4.302652729911275 * statistics.stdev(logs) / math.sqrt(3)
    return {'cohort_values': values, 'geomean': math.exp(center),
            'ci95_low': math.exp(center - half), 'ci95_high': math.exp(center + half),
            'df': 2, 'critical_t': 4.302652729911275}


cohorts, jit_rows = [], []
for source in review['per_cohort']:
    c = {'cohort': source['cohort'], 'variants': {},
         'candidate_over_baseline': source['candidate_over_baseline']}
    for name, v in source['variants'].items():
        measured = [r for r in rows if int(r['cohort']) == c['cohort']
                    and r['variant'] == name and r['kind'] == 'measure']
        assert len(measured) == 5
        deltas = [int(r['JIT_delta_ms']) for r in measured]
        walls = [float(r['wall_s']) for r in measured]
        # Preserve unavailable/decreasing observations explicitly if ever present.
        valid = all(d >= 0 for d in deltas)
        fraction = sum(deltas) / (1000 * sum(walls)) if valid else None
        flag = ('JIT_COUNTER_UNAVAILABLE_OR_DECREASING' if not valid else
                'MEASURED_JIT_ABOVE_2_PERCENT' if fraction > .02 else None)
        c['variants'][name] = {
            'warmup_requests': v['warm_requests'], 'measured_requests': 5,
            'wall_median_s': v['wall_seconds']['median'],
            'wall_mean_s': v['wall_seconds']['mean'],
            'wall_sample_SD_s': v['wall_seconds']['sample_SD'],
            'CPU_median_s': v['CPU_seconds']['median'],
            'allocated_median_bytes': v['allocated_bytes']['median'],
            'current_RSS_median_MiB': v['current_RSS_MiB']['median'],
            'lifetime_HWM_MiB': v['final_lifetime_HWM_MiB'],
            'measured_request_numbers': [int(r['request']) for r in measured],
            'measured_wall_s': walls, 'measured_JIT_deltas_ms': deltas,
            'measured_JIT_sum_ms': sum(deltas), 'measured_wall_sum_s': sum(walls),
            'measured_batch_JIT_fraction': fraction, 'measured_JIT_flag': flag}
        jit_rows.append({'cohort': c['cohort'], 'variant': name,
                         'measured_JIT_deltas_ms': ' '.join(map(str, deltas)),
                         'JIT_sum_ms': sum(deltas), 'wall_sum_s': sum(walls),
                         'JIT_fraction': fraction, 'flag': flag or ''})
    b, v = c['variants'][baseline], c['variants'][candidate]
    c['descriptive_candidate_over_baseline'] = {
        'mean_wall': v['wall_mean_s'] / b['wall_mean_s'],
        'current_RSS_median': v['current_RSS_median_MiB'] / b['current_RSS_median_MiB']}
    cohorts.append(c)

for metric, result in review['ratios'].items():
    calculated = interval([c['candidate_over_baseline'][metric] for c in cohorts])
    for key in ('geomean', 'ci95_low', 'ci95_high'):
        assert math.isclose(calculated[key], result['log_t95_df2_candidate_over_baseline'][key], rel_tol=1e-12)
sensitivities = {metric: interval([c['descriptive_candidate_over_baseline'][metric] for c in cohorts])
                 for metric in ('mean_wall', 'current_RSS_median')}
summary = {
    'status': review['status'], 'results': review['results'], 'review': str(ROOT / 'review.json'),
    'baseline': baseline, 'candidate': candidate, 'cohorts': 3, 'screen_only': False,
    'screen_observations_pooled': 0, 'settled_speed_gain_established': False,
    'per_cohort': cohorts, 'ratios': review['ratios'], 'descriptive_sensitivities': sensitivities,
    'all_six_digests_every_request': True,
    'all_graceful_driver_exit_codes_zero': review['all_graceful_driver_exit_codes_zero'],
    'telemetry_anomalies': review['telemetry_anomalies'],
    'JIT_flagged_batches': sum(v['measured_JIT_flag'] is not None for c in cohorts for v in c['variants'].values()),
    'JIT_policy': '../../MEASURED-JIT-REVIEW-NOTE.md',
    'interpretation': 'No demonstrated speed gain: median-wall log-t95(df2) includes 1, two of three cohort ratios exceed 1, and five measured batches exceed the predeclared JIT annotation. Approximate lifetime HWM is higher in all three cohorts. Preserve the original PASS and every outcome; no request or cohort discarded or retried.',
    'limits': review['limits']}
(ROOT / 'compact-summary.json').write_text(json.dumps(summary, indent=2, allow_nan=False) + '\n')
with (ROOT / 'measured-jit.csv').open('w', newline='') as stream:
    writer = csv.DictWriter(stream, fieldnames=list(jit_rows[0]))
    writer.writeheader()
    writer.writerows(jit_rows)

lines = [
    'Independent review passed for the three-cohort BigQuery traversal-name-switch / published STOCK confirmation. **No demonstrated speed gain:** elapsed cohort ratios were 1.025077180, 0.938362416 and 1.004951924; their geometric mean was 0.988759840 with 95% log-cohort t(df=2) interval [0.881161793, 1.109496609]. The interval includes useful gains and regressions. Two of three cohorts had higher candidate elapsed medians. Five of six measured batches exceeded the predeclared JIT annotation; a settled gain is not established.',
    '',
    '| Cohort | Variant | Warmups | Wall median s | Wall mean s | CPU median s | Allocated median bytes | Approx. lifetime HWM MiB | Current RSS median MiB |',
    '| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |']
for c in cohorts:
    for name, v in c['variants'].items():
        lines.append(f"| {c['cohort']} | {name} | {v['warmup_requests']} | {v['wall_median_s']:.9f} | {v['wall_mean_s']:.9f} | {v['CPU_median_s']:.6f} | {v['allocated_median_bytes']:,} | {v['lifetime_HWM_MiB']:.6f} | {v['current_RSS_median_MiB']:.6f} |")
lines += ['', '| Cohort | Elapsed ratio | CPU ratio | Allocated ratio | Approx. lifetime HWM ratio | Current RSS median ratio | Mean-wall ratio |',
          '| --- | ---: | ---: | ---: | ---: | ---: | ---: |']
for c in cohorts:
    r, s = c['candidate_over_baseline'], c['descriptive_candidate_over_baseline']
    lines.append(f"| {c['cohort']} | {r['elapsed']:.9f} | {r['CPU']:.9f} | {r['allocated_bytes']:.9f} | {r['RSS_lifetimeHWM']:.9f} | {s['current_RSS_median']:.9f} | {s['mean_wall']:.9f} |")
lines += ['', '| Metric | Geometric mean candidate/STOCK | 95% log-t(df=2) interval | Discrete whole-cohort bootstrap 95% |',
          '| --- | ---: | --- | --- |']
for metric, r in review['ratios'].items():
    t, b = r['log_t95_df2_candidate_over_baseline'], r['bootstrap95_candidate_over_baseline']
    lines.append(f"| {metric} | {t['geomean']:.9f} | [{t['ci95_low']:.9f}, {t['ci95_high']:.9f}] | [{b['ci95_low']:.9f}, {b['ci95_high']:.9f}] |")
s = sensitivities['mean_wall']
lines += ['', f"Mean-wall sensitivity: geometric mean {s['geomean']:.9f}, log-t95(df=2) [{s['ci95_low']:.9f}, {s['ci95_high']:.9f}]. Approximate lifetime HWM was 4.97–6.24% higher in every cohort; geometric mean 5.47% higher. Its descriptive interval is also above 1, but this secondary telemetry is approximate and not an exact physical-memory effect. Allocation and CPU show no established improvement. Current RSS snapshots differ from the lifetime window.",
          '', '| Cohort | Variant | All measured JIT deltas ms | Sum JIT ms | Sum elapsed s | JIT / elapsed | >2% flag |',
          '| --- | --- | --- | ---: | ---: | ---: | --- |']
for r in jit_rows:
    fraction = f"{100*r['JIT_fraction']:.6f}%" if r['JIT_fraction'] is not None else 'unavailable'
    lines.append(f"| {r['cohort']} | {r['variant']} | {r['measured_JIT_deltas_ms'].replace(' ', ', ')} | {r['JIT_sum_ms']} | {r['wall_sum_s']:.9f} | {fraction} | {'Yes' if r['flag'] else 'No'} |")
lines += ['',
    'Every warm and measured request passed all six oracle digests: 244 warm requests and 30 measured requests across six JVMs. Resource, pressure, loaded-origin and stopped/zero-new-OOM review passed. One existing-policy telemetry exception was admitted: STOCK cohort 1 VmHWM fell 168 KiB between warmup requests 6 and 7 (1,100,692 to 1,100,524 KiB), then rose to 1,288,384 KiB at request 8. The raw rows and later maximum are unchanged. No measured VmHWM dip was admitted, and no new rule was introduced.',
    '',
    'The three fresh independent pair-of-JVM cohorts are the inference units; the 30 measured requests are not 30 independent JVM comparisons. The log-t interval uses exp(mean(log ratios) ± 4.302652729911275 × sample-SD(log ratios)/sqrt(3)); it assumes approximately normal independent log effects, which three observations cannot establish. The bootstrap is discrete and confined to these three observed cohorts, so it cannot be the sole inference. Screens are not pooled. Secondary intervals are descriptive, without a multiplicity-adjusted memory claim.',
    '',
    'JIT flags preserve the original protocol PASS and all outcomes; no row or cohort was discarded, replaced, or retried. The JIT counter sums compiler work and is not elapsed compile time. Allocated bytes are VM-wide ThreadMXBean request deltas, not retained heap. RSS telemetry is approximate; lifetime HWM includes warmup/JIT, current RSS is a postrequest snapshot, and missing PSI remains unavailable rather than zero.',
    '',
    'Retained evidence: [independent review](review.json), [all request rows](rows.csv), [cohort medians](cohorts.csv), [ratios and intervals](ratios.csv), [measured JIT](measured-jit.csv), [compact exact summary](compact-summary.json), and [report transform](summarize-confirmation.py). Bundle generalization remains a separate workload with no result in this review.']
(ROOT / 'compact-summary.md').write_text('\n'.join(lines) + '\n')
print(json.dumps({'status': summary['status'], 'JIT_flagged_batches': summary['JIT_flagged_batches'],
                  'elapsed_log_t95': review['ratios']['elapsed']['log_t95_df2_candidate_over_baseline'],
                  'mean_wall_sensitivity': sensitivities['mean_wall'],
                  'outputs': ['compact-summary.json', 'compact-summary.md', 'measured-jit.csv']}))
