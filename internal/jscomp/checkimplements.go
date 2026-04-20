// Port of com/google/javascript/jscomp/TypeCheck.java's validateImplements
// branch + com/google/javascript/jscomp/CheckMissingOverrides.java
// @ closure-compiler v20260415 (narrow; no structural-interface check,
// no `@override`-required diagnostic).
//
// Fires JSC_INTERFACE_METHOD_NOT_IMPLEMENTED when a constructor
// annotated `@implements {I}` (or `@implements {I1} @implements {I2}`)
// omits a method declared on I (or on any interface I extends via
// `@extends`). The implementor's method set is gathered from its
// InstanceType's member map, walking the `@extends` super chain so
// an inherited implementation satisfies the interface slot.
//
// Two declaration shapes are recognised for the implementor:
//
//   /** @constructor @implements {Foo} */ function Bar() {};
//   /** @constructor @implements {Foo} */ var  Bar = function() {};
//
// Abstract classes (`@abstract`) are exempt — an abstract class may
// leave interface methods to concrete subclasses.

package jscomp

import (
	"fmt"
	"sort"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type CheckImplements struct {
	path       string
	creator    *TypedScopeCreator
	scope      *TypedScope
	scriptRoot *rhino.Node
	diags      []diag.Diagnostic
}

func NewCheckImplements(path string, creator *TypedScopeCreator, scope *TypedScope) *CheckImplements {
	return &CheckImplements{path: path, creator: creator, scope: scope}
}

func (c *CheckImplements) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.scriptRoot = script
	c.walk(script)
	return c.diags
}

func (c *CheckImplements) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.FUNCTION:
		c.checkFunction(n)
	case rhino.CLASS:
		c.checkClass(n)
	}
	for _, k := range n.Children() {
		c.walk(k)
	}
}

// checkFunction inspects a FUNCTION node that may be a `@constructor
// @implements {I}` declaration — either `function Foo(){}` or the
// value side of `var Foo = function(){}`.
func (c *CheckImplements) checkFunction(fn *rhino.Node) {
	info, name := constructorJSDocAndName(fn)
	if info == nil || !info.IsConstructor() {
		return
	}
	if len(info.Implementations()) == 0 {
		return
	}
	if name == "" {
		return
	}
	inst := c.scope.LookupNominalType(name)
	if inst == nil {
		return
	}
	c.verify(fn, info, inst, name)
}

// checkClass inspects an ES6 `class C implements I` (via JSDoc). Same
// logic as checkFunction but the name node is kids[0] on the CLASS.
func (c *CheckImplements) checkClass(cls *rhino.Node) {
	info := cls.JSDocInfo()
	if info == nil || len(info.Implementations()) == 0 {
		return
	}
	kids := cls.Children()
	if len(kids) < 1 {
		return
	}
	name := kids[0].String_()
	if name == "" {
		return
	}
	inst := c.scope.LookupNominalType(name)
	if inst == nil {
		return
	}
	c.verify(cls, info, inst, name)
}

// verify walks each `@implements` target and emits
// JSC_INTERFACE_METHOD_NOT_IMPLEMENTED for any slot the implementor
// (instance + super chain) is missing. Abstract classes are exempt.
func (c *CheckImplements) verify(pos *rhino.Node, info *rhino.JSDocInfo, inst *jstype.ObjectType, name string) {
	if info.IsAbstract() {
		return
	}
	stubbed := c.stubDeclaredProps(name)
	for _, impl := range info.Implementations() {
		iface := c.scope.LookupNominalType(impl.TypeName)
		if iface == nil {
			// Unresolved — that's P2-C's externs-hole job.
			continue
		}
		c.emitMissing(pos, inst, name, iface, stubbed)
	}
}

// stubDeclaredProps returns the set of property names declared on the
// implementor via stub-form statements — `X.prototype.m;` with no
// right-hand side. Java's TypedScopeCreator treats this as a valid
// declaration of the member; our scope.go only wires up ASSIGN-shaped
// prototype writes, so we recover the stub names here so they don't
// trigger a false JSC_INTERFACE_METHOD_NOT_IMPLEMENTED.
func (c *CheckImplements) stubDeclaredProps(implementorName string) map[string]bool {
	out := map[string]bool{}
	collectStubDecls(c.scriptRoot, implementorName, out)
	return out
}

func collectStubDecls(n *rhino.Node, implementorName string, out map[string]bool) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		if first := n.FirstChild(); first != nil && first.Token() == rhino.GETPROP {
			kids := first.Children()
			if len(kids) == 2 {
				inner := kids[0]
				propName := kids[1].String_()
				if inner.Token() == rhino.GETPROP {
					innerKids := inner.Children()
					if len(innerKids) == 2 &&
						innerKids[0].Token() == rhino.NAME &&
						innerKids[0].String_() == implementorName &&
						innerKids[1].String_() == "prototype" {
						out[propName] = true
					}
				}
			}
		}
	}
	for _, k := range n.Children() {
		collectStubDecls(k, implementorName, out)
	}
}

// emitMissing walks the interface's @extends chain in BFS order and
// emits the first method the implementor is missing. Mirrors Java's
// validateImplements which fires one diagnostic per unimplemented
// slot, ordered by the first interface that declares it.
func (c *CheckImplements) emitMissing(pos *rhino.Node, inst *jstype.ObjectType, name string, iface *jstype.ObjectType, stubbed map[string]bool) {
	visited := map[*jstype.ObjectType]bool{}
	queue := []*jstype.ObjectType{iface}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil || visited[cur] {
			continue
		}
		visited[cur] = true
		for _, prop := range declaredMemberOrder(cur) {
			if inst.GetPropertyType(prop) != nil {
				continue
			}
			if stubbed[prop] {
				continue
			}
			line, col := functionKeywordPos(pos)
			c.diags = append(c.diags, diag.Diagnostic{
				Path:  c.path,
				Line:  line,
				Col:   col,
				Code:  diag.JSC_INTERFACE_METHOD_NOT_IMPLEMENTED,
				Level: diag.Error,
				Message: fmt.Sprintf("property %s on interface %s is not implemented by type %s",
					prop, cur.TypeName(), name),
			})
			return
		}
		queue = append(queue, c.extendedInterfaceTypes(cur)...)
		if super := cur.Super(); super != nil {
			queue = append(queue, super)
		}
	}
}

// functionKeywordPos returns the (line, col) of the `function` keyword
// (for a FUNCTION node) or `class` keyword (for a CLASS node) that
// owns `pos`. The parser records position only on the NAME child, so
// we recover the keyword position by subtracting the fixed prefix
// length — `function ` (9) when the function is named,
// `function` (8) when anonymous.
func functionKeywordPos(pos *rhino.Node) (int, int) {
	if pos == nil {
		return 0, 0
	}
	kids := pos.Children()
	if len(kids) < 1 {
		return pos.Lineno(), pos.Charno()
	}
	name := kids[0]
	if name.Lineno() == 0 && name.Charno() == 0 {
		return pos.Lineno(), pos.Charno()
	}
	// CLASS keyword is "class " (6 when named) / "class" (5 when anon).
	if pos.Token() == rhino.CLASS {
		if name.String_() == "" {
			return name.Lineno(), name.Charno() - 5
		}
		return name.Lineno(), name.Charno() - 6
	}
	// FUNCTION keyword: 9 when named, 8 when anonymous.
	if name.String_() == "" {
		return name.Lineno(), name.Charno() - 8
	}
	return name.Lineno(), name.Charno() - 9
}

// extendedInterfaceTypes resolves every `@extends {…}` expression that
// appeared on the interface's declaration JSDoc to its instance type.
// Needed because interfaces may carry multiple `@extends`, but
// ObjectType.Super() only stores the last one wired by scope.go.
func (c *CheckImplements) extendedInterfaceTypes(iface *jstype.ObjectType) []*jstype.ObjectType {
	info := c.findInterfaceJSDoc(iface.Name)
	if info == nil {
		return nil
	}
	var out []*jstype.ObjectType
	for _, e := range info.ExtendedInterfaces() {
		if t := c.scope.LookupNominalType(e.TypeName); t != nil {
			out = append(out, t)
		}
	}
	return out
}

// findInterfaceJSDoc recovers the JSDocInfo attached to an interface
// declaration named `name` by scanning TypedScopeCreator's per-function
// map. The scope tree doesn't expose a FunctionType → JSDoc
// back-reference, so we match by FUNCTION-node name.
func (c *CheckImplements) findInterfaceJSDoc(name string) *rhino.JSDocInfo {
	for fn := range c.creator.Functions {
		if fn == nil || fn.Token() != rhino.FUNCTION {
			continue
		}
		kids := fn.Children()
		if len(kids) < 1 || kids[0].String_() != name {
			continue
		}
		if info := fn.JSDocInfo(); info != nil && info.IsInterface() {
			return info
		}
		// `var Foo = function(){}` — JSDoc lives on the VAR/LET/CONST.
		if p := fn.Parent(); p != nil && p.Token() == rhino.NAME && p.Parent() != nil {
			if d := p.Parent(); isDeclToken(d.Token()) {
				if di := d.JSDocInfo(); di != nil && di.IsInterface() {
					return di
				}
			}
		}
	}
	return nil
}

// declaredMemberOrder returns directly-declared member names sorted
// lexicographically. Go map iteration is non-deterministic; Java's
// LinkedHashMap preserves insertion order, but for a single-miss
// diagnostic deterministic lexicographic selection is good enough.
func declaredMemberOrder(o *jstype.ObjectType) []string {
	m := o.OwnMembers()
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// constructorJSDocAndName returns the (JSDocInfo, name) pair for a
// FUNCTION node when the function is part of a `@constructor` shape
// this pass recognises. Supports:
//
//   /** @constructor */ function Foo() { … }        // FUNCTION owns the JSDoc
//   /** @constructor */ var Foo = function() { … }; // VAR owns the JSDoc
//
// Returns (nil, "") when the node isn't a recognised constructor shape.
func constructorJSDocAndName(fn *rhino.Node) (*rhino.JSDocInfo, string) {
	// Case 1: `function Foo(){}` — JSDoc on the FUNCTION, name in kids[0].
	kids := fn.Children()
	if len(kids) < 1 {
		return nil, ""
	}
	if info := fn.JSDocInfo(); info != nil && kids[0].String_() != "" {
		return info, kids[0].String_()
	}
	// Case 2: `var Foo = function(){};` — the FUNCTION's parent is a
	// NAME ("Foo") whose parent is a VAR/LET/CONST carrying JSDoc.
	p := fn.Parent()
	if p == nil || p.Token() != rhino.NAME {
		return nil, ""
	}
	g := p.Parent()
	if g == nil || !isDeclToken(g.Token()) {
		return nil, ""
	}
	info := g.JSDocInfo()
	if info == nil {
		info = p.JSDocInfo()
	}
	if info == nil {
		info = fn.JSDocInfo()
	}
	return info, p.String_()
}
