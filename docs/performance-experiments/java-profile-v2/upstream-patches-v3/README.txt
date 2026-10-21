Closure v20261015: third normalized source-patch bundle

Promotion status: NONE_PROMOTED. Both candidates were unmeasured at task
assignment. The root agent's active empty-parameters timing screen was not read
or reviewed while preparing this bundle; no performance verdict is included.

Included independent patches

  empty-parameters.patch
    Return the already-computed empty immutable parameter-name set without
    walking the function body. Parameter extraction and its failures still run.
    For well-formed compiler ASTs, the skipped read-only walk cannot add a name.
    Argument side effects and separate receiver/closure analyses are unchanged.
    The existing malformed-external-AST exception boundary is intentional and
    covered explicitly by the candidate-only EmptyParametersProbe argument.

  scope-declare-once.patch
    Remove a repeated own-scope lookup during declaration by sharing the existing
    parent-shadow predicate. Existing canDeclare semantics, internal replacement,
    shadowing checks, ordering and allocation behavior remain unchanged.
    Actual Scope/TypedScope own-slot predicates are pure; future stateful
    same-package implementations would need a renewed review.

Both patches use a/src/com/... and b/src/com/... paths for the Closure Compiler
repository root. Each is based on the untouched pinned original for its file.
The two patches are validated individually, not in combination with each other
or with the earlier bundles. Existing candidate sources and patches are preserved.

Existing validation

Successful CONTROL/candidate build receipts are checked against each patch's
original and effective unshaded source hashes. Independent receipt review at
reviews/empty-scope-probes-v1/review.json passed all nine runs:
  ModifiedParametersProbe: STOCK / CONTROL / empty-parameters, no arguments.
  EmptyParametersProbe: STOCK / CONTROL without arguments; empty-parameters
    with --expect-shortcut only. Successful stdout is identical across all three.
  ScopeDeclareOnceProbe: STOCK / CONTROL / scope-declare-once, no arguments.

That review independently checks source snapshots/shading, flags/argv, narrow
overlay/probe class hashes, output bytes, exit statuses and class-load origins.
These focused probes do not establish full compiler-byte parity or performance.

Reproduction from a matching Closure Compiler checkout

  git apply --check --whitespace=error /absolute/path/to/selected.patch
  git apply --whitespace=error /absolute/path/to/selected.patch

prepare.py generates each normalized patch from original and candidate bytes.
It verifies the pinned source, closed source-hash record and existing successful
build receipts, then checks/applies/reverse-checks/reverses the patch in a fresh
temporary tree containing only the one original Java file. Forward bytes must
equal the candidate exactly, and reversed bytes must equal the original exactly.
The temporary tree is removed, and all small input hashes are rechecked.

proof.json records exact commands, statuses, output, source/patch/build hashes,
and the independent probe-review reference. No archive or class file was read
by the patch generator. No Java, compiler, benchmark or network action was run
while preparing this bundle. Earlier source-only preparation statuses remain
unchanged in frozen candidate documents; this is a separate later evidence record.
