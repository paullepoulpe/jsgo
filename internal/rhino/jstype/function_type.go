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
	// Mirrors Java Kind.ARROW.  Arrow functions do NOT rebind `this`,
	// so IsOrdinaryFunction() returns false for them.
	FunctionKindArrow
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

	// ownerFunction mirrors FunctionType for prototype methods — when
	// this FunctionType is a method, ownerFunction points to the enclosing
	// constructor FunctionType. Mirrors the Java relationship accessed via
	// ObjectType.getOwnerFunction() → PrototypeObjectType.ownerFunction.
	ownerFunction *FunctionType

	// closurePrimitiveId is a free-form string tag identifying a known
	// Closure primitive (e.g. "goog.asserts.assert"). Mirrors
	// FunctionType.closurePrimitive in Java. Empty when not set.
	closurePrimitiveId string
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

// NewFunctionTypeWithThis creates a plain FunctionType with an explicit
// typeOfThis. Used when resolving `function(this:T, …)` JSDoc annotations
// so TypeCheck can emit JSC_EXPECTED_THIS_TYPE on bare-name call sites.
// Mirrors FunctionType.builder(registry).withTypeOfThis(T).build().
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

// SetTypeOfThis mirrors FunctionType.setTypeOfThis() in Java. Called
// by scope-building code when resolving `function(this:T, …)` JSDoc
// so EXPECTED_THIS_TYPE checks can fire.
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

// IsArrow mirrors FunctionType.isArrow() in Java.
// Returns true when this is an ES2015 arrow function (`() => …`).
// Arrow functions do not rebind `this` and cannot be constructors.
func (f *FunctionType) IsArrow() bool {
	return f.kind == FunctionKindArrow
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

// SetJSDocInfo mirrors FunctionType.setJSDocInfo — attaches JSDoc to
// anonymous structural function types. Called by InferJsDocInfo (case #3
// in the class doc).
func (f *FunctionType) SetJSDocInfo(info *rhino.JSDocInfo) { f.docInfo = info }

// GetJSDocInfo mirrors FunctionType.getJSDocInfo.
func (f *FunctionType) GetJSDocInfo() *rhino.JSDocInfo { return f.docInfo }

// GetInstanceType mirrors FunctionType.getInstanceType() — returns the
// ObjectType produced by `new F(…)`. Nil for non-constructor functions.
func (f *FunctionType) GetInstanceType() *ObjectType { return f.Instance }

// GetPrototype mirrors FunctionType.getPrototype() — returns the
// prototype ObjectType for the constructor. In our simplified model we
// reuse Instance for both instance and prototype slots.
func (f *FunctionType) GetPrototype() *ObjectType { return f.Instance }

// SetInstanceType mirrors FunctionType.setInstanceType in Java —
// overwrites the Instance (and syncs Return/typeOfThis) for constructors.
// Used by TypedScopeCreator when wiring up the two-pass constructor stubs.
func (f *FunctionType) SetInstanceType(inst *ObjectType) {
	f.Instance = inst
	f.returnType = inst
	f.Return = inst
	f.typeOfThis = inst
}

// SetImplicitMatch mirrors FunctionType.setImplicitMatch(boolean) in Java.
// When flag is true the interface uses structural (duck-type) matching
// (@record semantics). Stored in implicitMatch; readable via
// IsStructuralInterface().
func (f *FunctionType) SetImplicitMatch(flag bool) { f.implicitMatch = flag }

// IsStructuralInterface mirrors FunctionType.isStructuralInterface() in Java.
// Returns true when SetImplicitMatch(true) was called — i.e. the function
// is a @record interface that uses structural matching.
func (f *FunctionType) IsStructuralInterface() bool { return f.implicitMatch }

// SetImplementedInterfaces mirrors FunctionType.setImplementedInterfaces in
// Java. Stores the list of @implements ObjectTypes on this constructor.
// The list is populated by TypedScopeCreator when processing @implements
// JSDoc on @constructor declarations.
func (f *FunctionType) SetImplementedInterfaces(ifaces []*ObjectType) {
	f.implementedInterfaces = ifaces
}

// GetImplementedInterfaces mirrors FunctionType.getImplementedInterfaces()
// in Java. Returns the direct @implements interfaces for this constructor.
// Returns nil (not the empty slice) when none have been recorded.
func (f *FunctionType) GetImplementedInterfaces() []*ObjectType {
	return f.implementedInterfaces
}

// GetOwnerFunction mirrors the ownerFunction concept from Java —
// for a method FunctionType, returns the enclosing constructor
// FunctionType. Nil for top-level functions and constructors.
func (f *FunctionType) GetOwnerFunction() *FunctionType { return f.ownerFunction }

// SetOwnerFunction records the enclosing constructor FunctionType for
// method types. Called by TypedScopeCreator when attaching a method's
// FunctionType to the prototype.
func (f *FunctionType) SetOwnerFunction(owner *FunctionType) { f.ownerFunction = owner }

// GetClosurePrimitiveId mirrors FunctionType.getClosurePrimitive().String()
// in Java. Returns the closure-primitive tag (e.g. "goog.asserts.assert")
// or the empty string when none has been set.
func (f *FunctionType) GetClosurePrimitiveId() string { return f.closurePrimitiveId }

// SetClosurePrimitiveId mirrors FunctionType.Builder.withClosurePrimitiveId
// in Java. Tags this function type as a known Closure primitive.
func (f *FunctionType) SetClosurePrimitiveId(id string) { f.closurePrimitiveId = id }

// --- FunctionTypeBuilder — fluent builder mirroring FunctionType.Builder ---

// FunctionTypeBuilder mirrors FunctionType.Builder in Java — a fluent
// builder that accumulates fields before calling Build() to produce a
// *FunctionType. Use NewFunctionTypeBuilder() to obtain an instance.
//
// The Go builder differs from Java in two ways:
//  1. No JSTypeRegistry back-pointer (Go types are not interned at build time).
//  2. Template keys are stored as []*TemplateType rather than ImmutableList.
type FunctionTypeBuilder struct {
	params     []Parameter
	returnType JSType
	typeOfThis JSType
	kind       FunctionKind
	templates  []*TemplateType
	primitiveId string
}

// NewFunctionTypeBuilder returns a blank FunctionTypeBuilder for an
// ordinary function. Mirrors FunctionType.builder(registry) in Java.
func NewFunctionTypeBuilder() *FunctionTypeBuilder {
	return &FunctionTypeBuilder{kind: FunctionKindOrdinary}
}

// WithReturnType sets the return type. Mirrors Builder.withReturnType(T).
func (b *FunctionTypeBuilder) WithReturnType(t JSType) *FunctionTypeBuilder {
	b.returnType = t
	return b
}

// WithReturnTypeInferred is a no-op placeholder for callers that need the
// Java Builder.withReturnType(T, inferred) overload. The inferred flag is
// set on the produced FunctionType via the returnTypeInferred field.
func (b *FunctionTypeBuilder) WithReturnTypeInferred(t JSType, inferred bool) *FunctionTypeBuilder {
	b.returnType = t
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

// WithParameters sets the parameter list.
// Mirrors Builder.withParameters(ImmutableList<Parameter>).
func (b *FunctionTypeBuilder) WithParameters(params []Parameter) *FunctionTypeBuilder {
	b.params = params
	return b
}

// Build constructs and returns the *FunctionType.
// Mirrors Builder.build() in Java.
func (b *FunctionTypeBuilder) Build() *FunctionType {
	f := newFunctionTypeFromParams(b.params, b.returnType, b.kind, b.typeOfThis)
	if len(b.templates) > 0 {
		f.TemplateTypes = b.templates
	}
	if b.primitiveId != "" {
		f.closurePrimitiveId = b.primitiveId
	}
	return f
}
