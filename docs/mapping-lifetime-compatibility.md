# Preserve source-map mapping lifetimes

This source-only compatibility patch starts at615e515b (new main abcf69a0 plus
four prior compatibility fixes). No Go build, test, race/checkptr run, compiled
layout check, Java gate or performance measurement has run for this patch.

Printer mappings exposed through public builders, CodePrinter_ToSource and
GetSourceMappings keep their original retained-pointer behavior. Public result
sequences capture fixed pointer-list bounds, remain repeatable after a complete
iteration, and preserve mapping pointee sharing through Go value copies. A
retained custom printer can append mappings without extending a prior result.
Custom factories and all public printers use owned dynamic mapping arrays; none
of those records are returned to the reuse pool. GetSourceMappings explicitly
revokes release before exporting pointers. Public Build honors source-map
options, including when a custom factory inspects its consumer.

Only Compiler.toSourceAndMappings enters the private compiler-owned builder
lane. A fresh builder with no custom factory owns its printer, and the compiler
consumes each mapping immediately into its SourceMap before the iteration
returns its slabs. A custom factory supplied even through that private helper
uses the stable public lane. An early break or callback panic does not return
slabs during unwinding. This does not promise repeatability for the private
single-consume result. The private compiler path retains new-main's chunk
iteration/pool optimization without assigning that ownership to public APIs.

Generator mappings now carry a GC-tracked immutable source-name pointer instead
of an index into the generator's source-name table. One cell is created for a
contiguous source-name run. Reset drops the generator's cached pointer without
mutating cells retained by dynamic mappings or their Go value copies. A retained
mapping can be replayed through another generator's public LineMapper. A fresh
public Mapping decodes its nil source pointer as the original empty source name.
The first16 inline mapping records are still cleared by Reset, as before; this
patch does not add a promise that those inline cells survive explicit reuse.

On amd64, two pointers plus six int32 positions suggest40 bytes, the same rounded
size as new-main's one pointer/seven-int32 record. This is unbuilt layout
arithmetic. A candidate-only test will log/check the actual size. Source-name
cells add one allocation per source-name run and remove the generator-owned
source-name slice/index cache; profiles must measure that tradeoff. No allocation,
RSS or speed gain is claimed. Java-int32 coordinates and wrapping prefix/offset
arithmetic remain unchanged. Existing FilePosition constructors and tests do not
establish a wider Go-only coordinate contract, so no speculative width expansion
is included.

Seven common regression roots use stock1fe APIs and can be overlaid without
production changes onto1fe (positive controls) and615e (negative controls):
repeat/retain/copy printer results, interrupted iteration, custom Build/build-
with-mappings/ToSource observers, concurrent public builds, fresh Mapping entry,
retained/copied/reset/other-generator entries, and Java32 source/prefix bounds.
Two separate candidate-only roots check the private ownership lane/export guard
and the source-cell representation. All are UNRUN; original tests are unchanged.
Negative controls must reproduce the concrete615e failures; an inability to
compile a positive control is not evidence of the regression.

Required admission: focused positive/negative controls; complete normal/race/
package-checkptr jscomp and sourcemap suites; existing Java source-map/print and
full parity gates; deterministic parallel compilation checks; strict six-artifact
CLI byte gates. Rebuild with the new-main shipped PGO and verify actual inlining
before alternating benchmark/RSS trials. Preserve every failure/inconclusive
result. Existing Node64/cap127 candidates and their pinned pending validation
plans remain untouched until the new base is verified.
