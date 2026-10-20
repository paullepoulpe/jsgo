Reproducing this snapshot requires path remapping and external pinned inputs.

The evidence archive contains the experiment's regular source/build/probe/diagnostic,
launch, raw-request, failure, review and frozen-version records. It does not contain
the external JDK, stock Closure jar, container image or benchmark dataset payloads.
Symlinks are metadata only and are never followed or stored as archive links.

reproduction-inputs/workloads/{bigquery,bundle}/ contains exact spec, six-output
oracle and oracle-metadata bytes verified against the actual completed runner's
metadata pins. helpers/ contains the two exact pinned helper sources; the experiment
already contains the shared CliFactors.java/classes, overlay sources/classes,
build receipts and frozen runner/reviewer/analysis sources. manifest.json maps copies
back to original absolute paths, hashes and runtime metadata.

Obtain the pinned Closure v20261015 Maven jar and matching sources, the recorded
Java21 installation/dependencies, immutable container-image digest and the exact
external corpus described by each spec, oracle metadata and runtime input_manifest.
Check their recorded hashes before running. Do not replace input files with different
versions or compare to old timing rows. Dependency relocation and all generated
nested classes must follow build receipts; STOCK remains the unmodified jar.

Original paths are intentionally retained in receipts. Remap spec argv/input/cwd
paths, ROOT/HELP/JDK/jar/image/overlay/driver paths and mounts into your isolated
workspace. Use a NEW version of each adapted helper/manifest and record its hashes;
existing runtime grants and original machine paths are not portable authorization.
Retain originals alongside adaptations. Whole upstream suites and custom-extension
boundaries are separate from these strict byte gates and narrow focused probes.

Follow report/methodology or methodology.md: actual4GiB/CPU2/default collector,
fresh compiler+extern parsing per request, matched settling, sequential alternating
JVM requests, strict six digests, fresh independent cohorts and stopped/zeroOOM.
Current RSS, process-lifetime HWM, container peak and allocated bytes have different
windows. Missing PSI is unavailable, not zero. Screens are excluded from holdout CI.
This local publisher neither runs those steps nor establishes a new speed claim.
