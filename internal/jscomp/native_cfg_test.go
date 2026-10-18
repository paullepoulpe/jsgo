package jscomp

import (
	"fmt"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The generic graph is constructed only as a differential oracle. The native
// path must preserve every node ID, edge, priority, and independent in/out order
// without retaining generic nodes or edges behind an adapter.
func TestNativeCFGConstructionDifferential(t *testing.T) {
	cases := []struct {
		body  string
		async bool
	}{
		{"var x=1; return x;", false},
		{"var x; if(param1){x=1;} else{x=2;} return x;", false},
		{"var x=0; while(param1){ if(x)continue; x++; } return x;", false},
		{"var x=0; do { x++; if(param1)break; } while(param2); return x;", false},
		{"var x=0; for(var i=0;i<10;i++){ if(i===2)continue; x+=i; } return x;", false},
		{"for(var key in param1) { if(key) break; } return key;", false},
		{"for(let [x,y] of param1) { param2(x,y); }", false},
		{"for await(const x of param1) { if(x)continue; param2(x); }", true},
		{"outer: for(;;) { for(;;) { if(param1)break outer; continue outer; } }", false},
		{"try { param1(); return 1; } catch(e){ param2(e); } finally { param2(); }", false},
		{"try { try { param1(); } finally { param2(); } } catch(e) { throw e; }", false},
		{"try { if(param1)return 1; throw 2; } finally { if(param2)return 3; }", false},
		{"var x=0; switch(param1){case 1:x=1;break;default:x=2;case 3:x++;} return x;", false},
		{"function nested(x){if(x)return x; return 0;} var x=param1?.(param2); return x ?? 0;", false},
		{"class C { static { if(param1) param2(); } method(){ return param1(); } } return C;", false},
		{"return 1; var dead=param1(); while(dead){ dead--; }", false},
	}
	for index, source := range cases {
		for _, traverse := range []bool{false, true} {
			t.Run(fmt.Sprintf("case=%d/functions=%t", index, traverse), func(t *testing.T) {
				tester := ReachingUseDefTester_Create()
				script := tester.parseScript(source.body, source.async)
				root := script.GetFirstChild()
				makeBuilder := func() *ControlFlowAnalysis_Builder {
					return ControlFlowAnalysis_Builder_().SetCompiler(tester.compiler).
						SetCfgRoot(root).SetTraverseFunctions(traverse).SetIncludeEdgeAnnotations(true)
				}
				generic := newIndexedFlowTopology[*rhino.Node](makeBuilder().ComputeCfg())
				native := makeBuilder().ComputeIndexedCfg()
				if generic == nil || native == nil {
					t.Fatal("valid fixture did not produce a topology")
				}
				if len(native.nodes) != 0 || len(native.originalEdges) != 0 {
					t.Fatal("native construction retained a generic object graph")
				}
				if generic.entry != native.entry || generic.implicit != native.implicit || !slices.Equal(generic.values, native.values) {
					t.Fatal("node creation/entry/implicit-return order differs")
				}
				if !slices.Equal(generic.edges, native.edges) || !slices.Equal(generic.records, native.records) ||
					!slices.Equal(generic.incoming, native.incoming) || !slices.Equal(generic.outgoing, native.outgoing) {
					t.Fatal("CFG branch records or incoming/outgoing encounter order differs")
				}
				for a := range native.values {
					for b := range native.values {
						if generic.forwardOrder(int32(a), int32(b)) != native.forwardOrder(int32(a), int32(b)) ||
							generic.backwardOrder(int32(a), int32(b)) != native.backwardOrder(int32(a), int32(b)) {
							t.Fatalf("priority comparison differs for nodes %d/%d", a, b)
						}
					}
				}
			})
		}
	}
}
