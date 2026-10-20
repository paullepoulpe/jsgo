# Bounded CommonJS script cohorts

Source prototype based on `837a92e7`, off by default. Enable it with
`JSGO_EXPERIMENT_PARALLEL_COMMONJS=1` together with the existing fresh CLI
source ownership experiment, `JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1`.
The existing `--num_parallel_threads=N` bounds work items and private script
headers to at most N adjacent scripts. It does not change GOMAXPROCS or cap
runtime goroutines. No benchmark filename or source-size threshold determines
admission.

## Opportunity and limit

The existing whitespace profile attributes about 7.9% of whole-compilation CPU
to CommonJS rewriting. Perfectly halving that work with two workers saves at
most about 4% of whole-compilation CPU, before coordination and admission
overhead. This phase alone cannot establish Java parity. There is no measured
speedup, memory result, or retention decision for this prototype yet.

CommonJS already builds finder and rewrite scopes separately for each SCRIPT.
Its cross-file decisions use the module-type map built before this pass and
the loader's immutable path/package tables; it does not inspect a sibling's
rewritten body. That makes ordinary modules a substantially stronger boundary
than arbitrary script concurrency in Normalize or class lowering.

Normalize includes globally unique renaming and duplicate-declaration handling.
Class lowering has a concrete sequential dependency: rewriting a constructor
can change whether a later subclass considers that constructor to return only
undefined. Neither whole pass was made concurrent here.

## Ownership and exact publication

Only the private DefaultPassConfig factory opts in. The public
`NewProcessCommonJSModules` constructor retains its original traversal.
Admission requires the exact stock compiler, RAW/untyped ASTs, NODE or BROWSER
resolution, known stock conventions, warning guards, change observers, error
manager and source loaders, and an active CLI ownership witness. Unknown
resolvers, Webpack shims, custom callbacks, typed-AST restore paths and public
passes retain the original route.

Each eligible script owns its AST mutations, finder/rewrite scopes, pass
instance, loader header, concrete Node/Browser resolver header and event buffer.
Resolver maps, package aliases, search-folder slices and root-path slices remain
shared read-only. Changing a worker loader's error handler changes only that
worker's copied headers; it never swaps a handler on the canonical loader.

CompilerInput objects are never copied. Traversals, scope Vars, source-map and
input hooks see the original compiler input. Its lazy ModulePath is warmed
serially; a private immutable path header with the same normalized string is
used only at the two CommonJS resolution call sites. Name-only path operations
retain the canonical path. If a cached path belongs to a replaced loader, that
file takes the original serial route.

At most N adjacent scripts form a window. Contiguous eligible groups of at
least two run concurrently; unsupported files and isolated eligible files run
serially at their original position. Every group joins before the next group
or unsupported file starts. Events replay in original file order and original
per-file report order. Externs still run after all JS scripts. Heavy-first
scheduling within a group changes execution order only, not publication order.

Scope changes capture the original enclosing scope at report time, before a
rewrite can detach its ancestor. FUNCTION deletion flags and change stamps
publish only in the coordinator. The worker's embedded AbstractCompiler is nil,
so a newly introduced unaudited compiler dependency cannot silently mutate
shared state.

Diagnostics need more care than postponing `Compiler.Report`. A stock
`@suppress {moduleLoad}` on `module.exports = {x: require.ensure([], 3)}` can
suppress the finder warning before object-literal expansion removes the
assignment carrying that annotation. Reevaluating suppression after the rewrite
would expose a formerly hidden warning. The worker evaluates the admitted
read-only warning-guard chain at the original report point, freezes the final
level including OFF, and the coordinator passes that result directly to the
known stock error manager. Custom option error handlers are excluded. The
loader's explicit-level callback retains the compiler's original behavior of
ignoring that supplied level and deriving it from the diagnostic and guards.

## Global-name fallback

The UMD/IIFE mutators consume shared unique-name suppliers, clone arbitrary
function subtrees and can change structural assumptions across several
traversals. These scripts remain serial. An admission walk conservatively
rejects every original UMD condition trigger: module/define names,
window.define, the `"amd" in ...` marker, and true conditions. It also rejects
nested occurrences and closure-unaware shadows. This overrejects some ordinary
conditions in exchange for keeping the exact original global-name behavior.

Outside UMD, the pass's sole CloneTree duplicates a qualified export target;
function values are moved, not cloned. A valid require.ensure callback is
reparented within its original script. Worker change reports must target the
original SCRIPT or a FUNCTION originally owned by it.

## Validation still required

`TestParallelCommonJS*` contains source-only tests for:

- Actual NODE/BROWSER rewrites and cross-file imports at N=2 and N=4,
  canonical input/Var identities, private resolver headers and retained tables.
- Exact AST, printed output, source maps, feature sets, lifecycle, unique-name
  state, diagnostic output, ordered changed scopes and function/script stamps.
- Suppressed versus visible warnings whose ancestors change during rewriting.
- Mixed ordinary/UMD files with real dispatch before and after serial fallback.
- Real simultaneous jobs, bounded by the existing N budget.
- Package-main aliases, public constructors, N=1, unknown resolver callbacks,
  and ownership-lifetime fallback before mutation.

No build, focused test, race test or benchmark has been run for this commit.
Central validation should first run the focused suite normally and under race,
then all legacy CommonJS/module-resolution tests and Java artifact gates at
N=1/2/4 with the flag off and on. Source-map and warning suppression fixtures
are correctness gates, not optional performance diagnostics. A real CLI run
must confirm that the private DefaultPassConfig factory actually dispatches.

If those gates pass, compare same-binary off/on with default GC, fixed CPU
affinity, repeated interleaved runs and peak RSS. Admission adds one AST walk,
small windows can leave a large file unpaired, and the facade does not currently
activate every owner-only scope-recipe optimization. These costs can erase the
small whole-task ceiling; neither a goroutine count nor a microbenchmark proves
an end-to-end win.
