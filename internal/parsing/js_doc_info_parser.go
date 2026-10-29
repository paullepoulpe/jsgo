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

	config                *Config
	suppressionNames      map[string]struct{}
	closurePrimitiveNames map[string]struct{}
	preserveWhitespace    bool
	onlyParseLicenseJsDoc bool

	licenseText    string
	hasLicenseText bool

	recordTsType bool

	// One token buffer.
	unreadToken JsDocToken
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

// Java: JsDocInfoParser.setLicenseTextMonotonic(String)
func (p *JsDocInfoParser) SetLicenseTextMonotonic(x string) {
	// (checkNotNull(x): a Go string is never null.)
	jutil.CheckState(!p.hasLicenseText)
	p.licenseText = x
	p.hasLicenseText = true
}

// Java: JsDocInfoParser.getLicenseText()
func (p *JsDocInfoParser) GetLicenseText() (string, bool) {
	return p.licenseText, p.hasLicenseText
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
	p.config = config
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
	p.licenseText = ""
	p.hasLicenseText = false
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

var jsDocInfoParser_defaultConfig = Config_Builder_().
	SetLanguageMode(Config_LanguageMode_ECMASCRIPT3).
	SetStrictMode(Config_StrictMode_SLOPPY).
	SetClosurePrimitiveNames(slices.Values([]string{"testPrimitive"})).
	Build()

// Java: JsDocInfoParser.getParser(String)
func JsDocInfoParser_GetParser(toParse string) *JsDocInfoParser {
	return NewJsDocInfoParser(
		NewJsDocTokenStream(toParse),
		toParse,
		0,
		rhino.NilNode,
		jsDocInfoParser_defaultConfig,
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

	var extendedTypesBuf [2]JsDocInfoParser_ExtendedTypeInfo
	extendedTypes := extendedTypesBuf[:0]
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
func (p *JsDocInfoParser) ParseHelperLoop(token JsDocToken, extendedTypes *[]JsDocInfoParser_ExtendedTypeInfo) bool {
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
func (p *JsDocInfoParser) ParseAnnotation(token JsDocToken, extendedTypes *[]JsDocInfoParser_ExtendedTypeInfo) JsDocToken {
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()

	annotationName := p.stream.GetString()
	annotation, annotationFound := p.config.LookupAnnotation(annotationName)

	if p.onlyParseLicenseJsDoc &&
		!(annotationFound && (annotation == Annotation_LICENSE || annotation == Annotation_PRESERVE)) {
		return p.Next()
	}

	if !annotationFound || annotationName == "" {
		p.AddParserWarningWithArg(rhino.Msg_BAD_JSDOC_TAG, annotationName)
		return p.Next()
	}

	// Mark the beginning of the annotation.
	p.jsdocBuilder.MarkAnnotation(annotationName, lineno, charno)

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
		return p.parseAuthorTag()
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
		return p.parseDeprecatedTag(token)
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
		return p.parseDescTag(token)
	case Annotation_TS_TYPE:
		return p.parseTsTypeTag(token)
	case Annotation_FILE_OVERVIEW:
		return p.parseFileOverviewTag(token)
	case Annotation_LICENSE, Annotation_PRESERVE:
		return p.parseLicenseOrPreserveTag(token)
	case Annotation_ENHANCE:
		enhanceInfo := p.ExtractSingleLineBlock()
		p.jsdocBuilder.RecordEnhance(enhanceInfo.string_)
		return enhanceInfo.token
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
		return p.parseEnumTag()
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
		return p.parseExtendsOrImplementsTag(annotation, extendedTypes)
	case Annotation_LENDS:
		return p.parseLendsTag()
	case Annotation_MEANING:
		meaningInfo := p.ExtractMultilineTextualBlock(token)
		if !p.jsdocBuilder.RecordMeaning(meaningInfo.string_) {
			p.AddParserWarning(rhino.Msg_JSDOC_MEANING_EXTRA)
		}
		return meaningInfo.token
	case Annotation_ALTERNATE_MESSAGE_ID:
		alternateMessageIdInfo := p.ExtractSingleLineBlock()
		if !p.jsdocBuilder.RecordAlternateMessageId(alternateMessageIdInfo.string_) {
			p.AddParserWarning(rhino.Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA)
		}
		return alternateMessageIdInfo.token
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
		return p.parseThrowsTag(token)
	case Annotation_PARAM:
		return p.parseParamTag()
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
		return p.ParseModifiesTag(p.Next())
	case Annotation_IMPLICIT_CAST:
		if !p.jsdocBuilder.RecordImplicitCast() {
			p.AddTypeWarning(rhino.Msg_JSDOC_IMPLICITCAST)
		}
		return p.EatUntilEOLIfNotAnnotation()
	case Annotation_SEE:
		return p.parseSeeTag()
	case Annotation_SUPPRESS:
		return p.ParseSuppressTag(p.Next())
	case Annotation_TEMPLATE:
		return p.parseTemplateTag(lineno, charno)
	case Annotation_IDGENERATOR:
		return p.ParseIdGeneratorTag(p.Next())
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
		return p.parseTypeBearingTag(token, annotation)
	}

	return p.Next()
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
