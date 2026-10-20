# Compiler-owned fixed dependency matcher

Source baseline: 9aee69c4. Internal experiment flag:
`JSGO_EXPERIMENT_NATIVE_DEPENDENCY_SCAN=1`.

The eight merged whitespace diagnostic profiles attribute 0.42/11.12 sampled
CPU seconds (3.78%) to JsFileRegexParser.ParseFile. This includes the line/comment
driver and other dependency work, so 3.78% is an upper bound, not the expected
gain from this matcher. SourceFile synchronization is a separate experiment.

The existing PrebuildDependencyInfo already scans independent stock files in a
bounded pool under num_parallel_threads, buffering diagnostics and publishing
caches in input order. This experiment adds no concurrency or runtime setting.

Only the compiler's private stock-parser construction enables a native matcher
for the fixed GOOG_PROVIDE_REQUIRE_PATTERN. Public constructors, non-stock
compiler implementations and custom virtual-dispatch parser selves retain the
original regular expression. The native path streams full-match/method/argument
byte boundaries and avoids regexp.FindAllStringSubmatchIndex's outer slice and
per-match capture-index slices. It does not build another dependency graph.
Native character work has a per-line linear budget. Pathological unmatched
nested starts switch once to the original regular expression, resuming after
the already-published match prefix. This prevents repeated failed suffix scans
from becoming quadratic and preserves callback/diagnostic encounter order.

The native recognizer preserves the historical regex grammar rather than
implementing JavaScript: only start-of-line or semicolon starts, the optional
greedy assignment prefix, its exact Unicode/ASCII classes, six method names,
ASCII Java regex whitespace, lazy arguments and the first closing parenthesis.
Semicolons and parentheses inside strings remain historical false positives.
Go Unicode classes use the same tables as the existing RE2 adapter. Newline,
carriage return, U+0085, U+2028 and U+2029 keep the original dot exclusions;
trailing Java whitespace before ')' may still include CR/LF. The original
ParseJsString validator and dependency/module-conflict publication code are
shared by both matchers.

The original buffered line driver, comment quirks, two-line retry buffer,
shortcut behavior, JSDoc annotations, module resolver, ES import/export matcher,
source lines and diagnostic UTF-16 locations remain intact. Neither the native
matcher nor this patch computes a new diagnostic column.

Source tests compare every match/capture boundary across thousands of grammar
and seeded malformed fixtures, then compare complete dependency objects,
success state, line counts and ordered diagnostic strings across comments,
module conflicts, invalid strings, CRLF, Unicode, retries and parser reuse.
Private probes prove actual native dispatch and custom-self fallback. A matched
microbenchmark reports allocations for ordinary, many-match and Unicode lines.

No tests, builds, microbenchmarks or whole CLI timings were run during the
source-only reservation. Before retention: run focused native/ported dependency
tests, race concurrent file prebuild tests, existing Java dependency oracles and
full CLI N=1/N=2 output gates with the flag; then measure fresh CLI wall time,
peak RSS and allocations against the same baseline without the flag. The
remaining line driver/ES6 matcher costs and matcher discovery overhead can make
the macro benefit much smaller than the microbenchmark result.
