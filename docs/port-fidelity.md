# Port Fidelity Report

Generated: auto-computed by tools/portcheck.

## Summary

| Metric | Value |
|--------|-------|
| Total files scored | 232 |
| Average score | 56.9 |
| Std deviation | 12.4 |
| Min score | 0.0 |
| Median score | 57.0 |
| Max score | 90.9 |

## Files by Fidelity Score (worst first)

| Rank | Score | Go File | Java Source | Method Overlap % | Stub Count | LoC Ratio | Notes |
|------|-------|---------|-------------|-----------------|------------|-----------|-------|
| 1 | 0 | `jscomp/default_externs.go` | `` | 0% | 0 | 0.00 | missing Port of header |
| 2 | 28 | `parser/tokentype.go` | `jscomp/parsing/parser/TokenType.java` | 0% | 0 | 1.36 | 17 go-only fns |
| 3 | 31 | `jscomp/closurereverseabstractinterpreter.go` | `jscomp/type/ClosureReverseAbstractInterpreter.java` | 12% | 4 | 1.50 | 4 stubs; 29 go-only fns |
| 4 | 33 | `jscomp/compiler.go` | `jscomp/Compiler.java` | 0% | 0 | 0.04 | 4 go-only fns |
| 5 | 33 | `jscomp/checkjsdoctypes.go` | `jscomp/TypeValidator.java` | 0% | 0 | 1.33 | 23 go-only fns |
| 6 | 33 | `jscomp/checkmissingrequires.go` | `jscomp/CheckMissingRequires.java` | 0% | 0 | 1.32 | 37 go-only fns |
| 7 | 34 | `jscomp/var.go` | `jscomp/Var.java` | 0% | 0 | 0.77 | 9 go-only fns |
| 8 | 34 | `jscomp/checkunreachablecode.go` | `jscomp/CheckUnreachableCode.java` | 0% | 0 | 1.50 | 8 go-only fns |
| 9 | 35 | `jscomp/checksuper.go` | `jscomp/CheckSuper.java` | 44% | 14 | 1.16 | 14 stubs; 5 go-only fns |
| 10 | 35 | `jscomp/dataflow.go` | `jscomp/DataFlowAnalysis.java` | 0% | 1 | 0.41 | 1 stubs; 4 go-only fns |
| 11 | 35 | `rhino/jstype/jstype.go` | `rhino/jstype/JSType.java` | 15% | 0 | 0.59 | 10 go-only fns |
| 12 | 36 | `jscomp/checksideeffects.go` | `jscomp/CheckSideEffects.java` | 0% | 0 | 1.50 | 8 go-only fns |
| 13 | 37 | `jscomp/checktypeimportcodereferences.go` | `jscomp/CheckTypeImportCodeReferences.java` | 0% | 0 | 1.50 | 7 go-only fns |
| 14 | 37 | `rhino/jstype/registry.go` | `rhino/jstype/JSTypeRegistry.java` | 7% | 0 | 0.20 | 21 go-only fns |
| 15 | 38 | `parser/trees/parse_tree.go` | `jscomp/parsing/parser/trees/ParseTree.java` | 0% | 0 | 0.01 | 1 go-only fns |
| 16 | 39 | `jscomp/checkaccesscontrols.go` | `jscomp/CheckAccessControls.java` | 0% | 0 | 0.83 | 47 go-only fns |
| 17 | 39 | `jscomp/typecheck.go` | `jscomp/TypeCheck.java` | 10% | 2 | 1.50 | 2 stubs; 132 go-only fns |
| 18 | 39 | `diag/diag.go` | `jscomp/JSError.java` | 0% | 0 | 0.03 | 1 go-only fns |
| 19 | 40 | `jscomp/flowscope.go` | `jscomp/type/FlowScope.java` | 13% | 0 | 1.50 | 10 go-only fns |
| 20 | 41 | `jscomp/nodetraversal.go` | `jscomp/NodeTraversal.java` | 12% | 3 | 0.16 | 3 stubs; 3 go-only fns |
| 21 | 41 | `jscomp/astvalidator.go` | `jscomp/AstValidator.java` | 6% | 0 | 0.04 | 6 go-only fns |
| 22 | 41 | `jscomp/compileroptions.go` | `jscomp/CompilerOptions.java` | 0% | 0 | 0.01 | 2 go-only fns |
| 23 | 42 | `jscomp/nodeutil.go` | `jscomp/NodeUtil.java` | 22% | 1 | 0.45 | 1 stubs; 31 go-only fns |
| 24 | 43 | `rhino/jstype/object_type.go` | `rhino/jstype/PrototypeObjectType.java` | 13% | 3 | 0.78 | 3 stubs; 32 go-only fns |
| 25 | 44 | `jscomp/checkconformance.go` | `jscomp/CheckConformance.java` | 5% | 0 | 0.02 | 2 go-only fns |
| 26 | 44 | `jscomp/varrefcheck.go` | `jscomp/VariableReferenceCheck.java` | 3% | 0 | 1.50 | 21 go-only fns |
| 27 | 44 | `jscomp/cfg.go` | `jscomp/ControlFlowGraph.java` | 0% | 0 | 0.32 | 1 go-only fns |
| 28 | 44 | `jscomp/inferconsts.go` | `jscomp/InferConsts.java` | 0% | 0 | 1.50 | 6 go-only fns |
| 29 | 44 | `rhino/jstype/enum_type.go` | `rhino/jstype/EnumType.java` | 3% | 0 | 0.27 | 7 go-only fns |
| 30 | 44 | `jscomp/astfactory.go` | `jscomp/AstFactory.java` | 13% | 0 | 0.07 | 6 go-only fns |
| 31 | 44 | `jscomp/inferjsdocinfo.go` | `jscomp/InferJsDocInfo.java` | 8% | 0 | 1.50 | 16 go-only fns |
| 32 | 44 | `jscomp/scopecreator.go` | `jscomp/ScopeCreator.java` | 0% | 0 | 1.50 | 4 go-only fns |
| 33 | 45 | `jscomp/parsing/jsdoctoken.go` | `jscomp/parsing/JsDocToken.java` | 0% | 0 | 1.50 |  |
| 34 | 45 | `rhino/jstype/proxy_object_type.go` | `rhino/jstype/ProxyObjectType.java` | 62% | 9 | 0.54 | 9 stubs; 5 go-only fns |
| 35 | 45 | `rhino/jstype/object_type_ext.go` | `rhino/jstype/ObjectType.java` | 4% | 0 | 0.07 | 4 go-only fns |
| 36 | 45 | `jscomp/checkduplicatemember.go` | `jscomp/StrictModeCheck.java` | 0% | 0 | 1.29 | 5 go-only fns |
| 37 | 45 | `jscomp/scope.go` | `jscomp/TypedScopeCreator.java` | 2% | 0 | 1.29 | 82 go-only fns |
| 38 | 45 | `jscomp/inference.go` | `jscomp/TypeInference.java` | 6% | 0 | 1.25 | 49 go-only fns |
| 39 | 46 | `jscomp/symboltable.go` | `jscomp/SymbolTable.java` | 51% | 5 | 0.77 | 5 stubs; 54 go-only fns |
| 40 | 46 | `diag/diagnostictype.go` | `jscomp/DiagnosticType.java` | 0% | 0 | 1.50 | 7 go-only fns |
| 41 | 47 | `jscomp/checkclosureimports.go` | `jscomp/CheckClosureImports.java` | 14% | 0 | 1.06 | 22 go-only fns |
| 42 | 48 | `rhino/closure_primitive.go` | `rhino/ClosurePrimitive.java` | 0% | 0 | 0.30 | 1 go-only fns |
| 43 | 49 | `rhino/jsdoc.go` | `rhino/JSDocInfo.java` | 31% | 1 | 0.34 | 1 stubs; 68 go-only fns |
| 44 | 49 | `jscomp/parsing/jsdocinfoparser.go` | `jscomp/parsing/JsDocInfoParser.java` | 27% | 0 | 0.57 | 21 go-only fns |
| 45 | 50 | `jscomp/diagnosticgroup.go` | `jscomp/DiagnosticGroup.java` | 0% | 0 | 0.00 |  |
| 46 | 50 | `jscomp/passnames.go` | `jscomp/PassNames.java` | 0% | 0 | 0.00 |  |
| 47 | 50 | `jscomp/reverseabstractinterpreter.go` | `jscomp/type/ReverseAbstractInterpreter.java` | 0% | 0 | 0.00 |  |
| 48 | 50 | `rhino/jstype/jstype_class.go` | `rhino/jstype/JSTypeClass.java` | 0% | 0 | 0.00 |  |
| 49 | 50 | `rhino/jstype/relationship_visitor.go` | `rhino/jstype/RelationshipVisitor.java` | 0% | 0 | 0.00 |  |
| 50 | 50 | `rhino/static_ref.go` | `rhino/StaticRef.java` | 0% | 0 | 0.00 |  |
| 51 | 50 | `rhino/static_slot.go` | `rhino/StaticSlot.java` | 0% | 0 | 0.00 |  |
| 52 | 50 | `rhino/static_symbol_table.go` | `rhino/StaticSymbolTable.java` | 0% | 0 | 0.00 |  |
| 53 | 50 | `jscomp/controlflowanalysis.go` | `jscomp/ControlFlowAnalysis.java` | 3% | 0 | 0.95 | 28 go-only fns |
| 54 | 50 | `diag/jserror.go` | `jscomp/JSError.java` | 9% | 0 | 0.25 | 4 go-only fns |
| 55 | 50 | `jscomp/abstractcompiler.go` | `jscomp/AbstractCompiler.java` | 13% | 1 | 0.18 | 1 stubs; 2 go-only fns |
| 56 | 50 | `rhino/token.go` | `rhino/Token.java` | 0% | 0 | 1.50 | 1 go-only fns |
| 57 | 51 | `jscomp/defaultpassconfig.go` | `jscomp/DefaultPassConfig.java` | 3% | 0 | 0.11 | 1 go-only fns |
| 58 | 51 | `jscomp/nodenameextractor.go` | `jscomp/NodeNameExtractor.java` | 14% | 0 | 1.50 | 5 go-only fns |
| 59 | 51 | `jscomp/scope_untyped.go` | `jscomp/Scope.java` | 20% | 0 | 1.23 | 12 go-only fns |
| 60 | 51 | `jscomp/abstractscope.go` | `jscomp/AbstractScope.java` | 11% | 0 | 0.33 | 4 go-only fns |
| 61 | 51 | `rhino/jstype/no_resolved_type.go` | `rhino/jstype/NoResolvedType.java` | 6% | 0 | 0.08 | 3 go-only fns |
| 62 | 51 | `jscomp/controlflowgraph.go` | `jscomp/ControlFlowGraph.java` | 0% | 0 | 0.45 | 3 go-only fns |
| 63 | 52 | `rhino/jstype/function_type.go` | `rhino/jstype/FunctionType.java` | 33% | 0 | 0.36 | 28 go-only fns |
| 64 | 52 | `jscomp/chainablereverseabstractinterpreter.go` | `jscomp/type/ChainableReverseAbstractInterpreter.java` | 65% | 3 | 0.74 | 3 stubs; 10 go-only fns |
| 65 | 52 | `parser/trees/update_expression_tree.go` | `jscomp/parsing/parser/trees/UpdateExpressionTree.java` | 0% | 0 | 0.30 | tiny-wrapper(parent bonus 0); 3 go-only fns |
| 66 | 52 | `rhino/jstype/template_type.go` | `rhino/jstype/TemplateType.java` | 9% | 0 | 1.50 | 8 go-only fns |
| 67 | 52 | `rhino/jstype/unknown_type.go` | `rhino/jstype/UnknownType.java` | 0% | 0 | 0.00 | tiny-wrapper(parent bonus 15); 2 go-only fns |
| 68 | 53 | `jscomp/checkjsdoc.go` | `jscomp/CheckJSDoc.java` | 32% | 2 | 1.16 | 2 stubs; 23 go-only fns |
| 69 | 53 | `rhino/jstype/templatized_type.go` | `rhino/jstype/TemplatizedType.java` | 15% | 0 | 0.49 | 9 go-only fns |
| 70 | 53 | `rhino/jstype/all_type.go` | `rhino/jstype/AllType.java` | 1% | 0 | 0.00 | tiny-wrapper(parent bonus 15); 2 go-only fns |
| 71 | 53 | `jscomp/typeinferencepass.go` | `jscomp/TypeInferencePass.java` | 0% | 0 | 1.50 |  |
| 72 | 54 | `rhino/jstype/property.go` | `rhino/jstype/Property.java` | 55% | 2 | 1.50 | 2 stubs; 14 go-only fns |
| 73 | 54 | `jscomp/varcheck.go` | `jscomp/VarCheck.java` | 8% | 0 | 0.48 | 4 go-only fns |
| 74 | 54 | `rhino/jstype/null_type.go` | `rhino/jstype/NullType.java` | 1% | 0 | 0.06 | tiny-wrapper(parent bonus 15); 2 go-only fns |
| 75 | 54 | `rhino/jstype/no_object_type.go` | `rhino/jstype/NoObjectType.java` | 25% | 1 | 0.00 | 1 stubs; 3 go-only fns |
| 76 | 54 | `rhino/jstype/union_type.go` | `rhino/jstype/UnionType.java` | 20% | 0 | 0.43 | 5 go-only fns |
| 77 | 54 | `jscomp/compilerinput.go` | `jscomp/CompilerInput.java` | 25% | 0 | 0.03 | 4 go-only fns |
| 78 | 55 | `rhino/node.go` | `rhino/Node.java` | 51% | 0 | 0.26 | 17 go-only fns |
| 79 | 55 | `parser/trees/formal_parameter_list_tree.go` | `jscomp/parsing/parser/trees/FormalParameterListTree.java` | 0% | 0 | 0.22 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 80 | 55 | `jscomp/typevalidator.go` | `jscomp/TypeValidator.java` | 43% | 2 | 0.72 | 2 stubs; 22 go-only fns |
| 81 | 56 | `parser/trees/comprehension_tree.go` | `jscomp/parsing/parser/trees/ComprehensionTree.java` | 0% | 0 | 0.27 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 82 | 56 | `rhino/jstype/named_type.go` | `rhino/jstype/NamedType.java` | 28% | 2 | 0.51 | 2 stubs; 4 go-only fns |
| 83 | 56 | `rhino/jstype/void_type.go` | `rhino/jstype/VoidType.java` | 1% | 0 | 0.06 | tiny-wrapper(parent bonus 15); 1 go-only fns |
| 84 | 56 | `parser/trees/argument_list_tree.go` | `jscomp/parsing/parser/trees/ArgumentListTree.java` | 0% | 0 | 0.32 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 85 | 56 | `parser/trees/class_declaration_tree.go` | `jscomp/parsing/parser/trees/ClassDeclarationTree.java` | 0% | 0 | 0.32 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 86 | 56 | `parser/trees/import_declaration_tree.go` | `jscomp/parsing/parser/trees/ImportDeclarationTree.java` | 0% | 0 | 0.32 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 87 | 56 | `parser/trees/export_declaration_tree.go` | `jscomp/parsing/parser/trees/ExportDeclarationTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 88 | 56 | `parser/trees/get_accessor_tree.go` | `jscomp/parsing/parser/trees/GetAccessorTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 89 | 56 | `parser/trees/object_literal_expression_tree.go` | `jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 90 | 56 | `parser/trees/optional_member_expression_tree.go` | `jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 91 | 56 | `parser/trees/program_tree.go` | `jscomp/parsing/parser/trees/ProgramTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 92 | 56 | `parser/trees/set_accessor_tree.go` | `jscomp/parsing/parser/trees/SetAccessorTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 93 | 56 | `parser/trees/variable_declaration_list_tree.go` | `jscomp/parsing/parser/trees/VariableDeclarationListTree.java` | 0% | 0 | 0.33 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 94 | 56 | `parser/trees/computed_property_setter_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertySetterTree.java` | 0% | 0 | 0.35 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 95 | 56 | `parser/trees/opt_chain_call_expression_tree.go` | `jscomp/parsing/parser/trees/OptChainCallExpressionTree.java` | 0% | 0 | 0.35 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 96 | 57 | `parser/trees/new_expression_tree.go` | `jscomp/parsing/parser/trees/NewExpressionTree.java` | 0% | 0 | 0.35 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 97 | 57 | `parser/trees/optional_member_lookup_expression_tree.go` | `jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java` | 0% | 0 | 0.35 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 98 | 57 | `parser/parser.go` | `jscomp/parsing/parser/Parser.java` | 8% | 0 | 1.15 | 2 go-only fns |
| 99 | 57 | `rhino/jstype/number_type.go` | `rhino/jstype/NumberType.java` | 1% | 0 | 0.07 | tiny-wrapper(parent bonus 15); 1 go-only fns |
| 100 | 57 | `rhino/jstype/string_type.go` | `rhino/jstype/StringType.java` | 1% | 0 | 0.07 | tiny-wrapper(parent bonus 15); 1 go-only fns |
| 101 | 57 | `parser/trees/case_clause_tree.go` | `jscomp/parsing/parser/trees/CaseClauseTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 102 | 57 | `parser/trees/comprehension_for_tree.go` | `jscomp/parsing/parser/trees/ComprehensionForTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 103 | 57 | `parser/trees/template_literal_expression_tree.go` | `jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 104 | 57 | `rhino/jstype/value_type.go` | `rhino/jstype/ValueType.java` | 8% | 0 | 0.22 | 3 go-only fns |
| 105 | 57 | `rhino/jstype/boolean_type.go` | `rhino/jstype/BooleanType.java` | 1% | 0 | 0.07 | tiny-wrapper(parent bonus 15); 1 go-only fns |
| 106 | 57 | `parser/trees/array_pattern_tree.go` | `jscomp/parsing/parser/trees/ArrayPatternTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 107 | 57 | `parser/trees/comma_expression_tree.go` | `jscomp/parsing/parser/trees/CommaExpressionTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 108 | 57 | `parser/trees/comprehension_if_tree.go` | `jscomp/parsing/parser/trees/ComprehensionIfTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 109 | 57 | `parser/trees/object_pattern_tree.go` | `jscomp/parsing/parser/trees/ObjectPatternTree.java` | 0% | 0 | 0.36 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 110 | 57 | `parser/trees/array_literal_expression_tree.go` | `jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java` | 0% | 0 | 0.38 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 111 | 57 | `parser/trees/export_specifier_tree.go` | `jscomp/parsing/parser/trees/ExportSpecifierTree.java` | 0% | 0 | 0.38 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 112 | 57 | `parser/trees/import_specifier_tree.go` | `jscomp/parsing/parser/trees/ImportSpecifierTree.java` | 0% | 0 | 0.38 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 113 | 57 | `parser/trees/switch_statement_tree.go` | `jscomp/parsing/parser/trees/SwitchStatementTree.java` | 0% | 0 | 0.38 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 114 | 57 | `jscomp/abstractvar.go` | `jscomp/AbstractVar.java` | 18% | 0 | 0.12 | 1 go-only fns |
| 115 | 57 | `parser/trees/block_tree.go` | `jscomp/parsing/parser/trees/BlockTree.java` | 0% | 0 | 0.40 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 116 | 57 | `parser/trees/default_clause_tree.go` | `jscomp/parsing/parser/trees/DefaultClauseTree.java` | 0% | 0 | 0.40 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 117 | 57 | `parser/trees/identifier_expression_tree.go` | `jscomp/parsing/parser/trees/IdentifierExpressionTree.java` | 0% | 0 | 0.40 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 118 | 57 | `parser/trees/literal_expression_tree.go` | `jscomp/parsing/parser/trees/LiteralExpressionTree.java` | 0% | 0 | 0.40 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 119 | 57 | `parser/trees/template_literal_portion_tree.go` | `jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java` | 0% | 0 | 0.40 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 120 | 57 | `parser/trees/computed_property_definition_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 121 | 57 | `parser/trees/default_parameter_tree.go` | `jscomp/parsing/parser/trees/DefaultParameterTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 122 | 57 | `parser/trees/do_while_statement_tree.go` | `jscomp/parsing/parser/trees/DoWhileStatementTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 123 | 57 | `parser/trees/labelled_statement_tree.go` | `jscomp/parsing/parser/trees/LabelledStatementTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 124 | 57 | `parser/trees/member_expression_tree.go` | `jscomp/parsing/parser/trees/MemberExpressionTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 125 | 57 | `parser/trees/member_lookup_expression_tree.go` | `jscomp/parsing/parser/trees/MemberLookupExpressionTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 126 | 57 | `parser/trees/property_name_assignment_tree.go` | `jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 127 | 57 | `parser/trees/unary_expression_tree.go` | `jscomp/parsing/parser/trees/UnaryExpressionTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 128 | 57 | `parser/trees/yield_expression_tree.go` | `jscomp/parsing/parser/trees/YieldExpressionTree.java` | 0% | 0 | 0.42 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 129 | 57 | `parser/trees/binary_operator_tree.go` | `jscomp/parsing/parser/trees/BinaryOperatorTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 130 | 57 | `parser/trees/computed_property_getter_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 131 | 57 | `parser/trees/conditional_expression_tree.go` | `jscomp/parsing/parser/trees/ConditionalExpressionTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 132 | 57 | `parser/trees/debugger_statement_tree.go` | `jscomp/parsing/parser/trees/DebuggerStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 133 | 57 | `parser/trees/empty_statement_tree.go` | `jscomp/parsing/parser/trees/EmptyStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 134 | 57 | `parser/trees/for_await_of_statement_tree.go` | `jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 135 | 57 | `parser/trees/for_in_statement_tree.go` | `jscomp/parsing/parser/trees/ForInStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 136 | 57 | `parser/trees/for_of_statement_tree.go` | `jscomp/parsing/parser/trees/ForOfStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 137 | 57 | `parser/trees/if_statement_tree.go` | `jscomp/parsing/parser/trees/IfStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 138 | 57 | `parser/trees/missing_primary_expression_tree.go` | `jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 139 | 57 | `parser/trees/new_target_expression_tree.go` | `jscomp/parsing/parser/trees/NewTargetExpressionTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 140 | 57 | `parser/trees/null_tree.go` | `jscomp/parsing/parser/trees/NullTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 141 | 57 | `parser/trees/super_expression_tree.go` | `jscomp/parsing/parser/trees/SuperExpressionTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 142 | 57 | `parser/trees/this_expression_tree.go` | `jscomp/parsing/parser/trees/ThisExpressionTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 143 | 57 | `parser/trees/try_statement_tree.go` | `jscomp/parsing/parser/trees/TryStatementTree.java` | 0% | 0 | 0.43 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 144 | 57 | `parser/trees/for_statement_tree.go` | `jscomp/parsing/parser/trees/ForStatementTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 145 | 57 | `parser/trees/await_expression_tree.go` | `jscomp/parsing/parser/trees/AwaitExpressionTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 146 | 57 | `parser/trees/dynamic_import_tree.go` | `jscomp/parsing/parser/trees/DynamicImportTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 147 | 57 | `parser/trees/expression_statement_tree.go` | `jscomp/parsing/parser/trees/ExpressionStatementTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 148 | 57 | `parser/trees/finally_tree.go` | `jscomp/parsing/parser/trees/FinallyTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 149 | 57 | `parser/trees/iter_rest_tree.go` | `jscomp/parsing/parser/trees/IterRestTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 150 | 57 | `parser/trees/iter_spread_tree.go` | `jscomp/parsing/parser/trees/IterSpreadTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 151 | 57 | `parser/trees/object_rest_tree.go` | `jscomp/parsing/parser/trees/ObjectRestTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 152 | 57 | `parser/trees/object_spread_tree.go` | `jscomp/parsing/parser/trees/ObjectSpreadTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 153 | 57 | `parser/trees/paren_expression_tree.go` | `jscomp/parsing/parser/trees/ParenExpressionTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 154 | 57 | `parser/trees/return_statement_tree.go` | `jscomp/parsing/parser/trees/ReturnStatementTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 155 | 57 | `parser/trees/template_substitution_tree.go` | `jscomp/parsing/parser/trees/TemplateSubstitutionTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 156 | 57 | `parser/trees/throw_statement_tree.go` | `jscomp/parsing/parser/trees/ThrowStatementTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 157 | 57 | `parser/trees/variable_statement_tree.go` | `jscomp/parsing/parser/trees/VariableStatementTree.java` | 0% | 0 | 0.44 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 158 | 58 | `parser/trees/call_expression_tree.go` | `jscomp/parsing/parser/trees/CallExpressionTree.java` | 0% | 0 | 0.45 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 159 | 58 | `parser/trees/catch_tree.go` | `jscomp/parsing/parser/trees/CatchTree.java` | 0% | 0 | 0.45 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 160 | 58 | `parser/trees/computed_property_method_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java` | 0% | 0 | 0.45 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 161 | 58 | `parser/trees/variable_declaration_tree.go` | `jscomp/parsing/parser/trees/VariableDeclarationTree.java` | 0% | 0 | 0.45 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 162 | 58 | `parser/trees/while_statement_tree.go` | `jscomp/parsing/parser/trees/WhileStatementTree.java` | 0% | 0 | 0.45 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 163 | 58 | `parser/trees/with_statement_tree.go` | `jscomp/parsing/parser/trees/WithStatementTree.java` | 0% | 0 | 0.45 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 164 | 58 | `jscomp/parserrunner.go` | `jscomp/parsing/ParserRunner.java` | 14% | 2 | 0.89 | 2 stubs; 1 go-only fns |
| 165 | 58 | `parser/trees/computed_property_field_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java` | 0% | 0 | 1.50 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 166 | 58 | `parser/trees/field_declaration_tree.go` | `jscomp/parsing/parser/trees/FieldDeclarationTree.java` | 0% | 0 | 1.50 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 167 | 58 | `parser/trees/import_meta_expression_tree.go` | `jscomp/parsing/parser/trees/ImportMetaExpressionTree.java` | 0% | 0 | 1.50 | tiny-wrapper(parent bonus 0); 1 go-only fns |
| 168 | 58 | `jscomp/checkregexp.go` | `jscomp/CheckRegExp.java` | 20% | 0 | 1.45 | 5 go-only fns |
| 169 | 58 | `parser/scanner.go` | `jscomp/parsing/parser/Scanner.java` | 51% | 0 | 1.42 | 43 go-only fns |
| 170 | 58 | `rhino/static_scope.go` | `rhino/StaticScope.java` | 20% | 0 | 0.53 |  |
| 171 | 59 | `rhino/jsdoc_info_builder.go` | `rhino/JSDocInfo.java` | 34% | 1 | 0.38 | 1 stubs; 3 go-only fns |
| 172 | 59 | `rhino/jstype/enum_element_type.go` | `rhino/jstype/EnumElementType.java` | 17% | 0 | 0.22 |  |
| 173 | 60 | `rhino/static_source_file.go` | `rhino/StaticSourceFile.java` | 45% | 0 | 0.19 |  |
| 174 | 61 | `parser/trees/comment.go` | `jscomp/parsing/parser/trees/Comment.java` | 20% | 0 | 0.04 | 1 go-only fns |
| 175 | 61 | `jscomp/checkmissingreturn.go` | `jscomp/CheckMissingReturn.java` | 38% | 0 | 1.50 | 6 go-only fns |
| 176 | 61 | `rhino/jstype/record_type.go` | `rhino/jstype/RecordType.java` | 31% | 0 | 1.50 | 5 go-only fns |
| 177 | 62 | `jscomp/parsing/irfactory.go` | `jscomp/parsing/IRFactory.java` | 30% | 0 | 0.34 | 2 go-only fns |
| 178 | 62 | `rhino/qualified_name.go` | `rhino/QualifiedName.java` | 52% | 1 | 1.25 | 1 stubs; 6 go-only fns |
| 179 | 62 | `rhino/jstype/symbol_type.go` | `rhino/jstype/SymbolType.java` | 1% | 0 | 0.18 | tiny-wrapper(parent bonus 20); 2 go-only fns |
| 180 | 62 | `diag/diagnosticgroup.go` | `jscomp/DiagnosticGroup.java` | 25% | 0 | 1.41 | 5 go-only fns |
| 181 | 62 | `rhino/jstype/arrow_type.go` | `rhino/jstype/ArrowType.java` | 24% | 0 | 0.51 | 4 go-only fns |
| 182 | 63 | `parser/parser_statements.go` | `jscomp/parsing/parser/Parser.java` | 23% | 0 | 0.36 |  |
| 183 | 63 | `jscomp/parsing/jsdoctokenstream.go` | `jscomp/parsing/JsDocTokenStream.java` | 26% | 0 | 0.56 | 6 go-only fns |
| 184 | 63 | `jscomp/checksuspiciouscode.go` | `jscomp/CheckSuspiciousCode.java` | 44% | 0 | 1.50 | 3 go-only fns |
| 185 | 63 | `rhino/jstype/known_symbol_type.go` | `rhino/jstype/KnownSymbolType.java` | 12% | 0 | 0.40 | 3 go-only fns |
| 186 | 63 | `jscomp/checkdebuggerstatement.go` | `jscomp/CheckDebuggerStatement.java` | 17% | 0 | 1.14 | 2 go-only fns |
| 187 | 64 | `rhino/jstype/no_type.go` | `rhino/jstype/NoType.java` | 40% | 1 | 0.00 | 1 stubs; 2 go-only fns |
| 188 | 64 | `rhino/jstype/prototype_object_type.go` | `rhino/jstype/PrototypeObjectType.java` | 44% | 2 | 0.11 | 2 stubs; 1 go-only fns |
| 189 | 64 | `diag/allowlistwarningsguard.go` | `jscomp/AllowlistWarningsGuard.java` | 33% | 0 | 0.50 | 4 go-only fns |
| 190 | 64 | `jscomp/checkglobalthis.go` | `jscomp/CheckGlobalThis.java` | 33% | 0 | 1.50 | 2 go-only fns |
| 191 | 65 | `parser/trees/parse_tree_type.go` | `jscomp/parsing/parser/trees/ParseTreeType.java` | 0% | 0 | 0.00 |  |
| 192 | 65 | `rhino/jstype/equivalence_method.go` | `rhino/jstype/EquivalenceMethod.java` | 0% | 0 | 0.00 |  |
| 193 | 65 | `rhino/jstype/jstype_native.go` | `rhino/jstype/JSTypeNative.java` | 0% | 0 | 0.00 |  |
| 194 | 67 | `jscomp/scoped_name.go` | `jscomp/ScopedName.java` | 20% | 0 | 1.50 | 4 go-only fns |
| 195 | 67 | `rhino/jstype/instance_object_type.go` | `rhino/jstype/InstanceObjectType.java` | 59% | 3 | 1.50 | 3 stubs; 2 go-only fns |
| 196 | 67 | `jscomp/ast_analyzer.go` | `jscomp/AstAnalyzer.java` | 40% | 0 | 0.37 | 2 go-only fns |
| 197 | 69 | `jscomp/passfactory.go` | `jscomp/PassFactory.java` | 43% | 0 | 1.50 | 2 go-only fns |
| 198 | 69 | `rhino/jstype/big_int_type.go` | `rhino/jstype/BigIntType.java` | 38% | 1 | 1.50 | 1 stubs; 1 go-only fns |
| 199 | 69 | `rhino/input_id.go` | `rhino/InputId.java` | 29% | 0 | 0.23 | 1 go-only fns |
| 200 | 70 | `rhino/jstype/template_type_replacer.go` | `rhino/jstype/TemplateTypeReplacer.java` | 43% | 0 | 0.89 | 5 go-only fns |
| 201 | 71 | `rhino/jstype/property_map.go` | `rhino/jstype/PropertyMap.java` | 62% | 0 | 1.19 | 7 go-only fns |
| 202 | 72 | `rhino/jstype/subtype_checker.go` | `rhino/jstype/SubtypeChecker.java` | 46% | 0 | 0.48 | 2 go-only fns |
| 203 | 72 | `rhino/jstype/equality_checker.go` | `rhino/jstype/EqualityChecker.java` | 38% | 0 | 0.98 | 4 go-only fns |
| 204 | 72 | `jscomp/nodeiterators.go` | `jscomp/NodeIterators.java` | 44% | 0 | 1.23 | 4 go-only fns |
| 205 | 72 | `rhino/jstype/template_type_map.go` | `rhino/jstype/TemplateTypeMap.java` | 61% | 0 | 0.99 | 7 go-only fns |
| 206 | 72 | `jscomp/astmanipulations.go` | `jscomp/AstManipulations.java` | 50% | 0 | 0.22 |  |
| 207 | 72 | `rhino/simple_source_file.go` | `rhino/SimpleSourceFile.java` | 50% | 2 | 1.50 | 2 stubs; 2 go-only fns |
| 208 | 73 | `rhino/jstype/jstype_iterations.go` | `rhino/jstype/JSTypeIterations.java` | 43% | 0 | 1.11 | 3 go-only fns |
| 209 | 73 | `rhino/jstype/type_string_builder.go` | `rhino/jstype/TypeStringBuilder.java` | 50% | 0 | 0.89 | 3 go-only fns |
| 210 | 74 | `jscomp/passconfig.go` | `jscomp/PassConfig.java` | 46% | 0 | 1.24 | 2 go-only fns |
| 211 | 74 | `rhino/jstype/boolean_literal_set.go` | `rhino/jstype/BooleanLiteralSet.java` | 43% | 0 | 0.53 | 1 go-only fns |
| 212 | 75 | `parser/parser_new.go` | `jscomp/parsing/parser/Parser.java` | 41% | 0 | 0.67 |  |
| 213 | 76 | `jscomp/typetransformation.go` | `jscomp/TypeTransformation.java` | 76% | 0 | 0.85 | 7 go-only fns |
| 214 | 76 | `diag/diagnosticgroups.go` | `jscomp/DiagnosticGroups.java` | 50% | 0 | 0.29 | 1 go-only fns |
| 215 | 77 | `jscomp/typedvar.go` | `jscomp/TypedVar.java` | 57% | 0 | 0.16 | 2 go-only fns |
| 216 | 79 | `jscomp/typemismatch.go` | `jscomp/TypeMismatch.java` | 56% | 0 | 1.45 | 1 go-only fns |
| 217 | 81 | `diag/abstractmessageformatter.go` | `jscomp/AbstractMessageFormatter.java` | 67% | 0 | 0.36 | 1 go-only fns |
| 218 | 81 | `rhino/jstype/contains_upper_bound_super_type_visitor.go` | `rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 60% | 0 | 1.50 | 2 go-only fns |
| 219 | 81 | `rhino/jstype_expr.go` | `rhino/JSTypeExpression.java` | 68% | 0 | 1.43 | 3 go-only fns |
| 220 | 81 | `rhino/jstype/record_type_builder.go` | `rhino/jstype/RecordTypeBuilder.java` | 50% | 0 | 0.78 | 1 go-only fns |
| 221 | 82 | `parser/trees/break_statement_tree.go` | `jscomp/parsing/parser/trees/BreakStatementTree.java` | 1% | 0 | 0.36 | tiny-wrapper(parent bonus 15); 1 go-only fns |
| 222 | 82 | `parser/trees/continue_statement_tree.go` | `jscomp/parsing/parser/trees/ContinueStatementTree.java` | 1% | 0 | 0.36 | tiny-wrapper(parent bonus 15); 1 go-only fns |
| 223 | 82 | `rhino/outcome.go` | `rhino/Outcome.java` | 60% | 0 | 0.40 | 1 go-only fns |
| 224 | 84 | `rhino/non_jsdoc_comment.go` | `rhino/NonJSDocComment.java` | 64% | 0 | 1.50 | 1 go-only fns |
| 225 | 84 | `parser/trees/function_declaration_tree.go` | `jscomp/parsing/parser/trees/FunctionDeclarationTree.java` | 75% | 0 | 0.37 | 1 go-only fns |
| 226 | 86 | `rhino/jstype/can_cast_to_visitor.go` | `rhino/jstype/CanCastToVisitor.java` | 79% | 0 | 1.17 | 2 go-only fns |
| 227 | 87 | `jscomp/basicblock.go` | `jscomp/BasicBlock.java` | 67% | 0 | 1.06 | 1 go-only fns |
| 228 | 88 | `rhino/jstype/jstype_resolver.go` | `rhino/jstype/JSTypeResolver.java` | 80% | 1 | 0.97 | 1 stubs |
| 229 | 88 | `rhino/jstype/function_param_builder.go` | `rhino/jstype/FunctionParamBuilder.java` | 75% | 0 | 0.97 | 2 go-only fns |
| 230 | 88 | `rhino/jstype/visitor.go` | `rhino/jstype/Visitor.java` | 95% | 0 | 0.03 |  |
| 231 | 89 | `rhino/source_position.go` | `rhino/SourcePosition.java` | 89% | 0 | 0.37 |  |
| 232 | 91 | `jscomp/passlistbuilder.go` | `jscomp/PassListBuilder.java` | 76% | 0 | 1.12 | 1 go-only fns |

