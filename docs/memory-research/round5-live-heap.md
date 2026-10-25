# Round 5: the live heap at the peaks (2026-10-25)

Builds on round 4 (`round4-heap-peaks.md`). The goal is the live heap
(what survives a collection), not RSS: peak RSS is about the live heap times
the GC's headroom, which is a setting.

How it was measured: `JSGO_MEMSTATS` (peak live heap, sampled after every GC
cycle; `GOGC=10` for dense samples), the live heap after the forced GC that
follows the types-to-colors step (typed-full; deterministic), and heap
profiles at the peak (`JSGO_PEAK_HEAP_PROFILE`).

## Results (against main at a8b05ce3)

| workload | peak live, main | peak live, now | allocated, main -> now |
|---|---|---|---|
| bundle (4 large untyped files) | 133-135 MB | 105-106 MB (-21%) | 1120 -> 1060 MB |
| checks-only | 162 MB | 147-149 MB (-9%) | 651 -> 551 MB |
| typed-full | 225-226 MB | 206 MB (-9%) | 2404 -> 2252 MB |

typed-full, live heap after types-to-colors: 172.5 -> 162.7 MB.
Output, source maps and stderr are byte-identical on all four gate
workloads after every step.

## The cuts

| cut | where it pays | measured |
|---|---|---|
| parser positions as offsets, lines from a line table (round 4, cut 3) | parse trees: 16 bytes less per tree and token | bundle peak ~-25 MB |
| 32-byte node record: the 64 property bits become a 32-bit word; the 32 rare boolean properties (each on under 1 % of the nodes; 15 never set) move to a side map | every node | typed-full -7 MB after colors, ~-9 MB at the peak |
| control-flow graph nodes keep up to two edges per direction inline, as pointers (144 -> 80 bytes) | typed-full peaks while inferring types in the global scope, whose CFG has ~100K nodes | typed-full peak ~-8 MB |
| identifier tokens keep a string-pool id (32 -> 20 bytes); the pool moved to `internal/strpool` so the parser can use it | ~420K tokens at the bundle parse peak | bundle peak ~-5 MB |
| the built-in externs' text is reloadable, so parsing frees it | every workload, all compile long | -1.7 MB |
| rare boolean properties in their own map (4-byte values) | nodes that have some | -1.4 MB |
| TypedScope's var index behind a pointer (96 -> 80 bytes) | ~90K typed scopes | ~-1.4 MB |

## Where the peaks are now

- **typed-full** (type inference of the global scope): node records 66 MB
  (2.1 M nodes), type and original-name columns 14 MB, typed variables
  14 MB (215K, most of them qualified names), typed scopes 9 MB, the global
  CFG and its data-flow states ~19 MB, types ~25 MB.
- **bundle** (end of `typescript.js`'s parse): its parse tree (~40 MB,
  identifier tokens and expressions the largest part), its source text
  (8.5 MB), and the ASTs of the three files parsed in parallel with it.

## Tried and dropped

- **Original names in the ext box instead of a column.** At the typed-full
  peak only 1 % of the column's slots are used, but renaming later gives a
  large share of the names one: the ext boxes cost 100 MB more allocation
  and no peak gain.
- **TypedVar names as string-pool ids.** Interning qualified names keeps
  them for the life of the process: +10 MB after types-to-colors, and the
  peak barely moved (the strings were alive there anyway).
- **Embedding the identifier token in IdentifierExpressionTree.** The
  scanner allocates the token before the parser knows what it is part of,
  so a copy saves nothing.
- **A memory limit (GOMEMLIMIT).** It cuts peak RSS a lot (typed-full
  473 -> 305 MB) but does not change the live heap; left as a setting
  (round 4, "GC settings").

## Next candidates

- TypedVar 64 -> 48 bytes (~-3.5 MB at the typed-full peak) needs 8 more
  bytes than the type pointer gives: the CompilerInput pointer as an index,
  for both Var and TypedVar.
- Pages for the original-name column (~-2 MB at the typed-full peak).
- Parse trees: the comma positions of argument lists (24 bytes per list)
  could be computed from the source when comments need them.
