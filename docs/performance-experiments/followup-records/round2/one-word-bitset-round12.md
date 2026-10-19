# Round12: one-word BitSet storage screen

Fresh non-worker CLI, default runtime GC, native9 plus CLI source ownership; native live-variable solver OFF. Three rotating repetitions per variant, all five fixed tasks, compiler threads=1 and affinity0,1 on the two-CPU quota. Same Go1.27.2, shipped PGO and package -l=4 flags. Source commits and binary hashes: `one-word-bitset-screen-provenance.json`. Full command: `one-word-bitset-screen-command.json`.

All45 measured outputs match pinned Java. One-word source7e0ba5fb passes923 focused compiler race-test cases,2SKIP, no FAIL/RACE, and all15 Java artifact gates at compiler threads1/2/4. Standalone jutil normal/race passed original vectors and frozen differential oracle.

| Task | Layout | Median s | Range s | Peak RSS MiB | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|---:|
| typed-app | slice | 4.931 | 4.928–5.158 | 281.2 | 0.758 | 3.88 |
| typed-app | inline2 | 5.000 | 4.900–5.052 | 274.9 | 0.748 | 3.97 |
| typed-app | inline1 | 5.017 | 4.971–5.059 | 276.3 | 0.745 | 3.95 |
| bundle | slice | 18.522 | 18.465–18.611 | 431.3 | 0.896 | 3.93 |
| bundle | inline2 | 18.121 | 17.839–18.895 | 425.8 | 0.916 | 3.98 |
| bundle | inline1 | 18.660 | 18.545–19.310 | 431.7 | 0.890 | 3.93 |
| downlevel-whitespace | slice | 1.253 | 1.196–1.329 | 163.4 | 0.749 | 5.69 |
| downlevel-whitespace | inline2 | 1.238 | 1.178–1.241 | 168.6 | 0.758 | 5.52 |
| downlevel-whitespace | inline1 | 1.193 | 1.165–1.303 | 168.0 | 0.786 | 5.54 |
| downlevel-simple | slice | 3.706 | 3.616–3.719 | 206.1 | 0.757 | 5.43 |
| downlevel-simple | inline2 | 3.804 | 3.683–3.806 | 215.6 | 0.738 | 5.19 |
| downlevel-simple | inline1 | 3.700 | 3.685–3.806 | 211.6 | 0.758 | 5.28 |
| checks-only | slice | 7.451 | 7.427–7.841 | 657.2 | 0.810 | 2.07 |
| checks-only | inline2 | 7.422 | 7.402–7.433 | 663.4 | 0.814 | 2.05 |
| checks-only | inline1 | 7.487 | 7.338–7.779 | 671.1 | 0.806 | 2.03 |

Ratios are Java/Go; above1 favors Go. Java is the stored warm cohort (fresh Compiler and reparsed externs every iteration), not newly paired timing. Go RSS is the median fresh-process lifetime peak; Java RSS is the warm JVM lifetime peak including warmup. These memory windows differ. Three repetitions are a screen, not a confidence interval.

## Interpretation

The 24-byte one-word layout does not show a broad improvement over either original slice or 32-byte two-word storage. It is not promoted. Against inline2, one-word typed time is0.35% higher and bundle2.98% higher; their RSS is0.49% and1.39% higher. Whitespace is faster, but its inactive solver path and prior control noise mean this is insufficient to declare a storage win. Simple is also mixed against the original slice, and checks RSS rises.

The untimed distribution diagnostic (`bitset-distribution-original/report.md`) correctly predicted that most constructor results reserve at most one word, but method-event requested-byte models are not actual allocation/RSS estimates: headers can be stack-resident, Clone/growth/large holders affect costs, and representation branches/code layout change CPU/GC behavior. This negative screen is retained rather than recasting the event count as a performance result.

## Further distinct hypotheses

A separate 24-byte two-word layout using real GC-tracked immutable capacity descriptors is source-ready; it has not been measured. A separate native-liveness scratch reservation change remains unmeasured. Neither is included in this one-word result. Larger native type-recipe, exact FSIV use-cardinality and parallel CommonJS prototypes are being validated independently. No worker/preparse cache or runtime GC tuning is involved.
