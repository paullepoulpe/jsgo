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

The initial grammar covers comment-free ES5 expression/data programs:

- Explicitly semicolon-terminated `var`, expression and empty statements.
- Ordinary unary, arithmetic, shift, relational, equality, bitwise and logical
  expressions; simple assignment and conditional expressions.
- Calls, ordinary property/index access, parentheses, arrays without holes and
  objects with explicit identifier/string/number keys.
- Decimal numbers, unescaped strings, ordinary identifiers, `this`, booleans
  and null. UTF-16 source positions use the existing scanner.

Functions, blocks, modules, ASI, patterns, trailing commas, spread, optional
chains, exotic numeric syntax, escaped strings and reserved/private identifiers
remain on the original path. TypeScript sources and configurations requesting
comment descriptions/statement ranges also use the original path. Native
recursion is bounded; deeper input falls back.

Fallback is for the entire file. All scanner and IR diagnostics are private
until acceptance; any warning or error rejects the attempt, and the original
parser then emits its complete diagnostic sequence. Any comment also rejects
an attempt, including comments encountered late in the file. Therefore licenses,
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
bundles containing functions/modules/comments. Production expansion should add
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

Five alternating original/native pairs on the supported 128-entry data program,
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
