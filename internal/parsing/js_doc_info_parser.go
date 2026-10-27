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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TODO(nicksantos): Unify all the JSDocInfo stuff into one package, instead of
// spreading it across multiple packages.

// JsDocInfoParser is a parser for JSDoc comments.
type JsDocInfoParser struct {
	stream          *JsDocTokenStream
	jsdocBuilder    *rhino.JSDocInfo_Builder
	jsdocBuilderVal rhino.JSDocInfo_Builder
	errorReporter   rhino.ErrorReporter
	jsDocSourceKind JsDocInfoParser_JsDocSourceKind

	// Use a template node for properties set on all nodes to minimize the
	// memory footprint associated with these (similar to IRFactory).
	templateNode rhino.Node

	// Go-only: performance — where the parser takes its node ids (nil: the
	// store's counter; see IRFactory.alloc).
	alloc *rhino.NodeAllocator

	// The DocInfo with the fileoverview tag for the whole file.
	fileOverviewJSDocInfo *rhino.JSDocInfo
	state                 JsDocInfoParser_State

	annotations           *jutil.LinkedHashMap[string, Annotation]
	suppressionNames      *jutil.LinkedHashSet[string]
	closurePrimitiveNames *jutil.LinkedHashSet[string]
	preserveWhitespace    bool
	onlyParseLicenseJsDoc bool

	licenseText *string

	recordTsType bool

	// One token buffer.
	unreadToken JsDocToken
}

// There's no colon because the colon is parsed as a token before the TTL is
// extracted.
const JsDocInfoParser_TTL_START_DELIMITER = "="
const JsDocInfoParser_TTL_END_DELIMITER = "=:"

// JsDocInfoParser_TEMPLATE_NAME_MATCHER is
// CharMatcher.javaLetterOrDigit().or(CharMatcher.anyOf("_$")).
//
// Go-only: language-adapter — a Guava CharMatcher, as a predicate on a UTF-16
// char.
//
// Java: JsDocInfoParser.TEMPLATE_NAME_MATCHER
var JsDocInfoParser_TEMPLATE_NAME_MATCHER = func(c uint16) bool {
	return jutil.CharacterIsLetter(c) || jutil.CharacterIsDigit(c) || c == '_' || c == '$'
}

const JsDocInfoParser_BAD_TYPE_WIKI_LINK = " See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler" +
	" for more information."

const JsDocInfoParser_TSICKLE_MISSING_TYPE_PLACEHOLDER = "JsDocInfoParser_TsickleMode_MissingSupertypePlaceholder"
const JsDocInfoParser_TS_TYPE_PLACEHOLDER = "JsDocInfoParser_tsType_PlaceHolder"

// Java: JsDocInfoParser.addParserWarning(Msg, String)
func (p *JsDocInfoParser) AddParserWarningWithArg(msg rhino.Msg, messageArg string) {
	p.AddParserWarningWithArgAt(msg, messageArg, p.stream.GetLineno(), p.stream.GetCharno())
}

// Java: JsDocInfoParser.addParserWarning(Msg, String, int, int)
func (p *JsDocInfoParser) AddParserWarningWithArgAt(msg rhino.Msg, messageArg string, lineno int, charno int) {
	p.errorReporter.Warning(msg.FormatArgs(messageArg), p.GetSourceName(), lineno, charno)
}

// Java: JsDocInfoParser.addParserWarning(Msg)
func (p *JsDocInfoParser) AddParserWarning(msg rhino.Msg) {
	p.AddParserWarningAt(msg, p.stream.GetLineno(), p.stream.GetCharno())
}

// Java: JsDocInfoParser.addParserWarning(Msg, int, int)
func (p *JsDocInfoParser) AddParserWarningAt(msg rhino.Msg, lineno int, charno int) {
	p.errorReporter.Warning(msg.Format(), p.GetSourceName(), lineno, charno)
}

// Java: JsDocInfoParser.addTypeWarning(Msg, String)
func (p *JsDocInfoParser) AddTypeWarningWithArg(msg rhino.Msg, messageArg string) {
	p.AddTypeWarningWithArgAt(msg, messageArg, p.stream.GetLineno(), p.stream.GetCharno())
}

// Java: JsDocInfoParser.addTypeWarning(Msg, String, int, int)
func (p *JsDocInfoParser) AddTypeWarningWithArgAt(msg rhino.Msg, messageArg string, lineno int, charno int) {
	p.errorReporter.Warning(
		"Bad type annotation. "+msg.FormatArgs(messageArg)+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		p.GetSourceName(),
		lineno,
		charno)
}

// Java: JsDocInfoParser.addTypeWarning(Msg)
func (p *JsDocInfoParser) AddTypeWarning(msg rhino.Msg) {
	p.AddTypeWarningAt(msg, p.stream.GetLineno(), p.stream.GetCharno())
}

// Java: JsDocInfoParser.addTypeWarning(Msg, int, int)
func (p *JsDocInfoParser) AddTypeWarningAt(msg rhino.Msg, lineno int, charno int) {
	p.errorReporter.Warning(
		"Bad type annotation. "+msg.Format()+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		p.GetSourceName(),
		lineno,
		charno)
}

// Java: JsDocInfoParser.addMissingTypeWarning(int, int)
func (p *JsDocInfoParser) AddMissingTypeWarning(lineno int, charno int) {
	p.errorReporter.Warning("Missing type declaration.", p.GetSourceName(), lineno, charno)
}

// Java: JsDocInfoParser.modifiesAnnotationKeywords
var JsDocInfoParser_modifiesAnnotationKeywords = jutil.NewLinkedHashSet("this", "arguments")

// Java: JsDocInfoParser.idGeneratorAnnotationKeywords
var JsDocInfoParser_idGeneratorAnnotationKeywords = jutil.NewLinkedHashSet("unique", "consistent", "stable", "mapped", "xid")

// Java: JsDocInfoParser.primitiveTypes
var JsDocInfoParser_primitiveTypes = jutil.NewLinkedHashSet("number", "string", "boolean", "symbol")

// Java: JsDocInfoParser.setLicenseTextMonotonic(String)
func (p *JsDocInfoParser) SetLicenseTextMonotonic(x string) {
	// (checkNotNull(x): a Go string is never null.)
	jutil.CheckState(p.licenseText == nil)
	p.licenseText = &x
}

// Java: JsDocInfoParser.getLicenseText()
func (p *JsDocInfoParser) GetLicenseText() *string {
	return p.licenseText
}

// SetFileOverviewJSDocInfo sets the file overview JSDocInfo, in order to warn
// about multiple uses of the @fileoverview tag in a file.
//
// Java: JsDocInfoParser.setFileOverviewJSDocInfo(JSDocInfo)
func (p *JsDocInfoParser) SetFileOverviewJSDocInfo(fileOverviewJSDocInfo *rhino.JSDocInfo) {
	p.fileOverviewJSDocInfo = fileOverviewJSDocInfo
}

// Java: JsDocInfoParser.getSourceFile()
func (p *JsDocInfoParser) GetSourceFile() rhino.StaticSourceFile {
	return p.templateNode.GetStaticSourceFile()
}

// JsDocInfoParser_State is JsDocInfoParser.State.
type JsDocInfoParser_State int32

const (
	JsDocInfoParser_State_SEARCHING_ANNOTATION JsDocInfoParser_State = iota
	JsDocInfoParser_State_SEARCHING_NEWLINE
	JsDocInfoParser_State_NEXT_IS_ANNOTATION
)

var jsDocInfoParser_State_names = [...]string{"SEARCHING_ANNOTATION", "SEARCHING_NEWLINE", "NEXT_IS_ANNOTATION"}

// Name is Enum.name().
//
// Java: JsDocInfoParser.State.name()
func (s JsDocInfoParser_State) Name() string { return jsDocInfoParser_State_names[s] }

// String is Enum.toString().
//
// Java: JsDocInfoParser.State.toString()
func (s JsDocInfoParser_State) String() string { return s.Name() }

// JsDocInfoParser_JsDocSourceKind says how to handle unexpected JSDoc format.
type JsDocInfoParser_JsDocSourceKind int32

const (
	JsDocInfoParser_JsDocSourceKind_NORMAL  JsDocInfoParser_JsDocSourceKind = iota // classic Closure Compiler
	JsDocInfoParser_JsDocSourceKind_TSICKLE                                        // make some extra allowances for tsickle-generated JSDoc
)

var jsDocInfoParser_JsDocSourceKind_names = [...]string{"NORMAL", "TSICKLE"}

// Name is Enum.name().
//
// Java: JsDocInfoParser.JsDocSourceKind.name()
func (k JsDocInfoParser_JsDocSourceKind) Name() string {
	return jsDocInfoParser_JsDocSourceKind_names[k]
}

// String is Enum.toString().
//
// Java: JsDocInfoParser.JsDocSourceKind.toString()
func (k JsDocInfoParser_JsDocSourceKind) String() string { return k.Name() }

// NewJsDocInfoParser constructs a parser for one JSDoc comment.
//
// Go-only: language-adapter — the Java `comment` parameter is nullable and
// guarded by `comment != null`; every Java caller passes a non-null comment,
// and a Go string cannot be null, so the comment is always recorded.
//
// Java: JsDocInfoParser(JsDocTokenStream, String, int, Node, Config, JsDocSourceKind, ErrorReporter)
func NewJsDocInfoParser(
	stream *JsDocTokenStream,
	comment string,
	commentPosition int,
	templateNode rhino.Node,
	config *Config,
	jsDocSourceKind JsDocInfoParser_JsDocSourceKind,
	errorReporter rhino.ErrorReporter,
) *JsDocInfoParser {
	p := &JsDocInfoParser{}
	p.Init(stream, comment, commentPosition, templateNode, config, jsDocSourceKind, errorReporter)
	return p
}

// Init resets the parser for a new JSDoc comment without allocating.
func (p *JsDocInfoParser) Init(
	stream *JsDocTokenStream,
	comment string,
	commentPosition int,
	templateNode rhino.Node,
	config *Config,
	jsDocSourceKind JsDocInfoParser_JsDocSourceKind,
	errorReporter rhino.ErrorReporter,
) {
	p.stream = stream
	p.jsdocBuilder = &p.jsdocBuilderVal
	p.jsdocBuilder.Reset()

	parseDocumentation := config.JsDocParsingMode().ShouldParseDescriptions()
	if parseDocumentation {
		p.jsdocBuilder.ParseDocumentation()
	}
	p.jsdocBuilder.RecordOriginalCommentString(comment)
	p.jsdocBuilder.RecordOriginalCommentPosition(commentPosition)
	p.recordTsType = parseDocumentation
	p.annotations = config.Annotations()
	p.suppressionNames = config.SuppressionNames()
	p.closurePrimitiveNames = config.ClosurePrimitiveNames()
	p.preserveWhitespace = config.JsDocParsingMode().ShouldPreserveWhitespace()
	p.jsDocSourceKind = jsDocSourceKind
	p.onlyParseLicenseJsDoc = config.JsDocParsingMode() == Config_JsDocParsing_LICENSE_COMMENTS_ONLY

	p.errorReporter = errorReporter
	if templateNode.IsNil() {
		p.templateNode = rhino.IR_Script()
	} else {
		p.templateNode = templateNode
	}
	p.fileOverviewJSDocInfo = nil
	p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
	p.licenseText = nil
	p.unreadToken = JsDocInfoParser_NO_UNREAD_TOKEN
}

// GetSourceName returns the source file name.
//
// Go-only: language-adapter — Java returns null when there is no source
// file; here that is "", as in rhino.JSTypeExpression and
// rhino.ErrorReporter, which take the name as a Go string.
//
// Java: JsDocInfoParser.getSourceName()
func (p *JsDocInfoParser) GetSourceName() string {
	sourceFile := p.GetSourceFile()
	if sourceFile == nil {
		return ""
	}
	return sourceFile.GetName()
}

// ParseInlineTypeDoc parses a description as a `@type`.
//
// Java: JsDocInfoParser.parseInlineTypeDoc()
func (p *JsDocInfoParser) ParseInlineTypeDoc() *rhino.JSDocInfo {
	p.SkipEOLs()

	token := p.Next()
	lineno := p.stream.GetLineno()
	startCharno := p.stream.GetCharno()
	typeAst := p.ParseParamTypeExpression(token)
	p.RecordTypeNode(lineno, startCharno, typeAst, token == JsDocToken_LEFT_CURLY)

	expr := p.CreateJSTypeExpression(typeAst)
	if expr != nil {
		p.jsdocBuilder.RecordType(expr)
		p.jsdocBuilder.RecordInlineType()
		return p.RetrieveAndResetParsedJSDocInfo()
	}
	return nil
}

// Java: JsDocInfoParser.recordTypeNode(int, int, Node, boolean)
func (p *JsDocInfoParser) RecordTypeNode(lineno int, startCharno int, typeAst rhino.Node, matchingLC bool) {
	if !typeAst.IsNil() {
		endLineno := p.stream.GetLineno()
		endCharno := p.stream.GetCharno()
		p.jsdocBuilder.MarkTypeNode(typeAst, lineno, startCharno, endLineno, endCharno, matchingLC)
	}
}

// JsDocInfoParser_ParseTypeString parses a string containing a JsDoc type
// declaration, returning the type if the parsing succeeded or nil if it
// failed.
//
// Java: JsDocInfoParser.parseTypeString(String)
func JsDocInfoParser_ParseTypeString(typeString string) rhino.Node {
	parser := JsDocInfoParser_GetParser(typeString)
	return parser.ParseTopLevelTypeExpression(parser.Next())
}

// JsDocInfoParser_ParseJsdoc parses a string containing a JsDoc declaration,
// returning the entire JSDocInfo if the parsing succeeded or nil if it
// failed.
//
// Java: JsDocInfoParser.parseJsdoc(String)
func JsDocInfoParser_ParseJsdoc(toParse string) *rhino.JSDocInfo {
	parser := JsDocInfoParser_GetParser(toParse)
	parser.Parse()
	return parser.RetrieveAndResetParsedJSDocInfo()
}

// Java: JsDocInfoParser.parseFileOverviewJsdoc(String)
func JsDocInfoParser_ParseFileOverviewJsdoc(toParse string) *rhino.JSDocInfo {
	parser := JsDocInfoParser_GetParser(toParse)
	parser.Parse()
	return parser.GetFileOverviewJSDocInfo()
}

// Java: JsDocInfoParser.getParser(String)
func JsDocInfoParser_GetParser(toParse string) *JsDocInfoParser {
	config := Config_Builder_().
		SetLanguageMode(Config_LanguageMode_ECMASCRIPT3).
		SetStrictMode(Config_StrictMode_SLOPPY).
		SetClosurePrimitiveNames(jutil.NewLinkedHashSet("testPrimitive").All()).
		Build()
	return NewJsDocInfoParser(
		NewJsDocTokenStream(toParse),
		toParse,
		0,
		rhino.NilNode,
		config,
		JsDocInfoParser_JsDocSourceKind_NORMAL,
		rhino.ErrorReporter_NULL_INSTANCE)
}

// Parse parses a JSDocInfo object. This parsing method reads all tokens
// returned by the JsDocTokenStream.GetJsDocToken() method until the
// JsDocToken_EOC is returned.
//
// Returns true if JSDoc information was correctly parsed, false otherwise.
//
// Java: JsDocInfoParser.parse()
func (p *JsDocInfoParser) Parse() bool {
	p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
	p.SkipEOLs()

	token := p.Next()

	// Always record that we have a comment.
	if p.jsdocBuilder.ShouldParseDocumentation() {
		blockInfo := p.ExtractBlockComment(token)
		token = blockInfo.token
		if blockInfo.string_ != "" {
			p.jsdocBuilder.RecordBlockDescription(blockInfo.string_)
		}
	} else {
		if token != JsDocToken_ANNOTATION && token != JsDocToken_EOC {
			// Mark that there was a description, but don't bother marking
			// what it was.
			p.jsdocBuilder.RecordBlockDescription("")
		}
	}

	var extendedTypes []*JsDocInfoParser_ExtendedTypeInfo
	return p.ParseHelperLoop(token, &extendedTypes)
}

// ParseImportantComment parses an important comment. Important comments begin
// with /*! They are treated as license blocks, but no further JSDoc parsing is
// performed
//
// Java: JsDocInfoParser.parseImportantComment()
func (p *JsDocInfoParser) ParseImportantComment() {
	p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
	p.SkipEOLs()

	token := p.Next()

	info := p.ExtractMultilineComment(token, JsDocInfoParser_WhitespaceOption_PRESERVE, false, true)

	// An extra space is added by the @license annotation
	// so we need to add one here so they will be identical
	license := " " + info.string_

	p.SetLicenseTextMonotonic(license)
	if p.jsdocBuilder.ShouldParseDocumentation() {
		p.jsdocBuilder.RecordBlockDescription(license)
	} else {
		p.jsdocBuilder.RecordBlockDescription("")
	}
}

// Java: JsDocInfoParser.parseHelperLoop(JsDocToken, List<ExtendedTypeInfo>)
func (p *JsDocInfoParser) ParseHelperLoop(token JsDocToken, extendedTypes *[]*JsDocInfoParser_ExtendedTypeInfo) bool {
	for {
		switch token {
		case JsDocToken_ANNOTATION:
			if p.state == JsDocInfoParser_State_SEARCHING_ANNOTATION {
				p.state = JsDocInfoParser_State_SEARCHING_NEWLINE
				token = p.ParseAnnotation(token, extendedTypes)
			} else {
				token = p.Next()
			}
		case JsDocToken_EOC:
			success := true
			// TODO(johnlenz): It should be a parse error to have an @extends
			// or similiar annotations in a file overview block.
			p.CheckExtendedTypes(*extendedTypes)
			if p.HasParsedFileOverviewDocInfo() {
				p.fileOverviewJSDocInfo = p.RetrieveAndResetParsedJSDocInfo()
				visibility := p.fileOverviewJSDocInfo.GetVisibility()
				switch visibility {
				case rhino.JSDocInfo_Visibility_PRIVATE, // fallthrough
					rhino.JSDocInfo_Visibility_PROTECTED:
					// PRIVATE and PROTECTED are not allowed in @fileoverview JsDoc.
					p.AddParserWarningWithArg(
						rhino.Msg_BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION,
						jutil.StringToLowerCaseASCII(visibility.String()))
					success = false
				default:
					// PACKAGE, PUBLIC, and (implicitly) INHERITED are allowed
					// in @fileoverview JsDoc.
				}
			}
			return success
		case JsDocToken_EOF:
			// discard any accumulated information
			p.jsdocBuilder.Build()
			p.AddParserWarning(rhino.Msg_UNEXPECTED_EOF)
			p.CheckExtendedTypes(*extendedTypes)
			return false
		case JsDocToken_EOL:
			if p.state == JsDocInfoParser_State_SEARCHING_NEWLINE {
				p.state = JsDocInfoParser_State_SEARCHING_ANNOTATION
			}
			token = p.Next()
		default:
			if token == JsDocToken_STAR && p.state == JsDocInfoParser_State_SEARCHING_ANNOTATION {
				token = p.Next()
			} else {
				p.state = JsDocInfoParser_State_SEARCHING_NEWLINE
				token = p.EatTokensUntilEOL()
			}
		}
	}
}

// Java: JsDocInfoParser.parseAnnotation(JsDocToken, List<ExtendedTypeInfo>)
func (p *JsDocInfoParser) ParseAnnotation(token JsDocToken, extendedTypes *[]*JsDocInfoParser_ExtendedTypeInfo) JsDocToken {
	// JSTypes are represented as Rhino AST nodes, and then resolved later.
	var type_ *rhino.JSTypeExpression
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()

	annotationName := p.stream.GetString()
	annotation, annotationFound := p.annotations.Lookup(annotationName)

	if p.onlyParseLicenseJsDoc &&
		!(annotationFound && (annotation == Annotation_LICENSE || annotation == Annotation_PRESERVE)) {
		return p.Next()
	}

	if !annotationFound || annotationName == "" {
		p.AddParserWarningWithArg(rhino.Msg_BAD_JSDOC_TAG, annotationName)
	} else {
		// Mark the beginning of the annotation.
		p.jsdocBuilder.MarkAnnotation(annotationName, lineno, charno)

		var matchingRc bool
		switch annotation {
		case Annotation_NG_INJECT:
			if p.jsdocBuilder.IsNgInjectRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_NGINJECT_EXTRA)
			} else {
				p.jsdocBuilder.RecordNgInject(true)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_ABSTRACT:
			if !p.jsdocBuilder.RecordAbstract() {
				p.AddTypeWarning(rhino.Msg_JSDOC_INCOMPAT_TYPE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_AUTHOR:
			if p.jsdocBuilder.ShouldParseDocumentation() {
				authorInfo := p.ExtractSingleLineBlock()
				author := authorInfo.string_

				if author == "" {
					p.AddParserWarning(rhino.Msg_JSDOC_AUTHORMISSING)
				} else {
					p.jsdocBuilder.RecordAuthor(author)
				}
				token = authorInfo.token
			} else {
				token = p.EatUntilEOLIfNotAnnotation()
			}
			return token
		case Annotation_UNRESTRICTED:
			if !p.jsdocBuilder.RecordUnrestricted() {
				p.AddTypeWarning(rhino.Msg_JSDOC_INCOMPAT_TYPE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_STRUCT:
			if !p.jsdocBuilder.RecordStruct() {
				p.AddTypeWarning(rhino.Msg_JSDOC_INCOMPAT_TYPE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_DICT:
			if !p.jsdocBuilder.RecordDict() {
				p.AddTypeWarning(rhino.Msg_JSDOC_INCOMPAT_TYPE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE:
			if !p.jsdocBuilder.RecordCollapsibleOrBreakMyCode() {
				p.AddParserWarning(rhino.Msg_JSDOC_COLLAPSIBLEORBREAKMYCODE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_CONSTRUCTOR:
			if !p.jsdocBuilder.RecordConstructor() {
				if p.jsdocBuilder.IsInterfaceRecorded() {
					p.AddTypeWarning(rhino.Msg_JSDOC_INTERFACE_CONSTRUCTOR)
				} else {
					p.AddTypeWarning(rhino.Msg_JSDOC_INCOMPAT_TYPE)
				}
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_RECORD:
			if !p.jsdocBuilder.RecordImplicitMatch() {
				p.AddTypeWarning(rhino.Msg_JSDOC_RECORD)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_DEPRECATED:
			if !p.jsdocBuilder.RecordDeprecated() {
				p.AddParserWarning(rhino.Msg_JSDOC_DEPRECATED)
			}

			// Find the reason/description, if any.
			reasonInfo := p.ExtractMultilineTextualBlock(token)

			reason := reasonInfo.string_

			if len(reason) > 0 {
				p.jsdocBuilder.RecordDeprecationReason(reason)
			}

			token = reasonInfo.token
			return token
		case Annotation_INTERFACE:
			if !p.jsdocBuilder.RecordInterface() {
				if p.jsdocBuilder.IsConstructorRecorded() {
					p.AddTypeWarning(rhino.Msg_JSDOC_INTERFACE_CONSTRUCTOR)
				} else {
					p.AddTypeWarning(rhino.Msg_JSDOC_INCOMPAT_TYPE)
				}
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_DESC:
			if p.jsdocBuilder.IsDescriptionRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_DESC_EXTRA)
				return p.EatUntilEOLIfNotAnnotation()
			} else {
				descriptionInfo := p.ExtractMultilineTextualBlock(token)

				description := descriptionInfo.string_

				p.jsdocBuilder.RecordDescription(description)
				token = descriptionInfo.token
				return token
			}
		case Annotation_TS_TYPE:
			tsTypeInfo := p.ExtractMultilineTextualBlock(token)

			tsType := tsTypeInfo.string_

			if p.recordTsType {
				p.jsdocBuilder.RecordTsType(tsType)
			} else {
				p.jsdocBuilder.RecordTsType(JsDocInfoParser_TS_TYPE_PLACEHOLDER)
			}
			token = tsTypeInfo.token
			return token
		case Annotation_FILE_OVERVIEW:
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
		case Annotation_LICENSE, Annotation_PRESERVE:
			// Always use PRESERVE for @license and @preserve blocks.
			preserveInfo := p.ExtractMultilineTextualBlockWithOption(token, JsDocInfoParser_WhitespaceOption_PRESERVE, true)

			preserve := preserveInfo.string_

			if len(preserve) > 0 {
				p.SetLicenseTextMonotonic(preserve)
			}

			token = preserveInfo.token
			return token
		case Annotation_ENHANCE:
			enhanceInfo := p.ExtractSingleLineBlock()
			enhance := enhanceInfo.string_
			token = enhanceInfo.token

			p.jsdocBuilder.RecordEnhance(enhance)
			return token
		case Annotation_MODS:
			if p.jsdocBuilder.IsModsRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_MODS_EXTRA)
			} else {
				token = p.ParseModsTag(p.Next())
			}
			if token != JsDocToken_EOC && token != JsDocToken_EOF {
				token = p.EatUntilEOLIfNotAnnotation()
			}
			return token
		case Annotation_ENUM:
			token = p.Next()
			lineno = p.stream.GetLineno()
			charno = p.stream.GetCharno()

			type_ = nil
			if token != JsDocToken_EOL && token != JsDocToken_EOC {
				typeNode := p.ParseAndRecordTypeNode(token)
				if !typeNode.IsNil() && typeNode.IsStringLit() {
					typeName := typeNode.GetString()
					if !JsDocInfoParser_primitiveTypes.Contains(typeName) {
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
		case Annotation_EXTERNS:
			if !p.jsdocBuilder.RecordExterns() {
				p.AddParserWarning(rhino.Msg_JSDOC_EXTERNS)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_NO_COVERAGE:
			if !p.jsdocBuilder.RecordNoCoverage() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOCOVERAGE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_TYPE_SUMMARY:
			if !p.jsdocBuilder.RecordTypeSummary() {
				p.AddParserWarning(rhino.Msg_JSDOC_TYPESUMMARY)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_EXTENDS, Annotation_IMPLEMENTS:
			p.SkipEOLs()
			token = p.Next()
			lineno = p.stream.GetLineno()
			charno = p.stream.GetCharno()
			matchingRc = false
			if token == JsDocToken_LEFT_CURLY {
				token = p.Next()
				matchingRc = true
			}

			if token == JsDocToken_STRING {
				typeNode := p.ParseAndRecordTypeNameNode(token, lineno, charno, matchingRc)

				lineno = p.stream.GetLineno()
				charno = p.stream.GetCharno()

				typeNode = p.WrapNode(rhino.Token_BANG, typeNode)
				type_ = p.CreateJSTypeExpression(typeNode)

				if annotation == Annotation_EXTENDS {
					// record the extended type, check later
					*extendedTypes = append(*extendedTypes, NewJsDocInfoParser_ExtendedTypeInfo(type_, p.stream.GetLineno(), p.stream.GetCharno()))
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
			token = p.EatUntilEOLIfNotAnnotationFrom(token)
			return token
		case Annotation_LENDS:
			p.SkipEOLs()

			matchingRc = false
			if p.Match(JsDocToken_LEFT_CURLY) {
				token = p.Next()
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
		case Annotation_MEANING:
			meaningInfo := p.ExtractMultilineTextualBlock(token)
			meaning := meaningInfo.string_
			token = meaningInfo.token
			if !p.jsdocBuilder.RecordMeaning(meaning) {
				p.AddParserWarning(rhino.Msg_JSDOC_MEANING_EXTRA)
			}
			return token
		case Annotation_ALTERNATE_MESSAGE_ID:
			alternateMessageIdInfo := p.ExtractSingleLineBlock()
			alternateMessageId := alternateMessageIdInfo.string_
			token = alternateMessageIdInfo.token
			if !p.jsdocBuilder.RecordAlternateMessageId(alternateMessageId) {
				p.AddParserWarning(rhino.Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA)
			}
			return token
		case Annotation_CLOSURE_PRIMITIVE:
			p.SkipEOLs()
			return p.ParseClosurePrimitiveTag(p.Next())
		case Annotation_NO_COMPILE:
			if !p.jsdocBuilder.RecordNoCompile() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOCOMPILE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_NO_DTS:
			if !p.jsdocBuilder.RecordNoDts() {
				p.AddParserWarning(rhino.Msg_JSDOC_NODTS)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_NO_COLLAPSE:
			if !p.jsdocBuilder.RecordNoCollapse() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOCOLLAPSE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_ENCOURAGE_INLINING:
			if !p.jsdocBuilder.RecordEncourageInlining() {
				p.AddParserWarning(rhino.Msg_JSDOC_ENCOURAGE_INLINING)
			}
			if p.jsdocBuilder.IsRequireInlining() {
				p.AddParserWarning(rhino.Msg_JSDOC_INCOMPAT_INLINING)
			}
			if p.jsdocBuilder.IsNoInline() {
				p.AddParserWarning(rhino.Msg_JSDOC_INCOMPAT_INLINING)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_REQUIRE_INLINING:
			if !p.jsdocBuilder.RecordRequireInlining() {
				p.AddParserWarning(rhino.Msg_JSDOC_REQUIRE_INLINING)
			}
			if p.jsdocBuilder.IsEncourageInlining() {
				p.AddParserWarning(rhino.Msg_JSDOC_INCOMPAT_INLINING)
			}
			if p.jsdocBuilder.IsNoInline() {
				p.AddParserWarning(rhino.Msg_JSDOC_INCOMPAT_INLINING)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_NO_INLINE:
			if !p.jsdocBuilder.RecordNoInline() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOINLINE)
			}
			if p.jsdocBuilder.IsEncourageInlining() {
				p.AddParserWarning(rhino.Msg_JSDOC_INCOMPAT_INLINING)
			}
			if p.jsdocBuilder.IsRequireInlining() {
				p.AddParserWarning(rhino.Msg_JSDOC_INCOMPAT_INLINING)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_PROVIDE_GOOG:
			if !p.jsdocBuilder.RecordProvideGoog() {
				p.AddParserWarning(rhino.Msg_JSDOC_PROVIDE_GOOG)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_PROVIDE_ALREADY_PROVIDED:
			if !p.jsdocBuilder.RecordProvideAlreadyProvided() {
				p.AddParserWarning(rhino.Msg_JSDOC_PROVIDE_ALREADY_PROVIDED)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_PURE_OR_BREAK_MY_CODE:
			if !p.jsdocBuilder.RecordPureOrBreakMyCode() {
				p.AddParserWarning(rhino.Msg_JSDOC_PUREORBREAKMYCODE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_INHERIT_DOC, Annotation_OVERRIDE:
			if !p.jsdocBuilder.RecordOverride() {
				p.AddTypeWarning(rhino.Msg_JSDOC_OVERRIDE)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_POLYMER:
			if p.jsdocBuilder.IsPolymerRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_POLYMER_EXTRA)
			} else {
				p.jsdocBuilder.RecordPolymer()
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_POLYMER_BEHAVIOR:
			if p.jsdocBuilder.IsPolymerBehaviorRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_POLYMERBEHAVIOR_EXTRA)
			} else {
				p.jsdocBuilder.RecordPolymerBehavior()
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_CUSTOM_ELEMENT:
			if p.jsdocBuilder.IsCustomElementRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_CUSTOMELEMENT_EXTRA)
			} else {
				p.jsdocBuilder.RecordCustomElement()
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_MIXIN_CLASS:
			if p.jsdocBuilder.IsMixinClassRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_MIXINCLASS_EXTRA)
			} else {
				p.jsdocBuilder.RecordMixinClass()
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_MIXIN_FUNCTION:
			if p.jsdocBuilder.IsMixinFunctionRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_MIXINFUNCTION_EXTRA)
			} else {
				p.jsdocBuilder.RecordMixinFunction()
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_SASS_GENERATED_CSS_TS:
			if p.jsdocBuilder.IsSassGeneratedCssTsRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_SASS_GENERATED_CSS_TS)
			} else {
				p.jsdocBuilder.RecordSassGeneratedCssTs()
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_CLOSURE_UNAWARE_CODE:
			return p.ParseClosureUnawareCode()
		case Annotation_THROWS:
			lineno = p.stream.GetLineno()
			charno = p.stream.GetCharno()
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
		case Annotation_PARAM:
			p.SkipEOLs()
			token = p.Next()
			lineno = p.stream.GetLineno()
			charno = p.stream.GetCharno()
			type_ = nil

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

			var name *string
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

				nameString := p.stream.GetString()
				name = &nameString

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
				if !rhino.TokenStream_IsJSIdentifier(*name) {
					p.AddParserWarningWithArgAt(rhino.Msg_INVALID_VARIABLE_NAME, *name, lineno, charno)
					name = nil
				} else if !p.jsdocBuilder.RecordParameter(*name, type_) {
					if p.jsdocBuilder.HasParameter(*name) {
						p.AddTypeWarningWithArgAt(rhino.Msg_DUP_VARIABLE_NAME, *name, lineno, charno)
					} else {
						p.AddTypeWarningWithArgAt(rhino.Msg_JSDOC_INCOMPAT_TYPE, *name, lineno, charno)
					}
				}
			}

			if name == nil {
				token = p.EatUntilEOLIfNotAnnotationFrom(token)
				return token
			}

			p.jsdocBuilder.MarkName(*name, p.templateNode, lineno, charno)

			// Find the parameter's description (if applicable).
			if p.jsdocBuilder.ShouldParseDocumentation() && token != JsDocToken_ANNOTATION {
				paramDescriptionInfo := p.ExtractMultilineTextualBlock(token)

				paramDescription := paramDescriptionInfo.string_

				if len(paramDescription) > 0 {
					p.jsdocBuilder.RecordParameterDescription(*name, paramDescription)
				}

				token = paramDescriptionInfo.token
			} else if token != JsDocToken_EOC && token != JsDocToken_EOF {
				token = p.EatUntilEOLIfNotAnnotation()
			}
			return token
		case Annotation_NO_SIDE_EFFECTS:
			if p.jsdocBuilder.IsModifiesRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES)
			} else if p.jsdocBuilder.IsThrowsRecorded() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS)
			} else if !p.jsdocBuilder.RecordNoSideEffects() {
				p.AddParserWarning(rhino.Msg_JSDOC_NOSIDEEFFECTS)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_MODIFIES:
			token = p.ParseModifiesTag(p.Next())
			return token
		case Annotation_IMPLICIT_CAST:
			if !p.jsdocBuilder.RecordImplicitCast() {
				p.AddTypeWarning(rhino.Msg_JSDOC_IMPLICITCAST)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_SEE:
			if p.jsdocBuilder.ShouldParseDocumentation() {
				referenceInfo := p.ExtractSingleLineBlock()
				reference := referenceInfo.string_

				if reference == "" {
					p.AddParserWarning(rhino.Msg_JSDOC_SEEMISSING)
				} else {
					p.jsdocBuilder.RecordReference(reference)
				}

				token = referenceInfo.token
			} else {
				token = p.EatUntilEOLIfNotAnnotation()
			}
			return token
		case Annotation_SUPPRESS:
			token = p.ParseSuppressTag(p.Next())
			return token
		case Annotation_TEMPLATE:
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
			var templateNames []string
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
		case Annotation_IDGENERATOR:
			token = p.ParseIdGeneratorTag(p.Next())
			return token
		case Annotation_LOG_TYPE_IN_COMPILER:
			_ = p.jsdocBuilder.RecordLogTypeInCompiler()
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_NOT_IMPLEMENTED,
			Annotation_JSX,
			Annotation_JSX_FRAGMENT,
			Annotation_SOY_MODULE,
			Annotation_SOY_TEMPLATE,
			Annotation_WIZ_ANALYZER,
			Annotation_MAY_HAVE_EXTRA_EDGE:
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_USED_VIA_DOT_CONSTRUCTOR:
			if !p.jsdocBuilder.RecordUsedViaDotConstructor() {
				p.AddParserWarning(rhino.Msg_JSDOC_USEDVIADOTCONSTRUCTOR)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_WIZACTION:
			if !p.jsdocBuilder.RecordWizaction() {
				p.AddParserWarning(rhino.Msg_JSDOC_WIZACTION)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_WIZCALLBACK:
			if !p.jsdocBuilder.RecordWizcallback() {
				p.AddParserWarning(rhino.Msg_JSDOC_WIZCALLBACK)
			}
			return p.EatUntilEOLIfNotAnnotation()
		case Annotation_CONSTANT,
			Annotation_FINAL,
			Annotation_DEFINE,
			Annotation_EXPORT,
			Annotation_RETURN,
			Annotation_PACKAGE,
			Annotation_PRIVATE,
			Annotation_PROTECTED,
			Annotation_PUBLIC,
			Annotation_THIS,
			Annotation_TYPE,
			Annotation_TYPEDEF:
			lineno = p.stream.GetLineno()
			charno = p.stream.GetCharno()

			var typeNode rhino.Node
			hasType := p.LookAheadForType()
			isAlternateTypeAnnotation := annotation == Annotation_PACKAGE ||
				annotation == Annotation_PRIVATE ||
				annotation == Annotation_PROTECTED ||
				annotation == Annotation_PUBLIC ||
				annotation == Annotation_CONSTANT ||
				annotation == Annotation_FINAL ||
				annotation == Annotation_EXPORT
			canSkipTypeAnnotation := isAlternateTypeAnnotation || annotation == Annotation_RETURN
			type_ = nil

			if annotation == Annotation_RETURN && !hasType {
				p.AddMissingTypeWarning(p.stream.GetLineno(), p.stream.GetCharno())
			}

			if hasType || !canSkipTypeAnnotation {
				p.SkipEOLs()
				token = p.Next()
				typeNode = p.ParseAndRecordTypeNode(token)

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

						returnDescription := returnDescriptionInfo.string_

						if len(returnDescription) > 0 {
							p.jsdocBuilder.RecordReturnDescription(returnDescription)
						}

						token = returnDescriptionInfo.token
					} else {
						token = p.EatUntilEOLIfNotAnnotation()
					}
					return token

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
	}

	return p.Next()
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
	for _, c := range jutil.StringToUTF16(name) {
		if !JsDocInfoParser_TEMPLATE_NAME_MATCHER(c) {
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
func (p *JsDocInfoParser) CheckExtendedTypes(extendedTypes []*JsDocInfoParser_ExtendedTypeInfo) {
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
	} else {
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
			suppressDescription := suppressDescriptionInfo.string_
			p.jsdocBuilder.RecordSuppressionSliceWithDescription(suppressions, suppressDescription)
			token = suppressDescriptionInfo.token
		} else {
			token = p.EatUntilEOLIfNotAnnotationFrom(p.Next())
			p.jsdocBuilder.RecordSuppressionSliceWithDescription(suppressions, "")
		}
		return token
	}
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
		modifies := jutil.NewLinkedHashSet[string]()
		for {
			if p.Match(JsDocToken_STRING) {
				name := p.stream.GetString()
				if !JsDocInfoParser_modifiesAnnotationKeywords.Contains(name) && !p.jsdocBuilder.HasParameter(name) {
					p.AddParserWarningWithArg(rhino.Msg_JSDOC_MODIFIES_UNKNOWN, name)
				}

				modifies.Add(p.stream.GetString())
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
			} else if !p.jsdocBuilder.RecordModifies(modifies) {
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
			if !JsDocInfoParser_idGeneratorAnnotationKeywords.Contains(name) && !p.jsdocBuilder.HasParameter(name) {
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

// ParseAndRecordTypeNode looks for a type expression at the current token and
// if found, returns it. Note that this method consumes input.
//
// token: The current token. Returns the type expression found or nil if none.
//
// Java: JsDocInfoParser.parseAndRecordTypeNode(JsDocToken)
func (p *JsDocInfoParser) ParseAndRecordTypeNode(token JsDocToken) rhino.Node {
	return p.ParseAndRecordTypeNodeWithPosition(
		token, p.stream.GetLineno(), p.stream.GetCharno(), token == JsDocToken_LEFT_CURLY, false)
}

// ParseAndRecordTypeNodeWithPosition looks for a parameter type expression at
// the current token and if found, returns it. Note that this method consumes
// input.
//
// token: The current token. lineno: The line of the type expression.
// startCharno: The starting character position of the type expression.
// matchingLC: Whether the type expression starts with a "{".
// onlyParseSimpleNames: If true, only simple type names are parsed (via a
// call to parseTypeNameAnnotation instead of parseTypeExpressionAnnotation).
// Returns the type expression found or nil if none.
//
// Java: JsDocInfoParser.parseAndRecordTypeNode(JsDocToken, int, int, boolean, boolean)
func (p *JsDocInfoParser) ParseAndRecordTypeNodeWithPosition(
	token JsDocToken,
	lineno int,
	startCharno int,
	matchingLC bool,
	onlyParseSimpleNames bool,
) rhino.Node {
	var typeNode rhino.Node

	if onlyParseSimpleNames {
		typeNode = p.ParseTypeNameAnnotation(token)
	} else {
		typeNode = p.ParseTypeExpressionAnnotation(token)
	}

	p.RecordTypeNode(lineno, startCharno, typeNode, matchingLC)
	return typeNode
}

// ParseAndRecordTypeNameNode looks for a type expression at the current token
// and if found, returns it. Note that this method consumes input.
//
// token: The current token. lineno: The line of the type expression.
// startCharno: The starting character position of the type expression.
// matchingLC: Whether the type expression starts with a "{". Returns the
// type expression found or nil if none.
//
// Java: JsDocInfoParser.parseAndRecordTypeNameNode(JsDocToken, int, int, boolean)
func (p *JsDocInfoParser) ParseAndRecordTypeNameNode(
	token JsDocToken, lineno int, startCharno int, matchingLC bool) rhino.Node {
	return p.ParseAndRecordTypeNodeWithPosition(token, lineno, startCharno, matchingLC, true)
}

// ParseAndRecordParamTypeNode looks for a type expression at the current
// token and if found, returns it. Note that this method consumes input.
//
// Parameter type expressions are special for two reasons:
//
//  1. They must begin with '{', to distinguish type names from param names.
//  2. They may end in '=', to denote optionality.
//
// token: The current token. Returns the type expression found or nil if none.
//
// Java: JsDocInfoParser.parseAndRecordParamTypeNode(JsDocToken)
func (p *JsDocInfoParser) ParseAndRecordParamTypeNode(token JsDocToken) rhino.Node {
	jutil.CheckArgument(token == JsDocToken_LEFT_CURLY)
	lineno := p.stream.GetLineno()
	startCharno := p.stream.GetCharno()

	typeNode := p.ParseParamTypeExpressionAnnotation(token)
	p.RecordTypeNode(lineno, startCharno, typeNode, true)
	return typeNode
}

// ToStringToken converts a JSDoc token to its string representation.
//
// Java: JsDocInfoParser.toString(JsDocToken)
func (p *JsDocInfoParser) ToStringToken(token JsDocToken) string {
	switch token {
	case JsDocToken_ANNOTATION:
		return "@" + p.stream.GetString()
	case JsDocToken_BANG:
		return "!"
	case JsDocToken_COMMA:
		return ","
	case JsDocToken_COLON:
		return ":"
	case JsDocToken_RIGHT_ANGLE:
		return ">"
	case JsDocToken_LEFT_SQUARE:
		return "["
	case JsDocToken_LEFT_CURLY:
		return "{"
	case JsDocToken_LEFT_PAREN:
		return "("
	case JsDocToken_LEFT_ANGLE:
		return "<"
	case JsDocToken_QMARK:
		return "?"
	case JsDocToken_PIPE:
		return "|"
	case JsDocToken_RIGHT_SQUARE:
		return "]"
	case JsDocToken_RIGHT_CURLY:
		return "}"
	case JsDocToken_RIGHT_PAREN:
		return ")"
	case JsDocToken_STAR:
		return "*"
	case JsDocToken_ITER_REST:
		return "..."
	case JsDocToken_EQUALS:
		return "="
	case JsDocToken_STRING:
		return p.stream.GetString()
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// CreateJSTypeExpression constructs a new JSTypeExpression.
//
// n: A node. May be nil.
//
// Java: JsDocInfoParser.createJSTypeExpression(Node)
func (p *JsDocInfoParser) CreateJSTypeExpression(n rhino.Node) *rhino.JSTypeExpression {
	if n.IsNil() {
		return nil
	}
	return rhino.NewJSTypeExpression(n, p.GetSourceName())
}
