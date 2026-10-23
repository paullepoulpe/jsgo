# Cluster: go-gc-techniques — what Go offers, and how memory-heavy Go systems keep GC cost down

Research only. Sources were read from primary code where possible: the Go 1.27.3 tree (sparse clone in
`research/src/goroot`, tag `go1.27.3`, the current release on 2026-10-23), the golang/website
release notes and blog posts (`research/src/gogc/*.md`), CockroachDB/Pebble (`research/src/gc-cockroach`,
`research/src/gc-pebble`), typescript-go and esbuild (`research/src/gc-typescript-go`, `gc-esbuild`),
and raw files of bigcache, freecache, fastcache, VictoriaMetrics, Prometheus, Badger and Ristretto
(`research/src/gogc/`). Size-class and map arithmetic: `research/scripts/sizeclass.py`, `mapcost.py`.

## 0. Summary (the five things that matter most)

1. **jsgo's parse-time AST slabs waste 11–16 % to size-class rounding.** Since Go 1.22, every
   pointerful object of 513 B…32 760 B carries an 8-byte malloc header (`internal/runtime/gc/malloc.go`,
   `MallocHeaderSize`, `MinSizeForMallocHeader = 512`). `jutil.Slab[rhino.Node]` grows 16→32→64→128
   nodes = 1 024/2 048/4 096/8 192 B, which become 1 152/2 304/4 864/**9 472 B** objects (13.5 % waste on the
   128-node chunk that holds most nodes). Same for `JSTypeExpression` slabs (1 024×24 B → 27 264 B, 9.9 %).
   Fix is a few lines (typescript-go already does it): allocate chunks with `slices.Grow([]T(nil), n)` and
   use the whole `cap` (147 nodes per 9 472 B chunk instead of 128), or size chunks as (class − 8 B)/sizeof.
   Estimate: **−18 to −22 MB on typed-full, −50 to −65 MB on bundle-mega, ≈ −0.35 GB on the 35M-node build**.
   Keeps every API, no output risk.
2. **On Go ≥ 1.26, re-test whether parse-time Node slabs pay off at all.** The Green Tea GC (default since
   1.26) batches marking only for objects of 16–512 B (`gcUsesSpanInlineMarkBits`); 1–8 KiB slab
   chunks use the old per-object path. Slabs also keep dead nodes and everything they point to alive
   while any neighbour is live (a removed function whose root shares a chunk with live code keeps its whole
   subtree). Allocating each `rhino.Node` separately (64 B class: no rounding, no tail waste, size-specialized
   malloc in 1.27) frees dead nodes exactly. The gclab's own locality data shows fresh per-node allocation
   walks as fast as slab layout (11.6–12.1 vs 12.9–14.9 ns/node). Linux numbers in the context were taken
   with Go 1.24.7 (no Green Tea): re-measure on 1.27.3 first.
3. **Pointer-free memory lowers GC CPU, not live bytes or RSS.** GC mark CPU ≈ allocated bytes ×
   scannable fraction of the live heap ÷ (GOGC/100). Peak heap ≈ peak live × (1 + GOGC/100 × f), with
   f ≈ 0.5–0.9 measured. noscan bytes still count toward the heap goal. So pointer-free work only lowers peak RSS
   through (a) smaller fields (int32 instead of 8 B pointers, StringID instead of string) or (b) spending
   the CPU saved on a lower GOGC. GOGC=50 would give **−17 to −23 % peak RSS** (typed-full 879→~715 MB,
   bundle-mega 1.3→~1.0 GB). Today that costs ~1.6–2× GC CPU; after a noscan AST it costs about what
   GOGC=100 costs today. **Arenas that never free nodes raise the peak** when the AST churns (gclab:
   typed-app allocates ~2.6× more nodes than are live), so a pointer-free AST needs a reclamation story.
4. **Add a GC policy layer** like cmd/compile and CockroachDB do: a starting-heap floor for small
   builds (cmd/compile uses 128 MiB, port of `base.AdjustStartingHeap`), a `--max_heap`-style flag →
   `debug.SetMemoryLimit` (soft, Java `-Xmx`-like) with GOGC raised to 200–300 only when a limit is set
   (CockroachDB's exact rule), env `GOGC`/`GOMEMLIMIT` always win, and later a `memory` policy with
   GOGC≈50.
5. **Measure pointer density per component**: add `/gc/scan/heap:bytes`, `/gc/scan/globals:bytes`,
   `/gc/heap/objects:objects` at the peak-live cycle to the existing `JSGO_MEMSTATS` hook. Count non-nil
   pointer words per component with a debug census. Count slabs kept alive by dead neighbours with
   `weak.Pointer` probes. Side finding: the string pool's global page directories put **2.5 MiB of pointer
   arrays in BSS**, rescanned every cycle (single-100 runs 1 899 cycles).

---

## 1. Go runtime facts that decide the design (Go 1.27.3 unless noted)

### 1.1 Versions
- Current release: **go1.27.3** (`git ls-remote` tags; 1.27.0 in Aug 2026). jsgo's go.mod says `go 1.24.0`;
  the Linux GC lab and heap numbers used **go1.24.7**; the README's macOS numbers used 1.27.1.
  GC/allocator features come from the toolchain's GOEXPERIMENT baseline, not the go.mod directive
  (`src/internal/buildcfg/exp.go`: `GreenTeaGC: true, SizeSpecializedMalloc: true, RandomizedHeapBase64: true`).
- 1.19 soft memory limit; 1.21 `/gc/scan/*` metrics; 1.22 allocation headers; 1.23 `unique`;
  1.24 Swiss-table maps, `weak`, `runtime.AddCleanup`; 1.25 Green Tea (experiment), container-aware
  GOMAXPROCS (CPU only; **no automatic memory limit from cgroups**); 1.26 Green Tea default + AVX-512 scan
  kernels, heap base randomization, more stack-allocated slices; 1.27 size-specialized malloc (“<80 B
  allocations up to 30 % cheaper, ~1 % overall”, `go1.27.md`).

### 1.2 What a GC cycle costs
- Pacer (`runtime/mgcpacer.go:1295`):
  `goal = heapMarked + (heapMarked + lastStackScan + globalsScan) × GOGC/100`, floor
  `heapMinimum = 4 MiB × GOGC/100`. **heapMarked includes noscan bytes**, so pointer-free memory does not
  change how often the GC runs, only how much each cycle costs.
- Steady state: cycles ≈ alloc_bytes / (live × GOGC/100); work per cycle ≈ scannable live; so
  **GC mark CPU ∝ alloc_bytes × scannable_fraction / (GOGC/100)**. This matches the gclab:
  typed-app 1.32 GiB allocated over ~36 MiB average live → 37 cycles.
  For a fixed RSS budget (fixed GOGC), the only levers are allocation volume and pointer density.
- ~90 % of GC cost is marking; ≥35 % of mark time is stalled on memory (Green Tea blog,
  `blog-greenteagc.md`). Mark cost scales with the **number of non-nil pointer words** visited (each needs
  findObject + mark-bit test), plus a per-object cost. Nil pointer words are cheap; pointers to noscan
  targets are cheap (mark only, no enqueue).
- `/gc/scan/heap:bytes` = “live heap omitting no-scan objects and **no-scan tails** of objects”
  (`heapScan`, mgcpacer.go:198). The GC stops scanning an object at its last pointer word (GC guide:
  “group pointer fields at the beginning”). `rhino.Node` already does this: 5 pointer words first, a
  40 B scan prefix out of 64 B.
- GC CPU limiter: when the memory limit forces cycles, GC CPU is capped at ~50 % over a
  2×GOMAXPROCS CPU-second window (GC guide). A limit set too low gives up to a 2× slowdown and lets the
  heap exceed the limit; it never fails hard.
- jsgo's measured split (gclab, Go 1.24.7): GC 22–29 % of CPU but 3–7 % of wall (mostly runs on idle
  cores); mallocgc 12–15 % of CPU; live heap 80–96 % pointerful.

### 1.3 Green Tea (default since 1.26)
- Span-centric marking: when a pointer is first found, the runtime sets a mark bit and queues the
  **span** (FIFO). Marked objects build up on the span, which is then scanned in one pass with inline
  mark/scan bitmaps (`runtime/mgcmark_greenteagc.go` header). AVX-512 “expander” kernels handle packed
  spans (`internal/runtime/gc/scan/scan_amd64.go`) on Ice Lake / Zen 4+.
- **Scope: objects of 16–512 B only** (`gcUsesSpanInlineMarkBits(size) = heapBitsInSpan(size) && size >= 16`).
  Objects with a malloc header (>512 B) and large objects keep the old graph-flood `scanObject` path.
- Reported effect: GC CPU −10 % to −40 % (modal −10 %), about another −10 % with AVX-512 (go1.26.md).
  Worst case: one object per span per pass (the runtime special-cases it).
- For jsgo: individually allocated Nodes (64 B), JSTypes (64–224 B), Vars and Scopes all qualify.
  Node slab chunks (1–32 KiB) do not, nor do LinkedHashMap slot arrays over 512 B. Slabs already scan sequentially
  (8 KiB at a time), so the difference may be small; it is not known until measured on 1.27.

### 1.4 Size classes, the malloc header, and jsgo's struct and slab sizes
- 68 size classes (`internal/runtime/gc/sizeclasses.go`). jsgo's struct sizes 64, 80, 96, 144, 176,
  208, 224 B are exact classes (span tail waste 0/32/32/128/96/80/128 B per 8 KiB span, ≤1.6 %).
  **152 B (InstanceObjectType) rounds to 160 (5 % waste)**: either add 8 B of useful fields or trim to 144.
- Objects of ≤512 B keep their pointer bitmap at the end of the span (no header). **Pointerful objects
  of 513 B…32 760 B get an 8-byte type header**, and the size class is chosen for size+8
  (`runtime/malloc.go:mallocgcSmallScanHeader`). Objects >32 760 B are “large”: no header, rounded to
  8 KiB pages, one span each. `roundupsize` (`runtime/msize.go`) gives the usable size, and
  `append`/`slices.Grow` use it to set cap.
- Slab chunks of power-of-two byte size are the worst case:

| chunk (elem × n) | requested | allocated | waste |
|---|---|---|---|
| Node 64 B × 16 | 1 024 | 1 152 | 11.1 % |
| Node × 32 | 2 048 | 2 304 | 11.1 % |
| Node × 64 | 4 096 | 4 864 | 15.8 % |
| **Node × 128 (jutil.Slab max)** | 8 192 | **9 472** | **13.5 %** |
| Node × 512 (prebuild fixed slab) | 32 768 | 32 768 (large) | 0 |
| Node × 127 / × 511 | 8 128 / 32 704 | 8 192 / 32 768 | 0.8 / 0.2 % |
| JSTypeExpression 24 B × 1024 | 24 576 | 27 264 | 9.9 % |
| JSTypeExpression × 1023 | 24 552 | 24 576 | 0.1 % |
| JSDocInfo 32 B × 1024 | 32 768 | 32 768 (large) | 0 |

  Best 8 KiB / 32 KiB chunk lengths (pointerful elem): 24 B→341/1365, 32→255/1023, 40→204/819,
  48→170/682, 56→146/585, 64→127/511, 80→102/409, 96→85/341, 112→73/292, 136→60/240, 152→53/215,
  176→46/186, 208→39/157, 224→36/146.
- The heap profiler records the rounded size (`profilealloc(mp, x, size)` after rounding), so jsgo's
  measured “AST 32 % / 160 MB” already includes this waste. It shows up as `inuse_space/inuse_objects
  = 9472` at `jutil.(*Slab).grow`.

### 1.5 Interior pointers and slabs
- Any pointer into an object keeps the **whole object** alive, and the whole object is scanned,
  including slots of dead nodes. Their pointers (children via `first`, `typeOrColor` → JSType/Color/JSDoc,
  ext box → property lists) keep further objects alive, transitively. A detached subtree whose root shares
  a 128-node chunk with one live node survives in full. In ADVANCED builds that delete much code
  (`removeUnusedCode`, inlining), this can retain a large share of removed code. Its size is unmeasured;
  see §3 R7 for a weak-pointer census.
- The `arena` experiment's own docs say the same: “Any pointer into a chunk will therefore always cause
  the whole chunk to be scanned while its corresponding arena is still live” (`runtime/arena.go`).
- Tiny allocator: noscan objects under 16 B are packed into 16 B blocks; one live tiny object keeps its block alive.

### 1.6 Maps (Swiss tables, 1.24+)
- Groups of 8 slots plus an 8 B control word; max load 7/8; a table holds at most 1 024 slots, then the
  directory splits; `Map` header 48 B; **the first insert allocates a full 8-slot group**
  (`internal/runtime/maps/{map,table,group}.go`). Bytes per entry (`scripts/mapcost.py`, keys/values
  inline, string bytes excluded):

| map | n=1 | n=8 | n=100 | n=10⁴ | n=10⁵ | scanned? |
|---|---|---|---|---|---|---|
| `map[string]int` | 256 | 32 | 35 | 44 | 35 | yes (key ptr) |
| `map[*Node]int32` / `map[*Node]*X` | 192 | 24 | 24 | 30 | 24 | yes |
| `map[int32]int32` | 128 | 16 | 12 | 15 | 12 | **no** |
| open-addressing `[]int32` index (load ≤ 0.75) | – | – | 10 | 6.6 | 10.5 | **no** |
| … plus a cached 32-bit hash per slot | – | – | 20 | 13 | 21 | no |

- A map whose key and value types contain no pointers is never scanned (since Go 1.5, issue 9477; the
  basis of bigcache). jsgo has 57 `map[*rhino.Node]…` and 112 `LinkedHashMap/Set[*rhino.Node…]` sites;
  all are scanned and keep their keys alive.

### 1.7 `unique` (1.23), `weak` and `AddCleanup` (1.24)
- `unique.Handle[T]` is `struct{ value *T }`: an 8 B **pointer**, scanned. The canonical map is a concurrent
  hash-trie with 16-way indirect nodes. Each entry holds a cloned value, a `weak.Pointer` (a runtime
  “special” plus an 8 B handle object) and a cleanup to delete the entry (`src/unique/{handle,canonmap}.go`).
  That is roughly 100 B or more of overhead per distinct value, and lookups hash the full value.
  **Worse than jsgo's existing noscan pool with 32-bit IDs** for every hot use; fine for low-cardinality
  cold values.
- `weak.Pointer[T]` (comparable, survives the object's death as nil): built for caches and
  canonicalization maps. For a batch compiler its most useful role is **diagnostics**: probe whether a
  buffer or slab was really freed (§3 R7).

### 1.8 Arenas, regions, freegc
- `GOEXPERIMENT=arenas` still exists (off) in 1.27. The proposal (#51317) is **on indefinite hold**: arenas
  need an extra parameter threaded through every API, and they do not compose with escape analysis
  (memory-regions design doc, `research/src/gogc/regions.md`). Measured value was 10–15 % on real
  apps. Arena chunks are 8 MiB, scanned whole while any pointer into them is live.
- “Memory regions” (`region.Do(f)`, #70257) is a design doc only; nothing ships in 1.27. If it ships,
  per-pass `region.Do(pass.Process)` fits jsgo's pass-local churn exactly.
- `runtime.freegc` (GOEXPERIMENT=runtimefreegc) is compiler-internal (“The intended caller is the
  compiler”), off by default and not callable by user code.

### 1.9 Rules for pointer-free (“noscan”) memory
- noscan objects: allocated from separate span classes, never scanned (marking = set a bit; Green Tea
  fast-tracks them), but **still counted in heapMarked and the heap goal**.
- Storing Go pointers in memory the GC thinks is noscan (for example a `[]byte` or `[]uint64` chunk cast to
  `*Node`) is legal only if every target is kept alive some other way, and it relies on Go's heap never moving
  objects (true in every release; never formally promised). Write barriers are still emitted on typed stores
  (harmless). `typedmemmove` uses the static type, so copies stay correct. Untyped bulk clears skip barriers
  on noscan spans (`bulkBarrierPreWrite` → `typePointersOfUnchecked` returns nothing for noscan), which is
  fine only under the same rule.
- **`uintptr`-hidden pointers break `-race`.** `(*T)(unsafe.Pointer(u))` for a stored `uintptr u` is
  instrumented by checkptr (`cmd/compile/internal/walk/convert.go:walkCheckPtrArithmetic`) and throws
  “pointer arithmetic result points to invalid allocation” when `u` points into the heap with no original
  pointer (`runtime/checkptr.go:50`). Only conversions of function-call results are exempt. Casting a
  noscan chunk element (`(*Node)(unsafe.Pointer(&chunk[i]))`) passes checkptr.
- Go 1.26+ randomizes the heap base on 64-bit platforms. If jsgo builds with ≥1.26 and outputs stay
  byte-identical across runs, no output depends on addresses (useful when changing allocation patterns).

### 1.10 Metrics to measure with
`/gc/scan/heap:bytes`, `/gc/scan/globals:bytes`, `/gc/scan/stack:bytes`, `/gc/scan/total:bytes`,
`/gc/heap/live:bytes`, `/gc/heap/goal:bytes`, `/gc/heap/objects:objects`, `/gc/heap/allocs-by-size:bytes`
(histogram), `/memory/classes/heap/{objects,unused,free,released}:bytes`, `/cpu/classes/gc/mark/{assist,dedicated,idle}:cpu-seconds`,
`/gc/gogc:percent`, `/gc/gomemlimit:bytes` (`runtime/metrics/description.go`). `GODEBUG=gctrace=1`
gives per-cycle heap and CPU figures.

---

## 2. Techniques in memory-heavy Go projects

### Go compiler (cmd/compile)
- **Starting-heap boost** (`cmd/compile/internal/base/startheap.go`, `gc/main.go:88`): unless `GOGC` or
  `GOMEMLIMIT` is set (or `-c=1`), it raises GOGC so the first GC happens near **128 MiB** (`-d=gcstart=N`).
  A re-arming finalizer reads `/gc/heap/goal:bytes`, estimates live as `100×goal/(100+gogc)`, re-derates
  (70 % below GOGC 600, 55 % above), and restores the original GOGC once live passes half the target.
  It was tuned on large packages, **including `typescript-go/internal/ast`**, to “resemble
  GOMEMLIMIT=2000MiB GOGC=10000” in user CPU.
- `src.XPos` is 8 B: an `int32` index into a PosBase table plus `lico uint32` (20-bit line, 8-bit column,
  flags) (`cmd/internal/src/xpos.go`, `pos.go:319`). Positions carry no pointers.
- `ssa.Cache` keeps preallocated `[2000]Value`, `[200]Block`, `[2000]Location` reused across functions and
  clears only the used prefix (`ssa/cache.go`). It uses int32 IDs, sparse sets and maps, and bitvectors
  instead of maps.

### typescript-go (GC-level facts only)
- `core.Arena[T]` (`internal/core/arena.go`): bump allocation of chunks that double to **≤256 elements**,
  allocated with `slices.Grow[[]T](nil, n)` so the **chunk fills its whole size class**. That is the fix jsgo
  needs (§1.4). Used per parser, binder and checker for nodes, symbols, signatures, flow nodes and so on.
- Default Go GC: no `SetGCPercent` or `SetMemoryLimit` in the binary. `--extendedDiagnostics` runs
  `runtime.GC()` twice, then `ReadMemStats`. `ast.Node` is 48 B with `Parent *Node` and a 16 B `data` interface
  placed last (pointers not grouped first).

### esbuild
- `cmd/esbuild/main.go:323`: **`debug.SetGCPercent(-1)` (GC off) for non-watch, non-serve runs with
  at most one entry point**. jsgo's gclab measured `GOGC=off` without a limit at +86 % wall on bundle
  (page-fault cost on the bench VM). `GOGC=off` with `GOMEMLIMIT` set to a budget is the safe variant
  (gclab: −1.6 % to −11.5 % wall).

### CockroachDB / Pebble
- `pkg/cli/start.go:533-610`: `--max-go-memory` → `debug.SetMemoryLimit`. The default is 2.25 × the SQL
  memory pool, capped at 0.9 × system/cgroup memory − 1.15 × cache, with a floor of 256 MiB. Env `GOMEMLIMIT`
  is respected. **GOGC defaults to 300, but only when a soft limit is in force** (“to avoid introducing
  OOMs”). This is a direct template for a jsgo `-Xmx`-like flag.
- Pebble `internal/manual`: the block cache and memtables live in **C malloc memory** (cgo), off the Go heap
  and refcounted. Under `-race`, half the allocations randomly use Go memory so the race detector sees them.
- Pebble `internal/arenaskl`: the memtable skiplist sits in **one `[]byte` arena**. Nodes hold `uint32`
  key/value offsets and towers of `uint32` prev/next links, and are truncated to their random height.
  Nodes are cast from arena bytes; they contain no Go pointers, so this is legal.

### Badger / Ristretto (Dgraph)
- Badger `skl/skl.go`: the same arena skiplist (value packed as offset|size in one `uint64`, `keyOffset uint32`,
  `tower [maxHeight]atomic.Uint32`).
- Ristretto `z.Calloc` (jemalloc via cgo, build tag) and `z.Allocator`: a bump allocator over Calloc'd
  buffers that are “safe to unsafe cast to Go struct pointers” only for pointer-free structs. Dgraph's
  blog itself recommends against this unless necessary, and the memory is invisible to pprof.

### bigcache / freecache / fastcache
- bigcache: per shard, `map[uint64]uint64` (hash → offset) plus a `BytesQueue` byte ring. The README says
  “if map without pointers in keys and values is used then GC will omit its content”. 20 M entries: GC pause
  1.5 ms vs 9.3 ms for a plain map.
- freecache: 256 segments, each one ring buffer plus one `[]entryPtr{offset int64, hash16, keyLen uint16}`
  slot array: “only 512 pointers” in total. Its README also says: **lower GOGC (e.g. 20) when you
  allocate large amounts of noscan memory**. The heap goal scales with all live bytes, so with a mostly
  noscan heap, frequent cycles are cheap.
- fastcache (VictoriaMetrics): buckets of 64 KiB chunks allocated **off-heap via mmap** plus a
  `map[uint64]uint64` index (hash → generation|offset); O(chunks) pointers per bucket.

### VictoriaMetrics, Prometheus TSDB
- VictoriaMetrics `lib/bytesutil/internstring.go`: interning through a mutable `map[string]string` plus a
  read-only snapshot map with expiry, and `strings.Clone` on insert so interned strings don't keep large
  buffers alive. Simpler than jsgo's pool and pointerful. `lib/uint64set`: roaring-like buckets with a
  56-entry inline sorted `[]uint16` pool that switch to an 8 KiB noscan bitmap per 2¹⁶ range.
- Prometheus `model/labels`: the default `stringlabels` stores **one length-prefixed `string` per label
  set** (one pointer instead of 2 per pair). `dedupelabels` encodes varint symbol IDs into a shared
  `SymbolTable` (`byNum []string`, `byName map[string]int`). Part of the “How Prometheus Halved Its Memory
  Usage” series (Bryan Boreham, ~30 PRs).

### gopls (2023 redesign)
- v0.12 replaced whole-program in-memory type-checked state with **separate compilation**: per-package
  export data, xref and method-set indexes serialized to a file cache, loaded on demand. Memory fell by
  ~75 % on average, more on larger repos. It notes that “typed syntax trees are typically 30× larger than the
  source text” (jsgo: 552 MB / 28.5 MB ≈ 19× at typed-full's peak). The analog for jsgo is serializing
  stages it doesn't need resident, as the existing `LazyAst` typed-AST deserializer already does.

---

## 3. Recommendations for jsgo (ranked by peak live heap / pointer impact, then simplicity)

Numbers use the context file. typed-full: live 552 MB, RSS 879 MB, AST 160 MB. typed-mega: 980 MB /
1.5 GB. bundle-mega: 700 MB / 1.3 GB, AST assumed 60–80 %. 35M-node staged build: 8.7 GB RSS.

### R1. Make every slab chunk fit its size class — default mode, ~hours, no API or output risk
- **Where**: `internal/jutil/slab.go` (`grow`: `make([]T, n)` → `c := slices.Grow([]T(nil), n); s.chunk = c[:cap(c)]`);
  `NewFixedSlab` callers (`rhino/js_doc_info.go` `TypeExprSlab` 1024 → 1023 or 1365; the deserializer's
  512/256 slabs, by elem size); the hand-rolled `min(max(cap*2,16),256)` chunk growers (coalesce graph,
  color graph node factory, disambiguation, scope/var/RC/GlobalNamespace/inliner arenas); and
  `rhino.cloneTreeWithSlab` (exact-count slabs of 9–511 nodes: also round up and keep the slack for the next clone).
- **Mechanism**: avoid pointerful chunk sizes in (512 B, 32 760 B] that sit just above a class after the
  8 B header. Rule: use `slices.Grow` and take `cap`, or pick lengths from the §1.4 table. A unit test can
  assert `cap(slices.Grow([]T(nil), n))*unsafe.Sizeof(T) >= 0.97 × roundupsize` for every slab type.
- **Gain**: parse-time nodes sit mostly in 128-node chunks, so AST bytes −12 % (≥85 % of nodes in such
  chunks × 13.5 %): **typed-full −18–22 MB (−3.5 % peak live, ≈ −30 MB RSS); bundle-mega −50–65 MB
  (−7–9 %); 35M-node build ≈ −0.35 GB**. Also: with per-file IRFactory slabs (the serial path passes nil
  slabs, `parser_runner.go:122`), each file leaves a half-used last chunk (~4–5 KiB per file plus the JSDoc
  node slab). Use one compilation- or worker-scoped slab, as the parallel prebuild already does (512-node,
  32 KiB chunks).
- Keeps `*rhino.Node` and all pass code.

### R2. A/B the Node slabs against per-node allocation on Go 1.27 — default mode, ~1 day plus measurement
- **Why**: on ≥1.26 individual 64 B nodes get Green Tea batching (and AVX-512 on the bench CPU if it
  has it), 1.27 makes <80 B allocations up to 30 % cheaper, the 64 B class has zero rounding and tail
  waste, and **dead nodes are freed exactly** (no retention through chunk neighbours, §1.5). The gclab measured
  per-node “fresh alloc” walks as fast as slab order. Allocation cost: bundle's ~1.1 M parse nodes × ~15–20 ns ≈
  20 ms on a 30 s run.
- **How**: `IRFactory` already calls `f.nodeSlab.New().Init…FromTemplate(...)`, so a build tag or flag that makes
  `Slab.New` return `new(T)` gives the A/B with no pass changes. Compare peak live, `/gc/scan/heap`, GC
  CPU and wall (typed-app, bundle, N=10 paired) on **go1.27.3**. Re-run the gclab baseline on 1.27.3
  first, since Green Tea changes the numbers.
- **Prerequisite (output safety)**: `rhino.nodeExtraShards` is a `map[uintptr]nodeExtraEntry` keyed by node
  address and never cleared (`rhino/node.go:3108`). `setExtraLength` and `setExtraOrigNameID` read-modify-write the
  existing entry, so a new node at a reused address that sets one field inherits the dead node's other field
  (for example a stale original name on a NUMBER node, which would change source-map names). Slabs make address reuse rare;
  per-node allocation makes it common. Fix: when the node's `kindFlagHasExtra` bit is not set, start
  from a zero entry instead of `s.m[k]`, and delete entries when a node is re-initialized. The same caution applies to the
  other address-hashed caches (`type_inference.go:3230`, `control_flow_analysis.go:347`,
  `reference_collector.go:34`, `jstype_js_type_registry.go:332`): they are safe only if they compare the
  stored pointer or are pass-local.
- **Expected**: at least R1's gain (no rounding at all) plus whatever chunk-neighbour retention exists
  (measure with R7). GC CPU likely neutral to better on 1.27. If per-node loses on wall time, keep R1.

### R3. GC policy: starting floor, `-Xmx`-like flag, and a memory-first mode — default mode, ~1 day
GC settings never change output, so all of this is allowed in default mode.
```go
// cmd/jsgo, before compiling. Explicit env vars always win.
if os.Getenv("GOGC") == "" && os.Getenv("GOMEMLIMIT") == "" {
    if maxHeap > 0 {                     // --max_heap=8g | 75% (of cgroup memory.max or RAM)
        debug.SetMemoryLimit(maxHeap)    // soft, like -Xmx but never an OOM error
        debug.SetGCPercent(gogcWhenLimited) // 200–300 now (CockroachDB uses 300 with a limit)
    }
    adjustStartingHeap(64 << 20)        // port of cmd/compile/internal/base/startheap.go
}
```
- **Starting floor** (64 MiB; cmd/compile uses 128): single-file builds run ~19 cycles per compile today
  (1 899 cycles over the 100-file run, GC 27 % of CPU). With a floor most of the early 4→8→16→32 MiB cycles
  disappear. Expect **GC CPU on single files to drop by more than half**, for +20–40 MB peak RSS per process
  (gclab GOGC=400 gives 381 cycles and 156 MB RSS, an upper bound). Large builds are unaffected, because
  the floor is released once live passes half of it.
- **`--max_heap`** (or document `GOMEMLIMIT`): the Java `-Xmx` analog. It is soft: below the live heap,
  the GC is capped at ~50 % CPU (≤2× slowdown) and memory overshoots instead of failing. If a hard limit is
  wanted, check `/gc/heap/live:bytes` after each GC in the existing finalizer hook and exit with a Java-like
  OOM message. A safety-net default of `0.9 × cgroup memory.max` (read `/sys/fs/cgroup/memory.max`) costs
  nothing until the limit is near and turns container OOM-kills into slower runs; Go does not do this by itself.
- **Memory-first mode** (`--gc_policy=memory`, or the default once the AST is noscan): GOGC=50.
  Predicted peak RSS = live × (1 + 0.5 × f) with f = RSS/live − 1 as measured:
  typed-full 552 × 1.30 ≈ **715 MB (−164 MB, −19 %)**; typed-mega 980 × 1.27 ≈ **1.24 GB (−17 %)**;
  bundle-mega 700 × 1.43 ≈ **1.0 GB (−23 %)**. Cost today: ~2× cycles, GC CPU 25 % → ~40 % of CPU,
  ≈ +3–7 % wall on an idle 4-core box (more on a loaded one). After Green Tea plus a noscan AST, roughly
  equal to today's GC CPU.
- **Throughput mode** (`--gc_policy=throughput`, requires a budget): `GOGC=off` + `GOMEMLIMIT=budget`
  (the gclab's “Java-RSS budget” row: 1–3 cycles, −1.6 % to −11.5 % wall). Never `GOGC=off` alone.

### R4. Pointer-free structures: GC-level guidance for the in-flight AST, CFG and collection work
- **Count non-nil pointer words, not bytes.** Node: up to 5 per 64 B (parent, next, previous incl. the
  circular first→last, first, typeOrColor). typed-full ≈ 2.2–2.5 M nodes ≈ 10–12 M pointer words per
  cycle. `previous`/`parent` are redundant for reachability, but Go cannot skip them without the
  hazards in §1.9 (rejected below).
- **noscan buys CPU; compact fields buy bytes.** An int32-linked node (4×int32 links, int32
  type/color id, payload, flags) is ~44 B, padded to the 48 B class (170 per span): **−25 % AST bytes** and 0 mark work.
  The same AST as a noscan slab with `*Node` fields keeps 64 B: 0 % bytes, all of the CPU win.
- **Reclamation is the hard part.** Any arena (noscan slab or `[]Node` + int32) frees nothing until the
  whole arena dies. The gclab saw typed-app allocate ~90 MiB of nodes for ~35 MiB live (96 B nodes then),
  so a never-free AST arena *raises* the peak on untyped ADVANCED builds where the optimizer churns the AST.
  Options: (a) arena only for parse-time nodes, with optimizer-created nodes in a second, GC-managed form
  (impossible with mixed `*Node` links unless every link target is an arena node); (b) epoch compaction:
  at a pass boundary, copy the reachable AST in pre-order into a fresh arena and remap every long-lived holder
  (roots, inputs, JSType `source` nodes, JSDoc, change-tracking, side tables). This is only feasible
  with int32 IDs and a registry of holders, and it also fixes locality (gclab: contiguous walks 2.5× faster);
  (c) a free list fed by explicit deletion points. Closure's `NodeUtil.markFunctionsDeleted` and
  `reportChangeToEnclosingScope` are not ownership-precise, so this is unsafe.
- **Minimal-unsafe variant that keeps `*Node`** (if chosen): chunks allocated as noscan (`new([N]nodeRaw)`
  with a layout identical to Node), cast per element; all chunks rooted in the compilation; `typeOrColor`
  moved to a **parallel scanned `[]unsafe.Pointer`** (8 B/node scanned instead of 40). A node needs its
  global index for that, which fits in 4 of the freed 8 bytes. Every Node allocation must go through the arena:
  35 `&Node{}`/`new(Node)` sites today, including `rhino.NewNode`. Never store a non-arena pointer in a node.
  Never free (see above). checkptr-compatible; `-race` must stay green.
- **Struct-of-arrays for everything else** (bigcache/freecache/fastcache/arenaskl pattern): bulk data in
  noscan slices, the few real pointers in one small parallel slice, int32 IDs between them; lifetime = the
  owning compilation or pass (no unsafe at all).
- **Interfaces in long-lived structs cost 2 words**: the itab word points outside the heap (cheap) but is
  still visited. `JSTypeBase.self JSType` (16 B) points the object at itself, and `registry` is the same
  pointer on every type: 3 of 6 pointer words in the 64 B base do no reachability work. Replacing `self`
  with a kind byte plus a container-of cast (base embedded at offset 0) saves 16 B and 2 words per type
  (part of the in-flight JSType shrinking).

### R5. Strings and IDs — default mode
- Keep the noscan string pool with 32-bit IDs. Don't adopt `unique` (§1.7: an 8 B scanned pointer per handle,
  ~100 B+ per entry).
- Push StringID into remaining string-keyed hot collections (`LinkedHashMap[string,*]` is the top
  allocation site): a `LhmSlot[string,V*]` slot is 24 B with 2 pointer words; `[StringID,int32]` is 8 B with none. Java
  `String.hashCode` order emulation can read a per-ID cached hash (4 B/ID in the pool). This changes the
  key type at call sites, so do it only for the top few sites.
- `strings.Clone` any substring that outlives parsing, or the whole source text stays alive (one interior
  pointer pins the whole string). Verify in a test: `w := weak.Make(unsafe.StringData(src))`, drop the
  compiler's references, `runtime.GC()`, expect `w.Value() == nil`.
- **String pool directories** (`rhino_string_pool.go:75`): 64 shards × (`[1024]*chunk` + `[4096]*page`)
  = **2.5 MiB of pointer arrays in BSS**, scanned as roots every cycle and added to the heap goal. Make
  them heap slices grown on demand. Small, but it matters for many-cycle single-file runs (~0.3 ms × 1 899 cycles).

### R6. Maps — default mode, follows the in-flight jutil work
- Replace Go maps used as indexes with noscan `[]int32` open addressing: 7–10 B/entry vs 24–44 B for
  pointer- or string-keyed Swiss maps, and no scanning (§1.6). Small Go maps are expensive (256 B for one
  string-keyed entry); jsgo's inline ≤8-entry LinkedHashMap already avoids that.
- `map[*Node]V` side tables (57) and `LinkedHashMap[*Node…]` (112): if node IDs appear (R4), turn
  pass-local ones into dense `[]V` indexed by ID, or noscan open addressing keyed by ID.

### R7. Measure pointer density and retention per component — do first, ~1 day
- Extend `cmd/jsgo/memstats.go`: at each cycle where peak live rises, also record
  `/gc/scan/heap:bytes`, `/gc/scan/globals:bytes`, `/gc/scan/stack:bytes`, `/gc/heap/objects:objects`,
  `/memory/classes/heap/unused:bytes`, plus GC mark CPU at exit. Report pointer density = scan/live at the
  peak.
- Per component: in a debug build, at the peak (after checkTypes) and mid-loop, walk known roots and count
  objects and **non-nil pointer words** per class: AST (walk `jsRoot`, `externsRoot`, JSDoc type
  expressions), types (registry), scopes, collections. Multiply by a reflect-computed scan prefix per
  struct type (offset of the last pointer word + 8; strings/slices/maps/pointers = 1 word, interfaces = 2).
  The gclab's `heapclass.py` per-site classification can then report scan bytes per class.
- **Chunk-neighbour retention**: in debug builds, keep a `weak.Pointer[rhino.Node]` to `&chunk[0]` of every
  Node slab chunk. At a safe point, `runtime.GC()`, count chunks still alive × chunk length, and compare with
  nodes reachable from the roots. The difference = dead nodes retained by live neighbours (plus their subtrees).
  This decides R2 and sizes the R4 reclamation problem.
- Rerun on go1.27.3 (Green Tea + AVX-512 if the bench CPU supports it) before deciding any AST redesign.

### R8. Retained pools and arenas
- `sync.Pool` contents survive at most two GC cycles (victim cache) but are live, scanned heap until then.
  The global `cachedScopeArenas`/`cachedRCArenas`/`cachedGNArenas`/`cachedReusableCfgs` atomics keep their
  arenas reachable from globals for the whole run, raising heapMarked and therefore the heap goal. The in-flight
  removal is right; with Go 1.27 allocation costs, per-pass arenas released at pass end (or plain
  allocation) beat cross-pass retention for peak memory.

---

## 4. Fast-mode ideas, and rejected ideas

**Fast mode only**: the GC techniques above need no relaxed semantics, since GC policy never changes output.
Candidates that do need `--fast_mode`: dropping data kept only for output parity (JSDoc description strings
when not printed, original-name side tables when no source map is requested), which lowers live bytes directly.

**Rejected or deferred**
- `GOEXPERIMENT=arenas`: on indefinite hold, may be removed; 8 MiB chunks scanned whole while any pointer is live; API threading.
- Memory regions (`region.Do`): not shipped as of 1.27; revisit if it lands (a natural fit for per-pass churn).
- cgo/jemalloc/mmap off-heap memory (Pebble, Ristretto, fastcache): same GC benefit as noscan Go slices,
  but needs manual free, is invisible to pprof, complicates the build, and still cannot hold Go pointers.
  Noscan Go memory is freed by the GC anyway.
- `unique` for interning: see §1.7.
- `GOGC=off` alone (esbuild): +86 % wall on bundle (page faults) and unbounded RSS.
- GC ballast (Twitch, 2019): obsolete since `GOMEMLIMIT`.
- Hiding back-pointers (`parent`/`previous`) as `uintptr` so the GC only traces first/next: a node reachable
  only through a hidden parent edge (a pass holding a child of a detached subtree) is freed while still
  referenced, which is nondeterministic use-after-free that tests won't reliably catch, and checkptr throws under `-race` (§1.9).
- `runtime.freegc`: compiler-internal and experimental; would need unsupported `go:linkname`.
- Weak-pointer caches in production: a batch compiler has no use for them beyond diagnostics (R7).

## Sources
- Go 1.27.3 source: `src/runtime/{malloc.go,msize.go,mgcpacer.go,mgcmark_greenteagc.go,arena.go,mbitmap.go,checkptr.go}`,
  `src/internal/runtime/gc/{sizeclasses.go,malloc.go}`, `src/internal/runtime/maps/`, `src/unique/`, `src/weak/pointer.go`,
  `src/runtime/metrics/description.go`, `src/cmd/compile/internal/{base/startheap.go,gc/main.go,ssa/cache.go,walk/convert.go}`,
  `src/cmd/internal/src/{xpos.go,pos.go}`, `src/internal/buildcfg/exp.go` — https://github.com/golang/go/tree/go1.27.3
- Release notes go1.23–go1.27 and blog posts (greenteagc 2025-10-29, swisstable 2025-02-26, unique 2024-08-27,
  gopls-scalability 2023-09-08), GC guide — https://github.com/golang/website/tree/master/_content
- Memory regions design — https://github.com/golang/proposal/blob/master/design/70257-memory-regions.md
- typescript-go `internal/core/arena.go`, `internal/ast/ast.go` — https://github.com/microsoft/typescript-go
- esbuild `cmd/esbuild/main.go` — https://github.com/evanw/esbuild
- CockroachDB `pkg/cli/start.go` — https://github.com/cockroachdb/cockroach ; Pebble `internal/manual`, `internal/arenaskl` — https://github.com/cockroachdb/pebble
- Badger `skl/` — https://github.com/dgraph-io/badger ; Ristretto `z/` — https://github.com/dgraph-io/ristretto
- bigcache https://github.com/allegro/bigcache ; freecache https://github.com/coocood/freecache ; fastcache https://github.com/VictoriaMetrics/fastcache
- VictoriaMetrics `lib/bytesutil/internstring.go`, `lib/uint64set` — https://github.com/VictoriaMetrics/VictoriaMetrics
- Prometheus `model/labels/labels_{stringlabels,dedupelabels}.go` — https://github.com/prometheus/prometheus ;
  “30 Pull Requests Later, Prometheus Memory Use Is Cut in Half” — https://thenewstack.io/30-pull-requests-later-prometheus-memory-use-is-cut-in-half
- jsgo: `docs/gclab-2026-10-18.md`, `docs/performance-roadmap.md`, `internal/jutil/slab.go`, `internal/rhino/node.go`,
  `internal/rhino/rhino_string_pool.go`, `internal/parsing/{ir_factory.go,parser_runner.go}`, `cmd/jsgo/memstats.go`
