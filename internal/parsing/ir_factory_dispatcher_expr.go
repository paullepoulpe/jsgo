// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261021

/*
 * Copyright 2013 The Closure Compiler Authors.
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

package parsing

import (
	"math"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Java: IRFactory.TransformDispatcher.newStringNodeFromIdentifier(Token, IdentifierType, IdentifierToken)
func (f *IRFactory) NewStringNodeFromIdentifier(
	type_ rhino.Token, identifierType IRFactory_IdentifierType, identifierToken parser.Token) rhino.Node {
	checkTokenNotNull(identifierToken)
	// Private properties can only ever be referenced from within the scope of
	// a class. If we're not in a class, we always use STANDARD rules.
	if !f.classScope.InScope() {
		identifierType = IRFactory_IdentifierType_STANDARD
	}

	// Go-only: performance — the token's name is already interned: the node
	// takes its ID (Java: the value string).
	var value strpool.ID
	if identifierToken.IsPrivateIdentifier() {
		if identifierType == IRFactory_IdentifierType_CAN_BE_PRIVATE {
			f.MaybeWarnForFeatureToken(identifierToken, parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
		} else {
			f.errorReporter.Error(
				"Private identifiers may not be used in this context",
				f.sourceName,
				f.LinenoOfToken(identifierToken),
				f.CharnoOfToken(identifierToken))
		}
		value = identifierToken.MaybePrivateValueID()
	} else {
		value = identifierToken.ValueID()
	}

	node := f.alloc.Alloc().InitStringIDFromTemplate(type_, value, f.templateNode)
	if identifierType == IRFactory_IdentifierType_CAN_BE_PRIVATE &&
		identifierToken.IsPrivateIdentifier() {
		node.SetPrivateIdentifier()
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processConditionalExpression(ConditionalExpressionTree)
func (f *IRFactory) ProcessConditionalExpression(exprNode *parser.ConditionalExpressionTree) rhino.Node {
	return f.NewNodeWith3Children(
		rhino.Token_HOOK,
		f.Transform(exprNode.Condition),
		f.Transform(exprNode.Left),
		f.Transform(exprNode.Right))
}

// Java: IRFactory.TransformDispatcher.processElementGet(MemberLookupExpressionTree)
func (f *IRFactory) ProcessElementGet(getNode *parser.MemberLookupExpressionTree) rhino.Node {
	return f.NewNodeWith2Children(
		rhino.Token_GETELEM, f.Transform(getNode.Operand), f.Transform(getNode.MemberExpression))
}

// Java: IRFactory.TransformDispatcher.processOptChainElementGet(OptionalMemberLookupExpressionTree)
func (f *IRFactory) ProcessOptChainElementGet(getNode *parser.OptionalMemberLookupExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(getNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	getElem := f.NewNodeWith2Children(
		rhino.Token_OPTCHAIN_GETELEM,
		f.Transform(getNode.Operand),
		f.Transform(getNode.MemberExpression))
	getElem.SetIsOptionalChainStart(getNode.IsStartOfOptionalChain)
	return getElem
}

// Java: IRFactory.TransformDispatcher.processFunctionCall(CallExpressionTree)
func (f *IRFactory) ProcessFunctionCall(callNode *parser.CallExpressionTree) rhino.Node {
	node := f.NewNodeWithChild(rhino.Token_CALL, f.Transform(callNode.Operand))
	node.SetTrailingComma(callNode.Arguments.HasTrailingComma)
	argumentsTree := callNode.Arguments
	argCount := 0
	for _, child := range argumentsTree.Arguments {
		childNode := f.Transform(child)
		node.AddChildToBack(childNode)
		// The non-trailing comments are already attached to `childNode` in
		// `transform(child)` call. Now we must attach possible trailing
		// comments to `childNode`.

		f.attachPossibleTrailingCommentsForArgAt(
			childNode, argumentsTree.Arguments, argCount, argumentsTree.CommaPositions, argumentsTree.Location.End)
	}
	f.AnnotateCalls(node)
	return node
}

// AnnotateCalls: there are two types of calls we are interested in calls
// without explicit "this" values (what we are call "free" calls) and direct
// call to eval.
//
// Java: IRFactory.TransformDispatcher.annotateCalls(Node)
func (f *IRFactory) AnnotateCalls(n rhino.Node) {
	jutil.CheckStateMsg1(n.IsCall() || n.IsOptChainCall() || n.IsTaggedTemplateLit(), "%s", n)

	// Keep track of of the "this" context of a call.  A call without an
	// explicit "this" is a free call.
	callee := n.GetFirstChild()

	// ignore cast nodes.
	for callee.IsCast() {
		callee = callee.GetFirstChild()
	}

	if !f.IsNormalOrOptChainGet(callee) {
		// This call originally was not passed a `this` value.
		// Inlining could change the callee into a property reference of some
		// kind. The code printer will recognize the `FREE_CALL` property and
		// wrap the real callee with `(0, real.callee)(args)` when necessary to
		// avoid changing the calling behavior.
		n.PutBooleanProp(rhino.Node_FREE_CALL, true)

		if callee.IsName() && callee.GetString() == "eval" {
			// Keep track of the context in which eval is called. It is
			// important to distinguish between "(0, eval)()" and "eval()".
			callee.PutBooleanProp(rhino.Node_DIRECT_EVAL, true)
		} else if callee.IsComma() && callee.GetFirstChild().IsNumber() {
			// The input code may actually already contain calls of the form
			// `(0, real.callee)(arg1, arg2)`. For example, the TypeScript
			// compiler outputs code like this in some cases.
			// This makes it hard for us to connect calls with the called
			// functions, though, so we will simplify these cases. The
			// `FREE_CALL` property we've just applied will tell the code
			// printer to put this wrapping back later.
			realCallee := callee.GetSecondChild()
			callee.ReplaceWith(realCallee.Detach())
			// TODO(bradfordcsmith): Why do I get an NPE if I try to report
			// this code change?
		}
	}
}

// IsNormalOrOptChainGet: is this a GETPROP, OPTCHAIN_GETPROP, GETELEM, or
// OPTCHAIN_GETELEM?
//
// Java: IRFactory.TransformDispatcher.isNormalOrOptChainGet(Node)
func (f *IRFactory) IsNormalOrOptChainGet(n rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem() || n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// AttachPossibleTrailingCommentsForArg attaches trailing comments
// associated with this arg or formal param to it.
//
// paramNode: The node to which we're attaching trailing comment. endZone:
// The end location until which we fetch pending comments for attachment.
//
// Java: IRFactory.TransformDispatcher.attachPossibleTrailingCommentsForArg(Node, SourcePosition)
func (f *IRFactory) AttachPossibleTrailingCommentsForArg(paramNode rhino.Node, endZone parser.LinePosition) {
	trailingComment := f.ParseNonJSDocCommentAt(endZone, true)
	if trailingComment == nil {
		return
	}
	paramNode.SetTrailingNonJSDocComment(trailingComment)
}

func (f *IRFactory) attachPossibleTrailingCommentsForArgAt(
	paramNode rhino.Node,
	args []parser.ParseTree,
	commaCount int,
	commaPositions []parser.SourcePosition,
	argListEndPosition parser.SourcePosition,
) {
	if !f.hasAnyPendingCommentBeforeOffset(argListEndPosition.Offset) {
		return
	}
	if len(args) > commaCount+1 {
		arg := args[commaCount]
		nextParam := args[commaCount+1]
		argEnd := f.linePos(arg.GetParseTreeBase().Location.End)
		if f.linePos(nextParam.GetParseTreeBase().Location.Start).Line > argEnd.Line {
			tempSourcePos := parser.LinePosition{
				Offset: math.MaxInt32,
				Line:   argEnd.Line,
				Column: math.MaxInt32,
			}
			f.AttachPossibleTrailingCommentsForArg(paramNode, tempSourcePos)
			return
		}
		f.AttachPossibleTrailingCommentsForArg(paramNode, f.linePos(commaPositions[commaCount]))
		return
	}
	f.AttachPossibleTrailingCommentsForArg(paramNode, f.linePos(argListEndPosition))
}

// AttachPossibleTrailingComment attaches trailing comments associated with
// this node.
//
// node: The node to which we're attaching trailing comment. tokenEnd: The
// end location after which we're looking for a trailing comment.
//
// Java: IRFactory.TransformDispatcher.attachPossibleTrailingComment(Node, SourcePosition)
func (f *IRFactory) AttachPossibleTrailingComment(node rhino.Node, tokenEndOffset parser.SourcePosition) {
	if !f.includeAllComments || f.nonJsdocTracker.index >= len(f.nonJsdocTracker.source) {
		return
	}
	tokenEnd := f.linePos(tokenEndOffset)
	// Check for pending comments on the rest of the current line.
	nextLine := parser.LinePosition{Offset: math.MaxInt32, Line: tokenEnd.Line + 1, Column: 0}
	trailingComment := f.ParseTrailingNonJSDocCommentAt(tokenEnd, nextLine)
	if trailingComment == nil {
		return
	}
	node.SetTrailingNonJSDocComment(trailingComment)
}

// Java: IRFactory.TransformDispatcher.addExtraTrailingComment(Node, NonJSDocComment)
func (f *IRFactory) AddExtraTrailingComment(node rhino.Node, lastComment *rhino.NonJSDocComment) {
	if lastComment == nil {
		return
	}
	if !node.HasChildren() {
		node.SetTrailingNonJSDocComment(lastComment)
		return
	}

	lastChild := node.GetLastChild()
	currentComment := lastChild.GetTrailingNonJSDocComment()

	if currentComment == nil {
		// We add a line break here as this is cannot be a trailing comment on
		// the same line as the last child.
		// TODO(user): pass in proper child tree end position to get the number
		// of line breaks right.
		newStart := parser.LinePosition{
			Offset: lastComment.GetEndPosition().Offset - 1,
			Line:   lastComment.GetEndPosition().Line - 1,
			Column: 0,
		}

		contents := "\n" + lastComment.GetCommentString()
		newlineComment := rhino.NewNonJSDocComment(newStart, lastComment.GetEndPosition(), &contents)
		node.GetLastChild().SetTrailingNonJSDocComment(newlineComment)
		return
	}

	blankLines := int(lastComment.GetStartPosition().Line - currentComment.GetEndPosition().Line)
	numWhiteSpace := 0
	if blankLines == 0 {
		numWhiteSpace = int(lastComment.GetStartPosition().Column - currentComment.GetEndPosition().Column - 1)
	}

	var comment strings.Builder
	comment.WriteString(currentComment.GetCommentString())
	for i := 0; i < numWhiteSpace; i++ {
		comment.WriteString(" ")
	}
	for i := 0; i < blankLines; i++ {
		comment.WriteString("\n")
	}
	comment.WriteString(lastComment.GetCommentString())

	contents := comment.String()
	allComments := rhino.NewNonJSDocComment(
		currentComment.GetStartPosition(), lastComment.GetEndPosition(), &contents)
	allComments.SetEndsAsLineComment(lastComment.IsEndingAsLineComment())
	allComments.SetIsInline(true)
	node.GetLastChild().SetTrailingNonJSDocComment(allComments)
}

// Java: IRFactory.TransformDispatcher.processOptChainFunctionCall(OptChainCallExpressionTree)
func (f *IRFactory) ProcessOptChainFunctionCall(callNode *parser.OptChainCallExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(callNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	node := f.NewNodeWithChild(rhino.Token_OPTCHAIN_CALL, f.Transform(callNode.Operand))
	node.SetTrailingComma(callNode.HasTrailingComma)
	for _, child := range callNode.Arguments.Arguments {
		node.AddChildToBack(f.Transform(child))
	}
	node.SetIsOptionalChainStart(callNode.IsStartOfOptionalChain)
	f.AnnotateCalls(node)
	return node
}

// Java: IRFactory.TransformDispatcher.processIterRest(IterRestTree)
func (f *IRFactory) ProcessIterRest(tree *parser.IterRestTree) rhino.Node {
	target := f.TransformNodeWithInlineComments(tree.AssignmentTarget)
	return f.NewNodeWithChild(rhino.Token_ITER_REST, target)
}

// Java: IRFactory.TransformDispatcher.processIterSpread(IterSpreadTree)
func (f *IRFactory) ProcessIterSpread(tree *parser.IterSpreadTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)

	return f.NewNodeWithChild(rhino.Token_ITER_SPREAD, f.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.processObjectSpread(ObjectSpreadTree)
func (f *IRFactory) ProcessObjectSpread(tree *parser.ObjectSpreadTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)

	return f.NewNodeWithChild(rhino.Token_OBJECT_SPREAD, f.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.markBinaryExpressionFeatures(BinaryOperatorTree)
func (f *IRFactory) MarkBinaryExpressionFeatures(exprNode *parser.BinaryOperatorTree) {
	switch exprNode.Operator.Type() {
	case parser.TokenType_STAR_STAR, parser.TokenType_STAR_STAR_EQUAL:
		f.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_EXPONENT_OP)
	case parser.TokenType_QUESTION_QUESTION:
		f.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_NULL_COALESCE_OP)
	case parser.TokenType_QUESTION_QUESTION_EQUAL:
		f.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_NULL_COALESCE_OP)
		fallthrough
	case parser.TokenType_OR_EQUAL, parser.TokenType_AND_EQUAL:
		f.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	default:
	}
}

// Java: IRFactory.TransformDispatcher.processBinaryExpression(BinaryOperatorTree)
func (f *IRFactory) ProcessBinaryExpression(exprNode *parser.BinaryOperatorTree) rhino.Node {
	rightStartPos := exprNode.Right.GetParseTreeBase().Location.Start
	if f.jsdocTracker.hasPendingBeforeOffset(rightStartPos.Offset) || f.nonJsdocTracker.hasPendingBeforeOffset(rightStartPos.Offset) {
		rightStart := f.linePos(rightStartPos)
		if f.jsdocTracker.HasPendingCommentBefore(rightStart) || f.nonJsdocTracker.HasPendingCommentBefore(rightStart) {
			f.MarkBinaryExpressionFeatures(exprNode)
			return f.NewNodeWith2Children(
				IRFactory_TransformBinaryTokenType(exprNode.Operator.Type()),
				f.TransformLhsOfBinaryTree(exprNode),
				f.Transform(exprNode.Right))
		}
	}
	// No pending comments, we can traverse out of order.
	return f.ProcessBinaryExpressionHelper(exprNode)
}

// ProcessBinaryExpressionHelper: deep binary ops (typical string
// concatentations) can cause stack overflows, avoid recursing in this case
// and loop instead.
//
// Java: IRFactory.TransformDispatcher.processBinaryExpressionHelper(BinaryOperatorTree)
func (f *IRFactory) ProcessBinaryExpressionHelper(exprTree *parser.BinaryOperatorTree) rhino.Node {
	var root rhino.Node
	var current rhino.Node
	var previous rhino.Node
	for exprTree != nil {
		f.MarkBinaryExpressionFeatures(exprTree)
		previous = current
		// Skip the first child but recurse normally into the right operand as
		// typically this isn't deep and because we have already checked that
		// there isn't any JSDoc we can traverse out of order.
		current = f.NewNodeWithChild(
			IRFactory_TransformBinaryTokenType(exprTree.Operator.Type()), f.Transform(exprTree.Right))
		// We have inlined "transform" here. Normally, we would need to handle
		// the JSDoc here but we know there is no JSDoc to attach, which
		// simplifies things.
		f.SetSourceInfo(current, exprTree)

		// As the iteration continues add the left operand.
		if !previous.IsNil() {
			previous.AddChildToFront(current)
		}

		if binaryOperatorTree, ok := exprTree.Left.(*parser.BinaryOperatorTree); ok {
			// continue with the left hand child
			exprTree = binaryOperatorTree
		} else {
			// Finish things off, add the left operand to the current node.
			leftNode := f.TransformLhsOfBinaryTree(exprTree)
			current.AddChildToFront(leftNode)
			// Nothing left to do.
			exprTree = nil
		}

		// Save the top binary op, this is the result to return.
		if root.IsNil() {
			root = current
		}
	}
	return root
}

// Java: IRFactory.TransformDispatcher.transformLhsOfBinaryTree(BinaryOperatorTree)
func (f *IRFactory) TransformLhsOfBinaryTree(exprTree *parser.BinaryOperatorTree) rhino.Node {
	if exprTree.Operator.Type() == parser.TokenType_IN {
		if identExpr, ok := exprTree.Left.(*parser.IdentifierExpressionTree); ok && identExpr.IdentifierToken.IsPrivateIdentifier() {
			ignored := f.privateIdLhsOfInScope.MaybeIncrement(true)
			defer ignored.Close()
			return f.Transform(exprTree.Left)
		}
	}
	return f.Transform(exprTree.Left)
}

// Java: IRFactory.TransformDispatcher.processThisExpression(ThisExpressionTree)
func (f *IRFactory) ProcessThisExpression(unused *parser.ThisExpressionTree) rhino.Node {
	return f.NewNode(rhino.Token_THIS)
}

// Java: IRFactory.TransformDispatcher.processName(IdentifierExpressionTree, IdentifierType)
func (f *IRFactory) ProcessName(nameNode *parser.IdentifierExpressionTree, identifierType IRFactory_IdentifierType) rhino.Node {
	return f.ProcessNameOfToken(nameNode.IdentifierToken, rhino.Token_NAME, identifierType)
}

// Java: IRFactory.TransformDispatcher.processName(IdentifierToken, Token, IdentifierType)
func (f *IRFactory) ProcessNameOfToken(
	identifierToken parser.Token, output rhino.Token, identifierType IRFactory_IdentifierType) rhino.Node {
	checkTokenNotNull(identifierToken)
	startPos := f.linePos(identifierToken.Location.Start)
	var comment *rhino.NonJSDocComment
	if f.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = f.ParseNonJSDocCommentAt(startPos, true)
	}

	node := f.NewStringNodeFromIdentifier(output, identifierType, identifierToken)

	if output == rhino.Token_NAME {
		f.MaybeWarnReservedKeyword(identifierToken)

		if f.jsdocTracker.hasPendingBeforeOffset(startPos.Offset) {
			info := f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(startPos))
			if info != nil {
				node.SetJSDocInfo(info)
			}
		}
	}

	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	f.setSourceInfoAt(node, startPos, identifierToken.Location.End)
	return node
}

// Java: IRFactory.TransformDispatcher.processNameWithInlineComments(IdentifierExpressionTree, IdentifierType)
func (f *IRFactory) ProcessNameWithInlineComments(
	identifierExpression *parser.IdentifierExpressionTree, identifierType IRFactory_IdentifierType) rhino.Node {
	return f.ProcessNameWithInlineCommentsOfToken(identifierExpression.IdentifierToken, identifierType)
}

// Java: IRFactory.TransformDispatcher.processNameWithInlineComments(IdentifierToken, IdentifierType)
func (f *IRFactory) ProcessNameWithInlineCommentsOfToken(
	identifierToken parser.Token, identifierType IRFactory_IdentifierType) rhino.Node {
	checkTokenNotNull(identifierToken)
	startPos := f.linePos(identifierToken.Location.Start)
	info := f.ParseInlineJSDocAt(startPos)
	var comment *rhino.NonJSDocComment
	if f.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = f.ParseNonJSDocCommentAt(startPos, false)
	}

	f.MaybeWarnReservedKeyword(identifierToken)
	node := f.NewStringNodeFromIdentifier(rhino.Token_NAME, identifierType, identifierToken)

	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	f.setSourceInfoAt(node, startPos, identifierToken.Location.End)
	return node
}

// Java: IRFactory.TransformDispatcher.maybeWarnKeywordProperty(Node)
func (f *IRFactory) MaybeWarnKeywordProperty(node rhino.Node) {
	if f.currentFileIsExterns {
		return
	}
	s := node.GetString()
	if len(s) < 2 || len(s) > 12 || s[0] < 'a' || s[0] > 'y' {
		return
	}
	if !rhino.TokenStream_IsKeyword(s) {
		return
	}

	f.features = f.features.With(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	if f.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		f.errorReporter.Warning(
			IRFactory_INVALID_ES3_PROP_NAME, f.sourceName, node.GetLineno(), node.GetCharno())
	}
}

// Java: IRFactory.TransformDispatcher.maybeWarnReservedKeyword(IdentifierToken)
func (f *IRFactory) MaybeWarnReservedKeyword(token parser.Token) {
	if token.IsPrivateIdentifier() {
		return
	}
	// We're just checking for keywords here, not actually using the
	// identifier so we don't care about potentially running into a misplaced
	// private identifier.
	identifier := token.GetMaybePrivateValue()
	if len(identifier) < 2 || len(identifier) > 12 || identifier[0] < 'a' || identifier[0] > 'y' {
		return
	}
	isIdentifier := false
	if rhino.TokenStream_IsKeyword(identifier) {
		f.features = f.features.With(parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS)
		isIdentifier = f.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3
	}
	if token.IsKeyword() && f.isReservedKeyword(identifier) {
		f.features = f.features.With(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
		isIdentifier = f.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3
	}
	if isIdentifier {
		f.errorReporter.Error(
			"identifier is a reserved word",
			f.sourceName,
			f.LinenoOfPosition(token.Location.Start),
			f.CharnoOfPosition(token.Location.Start))
	}
}

// Java: IRFactory.TransformDispatcher.processNewExpression(NewExpressionTree)
func (f *IRFactory) ProcessNewExpression(exprNode *parser.NewExpressionTree) rhino.Node {
	node := f.NewNodeWithChild(rhino.Token_NEW, f.Transform(exprNode.Operand))
	node.SetTrailingComma(exprNode.HasTrailingComma)
	if exprNode.Arguments != nil {
		for _, arg := range exprNode.Arguments.Arguments {
			node.AddChildToBack(f.Transform(arg))
		}
	}
	return node
}

// Java: IRFactory.TransformDispatcher.checkParenthesizedExpression(ParenExpressionTree)
func (f *IRFactory) CheckParenthesizedExpression(exprNode *parser.ParenExpressionTree) {
	if exprNode.Expression.GetParseTreeBase().Type == parser.ParseTreeType_COMMA_EXPRESSION {
		commaNodes := exprNode.Expression.(*parser.CommaExpressionTree).Expressions
		lastChild := commaNodes[len(commaNodes)-1]
		if lastChild.GetParseTreeBase().Type == parser.ParseTreeType_ITER_REST {
			f.errorReporter.Error(
				"A rest parameter must be in a parameter list.",
				f.sourceName,
				f.Lineno(lastChild),
				f.Charno(lastChild))
		}
	}
}

// Java: IRFactory.TransformDispatcher.processParenthesizedExpression(ParenExpressionTree)
func (f *IRFactory) ProcessParenthesizedExpression(exprNode *parser.ParenExpressionTree) rhino.Node {
	f.CheckParenthesizedExpression(exprNode)
	expr := f.Transform(exprNode.Expression)
	expr.SetIsParenthesized(true)
	return expr
}

// Java: IRFactory.TransformDispatcher.processPropertyGet(MemberExpressionTree)
func (f *IRFactory) ProcessPropertyGet(getNode *parser.MemberExpressionTree) rhino.Node {
	leftChild := f.Transform(getNode.Operand)
	propName := getNode.MemberName
	if propName.IsNull() {
		return leftChild
	}

	// Private element get prop not valid on super receiver.
	identifierType := IRFactory_IdentifierType_CAN_BE_PRIVATE
	if leftChild.IsSuper() {
		identifierType = IRFactory_IdentifierType_STANDARD
	}
	getProp := f.MaybeAddNonJsDocComment(
		f.NewStringNodeFromIdentifier(rhino.Token_GETPROP, identifierType, propName),
		propName.Location.Start)
	getProp.AddChildToBack(leftChild)
	f.SetSourceInfoFromToken(getProp, propName)
	f.MaybeWarnKeywordProperty(getProp)
	return getProp
}

// Java: IRFactory.TransformDispatcher.processOptChainPropertyGet(OptionalMemberExpressionTree)
func (f *IRFactory) ProcessOptChainPropertyGet(getNode *parser.OptionalMemberExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(getNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)

	leftChild := f.Transform(getNode.Operand)
	propName := getNode.MemberName
	if propName.IsNull() {
		return leftChild
	}

	// Private element get prop not valid on super receiver.
	identifierType := IRFactory_IdentifierType_CAN_BE_PRIVATE
	if leftChild.IsSuper() {
		identifierType = IRFactory_IdentifierType_STANDARD
	}
	getProp := f.MaybeAddNonJsDocComment(
		f.NewStringNodeFromIdentifier(rhino.Token_OPTCHAIN_GETPROP, identifierType, propName),
		propName.Location.Start)
	getProp.AddChildToBack(leftChild)
	getProp.SetIsOptionalChainStart(getNode.IsStartOfOptionalChain)
	f.SetSourceInfoFromToken(getProp, propName)
	f.MaybeWarnKeywordProperty(getProp)
	return getProp
}

// Java: IRFactory.TransformDispatcher.processUnaryExpression(UnaryExpressionTree)
func (f *IRFactory) ProcessUnaryExpression(exprNode *parser.UnaryExpressionTree) rhino.Node {
	type_ := IRFactory_TransformUnaryTokenType(exprNode.Operator.Type())
	operand := f.Transform(exprNode.Operand)

	switch type_ {
	case rhino.Token_DELPROP:
		if !(operand.IsGetProp() ||
			operand.IsGetElem() ||
			operand.IsName() ||
			operand.IsOptChainGetProp() ||
			operand.IsOptChainGetElem()) {
			msg := "Invalid delete operand. Only properties can be deleted."
			f.errorReporter.Error(msg, f.sourceName, operand.GetLineno(), 0)
		}
	case rhino.Token_POS:
		if operand.IsBigInt() {
			f.errorReporter.Error(
				"Cannot convert a BigInt value to a number", f.sourceName, operand.GetLineno(), 0)
		}
	default:
	}

	return f.NewNodeWithChild(type_, operand)
}

// Java: IRFactory.TransformDispatcher.processUpdateExpression(UpdateExpressionTree)
func (f *IRFactory) ProcessUpdateExpression(updateExpr *parser.UpdateExpressionTree) rhino.Node {
	type_ := IRFactory_TransformUpdateTokenType(updateExpr.Operator.Type())
	operand := f.Transform(updateExpr.Operand)
	return f.CreateUpdateNode(
		type_, updateExpr.OperatorPosition == parser.UpdateExpressionTree_OperatorPosition_POSTFIX, operand)
}

// Java: IRFactory.TransformDispatcher.createUpdateNode(Token, boolean, Node)
func (f *IRFactory) CreateUpdateNode(type_ rhino.Token, postfix bool, operand rhino.Node) rhino.Node {
	assignTarget := operand
	if operand.IsCast() {
		assignTarget = operand.GetFirstChild()
	}
	if !assignTarget.IsValidAssignmentTarget() {
		fix := "prefix"
		if postfix {
			fix = "postfix"
		}
		op := "decrement"
		if type_ == rhino.Token_INC {
			op = "increment"
		}
		f.errorReporter.Error(
			jutil.StringFormat("Invalid %s %s operand.", fix, op),
			f.sourceName,
			operand.GetLineno(),
			operand.GetCharno())
	}
	node := f.NewNodeWithChild(type_, operand)
	node.PutBooleanProp(rhino.Node_INCRDECR_PROP, postfix)
	return node
}

// Java: IRFactory.TransformDispatcher.processMissingExpression(MissingPrimaryExpressionTree)
func (f *IRFactory) ProcessMissingExpression(unused *parser.MissingPrimaryExpressionTree) rhino.Node {
	// This will already have been reported as an error by the parser.
	// Try to create something valid that ide mode might be able to
	// continue with.
	return f.CreateMissingExpressionNode()
}

// Java: IRFactory.TransformDispatcher.createMissingNameNode()
func (f *IRFactory) CreateMissingNameNode() rhino.Node {
	return f.NewStringNodeWithToken(rhino.Token_NAME, "__missing_name__")
}

// Java: IRFactory.TransformDispatcher.createMissingExpressionNode()
func (f *IRFactory) CreateMissingExpressionNode() rhino.Node {
	return f.NewStringNodeWithToken(rhino.Token_NAME, "__missing_expression__")
}

// Java: IRFactory.TransformDispatcher.processIllegalToken(ParseTree)
func (f *IRFactory) ProcessIllegalToken(node parser.ParseTree) rhino.Node {
	f.errorReporter.Error("Unsupported syntax: "+node.GetParseTreeBase().Type.String(), f.sourceName, f.Lineno(node), 0)
	return f.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processCommaExpression(CommaExpressionTree)
func (f *IRFactory) ProcessCommaExpression(tree *parser.CommaExpressionTree) rhino.Node {
	root := f.NewNode(rhino.Token_COMMA)
	start := tree.Expressions[0].GetParseTreeBase().Location.Start
	end := tree.Expressions[1].GetParseTreeBase().Location.End
	f.SetSourceInfoFromPositions(root, start, end)
	for _, expr := range tree.Expressions {
		count := root.GetChildCount()
		if count < 2 {
			root.AddChildToBack(f.Transform(expr))
		} else {
			end = expr.GetParseTreeBase().Location.End
			root = f.NewNodeWith2Children(rhino.Token_COMMA, root, f.Transform(expr))
			f.SetSourceInfoFromPositions(root, start, end)
		}
	}
	return root
}

// Java: IRFactory.TransformDispatcher.processSuper(SuperExpressionTree)
func (f *IRFactory) ProcessSuper(tree *parser.SuperExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_SUPER)
	return f.NewNode(rhino.Token_SUPER)
}

// Java: IRFactory.TransformDispatcher.processNewTarget(NewTargetExpressionTree)
func (f *IRFactory) ProcessNewTarget(tree *parser.NewTargetExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_NEW_TARGET)
	return f.NewNode(rhino.Token_NEW_TARGET)
}

// Java: IRFactory.TransformDispatcher.processYield(YieldExpressionTree)
func (f *IRFactory) ProcessYield(tree *parser.YieldExpressionTree) rhino.Node {
	yield := f.NewNode(rhino.Token_YIELD)
	if tree.Expression != nil {
		yield.AddChildToBack(f.Transform(tree.Expression))
	}
	yield.SetYieldAll(tree.IsYieldAll)
	return yield
}

// Java: IRFactory.TransformDispatcher.processAwait(AwaitExpressionTree)
func (f *IRFactory) ProcessAwait(tree *parser.AwaitExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	await := f.NewNode(rhino.Token_AWAIT)
	await.AddChildToBack(f.Transform(tree.Expression))
	return await
}

// Java: IRFactory.TransformDispatcher.processDynamicImport(DynamicImportTree)
func (f *IRFactory) ProcessDynamicImport(dynamicImportNode *parser.DynamicImportTree) rhino.Node {
	f.MaybeWarnForFeature(dynamicImportNode, parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	argument := f.Transform(dynamicImportNode.Argument)
	return f.NewNodeWithChild(rhino.Token_DYNAMIC_IMPORT, argument)
}

// Java: IRFactory.TransformDispatcher.processImportMeta(ImportMetaExpressionTree)
func (f *IRFactory) ProcessImportMeta(tree *parser.ImportMetaExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_IMPORT_META)
	return f.NewNode(rhino.Token_IMPORT_META)
}
