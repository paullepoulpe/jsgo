package jscomp

import (
	"os"
	"testing"
)

// Go-only: checks that ConvertTypesToColors really releases the JSType
// registry (every JSType points to it, so one JSType kept reachable keeps
// the whole type graph, its JSDocInfos and scopes alive through the
// optimizations).
//
// The check is strict (see typeRegistryReleased): a single stray JSType
// fails it. Opt-in until the known retainers are gone: the package-level
// rhino.propertySymbolList (KnownSymbolTypes), and detached Nodes that still
// carry a JSType and are kept reachable by their rhino.Node slab chunk.
// Run with JSGO_CHECK_TYPE_REGISTRY_RELEASE=1.
func TestTypeRegistryReleasedAfterTypesToColors(t *testing.T) {
	if os.Getenv("JSGO_CHECK_TYPE_REGISTRY_RELEASE") == "" {
		t.Skip("set JSGO_CHECK_TYPE_REGISTRY_RELEASE=1 (known retainers: rhino.propertySymbolList, detached typed Nodes in Node slabs)")
	}
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
