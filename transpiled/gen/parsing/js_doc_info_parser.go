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
// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.jscomp.parsing.JsDocInfoParser
type JsDocInfoParser struct {
	stream                *JsDocTokenStream                                 `gson:"stream"`
	jsdocBuilder          *pkg_rhino.JSDocInfo_Builder                      `gson:"jsdocBuilder"`
	errorReporter         pkg_rhino.ErrorReporter                           `gson:"errorReporter"`
	fileOverviewJSDocInfo *pkg_rhino.JSDocInfo                              `gson:"fileOverviewJSDocInfo"`
	annotations           j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"annotations"`
	suppressionNames      j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"suppressionNames"`
	closurePrimitiveNames j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"closurePrimitiveNames"`
	licenseText           string                                            `gson:"licenseText"`
	reusableExtendedTypes j2gort.Ext_java_util_ArrayList                    `gson:"reusableExtendedTypes"`
	templateNode          pkg_rhino.Node                                    `gson:"templateNode"`
	idHash_               atomic.Int32
	jsDocSourceKind       JsDocInfoParser_JsDocSourceKind `gson:"jsDocSourceKind"`
	state                 JsDocInfoParser_State           `gson:"state"`
	preserveWhitespace    bool                            `gson:"preserveWhitespace"`
	onlyParseLicenseJsDoc bool                            `gson:"onlyParseLicenseJsDoc"`
	recordTsType          bool                            `gson:"recordTsType"`
	unreadToken           JsDocToken                      `gson:"unreadToken"`
	nodeAlloc_            *nodestore.NodeAllocator
}

func (this_ *JsDocInfoParser) BindSelf_(s any) {
	_ = s
}
func (this_ *JsDocInfoParser) IsJavaNil_() bool { return this_ == nil }
func (*JsDocInfoParser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocInfoParser", "JsDocInfoParser"
}
func (this_ *JsDocInfoParser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsDocInfoParser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocInfoParser", "com.google.javascript.jscomp.parsing.JsDocInfoParser", func(v_ any) bool {
		return j2gort.InstanceOf[*JsDocInfoParser](v_)
	})
}

func StaticGet_JsDocInfoParser_TTL_START_DELIMITER() string { return "=" }
func StaticGet_JsDocInfoParser_TTL_END_DELIMITER() string   { return "=:" }

var static_JsDocInfoParser_TEMPLATE_NAME_MATCHER j2gort.Ext_com_google_common_base_CharMatcher

func StaticGet_JsDocInfoParser_TEMPLATE_NAME_MATCHER() j2gort.Ext_com_google_common_base_CharMatcher {
	clinit_JsDocInfoParser()
	return static_JsDocInfoParser_TEMPLATE_NAME_MATCHER
}
func StaticSet_JsDocInfoParser_TEMPLATE_NAME_MATCHER(v j2gort.Ext_com_google_common_base_CharMatcher) {
	clinit_JsDocInfoParser()
	static_JsDocInfoParser_TEMPLATE_NAME_MATCHER = v
}
func StaticGet_JsDocInfoParser_BAD_TYPE_WIKI_LINK() string {
	return " See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler for more information."
}
func StaticGet_JsDocInfoParser_TSICKLE_MISSING_TYPE_PLACEHOLDER() string {
	return "JsDocInfoParser_TsickleMode_MissingSupertypePlaceholder"
}
func StaticGet_JsDocInfoParser_TS_TYPE_PLACEHOLDER() string {
	return "JsDocInfoParser_tsType_PlaceHolder"
}

var static_JsDocInfoParser_modifiesAnnotationKeywords j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_JsDocInfoParser_ModifiesAnnotationKeywords() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_JsDocInfoParser()
	return static_JsDocInfoParser_modifiesAnnotationKeywords
}
func StaticSet_JsDocInfoParser_ModifiesAnnotationKeywords(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_JsDocInfoParser()
	static_JsDocInfoParser_modifiesAnnotationKeywords = v
}

var static_JsDocInfoParser_idGeneratorAnnotationKeywords j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_JsDocInfoParser_IdGeneratorAnnotationKeywords() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_JsDocInfoParser()
	return static_JsDocInfoParser_idGeneratorAnnotationKeywords
}
func StaticSet_JsDocInfoParser_IdGeneratorAnnotationKeywords(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_JsDocInfoParser()
	static_JsDocInfoParser_idGeneratorAnnotationKeywords = v
}

var static_JsDocInfoParser_primitiveTypes j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_JsDocInfoParser_PrimitiveTypes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_JsDocInfoParser()
	return static_JsDocInfoParser_primitiveTypes
}
func StaticSet_JsDocInfoParser_PrimitiveTypes(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_JsDocInfoParser()
	static_JsDocInfoParser_primitiveTypes = v
}

var static_JsDocInfoParser_NO_UNREAD_TOKEN JsDocToken

func StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN() JsDocToken {
	clinit_JsDocInfoParser()
	return static_JsDocInfoParser_NO_UNREAD_TOKEN
}
func StaticSet_JsDocInfoParser_NO_UNREAD_TOKEN(v JsDocToken) {
	clinit_JsDocInfoParser()
	static_JsDocInfoParser_NO_UNREAD_TOKEN = v
}

var clinitDone_JsDocInfoParser uint32
var clinitState_JsDocInfoParser j2gort.ClinitState

func clinit_JsDocInfoParser() {
	if atomic.LoadUint32(&clinitDone_JsDocInfoParser) != 0 {
		return
	}
	clinitSlow_JsDocInfoParser()
}

func clinitSlow_JsDocInfoParser() {
	if !j2gort.EnterClinit(&clinitState_JsDocInfoParser, &clinitDone_JsDocInfoParser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsDocInfoParser, &clinitDone_JsDocInfoParser)
	static_JsDocInfoParser_TEMPLATE_NAME_MATCHER = j2gort.Ext_com_google_common_base_CharMatcher_S_JavaLetterOrDigit().Or_1_Ext_com_google_common_base_CharMatcher(j2gort.Ext_com_google_common_base_CharMatcher_S_AnyOf("_$"))
	static_JsDocInfoParser_modifiesAnnotationKeywords = j2gort.ImmSetOfStr2("this", "arguments")
	static_JsDocInfoParser_idGeneratorAnnotationKeywords = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_5("unique", "consistent", "stable", "mapped", "xid")
	static_JsDocInfoParser_primitiveTypes = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("number", "string", "boolean", "symbol")
	static_JsDocInfoParser_NO_UNREAD_TOKEN = JsDocToken(0)
}

func (this_ *JsDocInfoParser) initFields_JsDocInfoParser() {
	_ = this_
	this_.jsdocBuilder = pkg_rhino.JSDocInfo_S_Builder()

	this_.licenseText = j2gort.NullString
	this_.reusableExtendedTypes = j2gort.NewExt_java_util_ArrayList()
	this_.unreadToken = StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN()
}

// Java constructor: JsDocInfoParser(com.google.javascript.jscomp.parsing.JsDocTokenStream stream, java.lang.String comment, int commentPosition, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node templateNode, com.google.javascript.jscomp.parsing.Config config, com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind jsDocSourceKind, com.google.javascript.rhino.ErrorReporter errorReporter)
func NewJsDocInfoParser(stream *JsDocTokenStream, comment string, commentPosition int32, templateNode pkg_rhino.Node, config Config, jsDocSourceKind JsDocInfoParser_JsDocSourceKind, errorReporter pkg_rhino.ErrorReporter) *JsDocInfoParser {
	clinit_JsDocInfoParser()
	this_ := &JsDocInfoParser{}
	InitJsDocInfoParser(this_, stream, comment, commentPosition, templateNode, config, jsDocSourceKind, errorReporter)
	return this_
}

func InitJsDocInfoParser(this_ *JsDocInfoParser, stream *JsDocTokenStream, comment string, commentPosition int32, templateNode pkg_rhino.Node, config Config, jsDocSourceKind JsDocInfoParser_JsDocSourceKind, errorReporter pkg_rhino.ErrorReporter) {
	_ = stream
	_ = comment
	_ = commentPosition
	_ = templateNode
	_ = config
	_ = jsDocSourceKind
	_ = errorReporter
	this_.initFields_JsDocInfoParser()
	this_.Reset(stream, comment, commentPosition, templateNode, config, jsDocSourceKind, errorReporter)
}

// Java method: JsDocInfoParser.addParserWarning(com.google.javascript.rhino.Msg msg, java.lang.String messageArg)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addParserWarning_7008(msg pkg_rhino.Msg, messageArg string) {
	_ = msg
	_ = messageArg
	this_.priv_JsDocInfoParser_addParserWarning_7009(msg, messageArg, this_.stream.GetLineno(), this_.stream.GetCharno())
}

// Java method: JsDocInfoParser.addParserWarning(com.google.javascript.rhino.Msg msg, java.lang.String messageArg, int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addParserWarning_7009(msg pkg_rhino.Msg, messageArg string, lineno int32, charno int32) {
	_ = msg
	_ = messageArg
	_ = lineno
	_ = charno
	this_.errorReporter.Warning(msg.Format_1([]any{j2gort.BoxString(messageArg)}), this_.priv_JsDocInfoParser_getSourceName_7023(), lineno, charno)
}

// Java method: JsDocInfoParser.addParserWarning(com.google.javascript.rhino.Msg msg)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addParserWarning_7010(msg pkg_rhino.Msg) {
	_ = msg
	this_.priv_JsDocInfoParser_addParserWarning_7011(msg, this_.stream.GetLineno(), this_.stream.GetCharno())
}

// Java method: JsDocInfoParser.addParserWarning(com.google.javascript.rhino.Msg msg, int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addParserWarning_7011(msg pkg_rhino.Msg, lineno int32, charno int32) {
	_ = msg
	_ = lineno
	_ = charno
	this_.errorReporter.Warning(msg.Format_0(), this_.priv_JsDocInfoParser_getSourceName_7023(), lineno, charno)
}

// Java method: JsDocInfoParser.addTypeWarning(com.google.javascript.rhino.Msg msg, java.lang.String messageArg)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addTypeWarning_7012(msg pkg_rhino.Msg, messageArg string) {
	_ = msg
	_ = messageArg
	this_.priv_JsDocInfoParser_addTypeWarning_7013(msg, messageArg, this_.stream.GetLineno(), this_.stream.GetCharno())
}

// Java method: JsDocInfoParser.addTypeWarning(com.google.javascript.rhino.Msg msg, java.lang.String messageArg, int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addTypeWarning_7013(msg pkg_rhino.Msg, messageArg string, lineno int32, charno int32) {
	_ = msg
	_ = messageArg
	_ = lineno
	_ = charno
	this_.errorReporter.Warning(("Bad type annotation. " + j2gort.StrValString(msg.Format_1([]any{j2gort.BoxString(messageArg)})) + " See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler for more information."), this_.priv_JsDocInfoParser_getSourceName_7023(), lineno, charno)
}

// Java method: JsDocInfoParser.addTypeWarning(com.google.javascript.rhino.Msg msg)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addTypeWarning_7014(msg pkg_rhino.Msg) {
	_ = msg
	this_.priv_JsDocInfoParser_addTypeWarning_7015(msg, this_.stream.GetLineno(), this_.stream.GetCharno())
}

// Java method: JsDocInfoParser.addTypeWarning(com.google.javascript.rhino.Msg msg, int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addTypeWarning_7015(msg pkg_rhino.Msg, lineno int32, charno int32) {
	_ = msg
	_ = lineno
	_ = charno
	this_.errorReporter.Warning(("Bad type annotation. " + j2gort.StrValString(msg.Format_0()) + " See https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler for more information."), this_.priv_JsDocInfoParser_getSourceName_7023(), lineno, charno)
}

// Java method: JsDocInfoParser.addMissingTypeWarning(int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_addMissingTypeWarning_7016(lineno int32, charno int32) {
	_ = lineno
	_ = charno
	this_.errorReporter.Warning("Missing type declaration.", this_.priv_JsDocInfoParser_getSourceName_7023(), lineno, charno)
}

// Java method: JsDocInfoParser.setLicenseTextMonotonic(java.lang.String x)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_setLicenseTextMonotonic_7017(x string) {
	_ = x
	j2gort.CheckNotNullString(x)
	if !(this_.licenseText == j2gort.NullString) {
		j2gort.CheckState(false)
	}
	this_.licenseText = x
}

// Java method: JsDocInfoParser.getLicenseText()
func (this_ *JsDocInfoParser) GetLicenseText() string {
	return this_.licenseText
}

// Java method: JsDocInfoParser.setFileOverviewJSDocInfo(com.google.javascript.rhino.JSDocInfo fileOverviewJSDocInfo)
func (this_ *JsDocInfoParser) SetFileOverviewJSDocInfo(fileOverviewJSDocInfo *pkg_rhino.JSDocInfo) {
	_ = fileOverviewJSDocInfo
	this_.fileOverviewJSDocInfo = fileOverviewJSDocInfo
}

// Java method: JsDocInfoParser.getSourceFile()
func (this_ *JsDocInfoParser) GetSourceFile() pkg_rhino.StaticSourceFile {
	return this_.templateNode.GetStaticSourceFile()
}

// Java method: JsDocInfoParser.reset(com.google.javascript.jscomp.parsing.JsDocTokenStream stream, java.lang.String comment, int commentPosition, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node templateNode, com.google.javascript.jscomp.parsing.Config config, com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind jsDocSourceKind, com.google.javascript.rhino.ErrorReporter errorReporter)
func (this_ *JsDocInfoParser) Reset(stream *JsDocTokenStream, comment string, commentPosition int32, templateNode pkg_rhino.Node, config Config, jsDocSourceKind JsDocInfoParser_JsDocSourceKind, errorReporter pkg_rhino.ErrorReporter) {
	_ = stream
	_ = comment
	_ = commentPosition
	_ = templateNode
	_ = config
	_ = jsDocSourceKind
	_ = errorReporter
	this_.stream = stream
	this_.fileOverviewJSDocInfo = nil
	this_.state = JsDocInfoParser_State(0)
	this_.licenseText = j2gort.NullString
	this_.unreadToken = StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN()
	var parseDocumentation bool = config.JsDocParsingMode().ShouldParseDescriptions()
	_ = parseDocumentation
	this_.jsdocBuilder.Reset()
	if parseDocumentation {
		this_.jsdocBuilder.ParseDocumentation()
	}
	if comment != j2gort.NullString {
		this_.jsdocBuilder.RecordOriginalCommentString(comment)
		this_.jsdocBuilder.RecordOriginalCommentPosition(commentPosition)
	}
	this_.recordTsType = parseDocumentation
	this_.annotations = config.Annotations()
	this_.suppressionNames = config.SuppressionNames()
	this_.closurePrimitiveNames = config.ClosurePrimitiveNames()
	this_.preserveWhitespace = config.JsDocParsingMode().ShouldPreserveWhitespace()
	this_.jsDocSourceKind = jsDocSourceKind
	this_.onlyParseLicenseJsDoc = (config.JsDocParsingMode() == StaticGet_Config_JsDocParsing_LICENSE_COMMENTS_ONLY())
	this_.errorReporter = errorReporter
	if templateNode == 0 {
		this_.templateNode = pkg_rhino.IR_Script_0()
	} else {
		this_.templateNode = templateNode
	}
}

// Java method: JsDocInfoParser.clearForReuse()
func (this_ *JsDocInfoParser) ClearForReuse() {
	this_.stream = nil
	this_.jsdocBuilder.Reset()
	this_.errorReporter = nil
	this_.templateNode = pkg_rhino.Node(0)
	this_.fileOverviewJSDocInfo = nil
	this_.state = JsDocInfoParser_State(0)
	this_.annotations = nil
	this_.suppressionNames = nil
	this_.closurePrimitiveNames = nil
	this_.licenseText = j2gort.NullString
	this_.unreadToken = StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN()
	if !(this_.reusableExtendedTypes.IsEmpty()) {
		this_.reusableExtendedTypes.Clear()
	}
}

// Java method: JsDocInfoParser.getSourceName()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_getSourceName_7023() string {
	var sourceFile pkg_rhino.StaticSourceFile = this_.GetSourceFile()
	_ = sourceFile
	if sourceFile == nil || sourceFile.IsJavaNil_() {
		return j2gort.NullString
	} else {
		return sourceFile.GetName()
	}
}

// Java method: JsDocInfoParser.parseInlineTypeDoc()
func (this_ *JsDocInfoParser) ParseInlineTypeDoc() *pkg_rhino.JSDocInfo {
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	var token JsDocToken = this_.priv_JsDocInfoParser_next_7086()
	_ = token
	var lineno int32 = this_.stream.GetLineno()
	_ = lineno
	var startCharno int32 = this_.stream.GetCharno()
	_ = startCharno
	var typeAst pkg_rhino.Node = this_.priv_JsDocInfoParser_parseParamTypeExpression_7052(token)
	_ = typeAst
	this_.priv_JsDocInfoParser_recordTypeNode_7025(lineno, startCharno, typeAst, (token == StaticGet_JsDocToken_LEFT_CURLY()))
	var expr pkg_rhino.JSTypeExpression = this_.CreateJSTypeExpression(typeAst)
	_ = expr
	if expr != (pkg_rhino.JSTypeExpression{}) {
		this_.jsdocBuilder.RecordType(expr)
		this_.jsdocBuilder.RecordInlineType()
		return this_.RetrieveAndResetParsedJSDocInfo()
	}
	return nil
}

// Java method: JsDocInfoParser.recordTypeNode(int lineno, int startCharno, com.google.javascript.rhino.Node typeAst, boolean matchingLC)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_recordTypeNode_7025(lineno int32, startCharno int32, typeAst pkg_rhino.Node, matchingLC bool) {
	_ = lineno
	_ = startCharno
	_ = typeAst
	_ = matchingLC
	if typeAst != 0 {
		var endLineno int32 = this_.stream.GetLineno()
		_ = endLineno
		var endCharno int32 = this_.stream.GetCharno()
		_ = endCharno
		this_.jsdocBuilder.MarkTypeNode(typeAst, lineno, startCharno, endLineno, endCharno, matchingLC)
	}
}

// Java static method: JsDocInfoParser.parseTypeString(java.lang.String typeString)
func JsDocInfoParser_ParseTypeString(typeString string) pkg_rhino.Node {
	clinit_JsDocInfoParser()
	_ = typeString
	var parser *JsDocInfoParser = JsDocInfoParser_GetParser(typeString)
	_ = parser
	return parser.priv_JsDocInfoParser_parseTopLevelTypeExpression_7055(parser.priv_JsDocInfoParser_next_7086())
}

// Java static method: JsDocInfoParser.parseJsdoc(java.lang.String toParse)
func JsDocInfoParser_ParseJsdoc(toParse string) *pkg_rhino.JSDocInfo {
	clinit_JsDocInfoParser()
	_ = toParse
	var parser *JsDocInfoParser = JsDocInfoParser_GetParser(toParse)
	_ = parser
	parser.Parse()
	return parser.RetrieveAndResetParsedJSDocInfo()
}

// Java static method: JsDocInfoParser.parseFileOverviewJsdoc(java.lang.String toParse)
func JsDocInfoParser_ParseFileOverviewJsdoc(toParse string) *pkg_rhino.JSDocInfo {
	clinit_JsDocInfoParser()
	_ = toParse
	var parser *JsDocInfoParser = JsDocInfoParser_GetParser(toParse)
	_ = parser
	parser.Parse()
	return parser.GetFileOverviewJSDocInfo()
}

// Java static method: JsDocInfoParser.getParser(java.lang.String toParse)
func JsDocInfoParser_GetParser(toParse string) *JsDocInfoParser {
	clinit_JsDocInfoParser()
	_ = toParse
	var config Config = Config_S_Builder().SetLanguageMode(StaticGet_Config_LanguageMode_ECMASCRIPT3()).SetStrictMode(StaticGet_Config_StrictMode_SLOPPY()).SetClosurePrimitiveNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmSetOfStr1("testPrimitive"))).Build()
	_ = config
	return NewJsDocInfoParser(NewJsDocTokenStream_1(toParse), toParse, int32(0), pkg_rhino.Node(0), config, StaticGet_JsDocInfoParser_JsDocSourceKind_NORMAL(), pkg_rhino.StaticGet_ErrorReporter_NULL_INSTANCE())
}

// Java method: JsDocInfoParser.parse()
func (this_ *JsDocInfoParser) Parse() bool {
	this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	var token JsDocToken = this_.priv_JsDocInfoParser_next_7086()
	_ = token
	if this_.jsdocBuilder.ShouldParseDocumentation() {
		var blockInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractBlockComment_7048(token)
		_ = blockInfo
		token = blockInfo.token
		if !(j2gort.CheckNotNullString(blockInfo.string_) == "") {
			this_.jsdocBuilder.RecordBlockDescription(blockInfo.string_)
		}
	} else {
		if (token != StaticGet_JsDocToken_ANNOTATION()) && (token != StaticGet_JsDocToken_EOC()) {
			this_.jsdocBuilder.RecordBlockDescription("")
		}
	}
	var extendedTypes j2gort.Ext_java_util_ArrayList = this_.reusableExtendedTypes
	_ = extendedTypes
	if !(extendedTypes.IsEmpty()) {
		extendedTypes.Clear()
	}
	defer func() {
		if !(extendedTypes.IsEmpty()) {
			extendedTypes.Clear()
		}
	}()
	{
		return this_.priv_JsDocInfoParser_parseHelperLoop_7028(token, extendedTypes)
	}
	panic("unreachable")
}

// Java method: JsDocInfoParser.parseImportantComment()
func (this_ *JsDocInfoParser) ParseImportantComment() {
	this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	var token JsDocToken = this_.priv_JsDocInfoParser_next_7086()
	_ = token
	var info JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineComment_7049(token, StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE(), false, true)
	_ = info
	var license string = (" " + j2gort.StrValString(info.string_))
	_ = license
	this_.priv_JsDocInfoParser_setLicenseTextMonotonic_7017(license)
	if this_.jsdocBuilder.ShouldParseDocumentation() {
		this_.jsdocBuilder.RecordBlockDescription(license)
	} else {
		this_.jsdocBuilder.RecordBlockDescription("")
	}
}

// Java method: JsDocInfoParser.parseHelperLoop(com.google.javascript.jscomp.parsing.JsDocToken token, java.util.ArrayList<com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo> extendedTypes)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseHelperLoop_7028(token JsDocToken, extendedTypes j2gort.Ext_java_util_ArrayList) bool {
	_ = token
	_ = extendedTypes
	for {
		switch (token).Ordinal() {
		case int32(14) /* JsDocToken.ANNOTATION */ :
			if this_.state == StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION() {
				this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_NEWLINE()
				token = this_.priv_JsDocInfoParser_parseAnnotation_7029(token, extendedTypes)
			} else {
				token = this_.priv_JsDocInfoParser_next_7086()
			}
		case int32(17) /* JsDocToken.EOC */ :
			var success bool = true
			_ = success
			this_.priv_JsDocInfoParser_checkExtendedTypes_7032(extendedTypes)
			if this_.priv_JsDocInfoParser_hasParsedFileOverviewDocInfo_7090() {
				this_.fileOverviewJSDocInfo = this_.RetrieveAndResetParsedJSDocInfo()
				var visibility pkg_rhino.JSDocInfo_Visibility = this_.fileOverviewJSDocInfo.GetVisibility()
				_ = visibility
				switch (visibility).Ordinal() {
				case int32(0) /* Visibility.PRIVATE */, int32(2) /* Visibility.PROTECTED */ :
					this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION(), j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.CheckNotNullEnum(visibility).String()))
					success = false
				default:
				}
			}
			return success
		case int32(0) /* JsDocToken.EOF */ :
			this_.jsdocBuilder.Reset()
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_UNEXPECTED_EOF())
			this_.priv_JsDocInfoParser_checkExtendedTypes_7032(extendedTypes)
			return false
		case int32(1) /* JsDocToken.EOL */ :
			if this_.state == StaticGet_JsDocInfoParser_State_SEARCHING_NEWLINE() {
				this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
			}
			token = this_.priv_JsDocInfoParser_next_7086()
		default:
			if (token == StaticGet_JsDocToken_STAR()) && (this_.state == StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()) {
				token = this_.priv_JsDocInfoParser_next_7086()
			} else {
				this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_NEWLINE()
				token = this_.priv_JsDocInfoParser_eatTokensUntilEOL_7080()
			}
		}
	}
	panic("unreachable")
}

// Java method: JsDocInfoParser.parseAnnotation(com.google.javascript.jscomp.parsing.JsDocToken token, java.util.ArrayList<com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo> extendedTypes)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseAnnotation_7029(token JsDocToken, extendedTypes j2gort.Ext_java_util_ArrayList) JsDocToken {
	_ = token
	_ = extendedTypes
	var type_ pkg_rhino.JSTypeExpression = (pkg_rhino.JSTypeExpression{})
	_ = type_
	var lineno int32 = this_.stream.GetLineno()
	_ = lineno
	var charno int32 = this_.stream.GetCharno()
	_ = charno
	var annotationName string = this_.stream.GetString()
	_ = annotationName
	var annotation Annotation = j2gort.Cast[Annotation](j2gort.MapGetStr(this_.annotations, annotationName))
	_ = annotation
	if this_.onlyParseLicenseJsDoc && !((annotation == StaticGet_Annotation_LICENSE()) || (annotation == StaticGet_Annotation_PRESERVE())) {
		return this_.priv_JsDocInfoParser_next_7086()
	}
	if (annotation == 0) || (j2gort.CheckNotNullString(annotationName) == "") {
		this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_BAD_JSDOC_TAG(), annotationName)
	} else {
		this_.jsdocBuilder.MarkAnnotation(annotationName, lineno, charno)
		var matchingRc bool = false
		_ = matchingRc
		switch (annotation).Ordinal() {
		case int32(38) /* Annotation.NG_INJECT */ :
			if this_.jsdocBuilder.IsNgInjectRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NGINJECT_EXTRA())
			} else {
				this_.jsdocBuilder.RecordNgInject(true)
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(0) /* Annotation.ABSTRACT */ :
			if !(this_.jsdocBuilder.RecordAbstract()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(2) /* Annotation.AUTHOR */ :
			if this_.jsdocBuilder.ShouldParseDocumentation() {
				var authorInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractSingleLineBlock_7044()
				_ = authorInfo
				var author string = authorInfo.string_
				_ = author
				if j2gort.CheckNotNullString(author) == "" {
					this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_AUTHORMISSING())
				} else {
					this_.jsdocBuilder.RecordAuthor(author)
				}
				token = authorInfo.token
			} else {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
			}
			return token
		case int32(73) /* Annotation.UNRESTRICTED */ :
			if !(this_.jsdocBuilder.RecordUnrestricted()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(64) /* Annotation.STRUCT */ :
			if !(this_.jsdocBuilder.RecordStruct()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(12) /* Annotation.DICT */ :
			if !(this_.jsdocBuilder.RecordDict()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(4) /* Annotation.COLLAPSIBLE_OR_BREAK_MY_CODE */ :
			if !(this_.jsdocBuilder.RecordCollapsibleOrBreakMyCode()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_COLLAPSIBLEORBREAKMYCODE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(6) /* Annotation.CONSTRUCTOR */ :
			if !(this_.jsdocBuilder.RecordConstructor()) {
				if this_.jsdocBuilder.IsInterfaceRecorded() {
					this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INTERFACE_CONSTRUCTOR())
				} else {
					this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE())
				}
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(8) /* Annotation.RECORD */ :
			if !(this_.jsdocBuilder.RecordImplicitMatch()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_RECORD())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(10) /* Annotation.DEPRECATED */ :
			if !(this_.jsdocBuilder.RecordDeprecated()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_DEPRECATED())
			}
			var reasonInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
			_ = reasonInfo
			var reason string = reasonInfo.string_
			_ = reason
			if j2gort.JStrLen(reason) > int32(0) {
				this_.jsdocBuilder.RecordDeprecationReason(reason)
			}
			token = reasonInfo.token
			return token
		case int32(25) /* Annotation.INTERFACE */ :
			if !(this_.jsdocBuilder.RecordInterface()) {
				if this_.jsdocBuilder.IsConstructorRecorded() {
					this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INTERFACE_CONSTRUCTOR())
				} else {
					this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE())
				}
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(11) /* Annotation.DESC */ :
			if this_.jsdocBuilder.IsDescriptionRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_DESC_EXTRA())
				return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
			} else {
				var descriptionInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
				_ = descriptionInfo
				var description string = descriptionInfo.string_
				_ = description
				this_.jsdocBuilder.RecordDescription(description)
				token = descriptionInfo.token
				return token
			}
		case int32(76) /* Annotation.TS_TYPE */ :
			var tsTypeInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
			_ = tsTypeInfo
			var tsType string = tsTypeInfo.string_
			_ = tsType
			if this_.recordTsType {
				this_.jsdocBuilder.RecordTsType(tsType)
			} else {
				this_.jsdocBuilder.RecordTsType("JsDocInfoParser_tsType_PlaceHolder")
			}
			token = tsTypeInfo.token
			return token
		case int32(19) /* Annotation.FILE_OVERVIEW */ :
			var fileOverview string = ""
			_ = fileOverview
			if this_.jsdocBuilder.ShouldParseDocumentation() && !(this_.priv_JsDocInfoParser_lookAheadForAnnotation_7094()) {
				var fileOverviewInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7046(token, this_.priv_JsDocInfoParser_getWhitespaceOption_7047(StaticGet_JsDocInfoParser_WhitespaceOption_TRIM()), false)
				_ = fileOverviewInfo
				fileOverview = fileOverviewInfo.string_
				token = fileOverviewInfo.token
			} else {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
			}
			if !(this_.jsdocBuilder.RecordFileOverview(fileOverview)) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_FILEOVERVIEW_EXTRA())
			}
			return token
		case int32(30) /* Annotation.LICENSE */, int32(51) /* Annotation.PRESERVE */ :
			var preserveInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7046(token, StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE(), true)
			_ = preserveInfo
			var preserve string = preserveInfo.string_
			_ = preserve
			if j2gort.JStrLen(preserve) > int32(0) {
				this_.priv_JsDocInfoParser_setLicenseTextMonotonic_7017(preserve)
			}
			token = preserveInfo.token
			return token
		case int32(14) /* Annotation.ENHANCE */ :
			var enhanceInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractSingleLineBlock_7044()
			_ = enhanceInfo
			var enhance string = enhanceInfo.string_
			_ = enhance
			token = enhanceInfo.token
			this_.jsdocBuilder.RecordEnhance(enhance)
			return token
		case int32(37) /* Annotation.MODS */ :
			if this_.jsdocBuilder.IsModsRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODS_EXTRA())
			} else {
				token = this_.priv_JsDocInfoParser_parseModsTag_7034(this_.priv_JsDocInfoParser_next_7086())
			}
			if (token != StaticGet_JsDocToken_EOC()) && (token != StaticGet_JsDocToken_EOF()) {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
			}
			return token
		case int32(15) /* Annotation.ENUM */ :
			token = this_.priv_JsDocInfoParser_next_7086()
			lineno = this_.stream.GetLineno()
			charno = this_.stream.GetCharno()
			type_ = (pkg_rhino.JSTypeExpression{})
			if (token != StaticGet_JsDocToken_EOL()) && (token != StaticGet_JsDocToken_EOC()) {
				var typeNode pkg_rhino.Node = this_.ParseAndRecordTypeNode_1(token)
				_ = typeNode
				if (typeNode != 0) && typeNode.IsStringLit() {
					var typeName string = typeNode.GetString()
					_ = typeName
					if !(j2gort.SetContainsStr(StaticGet_JsDocInfoParser_PrimitiveTypes(), typeName)) {
						typeNode = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_BANG(), typeNode)
					}
				}
				type_ = this_.CreateJSTypeExpression(typeNode)
			} else {
				this_.priv_JsDocInfoParser_restoreLookAhead_7082(token)
			}
			if type_ == (pkg_rhino.JSTypeExpression{}) {
				type_ = this_.CreateJSTypeExpression(this_.priv_JsDocInfoParser_newStringNode_7075("number"))
			}
			if !(this_.jsdocBuilder.RecordEnumParameterType(type_)) {
				this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), lineno, charno)
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(17) /* Annotation.EXTERNS */ :
			if !(this_.jsdocBuilder.RecordExterns()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_EXTERNS())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(41) /* Annotation.NO_COVERAGE */ :
			if !(this_.jsdocBuilder.RecordNoCoverage()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOCOVERAGE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(72) /* Annotation.TYPE_SUMMARY */ :
			if !(this_.jsdocBuilder.RecordTypeSummary()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_TYPESUMMARY())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(16) /* Annotation.EXTENDS */, int32(22) /* Annotation.IMPLEMENTS */ :
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			token = this_.priv_JsDocInfoParser_next_7086()
			lineno = this_.stream.GetLineno()
			charno = this_.stream.GetCharno()
			matchingRc = false
			if token == StaticGet_JsDocToken_LEFT_CURLY() {
				token = this_.priv_JsDocInfoParser_next_7086()
				matchingRc = true
			}
			if token == StaticGet_JsDocToken_STRING() {
				var typeNode_2 pkg_rhino.Node = this_.priv_JsDocInfoParser_parseAndRecordTypeNameNode_7040(token, lineno, charno, matchingRc)
				_ = typeNode_2
				lineno = this_.stream.GetLineno()
				charno = this_.stream.GetCharno()
				typeNode_2 = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_BANG(), typeNode_2)
				type_ = this_.CreateJSTypeExpression(typeNode_2)
				if annotation == StaticGet_Annotation_EXTENDS() {
					extendedTypes.Add_1(j2gort.BoxNullableVal(NewJsDocInfoParser_ExtendedTypeInfo(type_, this_.stream.GetLineno(), this_.stream.GetCharno())))
				} else {
					if !(annotation == StaticGet_Annotation_IMPLEMENTS()) {
						j2gort.CheckState(false)
					}
					if !(this_.jsdocBuilder.RecordImplementedInterface(type_)) {
						this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_IMPLEMENTS_DUPLICATE(), lineno, charno)
					}
				}
				token = this_.priv_JsDocInfoParser_next_7086()
				if matchingRc {
					if token != StaticGet_JsDocToken_RIGHT_CURLY() {
						this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
					} else {
						token = this_.priv_JsDocInfoParser_next_7086()
					}
				} else if ((token != StaticGet_JsDocToken_EOL()) && (token != StaticGet_JsDocToken_EOF())) && (token != StaticGet_JsDocToken_EOC()) {
					this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_END_ANNOTATION_EXPECTED())
				}
			} else if (token == StaticGet_JsDocToken_BANG()) || (token == StaticGet_JsDocToken_QMARK()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_IMPLEMENTS_EXTRAQUALIFIER(), lineno, charno)
			} else {
				this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_NO_TYPE_NAME(), lineno, charno)
			}
			token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(token)
			return token
		case int32(29) /* Annotation.LENDS */ :
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			matchingRc = false
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_LEFT_CURLY()) {
				token = this_.priv_JsDocInfoParser_next_7086()
				matchingRc = true
			}
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
				var lendsExpression pkg_rhino.JSTypeExpression = this_.CreateJSTypeExpression(this_.priv_JsDocInfoParser_parseNameExpression_7060(this_.priv_JsDocInfoParser_next_7086()))
				_ = lendsExpression
				if !(this_.jsdocBuilder.RecordLends(lendsExpression)) {
					this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_LENDS_INCOMPATIBLE())
				}
			} else {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_LENDS_MISSING())
			}
			if matchingRc && !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(33) /* Annotation.MEANING */ :
			var meaningInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
			_ = meaningInfo
			var meaning string = meaningInfo.string_
			_ = meaning
			token = meaningInfo.token
			if !(this_.jsdocBuilder.RecordMeaning(meaning)) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MEANING_EXTRA())
			}
			return token
		case int32(1) /* Annotation.ALTERNATE_MESSAGE_ID */ :
			var alternateMessageIdInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractSingleLineBlock_7044()
			_ = alternateMessageIdInfo
			var alternateMessageId string = alternateMessageIdInfo.string_
			_ = alternateMessageId
			token = alternateMessageIdInfo.token
			if !(this_.jsdocBuilder.RecordAlternateMessageId(alternateMessageId)) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA())
			}
			return token
		case int32(3) /* Annotation.CLOSURE_PRIMITIVE */ :
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			return this_.priv_JsDocInfoParser_parseClosurePrimitiveTag_7035(this_.priv_JsDocInfoParser_next_7086())
		case int32(40) /* Annotation.NO_COMPILE */ :
			if !(this_.jsdocBuilder.RecordNoCompile()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOCOMPILE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(42) /* Annotation.NO_DTS */ :
			if !(this_.jsdocBuilder.RecordNoDts()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NODTS())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(28) /* Annotation.LEGACY_NAMED_TS_EXPORT */ :
			if !(this_.jsdocBuilder.RecordLegacyNamedTsExport()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_LEGACY_TS_EXPORT())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(39) /* Annotation.NO_COLLAPSE */ :
			if !(this_.jsdocBuilder.RecordNoCollapse()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOCOLLAPSE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(13) /* Annotation.ENCOURAGE_INLINING */ :
			if !(this_.jsdocBuilder.RecordEncourageInlining()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_ENCOURAGE_INLINING())
			}
			if this_.jsdocBuilder.IsRequireInlining() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_INLINING())
			}
			if this_.jsdocBuilder.IsNoInline() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_INLINING())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(58) /* Annotation.REQUIRE_INLINING */ :
			if !(this_.jsdocBuilder.RecordRequireInlining()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_REQUIRE_INLINING())
			}
			if this_.jsdocBuilder.IsEncourageInlining() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_INLINING())
			}
			if this_.jsdocBuilder.IsNoInline() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_INLINING())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(43) /* Annotation.NO_INLINE */ :
			if !(this_.jsdocBuilder.RecordNoInline()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOINLINE())
			}
			if this_.jsdocBuilder.IsEncourageInlining() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_INLINING())
			}
			if this_.jsdocBuilder.IsRequireInlining() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_INLINING())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(54) /* Annotation.PROVIDE_GOOG */ :
			if !(this_.jsdocBuilder.RecordProvideGoog()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_PROVIDE_GOOG())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(55) /* Annotation.PROVIDE_ALREADY_PROVIDED */ :
			if !(this_.jsdocBuilder.RecordProvideAlreadyProvided()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_PROVIDE_ALREADY_PROVIDED())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(57) /* Annotation.PURE_OR_BREAK_MY_CODE */ :
			if !(this_.jsdocBuilder.RecordPureOrBreakMyCode()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_PUREORBREAKMYCODE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(24) /* Annotation.INHERIT_DOC */, int32(46) /* Annotation.OVERRIDE */ :
			if !(this_.jsdocBuilder.RecordOverride()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_OVERRIDE())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(49) /* Annotation.POLYMER */ :
			if this_.jsdocBuilder.IsPolymerRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_POLYMER_EXTRA())
			} else {
				this_.jsdocBuilder.RecordPolymer()
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(50) /* Annotation.POLYMER_BEHAVIOR */ :
			if this_.jsdocBuilder.IsPolymerBehaviorRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_POLYMERBEHAVIOR_EXTRA())
			} else {
				this_.jsdocBuilder.RecordPolymerBehavior()
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(7) /* Annotation.CUSTOM_ELEMENT */ :
			if this_.jsdocBuilder.IsCustomElementRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CUSTOMELEMENT_EXTRA())
			} else {
				this_.jsdocBuilder.RecordCustomElement()
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(34) /* Annotation.MIXIN_CLASS */ :
			if this_.jsdocBuilder.IsMixinClassRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MIXINCLASS_EXTRA())
			} else {
				this_.jsdocBuilder.RecordMixinClass()
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(35) /* Annotation.MIXIN_FUNCTION */ :
			if this_.jsdocBuilder.IsMixinFunctionRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MIXINFUNCTION_EXTRA())
			} else {
				this_.jsdocBuilder.RecordMixinFunction()
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(60) /* Annotation.SASS_GENERATED_CSS_TS */ :
			if this_.jsdocBuilder.IsSassGeneratedCssTsRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_SASS_GENERATED_CSS_TS())
			} else {
				this_.jsdocBuilder.RecordSassGeneratedCssTs()
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(67) /* Annotation.CLOSURE_UNAWARE_CODE */ :
			return this_.priv_JsDocInfoParser_parseClosureUnawareCode_7096()
		case int32(69) /* Annotation.THROWS */ :
			lineno = this_.stream.GetLineno()
			charno = this_.stream.GetCharno()
			var throwsAnnotation string = ""
			_ = throwsAnnotation
			if !(this_.priv_JsDocInfoParser_lookAheadForAnnotation_7094()) {
				var throwsInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
				_ = throwsInfo
				throwsAnnotation = javalang.StringTrim(j2gort.CheckNotNullString(throwsInfo.string_))
				token = throwsInfo.token
			}
			if this_.jsdocBuilder.IsNoSideEffectsRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS(), lineno, charno)
			} else if !(this_.jsdocBuilder.RecordThrowsAnnotation(throwsAnnotation)) {
				this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_THROWS(), lineno, charno)
			}
			return token
		case int32(48) /* Annotation.PARAM */ :
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			token = this_.priv_JsDocInfoParser_next_7086()
			lineno = this_.stream.GetLineno()
			charno = this_.stream.GetCharno()
			type_ = (pkg_rhino.JSTypeExpression{})
			var hasParamType bool = false
			_ = hasParamType
			if token == StaticGet_JsDocToken_LEFT_CURLY() {
				type_ = this_.CreateJSTypeExpression(this_.priv_JsDocInfoParser_parseAndRecordParamTypeNode_7041(token))
				if type_ == (pkg_rhino.JSTypeExpression{}) {
					return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
				}
				this_.priv_JsDocInfoParser_skipEOLs_7088()
				token = this_.priv_JsDocInfoParser_next_7086()
				lineno = this_.stream.GetLineno()
				charno = this_.stream.GetCharno()
				hasParamType = true
			}
			var name string = j2gort.NullString
			_ = name
			var isBracketedParam bool = (StaticGet_JsDocToken_LEFT_SQUARE() == token)
			_ = isBracketedParam
			if isBracketedParam {
				token = this_.priv_JsDocInfoParser_next_7086()
			}
			if StaticGet_JsDocToken_STRING() != token {
				this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_MISSING_VARIABLE_NAME(), lineno, charno)
			} else {
				if !(hasParamType) {
					this_.priv_JsDocInfoParser_addMissingTypeWarning_7016(this_.stream.GetLineno(), this_.stream.GetCharno())
				}
				name = this_.stream.GetString()
				if isBracketedParam {
					token = this_.priv_JsDocInfoParser_next_7086()
					if StaticGet_JsDocToken_EQUALS() == token {
						token = this_.priv_JsDocInfoParser_next_7086()
						if StaticGet_JsDocToken_STRING() == token {
							token = this_.priv_JsDocInfoParser_next_7086()
						}
					}
					if StaticGet_JsDocToken_RIGHT_SQUARE() != token {
						this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RB())
					} else if type_ != (pkg_rhino.JSTypeExpression{}) {
						type_ = pkg_rhino.JSTypeExpression_MakeOptionalArg(type_)
					}
				}
				if !(pkg_rhino.TokenStream_IsJSIdentifier(name)) {
					this_.priv_JsDocInfoParser_addParserWarning_7009(pkg_rhino.StaticGet_Msg_INVALID_VARIABLE_NAME(), name, lineno, charno)
					name = j2gort.NullString
				} else if !(this_.jsdocBuilder.RecordParameter(name, type_)) {
					if this_.jsdocBuilder.HasParameter(name) {
						this_.priv_JsDocInfoParser_addTypeWarning_7013(pkg_rhino.StaticGet_Msg_DUP_VARIABLE_NAME(), name, lineno, charno)
					} else {
						this_.priv_JsDocInfoParser_addTypeWarning_7013(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), name, lineno, charno)
					}
				}
			}
			if name == j2gort.NullString {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(token)
				return token
			}
			this_.jsdocBuilder.MarkName(name, this_.templateNode, lineno, charno)
			if this_.jsdocBuilder.ShouldParseDocumentation() && (token != StaticGet_JsDocToken_ANNOTATION()) {
				var paramDescriptionInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
				_ = paramDescriptionInfo
				var paramDescription string = paramDescriptionInfo.string_
				_ = paramDescription
				if j2gort.JStrLen(paramDescription) > int32(0) {
					this_.jsdocBuilder.RecordParameterDescription(name, paramDescription)
				}
				token = paramDescriptionInfo.token
			} else if (token != StaticGet_JsDocToken_EOC()) && (token != StaticGet_JsDocToken_EOF()) {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
			}
			return token
		case int32(44) /* Annotation.NO_SIDE_EFFECTS */ :
			if this_.jsdocBuilder.IsModifiesRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES())
			} else if this_.jsdocBuilder.IsThrowsRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS())
			} else if !(this_.jsdocBuilder.RecordNoSideEffects()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOSIDEEFFECTS())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(36) /* Annotation.MODIFIES */ :
			token = this_.priv_JsDocInfoParser_parseModifiesTag_7036(this_.priv_JsDocInfoParser_next_7086())
			return token
		case int32(23) /* Annotation.IMPLICIT_CAST */ :
			if !(this_.jsdocBuilder.RecordImplicitCast()) {
				this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_IMPLICITCAST())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(61) /* Annotation.SEE */ :
			if this_.jsdocBuilder.ShouldParseDocumentation() {
				var referenceInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractSingleLineBlock_7044()
				_ = referenceInfo
				var reference string = referenceInfo.string_
				_ = reference
				if j2gort.CheckNotNullString(reference) == "" {
					this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_SEEMISSING())
				} else {
					this_.jsdocBuilder.RecordReference(reference)
				}
				token = referenceInfo.token
			} else {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
			}
			return token
		case int32(65) /* Annotation.SUPPRESS */ :
			token = this_.priv_JsDocInfoParser_parseSuppressTag_7033(this_.priv_JsDocInfoParser_next_7086())
			return token
		case int32(66) /* Annotation.TEMPLATE */ :
			var boundTypeExpression pkg_rhino.JSTypeExpression = (pkg_rhino.JSTypeExpression{})
			_ = boundTypeExpression
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_LEFT_CURLY()) {
				this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_BOUNDEDGENERICS_USED(), lineno, charno)
				var boundTypeNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTypeExpressionAnnotation_7051(this_.priv_JsDocInfoParser_next_7086())
				_ = boundTypeNode
				if boundTypeNode != 0 {
					boundTypeExpression = this_.CreateJSTypeExpression(boundTypeNode)
				}
			}
			var templateNames []string = []string{}
			_ = templateNames
			if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COLON())) {
				for doFirst_1 := true; doFirst_1 || (this_.priv_JsDocInfoParser_eatIfMatch_7085(StaticGet_JsDocToken_COMMA())); doFirst_1 = false {
					var nameNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseNameExpression_7060(this_.priv_JsDocInfoParser_next_7086())
					_ = nameNode
					var templateName string
					_ = templateName
					if nameNode == 0 {
						templateName = j2gort.NullString
					} else {
						templateName = nameNode.GetString()
					}
					if JsDocInfoParser_ValidTemplateTypeName(templateName) {
						templateNames = append(templateNames, templateName)
					}
				}
			}
			var ttlAst pkg_rhino.Node = pkg_rhino.Node(0)
			_ = ttlAst
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COLON()) {
				var ttlExtraction JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7046(this_.priv_JsDocInfoParser_current_7087(), StaticGet_JsDocInfoParser_WhitespaceOption_TRIM(), false)
				_ = ttlExtraction
				token = ttlExtraction.token
				ttlAst = this_.priv_JsDocInfoParser_parseTtlAst_7030(ttlExtraction, lineno, charno)
			} else {
				token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(this_.priv_JsDocInfoParser_current_7087())
			}
			switch int32(len(templateNames)) {
			case int32(0):
				this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_NAME_MISSING(), lineno, charno)
				return token
			case int32(1):
			default:
				if (boundTypeExpression != (pkg_rhino.JSTypeExpression{})) || (ttlAst != 0) {
					this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_MULTIPLEDECLARATION(), lineno, charno)
				}
			}
			if (boundTypeExpression != (pkg_rhino.JSTypeExpression{})) && (ttlAst != 0) {
				this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_BOUNDSWITHTTL(), lineno, charno)
				return token
			}
			if ttlAst != 0 {
				if !(this_.jsdocBuilder.RecordTypeTransformation(templateNames[int32(0)], ttlAst)) {
					this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_NAME_REDECLARATION(), lineno, charno)
				}
			} else if boundTypeExpression != (pkg_rhino.JSTypeExpression{}) {
				if !(this_.jsdocBuilder.RecordTemplateTypeName_2(templateNames[int32(0)], boundTypeExpression)) {
					this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_NAME_REDECLARATION(), lineno, charno)
				}
			} else {
				for _, templateName_2 := range templateNames {
					_ = templateName_2
					if !(this_.jsdocBuilder.RecordTemplateTypeName_1(templateName_2)) {
						this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_NAME_REDECLARATION(), lineno, charno)
					}
				}

			}
			return token
		case int32(21) /* Annotation.IDGENERATOR */ :
			token = this_.priv_JsDocInfoParser_parseIdGeneratorTag_7037(this_.priv_JsDocInfoParser_next_7086())
			return token
		case int32(31) /* Annotation.LOG_TYPE_IN_COMPILER */ :
			var unused bool = this_.jsdocBuilder.RecordLogTypeInCompiler()
			_ = unused
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(45) /* Annotation.NOT_IMPLEMENTED */, int32(26) /* Annotation.JSX */, int32(27) /* Annotation.JSX_FRAGMENT */, int32(62) /* Annotation.SOY_MODULE */, int32(63) /* Annotation.SOY_TEMPLATE */, int32(77) /* Annotation.WIZ_ANALYZER */, int32(32) /* Annotation.MAY_HAVE_EXTRA_EDGE */ :
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(74) /* Annotation.USED_VIA_DOT_CONSTRUCTOR */ :
			if !(this_.jsdocBuilder.RecordUsedViaDotConstructor()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_USEDVIADOTCONSTRUCTOR())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(75) /* Annotation.WIZACTION */ :
			if !(this_.jsdocBuilder.RecordWizaction()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_WIZACTION())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(78) /* Annotation.WIZCALLBACK */ :
			if !(this_.jsdocBuilder.RecordWizcallback()) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_WIZCALLBACK())
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		case int32(5) /* Annotation.CONSTANT */, int32(20) /* Annotation.FINAL */, int32(9) /* Annotation.DEFINE */, int32(18) /* Annotation.EXPORT */, int32(59) /* Annotation.RETURN */, int32(47) /* Annotation.PACKAGE */, int32(52) /* Annotation.PRIVATE */, int32(53) /* Annotation.PROTECTED */, int32(56) /* Annotation.PUBLIC */, int32(68) /* Annotation.THIS */, int32(70) /* Annotation.TYPE */, int32(71) /* Annotation.TYPEDEF */ :
			lineno = this_.stream.GetLineno()
			charno = this_.stream.GetCharno()
			var typeNode_3 pkg_rhino.Node = pkg_rhino.Node(0)
			_ = typeNode_3
			var hasType bool = this_.priv_JsDocInfoParser_lookAheadForType_7093()
			_ = hasType
			var isAlternateTypeAnnotation bool = (((((((annotation == StaticGet_Annotation_PACKAGE()) || (annotation == StaticGet_Annotation_PRIVATE())) || (annotation == StaticGet_Annotation_PROTECTED())) || (annotation == StaticGet_Annotation_PUBLIC())) || (annotation == StaticGet_Annotation_CONSTANT())) || (annotation == StaticGet_Annotation_FINAL())) || (annotation == StaticGet_Annotation_EXPORT()))
			_ = isAlternateTypeAnnotation
			var canSkipTypeAnnotation bool = (isAlternateTypeAnnotation || (annotation == StaticGet_Annotation_RETURN()))
			_ = canSkipTypeAnnotation
			type_ = (pkg_rhino.JSTypeExpression{})
			if (annotation == StaticGet_Annotation_RETURN()) && !(hasType) {
				this_.priv_JsDocInfoParser_addMissingTypeWarning_7016(this_.stream.GetLineno(), this_.stream.GetCharno())
			}
			if hasType || !(canSkipTypeAnnotation) {
				this_.priv_JsDocInfoParser_skipEOLs_7088()
				token = this_.priv_JsDocInfoParser_next_7086()
				typeNode_3 = this_.ParseAndRecordTypeNode_1(token)
				if annotation == StaticGet_Annotation_THIS() {
					typeNode_3 = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_BANG(), typeNode_3)
				}
				type_ = this_.CreateJSTypeExpression(typeNode_3)
			}
			var hasError bool = ((type_ == (pkg_rhino.JSTypeExpression{})) && !(canSkipTypeAnnotation))
			_ = hasError
			if !(hasError) {
				if (((type_ != (pkg_rhino.JSTypeExpression{})) && isAlternateTypeAnnotation) || (annotation == StaticGet_Annotation_TYPE())) && !(this_.jsdocBuilder.RecordType(type_)) {
					this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), lineno, charno)
				}
				var isAnnotationNext bool = this_.priv_JsDocInfoParser_lookAheadForAnnotation_7094()
				_ = isAnnotationNext
				switch (annotation).Ordinal() {
				case int32(5) /* Annotation.CONSTANT */ :
					if !(this_.jsdocBuilder.RecordConstancy()) {
						this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CONST())
					}
				case int32(20) /* Annotation.FINAL */ :
					if !(this_.jsdocBuilder.RecordFinality()) {
						this_.priv_JsDocInfoParser_addTypeWarning_7014(pkg_rhino.StaticGet_Msg_JSDOC_FINAL())
					}
				case int32(9) /* Annotation.DEFINE */ :
					if !(this_.jsdocBuilder.RecordDefineType(type_)) {
						this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_DEFINE(), lineno, charno)
					}
					if !(isAnnotationNext) {
						return this_.priv_JsDocInfoParser_recordDescription_7031(token)
					}
				case int32(18) /* Annotation.EXPORT */ :
					if !(this_.jsdocBuilder.RecordExport()) {
						this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_EXPORT(), lineno, charno)
					}
					if !(isAnnotationNext) {
						return this_.priv_JsDocInfoParser_recordDescription_7031(token)
					}
				case int32(52) /* Annotation.PRIVATE */ :
					if !(this_.jsdocBuilder.RecordVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE())) {
						this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_EXTRA_VISIBILITY(), lineno, charno)
					}
					if !(isAnnotationNext) {
						return this_.priv_JsDocInfoParser_recordDescription_7031(token)
					}
				case int32(47) /* Annotation.PACKAGE */ :
					if !(this_.jsdocBuilder.RecordVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_PACKAGE())) {
						this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_EXTRA_VISIBILITY(), lineno, charno)
					}
					if !(isAnnotationNext) {
						return this_.priv_JsDocInfoParser_recordDescription_7031(token)
					}
				case int32(53) /* Annotation.PROTECTED */ :
					if !(this_.jsdocBuilder.RecordVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_PROTECTED())) {
						this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_EXTRA_VISIBILITY(), lineno, charno)
					}
					if !(isAnnotationNext) {
						return this_.priv_JsDocInfoParser_recordDescription_7031(token)
					}
				case int32(56) /* Annotation.PUBLIC */ :
					if !(this_.jsdocBuilder.RecordVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_PUBLIC())) {
						this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_EXTRA_VISIBILITY(), lineno, charno)
					}
					if !(isAnnotationNext) {
						return this_.priv_JsDocInfoParser_recordDescription_7031(token)
					}
				case int32(59) /* Annotation.RETURN */ :
					if type_ == (pkg_rhino.JSTypeExpression{}) {
						type_ = this_.CreateJSTypeExpression(this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_QMARK()))
					}
					if !(this_.jsdocBuilder.RecordReturnType(type_)) {
						this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), lineno, charno)
						break
					}
					if this_.jsdocBuilder.ShouldParseDocumentation() && !(isAnnotationNext) {
						var returnDescriptionInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
						_ = returnDescriptionInfo
						var returnDescription string = returnDescriptionInfo.string_
						_ = returnDescription
						if j2gort.JStrLen(returnDescription) > int32(0) {
							this_.jsdocBuilder.RecordReturnDescription(returnDescription)
						}
						token = returnDescriptionInfo.token
					} else {
						token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
					}
					return token
				case int32(68) /* Annotation.THIS */ :
					if !(this_.jsdocBuilder.RecordThisType(type_)) {
						this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), lineno, charno)
					}
				case int32(71) /* Annotation.TYPEDEF */ :
					if !(this_.jsdocBuilder.RecordTypedef(type_)) {
						this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), lineno, charno)
					}
				default:
				}
			}
			return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
		}
	}
	return this_.priv_JsDocInfoParser_next_7086()
}

// Java method: JsDocInfoParser.parseTtlAst(com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtractionInfo ttlExtraction, int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTtlAst_7030(ttlExtraction JsDocInfoParser_ExtractionInfo, lineno int32, charno int32) pkg_rhino.Node {
	_ = ttlExtraction
	_ = lineno
	_ = charno
	var ttlExpression string = ttlExtraction.string_
	_ = ttlExpression
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(ttlExpression), "=")) {
		return pkg_rhino.Node(0)
	}
	ttlExpression = j2gort.JStrSubFrom(ttlExpression, j2gort.JStrLen("="))
	var endIndex int32 = j2gort.JStrIndexOf(ttlExpression, "=:")
	_ = endIndex
	if endIndex >= int32(0) {
		ttlExpression = j2gort.JStrSub(ttlExpression, int32(0), endIndex)
	} else {
		this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER(), lineno, charno)
	}
	ttlExpression = javalang.StringTrim(j2gort.CheckNotNullString(ttlExpression))
	if j2gort.CheckNotNullString(ttlExpression) == "" {
		this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING(), lineno, charno)
		return pkg_rhino.Node(0)
	}
	clinit_TypeTransformationParser()
	var val_3 TypeTransformationParser
	InitTypeTransformationParser(&val_3, ttlExpression, this_.GetSourceFile(), this_.errorReporter, lineno, charno)
	var ttlParser *TypeTransformationParser = &val_3
	_ = ttlParser
	if !(ttlParser.ParseTypeTransformation()) {
		return pkg_rhino.Node(0)
	}
	return ttlParser.GetTypeTransformationAst()
}

// Java static method: JsDocInfoParser.validTemplateTypeName(java.lang.String name)
func JsDocInfoParser_ValidTemplateTypeName(name string) bool {
	clinit_JsDocInfoParser()
	_ = name
	return (((name != j2gort.NullString) && !(j2gort.CheckNotNullString(name) == "")) && StaticGet_JsDocInfoParser_TEMPLATE_NAME_MATCHER().MatchesAllOf(j2gort.BoxString(name)))
}

// Java method: JsDocInfoParser.recordDescription(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_recordDescription_7031(token JsDocToken) JsDocToken {
	_ = token
	if this_.jsdocBuilder.ShouldParseDocumentation() {
		var descriptionInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
		_ = descriptionInfo
		token = descriptionInfo.token
	} else {
		token = this_.priv_JsDocInfoParser_eatTokensUntilEOL_7081(token)
	}
	return token
}

// Java method: JsDocInfoParser.checkExtendedTypes(java.util.ArrayList<com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo> extendedTypes)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_checkExtendedTypes_7032(extendedTypes j2gort.Ext_java_util_ArrayList) {
	_ = extendedTypes
	if extendedTypes.Size() > int32(1) {
		var isInterface bool = this_.jsdocBuilder.IsInterfaceRecorded()
		_ = isInterface
		var isFromTsickle bool = func() bool {
			re_4 := this_.jsDocSourceKind
			ae_5 := StaticGet_JsDocInfoParser_JsDocSourceKind_TSICKLE()
			if re_4 == 0 {
				j2gort.ThrowNPE()
			}
			return re_4 == ae_5
		}()
		_ = isFromTsickle
		if !(isInterface) && isFromTsickle {
			var result bool = this_.jsdocBuilder.RecordBaseType(this_.CreateJSTypeExpression(this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_BANG(), this_.priv_JsDocInfoParser_newStringNode_7075("JsDocInfoParser_TsickleMode_MissingSupertypePlaceholder")).SrcrefTree(j2gort.ListGetVal[JsDocInfoParser_ExtendedTypeInfo](extendedTypes, int32(0)).type_.GetRoot())))
			_ = result
			if !(result) {
				j2gort.CheckStateStr(false, "Unexpected failure to record base type")
			}
			return
		}
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = extendedTypes.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var typeInfo JsDocInfoParser_ExtendedTypeInfo = j2gort.ListGetVal[JsDocInfoParser_ExtendedTypeInfo](extendedTypes, i)
			_ = typeInfo
			if this_.jsdocBuilder.IsInterfaceRecorded() {
				if !(this_.jsdocBuilder.RecordExtendedInterface(typeInfo.type_)) {
					this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_EXTENDS_DUPLICATE(), typeInfo.lineno, typeInfo.charno)
				}
			} else {
				if !(this_.jsdocBuilder.RecordBaseType(typeInfo.type_)) {
					this_.priv_JsDocInfoParser_addTypeWarning_7015(pkg_rhino.StaticGet_Msg_JSDOC_INCOMPAT_TYPE(), typeInfo.lineno, typeInfo.charno)
				}
			}
		}
	}
}

// Java method: JsDocInfoParser.parseSuppressTag(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseSuppressTag_7033(token JsDocToken) JsDocToken {
	_ = token
	if token != StaticGet_JsDocToken_LEFT_CURLY() {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_SUPPRESS())
		return token
	} else {
		var suppressions j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
		_ = suppressions
		for {
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
				var name string = this_.stream.GetString()
				_ = name
				if !(j2gort.SetContainsStr(this_.suppressionNames, name)) {
					this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_JSDOC_SUPPRESS_UNKNOWN(), name)
				}
				j2gort.SetAddStr(suppressions, this_.stream.GetString())
				token = this_.priv_JsDocInfoParser_next_7086()
			} else {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_SUPPRESS())
				return token
			}
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			if this_.priv_JsDocInfoParser_match_7084(StaticGet_JsDocToken_PIPE(), StaticGet_JsDocToken_COMMA()) {
				token = this_.priv_JsDocInfoParser_next_7086()
			} else {
				break
			}
		}
		token = this_.priv_JsDocInfoParser_next_7086()
		if token != StaticGet_JsDocToken_RIGHT_CURLY() {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_SUPPRESS())
			return token
		}
		if this_.jsdocBuilder.ShouldParseDocumentation() && !(this_.priv_JsDocInfoParser_lookAheadForAnnotation_7094()) {
			var suppressDescriptionInfo JsDocInfoParser_ExtractionInfo = this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token)
			_ = suppressDescriptionInfo
			var suppressDescription string = suppressDescriptionInfo.string_
			_ = suppressDescription
			this_.jsdocBuilder.RecordSuppressions_2(j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(suppressions), suppressDescription)
			token = suppressDescriptionInfo.token
		} else {
			token = this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(this_.priv_JsDocInfoParser_next_7086())
			this_.jsdocBuilder.RecordSuppressions_1(suppressions)
		}
		return token
	}
}

// Java method: JsDocInfoParser.parseModsTag(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseModsTag_7034(token JsDocToken) JsDocToken {
	_ = token
	if token != StaticGet_JsDocToken_LEFT_CURLY() {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODS())
		return token
	}
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING())) {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODS())
		return token
	}
	var namespace string = this_.stream.GetString()
	_ = namespace
	token = this_.priv_JsDocInfoParser_next_7086()
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODS())
	} else {
		token = this_.priv_JsDocInfoParser_next_7086()
		this_.jsdocBuilder.RecordMods(namespace)
	}
	return token
}

// Java method: JsDocInfoParser.parseClosurePrimitiveTag(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseClosurePrimitiveTag_7035(token JsDocToken) JsDocToken {
	_ = token
	if token != StaticGet_JsDocToken_LEFT_CURLY() {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_LC())
		return token
	} else if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
		var name string = this_.stream.GetString()
		_ = name
		if !(j2gort.SetContainsStr(this_.closurePrimitiveNames, name)) {
			this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_JSDOC_CLOSUREPRIMITIVE_INVALID(), name)
		} else if !(this_.jsdocBuilder.RecordClosurePrimitiveId(name)) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CLOSUREPRIMITIVE_EXTRA())
		}
		token = this_.priv_JsDocInfoParser_next_7086()
	} else {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CLOSUREPRIMITIVE_MISSING())
		return token
	}
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
		this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
	} else {
		token = this_.priv_JsDocInfoParser_next_7086()
	}
	return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078()
}

// Java method: JsDocInfoParser.parseModifiesTag(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseModifiesTag_7036(token JsDocToken) JsDocToken {
	_ = token
	if token == StaticGet_JsDocToken_LEFT_CURLY() {
		var modifies j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
		_ = modifies
		for {
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
				var name string = this_.stream.GetString()
				_ = name
				if !(j2gort.SetContainsStr(StaticGet_JsDocInfoParser_ModifiesAnnotationKeywords(), name)) && !(this_.jsdocBuilder.HasParameter(name)) {
					this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_JSDOC_MODIFIES_UNKNOWN(), name)
				}
				j2gort.SetAddStr(modifies, this_.stream.GetString())
				token = this_.priv_JsDocInfoParser_next_7086()
			} else {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODIFIES())
				return token
			}
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_PIPE()) {
				token = this_.priv_JsDocInfoParser_next_7086()
			} else {
				break
			}
		}
		if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODIFIES())
		} else {
			token = this_.priv_JsDocInfoParser_next_7086()
			if this_.jsdocBuilder.IsNoSideEffectsRecorded() {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES())
			} else if !(this_.jsdocBuilder.RecordModifies(modifies)) {
				this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_MODIFIES_DUPLICATE())
			}
		}
	}
	return token
}

// Java method: JsDocInfoParser.parseIdGeneratorTag(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseIdGeneratorTag_7037(token JsDocToken) JsDocToken {
	_ = token
	var idgenKind string = "unique"
	_ = idgenKind
	if token == StaticGet_JsDocToken_LEFT_CURLY() {
		if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
			var name string = this_.stream.GetString()
			_ = name
			if !(j2gort.SetContainsStr(StaticGet_JsDocInfoParser_IdGeneratorAnnotationKeywords(), name)) && !(this_.jsdocBuilder.HasParameter(name)) {
				this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_UNKNOWN(), name)
			}
			idgenKind = name
			token = this_.priv_JsDocInfoParser_next_7086()
		} else {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_BAD())
			return token
		}
		if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_BAD())
		} else {
			token = this_.priv_JsDocInfoParser_next_7086()
		}
	}
	switch j2gort.CheckNotNullString(idgenKind) {
	case "unique":
		if !(this_.jsdocBuilder.RecordIdGenerator()) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_DUPLICATE())
		}
	case "consistent":
		if !(this_.jsdocBuilder.RecordConsistentIdGenerator()) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_DUPLICATE())
		}
	case "stable":
		if !(this_.jsdocBuilder.RecordStableIdGenerator()) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_DUPLICATE())
		}
	case "xid":
		if !(this_.jsdocBuilder.RecordXidGenerator()) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_DUPLICATE())
		}
	case "mapped":
		if !(this_.jsdocBuilder.RecordMappedIdGenerator()) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_IDGEN_DUPLICATE())
		}
	}
	return token
}

// Java method: JsDocInfoParser.parseAndRecordTypeNode(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) ParseAndRecordTypeNode_1(token JsDocToken) pkg_rhino.Node {
	_ = token
	return this_.priv_JsDocInfoParser_parseAndRecordTypeNode_7039(token, this_.stream.GetLineno(), this_.stream.GetCharno(), (token == StaticGet_JsDocToken_LEFT_CURLY()), false)
}

// Java method: JsDocInfoParser.parseAndRecordTypeNode(com.google.javascript.jscomp.parsing.JsDocToken token, int lineno, int startCharno, boolean matchingLC, boolean onlyParseSimpleNames)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseAndRecordTypeNode_7039(token JsDocToken, lineno int32, startCharno int32, matchingLC bool, onlyParseSimpleNames bool) pkg_rhino.Node {
	_ = token
	_ = lineno
	_ = startCharno
	_ = matchingLC
	_ = onlyParseSimpleNames
	var typeNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = typeNode
	if onlyParseSimpleNames {
		typeNode = this_.priv_JsDocInfoParser_parseTypeNameAnnotation_7054(token)
	} else {
		typeNode = this_.priv_JsDocInfoParser_parseTypeExpressionAnnotation_7051(token)
	}
	this_.priv_JsDocInfoParser_recordTypeNode_7025(lineno, startCharno, typeNode, matchingLC)
	return typeNode
}

// Java method: JsDocInfoParser.parseAndRecordTypeNameNode(com.google.javascript.jscomp.parsing.JsDocToken token, int lineno, int startCharno, boolean matchingLC)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseAndRecordTypeNameNode_7040(token JsDocToken, lineno int32, startCharno int32, matchingLC bool) pkg_rhino.Node {
	_ = token
	_ = lineno
	_ = startCharno
	_ = matchingLC
	return this_.priv_JsDocInfoParser_parseAndRecordTypeNode_7039(token, lineno, startCharno, matchingLC, true)
}

// Java method: JsDocInfoParser.parseAndRecordParamTypeNode(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseAndRecordParamTypeNode_7041(token JsDocToken) pkg_rhino.Node {
	_ = token
	if !(token == StaticGet_JsDocToken_LEFT_CURLY()) {
		j2gort.CheckArgument(false)
	}
	var lineno int32 = this_.stream.GetLineno()
	_ = lineno
	var startCharno int32 = this_.stream.GetCharno()
	_ = startCharno
	var typeNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseParamTypeExpressionAnnotation_7053(token)
	_ = typeNode
	this_.priv_JsDocInfoParser_recordTypeNode_7025(lineno, startCharno, typeNode, true)
	return typeNode
}

// Java method: JsDocInfoParser.toString(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_toString_7042(token JsDocToken) string {
	_ = token
	switch (token).Ordinal() {
	case int32(14) /* JsDocToken.ANNOTATION */ :
		return ("@" + j2gort.StrValString(this_.stream.GetString()))
	case int32(19) /* JsDocToken.BANG */ :
		return "!"
	case int32(11) /* JsDocToken.COMMA */ :
		return ","
	case int32(12) /* JsDocToken.COLON */ :
		return ":"
	case int32(3) /* JsDocToken.RIGHT_ANGLE */ :
		return ">"
	case int32(5) /* JsDocToken.LEFT_SQUARE */ :
		return "["
	case int32(7) /* JsDocToken.LEFT_CURLY */ :
		return "{"
	case int32(9) /* JsDocToken.LEFT_PAREN */ :
		return "("
	case int32(2) /* JsDocToken.LEFT_ANGLE */ :
		return "<"
	case int32(18) /* JsDocToken.QMARK */ :
		return "?"
	case int32(15) /* JsDocToken.PIPE */ :
		return "|"
	case int32(6) /* JsDocToken.RIGHT_SQUARE */ :
		return "]"
	case int32(8) /* JsDocToken.RIGHT_CURLY */ :
		return "}"
	case int32(10) /* JsDocToken.RIGHT_PAREN */ :
		return ")"
	case int32(16) /* JsDocToken.STAR */ :
		return "*"
	case int32(13) /* JsDocToken.ITER_REST */ :
		return "..."
	case int32(20) /* JsDocToken.EQUALS */ :
		return "="
	case int32(4) /* JsDocToken.STRING */ :
		return this_.stream.GetString()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.CheckNotNullEnum(token).String()))))
	}
	panic("unreachable switch expression")
}

// Java method: JsDocInfoParser.createJSTypeExpression(com.google.javascript.rhino.Node n)
func (this_ *JsDocInfoParser) CreateJSTypeExpression(n pkg_rhino.Node) pkg_rhino.JSTypeExpression {
	_ = n
	if n == 0 {
		return (pkg_rhino.JSTypeExpression{})
	} else {
		return pkg_rhino.NewJSTypeExpression(n, this_.priv_JsDocInfoParser_getSourceName_7023())
	}
}

// Java method: JsDocInfoParser.extractSingleLineBlock()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_extractSingleLineBlock_7044() JsDocInfoParser_ExtractionInfo {
	this_.stream.Update()
	var lineno int32 = this_.stream.GetLineno()
	_ = lineno
	var charno int32 = (this_.stream.GetCharno() + int32(1))
	_ = charno
	var line string = javalang.StringTrim(j2gort.CheckNotNullString(this_.priv_JsDocInfoParser_getRemainingJSDocLine_7089()))
	_ = line
	if j2gort.JStrLen(line) > int32(0) {
		this_.jsdocBuilder.MarkText(line, lineno, charno, lineno, (charno + j2gort.JStrLen(line)))
	}
	return NewJsDocInfoParser_ExtractionInfo(line, this_.priv_JsDocInfoParser_next_7086())
}

// Java method: JsDocInfoParser.extractMultilineTextualBlock(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_extractMultilineTextualBlock_7045(token JsDocToken) JsDocInfoParser_ExtractionInfo {
	_ = token
	return this_.priv_JsDocInfoParser_extractMultilineTextualBlock_7046(token, this_.priv_JsDocInfoParser_getWhitespaceOption_7047(StaticGet_JsDocInfoParser_WhitespaceOption_SINGLE_LINE()), false)
}

// Java method: JsDocInfoParser.extractMultilineTextualBlock(com.google.javascript.jscomp.parsing.JsDocToken token, com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption option, boolean includeAnnotations)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_extractMultilineTextualBlock_7046(token JsDocToken, option JsDocInfoParser_WhitespaceOption, includeAnnotations bool) JsDocInfoParser_ExtractionInfo {
	_ = token
	_ = option
	_ = includeAnnotations
	if ((token == StaticGet_JsDocToken_EOC()) || (token == StaticGet_JsDocToken_EOL())) || (token == StaticGet_JsDocToken_EOF()) {
		return NewJsDocInfoParser_ExtractionInfo("", token)
	}
	return this_.priv_JsDocInfoParser_extractMultilineComment_7049(token, option, true, includeAnnotations)
}

// Java method: JsDocInfoParser.getWhitespaceOption(com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption defaultValue)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_getWhitespaceOption_7047(defaultValue JsDocInfoParser_WhitespaceOption) JsDocInfoParser_WhitespaceOption {
	_ = defaultValue
	if this_.preserveWhitespace {
		return StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE()
	} else {
		return defaultValue
	}
}

// Java method: JsDocInfoParser.extractBlockComment(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_extractBlockComment_7048(token JsDocToken) JsDocInfoParser_ExtractionInfo {
	_ = token
	return this_.priv_JsDocInfoParser_extractMultilineComment_7049(token, this_.priv_JsDocInfoParser_getWhitespaceOption_7047(StaticGet_JsDocInfoParser_WhitespaceOption_TRIM()), false, false)
}

// Java method: JsDocInfoParser.extractMultilineComment(com.google.javascript.jscomp.parsing.JsDocToken token, com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption option, boolean isMarker, boolean includeAnnotations)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_extractMultilineComment_7049(token JsDocToken, option JsDocInfoParser_WhitespaceOption, isMarker bool, includeAnnotations bool) JsDocInfoParser_ExtractionInfo {
	_ = token
	_ = option
	_ = isMarker
	_ = includeAnnotations
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = builder
	var startLineno int32 = int32(-1)
	_ = startLineno
	var startCharno int32 = int32(-1)
	_ = startCharno
	if isMarker {
		this_.stream.Update()
		startLineno = this_.stream.GetLineno()
		startCharno = (this_.stream.GetCharno() + int32(1))
		var line string = this_.priv_JsDocInfoParser_getRemainingJSDocLine_7089()
		_ = line
		if option != StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE() {
			line = javalang.StringTrim(j2gort.CheckNotNullString(line))
		}
		builder.Append_1_string(line)
		this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
		token = this_.priv_JsDocInfoParser_next_7086()
	}
	var ignoreStar bool = false
	_ = ignoreStar
	var lineStartChar int32 = int32(-1)
	_ = lineStartChar
	for doFirst_6 := true; doFirst_6 || (true); doFirst_6 = false {
		switch (token).Ordinal() {
		case int32(16) /* JsDocToken.STAR */ :
			if ignoreStar {
				lineStartChar = (this_.stream.GetCharno() + int32(1))
				ignoreStar = false
			} else {
				this_.priv_JsDocInfoParser_padLine_7050(builder, lineStartChar, option)
				lineStartChar = int32(-1)
				builder.Append_1_uint16(uint16(42))
			}
			token = this_.priv_JsDocInfoParser_next_7086()
			for token == StaticGet_JsDocToken_STAR() {
				if lineStartChar != int32(-1) {
					this_.priv_JsDocInfoParser_padLine_7050(builder, lineStartChar, option)
					lineStartChar = int32(-1)
				}
				builder.Append_1_uint16(uint16(42))
				token = this_.priv_JsDocInfoParser_next_7086()
			}
			continue
		case int32(1) /* JsDocToken.EOL */ :
			if option != StaticGet_JsDocInfoParser_WhitespaceOption_SINGLE_LINE() {
				builder.Append_1_uint16(uint16(10))
			}
			ignoreStar = true
			lineStartChar = int32(0)
			token = this_.priv_JsDocInfoParser_next_7086()
			continue
		default:
			ignoreStar = false
			this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
			var isEOC bool = (token == StaticGet_JsDocToken_EOC())
			_ = isEOC
			if !(isEOC) {
				this_.priv_JsDocInfoParser_padLine_7050(builder, lineStartChar, option)
				lineStartChar = int32(-1)
			}
			if ((token == StaticGet_JsDocToken_EOC()) || (token == StaticGet_JsDocToken_EOF())) || ((token == StaticGet_JsDocToken_ANNOTATION()) && !(includeAnnotations)) {
				var multilineText string = builder.String()
				_ = multilineText
				if option != StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE() {
					multilineText = javalang.StringTrim(j2gort.CheckNotNullString(multilineText))
				}
				if isMarker && !(j2gort.CheckNotNullString(multilineText) == "") {
					var endLineno int32 = this_.stream.GetLineno()
					_ = endLineno
					var endCharno int32 = this_.stream.GetCharno()
					_ = endCharno
					this_.jsdocBuilder.MarkText(multilineText, startLineno, startCharno, endLineno, endCharno)
				}
				return NewJsDocInfoParser_ExtractionInfo(multilineText, token)
			}
			builder.Append_1_string(this_.priv_JsDocInfoParser_toString_7042(token))
			var line_2 string = this_.priv_JsDocInfoParser_getRemainingJSDocLine_7089()
			_ = line_2
			if option != StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE() {
				line_2 = JsDocInfoParser_TrimEnd(line_2)
			}
			builder.Append_1_string(line_2)
			token = this_.priv_JsDocInfoParser_next_7086()
		}
	}
	panic("unreachable")
}

// Java method: JsDocInfoParser.padLine(java.lang.StringBuilder builder, int lineStartChar, com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption option)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_padLine_7050(builder *j2gort.Ext_java_lang_StringBuilder, lineStartChar int32, option JsDocInfoParser_WhitespaceOption) {
	_ = builder
	_ = lineStartChar
	_ = option
	if (lineStartChar != int32(-1)) && (option == StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE()) {
		var numSpaces int32 = (this_.stream.GetCharno() - lineStartChar)
		_ = numSpaces
		{
			var i int32 = int32(0)
			_ = i
			for ; i < numSpaces; i = i + 1 {
				builder.Append_1_uint16(uint16(32))
			}
		}
	} else if builder.Length() > int32(0) {
		if (int32(builder.CharAt((builder.Length() - int32(1)))) != int32(uint16(10))) || (option == StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE()) {
			builder.Append_1_uint16(uint16(32))
		}
	}
}

// Java static method: JsDocInfoParser.trimEnd(java.lang.String s)
func JsDocInfoParser_TrimEnd(s string) string {
	clinit_JsDocInfoParser()
	_ = s
	var trimCount int32 = int32(0)
	_ = trimCount
	for trimCount < j2gort.JStrLen(s) {
		var ch uint16 = j2gort.JStrCharAt(s, ((j2gort.JStrLen(s) - trimCount) - int32(1)))
		_ = ch
		if pkg_rhino.TokenUtil_IsWhitespace(int32(ch)) {
			trimCount = trimCount + 1
		} else {
			break
		}
	}
	if trimCount == int32(0) {
		return s
	}
	return j2gort.JStrSub(s, int32(0), (j2gort.JStrLen(s) - trimCount))
}

// Java method: JsDocInfoParser.parseTypeExpressionAnnotation(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTypeExpressionAnnotation_7051(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token == StaticGet_JsDocToken_LEFT_CURLY() {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		var typeNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTopLevelTypeExpression_7055(this_.priv_JsDocInfoParser_next_7086())
		_ = typeNode
		if typeNode != 0 {
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
				if typeNode.IsStringLit() && func() bool {
					sr_7 := "import"
					se_8 := typeNode.GetString()
					return j2gort.CheckNotNullString(sr_7) == se_8
				}() {
					this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_IMPORT())
				} else {
					this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
				}
			} else {
				this_.priv_JsDocInfoParser_next_7086()
			}
		}
		return typeNode
	} else {
		this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_BRACES())
		return this_.priv_JsDocInfoParser_parseTypeExpression_7057(token)
	}
}

// Java method: JsDocInfoParser.parseParamTypeExpression(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseParamTypeExpression_7052(token JsDocToken) pkg_rhino.Node {
	_ = token
	var restArg bool = false
	_ = restArg
	if token == StaticGet_JsDocToken_ITER_REST() {
		token = this_.priv_JsDocInfoParser_next_7086()
		if token == StaticGet_JsDocToken_RIGHT_CURLY() {
			this_.priv_JsDocInfoParser_restoreLookAhead_7082(token)
			return this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_ITER_REST(), this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_EMPTY()))
		}
		restArg = true
	}
	var typeNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTopLevelTypeExpression_7055(token)
	_ = typeNode
	if typeNode != 0 {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		if restArg {
			typeNode = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_ITER_REST(), typeNode)
		} else if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_EQUALS()) {
			this_.priv_JsDocInfoParser_next_7086()
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			typeNode = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_EQUALS(), typeNode)
		}
	}
	return typeNode
}

// Java method: JsDocInfoParser.parseParamTypeExpressionAnnotation(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseParamTypeExpressionAnnotation_7053(token JsDocToken) pkg_rhino.Node {
	_ = token
	if !(token == StaticGet_JsDocToken_LEFT_CURLY()) {
		j2gort.CheckArgument(false)
	}
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	var typeNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseParamTypeExpression_7052(this_.priv_JsDocInfoParser_next_7086())
	_ = typeNode
	if typeNode != 0 {
		if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
			this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
		} else {
			this_.priv_JsDocInfoParser_next_7086()
		}
	}
	return typeNode
}

// Java method: JsDocInfoParser.parseTypeNameAnnotation(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTypeNameAnnotation_7054(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token == StaticGet_JsDocToken_LEFT_CURLY() {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		var typeNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTypeName_7061(this_.priv_JsDocInfoParser_next_7086())
		_ = typeNode
		if typeNode != 0 {
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
				this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
			} else {
				this_.priv_JsDocInfoParser_next_7086()
			}
		}
		return typeNode
	} else {
		return this_.priv_JsDocInfoParser_parseTypeName_7061(token)
	}
}

// Java method: JsDocInfoParser.parseTopLevelTypeExpression(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTopLevelTypeExpression_7055(token JsDocToken) pkg_rhino.Node {
	_ = token
	var typeExpr pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTypeExpression_7057(token)
	_ = typeExpr
	if typeExpr != 0 {
		if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_PIPE()) {
			this_.priv_JsDocInfoParser_next_7086()
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			token = this_.priv_JsDocInfoParser_next_7086()
			return this_.priv_JsDocInfoParser_parseUnionTypeWithAlternate_7067(token, typeExpr)
		}
	}
	return typeExpr
}

// Java method: JsDocInfoParser.parseTypeExpressionList(java.lang.String typeName, com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTypeExpressionList_7056(typeName string, token JsDocToken) pkg_rhino.Node {
	_ = typeName
	_ = token
	var typeExpr pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTopLevelTypeExpression_7055(token)
	_ = typeExpr
	if typeExpr == 0 {
		return pkg_rhino.Node(0)
	}
	var typeList pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_BLOCK())
	_ = typeList
	var numTypeExprs int32 = int32(1)
	_ = numTypeExprs
	typeList.AddChildToBack(typeExpr)
	for this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COMMA()) {
		this_.priv_JsDocInfoParser_next_7086()
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		typeExpr = this_.priv_JsDocInfoParser_parseTopLevelTypeExpression_7055(this_.priv_JsDocInfoParser_next_7086())
		if typeExpr == 0 {
			return pkg_rhino.Node(0)
		}
		numTypeExprs = numTypeExprs + 1
		typeList.AddChildToBack(typeExpr)
	}
	if (j2gort.CheckNotNullString(typeName) == "Object") && (numTypeExprs == int32(1)) {
		typeList.AddChildToFront(this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_QMARK()))
	}
	return typeList
}

// Java method: JsDocInfoParser.parseTypeExpression(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTypeExpression_7057(token JsDocToken) pkg_rhino.Node {
	_ = token
	var lineno int32 = this_.stream.GetLineno()
	_ = lineno
	var charno int32 = this_.stream.GetCharno()
	_ = charno
	switch (token).Ordinal() {
	case int32(18) /* JsDocToken.QMARK */ :
		token = this_.priv_JsDocInfoParser_next_7086()
		if (((((((((token == StaticGet_JsDocToken_COMMA()) || (token == StaticGet_JsDocToken_EQUALS())) || (token == StaticGet_JsDocToken_RIGHT_SQUARE())) || (token == StaticGet_JsDocToken_RIGHT_CURLY())) || (token == StaticGet_JsDocToken_RIGHT_PAREN())) || (token == StaticGet_JsDocToken_PIPE())) || (token == StaticGet_JsDocToken_RIGHT_ANGLE())) || (token == StaticGet_JsDocToken_EOC())) || (token == StaticGet_JsDocToken_EOL())) || (token == StaticGet_JsDocToken_EOF()) {
			this_.priv_JsDocInfoParser_restoreLookAhead_7082(token)
			return this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_QMARK())
		}
		return this_.priv_JsDocInfoParser_wrapNode_7073(pkg_rhino.StaticGet_Token_QMARK(), this_.priv_JsDocInfoParser_parseBasicTypeExpression_7059(token), lineno, charno)
	case int32(19) /* JsDocToken.BANG */ :
		return this_.priv_JsDocInfoParser_wrapNode_7073(pkg_rhino.StaticGet_Token_BANG(), this_.priv_JsDocInfoParser_parseBasicTypeExpression_7059(this_.priv_JsDocInfoParser_next_7086()), lineno, charno)
	default:
		var basicTypeExpr pkg_rhino.Node = this_.priv_JsDocInfoParser_parseBasicTypeExpression_7059(token)
		_ = basicTypeExpr
		lineno = this_.stream.GetLineno()
		charno = this_.stream.GetCharno()
		if basicTypeExpr != 0 {
			if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_QMARK()) {
				this_.priv_JsDocInfoParser_next_7086()
				return this_.priv_JsDocInfoParser_wrapNode_7073(pkg_rhino.StaticGet_Token_QMARK(), basicTypeExpr, lineno, charno)
			} else if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_BANG()) {
				this_.priv_JsDocInfoParser_next_7086()
				return this_.priv_JsDocInfoParser_wrapNode_7073(pkg_rhino.StaticGet_Token_BANG(), basicTypeExpr, lineno, charno)
			}
		}
		return basicTypeExpr
	}
	panic("unreachable switch expression")
}

// Java method: JsDocInfoParser.parseContextTypeExpression(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseContextTypeExpression_7058(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token == StaticGet_JsDocToken_QMARK() {
		return this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_QMARK())
	} else {
		return this_.priv_JsDocInfoParser_parseBasicTypeExpression_7059(token)
	}
}

// Java method: JsDocInfoParser.parseBasicTypeExpression(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseBasicTypeExpression_7059(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token == StaticGet_JsDocToken_STAR() {
		return this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_STAR())
	} else if token == StaticGet_JsDocToken_LEFT_CURLY() {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		return this_.priv_JsDocInfoParser_parseRecordType_7068(this_.priv_JsDocInfoParser_next_7086())
	} else if token == StaticGet_JsDocToken_LEFT_PAREN() {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		return this_.priv_JsDocInfoParser_parseUnionType_7066(this_.priv_JsDocInfoParser_next_7086())
	} else if token == StaticGet_JsDocToken_STRING() {
		var string_ string = this_.stream.GetString()
		_ = string_
		switch j2gort.CheckNotNullString(string_) {
		case "function":
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			return this_.priv_JsDocInfoParser_parseFunctionType_7063(this_.priv_JsDocInfoParser_next_7086())
		case "null", "undefined":
			return this_.priv_JsDocInfoParser_newStringNode_7075(string_)
		case "typeof":
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			return this_.priv_JsDocInfoParser_parseTypeofType_7062(this_.priv_JsDocInfoParser_next_7086())
		default:
			return this_.priv_JsDocInfoParser_parseTypeName_7061(token)
		}
		panic("unreachable switch expression")
	}
	this_.priv_JsDocInfoParser_restoreLookAhead_7082(token)
	return this_.priv_JsDocInfoParser_reportGenericTypeSyntaxWarning_7077()
}

// Java method: JsDocInfoParser.parseNameExpression(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseNameExpression_7060(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token != StaticGet_JsDocToken_STRING() {
		this_.priv_JsDocInfoParser_addParserWarning_7011(pkg_rhino.StaticGet_Msg_JSDOC_NAME_SYNTAX(), this_.stream.GetLineno(), this_.stream.GetCharno())
		return pkg_rhino.Node(0)
	}
	var startLineno int32 = this_.stream.GetLineno()
	_ = startLineno
	var startCharno int32 = this_.stream.GetCharno()
	_ = startCharno
	var startOffset int32 = (this_.stream.GetCursor() - j2gort.JStrLen(this_.stream.GetString()))
	_ = startOffset
	var typeName j2gort.LocalStringBuilder
	typeName.InitStr(this_.stream.GetString())

	var endOffset int32 = this_.stream.GetCursor()
	_ = endOffset
	for this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_EOL()) && javalang.StringEndsWith(j2gort.CheckNotNullString(typeName.String()), ".") {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
			this_.priv_JsDocInfoParser_next_7086()
			endOffset = this_.stream.GetCursor()
			typeName.AppendStr(this_.stream.GetString())
		}
	}
	var str pkg_rhino.Node = this_.priv_JsDocInfoParser_newStringNode_7075(typeName.String())
	_ = str
	str.SetLinenoCharno(startLineno, startCharno)
	str.SetLength((endOffset - startOffset))
	return str
}

// Java method: JsDocInfoParser.parseTypeName(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTypeName_7061(token JsDocToken) pkg_rhino.Node {
	_ = token
	var typeNameNode pkg_rhino.Node = this_.priv_JsDocInfoParser_parseNameExpression_7060(token)
	_ = typeNameNode
	if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_LEFT_ANGLE()) {
		this_.priv_JsDocInfoParser_next_7086()
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		var memberType pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTypeExpressionList_7056(typeNameNode.GetString(), this_.priv_JsDocInfoParser_next_7086())
		_ = memberType
		if memberType != 0 {
			typeNameNode.AddChildToFront(memberType)
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_ANGLE())) {
				return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_GT())
			}
			this_.priv_JsDocInfoParser_next_7086()
		}
	}
	return typeNameNode
}

// Java method: JsDocInfoParser.parseTypeofType(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseTypeofType_7062(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token == StaticGet_JsDocToken_LEFT_CURLY() {
		return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_UNNECESSARY_BRACES())
	}
	var typeofType pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_TYPEOF())
	_ = typeofType
	var name pkg_rhino.Node = this_.priv_JsDocInfoParser_parseNameExpression_7060(token)
	_ = name
	if name == 0 {
		return pkg_rhino.Node(0)
	}
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	typeofType.AddChildToFront(name)
	return typeofType
}

// Java method: JsDocInfoParser.parseFunctionType(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseFunctionType_7063(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token != StaticGet_JsDocToken_LEFT_PAREN() {
		this_.priv_JsDocInfoParser_restoreLookAhead_7082(token)
		return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_LP())
	}
	var functionType pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_FUNCTION())
	_ = functionType
	var parameters pkg_rhino.Node = pkg_rhino.Node(0)
	_ = parameters
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_PAREN())) {
		token = this_.priv_JsDocInfoParser_next_7086()
		var hasParams bool = true
		_ = hasParams
		if token == StaticGet_JsDocToken_STRING() {
			var tokenStr string = this_.stream.GetString()
			_ = tokenStr
			var isThis bool = ("this" == tokenStr)
			_ = isThis
			var isNew bool = ("new" == tokenStr)
			_ = isNew
			if isThis || isNew {
				if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COLON()) {
					this_.priv_JsDocInfoParser_next_7086()
					this_.priv_JsDocInfoParser_skipEOLs_7088()
					var contextType pkg_rhino.Node = this_.priv_JsDocInfoParser_wrapNode_7072(func() pkg_rhino.Token {
						if isThis {
							return pkg_rhino.StaticGet_Token_THIS()
						}
						return pkg_rhino.StaticGet_Token_NEW()
					}(), this_.priv_JsDocInfoParser_parseContextTypeExpression_7058(this_.priv_JsDocInfoParser_next_7086()))
					_ = contextType
					if contextType == 0 {
						return pkg_rhino.Node(0)
					}
					functionType.AddChildToFront(contextType)
				} else {
					return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_COLON())
				}
				if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COMMA()) {
					this_.priv_JsDocInfoParser_next_7086()
					this_.priv_JsDocInfoParser_skipEOLs_7088()
					token = this_.priv_JsDocInfoParser_next_7086()
				} else {
					hasParams = false
				}
			}
		}
		if hasParams {
			parameters = this_.priv_JsDocInfoParser_parseParametersType_7064(token)
			if parameters == 0 {
				return pkg_rhino.Node(0)
			}
		}
	}
	if parameters != 0 {
		functionType.AddChildToBack(parameters)
	}
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_PAREN())) {
		return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RP())
	}
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	this_.priv_JsDocInfoParser_next_7086()
	var resultType pkg_rhino.Node = this_.priv_JsDocInfoParser_parseResultType_7065()
	_ = resultType
	if resultType == 0 {
		return pkg_rhino.Node(0)
	} else {
		functionType.AddChildToBack(resultType)
	}
	return functionType
}

// Java method: JsDocInfoParser.parseParametersType(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseParametersType_7064(token JsDocToken) pkg_rhino.Node {
	_ = token
	var paramsType pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_PARAM_LIST())
	_ = paramsType
	var isVarArgs bool = false
	_ = isVarArgs
	var paramType pkg_rhino.Node = pkg_rhino.Node(0)
	_ = paramType
	if token != StaticGet_JsDocToken_RIGHT_PAREN() {
		for doFirst_9 := true; doFirst_9 || (this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COMMA())); doFirst_9 = false {
			if paramType != 0 {
				this_.priv_JsDocInfoParser_next_7086()
				this_.priv_JsDocInfoParser_skipEOLs_7088()
				token = this_.priv_JsDocInfoParser_next_7086()
			}
			if token == StaticGet_JsDocToken_ITER_REST() {
				this_.priv_JsDocInfoParser_skipEOLs_7088()
				if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_PAREN()) {
					paramType = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_ITER_REST())
				} else {
					this_.priv_JsDocInfoParser_skipEOLs_7088()
					paramType = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_ITER_REST(), this_.priv_JsDocInfoParser_parseTypeExpression_7057(this_.priv_JsDocInfoParser_next_7086()))
					this_.priv_JsDocInfoParser_skipEOLs_7088()
				}
				isVarArgs = true
			} else {
				paramType = this_.priv_JsDocInfoParser_parseTypeExpression_7057(token)
				if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_EQUALS()) {
					this_.priv_JsDocInfoParser_skipEOLs_7088()
					this_.priv_JsDocInfoParser_next_7086()
					paramType = this_.priv_JsDocInfoParser_wrapNode_7072(pkg_rhino.StaticGet_Token_EQUALS(), paramType)
				}
			}
			if paramType == 0 {
				return pkg_rhino.Node(0)
			}
			paramsType.AddChildToBack(paramType)
			if isVarArgs {
				break
			}
		}
	}
	if isVarArgs && this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COMMA()) {
		return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_FUNCTION_VARARGS())
	}
	return paramsType
}

// Java method: JsDocInfoParser.parseResultType()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseResultType_7065() pkg_rhino.Node {
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COLON())) {
		return this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_EMPTY())
	}
	this_.priv_JsDocInfoParser_next_7086()
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) && func() bool {
		sr_10 := "void"
		se_11 := this_.stream.GetString()
		return j2gort.CheckNotNullString(sr_10) == se_11
	}() {
		this_.priv_JsDocInfoParser_next_7086()
		return this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_VOID())
	} else {
		return this_.priv_JsDocInfoParser_parseTypeExpression_7057(this_.priv_JsDocInfoParser_next_7086())
	}
}

// Java method: JsDocInfoParser.parseUnionType(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseUnionType_7066(token JsDocToken) pkg_rhino.Node {
	_ = token
	return this_.priv_JsDocInfoParser_parseUnionTypeWithAlternate_7067(token, pkg_rhino.Node(0))
}

// Java method: JsDocInfoParser.parseUnionTypeWithAlternate(com.google.javascript.jscomp.parsing.JsDocToken token, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node alternate)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseUnionTypeWithAlternate_7067(token JsDocToken, alternate pkg_rhino.Node) pkg_rhino.Node {
	_ = token
	_ = alternate
	var union pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_PIPE())
	_ = union
	if alternate != 0 {
		union.AddChildToBack(alternate)
	}
	var expr pkg_rhino.Node = pkg_rhino.Node(0)
	_ = expr
	for doFirst_12 := true; doFirst_12 || (this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_PIPE())); doFirst_12 = false {
		if expr != 0 {
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			token = this_.priv_JsDocInfoParser_next_7086()
			if !(token == StaticGet_JsDocToken_PIPE()) {
				j2gort.CheckState(false)
			}
			this_.priv_JsDocInfoParser_skipEOLs_7088()
			token = this_.priv_JsDocInfoParser_next_7086()
		}
		expr = this_.priv_JsDocInfoParser_parseTypeExpression_7057(token)
		if expr == 0 {
			return pkg_rhino.Node(0)
		}
		union.AddChildToBack(expr)
	}
	if alternate == 0 {
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_PAREN())) {
			return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RP())
		}
		this_.priv_JsDocInfoParser_next_7086()
	}
	if union.HasOneChild() {
		var firstChild pkg_rhino.Node = union.GetFirstChild()
		_ = firstChild
		firstChild.Detach()
		return firstChild
	}
	return union
}

// Java method: JsDocInfoParser.parseRecordType(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseRecordType_7068(token JsDocToken) pkg_rhino.Node {
	_ = token
	var recordType pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_LC())
	_ = recordType
	var fieldTypeList pkg_rhino.Node = this_.priv_JsDocInfoParser_parseFieldTypeList_7069(token)
	_ = fieldTypeList
	if fieldTypeList == 0 {
		return this_.priv_JsDocInfoParser_reportGenericTypeSyntaxWarning_7077()
	}
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY())) {
		return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_MISSING_RC())
	}
	this_.priv_JsDocInfoParser_next_7086()
	recordType.AddChildToBack(fieldTypeList)
	return recordType
}

// Java method: JsDocInfoParser.parseFieldTypeList(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseFieldTypeList_7069(token JsDocToken) pkg_rhino.Node {
	_ = token
	var fieldTypeList pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_LB())
	_ = fieldTypeList
	var names j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = names
	for doFirst_13 := true; doFirst_13 || (true); doFirst_13 = false {
		var fieldType pkg_rhino.Node = this_.priv_JsDocInfoParser_parseFieldType_7070(token)
		_ = fieldType
		if fieldType == 0 {
			return pkg_rhino.Node(0)
		}
		var name string
		_ = name
		if fieldType.IsStringKey() {
			name = fieldType.GetString()
		} else {
			name = fieldType.GetFirstChild().GetString()
		}
		if names.Add(name) {
			fieldTypeList.AddChildToBack(fieldType)
		} else {
			this_.priv_JsDocInfoParser_addTypeWarning_7012(pkg_rhino.StaticGet_Msg_JSDOC_TYPE_RECORD_DUPLICATE(), name)
		}
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COMMA())) {
			break
		}
		this_.priv_JsDocInfoParser_next_7086()
		this_.priv_JsDocInfoParser_skipEOLs_7088()
		if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_RIGHT_CURLY()) {
			break
		}
		token = this_.priv_JsDocInfoParser_next_7086()
	}
	return fieldTypeList
}

// Java method: JsDocInfoParser.parseFieldType(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseFieldType_7070(token JsDocToken) pkg_rhino.Node {
	_ = token
	var fieldName pkg_rhino.Node = this_.priv_JsDocInfoParser_parseFieldName_7071(token)
	_ = fieldName
	if fieldName == 0 {
		return pkg_rhino.Node(0)
	}
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_COLON())) {
		return fieldName
	}
	this_.priv_JsDocInfoParser_next_7086()
	this_.priv_JsDocInfoParser_skipEOLs_7088()
	var typeExpression pkg_rhino.Node = this_.priv_JsDocInfoParser_parseTypeExpression_7057(this_.priv_JsDocInfoParser_next_7086())
	_ = typeExpression
	if typeExpression == 0 {
		return pkg_rhino.Node(0)
	}
	var fieldType pkg_rhino.Node = this_.priv_JsDocInfoParser_newNode_7074(pkg_rhino.StaticGet_Token_COLON())
	_ = fieldType
	fieldType.AddChildToBack(fieldName)
	fieldType.AddChildToBack(typeExpression)
	return fieldType
}

// Java method: JsDocInfoParser.parseFieldName(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseFieldName_7071(token JsDocToken) pkg_rhino.Node {
	_ = token
	if token != StaticGet_JsDocToken_STRING() {
		return pkg_rhino.Node(0)
	}
	var s string = this_.stream.GetString()
	_ = s
	var n pkg_rhino.Node
	if this_.nodeAlloc_ != nil {
		if s == j2gort.NullString {
			j2gort.ThrowNPE()
		}
		n = pkg_rhino.Node(this_.nodeAlloc_.AllocStringNode(uint8(pkg_rhino.StaticGet_Token_STRING_KEY())-1, this_.nodeAlloc_.Intern(s), nodestore.Node(this_.templateNode))).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno())
	} else {
		n = pkg_rhino.Node_NewString_2(pkg_rhino.StaticGet_Token_STRING_KEY(), s).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
	}

	_ = n
	n.SetLength(j2gort.JStrLen(s))
	return n
}

// Java method: JsDocInfoParser.wrapNode(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_wrapNode_7072(type_ pkg_rhino.Token, n pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = n
	if n == 0 {
		return pkg_rhino.Node(0)
	} else {
		return this_.priv_JsDocInfoParser_wrapNode_7073(type_, n, n.GetLineno(), n.GetCharno())
	}
}

// Java method: JsDocInfoParser.wrapNode(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node n, int lineno, int charno)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_wrapNode_7073(type_ pkg_rhino.Token, n pkg_rhino.Node, lineno int32, charno int32) pkg_rhino.Node {
	if n == 0 {
		return pkg_rhino.Node(0)
	}
	if this_.nodeAlloc_ != nil {
		return pkg_rhino.Node(this_.nodeAlloc_.AllocNode2(uint8(type_)-1, nodestore.Node(n), nodestore.Node(this_.templateNode))).SetLinenoCharno(lineno, charno)
	}
	return pkg_rhino.NewNode_2(type_, n).SetLinenoCharno(lineno, charno).ClonePropsFrom(this_.templateNode)

}

// Java method: JsDocInfoParser.newNode(com.google.javascript.rhino.Token type)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_newNode_7074(type_ pkg_rhino.Token) pkg_rhino.Node {
	if this_.nodeAlloc_ != nil {
		return pkg_rhino.Node(this_.nodeAlloc_.AllocNode1(uint8(type_)-1, nodestore.Node(this_.templateNode))).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno())
	}
	return pkg_rhino.NewNode_1(type_).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)

}

// Java method: JsDocInfoParser.newStringNode(java.lang.String s)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_newStringNode_7075(s string) pkg_rhino.Node {
	var n pkg_rhino.Node
	if this_.nodeAlloc_ != nil {
		if s == j2gort.NullString {
			j2gort.ThrowNPE()
		}
		n = pkg_rhino.Node(this_.nodeAlloc_.AllocStringNode(uint8(pkg_rhino.StaticGet_Token_STRINGLIT())-1, this_.nodeAlloc_.Intern(s), nodestore.Node(this_.templateNode))).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno())
	} else {
		n = pkg_rhino.Node_NewString_1(s).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
	}
	n.SetLength(j2gort.JStrLen(s))
	return n

}

// Java method: JsDocInfoParser.reportTypeSyntaxWarning(com.google.javascript.rhino.Msg warning)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(warning pkg_rhino.Msg) pkg_rhino.Node {
	_ = warning
	this_.priv_JsDocInfoParser_addTypeWarning_7015(warning, this_.stream.GetLineno(), this_.stream.GetCharno())
	return pkg_rhino.Node(0)
}

// Java method: JsDocInfoParser.reportGenericTypeSyntaxWarning()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_reportGenericTypeSyntaxWarning_7077() pkg_rhino.Node {
	return this_.priv_JsDocInfoParser_reportTypeSyntaxWarning_7076(pkg_rhino.StaticGet_Msg_JSDOC_TYPE_SYNTAX())
}

// Java method: JsDocInfoParser.eatUntilEOLIfNotAnnotation()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7078() JsDocToken {
	return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(this_.priv_JsDocInfoParser_next_7086())
}

// Java method: JsDocInfoParser.eatUntilEOLIfNotAnnotation(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(token JsDocToken) JsDocToken {
	_ = token
	if token == StaticGet_JsDocToken_ANNOTATION() {
		this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
		return token
	}
	return this_.priv_JsDocInfoParser_eatTokensUntilEOL_7081(token)
}

// Java method: JsDocInfoParser.eatTokensUntilEOL()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_eatTokensUntilEOL_7080() JsDocToken {
	return this_.priv_JsDocInfoParser_eatTokensUntilEOL_7081(this_.priv_JsDocInfoParser_next_7086())
}

// Java method: JsDocInfoParser.eatTokensUntilEOL(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_eatTokensUntilEOL_7081(token JsDocToken) JsDocToken {
	_ = token
	for doFirst_14 := true; doFirst_14 || (true); doFirst_14 = false {
		if ((token == StaticGet_JsDocToken_EOL()) || (token == StaticGet_JsDocToken_EOC())) || (token == StaticGet_JsDocToken_EOF()) {
			this_.state = StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION()
			return token
		}
		token = this_.priv_JsDocInfoParser_next_7086()
	}
	panic("unreachable")
}

// Java method: JsDocInfoParser.restoreLookAhead(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_restoreLookAhead_7082(token JsDocToken) {
	_ = token
	this_.unreadToken = token
}

// Java method: JsDocInfoParser.match(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_match_7083(token JsDocToken) bool {
	_ = token
	this_.unreadToken = this_.priv_JsDocInfoParser_next_7086()
	return (this_.unreadToken == token)
}

// Java method: JsDocInfoParser.match(com.google.javascript.jscomp.parsing.JsDocToken token1, com.google.javascript.jscomp.parsing.JsDocToken token2)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_match_7084(token1 JsDocToken, token2 JsDocToken) bool {
	_ = token1
	_ = token2
	this_.unreadToken = this_.priv_JsDocInfoParser_next_7086()
	return ((this_.unreadToken == token1) || (this_.unreadToken == token2))
}

// Java method: JsDocInfoParser.eatIfMatch(com.google.javascript.jscomp.parsing.JsDocToken token)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_eatIfMatch_7085(token JsDocToken) bool {
	_ = token
	if this_.priv_JsDocInfoParser_match_7083(token) {
		this_.priv_JsDocInfoParser_next_7086()
		return true
	}
	return false
}

// Java method: JsDocInfoParser.next()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_next_7086() JsDocToken {
	if this_.unreadToken == StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN() {
		return this_.stream.GetJsDocToken()
	} else {
		return this_.priv_JsDocInfoParser_current_7087()
	}
}

// Java method: JsDocInfoParser.current()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_current_7087() JsDocToken {
	var t JsDocToken = this_.unreadToken
	_ = t
	this_.unreadToken = StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN()
	return t
}

// Java method: JsDocInfoParser.skipEOLs()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_skipEOLs_7088() {
	for this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_EOL()) {
		this_.priv_JsDocInfoParser_next_7086()
		if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STAR()) {
			this_.priv_JsDocInfoParser_next_7086()
		}
	}
}

// Java method: JsDocInfoParser.getRemainingJSDocLine()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_getRemainingJSDocLine_7089() string {
	var result string = this_.stream.GetRemainingJSDocLine()
	_ = result
	this_.unreadToken = StaticGet_JsDocInfoParser_NO_UNREAD_TOKEN()
	return result
}

// Java method: JsDocInfoParser.hasParsedFileOverviewDocInfo()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_hasParsedFileOverviewDocInfo_7090() bool {
	return this_.jsdocBuilder.IsPopulatedWithFileOverview()
}

// Java method: JsDocInfoParser.retrieveAndResetParsedJSDocInfo()
func (this_ *JsDocInfoParser) RetrieveAndResetParsedJSDocInfo() *pkg_rhino.JSDocInfo {
	return this_.jsdocBuilder.Build_0()
}

// Java method: JsDocInfoParser.getFileOverviewJSDocInfo()
func (this_ *JsDocInfoParser) GetFileOverviewJSDocInfo() *pkg_rhino.JSDocInfo {
	return this_.fileOverviewJSDocInfo
}

// Java method: JsDocInfoParser.lookAheadForType()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_lookAheadForType_7093() bool {
	return this_.priv_JsDocInfoParser_lookAheadFor_7095(uint16(123))
}

// Java method: JsDocInfoParser.lookAheadForAnnotation()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_lookAheadForAnnotation_7094() bool {
	return this_.priv_JsDocInfoParser_lookAheadFor_7095(uint16(64))
}

// Java method: JsDocInfoParser.lookAheadFor(char expect)
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_lookAheadFor_7095(expect uint16) bool {
	_ = expect
	var matched bool = false
	_ = matched
	var c int32 = int32(0)
	_ = c
	for {
		c = this_.stream.GetChar()
		if c == int32(uint16(32)) {
			continue
		} else if c == int32(expect) {
			matched = true
			break
		} else {
			break
		}
	}
	this_.stream.UngetChar(c)
	return matched
}

// Java method: JsDocInfoParser.parseClosureUnawareCode()
func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseClosureUnawareCode_7096() JsDocToken {
	this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID())
	if !(this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_LEFT_CURLY())) {
		if !(this_.jsdocBuilder.RecordClosureUnawareCode_0()) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA())
		}
		return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(this_.priv_JsDocInfoParser_current_7087())
	}
	var token JsDocToken = this_.priv_JsDocInfoParser_next_7086()
	_ = token
	if this_.priv_JsDocInfoParser_match_7083(StaticGet_JsDocToken_STRING()) {
		var value string = this_.stream.GetString()
		_ = value
		var mode pkg_rhino.JSDocInfo_PerFileClosureUnawareMode = func() pkg_rhino.JSDocInfo_PerFileClosureUnawareMode {
			switch j2gort.CheckNotNullString(value) {
			case "WHITESPACE":
				return pkg_rhino.StaticGet_JSDocInfo_PerFileClosureUnawareMode_WHITESPACE()
			case "SIMPLE":
				return pkg_rhino.StaticGet_JSDocInfo_PerFileClosureUnawareMode_SIMPLE()
			default:
				this_.priv_JsDocInfoParser_addParserWarning_7008(pkg_rhino.StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE(), value)
				return pkg_rhino.StaticGet_JSDocInfo_PerFileClosureUnawareMode_SIMPLE()
			}
			panic("unreachable switch expression")
		}()
		_ = mode
		if !(this_.jsdocBuilder.RecordClosureUnawareCode_1(mode)) {
			this_.priv_JsDocInfoParser_addParserWarning_7010(pkg_rhino.StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA())
		}
	}
	return this_.priv_JsDocInfoParser_eatUntilEOLIfNotAnnotation_7079(token)
}

// Java enum: com.google.javascript.jscomp.parsing.JsDocInfoParser.State
type JsDocInfoParser_State uint8

const (
	enum_JsDocInfoParser_State_SEARCHING_ANNOTATION JsDocInfoParser_State = 1
	enum_JsDocInfoParser_State_SEARCHING_NEWLINE    JsDocInfoParser_State = 2
	enum_JsDocInfoParser_State_NEXT_IS_ANNOTATION   JsDocInfoParser_State = 3
)

var enumBaseTable_JsDocInfoParser_State = [...]j2gort.EnumBase{
	{Name_: "SEARCHING_ANNOTATION", Ordinal_: 0},
	{Name_: "SEARCHING_NEWLINE", Ordinal_: 1},
	{Name_: "NEXT_IS_ANNOTATION", Ordinal_: 2},
}
var enumValues_JsDocInfoParser_State = []JsDocInfoParser_State{1, 2, 3}

func StaticGet_JsDocInfoParser_State_SEARCHING_ANNOTATION() JsDocInfoParser_State {
	return enum_JsDocInfoParser_State_SEARCHING_ANNOTATION
}
func StaticGet_JsDocInfoParser_State_SEARCHING_NEWLINE() JsDocInfoParser_State {
	return enum_JsDocInfoParser_State_SEARCHING_NEWLINE
}
func StaticGet_JsDocInfoParser_State_NEXT_IS_ANNOTATION() JsDocInfoParser_State {
	return enum_JsDocInfoParser_State_NEXT_IS_ANNOTATION
}

func JsDocInfoParser_State_Values() []JsDocInfoParser_State {
	return j2gort.CloneSlice(enumValues_JsDocInfoParser_State)
}

func JsDocInfoParser_State_ValueOf(name string) JsDocInfoParser_State {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsDocInfoParser_State {
		if enumBaseTable_JsDocInfoParser_State[i_].Name_ == name {
			return JsDocInfoParser_State(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.JsDocInfoParser.State." + name))
}

func (this_ JsDocInfoParser_State) IsJavaNil_() bool                              { return this_ == 0 }
func (JsDocInfoParser_State) HasIdentityEquals_()                                 {}
func (JsDocInfoParser_State) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JsDocInfoParser_State) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsDocInfoParser_State) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JsDocInfoParser_State) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocInfoParser$State", "State"
}
func (this_ JsDocInfoParser_State) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsDocInfoParser_State[this_-1]
}
func (this_ JsDocInfoParser_State) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsDocInfoParser_State) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsDocInfoParser_State) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocInfoParser_State[this_-1].Name_
}
func (this_ JsDocInfoParser_State) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsDocInfoParser_State[this_-1].Name_
}
func (this_ JsDocInfoParser_State) Equals(other_ any) bool {
	o_, ok_ := other_.(JsDocInfoParser_State)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsDocInfoParser_State) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsDocInfoParser_State[this_-1].HashCode()
}
func (this_ JsDocInfoParser_State) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocInfoParser_State[this_-1].CompareTo(other_)
}
func (this_ JsDocInfoParser_State) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsDocInfoParser_State[this_-1].Name_)
}
func init() {
	vals_ := []any{JsDocInfoParser_State(1), JsDocInfoParser_State(2), JsDocInfoParser_State(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocInfoParser.State", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocInfoParser$State", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocInfoParser.State", "com.google.javascript.jscomp.parsing.JsDocInfoParser$State", func(v_ any) bool {
		return j2gort.InstanceOf[JsDocInfoParser_State](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind
type JsDocInfoParser_JsDocSourceKind uint8

const (
	enum_JsDocInfoParser_JsDocSourceKind_NORMAL  JsDocInfoParser_JsDocSourceKind = 1
	enum_JsDocInfoParser_JsDocSourceKind_TSICKLE JsDocInfoParser_JsDocSourceKind = 2
)

var enumBaseTable_JsDocInfoParser_JsDocSourceKind = [...]j2gort.EnumBase{
	{Name_: "NORMAL", Ordinal_: 0},
	{Name_: "TSICKLE", Ordinal_: 1},
}
var enumValues_JsDocInfoParser_JsDocSourceKind = []JsDocInfoParser_JsDocSourceKind{1, 2}

func StaticGet_JsDocInfoParser_JsDocSourceKind_NORMAL() JsDocInfoParser_JsDocSourceKind {
	return enum_JsDocInfoParser_JsDocSourceKind_NORMAL
}
func StaticGet_JsDocInfoParser_JsDocSourceKind_TSICKLE() JsDocInfoParser_JsDocSourceKind {
	return enum_JsDocInfoParser_JsDocSourceKind_TSICKLE
}

func JsDocInfoParser_JsDocSourceKind_Values() []JsDocInfoParser_JsDocSourceKind {
	return j2gort.CloneSlice(enumValues_JsDocInfoParser_JsDocSourceKind)
}

func JsDocInfoParser_JsDocSourceKind_ValueOf(name string) JsDocInfoParser_JsDocSourceKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsDocInfoParser_JsDocSourceKind {
		if enumBaseTable_JsDocInfoParser_JsDocSourceKind[i_].Name_ == name {
			return JsDocInfoParser_JsDocSourceKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind." + name))
}

func (this_ JsDocInfoParser_JsDocSourceKind) IsJavaNil_() bool                              { return this_ == 0 }
func (JsDocInfoParser_JsDocSourceKind) HasIdentityEquals_()                                 {}
func (JsDocInfoParser_JsDocSourceKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JsDocInfoParser_JsDocSourceKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsDocInfoParser_JsDocSourceKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JsDocInfoParser_JsDocSourceKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocInfoParser$JsDocSourceKind", "JsDocSourceKind"
}
func (this_ JsDocInfoParser_JsDocSourceKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsDocInfoParser_JsDocSourceKind[this_-1]
}
func (this_ JsDocInfoParser_JsDocSourceKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsDocInfoParser_JsDocSourceKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsDocInfoParser_JsDocSourceKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocInfoParser_JsDocSourceKind[this_-1].Name_
}
func (this_ JsDocInfoParser_JsDocSourceKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsDocInfoParser_JsDocSourceKind[this_-1].Name_
}
func (this_ JsDocInfoParser_JsDocSourceKind) Equals(other_ any) bool {
	o_, ok_ := other_.(JsDocInfoParser_JsDocSourceKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsDocInfoParser_JsDocSourceKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsDocInfoParser_JsDocSourceKind[this_-1].HashCode()
}
func (this_ JsDocInfoParser_JsDocSourceKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocInfoParser_JsDocSourceKind[this_-1].CompareTo(other_)
}
func (this_ JsDocInfoParser_JsDocSourceKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsDocInfoParser_JsDocSourceKind[this_-1].Name_)
}
func init() {
	vals_ := []any{JsDocInfoParser_JsDocSourceKind(1), JsDocInfoParser_JsDocSourceKind(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocInfoParser$JsDocSourceKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind", "com.google.javascript.jscomp.parsing.JsDocInfoParser$JsDocSourceKind", func(v_ any) bool {
		return j2gort.InstanceOf[JsDocInfoParser_JsDocSourceKind](v_)
	})
}

// Java class: com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtractionInfo
type JsDocInfoParser_ExtractionInfo struct {
	string_ string     `gson:"string"`
	token   JsDocToken `gson:"token"`
}

func (this_ JsDocInfoParser_ExtractionInfo) BindSelf_(s any) {
	_ = s
}
func (this_ JsDocInfoParser_ExtractionInfo) IsJavaNil_() bool {
	return this_ == (JsDocInfoParser_ExtractionInfo{})
}
func (JsDocInfoParser_ExtractionInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocInfoParser$ExtractionInfo", "ExtractionInfo"
}
func (JsDocInfoParser_ExtractionInfo) HasIdentityEquals_() {}
func (this_ JsDocInfoParser_ExtractionInfo) IdentityHash_() int32 {
	if this_ == (JsDocInfoParser_ExtractionInfo{}) {
		return 0
	}
	var h_ int32 = 1
	{
		hs_15 := this_.string_
		h_ = 31 * h_
		if hs_15 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_15)
		}
	}
	h_ = 31*h_ + (this_.token).HashCode()
	return h_
}
func (this_ JsDocInfoParser_ExtractionInfo) FastValHash_() (int32, bool, bool) {
	if this_.IsJavaNil_() {
		return 0, false, true
	}
	return this_.IdentityHash_(), true, true
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtractionInfo", "com.google.javascript.jscomp.parsing.JsDocInfoParser$ExtractionInfo", func(v_ any) bool {
		return j2gort.InstanceOf[JsDocInfoParser_ExtractionInfo](v_)
	})
}

func (this_ *JsDocInfoParser_ExtractionInfo) initFields_JsDocInfoParser_ExtractionInfo() {
	_ = this_
	this_.string_ = j2gort.NullString
}

// Java constructor: ExtractionInfo(java.lang.String string, com.google.javascript.jscomp.parsing.JsDocToken token)
func NewJsDocInfoParser_ExtractionInfo(string_ string, token JsDocToken) JsDocInfoParser_ExtractionInfo {
	this_ := JsDocInfoParser_ExtractionInfo{}
	InitJsDocInfoParser_ExtractionInfo(&this_, string_, token)
	return this_
}

func InitJsDocInfoParser_ExtractionInfo(this_ *JsDocInfoParser_ExtractionInfo, string_ string, token JsDocToken) {
	_ = string_
	_ = token
	this_.initFields_JsDocInfoParser_ExtractionInfo()
	this_.string_ = string_
	this_.token = token
}

// Java class: com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo
type JsDocInfoParser_ExtendedTypeInfo struct {
	type_  pkg_rhino.JSTypeExpression `gson:"type"`
	lineno int32                      `gson:"lineno"`
	charno int32                      `gson:"charno"`
}

func (this_ JsDocInfoParser_ExtendedTypeInfo) BindSelf_(s any) {
	_ = s
}
func (this_ JsDocInfoParser_ExtendedTypeInfo) IsJavaNil_() bool {
	return this_ == (JsDocInfoParser_ExtendedTypeInfo{})
}
func (JsDocInfoParser_ExtendedTypeInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocInfoParser$ExtendedTypeInfo", "ExtendedTypeInfo"
}
func (JsDocInfoParser_ExtendedTypeInfo) HasIdentityEquals_() {}
func (this_ JsDocInfoParser_ExtendedTypeInfo) IdentityHash_() int32 {
	if this_ == (JsDocInfoParser_ExtendedTypeInfo{}) {
		return 0
	}
	var h_ int32 = 1
	h_ = 31*h_ + (this_.type_).IdentityHash_()
	h_ = 31*h_ + (this_.lineno)
	h_ = 31*h_ + (this_.charno)
	return h_
}
func (this_ JsDocInfoParser_ExtendedTypeInfo) FastValHash_() (int32, bool, bool) {
	if this_.IsJavaNil_() {
		return 0, false, true
	}
	return this_.IdentityHash_(), true, true
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo", "com.google.javascript.jscomp.parsing.JsDocInfoParser$ExtendedTypeInfo", func(v_ any) bool {
		return j2gort.InstanceOf[JsDocInfoParser_ExtendedTypeInfo](v_)
	})
}

func (this_ *JsDocInfoParser_ExtendedTypeInfo) initFields_JsDocInfoParser_ExtendedTypeInfo() {
	_ = this_
}

// Java constructor: ExtendedTypeInfo(com.google.javascript.rhino.JSTypeExpression type, int lineno, int charno)
func NewJsDocInfoParser_ExtendedTypeInfo(type_ pkg_rhino.JSTypeExpression, lineno int32, charno int32) JsDocInfoParser_ExtendedTypeInfo {
	this_ := JsDocInfoParser_ExtendedTypeInfo{}
	InitJsDocInfoParser_ExtendedTypeInfo(&this_, type_, lineno, charno)
	return this_
}

func InitJsDocInfoParser_ExtendedTypeInfo(this_ *JsDocInfoParser_ExtendedTypeInfo, type_ pkg_rhino.JSTypeExpression, lineno int32, charno int32) {
	_ = type_
	_ = lineno
	_ = charno
	this_.initFields_JsDocInfoParser_ExtendedTypeInfo()
	this_.type_ = type_
	this_.lineno = lineno
	this_.charno = charno
}

// Java enum: com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption
type JsDocInfoParser_WhitespaceOption uint8

const (
	enum_JsDocInfoParser_WhitespaceOption_PRESERVE    JsDocInfoParser_WhitespaceOption = 1
	enum_JsDocInfoParser_WhitespaceOption_TRIM        JsDocInfoParser_WhitespaceOption = 2
	enum_JsDocInfoParser_WhitespaceOption_SINGLE_LINE JsDocInfoParser_WhitespaceOption = 3
)

var enumBaseTable_JsDocInfoParser_WhitespaceOption = [...]j2gort.EnumBase{
	{Name_: "PRESERVE", Ordinal_: 0},
	{Name_: "TRIM", Ordinal_: 1},
	{Name_: "SINGLE_LINE", Ordinal_: 2},
}
var enumValues_JsDocInfoParser_WhitespaceOption = []JsDocInfoParser_WhitespaceOption{1, 2, 3}

func StaticGet_JsDocInfoParser_WhitespaceOption_PRESERVE() JsDocInfoParser_WhitespaceOption {
	return enum_JsDocInfoParser_WhitespaceOption_PRESERVE
}
func StaticGet_JsDocInfoParser_WhitespaceOption_TRIM() JsDocInfoParser_WhitespaceOption {
	return enum_JsDocInfoParser_WhitespaceOption_TRIM
}
func StaticGet_JsDocInfoParser_WhitespaceOption_SINGLE_LINE() JsDocInfoParser_WhitespaceOption {
	return enum_JsDocInfoParser_WhitespaceOption_SINGLE_LINE
}

func JsDocInfoParser_WhitespaceOption_Values() []JsDocInfoParser_WhitespaceOption {
	return j2gort.CloneSlice(enumValues_JsDocInfoParser_WhitespaceOption)
}

func JsDocInfoParser_WhitespaceOption_ValueOf(name string) JsDocInfoParser_WhitespaceOption {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsDocInfoParser_WhitespaceOption {
		if enumBaseTable_JsDocInfoParser_WhitespaceOption[i_].Name_ == name {
			return JsDocInfoParser_WhitespaceOption(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption." + name))
}

func (this_ JsDocInfoParser_WhitespaceOption) IsJavaNil_() bool                              { return this_ == 0 }
func (JsDocInfoParser_WhitespaceOption) HasIdentityEquals_()                                 {}
func (JsDocInfoParser_WhitespaceOption) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JsDocInfoParser_WhitespaceOption) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsDocInfoParser_WhitespaceOption) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JsDocInfoParser_WhitespaceOption) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.JsDocInfoParser$WhitespaceOption", "WhitespaceOption"
}
func (this_ JsDocInfoParser_WhitespaceOption) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsDocInfoParser_WhitespaceOption[this_-1]
}
func (this_ JsDocInfoParser_WhitespaceOption) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsDocInfoParser_WhitespaceOption) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsDocInfoParser_WhitespaceOption) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocInfoParser_WhitespaceOption[this_-1].Name_
}
func (this_ JsDocInfoParser_WhitespaceOption) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsDocInfoParser_WhitespaceOption[this_-1].Name_
}
func (this_ JsDocInfoParser_WhitespaceOption) Equals(other_ any) bool {
	o_, ok_ := other_.(JsDocInfoParser_WhitespaceOption)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsDocInfoParser_WhitespaceOption) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsDocInfoParser_WhitespaceOption[this_-1].HashCode()
}
func (this_ JsDocInfoParser_WhitespaceOption) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocInfoParser_WhitespaceOption[this_-1].CompareTo(other_)
}
func (this_ JsDocInfoParser_WhitespaceOption) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsDocInfoParser_WhitespaceOption[this_-1].Name_)
}
func init() {
	vals_ := []any{JsDocInfoParser_WhitespaceOption(1), JsDocInfoParser_WhitespaceOption(2), JsDocInfoParser_WhitespaceOption(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocInfoParser$WhitespaceOption", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocInfoParser.WhitespaceOption", "com.google.javascript.jscomp.parsing.JsDocInfoParser$WhitespaceOption", func(v_ any) bool {
		return j2gort.InstanceOf[JsDocInfoParser_WhitespaceOption](v_)
	})
}
