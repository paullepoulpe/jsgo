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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
