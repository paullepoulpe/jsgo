// Port of com/google/javascript/jscomp/parsing/parser/Parser.java @ closure-compiler v20261028

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// parseTreeSlabs groups the per-AST-node slab allocators and shared child-list
// chunks for Parser. A statement's parse trees are created together and die
// together once converted; zeroing p.parseTreeSlabs resets all 20 slabs and
// chunk slices in a single assignment.
type parseTreeSlabs struct {
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

func (p *Parser) resetSlabs() {
	p.parseTreeSlabs = parseTreeSlabs{}
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

func (p *Parser) popTrees(base int) []ParseTree {
	return popStackChunked(&p.treeStack, base, &p.treeLists, parser_treeListChunk)
}

func (p *Parser) popPos(base int) []SourcePosition {
	return popStackChunked(&p.posStack, base, &p.posLists, parser_posListChunk)
}

func (p *Parser) popVarDecls(base int) []*VariableDeclarationTree {
	return popStackChunked(&p.varDeclStack, base, &p.varDeclLists, parser_treeListChunk)
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

func allocTree[T any](s *jutil.Slab[T], v T) *T {
	t := s.New()
	*t = v
	return t
}

func baseTree(type_ ParseTreeType, loc SourceRange) ParseTreeBase {
	return ParseTreeBase{Type: type_, Location: loc}
}

func (p *Parser) newBlockTree(loc SourceRange, statements []ParseTree) *BlockTree {
	return allocTree(&p.blocks, BlockTree{ParseTreeBase: baseTree(ParseTreeType_BLOCK, loc), Statements: statements})
}

func (p *Parser) newVariableDeclarationTree(loc SourceRange, lvalue, initializer ParseTree) *VariableDeclarationTree {
	return allocTree(&p.varDecls, VariableDeclarationTree{ParseTreeBase: baseTree(ParseTreeType_VARIABLE_DECLARATION, loc), Lvalue: lvalue, Initializer: initializer})
}

func (p *Parser) newVariableDeclarationListTree(loc SourceRange, declType TokenType, decls []*VariableDeclarationTree) *VariableDeclarationListTree {
	return allocTree(&p.varDeclLsts, VariableDeclarationListTree{ParseTreeBase: baseTree(ParseTreeType_VARIABLE_DECLARATION_LIST, loc), DeclarationType: declType, Declarations: decls})
}

func (p *Parser) newVariableStatementTree(loc SourceRange, decls *VariableDeclarationListTree) *VariableStatementTree {
	return allocTree(&p.varStmts, VariableStatementTree{ParseTreeBase: baseTree(ParseTreeType_VARIABLE_STATEMENT, loc), Declarations: decls})
}

func (p *Parser) newMemberLookupExpressionTree(loc SourceRange, operand, memberExpr ParseTree) *MemberLookupExpressionTree {
	return allocTree(&p.memberLkps, MemberLookupExpressionTree{ParseTreeBase: baseTree(ParseTreeType_MEMBER_LOOKUP_EXPRESSION, loc), Operand: operand, MemberExpression: memberExpr})
}

func (p *Parser) newThisExpressionTree(loc SourceRange) *ThisExpressionTree {
	return allocTree(&p.thisExprs, ThisExpressionTree{ParseTreeBase: baseTree(ParseTreeType_THIS_EXPRESSION, loc)})
}

func (p *Parser) newExpressionStatementTree(loc SourceRange, expr ParseTree) *ExpressionStatementTree {
	return allocTree(&p.exprStmts, ExpressionStatementTree{ParseTreeBase: baseTree(ParseTreeType_EXPRESSION_STATEMENT, loc), Expression: expr})
}

func (p *Parser) newCallExpressionTree(loc SourceRange, operand ParseTree, args *ArgumentListTree) *CallExpressionTree {
	return allocTree(&p.calls, CallExpressionTree{ParseTreeBase: baseTree(ParseTreeType_CALL_EXPRESSION, loc), Operand: operand, Arguments: args})
}

func (p *Parser) newArgumentListTree(loc SourceRange, args []ParseTree, trailingComma bool, commas []SourcePosition) *ArgumentListTree {
	return allocTree(&p.argLists, ArgumentListTree{ParseTreeBase: baseTree(ParseTreeType_ARGUMENT_LIST, loc), Arguments: args, HasTrailingComma: trailingComma, CommaPositions: commas})
}

func (p *Parser) newIdentifierExpressionTree(loc SourceRange, tok Token) *IdentifierExpressionTree {
	return allocTree(&p.identExprs, IdentifierExpressionTree{ParseTreeBase: baseTree(ParseTreeType_IDENTIFIER_EXPRESSION, loc), IdentifierToken: tok})
}

func (p *Parser) newMemberExpressionTree(loc SourceRange, operand ParseTree, memberName Token) *MemberExpressionTree {
	return allocTree(&p.memberExprs, MemberExpressionTree{ParseTreeBase: baseTree(ParseTreeType_MEMBER_EXPRESSION, loc), Operand: operand, MemberName: memberName})
}

func (p *Parser) newBinaryOperatorTree(loc SourceRange, left ParseTree, op Token, right ParseTree) *BinaryOperatorTree {
	return allocTree(&p.binaryOps, BinaryOperatorTree{ParseTreeBase: baseTree(ParseTreeType_BINARY_OPERATOR, loc), Left: left, Operator: op, Right: right})
}

func (p *Parser) newLiteralExpressionTree(loc SourceRange, tok Token) *LiteralExpressionTree {
	return allocTree(&p.literals, LiteralExpressionTree{ParseTreeBase: baseTree(ParseTreeType_LITERAL_EXPRESSION, loc), LiteralToken: tok})
}

func (p *Parser) newPropertyNameAssignmentTree(loc SourceRange, name Token, val ParseTree) *PropertyNameAssignmentTree {
	return allocTree(&p.propAssigns, PropertyNameAssignmentTree{ParseTreeBase: baseTree(ParseTreeType_PROPERTY_NAME_ASSIGNMENT, loc), Name: name, Value: val})
}

func (p *Parser) newFunctionDeclarationTree(
	loc SourceRange, name Token, params *FormalParameterListTree, body ParseTree,
	kind FunctionDeclarationTree_Kind, isClassMember, isStatic, isGenerator, isAsync bool,
) *FunctionDeclarationTree {
	if params == nil || body == nil {
		panic(jutil.NewNullPointerException())
	}
	return allocTree(&p.funcDecls, FunctionDeclarationTree{
		ParseTreeBase: baseTree(ParseTreeType_FUNCTION_DECLARATION, loc),
		Name:          name, FormalParameterList: params, FunctionBody: body,
		IsClassMember: isClassMember, IsStatic: isStatic, IsGenerator: isGenerator, IsAsync: isAsync, Kind: kind,
	})
}

func (p *Parser) newFormalParameterListTree(loc SourceRange, params []ParseTree, trailingComma bool, commas []SourcePosition) *FormalParameterListTree {
	checkFormalParameterListCommas(params, commas)
	return allocTree(&p.paramLists, FormalParameterListTree{ParseTreeBase: baseTree(ParseTreeType_FORMAL_PARAMETER_LIST, loc), Parameters: params, HasTrailingComma: trailingComma, CommaPositions: commas})
}

func (p *Parser) newReturnStatementTree(loc SourceRange, expr ParseTree) *ReturnStatementTree {
	return allocTree(&p.returnStmts, ReturnStatementTree{ParseTreeBase: baseTree(ParseTreeType_RETURN_STATEMENT, loc), Expression: expr})
}

func (p *Parser) newIfStatementTree(loc SourceRange, cond, ifClause, elseClause ParseTree) *IfStatementTree {
	return allocTree(&p.ifStmts, IfStatementTree{ParseTreeBase: baseTree(ParseTreeType_IF_STATEMENT, loc), Condition: cond, IfClause: ifClause, ElseClause: elseClause})
}

func (p *Parser) newParenExpressionTree(loc SourceRange, expr ParseTree) *ParenExpressionTree {
	return allocTree(&p.parenExprs, ParenExpressionTree{ParseTreeBase: baseTree(ParseTreeType_PAREN_EXPRESSION, loc), Expression: expr})
}

func (p *Parser) newUnaryExpressionTree(loc SourceRange, op Token, operand ParseTree) *UnaryExpressionTree {
	return allocTree(&p.unaryExprs, UnaryExpressionTree{ParseTreeBase: baseTree(ParseTreeType_UNARY_EXPRESSION, loc), Operator: op, Operand: operand})
}
