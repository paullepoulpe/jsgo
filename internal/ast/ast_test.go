package ast

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func newStore(t *testing.T) {
	prev := rhino.SetCurrentNodeStore(rhino.NewNodeStore())
	t.Cleanup(func() { rhino.SetCurrentNodeStore(prev) })
}

func TestIfAccessors(t *testing.T) {
	newStore(t)
	cond := NewName("x")
	then := NewBlock(NewReturn().AsStatement())
	ifs := NewIf(cond.AsExpr(), then)

	if c, ok := ifs.Cond().AsName(); !ok || c.Ident() != "x" {
		t.Fatalf("Cond() = %v, want the name x", ifs.Cond().Node())
	}
	if ifs.Then() != then {
		t.Fatal("Then() is not the block passed to NewIf")
	}
	if _, ok := ifs.Else(); ok {
		t.Fatal("Else() reports an else on an if without one")
	}
	els := NewBlock()
	ifs.SetElse(els)
	if got, ok := ifs.Else(); !ok || got != els {
		t.Fatal("Else() does not return the block given to SetElse")
	}
	ifs.RemoveElse()
	if _, ok := ifs.Else(); ok {
		t.Fatal("RemoveElse() left the else")
	}

	// The categories: an If is a Statement, and back.
	s := ifs.AsStatement()
	if back, ok := s.AsIf(); !ok || back != ifs {
		t.Fatal("Statement.AsIf() does not give the If back")
	}
	if _, ok := s.AsBlock(); ok {
		t.Fatal("an If converted to a Block")
	}
	if _, ok := AsExpr(ifs.Node()); ok {
		t.Fatal("an If converted to an Expr")
	}
	if _, ok := AsIf(rhino.NilNode); ok {
		t.Fatal("the nil node converted to an If")
	}
}

func TestListsAndNavigation(t *testing.T) {
	newStore(t)
	a, b, c := NewEmpty(), NewDebugger(), NewEmpty()
	block := NewBlock(a.AsStatement(), b.AsStatement())
	block.AppendStatement(c.AsStatement())

	var got []rhino.Token
	for s := range block.Statements() {
		got = append(got, s.Node().GetToken())
	}
	want := []rhino.Token{rhino.Token_EMPTY, rhino.Token_DEBUGGER, rhino.Token_EMPTY}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("Statements() = %v, want %v", got, want)
	}
	if first, _ := block.FirstStatement(); first != a.AsStatement() {
		t.Fatal("FirstStatement() is not the first statement")
	}
	if last, _ := block.LastStatement(); last != c.AsStatement() {
		t.Fatal("LastStatement() is not the last statement")
	}
	if next, ok := a.AsStatement().Next(); !ok || next != b.AsStatement() {
		t.Fatal("Next() is not the following statement")
	}
	if _, ok := a.AsStatement().Previous(); ok {
		t.Fatal("Previous() of the first statement reports one")
	}
	if _, ok := NewBlock().LastStatement(); ok {
		t.Fatal("LastStatement() of an empty block reports one")
	}

	// A call's arguments start after its callee.
	call := NewCall(NewName("f").AsExpr(), NewName("x").AsArg())
	if last, ok := call.LastArg(); !ok || last.Node().GetString() != "x" {
		t.Fatal("LastArg() is not the argument")
	}
	if _, ok := NewCall(NewName("f").AsExpr()).LastArg(); ok {
		t.Fatal("LastArg() of a call without arguments reports the callee")
	}
	if n := call.NumArgs(); n != 1 {
		t.Fatalf("NumArgs() = %d, want 1 (the callee is not an argument)", n)
	}
	if n := NewCall(NewName("f").AsExpr()).NumArgs(); n != 0 {
		t.Fatalf("NumArgs() of f() = %d, want 0", n)
	}
	if n := block.NumStatements(); n != 3 {
		t.Fatalf("NumStatements() = %d, want 3", n)
	}
}

// A category has the accessors that all its kinds share: the object of a
// GetProp (its only child) and of a GetElem (its first one), the args of a
// Call and of an OptChainCall.
func TestCategorySharedAccessors(t *testing.T) {
	newStore(t)
	for _, n := range []rhino.Node{
		NewGetProp(NewName("o").AsExpr(), "p").Node(),
		NewGetElem(NewName("o").AsExpr(), NewStringLit("p").AsExpr()).Node(),
		NewOptChainGetElem(NewName("o").AsExpr(), NewStringLit("p").AsExpr()).Node(),
	} {
		get, ok := AsNormalOrOptChainGet(n)
		if !ok {
			t.Fatalf("AsNormalOrOptChainGet(%s) failed", n.GetToken())
		}
		if obj, ok := get.Object().AsName(); !ok || obj.Ident() != "o" {
			t.Fatalf("Object() of a %s is not its object", n.GetToken())
		}
	}

	for _, n := range []rhino.Node{
		NewCall(NewName("f").AsExpr(), NewName("x").AsArg(), NewName("y").AsArg()).Node(),
		NewOptChainCall(NewName("f").AsExpr(), NewName("x").AsArg(), NewName("y").AsArg()).Node(),
	} {
		call, ok := AsNormalOrOptChainCall(n)
		if !ok {
			t.Fatalf("AsNormalOrOptChainCall(%s) failed", n.GetToken())
		}
		var args []string
		for a := range call.Args() {
			args = append(args, a.Node().GetString())
		}
		if len(args) != 2 || args[0] != "x" || args[1] != "y" {
			t.Fatalf("Args() of a %s = %q, want [x y]", n.GetToken(), args)
		}
		first, _ := call.FirstArg()
		if second, ok := first.Next(); !ok || second.Node().GetString() != "y" {
			t.Fatalf("FirstArg().Next() of a %s is not the second argument", n.GetToken())
		}
	}
	if call, _ := AsNormalOrOptChainCall(NewCall(NewName("f").AsExpr()).Node()); call.HasArgs() {
		t.Fatal("HasArgs() of a call without arguments counts the callee")
	}
}

// NormalOrOptChainGetElem is both element accesses, with the object and the
// index they share; a property access is not one.
func TestNormalOrOptChainGetElem(t *testing.T) {
	newStore(t)
	for _, n := range []rhino.Node{
		NewGetElem(NewName("o").AsExpr(), NewStringLit("p").AsExpr()).Node(),
		NewOptChainGetElem(NewName("o").AsExpr(), NewStringLit("p").AsExpr()).Node(),
	} {
		get, ok := AsNormalOrOptChainGetElem(n)
		if !ok {
			t.Fatalf("AsNormalOrOptChainGetElem(%s) failed", n.GetToken())
		}
		if obj, ok := get.Object().AsName(); !ok || obj.Ident() != "o" {
			t.Fatalf("Object() of a %s is not its object", n.GetToken())
		}
		if index, ok := get.Index().AsStringLit(); !ok || index.Value() != "p" {
			t.Fatalf("Index() of a %s is not its index", n.GetToken())
		}
		if !get.AsNormalOrOptChainGet().Object().Is(get.Object().Node()) {
			t.Fatalf("AsNormalOrOptChainGet() of a %s has another object", n.GetToken())
		}
		if _, ok := get.AsExpr().AsNormalOrOptChainGetElem(); !ok {
			t.Fatalf("Expr.AsNormalOrOptChainGetElem() of a %s failed", n.GetToken())
		}
	}
	if _, ok := AsNormalOrOptChainGetElem(NewGetProp(NewName("o").AsExpr(), "p").Node()); ok {
		t.Fatal("o.p converted to a NormalOrOptChainGetElem")
	}
	if _, ok := AsNormalOrOptChainGetElem(rhino.NilNode); ok {
		t.Fatal("the nil node converted to a NormalOrOptChainGetElem")
	}
}

func TestOnlyStatementAndDeclarator(t *testing.T) {
	newStore(t)
	if _, ok := NewBlock().OnlyStatement(); ok {
		t.Fatal("OnlyStatement() of an empty block reports one")
	}
	ret := NewReturn().AsStatement()
	block := NewBlock(ret)
	if s, ok := block.OnlyStatement(); !ok || s != ret {
		t.Fatal("OnlyStatement() is not the statement of a block with one")
	}
	block.AppendStatement(NewEmpty().AsStatement())
	if _, ok := block.OnlyStatement(); ok {
		t.Fatal("OnlyStatement() of a block with two statements reports one")
	}

	x := NewName("x").AsDeclarator()
	decl := NewNameDeclaration(rhino.Token_VAR, x)
	if d, ok := decl.OnlyDeclarator(); !ok || d != x {
		t.Fatal("OnlyDeclarator() is not the declarator of var x")
	}
	decl.AppendDeclarator(NewName("y").AsDeclarator())
	if _, ok := decl.OnlyDeclarator(); ok {
		t.Fatal("OnlyDeclarator() of var x, y reports one")
	}
}

// Assignment is both kinds of assignment, with their shared target and
// value.
func TestAssignmentCategory(t *testing.T) {
	newStore(t)
	for _, e := range []Expr{
		NewAssign(NewName("a").AsTarget(), NewName("b").AsExpr()).AsExpr(),
		NewAssignOp(rhino.Token_ASSIGN_ADD, NewName("a").AsTarget(), NewName("b").AsExpr()).AsExpr(),
	} {
		n := e.Node()
		assign, ok := e.AsAssignment()
		if !ok {
			t.Fatalf("Expr.AsAssignment() of a %s failed", n.GetToken())
		}
		if assign.Token() != n.GetToken() {
			t.Fatalf("Token() = %s, want %s", assign.Token(), n.GetToken())
		}
		if target, ok := assign.Target().AsName(); !ok || target.Ident() != "a" {
			t.Fatalf("Target() of a %s is not its target", n.GetToken())
		}
		if value, ok := assign.Value().AsName(); !ok || value.Ident() != "b" {
			t.Fatalf("Value() of a %s is not its value", n.GetToken())
		}
	}
	bin := NewBinary(rhino.Token_EQ, NewName("a").AsExpr(), NewName("b").AsExpr())
	if _, ok := bin.AsExpr().AsAssignment(); ok {
		t.Fatal("a == b converted to an Assignment")
	}
}

// CallOrNew is the calls and new, with the callee and arguments they share.
func TestCallOrNew(t *testing.T) {
	newStore(t)
	callee := func() Expr { return NewName("f").AsExpr() }
	arg := func() Arg { return NewName("a").AsArg() }
	for _, e := range []Expr{
		NewCall(callee(), arg()).AsExpr(),
		NewOptChainCall(callee(), arg()).AsExpr(),
		NewNew(callee(), arg()).AsExpr(),
	} {
		n := e.Node()
		call, ok := e.AsCallOrNew()
		if !ok {
			t.Fatalf("Expr.AsCallOrNew() of a %s failed", n.GetToken())
		}
		if f, ok := call.Callee().AsName(); !ok || f.Ident() != "f" {
			t.Fatalf("Callee() of a %s is not its callee", n.GetToken())
		}
		if a, ok := call.OnlyArg(); !ok || !a.Is(n.GetLastChild()) {
			t.Fatalf("OnlyArg() of a %s is not its argument", n.GetToken())
		}
	}
	if call := NewCall(callee()).AsNormalOrOptChainCall().AsCallOrNew(); call.HasArgs() {
		t.Fatal("HasArgs() is true for f()")
	}
	tagged := NewTaggedTemplateLit(callee(), NewTemplateLit())
	if _, ok := tagged.AsExpr().AsCallOrNew(); ok {
		t.Fatal("a tagged template converted to a CallOrNew")
	}
}

// EnhancedFor is the for-in and for-of loops, with the target they share.
func TestEnhancedFor(t *testing.T) {
	newStore(t)
	target := func() ForTarget {
		return NewNameDeclaration(rhino.Token_VAR, NewName("x").AsDeclarator()).AsForTarget()
	}
	for _, n := range []rhino.Node{
		NewForIn(target(), NewName("o").AsExpr(), NewBlock()).Node(),
		NewForOf(target(), NewName("o").AsExpr(), NewBlock()).Node(),
		NewForAwaitOf(target(), NewName("o").AsExpr(), NewBlock()).Node(),
	} {
		loop, ok := AsEnhancedFor(n)
		if !ok {
			t.Fatalf("AsEnhancedFor(%s) failed", n.GetToken())
		}
		if decl, ok := loop.Target().AsNameDeclaration(); !ok || decl.Token() != rhino.Token_VAR {
			t.Fatalf("Target() of a %s is not its var", n.GetToken())
		}
		if back, ok := loop.AsLoop().AsEnhancedFor(); !ok || back != loop {
			t.Fatalf("Loop.AsEnhancedFor() does not give the %s back", n.GetToken())
		}
		if o, ok := loop.Collection().AsName(); !ok || o.Ident() != "o" {
			t.Fatalf("Collection() of a %s is not the name o", n.GetToken())
		}
	}
	if _, ok := AsEnhancedFor(NewWhile(NewName("c").AsExpr(), NewBlock()).Node()); ok {
		t.Fatal("a while converted to an EnhancedFor")
	}
}

// ForOfLoop is the for-of and for-await-of loops, with the target and
// iterable they share.
func TestForOfLoop(t *testing.T) {
	newStore(t)
	for _, n := range []rhino.Node{
		NewForOf(NewName("x").AsForTarget(), NewName("o").AsExpr(), NewBlock()).Node(),
		NewForAwaitOf(NewName("x").AsForTarget(), NewName("o").AsExpr(), NewBlock()).Node(),
	} {
		loop, ok := AsForOfLoop(n)
		if !ok {
			t.Fatalf("AsForOfLoop(%s) failed", n.GetToken())
		}
		if x, ok := loop.Target().AsName(); !ok || x.Ident() != "x" {
			t.Fatalf("Target() of a %s is not its target", n.GetToken())
		}
		if o, ok := loop.Iterable().AsName(); !ok || o.Ident() != "o" {
			t.Fatalf("Iterable() of a %s is not its iterable", n.GetToken())
		}
		if !loop.Body().Is(n.GetLastChild()) {
			t.Fatalf("Body() of a %s is not its body", n.GetToken())
		}
		if back, ok := loop.AsEnhancedFor().AsForOfLoop(); !ok || back != loop {
			t.Fatalf("EnhancedFor.AsForOfLoop() does not give the %s back", n.GetToken())
		}
	}
	forIn := NewForIn(NewName("x").AsForTarget(), NewName("o").AsExpr(), NewBlock())
	if _, ok := AsForOfLoop(forIn.Node()); ok {
		t.Fatal("a for-in converted to a ForOfLoop")
	}
}

func TestFunctionRemoveName(t *testing.T) {
	newStore(t)
	fn := NewFunction("f", NewParamList(), NewBlock().AsFunctionBody())
	if name, ok := fn.Name(); !ok || name.Ident() != "f" {
		t.Fatal("Name() is not the name given to NewFunction")
	}
	fn.RemoveName()
	if _, ok := fn.Name(); ok {
		t.Fatal("RemoveName() left the name")
	}
	if !fn.Node().GetFirstChild().IsName() {
		t.Fatal("RemoveName() removed the NAME child")
	}
}

func TestTryCatch(t *testing.T) {
	newStore(t)
	catch := NewCatchClause(t)
	try := NewTry(NewBlock())
	if _, ok := try.Catch(); ok {
		t.Fatal("Catch() reports a catch on a try without one")
	}
	holder := try.Catches()
	if !holder.Is(try.Node().GetSecondChild()) || holder.Node().HasChildren() {
		t.Fatal("Catches() of a try without catch is not its empty second child")
	}
	try.SetCatch(catch)
	if got, ok := try.Catch(); !ok || got != catch {
		t.Fatal("Catch() does not return the catch given to SetCatch")
	}
	if got, ok := try.Catches().Catch(); !try.Catches().Is(holder.Node()) || !ok || got != catch {
		t.Fatal("Catches() is not the block that holds the catch")
	}
	var problems []string
	Check(try.Node(), func(n rhino.Node, p string) { problems = append(problems, p) })
	if len(problems) != 1 || !strings.Contains(problems[0], "Catch has no param") {
		t.Fatalf("Check() = %q, want only the missing catch param", problems)
	}
}

// NewCatchClause returns a catch clause without its param (Check reports it).
func NewCatchClause(t *testing.T) Catch {
	c, ok := AsCatch(rhino.NewNode(rhino.Token_CATCH))
	if !ok {
		t.Fatal("AsCatch failed on a CATCH node")
	}
	return c
}

func TestCheckReportsProblems(t *testing.T) {
	newStore(t)
	// An if whose then is an expression instead of a block, and an extra child.
	ifs := rhino.NewNode(rhino.Token_IF)
	ifs.AddChildToBack(rhino.IR_Name("x"))
	ifs.AddChildToBack(rhino.IR_Name("y"))
	var problems []string
	Check(ifs, func(n rhino.Node, p string) { problems = append(problems, p) })
	if len(problems) != 1 || problems[0] != "If.then is a NAME, want Block" {
		t.Fatalf("Check() = %q", problems)
	}
	ifs.AddChildToBack(rhino.IR_Block())
	ifs.AddChildToBack(rhino.IR_Block())
	problems = nil
	Check(ifs, func(n rhino.Node, p string) { problems = append(problems, p) })
	if len(problems) != 2 || problems[1] != "If has an extra child BLOCK" {
		t.Fatalf("Check() = %q", problems)
	}
}

// nodes.gen.go must be what the generator makes from grammar.txt.
func TestGeneratedCodeIsUpToDate(t *testing.T) {
	if testing.Short() {
		t.Skip("runs the generator")
	}
	out := filepath.Join(t.TempDir(), "nodes.gen.go")
	cmd := exec.Command("go", "run", "../../tools/astgen", "-grammar", "grammar.txt", "-out", out)
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("astgen: %v\n%s", err, b)
	}
	want, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("nodes.gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("nodes.gen.go is out of date: run go generate ./internal/ast")
	}
}

// NormalGet is a GetProp or a GetElem (what an assignment to a property
// assigns), not a property access in an optional chain.
func TestNormalGet(t *testing.T) {
	newStore(t)
	for _, n := range []rhino.Node{
		NewGetProp(NewName("o").AsExpr(), "p").Node(),
		NewGetElem(NewName("o").AsExpr(), NewStringLit("p").AsExpr()).Node(),
	} {
		get, ok := AsNormalGet(n)
		if !ok {
			t.Fatalf("AsNormalGet(%s) failed", n.GetToken())
		}
		if obj, ok := get.Object().AsName(); !ok || obj.Ident() != "o" {
			t.Fatalf("Object() of a %s is not its object", n.GetToken())
		}
		if back, ok := get.AsTarget().AsNormalGet(); !ok || back != get {
			t.Fatalf("Target.AsNormalGet() does not give the %s back", n.GetToken())
		}
	}
	for _, n := range []rhino.Node{
		NewOptChainGetElem(NewName("o").AsExpr(), NewStringLit("p").AsExpr()).Node(),
		NewName("o").Node(),
	} {
		if _, ok := AsNormalGet(n); ok {
			t.Fatalf("a %s converted to a NormalGet", n.GetToken())
		}
	}
}

// Must conversions return the view for the right kind and panic, naming
// both kinds, for another.
func TestMust(t *testing.T) {
	newStore(t)
	name := NewName("x")
	if MustName(name.Node()) != name || MustExpr(name.Node()).Node() != name.Node() {
		t.Fatal("Must conversion of the right kind")
	}
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok || e.Error() != "ast.MustIf: a NAME node" {
			t.Fatalf("panic %v, want ast.MustIf: a NAME node", r)
		}
	}()
	MustIf(name.Node())
}

func TestTemplateLitString(t *testing.T) {
	newStore(t)
	cooked := "a\nb"
	s, ok := AsTemplateLitString(rhino.Node_NewTemplateLitString(&cooked, `a\nb`))
	if !ok {
		t.Fatal("AsTemplateLitString")
	}
	if c, ok := s.Cooked(); s.Raw() != `a\nb` || !ok || c != "a\nb" {
		t.Fatalf("Raw %q Cooked %q %v", s.Raw(), c, ok)
	}
	none, _ := AsTemplateLitString(rhino.Node_NewTemplateLitString(nil, `\u{`))
	if _, ok := none.Cooked(); ok {
		t.Fatal("Cooked of a string without one")
	}
}

func TestAsKindOf(t *testing.T) {
	newStore(t)
	not := NewUnary(rhino.Token_NOT, NewName("a").AsExpr()).AsExpr()
	if _, ok := not.AsUnaryOf(rhino.Token_NOT); !ok {
		t.Fatal("AsUnaryOf(NOT) of !a")
	}
	if _, ok := not.AsUnaryOf(rhino.Token_VOID); ok {
		t.Fatal("AsUnaryOf(VOID) of !a")
	}
	if _, ok := AsBinaryOf(not.Node(), rhino.Token_NOT); ok {
		t.Fatal("AsBinaryOf(NOT) of a unary")
	}
}

// ReadAhead lets the loop body move the current element: here every
// statement moves to another block, and all three are still visited.
func TestReadAhead(t *testing.T) {
	newStore(t)
	a, b, c := NewEmpty(), NewDebugger(), NewEmpty()
	from := NewBlock(a.AsStatement(), b.AsStatement(), c.AsStatement())
	to := NewBlock()
	n := 0
	for s := range ReadAhead(from.FirstStatement()) {
		to.AppendStatement(s.Detach())
		n++
	}
	if n != 3 || from.HasStatements() {
		t.Fatalf("moved %d statements, %v left behind", n, from.HasStatements())
	}
	if first, _ := to.FirstStatement(); first != a.AsStatement() {
		t.Fatal("the statements moved out of order")
	}

	// An empty list, and a loop body that stops early.
	for range ReadAhead(NewBlock().FirstStatement()) {
		t.Fatal("ReadAhead of an empty block yields an element")
	}
	seq := ReadAhead(to.FirstStatement())
	for range seq {
		break
	}
	n = 0
	for range seq {
		n++
	}
	if n != 3 {
		t.Fatalf("a second run of the same iterator yields %d elements, want 3", n)
	}
}

// IsQuoted tells {'a': b} from {a: b}.
func TestStringKeyIsQuoted(t *testing.T) {
	newStore(t)
	plain, ok := AsStringKey(rhino.IR_StringKeyWithValue("a", rhino.IR_Name("b")))
	if !ok || plain.IsQuoted() || plain.Key() != "a" {
		t.Fatal("{a: b}: want an unquoted key a")
	}
	quoted, ok := AsStringKey(rhino.IR_QuotedStringKey("a", rhino.IR_Name("b")))
	if !ok || !quoted.IsQuoted() || quoted.Key() != "a" {
		t.Fatal("{'a': b}: want a quoted key a")
	}
}
