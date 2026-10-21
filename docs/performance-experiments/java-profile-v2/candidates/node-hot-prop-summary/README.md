Prepared contingency: immutable summaries of two Node properties
===============================================================

Status: source-only; not compiled, run, or measured. No recommendation to adopt
this candidate yet. The accumulator confirmation remains the sole active runtime.
The earlier `../node-scope-flag-cache` draft is retained separately and deferred.

The change adds one final byte to the private immutable `PropListItem` base class.
It summarizes the nearest ARROW_FN and SYNTHETIC properties in that item's suffix:
two bits each, encoding absent (0), nonzero integer (1), or object (2). A new item
inherits its tail's summary, then replaces its own property's slot. Only
`isArrowFunction()` and `isSyntheticBlock()` use the summary; generic property
reads and all property-list writes stay unchanged. The patch contains one source
file, and no Node instance field or invalidation machinery.

Correctness argument
--------------------

The empty suffix has both properties absent. An unrelated item preserves the
tail's states; an item for either property replaces that state with its own kind.
By induction this selects the same nearest item as `lookupProperty()`, including
duplicate lists assembled through reflection. Every successfully constructed
IntPropListItem stores a nonzero integer, including negative integers. Zero is
rejected by its constructor and public zero writes remove the item. The object
state calls the original generic getter, preserving UnsupportedOperationException
instead of treating an object as a boolean. The FUNCTION gate remains first in
`isArrowFunction()`; `isSyntheticBlock()` remains ungated by token.

There are only two concrete item subclasses, both final. Their value, type, and
next fields are final. Constructors do not publish `this` or invoke callbacks;
the original nonnull/nonzero validation remains in its original order. `chain()`
constructs the same kind and value with a new tail. Removal either returns an
existing suffix or reconstructs the prefix through these constructors. Summary
values are at most 10, so byte sign extension cannot affect the two reads. The
pinned enum ordinals used here fit the existing byte representation.

All 12 head stores in the pinned Node.java were audited and remain unchanged:

| Original line(s) | Path | Why the summary remains valid |
| --- | --- | --- |
| 1030 | clonePropsFrom | Shares the immutable head and suffix |
| 1186, 1188 | putProp | Rebuilds/removes, then constructs an object item |
| 1198, 1200 | putIntProp | Rebuilds/removes, then constructs a nonzero int item |
| 1205 | removeProp | Returns/reconstructs an immutable suffix/prefix |
| 1345, 1352, 1358 | deserializeProperties | Constructs boolean, constant, and side-effect items |
| 1726 | setStaticSourceFileFrom | Shares an existing source-file tail |
| 1991 | setPropListHead | Replaces the head with an already valid immutable list |
| 2762 | copyBaseNodeFields | Shares the head for normal and specialized clones |

The private head has no external stores, and there are no in-place item updates
in the extracted sources. Deep clones of JSDoc type expressions and closure
shadows use the same ordinary construction/replacement paths. Deserialization
may fail after adding some items; those already-created suffixes remain valid.
Read-only concurrent access writes no cache state. This adds no guarantee for
concurrent Node mutation or reflection/Unsafe changes to existing final fields.
The probe does construct valid immutable duplicate lists through reflection.

Public and package method signatures remain unchanged. The private base-item
constructor changes, so any overlay must include the complete newly compiled
Node class family, including all three property-list classes. There is no Java
serialization contract for Node being changed; TypedAST property serialization,
property order, debug text, and validation paths are untouched.

Evidence and limits
-------------------

The completed steady BigQuery profile has 26/2221 sampled CPU leaves in
`lookupProperty()` reached through `isArrowFunction()` (1.1706%), and 21/2221
through `isSyntheticBlock()` (0.9455%), the latter via `createsBlockScope()`.
These are caller-prefix observations from a sampled profile, not call counts or
an expected speedup. `profile-evidence.tsv` retains those four rows. The summary
read still costs work, short property lists may already be cheap, and inlining
or generated code could change the result.

Every item construction now reads its tail summary, compares its type with the
two hot properties, and stores a summary, including construction of unrelated
properties and rebuilt prefixes. That cost can outweigh the saved walks. This
is deliberately a two-property experiment, not a generic property cache.

Adding a byte may fit existing item padding, but object sizes and field placement
have NOT been measured. Neither zero footprint growth nor a particular item
size is claimed. Node/StringNode receive no direct field; growth of either
property-item subclass would still increase compiler memory. Before promotion,
measure both item layouts under the actual JVM flags, then require output and
diagnostic byte equality and a worthwhile end-to-end result with construction
cost included. This patch has no measured benefit yet.

Focused probe and exact-source checks
------------------------------------

`NodeHotPropSummaryProbe` compares both dedicated getters against unchanged
generic reads, repeated three times per observation. It covers absent/nonzero/
zero/object values, public setters and token gates, every Token plus null,
unrelated prefix rebuilds, duplicate properties, sharing, all specialized clones,
deep type/shadow clones, source-tail transfers, and successful/partially failing
deserialization. It records property order, serialized bits, and validation
messages. Candidate-only `--require-summary` reads the final byte reflectively
and verifies it against the generic oracle; successful stdout should match STOCK
and CONTROL exactly. This is a source-reviewed fixture, not a runtime result.

`prepare-source.py` creates the normalized Closure-repository patch and records
all head stores, narrow source hashes, and exact forward/reverse application in
a temporary tree containing only the original Node.java. `patch-check.json` is
the receipt. No archives, compiler, JVM, or benchmark are needed for this check.

After the runtime owner schedules an exclusive slot, use fresh build/results
directories and the existing recorded harness (commands below were NOT run):

```sh
cd /workspace/jsgo-java-experiments/profile-v2
python3 build-candidate.py node-hot-prop-summary \
  --source-root candidates/node-hot-prop-summary/src
python3 run-probes.py \
  --manifest build/node-hot-prop-summary/manifest.json \
  --probe-source candidates/node-hot-prop-summary/probe/com/google/javascript/rhino/NodeHotPropSummaryProbe.java \
  --probe-class com.google.javascript.rhino.NodeHotPropSummaryProbe \
  --candidate-arg=--require-summary --compare-stdout \
  --origin-class com.google.javascript.rhino.Node \
  --origin-class 'com.google.javascript.rhino.Node$PropListItem' \
  --origin-class 'com.google.javascript.rhino.Node$IntPropListItem' \
  --origin-class 'com.google.javascript.rhino.Node$ObjectPropListItem' \
  --out results/node-hot-prop-summary-probes-attempt-1
```

Retain any first fixture/build failure and fix it in a separately recorded source
version. A passing focused probe alone does not establish compiler byte identity
or a performance win.
