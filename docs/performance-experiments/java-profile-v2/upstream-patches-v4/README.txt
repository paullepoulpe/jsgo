ReferenceCollection initial-capacity experiment: exact upstream patch

Status: NONE_PROMOTED. This one-line patch is retained as a completed negative
measurement candidate, not recommended for application. Existing upstream-patches,
upstream-patches-v2 and upstream-patches-v3 remain unchanged. Individual patches
have not been validated in combination.

reference-list-two.patch is generated from the exact frozen CONTROL and candidate
bytes, normalized to a/src/... and b/src/... in a Closure Compiler checkout.
prepare.py checks the pinned v20261015 source, frozen source hashes, build-source
receipts, existing probe outputs/class origins and tiny overlay classes. It then
applies and reverses this patch in an isolated temporary tree and compares exact
bytes. proof.json and ../reviews/reference-list-two-probes-attempt-1/review.json
record these checks. No Java, compiler, benchmark or archive access is performed.

All three focused STOCK/CONTROL/candidate probe runs passed with identical stdout.
They cover ordering/identity, nulls/duplicates, iterator removal/fail-fast, clearing,
empty queries and serialization; they are not the full upstream Java suite.

The completed BigQuery CONTROL screen summary reports candidate/baseline ratios:
wall 1.037192071 (+3.72%), CPU 1.012244898 (+1.22%), allocated bytes 1.001337697
(+0.134%), lifetime high-water memory 1.129029255 (+12.90%). This is one independent
JVM cohort, with five measured requests per variant after 62 warm requests; it is
a negative screen and provides no confidence interval or general speed claim.
These metrics are transcribed from the recorded summary, not independently rerun
or promoted by this source/probe receipt review.

Changing new ArrayList<>() to new ArrayList<>(2) preserves the logical compiler
operations examined but eagerly allocates empty lists and increases growth/copying
for lists longer than two. Capacity reflection, resource exhaustion timing and
arbitrary same-package capacity-changing calls are outside the equivalence claim.
The frozen candidate README retains its historical source-only preparation status;
later build/probe/screen evidence is linked here without mutating frozen inputs.
