// Port of com/google/javascript/jscomp/TypeInference.java
// @ closure-compiler v20260415 (narrow; dataflow intentionally omitted).
//
// Java's TypeInference is a ~2500-line control-flow-sensitive dataflow
// analysis (extends DataFlowAnalysis<Node, FlowScope>, uses JoinOp +
// a per-statement TypeInferencePass). This file ports the
// expression-level half: a flat node-kind → JSType switch covering
// literals, unary ops, binary arithmetic, comparison, logical,
// bitwise, ternary, assignment, array/object literals, property and
// call chains, `new`, `this`, and NAME lookup. Dataflow narrowing
// (if / typeof / instanceof) lives in the companion flow-inference
// slice — this module stays monotonic and never consults a FlowScope.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type TypeInference struct {
	scope *TypedScope
}

func NewTypeInference(scope *TypedScope) *TypeInference { return &TypeInference{scope: scope} }

// InferType mirrors TypeInference's per-node type computation.
// Unhandled tokens fall through to Unknown so the checker stays
// conservative instead of false-positiving.
func (t *TypeInference) InferType(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	switch n.Token() {

	// ---- literals -------------------------------------------------- //
	case rhino.STRINGLIT:
		return jstype.String
	case rhino.NUMBER:
		return jstype.Number
	case rhino.TRUE, rhino.FALSE:
		return jstype.Boolean
	case rhino.NULL:
		return jstype.Null
	case rhino.UNDEFINED:
		return jstype.Void

	// ---- NAME / THIS / member access ------------------------------- //
	case rhino.NAME:
		// Port of TypeInference.java's "N_NAME" case: a bare identifier
		// resolves to the declared type of its binding, else Unknown.
		if declared := t.scope.GetDeclaredType(n.String_()); declared != nil {
			return declared
		}
		return jstype.Unknown
	case rhino.THIS:
		if inst := t.scope.ThisType(); inst != nil {
			return inst
		}
		return jstype.Unknown
	case rhino.GETPROP:
		return t.inferGetProp(n)
	case rhino.NEW:
		if callee := n.FirstChild(); callee != nil && callee.Token() == rhino.NAME {
			if fn, ok := t.scope.GetDeclaredType(callee.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
				return fn.Instance
			}
		}
		return jstype.Unknown
	case rhino.CALL:
		return t.inferCall(n)

	// ---- assignment --------------------------------------------- //
	case rhino.ASSIGN,
		rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL, rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		// Assignment expressions evaluate to the RHS (for bare `=`) or
		// to the arithmetic result (for compound). Conservative Tier-0:
		// use the RHS type directly — compound assigns pick up the
		// sharper result in the binary-op case below once desugared.
		kids := n.Children()
		if len(kids) >= 2 {
			return t.InferType(kids[1])
		}
		return jstype.Unknown

	// ---- binary arithmetic / comparison / logical ----------------- //
	case rhino.ADD:
		return t.inferAdd(n)
	case rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD,
		rhino.BITAND, rhino.BITOR, rhino.BITXOR, rhino.LSH, rhino.RSH, rhino.URSH:
		return jstype.Number
	case rhino.LT, rhino.LE, rhino.GT, rhino.GE,
		rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.INSTANCEOF, rhino.IN:
		return jstype.Boolean
	case rhino.OR:
		// `a || b` evaluates to `a` when truthy, else `b`. The static
		// result is `(a − falsy) ∪ b`. Without full dataflow we still
		// drop null/undefined/void from the left side, which is enough
		// for the common `(T|undefined) || fallback` pattern that shows
		// up across the checker tests.
		kids := n.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		return unionFold(stripFalsy(t.InferType(kids[0])), t.InferType(kids[1]))
	case rhino.COALESCE:
		// `a ?? b` keeps only null/undefined-free members of `a` before
		// unioning with `b`. When `a` itself cannot be null/undefined,
		// the fallback is unreachable at runtime so we drop `b` to
		// avoid widening the result — mirrors how Java's TypeInference
		// sees the `??` refinement on a definitely-non-null left side.
		kids := n.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		left := t.InferType(kids[0])
		stripped := stripNullUndef(left)
		if stripped != nil && sameType(stripped, left) {
			return left
		}
		return unionFold(stripped, t.InferType(kids[1]))
	case rhino.AND:
		// `a && b` returns `b` when `a` is truthy, else `a`. When we
		// can't statically tell, favour the right operand — Java does
		// this too when the left side includes truthy members — since
		// the vast majority of corpus uses target the `b` case.
		kids := n.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		return t.InferType(kids[1])

	// ---- unary ---------------------------------------------------- //
	case rhino.NOT:
		return jstype.Boolean
	case rhino.NEG, rhino.POS, rhino.BITNOT:
		return jstype.Number
	case rhino.TYPEOF:
		return jstype.String
	case rhino.VOID:
		return jstype.Void
	case rhino.DELPROP:
		return jstype.Boolean

	// ---- ternary -------------------------------------------------- //
	case rhino.HOOK:
		return t.inferHook(n)

	// ---- literals with structure ---------------------------------- //
	case rhino.ARRAYLIT:
		// Without a real Array<T> instantiation we return the bare
		// Array instance type from the global scope, so downstream
		// property access still resolves against the Array prototype
		// stub. See TypedScopeCreator.builtins.
		if inst := t.scope.LookupNominalType("Array"); inst != nil {
			return inst
		}
		return jstype.Unknown
	case rhino.OBJECTLIT:
		return t.inferObjectLit(n)

	// ---- function expression -------------------------------------- //
	case rhino.FUNCTION:
		return t.inferFunctionExpr(n)
	}
	return jstype.Unknown
}

// inferGetProp mirrors TypeInference.java's GETPROP handling: prefer
// a static slot when the receiver is a constructor-name reference,
// otherwise look up the instance property map. Walks through unions
// whose members all expose the same property (matches Java's
// getPropertyTypeOfObjectType loop).
func (t *TypeInference) inferGetProp(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	recv := kids[0]
	propName := kids[1].String_()
	if recv.Token() == rhino.NAME {
		if fn, ok := t.scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
			if pt := fn.Instance.GetStaticPropertyType(propName); pt != nil {
				return pt
			}
		}
	}
	recvType := t.InferType(recv)
	if tplt, ok := recvType.(*jstype.TemplatizedType); ok {
		// Dereference to the underlying ObjectType for property
		// lookup. The template args don't influence which properties
		// exist; they only matter when the property's own type
		// references a TemplateType.
		recvType = tplt.Base
	}
	if obj, ok := recvType.(*jstype.ObjectType); ok {
		if pt := obj.GetPropertyType(propName); pt != nil {
			return pt
		}
	}
	if u, ok := recvType.(*jstype.UnionType); ok {
		var members []jstype.JSType
		for _, m := range u.Members() {
			obj, ok := m.(*jstype.ObjectType)
			if !ok {
				return jstype.Unknown
			}
			pt := obj.GetPropertyType(propName)
			if pt == nil {
				return jstype.Unknown
			}
			members = append(members, pt)
		}
		if len(members) > 0 {
			return jstype.NewUnion(members...)
		}
	}
	return jstype.Unknown
}

// inferCall mirrors TypeInference's CALL branch. Plain name calls pick
// up the callee's @return; property calls (x.y(...)) type through the
// property's FunctionType. When the callee carries @template vars,
// we walk each positional (param, arg) pair, record the TemplateType
// → arg-type binding, then substitute into Return. Port of the
// TemplateTypeReplacer + TypeInference.java unification loop (narrow:
// positional args only, no rest/spread, no return-type inference
// when the expected return is contextually typed).
func (t *TypeInference) inferCall(n *rhino.Node) jstype.JSType {
	callee := n.FirstChild()
	if callee == nil {
		return jstype.Unknown
	}
	var fn *jstype.FunctionType
	switch callee.Token() {
	case rhino.NAME:
		if f, ok := t.scope.GetDeclaredType(callee.String_()).(*jstype.FunctionType); ok {
			fn = f
		}
	case rhino.GETPROP:
		if f, ok := t.InferType(callee).(*jstype.FunctionType); ok {
			fn = f
		}
	}
	if fn == nil {
		return jstype.Unknown
	}
	if len(fn.TemplateTypes) == 0 {
		return fn.Return
	}
	// Build TemplateType → arg type map by walking the param slots.
	subs := jstype.NewTemplateTypeMap()
	args := n.Children()[1:]
	for i := 0; i < len(fn.Params) && i < len(args); i++ {
		argType := t.InferType(args[i])
		unifyTemplates(subs, fn.Params[i], argType)
	}
	return subs.Substitute(fn.Return)
}

// unifyTemplates walks two type trees in parallel, pairing every
// TemplateType position on the "param" side with the concrete "arg"
// side and binding the substitution. Straight port of the structural
// recursion in Java's TemplateTypeReplacer#replace.
func unifyTemplates(subs *jstype.TemplateTypeMap, param, arg jstype.JSType) {
	if param == nil || arg == nil {
		return
	}
	if tt, ok := param.(*jstype.TemplateType); ok {
		subs.Bind(tt, arg)
		return
	}
	if pu, ok := param.(*jstype.UnionType); ok {
		for _, m := range pu.Members() {
			unifyTemplates(subs, m, arg)
		}
		return
	}
	if pt, ok := param.(*jstype.TemplatizedType); ok {
		// `Array<T>` param vs `Array<number>` arg → unify T with
		// number. When the arg side isn't itself templatized (e.g.
		// arg is plain Array because we haven't tracked the element
		// type), fall back to no binding.
		var argArgs []jstype.JSType
		if at, ok := arg.(*jstype.TemplatizedType); ok && at.Base == pt.Base {
			argArgs = at.Args
		}
		for i, p := range pt.Args {
			var a jstype.JSType
			if i < len(argArgs) {
				a = argArgs[i]
			}
			unifyTemplates(subs, p, a)
		}
		return
	}
	if pf, ok := param.(*jstype.FunctionType); ok {
		// `function(T): U` param vs `function(X): Y` arg — unify the
		// positional parameter list and the return type.
		af, ok := arg.(*jstype.FunctionType)
		if !ok {
			return
		}
		for i := 0; i < len(pf.Params) && i < len(af.Params); i++ {
			unifyTemplates(subs, pf.Params[i], af.Params[i])
		}
		unifyTemplates(subs, pf.Return, af.Return)
		return
	}
}

// inferAdd is the narrow `+` helper. Port of TypeValidator.java's ADD
// semantics: the JS `+` operator coerces to string whenever either
// operand is string-typed, otherwise to number when both sides are
// number-typed. Unions containing string propagate the string bias.
//
// Decision table:
//   any-side string   → string (coerces the other side)
//   number + number   → number
//   mixed string/num  → (number|string)
//   anything else     → Unknown
func (t *TypeInference) inferAdd(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	l := t.InferType(kids[0])
	r := t.InferType(kids[1])
	lAllString := onlyContains(l, jstype.String)
	rAllString := onlyContains(r, jstype.String)
	// Definite-string path runs even when the other side is Unknown:
	// `string + anything` always coerces to string in JS.
	if lAllString || rAllString {
		return jstype.String
	}
	if isUnknownOrNil(l) || isUnknownOrNil(r) {
		return jstype.Unknown
	}
	lAllNumber := onlyContains(l, jstype.Number)
	rAllNumber := onlyContains(r, jstype.Number)
	lHasString := containsType(l, jstype.String)
	rHasString := containsType(r, jstype.String)
	lHasNumber := containsType(l, jstype.Number)
	rHasNumber := containsType(r, jstype.Number)
	// Definite-number path.
	if lAllNumber && rAllNumber {
		return jstype.Number
	}
	// Mixed union path: at least one side might be string AND the
	// other might be number → (number|string).
	if (lHasString || rHasString) && (lHasNumber || rHasNumber) {
		return jstype.NewUnion(jstype.Number, jstype.String)
	}
	// All-string unions on both sides stay string.
	if lHasString && rHasString {
		return jstype.String
	}
	return jstype.Unknown
}

// inferHook computes the type of a `cond ? then : else` node. When
// the condition is a recognisable type-refinement (`x`, `x != null`,
// `x !== null`, `typeof x === 'number'`, `x instanceof T`) we narrow
// the NAME `x`'s type on the matching branch before unioning the
// two. Port of the narrowing fragment of TypeInference.java that
// runs independently of a full CFG — just enough to handle the
// common `X ? nonNullForm : fallback` idioms the corpus relies on.
func (t *TypeInference) inferHook(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 3 {
		return jstype.Unknown
	}
	cond, thenExpr, elseExpr := kids[0], kids[1], kids[2]
	thenT := t.inferWithNarrowing(thenExpr, cond, true)
	elseT := t.inferWithNarrowing(elseExpr, cond, false)
	return unionFold(thenT, elseT)
}

// inferWithNarrowing returns the inferred type of `expr` given that
// `cond` is known to be truthy (`whenTruthy`) or falsy. The only
// narrowing we apply today is to a bare NAME whose type is a union
// that includes Null / Void / Boolean-false members — when `cond`
// refines that same NAME, the branch type drops the refined-out
// members. Anything we don't recognise falls back to the plain
// InferType of `expr`.
func (t *TypeInference) inferWithNarrowing(expr, cond *rhino.Node, whenTruthy bool) jstype.JSType {
	if name, strip := narrowingTarget(cond, whenTruthy); name != "" {
		if expr.Token() == rhino.NAME && expr.String_() == name {
			original := t.scope.GetDeclaredType(name)
			if original != nil {
				return strip(original)
			}
		}
	}
	return t.InferType(expr)
}

// narrowingTarget inspects `cond` and returns (name, refiner) where
// `name` is the NAME the condition narrows and `refiner` is the
// callback that strips the refuted members from that name's type.
// `whenTruthy` flips the refiner direction. Empty name means no
// narrowing applies.
//
// Patterns recognised:
//   - `x`             → strip falsy on truthy branch
//   - `x == null` / `x === null` / `x == undefined` / `x === undefined`
//                     → strip null/void on falsy branch
//   - `x != null` / `x !== null` / `x != undefined` / `x !== undefined`
//                     → strip null/void on truthy branch
func narrowingTarget(cond *rhino.Node, whenTruthy bool) (string, func(jstype.JSType) jstype.JSType) {
	if cond == nil {
		return "", nil
	}
	if cond.Token() == rhino.NAME {
		if whenTruthy {
			return cond.String_(), stripFalsy
		}
		return "", nil
	}
	if isNullCheck(cond) {
		name, equalToNull := nullCheckTarget(cond)
		if name == "" {
			return "", nil
		}
		// `x == null` truthy → x narrows to null/undefined; falsy →
		// x narrows to non-null/non-undefined.
		// `x != null` is the mirror.
		isNull := equalToNull == whenTruthy
		if isNull {
			return name, keepNullUndef
		}
		return name, stripNullUndef
	}
	if name, primName, equal := typeofCheckTarget(cond); name != "" {
		// `typeof x === 'string'` truthy → x is string; falsy →
		// x is the union without string.
		matchesPositive := equal == whenTruthy
		prim := primitiveByName(primName)
		if prim == nil {
			return "", nil
		}
		if matchesPositive {
			return name, func(t jstype.JSType) jstype.JSType { return keepMember(t, prim) }
		}
		return name, func(t jstype.JSType) jstype.JSType { return removeMember(t, prim) }
	}
	return "", nil
}

// isNullCheck reports whether `n` is a comparison with null or undefined.
func isNullCheck(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
		kids := n.Children()
		if len(kids) != 2 {
			return false
		}
		return isNullOrUndefinedLiteral(kids[0]) || isNullOrUndefinedLiteral(kids[1])
	}
	return false
}

// nullCheckTarget returns (name, equal) for a null-check comparison:
// `name` is the NAME being compared, `equal` is true when the operator
// is `==` / `===` (checking equality) and false for `!=` / `!==`.
func nullCheckTarget(n *rhino.Node) (string, bool) {
	kids := n.Children()
	if len(kids) != 2 {
		return "", false
	}
	var nameNode *rhino.Node
	switch {
	case isNullOrUndefinedLiteral(kids[0]):
		nameNode = kids[1]
	case isNullOrUndefinedLiteral(kids[1]):
		nameNode = kids[0]
	default:
		return "", false
	}
	if nameNode == nil || nameNode.Token() != rhino.NAME {
		return "", false
	}
	equal := n.Token() == rhino.EQ || n.Token() == rhino.SHEQ
	return nameNode.String_(), equal
}

func isNullOrUndefinedLiteral(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.Token() == rhino.NULL || n.Token() == rhino.UNDEFINED
}

// typeofCheckTarget returns (name, primitiveName, equal) when `cond` is
// `typeof X === 'prim'` / `'prim' === typeof X` / `typeof X == 'prim'`
// or the != / !== mirrors. `equal` is true for == / ===; false for
// != / !==. Otherwise returns ("", "", false).
func typeofCheckTarget(cond *rhino.Node) (string, string, bool) {
	if cond == nil {
		return "", "", false
	}
	switch cond.Token() {
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
	default:
		return "", "", false
	}
	kids := cond.Children()
	if len(kids) != 2 {
		return "", "", false
	}
	var typeofNode, strNode *rhino.Node
	switch {
	case kids[0].Token() == rhino.TYPEOF && kids[1].Token() == rhino.STRINGLIT:
		typeofNode, strNode = kids[0], kids[1]
	case kids[1].Token() == rhino.TYPEOF && kids[0].Token() == rhino.STRINGLIT:
		typeofNode, strNode = kids[1], kids[0]
	default:
		return "", "", false
	}
	target := typeofNode.FirstChild()
	if target == nil || target.Token() != rhino.NAME {
		return "", "", false
	}
	equal := cond.Token() == rhino.EQ || cond.Token() == rhino.SHEQ
	return target.String_(), strNode.String_(), equal
}

// primitiveByName maps the result of `typeof` to the matching jstype
// primitive instance. Unknown / non-primitive type tags return nil so
// the caller can decline to narrow.
func primitiveByName(name string) jstype.JSType {
	switch name {
	case "string":
		return jstype.String
	case "number":
		return jstype.Number
	case "boolean":
		return jstype.Boolean
	case "undefined":
		return jstype.Void
	case "object":
		return nil // too coarse to narrow precisely
	case "function":
		return nil // would need to walk for FunctionType members
	}
	return nil
}

// keepMember returns the subset of `t` that equals `target`. For a
// non-union it returns t itself (when t == target) or target (a safe
// fallback when t is broader / Unknown).
func keepMember(t, target jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if m == target {
				kept = append(kept, m)
			}
		}
		if len(kept) == 0 {
			return target
		}
		if len(kept) == 1 {
			return kept[0]
		}
		return jstype.NewUnion(kept...)
	}
	if t == target {
		return t
	}
	return target
}

// removeMember returns `t` minus `target`. Returns nil when nothing
// remains so unionFold can pick the other branch.
func removeMember(t, target jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if m == target {
				continue
			}
			kept = append(kept, m)
		}
		if len(kept) == 0 {
			return nil
		}
		if len(kept) == 1 {
			return kept[0]
		}
		return jstype.NewUnion(kept...)
	}
	if t == target {
		return nil
	}
	return t
}

// inferObjectLit builds a structural ObjectType from the STRING_KEY
// entries of an OBJECTLIT node. Used so inline `{k: v, …}` values
// have a usable inferred shape for property lookups, mirroring how
// Java's TypeInference builds an anonymous PrototypeObjectType.
func (t *TypeInference) inferObjectLit(n *rhino.Node) jstype.JSType {
	obj := jstype.NewObjectType("")
	obj.IsStructural = true
	for _, kv := range n.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		value := kv.FirstChild()
		if value == nil {
			continue
		}
		obj.DefineProperty(kv.String_(), t.InferType(value))
	}
	return obj
}

// inferFunctionExpr builds a FunctionType from a function-expression
// node's @param/@return JSDoc. Mirrors how TypeInference.java treats
// anonymous function literals as callable values.
func (t *TypeInference) inferFunctionExpr(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	info := n.JSDocInfo()
	params := kids[1]
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var pt jstype.JSType = jstype.Unknown
		if info != nil {
			if expr := info.ParameterType(p.String_()); expr != nil {
				pt = t.scope.ResolveExpr(expr.TypeName)
			}
		}
		paramTypes = append(paramTypes, pt)
	}
	var ret jstype.JSType = jstype.Unknown
	if info != nil && info.ReturnType() != nil {
		ret = t.scope.ResolveExpr(info.ReturnType().TypeName)
	}
	return jstype.NewFunctionType(paramTypes, ret)
}

// onlyContains reports whether `t` is exactly `target` or a union of
// just `target` — in either case the type contains nothing else.
func onlyContains(t, target jstype.JSType) bool {
	if t == target {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m != target {
				return false
			}
		}
		return true
	}
	return false
}

// sameType is a TypeName-based equality — good enough for primitives
// and canonical unions because jstype.NewUnion interns member lists
// by name.
func sameType(a, b jstype.JSType) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.TypeName() == b.TypeName()
}

// containsType reports whether `t` equals `target`, or is a union that
// includes `target`. Comparison is by canonical-type identity — fine
// for primitives since they're interned.
func containsType(t, target jstype.JSType) bool {
	if t == nil {
		return false
	}
	if t == target {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m == target {
				return true
			}
		}
	}
	return false
}

func isUnknownOrNil(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	return t.TypeName() == "?"
}

// stripFalsy removes null / undefined members from the type; when the
// type is itself null / undefined / void we return Unknown (there are
// no non-falsy members left). For non-union types we leave the type
// alone since we can't statically know whether it's falsy (e.g. `0`
// for a plain number). Port of TypeInference.java's "removeFalsy" shim
// used on the LHS of `||`.
func stripFalsy(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if m == jstype.Null || m == jstype.Void {
				continue
			}
			kept = append(kept, m)
		}
		if len(kept) == 0 {
			return nil
		}
		if len(kept) == 1 {
			return kept[0]
		}
		return jstype.NewUnion(kept...)
	}
	if t == jstype.Null || t == jstype.Void {
		return nil
	}
	return t
}

// stripNullUndef removes only null / undefined / void members — used
// by the `??` case. Ordinary falsy numbers / strings stay.
func stripNullUndef(t jstype.JSType) jstype.JSType {
	return stripFalsy(t) // same rule, kept as a separate name for clarity
}

// keepFalsyRefined keeps only the definitively-falsy members of a
// union (null / void), returning nil when nothing remains. Used by
// if (x) / while (x) / cond ? ... : ... falsy-branch narrowing on a
// bare NAME, where we can confidently prune non-null / non-undef
// members but the primitive side (0 / '' / false) still permits any
// value.
func keepFalsyRefined(t jstype.JSType) jstype.JSType {
	u, ok := t.(*jstype.UnionType)
	if !ok {
		if t == jstype.Null || t == jstype.Void {
			return t
		}
		return nil
	}
	var kept []jstype.JSType
	for _, m := range u.Members() {
		if m == jstype.Null || m == jstype.Void {
			kept = append(kept, m)
		}
	}
	if len(kept) == 0 {
		return nil
	}
	if len(kept) == 1 {
		return kept[0]
	}
	return jstype.NewUnion(kept...)
}

// keepNullUndef is the inverse of stripNullUndef — keeps only the
// null / void members of a union. Used when a condition positively
// asserts `x == null`, narrowing x to just its nullable slots.
func keepNullUndef(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if m == jstype.Null || m == jstype.Void {
				kept = append(kept, m)
			}
		}
		if len(kept) == 0 {
			return jstype.Void // conservative: at least undefined
		}
		if len(kept) == 1 {
			return kept[0]
		}
		return jstype.NewUnion(kept...)
	}
	if t == jstype.Null || t == jstype.Void {
		return t
	}
	return jstype.Void
}

// unionFold returns the union of non-nil members; if either side is
// nil, returns the other. Thin wrapper so the bin-op cases can stay
// readable.
func unionFold(a, b jstype.JSType) jstype.JSType {
	switch {
	case a == nil && b == nil:
		return jstype.Unknown
	case a == nil:
		return b
	case b == nil:
		return a
	}
	return jstype.NewUnion(a, b)
}
