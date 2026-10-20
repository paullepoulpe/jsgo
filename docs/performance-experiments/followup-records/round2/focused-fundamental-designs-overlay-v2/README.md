# Fundamental compiler designs: plain focused test overlay

Root-parameterized preparation:

```
python3 prepare.py --root /workspace/jsgo-fundamental-compiler-designs
```

This extends the frozen focused-fresh-template overlay with `native_array_unknown_recipe_test.go`, `parallel_commonjs_test.go`, and `fsiv_use_cardinality_test.go`. It keeps the same package import path and every original production source file. Every replacement path ends in `_test.go`; there are no production replacements, diagnostic counters, synchronization hooks, or mocks. It is suitable for normal and race builds of the actual production implementation.

The utility-only declarations from `reaching_use_def_tester_test.go` are extracted verbatim into `retained_helpers_test.go`, alongside the original overlay's helpers. Its complete declaration text, original source hash, and retained test hashes are recorded in `manifest.json`. `truth_test.go` is retained as original assertion-helper support (it declares no test suite); no compiler/JUnit ported suite is pulled in. Preparation also verified the extracted declarations remain verbatim after gofmt.

New suite regex:

```
^(TestNativeArrayUnknownRecipe|TestParallelCommonJS|TestFSIVUseCardinality)
```

No Go build or test has been run by this agent. Optional activation-only production instrumentation will be prepared in a separate directory and must not be merged into this race overlay.

## V2 helper dependency correction

The first plain/diagnostic compilation attempts found missing `jsaMlInputs` and `jsaMlPACKAGE_JSON_MAIN_ENTRIES` in the retained CommonJS suite. Those attempts and their original overlay files remain intact. V2 extracts those two functions and their exact dependencies `jsaMlInput` and `jsaMap` verbatim from the original `deps_module_loader_test.go`, with its license and source/function SHA256 hashes. No original fixture assertion or production source changes. The ModuleLoader test suite remains excluded.
