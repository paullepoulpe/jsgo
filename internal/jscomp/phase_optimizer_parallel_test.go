package jscomp

import (
	"fmt"
	"io"
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

func TestCountAstSizesRecoversWorkerPanic(t *testing.T) {
	rhino.ResetNodeStore()
	valid := rhino.IR_Empty()
	nodes := make([]rhino.Node, 260)
	for i := range nodes {
		nodes[i] = valid
	}
	// Inject a NilNode past the first chunk so child worker goroutines hit a panic
	// inside NodeUtil_CountAstSize; countAstSizes must recover and re-panic on caller.
	nodes[150] = rhino.NilNode

	var caught any
	func() {
		defer func() { caught = recover() }()
		_ = countAstSizes(nodes, 4)
	}()
	if caught == nil {
		t.Fatal("expected countAstSizes worker panic to be recovered and re-panicked on caller")
	}
}

func TestParallelFastGzipFlushBatchRecoversWorkerPanic(t *testing.T) {
	p := &parallelFastGzipOutputStream{
		out:        io.Discard,
		numThreads: 2,
		batch:      [][]byte{[]byte("chunk-0"), []byte("chunk-1")},
		// Intentionally undersized workers slice (len 1 < batch len 2) so worker 1
		// panics inside its child goroutine while worker 0 succeeds.
		workers: make([]parallelGzipWorker, 1),
	}

	var caught any
	func() {
		defer func() { caught = recover() }()
		p.flushBatch()
	}()
	if caught == nil {
		t.Fatal("expected parallelFastGzipOutputStream.flushBatch worker panic to be recovered and re-panicked on caller")
	}
	if len(p.batch) != 0 {
		t.Fatalf("expected p.batch to be cleared after recovered panic, got len=%d", len(p.batch))
	}
}
