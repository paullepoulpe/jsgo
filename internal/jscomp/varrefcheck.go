// Port of com/google/javascript/jscomp/VariableReferenceCheck.java +
// the duplicate-decl branch of com/google/javascript/jscomp/VarCheck.java
// @ closure-compiler v20260415 (narrow; intra-file scope only).
//
// Java's VariableReferenceCheck walks the AST per-scope and records
// every value-namespace declaration (var/let/const/function/class/
// @param). A second declaration of the same name in the same scope
// fires JSC_VAR_MULTIPLY_DECLARED_ERROR, with a message pointing at
// the first occurrence's file:line:col. jsgo doesn't load the real
// externs bundle yet, so we can only detect duplicates that both
// live in the source file under check (the externs-vs-source
// collisions — ~26 extracted fixtures against es6.js's Disposable
// etc. — need P2-1 to land first).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// VariableReferenceCheck is the NodeTraversal.Callback that reports
// intra-file duplicate `var` declarations. Mirrors the pass of the
// same name on the Java side.
type VariableReferenceCheck struct {
	path  string
	diags []diag.Diagnostic

	// Stack of scope maps — one entry per enclosing function /
	// script. Each map key is a declared name and the value is the
	// first occurrence's source position, used to build the "First
	// occurrence: <path>:<line>:<col>" message.
	scopes []map[string]declSite

	// suppressDepth is non-zero whenever we're inside a function
	// whose JSDoc carries `@suppress {duplicate}`. While > 0,
	// recordDecl skips emitting the duplicate diagnostic. Mirrors
	// Java's per-pass suppression-stack behaviour.
	suppressDepth int
}

type declSite struct {
	line int
	col  int
}

func NewVariableReferenceCheck(path string, externScope *TypedScope) *VariableReferenceCheck {
	v := &VariableReferenceCheck{
		path:   path,
		scopes: []map[string]declSite{{}},
	}
	// Seed the global frame with every name declared by the externs
	// scope. Re-declaring one of these in user source fires
	// JSC_VAR_MULTIPLY_DECLARED_ERROR at the source-side position,
	// which matches Java's oracle.
	if externScope != nil {
		for name := range externScope.declared {
			v.scopes[0][name] = declSite{line: 0, col: 0}
		}
	}
	return v
}

// Check runs the traversal and returns the collected diagnostics.
// Callers are expected to build a fresh VariableReferenceCheck per
// file, matching Java's Compiler pass-instance pattern.
func (v *VariableReferenceCheck) Check(script *rhino.Node) []diag.Diagnostic {
	v.diags = nil
	v.walk(script, nil)
	return v.diags
}

func (v *VariableReferenceCheck) walk(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	// Enter a new scope frame for every FUNCTION (including arrows);
	// CLASS bodies don't create a value-namespace scope of their own
	// in the Java model, so we don't push one here.
	if n.Token() == rhino.FUNCTION {
		v.scopes = append(v.scopes, map[string]declSite{})
		// `@suppress {duplicate}` on the enclosing function silences
		// duplicate-declaration diagnostics for everything inside its
		// body — matches the Java pass's gate.
		suppressed := false
		if info := n.JSDocInfo(); info != nil && info.HasSuppression("duplicate") {
			v.suppressDepth++
			suppressed = true
		}
		defer func() {
			v.scopes = v.scopes[:len(v.scopes)-1]
			if suppressed {
				v.suppressDepth--
			}
		}()
		// Function parameters are declarations in the new scope.
		kids := n.Children()
		if len(kids) >= 2 {
			for _, p := range kids[1].Children() {
				v.recordDecl(p, p.String_())
			}
		}
	}
	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		if name := n.FirstChild(); name != nil && name.Token() == rhino.NAME {
			v.recordDecl(name, name.String_())
		}
	case rhino.FUNCTION, rhino.CLASS:
		// Named function / class declarations create a binding in the
		// enclosing scope. Skip members — class methods (FUNCTION
		// children of a CLASS-owned BLOCK) and object-literal methods
		// (FUNCTION under STRING_KEY under OBJECTLIT) are properties,
		// not value-namespace declarations.
		if isMemberContext(parent) {
			break
		}
		if parent != nil && (parent.Token() == rhino.BLOCK || parent.Token() == rhino.SCRIPT) {
			if kids := n.Children(); len(kids) >= 1 && kids[0].String_() != "" {
				name := kids[0]
				// The declaration lands in the enclosing scope (the
				// one above the just-pushed function scope). Use
				// scopes[-2] when we just pushed a function frame.
				if n.Token() == rhino.FUNCTION {
					target := v.scopes[len(v.scopes)-2]
					v.recordDeclInto(target, name, name.String_())
				} else {
					v.recordDecl(name, name.String_())
				}
			}
		}
	}
	for _, c := range n.Children() {
		v.walk(c, n)
	}
}

func (v *VariableReferenceCheck) recordDecl(nameNode *rhino.Node, name string) {
	v.recordDeclInto(v.scopes[len(v.scopes)-1], nameNode, name)
}

// isMemberContext reports whether `parent` (the node above a FUNCTION /
// CLASS) is a class-body or object-literal slot, where the nested
// definition is a member rather than a var-like declaration.
func isMemberContext(parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	// Class methods live under the CLASS's BLOCK body (the third
	// child of a CLASS node).
	if parent.Token() == rhino.BLOCK && parent.Parent() != nil && parent.Parent().Token() == rhino.CLASS {
		return true
	}
	// Object-literal method shorthand wraps the FUNCTION under a
	// STRING_KEY.
	if parent.Token() == rhino.STRING_KEY {
		return true
	}
	return false
}

func (v *VariableReferenceCheck) recordDeclInto(scope map[string]declSite, nameNode *rhino.Node, name string) {
	if name == "" {
		return
	}
	if prev, seen := scope[name]; seen {
		if v.suppressDepth > 0 {
			return
		}
		v.diags = append(v.diags, diag.Diagnostic{
			Path:  v.path,
			Line:  nameNode.Lineno(),
			Col:   nameNode.Charno(),
			Code:  diag.JSC_VAR_MULTIPLY_DECLARED_ERROR,
			Level: diag.Error,
			Message: fmt.Sprintf("Variable %s declared more than once. First occurrence: %s:%d:%d",
				name, v.path, prev.line, prev.col),
		})
		return
	}
	scope[name] = declSite{line: nameNode.Lineno(), col: nameNode.Charno()}
}
