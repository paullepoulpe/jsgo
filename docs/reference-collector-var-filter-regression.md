# Restore per-reference public variable filters

Base `bb521c0f` moved the exact-Var cache check ahead of the public varFilter callback. After one accepted reference primes lastVar, another AddReference for that same Var bypasses the predicate entirely. For an alternating predicate the callback runs once and subsequent rejected references are accepted. This also suppresses callbacks that count, panic or reenter the collector.

The fix restores the pre-bb ordering: invoke varFilter exactly once first, return immediately when rejected, then call getOrCreateReferences. That method retains the existing exact-Var/table lookup optimization and consults current cache state after callback reentry. No inferred purity/private fast flag is added. Reference allocation, scope/input/block access, pooling/release/reset guards and all other collection code remain unchanged. This is a public callback correctness fix, not a performance claim.

New file `internal/jscomp/reference_collector_var_filter_regression_test.go` retains all existing tests and adds:

* TestReferenceCollectorVarFilterEverySameVarCall: alternating acceptance/callback count, exact Var identity and reference order/metadata with both small-buffer and populated-table storage.
* TestReferenceCollectorVarFilterRejectsBeforeTraversalAndAllocation: cached rejection occurs before invalid traversal/block access, with no reference/collection/arena storage change and pooling enabled.
* TestReferenceCollectorVarFilterReentrant: public callback reenters AddReference for the same or another Var, preserving nested callback/reference order and selecting the correct post-callback collection.
* TestReferenceCollectorVarFilterPanicOnCachedVar: the second callback panic propagates before traversal access or reference append.

Source formatting/diff checks only have run. No repository Go build/test or compiler runtime is authorized for this branch yet. The complete integrated normal suite should include every original and new test file. Focused pattern: `^TestReferenceCollectorVarFilter` (plus the existing TestReferenceCollectorTest suite).

The source counterexample is direct: base addReferenceWithScope takes the lastVar branch before varFilter; each regression primes that branch. Executed stock-red evidence is still pending. Running these new internal-package tests against stock requires a separately coordinated full jscomp test-package build with the new file present; no reduced-source substitute is claimed as stock proof. A tiny public-API probe importing an already cached normal jscomp package may be an alternative if root approves it, but has not been implemented/run. Do not label source inspection as executed failure.
