# Concrete traversal experiment

`JSGO_EXPERIMENT_CONCRETE_TRAVERSAL=1` specializes the recursive `NodeTraversal`
driver for `ReferenceCollector_CollectorCallback`. The default remains the
original interface-callback driver. One type check at each public traversal
entry selects the concrete driver; its recursive walk calls concrete callback
methods without generic dictionaries or per-node callback assertions.

The generated driver shares the exact original `*NodeTraversal` state. It keeps
cached-next handling for AST mutation, lexical scope creation and enter/exit
notifications, class member ordering, destructuring/default evaluation order,
current node/function/change-scope state, and exception wrapping. Callbacks may
re-enter ordinary public traversal entries; they select the same appropriate
driver. Scope creators and all accessor methods are shared.

Regenerate or check it with:

```sh
python3 tools/bench/gen_concrete_traversal.py
python3 tools/bench/gen_concrete_traversal.py --check
```

The generator derives its algorithm directly from the translated
`internal/jscomp/node_traversal.go` and copies its Closure copyright/license
notice. It changes only driver method names/recursive calls and concrete
callback dispatch. A second, test-only specialization records all callback,
scope and current-node events and mutates sibling nodes. Differential tests
compare that driver's complete event sequence, resulting AST and exception
messages against the ordinary driver across all five entry methods, including
class fields/computed keys, defaults/destructuring, skipped subtrees, removal,
replacement, insertion and both destructuring evaluation-order settings.

Before retaining this experiment, run the existing collector/traversal suites
with and without the environment flag, byte-check complete compiler outputs,
and measure paired end-to-end workloads on a quiet machine. The captured
profile's 59–66% cumulative traversal time includes pass work; it is not an
estimate of removable dispatch overhead. `GetNext` already is one pointer load,
so this experiment does not eliminate its memory-locality cost.
