The closed v3 warmups suggest that v4 admission may need appreciably more warmup, or fail within its bound. They do **not** establish that more than 90 requests will be needed. No v4 execution or protocol change follows from this retrospective check.

Scope: nine explicitly named, independently reviewed BigQuery jobs through reference-list-two, comprising 11 paired cohorts and 22 JVM series. Only their retained warmups were read. The active bundle screen and measured rows were excluded. The calculation uses eight-request windows, checks the frozen v4 pure stability helper at eligible prefixes, and does not import or execute the runner module.

| Closed job / cohort | Warmups per variant | Final-eight JIT proxy: baseline / candidate | Lowest observed paired maximum JIT proxy after request 16 |
| --- | ---: | --- | ---: |
| accumulator screen | 54 | 2.0566% / 2.4166% | 2.2222% |
| identifier dispatch screen | 45 | 1.8252% / 1.7521% | 1.8252% |
| identifier dispatch confirmation 1 | 41 | 2.5175% / 2.9521% | 2.9521% |
| identifier dispatch confirmation 2 | 41 | 2.5189% / 2.5143% | 2.5189% |
| identifier dispatch confirmation 3 | 40 | 3.2264% / 2.8342% | 3.2264% |
| callback bypass screen | 45 | 1.4524% / 1.9349% | 1.9349% |
| leaf identifier screen | 53 | 3.8442% / 2.9051% | 2.1100% |
| prototype/scalar equality screen | 40 | 2.9129% / 3.3857% | 3.3857% |
| empty-parameters screen | 41 | 2.9624% / 3.0400% | 3.0400% |
| source-position screen | 41 | 2.8314% / 3.0200% | 3.0200% |
| reference-list screen | 62 | 1.8110% / 1.5942% | 1.8110% |

None of the 22 final individual windows met the <=1% criterion. Across all observed paired prefixes after request 16, no pair met all three v4 window thresholds, even when ignoring the minimum-60 rule. Only reference-list reached 60; its three eligible prefixes, 60–62, all failed the v4 gate. The historical v3 gate used four-request windows and a 2% threshold, so these facts do not invalidate its admitted results.

All series stop after 40–62 warmups because of the historical v3 stopping rule. That censoring leaves no observations at 90 or 150. Expecting automatic admission at 60 is unsupported, but estimating a pass probability or calling “>90 likely” would also exceed the evidence. The best observed paired JIT maximum was still about 1.81%; this is a feasibility warning, not a stationary floor or an extrapolated settling curve.

CompilationMXBean measures approximate accumulated compilation elapsed time, potentially summed across threads. A compilation spanning a paused container or a later counter update may contaminate its request delta; see [counter interpretation](../JIT-COUNTER-INTERPRETATION.md). Consequently, <=1% can reject a prefix for reasons beyond actual JIT CPU consumption. More warmups may help ordinary settling, but the records cannot guarantee admission or separate that effect from counter accounting. The frozen 60–150 rule, 1% threshold, five cohorts, 20 measured pairs and retain-all-failures policy remain unchanged. There is no post-result relaxation, v3 relabeling or new performance claim.

Reproduction: `python3 publication-preparation/analyze-v4-closed-warmups.py`. The small [JSON record](v4-closed-warmup-feasibility.json) preserves per-cohort windows, observed prefix coverage and limitations; it contains no active benchmark evidence. This analysis started no JVM, build or benchmark.
