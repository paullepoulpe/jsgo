# Untimed BitSet opportunity diagnostic

This overlay instruments the original slice-backed BitSet. It changes no
production source and does not force tracing, alter compiler eligibility guards,
retain set pointers, or tune GC. Atomic counters are diagnostic overhead; never
use the resulting binary for performance comparisons.

Generate against the exact baseline checkout:

```sh
python3 generate.py --repo /absolute/baseline/checkout --out /absolute/output
```

Then, only in an authorized build slot, pass the resulting `overlay.json` to
`go build -overlay ... ./cmd/jsgo`. Run with an absolute
`JSGO_BITSET_DIAGNOSTIC_PATH=/.../counts.json` on the same CLI fixtures/flags as
the normal screen. Both normal main return and Java-style System.exit dump
once; a direct os.Exit in an uncaught Go exception can omit the report.

The report contains constructor call counts by default/reserved/Clone and word
capacity 0/1/2/>2, Clone source capacity versus source used length, and actual
capacity-growth transitions. The growth-word sum counts backing capacity added,
not bytes allocated: reallocations allocate the entire new backing, allocator
rounding and escape analysis also matter.

These are method/event distributions, **not heap-object distributions**. Calls
include package initialization. Bare zero-value creation is unobserved. In
particular, native liveness creates a constructor temporary then copies its
value into a contiguous array; constructor counts there do not imply separate
header allocations. Clone frequently discards reservation, so inspect both the
constructor and source/used matrix when estimating the empty/large tradeoff.
Use actual TotalAlloc/Mallocs and peak RSS from separately measured ordinary
binaries to decide whether the inline representation helps.

Generator smoke check and generated Go formatting passed. No diagnostic binary
has been built or run yet.
