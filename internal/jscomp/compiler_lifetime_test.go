package jscomp

import (
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

// Go-only: the node tables (rhino/node_tables.go) are released when the
// Compiler is collected, so nothing reachable from them may reach the
// Compiler. A type-checked compilation that never converts to colors is the
// case that exercises it: the node type table keeps every attached JSType,
// each JSType its registry, and the registry reports to the Compiler.
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
	}()
	for i := 0; i < 50; i++ {
		runtime.GC()
		select {
		case <-collected:
			return
		case <-time.After(100 * time.Millisecond):
		}
	}
	t.Fatal("the Compiler is still reachable (its node tables were never released)")
}
