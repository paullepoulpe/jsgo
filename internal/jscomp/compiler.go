// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20260415
//
// The real Compiler is a large orchestrator running many passes; this
// scaffold wires only the three pieces Tier-0 needs (parse → scope →
// check) and returns diagnostics. Kept intentionally narrow so we don't
// bake design decisions before the corpus demands them.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckSource runs the Tier-0 pipeline on one source file and returns
// accumulated diagnostics. A parse error is returned as err; the caller
// decides whether to surface it as a diagnostic or as a fatal.
func CheckSource(path, src string) ([]diag.Diagnostic, error) {
	script, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	scope := creator.CreateScope(script)
	check := NewTypeCheck(path, creator, scope)
	return check.Check(script), nil
}
