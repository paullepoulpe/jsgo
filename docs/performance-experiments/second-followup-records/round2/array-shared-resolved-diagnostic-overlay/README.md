# Untimed native Array resolved-vector activation

NEW output directory; previous Array diagnostic sources/outputs are untouched.
Generator refuses to overwrite a frozen manifest/overlay. It never builds or
runs a compiler. Await root's integrated source before invoking it:

    python3 prepare.py --root /workspace/jsgo-post-colors-array-sharing \
      --tasks typed-app,checks-only

For a root-owned focused-test overlay, add `--base-overlay PATH`. Retain these
source test files in that base overlay:

- native_array_unknown_recipe_test.go
- native_array_shared_resolved_test.go

The generator adds the test-only file `array_shared_resolved_diagnostic_test.go`
and helper `array_shared_resolved_diagnostic_probe.go` by overlay. Production
CLI must remain a separate, non-diagnostic binary. The focused diagnostic binary
uses the exact original sources apart from the numeric-probe entry/exit hook
and omission of AST-retaining probe observations. It changes no options,
compiler admission, custom guards, experiment flags or public constructors.

Exact diagnostic regex: `^TestArraySharedResolvedDiagnosticPinnedTasks$`.
Set `JSGO_ARRAY_SHARED_DIAGNOSTIC_JOBS` to this NEW directory's generated
`jobs.json`. Clear inherited experiment/runtime variables, then use the matched
native9+CLI_SOURCE_LOCK baseline plus ARRAY_UNKNOWN_RECIPE=1 and
ARRAY_SHARED_RESOLVED=1. Keep default runtime GC. The test itself injects no flags.
Do not run while timing is reserved.

The test invokes actual untouched stock CommandLineRunner.DoRun on the exact
pinned typed/checks jobs in independent workspaces, checks all artifacts against
pinned Java SHA-256 digests, then reports original seeds, recipe hits and separate
shared-resolved hits. With sharing enabled, shared hits must equal recipe hits;
with sharing disabled they must be zero. New counters appear as task JSON lines
prefixed `ARRAY_SHARED_RESOLVED_DIAGNOSTIC` and per-job
`array-shared-resolved-counts.json`. Method-event counts are not mallocs or RSS;
no diagnostic elapsed time or memory value may be reported as a benchmark.
