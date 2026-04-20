// Port of com/google/javascript/jscomp/CheckSuspiciousCode.java's
// duplicate-member branch (subset) @ closure-compiler v20260415.
//
// Fires JSC_DUPLICATE_MEMBER + JSC_DUP_VAR_DECLARATION on object
// literals that list the same key twice — `{BB: 1, BB: 2}`. Java
// emits two diagnostics: a WARNING at the first occurrence with
// the "duplicate member" message, and an ERROR at the second
// occurrence describing the redefined variable.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type CheckDuplicateMember struct {
	path  string
	diags []diag.Diagnostic
}

func NewCheckDuplicateMember(path string) *CheckDuplicateMember {
	return &CheckDuplicateMember{path: path}
}

// Check walks the AST flagging duplicate STRING_KEY entries in
// every OBJECTLIT.
func (c *CheckDuplicateMember) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script, "")
	return c.diags
}

func (c *CheckDuplicateMember) walk(n *rhino.Node, owner string) {
	if n == nil {
		return
	}
	if n.Token() == rhino.OBJECTLIT {
		// Determine the receiver name to qualify the DUP_VAR
		// diagnostic with — matches how Java uses the containing
		// declaration's fully-qualified name (`a.BB` / `goog.a.BB`).
		scope := enclosingOwnerName(n, owner)
		seen := map[string]*rhino.Node{}
		for _, kv := range n.Children() {
			if kv.Token() != rhino.STRING_KEY {
				continue
			}
			name := kv.String_()
			if first, ok := seen[name]; ok {
				c.diags = append(c.diags,
					diag.Diagnostic{
						Path: c.path, Line: first.Lineno(), Col: first.Charno(),
						Code:  diag.JSC_DUPLICATE_MEMBER,
						Level: diag.Warning,
						Message: fmt.Sprintf(
							"Class or object literal contains duplicate member %q. In non-strict code, the last duplicate will overwrite the others.",
							name),
					},
					diag.Diagnostic{
						Path: c.path, Line: kv.Lineno(), Col: kv.Charno(),
						Code:  diag.JSC_DUP_VAR_DECLARATION,
						Level: diag.Error,
						Message: fmt.Sprintf(
							"variable %s.%s redefined, original definition at %s:%d",
							scope, name, c.path, first.Lineno()),
					},
				)
				continue
			}
			seen[name] = kv
		}
	}
	for _, ch := range n.Children() {
		c.walk(ch, owner)
	}
}

// enclosingOwnerName walks parents looking for the VAR / ASSIGN
// LHS that owns this OBJECTLIT, returning the dotted name used to
// qualify the JSC_DUP_VAR_DECLARATION diagnostic. Falls back to
// `<anonymous>` when the literal isn't directly assigned to a
// named slot.
func enclosingOwnerName(n *rhino.Node, fallback string) string {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.NAME:
			// VAR name — treat as the owner when VAR-initialised.
			if p := cur.Parent(); p != nil {
				switch p.Token() {
				case rhino.VAR, rhino.LET, rhino.CONST:
					return cur.String_()
				}
			}
		case rhino.ASSIGN:
			if kids := cur.Children(); len(kids) >= 1 {
				return flattenGetProp(kids[0])
			}
		}
	}
	if fallback != "" {
		return fallback
	}
	return "<anonymous>"
}

// flattenGetProp renders a NAME / GETPROP chain as a dotted path
// (`goog.a.BB`). Unknown tokens fall through to `<anonymous>`.
func flattenGetProp(n *rhino.Node) string {
	if n == nil {
		return "<anonymous>"
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			return flattenGetProp(kids[0]) + "." + kids[1].String_()
		}
	}
	return "<anonymous>"
}
