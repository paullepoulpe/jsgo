Independent review: empty-parameters and scope-declare-once probes

PASS for all three completed probe receipts, each run against STOCK, rebuilt
CONTROL and candidate. All nine Java runs and all three probe compilations
record exit 0. Each probe's complete stdout matches STOCK byte-for-byte, and all
stderr is empty. The source/build/probe checks were performed without a new JVM.

  empty-parameters-modified-probes-attempt-1: ModifiedParametersProbe, no args.
  empty-parameters-boundary-probes-attempt-1: EmptyParametersProbe; only the
    candidate gets --expect-shortcut. This intentionally checks the documented
    malformed-body behavior boundary while keeping successful stdout identical.
  scope-declare-once-probes-attempt-1: ScopeDeclareOnceProbe, no args.

review.py independently checks original/current/frozen probe source hashes,
mechanical dependency shading, exact compilation/run argv and flags, timestamps,
output hashes and bytes, source snapshots, build manifest/receipt identities,
original/effective production source hashes, narrow overlay/probe class hashes,
and intended class-file sets. It parses the recorded class-load logs itself to
verify each changed production family and each probe main came from the expected
jar/overlay/probe directory. The completed receipts show no runtimes at their own
preflight/postflight boundaries; no claim is made about the currently active
root-owned benchmark process.

The scoped fixture covers actual untyped/typed declarations, duplicate handling,
implicit slots, parent shadowing/bleeding names, replacement, live views/order,
undeclare/clear, and typed reservation/error paths. The modified-parameter oracle
and empty-parameter fixture cover ordered names, nested closures/fields, rest and
unsupported parameters, AST stability, early failures and the explicit malformed
body boundary. Probe sources were unchanged from their previously closed review.

review.json records exact reviewed hashes and per-run evidence. Small explicit
source, receipt, log and overlay/probe files were read; archives were not opened
or rehashed. Jar hashes were compared only as recorded across the build/probe
receipts. No Java, build, benchmark or network operation was performed.

Both candidates were unmeasured at task assignment. This review did not read the
active empty-parameters performance screen. Focused probes establish neither
complete compiler output/diagnostic parity nor a performance benefit.
