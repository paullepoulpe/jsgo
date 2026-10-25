# Round 4: where the heap peaks, and three cuts (2026-10-24)

Builds on round 3 (32-bit node ids in per-compilation stores,
`round3-flat-slices.md`). Goal: lower the heap at its peak, which sets the
process's RSS (with the default GOGC=100 the heap may grow to about twice
the live heap of the previous cycle).

## Where the peak is

`GODEBUG=gctrace=1` (live heap after each cycle) and a heap profile at the
largest live heap (`JSGO_PEAK_HEAP_PROFILE`):

| Workload | Peak live | When | What is live |
|---|---|---|---|
| bundle (4 large untyped files) | ~135-140 MB | parsing, 0.8 s | parser objects (tokens, parse trees) ~70 MB, node blocks ~16 MB, source text ~8 MB |
| typed-full (closure library + app, ADVANCED, source map) | ~195-215 MB | type checking, ~4.3 s | node blocks 83 MB (2.2 M nodes) + type/color column 9.5 MB, types ~45 MB, typed vars 9 MB |

At the typed-full peak, 91% of the nodes are code reachable from the root,
7% JSDoc type expressions, 1% dead: there is little to reclaim, so the
record size is what counts.

On bundle the peak is the parse of `typescript.js` (one huge statement: an
esbuild bundle in a single IIFE), before its parse tree is converted to
nodes. Converting statement by statement would not help that file; making
the parser's tokens and trees smaller would (see "Next").

## Cut 1: compaction in place

Compaction used to copy the live tree into a second store, so both trees
were in memory during the copy. It now marks the live nodes (bit per id),
slides their records down within the store's own blocks (new id = number of
live ids before it), frees the blocks left empty, and forwards the
compiler's holders. The live nodes still get ids in a fresh range, so a
stale id fails like a nil pointer. The second store, the per-record
"forwarded" flag and the forced GC before the copy are gone.

typed-full: compaction takes 124 ms instead of 165, allocates 90 MB less,
one GC cycle fewer. The compaction was not the peak on these workloads
(see above), so peak RSS barely moved; it would be on workloads where the
AST dominates the heap after type checking.

Turning compaction off is not an option: typed-full peak RSS goes from
~410 to ~520 MB.

## Cut 2: 36-byte node record

| | before | after |
|---|---|---|
| payload | 64 bits: two packed ids (string + original name, cooked + raw) or float64 bits | 32 bits: one value per kind (string id, change time, BigInt id, raw string id, or a number: inline int31, else an index into a table of the other numbers) |
| original name | high half of the payload (string kinds), ext box (numbers, templates) | a kind flag when it is the node's own string (what SourceInformationAnnotator records for every name), else a per-block column |
| cooked template string | high half of the payload | ext box, only when it differs from the raw string |
| property bits | `uint64` (forces 8-byte alignment) | two `uint32` (same single 64-bit loads and stores) |
| record | 40 bytes | 36 bytes |

typed-full peak live (GOGC=25 for a dense sample): 212/221 MB -> 205/214 MB.
No time change. Output and source maps byte-identical.

## Cut 3: parser positions as offsets

Every token and parse tree had a `SourceRange` of two `SourcePosition`s,
each with an offset, a line and a column (24 bytes per range), and the
scanner tracked lines as it went (`LineNumberScanner`). A position is now
the offset alone; the line and column come from a table of line starts that
the `SourceFile` builds on first use (`SourceFile.LinePosition`, a cursor
for the forward scans plus a binary search). What needs lines keeps them:
comments (`LineRange`), parse errors (`LinePosition`) and `NonJSDocComment`.

| object | before | after |
|---|---|---|
| `SourceRange` | 24 B | 8 B |
| `IdentifierToken` | 48 B | 32 B |
| `IdentifierExpressionTree` | 40 B | 24 B |
| `BinaryOperatorTree` | 80 B | 64 B |
| `CallExpressionTree` | 56 B | 40 B |

bundle (GOGC=25, 3 runs): peak live 130-138 MB -> 105-113 MB, peak RSS
191-197 MB -> 165-168 MB, 46 MB less allocated. Output, source maps and
warnings (with their lines and columns) byte-identical.

The bundle peak is now the end of `typescript.js`'s parse: its whole parse
tree (~47 MB) plus the ASTs of the three other files, parsed in parallel
(node blocks ~36 MB). IRFactory already drops each statement's parse tree
once converted, but conversion starts after the parse, so only a smaller
parse tree (or AST) lowers that peak.

## A gap in the earlier gates

The output gate compared `out.js` and stderr but not the source map
(`out.js.map`, written by typed-app and typed-full). The first version of
cut 2 changed 44 mappings (an in-place rewrite read a node's original name
after clearing it, when the node was its own source reference); the gate
now compares the map too. The branch's maps were checked against main up to
cut 1: identical.

## Next

- **Parser objects at the parse peak** (bundle, after cut 3): identifier
  tokens ~9.5 MB, identifier expressions ~6 MB, `[]ParseTree` slices ~4 MB,
  argument lists ~4 MB. Storing the identifier's name in the tree node
  instead of a separate token object would save about 8 bytes and one object
  per identifier expression.
- **The record further**: 32 bytes needs another 4 bytes; the candidates
  (fewer inline boolean properties, the source-file id) cost more than they
  save today (about 2% of nodes would need a cold entry of ~100 B).
- **GC settings**: with GOGC=25, typed-full peaks at ~300 MB RSS instead of
  ~400 MB, for about 75 more GC cycles; a memory limit
  (`debug.SetMemoryLimit`) would bound RSS more directly.
