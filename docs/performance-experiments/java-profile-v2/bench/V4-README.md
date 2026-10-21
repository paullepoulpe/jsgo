Prospective source only. Read [CONFIRMATION-V4-PRESPECIFIED.md](../CONFIRMATION-V4-PRESPECIFIED.md) before choosing a run. No v4 JVM, benchmark, Docker action or result was created while preparing these files.

Files: `run-v4.py`, `review-results-v4.py`, and `check-v4-source.py`. The runner fixes direct published STOCK, five independent cohorts, 20 measured pairs, and 60–150 matched warmups with the prescribed eight-request windows. It uses the same CliFactors driver, helper pins, resources, outputs, provenance and teardown as v3. `runner_policy_version` remains `cache-aware-v3-PSIoptional` because that resource/pressure policy is unchanged; `benchmark_protocol.version` is separately `confirmation-v4`.

Source identities at closure:

| File | SHA-256 |
| --- | --- |
| run-v4.py | 9adb1c481ea19998d45a7e75a87e7670c1d5efafb381f2138c4d8f64bc18730d |
| review-results-v4.py | 5eae13886c251040b3d90762b69a4334891c30c3b6c7ece2aeb368b2f308c3d2 |

The existing v3 launcher is unchanged. The isolated `../launch-v4.py` pins this runner and creates the usual coordination receipt with the full v4 protocol, fixed STOCK/five cohorts/20 pairs/60–150 warmups, seed20261020 and request timeout300. Total timeout defaults to14400 and can be selected before execution with `--total-timeout`. BigQuery is the default; a bundle run requires explicitly selecting `bundle` and is never launched automatically. This is the same runtime coordination mechanism, with no extra approval layer.

Root's selected BigQuery confirmation command (starts a benchmark; not executed during source preparation):

```text
python3 launch-v4.py --manifest MATCHING_MANIFEST --candidate SELECTED_CANDIDATE --out FRESH_V4_RESULTS
```

The receipt includes runner/manifest/spec hashes, fresh output path, baseline/candidate, cohorts=5, pairs=20, min_warmup=60, max_warmup=150, seed=20261020, request_timeout=300, total_timeout and `benchmark_protocol` exactly equal to the literal `PROTOCOL` dictionary. The runner and reviewer verify those settings. No custom receipt is needed when using this isolated launcher.

The launch command otherwise retains v3's manifest/spec/oracle/oracle-metadata/jar/java-home/image/out/runtime-grant arguments and adds these explicit sampling values:

```text
python3 bench/run-v4.py [the pinned input/resource arguments]
  --baseline STOCK --candidate SELECTED_CANDIDATE
  --cohorts 5 --pairs 20 --min-warmup 60 --max-warmup 150
  --seed 20261020 --request-timeout 300 --total-timeout 14400
```

After a fully stopped completed run, review into a new directory:

```text
python3 bench/review-results-v4.py --results FRESH_V4_RESULTS --out FRESH_V4_REVIEW --manifest MATCHING_MANIFEST
```

The v4 summary status explicitly says `PASS warm-gated byteidentical Java-vs-Java v4`; it does not assert that the measured phase stayed settled. Review emits rows.csv with all raw JIT before/after counters and validity, cohorts.csv, ratios.csv using `log_t95_df4_*`, measured-jit.csv with all20 and four fixed5 blocks, and review.json. JIT >1%/>2% and first10/last10 drift are diagnostic; no adaptive rejection, exclusion, retry or replacement follows. Historical v3 reviews/results remain unchanged. The publisher now dispatches fixed v3/v4 counts using metadata's benchmark_protocol and verifies the matching runner; it never relabels or pools v4 results as v3.

Verification completed: all new sources parse; old runner/reviewer bytes retain their known hashes; every existing runner definition except stability/cohort_ci/main is AST-identical (including the full JVM class), and every other reviewer check helper is AST-identical. AST-isolated pure functions passed minimum-window/JIT/drift/spread/unavailable-counter checks, a masked first-block JIT wave case, the fixed-block CSV shape, a symmetric-log df4 interval reference, consistent protocol literals across launcher/runner/reviewer/publisher, and publisher rejection of mixed v3/v4 counts. The check script never imports a benchmark module or executes a JVM/Docker/helper, launcher or publisher. No end-to-end runtime validation is claimed.
