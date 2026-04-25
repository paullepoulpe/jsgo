// Package javaast provides a lightweight Java statement AST, a
// recursive-descent parser for method bodies extracted by javaport,
// and a translator that emits equivalent Go source text.
//
// The parser handles roughly 70% of real Closure Compiler method
// bodies: getters/setters, delegation, null-checks, simple control
// flow (if/else, for, while, switch), local vars, assignments, return,
// throw, try/catch/finally, break/continue, and object construction.
// Anything it cannot recognise falls back to a panic stub so the Go
// file still compiles.
package javaast

// ---- Stmt nodes --------------------------------------------------------

// Stmt is the sealed interface for all statement nodes.
type Stmt interface{ stmtNode() }

type (
	// BlockStmt is a `{ ... }` sequence of statements.
	BlockStmt struct{ Stmts []Stmt }

	// ExprStmt is a bare expression used as a statement.
	ExprStmt struct{ X Expr }

	// ReturnStmt represents `return [expr];`.
	ReturnStmt struct{ X Expr } // X may be nil for void returns

	// AssignStmt covers `lhs = rhs`, `lhs += rhs`, etc.
	AssignStmt struct {
		LHS Expr
		Op  string // "=", "+=", "-=", "|=", "&=", "^=", "*=", "/=", "%=", "<<=", ">>=", ">>>=", "||=", "&&="
		RHS Expr
	}

	// LocalVarStmt is `Type name [= init] ;`.
	LocalVarStmt struct {
		Type string
		Name string
		Init Expr // nil if absent
	}

	// IfStmt represents `if (cond) then [else alt]`.
	IfStmt struct {
		Cond Expr
		Then Stmt
		Else Stmt // nil if absent; may itself be an IfStmt (else-if chain)
	}

	// ForStmt is a classic `for (init; cond; update) body` loop.
	ForStmt struct {
		Init Stmt // LocalVarStmt, AssignStmt, ExprStmt, or nil
		Cond Expr // nil → for {}
		Post Stmt // AssignStmt, ExprStmt, or nil
		Body Stmt
	}

	// ForEachStmt is `for (Type name : coll) body`.
	ForEachStmt struct {
		Type string
		Name string
		Coll Expr
		Body Stmt
	}

	// WhileStmt represents `while (cond) body`.
	WhileStmt struct {
		Cond Expr
		Body Stmt
	}

	// DoWhileStmt represents `do body while (cond);`.
	DoWhileStmt struct {
		Body Stmt
		Cond Expr
	}

	// SwitchStmt represents `switch (x) { case a: ... default: ... }`.
	SwitchStmt struct {
		X     Expr
		Cases []SwitchCase
	}

	// TryStmt represents `try { body } catch(T e) { handler } finally { fin }`.
	TryStmt struct {
		Body    Stmt
		Catches []CatchClause
		Finally Stmt // nil if absent
	}

	// ThrowStmt represents `throw expr;`.
	ThrowStmt struct{ X Expr }

	// BreakStmt represents `break [label];`.
	BreakStmt struct{ Label string }

	// ContinueStmt represents `continue [label];`.
	ContinueStmt struct{ Label string }

	// LabeledStmt represents `label: stmt`.
	LabeledStmt struct {
		Label string
		Stmt  Stmt
	}

	// RawStmt holds an untranslatable Java snippet — emitted as a
	// `// TODO: <text>` comment in Go.
	RawStmt struct{ Text string }
)

func (*BlockStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()     {}
func (*ReturnStmt) stmtNode()   {}
func (*AssignStmt) stmtNode()   {}
func (*LocalVarStmt) stmtNode() {}
func (*IfStmt) stmtNode()       {}
func (*ForStmt) stmtNode()      {}
func (*ForEachStmt) stmtNode()  {}
func (*WhileStmt) stmtNode()    {}
func (*DoWhileStmt) stmtNode()  {}
func (*SwitchStmt) stmtNode()   {}
func (*TryStmt) stmtNode()      {}
func (*ThrowStmt) stmtNode()    {}
func (*BreakStmt) stmtNode()    {}
func (*ContinueStmt) stmtNode() {}
func (*LabeledStmt) stmtNode()  {}
func (*RawStmt) stmtNode()      {}

// SwitchCase holds one `case X:` or `default:` group.
type SwitchCase struct {
	Value    Expr   // nil → default
	Stmts    []Stmt
	FallThru bool // true when no break/return at end → emit fallthrough
}

// CatchClause is one `catch (Type name) { body }` block.
type CatchClause struct {
	Type string
	Name string
	Body Stmt
}

// ---- Expr nodes --------------------------------------------------------

// Expr is the sealed interface for expression nodes.
type Expr interface{ exprNode() }

type (
	// LitExpr is a literal: integer, float, string, char, bool, null.
	LitExpr struct{ Val string }

	// IdentExpr is a bare identifier.
	IdentExpr struct{ Name string }

	// SelectorExpr is `X.sel` — field access or method ref.
	SelectorExpr struct {
		X   Expr
		Sel string
	}

	// IndexExpr is `X[idx]`.
	IndexExpr struct {
		X   Expr
		Idx Expr
	}

	// CallExpr is `Fun(Args...)`.
	CallExpr struct {
		Fun  Expr
		Args []Expr
	}

	// NewExpr is `new Type(Args...)`.
	NewExpr struct {
		Type string
		Args []Expr
	}

	// NewArrayExpr is `new Type[size]` or `new Type[]{elems...}`.
	NewArrayExpr struct {
		ElemType string
		Size     Expr   // nil for initialiser form
		Elems    []Expr // non-nil for `{...}` initialiser form
	}

	// ArrayLitExpr is `{a, b, c}` (array initialiser without `new`).
	ArrayLitExpr struct{ Elems []Expr }

	// UnaryExpr is `op X`.
	UnaryExpr struct {
		Op string
		X  Expr
	}

	// BinaryExpr is `X op Y`.
	BinaryExpr struct {
		X  Expr
		Op string
		Y  Expr
	}

	// TernaryExpr is `Cond ? Then : Else`.
	TernaryExpr struct {
		Cond Expr
		Then Expr
		Else Expr
	}

	// CastExpr is `(Type) X`.
	CastExpr struct {
		Type string
		X    Expr
	}

	// TypeCheckExpr is `X instanceof Type`.
	TypeCheckExpr struct {
		X    Expr
		Type string
	}

	// LambdaExpr is a Java lambda `param -> body` or `(p1, p2) -> body`.
	// Translated to an inline Go func literal.
	LambdaExpr struct {
		Params string // raw parameter list text, e.g. "n" or "n1, n2"
		Body   Expr   // the expression body (right of ->)
	}

	// RawExpr holds an untranslatable Java expression snippet.
	RawExpr struct{ Text string }
)

func (*LitExpr) exprNode()       {}
func (*IdentExpr) exprNode()     {}
func (*SelectorExpr) exprNode()  {}
func (*IndexExpr) exprNode()     {}
func (*CallExpr) exprNode()      {}
func (*NewExpr) exprNode()       {}
func (*NewArrayExpr) exprNode()  {}
func (*ArrayLitExpr) exprNode()  {}
func (*UnaryExpr) exprNode()     {}
func (*BinaryExpr) exprNode()    {}
func (*TernaryExpr) exprNode()   {}
func (*CastExpr) exprNode()      {}
func (*TypeCheckExpr) exprNode() {}
func (*LambdaExpr) exprNode()    {}
func (*RawExpr) exprNode()       {}
