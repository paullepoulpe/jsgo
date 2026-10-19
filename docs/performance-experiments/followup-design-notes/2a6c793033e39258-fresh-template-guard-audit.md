# Fresh-template admission guard audit

Frozen candidate: `27c88f9a`, integrated combined source: `3ebe6449`. Source audit only; no implementation changes, guard profiling, new captures, or timings performed.

## Finding

`TypeInferencePass.InferScope` calls `currentFreshTemplateValues()` once per inferred CFG root. While the mode is present, this repeats all of `freshTemplateValuesSafe`, including scope-schedule, initial-analysis, observer, convention, warning, reporter, root and input-ownership admission. The stock private CLI window can avoid repeating most of this work after an additional validation at the phase boundary. Entry-only caching would not preserve all current error/probe behavior.

## Actual repeated work

For each of S inferred scopes, current code performs:

- `nativeReverseInterpreterSafe`: checks the concrete stock chain/self/next links.
- `typedInitialAnalysisFusionSafe`: root shape and identity, options/custom-pass/debug settings, three recursive convention-chain validations across its own and observer checks, the change-handler list, error-manager chain, reporter/registry identities and error count.
- Recursive checks of both compiler and option warning-guard composition trees, each depth limited to 16.
- Typed-AST filesystem/runtime exclusion and iteration of `inputsById.Values()` until the first owned source loader.

The rough structural cost is `S * O(C + H + W + B)`, where C is convention-chain work, H is admitted change/error-manager handler work, W is both warning trees, and B is JavaHashMap iteration to the first owned source. `Values()` is an `iter.Seq`, not an eager copied collection. All ordinary CLI source/extern files eligible for ownership are wrapped, so the first-owned search can be short; this is not evidence of an S-by-input-count full scan or of guard heap allocations.

No guard cost has been measured. It may be insignificant, or may offset a constructor improvement. The prior 1.6% checks-only CPU attributed to the whole original factory family is an upper bound on removing that original factory work; it is NOT a bound on newly added admission overhead. The factory attribution includes required resolution/replacer work that this prototype retains.

## Ownership and callback proof

`ownCLISourceFiles` admits only the exact stock `CommandLineRunner`/base pair, not test mode. It runs with `c.GetRoot()==nil`, fresh known source loaders, no unknown observers/custom passes, and no saved/typed-AST inputs. `cliSourceFileOwnership.release` restores surviving source loaders at the DoRun return/panic boundary. Thus an external caller cannot retain the private loader admission marker for a later public pass invocation after DoRun returns.

The existing full entry gate additionally requires exact compiler/creator/registry/root identities, stock RAI and convention chains, stock reporter/error manager, known change handlers, and explicitly admitted warning guards. The stock scope-building pass and resolver construct/resolve concrete compiler types and update declaration/type metadata. They do not install new options, observer, RAI or input-loader callbacks. Stock `Compiler.Report` consults the admitted warning guards and stores diagnostics through the admitted error manager; stream/formatter output occurs later during `GenerateReport`, outside inference. No arbitrary output-stream callback runs between inferred scopes through this route.

A real nonconstant admission condition is error count. Scope construction and earlier inference scopes can add errors. The current repeated guard then disables the fresh constructor for subsequent scopes. This behavior should be retained rather than silently extending the experimental dispatch window to errored scopes.

The relevant private callback exception is `TypedScopeCreator.initialAnalysisProbe.afterFacts`. It can mutate configuration or inject other adversarial state in internal tests. It is absent in stock CLI. Retaining the original slow full-check behavior whenever this probe is nonnil avoids weakening those internal mutation cases. `nativeFlowObserver` is already rejected by the entry gate. Schedule/fresh probes only retain/count observations and provide no callback.

Concurrent external mutation is already outside the supported unsynchronized compiler lifecycle; repeated checks do not make it race-safe. This proposal adds no synchronization and claims none.

## Recommended future experiment

Do not alter frozen 3ebe6449 before its first macro result.

If guard optimization is warranted afterward:

1. Keep the existing full admission gate at private-window entry.
2. Capture an invocation-owned token with exact owner/root/creator/registry identities and the optional counter probe.
3. Run one full revalidation immediately after `FinishAndFreeze`, before either the original second traversal or schedule replay. Only then arm the cheap scope-use token.
4. In the ordinary nil-initialAnalysisProbe path, require the token and exact current compiler/creator/registry/root identities plus `owner.GetErrorCount()==0` at each InferScope. These are bounded identity/count checks.
5. When `initialAnalysisProbe!=nil`, preserve the current full per-scope guard to support internal adversarial callbacks unchanged.
6. Keep the existing deferred restoration on success and panic. Public constructors, public InferAllScopes and direct InferScope calls outside the private window continue with nil mode. A failed boundary validation disables only this experimental constructor; original inference still executes.

This replaces S full composite validations with approximately two full validations plus S bounded identity/error checks. Four additional captured identity pointers are about 32 payload bytes once per invocation before allocator size-class rounding; exact final struct layout depends on how the existing mode/probe pointer is arranged. No per-scope retained state or global cache is required.

## Required focused coverage before changing guards

- A phase1 declaration/type error disables fresh dispatch after freeze exactly as frozen source does.
- An error emitted by an earlier real inference scope stops fresh dispatch for later scopes in the same order as frozen source.
- Private afterFacts mutations of RAI, warnings and root retain original slow fallback behavior.
- Direct/public InferScope after private success and panic sees no active mode.
- Original/fresh typed AST, scope, diagnostic and output differentials cover both native and generic CFG engines.
- Malformed roots and phase1/constructor failures preserve original panic/cleanup behavior.

## Validation already completed for unchanged candidate

The focused Rhino constructor/type-map/type-specialization selection passed 33 top-level tests normally and the same 33 with `-race` (5 new constructor families, 16 existing TemplateTypeMap tests, 12 existing TemplatizedType tests). The constructor tests include arbitrary padding-bound callbacks, input mutation, reentrancy and exception order. They do not measure admission-guard cost.
