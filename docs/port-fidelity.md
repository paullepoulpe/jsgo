# Port Fidelity Report

Generated: auto-computed by tools/portcheck.

## Summary

| Metric | Value |
|--------|-------|
| Total files scored | 244 |
| Average score | 61.7 |
| Std deviation | 11.2 |
| Min score | 0.0 |
| Median score | 59.3 |
| Max score | 90.9 |

## Files by Fidelity Score (worst first)

| Rank | Score | Go File | Java Source | Method Overlap % | Stub Count | LoC Ratio | Port Debt | Notes |
|------|-------|---------|-------------|-----------------|------------|-----------|-----------|-------|
| 1 | 0 | `jscomp/default_externs.go` | `` | 0% | 0 | 0.00 | 0 | missing Port of header |
| 2 | 33 | `jscomp/compiler.go` | `jscomp/Compiler.java` | 0% | 0 | 0.04 | 0 | 4 go-only fns |
| 3 | 33 | `jscomp/checkjsdoctypes.go` | `jscomp/TypeValidator.java` | 0% | 0 | 1.33 | 0 | 23 go-only fns |
| 4 | 33 | `jscomp/typedscope.go` | `jscomp/TypedScope.java` | 0% | 0 | 1.50 | 0 | 37 go-only fns |
| 5 | 34 | `jscomp/semanticreverseabstractinterpreter.go` | `jscomp/type/SemanticReverseAbstractInterpreter.java` | 0% | 0 | 0.73 | 0 | 8 go-only fns |
| 6 | 34 | `jscomp/var.go` | `jscomp/Var.java` | 0% | 0 | 0.77 | 0 | 9 go-only fns |
| 7 | 39 | `jscomp/typecheck.go` | `jscomp/TypeCheck.java` | 10% | 2 | 1.50 | 0 | 2 stubs; 138 go-only fns |
| 8 | 41 | `rhino/type_declarations_ir.go` | `rhino/TypeDeclarationsIR.java` | 0% | 0 | 0.72 | 0 | 17 go-only fns |
| 9 | 42 | `jscomp/nodeutil.go` | `jscomp/NodeUtil.java` | 22% | 1 | 0.45 | 0 | 1 stubs; 31 go-only fns |
| 10 | 42 | `jscomp/nodetraversal.go` | `jscomp/NodeTraversal.java` | 15% | 3 | 0.17 | 0 | 3 stubs; 3 go-only fns |
| 11 | 44 | `jscomp/typedscopecreator.go` | `jscomp/TypedScopeCreator.java` | 2% | 0 | 1.07 | 0 | 56 go-only fns |
| 12 | 44 | `jscomp/scopecreator.go` | `jscomp/ScopeCreator.java` | 0% | 0 | 1.50 | 0 | 4 go-only fns |
| 13 | 45 | `jscomp/parsing/jsdoctoken.go` | `jscomp/parsing/JsDocToken.java` | 0% | 0 | 1.50 | 0 |  |
| 14 | 46 | `jscomp/abstractcompiler.go` | `jscomp/AbstractCompiler.java` | 64% | 50 | 0.52 | 0 | 50 stubs; 3 go-only fns |
| 15 | 50 | `jscomp/scope_untyped.go` | `jscomp/Scope.java` | 19% | 0 | 1.27 | 0 | 13 go-only fns |
| 16 | 50 | `jscomp/dataflow.go` | `jscomp/DataFlowAnalysis.java` | 24% | 1 | 1.46 | 0 | 1 stubs; 7 go-only fns |
| 17 | 50 | `jscomp/diagnosticgroup.go` | `jscomp/DiagnosticGroup.java` | 0% | 0 | 0.00 | 0 |  |
| 18 | 50 | `jscomp/passnames.go` | `jscomp/PassNames.java` | 0% | 0 | 0.00 | 0 |  |
| 19 | 50 | `jscomp/reverseabstractinterpreter.go` | `jscomp/type/ReverseAbstractInterpreter.java` | 0% | 0 | 0.00 | 0 |  |
| 20 | 50 | `rhino/jstype/jstype_class.go` | `rhino/jstype/JSTypeClass.java` | 0% | 0 | 0.00 | 0 |  |
| 21 | 50 | `rhino/jstype/relationship_visitor.go` | `rhino/jstype/RelationshipVisitor.java` | 0% | 0 | 0.00 | 0 |  |
| 22 | 50 | `rhino/static_ref.go` | `rhino/StaticRef.java` | 0% | 0 | 0.00 | 0 |  |
| 23 | 50 | `rhino/static_slot.go` | `rhino/StaticSlot.java` | 0% | 0 | 0.00 | 0 |  |
| 24 | 50 | `rhino/static_symbol_table.go` | `rhino/StaticSymbolTable.java` | 0% | 0 | 0.00 | 0 |  |
| 25 | 50 | `jscomp/controlflowanalysis.go` | `jscomp/ControlFlowAnalysis.java` | 3% | 0 | 0.96 | 0 | 28 go-only fns |
| 26 | 50 | `rhino/token.go` | `rhino/Token.java` | 0% | 0 | 1.50 | 0 | 1 go-only fns |
| 27 | 51 | `jscomp/nodenameextractor.go` | `jscomp/NodeNameExtractor.java` | 14% | 0 | 1.50 | 0 | 5 go-only fns |
| 28 | 52 | `parser/trees/parse_tree.go` | `jscomp/parsing/parser/trees/ParseTree.java` | 3% | 0 | 0.02 | 0 |  |
| 29 | 52 | `rhino/jstype/function_type.go` | `rhino/jstype/FunctionType.java` | 33% | 0 | 0.36 | 0 | 28 go-only fns |
| 30 | 52 | `jscomp/chainablereverseabstractinterpreter.go` | `jscomp/type/ChainableReverseAbstractInterpreter.java` | 65% | 3 | 0.74 | 0 | 3 stubs; 10 go-only fns |
| 31 | 52 | `rhino/jstype/template_type.go` | `rhino/jstype/TemplateType.java` | 9% | 0 | 1.50 | 0 | 8 go-only fns |
| 32 | 53 | `jscomp/compilerinput.go` | `jscomp/CompilerInput.java` | 45% | 2 | 0.34 | 0 | 2 stubs; 15 go-only fns |
| 33 | 53 | `jscomp/checkjsdoc.go` | `jscomp/CheckJSDoc.java` | 32% | 2 | 1.19 | 0 | 2 stubs; 23 go-only fns |
| 34 | 53 | `rhino/jstype/object_type.go` | `rhino/jstype/ObjectType.java` | 42% | 8 | 0.65 | 0 | 8 stubs; 2 go-only fns |
| 35 | 53 | `rhino/jstype/templatized_type.go` | `rhino/jstype/TemplatizedType.java` | 15% | 0 | 0.49 | 0 | 9 go-only fns |
| 36 | 53 | `jscomp/typeinferencepass.go` | `jscomp/TypeInferencePass.java` | 0% | 0 | 1.50 | 0 |  |
| 37 | 53 | `rhino/jstype/object_type_ext.go` | `rhino/jstype/ObjectType.java` | 4% | 0 | 0.07 | 0 |  |
| 38 | 54 | `rhino/jstype/property.go` | `rhino/jstype/Property.java` | 55% | 2 | 1.50 | 0 | 2 stubs; 14 go-only fns |
| 39 | 54 | `jscomp/abstractscope.go` | `jscomp/AbstractScope.java` | 33% | 1 | 1.23 | 0 | 1 stubs; 15 go-only fns |
| 40 | 54 | `jscomp/varcheck.go` | `jscomp/VarCheck.java` | 8% | 0 | 0.49 | 0 | 4 go-only fns |
| 41 | 55 | `rhino/node.go` | `rhino/Node.java` | 51% | 0 | 0.26 | 0 | 17 go-only fns |
| 42 | 55 | `jscomp/typevalidator.go` | `jscomp/TypeValidator.java` | 43% | 2 | 0.72 | 0 | 2 stubs; 22 go-only fns |
| 43 | 55 | `rhino/jstype/jstype.go` | `rhino/jstype/JSType.java` | 15% | 0 | 0.59 | 0 |  |
| 44 | 56 | `jscomp/controlflowgraph.go` | `jscomp/ControlFlowGraph.java` | 0% | 0 | 0.67 | 0 | 3 go-only fns |
| 45 | 56 | `jscomp/inferconsts.go` | `jscomp/InferConsts.java` | 0% | 0 | 1.50 | 0 |  |
| 46 | 56 | `parser/parser.go` | `jscomp/parsing/parser/Parser.java` | 8% | 0 | 1.17 | 0 | 2 go-only fns |
| 47 | 56 | `rhino/jstype/unknown_type.go` | `rhino/jstype/UnknownType.java` | 0% | 0 | 0.00 | 0 | tiny-wrapper(parent bonus 15) |
| 48 | 57 | `rhino/ir.go` | `rhino/IR.java` | 0% | 0 | 0.30 | 0 |  |
| 49 | 57 | `rhino/jstype/no_resolved_type.go` | `rhino/jstype/NoResolvedType.java` | 6% | 0 | 0.08 | 0 |  |
| 50 | 57 | `rhino/jstype/all_type.go` | `rhino/jstype/AllType.java` | 1% | 0 | 0.00 | 0 | tiny-wrapper(parent bonus 15) |
| 51 | 57 | `parser/trees/formal_parameter_list_tree.go` | `jscomp/parsing/parser/trees/FormalParameterListTree.java` | 0% | 0 | 0.22 | 0 | tiny-wrapper(parent bonus 0) |
| 52 | 58 | `jscomp/checktypeimportcodereferences.go` | `jscomp/CheckTypeImportCodeReferences.java` | 11% | 0 | 1.50 | 4 | 4 port-debt |
| 53 | 58 | `parser/trees/comprehension_tree.go` | `jscomp/parsing/parser/trees/ComprehensionTree.java` | 0% | 0 | 0.27 | 0 | tiny-wrapper(parent bonus 0) |
| 54 | 58 | `parser/trees/update_expression_tree.go` | `jscomp/parsing/parser/trees/UpdateExpressionTree.java` | 0% | 0 | 0.30 | 0 | tiny-wrapper(parent bonus 0) |
| 55 | 58 | `jscomp/checkregexp.go` | `jscomp/CheckRegExp.java` | 20% | 0 | 1.45 | 0 | 5 go-only fns |
| 56 | 58 | `rhino/jstype/void_type.go` | `rhino/jstype/VoidType.java` | 1% | 0 | 0.06 | 0 | tiny-wrapper(parent bonus 15) |
| 57 | 58 | `parser/trees/argument_list_tree.go` | `jscomp/parsing/parser/trees/ArgumentListTree.java` | 0% | 0 | 0.32 | 0 | tiny-wrapper(parent bonus 0) |
| 58 | 58 | `parser/trees/class_declaration_tree.go` | `jscomp/parsing/parser/trees/ClassDeclarationTree.java` | 0% | 0 | 0.32 | 0 | tiny-wrapper(parent bonus 0) |
| 59 | 58 | `parser/trees/import_declaration_tree.go` | `jscomp/parsing/parser/trees/ImportDeclarationTree.java` | 0% | 0 | 0.32 | 0 | tiny-wrapper(parent bonus 0) |
| 60 | 58 | `parser/scanner.go` | `jscomp/parsing/parser/Scanner.java` | 51% | 0 | 1.42 | 0 | 43 go-only fns |
| 61 | 58 | `rhino/jstype/null_type.go` | `rhino/jstype/NullType.java` | 1% | 0 | 0.06 | 0 | tiny-wrapper(parent bonus 15) |
| 62 | 58 | `parser/trees/export_declaration_tree.go` | `jscomp/parsing/parser/trees/ExportDeclarationTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 63 | 58 | `parser/trees/get_accessor_tree.go` | `jscomp/parsing/parser/trees/GetAccessorTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 64 | 58 | `parser/trees/object_literal_expression_tree.go` | `jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 65 | 58 | `parser/trees/optional_member_expression_tree.go` | `jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 66 | 58 | `parser/trees/program_tree.go` | `jscomp/parsing/parser/trees/ProgramTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 67 | 58 | `parser/trees/set_accessor_tree.go` | `jscomp/parsing/parser/trees/SetAccessorTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 68 | 58 | `parser/trees/variable_declaration_list_tree.go` | `jscomp/parsing/parser/trees/VariableDeclarationListTree.java` | 0% | 0 | 0.33 | 0 | tiny-wrapper(parent bonus 0) |
| 69 | 58 | `rhino/static_scope.go` | `rhino/StaticScope.java` | 20% | 0 | 0.53 | 0 |  |
| 70 | 58 | `rhino/jstype/registry.go` | `rhino/jstype/JSTypeRegistry.java` | 14% | 0 | 0.22 | 1 | 1 go-only fns; 1 port-debt |
| 71 | 58 | `parser/trees/computed_property_setter_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertySetterTree.java` | 0% | 0 | 0.35 | 0 | tiny-wrapper(parent bonus 0) |
| 72 | 58 | `parser/trees/opt_chain_call_expression_tree.go` | `jscomp/parsing/parser/trees/OptChainCallExpressionTree.java` | 0% | 0 | 0.35 | 0 | tiny-wrapper(parent bonus 0) |
| 73 | 59 | `parser/trees/new_expression_tree.go` | `jscomp/parsing/parser/trees/NewExpressionTree.java` | 0% | 0 | 0.35 | 0 | tiny-wrapper(parent bonus 0) |
| 74 | 59 | `parser/trees/optional_member_lookup_expression_tree.go` | `jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java` | 0% | 0 | 0.35 | 0 | tiny-wrapper(parent bonus 0) |
| 75 | 59 | `jscomp/symboltable.go` | `jscomp/SymbolTable.java` | 50% | 1 | 0.88 | 0 | 1 stubs; 56 go-only fns |
| 76 | 59 | `jscomp/checkaccesscontrols.go` | `jscomp/CheckAccessControls.java` | 0% | 0 | 0.83 | 0 |  |
| 77 | 59 | `rhino/jstype/number_type.go` | `rhino/jstype/NumberType.java` | 1% | 0 | 0.07 | 0 | tiny-wrapper(parent bonus 15) |
| 78 | 59 | `rhino/jstype/string_type.go` | `rhino/jstype/StringType.java` | 1% | 0 | 0.07 | 0 | tiny-wrapper(parent bonus 15) |
| 79 | 59 | `parser/trees/case_clause_tree.go` | `jscomp/parsing/parser/trees/CaseClauseTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 80 | 59 | `parser/trees/comprehension_for_tree.go` | `jscomp/parsing/parser/trees/ComprehensionForTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 81 | 59 | `parser/trees/template_literal_expression_tree.go` | `jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 82 | 59 | `rhino/jstype/boolean_type.go` | `rhino/jstype/BooleanType.java` | 1% | 0 | 0.07 | 0 | tiny-wrapper(parent bonus 15) |
| 83 | 59 | `jscomp/astfactory.go` | `jscomp/AstFactory.java` | 39% | 0 | 0.24 | 0 | 6 go-only fns |
| 84 | 59 | `parser/trees/array_pattern_tree.go` | `jscomp/parsing/parser/trees/ArrayPatternTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 85 | 59 | `parser/trees/comma_expression_tree.go` | `jscomp/parsing/parser/trees/CommaExpressionTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 86 | 59 | `parser/trees/comprehension_if_tree.go` | `jscomp/parsing/parser/trees/ComprehensionIfTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 87 | 59 | `parser/trees/object_pattern_tree.go` | `jscomp/parsing/parser/trees/ObjectPatternTree.java` | 0% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 0) |
| 88 | 59 | `jscomp/parserrunner.go` | `jscomp/parsing/ParserRunner.java` | 14% | 2 | 0.97 | 0 | 2 stubs; 1 go-only fns |
| 89 | 59 | `parser/trees/array_literal_expression_tree.go` | `jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java` | 0% | 0 | 0.38 | 0 | tiny-wrapper(parent bonus 0) |
| 90 | 59 | `parser/trees/export_specifier_tree.go` | `jscomp/parsing/parser/trees/ExportSpecifierTree.java` | 0% | 0 | 0.38 | 0 | tiny-wrapper(parent bonus 0) |
| 91 | 59 | `parser/trees/import_specifier_tree.go` | `jscomp/parsing/parser/trees/ImportSpecifierTree.java` | 0% | 0 | 0.38 | 0 | tiny-wrapper(parent bonus 0) |
| 92 | 59 | `parser/trees/switch_statement_tree.go` | `jscomp/parsing/parser/trees/SwitchStatementTree.java` | 0% | 0 | 0.38 | 0 | tiny-wrapper(parent bonus 0) |
| 93 | 59 | `parser/trees/block_tree.go` | `jscomp/parsing/parser/trees/BlockTree.java` | 0% | 0 | 0.40 | 0 | tiny-wrapper(parent bonus 0) |
| 94 | 59 | `parser/trees/default_clause_tree.go` | `jscomp/parsing/parser/trees/DefaultClauseTree.java` | 0% | 0 | 0.40 | 0 | tiny-wrapper(parent bonus 0) |
| 95 | 59 | `parser/trees/identifier_expression_tree.go` | `jscomp/parsing/parser/trees/IdentifierExpressionTree.java` | 0% | 0 | 0.40 | 0 | tiny-wrapper(parent bonus 0) |
| 96 | 59 | `parser/trees/literal_expression_tree.go` | `jscomp/parsing/parser/trees/LiteralExpressionTree.java` | 0% | 0 | 0.40 | 0 | tiny-wrapper(parent bonus 0) |
| 97 | 59 | `parser/trees/template_literal_portion_tree.go` | `jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java` | 0% | 0 | 0.40 | 0 | tiny-wrapper(parent bonus 0) |
| 98 | 59 | `parser/trees/computed_property_definition_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 99 | 59 | `parser/trees/default_parameter_tree.go` | `jscomp/parsing/parser/trees/DefaultParameterTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 100 | 59 | `parser/trees/do_while_statement_tree.go` | `jscomp/parsing/parser/trees/DoWhileStatementTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 101 | 59 | `parser/trees/labelled_statement_tree.go` | `jscomp/parsing/parser/trees/LabelledStatementTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 102 | 59 | `parser/trees/member_expression_tree.go` | `jscomp/parsing/parser/trees/MemberExpressionTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 103 | 59 | `parser/trees/member_lookup_expression_tree.go` | `jscomp/parsing/parser/trees/MemberLookupExpressionTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 104 | 59 | `parser/trees/property_name_assignment_tree.go` | `jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 105 | 59 | `parser/trees/unary_expression_tree.go` | `jscomp/parsing/parser/trees/UnaryExpressionTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 106 | 59 | `parser/trees/yield_expression_tree.go` | `jscomp/parsing/parser/trees/YieldExpressionTree.java` | 0% | 0 | 0.42 | 0 | tiny-wrapper(parent bonus 0) |
| 107 | 59 | `parser/tokentype.go` | `jscomp/parsing/parser/TokenType.java` | 5% | 0 | 1.37 | 0 |  |
| 108 | 59 | `parser/trees/binary_operator_tree.go` | `jscomp/parsing/parser/trees/BinaryOperatorTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 109 | 59 | `parser/trees/computed_property_getter_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 110 | 59 | `parser/trees/conditional_expression_tree.go` | `jscomp/parsing/parser/trees/ConditionalExpressionTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 111 | 59 | `parser/trees/debugger_statement_tree.go` | `jscomp/parsing/parser/trees/DebuggerStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 112 | 59 | `parser/trees/empty_statement_tree.go` | `jscomp/parsing/parser/trees/EmptyStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 113 | 59 | `parser/trees/for_await_of_statement_tree.go` | `jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 114 | 59 | `parser/trees/for_in_statement_tree.go` | `jscomp/parsing/parser/trees/ForInStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 115 | 59 | `parser/trees/for_of_statement_tree.go` | `jscomp/parsing/parser/trees/ForOfStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 116 | 59 | `parser/trees/if_statement_tree.go` | `jscomp/parsing/parser/trees/IfStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 117 | 59 | `parser/trees/missing_primary_expression_tree.go` | `jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 118 | 59 | `parser/trees/new_target_expression_tree.go` | `jscomp/parsing/parser/trees/NewTargetExpressionTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 119 | 59 | `parser/trees/null_tree.go` | `jscomp/parsing/parser/trees/NullTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 120 | 59 | `parser/trees/super_expression_tree.go` | `jscomp/parsing/parser/trees/SuperExpressionTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 121 | 59 | `parser/trees/this_expression_tree.go` | `jscomp/parsing/parser/trees/ThisExpressionTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 122 | 59 | `parser/trees/try_statement_tree.go` | `jscomp/parsing/parser/trees/TryStatementTree.java` | 0% | 0 | 0.43 | 0 | tiny-wrapper(parent bonus 0) |
| 123 | 59 | `parser/trees/for_statement_tree.go` | `jscomp/parsing/parser/trees/ForStatementTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 124 | 59 | `parser/trees/await_expression_tree.go` | `jscomp/parsing/parser/trees/AwaitExpressionTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 125 | 59 | `parser/trees/dynamic_import_tree.go` | `jscomp/parsing/parser/trees/DynamicImportTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 126 | 59 | `parser/trees/expression_statement_tree.go` | `jscomp/parsing/parser/trees/ExpressionStatementTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 127 | 59 | `parser/trees/finally_tree.go` | `jscomp/parsing/parser/trees/FinallyTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 128 | 59 | `parser/trees/iter_rest_tree.go` | `jscomp/parsing/parser/trees/IterRestTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 129 | 59 | `parser/trees/iter_spread_tree.go` | `jscomp/parsing/parser/trees/IterSpreadTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 130 | 59 | `parser/trees/object_rest_tree.go` | `jscomp/parsing/parser/trees/ObjectRestTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 131 | 59 | `parser/trees/object_spread_tree.go` | `jscomp/parsing/parser/trees/ObjectSpreadTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 132 | 59 | `parser/trees/paren_expression_tree.go` | `jscomp/parsing/parser/trees/ParenExpressionTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 133 | 59 | `parser/trees/return_statement_tree.go` | `jscomp/parsing/parser/trees/ReturnStatementTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 134 | 59 | `parser/trees/template_substitution_tree.go` | `jscomp/parsing/parser/trees/TemplateSubstitutionTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 135 | 59 | `parser/trees/throw_statement_tree.go` | `jscomp/parsing/parser/trees/ThrowStatementTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 136 | 59 | `parser/trees/variable_statement_tree.go` | `jscomp/parsing/parser/trees/VariableStatementTree.java` | 0% | 0 | 0.44 | 0 | tiny-wrapper(parent bonus 0) |
| 137 | 60 | `parser/trees/call_expression_tree.go` | `jscomp/parsing/parser/trees/CallExpressionTree.java` | 0% | 0 | 0.45 | 0 | tiny-wrapper(parent bonus 0) |
| 138 | 60 | `parser/trees/catch_tree.go` | `jscomp/parsing/parser/trees/CatchTree.java` | 0% | 0 | 0.45 | 0 | tiny-wrapper(parent bonus 0) |
| 139 | 60 | `parser/trees/computed_property_method_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java` | 0% | 0 | 0.45 | 0 | tiny-wrapper(parent bonus 0) |
| 140 | 60 | `parser/trees/variable_declaration_tree.go` | `jscomp/parsing/parser/trees/VariableDeclarationTree.java` | 0% | 0 | 0.45 | 0 | tiny-wrapper(parent bonus 0) |
| 141 | 60 | `parser/trees/while_statement_tree.go` | `jscomp/parsing/parser/trees/WhileStatementTree.java` | 0% | 0 | 0.45 | 0 | tiny-wrapper(parent bonus 0) |
| 142 | 60 | `parser/trees/with_statement_tree.go` | `jscomp/parsing/parser/trees/WithStatementTree.java` | 0% | 0 | 0.45 | 0 | tiny-wrapper(parent bonus 0) |
| 143 | 60 | `jscomp/varrefcheck.go` | `jscomp/VariableReferenceCheck.java` | 3% | 0 | 1.50 | 1 | 2 go-only fns; 1 port-debt |
| 144 | 60 | `jscomp/defaultpassconfig.go` | `jscomp/DefaultPassConfig.java` | 17% | 0 | 0.61 | 0 |  |
| 145 | 60 | `parser/trees/computed_property_field_tree.go` | `jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java` | 0% | 0 | 1.50 | 0 | tiny-wrapper(parent bonus 0) |
| 146 | 60 | `parser/trees/field_declaration_tree.go` | `jscomp/parsing/parser/trees/FieldDeclarationTree.java` | 0% | 0 | 1.50 | 0 | tiny-wrapper(parent bonus 0) |
| 147 | 60 | `parser/trees/import_meta_expression_tree.go` | `jscomp/parsing/parser/trees/ImportMetaExpressionTree.java` | 0% | 0 | 1.50 | 0 | tiny-wrapper(parent bonus 0) |
| 148 | 60 | `rhino/static_source_file.go` | `rhino/StaticSourceFile.java` | 45% | 0 | 0.19 | 0 |  |
| 149 | 60 | `rhino/nominal_type_builder.go` | `rhino/NominalTypeBuilder.java` | 8% | 0 | 0.67 | 0 | 4 go-only fns |
| 150 | 60 | `rhino/jstype/no_object_type.go` | `rhino/jstype/NoObjectType.java` | 25% | 1 | 0.00 | 0 | 1 stubs |
| 151 | 61 | `jscomp/syntacticscopecreator.go` | `jscomp/SyntacticScopeCreator.java` | 26% | 2 | 1.41 | 0 | 2 stubs; 3 go-only fns |
| 152 | 61 | `rhino/jstype/enum_element_type.go` | `rhino/jstype/EnumElementType.java` | 16% | 0 | 0.30 | 0 |  |
| 153 | 61 | `jscomp/checkmissingreturn.go` | `jscomp/CheckMissingReturn.java` | 38% | 0 | 1.50 | 0 | 6 go-only fns |
| 154 | 61 | `rhino/jstype/record_type.go` | `rhino/jstype/RecordType.java` | 31% | 0 | 1.50 | 0 | 5 go-only fns |
| 155 | 61 | `jscomp/astvalidator.go` | `jscomp/AstValidator.java` | 31% | 0 | 0.24 | 0 | 2 go-only fns |
| 156 | 62 | `rhino/qualified_name.go` | `rhino/QualifiedName.java` | 52% | 1 | 1.25 | 0 | 1 stubs; 6 go-only fns |
| 157 | 62 | `rhino/jsdoc_info_builder.go` | `rhino/JSDocInfo.java` | 34% | 0 | 0.39 | 0 | 3 go-only fns |
| 158 | 62 | `rhino/jstype/symbol_type.go` | `rhino/jstype/SymbolType.java` | 1% | 0 | 0.18 | 0 | tiny-wrapper(parent bonus 20); 2 go-only fns |
| 159 | 62 | `jscomp/compileroptions.go` | `jscomp/CompilerOptions.java` | 21% | 0 | 0.33 | 0 |  |
| 160 | 62 | `diag/diagnosticgroup.go` | `jscomp/DiagnosticGroup.java` | 25% | 0 | 1.41 | 0 | 5 go-only fns |
| 161 | 62 | `rhino/jstype/arrow_type.go` | `rhino/jstype/ArrowType.java` | 24% | 0 | 0.51 | 0 | 4 go-only fns |
| 162 | 63 | `jscomp/dotformatter.go` | `jscomp/DotFormatter.java` | 24% | 0 | 1.50 | 0 | 6 go-only fns |
| 163 | 63 | `rhino/jstype/value_type.go` | `rhino/jstype/ValueType.java` | 8% | 0 | 0.22 | 0 |  |
| 164 | 63 | `parser/parser_statements.go` | `jscomp/parsing/parser/Parser.java` | 23% | 0 | 0.36 | 0 |  |
| 165 | 63 | `jscomp/parsing/jsdoctokenstream.go` | `jscomp/parsing/JsDocTokenStream.java` | 26% | 0 | 0.56 | 0 | 6 go-only fns |
| 166 | 63 | `jscomp/checksuspiciouscode.go` | `jscomp/CheckSuspiciousCode.java` | 44% | 0 | 1.50 | 0 | 3 go-only fns |
| 167 | 63 | `rhino/jstype/known_symbol_type.go` | `rhino/jstype/KnownSymbolType.java` | 12% | 0 | 0.40 | 0 | 3 go-only fns |
| 168 | 63 | `jscomp/checkdebuggerstatement.go` | `jscomp/CheckDebuggerStatement.java` | 17% | 0 | 1.14 | 0 | 2 go-only fns |
| 169 | 64 | `diag/diagnosticgroups.go` | `jscomp/DiagnosticGroups.java` | 33% | 0 | 0.53 | 0 | 5 go-only fns |
| 170 | 64 | `rhino/jstype/named_type.go` | `rhino/jstype/NamedType.java` | 28% | 2 | 0.51 | 0 | 2 stubs |
| 171 | 64 | `rhino/jstype/no_type.go` | `rhino/jstype/NoType.java` | 40% | 1 | 0.00 | 0 | 1 stubs; 2 go-only fns |
| 172 | 64 | `rhino/jstype/prototype_object_type.go` | `rhino/jstype/PrototypeObjectType.java` | 44% | 2 | 0.11 | 0 | 2 stubs; 1 go-only fns |
| 173 | 64 | `jscomp/checkglobalthis.go` | `jscomp/CheckGlobalThis.java` | 33% | 0 | 1.50 | 0 | 2 go-only fns |
| 174 | 64 | `jscomp/inferjsdocinfo.go` | `jscomp/InferJsDocInfo.java` | 8% | 0 | 1.50 | 0 |  |
| 175 | 65 | `rhino/jstype/enum_type.go` | `rhino/jstype/EnumType.java` | 28% | 1 | 0.28 | 0 | 1 stubs |
| 176 | 65 | `jscomp/closurereverseabstractinterpreter.go` | `jscomp/type/ClosureReverseAbstractInterpreter.java` | 22% | 0 | 1.50 | 0 |  |
| 177 | 65 | `parser/trees/parse_tree_type.go` | `jscomp/parsing/parser/trees/ParseTreeType.java` | 0% | 0 | 0.00 | 0 |  |
| 178 | 65 | `rhino/jstype/equivalence_method.go` | `rhino/jstype/EquivalenceMethod.java` | 0% | 0 | 0.00 | 0 |  |
| 179 | 65 | `rhino/jstype/jstype_native.go` | `rhino/jstype/JSTypeNative.java` | 0% | 0 | 0.00 | 0 |  |
| 180 | 65 | `rhino/jstype/union_type.go` | `rhino/jstype/UnionType.java` | 22% | 0 | 0.43 | 0 |  |
| 181 | 65 | `jscomp/checksideeffects.go` | `jscomp/CheckSideEffects.java` | 33% | 0 | 1.50 | 0 |  |
| 182 | 66 | `rhino/token_stream.go` | `rhino/TokenStream.java` | 33% | 0 | 1.50 | 0 | 3 go-only fns |
| 183 | 67 | `jscomp/parsing/irfactory.go` | `jscomp/parsing/IRFactory.java` | 29% | 0 | 0.42 | 0 |  |
| 184 | 67 | `jscomp/scoped_name.go` | `jscomp/ScopedName.java` | 20% | 0 | 1.50 | 0 | 4 go-only fns |
| 185 | 67 | `rhino/jstype/instance_object_type.go` | `rhino/jstype/InstanceObjectType.java` | 59% | 3 | 1.50 | 0 | 3 stubs; 2 go-only fns |
| 186 | 67 | `jscomp/ast_analyzer.go` | `jscomp/AstAnalyzer.java` | 40% | 0 | 0.37 | 0 | 2 go-only fns |
| 187 | 68 | `rhino/closure_primitive.go` | `rhino/ClosurePrimitive.java` | 20% | 0 | 0.52 | 0 |  |
| 188 | 68 | `jscomp/base/tri.go` | `jscomp/base/Tri.java` | 33% | 0 | 0.51 | 0 |  |
| 189 | 69 | `jscomp/checkduplicatemember.go` | `jscomp/StrictModeCheck.java` | 19% | 0 | 1.34 | 1 | 1 port-debt |
| 190 | 69 | `jscomp/linkedflowscope.go` | `jscomp/LinkedFlowScope.java` | 21% | 0 | 0.63 | 0 |  |
| 191 | 69 | `jscomp/checkunreachablecode.go` | `jscomp/CheckUnreachableCode.java` | 36% | 0 | 1.50 | 0 |  |
| 192 | 69 | `jscomp/passfactory.go` | `jscomp/PassFactory.java` | 43% | 0 | 1.50 | 0 | 2 go-only fns |
| 193 | 69 | `rhino/jstype/big_int_type.go` | `rhino/jstype/BigIntType.java` | 38% | 1 | 1.50 | 0 | 1 stubs; 1 go-only fns |
| 194 | 69 | `rhino/input_id.go` | `rhino/InputId.java` | 29% | 0 | 0.23 | 0 | 1 go-only fns |
| 195 | 70 | `rhino/jstype/template_type_replacer.go` | `rhino/jstype/TemplateTypeReplacer.java` | 43% | 0 | 0.89 | 0 | 5 go-only fns |
| 196 | 70 | `jscomp/parsing/jsdocinfoparser.go` | `jscomp/parsing/JsDocInfoParser.java` | 28% | 0 | 0.60 | 0 |  |
| 197 | 70 | `diag/diagnostictype.go` | `jscomp/DiagnosticType.java` | 16% | 0 | 1.50 | 0 |  |
| 198 | 71 | `jscomp/checkclosureimports.go` | `jscomp/CheckClosureImports.java` | 32% | 0 | 1.08 | 17 | 17 port-debt |
| 199 | 71 | `jscomp/checkmissingrequires.go` | `jscomp/CheckMissingRequires.java` | 39% | 0 | 1.33 | 8 | 8 port-debt |
| 200 | 71 | `diag/allowlistwarningsguard.go` | `jscomp/AllowlistWarningsGuard.java` | 58% | 0 | 1.10 | 0 | 7 go-only fns |
| 201 | 71 | `rhino/jstype/property_map.go` | `rhino/jstype/PropertyMap.java` | 62% | 0 | 1.19 | 0 | 7 go-only fns |
| 202 | 72 | `rhino/jsdoc.go` | `rhino/JSDocInfo.java` | 31% | 0 | 0.35 | 0 |  |
| 203 | 72 | `rhino/jstype/subtype_checker.go` | `rhino/jstype/SubtypeChecker.java` | 46% | 0 | 0.48 | 0 | 2 go-only fns |
| 204 | 72 | `rhino/jstype/equality_checker.go` | `rhino/jstype/EqualityChecker.java` | 38% | 0 | 0.98 | 0 | 4 go-only fns |
| 205 | 72 | `jscomp/nodeiterators.go` | `jscomp/NodeIterators.java` | 44% | 0 | 1.23 | 0 | 4 go-only fns |
| 206 | 72 | `rhino/jstype/template_type_map.go` | `rhino/jstype/TemplateTypeMap.java` | 61% | 0 | 0.99 | 0 | 7 go-only fns |
| 207 | 72 | `jscomp/astmanipulations.go` | `jscomp/AstManipulations.java` | 50% | 0 | 0.22 | 0 |  |
| 208 | 72 | `rhino/error_reporter.go` | `rhino/ErrorReporter.java` | 50% | 2 | 0.09 | 0 | 2 stubs |
| 209 | 72 | `rhino/simple_source_file.go` | `rhino/SimpleSourceFile.java` | 50% | 2 | 1.50 | 0 | 2 stubs; 2 go-only fns |
| 210 | 73 | `rhino/jstype/jstype_iterations.go` | `rhino/jstype/JSTypeIterations.java` | 43% | 0 | 1.11 | 0 | 3 go-only fns |
| 211 | 73 | `rhino/jstype/type_string_builder.go` | `rhino/jstype/TypeStringBuilder.java` | 50% | 0 | 0.89 | 0 | 3 go-only fns |
| 212 | 73 | `diag/jserror.go` | `jscomp/JSError.java` | 34% | 1 | 1.01 | 0 | 1 stubs |
| 213 | 73 | `jscomp/abstractvar.go` | `jscomp/AbstractVar.java` | 57% | 0 | 0.70 | 0 | 5 go-only fns |
| 214 | 74 | `jscomp/passconfig.go` | `jscomp/PassConfig.java` | 46% | 0 | 1.24 | 0 | 2 go-only fns |
| 215 | 74 | `jscomp/checksuper.go` | `jscomp/CheckSuper.java` | 50% | 0 | 1.37 | 2 | 2 port-debt |
| 216 | 74 | `rhino/jstype/boolean_literal_set.go` | `rhino/jstype/BooleanLiteralSet.java` | 43% | 0 | 0.53 | 0 | 1 go-only fns |
| 217 | 75 | `jscomp/inference.go` | `jscomp/TypeInference.java` | 47% | 0 | 1.50 | 2 | 2 port-debt |
| 218 | 75 | `parser/parser_new.go` | `jscomp/parsing/parser/Parser.java` | 41% | 0 | 0.68 | 0 |  |
| 219 | 76 | `jscomp/typetransformation.go` | `jscomp/TypeTransformation.java` | 76% | 0 | 0.85 | 0 | 7 go-only fns |
| 220 | 77 | `jscomp/typedvar.go` | `jscomp/TypedVar.java` | 57% | 0 | 0.16 | 0 | 2 go-only fns |
| 221 | 77 | `jscomp/checkconformance.go` | `jscomp/CheckConformance.java` | 67% | 0 | 0.24 | 0 |  |
| 222 | 77 | `jscomp/cfg.go` | `jscomp/ControlFlowGraph.java` | 44% | 1 | 0.95 | 0 | 1 stubs; 1 go-only fns |
| 223 | 79 | `jscomp/typemismatch.go` | `jscomp/TypeMismatch.java` | 56% | 0 | 1.45 | 0 | 1 go-only fns |
| 224 | 80 | `jscomp/abstractmessageformatter.go` | `jscomp/AbstractMessageFormatter.java` | 55% | 0 | 0.45 | 0 | 1 go-only fns |
| 225 | 80 | `rhino/jstype/proxy_object_type.go` | `rhino/jstype/ProxyObjectType.java` | 62% | 1 | 0.62 | 0 | 1 stubs |
| 226 | 81 | `diag/abstractmessageformatter.go` | `jscomp/AbstractMessageFormatter.java` | 67% | 0 | 0.36 | 0 | 1 go-only fns |
| 227 | 81 | `rhino/jstype/contains_upper_bound_super_type_visitor.go` | `rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java` | 60% | 0 | 1.50 | 0 | 2 go-only fns |
| 228 | 81 | `parser/trees/comment.go` | `jscomp/parsing/parser/trees/Comment.java` | 60% | 0 | 0.12 | 0 |  |
| 229 | 81 | `rhino/jstype_expr.go` | `rhino/JSTypeExpression.java` | 68% | 0 | 1.43 | 0 | 3 go-only fns |
| 230 | 81 | `rhino/jstype/record_type_builder.go` | `rhino/jstype/RecordTypeBuilder.java` | 50% | 0 | 0.78 | 0 | 1 go-only fns |
| 231 | 82 | `rhino/outcome.go` | `rhino/Outcome.java` | 60% | 0 | 0.40 | 0 | 1 go-only fns |
| 232 | 84 | `parser/trees/break_statement_tree.go` | `jscomp/parsing/parser/trees/BreakStatementTree.java` | 1% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 15) |
| 233 | 84 | `parser/trees/continue_statement_tree.go` | `jscomp/parsing/parser/trees/ContinueStatementTree.java` | 1% | 0 | 0.36 | 0 | tiny-wrapper(parent bonus 15) |
| 234 | 84 | `rhino/non_jsdoc_comment.go` | `rhino/NonJSDocComment.java` | 64% | 0 | 1.50 | 0 | 1 go-only fns |
| 235 | 86 | `rhino/jstype/can_cast_to_visitor.go` | `rhino/jstype/CanCastToVisitor.java` | 79% | 0 | 1.17 | 0 | 2 go-only fns |
| 236 | 86 | `parser/trees/function_declaration_tree.go` | `jscomp/parsing/parser/trees/FunctionDeclarationTree.java` | 75% | 0 | 0.37 | 0 |  |
| 237 | 87 | `jscomp/basicblock.go` | `jscomp/BasicBlock.java` | 67% | 0 | 1.06 | 0 | 1 go-only fns |
| 238 | 87 | `rhino/token_util.go` | `rhino/TokenUtil.java` | 80% | 0 | 1.47 | 0 |  |
| 239 | 88 | `rhino/jstype/jstype_resolver.go` | `rhino/jstype/JSTypeResolver.java` | 80% | 1 | 0.97 | 0 | 1 stubs |
| 240 | 88 | `rhino/jstype/function_param_builder.go` | `rhino/jstype/FunctionParamBuilder.java` | 75% | 0 | 0.97 | 0 | 2 go-only fns |
| 241 | 88 | `rhino/jstype/visitor.go` | `rhino/jstype/Visitor.java` | 95% | 0 | 0.03 | 0 |  |
| 242 | 89 | `rhino/source_position.go` | `rhino/SourcePosition.java` | 89% | 0 | 0.37 | 0 |  |
| 243 | 90 | `jscomp/flowscope.go` | `jscomp/type/FlowScope.java` | 80% | 0 | 1.00 | 0 |  |
| 244 | 91 | `jscomp/passlistbuilder.go` | `jscomp/PassListBuilder.java` | 76% | 0 | 1.12 | 0 | 1 go-only fns |

