"""Descriptive failed-run inventory after the independent no-restart decision.

Reads retained JSON only. No benchmark execution, imports, pooled inference,
confidence interval, accepted summary, exclusions, or replacement of raw rows.
"""
import csv
import json
import math
from pathlib import Path
import statistics as st

OUT = Path(__file__).resolve().parent
ROOT = OUT.parents[1]
RUN = ROOT / 'results' / OUT.name
read = lambda p: json.loads(Path(p).read_text())
audit = read(OUT / 'failure-audit.json')
assert audit['status'] == 'AUDIT_COMPLETED_FAILED_RUN_REMAINS_REJECTED'
names = ('STOCK', 'accumulator')
all_rows, measured_rows, cohorts = [], [], []

def jit(rows):
    valid = [r['jit_before_ms'] >= 0 and r['jit_after_ms'] >= r['jit_before_ms'] for r in rows]
    deltas = [r['jit_after_ms'] - r['jit_before_ms'] if ok else None for r, ok in zip(rows, valid)]
    total = sum(deltas) if all(valid) else None
    wall = sum(r['wall_s'] for r in rows)
    fraction = total / (1000 * wall) if total is not None else None
    return {'requests': [r['request'] for r in rows], 'deltas_ms': deltas, 'sum_ms': total,
            'wall_sum_s': wall, 'fraction': fraction, 'above_1_percent': fraction > .01 if fraction is not None else None,
            'above_2_percent': fraction > .02 if fraction is not None else None, 'invalid_intervals': valid.count(False)}

def stats(rows):
    result = {}
    for label, key in [('wall_s', 'wall_s'), ('CPU_s', 'cpu_s'), ('allocated_bytes', 'allocated_bytes'),
                       ('current_RSS_MiB', 'rss_mib')]:
        values = [r[key] for r in rows]
        result[label] = {'median': st.median(values), 'mean': st.mean(values), 'min': min(values), 'max': max(values),
                         'sample_SD': st.stdev(values) if len(values) > 1 else None}
    return result

for n in (1, 2):
    path = RUN / f'cohort-{n}'
    pairs, warm = read(path / 'pairs.json'), read(path / 'warmup.json')
    variants = {}
    for name in names:
        rows = [read(p) for p in sorted((path / name / 'protocol').glob('*/row.json'), key=lambda p: int(p.parent.name))]
        for r in rows:
            for a, b in [(r['wall_s'], r['wall_ns'] / 1e9), (r['cpu_s'], r['cpu_ns'] / 1e9),
                         (r['rss_mib'], r['rss_kib'] / 1024), (r['peak_rss_mib'], r['hwm_kib'] / 1024)]:
                assert math.isclose(a, b, rel_tol=1e-11, abs_tol=1e-11)
            assert r['wall_s'] > 0 and r['cpu_s'] >= 0 and r['rss_mib'] > 0 and r['peak_rss_mib'] > 0
            assert 0 <= r['allocation_before'] <= r['allocation_after']
            assert r['allocated_bytes'] == r['allocation_after'] - r['allocation_before']
            item = {'cohort': n, 'variant': name, 'kind': r['kind'], 'request': r['request'],
                    'pair': r['request'] - 60 if r['kind'] == 'measure' else '',
                    'retention': 'COMPLETE_UNPAIRED_REQUEST' if (n, name, r['request']) == (2, 'STOCK', 68) else 'COMPLETE_ROW',
                    **{k: r[k] for k in ['wall_ns', 'wall_s', 'cpu_ns', 'cpu_s', 'allocation_before', 'allocation_after',
                       'allocated_bytes', 'rss_kib', 'rss_mib', 'hwm_kib', 'peak_rss_mib', 'cgroup_current_bytes',
                       'cgroup_peak_bytes', 'gc_count', 'gc_ms', 'jit_before_ms', 'jit_after_ms', 'exit']},
                    'JIT_delta_ms': r['jit_after_ms'] - r['jit_before_ms'],
                    'JIT_counter_valid': r['jit_before_ms'] >= 0 and r['jit_after_ms'] >= r['jit_before_ms'],
                    'strict_recorded_digests': 6, 'raw_row': str(path / name / 'protocol' / str(r['request']) / 'row.json')}
            all_rows.append(item)
            if r['kind'] == 'measure':
                measured_rows.append(item)
        warms = warm['rows'][name]
        a, b = warms[-16:-8], warms[-8:]
        drift = st.median(r['wall_s'] for r in b) / st.median(r['wall_s'] for r in a) - 1
        warm_jit = jit(b)['fraction']
        spread = max(r['wall_s'] for r in b) / min(r['wall_s'] for r in b)
        assert len(warms) == 60 and abs(drift) <= .03 and warm_jit <= .01 and spread <= 1.15
        paired = [pair['rows'][name] for pair in pairs]
        measured = [r for r in rows if r['kind'] == 'measure']
        variants[name] = {'complete_pair_count': len(paired), 'available_measured_rows': len(measured),
                          'complete_pairs_descriptive': stats(paired),
                          'maximum_retained_lifetime_HWM_MiB': max(r['peak_rss_mib'] for r in rows),
                          'warm_gate': {'requests': 60, 'median_drift': drift, 'JIT_fraction': warm_jit, 'spread': spread},
                          'all_available_measured_JIT': jit(measured),
                          'complete_pairs_JIT': jit(paired),
                          'fixed5_JIT_blocks': [dict(jit(measured[i:i+5]), partial=len(measured[i:i+5]) != 5) for i in range(0, len(measured), 5)]}
        if n == 1:
            variants[name]['last10_vs_first10_median_wall_drift'] = st.median(r['wall_s'] for r in measured[10:]) / st.median(r['wall_s'] for r in measured[:10]) - 1
        if len(measured) > len(paired):
            variants[name]['unpaired_metric_row'] = measured[-1]
    ratios = {}
    for metric in ('wall_s', 'CPU_s', 'allocated_bytes', 'current_RSS_MiB'):
        for statistic in ('median', 'mean'):
            ratios[metric + '_' + statistic] = variants['accumulator']['complete_pairs_descriptive'][metric][statistic] / variants['STOCK']['complete_pairs_descriptive'][metric][statistic]
    ratios['maximum_retained_lifetime_HWM_MiB'] = variants['accumulator']['maximum_retained_lifetime_HWM_MiB'] / variants['STOCK']['maximum_retained_lifetime_HWM_MiB']
    cohorts.append({'cohort': n, 'status': 'COMPLETE_COHORT_OF_FAILED_RUN' if n == 1 else 'PARTIAL_COHORT_OF_FAILED_RUN',
                    'complete_pairs': len(pairs), 'variants': variants, 'descriptive_candidate_over_STOCK': ratios})
assert len(all_rows) == 295 and len(measured_rows) == 55
missing = {key: '' for key in all_rows[0]}
missing.update(cohort=2, variant='accumulator', kind='measure', request=68, pair=8,
               retention='FAILED_PRESSURE_NO_METRIC_ROW', strict_recorded_digests='NOT_COMPUTED',
               raw_row='NOT_SAVED')
for filename, rows in [('all-available-rows.csv', all_rows + [missing]), ('measured-rows-with-JIT.csv', measured_rows + [missing])]:
    with (OUT / filename).open('w', newline='') as stream:
        writer = csv.DictWriter(stream, fieldnames=list(all_rows[0])); writer.writeheader(); writer.writerows(rows)
report = {'status': 'INCOMPLETE_REJECTED_DESCRIPTIVE_ONLY', 'restart_decision': 'Root closed NO_RESTART before receiving partial performance statistics',
          'accepted_performance_claim': False, 'confidence_intervals': None, 'pooled_estimate': None,
          'complete_pairs': 27, 'available_measured_rows': 55, 'missing_measured_rows': [missing], 'per_cohort': cohorts,
          'limitations': ['No five-cohort confirmation, confidence interval, pooling with screens, or speed promotion.',
                         'Cohort2 comparisons use only its seven completed pairs, with STOCK request68 separately retained; no missing candidate values are imputed.',
                         'Lifetime HWM is maximum retained process high-water mark including warmup/JIT, not incremental request RSS. Cohort2 coverage is unequal and incomplete.',
                         'Current RSS is a postrequest snapshot; allocation is VM-wide allocated-byte delta, not retained heap.',
                         'CompilationMXBean is approximate accumulated compiler elapsed time, may sum threads and include pause-spanning counter effects; JIT/wall is not CPU%.',
                         'Measured >1% and >2% annotations are diagnostics, never selective row exclusions. Partial fixed5 blocks are explicitly marked.',
                         'The frozen reviewer rejected failure evidence as required. Failure/resource/provenance checks and their scoped hashing limitation are in failure-audit.json.']}
(OUT / 'partial-descriptive.json').write_text(json.dumps(report, indent=2, allow_nan=False) + '\n')
lines = ['# Failed v4 confirmation: retained partial evidence', '',
         'The fixed five-cohort confirmation failed in cohort 2 and remains rejected. No restart was selected; that decision closed before these partial performance summaries were inspected. These numbers describe retained observations only: no accepted benchmark summary, five-cohort interval, pooled estimate, or speed promotion.', '',
         'The independent failure audit checked the runner/manifest/grant/spec/driver identities, loaded class origin and target-class bytes, compiler/resource configuration, zero-pressure status of all saved measured rows, metric conversions, and all 295 recorded six-digest sets against the pinned oracle. The frozen reviewer rejected the run with `failure evidence present: no accepted result`; its rejection is retained.', '',
         'Cohort 1 contains 60 warmups per JVM and 20 complete measured pairs. Cohort 2 contains 60 warmups per JVM, seven complete measured pairs, one saved unpaired STOCK request 68, and failed accumulator request 68. Cohorts 3–5 never started. There are 296 protocol request directories, 295 complete metric/digest rows (240 warmup, 55 measured), and 27 complete measured pairs. No rows were dropped, replaced, smoothed, or imputed.', '',
         '## Failure and cleanup', '',
         'Accumulator request 68 was the second variant request within cohort 2, measured pair 8. The shared 8 GiB parent gained 387 local memory-max events and 14,879 direct-scan/direct-steal counts. The compiler returned exit 0, but the zero-pressure assertion ran before six-digest computation and row persistence. Its stdout/stderr/exit streams remain and independently hash to the oracle; timing, allocation, RSS and JIT metrics were not saved. The final output artifacts exist, but full six-digest parity for this request was not established. STOCK request 68 remains separately available below and in the CSV.', '',
         'The snapshots show parent memory.current 7,709,859,840 → 8,051,613,696 bytes, below its 8,589,934,592-byte cap at those two instants. They do not rule out transient limit encounters. The shared parent covers paused JVMs and filesystem/kernel/cache state; its counters do not identify an exclusive causal process or prove an accumulator leak or slowdown. PSI was unavailable, not zero. All request and outer OOM deltas were zero.', '',
         'All four owned JVM drivers exited 0; all four containers were removed without fallback. Cleanup errors, remaining owned containers/processes and active compiler lists were empty. No telemetry exceptions were needed. No JVM/build execution or large-JAR/bulk-artifact rehashing was performed by this audit. All stored six-digest sets were checked; only the failing request’s three streams were rehashed, so this is a scoped failed-run audit rather than a substitute PASS from the frozen reviewer.', '',
         '## Descriptive comparisons', '',
         'Each table uses the complete pairs within that cohort. The extra STOCK row is shown separately and never silently discarded. Lifetime HWM uses all retained rows, including warmup; cohort 2 has unequal and incomplete endpoint coverage.', '']
for cohort in cohorts:
    n, v = cohort['cohort'], cohort['variants']
    lines.extend([f"Cohort {n}: {cohort['complete_pairs']} complete pairs; {cohort['status']}.", '',
                  '| Metric | STOCK | accumulator | Candidate / STOCK |', '| --- | ---: | ---: | ---: |'])
    for label, key, statistic in [('Elapsed median s', 'wall_s', 'median'), ('Elapsed mean s', 'wall_s', 'mean'),
                                   ('CPU median s', 'CPU_s', 'median'), ('CPU mean s', 'CPU_s', 'mean'),
                                   ('Allocated median bytes', 'allocated_bytes', 'median'), ('Allocated mean bytes', 'allocated_bytes', 'mean'),
                                   ('Current RSS median MiB', 'current_RSS_MiB', 'median')]:
        b, c = [v[name]['complete_pairs_descriptive'][key][statistic] for name in names]
        lines.append(f'| {label} | {b:.9f} | {c:.9f} | {c/b:.9f} |')
    b, c = [v[name]['maximum_retained_lifetime_HWM_MiB'] for name in names]
    lines.extend([f'| Maximum retained lifetime HWM MiB | {b:.9f} | {c:.9f} | {c/b:.9f} |', ''])
    for name in names:
        j = v[name]['all_available_measured_JIT']; w = v[name]['warm_gate']
        lines.extend([f"{name}: all {v[name]['available_measured_rows']} available measured JIT deltas {j['deltas_ms']} ms; total {j['sum_ms']} ms / {j['wall_sum_s']:.9f} s = {j['fraction']:.6%}; >1%: {j['above_1_percent']}; >2%: {j['above_2_percent']}; invalid intervals: {j['invalid_intervals']}.",
                      f"Warm gate independently recomputed at request 60: drift {w['median_drift']:.6%}, JIT/wall {w['JIT_fraction']:.6%}, spread {w['spread']:.9f}.", ''])
        blocks = ', '.join(f"{block['requests'][0]}–{block['requests'][-1]}: {block['fraction']:.6%}" + (' (partial)' if block['partial'] else '') for block in v[name]['fixed5_JIT_blocks'])
        lines.extend([f'Fixed five-request JIT blocks: {blocks}.', ''])
    if n == 1:
        lines.extend([f"Cohort 1 last10/first10 median elapsed drift: STOCK {v['STOCK']['last10_vs_first10_median_wall_drift']:.6%}; accumulator {v['accumulator']['last10_vs_first10_median_wall_drift']:.6%}. Diagnostic only; no rows excluded.", ''])
    else:
        r = v['STOCK']['unpaired_metric_row']
        lines.extend([f"Unpaired STOCK request 68: elapsed {r['wall_s']:.9f} s; CPU {r['cpu_s']:.9f} s; allocated {r['allocated_bytes']} bytes; current RSS {r['rss_mib']:.9f} MiB; lifetime HWM {r['peak_rss_mib']:.9f} MiB; JIT delta {r['jit_after_ms']-r['jit_before_ms']} ms. Accumulator request 68 has no saved metric row.", ''])
lines.extend(['CompilationMXBean reports approximate accumulated compilation elapsed time, may sum compiler threads, and may reflect compiles spanning container pauses. The JIT/wall annotations are not CPU percentages or proof of a performance cause. Allocation is a VM-wide request delta; RSS/HWM measures have the different scopes described above.', '',
              'The screen findings remain exploratory, and this failed direct-STOCK attempt cannot establish a speed gain. Source correctness and separate upstream test results do not turn an incomplete performance run into confirmation.', '',
              'Evidence: `failure-audit.json`, `retained-row-index.csv`, `partial-descriptive.json`, `all-available-rows.csv` (all 295 saved rows plus the explicit missing row), `measured-rows-with-JIT.csv`, and `frozen-reviewer-rejection/review.json`. Original results, runner, prespecification, root README and publication config were not changed.'])
(OUT / 'failure-audit.md').write_text('\n'.join(lines) + '\n')
print('\n'.join(lines))
