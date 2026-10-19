# Source audit: scope storage and scope reuse

Source-only review, 2026-10-19. No performance claims are inferred from the
allocation models below. Existing linked scope storage remains the recommended
baseline. Candidate branches are separate from the validated CLI compound.

## Dense declaration map

Candidate: `/workspace/jsgo-round-scope-map`, commits `d4f93288`, `6ec207a1`.
Standalone correctness and race tests passed. Its borrowed immutable name cell
is not safe for arbitrary public `InitAbstractVarBase` reinitialization: changing
the Var's name changes the tiny map's key. A raw declaration can put one Var in
multiple scopes, so detaching only from the Var's original scope is insufficient.
Scope reinitialization also leaves old live map views potentially reachable.
Ownership restrictions plus cold constructor-detachment hooks could close these
paths without extra common-path allocations, but add lifecycle complexity.
Do not adopt those hooks for an unproven storage improvement.

An owned inline string key closes this correctness gap and requires no Var or
Scope field changes. The entry grows from 24 to 32 bytes on amd64. The current
48-byte map header and geometric vector growth then have these whole-map costs,
including Go small-object size classes:

| Live declarations | Original linked retained = cumulative | Owned dense retained | Owned dense cumulative |
| --- | ---: | ---: | ---: |
| 1 | 80 | 80 | 80 |
| 2 | 128 | 112 | 144 |
| 3 | 176 | 176 | 272 |
| 4 | 224 | 176 | 272 |
| 5 | 272 | 304 | 528 |
| 6 | 320 | 304 | 528 |
| 7 | 368 | 304 | 528 |
| 8 | 416 | 304 | 528 |

Original tiny linked maps allocate a 32-byte header and individually rounded
48-byte entries (40-byte payload). Owned dense maps allocate a 48-byte header
and backing arrays with capacity 1, 2, 4, 8, each entry 32 bytes. Empty scopes
retain a nil map in both versions. There is no cumulative byte-allocation win
at sizes 2 through 8; retained storage improves at some sizes but regresses at
five entries. At eight entries, object allocations fall from nine to five and
retained bytes fall 27%, while cumulative bytes rise 27%. Larger linked maps
already use an index and slab allocation, so this table must not be extrapolated.
Actual per-scope size/capacity distributions were not collected; generic linked
map histograms do not establish the distribution of scope declaration maps.

Recommendation: retain linked storage. If a dense version is macro-tested, use
owned string keys, and measure total allocation, object count, peak RSS and time.
Its plausible benefit is retained memory/locality, not an established allocation
reduction.

## Iterator costs

The dense iterator uses an atomic active count and deferred cleanup so mutation
cannot compact/reuse integer positions while a live iterator holds them. A
snapshot iterator would change observable append/remove/clear behavior. Removing
tracking while retaining tombstones indefinitely can cause unbounded storage
growth under repeated remove/reinsert; swapping generations needs another stable
storage object and does not solve removal compaction by itself.

Go 1.27.2 compiler diagnostics confirm the deferred iterator body cannot inline,
but its yield argument does not escape. Separate diagnostic commit `ceceacbb`
compares equivalent declaration-index sum loops at 0, 1, 4, 8 and 32 declarations:

* Direct full loops: dense 0 allocations, linked 0 allocations.
* Direct early-exit loops: dense 0 allocations, linked 0 allocations.
* Opaque retained `iter.Seq` consumer: dense 2 allocations, linked 2 allocations.

An opaque consumer must allow an arbitrary iterator to retain the yield closure;
the backend cannot remove these allocations without restoring a concrete known
call target at the consumer or changing its iteration API. Existing concrete
loops already allocate zero. Atomic/defer CPU cost remains unmeasured. Dense
generic `Get` additionally has inline cost 98 (budget 80), versus 62 for linked
`Get`; its concrete wrapper still calls the generic body. Do not rewrite iterator
semantics to chase an allocation regression that was not observed.

## Compiler-wide scope reuse

Candidate: `/workspace/jsgo-round-scope-reuse`, commit `56275ea4`.
Do not integrate in its current form.

The owner pointer grows Scope from 48 to 56 bytes on amd64, moving its allocation
class from 48 to 64 bytes: an unconditional 16-byte, 33% increase for every Scope,
including public creators, unrelated passes and runs with the experiment off.
The candidate also adds mutation hooks, compiler ownership, a cache map, a
64-byte cache-entry object, and 40-byte declaration-event snapshots with geometric
backing-array growth. Snapshots record scanner events, which can outnumber live
declarations. Validation replays the original scanner, preserving AST-walk cost.

Cache keys require exact parent Scope identity. Module scopes always stay fresh;
provide/legacy-module scans reject caching. Stock-compiler, default-handler and
non-debug guards further restrict eligibility; parallel FSIV cohort jobs use
uncached creators. Global stamp changes retire the entire cache. Every reported
AST rewrite increments that stamp. The fixed-point optimizer also increments it
before each pass and in change queries, preventing reuse across its pass runs
even when a particular scope is unchanged.

InlineVariables uses one ReferenceCollector traversal. NodeTraversal already
memoizes materialized scopes in its own scope stack, and the hoisted-function
set avoids re-traversing jumped-to function bodies. There is no obvious material
same-pass hit opportunity in that consumer. Its rewrites repeatedly invalidate
the cache for later scope requests.

Serial FSIV has a concrete same-function opportunity: its outer function-block
scope exists before `NodeUtil_GetAllVarsDeclaredInFunction`, and that helper and
`DataFlowAnalysis_ComputeEscaped` synchronously re-traverse the same function
before candidate rewriting begins. Function-body and lexical-child scopes can
hit with identical parents in this window. Subsequent rewriting retires the
global cache before later scope creation. Reusing scopes in this read-only window
is a narrower opportunity than compiler-wide or cross-pass reuse.

Invalidation is lazy at the next cache lookup. After the last eligible pass,
subsequent unrelated rewrites do not themselves clear retained entries. Scope,
Var and declaration snapshots can therefore stay live until compilation ends.
Total hit counts alone cannot establish a useful memory or speed saving.

Required untimed evidence before reconsideration:

* Requests/hits and saved declaration counts by consumer and phase: InlineVariables,
  FSIV all-vars collection, escaped-variable analysis, outer traversal and
  cross-pass requests.
* Miss/retirement reasons: stamp, exact-parent mismatch, module/provide exclusion,
  context or declaration validation mismatch, and direct collection mutation.
* Scanner replay visits, declaration events versus live Vars, cache/snapshot
  peak capacities, and entry lifetime after the final eligible pass.
* Rewrite counts and stamp changes per analyzed function, plus actual allocation
  and peak-RSS comparison against unchanged scope storage.

## Less invasive direction

If profiling confirms enough FSIV scope-allocation cost, investigate a private
per-function creator confined to the synchronous all-vars/escaped-analysis window.
Seed the already-created function-body scope, record lexical scopes during the
first helper, reuse those exact-parent scopes during the second, then discard the
creator before gathering/rewrite work. Misses during the second traversal (such
as nested functions) should delegate without creating persistent cache entries.
This could avoid compiler ownership, global invalidation, declaration snapshots,
validation replay and the owner field on every Scope. It needs an explicit audit
that the two helper traversals and default callbacks do not mutate declarations;
custom/error-observing implementations must retain the original path. This is a
source proposal, not implemented or measured. Parser work and compiler inlining
experiments currently have higher priority than another scope cache.
