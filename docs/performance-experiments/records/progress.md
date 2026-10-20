# jsgo performance experiment log

Status: warm Java has not been beaten. Preliminary results, all on the same two-core affinity. No default GC changes.

Baseline: Closure Compiler v20261015 on JDK21.0.12.1; Java 12 warm-up iterations then 10 measured. Go1.24.7 with repository PGO. Three interleaved runs per Go candidate. Every measured artifact must match Java, including diagnostics, exits, JS and source maps.

| Experiment | Typed app | Bundle | Whitespace |
|---|---:|---:|---:|
| first-candidates / baseline | 6.094s | 23.529s | 1.475s |
| first-candidates / collections | 5.938s | 23.302s | 1.497s |
| first-candidates / dataflow | 6.022s | 22.731s | 1.530s |
| first-candidates / node-chunks | 5.622s | 22.032s | 1.384s |
| second-candidates / baseline | 6.090s | 23.418s | 1.469s |
| second-candidates / compact-parser | 6.155s | 23.543s | 1.517s |
| second-candidates / var-hash | 6.152s | 22.944s | 1.478s |
| second-candidates / combined | 5.709s | 20.881s | 1.439s |
| memory-envelope / default | 5.573s | 21.100s | 1.423s |
| memory-envelope / gc200 | 5.274s | 20.449s | 1.373s |
| memory-envelope / gc400 | 5.285s | 20.449s | 1.335s |
| memory-envelope / gc800 | 5.152s | 20.904s | 1.301s |
| memory-envelope / gcoff | 5.312s | 21.284s | 1.337s |

The parser compact-self experiment is rejected: no useful elapsed-time gain, despite smaller structs. Eager variable-name hashing is marginal; a lazy hash and removal of duplicate inheritance storage are under investigation.

The first persistent worker results passed artifact gates, but memory rose strongly between iterations. They are excluded as evidence of a retained optimization until allocation ownership is fixed. Partial AST chunks can connect independent compilations transitively; new compiler boundaries and privately allocated immutable cache templates are being tested.

New candidates: batched scopes/vars/references; one-entry token position reuse; an 80-byte AST node (previously 96); 64-byte Var (previously 80); conflict-record deduplication in contextual rename inversion. These require both performance and parity validation.

A newer Go toolchain alone gave only 0–3% and changes Unicode behavior. Unicode15 tables are pinned in a separate BSD-licensed package; jutil and Unicode package tests pass on Go1.27. Full compiler validation remains before using that toolchain.

Final validation must rerun stable warm Java, test held-out projects and the full ported suite, check memory across mixed compilations, and measure the final combined implementation. No unqualified Java victory may be claimed from a weaker warm-up, extra heap, cached whole-output replay, or a different workload.
