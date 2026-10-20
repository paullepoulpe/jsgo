# Contextual rename inverter no-change proof

`JSGO_EXPERIMENT_INVERTER_NOOP=1` enables a private stock-pass wrapper. It scans
the current JS tree without scopes and skips ContextualRenameInverter only
when no AST name could pass the inverter's original rename predicate. Public
constructors retain their original traversal.

The exact predicate is **last `$jscomp$` separator at an index greater than
zero**. It does not require a numeric or nonempty suffix. Thus `$jscomp$temp`
alone cannot be inverted, but `x$jscomp$`, `x$jscomp$words` and
`$jscomp$a$jscomp$tail` can be considered. For this ASCII separator, UTF-8 byte
index positivity equals Java UTF-16 index positivity; the scan never slices a
name at that byte index and avoids converting every Unicode name to UTF-16.

Explicit untyped declarations originate in NAME or IMPORT_STAR nodes. Implicit
arguments/this/super/exports names contain no separator; implicit goog namespace
declarations are global, and the inverter never renames global variables.
Scanning every NAME/IMPORT_STAR, including references and globals, is therefore
conservative. One possible match retains the entire original scoped traversal,
including conflict resolution and original change-report order. Property names,
comments and strings do not cause false matches.

This proof is made immediately before the pass; it retains no fact across a
mutation, pass or compilation. The scan follows ordinary child links and stops
at the first possible match. It does not allocate scopes, maps or name slices.
The private wrapper requires the exact Compiler's current JS/extern roots,
stock convention/observers/redeclaration handling, no custom passes, and no
debugging/tracing. Custom behavior retains the original callback schedule.

Both stock callsites explicitly mark the AST unnormalized before invoking this
pass. Requiring a NORMALIZED lifecycle stage would make the experiment always
fall back. Tests therefore exercise the actual RAW-stage boundary after
normalization and prove actual skip/full-traversal dispatch, rather than merely
enabling the flag.

Differential fixtures compare output, AST, diagnostics, stamps and change order
with the original inverter. Predicate tests cover NAME/IMPORT_STAR, empty and
nonnumeric suffixes, repeated separators, Unicode and WTF-8 spelling. Additional
fixtures preserve custom redeclaration callback order and reject unsupported
observers, copied conventions, custom passes and foreign roots. A caller-owned
optional probe records scanned nodes/names and actual skip/full-traversal counts;
it is nil in timed CLI runs. Set `JSGO_INVERTER_COVERAGE_INPUT=/path/to/input.js`
for the optional non-timed real-input fixture.

The ceiling is the entire inverter traversal and its temporary scope/conflict
allocations on trees without a possible rename. A tree with even one candidate
gets a short preliminary scan followed by the original work. Normalized bundles
may all take that fallback, so actual coverage must precede any performance
claim. No timing or memory improvement is claimed by this source prototype.
