package jscomp

import (
	"runtime"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: every Compiler allocates its nodes in its own generation of
// chunks (rhino/node_store.go), released as a whole once nothing references
// a node of it. A single stale *Node held by anything long-lived (a pool, a
// package-level cache) would pin the whole AST of that compilation; this
// checks that finished compilations are released.
func TestNodeGenerationsAreReleased(t *testing.T) {
	if rhino.NodePoisonMode {
		t.Skip("poison mode keeps the old chunks alive on purpose")
	}
	externs := []*SourceFile{SourceFile_FromCode("externs.js", "/** @param {*} x */ function alert(x) {}")}
	src := `
/** @constructor @param {number} x */
function Point(x) { /** @type {number} */ this.x = x; }
/** @return {number} */
Point.prototype.twice = function() { return /** @type {number} */ (this.x * 2); };
alert(new Point(3).twice());
`
	compile := func() {
		options := NewCompilerOptions()
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
		options.SetCheckTypes(true)
		c := NewCompiler()
		if r := c.Compile(externs, []*SourceFile{SourceFile_FromCode("in.js", src)}, options); !r.Success {
			t.Fatalf("compilation failed: %v", r.Errors)
		}
	}
	settle := func() int {
		for i := 0; i < 3; i++ {
			runtime.GC()
		}
		runtime.Gosched()
		for i := 0; i < 3; i++ {
			runtime.GC()
		}
		return rhino.NodeChunksLive()
	}
	compile()
	before := settle()
	const n = 20
	for i := 0; i < n; i++ {
		compile()
	}
	after := settle()
	t.Logf("live chunks: %d before, %d after %d more compilations (%d allocated in all)", before, after, n, rhino.NodeChunksAllocated())
	if after > before+n/2 {
		t.Fatalf("node generations are not released: %d live chunks after %d compilations (was %d)", after, n, before)
	}
}
