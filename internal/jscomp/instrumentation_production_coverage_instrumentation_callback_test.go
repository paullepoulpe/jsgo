package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestProductionCoverageInstrumentationCallback(t *testing.T) {
	c := NewCompiler()
	c.InitOptions(NewCompilerOptions())
	jsRoot := c.ParseTestCode(
		`function foo(x) {
  if (x > 0) {
    return x ? 1 : 2;
  }
  switch (x) {
    case -1:
      return x || 0;
  }
  return x && (x ?? 3);
}`)
	root := rhino.IR_Root(jsRoot)

	cb := NewProductionCoverageInstrumentationCallback(c, "__cov")
	NodeTraversal_Traverse(c, root, cb)
	vm := cb.GetInstrumentationMapping()

	if got, ok := vm.LookupNewName(" FileNames"); !ok || got != `["testcode"]` {
		t.Fatalf("FileNames = (%q, %v), want (%q, true)", got, ok, `["testcode"]`)
	}
	if got, ok := vm.LookupNewName(" FunctionNames"); !ok || got != `["foo"]` {
		t.Fatalf("FunctionNames = (%q, %v), want (%q, true)", got, ok, `["foo"]`)
	}
	if got, ok := vm.LookupNewName(" Types"); !ok || got != `["BRANCH","BRANCH_DEFAULT","FUNCTION"]` {
		t.Fatalf("Types = (%q, %v), want (%q, true)", got, ok, `["BRANCH","BRANCH_DEFAULT","FUNCTION"]`)
	}

	// Empty callback mapping must serialize empty JSON arrays ("[]"), not "null".
	emptyCb := NewProductionCoverageInstrumentationCallback(c, "__cov")
	emptyVm := emptyCb.GetInstrumentationMapping()
	if got, ok := emptyVm.LookupNewName(" Types"); !ok || got != "[]" {
		t.Fatalf("empty Types = (%q, %v), want ([], true)", got, ok)
	}
}
