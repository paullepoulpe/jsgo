# Intrusive ready FIFO for private RemoveUnusedCode tasks

This source-only follow-up to `d4b1aeaa` selects a second native task mode with
`JSGO_EXPERIMENT_REMOVE_UNUSED_INTRUSIVE_FIFO=1`. It is active only inside an
already-eligible `JSGO_EXPERIMENT_REMOVE_UNUSED_TASKS=1` invocation. The copied
payload FIFO remains compiled in the same binary for differential comparisons;
public/custom paths retain the original continuations.

The first arena copies a list's node/scope payloads into a separate growing FIFO,
zeroes/frees list rows, then copies and clears FIFO slots when applying work.
The intrusive version instead splices the list's existing head/tail row IDs into
a ready chain. Ordinary ApplyList is O(1): no row walk, payload copy, allocation,
or pointer publication is required. A ready dequeue copies the payload to the
stack and clears/releases its row before traversal can enqueue more work.
List IDs, clear-in-place builder aliasing, separate removable applied bits, and
original breadth-first order are unchanged. Empty list application still clears
that same list and cannot leave a dangling queue tail.

All payloads remain actual GC-visible Node/Scope pointers. Only task-row links
are integer IDs. There is no global owner, cross-phase pool or integer AST
handle. A Process defer releases the arena on success or failure. No pointer to
a consumed row escapes its dequeue, so same-row reuse during the next Apply
cannot overwrite the stack-owned payload being applied.

## Allocation and default-GC tradeoff

The intrusive mode removes the ready payload array and its reallocations/copies;
the original ready ring stored 16 bytes per capacity slot. Ready tasks instead
occupy an existing 24-byte arena row until dequeued. Task storage is consequently
bounded by peak simultaneous pending-plus-ready work, rather than the first
arena's peak pending rows plus its separately rounded peak ready ring.

Ignoring rounding, the first arena uses approximately `24*max(P)+16*max(R)` task
bytes; the intrusive arena uses `24*max(P+R)`. If pending rows transfer together
into readiness and the old row high-water mark was retained, the intrusive
version can remove the additional ready array. If pending and ready peaks occur
together, it may instead cost eight more bytes per ready task. Neither direction
is a measured whole-compile RSS or elapsed improvement. Both modes retain the
existing 64-row task/list chunk policy. Additional arena queue metadata is
constant per pass, not per task or removable.

The optional private untimed probe reports peak pending, peak ready, peak
combined logical work, and maximum allocated arena rows. No probe or diagnostic
per-row counting runs in timed CLI use; the intrusive list splice counts its
rows only when a caller supplied that test/diagnostic probe. Measure this mode
with the user's default GC configuration, without GC tuning, and report elapsed,
peak RSS and allocation counts separately.

Source tests compare all three representations on complete passes (AST, output,
diagnostics, change stamps and VarInfo identities), plus shared builders, empty
lists, applied flags, enqueue-during-Apply, row release before reentrant work,
333 simultaneous tasks across chunk growth, 1000 reused work items, weak AST/
scope collection while the arena stays live, public/custom fallback and panic
cleanup. Normal/race tests and whole-CLI Java byte gates remain required before
any timing claim.
