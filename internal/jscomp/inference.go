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
	case rhino.SUPER:
		// `super` in a derived-class constructor is bound under the name
		// "super" in the ctor-local scope (see TypedScopeCreator.bindSuper).
		if declared := t.scope.GetDeclaredType("super"); declared != nil {
			return declared
		}
		return jstype.Unknown
	case rhino.GETPROP:
		return t.inferGetProp(n)
	case rhino.GETELEM:
		return t.inferGetElem(n)
	case rhino.NEW:
		if callee := n.FirstChild(); callee != nil {
			if ct := t.InferType(callee); ct != nil {
				if fn, ok := ct.(*jstype.FunctionType); ok && fn.Instance != nil {
					return fn.Instance
				}
			}
		}
		return jstype.Unknown
	case rhino.CALL:
		return t.inferCall(n)

	// ---- assignment --------------------------------------------- //
	case rhino.ASSIGN:
		// Bare `=` — the expression value is the RHS type. Matches
		// TypeInference.java's ASSIGN branch.
		kids := n.Children()
		if len(kids) >= 2 {
			return t.InferType(kids[1])
		}
		return jstype.Unknown
	case rhino.ASSIGN_ADD:
		// `x += y` has the same result type as `x + y`; delegate to
		// inferAdd so string-coercion is tracked. Mirrors how Java's
		// TypeInference sharpens compound-assign result types via the
		// underlying binary-op helper (getOpResultType in ADD path).
		return t.inferAddChildren(n)
	case rhino.ASSIGN_SUB, rhino.ASSIGN_MUL, rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		// Arithmetic / bitwise compound assigns always produce Number.
		// Port of the corresponding TypeInference.java branches.
		return jstype.Number

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
		// \`a && b\` returns \`b\` when \`a\` is truthy, else \`a\`. Result
		// is \`(a − truthy) ∪ b\` — the falsy members of \`a\` plus
		// whatever \`b\` evaluates to under the truthy-a refinement.
		kids := n.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		left := t.InferType(kids[0])
		refs := buildRefinements(t.scope, kids[0], true)
		right := t.inferUnder(kids[1], refs)
		falsy := keepFalsyRefined(left)
		if falsy == nil {
			return right
		}
		return unionFold(falsy, right)

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
		// Port of traverseArrayLiteral: returns Array<Unknown> so that
		// getElem / property lookup resolves through the Array prototype.
		// Java creates createTemplatizedType(ARRAY_TYPE, UNKNOWN_TYPE).
		if inst := t.scope.LookupNominalType("Array"); inst != nil {
			return jstype.NewTemplatizedType(inst, []jstype.JSType{jstype.Unknown})
		}
		return jstype.Unknown
	case rhino.OBJECTLIT:
		return t.inferObjectLit(n)
	case rhino.REGEXP:
		// traverseRegExp: a regex literal evaluates to a RegExp instance.
		// Port of TypeInference.java's REGEXP branch which sets the node's
		// type to the native REGEXP_TYPE.
		if inst := t.scope.LookupNominalType("RegExp"); inst != nil {
			return inst
		}
		return jstype.Unknown
	case rhino.TEMPLATELIT:
		// traverseTemplateLit: a template literal (non-tagged) always
		// evaluates to a string. Port of the TEMPLATELIT branch in
		// TypeInference.java which sets the node type to STRING_TYPE.
		return jstype.String
	case rhino.TAGGED_TEMPLATELIT:
		// traverseTaggedTemplateLit: a tagged template literal is a
		// special function call — the tag function is the first child and
		// the template literal is the second child. The expression type
		// equals the tag's return type, mirroring Java's
		// setCallNodeTypeAfterChildrenTraversed on TAGGED_TEMPLATELIT.
		return t.inferTaggedTemplateLit(n)

	// ---- increment / decrement ------------------------------------ //
	case rhino.INC, rhino.DEC:
		// traverseBigIntCompatibleUnaryOperator (INC/DEC): the result
		// type mirrors the operand — number for numeric operands. Since
		// BigInt is not yet tracked statically in this port, INC/DEC
		// conservatively return Number. Matches the NO_BIGINT branch of
		// Java's traverseBigIntCompatibleUnaryOperator.
		return jstype.Number

	// ---- yield / await ------------------------------------------- //
	case rhino.YIELD:
		// traverseYieldAll / traverseYield: the type of the yield
		// expression is the value sent back into the generator via
		// next(value). Without full dataflow we conservatively return
		// Unknown, matching Java's fallback when no TNext is resolved.
		return jstype.Unknown
	case rhino.AWAIT:
		// traverseAwait: resolves to the awaited type stripped of its
		// Promise wrapper. Mirrors Java's
		// Promises.getResolvedType(registry, exprType).
		return t.inferAwait(n)

	// ---- function expression -------------------------------------- //
	case rhino.FUNCTION:
		return t.inferFunctionExpr(n)

	// ---- cast expression ----------------------------------------- //
	// /** @type {T} */(expr) — the node's declared type is the cast
	// target resolved from its @type JSDoc annotation.
	case rhino.CAST:
		if info := n.JSDocInfo(); info != nil && info.HasType() {
			if ct := t.scope.ResolveExpr(info.Type().TypeName); ct != nil {
				return ct
			}
		}
		return jstype.Unknown
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

func isNullOrUndefinedLiteral(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.Token() == rhino.NULL || n.Token() == rhino.UNDEFINED
}

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
		// qualifiedName is defined in refinement.go; both files are in
		// the same package so it is accessible here.
		if qname := qualifiedName(target); qname != "" {
			return qname, strNode.String_(), equal
		}
	}
	return "", "", false
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

// inferObjectLit builds a structural ObjectType from the STRING_KEY
// entries of an OBJECTLIT node. Used so inline `{k: v, …}` values
// have a usable inferred shape for property lookups, mirroring how
// Java's TypeInference builds an anonymous PrototypeObjectType.
// If the OBJECTLIT carries a JSDoc with @struct or @dict, the
// corresponding flag is propagated so visitObjectLit can fire
// JSC_ILLEGAL_OBJLIT_KEY.
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
		kvInfo := kv.JSDocInfo()
		// @type {T} on the key: property is declared with type T.
		if kvInfo != nil && kvInfo.HasType() {
			if declared := t.scope.ResolveExpr(kvInfo.Type().TypeName); declared != nil {
				obj.DefineDeclaredProperty(propName, declared, kv)
				continue
			}
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

// traverseUnaryPlusOp infers the type of unary `+`. Port of
// TypeInference.java#traverseUnaryPlus: bigint operand → NO_TYPE in
// Java (we return Number conservatively).
func (t *TypeInference) traverseUnaryPlusOp(operand jstype.JSType) jstype.JSType {
	// unary-plus on bigint is a runtime error; conservatively Number.
	_ = getBigIntPresence(operand) // validate same as Java
	return jstype.Number
}

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

// inferCatchParam returns the type for a catch parameter variable.
// Port of TypeInference.java#traverseCatch: declared type if present,
// else Unknown (any value can be thrown).
func (t *TypeInference) inferCatchParam(catchVar string) jstype.JSType {
	if declared := t.scope.GetDeclaredType(catchVar); declared != nil {
		return declared
	}
	return jstype.Unknown
}

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
