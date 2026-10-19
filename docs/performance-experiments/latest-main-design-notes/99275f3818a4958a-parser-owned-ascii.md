# Stock-parser ASCII source backing

This clean-main experiment starts at ff06ac15. It does not port the previous
native/direct parser or its fragment grammar, caches, node allocators or flags.
The complete existing Parser.ParseProgram and IRFactory transform remain the
only grammar/AST pipeline. No performance claim is made.

The new compiler-owned parser entry can initialize that same full Parser with
an ASCII-backed scanner. It first verifies every source byte is below128;
otherwise it uses the original SourceFile/Parser constructors. A private source
stores the original string without a UTF16 copy. Scanner character helpers widen
bytes to uint16, LineNumberScanner keeps the same CR/LF/CRLF accounting, and all
token/comment substrings clone ownership as the old UTF16 conversion did.
Unicode escapes in ASCII source still use the existing decoders; raw Unicode,
malformed UTF8 and lone-surrogate encodings keep the original UTF16 path.

Public SourceFile, Scanner, LineNumberScanner, NewParser and ParserRunner_Parse
constructors/behavior remain unchanged. A public nil ContentsUTF16 does not
select the private backing. The new ParserRunner_ParseOwnedASCII entry selects
it only in TYPES_ONLY mode and excludes TypeScript; all description/comment
preservation modes fall back before creating a private source.

TYPES_ONLY result comments and top-level ranges are empty. IRFactory nodes use
the caller's original rhino.StaticSourceFile. JSDoc markers contain integer
coordinates/AST nodes. The only Rhino payload retaining parser.SourcePosition is
NonJSDocComment, which IRFactory returns only in INCLUDE_ALL_COMMENTS. Reporter
forwarding publishes message/name/line/column, never the private position. Thus
no private parser SourceFile may escape through accepted/partial/error results.
An optional caller-owned observation records ASCII-backed bytes; its private
weak SourceFile handle lets tests keep returned ASTs and exceptions alive while
proving parser-source collection. No global observation or source registry is
introduced.

Prepared differential tests compare the original public runner to the owned
entry across full stock grammar, JSDoc/licenses, module metadata, closure-unaware
shadows, lexical/syntax errors, strict/sloppy and STOP/KEEP modes, UTF16 source
ranges, features, comments, sourceMappingURL, diagnostic/exception order and
weak lifetime. Scanner tests separately cover exact tokens/comments/locations,
retained rewind positions, regex/templates, Unicode rejection, original public
constructor behavior and owned string backing. Full internal/parser normal tests and 103 original/owned parsing oracle cases
pass on clean main. Focused race and actual DoRun checks are tracked separately
below.

The separate shared fresh CLI lease foundation is65b4423f (cherry-picked here
as12e59f6f). CompilerInput.JsAst.parse and parallel prebuild both call the same
private helper, which requires nativeCLIOwnsSource and the independent
JSGO_EXPERIMENT_PARSER_ASCII_SOURCE=1 flag. No native-parser or source-lock flag is
needed. Public Compile lacks this provenance and retains the original path.
The lease contains no strong file registry/backpointer, leaves only immutable
inert per-file tags, and clears the Compiler tag after all DoRun jobs join.
Actual stock DoRun N1/N2 tests attach a private numeric probe through the exact
stock runner CreateCompiler and compare JS/AST/maps/diagnostics/exit status for
success and syntax failures. Public Compile fallback is separately asserted.
No custom factories or eligibility-only assertions substitute for activation.

Only UTF16 source conversion is targeted. Token/position/range slabs remain.
The ASCII scan and added backing branches may cost more CPU than the avoided
copy saves; real source byte coverage and independent default-GC CLI time/RSS
measurements are required. Historical older frontend captures are not current
clean-main allocation counts.

## Focused validation

Go1.27.2, GOMAXPROCS1, GOGC25 only to bound build memory, `-p=1 -pgo=off`;
no benchmark or GC-tuned performance result:

- Full internal/parser normal15 groups PASS0.027s and race PASS.
- Standalone parsing normal103 complete original/owned comparisons PASS0.028s;
  race PASS1.073s. These include96 grammar/strict/run-mode cases, four Unicode
  sources and three public comment/description preservation modes.
- Weak private-source lifetime, string ownership, closure-unaware shadows and
  lexical/syntax error/exception oracles pass. No production fixes were needed.

Logs are `/workspace/jsgo-results/round2/main-parser-ascii/{parser-normal,
parser-race,parsing-normal,parsing-race}.log`. The large legacy internal/parsing
suite and jscomp actual DoRun/public Compile tests await their scheduled slots.
Source peer review found no parser/CLI glue/probe blocker. The shared lease is
only provenance: this full parser keeps original reporter/callback chronology,
so old speculative-native callback guard machinery is neither needed nor ported.

The full INTERNAL legacy parsing suite also passes:1476 top-level tests and six
expected unconfigured corpus-oracle skips, both through original dispatch and
through an owned-ASCII test-only callsite overlay. The SAME binary ran both
modes. Its numeric dispatch witness records2919 legacy parser calls in each;
normal mode backs zero, owned mode backs2404 calls/783355 source bytes. Legacy
expectations/cases are unchanged. New103 full differential oracles retain their
untouched original baseline. Only external parsing_test file
js_doc_info_parser_position_test.go is excluded because it imports jscomp and
would require an unscheduled giant compiler build; that one position test and
the actual DoRun jscomp tests remain pending the central tranche.

Reproduction artifacts:
`/workspace/jsgo-results/round2/main-parser-ascii/generate_internal_overlay.py`,
`legacy-build-command.json`, and `{legacy-build,legacy-normal,legacy-owned}.log`.
The explicit-file build/run CWD is internal/parsing. Generated overlays and the
test binary live in /tmp; they are never product source or benchmark variants.
