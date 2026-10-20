# FSIV quotient activation-only overlay

This is an **untimed, non-race diagnostic overlay**. It adds numeric aggregation protected by a diagnostic mutex; it must never be used for race validation or reported performance. The plain focused race overlay remains unchanged.

Prepare against any checkout containing the frozen quotient prototype:

```
python3 prepare.py --root /workspace/jsgo-fundamental-compiler-designs \
  --base-overlay /workspace/jsgo-results/round2/focused-fundamental-designs-overlay/overlay.json
```

`--base-overlay` can instead be the Array/CommonJS diagnostic overlay; it preserves all prior replacements and composes the two narrow production hooks from their effective overlaid source. Original and resulting source hashes are recorded. Admission, experiment flags, compiler options, observers, warning guards, source ownership, pass factories and priority topology remain original.

The hooks attach a private probe over each function's original `EnterScope` analysis/rewrite window after scope facts are obtained, including worker-owned passes. The final deferred copy includes actual candidate queries and runs after original cleanup on normal/panic exit. Aggregation retains only integers, not AST/Scope/Var/CFG objects. No probe is attached over an existing unit-test probe. `DoRun` resets counts and writes one JSON report directly to the absolute file selected by `JSGO_FSIV_USE_CARDINALITY_DIAGNOSTIC_JSON`; stdout/stderr are unchanged. Report-file I/O failures are detected by the diagnostic test's required report read.

To exercise the exact pinned benchmark inputs through the actual stock CLI, reuse the Array diagnostic `jobs.json` (root/task args, workspaces, pinned Java artifact hashes):

```
JSGO_FSIV_USE_CARDINALITY_DIAGNOSTIC_JOBS=/workspace/jsgo-results/round2/array-recipe-diagnostic-overlay/jobs.json \
  <diagnostic-test-binary> -test.run '^TestFSIVUseCardinalityDiagnosticPinnedTasks$' -test.v
```

Set the intended experimental flags in the environment; the test does not enable any itself. It calls stock `CommandLineRunner.DoRun`, validates JS/stdout/stderr/exit/source-map hashes against pinned Java, writes per-job `fsiv-use-cardinality.json`, and logs its counts. The diagnostic can be run at N1/N2/N4 by preparing corresponding job argument manifests, keeping N identical to the intended comparison. No benchmark timings should be collected from this binary.

Counts distinguish all admitted-to-FSIV function windows from backward-analysis attempts (functions with no candidates do not need backward analysis), then actual quotient admissions and weighted rejected CFG nodes, stock/memory/step/collision/AST rejection reasons, generated facts/actions/solver steps, and actual singleton queries. `DenseBytes` is cumulative allocation for the two int32 state matrices; `MaxDenseBytes` is the largest individual matrix pair. They exclude preparatory programs/topology/rounding and are **not RSS or retained-memory measurements**. A low node-weighted admission fraction is a meaningful ceiling limitation even when function-count admission looks high.

Prepared source only. No Go builds or execution by this agent.
