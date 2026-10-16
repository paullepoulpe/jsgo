package parsing_test

import (
	"testing"
	"unsafe"

	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp/parsing"
	"github.com/paullepoulpe/jsgo/legacy/internal/parser/trees"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// stubErrorReporter is a minimal errorReporter for tests.
type stubErrorReporter struct{}

func (s *stubErrorReporter) Error(msg, sourceName string, lineno, charno int)   {}
func (s *stubErrorReporter) Warning(msg, sourceName string, lineno, charno int) {}

// TestAsParseTreeInvariant verifies that &concrete.ParseTree == pt when pt is
// the pointer returned by asParseTree (or equivalently, &tree.ParseTree directly).
// This is the invariant required by IRFactory's unsafe.Pointer dispatch:
//
//	(*BlockTree)(unsafe.Pointer(pt)) must equal the original *BlockTree.
func TestAsParseTreeInvariant(t *testing.T) {
	t.Run("BlockTree", func(t *testing.T) {
		loc := trees.SourceRange{}
		orig := trees.NewBlockTree(loc, nil)
		pt := &orig.ParseTree

		// The unsafe cast must recover the original pointer.
		recovered := (*trees.BlockTree)(unsafe.Pointer(pt))
		if recovered != orig {
			t.Fatalf("unsafe.Pointer round-trip failed: got %p, want %p", recovered, orig)
		}
	})

	t.Run("BinaryOperatorTree", func(t *testing.T) {
		loc := trees.SourceRange{}
		orig := trees.NewBinaryOperatorTree(loc, nil, nil, nil)
		pt := &orig.ParseTree

		recovered := (*trees.BinaryOperatorTree)(unsafe.Pointer(pt))
		if recovered != orig {
			t.Fatalf("unsafe.Pointer round-trip failed: got %p, want %p", recovered, orig)
		}
	})

	t.Run("IfStatementTree", func(t *testing.T) {
		loc := trees.SourceRange{}
		orig := trees.NewIfStatementTree(loc, nil, nil, nil)
		pt := &orig.ParseTree

		recovered := (*trees.IfStatementTree)(unsafe.Pointer(pt))
		if recovered != orig {
			t.Fatalf("unsafe.Pointer round-trip failed: got %p, want %p", recovered, orig)
		}
	})

	t.Run("ExpressionStatementTree", func(t *testing.T) {
		loc := trees.SourceRange{}
		orig := trees.NewExpressionStatementTree(loc, nil)
		pt := &orig.ParseTree

		recovered := (*trees.ExpressionStatementTree)(unsafe.Pointer(pt))
		if recovered != orig {
			t.Fatalf("unsafe.Pointer round-trip failed: got %p, want %p", recovered, orig)
		}
	})

	t.Run("VariableDeclarationTree", func(t *testing.T) {
		loc := trees.SourceRange{}
		orig := trees.NewVariableDeclarationTree(loc, nil, nil)
		pt := &orig.ParseTree

		recovered := (*trees.VariableDeclarationTree)(unsafe.Pointer(pt))
		if recovered != orig {
			t.Fatalf("unsafe.Pointer round-trip failed: got %p, want %p", recovered, orig)
		}
	})
}

// TestIRFactoryTransformBlock exercises IRFactory.Transform on a *BlockTree
// passed as &tree.ParseTree, verifying that dispatch works end-to-end.
func TestIRFactoryTransformBlock(t *testing.T) {
	sf := rhino.NewSimpleSourceFile("test.js", rhino.SourceKindStrong)
	factory := parsing.NewIRFactory(sf, &stubErrorReporter{})

	// Build a BlockTree with one identifier child.
	idLoc := trees.SourceRange{}
	idTok := &trees.IdentifierToken{Value: "x"}
	idTree := trees.NewIdentifierExpressionTree(idLoc, idTok)

	blockLoc := trees.SourceRange{}
	block := trees.NewBlockTree(blockLoc, []*trees.ParseTree{&idTree.ParseTree})

	// Pass the embedded ParseTree pointer — this is the invariant.
	pt := &block.ParseTree
	if pt.Type != trees.ParseTreeTypeBlock {
		t.Fatalf("expected ParseTreeTypeBlock, got %v", pt.Type)
	}

	node := factory.Transform(pt)
	if node == nil {
		t.Fatal("Transform returned nil for BlockTree")
	}
	if node.Token() != rhino.BLOCK {
		t.Fatalf("expected BLOCK node, got %v", node.Token())
	}
	if node.FirstChild() == nil {
		t.Fatal("expected one child in BLOCK node")
	}
}
