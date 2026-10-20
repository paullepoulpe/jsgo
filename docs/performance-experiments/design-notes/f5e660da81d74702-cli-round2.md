# Ordinary CLI performance, current-main round 2

Recorded from the October 18–19, 2026 screening artifacts. Warm Java remains
faster on all five tasks in this screen. The Go variants use less resident
memory under the process windows below. Native paths help several tasks, but
there is no universal win, and the compact-storage compound is not an automatic
keeper.

## Measurement contract

Each Go variant ran three fresh ordinary CLI processes per task. Variant order
was interleaved with seed 20261018; tables show median elapsed seconds and median
process peak RSS in MiB. The native/storage screen contains 45 runs. Compiler
threads were 1, affinity was CPUs 0 and 1, and the cgroup CPU quota was 2 CPUs
(`cpu.max=200000 100000`); the visible cpuset was 0–2. Runtime GC and Go CPU-budget
detection used their defaults, with benchmark GC/scheduler overrides removed.
Compiler work-item threads and GOMAXPROCS are distinct settings.

The environment was Debian 13, AMD EPYC 9V74, Go 1.27.2 linux/amd64, and OpenJDK
21.0.12.1. Builds use the repository's shipped PGO profile, whose SHA256 is
`f6aee9c5761f48f8e96beb0ee7cf5d4d8ef5f631460b2102d785d46885bdeaac`.
Build-only GC/CPU overrides are not runtime benchmark settings.

Java numbers reuse one stored cohort per task: 12 warmups followed by 10 measured
iterations in one JVM, with median measured time. Stock Java creates a fresh
compiler and performs normal extern parsing each iteration. These are ordinary
CLI compilations; worker/extern-preparse caching is outside this experiment.
Java was not rerun as a paired cohort for each Go variant, so the ratios provide
screening context, not independent-cohort confidence estimates.

Go RSS is the median peak of fresh CLI processes. Java RSS is the stored JVM's
whole-lifetime peak, including warmups and JIT/heap state. This is a process
footprint comparison with different lifetime windows, not incremental memory
per compilation or a heap-allocation comparison. Both ratio directions are
**Java / Go; above 1 favors Go**. A time ratio below 1 means warm Java is faster.

## Native and compact-storage screen

`main` is the current-main baseline. `native80` combines the earlier static
allocation/lookup batch with the experimental native CFG, local-flow,
reference-storage, inline-argument facts, and bounded FSIV paths, retaining
80-byte nodes. `storage64` includes those changes plus 64-byte per-kind AST
storage, adaptive tiny inferred-slot maps, and qualified-name lookup changes.
This is a compound comparison, not an isolated Node-layout ablation.

| Task | Java warm s | Java RSS MiB | main s / MiB | native80 s / MiB | storage64 s / MiB |
| --- | ---: | ---: | ---: | ---: | ---: |
| typed-app | 3.739 | 1090.2 | 5.731 / 260.2 | 5.689 / 257.3 | 5.628 / 237.8 |
| bundle | 16.600 | 1696.6 | 20.933 / 546.0 | 20.050 / 623.1 | 20.112 / 643.6 |
| downlevel-whitespace | 0.938 | 930.2 | 1.422 / 171.9 | 1.461 / 157.9 | 1.450 / 152.2 |
| downlevel-simple | 2.805 | 1118.4 | 4.416 / 202.0 | 4.059 / 201.0 | 4.255 / 189.2 |
| checks-only | 6.038 | 1360.0 | 8.525 / 659.8 | 8.242 / 668.1 | 8.466 / 604.9 |

| Task | main Java/Go time / RSS | native80 Java/Go time / RSS | storage64 Java/Go time / RSS |
| --- | ---: | ---: | ---: |
| typed-app | 0.65x / 4.19x | 0.66x / 4.24x | 0.66x / 4.58x |
| bundle | 0.79x / 3.11x | 0.83x / 2.72x | 0.83x / 2.64x |
| downlevel-whitespace | 0.66x / 5.41x | 0.64x / 5.89x | 0.65x / 6.11x |
| downlevel-simple | 0.64x / 5.54x | 0.69x / 5.56x | 0.66x / 5.91x |
| checks-only | 0.71x / 2.06x | 0.73x / 2.04x | 0.71x / 2.25x |

Compared with native80, storage64 is about 4.8% slower on downlevel-simple
(4.059 → 4.255 s), with lower RSS (201.0 → 189.2 MiB). Checks-only is also slower
(8.242 → 8.466 s) while RSS falls (668.1 → 604.9 MiB). Typed-app improves slightly in
both measures. These mixed results do not establish compact storage as a keeper.

Bundle RSS is particularly noisy: individual peaks in this screen span roughly
513–568 MiB for main, 557–690 MiB for native80, and 425–666 MiB for storage64. Its
median rises for both compounds despite some allocation reductions. Do not
infer a dependable bundle-memory improvement from a smaller Node or from one
median. Three runs are an initial screen, not a statistical confidence bound.

The native variants enabled exactly these flags:

```text
JSGO_EXPERIMENT_NATIVE_CFG=1
JSGO_EXPERIMENT_NATIVE_LOCAL_FLOW=1
JSGO_EXPERIMENT_NATIVE_REFERENCES=1
JSGO_EXPERIMENT_INLINE_ARGUMENT_FACTS=1
JSGO_EXPERIMENT_PARALLEL_FSIV=1
```

The native parser is **not included**. At compiler threads 1, a parallel-FSIV flag
being enabled does not demonstrate multi-worker speedup.

## Earlier allocation-only screen

This separate interleaved three-run screen compares main with the earlier
static micro-optimization batch, with experiment environment flags unset.
That batch is included in native80 and storage64: variable-inlining memo cells,
modified-parameter accumulation, empty peephole snapshots, skipping unnecessary
backward-use analysis, concrete Scope.GetVar lookup, and candidate-traversal
reuse. Its compound result cannot assign a gain to each individual change.
The main medians below belong to this earlier screen; do not mix them with the
main column above as if measured in the same cohort.

| Task | main s / MiB | allocations s / MiB | allocations Java/Go time | allocations Java/Go RSS |
| --- | ---: | ---: | ---: | ---: |
| typed-app | 5.687 / 256.5 | 5.817 / 256.6 | 0.64x | 4.25x |
| bundle | 20.523 / 558.5 | 19.658 / 630.1 | 0.84x | 2.69x |
| downlevel-whitespace | 1.446 / 168.3 | 1.412 / 168.7 | 0.66x | 5.52x |
| downlevel-simple | 4.192 / 203.0 | 4.161 / 195.1 | 0.67x | 5.73x |
| checks-only | 8.620 / 659.8 | 8.434 / 671.2 | 0.72x | 2.03x |

The allocation batch improves bundle elapsed time while increasing its RSS in
this cohort. Typed-app regresses slightly. The separate one-shot baseline
recorded 5.795/20.307/1.467/4.404/8.569 s for the five tasks, respectively; the
interleaved screen's own main measurements are the relevant Go-to-Go control.

## Source and binary provenance

| Variant | Compiler source checkpoint | Binary SHA256 |
| --- | --- | --- |
| main | da20b44f, compiler tree matching main 123b46e8 | `5f4dc046a4609c6ea27953c10da5d60bed265596a7ff48df712507e79ae9d5d9` |
| native80 | ef9657d7 | `503ae718e8cd225a9c21269b14e791284d6a53bcf20e17bca9a0c9316274f2aa` |
| storage64 | bd306adb | `f149fca868190a5465085225305006d7e54189ad383d299270937cb2f9a4d350` |
| allocations, earlier screen | earlier combined binary recorded by allocations-screen/machine.json | `d520e8c51a10aaa3a639879bbaa35c39f453ce18f7abb52e6702e7e40073bb7f` |

The native/storage screen's benchmark-runner commit is also bd306adb; it does not
mean every variant used that compiler source. The pinned Java jar SHA256 is
`d6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6`.

## Correctness and remaining work

Every summary cell reports `all_outputs_match=true`, including the 45
native/storage and 30 allocation-screen process runs. The baseline artifacts
passed all five stock-Java gates. Additional complete CLI gates recorded zero
artifact mismatches and successful exits for the combined allocation binary at
N=1/N=2, native80 at N=1, and storage64 at N=1/N=2. N=2 here is correctness coverage;
the timing tables use N=1 and make no thread-scaling claim.

Final focused logs include 437 PASS entries for storage64 native CFG/type
inference, 1164 for compiler/reference/scope coverage, and 135 for facts/FSIV
cohorts. The earlier allocation/reference log records 820 PASS entries. These
counts include nested subtests. A real-input FSIV coverage test was skipped when
its optional input environment variable was absent. Initial native and
facts/cohort logs contain fixture failures; repaired final logs pass. This is
focused coverage plus CLI parity, not a claim that the entire repository suite
was rerun here.

The native parser, scalar-first BitSet, additional PGO training/tuning, and newer
micro experiments remain pending and unmeasured in these tables. The latter
include RemoveUnused logging/prepass work 45fd0264/6f76aea7, native-query
8676d62d/61a0ec3d, and newer callee/scope-sharing phase work. Their unit/oracle
results do not establish a CLI speed or RSS gain. The next decision requires
fresh compound gates and matched interleaved measurements; no PR #23 edit or
retention decision is implied by this record.

Raw records are under `/workspace/jsgo-results/round2/`:

- `native-storage-screen/summary.json`, `machine.json`, and per-task runs;
- `allocations-screen/summary.json`, `machine.json`, and per-task runs;
- `baseline-summary.md`, `baseline/`, and `build-main-provenance.json`;
- `gates.log`, `gates-native80.log`, and `gates-storage64.log`;
- `tests-storage64-native.log`, `tests-storage64-compiler.log`,
  `tests-storage64-facts-cohorts.log`, and `tests-allocations-references.log`.
