// Port of com/google/javascript/rhino/jstype/FunctionType.java @ closure-compiler v20260415
//
// FunctionType in Java is a rich node: it carries parameters, return
// type, this-type, supers, template info, constructor/interface flag,
// and more. We now carry the full Parameter shape (optional/rest/
// default bits) matching the Java inner-class, along with TypeOfThis
// and Kind for constructor/interface/ordinary classification.
// The older flat Params []JSType / Return JSType fields are kept as
// shim accessors so that existing callers continue to compile.

package jstype

import "strings"

// FunctionKind mirrors FunctionType.Kind in Java.
// Determines whether this is a constructor, interface, or plain function.
type FunctionKind int

const (
	// FunctionKindOrdinary is a plain, non-constructor function.
	// Mirrors Java Kind.ORDINARY.
	FunctionKindOrdinary FunctionKind = iota
	// FunctionKindConstructor is a @constructor function.
	// Mirrors Java Kind.CONSTRUCTOR.
	FunctionKindConstructor
	// FunctionKindInterface is a @interface / @record function.
	// Mirrors Java Kind.INTERFACE.
	FunctionKindInterface
	// FunctionKindNone is the sentinel used for the synthetic NONE
	// kind in Java (e.g. the function type itself). Unused in Tier-1.
	FunctionKindNone
)

// Parameter models a single JavaScript function parameter.
//
// Port of the inner class FunctionType.Parameter (an @AutoValue) in
// the Java source. Carries the JSType together with the three boolean
// flags that distinguish the four parameter flavours:
//
//	required   — IsOptional=false, IsVariadic=false
//	optional   — IsOptional=true,  IsVariadic=false  (JSDoc {=T})
//	rest/var   — IsOptional=false, IsVariadic=true   (JSDoc {...T})
//	default    — IsOptional=true + non-nil DefaultValue            (JSDoc default-value)
//
// Java's AutoValue generates immutable value-semantics; in Go we use
// a plain struct. Create via NewParameter; do not take the address of
// a literal unless you know what you're doing.
type Parameter struct {
	// jsType is the declared type of the parameter.
	jsType JSType
	// isOptional mirrors Parameter.isOptional(). True when the
	// parameter may be omitted at the call site (JSDoc {=T} or
	// closure-convention opt_name prefix).
	isOptional bool
	// isVariadic mirrors Parameter.isVariadic(). True when the
	// parameter is a rest/var-args param (JSDoc {...T} or
	// closure-convention var_args name). At most one variadic
	// parameter per function, and it must be the last one.
	isVariadic bool
	// hasDefaultValue records whether the parameter was declared
	// with an ES2015 default-value expression. Java does not have
	// this field directly on Parameter but the builder tracks it;
	// we store it here to keep the data together.
	// TODO(go-port): Java FunctionParamBuilder does not expose a
	// hasDefaultValue bit on Parameter itself; this field is a
	// Go-side addition and may need revisiting when we port the
	// full type-inference logic.
	hasDefaultValue bool
}

// NewParameter creates a Parameter with the given type and flags.
// Mirrors Parameter.create(JSType, boolean, boolean) in Java.
func NewParameter(t JSType, isOptional, isVariadic bool) Parameter {
	return Parameter{jsType: t, isOptional: isOptional, isVariadic: isVariadic}
}

// NewParameterWithDefault creates a Parameter that has a default value.
// Callers should also pass isOptional=true (a param with a default is
// always optional from the caller's perspective).
func NewParameterWithDefault(t JSType, isOptional, isVariadic bool) Parameter {
	return Parameter{jsType: t, isOptional: isOptional, isVariadic: isVariadic, hasDefaultValue: true}
}

// GetJSType mirrors Parameter.getJSType() in Java.
func (p Parameter) GetJSType() JSType { return p.jsType }

// IsOptional mirrors Parameter.isOptional() in Java.
func (p Parameter) IsOptional() bool { return p.isOptional }

// IsVariadic mirrors Parameter.isVariadic() in Java.
// True for rest / var-args params ({...T} in JSDoc).
func (p Parameter) IsVariadic() bool { return p.isVariadic }

// HasDefaultValue reports whether the parameter was declared with an
// ES2015 default-value expression.
func (p Parameter) HasDefaultValue() bool { return p.hasDefaultValue }

// FunctionType represents the type of a JavaScript function value.
//
// Port of com/google/javascript/rhino/jstype/FunctionType.java.
// This is the primary bearer of per-function type information: its
// parameter list (with optional/rest/default bits), return type,
// this-type, constructor/interface classification, and template vars.
//
// Backward-compatible shim fields (Params, Return) are preserved so
// that existing callers continue to compile; new code should prefer
// the accessor methods (GetParameters, GetReturnType, etc.).
type FunctionType struct {
	// parameters is the ordered list of typed parameters with their
	// optional/rest/default flags. Mirrors the parameterList stored
	// on ArrowType in Java (accessed via getParameters()).
	parameters []Parameter

	// returnType is the declared or inferred return type.
	// Mirrors ArrowType.returnType in Java (accessed via getReturnType()).
	returnType JSType

	// returnTypeInferred is true when the return type was inferred
	// rather than explicitly declared. Mirrors
	// ArrowType.returnTypeInferred in Java.
	returnTypeInferred bool

	// typeOfThis is the type of `this` inside the function body.
	// Mirrors FunctionType.typeOfThis in Java (getTypeOfThis()).
	// For constructors this is the instance type; for ordinary
	// functions it is Unknown by default.
	typeOfThis JSType

	// kind classifies the function as ordinary / constructor /
	// interface / none. Mirrors FunctionType.kind in Java.
	kind FunctionKind

	// --- Shim fields kept for backward compatibility ---
	// These are the flat representation that existing callers still
	// read directly. They are kept in sync with parameters/returnType
	// by the constructors; callers that use the new accessor methods
	// do not need them.

	// Params is the flat ordered slice of parameter types. Existing
	// callers that access fn.Params[i] directly continue to work.
	// New code should prefer GetParameters().
	Params []JSType

	// Return is the return type. Existing callers that access
	// fn.Return directly continue to work. New code should prefer
	// GetReturnType().
	Return JSType

	// MinArity is the minimum number of required arguments.
	// Mirrors FunctionType.getMinArity() in Java.
	// Kept as a field so existing callers (fn.MinArity) compile.
	MinArity int

	// MaxArity is the maximum number of accepted arguments.
	// MaxArgsUnlimited when a variadic param is present.
	// Mirrors FunctionType.getMaxArity() in Java.
	// Kept as a field so existing callers (fn.MaxArity) compile.
	MaxArity int

	// IsConstructor is true when the function was annotated
	// @constructor. Kept as a field for existing callers.
	// New code should prefer the IsConstructor() method.
	IsConstructor bool

	// Instance is the ObjectType created for `new ThisCtor(args)`
	// when IsConstructor is true. Unset for plain functions.
	Instance *ObjectType

	// TemplateTypes are the function's @template parameters in
	// declaration order. Empty for non-generic functions.
	TemplateTypes []*TemplateType

	// IsES6Class is set when this FunctionType came from an ES6
	// `class` declaration.
	IsES6Class bool
}

// MaxArgsUnlimited is the sentinel MaxArity for rest-parameter
// functions — callers may pass arbitrarily many arguments.
// Mirrors Integer.MAX_VALUE returned by FunctionType.getMaxArity()
// in Java when the last parameter isVariadic.
const MaxArgsUnlimited = 1 << 30

// computeArity derives MinArity and MaxArity from the parameter list.
// This mirrors the logic in FunctionType.getMinArity() /
// getMaxArity() in Java.
func computeArity(params []Parameter) (min, max int) {
	min = 0
	for i, p := range params {
		if !p.isOptional && !p.isVariadic {
			min = i + 1
		}
	}
	if len(params) > 0 && params[len(params)-1].isVariadic {
		max = MaxArgsUnlimited
	} else {
		max = len(params)
	}
	return min, max
}

// flatTypes extracts the JSType slice from a Parameter slice for the
// shim Params field.
func flatTypes(params []Parameter) []JSType {
	if len(params) == 0 {
		return nil
	}
	out := make([]JSType, len(params))
	for i, p := range params {
		out[i] = p.jsType
	}
	return out
}

// newFunctionTypeFromParams is the internal constructor shared by the
// public factory functions.
func newFunctionTypeFromParams(params []Parameter, ret JSType, kind FunctionKind, typeOfThis JSType) *FunctionType {
	if ret == nil {
		ret = Unknown
	}
	minArity, maxArity := computeArity(params)
	f := &FunctionType{
		parameters: params,
		returnType: ret,
		typeOfThis: typeOfThis,
		kind:       kind,
		Params:     flatTypes(params),
		Return:     ret,
		MinArity:   minArity,
		MaxArity:   maxArity,
		IsConstructor: kind == FunctionKindConstructor,
	}
	return f
}

// NewFunctionType creates a plain (non-constructor) FunctionType from
// a flat []JSType slice. This is the backward-compatible factory used
// by existing callers that have not yet been migrated to Parameters.
//
// Mirrors FunctionType.builder(registry)…withParameters(…).build()
// for the ORDINARY kind.
func NewFunctionType(params []JSType, ret JSType) *FunctionType {
	ps := make([]Parameter, len(params))
	for i, t := range params {
		ps[i] = NewParameter(t, false, false)
	}
	return newFunctionTypeFromParams(ps, ret, FunctionKindOrdinary, nil)
}

// NewFunctionTypeFromParameters creates a plain FunctionType from a
// rich []Parameter slice. This is the new factory preferred by code
// that needs to distinguish optional/rest/default params.
func NewFunctionTypeFromParameters(params []Parameter, ret JSType) *FunctionType {
	return newFunctionTypeFromParams(params, ret, FunctionKindOrdinary, nil)
}

// NewConstructorType is the variant used when @constructor is present.
// Returns a FunctionType whose Instance is a freshly-named ObjectType
// — the type assigned to `new ThisCtor(args)` expressions.
//
// Mirrors FunctionType.builder(registry).forConstructor().build().
func NewConstructorType(name string, params []JSType) *FunctionType {
	inst := NewObjectType(name)
	ps := make([]Parameter, len(params))
	for i, t := range params {
		ps[i] = NewParameter(t, false, false)
	}
	minArity, maxArity := computeArity(ps)
	return &FunctionType{
		parameters:    ps,
		returnType:    inst,
		typeOfThis:    inst,
		kind:          FunctionKindConstructor,
		Params:        flatTypes(ps),
		Return:        inst,
		MinArity:      minArity,
		MaxArity:      maxArity,
		IsConstructor: true,
		Instance:      inst,
	}
}

// NewConstructorTypeFromParameters is the rich variant of
// NewConstructorType that accepts a []Parameter slice.
func NewConstructorTypeFromParameters(name string, params []Parameter) *FunctionType {
	inst := NewObjectType(name)
	minArity, maxArity := computeArity(params)
	return &FunctionType{
		parameters:    params,
		returnType:    inst,
		typeOfThis:    inst,
		kind:          FunctionKindConstructor,
		Params:        flatTypes(params),
		Return:        inst,
		MinArity:      minArity,
		MaxArity:      maxArity,
		IsConstructor: true,
		Instance:      inst,
	}
}

// --- Java-equivalent accessors (new API) --------------------------------

// GetParameters mirrors FunctionType.getParameters() in Java.
// Returns the ordered list of Parameters with optional/rest/default
// bits. Do NOT mutate the returned slice.
func (f *FunctionType) GetParameters() []Parameter {
	return f.parameters
}

// GetReturnType mirrors FunctionType.getReturnType() in Java.
func (f *FunctionType) GetReturnType() JSType {
	return f.returnType
}

// IsReturnTypeInferred mirrors FunctionType.isReturnTypeInferred().
func (f *FunctionType) IsReturnTypeInferred() bool {
	return f.returnTypeInferred
}

// GetTypeOfThis mirrors FunctionType.getTypeOfThis() in Java.
// For constructors returns the instance type; for ordinary functions
// returns nil (callers should treat nil as Unknown).
func (f *FunctionType) GetTypeOfThis() JSType {
	return f.typeOfThis
}

// GetKind mirrors FunctionType.getKind() in Java.
func (f *FunctionType) GetKind() FunctionKind {
	return f.kind
}

// IsOrdinaryFunction mirrors FunctionType.isOrdinaryFunction() in Java.
func (f *FunctionType) IsOrdinaryFunction() bool {
	return f.kind == FunctionKindOrdinary
}

// IsInterface mirrors FunctionType.isInterface() in Java.
func (f *FunctionType) IsInterface() bool {
	return f.kind == FunctionKindInterface
}

// GetMinArguments mirrors FunctionType.getMinArity() in Java.
// Returns the position of the last required (non-optional,
// non-variadic) parameter.
func (f *FunctionType) GetMinArguments() int {
	min := 0
	for i, p := range f.parameters {
		if !p.isOptional && !p.isVariadic {
			min = i + 1
		}
	}
	return min
}

// GetMaxArguments mirrors FunctionType.getMaxArity() in Java.
// Returns MaxArgsUnlimited when the last parameter is variadic.
func (f *FunctionType) GetMaxArguments() int {
	if len(f.parameters) > 0 && f.parameters[len(f.parameters)-1].isVariadic {
		return MaxArgsUnlimited
	}
	return len(f.parameters)
}

// HasRestParam reports whether this function has a variadic (rest)
// parameter. Equivalent to checking getMaxArity()==MAX_VALUE in Java.
func (f *FunctionType) HasRestParam() bool {
	return len(f.parameters) > 0 && f.parameters[len(f.parameters)-1].isVariadic
}

// --- TypeName / subtype ---------------------------------------------------

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
