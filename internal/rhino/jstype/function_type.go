// Port of com/google/javascript/rhino/jstype/FunctionType.java @ closure-compiler v20260415
//
// FunctionType in Java is a rich node: it carries parameters, return
// type, this-type, supers, template info, constructor/interface flag,
// and more. Tier-1 only exercises the ordered parameter list + return
// type, so that's all that used to be here; this port fills in the
// extra fields jsgo now needs to model optional / rest params, the
// THIS-binding slot, a closure-capture stub, and the three function
// "kinds" Closure distinguishes (named, anonymous, native).

package jstype

import "strings"

// FunctionKind mirrors FunctionType.Kind in Java: NAMED is a declared
// function (function Foo(){}), ANONYMOUS is a function expression
// without a name (var f = function(){}), NATIVE is a host-provided
// built-in that has no AST body (Array, Function itself, etc.).
type FunctionKind int

const (
	KindNamed     FunctionKind = iota // function Foo(){}
	KindAnonymous                     // function(){} / () => ...
	KindNative                        // built-in; no source body
)

type FunctionType struct {
	// Name is the source-level identifier for named functions; empty
	// for anonymous and usually the class name for native built-ins.
	// Mirrors FunctionType.getName().
	Name string

	// Kind is one of KindNamed / KindAnonymous / KindNative. Picked
	// by NewFunctionType / NewAnonFunctionType / NewNativeFunctionType.
	Kind FunctionKind

	// Ordered positional parameters. Unknown params (e.g. JSDoc left a
	// name unannotated) carry the Unknown JSType so call-site checking
	// stays conservative.
	Params []JSType

	// OptionalParams is a bit-parallel slice with Params: true means
	// the param was declared with a trailing `=` in the JSDoc type
	// expression (e.g. `@param {string=} x`). Used by call-site arity
	// checking to allow a shorter argument list. Mirrors
	// FunctionType.isOptionalParameter.
	OptionalParams []bool

	// RestParam, when >= 0, is the index of the rest (`...`) parameter.
	// The slot in Params at that index is the element type; callers
	// may pass zero or more values of that type at this position.
	// Mirrors FunctionType.hasRestParameter / getRestFormalsType.
	// -1 means no rest parameter.
	RestParam int

	Return JSType // Unknown when @return was absent / unrecognised.

	// This is the THIS-binding type, if any. Set from
	// `@this {Foo}` JSDoc or the enclosing class scope. nil when no
	// binding applies (free function on module-top). Mirrors
	// FunctionType.getTypeOfThis.
	This JSType

	// IsConstructor mirrors FunctionType.isConstructor() — set when the
	// owning function carried `@constructor`. Influences both call-site
	// shape checking and the nominal type produced by `new`.
	IsConstructor bool

	// Instance is the InstanceType created for `new ThisFunction()`
	// when IsConstructor is true. Unset for plain functions; the scope
	// creator wires this when it processes the @constructor JSDoc.
	Instance *ObjectType

	// TemplateTypes are the function's @template parameters, in the
	// order they were declared. Call-site unification uses these to
	// build a TemplateTypeMap from arg types before substituting into
	// Return. Empty for non-generic functions.
	TemplateTypes []*TemplateType

	// ClosureEnv is a stub for the lexical-scope capture set. Java's
	// FunctionType has no single analogue: closure semantics are
	// handled by the scope-chain, but jsgo carries this slot as a
	// forward-looking anchor so escape-analysis style passes can hang
	// off it without a schema change later. Intentionally unused in
	// the current call-site check path.
	ClosureEnv map[string]JSType
}

// NewFunctionType builds a plain (named) FunctionType from a param
// type list and return type. Matches FunctionType.forDeclaration in
// Java when called from a function declaration.
func NewFunctionType(params []JSType, ret JSType) *FunctionType {
	if ret == nil {
		ret = Unknown
	}
	return &FunctionType{
		Kind:      KindNamed,
		Params:    params,
		Return:    ret,
		RestParam: -1,
	}
}

// NewAnonFunctionType is the variant picked for function expressions
// without a name (the RHS of `var f = function(){}` or arrow-function
// bodies). Mirrors FunctionType.forExpression. Kept distinct from
// NewFunctionType so diagnostics can tell the two apart.
func NewAnonFunctionType(params []JSType, ret JSType) *FunctionType {
	f := NewFunctionType(params, ret)
	f.Kind = KindAnonymous
	return f
}

// NewNativeFunctionType is the variant for built-in host functions
// (Array, Function, Object …). Closure sets a different flag on these
// so subtype checks can treat them specially (e.g. any function is a
// subtype of the native `Function` type). Mirrors
// FunctionType.Builder.forNativeType.
func NewNativeFunctionType(name string, params []JSType, ret JSType) *FunctionType {
	f := NewFunctionType(params, ret)
	f.Kind = KindNative
	f.Name = name
	return f
}

// NewConstructorType is the variant used when @constructor is present.
// Returns a FunctionType whose Instance is a freshly-named ObjectType
// — the type assigned to `new ThisCtor(args)` expressions.
func NewConstructorType(name string, params []JSType) *FunctionType {
	inst := NewObjectType(name)
	return &FunctionType{
		Name:          name,
		Kind:          KindNamed,
		Params:        params,
		Return:        inst,
		RestParam:     -1,
		IsConstructor: true,
		Instance:      inst,
	}
}

// IsOptionalParam reports whether param i was declared optional
// (trailing `=` in the JSDoc). Safe for out-of-range / missing
// metadata: returns false. Mirrors FunctionType.isOptionalParameter.
func (f *FunctionType) IsOptionalParam(i int) bool {
	if i < 0 || i >= len(f.OptionalParams) {
		return false
	}
	return f.OptionalParams[i]
}

// HasRestParam reports whether the final formal is a `...` rest
// position. Mirrors FunctionType.hasRestFormals.
func (f *FunctionType) HasRestParam() bool { return f.RestParam >= 0 }

// MinArity is the smallest number of call-site arguments that
// satisfies this signature — params count minus the trailing run of
// optionals (and minus 1 for a rest position, which accepts zero
// elements). Mirrors FunctionType.getMinArity.
func (f *FunctionType) MinArity() int {
	n := len(f.Params)
	if f.HasRestParam() {
		n--
	}
	for n > 0 && f.IsOptionalParam(n-1) {
		n--
	}
	if n < 0 {
		return 0
	}
	return n
}

// MaxArity is the largest accepted arg count, or -1 for "unbounded"
// when a rest param is present. Mirrors FunctionType.getMaxArity.
func (f *FunctionType) MaxArity() int {
	if f.HasRestParam() {
		return -1
	}
	return len(f.Params)
}

func (f *FunctionType) TypeName() string {
	var names []string
	for i, p := range f.Params {
		nm := p.TypeName()
		if f.HasRestParam() && i == f.RestParam {
			nm = "..." + nm
		} else if f.IsOptionalParam(i) {
			nm = nm + "="
		}
		names = append(names, nm)
	}
	return "function(" + strings.Join(names, ", ") + "): " + f.Return.TypeName()
}

// IsSubtypeOf is deliberately permissive for Tier-1: two FunctionType
// values are equivalent only when every param and the return line up
// exactly, and otherwise we accept AllType/Unknown. A faithful variance
// implementation lands with the Tier-2 subtype work.
func (f *FunctionType) IsSubtypeOf(o JSType) bool {
	if g, ok := o.(*FunctionType); ok {
		if len(f.Params) != len(g.Params) {
			return false
		}
		for i := range f.Params {
			if !IsSubtype(f.Params[i], g.Params[i]) {
				return false
			}
		}
		return IsSubtype(f.Return, g.Return)
	}
	// Any FunctionType is assignment-compatible with the
	// built-in `Function` instance (Java's subtype relationship
	// for callable-vs-Function). Matches how
	// FunctionType.isSubtype returns true against the top
	// Function type.
	if obj, ok := o.(*ObjectType); ok && obj.Name == "Function" {
		return true
	}
	return isAll(o) || isUnknown(o)
}
