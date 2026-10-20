# Unchanged-function FSIV analysis reuse: source audit

This candidate stops at source assessment. It adds no result cache, runtime flag,
or claimed speed or memory improvement. The proposed optimization was to skip
FlowSensitiveInlineVariables (FSIV) analysis when a function appeared unchanged
since the previous invocation.

## Scheduling and the actual opportunity

The stock DefaultPassConfig does not mark FSIV's PassFactory as loopable. It
requests the pass before the main optimization loop and again in finalizations;
each request creates a new pass object. A pass-local cache would therefore not
provide reuse across the stock invocations. Reuse across them would require
compiler-owned state, with explicit initialization and restore boundaries.

One FSIV invocation is deliberately not a fixed-point computation. After an
inline introduces dependencies absent from another candidate's reaching-def
metadata, EnterScope discards that other candidate rather than rebuilding the
CFG. The source comment explicitly describes this conservative heuristic.
Recording the *resulting* function stamp as already optimized would incorrectly
assume that another invocation cannot do more work on that same resulting AST.

A possible certificate must instead describe a function on which the previous
invocation made no relevant rewrite at all. Rewrites in descendants matter too:
escape discovery traverses nested functions, and parent expressions may contain
those functions. Ancestor rewrites can change lexical binding or ancestry even
when the function's own stamp remains unchanged. Default-parameter nested
functions are processed before the enclosing body, so a certificate created
there cannot ignore changes made later to the enclosing function.

## Why a function stamp is not the dependency key

The necessary inputs extend beyond the function's own syntax:

- Own and descendant syntax, reference ownership, declaration names, and
  parent/scope identity determine the CFG, escape set and canonical bindings.
- Relevant enclosing declarations and script membership affect those bindings.
  ROOT is deliberately never assigned an AST change stamp.
- Call side-effect flags, colors or raw types, accessor facts and AstAnalyzer
  options affect whether moving an expression is safe.
- Compiler lifetime, restored ASTs and change-stamp reset are separate epochs.

PureFunctionIdentifier does report changes when it updates call side-effect
flags. That handles an important dependency; it does not make all semantic
metadata stamp-tracked. ChangeVerifier calls
Node.IsEquivalentWithSideEffectsToShallow, whose type and JSDoc comparison modes
are IGNORE. Thus passing ChangeVerifier is not a proof that every fact read by
AstAnalyzer is unchanged. This audit does not establish that a particular stock
pass sequence violates such a cache; it establishes that the existing stamp
contract alone is insufficient to justify it.

AccessorSummary also exposes its map through GetAccessors. Summary pointer
identity alone would not detect an in-place content change. Raw-type operation,
custom compilers/conventions/observers, and mutable user options would need a
proven contract or a conservative fallback. Color-based operation still needs
the per-node color dependency accounted for.

Compiler.restoreCompilerState resets ChangeTracker's counter after loading an
AST whose node stamps are zero. Comparing only the current numeric counter
cannot identify every reset: subsequent changes can reach an old value again.
Compiler initialization/restore must explicitly retire compiler-owned caches,
or supply an epoch included in the key. Detached functions must also be pruned
so a cache does not retain removed ASTs for the rest of the compilation.

## What would justify a future prototype

A future implementation should first measure how many costly analyses satisfy
all of the following, without reporting a performance gain from fallback work:

1. The first invocation actually reaches CFG/reaching-definition analysis.
2. Its entire relevant subtree and enclosing dependency set remain unchanged
   during and after that invocation.
3. Every analyzer input has a stable, audited generation or validated value.
4. The next invocation has the same inputs, compilation epoch and canonical
   binding context.

Only then should a clean certificate skip the CFG/analyses. Tests must compare
against recomputation after own, sibling, ancestor, descendant, accessor,
side-effect, color, source-membership and restore changes; dependency-chain
cases must exercise FSIV's deliberate candidate backoff. Actual hit/miss and
saved-analysis counts, ordinary CLI output/diagnostics, and retained memory
must accompany timings. Parallel jobs would need serialized certificate
publication or private ownership; sharing a mutable compiler cache implicitly
through their facade would invalidate the ownership proof.

There are only two stock opportunities, separated by many transformations.
Deep metadata validation and retained-node bookkeeping can consume the walks
and memory the cache aims to save. The bounded alternatives remain stronger
first steps: skip unused backward analysis when no candidate exists, fuse facts
within one unchanged analysis interval, or replace generic graph/flow machinery
while retaining the original recomputation boundary.
