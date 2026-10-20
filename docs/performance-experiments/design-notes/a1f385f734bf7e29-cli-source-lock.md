# CLI-owned source loading synchronization

`JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1` specializes synchronization for fresh files
created by the exact stock CommandLineRunner. It does not cache parsed externs
or compiler state. Raw source text keeps its existing load/clear lifecycle.

Repeated dependency-info queries dominated the observed SourceFile monitor
cost: CompilerInput.GetCode → generateDependencyInfo → GetProvides → module
discovery. A reentrant Java-style monitor obtains a goroutine ID through
runtime.Stack on every lock and unlock. The known stock preload/disk/zip loaders
do not call back into their SourceFile, so that owner boundary can use a plain
mutex instead.

The owner is a private CodeLoader wrapper, installed after stock CLI source
creation and before Compiler.Init publishes inputs or starts parallel work. It
adds no fields to SourceFile and uses no global file map. A wrapper-local atomic
pointer publishes the same immutable cached string that SourceFile already
owns; it is not another text cache. Existing loaded reads stay lock-free.
Cold loading and SetCode/ClearCachedSource use the wrapper's ordinary mutex.
The original bookkeeping is a shared private primitive, so GetCode→SetCode and
GetCodeReader→GetCode cannot reenter that nonrecursive lock. Cache invalidation,
BOM removal, UTF-16 metrics, line metadata, streaming-reader behavior and loader
errors retain the original sequence. Clear still reloads disk/zip text later.

DoRun restores every original loader by defer after compiler workers have
joined, including ordinary errors, output early returns and internal failures.
The compiler's final SourceFile state is retained; only the private ownership
wrapper is removed. Public constructors and custom loader monitor/reentrancy
contracts remain unchanged. The owner rejects test/custom runners, previously
initialized compilers, aliased files, custom observers/passes, unknown/stub
loaders and typed/saved-state lifecycles before installing anything.

The restoration list contains weak SourceFile handles. Dependency pruning calls
repartitionInputs → RebuildInputsFromChunks → InitInputsByIdMap, which clears
discarded compiler inputs from the map. A strong restoration list would retain
their raw code until DoRun returned. Weak handles restore surviving files while
allowing discarded files and their source text to collect at the original time.
They add no per-SourceFile fields or ownership backpointer. A GC regression keeps
the compiler and owner live after the actual chunk-removal/map-rebuild path,
proves the discarded SourceFile collects, then restores the surviving loader
through an internal exception.

The atomic publication covers code snapshots and owned GetCode/GetCodeReader
with Set/Clear. It does not redefine every SourceFile metadata API as an
arbitrarily concurrent interface. Ordinary CLI phase ownership still protects
line-offset computation: parsing workers receive distinct SourceFiles, buffered
diagnostics publish after joining, and function-worker source formatting already
serializes same-file lazy metadata access. Loader restoration occurs only at
the outer lifecycle boundary, never concurrently with active users.

Tests compare preload/disk/zip readers and bookkeeping with the original path,
including Unicode/BOM input, edits, cache clear, loader exceptions and reloads.
They exercise concurrent code-reader/Set/Clear snapshots, public monitor
blocking, an actual custom loader that reenters GetCodeReader, lifecycle fences,
real N2 compiler parsing, actual builtin extern-file acceptance, and stock DoRun
output/diagnostics parity. No timing improvement is claimed until the source
prototype passes normal/race gates and ordinary-CLI time/peak-RSS measurements.
