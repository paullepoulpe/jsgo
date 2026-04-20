// Port of com/google/javascript/jscomp/AccessControls.java
// (née CheckAccessControls) @ closure-compiler v20260415 (subset).
//
// AccessControls in Java is a ~900-line pass that enforces three kinds
// of property-access discipline: (a) `@private` / `@protected` /
// `@public` visibility gating, (b) `@dict` vs `@struct` shape
// discipline (no `.p` on dicts, no `x['p']` on structs), and (c) the
// @package / @final / @const read-only flags. This port covers a
// subset: JSC_ILLEGAL_PROPERTY_ACCESS for "Cannot do '.' access on a
// dict" — enough for the `testgetpropdictN` fixtures. The visibility
// metadata is plumbed into ObjectType and JSDocInfo so a later slice
// can emit JSC_BAD_PRIVATE_PROPERTY_ACCESS /
// JSC_BAD_PROTECTED_PROPERTY_ACCESS without re-parsing JSDoc, but the
// enforcement itself is deferred: the baseline scope doesn't carry
// enough nominal class identity at check time to fire without
// false-positives.
//
// The pass runs after TypeCheck and does not mutate scope state — it
// scans the AST once up front to flip `IsDict` / `IsStruct` on the
// matching InstanceObjectTypes, then relies on TypeInference + the
// active local scope (via NodeTraversal) to resolve each GETPROP's
// receiver type.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type CheckAccessControls struct {
	path    string
	creator *TypedScopeCreator
	scope   *TypedScope
	diags   []diag.Diagnostic
}

// NewCheckAccessControls mirrors Java's AccessControls(AbstractCompiler).
func NewCheckAccessControls(path string, creator *TypedScopeCreator, scope *TypedScope) *CheckAccessControls {
	return &CheckAccessControls{path: path, creator: creator, scope: scope}
}

// Check drives the pass. Mirrors AccessControls.process in Java.
func (c *CheckAccessControls) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.propagateDictStructFlags(script)
	NewNodeTraversal(c, c.creator, c.scope).Traverse(script)
	return c.diags
}

// ShouldTraverse / Visit implement the NodeTraversal.Callback contract.
// Mirrors AbstractPostOrderCallback in Java — all logic runs in Visit.
func (c *CheckAccessControls) ShouldTraverse(_ *NodeTraversal, _ *rhino.Node, _ *rhino.Node) bool {
	return true
}

func (c *CheckAccessControls) Visit(t *NodeTraversal, n *rhino.Node, _ *rhino.Node) {
	if n.Token() == rhino.GETPROP {
		c.checkDotAccess(n, t.CurrentScope())
	}
}

// propagateDictStructFlags walks the AST and copies `@dict` / `@struct`
// from each FUNCTION / CLASS / var-initializer JSDocInfo onto the
// matching InstanceObjectType. Closure does this inside
// TypedScopeCreator's FunctionTypeBuilder; we do it here to avoid
// mutating scope.go.
func (c *CheckAccessControls) propagateDictStructFlags(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.FUNCTION, rhino.CLASS:
		if kids := n.Children(); len(kids) >= 1 {
			c.markInstance(kids[0].String_(), n.JSDocInfo())
		}
	case rhino.VAR, rhino.LET, rhino.CONST:
		// `/** @constructor @dict */ var F = function(){}` — the tag
		// can sit on either the declaration or the initializer.
		if nm := n.FirstChild(); nm != nil {
			if init := nm.FirstChild(); init != nil &&
				(init.Token() == rhino.FUNCTION || init.Token() == rhino.CLASS) {
				info := n.JSDocInfo()
				if info == nil {
					info = init.JSDocInfo()
				}
				c.markInstance(nm.String_(), info)
			}
		}
	}
	for _, child := range n.Children() {
		c.propagateDictStructFlags(child)
	}
}

func (c *CheckAccessControls) markInstance(name string, info *rhino.JSDocInfo) {
	if name == "" || info == nil || (!info.MakesDicts() && !info.MakesStructs()) {
		return
	}
	inst := c.scope.LookupNominalType(name)
	if inst == nil {
		return
	}
	if info.MakesDicts() {
		inst.IsDict = true
	}
	if info.MakesStructs() {
		inst.IsStruct = true
	}
}

// checkDotAccess emits JSC_ILLEGAL_PROPERTY_ACCESS when the receiver
// of a `.p` access resolves to an ObjectType tagged `@dict`. Mirrors
// the dict-access branch of Java's TypeCheck.visitGetProp / AccessControls.
func (c *CheckAccessControls) checkDotAccess(getprop *rhino.Node, scope *TypedScope) {
	kids := getprop.Children()
	if len(kids) < 2 {
		return
	}
	propNode := kids[1]
	// Legacy-mode nominal nullability: `@param {Foo}` is `Foo|null` in
	// Closure's default config. Strip the null / undefined members so
	// the dict detection still fires on the underlying nominal type.
	recvType := stripNullUndef(NewTypeInference(scope).InferType(kids[0]))
	obj := asDictObjectType(recvType)
	if obj == nil {
		return
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    propNode.Lineno(),
		Col:     propNode.Charno(),
		Code:    diag.JSC_ILLEGAL_PROPERTY_ACCESS,
		Level:   diag.Error,
		Message: "Cannot do '.' access on a dict",
	})
}

// asDictObjectType unwraps Templatized / Union wrappers looking for a
// dict-tagged ObjectType (tag inherited through the @extends chain).
// For unions this matches Java's visitGetPropDict: any dict alternate
// makes `.prop` illegal on the whole value. Returns nil when the
// receiver isn't a dict.
func asDictObjectType(t jstype.JSType) *jstype.ObjectType {
	switch v := t.(type) {
	case *jstype.ObjectType:
		if isDictThroughChain(v) {
			return v
		}
	case *jstype.TemplatizedType:
		if v.Base != nil && isDictThroughChain(v.Base) {
			return v.Base
		}
	case *jstype.UnionType:
		for _, alt := range v.Members() {
			if obj := asDictObjectType(alt); obj != nil {
				return obj
			}
		}
	}
	return nil
}

// isDictThroughChain reports whether `obj` or any superclass is
// tagged `@dict`. Mirrors FunctionType.makesDicts's walk up the
// @extends chain in Java.
func isDictThroughChain(obj *jstype.ObjectType) bool {
	for cur := obj; cur != nil; cur = cur.Super() {
		if cur.IsDict {
			return true
		}
	}
	return false
}
