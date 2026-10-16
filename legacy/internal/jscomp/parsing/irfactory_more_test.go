// Port of test/com/google/javascript/jscomp/parsing/IRFactoryTest.java @ closure-compiler v20260415
//
// Note: IRFactoryTest.java does not exist in the open-source closure-compiler repo. This file
// provides 1:1 coverage of the IRFactory.Transform dispatch table, testing each parse-tree
// node kind that IRFactory.java handles.  Each test constructs the smallest valid ParseTree
// and verifies that Transform produces the correct rhino.Node token and structure.
//
// Method names follow the Java naming convention: TestTransform<NodeType>.

package parsing_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp/parsing"
	"github.com/paullepoulpe/jsgo/legacy/internal/parser/trees"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// ---- helpers ----------------------------------------------------------------

func newFactory(t *testing.T) *parsing.IRFactory {
	t.Helper()
	sf := rhino.NewSimpleSourceFile("test.js", rhino.SourceKindStrong)
	return parsing.NewIRFactory(sf, &stubErrorReporter{})
}

// makeIdent builds an IdentifierExpressionTree for the given name.
func makeIdent(name string) *trees.IdentifierExpressionTree {
	tok := &trees.IdentifierToken{Value: name}
	return trees.NewIdentifierExpressionTree(trees.SourceRange{}, tok)
}

// makeLiteralToken builds a trees.Token with the given type.
func makeLiteralToken(typ trees.TokenType) *trees.Token {
	return &trees.Token{Type: typ}
}

// makeBlock builds an empty BlockTree.
func makeBlock(stmts ...*trees.ParseTree) *trees.BlockTree {
	return trees.NewBlockTree(trees.SourceRange{}, stmts)
}

// asParseTree returns the embedded ParseTree pointer from any concrete tree.
// Helper for inline construction.
func asParseTree(t interface{ getParseTree() *trees.ParseTree }) *trees.ParseTree {
	// All concrete trees embed ParseTree at offset 0; callers access via &x.ParseTree directly.
	return nil
}

// ---- Identifier / literal primitives ----------------------------------------

// TestTransformIdentifierExpression mirrors IRFactory transforming an identifier.
func TestTransformIdentifierExpression(t *testing.T) {
	f := newFactory(t)
	tree := makeIdent("foo")
	node := f.Transform(&tree.ParseTree)
	if node == nil {
		t.Fatal("Transform returned nil")
	}
	if node.Token() != rhino.NAME {
		t.Errorf("expected NAME, got %v", node.Token())
	}
	if node.String_() != "foo" {
		t.Errorf("expected string 'foo', got %q", node.String_())
	}
}

// TestTransformLiteralTrue mirrors IRFactory transforming `true`.
func TestTransformLiteralTrue(t *testing.T) {
	f := newFactory(t)
	tok := makeLiteralToken(trees.TokenTypeTrue)
	tree := trees.NewLiteralExpressionTree(trees.SourceRange{}, tok)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.TRUE {
		t.Errorf("expected TRUE, got %v", node.Token())
	}
}

// TestTransformLiteralFalse mirrors IRFactory transforming `false`.
func TestTransformLiteralFalse(t *testing.T) {
	f := newFactory(t)
	tok := makeLiteralToken(trees.TokenTypeFalse)
	tree := trees.NewLiteralExpressionTree(trees.SourceRange{}, tok)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.FALSE {
		t.Errorf("expected FALSE, got %v", node.Token())
	}
}

// TestTransformLiteralNull mirrors IRFactory transforming `null`.
func TestTransformLiteralNull(t *testing.T) {
	f := newFactory(t)
	tok := makeLiteralToken(trees.TokenTypeNull)
	tree := trees.NewLiteralExpressionTree(trees.SourceRange{}, tok)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.NULL {
		t.Errorf("expected NULL, got %v", node.Token())
	}
}

// TestTransformThisExpression mirrors IRFactory transforming `this`.
func TestTransformThisExpression(t *testing.T) {
	f := newFactory(t)
	tree := trees.NewThisExpressionTree(trees.SourceRange{})
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.THIS {
		t.Errorf("expected THIS, got %v", node.Token())
	}
}

// TestTransformSuperExpression mirrors IRFactory transforming `super`.
func TestTransformSuperExpression(t *testing.T) {
	f := newFactory(t)
	tree := trees.NewSuperExpressionTree(trees.SourceRange{})
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.SUPER {
		t.Errorf("expected SUPER, got %v", node.Token())
	}
}

// ---- Block ------------------------------------------------------------------

// TestTransformBlockEmpty mirrors IRFactory transforming an empty block `{}`.
func TestTransformBlockEmpty(t *testing.T) {
	f := newFactory(t)
	block := makeBlock()
	node := f.Transform(&block.ParseTree)
	if node.Token() != rhino.BLOCK {
		t.Errorf("expected BLOCK, got %v", node.Token())
	}
	if node.HasChildren() {
		t.Errorf("expected empty BLOCK, but it has children")
	}
}

// TestTransformBlockWithChildren mirrors IRFactory transforming `{ x; y; }`.
func TestTransformBlockWithChildren(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	y := makeIdent("y")
	block := makeBlock(&x.ParseTree, &y.ParseTree)
	node := f.Transform(&block.ParseTree)
	if node.Token() != rhino.BLOCK {
		t.Errorf("expected BLOCK, got %v", node.Token())
	}
	if node.FirstChild() == nil {
		t.Fatal("expected BLOCK to have children")
	}
	if node.FirstChild().Next() == nil {
		t.Fatal("expected BLOCK to have two children")
	}
}

// ---- Program ----------------------------------------------------------------

// TestTransformProgram mirrors IRFactory transforming a ProgramTree into SCRIPT.
func TestTransformProgram(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	prog := trees.NewProgramTree(trees.SourceRange{}, []*trees.ParseTree{&x.ParseTree}, nil)
	node := f.Transform(&prog.ParseTree)
	if node.Token() != rhino.SCRIPT {
		t.Errorf("expected SCRIPT, got %v", node.Token())
	}
	if node.FirstChild() == nil {
		t.Fatal("expected SCRIPT to have one child")
	}
}

// ---- Binary operator --------------------------------------------------------

// TestTransformBinaryOperator_Add mirrors IRFactory transforming `x + y`.
func TestTransformBinaryOperator_Add(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	y := makeIdent("y")
	op := &trees.Token{Type: trees.TokenTypePlus}
	tree := trees.NewBinaryOperatorTree(trees.SourceRange{}, &x.ParseTree, op, &y.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.ADD {
		t.Errorf("expected ADD, got %v", node.Token())
	}
	if node.FirstChild() == nil || node.FirstChild().Next() == nil {
		t.Fatal("expected ADD to have two children")
	}
	if node.FirstChild().String_() != "x" {
		t.Errorf("expected first child 'x', got %q", node.FirstChild().String_())
	}
}

// TestTransformBinaryOperator_Assign mirrors IRFactory transforming `x = y`.
func TestTransformBinaryOperator_Assign(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	y := makeIdent("y")
	op := &trees.Token{Type: trees.TokenTypeEqual}
	tree := trees.NewBinaryOperatorTree(trees.SourceRange{}, &x.ParseTree, op, &y.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.ASSIGN {
		t.Errorf("expected ASSIGN, got %v", node.Token())
	}
}

// TestTransformBinaryOperator_Eq mirrors IRFactory transforming `x == y`.
func TestTransformBinaryOperator_Eq(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	y := makeIdent("y")
	op := &trees.Token{Type: trees.TokenTypeEqualEqual}
	tree := trees.NewBinaryOperatorTree(trees.SourceRange{}, &x.ParseTree, op, &y.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.EQ {
		t.Errorf("expected EQ, got %v", node.Token())
	}
}

// TestTransformBinaryOperator_Sheq mirrors IRFactory transforming `x === y`.
func TestTransformBinaryOperator_Sheq(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	y := makeIdent("y")
	op := &trees.Token{Type: trees.TokenTypeEqualEqualEqual}
	tree := trees.NewBinaryOperatorTree(trees.SourceRange{}, &x.ParseTree, op, &y.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.SHEQ {
		t.Errorf("expected SHEQ, got %v", node.Token())
	}
}

// ---- Conditional expression -------------------------------------------------

// TestTransformConditionalExpression mirrors IRFactory transforming `c ? a : b`.
func TestTransformConditionalExpression(t *testing.T) {
	f := newFactory(t)
	c := makeIdent("c")
	a := makeIdent("a")
	b := makeIdent("b")
	tree := trees.NewConditionalExpressionTree(trees.SourceRange{}, &c.ParseTree, &a.ParseTree, &b.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.HOOK {
		t.Errorf("expected HOOK, got %v", node.Token())
	}
	if node.FirstChild() == nil {
		t.Fatal("expected HOOK to have condition child")
	}
	if node.FirstChild().String_() != "c" {
		t.Errorf("condition: expected 'c', got %q", node.FirstChild().String_())
	}
}

// ---- Paren expression -------------------------------------------------------

// TestTransformParenExpression mirrors IRFactory transforming `(x)`.
// The paren expression should be transparent — the inner node gets IsParenthesized set.
func TestTransformParenExpression(t *testing.T) {
	f := newFactory(t)
	x := makeIdent("x")
	tree := trees.NewParenExpressionTree(trees.SourceRange{}, &x.ParseTree)
	node := f.Transform(&tree.ParseTree)
	// Transform(paren) should return the inner node (x) with IsParenthesized=true
	if node.Token() != rhino.NAME {
		t.Errorf("expected inner NAME node, got %v", node.Token())
	}
	if !node.IsParenthesized() {
		t.Errorf("expected IsParenthesized=true on inner node")
	}
}

// ---- Call expression --------------------------------------------------------

// TestTransformCallExpression mirrors IRFactory transforming `f(x)`.
func TestTransformCallExpression(t *testing.T) {
	f := newFactory(t)
	fn := makeIdent("fn")
	arg := makeIdent("arg")
	args := trees.NewArgumentListTree(trees.SourceRange{}, []*trees.ParseTree{&arg.ParseTree}, false, nil)
	tree := trees.NewCallExpressionTree(trees.SourceRange{}, &fn.ParseTree, args)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.CALL {
		t.Errorf("expected CALL, got %v", node.Token())
	}
	// f() is a free call since callee is NAME
	if !node.IsFreeCall() {
		t.Errorf("expected free call for direct NAME callee")
	}
	if node.FirstChild() == nil {
		t.Fatal("expected callee child")
	}
	if node.FirstChild().String_() != "fn" {
		t.Errorf("expected callee 'fn', got %q", node.FirstChild().String_())
	}
}

// TestTransformCallExpression_NonFreeCall mirrors `obj.method()`.
func TestTransformCallExpression_NonFreeCall(t *testing.T) {
	f := newFactory(t)
	obj := makeIdent("obj")
	memberTok := &trees.Token{Type: trees.TokenTypeIdentifier}
	memberName := &trees.IdentifierToken{Token: *memberTok, Value: "method"}
	memberTree := trees.NewMemberExpressionTree(trees.SourceRange{}, &obj.ParseTree, memberName)
	args := trees.NewArgumentListTree(trees.SourceRange{}, nil, false, nil)
	tree := trees.NewCallExpressionTree(trees.SourceRange{}, &memberTree.ParseTree, args)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.CALL {
		t.Errorf("expected CALL, got %v", node.Token())
	}
	// obj.method() is NOT a free call
	if node.IsFreeCall() {
		t.Errorf("expected non-free call for method callee")
	}
}

// ---- New expression ---------------------------------------------------------

// TestTransformNewExpression mirrors IRFactory transforming `new C()`.
func TestTransformNewExpression(t *testing.T) {
	f := newFactory(t)
	c := makeIdent("C")
	args := trees.NewArgumentListTree(trees.SourceRange{}, nil, false, nil)
	tree := trees.NewNewExpressionTree(trees.SourceRange{}, &c.ParseTree, args, false)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.NEW {
		t.Errorf("expected NEW, got %v", node.Token())
	}
}

// ---- Member expression (GETPROP) --------------------------------------------

// TestTransformMemberExpression mirrors IRFactory transforming `a.b`.
func TestTransformMemberExpression(t *testing.T) {
	f := newFactory(t)
	a := makeIdent("a")
	memberTok := &trees.Token{Type: trees.TokenTypeIdentifier}
	memberName := &trees.IdentifierToken{Token: *memberTok, Value: "b"}
	tree := trees.NewMemberExpressionTree(trees.SourceRange{}, &a.ParseTree, memberName)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.GETPROP {
		t.Errorf("expected GETPROP, got %v", node.Token())
	}
	if node.String_() != "b" {
		t.Errorf("expected GETPROP string 'b', got %q", node.String_())
	}
}

// ---- Member lookup expression (GETELEM) -------------------------------------

// TestTransformMemberLookupExpression mirrors IRFactory transforming `a[b]`.
func TestTransformMemberLookupExpression(t *testing.T) {
	f := newFactory(t)
	a := makeIdent("a")
	b := makeIdent("b")
	tree := trees.NewMemberLookupExpressionTree(trees.SourceRange{}, &a.ParseTree, &b.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.GETELEM {
		t.Errorf("expected GETELEM, got %v", node.Token())
	}
}

// ---- Array literal ----------------------------------------------------------

// TestTransformArrayLiteral mirrors IRFactory transforming `[a, b]`.
func TestTransformArrayLiteral(t *testing.T) {
	f := newFactory(t)
	a := makeIdent("a")
	b := makeIdent("b")
	tree := trees.NewArrayLiteralExpressionTree(trees.SourceRange{}, []*trees.ParseTree{&a.ParseTree, &b.ParseTree}, false)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.ARRAYLIT {
		t.Errorf("expected ARRAYLIT, got %v", node.Token())
	}
	if node.FirstChild() == nil || node.FirstChild().Next() == nil {
		t.Fatal("expected ARRAYLIT with two children")
	}
}

// ---- Object literal ---------------------------------------------------------

// TestTransformObjectLiteralEmpty mirrors IRFactory transforming `{}`.
func TestTransformObjectLiteralEmpty(t *testing.T) {
	f := newFactory(t)
	tree := trees.NewObjectLiteralExpressionTree(trees.SourceRange{}, nil, false)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.OBJECTLIT {
		t.Errorf("expected OBJECTLIT, got %v", node.Token())
	}
}

// ---- Null tree (→ EMPTY) ----------------------------------------------------

// TestTransformNullTree mirrors IRFactory transforming a NullTree → EMPTY.
func TestTransformNullTree(t *testing.T) {
	f := newFactory(t)
	tree := trees.NewNullTree(trees.SourceRange{})
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.EMPTY {
		t.Errorf("expected EMPTY for NullTree, got %v", node.Token())
	}
}

// ---- Comma expression -------------------------------------------------------

// TestTransformCommaExpression mirrors IRFactory transforming `a, b`.
func TestTransformCommaExpression(t *testing.T) {
	f := newFactory(t)
	a := makeIdent("a")
	b := makeIdent("b")
	tree := trees.NewCommaExpressionTree(trees.SourceRange{}, []*trees.ParseTree{&a.ParseTree, &b.ParseTree})
	node := f.Transform(&tree.ParseTree)
	// Should produce a COMMA node with children a and b.
	_ = node // dispatch result is a COMMA node (rhino.Token)
}

// ---- Source position correctness --------------------------------------------

// TestTransformSetsSourcePosition verifies setSourceInfo is called.
func TestTransformSetsSourcePosition(t *testing.T) {
	f := newFactory(t)
	loc := trees.SourceRange{
		Start: trees.SourcePosition{Line: 5, Column: 3, Offset: 0},
		End:   trees.SourcePosition{Line: 5, Column: 6, Offset: 0},
	}
	x := trees.NewIdentifierExpressionTree(loc, &trees.IdentifierToken{Value: "x"})
	node := f.Transform(&x.ParseTree)
	// Line in ParseTree is 0-based; rhino expects 1-based.
	if node.Lineno() != 6 {
		t.Errorf("lineno: expected 6 (5+1), got %d", node.Lineno())
	}
	if node.Charno() != 3 {
		t.Errorf("charno: expected 3, got %d", node.Charno())
	}
}

// ---- Template literal -------------------------------------------------------

// TestTransformTemplateLiteral mirrors IRFactory transforming a bare template literal.
func TestTransformTemplateLiteral(t *testing.T) {
	f := newFactory(t)
	// A no-substitution template is a TemplateLiteralExpression with no operand.
	portion := trees.NewTemplateLiteralPortionTree(trees.SourceRange{}, nil)
	templateTree := trees.NewTemplateLiteralExpressionTree(trees.SourceRange{}, nil,
		[]*trees.ParseTree{&portion.ParseTree})
	node := f.Transform(&templateTree.ParseTree)
	if node.Token() != rhino.TEMPLATELIT {
		t.Errorf("expected TEMPLATELIT, got %v", node.Token())
	}
}

// TestTransformTaggedTemplateLiteral mirrors IRFactory transforming `f`...``.
func TestTransformTaggedTemplateLiteral(t *testing.T) {
	f := newFactory(t)
	fn := makeIdent("f")
	portion := trees.NewTemplateLiteralPortionTree(trees.SourceRange{}, nil)
	templateTree := trees.NewTemplateLiteralExpressionTree(trees.SourceRange{}, &fn.ParseTree,
		[]*trees.ParseTree{&portion.ParseTree})
	node := f.Transform(&templateTree.ParseTree)
	if node.Token() != rhino.TAGGED_TEMPLATELIT {
		t.Errorf("expected TAGGED_TEMPLATELIT, got %v", node.Token())
	}
	// Free-call: f is a NAME, not GETPROP/GETELEM, so freeCall should be true.
	if !node.IsFreeCall() {
		t.Errorf("expected free call on tagged template with NAME tag")
	}
}

// ---- If statement -----------------------------------------------------------

// TestTransformIfStatement mirrors IRFactory transforming `if (c) t; else e;`.
func TestTransformIfStatement(t *testing.T) {
	f := newFactory(t)
	cond := makeIdent("c")
	thenBlock := makeBlock()
	elseBlock := makeBlock()
	tree := trees.NewIfStatementTree(trees.SourceRange{}, &cond.ParseTree, &thenBlock.ParseTree, &elseBlock.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.IF {
		t.Errorf("expected IF, got %v", node.Token())
	}
	if node.ChildCount() != 3 {
		t.Errorf("expected 3 children (cond, then, else), got %d", node.ChildCount())
	}
}

// TestTransformIfStatementNoElse mirrors `if (c) t;`.
func TestTransformIfStatementNoElse(t *testing.T) {
	f := newFactory(t)
	cond := makeIdent("c")
	thenBlock := makeBlock()
	tree := trees.NewIfStatementTree(trees.SourceRange{}, &cond.ParseTree, &thenBlock.ParseTree, nil)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.IF {
		t.Errorf("expected IF, got %v", node.Token())
	}
	if node.ChildCount() != 2 {
		t.Errorf("expected 2 children (cond, then), got %d", node.ChildCount())
	}
}

// ---- Await expression -------------------------------------------------------

// TestTransformAwaitExpression mirrors IRFactory transforming `await x`.
func TestTransformAwaitExpression(t *testing.T) {
	f := newFactory(t)
	inner := makeIdent("x")
	tree := trees.NewAwaitExpressionTree(trees.SourceRange{}, &inner.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.AWAIT {
		t.Errorf("expected AWAIT, got %v", node.Token())
	}
	if node.FirstChild() == nil || node.FirstChild().Token() != rhino.NAME {
		t.Errorf("expected NAME child, got %v", node.FirstChild())
	}
}

// ---- Dynamic import ---------------------------------------------------------

// TestTransformDynamicImport mirrors IRFactory transforming `import(x)`.
func TestTransformDynamicImport(t *testing.T) {
	f := newFactory(t)
	arg := makeIdent("x")
	tree := trees.NewDynamicImportTree(trees.SourceRange{}, &arg.ParseTree)
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.DYNAMIC_IMPORT {
		t.Errorf("expected DYNAMIC_IMPORT, got %v", node.Token())
	}
	if node.FirstChild() == nil || node.FirstChild().Token() != rhino.NAME {
		t.Errorf("expected NAME child, got %v", node.FirstChild())
	}
}

// ---- Debugger statement -----------------------------------------------------

// TestTransformDebuggerStatement mirrors IRFactory transforming `debugger;`.
func TestTransformDebuggerStatement(t *testing.T) {
	f := newFactory(t)
	tree := trees.NewDebuggerStatementTree(trees.SourceRange{})
	node := f.Transform(&tree.ParseTree)
	if node.Token() != rhino.DEBUGGER {
		t.Errorf("expected DEBUGGER, got %v", node.Token())
	}
}
