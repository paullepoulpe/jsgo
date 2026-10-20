# Pass-owned reference batches

`JSGO_EXPERIMENT_REFERENCE_BATCHES=1` enables a private record allocator in the
stock `InlineVariables` and `VariableReferenceCheck` passes. Public
`NewReference` and `NewReferenceCollector` retain their existing behavior.
Custom compilers, compiler observers, and collector behaviors fall back.

Both stock collectors retain all their `ReferenceCollection` objects and
references in their variable map through traversal. Neither stock behavior
removes map entries or exposes that mutable map. The allocator retains only its
current GC-tracked `[]Reference`: capacities are 8, 16, 32, 64, then 128. It never
recycles or clears a published record and has no global owner or slab list.
Every published pointer remains unique and mutable as before. Keeping one
record alive also keeps at most 127 neighboring records alive; this is why
arbitrary public collectors do not opt in. Unused current tail storage is at
most 127 * 32 = 4064 bytes per collector on the benchmark's 64-bit platform.

Existing parser-enabled allocation profiles attribute these sampled estimates
to `newReferenceWithScopeAndInputId`:

| Workload | Objects | Bytes | Share of allocated objects | Share of allocated bytes |
| --- | ---: | ---: | ---: | ---: |
| bundle | 4,800,527 | 153,621,072 | 5.59% | 3.11% |
| typed-app | 901,134 | 28,836,720 | 4.89% | 2.40% |
| checks-only | 262,147 | 8,388,864 | 1.37% | 0.56% |

These constructor totals include other collectors and cloned references, so
they are an upper bound on eligible allocations, not measured savings. Batching
primarily reduces allocation count; payload bytes and pointer scanning remain.
Reduced allocator/GC overhead must be checked against possible extra retention.
The corresponding flow-scope/overlay/slot constructor totals were substantially
smaller and absent from bundle, so a flow allocator facade was deferred.

Focused tests check published identity/fields through repeated slab replacement,
bounded sizes, public/custom fallback, and full pass AST/output/diagnostic/change
ordering against individual allocations with generic and native reference maps.
No performance improvement is established by source inspection.
