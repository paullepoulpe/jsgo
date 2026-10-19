# Bounded parser front cache of global string cells

Source-only audit, 2026-10-19. No build, test, profile capture/decoding or timing
was run. Current source inspected: `/workspace/jsgo-structural-bitset-template`
at `63dc7116`. Existing attribution only:
`frontend-allocation-audit.json`, `node-string-cell-layout-audit.md`, and
`structural-profiles/environment.json`. No production implementation is proposed
as ready or measured.

## Feasibility and actual routes

The earlier warning about a separate parser-local intern pool does not reject
this different design. A front cache can store the **exact `*string` returned by
the original global `rhinoStringPoolCell`**. A local hit publishes that same
cell; a miss calls the unchanged global pool. It creates no per-spelling string
clone or new escaped string-header cell, and keeps cross-file sharing on global
misses/hits. Node layout and public SetString need not change.

`NodeAllocator.NewString` currently creates a StringNode and calls SetString
(`internal/rhino/node_allocator.go:36`). SetString checks the fresh node's kind,
then calls the mutex-protected global pool. Its kind is already set immediately
before this call, so a private cache path can preserve that check before
assigning the resulting cell. Pool lookup has no callback or custom equality.

The allocator is constructed once per native speculative parse, when the owned
allocator option is enabled (`internal/parsing/native_parser_experiment.go:196`).
It is owned by the parse's IRFactory. Nodes do not retain it. Separate input
parsers get distinct allocators and can run concurrently; cache lookup itself
needs no synchronization. Discarding a rejected attempt releases its allocator
and cache. No global cache/list or AST-to-allocator backlink is needed.

Both IRFactory.NewStringNode and NewStringNodeWithToken use this allocator.
Their shared factory also transforms original-grammar fragments, so accepted
native attempts cover strings inside those fragments too. It does **not** cover
whole-file fallback: that discards the attempt and uses the ordinary factory.
JSDoc parser constructors that directly call public Node_NewString also remain
outside the cache. Existing StringNode cloning already copies the shared cell
without looking it up, so it offers no lookup to eliminate. Template raw/cooked
strings take separate methods; keep them unchanged in an initial bounded
candidate rather than assuming every interning call routes through NewString.

## Value, reset and ownership semantics

Global cells are immutable. SetString replaces a node's cell pointer, never the
cell contents. Current-name getters return string values; original-name mutable
headers have separate storage and must remain unchanged. Equality is by Java
string value, not cell address. A local comparison must therefore compare the
spelling (including arbitrary valid Unicode/WTF8 bytes), not input-header or
backing-pointer identity. Ordinary Go string equality is the same equality used
as the global map's key.

The global map clears at 1,048,576 names. A front cache can still retain an old
cell safely across such a reset, just as existing live nodes and clones do.
Values remain correct; no stale mutable state exists. This does **not** promise
the same generation's canonical pointer or identical map occupancy/reset timing:
a local hit can avoid reinserting its spelling after a clear. Pointer identity
and exact pool occupancy are not public compiler contracts. If exact current
generation pointer equality were required as an internal invariant, a cache
would need generation checks; no such invariant was found. Adding an epoch only
for this unobservable behavior would add synchronization and new global state.

A hit can share an old-generation cell while another input obtains a newer cell
of the same spelling. This is existing legal pool behavior, not newly allocated
per-parser interning. However, a bounded cache may keep an otherwise dead old
cell alive until its parse finishes. Bound both entry count and admitted string
length if a hard retained-byte bound is desired.

On a map-cache miss, **store `*globalCell` as the key**, not the incoming scanner
slice. The global pool cloned that spelling precisely to avoid retaining a
larger source buffer. Saving the original input slice as a cache key would undo
that ownership protection. A fixed pointer-only cache avoids key headers and
source retention entirely. Caches may hold only string cells, never tokens,
nodes, source files or scopes.

## Candidate shapes and costs

The cache removes no normal per-node allocation: cells are already globally
shared. Its target is repeated lock acquisition plus global lookup, particularly
parallel parser contention. Every miss adds a front lookup before the original
lookup. A local hash map still hashes strings and introduces its own allocation
and scan cost; it is not a free replacement for the existing fast native map.

Two plausible bounded variants, to select by an untimed hit inventory first:

* A two/four-cell move-to-front array stores only original `*string` pointers.
  Hits use at most two/four value comparisons, then reorder pointers. Misses call
  the global pool and replace the oldest entry. Two pointers add 16 bytes of
  payload; explicit activation can grow the 32-byte allocator to the 64-byte
  allocation class. This avoids per-cache map allocation and extra string-key
  headers, but limited recency may give too few hits.
* A lazily allocated, capped `map[string]*string` of 64 short spellings gives
  wider locality. Stop admitting or use bounded eviction; never grow forever.
  The map's copied key headers are additional index storage, though it creates
  no separate canonical header cells. For Go's eight-slot Swiss groups, 64-entry
  reservation rounds to 128 slots of roughly 24 bytes plus control bytes: about
  3.2 KiB of group payload, plus map/table/directory and allocator overhead.
  Reservation128 is about 6.4 KiB of group payload. These are derived layout
  estimates, not measured allocations. Per-active-parser retention scales with
  compiler thread count. Admit only short spellings to avoid retaining huge
  otherwise-dead literals after a pool clear.

A two-cell MRU is the simpler first implementation **if its real hit rate is
useful**. No new node layouts, source-local canonical cells, persistent/global
state, observer callbacks or compiler mutation hooks are required. A private
parser constructor can opt in only the native owned allocator; zero-value/public
allocator construction and public SetString can keep the original path. Reuse
existing native parser admission/fallback rather than changing custom reporter
or source-file behavior. Do not add several production knobs to compare cache
policies; a temporary isolated binary/prototype is sufficient.

## Existing attribution and honest ceiling

Existing default-GC structural CPU captures have the native parser and owned
allocator enabled, with affinity to two CPUs. Their frontend table reports:

| Fixture | Whole sampled CPU | Retained frontend ancestry | Native-direct context |
| --- | ---: | ---: | ---: |
| bundle | 21.18 s | 0.81 s / 3.82% | 0.46 s / 2.17% |
| typed-app | 5.93 s | 0.26 s / 4.38% | 0.23 s / 3.88% |
| checks-only | 8.57 s | 1.05 s / 12.25% | 0.88 s / 10.27% |
| downlevel-whitespace | 1.39 s | 0.23 s / 16.55% | 0.10 s / 7.19% |

The cache can eliminate only a subset of these parser costs. Scanner work,
UTF16 conversions, positions, AST construction and JSDoc remain. The
rhinoStringPoolCell leaf itself is absent from the frontend top-flat tables
(all sampled frontend leaves for typed/whitespace, top50 for bundle/checks).
Native map lookup contributes 0.02 s frontend flat in bundle and 0.03 s in
checks, but these sites serve other frontend maps too and cannot all be assigned
to the string pool. Mutex wait time and global-cell ancestry are not separately
reported. The 64-frame CPU stack limit can lose frontend ancestry, so these
figures are **observed envelopes, not rigorous whole-program upper bounds**.
They do not establish a multi-percent string-pool bottleneck or parallel
contention.

Existing frontend sampled pool allocations are only 3.75 MB bundle, 1.61 MB
typed, 5.85 MB checks, and 2.15 MB whitespace. These include unique cells/map
growth, which the front cache normally cannot remove. Cache storage therefore
initially adds allocation/retention. No memory improvement should be promised.

## Recommendation and next evidence

Semantically viable and distinct from a parser-local intern pool, but currently
**low-priority/unjustified as a major speed redesign**. Do not implement a wide
cache on the strength of repeated-name intuition. First collect untimed actual
NewString call counts and simulated two/four/64-entry hit rates, separately for
accepted attempts and rejected/fallback attempts, with short-string exclusions
and per-file high-water cache size. Aggregate per-parser owned counters after
the attempt, not global per-call locks or retained node maps. Do not interpret
sampled allocation objects as exact unique spellings.

If hit rates are high, check isolated parallel-parse contention at the same
thread counts and default GC, then compare complete fresh CLI time, TotalAlloc,
Mallocs and peak RSS for one bounded policy against the unchanged baseline.
Required differential cases include cross-allocator sharing on misses,
SetString replacement isolation, equal spellings from different backing
strings, empty/Unicode/WTF8 strings, cache eviction, speculative fallback,
concurrent owners, and cache hits across a forced global-pool clear. Existing
allocator clone/discard/ownership tests remain relevant. No implementation,
runtime inventory, performance win or final threshold has been asserted.
