// Port of com/google/javascript/jscomp/serialization/TypedAstSerializer.java @ closure-compiler v20261015

/*
 * Copyright 2021 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypedAstSerializer transforms a compiler AST into a serialized TypedAst
// object.
//
// Java: final class TypedAstSerializer
type TypedAstSerializer struct {
	compiler          AbstractCompiler
	serializationMode *SerializationOptions
	stringPool        *StringPool_Builder
	previousLine      int
	previousColumn    int
	// Go-only: language-adapter — the ArrayDeque<SourceFile> is used only as
	// a stack (addLast, removeLast, peekLast, clear): a slice, top last.
	subtreeSourceFiles []*SourceFile
	sourceFilePointers *jutil.LinkedHashMap[*SourceFile, int]

	typeSerializer TypedAstSerializer_TypeSerializer
}

// Java: TypedAstSerializer(AbstractCompiler, SerializationOptions)
func NewTypedAstSerializer(compiler AbstractCompiler, serializationMode *SerializationOptions) *TypedAstSerializer {
	s := &TypedAstSerializer{
		stringPool:         StringPool_Builder_(),
		sourceFilePointers: jutil.NewLinkedHashMap[*SourceFile, int](),
		typeSerializer:     nil,
	}
	s.compiler = compiler
	s.serializationMode = serializationMode
	return s
}

// SerializeRoots transforms the given compiler AST root nodes into into a
// serialized TypedAst object.
//
// Java: TypedAstSerializer.serializeRoots(Node, Node)
func (s *TypedAstSerializer) SerializeRoots(externsRoot *rhino.Node, jsRoot *rhino.Node) *TypedAst {
	jutil.CheckArgument(externsRoot.IsRoot())
	jutil.CheckArgument(jsRoot.IsRoot())

	if s.compiler.HasOptimizationColors() {
		s.typeSerializer = s.CreateColorTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
	} else if s.compiler.HasTypeCheckingRun() {
		// The AST has JSTypes, but we want to serialize colors instead.
		// TODO(bradfordcsmith): Change this branch to throw an error and delete the related logic.
		//     Nothing should be using it anymore.
		s.typeSerializer = TypedAstSerializer_CreateJSTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
	} else {
		s.typeSerializer = NewTypedAstSerializer_NoOpTypeSerializer()
	}

	builder := TypedAst_NewBuilder()
	for script := externsRoot.GetFirstChild(); script != nil; script = script.GetNext() {
		if NodeUtil_IsFromTypeSummary(script) {
			continue
		}
		builder.AddExternAst(s.SerializeScriptNode(script))
	}
	for script := jsRoot.GetFirstChild(); script != nil; script = script.GetNext() {
		if NodeUtil_IsFromTypeSummary(script) {
			continue
		}
		builder.AddCodeAst(s.SerializeScriptNode(script))
	}

	sourceFileProtos := make([]*SourceFileProto, 0, s.sourceFilePointers.Size())
	for f := range s.sourceFilePointers.Keys() {
		sourceFileProtos = append(sourceFileProtos, f.GetProto())
	}
	sourceFiles := SourceFilePool_NewBuilder().
		AddAllSourceFile(sourceFileProtos).
		Build()

	if s.compiler.GetExternProperties() != nil {
		externsSummary := ExternsSummary_NewBuilder()
		for prop := range s.compiler.GetExternProperties().All() {
			externsSummary.AddPropNamePtr(int32(s.stringPool.Put(prop)))
		}
		builder.SetExternsSummaryFromBuilder(externsSummary)
	}

	return builder.
		AddAllRuntimeLibraryToInject(s.serializationMode.RuntimeLibraries()).
		SetTypePool(s.typeSerializer.GenerateTypePool()).
		SetStringPool(s.stringPool.Build().ToProto()).
		SetSourceFilePool(sourceFiles).
		Build()
}

// Java: TypedAstSerializer.serializeScriptNode(Node)
func (s *TypedAstSerializer) SerializeScriptNode(script *rhino.Node) *LazyAst {
	jutil.CheckState(script.IsScript())
	s.previousColumn = 0
	s.previousLine = s.previousColumn

	sourceFile := s.GetSourceFilePointer(script)
	scriptProto := s.Visit(script)
	s.subtreeSourceFiles = s.subtreeSourceFiles[:0]

	encodedSourceMap := s.compiler.GetBase64SourceMapContents(*script.GetSourceFileName())

	lazyAstBuilder := LazyAst_NewBuilder().SetScript(scriptProto.ToByteString()).SetSourceFile(int32(sourceFile))

	if encodedSourceMap != nil {
		// This is the encoded source map taken from the inline sourcemap comment. It does not include
		// the base64 prefix.
		// E.g. We serialize "eyJ2ZXJzaW9uI..." from the "//# sourceMappingURL=
		// data:application/json;base64,eyJ2ZXJzaW9uI..." comment.
		lazyAstBuilder.SetSourceMappingUrl(*encodedSourceMap)
	}

	return lazyAstBuilder.Build()
}

// Java: TypedAstSerializer.createWithPositionInfo(Node)
func (s *TypedAstSerializer) CreateWithPositionInfo(n *rhino.Node) *AstNode_Builder {
	jutil.CheckState(n.GetLength() >= 0)
	currentLine := n.GetLineno()
	currentColumn := n.GetCharno()
	builder := AstNode_NewBuilder().
		SetRelativeLine(int32(currentLine - s.previousLine)).
		SetRelativeColumn(int32(currentColumn - s.previousColumn))
	s.previousLine = currentLine
	s.previousColumn = currentColumn
	return builder
}

// Java: TypedAstSerializer.visit(Node)
func (s *TypedAstSerializer) Visit(n *rhino.Node) *AstNode {
	builder := s.CreateWithPositionInfo(n)
	s.AddType(n, builder)
	serializedJsdoc := JSDocSerializer_SerializeJsdoc(n.GetJSDocInfo(), s.stringPool)
	if serializedJsdoc != nil {
		builder.SetJsdoc(serializedJsdoc)
	}
	builder.SetKind(s.KindTranslator(n))
	s.ValueTranslator(builder, n)
	builder.SetBooleanProperties(n.SerializeProperties())
	sourceFile := s.GetSourceFilePointer(n)
	builder.SetSourceFile(int32(sourceFile))

	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		builder.AddChild(s.Visit(child))
	}
	shadowedCode := n.GetClosureUnawareShadow()
	if shadowedCode != nil {
		// For Closure shadow hosts, the ASTNode will get a boolean property bit set
		// (CLOSURE_UNAWARE_SHADOW) that indicates this child ASTNode is not a normal child, but is
		// the contents of the shadow.
		// Shadow roots are structured as
		// ROOT -> SCRIPT -> EXPR_RESULT -> FUNCTION
		// We avoid including the ROOT -> SCRIPT -> EXPR_RESULT structure in the serialization as it
		// is all synthetic code that would unnecessarily bloat the TypedAST and is instead recreated
		// upon deserialization.
		// The child ASTNode is just the FUNCTION.
		builder.AddChild(s.Visit(shadowedCode.GetFirstFirstChild().GetFirstChild()))
	}

	if sourceFile != 0 {
		// ArrayDeque.removeLast() throws NoSuchElementException on an empty deque.
		if len(s.subtreeSourceFiles) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		s.subtreeSourceFiles = s.subtreeSourceFiles[:len(s.subtreeSourceFiles)-1]
	}
	s.SetOriginalName(builder, n)

	return builder.Build()
}

// peekLastSubtreeSourceFile is subtreeSourceFiles.peekLast(): the top of the
// stack, or null when it is empty.
//
// Go-only: language-adapter — ArrayDeque.peekLast on the slice stack.
func (s *TypedAstSerializer) peekLastSubtreeSourceFile() *SourceFile {
	if len(s.subtreeSourceFiles) == 0 {
		return nil
	}
	return s.subtreeSourceFiles[len(s.subtreeSourceFiles)-1]
}

// Java: TypedAstSerializer.getSourceFilePointer(Node)
func (s *TypedAstSerializer) GetSourceFilePointer(n *rhino.Node) int {
	var sourceFile *SourceFile
	if ssf := n.GetStaticSourceFile(); ssf != nil {
		sourceFile = ssf.(*SourceFile)
	}
	if sourceFile == nil {
		// TODO(b/186056977): enforce that SourceFile is not null in externs as well as code.
		if !s.peekLastSubtreeSourceFile().IsExtern() {
			jutil.CheckStateMsg2(
				false,
				"Unexpected null SourceFile for node %s with parent %s",
				n.ToStringTree(),
				n.GetParent())
		}
		return 0 // not set
	}

	if sourceFile == s.peekLastSubtreeSourceFile() {
		// To save space, only serialize a SourceFile for a node if it is different than the parent's
		// source.
		return 0 // not set
	}

	s.subtreeSourceFiles = append(s.subtreeSourceFiles, sourceFile)
	return s.sourceFilePointers.ComputeIfAbsent(
		sourceFile, func(f *SourceFile) int { return 1 + s.sourceFilePointers.Size() })
}

// Java: TypedAstSerializer.addType(Node, AstNode.Builder)
func (s *TypedAstSerializer) AddType(n *rhino.Node, builder *AstNode_Builder) {
	s.typeSerializer.AddTypeForNode(n, builder)
}

// Java: TypedAstSerializer.valueTranslator(AstNode.Builder, Node)
func (s *TypedAstSerializer) ValueTranslator(builder *AstNode_Builder, n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_NAME,
		rhino.Token_STRINGLIT,
		rhino.Token_STRING_KEY,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_LABEL_NAME,
		rhino.Token_IMPORT_STAR:
		builder.SetStringValuePointer(int32(s.stringPool.Put(n.GetString())))
		return
	case rhino.Token_TEMPLATELIT_STRING:
		// The raw string is put in the pool before the cooked one (Java's
		// evaluation order).
		templateStringValue := TemplateStringValue_NewBuilder().
			SetRawStringPointer(int32(s.stringPool.Put(n.GetRawString())))
		cookedStringPointer := -1
		if n.GetCookedString() != nil {
			cookedStringPointer = s.stringPool.Put(*n.GetCookedString())
		}
		builder.SetTemplateStringValue(
			templateStringValue.
				SetCookedStringPointer(int32(cookedStringPointer)).
				Build())
		return
	case rhino.Token_NUMBER:
		builder.SetDoubleValue(n.GetDouble())
		return
	case rhino.Token_BIGINT:
		builder.SetStringValuePointer(int32(s.stringPool.Put(n.GetBigInt().String())))
		return
	default:
		// No value
		return
	}
}

// Java: TypedAstSerializer.setOriginalName(AstNode.Builder, Node)
func (s *TypedAstSerializer) SetOriginalName(builder *AstNode_Builder, n *rhino.Node) {
	originalName := n.GetOriginalName()
	if originalName == nil {
		builder.SetOriginalNamePointer(0) // equivalent to 'not set'
		return
	}

	builder.SetOriginalNamePointer(int32(s.stringPool.Put(*originalName)))
}

// Java: TypedAstSerializer.kindTranslator(Node)
func (s *TypedAstSerializer) KindTranslator(n *rhino.Node) NodeKind {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return NodeKind_SOURCE_FILE
	case rhino.Token_NUMBER:
		return NodeKind_NUMBER_LITERAL
	case rhino.Token_STRINGLIT:
		return NodeKind_STRING_LITERAL
	case rhino.Token_BIGINT:
		return NodeKind_BIGINT_LITERAL
	case rhino.Token_REGEXP:
		return NodeKind_REGEX_LITERAL
	case rhino.Token_FALSE:
		return NodeKind_FALSE
	case rhino.Token_TRUE:
		return NodeKind_TRUE
	case rhino.Token_NULL:
		return NodeKind_NULL
	case rhino.Token_THIS:
		return NodeKind_THIS
	case rhino.Token_VOID:
		return NodeKind_VOID
	case rhino.Token_ARRAYLIT:
		return NodeKind_ARRAY_LITERAL
	case rhino.Token_OBJECTLIT:
		return NodeKind_OBJECT_LITERAL
	case rhino.Token_NAME:
		return NodeKind_IDENTIFIER
	case rhino.Token_ASSIGN:
		return NodeKind_ASSIGNMENT
	case rhino.Token_CALL:
		return NodeKind_CALL
	case rhino.Token_NEW:
		return NodeKind_NEW
	case rhino.Token_GETPROP:
		return NodeKind_PROPERTY_ACCESS
	case rhino.Token_GETELEM:
		return NodeKind_ELEMENT_ACCESS
	case rhino.Token_COMMA:
		return NodeKind_COMMA
	case rhino.Token_OR:
		return NodeKind_BOOLEAN_OR
	case rhino.Token_AND:
		return NodeKind_BOOLEAN_AND
	case rhino.Token_HOOK:
		return NodeKind_HOOK
	case rhino.Token_EQ:
		return NodeKind_EQUAL
	case rhino.Token_NE:
		return NodeKind_NOT_EQUAL
	case rhino.Token_LT:
		return NodeKind_LESS_THAN
	case rhino.Token_LE:
		return NodeKind_LESS_THAN_EQUAL
	case rhino.Token_GT:
		return NodeKind_GREATER_THAN
	case rhino.Token_GE:
		return NodeKind_GREATER_THAN_EQUAL
	case rhino.Token_SHEQ:
		return NodeKind_TRIPLE_EQUAL
	case rhino.Token_SHNE:
		return NodeKind_NOT_TRIPLE_EQUAL
	case rhino.Token_NOT:
		return NodeKind_NOT
	case rhino.Token_POS:
		return NodeKind_POSITIVE
	case rhino.Token_NEG:
		return NodeKind_NEGATIVE
	case rhino.Token_TYPEOF:
		return NodeKind_TYPEOF
	case rhino.Token_INSTANCEOF:
		return NodeKind_INSTANCEOF
	case rhino.Token_IN:
		return NodeKind_IN
	case rhino.Token_ADD:
		return NodeKind_ADD
	case rhino.Token_SUB:
		return NodeKind_SUBTRACT
	case rhino.Token_MUL:
		return NodeKind_MULTIPLY
	case rhino.Token_DIV:
		return NodeKind_DIVIDE
	case rhino.Token_MOD:
		return NodeKind_MODULO
	case rhino.Token_EXPONENT:
		return NodeKind_EXPONENT
	case rhino.Token_BITNOT:
		return NodeKind_BITWISE_NOT
	case rhino.Token_BITOR:
		return NodeKind_BITWISE_OR
	case rhino.Token_BITAND:
		return NodeKind_BITWISE_AND
	case rhino.Token_BITXOR:
		return NodeKind_BITWISE_XOR
	case rhino.Token_LSH:
		return NodeKind_LEFT_SHIFT
	case rhino.Token_RSH:
		return NodeKind_RIGHT_SHIFT
	case rhino.Token_URSH:
		return NodeKind_UNSIGNED_RIGHT_SHIFT
	case rhino.Token_INC:
		if n.GetBooleanProp(rhino.Node_INCRDECR_PROP) {
			return NodeKind_POST_INCREMENT
		}
		return NodeKind_PRE_INCREMENT
	case rhino.Token_DEC:
		if n.GetBooleanProp(rhino.Node_INCRDECR_PROP) {
			return NodeKind_POST_DECREMENT
		}
		return NodeKind_PRE_DECREMENT
	case rhino.Token_ASSIGN_ADD:
		return NodeKind_ASSIGN_ADD
	case rhino.Token_ASSIGN_SUB:
		return NodeKind_ASSIGN_SUBTRACT
	case rhino.Token_ASSIGN_MUL:
		return NodeKind_ASSIGN_MULTIPLY
	case rhino.Token_ASSIGN_DIV:
		return NodeKind_ASSIGN_DIVIDE
	case rhino.Token_ASSIGN_MOD:
		return NodeKind_ASSIGN_MODULO
	case rhino.Token_ASSIGN_EXPONENT:
		return NodeKind_ASSIGN_EXPONENT
	case rhino.Token_ASSIGN_BITOR:
		return NodeKind_ASSIGN_BITWISE_OR
	case rhino.Token_ASSIGN_BITAND:
		return NodeKind_ASSIGN_BITWISE_AND
	case rhino.Token_ASSIGN_BITXOR:
		return NodeKind_ASSIGN_BITWISE_XOR
	case rhino.Token_ASSIGN_LSH:
		return NodeKind_ASSIGN_LEFT_SHIFT
	case rhino.Token_ASSIGN_RSH:
		return NodeKind_ASSIGN_RIGHT_SHIFT
	case rhino.Token_ASSIGN_URSH:
		return NodeKind_ASSIGN_UNSIGNED_RIGHT_SHIFT
	case rhino.Token_YIELD:
		return NodeKind_YIELD
	case rhino.Token_AWAIT:
		return NodeKind_AWAIT
	case rhino.Token_DELPROP:
		return NodeKind_DELETE
	case rhino.Token_TAGGED_TEMPLATELIT:
		return NodeKind_TAGGED_TEMPLATELIT
	case rhino.Token_TEMPLATELIT:
		return NodeKind_TEMPLATELIT
	case rhino.Token_TEMPLATELIT_SUB:
		return NodeKind_TEMPLATELIT_SUB
	case rhino.Token_TEMPLATELIT_STRING:
		return NodeKind_TEMPLATELIT_STRING
	case rhino.Token_NEW_TARGET:
		return NodeKind_NEW_TARGET
	case rhino.Token_COMPUTED_PROP:
		return NodeKind_COMPUTED_PROP
	case rhino.Token_IMPORT_META:
		return NodeKind_IMPORT_META
	case rhino.Token_OPTCHAIN_GETPROP:
		return NodeKind_OPTCHAIN_PROPERTY_ACCESS
	case rhino.Token_OPTCHAIN_CALL:
		return NodeKind_OPTCHAIN_CALL
	case rhino.Token_OPTCHAIN_GETELEM:
		return NodeKind_OPTCHAIN_ELEMENT_ACCESS
	case rhino.Token_COALESCE:
		return NodeKind_COALESCE
	case rhino.Token_DYNAMIC_IMPORT:
		return NodeKind_DYNAMIC_IMPORT
	case rhino.Token_ASSIGN_OR:
		return NodeKind_ASSIGN_OR
	case rhino.Token_ASSIGN_AND:
		return NodeKind_ASSIGN_AND
	case rhino.Token_ASSIGN_COALESCE:
		return NodeKind_ASSIGN_COALESCE
	case rhino.Token_EXPR_RESULT:
		return NodeKind_EXPRESSION_STATEMENT
	case rhino.Token_BREAK:
		return NodeKind_BREAK_STATEMENT
	case rhino.Token_CONTINUE:
		return NodeKind_CONTINUE_STATEMENT
	case rhino.Token_DEBUGGER:
		return NodeKind_DEBUGGER_STATEMENT
	case rhino.Token_DO:
		return NodeKind_DO_STATEMENT
	case rhino.Token_FOR:
		return NodeKind_FOR_STATEMENT
	case rhino.Token_FOR_IN:
		return NodeKind_FOR_IN_STATEMENT
	case rhino.Token_FOR_OF:
		return NodeKind_FOR_OF_STATEMENT
	case rhino.Token_FOR_AWAIT_OF:
		return NodeKind_FOR_AWAIT_OF_STATEMENT
	case rhino.Token_IF:
		return NodeKind_IF_STATEMENT
	case rhino.Token_RETURN:
		return NodeKind_RETURN_STATEMENT
	case rhino.Token_SWITCH:
		return NodeKind_SWITCH_STATEMENT
	case rhino.Token_THROW:
		return NodeKind_THROW_STATEMENT
	case rhino.Token_TRY:
		return NodeKind_TRY_STATEMENT
	case rhino.Token_WHILE:
		return NodeKind_WHILE_STATEMENT
	case rhino.Token_EMPTY:
		return NodeKind_EMPTY
	case rhino.Token_WITH:
		return NodeKind_WITH
	case rhino.Token_IMPORT:
		return NodeKind_IMPORT
	case rhino.Token_EXPORT:
		return NodeKind_EXPORT
	case rhino.Token_VAR:
		return NodeKind_VAR_DECLARATION
	case rhino.Token_CONST:
		return NodeKind_CONST_DECLARATION
	case rhino.Token_LET:
		return NodeKind_LET_DECLARATION
	case rhino.Token_FUNCTION:
		return NodeKind_FUNCTION_LITERAL
	case rhino.Token_CLASS:
		return NodeKind_CLASS_LITERAL
	case rhino.Token_BLOCK:
		return NodeKind_BLOCK
	case rhino.Token_LABEL:
		return NodeKind_LABELED_STATEMENT
	case rhino.Token_LABEL_NAME:
		return NodeKind_LABELED_NAME
	case rhino.Token_CLASS_MEMBERS:
		return NodeKind_CLASS_MEMBERS
	case rhino.Token_MEMBER_FUNCTION_DEF:
		return NodeKind_METHOD_DECLARATION
	case rhino.Token_MEMBER_FIELD_DEF:
		return NodeKind_FIELD_DECLARATION
	case rhino.Token_COMPUTED_FIELD_DEF:
		return NodeKind_COMPUTED_PROP_FIELD
	case rhino.Token_PARAM_LIST:
		return NodeKind_PARAMETER_LIST
	case rhino.Token_STRING_KEY:
		if n.IsQuotedStringKey() {
			return NodeKind_QUOTED_STRING_KEY
		}
		return NodeKind_RENAMABLE_STRING_KEY
	case rhino.Token_CASE:
		return NodeKind_CASE
	case rhino.Token_DEFAULT_CASE:
		return NodeKind_DEFAULT_CASE
	case rhino.Token_CATCH:
		return NodeKind_CATCH
	case rhino.Token_SUPER:
		return NodeKind_SUPER
	case rhino.Token_ARRAY_PATTERN:
		return NodeKind_ARRAY_PATTERN
	case rhino.Token_OBJECT_PATTERN:
		return NodeKind_OBJECT_PATTERN
	case rhino.Token_DESTRUCTURING_LHS:
		return NodeKind_DESTRUCTURING_LHS
	case rhino.Token_DEFAULT_VALUE:
		return NodeKind_DEFAULT_VALUE
	case rhino.Token_GETTER_DEF:
		if n.IsQuotedStringKey() {
			return NodeKind_QUOTED_GETTER_DEF
		}
		return NodeKind_RENAMABLE_GETTER_DEF
	case rhino.Token_SETTER_DEF:
		if n.IsQuotedStringKey() {
			return NodeKind_QUOTED_SETTER_DEF
		}
		return NodeKind_RENAMABLE_SETTER_DEF
	case rhino.Token_IMPORT_SPECS:
		return NodeKind_IMPORT_SPECS
	case rhino.Token_IMPORT_SPEC:
		return NodeKind_IMPORT_SPEC
	case rhino.Token_IMPORT_STAR:
		return NodeKind_IMPORT_STAR
	case rhino.Token_EXPORT_SPECS:
		return NodeKind_EXPORT_SPECS
	case rhino.Token_EXPORT_SPEC:
		return NodeKind_EXPORT_SPEC
	case rhino.Token_MODULE_BODY:
		return NodeKind_MODULE_BODY
	case rhino.Token_ITER_REST:
		return NodeKind_ITER_REST
	case rhino.Token_ITER_SPREAD:
		return NodeKind_ITER_SPREAD
	case rhino.Token_OBJECT_REST:
		return NodeKind_OBJECT_REST
	case rhino.Token_OBJECT_SPREAD:
		return NodeKind_OBJECT_SPREAD
	case rhino.Token_SWITCH_BODY:
		return NodeKind_SWITCH_BODY
	// Explicitly unsupported token types. Not serialized.
	case rhino.Token_ROOT,
		// TS type tokens
		rhino.Token_STRING_TYPE,
		rhino.Token_BOOLEAN_TYPE,
		rhino.Token_NUMBER_TYPE,
		rhino.Token_FUNCTION_TYPE,
		rhino.Token_PARAMETERIZED_TYPE,
		rhino.Token_UNION_TYPE,
		rhino.Token_ANY_TYPE,
		rhino.Token_NULLABLE_TYPE,
		rhino.Token_VOID_TYPE,
		rhino.Token_REST_PARAMETER_TYPE,
		rhino.Token_NAMED_TYPE,
		rhino.Token_RECORD_TYPE,
		rhino.Token_UNDEFINED_TYPE,
		rhino.Token_OPTIONAL_PARAMETER,
		rhino.Token_ARRAY_TYPE,
		rhino.Token_GENERIC_TYPE,
		rhino.Token_GENERIC_TYPE_LIST,
		rhino.Token_INTERFACE,
		rhino.Token_INTERFACE_EXTENDS,
		rhino.Token_INTERFACE_MEMBERS,
		rhino.Token_MEMBER_VARIABLE_DEF,
		rhino.Token_ENUM,
		rhino.Token_ENUM_MEMBERS,
		rhino.Token_IMPLEMENTS,
		rhino.Token_TYPE_ALIAS,
		rhino.Token_DECLARE,
		rhino.Token_INDEX_SIGNATURE,
		rhino.Token_CALL_SIGNATURE,
		rhino.Token_NAMESPACE,
		rhino.Token_NAMESPACE_ELEMENTS,
		// JSDoc tokens
		rhino.Token_ANNOTATION,
		rhino.Token_PIPE,
		rhino.Token_STAR,
		rhino.Token_EOC,
		rhino.Token_QMARK,
		rhino.Token_BANG,
		rhino.Token_EQUALS,
		rhino.Token_LB,
		rhino.Token_LC,
		rhino.Token_COLON,
		// Closure-type-system-specific token
		rhino.Token_CAST,
		// Unused tokens
		rhino.Token_PLACEHOLDER1,
		rhino.Token_PLACEHOLDER2,
		rhino.Token_PLACEHOLDER3:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unserializable token for node: " + n.String()))
}

// TypedAstSerializer_TypeSerializer is used to provide TypePointers for
// serializing Nodes and to generate the TypePool.
//
// Java: interface TypedAstSerializer.TypeSerializer
type TypedAstSerializer_TypeSerializer interface {
	// AddTypeForNode adds, if appropriate for `node`, an `Integer` to
	// `astNodeBuilder`.
	//
	// Java: TypedAstSerializer.TypeSerializer.addTypeForNode(Node, AstNode.Builder)
	AddTypeForNode(node *rhino.Node, astNodeBuilder *AstNode_Builder)

	// GenerateTypePool returns a `TypePool` containing the types used by
	// `addTypeForNode()`.
	//
	// Java: TypedAstSerializer.TypeSerializer.generateTypePool()
	GenerateTypePool() *TypePool
}

// TypedAstSerializer_NoOpTypeSerializer is used when type checking has not
// been done.
//
// Java: private static class TypedAstSerializer.NoOpTypeSerializer
type TypedAstSerializer_NoOpTypeSerializer struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: TypedAstSerializer.NoOpTypeSerializer()
func NewTypedAstSerializer_NoOpTypeSerializer() *TypedAstSerializer_NoOpTypeSerializer {
	return &TypedAstSerializer_NoOpTypeSerializer{}
}

// Java: TypedAstSerializer.NoOpTypeSerializer.addTypeForNode(Node, AstNode.Builder)
func (*TypedAstSerializer_NoOpTypeSerializer) AddTypeForNode(node *rhino.Node, astNodeBuilder *AstNode_Builder) {
	// Do nothing.
}

// Java: TypedAstSerializer.NoOpTypeSerializer.generateTypePool()
func (*TypedAstSerializer_NoOpTypeSerializer) GenerateTypePool() *TypePool {
	return TypePool_GetDefaultInstance() // empty TypePool
}

// TypedAstSerializer_CreateJSTypeSerializer creates the `TypeSerializer`
// appropriate for an AST that contains JSTypes.
//
// Java: TypedAstSerializer.createJSTypeSerializer(AbstractCompiler, SerializationOptions, StringPool.Builder)
func TypedAstSerializer_CreateJSTypeSerializer(
	compiler AbstractCompiler,
	serializationMode *SerializationOptions,
	stringPoolBuilder *StringPool_Builder) *TypedAstSerializer_JSTypeSerializer {
	serializeTypesToPointers := SerializeTypesToPointers_Create(compiler, stringPoolBuilder, serializationMode)
	// Gather and serialize all the types now.
	serializeTypesToPointers.GatherTypesOnAst(compiler.GetRoot())
	return NewTypedAstSerializer_JSTypeSerializer(
		serializeTypesToPointers.GetTypePointersByJstype(), serializeTypesToPointers.GetTypePool())
}

// TypedAstSerializer_JSTypeSerializer is used when the AST's JSTypes have not
// been converted to Colors.
//
// Java: private static class TypedAstSerializer.JSTypeSerializer
type TypedAstSerializer_JSTypeSerializer struct {
	// Everything is pre-calculated with this form of serialization.
	typesToPointers *base.LinkedIdentityHashMap[rhino.JSType, int]
	typePool        *TypePool
}

// Java: TypedAstSerializer.JSTypeSerializer(LinkedIdentityHashMap<JSType, Integer>, TypePool)
func NewTypedAstSerializer_JSTypeSerializer(
	typesToPointers *base.LinkedIdentityHashMap[rhino.JSType, int], typePool *TypePool) *TypedAstSerializer_JSTypeSerializer {
	return &TypedAstSerializer_JSTypeSerializer{typesToPointers: typesToPointers, typePool: typePool}
}

// Java: TypedAstSerializer.JSTypeSerializer.addTypeForNode(Node, AstNode.Builder)
func (s *TypedAstSerializer_JSTypeSerializer) AddTypeForNode(node *rhino.Node, astNodeBuilder *AstNode_Builder) {
	type_ := node.GetJSType()
	if type_ != nil {
		pointer, ok := s.typesToPointers.Lookup(type_)
		if !ok {
			// checkNotNull(typesToPointers.get(type), "cannot find pointer for %s", type)
			jutil.CheckNotNullMsg1[any](nil, "cannot find pointer for %s", type_)
		}
		astNodeBuilder.SetType(int32(pointer))
	}
}

// Java: TypedAstSerializer.JSTypeSerializer.generateTypePool()
func (s *TypedAstSerializer_JSTypeSerializer) GenerateTypePool() *TypePool {
	return s.typePool
}

// CreateColorTypeSerializer creates a `TypeSerializer` that knows how to
// serialize `Color`s from the AST.
//
// Java: TypedAstSerializer.createColorTypeSerializer(AbstractCompiler, SerializationOptions, StringPool.Builder)
func (s *TypedAstSerializer) CreateColorTypeSerializer(
	compiler AbstractCompiler,
	serializationMode *SerializationOptions,
	stringPoolBuilder *StringPool_Builder) *TypedAstSerializer_ColorTypeSerializer {
	// Gather all the property names that are actually used in the AST.
	usedPropertyNames := TypedAstSerializer_CollectUsedPropertyNames(compiler)
	colorSerializer := NewColorSerializer(
		serializationMode,
		// lookup / allocate strings from the shared string pool
		stringPoolBuilder.Put,
		// only include property names known to be used in the AST
		usedPropertyNames.Contains)
	return NewTypedAstSerializer_ColorTypeSerializer(colorSerializer, compiler.GetColorRegistry())
}

// TypedAstSerializer_CollectUsedPropertyNames adds all unquoted property
// names appearing in the AST.
//
// Java: TypedAstSerializer.collectUsedPropertyNames(AbstractCompiler)
func TypedAstSerializer_CollectUsedPropertyNames(compiler AbstractCompiler) *jutil.LinkedHashSet[string] {
	propertyNamesBuilder := jutil.NewLinkedHashSet[string]()

	NodeTraversal_Traverse(
		compiler,
		compiler.GetRoot(),
		&typedAstSerializer_collectUsedPropertyNamesCallback{propertyNamesBuilder: propertyNamesBuilder})
	return propertyNamesBuilder
}

// typedAstSerializer_collectUsedPropertyNamesCallback is the anonymous
// NodeTraversal.Callback of collectUsedPropertyNames.
type typedAstSerializer_collectUsedPropertyNamesCallback struct {
	propertyNamesBuilder *jutil.LinkedHashSet[string]
}

// Java: TypedAstSerializer.collectUsedPropertyNames(...).new Callback() {...}.shouldTraverse(NodeTraversal, Node, Node)
func (c *typedAstSerializer_collectUsedPropertyNamesCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// We won't serialize type summary files, so we don't care about property names
	// appearing in them.
	return !n.IsScript() || !NodeUtil_IsFromTypeSummary(n)
}

// Java: TypedAstSerializer.collectUsedPropertyNames(...).new Callback() {...}.visit(NodeTraversal, Node, Node)
func (c *typedAstSerializer_collectUsedPropertyNamesCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP, // "name" from (someObject.name)
		rhino.Token_OPTCHAIN_GETPROP: // "name" from (someObject?.name)
		c.propertyNamesBuilder.Add(n.GetString())
	case rhino.Token_STRING_KEY, // "name" from obj = {name: 0}
		rhino.Token_MEMBER_FUNCTION_DEF, // "name" from class C { name() {} }
		rhino.Token_MEMBER_FIELD_DEF,    // "name" from class C { name = 0; }
		rhino.Token_GETTER_DEF,          // "name" from class C { get name() {} }
		rhino.Token_SETTER_DEF:          // "name" from class C { set name(n) {} }
		if !n.IsQuotedStringKey() {
			c.propertyNamesBuilder.Add(n.GetString())
		}
	default:
	}
}

// TypedAstSerializer_ColorTypeSerializer is used when the AST has `Color`s
// rather than `JSType`s.
//
// Java: private static class TypedAstSerializer.ColorTypeSerializer
type TypedAstSerializer_ColorTypeSerializer struct {
	colorSerializer *ColorSerializer
	colorRegistry   *colors.ColorRegistry
}

// Java: TypedAstSerializer.ColorTypeSerializer(ColorSerializer, ColorRegistry)
func NewTypedAstSerializer_ColorTypeSerializer(colorSerializer *ColorSerializer, colorRegistry *colors.ColorRegistry) *TypedAstSerializer_ColorTypeSerializer {
	return &TypedAstSerializer_ColorTypeSerializer{colorSerializer: colorSerializer, colorRegistry: colorRegistry}
}

// Java: TypedAstSerializer.ColorTypeSerializer.addTypeForNode(Node, AstNode.Builder)
func (s *TypedAstSerializer_ColorTypeSerializer) AddTypeForNode(node *rhino.Node, astNodeBuilder *AstNode_Builder) {
	color := node.GetColor()
	if color != nil {
		astNodeBuilder.SetType(int32(s.colorSerializer.AddColor(color)))
	}
}

// Java: TypedAstSerializer.ColorTypeSerializer.generateTypePool()
func (s *TypedAstSerializer_ColorTypeSerializer) GenerateTypePool() *TypePool {
	mismatchLocationsForDebugging := s.colorRegistry.GetMismatchLocationsForDebugging()
	return s.colorSerializer.GenerateTypePool(
		s.colorRegistry.GetDisambiguationSupertypes,
		func(color *colors.Color) *jutil.LinkedHashSet[string] {
			// Go-only: language-adapter — the ColorSerializer port takes the
			// ImmutableSet<String> as a LinkedHashSet, the SetMultimap port
			// returns an ObjectLinkedHashSet: copied in iteration order.
			locations := jutil.NewLinkedHashSet[string]()
			for location := range mismatchLocationsForDebugging.Get(color.GetId()).All() {
				locations.Add(location)
			}
			return locations
		})
}
