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
	try.SetCatch(catch)
	if got, ok := try.Catch(); !ok || got != catch {
		t.Fatal("Catch() does not return the catch given to SetCatch")
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
