# Compact traversal scope stacks: source proof and fresh profile owners

Frozen comparison source: a3a2f8f3293b293571fb835f6234f7deffc2fac3. Common test commit 3184013c; stack4 67ccaeb5 and stack8 edb1ef24. Each production commit changes only the private scopesBuf dimension. No builds/tests have run on these candidates.

## Observable behavior

The repository has no reference to scopesBuf outside its field and initialization, and no cap(scopes), element-pointer exposure, reflection, or unsafe traversal-layout dependency. scopes is private. All public scope root/depth/nearest-parent queries read logical length and entries. PushScope and PushScopeOfScopeWithQuietly append the same node/scope pair before EnterScope, which may materialize scopes. PopScope invokes ExitScope before clearing the exact final entry and shrinking. Quiet operations only suppress callbacks. GetAbstractScopeAt recursively obtains the identical parent then stores the exact returned scope by index. Go append copies interface/node values on growth, retaining Scope/Var identity; no internal element address survives a normal callback or growth boundary.

Reduced inline capacity changes allocation timing, growth counts and GC retention only, not public capacity or semantic values. A custom ScopeCreator can reenter the existing public traversal methods; any original element-pointer/reentrant-growth limitation already exists at depth16 and is not relied on by supported stock traversal. Tests cover valid callback queries and materialization; they do not claim new reentrant traversal semantics.

After overflow, old inline entries can retain references until traversal collection, while active heap entries are cleared on pop. This already exists in the 16-entry control. A smaller inline buffer retains fewer such stale inline references but may allocate more/larger heap backing arrays. No pool, reuse window or global production state is introduced.

## Layout model, not measured RSS

scopeStackEntry is one node pointer plus AbstractScope interface =24 bytes on linux/amd64. The rest of NodeTraversal occupies160 bytes. Control16:544-byte logical struct, expected576-byte allocator class; stack8:352 bytes/class352; stack4:256 bytes/class256. Construction-only savings are224/320 bytes respectively before overflow allocation and retention. These models do not establish end-to-end gains.

Common tests compare full pre/post/Enter/Exit/scope-creation order and root/depth/local declaration/ancestor Var identities to a separate nongrowing256-entry stack on the exact same AST at nesting depths0,3,7,15,33,80. A separate quiet prebuilt-scope test crosses all capacities and checks identity/depth, callback suppression, and zeroed popped backing slots. Existing NodeTraversal tests remain unchanged. Suggested central regex: ^TestNodeTraversal (including new ^TestNodeTraversalStackCapacity).

## Existing fresh ff06 allocation samples

Profiles predate the combined candidate and have all new experiments OFF. Each sample is assigned once within a family. Nearest concrete pass excludes generic traversal, collector and compiler wrappers; outer ancestry may be truncated. Callback and pass tables are alternative overlapping partitions, never additive. Raw per-owner/callback/top-stack data: latest-main-profiles-v2/allocation-owners.json.

| Task | Traversal constructor sampled MiB | SSC creation sampled MiB |
|---|---:|---:|
| typed-app | 50.53 | 149.44 |
| bundle | 390.71 | 694.94 |
| downlevel-whitespace | 4.50 | 49.68 |
| downlevel-simple | 71.54 | 188.40 |
| checks-only | 41.52 | 53.57 |

Bundle traversal construction: FSIV266.65MiB (68.2% of390.71), coalesce42.52MiB(10.9%), DAE13.01MiB(3.3%). Nearest concrete callbacks: GatherCandidates.Visit159.59MiB, MustBeReachingVariableDef.ComputeDependence42.02MiB, CFA.computeCfg39.02MiB, coalesce.EnterScopeWithCfg32.52MiB, Candidate.GetDefinition26.51MiB.

The existing lazy cardinality quotient changes whether/when backward MaybeReachingVariableUse is built. It does not eliminate the preceding forward reaching-def, GatherCandidates per-CFG-node traversal, Candidate.GetDefinition, or shared CFA/coalesce traversal constructors. Therefore the largest constructor sites are complementary to the quotient, rather than mostly already addressed by it. This conclusion concerns source coverage; the new quotient's measured byte savings still need its own activated run.

Bundle SSC creation is distributed: RUC87.17MiB(12.5%), FSIV77.08(11.1%), InlineVariables61.64(8.9%), InlineObjectLiterals60.64(8.7%), InlineFunctions46.63(6.7%), Normalize32.07(4.6%). Nearest callbacks identify RUC.TraverseFunction65.10MiB and escapedFinder.Visit36.51MiB. No single existing quotient covers most of these declarations. Any cross-pass reuse still needs declaration freshness/change-proof; the compact inline buffer requires none of that complexity.

Checks constructors mainly originate in TypedScopeCreator_AbstractScopeBuilderBase.Build27.01MiB and CFA.computeCfg12.51MiB, whereas bundle constructors mostly belong to FSIV. This supports measuring both typed/checks and optimization-heavy workloads.
