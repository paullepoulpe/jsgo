Untimed capture-rejection diagnosis only. Frozen production source is never edited.

Generator: prepare.py --root MERGED_CHECKOUT --base-overlay BASE_FOCUSED_OVERLAY --tasks typed-app,checks-only

Use the generated overlay.json in a centrally scheduled focused jscomp test build. Run TestArrayRecipeDiagnosticPinnedTasks with JSGO_ARRAY_RECIPE_DIAGNOSTIC_JOBS pointing to jobs.json and exactly the intended runtime experiment flags; neither generator nor test forces admission/flags. Never use the diagnostic binary for timings/RSS or performance claims.

The existing numeric reports now distinguish exact native Array<?> calls, Capture attempts, and Capture rejection. ARRAY_RECIPE_CAPTURE_DETAIL prints rejection labels and at most the first capture failure per task. Snapshots store scalars, concrete dynamic type names (%T), key names and at most sixteen vector items; they retain no AST/registry/type pointers and invoke no virtual type String/equality/hash/resolve operations. Capture continues to return nil at exactly its original branch; no admission guard is changed.

The pinned task arguments and Java status/stdout/stderr/JS/map SHA256 comparisons are unchanged. The base overlay excludes only test files; no production stubs are introduced. Manifest includes source revision and the diagnostic production replacement hashes.
