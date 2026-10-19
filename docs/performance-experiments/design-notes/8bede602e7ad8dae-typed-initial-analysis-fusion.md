# Private initial typed-scope analysis fusion

`JSGO_EXPERIMENT_TYPED_INITIAL_FUSION=1` lets the private default inference factory
combine TypedScopeCreator's first-order function analysis and nonnullable
enum/typedef identification into one AST traversal. Public TypeInferencePass,
TypedScopeCreator and callback constructors retain their original two walks.
The private enablement window closes on return and panic; it does not permanently
change a creator returned earlier through a public API.

Both analyses still run their original methods on the same memoized untyped
scopes. The first analysis writes only reserved names, functions with returns,
assignment counts and escaped names. It does not read or write JSTypeRegistry.
The second writes only the registry's nonnullable-name sets and consults AST
annotations, existing constant flags and untyped declarations. It does not read
any first-order output. First.EnterScope materializes declarations before either
callback's visits. The stock syntactic creator reads AST declarations/input
bindings, and its default redeclaration handler is empty. Earlier second-analysis
registry writes therefore cannot change later first-analysis fact collection or
stock scope construction; earlier first-analysis writes cannot change the second.

Each callback's own preorder/postorder and scoped Enter/Exit sequence is preserved.
The original first analysis uses TraverseRoots and omits the common ROOT callback;
the second uses Traverse(root) and includes it. The fused walk uses Traverse(root)
and explicitly omits the first callback's common-ROOT preorder/postorder. It still
delivers the second callback's final ROOT visit. The first's root ExitScope is an
empty method; its root transient CurrentNode differs, so tests compare the actual
scope root that identifies each Enter/Exit event. No public callback observes the
combined traversal.

An unexpected malformed-AST internal failure can leave a different partial
registry: the fused walk may have registered enums before the original first
walk would fail. Public/custom paths retain the original schedule. This private
default-factory experiment makes no promise about inspectable partial registry
state after a failed internal compilation; it must not publish a CLI artifact
from that failed execution.

Selection requires the exact stock compiler and current two-child root pair,
stock coding convention including copied-proxy self fences, known compiler
change/error observers, the compiler-owned type registry and standard Rhino
reporter, no custom passes, no debug/tracing/determinism observers and no errors.
Unknown/custom contexts keep the original two walks. In particular, a custom
redeclaration handler can inspect registry state while scopes are built; fusing
would change what it sees and is prohibited. Type resolution and property caches
remain sequential and unchanged: TypedScopeCreator.FROZEN does not make the
entire JSType graph immutable.

This structurally removes one full AST traversal and its traversal bookkeeping.
It does not reuse cached parse/type results or omit either analysis. Existing
filtered samples put CreateInitialScope at roughly .14s typed/.25s checks, with
the enum callback itself much smaller; attribution loses some ancestor context
at deep stacks, so neither number is a measured improvement or guaranteed ceiling.

Optional private test probes capture both callbacks' exact own event sequences
and pre-build first-order facts. Differential fixtures include enum/typedef
aliases, nested functions/default parameters, class fields/static blocks,
catch/block scopes, ES imports/exports and goog.loadModule. Full inference tests
compare typed AST, ordered scopes/Vars, inferred/escaped/assignment flags, output
and diagnostics, plus public/custom fallback and panic lifetime.

Source-only prototype: formatted and reviewed, not yet built, tested or timed.
