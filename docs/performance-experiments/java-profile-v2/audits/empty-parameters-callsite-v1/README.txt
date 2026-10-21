Empty-formals shortcut: compiler-callsite review

Status: source review complete; runtime/build/byte-equivalence/performance work
belongs to the root agent. This audit did not run Java or rebuild any source.
No original candidate, patch, or probe was modified. source-review.json records
the exact frozen inputs and confirms that the candidate differs from its pinned
control by only the three-line guard after getFunctionParameterSet(fnNode).

Candidate

  candidates/scope-maps/empty-parameters-only/src/com/google/javascript/jscomp/FunctionArgumentInjector.java
  candidates/scope-maps/empty-parameters-only/candidate.patch

  if (names.isEmpty()) {
    return names;
  }

Proof for compiler-owned, well-formed ASTs

getFunctionParameterSet still executes before the guard. Every ordinary formal
adds a string, even an empty string. A rest formal adds REST_MARKER. Default,
object-pattern, and array-pattern formals retain their existing rejection.
Consequently an empty result here means there are zero formal parameter nodes;
rest parameters do not take the shortcut.

The recursive analysis has only one local insertion, guarded by membership in
names. An empty names set prevents every insertion. Inductively, each descendant
returns an empty set and every union is empty, including descendants within
functions, arrows, or nonstatic class field initializers. The outer result is
therefore exactly the same empty immutable set of names and has the same order.

The skipped code only reads Node predicates/accessors, changes a local closure
boolean, and allocates/merges local LinkedHashSets. getString reads StringNode
storage. canNameValueChange is already short-circuited by names.contains and
never executes when names is empty. There are no AST writes, compiler callbacks,
diagnostics, scope creation, analysis-state changes, or required validation in
this private walk. The cost removed scales with body size; neither its frequency
nor an end-to-end speedup is established by this source proof.

Actual production callers in pinned sources

1. InlineFunctions.updateFunctionStateForInlining, line 321, after parameter and
   candidate screening. The empty result preserves the direct-inline decision
   and names-to-alias state. Separate receiver, inner-closure, and local-name
   checks at lines 327 onward still execute unchanged.

2. FunctionToBlockMutator.mutateInternal, line 152, on the cloned function after
   optional local-name rewriting. The complete call-argument map is still built
   at lines 153-155; gatherCallArgumentsNeedingTemps still runs for a nonempty
   argument map at lines 158-163. Argument aliasing/injection, loop declaration
   repairs, and return rewriting are unchanged. Extra actual arguments and this
   remain represented by the map's anonymous placeholders and THIS_MARKER, so
   zero formal parameters does not skip their side-effect handling.

The latter route is reached through FunctionInjector and recognized IIFE/UMD
rewrites in ProcessCommonJSModules. Some J2CL bootstrap/property inlining routes
run before normalization. This proof does not require normalization; it requires
the well-formed FUNCTION/PARAM_LIST/body nodes supplied by those compiler paths.

Boundary

This is not universal exception equivalence for arbitrary malformed external
ASTs. For example, a raw new Node(Token.NAME) is not a StringNode: the old empty-
formal traversal throws ClassCastException at getString, whereas the guard never
visits it. Corrupt links/properties can likewise hide traversal failures. The
helper is package-private and its private walk is not a validator or diagnostic
pass. This boundary does not change the proved result or required checks for
compiler-owned, well-formed ASTs; the explicit boundary probe retains it for
review. Null/nonfunction/missing-parameter failures before the guard and malformed
body failures with a real/rest formal remain checked.

Probe source review

ModifiedParametersProbe compares ordered results with a copied pinned algorithm
and checks AST stability, unsupported formals, nested closures/fields, ordinary
assignments, and rest handling. EmptyParametersProbe adds no-formal bodies,
arguments/this, generator, repeated-body, early-error, and explicit malformed-
body boundary cases. Both frozen probes were already compiled and passed for
STOCK/CONTROL/accumulator in results/accumulator-probes-v1/attempt-2; that prior
result is not a runtime validation of the new empty-parameters candidate.

Reproduction after build/empty-parameters/manifest.json exists

These are two separate runner invocations. Only EmptyParametersProbe accepts
the candidate boundary argument; --candidate-arg would apply to every selected
probe in one invocation. Use fresh output directories, retaining every failed
attempt. Commands below were not executed by this source audit.

python3 /workspace/jsgo-java-experiments/profile-v2/run-probes.py \
  --manifest /workspace/jsgo-java-experiments/profile-v2/build/empty-parameters/manifest.json \
  --probe-source /workspace/jsgo-java-experiments/profile-v2/candidates/scope-maps/probe/com/google/javascript/jscomp/ModifiedParametersProbe.java \
  --probe-class com.google.javascript.jscomp.ModifiedParametersProbe \
  --compare-stdout \
  --out /workspace/jsgo-java-experiments/profile-v2/results/empty-parameters-modified-probe-reproduction-1

python3 /workspace/jsgo-java-experiments/profile-v2/run-probes.py \
  --manifest /workspace/jsgo-java-experiments/profile-v2/build/empty-parameters/manifest.json \
  --probe-source /workspace/jsgo-java-experiments/profile-v2/candidates/scope-maps/empty-parameters-only/probe/com/google/javascript/jscomp/EmptyParametersProbe.java \
  --probe-class com.google.javascript.jscomp.EmptyParametersProbe \
  --candidate-arg=--expect-shortcut \
  --compare-stdout \
  --out /workspace/jsgo-java-experiments/profile-v2/results/empty-parameters-boundary-probe-reproduction-1

The runner records source hashes, shading/build arguments, class origins, all
STOCK/CONTROL/candidate outputs and statuses, and postflight integrity. Root's
compiler-byte gates must then compare all existing output/diagnostic artifacts
with the same pinned inputs and options before any timing result is promoted.
Neither semantic probes nor this proof establish a worthwhile performance win.
