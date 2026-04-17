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
		switch {
		case isDeclToken(child.Token()):
			out = append(out, c.checkVar(child)...)
		case child.Token() == rhino.EXPR_RESULT:
			if a := child.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
				out = append(out, c.checkAssign(a)...)
			}
		}
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

// checkAssign handles bare `IDENT = LITERAL;` statements. Matches
// TypeCheck.java's "assignment" JSC_TYPE_MISMATCH case.
func (c *TypeCheck) checkAssign(a *rhino.Node) []diag.Diagnostic {
	name := a.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return nil
	}
	rhs := name
	// ASSIGN has two children: lhs NAME, rhs literal. NAME is first; rhs
	// is the second child — access via Children().
	kids := a.Children()
	if len(kids) < 2 {
		return nil
	}
	rhs = kids[1]
	declared := c.scope.GetDeclaredType(name.String_())
	if declared == nil {
		return nil
	}
	found := c.infer.InferType(rhs)
	if found.IsSubtypeOf(declared) {
		return nil
	}
	// TypeCheck.java reports the "assignment" variant at the LHS
	// position, not the RHS (unlike "initializing variable").
	msg := fmt.Sprintf("assignment\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    name.Lineno(),
		Col:     name.Charno(),
		Code:    diag.JSC_TYPE_MISMATCH,
		Level:   diag.Error,
		Message: msg,
	}}
}
