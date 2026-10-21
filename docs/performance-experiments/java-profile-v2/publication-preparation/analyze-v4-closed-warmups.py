"""Tiny retrospective feasibility check of explicitly named closed v3 warmups.

No active-result discovery, benchmark import, JVM, Docker or executable hashing.
This does not relabel v3 outcomes or change the frozen v4 protocol.
"""
import ast
import json
import statistics
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
JOBS = (
    'screen-bigquery-control-accumulator-v1',
    'screen-bigquery-control-traversal-name-switch-v1',
    'confirm-bigquery-stock-traversal-name-switch-v1',
    'screen-bigquery-control-traversal-noop-v1',
    'screen-bigquery-control-traversal-leaf-name-v1',
    'screen-bigquery-control-equality-prototype-scalar-v1',
    'screen-bigquery-control-empty-parameters-v1',
    'screen-bigquery-control-source-position-last-v1',
    'screen-bigquery-control-reference-list-two-v1',
)
read = lambda path: json.loads(path.read_text())
# Execute only the isolated pure stability helper, never the runner module.
tree = ast.parse((ROOT / 'bench/run-v4.py').read_text())
function = next(node for node in tree.body if isinstance(node, ast.FunctionDef) and node.name == 'stability')
scope = {'statistics': statistics}
exec(compile(ast.Module(body=[function], type_ignores=[]), '<isolated-v4-stability>', 'exec'), scope)

def window(rows):
    old, new = rows[-16:-8], rows[-8:]
    assert len(old) == len(new) == 8
    valid = all(r['jit_before_ms'] >= 0 and r['jit_after_ms'] >= r['jit_before_ms'] for r in new)
    jit = sum(r['jit_after_ms'] - r['jit_before_ms'] for r in new) / (1000 * sum(r['wall_s'] for r in new)) if valid else None
    drift = statistics.median(r['wall_s'] for r in new) / statistics.median(r['wall_s'] for r in old) - 1
    spread = max(r['wall_s'] for r in new) / min(r['wall_s'] for r in new)
    stable = abs(drift) <= .03 and jit is not None and jit <= .01 and spread <= 1.15
    if len(rows) >= 60:
        exact = scope['stability'](rows)
        assert (exact['stable'], exact['last8_vs_prior8_median_drift'], exact['last8_JIT_time_over_wall'], exact['last8_max_over_min']) == (stable, drift, jit, spread)
    return {'JIT_fraction': jit, 'median_drift': drift, 'spread': spread, 'window_thresholds_pass_ignoring_minimum60': stable}

cohorts = []
for job in JOBS:
    root = ROOT / 'results' / job
    review = read(ROOT / 'reviews' / job / 'review.json')
    assert review['status'] == 'PASS independent Java-factor review' and review['results'] == str(root)
    assert read(root / 'summary.json')['complete'] is True
    for index in range(1, review['cohorts'] + 1):
        rows = read(root / f'cohort-{index}/warmup.json')['rows']
        assert len({len(v) for v in rows.values()}) == 1
        count = len(next(iter(rows.values())))
        windows = {n: {name: window(v[:n]) for name, v in rows.items()} for n in range(16, count + 1)}
        full = [n for n, variants in windows.items() if n >= 60 and all(v['window_thresholds_pass_ignoring_minimum60'] for v in variants.values())]
        below_min = [n for n, variants in windows.items() if all(v['window_thresholds_pass_ignoring_minimum60'] for v in variants.values())]
        cohorts.append({'job': job, 'cohort': index, 'observed_warmups_per_variant': count,
                        'final_window': windows[count],
                        'first_observed_paired_v4_pass': min(full) if full else None,
                        'first_paired_window_thresholds_pass_ignoring_minimum60': min(below_min) if below_min else None,
                        'minimum_observed_paired_max_JIT_fraction_after16': min(max(v['JIT_fraction'] for v in variants.values()) for variants in windows.values()),
                        'first_individual_JIT_at_most1pct': {name: next((n for n, variants in windows.items() if variants[name]['JIT_fraction'] <= .01), None) for name in rows},
                        'eligible_observed_prefix_count': max(0, count - 59)})
summary = {'scope': 'Nine explicitly named completed reviewed v3 BigQuery jobs; warmups only, no measured or active rows.',
           'protocol_unchanged': True, 'v4_runtime_started': False, 'cohort_pairs': len(cohorts),
           'JVM_series': 2 * len(cohorts), 'maximum_observed_warmups': max(c['observed_warmups_per_variant'] for c in cohorts),
           'pairs_reaching_minimum60': sum(c['observed_warmups_per_variant'] >= 60 for c in cohorts),
           'eligible_observed_paired_prefixes': sum(c['eligible_observed_prefix_count'] for c in cohorts),
           'pairs_with_an_observed_v4_gate_pass': sum(c['first_observed_paired_v4_pass'] is not None for c in cohorts),
           'pairs_ever_meeting_window_thresholds_ignoring_minimum60': sum(c['first_paired_window_thresholds_pass_ignoring_minimum60'] is not None for c in cohorts),
           'final_individual_windows_JIT_at_most1pct': sum(v['JIT_fraction'] <= .01 for c in cohorts for v in c['final_window'].values()),
           'cohorts': cohorts,
           'limits': ['Historical v3 stopping censors later warmups; no data beyond the observed maximum, so >90 need or pass probability cannot be estimated reliably.',
                      'CompilationMXBean is approximate accumulated compilation elapsed time, potentially contaminated by pause-spanning timers and delayed updates; not CPU%.',
                      'Descriptive feasibility only; no altered thresholds, v3 rejection, favorable subset or v4 runtime.']}
(Path(__file__).with_name('v4-closed-warmup-feasibility.json')).write_text(json.dumps(summary, indent=2, allow_nan=False) + '\n')
print(json.dumps({k: v for k, v in summary.items() if k != 'cohorts'}, indent=2))
for c in cohorts:
    print(c['job'], c['cohort'], c['observed_warmups_per_variant'],
          {name: round(v['JIT_fraction'] * 100, 6) for name, v in c['final_window'].items()},
          'first v4 paired pass', c['first_observed_paired_v4_pass'],
          'best paired JIT%', round(c['minimum_observed_paired_max_JIT_fraction_after16'] * 100, 6))
