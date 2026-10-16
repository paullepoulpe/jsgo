# jsgo

A Go port of [Google Closure Compiler](https://github.com/google/closure-compiler),
aiming to be a drop-in replacement: same flags, same output, same diagnostics.

The port is a whole-file, 1:1 translation of the pinned Java release
(**v20260415**), built bottom-up and verified layer by layer against the
Java compiler itself on large corpora.

- [PLAN.md](PLAN.md) — feasibility assessment, strategy and roadmap.
- [PORTING.md](PORTING.md) — the translation contract every ported file follows.
- [tools/oracle/](tools/oracle/) — tools that extract ground truth from the Java compiler.
- [legacy/](legacy/) — the earlier type-checker-only attempt, frozen for reference.

```sh
make            # build, vet, test
make oracle-tools   # fetch the pinned jar/sources and build the Java oracle tools (needs a JDK)
```
