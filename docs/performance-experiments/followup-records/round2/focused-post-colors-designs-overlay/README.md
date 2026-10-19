# Plain focused post-colors compiler-design overlay

Frozen source: `/workspace/jsgo-post-colors-design-integration`, `3750c2599170db23522d32d8515118c2d1cbe826`.

```
python3 prepare.py --root /workspace/jsgo-post-colors-design-integration
```

Built from the preserved V2 fundamental overlay helper closure. Retains the three fundamental suites plus `native_live_phase_admission_test.go`, `post_colors_native_admission_test.go`, and `remove_unused_stock_colors_test.go`. `native_output_buffer_test.go` was already retained. Total44 original test/helper files,420 test-only replacements.

This overlay changes no production file, guard, factory, callback or synchronization. It has no diagnostic instrumentation. Every replacement key ends in `_test.go`; excluded suites become a `package jscomp` stub. Suitable for ordinary and race validation of the original implementation. No builds/tests run by this agent.

The eight original helper-extraction source hashes, resulting overlay-helper hash, retained test-file hashes, verbatim ReachingUseDefTester declaration hash and four verbatim CommonJS helper hashes are recorded in manifest.json. Original declarations/licenses remain unchanged. Original Google Truth assertion helper source is retained without importing JUnit/compiler suites. All hashes, expected frozen source and declaration contents were verified after generation.

All earlier overlays and failed-attempt evidence remain intact. Separate activation-only overlays must be generated separately and must never be composed into this plain race overlay.
