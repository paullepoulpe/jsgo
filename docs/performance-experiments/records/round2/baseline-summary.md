# Current-main one-shot CLI baseline

Source da20b44f (compiler tree matches origin/main 123b46e8), Go1.27.2 with shipped PGO. Stock pinned Java jar; fresh compiler and normal extern parsing each iteration. CPU affinity0,1, cgroup quota2CPUs. Compiler threads1, runtime defaultGC/CPU detection.

Go: median of3 fresh processes. Java: median10 measured iterations after12 warmups, one JVM per task. RSS: median fresh Go process peaks versus whole Java JVM lifetime peak including warmups. All five artifact gates passed. Java ratios reference one cohort and are screening context, not independent-cohort confidence estimates.

| Task | Go s | Java warm s | Java/Go time | Go RSS MiB | Java RSS MiB | Java/Go RSS |
|---|---:|---:|---:|---:|---:|---:|
| typed-app | 5.795 | 3.739 | 0.65x | 257.3 | 1090.2 | 4.24x |
| bundle | 20.307 | 16.600 | 0.82x | 612.5 | 1696.6 | 2.77x |
| downlevel-whitespace | 1.467 | 0.938 | 0.64x | 164.9 | 930.2 | 5.64x |
| downlevel-simple | 4.404 | 2.805 | 0.64x | 199.6 | 1118.4 | 5.60x |
| checks-only | 8.569 | 6.038 | 0.70x | 657.1 | 1360.0 | 2.07x |

Raw observations: baseline/; profiles: profiles/; binary/build provenance: build-main-provenance.json.
