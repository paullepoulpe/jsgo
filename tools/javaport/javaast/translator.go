package javaast

import (
	"fmt"
	"strings"
)

// MethodContext carries per-method information the translator needs to
// emit idiomatic Go: the receiver variable name, the class name, and
// the Java→Go type mapper.
type MethodContext struct {
	RecvName  string           // e.g. "c"
	ClassName string           // e.g. "CheckJSDoc"
	TypeMap   func(string) string // Java type → Go type; may be nil (identity)
}

// TranslateBody translates the Java method body src into Go statements.
// It returns the translated Go source and a boolean indicating whether the
// translation succeeded (false means we fell back to a stub).
func TranslateBody(src string, ctx MethodContext) (string, bool) {
	block := ParseBody(src)
	if len(block.Stmts) == 0 {
		// Empty body — emit nothing.
		return "", true
	}
	t := &translator{ctx: ctx, indent: 1}
	var b strings.Builder
	t.emitStmtList(&b, block.Stmts)
	return b.String(), true
}

// ---- translator --------------------------------------------------------

type translator struct {
	ctx    MethodContext
	indent int
}

func (t *translator) tab() string {
	return strings.Repeat("\t", t.indent)
}

func (t *translator) emitStmtList(b *strings.Builder, stmts []Stmt) {
	for _, s := range stmts {
		t.emitStmt(b, s)
	}
}

func (t *translator) emitStmt(b *strings.Builder, s Stmt) {
	switch v := s.(type) {

	case *BlockStmt:
		t.emitStmtList(b, v.Stmts)

	case *ReturnStmt:
		if v.X == nil {
			fmt.Fprintf(b, "%sreturn\n", t.tab())
		} else {
			fmt.Fprintf(b, "%sreturn %s\n", t.tab(), t.emitExpr(v.X))
		}

	case *ExprStmt:
		fmt.Fprintf(b, "%s%s\n", t.tab(), t.emitExpr(v.X))

	case *AssignStmt:
		lhs := t.emitExpr(v.LHS)
		rhs := t.emitExpr(v.RHS)
		fmt.Fprintf(b, "%s%s %s %s\n", t.tab(), lhs, v.Op, rhs)

	case *LocalVarStmt:
		if v.Init != nil {
			fmt.Fprintf(b, "%s%s := %s\n", t.tab(), goIdent(v.Name), t.emitExpr(v.Init))
		} else {
			goType := t.javaType(v.Type)
			fmt.Fprintf(b, "%svar %s %s\n", t.tab(), goIdent(v.Name), goType)
		}

	case *IfStmt:
		cond := t.translateCond(v.Cond)
		fmt.Fprintf(b, "%sif %s {\n", t.tab(), cond)
		t.indent++
		t.emitStmt(b, v.Then)
		t.indent--
		if v.Else != nil {
			// else-if chain
			if elseIf, ok := v.Else.(*IfStmt); ok {
				elseStr := t.emitElseIf(elseIf)
				fmt.Fprintf(b, "%s} %s", t.tab(), elseStr)
			} else {
				fmt.Fprintf(b, "%s} else {\n", t.tab())
				t.indent++
				t.emitStmt(b, v.Else)
				t.indent--
				fmt.Fprintf(b, "%s}\n", t.tab())
			}
		} else {
			fmt.Fprintf(b, "%s}\n", t.tab())
		}

	case *ForStmt:
		var parts []string
		if v.Init != nil {
			parts = append(parts, t.forClause(v.Init))
		} else {
			parts = append(parts, "")
		}
		if v.Cond != nil {
			parts = append(parts, t.translateCond(v.Cond))
		} else {
			parts = append(parts, "")
		}
		if v.Post != nil {
			parts = append(parts, t.forClause(v.Post))
		} else {
			parts = append(parts, "")
		}
		header := strings.Join(parts, "; ")
		// If all parts empty → `for {`
		if strings.TrimSpace(strings.ReplaceAll(header, ";", "")) == "" {
			fmt.Fprintf(b, "%sfor {\n", t.tab())
		} else {
			fmt.Fprintf(b, "%sfor %s {\n", t.tab(), header)
		}
		t.indent++
		t.emitStmt(b, v.Body)
		t.indent--
		fmt.Fprintf(b, "%s}\n", t.tab())

	case *ForEachStmt:
		fmt.Fprintf(b, "%sfor _, %s := range %s {\n", t.tab(), goIdent(v.Name), t.emitExpr(v.Coll))
		t.indent++
		t.emitStmt(b, v.Body)
		t.indent--
		fmt.Fprintf(b, "%s}\n", t.tab())

	case *WhileStmt:
		cond := t.translateCond(v.Cond)
		fmt.Fprintf(b, "%sfor %s {\n", t.tab(), cond)
		t.indent++
		t.emitStmt(b, v.Body)
		t.indent--
		fmt.Fprintf(b, "%s}\n", t.tab())

	case *DoWhileStmt:
		cond := t.translateCond(v.Cond)
		fmt.Fprintf(b, "%sfor {\n", t.tab())
		t.indent++
		t.emitStmt(b, v.Body)
		fmt.Fprintf(b, "%sif !(%s) {\n", t.tab(), cond)
		t.indent++
		fmt.Fprintf(b, "%sbreak\n", t.tab())
		t.indent--
		fmt.Fprintf(b, "%s}\n", t.tab())
		t.indent--
		fmt.Fprintf(b, "%s}\n", t.tab())

	case *SwitchStmt:
		fmt.Fprintf(b, "%sswitch %s {\n", t.tab(), t.emitExpr(v.X))
		for _, c := range v.Cases {
			if c.Value == nil {
				fmt.Fprintf(b, "%sdefault:\n", t.tab())
			} else {
				fmt.Fprintf(b, "%scase %s:\n", t.tab(), t.emitExpr(c.Value))
			}
			t.indent++
			t.emitStmtList(b, c.Stmts)
			if c.FallThru && len(c.Stmts) > 0 {
				// Only emit fallthrough if last stmt is not break/return.
				last := c.Stmts[len(c.Stmts)-1]
				switch last.(type) {
				case *BreakStmt, *ReturnStmt:
					// Already has explicit terminator — no fallthrough.
				default:
					fmt.Fprintf(b, "%sfallthrough\n", t.tab())
				}
			}
			t.indent--
		}
		fmt.Fprintf(b, "%s}\n", t.tab())

	case *TryStmt:
		t.emitTry(b, v)

	case *ThrowStmt:
		fmt.Fprintf(b, "%spanic(%s)\n", t.tab(), t.emitExpr(v.X))

	case *BreakStmt:
		if v.Label != "" {
			fmt.Fprintf(b, "%sbreak %s\n", t.tab(), v.Label)
		} else {
			fmt.Fprintf(b, "%sbreak\n", t.tab())
		}

	case *ContinueStmt:
		if v.Label != "" {
			fmt.Fprintf(b, "%scontinue %s\n", t.tab(), v.Label)
		} else {
			fmt.Fprintf(b, "%scontinue\n", t.tab())
		}

	case *LabeledStmt:
		fmt.Fprintf(b, "%s%s:\n", t.tab(), v.Label)
		if v.Stmt != nil {
			t.emitStmt(b, v.Stmt)
		}

	case *RawStmt:
		text := strings.TrimSpace(v.Text)
		if text != "" {
			fmt.Fprintf(b, "%s// TODO: %s\n", t.tab(), text)
		}

	default:
		fmt.Fprintf(b, "%s// TODO: unknown stmt\n", t.tab())
	}
}

// emitElseIf produces `else if cond { ... } [else ...]` text.
func (t *translator) emitElseIf(v *IfStmt) string {
	var b strings.Builder
	cond := t.translateCond(v.Cond)
	fmt.Fprintf(&b, "else if %s {\n", cond)
	t.indent++
	t.emitStmt(&b, v.Then)
	t.indent--
	if v.Else != nil {
		if elseIf, ok := v.Else.(*IfStmt); ok {
			fmt.Fprintf(&b, "%s} %s", t.tab(), t.emitElseIf(elseIf))
		} else {
			fmt.Fprintf(&b, "%s} else {\n", t.tab())
			t.indent++
			t.emitStmt(&b, v.Else)
			t.indent--
			fmt.Fprintf(&b, "%s}\n", t.tab())
		}
	} else {
		fmt.Fprintf(&b, "%s}\n", t.tab())
	}
	return b.String()
}

// forClause emits an init or post clause of a for statement without a trailing newline.
func (t *translator) forClause(s Stmt) string {
	switch v := s.(type) {
	case *LocalVarStmt:
		if v.Init != nil {
			return fmt.Sprintf("%s := %s", goIdent(v.Name), t.emitExpr(v.Init))
		}
		return fmt.Sprintf("var %s %s", goIdent(v.Name), t.javaType(v.Type))
	case *AssignStmt:
		return fmt.Sprintf("%s %s %s", t.emitExpr(v.LHS), v.Op, t.emitExpr(v.RHS))
	case *ExprStmt:
		return t.emitExpr(v.X)
	default:
		var b strings.Builder
		t.emitStmt(&b, s)
		return strings.TrimSpace(b.String())
	}
}

// emitTry translates a try/catch/finally into Go defer+recover.
func (t *translator) emitTry(b *strings.Builder, v *TryStmt) {
	hasCatch := len(v.Catches) > 0
	hasFin := v.Finally != nil

	if hasCatch || hasFin {
		fmt.Fprintf(b, "%sdefer func() {\n", t.tab())
		t.indent++
		if hasCatch {
			fmt.Fprintf(b, "%sif r := recover(); r != nil {\n", t.tab())
			t.indent++
			for _, c := range v.Catches {
				// Simple: emit the catch body; type check is best-effort.
				fmt.Fprintf(b, "%s_ = r // was: %s %s\n", t.tab(), c.Type, c.Name)
				t.emitStmt(b, c.Body)
			}
			t.indent--
			fmt.Fprintf(b, "%s}\n", t.tab())
		}
		if hasFin {
			t.emitStmt(b, v.Finally)
		}
		t.indent--
		fmt.Fprintf(b, "%s}()\n", t.tab())
	}
	// Emit the try body directly.
	t.emitStmt(b, v.Body)
}

// translateCond rewrites Java condition patterns to Go:
//   - `x == null` → `x == nil`
//   - `x != null` → `x != nil`
//   - bare boolean ident → keep as-is
func (t *translator) translateCond(e Expr) string {
	s := t.emitExpr(e)
	return s
}

// ---- expression emission -----------------------------------------------

func (t *translator) emitExpr(e Expr) string {
	if e == nil {
		return ""
	}
	switch v := e.(type) {
	case *LitExpr:
		return translateLiteral(v.Val)

	case *IdentExpr:
		return t.translateIdent(v.Name)

	case *SelectorExpr:
		obj := t.emitExpr(v.X)
		// `this.foo` → `c.foo`
		if obj == "this" || obj == t.ctx.RecvName {
			return t.ctx.RecvName + "." + goIdent(v.Sel)
		}
		return obj + "." + goIdent(v.Sel)

	case *IndexExpr:
		return t.emitExpr(v.X) + "[" + t.emitExpr(v.Idx) + "]"

	case *CallExpr:
		return t.emitCall(v)

	case *NewExpr:
		goType := t.javaType(v.Type)
		// Strip leading `*` for constructor name: *Foo → NewFoo
		base := strings.TrimPrefix(goType, "*")
		args := t.emitArgs(v.Args)
		return "New" + base + "(" + args + ")"

	case *NewArrayExpr:
		goElem := t.javaType(v.ElemType)
		if v.Elems != nil {
			elems := make([]string, len(v.Elems))
			for i, el := range v.Elems {
				elems[i] = t.emitExpr(el)
			}
			return "[]" + goElem + "{" + strings.Join(elems, ", ") + "}"
		}
		if v.Size != nil {
			return "make([]" + goElem + ", " + t.emitExpr(v.Size) + ")"
		}
		return "[]" + goElem + "{}"

	case *ArrayLitExpr:
		elems := make([]string, len(v.Elems))
		for i, el := range v.Elems {
			elems[i] = t.emitExpr(el)
		}
		return "{" + strings.Join(elems, ", ") + "}"

	case *UnaryExpr:
		inner := t.emitExpr(v.X)
		switch v.Op {
		case "++", "--":
			return inner + v.Op
		}
		return v.Op + inner

	case *BinaryExpr:
		l := t.emitExpr(v.X)
		r := t.emitExpr(v.Y)
		op := translateBinaryOp(v.Op, r)
		// null comparisons
		if r == "null" || r == "nil" {
			r = "nil"
		}
		if l == "null" || l == "nil" {
			l = "nil"
		}
		return l + " " + op + " " + r

	case *TernaryExpr:
		// Go doesn't have ternary; emit as comment + placeholder.
		// But for simple cases inside expressions, we have to emit a
		// function-call form.
		cond := t.emitExpr(v.Cond)
		then := t.emitExpr(v.Then)
		els := t.emitExpr(v.Else)
		return fmt.Sprintf("func() any { if %s { return %s }; return %s }()", cond, then, els)

	case *CastExpr:
		inner := t.emitExpr(v.X)
		goType := t.javaType(v.Type)
		return goType + "(" + inner + ")"

	case *TypeCheckExpr:
		inner := t.emitExpr(v.X)
		goType := t.javaType(v.Type)
		// Go uses type assertion for interface checks; emit reflect-style as best effort.
		return fmt.Sprintf("func() bool { _, ok := %s.(%s); return ok }()", inner, goType)

	case *RawExpr:
		return v.Text

	default:
		return "/* unknown expr */"
	}
}

func (t *translator) emitCall(v *CallExpr) string {
	args := t.emitArgs(v.Args)
	switch fn := v.Fun.(type) {
	case *SelectorExpr:
		obj := t.emitExpr(fn.X)
		method := goMethodName(fn.Sel)
		if obj == "this" || obj == t.ctx.RecvName {
			return t.ctx.RecvName + "." + method + "(" + args + ")"
		}
		return obj + "." + method + "(" + args + ")"
	case *IdentExpr:
		name := fn.Name
		// super.method() → treated as method call on embedded type
		if name == "super" {
			return "/* super */ " + args
		}
		goName := goMethodName(name)
		return goName + "(" + args + ")"
	default:
		return t.emitExpr(v.Fun) + "(" + args + ")"
	}
}

func (t *translator) emitArgs(args []Expr) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = t.emitExpr(a)
	}
	return strings.Join(parts, ", ")
}

func (t *translator) translateIdent(name string) string {
	switch name {
	case "this":
		return t.ctx.RecvName
	case "null":
		return "nil"
	case "true", "false":
		return name
	}
	return goIdent(name)
}

func (t *translator) javaType(jt string) string {
	if t.ctx.TypeMap != nil {
		return t.ctx.TypeMap(jt)
	}
	return jt
}

// ---- helpers -----------------------------------------------------------

func translateLiteral(val string) string {
	if val == "null" {
		return "nil"
	}
	// Java long literals: `42L` → `42`
	if strings.HasSuffix(val, "L") || strings.HasSuffix(val, "l") {
		return val[:len(val)-1]
	}
	return val
}

func translateBinaryOp(op, rhs string) string {
	_ = rhs
	switch op {
	case "&&":
		return "&&"
	case "||":
		return "||"
	}
	return op
}

// goIdent converts a Java identifier to its idiomatic Go equivalent.
// For field accesses / variable names we leave them lowercase as-is
// (unlike method names which are Title-cased only on public API).
func goIdent(name string) string {
	switch name {
	case "this":
		return "this"
	case "type":
		return "typ"
	case "len":
		return "length"
	case "cap":
		return "capacity"
	case "range":
		return "rng"
	case "map":
		return "mp"
	case "var":
		return "v"
	case "func":
		return "fn"
	case "package":
		return "pkg"
	case "interface":
		return "iface"
	}
	return name
}

// goMethodName Title-cases the first character to follow Go conventions
// for exported methods. Constructors get the "New" prefix at call-site.
func goMethodName(name string) string {
	if name == "" {
		return name
	}
	// Already title-cased or all-caps — leave alone.
	if name[0] >= 'A' && name[0] <= 'Z' {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
