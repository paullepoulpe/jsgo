Historical patch-preparation notes. Final scratch git-apply and exact-source validation passed for all six patches; see [validation receipt](../patch-apply-validation.json) and the [completed study](../README.md). Full upstream build/tests remain unrun. The original preparation note follows unchanged and is preserved verbatim in the evidence archive.

Upstream-layout source patches (v20261015)

Each .patch uses git-style a/src/com/google/... and b/src/com/google/... headers.
Apply from the upstream Closure Compiler v20261015 repository root, where Java
production files live under src/. The originals in ../../variants/ remain unchanged.
These new patches were regenerated from exact unshaded candidate source versus
/workspace/jsgo/.cache/closure-src originals, not from mechanically shaded build sources.

Each of five patches changes only that factor. all.patch combines all five over
the same stock base: nine disjoint production paths, with no conflict resolution.
Apply either all.patch OR selected single-factor patches; do not apply both.
No focused test sources are hidden inside production patches; their separate
paths/receipts remain in the report evidence index.

Validation performed: unified hunk application reconstructed exact candidate
source in memory. Actual git apply --check, upstream build/full test suites,
public-extension parity and benchmark acceptance remain future checks. Original
factor-patch hashes were compared before/after generation. No compiler, JVM,
Docker, archive or timing-output read occurred. manifest.json records small-patch
hashes and source paths. This is not a correctness/performance endorsement.
