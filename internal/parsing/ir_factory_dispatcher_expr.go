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

	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// Java: IRFactory.TransformDispatcher.newStringNodeFromIdentifier(Token, IdentifierType, IdentifierToken)
func (d *IRFactory_TransformDispatcher) NewStringNodeFromIdentifier(
	type_ rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType, identifierToken parser.Token) rhino.Node {
	checkTokenNotNull(identifierToken)
	// Private properties can only ever be referenced from within the scope of
	// a class. If we're not in a class, we always use STANDARD rules.
	if !d.classScope.InScope() {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_STANDARD
	}

	// Go-only: performance — the token's name is already interned: the node
	// takes its ID (Java: the value string).
	var value strpool.ID
	if identifierToken.IsPrivateIdentifier() {
		if identifierType == IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE {
			d.MaybeWarnForFeatureToken(identifierToken, parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
		} else {
			d.errorReporter.Error(
				"Private identifiers may not be used in this context",
				d.sourceName,
				d.LinenoOfToken(identifierToken),
				d.CharnoOfToken(identifierToken))
		}
		value = identifierToken.MaybePrivateValueID()
	} else {
		value = identifierToken.ValueID()
	}

	node := d.alloc.Alloc().InitStringIDFromTemplate(type_, value, d.templateNode)
	if identifierType == IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE &&
		identifierToken.IsPrivateIdentifier() {
		node.SetPrivateIdentifier()
	}
	return node
}

// ProcessObjectLitKey transforms the given object key `input` into a Node
// with token `output`.
//
// Depending on `input`, this may add quoting, such as for numerical values.
// For example, in `{2: null}`, `2` will be transformed into a quoted string.
// This adjustment loses some accuracy about the source code, but simplifies
// the AST.
//
// Get and set accessors process their name using this method. For classes,
// this can be a private property in which case identifierType will be set to
// CAN_BE_PRIVATE. E.g. `class C { get #f() { return value; } }`.
//
// Java: IRFactory.TransformDispatcher.processObjectLitKey(parser.Token, Token, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessObjectLitKey(
	input parser.Token,
	output rhino.Token,
	identifierType IRFactory_TransformDispatcher_IdentifierType) rhino.Node {
	if input.IsNull() {
		return d.CreateMissingExpressionNode()
	}

	if input.Type() == parser.TokenType_IDENTIFIER {
		return d.ProcessNameOfToken(input, output, identifierType)
	}

	literal := input
	startPos := d.linePos(literal.Location.Start)
	var jsDocInfo *rhino.JSDocInfo
	if d.jsdocTracker.hasPendingBeforeOffset(startPos.Offset) {
		jsDocInfo = d.ParseJSDocInfoFrom(d.GetJSDocCommentAt(startPos))
	}
	var comment *rhino.NonJSDocComment
	if d.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = d.ParseNonJSDocCommentAt(startPos, true)
	}

	var node rhino.Node
	switch input.Type() {
	case parser.TokenType_NUMBER:
		node = d.NewStringNodeWithToken(output, dtoa.DToA_NumberToString(d.NormalizeNumber(literal)))
	case parser.TokenType_BIGINT:
		node = d.NewStringNodeWithToken(output, d.NormalizeBigInt(literal).String())
	default:
		node = d.alloc.Alloc().InitStringIDFromTemplate(output, d.normalizeStringToID(literal), d.templateNode)
	}

	if jsDocInfo != nil {
		node.SetJSDocInfo(jsDocInfo)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}

	d.setSourceInfoAt(node, startPos, literal.Location.End)
	node.PutBooleanProp(rhino.Node_QUOTED_PROP, true)
	return node
}

// Java: IRFactory.TransformDispatcher.processArrayLiteral(ArrayLiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessArrayLiteral(tree *parser.ArrayLiteralExpressionTree) rhino.Node {
	node := d.NewNode(rhino.Token_ARRAYLIT)
	node.SetTrailingComma(tree.HasTrailingComma)
	for _, child := range tree.Elements {
		c := d.Transform(child)
		node.AddChildToBack(c)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processArrayPattern(ArrayPatternTree)
func (d *IRFactory_TransformDispatcher) ProcessArrayPattern(tree *parser.ArrayPatternTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)

	node := d.NewNode(rhino.Token_ARRAY_PATTERN)
	for _, child := range tree.Elements {
		var elementNode rhino.Node
		switch child.GetParseTreeBase().Type {
		case parser.ParseTreeType_DEFAULT_PARAMETER:
			// processDefaultParameter() knows how to find and apply inline
			// JSDoc to the right node
			elementNode = d.ProcessDefaultParameter(child.(*parser.DefaultParameterTree))
		case parser.ParseTreeType_ITER_REST:
			d.MaybeWarnForFeature(child, parser.FeatureSet_Feature_ARRAY_PATTERN_REST)
			elementNode = d.TransformNodeWithInlineComments(child)
		default:
			elementNode = d.TransformNodeWithInlineComments(child)
		}
		node.AddChildToBack(elementNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processObjectPattern(ObjectPatternTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPattern(tree *parser.ObjectPatternTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)

	node := d.NewNode(rhino.Token_OBJECT_PATTERN)
	for _, child := range tree.Fields {
		childNode := d.ProcessObjectPatternElement(child)
		node.AddChildToBack(childNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processObjectPatternElement(ParseTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternElement(child parser.ParseTree) rhino.Node {
	switch child.GetParseTreeBase().Type {
	case parser.ParseTreeType_DEFAULT_PARAMETER:
		// shorthand with a default value
		// let { /** inlineType */ name = default } = something;
		return d.ProcessObjectPatternShorthandWithDefault(child.(*parser.DefaultParameterTree))
	case parser.ParseTreeType_PROPERTY_NAME_ASSIGNMENT:
		return d.ProcessObjectPatternPropertyNameAssignment(child.(*parser.PropertyNameAssignmentTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_DEFINITION:
		// let {[expression]: /** inlineType */ name} = something;
		computedPropertyDefinition := child.(*parser.ComputedPropertyDefinitionTree)
		return d.ProcessObjectPatternComputedPropertyDefinition(computedPropertyDefinition)
	case parser.ParseTreeType_OBJECT_REST:
		// let {...restObject} = someObject;
		d.MaybeWarnForFeature(child, parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
		target := d.TransformNodeWithInlineComments(child.(*parser.ObjectRestTree).AssignmentTarget)
		rest := d.NewNodeWithChild(rhino.Token_OBJECT_REST, target)
		d.SetSourceInfo(rest, child)
		return rest
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected object pattern element: " + child.String()))
	}
}

// ProcessObjectPatternShorthandWithDefault processes an object pattern
// element using shorthand and a default value.
//
// e.g. the `/** inlineType * / name = default` part of this code
//
//	let { /** inlineType * / name = default } = something;
//
// Java: IRFactory.TransformDispatcher.processObjectPatternShorthandWithDefault(DefaultParameterTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternShorthandWithDefault(defaultParameter *parser.DefaultParameterTree) rhino.Node {
	defaultValueNode := d.ProcessDefaultParameter(defaultParameter)
	// Store in AST as non-shorthand form & just note it was originally
	// shorthand {name: /**inlineType */ name = default }
	nameNode := defaultValueNode.GetFirstChild()
	stringKeyNode := d.MaybeAddNonJsDocComment(
		d.NewStringNodeWithToken(rhino.Token_STRING_KEY, nameNode.GetString()), defaultParameter.GetStart())
	d.SetSourceInfoFromNode(stringKeyNode, nameNode)
	stringKeyNode.SetShorthandProperty(true)
	stringKeyNode.AddChildToBack(defaultValueNode)
	return stringKeyNode
}

// ProcessObjectPatternPropertyNameAssignment processes property name
// assignments in an object pattern.
//
// Covers these cases.
//
//	let {name} = someObject;
//	let {key: name} = someObject;
//	let {key: name = something} = someObject;
//
// Java: IRFactory.TransformDispatcher.processObjectPatternPropertyNameAssignment(PropertyNameAssignmentTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternPropertyNameAssignment(
	propertyNameAssignment *parser.PropertyNameAssignmentTree) rhino.Node {
	key := d.ProcessObjectLitKey(
		propertyNameAssignment.Name, rhino.Token_STRING_KEY, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	targetTree := propertyNameAssignment.Value
	var valueNode rhino.Node
	if targetTree == nil {
		// `let { /** inlineType */ key } = something;`
		// The key is also the target name.
		valueNode = d.ProcessNameWithInlineCommentsOfToken(
			propertyNameAssignment.Name, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
		key.SetShorthandProperty(true)
	} else {
		valueNode = d.ProcessDestructuringElementTarget(targetTree)
	}
	key.AddChildToFront(valueNode)
	return key
}

// Java: IRFactory.TransformDispatcher.processDestructuringElementTarget(ParseTree)
func (d *IRFactory_TransformDispatcher) ProcessDestructuringElementTarget(targetTree parser.ParseTree) rhino.Node {
	var valueNode rhino.Node
	if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_DEFAULT_PARAMETER {
		// let {key: /** inlineType */ name = default} = something;
		// let [/** inlineType */ name = default] = something;
		// processDefaultParameter() knows how to apply the inline JSDoc, if
		// any, to the right node
		valueNode = d.ProcessDefaultParameter(targetTree.(*parser.DefaultParameterTree))
	} else if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_IDENTIFIER_EXPRESSION {
		// let {key: /** inlineType */ name} = something
		// let [/** inlineType */ name] = something
		// Allow inline JSDoc on the name, since we may well be declaring it
		// here.
		valueNode = d.ProcessNameWithInlineComments(
			targetTree.(*parser.IdentifierExpressionTree), IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	} else {
		// ({prop: /** string */ ns.a.b} = someObject);
		// NOTE: CheckJSDoc will report an error for this case, since we want
		// qualified names to be declared with individual statements, like
		// `/** @type {string} */ ns.a.b;`
		valueNode = d.TransformNodeWithInlineComments(targetTree)
	}
	return valueNode
}

// ProcessObjectPatternComputedPropertyDefinition processes a computed
// property in an object pattern.
//
// Covers these cases:
//
//	let {[expression]: name} = someObject;
//	let {[expression]: name = defaultValue} = someObject;
//
// Java: IRFactory.TransformDispatcher.processObjectPatternComputedPropertyDefinition(ComputedPropertyDefinitionTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternComputedPropertyDefinition(
	computedPropertyDefinition *parser.ComputedPropertyDefinitionTree) rhino.Node {
	d.MaybeWarnForFeature(computedPropertyDefinition, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	expressionNode := d.Transform(computedPropertyDefinition.Property)

	valueTree := computedPropertyDefinition.Value
	valueNode := d.ProcessDestructuringElementTarget(valueTree)
	computedPropertyNode := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, expressionNode, valueNode)
	d.SetSourceInfo(computedPropertyNode, computedPropertyDefinition)
	return computedPropertyNode
}

// Java: IRFactory.TransformDispatcher.processConditionalExpression(ConditionalExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessConditionalExpression(exprNode *parser.ConditionalExpressionTree) rhino.Node {
	return d.NewNodeWith3Children(
		rhino.Token_HOOK,
		d.Transform(exprNode.Condition),
		d.Transform(exprNode.Left),
		d.Transform(exprNode.Right))
}

// Java: IRFactory.TransformDispatcher.processElementGet(MemberLookupExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessElementGet(getNode *parser.MemberLookupExpressionTree) rhino.Node {
	return d.NewNodeWith2Children(
		rhino.Token_GETELEM, d.Transform(getNode.Operand), d.Transform(getNode.MemberExpression))
}

// Java: IRFactory.TransformDispatcher.processOptChainElementGet(OptionalMemberLookupExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessOptChainElementGet(getNode *parser.OptionalMemberLookupExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(getNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	getElem := d.NewNodeWith2Children(
		rhino.Token_OPTCHAIN_GETELEM,
		d.Transform(getNode.Operand),
		d.Transform(getNode.MemberExpression))
	getElem.SetIsOptionalChainStart(getNode.IsStartOfOptionalChain)
	return getElem
}

// Java: IRFactory.TransformDispatcher.processFunctionCall(CallExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessFunctionCall(callNode *parser.CallExpressionTree) rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_CALL, d.Transform(callNode.Operand))
	node.SetTrailingComma(callNode.Arguments.HasTrailingComma)
	argumentsTree := callNode.Arguments
	argCount := 0
	for _, child := range argumentsTree.Arguments {
		childNode := d.Transform(child)
		node.AddChildToBack(childNode)
		// The non-trailing comments are already attached to `childNode` in
		// `transform(child)` call. Now we must attach possible trailing
		// comments to `childNode`.

		d.attachPossibleTrailingCommentsForArgAt(
			childNode, argumentsTree.Arguments, argCount, argumentsTree.CommaPositions, argumentsTree.Location.End)
	}
	d.AnnotateCalls(node)
	return node
}

// AnnotateCalls: there are two types of calls we are interested in calls
// without explicit "this" values (what we are call "free" calls) and direct
// call to eval.
//
// Java: IRFactory.TransformDispatcher.annotateCalls(Node)
func (d *IRFactory_TransformDispatcher) AnnotateCalls(n rhino.Node) {
	jutil.CheckStateMsg1(n.IsCall() || n.IsOptChainCall() || n.IsTaggedTemplateLit(), "%s", n)

	// Keep track of of the "this" context of a call.  A call without an
	// explicit "this" is a free call.
	callee := n.GetFirstChild()

	// ignore cast nodes.
	for callee.IsCast() {
		callee = callee.GetFirstChild()
	}

	if !d.IsNormalOrOptChainGet(callee) {
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
func (d *IRFactory_TransformDispatcher) IsNormalOrOptChainGet(n rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem() || n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// GetEndOfArgCommentZones calculates, for each arg, a location (end
// SourcePosition) such that all trailing comments before this location
// correspond to that arg. Can be used while processing both function calls
// (ArgsList) as well as declarations (ParamList).
//
// args: list of arguments or formal parameters (ParseTree nodes).
// commaPositions: list of SourcePositions corresponding to commas in the
// argsList or formal parameter list. argListEndPosition: SourcePosition of
// the end of argsList or paramList.
//
// Java: IRFactory.TransformDispatcher.getEndOfArgCommentZones(ImmutableList<ParseTree>, ImmutableList<SourcePosition>, SourcePosition)
func (d *IRFactory_TransformDispatcher) GetEndOfArgCommentZones(
	args []parser.ParseTree,
	commaPositions []parser.SourcePosition,
	argListEndPosition parser.SourcePosition) []parser.LinePosition {
	zones := make([]parser.LinePosition, 0, len(args))
	commaCount := 0
	for _, arg := range args {
		if len(args) > commaCount+1 {
			// there is a next arg after this arg
			nextParam := args[commaCount+1]
			argEnd := d.linePos(arg.GetParseTreeBase().Location.End)
			if d.linePos(nextParam.GetParseTreeBase().Location.Start).Line > argEnd.Line {
				// Next arg is on a new line; all trailing comments on this line
				// belong to this arg create a source position to represent the
				// end of current line
				tempSourcePos := parser.LinePosition{
					Offset: math.MaxInt32,
					Line:   argEnd.Line,
					Column: math.MaxInt32,
				}
				zones = append(zones, tempSourcePos)
			} else {
				// Next arg is on the same line; trailing comments before the
				// comma belong to this arg
				commaPosition := d.linePos(commaPositions[commaCount])
				zones = append(zones, commaPosition)
			}
		} else {
			// last arg; trailing comments till the end of argList belong to
			// this arg
			zones = append(zones, d.linePos(argListEndPosition))
		}
		commaCount++
	}
	return zones
}

// AttachPossibleTrailingCommentsForArg attaches trailing comments
// associated with this arg or formal param to it.
//
// paramNode: The node to which we're attaching trailing comment. endZone:
// The end location until which we fetch pending comments for attachment.
//
// Java: IRFactory.TransformDispatcher.attachPossibleTrailingCommentsForArg(Node, SourcePosition)
func (d *IRFactory_TransformDispatcher) AttachPossibleTrailingCommentsForArg(paramNode rhino.Node, endZone parser.LinePosition) {
	trailingComment := d.ParseNonJSDocCommentAt(endZone, true)
	if trailingComment == nil {
		return
	}
	paramNode.SetTrailingNonJSDocComment(trailingComment)
}

func (d *IRFactory_TransformDispatcher) attachPossibleTrailingCommentsForArgAt(
	paramNode rhino.Node,
	args []parser.ParseTree,
	commaCount int,
	commaPositions []parser.SourcePosition,
	argListEndPosition parser.SourcePosition,
) {
	if !d.hasAnyPendingCommentBeforeOffset(argListEndPosition.Offset) {
		return
	}
	if len(args) > commaCount+1 {
		arg := args[commaCount]
		nextParam := args[commaCount+1]
		argEnd := d.linePos(arg.GetParseTreeBase().Location.End)
		if d.linePos(nextParam.GetParseTreeBase().Location.Start).Line > argEnd.Line {
			tempSourcePos := parser.LinePosition{
				Offset: math.MaxInt32,
				Line:   argEnd.Line,
				Column: math.MaxInt32,
			}
			d.AttachPossibleTrailingCommentsForArg(paramNode, tempSourcePos)
			return
		}
		d.AttachPossibleTrailingCommentsForArg(paramNode, d.linePos(commaPositions[commaCount]))
		return
	}
	d.AttachPossibleTrailingCommentsForArg(paramNode, d.linePos(argListEndPosition))
}

// AttachPossibleTrailingComment attaches trailing comments associated with
// this node.
//
// node: The node to which we're attaching trailing comment. tokenEnd: The
// end location after which we're looking for a trailing comment.
//
// Java: IRFactory.TransformDispatcher.attachPossibleTrailingComment(Node, SourcePosition)
func (d *IRFactory_TransformDispatcher) AttachPossibleTrailingComment(node rhino.Node, tokenEndOffset parser.SourcePosition) {
	if !d.includeAllComments || d.nonJsdocTracker.index >= len(d.nonJsdocTracker.source) {
		return
	}
	tokenEnd := d.linePos(tokenEndOffset)
	// Check for pending comments on the rest of the current line.
	nextLine := parser.LinePosition{Offset: math.MaxInt32, Line: tokenEnd.Line + 1, Column: 0}
	trailingComment := d.ParseTrailingNonJSDocCommentAt(tokenEnd, nextLine)
	if trailingComment == nil {
		return
	}
	node.SetTrailingNonJSDocComment(trailingComment)
}

// Java: IRFactory.TransformDispatcher.addExtraTrailingComment(Node, NonJSDocComment)
func (d *IRFactory_TransformDispatcher) AddExtraTrailingComment(node rhino.Node, lastComment *rhino.NonJSDocComment) {
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
func (d *IRFactory_TransformDispatcher) ProcessOptChainFunctionCall(callNode *parser.OptChainCallExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(callNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	node := d.NewNodeWithChild(rhino.Token_OPTCHAIN_CALL, d.Transform(callNode.Operand))
	node.SetTrailingComma(callNode.HasTrailingComma)
	for _, child := range callNode.Arguments.Arguments {
		node.AddChildToBack(d.Transform(child))
	}
	node.SetIsOptionalChainStart(callNode.IsStartOfOptionalChain)
	d.AnnotateCalls(node)
	return node
}

// Java: IRFactory.TransformDispatcher.processIterRest(IterRestTree)
func (d *IRFactory_TransformDispatcher) ProcessIterRest(tree *parser.IterRestTree) rhino.Node {
	target := d.TransformNodeWithInlineComments(tree.AssignmentTarget)
	return d.NewNodeWithChild(rhino.Token_ITER_REST, target)
}

// Java: IRFactory.TransformDispatcher.processIterSpread(IterSpreadTree)
func (d *IRFactory_TransformDispatcher) ProcessIterSpread(tree *parser.IterSpreadTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)

	return d.NewNodeWithChild(rhino.Token_ITER_SPREAD, d.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.processObjectSpread(ObjectSpreadTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectSpread(tree *parser.ObjectSpreadTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)

	return d.NewNodeWithChild(rhino.Token_OBJECT_SPREAD, d.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.markBinaryExpressionFeatures(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) MarkBinaryExpressionFeatures(exprNode *parser.BinaryOperatorTree) {
	switch exprNode.Operator.Type() {
	case parser.TokenType_STAR_STAR, parser.TokenType_STAR_STAR_EQUAL:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_EXPONENT_OP)
	case parser.TokenType_QUESTION_QUESTION:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_NULL_COALESCE_OP)
	case parser.TokenType_QUESTION_QUESTION_EQUAL:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_NULL_COALESCE_OP)
		fallthrough
	case parser.TokenType_OR_EQUAL, parser.TokenType_AND_EQUAL:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	default:
	}
}

// Java: IRFactory.TransformDispatcher.processBinaryExpression(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) ProcessBinaryExpression(exprNode *parser.BinaryOperatorTree) rhino.Node {
	rightStartPos := exprNode.Right.GetParseTreeBase().Location.Start
	if d.jsdocTracker.hasPendingBeforeOffset(rightStartPos.Offset) || d.nonJsdocTracker.hasPendingBeforeOffset(rightStartPos.Offset) {
		rightStart := d.linePos(rightStartPos)
		if d.jsdocTracker.HasPendingCommentBefore(rightStart) || d.nonJsdocTracker.HasPendingCommentBefore(rightStart) {
			d.MarkBinaryExpressionFeatures(exprNode)
			return d.NewNodeWith2Children(
				IRFactory_TransformBinaryTokenType(exprNode.Operator.Type()),
				d.TransformLhsOfBinaryTree(exprNode),
				d.Transform(exprNode.Right))
		}
	}
	// No pending comments, we can traverse out of order.
	return d.ProcessBinaryExpressionHelper(exprNode)
}

// ProcessBinaryExpressionHelper: deep binary ops (typical string
// concatentations) can cause stack overflows, avoid recursing in this case
// and loop instead.
//
// Java: IRFactory.TransformDispatcher.processBinaryExpressionHelper(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) ProcessBinaryExpressionHelper(exprTree *parser.BinaryOperatorTree) rhino.Node {
	var root rhino.Node
	var current rhino.Node
	var previous rhino.Node
	for exprTree != nil {
		d.MarkBinaryExpressionFeatures(exprTree)
		previous = current
		// Skip the first child but recurse normally into the right operand as
		// typically this isn't deep and because we have already checked that
		// there isn't any JSDoc we can traverse out of order.
		current = d.NewNodeWithChild(
			IRFactory_TransformBinaryTokenType(exprTree.Operator.Type()), d.Transform(exprTree.Right))
		// We have inlined "transform" here. Normally, we would need to handle
		// the JSDoc here but we know there is no JSDoc to attach, which
		// simplifies things.
		d.SetSourceInfo(current, exprTree)

		// As the iteration continues add the left operand.
		if !previous.IsNil() {
			previous.AddChildToFront(current)
		}

		if binaryOperatorTree, ok := exprTree.Left.(*parser.BinaryOperatorTree); ok {
			// continue with the left hand child
			exprTree = binaryOperatorTree
		} else {
			// Finish things off, add the left operand to the current node.
			leftNode := d.TransformLhsOfBinaryTree(exprTree)
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
func (d *IRFactory_TransformDispatcher) TransformLhsOfBinaryTree(exprTree *parser.BinaryOperatorTree) rhino.Node {
	if exprTree.Operator.Type() == parser.TokenType_IN {
		if identExpr, ok := exprTree.Left.(*parser.IdentifierExpressionTree); ok && identExpr.IdentifierToken.IsPrivateIdentifier() {
			ignored := d.privateIdLhsOfInScope.MaybeIncrement(true)
			defer ignored.Close()
			return d.Transform(exprTree.Left)
		}
	}
	return d.Transform(exprTree.Left)
}

// Java: IRFactory.TransformDispatcher.processThisExpression(ThisExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessThisExpression(unused *parser.ThisExpressionTree) rhino.Node {
	return d.NewNode(rhino.Token_THIS)
}

// Java: IRFactory.TransformDispatcher.processName(IdentifierExpressionTree, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessName(nameNode *parser.IdentifierExpressionTree, identifierType IRFactory_TransformDispatcher_IdentifierType) rhino.Node {
	return d.ProcessNameOfToken(nameNode.IdentifierToken, rhino.Token_NAME, identifierType)
}

// Java: IRFactory.TransformDispatcher.processName(IdentifierToken, Token, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessNameOfToken(
	identifierToken parser.Token, output rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType) rhino.Node {
	checkTokenNotNull(identifierToken)
	startPos := d.linePos(identifierToken.Location.Start)
	var comment *rhino.NonJSDocComment
	if d.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = d.ParseNonJSDocCommentAt(startPos, true)
	}

	node := d.NewStringNodeFromIdentifier(output, identifierType, identifierToken)

	if output == rhino.Token_NAME {
		d.MaybeWarnReservedKeyword(identifierToken)

		if d.jsdocTracker.hasPendingBeforeOffset(startPos.Offset) {
			info := d.ParseJSDocInfoFrom(d.GetJSDocCommentAt(startPos))
			if info != nil {
				node.SetJSDocInfo(info)
			}
		}
	}

	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	d.setSourceInfoAt(node, startPos, identifierToken.Location.End)
	return node
}

// Java: IRFactory.TransformDispatcher.processString(LiteralToken)
func (d *IRFactory_TransformDispatcher) ProcessString(token parser.Token) rhino.Node {
	checkTokenNotNull(token)
	jutil.CheckArgument(token.Type() == parser.TokenType_STRING)
	startPos := d.linePos(token.Location.Start)
	var comment *rhino.NonJSDocComment
	if d.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = d.ParseNonJSDocCommentAt(startPos, true)
	}

	node := d.alloc.Alloc().InitStringIDFromTemplate(rhino.Token_STRINGLIT, d.normalizeStringToID(token), d.templateNode)

	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	d.setSourceInfoAt(node, startPos, token.Location.End)
	return node
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteralToken(TemplateLiteralToken)
func (d *IRFactory_TransformDispatcher) ProcessTemplateLiteralToken(token parser.Token) rhino.Node {
	checkTokenNotNull(token)
	jutil.CheckArgument(token.IsTemplateLiteral())
	raw := d.rawTokenValue(token)
	var node rhino.Node
	if token.HasError() {
		node = d.NewTemplateLitStringNode(nil, raw)
	} else {
		cooked := d.NormalizeString(token, true)
		node = d.NewTemplateLitStringNode(&cooked, raw)
	}
	d.SetSourceInfoFromToken(node, token)
	return node
}

// Java: IRFactory.TransformDispatcher.processNameWithInlineComments(IdentifierExpressionTree, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessNameWithInlineComments(
	identifierExpression *parser.IdentifierExpressionTree, identifierType IRFactory_TransformDispatcher_IdentifierType) rhino.Node {
	return d.ProcessNameWithInlineCommentsOfToken(identifierExpression.IdentifierToken, identifierType)
}

// Java: IRFactory.TransformDispatcher.processNameWithInlineComments(IdentifierToken, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessNameWithInlineCommentsOfToken(
	identifierToken parser.Token, identifierType IRFactory_TransformDispatcher_IdentifierType) rhino.Node {
	checkTokenNotNull(identifierToken)
	startPos := d.linePos(identifierToken.Location.Start)
	info := d.ParseInlineJSDocAt(startPos)
	var comment *rhino.NonJSDocComment
	if d.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = d.ParseNonJSDocCommentAt(startPos, false)
	}

	d.MaybeWarnReservedKeyword(identifierToken)
	node := d.NewStringNodeFromIdentifier(rhino.Token_NAME, identifierType, identifierToken)

	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	d.setSourceInfoAt(node, startPos, identifierToken.Location.End)
	return node
}

// Java: IRFactory.TransformDispatcher.maybeWarnKeywordProperty(Node)
func (d *IRFactory_TransformDispatcher) MaybeWarnKeywordProperty(node rhino.Node) {
	if d.currentFileIsExterns {
		return
	}
	s := node.GetString()
	if len(s) < 2 || len(s) > 12 || s[0] < 'a' || s[0] > 'y' {
		return
	}
	if !rhino.TokenStream_IsKeyword(s) {
		return
	}

	d.features = d.features.With(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	if d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		d.errorReporter.Warning(
			IRFactory_INVALID_ES3_PROP_NAME, d.sourceName, node.GetLineno(), node.GetCharno())
	}
}

// Java: IRFactory.TransformDispatcher.maybeWarnReservedKeyword(IdentifierToken)
func (d *IRFactory_TransformDispatcher) MaybeWarnReservedKeyword(token parser.Token) {
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
		d.features = d.features.With(parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS)
		isIdentifier = d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3
	}
	if token.IsKeyword() && d.isReservedKeyword(identifier) {
		d.features = d.features.With(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
		isIdentifier = d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3
	}
	if isIdentifier {
		d.errorReporter.Error(
			"identifier is a reserved word",
			d.sourceName,
			d.LinenoOfPosition(token.Location.Start),
			d.CharnoOfPosition(token.Location.Start))
	}
}

// Java: IRFactory.TransformDispatcher.processNewExpression(NewExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNewExpression(exprNode *parser.NewExpressionTree) rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_NEW, d.Transform(exprNode.Operand))
	node.SetTrailingComma(exprNode.HasTrailingComma)
	if exprNode.Arguments != nil {
		for _, arg := range exprNode.Arguments.Arguments {
			node.AddChildToBack(d.Transform(arg))
		}
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processNumberLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNumberLiteral(literalNode *parser.LiteralExpressionTree) rhino.Node {
	value := d.NormalizeNumber(literalNode.LiteralToken)
	number := d.NewNumberNode(value)
	d.SetSourceInfo(number, literalNode)
	return number
}

// Java: IRFactory.TransformDispatcher.processBigIntLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessBigIntLiteral(literalNode *parser.LiteralExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(literalNode, parser.FeatureSet_Feature_BIGINT)
	value := d.NormalizeBigInt(literalNode.LiteralToken)
	bigint := d.NewBigIntNode(value)
	d.SetSourceInfo(bigint, literalNode)
	return bigint
}

// Java: IRFactory.TransformDispatcher.processObjectLiteral(ObjectLiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectLiteral(objTree *parser.ObjectLiteralExpressionTree) rhino.Node {
	node := d.NewNode(rhino.Token_OBJECTLIT)
	node.SetTrailingComma(objTree.HasTrailingComma)
	maybeWarn := false
	for _, el := range objTree.PropertyNameAndValues {
		if el.GetParseTreeBase().Type == parser.ParseTreeType_DEFAULT_PARAMETER {
			// (e.g. var o = { x=4 };) This is only parsed for compatibility
			// with object patterns.
			d.errorReporter.Error(
				"Default value cannot appear at top level of an object literal.",
				d.sourceName,
				d.Lineno(el),
				0)
			continue
		} else if el.GetParseTreeBase().Type == parser.ParseTreeType_GET_ACCESSOR && d.MaybeReportGetter(el) {
			continue
		} else if el.GetParseTreeBase().Type == parser.ParseTreeType_SET_ACCESSOR && d.MaybeReportSetter(el) {
			continue
		}

		key := d.Transform(el)
		if !key.IsComputedProp() &&
			!key.IsQuotedStringKey() &&
			!key.IsSpread() &&
			!d.currentFileIsExterns {
			d.MaybeWarnKeywordProperty(key)
		}
		if key.IsShorthandProperty() {
			maybeWarn = true
		}

		node.AddChildToBack(key)
	}
	if maybeWarn {
		d.MaybeWarnForFeature(objTree, parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyDefinition(ComputedPropertyDefinitionTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyDefinition(tree *parser.ComputedPropertyDefinitionTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	return d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, d.Transform(tree.Property), d.Transform(tree.Value))
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyMethod(ComputedPropertyMethodTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyMethod(tree *parser.ComputedPropertyMethodTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	n := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, d.Transform(tree.Property), d.Transform(tree.Method))
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_METHOD, true)
	if tree.Method.(*parser.FunctionDeclarationTree).IsStatic {
		n.SetStaticMember(true)
	}
	return n
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyGetter(ComputedPropertyGetterTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyGetter(tree *parser.ComputedPropertyGetterTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	key := d.Transform(tree.Property)
	body := d.Transform(tree.Body)
	function := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), body)
	function.SrcrefTreeIfMissing(body)
	n := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, key, function)
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, true)
	n.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)
	return n
}

// Java: IRFactory.TransformDispatcher.processComputedPropertySetter(ComputedPropertySetterTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertySetter(tree *parser.ComputedPropertySetterTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	key := d.Transform(tree.Property)

	paramList := d.ProcessFormalParameterList(tree.Parameter)
	d.SetSourceInfo(paramList, tree.Parameter)

	body := d.Transform(tree.Body)

	function := rhino.IR_Function(rhino.IR_Name(""), paramList, body)
	function.SrcrefTreeIfMissing(body)

	n := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, key, function)
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, true)
	n.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)
	return n
}

// Java: IRFactory.TransformDispatcher.processGetAccessor(GetAccessorTree)
func (d *IRFactory_TransformDispatcher) ProcessGetAccessor(tree *parser.GetAccessorTree) rhino.Node {
	identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
	if tree.IsClassMember {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	}
	key := d.ProcessObjectLitKey(tree.PropertyName, rhino.Token_GETTER_DEF, identifierType)
	body := d.Transform(tree.Body)
	dummyName := d.NewStringNodeWithToken(rhino.Token_NAME, "")
	d.SetSourceInfo(dummyName, tree.Body)
	paramList := d.NewNode(rhino.Token_PARAM_LIST)
	d.SetSourceInfo(paramList, tree.Body)
	value := d.NewNodeWith3Children(rhino.Token_FUNCTION, dummyName, paramList, body)
	d.SetSourceInfo(value, tree.Body)
	key.AddChildToFront(value)
	key.SetStaticMember(tree.IsStatic)
	return key
}

// Java: IRFactory.TransformDispatcher.processSetAccessor(SetAccessorTree)
func (d *IRFactory_TransformDispatcher) ProcessSetAccessor(tree *parser.SetAccessorTree) rhino.Node {
	identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
	if tree.IsClassMember {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	}
	key := d.ProcessObjectLitKey(tree.PropertyName, rhino.Token_SETTER_DEF, identifierType)

	paramList := d.ProcessFormalParameterList(tree.Parameter)
	d.SetSourceInfo(paramList, tree.Parameter)

	body := d.Transform(tree.Body)

	dummyName := d.NewStringNodeWithToken(rhino.Token_NAME, "")
	d.SetSourceInfoFromToken(dummyName, tree.PropertyName)

	value := d.NewNodeWith3Children(rhino.Token_FUNCTION, dummyName, paramList, body)
	d.SetSourceInfo(value, tree.Body)

	key.AddChildToFront(value)
	key.SetStaticMember(tree.IsStatic)
	return key
}

// Java: IRFactory.TransformDispatcher.processPropertyNameAssignment(PropertyNameAssignmentTree)
func (d *IRFactory_TransformDispatcher) ProcessPropertyNameAssignment(tree *parser.PropertyNameAssignmentTree) rhino.Node {
	key := d.ProcessObjectLitKey(tree.Name, rhino.Token_STRING_KEY, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	if tree.Value != nil {
		key.AddChildToFront(d.Transform(tree.Value))
	} else {
		value := d.MaybeAddNonJsDocComment(
			d.NewStringNodeWithToken(rhino.Token_NAME, key.GetString()), tree.Name.GetStart()).
			Srcref(key)
		key.SetShorthandProperty(true)
		key.AddChildToFront(value)
	}
	return key
}

// Java: IRFactory.TransformDispatcher.checkParenthesizedExpression(ParenExpressionTree)
func (d *IRFactory_TransformDispatcher) CheckParenthesizedExpression(exprNode *parser.ParenExpressionTree) {
	if exprNode.Expression.GetParseTreeBase().Type == parser.ParseTreeType_COMMA_EXPRESSION {
		commaNodes := exprNode.Expression.(*parser.CommaExpressionTree).Expressions
		lastChild := commaNodes[len(commaNodes)-1]
		if lastChild.GetParseTreeBase().Type == parser.ParseTreeType_ITER_REST {
			d.errorReporter.Error(
				"A rest parameter must be in a parameter list.",
				d.sourceName,
				d.Lineno(lastChild),
				d.Charno(lastChild))
		}
	}
}

// Java: IRFactory.TransformDispatcher.processParenthesizedExpression(ParenExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessParenthesizedExpression(exprNode *parser.ParenExpressionTree) rhino.Node {
	d.CheckParenthesizedExpression(exprNode)
	expr := d.Transform(exprNode.Expression)
	expr.SetIsParenthesized(true)
	return expr
}

// Java: IRFactory.TransformDispatcher.processPropertyGet(MemberExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessPropertyGet(getNode *parser.MemberExpressionTree) rhino.Node {
	leftChild := d.Transform(getNode.Operand)
	propName := getNode.MemberName
	if propName.IsNull() {
		return leftChild
	}

	// Private element get prop not valid on super receiver.
	identifierType := IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	if leftChild.IsSuper() {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_STANDARD
	}
	getProp := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(rhino.Token_GETPROP, identifierType, propName),
		propName.Location.Start)
	getProp.AddChildToBack(leftChild)
	d.SetSourceInfoFromToken(getProp, propName)
	d.MaybeWarnKeywordProperty(getProp)
	return getProp
}

// Java: IRFactory.TransformDispatcher.processOptChainPropertyGet(OptionalMemberExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessOptChainPropertyGet(getNode *parser.OptionalMemberExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(getNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)

	leftChild := d.Transform(getNode.Operand)
	propName := getNode.MemberName
	if propName.IsNull() {
		return leftChild
	}

	// Private element get prop not valid on super receiver.
	identifierType := IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	if leftChild.IsSuper() {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_STANDARD
	}
	getProp := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(rhino.Token_OPTCHAIN_GETPROP, identifierType, propName),
		propName.Location.Start)
	getProp.AddChildToBack(leftChild)
	getProp.SetIsOptionalChainStart(getNode.IsStartOfOptionalChain)
	d.SetSourceInfoFromToken(getProp, propName)
	d.MaybeWarnKeywordProperty(getProp)
	return getProp
}

// Java: IRFactory.TransformDispatcher.processRegExpLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessRegExpLiteral(literalTree *parser.LiteralExpressionTree) rhino.Node {
	token := literalTree.LiteralToken
	literalStringNode := d.NewStringNode(d.NormalizeRegex(token))
	// TODO(johnlenz): fix the source location.
	d.SetSourceInfoFromToken(literalStringNode, token)
	node := d.NewNodeWithChild(rhino.Token_REGEXP, literalStringNode)

	rawRegex := d.rawTokenValue(token)
	// Go-only: language-adapter — '/' is ASCII, so the byte index of the last
	// '/' delimits the same substring as Java's UTF-16 index.
	lastSlash := strings.LastIndexByte(rawRegex, '/')
	flags := ""
	if lastSlash < len(rawRegex) {
		flags = rawRegex[lastSlash+1:]
	}
	d.ValidateRegExpFlags(literalTree, flags)

	if flags != "" {
		flagsNode := d.NewStringNode(flags)
		// TODO(johnlenz): fix the source location.
		d.SetSourceInfoFromToken(flagsNode, token)
		node.AddChildToBack(flagsNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.validateRegExpFlags(LiteralExpressionTree, String)
func (d *IRFactory_TransformDispatcher) ValidateRegExpFlags(tree *parser.LiteralExpressionTree, flags string) {
	if flags == "" {
		return
	}
	allASCII := true
	for i := 0; i < len(flags); i++ {
		if flags[i] >= 0x80 {
			allASCII = false
			break
		}
	}
	if allASCII {
		for i := 0; i < len(flags); i++ {
			flag := flags[i]
			switch flag {
			case 'g', 'i', 'm':
			case 'u', 'y':
				feature := parser.FeatureSet_Feature_REGEXP_FLAG_Y
				if flag == 'u' {
					feature = parser.FeatureSet_Feature_REGEXP_FLAG_U
				}
				d.MaybeWarnForFeature(tree, feature)
			case 's':
				d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_S)
			case 'd':
				d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_D)
			default:
				d.errorReporter.Error(
					"Invalid RegExp flag '"+string(flag)+"'", d.sourceName, d.Lineno(tree), d.Charno(tree))
			}
		}
		return
	}
	for _, flag := range jutil.StringToUTF16(flags) {
		switch flag {
		case 'g', 'i', 'm':
		case 'u', 'y':
			feature := parser.FeatureSet_Feature_REGEXP_FLAG_Y
			if flag == 'u' {
				feature = parser.FeatureSet_Feature_REGEXP_FLAG_U
			}
			d.MaybeWarnForFeature(tree, feature)
		case 's':
			d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_S)
		case 'd':
			d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_D)
		default:
			d.errorReporter.Error(
				"Invalid RegExp flag '"+jutil.StringFromChar(flag)+"'", d.sourceName, d.Lineno(tree), d.Charno(tree))
		}
	}
}

// Java: IRFactory.TransformDispatcher.processStringLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessStringLiteral(literalTree *parser.LiteralExpressionTree) rhino.Node {
	return d.ProcessString(literalTree.LiteralToken)
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteral(TemplateLiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessTemplateLiteral(tree *parser.TemplateLiteralExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	templateLitNode := d.NewNode(rhino.Token_TEMPLATELIT)
	d.SetSourceInfo(templateLitNode, tree)
	var node rhino.Node
	if tree.Operand == nil {
		node = templateLitNode
	} else {
		node = d.NewNodeWith2Children(rhino.Token_TAGGED_TEMPLATELIT, d.Transform(tree.Operand), templateLitNode)
	}
	for _, child := range tree.Elements {
		templateLitNode.AddChildToBack(d.Transform(child))
	}
	if node.IsTaggedTemplateLit() {
		d.AnnotateCalls(node)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteralPortion(TemplateLiteralPortionTree)
func (d *IRFactory_TransformDispatcher) ProcessTemplateLiteralPortion(tree *parser.TemplateLiteralPortionTree) rhino.Node {
	return d.ProcessTemplateLiteralToken(tree.Value)
}

// Java: IRFactory.TransformDispatcher.processTemplateSubstitution(TemplateSubstitutionTree)
func (d *IRFactory_TransformDispatcher) ProcessTemplateSubstitution(tree *parser.TemplateSubstitutionTree) rhino.Node {
	return d.NewNodeWithChild(rhino.Token_TEMPLATELIT_SUB, d.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.processUnaryExpression(UnaryExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessUnaryExpression(exprNode *parser.UnaryExpressionTree) rhino.Node {
	type_ := IRFactory_TransformUnaryTokenType(exprNode.Operator.Type())
	operand := d.Transform(exprNode.Operand)

	switch type_ {
	case rhino.Token_DELPROP:
		if !(operand.IsGetProp() ||
			operand.IsGetElem() ||
			operand.IsName() ||
			operand.IsOptChainGetProp() ||
			operand.IsOptChainGetElem()) {
			msg := "Invalid delete operand. Only properties can be deleted."
			d.errorReporter.Error(msg, d.sourceName, operand.GetLineno(), 0)
		}
	case rhino.Token_POS:
		if operand.IsBigInt() {
			d.errorReporter.Error(
				"Cannot convert a BigInt value to a number", d.sourceName, operand.GetLineno(), 0)
		}
	default:
	}

	return d.NewNodeWithChild(type_, operand)
}

// Java: IRFactory.TransformDispatcher.processUpdateExpression(UpdateExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessUpdateExpression(updateExpr *parser.UpdateExpressionTree) rhino.Node {
	type_ := IRFactory_TransformUpdateTokenType(updateExpr.Operator.Type())
	operand := d.Transform(updateExpr.Operand)
	return d.CreateUpdateNode(
		type_, updateExpr.OperatorPosition == parser.UpdateExpressionTree_OperatorPosition_POSTFIX, operand)
}

// Java: IRFactory.TransformDispatcher.createUpdateNode(Token, boolean, Node)
func (d *IRFactory_TransformDispatcher) CreateUpdateNode(type_ rhino.Token, postfix bool, operand rhino.Node) rhino.Node {
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
		d.errorReporter.Error(
			jutil.StringFormat("Invalid %s %s operand.", fix, op),
			d.sourceName,
			operand.GetLineno(),
			operand.GetCharno())
	}
	node := d.NewNodeWithChild(type_, operand)
	node.PutBooleanProp(rhino.Node_INCRDECR_PROP, postfix)
	return node
}

// Java: IRFactory.TransformDispatcher.processMissingExpression(MissingPrimaryExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessMissingExpression(unused *parser.MissingPrimaryExpressionTree) rhino.Node {
	// This will already have been reported as an error by the parser.
	// Try to create something valid that ide mode might be able to
	// continue with.
	return d.CreateMissingExpressionNode()
}

// Java: IRFactory.TransformDispatcher.createMissingNameNode()
func (d *IRFactory_TransformDispatcher) CreateMissingNameNode() rhino.Node {
	return d.NewStringNodeWithToken(rhino.Token_NAME, "__missing_name__")
}

// Java: IRFactory.TransformDispatcher.createMissingExpressionNode()
func (d *IRFactory_TransformDispatcher) CreateMissingExpressionNode() rhino.Node {
	return d.NewStringNodeWithToken(rhino.Token_NAME, "__missing_expression__")
}

// Java: IRFactory.TransformDispatcher.processIllegalToken(ParseTree)
func (d *IRFactory_TransformDispatcher) ProcessIllegalToken(node parser.ParseTree) rhino.Node {
	d.errorReporter.Error("Unsupported syntax: "+node.GetParseTreeBase().Type.String(), d.sourceName, d.Lineno(node), 0)
	return d.NewNode(rhino.Token_EMPTY)
}

// MaybeReportGetter reports an illegal getter and returns true if the
// language mode is too low.
//
// Java: IRFactory.TransformDispatcher.maybeReportGetter(ParseTree)
func (d *IRFactory_TransformDispatcher) MaybeReportGetter(node parser.ParseTree) bool {
	d.features = d.features.With(parser.FeatureSet_Feature_GETTER)
	if d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		d.errorReporter.Error(IRFactory_GETTER_ERROR_MESSAGE, d.sourceName, d.Lineno(node), 0)
		return true
	}
	return false
}

// MaybeReportSetter reports an illegal setter and returns true if the
// language mode is too low.
//
// Java: IRFactory.TransformDispatcher.maybeReportSetter(ParseTree)
func (d *IRFactory_TransformDispatcher) MaybeReportSetter(node parser.ParseTree) bool {
	d.features = d.features.With(parser.FeatureSet_Feature_SETTER)
	if d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		d.errorReporter.Error(IRFactory_SETTER_ERROR_MESSAGE, d.sourceName, d.Lineno(node), 0)
		return true
	}
	return false
}

// Java: IRFactory.TransformDispatcher.processBooleanLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessBooleanLiteral(literal *parser.LiteralExpressionTree) rhino.Node {
	return d.NewNode(IRFactory_TransformBooleanTokenType(literal.LiteralToken.Type()))
}

// Java: IRFactory.TransformDispatcher.processNullLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNullLiteral(unused *parser.LiteralExpressionTree) rhino.Node {
	return d.NewNode(rhino.Token_NULL)
}

// Java: IRFactory.TransformDispatcher.processNull(NullTree)
func (d *IRFactory_TransformDispatcher) ProcessNull(unused *parser.NullTree) rhino.Node {
	// NOTE: This is not a NULL literal but a placeholder node such as in
	// an array with "holes".
	return d.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processCommaExpression(CommaExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessCommaExpression(tree *parser.CommaExpressionTree) rhino.Node {
	root := d.NewNode(rhino.Token_COMMA)
	start := tree.Expressions[0].GetParseTreeBase().Location.Start
	end := tree.Expressions[1].GetParseTreeBase().Location.End
	d.SetSourceInfoFromPositions(root, start, end)
	for _, expr := range tree.Expressions {
		count := root.GetChildCount()
		if count < 2 {
			root.AddChildToBack(d.Transform(expr))
		} else {
			end = expr.GetParseTreeBase().Location.End
			root = d.NewNodeWith2Children(rhino.Token_COMMA, root, d.Transform(expr))
			d.SetSourceInfoFromPositions(root, start, end)
		}
	}
	return root
}

// Java: IRFactory.TransformDispatcher.processSuper(SuperExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessSuper(tree *parser.SuperExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_SUPER)
	return d.NewNode(rhino.Token_SUPER)
}

// Java: IRFactory.TransformDispatcher.processNewTarget(NewTargetExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNewTarget(tree *parser.NewTargetExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_NEW_TARGET)
	return d.NewNode(rhino.Token_NEW_TARGET)
}

// Java: IRFactory.TransformDispatcher.processYield(YieldExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessYield(tree *parser.YieldExpressionTree) rhino.Node {
	yield := d.NewNode(rhino.Token_YIELD)
	if tree.Expression != nil {
		yield.AddChildToBack(d.Transform(tree.Expression))
	}
	yield.SetYieldAll(tree.IsYieldAll)
	return yield
}

// Java: IRFactory.TransformDispatcher.processAwait(AwaitExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessAwait(tree *parser.AwaitExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	await := d.NewNode(rhino.Token_AWAIT)
	await.AddChildToBack(d.Transform(tree.Expression))
	return await
}

// Java: IRFactory.TransformDispatcher.processDynamicImport(DynamicImportTree)
func (d *IRFactory_TransformDispatcher) ProcessDynamicImport(dynamicImportNode *parser.DynamicImportTree) rhino.Node {
	d.MaybeWarnForFeature(dynamicImportNode, parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	argument := d.Transform(dynamicImportNode.Argument)
	return d.NewNodeWithChild(rhino.Token_DYNAMIC_IMPORT, argument)
}

// Java: IRFactory.TransformDispatcher.processImportMeta(ImportMetaExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessImportMeta(tree *parser.ImportMetaExpressionTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_IMPORT_META)
	return d.NewNode(rhino.Token_IMPORT_META)
}

// Java: IRFactory.TransformDispatcher.processLiteralExpression(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessLiteralExpression(expr *parser.LiteralExpressionTree) rhino.Node {
	switch expr.LiteralToken.Type() {
	case parser.TokenType_NUMBER:
		return d.ProcessNumberLiteral(expr)
	case parser.TokenType_STRING:
		return d.ProcessStringLiteral(expr)
	case parser.TokenType_BIGINT:
		return d.ProcessBigIntLiteral(expr)
	case parser.TokenType_FALSE, parser.TokenType_TRUE:
		return d.ProcessBooleanLiteral(expr)
	case parser.TokenType_NULL:
		return d.ProcessNullLiteral(expr)
	case parser.TokenType_REGULAR_EXPRESSION:
		return d.ProcessRegExpLiteral(expr)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected literal type: " +
				irFactory_javaClassOf(expr.LiteralToken) +
				" type: " +
				expr.LiteralToken.Type().String()))
	}
}
