# Shared Java-emulation overhead audit

Read-only audit on 2026-10-19. No new workload captures, builds or timings.
Existing input: `next-balanced-pgo/{task}.merged.pprof`, originating from source
`eecbcf44`, binary SHA `acefcb56f439a9288cf2827a34ca96a9f36c74acf5cc36eb7c0c6d80a9bb7c30`.
These are diagnostic GOGC=200 captures, not current default-GC performance
numbers. They predate the latest typed declaration, RUC and inlining prototypes.

| Task | Captured CPU s | GetNext flat | TraverseBranch flat | String-key LinkedHashMap.find retained union | ObjectHashCode retained union |
| --- | ---: | ---: | ---: | ---: | ---: |
| typed-app | 11.92 | 8.47% | 7.05% | 3.78% | 0.92% |
| bundle | 19.33 | 8.38% | 7.40% | 5.33% | 1.66% |
| downlevel-whitespace | 11.12 | 4.32% | 3.69% | 3.15% | 1.17% |
| downlevel-simple | 12.12 | 6.77% | 7.43% | 5.94% | 1.49% |
| checks-only | 7.93 | 8.83% | 7.31% | 2.90% | 1.51% |

Flat samples identify actual innermost symbols, including inline attribution.
Retained unions remain lower bounds because CPU stacks stop at64 frames; they
are not hard optimization ceilings. Columns overlap through callers and are not
additive. Generic shapes conflate different concrete maps and passes. The
reported map union is one string/pointer generic shape, not every map cost.

`Node.GetNext` is literally `return n.next`. More callback specialization does
not remove the underlying dependent pointer loads. Eliminating whole repeated
walks or replaying a proven shape-stable itinerary is the larger design lever.
The current inliner discovery tape is already a test of that idea; copying
another callback driver without evidence would repeat a track that did not
produce broad wins. A reusable checking itinerary would need exact scope and
script events, per-pass ShouldTraverse behavior, no AST-shape mutations and
strict mutation/observer ownership. Current TypeCheck annotates types consumed
by later checks, so naive per-node pass fusion would change semantics.

String-key ordered-map lookup is the broad remaining emulation family. Node.str
already points at an interned immutable string header. However,
RhinoStringPoolCell clears the global pool at1M names; equal spellings can then
have different cell identities inside a live AST. Raw pointer-key substitution
is invalid. A compiler-owned atom/ID representation could avoid repeated string
hashes, but would need value fallback or a new canonicalization lifetime and
change many node-to-scope lookup boundaries. A pointer-fast-lane plus value map
adds another map, memory and miss work. Existing native scope/declaration/RUC
storage addresses part of this old profile already. Short strings also use a
very fast native hash path, while integer maps retain their own hash/probe cost.
There is not yet evidence of a new shared >5% win from that rewrite.

Generic ObjectHashCode/ObjectsEquals are smaller than the broad-map/traversal
families here. Their public contracts include custom equality, typed nils,
float bit equality and nullable strings. Removing those contracts globally for
an observed1–2% retained family would be an unjustified semantic tradeoff.

The audit therefore recommends fresh attribution of the current candidate before
another shared collection/AST rewrite. No new implementation is proposed from
these old profiles alone. Native parser acceptance on freshly parsed stock
externs is a coverage question sent to the parser owner; it is unrelated to
preparse caching and has not been asserted as an unimplemented gain.

Parser-owner follow-up: no existing report separately inventories native
acceptance/original-fragment fraction for the freshly parsed builtin extern
corpus. Prior coverage covered real bundles, core-js and Closure test inputs;
there is no explicit extern exclusion in the current sink. This is an evidence
gap, not proof of a remaining fallback or speed opportunity. No inventory job
was launched during the timing reservation.
