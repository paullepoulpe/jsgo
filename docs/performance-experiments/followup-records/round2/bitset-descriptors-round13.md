# Round13: GC-safe descriptor BitSet screen

Independent fresh CLI comparison: original two-inline32-byte storage a6afc95f versus two-inline24-byte GC-tracked capacity-descriptor storage3c7ac3c0. Same Go1.27.2/shippedPGO/l4/default runtimeGC/native9+CLI ownership; native LVA and newer factors OFF. N1,affinity0,1,three rotating repetitions per variant per task. No build/test/profile jobs overlap timing. Full command and binary hashes are in bitset-descriptors-screen-{command,provenance}.json.

All30 measured artifacts match pinned Java; all15 N1/N2/N4 production compatibility gates match. Standalone jutil normal/race passes original vectors, frozen randomized/capacity/clone/nil oracle,24B layout/allocation tests and shared-descriptor concurrency tests. Two independent source reviews found no blocker. Full compiler race was not repeated specifically with this descriptor layout; no default promotion.

| Task | Variant | Median s | Range s | Peak RSS MiB | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|---:|
| typed-app | inline2 | 5.007 | 4.912–5.166 | 276.2 | 0.747 | 3.95 |
| typed-app | descriptors | 4.992 | 4.929–5.068 | 278.5 | 0.749 | 3.92 |
| bundle | inline2 | 18.267 | 18.251–18.453 | 435.2 | 0.909 | 3.90 |
| bundle | descriptors | 18.139 | 18.053–18.413 | 439.0 | 0.915 | 3.86 |
| downlevel-whitespace | inline2 | 1.221 | 1.169–1.229 | 170.6 | 0.769 | 5.45 |
| downlevel-whitespace | descriptors | 1.168 | 1.156–1.237 | 171.4 | 0.803 | 5.43 |
| downlevel-simple | inline2 | 3.824 | 3.635–3.906 | 209.2 | 0.733 | 5.35 |
| downlevel-simple | descriptors | 3.898 | 3.896–3.898 | 211.7 | 0.720 | 5.28 |
| checks-only | inline2 | 7.429 | 7.344–7.491 | 648.4 | 0.813 | 2.10 |
| checks-only | descriptors | 7.694 | 7.561–7.960 | 656.8 | 0.785 | 2.07 |

Ratios above1 favor Go. Java uses stored warm cohort with fresh Compiler and extern parsing per iteration; no new paired-Java confidence estimate. Go RSS is median fresh-process lifetime peak, Java RSS includes JVM warmup lifetime peak. These windows differ. Three observations are a screen, not a confidence interval.

## Interpretation

No broad speed/RSS win. Descriptor medians are typed0.3% faster,bundle0.7% faster,whitespace4.3% faster,simple1.9% slower,checks3.6% slower. Median RSS is higher on all five tasks in this cohort. Payload/object savings alone do not establish compiler improvement: descriptor checks/derived used length, GC scheduling and code placement also matter. Keep this experimental layout isolated rather than selecting its winning task as an overall improvement.

The post-colorization admission finding affects different optional native backends; it does NOT invalidate this unconditional BitSet representation comparison. Those backends were explicitly off here. See runtime-library-admission-audit.md for which earlier factor interpretations require retesting.
