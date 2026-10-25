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

typed-full, live heap after types-to-colors: 172.5 -> 160.0 MB (with the paged original-name column).
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
| variables keep their CompilerInput as a 4-byte index (weak table) and TypedVar its type as a base pointer: Var 48 -> 32, TypedVar 64 -> 48 bytes | 215K typed variables at the typed-full peak | ~-2.4 MB at the peak, 32 MB less allocated |
| the original-name column is paged (64 values per page, allocated on first write) | typed-full: ~3 % of the nodes have one at the peak, most after renaming | -2.7 MB after colors, ~-3 MB at the peak |

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

## Checked and left (each about 1 % or inherent)

Measured on the head after the cuts above:

- **Dead nodes before compaction.** Compaction cannot run before
  types-to-colors (the type graph holds nodes). At the start of type
  inference, the nodes reachable from the root are 91 % of the store on
  typed-full (1.99 M of 2.18 M) and checks-only (1.38 M of 1.51 M), 82 % on
  typed-app. The rest includes the JSDoc type expressions, which are live,
  so there is less than 9 % to reclaim, and reclaiming it would need a
  compaction that forwards the type graph.
- **The node record** (32 bytes: four links, payload, property word,
  position, length, token, flags) has no field left to move out. Smaller
  would mean bit-packing the links, which every tree walk reads.
- **The global CFG** of typed-full (~107K nodes): nodes 8.7 MB (80 bytes),
  edges 5.2 MB (40 bytes), data-flow states 5.6 MB. Edge annotations hold
  the branched flow of type inference, so they stay on the edge.
- **TypedScope** structs: 6.7 MB at the typed-full peak (~88K scopes, 80
  bytes). **TypedVar**: 10.6 MB (48 bytes); the next size class (32) would
  need both the name string and the scope pointer gone.
- **Qualified names of object literal keys** (`owner + "." + key`, kept by
  their TypedVars): 6.6 MB on typed-full, 6 MB on checks-only. Packing each
  literal's names into one string would save the size-class rounding, about
  1 MB.
- **Parse trees** at the bundle peak: ~45 MB, already slab-allocated, lists
  copied out at their exact length.

## Next candidates

- **Convert each statement as soon as it is parsed, at any depth.** The
  bundle peak is `typescript.js`'s whole parse tree (~45 MB): the file is
  one IIFE, so converting top-level statements one by one does not help,
  but converting every statement of a function body once it is parsed would
  keep only one statement's parse tree alive (bundle peak ~-40 MB). The
  IRFactory attaches JSDoc with a forward-moving comment iterator and
  reports some warnings while converting, so children converted before
  their parent change both orders; it needs a position-keyed comment lookup
  and a check that warnings come out in the same order.
- Parse trees: the comma positions of argument lists (24 bytes per list)
  could be computed from the source when comments need them.
