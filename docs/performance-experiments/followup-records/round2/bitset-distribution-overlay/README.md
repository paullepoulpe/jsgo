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

The report contains constructor call counts by default/reserved/Clone/owned-word
helper and word capacity 0/1/2/>2, Clone source capacity versus source used length,
and actual capacity-growth transitions. It also preserves exact requested-bit /
capacity constructor histograms, exact Clone source-capacity / used-length /
result-capacity triples, and exact from/to capacity transitions. The growth-word
sum counts backing capacity added; the separately named allocated-backing-word
sum counts each whole new backing. Neither is malloc bytes: allocator rounding
and escape analysis also matter.

These are method/event distributions, **not heap-object distributions**. Calls
include package initialization. Bare zero-value creation is unobserved. In
particular, native liveness creates a constructor temporary then copies its
value into a contiguous array; constructor counts there do not imply separate
header allocations. Clone frequently discards reservation, so inspect both the
constructor and source/used matrix when estimating the empty/large tradeoff.
Use actual TotalAlloc/Mallocs and peak RSS from separately measured ordinary
binaries to decide whether the inline representation helps.

The frozen inventory source for `f2c17056` is in
`native-analysis-f2c17056/overlay.json`; its manifest hashes all original and
replacement files, and the generating Python source is copied beside it.
Native liveness must be OFF for the first inventory, matching the generic-engine
control in recent screens. Its publisher unconditionally invokes the V2 owned-
word helper when native liveness is on; assert the helper row is0 in every first
inventory report. Keep the existing native9 + CLI_SOURCE_LOCK controls with all
other new factors off. One-word representation and actual-bound native scratch
are separate later factors.

The one-word24B alternative's optimistic constructor-byte model versus original
is `-8*C1 +24*C2plus`, where C1 counts capacity1 and C2plus counts capacity≥2
constructor results (including Clone). Empty constructors are byte-neutral.
Capacity0→1 growth saves another8 bytes; first promotion0/1→≥2 adds24 bytes for a
holder. Growth that was already large has the same backing allocation. This is
an event model, not actual heap savings: nonescaping allocations, header arrays,
retention and allocator accounting still require whole-binary measurement.

Generator smoke check and generated Go formatting passed. No diagnostic binary
has been built or run yet.
