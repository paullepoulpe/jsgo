# Private inlining clone allocator feasibility

2026-10-19. Existing default-GC structural heap512 profiles were read after
forced GC; no new workload or profile was captured. All CloneTree/CloneNode
ancestry is inclusive. CPU64 results can lose outer ancestry; these captures
are particularly sparse for cloning, so zero samples are not proof of zero CPU.

| Task | Sampled clone MiB | Sampled clone objects | Whole allocated bytes | Inlining clone MiB | Observed clone CPU |
|---|---:|---:|---:|---:|---:|
| bundle | 24.50 | 321,146 | 0.51% | 22.50 | 0.01/21.18s |
| typed-app | 7.50 | 108,140 | 0.60% | 6.00 | No samples |
| checks-only | 0.50 | 6,554 | 0.03% | 0 | No samples |
| downlevel-whitespace | 0 | 0 | 0% | 0 | No samples |
| downlevel-simple | 2.00 | 36,046 | 0.18% | 0.50 | No samples |

Bundle cloning accounts for 0.42% of sampled allocation objects. Its clone
allocation leaf is Node construction, not a hidden large metadata-copy tree.
Typed clone allocation includes 0.5 MiB of Node extension records as well as
Node payloads. These are allocation estimates, not exact counters or peak RSS.
The old canceled worker's extern clone hotspot is a different workload and
must not be carried into ordinary CLI conclusions.

## Existing source support

Rhino already has a typed, owner-local NodeAllocator.CloneNode. It preserves
the shallow, detached, false-type-expression-clone semantics of CloneNode:
folded Number, BigInt, String and TemplateLiteral subclass payloads, type/color,
source location, original name and shared immutable property list. Public
CloneTree recursively creates children, reconstructs parent/next/previous
links, and separately clones a closure-unaware shadow. A private
CloneTreeWithAllocator can reuse that same recursion with allocator.CloneNode
and preserve the public default API. It would need explicit shadow recursion
and original child encounter order, not just a plain child walk.

Potential stock inlining sites are FunctionInjector's safe return-node clone,
FunctionToBlockMutator's function/argument/return clones, FunctionArgumentInjector
replacement templates and rest arguments, InlineFunctions' safe function copy,
and InlineVariables' repeated value clones. Public injector/mutator constructors,
custom observers/conventions/callbacks, and false/true type-expression variants
need separate eligibility. The compiler must never install an allocator in a
process-global Node constructor or reuse a published node.

## Lifetime and memory tradeoff

Current NodeAllocator grows 8,16,32,64,128 and retains only its current block.
Nodes do not retain their allocator. Discarding an owner releases unused blocks,
but any published interior node keeps its whole backing block alive.

An allocator per inlining invocation avoids cross-call lifetime coupling, but
has at least eight Node slots. Tiny return/name clones can therefore allocate
many more payload bytes than individual nodes. An allocator per InlineFunctions
Process amortizes tail waste, but can mix long-lived safe function templates or
selected output nodes with deleted branches, wrappers and temporary clones.
One survivor can retain up to 127 neighbors. Those neighbors can retain their
own pointers, descendants and metadata. Peak RSS can increase despite fewer
allocation objects. Ending a pass does not release a chunk still referenced by
published AST nodes. Separate owner allocators would remain concurrency-safe;
a shared compiler allocator would need an explicit parallel phase contract.

Exact-size allocation per cloned tree could avoid unused tails, but requires
an extra node-count walk (including hidden shadows), or a trustworthy frozen
summary. Partial subtree deletion still co-retains neighbors in the same tree.
That extra traversal needs an elapsed-time justification rather than an
objects-only microbenchmark.

## Decision

No prototype was implemented. The existing ordinary-CLI profiles place this
below 0.6% of allocated bytes and below meaningful direct CPU resolution. The
allocator would not remove the node payload bytes and could increase retention.
It is technically feasible but weaker than the measured contextual-inverter
representation opportunity. Revisit only if a newer ordinary-CLI workload
shows materially larger clone ancestry or an independently justified clone
elimination strategy.

Raw counts and allocation leaf/owning-pass breakdowns:
`clone-allocation-audit.json`.
