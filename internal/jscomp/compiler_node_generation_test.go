package jscomp

import (
	"runtime"
	"testing"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: a compilation's nodes live in its node store (rhino/node_store.go),
// freed by the GC once nothing references a node of it. A single stale
// node handle held by anything long-lived (a pool, a package-level cache)
// would pin the whole AST of that compilation; this checks that finished
// compilations are released.
func TestNodeStoresAreReleased(t *testing.T) {
	externs := []*SourceFile{SourceFile_FromCode("externs.js", "/** @param {*} x */ function alert(x) {}")}
	src := `
/** @constructor @param {number} x */
function Point(x) { /** @type {number} */ this.x = x; }
/** @return {number} */
Point.prototype.twice = function() { return /** @type {number} */ (this.x * 2); };
alert(new Point(3).twice());
`
	var stores []weak.Pointer[rhino.NodeStore]
	compile := func() {
		rhino.ResetNodeStore()
		options := NewCompilerOptions()
		CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
		options.SetCheckTypes(true)
		c := NewCompiler()
		stores = append(stores, weak.Make(rhino.CurrentNodeStore()))
		if r := c.Compile(externs, []*SourceFile{SourceFile_FromCode("in.js", src)}, options); !r.Success {
			t.Fatalf("compilation failed: %v", r.Errors)
		}
	}
	const n = 20
	for i := 0; i < n; i++ {
		compile()
	}
	rhino.ResetNodeStore()
	for i := 0; i < 3; i++ {
		runtime.GC()
	}
	live := 0
	for _, w := range stores {
		if w.Value() != nil {
			live++
		}
	}
	t.Logf("%d of %d node stores still live", live, n)
	if live > 1 {
		t.Fatalf("node stores are not released: %d of %d still live", live, n)
	}
}
