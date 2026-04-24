// Port of com/google/javascript/jscomp/parsing/IRFactory.java
// @ closure-compiler v20260415
//
// IRFactory transforms ParseTree nodes (from internal/parser/trees) into
// rhino IR Nodes (from internal/rhino). Only the subset of transforms
// needed for the current corpus tier is implemented; stub methods are
// marked with a TODO comment.

package parsing

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/parser/trees"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// IRFactory transforms the external parse-tree AST (trees.ParseTree) into
// the internal IR AST (rhino.Node). Mirrors
// com.google.javascript.jscomp.parsing.IRFactory.
type IRFactory struct {
	sourceFile    rhino.StaticSourceFile
	errorReporter errorReporter
}

// errorReporter is a minimal error-reporting interface used by IRFactory.
// Mirrors com.google.javascript.rhino.ErrorReporter.
type errorReporter interface {
	// Error reports a hard error at the given source position.
	Error(msg, sourceName string, lineno, charno int)
	// Warning reports a warning at the given source position.
	Warning(msg, sourceName string, lineno, charno int)
}

// NewIRFactory constructs an IRFactory for the given source file.
func NewIRFactory(sourceFile rhino.StaticSourceFile, er errorReporter) *IRFactory {
	return &IRFactory{
		sourceFile:    sourceFile,
		errorReporter: er,
	}
}

// Transform converts a ParseTree node into a rhino IR Node.
// Mirrors IRFactory.transform(ParseTree) in Java.
func (f *IRFactory) Transform(tree *trees.ParseTree) *rhino.Node {
	if tree == nil {
		return nil
	}
	node := f.dispatch(tree)
	if node != nil {
		f.setSourceInfo(node, tree)
	}
	return node
}

// dispatch routes a ParseTree to the appropriate transform* method based on
// its type. Mirrors the TransformDispatcher.process(ParseTree) switch in Java.
func (f *IRFactory) dispatch(tree *trees.ParseTree) *rhino.Node {
	switch tree.Type {
	// ---- programme / block -----------------------------------------
	case trees.ParseTreeTypeProgram:
		return f.transformProgram((*trees.ProgramTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBlock:
		return f.transformBlock((*trees.BlockTree)(unsafe.Pointer(tree)))
	// NOTE: ParseTreeTypeExpressionStatement is intentionally omitted here.
	// ExpressionStatementTree does not embed ParseTree, so it cannot be
	// addressed via *ParseTree from the parse-tree slice. The dispatch will be
	// wired when the parser produces a proper embedded tree.

	// ---- primitive expressions -------------------------------------
	case trees.ParseTreeTypeIdentifierExpression:
		return f.transformIdentifierExpression((*trees.IdentifierExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeLiteralExpression:
		return f.transformLiteralExpression((*trees.LiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeThisExpression:
		return f.transformThisExpression()
	case trees.ParseTreeTypeSuperExpression:
		return f.transformSuperExpression()

	// ---- unary / binary / update -----------------------------------
	case trees.ParseTreeTypeUnaryExpression:
		return f.transformUnaryExpression((*trees.UnaryExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeUpdateExpression:
		return f.transformUpdateExpression((*trees.UpdateExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBinaryOperator:
		return f.transformBinaryOperator((*trees.BinaryOperatorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeConditionalExpression:
		return f.transformConditionalExpression((*trees.ConditionalExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCommaExpression:
		return f.transformCommaExpression((*trees.CommaExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeParenExpression:
		return f.transformParenExpression((*trees.ParenExpressionTree)(unsafe.Pointer(tree)))

	// ---- calls / member access / new -------------------------------
	case trees.ParseTreeTypeCallExpression:
		return f.transformCallExpression((*trees.CallExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeNewExpression:
		return f.transformNewExpression((*trees.NewExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeMemberExpression:
		return f.transformMemberExpression((*trees.MemberExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeMemberLookupExpression:
		return f.transformMemberLookupExpression((*trees.MemberLookupExpressionTree)(unsafe.Pointer(tree)))

	// ---- structured literals ---------------------------------------
	case trees.ParseTreeTypeArrayLiteralExpression:
		return f.transformArrayLiteralExpression((*trees.ArrayLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeObjectLiteralExpression:
		return f.transformObjectLiteralExpression((*trees.ObjectLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypePropertyNameAssignment:
		return f.transformPropertyNameAssignment((*trees.PropertyNameAssignmentTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeGetAccessor:
		return f.transformGetAccessor((*trees.GetAccessorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeSetAccessor:
		return f.transformSetAccessor((*trees.SetAccessorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateLiteralExpression:
		return f.transformTemplateLiteralExpression((*trees.TemplateLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateLiteralPortion:
		return f.transformTemplateLiteralPortion((*trees.TemplateLiteralPortionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateSubstitution:
		return f.transformTemplateSubstitution((*trees.TemplateSubstitutionTree)(unsafe.Pointer(tree)))

	// ---- null / empty placeholder ----------------------------------
	case trees.ParseTreeTypeNull:
		return rhino.New(rhino.EMPTY)

	default:
		return rhino.New(rhino.EMPTY)
	}
}

// setSourceInfo copies source-position information from tree to node.
// Mirrors IRFactory.setSourceInfo(Node, ParseTree) in Java.
func (f *IRFactory) setSourceInfo(node *rhino.Node, tree *trees.ParseTree) {
	if node.Lineno() == 0 {
		// Location lines in ParseTree are 0-based; rhino expects 1-based.
		node.SetPos(tree.Location.Start.Line+1, tree.Location.Start.Column)
	}
}

// newNode creates a new rhino Node with the given token. Mirrors
// IRFactory.newNode(Token) in Java.
func (f *IRFactory) newNode(tok rhino.Token) *rhino.Node {
	return rhino.New(tok)
}

// newStringNode creates a new rhino string Node. Mirrors
// IRFactory.newStringNode(Token, String) in Java.
func (f *IRFactory) newStringNode(tok rhino.Token, s string) *rhino.Node {
	return rhino.NewString(tok, s)
}

// newNumberNode creates a new rhino number Node. Mirrors
// IRFactory.newNumberNode(double) in Java.
func (f *IRFactory) newNumberNode(v float64) *rhino.Node {
	return rhino.NewNumber(v)
}

// ---- programme / block ---------------------------------------------

// transformProgram mirrors TransformDispatcher.processAstRoot(ProgramTree).
func (f *IRFactory) transformProgram(tree *trees.ProgramTree) *rhino.Node {
	script := f.newNode(rhino.SCRIPT)
	for _, child := range tree.SourceElements {
		script.AddChildToBack(f.Transform(child))
	}
	return script
}

// transformBlock mirrors TransformDispatcher.processBlock(BlockTree).
func (f *IRFactory) transformBlock(tree *trees.BlockTree) *rhino.Node {
	node := f.newNode(rhino.BLOCK)
	for _, child := range tree.Statements {
		node.AddChildToBack(f.Transform(child))
	}
	return node
}

// transformExpressionStatementExpr mirrors
// TransformDispatcher.processExpressionStatement, given the inner expression
// tree. Used when the expression-statement's child is dispatched directly.
func (f *IRFactory) transformExpressionStatementExpr(expr *trees.ParseTree) *rhino.Node {
	node := f.newNode(rhino.EXPR_RESULT)
	node.AddChildToBack(f.Transform(expr))
	return node
}

// ---- primitive expressions -----------------------------------------

// transformIdentifierExpression mirrors
// TransformDispatcher.processName(IdentifierExpressionTree, IdentifierType).
func (f *IRFactory) transformIdentifierExpression(tree *trees.IdentifierExpressionTree) *rhino.Node {
	if tree.IdentifierToken == nil {
		return f.newStringNode(rhino.NAME, "__missing_name__")
	}
	return f.newStringNode(rhino.NAME, tree.IdentifierToken.Value)
}

// transformLiteralExpression dispatches on the literal token type.
// Mirrors the per-type literal processing in TransformDispatcher.
func (f *IRFactory) transformLiteralExpression(tree *trees.LiteralExpressionTree) *rhino.Node {
	if tree.LiteralToken == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := tree.LiteralToken
	switch tok.Type {
	case trees.TokenTypeTrue:
		return f.newNode(rhino.TRUE)
	case trees.TokenTypeFalse:
		return f.newNode(rhino.FALSE)
	case trees.TokenTypeNull:
		return f.newNode(rhino.NULL)
	case trees.TokenTypeNumber:
		// trees.Token has no Value field; number value is not accessible via the
		// stub. Return 0 as a placeholder until the full token is wired.
		return f.newNumberNode(0)
	case trees.TokenTypeString:
		return f.newStringNode(rhino.STRINGLIT, "")
	case trees.TokenTypeRegularExpression:
		return f.transformRegexpLiteral("")
	default:
		return f.newNode(rhino.EMPTY)
	}
}

// transformThisExpression mirrors
// TransformDispatcher.processThisExpression(ThisExpressionTree).
func (f *IRFactory) transformThisExpression() *rhino.Node {
	return f.newNode(rhino.THIS)
}

// transformSuperExpression mirrors Java super expression → Token.SUPER.
func (f *IRFactory) transformSuperExpression() *rhino.Node {
	return f.newNode(rhino.SUPER)
}

// ---- unary / binary / update ---------------------------------------

// transformUnaryExpression mirrors
// TransformDispatcher.processUnaryExpression(UnaryExpressionTree).
func (f *IRFactory) transformUnaryExpression(tree *trees.UnaryExpressionTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformUnaryTokenType(tree.Operator.Type)
	operand := f.Transform(tree.Operand)
	node := f.newNode(tok)
	node.AddChildToBack(operand)
	return node
}

// transformUpdateExpression mirrors
// TransformDispatcher.processUpdateExpression(UpdateExpressionTree).
func (f *IRFactory) transformUpdateExpression(tree *trees.UpdateExpressionTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformUpdateTokenType(tree.Operator.Type)
	operand := f.Transform(tree.Operand)
	node := f.newNode(tok)
	node.AddChildToBack(operand)
	postfix := tree.OperatorPosition == trees.OperatorPositionPostfix
	node.PutBooleanProp(rhino.PropIncrDecr, postfix)
	return node
}

// transformBinaryOperator mirrors
// TransformDispatcher.processBinaryExpression(BinaryOperatorTree).
func (f *IRFactory) transformBinaryOperator(tree *trees.BinaryOperatorTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformBinaryTokenType(tree.Operator.Type)
	left := f.Transform(tree.Left)
	right := f.Transform(tree.Right)
	node := f.newNode(tok)
	node.AddChildToBack(left)
	node.AddChildToBack(right)
	return node
}

// transformConditionalExpression mirrors
// TransformDispatcher.processConditionalExpression(ConditionalExpressionTree).
func (f *IRFactory) transformConditionalExpression(tree *trees.ConditionalExpressionTree) *rhino.Node {
	node := f.newNode(rhino.HOOK)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.Transform(tree.Left))
	node.AddChildToBack(f.Transform(tree.Right))
	return node
}

// transformCommaExpression mirrors
// TransformDispatcher.processCommaExpression(CommaExpressionTree).
func (f *IRFactory) transformCommaExpression(tree *trees.CommaExpressionTree) *rhino.Node {
	if len(tree.Expressions) == 0 {
		return f.newNode(rhino.EMPTY)
	}
	// Java builds a left-associative chain of COMMA nodes.
	root := f.newNode(rhino.COMMA)
	for i, expr := range tree.Expressions {
		if i < 2 {
			root.AddChildToBack(f.Transform(expr))
		} else {
			parent := f.newNode(rhino.COMMA)
			parent.AddChildToBack(root)
			parent.AddChildToBack(f.Transform(expr))
			root = parent
		}
	}
	return root
}

// transformParenExpression mirrors
// TransformDispatcher.processParenthesizedExpression(ParenExpressionTree).
func (f *IRFactory) transformParenExpression(tree *trees.ParenExpressionTree) *rhino.Node {
	expr := f.Transform(tree.Expression)
	if expr != nil {
		expr.SetIsParenthesized(true)
	}
	return expr
}

// ---- calls / member access / new -----------------------------------

// transformCallExpression mirrors
// TransformDispatcher.processFunctionCall(CallExpressionTree).
func (f *IRFactory) transformCallExpression(tree *trees.CallExpressionTree) *rhino.Node {
	node := f.newNode(rhino.CALL)
	node.AddChildToBack(f.Transform(tree.Operand))
	if tree.Arguments != nil {
		for _, arg := range tree.Arguments.Arguments {
			node.AddChildToBack(f.Transform(arg))
		}
	}
	// Free-call annotation: if callee is not a GETPROP/GETELEM, mark FREE_CALL.
	callee := node.FirstChild()
	if callee != nil && !callee.IsGetProp() && !callee.IsGetElem() {
		node.PutBooleanProp(rhino.PropFreeCall, true)
	}
	return node
}

// transformNewExpression mirrors
// TransformDispatcher.processNewExpression(NewExpressionTree).
func (f *IRFactory) transformNewExpression(tree *trees.NewExpressionTree) *rhino.Node {
	node := f.newNode(rhino.NEW)
	node.AddChildToBack(f.Transform(tree.Operand))
	if tree.Arguments != nil {
		for _, arg := range tree.Arguments.Arguments {
			node.AddChildToBack(f.Transform(arg))
		}
	}
	return node
}

// transformMemberExpression mirrors
// TransformDispatcher.processPropertyGet(MemberExpressionTree).
func (f *IRFactory) transformMemberExpression(tree *trees.MemberExpressionTree) *rhino.Node {
	leftChild := f.Transform(tree.Operand)
	if tree.MemberName == nil {
		return leftChild
	}
	getProp := f.newStringNode(rhino.GETPROP, tree.MemberName.Value)
	getProp.AddChildToBack(leftChild)
	return getProp
}

// transformMemberLookupExpression mirrors
// TransformDispatcher.processElementGet(MemberLookupExpressionTree).
func (f *IRFactory) transformMemberLookupExpression(tree *trees.MemberLookupExpressionTree) *rhino.Node {
	node := f.newNode(rhino.GETELEM)
	node.AddChildToBack(f.Transform(tree.Operand))
	node.AddChildToBack(f.Transform(tree.MemberExpression))
	return node
}

// ---- structured literals -------------------------------------------

// transformArrayLiteralExpression mirrors
// TransformDispatcher.processArrayLiteral(ArrayLiteralExpressionTree).
func (f *IRFactory) transformArrayLiteralExpression(tree *trees.ArrayLiteralExpressionTree) *rhino.Node {
	node := f.newNode(rhino.ARRAYLIT)
	for _, el := range tree.Elements {
		node.AddChildToBack(f.Transform(el))
	}
	return node
}

// transformObjectLiteralExpression mirrors
// TransformDispatcher.processObjectLiteral(ObjectLiteralExpressionTree).
func (f *IRFactory) transformObjectLiteralExpression(tree *trees.ObjectLiteralExpressionTree) *rhino.Node {
	node := f.newNode(rhino.OBJECTLIT)
	for _, el := range tree.PropertyNameAndValues {
		key := f.Transform(el)
		if key != nil {
			node.AddChildToBack(key)
		}
	}
	return node
}

// transformPropertyNameAssignment mirrors
// TransformDispatcher.processPropertyNameAssignment(PropertyNameAssignmentTree).
func (f *IRFactory) transformPropertyNameAssignment(tree *trees.PropertyNameAssignmentTree) *rhino.Node {
	key := f.processObjectLitKey(tree.Name, rhino.STRING_KEY)
	if tree.Value != nil {
		key.AddChildToFront(f.Transform(tree.Value))
	} else if key.Token() == rhino.STRING_KEY {
		// Shorthand property: value is a NAME node with the same string.
		val := f.newStringNode(rhino.NAME, key.String_())
		key.AddChildToFront(val)
		key.PutBooleanProp(rhino.PropShorthandProperty, true)
	}
	return key
}

// transformGetAccessor mirrors
// TransformDispatcher.processGetAccessor(GetAccessorTree).
func (f *IRFactory) transformGetAccessor(tree *trees.GetAccessorTree) *rhino.Node {
	key := f.processObjectLitKey(tree.PropertyName, rhino.GETTER_DEF)
	body := f.Transform(&tree.Body.ParseTree)
	dummyName := f.newStringNode(rhino.NAME, "")
	f.setSourceInfo(dummyName, &tree.Body.ParseTree)
	paramList := f.newNode(rhino.PARAM_LIST)
	f.setSourceInfo(paramList, &tree.Body.ParseTree)
	fn := f.newNode(rhino.FUNCTION)
	fn.AddChildToBack(dummyName)
	fn.AddChildToBack(paramList)
	fn.AddChildToBack(body)
	f.setSourceInfo(fn, &tree.Body.ParseTree)
	key.AddChildToFront(fn)
	if tree.IsStatic {
		key.PutBooleanProp(rhino.PropStaticMember, true)
	}
	return key
}

// transformSetAccessor mirrors
// TransformDispatcher.processSetAccessor(SetAccessorTree).
func (f *IRFactory) transformSetAccessor(tree *trees.SetAccessorTree) *rhino.Node {
	key := f.processObjectLitKey(tree.PropertyName, rhino.SETTER_DEF)
	// Build a minimal PARAM_LIST from the formal parameter.
	paramList := f.newNode(rhino.PARAM_LIST)
	if tree.Parameter != nil {
		for _, p := range tree.Parameter.Parameters {
			paramList.AddChildToBack(f.Transform(p))
		}
		f.setSourceInfo(paramList, &tree.Parameter.ParseTree)
	}
	body := f.Transform(&tree.Body.ParseTree)
	dummyName := f.newStringNode(rhino.NAME, "")
	if tree.PropertyName != nil {
		f.setSourceInfoFromToken(dummyName, tree.PropertyName)
	}
	fn := f.newNode(rhino.FUNCTION)
	fn.AddChildToBack(dummyName)
	fn.AddChildToBack(paramList)
	fn.AddChildToBack(body)
	f.setSourceInfo(fn, &tree.Body.ParseTree)
	key.AddChildToFront(fn)
	if tree.IsStatic {
		key.PutBooleanProp(rhino.PropStaticMember, true)
	}
	return key
}

// transformTemplateLiteralExpression mirrors
// TransformDispatcher.processTemplateLiteral(TemplateLiteralExpressionTree).
func (f *IRFactory) transformTemplateLiteralExpression(tree *trees.TemplateLiteralExpressionTree) *rhino.Node {
	templateLitNode := f.newNode(rhino.TEMPLATELIT)
	f.setSourceInfo(templateLitNode, &tree.ParseTree)
	var node *rhino.Node
	if tree.Operand == nil {
		node = templateLitNode
	} else {
		node = f.newNode(rhino.TAGGED_TEMPLATELIT)
		node.AddChildToBack(f.Transform(tree.Operand))
		node.AddChildToBack(templateLitNode)
	}
	for _, child := range tree.Elements {
		templateLitNode.AddChildToBack(f.Transform(child))
	}
	if node.Token() == rhino.TAGGED_TEMPLATELIT {
		// Free-call annotation for tagged template literals.
		callee := node.FirstChild()
		if callee != nil && !callee.IsGetProp() && !callee.IsGetElem() {
			node.PutBooleanProp(rhino.PropFreeCall, true)
		}
	}
	return node
}

// transformTemplateLiteralPortion mirrors
// TransformDispatcher.processTemplateLiteralPortion(TemplateLiteralPortionTree).
// Note: trees.Token has no Value field in the stub; cooked/raw are empty stubs.
func (f *IRFactory) transformTemplateLiteralPortion(tree *trees.TemplateLiteralPortionTree) *rhino.Node {
	return rhino.NewTemplateLitString("", "")
}

// transformTemplateSubstitution mirrors
// TransformDispatcher.processTemplateSubstitution(TemplateSubstitutionTree).
func (f *IRFactory) transformTemplateSubstitution(tree *trees.TemplateSubstitutionTree) *rhino.Node {
	node := f.newNode(rhino.TEMPLATELIT_SUB)
	node.AddChildToBack(f.Transform(tree.Expression))
	return node
}

// ---- helpers -------------------------------------------------------

// processObjectLitKey creates a STRING_KEY / GETTER_DEF / SETTER_DEF node
// from the property token. Mirrors IRFactory.TransformDispatcher.processObjectLitKey.
// Note: trees.Token is a stub without a Value field; the property name
// is not accessible through the current trees.Token type.
func (f *IRFactory) processObjectLitKey(tok *trees.Token, output rhino.Token) *rhino.Node {
	if tok == nil {
		return f.newStringNode(output, "__missing_key__")
	}
	switch tok.Type {
	case trees.TokenTypeIdentifier:
		// Identifier key: name carried by IdentifierToken.Value (not on trees.Token stub).
		node := f.newStringNode(output, "")
		f.setSourceInfoFromToken(node, tok)
		return node
	case trees.TokenTypeNumber:
		node := f.newStringNode(output, "0")
		node.PutBooleanProp(rhino.PropQuotedStringKey, true)
		f.setSourceInfoFromToken(node, tok)
		return node
	default:
		node := f.newStringNode(output, "")
		node.PutBooleanProp(rhino.PropQuotedStringKey, true)
		f.setSourceInfoFromToken(node, tok)
		return node
	}
}

// transformRegexpLiteral builds a REGEXP node from the raw literal string.
// Mirrors TransformDispatcher.processRegExpLiteral.
func (f *IRFactory) transformRegexpLiteral(raw string) *rhino.Node {
	lastSlash := len(raw) - 1
	for lastSlash > 0 && raw[lastSlash] != '/' {
		lastSlash--
	}
	pattern := raw[1:lastSlash]
	flags := raw[lastSlash+1:]
	patternNode := f.newStringNode(rhino.STRINGLIT, pattern)
	node := f.newNode(rhino.REGEXP)
	node.AddChildToBack(patternNode)
	if flags != "" {
		flagsNode := f.newStringNode(rhino.STRINGLIT, flags)
		node.AddChildToBack(flagsNode)
	}
	return node
}

// setSourceInfoFromToken copies position from a trees.Token to a rhino Node.
func (f *IRFactory) setSourceInfoFromToken(node *rhino.Node, tok *trees.Token) {
	if node.Lineno() == 0 {
		node.SetPos(tok.Location.Start.Line+1, tok.Location.Start.Column)
	}
}

// ---- token-type mapping functions ---------------------------------

// transformUnaryTokenType maps a parse-tree TokenType to the rhino unary
// operator token. Mirrors IRFactory.transformUnaryTokenType(TokenType).
func transformUnaryTokenType(t trees.TokenType) rhino.Token {
	switch t {
	case trees.TokenTypeBang:
		return rhino.NOT
	case trees.TokenTypeTilde:
		return rhino.BITNOT
	case trees.TokenTypePlus:
		return rhino.POS
	case trees.TokenTypeMinus:
		return rhino.NEG
	case trees.TokenTypeDelete:
		return rhino.DELPROP
	case trees.TokenTypeTypeof:
		return rhino.TYPEOF
	case trees.TokenTypeVoid:
		return rhino.VOID
	default:
		return rhino.EMPTY
	}
}

// transformUpdateTokenType maps a parse-tree TokenType to the rhino
// increment/decrement token.
func transformUpdateTokenType(t trees.TokenType) rhino.Token {
	switch t {
	case trees.TokenTypePlusPlus:
		return rhino.INC
	case trees.TokenTypeMinusMinus:
		return rhino.DEC
	default:
		return rhino.EMPTY
	}
}

// transformBinaryTokenType maps a parse-tree TokenType to the rhino binary
// operator token. Mirrors IRFactory.transformBinaryTokenType(TokenType).
func transformBinaryTokenType(t trees.TokenType) rhino.Token {
	switch t {
	case trees.TokenTypeBar:
		return rhino.BITOR
	case trees.TokenTypeCaret:
		return rhino.BITXOR
	case trees.TokenTypeAmpersand:
		return rhino.BITAND
	case trees.TokenTypeEqualEqual:
		return rhino.EQ
	case trees.TokenTypeNotEqual:
		return rhino.NE
	case trees.TokenTypeOpenAngle:
		return rhino.LT
	case trees.TokenTypeLessEqual:
		return rhino.LE
	case trees.TokenTypeCloseAngle:
		return rhino.GT
	case trees.TokenTypeGreaterEqual:
		return rhino.GE
	case trees.TokenTypeLeftShift:
		return rhino.LSH
	case trees.TokenTypeRightShift:
		return rhino.RSH
	case trees.TokenTypeUnsignedRightShift:
		return rhino.URSH
	case trees.TokenTypePlus:
		return rhino.ADD
	case trees.TokenTypeMinus:
		return rhino.SUB
	case trees.TokenTypeStar:
		return rhino.MUL
	case trees.TokenTypeSlash:
		return rhino.DIV
	case trees.TokenTypePercent:
		return rhino.MOD
	case trees.TokenTypeStarStar:
		return rhino.EXPONENT
	case trees.TokenTypeEqualEqualEqual:
		return rhino.SHEQ
	case trees.TokenTypeNotEqualEqual:
		return rhino.SHNE
	case trees.TokenTypeIn:
		return rhino.IN
	case trees.TokenTypeInstanceof:
		return rhino.INSTANCEOF
	case trees.TokenTypeComma:
		return rhino.COMMA
	case trees.TokenTypeEqual:
		return rhino.ASSIGN
	case trees.TokenTypeBarEqual:
		return rhino.ASSIGN_BITOR
	case trees.TokenTypeCaretEqual:
		return rhino.ASSIGN_BITXOR
	case trees.TokenTypeAmpersandEqual:
		return rhino.ASSIGN_BITAND
	case trees.TokenTypeLeftShiftEqual:
		return rhino.ASSIGN_LSH
	case trees.TokenTypeRightShiftEqual:
		return rhino.ASSIGN_RSH
	case trees.TokenTypeUnsignedRightShiftEqual:
		return rhino.ASSIGN_URSH
	case trees.TokenTypePlusEqual:
		return rhino.ASSIGN_ADD
	case trees.TokenTypeMinusEqual:
		return rhino.ASSIGN_SUB
	case trees.TokenTypeStarEqual:
		return rhino.ASSIGN_MUL
	case trees.TokenTypeStarStarEqual:
		return rhino.ASSIGN_EXPONENT
	case trees.TokenTypeSlashEqual:
		return rhino.ASSIGN_DIV
	case trees.TokenTypePercentEqual:
		return rhino.ASSIGN_MOD
	case trees.TokenTypeOr:
		return rhino.OR
	case trees.TokenTypeAnd:
		return rhino.AND
	case trees.TokenTypeQuestionQuestion:
		return rhino.COALESCE
	case trees.TokenTypeOrEqual:
		return rhino.ASSIGN_OR
	case trees.TokenTypeAndEqual:
		return rhino.ASSIGN_AND
	case trees.TokenTypeQuestionQuestionEqual:
		return rhino.ASSIGN_COALESCE
	default:
		return rhino.EMPTY
	}
}

// ---- string/number utilities ---------------------------------------

// normalizeString does minimal string literal normalization (strip quotes).
// A full implementation would handle escape sequences; this stub is
// sufficient for the current corpus tier.
func normalizeString(raw string) string {
	if len(raw) >= 2 {
		q := raw[0]
		if (q == '"' || q == '\'' || q == '`') && raw[len(raw)-1] == q {
			return raw[1 : len(raw)-1]
		}
	}
	return raw
}

// parseFloat converts a JS number literal string to float64.
// A full implementation would handle hex/octal/binary; this handles the
// common decimal case.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// floatToString converts a float64 to a JS-style string (used for numeric
// property keys). Mirrors DToA.numberToString.
func floatToString(v float64) string {
	return fmt.Sprintf("%g", v)
}
