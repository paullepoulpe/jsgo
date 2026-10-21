Measured-phase JIT reporting rule, declared before the first v2 candidate performance result while the first CONTROL/accumulator screen was still warming up.

The warmup gate alone does not establish that the five measured requests remain settled. The unchanged reviewer already writes every warm and measured request's `JIT_delta_ms`, `wall_s`, cohort, variant, and request number to `rows.csv`. Preserve and report **all measured JIT deltas** for each variant in every cohort, in request order; also retain the raw before/after counters in the per-request evidence.

For each cohort/variant batch of five measured requests, calculate:

`measured_batch_JIT_fraction = sum(jit_after_ms - jit_before_ms) / (1000 * sum(wall_s))`

Flag **strictly greater than 0.02** as `MEASURED_JIT_ABOVE_2_PERCENT`. Show the five deltas, summed JIT milliseconds, summed elapsed seconds, fraction, and flag beside the cohort effects and confidence intervals. Mark unavailable or decreasing raw counter intervals explicitly rather than treating them as zero. The JIT counter aggregates compiler work and is not necessarily elapsed compilation time.

This is a diagnostic annotation, not a new acceptance gate. Do not discard, replace, retry, or selectively omit requests, cohorts, or unfavorable outcomes because of this flag. Preserve the original protocol result, ratios, and uncertainty calculations. Inspect the ordered rows for a large measurement-phase JIT wave; such a wave prevents describing the observed elapsed difference as an established steady-state gain, even though the raw comparison remains reportable with its limitation.

The runner, warmup thresholds, byte checks, resource checks, selection plan, and existing rejection rules remain unchanged. This declaration applies uniformly to CONTROL, STOCK, and every candidate on both workloads.
