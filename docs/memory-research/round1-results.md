# Memory round 1 (2026-10-23): results

Goal set by the owner: stop optimizing for wall time alone. For workloads
larger than the benchmarks, Go's GC cost grows with every long-lived pointer
it re-traces, so the targets were **peak live heap / RSS and pointer count
first, code simplicity second**, with wall-time regressions accepted unless
catastrophic. Default output stays byte-identical with Java.

Nine agents each took one area (parser, JSDoc, JSTypes, jutil collections,
scopes, per-pass pools/GlobalNamespace, CFG/data flow, rhino.Node, and
compilation-wide lifetimes); six research agents studied other compilers
(reports in this directory, synthesis in [README.md](README.md) and
[round2-design.md](round2-design.md)). Everything was merged and verified on
`claude/go-gc-overhead-opt-yleif9`.

## Results

Base = 8ae12b3 (Go 1.24.7). Final = all merges + Go 1.27.3. Linux, 4 vCPU,
`--num_parallel_threads=1`. **All 11 tasks byte-identical with Java** (memgate
MATCH). `go vet` and `go test ./...` pass except `TestInlineGolden`, which fails
identically at the base commit.

### Peak live heap (GOGC=10, so the per-cycle sample lands near the true peak)

| Task | Base MB | Final MB | Change | Pointerful ("scan") MB base → final |
|---|---:|---:|---:|---|
| bundle (12.6 MB, SIMPLE) | 333 | 128 | **−62 %** | 290 → 89 |
| bundle-advanced-ts | 345 | 154 | **−55 %** | 306 → 99 |
| bundle-advanced | 169 | 84 | **−50 %** | 143 → 61 |
| shaka-blockly | 234 | 119 | **−49 %** | 201 → 92 |
| checks-only | 362 | 190 | **−48 %** | 311 → 150 |
| typed-mega (44.6 MB) | 842 | 451 | **−46 %** | 744 → 325 |
| bundle-mega (36.1 MB) | 811 | 441 | **−46 %** | 516 → 243 |
| downlevel-whitespace | 82 | 44 | −46 % | 67 → 28 |
| typed-full (28.5 MB) | 507 | 280 | **−45 %** | 434 → 202 |
| downlevel-simple | 80 | 51 | −36 % | 60 → 22 |
| typed-app | 129 | 85 | −34 % | 104 → 57 |

### Default GOGC: RSS and wall time

Mean of 2 interleaved runs on an otherwise idle machine; all outputs MATCH.

| Task | Peak RSS base (Go 1.24.7) | base rebuilt with Go 1.27.3 | final | RSS change | Wall base | Wall final | Wall change |
|---|---:|---:|---:|---:|---:|---:|---:|
| bundle | 592 MB | 562 MB | 281 MB | **−53 %** | 11.9 s | 13.7 s | +15 % |
| bundle-advanced-ts | 604 MB | 548 MB | 286 MB | **−53 %** | 18.4 s | 20.4 s | +11 % |
| typed-full | 949 MB | 905 MB | 531 MB | **−44 %** | 30.2 s | 35.9 s | +19 % |
| typed-mega | 1,487 MB | 1,413 MB | 862 MB | **−42 %** | 52.7 s | 60.3 s | +15 % |
| bundle-mega | 1,408 MB | 1,362 MB | 856 MB | **−39 %** | 64.2 s | 73.4 s | +14 % |
| typed-app | 230 MB | 225 MB | 144 MB | **−37 %** | 3.9 s | 4.3 s | +12 % |
| checks-only | 424 MB | 412 MB | 314 MB | **−26 %** | 6.6 s | 7.2 s | +9 % |

The Go 1.27.3 toolchain alone (same code) gives −3 to −9 % RSS and −2 to −6 %
wall.

### googleapis: 44,240 files, 7,483 chunks, ~35 M AST nodes (staged, 1 thread)

macOS arm64, 8 cores, 24 GB, Go 1.27.1, single runs. The final build's
7,483 chunk outputs, both state files and the diagnostics are byte-identical
with the base build's.

| Stage | Peak RSS base → final | Wall base → final | CPU base → final |
|---|---|---|---|
| CHECKS | 8,468 → 5,935 MB (**−30 %**) | 262 → 138 s (**−47 %**) | 150 → 135 s |
| OPTIMIZATIONS | 5,808 → 5,110 MB (−12 %) | 48 → 71 s (+48 %) | 51 → 78 s |
| FINALIZATIONS | 1,778 → 1,367 MB (−23 %) | 10 → 11 s | 11 → 12.5 s |

Peak live heap during FINALIZATIONS: 1,125 → 697 MB (−38 %). This is the
"larger than the benchmarks" regime: base's CHECKS stage used only 150 s of
CPU in 262 s of wall, i.e. it stalled on memory at 8.5 GB RSS, and the
smaller final heap removes that. OPTIMIZATIONS pays for the extra allocation
(+27 s CPU), the round-2 target.

### Against `main` after the merge (Rounds 2P–2U, commit c3830bd)

The other optimization stream re-tuned its slabs in the meantime; the merged
branch keeps plain allocation and main's algorithmic changes. Same machine,
default GOGC, 2 runs each, all byte-identical with Java:

| Task | Peak RSS main → merged | Wall main → merged |
|---|---|---|
| bundle | 616 → 307 MB (−50 %) | 4.2 → 4.7 s (+12 %) |
| typed-full | 872 → 578 MB (−34 %) | 10.7 → 11.6 s (+9 %) |
| typed-mega | 1,333 → 890 MB (−33 %) | 18.3 → 20.2 s (+10 %) |
| bundle-mega | 1,300 → 912 MB (−30 %) | 23.0 → 26.4 s (+15 %) |
| typed-app | 234 → 176 MB (−25 %) | 1.7 → 1.5 s (−12 %) |
| checks-only | 397 → 344 MB (−13 %) | 2.3 → 2.4 s (+4 %) |

googleapis, single runs, all 7,483 chunk outputs and stage logs
byte-identical across the three (peak RSS / wall):

| Stage | Java (`-Xmx16g`) | main | merged |
|---|---|---|---|
| CHECKS | 9,182 MB / 248 s | 6,987 MB / 95 s | 5,677 MB / 75 s |
| OPTIMIZATIONS | 10,569 MB / 165 s | 5,781 MB / 63 s | 4,965 MB / 50 s |
| FINALIZATIONS | 5,772 MB / 49 s | 1,742 MB / 8 s | 1,392 MB / 11 s |

At that scale the smaller heap wins wall time back despite the higher CPU:
1.9× less memory and 3.4× faster than Java (main: 1.5× and 2.8×).

### Cost

- Allocation volume rose (typed-full 1.38 → 2.33 GB, objects 4.7 M → 32 M):
  slabs, arenas and pools were replaced by plain allocation, as Java does.
  GC CPU roughly tripled (typed-full 1.3 → 4.7 s, bundle-mega 2.2 → 9.7 s)
  and total CPU rose 21–29 %.
- Wall time +9 to +19 % (table above). The same code on Go 1.24 was slower
  still.
- **The cost is allocation itself, not GC.** Raising GOGC buys back almost no
  wall time but gives back much of the RSS (final binary, single runs):

  | Task | GOGC=100 | GOGC=150 | GOGC=200 |
  |---|---|---|---|
  | typed-full | 531 MB / 35.9 s | 630 MB / 35.8 s | 715 MB / 35.0 s |
  | typed-mega | 862 MB / 60.3 s | 1,036 MB / 61.6 s | 1,244 MB / 60.7 s |
  | bundle-mega | 856 MB / 73.4 s | 1,069 MB / 70.5 s | 1,124 MB / 75.3 s |

  GC marking runs on otherwise idle cores; what the main thread pays for is
  ~7× more allocated objects (as the GC lab found, allocation is 12–15 % of
  CPU whatever the GC settings). So the default GOGC stays at 100, and the
  lever for wall time is allocation volume: round-2 item 1.
- Code: 49 commits, **net −6,991 lines** (+3,777 / −10,768 over 231 files);
  `sync.Pool` uses 33 → 10.

## What changed, by area

- **Collections (jutil):** LinkedHashMap / ObjectLinkedHashMap index is a
  pointer-free open-addressing `[]uint32` table instead of a Go map (5–16 B
  per entry instead of 25–60, never scanned); slot layout fixed for sets
  (16 → 8 B). Slab chunks fill their size class.
- **AST (rhino.Node):** nodes allocated one by one (no slabs: a live node kept
  its whole chunk and every dead neighbour alive); overflow lengths and
  original names moved from a global `map[uintptr]` (leaky; could hand a dead
  node's data to a new node at the same address) into the node's ext box;
  qualified-name strings no longer pin 64 KiB slabs; string-pool tables 3/4
  load.
- **Parser:** SourcePosition/SourceRange are pointer-free values; ParseArena and
  its pools deleted; all-ASCII files need no UTF-16 copy; `self` back-pointers
  removed (85 As* methods deleted); each statement's parse tree is freed once
  converted (typescript.js parse tree 241 → 90 MB).
- **JSDoc:** builder slabs deleted (a fresh 32 KiB `[]any` chunk per file held
  46.7 MB for ~1 MB of values at typed-full's peak); `JSDocInfo.Property` is
  one concrete struct (no interface dispatch).
- **JSTypes:** registry slabs deleted; `propertySymbolList` (a package-level
  list that kept the whole registry alive after types→colors and across
  compilations) is per registry; `JSTypeBase` 64 → 32 B (no `self` interface,
  no per-type registry pointer, serial ids moved to the serializer);
  FunctionType 224 → 160, NamedType 208 → 176, InstanceObjectType 152 → 104,
  UnionType 128 → 96, Property 40 → 32.
- **Scopes:** Scope/Var/TypedScope/TypedVar/Reference allocated plainly; the
  scopeArena / typedScopeArena / ReferenceCollector arenas and their pools,
  reset-to-base machinery and per-pass hooks deleted (−2,700 lines);
  TypedScope vars in a slice with a pointer-free index.
- **Per-pass state:** GlobalNamespace, RemoveUnusedCode, InferConsts, the
  inliners, OptimizeCalls and renaming no longer keep pooled/cached state
  after their pass (the GlobalNamespace from processDefines held +73 MB on
  top of typed-full's peak); GlobalNamespace Name 200 → 160 B.
- **CFG / data flow:** the global CFG caches that pinned TypeInference's flow
  scopes and every TypedScope through the checks peak are gone; flow states in
  one slice; one documented per-owner `reuseSlab` replaces several arenas.
- **Lifetimes:** dependency info no longer pins whole source texts through
  substrings; report/source-map texts dropped after use; source-map mappings
  delta-encoded (~10 B instead of 32 B each); runtime libs read in place.
- **Toolchain:** Go 1.27.3 (`toolchain` directive).
- **jutil.Slab** remains only for objects that die together (one pass's Vars,
  union-find nodes, deserializer protos), documented as such.

## Lessons

1. **Slabs and pools were the main source of excess heap.** In Go a pointer to
   any element keeps the whole chunk alive and scanned, and a pointerful
   allocation over 512 B carries an 8-byte header that pushed power-of-two
   chunks into the next size class (128 × 64 B nodes cost 74 B each). Small
   objects allocated one by one have no header and die individually.
2. **Cross-pass retention dominated the typed peaks**, not object sizes: global
   CFG caches, cached GlobalNamespace arenas, a package-level symbol list, and
   dead typed nodes pinned by slab neighbours each kept tens of MB (and the
   whole type registry) alive.
3. **Measure peak live at a low GOGC.** At the default GOGC a run gets ~10 GC
   samples and the "peak live" reading varied by ±50 MB between identical
   runs; GOGC=10 is stable to ~±10 MB.
4. Heap fragmentation (holes in partially used spans) stayed at 2–7 % through
   the pipeline, so rebuilding the AST through Go's allocator to "defragment"
   is not worth it; see [round2-design.md](round2-design.md).

## Tooling added

- `JSGO_MEMSTATS=<file>`: alloc bytes/objects, GC cycles/CPU, peak live heap,
  scannable heap and live objects at the peak (`cmd/jsgo/memstats.go`).
- `JSGO_PEAK_HEAP_PROFILE=<file>`: heap profile at the live-heap peak.
- `tools/bench/memgate.py`: byte-identity gate against stored Java outputs plus
  those numbers, per task; `--gogc 10` for stable peaks.

## Next (round 2), ranked

1. **Recover the allocation churn** that plain allocation added, without
   bringing back retention: Var 40 → 32 B (per-pass `nameHash` scratch into
   pass-local tables), Reference 32 → 24 B, pre-sized root-CFG maps, no heap
   token for scanner punctuators, measure a growth policy for long
   `ArrayList`s.
2. **typeOrColor as a 32-bit handle** into registry/color tables: −20 % AST
   pointer words, and the type table can be dropped wholesale at colors.
3. **Remaining after-colors retention:** ChangeTracker timeline and ModuleMap
   keep detached typed nodes; TypedScope `reservedNames`/`module` only needed
   while building (6.7 MB at bundle-advanced's peak).
4. **typed-full's peak is now processDefines_CHECK's GlobalNamespace**: a
   defines-only namespace for ProcessDefines (−10–15 MB est.).
5. **Colors as compact id lists** (15–20 MB → 3–5 MB on typed-full post-check).
6. **Parser:** offset-only ranges plus a line table (~20 MB on large files).
7. **Input source-map parsers** keep ~45 MB through output on bundle-mega.
8. **Pointer-free AST with compaction** (see [round2-design.md](round2-design.md)):
   only together with a complete census of cross-pass `*Node` holders.
9. **googleapis with Java and `--num_parallel_threads=4`**, and repeated runs
   (the numbers above are single runs).

Fast-mode candidates (need relaxed semantics): hash-consing nullable/optional
unions and templatized types (70 % of live unions are duplicates); dropping
JSDoc type trees after checks; Go-map iteration plus sorting where order is
visible, instead of emulating Java HashMap order.
