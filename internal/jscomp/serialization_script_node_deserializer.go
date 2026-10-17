// Port of com/google/javascript/jscomp/serialization/ScriptNodeDeserializer.java @ closure-compiler v20260415

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ScriptNodeDeserializer deserializes an AstNode-tree representing a SCRIPT
// into a Node-tree.
//
// This process depends on other information from the TypedAST format, but
// the output it limited to only a single SCRIPT. The other deserialized
// content must be provided beforehand.
type ScriptNodeDeserializer struct {
	sourceFile       *SourceFile
	scriptBytes      []byte
	sourceMappingURL string
	colorPoolShard   *ColorPool_ShardView // Optional: nil when absent
	stringPool       *StringPool
	filePool         []*SourceFile
}

// Java: ScriptNodeDeserializer(LazyAst, StringPool, Optional<ColorPool.ShardView>, ImmutableList<SourceFile>)
func NewScriptNodeDeserializer(
	ast *LazyAst,
	stringPool *StringPool,
	colorPoolShard *ColorPool_ShardView,
	filePool []*SourceFile) *ScriptNodeDeserializer {
	return &ScriptNodeDeserializer{
		scriptBytes:      ast.GetScript(),
		sourceFile:       filePool[ast.GetSourceFile()-1],
		sourceMappingURL: ast.GetSourceMappingUrl(),
		colorPoolShard:   colorPoolShard,
		stringPool:       stringPool,
		filePool:         filePool,
	}
}

// Java: ScriptNodeDeserializer.deserializeNew()
func (d *ScriptNodeDeserializer) DeserializeNew() *rhino.Node {
	return newScriptNodeDeserializer_Runner(d).Run()
}

// Java: ScriptNodeDeserializer.Runner
type scriptNodeDeserializer_Runner struct {
	owner_           *ScriptNodeDeserializer
	featureCollector *parsing.FeatureCollector
	previousLine     int
	previousColumn   int
}

func newScriptNodeDeserializer_Runner(owner *ScriptNodeDeserializer) *scriptNodeDeserializer_Runner {
	return &scriptNodeDeserializer_Runner{owner_: owner, featureCollector: parsing.NewFeatureCollector()}
}

// Java: ScriptNodeDeserializer.Runner.run()
func (r *scriptNodeDeserializer_Runner) Run() (result *rhino.Node) {
	jutil.Try(func() {
		// The real limit is stack space.
		astStream := &protoReader{b: r.Owner().scriptBytes, limit: math.MaxInt32}

		contextStack := jutil.NewArrayList[parsing.FeatureCollector_FeatureContext]()
		contextStack.Add(parsing.FeatureCollector_FeatureContext_NONE)

		astNode := &AstNode{memoizedSize: -1}
		astNode.mergeFromReader(astStream)
		astNode.frozen = true
		scriptNode := r.Visit(
			astNode,
			contextStack,
			r.Owner().CreateSourceInfoTemplate(r.Owner().sourceFile))
		scriptNode.PutProp(rhino.Node_FEATURE_SET, r.featureCollector.AllFeatures())
		result = scriptNode
	}, func(ex jutil.IIOException) {
		panic(NewMalformedTypedAstExceptionWithCause(r.Owner().sourceFile, ex))
	})
	return result
}

// Java: ScriptNodeDeserializer.Runner.owner()
func (r *scriptNodeDeserializer_Runner) Owner() *ScriptNodeDeserializer {
	return r.owner_
}

// Visit visits a single AST node and returns the corresponding Node.
// contextStack is the context stack to use for the node; sourceFileTemplate
// (may be nil) is a template node to use for the node.
//
// Java: ScriptNodeDeserializer.Runner.visit(AstNode, List<FeatureContext>, Node)
func (r *scriptNodeDeserializer_Runner) Visit(
	astNode *AstNode, contextStack *jutil.ArrayList[parsing.FeatureCollector_FeatureContext], sourceFileTemplate *rhino.Node) *rhino.Node {
	if sourceFileTemplate == nil || astNode.GetSourceFile() != 0 {
		// 0 == 'not set'
		sourceFileTemplate =
			r.Owner().
				CreateSourceInfoTemplate(r.Owner().filePool[astNode.GetSourceFile()-1])
	}

	currentLine := r.previousLine + int(astNode.GetRelativeLine())
	currentColumn := r.previousColumn + int(astNode.GetRelativeColumn())

	n := r.Owner().DeserializeSingleNode(astNode, sourceFileTemplate)
	if astNode.HasType() && r.Owner().colorPoolShard != nil {
		n.SetColor(r.Owner().colorPoolShard.GetColor(int(astNode.GetType())))
	}
	properties := astNode.GetBooleanProperties()
	if properties > 0 {
		n.DeserializeProperties(
			r.Owner().FilterOutCastProp(properties), sourceFileTemplate.GetIsInClosureUnawareSubtree())
	}
	n.SetJSDocInfo(JSDocSerializer_DeserializeJsdoc(astNode.GetJsdoc(), r.Owner().stringPool))
	n.SetLinenoCharno(currentLine, currentColumn)
	r.previousLine = currentLine
	r.previousColumn = currentColumn
	if contextStack != nil {
		// Note: We use unmodifiableList() instead of ImmutableList for performance reasons.
		// (Go-only: language-adapter — the list itself is passed; FeatureCollector does not
		// modify it.)
		contextStack.Add(
			r.featureCollector.VisitSingleNode(contextStack, n))
	}

	if rhino.Node_HasBitSet(properties, int(rhino.NodeProperty_CLOSURE_UNAWARE_SHADOW.GetNumber())) {
		serializedShadowChild := astNode.GetChild(0)
		closureUnawareSubtreeTemplateNode := sourceFileTemplate.CloneTree()
		closureUnawareSubtreeTemplateNode.SetIsInClosureUnawareSubtree(true)
		// Note: we intentionally pass `contextStack` to visit() in order to capture any
		// language features present in the shadow AST, like const syntax, into the main AST
		// featureset. We are doing this to align ScriptNodeDeserializer with the existing behavior
		// in IRFactory. (In the future we could consider instead giving each shadow AST script
		// its own separate featureset, but that would be a larger change.)
		shadowedCode := r.Visit(serializedShadowChild, contextStack, closureUnawareSubtreeTemplateNode)
		r.Owner().SetOriginalNameIfPresent(serializedShadowChild, shadowedCode)
		// The shadowed code is only the "source" parts of the shadow structure, and does not
		// include the synthetic code that is needed for the compiler to consider it a valid
		// standalone AST. We recreate that here.
		// This must be kept in sync with the shadow structure created by TypedAstSerializer.
		shadowRoot := rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(shadowedCode)))
		shadowRoot.ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		shadowRoot.GetFirstChild().ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		shadowRoot.GetFirstFirstChild().ClonePropsFrom(closureUnawareSubtreeTemplateNode)

		n.SetClosureUnawareShadow(shadowRoot)
		return n
	}
	children := astNode.GetChildCount()

	for i := 0; i < children; i++ {
		child := astNode.GetChild(i)
		deserializedChild := r.Visit(child, contextStack, sourceFileTemplate)
		n.AddChildToBack(deserializedChild)
		r.Owner().SetOriginalNameIfPresent(child, deserializedChild)
	}

	contextStack.RemoveLast()

	return n
}

// Java: ScriptNodeDeserializer.getSourceMappingURL()
func (d *ScriptNodeDeserializer) GetSourceMappingURL() string {
	return d.sourceMappingURL
}

// Java: ScriptNodeDeserializer.getSourceFile()
func (d *ScriptNodeDeserializer) GetSourceFile() *SourceFile {
	return d.sourceFile
}

// CreateSourceInfoTemplate creates a template node to use as a source of
// common attributes.
//
// This allows the prop structure to be shared among all the node from this
// source file. This reduces the cost of these properties to O(nodes) to
// O(files).
//
// Java: ScriptNodeDeserializer.createSourceInfoTemplate(SourceFile)
func (d *ScriptNodeDeserializer) CreateSourceInfoTemplate(file *SourceFile) *rhino.Node {
	// The Node type choice is arbitrary.
	sourceInfoTemplate := rhino.NewNode(rhino.Token_SCRIPT)
	sourceInfoTemplate.SetStaticSourceFile(file)
	return sourceInfoTemplate
}

// Java: ScriptNodeDeserializer.setOriginalNameIfPresent(AstNode, Node)
func (d *ScriptNodeDeserializer) SetOriginalNameIfPresent(astNode *AstNode, n *rhino.Node) {
	if astNode.GetOriginalNamePointer() != 0 {
		n.SetOriginalNameFromStringPool(
			d.stringPool.GetInternedStrings(), int(astNode.GetOriginalNamePointer()))
	}
}

// StringNode creates a new string node with the given token & string value
// of the AstNode.
//
// Prefer calling this method over calling a regular Node.* or IR.* method
// when possible. This method integrates with RhinoStringPool to cache String
// interning results.
//
// Java: ScriptNodeDeserializer.stringNode(Token, AstNode, Node)
func (d *ScriptNodeDeserializer) StringNode(token rhino.Token, n *AstNode, templateNode *rhino.Node) *rhino.Node {
	str := rhino.Node_NewStringFromStringPool(token, d.stringPool.GetInternedStrings(), int(n.GetStringValuePointer()))
	str.ClonePropsFrom(templateNode)
	return str
}

// Java: ScriptNodeDeserializer.newNodeFromTemplate(Token, Node)
func (d *ScriptNodeDeserializer) NewNodeFromTemplate(token rhino.Token, templateNode *rhino.Node) *rhino.Node {
	n := templateNode.CloneTree()
	n.SetToken(token)
	return n
}

// Java: ScriptNodeDeserializer.deserializeSingleNode(AstNode, Node)
func (d *ScriptNodeDeserializer) DeserializeSingleNode(n *AstNode, templateNode *rhino.Node) *rhino.Node {
	switch n.GetKind() {
	case NodeKind_SOURCE_FILE:
		return d.NewNodeFromTemplate(rhino.Token_SCRIPT, templateNode)

	case NodeKind_NUMBER_LITERAL:
		numberNode := rhino.IR_Number(n.GetDoubleValue())
		numberNode.ClonePropsFrom(templateNode)
		return numberNode
	case NodeKind_STRING_LITERAL:
		return d.StringNode(rhino.Token_STRINGLIT, n, templateNode)
	case NodeKind_IDENTIFIER:
		return d.StringNode(rhino.Token_NAME, n, templateNode)
	case NodeKind_FALSE:
		return d.NewNodeFromTemplate(rhino.Token_FALSE, templateNode)
	case NodeKind_TRUE:
		return d.NewNodeFromTemplate(rhino.Token_TRUE, templateNode)
	case NodeKind_NULL:
		return d.NewNodeFromTemplate(rhino.Token_NULL, templateNode)
	case NodeKind_THIS:
		return d.NewNodeFromTemplate(rhino.Token_THIS, templateNode)
	case NodeKind_VOID:
		return d.NewNodeFromTemplate(rhino.Token_VOID, templateNode)
	case NodeKind_BIGINT_LITERAL:
		bigintString := d.stringPool.Get(int(n.GetStringValuePointer()))
		bi := rhino.IR_Bigint(jutil.NewBigIntegerFromString(bigintString))
		bi.ClonePropsFrom(templateNode)
		return bi
	case NodeKind_REGEX_LITERAL:
		return d.NewNodeFromTemplate(rhino.Token_REGEXP, templateNode)
	case NodeKind_ARRAY_LITERAL:
		return d.NewNodeFromTemplate(rhino.Token_ARRAYLIT, templateNode)
	case NodeKind_OBJECT_LITERAL:
		return d.NewNodeFromTemplate(rhino.Token_OBJECTLIT, templateNode)

	case NodeKind_ASSIGNMENT:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN, templateNode)
	case NodeKind_CALL:
		return d.NewNodeFromTemplate(rhino.Token_CALL, templateNode)
	case NodeKind_NEW:
		return d.NewNodeFromTemplate(rhino.Token_NEW, templateNode)
	case NodeKind_PROPERTY_ACCESS:
		return d.StringNode(rhino.Token_GETPROP, n, templateNode)
	case NodeKind_ELEMENT_ACCESS:
		return d.NewNodeFromTemplate(rhino.Token_GETELEM, templateNode)

	case NodeKind_COMMA:
		return d.NewNodeFromTemplate(rhino.Token_COMMA, templateNode)
	case NodeKind_BOOLEAN_OR:
		return d.NewNodeFromTemplate(rhino.Token_OR, templateNode)
	case NodeKind_BOOLEAN_AND:
		return d.NewNodeFromTemplate(rhino.Token_AND, templateNode)
	case NodeKind_HOOK:
		return d.NewNodeFromTemplate(rhino.Token_HOOK, templateNode)
	case NodeKind_EQUAL:
		return d.NewNodeFromTemplate(rhino.Token_EQ, templateNode)
	case NodeKind_NOT_EQUAL:
		return d.NewNodeFromTemplate(rhino.Token_NE, templateNode)
	case NodeKind_LESS_THAN:
		return d.NewNodeFromTemplate(rhino.Token_LT, templateNode)
	case NodeKind_LESS_THAN_EQUAL:
		return d.NewNodeFromTemplate(rhino.Token_LE, templateNode)
	case NodeKind_GREATER_THAN:
		return d.NewNodeFromTemplate(rhino.Token_GT, templateNode)
	case NodeKind_GREATER_THAN_EQUAL:
		return d.NewNodeFromTemplate(rhino.Token_GE, templateNode)
	case NodeKind_TRIPLE_EQUAL:
		return d.NewNodeFromTemplate(rhino.Token_SHEQ, templateNode)
	case NodeKind_NOT_TRIPLE_EQUAL:
		return d.NewNodeFromTemplate(rhino.Token_SHNE, templateNode)
	case NodeKind_NOT:
		return d.NewNodeFromTemplate(rhino.Token_NOT, templateNode)
	case NodeKind_POSITIVE:
		return d.NewNodeFromTemplate(rhino.Token_POS, templateNode)
	case NodeKind_NEGATIVE:
		return d.NewNodeFromTemplate(rhino.Token_NEG, templateNode)
	case NodeKind_TYPEOF:
		return d.NewNodeFromTemplate(rhino.Token_TYPEOF, templateNode)
	case NodeKind_INSTANCEOF:
		return d.NewNodeFromTemplate(rhino.Token_INSTANCEOF, templateNode)
	case NodeKind_IN:
		return d.NewNodeFromTemplate(rhino.Token_IN, templateNode)

	case NodeKind_ADD:
		return d.NewNodeFromTemplate(rhino.Token_ADD, templateNode)
	case NodeKind_SUBTRACT:
		return d.NewNodeFromTemplate(rhino.Token_SUB, templateNode)
	case NodeKind_MULTIPLY:
		return d.NewNodeFromTemplate(rhino.Token_MUL, templateNode)
	case NodeKind_DIVIDE:
		return d.NewNodeFromTemplate(rhino.Token_DIV, templateNode)
	case NodeKind_MODULO:
		return d.NewNodeFromTemplate(rhino.Token_MOD, templateNode)
	case NodeKind_EXPONENT:
		return d.NewNodeFromTemplate(rhino.Token_EXPONENT, templateNode)
	case NodeKind_BITWISE_NOT:
		return d.NewNodeFromTemplate(rhino.Token_BITNOT, templateNode)
	case NodeKind_BITWISE_OR:
		return d.NewNodeFromTemplate(rhino.Token_BITOR, templateNode)
	case NodeKind_BITWISE_AND:
		return d.NewNodeFromTemplate(rhino.Token_BITAND, templateNode)
	case NodeKind_BITWISE_XOR:
		return d.NewNodeFromTemplate(rhino.Token_BITXOR, templateNode)
	case NodeKind_LEFT_SHIFT:
		return d.NewNodeFromTemplate(rhino.Token_LSH, templateNode)
	case NodeKind_RIGHT_SHIFT:
		return d.NewNodeFromTemplate(rhino.Token_RSH, templateNode)
	case NodeKind_UNSIGNED_RIGHT_SHIFT:
		return d.NewNodeFromTemplate(rhino.Token_URSH, templateNode)
	case NodeKind_PRE_INCREMENT:
		return d.NewNodeFromTemplate(rhino.Token_INC, templateNode)
	case NodeKind_POST_INCREMENT:
		postInc := d.NewNodeFromTemplate(rhino.Token_INC, templateNode)
		postInc.PutBooleanProp(rhino.Node_INCRDECR_PROP, true)
		return postInc
	case NodeKind_PRE_DECREMENT:
		return d.NewNodeFromTemplate(rhino.Token_DEC, templateNode)
	case NodeKind_POST_DECREMENT:
		postDec := d.NewNodeFromTemplate(rhino.Token_DEC, templateNode)
		postDec.PutBooleanProp(rhino.Node_INCRDECR_PROP, true)
		return postDec

	case NodeKind_ASSIGN_ADD:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_ADD, templateNode)
	case NodeKind_ASSIGN_SUBTRACT:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_SUB, templateNode)
	case NodeKind_ASSIGN_MULTIPLY:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_MUL, templateNode)
	case NodeKind_ASSIGN_DIVIDE:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_DIV, templateNode)
	case NodeKind_ASSIGN_MODULO:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_MOD, templateNode)
	case NodeKind_ASSIGN_EXPONENT:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_EXPONENT, templateNode)
	case NodeKind_ASSIGN_BITWISE_OR:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_BITOR, templateNode)
	case NodeKind_ASSIGN_BITWISE_AND:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_BITAND, templateNode)
	case NodeKind_ASSIGN_BITWISE_XOR:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_BITXOR, templateNode)
	case NodeKind_ASSIGN_LEFT_SHIFT:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_LSH, templateNode)
	case NodeKind_ASSIGN_RIGHT_SHIFT:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_RSH, templateNode)
	case NodeKind_ASSIGN_UNSIGNED_RIGHT_SHIFT:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_URSH, templateNode)

	case NodeKind_YIELD:
		return d.NewNodeFromTemplate(rhino.Token_YIELD, templateNode)
	case NodeKind_AWAIT:
		return d.NewNodeFromTemplate(rhino.Token_AWAIT, templateNode)
	case NodeKind_DELETE:
		return d.NewNodeFromTemplate(rhino.Token_DELPROP, templateNode)
	case NodeKind_TAGGED_TEMPLATELIT:
		return d.NewNodeFromTemplate(rhino.Token_TAGGED_TEMPLATELIT, templateNode)
	case NodeKind_TEMPLATELIT:
		return d.NewNodeFromTemplate(rhino.Token_TEMPLATELIT, templateNode)
	case NodeKind_TEMPLATELIT_SUB:
		return d.NewNodeFromTemplate(rhino.Token_TEMPLATELIT_SUB, templateNode)
	case NodeKind_TEMPLATELIT_STRING:
		templateStringValue := n.GetTemplateStringValue()
		templateLitString :=
			rhino.Node_NewTemplateLitStringFromStringPool(
				d.stringPool.GetInternedStrings(),
				int(templateStringValue.GetCookedStringPointer()),
				int(templateStringValue.GetRawStringPointer()))
		templateLitString.ClonePropsFrom(templateNode)
		return templateLitString
	case NodeKind_NEW_TARGET:
		return d.NewNodeFromTemplate(rhino.Token_NEW_TARGET, templateNode)
	case NodeKind_COMPUTED_PROP:
		return d.NewNodeFromTemplate(rhino.Token_COMPUTED_PROP, templateNode)
	case NodeKind_IMPORT_META:
		return d.NewNodeFromTemplate(rhino.Token_IMPORT_META, templateNode)
	case NodeKind_OPTCHAIN_PROPERTY_ACCESS:
		return d.StringNode(rhino.Token_OPTCHAIN_GETPROP, n, templateNode)
	case NodeKind_OPTCHAIN_CALL:
		return d.NewNodeFromTemplate(rhino.Token_OPTCHAIN_CALL, templateNode)
	case NodeKind_OPTCHAIN_ELEMENT_ACCESS:
		return d.NewNodeFromTemplate(rhino.Token_OPTCHAIN_GETELEM, templateNode)
	case NodeKind_COALESCE:
		return d.NewNodeFromTemplate(rhino.Token_COALESCE, templateNode)
	case NodeKind_DYNAMIC_IMPORT:
		return d.NewNodeFromTemplate(rhino.Token_DYNAMIC_IMPORT, templateNode)

	case NodeKind_ASSIGN_OR:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_OR, templateNode)
	case NodeKind_ASSIGN_AND:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_AND, templateNode)
	case NodeKind_ASSIGN_COALESCE:
		return d.NewNodeFromTemplate(rhino.Token_ASSIGN_COALESCE, templateNode)

	case NodeKind_EXPRESSION_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_EXPR_RESULT, templateNode)
	case NodeKind_BREAK_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_BREAK, templateNode)
	case NodeKind_CONTINUE_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_CONTINUE, templateNode)
	case NodeKind_DEBUGGER_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_DEBUGGER, templateNode)
	case NodeKind_DO_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_DO, templateNode)
	case NodeKind_FOR_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_FOR, templateNode)
	case NodeKind_FOR_IN_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_FOR_IN, templateNode)
	case NodeKind_FOR_OF_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_FOR_OF, templateNode)
	case NodeKind_FOR_AWAIT_OF_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_FOR_AWAIT_OF, templateNode)
	case NodeKind_IF_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_IF, templateNode)
	case NodeKind_RETURN_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_RETURN, templateNode)
	case NodeKind_SWITCH_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_SWITCH, templateNode)
	case NodeKind_SWITCH_BODY:
		return d.NewNodeFromTemplate(rhino.Token_SWITCH_BODY, templateNode)
	case NodeKind_THROW_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_THROW, templateNode)
	case NodeKind_TRY_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_TRY, templateNode)
	case NodeKind_WHILE_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_WHILE, templateNode)
	case NodeKind_EMPTY:
		return d.NewNodeFromTemplate(rhino.Token_EMPTY, templateNode)
	case NodeKind_WITH:
		return d.NewNodeFromTemplate(rhino.Token_WITH, templateNode)
	case NodeKind_IMPORT:
		return d.NewNodeFromTemplate(rhino.Token_IMPORT, templateNode)
	case NodeKind_EXPORT:
		return d.NewNodeFromTemplate(rhino.Token_EXPORT, templateNode)

	case NodeKind_VAR_DECLARATION:
		return d.NewNodeFromTemplate(rhino.Token_VAR, templateNode)
	case NodeKind_CONST_DECLARATION:
		return d.NewNodeFromTemplate(rhino.Token_CONST, templateNode)
	case NodeKind_LET_DECLARATION:
		return d.NewNodeFromTemplate(rhino.Token_LET, templateNode)
	case NodeKind_FUNCTION_LITERAL:
		return d.NewNodeFromTemplate(rhino.Token_FUNCTION, templateNode)
	case NodeKind_CLASS_LITERAL:
		return d.NewNodeFromTemplate(rhino.Token_CLASS, templateNode)

	case NodeKind_BLOCK:
		return d.NewNodeFromTemplate(rhino.Token_BLOCK, templateNode)
	case NodeKind_LABELED_STATEMENT:
		return d.NewNodeFromTemplate(rhino.Token_LABEL, templateNode)
	case NodeKind_LABELED_NAME:
		return d.StringNode(rhino.Token_LABEL_NAME, n, templateNode)
	case NodeKind_CLASS_MEMBERS:
		return d.NewNodeFromTemplate(rhino.Token_CLASS_MEMBERS, templateNode)
	case NodeKind_METHOD_DECLARATION:
		return d.StringNode(rhino.Token_MEMBER_FUNCTION_DEF, n, templateNode)
	case NodeKind_FIELD_DECLARATION:
		return d.StringNode(rhino.Token_MEMBER_FIELD_DEF, n, templateNode)
	case NodeKind_COMPUTED_PROP_FIELD:
		return d.NewNodeFromTemplate(rhino.Token_COMPUTED_FIELD_DEF, templateNode)
	case NodeKind_PARAMETER_LIST:
		return d.NewNodeFromTemplate(rhino.Token_PARAM_LIST, templateNode)
	case NodeKind_RENAMABLE_STRING_KEY:
		return d.StringNode(rhino.Token_STRING_KEY, n, templateNode)
	case NodeKind_QUOTED_STRING_KEY:
		quotedStringKey := d.StringNode(rhino.Token_STRING_KEY, n, templateNode)
		quotedStringKey.SetQuotedStringKey()
		return quotedStringKey
	case NodeKind_CASE:
		return d.NewNodeFromTemplate(rhino.Token_CASE, templateNode)
	case NodeKind_DEFAULT_CASE:
		return d.NewNodeFromTemplate(rhino.Token_DEFAULT_CASE, templateNode)
	case NodeKind_CATCH:
		return d.NewNodeFromTemplate(rhino.Token_CATCH, templateNode)
	case NodeKind_SUPER:
		return d.NewNodeFromTemplate(rhino.Token_SUPER, templateNode)
	case NodeKind_ARRAY_PATTERN:
		return d.NewNodeFromTemplate(rhino.Token_ARRAY_PATTERN, templateNode)
	case NodeKind_OBJECT_PATTERN:
		return d.NewNodeFromTemplate(rhino.Token_OBJECT_PATTERN, templateNode)
	case NodeKind_DESTRUCTURING_LHS:
		return d.NewNodeFromTemplate(rhino.Token_DESTRUCTURING_LHS, templateNode)
	case NodeKind_DEFAULT_VALUE:
		return d.NewNodeFromTemplate(rhino.Token_DEFAULT_VALUE, templateNode)

	case NodeKind_RENAMABLE_GETTER_DEF, NodeKind_QUOTED_GETTER_DEF:
		getterDef := d.StringNode(rhino.Token_GETTER_DEF, n, templateNode)
		if n.GetKind() == NodeKind_QUOTED_GETTER_DEF {
			getterDef.SetQuotedStringKey()
		}
		return getterDef
	case NodeKind_RENAMABLE_SETTER_DEF, NodeKind_QUOTED_SETTER_DEF:
		setterDef := d.StringNode(rhino.Token_SETTER_DEF, n, templateNode)
		if n.GetKind() == NodeKind_QUOTED_SETTER_DEF {
			setterDef.SetQuotedStringKey()
		}
		return setterDef

	case NodeKind_IMPORT_SPECS:
		return d.NewNodeFromTemplate(rhino.Token_IMPORT_SPECS, templateNode)
	case NodeKind_IMPORT_SPEC:
		return d.NewNodeFromTemplate(rhino.Token_IMPORT_SPEC, templateNode)
	case NodeKind_IMPORT_STAR:
		return d.StringNode(rhino.Token_IMPORT_STAR, n, templateNode)
	case NodeKind_EXPORT_SPECS:
		return d.NewNodeFromTemplate(rhino.Token_EXPORT_SPECS, templateNode)
	case NodeKind_EXPORT_SPEC:
		return d.NewNodeFromTemplate(rhino.Token_EXPORT_SPEC, templateNode)
	case NodeKind_MODULE_BODY:
		return d.NewNodeFromTemplate(rhino.Token_MODULE_BODY, templateNode)
	case NodeKind_ITER_REST:
		return d.NewNodeFromTemplate(rhino.Token_ITER_REST, templateNode)
	case NodeKind_ITER_SPREAD:
		return d.NewNodeFromTemplate(rhino.Token_ITER_SPREAD, templateNode)
	case NodeKind_OBJECT_REST:
		return d.NewNodeFromTemplate(rhino.Token_OBJECT_REST, templateNode)
	case NodeKind_OBJECT_SPREAD:
		return d.NewNodeFromTemplate(rhino.Token_OBJECT_SPREAD, templateNode)

	case NodeKind_NODE_KIND_UNSPECIFIED, NodeKind_UNRECOGNIZED:
		break

	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected serialized kind for AstNode: " + n.String()))
}

// FilterOutCastProp filters out any NodeProperty.COLOR_BEFORE_CASTs if no
// colors are being deserialized.
//
// This is because it doesn't make sense to have that property present on
// nodes that don't have colors.
//
// Java: ScriptNodeDeserializer.filterOutCastProp(long)
func (d *ScriptNodeDeserializer) FilterOutCastProp(nodeProperties int64) int64 {
	if d.colorPoolShard != nil {
		return nodeProperties // we are deserializing colors, so this is fine.
	}
	return nodeProperties &^ (int64(1) << rhino.NodeProperty_COLOR_FROM_CAST.GetNumber())
}
