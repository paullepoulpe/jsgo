// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20260415
//
// TypeCheck is the big Java pass that walks the AST after scope+inference
// and emits JSError objects. Tier-0 implements only the
// "initializing variable" case: a VAR with a JSDoc-declared type whose
// initializer has an inferred type that is not a subtype of the declared
// type. Matches DiagnosticType JSC_TYPE_MISMATCH.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type TypeCheck struct {
	scope *TypedScope
	infer *TypeInference
	path  string
}

func NewTypeCheck(path string, scope *TypedScope, infer *TypeInference) *TypeCheck {
	return &TypeCheck{scope: scope, infer: infer, path: path}
}

// Check walks the SCRIPT and returns accumulated diagnostics.
func (c *TypeCheck) Check(script *rhino.Node) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, child := range script.Children() {
		if child.Token() != rhino.VAR {
			continue
		}
		out = append(out, c.checkVar(child)...)
	}
	return out
}

func (c *TypeCheck) checkVar(v *rhino.Node) []diag.Diagnostic {
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return nil
	}
	init := name.FirstChild()
	if init == nil {
		return nil
	}
	declared := c.scope.GetDeclaredType(name.String_())
	if declared == nil {
		return nil
	}
	found := c.infer.InferType(init)
	if found.IsSubtypeOf(declared) {
		return nil
	}
	// Matches TypeCheck.java "initializing variable" JSC_TYPE_MISMATCH.
	msg := fmt.Sprintf("initializing variable\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    init.Lineno(),
		Col:     init.Charno(),
		Code:    diag.JSC_TYPE_MISMATCH,
		Level:   diag.Error,
		Message: msg,
	}}
}
