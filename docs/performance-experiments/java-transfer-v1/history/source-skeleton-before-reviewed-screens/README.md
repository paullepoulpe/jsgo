Draft: transferring Go optimizations to warm Java

Measurement results: PENDING root review. No timing, CPU, RSS or allocation result is admitted by this skeleton. Preparing it runs no compiler/benchmark and changes no published record or stock oracle.

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

| Planned comparison | Task | Admission | Elapsed ratio | RSS lifetime-HWM ratio | Allocation ratio | Decision |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL / STOCK | BigQuery | PENDING | PENDING | PENDING | PENDING | PENDING |
| Each factor / CONTROL | BigQuery | PENDING | PENDING | PENDING | PENDING | PENDING |
| ALL / CONTROL | BigQuery | PENDING | PENDING | PENDING | PENDING | PENDING |
| Selected factors / CONTROL | bundle-n1 | PENDING selection | PENDING | PENDING | PENDING | PENDING |
| Selected factors or ALL / STOCK | Prespecified confirmation | Three independent JVM-pair cohorts planned | PENDING | PENDING | PENDING | PENDING |

Ratios use candidate/baseline; below1 means lower cost. Fill only actual accepted measurements, with absolute values, variance, cohort counts and exact memory windows. Build/probe durations and earlier Go/Java results cannot fill these cells. No Java-versus-Go conclusion follows from this Java-versus-Java study.

Failed toolchain history is retained. build/control-launcher-missing encountered an absent javac executable (parent-reported); that directory has source/classes directories but no build.json/exit-code receipt, so no code or exit value is invented. build/control-release-missing has actual FAILED/exit2 and stderr “release version 21 not supported”; root diagnosed missing ct.sym support. Successful retries use the existing Java21 jdk.compiler module with matched -source21/-target21 flags for all overlays. Neither attempt is a factor-source rejection or a successful build; both remain in the evidence index.

Read methodology.md and publication-plan.md before admitting or publishing results. Root owns measurement scheduling, review and publication.

The first BigQuery STOCK/CONTROL screen (results/screen-bigquery-stock-control-v1) was REJECTED for “outer memory pressure guard”; no accepted performance claim follows. Its stopped receipt has zero new OOM, no cleanup errors, no owned containers/drivers and no active compilers. Historical aggregate OOM counters are not new OOM events.

Before a separately named retry, root recorded POSIX_FADV_DONTNEED for ordinary non-symlink workspace files (23,289 files;24,460,575,574 logical bytes; no content changes/errors) in results/cache-reclaim-before-screen-v1.json. Recorded parent memory.current fell from6,426,492,928 to4,449,914,880 bytes; memory.events stayed unchanged. This is pressure recovery, not a compiler RSS/speed result or GC setting. Retry output uses a new path; its measurements are not read/admitted by this skeleton.

Upstream-layout production patches are ready in patches/: five individual .patch files and all.patch, using src/com/google/... repository paths. They preserve exact unshaded candidate text against cached v20261015 stock sources. Each hunk was reconstructed in memory; actual upstream git-apply/build/full-suite checks remain pending. Original factor patches are unchanged.
