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
func (f *IRFactory) ProcessObjectLitKey(
	input parser.Token,
	output rhino.Token,
	identifierType IRFactory_IdentifierType) rhino.Node {
	if input.IsNull() {
		return f.CreateMissingExpressionNode()
	}

	if input.Type() == parser.TokenType_IDENTIFIER {
		return f.ProcessNameOfToken(input, output, identifierType)
	}

	literal := input
	startPos := f.linePos(literal.Location.Start)
	var jsDocInfo *rhino.JSDocInfo
	if f.jsdocTracker.hasPendingBeforeOffset(startPos.Offset) {
		jsDocInfo = f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(startPos))
	}
	var comment *rhino.NonJSDocComment
	if f.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = f.ParseNonJSDocCommentAt(startPos, true)
	}

	var node rhino.Node
	switch input.Type() {
	case parser.TokenType_NUMBER:
		node = f.NewStringNodeWithToken(output, dtoa.DToA_NumberToString(f.NormalizeNumber(literal)))
	case parser.TokenType_BIGINT:
		node = f.NewStringNodeWithToken(output, f.NormalizeBigInt(literal).String())
	default:
		node = f.alloc.Alloc().InitStringIDFromTemplate(output, f.normalizeStringToID(literal), f.templateNode)
	}

	if jsDocInfo != nil {
		node.SetJSDocInfo(jsDocInfo)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}

	f.setSourceInfoAt(node, startPos, literal.Location.End)
	node.PutBooleanProp(rhino.Node_QUOTED_PROP, true)
	return node
}

// Java: IRFactory.TransformDispatcher.processArrayLiteral(ArrayLiteralExpressionTree)
func (f *IRFactory) ProcessArrayLiteral(tree *parser.ArrayLiteralExpressionTree) rhino.Node {
	node := f.NewNode(rhino.Token_ARRAYLIT)
	node.SetTrailingComma(tree.HasTrailingComma)
	for _, child := range tree.Elements {
		c := f.Transform(child)
		node.AddChildToBack(c)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processArrayPattern(ArrayPatternTree)
func (f *IRFactory) ProcessArrayPattern(tree *parser.ArrayPatternTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)

	node := f.NewNode(rhino.Token_ARRAY_PATTERN)
	for _, child := range tree.Elements {
		var elementNode rhino.Node
		switch child.GetParseTreeBase().Type {
		case parser.ParseTreeType_DEFAULT_PARAMETER:
			// processDefaultParameter() knows how to find and apply inline
			// JSDoc to the right node
			elementNode = f.ProcessDefaultParameter(child.(*parser.DefaultParameterTree))
		case parser.ParseTreeType_ITER_REST:
			f.MaybeWarnForFeature(child, parser.FeatureSet_Feature_ARRAY_PATTERN_REST)
			elementNode = f.TransformNodeWithInlineComments(child)
		default:
			elementNode = f.TransformNodeWithInlineComments(child)
		}
		node.AddChildToBack(elementNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processObjectPattern(ObjectPatternTree)
func (f *IRFactory) ProcessObjectPattern(tree *parser.ObjectPatternTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)

	node := f.NewNode(rhino.Token_OBJECT_PATTERN)
	for _, child := range tree.Fields {
		childNode := f.ProcessObjectPatternElement(child)
		node.AddChildToBack(childNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processObjectPatternElement(ParseTree)
func (f *IRFactory) ProcessObjectPatternElement(child parser.ParseTree) rhino.Node {
	switch child.GetParseTreeBase().Type {
	case parser.ParseTreeType_DEFAULT_PARAMETER:
		// shorthand with a default value
		// let { /** inlineType */ name = default } = something;
		return f.ProcessObjectPatternShorthandWithDefault(child.(*parser.DefaultParameterTree))
	case parser.ParseTreeType_PROPERTY_NAME_ASSIGNMENT:
		return f.ProcessObjectPatternPropertyNameAssignment(child.(*parser.PropertyNameAssignmentTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_DEFINITION:
		// let {[expression]: /** inlineType */ name} = something;
		computedPropertyDefinition := child.(*parser.ComputedPropertyDefinitionTree)
		return f.ProcessObjectPatternComputedPropertyDefinition(computedPropertyDefinition)
	case parser.ParseTreeType_OBJECT_REST:
		// let {...restObject} = someObject;
		f.MaybeWarnForFeature(child, parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
		target := f.TransformNodeWithInlineComments(child.(*parser.ObjectRestTree).AssignmentTarget)
		rest := f.NewNodeWithChild(rhino.Token_OBJECT_REST, target)
		f.SetSourceInfo(rest, child)
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
func (f *IRFactory) ProcessObjectPatternShorthandWithDefault(defaultParameter *parser.DefaultParameterTree) rhino.Node {
	defaultValueNode := f.ProcessDefaultParameter(defaultParameter)
	// Store in AST as non-shorthand form & just note it was originally
	// shorthand {name: /**inlineType */ name = default }
	nameNode := defaultValueNode.GetFirstChild()
	stringKeyNode := f.MaybeAddNonJsDocComment(
		f.NewStringNodeWithToken(rhino.Token_STRING_KEY, nameNode.GetString()), defaultParameter.GetStart())
	f.SetSourceInfoFromNode(stringKeyNode, nameNode)
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
func (f *IRFactory) ProcessObjectPatternPropertyNameAssignment(
	propertyNameAssignment *parser.PropertyNameAssignmentTree) rhino.Node {
	key := f.ProcessObjectLitKey(
		propertyNameAssignment.Name, rhino.Token_STRING_KEY, IRFactory_IdentifierType_STANDARD)
	targetTree := propertyNameAssignment.Value
	var valueNode rhino.Node
	if targetTree == nil {
		// `let { /** inlineType */ key } = something;`
		// The key is also the target name.
		valueNode = f.ProcessNameWithInlineCommentsOfToken(
			propertyNameAssignment.Name, IRFactory_IdentifierType_STANDARD)
		key.SetShorthandProperty(true)
	} else {
		valueNode = f.ProcessDestructuringElementTarget(targetTree)
	}
	key.AddChildToFront(valueNode)
	return key
}

// Java: IRFactory.TransformDispatcher.processDestructuringElementTarget(ParseTree)
func (f *IRFactory) ProcessDestructuringElementTarget(targetTree parser.ParseTree) rhino.Node {
	var valueNode rhino.Node
	if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_DEFAULT_PARAMETER {
		// let {key: /** inlineType */ name = default} = something;
		// let [/** inlineType */ name = default] = something;
		// processDefaultParameter() knows how to apply the inline JSDoc, if
		// any, to the right node
		valueNode = f.ProcessDefaultParameter(targetTree.(*parser.DefaultParameterTree))
	} else if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_IDENTIFIER_EXPRESSION {
		// let {key: /** inlineType */ name} = something
		// let [/** inlineType */ name] = something
		// Allow inline JSDoc on the name, since we may well be declaring it
		// here.
		valueNode = f.ProcessNameWithInlineComments(
			targetTree.(*parser.IdentifierExpressionTree), IRFactory_IdentifierType_STANDARD)
	} else {
		// ({prop: /** string */ ns.a.b} = someObject);
		// NOTE: CheckJSDoc will report an error for this case, since we want
		// qualified names to be declared with individual statements, like
		// `/** @type {string} */ ns.a.b;`
		valueNode = f.TransformNodeWithInlineComments(targetTree)
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
func (f *IRFactory) ProcessObjectPatternComputedPropertyDefinition(
	computedPropertyDefinition *parser.ComputedPropertyDefinitionTree) rhino.Node {
	f.MaybeWarnForFeature(computedPropertyDefinition, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	expressionNode := f.Transform(computedPropertyDefinition.Property)

	valueTree := computedPropertyDefinition.Value
	valueNode := f.ProcessDestructuringElementTarget(valueTree)
	computedPropertyNode := f.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, expressionNode, valueNode)
	f.SetSourceInfo(computedPropertyNode, computedPropertyDefinition)
	return computedPropertyNode
}

// Java: IRFactory.TransformDispatcher.processString(LiteralToken)
func (f *IRFactory) ProcessString(token parser.Token) rhino.Node {
	checkTokenNotNull(token)
	jutil.CheckArgument(token.Type() == parser.TokenType_STRING)
	startPos := f.linePos(token.Location.Start)
	var comment *rhino.NonJSDocComment
	if f.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = f.ParseNonJSDocCommentAt(startPos, true)
	}

	node := f.alloc.Alloc().InitStringIDFromTemplate(rhino.Token_STRINGLIT, f.normalizeStringToID(token), f.templateNode)

	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	f.setSourceInfoAt(node, startPos, token.Location.End)
	return node
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteralToken(TemplateLiteralToken)
func (f *IRFactory) ProcessTemplateLiteralToken(token parser.Token) rhino.Node {
	checkTokenNotNull(token)
	jutil.CheckArgument(token.IsTemplateLiteral())
	raw := f.rawTokenValue(token)
	var node rhino.Node
	if token.HasError() {
		node = f.NewTemplateLitStringNode(nil, raw)
	} else {
		cooked := f.NormalizeString(token, true)
		node = f.NewTemplateLitStringNode(&cooked, raw)
	}
	f.SetSourceInfoFromToken(node, token)
	return node
}

// Java: IRFactory.TransformDispatcher.processNumberLiteral(LiteralExpressionTree)
func (f *IRFactory) ProcessNumberLiteral(literalNode *parser.LiteralExpressionTree) rhino.Node {
	value := f.NormalizeNumber(literalNode.LiteralToken)
	number := f.NewNumberNode(value)
	f.SetSourceInfo(number, literalNode)
	return number
}

// Java: IRFactory.TransformDispatcher.processBigIntLiteral(LiteralExpressionTree)
func (f *IRFactory) ProcessBigIntLiteral(literalNode *parser.LiteralExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(literalNode, parser.FeatureSet_Feature_BIGINT)
	value := f.NormalizeBigInt(literalNode.LiteralToken)
	bigint := f.NewBigIntNode(value)
	f.SetSourceInfo(bigint, literalNode)
	return bigint
}

// Java: IRFactory.TransformDispatcher.processObjectLiteral(ObjectLiteralExpressionTree)
func (f *IRFactory) ProcessObjectLiteral(objTree *parser.ObjectLiteralExpressionTree) rhino.Node {
	node := f.NewNode(rhino.Token_OBJECTLIT)
	node.SetTrailingComma(objTree.HasTrailingComma)
	maybeWarn := false
	for _, el := range objTree.PropertyNameAndValues {
		if el.GetParseTreeBase().Type == parser.ParseTreeType_DEFAULT_PARAMETER {
			// (e.g. var o = { x=4 };) This is only parsed for compatibility
			// with object patterns.
			f.errorReporter.Error(
				"Default value cannot appear at top level of an object literal.",
				f.sourceName,
				f.Lineno(el),
				0)
			continue
		} else if el.GetParseTreeBase().Type == parser.ParseTreeType_GET_ACCESSOR && f.MaybeReportGetter(el) {
			continue
		} else if el.GetParseTreeBase().Type == parser.ParseTreeType_SET_ACCESSOR && f.MaybeReportSetter(el) {
			continue
		}

		key := f.Transform(el)
		if !key.IsComputedProp() &&
			!key.IsQuotedStringKey() &&
			!key.IsSpread() &&
			!f.currentFileIsExterns {
			f.MaybeWarnKeywordProperty(key)
		}
		if key.IsShorthandProperty() {
			maybeWarn = true
		}

		node.AddChildToBack(key)
	}
	if maybeWarn {
		f.MaybeWarnForFeature(objTree, parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyDefinition(ComputedPropertyDefinitionTree)
func (f *IRFactory) ProcessComputedPropertyDefinition(tree *parser.ComputedPropertyDefinitionTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	return f.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, f.Transform(tree.Property), f.Transform(tree.Value))
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyMethod(ComputedPropertyMethodTree)
func (f *IRFactory) ProcessComputedPropertyMethod(tree *parser.ComputedPropertyMethodTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	n := f.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, f.Transform(tree.Property), f.Transform(tree.Method))
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_METHOD, true)
	if tree.Method.(*parser.FunctionDeclarationTree).IsStatic {
		n.SetStaticMember(true)
	}
	return n
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyGetter(ComputedPropertyGetterTree)
func (f *IRFactory) ProcessComputedPropertyGetter(tree *parser.ComputedPropertyGetterTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	key := f.Transform(tree.Property)
	body := f.Transform(tree.Body)
	function := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), body)
	function.SrcrefTreeIfMissing(body)
	n := f.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, key, function)
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, true)
	n.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)
	return n
}

// Java: IRFactory.TransformDispatcher.processComputedPropertySetter(ComputedPropertySetterTree)
func (f *IRFactory) ProcessComputedPropertySetter(tree *parser.ComputedPropertySetterTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	key := f.Transform(tree.Property)

	paramList := f.ProcessFormalParameterList(tree.Parameter)
	f.SetSourceInfo(paramList, tree.Parameter)

	body := f.Transform(tree.Body)

	function := rhino.IR_Function(rhino.IR_Name(""), paramList, body)
	function.SrcrefTreeIfMissing(body)

	n := f.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, key, function)
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, true)
	n.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)
	return n
}

// Java: IRFactory.TransformDispatcher.processGetAccessor(GetAccessorTree)
func (f *IRFactory) ProcessGetAccessor(tree *parser.GetAccessorTree) rhino.Node {
	identifierType := IRFactory_IdentifierType_STANDARD
	if tree.IsClassMember {
		identifierType = IRFactory_IdentifierType_CAN_BE_PRIVATE
	}
	key := f.ProcessObjectLitKey(tree.PropertyName, rhino.Token_GETTER_DEF, identifierType)
	body := f.Transform(tree.Body)
	dummyName := f.NewStringNodeWithToken(rhino.Token_NAME, "")
	f.SetSourceInfo(dummyName, tree.Body)
	paramList := f.NewNode(rhino.Token_PARAM_LIST)
	f.SetSourceInfo(paramList, tree.Body)
	value := f.NewNodeWith3Children(rhino.Token_FUNCTION, dummyName, paramList, body)
	f.SetSourceInfo(value, tree.Body)
	key.AddChildToFront(value)
	key.SetStaticMember(tree.IsStatic)
	return key
}

// Java: IRFactory.TransformDispatcher.processSetAccessor(SetAccessorTree)
func (f *IRFactory) ProcessSetAccessor(tree *parser.SetAccessorTree) rhino.Node {
	identifierType := IRFactory_IdentifierType_STANDARD
	if tree.IsClassMember {
		identifierType = IRFactory_IdentifierType_CAN_BE_PRIVATE
	}
	key := f.ProcessObjectLitKey(tree.PropertyName, rhino.Token_SETTER_DEF, identifierType)

	paramList := f.ProcessFormalParameterList(tree.Parameter)
	f.SetSourceInfo(paramList, tree.Parameter)

	body := f.Transform(tree.Body)

	dummyName := f.NewStringNodeWithToken(rhino.Token_NAME, "")
	f.SetSourceInfoFromToken(dummyName, tree.PropertyName)

	value := f.NewNodeWith3Children(rhino.Token_FUNCTION, dummyName, paramList, body)
	f.SetSourceInfo(value, tree.Body)

	key.AddChildToFront(value)
	key.SetStaticMember(tree.IsStatic)
	return key
}

// Java: IRFactory.TransformDispatcher.processPropertyNameAssignment(PropertyNameAssignmentTree)
func (f *IRFactory) ProcessPropertyNameAssignment(tree *parser.PropertyNameAssignmentTree) rhino.Node {
	key := f.ProcessObjectLitKey(tree.Name, rhino.Token_STRING_KEY, IRFactory_IdentifierType_STANDARD)
	if tree.Value != nil {
		key.AddChildToFront(f.Transform(tree.Value))
	} else {
		value := f.MaybeAddNonJsDocComment(
			f.NewStringNodeWithToken(rhino.Token_NAME, key.GetString()), tree.Name.GetStart()).
			Srcref(key)
		key.SetShorthandProperty(true)
		key.AddChildToFront(value)
	}
	return key
}

// Java: IRFactory.TransformDispatcher.processRegExpLiteral(LiteralExpressionTree)
func (f *IRFactory) ProcessRegExpLiteral(literalTree *parser.LiteralExpressionTree) rhino.Node {
	token := literalTree.LiteralToken
	literalStringNode := f.NewStringNode(f.NormalizeRegex(token))
	// TODO(johnlenz): fix the source location.
	f.SetSourceInfoFromToken(literalStringNode, token)
	node := f.NewNodeWithChild(rhino.Token_REGEXP, literalStringNode)

	rawRegex := f.rawTokenValue(token)
	// Go-only: language-adapter — '/' is ASCII, so the byte index of the last
	// '/' delimits the same substring as Java's UTF-16 index.
	lastSlash := strings.LastIndexByte(rawRegex, '/')
	flags := ""
	if lastSlash < len(rawRegex) {
		flags = rawRegex[lastSlash+1:]
	}
	f.ValidateRegExpFlags(literalTree, flags)

	if flags != "" {
		flagsNode := f.NewStringNode(flags)
		// TODO(johnlenz): fix the source location.
		f.SetSourceInfoFromToken(flagsNode, token)
		node.AddChildToBack(flagsNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.validateRegExpFlags(LiteralExpressionTree, String)
func (f *IRFactory) ValidateRegExpFlags(tree *parser.LiteralExpressionTree, flags string) {
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
				f.MaybeWarnForFeature(tree, feature)
			case 's':
				f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_S)
			case 'd':
				f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_D)
			default:
				f.errorReporter.Error(
					"Invalid RegExp flag '"+string(flag)+"'", f.sourceName, f.Lineno(tree), f.Charno(tree))
			}
		}
		return
	}
	for i := 0; i < len(flags); {
		r, width := decodeWTF8Rune(flags[i:])
		switch r {
		case 'g', 'i', 'm':
		case 'u', 'y':
			feature := parser.FeatureSet_Feature_REGEXP_FLAG_Y
			if r == 'u' {
				feature = parser.FeatureSet_Feature_REGEXP_FLAG_U
			}
			f.MaybeWarnForFeature(tree, feature)
		case 's':
			f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_S)
		case 'd':
			f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_D)
		default:
			if r >= 0x10000 {
				hi := uint16(0xD800 + ((r - 0x10000) >> 10))
				lo := uint16(0xDC00 + ((r - 0x10000) & 0x3FF))
				f.errorReporter.Error(
					"Invalid RegExp flag '"+jutil.StringFromChar(hi)+"'", f.sourceName, f.Lineno(tree), f.Charno(tree))
				f.errorReporter.Error(
					"Invalid RegExp flag '"+jutil.StringFromChar(lo)+"'", f.sourceName, f.Lineno(tree), f.Charno(tree))
			} else {
				f.errorReporter.Error(
					"Invalid RegExp flag '"+flags[i:i+width]+"'", f.sourceName, f.Lineno(tree), f.Charno(tree))
			}
		}
		i += width
	}
}

// Java: IRFactory.TransformDispatcher.processStringLiteral(LiteralExpressionTree)
func (f *IRFactory) ProcessStringLiteral(literalTree *parser.LiteralExpressionTree) rhino.Node {
	return f.ProcessString(literalTree.LiteralToken)
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteral(TemplateLiteralExpressionTree)
func (f *IRFactory) ProcessTemplateLiteral(tree *parser.TemplateLiteralExpressionTree) rhino.Node {
	f.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	templateLitNode := f.NewNode(rhino.Token_TEMPLATELIT)
	f.SetSourceInfo(templateLitNode, tree)
	var node rhino.Node
	if tree.Operand == nil {
		node = templateLitNode
	} else {
		node = f.NewNodeWith2Children(rhino.Token_TAGGED_TEMPLATELIT, f.Transform(tree.Operand), templateLitNode)
	}
	for _, child := range tree.Elements {
		templateLitNode.AddChildToBack(f.Transform(child))
	}
	if node.IsTaggedTemplateLit() {
		f.AnnotateCalls(node)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteralPortion(TemplateLiteralPortionTree)
func (f *IRFactory) ProcessTemplateLiteralPortion(tree *parser.TemplateLiteralPortionTree) rhino.Node {
	return f.ProcessTemplateLiteralToken(tree.Value)
}

// Java: IRFactory.TransformDispatcher.processTemplateSubstitution(TemplateSubstitutionTree)
func (f *IRFactory) ProcessTemplateSubstitution(tree *parser.TemplateSubstitutionTree) rhino.Node {
	return f.NewNodeWithChild(rhino.Token_TEMPLATELIT_SUB, f.Transform(tree.Expression))
}

// MaybeReportGetter reports an illegal getter and returns true if the
// language mode is too low.
//
// Java: IRFactory.TransformDispatcher.maybeReportGetter(ParseTree)
func (f *IRFactory) MaybeReportGetter(node parser.ParseTree) bool {
	f.features = f.features.With(parser.FeatureSet_Feature_GETTER)
	if f.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		f.errorReporter.Error(IRFactory_GETTER_ERROR_MESSAGE, f.sourceName, f.Lineno(node), 0)
		return true
	}
	return false
}

// MaybeReportSetter reports an illegal setter and returns true if the
// language mode is too low.
//
// Java: IRFactory.TransformDispatcher.maybeReportSetter(ParseTree)
func (f *IRFactory) MaybeReportSetter(node parser.ParseTree) bool {
	f.features = f.features.With(parser.FeatureSet_Feature_SETTER)
	if f.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		f.errorReporter.Error(IRFactory_SETTER_ERROR_MESSAGE, f.sourceName, f.Lineno(node), 0)
		return true
	}
	return false
}

// Java: IRFactory.TransformDispatcher.processBooleanLiteral(LiteralExpressionTree)
func (f *IRFactory) ProcessBooleanLiteral(literal *parser.LiteralExpressionTree) rhino.Node {
	return f.NewNode(IRFactory_TransformBooleanTokenType(literal.LiteralToken.Type()))
}

// Java: IRFactory.TransformDispatcher.processNullLiteral(LiteralExpressionTree)
func (f *IRFactory) ProcessNullLiteral(unused *parser.LiteralExpressionTree) rhino.Node {
	return f.NewNode(rhino.Token_NULL)
}

// Java: IRFactory.TransformDispatcher.processNull(NullTree)
func (f *IRFactory) ProcessNull(unused *parser.NullTree) rhino.Node {
	// NOTE: This is not a NULL literal but a placeholder node such as in
	// an array with "holes".
	return f.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processLiteralExpression(LiteralExpressionTree)
func (f *IRFactory) ProcessLiteralExpression(expr *parser.LiteralExpressionTree) rhino.Node {
	switch expr.LiteralToken.Type() {
	case parser.TokenType_NUMBER:
		return f.ProcessNumberLiteral(expr)
	case parser.TokenType_STRING:
		return f.ProcessStringLiteral(expr)
	case parser.TokenType_BIGINT:
		return f.ProcessBigIntLiteral(expr)
	case parser.TokenType_FALSE, parser.TokenType_TRUE:
		return f.ProcessBooleanLiteral(expr)
	case parser.TokenType_NULL:
		return f.ProcessNullLiteral(expr)
	case parser.TokenType_REGULAR_EXPRESSION:
		return f.ProcessRegExpLiteral(expr)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected literal type: " +
				irFactory_javaClassOf(expr.LiteralToken) +
				" type: " +
				expr.LiteralToken.Type().String()))
	}
}
