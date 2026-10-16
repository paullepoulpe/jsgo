# Oracle tooling

The pinned Java Closure Compiler (v20260415) is the specification. These tools
extract ground truth from it, layer by layer, for any JavaScript input, so the
Go port can be diffed against Java on arbitrarily large corpora instead of a
fixed set of hand-picked fixtures.

```sh
tools/oracle/fetch.sh    # .cache/closure-compiler.jar (sha256-checked) + .cache/closure-src/
tools/oracle/build.sh    # compiles java/*.java into .cache/oracle-build/
CP=.cache/closure-compiler.jar:.cache/oracle-build
```

| Tool | Layer | What it emits |
|---|---|---|
| `java/AstDump.java` | parser (`ParserRunner` → `IRFactory` → `rhino.Node`) | per file: parse diagnostics + the full `Node` tree, one node per line via `Node.toString(true, true, false)`, plus `JSDocInfo.toStringVerbose()` for nodes that carry JSDoc. `java -cp $CP AstDump <outdir> <root> <rel-file>...` |
| `java/CliBatch.java` | whole compiler (CLI) | the exact `java -jar closure-compiler.jar ...` output and exit code, for many invocations in one JVM. Jobs on stdin: `<outfile>\t<arg>\t<arg>...`. Byte-identical to the per-process oracles in `testdata/oracle/`. |
| `java/ParseTreeDump.java` | parser (`parsing.parser.Parser`, before IRFactory) | per file: the parser config, every parse error/warning with its position, the feature set, the source map URL, the comment list and the whole `ProgramTree`, dumped field by field by reflection. Configured exactly as `ParserRunner.parse` does for `--language_in=ECMASCRIPT_NEXT`. `java -cp $CP ParseTreeDump [--gzip] [--lang=<LanguageMode>] <outdir> <root> <rel-file>...\|@<listfile>`; the Go side is `internal/parser/parsetreedump_test.go` (`PTREE_CORPUS`, `PTREE_ORACLE_DIR`, `PTREE_FILES`, `PTREE_DIFF_DIR`). |
| `java/ParserEnumDump.java` | parser | every parser enum (ordinal, name, toString) and static `FeatureSet` constant → `internal/parser/testdata/enums.txt` |
| `java/JavaCharacterTable.java` | jutil | `java.lang.Character` classification of all 65536 chars → `internal/jutil/testdata/character.tsv` |
| `pkgscc.py` | — | the Java package SCCs that define the Go package layout (PORTING.md §2) |

The Go side of each oracle (e.g. the Go `rhino.Node` dump) must reproduce the
Java format byte for byte; it is produced by ported code (`Node.toString`,
`JSDocInfo.toStringVerbose`), not by a Go re-implementation.
