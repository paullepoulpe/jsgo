# Mechanical Java-Prep Infrastructure

This directory contains the 70 mechanical Java-prep transformations (PRs 1–70) applied to upstream Google Closure Compiler (`v20261104`) prior to mechanical transpilation (`j2go`).

## Overview of Transformations

| Patch | Name | Description |
|---|---|---|
| `01_dereflect_command_line_runner.patch` | De-reflect CommandLineRunner | Replaces args4j reflection with deterministic switch dispatch in `CommandLineRunner.java`. |
| `02_dereflect_gson_serialization.patch` | De-reflect Gson | Replaces reflective Gson serialization with explicit streaming `JsonReader`/`JsonWriter`. |
| `03_desugar_stream_apis.patch` | Desugar Stream APIs | Converts `definitions.stream().allMatch(...)` and `map(...)` to standard loops. |
| `04_replace_guava_immutable_list.patch` | Replace Guava `toImmutableList` | Lowers stream collectors to standard `ArrayList` and `ImmutableList.copyOf`. |
| `05_replace_guava_multimap.patch` | Replace Guava Multimap | Uses `SimpleMultimap.java` with standard `computeIfAbsent`. |
| `06_replace_guava_table.patch` | Replace Guava Table | Uses `SimpleTable.java` backed by nested `LinkedHashMap`. |
| `07_convert_autovalue_to_pojos.patch` | AutoValue to POJOs | Replaces `@AutoValue` annotations with clean immutable classes. |
| `08_desugar_try_with_resources.patch` | Desugar try-with-resources | Desugars `try (...)` statements to explicit `try-finally` blocks. |
| `09_replace_string_format_join.patch` | String.format to StringBuilder | Replaces `String.format` call sites with direct string concatenations. |
| `10_replace_optional.patch` | Replace Optional<T> | Replaces `Optional<T>` usage with nullable references. |
| `11_normalize_switch_statements.patch` | Normalize switch statements | Converts Java 14+ switch expressions (`->`) to standard Java switch statements. |
| `12_normalize_varargs_diamond.patch` | Normalize varargs & diamond | Adds explicit generic type arguments to remove type-inference ambiguities. |
| `13_cleanup_annotations_imports.patch` | Clean up annotations & imports | Removes unused third-party annotations and dead imports. |
| `14`–`68` | Representation, concrete types & allocation slimming | Mechanical type devirtualization, unboxed primitives, inline small-collection slots, and pass-local scratch pooling. |
| `69_jstype_featureset_parser_and_ref_slimming.patch` | JSType, FeatureSet, parser & ref slimming | 64-bit `long` bitmask `FeatureSet` + allocation-free `ExpressionDecomposer` IE11 workaround flag; `byte` `JSType.templateCheckResult`, lazy `PrototypeObjectType.NamingExtra`, and `final` `FunctionType` (folding `NoObjectType`/`NoType`); pooled/unboxed `JsDocTokenStream` + 4-slot `JSDocInfo.Builder`; 4-slot `ReferenceCollection`, 2-slot `GlobalNamespace.Name` & `InlineFunctions.FunctionState`, and `MemoizedScopeCreator` capacity hints. |
| `70_jstype_jsdoc_templatemap_and_dfa_unbox.patch` | JSType, JSDoc, TemplateTypeMap & DFA unboxing | Inline `PrototypeObjectType` `className`/`ownerFunction`/`templateParamCount` (eliminating `NamingExtra`); lazy `TemplateTypeMap.resolvedTemplateValues`, `SINGLE_SUBMAP_START`, and direct slice/builder reuse; `JSDocInfo` `TYPE_EXPR_PROPERTIES_MASK` + `hasTypeNodes()` + direct `collectTypeNodes`; `String[]`/`ReservedNames` in `TypedScopeCreator`/`TypedScope`, empty/singleton fast paths in `FunctionType.getAllImplementedInterfaces` & `TemplatizedType`, scratch sets + 4-slot union buffer in `JSTypeReconserializer`/`ColorId`; dense `0..N-1` array lattices + 2-slot `Definition` deps in `MustBeReachingVariableDef` and `MaybeReachingVariableUse`. |

## Pure Java Helpers
- `com/google/javascript/jscomp/SimpleMultimap.java`: Pure-Java multimap with zero external dependencies.
- `com/google/javascript/jscomp/SimpleTable.java`: Pure-Java 2D table with zero external dependencies.

## Usage

```bash
# Verify all patches apply cleanly (dry-run):
./transpiled/prep/apply_prep.sh --check

# Reset .cache/closure-src/ to pristine state, apply all patches, and verify:
./transpiled/prep/apply_prep.sh --reset --verify
```
