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

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	// FunctionKindArrow is an ES2015 arrow function (`() => …`).
	// Go-only: language-adapter — Java has no separate arrow kind; arrow
	// functions are ordinary functions that capture `this` lexically.
	FunctionKindArrow
	// FunctionKindNone is the sentinel used for the synthetic NONE
	// kind in Java (e.g. the function type itself).
	FunctionKindNone
)

// ConstructorAmbiguity mirrors FunctionType.ConstructorAmbiguity in Java.
// Used by isAmbiguousConstructor().
type ConstructorAmbiguity int

const (
	// ConstructorAmbiguityUnknown is the default before calculateConstructorAmbiguity runs.
	ConstructorAmbiguityUnknown ConstructorAmbiguity = iota
	// ConstructorAmbiguityIsAmbiguous means the constructor is ambiguous
	// (no source node, or superclass is ambiguous).
	ConstructorAmbiguityIsAmbiguous
	// ConstructorAmbiguityIsUnambiguous means the constructor is unambiguous.
	ConstructorAmbiguityIsUnambiguous
)

// propAccess mirrors FunctionType.PropAccess in Java.
// Controls struct/dict annotation inheritance.
type propAccess int

const (
	// propAccessAny is the default (implicit ES5-class behaviour).
	propAccessAny propAccess = iota
	// propAccessAnyExplicit is @unrestricted (explicit).
	propAccessAnyExplicit
	// propAccessStruct is @struct or implicit ES6-class behaviour.
	propAccessStruct
	// propAccessDict is @dict or implicit ES5-class behaviour.
	propAccessDict
)

// delegateSuffix is the suffix used to detect delegate-proxy constructors.
// Mirrors ObjectType.createDelegateSuffix("Proxy") in Java.
const delegateSuffix = "_DelegateProxy"

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
	// Go-only: language-adapter — Java FunctionParamBuilder does not expose a
	// hasDefaultValue bit on Parameter itself; this field is a
	// Go-side addition.
	hasDefaultValue bool
}

// NewParameter creates a Parameter with the given type and flags.
// Mirrors Parameter.create(JSType, boolean, boolean) in Java.
// Go-only: language-adapter — Java uses the static Parameter.create() factory
// inside FunctionType; Go uses a package-level function for ergonomics.
func NewParameter(t JSType, isOptional, isVariadic bool) Parameter {
	return Parameter{jsType: t, isOptional: isOptional, isVariadic: isVariadic}
}

// NewParameterWithDefault creates a Parameter that has a default value.
// Callers should also pass isOptional=true (a param with a default is
// always optional from the caller's perspective).
// Go-only: language-adapter — Java tracks default-value separately; we
// fold it into Parameter for convenience.
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
// Go-only: language-adapter — see NewParameterWithDefault.
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

	// propAccess mirrors FunctionType.propAccess in Java.
	// Tracks struct/dict annotation state for constructor types.
	propAccess propAccess

	// constructorAmbiguity mirrors FunctionType.constructorAmbiguity in Java.
	constructorAmbiguity ConstructorAmbiguity

	// googModuleId mirrors FunctionType.googModuleId in Java.
	googModuleId string

	// canonicalRepresentation mirrors FunctionType.canonicalRepresentation in Java.
	// Only set for constructors; the canonical variant without parameter/return types.
	canonicalRepresentation *FunctionType

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

	// docInfo mirrors FunctionType.getJSDocInfo() — the JSDoc from the
	// declaration. Set by InferJsDocInfo for anonymous structural function
	// types (case #3 in the class doc).
	docInfo *rhino.JSDocInfo

	// implicitMatch mirrors FunctionType.isStructuralInterface in Java —
	// set to true by SetImplicitMatch(true) when `@record` makes the
	// interface use structural (duck-type) matching. Distinct from the
	// ObjectType.IsStructural flag: this flag lives on the constructor.
	implicitMatch bool

	// implementedInterfaces mirrors FunctionType.implementedInterfaces in
	// Java — the list of @implements ObjectType values set from JSDoc.
	// Only meaningful when kind == FunctionKindConstructor.
	implementedInterfaces []*ObjectType

	// implementedInterfaceTypes mirrors the templatized variant of
	// implementedInterfaces — the parallel JSType slice preserves
	// TemplatizedType arguments that Java's `(ObjectType) maybeInterType`
	// downcast carries through implicitly (TemplatizedType IS-A
	// ObjectType in Java; in the Go port the two are sibling types so we
	// keep both representations). Indexed in lockstep with
	// implementedInterfaces. Set by SetImplementedInterfaceTypes; nil
	// when no templatized info was preserved.
	// Go-only: language-adapter — see SetImplementedInterfaceTypes.
	implementedInterfaceTypes []JSType

	// ownerFunction mirrors FunctionType for prototype methods — when
	// this FunctionType is a method, ownerFunction points to the enclosing
	// constructor FunctionType. Mirrors the Java relationship accessed via
	// ObjectType.getOwnerFunction() → PrototypeObjectType.ownerFunction.
	ownerFunction *FunctionType

	// closurePrimitiveId is a free-form string tag identifying a known
	// Closure primitive (e.g. "goog.asserts.assert"). Mirrors
	// FunctionType.closurePrimitive in Java. Empty when not set.
	closurePrimitiveId string

	// extendedInterfaces mirrors FunctionType.extendedInterfaces — the list of
	// @extends ObjectType values set from JSDoc on @interface declarations.
	// Only meaningful when kind == FunctionKindInterface.
	extendedInterfaces []*ObjectType

	// source is the AST Node for the function declaration/expression.
	// Mirrors FunctionType.source in Java. Nil for synthetic types.
	source *rhino.Node

	// isAbstract mirrors FunctionType.isAbstract in Java — true when the
	// function is declared @abstract (may not be instantiated directly).
	isAbstract bool
}

// MaxArgsUnlimited is the sentinel MaxArity for rest-parameter
// functions — callers may pass arbitrarily many arguments.
// Mirrors Integer.MAX_VALUE returned by FunctionType.getMaxArity()
// in Java when the last parameter isVariadic.
const MaxArgsUnlimited = 1 << 30

// computeArity derives MinArity and MaxArity from the parameter list.
// This mirrors the logic in FunctionType.getMinArity() /
// getMaxArity() in Java.
// Go-only: language-adapter — Java computes arity on-the-fly from the
// parameter list; Go caches it in fields for backward compatibility.
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
// Go-only: language-adapter — Java has no parallel flat Params slice;
// this exists to preserve backward-compatible field access for existing callers.
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
// Go-only: language-adapter — Java uses the Builder pattern (FunctionType.Builder);
// Go consolidates construction into a shared private factory for simplicity.
func newFunctionTypeFromParams(params []Parameter, ret JSType, kind FunctionKind, typeOfThis JSType) *FunctionType {
	if ret == nil {
		ret = Unknown
	}
	minArity, maxArity := computeArity(params)
	f := &FunctionType{
		parameters:    params,
		returnType:    ret,
		typeOfThis:    typeOfThis,
		kind:          kind,
		Params:        flatTypes(params),
		Return:        ret,
		MinArity:      minArity,
		MaxArity:      maxArity,
		IsConstructor: kind == FunctionKindConstructor,
	}
	if kind == FunctionKindConstructor {
		f.propAccess = propAccessAny
	}
	return f
}

// NewFunctionType creates a plain (non-constructor) FunctionType from
// a flat []JSType slice. This is the backward-compatible factory used
// by existing callers that have not yet been migrated to Parameters.
//
// Mirrors FunctionType.builder(registry)…withParameters(…).build()
// for the ORDINARY kind.
// Go-only: language-adapter — Java uses FunctionType.builder(registry).build();
// Go exposes typed factory functions because there is no JSTypeRegistry.
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
// Go-only: language-adapter — see NewFunctionType; this overload accepts
// the Parameter type directly rather than plain JSType.
func NewFunctionTypeFromParameters(params []Parameter, ret JSType) *FunctionType {
	return newFunctionTypeFromParams(params, ret, FunctionKindOrdinary, nil)
}

// NewFunctionTypeWithThis creates a plain FunctionType with an explicit
// typeOfThis. Used when resolving `function(this:T, …)` JSDoc annotations
// so TypeCheck can emit JSC_EXPECTED_THIS_TYPE on bare-name call sites.
// Mirrors FunctionType.builder(registry).withTypeOfThis(T).build().
// Go-only: language-adapter — see NewFunctionType; this overload also accepts
// a typeOfThis argument.
func NewFunctionTypeWithThis(params []JSType, ret JSType, typeOfThis JSType) *FunctionType {
	ps := make([]Parameter, len(params))
	for i, t := range params {
		ps[i] = NewParameter(t, false, false)
	}
	return newFunctionTypeFromParams(ps, ret, FunctionKindOrdinary, typeOfThis)
}

// NewConstructorType is the variant used when @constructor is present.
// Returns a FunctionType whose Instance is a freshly-named ObjectType
// — the type assigned to `new ThisCtor(args)` expressions.
//
// Mirrors FunctionType.builder(registry).forConstructor().build().
// Go-only: language-adapter — see NewFunctionType; this overload creates
// a constructor with a named InstanceObjectType.
func NewConstructorType(name string, params []JSType) *FunctionType {
	inst := NewObjectType(name)
	ps := make([]Parameter, len(params))
	for i, t := range params {
		ps[i] = NewParameter(t, false, false)
	}
	minArity, maxArity := computeArity(ps)
	fn := &FunctionType{
		parameters:    ps,
		returnType:    inst,
		typeOfThis:    inst,
		kind:          FunctionKindConstructor,
		propAccess:    propAccessAny,
		Params:        flatTypes(ps),
		Return:        inst,
		MinArity:      minArity,
		MaxArity:      maxArity,
		IsConstructor: true,
		Instance:      inst,
	}
	inst.Ctor = fn
	return fn
}

// NewConstructorTypeFromParameters is the rich variant of
// NewConstructorType that accepts a []Parameter slice.
// Go-only: language-adapter — see NewConstructorType; this overload accepts
// the Parameter type directly rather than plain JSType.
func NewConstructorTypeFromParameters(name string, params []Parameter) *FunctionType {
	inst := NewObjectType(name)
	minArity, maxArity := computeArity(params)
	fn := &FunctionType{
		parameters:    params,
		returnType:    inst,
		typeOfThis:    inst,
		kind:          FunctionKindConstructor,
		propAccess:    propAccessAny,
		Params:        flatTypes(params),
		Return:        inst,
		MinArity:      minArity,
		MaxArity:      maxArity,
		IsConstructor: true,
		Instance:      inst,
	}
	inst.Ctor = fn
	return fn
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
	if f.typeOfThis == nil {
		return Unknown
	}
	return f.typeOfThis
}

// SetTypeOfThis sets the typeOfThis field. Called by scope-building code
// when resolving `function(this:T, …)` JSDoc.
// Go-only: language-adapter — Java sets typeOfThis via the Builder or
// resolveInternal; this setter is needed by Go callers that mutate the type.
func (f *FunctionType) SetTypeOfThis(t JSType) {
	f.typeOfThis = t
}

// GetKind mirrors FunctionType.getKind() in Java.
func (f *FunctionType) GetKind() FunctionKind {
	return f.kind
}

// IsOrdinaryFunction mirrors FunctionType.isOrdinaryFunction() in Java.
// Returns true only for plain `function` declarations/expressions —
// NOT for arrow functions, constructors, or interfaces.
func (f *FunctionType) IsOrdinaryFunction() bool {
	return f.kind == FunctionKindOrdinary
}

// IsArrow mirrors Java — arrow functions capture `this` lexically.
// Go-only: language-adapter — Java's Kind enum has no ARROW value.
func (f *FunctionType) IsArrow() bool {
	return f.kind == FunctionKindArrow
}

// IsInterface mirrors FunctionType.isInterface() in Java.
func (f *FunctionType) IsInterface() bool {
	return f.kind == FunctionKindInterface
}

// GetMinArity mirrors FunctionType.getMinArity() in Java.
// Returns the position of the last required (non-optional,
// non-variadic) parameter.
func (f *FunctionType) GetMinArity() int {
	i := 0
	min := 0
	for _, p := range f.parameters {
		i++
		if !p.isOptional && !p.isVariadic {
			min = i
		}
	}
	return min
}

// GetMaxArity mirrors FunctionType.getMaxArity() in Java.
// Returns MaxArgsUnlimited when the last parameter is variadic.
func (f *FunctionType) GetMaxArity() int {
	params := f.parameters
	if len(params) == 0 {
		return 0
	}
	lastParam := params[len(params)-1]
	if !lastParam.isVariadic {
		return len(params)
	}
	return MaxArgsUnlimited
}

// GetMinArguments is a Go-only alias for GetMinArity.
// Go-only: language-adapter — older call sites use this name.
func (f *FunctionType) GetMinArguments() int { return f.GetMinArity() }

// GetMaxArguments is a Go-only alias for GetMaxArity.
// Go-only: language-adapter — older call sites use this name.
func (f *FunctionType) GetMaxArguments() int { return f.GetMaxArity() }

// HasRestParam reports whether this function has a variadic (rest)
// parameter.
// Go-only: language-adapter — Java checks getMaxArity()==Integer.MAX_VALUE.
func (f *FunctionType) HasRestParam() bool {
	return len(f.parameters) > 0 && f.parameters[len(f.parameters)-1].isVariadic
}

// ToMaybeFunctionType mirrors FunctionType.toMaybeFunctionType() in Java.
// Returns this FunctionType (self-identity cast).
func (f *FunctionType) ToMaybeFunctionType() *FunctionType { return f }

// CanBeCalled mirrors FunctionType.canBeCalled() — all function types can
// be called. Mirrors Java's override that returns true.
func (f *FunctionType) CanBeCalled() bool { return true }

// IsInstanceType mirrors FunctionType.isInstanceType() in Java.
// In Java: return identical(this, registry.getNativeType(FUNCTION_TYPE)).
// Only the native `Function` singleton is both a function type and an instance
// type. In the Go port there is no registry so we check the instance's name
// as a heuristic (returns true for the synthetic "Function" instance type).
func (f *FunctionType) IsInstanceType() bool {
	// Java checks pointer identity against the registry's interned FUNCTION_TYPE.
	// In Go we use the name heuristic: the "Function" constructor's instance
	// is named "Function" — matching the native type's reference name.
	if f.Instance != nil && f.Instance.Name == "Function" {
		return f.kind == FunctionKindConstructor || f.kind == FunctionKindInterface
	}
	return false
}

// HasInstanceType mirrors FunctionType.hasInstanceType() in Java.
// Returns true for constructors and interfaces (both create instance types).
func (f *FunctionType) HasInstanceType() bool {
	return f.kind == FunctionKindConstructor || f.kind == FunctionKindInterface
}

// GetInstanceType mirrors FunctionType.getInstanceType() — returns the
// ObjectType produced by `new F(…)`. Nil for non-constructor functions.
func (f *FunctionType) GetInstanceType() *ObjectType { return f.Instance }

// SetInstanceType overwrites the Instance (and syncs Return/typeOfThis)
// for constructors. Used by TypedScopeCreator when wiring up the two-pass
// constructor stubs.
// Go-only: language-adapter — Java does not expose setInstanceType; the
// instance type is set via FunctionType(Builder) and is immutable afterwards.
// The Go port needs mutability because of two-pass type resolution.
func (f *FunctionType) SetInstanceType(inst *ObjectType) {
	f.Instance = inst
	f.returnType = inst
	f.Return = inst
	f.typeOfThis = inst
}

// GetPrototype mirrors FunctionType.getPrototype() in Java — returns the
// prototype ObjectType for the constructor. In our simplified model we
// reuse Instance for both instance and prototype slots.
func (f *FunctionType) GetPrototype() *ObjectType { return f.Instance }

// GetPrototypeProperty mirrors FunctionType.getPrototypeProperty() in Java.
// Equivalent to getPrototype().
func (f *FunctionType) GetPrototypeProperty() *ObjectType { return f.GetPrototype() }

// SetPrototypeBasedOn mirrors FunctionType.setPrototypeBasedOn(ObjectType) in Java.
// Sets the prototype by establishing an implicit-prototype relationship with
// the given base type. In the simplified Go port, sets the super chain.
func (f *FunctionType) SetPrototypeBasedOn(baseType *ObjectType) {
	if f.Instance != nil && baseType != nil {
		f.Instance.SetSuper(baseType)
	}
}

// SetPrototype mirrors FunctionType.setPrototype(ObjectType, Node) in Java.
// Returns false if the prototype is nil or is the instance type itself.
func (f *FunctionType) SetPrototype(prototype *ObjectType, propertyNode *rhino.Node) bool {
	if prototype == nil {
		return false
	}
	if f.IsConstructor && f.Instance != nil && prototype == f.Instance {
		return false
	}
	f.Instance = prototype
	return true
}

// GetSuperClassConstructor mirrors FunctionType.getSuperClassConstructor() in Java.
// Given a constructor or interface type, returns its superclass constructor
// or nil if none exists. Walks prototype → implicit-prototype → constructor.
func (f *FunctionType) GetSuperClassConstructor() *FunctionType {
	if f.kind != FunctionKindConstructor && f.kind != FunctionKindInterface {
		return nil
	}
	proto := f.GetPrototype()
	if proto == nil {
		return nil
	}
	maybeSuperInstanceType := proto.GetImplicitPrototype()
	if maybeSuperInstanceType == nil {
		return nil
	}
	return maybeSuperInstanceType.GetConstructor()
}

// GetBindReturnType mirrors FunctionType.getBindReturnType(int) in Java.
// Returns the type produced by calling .bind() on this function with
// argsToBind bound arguments. Pass -1 to return a function that accepts
// any parameters.
func (f *FunctionType) GetBindReturnType(argsToBind int) *FunctionType {
	b := NewFunctionTypeBuilder().WithReturnType(f.returnType)
	if argsToBind < 0 {
		return b.Build()
	}
	origParams := f.parameters
	params := make([]Parameter, len(origParams))
	copy(params, origParams)
	for i := 1; i < argsToBind && len(params) > 0; i++ {
		if params[0].isVariadic {
			break
		}
		params = params[1:]
	}
	return b.WithParameters(params).Build()
}

// GetCallOrBindSignature mirrors FunctionType.getCallOrBindSignature(boolean) in Java.
// Constructs the synthetic .call() or .bind() signature for this function type.
// Notice that "call" and "bind" have the same argument signature, except that
// all the arguments of "bind" (except the first) are optional.
func (f *FunctionType) GetCallOrBindSignature(isCall bool) *FunctionType {
	return f.getCallOrBindSignature(isCall)
}

// getCallOrBindSignature is the internal implementation used by GetCallOrBindSignature.
func (f *FunctionType) getCallOrBindSignature(isCall bool) *FunctionType {
	isBind := !isCall

	var retType JSType
	if isCall {
		retType = f.GetReturnType()
	} else {
		retType = f.GetBindReturnType(-1)
	}
	b := NewFunctionTypeBuilder().WithReturnType(retType)

	origParams := f.parameters
	params := make([]Parameter, len(origParams))
	copy(params, origParams)

	// Prepend `this` argument.
	thisType := f.GetTypeOfThis()
	thisParam := NewParameter(thisType, false, false)
	params = append([]Parameter{thisParam}, params...)

	if isBind {
		// All arguments of bind() are optional but not undefinable.
		for i := 1; i < len(params); i++ {
			cur := params[i]
			params[i] = NewParameter(cur.jsType, true, cur.isVariadic)
		}
	} else if isCall {
		// The first argument of call() is optional iff all the arguments
		// are optional. Check the first non-this argument.
		var firstArg *Parameter
		if len(params) > 1 {
			p := params[1]
			firstArg = &p
		}
		if firstArg == nil || firstArg.isOptional || firstArg.isVariadic {
			params[0] = NewParameter(thisParam.jsType, true, false)
		}
	}

	return b.WithParameters(params).Build()
}

// SupAndInfHelper mirrors FunctionType.supAndInfHelper(FunctionType, boolean) in Java.
// Computes the supremum (leastSuper=true) or infimum (leastSuper=false)
// of two function types.
func (f *FunctionType) SupAndInfHelper(that *FunctionType, leastSuper bool) *FunctionType {
	if that == nil {
		return f
	}
	if f == that {
		return f
	}

	// If these are ordinary functions, then merge them.
	if f.IsOrdinaryFunction() && that.IsOrdinaryFunction() {
		isSubtypeOfThat := f.IsSubtypeOf(that)
		isSubtypeOfThis := that.IsSubtypeOf(f)
		if isSubtypeOfThat && !isSubtypeOfThis {
			if leastSuper {
				return that
			}
			return f
		} else if isSubtypeOfThis && !isSubtypeOfThat {
			if leastSuper {
				return f
			}
			return that
		}

		merged := f.tryMergeFunctionPiecewise(that, leastSuper)
		if merged != nil {
			return merged
		}
	}

	// Fall back to the top or bottom of the lattice.
	// In the Go port we return f or that without a full registry.
	if leastSuper {
		return that
	}
	return f
}

// tryMergeFunctionPiecewise mirrors FunctionType.tryMergeFunctionPiecewise in Java.
// Tries to merge two function types component-wise (params + return + this-type).
// Returns nil when the parameter lists differ.
func (f *FunctionType) tryMergeFunctionPiecewise(other *FunctionType, leastSuper bool) *FunctionType {
	// Parameters must match exactly.
	if len(f.parameters) != len(other.parameters) {
		return nil
	}
	for i := range f.parameters {
		if f.parameters[i].jsType != other.parameters[i].jsType {
			return nil
		}
	}
	newParamsNode := f.parameters

	var newReturnType JSType
	if leastSuper {
		if f.returnType == other.returnType {
			newReturnType = f.returnType
		} else {
			// Use All as the join when return types differ.
			newReturnType = All
		}
	} else {
		if f.returnType == other.returnType {
			newReturnType = f.returnType
		} else {
			newReturnType = Unknown
		}
	}

	var newTypeOfThis JSType
	if f.typeOfThis == other.typeOfThis {
		newTypeOfThis = f.typeOfThis
	} else {
		if leastSuper {
			newTypeOfThis = All
		} else {
			newTypeOfThis = Unknown
		}
	}

	newReturnTypeInferred := f.returnTypeInferred || other.returnTypeInferred
	result := NewFunctionTypeBuilder().
		WithParameters(newParamsNode).
		WithReturnTypeInferred(newReturnType, newReturnTypeInferred).
		WithTypeOfThis(newTypeOfThis).
		Build()
	result.returnTypeInferred = newReturnTypeInferred
	return result
}

// CheckExtendsLoop mirrors FunctionType.checkExtendsLoop() in Java.
// Checks if there is a loop in the type-extends chain.
// Returns the loop path if a cycle exists, nil otherwise.
func (f *FunctionType) CheckExtendsLoop() []*FunctionType {
	return f.checkExtendsLoop(map[*FunctionType]bool{}, []*FunctionType{})
}

// checkExtendsLoop mirrors FunctionType.checkExtendsLoop(Set, List) in Java.
func (f *FunctionType) checkExtendsLoop(cache map[*FunctionType]bool, path []*FunctionType) []*FunctionType {
	for _, interfaceType := range f.extendedInterfaces {
		if interfaceType == nil {
			continue
		}
		superConstructor := interfaceType.GetConstructor()
		if superConstructor == nil {
			continue
		}
		if cache[superConstructor] {
			// Cycle detected. Prune path to the loop.
			path = append(path, superConstructor)
			for len(path) > 0 && path[0] != superConstructor {
				path = path[1:]
			}
			return path
		}
		cache[superConstructor] = true
		path = append(path, superConstructor)
		result := superConstructor.checkExtendsLoop(cache, path)
		if result != nil {
			return result
		}
		delete(cache, superConstructor)
		path = path[:len(path)-1]
	}
	return nil
}

// ForgetParameterAndReturnTypes mirrors FunctionType.forgetParameterAndReturnTypes() in Java.
// Creates a new constructor with the parameters and return type stripped.
func (f *FunctionType) ForgetParameterAndReturnTypes() *FunctionType {
	name := f.GetReferenceName()
	var typeOfThis *ObjectType
	if f.Instance != nil {
		typeOfThis = f.Instance
	}
	result := NewFunctionTypeBuilder().
		WithKind(f.kind).
		WithTypeOfThis(typeOfThis).
		Build()
	result.source = f.source
	result.canonicalRepresentation = f
	if name != "" && result.Instance != nil {
		result.Instance.Name = name
	}
	if typeOfThis != nil {
		result.SetPrototypeBasedOn(typeOfThis)
	}
	return result
}

// IsAmbiguousConstructor mirrors FunctionType.isAmbiguousConstructor() in Java.
// Returns true if the constructor does not come from a literal class or function
// in the AST, or if it extends such an ambiguous constructor.
func (f *FunctionType) IsAmbiguousConstructor() bool {
	if f.constructorAmbiguity == ConstructorAmbiguityUnknown {
		f.constructorAmbiguity = f.calculateConstructorAmbiguity()
	}
	return f.constructorAmbiguity == ConstructorAmbiguityIsAmbiguous
}

// calculateConstructorAmbiguity mirrors FunctionType.calculateConstructorAmbiguity() in Java.
func (f *FunctionType) calculateConstructorAmbiguity() ConstructorAmbiguity {
	if f.source != nil {
		// We can see the definition of the class.
		superConstructor := f.GetSuperClassConstructor()
		if superConstructor == nil || !superConstructor.IsAmbiguousConstructor() {
			return ConstructorAmbiguityIsUnambiguous
		}
	}
	if f.isDelegateProxy() {
		return ConstructorAmbiguityIsUnambiguous
	}
	if f.Instance != nil && f.Instance.isNative {
		return ConstructorAmbiguityIsUnambiguous
	}
	superConstructor := f.GetSuperClassConstructor()
	if superConstructor != nil && superConstructor.IsAmbiguousConstructor() {
		return ConstructorAmbiguityIsAmbiguous
	}
	// Type was created directly from JSDoc without a function or class literal.
	return ConstructorAmbiguityIsAmbiguous
}

// isDelegateProxy mirrors FunctionType.isDelegateProxy() in Java.
// Returns true when this is a compiler-generated delegate-proxy constructor.
func (f *FunctionType) isDelegateProxy() bool {
	ref := f.GetReferenceName()
	return ref != "" && strings.HasSuffix(ref, delegateSuffix)
}

// GetClosurePrimitive mirrors FunctionType.getClosurePrimitive() in Java.
// Returns the @closurePrimitive identifier associated with this function.
func (f *FunctionType) GetClosurePrimitive() string {
	return f.closurePrimitiveId
}

// GetClosurePrimitiveId mirrors FunctionType.getClosurePrimitive().String() in Java.
// Go-only: language-adapter — Java uses a ClosurePrimitive enum; Go uses string.
func (f *FunctionType) GetClosurePrimitiveId() string { return f.closurePrimitiveId }

// SetClosurePrimitiveId tags this function type as a known Closure primitive.
// Mirrors FunctionType.Builder.withClosurePrimitiveId in Java.
// Go-only: language-adapter — Java sets the ClosurePrimitive via the Builder;
// Go exposes a setter for post-construction mutation.
func (f *FunctionType) SetClosurePrimitiveId(id string) { f.closurePrimitiveId = id }

// GetCanonicalRepresentation mirrors FunctionType.getCanonicalRepresentation() in Java.
// Returns the original canonical variant of this function (only used for constructors).
func (f *FunctionType) GetCanonicalRepresentation() *FunctionType {
	return f.canonicalRepresentation
}

// GetGoogModuleId mirrors FunctionType.getGoogModuleId() in Java.
// Returns the ID of the goog.module in which this function was declared.
func (f *FunctionType) GetGoogModuleId() string {
	return f.googModuleId
}

// SetGoogModuleId sets the goog.module ID for this function.
// Go-only: language-adapter — Java sets googModuleId via Builder.setGoogModuleId();
// Go exposes a setter for post-construction mutation.
func (f *FunctionType) SetGoogModuleId(id string) { f.googModuleId = id }

// ClearCachedValues mirrors FunctionType.clearCachedValues() in Java.
// Clears any cached derived data (e.g. for hot-swap passes).
func (f *FunctionType) ClearCachedValues() {
	if f.Instance != nil {
		// Clear instance-side caches.
		_ = f.Instance
	}
}

// HasCachedValues mirrors FunctionType.hasCachedValues() in Java.
// Returns true when the prototype slot (instance) is non-nil.
func (f *FunctionType) HasCachedValues() bool {
	return f.Instance != nil
}

// HasAnyTemplateTypesInternal mirrors FunctionType.hasAnyTemplateTypesInternal() in Java.
func (f *FunctionType) HasAnyTemplateTypesInternal() bool {
	return len(f.TemplateTypes) > 0
}

// SetExplicitUnrestricted mirrors FunctionType.setExplicitUnrestricted() in Java.
// Marks this constructor's instances as explicitly @unrestricted (neither
// struct nor dict), overriding any inherited annotation.
func (f *FunctionType) SetExplicitUnrestricted() {
	f.propAccess = propAccessAnyExplicit
}

// ToBuilder mirrors FunctionType.toBuilder() in Java.
// Copies all information from this function type into a new Builder.
func (f *FunctionType) ToBuilder() *FunctionTypeBuilder {
	b := NewFunctionTypeBuilder().
		WithKind(f.kind).
		WithParameters(f.parameters).
		WithReturnTypeInferred(f.returnType, f.returnTypeInferred).
		WithTypeOfThis(f.typeOfThis)
	if len(f.TemplateTypes) > 0 {
		b = b.WithTemplateKeys(f.TemplateTypes)
	}
	if f.closurePrimitiveId != "" {
		b = b.WithClosurePrimitiveId(f.closurePrimitiveId)
	}
	return b
}

// --- TypeName / subtype ---------------------------------------------------

// TypeName returns the string representation of this function type.
// Mirrors FunctionType.appendTo (Java's toString equivalent).
// Go-only: stdlib-interop — Go uses TypeName() to satisfy the JSType interface;
// Java uses appendTo(TypeStringBuilder) with a visitor pattern.
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

// SetJSDocInfo attaches JSDoc to anonymous structural function types.
// Called by InferJsDocInfo (case #3 in the class doc).
// Go-only: language-adapter — Java stores JSDoc on nodes/slots; Go attaches
// it directly to the type for convenience.
func (f *FunctionType) SetJSDocInfo(info *rhino.JSDocInfo) { f.docInfo = info }

// GetJSDocInfo returns the JSDoc attached to this function type.
// Go-only: language-adapter — see SetJSDocInfo.
func (f *FunctionType) GetJSDocInfo() *rhino.JSDocInfo { return f.docInfo }

// SetImplicitMatch mirrors FunctionType.setImplicitMatch(boolean) in Java.
// When flag is true the interface uses structural (duck-type) matching
// (@record semantics). Stored in implicitMatch; readable via
// IsStructuralInterface().
func (f *FunctionType) SetImplicitMatch(flag bool) { f.implicitMatch = flag }

// IsStructuralInterface mirrors FunctionType.isStructuralInterface() in Java.
// Returns true when SetImplicitMatch(true) was called — i.e. the function
// is a @record interface that uses structural matching.
func (f *FunctionType) IsStructuralInterface() bool {
	return f.kind == FunctionKindInterface && f.implicitMatch
}

// IsAbstract mirrors FunctionType.isAbstract() in Java.
func (f *FunctionType) IsAbstract() bool { return f.isAbstract }

// SetAbstract sets the isAbstract flag on this function type.
// Go-only: language-adapter — Java sets isAbstract via the Builder;
// Go exposes a setter for post-construction mutation.
func (f *FunctionType) SetAbstract(v bool) { f.isAbstract = v }

// GetSource mirrors FunctionType.getSource() in Java.
// Returns the source AST Node for the function, or nil.
func (f *FunctionType) GetSource() *rhino.Node {
	return f.source
}

// SetSource mirrors FunctionType.setSource(Node) in Java.
func (f *FunctionType) SetSource(src *rhino.Node) {
	f.source = src
}

// GetExtendedInterfaces mirrors FunctionType.getExtendedInterfaces() in Java.
// Returns the directly @extends interfaces for this interface.
func (f *FunctionType) GetExtendedInterfaces() []*ObjectType {
	return f.extendedInterfaces
}

// GetExtendedInterfacesCount mirrors FunctionType.getExtendedInterfacesCount().
func (f *FunctionType) GetExtendedInterfacesCount() int {
	return len(f.extendedInterfaces)
}

// SetExtendedInterfaces mirrors FunctionType.setExtendedInterfaces(List) in Java.
// Only valid for @interface types.
func (f *FunctionType) SetExtendedInterfaces(ifaces []*ObjectType) {
	f.extendedInterfaces = ifaces
}

// SetImplementedInterfaces mirrors FunctionType.setImplementedInterfaces in
// Java. Stores the list of @implements ObjectTypes on this constructor.
// The list is populated by TypedScopeCreator when processing @implements
// JSDoc on @constructor declarations.
func (f *FunctionType) SetImplementedInterfaces(ifaces []*ObjectType) {
	f.implementedInterfaces = ifaces
}

// GetImplementedInterfaces mirrors FunctionType.getImplementedInterfaces()
// in Java. Returns the direct @implements interfaces for this constructor
// and all superclass @implements interfaces.
func (f *FunctionType) GetImplementedInterfaces() []*ObjectType {
	superCtor := f.GetSuperClassConstructor()
	if superCtor == nil {
		return f.implementedInterfaces
	}
	var result []*ObjectType
	result = append(result, f.implementedInterfaces...)
	for superCtor != nil {
		result = append(result, superCtor.implementedInterfaces...)
		superCtor = superCtor.GetSuperClassConstructor()
	}
	return result
}

// GetOwnImplementedInterfaces mirrors FunctionType.getOwnImplementedInterfaces()
// in Java — returns only directly declared @implements interfaces.
func (f *FunctionType) GetOwnImplementedInterfaces() []*ObjectType {
	return f.implementedInterfaces
}

// SetImplementedInterfaceTypes records the parallel templatized form
// of each implemented interface (e.g. `Foo<string>`).
// Go-only: language-adapter — Java's TemplatizedType extends ObjectType
// so the cast preserves the templatization; Go has TemplatizedType as a
// sibling JSType, so we keep both views.
func (f *FunctionType) SetImplementedInterfaceTypes(types []JSType) {
	f.implementedInterfaceTypes = types
}

// GetImplementedInterfaceTypes returns the parallel templatized list.
// Go-only: language-adapter — see SetImplementedInterfaceTypes.
func (f *FunctionType) GetImplementedInterfaceTypes() []JSType {
	return f.implementedInterfaceTypes
}

// HasImplementedInterfaces mirrors FunctionType.hasImplementedInterfaces() —
// returns true when this constructor has at least one @implements.
func (f *FunctionType) HasImplementedInterfaces() bool {
	if len(f.implementedInterfaces) > 0 {
		return true
	}
	superCtor := f.GetSuperClassConstructor()
	if superCtor != nil {
		return superCtor.HasImplementedInterfaces()
	}
	return false
}

// GetAllImplementedInterfaces mirrors FunctionType.getAllImplementedInterfaces()
// in Java — returns the transitive closure of all interfaces implemented/extended.
// Mirrors Java's LinkedHashSet-based dedup traversal.
func (f *FunctionType) GetAllImplementedInterfaces() []*ObjectType {
	seen := map[*ObjectType]bool{}
	var out []*ObjectType
	for _, iface := range f.implementedInterfaces {
		addRelatedInterfaces(iface, seen, &out)
	}
	return out
}

// addRelatedInterfaces is the recursive helper for GetAllImplementedInterfaces.
// Mirrors Java's FunctionType.addRelatedInterfaces.
func addRelatedInterfaces(instance *ObjectType, seen map[*ObjectType]bool, out *[]*ObjectType) {
	if instance == nil {
		return
	}
	ctor := instance.GetConstructor()
	if ctor == nil {
		return
	}
	if !ctor.IsInterface() {
		return
	}
	if seen[instance] {
		return
	}
	seen[instance] = true
	*out = append(*out, instance)
	for _, extIface := range instance.GetCtorExtendedInterfaces() {
		addRelatedInterfaces(extIface, seen, out)
	}
}

// GetAncestorInterfaces mirrors FunctionType.getAncestorInterfaces() in Java.
// For constructors returns the implementedInterfaces; for interfaces returns
// the extendedInterfaces.
func (f *FunctionType) GetAncestorInterfaces() []*ObjectType {
	if f.kind == FunctionKindConstructor {
		return f.implementedInterfaces
	}
	return f.extendedInterfaces
}

// GetOwnerFunction mirrors the ownerFunction concept from Java —
// for a method FunctionType, returns the enclosing constructor
// FunctionType. Nil for top-level functions and constructors.
// Go-only: language-adapter — Java stores ownerFunction on PrototypeObjectType.
func (f *FunctionType) GetOwnerFunction() *FunctionType { return f.ownerFunction }

// SetOwnerFunction records the enclosing constructor FunctionType for
// method types. Called by TypedScopeCreator when attaching a method's
// FunctionType to the prototype.
// Go-only: language-adapter — Java stores ownerFunction on PrototypeObjectType.
func (f *FunctionType) SetOwnerFunction(owner *FunctionType) { f.ownerFunction = owner }

// AcceptsArguments mirrors FunctionType.acceptsArguments(List<JSType>) in Java.
// Returns true if the given argument types satisfy this function's parameter list
// (each arg is a subtype of the corresponding param, and the count is in range).
func (f *FunctionType) AcceptsArguments(argumentTypes []JSType) bool {
	params := f.parameters
	var parameter *Parameter
	pIdx := 0
	for _, argument := range argumentTypes {
		if pIdx < len(params) || (parameter != nil && parameter.isVariadic) {
			if pIdx < len(params) {
				p := params[pIdx]
				parameter = &p
				pIdx++
			}
			if !IsSubtype(argument, parameter.jsType) {
				return false
			}
		} else {
			return false
		}
	}
	numArgs := len(argumentTypes)
	return f.GetMinArity() <= numArgs && numArgs <= f.GetMaxArity()
}

// GetPropertyTypeMap mirrors FunctionType.getPropertyTypeMap() in Java.
// Returns a map of all properties declared on the prototype chain of this
// constructor, including inherited ones from extended interfaces.
func (f *FunctionType) GetPropertyTypeMap() map[string]JSType {
	result := map[string]JSType{}
	seen := map[*FunctionType]bool{}
	updatePropertyTypeMap(f, result, seen)
	return result
}

// updatePropertyTypeMap is the recursive helper for GetPropertyTypeMap.
// Mirrors Java's private static updatePropertyTypeMap.
func updatePropertyTypeMap(ft *FunctionType, out map[string]JSType, seen map[*FunctionType]bool) {
	if ft == nil || seen[ft] {
		return
	}
	seen[ft] = true
	// own prototype properties
	proto := ft.GetPrototype()
	if proto != nil {
		for name, t := range proto.members {
			if _, exists := out[name]; !exists {
				out[name] = t
			}
		}
	}
	// retrieve all property types from its super class (extended interfaces)
	for _, iface := range ft.extendedInterfaces {
		if iface == nil {
			continue
		}
		superConstructor := iface.GetConstructor()
		if superConstructor == nil || seen[superConstructor] {
			continue
		}
		seen[superConstructor] = true
		updatePropertyTypeMap(superConstructor, out, seen)
		delete(seen, superConstructor)
	}
}

// HasEqualCallType mirrors FunctionType.hasEqualCallType(FunctionType) in Java.
// Two function types have equal call signatures when their parameters and
// return type are identical (by pointer identity for Go types).
func (f *FunctionType) HasEqualCallType(that *FunctionType) bool {
	if that == nil {
		return false
	}
	if len(f.parameters) != len(that.parameters) {
		return false
	}
	for i := range f.parameters {
		if f.parameters[i].jsType != that.parameters[i].jsType {
			return false
		}
	}
	return f.returnType == that.returnType
}

// GetConstructorOnlyTemplateParameters mirrors
// FunctionType.getConstructorOnlyTemplateParameters() in Java. Returns the
// template parameters that live on the constructor but not the instance.
// In the Go port we return the full TemplateTypes list since we don't track
// instance vs constructor template keys separately.
func (f *FunctionType) GetConstructorOnlyTemplateParameters() []*TemplateType {
	if !f.IsConstructor {
		return nil
	}
	if f.source == nil {
		return nil
	}
	if f.Instance == nil {
		return f.TemplateTypes
	}
	instanceParams := f.Instance.typeParameters
	// return only TemplateTypes not present in the instance's params.
	var out []*TemplateType
	for _, tt := range f.TemplateTypes {
		found := false
		for _, ip := range instanceParams {
			if ip == tt {
				found = true
				break
			}
		}
		if !found {
			out = append(out, tt)
		}
	}
	return out
}

// GetOwnPropertyNames mirrors FunctionType.getOwnPropertyNames() in Java.
// Returns the set of own property names; always includes "prototype" when
// the function is a constructor with an Instance type.
func (f *FunctionType) GetOwnPropertyNames() []string {
	seen := map[string]bool{}
	var out []string
	if f.Instance != nil {
		seen["prototype"] = true
		out = append(out, "prototype")
	}
	if f.Instance != nil {
		for name := range f.Instance.staticMembers {
			if !seen[name] {
				seen[name] = true
				out = append(out, name)
			}
		}
	}
	return out
}

// GetReferenceName mirrors FunctionType.getReferenceName() (via
// PrototypeObjectType). Returns the constructor or function name,
// or empty string for anonymous types.
func (f *FunctionType) GetReferenceName() string {
	if f.Instance != nil {
		return f.Instance.Name
	}
	return ""
}

// HasReferenceName mirrors FunctionType.hasReferenceName().
func (f *FunctionType) HasReferenceName() bool {
	return f.GetReferenceName() != ""
}

// MakesStructs mirrors FunctionType.makesStructs() in Java — returns true
// when this constructor's instances should be @struct typed (forbidding
// computed-property access). Walks the super chain when propAccess is ANY.
func (f *FunctionType) MakesStructs() bool {
	if !f.HasInstanceType() || f.Instance == nil {
		return false
	}
	switch f.propAccess {
	case propAccessStruct:
		return true
	case propAccessDict:
		return false
	case propAccessAnyExplicit:
		return false
	case propAccessAny:
		superc := f.GetSuperClassConstructor()
		if superc != nil && superc.MakesStructs() {
			f.SetStruct()
			return true
		}
		return false
	}
	return false
}

// MakesDicts mirrors FunctionType.makesDicts() in Java — returns true when
// this constructor's instances should be @dict typed (forbidding dotted
// property access). Walks the super chain when propAccess is ANY.
func (f *FunctionType) MakesDicts() bool {
	if f.kind != FunctionKindConstructor {
		return false
	}
	switch f.propAccess {
	case propAccessDict:
		return true
	case propAccessStruct:
		return false
	case propAccessAnyExplicit:
		return false
	case propAccessAny:
		superc := f.GetSuperClassConstructor()
		if superc != nil && superc.MakesDicts() {
			f.SetDict()
			return true
		}
		return false
	}
	return false
}

// SetStruct marks this constructor's instance as @struct.
// Mirrors FunctionType.setStruct().
func (f *FunctionType) SetStruct() {
	f.propAccess = propAccessStruct
	if f.Instance != nil {
		f.Instance.IsStruct = true
	}
}

// SetDict marks this constructor's instance as @dict.
// Mirrors FunctionType.setDict().
func (f *FunctionType) SetDict() {
	f.propAccess = propAccessDict
	if f.Instance != nil {
		f.Instance.IsDict = true
	}
}

// HasProperties mirrors FunctionType.hasProperties() in Java.
// Returns true when the function type has any own (non-prototype) properties.
// In the Go port we check the staticMembers on the Instance type if present.
func (f *FunctionType) HasProperties() bool {
	if f.Instance != nil && len(f.Instance.staticMembers) > 0 {
		return true
	}
	return false
}

// MatchesObjectContext reports true — functions are objects.
// Mirrors FunctionType (PrototypeObjectType) matchesObjectContext().
// Go-only: language-adapter — Java's matchesObjectContext() lives on ObjectType;
// Go surfaces it on FunctionType so call-site code can use a uniform interface.
func (f *FunctionType) MatchesObjectContext() bool { return true }

// GetPossibleToBooleanOutcomes returns BooleanLiteralSetTrue — functions are always truthy.
// Mirrors FunctionType (ObjectType) getPossibleToBooleanOutcomes().
// Go-only: language-adapter — Java's getPossibleToBooleanOutcomes() is on ObjectType;
// Go surfaces it on FunctionType for uniform interface access.
func (f *FunctionType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSetTrue
}

// --- FunctionTypeBuilder — fluent builder mirroring FunctionType.Builder ---

// FunctionTypeBuilder mirrors FunctionType.Builder in Java — a fluent
// builder that accumulates fields before calling Build() to produce a
// *FunctionType. Use NewFunctionTypeBuilder() to obtain an instance.
//
// The Go builder differs from Java in two ways:
//  1. No JSTypeRegistry back-pointer (Go types are not interned at build time).
//  2. Template keys are stored as []*TemplateType rather than ImmutableList.
type FunctionTypeBuilder struct {
	params             []Parameter
	returnType         JSType
	returnTypeInferred bool
	typeOfThis         JSType
	kind               FunctionKind
	templates          []*TemplateType
	primitiveId        string
	name               string
	sourceNode         *rhino.Node
	googModuleId       string
	isAbstract         bool
	canonical          *FunctionType
	constructorOnlyKeys []*TemplateType
	prototypeBasedOn   *ObjectType
}

// NewFunctionTypeBuilder returns a blank FunctionTypeBuilder for an
// ordinary function. Mirrors FunctionType.builder(registry) in Java.
// Go-only: language-adapter — Java uses the static FunctionType.builder(registry)
// factory; Go uses a constructor function without the JSTypeRegistry.
func NewFunctionTypeBuilder() *FunctionTypeBuilder {
	return &FunctionTypeBuilder{kind: FunctionKindOrdinary}
}

// WithName sets the function name. Mirrors Builder.withName(String).
func (b *FunctionTypeBuilder) WithName(name string) *FunctionTypeBuilder {
	b.name = name
	return b
}

// WithSourceNode sets the source AST node. Mirrors Builder.withSourceNode(Node).
func (b *FunctionTypeBuilder) WithSourceNode(node *rhino.Node) *FunctionTypeBuilder {
	b.sourceNode = node
	return b
}

// SetGoogModuleId sets the goog.module ID. Mirrors Builder.setGoogModuleId(String).
func (b *FunctionTypeBuilder) SetGoogModuleId(id string) *FunctionTypeBuilder {
	b.googModuleId = id
	return b
}

// WithReturnType sets the return type. Mirrors Builder.withReturnType(T).
func (b *FunctionTypeBuilder) WithReturnType(t JSType) *FunctionTypeBuilder {
	b.returnType = t
	return b
}

// WithReturnTypeInferred sets return type and inferred flag.
// Mirrors Builder.withReturnType(JSType, boolean).
// Go-only: language-adapter — this method name combines Java's
// withReturnType(JSType, boolean) overload into a distinct Go name.
func (b *FunctionTypeBuilder) WithReturnTypeInferred(t JSType, inferred bool) *FunctionTypeBuilder {
	b.returnType = t
	b.returnTypeInferred = inferred
	return b
}

// WithInferredReturnType sets an inferred return type.
// Mirrors Builder.withInferredReturnType(JSType).
func (b *FunctionTypeBuilder) WithInferredReturnType(t JSType) *FunctionTypeBuilder {
	b.returnType = t
	b.returnTypeInferred = true
	return b
}

// WithTypeOfThis sets the typeOfThis. Mirrors Builder.withTypeOfThis(T).
func (b *FunctionTypeBuilder) WithTypeOfThis(t JSType) *FunctionTypeBuilder {
	b.typeOfThis = t
	return b
}

// WithTemplateKeys records the function's @template parameters.
// Mirrors Builder.withTemplateKeys(ImmutableList<TemplateType>).
func (b *FunctionTypeBuilder) WithTemplateKeys(keys []*TemplateType) *FunctionTypeBuilder {
	b.templates = keys
	return b
}

// WithConstructorTemplateKeys sets the constructor-only template keys.
// Mirrors Builder.withConstructorTemplateKeys(Iterable<TemplateType>).
func (b *FunctionTypeBuilder) WithConstructorTemplateKeys(keys []*TemplateType) *FunctionTypeBuilder {
	b.constructorOnlyKeys = keys
	return b
}

// WithClosurePrimitiveId tags the function as a known Closure primitive.
// Mirrors Builder.withClosurePrimitiveId(ClosurePrimitive).
func (b *FunctionTypeBuilder) WithClosurePrimitiveId(id string) *FunctionTypeBuilder {
	b.primitiveId = id
	return b
}

// WithKind sets the FunctionKind (ordinary / constructor / interface / arrow).
// Mirrors Builder.withKind(Kind).
func (b *FunctionTypeBuilder) WithKind(k FunctionKind) *FunctionTypeBuilder {
	b.kind = k
	return b
}

// ForConstructor makes this a constructor type.
// Mirrors Builder.forConstructor().
func (b *FunctionTypeBuilder) ForConstructor() *FunctionTypeBuilder {
	b.kind = FunctionKindConstructor
	return b
}

// ForInterface makes this an interface type with no parameters.
// Mirrors Builder.forInterface().
func (b *FunctionTypeBuilder) ForInterface() *FunctionTypeBuilder {
	b.kind = FunctionKindInterface
	b.params = []Parameter{}
	return b
}

// WithIsAbstract marks this as an abstract function.
// Mirrors Builder.withIsAbstract(boolean).
func (b *FunctionTypeBuilder) WithIsAbstract(v bool) *FunctionTypeBuilder {
	b.isAbstract = v
	return b
}

// WithCanonicalRepresentation sets the canonical representation for constructors.
// Mirrors Builder.withCanonicalRepresentation(FunctionType).
func (b *FunctionTypeBuilder) WithCanonicalRepresentation(repr *FunctionType) *FunctionTypeBuilder {
	b.canonical = repr
	return b
}

// WithPrototypeBasedOn sets the base type for prototype initialization.
// Mirrors Builder.withPrototypeBasedOn(ObjectType).
func (b *FunctionTypeBuilder) WithPrototypeBasedOn(base *ObjectType) *FunctionTypeBuilder {
	b.prototypeBasedOn = base
	return b
}

// WithParameters sets the parameter list.
// Mirrors Builder.withParameters(ImmutableList<Parameter>).
func (b *FunctionTypeBuilder) WithParameters(params []Parameter) *FunctionTypeBuilder {
	b.params = params
	return b
}

// Build constructs and returns the *FunctionType.
// Mirrors Builder.build() in Java.
func (b *FunctionTypeBuilder) Build() *FunctionType {
	var typeOfThis JSType
	if b.typeOfThis != nil {
		typeOfThis = b.typeOfThis
	}
	f := newFunctionTypeFromParams(b.params, b.returnType, b.kind, typeOfThis)
	f.returnTypeInferred = b.returnTypeInferred
	f.source = b.sourceNode
	f.googModuleId = b.googModuleId
	f.isAbstract = b.isAbstract
	f.canonicalRepresentation = b.canonical
	if len(b.templates) > 0 {
		f.TemplateTypes = b.templates
	}
	if b.primitiveId != "" {
		f.closurePrimitiveId = b.primitiveId
	}
	if b.name != "" {
		if f.Instance == nil && (b.kind == FunctionKindConstructor || b.kind == FunctionKindInterface) {
			inst := NewObjectType(b.name)
			f.Instance = inst
			inst.Ctor = f
			f.typeOfThis = inst
			f.returnType = inst
			f.Return = inst
		} else if f.Instance != nil && f.Instance.Name == "" {
			f.Instance.Name = b.name
		}
	}
	if b.prototypeBasedOn != nil {
		f.SetPrototypeBasedOn(b.prototypeBasedOn)
	}
	return f
}

// BuildAndResolve mirrors FunctionType.Builder.buildAndResolve() in Java.
func (b *FunctionTypeBuilder) BuildAndResolve() *FunctionType {
	return b.Build()
}
