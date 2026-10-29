// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// There's no colon because the colon is parsed as a token before the TTL is
// extracted.
const JsDocInfoParser_TTL_START_DELIMITER = "="
const JsDocInfoParser_TTL_END_DELIMITER = "=:"

// isTemplateNameChar is CharMatcher.javaLetterOrDigit().or(CharMatcher.anyOf("_$")).
func isTemplateNameChar(c uint16) bool {
	return jutil.CharacterIsLetter(c) || jutil.CharacterIsDigit(c) || c == '_' || c == '$'
}

func isModifiesAnnotationKeyword(name string) bool {
	return name == "this" || name == "arguments"
}

func isIdGeneratorAnnotationKeyword(name string) bool {
	switch name {
	case "unique", "consistent", "stable", "mapped", "xid":
		return true
	default:
		return false
	}
}

func isJsDocPrimitiveType(name string) bool {
	switch name {
	case "number", "string", "boolean", "symbol":
		return true
	default:
		return false
	}
}

func (p *JsDocInfoParser) parseAuthorTag() JsDocToken {
	if p.jsdocBuilder.ShouldParseDocumentation() {
		authorInfo := p.ExtractSingleLineBlock()
		if authorInfo.string_ == "" {
			p.AddParserWarning(rhino.Msg_JSDOC_AUTHORMISSING)
		} else {
			p.jsdocBuilder.RecordAuthor(authorInfo.string_)
		}
		return authorInfo.token
	}
	return p.EatUntilEOLIfNotAnnotation()
}

func (p *JsDocInfoParser) parseDeprecatedTag(token JsDocToken) JsDocToken {
	if !p.jsdocBuilder.RecordDeprecated() {
		p.AddParserWarning(rhino.Msg_JSDOC_DEPRECATED)
	}
	// Find the reason/description, if any.
	reasonInfo := p.ExtractMultilineTextualBlock(token)
	if len(reasonInfo.string_) > 0 {
		p.jsdocBuilder.RecordDeprecationReason(reasonInfo.string_)
	}
	return reasonInfo.token
}

func (p *JsDocInfoParser) parseDescTag(token JsDocToken) JsDocToken {
	if p.jsdocBuilder.IsDescriptionRecorded() {
		p.AddParserWarning(rhino.Msg_JSDOC_DESC_EXTRA)
		return p.EatUntilEOLIfNotAnnotation()
	}
	descriptionInfo := p.ExtractMultilineTextualBlock(token)
	p.jsdocBuilder.RecordDescription(descriptionInfo.string_)
	return descriptionInfo.token
}

func (p *JsDocInfoParser) parseTsTypeTag(token JsDocToken) JsDocToken {
	tsTypeInfo := p.ExtractMultilineTextualBlock(token)
	if p.recordTsType {
		p.jsdocBuilder.RecordTsType(tsTypeInfo.string_)
	} else {
		p.jsdocBuilder.RecordTsType(JsDocInfoParser_TS_TYPE_PLACEHOLDER)
	}
	return tsTypeInfo.token
}

func (p *JsDocInfoParser) parseFileOverviewTag(token JsDocToken) JsDocToken {
	fileOverview := ""
	if p.jsdocBuilder.ShouldParseDocumentation() && !p.LookAheadForAnnotation() {
		fileOverviewInfo := p.ExtractMultilineTextualBlockWithOption(
			token, p.GetWhitespaceOption(JsDocInfoParser_WhitespaceOption_TRIM), false)
		fileOverview = fileOverviewInfo.string_
		token = fileOverviewInfo.token
	} else {
		token = p.EatUntilEOLIfNotAnnotation()
	}

	if !p.jsdocBuilder.RecordFileOverview(fileOverview) {
		p.AddParserWarning(rhino.Msg_JSDOC_FILEOVERVIEW_EXTRA)
	}
	return token
}

func (p *JsDocInfoParser) parseLicenseOrPreserveTag(token JsDocToken) JsDocToken {
	// Always use PRESERVE for @license and @preserve blocks.
	preserveInfo := p.ExtractMultilineTextualBlockWithOption(token, JsDocInfoParser_WhitespaceOption_PRESERVE, true)
	if len(preserveInfo.string_) > 0 {
		p.SetLicenseTextMonotonic(preserveInfo.string_)
	}
	return preserveInfo.token
}

func (p *JsDocInfoParser) parseEnumTag() JsDocToken {
	token := p.Next()
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()

	var type_ *rhino.JSTypeExpression
	if token != JsDocToken_EOL && token != JsDocToken_EOC {
		typeNode := p.ParseAndRecordTypeNode(token)
		if !typeNode.IsNil() && typeNode.IsStringLit() {
			if !isJsDocPrimitiveType(typeNode.GetString()) {
				typeNode = p.WrapNode(rhino.Token_BANG, typeNode)
			}
		}
		type_ = p.CreateJSTypeExpression(typeNode)
	} else {
		p.RestoreLookAhead(token)
	}

	if type_ == nil {
		type_ = p.CreateJSTypeExpression(p.NewStringNode("number"))
	}
	if !p.jsdocBuilder.RecordEnumParameterType(type_) {
		p.AddTypeWarningAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, lineno, charno)
	}
	return p.EatUntilEOLIfNotAnnotation()
}

func (p *JsDocInfoParser) parseExtendsOrImplementsTag(
	annotation Annotation, extendedTypes *[]JsDocInfoParser_ExtendedTypeInfo) JsDocToken {
	p.SkipEOLs()
	token := p.Next()
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()
	matchingRc := false
	if token == JsDocToken_LEFT_CURLY {
		token = p.Next()
		matchingRc = true
	}

	if token == JsDocToken_STRING {
		typeNode := p.ParseAndRecordTypeNameNode(token, lineno, charno, matchingRc)

		lineno = p.stream.GetLineno()
		charno = p.stream.GetCharno()

		typeNode = p.WrapNode(rhino.Token_BANG, typeNode)
		type_ := p.CreateJSTypeExpression(typeNode)

		if annotation == Annotation_EXTENDS {
			// record the extended type, check later
			*extendedTypes = append(*extendedTypes, JsDocInfoParser_ExtendedTypeInfo{
				type_:  type_,
				lineno: p.stream.GetLineno(),
				charno: p.stream.GetCharno(),
			})
		} else {
			jutil.CheckState(annotation == Annotation_IMPLEMENTS)
			if !p.jsdocBuilder.RecordImplementedInterface(type_) {
				p.AddTypeWarningAt(rhino.Msg_JSDOC_IMPLEMENTS_DUPLICATE, lineno, charno)
			}
		}
		token = p.Next()
		if matchingRc {
			if token != JsDocToken_RIGHT_CURLY {
				p.AddTypeWarning(rhino.Msg_JSDOC_MISSING_RC)
			} else {
				token = p.Next()
			}
		} else if token != JsDocToken_EOL &&
			token != JsDocToken_EOF &&
			token != JsDocToken_EOC {
			p.AddTypeWarning(rhino.Msg_END_ANNOTATION_EXPECTED)
		}
	} else if token == JsDocToken_BANG || token == JsDocToken_QMARK {
		p.AddTypeWarningAt(rhino.Msg_JSDOC_IMPLEMENTS_EXTRAQUALIFIER, lineno, charno)
	} else {
		p.AddTypeWarningAt(rhino.Msg_NO_TYPE_NAME, lineno, charno)
	}
	return p.EatUntilEOLIfNotAnnotationFrom(token)
}

func (p *JsDocInfoParser) parseLendsTag() JsDocToken {
	p.SkipEOLs()

	matchingRc := false
	if p.Match(JsDocToken_LEFT_CURLY) {
		p.Next()
		matchingRc = true
	}

	if p.Match(JsDocToken_STRING) {
		lendsExpression := p.CreateJSTypeExpression(p.ParseNameExpression(p.Next()))
		if !p.jsdocBuilder.RecordLends(lendsExpression) {
			p.AddTypeWarning(rhino.Msg_JSDOC_LENDS_INCOMPATIBLE)
		}
	} else {
		p.AddTypeWarning(rhino.Msg_JSDOC_LENDS_MISSING)
	}

	if matchingRc && !p.Match(JsDocToken_RIGHT_CURLY) {
		p.AddTypeWarning(rhino.Msg_JSDOC_MISSING_RC)
	}
	return p.EatUntilEOLIfNotAnnotation()
}

func (p *JsDocInfoParser) parseThrowsTag(token JsDocToken) JsDocToken {
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()
	throwsAnnotation := ""
	if !p.LookAheadForAnnotation() {
		throwsInfo := p.ExtractMultilineTextualBlock(token)
		throwsAnnotation = jutil.StringTrim(throwsInfo.string_)
		token = throwsInfo.token
	}
	if p.jsdocBuilder.IsNoSideEffectsRecorded() {
		p.AddParserWarningAt(rhino.Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS, lineno, charno)
	} else if !p.jsdocBuilder.RecordThrowsAnnotation(throwsAnnotation) {
		p.AddParserWarningAt(rhino.Msg_JSDOC_THROWS, lineno, charno)
	}
	return token
}

func (p *JsDocInfoParser) parseParamTag() JsDocToken {
	p.SkipEOLs()
	token := p.Next()
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()
	var type_ *rhino.JSTypeExpression

	hasParamType := false

	if token == JsDocToken_LEFT_CURLY {
		type_ = p.CreateJSTypeExpression(p.ParseAndRecordParamTypeNode(token))

		if type_ == nil {
			// parsing error reported during recursive descent
			// recovering parsing
			return p.EatUntilEOLIfNotAnnotation()
		}
		p.SkipEOLs()
		token = p.Next()
		lineno = p.stream.GetLineno()
		charno = p.stream.GetCharno()
		hasParamType = true
	}

	var name string
	hasName := false
	isBracketedParam := JsDocToken_LEFT_SQUARE == token
	if isBracketedParam {
		token = p.Next()
	}

	if JsDocToken_STRING != token {
		p.AddTypeWarningAt(rhino.Msg_MISSING_VARIABLE_NAME, lineno, charno)
	} else {
		if !hasParamType {
			p.AddMissingTypeWarning(p.stream.GetLineno(), p.stream.GetCharno())
		}

		name = p.stream.GetString()
		hasName = true

		if isBracketedParam {
			token = p.Next()

			// Throw out JsDocToolkit's "default" parameter
			// annotation.  It makes no sense under our type
			// system.
			if JsDocToken_EQUALS == token {
				token = p.Next()
				if JsDocToken_STRING == token {
					token = p.Next()
				}
			}

			if JsDocToken_RIGHT_SQUARE != token {
				p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RB)
			} else if type_ != nil {
				// Make the type expression optional, if it isn't
				// already.
				type_ = rhino.JSTypeExpression_MakeOptionalArg(type_)
			}
		}

		// We do not handle the JsDocToolkit method
		// for handling properties of params, so if the param name has a DOT
		// in it, report a warning and throw it out.
		// See https://github.com/google/closure-compiler/issues/499
		if !rhino.TokenStream_IsJSIdentifier(name) {
			p.AddParserWarningWithArgAt(rhino.Msg_INVALID_VARIABLE_NAME, name, lineno, charno)
			hasName = false
		} else if !p.jsdocBuilder.RecordParameter(name, type_) {
			if p.jsdocBuilder.HasParameter(name) {
				p.AddTypeWarningWithArgAt(rhino.Msg_DUP_VARIABLE_NAME, name, lineno, charno)
			} else {
				p.AddTypeWarningWithArgAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, name, lineno, charno)
			}
		}
	}

	if !hasName {
		return p.EatUntilEOLIfNotAnnotationFrom(token)
	}

	p.jsdocBuilder.MarkName(name, p.templateNode, lineno, charno)

	// Find the parameter's description (if applicable).
	if p.jsdocBuilder.ShouldParseDocumentation() && token != JsDocToken_ANNOTATION {
		paramDescriptionInfo := p.ExtractMultilineTextualBlock(token)
		if len(paramDescriptionInfo.string_) > 0 {
			p.jsdocBuilder.RecordParameterDescription(name, paramDescriptionInfo.string_)
		}
		token = paramDescriptionInfo.token
	} else if token != JsDocToken_EOC && token != JsDocToken_EOF {
		token = p.EatUntilEOLIfNotAnnotation()
	}
	return token
}

func (p *JsDocInfoParser) parseSeeTag() JsDocToken {
	if p.jsdocBuilder.ShouldParseDocumentation() {
		referenceInfo := p.ExtractSingleLineBlock()
		if referenceInfo.string_ == "" {
			p.AddParserWarning(rhino.Msg_JSDOC_SEEMISSING)
		} else {
			p.jsdocBuilder.RecordReference(referenceInfo.string_)
		}
		return referenceInfo.token
	}
	return p.EatUntilEOLIfNotAnnotation()
}

func (p *JsDocInfoParser) parseTemplateTag(lineno int, charno int) JsDocToken {
	// Attempt to parse a template bound type.
	var boundTypeExpression *rhino.JSTypeExpression
	if p.Match(JsDocToken_LEFT_CURLY) {
		p.AddParserWarningAt(rhino.Msg_JSDOC_TEMPLATE_BOUNDEDGENERICS_USED, lineno, charno)

		boundTypeNode := p.ParseTypeExpressionAnnotation(p.Next())
		if !boundTypeNode.IsNil() {
			boundTypeExpression = p.CreateJSTypeExpression(boundTypeNode)
		}
	}

	// Collect any names of template types.
	var templateNamesBuf [4]string
	templateNames := templateNamesBuf[:0]
	if !p.Match(JsDocToken_COLON) {
		// Skip colons so as not to consume TTLs accidentally.
		// Parse a list of comma separated names.
		for {
			nameNode := p.ParseNameExpression(p.Next())
			var templateName string
			if !nameNode.IsNil() {
				templateName = nameNode.GetString()
			}

			if JsDocInfoParser_ValidTemplateTypeName(templateName) {
				templateNames = append(templateNames, templateName)
			}
			if !p.EatIfMatch(JsDocToken_COMMA) {
				break
			}
		}
	}

	// Look for some TTL. Always use TRIM for template TTL expressions.
	var token JsDocToken
	var ttlAst rhino.Node
	if p.Match(JsDocToken_COLON) {
		// TODO(nickreid): Fix `extractMultilineTextualBlock` so the result includes the
		// current token. At present, it reads the stream directly and bypasses any previously
		// read tokens.
		ttlExtraction := p.ExtractMultilineTextualBlockWithOption(p.Current(), JsDocInfoParser_WhitespaceOption_TRIM, false)
		token = ttlExtraction.token
		ttlAst = p.ParseTtlAst(ttlExtraction, lineno, charno)
	} else {
		token = p.EatUntilEOLIfNotAnnotationFrom(p.Current())
	}

	// Validate the number of declared template names, which depends on bounds and TTL.
	switch len(templateNames) {
	case 0:
		p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_NAME_MISSING, lineno, charno)
		return token // There's nothing we can connect a bound or TTL to.
	case 1:
	default:
		if boundTypeExpression != nil || !ttlAst.IsNil() {
			p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_MULTIPLEDECLARATION, lineno, charno)
		}
	}

	if boundTypeExpression != nil && !ttlAst.IsNil() {
		p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_BOUNDSWITHTTL, lineno, charno)
		return token // It's undecidable what the user intent was.
	}

	// Based on the form form of the declaration, record the information in the JsDocInfo.
	if !ttlAst.IsNil() {
		if !p.jsdocBuilder.RecordTypeTransformation(templateNames[0], ttlAst) {
			p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_NAME_REDECLARATION, lineno, charno)
		}
	} else if boundTypeExpression != nil {
		if !p.jsdocBuilder.RecordTemplateTypeNameWithBound(templateNames[0], boundTypeExpression) {
			p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_NAME_REDECLARATION, lineno, charno)
		}
	} else {
		for _, templateName := range templateNames {
			if !p.jsdocBuilder.RecordTemplateTypeName(templateName) {
				p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_NAME_REDECLARATION, lineno, charno)
			}
		}
	}

	return token
}

func (p *JsDocInfoParser) parseTypeBearingTag(token JsDocToken, annotation Annotation) JsDocToken {
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()

	hasType := p.LookAheadForType()
	isAlternateTypeAnnotation := annotation == Annotation_PACKAGE ||
		annotation == Annotation_PRIVATE ||
		annotation == Annotation_PROTECTED ||
		annotation == Annotation_PUBLIC ||
		annotation == Annotation_CONSTANT ||
		annotation == Annotation_FINAL ||
		annotation == Annotation_EXPORT
	canSkipTypeAnnotation := isAlternateTypeAnnotation || annotation == Annotation_RETURN
	var type_ *rhino.JSTypeExpression

	if annotation == Annotation_RETURN && !hasType {
		p.AddMissingTypeWarning(p.stream.GetLineno(), p.stream.GetCharno())
	}

	if hasType || !canSkipTypeAnnotation {
		p.SkipEOLs()
		token = p.Next()
		typeNode := p.ParseAndRecordTypeNode(token)

		if annotation == Annotation_THIS {
			typeNode = p.WrapNode(rhino.Token_BANG, typeNode)
		}
		type_ = p.CreateJSTypeExpression(typeNode)
	}

	// The error was reported during recursive descent
	// recovering parsing
	hasError := type_ == nil && !canSkipTypeAnnotation
	if !hasError {
		// Record types for @type.
		// If the @package, @private, @protected, or @public annotations
		// have a type attached, pretend that they actually wrote:
		// @type {type}\n@private
		// This will have some weird behavior in some cases
		// (for example, @private can now be used as a type-cast),
		// but should be mostly OK.
		if ((type_ != nil && isAlternateTypeAnnotation) || annotation == Annotation_TYPE) &&
			!p.jsdocBuilder.RecordType(type_) {
			p.AddTypeWarningAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, lineno, charno)
		}

		isAnnotationNext := p.LookAheadForAnnotation()

		switch annotation {
		case Annotation_CONSTANT:
			if !p.jsdocBuilder.RecordConstancy() {
				p.AddParserWarning(rhino.Msg_JSDOC_CONST)
			}

		case Annotation_FINAL:
			if !p.jsdocBuilder.RecordFinality() {
				p.AddTypeWarning(rhino.Msg_JSDOC_FINAL)
			}

		case Annotation_DEFINE:
			if !p.jsdocBuilder.RecordDefineType(type_) {
				p.AddParserWarningAt(rhino.Msg_JSDOC_DEFINE, lineno, charno)
			}
			if !isAnnotationNext {
				return p.RecordDescription(token)
			}

		case Annotation_EXPORT:
			if !p.jsdocBuilder.RecordExport() {
				p.AddParserWarningAt(rhino.Msg_JSDOC_EXPORT, lineno, charno)
			}
			if !isAnnotationNext {
				return p.RecordDescription(token)
			}

		case Annotation_PRIVATE:
			if !p.jsdocBuilder.RecordVisibility(rhino.JSDocInfo_Visibility_PRIVATE) {
				p.AddParserWarningAt(rhino.Msg_JSDOC_EXTRA_VISIBILITY, lineno, charno)
			}
			if !isAnnotationNext {
				return p.RecordDescription(token)
			}

		case Annotation_PACKAGE:
			if !p.jsdocBuilder.RecordVisibility(rhino.JSDocInfo_Visibility_PACKAGE) {
				p.AddParserWarningAt(rhino.Msg_JSDOC_EXTRA_VISIBILITY, lineno, charno)
			}
			if !isAnnotationNext {
				return p.RecordDescription(token)
			}

		case Annotation_PROTECTED:
			if !p.jsdocBuilder.RecordVisibility(rhino.JSDocInfo_Visibility_PROTECTED) {
				p.AddParserWarningAt(rhino.Msg_JSDOC_EXTRA_VISIBILITY, lineno, charno)
			}
			if !isAnnotationNext {
				return p.RecordDescription(token)
			}

		case Annotation_PUBLIC:
			if !p.jsdocBuilder.RecordVisibility(rhino.JSDocInfo_Visibility_PUBLIC) {
				p.AddParserWarningAt(rhino.Msg_JSDOC_EXTRA_VISIBILITY, lineno, charno)
			}
			if !isAnnotationNext {
				return p.RecordDescription(token)
			}

		case Annotation_RETURN:
			if type_ == nil {
				type_ = p.CreateJSTypeExpression(p.NewNode(rhino.Token_QMARK))
			}

			if !p.jsdocBuilder.RecordReturnType(type_) {
				p.AddTypeWarningAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, lineno, charno)
				break
			}

			// TODO(johnlenz): The extractMultilineTextualBlock method
			// and friends look directly at the stream, regardless of
			// last token read, so we don't want to read the first
			// "STRING" out of the stream.

			// Find the return's description (if applicable).
			if p.jsdocBuilder.ShouldParseDocumentation() && !isAnnotationNext {
				returnDescriptionInfo := p.ExtractMultilineTextualBlock(token)
				if len(returnDescriptionInfo.string_) > 0 {
					p.jsdocBuilder.RecordReturnDescription(returnDescriptionInfo.string_)
				}
				return returnDescriptionInfo.token
			}
			return p.EatUntilEOLIfNotAnnotation()

		case Annotation_THIS:
			if !p.jsdocBuilder.RecordThisType(type_) {
				p.AddTypeWarningAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, lineno, charno)
			}

		case Annotation_TYPEDEF:
			if !p.jsdocBuilder.RecordTypedef(type_) {
				p.AddTypeWarningAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, lineno, charno)
			}

		default:
		}
	}

	return p.EatUntilEOLIfNotAnnotation()
}

// Java: JsDocInfoParser.parseTtlAst(ExtractionInfo, int, int)
func (p *JsDocInfoParser) ParseTtlAst(ttlExtraction JsDocInfoParser_ExtractionInfo, lineno int, charno int) rhino.Node {
	ttlExpression := ttlExtraction.string_
	if !strings.HasPrefix(ttlExpression, JsDocInfoParser_TTL_START_DELIMITER) {
		return rhino.NilNode
	}

	// (Both delimiters are ASCII, so byte offsets are char boundaries.)
	ttlExpression = ttlExpression[len(JsDocInfoParser_TTL_START_DELIMITER):]

	endIndex := strings.Index(ttlExpression, JsDocInfoParser_TTL_END_DELIMITER)
	if endIndex >= 0 {
		ttlExpression = ttlExpression[0:endIndex]
	} else {
		p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER, lineno, charno)
	}
	ttlExpression = jutil.StringTrim(ttlExpression)

	if ttlExpression == "" {
		p.AddTypeWarningAt(rhino.Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING, lineno, charno)
		return rhino.NilNode
	}

	// Build the AST for the type transformation
	ttlParser := NewTypeTransformationParser(ttlExpression, p.GetSourceFile(), p.errorReporter, lineno, charno)
	if !ttlParser.ParseTypeTransformation() {
		return rhino.NilNode
	}

	return ttlParser.GetTypeTransformationAst()
}

// JsDocInfoParser_ValidTemplateTypeName reports whether name is a valid
// template type name: the types in @template annotations must contain only
// letters, digits, underscores and dollar signs.
//
// Java: JsDocInfoParser.validTemplateTypeName(String)
func JsDocInfoParser_ValidTemplateTypeName(name string) bool {
	if name == "" {
		return false
	}
	// CharMatcher.matchesAllOf: every UTF-16 char matches.
	for i := 0; i < len(name); i++ {
		b := name[i]
		if b >= 0x80 {
			for _, r := range name[i:] {
				if r > 0xFFFF || !isTemplateNameChar(uint16(r)) {
					return false
				}
			}
			return true
		}
		if !((b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '_' || b == '$') {
			return false
		}
	}
	return true
}

// RecordDescription records a marker's description if there is one available
// and record it in the current marker.
//
// TODO(rishipal): This function is a misnomer as it does not record anything.
//
// Java: JsDocInfoParser.recordDescription(JsDocToken)
func (p *JsDocInfoParser) RecordDescription(token JsDocToken) JsDocToken {
	// Find marker's description (if applicable).
	if p.jsdocBuilder.ShouldParseDocumentation() {
		descriptionInfo := p.ExtractMultilineTextualBlock(token)
		token = descriptionInfo.token
	} else {
		token = p.EatTokensUntilEOLFrom(token)
	}
	return token
}

// Java: JsDocInfoParser.checkExtendedTypes(List<ExtendedTypeInfo>)
func (p *JsDocInfoParser) CheckExtendedTypes(extendedTypes []JsDocInfoParser_ExtendedTypeInfo) {
	if len(extendedTypes) > 1 {
		// Multiple extends is always ok for an interface.
		isInterface := p.jsdocBuilder.IsInterfaceRecorded()
		// Normally having multiple extended types is an error for a non-interface, but in TSICKLE
		// mode, we instead record there being a single @extends with a placeholder name. This
		// works around some weirdness with TypeScript allowing extending classes.
		isFromTsickle := p.jsDocSourceKind == JsDocInfoParser_JsDocSourceKind_TSICKLE
		if !isInterface && isFromTsickle {
			result := p.jsdocBuilder.RecordBaseType(
				p.CreateJSTypeExpression(
					p.WrapNode(rhino.Token_BANG, p.NewStringNode(JsDocInfoParser_TSICKLE_MISSING_TYPE_PLACEHOLDER)).
						SrcrefTree(extendedTypes[0].type_.GetRoot())))
			jutil.CheckStateMsg(result, "Unexpected failure to record base type")
			return
		}
	}
	for _, typeInfo := range extendedTypes {
		// If interface, record the multiple extended interfaces
		if p.jsdocBuilder.IsInterfaceRecorded() {
			if !p.jsdocBuilder.RecordExtendedInterface(typeInfo.type_) {
				p.AddParserWarningAt(rhino.Msg_JSDOC_EXTENDS_DUPLICATE, typeInfo.lineno, typeInfo.charno)
			}
		} else {
			if !p.jsdocBuilder.RecordBaseType(typeInfo.type_) {
				p.AddTypeWarningAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, typeInfo.lineno, typeInfo.charno)
			}
		}
	}
}

// ParseSuppressTag parses a `@suppress` tag of the form
// `@suppress{warning1|warning2}`.
//
// token: The current token.
//
// Java: JsDocInfoParser.parseSuppressTag(JsDocToken)
func (p *JsDocInfoParser) ParseSuppressTag(token JsDocToken) JsDocToken {
	if token != JsDocToken_LEFT_CURLY {
		p.AddParserWarning(rhino.Msg_JSDOC_SUPPRESS)
		return token
	}
	var suppressBuf [8]string
	suppressions := suppressBuf[:0]
	for {
		p.SkipEOLs()
		if p.Match(JsDocToken_STRING) {
			name := p.stream.GetString()
			if !p.suppressionNames.Contains(name) {
				p.AddParserWarningWithArg(rhino.Msg_JSDOC_SUPPRESS_UNKNOWN, name)
			}

			suppressions = append(suppressions, name)
			token = p.Next()
		} else {
			p.AddParserWarning(rhino.Msg_JSDOC_SUPPRESS)
			return token
		}

		p.SkipEOLs()
		if p.MatchEither(JsDocToken_PIPE, JsDocToken_COMMA) {
			token = p.Next()
		} else {
			break
		}
	}

	token = p.Next()
	if token != JsDocToken_RIGHT_CURLY {
		p.AddParserWarning(rhino.Msg_JSDOC_SUPPRESS)
		return token
	}

	// Find the suppressions' description (if applicable).
	if p.jsdocBuilder.ShouldParseDocumentation() && !p.LookAheadForAnnotation() {
		suppressDescriptionInfo := p.ExtractMultilineTextualBlock(token)
		p.jsdocBuilder.RecordSuppressionSliceWithDescription(suppressions, suppressDescriptionInfo.string_)
		token = suppressDescriptionInfo.token
	} else {
		token = p.EatUntilEOLIfNotAnnotationFrom(p.Next())
		p.jsdocBuilder.RecordSuppressionSliceWithDescription(suppressions, "")
	}
	return token
}

// ParseModsTag parses a `@mods` tag of the form `@mods {google3.path.to.file}`.
//
// token: The current token.
//
// Java: JsDocInfoParser.parseModsTag(JsDocToken)
func (p *JsDocInfoParser) ParseModsTag(token JsDocToken) JsDocToken {
	if token != JsDocToken_LEFT_CURLY {
		p.AddParserWarning(rhino.Msg_JSDOC_MODS)
		return token
	}
	if !p.Match(JsDocToken_STRING) {
		p.AddParserWarning(rhino.Msg_JSDOC_MODS)
		return token
	}
	namespace := p.stream.GetString()
	token = p.Next()
	if !p.Match(JsDocToken_RIGHT_CURLY) {
		p.AddParserWarning(rhino.Msg_JSDOC_MODS)
	} else {
		token = p.Next()
		p.jsdocBuilder.RecordMods(namespace)
	}
	return token
}

// ParseClosurePrimitiveTag parses a `@closurePrimitive` tag
//
// token: The current token.
//
// Java: JsDocInfoParser.parseClosurePrimitiveTag(JsDocToken)
func (p *JsDocInfoParser) ParseClosurePrimitiveTag(token JsDocToken) JsDocToken {
	if token != JsDocToken_LEFT_CURLY {
		p.AddParserWarning(rhino.Msg_JSDOC_MISSING_LC)
		return token
	} else if p.Match(JsDocToken_STRING) {
		name := p.stream.GetString()
		if !p.closurePrimitiveNames.Contains(name) {
			p.AddParserWarningWithArg(rhino.Msg_JSDOC_CLOSUREPRIMITIVE_INVALID, name)
		} else if !p.jsdocBuilder.RecordClosurePrimitiveId(name) {
			p.AddParserWarning(rhino.Msg_JSDOC_CLOSUREPRIMITIVE_EXTRA)
		}
		token = p.Next()
	} else {
		p.AddParserWarning(rhino.Msg_JSDOC_CLOSUREPRIMITIVE_MISSING)
		return token
	}

	if !p.Match(JsDocToken_RIGHT_CURLY) {
		p.AddParserWarning(rhino.Msg_JSDOC_MISSING_RC)
	} else {
		token = p.Next()
	}
	return p.EatUntilEOLIfNotAnnotation()
}

// ParseModifiesTag parses a `@modifies` tag of the form
// `@modifies{this|arguments|param}`.
//
// token: The current token.
//
// Java: JsDocInfoParser.parseModifiesTag(JsDocToken)
func (p *JsDocInfoParser) ParseModifiesTag(token JsDocToken) JsDocToken {
	if token == JsDocToken_LEFT_CURLY {
		var modifiesBuf [4]string
		modifies := modifiesBuf[:0]
		for {
			if p.Match(JsDocToken_STRING) {
				name := p.stream.GetString()
				if !isModifiesAnnotationKeyword(name) && !p.jsdocBuilder.HasParameter(name) {
					p.AddParserWarningWithArg(rhino.Msg_JSDOC_MODIFIES_UNKNOWN, name)
				}

				if !slices.Contains(modifies, name) {
					modifies = append(modifies, name)
				}
				token = p.Next()
			} else {
				p.AddParserWarning(rhino.Msg_JSDOC_MODIFIES)
				return token
			}

			if p.Match(JsDocToken_PIPE) {
				token = p.Next()
			} else {
				break
			}
		}

		if !p.Match(JsDocToken_RIGHT_CURLY) {
			p.AddParserWarning(rhino.Msg_JSDOC_MODIFIES)
		} else {
			token = p.Next()
			if p.jsdocBuilder.IsNoSideEffectsRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES)
			} else if !p.jsdocBuilder.RecordModifiesSlice(modifies) {
				p.AddParserWarning(rhino.Msg_JSDOC_MODIFIES_DUPLICATE)
			}
		}
	}
	return token
}

// ParseIdGeneratorTag parses a `@idgenerator` tag of the form `@idgenerator`
// or `@idgenerator{consistent}`.
//
// token: The current token.
//
// Java: JsDocInfoParser.parseIdGeneratorTag(JsDocToken)
func (p *JsDocInfoParser) ParseIdGeneratorTag(token JsDocToken) JsDocToken {
	idgenKind := "unique"
	if token == JsDocToken_LEFT_CURLY {
		if p.Match(JsDocToken_STRING) {
			name := p.stream.GetString()
			if !isIdGeneratorAnnotationKeyword(name) && !p.jsdocBuilder.HasParameter(name) {
				p.AddParserWarningWithArg(rhino.Msg_JSDOC_IDGEN_UNKNOWN, name)
			}

			idgenKind = name
			token = p.Next()
		} else {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_BAD)
			return token
		}

		if !p.Match(JsDocToken_RIGHT_CURLY) {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_BAD)
		} else {
			token = p.Next()
		}
	}

	switch idgenKind {
	case "unique":
		if !p.jsdocBuilder.RecordIdGenerator() {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_DUPLICATE)
		}
	case "consistent":
		if !p.jsdocBuilder.RecordConsistentIdGenerator() {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_DUPLICATE)
		}
	case "stable":
		if !p.jsdocBuilder.RecordStableIdGenerator() {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_DUPLICATE)
		}
	case "xid":
		if !p.jsdocBuilder.RecordXidGenerator() {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_DUPLICATE)
		}
	case "mapped":
		if !p.jsdocBuilder.RecordMappedIdGenerator() {
			p.AddParserWarning(rhino.Msg_JSDOC_IDGEN_DUPLICATE)
		}
	}

	return token
}
