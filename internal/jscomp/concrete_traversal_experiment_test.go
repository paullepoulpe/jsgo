package jscomp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// A test-only concrete callback records the complete control-flow contract and
// performs AST mutations that require the driver's cached-next rule. The same
// generator emits this driver and the production ReferenceCollector driver.
type concreteTraversalEvents struct {
	events []string
	mutate bool
	fail   bool
}

func concreteEventNode(n *rhino.Node) string {
	if n == nil {
		return "nil"
	}
	label := n.GetToken().String()
	if n.IsName() || n.IsGetProp() || n.IsStringLit() || n.IsStringKey() {
		label += ":" + n.GetString()
	}
	return label
}
func (c *concreteTraversalEvents) record(kind string, nt *NodeTraversal, n, parent *rhino.Node) {
	c.events = append(c.events, fmt.Sprintf("%s n=%s parent=%s current=%s scope=%s depth=%d fn=%s change=%s",
		kind, concreteEventNode(n), concreteEventNode(parent), concreteEventNode(nt.GetCurrentNode()),
		concreteEventNode(nt.GetScopeRoot()), nt.GetScope().GetDepth(), concreteEventNode(nt.GetEnclosingFunction()),
		concreteEventNode(nt.currentChangeScope)))
}
func (c *concreteTraversalEvents) ShouldTraverse(nt *NodeTraversal, n, parent *rhino.Node) bool {
	c.record("pre", nt, n, parent)
	return !(n.IsCall() && n.GetFirstChild().MatchesQualifiedName("skip_call"))
}
func (c *concreteTraversalEvents) Visit(nt *NodeTraversal, n, parent *rhino.Node) {
	c.record("post", nt, n, parent)
	if c.fail && n.IsName() && n.GetString() == "throw_here" {
		panic(jutil.NewRuntimeException("concrete-event-failure"))
	}
	if !c.mutate {
		return
	}
	if n.IsName() && n.GetString() == "replace_me" && parent.IsExprResult() {
		n.ReplaceWith(rhino.IR_Name("replaced").Srcref(n))
	}
	if n.IsExprResult() && n.GetFirstChild().IsName() {
		switch n.GetFirstChild().GetString() {
		case "remove_me":
			n.Detach()
		case "insert_after":
			rhino.IR_ExprResult(rhino.IR_Name("inserted")).InsertAfter(n)
		}
	}
}
func (c *concreteTraversalEvents) EnterScope(nt *NodeTraversal) {
	c.record("enter", nt, nt.GetScopeRoot(), nil)
}
func (c *concreteTraversalEvents) ExitScope(nt *NodeTraversal) {
	c.record("exit", nt, nt.GetScopeRoot(), nil)
}

func runConcreteTraversalEvents(t *testing.T, code, entry string, concrete, mutate, fail, obey bool) ([]string, string, string) {
	t.Helper()
	compiler := scopeTestNewCompiler(NewCompilerOptions())
	tree := scopeTestParse(t, compiler, code)
	creator := NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))
	cb := &concreteTraversalEvents{mutate: mutate, fail: fail}
	traversal := NodeTraversal_Builder_().SetCompiler(compiler).SetCallback(cb).SetScopeCreator(creator).Build()
	traversal.obeyDestructuringAndDefaultValueExecutionOrder = obey
	thrown := scopeTestCatch(func() {
		switch entry {
		case "roots":
			externs := scopeTestParse(t, compiler, "var externName;")
			rhino.IR_Root(externs.Detach(), tree.Detach())
			if concrete {
				traversal.eventsConcreteTraverseRoots(cb, externs, tree)
			} else {
				traversal.TraverseRoots(externs, tree)
			}
		case "with-scope":
			scope := creator.CreateScope(tree, nil)
			if concrete {
				traversal.eventsConcreteTraverseWithScope(cb, tree, scope)
			} else {
				traversal.TraverseWithScope(tree, scope)
			}
		case "at-scope":
			global := creator.CreateScope(tree, nil)
			scope := creator.CreateScope(tree.GetFirstChild(), global)
			if concrete {
				traversal.eventsConcreteTraverseAtScope(cb, scope)
			} else {
				traversal.TraverseAtScope(scope)
			}
		case "scope-root":
			root := tree.GetFirstChild()
			if concrete {
				traversal.eventsConcreteTraverseScopeRoot(cb, root)
			} else {
				traversal.TraverseScopeRoot(root)
			}
		case "tree":
			if concrete {
				traversal.eventsConcreteTraverse(cb, tree)
			} else {
				traversal.Traverse(tree)
			}
		default:
			t.Fatal("bad entry " + entry)
		}
	})
	message := ""
	if thrown != nil {
		message = fmt.Sprint(thrown)
	}
	return cb.events, tree.ToStringTree(), message
}

func TestConcreteTraversalEventAndMutationParity(t *testing.T) {
	code := `function outer(a = init()) {
   let [x = fallback(), y] = source;
   class C extends base() {
     [key()]() { return x; }
     [field()] = value();
     ordinary = other();
     static { let z = staticValue(); }
   }
   try { skip_call(hidden()); } catch (e) { use(e); }
   remove_me; replace_me; insert_after; keep_me;
   return function inner() { return a + x + y; };
 }`
	for _, entry := range []string{"tree", "roots", "with-scope", "at-scope", "scope-root"} {
		for _, obey := range []bool{false, true} {
			for _, mutate := range []bool{false, true} {
				t.Run(fmt.Sprint(entry, "/obey=", obey, "/mutate=", mutate), func(t *testing.T) {
					want, wantTree, wantError := runConcreteTraversalEvents(t, code, entry, false, mutate, false, obey)
					got, gotTree, gotError := runConcreteTraversalEvents(t, code, entry, true, mutate, false, obey)
					if !reflect.DeepEqual(got, want) {
						t.Fatalf("callback/state sequence differs\nwant %v\ngot %v", want, got)
					}
					if gotTree != wantTree || gotError != wantError {
						t.Fatalf("AST/error differs\nwant %s %q\ngot %s %q", wantTree, wantError, gotTree, gotError)
					}
					if len(got) == 0 || gotError != "" {
						t.Fatalf("did not exercise successful traversal: %q", gotError)
					}
				})
			}
		}
	}
}

func TestConcreteTraversalExceptionParity(t *testing.T) {
	code := "function outer() { var throw_here = 1; }"
	for _, entry := range []string{"tree", "roots", "with-scope", "at-scope", "scope-root"} {
		want, _, wantError := runConcreteTraversalEvents(t, code, entry, false, false, true, false)
		got, _, gotError := runConcreteTraversalEvents(t, code, entry, true, false, true, false)
		if !reflect.DeepEqual(got, want) || gotError != wantError || !strings.Contains(gotError, "concrete-event-failure") {
			t.Fatalf("%s exception contract differs: want %q got %q", entry, wantError, gotError)
		}
	}
}
