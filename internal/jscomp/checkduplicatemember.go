// Port of com/google/javascript/jscomp/StrictModeCheck.java's
// checkObjectLiteralOrClass (subset) @ closure-compiler v20260415.
//
// Object literals: fires JSC_DUPLICATE_MEMBER + JSC_DUP_VAR_DECLARATION
// when the same key appears twice in an OBJECTLIT — `{BB: 1, BB: 2}`.
// Java emits a WARNING at the first occurrence and an ERROR at the
// second.
//
// Class bodies: fires JSC_DUPLICATE_MEMBER (WARNING only) when a
// class declares two members with the same name, mirroring Java's
// StrictModeCheck.checkObjectLiteralOrClass. Getter + setter for the
// same name are compatible and do not trigger the diagnostic. Static
// and instance slots are tracked separately.

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

// Check walks the AST flagging duplicate members in every OBJECTLIT
// and every class body (CLASS → BLOCK).
func (c *CheckDuplicateMember) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script, "")
	return c.diags
}

func (c *CheckDuplicateMember) walk(n *rhino.Node, owner string) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.OBJECTLIT:
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
	case rhino.CLASS:
		// Mirrors StrictModeCheck.visit: n.isClass() →
		// checkObjectLiteralOrClass(n.getLastChild()).
		// CLASS children: [0]=name, [1]=superName, [2]=body BLOCK.
		if kids := n.Children(); len(kids) >= 3 {
			c.checkClassBody(kids[2])
		}
	}
	for _, ch := range n.Children() {
		c.walk(ch, owner)
	}
}

// checkClassBody ports StrictModeCheck.checkObjectLiteralOrClass for
// ES6 class bodies. It iterates the body BLOCK's children backwards
// (so the last declaration wins) and fires JSC_DUPLICATE_MEMBER on
// any earlier declaration with the same name.
//
// Four independent name-sets track instance-getter, instance-setter,
// static-getter, and static-setter slots — matching the Java 4-set
// design. A getter and a setter with the same name are compatible and
// do NOT trigger the diagnostic; two getters (or two non-setter
// members) DO.
func (c *CheckDuplicateMember) checkClassBody(body *rhino.Node) {
	// Four sets: getters / setters × instance / static.
	// key = member name, value = first-seen node (last in source).
	getters := map[string]*rhino.Node{}
	setters := map[string]*rhino.Node{}
	staticGetters := map[string]*rhino.Node{}
	staticSetters := map[string]*rhino.Node{}

	kids := body.Children()
	for i := len(kids) - 1; i >= 0; i-- {
		m := kids[i]
		switch m.Token() {
		case rhino.FUNCTION, rhino.CLASS_FIELD_DEF:
			// proceed
		default:
			continue // skip COMPUTED_PROP, EMPTY, etc.
		}
		first := m.FirstChild()
		if first == nil {
			continue
		}
		name := first.String_()
		if name == "" {
			continue
		}
		isStatic := m.JSDocInfo().IsStatic()
		isGetter := m.IsGetterDef()
		isSetter := m.IsSetterDef()

		if !isSetter {
			gs := getters
			if isStatic {
				gs = staticGetters
			}
			if _, dup := gs[name]; dup {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: m.Lineno(), Col: m.Charno(),
					Code:  diag.JSC_DUPLICATE_MEMBER,
					Level: diag.Warning,
					Message: fmt.Sprintf(
						"Class or object literal contains duplicate member %q. In non-strict code, the last duplicate will overwrite the others.",
						name),
				})
			} else {
				gs[name] = m
			}
		}
		if !isGetter {
			ss := setters
			if isStatic {
				ss = staticSetters
			}
			if _, dup := ss[name]; dup {
				c.diags = append(c.diags, diag.Diagnostic{
					Path: c.path, Line: m.Lineno(), Col: m.Charno(),
					Code:  diag.JSC_DUPLICATE_MEMBER,
					Level: diag.Warning,
					Message: fmt.Sprintf(
						"Class or object literal contains duplicate member %q. In non-strict code, the last duplicate will overwrite the others.",
						name),
				})
			} else {
				ss[name] = m
			}
		}
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
