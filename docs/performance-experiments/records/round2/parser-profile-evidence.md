# Activated parser/compiler profile evidence — 2026-10-19

> Attribution correction: retained stack-family totals are lower bounds, not hard performance ceilings. CPU stacks are limited to64 PCs; owning pass frames can be lost during generic recursion. See `profile-stack-depth-audit.md`.

These are single diagnostic captures, not repeated timing results. All five captures passed the stock-Java byte gates. The enabled and inactive captures use the same frozen source `11534351`, Go1.27.2, Node80 and default runtime GC; `profiles` is older `da20b44f` and includes source changes as a confound. Enabled provenance is in `parser-profiles-enabled/environment.json`.

## Activation is confirmed

Enabled raw locations contain `nativeParserAttempt`, `(*nativeParser).expressionBody`, `originalFragmentAt`, `(*nativeMustFlow).compute`, `(*nativeMaybeFlow).compute`, `(*nativeCFGBuilder)` and native reference-key/storage lines. Representative frames are saved in `parser-profile-activation-frames.txt`. The inactive archive instead records ParserRunner_Parse at line129 (the explicit disabled-native branch), so it must not be described as an activated-parser run. Thread budget1 means the parallel-FSIV flag can be enabled while its cohort guard intentionally falls back to serial.

## Exact allocation and GC deltas

TotalAlloc/Mallocs/NumGC are after-minus-before MemStats, measured before the overlay’s final diagnostic GC. MiB are binary units. RSS is not measured here. GC-count changes of one cycle are normal threshold effects.

| Task | da20 alloc MiB | same-source OFF MiB | enabled MiB | enabled vs OFF | mallocs OFF→ON, million | GC cycles OFF→ON | GC cycles da20 |
|---|---:|---:|---:|---:|---:|---:|---:|
| typed-app | 1313.7 | 1290.4 | 1245.8 | -3.45% | 25.32→23.88 | 32→31 | 32 |
| bundle | 5326.5 | 4983.5 | 4831.5 | -3.05% | 97.67→92.05 | 54→55 | 57 |
| downlevel-whitespace | 426.9 | 427.2 | 419.0 | -1.93% | 5.75→5.55 | 21→21 | 22 |
| downlevel-simple | 1152.4 | 1127.3 | 1096.8 | -2.70% | 19.41→18.28 | 36→35 | 37 |
| checks-only | 1965.3 | 1966.7 | 1782.5 | -9.37% | 32.15→28.33 | 29→27 | 29 |


## Sampled frontend CPU distribution

CPU seconds aggregate samples across runtime threads, not elapsed seconds. Frontend is the union of parsing/parser package frames, so shared scanner/IRFactory/JSDoc work is included. Native includes its fragment calls; the next column separates those calls. Full-original and fragment figures are sampled lower bounds when deep stacks lose outer context; zero is not proof that no fallback occurred. Categories such as scanner and IRFactory overlap and must not be summed.

| Task | all CPU s | frontend OFF→ON s | frontend ON/all | native excluding fragment s | original fragment s | full original s | GC CPU s |
|---|---:|---:|---:|---:|---:|---:|---:|
| typed-app | 6.87 | 0.43→0.31 | 4.5% | 0.28 | 0.02 | 0.01 | 1.45 |
| bundle | 26.59 | 1.14→0.88 | 3.3% | 0.47 | 0.33 | 0.00 | 6.12 |
| downlevel-whitespace | 1.91 | 0.30→0.26 | 13.6% | 0.12 | 0.01 | 0.12 | 0.48 |
| downlevel-simple | 5.25 | 0.30→0.32 | 6.1% | 0.17 | 0.02 | 0.11 | 1.26 |
| checks-only | 10.20 | 1.44→1.07 | 10.5% | 0.82 | 0.22 | 0.02 | 2.15 |


The activated parser removes much original parser work, but the whole-program CPU ceiling is modest for typed-app and bundle: frontend is only4.5% and3.3% of their current sampled CPU. Further bundle fragment elimination targets0.33s out of26.59s (~1.2%). Checks-only remains a larger frontend opportunity (10.5% CPU; sampled frontend allocation534.6MiB), including scanner295.1MiB and native nonfragment466.9MiB, with overlapping categories. Whitespace still samples substantial whole-original work (0.12s of0.26s frontend); sample stacks alone cannot identify which sources or rejection reasons.

## Current larger costs

`GetNext` is already inlined yet accounts for10.6–11.4% flat CPU on typed/bundle/checks. This points to visit count and pointer access locality, not a missing accessor inline. GC/assist stacks account for~21–25% of sampled CPU. Bundle scope scanning is2.36s (8.9% CPU), RemoveUnusedCode2.35s (8.8%), ReferenceCollector including its behaviors3.83s (14.4%). These cumulative areas overlap. Checks-only TypeInference is2.42s (23.7%) and615.2MiB sampled allocation.

Bundle sampled allocation still includes NewVar~0.23GiB, scope objects~0.14GiB, string-keyed ordered-map insertion~0.30GiB cumulative, native reference collection storage~0.24GiB cumulative and separate Reference objects~0.14GiB. This supports evaluating the already-frozen scope-reuse candidate and emitter/backend work, rather than assuming more parser syntax yields a large compiler-wide win. Scope reuse must prove useful hit counts and retained-memory cost; its validation still traverses declarations.

Sampled allocation totals are estimates and can lag recent allocations; use exact MemStats for total-byte comparisons. No Java memory/time ratio or warmed-Java victory can be inferred from these diagnostic captures. Raw CPU/allocation tables and component sums are in `parser-profile-analysis-enabled.json`; reproduce with `analyze_parser_profiles.py --candidate parser-profiles-enabled --output parser-profile-analysis-enabled.json`.

Next source investigation: whitespace emission and conversion costs, as requested by root.
