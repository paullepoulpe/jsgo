# Ordinary CLI redesign screen

These experiments are isolated from the draft PR's retained implementation.
Worker processes and parsed-extern caches are excluded. This screen is one
invocation per variant and workload, so it cannot establish small gains or
stable memory differences.

| Workload | Draft PR CLI (time / peak RSS) | Combined parser + native CFG (time / peak RSS) |
|---|---:|---:|
| typed-app | 5.224 s / 268.2 MiB | 5.093 s / 267.5 MiB |
| bundle | 18.610 s / 643.0 MiB | 18.782 s / 497.7 MiB |
| downlevel-whitespace | 1.341 s / 158.9 MiB | 1.252 s / 171.1 MiB |

The architecture control, with both experiments disabled, took 5.303 s,
18.739 s and 1.238 s respectively. The direct parser and CFG therefore show
no demonstrated incremental elapsed-time win on these full builds. Individual
and combined results, including reference Java time and RSS ratios, are in
[cli-redesign-screen.json](cli-redesign-screen.json). Do not add the isolated
parser microbenchmark speedup to these end-to-end results.

Before allocator isolation, the combined parser/CFG passed native CFG topology and data-flow differential tests,
the existing type-inference/type-checking/analysis/traversal suites with the
experiments enabled, and 36 exact CLI Java comparisons at 1, 2 and 4 compiler
jobs. After fixing speculative parser allocation, six additional CLI cases
cover ZIP inputs, parse-error order, comments and source maps at 1 and 4 jobs.
The parser allocator passed focused normal and race checks, including an exact
oracle for 671 core-js files, concurrent fallback and a regression which proves
discarded ASTs and templates are collectible while another parse remains live.

Paged type-state storage is excluded from further retention: its first screen
had no speed benefit and higher bundle RSS. Its allocation audit is recorded in
[experiments/indexed-flow-slots.md](experiments/indexed-flow-slots.md).

Current evidence supports keeping these redesigns experimental. It does not
establish that Go beats warmed Java or that these designs have reached their
performance ceiling.
