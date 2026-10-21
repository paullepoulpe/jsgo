Closure v20261015: normalized source patches

Status: none promoted. These are five independent, already-built candidates.
This package establishes patch applicability and exact source identity only;
it does not claim a confirmed performance benefit or upstream acceptance.
The candidate directories and their original patches remain unchanged.

Files
  accumulator.patch
    Reuse one result set through findModifiedParameters recursion.
  identity-differs.patch
    Return false immediately when JSType.differsFrom compares an object to itself.
  traversal-noop.patch
    Classify callbacks whose inherited final methods do no work, and bypass
    those method calls in the ordinary traversal path.
  traversal-name-switch.patch
    Bypass the initial special-handler switch for NAME nodes, which otherwise
    take its empty default arm. This is independent of traversal-noop.
  leaf-children-return.patch
    Return before allocating traversal collections when a property's children
    set is empty in DeadPropertyAssignmentElimination.Property.markChildrenRead.

Every patch uses a/src/com/... and b/src/com/... paths for the Closure Compiler
repository root. Each is against its own untouched pinned source, not a stack
of preceding candidates. In particular, do not concatenate the two traversal
patches and treat the result as a validated combined candidate.

Usage from a matching Closure Compiler checkout (one selected patch at a time):
  git apply --check /absolute/path/to/selected.patch
  git apply /absolute/path/to/selected.patch

proof.json records source and patch SHA-256 hashes, links to existing successful
build receipts, exact git argv/output/exit statuses, and applied-source hashes.
prepare.py generates these diffs from frozen control-src and src snapshots.
For each candidate it creates a fresh temporary tree containing only the
original changed Java source, runs git apply --check and git apply with
--whitespace=error, checks the resulting bytes against the frozen candidate,
then checks reverse applicability. The temporary trees are removed afterward.
The originals also match the local pinned extracted source and the original
source hashes in the existing build receipts; candidates match the effective
unshaded source hashes in those same receipts.

No compiler, JVM, benchmark, network request, or archive read/hash is involved
in preparing this package. Existing build/probe/timing evidence lives outside
this directory; root controls promotion after the required validation.
