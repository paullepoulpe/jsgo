package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestStripCodeOptChainCallParity(t *testing.T) {
	opts := NewCompilerOptions()
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetStripTypes(NewOrderedSet("goog.debug.Logger"))

	src := `
var goog = {debug: {Logger: {getLogger: function(name) { return null; }}}};
function wrap(testNameProvider, obj) {
  const originalTest = testNameProvider?.();
  let methodResult = obj?.run();
  var collapsed$var = testNameProvider?.();
  var logger = goog.debug.Logger.getLogger('demo');
  const commaOpt = (0, logger)?.();
  return [originalTest, methodResult, collapsed$var, commaOpt];
}
`
	comp := NewCompiler()
	res := comp.Compile(
		[]*SourceFile{},
		[]*SourceFile{SourceFile_FromCode("optchain_strip.js", src)},
		opts,
	)
	if !res.Success {
		t.Fatalf("Compile failed: %v", res.Errors)
	}
	got := comp.ToSource()
	if got == "" {
		t.Fatalf("expected non-empty output")
	}

	// Also verify direct exported helper calls on an OPTCHAIN_CALL AST node do not panic.
	sc := StripCodeTest_GetStripCodeInstance(comp)
	script := comp.ParseTestCode("const x = fn?.();")
	var optChainCall rhino.Node
	var walk func(rhino.Node)
	walk = func(n rhino.Node) {
		if n.IsNil() || !optChainCall.IsNil() {
			return
		}
		if n.IsOptChainCall() {
			optChainCall = n
			return
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
	}
	walk(script)
	if optChainCall.IsNil() {
		t.Fatalf("expected OPTCHAIN_CALL node in parsed AST")
	}
	if sc.IsCallWhoseReturnValueShouldBeStripped(optChainCall) {
		t.Fatalf("expected IsCallWhoseReturnValueShouldBeStripped(OPTCHAIN_CALL) == false")
	}
	if sc.IsMethodOrCtorCallThatTriggersRemoval(nil, optChainCall, optChainCall.GetParent()) {
		t.Fatalf("expected IsMethodOrCtorCallThatTriggersRemoval(OPTCHAIN_CALL) == false")
	}
}
