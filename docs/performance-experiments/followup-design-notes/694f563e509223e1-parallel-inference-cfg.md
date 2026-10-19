# Parallel inference CFG preparation

Status: isolated source prototype; no build, test, or performance claim yet.
Enable `JSGO_EXPERIMENT_PARALLEL_INFERENCE_CFG=1` with existing
`--num_parallel_threads=N` and native CFG enabled. The ordinary/default binary
path is unchanged. N=1, public constructors and unknown/custom observers keep
original scheduling. Runtime GC settings are not changed.

## Separation boundary

TypedScopeCreator.FinishAndFreeze prevents new scope creation, not mutations to
the type graph. Inference still updates canonical TypedVars, object properties,
registry and module-export facts, function signatures and AST JSType fields.
TypeCheck also updates casts, inferred annotations, mismatches, prototype
validation and JSDoc. Lazy type queries can mutate caches. These complete passes
cannot safely run per-function concurrently using the current representations.

ControlFlowAnalysis is a narrower independent phase: its scope, finally/exception
state, topology, priorities and work queues are per construction; it reads AST
syntax/parent links and options, not JSType annotations. It does not report
compiler diagnostics or change the AST. Stock inference changes type facts, not
AST syntax. All original native graph storage, edge ordering and prioritization
are reused unchanged.

The private builtin factory records exactly the original eligible scope work,
including the original current node, typed scope and script. After resolution
and FinishAndFreeze, it prepares at most N native CFGs concurrently with the
existing bounded parallelPrebuild helper. It then runs all inference in that
cohort serially in the original scope order, before preparing the next cohort.
No inference or type-graph write overlaps the CFG workers. Consumed graph
references are released, with deferred clearing on failure. Peak additional
prepared topology is bounded by N graphs, not all scopes; a single huge graph
can still dominate. The recorded work itself costs three pointers per eligible
scope. This storage/replay cost and scheduler overhead must be included in
whole-process measurements.

Source-context exception formatting calls SourceFile.GetLine, which lazily
writes line tables. Workers therefore use the exact successful traversal body
without its exception-context formatting. A failed preparation is discarded;
at its original serial ordinal the original complete InferScope builds the CFG
again and generates its normal exception/cause/context chain. Preceding semantic
inference and reporting still occur first. Test-only malformed CONTINUE fixtures
compare that full exception chain. Malformed/unknown public callback use remains
outside the admitted private fresh-CLI path.

## Existing parallelism is different

Compiler parsing/prebuild and dependency-info tasks already use N. Parallel FSIV
constructs separate native CFGs inside later FSIV jobs. Parallel coalescing is
another later phase. None supplies TypeInferencePass's earlier CFGs; structural
passes between them prevent simply reusing the later graphs. This prototype
uses the same N budget for compiler work items, without changing GOMAXPROCS or
limiting runtime/GC goroutines.

## Coverage and upper bounds

Only compilations that execute type inference benefit, primarily typed-app and
checks-only here. Ordinary whitespace/SIMPLE workloads without inference have
no benefit. The whole inference phase's reported 25.4% checks CPU share gives a
very loose ideal maximum of 1.34x even if that entire phase were parallel; this
prototype covers only its CFG subset. At N=2/4/8 that whole-phase ideal would be
1.146x/1.235x/1.286x. These are conditional Amdahl calculations, not elapsed-time
predictions, and exclude scheduler, retention, GC and CPU quota costs.

The earlier cached 0.85/7.93s number describes native **solver** compute, not CFG
construction; using its 10.7% as this prototype's parallel fraction would be
incorrect. A CFG-specific cumulative profile attribution is still required.

## Validation required

Source tests compare complete typed AST/scopes/output/diagnostics and exact
scope entry order for nested closures, defaults, fields/static blocks, modules
and try/finally at N=1/2/4. Independent concurrent speculative CFG construction
is compared against the original generic graph (node/edge/priority/order).
Public/custom/unowned paths must not prepare work. Malformed graph construction
must regenerate the original serial exception chain and clear recorded work.
Then run ordinary full compiler/race gates and all Java artifact gates before
any isolated N-scaling timing and peak-RSS comparison. Default runtime GC must
be retained for those measurements.

## Existing default-GC profile attribution

Decoded with Go 1.27.2 `go tool pprof -top -cum -nodefraction=0
-nodecount=10000 -tagfocus=jsgo_pass=inferTypes` on the existing
`structural-profiles/{task}/profile.cpu.pprof` captures, source 5b720eade6.
Their saved environment records default runtime GC, affinity CPUs 0/1, native
CFG enabled and the structural experiment bundle. No compiler or benchmark was
run for this analysis. Exact decoded outputs are
`round2/parallel-inference-cfg-{checks-only,typed-app}-cpu.txt` in the result root.

| Task | Whole sampled CPU | inferTypes CPU | CFG builder cumulative CPU | CFG / whole | CFG / inference |
| --- | ---: | ---: | ---: | ---: | ---: |
| checks-only | 8.57 s | 2.18 s | 0.11 s | 1.28% | 5.05% |
| typed-app | 5.93 s | 0.68 s | 0.04 s | 0.67% | 5.88% |

The builder is `ControlFlowAnalysis_Builder.computeNativeCfgWithQueries`; its
inline computeNativeCfg wrapper has the same cumulative count and is not added
again. Native solver compute separately accounts for 0.93s checks / 0.24s typed
and is not covered by this prototype. CPU stacks retain at most 64 PCs, so the
builder-inclusive ancestry counts are observable lower bounds, not rigorous
upper bounds on all construction cost; the phase label itself is exact even
when outer frames truncate. Sampling is coarse relative to these small masses.

Using only the observed fractions, an overhead-free N=2/4/infinite calculation
would give checks 1.006x/1.010x/1.013x and typed 1.003x/1.005x/1.007x. This is a
priority estimate, not a performance prediction: cohort retention, scheduler,
extra recording and GC can erase it. Defer this prototype behind larger proven
allocation targets and the phase-admission fixes. No source integration or
heavy test build is justified by the whole-inference 25.4% number alone.

Independent source review found no valid-stock-AST blocker. A narrow malformed
AST caveat remains: a few existing CFA assertions format Node.String(), which
can inspect type metadata even though normal graph construction is syntax-only.
The serial rebuild preserves the published exception chain but does not by
itself prove absence of every speculative formatting side effect on malformed
private ASTs. Keep this candidate isolated until stronger validation is useful.
