Neither selected candidate cleared the prespecified familywise-adjusted median-wall criterion against stock Java. The promising single-cohort screens therefore do not establish a reliably faster warm Java compiler on this workload.

These are Java-versus-Java experiments on Closure Compiler v20261015, with fresh compiler construction and normal extern parsing on every request. They do not update any Go result or establish a new Java/Go ratio. The completed analysis and report have been reviewed; the publication receipt records archive verification.

The observed ratios below compare each candidate with stock Java; lower is better. They summarize the completed follow-up runs, not the exploratory screens.

| Candidate and workload | Fresh paired-JVM cohorts | Median time ratio | Approximate peak RSS ratio | Allocated-byte ratio |
| --- | ---: | ---: | ---: | ---: |
| Color cache, BigQuery | 3 | 1.001 | 1.111 | 0.996 |
| All five changes, BigQuery | 3 | 1.020 | 0.918 | 0.989 |
| All five changes, bundle | 1 | 1.050 | 1.020 | 0.982 |

The adjusted time intervals are 0.876–1.145 for color and 0.906–1.149 for ALL. Both include no change, and these small samples cannot rule out useful gains or regressions. RSS changes are also uncertain; the bundle is a single-cohort screen. No optimization is recommended as a demonstrated speed win from this study.

Download the [complete evidence archive](java-transfer-v1-evidence.tar.gz), [screen data](completed-screens.json), [final analysis](confirmation-analysis.json), [comparison CSV](confirmation-comparisons.csv) and [raw cohort ratios](confirmation-cohort-ratios.csv). Receipt paths shown in code below are paths inside the evidence archive, not web links.

BigQuery uses ADVANCED with dependency pruning: 1,270 supplied input files total 22,180,254 bytes, including 1,209 JavaScript sources totaling 21,607,844 bytes and 227 generated protobuf-model JavaScript files. The retained manifest has 251 inputs; the emitted JavaScript is 161,467 bytes. Supplied bytes therefore overstate the retained application size. The separate ADVANCED bundle contains 11 pinned library inputs totaling 5,677,811 bytes (including Three.js, Lodash, React DOM, RxJS, D3, Moment, Vue, jQuery, Handlebars, Immutable and Underscore). Both use compiler N=1 inside the same two-CPU, 4 GiB limits. These workloads do not cover every large chunked enterprise build.

| Factor | Change | Boundary |
| --- | --- | --- |
| [Lazy FSIV analysis](patches/fsiv-lazy.patch) | Initialize backward use analysis only at the original use query; full variable map and candidate order unchanged. | Initialization precedes first successful inline. Skipping unused analysis may omit its pathological error; broader failure-path parity is unproved. |
| [Lazy namespace context](patches/namespace-lazy.patch) | Request Scope/JSChunk only after existing irrelevant-token/name exits; preserve explicit incremental context. | Custom compiler callbacks or input errors may observe changed getter timing; full extension equivalence is unproved. |
| [BitSet interference](patches/coalesce-bitset.patch) | Bulk-OR live sets instead of repeatedly setting variable pairs. | Upper triangle/diagonal/order preserved; extra lower bits unobserved. Rows do not alias lattice storage; ordinary analysis caps at 100 variables. |
| [Color identity cache](patches/color-cache.patch) | Per-factory cache before repeated normalization; protected external-map constructor remains uncached. | Preserve canonical insertion/IDs and errors; extra retained original-color entries may increase memory. |
| [Default polyfill table](patches/polyfill-cache.patch) | Lazily share immutable resource metadata; custom table injection unchanged. | No compiler/AST/extern cache. One table remains per classloader; successful resource immutability and normal loader anchors assumed. |

[ALL](patches/all.patch) combines all five over nine disjoint production source files. STOCK is the unchanged published jar; CONTROL recompiles those same nine source families with matched flags/relocation, and every factor uses the same compiler setup. CONTROL/STOCK is an experimental-artifact control, not a fixed rebuilding penalty to subtract from another cohort. Final holdouts compare candidates directly with STOCK.

Exploratory BigQuery screens: one pair of JVMs and five measured pairs each; no confidence intervals.

| Candidate / baseline | Median wall ratio | Mean wall ratio (sensitivity) | Median CPU ratio | Max observed approximate VmHWM ratio | Median allocated-byte ratio |
| --- | --- | --- | --- | --- | --- |
| CONTROL / STOCK | 1.0204666067035644 | 1.0264733296906947 | 1.01875 | 1.0004940367233965 | 0.9944049755826357 |
| ALL / CONTROL | 0.9740394275639087 | 0.995963526747487 | 1.0082815734989647 | 0.9245254011288077 | 0.9953882323641364 |
| fsiv-lazy / CONTROL | 1.0331203821346027 | 1.0350554318216456 | 1.028747433264887 | 0.8373530481153799 | 0.9854329236494126 |
| namespace-lazy / CONTROL | 0.9482400327483805 | 0.9864669318690914 | 0.9542743538767395 | 0.9414526973669367 | 0.9997582269845262 |
| color-cache / CONTROL | 0.9475562607922867 | 0.9617243024010067 | 0.95703125 | 0.9779596635615447 | 1.0065805416923779 |
| coalesce-bitset / CONTROL | 0.9936485460026859 | 1.0094435310257523 | 0.9758551307847082 | 0.9542267185286128 | 0.9962961511652734 |
| polyfill-cache / CONTROL | 0.993503390650036 | 0.9904838640005559 | 1.019189765458422 | 0.8964660068858583 | 0.9946304850967201 |

All ratios use candidate/baseline; below1 means lower cost. Means, medians of paired ratios and ratios of medians are different estimators. Namespace and ALL screen reductions weaken under mean-wall sensitivity; bitset changes direction. All five factors, neutral/negative screens and ALL are retained. Sampled RSS reductions from one JVM pair have no significance claim, and factor percentages cannot be added or multiplied.

Prespecified BigQuery holdouts: two selected candidates, three fresh independent paired-JVM cohorts each, five measured pairs per cohort; no screen rows pooled.

| Candidate | All 3 raw candidate/STOCK median-wall ratios | Geometric mean | Familywise 95% adjusted interval | Inverse STOCK/candidate | Inverse adjusted interval | Finding |
| --- | --- | --- | --- | --- | --- | --- |
| color-cache | 1.0265508678425144, 0.959121098430544, 1.019954879937816 | 1.001409313991265 | [0.8755022882650035, 1.1454231788882712] | 0.9985926693794689 | [0.8730397799096257, 1.1422014692636788] | did not clear adjusted criterion |
| ALL | 0.9819304407735583, 1.0382761312206144, 1.041879091168268 | 1.0203213450731439 | [0.905707265854835, 1.149439434196532] | 0.9800833872864954 | [0.8699892923884316, 1.1041095039203077] | did not clear adjusted criterion |

The overall finding uses adjusted upper bound < 1, not the most favorable screen or unadjusted interval. Bonferroni covers two prespecified median-wall comparisons: individual 97.5% log-t intervals with 2 degrees of freedom. Three cohort effects are the inference units; persistent-JVM requests are correlated. Log-t assumes independent approximately normal log effects, which three cohorts cannot test. Whole-cohort bootstrap intervals are discrete and confined to the three observed cohorts.

color-cache: all metrics and interval sensitivity

| Metric | All 3 raw candidate/STOCK ratios | Geometric mean | Bootstrap 95% | Log-t 95% | Inverse STOCK/candidate | Inverse bootstrap 95% | Inverse log-t 95% |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Median wall time (primary) | 1.0265508678425144, 0.959121098430544, 1.019954879937816 | 1.001409313991265 | [0.959121098430544, 1.0265508678425144] | [0.9123258777687, 1.0991912414027774] | 0.9985926693794689 | [0.974135847843258, 1.0426212098100522] | [0.9097597964152349, 1.0960995674546983] |
| Mean wall time | 1.0254099705557795, 0.9481094405712349, 1.0136548131922514 | 0.9951350691122385 | [0.9481094405712349, 1.0254099705557795] | [0.8958346387255823, 1.105442637478068] | 1.004888714144203 | [0.9752196962332957, 1.0547305587396125] | [0.9046150076872171, 1.116277443147994] |
| Median process CPU | 1.002020202020202, 0.99, 1.029598308668076 | 1.0070703990393073 | [0.99, 1.029598308668076] | [0.9580517541067898, 1.0585970792013626] | 0.9929792405316926 | [0.9712525667351131, 1.0101010101010102] | [0.9446464756491015, 1.0437849476433758] |
| Mean process CPU | 0.9967186218211647, 0.9832134292565946, 1.0185263157894737 | 0.9993805114070797 | [0.9832134292565946, 1.0185263157894737] | [0.9561684899386388, 1.0445454091928617] | 1.000619872596923 | [0.9818106655642828, 1.0170731707317073] | [0.9573542626286753, 1.0458407807019179] |
| Max observed approximate VmHWM RSS | 1.2005074445457309, 1.0682272549032212, 1.0689889811495168 | 1.1108802531153006 | [1.0682272549032212, 1.2005074445457309] | [0.9400948398917621, 1.3126919587215233] | 0.9001870338370377 | [0.8329810902408836, 0.9361303930506787] | [0.7617933463795535, 1.0637224645495715] |
| Median allocated bytes | 0.994013272366312, 1.0008715042922718, 0.9924597659956866 | 0.9957748178175004 | [0.9924597659956866, 1.0008715042922718] | [0.9847300849938807, 1.0069434283666048] | 1.0042431100956741 | [0.9991292545661113, 1.0075975210911936] | [0.9931044503881733, 1.015506701012607] |
| Mean allocated bytes | 0.9939636909797258, 1.0010431824834625, 0.9925778072578653 | 0.9958546701290786 | [0.9925778072578653, 1.0010431824834625] | [0.9846519784790034, 1.0071848182845482] | 1.0041625851595233 | [0.9989579046121921, 1.0074776936254897] | [0.9928664350831008, 1.0155872550469098] |

ALL: all metrics and interval sensitivity

| Metric | All 3 raw candidate/STOCK ratios | Geometric mean | Bootstrap 95% | Log-t 95% | Inverse STOCK/candidate | Inverse bootstrap 95% | Inverse log-t 95% |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Median wall time (primary) | 0.9819304407735583, 1.0382761312206144, 1.041879091168268 | 1.0203213450731439 | [0.9819304407735583, 1.041879091168268] | [0.9394101664202231, 1.1082013846826708] | 0.9800833872864954 | [0.9598042694941611, 1.018402076640181] | [0.9023630666968947, 1.0644977409715124] |
| Mean wall time | 0.9818025713311562, 1.0087367695021743, 1.0267136203468665 | 1.005581126040603 | [0.9818025713311562, 1.0267136203468665] | [0.9508516992169832, 1.0634606867525132] | 0.9944498500458354 | [0.9739814298578784, 1.0185347127826028] | [0.9403262503795012, 1.0516887132067914] |
| Median process CPU | 0.956777996070727, 0.9978813559322034, 1.0063025210084033 | 0.9867479118615834 | [0.956777996070727, 1.0063025210084033] | [0.922645256948387, 1.0553042290420214] | 1.0134300645373704 | [0.9937369519832986, 1.0451745379876796] | [0.9475940420591082, 1.0838401785182974] |
| Mean process CPU | 0.972520908004779, 1.0109565950273913, 0.9995862639635913 | 0.9942232434192873 | [0.972520908004779, 1.0109565950273913] | [0.9461931618386862, 1.0446913987777164] | 1.0058103213930563 | [0.989162150896207, 1.0282555282555284] | [0.9572204779038047, 1.0568666529535617] |
| Max observed approximate VmHWM RSS | 0.9969572945857555, 0.8277269437551695, 0.9365546559643669 | 0.9176961822859488 | [0.8277269437551695, 0.9969572945857555] | [0.7253856109047707, 1.160991161558574] | 1.0896852567360968 | [1.0030519917259935, 1.208127882684687] | [0.8613329998632796, 1.37857711121772] |
| Median allocated bytes | 0.9802020594711054, 0.989922431381145, 0.9961022208002253 | 0.9887205514632548 | [0.9802020594711054, 0.9961022208002253] | [0.9689909386673818, 1.0088518786669118] | 1.0114081259058003 | [1.0039130313319082, 1.0201978156825922] | [0.9912257895791317, 1.0320013945385949] |
| Mean allocated bytes | 0.980248799320546, 0.990025100634799, 0.9960968831813698 | 0.988768681343332 | [0.980248799320546, 0.9960968831813698] | [0.9690861538117413, 1.0088509688843996] | 1.0113588940149372 | [1.0039184108338581, 1.0201491709993875] | [0.9912266834672447, 1.0318999978140893] |

Mean wall/CPU, RSS and allocation are descriptive sensitivities, not additional adjusted confirmatory endpoints. No memory significance claim is inferred from a screen or selected secondary interval. All exact analysis values are copied without recomputing intervals.

Exact absolute measured statistics for each holdout JVM

| Archive result job | Cohort | Variant | Wall median s | Wall sample SD s | CPU median s | CPU sample SD s | Max observed approximate VmHWM MiB | Current RSS median MiB | Current RSS sample SD MiB | Allocated median bytes | Allocated sample SD bytes | Warm requests |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| confirm-bigquery-stock-color-cache-v1 | 1 | STOCK | 4.186277626 | 0.12572841628434128 | 4.95 | 0.20863844324572586 | 1285.546875 | 725.296875 | 173.92539070351998 | 2082470552 | 544510.831610171 | 41 |
| confirm-bigquery-stock-color-cache-v1 | 1 | color-cache | 4.29742693 | 0.13884231366187225 | 4.96 | 0.31408597549078826 | 1543.30859375 | 1229.40625 | 408.98016001141053 | 2070003368 | 631404.5685637696 | 41 |
| confirm-bigquery-stock-color-cache-v1 | 2 | STOCK | 4.301241062 | 0.1907801520271909 | 5.0 | 0.13501851724856134 | 1215.3203125 | 903.0078125 | 143.01520552455418 | 2080375296 | 611291.2508426732 | 45 |
| confirm-bigquery-stock-color-cache-v1 | 2 | color-cache | 4.125411052 | 0.17512708717774111 | 4.95 | 0.1835755975068582 | 1298.23828125 | 585.34765625 | 54.23986468588005 | 2082188352 | 112600.8729486588 | 45 |
| confirm-bigquery-stock-color-cache-v1 | 3 | STOCK | 4.117783332 | 0.0749324507878079 | 4.73 | 0.21365860619221497 | 1217.01953125 | 658.8515625 | 173.37795240422818 | 2091751528 | 398754.23341802903 | 58 |
| confirm-bigquery-stock-color-cache-v1 | 3 | color-cache | 4.199953204 | 0.09613404837146712 | 4.87 | 0.14237275020171522 | 1300.98046875 | 729.34765625 | 48.931847677796576 | 2075979232 | 295025.9926501392 | 58 |
| confirm-bigquery-stock-all-v1 | 1 | STOCK | 4.296817594 | 0.2017229443352065 | 5.09 | 0.3172853605195173 | 1235.0234375 | 604.3203125 | 71.17004617830777 | 2094515232 | 147940.34717817855 | 41 |
| confirm-bigquery-stock-all-v1 | 1 | ALL | 4.219175994 | 0.0530178869786431 | 4.87 | 0.061886993787063346 | 1231.265625 | 630.5625 | 97.3488564963401 | 2053048144 | 490802.87381717726 | 41 |
| confirm-bigquery-stock-all-v1 | 2 | STOCK | 4.072717462 | 0.10090764025044381 | 4.72 | 0.199449241663136 | 1511.25 | 697.765625 | 189.25082351727121 | 2087460656 | 640321.5681240169 | 47 |
| confirm-bigquery-stock-all-v1 | 2 | ALL | 4.22860533 | 0.1533579483263646 | 4.71 | 0.22928148638736623 | 1250.90234375 | 576.421875 | 241.85662243009486 | 2066424128 | 108956.84988104236 | 47 |
| confirm-bigquery-stock-all-v1 | 3 | STOCK | 4.045122668 | 0.11497638390029272 | 4.76 | 0.15043270920913454 | 1347.0625 | 589.10546875 | 30.284554788340994 | 2073328320 | 120272.13116927794 | 58 |
| confirm-bigquery-stock-all-v1 | 3 | ALL | 4.214528729 | 0.12175007759052227 | 4.79 | 0.17137677789012148 | 1261.59765625 | 904.578125 | 173.87733536666423 | 2065246944 | 122670.23963781925 | 58 |
| confirm-bundle-stock-all-v1 | 1 | STOCK | 10.594292571 | 0.3121224144549364 | 12.42 | 0.4496109429273272 | 1313.171875 | 693.82421875 | 26.565006114543294 | 4743941984 | 205467.73898400692 | 20 |
| confirm-bundle-stock-all-v1 | 1 | ALL | 11.126084782 | 0.3042687513944772 | 12.79 | 0.42811213484319766 | 1339.28515625 | 672.515625 | 33.11005884016103 | 4658801896 | 372118.24347000243 | 20 |

Each absolute median/SD uses five measured requests in that JVM. The RSS summary is the maximum observed approximate /proc VmHWM across warmup, JIT and measurement, not an exact physical lifetime peak. Legacy JSON fields retain their original HWM names. There is no per-request peak SD or proved accuracy bound. Intervals describe this approximate telemetry, not bounded error in true memory usage. Current RSS is an approximate post-request snapshot; cgroup peak includes charged cache/helpers. Allocated bytes are ThreadMXBean totals with management-boundary overhead, not live heap or RSS. Exact exploratory absolute samples/statistics remain in completed-screens.json.

Bundle transfer check: one cohort, external-validity screen only; no CI and not pooled into BigQuery.

| Metric | ALL / STOCK ratio |
| --- | --- |
| median_wall | 1.0501961039338943 |
| mean_wall | 1.0352914407649476 |
| median_CPU | 1.0297906602254427 |
| mean_CPU | 1.027652733118971 |
| lifetime_HWM | 1.0198856537724736 |
| median_allocated_bytes | 0.982052881698985 |
| mean_allocated_bytes | 0.9821072868373226 |

Bundle absolute JVM values are in the preceding table. One cohort cannot establish significance or a universal workload improvement.

Separate instrumented ALL diagnostics: activation, not timing evidence

| Counter | BigQuery | Bundle |
| --- | --- | --- |
| fsiv.eligibleFunctions | 1110 | 9766 |
| fsiv.backwardAnalyses | 82 | 2449 |
| fsiv.noCandidates | 8 | 1123 |
| namespace.traversalCollectCalls | 1196200 | 1904235 |
| namespace.scopeResolutions | 273530 | 654683 |
| namespace.chunkResolutions | 128894 | 46048 |
| coalesce.orIntoCalls | 63127 | 593103 |
| color.lookupCalls | 58340 | 142594 |
| color.cacheHits | 52260 | 135356 |
| color.cacheMisses | 6080 | 7238 |
| polyfill.getDefaultCalls | 13 | 12 |
| polyfill.resourceParseCalls | 1 | 1 |

All five paths activated on both workloads with matching six artifact/provenance/cap/stopped receipts. FSIV eligible counts are visits, not necessarily unique functions. These two instrumented requests were outside timings and are never used to estimate speed or memory benefit. Archive receipt: `diagnostic-activation/results-review-v1.json`.

Post hoc telemetry-review amendment

The color three-cohort benchmark runner completed successfully, but the original strict independent reviewer failed on STOCK cohort 3: reported VmHWM fell by 60 KiB between warmup requests 7 and 8 (1107716 to 1107656 KiB). The later maximum was 1246228 KiB (1217.01953125 MiB), strictly above the earlier sample. The numeric maximum, raw six-digest rows, timing and cgroup peak were unchanged.

After observing this anomaly, the review policy was amended to permit only recorded warmup-only VmHWM dips later strictly exceeded. Measured VmHWM dips and any cgroup memory.peak dip remain fatal; every independent admission check was rerun. This was a post hoc REVIEW amendment, not a prespecified rule or benchmark-protocol tuning. The benchmark runner, flags, inputs and raw values were unchanged. Linux documents VmHWM as inaccurate; the exact cause and numerical error bound remain unproved. The absence of a change in this maximum does not establish precise RSS measurement.

The original failed review and failed progress snapshot remain archive evidence at `reviews/confirm-bigquery-stock-color-cache-v1-review-v1-failed/review.json` and `results/review-telemetry-amendment-v1/`. The fixed plan resumed through `confirm_resume.py`; the mutable confirmation sequence retains the final closed state while its failure-time snapshot is preserved. The amendment note now records the completed full v2 PASS and links its accepted review and reviewer hashes. Documentation receipt: `bench/proc-status-documentation-v2.json`.

Correctness and upstream readiness

All six normalized patches passed `git apply --check`, applied successfully to scratch trees of the pinned original sources, and reproduced the exact unshaded candidate files ([receipt](patch-apply-validation.json)). This does not substitute for an upstream build. All seven compiled overlays and the shared driver have build receipts; focused probes passed for all five factors (coalescing: 1,497 actual-helper cases; polyfill: 4 controls). Every admitted warmup and measurement passed the stock stdout/stderr/exit/JS/map/manifest digests. Full upstream Closure suites, broad custom-extension behavior and universal error-path equivalence have not been run or established. ALL inherits each factor’s limitations above. Read the [historical upstream-readiness review](upstream-readiness.md), written before these holdouts; its pending-screen interpretations do not override the completed analysis here.

Failure and recovery history

The first final patch-validation helper assumed the individual-factor manifest schema for ALL and stopped with `KeyError: source_files`. The corrected helper validated all six patches; both receipts are retained under `results/upstream-patch-apply-v1/` and `results/upstream-patch-apply-v2/`. No measured implementation changed.

The absent javac launcher attempt has retained directories but no build.json or captured exit value; no missing receipt is invented. The subsequent --release 21 attempt failed with exit 2 because release 21 support was unavailable (root diagnosed missing ct.sym). Matched Java21 compiler-module -source 21/-target 21 builds then passed; neither failure is silently converted to a factor-source result.

STOCK/CONTROL screen v1 and v2 were rejected by the parent-memory guard with zero new OOM and clean stopped receipts. Before a retry, recorded POSIX_FADV_DONTNEED reclaimed file-cache pressure without changing content; it is not compiler RSS or GC tuning. v3 failed at initial memory.pressure lookup; no benchmark summary/stopped receipt exists for it, and missing PSI cannot be called zero. Completed runs use the frozen cache-aware-v3-PSIoptional policy, with measured local/direct-reclaim/available fault checks as proxies. PSI remained unavailable; absence of pressure is not proved. Historical aggregate OOM counters are distinct from new deltas.

Color’s underlying screen runner exited 0 and independent review passed despite a failed outer sequence wrapper after transport loss. Its successful byte/cleanup evidence is retained alongside the failed wrapper; a final stdout-write failure is inferred, not proven. Remaining factors restarted with durable logs. Archive paths: `results/transport-recovery-v1.json`, `results/cache-reclaim-before-screen-v1.json`, `screen-sequence-v1.json`, `build/control-launcher-missing/` and `build/control-release-missing/`.

Open work

Memoizing immutable Color.hashCode is an explicitly unimplemented, unmeasured idea: first profile actual repeated hash work and verify immutability/equality/ownership, then test and measure any patch separately. Other open work includes upstream pass/error/custom-callback suites, dense coalescing workloads, retained-heap measurements for caches, more independent JVM cohorts and additional generated/large/chunked workloads. Existing factor screens do not justify assuming any unmeasured interaction or new gain.

Method: same JDK 21, -Xmx3g and default collector, with actual 4 GiB / 2 CPUs per JVM, alternating sequential requests, matched 12–90 warmups, unchanged drift/JIT/spread criteria, fresh compiler/normal extern parsing every request, strict six digests and zero-new-OOM/stopped evidence. The inactive JVM is paused but its pages remain charged. No worker, extern AST or result cache is introduced. See [methodology](methodology.md), [evidence index](evidence-index.json) and the archive’s raw records. Paths require remapping and external pinned JDK/jar/corpus inputs; [reproduction notes](REPRODUCE.md) accompany the published snapshot.
