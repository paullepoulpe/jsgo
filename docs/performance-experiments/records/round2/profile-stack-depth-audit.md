# Profile attribution correction

Prior retained-stack family totals are observed attribution, **not hard whole-program performance ceilings**. Shared generic traversals can lose the owning compiler pass when the stack is truncated. Indirect GC costs are also outside a pass's retained stack. This corrects earlier ceiling wording in the parser, emitter and type-dispatch investigations.

## Go1.27.2 source facts

- `runtime/cpuprof.go:22`: CPU `maxCPUProfStack = 64`.
- `runtime/proc.go:5842`: the signal handler captures into `[maxCPUProfStack]uintptr`.
- `runtime/extern.go:167`: `GODEBUG=profstackdepth` explicitly applies to pprof profilers **except CPU**. Default128, maximum1024; `runtime/runtime1.go:429` clamps it to1024.
- `runtime/pprof/proto.go:406`: profile encoding groups inline PCs into locations and expands a final truncated inline frame. Therefore pprof location-ID count is not raw captured-PC count. Expanded logical stacks can slightly exceed64 (CPU) or128 (heap).

Use `GODEBUG=profstackdepth=512` for future heap allocation captures. The attempted depth1024 capture failed inside Go1.27.2 runtime.mprof with `bad profile stack count` and a1026-frame stack; it is archived under `next-profiles-depth1024-failed`. This is a diagnostic runtime limit failure, not a compiler production result. It does not solve CPU attribution. CPU pass labels are the low-impact diagnostic remedy; a custom deeper-stack runtime build is another option but changes the profiling toolchain.

## Existing enabled11534351 captures

| Task | CPU total | Expanded logical depth>=64 | Generic traversal, no retained Process entry | GetNext leaf under traversal, no Process entry |
|---|---:|---:|---:|---:|
| typed-app | 6.87s | 0.76s (11.1%) | 0.90s | 0.15s |
| bundle | 26.59s | 8.08s (30.4%) | 1.88s | 0.16s |
| downlevel-whitespace | 1.91s | 0.33s (17.3%) | 0.02s | 0s |
| downlevel-simple | 5.25s | 0.61s (11.6%) | 0.10s | 0s |
| checks-only | 10.20s | 1.76s (17.3%) | 2.89s | 0.25s |

Depth>=64 is a saturation-candidate measure, not proof each sample lost a parent. Missing Process entries can also reflect traversal work outside a Process method, so those counts are not all truncation loss. Source-scanned entry recognition is broad and cannot uniquely identify the actual factory pass. Broad callback context (Visit/ShouldTraverse/EnterScope/ExitScope) can survive even while the factory/Process owner is missing; retaining such a callback must not be mistaken for full attribution.

There is direct truncation evidence: **1.05s bundle CPU ends at a generic NodeTraversal frame**, so the older outer-owner frames are absent. The corresponding heap profile has **150.4MiB sampled allocation ending at a NodeTraversal frame**. Thus compiler-pass totals can undercount generic traversal cost. CommonJS6.8% retained CPU is a lower bound, not an upper bound that rules out a larger redesign.

Bundle has2248.7MiB sampled allocation at expanded depth>=64, including208.1MiB at>=128. The128-depth heap issue is smaller than the CPU64 issue but still material. Exact total allocation comes from MemStats, not these weighted heap samples.

Machine-readable buckets/histograms/tail frames: `profile-stack-depth-audit.json`. Reproducer: `audit_profile_stack_depth.py`. No new workload/profile was run for this audit.

## Prepared diagnostic overlay

`next-profile-overlay/overlay.json` is bound to clean eecbcf44af08ed6743de5ce69eddce675a71ee93 in `/workspace/jsgo-round-compound-next`. It combines the existing main CPU/alloc/live/MemStats capture with `jsgo_pass` CPU labels at one named factory pass site, two compiler factory loops, one direct Compiler.Process site, and the dependency-scan callsite. AST retrieval is evaluated before the dependency label. BeforePass/AfterPass/tracer options and experiment guards are unchanged.

The instrumentation boundaries are disjoint for the supplied stock benchmark workloads. Custom passes that recursively enter another labeled boundary require explicit parent-context propagation; they are outside this diagnostic harness. New goroutines inherit CPU labels; heap profiles do not support pprof labels, so heap depth must be enlarged separately. Capture and byte-gate first; inspect `go tool pprof -tags`, then use `-tagfocus=jsgo_pass=<name>` to attribute generic walking. Diagnostic profile durations are not production benchmark results.
