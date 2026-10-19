# Performance experiment record

This publishes the saved experiments through the completed template-guard V2 screen (2026-10-19). It includes positive, negative, rejected, superseded and inconclusive outcomes, validation failures, source-only audits, measurements and prototype source. Continued exploration after this snapshot is not represented as completed evidence.

**Warm Java has not been beaten across the workloads. The 3.7% figure was the latest isolated typed-app improvement, not the overall improvement. No exhaustive search or proof that further gains are impossible is claimed.** This PR publishes evidence and reproduction material; prototype compiler changes in archives are not applied to the PR's compiler/library tree.

## How to read the numbers

The latest same-binary template experiment measured typed-app 5.123 → 4.934 s (3.7% less elapsed time), RSS 275.7 → 252.6 MiB. Checks-only measured 7.504 → 7.405 s (1.3% less time), RSS 653.4 → 658.8 MiB. These five-sample ranges overlap. [Detailed round 11](records/round2/template-guard-v2-round11.md).

For a broader comparison, the following table uses ONE complete five-task candidate, a6afc95f (`inline` in round 10), against the recorded current-main baseline da20b44f (compiler tree matching main 123b46e8). Every row uses that same candidate and configuration; this does not combine different per-task winning binaries.

| Workload | Main → candidate Go seconds | Less elapsed time | Main → candidate peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|
| typed-app | 5.795 → 4.954 | 14.5% | 257.3 → 250.7 | 0.755 | 4.35 |
| bundle | 20.307 → 18.449 | 9.1% | 612.5 → 432.5 | 0.900 | 3.92 |
| downlevel-whitespace | 1.467 → 1.199 | 18.3% | 164.9 → 170.3 | 0.783 | 5.46 |
| downlevel-simple | 4.404 → 3.749 | 14.9% | 199.6 → 212.6 | 0.748 | 5.26 |
| checks-only | 8.569 → 7.336 | 14.4% | 657.1 → 664.4 | 0.823 | 2.05 |

The equal-weight geometric mean of these five latency ratios corresponds to **14.3% less elapsed time**. That arbitrary equal-weight aggregate is a descriptive screening statistic, not a confirmed overall improvement, workload-weighted throughput claim, or improvement introduced by this documentation PR. The baseline and candidate are separate cohorts. Candidate source includes historical experimental substrate choices, including parser-owned node batches and vendored Unicode tables; publication does not reinstate rejected production choices. Later c7998d45 was remeasured only on two typed tasks, so there is no complete five-task result for that latest binary. Earlier work already integrated into main is also outside this main-to-prototype comparison.

Both comparison cohorts use Go 1.27.2, default runtime GC, the same five inputs, shipped PGO and the recorded build options/CPU setup; candidate experiment flags and source differ and are retained in metadata. Three Go observations per cell are a screen. Java is the older stored warm JVM reference (12 warmups/10 measured compiles), with fresh compiler/extern parsing each iteration. Ratios above1 favor Go. Go RSS is median fresh CLI lifetime peak; Java RSS is whole warm JVM lifetime peak including warmup. These memory windows differ. Never treat the ratio as equivalent per-compilation live heap. [Baseline](records/round2/baseline-summary.md) · [complete round10](records/round2/inline-bitset-template-round10.md).

## Archive contents

- **63 saved summary files**: [dataset index](datasets.json) and the list below. Some are validation or protocol summaries, not performance comparisons.
- **553 normalized task/variant metric rows** where the schema supplies them: [CSV](metrics.csv), [JSON](metrics.json). Missing historical RSS/ratio fields remain blank/null; they are not invented. Nonstandard schemas remain available in full in each original summary.
- **9038 measurement/support files**: [measurement-records.tar.gz](measurement-records.tar.gz), [per-file SHA-256 manifest](measurement-manifest.json). Includes raw per-run timing/RSS/digests, machine/configuration/provenance, scripts, build/validation/failure logs, CPU/heap profiles, reports and diagnostic overlays.
- **120 source worktree snapshots**: [prototype-sources.tar.gz](prototype-sources.tar.gz), [source manifest](source-manifest.json). Each tracked patch applies to main123b46e8; listed untracked source files are preserved too. Some snapshots are unfinished/unmeasured source, explicitly not endorsements.
- **Exact committed history**: [prototype-history.bundle](prototype-history.bundle) preserves experimental commits excluded from main, including intermediate measured implementations. It is a Git data file, not an applied merge. It complements working-tree patches so SHA references remain recoverable after the cloud workspace expires.
- **71 distinct design-note versions**: [index and source attribution](design-notes.json), browsable under `design-notes/`.
- **Archive hashes/sizes**: [archives.json](archives.json). [Verification utility](verify_archive.py).

Downloaded corpora, repeated generated JS/maps, compiled CLI/test executables/classes, symlink targets and PID files are excluded. Their pinned source references and measured artifact hashes remain in metadata; builds reproduce executables. Credential files and full process-environment dumps are excluded; retained text and non-main source blobs passed the recorded credential-pattern scan. Historical `state.md`/`progress.md` and old design notes are retained verbatim as contemporaneous records, not current instructions or validated estimates. The current policies and conclusions in this index take precedence over stale expectations in those records.

## Current conclusions and remaining opportunities

- Runtime GOGC tuning is no longer an exploration target. Older GC datasets are retained only as historical evidence. New code comparisons use default runtime GC/CPU settings.
- Persistent workers and parsed-extern reuse have been dropped from the active target. Earlier unfair cached-Go/fresh-extern-Java comparisons are not valid optimization evidence; matching-Java worker followups and that correction remain archived. The target is fresh non-worker CLI compilation.
- Two-inline-word BitSets are mixed: typed/simple screens improve, bundle regresses; no default promotion. Fresh-template V1 regressed badly because per-scope admission introduced expensive monitor stack walks. V2 removes that regression and shows modest net gains against its own control. Other pass/storage/scope/parallel factors and their compound screens remain recorded, including combinations that fail to add their isolated gains.
- Compatibility tests gate the exercised paths, not all possible Java/Go programs or external API behavior. Historical broad suites and newer focused suites have distinct scopes. Skips and failed fixture/build attempts remain visible. No full-port parity claim is inferred from a focused suite.
- **Further gains remain plausible.** A consistent five-task candidate still needs roughly10–25% less wall time, depending on task, to equal its stored warm-Java reference. That does not prove the remaining gap can be eliminated. Comfortable wins need more than matching Java.
- Next evidence-driven experiments: exact BitSet capacity/Clone occupancy before a one-word layout; local parser reuse of existing global string cells with actual hit/contention counts; larger pass-specific native data structures/algorithms; safe additional per-script parallelism; semantic-type construction/storage reuse with explicit aliasing/mutation/resolution proofs. Existing native CFG/array solver state and shared string headers must not be counted again as new ideas.
- Ambitious options remain research, not forecast speedups: typed/indexed internal representations, denser pass-local work/state, ownership-aware builders and broad traversal redesign. Any must preserve compiler artifacts, diagnostic/order/failure semantics and supported public/custom fallbacks, and must be tested against fresh matched controls. No fixed theoretical maximum is supported by the profiles.
- More cores may help some phases, but this machine has a two-CPU quota. Four compiler jobs are oversubscribed. Compiler job counts bound work items and do not globally cap goroutines or change GOMAXPROCS. [Fresh default-GC whitespace thread sweep](records/round2/default-gc-white-traversal-thread-sweep/summary.json). Full performance matrices for other configurations are not invented from compatibility gates.

Useful current audits: [type/flow design](records/round2/type-flow-design-audit.md), [node/string layout](records/round2/node-string-cell-layout-audit.md), [parser front cache](records/round2/parser-string-front-cache-audit.md), [one-word BitSet](records/round2/oneword-bitset-alternative-audit.md), [template guard](records/round2/fresh-template-guard-audit.md).

## Reproduce and inspect

Run `python3 docs/performance-experiments/verify_archive.py` from a normal full clone to verify archive hashes and the measurement manifest without executing archived code. Extract the tarballs into a separate directory. Their paths preserve the original record/worktree relationships; original absolute paths are historical and must be rebased for local reproduction.

Fetch the repository's complete main history before importing the bundle (a shallow clone may lack prerequisites). `git bundle verify docs/performance-experiments/prototype-history.bundle` checks prerequisites. Fetch the specific desired archived branch into a new local experiment branch, or use the tracked patch plus listed untracked files against the exact source-manifest base. Do not overwrite an active checkout with arbitrary archived patches. Build/test commands, runtime flags, toolchains, pinned Java jar and input references are in each dataset's machine/command/provenance files. Binary identities are retained even though binaries are not.

## Dataset inventory

Families marked historical-worker or historical-GC are excluded from the current optimization claim. Validation summaries contain compatibility results only; microbenchmarks are not whole-compiler speed claims. Consult the accompanying reports/raw metadata before comparing cohorts.

| Saved dataset | Family / qualification |
|---|---|
| [arena-smoke/summary.json](records/arena-smoke/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [cached-worker-profile-20261018/summary.json](records/cached-worker-profile-20261018/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [delivery-paired/summary.json](records/delivery-paired/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [first-candidates/summary.json](records/first-candidates/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [java-worker-compare-smoke/summary.json](records/java-worker-compare-smoke/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [java-worker-equivalent-protocol/summary.json](records/java-worker-equivalent-protocol/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [layout-smoke/summary.json](records/layout-smoke/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [matched-worker-whitespace-repeat/summary.json](records/matched-worker-whitespace-repeat/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [matched-worker-whitespace-threads/summary.json](records/matched-worker-whitespace-threads/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [memory-envelope/summary.json](records/memory-envelope/summary.json) | Historical GC settings; exploration suspended |
| [native-parser-micro/summary.json](records/native-parser-micro/summary.json) | Diagnostic / microbenchmark; not whole-compiler throughput |
| [next-smoke/summary.json](records/next-smoke/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [node64-paired/summary.json](records/node64-paired/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [node64-pgo-smoke/summary.json](records/node64-pgo-smoke/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [ordinary-cli-owned-screen/summary.json](records/ordinary-cli-owned-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [ordinary-cli-screen/summary.json](records/ordinary-cli-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [ordinary-cli-threads/summary.json](records/ordinary-cli-threads/summary.json) | Thread scaling; inspect environment and Java cohort |
| [parallel-cli-owned/summary.json](records/parallel-cli-owned/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parallel-compound-final/summary.json](records/parallel-compound-final/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parallel-expanded/summary.json](records/parallel-expanded/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parallel-fsiv-cli/summary.json](records/parallel-fsiv-cli/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parallel-gates/summary.json](records/parallel-gates/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parallel-inlining-only-screen/summary.json](records/parallel-inlining-only-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [parallel-inlining-screen/summary.json](records/parallel-inlining-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [parallel-modules-checked/summary.json](records/parallel-modules-checked/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parallel-zip/summary.json](records/parallel-zip/summary.json) | Compatibility / validation; not a timed speedup claim |
| [parser-tree-storage/summary.json](records/parser-tree-storage/summary.json) | Diagnostic / microbenchmark; not whole-compiler throughput |
| [production-worker/summary.json](records/production-worker/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [protocol-harness-smoke/summary.json](records/protocol-harness-smoke/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [redesigns-first-worker-protocol/summary.json](records/redesigns-first-worker-protocol/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [refreshed-smoke/summary.json](records/refreshed-smoke/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/allocations-screen/summary.json](records/round2/allocations-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/code-design-screen/summary.json](records/round2/code-design-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/compound-screen/summary.json](records/round2/compound-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/default-gc-white-traversal-thread-sweep/summary.json](records/round2/default-gc-white-traversal-thread-sweep/summary.json) | Thread scaling; inspect environment and Java cohort |
| [round2/followup-screen/summary.json](records/round2/followup-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/inline-bitset-template-screen/summary.json](records/round2/inline-bitset-template-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/native-live-words-screen/summary.json](records/round2/native-live-words-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/native-storage-code-screen/summary.json](records/round2/native-storage-code-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/native-storage-screen/summary.json](records/round2/native-storage-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/next-screen/summary.json](records/round2/next-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/parser-gc-screen/summary.json](records/round2/parser-gc-screen/summary.json) | Historical GC settings; exploration suspended |
| [round2/recipes-externs-screen/summary.json](records/round2/recipes-externs-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/structural-screen/summary.json](records/round2/structural-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [round2/template-guard-v2-screen/summary.json](records/round2/template-guard-v2-screen/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [second-candidates/summary.json](records/second-candidates/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [small-worker/summary.json](records/small-worker/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-cache-off/summary.json](records/thread-cache-off/summary.json) | Thread scaling; inspect environment and Java cohort |
| [thread-cache-on/summary.json](records/thread-cache-on/summary.json) | Thread scaling; inspect environment and Java cohort |
| [thread-worker/protocol/small-whitespace/threads-1/summary.json](records/thread-worker/protocol/small-whitespace/threads-1/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker/protocol/small-whitespace/threads-2/summary.json](records/thread-worker/protocol/small-whitespace/threads-2/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker/protocol/small-whitespace/threads-4/summary.json](records/thread-worker/protocol/small-whitespace/threads-4/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker/summary.json](records/thread-worker/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker-repeat/protocol/small-whitespace/threads-1/summary.json](records/thread-worker-repeat/protocol/small-whitespace/threads-1/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker-repeat/protocol/small-whitespace/threads-2/summary.json](records/thread-worker-repeat/protocol/small-whitespace/threads-2/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker-repeat/protocol/small-whitespace/threads-4/summary.json](records/thread-worker-repeat/protocol/small-whitespace/threads-4/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [thread-worker-repeat/summary.json](records/thread-worker-repeat/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [toolchain/summary.json](records/toolchain/summary.json) | CLI candidate/toolchain screen; inspect per-variant controls and sample count |
| [type-flow-cli-next/summary.json](records/type-flow-cli-next/summary.json) | Compatibility / validation; not a timed speedup claim |
| [warm-combined/summary.json](records/warm-combined/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [warm-next127/summary.json](records/warm-next127/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [worker-protocol/summary.json](records/worker-protocol/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |
| [worker-retention/summary.json](records/worker-retention/summary.json) | Historical worker / protocol; superseded target, inspect Java matching |

## Other saved baseline/report entry points

- [baseline/report.md](records/baseline/report.md) — historical snapshot; see current policies above.
- [progress.md](records/progress.md) — historical snapshot; see current policies above.
- [state.md](records/state.md) — historical snapshot; see current policies above.
