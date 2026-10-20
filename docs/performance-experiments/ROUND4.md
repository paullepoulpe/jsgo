# Round 4: verified baselines and enterprise benchmark protocol

The validated integration now follows main`abcf69a0` at`42aaa853`. Its full
normal suite passes570 effective files/3885 top-level roots (3848PASS/37original
SKIP), all28 new regressions,21 packages and the exact original nested inventory.
BundleN1/N2 and typedN1 CLI gates match all6 Java digests, with default runtime
settings, actual4GiB caps, zero OOM deltas and stopped-process proof.

A fresh18-row balanced Go-only screen measures the new integration against
previously validated`1fe9acbb`: median same-round elapsed ratios are0.965857
(typed;3.41% less) and0.967770 (bundle;3.22% less). RSS ratios are0.998205
(typed;flat) and0.920322 (bundle;7.97% less). Three samples/build/task make this
a screen, with raw SD/ranges retained, not a confidence or fresh-JVM result.
This comparison includes changed main code and shipped PGO together.

The completed fresh typed warm-JVM results immediately below belong to previous
main`bb521c0f` plus fixes`1fe9acbb`, notabcf. No historical Java times are combined
with the new Go-only screen. Experimental redesign factors remain unmeasured.

The completed post-screen sync was `bb521c0f26b2fdeec9622247531174826e9fd23c`.
The confirmed warm-JVM bundle win below belongs to frozen `5de6e335`.
Both stock main and the integrated compatibility fixes pass bundle N1/N2
and typed N1 CLI byte gates. The fixes also pass the complete normal suite.
The latest-main Go-only screen, profiles and fresh three-JVM typed comparison
are complete. Fixed Go takes **9.74% less elapsed time** and **53.43% less
peak process RSS** than warm Java on this typed workload; all30 pairs match
Java's bytes and favor Go. Experimental redesign factors remain unmeasured.

The compatibility negative control uses unchanged `bb521c0f` production
code, all original tests, and 17 new public-API regression roots. Fourteen
roots reproduce the intended regressions; three positive controls pass.
Three further isolated scope-handler subcases also reproduce their expected
failures. All 20 dispositions match, with no fixture/compilation failures or
OOM-counter increase. The failures cover retained scope-creator and variable-
expert ownership, custom scope conversion and variable-filter callbacks,
overridden removable callbacks, custom implicit-property protection, and
redeclaration-handler visibility/comparability.

The integrated fixes are frozen at `1fe9acbb7d3bee4ae4d31849c5a2882cfb3c86ff`.
The complete normal suite passed: **29,535 PASS events / 3,820 top-level**,
103 skipped events / 37 top-level, zero failures across 21 packages. All
3,857 expected top-level names/statuses match; all 18 new regression roots
pass. Every original nested test name/status event matches the completed
5de run exactly, with 27 new passing events added. All 559 effective source
hashes match the frozen inventory, with only corpus-path test overlays and
no production overlay, filter, `-short`, cached result or runtime GC tuning.
There were no pauses or OOM-counter increases. Full live-Java corpus parity
and the unfiltered race suite remain outstanding.

The separately built normal fixed CLI now passes all six stock-Java digests
on bundle N1/N2 and typed N1, with verified 4 GiB caps and zero OOM deltas.
Its SHA256 is `052ad6581befb4c011ba5e1d9150e0ab61451ee032ab47fc456a6577ff308df9`.
The four compatibility fixes are retained as separate PR commits. Every one
of the PR's 1,037 production Go sources and 559 effective test files matches
the validated frozen source byte for byte, including the 12 corpus-path
test overlays; no production differences are hidden by the cherry-picks.

The subsequent screen retains stock main and the fixed build as separate
variants so that compatibility costs remain visible. The byte-check clocks
are diagnostic only, not a fresh timing comparison.

An isolated indexed-AST model is frozen at
`0ec2cc95065cc012eb949d12360a48f31d822649`. All nine focused roots pass in
normal, race and package-only `checkptr=2` builds, with no skips or OOMs.
The compiled layout is 56 bytes, and tests cover copied/foreign/overwritten
identity, metadata aliases, capacity fallback and GC lifetimes. This does
**not** change Rhino or the CLI and establishes no compiler speed/RSS gain.
Its complete-owner ledger exposes retained native chunks, historical type
graphs, pointer registries, fallback snapshots and unknown allocator costs.
Owner-wide history and additional getter/registration work may erase the
header saving; actual latest-main retained/allocation evidence is required
before any compiler integration.

### Next main integration: abcf69a0 / 42aaa853

Main advanced after the completed typed timing block. Stock `abcf69a0` builds
with shipped PGO SHA256`4cdc6cd7b1a056ff180beba14cd600f9afc03dca0df1d3ffe5de5f7d99ecedd4`
and passes all six Java digests on bundle N1/N2 and typed N1, with default
runtime settings, actual4GiB caps, zero OOM deltas and stopped-process proof.
Its binary SHA is`9376ee254b2060b589cf97132f0a893bbca91ca6d27f8e5f2844d4566892c093`.
These are correctness runs; neither their clocks nor upstream macOS percentages
are local performance measurements. The new PGO records Darwin/arm64 Go1.27.1
paths but lacks an exact collection-task/weight/revision manifest. Code and
shipped PGO must therefore be treated together in the next comparison.

The integration at`ae9008f20f505b457aacece60b1b9063bfe78b36` retains the previous
four compatibility fixes and adds separate corrections for published scope/
creator/cache lifetimes, owned CFGs and analysis annotations, repeatable public
printer mappings, self-contained source-map names, mutable namespace map and
snapshot handoff, owned public collector drafts/equal-Var lookup, and original
export-classification callbacks. Private stock compiler paths retain arena,
CFG/liveness and mapping reuse. Independent source reviews and activation
controls accompany these changes; source review is not runtime validation.

The isolated mapping fix`079a0934` passes its complete53-root sourcemap package,
including three common lifetime/arithmetic controls and an actual40-byte compact
record assertion. The older validated1fe passes52 roots. Pre-fix615 reproduces
fresh and retained replay index panics, while the arithmetic positive passes.
The original working-directory error and administrative disk stops are retained;
the corrected run uses the same binary, and no test was omitted. This validates
the isolated component branch, not the entire new integration. Existing inline16
generator Mapping clearing on Reset predates abcf and is unchanged; the replay
repair covers dynamic records and copies. Public printer ownership is separate.

The final integration inventory is frozen:570 effective test files,3885 expected
top-level roots (3848PASS/37SKIP),21 packages. All559 prior effective files and
3857 prior roots are preserved, with28 new roots added. The first615 pre-fix compiler run reproduced19 intended failures. One new
scope test instead failed because its function fixture lacked an AST parent;
that attempt is retained and excluded as regression evidence. A test-only
correction is frozen at`f5ba292ae7c7b1ce00a6463c5d2cbd4841aa4e5f`; all1037
production files remain byte-identical toae9008. The corrected case is being
rebuilt, with19 unchanged negative test bodies explicitly verified for reuse.
The older1fe positive control then exposed a second test error: traversal
wraps callback panics, while the new test expected the raw sentinel. The test
now asserts the established RuntimeException wrapper/cause contract. Sixteen
positive roots passed before that stop; the remaining root was not executed.
The full unfiltered candidate suite then ran first at`42aaa853`:3848PASS,
37unchanged skips, all28new roots PASS and21packages, with every original nested
name/status event matching1fe exactly and no OOM/still-running processes. This
includes actual private-reuse activation controls. The combined historical
negative/positive proof remains explicitly incomplete: the corrected test version
has not been rerun on the two control compilers. The prior raw dispositions,
fixture failures and launcher errors remain preserved without relabelling them.

The integrated normal CLI SHA256 is
`c1ce317e909b539e780ce16bb798ea233e7632195e7be13daf62c5115f7472fb`.
It passes bundleN1/N2 andtypedN1 strict6-digest gates with actual4GiB caps,
default runtime settings, zero OOM and stopped-process proof. Current PR
production/effective tests match the validated42aaa source exactly. Local PR production/test sources
match the frozen integration exactly (1037 production Go files and570 effective
test files, including the physical counterpart of the overlay-only corpus helper).
No new full-suite or performance claim follows from this equivalence.

The new base changes the open-factor queue. Scope/Var caps are now1024/2048,
whose large backing arrays already occupy exact allocator pages. The old186/341
factor loses its proposed density advantage and would reduce80-chunk retention
by roughly5–6 times; it is deferred pending fresh measured tail evidence. New
source-only Node64/default128, Node64/cap127, Node80/cap127, ordered-map growth
and traversal-layout variants retain byte-identical factor code/tests on the
new integration. Their builds, actual inlining, Java gates and time/RSS results
remain pending. Older frozen variants stay preserved.

Resource housekeeping preserves the evidence: only the rebuildable Go cache was
cleaned after a build disk guard; subsequently closed published CLI outputs were hash-checked before and after
replacing duplicate files with hard links. The latest pass verifies750
JS/map/manifest paths against cuts1–7 and retains23 unique contents, recovering
about1.49GB in that pass (earlier dedup passes are retained separately). Every
path, byte and mode remains unchanged. Logs, binaries, source and active outputs are untouched.
This is disk storage housekeeping, not a compiler memory improvement.

### Fresh42 CPU/allocation profiles

Separate main-only diagnostic build`361a7851...` retains all1036 other production
Go blobs and the same4cdc PGO as normalCLI`c1ce317e...`. TypedN1 andbundleN1
both match all6 Java artifacts under actual4GiB/defaultGC/CPU2. Captures,
12pprof tables and custom decoders are complete/stopped with zero OOM deltas.
The initial analysis-only read-only-cache error is retained; a workspace-cache
runner decodes the same captures successfully. No capture was rerun or selected
because of its timing. Instrumented clocks/RSS are diagnostic, not benchmarks.

| Profile | Sampled CPU seconds | Sampled cumulative alloc MiB | Sampled objects |
|---|---:|---:|---:|
| Typed |34.61|4182.08|25,418,477|
| Bundle |8.90|1358.99|7,795,218|

| Attribution view | Typed CPU | Bundle CPU |
|---|---:|---:|
| GetNext flat |16.32%|7.53%|
| NodeTraversal methods flat |7.11%|9.21%|
| GC-scan stack union |9.77%|10.00%|
| madvise |0.087%|0%|
| futex |0%|0%|
| Other syscall stack union |0.231%|0.112%|

The disjoint-priority tables, overlapping unions/intersections and exact source
owners are retained in the profile analysis. These views are not a sum of
removable costs. GetNext is already one inlined pointer load; no hardware-counter
cache-miss attribution was performed. Traversal-inclusive66% includes actual
pass callback work; callback-source-line flat1.24%/1.80% is not pure dispatch.

Typed allocation sites include AST/IR slab191.32MiB, scopeArena.allocScope187.86,
SourcePosition slab170.77, ArrayList.Add156.83 and string/pointer ordered-map
addSlow149.15. Cumulative sampled bytes are about28% below the earlier1fe profile,
but live retention/GC pacing and separate captures prevent treating that as an
RSS or elapsed result. Scope-tail retention needs a census before changing the
80-chunk policy. Pointer locality is a current experimental lead; the physical
Node64 factor still has no measured time/RSS result at this cut.

### abcf integration: completed Go-only screen

All18 rows (three balanced rounds per workload across oldfixed1fe, stockabcf
and fixed42aaa) match the same pinned Java bytes. Each variant occupies each
order position once per task. No builds, profiles, archives or heavy audits
overlap the block. Every row uses fresh CLI/defaultGC/actual4GiB/CPU2; all
inner/outer OOM deltas are zero and all owned runtimes/containers are stopped.

| Workload / source | Median seconds (sample SD) | Median peak RSS MiB (sample SD) |
|---|---:|---:|
| Typed / old fixed1fe |32.877 (0.328)|1097.15 (47.19)|
| Typed / stockabcf |31.222 (0.822)|1145.75 (17.75)|
| Typed / fixed42aaa |31.647 (0.319)|1093.79 (45.43)|
| Bundle / old fixed1fe |8.650 (0.030)|390.34 (5.15)|
| Bundle / stockabcf |8.385 (0.221)|381.82 (11.74)|
| Bundle / fixed42aaa |8.413 (0.093)|359.24 (11.32)|

| Workload / numerator ÷ denominator | Median same-round time ratio | Median same-round RSS ratio |
|---|---:|---:|
| Typed / stockabcf ÷ oldfixed1fe |0.958561|1.070808|
| Typed / fixed42aaa ÷ oldfixed1fe |0.965857|0.998205|
| Typed / fixed42aaa ÷ stockabcf |0.999007|0.957221|
| Bundle / stockabcf ÷ oldfixed1fe |0.969391|0.992517|
| Bundle / fixed42aaa ÷ oldfixed1fe |0.967770|0.920322|
| Bundle / fixed42aaa ÷ stockabcf |0.998334|0.926913|

Ratios use same-round observations and therefore need not equal ratios of the
separately displayed medians. Raw samples, ranges, variance, CPU data and row
mapping are retained in`main-abcf69a0-go-screen-preparation-v2`. These small
samples show no clear elapsed penalty for the compatibility fixes and lower RSS
in this screen; individual fixes have not been isolated as causes. Main code
and shipped PGO changed together. There are no Java timing observations or
confidence intervals in this screen; fresh Java must be measured separately.

### Previous bb main Go-only screen

The complete 18-row screen compares frozen `5de6e335`, stock `bb521c0f`
and fixed `1fe9acbb` in three balanced rounds on each workload. Each binary
occupies each order position once. All rows match the pinned Java artifacts,
use default Go runtime settings and verified 4 GiB caps, and have zero inner/
outer OOM deltas. No builds, profiles or packaging run during this block.

| Workload / source | Median s (sample SD) | Median peak RSS MiB (sample SD) |
|---|---:|---:|
| Bundle / 5de | 9.514 (0.117) | 371.43 (4.16) |
| Bundle / bb | 8.645 (0.245) | 370.76 (7.62) |
| Bundle / fixed | 8.763 (0.084) | 386.02 (5.48) |
| Typed / 5de | 34.823 (0.406) | 1136.91 (12.25) |
| Typed / bb | 32.578 (0.321) | 1108.98 (5.17) |
| Typed / fixed | 32.993 (0.553) | 1129.97 (30.03) |

| Workload / comparison | Median same-round time ratio | Median same-round RSS ratio |
|---|---:|---:|
| Bundle / bb ÷ 5de | 0.903739 | 1.002459 |
| Typed / bb ÷ 5de | 0.936104 | 0.980206 |
| Bundle / fixed ÷ bb | 1.010689 | 1.032162 |
| Typed / fixed ÷ bb | 1.012752 | 1.013961 |

Latest main takes 9.6% less bundle time and 6.4% less typed time in this
screen. Its RSS is effectively flat on bundle and about 2.0% lower on typed.
The compatibility corrections add roughly 1.1–1.3% time and 1.4–3.2% RSS in
these samples. Three observations per variant make this a screening result,
not a settled cost estimate or a confidence interval. All ranges, CPU data,
variance and per-round mappings remain in `main-bb-api-fixes-v1/go-screen-*`.
No historical Java time is combined with these measurements.

### Completed latest-main typed warm-JVM comparison

This is a fresh same-session comparison of stock `bb521c0f` and the validated
compatibility-fixed `1fe9acbb`, sharing each cohort's stock warm JVM. Inputs
are 1,820 files / 25,919,443 bytes, with 1,747 emitted input identities. Both
use compiler N1, two CPUs, real 4 GiB limits and zero swap. Go is a fresh CLI
with default runtime settings; Java uses `-Xmx3g`, retaining JVM/JIT warmth
while creating a fresh stock compiler and reparsing externs every request.
No worker, AST cache or historical Java timings enter these results.

The final wrapper admission passes: 30 measured pairs, 53 Java warmups,
18 fresh-Go warmups, all six artifact/stream/exit digests matching, zero OOM
deltas, and no remaining compiler processes or containers. No build, profile,
archive or heavy source audit overlaps the timed block. First fresh Java
oracle admission precedes every measured pair. The final three-cohort results
are below; absolute aggregates are geometric means of cohort medians, not
pooled medians.

Three restarted JVMs; warmups 18, 16, 19; 30 alternating pairs; all bytes match.

| Cohort | Go variant | Java s (SD) | Go s (SD) | Java RSS MiB | Go RSS MiB (SD) | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|---:|---:|
| 1 | stock-bb521c0 | 33.160 (0.330) | 30.957 (0.679) | 2418.34 | 1091.57 (28.02) | 1.07118 | 2.21547 |
| 1 | fixed-bb-1fe9acbb | 33.688 (0.389) | 30.452 (0.702) | 2418.34 | 1113.22 (6.33) | 1.10629 | 2.17238 |
| 2 | stock-bb521c0 | 33.882 (0.573) | 30.547 (1.243) | 2406.15 | 1113.09 (39.52) | 1.10920 | 2.16168 |
| 2 | fixed-bb-1fe9acbb | 33.498 (0.283) | 30.203 (0.367) | 2406.15 | 1122.17 (9.95) | 1.10906 | 2.14420 |
| 3 | stock-bb521c0 | 34.564 (0.615) | 30.943 (1.277) | 2376.02 | 1112.50 (11.54) | 1.11703 | 2.13575 |
| 3 | fixed-bb-1fe9acbb | 34.338 (0.511) | 30.982 (0.274) | 2376.02 | 1117.68 (27.18) | 1.10831 | 2.12585 |

| Go variant | Java/Go time (95% CI) | Java/Go RSS (95% CI) | Java / Go absolute s | Java / Go absolute RSS MiB |
|---|---:|---:|---:|---:|
| stock-bb521c0 | 1.09895 (1.07118–1.11703) | 2.17071 (2.13575–2.21547) | 33.864 / 30.815 | 2400.10 / 1105.67 |
| fixed-bb-1fe9acbb | 1.10789 (1.10629–1.10906) | 2.14739 (2.12585–2.17238) | 33.839 / 30.544 | 2400.10 / 1117.68 |

| Variant / runtime | Pooled elapsed median (range; SD) s | Pooled RSS median (range; SD) MiB |
|---|---:|---:|
| stock-bb521c0 / java | 33.726 (32.823–35.173; 0.745) | 2406.15 (2296.89–2418.34; 40.74) |
| stock-bb521c0 / go | 30.719 (30.111–33.435; 1.023) | 1106.74 (1070.38–1179.68; 27.72) |
| fixed-bb-1fe9acbb / java | 33.688 (33.096–34.650; 0.453) | 2406.15 (2296.89–2418.34; 40.74) |
| fixed-bb-1fe9acbb / go | 30.519 (30.028–31.915; 0.505) | 1117.68 (1101.54–1177.92; 17.86) |

RSS windows: Java common final warm-JVM lifetime VmHWM including warmup/JIT; Go median fresh-process lifetime wait4 peak. Not per-compilation live heap.

One typed workload, N1 compiler work, two CPUs and actual 4GiB caps. Three restarted cohorts underpin the CI; 30 requests are not 30 independent JVMs.

Outer max events are retained as resource-pressure diagnostics, distinct from OOM. No performance cause is inferred from the event counter alone; caps and heap settings were not changed.

Bootstrap: Take natural logs of3 cohort ratios; draw3 logs with replacement for each of10,000 bootstrap replicates; exponentiate their arithmetic mean; sort the replicates. Nominal95% percentile interval: sorted bootstrap elements250 and9749 (zero-based), as in the frozen helper.

All15 matched requests per runtime/variant supply descriptive min/max/median/mean/sampleSD/samplevariance only; they are not treated as15 independent JVMs for confidence.

Aggregate ratios, restarted-cohort bootstrap intervals, ranges, variance and direct Go-control effects are retained in summary.json and the frozen runner output.

The fixed build saves 9.738% elapsed time and 53.432% RSS against Java; stock
main saves 9.004% elapsed time and 53.932% RSS. The direct stock-Go/fixed-Go
time ratio is 1.00886 (95% cohort interval0.99873–1.01659), so this block does
not establish that the compatibility fixes speed up the compiler. Its direct
RSS ratio is0.98926 (0.98055–0.99537), about1.1% more RSS for fixed Go. These
are separate from the earlier three-observation screen; neither direction is
silently substituted. The performance lead includes optimizations landed on
main, not gains from the new unrun source candidates.

Raw requests, individual timings/CPU/RSS, warmup decisions, container records,
helper/source/runtime pins, first-oracle admission and shutdown proof are under
`main-bb-api-fixes-v1/paired-java/typed-n1/` and
`main-bb-api-fixes-v1/typed-warm-wrapper-v2/`. Full external Java corpus parity,
unfiltered compiler race and new enterprise warm timings remain open.

### Fresh fixed-main profile

Separate instrumented bundle/typed invocations also pass all six Java digests.
Typed samples attribute 13.40% flat CPU to traversal methods, 8.99% to the
inlined `GetNext` pointer load and a 14.16% union to GC scan stacks. These
are different, non-additive attribution views; none supplies a predicted
optimization gain. Inclusive traversal is 62.59% because it includes callback
work. Samples at callback-call source lines are 1.28%, not an isolated
interface-dispatch measurement. `madvise`, futex and Syscall6 account for
0.03%, 0% and 0.08% respectively, so the old syscall targets remain small.

Typed allocation samples total 5788.92 MiB / 29.16 million objects, including
382.12 MiB at scope allocation, 346.69 MiB in string/pointer LinkedHashMap
`addSlow`, 195.85 MiB at Var allocation, and 190.09 MiB at CFG node allocation.
Heap samples are not RSS, final live ownership or a scan-by-type breakdown.
Raw profiles, exact caller/owner attribution and the preserved wrapper setup
failure/retry are under `main-bb-fixed-profile-v1/`; diagnostic clocks are
not performance measurements. New isolated scope, traversal and container
pilots remain source-only after the completed typed JVM comparison.

### Isolated next candidates: source only

None of these candidates has been compiled, tested or timed yet. They are
separate branches based on the validated fixed compiler, rather than changes
to the active benchmark binaries. Independent source reviews establish no
performance gain.

| Candidate | Frozen source | Proposed saving and material risk |
|---|---|---|
| Large ordered-map growth | `269494d6` | Reserve twice the capacity at large growth boundaries to reduce repeated copies; extra slack may increase RSS and GC work. |
| Traversal ordinary path first | `58dd0503` | Change branch layout without changing callback order; discard if ordinary PGO assembly is already equivalent. |
| Scope/Var slab page fit | `fb32e356` | Caps 186/341 target 16 KiB backing allocations; allocation count increases, and rollover sizes can reserve more memory than stock. The earlier 340-slot source `7e607d7c` is retained. |
| Physical Node64 metadata union | `99f9fc7b` | Fold three cold metadata pointers into one tagged, GC-tracked immutable representation while preserving the hot links. Extra tuple allocations and getter work may erase the header saving. |

Node64 is a real Rhino prototype, distinct from the isolated indexed model.
Its intended 64-byte layout is still unverified by a build. It preserves
full token/type/source metadata and copied-Node field replacement semantics;
it changes exactly two original sizeof expectations from 80 to 64, leaving
their behavioral and GC assertions intact. Its seven new roots include five
stock-compatible semantic controls. Simultaneous disjoint metadata writes
to the same Node would now touch one shared representation; source review
found no such supported compiler path, but independent-clone tests do not
prove that case safe. Whole-compiler race validation remains necessary.

A separate 127-node slab experiment is being prepared for both Node80 and
Node64. This separates allocator packing from representation changes. The
size-class arithmetic is a hypothesis, not measured RSS or a speedup. Actual
factory/JSDoc inlining, parser/GC tests and strict Java byte gates precede
any timed comparison. All variants retain normal runtime GC settings.

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
artifacts under the actual 4 GiB cap. Its full normal suite has passed.
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

The completed full normal run passes **29,508 checks** (3,802 top-level),
with 103 skipped events (37 top-level) and zero failures across 21 packages.
All 3,839 original top-level names/statuses match the actual completed 289
run exactly. All 553 effective test-source files remain present; the only
12 test overlays change corpus paths. No production overlay, `-short`,
`-run` or cached test result was used. Runtime GC and CPU settings are default.
Full live-Java CLI/corpus parity and unfiltered race coverage remain separate,
outstanding gates.

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
Soy and the modern J2CL JRE fixture now also pass all Java/Go bytes and
independent execution of both compiled outputs. The first 2022 JRE fixture
was incompatible with the 2026 Java compiler's registration ABI; the new
fixture uses the unmodified published Kie J2CL packaging fork's May 2026
JRE/bootstrap pair, with matching ABI and verified checksums. This is a JRE
fixture, not a large production Java application. Soy includes the genuine
same-release `soy.checks` module missing from the first failed preparation.
All setup failures remain retained.

The 32-chunk ADVANCED gate passes all 103 digest entries: 98 output artifacts,
two expected absences, stdout/stderr/exit. Both engines use the same 1,820
provided files, emitting 1,747 source-manifest entries exactly once. The 68
`@externs` and five `@nocompile` files account for the difference, matching
the pinned Java AST inventory. All 32 named chunks and the weak dependency
record match. These correctness checks ran alongside the normal suite;
their clocks and diagnostic memory readings are not performance comparisons.

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

## Fourth completed evidence cut

[round4-measurement-records-v4.tar.gz](round4-measurement-records-v4.tar.gz)
adds 778 completed new/changed records (239,342,479 original bytes), with
473 content-addressed blobs. Eight rebuildable binaries/classes are excluded
with hashes. It preserves the complete 5de normal run and inventory check,
fresh byte-gated profiles, modern J2CL/Soy/32-chunk/typed correctness gates,
and the COW mutation diagnostic, including its blocked first Docker attempt.
Input symlinks are metadata only. Active optimization builds are excluded.

Archive SHA256: `2728b3efa9c3617129956e2d535f1feb0f26bf6ab9a78538ac71da025611efe1`.
All payloads passed the same verifier. The publication tree's 29 benchmark
tool checks pass.

The current profiles show GC-scan stack unions at 17.10% bundle / 15.75%
typed, and `madvise` at zero / 0.04 sampled seconds. Typed source-map record
allocation is now 39.69 MiB versus 2,377.65 MiB in the earlier profile.
These are single diagnostic profiles, not paired speed measurements.
Scope/Var/slot arenas still allocate about 754.61 MiB typed / 388.66 MiB
bundle. Their observed requests and lifetimes need a census before reuse or
growth-policy changes; allocation bytes alone do not establish retained RSS.

The naive COW72 diagnostic passes 18 focused Rhino tests and both Java
artifact gates, with four complete population and mutation snapshots per
task. It models 1,342,450 bundle / 3,929,488 typed wrapper allocations,
or 20.48 / 59.96 MiB, versus favorable full-slab savings of at most about
7.19 / 25.08 MiB in sampled node-slab allocations. This argues against a
naive per-store pair rewrite; it is **not** measured performance/RSS for an
implemented smaller Node. A direct-pointer union could avoid allocations
when only one metadata field exists. Its separate counter extension is
source-only; the old counters cannot decide that question.

[round4-source-manifest-v4.json](round4-source-manifest-v4.json) fingerprints
20 frozen states in [round4-prototype-sources-v4.bundle](round4-prototype-sources-v4.bundle).
New Sass value-result, node-slab-capacity and combined candidates remain
unpromoted with validation/timing gates pending. The direct-union diagnostic
is unbuilt. A complete genuine BigQuery API family has 227 generated JS
files from 48 root schemas and 13 dependencies; raw JavaScript and native
protobuf payload checks pass, while compiler gates and actual manifest
counts remain pending at this cut. Its source/generator pins and semantic
checker are included in the benchmark tools. An automatic approval rejection
of the concurrent compiler launch is retained; no compiler ran, and that
queue is deferred until the active build finishes. Earlier cuts are unchanged.

## Results completed after the fourth evidence cut

The BigQuery compiler workflow now passes. Both stock Java and main
`5de6e335` produce identical bytes for all six required artifact/stream/exit
digests. Both compiled programs execute independently, and native `protoc`
decodes their actual Table and Job payloads into the complete expected
messages. This verifies real generated-code behavior as well as compiler
output identity; these correctness clocks are not warm-JVM measurements.

The exact manifests select 251 of 1,209 supplied JavaScript candidates:
199 generated sources, nine protobuf runtime sources, 42 Closure Library
sources and one entry. All 194 BigQuery family providers remain selected.
The 958 omitted candidates are normal dependency pruning in both compilers;
the benchmark does not claim all supplied files become compiled output.

The first concurrent launch was rejected by automatic approval review
because a large Go build left insufficient assured outer-memory headroom.
No compiler child ran. After that build finished, a safer sequential launch
passed fresh total-headroom and no-other-workload checks before each step,
under the same actual 4 GiB cap. The reviewer accepted that launch; all
seven steps passed with no outer OOM-counter increase. The original
rejection and the successful serial resource proofs are retained separately.

Sass value-result and 102-node-slab candidates are still experimental.
Their combined production source is frozen at `8fefd1cb`; later commits
`bc88a93c` and `9d9d70b2` correct test fixtures only. The first new parsing
fixture incorrectly constructed partial parser objects; another comparison
ran type evaluation between raw AST snapshots even though type evaluation
mutates shared implicit-bound metadata. Both failed attempts are retained.
The final fixtures use complete constructors and compare both raw ASTs
before type evaluation. All eight selected parsing targets are accounted
for with the slab off/on: five pass, including the existing parsing/JSDoc
goldens, and three unchanged optional external-oracle targets explicitly
skip. Full-source jutil tests and CSS semantic tests also pass. These are
focused checks, not an experimental-candidate full-suite claim.

Actual candidate assembly shows the new slab allocation method did not
inline at a hot node allocation callsite, unlike stock main. The planned
five-variant screen therefore includes both stock main and the combined
binary with both factors disabled, so default-off code-generation cost is
visible. A tiny isolated Go compiler probe found that the explicit
`noinline` growth helper pushed the allocation method to cost 87, above the
ordinary budget of 80. Removing the pragma and using equivalent builtin
min/max clamps reduces its probe cost to 59 (stock `New`: 61). The separate
refinement is frozen at `4e3e56b2`; actual whole-binary inlining and its own
tests/artifact gates remain required. The earlier candidate is preserved
rather than silently replaced. No performance improvement is claimed for
either candidate yet.

The frozen original candidate passes all eight bundle/typed Java artifact
gates: Sass off/on crossed with slab off/on. A same-diagnostic-binary Sass
off/on bundle pair also passes every byte gate, with the slab disabled in
both. The default sampled allocation profile attributes 425,987 objects /
13,631,904 bytes to the old expert constructor when off and has zero samples
at that site when on. This demonstrates the intended allocation-site
change; it is not an exact allocation counter, total-heap reduction,
peak-RSS improvement or elapsed-time result. The later refinement changes
only the equivalent slab growth helper, not the Sass implementation.

## Five-variant allocation-candidate screen

Source `4e3e56b2`, normal CLI SHA256
`a2964efec241cac2d11ae013e1a17964f63436691dce18936d632e42de130d44`,
passes all 165 jutil tests, the focused full-source parsing checks off/on,
and all eight Java artifact gates. Actual assembly confirms inline allocation
at all nine factory and four JSDoc sites; capacity loads occur on rollover only.
The candidate still has no full normal/race validation or paired-JVM result.

The quiet Go-only screen runs five variants in three balanced rounds on each
workload, under the actual 4 GiB cap and default runtime GC/CPU policy. All
30 fresh invocations match the strict stock-Java artifact oracle for identical
argv, inputs and externs. This is exploratory evidence from three rows per
variant, not a settled-JVM comparison or confidence interval.

Each table reports median (sample SD); ratios are medians of same-round
variant/stock ratios, rather than ratios of independently selected medians.
All ranges, variance, individual rows, order and the alternative default-off
control summaries are retained in `sass-slab-inline-go-screen-v1/`.

### bundle N1

| Variant | Seconds (SD) | Peak RSS MiB (SD) | Time / stock | RSS / stock |
|---|---:|---:|---:|---:|
| stock-5de | 9.295 (0.015) | 389.18 (13.50) | 1.00000 | 1.00000 |
| integrated-off | 9.601 (0.237) | 373.78 (11.60) | 1.03213 | 1.00427 |
| sass-only | 9.524 (0.177) | 386.31 (5.74) | 1.02385 | 0.99261 |
| slab-only | 9.238 (0.145) | 379.55 (5.12) | 0.99628 | 0.95399 |
| sass-and-slab | 9.186 (0.055) | 372.29 (1.20) | 0.99068 | 0.95292 |

### typed N1

| Variant | Seconds (SD) | Peak RSS MiB (SD) | Time / stock | RSS / stock |
|---|---:|---:|---:|---:|
| stock-5de | 34.022 (0.157) | 1118.26 (10.30) | 1.00000 | 1.00000 |
| integrated-off | 33.096 (0.508) | 1134.68 (13.71) | 0.97913 | 1.01469 |
| sass-only | 34.140 (0.568) | 1131.23 (18.76) | 1.01000 | 1.01160 |
| slab-only | 32.473 (0.675) | 1117.79 (15.18) | 0.95446 | 1.00804 |
| sass-and-slab | 32.744 (0.898) | 1102.77 (22.73) | 0.96871 | 0.98615 |

Slab-only is the better typed lead: median same-round time/stock is 0.95446,
but time/default-off is 0.98298. Default-off itself is 2.1% faster than stock
on typed input and 3.2% slower on bundle, showing why the control is required.
Typed slab RSS/stock is 1.00804, so this does not establish a memory saving.
Combined typed time/stock is 0.96871, weaker than slab-only. The combined
bundle ratios are 0.99068 time and 0.95292 RSS; the RSS spread is substantial.
Sass removes its targeted allocation-site samples but does not demonstrate a
macro speed win here. No factor is promoted from this small screen.

Matched PGO retraining remains an open discriminator. Changed concrete callback
identities and relative call-site lines can change Go PGO decisions in either
direction. There is no existing proof that stock PGO devirtualized the Sass
callback or that PGO caused this result. Before training, the old `pgo.py`
needs fail-closed workload/exit/profile checks; skipped training must never be
silently accepted. The archived audit specifies stock/default-off/enabled
training controls. No such retraining was performed in this screen.

The test-only corpus-root correction is now checked into the PR. Its 12 source
files match every byte of the overlays used for the completed full main5de
normal run; production compiler/library source remains stock main. Use
`JSGO_ORACLE_CORPUS=/path/to/the/pinned/corpus` to request that coverage; an
explicit missing/empty/non-directory root fails instead of silently reducing
coverage. Full external Java corpus and unfiltered race parity remain open.

## Fifth completed evidence cut

[round4-measurement-records-v5.tar.gz](round4-measurement-records-v5.tar.gz)
adds 2,142 new/changed closed records (720,724,211 original
bytes), 1,045 content-addressed blobs, and 11 excluded
rebuildable binaries/classes with hashes. Archive SHA256:
`a29fcd48d4c8ef3c770e4cf4514e7e9e613a6658415553c6b5ff40ee8c14645e`.
The verifier passes every payload hash. This preserves the successful serial
BigQuery workflow, original/refined candidate builds and tests, failed
fixtures and launch attempts, byte gates, allocation-site/inlining evidence,
all30 exploratory rows, numeric/source-only design audits and PGO caveats.

[round4-source-manifest-v5.json](round4-source-manifest-v5.json) fingerprints
24 source states in [round4-prototype-sources-v5.bundle](round4-prototype-sources-v5.bundle)
(SHA256 `4bd4386e6762d9970c9939f626a9d0a2ee63f332dcd83383dc228a434a034fcd`). Earlier cuts remain unchanged.

At the post-screen fetch, main advanced to
`bb521c0f26b2fdeec9622247531174826e9fd23c` (23 changed files, new PGO).
All completed performance numbers above remain measurements of their named
frozen sources, including the5de warm-JVM bundle win. This new main has not
yet been measured by these experiments. A fresh-base validation and
compatibility review precede another long warm-JVM block; the commit title's
claimed gains are not imported as our measured results.

## Sixth completed evidence cut

[round4-measurement-records-v6.tar.gz](round4-measurement-records-v6.tar.gz)
adds 288 closed records (111,835,792 source bytes), 213 content-addressed
blobs, and five excluded rebuildable artifacts with hashes. Archive SHA256:
`4050cfca6f4bd05fdf836dda4cbd43fbfbff833ed8fe14522abe9c6ec33f6cce`.
Every payload hash passes verification. This cut preserves stock-main byte
gates, all 20 stock negative/control dispositions, complete fixed normal
inventory and raw logs, fixed CLI byte gates, independent source reviews,
source-equivalence proof, rebuildable-cache cleanup, isolated model normal/
race/checkptr validation, and prepared fresh benchmark/profile commands.
Prepared commands establish no timing result.

[round4-source-manifest-v6.json](round4-source-manifest-v6.json) fingerprints
30 frozen source states in [round4-prototype-sources-v6.bundle](round4-prototype-sources-v6.bundle)
(SHA256 `4e080067df3bd4051be09046ad4981687b4a9729102dee0344da1453024b1bd4`).
The four compatibility corrections are now separate PR commits; experimental
performance factors and the isolated model remain unpromoted. All earlier
archives remain unchanged.

## Seventh completed evidence cut

[round4-measurement-records-v7.tar.gz](round4-measurement-records-v7.tar.gz)
adds1,251 closed records (1,339,343,355 original bytes),499 unique blobs
(50,423,827 bytes) and one excluded rebuildable diagnostic binary with its hash.
Archive SHA256:
`f1b631e27c7c519362d7e09546c5f793a1bcabeb186272fe8c63871a997fb67e`.
Every payload hash passes verification. This cut preserves the18-row latest-
main screen, complete typed warm-JVM comparison and bootstrap/variance review,
fresh profiles and preserved setup failure, unrun source candidates and peer
reviews, enterprise benchmark preparations, validation plans and environment
checkpoints. All prior cuts remain unchanged.

[round4-source-manifest-v7.json](round4-source-manifest-v7.json) fingerprints
38 recoverable source states in
[round4-prototype-sources-v7.bundle](round4-prototype-sources-v7.bundle),
SHA256`6f14665a274064313d3bf727f313e75ad3b4a20b46fabed275a5ba2114860425`.
No experimental factor is promoted.

A new fetch after the timing block finds main`abcf69a0`, a further32-file
optimization commit with new PGO, inter-pass caches and reusable analysis/pass
storage. Synchronization and source review are next. All completed timing and
full-normal results on this page remain attached to their named frozen source;
they do not validate or measure this new main or import its commit-title gains.


## Eighth completed evidence cut

[round4-measurement-records-v8.tar.gz](round4-measurement-records-v8.tar.gz)
adds1,271 new/changed closed records (424,833,516 original bytes),858 unique
blobs (97,239,454 bytes) and9 excluded rebuildable binaries/classes with hashes.
Archive SHA256:
`db0c601ef14aa9651d74ef40c161a011979234c306d5d5081102f54959c6b153`.
Every payload hash passes verification. Earlier cuts remain unchanged. This
cut includes the latest main merge/fixes, full570-file normal inventory,
3CLI byte gates,18-row Go-only screen with RSS/variance, fresh diagnostic
profiles, native string-map assembly audit, full live-parity/enterprise
preparations, source-equivalence/reclamation records, and all retained fixture,
launcher, analysis and administrative failures.

[round4-source-manifest-v8.json](round4-source-manifest-v8.json) fingerprints
64 recoverable source states in
[round4-prototype-sources-v8.bundle](round4-prototype-sources-v8.bundle),
SHA256`534a130efd6f0f0f3f0e4d42b96a33f34c83c66a383742c2c9acd162a9c06263`.

The admission is deliberately a partial review cut: completed currentnormal,
CLI and Go-only measurements pass, while corrected-source combined control
rebuilds, broader live-Java/race parity, fresh current JVM timings and factor
macro results remain open. Earlier20 negative/16positive control results retain
theirf5ba source identity; they are not relabelled42aaa. At cut time the separate
Node64 controller was prepared; its actual experiment starts after this cut.

Java21.0.12.1 is the measured JVM. Artifact/helper class files require21;
Java25LTS is a useful further stock comparator, but only21 is installed and no
25/27 runtime or download was performed. No result here claims superiority over
all JVM versions or representative workloads. Current generated/chunked
benchmark and complete live parity commands are preserved for actual execution.
