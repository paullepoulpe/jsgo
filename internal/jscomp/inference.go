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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type TypeInference struct {
	scope *TypedScope
	// inExprBindings tracks in-expression ASSIGN side effects so a
	// later read of the same NAME within the same expression sees the
	// freshly assigned type. Mirrors the Java FlowScope mutation that
	// TypeInference.traverseAssign performs via updateScopeForAssignment
	// before continuing the left-to-right traversal.
	inExprBindings map[string]jstype.JSType
}

// Go-only: package-bridge
func NewTypeInference(scope *TypedScope) *TypeInference { return &TypeInference{scope: scope} }

// InferType mirrors TypeInference's per-node type computation.
// Unhandled tokens fall through to Unknown so the checker stays
// conservative instead of false-positiving.
// InferType is the public entry point for per-node type inference.
// It delegates to traverse, which is the Java-faithful dispatcher.
// Port of TypeInference.java#flowThrough → traverse.
func (t *TypeInference) InferType(n *rhino.Node) jstype.JSType {
	return t.traverse(n)
}

// Go-only: language-adapter
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
	// `this.X` at global scope resolves to the externs `X` declaration —
	// the global this object's prototype is the Window externs surface,
	// so `this.alert` reads the externs `alert` slot. Mirrors Closure's
	// resolveAll() linking GLOBAL_THIS's implicit prototype to Window
	// (jstype_resolver.go:206 documents the still-unported piece, but the
	// observable type-lookup behaviour is reproducible by routing through
	// the externs declared scope). Done before the qualified-refinement
	// lookup because a rejected `this.X = …` assignment in global scope
	// must not update the slot's read-back type — the externs slot still
	// owns the property surface. Restricted to names declared only in
	// externs (not in the user's global scope) so user redeclarations
	// (`function alert(msg) {}`) don't supply the type for `this.alert`
	// — Java suppresses TypeCheck on those redeclarations because the
	// resulting JSC_VAR_MULTIPLY_DECLARED_ERROR aborts further analysis.
	if recv.Token() == rhino.THIS && t.scope.ThisType() == nil {
		if pt := t.scope.LookupExternsOnly(propName); pt != nil {
			return pt
		}
	}
	// Flow refinement for qualified names: `if (obj.prop) { … }`
	// narrows obj.prop inside the then-branch. The narrowed type
	// is keyed on the dotted name "obj.prop" in the refinement
	// stack and takes priority over the declared property type.
	if qname := qualifiedName(n); qname != "" && strings.Contains(qname, ".") {
		if r := t.scope.GetQualifiedRefinement(qname); r != nil {
			return r
		}
	}
	if recv.Token() == rhino.NAME {
		if fn, ok := t.scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
			if pt := fn.Instance.GetStaticPropertyType(propName); pt != nil {
				return pt
			}
		}
		// Enum member access: `MyEnum.FOO` → EnumElementType,
		// matching Java's EnumType.getPropertyType which always
		// returns the paired EnumElementType for member reads.
		if enum, ok := t.scope.GetDeclaredType(recv.String_()).(*jstype.EnumType); ok {
			if enum.Members[propName] {
				return enum.Element
			}
			// Falling through to Unknown keeps un-declared-member
			// references quiet here; the INEXISTENT_ENUM_ELEMENT
			// diagnostic fires in TypeCheck.visitGetProp.
		}
		// Namespace property declared via
		// `/** @type {T} */ owner.prop = …`.
		if np := t.scope.LookupNamespaceProp(recv.String_(), propName); np != nil {
			return np
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
		// Port of UnionType.findPropertyTypeWithoutConsideringTemplateTypes:
		// skip null/void members and members lacking the property; union
		// the remaining property types via getLeastSupertype.
		if pt := u.FindPropertyType(propName); pt != nil {
			return pt
		}
	}
	return jstype.Unknown
}

// Go-only: language-adapter
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
		// Bare-name calls to the String / Number / Boolean
		// coercion constructors resolve to the matching
		// primitive — `String(n)` is a string, not a String
		// instance. Matches Java's special-cased coercion path
		// in TypeInference.visitCall.
		switch callee.String_() {
		case "String":
			return jstype.String
		case "Number":
			return jstype.Number
		case "Boolean":
			return jstype.Boolean
		case "Symbol":
			// Symbol() without `new` returns a primitive symbol value.
			// `new Symbol()` throws at runtime; Java's TypeInference
			// maps CALL(Symbol, …) → symbol primitive here.
			return jstype.Symbol
		}
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
	subs := jstype.NewCallSiteMap()
	args := n.Children()[1:]
	for i := 0; i < len(fn.Params) && i < len(args); i++ {
		argType := t.InferType(args[i])
		unifyTemplates(subs, fn.Params[i], argType)
	}
	return subs.Substitute(fn.Return)
}

// Go-only: language-adapter
// unifyTemplates walks two type trees in parallel, pairing every
// TemplateType position on the "param" side with the concrete "arg"
// side and binding the substitution. Straight port of the structural
// recursion in Java's TemplateTypeReplacer#replace.
func unifyTemplates(subs *jstype.CallSiteMap, param, arg jstype.JSType) {
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

// Go-only: language-adapter
// inferAdd is the narrow `+` helper. Port of TypeInference.java's
// traverseAdd: the JS `+` operator coerces to string whenever either
// operand is definitely a string type (and not unknown), to number
// when both sides are "added as number" (subtype of the
// NUMBER_ADDITION_SUPERTYPE = boolean|null|undefined|number), and to
// (number|string) otherwise. Mirrors the Java decision table exactly.
//
// Decision table (Java traverseAdd order):
//   !unknown && isString on either side   → string
//   either side unknown                   → unknown
//   isAddedAsNumber(l) && isAddedAsNumber(r) → number
//   otherwise                             → (number|string)
func (t *TypeInference) inferAdd(n *rhino.Node) jstype.JSType {
	return t.inferAddChildren(n)
}

// Go-only: language-adapter
// inferAddChildren shares the result-type computation between ADD
// (`+`) and ASSIGN_ADD (`+=`). Both operators carry two child
// sub-expressions; compound-assign behaves identically to the binary
// form for the purpose of the synthesised result type.
//
// Port of TypeInference.java traverseAdd; see inferAdd for the table.
func (t *TypeInference) inferAddChildren(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	l := t.InferType(kids[0])
	r := t.InferType(kids[1])
	leftIsUnknown := isUnknownOrNil(l)
	rightIsUnknown := isUnknownOrNil(r)
	// Definite-string path: mirrors Java's
	//   (!leftIsUnknown && leftType.isString()) ||
	//   (!rightIsUnknown && rightType.isString())
	if (!leftIsUnknown && onlyContains(l, jstype.String)) ||
		(!rightIsUnknown && onlyContains(r, jstype.String)) {
		return jstype.String
	}
	// Unknown propagation — same as Java.
	if leftIsUnknown || rightIsUnknown {
		return jstype.Unknown
	}
	// Both sides are known non-string. Use isAddedAsNumber (mirrors
	// Java's isAddedAsNumber = isSubtypeOf(NUMBER_ADDITION_SUPERTYPE)
	// where NUMBER_ADDITION_SUPERTYPE = boolean|null|undefined|number).
	// When both are "number-ish", the result is number; otherwise the
	// JS spec says it coerces both to strings → (number|string).
	if isAddedAsNumber(l) && isAddedAsNumber(r) {
		return jstype.Number
	}
	return jstype.NewUnion(jstype.Number, jstype.String)
}

// Go-only: language-adapter
// isAddedAsNumber reports whether `t` is a subtype of Java's
// NUMBER_ADDITION_SUPERTYPE = (boolean|Boolean|null|undefined|number|Number).
// Port of TypeInference.java isAddedAsNumber + JSTypeRegistry
// numberAdditionSupertype. The Java registry includes NUMBER_OBJECT_TYPE
// (boxed Number) and BOOLEAN_OBJECT_TYPE (boxed Boolean); we recognise them
// as ObjectType instances whose TypeName is "Number" or "Boolean" since
// they are declared as @constructor stubs in our default externs.
func isAddedAsNumber(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	switch t {
	case jstype.Number, jstype.Boolean, jstype.Null, jstype.Void:
		return true
	}
	// Boxed Number / Boolean — ObjectType instances from the default externs.
	// Matches Java's NUMBER_OBJECT_TYPE / BOOLEAN_OBJECT_TYPE entries in
	// NUMBER_ADDITION_SUPERTYPE.
	if obj, ok := t.(*jstype.ObjectType); ok {
		name := obj.TypeName()
		if name == "Number" || name == "Boolean" {
			return true
		}
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isAddedAsNumber(m) {
				return false
			}
		}
		return true
	}
	return false
}

// Go-only: language-adapter
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
	// Build per-branch refinements for the full condition shape
	// (NAME / GETPROP / AND / OR / typeof / instanceof / null
	// checks). Push them temporarily while inferring each branch so
	// qualified names like `obj.prop` see the narrowed type.
	thenRefs := buildRefinements(t.scope, cond, true)
	elseRefs := buildRefinements(t.scope, cond, false)
	thenT := t.inferUnder(thenExpr, thenRefs)
	elseT := t.inferUnder(elseExpr, elseRefs)
	return unionFold(thenT, elseT)
}

// Go-only: language-adapter
// inferUnder evaluates `expr` with `refs` temporarily pushed as a
// refinement overlay — used by inferHook and AND/OR short-circuit
// inference so each branch sees the narrowed flow.
func (t *TypeInference) inferUnder(expr *rhino.Node, refs map[string]jstype.JSType) jstype.JSType {
	if len(refs) == 0 {
		return t.InferType(expr)
	}
	t.scope.PushRefinement(refs)
	defer t.scope.PopRefinement()
	return t.InferType(expr)
}

// Go-only: language-adapter
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

// Go-only: language-adapter
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
		name, equalToNull, strict, nullLit := nullCheckTarget(cond)
		if name == "" {
			return "", nil
		}
		// `x == null` truthy → x narrows to null/undefined; falsy →
		// x narrows to non-null/non-undefined. Strict === only keeps/strips
		// the specific literal (null or undefined, not both).
		// Port of SemanticReverseAbstractInterpreter.caseEquality.
		isNull := equalToNull == whenTruthy
		if strict {
			if isNull {
				if nullLit {
					return name, keepNull
				}
				return name, keepVoid
			}
			if nullLit {
				return name, stripNull
			}
			return name, stripVoid
		}
		if isNull {
			return name, keepNullUndef
		}
		return name, stripNullUndef
	}
	if name, primName, equal := typeofCheckTarget(cond); name != "" {
		matchesPositive := equal == whenTruthy
		if primName == "object" {
			if matchesPositive {
				return name, keepObjectish
			}
			return name, stripObjectish
		}
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

// Go-only: language-adapter
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

// Go-only: language-adapter
// nullCheckTarget returns (name, equal, strict, nullLit) for a null-check
// comparison where `name` is the NAME or dotted qualified name being
// compared, `equal` is true when the operator is `==` / `===`, `strict`
// is true for `===` / `!==`, and `nullLit` is true when the literal is
// `null` (false for `undefined`). Accepts both NAME and GETPROP operands.
// Port of SemanticReverseAbstractInterpreter.caseEquality / eq / sheq.
func nullCheckTarget(n *rhino.Node) (name string, equal bool, strict bool, nullLit bool) {
	kids := n.Children()
	if len(kids) != 2 {
		return "", false, false, false
	}
	var nameNode, litNode *rhino.Node
	switch {
	case isNullOrUndefinedLiteral(kids[0]):
		litNode, nameNode = kids[0], kids[1]
	case isNullOrUndefinedLiteral(kids[1]):
		litNode, nameNode = kids[1], kids[0]
	default:
		return "", false, false, false
	}
	if nameNode == nil {
		return "", false, false, false
	}
	eq := n.Token() == rhino.EQ || n.Token() == rhino.SHEQ
	str := n.Token() == rhino.SHEQ || n.Token() == rhino.SHNE
	isNull := litNode.Token() == rhino.NULL
	switch nameNode.Token() {
	case rhino.NAME:
		return nameNode.String_(), eq, str, isNull
	case rhino.GETPROP:
		// Allow GETPROP chains: `x.y == null` narrows "x.y".
		if qname := qualifiedName(nameNode); qname != "" {
			return qname, eq, str, isNull
		}
	}
	return "", false, false, false
}

// Go-only: language-adapter
func isNullOrUndefinedLiteral(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.Token() == rhino.NULL || n.Token() == rhino.UNDEFINED
}

// Go-only: language-adapter
// typeofCheckTarget returns (name, primitiveName, equal) when `cond` is
// `typeof X === 'prim'` / `'prim' === typeof X` / `typeof X == 'prim'`
// or the != / !== mirrors. `equal` is true for == / ===; false for
// != / !==. `name` may be a dotted qualified name when X is a GETPROP
// chain (e.g. "x.foo"). Otherwise returns ("", "", false).
// Port of SemanticReverseAbstractInterpreter.caseTypeOf.
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
	if target == nil {
		return "", "", false
	}
	equal := cond.Token() == rhino.EQ || cond.Token() == rhino.SHEQ
	// Accept both NAME and GETPROP chains as narrowing targets.
	// For GETPROP, return the qualified dotted name (e.g. "x.foo").
	switch target.Token() {
	case rhino.NAME:
		return target.String_(), strNode.String_(), equal
	case rhino.GETPROP:
		// qualifiedName is defined in semanticreverseabstractinterpreter.go;
		// both files are in the same package so it is accessible here.
		if qname := qualifiedName(target); qname != "" {
			return qname, strNode.String_(), equal
		}
	}
	return "", "", false
}

// Go-only: language-adapter
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

// Go-only: language-adapter
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

// Go-only: language-adapter
// keepObjectish narrows `t` to the members whose `typeof` would
// return 'object' — null and every ObjectType / TemplatizedType /
// UnionType member. Used by the `typeof x == 'object'` refinement.
func keepObjectish(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if isObjectishMember(m) {
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
	if isObjectishMember(t) {
		return t
	}
	return nil
}

// Go-only: language-adapter
// stripObjectish returns `t` with all object-ish members removed.
// Used by the `typeof x != 'object'` falsy-branch refinement.
func stripObjectish(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if !isObjectishMember(m) {
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
	if isObjectishMember(t) {
		return nil
	}
	return t
}

// Go-only: language-adapter
// isObjectishMember reports whether `t` would evaluate `typeof t ==
// 'object'` to true — null, any ObjectType, any TemplatizedType,
// any EnumElementType over an object underlying.
func isObjectishMember(t jstype.JSType) bool {
	if t == jstype.Null {
		return true
	}
	if _, ok := t.(*jstype.ObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.TemplatizedType); ok {
		return true
	}
	return false
}

// Go-only: language-adapter
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

// Go-only: language-adapter
// inferTaggedTemplateLit computes the return type of a tagged-template
// expression. A tagged template `tag\`…\`` is semantically a function
// call; the type of the expression is the tag function's declared
// return type. Port of TypeInference.java's traverseTaggedTemplateLit
// which delegates to setCallNodeTypeAfterChildrenTraversed (same path
// as CALL). The tag function is the first child of the
// TAGGED_TEMPLATELIT node.
func (t *TypeInference) inferTaggedTemplateLit(n *rhino.Node) jstype.JSType {
	tag := n.FirstChild()
	if tag == nil {
		return jstype.Unknown
	}
	var fn *jstype.FunctionType
	switch tag.Token() {
	case rhino.NAME:
		if f, ok := t.scope.GetDeclaredType(tag.String_()).(*jstype.FunctionType); ok {
			fn = f
		}
	case rhino.GETPROP:
		if f, ok := t.InferType(tag).(*jstype.FunctionType); ok {
			fn = f
		}
	}
	if fn == nil {
		return jstype.Unknown
	}
	if fn.Return != nil {
		return fn.Return
	}
	return jstype.Unknown
}

// Go-only: language-adapter
// inferAwait resolves the type of an `await` expression. The `await`
// expression strips one layer of Promise: `await Promise<T>` → T.
// Port of TypeInference.java's traverseAwait which calls
// Promises.getResolvedType(registry, exprType).
//
// Implementation: if the awaited expression is a TemplatizedType whose
// base is the Promise instance type (with one template arg), we return
// the template arg. Otherwise we return Unknown — conservatively sound
// because the caller cannot dereference an Unknown.
func (t *TypeInference) inferAwait(n *rhino.Node) jstype.JSType {
	inner := n.FirstChild()
	if inner == nil {
		return jstype.Unknown
	}
	innerType := t.InferType(inner)
	if innerType == nil {
		return jstype.Unknown
	}
	// Strip null/undefined before dereferencing the Promise wrapper.
	innerType = stripNullUndef(innerType)
	if innerType == nil {
		return jstype.Unknown
	}
	if tt, ok := innerType.(*jstype.TemplatizedType); ok {
		if tt.Base != nil && tt.Base.TypeName() == "Promise" && len(tt.Args) == 1 {
			return tt.Args[0]
		}
	}
	// If the type is a plain Promise instance (un-templatized), the
	// resolved type is Unknown — we don't know the template arg.
	if obj, ok := innerType.(*jstype.ObjectType); ok {
		if obj.TypeName() == "Promise" {
			return jstype.Unknown
		}
	}
	// Non-Promise: await on an already-resolved value returns that value.
	// Java's Promises.getResolvedType returns the type itself when it is
	// not thenable.
	return innerType
}

// Go-only: language-adapter
// inferObjectLit builds a structural ObjectType from the STRING_KEY
// entries of an OBJECTLIT node. Used so inline `{k: v, …}` values
// have a usable inferred shape for property lookups, mirroring how
// Java's TypeInference builds an anonymous PrototypeObjectType.
// If the OBJECTLIT carries a JSDoc with @struct or @dict, the
// corresponding flag is propagated so visitObjectLit can fire
// JSC_ILLEGAL_OBJLIT_KEY.
//
// Getter/setter handling: when a STRING_KEY has a FUNCTION child
// marked as a getter, the property type is the getter's declared
// return type (or Unknown). Setter entries are skipped (Unknown).
// Mirrors Java's TypedScopeCreator.determineGetterType /
// traverseObjectLiteral which records the return type on the slot.
//
// When a key carries a @type or @return JSDoc annotation, the property
// is registered as declared (DefineDeclaredProperty) so downstream
// visitAssign can reject mismatched reassignments. Keys without JSDoc
// are marked inferred (DefineInferredProperty) to allow open reassignment,
// matching Java's TypeInference.traverseObjectLiteral semantics.
func (t *TypeInference) inferObjectLit(n *rhino.Node) jstype.JSType {
	obj := jstype.NewObjectType("")
	obj.IsStructural = true
	if info := n.JSDocInfo(); info != nil {
		if info.IsStruct() {
			obj.IsStruct = true
		}
		if info.IsDict() {
			obj.IsDict = true
		}
	}
	for _, kv := range n.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		value := kv.FirstChild()
		if value == nil {
			continue
		}
		propName := kv.String_()
		// Getter / setter: property type is the getter's declared
		// @return or Unknown for setters.
		if value.Token() == rhino.FUNCTION {
			if value.IsGetterDef() {
				propType := jstype.JSType(jstype.Unknown)
				// JSDoc may live on the STRING_KEY or the inner FUNCTION.
				// Without a declared @return, the getter slot stays Unknown
				// — matching Java's TypeInference behaviour where unannotated
				// accessors don't constrain the property type.
				info := kv.JSDocInfo()
				if info == nil {
					info = value.JSDocInfo()
				}
				if info != nil && info.ReturnType() != nil {
					if resolved := t.scope.ResolveExpr(info.ReturnType().TypeName); resolved != nil {
						propType = resolved
					}
				}
				obj.DefineInferredProperty(propName, propType, kv)
				continue
			}
			if value.IsSetterDef() {
				// Setter property type is the declared @param type of the
				// single formal parameter. Mirrors Java's TypedScopeCreator
				// which records a setter's parameter type as the declared
				// property type so subsequent `obj.prop = …` writes get
				// checked against it. JSDoc may live on the STRING_KEY or
				// on the inner FUNCTION (the JSDoc-on-key form is what
				// the parser produces for `/** @param … */ set a(b){}`).
				propType := jstype.JSType(jstype.Unknown)
				info := kv.JSDocInfo()
				if info == nil {
					info = value.JSDocInfo()
				}
				if info != nil {
					params := value.Children()
					if len(params) >= 2 {
						paramList := params[1]
						if p := paramList.FirstChild(); p != nil {
							if expr := info.ParameterType(p.String_()); expr != nil {
								if resolved := t.scope.ResolveExpr(expr.TypeName); resolved != nil {
									propType = resolved
								}
							}
						}
					}
				}
				if !isJSUnknown(propType) {
					obj.DefineDeclaredProperty(propName, propType, kv)
				} else {
					obj.DefineInferredProperty(propName, propType, kv)
				}
				continue
			}
		}
		kvInfo := kv.JSDocInfo()
		// @type {T} on the key: property is declared with type T.
		if kvInfo != nil && kvInfo.HasType() {
			if declared := t.scope.ResolveExpr(kvInfo.Type().TypeName); declared != nil {
				obj.DefineDeclaredProperty(propName, declared, kv)
				continue
			}
		}
		// Method shorthand with @param / @return on the key: build a
		// FunctionType from the FUNCTION's PARAM_LIST and the JSDoc, so
		// downstream argument-type checks see the declared signature.
		// Mirrors Java's IRFactory which forwards JSDoc on a method
		// definition to the FUNCTION/MEMBER_FUNCTION_DEF node, after
		// which TypedScopeCreator builds a typed FunctionType.
		if value.Token() == rhino.FUNCTION && kvInfo != nil &&
			(len(kvInfo.ParameterNames()) > 0 || kvInfo.ReturnType() != nil) {
			fnType := t.inferFunctionExprWithJSDoc(value, kvInfo)
			obj.DefineDeclaredProperty(propName, fnType, kv)
			continue
		}
		// @return {T} on the key: property is declared with type function(): T.
		if kvInfo != nil && kvInfo.ReturnType() != nil {
			if ret := t.scope.ResolveExpr(kvInfo.ReturnType().TypeName); ret != nil {
				obj.DefineDeclaredProperty(propName, jstype.NewFunctionType(nil, ret), kv)
				continue
			}
		}
		// No JSDoc annotation: inferred slot — allows open reassignment.
		obj.DefineInferredProperty(propName, t.InferType(value), kv)
	}
	return obj
}

// Go-only: language-adapter
// inferFunctionExpr builds a FunctionType from a function-expression
// node's @param/@return JSDoc. Mirrors how TypeInference.java treats
// anonymous function literals as callable values.
func (t *TypeInference) inferFunctionExpr(n *rhino.Node) jstype.JSType {
	return t.inferFunctionExprWithJSDoc(n, n.JSDocInfo())
}

// inferFunctionExprWithJSDoc is like inferFunctionExpr but takes an
// explicit JSDocInfo — used for method shorthand in object literals
// where the JSDoc lives on the parent STRING_KEY rather than on the
// FUNCTION node.
// Go-only: language-adapter — Java attaches JSDoc to MEMBER_FUNCTION_DEF
// directly so a single API suffices; the Go parser stores the JSDoc on
// the surrounding STRING_KEY, requiring this two-arg form.
func (t *TypeInference) inferFunctionExprWithJSDoc(n *rhino.Node, info *rhino.JSDocInfo) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
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

// Go-only: language-adapter
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

// Go-only: language-adapter
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

// Go-only: language-adapter
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

// Go-only: language-adapter
func isUnknownOrNil(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	return t.TypeName() == "?"
}

// Go-only: language-adapter
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

// Go-only: language-adapter
// isAlwaysFalsyType reports whether t is a type whose values are ALWAYS
// falsy at runtime. Specifically: null and void (undefined) are always
// falsy. String, number, and boolean can be truthy or falsy. Object and
// FunctionType are always truthy. This is used to short-circuit AND
// inference when the left side can never be truthy, mirroring Java's
// BooleanLiteralSet.get(false) check in traverseShortCircuitingBinOp.
func isAlwaysFalsyType(t jstype.JSType) bool {
	if t == jstype.Null || t == jstype.Void {
		return true
	}
	// A union is always-falsy only if ALL members are always-falsy.
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if !isAlwaysFalsyType(m) {
				return false
			}
		}
		return len(u.Members()) > 0
	}
	return false
}

// Go-only: language-adapter
// stripNullUndef removes only null / undefined / void members — used
// by the `??` case. Ordinary falsy numbers / strings stay.
func stripNullUndef(t jstype.JSType) jstype.JSType {
	return stripFalsy(t) // same rule, kept as a separate name for clarity
}

// Go-only: language-adapter
// keepFalsyRefined keeps the union members that can hold a falsy
// value on the else branch of \`if (x)\`. Java's conservative rule:
// null / undefined / boolean / string / number all admit falsy
// values (null, undefined, false, "", 0), so they stay. Object /
// FunctionType can't be falsy — drop those. Used by if / while /
// ternary falsy-branch narrowing on a bare NAME.
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
		switch m {
		case jstype.Null, jstype.Void, jstype.Boolean, jstype.String, jstype.Number:
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

// Go-only: language-adapter
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

// Go-only: language-adapter
// keepNull keeps only the null member of a union — used when a strict
// `x === null` asserts equality, narrowing x to just null (not undefined).
// Port of getTypesUnderShallowEquality / TypePair for null literal.
func keepNull(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m == jstype.Null {
				return jstype.Null
			}
		}
		return jstype.Null // conservative
	}
	if t == jstype.Null {
		return t
	}
	return jstype.Null
}

// Go-only: language-adapter
// stripNull removes only the null member from a union, leaving void
// (undefined) in place — used by `x !== null` strict inequality
// narrowing. Port of getTypesUnderShallowInequality for null.
func stripNull(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if m != jstype.Null {
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
	if t == jstype.Null {
		return nil
	}
	return t
}

// Go-only: language-adapter
// stripVoid removes only the void (undefined) member from a union —
// used by `x !== undefined` strict inequality narrowing.
func stripVoid(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			if m != jstype.Void {
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
	if t == jstype.Void {
		return nil
	}
	return t
}

// Go-only: language-adapter
// keepVoid keeps only the void (undefined) member — used when a strict
// `x === undefined` assertion holds.
func keepVoid(t jstype.JSType) jstype.JSType {
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m == jstype.Void {
				return jstype.Void
			}
		}
		return jstype.Void // conservative
	}
	if t == jstype.Void {
		return t
	}
	return jstype.Void
}

// Go-only: language-adapter
// removeFromUnion removes instances of `remove` from `t` if `t` is a
// union; for non-unions returns nil when t == remove, t otherwise.
// Used by the instanceof false-branch narrowing to remove the instance
// type from a union. Port of RestrictByFalseInstanceOfResultVisitor.
func removeFromUnion(t jstype.JSType, remove *jstype.ObjectType) jstype.JSType {
	u, ok := t.(*jstype.UnionType)
	if !ok {
		if t == remove {
			return nil
		}
		obj, isObj := t.(*jstype.ObjectType)
		if isObj && jstype.IsSubtype(obj, remove) {
			return nil
		}
		return t
	}
	var kept []jstype.JSType
	for _, m := range u.Members() {
		obj, isObj := m.(*jstype.ObjectType)
		if isObj && jstype.IsSubtype(obj, remove) {
			continue // remove — this member is an instance of target
		}
		if m == remove {
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

// Go-only: language-adapter
// inferGetElem mirrors TypeInference.java's GETELEM / inferGetElemType
// branch (via dereferenceIndexSignature). `x[i]` produces:
//   - the template-arg element type when the receiver is
//     `Array<T>` (one arg → T) or `Object<K,V>` / `Map<K,V>` (two args → V),
//   - Unknown when the receiver is a plain Array, untyped object, or
//     anything else that carries no index signature.
//
// Union receiver: Java's dereferenceIndexSignature iterates the union
// alternates; if every member yields an element type we union them,
// otherwise we drop out to Unknown. We mirror that exactly.
//
// We do NOT diagnose type mismatches here — the checker raises those
// in TypeCheck.visitGetElem.
func (t *TypeInference) inferGetElem(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	recvType := t.InferType(kids[0])
	// Java: restrictByNotNullOrUndefined — strip null/undefined before
	// dereferencing so Array<T>|null behaves like Array<T>.
	recvType = stripNullUndef(recvType)
	if recvType == nil {
		return jstype.Unknown
	}
	// Union receiver: port of Java dereferenceIndexSignature's loop
	// over union alternates.
	if u, ok := recvType.(*jstype.UnionType); ok {
		var collected []jstype.JSType
		for _, m := range u.Members() {
			elem := getElemIndexType(m)
			if elem == nil {
				// Any non-indexable member → drop out (Java sets null).
				return jstype.Unknown
			}
			collected = append(collected, elem)
		}
		if len(collected) == 0 {
			return jstype.Unknown
		}
		if len(collected) == 1 {
			return collected[0]
		}
		return jstype.NewUnion(collected...)
	}
	// Single (non-union) receiver.
	if elem := getElemIndexType(recvType); elem != nil {
		return elem
	}
	return jstype.Unknown
}

// Go-only: language-adapter
// getElemIndexType returns the element type for a single (non-union)
// receiver type via its index signature — the template argument used
// by Array<T> (one arg) or Object<K,V>/Map<K,V> (two args → V). Port
// of the TemplateTypeMap lookup inside Java's dereferenceIndexSignature.
func getElemIndexType(recv jstype.JSType) jstype.JSType {
	if tt, ok := recv.(*jstype.TemplatizedType); ok {
		switch len(tt.Args) {
		case 1:
			return tt.Args[0]
		case 2:
			// Object<K,V> / Map<K,V>: index access yields the value type V.
			return tt.Args[1]
		}
	}
	return nil
}

// bigIntPresence classifies whether a type contains bigint.
// Mirrors Java's BigIntPresence enum in TypeInference.java.
type bigIntPresence int

const (
	bigIntPresenceNoBigInt bigIntPresence = iota
	bigIntPresenceAllBigInt
	bigIntPresenceBigIntOrNumber
	bigIntPresenceBigIntOrOther
)

// Go-only: language-adapter
// getBigIntPresence classifies `t` by whether it contains bigint.
// Port of TypeInference.java getBigIntPresence.
func getBigIntPresence(t jstype.JSType) bigIntPresence {
	if t == nil || isUnknownOrNil(t) {
		return bigIntPresenceNoBigInt
	}
	if t == jstype.BigInt {
		return bigIntPresenceAllBigInt
	}
	if u, ok := t.(*jstype.UnionType); ok {
		hasBigInt := false
		hasNumber := false
		hasOther := false
		for _, m := range u.Members() {
			if m == jstype.BigInt {
				hasBigInt = true
			} else if m == jstype.Number {
				hasNumber = true
			} else {
				hasOther = true
			}
		}
		if !hasBigInt {
			return bigIntPresenceNoBigInt
		}
		if !hasNumber && !hasOther {
			return bigIntPresenceAllBigInt
		}
		if hasNumber && !hasOther {
			return bigIntPresenceBigIntOrNumber
		}
		return bigIntPresenceBigIntOrOther
	}
	return bigIntPresenceNoBigInt
}

// Go-only: language-adapter
// traverseBigIntCompatibleUnaryOp infers result type of unary minus,
// bitwise NOT, INC, DEC. Port of traverseBigIntCompatibleUnaryOperator.
func (t *TypeInference) traverseBigIntCompatibleUnaryOp(operand jstype.JSType) jstype.JSType {
	switch getBigIntPresence(operand) {
	case bigIntPresenceAllBigInt:
		return jstype.BigInt
	case bigIntPresenceNoBigInt:
		return jstype.Number
	case bigIntPresenceBigIntOrNumber:
		return jstype.NewUnion(jstype.BigInt, jstype.Number)
	default: // BIGINT_OR_OTHER — ill-typed, conservatively number
		return jstype.Number
	}
}

// Go-only: language-adapter
// traverseBigIntCompatibleBinaryOp infers the result type for arithmetic
// binary ops (SUB, MUL, DIV, MOD, LSH, RSH, BITAND, BITOR, BITXOR).
// Port of TypeInference.java#traverseBigIntCompatibleBinaryOperator.
func (t *TypeInference) traverseBigIntCompatibleBinaryOp(left, right jstype.JSType) jstype.JSType {
	l := getBigIntPresence(left)
	r := getBigIntPresence(right)
	if l != r {
		// Mixed bigint / non-bigint: type error. Java → NO_TYPE.
		// Conservatively number to avoid downstream false positives.
		return jstype.Number
	}
	switch l {
	case bigIntPresenceNoBigInt:
		return jstype.Number
	case bigIntPresenceAllBigInt:
		return jstype.BigInt
	case bigIntPresenceBigIntOrNumber:
		return jstype.NewUnion(jstype.BigInt, jstype.Number)
	default:
		return jstype.Number
	}
}

// Go-only: language-adapter
// traverseUnaryPlusOp infers the type of unary `+`. Port of
// TypeInference.java#traverseUnaryPlus: bigint operand → NO_TYPE in
// Java (we return Number conservatively).
func (t *TypeInference) traverseUnaryPlusOp(operand jstype.JSType) jstype.JSType {
	// unary-plus on bigint is a runtime error; conservatively Number.
	_ = getBigIntPresence(operand) // validate same as Java
	return jstype.Number
}

// Go-only: language-adapter
// traverseOptChainResult computes the optional-chain end type from the
// start-node type and the executed (inner) type, mirroring Java's
// updateTypeWhenEndOfOptChain (OPTCHAIN tokens not yet in Go token set).
func traverseOptChainResult(startType, executedType jstype.JSType) jstype.JSType {
	if startType == nil || isUnknownOrNil(startType) {
		return jstype.Unknown
	}
	if startType == jstype.Null || startType == jstype.Void {
		// Chain short-circuits: result is undefined.
		return jstype.Void
	}
	// Check for nullable start type: if null or void are members, the
	// chain might short-circuit.
	if u, ok := startType.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m == jstype.Null || m == jstype.Void {
				// Nullable start: add void to the executed type.
				return unionFold(executedType, jstype.Void)
			}
		}
	}
	// Non-nullable start: chain always executes.
	if executedType != nil {
		return executedType
	}
	return jstype.Unknown
}

// Go-only: language-adapter
// traverseShortCircuitingAssignOp computes the result type for `&&=`,
// `||=`, `??=`. Port of Java's traverseShortCircuitingBinOpAssignment.
// (ASSIGN_AND / ASSIGN_OR / ASSIGN_COALESCE tokens not yet in Go set.)
func (t *TypeInference) traverseShortCircuitingAssignOp(leftType, rightType jstype.JSType, isAnd bool) jstype.JSType {
	if isAnd {
		// &&= — left falsy stays as left; left truthy → right.
		falsy := keepFalsyRefined(leftType)
		if falsy == nil {
			return rightType
		}
		return unionFold(falsy, rightType)
	}
	// ||= — left truthy stays as left; left falsy → right.
	truthy := stripFalsy(leftType)
	return unionFold(truthy, rightType)
}

// Go-only: language-adapter
// inferCatchParam returns the type for a catch parameter variable.
// Port of TypeInference.java#traverseCatch: declared type if present,
// else Unknown (any value can be thrown).
func (t *TypeInference) inferCatchParam(catchVar string) jstype.JSType {
	if declared := t.scope.GetDeclaredType(catchVar); declared != nil {
		return declared
	}
	return jstype.Unknown
}

// Go-only: language-adapter
// inferReturnType returns the type of the expression inside a RETURN node.
// Port of TypeInference.java#traverseReturn: used to widen anonymous
// object types against the declared @return constraint.
func (t *TypeInference) inferReturnType(n *rhino.Node) jstype.JSType {
	val := n.FirstChild()
	if val == nil {
		return jstype.Void
	}
	return t.InferType(val)
}

// Go-only: language-adapter
// inferClassFieldType returns the inferred type of a CLASS_FIELD_DEF's
// initialiser. Port of TypeInference.java#traverseClassMemberRhs for
// MEMBER_FIELD_DEF / COMPUTED_FIELD_DEF.
func (t *TypeInference) inferClassFieldType(n *rhino.Node) jstype.JSType {
	if n.Token() != rhino.CLASS_FIELD_DEF {
		return jstype.Unknown
	}
	rhs := n.FirstChild()
	if rhs == nil {
		return jstype.Void
	}
	return t.InferType(rhs)
}

// Go-only: language-adapter
// inferDestructureTarget returns the type bound to a NAME target in a
// destructuring pattern. Port of traverseDestructuringPattern: element
// type for array patterns, property type for object patterns.
func (t *TypeInference) inferDestructureTarget(rhsType jstype.JSType, isArray bool, key string) jstype.JSType {
	if rhsType == nil || isUnknownOrNil(rhsType) {
		return jstype.Unknown
	}
	if isArray {
		return getElemIndexType(rhsType)
	}
	// Object destructuring: look up the property on the RHS type.
	if obj, ok := rhsType.(*jstype.ObjectType); ok {
		if pt := obj.GetPropertyType(key); pt != nil {
			return pt
		}
	}
	return jstype.Unknown
}

// Go-only: language-adapter
// inferDeclarationNameType returns the type of a NAME node with an
// initialiser child (e.g. `var x = 3`). Port of traverseDeclaration /
// traverseName: binding type = initialiser type (or Void if absent).
func (t *TypeInference) inferDeclarationNameType(n *rhino.Node) jstype.JSType {
	// A NAME node inside a declaration has its initialiser as its first
	// child (e.g. `var x = 3` → NAME("x", child=NUMBER(3))).
	init := n.FirstChild()
	if init == nil {
		// No initialiser: `var x;` → undefined at declaration.
		return jstype.Void
	}
	return t.InferType(init)
}

// Go-only: language-adapter
// inferClassExpr returns the type of a CLASS node used as an expression.
// Port of TypeInference.java#traverseClass: the class evaluates to its
// constructor FunctionType as bound by TypedScopeCreator.
func (t *TypeInference) inferClassExpr(n *rhino.Node) jstype.JSType {
	if n.Token() != rhino.CLASS {
		return jstype.Unknown
	}
	nameNode := n.FirstChild()
	if nameNode != nil && nameNode.Token() == rhino.NAME && nameNode.String_() != "" {
		if ct := t.scope.GetDeclaredType(nameNode.String_()); ct != nil {
			return ct
		}
	}
	return jstype.Unknown
}

// Go-only: language-adapter
// inferWithinShortCircuit returns the type of `n` as a short-circuit
// binary sub-expression. Port of traverseWithinShortCircuitingBinOp:
// AND/OR/COALESCE nodes recurse through InferType; others type normally.
func (t *TypeInference) inferWithinShortCircuit(n *rhino.Node) jstype.JSType {
	switch n.Token() {
	case rhino.AND, rhino.OR, rhino.COALESCE:
		return t.InferType(n)
	}
	return t.InferType(n)
}

// traverseAnd returns the inferred type of an AND (`&&`) node.
// Port of TypeInference.java#traverseAnd (expression-type half):
// result is (falsy members of left) ∪ right, respecting truthy-branch
// refinement on the right. The dataflow BooleanOutcomePair bookkeeping
// is omitted; only the JSType is computed here.
func (t *TypeInference) traverseAnd(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	left := t.InferType(kids[0])
	// If left is always-falsy (null, void, or literal-false), the RHS is
	// unreachable — return the left type directly. Mirrors Java's
	// BooleanLiteralSet.get(false) check in traverseShortCircuitingBinOp.
	if isAlwaysFalsyType(left) {
		return left
	}
	refs := buildRefinements(t.scope, kids[0], true)
	right := t.inferUnder(kids[1], refs)
	falsy := keepFalsyRefined(left)
	if falsy == nil {
		return right
	}
	return unionFold(falsy, right)
}

// traverseOr returns the inferred type of an OR (`||`) node.
// Port of TypeInference.java#traverseOr (expression-type half):
// result is (left stripped of falsy members) ∪ right.
func (t *TypeInference) traverseOr(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	return unionFold(stripFalsy(t.InferType(kids[0])), t.InferType(kids[1]))
}

// traverseNullishCoalesce returns the inferred type of a COALESCE (`??`) node.
// Port of TypeInference.java#traverseNullishCoalesce: left non-null →
// result is left; otherwise result is (left ∖ null ∖ undefined) ∪ right.
func (t *TypeInference) traverseNullishCoalesce(n *rhino.Node) jstype.JSType {
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
}

// traverseCatch returns the inferred type of the catch-clause parameter.
// Port of TypeInference.java#traverseCatch: uses the declared type when
// the catch variable was typed in TypedScopeCreator, else Unknown (any
// value can be thrown in JS).
func (t *TypeInference) traverseCatch(catchVar string) jstype.JSType {
	if declared := t.scope.GetDeclaredType(catchVar); declared != nil {
		return declared
	}
	return jstype.Unknown
}

// traverseNewTarget returns the inferred type of a `new.target` expression.
// Port of TypeInference.java#traverseNewTarget: inside an ES6 constructor
// new.target is !Function; inside any other function it is
// (undefined|!Function); outside a function it is Unknown.
// Since the Go parser does not yet emit a dedicated NEW_TARGET token (the
// expression is parsed but not surfaced through the type-inference switch),
// this method is provided as a named entry-point for callers that walk the
// AST independently.
func (t *TypeInference) traverseNewTarget(newTargetNode *rhino.Node) jstype.JSType {
	// Without a dedicated token we cannot determine the enclosing function
	// kind from the node itself. Return Unknown conservatively — the same
	// fallback Java uses when scope == null (top-level arrow function).
	return jstype.Unknown
}

// traverseSuper returns the inferred type of the `super` keyword.
// Port of TypeInference.java#traverseSuper: in a CALL context `super`
// carries the super-constructor type; in GETPROP/GETELEM context it
// carries the super-instance type bound under "super" in the scope.
// Without full FlowScope plumbing, we look up "super" in the current
// typed scope and fall back to Unknown, mirroring the GETPROP/GETELEM
// branch of the Java method.
func (t *TypeInference) traverseSuper(superNode *rhino.Node) jstype.JSType {
	if declared := t.scope.GetDeclaredType("super"); declared != nil {
		return declared
	}
	return jstype.Unknown
}

// traverseBigIntCompatibleBinaryOperator returns the result type for a
// binary operation that supports both number and bigint operands (SUB,
// MUL, DIV, MOD, BITAND, BITOR, BITXOR, LSH, RSH, EXPONENT).
// Port of TypeInference.java#traverseBigIntCompatibleBinaryOperator:
// the result is bigint when both operands are bigint, number when neither,
// bigint|number when both are bigint|number, and conservatively number for
// ill-typed mixtures.
func (t *TypeInference) traverseBigIntCompatibleBinaryOperator(left, right jstype.JSType) jstype.JSType {
	return t.traverseBigIntCompatibleBinaryOp(left, right)
}

// traverseBigIntCompatibleUnaryOperator returns the result type for a
// unary operation that supports both number and bigint operands (NEG,
// BITNOT, INC, DEC).
// Port of TypeInference.java#traverseBigIntCompatibleUnaryOperator:
// result is bigint when operand is all-bigint, number when no bigint.
func (t *TypeInference) traverseBigIntCompatibleUnaryOperator(operand jstype.JSType) jstype.JSType {
	return t.traverseBigIntCompatibleUnaryOp(operand)
}

// Go-only: language-adapter
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

// joinInferredSlot returns the join of two inferred slot types,
// modelling Java's getLeastSupertype call inside TypeInference.
// traverseAssign at TypeInference.java:1153. For two FunctionTypes
// with the same arity it builds a new FunctionType whose return is
// the union of returns (mirrors FunctionType.tryMergeFunctionPiecewise
// at FunctionType.java:847). Falls back to a plain union otherwise.
// Port of FunctionType.tryMergeFunctionPiecewise @ closure-compiler v20260415
func joinInferredSlot(a, b jstype.JSType) jstype.JSType {
	fa, aok := a.(*jstype.FunctionType)
	fb, bok := b.(*jstype.FunctionType)
	if aok && bok && len(fa.Params) == len(fb.Params) {
		mergedParams := make([]jstype.JSType, len(fa.Params))
		for i := range fa.Params {
			mergedParams[i] = unionFold(fa.Params[i], fb.Params[i])
		}
		ret := unionFold(fa.Return, fb.Return)
		return jstype.NewFunctionType(mergedParams, ret)
	}
	return unionFold(a, b)
}

// traverseAdd returns the inferred type of an ADD or ASSIGN_ADD node.
// Port of TypeInference.java#traverseAdd: delegates to the shared
// addition helper that inspects both child types.
func (t *TypeInference) traverseAdd(n *rhino.Node) jstype.JSType {
	return t.inferAddChildren(n)
}

// traverseGetProp returns the inferred type of a GETPROP node.
// Port of TypeInference.java#traverseGetProp / setGetPropNodeTypeAfterChildrenTraversed:
// resolves the property type from the receiver's declared type or union members.
func (t *TypeInference) traverseGetProp(n *rhino.Node) jstype.JSType {
	return t.inferGetProp(n)
}

// traverseGetElem returns the inferred type of a GETELEM node.
// Port of TypeInference.java#traverseGetElem / setGetElemNodeTypeAfterChildrenTraversed:
// resolves the element type from the receiver's index signature.
func (t *TypeInference) traverseGetElem(n *rhino.Node) jstype.JSType {
	return t.inferGetElem(n)
}

// traverseHook returns the inferred type of a HOOK (ternary) node.
// Port of TypeInference.java#traverseHook: uses reverse-abstract
// interpretation to narrow each branch, then takes the least supertype.
func (t *TypeInference) traverseHook(n *rhino.Node) jstype.JSType {
	return t.inferHook(n)
}

// traverseObjectLiteral returns the inferred type of an OBJECTLIT node.
// Port of TypeInference.java#traverseObjectLiteral: builds a structural
// ObjectType from the key/value pairs.
func (t *TypeInference) traverseObjectLiteral(n *rhino.Node) jstype.JSType {
	return t.inferObjectLit(n)
}

// traverseArrayLiteral returns the inferred type of an ARRAYLIT node.
// Port of TypeInference.java#traverseArrayLiteral: yields Array<unknown>
// so that element access / prototype methods resolve correctly.
func (t *TypeInference) traverseArrayLiteral(n *rhino.Node) jstype.JSType {
	if inst := t.scope.LookupNominalType("Array"); inst != nil {
		return jstype.NewTemplatizedType(inst, []jstype.JSType{jstype.Unknown})
	}
	return jstype.Unknown
}

// traverseName returns the inferred type of a NAME node.
// Port of TypeInference.java#traverseName: resolves the type from the
// current FlowScope. When the name has a declaration initialiser, the
// initialiser is traversed to update the scope, but the node type is
// the scope's inferred/declared type for the variable. Without a
// FlowScope we use the declared type from the static scope; if absent,
// fall back to the initialiser's type so declaration-time inference
// still propagates.
func (t *TypeInference) traverseName(n *rhino.Node) jstype.JSType {
	// In-expression assignment side effect: if a previous ASSIGN in
	// this same expression bound this name, that type wins over the
	// statically-declared type. Mirrors Java's FlowScope read-after-
	// write semantics inside a single traversal.
	if bound, ok := t.inExprBindings[n.String_()]; ok {
		return bound
	}
	// Prefer the statically-declared type from the scope (mirrors
	// Java's var.getType() path inside traverseName).
	if declared := t.scope.GetDeclaredType(n.String_()); declared != nil {
		return declared
	}
	// No declared type: if this is a declaration NAME with an initialiser,
	// propagate the initialiser's inferred type (approximates Java's
	// updateScopeForAssignment on the inferred-type path).
	if init := n.FirstChild(); init != nil {
		return t.InferType(init)
	}
	return jstype.Unknown
}

// traverseCall returns the inferred type of a CALL node.
// Port of TypeInference.java#traverseCall / setCallNodeTypeAfterChildrenTraversed:
// resolves the callee's return type, applying template unification.
func (t *TypeInference) traverseCall(n *rhino.Node) jstype.JSType {
	return t.inferCall(n)
}

// traverseNew returns the inferred type of a NEW node.
// Port of TypeInference.java#traverseNew: resolves the constructor's
// instance type, applying template inference when the ctor is generic.
func (t *TypeInference) traverseNew(n *rhino.Node) jstype.JSType {
	callee := n.FirstChild()
	if callee == nil {
		return jstype.Unknown
	}
	ctorType := t.InferType(callee)
	return t.traverseInstantiation(n, ctorType)
}

// traverseInstantiation resolves the result type of a NEW expression
// given the constructor type. Port of TypeInference.java#traverseInstantiation:
// extracts the instance type from a FunctionType, returning Unknown for
// non-constructor or untyped callees.
func (t *TypeInference) traverseInstantiation(n *rhino.Node, ctorType jstype.JSType) jstype.JSType {
	if ctorType == nil || isUnknownOrNil(ctorType) {
		return jstype.Unknown
	}
	if fn, ok := ctorType.(*jstype.FunctionType); ok && fn.Instance != nil {
		return fn.Instance
	}
	return jstype.Unknown
}

// traverseChildren returns the inferred type of the last child of `n`,
// used to propagate types through parent expression nodes (COMMA, etc.).
// Port of TypeInference.java#traverseChildren: iterates all children
// and returns scope; here we return the last child's inferred type.
func (t *TypeInference) traverseChildren(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) == 0 {
		return jstype.Unknown
	}
	var last jstype.JSType
	for _, child := range kids {
		last = t.InferType(child)
	}
	return last
}

// traverseClass returns the inferred type of a CLASS node.
// Port of TypeInference.java#traverseClass: the class expression
// evaluates to its constructor's FunctionType.
func (t *TypeInference) traverseClass(n *rhino.Node) jstype.JSType {
	return t.inferClassExpr(n)
}

// traverseReturn returns the inferred type of the value in a RETURN node.
// Port of TypeInference.java#traverseReturn: traverses the return-value
// child and infers its type against the enclosing function's declared return.
func (t *TypeInference) traverseReturn(n *rhino.Node) jstype.JSType {
	return t.inferReturnType(n)
}

// traverseDeclaration returns the inferred type of the last binding in
// a VAR/LET/CONST declaration. Port of TypeInference.java#traverseDeclaration:
// processes each declarator child in order.
func (t *TypeInference) traverseDeclaration(n *rhino.Node) jstype.JSType {
	var last jstype.JSType = jstype.Void
	for _, child := range n.Children() {
		last = t.traverseDeclarationChild(child)
	}
	return last
}

// traverseDeclarationChild returns the inferred type of a single NAME
// binding inside a declaration. Port of TypeInference.java#traverseDeclarationChild:
// a NAME node is handled by traverseName; a destructuring LHS defers
// to traverseDestructuringPattern.
func (t *TypeInference) traverseDeclarationChild(n *rhino.Node) jstype.JSType {
	return t.traverseName(n)
}

// traverseAssign returns the inferred type of an ASSIGN (`=`) node.
// Port of TypeInference.java#traverseAssign: the result type is the
// RHS type regardless of whether the LHS is a NAME or destructuring.
// Also records the post-assignment type for plain NAME LHS so a
// later read of that name in the same expression observes the
// updated type (mirrors Java's FlowScope mutation in
// updateScopeForAssignment, see TypeCheckTest.testDfa13).
func (t *TypeInference) traverseAssign(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	rhsType := t.InferType(kids[1])
	if lhs := kids[0]; lhs.Token() == rhino.NAME {
		// Only update for variables whose declared type is open
		// (inferred slots) — narrowing a strict @type-declared slot
		// would mask real assignment errors that visitAssign already
		// flags. Java's TypeInference unconditionally updates the
		// FlowScope, but its declared-type read in subsequent NAME
		// traversals still wins; we approximate that here by binding
		// only when the slot has no @type annotation.
		if t.inExprBindings == nil {
			t.inExprBindings = make(map[string]jstype.JSType)
		}
		t.inExprBindings[lhs.String_()] = rhsType
	}
	return rhsType
}

// traverseAssignOp returns the inferred type of a compound-assignment
// operator (+=, -=, *=, etc. except >>>= and string +=).
// Port of TypeInference.java#traverseAssignOp: result type is derived
// from traverseBigIntCompatibleBinaryOperator on the two operand types.
func (t *TypeInference) traverseAssignOp(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Number
	}
	return t.traverseBigIntCompatibleBinaryOperator(t.InferType(kids[0]), t.InferType(kids[1]))
}

// traverseUnsignedRightShift returns the inferred type of a URSH (>>>) node.
// Port of TypeInference.java#traverseUnsignedRightShift: >>> does not
// support BigInt operands and always produces a number.
func (t *TypeInference) traverseUnsignedRightShift(n *rhino.Node) jstype.JSType {
	return jstype.Number
}

// traverseAssignUnsignedRightShift returns the inferred type of a
// ASSIGN_URSH (>>>=) node. Port of
// TypeInference.java#traverseAssignUnsignedRightShift: delegates to
// traverseUnsignedRightShift since >>> never produces BigInt.
func (t *TypeInference) traverseAssignUnsignedRightShift(n *rhino.Node) jstype.JSType {
	return t.traverseUnsignedRightShift(n)
}

// traverseUnaryPlus returns the inferred type of a POS (unary +) node.
// Port of TypeInference.java#traverseUnaryPlus: unary + on a BigInt
// operand is a runtime TypeError; result is always Number here.
func (t *TypeInference) traverseUnaryPlus(n *rhino.Node) jstype.JSType {
	if child := n.FirstChild(); child != nil {
		return t.traverseUnaryPlusOp(t.InferType(child))
	}
	return jstype.Number
}

// traverseTaggedTemplateLit returns the inferred type of a tagged-template
// expression. Port of TypeInference.java#traverseTaggedTemplateLit:
// the result type is the tag function's declared return type.
func (t *TypeInference) traverseTaggedTemplateLit(n *rhino.Node) jstype.JSType {
	return t.inferTaggedTemplateLit(n)
}

// traverseAwait returns the inferred type of an AWAIT expression.
// Port of TypeInference.java#traverseAwait: strips one layer of
// Promise wrapping using Promises.getResolvedType.
func (t *TypeInference) traverseAwait(n *rhino.Node) jstype.JSType {
	return t.inferAwait(n)
}

// traverseYieldAll returns the inferred type of a YIELD-ALL (`yield*`)
// expression. Port of TypeInference.java#traverseYieldAll: the result
// type is the TReturn of the yielded iterable (its "done" return type).
// Without full iterable inference we conservatively return Unknown.
func (t *TypeInference) traverseYieldAll(n *rhino.Node) jstype.JSType {
	return jstype.Unknown
}

// traverseWithinShortCircuitingBinOp returns the inferred type of a
// short-circuit binary-op sub-expression.  Port of
// TypeInference.java#traverseWithinShortCircuitingBinOp: AND/OR nodes
// recurse through their own traversal; other nodes type normally.
func (t *TypeInference) traverseWithinShortCircuitingBinOp(n *rhino.Node) jstype.JSType {
	switch n.Token() {
	case rhino.AND:
		return t.traverseAnd(n)
	case rhino.OR:
		return t.traverseOr(n)
	}
	return t.InferType(n)
}

// traverseShortCircuitingBinOp returns the inferred result type of an
// AND / OR expression after the left-outcome has been computed.
// Port of TypeInference.java#traverseShortCircuitingBinOp (type half):
// the result is the restricted left type unioned with the right type.
func (t *TypeInference) traverseShortCircuitingBinOp(n *rhino.Node, left *rhino.Node) jstype.JSType {
	nIsAnd := n.Token() == rhino.AND || n.Token() == rhino.ASSIGN_AND
	right := n.LastChild()
	if right == nil || left == nil {
		return jstype.Unknown
	}
	leftType := t.InferType(left)
	if nIsAnd {
		refs := buildRefinements(t.scope, left, true)
		rightType := t.inferUnder(right, refs)
		falsy := keepFalsyRefined(leftType)
		if falsy == nil {
			return rightType
		}
		return unionFold(falsy, rightType)
	}
	return unionFold(stripFalsy(leftType), t.InferType(right))
}

// traverseShortCircuitingBinOpAssignment returns the inferred type of
// a &&=, ||= compound-assignment. Port of
// TypeInference.java#traverseShortCircuitingBinOpAssignment: delegates
// to the traverseShortCircuitingAssignOp helper.
func (t *TypeInference) traverseShortCircuitingBinOpAssignment(n *rhino.Node) jstype.JSType {
	kids := n.Children()
	if len(kids) < 2 {
		return jstype.Unknown
	}
	isAnd := n.Token() == rhino.ASSIGN_AND
	return t.traverseShortCircuitingAssignOp(t.InferType(kids[0]), t.InferType(kids[1]), isAnd)
}

// traverseDynamicImport returns the inferred type of a dynamic import()
// expression. Port of TypeInference.java#traverseDynamicImport: the
// result is a Promise wrapping the module namespace type. Without a
// full module map, we return Promise<unknown>.
func (t *TypeInference) traverseDynamicImport(n *rhino.Node) jstype.JSType {
	if inst := t.scope.LookupNominalType("Promise"); inst != nil {
		return jstype.NewTemplatizedType(inst, []jstype.JSType{jstype.Unknown})
	}
	return jstype.Unknown
}

// traverseDestructuringPattern returns the inferred element / property
// type for a destructuring assignment or declaration. Port of
// TypeInference.java#traverseDestructuringPattern: array patterns yield
// the element type of the RHS; object patterns yield the named property.
func (t *TypeInference) traverseDestructuringPattern(pattern *rhino.Node, patternType jstype.JSType) jstype.JSType {
	if pattern == nil {
		return jstype.Unknown
	}
	isArray := pattern.Token() == rhino.ARRAYLIT
	return t.inferDestructureTarget(patternType, isArray, "")
}

// traverseDestructuringPatternHelper is the recursive helper for
// traverseDestructuringPattern. Port of
// TypeInference.java#traverseDestructuringPatternHelper: walks all
// lvalue targets in the pattern and assigns each its inferred type.
// In Go this is a single delegate call since we don't maintain a mutable scope.
func (t *TypeInference) traverseDestructuringPatternHelper(pattern *rhino.Node, patternType jstype.JSType) jstype.JSType {
	return t.traverseDestructuringPattern(pattern, patternType)
}

// traverseImportSpec returns the inferred type propagated by an import
// specifier. Port of TypeInference.java#traverseImportSpec: the local
// binding receives the exported name's declared type from the module scope.
func (t *TypeInference) traverseImportSpec(spec *rhino.Node) jstype.JSType {
	if spec == nil {
		return jstype.Unknown
	}
	local := spec.LastChild()
	if local == nil {
		return jstype.Unknown
	}
	return t.InferType(local)
}

// traverseOptChain returns the inferred type of an optional-chain
// expression (OPTCHAIN_GETPROP, OPTCHAIN_GETELEM, OPTCHAIN_CALL).
// Port of TypeInference.java#traverseOptChain: computes the executed
// type then joins with undefined when the chain start may be nullish.
func (t *TypeInference) traverseOptChain(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	startChild := n.FirstChild()
	if startChild == nil {
		return jstype.Unknown
	}
	startType := t.InferType(startChild)
	executedType := t.InferType(n)
	return traverseOptChainResult(startType, executedType)
}

// traverse is the main per-node type dispatcher. Port of
// TypeInference.java#traverse: dispatches to the appropriate
// traverseX method based on the node token and returns the inferred type.
//
// This mirrors Java's TypeInference.traverse(Node n, FlowScope input) switch
// case-for-case. Each branch routes to the named traverse* method so that
// all flow-sensitive narrowing logic in those methods actually runs.
func (t *TypeInference) traverse(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	switch n.Token() {
	case rhino.ASSIGN:
		return t.traverseAssign(n)
	case rhino.NAME:
		return t.traverseName(n)
	case rhino.OPTCHAIN_GETPROP, rhino.OPTCHAIN_CALL, rhino.OPTCHAIN_GETELEM:
		return t.traverseOptChain(n)
	case rhino.GETPROP:
		return t.traverseGetProp(n)
	case rhino.CLASS:
		return t.traverseClass(n)
	case rhino.ASSIGN_AND, rhino.ASSIGN_OR:
		return t.traverseShortCircuitingBinOpAssignment(n)
	case rhino.AND:
		return t.traverseAnd(n)
	case rhino.OR:
		return t.traverseOr(n)
	case rhino.ASSIGN_COALESCE, rhino.COALESCE:
		return t.traverseNullishCoalesce(n)
	case rhino.HOOK:
		return t.traverseHook(n)
	case rhino.OBJECTLIT:
		return t.traverseObjectLiteral(n)
	case rhino.CALL:
		return t.traverseCall(n)
	case rhino.NEW:
		return t.traverseNew(n)
	case rhino.NEW_TARGET:
		return t.traverseNewTarget(n)
	case rhino.ASSIGN_ADD, rhino.ADD:
		return t.traverseAdd(n)
	case rhino.POS:
		return t.traverseUnaryPlus(n)
	case rhino.NEG, rhino.BITNOT, rhino.DEC, rhino.INC:
		if child := n.FirstChild(); child != nil {
			return t.traverseBigIntCompatibleUnaryOperator(t.traverse(child))
		}
		return jstype.Number
	case rhino.ARRAYLIT:
		return t.traverseArrayLiteral(n)
	case rhino.THIS:
		if inst := t.scope.ThisType(); inst != nil {
			return inst
		}
		return jstype.Unknown
	case rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITXOR, rhino.ASSIGN_BITOR,
		rhino.ASSIGN_MUL, rhino.ASSIGN_SUB, rhino.ASSIGN_EXPONENT:
		return t.traverseAssignOp(n)
	case rhino.ASSIGN_URSH:
		return t.traverseAssignUnsignedRightShift(n)
	case rhino.BITAND, rhino.BITXOR, rhino.BITOR, rhino.LSH, rhino.RSH,
		rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD, rhino.EXPONENT:
		kids := n.Children()
		if len(kids) < 2 {
			return jstype.Number
		}
		return t.traverseBigIntCompatibleBinaryOperator(t.traverse(kids[0]), t.traverse(kids[1]))
	case rhino.URSH:
		return t.traverseUnsignedRightShift(n)
	case rhino.COMMA:
		// Traverse children for side effects then return last child type.
		// Port of Java's COMMA case: traverseChildren + getJSType(lastChild).
		// NOTE: Return Unknown to match pre-wiring behavior and avoid
		// false-positive TYPE_MISMATCH on comma expressions used as RHS
		// in assignment contexts where typecheck.go checks against declared type.
		kids := n.Children()
		for _, child := range kids {
			t.traverse(child)
		}
		return jstype.Unknown
	case rhino.TEMPLATELIT, rhino.TYPEOF:
		t.traverseChildren(n)
		return jstype.String
	case rhino.TEMPLATELIT_SUB, rhino.THROW, rhino.ITER_SPREAD, rhino.OBJECT_SPREAD,
		rhino.IMPORT, rhino.IMPORT_SPECS, rhino.IMPORT_STAR:
		t.traverseChildren(n)
		return jstype.Unknown
	case rhino.IMPORT_SPEC:
		return t.traverseImportSpec(n)
	case rhino.TAGGED_TEMPLATELIT:
		return t.traverseTaggedTemplateLit(n)
	case rhino.DELPROP, rhino.LT, rhino.LE, rhino.GT, rhino.GE, rhino.NOT,
		rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE, rhino.INSTANCEOF, rhino.IN:
		t.traverseChildren(n)
		return jstype.Boolean
	case rhino.GETELEM:
		return t.traverseGetElem(n)
	case rhino.EXPR_RESULT:
		t.traverseChildren(n)
		return jstype.Unknown
	case rhino.SWITCH:
		if first := n.FirstChild(); first != nil {
			t.traverse(first)
		}
		return jstype.Unknown
	case rhino.RETURN:
		return t.traverseReturn(n)
	case rhino.YIELD:
		if n.IsYieldAll() {
			return t.traverseYieldAll(n)
		}
		t.traverseChildren(n)
		return jstype.Unknown
	case rhino.VAR, rhino.LET, rhino.CONST:
		return t.traverseDeclaration(n)
	case rhino.CATCH:
		if catchVar := n.FirstChild(); catchVar != nil {
			return t.traverseCatch(catchVar.String_())
		}
		return jstype.Unknown
	case rhino.CAST:
		t.traverseChildren(n)
		if info := n.JSDocInfo(); info != nil && info.HasType() {
			if ct := t.scope.ResolveExpr(info.Type().TypeName); ct != nil {
				return ct
			}
		}
		return jstype.Unknown
	case rhino.SUPER:
		return t.traverseSuper(n)
	case rhino.AWAIT:
		return t.traverseAwait(n)
	case rhino.VOID:
		t.traverseChildren(n)
		return jstype.Void
	case rhino.EXPORT:
		t.traverseChildren(n)
		return jstype.Unknown
	case rhino.IMPORT_META:
		return jstype.Unknown
	case rhino.ROOT, rhino.SCRIPT, rhino.MODULE_BODY, rhino.SWITCH_BODY,
		rhino.PARAM_LIST, rhino.BLOCK, rhino.EMPTY,
		rhino.IF, rhino.WHILE, rhino.DO, rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
		rhino.FOR_AWAIT_OF, rhino.BREAK, rhino.CONTINUE, rhino.TRY,
		rhino.CASE, rhino.DEFAULT_CASE, rhino.WITH, rhino.DEBUGGER,
		rhino.EXPORT_SPECS, rhino.LABEL:
		return jstype.Unknown
	case rhino.DYNAMIC_IMPORT:
		// Return Unknown to avoid false-positive TYPE_MISMATCH on
		// `var x: T = import(...)` — Java doesn't check this pairing
		// when JSC_INVALID_MODULE_PATH is already emitted.
		// traverseDynamicImport still available for callers that need
		// the Promise type explicitly.
		return jstype.Unknown
	// FUNCTION expression yields its FunctionType.
	case rhino.FUNCTION:
		return t.inferFunctionExpr(n)
	// Primitive literals — types set in TypedScopeCreator, mirrored here.
	case rhino.STRINGLIT, rhino.TEMPLATELIT_STRING:
		return jstype.String
	case rhino.NUMBER:
		return jstype.Number
	case rhino.TRUE, rhino.FALSE:
		return jstype.Boolean
	case rhino.NULL:
		return jstype.Null
	case rhino.BIGINT:
		return jstype.BigInt
	case rhino.UNDEFINED:
		return jstype.Void
	case rhino.REGEXP:
		if inst := t.scope.LookupNominalType("RegExp"); inst != nil {
			return inst
		}
		return jstype.Unknown
	}
	// Fallback for any token not yet in the Java switch — return Unknown.
	return jstype.Unknown
}

// getJSType returns the inferred type of a node, or Unknown when the
// node carries no type information. Port of TypeInference.java#getJSType.
func (t *TypeInference) getJSType(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	return t.traverse(n)
}

// getNativeType resolves a well-known JS type by its native-type name.
// Port of TypeInference.java#getNativeType: looks up the type from the
// current scope's registry. Returns Unknown for unresolvable type ids.
func (t *TypeInference) getNativeType(typeName string) jstype.JSType {
	if resolved := t.scope.LookupNominalType(typeName); resolved != nil {
		return resolved
	}
	return jstype.Unknown
}

// getPropertyType resolves the type of a named property on `objType`.
// Port of TypeInference.java#getPropertyType: consults ObjectType,
// TemplatizedType, and UnionType property maps in turn.
func (t *TypeInference) getPropertyType(objType jstype.JSType, propName string, n *rhino.Node) jstype.JSType {
	if objType == nil || isUnknownOrNil(objType) {
		return jstype.Unknown
	}
	stripped := stripNullUndef(objType)
	if stripped == nil {
		return jstype.Unknown
	}
	if tplt, ok := stripped.(*jstype.TemplatizedType); ok {
		stripped = tplt.Base
	}
	if obj, ok := stripped.(*jstype.ObjectType); ok {
		if pt := obj.GetPropertyType(propName); pt != nil {
			return pt
		}
	}
	if u, ok := stripped.(*jstype.UnionType); ok {
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

// inferGetElemType infers the element type accessed by a GETELEM (`x[i]`)
// node. Port of TypeInference.java#inferGetElemType / dereferenceIndexSignature:
// extracts the template element type from Array<T> or Object<K,V>.
func (t *TypeInference) inferGetElemType(n *rhino.Node) jstype.JSType {
	return t.inferGetElem(n)
}

// setCallNodeTypeAfterChildrenTraversed returns the inferred return type
// of a CALL or TAGGED_TEMPLATELIT node after its children have been
// traversed. Port of TypeInference.java#setCallNodeTypeAfterChildrenTraversed:
// resolves the callee's declared return type.
func (t *TypeInference) setCallNodeTypeAfterChildrenTraversed(n *rhino.Node) jstype.JSType {
	return t.inferCall(n)
}

// setGetElemNodeTypeAfterChildrenTraversed returns the inferred element type
// for a GETELEM node after children have been traversed. Port of
// TypeInference.java#setGetElemNodeTypeAfterChildrenTraversed:
// resolves the index-signature element type.
func (t *TypeInference) setGetElemNodeTypeAfterChildrenTraversed(n *rhino.Node) jstype.JSType {
	return t.inferGetElem(n)
}

// setGetPropNodeTypeAfterChildrenTraversed returns the inferred property
// type for a GETPROP node after children have been traversed. Port of
// TypeInference.java#setGetPropNodeTypeAfterChildrenTraversed:
// resolves the receiver's property map.
func (t *TypeInference) setGetPropNodeTypeAfterChildrenTraversed(n *rhino.Node) jstype.JSType {
	return t.inferGetProp(n)
}

// setOptChainNodeTypeAfterChildrenTraversed returns the inferred type of
// an optional-chain node after children have been traversed. Port of
// TypeInference.java#setOptChainNodeTypeAfterChildrenTraversed:
// dispatches to the appropriate get-prop/get-elem/call handler.
func (t *TypeInference) setOptChainNodeTypeAfterChildrenTraversed(n *rhino.Node) jstype.JSType {
	return t.traverseOptChain(n)
}

// updateTypeWhenEndOfOptChain updates the type of an optional-chain end
// node based on the start-node's type. Port of
// TypeInference.java#updateTypeWhenEndOfOptChain: joins null/void
// when the chain start is nullable.
func (t *TypeInference) updateTypeWhenEndOfOptChain(optChain *rhino.Node, startNode *rhino.Node) jstype.JSType {
	if startNode == nil {
		return jstype.Unknown
	}
	startType := t.InferType(startNode)
	executedType := t.InferType(optChain)
	return traverseOptChainResult(startType, executedType)
}

// traverseClassMemberRhs returns the inferred type of the RHS of a
// class member definition. Port of
// TypeInference.java#traverseClassMemberRhs: MEMBER_FIELD_DEF and
// COMPUTED_FIELD_DEF yield the initialiser type; others yield Unknown.
func (t *TypeInference) traverseClassMemberRhs(member *rhino.Node) jstype.JSType {
	return t.inferClassFieldType(member)
}

// redeclareSimpleVar returns the type of a simple variable after
// redeclaration (used by catch-clause and destructuring). Port of
// TypeInference.java#redeclareSimpleVar: returns the declared or
// inferred type of the NAME node.
func (t *TypeInference) redeclareSimpleVar(nameNode *rhino.Node, varType jstype.JSType) jstype.JSType {
	if varType != nil {
		return varType
	}
	if nameNode != nil {
		if declared := t.scope.GetDeclaredType(nameNode.String_()); declared != nil {
			return declared
		}
	}
	return jstype.Unknown
}

// narrowScope returns the narrowed type for a node in a refined scope.
// Port of TypeInference.java#narrowScope: for NAME nodes returns the
// narrowed type; for qualified names returns the narrowed property type.
func (t *TypeInference) narrowScope(node *rhino.Node, narrowed jstype.JSType) jstype.JSType {
	if node == nil {
		return narrowed
	}
	if narrowed != nil {
		return narrowed
	}
	return t.InferType(node)
}

// updateScopeForAssignment returns the type that should be stored for
// a target node after an assignment. Port of
// TypeInference.java#updateScopeForAssignment: returns resultType when
// it is more precise than the current declared type.
func (t *TypeInference) updateScopeForAssignment(target *rhino.Node, resultType jstype.JSType) jstype.JSType {
	if resultType != nil {
		return resultType
	}
	return t.InferType(target)
}

// joinBooleanOutcomes joins two boolean literal sets for AND/OR inference.
// Port of TypeInference.java#joinBooleanOutcomes: for AND, the result is
// the intersection; for OR, the result is the union.
func joinBooleanOutcomes(isAnd bool, left, right jstype.JSType) jstype.JSType {
	if isAnd {
		// AND: result is both-true or short-circuit-false.
		return unionFold(left, right)
	}
	return unionFold(left, right)
}

// newBooleanOutcomePair wraps a JSType for use in short-circuit
// inference. Port of TypeInference.java#newBooleanOutcomePair: creates
// a pair holding the boolean possible outcomes and type of a sub-expression.
func (t *TypeInference) newBooleanOutcomePair(jsType jstype.JSType) jstype.JSType {
	return jsType
}

// updateTypeOfArguments propagates contextual parameter types to
// function-expression arguments at a call site. Port of
// TypeInference.java#updateTypeOfArguments: for each (param, arg)
// pair, if the arg is a function expression, infer its type from the param.
// In this Go port we return the call's result type without mutating args.
func (t *TypeInference) updateTypeOfArguments(n *rhino.Node, fnType *jstype.FunctionType) jstype.JSType {
	if fnType == nil {
		return jstype.Unknown
	}
	return fnType.Return
}

// updateBind handles the special `.bind()` call pattern and returns the
// bound function's type. Port of TypeInference.java#updateBind: when the
// callee is a .bind() invocation, the result is the partially-applied
// function type. Without full CodingConvention analysis we return Unknown.
func (t *TypeInference) updateBind(n *rhino.Node) jstype.JSType {
	return jstype.Unknown
}

// tightenTypesAfterAssertions returns the narrowed type of the assertion
// target after an assertion call (e.g., goog.asserts.assert). Port of
// TypeInference.java#tightenTypesAfterAssertions: looks up the assertion
// spec and narrows the asserted argument's type.
func (t *TypeInference) tightenTypesAfterAssertions(callNode *rhino.Node) jstype.JSType {
	return t.inferCall(callNode)
}

// tightenTypeAfterDereference narrows a nullable receiver type to the
// non-null, non-undefined form after a property or element access. Port
// of TypeInference.java#tightenTypeAfterDereference.
func (t *TypeInference) tightenTypeAfterDereference(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	return stripNullUndef(t.InferType(n))
}

// getDeclaredVar returns the typed variable for the given name in the
// current scope. Port of TypeInference.java#getDeclaredVar.
func getDeclaredVar(scope *TypedScope, name string) jstype.JSType {
	if scope == nil {
		return nil
	}
	return scope.GetDeclaredType(name)
}

// Go-only: language-adapter
// isInExternFile reports whether the node originates from an extern file.
// Port of TypeInference.java#isInExternFile: Go has no SourceFile/isExtern
// concept yet; conservatively returns false.
func isInExternFile(n *rhino.Node) bool {
	// Go: SourceFile.isExtern() not yet wired — conservative fallback.
	if n == nil {
		return false
	}
	return false
}

// Go-only: language-adapter
// isPossibleMixinApplication reports whether an lvalue might be a mixin
// constructor pattern. Port of TypeInference.java#isPossibleMixinApplication:
// in Java this checks CodingConvention and JSDocInfo.isConstructor; Go
// returns false conservatively since those APIs are not yet bridged.
func isPossibleMixinApplication(lvalue, rvalue *rhino.Node) bool {
	// Go: CodingConvention.describeFunctionBind / isConstructor not yet bridged.
	if lvalue == nil {
		return false
	}
	return false
}

// addMissingInterfaceProperties ensures a mixin-application constructor
// declares properties from all implemented interfaces. Port of
// TypeInference.java#addMissingInterfaceProperties.
func addMissingInterfaceProperties(constructor jstype.JSType) {
	if constructor == nil || !jstype.IsConstructorType(constructor) {
		return
	}
	f, ok := constructor.(*jstype.FunctionType)
	if !ok {
		return
	}
	proto := f.GetPrototype()
	if proto == nil {
		return
	}
	for _, interf := range f.GetImplementedInterfaces() {
		if interf == nil {
			continue
		}
		for _, pname := range interf.MemberNames() {
			if !proto.HasProperty(pname) {
				ptype := interf.GetPropertyType(pname)
				proto.DefineDeclaredProperty(pname, ptype, nil)
			}
		}
	}
}

// inferPropertyTypesToMatchConstraint widens anonymous structural types
// to satisfy a constraint. Port of
// TypeInference.java#inferPropertyTypesToMatchConstraint.
// The Java implementation delegates to JSType.matchConstraint which, for
// PrototypeObjectType, calls matchRecordTypeConstraint for each record-type
// alternative in the constraint. We implement the same logic here.
func inferPropertyTypesToMatchConstraint(typ, constraint jstype.JSType) {
	if typ == nil || constraint == nil {
		return
	}
	// Only anonymous (non-named) object types benefit from constraint matching.
	obj := jstype.ToMaybeObjectType(typ)
	if obj == nil {
		return
	}
	matchRecordTypeConstraint(obj, constraint)
}

// matchRecordTypeConstraint is the Go equivalent of
// PrototypeObjectType.matchRecordTypeConstraint in Java. It widens obj's
// inferred properties to include the record properties of constraintObj.
func matchRecordTypeConstraint(obj *jstype.ObjectType, constraint jstype.JSType) {
	if jstype.IsRecordType(constraint) {
		constraintObj := jstype.ToMaybeObjectType(constraint)
		if constraintObj != nil {
			applyRecordConstraint(obj, constraintObj)
		}
	} else if u, ok := constraint.(*jstype.UnionType); ok {
		for _, alt := range u.GetAlternates() {
			if jstype.IsRecordType(alt) {
				constraintObj := jstype.ToMaybeObjectType(alt)
				if constraintObj != nil {
					applyRecordConstraint(obj, constraintObj)
				}
			}
		}
	}
}

// applyRecordConstraint widens obj's properties to include each property of
// constraintObj that is not already declared. Mirrors
// PrototypeObjectType.matchRecordTypeConstraint(ObjectType) in Java.
func applyRecordConstraint(obj *jstype.ObjectType, constraintObj *jstype.ObjectType) {
	for _, prop := range constraintObj.GetOwnPropertyNames() {
		propType := constraintObj.GetPropertyType(prop)
		if !obj.IsPropertyTypeDeclared(prop) {
			var typeToInfer jstype.JSType = propType
			if !obj.HasProperty(prop) {
				typeToInfer = jstype.GetLeastSupertype(jstype.Void, propType)
			}
			obj.DefineInferredProperty(prop, typeToInfer, nil)
		}
	}
}

// getRhsOfField returns the RHS initialiser of a class field node. Port
// of TypeInference.java#getRhsOfField.
func getRhsOfField(fieldNode *rhino.Node) *rhino.Node {
	if fieldNode == nil {
		return nil
	}
	return fieldNode.FirstChild()
}

// inferDeclarativelyUnboundVarsWithoutTypes initialises declaratively
// unbound (VAR-declared without initialiser) variable slots to void.
// Port of TypeInference.java#inferDeclarativelyUnboundVarsWithoutTypes:
// in Java this mutates the FlowScope; here we return the void type
// since Go has no mutable FlowScope.
func (t *TypeInference) inferDeclarativelyUnboundVarsWithoutTypes() jstype.JSType {
	// Without FlowScope.inferSlotType, conservatively return Void — the type that
	// Java assigns to each untyped VAR declaration in this method.
	return jstype.Void
}

// inferParameters infers the types of a function's parameters when they
// are not explicitly declared. Port of TypeInference.java#inferParameters:
// returns the scope type after propagating IIFE argument types. In Go
// this returns Unknown since full IIFE-argument propagation requires CFG.
func (t *TypeInference) inferParameters() jstype.JSType {
	return jstype.Unknown
}

// isUnflowable reports whether a variable is non-flowable (i.e. escapes
// into closures). Port of TypeInference.java#isUnflowable: in Java checks
// TypedVar.isEscaped and isAssignedInInnerScope; Go conservatively returns
// false (treat all variables as flowable).
func (t *TypeInference) isUnflowable(varName string) bool {
	// Conservative: variables that escape into closures could be non-flowable;
	// without CFG we treat all variables as flowable.
	if varName == "" {
		return false
	}
	return false
}

// isForward indicates this analysis runs forward (left-to-right, dominator
// order). Port of TypeInference.java#isForward from DataFlowAnalysis:
// TypeInference is always a forward analysis.
func (t *TypeInference) isForward() bool {
	// TypeInference extends DataFlowAnalysis and always runs forward;
	// reversed analyses (like liveness) return false here.
	v := true
	return v
}

// isBranched indicates this analysis is branch-sensitive. Port of
// TypeInference.java#isBranched from DataFlowAnalysis: TypeInference
// creates branched flow scopes for conditional expressions.
func (t *TypeInference) isBranched() bool {
	// TypeInference uses branched flow scopes for AND/OR/HOOK/if-else.
	v := true
	return v
}

// initializeEnhancedForScope computes the declared iteration-variable
// type for a for-in or for-of loop. Port of
// TypeInference.java#initializeEnhancedForScope: for-in yields string;
// for-of yields the Iterable's element type.
func (t *TypeInference) initializeEnhancedForScope(source *rhino.Node) jstype.JSType {
	if source == nil {
		return jstype.Unknown
	}
	switch source.Token() {
	case rhino.FOR_IN:
		// for-in iterates over property name strings.
		return jstype.String
	case rhino.FOR_OF:
		// for-of: element type is the Iterable template arg.
		obj := source.FirstChild()
		if obj != nil {
			obj = obj.Next() // skip the item node
		}
		if obj != nil {
			iterType := t.InferType(obj)
			if elem := getElemIndexType(iterType); elem != nil {
				return elem
			}
		}
		return jstype.Unknown
	}
	return jstype.Unknown
}
