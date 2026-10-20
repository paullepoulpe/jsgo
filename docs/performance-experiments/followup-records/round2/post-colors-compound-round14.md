# Round14: corrected post-colorization compound screen

Same binary source22358ec4, native9 + CLI source ownership control versus eight actual-admission corrected factors. Fresh non-worker CLI, default runtime GC, compiler threads1, affinity0,1 on two-CPU quota; Go1.27.2/shippedPGO/l4. Three alternating repetitions per variant; all30 measured outputs match pinned Java. All15 preflight Java artifact gates at threads1/2/4 match; focused compiler race tests1299PASS/2SKIP, no failures/races. Exact command and logs alongside this report.

Compound factors: native live-variable solver, compact inverter, first scope allocation, ASCII output, RUC native tasks + intrusive FIFO (FIFO requires tasks), exact FSIV use-cardinality, InlineFunctions candidate-discovery kernel. Generalized type recipe and fresh-template vector flag OFF.

| Task | Control s | Compound s | Time change | Control RSS MiB | Compound RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|---:|---:|
| typed-app | 5.076 | 4.892 | -3.62% | 276.0 | 246.7 | 0.764 | 4.42 |
| bundle | 18.407 | 17.813 | -3.22% | 421.1 | 434.2 | 0.932 | 3.91 |
| downlevel-whitespace | 1.230 | 1.231 | +0.05% | 173.7 | 159.6 | 0.762 | 5.83 |
| downlevel-simple | 3.747 | 3.641 | -2.84% | 212.2 | 203.7 | 0.771 | 5.49 |
| checks-only | 7.277 | 7.499 | +3.05% | 652.6 | 648.2 | 0.805 | 2.10 |

Ratios Java/Go, above1 favorsGo. Java is stored warm cohort with fresh Compiler/reparsed externs per iteration, not fresh paired confidence estimate. Go RSS median fresh-process lifetime peak; Java warm-JVM lifetime peak including warmup; windows differ. Three runs screen only; retain min/max in summary, do not infer confidence. No broad warm-Java parity achieved. Checks-only slowdown and bundle RSS increase require factor attribution; do not promote all factors as a universal win.

Activation evidence: post-colors-all-factor-activation.log and per-task post-colors-activation.json; all5 real inputs matched Java with positive backend counts where each phase is present. Counts/dense cumulative bytes are not RSS or speed. Earlier optimized-CLI negatives affected by the irrelevant runtime-library-table admission fence remain withdrawn; raw historical data retained.
