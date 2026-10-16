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
| `pkgscc.py` | — | the Java package SCCs that define the Go package layout (PORTING.md §2) |

The Go side of each oracle (e.g. the Go `rhino.Node` dump) must reproduce the
Java format byte for byte; it is produced by ported code (`Node.toString`,
`JSDocInfo.toStringVerbose`), not by a Go re-implementation.
