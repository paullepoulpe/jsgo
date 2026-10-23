# Memory round 2 (2026-10-23/24): results

Round 1 ([round1-results.md](round1-results.md)) removed slabs, pools and
cross-pass caches and cut peak heap by 34–62 % at a wall-time cost of
9–19 %. Round 2 had three goals: make the AST invisible to the GC, recover
the wall-time cost by cutting allocation *count* (not by bringing retention
back), and clear the remaining known retention. Four agents worked in
parallel, each in its own worktree, on top of `main` after it had absorbed
round 1 (`20b2f6db`); everything was merged onto
`claude/go-gc-overhead-opt-yleif9` and verified byte-identical with Java
on all 11 benchmark tasks, with `go vet`, the licence check and the full
test suite green.

## What changed

1. **Pointer-free AST** (`internal/rhino/node_store.go`, `node_tables.go`,
   `node_compact.go`, `internal/jscomp/compiler_compact.go`;
   [round2-ast-links.md](round2-ast-links.md)). `Node.typeOrColor` became
   two 32-bit handles (`tc`: JSType or Color; `aux`: JSDocInfo or ext box)
   into per-compilation tables, and the four links became 32-bit ids into
   480-node chunks behind a process-wide chunk directory. A `Node` is 56 B
   and holds no Go pointer; `*Node` stays the handle passes use. Copying
   compaction runs after `typesToColors` (after a forced GC) and whenever
   chunks allocated since the last compaction reach a quarter of the live
   tree, forwarding every cross-pass holder from the
   [census](node-holder-census.md) and rebuilding the aux table. A poison
   mode (`JSGO_AST_POISON=1`) keeps old chunks and makes stale use panic;
   the full jscomp suite and six tasks are clean under it. The type column
   is released wholesale at `typesToColors` regardless of dead typed nodes.
2. **Allocation count** (`internal/jscomp/reference_collector.go`,
   `scope.go`, `syntactic_scope_creator.go`, the parser, `jutil/slab.go`).
   Per-pass slabs (dying with their pass) for references, reference
   collections, basic blocks, scopes, scope var lists, tokens and parse
   trees; `Scope.varIndex` as a pointer-free open-addressing table. Slab
   chunks are one-page size classes (256 B–8 KB): a first version with
   10–21 KB chunks grew the page heap by 14 % on googleapis because
   multi-page spans cannot reuse scattered freed pages.
3. **Retention** (`compiler.go`, `convert_to_colors`, `typed_scope.go`,
   `global_namespace.go`, `sourcemap`). Types on detached nodes held by the
   change timeline, module metadata, module bindings and diagnostics are
   converted at `typesToColors`, so the JSType graph is released (the strict
   weak-pointer check passes); TypedScope's build-only fields moved to a
   side struct dropped after validation (120 → 88 B); ProcessDefines'
   GlobalNamespace builds only define/const candidates (checks-only: 141 k →
   10 k names, 43 → 8 MB); input source-map entries are per-line delta
   varints (bundle-mega: 49 → 12 MB kept through output); `JSGO_MEMSTATS`
   reports live-after-colors and a mapped-memory breakdown.
4. **Node tables fixed two Compiler leaks** (test processes only): the
   registry's error reporter and `CompilerInput.inlineId` kept whole
   compilers alive.

Not done, measured instead: Colors (the pool already keeps them at ~5 MB,
the 15–20 MB estimate was stale); Var/Reference byte shrinking (≤ 1.5 % of
allocated bytes each).

## Results: `main` (20b2f6db) vs final (8ac5b399)

macOS arm64, 8 cores, 24 GB, Go 1.27.1, default GOGC, 2 interleaved runs on
an idle machine, all outputs byte-identical with Java.

| Task | Peak RSS | Wall | GC CPU (final) | Objects allocated |
|---|---|---|---|---|
| typed-full | 908 → 466 MB (**−49 %**) | 10.5 → 12.7 s (+20 %) | 0.55 s | 30 → 7 M |
| typed-mega | 1,452 → 791 MB (**−46 %**) | 17.7 → 20.7 s (+17 %) | 0.82 s | 49 → 11 M |
| bundle | 510 → 292 MB (−43 %) | 4.3 → 5.2 s (+22 %) | 0.17 s | 15 → 2 M |
| typed-app | 244 → 166 MB (−32 %) | 1.5 → 1.6 s (+2 %) | 0.12 s | 5 → 1.3 M |
| bundle-mega | 1,260 → 894 MB (−29 %) | 22.8 → 31.0 s (+36 %) | 0.86 s | 69 → 10 M |
| checks-only | 403 → 304 MB (−24 %) | 2.4 → 2.6 s (+9 %) | 0.11 s | 9 → 3 M |

Per-step attribution (same pair, before the AST links step): the first
three changes together measured −7 % to +12 % wall against `main`; the
pointer-free AST added 15–25 % mutator CPU (id → pointer translation in
`GetNext`/`GetParent`/… on hundreds of millions of calls, plus no per-node
freeing between compactions) while taking GC CPU from 2–5 s to under 1 s.
A follow-up to cache the chunk base on the hot traversal paths is in
progress; see the end of this file.

googleapis (44k files, 7,483 chunks, ~35 M nodes), single runs, all 7,483
chunk outputs and stage logs byte-identical (peak RSS / wall):

| Stage | main | final |
|---|---|---|
| CHECKS | 7,368 MB / 71 s | 5,381 MB / 102 s (−27 % RSS, +43 % wall) |
| OPTIMIZATIONS | 5,625 MB / 54 s | 5,616 MB / 77 s (0 %, +43 %) |
| FINALIZATIONS | 1,886 MB / 10 s | 1,278 MB / 12 s (−32 %, +11 %) |

Whole-run GC CPU on the final build: 0.74 s (3.7 M objects allocated).
OPTIMIZATIONS did not shed RSS here, and the mutator cost is larger than on
the standard tasks; both are under investigation in the follow-up below
(compaction transient at that stage's peak, and id translation in the
traversal loops).

Since round 1's start (8ae12b3, Go 1.24.7): typed-full peak RSS 949 →
466 MB (−51 %), typed-mega 1,487 → 791 MB (−47 %), bundle-mega 1,408 →
894 MB (−37 %); wall typed-full 30.2 → 12.7 s on this machine is not
comparable (different machine), see round 1 for the Linux numbers.

## GC policy

With the AST untraced, a GC cycle costs little, so a lower GOGC turns the
saved CPU into RSS: at GOGC=50 the final build measured typed-app 140,
bundle 241, checks-only 294, typed-full 466, typed-mega 781, bundle-mega
689 MB (−9 to −25 % vs GOGC=100) with GC CPU 0.2–1.5 s and mutator CPU
unchanged. The default is left at 100; switching the CLI default to 50
(env `GOGC` still wins) is recommended and is a one-line change.

## Lessons

- A node with no pointers in a chunk with one pointer (its generation
  owner) is traced in O(1), but Go's `/gc/scan/heap` accounts the whole
  chunk, so the memgate "scan" column no longer measures AST tracing work;
  GC CPU does.
- Compaction needs a forced GC first: the copy otherwise sits on top of
  the not-yet-collected type graph (typed-full peak 731 vs 535 MB).
- Chunk allocations must stay in one-page size classes; the page heap,
  not live data, grew when they did not.
- Write chunk-directory slots as integers: a pointer store's write
  barrier shades the old value, which may be a chunk the GC already freed.

## Next

1. Shave the id-translation cost on the hot accessors (in progress).
2. GOGC default 50.
3. Chunk-by-chunk release during the compaction copy (removes the 2× AST
   transient at `--gogc 10`).
4. Dense node-id-indexed side tables (Timeline, refsIndex, scriptSizes).
5. Remaining allocation sites: AST nodes themselves (IRFactory), identifier
   strings in non-ASCII files, qualified-name strings, TypedVars.
6. typed-full/checks-only now peak in TypeInference (TypedVars, typed
   scopes, flow states): the next memory lever.
