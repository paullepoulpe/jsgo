# Next-candidate activation and validation

Source **baf328cba812a7ebadd2eb3f9912107382805e8a**. All five fresh stock DoRun diagnostics exited 0 and closed their CLI ownership lease. These counters come from the actual pipeline without forcing admission.

| Task | Native gather CFG nodes / attempted | Original fallback nodes | Captured-validation hits | Lazy-use admitted / attempted |
| --- | ---: | ---: | ---: | ---: |
| typed-app | 19,591 / 19,591 | 0 | 2,558 | 55 / 55 |
| bundle | 300,227 / 300,231 | 4 | 0 | 795 / 795 |
| downlevel-whitespace | 0 / 0 | 0 | 0 | 0 / 0 |
| downlevel-simple | 51,764 / 51,764 | 0 | 0 | 279 / 279 |
| checks-only | 0 / 0 | 0 | 21,685 | 0 / 0 |

Bundle gathered **300,227 CFG nodes natively**, with four original-traversal fallbacks. Typed-app and SIMPLE gathered every attempted CFG node natively. Their function-window visits were 1,460, 17,685 and 4,311 respectively, across two FSIV pass invocations each; these are not distinct source functions. Whitespace and checks-only do not run this phase.

Captured validation recorded **2,558 typed-app hits** and **21,685 checks-only hits**, one seed each and zero captured fallbacks. Other tasks did not activate that inference route.

The diagnostics enabled recipe/shared/captured, native gather and lazy-use factors; parser/CommonJS factors were off. Five isolated workspaces use N2 and affinity CPUs 0,1: [exact flags and immutable manifest](latest-main-next-captured-gather-activation-v3/manifest.json), [case/run commands](latest-main-next-captured-gather-activation-v3/run-commands.json).

## Correctness evidence

- Plain production race suite: **205 PASS entries and one optional Angular SKIP**, with no failure or race report. The Angular fixture separately passed, bringing the total to **206 distinct cases exercised**. See [race log](tests-latest-main-next-focused-race.log), [race command](tests-latest-main-next-focused-race-command.json), and [Angular rerun](tests-latest-main-next-angular-race.log).
- All **15 CLI artifact gates MATCH** across five tasks and N=1/2/4: [gate log](gates-latest-main-next-standard.log).
- All five actual diagnostics match pinned **stdout, stderr, exit and every emitted JS/map digest**, including nonempty warnings: [verified per-task results and counters](latest-main-next-activation-verified.json).

The standard CLI and plain race binaries use the same baf source: [build commands](latest-main-next-build-commands.json), [source/commit provenance](latest-main-next-integration-source.json), and [plain focused-overlay manifest](focused-latest-main-next-overlay/manifest.json). The source record described the initial unbuilt checkpoint; completed validation is linked above. Observation helpers are separate ordinary-package overlays, excluded from production/race binaries and timed runs: [diagnostic build command](latest-main-next-captured-gather-activation-v3/build-command.json).

Counts establish activation, not CPU coverage or speedup. Function visits are per pass; modeled/accounted byte counters are not live memory, peak heap or RSS. This report makes no timing, RSS or Java/Go ratio claim; those belong to the separate macro measurements. Build-only memory controls are not runtime GC tuning.
