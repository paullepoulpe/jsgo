# CommonJS concurrency on current main

Source-only port onto `ff06ac15`, depending only on the scoped stock CLI source
lease (`65b4423f`). This keeps the original CommonJS cohort experiment and two
refinements as independent, default-off factors:

| Environment flag | Behavior |
| --- | --- |
| `JSGO_EXPERIMENT_PARALLEL_COMMONJS=1` | Original guarded adjacent N-script cohorts |
| `JSGO_EXPERIMENT_COMMONJS_EXACT_TRUE_GUARD=1` | Recognize literal TRUE as UMD only when it is the complete condition |
| `JSGO_EXPERIMENT_COMMONJS_RUN_POOL=1` | Schedule maximal eligible contiguous runs through main's existing pool |

The latter two flags require the first. The compiler's existing
`--num_parallel_threads=N` bounds active jobs; these factors do not change
GOMAXPROCS, runtime goroutines, the shared pool implementation or public API.
No benchmark filename or size threshold determines eligibility.

## Safety boundary

Only the private default-pass factory opts in. Public constructors, N=1,
typed ASTs/raw type annotations, custom passes/conventions/warning callbacks,
unknown change observers/resolvers and unknown source provenance retain the
original path. A live stock CLI source lease and an owned source witness are
required; compiler-generated exact preloaded sources are also allowed.

The phase-local facade contains only audited CommonJS methods. Unsupported
compiler dependencies are intentionally nil. It has no dependency on other
experiments' scope recipes, native facts, compiler facades or proof helpers.
Each worker owns one SCRIPT's mutations and its scopes, pass, loader/resolver
headers, ownership map and events. Original CompilerInput/Var identities and
canonical lazy input paths remain intact. Resolver lookup tables are shared
read-only; the private path is used only for actual module-resolution calls.

Diagnostics freeze the stock warning-guard result at the original report
point, before rewriting can remove an ancestor suppression annotation.
The coordinator publishes this final level directly to the proven stock error
manager. Changes, deletion flags and failures publish in original input and
per-input order. Externs run last. UMD and IIFE/name-generating paths remain
serial barriers, preserving shared unique-name supplier state.

The exact-TRUE factor fixes an overly conservative admission check on
ordinary conditions such as `value === true`. Other module/define/amd triggers
remain conservative. It does not itself prove more parallel CPU coverage.
`TestParallelCommonJSExactTruePinnedAngularAdmission` optionally reports actual
parsed AST admission/counts when `JSGO_COMMONJS_ANGULAR_AST_SOURCE` names the
local pinned Angular core source; no numeric claim precedes that test.

## Pool and memory tradeoff

The original path joins after at most N adjacent scripts. The run factor calls
the existing shared `parallelPrebuild` once per maximal eligible contiguous
run, schedules heavy files first and joins before ordered publication.
Unsupported scripts drain the preceding run, execute serially at their
original position, and then permit a new run. Singleton runs stay serial.

Queued descriptors contain original input/SCRIPT/path identities and FUNCTION
pointer slices. Workers construct facades/maps only when starting a job and
retain only events/failure afterward. Private analysis state remains bounded
by N; queued descriptors, function pointers and completed events can exceed N.
Event pointers can retain detached nodes until the run joins. This avoids
N-file head-of-line barriers but may raise peak RSS. It constructs no duplicate
AST; both elapsed time and RSS require measurement before retention.

## Evidence and remaining gates

On the older source, five interleaved matched whitespace runs showed the
original factor slower at N=2 (1.270s on / 1.176s off) and N=4
(1.187s on / 1.170s off). Those measurements are not current-main results.
The N=1 factor is inactive and its apparent difference is control variation.
Neither original admission coverage nor a scheduling redesign establishes
a speedup on current main.

The port and both refinements are unbuilt and unmeasured at this checkpoint.
Focused `^TestParallelCommonJS` tests cover exact output/AST/maps/features,
canonical identity, private loader state, names/counters, frozen diagnostics,
change-stamp order, custom fallback and expired leases. Pool tests require
real overlap within N, progress through more than N jobs behind a blocked
first file, serial UMD barriers, and original failure-prefix publication after
out-of-order failures. Central normal/race tests and all CLI/Java artifact
gates must pass before matched current-main timing and peak-RSS comparisons.

`tools/bench/commonjs_main_activation.py --root CHECKOUT --out NEW_DIRECTORY`
writes an untimed overlay without building or running Go. Optional repeated
`--case EXISTING_CASE.json` arguments copy pinned cwd/args into fresh report
paths. The central test build can include that overlay, then run only
`^TestParallelCommonJSMainBenchmarkActivation$` with
`JSGO_COMMONJS_ACTIVATION_CASE` set to a generated case. The test enables the
original factor and leaves exact-TRUE/run-pool selection to the caller.
It records actual private-factory execution, acceptance, prospective admission
and actual dispatch/launch counts, including whole-phase fallback and absent
passes. The extra diagnostic scans and elapsed values are untimed activation
evidence, never benchmark results. Production source remains unchanged.
When the checkout also includes the clean-main ASCII parser factor, optional
`--ascii` attaches its existing numeric runner probe before that same DoRun and
includes parse/source-byte/backing counters afterward. No separate parse or
compilation is performed, and the ASCII factor's enable flag stays caller-owned.
Optional `--fsiv` calls the separate untimed FSIV diagnostic overlay's numeric
reset/snapshot around that same DoRun and includes its counters in `extras.fsiv`.
The resulting overlay must be composed with that FSIV helper overlay; this
option adds no FSIV production hook and does not turn its optimization on.
