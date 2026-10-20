Draft: transferring Go optimizations to warm Java

Seven independently reviewed one-cohort screens are recorded below; confirmatory effects remain PENDING. This source-only report update runs no compiler/benchmark and changes no published record or stock oracle.

These experiments test five ideas on Closure Compiler v20261015. STOCK is the unchanged pinned Maven jar. CONTROL recompiles the same nine Java source families without algorithm changes. The five individual factors and ALL compile those same families with identical flags and dependency relocation. CONTROL/STOCK measures recompilation effects; candidate/CONTROL isolates source changes. ALL must be measured directly, never estimated by multiplying factor gains.

| Factor | Change | Correctness and memory boundary |
| --- | --- | --- |
| fsiv-lazy | Initialize backward reaching-use analysis only when a candidate reaches the original use query; reuse it within that function. | Full variable map and rejection order unchanged; initialization precedes first successful inline/AST mutation. Omitting an unused analysis can omit its pathological errors, such as iteration-limit failures; universal error-path equivalence remains open. |
| namespace-lazy | Acquire traversal Scope/JSChunk after existing irrelevant-token/name exits, preserving explicit incremental arguments. | Namespace/module/source state and incremental context preserved. Custom AbstractCompiler callbacks or failures may observe different lazy getInput/getChunk timing; universal extension-callback equivalence is not claimed. |
| coalesce-bitset | OR full live-in/out BitSets into interference rows instead of setting each pair. | Upper triangle, diagonal and ordered graph extraction preserved; additional lower-triangle bits unobserved. Destination rows do not alias lattice storage. The ordinary100-variable bound limits potential benefit; sparse cases can do extra work. |
| color-cache | Per-factory identity cache before repeated Color normalization. | Canonical order/IDs and miss/errors preserved; protected externally supplied-map constructor stays uncached. Extra entries retain original colors for the factory lifetime and can increase memory. |
| polyfill-cache | Lazily share the immutable default Polyfills resource table across RUC/Rewrite/Isolate pass instances. | Final Strings/enum entries and Guava immutable collections only; no compiler/AST/extern cache. Custom tables remain independent and failed loads retry. One table retained per defining classloader; dynamic resources/exotic loader anchors outside equivalence claim. |

Original/candidate sources, patches and provenance remain under ../variants/. Their historical source-only records are retained; later build/probe receipts describe subsequent execution separately.

Seven overlays—CONTROL, five factors and ALL—have PASS build receipts, each recording71 generated classes. Compilation is not whole-compiler parity. Focused results are separately recorded in ../results/focused-probes-v1/results.json:

| Probe | Actual recorded runs | Evidence boundary |
| --- | --- | --- |
| fsiv-lazy | PASS STOCK, CONTROL, factor, ALL | Seven activation/output fixtures, output-byte comparison recorded. |
| color-cache | PASS STOCK, CONTROL, factor, ALL | Canonicalization/order/isolation/external-map/error controls. |
| coalesce-bitset | PASS factor, ALL | 1,497 actual-helper cases; candidate helper absent in STOCK. |
| polyfill-cache | PASS factor, ALL | Four controls: concurrent default, ordered query equivalence, separate pass/compiler state, custom-table/failure isolation. |
| namespace-lazy | PASS STOCK, CONTROL, factor, ALL | Separate namespace-probe-v1 receipt records successful compilation/runs and output-byte gate. Custom compiler callback/error boundaries remain open. |

The namespace update is recorded separately in ../results/namespace-probe-v1/results.json; the earlier focused-probes-v1 receipt remains unchanged. These controls do not establish full upstream suites, public-extension/error-path parity or performance. Every benchmark request, including warmup, must match the exact stock six-artifact oracle.

Seven completed BigQuery screens have independent PASS reviews, each from one pair of persistent JVMs and five measured pairs. All warmup/measurement requests passed six artifact digests; completed stages are stopped with zero new OOM. These are exploratory observations with no confidence interval, not confirmed treatment effects. All five factors and ALL now have exploratory screens; bundle validation and three-cohort confirmation remain pending at this cut.

Ratios below are candidate/baseline; lower is less cost. Median elapsed/CPU and allocation ratios divide within-JVM medians. Mean-wall sensitivity divides means and is a separate estimator. Lifetime-HWM RSS compares one final process peak per JVM, including warmup/JIT; it is not per-compilation live heap.

| Candidate / baseline | Elapsed median ratio | Elapsed mean ratio (sensitivity) | CPU median ratio | Lifetime-HWM RSS ratio | Allocated-byte median ratio |
| --- | ---: | ---: | ---: | ---: | ---: |
| CONTROL / STOCK | 1.020467 | 1.026473 | 1.018750 | 1.000494 | 0.994405 |
| ALL / CONTROL | 0.974039 | 0.995964 | 1.008282 | 0.924525 | 0.995388 |
| fsiv-lazy / CONTROL | 1.033120 | 1.035055 | 1.028747 | 0.837353 | 0.985433 |
| namespace-lazy / CONTROL | 0.948240 | 0.986467 | 0.954274 | 0.941453 | 0.999758 |
| color-cache / CONTROL | 0.947556 | 0.961724 | 0.957031 | 0.977960 | 1.006581 |
| coalesce-bitset / CONTROL | 0.993649 | 1.009444 | 0.975855 | 0.954227 | 0.996296 |
| polyfill-cache / CONTROL | 0.993503 | 0.990484 | 1.019190 | 0.896466 | 0.994630 |

ALL’s median time ratio0.974039 weakens to mean0.995964, with median CPU1.008282; this is a weak pilot signal. Namespace’s median0.948240 weakens to mean0.986467. Color’s median0.947556 and mean0.961724 have clearer mean/CPU corroboration in this screen, with allocated bytes1.006581. FSIV is slower in both median/mean elapsed despite lower sampled lifetime HWM. Bitset’s median0.993649 changes to mean1.009444 and same-pair mean1.011285, with substantial within-cohort noise. None supplies a replicated speed claim. CONTROL/STOCK’s1.020467 median is one observation, not a constant penalty to subtract from other comparisons. Separate factor gains cannot be added or multiplied.

The table below gives each JVM’s absolute median and sample SD across its five measured requests. Its lifetime HWM is one final value per JVM, without a per-request HWM SD. Allocations are median MiB allocated per request, not live heap; exact byte values and all current-RSS samples are in completed-screens.json.

| Screen | Variant | Wall median s | Wall sample SD s | CPU median s | Lifetime HWM MiB | Allocated median MiB |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| CONTROL/ STOCK | STOCK | 4.111975 | 0.126675 | 4.800000 | 1233.460938 | 1994.381844 |
| CONTROL/ STOCK | CONTROL | 4.196134 | 0.111013 | 4.890000 | 1234.070312 | 1983.223228 |
| ALL/ CONTROL | CONTROL | 4.151817 | 0.081830 | 4.830000 | 1321.222656 | 1989.261375 |
| ALL/ CONTROL | ALL | 4.044033 | 0.108065 | 4.870000 | 1221.503906 | 1980.087364 |
| fsiv-lazy/ CONTROL | CONTROL | 4.134669 | 0.096786 | 4.870000 | 1579.292969 | 1997.130592 |
| fsiv-lazy/ CONTROL | fsiv-lazy | 4.271611 | 0.123061 | 5.010000 | 1322.425781 | 1968.038239 |
| namespace-lazy/ CONTROL | CONTROL | 4.308322 | 0.195635 | 5.030000 | 1385.031250 | 1979.635468 |
| namespace-lazy/ CONTROL | namespace-lazy | 4.085323 | 0.201763 | 4.800000 | 1303.941406 | 1979.156845 |
| color-cache/ CONTROL | CONTROL | 4.443304 | 0.199208 | 5.120000 | 1295.742188 | 1987.839745 |
| color-cache/ CONTROL | color-cache | 4.210281 | 0.237803 | 4.900000 | 1267.183594 | 2000.920807 |
| coalesce-bitset/ CONTROL | CONTROL | 4.231353 | 0.172440 | 4.970000 | 1298.519531 | 1983.604065 |
| coalesce-bitset/ CONTROL | coalesce-bitset | 4.204478 | 0.201887 | 4.850000 | 1239.082031 | 1976.257095 |
| polyfill-cache/ CONTROL | CONTROL | 4.193790 | 0.086109 | 4.690000 | 1428.425781 | 1985.472839 |
| polyfill-cache/ CONTROL | polyfill-cache | 4.166544 | 0.105407 | 4.780000 | 1280.535156 | 1974.811813 |

Five requests share one adaptive JVM state, so their SD/ratios are descriptive, not five independent treatment replicates. Identical CONTROL medians vary across separate cohorts; do not normalize other screens by the one STOCK/CONTROL observation. Confirmatory holdouts must exclude every screening observation. The existing median estimator remains primary, with predeclared mean-wall/summed-wall sensitivity. Pending team plan: color and ALL directly versus STOCK using three new independent pair-of-JVM cohorts each; a separate ALL bundle screen tests a second workload. These are planned checks, with no admitted outcome. Selection, multiple-testing limits and any later batch must remain explicit.

Polyfill screen is now independently reviewed PASS: exact median elapsed ratio0.993503390650036, mean-wall ratio0.9904838640005559, median CPU1.019189765458422, lifetime-HWM RSS0.8964660068858583 and allocation0.9946304850967201. Its118 total warmup requests and10 measured requests all passed byte gates; this one-cohort near-neutral timing result is not a reliable speed claim. Three-cohort confirmation and bundle checks: PENDING. All completed screen paths, exact ratios, absolute sample statistics and source-reviewed limits are in completed-screens.json/evidence-index.json. The recommendation predates bitset/polyfill completion; its original five-screen scope is preserved, and bitset/polyfill here are derived separately from their actual independent reviews. No Java-versus-Go result follows.

Failed toolchain history is retained. build/control-launcher-missing encountered an absent javac executable (parent-reported); that directory has source/classes directories but no build.json/exit-code receipt, so no code or exit value is invented. build/control-release-missing has actual FAILED/exit2 and stderr “release version 21 not supported”; root diagnosed missing ct.sym support. Successful retries use the existing Java21 jdk.compiler module with matched -source21/-target21 flags for all overlays. Neither attempt is a factor-source rejection or a successful build; both remain in the evidence index.

Read methodology.md and publication-plan.md before admitting or publishing results. Root owns measurement scheduling, review and publication.

The first BigQuery STOCK/CONTROL screen (results/screen-bigquery-stock-control-v1) was REJECTED for “outer memory pressure guard”; no accepted performance claim follows. Its stopped receipt has zero new OOM, no cleanup errors, no owned containers/drivers and no active compilers. Historical aggregate OOM counters are not new OOM events.

Before a separately named retry, root recorded POSIX_FADV_DONTNEED for ordinary non-symlink workspace files (23,289 files;24,460,575,574 logical bytes; no content changes/errors) in results/cache-reclaim-before-screen-v1.json. Recorded parent memory.current fell from6,426,492,928 to4,449,914,880 bytes; memory.events stayed unchanged. This is pressure recovery, not a compiler RSS/speed result or GC setting. The later STOCK/CONTROL v4 retry completed and is one of the seven independently reviewed screens above; failed earlier versions stay separate.

Upstream-layout production patches are ready in patches/: five individual .patch files and all.patch, using src/com/google/... repository paths. They preserve exact unshaded candidate text against cached v20261015 stock sources. Each hunk was reconstructed in memory; actual upstream git-apply/build/full-suite checks remain pending. Original factor patches are unchanged.

Additional failed-screen and transport history

STOCK/CONTROL v1 and v2 were rejected by the outer-memory-pressure guard, both with zero new OOM and clean stopped receipts. No accepted performance result is substituted from them. v3 failed before normal benchmark initialization because the cgroup memory.pressure file was absent (launch exit1); only helper snapshots and launcher stderr/result exist, with no benchmark summary or stopped/OOM receipt. Missing PSI is not zero PSI.

Completed v4/factor screens use the explicitly recorded cache-aware-v3-PSIoptional runner. Headroom includes capped reclaimable inactive-file/slab accounting; warm reclamation is recorded/allowed, while measured local high/max/direct-scan/direct-steal and available major-fault deltas must be zero. PSI was unavailable in these reviews, so zero pressure is not proved by the remaining proxies. Protocol/runner versions and telemetry limits must accompany results.

Color’s underlying runner exited0 and its independent review passed, despite a failed outer sequence wrapper after terminal transport loss. results/transport-recovery-v1.json preserves the failed wrapper/sequence and acceptance evidence. A final stdout-write failure is inferred, not directly captured; the cause is not proven. Remaining factors restarted with durable logs under the same frozen protocol. The wrapper failure itself supplies no performance claim; completed/reviewed color evidence remains exploratory.

The fixed confirmation-plan-v1.json predates confirmation observations: color and ALL versus STOCK, three fresh paired-JVM cohorts and five measured pairs each; separate ALL-versus-STOCK bundle cohort. Primary is the geometric mean of cohort median-wall ratios. Mean/summed wall and CPU/RSS/allocation remain sensitivities; selected screens are excluded. The plan includes whole-cohort bootstrap, log-t(df2) and Bonferroni familywise adjustment for two primary comparisons; no reliable speed win if the adjusted interval includes1. No confirmation result is admitted here.
