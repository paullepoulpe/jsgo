package parser

import (
	"sync"
	"sync/atomic"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	arenaSlabMinChunk  = 128
	arenaSlabMaxChunk  = 8192
	arenaSlabMaxChunks = 256
	arenaMaxUTF16Cap   = 8 * 1024 * 1024
)

type arenaSlab[T any] struct {
	chunks   [][]T
	chunk    []T
	chunkIdx int
	used     int
}

func (s *arenaSlab[T]) New() *T {
	if s.used == len(s.chunk) {
		s.grow()
	}
	p := &s.chunk[s.used]
	s.used++
	return p
}

func (s *arenaSlab[T]) grow() {
	if s.chunk != nil {
		s.chunkIdx++
	}
	if s.chunkIdx < len(s.chunks) {
		s.chunk = s.chunks[s.chunkIdx]
		s.used = 0
		return
	}
	n := arenaSlabMinChunk
	if s.chunkIdx > 0 {
		n = 2 * len(s.chunks[s.chunkIdx-1])
		if n > arenaSlabMaxChunk {
			n = arenaSlabMaxChunk
		}
	}
	c := make([]T, n)
	s.chunks = append(s.chunks, c)
	s.chunk = c
	s.used = 0
}

func (s *arenaSlab[T]) reset() {
	if s.chunk == nil {
		return
	}
	for i := 0; i < s.chunkIdx; i++ {
		clear(s.chunks[i])
	}
	clear(s.chunk[:s.used])
	if len(s.chunks) > arenaSlabMaxChunks {
		clear(s.chunks[arenaSlabMaxChunks:])
		s.chunks = s.chunks[:arenaSlabMaxChunks]
	}
	s.chunk = nil
	s.chunkIdx = 0
	s.used = 0
}

// ParseArena holds reusable buffers and slabs for a single file parse in
// ParserRunner_Parse. All objects allocated from ParseArena are transient to
// the parse and must not be retained after ParseArena.Release().
type ParseArena struct {
	utf16Buf []uint16

	// Scanner / LineNumberScanner / Parser position, range, and token slabs.
	positions      arenaSlab[SourcePosition]
	ranges         arenaSlab[SourceRange]
	tokens         arenaSlab[TokenBase]
	identifiers    arenaSlab[IdentifierToken]
	literals       arenaSlab[LiteralTokenBase]
	stringLiterals arenaSlab[StringLiteralToken]

	// High-frequency ParseTree slabs.
	programs                    arenaSlab[ProgramTree]
	functionBuilders            arenaSlab[FunctionDeclarationTree_Builder]
	functionDeclarations        arenaSlab[FunctionDeclarationTree]
	formalParameterLists        arenaSlab[FormalParameterListTree]
	blocks                      arenaSlab[BlockTree]
	variableStatements          arenaSlab[VariableStatementTree]
	variableDeclarationLists    arenaSlab[VariableDeclarationListTree]
	variableDeclarations        arenaSlab[VariableDeclarationTree]
	emptyStatements             arenaSlab[EmptyStatementTree]
	expressionStatements        arenaSlab[ExpressionStatementTree]
	ifStatements                arenaSlab[IfStatementTree]
	nullTrees                   arenaSlab[NullTree]
	returnStatements            arenaSlab[ReturnStatementTree]
	thisExpressions             arenaSlab[ThisExpressionTree]
	identifierExpressions       arenaSlab[IdentifierExpressionTree]
	literalExpressions          arenaSlab[LiteralExpressionTree]
	arrayLiteralExpressions     arenaSlab[ArrayLiteralExpressionTree]
	objectLiteralExpressions    arenaSlab[ObjectLiteralExpressionTree]
	propertyNameAssignments     arenaSlab[PropertyNameAssignmentTree]
	commaExpressions            arenaSlab[CommaExpressionTree]
	parenExpressions            arenaSlab[ParenExpressionTree]
	binaryOperators             arenaSlab[BinaryOperatorTree]
	conditionalExpressions      arenaSlab[ConditionalExpressionTree]
	unaryExpressions            arenaSlab[UnaryExpressionTree]
	updateExpressions           arenaSlab[UpdateExpressionTree]
	callExpressions             arenaSlab[CallExpressionTree]
	memberLookupExpressions     arenaSlab[MemberLookupExpressionTree]
	memberExpressions           arenaSlab[MemberExpressionTree]
	newExpressions              arenaSlab[NewExpressionTree]
	argumentLists               arenaSlab[ArgumentListTree]
	classDeclarations           arenaSlab[ClassDeclarationTree]
	computedPropertyMethods     arenaSlab[ComputedPropertyMethodTree]
	whileStatements             arenaSlab[WhileStatementTree]
	forOfStatements             arenaSlab[ForOfStatementTree]
	forStatements               arenaSlab[ForStatementTree]
	forInStatements             arenaSlab[ForInStatementTree]
	continueStatements          arenaSlab[ContinueStatementTree]
	breakStatements             arenaSlab[BreakStatementTree]
	switchStatements            arenaSlab[SwitchStatementTree]
	caseClauses                 arenaSlab[CaseClauseTree]
	throwStatements             arenaSlab[ThrowStatementTree]
	tryStatements               arenaSlab[TryStatementTree]
	catchTrees                  arenaSlab[CatchTree]
}

var (
	parseArenaSlots [4]atomic.Pointer[ParseArena]
	parseArenaPool  = sync.Pool{
		New: func() any {
			return &ParseArena{}
		},
	}
)

// AcquireParseArena acquires a reusable ParseArena from the pool.
func AcquireParseArena() *ParseArena {
	for i := range parseArenaSlots {
		if a := parseArenaSlots[i].Swap(nil); a != nil {
			return a
		}
	}
	return parseArenaPool.Get().(*ParseArena)
}

// Release resets the arena and returns it to the pool.
func (a *ParseArena) Release() {
	if a == nil {
		return
	}
	a.reset()
	for i := range parseArenaSlots {
		if parseArenaSlots[i].CompareAndSwap(nil, a) {
			return
		}
	}
	parseArenaPool.Put(a)
}

// DrainParseArenaPool releases any cached ParseArenas from the global slot cache
// so the garbage collector can reclaim their backing slabs after batch parsing completes.
func DrainParseArenaPool() {
	for i := range parseArenaSlots {
		parseArenaSlots[i].Store(nil)
	}
}

func (a *ParseArena) reset() {
	if cap(a.utf16Buf) <= arenaMaxUTF16Cap {
		a.utf16Buf = a.utf16Buf[:0]
	} else {
		a.utf16Buf = nil
	}
	a.positions.reset()
	a.ranges.reset()
	a.tokens.reset()
	a.identifiers.reset()
	a.literals.reset()
	a.stringLiterals.reset()

	a.programs.reset()
	a.functionBuilders.reset()
	a.functionDeclarations.reset()
	a.formalParameterLists.reset()
	a.blocks.reset()
	a.variableStatements.reset()
	a.variableDeclarationLists.reset()
	a.variableDeclarations.reset()
	a.emptyStatements.reset()
	a.expressionStatements.reset()
	a.ifStatements.reset()
	a.nullTrees.reset()
	a.returnStatements.reset()
	a.thisExpressions.reset()
	a.identifierExpressions.reset()
	a.literalExpressions.reset()
	a.arrayLiteralExpressions.reset()
	a.objectLiteralExpressions.reset()
	a.propertyNameAssignments.reset()
	a.commaExpressions.reset()
	a.parenExpressions.reset()
	a.binaryOperators.reset()
	a.conditionalExpressions.reset()
	a.unaryExpressions.reset()
	a.updateExpressions.reset()
	a.callExpressions.reset()
	a.memberLookupExpressions.reset()
	a.memberExpressions.reset()
	a.newExpressions.reset()
	a.argumentLists.reset()
	a.classDeclarations.reset()
	a.computedPropertyMethods.reset()
	a.whileStatements.reset()
	a.forOfStatements.reset()
	a.forStatements.reset()
	a.forInStatements.reset()
	a.continueStatements.reset()
	a.breakStatements.reset()
	a.switchStatements.reset()
	a.caseClauses.reset()
	a.throwStatements.reset()
	a.tryStatements.reset()
	a.catchTrees.reset()
}

func (a *ParseArena) decodeUTF16(s string) ([]uint16, int) {
	prefix := jutil.ASCIIPrefixLen(s)
	n := len(s)
	buf := a.utf16Buf
	if cap(buf) < n {
		if n <= arenaMaxUTF16Cap {
			newCap := min(max(n, max(cap(buf)*2, 4096)), arenaMaxUTF16Cap)
			buf = make([]uint16, n, newCap)
			a.utf16Buf = buf
		} else {
			buf = make([]uint16, n)
		}
	} else {
		buf = buf[:n]
	}
	for i := 0; i < prefix; i++ {
		buf[i] = uint16(s[i])
	}
	if prefix == n {
		if cap(buf) <= arenaMaxUTF16Cap {
			a.utf16Buf = buf
		}
		return buf, prefix
	}
	out := buf[:prefix]
	for i := prefix; i < n; {
		r, width := decodeWTF8InArena(s[i:])
		i += width
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			out = append(out, uint16(hi), uint16(lo))
		} else {
			out = append(out, uint16(r))
		}
	}
	if cap(out) <= arenaMaxUTF16Cap {
		a.utf16Buf = out
	}
	return out, prefix
}

func decodeWTF8InArena(s string) (r rune, size int) {
	c := s[0]
	if c < 0x80 {
		return rune(c), 1
	}
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

func (p *Parser) newSourceRange(start, end *SourcePosition) *SourceRange {
	r := p.arena.ranges.New()
	r.Start = start
	r.End = end
	return r
}

func (p *Parser) newProgramTree(location *SourceRange, sourceElements []ParseTree, sourceComments []*Comment) *ProgramTree {
	t := p.arena.programs.New()
	t.InitParseTreeBase(t, ParseTreeType_PROGRAM, location)
	t.SourceElements = sourceElements
	t.SourceComments = sourceComments
	return t
}

func (p *Parser) newFunctionBuilder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	b := p.arena.functionBuilders.New()
	b.kind = kind
	b.arena = p.arena
	return b
}

func (p *Parser) newFormalParameterListTree(location *SourceRange, parameters []ParseTree, hasTrailingComma bool, commaPositions []*SourcePosition) *FormalParameterListTree {
	t := p.arena.formalParameterLists.New()
	t.InitParseTreeBase(t, ParseTreeType_FORMAL_PARAMETER_LIST, location)
	numParams := len(parameters)
	numCommas := len(commaPositions)
	if !(numCommas <= numParams && numCommas >= numParams-1) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat(
			"Unexpected # of comma and formal params.\nparams: %s\ncomma positions: %s",
			jutil.CollectionToString(parameters),
			jutil.CollectionToString(commaPositions))))
	}
	t.Parameters = parameters
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}

func (p *Parser) newBlockTree(location *SourceRange, statements []ParseTree) *BlockTree {
	t := p.arena.blocks.New()
	t.InitParseTreeBase(t, ParseTreeType_BLOCK, location)
	t.Statements = statements
	return t
}

func (p *Parser) newVariableStatementTree(location *SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	t := p.arena.variableStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_STATEMENT, location)
	t.Declarations = declarations
	return t
}

func (p *Parser) newVariableDeclarationListTree(location *SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	t := p.arena.variableDeclarationLists.New()
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_DECLARATION_LIST, location)
	t.DeclarationType = declarationType
	t.Declarations = declarations
	return t
}

func (p *Parser) newVariableDeclarationTree(location *SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	t := p.arena.variableDeclarations.New()
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_DECLARATION, location)
	t.Lvalue = lvalue
	t.Initializer = initializer
	return t
}

func (p *Parser) newEmptyStatementTree(location *SourceRange) *EmptyStatementTree {
	t := p.arena.emptyStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_EMPTY_STATEMENT, location)
	return t
}

func (p *Parser) newExpressionStatementTree(location *SourceRange, expression ParseTree) *ExpressionStatementTree {
	t := p.arena.expressionStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_EXPRESSION_STATEMENT, location)
	t.Expression = expression
	return t
}

func (p *Parser) newIfStatementTree(location *SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	t := p.arena.ifStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_IF_STATEMENT, location)
	t.Condition = condition
	t.IfClause = ifClause
	t.ElseClause = elseClause
	return t
}

func (p *Parser) newNullTree(location *SourceRange) *NullTree {
	t := p.arena.nullTrees.New()
	t.InitParseTreeBase(t, ParseTreeType_NULL, location)
	return t
}

func (p *Parser) newReturnStatementTree(location *SourceRange, expression ParseTree) *ReturnStatementTree {
	t := p.arena.returnStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_RETURN_STATEMENT, location)
	t.Expression = expression
	return t
}

func (p *Parser) newThisExpressionTree(location *SourceRange) *ThisExpressionTree {
	t := p.arena.thisExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_THIS_EXPRESSION, location)
	return t
}

func (p *Parser) newIdentifierExpressionTree(location *SourceRange, identifierToken *IdentifierToken) *IdentifierExpressionTree {
	t := p.arena.identifierExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_IDENTIFIER_EXPRESSION, location)
	t.IdentifierToken = identifierToken
	return t
}

func (p *Parser) newLiteralExpressionTree(location *SourceRange, literalToken Token) *LiteralExpressionTree {
	t := p.arena.literalExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_LITERAL_EXPRESSION, location)
	t.LiteralToken = literalToken
	return t
}

func (p *Parser) newArrayLiteralExpressionTree(location *SourceRange, elements []ParseTree, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	t := p.arena.arrayLiteralExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_ARRAY_LITERAL_EXPRESSION, location)
	t.Elements = elements
	t.HasTrailingComma = hasTrailingComma
	return t
}

func (p *Parser) newObjectLiteralExpressionTree(location *SourceRange, propertyNameAndValues []ParseTree, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	t := p.arena.objectLiteralExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_OBJECT_LITERAL_EXPRESSION, location)
	t.PropertyNameAndValues = propertyNameAndValues
	t.HasTrailingComma = hasTrailingComma
	return t
}

func (p *Parser) newPropertyNameAssignmentTree(location *SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	t := p.arena.propertyNameAssignments.New()
	t.InitParseTreeBase(t, ParseTreeType_PROPERTY_NAME_ASSIGNMENT, location)
	t.Name = name
	t.Value = value
	return t
}

func (p *Parser) newCommaExpressionTree(location *SourceRange, expressions []ParseTree) *CommaExpressionTree {
	t := p.arena.commaExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_COMMA_EXPRESSION, location)
	t.Expressions = expressions
	return t
}

func (p *Parser) newParenExpressionTree(location *SourceRange, expression ParseTree) *ParenExpressionTree {
	t := p.arena.parenExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_PAREN_EXPRESSION, location)
	t.Expression = expression
	return t
}

func (p *Parser) newBinaryOperatorTree(location *SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	t := p.arena.binaryOperators.New()
	t.InitParseTreeBase(t, ParseTreeType_BINARY_OPERATOR, location)
	t.Left = left
	t.Operator = operator
	t.Right = right
	return t
}

func (p *Parser) newConditionalExpressionTree(location *SourceRange, condition ParseTree, left ParseTree, right ParseTree) *ConditionalExpressionTree {
	t := p.arena.conditionalExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_CONDITIONAL_EXPRESSION, location)
	t.Condition = condition
	t.Left = left
	t.Right = right
	return t
}

func (p *Parser) newUnaryExpressionTree(location *SourceRange, operator Token, operand ParseTree) *UnaryExpressionTree {
	t := p.arena.unaryExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_UNARY_EXPRESSION, location)
	t.Operator = operator
	t.Operand = operand
	return t
}

func (p *Parser) newUpdateExpressionTreePrefix(location *SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	t := p.arena.updateExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_UPDATE_EXPRESSION, location)
	t.Operator = operator
	t.OperatorPosition = UpdateExpressionTree_OperatorPosition_PREFIX
	t.Operand = operand
	return t
}

func (p *Parser) newUpdateExpressionTreePostfix(location *SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	t := p.arena.updateExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_UPDATE_EXPRESSION, location)
	t.Operator = operator
	t.OperatorPosition = UpdateExpressionTree_OperatorPosition_POSTFIX
	t.Operand = operand
	return t
}

func (p *Parser) newCallExpressionTree(location *SourceRange, operand ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	t := p.arena.callExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	return t
}

func (p *Parser) newMemberLookupExpressionTree(location *SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	t := p.arena.memberLookupExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	return t
}

func (p *Parser) newMemberExpressionTree(location *SourceRange, operand ParseTree, memberName *IdentifierToken) *MemberExpressionTree {
	t := p.arena.memberExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	return t
}

func (p *Parser) newNewExpressionTree(location *SourceRange, operand ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) *NewExpressionTree {
	t := p.arena.newExpressions.New()
	t.InitParseTreeBase(t, ParseTreeType_NEW_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	return t
}

func (p *Parser) newArgumentListTree(location *SourceRange, arguments []ParseTree, hasTrailingComma bool, commaPositions []*SourcePosition) *ArgumentListTree {
	t := p.arena.argumentLists.New()
	t.InitParseTreeBase(t, ParseTreeType_ARGUMENT_LIST, location)
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}

func (p *Parser) newClassDeclarationTree(location *SourceRange, name *IdentifierToken, superClass ParseTree, elements []ParseTree) *ClassDeclarationTree {
	t := p.arena.classDeclarations.New()
	t.InitParseTreeBase(t, ParseTreeType_CLASS_DECLARATION, location)
	t.Name = name
	t.SuperClass = superClass
	t.Elements = elements
	return t
}

func (p *Parser) newComputedPropertyMethodTree(location *SourceRange, property ParseTree, method ParseTree) *ComputedPropertyMethodTree {
	t := p.arena.computedPropertyMethods.New()
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_METHOD, location)
	t.Property = property
	t.Method = method
	return t
}

func (p *Parser) newWhileStatementTree(location *SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	t := p.arena.whileStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_WHILE_STATEMENT, location)
	t.Condition = condition
	t.Body = body
	return t
}

func (p *Parser) newForOfStatementTree(location *SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForOfStatementTree {
	t := p.arena.forOfStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_FOR_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}

func (p *Parser) newForStatementTree(location *SourceRange, variables ParseTree, condition ParseTree, increment ParseTree, body ParseTree) *ForStatementTree {
	t := p.arena.forStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_FOR_STATEMENT, location)
	t.Initializer = variables
	t.Condition = condition
	t.Increment = increment
	t.Body = body
	return t
}

func (p *Parser) newForInStatementTree(location *SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForInStatementTree {
	t := p.arena.forInStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_FOR_IN_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}

func (p *Parser) newContinueStatementTree(location *SourceRange, name *IdentifierToken) *ContinueStatementTree {
	t := p.arena.continueStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_CONTINUE_STATEMENT, location)
	t.Name = name
	return t
}

func (p *Parser) newBreakStatementTree(location *SourceRange, name *IdentifierToken) *BreakStatementTree {
	t := p.arena.breakStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_BREAK_STATEMENT, location)
	t.Name = name
	return t
}

func (p *Parser) newSwitchStatementTree(location *SourceRange, expression ParseTree, caseClauses []ParseTree) *SwitchStatementTree {
	t := p.arena.switchStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_SWITCH_STATEMENT, location)
	t.Expression = expression
	t.CaseClauses = caseClauses
	return t
}

func (p *Parser) newCaseClauseTree(location *SourceRange, expression ParseTree, statements []ParseTree) *CaseClauseTree {
	t := p.arena.caseClauses.New()
	t.InitParseTreeBase(t, ParseTreeType_CASE_CLAUSE, location)
	t.Expression = expression
	t.Statements = statements
	return t
}

func (p *Parser) newThrowStatementTree(location *SourceRange, value ParseTree) *ThrowStatementTree {
	t := p.arena.throwStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_THROW_STATEMENT, location)
	t.Value = value
	return t
}

func (p *Parser) newTryStatementTree(location *SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) *TryStatementTree {
	t := p.arena.tryStatements.New()
	t.InitParseTreeBase(t, ParseTreeType_TRY_STATEMENT, location)
	t.Body = body
	t.CatchBlock = catchBlock
	t.FinallyBlock = finallyBlock
	return t
}

func (p *Parser) newCatchTree(location *SourceRange, exception ParseTree, catchBody ParseTree) *CatchTree {
	t := p.arena.catchTrees.New()
	t.InitParseTreeBase(t, ParseTreeType_CATCH, location)
	t.Exception = exception
	t.CatchBody = catchBody
	return t
}
