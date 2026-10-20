# Private compact contextual name restoration

This is an unmeasured, default-off CLI experiment. Set both
`JSGO_EXPERIMENT_COMPACT_INVERTER=1` and
`JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1`. The latter provides the existing fresh
stock-CLI ownership proof. `INVERTER_NOOP` is independent; when it skips the
pass, there is no compact invocation.

The current default-GC bundle heap profile attributes 313.56 MiB and 1.88 million
sampled allocations to contextual name restoration. Adding potential shadow
variables accounts for 135.88 MiB; constructing the temporary disallowed-name
sets accounts for 125.15 MiB. The pass owns 0.76/21.18 seconds of sampled CPU.
Those 6.5% of sampled bytes and 3.6% of CPU are upper bounds on the opportunity,
not predicted savings. Sampled heap estimates differ from exact MemStats totals.

Each rename record stores its shadow encounter sequence as contiguous uint32
identity IDs. An invocation-private pointer table retains the real GC-visible
VariableInfo objects. Short shadow lists deduplicate linearly; larger lists add
a native uint32 membership index. Neither shadow links nor these membership
maps contain AST pointers. This removes generic linked-set entries and their
linked-list/hash storage. IDs are never published to AST nodes, Scope, or Var.

The temporary disallowed-name collection needs membership only. It uses a native
string set for valid UTF-8 names, whose byte encoding is unique for their Java
UTF-16 sequence. Malformed/WTF-8 inserted names or membership queries select the
existing port's linked set, preserving its original Go-byte key policy. This
does not repair the existing collection's noncanonical UTF-16 equivalence gap.
Names are read from the current
VariableInfo once per distinct shadow, in their original insertion order,
immediately before choosing a replacement. Earlier renames therefore remain
visible. The string set is not retained after the decision.

Storage belongs to one outer local-scope cohort. At the exact point the original
pass clears its name stacks, the experiment drops all records, identity tables,
and membership indexes. This can retain some already-renamed sibling records
longer than their original linked sets, until that cohort ends. It does not
retain old cohorts or use a process/global pool. Admission is capped at 65,536
variables, 65,536 rename records and 1,048,576 distinct shadow links per cohort.
Reaching a cap, or observing a custom VariableInfo implementation, materializes
all live original linked sets in encounter order without reading names and
retires the experiment until the next independent cohort.

Scope and Var layouts do not change. RenamingInfo gains a uint32 identity handle;
its source layout grows from 72 to 80 bytes, within the same 80-byte allocation
class. Contiguous record rows contain an owner pointer, a uint32 slice and an
optional native membership map. Their owner pointer permits bounded one-way
fallback without a global registry. Optional caller-owned counters record
actual dispatch, retirement, name observation order and storage peaks; they are
nil in timed compilation.

Only the private default-pass factory can enable this path. Public constructors
retain the original representation. Eligibility excludes foreign compilers and
roots, custom passes, reporters, conventions, redeclaration callbacks, warning
guards, observers, debug/tracing and typed-AST injection. Success and exception
paths clear the invocation via defer. Tests compare full AST/output/diagnostics
and change ordering, exercise current-name mutation and identity deduplication,
caps/custom fallback, malformed names, weak graph lifetime, exceptional cleanup,
and fresh N=1/N=2 DoRun output parity. Full Java CLI and race gates remain
required before any performance claim.
