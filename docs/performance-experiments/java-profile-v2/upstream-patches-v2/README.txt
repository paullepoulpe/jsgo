Closure v20261015: second normalized source-patch bundle

Promotion status: NONE_PROMOTED.
This bundle establishes individual patch applicability and exact source identity.
It does not establish a performance benefit, universal API equivalence, upstream
acceptance, or a validated combination. Existing upstream-patches and every
candidate/control/frozen artifact remain unchanged.

Status snapshot at preparation
  traversal-leaf-name: MEASURED_SCREEN_NEGATIVE_NO_PROMOTION. The completed
    CONTROL comparison did not establish a useful gain; this is not a promoted
    optimization. See results/screen-bigquery-control-traversal-leaf-name-v1.
  equality-prototype-scalar: TIMING_PENDING. A source-guard opportunity and a
    successful build are not performance evidence. No timing verdict is included.
  reference-function-first-child, graph-edge-copy, source-position-last:
    NOT_MEASURED. No elapsed-time or memory benefit is claimed.

All five have existing successful build-validation receipts at preparation.
This source-only bundle verifies their source hashes against those receipts;
it does not rerun the build, probes, compiler byte gates, or benchmarks. Candidate
README/manifests retain their original preparation-time status; later runtime
evidence lives in the separate build/results directories. Status here is a
dated snapshot, not an automatically updated verdict on future runs.

Included independent patches
  traversal-leaf-name.patch
    Inline visits of leaf NAME children for callback bases with final always-true
    shouldTraverse methods. Raw Java stack shape differs for eligible visits.
  reference-function-first-child.patch
    Avoid the arrow-function flag lookup for a function's first child in
    Reference.isDeclarationHelper; preserve the existing boolean result.
  graph-edge-copy.patch
    Copy each node's outgoing edges by index into a capacity-reserved snapshot,
    avoiding addAll's temporary input array. Concurrent mutation and resource-
    failure boundaries are outside the ordinary snapshot equivalence proof.
  source-position-last.patch
    Reuse the last scanner-created position when its full position tuple matches,
    after all original scanning/validation. Public SourcePosition identity and
    aliasing can change; universal third-party object-identity parity is not claimed.
  equality-prototype-scalar.patch
    Skip allocating EqualityChecker for an exact PrototypeObjectType against
    final null/void/unknown types only when the left unknown cache is already
    known and its template map is empty; preserve the original fallback otherwise.

Each patch has a/src/com/... and b/src/com/... paths for the Closure Compiler
repository root and is based on its own untouched pinned v20261015 source.
These patches were NOT validated in combination, either with one another or
with any patch from the earlier upstream-patches bundle. Apply one selected
patch to a matching checkout at a time; do not infer combination validity from
the fact that some patches touch different Java files.

Usage from a matching Closure Compiler checkout:
  git apply --check --whitespace=error /absolute/path/to/selected.patch
  git apply --whitespace=error /absolute/path/to/selected.patch

Source proof
  prepare.py generates diffs from the actual frozen control-src and src bytes,
  never by rewriting an earlier patch's text. For each of the five candidates:
  - Verify control/candidate hashes against its frozen source manifest.
  - Verify the control bytes against the pinned extracted Java source.
  - Verify both original/effective unshaded source hashes in the existing
    successful CONTROL/candidate build receipts, including their manifest hashes.
  - In a fresh tiny temporary tree, run forward check/apply and compare exact
    candidate bytes; run reverse check/apply and compare exact original bytes.
  - Verify only the intended repository-relative Java file exists in that tree,
    remove the tree, and recheck the source/receipt inputs after generation.

proof.json records command argv/output/exit codes, normalized patch hashes,
frozen/pinned/applied/restored source hashes, build-receipt links and hashes, the
status snapshot and source-only validation boundary. No jar/class/archive was
opened or hashed. Only these small source, manifest, receipt and patch files were
read/hashed. No compiler, JVM, network request or benchmark was run for this bundle.
