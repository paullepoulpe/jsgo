# Native analysis and output storage screen

All70fresh CLI runs (seven variants × five unchanged tasks × two runs) matched pinned Java output, diagnostics, exit status and source maps. No worker or parsed extern cache. CPUaffinity0,1 under two-CPU quota; compiler jobs=1; runtime GOGC/GOMEMLIMIT/GODEBUG/GOMAXPROCS unset. All builds/tests/profiling stopped before timing.

CandidateCLIproduction b2e407f6, finaltest-only61431bde. Go1.27.2, shipped PGO, four package -l=4 flags. PriorCLI641cb85d. Explicit environments, hashes, execution order and every raw sample: native-storage-code-screen/. Each candidate uses prior nine native/inline facts plus CLI_SOURCE_LOCK ownership witness. Scope recipes stay off. No default is changed.

Factors: compact contextual rename shadow records; native two-word live-variable solver with lazy escape-sensitive gen/kill; firstVar/map/entry coallocation (144 bytes/three objects ->144bytes/one object on this platform); ASCII printer/join buffer with permanent UTF16 promotion. Compound enables all four. Same-binary base isolates flag effects; prior binary also captures unconditional field/code-generation differences.

No broad improvement is established. Compound is neutral/slower on typed/bundle/whitespace/checks, improves simple by3.5% vs same-binary base but only0.4% vs prior. Typed compound RSS rises10.1%, bundle4.6% vs same-binary base. Isolated apparent wins do not add together. Two samples are a screen, not a retained win or confidence estimate. Notably the inactive liveness flag appears4.2% faster on whitespace, calibrating sampling variation rather than a valid optimization gain.

Java/Go ratios above1favorGo. The Java cohort is older stored warm JVM data, not a fresh paired confidence estimate. Go RSS is median fresh CLI process peak; Java is whole warm JVM-lifetime peak including warmup. None beats warm Java.

| Task | Variant | Go s | Peak MiB | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|
| typed-app | prior | 5.053 | 270.2 | 0.740 | 4.04 |
| typed-app | base | 4.929 | 251.2 | 0.759 | 4.34 |
| typed-app | compact | 4.974 | 266.7 | 0.752 | 4.09 |
| typed-app | live | 5.118 | 280.4 | 0.731 | 3.89 |
| typed-app | scope | 4.944 | 277.5 | 0.756 | 3.93 |
| typed-app | output | 5.006 | 261.8 | 0.747 | 4.16 |
| typed-app | compound | 4.930 | 276.5 | 0.758 | 3.94 |
| bundle | prior | 18.768 | 432.1 | 0.884 | 3.93 |
| bundle | base | 18.530 | 415.6 | 0.896 | 4.08 |
| bundle | compact | 18.947 | 433.6 | 0.876 | 3.91 |
| bundle | live | 18.721 | 431.3 | 0.887 | 3.93 |
| bundle | scope | 18.743 | 437.6 | 0.886 | 3.88 |
| bundle | output | 18.906 | 435.7 | 0.878 | 3.89 |
| bundle | compound | 18.898 | 434.8 | 0.878 | 3.90 |
| downlevel-whitespace | prior | 1.228 | 169.5 | 0.764 | 5.49 |
| downlevel-whitespace | base | 1.224 | 169.8 | 0.766 | 5.48 |
| downlevel-whitespace | compact | 1.176 | 159.1 | 0.798 | 5.85 |
| downlevel-whitespace | live | 1.173 | 160.6 | 0.800 | 5.79 |
| downlevel-whitespace | scope | 1.259 | 168.1 | 0.745 | 5.53 |
| downlevel-whitespace | output | 1.209 | 160.9 | 0.776 | 5.78 |
| downlevel-whitespace | compound | 1.253 | 161.2 | 0.749 | 5.77 |
| downlevel-simple | prior | 3.744 | 215.2 | 0.749 | 5.20 |
| downlevel-simple | base | 3.864 | 213.4 | 0.726 | 5.24 |
| downlevel-simple | compact | 3.834 | 213.7 | 0.732 | 5.23 |
| downlevel-simple | live | 3.640 | 211.6 | 0.771 | 5.29 |
| downlevel-simple | scope | 3.688 | 211.4 | 0.761 | 5.29 |
| downlevel-simple | output | 3.595 | 212.6 | 0.780 | 5.26 |
| downlevel-simple | compound | 3.729 | 214.1 | 0.752 | 5.22 |
| checks-only | prior | 7.653 | 671.9 | 0.789 | 2.02 |
| checks-only | base | 7.755 | 677.4 | 0.779 | 2.01 |
| checks-only | compact | 7.677 | 685.2 | 0.787 | 1.98 |
| checks-only | live | 7.743 | 656.6 | 0.780 | 2.07 |
| checks-only | scope | 7.731 | 661.9 | 0.781 | 2.05 |
| checks-only | output | 7.605 | 652.1 | 0.794 | 2.09 |
| checks-only | compound | 8.074 | 656.6 | 0.748 | 2.07 |

Validation:923PASS,2SKIP in BOTH normal and race focused suites; no races. Includes exact liveness masks/hash/steps/escaped order/arguments transitions; missing-priority original fallback and exact partial failures; current rename observations and retirement; UTF16/WTF8/source-map/zero-value output behavior; scope identity/check/publication/capture guards. Actual native liveness inside N2/N4 worker jobs is observed with private probes and full serial snapshot parity. All five CLI tasks match Java atN1/N2/N4 with the four factors plus parallel coalescing. Jutil allocation test confirms3->1 objects, normal+race8checks passed.

Followups: native publication still allocates tiny backing words separately for every result. Invocation-owned exact-capacity word windows are the smaller V2 alternative. A standalone generic inline-small-BitSet representation is a separate larger alternative with empty/large-set memory tradeoffs. Neither followup is measured here; no GC tuning is involved.
