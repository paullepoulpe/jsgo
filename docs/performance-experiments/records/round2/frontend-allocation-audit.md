# Remaining frontend allocation representation audit

This analyzes existing default-GC structural captures after the timing
reservation ended on 2026-10-19. No new workload or profile was captured.
Heap allocation columns come from `profile.live.pprof` after forced GC and
retain 512 PCs. CPU profiles retain 64 PCs and the white capture has only 139
10-ms samples, so ancestry and percentage estimates are coarse. Inclusive
categories below overlap.

| White category | Sampled allocated MiB | Whole allocation share | Sampled CPU seconds | Whole CPU share |
|---|---:|---:|---:|---:|
| Whole task | 400.90 | 100% | 1.39 | 100% |
| Frontend ancestry | 124.73 | 31.11% | 0.23 | 16.55% |
| Scanner ancestry | 52.68 | 13.14% | 0.06 | 4.32% |
| Position/range ancestry | 40.64 | 10.14% | 0.03 | 2.16% |
| SourceFile UTF-16 conversion | 11.38 | 2.84% | 0.02 | 1.44% |
| Scanner UTF-16 to string conversion | 2.00 | 0.50% | 0.01 | 0.72% |

The actual SourcePosition slab leaf accounts for 24.09 MiB and SourceRange slab
leaf for 8.01 MiB on white. Position/range ancestry additionally includes other
allocations called below those functions, so 40.64 MiB is not a standalone
slab-size saving. Bundle's corresponding slab leafs are much larger, but their
absolute size cannot be substituted for white's opportunity.

White frontend allocations divide into 75.55 MiB direct native attempts,
5.02 MiB native original-grammar fragments, 42.65 MiB whole-original parses,
and 1.50 MiB other paths. Attempts may later reject, so these do not measure
accepted-source byte percentages. Heap stacks attribute SourceFile UTF-16
creation to 8.32 MiB under native attempts and 3.06 MiB under whole-original
fallback. Sharing a decoded SourceFile across rejected/native and original
attempts cannot eliminate the whole 11.38 MiB: accepted attempts still need
their source representation. The current extern-native followup is newer than
this capture and reduces some fallbacks further.

## Source findings

`parser.SourcePosition` is 32 bytes: one source pointer and three machine ints.
`SourceRange` is two position pointers, 16 bytes. Scanner tokens have a SourceRange
pointer and a self interface. Identifier tokens add their string and metadata.
The native direct parser uses these same public scanner/token objects and keeps
start/end position pointers in its nativeExpression values. Original grammar
fragments share the Scanner and rewind it to preserve lookahead and comments.

The scanner already batches tokens, identifiers, positions and ranges, and
already reuses its last immutable position at the same offset. Slab only retains
its current block; it does not keep a list of every prior block. Repeating those
batching/cache ideas would not be a new representation change.

SourceFile already converts its text to UTF-16 once per parser source instance.
StringToUTF16 and StringFromUTF16 already have ASCII allocation fast paths;
the latter allocates one byte backing array and exposes it as an immutable
string. Keyword identification still briefly builds a string, but all scanner
UTF16-to-string conversion together is only 2 MiB / one CPU sample on white.
That does not support a large new keyword helper as the next architectural step.

## Substantive design options and limits

The strongest frontend-specific storage redesign is a private value-token
stream: compact offset/line/column records with the source identity stored once
per scanner. A six-int32 source span is 24 bytes versus two 32-byte positions
plus their 16-byte range, before token objects. Token kinds and flags can be
concrete fields rather than a self interface. Real SourcePosition/SourceRange
objects would materialize only at public/original-grammar boundaries and retained
comment/diagnostic positions. This removes pointer scanning and repeated metadata,
rather than adding another pool.

It requires a genuine scanner API/kernel change. Public Token.GetTokenBase
exposes Location as a pointer graph, and the original Parser stores *Scanner,
not an interchangeable lexer interface. Returning wrappers that eagerly build
those public locations would retain the current allocation cost. Recycling
cells behind nativeExpression.start/end is unsafe: nested parsing, lookahead,
rewinds and fragments can still refer to old positions. SourcePosition's public
machine-int fields also prevent a transparent private 32-bit backing substitute.
Changing all public field types would introduce unrelated Go API changes.

A private ASCII source view could avoid a full UTF-16 copy while retaining Java
offsets for fully ASCII files. The current scanner and line scanner index their
UTF-16 slices directly throughout, and IRFactory reads that slice for comment
handling. A real dual-source character kernel is needed, with an original UTF-16
fallback for Unicode, CRLF, escaped/lone surrogates, comments, templates and
fragment rewinds. Simply borrowing string substrings can retain whole input
buffers through published identifiers; ownership must be explicit. The measured
source-conversion fraction is small, so this is weaker than compact value tokens.

The measured direct scanner CPU share is about 4%, not evidence that deleting
scanner allocations alone closes the roughly 25% elapsed gap to warmed Java.
Indirect default-GC effects might add value, but require measurement. The white
capture's unlabeled CPU includes substantial GC scanning, which this pointer
reduction would target; that is an experiment hypothesis, not a predicted win.
Given the bridge complexity and coarse CPU evidence, no new frontend prototype
was implemented in this audit.

Raw results: `frontend-allocation-audit.json`; reader:
`frontend_alloc_audit.py`. Existing heap totals/ancestry and source inspection
are reproducible independently of any future capture or GC tuning.
