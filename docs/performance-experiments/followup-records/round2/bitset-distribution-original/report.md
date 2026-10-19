# Original-slice BitSet event distribution

Untimed diagnostic inventory, 2026-10-19. Original source f2c17056, shipped PGO and four package -l4 options, frozen overlay hashes verified before every task. Runtime uses default GC/GOMAXPROCS/GOMEMLIMIT/GODEBUG with only native9 + CLI_SOURCE_LOCK enabled. Native liveness is OFF. All five tasks reproduce every stored Java artifact byte for byte; owned-word helper calls are zero.

**Method/event counts are not mallocs, retained objects, bytes or RSS.** Counters include package initialization and copied constructor temporaries, omit bare zero-value creation, and add substantial diagnostic overhead. This binary must never supply performance measurements. No live BitSet/AST/scope references are retained by the counters.

| Task | Result cap0 | cap1 | cap2 | >2 | Clone0/1/2/>2 | 0→1 growth | First small→large promotions | Java gate |
| --- | ---: | ---: | ---: | ---: | --- | ---: | ---: | --- |
| downlevel-whitespace | 2 | 54 | 0 | 0 | 0/0/0/0 | 2 | 0 | MATCH |
| downlevel-simple | 60,961 | 348,016 | 0 | 0 | 18,032/59,490/0/0 | 22,746 | 0 | MATCH |
| typed-app | 39,108 | 195,141 | 0 | 0 | 8,427/37,879/0/0 | 14,663 | 1,134 | MATCH |
| checks-only | 2 | 30,836 | 0 | 0 | 0/0/0/0 | 2 | 0 | MATCH |
| bundle | 344,597 | 2,082,265 | 39,779 | 0 | 99,421/361,356/3,704/0 | 123,546 | 627 | MATCH |

Exact per-request capacities, Clone source-capacity / used-length / result-capacity triples and from/to growth capacities remain in each task's `counts-0.json`. Typed-app has no large constructor results but does promote 1,134 initially small sets to large storage; constructor bins alone would miss that cost.

An optimistic one-word constructor-plus-growth payload model versus original is `-8*C1 +24*C2plus -8*growth(0→1) +24*growth(0/1→large)`. This predicts about15.90 MiB fewer payload bytes for bundle,2.83 MiB for simple,1.58 MiB for typed,0.24 MiB for checks and448 bytes for whitespace. It is not a heap-allocation result: nonescaping temporaries, allocator behavior, array-cell headers, live graph retention and RSS need independent measurement. The byte opportunity is modest relative to whole-compilation allocations; removing small backing-object traffic may have a different CPU effect, which is unmeasured.

Conclusion: the distribution justifies validating and independently measuring the one-word24B alternative. It does not justify retaining it yet or claiming Java parity, a memory ratio improvement or any elapsed-time gain. The native-liveness128→actual scratch factor is separate and inactive in this inventory.

Artifacts: `provenance.json`, `summary.json`, `event-model.json`, per-task `gate.json`/`counts-0.json`, and original stdout/stderr/output/source-map files. The first discarded attempt had a Python null-array summarizer assumption (zero Clone events); that runner bug was fixed without changing the frozen overlay, binary or source. Failed-attempt evidence is preserved in the sibling `bitset-distribution-original-empty-clone-runner-fix` directory.
