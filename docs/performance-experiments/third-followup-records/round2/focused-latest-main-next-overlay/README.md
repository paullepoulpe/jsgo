Plain focused overlay for latest-main compiler designs.

Every production source remains unchanged. All new jscomp tests since ff06ac15 are retained verbatim, together with original ReachingUseDefTester and Truth utility files. Four original ModuleLoader helper declarations are copied verbatim into an overlay-only test file, replacing only their original test file. The exact original scopeTestParse function and contiguous nodeTraversalTestCallback/nodeTraversalTestScopedCallback type+four-method prefix are copied verbatim into another overlay-only helper replacing scope_test_helpers_test.go. All other original jscomp test files become a package-only stub. The original package/import path remains unchanged. Manifest records source and declaration SHA-256 provenance. No diagnostic hooks, counters, synchronization, admission changes, or mocks are injected.

Generate with `python3 prepare.py --root /path/to/worktree`; do not regenerate files while a build uses this overlay.
