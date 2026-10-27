package jscomp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestCountAstSizesMatchesSequential(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("(function(){")
	for i := 0; i < 300; i++ {
		fmt.Fprintf(&sb, "function f%d(a,b){ if (a) { return [a, b, {x: a + %d}]; } for (var i = 0; i < b; i++) g(i, function(){ return i * %d; }); }\n", i, i, i)
	}
	sb.WriteString("})();\nvar z = 1;\n")
	var nodes []rhino.Node
	for _, js := range []string{sb.String(), "var a;", "", "f(1,2,3); var o = {a: [1, [2, [3]]]};"} {
		nodes = append(nodes, nuParse(t, js))
	}
	nodes = append(nodes, nodes[0].GetFirstChild(), rhino.NewNode(rhino.Token_EMPTY))
	want := make([]int, len(nodes))
	for i, n := range nodes {
		want[i] = NodeUtil_CountAstSize(n)
	}
	for _, workers := range []int{1, 2, 3, 4, 8} {
		got := countAstSizes(nodes, workers)
		for i := range nodes {
			if got[i] != want[i] {
				t.Fatalf("workers=%d node %d: got %d, want %d", workers, i, got[i], want[i])
			}
		}
	}
}
