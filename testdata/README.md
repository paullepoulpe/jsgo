# testdata

Closure type-checking fixtures with their Java oracles. They are kept as
the acceptance set for the type checker (milestone M7 in `../PLAN.md`).

- `fixtures/tier0..tier3/` — 44 small hand-written programs.
- `fixtures/extracted/` — 1,444 programs extracted from Closure's
  `TypeCheckTest.java` (v20260415), one per test method that has a single
  plain `addSource` string.
- `oracle/`, `oracle-strict/` — the pinned Java compiler's CLI output and
  exit code for each fixture, under
  `--checks_only --jscomp_error=checkTypes --warning_level=VERBOSE --language_in=ECMASCRIPT_NEXT`,
  plus `--jscomp_warning=strictCheckTypes` for `oracle-strict/`.

To regenerate an oracle, feed `tools/oracle/java/CliBatch` one job per
line: `<outfile>\t<flags…>\t--js\t<fixture>`. Its output is byte-identical
to these files. The extraction tool and the original per-process scripts
belonged to the deleted first attempt; they are in git history at
`c05fc09` (`tools/extract_fixtures`, `scripts/gen-oracles.sh`).
