// Port of com/google/javascript/rhino/jstype/FunctionType.java @ closure-compiler v20260415
//
// FunctionType in Java is a rich node: it carries parameters, return
// type, this-type, supers, template info, constructor/interface flag,
// and more. Tier-1 only exercises the ordered parameter list + return
// type, so that's all we record here; the Java class's extra fields
// will be added as later tiers require them.

package jstype

import "strings"

type FunctionType struct {
	// Ordered positional parameters. Unknown params (e.g. JSDoc left a
	// name unannotated) carry the Unknown JSType so call-site checking
	// stays conservative.
	Params []JSType
	Return JSType // Unknown when @return was absent / unrecognised.

	// MinArity is the minimum number of arguments callers must pass.
	// Equals len(Params) when every param is required; decreases when
	// closure-convention `opt_<name>` or `var_args` prefixes mark
	// params optional. Mirrors FunctionType.getMinArity() in Java.
	MinArity int
	// MaxArity is the maximum number of arguments callers may pass.
	// Equals len(Params) normally; MaxArgsUnlimited (math.MaxInt32)
	// when the function has a rest param (`var_args`). Mirrors
	// FunctionType.getMaxArity().
	MaxArity int

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
}

// MaxArgsUnlimited is the sentinel MaxArity for rest-parameter
// functions — callers may pass arbitrarily many arguments.
const MaxArgsUnlimited = 1 << 30

func NewFunctionType(params []JSType, ret JSType) *FunctionType {
	if ret == nil {
		ret = Unknown
	}
	return &FunctionType{
		Params:   params,
		Return:   ret,
		MinArity: len(params),
		MaxArity: len(params),
	}
}

// NewConstructorType is the variant used when @constructor is present.
// Returns a FunctionType whose Instance is a freshly-named ObjectType
// — the type assigned to `new ThisCtor(args)` expressions.
func NewConstructorType(name string, params []JSType) *FunctionType {
	inst := NewObjectType(name)
	return &FunctionType{
		Params:        params,
		Return:        inst,
		MinArity:      len(params),
		MaxArity:      len(params),
		IsConstructor: true,
		Instance:      inst,
	}
}

func (f *FunctionType) TypeName() string {
	var names []string
	for _, p := range f.Params {
		names = append(names, p.TypeName())
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
