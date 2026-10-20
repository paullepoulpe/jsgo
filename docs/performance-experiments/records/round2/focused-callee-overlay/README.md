# Focused callee cleanup validation overlay

This overlay excludes unrelated **test declarations only**. All production Go
files and the actual package/import path remain unchanged. It retains callee,
argument-facts, function-FSIV, scope-window and copied-convention tests. Three
original function-argument test helper functions are copied verbatim into a
replacement for their otherwise large JUnit suite. The original license notice
is retained. No production mocks or compile flags changing semantics are used.

Prepare (source only):

```sh
python3 /workspace/jsgo-results/round2/focused-callee-overlay/prepare.py --root /workspace/jsgo-round-work-elision
```

When root grants the build reservation, from that worktree run:

```sh
env GOGC=25 GOMAXPROCS=1 /workspace/toolchains/go1.27.2/bin/go test -p=1 -pgo=off -overlay=/workspace/jsgo-results/round2/focused-callee-overlay/overlay.json -c -o /workspace/jsgo-results/round2/bin/jscomp-callee-focused.test ./internal/jscomp
```

Run from the worktree's `internal/jscomp` directory. For the immediate cleanup
regression use `-test.run '^TestReadOnlyCalleeFactsExceptionalCleanupAfterPopulation$'`.
Broader retained validation uses
`-test.run '^(TestReadOnlyCalleeFacts|TestInlineArgumentFacts|TestFunctionFSIV|TestAuditedConvention|TestFSIVScopeSharing)'`.

This is a focused supplementary validation binary; it does not replace the full
compiler/traversal suite or benchmark binary. No build was started while
preparing this overlay.
