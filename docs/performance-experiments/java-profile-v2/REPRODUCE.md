This snapshot preserves Java-versus-Java experiment evidence; it makes no new speed claim.

Read README.md and publication-jobs.json for the final reviewed jobs and retained
incomplete/diagnostic attempts. All ordinary experiment records are archived,
including rejected/unexecuted proposals, failed probes/builds, interrupted runs,
profiles' aggregate summaries and JIT annotations. evidence-members.json lists
every included file and every exclusion. Raw JFR files without an explicit passing
audit clearance are omitted from the public archive and retained locally, because
the profile configuration captured environment/system-property strings. Python
caches, symlinks and nonregular files are excluded; links are never followed.

reproduction-inputs/ contains exact per-job spec/oracle/oracle-metadata/helper
snapshots and the shared CliFactors source, externally referenced v1 driver classes
and build receipt. Its manifest identifies original paths and verified hashes.
Optional reference-only metadata supplies workload inputs without a v2 result.
Candidate/control source, overlays, class inventories, probes, build records,
controllers, reviewers and selection notes remain in the main evidence tree.

The published Closure jar, source jar, JDK/dependencies, immutable container image
and workload payloads are external dependencies identified by the recorded pins;
they are not copied into this archive. Remap original absolute paths in a new
workspace, preserve original versions, and record hashes for adapted helpers.
Old runtime coordination receipts are historical evidence, not portable grants.

Use the original method: fresh compiler and normal extern parsing each request,
same flags/default collector, actual 4 GiB/no swap/two CPUs per JVM, matched settled
warmups, sequential alternating requests, six output digests, independent cohorts,
zero new OOM and stopped cleanup. Performance tests use no JFR. Screens remain
separate from confirmation. JFR samples are diagnostic, allocated bytes are not
retained heap, RSS windows differ, and unavailable PSI is not zero. Full upstream
unit suites were not implied by the focused probes or workload parity checks.
