# Private candidate-discovery kernel

This source-only experiment is disabled unless
`JSGO_EXPERIMENT_INLINE_CANDIDATE_KERNEL=1`. It has not been built, tested, or
benchmarked yet. It uses the existing `stockReferenceDiscoverySafe` guard:
private builtin name supplier, current compiler JS root, stock conventions,
side-effect analyzer, observers/reporting, and colors/no raw-type callbacks.
Public inliners, custom callbacks, custom conventions, and noncurrent roots use
the original traversal. Reference discovery still uses its original real lazy
scopes, with or without the existing optional itinerary.

`FindCandidateFunctions.Visit` and both of its helpers read exactly
`InGlobalHoistScope` (the current hoist root is nil) and `GetChunk` (the canonical
input lookup after the original script-cache reset). Their transitive function
classification, size checks, annotation checks, and state updates receive no
NodeTraversal. They never call GetScope, GetScopeRoot, GetScopeDepth, or the
scope creator. Only VAR/LET/CONST/FUNCTION/CALL/OPTCHAIN_CALL can change candidate
state; every other Visit is a no-op on a valid compiler AST.

The private kernel keeps the original live recursive ordered handler graph,
saved-next semantics, all currentNode assignments, script interruption/cache,
function/hoist/change context, and Try/ThrowUnexpectedException boundary. It
omits real lexical scope storage and inherited always-true preorder dispatch.
When the reference itinerary is enabled, known scope roots are emitted at the
original enter/exit boundaries directly, without scope stack writes or recorder
interface calls. Class extends/name/computed-key/field/static ordering and
optional default/destructuring execution order remain unchanged. It retains
NodeTraversal as a context object so GetChunk and exception formatting are the
same implementations. It does not create an alternative reference solver.

Tests compare every postorder event/context and every itinerary opcode/root
against the original recorder, replay through the unchanged real-scope engine,
compare complete pass output/AST/diagnostics/change stamps, assert actual
private dispatch and public fallback, compare malformed-class exception context
and tape publication, and exercise the actual default factory in a full stock
compiler lifecycle with source maps. Test instrumentation is caller-owned and
nil in production. All tests remain unrun pending the central build slot.

The last labeled bundle profile attributed 10.9% of CPU to InlineFunctions;
that is the entire semantic pass, not removable dispatch overhead. Older generic
TraverseBranch flat attribution was about 2.1% whitespace / 3.2% typed / 4.2%
bundle across ALL generic traversals. Neither figure establishes this kernel's
speedup. Discovery remains a full linked-tree scan; classifying function bodies,
side effects, cost checks and reference scope construction remain real work.
This experiment targets broader no-op callback/scope machinery than intrinsic
leaf filtering or the earlier same-state concrete callback generator, but no
win should be inferred before real CLI activation and isolated timing.

RUC already has its own concrete TraverseNode and TypeInference has a concrete
recursive token switch plus an existing private concrete CFG solver. Applying a
NodeTraversal specialization to those paths would miss their main work. Their
remaining interfaces are removable/variable records and FlowScope/JSType calls;
those need separate representation/semantic evidence, not another generic
walker rewrite. Go generic callback constraints alone do not ensure monomorphic
code, so this kernel uses concrete calls rather than per-node type assertions.
