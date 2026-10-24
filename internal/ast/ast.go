// Package ast is a typed view of the JavaScript AST.
//
// The compiler stores the AST as rhino.Nodes: every node has a token, a
// list of children and generic accessors (GetFirstChild, GetSecondChild,
// GetString...), so what a child means is a matter of convention: the
// second child of an IF is its then-block. This package gives each node
// kind a Go type whose methods name those children:
//
//	ifs, ok := ast.AsIf(n)  // ok is false if n is not an if statement
//	cond := ifs.Cond()      // an ast.Expr
//	if els, ok := ifs.Else(); ok {
//		...                 // els is an ast.Block
//	}
//
// A category (Expr, Statement, Target...) is the set of kinds that can
// appear in one place; As methods convert between kinds and categories,
// checking the token when the conversion can fail.
//
// The types are views: each holds the rhino.Node it describes (4 bytes),
// so they cost nothing to create or pass, and Node() goes back to the
// rhino API for anything this package does not cover.
//
// The kinds, their children and the categories are listed in grammar.txt;
// nodes.gen.go is generated from it (go generate ./internal/ast). Check
// verifies a tree against the grammar.
package ast

//go:generate go run ../../tools/astgen -grammar grammar.txt -out nodes.gen.go
