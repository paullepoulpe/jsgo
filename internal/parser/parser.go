// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20261021

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
	config               Parser_Config
	commentRecorder      Parser_CommentRecorder
	functionContextStack []Parser_FunctionFlavor
	functionContextBuf   [8]Parser_FunctionFlavor
	featureBits          uint64
	features             *FeatureSet
	lastSourcePosition   SourcePosition
	// Go-only: language-adapter — nil is Java null.
	sourceMapURL *string

	// Go-only: performance — scratch stacks for the lists being parsed (Java
	// uses an ImmutableList.Builder per list); popStackChunked copies each
	// finished list out at its exact length.
	treeStack    []ParseTree
	posStack     []SourcePosition
	varDeclStack []*VariableDeclarationTree

	// Go-only: performance — the commonest trees, and the child lists
	// popped from treeStack, are allocated in slabs (see jutil.Slab): a
	// statement's parse trees are created together and die together once
	// it is converted.
	identExprs   jutil.Slab[IdentifierExpressionTree]
	memberExprs  jutil.Slab[MemberExpressionTree]
	binaryOps    jutil.Slab[BinaryOperatorTree]
	literals     jutil.Slab[LiteralExpressionTree]
	propAssigns  jutil.Slab[PropertyNameAssignmentTree]
	exprStmts    jutil.Slab[ExpressionStatementTree]
	calls        jutil.Slab[CallExpressionTree]
	argLists     jutil.Slab[ArgumentListTree]
	blocks       jutil.Slab[BlockTree]
	varDecls     jutil.Slab[VariableDeclarationTree]
	varDeclLsts  jutil.Slab[VariableDeclarationListTree]
	varStmts     jutil.Slab[VariableStatementTree]
	memberLkps   jutil.Slab[MemberLookupExpressionTree]
	thisExprs    jutil.Slab[ThisExpressionTree]
	funcDecls    jutil.Slab[FunctionDeclarationTree]
	paramLists   jutil.Slab[FormalParameterListTree]
	returnStmts  jutil.Slab[ReturnStatementTree]
	ifStmts      jutil.Slab[IfStatementTree]
	parenExprs   jutil.Slab[ParenExpressionTree]
	unaryExprs   jutil.Slab[UnaryExpressionTree]
	treeLists    []ParseTree
	posLists     []SourcePosition
	varDeclLists []*VariableDeclarationTree
}

func (p *Parser) newBlockTree(location SourceRange, statements []ParseTree) *BlockTree {
	t := p.blocks.New()
	t.InitParseTreeBase(ParseTreeType_BLOCK, location)
	t.Statements = statements
	return t
}

func (p *Parser) newVariableDeclarationTree(location SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	t := p.varDecls.New()
	t.InitParseTreeBase(ParseTreeType_VARIABLE_DECLARATION, location)
	t.Lvalue = lvalue
	t.Initializer = initializer
	return t
}

func (p *Parser) newVariableDeclarationListTree(location SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	t := p.varDeclLsts.New()
	t.InitParseTreeBase(ParseTreeType_VARIABLE_DECLARATION_LIST, location)
	t.DeclarationType = declarationType
	t.Declarations = declarations
	return t
}

func (p *Parser) newVariableStatementTree(location SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	t := p.varStmts.New()
	t.InitParseTreeBase(ParseTreeType_VARIABLE_STATEMENT, location)
	t.Declarations = declarations
	return t
}

func (p *Parser) newMemberLookupExpressionTree(location SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	t := p.memberLkps.New()
	t.InitParseTreeBase(ParseTreeType_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	return t
}

func (p *Parser) newThisExpressionTree(location SourceRange) *ThisExpressionTree {
	t := p.thisExprs.New()
	t.InitParseTreeBase(ParseTreeType_THIS_EXPRESSION, location)
	return t
}

func (p *Parser) newExpressionStatementTree(location SourceRange, expression ParseTree) *ExpressionStatementTree {
	t := p.exprStmts.New()
	t.InitParseTreeBase(ParseTreeType_EXPRESSION_STATEMENT, location)
	t.Expression = expression
	return t
}

func (p *Parser) newCallExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	t := p.calls.New()
	t.InitParseTreeBase(ParseTreeType_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	return t
}

func (p *Parser) newArgumentListTree(location SourceRange, arguments []ParseTree, hasTrailingComma bool, commaPositions []SourcePosition) *ArgumentListTree {
	t := p.argLists.New()
	t.InitParseTreeBase(ParseTreeType_ARGUMENT_LIST, location)
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}

// parser_treeListChunk and parser_posListChunk are the sizes of the chunks
// child lists are carved from (4 KB, counting the 8-byte header of a
// pointerful allocation); a list longer than an eighth of a chunk is
// allocated on its own.
const (
	parser_treeListChunk = 254
	parser_posListChunk  = 340
)

// popStackChunked returns the elements pushed on *stack since base, carved
// from *chunk (a shared chunk of chunkSize elements), and pops them.
func popStackChunked[T any](stack *[]T, base int, chunk *[]T, chunkSize int) []T {
	s := (*stack)[base:]
	n := len(s)
	var res []T
	switch {
	case n == 0:
		res = []T{}
	case n > chunkSize/8:
		res = make([]T, n)
	default:
		if len(*chunk) < n {
			*chunk = jutil.SlabChunk[T](chunkSize)
		}
		res = (*chunk)[:n:n]
		*chunk = (*chunk)[n:]
	}
	copy(res, s)
	clear(s)
	*stack = (*stack)[:base]
	return res
}

// popTrees pops the elements pushed on p.treeStack since base into a slice
// carved from a shared chunk.
func (p *Parser) popTrees(base int) []ParseTree {
	return popStackChunked(&p.treeStack, base, &p.treeLists, parser_treeListChunk)
}

func (p *Parser) singleTreeList(t ParseTree) []ParseTree {
	if len(p.treeLists) < 1 {
		p.treeLists = jutil.SlabChunk[ParseTree](parser_treeListChunk)
	}
	res := p.treeLists[:1:1]
	p.treeLists = p.treeLists[1:]
	res[0] = t
	return res
}

func (p *Parser) newIdentifierExpressionTree(location SourceRange, identifierToken Token) *IdentifierExpressionTree {
	t := p.identExprs.New()
	t.InitParseTreeBase(ParseTreeType_IDENTIFIER_EXPRESSION, location)
	t.IdentifierToken = identifierToken
	return t
}

func (p *Parser) newMemberExpressionTree(location SourceRange, operand ParseTree, memberName Token) *MemberExpressionTree {
	t := p.memberExprs.New()
	t.InitParseTreeBase(ParseTreeType_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	return t
}

func (p *Parser) newBinaryOperatorTree(location SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	t := p.binaryOps.New()
	t.InitParseTreeBase(ParseTreeType_BINARY_OPERATOR, location)
	t.Left = left
	t.Operator = operator
	t.Right = right
	return t
}

func (p *Parser) newLiteralExpressionTree(location SourceRange, literalToken Token) *LiteralExpressionTree {
	t := p.literals.New()
	t.InitParseTreeBase(ParseTreeType_LITERAL_EXPRESSION, location)
	t.LiteralToken = literalToken
	return t
}

func (p *Parser) newPropertyNameAssignmentTree(location SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	t := p.propAssigns.New()
	t.InitParseTreeBase(ParseTreeType_PROPERTY_NAME_ASSIGNMENT, location)
	t.Name = name
	t.Value = value
	return t
}

func (p *Parser) newFunctionDeclarationTree(
	location SourceRange,
	name Token,
	formalParameterList *FormalParameterListTree,
	functionBody ParseTree,
	kind FunctionDeclarationTree_Kind,
	isClassMember, isStatic, isGenerator, isAsync bool,
) *FunctionDeclarationTree {
	if formalParameterList == nil || functionBody == nil {
		panic(jutil.NewNullPointerException())
	}
	t := p.funcDecls.New()
	t.InitParseTreeBase(ParseTreeType_FUNCTION_DECLARATION, location)
	t.Name = name
	t.FormalParameterList = formalParameterList
	t.FunctionBody = functionBody
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.IsGenerator = isGenerator
	t.IsAsync = isAsync
	t.Kind = kind
	return t
}

func (p *Parser) newFormalParameterListTree(location SourceRange, parameters []ParseTree, hasTrailingComma bool, commaPositions []SourcePosition) *FormalParameterListTree {
	numParams := len(parameters)
	numCommas := len(commaPositions)
	if !(numCommas <= numParams && numCommas >= numParams-1) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat(
			"Unexpected # of comma and formal params.\nparams: %s\ncomma positions: %s",
			jutil.CollectionToString(parameters),
			jutil.CollectionToString(commaPositions))))
	}
	t := p.paramLists.New()
	t.InitParseTreeBase(ParseTreeType_FORMAL_PARAMETER_LIST, location)
	t.Parameters = parameters
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}

func (p *Parser) newReturnStatementTree(location SourceRange, expression ParseTree) *ReturnStatementTree {
	t := p.returnStmts.New()
	t.InitParseTreeBase(ParseTreeType_RETURN_STATEMENT, location)
	t.Expression = expression
	return t
}

func (p *Parser) newIfStatementTree(location SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	t := p.ifStmts.New()
	t.InitParseTreeBase(ParseTreeType_IF_STATEMENT, location)
	t.Condition = condition
	t.IfClause = ifClause
	t.ElseClause = elseClause
	return t
}

func (p *Parser) newParenExpressionTree(location SourceRange, expression ParseTree) *ParenExpressionTree {
	t := p.parenExprs.New()
	t.InitParseTreeBase(ParseTreeType_PAREN_EXPRESSION, location)
	t.Expression = expression
	return t
}

func (p *Parser) newUnaryExpressionTree(location SourceRange, operator Token, operand ParseTree) *UnaryExpressionTree {
	t := p.unaryExprs.New()
	t.InitParseTreeBase(ParseTreeType_UNARY_EXPRESSION, location)
	t.Operator = operator
	t.Operand = operand
	return t
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
func NewParser(config Parser_Config, errorReporter ErrorReporter, source *SourceFile) *Parser {
	p := &Parser{}
	p.commentRecorder = Parser_CommentRecorder{this: p, lastCommentEndOffset: -1}
	p.functionContextStack = p.functionContextBuf[:0]
	p.config = config
	p.errorReporter = errorReporter
	p.scanner = NewScanner(errorReporter, &p.commentRecorder, source, 0)
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
func NewParser_Config() Parser_Config {
	return NewParser_ConfigWithMode(Parser_Config_Mode_ES8_OR_GREATER /* isStrictMode= */, false)
}

// NewParser_ConfigWithMode constructs a Config.
//
// Java: Parser.Config(Mode, boolean)
func NewParser_ConfigWithMode(mode Parser_Config_Mode, isStrictMode bool) Parser_Config {
	return Parser_Config{
		atLeast6:     !(mode == Parser_Config_Mode_ES3 || mode == Parser_Config_Mode_ES5),
		atLeast8:     mode == Parser_Config_Mode_ES8_OR_GREATER,
		isStrictMode: isStrictMode,
		// Generally, we allow everything that is valid in any mode
		// we only warn about things that are not represented in the AST.
		warnTrailingCommas: mode == Parser_Config_Mode_ES3,
	}
}

// Parser_SOURCE_MAPPING_URL_PREFIX is Parser.SOURCE_MAPPING_URL_PREFIX.
const Parser_SOURCE_MAPPING_URL_PREFIX = "//# sourceMappingURL="

// Parser_CommentRecorder is the inner class Parser.CommentRecorder.
type Parser_CommentRecorder struct {
	// Go-only: language-adapter — the enclosing Parser instance of the Java
	// inner class.
	this *Parser

	comments []Comment
	// Go-only: language-adapter — the offset of Java's
	// lastCommentEndPosition, -1 for null.
	lastCommentEndOffset int32
}

// NewParser_CommentRecorder constructs the comment recorder of parser.
//
// Java: Parser.CommentRecorder()
func NewParser_CommentRecorder(this *Parser) *Parser_CommentRecorder {
	return &Parser_CommentRecorder{this: this, comments: []Comment{}, lastCommentEndOffset: -1}
}

// RecordComment records a comment the scanner passed.
//
// Java: Parser.CommentRecorder.recordComment(Comment.Type, SourceRange, String)
func (c *Parser_CommentRecorder) RecordComment(type_ Comment_Type, range_ LineRange, value string) {
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
func (c *Parser_CommentRecorder) GetComments() []Comment {
	if len(c.comments) == 0 {
		return []Comment{}
	}
	return c.comments
}

// GetComments returns the comments of the source file.
//
// Java: Parser.getComments()
func (p *Parser) GetComments() []Comment {
	return p.commentRecorder.GetComments()
}

// GetFeatures returns the language features used by the source file.
//
// Java: Parser.getFeatures()
func (p *Parser) GetFeatures() *FeatureSet {
	if p.featureBits == 0 {
		return FeatureSet_BARE_MINIMUM
	}
	if p.features == nil || p.features.bits != p.featureBits {
		p.features = &FeatureSet{bits: p.featureBits}
	}
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

	return p.popTrees(base)
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
		moduleSpecifier := p.Eat(TokenType_STRING)
		p.EatPossiblyImplicitSemiColon()

		return NewImportDeclarationTree(p.GetTreeLocation(start), Token{}, nil, Token{}, moduleSpecifier)
	}

	// import ImportedDefaultBinding from ModuleSpecifier
	// import NameSpaceImport from ModuleSpecifier
	// import NamedImports from ModuleSpecifier ;
	// import ImportedDefaultBinding , NameSpaceImport from ModuleSpecifier ;
	// import ImportedDefaultBinding , NamedImports from ModuleSpecifier ;
	var defaultBindingIdentifier Token
	var nameSpaceImportIdentifier Token
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
	moduleSpecifier := p.Eat(TokenType_STRING)
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
	base := len(p.treeStack)
	p.Eat(TokenType_OPEN_CURLY)
	for p.PeekIdOrKeyword() {
		spec := p.ParseImportSpecifier()
		p.treeStack = append(p.treeStack, spec)
		if !p.Peek(TokenType_CLOSE_CURLY) {
			p.Eat(TokenType_COMMA)
		}
	}
	p.Eat(TokenType_CLOSE_CURLY)
	return p.popTrees(base)
}

// ParseImportSpecifier parses ImportSpecifier ::= Identifier ('as' Identifier)?
//
// Java: Parser.parseImportSpecifier()
func (p *Parser) ParseImportSpecifier() ParseTree {
	start := p.GetTreeStartLocation()
	importedName := p.EatIdOrKeywordAsId()
	var destinationName Token
	if p.PeekPredefinedString(PredefinedName_AS) {
		p.EatPredefinedString(PredefinedName_AS)
		destinationName = p.EatId()
	} else if importedName.IsKeyword() {
		p.ReportExpectedError(Token{}, PredefinedName_AS)
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

	var moduleSpecifier Token
	if isExportAll || (isExportSpecifier && p.PeekPredefinedString(PredefinedName_FROM)) {
		p.EatPredefinedString(PredefinedName_FROM)
		moduleSpecifier = p.Eat(TokenType_STRING)
	} else if isExportSpecifier {
		for _, tree := range exportSpecifierList {
			importedName := tree.(*ExportSpecifierTree).ImportedName
			if importedName.IsKeyword() {
				p.ReportErrorAtToken(importedName, "cannot use keyword '%s' here.", importedName)
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
	base := len(p.treeStack)
	p.Eat(TokenType_OPEN_CURLY)
	for p.PeekIdOrKeyword() {
		spec := p.ParseExportSpecifier()
		p.treeStack = append(p.treeStack, spec)
		if !p.Peek(TokenType_CLOSE_CURLY) {
			p.Eat(TokenType_COMMA)
		}
	}
	p.Eat(TokenType_CLOSE_CURLY)
	return p.popTrees(base)
}

// ParseExportSpecifier parses ExportSpecifier ::= Identifier ('as' Identifier)?
//
// Java: Parser.parseExportSpecifier()
func (p *Parser) ParseExportSpecifier() ParseTree {
	start := p.GetTreeStartLocation()
	importedName := p.EatIdOrKeywordAsId()
	var destinationName Token
	if p.PeekPredefinedString(PredefinedName_AS) {
		p.EatPredefinedString(PredefinedName_AS)
		destinationName = p.EatIdOrKeywordAsId()
	}
	return NewExportSpecifierTree(p.GetTreeLocation(start), importedName, destinationName)
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

// Java: Parser.parseFunctionTail(FunctionDeclarationTree.Builder, FunctionFlavor)
func (p *Parser) parseFunctionTail(functionFlavor Parser_FunctionFlavor) (*FormalParameterListTree, *BlockTree) {
	p.functionContextStack = append(p.functionContextStack, functionFlavor)
	params := p.ParseFormalParameterList()
	body := p.ParseFunctionBody()
	p.functionContextStack = p.functionContextStack[:len(p.functionContextStack)-1]
	return params, body
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
	isGenerator := p.EatOpt(TokenType_STAR).IsPresent()

	name := p.EatId()
	params, body := p.parseFunctionTail(functionFlavorOf(isGenerator, Parser_FunctionFlavor_GENERATOR, Parser_FunctionFlavor_NORMAL))
	return p.newFunctionDeclarationTree(
		p.GetTreeLocation(start), name, params, body,
		FunctionDeclarationTree_Kind_DECLARATION, false, false, isGenerator, false)
}

// Java: Parser.parseAsyncFunctionDeclaration()
func (p *Parser) ParseAsyncFunctionDeclaration() ParseTree {
	start := p.GetTreeStartLocation()
	p.EatAsyncFunctionStart()

	generator := p.Peek(TokenType_STAR)
	if generator {
		p.Eat(TokenType_STAR)
	}

	name := p.EatId()
	params, body := p.parseFunctionTail(
		functionFlavorOf(generator, Parser_FunctionFlavor_ASYNCHRONOUS_GENERATOR, Parser_FunctionFlavor_ASYNCHRONOUS))
	return p.newFunctionDeclarationTree(
		p.GetTreeLocation(start), name, params, body,
		FunctionDeclarationTree_Kind_DECLARATION, false, false, generator, true)
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
			if comma.IsPresent() {
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
	result := p.popTrees(treeBase)
	commaPositions := popStackChunked(&p.posStack, posBase, &p.posLists, parser_posListChunk)
	return p.newFormalParameterListTree(
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
	return p.newBlockTree(p.GetTreeLocation(start), result)
}

// Java: Parser.parseSourceElementList()
func (p *Parser) ParseSourceElementList() []ParseTree {
	base := len(p.treeStack)

	for p.PeekSourceElement() {
		elem := p.ParseSourceElement()
		p.treeStack = append(p.treeStack, elem)
	}

	return p.popTrees(base)
}

// Java: Parser.parseIterSpread()
func (p *Parser) ParseIterSpread() *IterSpreadTree {
	start := p.GetTreeStartLocation()
	p.Eat(TokenType_ELLIPSIS)
	operand := p.ParseAssignmentExpression()
	return NewIterSpreadTree(p.GetTreeLocation(start), operand)
}
