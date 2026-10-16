// Port of com/google/javascript/jscomp/serialization/ScriptNodeDeserializer.java @ closure-compiler v20260415

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
	return newScriptNodeDeserializer_Runner(d).run()
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
func (r *scriptNodeDeserializer_Runner) run() (result *rhino.Node) {
	jutil.Try(func() {
		// The real limit is stack space.
		astStream := &protoReader{b: r.owner().scriptBytes, limit: math.MaxInt32}

		contextStack := jutil.NewArrayList[parsing.FeatureCollector_FeatureContext]()
		contextStack.Add(parsing.FeatureCollector_FeatureContext_NONE)

		astNode := &AstNode{memoizedSize: -1}
		astNode.mergeFromReader(astStream)
		astNode.frozen = true
		scriptNode := r.visit(
			astNode,
			contextStack,
			r.owner().createSourceInfoTemplate(r.owner().sourceFile))
		scriptNode.PutProp(rhino.Node_FEATURE_SET, r.featureCollector.AllFeatures())
		result = scriptNode
	}, func(ex jutil.IIOException) {
		panic(NewMalformedTypedAstExceptionWithCause(r.owner().sourceFile, ex))
	})
	return result
}

// Java: ScriptNodeDeserializer.Runner.owner()
func (r *scriptNodeDeserializer_Runner) owner() *ScriptNodeDeserializer {
	return r.owner_
}

// visit visits a single AST node and returns the corresponding Node.
// contextStack is the context stack to use for the node; sourceFileTemplate
// (may be nil) is a template node to use for the node.
//
// Java: ScriptNodeDeserializer.Runner.visit(AstNode, List<FeatureContext>, Node)
func (r *scriptNodeDeserializer_Runner) visit(
	astNode *AstNode, contextStack *jutil.ArrayList[parsing.FeatureCollector_FeatureContext], sourceFileTemplate *rhino.Node) *rhino.Node {
	if sourceFileTemplate == nil || astNode.GetSourceFile() != 0 {
		// 0 == 'not set'
		sourceFileTemplate =
			r.owner().
				createSourceInfoTemplate(r.owner().filePool[astNode.GetSourceFile()-1])
	}

	currentLine := r.previousLine + int(astNode.GetRelativeLine())
	currentColumn := r.previousColumn + int(astNode.GetRelativeColumn())

	n := r.owner().deserializeSingleNode(astNode, sourceFileTemplate)
	if astNode.HasType() && r.owner().colorPoolShard != nil {
		n.SetColor(r.owner().colorPoolShard.GetColor(int(astNode.GetType())))
	}
	properties := astNode.GetBooleanProperties()
	if properties > 0 {
		n.DeserializeProperties(
			r.owner().filterOutCastProp(properties), sourceFileTemplate.GetIsInClosureUnawareSubtree())
	}
	n.SetJSDocInfo(JSDocSerializer_DeserializeJsdoc(astNode.GetJsdoc(), r.owner().stringPool))
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
		shadowedCode := r.visit(serializedShadowChild, contextStack, closureUnawareSubtreeTemplateNode)
		r.owner().setOriginalNameIfPresent(serializedShadowChild, shadowedCode)
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
		deserializedChild := r.visit(child, contextStack, sourceFileTemplate)
		n.AddChildToBack(deserializedChild)
		r.owner().setOriginalNameIfPresent(child, deserializedChild)
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

// createSourceInfoTemplate creates a template node to use as a source of
// common attributes.
//
// This allows the prop structure to be shared among all the node from this
// source file. This reduces the cost of these properties to O(nodes) to
// O(files).
//
// Java: ScriptNodeDeserializer.createSourceInfoTemplate(SourceFile)
func (d *ScriptNodeDeserializer) createSourceInfoTemplate(file *SourceFile) *rhino.Node {
	// The Node type choice is arbitrary.
	sourceInfoTemplate := rhino.NewNode(rhino.Token_SCRIPT)
	sourceInfoTemplate.SetStaticSourceFile(file)
	return sourceInfoTemplate
}

// Java: ScriptNodeDeserializer.setOriginalNameIfPresent(AstNode, Node)
func (d *ScriptNodeDeserializer) setOriginalNameIfPresent(astNode *AstNode, n *rhino.Node) {
	if astNode.GetOriginalNamePointer() != 0 {
		n.SetOriginalNameFromStringPool(
			d.stringPool.GetInternedStrings(), int(astNode.GetOriginalNamePointer()))
	}
}

// stringNode creates a new string node with the given token & string value
// of the AstNode.
//
// Prefer calling this method over calling a regular Node.* or IR.* method
// when possible. This method integrates with RhinoStringPool to cache String
// interning results.
//
// Java: ScriptNodeDeserializer.stringNode(Token, AstNode, Node)
func (d *ScriptNodeDeserializer) stringNode(token rhino.Token, n *AstNode, templateNode *rhino.Node) *rhino.Node {
	str := rhino.Node_NewStringFromStringPool(token, d.stringPool.GetInternedStrings(), int(n.GetStringValuePointer()))
	str.ClonePropsFrom(templateNode)
	return str
}

// Java: ScriptNodeDeserializer.newNodeFromTemplate(Token, Node)
func (d *ScriptNodeDeserializer) newNodeFromTemplate(token rhino.Token, templateNode *rhino.Node) *rhino.Node {
	n := templateNode.CloneTree()
	n.SetToken(token)
	return n
}

// Java: ScriptNodeDeserializer.deserializeSingleNode(AstNode, Node)
func (d *ScriptNodeDeserializer) deserializeSingleNode(n *AstNode, templateNode *rhino.Node) *rhino.Node {
	switch n.GetKind() {
	case NodeKind_SOURCE_FILE:
		return d.newNodeFromTemplate(rhino.Token_SCRIPT, templateNode)

	case NodeKind_NUMBER_LITERAL:
		numberNode := rhino.IR_Number(n.GetDoubleValue())
		numberNode.ClonePropsFrom(templateNode)
		return numberNode
	case NodeKind_STRING_LITERAL:
		return d.stringNode(rhino.Token_STRINGLIT, n, templateNode)
	case NodeKind_IDENTIFIER:
		return d.stringNode(rhino.Token_NAME, n, templateNode)
	case NodeKind_FALSE:
		return d.newNodeFromTemplate(rhino.Token_FALSE, templateNode)
	case NodeKind_TRUE:
		return d.newNodeFromTemplate(rhino.Token_TRUE, templateNode)
	case NodeKind_NULL:
		return d.newNodeFromTemplate(rhino.Token_NULL, templateNode)
	case NodeKind_THIS:
		return d.newNodeFromTemplate(rhino.Token_THIS, templateNode)
	case NodeKind_VOID:
		return d.newNodeFromTemplate(rhino.Token_VOID, templateNode)
	case NodeKind_BIGINT_LITERAL:
		bigintString := d.stringPool.Get(int(n.GetStringValuePointer()))
		bi := rhino.IR_Bigint(jutil.NewBigIntegerFromString(bigintString))
		bi.ClonePropsFrom(templateNode)
		return bi
	case NodeKind_REGEX_LITERAL:
		return d.newNodeFromTemplate(rhino.Token_REGEXP, templateNode)
	case NodeKind_ARRAY_LITERAL:
		return d.newNodeFromTemplate(rhino.Token_ARRAYLIT, templateNode)
	case NodeKind_OBJECT_LITERAL:
		return d.newNodeFromTemplate(rhino.Token_OBJECTLIT, templateNode)

	case NodeKind_ASSIGNMENT:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN, templateNode)
	case NodeKind_CALL:
		return d.newNodeFromTemplate(rhino.Token_CALL, templateNode)
	case NodeKind_NEW:
		return d.newNodeFromTemplate(rhino.Token_NEW, templateNode)
	case NodeKind_PROPERTY_ACCESS:
		return d.stringNode(rhino.Token_GETPROP, n, templateNode)
	case NodeKind_ELEMENT_ACCESS:
		return d.newNodeFromTemplate(rhino.Token_GETELEM, templateNode)

	case NodeKind_COMMA:
		return d.newNodeFromTemplate(rhino.Token_COMMA, templateNode)
	case NodeKind_BOOLEAN_OR:
		return d.newNodeFromTemplate(rhino.Token_OR, templateNode)
	case NodeKind_BOOLEAN_AND:
		return d.newNodeFromTemplate(rhino.Token_AND, templateNode)
	case NodeKind_HOOK:
		return d.newNodeFromTemplate(rhino.Token_HOOK, templateNode)
	case NodeKind_EQUAL:
		return d.newNodeFromTemplate(rhino.Token_EQ, templateNode)
	case NodeKind_NOT_EQUAL:
		return d.newNodeFromTemplate(rhino.Token_NE, templateNode)
	case NodeKind_LESS_THAN:
		return d.newNodeFromTemplate(rhino.Token_LT, templateNode)
	case NodeKind_LESS_THAN_EQUAL:
		return d.newNodeFromTemplate(rhino.Token_LE, templateNode)
	case NodeKind_GREATER_THAN:
		return d.newNodeFromTemplate(rhino.Token_GT, templateNode)
	case NodeKind_GREATER_THAN_EQUAL:
		return d.newNodeFromTemplate(rhino.Token_GE, templateNode)
	case NodeKind_TRIPLE_EQUAL:
		return d.newNodeFromTemplate(rhino.Token_SHEQ, templateNode)
	case NodeKind_NOT_TRIPLE_EQUAL:
		return d.newNodeFromTemplate(rhino.Token_SHNE, templateNode)
	case NodeKind_NOT:
		return d.newNodeFromTemplate(rhino.Token_NOT, templateNode)
	case NodeKind_POSITIVE:
		return d.newNodeFromTemplate(rhino.Token_POS, templateNode)
	case NodeKind_NEGATIVE:
		return d.newNodeFromTemplate(rhino.Token_NEG, templateNode)
	case NodeKind_TYPEOF:
		return d.newNodeFromTemplate(rhino.Token_TYPEOF, templateNode)
	case NodeKind_INSTANCEOF:
		return d.newNodeFromTemplate(rhino.Token_INSTANCEOF, templateNode)
	case NodeKind_IN:
		return d.newNodeFromTemplate(rhino.Token_IN, templateNode)

	case NodeKind_ADD:
		return d.newNodeFromTemplate(rhino.Token_ADD, templateNode)
	case NodeKind_SUBTRACT:
		return d.newNodeFromTemplate(rhino.Token_SUB, templateNode)
	case NodeKind_MULTIPLY:
		return d.newNodeFromTemplate(rhino.Token_MUL, templateNode)
	case NodeKind_DIVIDE:
		return d.newNodeFromTemplate(rhino.Token_DIV, templateNode)
	case NodeKind_MODULO:
		return d.newNodeFromTemplate(rhino.Token_MOD, templateNode)
	case NodeKind_EXPONENT:
		return d.newNodeFromTemplate(rhino.Token_EXPONENT, templateNode)
	case NodeKind_BITWISE_NOT:
		return d.newNodeFromTemplate(rhino.Token_BITNOT, templateNode)
	case NodeKind_BITWISE_OR:
		return d.newNodeFromTemplate(rhino.Token_BITOR, templateNode)
	case NodeKind_BITWISE_AND:
		return d.newNodeFromTemplate(rhino.Token_BITAND, templateNode)
	case NodeKind_BITWISE_XOR:
		return d.newNodeFromTemplate(rhino.Token_BITXOR, templateNode)
	case NodeKind_LEFT_SHIFT:
		return d.newNodeFromTemplate(rhino.Token_LSH, templateNode)
	case NodeKind_RIGHT_SHIFT:
		return d.newNodeFromTemplate(rhino.Token_RSH, templateNode)
	case NodeKind_UNSIGNED_RIGHT_SHIFT:
		return d.newNodeFromTemplate(rhino.Token_URSH, templateNode)
	case NodeKind_PRE_INCREMENT:
		return d.newNodeFromTemplate(rhino.Token_INC, templateNode)
	case NodeKind_POST_INCREMENT:
		postInc := d.newNodeFromTemplate(rhino.Token_INC, templateNode)
		postInc.PutBooleanProp(rhino.Node_INCRDECR_PROP, true)
		return postInc
	case NodeKind_PRE_DECREMENT:
		return d.newNodeFromTemplate(rhino.Token_DEC, templateNode)
	case NodeKind_POST_DECREMENT:
		postDec := d.newNodeFromTemplate(rhino.Token_DEC, templateNode)
		postDec.PutBooleanProp(rhino.Node_INCRDECR_PROP, true)
		return postDec

	case NodeKind_ASSIGN_ADD:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_ADD, templateNode)
	case NodeKind_ASSIGN_SUBTRACT:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_SUB, templateNode)
	case NodeKind_ASSIGN_MULTIPLY:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_MUL, templateNode)
	case NodeKind_ASSIGN_DIVIDE:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_DIV, templateNode)
	case NodeKind_ASSIGN_MODULO:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_MOD, templateNode)
	case NodeKind_ASSIGN_EXPONENT:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_EXPONENT, templateNode)
	case NodeKind_ASSIGN_BITWISE_OR:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_BITOR, templateNode)
	case NodeKind_ASSIGN_BITWISE_AND:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_BITAND, templateNode)
	case NodeKind_ASSIGN_BITWISE_XOR:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_BITXOR, templateNode)
	case NodeKind_ASSIGN_LEFT_SHIFT:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_LSH, templateNode)
	case NodeKind_ASSIGN_RIGHT_SHIFT:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_RSH, templateNode)
	case NodeKind_ASSIGN_UNSIGNED_RIGHT_SHIFT:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_URSH, templateNode)

	case NodeKind_YIELD:
		return d.newNodeFromTemplate(rhino.Token_YIELD, templateNode)
	case NodeKind_AWAIT:
		return d.newNodeFromTemplate(rhino.Token_AWAIT, templateNode)
	case NodeKind_DELETE:
		return d.newNodeFromTemplate(rhino.Token_DELPROP, templateNode)
	case NodeKind_TAGGED_TEMPLATELIT:
		return d.newNodeFromTemplate(rhino.Token_TAGGED_TEMPLATELIT, templateNode)
	case NodeKind_TEMPLATELIT:
		return d.newNodeFromTemplate(rhino.Token_TEMPLATELIT, templateNode)
	case NodeKind_TEMPLATELIT_SUB:
		return d.newNodeFromTemplate(rhino.Token_TEMPLATELIT_SUB, templateNode)
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
		return d.newNodeFromTemplate(rhino.Token_NEW_TARGET, templateNode)
	case NodeKind_COMPUTED_PROP:
		return d.newNodeFromTemplate(rhino.Token_COMPUTED_PROP, templateNode)
	case NodeKind_IMPORT_META:
		return d.newNodeFromTemplate(rhino.Token_IMPORT_META, templateNode)
	case NodeKind_OPTCHAIN_PROPERTY_ACCESS:
		return d.stringNode(rhino.Token_OPTCHAIN_GETPROP, n, templateNode)
	case NodeKind_OPTCHAIN_CALL:
		return d.newNodeFromTemplate(rhino.Token_OPTCHAIN_CALL, templateNode)
	case NodeKind_OPTCHAIN_ELEMENT_ACCESS:
		return d.newNodeFromTemplate(rhino.Token_OPTCHAIN_GETELEM, templateNode)
	case NodeKind_COALESCE:
		return d.newNodeFromTemplate(rhino.Token_COALESCE, templateNode)
	case NodeKind_DYNAMIC_IMPORT:
		return d.newNodeFromTemplate(rhino.Token_DYNAMIC_IMPORT, templateNode)

	case NodeKind_ASSIGN_OR:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_OR, templateNode)
	case NodeKind_ASSIGN_AND:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_AND, templateNode)
	case NodeKind_ASSIGN_COALESCE:
		return d.newNodeFromTemplate(rhino.Token_ASSIGN_COALESCE, templateNode)

	case NodeKind_EXPRESSION_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_EXPR_RESULT, templateNode)
	case NodeKind_BREAK_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_BREAK, templateNode)
	case NodeKind_CONTINUE_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_CONTINUE, templateNode)
	case NodeKind_DEBUGGER_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_DEBUGGER, templateNode)
	case NodeKind_DO_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_DO, templateNode)
	case NodeKind_FOR_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_FOR, templateNode)
	case NodeKind_FOR_IN_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_FOR_IN, templateNode)
	case NodeKind_FOR_OF_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_FOR_OF, templateNode)
	case NodeKind_FOR_AWAIT_OF_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_FOR_AWAIT_OF, templateNode)
	case NodeKind_IF_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_IF, templateNode)
	case NodeKind_RETURN_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_RETURN, templateNode)
	case NodeKind_SWITCH_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_SWITCH, templateNode)
	case NodeKind_SWITCH_BODY:
		return d.newNodeFromTemplate(rhino.Token_SWITCH_BODY, templateNode)
	case NodeKind_THROW_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_THROW, templateNode)
	case NodeKind_TRY_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_TRY, templateNode)
	case NodeKind_WHILE_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_WHILE, templateNode)
	case NodeKind_EMPTY:
		return d.newNodeFromTemplate(rhino.Token_EMPTY, templateNode)
	case NodeKind_WITH:
		return d.newNodeFromTemplate(rhino.Token_WITH, templateNode)
	case NodeKind_IMPORT:
		return d.newNodeFromTemplate(rhino.Token_IMPORT, templateNode)
	case NodeKind_EXPORT:
		return d.newNodeFromTemplate(rhino.Token_EXPORT, templateNode)

	case NodeKind_VAR_DECLARATION:
		return d.newNodeFromTemplate(rhino.Token_VAR, templateNode)
	case NodeKind_CONST_DECLARATION:
		return d.newNodeFromTemplate(rhino.Token_CONST, templateNode)
	case NodeKind_LET_DECLARATION:
		return d.newNodeFromTemplate(rhino.Token_LET, templateNode)
	case NodeKind_FUNCTION_LITERAL:
		return d.newNodeFromTemplate(rhino.Token_FUNCTION, templateNode)
	case NodeKind_CLASS_LITERAL:
		return d.newNodeFromTemplate(rhino.Token_CLASS, templateNode)

	case NodeKind_BLOCK:
		return d.newNodeFromTemplate(rhino.Token_BLOCK, templateNode)
	case NodeKind_LABELED_STATEMENT:
		return d.newNodeFromTemplate(rhino.Token_LABEL, templateNode)
	case NodeKind_LABELED_NAME:
		return d.stringNode(rhino.Token_LABEL_NAME, n, templateNode)
	case NodeKind_CLASS_MEMBERS:
		return d.newNodeFromTemplate(rhino.Token_CLASS_MEMBERS, templateNode)
	case NodeKind_METHOD_DECLARATION:
		return d.stringNode(rhino.Token_MEMBER_FUNCTION_DEF, n, templateNode)
	case NodeKind_FIELD_DECLARATION:
		return d.stringNode(rhino.Token_MEMBER_FIELD_DEF, n, templateNode)
	case NodeKind_COMPUTED_PROP_FIELD:
		return d.newNodeFromTemplate(rhino.Token_COMPUTED_FIELD_DEF, templateNode)
	case NodeKind_PARAMETER_LIST:
		return d.newNodeFromTemplate(rhino.Token_PARAM_LIST, templateNode)
	case NodeKind_RENAMABLE_STRING_KEY:
		return d.stringNode(rhino.Token_STRING_KEY, n, templateNode)
	case NodeKind_QUOTED_STRING_KEY:
		quotedStringKey := d.stringNode(rhino.Token_STRING_KEY, n, templateNode)
		quotedStringKey.SetQuotedStringKey()
		return quotedStringKey
	case NodeKind_CASE:
		return d.newNodeFromTemplate(rhino.Token_CASE, templateNode)
	case NodeKind_DEFAULT_CASE:
		return d.newNodeFromTemplate(rhino.Token_DEFAULT_CASE, templateNode)
	case NodeKind_CATCH:
		return d.newNodeFromTemplate(rhino.Token_CATCH, templateNode)
	case NodeKind_SUPER:
		return d.newNodeFromTemplate(rhino.Token_SUPER, templateNode)
	case NodeKind_ARRAY_PATTERN:
		return d.newNodeFromTemplate(rhino.Token_ARRAY_PATTERN, templateNode)
	case NodeKind_OBJECT_PATTERN:
		return d.newNodeFromTemplate(rhino.Token_OBJECT_PATTERN, templateNode)
	case NodeKind_DESTRUCTURING_LHS:
		return d.newNodeFromTemplate(rhino.Token_DESTRUCTURING_LHS, templateNode)
	case NodeKind_DEFAULT_VALUE:
		return d.newNodeFromTemplate(rhino.Token_DEFAULT_VALUE, templateNode)

	case NodeKind_RENAMABLE_GETTER_DEF, NodeKind_QUOTED_GETTER_DEF:
		getterDef := d.stringNode(rhino.Token_GETTER_DEF, n, templateNode)
		if n.GetKind() == NodeKind_QUOTED_GETTER_DEF {
			getterDef.SetQuotedStringKey()
		}
		return getterDef
	case NodeKind_RENAMABLE_SETTER_DEF, NodeKind_QUOTED_SETTER_DEF:
		setterDef := d.stringNode(rhino.Token_SETTER_DEF, n, templateNode)
		if n.GetKind() == NodeKind_QUOTED_SETTER_DEF {
			setterDef.SetQuotedStringKey()
		}
		return setterDef

	case NodeKind_IMPORT_SPECS:
		return d.newNodeFromTemplate(rhino.Token_IMPORT_SPECS, templateNode)
	case NodeKind_IMPORT_SPEC:
		return d.newNodeFromTemplate(rhino.Token_IMPORT_SPEC, templateNode)
	case NodeKind_IMPORT_STAR:
		return d.stringNode(rhino.Token_IMPORT_STAR, n, templateNode)
	case NodeKind_EXPORT_SPECS:
		return d.newNodeFromTemplate(rhino.Token_EXPORT_SPECS, templateNode)
	case NodeKind_EXPORT_SPEC:
		return d.newNodeFromTemplate(rhino.Token_EXPORT_SPEC, templateNode)
	case NodeKind_MODULE_BODY:
		return d.newNodeFromTemplate(rhino.Token_MODULE_BODY, templateNode)
	case NodeKind_ITER_REST:
		return d.newNodeFromTemplate(rhino.Token_ITER_REST, templateNode)
	case NodeKind_ITER_SPREAD:
		return d.newNodeFromTemplate(rhino.Token_ITER_SPREAD, templateNode)
	case NodeKind_OBJECT_REST:
		return d.newNodeFromTemplate(rhino.Token_OBJECT_REST, templateNode)
	case NodeKind_OBJECT_SPREAD:
		return d.newNodeFromTemplate(rhino.Token_OBJECT_SPREAD, templateNode)

	case NodeKind_NODE_KIND_UNSPECIFIED, NodeKind_UNRECOGNIZED:
		break

	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected serialized kind for AstNode: " + n.String()))
}

// filterOutCastProp filters out any NodeProperty.COLOR_BEFORE_CASTs if no
// colors are being deserialized.
//
// This is because it doesn't make sense to have that property present on
// nodes that don't have colors.
//
// Java: ScriptNodeDeserializer.filterOutCastProp(long)
func (d *ScriptNodeDeserializer) filterOutCastProp(nodeProperties int64) int64 {
	if d.colorPoolShard != nil {
		return nodeProperties // we are deserializing colors, so this is fine.
	}
	return nodeProperties &^ (int64(1) << rhino.NodeProperty_COLOR_FROM_CAST.GetNumber())
}
