/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/Parser.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.Parser
type Parser struct {
	scanner                              *Scanner                                                    `gson:"scanner"`
	errorReporter                        ErrorReporter                                               `gson:"errorReporter"`
	config                               *Parser_Config                                              `gson:"config"`
	comments                             *j2gort.Ext_com_google_common_collect_ImmutableList_Builder `gson:"comments"`
	features                             *FeatureSet                                                 `gson:"features"`
	sourceMapURL                         string                                                      `gson:"sourceMapURL"`
	treeScratch                          []ParseTree                                                 `gson:"treeScratch"`
	commaScratch                         []SourcePosition                                            `gson:"commaScratch"`
	lastCommentEndPosition               SourcePosition                                              `gson:"lastCommentEndPosition"`
	lastSourcePosition                   SourcePosition                                              `gson:"lastSourcePosition"`
	treeScratchTop                       int32                                                       `gson:"treeScratchTop"`
	commaScratchTop                      int32                                                       `gson:"commaScratchTop"`
	idHash_                              atomic.Int32
	inGenerator                          bool `gson:"inGenerator"`
	slab_ImportDeclarationTree           j2gort.PooledSlab[ImportDeclarationTree]
	slab_ImportSpecifierTree             j2gort.PooledSlab[ImportSpecifierTree]
	slab_ExportDeclarationTree           j2gort.PooledSlab[ExportDeclarationTree]
	slab_ExportSpecifierTree             j2gort.PooledSlab[ExportSpecifierTree]
	slab_ClassDeclarationTree            j2gort.PooledSlab[ClassDeclarationTree]
	slab_FieldDeclarationTree            j2gort.PooledSlab[FieldDeclarationTree]
	slab_ComputedPropertyFieldTree       j2gort.PooledSlab[ComputedPropertyFieldTree]
	slab_FunctionDeclarationTree_Builder j2gort.PooledSlab[FunctionDeclarationTree_Builder]
	slab_FunctionDeclarationTree         j2gort.PooledSlab[FunctionDeclarationTree]
	slab_ComputedPropertyMethodTree      j2gort.PooledSlab[ComputedPropertyMethodTree]
	slab_DefaultParameterTree            j2gort.PooledSlab[DefaultParameterTree]
	slab_FormalParameterListTree         j2gort.PooledSlab[FormalParameterListTree]
	slab_BlockTree                       j2gort.PooledSlab[BlockTree]
	slab_IterSpreadTree                  j2gort.PooledSlab[IterSpreadTree]
	slab_VariableStatementTree           j2gort.PooledSlab[VariableStatementTree]
	slab_VariableDeclarationListTree     j2gort.PooledSlab[VariableDeclarationListTree]
	slab_VariableDeclarationTree         j2gort.PooledSlab[VariableDeclarationTree]
	slab_EmptyStatementTree              j2gort.PooledSlab[EmptyStatementTree]
	slab_ExpressionStatementTree         j2gort.PooledSlab[ExpressionStatementTree]
	slab_IfStatementTree                 j2gort.PooledSlab[IfStatementTree]
	slab_DoWhileStatementTree            j2gort.PooledSlab[DoWhileStatementTree]
	slab_WhileStatementTree              j2gort.PooledSlab[WhileStatementTree]
	slab_ForOfStatementTree              j2gort.PooledSlab[ForOfStatementTree]
	slab_ForAwaitOfStatementTree         j2gort.PooledSlab[ForAwaitOfStatementTree]
	slab_NullTree                        j2gort.PooledSlab[NullTree]
	slab_ForStatementTree                j2gort.PooledSlab[ForStatementTree]
	slab_ForInStatementTree              j2gort.PooledSlab[ForInStatementTree]
	slab_ContinueStatementTree           j2gort.PooledSlab[ContinueStatementTree]
	slab_BreakStatementTree              j2gort.PooledSlab[BreakStatementTree]
	slab_ReturnStatementTree             j2gort.PooledSlab[ReturnStatementTree]
	slab_WithStatementTree               j2gort.PooledSlab[WithStatementTree]
	slab_SwitchStatementTree             j2gort.PooledSlab[SwitchStatementTree]
	slab_CaseClauseTree                  j2gort.PooledSlab[CaseClauseTree]
	slab_DefaultClauseTree               j2gort.PooledSlab[DefaultClauseTree]
	slab_LabelledStatementTree           j2gort.PooledSlab[LabelledStatementTree]
	slab_ThrowStatementTree              j2gort.PooledSlab[ThrowStatementTree]
	slab_TryStatementTree                j2gort.PooledSlab[TryStatementTree]
	slab_CatchTree                       j2gort.PooledSlab[CatchTree]
	slab_FinallyTree                     j2gort.PooledSlab[FinallyTree]
	slab_DebuggerStatementTree           j2gort.PooledSlab[DebuggerStatementTree]
	slab_SuperExpressionTree             j2gort.PooledSlab[SuperExpressionTree]
	slab_ThisExpressionTree              j2gort.PooledSlab[ThisExpressionTree]
	slab_DynamicImportTree               j2gort.PooledSlab[DynamicImportTree]
	slab_IdentifierExpressionTree        j2gort.PooledSlab[IdentifierExpressionTree]
	slab_LiteralExpressionTree           j2gort.PooledSlab[LiteralExpressionTree]
	slab_TemplateLiteralPortionTree      j2gort.PooledSlab[TemplateLiteralPortionTree]
	slab_TemplateLiteralExpressionTree   j2gort.PooledSlab[TemplateLiteralExpressionTree]

	// Java constructor: Parser(com.google.javascript.jscomp.parsing.parser.Parser.Config config, com.google.javascript.jscomp.parsing.parser.util.ErrorReporter errorReporter, com.google.javascript.jscomp.parsing.parser.SourceFile source)
	slab_TemplateSubstitutionTree       j2gort.PooledSlab[TemplateSubstitutionTree]
	slab_ArrayLiteralExpressionTree     j2gort.PooledSlab[ArrayLiteralExpressionTree]
	slab_ObjectLiteralExpressionTree    j2gort.PooledSlab[ObjectLiteralExpressionTree]
	slab_ObjectSpreadTree               j2gort.PooledSlab[ObjectSpreadTree]
	slab_ComputedPropertyDefinitionTree j2gort.PooledSlab[ComputedPropertyDefinitionTree]
	slab_GetAccessorTree                j2gort.PooledSlab[GetAccessorTree]
	slab_ComputedPropertyGetterTree     j2gort.PooledSlab[ComputedPropertyGetterTree]
	slab_SetAccessorTree                j2gort.PooledSlab[SetAccessorTree]
	slab_ComputedPropertySetterTree     j2gort.PooledSlab[ComputedPropertySetterTree]
	slab_PropertyNameAssignmentTree     j2gort.PooledSlab[PropertyNameAssignmentTree]
	slab_MissingPrimaryExpressionTree   j2gort.PooledSlab[MissingPrimaryExpressionTree]
	slab_CommaExpressionTree            j2gort.

		// Java method: Parser.recordComment(com.google.javascript.jscomp.parsing.parser.trees.Comment.Type type, com.google.javascript.jscomp.parsing.parser.util.SourceRange range, java.lang.String value)
		PooledSlab[CommaExpressionTree]
	slab_ParenExpressionTree                j2gort.PooledSlab[ParenExpressionTree]
	slab_BinaryOperatorTree                 j2gort.PooledSlab[BinaryOperatorTree]
	slab_YieldExpressionTree                j2gort.PooledSlab[YieldExpressionTree]
	slab_ConditionalExpressionTree          j2gort.PooledSlab[ConditionalExpressionTree]
	slab_UnaryExpressionTree                j2gort.PooledSlab[UnaryExpressionTree]
	slab_AwaitExpressionTree                j2gort.PooledSlab[AwaitExpressionTree]
	slab_UpdateExpressionTree               j2gort.PooledSlab[UpdateExpressionTree]
	slab_CallExpressionTree                 j2gort.PooledSlab[CallExpressionTree]
	slab_MemberLookupExpressionTree         j2gort.PooledSlab[MemberLookupExpressionTree]
	slab_MemberExpressionTree               j2gort.PooledSlab[MemberExpressionTree]
	slab_OptChainCallExpressionTree         j2gort.PooledSlab[OptChainCallExpressionTree]
	slab_OptionalMemberLookupExpressionTree j2gort.PooledSlab[OptionalMemberLookupExpressionTree]
	slab_OptionalMemberExpressionTree       j2gort.PooledSlab[OptionalMemberExpressionTree]
	slab_NewExpressionTree                  j2gort.PooledSlab[NewExpressionTree]
	slab_NewTargetExpressionTree            j2gort.PooledSlab[NewTargetExpressionTree]
	slab_ImportMetaExpressionTree           j2gort.

		// Java method: Parser.buildComments()
		PooledSlab[ImportMetaExpressionTree]
	slab_ArgumentListTree  j2gort.PooledSlab[ArgumentListTree]
	slab_IterRestTree      j2gort.PooledSlab[IterRestTree]
	slab_ArrayPatternTree  j2gort.PooledSlab[ArrayPatternTree]
	slab_ObjectRestTree    j2gort.PooledSlab[ObjectRestTree]
	slab_ObjectPatternTree j2gort.PooledSlab[

	// Java method: Parser.getComments()
	ObjectPatternTree]
}

func (this_ *Parser) BindSelf_(s any) {
	_ = s
}
func (this_ *Parser) IsJavaNil_() bool { return this_ == nil }
func (*Parser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Parser", "Parser"
}
func (this_ *Parser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Parser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Parser", "com.google.javascript.jscomp.parsing.parser.Parser", func(v_ any) bool {
		return j2gort.InstanceOf[*Parser](v_)
	})
}
func (*Parser) Implements_parser_Scanner_CommentRecorder() {}

func StaticGet_Parser_SOURCE_MAPPING_URL_PREFIX() string { return "//# sourceMappingURL=" }

var static_Parser_EMPTY_PARSE_TREES j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_Parser_EMPTY_PARSE_TREES() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_Parser()
	return static_Parser_EMPTY_PARSE_TREES
}
func StaticSet_Parser_EMPTY_PARSE_TREES(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_Parser()
	static_Parser_EMPTY_PARSE_TREES = v
}

var static_Parser_EMPTY_PARSE_TREE_ARRAY []ParseTree

func StaticGet_Parser_EMPTY_PARSE_TREE_ARRAY() []ParseTree {
	clinit_Parser()
	return static_Parser_EMPTY_PARSE_TREE_ARRAY
}
func StaticSet_Parser_EMPTY_PARSE_TREE_ARRAY(v []ParseTree) {
	clinit_Parser()
	static_Parser_EMPTY_PARSE_TREE_ARRAY = v
}

var static_Parser_EMPTY_SOURCE_POSITION_ARRAY []SourcePosition

func StaticGet_Parser_EMPTY_SOURCE_POSITION_ARRAY() []SourcePosition {
	clinit_Parser()
	return static_Parser_EMPTY_SOURCE_POSITION_ARRAY
}
func StaticSet_Parser_EMPTY_SOURCE_POSITION_ARRAY(v []SourcePosition) {
	clinit_Parser()
	static_Parser_EMPTY_SOURCE_POSITION_ARRAY = v
}
func StaticGet_Parser_AWAIT() string { return "await" }
func StaticGet_Parser_ASYNC() string { return "async" }

var clinitDone_Parser uint32
var clinitState_Parser j2gort.ClinitState

func clinit_Parser() {
	if atomic.LoadUint32(&clinitDone_Parser) != 0 {
		return
	}
	clinitSlow_Parser()
}

func clinitSlow_Parser() {
	if !j2gort.EnterClinit(&clinitState_Parser, &clinitDone_Parser) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Parser, &clinitDone_Parser)
	static_Parser_EMPTY_PARSE_TREES = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_Parser_EMPTY_PARSE_TREE_ARRAY = make([]ParseTree, int32(0))
	static_Parser_EMPTY_SOURCE_POSITION_ARRAY = make([]SourcePosition, int32(0))
}

func (this_ *Parser) initFields_Parser() {
	_ = this_

	this_.features = StaticGet_FeatureSet_BARE_MINIMUM()
	this_.sourceMapURL = j2gort.NullString
	this_.treeScratch = make([]ParseTree, int32(32))

	this_.commaScratch = make([]SourcePosition, int32(32))

}

func NewParser(config *Parser_Config, errorReporter ErrorReporter, source *SourceFile) *Parser {
	clinit_Parser()
	this_ := &Parser{}
	InitParser(this_, config, errorReporter, source)
	return this_
}

func InitParser(this_ *Parser, config *Parser_Config, errorReporter ErrorReporter, source *SourceFile) {
	_ = config
	_ = errorReporter
	_ = source
	this_.initFields_Parser()
	this_.config = config
	this_.errorReporter = errorReporter
	this_.scanner = NewScanner(errorReporter, ToIface_Scanner_CommentRecorder(this_), source, int32(0))
	this_.lastSourcePosition = this_.scanner.GetPosition_0()
}

func (this_ *Parser) RecordComment(type_ Comment_Type, range_ SourceRange, value string) {
	_ = type_
	_ = range_
	_ = value
	if (this_.lastCommentEndPosition == (SourcePosition{})) || (range_.end.offset > this_.lastCommentEndPosition.offset) {
		value = javalang.StringTrim(j2gort.CheckNotNullString(value))
		if javalang.StringStartsWith(j2gort.CheckNotNullString(value), "//# sourceMappingURL=") {
			this_.sourceMapURL = j2gort.JStrSubFrom(value, j2gort.JStrLen("//# sourceMappingURL="))
		}
		var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = this_.comments
		_ = builder
		if builder == nil {
			builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
			this_.comments = builder
		}
		builder.AddVoid_(NewComment(value, range_, type_))
		this_.lastCommentEndPosition = range_.end
	}
}

func (this_ *Parser) priv_Parser_buildComments_7212() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.comments != nil {
		return this_.comments.BuildImmList_()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
}

func (this_ *Parser) GetComments() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](this_.priv_Parser_buildComments_7212())
}

// Java method: Parser.getFeatures()
func (this_ *Parser) GetFeatures() *FeatureSet {
	return this_.features
}

// Java method: Parser.getSourceMapURL()
func (this_ *Parser) GetSourceMapURL() string {
	return this_.sourceMapURL
}

// Java method: Parser.parseProgram()
func (this_ *Parser) ParseProgram() *ProgramTree {
	var didRet_1 bool
	var retVal_2 *ProgramTree
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javalang.IError](jex_4) {
					e := j2gort.AsIface[javalang.IError](jex_4)
					_ = e
					if j2gort.JStrIndexOf(e.Error(), "java.lang.StackOverflowError") >= 0 {
						this_.priv_Parser_reportError_7441("Too deep recursion while parsing", []any{})
						retVal_2 = nil
						didRet_1 = true
						return
					}
					panic(j2gort.CheckNotNull(e))
				} else {
					panic(rec_3)
				}
			}
		}()
		var start SourcePosition = this_.lastSourcePosition
		_ = start
		var sourceElements j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_Parser_parseGlobalSourceElements_7217()
		_ = sourceElements
		this_.priv_Parser_eat_7425(StaticGet_TokenType_END_OF_FILE())
		retVal_2 = NewProgramTree(this_.priv_Parser_getTreeLocation_7429(start), sourceElements, this_.priv_Parser_buildComments_7212())
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java method: Parser.parseGlobalSourceElements()
func (this_ *Parser) priv_Parser_parseGlobalSourceElements_7217() j2gort.Ext_com_google_common_collect_ImmutableList {
	var result []ParseTree = []ParseTree{}
	_ = result
	for !(this_.priv_Parser_peek_7433(StaticGet_TokenType_END_OF_FILE())) {
		result = append(result, j2gort.CheckNotNullElem(this_.priv_Parser_parseScriptElement_7218()))
	}
	return j2gort.ImmListFromTypedSlice(result)
}

// Java method: Parser.parseScriptElement()
func (this_ *Parser) priv_Parser_parseScriptElement_7218() ParseTree {
	if this_.priv_Parser_peekImportDeclaration_7219() {
		return this_.priv_Parser_parseImportDeclaration_7220()
	}
	if this_.priv_Parser_peekExportDeclaration_7223() {
		return this_.priv_Parser_parseExportDeclaration_7224()
	}
	return this_.priv_Parser_parseSourceElement_7248()
}

// Java method: Parser.peekImportDeclaration()
func (this_ *Parser) priv_Parser_peekImportDeclaration_7219() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_IMPORT()) && (((this_.priv_Parser_peekIdOrKeyword_7420(int32(1)) || this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_STRING())) || this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_OPEN_CURLY())) || this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_STAR())))
}

// Java method: Parser.parseImportDeclaration()
func (this_ *Parser) priv_Parser_parseImportDeclaration_7220() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_IMPORT())
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_STRING()) {
		var moduleSpecifier Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_STRING())
		_ = moduleSpecifier
		this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
		return ToIface_ParseTree(allocPooled_ImportDeclarationTree(&this_.slab_ImportDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start), (Token{}), nil, (Token{}), moduleSpecifier))
	}
	var defaultBindingIdentifier Token = (Token{})
	_ = defaultBindingIdentifier
	var nameSpaceImportIdentifier Token = (Token{})
	_ = nameSpaceImportIdentifier
	var identifierSet j2gort.Ext_com_google_common_collect_ImmutableList = nil
	_ = identifierSet
	var parseExplicitNames bool = true
	_ = parseExplicitNames
	if this_.priv_Parser_peekId_7417() {
		defaultBindingIdentifier = this_.priv_Parser_eatId_7422()
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
		} else {
			parseExplicitNames = false
		}
	} else if Keywords_IsKeyword_1_TokenType(this_.priv_Parser_peekType_7435()) {
		var keyword Token = this_.priv_Parser_nextToken_7430()
		_ = keyword
		this_.priv_Parser_reportError_7439(keyword, "cannot use keyword '%s' here.", []any{j2gort.BoxNullableVal(keyword)})
	}
	if parseExplicitNames {
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_STAR()) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_STAR())
			this_.priv_Parser_eatPredefinedString_7326("as")
			nameSpaceImportIdentifier = this_.priv_Parser_eatId_7422()
		} else {
			identifierSet = this_.priv_Parser_parseImportSpecifierSet_7221()
		}
	}
	this_.priv_Parser_eatPredefinedString_7326("from")
	var moduleSpecifier_2 Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_STRING())
	_ = moduleSpecifier_2
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return ToIface_ParseTree(allocPooled_ImportDeclarationTree(&this_.slab_ImportDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start), defaultBindingIdentifier, identifierSet, nameSpaceImportIdentifier, moduleSpecifier_2))
}

// Java method: Parser.parseImportSpecifierSet()
func (this_ *Parser) priv_Parser_parseImportSpecifierSet_7221() j2gort.Ext_com_google_common_collect_ImmutableList {
	var elements []ParseTree = []ParseTree{}
	_ = elements
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	for this_.priv_Parser_peekIdOrKeyword_7419() {
		elements = append(elements, j2gort.CheckNotNullElem(this_.priv_Parser_parseImportSpecifier_7222()))
		if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_CURLY())) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
		}
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return j2gort.ImmListFromTypedSlice(elements)
}

// Java method: Parser.parseImportSpecifier()
func (this_ *Parser) priv_Parser_parseImportSpecifier_7222() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var importedName Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
	_ = importedName
	var destinationName Token = (Token{})
	_ = destinationName
	if this_.priv_Parser_peekPredefinedString_7325("as") {
		this_.priv_Parser_eatPredefinedString_7326("as")
		destinationName = this_.priv_Parser_eatId_7422()
	} else if importedName.IsKeyword() {
		this_.priv_Parser_reportExpectedError_7426((Token{}), "as")
	}
	return ToIface_ParseTree(allocPooled_ImportSpecifierTree(&this_.slab_ImportSpecifierTree, this_.priv_Parser_getTreeLocation_7429(start), importedName, destinationName))
}

// Java method: Parser.peekExportDeclaration()
func (this_ *Parser) priv_Parser_peekExportDeclaration_7223() bool {
	return this_.priv_Parser_peek_7433(StaticGet_TokenType_EXPORT())
}

// Java method: Parser.parseExportDeclaration()
func (this_ *Parser) priv_Parser_parseExportDeclaration_7224() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var isDefault bool = false
	_ = isDefault
	var isExportAll bool = false
	_ = isExportAll
	var isExportSpecifier bool = false
	_ = isExportSpecifier
	var needsSemiColon bool = true
	_ = needsSemiColon
	this_.priv_Parser_eat_7425(StaticGet_TokenType_EXPORT())
	var export ParseTree = nil
	_ = export
	var exportSpecifierList []ParseTree
	_ = exportSpecifierList
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(65) /* TokenType.STAR */ :
		isExportAll = true
		this_.priv_Parser_nextToken_7430()
	case int32(2) /* TokenType.IDENTIFIER */ :
		export = this_.priv_Parser_parseAsyncFunctionDeclaration_7257()
	case int32(14) /* TokenType.FUNCTION */ :
		export = this_.priv_Parser_parseFunctionDeclaration_7255()
		needsSemiColon = false
	case int32(29) /* TokenType.CLASS */ :
		export = this_.priv_Parser_parseClassDeclaration_7228()
		needsSemiColon = false
	case int32(8) /* TokenType.DEFAULT */ :
		isDefault = true
		this_.priv_Parser_nextToken_7430()
		export = this_.priv_Parser_parseExpression_7337()
		needsSemiColon = false
	case int32(45) /* TokenType.OPEN_CURLY */ :
		isExportSpecifier = true
		exportSpecifierList = j2gort.ImmListToSliceNullable[ParseTree](this_.priv_Parser_parseExportSpecifierSet_7225())
	case int32(25) /* TokenType.VAR */, int32(38) /* TokenType.LET */, int32(30) /* TokenType.CONST */ :
		fallthrough
	default:
		export = ToIface_ParseTree(this_.priv_Parser_parseVariableDeclarationList_7275())
	}
	var moduleSpecifier Token = (Token{})
	_ = moduleSpecifier
	if isExportAll || (isExportSpecifier && this_.priv_Parser_peekPredefinedString_7325("from")) {
		this_.priv_Parser_eatPredefinedString_7326("from")
		moduleSpecifier = this_.priv_Parser_eat_7425(StaticGet_TokenType_STRING())
	} else if isExportSpecifier {
		{

			for _, tree := range exportSpecifierList {

				_ = tree
				var importedName Token = tree.GetParseTreeBase().AsExportSpecifier().importedName
				_ = importedName
				if importedName.IsKeyword() {
					this_.priv_Parser_reportError_7439(importedName, "cannot use keyword '%s' here.", []any{j2gort.BoxNullableVal(importedName)})
				}
			}
		}
	}
	if needsSemiColon || this_.priv_Parser_peekImplicitSemiColon_7411() {
		this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	}
	return ToIface_ParseTree(allocPooled_ExportDeclarationTree(&this_.slab_ExportDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start), isDefault, isExportAll, export, j2gort.ImmListFromTypedSlice(exportSpecifierList), moduleSpecifier))
}

// Java method: Parser.parseExportSpecifierSet()
func (this_ *Parser) priv_Parser_parseExportSpecifierSet_7225() j2gort.Ext_com_google_common_collect_ImmutableList {
	var elements []ParseTree = []ParseTree{}
	_ = elements
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	for this_.priv_Parser_peekIdOrKeyword_7419() {
		elements = append(elements, j2gort.CheckNotNullElem(this_.priv_Parser_parseExportSpecifier_7226()))
		if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_CURLY())) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
		}
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return j2gort.ImmListFromTypedSlice(elements)
}

// Java method: Parser.parseExportSpecifier()
func (this_ *Parser) priv_Parser_parseExportSpecifier_7226() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var importedName Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
	_ = importedName
	var destinationName Token = (Token{})
	_ = destinationName
	if this_.priv_Parser_peekPredefinedString_7325("as") {
		this_.priv_Parser_eatPredefinedString_7326("as")
		destinationName = this_.priv_Parser_eatIdOrKeywordAsId_7424()
	}
	return ToIface_ParseTree(allocPooled_ExportSpecifierTree(&this_.slab_ExportSpecifierTree, this_.priv_Parser_getTreeLocation_7429(start), importedName, destinationName))
}

// Java method: Parser.peekClassDeclaration()
func (this_ *Parser) priv_Parser_peekClassDeclaration_7227() bool {
	return this_.priv_Parser_peek_7433(StaticGet_TokenType_CLASS())
}

// Java method: Parser.parseClassDeclaration()
func (this_ *Parser) priv_Parser_parseClassDeclaration_7228() ParseTree {
	return this_.priv_Parser_parseClass_7230(false)
}

// Java method: Parser.parseClassExpression()
func (this_ *Parser) priv_Parser_parseClassExpression_7229() ParseTree {
	return this_.priv_Parser_parseClass_7230(true)
}

// Java method: Parser.parseClass(boolean isExpression)
func (this_ *Parser) priv_Parser_parseClass_7230(isExpression bool) ParseTree {
	_ = isExpression
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLASS())
	var name Token = (Token{})
	_ = name
	if !(isExpression) || this_.priv_Parser_peekId_7417() {
		name = this_.priv_Parser_eatId_7422()
	}
	var superClass ParseTree = nil
	_ = superClass
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_EXTENDS()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_EXTENDS())
		superClass = this_.priv_Parser_parseLeftHandSideExpression_7386()
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	var elements j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_Parser_parseClassElements_7231()
	_ = elements
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return ToIface_ParseTree(allocPooled_ClassDeclarationTree(&this_.slab_ClassDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start), name, superClass, elements))
}

// Java method: Parser.parseClassElements()
func (this_ *Parser) priv_Parser_parseClassElements_7231() j2gort.Ext_com_google_common_collect_ImmutableList {
	var result *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = result
	for {
		var token Token = this_.priv_Parser_peekToken_7437()
		_ = token
		if token.type_ == StaticGet_TokenType_SEMI_COLON() {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_SEMI_COLON())
			continue
		} else {
			if this_.priv_Parser_isClassElementStart_7232(token) {
			} else {
				return result.BuildImmList_()
			}
		}
		result.AddVoid_(this_.priv_Parser_parseClassElement_7234())
	}
	panic("unreachable")
}

// Java method: Parser.isClassElementStart(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *Parser) priv_Parser_isClassElementStart_7232(token Token) bool {
	_ = token
	switch (token.type_).Ordinal() {
	case int32(2) /* TokenType.IDENTIFIER */, int32(104) /* TokenType.NUMBER */, int32(106) /* TokenType.BIGINT */, int32(65) /* TokenType.STAR */, int32(43) /* TokenType.STATIC */, int32(105) /* TokenType.STRING */, int32(49) /* TokenType.OPEN_SQUARE */ :
		return true
	default:
		if Keywords_IsKeyword_1_TokenType(token.type_) {
			return true
		}
	}
	return false
}

// Java method: Parser.createObjectLiteralElementInfo()
func (this_ *Parser) priv_Parser_createObjectLiteralElementInfo_7233() *ClassOrObjectElementInfo {
	return ClassOrObjectElementInfo_CreateObjectLiteralElementInfo(this_.priv_Parser_getTreeStartLocation_7427())
}

// Java method: Parser.parseClassElement()
func (this_ *Parser) priv_Parser_parseClassElement_7234() ParseTree {
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_SEMI_COLON()) {
		return ToIface_ParseTree(this_.priv_Parser_parseEmptyStatement_7280())
	}
	if this_.priv_Parser_peekClassStaticInitializerBlock_7237() {
		return this_.priv_Parser_parseClassStaticInitializerBlock_7246()
	}
	var elementInfo *ClassOrObjectElementInfo = ClassOrObjectElementInfo_CreateClassMemberInfo(this_.priv_Parser_getTreeStartLocation_7427(), this_.priv_Parser_eatStaticIfNotElementName_7235())
	_ = elementInfo
	if this_.priv_Parser_peekGetAccessor_7324() {
		return this_.priv_Parser_parseGetAccessor_7329(elementInfo)
	} else if this_.priv_Parser_peekSetAccessor_7330() {
		return this_.priv_Parser_parseSetAccessor_7332(elementInfo)
	} else if this_.priv_Parser_peekAsyncMethod_7236() {
		return this_.priv_Parser_parseAsyncMethod_7245(elementInfo)
	} else {
		return this_.priv_Parser_parseClassMemberDeclaration_7243(elementInfo)
	}
}

// Java method: Parser.eatStaticIfNotElementName()
func (this_ *Parser) priv_Parser_eatStaticIfNotElementName_7235() bool {
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_STATIC()) && this_.priv_Parser_isClassElementStart_7232(this_.priv_Parser_peekToken_7438(int32(1))) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_STATIC())
		return true
	}
	return false
}

// Java method: Parser.peekAsyncMethod()
func (this_ *Parser) priv_Parser_peekAsyncMethod_7236() bool {
	return ((this_.priv_Parser_peekPredefinedString_7325("async") && !(this_.priv_Parser_peekImplicitSemiColon_7412(int32(1)))) && (this_.priv_Parser_peekPropertyNameOrComputedProp_7319(int32(1)) || (this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_STAR()) && this_.priv_Parser_peekPropertyNameOrComputedProp_7319(int32(2)))))
}

// Java method: Parser.peekClassStaticInitializerBlock()
func (this_ *Parser) priv_Parser_peekClassStaticInitializerBlock_7237() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_STATIC()) && this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_OPEN_CURLY()))
}

// Java method: Parser.parseClassElementName(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseClassElementName_7238(elementInfo *ClassOrObjectElementInfo) {
	_ = elementInfo
	if this_.priv_Parser_peekPropertyName_7320(int32(0)) {
		if this_.priv_Parser_peekIdOrKeyword_7419() {
			elementInfo.SetName(this_.priv_Parser_eatIdOrKeywordAsId_7424())
			if elementInfo.GetName().IsKeyword() {
				this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES())
			}
		} else {
			elementInfo.SetNameExpr(ToIface_ParseTree(this_.priv_Parser_parseLiteralExpression_7312()))
		}
	} else {
		elementInfo.SetNameExpr(this_.priv_Parser_parseComputedPropertyName_7323())
	}
}

// Java method: Parser.parseFieldDefinition(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseFieldDefinition_7239(elementInfo *ClassOrObjectElementInfo) ParseTree {
	_ = elementInfo
	var initializer ParseTree = nil
	_ = initializer
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
		initializer = this_.priv_Parser_parseInitializer_7279(StaticGet_Parser_Expression_NORMAL())
	}
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	if elementInfo.HasName() {
		return ToIface_ParseTree(allocPooled_FieldDeclarationTree(&this_.slab_FieldDeclarationTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), elementInfo.GetName(), elementInfo.isStatic, initializer))
	} else {
		if !(elementInfo.HasNameExpr()) {
			j2gort.CheckState(false)
		}
		return ToIface_ParseTree(allocPooled_ComputedPropertyFieldTree(&this_.slab_ComputedPropertyFieldTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), elementInfo.GetNameExpr(), elementInfo.isStatic, initializer))
	}
}

// Java method: Parser.parseMethodDefinition(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo, boolean isGenerator)
func (this_ *Parser) priv_Parser_parseMethodDefinition_7240(elementInfo *ClassOrObjectElementInfo, isGenerator bool) ParseTree {
	_ = elementInfo
	_ = isGenerator
	var kind FunctionDeclarationTree_Kind = FunctionDeclarationTree_Kind(0)
	_ = kind
	if elementInfo.HasNameExpr() {
		kind = StaticGet_FunctionDeclarationTree_Kind_EXPRESSION()
	} else {
		kind = StaticGet_FunctionDeclarationTree_Kind_MEMBER()
	}
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, kind).SetIsClassMember(elementInfo.isClassMember).SetStatic(elementInfo.isStatic)
	_ = builder
	if elementInfo.HasName() {
		builder.SetName(elementInfo.GetName())
	}
	this_.priv_Parser_parseFunctionTail_7247(builder, func() Parser_FunctionFlavor {
		if isGenerator {
			return StaticGet_Parser_FunctionFlavor_GENERATOR()
		}
		return StaticGet_Parser_FunctionFlavor_NORMAL()
	}())
	var function ParseTree = ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start)))
	_ = function
	if kind == StaticGet_FunctionDeclarationTree_Kind_MEMBER() {
		return function
	} else {
		if !(elementInfo.HasNameExpr()) {
			j2gort.CheckState(false)
		}
		return ToIface_ParseTree(allocPooled_ComputedPropertyMethodTree(&this_.slab_ComputedPropertyMethodTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), elementInfo.GetNameExpr(), function))
	}
}

// Java method: Parser.parseObjectLiteralMethodDeclaration()
func (this_ *Parser) priv_Parser_parseObjectLiteralMethodDeclaration_7241() ParseTree {
	return this_.priv_Parser_parseMethodDeclaration_7242(this_.priv_Parser_createObjectLiteralElementInfo_7233())
}

// Java method: Parser.parseMethodDeclaration(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseMethodDeclaration_7242(elementInfo *ClassOrObjectElementInfo) ParseTree {
	_ = elementInfo
	var isGenerator bool = !(this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_STAR())).IsJavaNil_()
	_ = isGenerator
	this_.priv_Parser_parseClassElementName_7238(elementInfo)
	return this_.priv_Parser_parseMethodDefinition_7240(elementInfo, isGenerator)
}

// Java method: Parser.parseClassMemberDeclaration(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseClassMemberDeclaration_7243(elementInfo *ClassOrObjectElementInfo) ParseTree {
	_ = elementInfo
	var isGenerator bool = !(this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_STAR())).IsJavaNil_()
	_ = isGenerator
	this_.priv_Parser_parseClassElementName_7238(elementInfo)
	if this_.priv_Parser_peekType_7436(int32(0)) == StaticGet_TokenType_OPEN_PAREN() {
		return this_.priv_Parser_parseMethodDefinition_7240(elementInfo, isGenerator)
	} else {
		return this_.priv_Parser_parseFieldDefinition_7239(elementInfo)
	}
}

// Java method: Parser.parseObjectLiteralAsyncMethod()
func (this_ *Parser) priv_Parser_parseObjectLiteralAsyncMethod_7244() ParseTree {
	return this_.priv_Parser_parseAsyncMethod_7245(this_.priv_Parser_createObjectLiteralElementInfo_7233())
}

// Java method: Parser.parseAsyncMethod(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseAsyncMethod_7245(elementInfo *ClassOrObjectElementInfo) ParseTree {
	_ = elementInfo
	this_.priv_Parser_eatPredefinedString_7326("async")
	var generator bool = this_.priv_Parser_peek_7433(StaticGet_TokenType_STAR())
	_ = generator
	if generator {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_STAR())
	}
	if this_.priv_Parser_peekPropertyName_7320(int32(0)) {
		if this_.priv_Parser_peekIdOrKeyword_7419() {
			var name Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
			_ = name
			var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_MEMBER()).SetAsync(true).SetGenerator(generator).SetIsClassMember(elementInfo.isClassMember).SetStatic(elementInfo.isStatic).SetName(name)
			_ = builder
			this_.priv_Parser_parseFunctionTail_7247(builder, func() Parser_FunctionFlavor {
				if generator {
					return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR()
				}
				return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS()
			}())
			return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(name.GetStart())))
		} else {
			var nameExpr ParseTree = ToIface_ParseTree(this_.priv_Parser_parseLiteralExpression_7312())
			_ = nameExpr
			var builder_2 *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_EXPRESSION()).SetAsync(true).SetGenerator(generator).SetIsClassMember(elementInfo.isClassMember).SetStatic(elementInfo.isStatic)
			_ = builder_2
			this_.priv_Parser_parseFunctionTail_7247(builder_2, func() Parser_FunctionFlavor {
				if generator {
					return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR()
				}
				return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS()
			}())
			var function ParseTree = ToIface_ParseTree(builder_2.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(nameExpr.GetParseTreeBase().GetStart())))
			_ = function
			return ToIface_ParseTree(allocPooled_ComputedPropertyMethodTree(&this_.slab_ComputedPropertyMethodTree, this_.priv_Parser_getTreeLocation_7429(nameExpr.GetParseTreeBase().GetStart()), nameExpr, function))
		}
	} else {
		var nameExpr_2 ParseTree = this_.priv_Parser_parseComputedPropertyName_7323()
		_ = nameExpr_2
		var builder_3 *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_EXPRESSION()).SetAsync(true).SetGenerator(generator).SetIsClassMember(elementInfo.isClassMember).SetStatic(elementInfo.isStatic)
		_ = builder_3
		this_.priv_Parser_parseFunctionTail_7247(builder_3, func() Parser_FunctionFlavor {
			if generator {
				return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR()
			}
			return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS()
		}())
		var function_2 ParseTree = ToIface_ParseTree(builder_3.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(nameExpr_2.GetParseTreeBase().GetStart())))
		_ = function_2
		return ToIface_ParseTree(allocPooled_ComputedPropertyMethodTree(&this_.slab_ComputedPropertyMethodTree, this_.priv_Parser_getTreeLocation_7429(nameExpr_2.GetParseTreeBase().GetStart()), nameExpr_2, function_2))
	}
}

// Java method: Parser.parseClassStaticInitializerBlock()
func (this_ *Parser) priv_Parser_parseClassStaticInitializerBlock_7246() ParseTree {
	this_.priv_Parser_eat_7425(StaticGet_TokenType_STATIC())
	return ToIface_ParseTree(this_.priv_Parser_parseBlock_7272())
}

// Java method: Parser.parseFunctionTail(com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Builder builder, com.google.javascript.jscomp.parsing.parser.Parser.FunctionFlavor functionFlavor)
func (this_ *Parser) priv_Parser_parseFunctionTail_7247(builder *FunctionDeclarationTree_Builder, functionFlavor Parser_FunctionFlavor) {
	_ = builder
	_ = functionFlavor
	var prevInGenerator bool = this_.inGenerator
	_ = prevInGenerator
	this_.inGenerator = functionFlavor.FieldGet_isGenerator()
	builder.SetGenerator(functionFlavor.FieldGet_isGenerator()).SetFormalParameterList(this_.priv_Parser_parseFormalParameterList_7263()).SetFunctionBody(ToIface_ParseTree(this_.priv_Parser_parseFunctionBody_7265()))
	this_.inGenerator = prevInGenerator
}

// Java method: Parser.parseSourceElement()
func (this_ *Parser) priv_Parser_parseSourceElement_7248() ParseTree {
	if this_.priv_Parser_peekAsyncFunctionStart_7250() {
		return this_.priv_Parser_parseAsyncFunctionDeclaration_7257()
	}
	if this_.priv_Parser_peekFunction_7252() {
		return this_.priv_Parser_parseFunctionDeclaration_7255()
	}
	if this_.priv_Parser_peekClassDeclaration_7227() {
		return this_.priv_Parser_parseClassDeclaration_7228()
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_LET()) {
		return ToIface_ParseTree(this_.priv_Parser_parseVariableStatement_7274())
	}
	return this_.priv_Parser_parseStatementStandard_7269()
}

// Java method: Parser.peekSourceElement()
func (this_ *Parser) priv_Parser_peekSourceElement_7249() bool {
	return ((this_.priv_Parser_peekFunction_7252() || this_.priv_Parser_peekStatementStandard_7271()) || this_.priv_Parser_peekDeclaration_7253())
}

// Java method: Parser.peekAsyncFunctionStart()
func (this_ *Parser) priv_Parser_peekAsyncFunctionStart_7250() bool {
	return ((this_.priv_Parser_peekPredefinedString_7325("async") && !(this_.priv_Parser_peekImplicitSemiColon_7412(int32(1)))) && this_.priv_Parser_peekFunction_7254(int32(1)))
}

// Java method: Parser.eatAsyncFunctionStart()
func (this_ *Parser) priv_Parser_eatAsyncFunctionStart_7251() {
	this_.priv_Parser_eatPredefinedString_7326("async")
	this_.priv_Parser_eat_7425(StaticGet_TokenType_FUNCTION())
}

// Java method: Parser.peekFunction()
func (this_ *Parser) priv_Parser_peekFunction_7252() bool {
	return this_.priv_Parser_peekFunction_7254(int32(0))
}

// Java method: Parser.peekDeclaration()
func (this_ *Parser) priv_Parser_peekDeclaration_7253() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_LET()) || this_.priv_Parser_peekClassDeclaration_7227())
}

// Java method: Parser.peekFunction(int index)
func (this_ *Parser) priv_Parser_peekFunction_7254(index int32) bool {
	_ = index
	return this_.priv_Parser_peek_7434(index, StaticGet_TokenType_FUNCTION())
}

// Java method: Parser.parseFunctionDeclaration()
func (this_ *Parser) priv_Parser_parseFunctionDeclaration_7255() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_Keywords_FUNCTION().type_)
	var isGenerator bool = !(this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_STAR())).IsJavaNil_()
	_ = isGenerator
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_DECLARATION()).SetName(this_.priv_Parser_eatId_7422())
	_ = builder
	this_.priv_Parser_parseFunctionTail_7247(builder, func() Parser_FunctionFlavor {
		if isGenerator {
			return StaticGet_Parser_FunctionFlavor_GENERATOR()
		}
		return StaticGet_Parser_FunctionFlavor_NORMAL()
	}())
	return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseFunctionExpression()
func (this_ *Parser) priv_Parser_parseFunctionExpression_7256() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_Keywords_FUNCTION().type_)
	var isGenerator bool = !(this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_STAR())).IsJavaNil_()
	_ = isGenerator
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_EXPRESSION()).SetName(this_.priv_Parser_eatIdOpt_7421())
	_ = builder
	this_.priv_Parser_parseFunctionTail_7247(builder, func() Parser_FunctionFlavor {
		if isGenerator {
			return StaticGet_Parser_FunctionFlavor_GENERATOR()
		}
		return StaticGet_Parser_FunctionFlavor_NORMAL()
	}())
	return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseAsyncFunctionDeclaration()
func (this_ *Parser) priv_Parser_parseAsyncFunctionDeclaration_7257() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eatAsyncFunctionStart_7251()
	var generator bool = this_.priv_Parser_peek_7433(StaticGet_TokenType_STAR())
	_ = generator
	if generator {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_STAR())
	}
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_DECLARATION()).SetName(this_.priv_Parser_eatId_7422()).SetAsync(true)
	_ = builder
	this_.priv_Parser_parseFunctionTail_7247(builder, func() Parser_FunctionFlavor {
		if generator {
			return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR()
		}
		return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS()
	}())
	return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseAsyncFunctionExpression()
func (this_ *Parser) priv_Parser_parseAsyncFunctionExpression_7258() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eatAsyncFunctionStart_7251()
	var generator bool = this_.priv_Parser_peek_7433(StaticGet_TokenType_STAR())
	_ = generator
	if generator {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_STAR())
	}
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_EXPRESSION()).SetName(this_.priv_Parser_eatIdOpt_7421()).SetAsync(true)
	_ = builder
	this_.priv_Parser_parseFunctionTail_7247(builder, func() Parser_FunctionFlavor {
		if generator {
			return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR()
		}
		return StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS()
	}())
	return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.peekParameter()
func (this_ *Parser) priv_Parser_peekParameter_7259() bool {
	if this_.priv_Parser_peekId_7417() || this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
		return true
	}
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_SQUARE()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_CURLY()))
}

// Java method: Parser.parseParameter()
func (this_ *Parser) priv_Parser_parseParameter_7260() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var parameter ParseTree = nil
	_ = parameter
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
		parameter = this_.priv_Parser_parseIterRest_7402(StaticGet_Parser_PatternKind_INITIALIZER())
	} else if this_.priv_Parser_peekId_7417() {
		parameter = ToIface_ParseTree(this_.priv_Parser_parseIdentifierExpression_7311())
	} else if this_.priv_Parser_peekPatternStart_7399() {
		parameter = this_.priv_Parser_parsePattern_7400(StaticGet_Parser_PatternKind_INITIALIZER())
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "parseParameterCalled() without confirming a parameter exists.")))
	}
	if !(parameter.GetParseTreeBase().IsRestParameter()) && this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_EQUAL())
		var defaultValue ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
		_ = defaultValue
		parameter = ToIface_ParseTree(allocPooled_DefaultParameterTree(&this_.slab_DefaultParameterTree, this_.priv_Parser_getTreeLocation_7429(start), parameter, defaultValue))
	}
	return parameter
}

// Java method: Parser.pushTreeScratch(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func (this_ *Parser) priv_Parser_pushTreeScratch_7261(tree ParseTree) {
	_ = tree
	var top int32 = this_.treeScratchTop
	_ = top
	if top == int32(len(this_.treeScratch)) {
		this_.treeScratch = func() []ParseTree {
			arrSrc_7 := this_.treeScratch
			arrLen_8 := (top * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_7)
			if arrLen_8 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_9 := make([]ParseTree, arrLen_8)
			copy(arrOut_9, arrSrc_7)
			return arrOut_9
		}()
	}
	{
		a_10 := this_.treeScratch
		i_11 := top
		a_10[i_11] = tree
	}
	this_.treeScratchTop = (top + int32(1))
}

// Java method: Parser.pushCommaScratch(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos)
func (this_ *Parser) priv_Parser_pushCommaScratch_7262(pos SourcePosition) {
	_ = pos
	var top int32 = this_.commaScratchTop
	_ = top
	if top == int32(len(this_.commaScratch)) {
		this_.commaScratch = func() []SourcePosition {
			arrSrc_12 := this_.commaScratch
			arrLen_13 := (top * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_12)
			if arrLen_13 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_14 := make([]SourcePosition, arrLen_13)
			copy(arrOut_14, arrSrc_12)
			return arrOut_14
		}()
	}
	{
		a_15 := this_.commaScratch
		i_16 := top
		a_15[i_16] = pos
	}
	this_.commaScratchTop = (top + int32(1))
}

// Java method: Parser.parseFormalParameterList()
func (this_ *Parser) priv_Parser_parseFormalParameterList_7263() *FormalParameterListTree {
	var listStart SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = listStart
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	if !(this_.priv_Parser_peekParameter_7259()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_1(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart))
	}
	var first ParseTree = this_.priv_Parser_parseParameter_7260()
	_ = first
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_5(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, false, int32(0), (SourcePosition{}))
	}
	var trailingComma bool = false
	_ = trailingComma
	var firstComma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
	_ = firstComma
	var firstCommaPos SourcePosition
	_ = firstCommaPos
	if !(firstComma).IsJavaNil_() {
		firstCommaPos = firstComma.GetStart()
	} else {
		firstCommaPos = this_.priv_Parser_getTreeEndLocation_7428()
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		if !(this_.config.atLeast8) {
			this_.priv_Parser_reportError_7439(firstComma, "Invalid trailing comma in formal parameter list", []any{})
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_5(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, true, int32(1), firstCommaPos)
	}
	if !(this_.priv_Parser_peekParameter_7259()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_5(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, false, int32(1), firstCommaPos)
	}
	var second ParseTree = this_.priv_Parser_parseParameter_7260()
	_ = second
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_7(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, second, false, int32(1), firstCommaPos, (SourcePosition{}))
	}
	var secondComma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
	_ = secondComma
	var secondCommaPos SourcePosition
	_ = secondCommaPos
	if !(secondComma).IsJavaNil_() {
		secondCommaPos = secondComma.GetStart()
	} else {
		secondCommaPos = this_.priv_Parser_getTreeEndLocation_7428()
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		if !(this_.config.atLeast8) {
			this_.priv_Parser_reportError_7439(secondComma, "Invalid trailing comma in formal parameter list", []any{})
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_7(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, second, true, int32(2), firstCommaPos, secondCommaPos)
	}
	if !(this_.priv_Parser_peekParameter_7259()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_7(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, second, false, int32(2), firstCommaPos, secondCommaPos)
	}
	var third ParseTree = this_.priv_Parser_parseParameter_7260()
	_ = third
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_9(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, second, third, false, int32(2), firstCommaPos, secondCommaPos, (SourcePosition{}))
	}
	var thirdComma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
	_ = thirdComma
	var thirdCommaPos SourcePosition
	_ = thirdCommaPos
	if !(thirdComma).IsJavaNil_() {
		thirdCommaPos = thirdComma.GetStart()
	} else {
		thirdCommaPos = this_.priv_Parser_getTreeEndLocation_7428()
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		if !(this_.config.atLeast8) {
			this_.priv_Parser_reportError_7439(thirdComma, "Invalid trailing comma in formal parameter list", []any{})
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_9(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, second, third, true, int32(3), firstCommaPos, secondCommaPos, thirdCommaPos)
	}
	if !(this_.priv_Parser_peekParameter_7259()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_FormalParameterListTree_9(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), first, second, third, false, int32(3), firstCommaPos, secondCommaPos, thirdCommaPos)
	}
	var savedTreeTop int32 = this_.treeScratchTop
	_ = savedTreeTop
	var savedCommaTop int32 = this_.commaScratchTop
	_ = savedCommaTop
	defer func() {
		{
			var i int32 = savedTreeTop
			_ = i
			var end int32 = this_.treeScratchTop
			_ = end
			for ; i < end; i = i + 1 {
				{
					a_17 := this_.treeScratch
					i_18 := i
					a_17[i_18] = nil
				}
			}
		}
		this_.treeScratchTop = savedTreeTop
		this_.commaScratchTop = savedCommaTop
	}()
	{
		var paramCount int32 = int32(3)
		_ = paramCount
		var commaCount int32 = int32(3)
		_ = commaCount
		for this_.priv_Parser_peekParameter_7259() {
			this_.priv_Parser_pushTreeScratch_7261(this_.priv_Parser_parseParameter_7260())
			paramCount = paramCount + 1
			if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN())) {
				var comma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
				_ = comma
				this_.priv_Parser_pushCommaScratch_7262(func() SourcePosition {
					if !(comma).IsJavaNil_() {
						return comma.GetStart()
					}
					return this_.priv_Parser_getTreeEndLocation_7428()
				}())
				commaCount = commaCount + 1
				if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
					if !(this_.config.atLeast8) {
						this_.priv_Parser_reportError_7439(comma, "Invalid trailing comma in formal parameter list", []any{})
					}
					trailingComma = true
				}
			}
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		var restParams []ParseTree
		{
			arrSrc_19 := this_.treeScratch
			arrFrom_20 := savedTreeTop
			arrTo_21 := this_.treeScratchTop
			j2gort.CheckNotNullSlice(arrSrc_19)
			if arrFrom_20 > arrTo_21 {
				panic(javalang.NewIllegalArgumentException("fromIndex > toIndex"))
			}
			if arrFrom_20 < 0 || int(arrFrom_20) > len(arrSrc_19) {
				panic(javalang.NewArrayIndexOutOfBoundsException(""))
			}
			arrLen_22 := arrTo_21 - arrFrom_20
			arrOut_23 := make([]ParseTree, arrLen_22)
			copy(arrOut_23, arrSrc_19[arrFrom_20:])
			restParams = arrOut_23
		}

		_ = restParams
		var restCommas []SourcePosition
		_ = restCommas
		if commaCount > int32(3) {
			{

				arrSrc_24 := this_.commaScratch
				arrFrom_25 := savedCommaTop
				arrTo_26 := this_.commaScratchTop
				j2gort.CheckNotNullSlice(arrSrc_24)
				if arrFrom_25 > arrTo_26 {
					panic(javalang.NewIllegalArgumentException("fromIndex > toIndex"))
				}
				if arrFrom_25 < 0 || int(arrFrom_25) > len(arrSrc_24) {
					panic(javalang.NewArrayIndexOutOfBoundsException(""))
				}
				arrLen_27 := arrTo_26 - arrFrom_25
				arrOut_28 := make([]SourcePosition, arrLen_27)
				copy(arrOut_28, arrSrc_24[arrFrom_25:])
				restCommas = arrOut_28
			}

		} else {
			restCommas = nil
		}
		return allocPooled_FormalParameterListTree_12(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(listStart), paramCount, first, second, third, restParams, trailingComma, commaCount, firstCommaPos, secondCommaPos, thirdCommaPos, restCommas)
	}
	panic("unreachable")
}

// Java method: Parser.parseSetterParameterList()
func (this_ *Parser) priv_Parser_parseSetterParameterList_7264() *FormalParameterListTree {
	var parameterList *FormalParameterListTree = this_.priv_Parser_parseFormalParameterList_7263()
	_ = parameterList
	if parameterList.parameterCount != int32(1) {
		this_.priv_Parser_reportError_7440(ToIface_ParseTree(parameterList), "Setter must have exactly 1 parameter, found %d", []any{j2gort.ValueOfInt(parameterList.parameterCount)})
	}
	if parameterList.parameterCount >= int32(1) {
		var parameter ParseTree = parameterList.param0
		_ = parameter
		if parameter.GetParseTreeBase().IsRestParameter() {
			this_.priv_Parser_reportError_7440(parameter, "Setter must not have a rest parameter", []any{})
		}
	}
	return parameterList
}

// Java method: Parser.parseFunctionBody()
func (this_ *Parser) priv_Parser_parseFunctionBody_7265() *BlockTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	var result j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_Parser_parseSourceElementList_7266()
	_ = result
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return allocPooled_BlockTree(&this_.slab_BlockTree, this_.priv_Parser_getTreeLocation_7429(start), result)
}

// Java method: Parser.parseSourceElementList()
func (this_ *Parser) priv_Parser_parseSourceElementList_7266() j2gort.Ext_com_google_common_collect_ImmutableList {
	if !(this_.priv_Parser_peekSourceElement_7249()) {
		return StaticGet_Parser_EMPTY_PARSE_TREES()
	}
	var first ParseTree = this_.priv_Parser_parseSourceElement_7248()
	_ = first
	if !(this_.priv_Parser_peekSourceElement_7249()) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(first)
	}
	var second ParseTree = this_.priv_Parser_parseSourceElement_7248()
	_ = second
	if !(this_.priv_Parser_peekSourceElement_7249()) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(first, second)
	}
	var result []ParseTree = []ParseTree{}
	_ = result
	result = append(result, j2gort.CheckNotNullElem(first))
	result = append(result, j2gort.CheckNotNullElem(second))
	for this_.priv_Parser_peekSourceElement_7249() {
		result = append(result, j2gort.CheckNotNullElem(this_.priv_Parser_parseSourceElement_7248()))
	}
	return j2gort.ImmListFromTypedSlice(result)
}

// Java method: Parser.parseIterSpread()
func (this_ *Parser) priv_Parser_parseIterSpread_7267() *IterSpreadTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_ELLIPSIS())
	var operand ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
	_ = operand
	return allocPooled_IterSpreadTree(&this_.slab_IterSpreadTree, this_.priv_Parser_getTreeLocation_7429(start), operand)
}

// Java method: Parser.parseStatement()
func (this_ *Parser) priv_Parser_parseStatement_7268() ParseTree {
	return this_.priv_Parser_parseSourceElement_7248()
}

// Java method: Parser.parseStatementStandard()
func (this_ *Parser) priv_Parser_parseStatementStandard_7269() ParseTree {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(45) /* TokenType.OPEN_CURLY */ :
		return ToIface_ParseTree(this_.priv_Parser_parseBlock_7272())
	case int32(30) /* TokenType.CONST */, int32(25) /* TokenType.VAR */ :
		return ToIface_ParseTree(this_.priv_Parser_parseVariableStatement_7274())
	case int32(52) /* TokenType.SEMI_COLON */ :
		return ToIface_ParseTree(this_.priv_Parser_parseEmptyStatement_7280())
	case int32(15) /* TokenType.IF */ :
		return ToIface_ParseTree(this_.priv_Parser_parseIfStatement_7282())
	case int32(10) /* TokenType.DO */ :
		return this_.priv_Parser_parseDoWhileStatement_7283()
	case int32(27) /* TokenType.WHILE */ :
		return this_.priv_Parser_parseWhileStatement_7284()
	case int32(13) /* TokenType.FOR */ :
		return this_.priv_Parser_parseForStatement_7285()
	case int32(6) /* TokenType.CONTINUE */ :
		return this_.priv_Parser_parseContinueStatement_7293()
	case int32(3) /* TokenType.BREAK */ :
		return this_.priv_Parser_parseBreakStatement_7294()
	case int32(19) /* TokenType.RETURN */ :
		return this_.priv_Parser_parseReturnStatement_7295()
	case int32(28) /* TokenType.WITH */ :
		return this_.priv_Parser_parseWithStatement_7296()
	case int32(20) /* TokenType.SWITCH */ :
		return this_.priv_Parser_parseSwitchStatement_7297()
	case int32(22) /* TokenType.THROW */ :
		return this_.priv_Parser_parseThrowStatement_7302()
	case int32(23) /* TokenType.TRY */ :
		return this_.priv_Parser_parseTryStatement_7303()
	case int32(7) /* TokenType.DEBUGGER */ :
		return this_.priv_Parser_parseDebuggerStatement_7306()
	default:
		if this_.priv_Parser_peekLabelledStatement_7301() {
			return this_.priv_Parser_parseLabelledStatement_7300()
		}
		return ToIface_ParseTree(this_.priv_Parser_parseExpressionStatement_7281())
	}
	panic("unreachable switch expression")
}

// Java method: Parser.peekStatement()
func (this_ *Parser) priv_Parser_peekStatement_7270() bool {
	return this_.priv_Parser_peekSourceElement_7249()
}

// Java method: Parser.peekStatementStandard()
func (this_ *Parser) priv_Parser_peekStatementStandard_7271() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(45) /* TokenType.OPEN_CURLY */, int32(25) /* TokenType.VAR */, int32(30) /* TokenType.CONST */, int32(52) /* TokenType.SEMI_COLON */, int32(15) /* TokenType.IF */, int32(10) /* TokenType.DO */, int32(27) /* TokenType.WHILE */, int32(13) /* TokenType.FOR */, int32(6) /* TokenType.CONTINUE */, int32(3) /* TokenType.BREAK */, int32(19) /* TokenType.RETURN */, int32(28) /* TokenType.WITH */, int32(20) /* TokenType.SWITCH */, int32(22) /* TokenType.THROW */, int32(23) /* TokenType.TRY */, int32(7) /* TokenType.DEBUGGER */, int32(44) /* TokenType.YIELD */, int32(2) /* TokenType.IDENTIFIER */, int32(113) /* TokenType.TYPE */, int32(114) /* TokenType.DECLARE */, int32(115) /* TokenType.MODULE */, int32(116) /* TokenType.NAMESPACE */, int32(21) /* TokenType.THIS */, int32(29) /* TokenType.CLASS */, int32(35) /* TokenType.SUPER */, int32(104) /* TokenType.NUMBER */, int32(106) /* TokenType.BIGINT */, int32(105) /* TokenType.STRING */, int32(112) /* TokenType.NO_SUBSTITUTION_TEMPLATE */, int32(109) /* TokenType.TEMPLATE_HEAD */, int32(101) /* TokenType.NULL */, int32(102) /* TokenType.TRUE */, int32(96) /* TokenType.SLASH */, int32(97) /* TokenType.SLASH_EQUAL */, int32(103) /* TokenType.FALSE */, int32(49) /* TokenType.OPEN_SQUARE */, int32(47) /* TokenType.OPEN_PAREN */, int32(18) /* TokenType.NEW */, int32(9) /* TokenType.DELETE */, int32(26) /* TokenType.VOID */, int32(24) /* TokenType.TYPEOF */, int32(68) /* TokenType.PLUS_PLUS */, int32(69) /* TokenType.MINUS_MINUS */, int32(63) /* TokenType.PLUS */, int32(64) /* TokenType.MINUS */, int32(77) /* TokenType.TILDE */, int32(76) /* TokenType.BANG */, int32(34) /* TokenType.IMPORT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseBlock()
func (this_ *Parser) priv_Parser_parseBlock_7272() *BlockTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	var result j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_Parser_parseSourceElementList_7266()
	_ = result
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return allocPooled_BlockTree(&this_.slab_BlockTree, this_.priv_Parser_getTreeLocation_7429(start), result)
}

// Java method: Parser.parseStatementList()
func (this_ *Parser) priv_Parser_parseStatementList_7273() j2gort.Ext_com_google_common_collect_ImmutableList {
	if !(this_.priv_Parser_peekStatement_7270()) {
		return StaticGet_Parser_EMPTY_PARSE_TREES()
	}
	var first ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = first
	if !(this_.priv_Parser_peekStatement_7270()) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(first)
	}
	var result []ParseTree = []ParseTree{}
	_ = result
	result = append(result, j2gort.CheckNotNullElem(first))
	for this_.priv_Parser_peekStatement_7270() {
		result = append(result, j2gort.CheckNotNullElem(this_.priv_Parser_parseStatement_7268()))
	}
	return j2gort.ImmListFromTypedSlice(result)
}

// Java method: Parser.parseVariableStatement()
func (this_ *Parser) priv_Parser_parseVariableStatement_7274() *VariableStatementTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var declarations *VariableDeclarationListTree = this_.priv_Parser_parseVariableDeclarationList_7275()
	_ = declarations
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return allocPooled_VariableStatementTree(&this_.slab_VariableStatementTree, this_.priv_Parser_getTreeLocation_7429(start), declarations)
}

// Java method: Parser.parseVariableDeclarationList()
func (this_ *Parser) priv_Parser_parseVariableDeclarationList_7275() *VariableDeclarationListTree {
	return this_.priv_Parser_parseVariableDeclarationList_7277(StaticGet_Parser_Expression_NORMAL())
}

// Java method: Parser.parseVariableDeclarationListNoIn()
func (this_ *Parser) priv_Parser_parseVariableDeclarationListNoIn_7276() *VariableDeclarationListTree {
	return this_.priv_Parser_parseVariableDeclarationList_7277(StaticGet_Parser_Expression_NO_IN())
}

// Java method: Parser.parseVariableDeclarationList(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseVariableDeclarationList_7277(expressionIn Parser_Expression) *VariableDeclarationListTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var token TokenType = this_.priv_Parser_peekType_7435()
	_ = token
	switch (token).Ordinal() {
	case int32(30) /* TokenType.CONST */, int32(38) /* TokenType.LET */, int32(25) /* TokenType.VAR */ :
		this_.priv_Parser_eat_7425(token)
	default:
		this_.priv_Parser_reportError_7439(this_.priv_Parser_peekToken_7437(), "expected declaration", []any{})
		return nil
	}
	var declarations []*VariableDeclarationTree = []*VariableDeclarationTree{}
	_ = declarations
	declarations = append(declarations, j2gort.CheckNotNullElem(this_.priv_Parser_parseVariableDeclaration_7278(token, expressionIn)))
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
		declarations = append(declarations, j2gort.CheckNotNullElem(this_.priv_Parser_parseVariableDeclaration_7278(token, expressionIn)))
	}
	return allocPooled_VariableDeclarationListTree(&this_.slab_VariableDeclarationListTree, this_.priv_Parser_getTreeLocation_7429(start), token, j2gort.ImmListFromTypedSlice(declarations))
}

// Java method: Parser.parseVariableDeclaration(com.google.javascript.jscomp.parsing.parser.TokenType binding, com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseVariableDeclaration_7278(binding TokenType, expressionIn Parser_Expression) *VariableDeclarationTree {
	_ = binding
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var lvalue ParseTree = nil
	_ = lvalue
	if this_.priv_Parser_peekPatternStart_7399() {
		lvalue = this_.priv_Parser_parsePattern_7400(StaticGet_Parser_PatternKind_INITIALIZER())
	} else {
		lvalue = ToIface_ParseTree(this_.priv_Parser_parseIdentifierExpression_7311())
	}
	var initializer ParseTree = nil
	_ = initializer
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
		initializer = this_.priv_Parser_parseInitializer_7279(expressionIn)
	} else if expressionIn != StaticGet_Parser_Expression_NO_IN() {
		this_.priv_Parser_maybeReportNoInitializer_7289(binding, lvalue)
	}
	return allocPooled_VariableDeclarationTree(&this_.slab_VariableDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start), lvalue, initializer)
}

// Java method: Parser.parseInitializer(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseInitializer_7279(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	this_.priv_Parser_eat_7425(StaticGet_TokenType_EQUAL())
	return this_.priv_Parser_parseAssignment_7342(expressionIn)
}

// Java method: Parser.parseEmptyStatement()
func (this_ *Parser) priv_Parser_parseEmptyStatement_7280() *EmptyStatementTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_SEMI_COLON())
	return allocPooled_EmptyStatementTree(&this_.slab_EmptyStatementTree, this_.priv_Parser_getTreeLocation_7429(start))
}

// Java method: Parser.parseExpressionStatement()
func (this_ *Parser) priv_Parser_parseExpressionStatement_7281() *ExpressionStatementTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var expression ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = expression
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return allocPooled_ExpressionStatementTree(&this_.slab_ExpressionStatementTree, this_.priv_Parser_getTreeLocation_7429(start), expression)
}

// Java method: Parser.parseIfStatement()
func (this_ *Parser) priv_Parser_parseIfStatement_7282() *IfStatementTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_IF())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	var condition ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = condition
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var ifClause ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = ifClause
	var elseClause ParseTree = nil
	_ = elseClause
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELSE()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_ELSE())
		elseClause = this_.priv_Parser_parseStatement_7268()
	}
	return allocPooled_IfStatementTree(&this_.slab_IfStatementTree, this_.priv_Parser_getTreeLocation_7429(start), condition, ifClause, elseClause)
}

// Java method: Parser.parseDoWhileStatement()
func (this_ *Parser) priv_Parser_parseDoWhileStatement_7283() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_DO())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	this_.priv_Parser_eat_7425(StaticGet_TokenType_WHILE())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	var condition ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = condition
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_SEMI_COLON()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_SEMI_COLON())
	}
	return ToIface_ParseTree(allocPooled_DoWhileStatementTree(&this_.slab_DoWhileStatementTree, this_.priv_Parser_getTreeLocation_7429(start), body, condition))
}

// Java method: Parser.parseWhileStatement()
func (this_ *Parser) priv_Parser_parseWhileStatement_7284() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_WHILE())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	var condition ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = condition
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	return ToIface_ParseTree(allocPooled_WhileStatementTree(&this_.slab_WhileStatementTree, this_.priv_Parser_getTreeLocation_7429(start), condition, body))
}

// Java method: Parser.parseForStatement()
func (this_ *Parser) priv_Parser_parseForStatement_7285() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_FOR())
	var awaited bool = this_.priv_Parser_peekPredefinedString_7325("await")
	_ = awaited
	if awaited {
		this_.priv_Parser_eatPredefinedString_7326("await")
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	if this_.priv_Parser_peekVariableDeclarationList_7290() {
		var variables *VariableDeclarationListTree = this_.priv_Parser_parseVariableDeclarationListNoIn_7276()
		_ = variables
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_IN()) {
			if awaited {
				this_.priv_Parser_reportError_7441("for-await-of is the only allowed asynchronous iteration", []any{})
			}
			if int32(len(variables.declarations)) > int32(1) {
				this_.priv_Parser_reportError_7441("for-in statement may not have more than one variable declaration", []any{})
			}
			var declaration *VariableDeclarationTree = variables.declarations[int32(0)]
			_ = declaration
			if declaration.initializer != nil && !declaration.initializer.IsJavaNil_() {
				if this_.config.atLeast6 {
					this_.priv_Parser_reportError_7441("for-in statement may not have initializer", []any{})
				} else {
					this_.errorReporter.GetErrorReporterBase().ReportWarning_3(declaration.GetParseTreeBase().location.start, "for-in statement should not have initializer", []any{})
				}
			}
			return this_.priv_Parser_parseForInStatement_7292(start, ToIface_ParseTree(variables))
		} else if this_.priv_Parser_peekPredefinedString_7325("of") {
			if int32(len(variables.declarations)) > int32(1) {
				if awaited {
					this_.priv_Parser_reportError_7441("for-await-of statement may not have more than one variable declaration", []any{})
				} else {
					this_.priv_Parser_reportError_7441("for-of statement may not have more than one variable declaration", []any{})
				}
			}
			var declaration_2 *VariableDeclarationTree = variables.declarations[int32(0)]
			_ = declaration_2
			if declaration_2.initializer != nil && !declaration_2.initializer.IsJavaNil_() {
				if awaited {
					this_.priv_Parser_reportError_7441("for-await-of statement may not have initializer", []any{})
				} else {
					this_.priv_Parser_reportError_7441("for-of statement may not have initializer", []any{})
				}
			}
			if awaited {
				return this_.priv_Parser_parseForAwaitOfStatement_7287(start, ToIface_ParseTree(variables))
			} else {
				return this_.priv_Parser_parseForOfStatement_7286(start, ToIface_ParseTree(variables))
			}
		} else {
			this_.priv_Parser_checkVanillaForInitializers_7288(variables)
			return this_.priv_Parser_parseForStatement_7291(start, ToIface_ParseTree(variables))
		}
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_SEMI_COLON()) {
		return this_.priv_Parser_parseForStatement_7291(start, nil)
	}
	var initializer ParseTree = this_.priv_Parser_parseExpressionNoIn_7336()
	_ = initializer
	if (this_.priv_Parser_peek_7433(StaticGet_TokenType_IN()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL())) || this_.priv_Parser_peekPredefinedString_7325("of") {
		initializer = this_.priv_Parser_transformLeftHandSideExpression_7351(initializer)
		if !(initializer.GetParseTreeBase().IsValidAssignmentTarget()) {
			this_.priv_Parser_reportError_7441("invalid assignment target", []any{})
		}
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_IN()) || this_.priv_Parser_peekPredefinedString_7325("of") {
		if (initializer.GetParseTreeBase().type_ != StaticGet_ParseTreeType_BINARY_OPERATOR()) && (initializer.GetParseTreeBase().type_ != StaticGet_ParseTreeType_COMMA_EXPRESSION()) {
			if this_.priv_Parser_peek_7433(StaticGet_TokenType_IN()) {
				return this_.priv_Parser_parseForInStatement_7292(start, initializer)
			} else {
				if awaited {
					return this_.priv_Parser_parseForAwaitOfStatement_7287(start, initializer)
				} else {
					return this_.priv_Parser_parseForOfStatement_7286(start, initializer)
				}
			}
		}
	}
	return this_.priv_Parser_parseForStatement_7291(start, initializer)
}

// Java method: Parser.parseForOfStatement(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer)
func (this_ *Parser) priv_Parser_parseForOfStatement_7286(start SourcePosition, initializer ParseTree) ParseTree {
	_ = start
	_ = initializer
	this_.priv_Parser_eatPredefinedString_7326("of")
	var collection ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = collection
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	return ToIface_ParseTree(allocPooled_ForOfStatementTree(&this_.slab_ForOfStatementTree, this_.priv_Parser_getTreeLocation_7429(start), initializer, collection, body))
}

// Java method: Parser.parseForAwaitOfStatement(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer)
func (this_ *Parser) priv_Parser_parseForAwaitOfStatement_7287(start SourcePosition, initializer ParseTree) ParseTree {
	_ = start
	_ = initializer
	this_.priv_Parser_eatPredefinedString_7326("of")
	var collection ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = collection
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	return ToIface_ParseTree(allocPooled_ForAwaitOfStatementTree(&this_.slab_ForAwaitOfStatementTree, this_.priv_Parser_getTreeLocation_7429(start), initializer, collection, body))
}

// Java method: Parser.checkVanillaForInitializers(com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree variables)
func (this_ *Parser) priv_Parser_checkVanillaForInitializers_7288(variables *VariableDeclarationListTree) {
	_ = variables
	{

		for _, declaration := range variables.declarations {

			_ = declaration
			if declaration.initializer == nil || declaration.initializer.IsJavaNil_() {
				this_.priv_Parser_maybeReportNoInitializer_7289(variables.declarationType, declaration.lvalue)
			}
		}
	}
}

// Java method: Parser.maybeReportNoInitializer(com.google.javascript.jscomp.parsing.parser.TokenType token, com.google.javascript.jscomp.parsing.parser.trees.ParseTree lvalue)
func (this_ *Parser) priv_Parser_maybeReportNoInitializer_7289(token TokenType, lvalue ParseTree) {
	_ = token
	_ = lvalue
	if token == StaticGet_TokenType_CONST() {
		this_.priv_Parser_reportError_7441("const variables must have an initializer", []any{})
	} else if lvalue.GetParseTreeBase().IsPattern() {
		this_.priv_Parser_reportError_7441("destructuring must have an initializer", []any{})
	}
}

// Java method: Parser.peekVariableDeclarationList()
func (this_ *Parser) priv_Parser_peekVariableDeclarationList_7290() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(25) /* TokenType.VAR */, int32(30) /* TokenType.CONST */, int32(38) /* TokenType.LET */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseForStatement(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.trees.@org.jspecify.annotations.Nullable ParseTree initializer)
func (this_ *Parser) priv_Parser_parseForStatement_7291(start SourcePosition, initializer ParseTree) ParseTree {
	_ = start
	_ = initializer
	if initializer == nil || initializer.IsJavaNil_() {
		initializer = ToIface_ParseTree(allocPooled_NullTree(&this_.slab_NullTree, NewSourceRange(this_.priv_Parser_getTreeEndLocation_7428(), this_.priv_Parser_getTreeStartLocation_7427())))
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_SEMI_COLON())
	var condition ParseTree = nil
	_ = condition
	if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_SEMI_COLON())) {
		condition = this_.priv_Parser_parseExpression_7337()
	} else {
		condition = ToIface_ParseTree(allocPooled_NullTree(&this_.slab_NullTree, NewSourceRange(this_.priv_Parser_getTreeEndLocation_7428(), this_.priv_Parser_getTreeStartLocation_7427())))
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_SEMI_COLON())
	var increment ParseTree = nil
	_ = increment
	if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN())) {
		increment = this_.priv_Parser_parseExpression_7337()
	} else {
		increment = ToIface_ParseTree(allocPooled_NullTree(&this_.slab_NullTree, NewSourceRange(this_.priv_Parser_getTreeEndLocation_7428(), this_.priv_Parser_getTreeStartLocation_7427())))
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	return ToIface_ParseTree(allocPooled_ForStatementTree(&this_.slab_ForStatementTree, this_.priv_Parser_getTreeLocation_7429(start), initializer, condition, increment, body))
}

// Java method: Parser.parseForInStatement(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer)
func (this_ *Parser) priv_Parser_parseForInStatement_7292(start SourcePosition, initializer ParseTree) ParseTree {
	_ = start
	_ = initializer
	this_.priv_Parser_eat_7425(StaticGet_TokenType_IN())
	var collection ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = collection
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	return ToIface_ParseTree(allocPooled_ForInStatementTree(&this_.slab_ForInStatementTree, this_.priv_Parser_getTreeLocation_7429(start), initializer, collection, body))
}

// Java method: Parser.parseContinueStatement()
func (this_ *Parser) priv_Parser_parseContinueStatement_7293() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CONTINUE())
	var name Token = (Token{})
	_ = name
	if !(this_.priv_Parser_peekImplicitSemiColon_7411()) {
		name = this_.priv_Parser_eatIdOpt_7421()
	}
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return ToIface_ParseTree(allocPooled_ContinueStatementTree(&this_.slab_ContinueStatementTree, this_.priv_Parser_getTreeLocation_7429(start), name))
}

// Java method: Parser.parseBreakStatement()
func (this_ *Parser) priv_Parser_parseBreakStatement_7294() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_BREAK())
	var name Token = (Token{})
	_ = name
	if !(this_.priv_Parser_peekImplicitSemiColon_7411()) {
		name = this_.priv_Parser_eatIdOpt_7421()
	}
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return ToIface_ParseTree(allocPooled_BreakStatementTree(&this_.slab_BreakStatementTree, this_.priv_Parser_getTreeLocation_7429(start), name))
}

// Java method: Parser.parseReturnStatement()
func (this_ *Parser) priv_Parser_parseReturnStatement_7295() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_RETURN())
	var expression ParseTree = nil
	_ = expression
	if !(this_.priv_Parser_peekImplicitSemiColon_7411()) {
		expression = this_.priv_Parser_parseExpression_7337()
	}
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return ToIface_ParseTree(allocPooled_ReturnStatementTree(&this_.slab_ReturnStatementTree, this_.priv_Parser_getTreeLocation_7429(start), expression))
}

// Java method: Parser.parseWithStatement()
func (this_ *Parser) priv_Parser_parseWithStatement_7296() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_WITH())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	var expression ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = expression
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	var body ParseTree = this_.priv_Parser_parseStatement_7268()
	_ = body
	return ToIface_ParseTree(allocPooled_WithStatementTree(&this_.slab_WithStatementTree, this_.priv_Parser_getTreeLocation_7429(start), expression, body))
}

// Java method: Parser.parseSwitchStatement()
func (this_ *Parser) priv_Parser_parseSwitchStatement_7297() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_SWITCH())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	var expression ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = expression
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	var caseClauses j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_Parser_parseCaseClauses_7298()
	_ = caseClauses
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return ToIface_ParseTree(allocPooled_SwitchStatementTree(&this_.slab_SwitchStatementTree, this_.priv_Parser_getTreeLocation_7429(start), expression, caseClauses))
}

// Java method: Parser.parseCaseClauses()
func (this_ *Parser) priv_Parser_parseCaseClauses_7298() j2gort.Ext_com_google_common_collect_ImmutableList {
	var foundDefaultClause bool = false
	_ = foundDefaultClause
	var result *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = result
	for {
		var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
		_ = start
		switch (this_.priv_Parser_peekType_7435()).Ordinal() {
		case int32(4) /* TokenType.CASE */ :
			this_.priv_Parser_eat_7425(StaticGet_TokenType_CASE())
			var expression ParseTree = this_.priv_Parser_parseExpression_7337()
			_ = expression
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
			var statements j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_Parser_parseCaseStatementsOpt_7299()
			_ = statements
			result.AddVoid_(j2gort.PtrToAny(allocPooled_CaseClauseTree(&this_.slab_CaseClauseTree, this_.priv_Parser_getTreeLocation_7429(start), expression, statements)))
		case int32(8) /* TokenType.DEFAULT */ :
			if foundDefaultClause {
				this_.priv_Parser_reportError_7441("Switch statements may have at most one default clause", []any{})
			} else {
				foundDefaultClause = true
			}
			this_.priv_Parser_eat_7425(StaticGet_TokenType_DEFAULT())
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
			result.AddVoid_(j2gort.PtrToAny(allocPooled_DefaultClauseTree(&this_.slab_DefaultClauseTree, this_.priv_Parser_getTreeLocation_7429(start), this_.priv_Parser_parseCaseStatementsOpt_7299())))
		default:
			return result.BuildImmList_()
		}
	}
	panic("unreachable")
}

// Java method: Parser.parseCaseStatementsOpt()
func (this_ *Parser) priv_Parser_parseCaseStatementsOpt_7299() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.priv_Parser_parseStatementList_7273()
}

// Java method: Parser.parseLabelledStatement()
func (this_ *Parser) priv_Parser_parseLabelledStatement_7300() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var name Token = this_.priv_Parser_eatId_7422()
	_ = name
	this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
	return ToIface_ParseTree(allocPooled_LabelledStatementTree(&this_.slab_LabelledStatementTree, this_.priv_Parser_getTreeLocation_7429(start), name, this_.priv_Parser_parseStatement_7268()))
}

// Java method: Parser.peekLabelledStatement()
func (this_ *Parser) priv_Parser_peekLabelledStatement_7301() bool {
	return (this_.priv_Parser_peekId_7417() && this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_COLON()))
}

// Java method: Parser.parseThrowStatement()
func (this_ *Parser) priv_Parser_parseThrowStatement_7302() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_THROW())
	var value ParseTree = nil
	_ = value
	if this_.priv_Parser_peekImplicitSemiColon_7411() {
		this_.priv_Parser_reportError_7441("semicolon/newline not allowed after 'throw'", []any{})
	} else {
		value = this_.priv_Parser_parseExpression_7337()
	}
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return ToIface_ParseTree(allocPooled_ThrowStatementTree(&this_.slab_ThrowStatementTree, this_.priv_Parser_getTreeLocation_7429(start), value))
}

// Java method: Parser.parseTryStatement()
func (this_ *Parser) priv_Parser_parseTryStatement_7303() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_TRY())
	var body ParseTree = ToIface_ParseTree(this_.priv_Parser_parseBlock_7272())
	_ = body
	var catchBlock ParseTree = nil
	_ = catchBlock
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CATCH()) {
		catchBlock = ToIface_ParseTree(this_.priv_Parser_parseCatch_7304())
	}
	var finallyBlock ParseTree = nil
	_ = finallyBlock
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_FINALLY()) {
		finallyBlock = ToIface_ParseTree(this_.priv_Parser_parseFinallyBlock_7305())
	}
	if (catchBlock == nil || catchBlock.IsJavaNil_()) && (finallyBlock == nil || finallyBlock.IsJavaNil_()) {
		this_.priv_Parser_reportError_7441("'catch' or 'finally' expected.", []any{})
	}
	return ToIface_ParseTree(allocPooled_TryStatementTree(&this_.slab_TryStatementTree, this_.priv_Parser_getTreeLocation_7429(start), body, catchBlock, finallyBlock))
}

// Java method: Parser.parseCatch()
func (this_ *Parser) priv_Parser_parseCatch_7304() *CatchTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CATCH())
	var exception ParseTree = ToIface_ParseTree(allocPooled_EmptyStatementTree(&this_.slab_EmptyStatementTree, NewSourceRange(this_.priv_Parser_getTreeEndLocation_7428(), this_.priv_Parser_getTreeStartLocation_7427())))
	_ = exception
	if this_.priv_Parser_peekToken_7437().type_ == StaticGet_TokenType_OPEN_PAREN() {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
		if this_.priv_Parser_peekPatternStart_7399() {
			exception = this_.priv_Parser_parsePattern_7400(StaticGet_Parser_PatternKind_INITIALIZER())
		} else {
			exception = ToIface_ParseTree(this_.priv_Parser_parseIdentifierExpression_7311())
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	} else {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING())
	}
	var catchBody *BlockTree = this_.priv_Parser_parseBlock_7272()
	_ = catchBody
	return allocPooled_CatchTree(&this_.slab_CatchTree, this_.priv_Parser_getTreeLocation_7429(start), exception, ToIface_ParseTree(catchBody))
}

// Java method: Parser.parseFinallyBlock()
func (this_ *Parser) priv_Parser_parseFinallyBlock_7305() *FinallyTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_FINALLY())
	var finallyBlock *BlockTree = this_.priv_Parser_parseBlock_7272()
	_ = finallyBlock
	return allocPooled_FinallyTree(&this_.slab_FinallyTree, this_.priv_Parser_getTreeLocation_7429(start), ToIface_ParseTree(finallyBlock))
}

// Java method: Parser.parseDebuggerStatement()
func (this_ *Parser) priv_Parser_parseDebuggerStatement_7306() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_DEBUGGER())
	this_.priv_Parser_eatPossiblyImplicitSemiColon_7410()
	return ToIface_ParseTree(allocPooled_DebuggerStatementTree(&this_.slab_DebuggerStatementTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parsePrimaryExpression()
func (this_ *Parser) priv_Parser_parsePrimaryExpression_7307() ParseTree {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(29) /* TokenType.CLASS */ :
		return this_.priv_Parser_parseClassExpression_7229()
	case int32(35) /* TokenType.SUPER */ :
		return ToIface_ParseTree(this_.priv_Parser_parseSuperExpression_7308())
	case int32(21) /* TokenType.THIS */ :
		return ToIface_ParseTree(this_.priv_Parser_parseThisExpression_7309())
	case int32(34) /* TokenType.IMPORT */ :
		return ToIface_ParseTree(this_.priv_Parser_parseDynamicImportExpression_7310())
	case int32(2) /* TokenType.IDENTIFIER */, int32(113) /* TokenType.TYPE */, int32(114) /* TokenType.DECLARE */, int32(115) /* TokenType.MODULE */, int32(116) /* TokenType.NAMESPACE */ :
		return ToIface_ParseTree(this_.priv_Parser_parseIdentifierExpression_7311())
	case int32(104) /* TokenType.NUMBER */, int32(105) /* TokenType.STRING */, int32(106) /* TokenType.BIGINT */, int32(102) /* TokenType.TRUE */, int32(103) /* TokenType.FALSE */, int32(101) /* TokenType.NULL */ :
		return ToIface_ParseTree(this_.priv_Parser_parseLiteralExpression_7312())
	case int32(112) /* TokenType.NO_SUBSTITUTION_TEMPLATE */, int32(109) /* TokenType.TEMPLATE_HEAD */ :
		return ToIface_ParseTree(this_.priv_Parser_parseTemplateLiteral_7313(nil))
	case int32(49) /* TokenType.OPEN_SQUARE */ :
		return this_.priv_Parser_parseArrayLiteral_7316()
	case int32(45) /* TokenType.OPEN_CURLY */ :
		return this_.priv_Parser_parseObjectLiteral_7317()
	case int32(47) /* TokenType.OPEN_PAREN */ :
		return this_.priv_Parser_parseCoverParenthesizedExpressionAndArrowParameterList_7334()
	case int32(96) /* TokenType.SLASH */, int32(97) /* TokenType.SLASH_EQUAL */ :
		return this_.priv_Parser_parseRegularExpressionLiteral_7315()
	default:
		return this_.priv_Parser_parseMissingPrimaryExpression_7335()
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseSuperExpression()
func (this_ *Parser) priv_Parser_parseSuperExpression_7308() *SuperExpressionTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_SUPER())
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_DOT()) {
		this_.priv_Parser_reportError_7441("Optional chaining is forbidden in super?.", []any{})
	}
	return allocPooled_SuperExpressionTree(&this_.slab_SuperExpressionTree, this_.priv_Parser_getTreeLocation_7429(start))
}

// Java method: Parser.parseThisExpression()
func (this_ *Parser) priv_Parser_parseThisExpression_7309() *ThisExpressionTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_THIS())
	return allocPooled_ThisExpressionTree(&this_.slab_ThisExpressionTree, this_.priv_Parser_getTreeLocation_7429(start))
}

// Java method: Parser.parseDynamicImportExpression()
func (this_ *Parser) priv_Parser_parseDynamicImportExpression_7310() *DynamicImportTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_IMPORT())
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_DOT()) {
		this_.priv_Parser_reportError_7441("Optional chaining is forbidden in import?.", []any{})
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	var argument ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
	_ = argument
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
	return allocPooled_DynamicImportTree(&this_.slab_DynamicImportTree, this_.priv_Parser_getTreeLocation_7429(start), argument)
}

// Java method: Parser.parseIdentifierExpression()
func (this_ *Parser) priv_Parser_parseIdentifierExpression_7311() *IdentifierExpressionTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var identifier Token = this_.priv_Parser_eatId_7422()
	_ = identifier
	return allocPooled_IdentifierExpressionTree(&this_.slab_IdentifierExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), identifier)
}

// Java method: Parser.parseLiteralExpression()
func (this_ *Parser) priv_Parser_parseLiteralExpression_7312() *LiteralExpressionTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var literal Token = this_.priv_Parser_nextLiteralToken_7314()
	_ = literal
	if (literal.type_ == StaticGet_TokenType_STRING()) && literal.HasUnescapedUnicodeLineOrParagraphSeparator() {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP())
	}
	if (literal.type_ == StaticGet_TokenType_NUMBER()) && (j2gort.JStrIndexOf(literal.String(), "_") >= 0) {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_NUMERIC_SEPARATOR())
	}
	if literal.type_ == StaticGet_TokenType_BIGINT() {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_BIGINT())
	}
	return allocPooled_LiteralExpressionTree(&this_.slab_LiteralExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), literal)
}

// Java method: Parser.parseTemplateLiteral(com.google.javascript.jscomp.parsing.parser.trees.@org.jspecify.annotations.Nullable ParseTree operand)
func (this_ *Parser) priv_Parser_parseTemplateLiteral_7313(operand ParseTree) *TemplateLiteralExpressionTree {
	_ = operand
	var start SourcePosition
	_ = start
	if operand == nil || operand.IsJavaNil_() {
		start = this_.priv_Parser_getTreeStartLocation_7427()
	} else {
		start = operand.GetParseTreeBase().location.start
	}
	var token Token = this_.priv_Parser_nextToken_7430()
	_ = token
	if !(token.IsTemplateLiteralToken()) {
		this_.priv_Parser_reportError_7439(token, "Unexpected template literal token %s.", []any{j2gort.BoxString(token.type_.String())})
	}
	var isTaggedTemplate bool = (operand != nil && !operand.IsJavaNil_())
	_ = isTaggedTemplate
	var templateToken Token = token
	_ = templateToken
	if !(isTaggedTemplate) {
		this_.priv_Parser_reportTemplateErrorIfPresent_7443(templateToken)
	}
	var elements *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = elements
	elements.AddVoid_(j2gort.PtrToAny(allocPooled_TemplateLiteralPortionTree(&this_.slab_TemplateLiteralPortionTree, templateToken.location, templateToken)))
	if templateToken.type_ == StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE() {
		return allocPooled_TemplateLiteralExpressionTree(&this_.slab_TemplateLiteralExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, elements.BuildImmList_())
	}
	var expression ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = expression
	elements.AddVoid_(j2gort.PtrToAny(allocPooled_TemplateSubstitutionTree(&this_.slab_TemplateSubstitutionTree, expression.GetParseTreeBase().location, expression)))
	for !(this_.errorReporter.GetErrorReporterBase().HadError()) {
		templateToken = this_.priv_Parser_nextTemplateLiteralToken_7432()
		if (templateToken.type_ == StaticGet_TokenType_ERROR()) || (templateToken.type_ == StaticGet_TokenType_END_OF_FILE()) {
			break
		}
		if !(isTaggedTemplate) {
			this_.priv_Parser_reportTemplateErrorIfPresent_7443(templateToken)
		}
		elements.AddVoid_(j2gort.PtrToAny(allocPooled_TemplateLiteralPortionTree(&this_.slab_TemplateLiteralPortionTree, templateToken.location, templateToken)))
		if templateToken.type_ == StaticGet_TokenType_TEMPLATE_TAIL() {
			break
		}
		expression = this_.priv_Parser_parseExpression_7337()
		elements.AddVoid_(j2gort.PtrToAny(allocPooled_TemplateSubstitutionTree(&this_.slab_TemplateSubstitutionTree, expression.GetParseTreeBase().location, expression)))
	}
	return allocPooled_TemplateLiteralExpressionTree(&this_.slab_TemplateLiteralExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, elements.BuildImmList_())
}

// Java method: Parser.nextLiteralToken()
func (this_ *Parser) priv_Parser_nextLiteralToken_7314() Token {
	return this_.priv_Parser_nextToken_7430()
}

// Java method: Parser.parseRegularExpressionLiteral()
func (this_ *Parser) priv_Parser_parseRegularExpressionLiteral_7315() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var literal Token = this_.priv_Parser_nextRegularExpressionLiteralToken_7431()
	_ = literal
	this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_REGEXP_SYNTAX())
	return ToIface_ParseTree(allocPooled_LiteralExpressionTree(&this_.slab_LiteralExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), literal))
}

// Java method: Parser.parseArrayLiteral()
func (this_ *Parser) priv_Parser_parseArrayLiteral_7316() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var elements []ParseTree = []ParseTree{}
	_ = elements
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
	var trailingCommaToken Token = (Token{})
	_ = trailingCommaToken
	for (this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS())) || this_.priv_Parser_peekAssignmentExpression_7341() {
		trailingCommaToken = (Token{})
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
			var commaStart SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
			_ = commaStart
			trailingCommaToken = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
			elements = append(elements, allocPooled_NullTree(&this_.slab_NullTree, NewSourceRange(commaStart, commaStart)))
		} else {
			if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
				this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS())
				elements = append(elements, j2gort.CheckNotNullElem(this_.priv_Parser_parseIterSpread_7267()))
			} else {
				elements = append(elements, j2gort.CheckNotNullElem(this_.priv_Parser_parseAssignmentExpression_7340()))
			}
			if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_SQUARE())) {
				trailingCommaToken = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
			}
		}
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())
	this_.MaybeReportTrailingComma(trailingCommaToken)
	return ToIface_ParseTree(allocPooled_ArrayLiteralExpressionTree(&this_.slab_ArrayLiteralExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), j2gort.ImmListFromTypedSlice(elements), !(trailingCommaToken).IsJavaNil_()))
}

// Java method: Parser.parseObjectLiteral()
func (this_ *Parser) priv_Parser_parseObjectLiteral_7317() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var result []ParseTree = []ParseTree{}
	_ = result
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	var commaToken Token = (Token{})
	_ = commaToken
	for (this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) || this_.priv_Parser_peekPropertyNameOrComputedProp_7319(int32(0))) || this_.priv_Parser_peek_7433(StaticGet_TokenType_STAR()) {
		result = append(result, j2gort.CheckNotNullElem(this_.priv_Parser_parseObjectLiteralPropertyAssignment_7321()))
		commaToken = this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_COMMA())
		if (commaToken).IsJavaNil_() {
			break
		}
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	this_.MaybeReportTrailingComma(commaToken)
	return ToIface_ParseTree(allocPooled_ObjectLiteralExpressionTree(&this_.slab_ObjectLiteralExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), j2gort.ImmListFromTypedSlice(result), !(commaToken).IsJavaNil_()))
}

// Java method: Parser.maybeReportTrailingComma(com.google.javascript.jscomp.parsing.parser.Token commaToken)
func (this_ *Parser) MaybeReportTrailingComma(commaToken Token) {
	_ = commaToken
	if !(commaToken).IsJavaNil_() {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_TRAILING_COMMA())
		if this_.config.warnTrailingCommas {
			this_.errorReporter.GetErrorReporterBase().ReportWarning_3(commaToken.location.start, "Trailing comma is not legal in an ECMA-262 object initializer", []any{})
		}
	}
}

// Java method: Parser.peekPropertyNameOrComputedProp(int tokenIndex)
func (this_ *Parser) priv_Parser_peekPropertyNameOrComputedProp_7319(tokenIndex int32) bool {
	_ = tokenIndex
	return (this_.priv_Parser_peekPropertyName_7320(tokenIndex) || (this_.priv_Parser_peekType_7436(tokenIndex) == StaticGet_TokenType_OPEN_SQUARE()))
}

// Java method: Parser.peekPropertyName(int tokenIndex)
func (this_ *Parser) priv_Parser_peekPropertyName_7320(tokenIndex int32) bool {
	_ = tokenIndex
	var type_ TokenType = this_.priv_Parser_peekType_7436(tokenIndex)
	_ = type_
	switch (type_).Ordinal() {
	case int32(2) /* TokenType.IDENTIFIER */, int32(105) /* TokenType.STRING */, int32(104) /* TokenType.NUMBER */, int32(106) /* TokenType.BIGINT */ :
		return true
	default:
		return Keywords_IsKeyword_1_TokenType(type_)
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseObjectLiteralPropertyAssignment()
func (this_ *Parser) priv_Parser_parseObjectLiteralPropertyAssignment_7321() ParseTree {
	var type_ TokenType = this_.priv_Parser_peekType_7435()
	_ = type_
	if type_ == StaticGet_TokenType_STAR() {
		return this_.priv_Parser_parseObjectLiteralPropertyAssignmentGenerator_7322()
	} else if type_ == StaticGet_TokenType_ELLIPSIS() {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD())
		var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
		_ = start
		this_.priv_Parser_eat_7425(StaticGet_TokenType_ELLIPSIS())
		var operand ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
		_ = operand
		return ToIface_ParseTree(allocPooled_ObjectSpreadTree(&this_.slab_ObjectSpreadTree, this_.priv_Parser_getTreeLocation_7429(start), operand))
	} else if ((((type_ == StaticGet_TokenType_STRING()) || (type_ == StaticGet_TokenType_NUMBER())) || (type_ == StaticGet_TokenType_BIGINT())) || (type_ == StaticGet_TokenType_IDENTIFIER())) || Keywords_IsKeyword_1_TokenType(type_) {
		if this_.priv_Parser_peekGetAccessor_7324() {
			return this_.priv_Parser_parseObjectLiteralGetAccessor_7328()
		} else if this_.priv_Parser_peekSetAccessor_7330() {
			return this_.priv_Parser_parseObjectLiteralSetAccessor_7331()
		} else if this_.priv_Parser_peekAsyncMethod_7236() {
			return this_.priv_Parser_parseObjectLiteralAsyncMethod_7244()
		} else if this_.priv_Parser_peekType_7436(int32(1)) == StaticGet_TokenType_OPEN_PAREN() {
			return this_.priv_Parser_parseObjectLiteralMethodDeclaration_7241()
		} else {
			return this_.priv_Parser_parseObjectLiteralPropertyNameAssignment_7333()
		}
	} else if type_ == StaticGet_TokenType_OPEN_SQUARE() {
		var start_2 SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
		_ = start_2
		var name ParseTree = this_.priv_Parser_parseComputedPropertyName_7323()
		_ = name
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_COLON()) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
			var value ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
			_ = value
			return ToIface_ParseTree(allocPooled_ComputedPropertyDefinitionTree(&this_.slab_ComputedPropertyDefinitionTree, this_.priv_Parser_getTreeLocation_7429(start_2), name, value))
		} else {
			var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_EXPRESSION())
			_ = builder
			this_.priv_Parser_parseFunctionTail_7247(builder, StaticGet_Parser_FunctionFlavor_NORMAL())
			var value_2 ParseTree = ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start_2)))
			_ = value_2
			return ToIface_ParseTree(allocPooled_ComputedPropertyMethodTree(&this_.slab_ComputedPropertyMethodTree, this_.priv_Parser_getTreeLocation_7429(start_2), name, value_2))
		}
	} else {
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "unreachable")))
	}
}

// Java method: Parser.parseObjectLiteralPropertyAssignmentGenerator()
func (this_ *Parser) priv_Parser_parseObjectLiteralPropertyAssignmentGenerator_7322() ParseTree {
	var type_ TokenType = this_.priv_Parser_peekType_7436(int32(1))
	_ = type_
	if (((type_ == StaticGet_TokenType_STRING()) || (type_ == StaticGet_TokenType_NUMBER())) || (type_ == StaticGet_TokenType_IDENTIFIER())) || Keywords_IsKeyword_1_TokenType(type_) {
		return this_.priv_Parser_parseObjectLiteralMethodDeclaration_7241()
	} else {
		var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
		_ = start
		this_.priv_Parser_eat_7425(StaticGet_TokenType_STAR())
		var name ParseTree = this_.priv_Parser_parseComputedPropertyName_7323()
		_ = name
		var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_EXPRESSION())
		_ = builder
		this_.priv_Parser_parseFunctionTail_7247(builder, StaticGet_Parser_FunctionFlavor_GENERATOR())
		var value ParseTree = ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start)))
		_ = value
		return ToIface_ParseTree(allocPooled_ComputedPropertyMethodTree(&this_.slab_ComputedPropertyMethodTree, this_.priv_Parser_getTreeLocation_7429(start), name, value))
	}
}

// Java method: Parser.parseComputedPropertyName()
func (this_ *Parser) priv_Parser_parseComputedPropertyName_7323() ParseTree {
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
	var assign ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
	_ = assign
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())
	return assign
}

// Java method: Parser.peekGetAccessor()
func (this_ *Parser) priv_Parser_peekGetAccessor_7324() bool {
	return (this_.priv_Parser_peekPredefinedString_7325("get") && this_.priv_Parser_peekPropertyNameOrComputedProp_7319(int32(1)))
}

// Java method: Parser.peekPredefinedString(java.lang.String string)
func (this_ *Parser) priv_Parser_peekPredefinedString_7325(string_ string) bool {
	_ = string_
	return this_.priv_Parser_peekPredefinedString_7327(int32(0), string_)
}

// Java method: Parser.eatPredefinedString(java.lang.String string)
func (this_ *Parser) priv_Parser_eatPredefinedString_7326(string_ string) Token {
	_ = string_
	var token Token = this_.priv_Parser_eatId_7422()
	_ = token
	if (token).IsJavaNil_() || !(token.ValueEquals(string_)) {
		this_.priv_Parser_reportExpectedError_7426(token, j2gort.BoxString(string_))
		return (Token{})
	}
	return token
}

// Java method: Parser.peekPredefinedString(int index, java.lang.String string)
func (this_ *Parser) priv_Parser_peekPredefinedString_7327(index int32, string_ string) bool {
	_ = index
	_ = string_
	return (this_.priv_Parser_peek_7434(index, StaticGet_TokenType_IDENTIFIER()) && this_.priv_Parser_peekToken_7438(index).ValueEquals(string_))
}

// Java method: Parser.parseObjectLiteralGetAccessor()
func (this_ *Parser) priv_Parser_parseObjectLiteralGetAccessor_7328() ParseTree {
	return this_.priv_Parser_parseGetAccessor_7329(this_.priv_Parser_createObjectLiteralElementInfo_7233())
}

// Java method: Parser.parseGetAccessor(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseGetAccessor_7329(elementInfo *ClassOrObjectElementInfo) ParseTree {
	_ = elementInfo
	this_.priv_Parser_eatPredefinedString_7326("get")
	if this_.priv_Parser_peekPropertyName_7320(int32(0)) {
		var propertyName Token = this_.priv_Parser_eatObjectLiteralPropertyName_7423()
		_ = propertyName
		this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		var body *BlockTree = this_.priv_Parser_parseFunctionBody_7265()
		_ = body
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_GETTER())
		return ToIface_ParseTree(allocPooled_GetAccessorTree(&this_.slab_GetAccessorTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), propertyName, elementInfo.isClassMember, elementInfo.isStatic, body))
	} else {
		var property ParseTree = this_.priv_Parser_parseComputedPropertyName_7323()
		_ = property
		this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		var body_2 *BlockTree = this_.priv_Parser_parseFunctionBody_7265()
		_ = body_2
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_GETTER())
		return ToIface_ParseTree(allocPooled_ComputedPropertyGetterTree(&this_.slab_ComputedPropertyGetterTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), property, elementInfo.isStatic, body_2))
	}
}

// Java method: Parser.peekSetAccessor()
func (this_ *Parser) priv_Parser_peekSetAccessor_7330() bool {
	return (this_.priv_Parser_peekPredefinedString_7325("set") && this_.priv_Parser_peekPropertyNameOrComputedProp_7319(int32(1)))
}

// Java method: Parser.parseObjectLiteralSetAccessor()
func (this_ *Parser) priv_Parser_parseObjectLiteralSetAccessor_7331() ParseTree {
	return this_.priv_Parser_parseSetAccessor_7332(this_.priv_Parser_createObjectLiteralElementInfo_7233())
}

// Java method: Parser.parseSetAccessor(com.google.javascript.jscomp.parsing.parser.ClassOrObjectElementInfo elementInfo)
func (this_ *Parser) priv_Parser_parseSetAccessor_7332(elementInfo *ClassOrObjectElementInfo) ParseTree {
	_ = elementInfo
	this_.priv_Parser_eatPredefinedString_7326("set")
	if this_.priv_Parser_peekPropertyName_7320(int32(0)) {
		var propertyName Token = this_.priv_Parser_eatObjectLiteralPropertyName_7423()
		_ = propertyName
		var parameter *FormalParameterListTree = this_.priv_Parser_parseSetterParameterList_7264()
		_ = parameter
		var body *BlockTree = this_.priv_Parser_parseFunctionBody_7265()
		_ = body
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_SETTER())
		return ToIface_ParseTree(allocPooled_SetAccessorTree(&this_.slab_SetAccessorTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), propertyName, elementInfo.isClassMember, elementInfo.isStatic, parameter, body))
	} else {
		var property ParseTree = this_.priv_Parser_parseComputedPropertyName_7323()
		_ = property
		var parameter_2 *FormalParameterListTree = this_.priv_Parser_parseSetterParameterList_7264()
		_ = parameter_2
		var body_2 *BlockTree = this_.priv_Parser_parseFunctionBody_7265()
		_ = body_2
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_SETTER())
		return ToIface_ParseTree(allocPooled_ComputedPropertySetterTree(&this_.slab_ComputedPropertySetterTree, this_.priv_Parser_getTreeLocation_7429(elementInfo.start), property, elementInfo.isStatic, parameter_2, body_2))
	}
}

// Java method: Parser.parseObjectLiteralPropertyNameAssignment()
func (this_ *Parser) priv_Parser_parseObjectLiteralPropertyNameAssignment_7333() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var name Token = this_.priv_Parser_eatObjectLiteralPropertyName_7423()
	_ = name
	var colon Token = this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_COLON())
	_ = colon
	if (colon).IsJavaNil_() {
		if name.type_ != StaticGet_TokenType_IDENTIFIER() {
			this_.priv_Parser_reportExpectedError_7426(this_.priv_Parser_peekToken_7437(), j2gort.BoxNullableVal(StaticGet_TokenType_COLON()))
		} else if name.IsKeyword() {
			this_.priv_Parser_reportError_7439(name, "Cannot use keyword in short object literal", []any{})
		} else if this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
			var idTree *IdentifierExpressionTree = allocPooled_IdentifierExpressionTree(&this_.slab_IdentifierExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), name)
			_ = idTree
			this_.priv_Parser_eat_7425(StaticGet_TokenType_EQUAL())
			var defaultValue ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
			_ = defaultValue
			return ToIface_ParseTree(allocPooled_DefaultParameterTree(&this_.slab_DefaultParameterTree, this_.priv_Parser_getTreeLocation_7429(start), ToIface_ParseTree(idTree), defaultValue))
		}
	}
	var value ParseTree
	_ = value
	if (colon).IsJavaNil_() {
		value = nil
	} else {
		value = this_.priv_Parser_parseAssignmentExpression_7340()
	}
	return ToIface_ParseTree(allocPooled_PropertyNameAssignmentTree(&this_.slab_PropertyNameAssignmentTree, this_.priv_Parser_getTreeLocation_7429(start), name, value))
}

// Java method: Parser.parseCoverParenthesizedExpressionAndArrowParameterList()
func (this_ *Parser) priv_Parser_parseCoverParenthesizedExpressionAndArrowParameterList_7334() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_ARROW()) {
			return ToIface_ParseTree(allocPooled_FormalParameterListTree_1(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(start)))
		} else {
			this_.priv_Parser_reportError_7441("invalid parenthesized expression", []any{})
			return ToIface_ParseTree(allocPooled_MissingPrimaryExpressionTree(&this_.slab_MissingPrimaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(start)))
		}
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
		var restParam ParseTree = this_.priv_Parser_parseParameter_7260()
		_ = restParam
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_ARROW()) {
			return ToIface_ParseTree(allocPooled_FormalParameterListTree_5(&this_.slab_FormalParameterListTree, this_.priv_Parser_getTreeLocation_7429(start), restParam, false, int32(0), (SourcePosition{})))
		} else {
			this_.priv_Parser_reportError_7441("invalid parenthesized expression", []any{})
			return ToIface_ParseTree(allocPooled_MissingPrimaryExpressionTree(&this_.slab_MissingPrimaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(start)))
		}
	}
	var result ParseTree = this_.priv_Parser_parseExpression_7337()
	_ = result
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
		if this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_CLOSE_PAREN()) {
			this_.priv_Parser_resetScanner_7353(start)
			return ToIface_ParseTree(this_.priv_Parser_parseFormalParameterList_7263())
		} else {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
			result = ToIface_ParseTree(allocPooled_CommaExpressionTree(&this_.slab_CommaExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(result, this_.priv_Parser_parseParameter_7260())))
		}
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
	return ToIface_ParseTree(allocPooled_ParenExpressionTree(&this_.slab_ParenExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), result))
}

// Java method: Parser.parseMissingPrimaryExpression()
func (this_ *Parser) priv_Parser_parseMissingPrimaryExpression_7335() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_nextToken_7430()
	this_.priv_Parser_reportError_7441("primary expression expected", []any{})
	return ToIface_ParseTree(allocPooled_MissingPrimaryExpressionTree(&this_.slab_MissingPrimaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseExpressionNoIn()
func (this_ *Parser) priv_Parser_parseExpressionNoIn_7336() ParseTree {
	return this_.priv_Parser_parse_7339(StaticGet_Parser_Expression_NO_IN())
}

// Java method: Parser.parseExpression()
func (this_ *Parser) priv_Parser_parseExpression_7337() ParseTree {
	return this_.priv_Parser_parse_7339(StaticGet_Parser_Expression_NORMAL())
}

// Java method: Parser.peekExpression()
func (this_ *Parser) priv_Parser_peekExpression_7338() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(76) /* TokenType.BANG */, int32(29) /* TokenType.CLASS */, int32(9) /* TokenType.DELETE */, int32(103) /* TokenType.FALSE */, int32(14) /* TokenType.FUNCTION */, int32(2) /* TokenType.IDENTIFIER */, int32(113) /* TokenType.TYPE */, int32(114) /* TokenType.DECLARE */, int32(115) /* TokenType.MODULE */, int32(116) /* TokenType.NAMESPACE */, int32(64) /* TokenType.MINUS */, int32(69) /* TokenType.MINUS_MINUS */, int32(18) /* TokenType.NEW */, int32(101) /* TokenType.NULL */, int32(104) /* TokenType.NUMBER */, int32(106) /* TokenType.BIGINT */, int32(45) /* TokenType.OPEN_CURLY */, int32(47) /* TokenType.OPEN_PAREN */, int32(49) /* TokenType.OPEN_SQUARE */, int32(63) /* TokenType.PLUS */, int32(68) /* TokenType.PLUS_PLUS */, int32(96) /* TokenType.SLASH */, int32(97) /* TokenType.SLASH_EQUAL */, int32(105) /* TokenType.STRING */, int32(112) /* TokenType.NO_SUBSTITUTION_TEMPLATE */, int32(109) /* TokenType.TEMPLATE_HEAD */, int32(35) /* TokenType.SUPER */, int32(21) /* TokenType.THIS */, int32(77) /* TokenType.TILDE */, int32(102) /* TokenType.TRUE */, int32(24) /* TokenType.TYPEOF */, int32(26) /* TokenType.VOID */, int32(44) /* TokenType.YIELD */ :
		return true
	case int32(34) /* TokenType.IMPORT */ :
		return (this_.priv_Parser_peekImportCall_7384() || this_.priv_Parser_peekImportDot_7385())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parse(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parse_7339(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var result ParseTree = this_.priv_Parser_parseAssignment_7342(expressionIn)
	_ = result
	if (this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) && !(this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_ELLIPSIS()))) && !(this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_CLOSE_PAREN())) {
		var exprs []ParseTree = []ParseTree{}
		_ = exprs
		exprs = append(exprs, j2gort.CheckNotNullElem(result))
		for (this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) && !(this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_ELLIPSIS()))) && !(this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_CLOSE_PAREN())) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
			exprs = append(exprs, j2gort.CheckNotNullElem(this_.priv_Parser_parseAssignment_7342(expressionIn)))
		}
		return ToIface_ParseTree(allocPooled_CommaExpressionTree(&this_.slab_CommaExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), j2gort.ImmListFromTypedSlice(exprs)))
	}
	return result
}

// Java method: Parser.parseAssignmentExpression()
func (this_ *Parser) priv_Parser_parseAssignmentExpression_7340() ParseTree {
	return this_.priv_Parser_parseAssignment_7342(StaticGet_Parser_Expression_NORMAL())
}

// Java method: Parser.peekAssignmentExpression()
func (this_ *Parser) priv_Parser_peekAssignmentExpression_7341() bool {
	return this_.priv_Parser_peekExpression_7338()
}

// Java method: Parser.parseAssignment(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseAssignment_7342(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_YIELD()) && this_.priv_Parser_inGeneratorContext_7357() {
		return this_.priv_Parser_parseYield_7358(expressionIn)
	}
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseConditional_7359(expressionIn)
	_ = left
	if this_.priv_Parser_isStartOfAsyncArrowFunction_7343(left) {
		this_.priv_Parser_resetScanner_7354(left)
		return this_.priv_Parser_parseAsyncArrowFunction_7348(expressionIn)
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ARROW()) {
		return this_.priv_Parser_completeAssignmentExpressionParseAtArrow_7344(left, expressionIn)
	}
	if this_.priv_Parser_peekAssignmentOperator_7356() {
		if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL())) {
			if !(left.GetParseTreeBase().IsValidNonVanillaAssignmentTarget()) {
				this_.priv_Parser_reportError_7441("invalid assignment target", []any{})
				return ToIface_ParseTree(allocPooled_MissingPrimaryExpressionTree(&this_.slab_MissingPrimaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(this_.priv_Parser_getTreeStartLocation_7427())))
			}
		}
		left = this_.priv_Parser_transformLeftHandSideExpression_7351(left)
		if !(left.GetParseTreeBase().IsValidAssignmentTarget()) {
			this_.priv_Parser_reportError_7441("invalid assignment target", []any{})
			return ToIface_ParseTree(allocPooled_MissingPrimaryExpressionTree(&this_.slab_MissingPrimaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(this_.priv_Parser_getTreeStartLocation_7427())))
		}
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseAssignment_7342(expressionIn)
		_ = right
		return ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.isStartOfAsyncArrowFunction(com.google.javascript.jscomp.parsing.parser.trees.ParseTree partialExpression)
func (this_ *Parser) priv_Parser_isStartOfAsyncArrowFunction_7343(partialExpression ParseTree) bool {
	_ = partialExpression
	base_partialExpression_ := partialExpression.GetParseTreeBase()
	if partialExpression.GetParseTreeBase().type_ == StaticGet_ParseTreeType_IDENTIFIER_EXPRESSION() {
		var identifierToken Token = base_partialExpression_.AsIdentifierExpression().identifierToken
		_ = identifierToken
		return ((((identifierToken.ValueEquals("async") && !(this_.priv_Parser_peekImplicitSemiColon_7412(int32(0)))) && this_.priv_Parser_peekId_7417()) && !(this_.priv_Parser_peekImplicitSemiColon_7412(int32(1)))) && this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_ARROW()))
	} else if partialExpression.GetParseTreeBase().type_ == StaticGet_ParseTreeType_CALL_EXPRESSION() {
		var callExpression *CallExpressionTree = base_partialExpression_.AsCallExpression()
		_ = callExpression
		var callee ParseTree = callExpression.operand
		_ = callee
		var arguments ParseTree = ToIface_ParseTree(callExpression.arguments)
		_ = arguments
		return (((((callee.GetParseTreeBase().type_ == StaticGet_ParseTreeType_IDENTIFIER_EXPRESSION()) && callee.GetParseTreeBase().AsIdentifierExpression().identifierToken.ValueEquals("async")) && (callee.GetParseTreeBase().location.end.line == arguments.GetParseTreeBase().location.start.line)) && !(this_.priv_Parser_peekImplicitSemiColon_7411())) && this_.priv_Parser_peek_7433(StaticGet_TokenType_ARROW()))
	} else {
		return false
	}
}

// Java method: Parser.completeAssignmentExpressionParseAtArrow(com.google.javascript.jscomp.parsing.parser.trees.ParseTree leftOfArrow, com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_completeAssignmentExpressionParseAtArrow_7344(leftOfArrow ParseTree, expressionIn Parser_Expression) ParseTree {
	_ = leftOfArrow
	_ = expressionIn
	if leftOfArrow.GetParseTreeBase().type_ == StaticGet_ParseTreeType_CALL_EXPRESSION() {
		return this_.priv_Parser_completeAssignmentExpressionParseAtArrow_7347(leftOfArrow.GetParseTreeBase().AsCallExpression())
	} else {
		return this_.priv_Parser_completeArrowFunctionParseAtArrow_7345(leftOfArrow, expressionIn)
	}
}

// Java method: Parser.completeArrowFunctionParseAtArrow(com.google.javascript.jscomp.parsing.parser.trees.ParseTree leftOfArrow, com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_completeArrowFunctionParseAtArrow_7345(leftOfArrow ParseTree, expressionIn Parser_Expression) ParseTree {
	_ = leftOfArrow
	_ = expressionIn
	var arrowFormalParameters *FormalParameterListTree = this_.priv_Parser_transformToArrowFormalParameters_7346(leftOfArrow)
	_ = arrowFormalParameters
	if this_.priv_Parser_peekImplicitSemiColon_7411() {
		this_.priv_Parser_reportError_7441("No newline allowed before '=>'", []any{})
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_ARROW())
	var arrowFunctionBody ParseTree = this_.priv_Parser_parseArrowFunctionBody_7349(expressionIn, StaticGet_Parser_FunctionFlavor_NORMAL())
	_ = arrowFunctionBody
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_ARROW()).SetFormalParameterList(arrowFormalParameters).SetFunctionBody(arrowFunctionBody)
	_ = builder
	return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(arrowFormalParameters.GetParseTreeBase().location.start)))
}

// Java method: Parser.transformToArrowFormalParameters(com.google.javascript.jscomp.parsing.parser.trees.ParseTree leftOfArrow)
func (this_ *Parser) priv_Parser_transformToArrowFormalParameters_7346(leftOfArrow ParseTree) *FormalParameterListTree {
	_ = leftOfArrow
	var arrowParameterList *FormalParameterListTree = nil
	_ = arrowParameterList
	switch (leftOfArrow.GetParseTreeBase().type_).Ordinal() {
	case int32(57) /* ParseTreeType.FORMAL_PARAMETER_LIST */ :
		arrowParameterList = leftOfArrow.GetParseTreeBase().AsFormalParameterList()
	case int32(26) /* ParseTreeType.IDENTIFIER_EXPRESSION */ :
		arrowParameterList = allocPooled_FormalParameterListTree_5(&this_.slab_FormalParameterListTree, leftOfArrow.GetParseTreeBase().location, leftOfArrow, false, int32(0), (SourcePosition{}))
	case int32(46) /* ParseTreeType.ARGUMENT_LIST */, int32(52) /* ParseTreeType.PAREN_EXPRESSION */ :
		this_.priv_Parser_resetScanner_7354(leftOfArrow)
		arrowParameterList = this_.priv_Parser_parseFormalParameterList_7263()
	default:
		this_.priv_Parser_reportError_7440(leftOfArrow, "invalid arrow function parameters", []any{})
		arrowParameterList = this_.priv_Parser_newEmptyFormalParameterList_7350(leftOfArrow.GetParseTreeBase().location)
	}
	return arrowParameterList
}

// Java method: Parser.completeAssignmentExpressionParseAtArrow(com.google.javascript.jscomp.parsing.parser.trees.CallExpressionTree callExpression)
func (this_ *Parser) priv_Parser_completeAssignmentExpressionParseAtArrow_7347(callExpression *CallExpressionTree) ParseTree {
	_ = callExpression
	var operand ParseTree = callExpression.operand
	_ = operand
	var arguments ParseTree = ToIface_ParseTree(callExpression.arguments)
	_ = arguments
	var result ParseTree = nil
	_ = result
	if operand.GetParseTreeBase().location.end.line < arguments.GetParseTreeBase().location.start.line {
		this_.priv_Parser_resetScannerAfter_7355(operand)
		result = operand
	} else {
		this_.priv_Parser_reportError_7441("'=>' unexpected", []any{})
		result = ToIface_ParseTree(callExpression)
	}
	return result
}

// Java method: Parser.parseAsyncArrowFunction(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseAsyncArrowFunction_7348(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eatPredefinedString_7326("async")
	if this_.priv_Parser_peekImplicitSemiColon_7411() {
		this_.priv_Parser_reportError_7441("No newline allowed between `async` and arrow function parameter list", []any{})
	}
	var arrowParameterList *FormalParameterListTree = nil
	_ = arrowParameterList
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_PAREN()) {
		arrowParameterList = this_.priv_Parser_parseFormalParameterList_7263()
	} else {
		var singleParameter *IdentifierExpressionTree = this_.priv_Parser_parseIdentifierExpression_7311()
		_ = singleParameter
		arrowParameterList = allocPooled_FormalParameterListTree_5(&this_.slab_FormalParameterListTree, singleParameter.GetParseTreeBase().location, ToIface_ParseTree(singleParameter), false, int32(0), (SourcePosition{}))
	}
	if this_.priv_Parser_peekImplicitSemiColon_7411() {
		this_.priv_Parser_reportError_7441("No newline allowed before '=>'", []any{})
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_ARROW())
	var arrowFunctionBody ParseTree = this_.priv_Parser_parseArrowFunctionBody_7349(expressionIn, StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS())
	_ = arrowFunctionBody
	var builder *FunctionDeclarationTree_Builder = allocPooledVia_FunctionDeclarationTree_S_Builder(&this_.slab_FunctionDeclarationTree_Builder, StaticGet_FunctionDeclarationTree_Kind_ARROW()).SetAsync(true).SetFormalParameterList(arrowParameterList).SetFunctionBody(arrowFunctionBody)
	_ = builder
	return ToIface_ParseTree(builder.buildWithPooledSlab_FunctionDeclarationTree(&this_.slab_FunctionDeclarationTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseArrowFunctionBody(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn, com.google.javascript.jscomp.parsing.parser.Parser.FunctionFlavor functionFlavor)
func (this_ *Parser) priv_Parser_parseArrowFunctionBody_7349(expressionIn Parser_Expression, functionFlavor Parser_FunctionFlavor) ParseTree {
	_ = expressionIn
	_ = functionFlavor
	var prevInGenerator bool = this_.inGenerator
	_ = prevInGenerator
	this_.inGenerator = functionFlavor.FieldGet_isGenerator()
	var arrowFunctionBody ParseTree = nil
	_ = arrowFunctionBody
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_CURLY()) {
		arrowFunctionBody = ToIface_ParseTree(this_.priv_Parser_parseFunctionBody_7265())
	} else {
		arrowFunctionBody = this_.priv_Parser_parseAssignment_7342(expressionIn)
	}
	this_.inGenerator = prevInGenerator
	return arrowFunctionBody
}

// Java method: Parser.newEmptyFormalParameterList(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func (this_ *Parser) priv_Parser_newEmptyFormalParameterList_7350(location SourceRange) *FormalParameterListTree {
	_ = location
	return allocPooled_FormalParameterListTree_1(&this_.slab_FormalParameterListTree, location)
}

// Java method: Parser.transformLeftHandSideExpression(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func (this_ *Parser) priv_Parser_transformLeftHandSideExpression_7351(tree ParseTree) ParseTree {
	_ = tree
	switch (tree.GetParseTreeBase().type_).Ordinal() {
	case int32(28) /* ParseTreeType.ARRAY_LITERAL_EXPRESSION */, int32(29) /* ParseTreeType.OBJECT_LITERAL_EXPRESSION */ :
		this_.priv_Parser_resetScanner_7354(tree)
		return this_.priv_Parser_parseLeftHandSidePattern_7352()
	default:
		return tree
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseLeftHandSidePattern()
func (this_ *Parser) priv_Parser_parseLeftHandSidePattern_7352() ParseTree {
	return this_.priv_Parser_parsePattern_7400(StaticGet_Parser_PatternKind_ANY())
}

// Java method: Parser.resetScanner(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start)
func (this_ *Parser) priv_Parser_resetScanner_7353(start SourcePosition) {
	_ = start
	this_.lastSourcePosition = start
	this_.scanner.SetPosition(this_.lastSourcePosition)
}

// Java method: Parser.resetScanner(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func (this_ *Parser) priv_Parser_resetScanner_7354(tree ParseTree) {
	_ = tree
	this_.scanner.SetPosition(tree.GetParseTreeBase().location.start)
}

// Java method: Parser.resetScannerAfter(com.google.javascript.jscomp.parsing.parser.trees.ParseTree parseTree)
func (this_ *Parser) priv_Parser_resetScannerAfter_7355(parseTree ParseTree) {
	_ = parseTree
	this_.lastSourcePosition = parseTree.GetParseTreeBase().location.end
	this_.scanner.SetPosition(this_.lastSourcePosition)
}

// Java method: Parser.peekAssignmentOperator()
func (this_ *Parser) priv_Parser_peekAssignmentOperator_7356() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(84) /* TokenType.EQUAL */, int32(87) /* TokenType.STAR_EQUAL */, int32(88) /* TokenType.STAR_STAR_EQUAL */, int32(97) /* TokenType.SLASH_EQUAL */, int32(89) /* TokenType.PERCENT_EQUAL */, int32(85) /* TokenType.PLUS_EQUAL */, int32(86) /* TokenType.MINUS_EQUAL */, int32(90) /* TokenType.LEFT_SHIFT_EQUAL */, int32(91) /* TokenType.RIGHT_SHIFT_EQUAL */, int32(92) /* TokenType.UNSIGNED_RIGHT_SHIFT_EQUAL */, int32(93) /* TokenType.AMPERSAND_EQUAL */, int32(95) /* TokenType.CARET_EQUAL */, int32(94) /* TokenType.BAR_EQUAL */, int32(99) /* TokenType.OR_EQUAL */, int32(98) /* TokenType.AND_EQUAL */, int32(100) /* TokenType.QUESTION_QUESTION_EQUAL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.inGeneratorContext()
func (this_ *Parser) priv_Parser_inGeneratorContext_7357() bool {
	return this_.inGenerator
}

// Java method: Parser.parseYield(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseYield_7358(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_YIELD())
	var isYieldAll bool = false
	_ = isYieldAll
	var expression ParseTree = nil
	_ = expression
	if !(this_.priv_Parser_peekImplicitSemiColon_7411()) {
		isYieldAll = !(this_.priv_Parser_eatOpt_7415(StaticGet_TokenType_STAR())).IsJavaNil_()
		if this_.priv_Parser_peekAssignmentExpression_7341() {
			expression = this_.priv_Parser_parseAssignment_7342(expressionIn)
		} else if isYieldAll {
			this_.priv_Parser_reportError_7441("yield* requires an expression", []any{})
		}
	}
	return ToIface_ParseTree(allocPooled_YieldExpressionTree(&this_.slab_YieldExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), isYieldAll, expression))
}

// Java method: Parser.parseConditional(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseConditional_7359(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var condition ParseTree = this_.priv_Parser_parseShortCircuit_7360(expressionIn)
	_ = condition
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_QUESTION())
		var left ParseTree = this_.priv_Parser_parseAssignment_7342(expressionIn)
		_ = left
		this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
		var right ParseTree = this_.priv_Parser_parseAssignment_7342(expressionIn)
		_ = right
		return ToIface_ParseTree(allocPooled_ConditionalExpressionTree(&this_.slab_ConditionalExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), condition, left, right))
	}
	return condition
}

// Java method: Parser.parseShortCircuit(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseShortCircuit_7360(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseLogicalOR_7362(expressionIn)
	_ = left
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_QUESTION()) {
		if left.GetParseTreeBase().type_ == StaticGet_ParseTreeType_BINARY_OPERATOR() {
			var binaryTree *BinaryOperatorTree = left.GetParseTreeBase().AsBinaryOperator()
			_ = binaryTree
			if (binaryTree.operator.type_ == StaticGet_TokenType_AND()) || (binaryTree.operator.type_ == StaticGet_TokenType_OR()) {
				this_.priv_Parser_reportError_7441("Logical OR and logical AND require parentheses when used with '??'", []any{})
			}
		}
		return this_.priv_Parser_parseNullishCoalesce_7361(expressionIn, left, start)
	} else {
		return left
	}
}

// Java method: Parser.parseNullishCoalesce(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn, com.google.javascript.jscomp.parsing.parser.trees.ParseTree left, com.google.javascript.jscomp.parsing.parser.util.SourcePosition start)
func (this_ *Parser) priv_Parser_parseNullishCoalesce_7361(expressionIn Parser_Expression, left ParseTree, start SourcePosition) ParseTree {
	_ = expressionIn
	_ = left
	_ = start
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_QUESTION()) {
		var operator Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_QUESTION_QUESTION())
		_ = operator
		var right ParseTree = this_.priv_Parser_parseBitwiseOR_7364(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_AND()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_OR()) {
		this_.priv_Parser_reportError_7441("Logical OR and logical AND require parentheses when used with '??'", []any{})
	}
	return left
}

// Java method: Parser.parseLogicalOR(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseLogicalOR_7362(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseLogicalAND_7363(expressionIn)
	_ = left
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_OR()) {
		var operator Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_OR())
		_ = operator
		var right ParseTree = this_.priv_Parser_parseLogicalAND_7363(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.parseLogicalAND(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseLogicalAND_7363(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseBitwiseOR_7364(expressionIn)
	_ = left
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_AND()) {
		var operator Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_AND())
		_ = operator
		var right ParseTree = this_.priv_Parser_parseBitwiseOR_7364(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.parseBitwiseOR(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseBitwiseOR_7364(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseBitwiseXOR_7365(expressionIn)
	_ = left
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_BAR()) {
		var operator Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_BAR())
		_ = operator
		var right ParseTree = this_.priv_Parser_parseBitwiseXOR_7365(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.parseBitwiseXOR(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseBitwiseXOR_7365(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseBitwiseAND_7366(expressionIn)
	_ = left
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_CARET()) {
		var operator Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_CARET())
		_ = operator
		var right ParseTree = this_.priv_Parser_parseBitwiseAND_7366(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.parseBitwiseAND(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseBitwiseAND_7366(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseEquality_7367(expressionIn)
	_ = left
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_AMPERSAND()) {
		var operator Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_AMPERSAND())
		_ = operator
		var right ParseTree = this_.priv_Parser_parseEquality_7367(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.parseEquality(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseEquality_7367(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseRelational_7369(expressionIn)
	_ = left
	for this_.priv_Parser_peekEqualityOperator_7368() {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseRelational_7369(expressionIn)
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.peekEqualityOperator()
func (this_ *Parser) priv_Parser_peekEqualityOperator_7368() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(59) /* TokenType.EQUAL_EQUAL */, int32(60) /* TokenType.NOT_EQUAL */, int32(61) /* TokenType.EQUAL_EQUAL_EQUAL */, int32(62) /* TokenType.NOT_EQUAL_EQUAL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseRelational(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_parseRelational_7369(expressionIn Parser_Expression) ParseTree {
	_ = expressionIn
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseShiftExpression_7371()
	_ = left
	for this_.priv_Parser_peekRelationalOperator_7370(expressionIn) {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseShiftExpression_7371()
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.peekRelationalOperator(com.google.javascript.jscomp.parsing.parser.Parser.Expression expressionIn)
func (this_ *Parser) priv_Parser_peekRelationalOperator_7370(expressionIn Parser_Expression) bool {
	_ = expressionIn
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(54) /* TokenType.OPEN_ANGLE */, int32(55) /* TokenType.CLOSE_ANGLE */, int32(57) /* TokenType.GREATER_EQUAL */, int32(56) /* TokenType.LESS_EQUAL */, int32(17) /* TokenType.INSTANCEOF */ :
		return true
	case int32(16) /* TokenType.IN */ :
		return (expressionIn == StaticGet_Parser_Expression_NORMAL())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseShiftExpression()
func (this_ *Parser) priv_Parser_parseShiftExpression_7371() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseAdditiveExpression_7373()
	_ = left
	for this_.priv_Parser_peekShiftOperator_7372() {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseAdditiveExpression_7373()
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.peekShiftOperator()
func (this_ *Parser) priv_Parser_peekShiftOperator_7372() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(70) /* TokenType.LEFT_SHIFT */, int32(71) /* TokenType.RIGHT_SHIFT */, int32(72) /* TokenType.UNSIGNED_RIGHT_SHIFT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseAdditiveExpression()
func (this_ *Parser) priv_Parser_parseAdditiveExpression_7373() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseMultiplicativeExpression_7375()
	_ = left
	for this_.priv_Parser_peekAdditiveOperator_7374() {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseMultiplicativeExpression_7375()
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.peekAdditiveOperator()
func (this_ *Parser) priv_Parser_peekAdditiveOperator_7374() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(63) /* TokenType.PLUS */, int32(64) /* TokenType.MINUS */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseMultiplicativeExpression()
func (this_ *Parser) priv_Parser_parseMultiplicativeExpression_7375() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseExponentiationExpression_7377()
	_ = left
	for this_.priv_Parser_peekMultiplicativeOperator_7376() {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseExponentiationExpression_7377()
		_ = right
		left = ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	}
	return left
}

// Java method: Parser.peekMultiplicativeOperator()
func (this_ *Parser) priv_Parser_peekMultiplicativeOperator_7376() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(65) /* TokenType.STAR */, int32(96) /* TokenType.SLASH */, int32(67) /* TokenType.PERCENT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseExponentiationExpression()
func (this_ *Parser) priv_Parser_parseExponentiationExpression_7377() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var left ParseTree = this_.priv_Parser_parseUnaryExpression_7378()
	_ = left
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_STAR_STAR()) {
		if left.GetParseTreeBase().type_ == StaticGet_ParseTreeType_UNARY_EXPRESSION() {
			this_.priv_Parser_reportError_7441("Unary operator '%s' requires parentheses before '**'", []any{j2gort.BoxNullableVal(left.GetParseTreeBase().AsUnaryExpression().operator)})
		} else if left.GetParseTreeBase().type_ == StaticGet_ParseTreeType_AWAIT_EXPRESSION() {
			this_.priv_Parser_reportError_7441("Unary operator 'await' requires parentheses before '**'", []any{})
		}
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var right ParseTree = this_.priv_Parser_parseExponentiationExpression_7377()
		_ = right
		return ToIface_ParseTree(allocPooled_BinaryOperatorTree(&this_.slab_BinaryOperatorTree, this_.priv_Parser_getTreeLocation_7429(start), left, operator, right))
	} else {
		return left
	}
}

// Java method: Parser.parseUnaryExpression()
func (this_ *Parser) priv_Parser_parseUnaryExpression_7378() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	if this_.priv_Parser_peekUnaryOperator_7379() {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var operand ParseTree = this_.priv_Parser_parseUnaryExpression_7378()
		_ = operand
		return ToIface_ParseTree(allocPooled_UnaryExpressionTree(&this_.slab_UnaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operator, operand))
	} else if this_.priv_Parser_peekAwaitExpression_7380() {
		return this_.priv_Parser_parseAwaitExpression_7381()
	} else {
		return this_.priv_Parser_parseUpdateExpression_7382()
	}
}

// Java method: Parser.peekUnaryOperator()
func (this_ *Parser) priv_Parser_peekUnaryOperator_7379() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(9) /* TokenType.DELETE */, int32(26) /* TokenType.VOID */, int32(24) /* TokenType.TYPEOF */, int32(63) /* TokenType.PLUS */, int32(64) /* TokenType.MINUS */, int32(77) /* TokenType.TILDE */, int32(76) /* TokenType.BANG */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.peekAwaitExpression()
func (this_ *Parser) priv_Parser_peekAwaitExpression_7380() bool {
	return this_.priv_Parser_peekPredefinedString_7325("await")
}

// Java method: Parser.parseAwaitExpression()
func (this_ *Parser) priv_Parser_parseAwaitExpression_7381() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eatPredefinedString_7326("await")
	var expression ParseTree = this_.priv_Parser_parseUnaryExpression_7378()
	_ = expression
	return ToIface_ParseTree(allocPooled_AwaitExpressionTree(&this_.slab_AwaitExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), expression))
}

// Java method: Parser.parseUpdateExpression()
func (this_ *Parser) priv_Parser_parseUpdateExpression_7382() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	if this_.priv_Parser_peekUpdateOperator_7383() {
		var operator Token = this_.priv_Parser_nextToken_7430()
		_ = operator
		var operand ParseTree = this_.priv_Parser_parseUnaryExpression_7378()
		_ = operand
		return ToIface_ParseTree(allocPooledVia_UpdateExpressionTree_Prefix(&this_.slab_UpdateExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operator, operand))
	} else {
		var lhs ParseTree = this_.priv_Parser_parseLeftHandSideExpression_7386()
		_ = lhs
		if this_.priv_Parser_peekUpdateOperator_7383() && !(this_.priv_Parser_peekImplicitSemiColon_7411()) {
			var operator_2 Token = this_.priv_Parser_nextToken_7430()
			_ = operator_2
			return ToIface_ParseTree(allocPooledVia_UpdateExpressionTree_Postfix(&this_.slab_UpdateExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operator_2, lhs))
		} else {
			return lhs
		}
	}
}

// Java method: Parser.peekUpdateOperator()
func (this_ *Parser) priv_Parser_peekUpdateOperator_7383() bool {
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(68) /* TokenType.PLUS_PLUS */, int32(69) /* TokenType.MINUS_MINUS */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.peekImportCall()
func (this_ *Parser) priv_Parser_peekImportCall_7384() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_IMPORT()) && this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_OPEN_PAREN()))
}

// Java method: Parser.peekImportDot()
func (this_ *Parser) priv_Parser_peekImportDot_7385() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_IMPORT()) && this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_PERIOD()))
}

// Java method: Parser.parseLeftHandSideExpression()
func (this_ *Parser) priv_Parser_parseLeftHandSideExpression_7386() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var operand ParseTree = this_.priv_Parser_parseNewExpression_7393()
	_ = operand
	var newExpressionTree *NewExpressionTree
	_ = newExpressionTree
	var ok_32 bool
	newExpressionTree, ok_32 = j2gort.TryCast[*NewExpressionTree](operand)
	ok_31 := !ok_32
	if ok_31 || (newExpressionTree.arguments != nil) {
		for this_.priv_Parser_peekCallSuffix_7387() {
			switch (this_.priv_Parser_peekType_7435()).Ordinal() {
			case int32(47) /* TokenType.OPEN_PAREN */ :
				var arguments *ArgumentListTree = this_.priv_Parser_parseArguments_7396()
				_ = arguments
				operand = ToIface_ParseTree(allocPooled_CallExpressionTree(&this_.slab_CallExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, arguments))
			case int32(49) /* TokenType.OPEN_SQUARE */ :
				this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
				var member ParseTree = this_.priv_Parser_parseExpression_7337()
				_ = member
				this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())
				operand = ToIface_ParseTree(allocPooled_MemberLookupExpressionTree(&this_.slab_MemberLookupExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, member))
			case int32(51) /* TokenType.PERIOD */ :
				this_.priv_Parser_eat_7425(StaticGet_TokenType_PERIOD())
				var id Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
				_ = id
				operand = ToIface_ParseTree(allocPooled_MemberExpressionTree(&this_.slab_MemberExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, id))
			case int32(112) /* TokenType.NO_SUBSTITUTION_TEMPLATE */, int32(109) /* TokenType.TEMPLATE_HEAD */ :
				operand = ToIface_ParseTree(this_.priv_Parser_parseTemplateLiteral_7313(operand))
			default:
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "unexpected case: "+j2gort.StrVal(this_.priv_Parser_peekType_7435()))))
			}
		}
		operand = this_.priv_Parser_maybeParseOptionalExpression_7388(operand)
	}
	return operand
}

// Java method: Parser.peekCallSuffix()
func (this_ *Parser) priv_Parser_peekCallSuffix_7387() bool {
	return ((((this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_PAREN()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_SQUARE())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_PERIOD())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_TEMPLATE_HEAD()))
}

// Java method: Parser.maybeParseOptionalExpression(com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand)
func (this_ *Parser) priv_Parser_maybeParseOptionalExpression_7388(operand ParseTree) ParseTree {
	_ = operand
	var start SourcePosition = operand.GetParseTreeBase().location.start
	_ = start
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_DOT()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_QUESTION_DOT())
		switch (this_.priv_Parser_peekType_7435()).Ordinal() {
		case int32(47) /* TokenType.OPEN_PAREN */ :
			var arguments *ArgumentListTree = this_.priv_Parser_parseArguments_7396()
			_ = arguments
			operand = ToIface_ParseTree(allocPooled_OptChainCallExpressionTree(&this_.slab_OptChainCallExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, arguments, true, arguments.hasTrailingComma))
		case int32(49) /* TokenType.OPEN_SQUARE */ :
			this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
			var member ParseTree = this_.priv_Parser_parseExpression_7337()
			_ = member
			this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())
			operand = ToIface_ParseTree(allocPooled_OptionalMemberLookupExpressionTree(&this_.slab_OptionalMemberLookupExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, member, true))
		case int32(112) /* TokenType.NO_SUBSTITUTION_TEMPLATE */, int32(109) /* TokenType.TEMPLATE_HEAD */ :
			this_.priv_Parser_reportError_7441("template literal cannot be used within optional chaining", []any{})
		default:
			if this_.priv_Parser_peekIdOrKeyword_7419() {
				var id Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
				_ = id
				operand = ToIface_ParseTree(allocPooled_OptionalMemberExpressionTree(&this_.slab_OptionalMemberExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, id, true))
			} else {
				this_.priv_Parser_reportError_7441("syntax error: %s not allowed in optional chain", []any{j2gort.BoxNullableVal(this_.priv_Parser_peekType_7435())})
			}
		}
		operand = this_.priv_Parser_parseRemainingOptionalChainSegment_7389(operand)
	}
	return operand
}

// Java method: Parser.parseRemainingOptionalChainSegment(com.google.javascript.jscomp.parsing.parser.trees.ParseTree optionalExpression)
func (this_ *Parser) priv_Parser_parseRemainingOptionalChainSegment_7389(optionalExpression ParseTree) ParseTree {
	_ = optionalExpression
	var start SourcePosition = optionalExpression.GetParseTreeBase().location.start
	_ = start
	for this_.priv_Parser_peekOptionalChainSuffix_7390() {
		if (this_.priv_Parser_peekType_7435() == StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE()) || (this_.priv_Parser_peekType_7435() == StaticGet_TokenType_TEMPLATE_HEAD()) {
			this_.priv_Parser_reportError_7441("template literal cannot be used within optional chaining", []any{})
			break
		}
		switch (this_.priv_Parser_peekType_7435()).Ordinal() {
		case int32(51) /* TokenType.PERIOD */ :
			this_.priv_Parser_eat_7425(StaticGet_TokenType_PERIOD())
			var id Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
			_ = id
			optionalExpression = ToIface_ParseTree(allocPooled_OptionalMemberExpressionTree(&this_.slab_OptionalMemberExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), optionalExpression, id, false))
		case int32(47) /* TokenType.OPEN_PAREN */ :
			var arguments *ArgumentListTree = this_.priv_Parser_parseArguments_7396()
			_ = arguments
			optionalExpression = ToIface_ParseTree(allocPooled_OptChainCallExpressionTree(&this_.slab_OptChainCallExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), optionalExpression, arguments, false, arguments.hasTrailingComma))
		case int32(49) /* TokenType.OPEN_SQUARE */ :
			this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
			var member ParseTree = this_.priv_Parser_parseExpression_7337()
			_ = member
			this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())
			optionalExpression = ToIface_ParseTree(allocPooled_OptionalMemberLookupExpressionTree(&this_.slab_OptionalMemberLookupExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), optionalExpression, member, false))
		default:
			panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "unexpected case: "+j2gort.StrVal(this_.priv_Parser_peekType_7435()))))
		}
	}
	return optionalExpression
}

// Java method: Parser.peekOptionalChainSuffix()
func (this_ *Parser) priv_Parser_peekOptionalChainSuffix_7390() bool {
	return ((((this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_PAREN()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_SQUARE())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_PERIOD())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_TEMPLATE_HEAD()))
}

// Java method: Parser.parseMemberExpressionNoNew()
func (this_ *Parser) priv_Parser_parseMemberExpressionNoNew_7391() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var operand ParseTree = nil
	_ = operand
	if this_.priv_Parser_peekImportDot_7385() {
		operand = this_.priv_Parser_parseImportDotMeta_7395()
	} else if this_.priv_Parser_peekAsyncFunctionStart_7250() {
		operand = this_.priv_Parser_parseAsyncFunctionExpression_7258()
	} else if this_.priv_Parser_peekFunction_7252() {
		operand = this_.priv_Parser_parseFunctionExpression_7256()
	} else {
		operand = this_.priv_Parser_parsePrimaryExpression_7307()
	}
	for this_.priv_Parser_peekMemberExpressionSuffix_7392() {
		switch (this_.priv_Parser_peekType_7435()).Ordinal() {
		case int32(49) /* TokenType.OPEN_SQUARE */ :
			this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
			var member ParseTree = this_.priv_Parser_parseExpression_7337()
			_ = member
			this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())
			operand = ToIface_ParseTree(allocPooled_MemberLookupExpressionTree(&this_.slab_MemberLookupExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, member))
		case int32(51) /* TokenType.PERIOD */ :
			this_.priv_Parser_eat_7425(StaticGet_TokenType_PERIOD())
			var id Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
			_ = id
			operand = ToIface_ParseTree(allocPooled_MemberExpressionTree(&this_.slab_MemberExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, id))
		case int32(112) /* TokenType.NO_SUBSTITUTION_TEMPLATE */, int32(109) /* TokenType.TEMPLATE_HEAD */ :
			operand = ToIface_ParseTree(this_.priv_Parser_parseTemplateLiteral_7313(operand))
		default:
			panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "unreachable")))
		}
	}
	return operand
}

// Java method: Parser.peekMemberExpressionSuffix()
func (this_ *Parser) priv_Parser_peekMemberExpressionSuffix_7392() bool {
	return (((this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_SQUARE()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_PERIOD())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE())) || this_.priv_Parser_peek_7433(StaticGet_TokenType_TEMPLATE_HEAD()))
}

// Java method: Parser.parseNewExpression()
func (this_ *Parser) priv_Parser_parseNewExpression_7393() ParseTree {
	if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_NEW())) {
		return this_.priv_Parser_parseMemberExpressionNoNew_7391()
	} else if this_.priv_Parser_peek_7434(int32(1), StaticGet_TokenType_PERIOD()) {
		return this_.priv_Parser_parseNewDotSomething_7394()
	} else {
		var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
		_ = start
		this_.priv_Parser_eat_7425(StaticGet_TokenType_NEW())
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_DOT()) {
			this_.priv_Parser_reportError_7441("Optional chaining is forbidden in `new?.target` contexts.", []any{})
		}
		var operand ParseTree = this_.priv_Parser_parseNewExpression_7393()
		_ = operand
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_QUESTION_DOT()) {
			this_.priv_Parser_reportError_7441("Optional chaining is forbidden in construction contexts.", []any{})
		}
		var arguments *ArgumentListTree = nil
		_ = arguments
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_PAREN()) {
			arguments = this_.priv_Parser_parseArguments_7396()
		}
		return ToIface_ParseTree(allocPooled_NewExpressionTree(&this_.slab_NewExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), operand, arguments, ((arguments != nil) && arguments.hasTrailingComma)))
	}
}

// Java method: Parser.parseNewDotSomething()
func (this_ *Parser) priv_Parser_parseNewDotSomething_7394() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_NEW())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_PERIOD())
	this_.priv_Parser_eatPredefinedString_7326("target")
	return ToIface_ParseTree(allocPooled_NewTargetExpressionTree(&this_.slab_NewTargetExpressionTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseImportDotMeta()
func (this_ *Parser) priv_Parser_parseImportDotMeta_7395() ParseTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_IMPORT())
	this_.priv_Parser_eat_7425(StaticGet_TokenType_PERIOD())
	this_.priv_Parser_eatPredefinedString_7326("meta")
	return ToIface_ParseTree(allocPooled_ImportMetaExpressionTree(&this_.slab_ImportMetaExpressionTree, this_.priv_Parser_getTreeLocation_7429(start)))
}

// Java method: Parser.parseArguments()
func (this_ *Parser) priv_Parser_parseArguments_7396() *ArgumentListTree {
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_PAREN())
	if !(this_.priv_Parser_peekAssignmentOrSpread_7397()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_1(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start))
	}
	var first ParseTree = this_.priv_Parser_parseAssignmentOrSpread_7398()
	_ = first
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_5(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, false, int32(0), (SourcePosition{}))
	}
	var trailingComma bool = false
	_ = trailingComma
	var firstComma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
	_ = firstComma
	var c0 SourcePosition
	_ = c0
	if !(firstComma).IsJavaNil_() {
		c0 = firstComma.GetStart()
	} else {
		c0 = (SourcePosition{})
	}
	var commaCount int32
	_ = commaCount
	if !(firstComma).IsJavaNil_() {
		commaCount = int32(1)
	} else {
		commaCount = int32(0)
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		if !(this_.config.atLeast8) {
			this_.priv_Parser_reportError_7439(firstComma, "Invalid trailing comma in arguments list", []any{})
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_5(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, true, commaCount, c0)
	}
	if !(this_.priv_Parser_peekAssignmentOrSpread_7397()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_5(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, false, commaCount, c0)
	}
	var second ParseTree = this_.priv_Parser_parseAssignmentOrSpread_7398()
	_ = second
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_7(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, second, false, commaCount, c0, (SourcePosition{}))
	}
	var secondComma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
	_ = secondComma
	var c1 SourcePosition = (SourcePosition{})
	_ = c1
	if !(secondComma).IsJavaNil_() {
		var pos SourcePosition = secondComma.GetStart()
		_ = pos
		if commaCount == int32(0) {
			c0 = pos
		} else {
			c1 = pos
		}
		commaCount = commaCount + 1
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		if !(this_.config.atLeast8) {
			this_.priv_Parser_reportError_7439(secondComma, "Invalid trailing comma in arguments list", []any{})
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_7(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, second, true, commaCount, c0, c1)
	}
	if !(this_.priv_Parser_peekAssignmentOrSpread_7397()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_7(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, second, false, commaCount, c0, c1)
	}
	var third ParseTree = this_.priv_Parser_parseAssignmentOrSpread_7398()
	_ = third
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_9(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, second, third, false, commaCount, c0, c1, (SourcePosition{}))
	}
	var thirdComma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
	_ = thirdComma
	var c2 SourcePosition = (SourcePosition{})
	_ = c2
	if !(thirdComma).IsJavaNil_() {
		var pos_2 SourcePosition = thirdComma.GetStart()
		_ = pos_2
		if commaCount == int32(0) {
			c0 = pos_2
		} else if commaCount == int32(1) {
			c1 = pos_2
		} else {
			c2 = pos_2
		}
		commaCount = commaCount + 1
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
		if !(this_.config.atLeast8) {
			this_.priv_Parser_reportError_7439(thirdComma, "Invalid trailing comma in arguments list", []any{})
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_9(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, second, third, true, commaCount, c0, c1, c2)
	}
	if !(this_.priv_Parser_peekAssignmentOrSpread_7397()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		return allocPooled_ArgumentListTree_9(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), first, second, third, false, commaCount, c0, c1, c2)
	}
	var savedTreeTop int32 = this_.treeScratchTop
	_ = savedTreeTop
	var savedCommaTop int32 = this_.commaScratchTop
	_ = savedCommaTop
	defer func() {
		{
			var i int32 = savedTreeTop
			_ = i
			var end int32 = this_.treeScratchTop
			_ = end
			for ; i < end; i = i + 1 {
				{
					a_33 := this_.treeScratch
					i_34 := i
					a_33[i_34] = nil
				}
			}
		}
		this_.treeScratchTop = savedTreeTop
		this_.commaScratchTop = savedCommaTop
	}()
	{
		var argCount int32 = int32(3)
		_ = argCount
		for this_.priv_Parser_peekAssignmentOrSpread_7397() {
			this_.priv_Parser_pushTreeScratch_7261(this_.priv_Parser_parseAssignmentOrSpread_7398())
			argCount = argCount + 1
			if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN())) {
				var comma Token = this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
				_ = comma
				if !(comma).IsJavaNil_() {
					var pos_3 SourcePosition = comma.GetStart()
					_ = pos_3
					if commaCount == int32(0) {
						c0 = pos_3
					} else if commaCount == int32(1) {
						c1 = pos_3
					} else if commaCount == int32(2) {
						c2 = pos_3
					} else {
						this_.priv_Parser_pushCommaScratch_7262(pos_3)
					}
					commaCount = commaCount + 1
				}
				if this_.priv_Parser_peek_7433(StaticGet_TokenType_CLOSE_PAREN()) {
					if !(this_.config.atLeast8) {
						this_.priv_Parser_reportError_7439(comma, "Invalid trailing comma in arguments list", []any{})
					}
					trailingComma = true
				}
			}
		}
		this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_PAREN())
		var restArgs []ParseTree
		{
			arrSrc_35 := this_.treeScratch
			arrFrom_36 := savedTreeTop
			arrTo_37 := this_.treeScratchTop
			j2gort.CheckNotNullSlice(arrSrc_35)
			if arrFrom_36 > arrTo_37 {
				panic(javalang.NewIllegalArgumentException("fromIndex > toIndex"))
			}
			if arrFrom_36 < 0 || int(arrFrom_36) > len(arrSrc_35) {
				panic(javalang.NewArrayIndexOutOfBoundsException(""))
			}
			arrLen_38 := arrTo_37 - arrFrom_36
			arrOut_39 := make([]ParseTree, arrLen_38)
			copy(arrOut_39, arrSrc_35[arrFrom_36:])
			restArgs = arrOut_39
		}

		_ = restArgs
		var restCommas []SourcePosition
		_ = restCommas
		if commaCount > int32(3) {
			{

				arrSrc_40 := this_.commaScratch
				arrFrom_41 := savedCommaTop
				arrTo_42 := this_.commaScratchTop
				j2gort.CheckNotNullSlice(arrSrc_40)
				if arrFrom_41 > arrTo_42 {
					panic(javalang.NewIllegalArgumentException("fromIndex > toIndex"))
				}
				if arrFrom_41 < 0 || int(arrFrom_41) > len(arrSrc_40) {
					panic(javalang.NewArrayIndexOutOfBoundsException(""))
				}
				arrLen_43 := arrTo_42 - arrFrom_41
				arrOut_44 := make([]SourcePosition, arrLen_43)
				copy(arrOut_44, arrSrc_40[arrFrom_41:])
				restCommas = arrOut_44
			}

		} else {
			restCommas = nil
		}
		return allocPooled_ArgumentListTree_12(&this_.slab_ArgumentListTree, this_.priv_Parser_getTreeLocation_7429(start), argCount, first, second, third, restArgs, trailingComma, commaCount, c0, c1, c2, restCommas)
	}
	panic("unreachable")
}

// Java method: Parser.peekAssignmentOrSpread()
func (this_ *Parser) priv_Parser_peekAssignmentOrSpread_7397() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) || this_.priv_Parser_peekAssignmentExpression_7341())
}

// Java method: Parser.parseAssignmentOrSpread()
func (this_ *Parser) priv_Parser_parseAssignmentOrSpread_7398() ParseTree {
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
		return ToIface_ParseTree(this_.priv_Parser_parseIterSpread_7267())
	}
	return this_.priv_Parser_parseAssignmentExpression_7340()
}

// Java method: Parser.peekPatternStart()
func (this_ *Parser) priv_Parser_peekPatternStart_7399() bool {
	return (this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_SQUARE()) || this_.priv_Parser_peek_7433(StaticGet_TokenType_OPEN_CURLY()))
}

// Java method: Parser.parsePattern(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind kind)
func (this_ *Parser) priv_Parser_parsePattern_7400(kind Parser_PatternKind) ParseTree {
	_ = kind
	switch (this_.priv_Parser_peekType_7435()).Ordinal() {
	case int32(49) /* TokenType.OPEN_SQUARE */ :
		return this_.priv_Parser_parseArrayPattern_7404(kind)
	case int32(45) /* TokenType.OPEN_CURLY */ :
		fallthrough
	default:
		return this_.priv_Parser_parseObjectPattern_7405(kind)
	}
	panic("unreachable")
}

// Java method: Parser.peekArrayPatternElement()
func (this_ *Parser) priv_Parser_peekArrayPatternElement_7401() bool {
	return this_.priv_Parser_peekExpression_7338()
}

// Java method: Parser.parseIterRest(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind patternKind)
func (this_ *Parser) priv_Parser_parseIterRest_7402(patternKind Parser_PatternKind) ParseTree {
	_ = patternKind
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	this_.priv_Parser_eat_7425(StaticGet_TokenType_ELLIPSIS())
	var patternAssignmentTarget ParseTree = this_.priv_Parser_parseRestAssignmentTarget_7403(patternKind)
	_ = patternAssignmentTarget
	return ToIface_ParseTree(allocPooled_IterRestTree(&this_.slab_IterRestTree, this_.priv_Parser_getTreeLocation_7429(start), patternAssignmentTarget))
}

// Java method: Parser.parseRestAssignmentTarget(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind patternKind)
func (this_ *Parser) priv_Parser_parseRestAssignmentTarget_7403(patternKind Parser_PatternKind) ParseTree {
	_ = patternKind
	var patternAssignmentTarget ParseTree = this_.priv_Parser_parsePatternAssignmentTargetNoDefault_7409(patternKind)
	_ = patternAssignmentTarget
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
		this_.priv_Parser_reportError_7441("A default value cannot be specified after '...'", []any{})
	}
	return patternAssignmentTarget
}

// Java method: Parser.parseArrayPattern(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind kind)
func (this_ *Parser) priv_Parser_parseArrayPattern_7404(kind Parser_PatternKind) ParseTree {
	_ = kind
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var elements []ParseTree = []ParseTree{}
	_ = elements
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_SQUARE())
	for this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) || this_.priv_Parser_peekArrayPatternElement_7401() {
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
			var nullStart SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
			_ = nullStart
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
			elements = append(elements, allocPooled_NullTree(&this_.slab_NullTree, this_.priv_Parser_getTreeLocation_7429(nullStart)))
		} else {
			elements = append(elements, j2gort.CheckNotNullElem(this_.priv_Parser_parsePatternAssignmentTarget_7408(kind)))
			if this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
				this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
			} else {
				break
			}
		}
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST())
		elements = append(elements, j2gort.CheckNotNullElem(this_.priv_Parser_parseIterRest_7402(kind)))
	}
	if (this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_SQUARE())).IsJavaNil_() {
		return ToIface_ParseTree(allocPooled_MissingPrimaryExpressionTree(&this_.slab_MissingPrimaryExpressionTree, this_.priv_Parser_getTreeLocation_7429(this_.priv_Parser_getTreeStartLocation_7427())))
	}
	return ToIface_ParseTree(allocPooled_ArrayPatternTree(&this_.slab_ArrayPatternTree, this_.priv_Parser_getTreeLocation_7429(start), j2gort.ImmListFromTypedSlice(elements)))
}

// Java method: Parser.parseObjectPattern(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind kind)
func (this_ *Parser) priv_Parser_parseObjectPattern_7405(kind Parser_PatternKind) ParseTree {
	_ = kind
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var fields []ParseTree = []ParseTree{}
	_ = fields
	this_.priv_Parser_eat_7425(StaticGet_TokenType_OPEN_CURLY())
	for this_.priv_Parser_peekObjectPatternField_7406() {
		fields = append(fields, j2gort.CheckNotNullElem(this_.priv_Parser_parseObjectPatternField_7407(kind)))
		if this_.priv_Parser_peek_7433(StaticGet_TokenType_COMMA()) {
			this_.priv_Parser_eat_7425(StaticGet_TokenType_COMMA())
		} else {
			break
		}
	}
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_ELLIPSIS()) {
		this_.priv_Parser_recordFeatureUsed_7444(StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST())
		var restStart SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
		_ = restStart
		this_.priv_Parser_eat_7425(StaticGet_TokenType_ELLIPSIS())
		var patternAssignmentTarget ParseTree = this_.priv_Parser_parseRestAssignmentTarget_7403(kind)
		_ = patternAssignmentTarget
		fields = append(fields, allocPooled_ObjectRestTree(&this_.slab_ObjectRestTree, this_.priv_Parser_getTreeLocation_7429(restStart), patternAssignmentTarget))
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_CLOSE_CURLY())
	return ToIface_ParseTree(allocPooled_ObjectPatternTree(&this_.slab_ObjectPatternTree, this_.priv_Parser_getTreeLocation_7429(start), j2gort.ImmListFromTypedSlice(fields)))
}

// Java method: Parser.peekObjectPatternField()
func (this_ *Parser) priv_Parser_peekObjectPatternField_7406() bool {
	return this_.priv_Parser_peekPropertyNameOrComputedProp_7319(int32(0))
}

// Java method: Parser.parseObjectPatternField(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind kind)
func (this_ *Parser) priv_Parser_parseObjectPatternField_7407(kind Parser_PatternKind) ParseTree {
	_ = kind
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	if this_.priv_Parser_peekType_7435() == StaticGet_TokenType_OPEN_SQUARE() {
		var key ParseTree = this_.priv_Parser_parseComputedPropertyName_7323()
		_ = key
		this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
		var value ParseTree = this_.priv_Parser_parsePatternAssignmentTarget_7408(kind)
		_ = value
		return ToIface_ParseTree(allocPooled_ComputedPropertyDefinitionTree(&this_.slab_ComputedPropertyDefinitionTree, this_.priv_Parser_getTreeLocation_7429(start), key, value))
	}
	var name Token = (Token{})
	_ = name
	if this_.priv_Parser_peekIdOrKeyword_7419() {
		var idToken Token = this_.priv_Parser_eatIdOrKeywordAsId_7424()
		_ = idToken
		if !(this_.priv_Parser_peek_7433(StaticGet_TokenType_COLON())) {
			if idToken.IsKeyword() {
				this_.priv_Parser_reportError_7441("cannot use keyword '%s' here.", []any{j2gort.BoxNullableVal(idToken)})
			}
			if this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
				var idTree *IdentifierExpressionTree = allocPooled_IdentifierExpressionTree(&this_.slab_IdentifierExpressionTree, this_.priv_Parser_getTreeLocation_7429(start), idToken)
				_ = idTree
				this_.priv_Parser_eat_7425(StaticGet_TokenType_EQUAL())
				var defaultValue ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
				_ = defaultValue
				return ToIface_ParseTree(allocPooled_DefaultParameterTree(&this_.slab_DefaultParameterTree, this_.priv_Parser_getTreeLocation_7429(start), ToIface_ParseTree(idTree), defaultValue))
			}
			return ToIface_ParseTree(allocPooled_PropertyNameAssignmentTree(&this_.slab_PropertyNameAssignmentTree, this_.priv_Parser_getTreeLocation_7429(start), idToken, nil))
		}
		name = idToken
	} else {
		name = this_.priv_Parser_parseLiteralExpression_7312().literalToken
	}
	this_.priv_Parser_eat_7425(StaticGet_TokenType_COLON())
	var value_2 ParseTree = this_.priv_Parser_parsePatternAssignmentTarget_7408(kind)
	_ = value_2
	return ToIface_ParseTree(allocPooled_PropertyNameAssignmentTree(&this_.slab_PropertyNameAssignmentTree, this_.priv_Parser_getTreeLocation_7429(start), name, value_2))
}

// Java method: Parser.parsePatternAssignmentTarget(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind patternKind)
func (this_ *Parser) priv_Parser_parsePatternAssignmentTarget_7408(patternKind Parser_PatternKind) ParseTree {
	_ = patternKind
	var start SourcePosition = this_.priv_Parser_getTreeStartLocation_7427()
	_ = start
	var assignmentTarget ParseTree = nil
	_ = assignmentTarget
	assignmentTarget = this_.priv_Parser_parsePatternAssignmentTargetNoDefault_7409(patternKind)
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_EQUAL()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_EQUAL())
		var defaultValue ParseTree = this_.priv_Parser_parseAssignmentExpression_7340()
		_ = defaultValue
		assignmentTarget = ToIface_ParseTree(allocPooled_DefaultParameterTree(&this_.slab_DefaultParameterTree, this_.priv_Parser_getTreeLocation_7429(start), assignmentTarget, defaultValue))
	}
	return assignmentTarget
}

// Java method: Parser.parsePatternAssignmentTargetNoDefault(com.google.javascript.jscomp.parsing.parser.Parser.PatternKind kind)
func (this_ *Parser) priv_Parser_parsePatternAssignmentTargetNoDefault_7409(kind Parser_PatternKind) ParseTree {
	_ = kind
	var assignmentTarget ParseTree = nil
	_ = assignmentTarget
	if this_.priv_Parser_peekPatternStart_7399() {
		assignmentTarget = this_.priv_Parser_parsePattern_7400(kind)
	} else {
		assignmentTarget = this_.priv_Parser_parseLeftHandSideExpression_7386()
		if !(assignmentTarget.GetParseTreeBase().IsValidAssignmentTarget()) {
			this_.priv_Parser_reportError_7441("invalid assignment target", []any{})
		}
		if (kind == StaticGet_Parser_PatternKind_INITIALIZER()) && (assignmentTarget.GetParseTreeBase().type_ != StaticGet_ParseTreeType_IDENTIFIER_EXPRESSION()) {
			this_.priv_Parser_reportError_7441("Only an identifier or destructuring pattern is allowed here.", []any{})
		}
	}
	return assignmentTarget
}

// Java method: Parser.eatPossiblyImplicitSemiColon()
func (this_ *Parser) priv_Parser_eatPossiblyImplicitSemiColon_7410() {
	if this_.priv_Parser_peek_7433(StaticGet_TokenType_SEMI_COLON()) {
		this_.priv_Parser_eat_7425(StaticGet_TokenType_SEMI_COLON())
		return
	}
	if this_.priv_Parser_peekImplicitSemiColon_7411() {
		return
	}
	this_.priv_Parser_reportError_7441("Semi-colon expected", []any{})
}

// Java method: Parser.peekImplicitSemiColon()
func (this_ *Parser) priv_Parser_peekImplicitSemiColon_7411() bool {
	return this_.priv_Parser_peekImplicitSemiColon_7412(int32(0))
}

// Java method: Parser.peekImplicitSemiColon(int index)
func (this_ *Parser) priv_Parser_peekImplicitSemiColon_7412(index int32) bool {
	_ = index
	var lineAdvanced bool = false
	_ = lineAdvanced
	if index == int32(0) {
		lineAdvanced = (this_.priv_Parser_getNextLine_7414() > this_.priv_Parser_getLastLine_7413())
	} else {
		lineAdvanced = (this_.priv_Parser_peekToken_7438(index).location.start.line > this_.priv_Parser_peekToken_7438((index - int32(1))).location.end.line)
	}
	return (((lineAdvanced || this_.priv_Parser_peek_7434(index, StaticGet_TokenType_SEMI_COLON())) || this_.priv_Parser_peek_7434(index, StaticGet_TokenType_CLOSE_CURLY())) || this_.priv_Parser_peek_7434(index, StaticGet_TokenType_END_OF_FILE()))
}

// Java method: Parser.getLastLine()
func (this_ *Parser) priv_Parser_getLastLine_7413() int32 {
	return this_.lastSourcePosition.line
}

// Java method: Parser.getNextLine()
func (this_ *Parser) priv_Parser_getNextLine_7414() int32 {
	return this_.priv_Parser_peekToken_7437().location.start.line
}

// Java method: Parser.eatOpt(com.google.javascript.jscomp.parsing.parser.TokenType expectedTokenType)
func (this_ *Parser) priv_Parser_eatOpt_7415(expectedTokenType TokenType) Token {
	_ = expectedTokenType
	if this_.priv_Parser_peek_7433(expectedTokenType) {
		return this_.priv_Parser_eat_7425(expectedTokenType)
	}
	return (Token{})
}

// Java method: Parser.inStrictContext()
func (this_ *Parser) priv_Parser_inStrictContext_7416() bool {
	return this_.config.isStrictMode
}

// Java method: Parser.peekId()
func (this_ *Parser) priv_Parser_peekId_7417() bool {
	return this_.priv_Parser_peekId_7418(int32(0))
}

// Java method: Parser.peekId(int index)
func (this_ *Parser) priv_Parser_peekId_7418(index int32) bool {
	_ = index
	var token Token = this_.priv_Parser_peekToken_7438(index)
	_ = token
	var typeToCheck TokenType = token.type_
	_ = typeToCheck
	if typeToCheck == StaticGet_TokenType_IDENTIFIER() {
		var value string = token.GetMaybePrivateValue()
		_ = value
		var k *Keywords = Keywords_Get_1_string(value)
		_ = k
		if k == nil {
			return true
		}
		typeToCheck = k.type_
	}
	return (!(this_.priv_Parser_inStrictContext_7416()) && Keywords_IsStrictKeyword(typeToCheck))
}

// Java method: Parser.peekIdOrKeyword()
func (this_ *Parser) priv_Parser_peekIdOrKeyword_7419() bool {
	return this_.priv_Parser_peekIdOrKeyword_7420(int32(0))
}

// Java method: Parser.peekIdOrKeyword(int index)
func (this_ *Parser) priv_Parser_peekIdOrKeyword_7420(index int32) bool {
	_ = index
	var type_ TokenType = this_.priv_Parser_peekType_7436(index)
	_ = type_
	return ((StaticGet_TokenType_IDENTIFIER() == type_) || Keywords_IsKeyword_1_TokenType(type_))
}

// Java method: Parser.eatIdOpt()
func (this_ *Parser) priv_Parser_eatIdOpt_7421() Token {
	if this_.priv_Parser_peekId_7417() {
		return this_.priv_Parser_eatIdOrKeywordAsId_7424()
	} else {
		return (Token{})
	}
}

// Java method: Parser.eatId()
func (this_ *Parser) priv_Parser_eatId_7422() Token {
	if this_.priv_Parser_peekId_7417() {
		return this_.priv_Parser_eatIdOrKeywordAsId_7424()
	} else {
		this_.priv_Parser_reportExpectedError_7426(this_.priv_Parser_peekToken_7437(), j2gort.BoxNullableVal(StaticGet_TokenType_IDENTIFIER()))
		if this_.priv_Parser_peekIdOrKeyword_7419() {
			return this_.priv_Parser_eatIdOrKeywordAsId_7424()
		} else {
			return (Token{})
		}
	}
}

// Java method: Parser.eatObjectLiteralPropertyName()
func (this_ *Parser) priv_Parser_eatObjectLiteralPropertyName_7423() Token {
	var token Token = this_.priv_Parser_peekToken_7437()
	_ = token
	switch (token.type_).Ordinal() {
	case int32(105) /* TokenType.STRING */, int32(104) /* TokenType.NUMBER */, int32(106) /* TokenType.BIGINT */ :
		return this_.priv_Parser_nextToken_7430()
	case int32(2) /* TokenType.IDENTIFIER */ :
		fallthrough
	default:
		return this_.priv_Parser_eatIdOrKeywordAsId_7424()
	}
	panic("unreachable")
}

// Java method: Parser.eatIdOrKeywordAsId()
func (this_ *Parser) priv_Parser_eatIdOrKeywordAsId_7424() Token {
	var token Token = this_.priv_Parser_nextToken_7430()
	_ = token
	if token.type_ == StaticGet_TokenType_IDENTIFIER() {
		return token
	} else if Keywords_IsKeyword_1_TokenType(token.type_) {
		return Token_Identifier(token.location, Keywords_Get_1_TokenType(token.type_).String())
	} else {
		this_.priv_Parser_reportExpectedError_7426(token, j2gort.BoxNullableVal(StaticGet_TokenType_IDENTIFIER()))
	}
	return (Token{})
}

// Java method: Parser.eat(com.google.javascript.jscomp.parsing.parser.TokenType expectedTokenType)
func (this_ *Parser) priv_Parser_eat_7425(expectedTokenType TokenType) Token {
	_ = expectedTokenType
	var token Token = this_.priv_Parser_nextToken_7430()
	_ = token
	if token.type_ != expectedTokenType {
		this_.priv_Parser_reportExpectedError_7426(token, j2gort.BoxNullableVal(expectedTokenType))
		return (Token{})
	}
	return token
}

// Java method: Parser.reportExpectedError(com.google.javascript.jscomp.parsing.parser.@org.jspecify.annotations.Nullable Token token, java.lang.Object expected)
func (this_ *Parser) priv_Parser_reportExpectedError_7426(token Token, expected any) {
	_ = token
	_ = expected
	this_.priv_Parser_reportError_7439(token, "'%s' expected", []any{expected})
}

// Java method: Parser.getTreeStartLocation()
func (this_ *Parser) priv_Parser_getTreeStartLocation_7427() SourcePosition {
	return this_.priv_Parser_peekToken_7437().location.start
}

// Java method: Parser.getTreeEndLocation()
func (this_ *Parser) priv_Parser_getTreeEndLocation_7428() SourcePosition {
	return this_.lastSourcePosition
}

// Java method: Parser.getTreeLocation(com.google.javascript.jscomp.parsing.parser.util.SourcePosition start)
func (this_ *Parser) priv_Parser_getTreeLocation_7429(start SourcePosition) SourceRange {
	_ = start
	return NewSourceRange(start, this_.priv_Parser_getTreeEndLocation_7428())
}

// Java method: Parser.nextToken()
func (this_ *Parser) priv_Parser_nextToken_7430() Token {
	var token Token = this_.scanner.NextToken()
	_ = token
	this_.lastSourcePosition = token.location.end
	return token
}

// Java method: Parser.nextRegularExpressionLiteralToken()
func (this_ *Parser) priv_Parser_nextRegularExpressionLiteralToken_7431() Token {
	var token Token = this_.scanner.NextRegularExpressionLiteralToken()
	_ = token
	this_.lastSourcePosition = token.location.end
	return token
}

// Java method: Parser.nextTemplateLiteralToken()
func (this_ *Parser) priv_Parser_nextTemplateLiteralToken_7432() Token {
	var token Token = this_.scanner.NextTemplateLiteralToken()
	_ = token
	this_.lastSourcePosition = token.location.end
	return token
}

// Java method: Parser.peek(com.google.javascript.jscomp.parsing.parser.TokenType expectedType)
func (this_ *Parser) priv_Parser_peek_7433(expectedType TokenType) bool {
	_ = expectedType
	return this_.priv_Parser_peek_7434(int32(0), expectedType)
}

// Java method: Parser.peek(int index, com.google.javascript.jscomp.parsing.parser.TokenType expectedType)
func (this_ *Parser) priv_Parser_peek_7434(index int32, expectedType TokenType) bool {
	_ = index
	_ = expectedType
	return (this_.priv_Parser_peekType_7436(index) == expectedType)
}

// Java method: Parser.peekType()
func (this_ *Parser) priv_Parser_peekType_7435() TokenType {
	return this_.priv_Parser_peekType_7436(int32(0))
}

// Java method: Parser.peekType(int index)
func (this_ *Parser) priv_Parser_peekType_7436(index int32) TokenType {
	_ = index
	return this_.priv_Parser_peekToken_7438(index).type_
}

// Java method: Parser.peekToken()
func (this_ *Parser) priv_Parser_peekToken_7437() Token {
	return this_.priv_Parser_peekToken_7438(int32(0))
}

// Java method: Parser.peekToken(int index)
func (this_ *Parser) priv_Parser_peekToken_7438(index int32) Token {
	_ = index
	return this_.scanner.PeekToken_1(index)
}

// Java method: Parser.reportError(com.google.javascript.jscomp.parsing.parser.Token token, java.lang.String message, java.lang.Object[] arguments)
func (this_ *Parser) priv_Parser_reportError_7439(token Token, message string, arguments []any) {
	_ = token
	_ = message
	_ = arguments
	if (token).IsJavaNil_() {
		this_.priv_Parser_reportError_7441(message, arguments)
	} else {
		this_.errorReporter.GetErrorReporterBase().ReportError_3(token.GetStart(), message, arguments)
	}
}

// Java method: Parser.reportError(com.google.javascript.jscomp.parsing.parser.trees.ParseTree parseTree, java.lang.String message, java.lang.Object[] arguments)
func (this_ *Parser) priv_Parser_reportError_7440(parseTree ParseTree, message string, arguments []any) {
	_ = parseTree
	_ = message
	_ = arguments
	if parseTree == nil || parseTree.IsJavaNil_() {
		this_.priv_Parser_reportError_7441(message, arguments)
	} else {
		this_.errorReporter.GetErrorReporterBase().ReportError_3(parseTree.GetParseTreeBase().location.start, message, arguments)
	}
}

// Java method: Parser.reportError(java.lang.String message, java.lang.Object[] arguments)
func (this_ *Parser) priv_Parser_reportError_7441(message string, arguments []any) {
	_ = message
	_ = arguments
	this_.errorReporter.GetErrorReporterBase().ReportError_3(this_.scanner.GetPosition_0(), message, arguments)
}

// Java method: Parser.reportError(com.google.javascript.jscomp.parsing.parser.util.SourcePosition position, java.lang.String message, java.lang.Object[] arguments)
func (this_ *Parser) priv_Parser_reportError_7442(position SourcePosition, message string, arguments []any) {
	_ = position
	_ = message
	_ = arguments
	this_.errorReporter.GetErrorReporterBase().ReportError_3(position, message, arguments)
}

// Java method: Parser.reportTemplateErrorIfPresent(com.google.javascript.jscomp.parsing.parser.Token templateToken)
func (this_ *Parser) priv_Parser_reportTemplateErrorIfPresent_7443(templateToken Token) {
	_ = templateToken
	var errorMessage string = templateToken.GetErrorMessage()
	_ = errorMessage
	if errorMessage == j2gort.NullString {
		return
	}
	switch (templateToken.GetErrorLevel()).Ordinal() {
	case int32(0) /* ErrorLevel.WARNING */ :
		this_.errorReporter.GetErrorReporterBase().ReportWarning_3(templateToken.GetErrorPosition(), "%s", []any{j2gort.BoxString(errorMessage)})
		return
	case int32(1) /* ErrorLevel.ERROR */ :
		this_.priv_Parser_reportError_7442(templateToken.GetErrorPosition(), "%s", []any{j2gort.BoxString(errorMessage)})
		return
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: Parser.recordFeatureUsed(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *Parser) priv_Parser_recordFeatureUsed_7444(feature *FeatureSet_Feature) *Parser {
	_ = feature
	this_.features = this_.features.With_1_FeatureSet_Feature(feature)
	return this_
}

// Java enum: com.google.javascript.jscomp.parsing.parser.Parser.FunctionFlavor
type Parser_FunctionFlavor uint8

const (
	enum_Parser_FunctionFlavor_NORMAL                 Parser_FunctionFlavor = 1
	enum_Parser_FunctionFlavor_GENERATOR              Parser_FunctionFlavor = 2
	enum_Parser_FunctionFlavor_ASYNCHRONOUS           Parser_FunctionFlavor = 3
	enum_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR Parser_FunctionFlavor = 4
)

var enumBaseTable_Parser_FunctionFlavor = [...]j2gort.EnumBase{
	{Name_: "NORMAL", Ordinal_: 0},
	{Name_: "GENERATOR", Ordinal_: 1},
	{Name_: "ASYNCHRONOUS", Ordinal_: 2},
	{Name_: "ASYNCHRONOUS_GENERATOR", Ordinal_: 3},
}
var enumValues_Parser_FunctionFlavor = []Parser_FunctionFlavor{1, 2, 3, 4}

var enumField_Parser_FunctionFlavor_isGenerator = [...]bool{
	false,
	true,
	false,
	true,
}

func (this_ Parser_FunctionFlavor) FieldGet_isGenerator() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_Parser_FunctionFlavor_isGenerator[this_-1]
}

func StaticGet_Parser_FunctionFlavor_NORMAL() Parser_FunctionFlavor {
	return enum_Parser_FunctionFlavor_NORMAL
}
func StaticGet_Parser_FunctionFlavor_GENERATOR() Parser_FunctionFlavor {
	return enum_Parser_FunctionFlavor_GENERATOR
}
func StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS() Parser_FunctionFlavor {
	return enum_Parser_FunctionFlavor_ASYNCHRONOUS
}
func StaticGet_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR() Parser_FunctionFlavor {
	return enum_Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR
}

func Parser_FunctionFlavor_Values() []Parser_FunctionFlavor {
	return j2gort.CloneSlice(enumValues_Parser_FunctionFlavor)
}

func Parser_FunctionFlavor_ValueOf(name string) Parser_FunctionFlavor {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Parser_FunctionFlavor {
		if enumBaseTable_Parser_FunctionFlavor[i_].Name_ == name {
			return Parser_FunctionFlavor(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.Parser.FunctionFlavor." + name))
}

func (this_ Parser_FunctionFlavor) IsJavaNil_() bool                              { return this_ == 0 }
func (Parser_FunctionFlavor) HasIdentityEquals_()                                 {}
func (Parser_FunctionFlavor) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Parser_FunctionFlavor) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Parser_FunctionFlavor) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Parser_FunctionFlavor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Parser$FunctionFlavor", "FunctionFlavor"
}
func (this_ Parser_FunctionFlavor) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Parser_FunctionFlavor[this_-1]
}
func (this_ Parser_FunctionFlavor) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Parser_FunctionFlavor) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Parser_FunctionFlavor) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_FunctionFlavor[this_-1].Name_
}
func (this_ Parser_FunctionFlavor) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Parser_FunctionFlavor[this_-1].Name_
}
func (this_ Parser_FunctionFlavor) Equals(other_ any) bool {
	o_, ok_ := other_.(Parser_FunctionFlavor)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Parser_FunctionFlavor) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Parser_FunctionFlavor[this_-1].HashCode()
}
func (this_ Parser_FunctionFlavor) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_FunctionFlavor[this_-1].CompareTo(other_)
}
func (this_ Parser_FunctionFlavor) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Parser_FunctionFlavor[this_-1].Name_)
}
func init() {
	vals_ := []any{Parser_FunctionFlavor(1), Parser_FunctionFlavor(2), Parser_FunctionFlavor(3), Parser_FunctionFlavor(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser.FunctionFlavor", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser$FunctionFlavor", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Parser.FunctionFlavor", "com.google.javascript.jscomp.parsing.parser.Parser$FunctionFlavor", func(v_ any) bool {
		return j2gort.InstanceOf[Parser_FunctionFlavor](v_)
	})
}

// Java class: com.google.javascript.jscomp.parsing.parser.Parser.Config
type Parser_Config struct {
	idHash_            atomic.Int32
	atLeast6           bool `gson:"atLeast6"`
	atLeast8           bool `gson:"atLeast8"`
	isStrictMode       bool `gson:"isStrictMode"`
	warnTrailingCommas bool `gson:"warnTrailingCommas"`
}

func (this_ *Parser_Config) BindSelf_(s any) {
	_ = s
}
func (this_ *Parser_Config) IsJavaNil_() bool { return this_ == nil }
func (*Parser_Config) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Parser$Config", "Config"
}
func (this_ *Parser_Config) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Parser_Config) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Parser.Config", "com.google.javascript.jscomp.parsing.parser.Parser$Config", func(v_ any) bool {
		return j2gort.InstanceOf[*Parser_Config](v_)
	})
}

func (this_ *Parser_Config) initFields_Parser_Config() {
	_ = this_
}

// Java constructor: Config()
func NewParser_Config() *Parser_Config {
	this_ := &Parser_Config{}
	InitParser_Config(this_)
	return this_
}

func InitParser_Config(this_ *Parser_Config) {
	InitParser_Config_2(this_, StaticGet_Parser_Config_Mode_ES8_OR_GREATER(), false)
}

// Java constructor: Config(com.google.javascript.jscomp.parsing.parser.Parser.Config.Mode mode, boolean isStrictMode)
func NewParser_Config_2(mode Parser_Config_Mode, isStrictMode bool) *Parser_Config {
	this_ := &Parser_Config{}
	InitParser_Config_2(this_, mode, isStrictMode)
	return this_
}

func InitParser_Config_2(this_ *Parser_Config, mode Parser_Config_Mode, isStrictMode bool) {
	_ = mode
	_ = isStrictMode
	this_.initFields_Parser_Config()
	this_.atLeast6 = !((mode == StaticGet_Parser_Config_Mode_ES3()) || (mode == StaticGet_Parser_Config_Mode_ES5()))
	this_.atLeast8 = (mode == StaticGet_Parser_Config_Mode_ES8_OR_GREATER())
	this_.isStrictMode = isStrictMode
	this_.warnTrailingCommas = (mode == StaticGet_Parser_Config_Mode_ES3())
}

// Java enum: com.google.javascript.jscomp.parsing.parser.Parser.Config.Mode
type Parser_Config_Mode uint8

const (
	enum_Parser_Config_Mode_ES3            Parser_Config_Mode = 1
	enum_Parser_Config_Mode_ES5            Parser_Config_Mode = 2
	enum_Parser_Config_Mode_ES6_OR_ES7     Parser_Config_Mode = 3
	enum_Parser_Config_Mode_ES8_OR_GREATER Parser_Config_Mode = 4
)

var enumBaseTable_Parser_Config_Mode = [...]j2gort.EnumBase{
	{Name_: "ES3", Ordinal_: 0},
	{Name_: "ES5", Ordinal_: 1},
	{Name_: "ES6_OR_ES7", Ordinal_: 2},
	{Name_: "ES8_OR_GREATER", Ordinal_: 3},
}
var enumValues_Parser_Config_Mode = []Parser_Config_Mode{1, 2, 3, 4}

func StaticGet_Parser_Config_Mode_ES3() Parser_Config_Mode {
	return enum_Parser_Config_Mode_ES3
}
func StaticGet_Parser_Config_Mode_ES5() Parser_Config_Mode {
	return enum_Parser_Config_Mode_ES5
}
func StaticGet_Parser_Config_Mode_ES6_OR_ES7() Parser_Config_Mode {
	return enum_Parser_Config_Mode_ES6_OR_ES7
}
func StaticGet_Parser_Config_Mode_ES8_OR_GREATER() Parser_Config_Mode {
	return enum_Parser_Config_Mode_ES8_OR_GREATER
}

func Parser_Config_Mode_Values() []Parser_Config_Mode {
	return j2gort.CloneSlice(enumValues_Parser_Config_Mode)
}

func Parser_Config_Mode_ValueOf(name string) Parser_Config_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Parser_Config_Mode {
		if enumBaseTable_Parser_Config_Mode[i_].Name_ == name {
			return Parser_Config_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.Parser.Config.Mode." + name))
}

func (this_ Parser_Config_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (Parser_Config_Mode) HasIdentityEquals_()                                 {}
func (Parser_Config_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Parser_Config_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Parser_Config_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Parser_Config_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Parser$Config$Mode", "Mode"
}
func (this_ Parser_Config_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Parser_Config_Mode[this_-1]
}
func (this_ Parser_Config_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Parser_Config_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Parser_Config_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_Config_Mode[this_-1].Name_
}
func (this_ Parser_Config_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Parser_Config_Mode[this_-1].Name_
}
func (this_ Parser_Config_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(Parser_Config_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Parser_Config_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Parser_Config_Mode[this_-1].HashCode()
}
func (this_ Parser_Config_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_Config_Mode[this_-1].CompareTo(other_)
}
func (this_ Parser_Config_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Parser_Config_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{Parser_Config_Mode(1), Parser_Config_Mode(2), Parser_Config_Mode(3), Parser_Config_Mode(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser.Config.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser$Config$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Parser.Config.Mode", "com.google.javascript.jscomp.parsing.parser.Parser$Config$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[Parser_Config_Mode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.parser.Parser.Expression
type Parser_Expression uint8

const (
	enum_Parser_Expression_NO_IN  Parser_Expression = 1
	enum_Parser_Expression_NORMAL Parser_Expression = 2
)

var enumBaseTable_Parser_Expression = [...]j2gort.EnumBase{
	{Name_: "NO_IN", Ordinal_: 0},
	{Name_: "NORMAL", Ordinal_: 1},
}
var enumValues_Parser_Expression = []Parser_Expression{1, 2}

func StaticGet_Parser_Expression_NO_IN() Parser_Expression {
	return enum_Parser_Expression_NO_IN
}
func StaticGet_Parser_Expression_NORMAL() Parser_Expression {
	return enum_Parser_Expression_NORMAL
}

func Parser_Expression_Values() []Parser_Expression {
	return j2gort.CloneSlice(enumValues_Parser_Expression)
}

func Parser_Expression_ValueOf(name string) Parser_Expression {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Parser_Expression {
		if enumBaseTable_Parser_Expression[i_].Name_ == name {
			return Parser_Expression(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.Parser.Expression." + name))
}

func (this_ Parser_Expression) IsJavaNil_() bool                              { return this_ == 0 }
func (Parser_Expression) HasIdentityEquals_()                                 {}
func (Parser_Expression) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Parser_Expression) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Parser_Expression) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Parser_Expression) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Parser$Expression", "Expression"
}
func (this_ Parser_Expression) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Parser_Expression[this_-1]
}
func (this_ Parser_Expression) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Parser_Expression) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Parser_Expression) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_Expression[this_-1].Name_
}
func (this_ Parser_Expression) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Parser_Expression[this_-1].Name_
}
func (this_ Parser_Expression) Equals(other_ any) bool {
	o_, ok_ := other_.(Parser_Expression)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Parser_Expression) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Parser_Expression[this_-1].HashCode()
}
func (this_ Parser_Expression) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_Expression[this_-1].CompareTo(other_)
}
func (this_ Parser_Expression) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Parser_Expression[this_-1].Name_)
}
func init() {
	vals_ := []any{Parser_Expression(1), Parser_Expression(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser.Expression", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser$Expression", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Parser.Expression", "com.google.javascript.jscomp.parsing.parser.Parser$Expression", func(v_ any) bool {
		return j2gort.InstanceOf[Parser_Expression](v_)
	})
}

// Java enum: com.google.javascript.jscomp.parsing.parser.Parser.PatternKind
type Parser_PatternKind uint8

const (
	enum_Parser_PatternKind_INITIALIZER Parser_PatternKind = 1
	enum_Parser_PatternKind_ANY         Parser_PatternKind = 2
)

var enumBaseTable_Parser_PatternKind = [...]j2gort.EnumBase{
	{Name_: "INITIALIZER", Ordinal_: 0},
	{Name_: "ANY", Ordinal_: 1},
}
var enumValues_Parser_PatternKind = []Parser_PatternKind{1, 2}

func StaticGet_Parser_PatternKind_INITIALIZER() Parser_PatternKind {
	return enum_Parser_PatternKind_INITIALIZER
}
func StaticGet_Parser_PatternKind_ANY() Parser_PatternKind {
	return enum_Parser_PatternKind_ANY
}

func Parser_PatternKind_Values() []Parser_PatternKind {
	return j2gort.CloneSlice(enumValues_Parser_PatternKind)
}

func Parser_PatternKind_ValueOf(name string) Parser_PatternKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Parser_PatternKind {
		if enumBaseTable_Parser_PatternKind[i_].Name_ == name {
			return Parser_PatternKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.Parser.PatternKind." + name))
}

func (this_ Parser_PatternKind) IsJavaNil_() bool                              { return this_ == 0 }
func (Parser_PatternKind) HasIdentityEquals_()                                 {}
func (Parser_PatternKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Parser_PatternKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Parser_PatternKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Parser_PatternKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.Parser$PatternKind", "PatternKind"
}
func (this_ Parser_PatternKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Parser_PatternKind[this_-1]
}
func (this_ Parser_PatternKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Parser_PatternKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Parser_PatternKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_PatternKind[this_-1].Name_
}
func (this_ Parser_PatternKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Parser_PatternKind[this_-1].Name_
}
func (this_ Parser_PatternKind) Equals(other_ any) bool {
	o_, ok_ := other_.(Parser_PatternKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Parser_PatternKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Parser_PatternKind[this_-1].HashCode()
}
func (this_ Parser_PatternKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Parser_PatternKind[this_-1].CompareTo(other_)
}
func (this_ Parser_PatternKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Parser_PatternKind[this_-1].Name_)
}
func init() {
	vals_ := []any{Parser_PatternKind(1), Parser_PatternKind(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser.PatternKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.Parser$PatternKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Parser.PatternKind", "com.google.javascript.jscomp.parsing.parser.Parser$PatternKind", func(v_ any) bool {
		return j2gort.InstanceOf[Parser_PatternKind](v_)
	})
}
func (this_ *Parser) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_ImportDeclarationTree.Release()
	this_.slab_ImportSpecifierTree.Release()
	this_.slab_ExportDeclarationTree.Release()
	this_.slab_ExportSpecifierTree.Release()
	this_.slab_ClassDeclarationTree.Release()
	this_.slab_FieldDeclarationTree.Release()
	this_.slab_ComputedPropertyFieldTree.Release()
	this_.slab_FunctionDeclarationTree_Builder.Release()
	this_.slab_FunctionDeclarationTree.Release()
	this_.slab_ComputedPropertyMethodTree.Release()
	this_.slab_DefaultParameterTree.Release()
	this_.slab_FormalParameterListTree.Release()
	this_.slab_BlockTree.Release()
	this_.slab_IterSpreadTree.Release()
	this_.slab_VariableStatementTree.Release()
	this_.slab_VariableDeclarationListTree.Release()
	this_.slab_VariableDeclarationTree.Release()
	this_.slab_EmptyStatementTree.Release()
	this_.slab_ExpressionStatementTree.Release()
	this_.slab_IfStatementTree.Release()
	this_.slab_DoWhileStatementTree.Release()
	this_.slab_WhileStatementTree.Release()
	this_.slab_ForOfStatementTree.Release()
	this_.slab_ForAwaitOfStatementTree.Release()
	this_.slab_NullTree.Release()
	this_.slab_ForStatementTree.Release()
	this_.slab_ForInStatementTree.Release()
	this_.slab_ContinueStatementTree.Release()
	this_.slab_BreakStatementTree.Release()
	this_.slab_ReturnStatementTree.Release()
	this_.slab_WithStatementTree.Release()
	this_.slab_SwitchStatementTree.Release()
	this_.slab_CaseClauseTree.Release()
	this_.slab_DefaultClauseTree.Release()
	this_.slab_LabelledStatementTree.Release()
	this_.slab_ThrowStatementTree.Release()
	this_.slab_TryStatementTree.Release()
	this_.slab_CatchTree.Release()
	this_.slab_FinallyTree.Release()
	this_.slab_DebuggerStatementTree.Release()
	this_.slab_SuperExpressionTree.Release()
	this_.slab_ThisExpressionTree.Release()
	this_.slab_DynamicImportTree.Release()
	this_.slab_IdentifierExpressionTree.Release()
	this_.slab_LiteralExpressionTree.Release()
	this_.slab_TemplateLiteralPortionTree.Release()
	this_.slab_TemplateLiteralExpressionTree.Release()
	this_.slab_TemplateSubstitutionTree.Release()
	this_.slab_ArrayLiteralExpressionTree.Release()
	this_.slab_ObjectLiteralExpressionTree.Release()
	this_.slab_ObjectSpreadTree.Release()
	this_.slab_ComputedPropertyDefinitionTree.Release()
	this_.slab_GetAccessorTree.Release()
	this_.slab_ComputedPropertyGetterTree.Release()
	this_.slab_SetAccessorTree.Release()
	this_.slab_ComputedPropertySetterTree.Release()
	this_.slab_PropertyNameAssignmentTree.Release()
	this_.slab_MissingPrimaryExpressionTree.Release()
	this_.slab_CommaExpressionTree.Release()
	this_.slab_ParenExpressionTree.Release()
	this_.slab_BinaryOperatorTree.Release()
	this_.slab_YieldExpressionTree.Release()
	this_.slab_ConditionalExpressionTree.Release()
	this_.slab_UnaryExpressionTree.Release()
	this_.slab_AwaitExpressionTree.Release()
	this_.slab_UpdateExpressionTree.Release()
	this_.slab_CallExpressionTree.Release()
	this_.slab_MemberLookupExpressionTree.Release()
	this_.slab_MemberExpressionTree.Release()
	this_.slab_OptChainCallExpressionTree.Release()
	this_.slab_OptionalMemberLookupExpressionTree.Release()
	this_.slab_OptionalMemberExpressionTree.Release()
	this_.slab_NewExpressionTree.Release()
	this_.slab_NewTargetExpressionTree.Release()
	this_.slab_ImportMetaExpressionTree.Release()
	this_.slab_ArgumentListTree.Release()
	this_.slab_IterRestTree.Release()
	this_.slab_ArrayPatternTree.Release()
	this_.slab_ObjectRestTree.Release()
	this_.slab_ObjectPatternTree.Release()
}
func (this_ *Parser) ResetSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_ImportDeclarationTree.Reset()
	this_.slab_ImportSpecifierTree.Reset()
	this_.slab_ExportDeclarationTree.Reset()
	this_.slab_ExportSpecifierTree.Reset()
	this_.slab_ClassDeclarationTree.Reset()
	this_.slab_FieldDeclarationTree.Reset()
	this_.slab_ComputedPropertyFieldTree.Reset()
	this_.slab_FunctionDeclarationTree_Builder.Reset()
	this_.slab_FunctionDeclarationTree.Reset()
	this_.slab_ComputedPropertyMethodTree.Reset()
	this_.slab_DefaultParameterTree.Reset()
	this_.slab_FormalParameterListTree.Reset()
	this_.slab_BlockTree.Reset()
	this_.slab_IterSpreadTree.Reset()
	this_.slab_VariableStatementTree.Reset()
	this_.slab_VariableDeclarationListTree.Reset()
	this_.slab_VariableDeclarationTree.Reset()
	this_.slab_EmptyStatementTree.Reset()
	this_.slab_ExpressionStatementTree.Reset()
	this_.slab_IfStatementTree.Reset()
	this_.slab_DoWhileStatementTree.Reset()
	this_.slab_WhileStatementTree.Reset()
	this_.slab_ForOfStatementTree.Reset()
	this_.slab_ForAwaitOfStatementTree.Reset()
	this_.slab_NullTree.Reset()
	this_.slab_ForStatementTree.Reset()
	this_.slab_ForInStatementTree.Reset()
	this_.slab_ContinueStatementTree.Reset()
	this_.slab_BreakStatementTree.Reset()
	this_.slab_ReturnStatementTree.Reset()
	this_.slab_WithStatementTree.Reset()
	this_.slab_SwitchStatementTree.Reset()
	this_.slab_CaseClauseTree.Reset()
	this_.slab_DefaultClauseTree.Reset()
	this_.slab_LabelledStatementTree.Reset()
	this_.slab_ThrowStatementTree.Reset()
	this_.slab_TryStatementTree.Reset()
	this_.slab_CatchTree.Reset()
	this_.slab_FinallyTree.Reset()
	this_.slab_DebuggerStatementTree.Reset()
	this_.slab_SuperExpressionTree.Reset()
	this_.slab_ThisExpressionTree.Reset()
	this_.slab_DynamicImportTree.Reset()
	this_.slab_IdentifierExpressionTree.Reset()
	this_.slab_LiteralExpressionTree.Reset()
	this_.slab_TemplateLiteralPortionTree.Reset()
	this_.slab_TemplateLiteralExpressionTree.Reset()
	this_.slab_TemplateSubstitutionTree.Reset()
	this_.slab_ArrayLiteralExpressionTree.Reset()
	this_.slab_ObjectLiteralExpressionTree.Reset()
	this_.slab_ObjectSpreadTree.Reset()
	this_.slab_ComputedPropertyDefinitionTree.Reset()
	this_.slab_GetAccessorTree.Reset()
	this_.slab_ComputedPropertyGetterTree.Reset()
	this_.slab_SetAccessorTree.Reset()
	this_.slab_ComputedPropertySetterTree.Reset()
	this_.slab_PropertyNameAssignmentTree.Reset()
	this_.slab_MissingPrimaryExpressionTree.Reset()
	this_.slab_CommaExpressionTree.Reset()
	this_.slab_ParenExpressionTree.Reset()
	this_.slab_BinaryOperatorTree.Reset()
	this_.slab_YieldExpressionTree.Reset()
	this_.slab_ConditionalExpressionTree.Reset()
	this_.slab_UnaryExpressionTree.Reset()
	this_.slab_AwaitExpressionTree.Reset()
	this_.slab_UpdateExpressionTree.Reset()
	this_.slab_CallExpressionTree.Reset()
	this_.slab_MemberLookupExpressionTree.Reset()
	this_.slab_MemberExpressionTree.Reset()
	this_.slab_OptChainCallExpressionTree.Reset()
	this_.slab_OptionalMemberLookupExpressionTree.Reset()
	this_.slab_OptionalMemberExpressionTree.Reset()
	this_.slab_NewExpressionTree.Reset()
	this_.slab_NewTargetExpressionTree.Reset()
	this_.slab_ImportMetaExpressionTree.Reset()
	this_.slab_ArgumentListTree.Reset()
	this_.slab_IterRestTree.Reset()
	this_.slab_ArrayPatternTree.Reset()
	this_.slab_ObjectRestTree.Reset()
	this_.slab_ObjectPatternTree.Reset()
}
func (this_ *Parser) ClearForPool_() {
	if this_ == nil {
		return
	}
	this_.ResetSlabs_()
	this_.errorReporter = nil
	this_.config = nil
	this_.comments = nil
	this_.lastCommentEndPosition = SourcePosition{}
	this_.features = nil
	this_.lastSourcePosition = SourcePosition{}
	this_.sourceMapURL = j2gort.NullString
	if len(this_.treeScratch) > 256 {
		this_.treeScratch = make([]ParseTree, 32)
	} else {
		clear(this_.treeScratch)
	}
	if len(this_.commaScratch) > 256 {
		this_.commaScratch = make([]SourcePosition, 32)
	} else {
		clear(this_.commaScratch)
	}
	this_.treeScratchTop = 0
	this_.commaScratchTop = 0
	this_.idHash_.Store(0)
	this_.inGenerator = false
	if this_.scanner != nil {
		sc := this_.scanner
		sc.errorReporter = nil
		sc.source = nil
		sc.contents = j2gort.NullString
		if cap(sc.currentTokens) > 64 {
			sc.currentTokens = nil
		} else {
			clear(sc.currentTokens)
			sc.currentTokens = sc.currentTokens[:0]
		}
		sc.index = 0
		sc.typeParameterLevel = 0
		sc.idHash_.Store(0)
		if sc.lineNumberScanner != nil {
			lns := sc.lineNumberScanner
			lns.sourceFile = nil
			lns.contents = j2gort.NullString
			lns.cachedPosition = SourcePosition{}
			lns.idHash_.Store(0)
		}
	}
}

func (this_ *Parser) ResetForReuse_(config *Parser_Config, errorReporter ErrorReporter, source *SourceFile) {
	clinit_Parser()
	this_.ResetSlabs_()
	this_.features = StaticGet_FeatureSet_BARE_MINIMUM()
	this_.sourceMapURL = j2gort.NullString
	if len(this_.treeScratch) < 32 {
		this_.treeScratch = make([]ParseTree, 32)
	}
	if len(this_.commaScratch) < 32 {
		this_.commaScratch = make([]SourcePosition, 32)
	}
	this_.treeScratchTop = 0
	this_.commaScratchTop = 0
	this_.comments = nil
	this_.lastCommentEndPosition = SourcePosition{}
	this_.idHash_.Store(0)
	this_.inGenerator = false
	this_.config = config
	this_.errorReporter = errorReporter
	if this_.scanner == nil {
		this_.scanner = NewScanner(errorReporter, this_, source, int32(0))
	} else {
		sc := this_.scanner
		sc.errorReporter = errorReporter
		sc.commentRecorder = this_
		sc.source = source
		if sc.lineNumberScanner == nil {
			sc.lineNumberScanner = NewLineNumberScanner(source)
		} else {
			lns := sc.lineNumberScanner
			lns.sourceFile = source
			lns.contents = source.contents
			lns.cachedPosition = SourcePosition{}
			lns.sourceLength = j2gort.JStrLen(lns.contents)
			lns.lastLine = -1
			lns.lastLineStart = -1
			lns.nextLineStart = 0
			lns.cachedOffset = -1
			lns.idHash_.Store(0)
		}
		sc.contents = source.contents
		sc.contentsLength = j2gort.JStrLen(sc.contents)
		clear(sc.currentTokens)
		sc.currentTokens = sc.currentTokens[:0]
		sc.index = 0
		sc.typeParameterLevel = 0
		sc.idHash_.Store(0)
	}
	this_.lastSourcePosition = this_.scanner.GetPosition_0()
}
