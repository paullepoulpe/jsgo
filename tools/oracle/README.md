# Oracle tooling

The pinned Java Closure Compiler (v20260415) is the specification. These tools
extract ground truth from it, layer by layer, for any JavaScript input, so the
Go port can be diffed against Java on arbitrarily large corpora instead of a
fixed set of hand-picked fixtures.

```sh
tools/oracle/fetch.sh    # .cache/closure-compiler.jar (sha256-checked) + .cache/closure-src/
tools/oracle/build.sh    # compiles java/*.java into .cache/oracle-build/
tools/oracle/corpus.sh   # pinned real-world corpus (corpus.txt) into .cache/corpus/
CP=.cache/closure-compiler.jar:.cache/oracle-build
```

| Tool | Layer | What it emits |
|---|---|---|
| `java/AstDump.java` | parser (`ParserRunner` → `IRFactory` → `rhino.Node`) | per file: parse diagnostics + the full `Node` tree, one node per line via `Node.toString(true, true, false)`, plus `JSDocInfo.toStringVerbose()` for nodes that carry JSDoc. `java -cp $CP AstDump <outdir> <root> <rel-file>...` |
| `java/ParsingDump.java` | parsing (`ParserRunner.parse`: Parser → IRFactory → `rhino.Node`, JsDocInfoParser) | per file: the `parsing.Config` the compiler uses, every message reported to the rhino `ErrorReporter` (with line/offset), the `FeatureSet`, source map URL, comments, and the `Node` tree (`Node.toString(true, true, false)` per node + `JSDocInfo.toStringVerbose()`), every line escaped to printable ASCII. This is the exact boundary of `internal/parsing`. `java -Xss1g -cp $CP ParsingDump [--gzip] [--lang=..] [--jsdoc] <outdir> <root> <rel-file>...\|@<listfile>` |
| `java/CliBatch.java` | whole compiler (CLI) | the exact `java -jar closure-compiler.jar ...` output and exit code, for many invocations in one JVM. Jobs on stdin: `<outfile>\t<arg>\t<arg>...`. Byte-identical to the per-process oracles in `testdata/oracle/`. |
| `java/ParseTreeDump.java` | parser (`parsing.parser.Parser`, before IRFactory) | per file: the parser config, every parse error/warning with its position, the feature set, the source map URL, the comment list and the whole `ProgramTree`, dumped field by field by reflection. Configured exactly as `ParserRunner.parse` does for `--language_in=ECMASCRIPT_NEXT`. `java -cp $CP ParseTreeDump [--gzip] [--lang=<LanguageMode>] <outdir> <root> <rel-file>...\|@<listfile>`; the Go side is `internal/parser/parsetreedump_test.go` (`PTREE_CORPUS`, `PTREE_ORACLE_DIR`, `PTREE_FILES`, `PTREE_DIFF_DIR`). |
| `java/ParserEnumDump.java` | parser | every parser enum (ordinal, name, toString) and static `FeatureSet` constant → `internal/parser/testdata/enums.txt` |
| `java/JavaCharacterTable.java` | jutil | `java.lang.Character` classification of all 65536 chars → `internal/jutil/testdata/character.tsv` |
| `java/JavaCaseTable.java` | jutil | `Character.toUpperCase/toLowerCase(char)` for every char whose case mapping is not the identity → `internal/jutil/testdata/character_case.tsv` |
| `java/GsonVectors.java` | jutil (Gson) | `Gson.fromJson` (lenient `JsonReader` → `JsonElement` tree → `toString` via `JsonWriter`, with every exception's `toString`) on ~1,800 JSON inputs (edge cases, 1024-char buffer boundaries, random token soups and mutated documents), the `getAs*` accessors, and `LazilyParsedNumber.intValue/longValue` → `internal/jutil/testdata/gson.tsv.gz`. `java -cp $CP GsonVectors <out.tsv.gz>` |
| `corpus.txt`, `corpus.sh` | — | 23 pinned npm packages (14,690 .js/.mjs/.cjs files: Closure Library, test262 parser tests, three.js, core-js, TypeScript, React, Angular, …), sha512-verified against the registry |
| `pkgscc.py` | — | the Java package SCCs that define the Go package layout (PORTING.md §2) |

The Go side of each oracle (e.g. the Go `rhino.Node` dump) must reproduce the
Java format byte for byte; it is produced by ported code (`Node.toString`,
`JSDocInfo.toStringVerbose`), not by a Go re-implementation.

Generating the AST oracle for the whole corpus (two JVMs, ~5 minutes, 1.4 GB):

```sh
(cd .cache/corpus && find . -type f | sed 's|^./||' | sort) > .cache/corpus-files.txt
xargs -a .cache/corpus-files.txt -d '\n' java -Xmx3g -Xss16m -cp $CP AstDump .cache/oracle/ast/corpus .cache/corpus
```

Note for Go-side dumpers: Java's `PrintStream` encodes a lone surrogate as `?`,
so WTF-8 strings must be written with lone surrogates replaced by `?`.
