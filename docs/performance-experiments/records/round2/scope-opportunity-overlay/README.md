# Untimed scope reconstruction opportunity diagnostic

Separate source overlay; not a production cache or performance candidate.
Generated from `/workspace/jsgo-round-compound-next` (manifest records its exact
commit and each input hash). No source worktree file is edited. No Go build or
counter run has been performed yet.

Prepare, when needed:

```sh
python3 /workspace/jsgo-results/round2/scope-opportunity-overlay/prepare.py --root /workspace/jsgo-round-compound-next
```

The generated `overlay.json` is independently usable with Go `-overlay` for an
ordinary CLI build. Set `JSGO_SCOPE_OPPORTUNITY_JSONL=/absolute/path/output.jsonl`
only for an **untimed** diagnostic run. One JSON record is appended per
compiler-owned PhaseOptimizer invocation, including panic/early-return paths.
Preserve ordinary CLI flags and tracer mode. No TIMING_ONLY option, extra AST
traversal, compiler callback, change report, scope reuse or declaration suppression
is added. Existing byte/diagnostic gates are still required before collecting data.

An optional `--base-overlay /path/to/overlay.json` composes another overlay
without modifying it. For example, the next-profile-overlay replacements can be
included if root later explicitly wants a combined untimed capture. The default
does not enable CPU/heap profiling or change CLI main.

## Observations

The real Scope/Var constructors and scanner execute fully. Observations attach
to the existing fresh scope construction and original DeclareVar cut points,
after its NAME/IMPORT_STAR invariant check. Exact stock GetInput is an additional
read-only map lookup. Implicit-goog-namespace scans are conservatively excluded
from opportunity prediction. All successful SyntacticScopeCreator constructions
inside compiler-owned PhaseOptimizer NamedPasses are counted by pass/token and
live declaration count; scopes outside those boundaries and facade compilers
are explicitly outside the denominator.

Only InlineVariables-owned creators (including inlineConstants and
earlyInlineVariables) and serial FSIV's outer/all-vars/escaped creator participate
in shadow reuse. FSIV's separate gather traversal remains unselected. Public
constructors return fresh objects; custom compiler/redeclaration/convention/error
observers, custom pass options, unsupported pass names/stages, debug/tracing and
other unsupported contexts are not admitted by the diagnostic's own selection
test. Existing production guards and behavior remain unchanged. Parallel FSIV
facades are not observed as candidates. Existing bounded FSIV-window hits that
already avoid construction are not counted as new opportunities.

Two independent shadow policies assign generation numbers to actual fresh
scopes. A canonical parent must have been observed in the current pass:

* `function_epoch`: exact global declaration sequence, followed by exact parent
  generation and unchanged local enclosing FUNCTION/SCRIPT identity/time.
* `exact_declarations`: exact global and local declaration sequence regardless
  of body-only epoch changes.

Both compare the actual fresh scan's ordered `(node,name,input)` events and
validate the originally retained Scope/Var snapshot. No hash serves as equality.
When a root/global proof fails, its new shadow generation causes all descendants
to miss even when their own epochs are unchanged. Real returned scopes and Vars
always remain newly allocated. `same_pass_opportunity` and
`cross_pass_opportunity` are separate; `from_pass/...` records transitions.
These are reconstruction opportunities, not implemented cache hits.

Miss/rejection buckets distinguish first construction, unselected parent,
root context, changed parent generation, direct Scope/Var mutation, declaration
changes with/without an epoch change, epoch-only changes, module/provide and
ineligible contexts. Per-pass snapshots include retained entry counts, stale
epochs/parents/deleted entries, unique reachable Scope-parent graph weight,
actual diagnostic slice capacities and explicit snapshot/entry payloads.
The reachable graph is the diagnostic's actual retained graph, an upper bound
for a future cache: real construction still creates fresh parents even where
the shadow policy predicts a canonical-parent hit, and retained fresh children
can keep those otherwise avoidable physical parents alive.
Old parent generations are intentionally visible in retained-graph accounting;
only one latest entry per AST root is stored. The counters expose retention cost
instead of silently assuming stale generations are free.

## Byte definitions and lifetime

Counts are exact observations of successful original constructions, live explicit
Vars and scanner events. Scope and live-Var payload bytes use `unsafe.Sizeof`.
The ordered-map byte field is a **source-derived model**, including header,
the first nine individual entries, large-map metadata and 16/32/64/128-entry
slabs. It excludes opaque runtime hash-index buckets, allocator size-class
rounding, superseded implicit-namespace declarations and transient collections.
There is no RSS, malloc-byte or elapsed-time claim. JSON includes current layout
sizes. Diagnostic snapshots themselves allocate/retain substantial memory.

Counter state is confined to overlay-only PhaseOptimizer, creator and scanner
fields. Scope, Var and Compiler have no added fields or lifecycle mutation hooks.
There is no global registry. Physical Scope-to-generation maps are reset at every
NamedPass boundary; retained shadow entries/snapshots are released in the
PhaseOptimizer Process defer on success, early return and panic. A mutex protects
record updates from independent constructors; stock pass boundaries already join
their child work before changing active-pass state. Unsupported custom recursive
optimizer schedules are not predictive diagnostic targets.

Do not use this build's time, heap/RSS or allocation totals to report production
performance. First inspect weighted cross-pass opportunities, global declaration
churn, parent cascades, empty-scope anchors, and retained/snapshot costs. Only then
decide whether a separately reviewed bounded cache implementation is worthwhile.

Overlay-only tests exercise actual hooked constructors, fresh identity,
cross-pass counts, epoch versus exact-proof divergence, global invalidation
cascades, direct map mutation, public Var reinitialization and panic cleanup.
They have been written and formatted but **not compiled or executed**.

To aggregate a future JSONL capture without double-counting overlapping views:

```sh
python3 /workspace/jsgo-results/round2/scope-opportunity-overlay/summarize.py /absolute/path/output.jsonl
```
