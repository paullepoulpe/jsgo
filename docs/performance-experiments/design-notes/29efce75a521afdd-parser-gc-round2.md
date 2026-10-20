# Direct parser and GC-policy screen, round 2

Frozen compiler sources: native80 `ef9657d7`, parser80 `11534351`. Every parser row uses the same binary; only GOGC differs. Private node batching is disabled. Go 1.27.2 and existing checked-in PGO throughout.

Three interleaved fresh CLI runs per cell, compiler threads 1, affinity 0,1 on a two-CPU cgroup. Each cell is median seconds / median peak RSS MiB. Every measured output matched pinned stock Java artifacts, diagnostics and exit status. Java is one stored JVM cohort (12 warmups, 10 measurements), not a fresh paired-Java confidence estimate. Java RSS covers the whole JVM lifetime including warmup; Go RSS is the median fresh-process peak.

| Workload | Warm Java s / MiB | Native analyses, default GC | + parser, default GC | + parser, GOGC=200 | + parser, GOGC=400 |
| --- | ---: | ---: | ---: | ---: | ---: |
| typed-app | 3.739 / 1090.2 | 5.790 / 254.6 | 5.802 / 256.1 | 5.433 / 304.3 | 5.541 / 485.8 |
| bundle | 16.600 / 1696.6 | 19.749 / 565.7 | 20.734 / 412.6 | 19.429 / 563.2 | 18.316 / 898.7 |
| downlevel-whitespace | 0.938 / 930.2 | 1.472 / 165.1 | 1.400 / 169.7 | 1.289 / 202.2 | 1.291 / 237.1 |
| downlevel-simple | 2.805 / 1118.4 | 4.153 / 203.0 | 4.261 / 200.1 | 3.860 / 248.3 | 3.847 / 339.1 |
| checks-only | 6.038 / 1360.0 | 8.276 / 658.7 | 8.176 / 652.3 | 7.748 / 882.2 | 7.398 / 1001.5 |

Fixed GOGC=200 ratios follow. This is an experimental configuration, not a proposed default. Java / Go values above 1 favor Go.

| Workload | Java/Go elapsed | Java/Go peak RSS | Native80 / candidate elapsed |
| --- | ---: | ---: | ---: |
| typed-app | 0.688 | 3.583 | 1.066 |
| bundle | 0.854 | 3.012 | 1.016 |
| downlevel-whitespace | 0.728 | 4.600 | 1.142 |
| downlevel-simple | 0.727 | 4.503 | 1.076 |
| checks-only | 0.779 | 1.542 | 1.068 |

## Interpretation

The parser alone has no broad speed win: typed-app is effectively unchanged, bundle is about 5% slower, whitespace about 5% faster, simple about 2.5% slower, and checks-only about 1.2% faster. Small differences remain screening evidence with only three repeats. Bundle RSS varies widely; inspect raw samples before claiming a memory win.

Raising GOGC buys speed with more memory. GOGC=400 is not uniformly better than 200. No row beats warmed Java on any of these five workloads. Native80 with GOGC=200/400 remains an unmeasured control: the combined rows cannot establish how much improvement requires the parser. Parser-only comparisons use identical default GC.

Removing transient ParseTrees may change GC pacing during later passes; this remains a hypothesis pending CPU/heap profiles. Separate non-tracing post-screen profiles were captured from the exact parser source. New fragment-config/comment refinements and work-elision prototypes are not included here.

## Validation

The full original parsing suite passed with the parser disabled and enabled. Exact AST/JSDoc/metadata/ordered-diagnostic oracles covered four bundles, all 897 core-js files, and 971 Closure files: 968 native acceptances, including 148 goog.modules. Full CLI gates passed all five tasks at threads 1 and 2.

## Provenance

- `native80`: `503ae718e8cd225a9c21269b14e791284d6a53bcf20e17bca9a0c9316274f2aa`
- `parser80`: `f9ad4316496736fa48dc871617680ac83588de78cdf992fd90b39ad5c8cc175d`
- `parser200`: `f9ad4316496736fa48dc871617680ac83588de78cdf992fd90b39ad5c8cc175d`
- `parser400`: `f9ad4316496736fa48dc871617680ac83588de78cdf992fd90b39ad5c8cc175d`

Raw per-run records, environment and JVM reference metadata: `/workspace/jsgo-results/round2/parser-gc-screen/`. Full parser validation: `/workspace/jsgo-results/round-parser-native80-validation.json`. CLI gates: `/workspace/jsgo-results/round2/gates-parser80.log`.
