# Round 4: current main, verified base and enterprise benchmark protocol

Investigation resumed on 2026-10-20 from main
`2262bb6a8a9d24ee34fdb52bd30c49269e4b5373`. This page records completed baseline measurements and ongoing experiments;
no experimental optimization has been promoted. Earlier rounds retain their original
source, machine, PGO and measurement cohorts; their percentages are not
additional improvements over this main.

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
42 top-level, with no skips/failures; its full unfiltered jscomp suite is
still running. No optimization is promoted based on this screen.

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
The earlier madvise percentage does not appear as a leading flat symbol
in this Linux cohort; fresh syscall attribution is still outstanding.

Read-only disassembly also rejects two tempting duplicate optimizations:
ordinary string maps already inline to Go fast-string operations, and PGO
already guards and inlines ReferenceCollector.Visit in the main traversal.
Current Node headers are 80 bytes. A bounded census of actual property,
string-cell and extension identities found weak sharing. Even ideal 24-byte
shared cold records save only 2.96% of finalized bundle structural bytes and
1.95% of typed-full structural bytes; 32-byte records increase bytes. These
are structural projections, not RSS or speed gains. Added objects, indirection
and copy-on-write costs make a wholesale AST64 rewrite unattractive on this
evidence. A narrower property/string-only census remains open.

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
