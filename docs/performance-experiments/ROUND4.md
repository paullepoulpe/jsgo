# Round 4: verified baselines and enterprise benchmark protocol

Investigation resumed on 2026-10-20 from main
`2262bb6a8a9d24ee34fdb52bd30c49269e4b5373`. At 08:25 UTC, remote main
had advanced to `289058def0873f349f80e49e6777afe3091e18a9`: five
optimization commits, 49 changed files and a new shipped PGO profile. That
main is now merged into PR 23. The original baseline and seven-factor screen
belong to `2262bb6a`; later sync results below identify their own source.
This page records completed baseline measurements and ongoing experiments;
no experimental optimization has been promoted. Earlier rounds retain their original
source, machine, PGO and measurement cohorts; their percentages are not
additional improvements over this main.

## Subsequent main syncs and completed 289058de results

Main advanced again to `5de6e3359349d70e60c98788e5e62b78cf1cb0e7`:
six further commits, 32 changed files and PGO SHA256
`e11b868352b3ed0d585686fa63b47f3d04e62259535497237d4804c78a060462`.
It is now merged into PR 23 and frozen for the next paired-JVM block.
Its bundle N1 and N2 correctness invocations match all six stock-Java
artifacts under the actual 4 GiB cap. Its full suite is in progress.
The completed N1 paired-JVM block below establishes a win on this bundle.

### Completed 5de6e335 warm-JVM comparison

Three independently restarted JVMs each settled after 21 warmup requests,
then completed five alternating Java/Go pairs. All 15 pairs match every
required output byte, stream and exit. Each Java request creates a fresh
stock compiler and reparses externs; Go is a fresh CLI. Both use compiler
N=1, default runtime GC, a two-CPU budget and a verified 4 GiB kernel cap.

| JVM cohort | Java median s (SD) | Go median s (SD) | Java peak RSS MiB | Go median peak RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|---:|
| 1 | 11.573 (0.323) | 9.567 (0.190) | 1555.47 | 381.54 | 1.2096 | 4.0768 |
| 2 | 11.334 (0.178) | 9.304 (0.187) | 1332.92 | 381.93 | 1.2182 | 3.4900 |
| 3 | 11.596 (0.115) | 9.252 (0.321) | 1477.89 | 384.56 | 1.2534 | 3.8431 |

The geometric mean Java/Go time ratio is **1.2269**, with a 95% restarted-
cohort bootstrap interval of **1.2096–1.2534**. Go takes **18.5% less elapsed
time** (geometric means of cohort medians: Java 11.500 s, Go 9.373 s).
Java/Go peak RSS is **3.7956** (interval **3.4900–4.0768**), corresponding to
**73.7% less peak RSS** for Go. Java RSS covers the warm JVM lifetime,
including warmup and JIT; Go RSS is the median fresh-process lifetime peak.
These are process-memory comparisons, not per-compilation live-heap ratios.

All three cohorts favor Go, with no outer OOM-counter increase. Only three
independent JVM cohorts underpin the intervals; this is a confirmed result
on the 11-library ADVANCED bundle, not a result for the larger typed,
many-chunk or generated-code workloads. The compiler and PGO are unchanged
stock main: **the measured win belongs to changes landed on main**, not to
the unpromoted prototypes archived in this PR. Raw evidence is
`paired-main5de-bundle-advanced-t1-v1/` in the third measurement delta below.

A separate three-round Go-only screen compares `289058de` with `5de6e335`:
median elapsed time 10.260 → 8.884 s (SD 0.213 → 0.124 s), and median peak
RSS 378.92 → 390.05 MiB (SD 0.35 → 8.16 MiB). Median same-round new/old
ratios are **0.86585 time** and **1.02989 RSS**: 13.4% less time with about
3.0% higher RSS. All six invocations match the same Java oracle. This screen
and the earlier 2262→289 screen are separate sessions; their percentages
must not be multiplied into a purported measured overall improvement.

The full `5de6e335` suite was explicitly paused from 09:18:49 to 09:42:48 UTC
for the quiet timing block, then resumed automatically. Its wall time is
not a performance measurement. A post-block fetch confirmed main still
points to this measured source.

The intervening main `289058de` passed the complete all-package normal
suite, including the test-only fail-closed corpus paths: 29,508 PASS events
(3,802 top-level), 103 skipped events (37 top-level), zero failures across
21 packages. All 3,839 expected top-level names/statuses match exactly;
no missing, extra or changed outcomes. All 553 effective test-source files
were retained. Its explicit 08:43:27–08:48:25 pause is recorded; unit-suite
wall time is not a performance result. Full live-Java parity remains open.

A quiet three-round alternating **Go-versus-Go** bundle screen found:

| Source | Median seconds (min–max) | Median peak RSS MiB (min–max) | Time SD seconds | RSS SD MiB |
|---|---:|---:|---:|---:|
| 2262bb6a | 11.757 (11.730–11.760) | 423.05 (418.14–436.61) | 0.0169 | 9.57 |
| 289058de | 10.000 (9.934–10.056) | 379.67 (376.25–381.90) | 0.0606 | 2.84 |

The median same-round new/old time ratio is 0.8525 (**14.7% less time**);
RSS is 0.8975 (**10.3% less peak RSS**). These are improvements already
landed on main, not wins from our unpromoted experimental factors. All
six runs match the same stock-Java oracle with identical argv/input identity.
This screen supplies no fresh Java timing/RSS ratio; combining it with
a prior Java cohort would not establish a new paired result.

Fresh 289058de typed profiles and full/pair AST censuses now also match all
six fresh stock-Java artifact digests; input manifests and argv match.
The new metadata sharing reopens compact-node investigation. A layout model
projects only a few MiB of slab savings for full 64+24 records versus
8.60 MiB typed / 4.21 MiB bundle for 72+16 pair records, depending on slab
versus standalone retention. These are structural models, not RSS/speed
wins. A bounded COW diagnostic is source-only. Subsequent 5de6e335 changes
require refreshed performance weights.

SourceFile focused race validation passed all 42 top-level / 102 total
checks with zero race warnings. A hex-spelling test expectation was corrected
to supported decimal exponent spelling; the failed attempt remains retained.
The combined DAE/scope/CFG candidate built normally and reached every one
of 80 selected top-level tests: 79 passed and one test fixture panicked
while bypassing warning-guard public ordering bookkeeping. Corrected public
AddGuard and custom-source escape regression fixtures are frozen, unrun.
No experimental compiler optimization is promoted.

Generated protobuf CommonJS **SIMPLE** now passes Java/Go byte gates and
independent execution of both compiled outputs, including real descriptor,
Timestamp and Struct binary roundtrips. The earlier ADVANCED CommonJS
fixture emitted byte-identical but broken code in both engines: preserve
that failure, exclude it from performance claims. A separate genuine
Closure-style protoc 3.20.3 fixture has now passed generation, raw execution,
all six ADVANCED Java/Go byte gates and independent execution of both
compiled outputs under the verified cap (zero OOM events). This is a
correctness diagnostic; it supplies no settled-JVM performance ratio.
Soy generation, J2CL inputs and dependency-preserving 32-chunk construction
are prepared; their remaining gates are explicitly pending.

## Verified base and remaining parity work

The full, unfiltered `go test ./...` completed with exit 0 using Go 1.27.2:
29,507 passing test/subtest events (3,801 top-level), zero failures and
104 skipped events (38 top-level). The command included every package and
all original test files, without `-short`, `-run`, source overlays or cached
test results. Build resource controls were removed before executing test
binaries; the runtime used default GC settings. Shipped PGO SHA256 is
`8befcc7bbc99bdc3452102d14115ef7db0ceb69d097535fe332465d804931f67`.

This is not complete live-Java parity. Some checked-in golden helpers use
an absolute corpus directory and silently omit corpus cases on another
machine. A test-only, fail-closed path override for eight affected suites has now
passed all 556 additional corpus-file/suite pairs (zero missing inputs,
failures or skips). An earlier progress count of 656 was an addition error;
556 is the verified total. All original test files and vectors are retained.
The initial command and its reduced corpus coverage remain recorded.
The full race suite and live Java CLI/corpus oracles remain outstanding.

The complete pinned held-out corpus is now fetched and checksum-verified:
16 packages, 10,101 allowed JS/MJS/CJS cases after the repository's 44 GPL
exclusions. Development corpus has 14,690 files. Coverage inventories count
expected jobs explicitly: the CLI harness can return success even when
its summary contains mismatches, and a default sample is not a full run.

## Benchmark protocol

See [the runner specification](../enterprise-benchmark.md). New paired trials
use the pinned, unchanged Java compiler, fresh CommandLineRunner/Compiler and
normal extern parsing on every JVM request. Go uses a fresh CLI process.
There is no worker/extern cache asymmetry and no runtime GC tuning.

Both engines run in separate kernel-enforced 4 GiB, zero-swap cgroups with
two CPU quota and affinity 0,1. A small independent OOM canary verified actual
enforcement. Java has an explicit 3 GiB maximum heap within its total cap.
The runner verifies observed kernel settings, not just Docker metadata.
The pinned image and external Debian JDK configuration files are fingerprinted.
An actual Java READY/STOP protocol smoke check now precedes compiler trials.

Trials alternate Java and Go in the same session, with randomized initial
order and at least three independently restarted JVM cohorts. Warmup starts
at 12 requests and extends to at most 30; rolling time drift, spread and JIT
compilation activity determine whether a cohort settled. Unsettled warmup
is inconclusive. The idle JVM is paused during Go trials. Its resident pages
remain charged; outer memory pressure and OOM events are retained.

Reports retain all samples, elapsed-time variance, peak process RSS, cgroup
peaks, and Java/Go ratios for time and RSS. Confidence intervals use restarted
JVM cohort medians. Java RSS covers its warmed process lifetime, including
warmup/JIT; Go RSS covers each fresh CLI lifetime. These windows are explicit.
All stdout, stderr, exits, output files, source maps and manifests must match;
expected absent artifacts are also checked. OOM/failure is not a speedup.

The large ADVANCED fixture has 1,820 actual inputs, 25,919,443 bytes: 971
Closure Library, 522 Shaka, 324 Blockly and three app files. A separate
four-chunk task preserves all inputs and records dependency/output manifests.
These exact source counts supersede descriptive counts in older benchmark
documentation. Genuine J2CL, Soy and protobuf generated-code fixtures are
being prepared; incomplete runtime dependencies are not replaced with stubs.

## Completed warm-JVM baseline

The `bundle-advanced` task compiles 11 genuine UMD libraries totaling
5,677,811 input bytes, with ADVANCED, checkVars, source maps and manifests.
Both engines use compiler thread count 1 under the two-CPU budget. Three
independently restarted JVM cohorts settled after 20, 20 and 21 requests;
each then ran five alternating Java/Go pairs. Every output byte matched.

| JVM cohort | Java median seconds | Go median seconds | Java peak RSS MiB | Go median peak RSS MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|---:|
| 1 | 11.190 | 11.474 | 1460.9 | 431.4 | 0.9752 | 3.3866 |
| 2 | 10.884 | 11.618 | 1707.5 | 419.6 | 0.9368 | 4.0689 |
| 3 | 10.982 | 11.518 | 1411.3 | 427.1 | 0.9535 | 3.3042 |

The geometric mean Java/Go time ratio is **0.9551** (95% restarted-cohort
bootstrap interval 0.9368–0.9752): Go is about **4.7% slower**. Java/Go RSS
is **3.5708** (3.3042–4.0689): Go uses about **72% less peak process RSS**.
Greater than one favors Go for both ratios. Only three independent JVM
cohorts underpin these intervals; this is one workload, not broad parity or
a claim of enterprise-wide reliability. Raw samples retain within-cohort
variance. The RSS windows described above differ intentionally: warmed JVM
lifetime versus each fresh Go CLI lifetime. No outer OOM counter changed.

## Completed seven-factor screen

Thirty fresh capped Go CLI runs tested stock main, an all-off experimental
binary, seven independent factors, and their compound configuration. Three
randomized balanced rounds retained every sample. Every row matched the
saved stock-Java artifact oracle, with exact argv and input identity checks.
This is a **Go-versus-Go exploratory screen**; it supplies no new Java timing
or memory ratios and cannot substitute for live paired JVM confirmation.

| Factor | Median seconds | Median peak RSS MiB | Median same-round time/control | Median same-round RSS/control |
|---|---:|---:|---:|---:|
| stock-main | 11.754 | 430.5 | 0.9975 | 1.0041 |
| factor-off | 11.737 | 426.3 | 1.0000 | 1.0000 |
| queue | 11.883 | 427.5 | 1.0046 | 1.0003 |
| typed-scope | 11.814 | 429.1 | 1.0066 | 1.0064 |
| preorder | 11.724 | 426.1 | 0.9939 | 0.9974 |
| inline-cost | 11.763 | 411.6 | 0.9998 | 0.9584 |
| interference | 11.776 | 425.6 | 1.0090 | 1.0117 |
| scope-slabs | 11.609 | 426.7 | 0.9951 | 0.9982 |
| inline-schedule | 11.770 | 424.2 | 0.9967 | 0.9955 |
| compound | 11.772 | 428.7 | 1.0029 | 0.9950 |

Control ratios here use **Go variant / Go all-off control**, so smaller is
better; they are medians of same-round ratios, not ratios of table medians.
There is no clear compound speed gain. Inline-cost views show a roughly 4%
RSS reduction in this small screen, warranting confirmation on other tasks.
Small time differences remain exploratory, including the Scope/Var slabs.
The seven-factor targeted plus corpus gate passed 310 test/subtest events,
42 top-level, with no skips/failures; its full unfiltered jscomp suite subsequently passed with 27,434 passing
events, 1,540 passing top-level tests, 95 skipped events and zero failures.
All 1,535 original top-level tests were present, plus 34 new tests; original
statuses match the stock baseline after the supplemental corpus gate. No optimization is promoted based on this screen.

## Fresh diagnostic evidence

Quiet, capped CPU/allocation profiles of current main were collected for
typed-full, bundle-advanced and checks-only. The binary replaces only main
to flush diagnostic profiles before the original exit receiver; compiler
libraries and shipped PGO are unchanged. These are diagnostic runs, not
timed Java comparison results. Bundle profile artifacts subsequently match
the paired stock-Java oracle; the typed/checks-only profiles are not yet
Java artifact-gated.

In typed-full, `Node.GetNext` accounts for 14.09% flat CPU samples,
`NodeTraversal.TraverseBranch` 4.79% flat / 58.47% cumulative, and `gcDrain`
15.05% cumulative. Cumulative figures overlap and cannot be added.
`GetNext` is already one inlined pointer load, so its attribution motivates
locality investigation rather than predicting a removable-call win.
Read-only decoding identifies madvise at 0.02 sampled CPU seconds in each
profile: 0.0417% of typed CPU and 0.1405% of bundle CPU, all through the
background Go scavenger. Generic syscall samples mostly attribute to source
file reads/I/O. This does not reproduce the earlier 14% madvise hypothesis.
A capped strace capture remains outstanding; old profiles are not evidence
for syscall weights on the new 289058de source.

Read-only disassembly also rejects two tempting duplicate optimizations:
ordinary string maps already inline to Go fast-string operations, and PGO
already guards and inlines ReferenceCollector.Visit in the main traversal.
Current Node headers are 80 bytes. A bounded census of actual property,
string-cell and extension identities found weak sharing. Even ideal 24-byte
shared cold records save only 2.96% of finalized bundle structural bytes and
1.95% of typed-full structural bytes; 32-byte records increase bytes. These
are structural projections, not RSS or speed gains. Added objects, indirection
and copy-on-write costs make a wholesale AST64 rewrite unattractive on this
evidence. A narrower property/string-only census was prepared but is not built.
Main 289058de materially changes Node original-name storage, property-list
sharing and traversal pooling while retaining the 80-byte header. The old
extension-driven census cannot be used to reject layouts on that new source.
It needs a new kind-aware census and fresh profile.

## Candidates and failures retained

No compiler optimization is promoted by this progress update. Independently
switchable source candidates cover initial data-flow queue reservation,
redundant TypedScope own-miss lookup, audited constant-true preorder dispatch,
InlineFunctions cost snapshots, bounded interference scratch and stock
Scope/Var allocation. Tests/builds and macro measurements are separate gates;
source plausibility is not a measured improvement. A shared SourceFile
line-cache race fix is kept separate from performance factors. Java probes
also exposed restored-metadata edge behavior that its regression tests must
preserve.

Two harness failures are preserved: checks-only does not emit its requested
output manifest, and the first containerized Java startup lacked external
Debian JDK configuration symlink targets. The fixes preserve compiler argv,
check the manifest's expected absence, mount/fingerprint exact dependencies,
and verify the real protocol startup. Neither failed attempt produced a
valid Java speed ratio.

## Reproducible first evidence snapshot

[round4-measurement-records-v1.tar.gz](round4-measurement-records-v1.tar.gz)
preserves 1,005 completed raw files (672,050,356 source bytes), including
compiled JS/maps, individual samples, profiles, original failures and full
unit logs. SHA256 blobs deduplicate identical artifacts; no samples are
discarded. Six rebuildable binary/class files are excluded with their hashes
and sizes recorded. Input symlinks are recorded as metadata, without following
or recreating them; payload inventories remain pinned in the run records.
Active SourceFile race and new-main runs are excluded from this first cut.

Archive SHA256: `3c7512fe7197ed7d4a6de18ef996ae69ce979d12ef47b70b498b5759b82b5008`.

```sh
python3 -B docs/performance-experiments/verify_round4_records.py \
  docs/performance-experiments/round4-measurement-records-v1.tar.gz
# Optional: append --restore /tmp/jsgo-round4-records (must not exist).
```

The verifier validates every path/blob/hash and never executes archived code.
The publisher's six synthetic tests cover roundtrip/deduplication, changing
inputs and rejected malformed archives. The benchmark tooling's 15 tests pass.

[round4-source-manifest.json](round4-source-manifest.json) fingerprints nine
frozen source snapshots in
[round4-prototype-sources.bundle](round4-prototype-sources.bundle), with exact
commits, trees, changed-file hashes and recovery commands. The prerequisite
is 2262bb6a, already in this PR. These include tested seven-factor code and
separate unvalidated concurrency/layout diagnostics; a snapshot is not a
promotion or proof of correctness. The SourceFile snapshot predates the
new race-run fixture correction. Later cuts will retain that correction and
the final race evidence. Earlier archives remain byte-for-byte unchanged.

## Second completed evidence cut

[round4-measurement-records-v2.tar.gz](round4-measurement-records-v2.tar.gz)
adds 771 completed new/changed files (198,900,043 original bytes), with
480 SHA256 blobs; 16 rebuildable binaries/classes are excluded with hashes.
Active 5de full-suite logs and upcoming timed trials are excluded.
SHA256: `11f0f2071d966c3b65d62efae6db2b83614c9a2162de9ad48dbf30f5b3956375`.
Use the same verifier as v1, changing the archive filename. Both cuts were
fully verified; the benchmark tooling now passes 22 lightweight tests.

[round4-source-manifest-v2.json](round4-source-manifest-v2.json) records
eleven additional frozen source snapshots in
[round4-prototype-sources-v2.bundle](round4-prototype-sources-v2.bundle).
The manifest labels source-only hypotheses and validated preparation
separately; recovery prerequisites 2262bb6a, 289058de and 5de6e335 are all
in PR 23 history. Earlier archives and source bundles remain unchanged.

## Third completed evidence cut

[round4-measurement-records-v3.tar.gz](round4-measurement-records-v3.tar.gz)
adds 719 completed new/changed files (364,561,319 original bytes), deduplicated
to 240 blobs. It contains every sample and artifact from the completed 5de
screen and three-cohort warm-JVM block, the genuine Closure protobuf
ADVANCED runtime gates, and the latest source/profile audits. One rebuildable
class file is excluded with its hash. Active full-suite and new generated
fixture attempts are excluded from this cut.

Archive SHA256: `02fe85d92d23deae0e4bb01a92e6a8afc994de886c623bc4dc32302b7315db9f`.
The same verifier validates all paths and payload hashes. All 28 benchmark
tool checks pass on the publication tree.

[round4-source-manifest-v3.json](round4-source-manifest-v3.json) records 14
frozen snapshots in [round4-prototype-sources-v3.bundle](round4-prototype-sources-v3.bundle),
including the earlier states, new post-callback DAE admission checks and the
5de COW diagnostic port. Both new compiler variants are **source-only**:
they have not passed Go compilation, normal/race tests, byte gates or timing.
The benchmark queue helper has passed its six Python model checks.

The new audits put historical DAE at only 1.27–1.79% of sampled CPU, and
Java-style string hashing below 1%. Neither cumulative traversal time nor
asynchronous GC can be credited wholesale to one proposed optimization.
An index-linked AST audit records a bounded-region design, identity,
value-copy and strong-edge constraints, and a measurable break-even test;
it is not an implemented layout or an RSS result. Current-main profiles and
actual COW allocation-event counts are the next useful gates before larger
rewrites. Earlier archives remain unchanged.
