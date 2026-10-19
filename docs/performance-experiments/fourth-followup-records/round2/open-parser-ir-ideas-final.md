# Open frontend architecture ideas

These remain substantial, unproven designs. The experiments so far do not
establish that full parser/IR redesign is exhausted, and do not predict a win.
No new code, build, profile or benchmark accompanies this note.

## What was actually explored

The current clean-main owned-ASCII factor retains the complete original
Parser.ParseProgram → ProgramTree → IRFactory → Rhino AST pipeline. It changes
only source backing for eligible fresh types-only CLI inputs. Its current
screen was flat/slower; that does **not** measure eliminating the intermediate
parse tree or location graph. Its full legacy parsing and focused differential
passes establish backing equivalence, not equivalence of a future fused parser.

Earlier isolated direct-Rhino experiments initially supported a restricted
expression/statement subset and rejected files containing functions/comments.
Later versions covered substantially more real programs, reused original
annotation/validation helpers, and delegated unsupported syntax to original
parser fragments. Thus accepted-file coverage was not wholly direct parsing.
They still consumed the original Scanner's eagerly allocated tokens, positions
and ranges, and retained whole-file fallback. Some early results were only
component/microbenchmarks; later macro results were mixed. Parser allocation
batching, token-self compaction and last-position reuse tested narrower storage
changes. None was a complete shared-grammar value/event frontend.

## 1. One grammar, alternative parse-tree/event/IR sinks

A compiler-private sink could receive grammar production results as compact
value records/events, then construct ordinary Rhino nodes without a parallel
ProgramTree object graph. Keep the original grammar and parse-tree sink as the
public implementation/oracle. A private event sink is a plausible intermediate
step before attempting immediate Rhino construction; it can defer decisions
without retaining every generated Java-style tree object.

This is more than replacing constructors. Current grammar routines inspect
intermediate-tree kinds and payloads while resolving arrow/parenthesized cover
grammar, patterns, assignments, functions and context. Such inspection needs an
explicit representation/contract. Source comments are available to the original
IR transform after parsing the whole file; fileoverview, licenses, JSDoc claims,
closure-unaware regions and warning order cannot simply move earlier. Bounded
transactional statement/function event windows or a separate attachment stage
must preserve those dependencies. Existing validation walks may be fused only
after proving the same diagnostic chronology and resulting metadata, including
failure/KEEP_GOING paths. Retain fallback; count delegated work explicitly.

## 2. Private value tokens and lazy public positions/ranges

Source identity could live once per lexer/file, with immutable value spans
(offset/line/column endpoints) carried by private tokens/events. Public
SourcePosition/SourceRange objects would materialize only where the existing
public parser, retained comments, diagnostics or other exposed APIs require
them. Narrow coordinates need checked overflow and an original/full-width path.
Ordinary GC-tracked Rhino nodes remain unchanged; this is not an integer-handle
AST, global pool or cross-compilation cache.

Today public token/location fields are read directly, and native expressions
retain position pointers across lookahead and fragment rewinds. Merely changing
a getter, compacting nativeExpression, or recycling existing position cells
cannot remove the scanner's already-created location objects. A genuine shared
lexical kernel plus public/private adapters is required. Bridges must preserve
lookahead payloads, comment callback chronology, rewinds, UTF16 indices and source
identity. Retained comments need independent stable positions; temporary region
storage must not retain unrelated syntax, and identifier/comment strings must
not borrow large source buffers. Existing public APIs must keep their fields,
identity and lifetime behavior.

## Evidence and decision gates

Fresh profiles containing parser token/tree/position/range allocation sites
motivate investigating representation. Inclusive sites overlap; adding them
would double-count allocation, and neither cumulative allocated bytes nor a
simple object-size model predicts live memory, peak RSS or elapsed improvement.
Old corpus acceptance/component timings are not current full-CLI performance.

Any new implementation should stay behind the fresh CLI boundary and keep
public/custom/IDE paths intact. First compare complete AST properties, JSDoc,
features, UTF16 source ranges, diagnostics/order/exceptions, comments and source
maps against the original grammar across valid, ambiguous and malformed input.
Then prove actual ordinary CLI activation/optimized source bytes and remaining
fragment/fallback work, including N1/N2 and race gates. Only afterward compare
paired default-GC CLI time and peak RSS independently. No cached-worker baseline,
GC tuning, semantic shortcut or speed prediction is implied by these ideas.
