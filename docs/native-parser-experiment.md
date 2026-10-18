# Direct AST parser experiment

This branch contains an opt-in architectural prototype. Set
`JSGO_EXPERIMENT_NATIVE_PARSER=1` when starting the CLI to try it.
The default compiler continues to use the original parser. This experiment has
not established an end-to-end performance gain.

The native path consumes the existing Scanner's tokens and constructs final
Rhino nodes directly. It creates neither a ProgramTree nor intermediate
expression/statement tree objects. Expression frames are small values containing
a final node and its source-position bounds. Existing IRFactory helpers still
provide literal normalization, source metadata, property keys, call flags,
directives, file-level setup and final validation.

The extended grammar covers ordinary ES5 expression/function programs:

- Explicitly semicolon-terminated `var`, expression, empty, return and throw
  statements; blocks, if/else and while statements.
- Plain function declarations/expressions with simple named parameters.
- Ordinary unary/update, arithmetic, shift, relational, equality, bitwise and
  logical expressions; simple assignment and conditional expressions.
- Calls, ordinary construction, property/index access, parentheses, arrays
  without holes and objects with explicit identifier/string/number keys.
- Decimal numbers, unescaped strings, ordinary identifiers, `this`, booleans
  and null. UTF-16 source positions use the existing scanner.
- Ordinary line/block comments in TYPES_ONLY mode, where the original pipeline
  emits no attached or returned ordinary-comment metadata. sourceMappingURL
  extraction uses the original Java-trim and exact prefix rules.

JSDoc/important comments, modules, ASI, patterns, trailing commas, spread,
optional chains, exotic numeric syntax, escaped strings and reserved/private
identifiers remain on the original path. Async/generator/arrow functions,
parameter defaults/rest/patterns, legacy unbraced function declarations, for
loops and try/catch are also unsupported. Binding names eval/arguments remain
conservatively rejected; references and direct eval calls use the original
call-annotation helper. TypeScript sources and configurations requesting comment
descriptions/statement ranges use the original path. Native recursion is bounded;
deeper input falls back.

Fallback is for the entire file. All scanner and IR diagnostics are private
until acceptance; any warning or error rejects the attempt, and the original
parser then emits its complete diagnostic sequence. Any JSDoc or important comment rejects
an attempt, including annotations encountered late in the file. Therefore licenses,
JSDoc, fileoverview, closure-unaware annotations and source-map URL comments
retain the original whole-file discovery/attachment behavior. The prototype
intentionally does not attempt to approximate those rules with an incomplete
stream of comments.

This makes the boundary explicit and differential testing straightforward, but
unsupported files pay for speculative work before fallback. Speculative nodes
also currently use IRFactory's shared weak allocation chunks. A failed large
prefix can share a chunk with a retained AST, extending the aborted graph's
lifetime; a production sink should own its allocation group per parse, including
under concurrent parsing. The parser-only benchmark starts an allocation group
per iteration for both variants, matching a normal Compiler job. In particular,
this prototype is not ready to replace the complete grammar or demonstrate a win on
general bundles containing unsupported functions/modules/annotations. Production expansion should add
full grammar support and whole-file comment discovery or a deferred attachment
stage, plus preserve parser/IR diagnostic ordering without running the original
parser in accepted cases.

Run `go test ./internal/parsing -run NativeParser` for AST/property/source-range
and diagnostic comparisons against the original pipeline. The generated tests
exercise many expression combinations. `BenchmarkNativeParserDataProgram`
compares the same data program through both implementations and reports
allocations. A useful parser microbenchmark win would still need full CLI
byte comparisons, source-map checks, real coverage counts and end-to-end timing
and peak-memory measurements before retention.

A prior whitespace CPU profile attributed about 16% directly to parsing plus IR
conversion. Removing that directly attributed work entirely corresponds to about
1.19x speedup on that particular profile; indirect garbage-collection work is
not included in that estimate. This prototype removes only part of front-end
work and accepts only a subset of files, so it cannot by itself establish broad
parity with Java.

## Component measurement

Five alternating original/native pairs on the original supported 128-entry data program,
with Go 1.27.2, two pinned CPUs, GOMAXPROCS=2 and default GC:

| Measure | Original parser + IR | Direct AST | Original / direct |
| --- | ---: | ---: | ---: |
| Median time | 1.752 ms | 0.871 ms | 2.01x |
| Allocated bytes/op | 1,093,121 | 829,637 | 1.32x |
| Allocations/op | 7,523 | 4,026 | 1.87x |
| Median benchmark-process peak RSS | 19,596 KiB | 18,480 KiB | 1.06x |

The RSS row is the peak of the entire 200 ms adaptive benchmark process,
including calibration. Allocated bytes/op is the more useful allocation measure
for this component experiment; neither number is a full CLI compilation's peak.
Both variants use the same test binary and no PGO profile. The final pair slowed
for both variants; all five pairs are retained in the raw results. The other
four pairs showed approximately 1.96–2.05x time ratios.

The full parsing test suite passed with the experiment disabled and enabled,
using the same compiled test binary. Dedicated generated-expression and fallback
comparisons passed as well. These results establish a promising direction for
accepted syntax. They do not establish a gain on general projects or a win over
Java. See [raw samples and build provenance](native-parser-results.json).

## Real-workload extension

Source inspection found 897 core-js internals/modules files (889,310 bytes) in
the downlevel benchmark workspace: 651 contain plain functions, 794 ordinary
line comments, 176 block comments, and only one JSDoc opener. This motivated the
function/statement/ordinary-comment extension rather than another data-only
specialization. These syntax observations motivated the implementation; coverage was then measured
separately from speed.

`TestNativeParserCoreJSWorkloadCoverage`, enabled with
`NATIVE_PARSER_CORPUS_ROOT=<benchmark workspace>`, tries every internals/modules
file and compares the complete AST/source metadata and diagnostics of each
accepted file against the original pipeline. It reports both file and byte
coverage. The downloaded inputs remain outside the repository; the production
parser has no workload-specific paths or names.

With ES_NEXT/strict configuration, the extended sink accepted 671 of 897 files
(74.8%), accounting for 358,983 of 889,310 source bytes (40.4%). Every accepted
file matched the original complete AST dump, node/source properties, feature
set, comments/statement ranges, source-map URL and ordered diagnostics. The
focused suite also compares original exception/cause chains when malformed
input falls back. This establishes syntax coverage and equivalence on that
corpus, not a predicted time improvement: parsing cost is not proportional to
file count or source bytes, and rejected files incur speculative overhead.
The full original `internal/parsing` suite also passed with the extension
explicitly disabled and enabled, using one compiled test binary for both runs.
End-to-end timing is required separately from the earlier data-program results.
