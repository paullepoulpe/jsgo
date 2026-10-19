# Inline bitset and fresh-template vector experiment protocol

Runtime GC settings are defaults, with GOGC/GOMEMLIMIT/GODEBUG/GOMAXPROCS unset. No worker process, cached parsed externs or new GC tuning. All variants retain identical nine earlier native factors plus CLI_SOURCE_LOCK; newer unrelated factors are off. Compiler job count1 and CPU affinity0,1 under the2CPU quota. Timing starts only after all builds/tests/gates complete.

Six variants, five fixed tasks, three randomized/rotated fresh CLI runs per cell:

| Variant | Binary source | BitSet | Native liveness | Fresh template values |
|---|---|---|---|---|
| prior | b2e407f6 | old slice | off | absent |
| inline | a6afc95f | inline two words | off | absent |
| inline_live | a6afc95f | inline two words | on | absent |
| template_control | 3ebe6449 | inline two words | off | off |
| template | 3ebe6449 | inline two words | off | on |
| compound | 3ebe6449 | inline two words | on | on |

The representation changes are across binaries; live on/off and template on/off have same-binary controls. Do not attribute template_control/inline code-layout/noise changes to enabling the vector path. Native liveness is inactive on whitespace/checks-only; template inference is inactive in all three untyped tasks (bundle SIMPLE, downlevel-simple SIMPLE and whitespace). Do not call inactive-control fluctuations wins. This does not independently isolate every interaction of the full three-factor matrix.

Every measured output/diagnostic/exit/source-map digest must match pinned stock Java. Warm Java ratios use the older stored12-warmup/10-measure cohort, not a fresh paired confidence estimate. RatiosJava/Go above1 favorGo. Go RSS is median fresh CLI lifetime peak; Java RSS is entire stored warm JVM peak, including warmup. Memory ratios reflect these different measurement windows. Three observations per cell are a screen, not a confidence interval. No default promotion from a small noisy minimum.

Raw command: inline-bitset-template-screen-command.json. Production provenance: build-inline-bitset-provenance.json and build-bitset-template-provenance.json. All compiled with Go1.27.2, shipped PGO and same four-package -l=4 flags. Build-only low-GC settings control compiler RAM, never benchmark execution. No source changes while compiling.
