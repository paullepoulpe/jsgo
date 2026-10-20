# Fresh extern parsing and scope declaration recipe screen

All 45 fresh CLI runs (3 variants × five unchanged tasks × 3 runs) matched pinned Java output, ordered diagnostics, exit status and source maps. Additional CLI gates passed all five tasks at N=1 and N=2. No worker or pre-parsed extern AST is used. Runtime GC and CPU settings are defaults; inherited GOGC/GOMEMLIMIT/GODEBUG/GOMAXPROCS are removed. Affinity 0,1 under the two-CPU quota; compiler jobs=1.

Prior production source 0fafa262; new CLI frozen641cb85d. Later7f3e971e corrects a RAW-stage test fixture and20c8452d updates docs, with no production difference. Go1.27.2, shipped PGO, four package -l=4 flags. Source hashes in build-recipes-externs-provenance.json; environments, alternating order and raw samples in recipes-externs-screen/.

The matched environment enables only the prior nine native/inline facts flags. New parser handles es3 special bindings using two original statement fragments (36 bytes); all source is parsed afresh. Scope recipes additionally enable only SCOPE_RECIPES. New parser and prior binary differences can include generated code layout; recipe on/off uses the same binary.

No broad gain is established. Parser medians are 0–2% faster but three samples and broad sample ranges cannot establish that small a gain. Recipes improve downlevel-simple about3.2% against same-binary control but regress the other tasks; typed-app RSS rises about7.2%. Recipes remain off. Neither candidate beats warm Java.

Ratios are Java / Go; above one favors Go. Java times/RSS here are the older stored warm cohort, not a fresh paired confidence estimate. Go RSS is median fresh CLI process peak; Java RSS is the whole warm JVM lifetime peak including warmup.

| Task | Variant | Go s | Go peak MiB | Java/Go time | Java/Go RSS |
|---|---|---:|---:|---:|---:|
| typed-app | prior | 5.256 | 258.3 | 0.711 | 4.22 |
| typed-app | externs | 5.226 | 255.3 | 0.715 | 4.27 |
| typed-app | recipes | 5.371 | 273.7 | 0.696 | 3.98 |
| bundle | prior | 18.392 | 438.1 | 0.903 | 3.87 |
| bundle | externs | 18.121 | 434.8 | 0.916 | 3.90 |
| bundle | recipes | 18.590 | 442.0 | 0.893 | 3.84 |
| downlevel-whitespace | prior | 1.292 | 164.4 | 0.726 | 5.66 |
| downlevel-whitespace | externs | 1.269 | 166.8 | 0.739 | 5.58 |
| downlevel-whitespace | recipes | 1.313 | 165.5 | 0.715 | 5.62 |
| downlevel-simple | prior | 3.764 | 211.0 | 0.745 | 5.30 |
| downlevel-simple | externs | 3.748 | 210.0 | 0.748 | 5.33 |
| downlevel-simple | recipes | 3.627 | 211.8 | 0.773 | 5.28 |
| checks-only | prior | 7.566 | 666.6 | 0.798 | 2.04 |
| checks-only | externs | 7.547 | 659.0 | 0.800 | 2.06 |
| checks-only | recipes | 7.599 | 660.7 | 0.795 | 2.06 |

Validation:807 focused checks PASS,2 SKIP, including exact fresh Scope/Var bindings, rename/move/clone/input invalidation, exception retirement and full SIMPLE output/AST/change-order oracle with actual recipe hits. Recipe race validation is pending the following combined source batch. Parsing differential family plus70 strict/sloppy special-binding fixtures and full fresh extern inventory pass:128/128 BROWSER and5/5 CUSTOM files accepted, zero full-file fallbacks, exact original AST/JSDoc/source/property/feature/diagnostic/exception parity.

Actual small SIMPLE validation records3scopehits and19savedscan visits; this is activation evidence, not a benchmark-workload hit-rate estimate. Logical recipe payload stats exclude map overhead, retained AST graphs and RSS.

New follow-up work targets compact contextual rename conflict storage, native live-variable solving, private first variable/map allocation, and ASCII printer/join storage. No runtime GC tuning is part of this work.
