package jscomp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func referenceCollectorVarFilterFixture(filter func(*Var) bool) (*ReferenceCollector, *Scope, *NodeTraversal, *CompilerInput) {
	root := rhino.NewNode(rhino.Token_ROOT)
	scope := Scope_CreateGlobalScope(root)
	input := &CompilerInput{id: rhino.NewInputId("reference-filter.js")}
	traversal := &NodeTraversal{
		scopes:        []scopeStackEntry{{node: root, untypedScope: scope}},
		compilerInput: input,
	}
	collector := NewReferenceCollectorWithVarFilter(nil, nil, nil, filter)
	collector.blockStack = append(collector.blockStack, &BasicBlock{})
	return collector, scope, traversal, input
}

func referenceCollectorAssertNodes(t *testing.T, collection *ReferenceCollection, want ...*rhino.Node) {
	t.Helper()
	var got []*rhino.Node
	if collection != nil {
		for ref := range collection.All() {
			got = append(got, ref.GetNode())
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("reference node order = %v, want %v", got, want)
	}
}

func TestReferenceCollectorVarFilterEverySameVarCall(t *testing.T) {
	for _, fillers := range []int{0, 17} {
		t.Run(fmt.Sprintf("other_variables_%d", fillers), func(t *testing.T) {
			active, calls := false, 0
			var filtered []*Var
			collector, scope, traversal, input := referenceCollectorVarFilterFixture(func(v *Var) bool {
				if !active {
					return true
				}
				calls++
				filtered = append(filtered, v)
				return calls%2 == 1
			})
			v := NewVar("x", rhino.IR_Name("x"), scope, 0, input, nil)
			first := rhino.IR_Name("x")
			collector.AddReference(v, first, traversal)
			for i := 0; i < fillers; i++ {
				name := fmt.Sprintf("filler%d", i)
				collector.AddReference(NewVar(name, rhino.IR_Name(name), scope, i+1, input, nil), rhino.IR_Name(name), traversal)
			}
			collection := collector.GetReferences(v) // prime the exact-Var cache
			active = true
			want := []*rhino.Node{first}
			for i := 1; i <= 6; i++ {
				n := rhino.IR_Name("x")
				collector.AddReference(v, n, traversal)
				if i%2 == 1 {
					want = append(want, n)
				}
				referenceCollectorAssertNodes(t, collection, want...)
			}
			if calls != 6 || len(filtered) != 6 {
				t.Fatalf("predicate calls = %d, want 6 on repeated cached Var", calls)
			}
			for _, got := range filtered {
				if got != v {
					t.Fatal("predicate received a different variable identity")
				}
			}
			for ref := range collection.All() {
				if ref.scope != scope || ref.inputId != input.id || ref.basicBlock != collector.blockStack[0] {
					t.Fatal("accepted reference lost traversal scope/input/block")
				}
			}
		})
	}
}

func TestReferenceCollectorVarFilterRejectsBeforeTraversalAndAllocation(t *testing.T) {
	accept, calls := true, 0
	collector, scope, traversal, input := referenceCollectorVarFilterFixture(func(*Var) bool {
		calls++
		return accept
	})
	collector.enableArenaPool()
	v := NewVar("x", rhino.IR_Name("x"), scope, 0, input, nil)
	first := rhino.IR_Name("x")
	collector.AddReference(v, first, traversal)
	collection := collector.GetReferences(v)
	refs, cols, vars, arena := len(collector.refSlab), len(collector.colSlab), len(collector.vars), collector.arena
	accept = false
	collector.blockStack = nil
	collector.AddReference(v, nil, nil) // rejected before invalid traversal/block access
	if calls != 2 || len(collector.refSlab) != refs || len(collector.colSlab) != cols || len(collector.vars) != vars || collector.arena != arena {
		t.Fatal("cached rejection skipped predicate or allocated/changed collector storage")
	}
	referenceCollectorAssertNodes(t, collection, first)
}

func TestReferenceCollectorVarFilterReentrant(t *testing.T) {
	for _, sameVar := range []bool{false, true} {
		t.Run(fmt.Sprintf("same_variable_%t", sameVar), func(t *testing.T) {
			var collector *ReferenceCollector
			var traversal *NodeTraversal
			var x, nestedVar *Var
			var events []string
			active := false
			outerRunning := false
			inner := rhino.IR_Name("nested")
			filter := func(v *Var) bool {
				if !active {
					return true
				}
				events = append(events, "enter:"+v.GetName())
				if v == x && !outerRunning {
					outerRunning = true
					collector.AddReference(nestedVar, inner, traversal)
				}
				events = append(events, "exit:"+v.GetName())
				return true
			}
			var scope *Scope
			var input *CompilerInput
			collector, scope, traversal, input = referenceCollectorVarFilterFixture(filter)
			x = NewVar("x", rhino.IR_Name("x"), scope, 0, input, nil)
			y := NewVar("y", rhino.IR_Name("y"), scope, 1, input, nil)
			firstX, firstY := rhino.IR_Name("x"), rhino.IR_Name("y")
			collector.AddReference(x, firstX, traversal)
			collector.AddReference(y, firstY, traversal)
			xCollection, yCollection := collector.GetReferences(x), collector.GetReferences(y)
			collector.GetReferences(x) // callback must run even on this cache hit
			nestedVar = y
			if sameVar {
				nestedVar = x
			}
			active = true
			outer := rhino.IR_Name("outer")
			collector.AddReference(x, outer, traversal)
			wantEvents := []string{"enter:x", "enter:" + nestedVar.GetName(), "exit:" + nestedVar.GetName(), "exit:x"}
			if !reflect.DeepEqual(events, wantEvents) {
				t.Fatalf("predicate reentry order = %v, want %v", events, wantEvents)
			}
			if sameVar {
				referenceCollectorAssertNodes(t, xCollection, firstX, inner, outer)
				referenceCollectorAssertNodes(t, yCollection, firstY)
			} else {
				referenceCollectorAssertNodes(t, xCollection, firstX, outer)
				referenceCollectorAssertNodes(t, yCollection, firstY, inner)
			}
		})
	}
}

func TestReferenceCollectorVarFilterPanicOnCachedVar(t *testing.T) {
	sentinel := new(int)
	calls := 0
	collector, scope, traversal, input := referenceCollectorVarFilterFixture(func(*Var) bool {
		calls++
		if calls == 2 {
			panic(sentinel)
		}
		return true
	})
	v := NewVar("x", rhino.IR_Name("x"), scope, 0, input, nil)
	first := rhino.IR_Name("x")
	collector.AddReference(v, first, traversal)
	collection := collector.GetReferences(v)
	func() {
		defer func() {
			if got := recover(); got != sentinel {
				t.Fatalf("cached predicate panic = %v, want sentinel", got)
			}
		}()
		collector.AddReference(v, nil, nil)
	}()
	if calls != 2 {
		t.Fatalf("predicate calls = %d, want 2", calls)
	}
	referenceCollectorAssertNodes(t, collection, first)
}
