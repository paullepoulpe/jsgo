# Frozen fresh-template focused overlay

Candidate: `27c88f9a`, integrated source: `c7998d45` (V2 `00f707bb`, fixture-only fix `22fbda3a`).

Regenerate with `python3 prepare.py --root /workspace/jsgo-structural-bitset-template`.
This preserves the package import path and all production files. It replaces only excluded test declarations with an empty `package jscomp`; retained helpers are copied verbatim from original test sources. No production stubs are present.

## Candidate changed files

- `docs/experiments/fresh-template-values.md`
- `internal/jscomp/fresh_template_values.go`
- `internal/jscomp/fresh_template_values_test.go`
- `internal/jscomp/native_type_inference.go`
- `internal/jscomp/type_inference.go`
- `internal/jscomp/type_inference_pass.go`
- `internal/jscomp/typed_initial_analysis_fusion.go`
- `internal/rhino/jstype_fresh_template_values.go`
- `internal/rhino/jstype_fresh_template_values_test.go`

## Retained tests

Same base 35 focused test files, plus the fresh-template constructor and guard test files:

- `access_control_property_facts_test.go`
- `coding_convention_safety_test.go`
- `contextual_inverter_compact_test.go`
- `contextual_inverter_noop_test.go`
- `deps_module_names_optimization_test.go`
- `empty_block_scoped_test.go`
- `empty_denormalize_test.go`
- `fresh_template_guard_test.go`
- `fresh_template_values_test.go`
- `fsiv_scope_sharing_test.go`
- `inline_argument_facts_test.go`
- `inline_callee_facts_test.go`
- `inline_direct_reference_worklist_test.go`
- `inline_discovery_chunks_test.go`
- `inline_reference_discovery_tape_test.go`
- `native_dependency_scan_test.go`
- `native_leaf_traversal_test.go`
- `native_live_variables_test.go`
- `native_output_buffer_test.go`
- `native_typed_variables_test.go`
- `parallel_coalesce_test.go`
- `parallel_fsiv_test.go`
- `parallel_peephole_test.go`
- `reference_batches_lifetime_test.go`
- `reference_batches_test.go`
- `reference_collector_native_storage_test.go`
- `remove_unused_code_work_elision_test.go`
- `remove_unused_intrusive_test.go`
- `remove_unused_tasks_test.go`
- `scope_first_allocation_test.go`
- `scope_recipes_test.go`
- `source_file_cli_lock_test.go`
- `stock_scope_insertion_test.go`
- `test_compiler_test.go`
- `testing_test_externs_builder_test.go`
- `typed_inference_scope_schedule_test.go`
- `typed_initial_analysis_fusion_test.go`

## Required helpers

- `typed_initial_analysis_fusion_test.go`: compiler initialization, typed-node IDs, type/scope/output/diagnostic snapshots, unknown compiler wrapper.
- `typed_inference_scope_schedule_test.go`: unknown reverse-interpreter and warning guard wrappers.
- `parallel_fsiv_test.go`: unknown observer wrapper.
- Original function-argument injector name/call/function helpers, coreG3Collector, typedScopeDump variable/class/type/position functions, removeUnusedOracle options/compiler/prepare/extern helpers, original CompilerTestCase extern definitions, and scopeTestNewCompiler are extracted unchanged by prepare.py.

## Separate Rhino validation

Normal and race both passed the same 33 top-level tests: five new constructor families, sixteen TemplateTypeMap tests, twelve TemplatizedType tests. These exercise original and fresh TemplateTypeReplacer/resolver behavior; no separate TemplateTypeReplacer test file exists. Logs are `../tests-fresh-template-rhino-normal.log` and `../tests-fresh-template-rhino-race.log`.

The jscomp focused overlay has not been built or executed by this agent; root owns compiler validation and measurements.

The V2 guard and phase hooks, bounded counter plan, additive Monitor API and new tests are source-frozen for central validation. Monitor normal passed five tests; race was still running at regeneration. No production file is replaced by this overlay.
