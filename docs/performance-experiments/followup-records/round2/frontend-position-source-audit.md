# Frontend position/source representation audit

2026-10-19. Source-only against `22358ec4` in
`/workspace/jsgo-post-colors-candidate-kernel`. No builds, tests, profiles,
benchmarks, input inventories, or implementation changes in this audit.

## Existing work and exact current layouts

- `parser.SourcePosition` remains source pointer + Offset/Line/Column machine
  ints: 32 bytes on this 64-bit target. Its fields are exported and read directly
  throughout Parser/IRFactory/nativeParser. `SourceRange` is two pointers: 16 B.
- Current TokenBase is self interface (16 B), range pointer (8 B), int32 kind
  plus alignment (8 B): 32 B. IdentifierToken adds string16 + bool/padding8:
  56 B. These are source layout deductions, not fresh unsafe.Sizeof test results.
- Plain tokens, identifiers, positions, and ranges already use parser-owned
  slabs (16/32/64/128 elements, only current chunk rooted by the slab). Positions
  already reuse the most recent immutable token position at the same offset.
  Comments use independent positions to avoid retaining token slabs.
- `/workspace/jsgo-position-cache` commit75eb890e is exactly that last-position
  cache; it is already present. Recycling those cells would break retained
  lookahead, nativeExpression start/end, parser rewind, and fragment positions.
- `/workspace/jsgo-source-position-cache` is SourceFile lazy line-cache
  concurrency/loader locking work; it does not compact parser SourcePosition.
- `/workspace/jsgo-compact-parser` tests TokenBase/ParseTreeBase32->16 by removing
  generated self-interface storage with concrete reconstruction. It preserves
  the SourcePosition32/SourceRange16 layouts. The old measured experiment was
  rejected: typed6.155 vs6.090s, bundle23.543 vs23.418s, white1.517 vs1.469s
  (`/workspace/jsgo-results/progress.md`). Do not repeat that as a new design.
- `/workspace/jsgo-native-parser-next` extends direct AST grammar/comments;
  it still consumes public Scanner/Token/location pointer graphs. Current
  nativeExpression retains two SourcePosition pointers. Changing those fields
  to values alone cannot remove scanner allocations already made eagerly.

## What the existing white capture establishes

The retained `frontend-allocation-audit.{md,json}` reports whole400.90 MiB,
frontend124.73, Scanner ancestry52.68, position/range ancestry40.64. These are
inclusive overlapping ancestry categories. Actual position slab leaf24.09 plus
range slab leaf8.01 totals32.10 MiB, not40.64 MiB of independently removable
storage. SourceFile UTF16 conversion11.38 MiB divides8.32 under native attempts
and3.06 under original parses. The capture predates newer extern fragment work;
these are historical measurements, not current constructor counts or accepted
byte coverage. Scanner CPU6/139 samples and positions3/139 samples are coarse;
64-PC truncated ancestry cannot prove a hard ceiling or a large speedup.

## Genuine position compaction needs a value-token boundary

The useful design is source identity once per lexer + private span of six
int32 values (24 B, no pointers), with kind/flags/payload in a value token.
Public SourcePosition/Range materialize only for retained comments, diagnostics,
and original grammar fragments. This can remove repeated source pointers,
location pointer chasing and much of the position/range allocation, unlike
another pool or self-interface compaction.

It is not a transparent SourcePosition field edit: exported machine-int fields,
Source pointer identity, and GetTokenBase.Location are current API contracts.
A private 32-bit representation must guard file length, line and column overflow
and use original scanner for out-of-range inputs. No narrowing of public fields.

A feasible architecture extracts one shared lexical cursor/kernel that returns
private value-token records. A public Scanner adapter materializes original
Token/location objects; the native sink consumes values. Native originalFragment
currently rewinds the SAME Scanner to current token start, invokes original
Parser, then resumes at actual last-consumed location. A new bridge must transfer
lookahead (including six-token extern lookahead and arbitrarily longer head
checks), cursor/line state, token payloads, and comment callback chronology
exactly. Original fragment lookahead can itself scan later comments; conversion
must not publish/record them twice. Native expression spans are copied values
that survive these transitions. A generated lexical kernel is preferable to
maintaining two independent grammars. Returning a wrapper which eagerly creates
public ranges simply retains today's allocation cost.

This is a substantial scanner/API-private redesign (~1500-line lexical source
plus adapters, native helper migration and oracle tests), not an appropriate
unvalidated small patch. It requires scanner-level token/value/range/error/comment
oracles, fallback AST/diagnostic/UTF16/source-map oracles, retained-comment GC
checks, and real CLI activation/bytes inventory before timing.

## Smaller independent first experiment: ASCII source backing

A separate source-only finding narrows the old audit's complexity estimate.
Scanner has only eight direct `s.contents[...]` uses: four substring conversions,
two comment-marker character reads, NextChar, and PeekCharAt. LineNumberScanner
has two direct character reads in AdvanceLine. Only two parsing-side direct
ContentsUTF16 uses exist: trailing non-JSDoc whitespace and diagnostic fallback
byte inventory. A private ASCII source backing can therefore reuse the EXACT
existing Scanner and original-fragment grammar via character/range helpers;
it does not need a second lexer or token API.

Possible bounded design (not implemented):
1. Leave public NewSourceFile/ContentsUTF16/Scanner behavior unchanged. Native
   types-only attempt may create a private source file after full ASCII check;
   retain original string and omit UTF16 backing. Any byte>=128 uses original
   source constructor. This also preserves invalid UTF8/lone-surrogate decoding
   on the original path, without a second Unicode implementation.
2. Scanner/LineNumberScanner constructors select private ASCII backing. Their
   small character helpers widen byte to uint16, keeping all lexical branches,
   UTF16 offsets and line/CRLF semantics unchanged. AdvanceLine chooses its
   backing once outside the loop. Escaped identifiers/string Unicode sequences
   continue through existing decoders; only SOURCE characters are ASCII.
3. Token/comment substring extraction uses strings.Clone, not borrowed slices:
   current StringFromUTF16 creates owned strings, and changing that would pin
   entire sources through retained tokens/AST/comments. Existing global
   canonical interning behavior remains unchanged.
4. Original fragments reuse the same dual-backed Scanner. A whole-file failed
   native attempt invokes the unchanged original entry point with its normal
   UTF16 source. No cached parse/decoded source is kept across compilations.
5. Audit private source escape: native runner rejects all JsDoc description/
   INCLUDE_ALL_COMMENTS modes, returns empty public token comments/ranges, and
   regular AST nodes retain the ORIGINAL rhino.StaticSourceFile, not this parser
   source. NonJSDocComment is the only rhino field holding parser positions;
   it is only returned for INCLUDE_ALL_COMMENTS, excluded here. Keep those guards
   and add an explicit retained-public-position test before permitting nil
   ContentsUTF16 on any private source. Public standalone parser APIs must never
   receive or expose the private alternate backing.

This may remove only the ASCII subset of the historical native-attempt8.32 MiB
UTF16 conversions (~2.08% whole allocated bytes), not all11.38 MiB. It keeps all
position/range/token slabs. An ASCII pre-scan and character helper branches add
CPU work; real files may contain non-ASCII comments and fall back, so coverage
must be counted, not assumed. It is independently measurable and smaller risk
than value tokens, but no speedup follows from the allocation calculation.

## Recommended sequence

If a bounded first prototype is desired, compare private ASCII source backing
against the exact frozen CLI binary with default GC, real ASCII/native accepted
and rejected byte counters, all fragment/comment/UTF16/error oracles. Do not mix
it with token self compaction. Decide afterward whether the larger value-token
bridge is justified. The latter is the only proposal here that targets the
position/range slabs; ASCII backing targets source decoding instead.
