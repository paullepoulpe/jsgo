// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20261015

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

package parser

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Parser parses a javascript file.
//
// The various parseX() methods never return null - even when parse errors
// are encountered. Typically parseX() will return a XTree ParseTree. Each
// ParseTree that is created includes its source location. The typical
// pattern for a parseX() method is:
//
//	XTree parseX() {
//	  SourcePosition start = getTreeStartLocation();
//	  parse X grammar element and its children
//	  return new XTree(getTreeLocation(start), children);
//	}
//
// parseX() methods must consume at least 1 token - even in error cases.
// This prevents infinite loops in the parser.
//
// Many parseX() methods are matched by a 'boolean peekX()' method which will
// return true if the beginning of an X appears at the current location.
// There are also peek() methods which examine the next token. peek() methods
// must not consume any tokens.
//
// The eat() method consumes a token and reports an error if the consumed
// token is not of the expected type. The eatOpt() methods consume the next
// token iff the next token is of the expected type and return the consumed
// token or null if no token was consumed.
//
// When parse errors are encountered, an error should be reported and the
// parse should return a best guess at the current parse tree.
//
// When parsing lists, the preferred pattern is:
//
//	eat(LIST_START);
//	ImmutableList.Builder<ParseTree> elements = ImmutableList.builder();
//	while (peekListElement()) {
//	  elements.add(parseListElement());
//	}
//	eat(LIST_END);
type Parser struct {
	scanner              *Scanner
	errorReporter        ErrorReporter
	config               *Parser_Config
	commentRecorder      *Parser_CommentRecorder
	functionContextStack []Parser_FunctionFlavor
	functionContextBuf   [8]Parser_FunctionFlavor
	features             *FeatureSet
	lastSourcePosition   SourcePosition
	// Go-only: language-adapter — nil is Java null.
	sourceMapURL *string

	// Go-only: performance — scratch stacks for the lists being parsed (Java
	// uses an ImmutableList.Builder per list); popStack copies each finished
	// list out at its exact length.
	treeStack    []ParseTree
	posStack     []SourcePosition
	varDeclStack []*VariableDeclarationTree
}

// popStack returns the elements pushed on *stack since base as a new slice
// of exactly that length (never nil), and pops them.
func popStack[T any](stack *[]T, base int) []T {
	s := (*stack)[base:]
	res := make([]T, len(s))
	copy(res, s)
	clear(s)
	*stack = (*stack)[:base]
	return res
}

// Parser_FunctionFlavor is Parser.FunctionFlavor: indicates the type of
// function currently being parsed.
type Parser_FunctionFlavor int32

const (
	Parser_FunctionFlavor_NORMAL Parser_FunctionFlavor = iota
	Parser_FunctionFlavor_GENERATOR
	Parser_FunctionFlavor_ASYNCHRONOUS
	Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR
)

var parser_FunctionFlavor_names = [...]string{"NORMAL", "GENERATOR", "ASYNCHRONOUS", "ASYNCHRONOUS_GENERATOR"}

// parser_FunctionFlavor_isGenerator is the Java field FunctionFlavor.isGenerator.
var parser_FunctionFlavor_isGenerator = [...]bool{false, true, false, true}

// Parser_FunctionFlavor_values is FunctionFlavor.values().
//
// Java: Parser.FunctionFlavor.values()
func Parser_FunctionFlavor_values() []Parser_FunctionFlavor {
	return []Parser_FunctionFlavor{
		Parser_FunctionFlavor_NORMAL,
		Parser_FunctionFlavor_GENERATOR,
		Parser_FunctionFlavor_ASYNCHRONOUS,
		Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR,
	}
}

// Name is Enum.name().
//
// Java: Parser.FunctionFlavor.name()
func (f Parser_FunctionFlavor) Name() string { return parser_FunctionFlavor_names[f] }

// String is Enum.toString().
//
// Java: Parser.FunctionFlavor.toString()
func (f Parser_FunctionFlavor) String() string { return f.Name() }

// IsGenerator is the Java field FunctionFlavor.isGenerator.
//
// Java: field Parser.FunctionFlavor.isGenerator
func (f Parser_FunctionFlavor) IsGenerator() bool { return parser_FunctionFlavor_isGenerator[f] }

// NewParser constructs a Parser.
//
// Java: Parser(Config, ErrorReporter, SourceFile)
func NewParser(config *Parser_Config, errorReporter ErrorReporter, source *SourceFile) *Parser {
	p := &Parser{}
	p.commentRecorder = NewParser_CommentRecorder(p)
	p.functionContextStack = p.functionContextBuf[:0]
	p.features = FeatureSet_BARE_MINIMUM
	p.config = config
	p.errorReporter = errorReporter
	p.scanner = NewScanner(errorReporter, p.commentRecorder, source, 0)
	p.functionContextStack = append(p.functionContextStack, Parser_FunctionFlavor_NORMAL)
	p.lastSourcePosition = p.scanner.GetPosition()
	return p
}

// Parser_Config is Parser.Config.
type Parser_Config struct {
	atLeast6           bool
	atLeast8           bool
	isStrictMode       bool
	warnTrailingCommas bool
}

// Parser_Config_Mode is Parser.Config.Mode.
type Parser_Config_Mode int32

const (
	Parser_Config_Mode_ES3 Parser_Config_Mode = iota
	Parser_Config_Mode_ES5
	Parser_Config_Mode_ES6_OR_ES7
	Parser_Config_Mode_ES8_OR_GREATER
)

var parser_Config_Mode_names = [...]string{"ES3", "ES5", "ES6_OR_ES7", "ES8_OR_GREATER"}

// Parser_Config_Mode_values is Mode.values().
//
// Java: Parser.Config.Mode.values()
func Parser_Config_Mode_values() []Parser_Config_Mode {
	return []Parser_Config_Mode{
		Parser_Config_Mode_ES3,
		Parser_Config_Mode_ES5,
		Parser_Config_Mode_ES6_OR_ES7,
		Parser_Config_Mode_ES8_OR_GREATER,
	}
}

// Name is Enum.name().
//
// Java: Parser.Config.Mode.name()
func (m Parser_Config_Mode) Name() string { return parser_Config_Mode_names[m] }

// String is Enum.toString().
//
// Java: Parser.Config.Mode.toString()
func (m Parser_Config_Mode) String() string { return m.Name() }

// NewParser_Config constructs the default Config (ES8_OR_GREATER, sloppy).
//
// Java: Parser.Config()
func NewParser_Config() *Parser_Config {
	return NewParser_ConfigWithMode(Parser_Config_Mode_ES8_OR_GREATER /* isStrictMode= */, false)
}

// NewParser_ConfigWithMode constructs a Config.
//
// Java: Parser.Config(Mode, boolean)
func NewParser_ConfigWithMode(mode Parser_Config_Mode, isStrictMode bool) *Parser_Config {
	c := &Parser_Config{}
	c.atLeast6 = !(mode == Parser_Config_Mode_ES3 || mode == Parser_Config_Mode_ES5)
	c.atLeast8 = mode == Parser_Config_Mode_ES8_OR_GREATER
	c.isStrictMode = isStrictMode

	// Generally, we allow everything that is valid in any mode
	// we only warn about things that are not represented in the AST.
	c.warnTrailingCommas = mode == Parser_Config_Mode_ES3
	return c
}

// Parser_SOURCE_MAPPING_URL_PREFIX is Parser.SOURCE_MAPPING_URL_PREFIX.
const Parser_SOURCE_MAPPING_URL_PREFIX = "//# sourceMappingURL="

// Parser_CommentRecorder is the inner class Parser.CommentRecorder.
type Parser_CommentRecorder struct {
	// Go-only: language-adapter — the enclosing Parser instance of the Java
	// inner class.
	this *Parser

	comments []*Comment
	// Go-only: language-adapter — the offset of Java's
	// lastCommentEndPosition, -1 for null.
	lastCommentEndOffset int32
}

// NewParser_CommentRecorder constructs the comment recorder of parser.
//
// Java: Parser.CommentRecorder()
func NewParser_CommentRecorder(this *Parser) *Parser_CommentRecorder {
	return &Parser_CommentRecorder{this: this, comments: []*Comment{}, lastCommentEndOffset: -1}
}

// RecordComment records a comment the scanner passed.
//
// Java: Parser.CommentRecorder.recordComment(Comment.Type, SourceRange, String)
func (c *Parser_CommentRecorder) RecordComment(type_ Comment_Type, range_ SourceRange, value string) {
	// If we rewind the token stream, the scanner might pass comments that we've already seen.
	// Only record comments past the furthest comment end position we've seen.
	// NB: this assumes the CommentRecorder is used for at most one source file.
	if range_.End.Offset > c.lastCommentEndOffset {
		value = jutil.StringTrim(value)
		if strings.HasPrefix(value, Parser_SOURCE_MAPPING_URL_PREFIX) {
			sourceMapURL := value[len(Parser_SOURCE_MAPPING_URL_PREFIX):]
			c.this.sourceMapURL = &sourceMapURL
		}
		c.comments = append(c.comments, NewComment(value, range_, type_))
		c.lastCommentEndOffset = range_.End.Offset
	}
}

// GetComments builds the list of recorded comments.
//
// Java: Parser.CommentRecorder.getComments()
func (c *Parser_CommentRecorder) GetComments() []*Comment {
	return append([]*Comment{}, c.comments...)
}

// GetComments returns the comments of the source file.
//
// Java: Parser.getComments()
func (p *Parser) GetComments() []*Comment {
	return p.commentRecorder.GetComments()
}

// GetFeatures returns the language features used by the source file.
//
// Java: Parser.getFeatures()
func (p *Parser) GetFeatures() *FeatureSet {
	return p.features
}

// GetSourceMapURL returns the url provided by the sourceMappingURL if any
// was found (nil otherwise).
//
// Java: Parser.getSourceMapURL()
func (p *Parser) GetSourceMapURL() *string {
	return p.sourceMapURL
}

// 14 Program

// ParseProgram parses the whole source file.
//
// Java: Parser.parseProgram()
func (p *Parser) ParseProgram() *ProgramTree {
	// Go-only: language-adapter — Java catches the StackOverflowError of
	// too-deep recursion here (`catch (Error e)`, reporting "Too deep
	// recursion while parsing" and returning null) and rethrows any other
	// Error. Go cannot recover from exhausting the goroutine stack (it is a
	// fatal runtime error, and the default 1 GB limit is far beyond the
	// JVM's), so only the rethrow path exists here.

	// Set the start location at the beginning of the file rather than the beginning of the first
	// token.  This ensures that it accounts for leading comments.
	start := p.lastSourcePosition
	sourceElements := p.ParseGlobalSourceElements()
	p.Eat(TokenType_END_OF_FILE)
	return NewProgramTree(p.GetTreeLocation(start), sourceElements, p.commentRecorder.GetComments())
}

// Java: Parser.parseGlobalSourceElements()
func (p *Parser) ParseGlobalSourceElements() []ParseTree {
	base := len(p.treeStack)

	for !p.Peek(TokenType_END_OF_FILE) {
		elem := p.ParseScriptElement()
		p.treeStack = append(p.treeStack, elem)
	}

	return popStack(&p.treeStack, base)
}

// ParseScriptElement parses an ImportDeclaration, ExportDeclaration or
// SourceElement.
//
// Java: Parser.parseScriptElement()
func (p *Parser) ParseScriptElement() ParseTree {
	if p.PeekImportDeclaration() {
		return p.ParseImportDeclaration()
	}

	if p.PeekExportDeclaration() {
		return p.ParseExportDeclaration()
	}

	return p.ParseSourceElement()
}

// Java: Parser.peekImportDeclaration()
func (p *Parser) PeekImportDeclaration() bool {
	return p.Peek(TokenType_IMPORT) &&
		(p.PeekIdOrKeywordAt(1) ||
			p.PeekAt(1, TokenType_STRING) ||
			p.PeekAt(1, TokenType_OPEN_CURLY) ||
			p.PeekAt(1, TokenType_STAR))
}

// Java: Parser.parseImportDeclaration()
func (p *Parser) ParseImportDeclaration() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IMPORT)

	// import ModuleSpecifier ;
	if p.Peek(TokenType_STRING) {
		moduleSpecifier := p.Eat(TokenType_STRING).(LiteralToken)
		p.EatPossiblyImplicitSemiColon()

		return NewImportDeclarationTree(p.GetTreeLocation(start), nil, nil, nil, moduleSpecifier)
	}

	// import ImportedDefaultBinding from ModuleSpecifier
	// import NameSpaceImport from ModuleSpecifier
	// import NamedImports from ModuleSpecifier ;
	// import ImportedDefaultBinding , NameSpaceImport from ModuleSpecifier ;
	// import ImportedDefaultBinding , NamedImports from ModuleSpecifier ;
	var defaultBindingIdentifier *IdentifierToken
	var nameSpaceImportIdentifier *IdentifierToken
	var identifierSet []ParseTree

	parseExplicitNames := true
	if p.PeekId() {
		defaultBindingIdentifier = p.EatId()
		if p.Peek(TokenType_COMMA) {
			p.Eat(TokenType_COMMA)
		} else {
			parseExplicitNames = false
		}
	} else if Keywords_IsKeywordTokenType(p.PeekType()) {
		keyword := p.NextToken()
		p.ReportErrorAtToken(keyword, "cannot use keyword '%s' here.", keyword)
	}

	if parseExplicitNames {
		if p.Peek(TokenType_STAR) {
			p.Eat(TokenType_STAR)
			p.EatPredefinedString(PredefinedName_AS)
			nameSpaceImportIdentifier = p.EatId()
		} else {
			identifierSet = p.ParseImportSpecifierSet()
		}
	}

	p.EatPredefinedString(PredefinedName_FROM)
	moduleStr := p.Eat(TokenType_STRING)
	var moduleSpecifier LiteralToken
	if moduleStr != nil {
		moduleSpecifier = moduleStr.(LiteralToken)
	}
	p.EatPossiblyImplicitSemiColon()

	return NewImportDeclarationTree(
		p.GetTreeLocation(start),
		defaultBindingIdentifier,
		identifierSet,
		nameSpaceImportIdentifier,
		moduleSpecifier)
}

// ParseImportSpecifierSet parses
// ImportSpecifierSet ::= '{' (ImportSpecifier (',' ImportSpecifier)* (,)? )?  '}'
//
// Java: Parser.parseImportSpecifierSet()
func (p *Parser) ParseImportSpecifierSet() []ParseTree {
	elements := []ParseTree{}
	p.Eat(TokenType_OPEN_CURLY)
	for p.PeekIdOrKeyword() {
		elements = append(elements, p.ParseImportSpecifier())
		if !p.Peek(TokenType_CLOSE_CURLY) {
			p.Eat(TokenType_COMMA)
		}
	}
	p.Eat(TokenType_CLOSE_CURLY)
	return elements
}

// ParseImportSpecifier parses ImportSpecifier ::= Identifier ('as' Identifier)?
//
// Java: Parser.parseImportSpecifier()
func (p *Parser) ParseImportSpecifier() ParseTree {
	start := p.GetTreeStartLocation()
	importedName := p.EatIdOrKeywordAsId()
	var destinationName *IdentifierToken
	if p.PeekPredefinedString(PredefinedName_AS) {
		p.EatPredefinedString(PredefinedName_AS)
		destinationName = p.EatId()
	} else if importedName.IsKeyword() {
		p.ReportExpectedError(nil, PredefinedName_AS)
	}
	return NewImportSpecifierTree(p.GetTreeLocation(start), importedName, destinationName)
}

// export  VariableStatement
// export  FunctionDeclaration
// export  ConstStatement
// export  ClassDeclaration
// export  default expression
// etc

// Java: Parser.peekExportDeclaration()
func (p *Parser) PeekExportDeclaration() bool {
	return p.Peek(TokenType_EXPORT)
}

/*
ExportDeclaration :

	export * FromClause ;
	export ExportClause [NoReference] FromClause ;
	export ExportClause ;
	export VariableStatement
	export Declaration[Default]
	export default AssignmentExpression ;

ExportClause [NoReference] :

	{ }
	{ ExportsList [?NoReference] }
	{ ExportsList [?NoReference] , }

ExportsList [NoReference] :

	ExportSpecifier [?NoReference]
	ExportsList [?NoReference] , ExportSpecifier [?NoReference]

ExportSpecifier [NoReference] :

	[~NoReference] IdentifierReference
	[~NoReference] IdentifierReference as IdentifierName
	[+NoReference] IdentifierName
	[+NoReference] IdentifierName as IdentifierName
*/

// Java: Parser.parseExportDeclaration()
func (p *Parser) ParseExportDeclaration() ParseTree {
	start := p.GetTreeStartLocation()
	isDefault := false
	isExportAll := false
	isExportSpecifier := false
	needsSemiColon := true
	p.Eat(TokenType_EXPORT)
	var export ParseTree
	var exportSpecifierList []ParseTree
	switch p.PeekType() {
	case TokenType_STAR:
		isExportAll = true
		p.NextToken()
	case TokenType_IDENTIFIER:
		export = p.ParseAsyncFunctionDeclaration()
	case TokenType_FUNCTION:
		export = p.ParseFunctionDeclaration()
		needsSemiColon = false
	case TokenType_CLASS:
		export = p.ParseClassDeclaration()
		needsSemiColon = false
	case TokenType_DEFAULT:
		isDefault = true
		p.NextToken()
		export = p.ParseExpression()
		needsSemiColon = false
	case TokenType_OPEN_CURLY:
		isExportSpecifier = true
		exportSpecifierList = p.ParseExportSpecifierSet()
	case TokenType_VAR, TokenType_LET, TokenType_CONST:
		fallthrough
	default: // unreachable, parse as a var decl to get a parse error.
		export = variableDeclarationListAsParseTree(p.ParseVariableDeclarationList())
	}

	var moduleSpecifier LiteralToken
	if isExportAll || (isExportSpecifier && p.PeekPredefinedString(PredefinedName_FROM)) {
		p.EatPredefinedString(PredefinedName_FROM)
		if token := p.Eat(TokenType_STRING); token != nil {
			moduleSpecifier = token.(LiteralToken)
		}
	} else if isExportSpecifier {
		for _, tree := range exportSpecifierList {
			importedName := tree.(*ExportSpecifierTree).ImportedName
			if importedName.IsKeyword() {
				p.ReportErrorAtToken(identifierTokenAsToken(importedName), "cannot use keyword '%s' here.", importedName)
			}
		}
	}

	if needsSemiColon || p.PeekImplicitSemiColon() {
		p.EatPossiblyImplicitSemiColon()
	}

	return NewExportDeclarationTree(
		p.GetTreeLocation(start),
		isDefault,
		isExportAll,
		export,
		exportSpecifierList,
		moduleSpecifier)
}

// ParseExportSpecifierSet parses
// ExportSpecifierSet ::= '{' (ExportSpecifier (',' ExportSpecifier)* (,)? )?  '}'
//
// Java: Parser.parseExportSpecifierSet()
func (p *Parser) ParseExportSpecifierSet() []ParseTree {
	elements := []ParseTree{}
	p.Eat(TokenType_OPEN_CURLY)
	for p.PeekIdOrKeyword() {
		elements = append(elements, p.ParseExportSpecifier())
		if !p.Peek(TokenType_CLOSE_CURLY) {
			p.Eat(TokenType_COMMA)
		}
	}
	p.Eat(TokenType_CLOSE_CURLY)
	return elements
}

// ParseExportSpecifier parses ExportSpecifier ::= Identifier ('as' Identifier)?
//
// Java: Parser.parseExportSpecifier()
func (p *Parser) ParseExportSpecifier() ParseTree {
	start := p.GetTreeStartLocation()
	importedName := p.EatIdOrKeywordAsId()
	var destinationName *IdentifierToken
	if p.PeekPredefinedString(PredefinedName_AS) {
		p.EatPredefinedString(PredefinedName_AS)
		destinationName = p.EatIdOrKeywordAsId()
	}
	return NewExportSpecifierTree(p.GetTreeLocation(start), importedName, destinationName)
}

// Java: Parser.peekClassDeclaration()
func (p *Parser) PeekClassDeclaration() bool {
	return p.Peek(TokenType_CLASS)
}

// Java: Parser.parseClassDeclaration()
func (p *Parser) ParseClassDeclaration() ParseTree {
	return p.ParseClass( /* isExpression= */ false)
}

// Java: Parser.parseClassExpression()
func (p *Parser) ParseClassExpression() ParseTree {
	return p.ParseClass( /* isExpression= */ true)
}

// Java: Parser.parseClass(boolean)
func (p *Parser) ParseClass(isExpression bool) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_CLASS)
	var name *IdentifierToken
	if !isExpression || p.PeekId() {
		name = p.EatId()
	}

	var superClass ParseTree
	if p.Peek(TokenType_EXTENDS) {
		p.Eat(TokenType_EXTENDS)
		superClass = p.ParseLeftHandSideExpression()
	}

	p.Eat(TokenType_OPEN_CURLY)
	elements := p.ParseClassElements()
	p.Eat(TokenType_CLOSE_CURLY)
	return NewClassDeclarationTree(p.GetTreeLocation(start), name, superClass, elements)
}

// Java: Parser.parseClassElements()
func (p *Parser) ParseClassElements() []ParseTree {
	result := []ParseTree{}
	for {
		token := p.PeekToken()
		if token.GetTokenBase().Type == TokenType_SEMI_COLON {
			p.Eat(TokenType_SEMI_COLON)
			continue
		} else {
			if p.IsClassElementStart(token) {
			} else {
				return result
			}
		}
		result = append(result, p.ParseClassElement())
	}
}

// Java: Parser.isClassElementStart(Token)
func (p *Parser) IsClassElementStart(token Token) bool {
	switch token.GetTokenBase().Type {
	case TokenType_IDENTIFIER, TokenType_NUMBER, TokenType_BIGINT, TokenType_STAR, TokenType_STATIC, TokenType_STRING, TokenType_OPEN_SQUARE:
		return true
	default:
		if Keywords_IsKeywordTokenType(token.GetTokenBase().Type) {
			return true
		}
	}
	return false
}

// Java: Parser.createObjectLiteralElementInfo()
func (p *Parser) CreateObjectLiteralElementInfo() *ClassOrObjectElementInfo {
	return ClassOrObjectElementInfo_CreateObjectLiteralElementInfo(p.GetTreeStartLocation())
}

// Java: Parser.parseClassElement()
func (p *Parser) ParseClassElement() ParseTree {
	if p.Peek(TokenType_SEMI_COLON) {
		return p.ParseEmptyStatement()
	}

	if p.PeekClassStaticInitializerBlock() {
		return p.ParseClassStaticInitializerBlock()
	}

	elementInfo := ClassOrObjectElementInfo_CreateClassMemberInfo(
		p.GetTreeStartLocation(), p.EatStaticIfNotElementName())
	if p.PeekGetAccessor() {
		return p.ParseGetAccessor(elementInfo)
	} else if p.PeekSetAccessor() {
		return p.ParseSetAccessor(elementInfo)
	} else if p.PeekAsyncMethod() {
		return p.ParseAsyncMethod(elementInfo)
	} else {
		return p.ParseClassMemberDeclaration(elementInfo)
	}
}

// Java: Parser.eatStaticIfNotElementName()
func (p *Parser) EatStaticIfNotElementName() bool {
	// only eat `static` if it being used as a keyword and not
	// a member name.
	if p.Peek(TokenType_STATIC) && p.IsClassElementStart(p.PeekTokenAt(1)) {
		p.Eat(TokenType_STATIC)
		return true
	}
	return false
}

// Java: Parser.peekAsyncMethod()
func (p *Parser) PeekAsyncMethod() bool {
	return p.PeekPredefinedString(Parser_ASYNC) &&
		!p.PeekImplicitSemiColonAt(1) &&
		(p.PeekPropertyNameOrComputedProp(1) ||
			(p.PeekAt(1, TokenType_STAR) && p.PeekPropertyNameOrComputedProp(2)))
}

// Java: Parser.peekClassStaticInitializerBlock()
func (p *Parser) PeekClassStaticInitializerBlock() bool {
	return p.Peek(TokenType_STATIC) && p.PeekAt(1, TokenType_OPEN_CURLY)
}

// Java: Parser.parseClassElementName(ClassOrObjectElementInfo)
func (p *Parser) ParseClassElementName(elementInfo *ClassOrObjectElementInfo) {
	if p.PeekPropertyName(0) {
		if p.PeekIdOrKeyword() {
			elementInfo.SetName(p.EatIdOrKeywordAsId())
			if elementInfo.GetName().IsKeyword() {
				p.RecordFeatureUsed(FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
			}
		} else {
			// { 'str'() {} }
			// { 123() {} }
			// Treat these as if they were computed properties.
			// TODO(b/123769080): Stop making this assumption!
			elementInfo.SetNameExpr(p.ParseLiteralExpression())
		}
	} else {
		elementInfo.SetNameExpr(p.ParseComputedPropertyName())
	}
}

// Java: Parser.parseFieldDefinition(ClassOrObjectElementInfo)
func (p *Parser) ParseFieldDefinition(elementInfo *ClassOrObjectElementInfo) ParseTree {
	var initializer ParseTree
	if p.Peek(TokenType_EQUAL) {
		initializer = p.ParseInitializer(Parser_Expression_NORMAL)
	}

	p.EatPossiblyImplicitSemiColon()
	if elementInfo.HasName() {
		return NewFieldDeclarationTree(
			p.GetTreeLocation(elementInfo.start),
			elementInfo.GetName(),
			elementInfo.isStatic,
			initializer)
	} else {
		jutil.CheckState(elementInfo.HasNameExpr())
		return NewComputedPropertyFieldTree(
			p.GetTreeLocation(elementInfo.start),
			elementInfo.GetNameExpr(),
			elementInfo.isStatic,
			initializer)
	}
}

// Java: Parser.parseMethodDefinition(ClassOrObjectElementInfo, boolean)
func (p *Parser) ParseMethodDefinition(elementInfo *ClassOrObjectElementInfo, isGenerator bool) ParseTree {
	var kind FunctionDeclarationTree_Kind
	if elementInfo.HasNameExpr() {
		kind = FunctionDeclarationTree_Kind_EXPRESSION
	} else {
		kind = FunctionDeclarationTree_Kind_MEMBER
	}

	builder := FunctionDeclarationTree_NewBuilder(kind).
		SetIsClassMember(elementInfo.isClassMember).
		SetStatic(elementInfo.isStatic)
	if elementInfo.HasName() {
		builder.SetName(elementInfo.GetName())
	}
	p.ParseFunctionTail(builder, functionFlavorOf(isGenerator, Parser_FunctionFlavor_GENERATOR, Parser_FunctionFlavor_NORMAL))

	var function ParseTree = builder.Build(p.GetTreeLocation(elementInfo.start))
	if kind == FunctionDeclarationTree_Kind_MEMBER {
		return function
	} else {
		jutil.CheckState(elementInfo.HasNameExpr())
		return NewComputedPropertyMethodTree(
			p.GetTreeLocation(elementInfo.start), elementInfo.GetNameExpr(), function)
	}
}

// functionFlavorOf is the Java conditional `cond ? ifTrue : ifFalse` over
// FunctionFlavor values.
//
// Go-only: language-adapter — Java's conditional operator.
func functionFlavorOf(cond bool, ifTrue, ifFalse Parser_FunctionFlavor) Parser_FunctionFlavor {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// Java: Parser.parseObjectLiteralMethodDeclaration()
func (p *Parser) ParseObjectLiteralMethodDeclaration() ParseTree {
	return p.ParseMethodDeclaration(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseMethodDeclaration(ClassOrObjectElementInfo)
func (p *Parser) ParseMethodDeclaration(elementInfo *ClassOrObjectElementInfo) ParseTree {
	isGenerator := p.EatOpt(TokenType_STAR) != nil
	p.ParseClassElementName(elementInfo)
	return p.ParseMethodDefinition(elementInfo, isGenerator)
}

// Java: Parser.parseClassMemberDeclaration(ClassOrObjectElementInfo)
func (p *Parser) ParseClassMemberDeclaration(elementInfo *ClassOrObjectElementInfo) ParseTree {
	isGenerator := p.EatOpt(TokenType_STAR) != nil
	p.ParseClassElementName(elementInfo)
	if p.PeekTypeAt(0) == TokenType_OPEN_PAREN {
		return p.ParseMethodDefinition(elementInfo, isGenerator)
	} else {
		return p.ParseFieldDefinition(elementInfo)
	}
}

// Java: Parser.parseObjectLiteralAsyncMethod()
func (p *Parser) ParseObjectLiteralAsyncMethod() ParseTree {
	return p.ParseAsyncMethod(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseAsyncMethod(ClassOrObjectElementInfo)
func (p *Parser) ParseAsyncMethod(elementInfo *ClassOrObjectElementInfo) ParseTree {
	p.EatPredefinedString(Parser_ASYNC)
	generator := p.Peek(TokenType_STAR)
	if generator {
		p.Eat(TokenType_STAR)
	}
	if p.PeekPropertyName(0) {
		if p.PeekIdOrKeyword() {
			name := p.EatIdOrKeywordAsId()
			builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_MEMBER).
				SetAsync(true).
				SetGenerator(generator).
				SetIsClassMember(elementInfo.isClassMember).
				SetStatic(elementInfo.isStatic).
				SetName(name)
			p.ParseFunctionTail(
				builder,
				functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))

			return builder.Build(p.GetTreeLocation(name.GetStart()))
		} else {
			// { 'str'() {} }
			// { 123() {} }
			// Treat these as if they were computed properties.
			var nameExpr ParseTree = p.ParseLiteralExpression()
			builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
				SetAsync(true).
				SetGenerator(generator).
				SetIsClassMember(elementInfo.isClassMember).
				SetStatic(elementInfo.isStatic)
			p.ParseFunctionTail(
				builder,
				functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))

			var function ParseTree = builder.Build(p.GetTreeLocation(nameExpr.GetStart()))
			return NewComputedPropertyMethodTree(
				p.GetTreeLocation(nameExpr.GetStart()), nameExpr, function)
		}
	} else {
		// expect '[' to start computed property name
		nameExpr := p.ParseComputedPropertyName()
		builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
			SetAsync(true).
			SetGenerator(generator).
			SetIsClassMember(elementInfo.isClassMember).
			SetStatic(elementInfo.isStatic)
		p.ParseFunctionTail(
			builder, functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))

		var function ParseTree = builder.Build(p.GetTreeLocation(nameExpr.GetStart()))
		return NewComputedPropertyMethodTree(
			p.GetTreeLocation(nameExpr.GetStart()), nameExpr, function)
	}
}

// Java: Parser.parseClassStaticInitializerBlock()
func (p *Parser) ParseClassStaticInitializerBlock() ParseTree {
	p.Eat(TokenType_STATIC)
	return p.ParseBlock()
}

// Java: Parser.parseFunctionTail(FunctionDeclarationTree.Builder, FunctionFlavor)
func (p *Parser) ParseFunctionTail(builder *FunctionDeclarationTree_Builder, functionFlavor Parser_FunctionFlavor) {
	p.functionContextStack = append(p.functionContextStack, functionFlavor)
	builder.
		SetGenerator(functionFlavor.IsGenerator()).
		SetFormalParameterList(p.ParseFormalParameterList()).
		SetFunctionBody(p.ParseFunctionBody())
	p.functionContextStack = p.functionContextStack[:len(p.functionContextStack)-1]
}

// Java: Parser.parseSourceElement()
func (p *Parser) ParseSourceElement() ParseTree {
	if p.PeekAsyncFunctionStart() {
		return p.ParseAsyncFunctionDeclaration()
	}

	if p.PeekFunction() {
		return p.ParseFunctionDeclaration()
	}

	if p.PeekClassDeclaration() {
		return p.ParseClassDeclaration()
	}

	// Harmony let block scoped bindings. let can only appear in
	// a block, not as a standalone statement: if() let x ... illegal
	if p.Peek(TokenType_LET) {
		return p.ParseVariableStatement()
	}
	// const and var are handled inside parseStatement

	return p.ParseStatementStandard()
}

// Java: Parser.peekSourceElement()
func (p *Parser) PeekSourceElement() bool {
	return p.PeekFunction() || p.PeekStatementStandard() || p.PeekDeclaration()
}

// Java: Parser.peekAsyncFunctionStart()
func (p *Parser) PeekAsyncFunctionStart() bool {
	return p.PeekPredefinedString(Parser_ASYNC) && !p.PeekImplicitSemiColonAt(1) && p.PeekFunctionAt(1)
}

// Java: Parser.eatAsyncFunctionStart()
func (p *Parser) EatAsyncFunctionStart() {
	p.EatPredefinedString(Parser_ASYNC)
	p.Eat(TokenType_FUNCTION)
}

// Java: Parser.peekFunction()
func (p *Parser) PeekFunction() bool {
	return p.PeekFunctionAt(0)
}

// Java: Parser.peekDeclaration()
func (p *Parser) PeekDeclaration() bool {
	return p.Peek(TokenType_LET) || p.PeekClassDeclaration()
}

// Java: Parser.peekFunction(int)
func (p *Parser) PeekFunctionAt(index int) bool {
	return p.PeekAt(index, TokenType_FUNCTION)
}

// 13 Function Definition

// Java: Parser.parseFunctionDeclaration()
func (p *Parser) ParseFunctionDeclaration() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(Keywords_FUNCTION.Type())
	isGenerator := p.EatOpt(TokenType_STAR) != nil

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_DECLARATION).SetName(p.EatId())
	p.ParseFunctionTail(builder, functionFlavorOf(isGenerator, Parser_FunctionFlavor_GENERATOR, Parser_FunctionFlavor_NORMAL))
	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.parseFunctionExpression()
func (p *Parser) ParseFunctionExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(Keywords_FUNCTION.Type())
	isGenerator := p.EatOpt(TokenType_STAR) != nil

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
		SetName(p.EatIdOpt())
	p.ParseFunctionTail(builder, functionFlavorOf(isGenerator, Parser_FunctionFlavor_GENERATOR, Parser_FunctionFlavor_NORMAL))

	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.parseAsyncFunctionDeclaration()
func (p *Parser) ParseAsyncFunctionDeclaration() ParseTree {
	start := p.GetTreeStartLocation()
	p.EatAsyncFunctionStart()

	generator := p.Peek(TokenType_STAR)
	if generator {
		p.Eat(TokenType_STAR)
	}

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_DECLARATION).
		SetName(p.EatId()).
		SetAsync(true)

	p.ParseFunctionTail(
		builder, functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))
	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.parseAsyncFunctionExpression()
func (p *Parser) ParseAsyncFunctionExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.EatAsyncFunctionStart()

	generator := p.Peek(TokenType_STAR)
	if generator {
		p.Eat(TokenType_STAR)
	}

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION).
		SetName(p.EatIdOpt()).
		SetAsync(true)

	p.ParseFunctionTail(
		builder, functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))
	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.peekParameter()
func (p *Parser) PeekParameter() bool {
	if p.PeekId() || p.Peek(TokenType_ELLIPSIS) {
		return true
	}
	return p.Peek(TokenType_OPEN_SQUARE) || p.Peek(TokenType_OPEN_CURLY)
}

// Java: Parser.parseParameter()
func (p *Parser) ParseParameter() ParseTree {
	start := p.GetTreeStartLocation()
	var parameter ParseTree

	if p.Peek(TokenType_ELLIPSIS) {
		parameter = p.ParseIterRest(Parser_PatternKind_INITIALIZER)
	} else if p.PeekId() {
		parameter = p.ParseIdentifierExpression()
	} else if p.PeekPatternStart() {
		parameter = p.ParsePattern(Parser_PatternKind_INITIALIZER)
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"parseParameterCalled() without confirming a parameter exists."))
	}

	if !parameter.IsRestParameter() && p.Peek(TokenType_EQUAL) {
		p.Eat(TokenType_EQUAL)
		defaultValue := p.ParseAssignmentExpression()
		parameter = NewDefaultParameterTree(p.GetTreeLocation(start), parameter, defaultValue)
	}

	return parameter
}

// Java: Parser.parseFormalParameterList()
func (p *Parser) ParseFormalParameterList() *FormalParameterListTree {
	listStart := p.GetTreeStartLocation()
	p.Eat(TokenType_OPEN_PAREN)

	treeBase := len(p.treeStack)
	trailingComma := false
	posBase := len(p.posStack)

	for p.PeekParameter() {
		param := p.ParseParameter()
		p.treeStack = append(p.treeStack, param)

		if !p.Peek(TokenType_CLOSE_PAREN) {
			comma := p.Eat(TokenType_COMMA)
			if comma != nil {
				p.posStack = append(p.posStack, comma.GetStart())
			} else {
				// semi-arbitrary comma position in case the code is syntactially invalid & missing one
				p.posStack = append(p.posStack, p.GetTreeEndLocation())
			}
			if p.Peek(TokenType_CLOSE_PAREN) {
				if !p.config.atLeast8 {
					p.ReportErrorAtToken(comma, "Invalid trailing comma in formal parameter list")
				}
				trailingComma = true
			}
		}
	}

	p.Eat(TokenType_CLOSE_PAREN)
	result := popStack(&p.treeStack, treeBase)
	commaPositions := popStack(&p.posStack, posBase)
	return NewFormalParameterListTree(
		p.GetTreeLocation(listStart), result, trailingComma, commaPositions)
}

// Java: Parser.parseSetterParameterList()
func (p *Parser) ParseSetterParameterList() *FormalParameterListTree {
	parameterList := p.ParseFormalParameterList()

	if len(parameterList.Parameters) != 1 {
		p.ReportErrorAtTree(
			parameterList,
			"Setter must have exactly 1 parameter, found %d",
			len(parameterList.Parameters))
	}

	if len(parameterList.Parameters) >= 1 {
		parameter := parameterList.Parameters[0]
		if parameter.IsRestParameter() {
			p.ReportErrorAtTree(parameter, "Setter must not have a rest parameter")
		}
	}

	return parameterList
}

// Java: Parser.parseFunctionBody()
func (p *Parser) ParseFunctionBody() *BlockTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_OPEN_CURLY)
	result := p.ParseSourceElementList()
	p.Eat(TokenType_CLOSE_CURLY)
	return NewBlockTree(p.GetTreeLocation(start), result)
}

// Java: Parser.parseSourceElementList()
func (p *Parser) ParseSourceElementList() []ParseTree {
	base := len(p.treeStack)

	for p.PeekSourceElement() {
		elem := p.ParseSourceElement()
		p.treeStack = append(p.treeStack, elem)
	}

	return popStack(&p.treeStack, base)
}

// Java: Parser.parseIterSpread()
func (p *Parser) ParseIterSpread() *IterSpreadTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_ELLIPSIS)
	operand := p.ParseAssignmentExpression()
	return NewIterSpreadTree(p.GetTreeLocation(start), operand)
}

// 12 Statements

// ParseStatement: in V8, all source elements may appear where statements
// occur in the grammar.
//
// Java: Parser.parseStatement()
func (p *Parser) ParseStatement() ParseTree {
	return p.ParseSourceElement()
}

// ParseStatementStandard reflects the ECMA standard. Most places use
// peekStatement instead.
//
// Java: Parser.parseStatementStandard()
func (p *Parser) ParseStatementStandard() ParseTree {
	switch p.PeekType() {
	case TokenType_OPEN_CURLY:
		return p.ParseBlock()
	case TokenType_CONST, TokenType_VAR:
		return p.ParseVariableStatement()
	case TokenType_SEMI_COLON:
		return p.ParseEmptyStatement()
	case TokenType_IF:
		return p.ParseIfStatement()
	case TokenType_DO:
		return p.ParseDoWhileStatement()
	case TokenType_WHILE:
		return p.ParseWhileStatement()
	case TokenType_FOR:
		return p.ParseForStatement()
	case TokenType_CONTINUE:
		return p.ParseContinueStatement()
	case TokenType_BREAK:
		return p.ParseBreakStatement()
	case TokenType_RETURN:
		return p.ParseReturnStatement()
	case TokenType_WITH:
		return p.ParseWithStatement()
	case TokenType_SWITCH:
		return p.ParseSwitchStatement()
	case TokenType_THROW:
		return p.ParseThrowStatement()
	case TokenType_TRY:
		return p.ParseTryStatement()
	case TokenType_DEBUGGER:
		return p.ParseDebuggerStatement()
	default:
		if p.PeekLabelledStatement() {
			return p.ParseLabelledStatement()
		}
		return p.ParseExpressionStatement()
	}
}

// PeekStatement: in V8 all source elements may appear where statements
// appear in the grammar.
//
// Java: Parser.peekStatement()
func (p *Parser) PeekStatement() bool {
	return p.PeekSourceElement()
}

// PeekStatementStandard reflects the ECMA standard. Most places use
// peekStatement instead.
//
// Java: Parser.peekStatementStandard()
func (p *Parser) PeekStatementStandard() bool {
	switch p.PeekType() {
	case TokenType_OPEN_CURLY,
		TokenType_VAR,
		TokenType_CONST,
		TokenType_SEMI_COLON,
		TokenType_IF,
		TokenType_DO,
		TokenType_WHILE,
		TokenType_FOR,
		TokenType_CONTINUE,
		TokenType_BREAK,
		TokenType_RETURN,
		TokenType_WITH,
		TokenType_SWITCH,
		TokenType_THROW,
		TokenType_TRY,
		TokenType_DEBUGGER,
		TokenType_YIELD,
		TokenType_IDENTIFIER,
		TokenType_TYPE,
		TokenType_DECLARE,
		TokenType_MODULE,
		TokenType_NAMESPACE,
		TokenType_THIS,
		TokenType_CLASS,
		TokenType_SUPER,
		TokenType_NUMBER,
		TokenType_BIGINT,
		TokenType_STRING,
		TokenType_NO_SUBSTITUTION_TEMPLATE,
		TokenType_TEMPLATE_HEAD,
		TokenType_NULL,
		TokenType_TRUE,
		TokenType_SLASH,       // regular expression literal
		TokenType_SLASH_EQUAL, // regular expression literal
		TokenType_FALSE,
		TokenType_OPEN_SQUARE,
		TokenType_OPEN_PAREN,
		TokenType_NEW,
		TokenType_DELETE,
		TokenType_VOID,
		TokenType_TYPEOF,
		TokenType_PLUS_PLUS,
		TokenType_MINUS_MINUS,
		TokenType_PLUS,
		TokenType_MINUS,
		TokenType_TILDE,
		TokenType_BANG,
		TokenType_IMPORT:
		return true
	default:
		return false
	}
}

// 12.1 Block

// Java: Parser.parseBlock()
func (p *Parser) ParseBlock() *BlockTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_OPEN_CURLY)
	// Spec says Statement list. However functions are also embedded in the wild.
	result := p.ParseSourceElementList()
	p.Eat(TokenType_CLOSE_CURLY)
	return NewBlockTree(p.GetTreeLocation(start), result)
}

// Java: Parser.parseStatementList()
func (p *Parser) ParseStatementList() []ParseTree {
	base := len(p.treeStack)
	for p.PeekStatement() {
		stmt := p.ParseStatement()
		p.treeStack = append(p.treeStack, stmt)
	}
	return popStack(&p.treeStack, base)
}

// 12.2 Variable Statement

// Java: Parser.parseVariableStatement()
func (p *Parser) ParseVariableStatement() *VariableStatementTree {
	start := p.GetTreeStartLocation()
	declarations := p.ParseVariableDeclarationList()
	p.EatPossiblyImplicitSemiColon()
	return NewVariableStatementTree(p.GetTreeLocation(start), declarations)
}

// Java: Parser.parseVariableDeclarationList()
func (p *Parser) ParseVariableDeclarationList() *VariableDeclarationListTree {
	return p.ParseVariableDeclarationListWithExpression(Parser_Expression_NORMAL)
}

// Java: Parser.parseVariableDeclarationListNoIn()
func (p *Parser) ParseVariableDeclarationListNoIn() *VariableDeclarationListTree {
	return p.ParseVariableDeclarationListWithExpression(Parser_Expression_NO_IN)
}

// ParseVariableDeclarationListWithExpression returns nil (Java null) after
// reporting an error if no declaration keyword is next.
//
// Java: Parser.parseVariableDeclarationList(Expression)
func (p *Parser) ParseVariableDeclarationListWithExpression(expressionIn Parser_Expression) *VariableDeclarationListTree {
	start := p.GetTreeStartLocation()
	token := p.PeekType()

	switch token {
	case TokenType_CONST, TokenType_LET, TokenType_VAR:
		p.Eat(token)
	default:
		p.ReportErrorAtToken(p.PeekToken(), "expected declaration")
		return nil
	}

	base := len(p.varDeclStack)

	decl := p.ParseVariableDeclaration(token, expressionIn)
	p.varDeclStack = append(p.varDeclStack, decl)
	for p.Peek(TokenType_COMMA) {
		p.Eat(TokenType_COMMA)
		decl = p.ParseVariableDeclaration(token, expressionIn)
		p.varDeclStack = append(p.varDeclStack, decl)
	}
	declarations := popStack(&p.varDeclStack, base)
	return NewVariableDeclarationListTree(p.GetTreeLocation(start), token, declarations)
}

// Java: Parser.parseVariableDeclaration(TokenType, Expression)
func (p *Parser) ParseVariableDeclaration(binding TokenType, expressionIn Parser_Expression) *VariableDeclarationTree {

	start := p.GetTreeStartLocation()
	var lvalue ParseTree
	if p.PeekPatternStart() {
		lvalue = p.ParsePattern(Parser_PatternKind_INITIALIZER)
	} else {
		lvalue = p.ParseIdentifierExpression()
	}

	var initializer ParseTree
	if p.Peek(TokenType_EQUAL) {
		initializer = p.ParseInitializer(expressionIn)
	} else if expressionIn != Parser_Expression_NO_IN {
		// NOTE(blickly): this is a bit of a hack, declarations outside of for statements allow "in",
		// and by chance, also must have initializers for const/destructuring. Vanilla for loops
		// also require intializers, but are handled separately in checkVanillaForInitializers
		p.MaybeReportNoInitializer(binding, lvalue)
	}
	return NewVariableDeclarationTree(p.GetTreeLocation(start), lvalue, initializer)
}

// Java: Parser.parseInitializer(Expression)
func (p *Parser) ParseInitializer(expressionIn Parser_Expression) ParseTree {
	p.Eat(TokenType_EQUAL)
	return p.ParseAssignment(expressionIn)
}

// 12.3 Empty Statement

// Java: Parser.parseEmptyStatement()
func (p *Parser) ParseEmptyStatement() *EmptyStatementTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_SEMI_COLON)
	return NewEmptyStatementTree(p.GetTreeLocation(start))
}

// 12.4 Expression Statement

// Java: Parser.parseExpressionStatement()
func (p *Parser) ParseExpressionStatement() *ExpressionStatementTree {
	start := p.GetTreeStartLocation()
	expression := p.ParseExpression()
	p.EatPossiblyImplicitSemiColon()
	return NewExpressionStatementTree(p.GetTreeLocation(start), expression)
}

// 12.5 If Statement

// Java: Parser.parseIfStatement()
func (p *Parser) ParseIfStatement() *IfStatementTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IF)
	p.Eat(TokenType_OPEN_PAREN)
	condition := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	ifClause := p.ParseStatement()
	var elseClause ParseTree
	if p.Peek(TokenType_ELSE) {
		p.Eat(TokenType_ELSE)
		elseClause = p.ParseStatement()
	}
	return NewIfStatementTree(p.GetTreeLocation(start), condition, ifClause, elseClause)
}

// 12.6 Iteration Statements

// 12.6.1 The do-while Statement

// Java: Parser.parseDoWhileStatement()
func (p *Parser) ParseDoWhileStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_DO)
	body := p.ParseStatement()
	p.Eat(TokenType_WHILE)
	p.Eat(TokenType_OPEN_PAREN)
	condition := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	// The semicolon after the "do-while" is optional.
	if p.Peek(TokenType_SEMI_COLON) {
		p.Eat(TokenType_SEMI_COLON)
	}
	return NewDoWhileStatementTree(p.GetTreeLocation(start), body, condition)
}

// 12.6.2 The while Statement

// Java: Parser.parseWhileStatement()
func (p *Parser) ParseWhileStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_WHILE)
	p.Eat(TokenType_OPEN_PAREN)
	condition := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewWhileStatementTree(p.GetTreeLocation(start), condition, body)
}

// 12.6.3 The for Statement
// 12.6.4 The for-in Statement
// The for-of Statement
// The for-await-of Statement

// Java: Parser.parseForStatement()
func (p *Parser) ParseForStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_FOR)
	awaited := p.PeekPredefinedString(Parser_AWAIT)
	if awaited {
		p.EatPredefinedString(Parser_AWAIT)
	}
	p.Eat(TokenType_OPEN_PAREN)
	if p.PeekVariableDeclarationList() {
		variables := p.ParseVariableDeclarationListNoIn()
		if p.Peek(TokenType_IN) {
			if awaited {
				p.ReportError("for-await-of is the only allowed asynchronous iteration")
			}
			// for-in: only one declaration allowed
			if len(variables.Declarations) > 1 {
				p.ReportError("for-in statement may not have more than one variable declaration")
			}
			declaration := variables.Declarations[0]
			if declaration.Initializer != nil {
				// An initializer is allowed here in ES5 and below, but not in ES6.
				// Warn about it, to encourage people to eliminate it from their code.
				// http://esdiscuss.org/topic/initializer-expression-on-for-in-syntax-subject
				if p.config.atLeast6 {
					p.ReportError("for-in statement may not have initializer")
				} else {
					p.errorReporter.ReportWarningFormat(
						declaration.Location.Start, "for-in statement should not have initializer")
				}
			}

			return p.ParseForInStatement(start, variableDeclarationListAsParseTree(variables))
		} else if p.PeekPredefinedString(PredefinedName_OF) {
			// for-of: only one declaration allowed
			if len(variables.Declarations) > 1 {
				if awaited {
					p.ReportError("for-await-of statement may not have more than one variable declaration")
				} else {
					p.ReportError("for-of statement may not have more than one variable declaration")
				}
			}
			// for-of: initializer is illegal
			declaration := variables.Declarations[0]
			if declaration.Initializer != nil {
				if awaited {
					p.ReportError("for-await-of statement may not have initializer")
				} else {
					p.ReportError("for-of statement may not have initializer")
				}
			}

			if awaited {
				return p.ParseForAwaitOfStatement(start, variableDeclarationListAsParseTree(variables))
			} else {
				return p.ParseForOfStatement(start, variableDeclarationListAsParseTree(variables))
			}
		} else {
			// "Vanilla" for statement: const/destructuring must have initializer
			p.CheckVanillaForInitializers(variables)
			return p.ParseForStatementWithInitializer(start, variableDeclarationListAsParseTree(variables))
		}
	}

	if p.Peek(TokenType_SEMI_COLON) {
		return p.ParseForStatementWithInitializer(start, nil)
	}

	initializer := p.ParseExpressionNoIn()
	if p.Peek(TokenType_IN) || p.Peek(TokenType_EQUAL) || p.PeekPredefinedString(PredefinedName_OF) {
		initializer = p.TransformLeftHandSideExpression(initializer)
		if !initializer.IsValidAssignmentTarget() {
			p.ReportError("invalid assignment target")
		}
	}

	if p.Peek(TokenType_IN) || p.PeekPredefinedString(PredefinedName_OF) {
		if initializer.GetParseTreeBase().Type != ParseTreeType_BINARY_OPERATOR &&
			initializer.GetParseTreeBase().Type != ParseTreeType_COMMA_EXPRESSION {
			if p.Peek(TokenType_IN) {
				return p.ParseForInStatement(start, initializer)
			} else {
				// for {await}? ( _ of _ )
				if awaited {
					return p.ParseForAwaitOfStatement(start, initializer)
				} else {
					return p.ParseForOfStatement(start, initializer)
				}
			}
		}
	}

	return p.ParseForStatementWithInitializer(start, initializer)
}

// The for-of Statement
// for  (  { let | var }?  identifier  of  expression  )  statement

// Java: Parser.parseForOfStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForOfStatement(start SourcePosition, initializer ParseTree) ParseTree {
	p.EatPredefinedString(PredefinedName_OF)
	collection := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForOfStatementTree(p.GetTreeLocation(start), initializer, collection, body)
}

// Java: Parser.parseForAwaitOfStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForAwaitOfStatement(start SourcePosition, initializer ParseTree) ParseTree {
	p.EatPredefinedString(PredefinedName_OF)
	collection := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForAwaitOfStatementTree(p.GetTreeLocation(start), initializer, collection, body)
}

// CheckVanillaForInitializers checks variable declarations in for
// statements.
//
// Java: Parser.checkVanillaForInitializers(VariableDeclarationListTree)
func (p *Parser) CheckVanillaForInitializers(variables *VariableDeclarationListTree) {
	for _, declaration := range variables.Declarations {
		if declaration.Initializer == nil {
			p.MaybeReportNoInitializer(variables.DeclarationType, declaration.Lvalue)
		}
	}
}

// MaybeReportNoInitializer reports if declaration requires an initializer,
// assuming initializer is absent.
//
// Java: Parser.maybeReportNoInitializer(TokenType, ParseTree)
func (p *Parser) MaybeReportNoInitializer(token TokenType, lvalue ParseTree) {
	if token == TokenType_CONST {
		p.ReportError("const variables must have an initializer")
	} else if lvalue.IsPattern() {
		p.ReportError("destructuring must have an initializer")
	}
}

// Java: Parser.peekVariableDeclarationList()
func (p *Parser) PeekVariableDeclarationList() bool {
	switch p.PeekType() {
	case TokenType_VAR, TokenType_CONST, TokenType_LET:
		return true
	default:
		return false
	}
}

// 12.6.3 The for Statement

// Java: Parser.parseForStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForStatementWithInitializer(start SourcePosition, initializer ParseTree) ParseTree {
	if initializer == nil {
		initializer = NewNullTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))
	}
	p.Eat(TokenType_SEMI_COLON)

	var condition ParseTree
	if !p.Peek(TokenType_SEMI_COLON) {
		condition = p.ParseExpression()
	} else {
		condition = NewNullTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))
	}
	p.Eat(TokenType_SEMI_COLON)

	var increment ParseTree
	if !p.Peek(TokenType_CLOSE_PAREN) {
		increment = p.ParseExpression()
	} else {
		increment = NewNullTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))
	}
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForStatementTree(p.GetTreeLocation(start), initializer, condition, increment, body)
}

// 12.6.4 The for-in Statement

// Java: Parser.parseForInStatement(SourcePosition, ParseTree)
func (p *Parser) ParseForInStatement(start SourcePosition, initializer ParseTree) ParseTree {
	p.Eat(TokenType_IN)
	collection := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewForInStatementTree(p.GetTreeLocation(start), initializer, collection, body)
}

// 12.7 The continue Statement

// Java: Parser.parseContinueStatement()
func (p *Parser) ParseContinueStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_CONTINUE)
	var name *IdentifierToken
	if !p.PeekImplicitSemiColon() {
		name = p.EatIdOpt()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewContinueStatementTree(p.GetTreeLocation(start), name)
}

// 12.8 The break Statement

// Java: Parser.parseBreakStatement()
func (p *Parser) ParseBreakStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_BREAK)
	var name *IdentifierToken
	if !p.PeekImplicitSemiColon() {
		name = p.EatIdOpt()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewBreakStatementTree(p.GetTreeLocation(start), name)
}

// 12.9 The return Statement

// Java: Parser.parseReturnStatement()
func (p *Parser) ParseReturnStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_RETURN)
	var expression ParseTree
	if !p.PeekImplicitSemiColon() {
		expression = p.ParseExpression()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewReturnStatementTree(p.GetTreeLocation(start), expression)
}

// 12.10 The with Statement

// Java: Parser.parseWithStatement()
func (p *Parser) ParseWithStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_WITH)
	p.Eat(TokenType_OPEN_PAREN)
	expression := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	body := p.ParseStatement()
	return NewWithStatementTree(p.GetTreeLocation(start), expression, body)
}

// 12.11 The switch Statement

// Java: Parser.parseSwitchStatement()
func (p *Parser) ParseSwitchStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_SWITCH)
	p.Eat(TokenType_OPEN_PAREN)
	expression := p.ParseExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	p.Eat(TokenType_OPEN_CURLY)
	caseClauses := p.ParseCaseClauses()
	p.Eat(TokenType_CLOSE_CURLY)
	return NewSwitchStatementTree(p.GetTreeLocation(start), expression, caseClauses)
}

// Java: Parser.parseCaseClauses()
func (p *Parser) ParseCaseClauses() []ParseTree {
	foundDefaultClause := false
	base := len(p.treeStack)

	for {
		start := p.GetTreeStartLocation()
		switch p.PeekType() {
		case TokenType_CASE:
			p.Eat(TokenType_CASE)
			expression := p.ParseExpression()
			p.Eat(TokenType_COLON)
			statements := p.ParseCaseStatementsOpt()
			clause := NewCaseClauseTree(p.GetTreeLocation(start), expression, statements)
			p.treeStack = append(p.treeStack, clause)
		case TokenType_DEFAULT:
			if foundDefaultClause {
				p.ReportError("Switch statements may have at most one default clause")
			} else {
				foundDefaultClause = true
			}
			p.Eat(TokenType_DEFAULT)
			p.Eat(TokenType_COLON)
			clause := NewDefaultClauseTree(p.GetTreeLocation(start), p.ParseCaseStatementsOpt())
			p.treeStack = append(p.treeStack, clause)
		default:
			return popStack(&p.treeStack, base)
		}
	}
}

// Java: Parser.parseCaseStatementsOpt()
func (p *Parser) ParseCaseStatementsOpt() []ParseTree {
	return p.ParseStatementList()
}

// 12.12 Labelled Statement

// Java: Parser.parseLabelledStatement()
func (p *Parser) ParseLabelledStatement() ParseTree {
	start := p.GetTreeStartLocation()
	name := p.EatId()
	p.Eat(TokenType_COLON)
	return NewLabelledStatementTree(p.GetTreeLocation(start), name, p.ParseStatement())
}

// Java: Parser.peekLabelledStatement()
func (p *Parser) PeekLabelledStatement() bool {
	return p.PeekId() && p.PeekAt(1, TokenType_COLON)
}

// 12.13 Throw Statement

// Java: Parser.parseThrowStatement()
func (p *Parser) ParseThrowStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_THROW)
	var value ParseTree
	if p.PeekImplicitSemiColon() {
		p.ReportError("semicolon/newline not allowed after 'throw'")
	} else {
		value = p.ParseExpression()
	}
	p.EatPossiblyImplicitSemiColon()
	return NewThrowStatementTree(p.GetTreeLocation(start), value)
}

// 12.14 Try Statement

// Java: Parser.parseTryStatement()
func (p *Parser) ParseTryStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_TRY)
	var body ParseTree = p.ParseBlock()
	var catchBlock ParseTree
	if p.Peek(TokenType_CATCH) {
		catchBlock = p.ParseCatch()
	}
	var finallyBlock ParseTree
	if p.Peek(TokenType_FINALLY) {
		finallyBlock = p.ParseFinallyBlock()
	}
	if catchBlock == nil && finallyBlock == nil {
		p.ReportError("'catch' or 'finally' expected.")
	}
	return NewTryStatementTree(p.GetTreeLocation(start), body, catchBlock, finallyBlock)
}

// Java: Parser.parseCatch()
func (p *Parser) ParseCatch() *CatchTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_CATCH)

	var exception ParseTree = NewEmptyStatementTree(NewSourceRange(p.GetTreeEndLocation(), p.GetTreeStartLocation()))

	if p.PeekToken().GetTokenBase().Type == TokenType_OPEN_PAREN {
		p.Eat(TokenType_OPEN_PAREN)
		if p.PeekPatternStart() {
			exception = p.ParsePattern(Parser_PatternKind_INITIALIZER)
		} else {
			exception = p.ParseIdentifierExpression()
		}
		p.Eat(TokenType_CLOSE_PAREN)
	} else {
		p.RecordFeatureUsed(FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
	}

	catchBody := p.ParseBlock()
	return NewCatchTree(p.GetTreeLocation(start), exception, catchBody)
}

// Java: Parser.parseFinallyBlock()
func (p *Parser) ParseFinallyBlock() *FinallyTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_FINALLY)
	finallyBlock := p.ParseBlock()
	return NewFinallyTree(p.GetTreeLocation(start), finallyBlock)
}

// 12.15 The Debugger Statement

// Java: Parser.parseDebuggerStatement()
func (p *Parser) ParseDebuggerStatement() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_DEBUGGER)
	p.EatPossiblyImplicitSemiColon()

	return NewDebuggerStatementTree(p.GetTreeLocation(start))
}

// 11.1 Primary Expressions

// Java: Parser.parsePrimaryExpression()
func (p *Parser) ParsePrimaryExpression() ParseTree {
	switch p.PeekType() {
	case TokenType_CLASS:
		return p.ParseClassExpression()
	case TokenType_SUPER:
		return p.ParseSuperExpression()
	case TokenType_THIS:
		return p.ParseThisExpression()
	case TokenType_IMPORT:
		return p.ParseDynamicImportExpression()
	case TokenType_IDENTIFIER, TokenType_TYPE, TokenType_DECLARE, TokenType_MODULE, TokenType_NAMESPACE:
		return p.ParseIdentifierExpression()
	case TokenType_NUMBER, TokenType_STRING, TokenType_BIGINT, TokenType_TRUE, TokenType_FALSE, TokenType_NULL:
		return p.ParseLiteralExpression()
	case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
		return p.ParseTemplateLiteral(nil)
	case TokenType_OPEN_SQUARE:
		return p.ParseArrayInitializer()
	case TokenType_OPEN_CURLY:
		return p.ParseObjectLiteral()
	case TokenType_OPEN_PAREN:
		return p.ParseCoverParenthesizedExpressionAndArrowParameterList()
	case TokenType_SLASH, TokenType_SLASH_EQUAL:
		return p.ParseRegularExpressionLiteral()
	default:
		return p.ParseMissingPrimaryExpression()
	}
}

// Java: Parser.parseSuperExpression()
func (p *Parser) ParseSuperExpression() *SuperExpressionTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_SUPER)
	if p.Peek(TokenType_QUESTION_DOT) { // super?.() not allowed
		p.ReportError("Optional chaining is forbidden in super?.")
	}
	return NewSuperExpressionTree(p.GetTreeLocation(start))
}

// Java: Parser.parseThisExpression()
func (p *Parser) ParseThisExpression() *ThisExpressionTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_THIS)
	return NewThisExpressionTree(p.GetTreeLocation(start))
}

// ParseDynamicImportExpression parses import(...).
// https://tc39.github.io/proposal-dynamic-import
//
// Java: Parser.parseDynamicImportExpression()
func (p *Parser) ParseDynamicImportExpression() *DynamicImportTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IMPORT)
	if p.Peek(TokenType_QUESTION_DOT) { // import?.() not allowed
		p.ReportError("Optional chaining is forbidden in import?.")
	}
	p.Eat(TokenType_OPEN_PAREN)
	argument := p.ParseAssignmentExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	p.RecordFeatureUsed(FeatureSet_Feature_DYNAMIC_IMPORT)
	return NewDynamicImportTree(p.GetTreeLocation(start), argument)
}

// Java: Parser.parseIdentifierExpression()
func (p *Parser) ParseIdentifierExpression() *IdentifierExpressionTree {
	start := p.GetTreeStartLocation()
	identifier := p.EatId()
	return NewIdentifierExpressionTree(p.GetTreeLocation(start), identifier)
}

// Java: Parser.parseLiteralExpression()
func (p *Parser) ParseLiteralExpression() *LiteralExpressionTree {
	start := p.GetTreeStartLocation()
	literal := p.NextLiteralToken()

	if literal.GetTokenBase().Type == TokenType_STRING &&
		literal.(*StringLiteralToken).HasUnescapedUnicodeLineOrParagraphSeparator() {
		p.RecordFeatureUsed(FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP)
	}

	if literal.GetTokenBase().Type == TokenType_NUMBER && strings.Contains(literal.String(), "_") {
		p.RecordFeatureUsed(FeatureSet_Feature_NUMERIC_SEPARATOR)
	}

	if literal.GetTokenBase().Type == TokenType_BIGINT {
		p.RecordFeatureUsed(FeatureSet_Feature_BIGINT)
	}

	return NewLiteralExpressionTree(p.GetTreeLocation(start), literal)
}

// ParseTemplateLiteral constructs a template literal expression tree.
// "operand" is used to handle the case like "foo`bar`", which is a
// CallExpression or MemberExpression that calls the function foo() with the
// template literal as the argument (with extra handling). In this case,
// operand would be "foo", which is the callsite.
//
// We store this operand in the TemplateLiteralExpressionTree and generate a
// TAGGED_TEMPLATELIT node if it's not null later when transpiling.
//
// operand: a non-null value would represent the callsite. Returns the
// template literal expression.
//
// Java: Parser.parseTemplateLiteral(ParseTree)
func (p *Parser) ParseTemplateLiteral(operand ParseTree) *TemplateLiteralExpressionTree {
	var start SourcePosition
	if operand == nil {
		start = p.GetTreeStartLocation()
	} else {
		start = operand.GetParseTreeBase().Location.Start
	}
	token := p.NextToken()
	if _, ok := token.(*TemplateLiteralToken); !ok {
		p.ReportErrorAtToken(token, "Unexpected template literal token %s.", token.GetTokenBase().Type.String())
	}
	isTaggedTemplate := operand != nil
	templateToken := token.(*TemplateLiteralToken)
	if !isTaggedTemplate {
		p.ReportTemplateErrorIfPresent(templateToken)
	}
	elements := []ParseTree{}
	elements = append(elements, NewTemplateLiteralPortionTree(templateToken.Location, templateToken))
	if templateToken.Type == TokenType_NO_SUBSTITUTION_TEMPLATE {
		return NewTemplateLiteralExpressionTree(p.GetTreeLocation(start), operand, elements)
	}

	// `abc${
	expression := p.ParseExpression()
	elements = append(elements, NewTemplateSubstitutionTree(expression.GetParseTreeBase().Location, expression))
	for !p.errorReporter.HadError() {
		templateToken = p.NextTemplateLiteralToken()
		if templateToken.Type == TokenType_ERROR || templateToken.Type == TokenType_END_OF_FILE {
			break
		}
		if !isTaggedTemplate {
			p.ReportTemplateErrorIfPresent(templateToken)
		}
		elements = append(elements, NewTemplateLiteralPortionTree(templateToken.Location, templateToken))
		if templateToken.Type == TokenType_TEMPLATE_TAIL {
			break
		}

		expression = p.ParseExpression()
		elements = append(elements, NewTemplateSubstitutionTree(expression.GetParseTreeBase().Location, expression))
	}

	return NewTemplateLiteralExpressionTree(p.GetTreeLocation(start), operand, elements)
}

// Java: Parser.nextLiteralToken()
func (p *Parser) NextLiteralToken() Token {
	return p.NextToken()
}

// Java: Parser.parseRegularExpressionLiteral()
func (p *Parser) ParseRegularExpressionLiteral() ParseTree {
	start := p.GetTreeStartLocation()
	literal := p.NextRegularExpressionLiteralToken()
	p.RecordFeatureUsed(FeatureSet_Feature_REGEXP_SYNTAX)
	return NewLiteralExpressionTree(p.GetTreeLocation(start), literal)
}

// Java: Parser.parseArrayInitializer()
func (p *Parser) ParseArrayInitializer() ParseTree {
	if p.PeekTypeAt(1) == TokenType_FOR {
		return p.ParseArrayComprehension()
	} else {
		return p.ParseArrayLiteral()
	}
}

// Java: Parser.parseGeneratorComprehension()
func (p *Parser) ParseGeneratorComprehension() ParseTree {
	return p.ParseComprehension(
		ComprehensionTree_ComprehensionType_GENERATOR, TokenType_OPEN_PAREN, TokenType_CLOSE_PAREN)
}

// Java: Parser.parseArrayComprehension()
func (p *Parser) ParseArrayComprehension() ParseTree {
	return p.ParseComprehension(
		ComprehensionTree_ComprehensionType_ARRAY, TokenType_OPEN_SQUARE, TokenType_CLOSE_SQUARE)
}

// Java: Parser.parseComprehension(ComprehensionTree.ComprehensionType, TokenType, TokenType)
func (p *Parser) ParseComprehension(type_ ComprehensionTree_ComprehensionType, startToken, endToken TokenType) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(startToken)

	children := []ParseTree{}
	for p.Peek(TokenType_FOR) || p.Peek(TokenType_IF) {
		if p.Peek(TokenType_FOR) {
			children = append(children, p.ParseComprehensionFor())
		} else {
			children = append(children, p.ParseComprehensionIf())
		}
	}

	tailExpression := p.ParseAssignmentExpression()
	p.Eat(endToken)

	return NewComprehensionTree(p.GetTreeLocation(start), type_, children, tailExpression)
}

// Java: Parser.parseComprehensionFor()
func (p *Parser) ParseComprehensionFor() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_FOR)
	p.Eat(TokenType_OPEN_PAREN)

	var initializer ParseTree
	if p.PeekId() {
		initializer = p.ParseIdentifierExpression()
	} else {
		initializer = p.ParsePattern(Parser_PatternKind_ANY)
	}

	p.EatPredefinedString(PredefinedName_OF)
	collection := p.ParseAssignmentExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	return NewComprehensionForTree(p.GetTreeLocation(start), initializer, collection)
}

// Java: Parser.parseComprehensionIf()
func (p *Parser) ParseComprehensionIf() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IF)
	p.Eat(TokenType_OPEN_PAREN)
	initializer := p.ParseAssignmentExpression()
	p.Eat(TokenType_CLOSE_PAREN)
	return NewComprehensionIfTree(p.GetTreeLocation(start), initializer)
}

// 11.1.4 Array Literal Expression

// Java: Parser.parseArrayLiteral()
func (p *Parser) ParseArrayLiteral() ParseTree {
	// ArrayLiteral :
	//   [ Elisionopt ]
	//   [ ElementList ]
	//   [ ElementList , Elisionopt ]
	//
	// ElementList :
	//   Elisionopt AssignmentOrSpreadExpression
	//   ElementList , Elisionopt AssignmentOrSpreadExpression
	//
	// Elision :
	//   ,
	//   Elision ,

	start := p.GetTreeStartLocation()
	elements := []ParseTree{}

	p.Eat(TokenType_OPEN_SQUARE)
	var trailingCommaToken Token
	for p.Peek(TokenType_COMMA) || p.Peek(TokenType_ELLIPSIS) || p.PeekAssignmentExpression() {
		trailingCommaToken = nil
		if p.Peek(TokenType_COMMA) {
			commaStart := p.GetTreeStartLocation()
			trailingCommaToken = p.Eat(TokenType_COMMA)
			// Consider the empty element to start & end immediately before the comma token.
			elements = append(elements, NewNullTree(NewSourceRange(commaStart, commaStart)))

		} else {
			if p.Peek(TokenType_ELLIPSIS) {
				p.RecordFeatureUsed(FeatureSet_Feature_SPREAD_EXPRESSIONS)
				elements = append(elements, p.ParseIterSpread())
			} else {
				elements = append(elements, p.ParseAssignmentExpression())
			}
			if !p.Peek(TokenType_CLOSE_SQUARE) {
				trailingCommaToken = p.Eat(TokenType_COMMA)
			}
		}
	}
	p.Eat(TokenType_CLOSE_SQUARE)

	p.MaybeReportTrailingComma(trailingCommaToken)

	return NewArrayLiteralExpressionTree(
		p.GetTreeLocation(start), elements, trailingCommaToken != nil)
}

// 11.1.4 Object Literal Expression

// Java: Parser.parseObjectLiteral()
func (p *Parser) ParseObjectLiteral() ParseTree {
	start := p.GetTreeStartLocation()
	base := len(p.treeStack)

	p.Eat(TokenType_OPEN_CURLY)
	var commaToken Token
	for p.Peek(TokenType_ELLIPSIS) || p.PeekPropertyNameOrComputedProp(0) || p.Peek(TokenType_STAR) {
		prop := p.ParseObjectLiteralPropertyAssignment()
		p.treeStack = append(p.treeStack, prop)
		commaToken = p.EatOpt(TokenType_COMMA)
		if commaToken == nil {
			break
		}
	}
	p.Eat(TokenType_CLOSE_CURLY)

	p.MaybeReportTrailingComma(commaToken)

	result := popStack(&p.treeStack, base)
	return NewObjectLiteralExpressionTree(
		p.GetTreeLocation(start), result, commaToken != nil)
}

// Java: Parser.maybeReportTrailingComma(Token)
func (p *Parser) MaybeReportTrailingComma(commaToken Token) {
	if commaToken != nil {
		p.RecordFeatureUsed(FeatureSet_Feature_TRAILING_COMMA)
		if p.config.warnTrailingCommas {
			// In ES3 mode warn about trailing commas which aren't accepted by
			// older browsers (such as IE8).
			p.errorReporter.ReportWarningFormat(
				commaToken.GetTokenBase().Location.Start,
				"Trailing comma is not legal in an ECMA-262 object initializer")
		}
	}
}

// Java: Parser.peekPropertyNameOrComputedProp(int)
func (p *Parser) PeekPropertyNameOrComputedProp(tokenIndex int) bool {
	return p.PeekPropertyName(tokenIndex) || p.PeekTypeAt(tokenIndex) == TokenType_OPEN_SQUARE
}

// Java: Parser.peekPropertyName(int)
func (p *Parser) PeekPropertyName(tokenIndex int) bool {
	type_ := p.PeekTypeAt(tokenIndex)
	switch type_ {
	case TokenType_IDENTIFIER, TokenType_STRING, TokenType_NUMBER, TokenType_BIGINT:
		return true
	default:
		return Keywords_IsKeywordTokenType(type_)
	}
}

// Java: Parser.parseObjectLiteralPropertyAssignment()
func (p *Parser) ParseObjectLiteralPropertyAssignment() ParseTree {
	type_ := p.PeekType()
	if type_ == TokenType_STAR {
		return p.ParseObjectLiteralPropertyAssignmentGenerator()
	} else if type_ == TokenType_ELLIPSIS {
		p.RecordFeatureUsed(FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_ELLIPSIS)
		operand := p.ParseAssignmentExpression()
		return NewObjectSpreadTree(p.GetTreeLocation(start), operand)
	} else if type_ == TokenType_STRING ||
		type_ == TokenType_NUMBER ||
		type_ == TokenType_BIGINT ||
		type_ == TokenType_IDENTIFIER ||
		Keywords_IsKeywordTokenType(type_) {
		if p.PeekGetAccessor() {
			return p.ParseObjectLiteralGetAccessor()
		} else if p.PeekSetAccessor() {
			return p.ParseObjectLiteralSetAccessor()
		} else if p.PeekAsyncMethod() {
			return p.ParseObjectLiteralAsyncMethod()
		} else if p.PeekTypeAt(1) == TokenType_OPEN_PAREN {
			return p.ParseObjectLiteralMethodDeclaration()
		} else {
			return p.ParseObjectLiteralPropertyNameAssignment()
		}
	} else if type_ == TokenType_OPEN_SQUARE {
		start := p.GetTreeStartLocation()
		name := p.ParseComputedPropertyName()

		if p.Peek(TokenType_COLON) {
			p.Eat(TokenType_COLON)
			value := p.ParseAssignmentExpression()
			return NewComputedPropertyDefinitionTree(p.GetTreeLocation(start), name, value)
		} else {
			builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION)
			p.ParseFunctionTail(builder, Parser_FunctionFlavor_NORMAL)
			var value ParseTree = builder.Build(p.GetTreeLocation(start))
			return NewComputedPropertyMethodTree(p.GetTreeLocation(start), name, value)
		}
	} else {
		panic(jutil.NewRuntimeException("unreachable"))
	}
}

// Java: Parser.parseObjectLiteralPropertyAssignmentGenerator()
func (p *Parser) ParseObjectLiteralPropertyAssignmentGenerator() ParseTree {
	type_ := p.PeekTypeAt(1)
	if type_ == TokenType_STRING ||
		type_ == TokenType_NUMBER ||
		type_ == TokenType_IDENTIFIER ||
		Keywords_IsKeywordTokenType(type_) {
		// parseMethodDeclaration will consume the '*'.
		return p.ParseObjectLiteralMethodDeclaration()
	} else {
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_STAR)

		name := p.ParseComputedPropertyName()
		builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_EXPRESSION)

		p.ParseFunctionTail(builder, Parser_FunctionFlavor_GENERATOR)
		var value ParseTree = builder.Build(p.GetTreeLocation(start))
		return NewComputedPropertyMethodTree(p.GetTreeLocation(start), name, value)
	}
}

// Java: Parser.parseComputedPropertyName()
func (p *Parser) ParseComputedPropertyName() ParseTree {

	p.Eat(TokenType_OPEN_SQUARE)
	assign := p.ParseAssignmentExpression()
	p.Eat(TokenType_CLOSE_SQUARE)
	return assign
}

// Java: Parser.peekGetAccessor()
func (p *Parser) PeekGetAccessor() bool {
	return p.PeekPredefinedString(PredefinedName_GET) && p.PeekPropertyNameOrComputedProp(1)
}

// Java: Parser.peekPredefinedString(String)
func (p *Parser) PeekPredefinedString(string_ string) bool {
	return p.PeekPredefinedStringAt(0, string_)
}

// EatPredefinedString consumes an identifier with the given value, or
// reports an error and returns nil.
//
// Java: Parser.eatPredefinedString(String)
func (p *Parser) EatPredefinedString(string_ string) Token {
	token := identifierTokenAsToken(p.EatId())
	if token == nil || !token.(*IdentifierToken).ValueEquals(string_) {
		p.ReportExpectedError(token, string_)
		return nil
	}
	return token
}

// Java: Parser.peekPredefinedString(int, String)
func (p *Parser) PeekPredefinedStringAt(index int, string_ string) bool {
	return p.PeekAt(index, TokenType_IDENTIFIER) &&
		p.PeekTokenAt(index).(*IdentifierToken).ValueEquals(string_)
}

// Java: Parser.parseObjectLiteralGetAccessor()
func (p *Parser) ParseObjectLiteralGetAccessor() ParseTree {
	return p.ParseGetAccessor(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseGetAccessor(ClassOrObjectElementInfo)
func (p *Parser) ParseGetAccessor(elementInfo *ClassOrObjectElementInfo) ParseTree {
	p.EatPredefinedString(PredefinedName_GET)

	if p.PeekPropertyName(0) {
		propertyName := p.EatObjectLiteralPropertyName()
		p.Eat(TokenType_OPEN_PAREN)
		p.Eat(TokenType_CLOSE_PAREN)
		body := p.ParseFunctionBody()
		p.RecordFeatureUsed(FeatureSet_Feature_GETTER)
		return NewGetAccessorTree(
			p.GetTreeLocation(elementInfo.start),
			propertyName,
			elementInfo.isClassMember,
			elementInfo.isStatic,
			body)
	} else {
		property := p.ParseComputedPropertyName()
		p.Eat(TokenType_OPEN_PAREN)
		p.Eat(TokenType_CLOSE_PAREN)
		body := p.ParseFunctionBody()
		p.RecordFeatureUsed(FeatureSet_Feature_GETTER)
		return NewComputedPropertyGetterTree(
			p.GetTreeLocation(elementInfo.start), property, elementInfo.isStatic, body)
	}
}

// Java: Parser.peekSetAccessor()
func (p *Parser) PeekSetAccessor() bool {
	return p.PeekPredefinedString(PredefinedName_SET) && p.PeekPropertyNameOrComputedProp(1)
}

// Java: Parser.parseObjectLiteralSetAccessor()
func (p *Parser) ParseObjectLiteralSetAccessor() ParseTree {
	return p.ParseSetAccessor(p.CreateObjectLiteralElementInfo())
}

// Java: Parser.parseSetAccessor(ClassOrObjectElementInfo)
func (p *Parser) ParseSetAccessor(elementInfo *ClassOrObjectElementInfo) ParseTree {
	p.EatPredefinedString(PredefinedName_SET)
	if p.PeekPropertyName(0) {
		propertyName := p.EatObjectLiteralPropertyName()
		parameter := p.ParseSetterParameterList()

		body := p.ParseFunctionBody()

		p.RecordFeatureUsed(FeatureSet_Feature_SETTER)
		return NewSetAccessorTree(
			p.GetTreeLocation(elementInfo.start),
			propertyName,
			elementInfo.isClassMember,
			elementInfo.isStatic,
			parameter,
			body)
	} else {
		property := p.ParseComputedPropertyName()
		parameter := p.ParseSetterParameterList()
		body := p.ParseFunctionBody()

		p.RecordFeatureUsed(FeatureSet_Feature_SETTER)
		return NewComputedPropertySetterTree(
			p.GetTreeLocation(elementInfo.start), property, elementInfo.isStatic, parameter, body)
	}
}

// Java: Parser.parseObjectLiteralPropertyNameAssignment()
func (p *Parser) ParseObjectLiteralPropertyNameAssignment() ParseTree {
	start := p.GetTreeStartLocation()
	name := p.EatObjectLiteralPropertyName()
	colon := p.EatOpt(TokenType_COLON)
	if colon == nil {
		if name.GetTokenBase().Type != TokenType_IDENTIFIER {
			p.ReportExpectedError(p.PeekToken(), TokenType_COLON)
		} else if name.(*IdentifierToken).IsKeyword() {
			p.ReportErrorAtToken(name, "Cannot use keyword in short object literal")
		} else if p.Peek(TokenType_EQUAL) {
			idTree := NewIdentifierExpressionTree(p.GetTreeLocation(start), name.(*IdentifierToken))
			p.Eat(TokenType_EQUAL)
			defaultValue := p.ParseAssignmentExpression()
			return NewDefaultParameterTree(p.GetTreeLocation(start), idTree, defaultValue)
		}
	}
	var value ParseTree
	if colon != nil {
		value = p.ParseAssignmentExpression()
	}
	return NewPropertyNameAssignmentTree(p.GetTreeLocation(start), name, value)
}

// 12.2 Primary Expression
//   CoverParenthesizedExpressionAndArrowParameterList ::=
//     ( Expression )
//     ( Expression, )
//     ( )
//     ( ... BindingIdentifier )
//     ( Expression , ... BindingIdentifier )

// Java: Parser.parseCoverParenthesizedExpressionAndArrowParameterList()
func (p *Parser) ParseCoverParenthesizedExpressionAndArrowParameterList() ParseTree {
	if p.PeekTypeAt(1) == TokenType_FOR {
		return p.ParseGeneratorComprehension()
	}

	start := p.GetTreeStartLocation()
	p.Eat(TokenType_OPEN_PAREN)
	// Case ( )
	if p.Peek(TokenType_CLOSE_PAREN) {
		p.Eat(TokenType_CLOSE_PAREN)
		if p.Peek(TokenType_ARROW) {
			return NewFormalParameterListTree(
				p.GetTreeLocation(start),
				[]ParseTree{},
				/* hasTrailingComma= */ false,
				[]SourcePosition{})
		} else {
			p.ReportError("invalid parenthesized expression")
			return NewMissingPrimaryExpressionTree(p.GetTreeLocation(start))
		}
	}
	// Case ( ... BindingIdentifier )
	if p.Peek(TokenType_ELLIPSIS) {
		params := []ParseTree{p.ParseParameter()}
		p.Eat(TokenType_CLOSE_PAREN)
		if p.Peek(TokenType_ARROW) {
			return NewFormalParameterListTree(
				p.GetTreeLocation(start),
				params,
				/* hasTrailingComma= */ false,
				[]SourcePosition{})
		} else {
			p.ReportError("invalid parenthesized expression")
			return NewMissingPrimaryExpressionTree(p.GetTreeLocation(start))
		}
	}
	// For either of the three remaining cases:
	//     ( Expression )
	//     ( Expression, )
	//     ( Expression, ...BindingIdentifier )
	// we can parse as an expression.
	result := p.ParseExpression()
	// If it follows with a comma, we must be in either of two cases
	//     ( Expression, )
	//     ( Expression, ...BindingIdentifier )
	// case.
	if p.Peek(TokenType_COMMA) {
		if p.PeekAt(1, TokenType_CLOSE_PAREN) {
			// Create the formal parameter list here so we can record
			// the trailing comma
			p.ResetScanner(start)
			// If we fail to parse as an ArrowFunction parameter list then
			// parseFormalParameterList will take care of reporting errors.
			return p.ParseFormalParameterList()
		} else {
			p.Eat(TokenType_COMMA)
			// Since we already parsed as an expression, we will guaranteed reparse this expression
			// as an arrow function parameter list, but just leave it as a comma expression for now.
			result = NewCommaExpressionTree(
				p.GetTreeLocation(start), []ParseTree{result, p.ParseParameter()})
		}
	}
	p.Eat(TokenType_CLOSE_PAREN)
	return NewParenExpressionTree(p.GetTreeLocation(start), result)
}

// Java: Parser.parseMissingPrimaryExpression()
func (p *Parser) ParseMissingPrimaryExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.NextToken()
	p.ReportError("primary expression expected")
	return NewMissingPrimaryExpressionTree(p.GetTreeLocation(start))
}

// Parser_Expression is Parser.Expression: differentiates between parsing for
// 'In' vs. 'NoIn' Variants of expression grammars.
type Parser_Expression int32

const (
	Parser_Expression_NO_IN Parser_Expression = iota
	Parser_Expression_NORMAL
)

var parser_Expression_names = [...]string{"NO_IN", "NORMAL"}

// Parser_Expression_values is Expression.values().
//
// Java: Parser.Expression.values()
func Parser_Expression_values() []Parser_Expression {
	return []Parser_Expression{Parser_Expression_NO_IN, Parser_Expression_NORMAL}
}

// Name is Enum.name().
//
// Java: Parser.Expression.name()
func (e Parser_Expression) Name() string { return parser_Expression_names[e] }

// String is Enum.toString().
//
// Java: Parser.Expression.toString()
func (e Parser_Expression) String() string { return e.Name() }

// 11.14 Expressions

// Java: Parser.parseExpressionNoIn()
func (p *Parser) ParseExpressionNoIn() ParseTree {
	return p.Parse(Parser_Expression_NO_IN)
}

// Java: Parser.parseExpression()
func (p *Parser) ParseExpression() ParseTree {
	return p.Parse(Parser_Expression_NORMAL)
}

// Java: Parser.peekExpression()
func (p *Parser) PeekExpression() bool {
	switch p.PeekType() {
	case TokenType_BANG,
		TokenType_CLASS,
		TokenType_DELETE,
		TokenType_FALSE,
		TokenType_FUNCTION,
		TokenType_IDENTIFIER,
		TokenType_TYPE,
		TokenType_DECLARE,
		TokenType_MODULE,
		TokenType_NAMESPACE,
		TokenType_MINUS,
		TokenType_MINUS_MINUS,
		TokenType_NEW,
		TokenType_NULL,
		TokenType_NUMBER,
		TokenType_BIGINT,
		TokenType_OPEN_CURLY,
		TokenType_OPEN_PAREN,
		TokenType_OPEN_SQUARE,
		TokenType_PLUS,
		TokenType_PLUS_PLUS,
		TokenType_SLASH,       // regular expression literal
		TokenType_SLASH_EQUAL, // regular expression literal
		TokenType_STRING,
		TokenType_NO_SUBSTITUTION_TEMPLATE,
		TokenType_TEMPLATE_HEAD,
		TokenType_SUPER,
		TokenType_THIS,
		TokenType_TILDE,
		TokenType_TRUE,
		TokenType_TYPEOF,
		TokenType_VOID,
		TokenType_YIELD:
		return true
	case TokenType_IMPORT:
		return p.PeekImportCall() || p.PeekImportDot()
	default:
		return false
	}
}

// Java: Parser.parse(Expression)
func (p *Parser) Parse(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	result := p.ParseAssignment(expressionIn)
	if p.Peek(TokenType_COMMA) && !p.PeekAt(1, TokenType_ELLIPSIS) && !p.PeekAt(1, TokenType_CLOSE_PAREN) {
		exprs := []ParseTree{}
		exprs = append(exprs, result)
		for p.Peek(TokenType_COMMA) &&
			!p.PeekAt(1, TokenType_ELLIPSIS) &&
			!p.PeekAt(1, TokenType_CLOSE_PAREN) {
			p.Eat(TokenType_COMMA)
			exprs = append(exprs, p.ParseAssignment(expressionIn))
		}
		return NewCommaExpressionTree(p.GetTreeLocation(start), exprs)
	}
	return result
}

// 12.14 Assignment operators

// Java: Parser.parseAssignmentExpression()
func (p *Parser) ParseAssignmentExpression() ParseTree {
	return p.ParseAssignment(Parser_Expression_NORMAL)
}

// Java: Parser.peekAssignmentExpression()
func (p *Parser) PeekAssignmentExpression() bool {
	return p.PeekExpression()
}

// Java: Parser.parseAssignment(Expression)
func (p *Parser) ParseAssignment(expressionIn Parser_Expression) ParseTree {
	if p.Peek(TokenType_YIELD) && p.InGeneratorContext() {
		return p.ParseYield(expressionIn)
	}

	start := p.GetTreeStartLocation()
	left := p.ParseConditional(expressionIn)

	if p.IsStartOfAsyncArrowFunction(left) {
		// re-evaluate as an async arrow function.
		p.ResetScannerToTree(left)
		return p.ParseAsyncArrowFunction(expressionIn)
	}
	if p.Peek(TokenType_ARROW) {
		return p.CompleteAssignmentExpressionParseAtArrowWithExpression(left, expressionIn)
	}

	if p.PeekAssignmentOperator() {
		if !p.Peek(TokenType_EQUAL) {
			// not the vanilla assignment operator `=`, but a special equals operator (`+=`, `-=`,
			// `**=`, etc)
			if !left.IsValidNonVanillaAssignmentTarget() {
				p.ReportError("invalid assignment target")
				return NewMissingPrimaryExpressionTree(p.GetTreeLocation(p.GetTreeStartLocation()))
			}
		}
		left = p.TransformLeftHandSideExpression(left)
		if !left.IsValidAssignmentTarget() {
			p.ReportError("invalid assignment target")
			return NewMissingPrimaryExpressionTree(p.GetTreeLocation(p.GetTreeStartLocation()))
		}
		operator := p.NextToken()
		right := p.ParseAssignment(expressionIn)
		return NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.isStartOfAsyncArrowFunction(ParseTree)
func (p *Parser) IsStartOfAsyncArrowFunction(partialExpression ParseTree) bool {
	if partialExpression.GetParseTreeBase().Type == ParseTreeType_IDENTIFIER_EXPRESSION {
		identifierToken := partialExpression.(*IdentifierExpressionTree).IdentifierToken
		// partialExpression is `async`
		// followed by `[no newline] bindingIdentifier [no newline] =>`
		return identifierToken.ValueEquals(Parser_ASYNC) &&
			!p.PeekImplicitSemiColonAt(0) &&
			p.PeekId() &&
			!p.PeekImplicitSemiColonAt(1) &&
			p.PeekAt(1, TokenType_ARROW)
	} else if partialExpression.GetParseTreeBase().Type == ParseTreeType_CALL_EXPRESSION {
		callExpression := partialExpression.(*CallExpressionTree)
		callee := callExpression.Operand
		var arguments ParseTree = callExpression.Arguments
		// partialExpression is `async [no newline] (parameters)`
		// followed by `[no newline] =>`
		return callee.GetParseTreeBase().Type == ParseTreeType_IDENTIFIER_EXPRESSION &&
			callee.(*IdentifierExpressionTree).IdentifierToken.ValueEquals(Parser_ASYNC) &&
			callee.GetParseTreeBase().Location.End.Line == arguments.GetParseTreeBase().Location.Start.Line &&
			!p.PeekImplicitSemiColon() &&
			p.Peek(TokenType_ARROW)
	} else {
		return false
	}
}

// Java: Parser.completeAssignmentExpressionParseAtArrow(ParseTree, Expression)
func (p *Parser) CompleteAssignmentExpressionParseAtArrowWithExpression(leftOfArrow ParseTree, expressionIn Parser_Expression) ParseTree {
	if leftOfArrow.GetParseTreeBase().Type == ParseTreeType_CALL_EXPRESSION {
		//   ... someAssignmentExpression // implicit semicolon
		//   (args) =>
		return p.CompleteAssignmentExpressionParseAtArrow(leftOfArrow.(*CallExpressionTree))
	} else {
		return p.CompleteArrowFunctionParseAtArrow(leftOfArrow, expressionIn)
	}
}

// Java: Parser.completeArrowFunctionParseAtArrow(ParseTree, Expression)
func (p *Parser) CompleteArrowFunctionParseAtArrow(leftOfArrow ParseTree, expressionIn Parser_Expression) ParseTree {
	arrowFormalParameters := p.TransformToArrowFormalParameters(leftOfArrow)
	if p.PeekImplicitSemiColon() {
		p.ReportError("No newline allowed before '=>'")
	}
	p.Eat(TokenType_ARROW)
	arrowFunctionBody := p.ParseArrowFunctionBody(expressionIn, Parser_FunctionFlavor_NORMAL)

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_ARROW).
		SetFormalParameterList(arrowFormalParameters).
		SetFunctionBody(arrowFunctionBody)
	return builder.Build(p.GetTreeLocation(arrowFormalParameters.Location.Start))
}

// Java: Parser.transformToArrowFormalParameters(ParseTree)
func (p *Parser) TransformToArrowFormalParameters(leftOfArrow ParseTree) *FormalParameterListTree {
	var arrowParameterList *FormalParameterListTree
	switch leftOfArrow.GetParseTreeBase().Type {
	case ParseTreeType_FORMAL_PARAMETER_LIST:
		arrowParameterList = leftOfArrow.(*FormalParameterListTree)
	case ParseTreeType_IDENTIFIER_EXPRESSION:
		// e.g. x => x + 1
		arrowParameterList = NewFormalParameterListTree(
			leftOfArrow.GetParseTreeBase().Location,
			[]ParseTree{leftOfArrow},
			/* hasTrailingComma= */ false,
			[]SourcePosition{})
	case ParseTreeType_ARGUMENT_LIST, ParseTreeType_PAREN_EXPRESSION:
		// e.g. (x) => x + 1
		p.ResetScannerToTree(leftOfArrow)
		// If we fail to parse as an ArrowFunction parameter list then
		// parseFormalParameterList will take care of reporting errors.
		arrowParameterList = p.ParseFormalParameterList()
	default:
		p.ReportErrorAtTree(leftOfArrow, "invalid arrow function parameters")
		arrowParameterList = Parser_NewEmptyFormalParameterList(leftOfArrow.GetParseTreeBase().Location)
	}
	return arrowParameterList
}

// Java: Parser.completeAssignmentExpressionParseAtArrow(CallExpressionTree)
func (p *Parser) CompleteAssignmentExpressionParseAtArrow(callExpression *CallExpressionTree) ParseTree {
	operand := callExpression.Operand
	var arguments ParseTree = callExpression.Arguments
	var result ParseTree
	if operand.GetParseTreeBase().Location.End.Line < arguments.GetParseTreeBase().Location.Start.Line {
		// break at the implicit semicolon
		// Example:
		// foo.bar // operand and implicit semicolon
		// () => { doSomething; };
		p.ResetScannerAfter(operand)
		result = operand
	} else {
		p.ReportError("'=>' unexpected")
		result = callExpression
	}
	return result
}

// Java: Parser.parseAsyncArrowFunction(Expression)
func (p *Parser) ParseAsyncArrowFunction(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	p.EatPredefinedString(Parser_ASYNC)
	if p.PeekImplicitSemiColon() {
		p.ReportError("No newline allowed between `async` and arrow function parameter list")
	}
	var arrowParameterList *FormalParameterListTree
	if p.Peek(TokenType_OPEN_PAREN) {
		// async (...) =>
		arrowParameterList = p.ParseFormalParameterList()
	} else {
		// async arg =>
		singleParameter := p.ParseIdentifierExpression()
		arrowParameterList = NewFormalParameterListTree(
			singleParameter.Location,
			[]ParseTree{singleParameter},
			/* hasTrailingComma= */ false,
			[]SourcePosition{})
	}
	if p.PeekImplicitSemiColon() {
		p.ReportError("No newline allowed before '=>'")
	}
	p.Eat(TokenType_ARROW)
	arrowFunctionBody := p.ParseArrowFunctionBody(expressionIn, Parser_FunctionFlavor_ASYNCHRONOUS)

	builder := FunctionDeclarationTree_NewBuilder(FunctionDeclarationTree_Kind_ARROW).
		SetAsync(true).
		SetFormalParameterList(arrowParameterList).
		SetFunctionBody(arrowFunctionBody)
	return builder.Build(p.GetTreeLocation(start))
}

// Java: Parser.parseArrowFunctionBody(Expression, FunctionFlavor)
func (p *Parser) ParseArrowFunctionBody(expressionIn Parser_Expression, functionFlavor Parser_FunctionFlavor) ParseTree {
	p.functionContextStack = append(p.functionContextStack, functionFlavor)
	var arrowFunctionBody ParseTree
	if p.Peek(TokenType_OPEN_CURLY) {
		arrowFunctionBody = p.ParseFunctionBody()
	} else {
		arrowFunctionBody = p.ParseAssignment(expressionIn)
	}
	p.functionContextStack = p.functionContextStack[:len(p.functionContextStack)-1]
	return arrowFunctionBody
}

// Java: Parser.newEmptyFormalParameterList(SourceRange)
func Parser_NewEmptyFormalParameterList(location SourceRange) *FormalParameterListTree {
	return NewFormalParameterListTree(
		location,
		[]ParseTree{},
		/* hasTrailingComma= */ false,
		[]SourcePosition{})
}

// TransformLeftHandSideExpression transforms a LeftHandSideExpression into
// a LeftHandSidePattern if possible. This returns the transformed tree if it
// parses as a LeftHandSidePattern, otherwise it returns the original tree.
//
// Java: Parser.transformLeftHandSideExpression(ParseTree)
func (p *Parser) TransformLeftHandSideExpression(tree ParseTree) ParseTree {
	switch tree.GetParseTreeBase().Type {
	case ParseTreeType_ARRAY_LITERAL_EXPRESSION, ParseTreeType_OBJECT_LITERAL_EXPRESSION:
		p.ResetScannerToTree(tree)
		// If we fail to parse as an LeftHandSidePattern then
		// parseLeftHandSidePattern will take care reporting errors.
		return p.ParseLeftHandSidePattern()
	default:
		return tree
	}
}

// Java: Parser.parseLeftHandSidePattern()
func (p *Parser) ParseLeftHandSidePattern() ParseTree {
	return p.ParsePattern(Parser_PatternKind_ANY)
}

// Java: Parser.resetScanner(SourcePosition)
func (p *Parser) ResetScanner(start SourcePosition) {
	// TODO(bradfordcsmith): lastSourcePosition should really point to the end of the last token
	//     before the tree to correctly detect implicit semicolons, but it doesn't matter for the
	//     current use case.
	p.lastSourcePosition = start
	p.scanner.SetPosition(p.lastSourcePosition)
}

// Java: Parser.resetScanner(ParseTree)
func (p *Parser) ResetScannerToTree(tree ParseTree) {
	p.scanner.SetPosition(tree.GetParseTreeBase().Location.Start)
}

// Java: Parser.resetScannerAfter(ParseTree)
func (p *Parser) ResetScannerAfter(parseTree ParseTree) {
	p.lastSourcePosition = parseTree.GetParseTreeBase().Location.End
	// NOTE: The "end" position for a parseTree actually points to the first character after the
	//     last token in the tree, so this is not an off-by-one error.
	p.scanner.SetPosition(p.lastSourcePosition)
}

// Java: Parser.peekAssignmentOperator()
func (p *Parser) PeekAssignmentOperator() bool {
	switch p.PeekType() {
	case TokenType_EQUAL,
		TokenType_STAR_EQUAL,
		TokenType_STAR_STAR_EQUAL,
		TokenType_SLASH_EQUAL,
		TokenType_PERCENT_EQUAL,
		TokenType_PLUS_EQUAL,
		TokenType_MINUS_EQUAL,
		TokenType_LEFT_SHIFT_EQUAL,
		TokenType_RIGHT_SHIFT_EQUAL,
		TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL,
		TokenType_AMPERSAND_EQUAL,
		TokenType_CARET_EQUAL,
		TokenType_BAR_EQUAL,
		TokenType_OR_EQUAL,
		TokenType_AND_EQUAL,
		TokenType_QUESTION_QUESTION_EQUAL:
		return true
	default:
		return false
	}
}

// Java: Parser.inGeneratorContext()
func (p *Parser) InGeneratorContext() bool {
	// disallow yield outside of generators
	return p.functionContextStack[len(p.functionContextStack)-1].IsGenerator()
}

// ParseYield parses yield [no line terminator] (*)? AssignExpression
// https://people.mozilla.org/~jorendorff/es6-draft.html#sec-generator-function-definitions-runtime-semantics-evaluation
//
// Java: Parser.parseYield(Expression)
func (p *Parser) ParseYield(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_YIELD)
	isYieldAll := false
	var expression ParseTree
	if !p.PeekImplicitSemiColon() {
		isYieldAll = p.EatOpt(TokenType_STAR) != nil
		if p.PeekAssignmentExpression() {
			expression = p.ParseAssignment(expressionIn)
		} else if isYieldAll {
			p.ReportError("yield* requires an expression")
		}
	}
	return NewYieldExpressionTree(p.GetTreeLocation(start), isYieldAll, expression)
}

// 11.12 Conditional Expression

// Java: Parser.parseConditional(Expression)
func (p *Parser) ParseConditional(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	condition := p.ParseShortCircuit(expressionIn)
	if p.Peek(TokenType_QUESTION) {
		p.Eat(TokenType_QUESTION)
		left := p.ParseAssignment(expressionIn)
		p.Eat(TokenType_COLON)
		right := p.ParseAssignment(expressionIn)
		return NewConditionalExpressionTree(p.GetTreeLocation(start), condition, left, right)
	}
	return condition
}

// Java: Parser.parseShortCircuit(Expression)
func (p *Parser) ParseShortCircuit(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseLogicalOR(expressionIn)
	if p.Peek(TokenType_QUESTION_QUESTION) {
		if left.GetParseTreeBase().Type == ParseTreeType_BINARY_OPERATOR {
			binaryTree := left.(*BinaryOperatorTree)
			if binaryTree.Operator.GetTokenBase().Type == TokenType_AND || binaryTree.Operator.GetTokenBase().Type == TokenType_OR {
				p.ReportError("Logical OR and logical AND require parentheses when used with '??'")
			}
		}
		return p.ParseNullishCoalesce(expressionIn, left, start)
	} else {
		return left
	}
}

// Java: Parser.parseNullishCoalesce(Expression, ParseTree, SourcePosition)
func (p *Parser) ParseNullishCoalesce(expressionIn Parser_Expression, left ParseTree, start SourcePosition) ParseTree {
	for p.Peek(TokenType_QUESTION_QUESTION) {
		operator := p.Eat(TokenType_QUESTION_QUESTION)
		right := p.ParseBitwiseOR(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	if p.Peek(TokenType_AND) || p.Peek(TokenType_OR) {
		p.ReportError("Logical OR and logical AND require parentheses when used with '??'")
	}
	return left
}

// 11.11 Logical OR

// Java: Parser.parseLogicalOR(Expression)
func (p *Parser) ParseLogicalOR(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseLogicalAND(expressionIn)
	for p.Peek(TokenType_OR) {
		operator := p.Eat(TokenType_OR)
		right := p.ParseLogicalAND(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.11 Logical AND

// Java: Parser.parseLogicalAND(Expression)
func (p *Parser) ParseLogicalAND(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseBitwiseOR(expressionIn)
	for p.Peek(TokenType_AND) {
		operator := p.Eat(TokenType_AND)
		right := p.ParseBitwiseOR(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.10 Bitwise OR

// Java: Parser.parseBitwiseOR(Expression)
func (p *Parser) ParseBitwiseOR(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseBitwiseXOR(expressionIn)
	for p.Peek(TokenType_BAR) {
		operator := p.Eat(TokenType_BAR)
		right := p.ParseBitwiseXOR(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.10 Bitwise XOR

// Java: Parser.parseBitwiseXOR(Expression)
func (p *Parser) ParseBitwiseXOR(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseBitwiseAND(expressionIn)
	for p.Peek(TokenType_CARET) {
		operator := p.Eat(TokenType_CARET)
		right := p.ParseBitwiseAND(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.10 Bitwise AND

// Java: Parser.parseBitwiseAND(Expression)
func (p *Parser) ParseBitwiseAND(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseEquality(expressionIn)
	for p.Peek(TokenType_AMPERSAND) {
		operator := p.Eat(TokenType_AMPERSAND)
		right := p.ParseEquality(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// 11.9 Equality Expression

// Java: Parser.parseEquality(Expression)
func (p *Parser) ParseEquality(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseRelational(expressionIn)
	for p.PeekEqualityOperator() {
		operator := p.NextToken()
		right := p.ParseRelational(expressionIn)
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekEqualityOperator()
func (p *Parser) PeekEqualityOperator() bool {
	switch p.PeekType() {
	case TokenType_EQUAL_EQUAL, TokenType_NOT_EQUAL, TokenType_EQUAL_EQUAL_EQUAL, TokenType_NOT_EQUAL_EQUAL:
		return true
	default:
		return false
	}
}

// 11.8 Relational

// Java: Parser.parseRelational(Expression)
func (p *Parser) ParseRelational(expressionIn Parser_Expression) ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseShiftExpression()
	for p.PeekRelationalOperator(expressionIn) {
		operator := p.NextToken()
		right := p.ParseShiftExpression()
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekRelationalOperator(Expression)
func (p *Parser) PeekRelationalOperator(expressionIn Parser_Expression) bool {
	switch p.PeekType() {
	case TokenType_OPEN_ANGLE, TokenType_CLOSE_ANGLE, TokenType_GREATER_EQUAL, TokenType_LESS_EQUAL, TokenType_INSTANCEOF:
		return true
	case TokenType_IN:
		return expressionIn == Parser_Expression_NORMAL
	default:
		return false
	}
}

// 11.7 Shift Expression

// Java: Parser.parseShiftExpression()
func (p *Parser) ParseShiftExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseAdditiveExpression()
	for p.PeekShiftOperator() {
		operator := p.NextToken()
		right := p.ParseAdditiveExpression()
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekShiftOperator()
func (p *Parser) PeekShiftOperator() bool {
	switch p.PeekType() {
	case TokenType_LEFT_SHIFT, TokenType_RIGHT_SHIFT, TokenType_UNSIGNED_RIGHT_SHIFT:
		return true
	default:
		return false
	}
}

// 11.6 Additive Expression

// Java: Parser.parseAdditiveExpression()
func (p *Parser) ParseAdditiveExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseMultiplicativeExpression()
	for p.PeekAdditiveOperator() {
		operator := p.NextToken()
		right := p.ParseMultiplicativeExpression()
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekAdditiveOperator()
func (p *Parser) PeekAdditiveOperator() bool {
	switch p.PeekType() {
	case TokenType_PLUS, TokenType_MINUS:
		return true
	default:
		return false
	}
}

// 11.5 Multiplicative Expression

// Java: Parser.parseMultiplicativeExpression()
func (p *Parser) ParseMultiplicativeExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseExponentiationExpression()
	for p.PeekMultiplicativeOperator() {
		operator := p.NextToken()
		right := p.ParseExponentiationExpression()
		left = NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	}
	return left
}

// Java: Parser.peekMultiplicativeOperator()
func (p *Parser) PeekMultiplicativeOperator() bool {
	switch p.PeekType() {
	case TokenType_STAR, TokenType_SLASH, TokenType_PERCENT:
		return true
	default:
		return false
	}
}

// Java: Parser.parseExponentiationExpression()
func (p *Parser) ParseExponentiationExpression() ParseTree {
	start := p.GetTreeStartLocation()
	left := p.ParseUnaryExpression()
	if p.Peek(TokenType_STAR_STAR) {
		// ExponentiationExpression does not allow a UnaryExpression before '**'.
		// Parentheses are required to disambiguate:
		//   (-x)**y is valid
		//   -(x**y) is valid
		//   -x**y is a syntax error
		if left.GetParseTreeBase().Type == ParseTreeType_UNARY_EXPRESSION {
			p.ReportError(
				"Unary operator '%s' requires parentheses before '**'",
				left.(*UnaryExpressionTree).Operator)
		} else if left.GetParseTreeBase().Type == ParseTreeType_AWAIT_EXPRESSION {
			p.ReportError("Unary operator 'await' requires parentheses before '**'")
		}
		operator := p.NextToken()
		right := p.ParseExponentiationExpression()
		return NewBinaryOperatorTree(p.GetTreeLocation(start), left, operator, right)
	} else {
		return left
	}
}

// 11.4 Unary Operator

// Java: Parser.parseUnaryExpression()
func (p *Parser) ParseUnaryExpression() ParseTree {
	start := p.GetTreeStartLocation()
	if p.PeekUnaryOperator() {
		operator := p.NextToken()
		operand := p.ParseUnaryExpression()
		return NewUnaryExpressionTree(p.GetTreeLocation(start), operator, operand)
	} else if p.PeekAwaitExpression() {
		return p.ParseAwaitExpression()
	} else {
		return p.ParseUpdateExpression()
	}
}

// Java: Parser.peekUnaryOperator()
func (p *Parser) PeekUnaryOperator() bool {
	switch p.PeekType() {
	case TokenType_DELETE, TokenType_VOID, TokenType_TYPEOF, TokenType_PLUS, TokenType_MINUS, TokenType_TILDE, TokenType_BANG:
		return true
	default:
		return false
	}
}

// Parser_AWAIT is Parser.AWAIT.
const Parser_AWAIT = "await"

// Java: Parser.peekAwaitExpression()
func (p *Parser) PeekAwaitExpression() bool {
	return p.PeekPredefinedString(Parser_AWAIT)
}

// Java: Parser.parseAwaitExpression()
func (p *Parser) ParseAwaitExpression() ParseTree {
	start := p.GetTreeStartLocation()
	p.EatPredefinedString(Parser_AWAIT)
	expression := p.ParseUnaryExpression()
	return NewAwaitExpressionTree(p.GetTreeLocation(start), expression)
}

// Java: Parser.parseUpdateExpression()
func (p *Parser) ParseUpdateExpression() ParseTree {
	start := p.GetTreeStartLocation()
	if p.PeekUpdateOperator() {
		operator := p.NextToken()
		operand := p.ParseUnaryExpression()
		return UpdateExpressionTree_Prefix(p.GetTreeLocation(start), operator, operand)
	} else {
		lhs := p.ParseLeftHandSideExpression()
		if p.PeekUpdateOperator() && !p.PeekImplicitSemiColon() {
			// newline not allowed before an update operator.
			operator := p.NextToken()
			return UpdateExpressionTree_Postfix(p.GetTreeLocation(start), operator, lhs)
		} else {
			return lhs
		}
	}
}

// Java: Parser.peekUpdateOperator()
func (p *Parser) PeekUpdateOperator() bool {
	switch p.PeekType() {
	case TokenType_PLUS_PLUS, TokenType_MINUS_MINUS:
		return true
	default:
		return false
	}
}

// Java: Parser.peekImportCall()
func (p *Parser) PeekImportCall() bool {
	return p.Peek(TokenType_IMPORT) && p.PeekAt(1, TokenType_OPEN_PAREN)
}

// Java: Parser.peekImportDot()
func (p *Parser) PeekImportDot() bool {
	return p.Peek(TokenType_IMPORT) && p.PeekAt(1, TokenType_PERIOD)
}

// ParseLeftHandSideExpression parses a LeftHandSideExpression.
//
// Java: Parser.parseLeftHandSideExpression()
func (p *Parser) ParseLeftHandSideExpression() ParseTree {
	start := p.GetTreeStartLocation()
	// We have these possible productions.
	// LeftHandSideExpression -> NewExpression
	//                        -> CallExpression
	//                        -> MemberExpression
	//                        -> OptionalExpression
	//
	// NewExpression -> new NewExpression
	//               -> MemberExpression
	//
	// CallExpression -> MemberExpression Arguments
	//                -> CallExpression ... see below
	//
	// OptionalExpression -> MemberExpression OptionalChain
	//                    -> CallExpression OptionalChain
	//                    -> OptionalExpression OptionalChain
	//
	// We try parsing a NewExpression, here, because that will include parsing MemberExpression.
	// If what we really have is a CallExpression or OptionalExpression, then the MemberExpression
	// we get back from parseNewExpression will be the first part of it, and we'll build the
	// rest later.
	operand := p.ParseNewExpression()

	// this test is equivalent to is member expression
	if newExpressionTree, ok := operand.(*NewExpressionTree); !ok || newExpressionTree.Arguments != nil {
		// We have a MemberExpression, but it may actually be just the first part of a CallExpression
		// Attempt to gather the rest of the CallExpression, if so.
		for p.PeekCallSuffix() {
			switch p.PeekType() {
			case TokenType_OPEN_PAREN:
				arguments := p.ParseArguments()
				operand = NewCallExpressionTree(p.GetTreeLocation(start), operand, arguments)
			case TokenType_OPEN_SQUARE:
				p.Eat(TokenType_OPEN_SQUARE)
				member := p.ParseExpression()
				p.Eat(TokenType_CLOSE_SQUARE)
				operand = NewMemberLookupExpressionTree(p.GetTreeLocation(start), operand, member)
			case TokenType_PERIOD:
				p.Eat(TokenType_PERIOD)
				id := p.EatIdOrKeywordAsId()
				operand = NewMemberExpressionTree(p.GetTreeLocation(start), operand, id)
			case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
				operand = p.ParseTemplateLiteral(operand)
			default:
				panic(jutil.NewAssertionError("unexpected case: " + p.PeekType().String()))
			}
		}
		operand = p.MaybeParseOptionalExpression(operand)
	}
	return operand
}

// Java: Parser.peekCallSuffix()
func (p *Parser) PeekCallSuffix() bool {
	return p.Peek(TokenType_OPEN_PAREN) ||
		p.Peek(TokenType_OPEN_SQUARE) ||
		p.Peek(TokenType_PERIOD) ||
		p.Peek(TokenType_NO_SUBSTITUTION_TEMPLATE) ||
		p.Peek(TokenType_TEMPLATE_HEAD)
}

// MaybeParseOptionalExpression tries to parse the expression as an optional
// expression.
//
// `operand?.identifier` or `operand?.[expression]` or `operand?.(arg1, arg2)`
//
// returns parse tree after trying to parse it as an optional expression
//
// Java: Parser.maybeParseOptionalExpression(ParseTree)
func (p *Parser) MaybeParseOptionalExpression(operand ParseTree) ParseTree {
	// The optional chain's source info should cover the lhs operand also
	start := operand.GetParseTreeBase().Location.Start

	for p.Peek(TokenType_QUESTION_DOT) {
		p.Eat(TokenType_QUESTION_DOT)
		switch p.PeekType() {
		case TokenType_OPEN_PAREN:
			arguments := p.ParseArguments()
			operand = NewOptChainCallExpressionTree(
				p.GetTreeLocation(start),
				operand,
				arguments,
				/* isStartOfOptionalChain= */ true,
				arguments.HasTrailingComma)
		case TokenType_OPEN_SQUARE:
			p.Eat(TokenType_OPEN_SQUARE)
			member := p.ParseExpression()
			p.Eat(TokenType_CLOSE_SQUARE)
			operand = NewOptionalMemberLookupExpressionTree(
				p.GetTreeLocation(start), operand, member /* isStartOfOptionalChain= */, true)
		case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
			p.ReportError("template literal cannot be used within optional chaining")
		default:
			if p.PeekIdOrKeyword() {
				id := p.EatIdOrKeywordAsId()
				operand = NewOptionalMemberExpressionTree(
					p.GetTreeLocation(start), operand, id /* isStartOfOptionalChain= */, true)
			} else {
				p.ReportError("syntax error: %s not allowed in optional chain", p.PeekType())
			}
		}
		operand = p.ParseRemainingOptionalChainSegment(operand)
	}
	return operand
}

// ParseRemainingOptionalChainSegment parses the remaining components of an
// optional chain till the current chain's end, or a new chain's start.
//
// `optionalExpression.identifier`, `optionalExpression[expression]`,
// `optionalExpression(arg1, arg2)`, or
// `optionalExpression?.optionalExpression`
//
// returns parse tree after trying to parse it as an optional chain
//
// Java: Parser.parseRemainingOptionalChainSegment(ParseTree)
func (p *Parser) ParseRemainingOptionalChainSegment(optionalExpression ParseTree) ParseTree {
	// The optional chain's source info should cover the lhs operand also
	start := optionalExpression.GetParseTreeBase().Location.Start
	for p.PeekOptionalChainSuffix() {
		if p.PeekType() == TokenType_NO_SUBSTITUTION_TEMPLATE ||
			p.PeekType() == TokenType_TEMPLATE_HEAD {
			p.ReportError("template literal cannot be used within optional chaining")
			break
		}
		switch p.PeekType() {
		case TokenType_PERIOD:
			p.Eat(TokenType_PERIOD)
			id := p.EatIdOrKeywordAsId()
			optionalExpression = NewOptionalMemberExpressionTree(
				p.GetTreeLocation(start),
				optionalExpression,
				id,
				/* isStartOfOptionalChain= */ false)
		case TokenType_OPEN_PAREN:
			arguments := p.ParseArguments()
			optionalExpression = NewOptChainCallExpressionTree(
				p.GetTreeLocation(start),
				optionalExpression,
				arguments,
				/* isStartOfOptionalChain= */ false,
				arguments.HasTrailingComma)
		case TokenType_OPEN_SQUARE:
			p.Eat(TokenType_OPEN_SQUARE)
			member := p.ParseExpression()
			p.Eat(TokenType_CLOSE_SQUARE)
			optionalExpression = NewOptionalMemberLookupExpressionTree(
				p.GetTreeLocation(start),
				optionalExpression,
				member,
				/* isStartOfOptionalChain= */ false)
		default:
			panic(jutil.NewAssertionError("unexpected case: " + p.PeekType().String()))
		}
	}
	return optionalExpression
}

// PeekOptionalChainSuffix: tokens that indicate continuation of an optional
// chain.
//
// Java: Parser.peekOptionalChainSuffix()
func (p *Parser) PeekOptionalChainSuffix() bool {
	return p.Peek(TokenType_OPEN_PAREN) || // a?.b( ...
		p.Peek(TokenType_OPEN_SQUARE) || // a?.b[ ...
		p.Peek(TokenType_PERIOD) || // a?.b. ...
		// TEMPLATE_HEAD and NO_SUBSTITUTION_TEMPLATE are actually not allowed within optional
		// chaining and leads to an early error as dictated by the spec.
		// https://tc39.es/proposal-optional-chaining/#sec-left-hand-side-expressions-static-semantics-early-errors
		p.Peek(TokenType_NO_SUBSTITUTION_TEMPLATE) || // a?.b`text`
		p.Peek(TokenType_TEMPLATE_HEAD) // a?.b`text ${substitution} text`
}

// Parser_ASYNC is Parser.ASYNC.
const Parser_ASYNC = "async"

// 11.2 Member Expression without the new production

// Java: Parser.parseMemberExpressionNoNew()
func (p *Parser) ParseMemberExpressionNoNew() ParseTree {
	start := p.GetTreeStartLocation()
	var operand ParseTree
	if p.PeekImportDot() {
		operand = p.ParseImportDotMeta()
	} else if p.PeekAsyncFunctionStart() {
		operand = p.ParseAsyncFunctionExpression()
	} else if p.PeekFunction() {
		operand = p.ParseFunctionExpression()
	} else {
		operand = p.ParsePrimaryExpression()
	}
	for p.PeekMemberExpressionSuffix() {
		switch p.PeekType() {
		case TokenType_OPEN_SQUARE:
			p.Eat(TokenType_OPEN_SQUARE)
			member := p.ParseExpression()
			p.Eat(TokenType_CLOSE_SQUARE)
			operand = NewMemberLookupExpressionTree(p.GetTreeLocation(start), operand, member)
		case TokenType_PERIOD:
			p.Eat(TokenType_PERIOD)
			id := p.EatIdOrKeywordAsId()
			operand = NewMemberExpressionTree(p.GetTreeLocation(start), operand, id)
		case TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD:
			operand = p.ParseTemplateLiteral(operand)
		default:
			panic(jutil.NewRuntimeException("unreachable"))
		}
	}
	return operand
}

// Java: Parser.peekMemberExpressionSuffix()
func (p *Parser) PeekMemberExpressionSuffix() bool {
	return p.Peek(TokenType_OPEN_SQUARE) ||
		p.Peek(TokenType_PERIOD) ||
		p.Peek(TokenType_NO_SUBSTITUTION_TEMPLATE) ||
		p.Peek(TokenType_TEMPLATE_HEAD)
}

// Java: Parser.parseNewExpression()
func (p *Parser) ParseNewExpression() ParseTree {
	if !p.Peek(TokenType_NEW) {
		return p.ParseMemberExpressionNoNew()
	} else if p.PeekAt(1, TokenType_PERIOD) {
		return p.ParseNewDotSomething()
	} else {
		start := p.GetTreeStartLocation()
		p.Eat(TokenType_NEW)
		if p.Peek(TokenType_QUESTION_DOT) { // new?.target not allowed
			p.ReportError("Optional chaining is forbidden in `new?.target` contexts.")
		}
		operand := p.ParseNewExpression()
		if p.Peek(TokenType_QUESTION_DOT) { // new a?.() not allowed
			p.ReportError("Optional chaining is forbidden in construction contexts.")
		}
		var arguments *ArgumentListTree
		if p.Peek(TokenType_OPEN_PAREN) {
			arguments = p.ParseArguments()
		}
		return NewNewExpressionTree(
			p.GetTreeLocation(start),
			operand,
			arguments,
			arguments != nil && arguments.HasTrailingComma)
	}
}

// Java: Parser.parseNewDotSomething()
func (p *Parser) ParseNewDotSomething() ParseTree {
	// currently only "target" is valid after "new."
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_NEW)
	p.Eat(TokenType_PERIOD)
	p.EatPredefinedString("target")
	return NewNewTargetExpressionTree(p.GetTreeLocation(start))
}

// Java: Parser.parseImportDotMeta()
func (p *Parser) ParseImportDotMeta() ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_IMPORT)
	p.Eat(TokenType_PERIOD)
	p.EatPredefinedString("meta")
	return NewImportMetaExpressionTree(p.GetTreeLocation(start))
}

// Java: Parser.parseArguments()
func (p *Parser) ParseArguments() *ArgumentListTree {
	// ArgumentList :
	//   AssignmentOrSpreadExpression
	//   ArgumentList , AssignmentOrSpreadExpression
	//
	// AssignmentOrSpreadExpression :
	//   ... AssignmentExpression
	//   AssignmentExpression

	start := p.GetTreeStartLocation()
	treeBase := len(p.treeStack)
	trailingComma := false
	posBase := len(p.posStack)

	p.Eat(TokenType_OPEN_PAREN)
	for p.PeekAssignmentOrSpread() {
		arg := p.ParseAssignmentOrSpread()
		p.treeStack = append(p.treeStack, arg)

		if !p.Peek(TokenType_CLOSE_PAREN) {
			comma := p.Eat(TokenType_COMMA)
			if comma != nil {
				p.posStack = append(p.posStack, comma.GetStart())
			}
			if p.Peek(TokenType_CLOSE_PAREN) {
				if !p.config.atLeast8 {
					p.ReportErrorAtToken(comma, "Invalid trailing comma in arguments list")
				}
				trailingComma = true
			}
		}
	}
	p.Eat(TokenType_CLOSE_PAREN)
	arguments := popStack(&p.treeStack, treeBase)
	commaPositions := popStack(&p.posStack, posBase)
	return NewArgumentListTree(
		p.GetTreeLocation(start), arguments, trailingComma, commaPositions)
}

// PeekAssignmentOrSpread reports whether we have a spread expression or an
// assignment next.
//
// This does not peek the operand for the spread expression. This means that
// ParseAssignmentOrSpread might still fail when this returns true.
//
// Java: Parser.peekAssignmentOrSpread()
func (p *Parser) PeekAssignmentOrSpread() bool {
	return p.Peek(TokenType_ELLIPSIS) || p.PeekAssignmentExpression()
}

// Java: Parser.parseAssignmentOrSpread()
func (p *Parser) ParseAssignmentOrSpread() ParseTree {
	if p.Peek(TokenType_ELLIPSIS) {
		return p.ParseIterSpread()
	}
	return p.ParseAssignmentExpression()
}

// Destructuring (aka pattern matching); see
// http://wiki.ecmascript.org/doku.php?id=harmony:destructuring

// Parser_PatternKind is Parser.PatternKind: kinds of destructuring patterns.
type Parser_PatternKind int32

const (
	// A var, let, const; catch head; or formal parameter list--only
	// identifiers are allowed as lvalues
	Parser_PatternKind_INITIALIZER Parser_PatternKind = iota
	// An assignment or for-in initializer--any lvalue is allowed
	Parser_PatternKind_ANY
)

var parser_PatternKind_names = [...]string{"INITIALIZER", "ANY"}

// Parser_PatternKind_values is PatternKind.values().
//
// Java: Parser.PatternKind.values()
func Parser_PatternKind_values() []Parser_PatternKind {
	return []Parser_PatternKind{Parser_PatternKind_INITIALIZER, Parser_PatternKind_ANY}
}

// Name is Enum.name().
//
// Java: Parser.PatternKind.name()
func (k Parser_PatternKind) Name() string { return parser_PatternKind_names[k] }

// String is Enum.toString().
//
// Java: Parser.PatternKind.toString()
func (k Parser_PatternKind) String() string { return k.Name() }

// Java: Parser.peekPatternStart()
func (p *Parser) PeekPatternStart() bool {
	return p.Peek(TokenType_OPEN_SQUARE) || p.Peek(TokenType_OPEN_CURLY)
}

// Java: Parser.parsePattern(PatternKind)
func (p *Parser) ParsePattern(kind Parser_PatternKind) ParseTree {
	switch p.PeekType() {
	case TokenType_OPEN_SQUARE:
		return p.ParseArrayPattern(kind)
	case TokenType_OPEN_CURLY:
		fallthrough
	default:
		return p.ParseObjectPattern(kind)
	}
}

// Java: Parser.peekArrayPatternElement()
func (p *Parser) PeekArrayPatternElement() bool {
	return p.PeekExpression()
}

// Java: Parser.parseIterRest(PatternKind)
func (p *Parser) ParseIterRest(patternKind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_ELLIPSIS)
	patternAssignmentTarget := p.ParseRestAssignmentTarget(patternKind)
	return NewIterRestTree(p.GetTreeLocation(start), patternAssignmentTarget)
}

// Java: Parser.parseRestAssignmentTarget(PatternKind)
func (p *Parser) ParseRestAssignmentTarget(patternKind Parser_PatternKind) ParseTree {
	patternAssignmentTarget := p.ParsePatternAssignmentTargetNoDefault(patternKind)
	if p.Peek(TokenType_EQUAL) {
		p.ReportError("A default value cannot be specified after '...'")
	}
	return patternAssignmentTarget
}

// ParseArrayPattern parses Pattern ::= ... | "[" Element? ("," Element?)* "]"
//
// Java: Parser.parseArrayPattern(PatternKind)
func (p *Parser) ParseArrayPattern(kind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	elements := []ParseTree{}
	p.Eat(TokenType_OPEN_SQUARE)
	for p.Peek(TokenType_COMMA) || p.PeekArrayPatternElement() {
		if p.Peek(TokenType_COMMA) {
			nullStart := p.GetTreeStartLocation()
			p.Eat(TokenType_COMMA)
			elements = append(elements, NewNullTree(p.GetTreeLocation(nullStart)))
		} else {
			elements = append(elements, p.ParsePatternAssignmentTarget(kind))

			if p.Peek(TokenType_COMMA) {
				// Consume the comma separator
				p.Eat(TokenType_COMMA)
			} else {
				// Otherwise we must be done
				break
			}
		}
	}
	if p.Peek(TokenType_ELLIPSIS) {
		p.RecordFeatureUsed(FeatureSet_Feature_ARRAY_PATTERN_REST)
		elements = append(elements, p.ParseIterRest(kind))
	}
	if p.Eat(TokenType_CLOSE_SQUARE) == nil {
		// If we get no closing bracket then return invalid tree to avoid compiler tripping
		// downstream. It's needed only for IDE mode where compiler continues processing even if
		// source has syntactic errors.
		return NewMissingPrimaryExpressionTree(p.GetTreeLocation(p.GetTreeStartLocation()))
	}
	return NewArrayPatternTree(p.GetTreeLocation(start), elements)
}

// ParseObjectPattern parses
// Pattern ::= "{" (Field ("," Field)* ","?)? "}" | ...
//
// Java: Parser.parseObjectPattern(PatternKind)
func (p *Parser) ParseObjectPattern(kind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	fields := []ParseTree{}
	p.Eat(TokenType_OPEN_CURLY)
	for p.PeekObjectPatternField() {
		fields = append(fields, p.ParseObjectPatternField(kind))

		if p.Peek(TokenType_COMMA) {
			// Consume the comma separator
			p.Eat(TokenType_COMMA)
		} else {
			// Otherwise we must be done
			break
		}
	}
	if p.Peek(TokenType_ELLIPSIS) {
		p.RecordFeatureUsed(FeatureSet_Feature_OBJECT_PATTERN_REST)
		restStart := p.GetTreeStartLocation()
		p.Eat(TokenType_ELLIPSIS)
		patternAssignmentTarget := p.ParseRestAssignmentTarget(kind)
		fields = append(fields, NewObjectRestTree(p.GetTreeLocation(restStart), patternAssignmentTarget))
	}
	p.Eat(TokenType_CLOSE_CURLY)
	return NewObjectPatternTree(p.GetTreeLocation(start), fields)
}

// Java: Parser.peekObjectPatternField()
func (p *Parser) PeekObjectPatternField() bool {
	return p.PeekPropertyNameOrComputedProp(0)
}

// Java: Parser.parseObjectPatternField(PatternKind)
func (p *Parser) ParseObjectPatternField(kind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	if p.PeekType() == TokenType_OPEN_SQUARE {
		key := p.ParseComputedPropertyName()
		p.Eat(TokenType_COLON)
		value := p.ParsePatternAssignmentTarget(kind)
		return NewComputedPropertyDefinitionTree(p.GetTreeLocation(start), key, value)
	}

	var name Token
	if p.PeekIdOrKeyword() {
		idToken := p.EatIdOrKeywordAsId()
		if !p.Peek(TokenType_COLON) {
			if idToken.IsKeyword() {
				p.ReportError("cannot use keyword '%s' here.", idToken)
			}
			if p.Peek(TokenType_EQUAL) {
				idTree := NewIdentifierExpressionTree(p.GetTreeLocation(start), idToken)
				p.Eat(TokenType_EQUAL)
				defaultValue := p.ParseAssignmentExpression()
				return NewDefaultParameterTree(p.GetTreeLocation(start), idTree, defaultValue)
			}
			return NewPropertyNameAssignmentTree(p.GetTreeLocation(start), identifierTokenAsToken(idToken), nil)
		}
		name = identifierTokenAsToken(idToken)
	} else {
		name = p.ParseLiteralExpression().LiteralToken
	}

	p.Eat(TokenType_COLON)
	value := p.ParsePatternAssignmentTarget(kind)
	return NewPropertyNameAssignmentTree(p.GetTreeLocation(start), name, value)
}

// ParsePatternAssignmentTarget parses a PatternAssignmentTarget: the
// location where the assigned value gets stored, including an optional
// default value.
//
//	Spec AssignmentElement === PatternAssignmentTarget(PatternKind.ANY)
//	  Valid in an assignment that is not a formal parameter list or variable
//	  declaration. Sub-patterns and arbitrary left hand side expressions are
//	  allowed.
//	Spec BindingElement === PatternAssignmentElement(PatternKind.INITIALIZER)
//	  Valid in a formal parameter list or variable declaration statement. Only
//	  sub-patterns and identifiers are allowed.
//
// Examples:
//
//	[a, {foo: b = 'default'}] = someArray;          // valid
//	[x.a, {foo: x.b = 'default'}] = someArray;      // valid
//
//	let [a, {foo: b = 'default'}] = someArray;      // valid
//	let [x.a, {foo: x.b = 'default'}] = someArray;  // invalid
//
//	function f([a, {foo: b = 'default'}]) {...}     // valid
//	function f([x.a, {foo: x.b = 'default'}]) {...} // invalid
//
// Java: Parser.parsePatternAssignmentTarget(PatternKind)
func (p *Parser) ParsePatternAssignmentTarget(patternKind Parser_PatternKind) ParseTree {
	start := p.GetTreeStartLocation()
	var assignmentTarget ParseTree
	assignmentTarget = p.ParsePatternAssignmentTargetNoDefault(patternKind)

	if p.Peek(TokenType_EQUAL) {
		p.Eat(TokenType_EQUAL)
		defaultValue := p.ParseAssignmentExpression()
		assignmentTarget = NewDefaultParameterTree(p.GetTreeLocation(start), assignmentTarget, defaultValue)
	}
	return assignmentTarget
}

// Java: Parser.parsePatternAssignmentTargetNoDefault(PatternKind)
func (p *Parser) ParsePatternAssignmentTargetNoDefault(kind Parser_PatternKind) ParseTree {
	var assignmentTarget ParseTree
	if p.PeekPatternStart() {
		assignmentTarget = p.ParsePattern(kind)
	} else {
		assignmentTarget = p.ParseLeftHandSideExpression()
		if !assignmentTarget.IsValidAssignmentTarget() {
			p.ReportError("invalid assignment target")
		}
		if kind == Parser_PatternKind_INITIALIZER &&
			assignmentTarget.GetParseTreeBase().Type != ParseTreeType_IDENTIFIER_EXPRESSION {
			// We're in the context of a formal parameter list or a variable declaration statement
			p.ReportError("Only an identifier or destructuring pattern is allowed here.")
		}
	}
	return assignmentTarget
}

// EatPossiblyImplicitSemiColon consumes a (possibly implicit) semi-colon.
// Reports an error if a semi-colon is not present.
//
// Java: Parser.eatPossiblyImplicitSemiColon()
func (p *Parser) EatPossiblyImplicitSemiColon() {
	if p.Peek(TokenType_SEMI_COLON) {
		p.Eat(TokenType_SEMI_COLON)
		return
	}
	if p.PeekImplicitSemiColon() {
		return
	}

	p.ReportError("Semi-colon expected")
}

// PeekImplicitSemiColon returns true if an implicit or explicit semi colon
// is at the current location.
//
// Java: Parser.peekImplicitSemiColon()
func (p *Parser) PeekImplicitSemiColon() bool {
	return p.PeekImplicitSemiColonAt(0)
}

// Java: Parser.peekImplicitSemiColon(int)
func (p *Parser) PeekImplicitSemiColonAt(index int) bool {
	var lineAdvanced bool
	if index == 0 {
		lineAdvanced = p.GetNextLine() > p.GetLastLine()
	} else {
		lineAdvanced = p.PeekTokenAt(index).GetTokenBase().Location.Start.Line > p.PeekTokenAt(index-1).GetTokenBase().Location.End.Line
	}
	return lineAdvanced ||
		p.PeekAt(index, TokenType_SEMI_COLON) ||
		p.PeekAt(index, TokenType_CLOSE_CURLY) ||
		p.PeekAt(index, TokenType_END_OF_FILE)
}

// GetLastLine returns the line number of the most recently consumed token.
//
// Java: Parser.getLastLine()
func (p *Parser) GetLastLine() int {
	return int(p.lastSourcePosition.Line)
}

// GetNextLine returns the line number of the next token.
//
// Java: Parser.getNextLine()
func (p *Parser) GetNextLine() int {
	return int(p.PeekToken().GetTokenBase().Location.Start.Line)
}

// EatOpt consumes the next token if it is of the expected type. Otherwise
// returns nil. Never reports errors.
//
// Returns the consumed token, or nil if the next token is not of the
// expected type.
//
// Java: Parser.eatOpt(TokenType)
func (p *Parser) EatOpt(expectedTokenType TokenType) Token {
	if p.Peek(expectedTokenType) {
		return p.Eat(expectedTokenType)
	}
	return nil
}

// Java: Parser.inStrictContext()
func (p *Parser) InStrictContext() bool {
	// TODO(johnlenz): track entering strict scripts/modules/functions.
	return p.config.isStrictMode
}

// Java: Parser.peekId()
func (p *Parser) PeekId() bool {
	return p.PeekIdAt(0)
}

// PeekIdAt returns whether the next token is an identifier.
//
// Java: Parser.peekId(int)
func (p *Parser) PeekIdAt(index int) bool {
	token := p.PeekTokenAt(index)
	typeToCheck := token.GetTokenBase().Type

	if typeToCheck == TokenType_IDENTIFIER {
		value := token.(*IdentifierToken).GetMaybePrivateValue()
		k := Keywords_Get(value)
		if k == nil {
			return true
		}
		// Swap the type to the matched keyword's type
		typeToCheck = k.Type()
	}

	// There is one special case to handle here: outside of strict-mode code, strict-mode keywords
	// can be used as identifiers
	return !p.InStrictContext() && Keywords_IsStrictKeyword(typeToCheck)
}

// Java: Parser.peekIdOrKeyword()
func (p *Parser) PeekIdOrKeyword() bool {
	return p.PeekIdOrKeywordAt(0)
}

// Java: Parser.peekIdOrKeyword(int)
func (p *Parser) PeekIdOrKeywordAt(index int) bool {
	type_ := p.PeekTypeAt(index)
	return TokenType_IDENTIFIER == type_ || Keywords_IsKeywordTokenType(type_)
}

// EatIdOpt is shorthand for eatOpt(TokenType.IDENTIFIER).
//
// Java: Parser.eatIdOpt()
func (p *Parser) EatIdOpt() *IdentifierToken {
	if p.PeekId() {
		return p.EatIdOrKeywordAsId()
	}
	return nil
}

// EatId consumes an identifier token that is not a reserved word.
//
// See http://www.ecma-international.org/ecma-262/5.1/#sec-7.6
//
// Java: Parser.eatId()
func (p *Parser) EatId() *IdentifierToken {
	if p.PeekId() {
		return p.EatIdOrKeywordAsId()
	} else {
		p.ReportExpectedError(p.PeekToken(), TokenType_IDENTIFIER)
		if p.PeekIdOrKeyword() {
			return p.EatIdOrKeywordAsId()
		} else {
			return nil
		}
	}
}

// Java: Parser.eatObjectLiteralPropertyName()
func (p *Parser) EatObjectLiteralPropertyName() Token {
	token := p.PeekToken()
	switch token.GetTokenBase().Type {
	case TokenType_STRING, TokenType_NUMBER, TokenType_BIGINT:
		return p.NextToken()
	case TokenType_IDENTIFIER:
		fallthrough
	default:
		return identifierTokenAsToken(p.EatIdOrKeywordAsId())
	}
}

// EatIdOrKeywordAsId consumes an identifier token that may be a reserved
// word, i.e. an IdentifierName, not necessarily an Identifier.
//
// See http://www.ecma-international.org/ecma-262/5.1/#sec-7.6
//
// Java: Parser.eatIdOrKeywordAsId()
func (p *Parser) EatIdOrKeywordAsId() *IdentifierToken {
	token := p.NextToken()
	if token.GetTokenBase().Type == TokenType_IDENTIFIER {
		return token.(*IdentifierToken)
	} else if Keywords_IsKeywordTokenType(token.GetTokenBase().Type) {
		return NewIdentifierToken(token.GetTokenBase().Location, Keywords_GetByTokenType(token.GetTokenBase().Type).String())
	} else {
		p.ReportExpectedError(token, TokenType_IDENTIFIER)
	}
	return nil
}

// Eat consumes the next token. If the consumed token is not of the expected
// type then report an error and return nil. Otherwise return the consumed
// token.
//
// Returns the consumed token, or nil if the next token is not of the
// expected type.
//
// Java: Parser.eat(TokenType)
func (p *Parser) Eat(expectedTokenType TokenType) Token {
	token := p.NextToken()
	if token.GetTokenBase().Type != expectedTokenType {
		p.ReportExpectedError(token, expectedTokenType)
		return nil
	}
	return token
}

// ReportExpectedError reports a 'X' expected error message.
//
// token: the location to report the message at. expected: the thing that
// was expected.
//
// Java: Parser.reportExpectedError(Token, Object)
func (p *Parser) ReportExpectedError(token Token, expected any) {
	p.ReportErrorAtToken(token, "'%s' expected", expected)
}

// GetTreeStartLocation returns a SourcePosition for the start of a parse
// tree that starts at the current location.
//
// Java: Parser.getTreeStartLocation()
func (p *Parser) GetTreeStartLocation() SourcePosition {
	return p.PeekToken().GetTokenBase().Location.Start
}

// GetTreeEndLocation returns a SourcePosition for the end of a parse tree
// that ends at the current location.
//
// Java: Parser.getTreeEndLocation()
func (p *Parser) GetTreeEndLocation() SourcePosition {
	return p.lastSourcePosition
}

// GetTreeLocation returns a SourceRange for a parse tree that starts at
// {start} and ends at the current location.
//
// Java: Parser.getTreeLocation(SourcePosition)
func (p *Parser) GetTreeLocation(start SourcePosition) SourceRange {
	return NewSourceRange(start, p.GetTreeEndLocation())
}

// NextToken consumes the next token and returns it. Will return a never
// ending stream of TokenType.END_OF_FILE at the end of the file so callers
// don't have to check for EOF explicitly.
//
// Tokenizing is contextual. NextToken() will never return a regular
// expression literal.
//
// Java: Parser.nextToken()
func (p *Parser) NextToken() Token {
	token := p.scanner.NextToken()
	p.lastSourcePosition = token.GetTokenBase().Location.End
	return token
}

// NextRegularExpressionLiteralToken consumes a regular expression literal
// token and returns it.
//
// Java: Parser.nextRegularExpressionLiteralToken()
func (p *Parser) NextRegularExpressionLiteralToken() LiteralToken {
	token := p.scanner.NextRegularExpressionLiteralToken()
	p.lastSourcePosition = token.GetTokenBase().Location.End
	return token
}

// NextTemplateLiteralToken consumes a template literal token and returns
// it.
//
// Java: Parser.nextTemplateLiteralToken()
func (p *Parser) NextTemplateLiteralToken() *TemplateLiteralToken {
	token := p.scanner.NextTemplateLiteralToken()
	p.lastSourcePosition = token.Location.End
	return token
}

// Peek returns true if the next token is of the expected type. Does not
// consume the token.
//
// Java: Parser.peek(TokenType)
func (p *Parser) Peek(expectedType TokenType) bool {
	return p.PeekAt(0, expectedType)
}

// PeekAt returns true if the index-th next token is of the expected type.
// Does not consume any tokens.
//
// Java: Parser.peek(int, TokenType)
func (p *Parser) PeekAt(index int, expectedType TokenType) bool {
	return p.PeekTypeAt(index) == expectedType
}

// PeekType returns the TokenType of the next token. Does not consume any
// tokens.
//
// Java: Parser.peekType()
func (p *Parser) PeekType() TokenType {
	return p.PeekTypeAt(0)
}

// PeekTypeAt returns the TokenType of the index-th next token. Does not
// consume any tokens.
//
// Java: Parser.peekType(int)
func (p *Parser) PeekTypeAt(index int) TokenType {
	return p.PeekTokenAt(index).GetTokenBase().Type
}

// PeekToken returns the next token. Does not consume any tokens.
//
// Java: Parser.peekToken()
func (p *Parser) PeekToken() Token {
	return p.PeekTokenAt(0)
}

// PeekTokenAt returns the index-th next token. Does not consume any tokens.
//
// Java: Parser.peekToken(int)
func (p *Parser) PeekTokenAt(index int) Token {
	return p.scanner.PeekTokenAt(index)
}

// ReportErrorAtToken reports an error message at a given token.
//
// token: the location to report the message at. message: the message to
// report in String.format style. arguments: the arguments to fill in the
// message format.
//
// Java: Parser.reportError(Token, String, Object...)
func (p *Parser) ReportErrorAtToken(token Token, message string, arguments ...any) {
	if token == nil {
		p.ReportError(message, arguments...)
	} else {
		p.errorReporter.ReportErrorFormat(token.GetStart(), message, arguments...)
	}
}

// ReportErrorAtTree reports an error message at a given parse tree's
// location.
//
// parseTree: the location to report the message at. message: the message to
// report in String.format style. arguments: the arguments to fill in the
// message format.
//
// Java: Parser.reportError(ParseTree, String, Object...)
func (p *Parser) ReportErrorAtTree(parseTree ParseTree, message string, arguments ...any) {
	if parseTree == nil {
		p.ReportError(message, arguments...)
	} else {
		p.errorReporter.ReportErrorFormat(parseTree.GetParseTreeBase().Location.Start, message, arguments...)
	}
}

// ReportError reports an error at the current location.
//
// message: the message to report in String.format style. arguments: the
// arguments to fill in the message format.
//
// Java: Parser.reportError(String, Object...)
func (p *Parser) ReportError(message string, arguments ...any) {
	p.errorReporter.ReportErrorFormat(p.scanner.GetPosition(), message, arguments...)
}

// ReportErrorAtPosition reports an error at the specified location.
//
// position: the position of the error. message: the message to report in
// String.format style. arguments: the arguments to fill in the message
// format.
//
// Java: Parser.reportError(SourcePosition, String, Object...)
func (p *Parser) ReportErrorAtPosition(position SourcePosition, message string, arguments ...any) {
	p.errorReporter.ReportErrorFormat(position, message, arguments...)
}

// Java: Parser.reportTemplateErrorIfPresent(TemplateLiteralToken)
func (p *Parser) ReportTemplateErrorIfPresent(templateToken *TemplateLiteralToken) {
	if templateToken.ErrorMessage == nil {
		return
	}
	// Java: switch (templateToken.errorLevel) — a null errorLevel would throw
	// a NullPointerException here.
	switch *templateToken.ErrorLevel {
	case TemplateLiteralToken_ErrorLevel_WARNING:
		p.errorReporter.ReportWarningFormat(templateToken.ErrorPosition, "%s", *templateToken.ErrorMessage)
		return
	case TemplateLiteralToken_ErrorLevel_ERROR:
		p.ReportErrorAtPosition(templateToken.ErrorPosition, "%s", *templateToken.ErrorMessage)
		return
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: Parser.recordFeatureUsed(Feature)
func (p *Parser) RecordFeatureUsed(feature FeatureSet_Feature) *Parser {
	p.features = p.features.With(feature)
	return p
}

// Go-only: language-adapter — conversions from possibly-nil concrete
// pointers to interfaces, so that Java null stays a nil interface (a typed
// nil pointer in an interface would compare non-nil).

// identifierTokenAsToken widens a possibly-nil *IdentifierToken to Token.
func identifierTokenAsToken(t *IdentifierToken) Token {
	if t == nil {
		return nil
	}
	return t
}

// variableDeclarationListAsParseTree widens a possibly-nil
// *VariableDeclarationListTree to ParseTree.
func variableDeclarationListAsParseTree(t *VariableDeclarationListTree) ParseTree {
	if t == nil {
		return nil
	}
	return t
}
