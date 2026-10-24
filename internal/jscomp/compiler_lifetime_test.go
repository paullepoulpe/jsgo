package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"runtime"
	"testing"
	"time"
)

// Go-only: checks that ConvertTypesToColors really releases the JSType
// registry (every JSType points to it, so one JSType kept reachable keeps
// the whole type graph, its JSDocInfos and scopes alive through the
// optimizations).
//
// The check is strict (see typeRegistryReleased): a single stray JSType
// fails it. Nodes refer to their types through handles that
// ClearJSTypeRegistry releases wholesale (rhino/node_tables.go), so even a
// detached typed Node that is still reachable does not retain the registry.
func TestTypeRegistryReleasedAfterTypesToColors(t *testing.T) {
	options := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	options.SetCheckTypes(true)
	externs := []*SourceFile{SourceFile_FromCode("externs.js", `
/** @param {*} x */ function alert(x) {}
`)}
	inputs := []*SourceFile{SourceFile_FromCode("in.js", `
/** @constructor @param {number} x */
function Point(x) { /** @type {number} */ this.x = x; }
/** @return {number} */
Point.prototype.twice = function() { return /** @type {number} */ (this.x * 2); };
alert(new Point(3).twice());
`)}
	c := NewCompiler()
	result := c.Compile(externs, inputs, options)
	if !result.Success {
		t.Fatalf("compilation failed: %v", result.Errors)
	}
	if !c.HasTypeCheckingRun() || !c.IsTypeRegistryCleared() {
		t.Fatal("expected type checking to have run and the registry to be cleared")
	}
	if !c.typeRegistryReleased() {
		t.Error("the JSTypeRegistry is still reachable after typesToColors")
	}
}

// Go-only: a Compiler is collectable once its node store is no longer
// current. A type-checked compilation that never converts to colors is the
// case that exercises it: the store's type table keeps every attached
// JSType, each JSType its registry, and the registry reports to the
// Compiler, so nothing else may keep the store or the Compiler alive.
func TestCompilerIsCollectableWithTypesAttached(t *testing.T) {
	collected := make(chan struct{}, 1)
	func() {
		options := NewCompilerOptions()
		options.SetCheckTypes(true)
		externs := []*SourceFile{SourceFile_FromCode("externs.js", "/** @param {*} x */ function alert(x) {}")}
		inputs := []*SourceFile{SourceFile_FromCode("in.js", `
/** @constructor @param {number} x */
function Point(x) { /** @type {number} */ this.x = x; }
alert(new Point(3).x);
`)}
		c := NewCompiler()
		c.Init(externs, inputs, options)
		c.Parse()
		c.Check()
		if !c.HasTypeCheckingRun() || c.IsTypeRegistryCleared() {
			t.Fatal("expected type checking to have run without clearing the registry")
		}
		runtime.AddCleanup(c, func(ch chan struct{}) { ch <- struct{}{} }, collected)
		rhino.ResetNodeStore()
	}()
	for i := 0; i < 50; i++ {
		runtime.GC()
		select {
		case <-collected:
			return
		case <-time.After(100 * time.Millisecond):
		}
	}
	t.Fatal("the Compiler is still reachable after its node store stopped being current")
}

// Go-only: a caller may keep the AST after the last use of its Compiler (a
// test harness that prints the tree after reading the diagnostics). The
// nodes' JSDoc, colors and cold data live in the Compiler's node tables, so
// those must stay until the nodes are unreachable too, not only the
// Compiler: once the Compiler is collected, a kept node still has its JSDoc.
func TestNodeTablesOutliveTheCompilerWhileNodesAreReachable(t *testing.T) {
	collected := make(chan struct{}, 1)
	var root rhino.Node
	func() {
		c := NewCompiler()
		c.Init([]*SourceFile{SourceFile_FromCode("externs.js", "")},
			[]*SourceFile{SourceFile_FromCode("in.js", "/** @const */ var a = 1;\n/** @return {number} */ function f() { return a; }\n")},
			NewCompilerOptions())
		root = c.ParseInputs()
		if c.GetErrorCount() > 0 {
			t.Fatalf("parse errors: %v", c.GetErrors())
		}
		runtime.AddCleanup(c, func(ch chan struct{}) { ch <- struct{}{} }, collected)
	}()
	done := false
	for i := 0; i < 50 && !done; i++ {
		runtime.GC()
		select {
		case <-collected:
			done = true
		case <-time.After(20 * time.Millisecond):
		}
	}
	if !done {
		t.Fatal("the Compiler was not collected")
	}
	// Give the tables' cleanup, which would run right after the Compiler's
	// in the broken lifetime, every chance to run.
	for i := 0; i < 5; i++ {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	script := root.GetLastChild().GetFirstChild()
	v := script.GetFirstChild()
	fn := v.GetNext()
	if v.GetJSDocInfo() == nil || !v.GetJSDocInfo().IsConstant() {
		t.Errorf("var a lost its JSDoc after the Compiler was collected: %v", v.GetJSDocInfo())
	}
	if fn.GetJSDocInfo() == nil || !fn.GetJSDocInfo().HasReturnType() {
		t.Errorf("function f lost its JSDoc after the Compiler was collected: %v", fn.GetJSDocInfo())
	}
	runtime.KeepAlive(root)
}
