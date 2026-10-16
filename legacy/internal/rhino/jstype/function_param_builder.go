// Port of com/google/javascript/rhino/jstype/FunctionParamBuilder.java @ closure-compiler v20260415
//
// FunctionParamBuilder is a convenience builder for the ordered
// Parameter list stored on FunctionType. It enforces the rule that
// optional parameters must precede variadic ones, and that variadic
// parameters may only appear once (at the end).

package jstype

// FunctionParamBuilder accumulates typed parameters for a FunctionType
// and produces an immutable []Parameter slice via Build().
//
// Port of com/google/javascript/rhino/jstype/FunctionParamBuilder.
// The Java class carries a JSTypeRegistry reference so that optional
// params can be wrapped with createOptionalType; in Go we do not yet
// have a registry singleton, so callers must pre-wrap the type with
// the optional-union themselves when needed (the builder still sets
// the isOptional bit on the Parameter).
//
// TODO(go-port): Wire in the type registry when it is available so
// that AddOptionalParams can automatically wrap types with |undefined
// the same way the Java implementation does via
// registry.createOptionalType(type).
type FunctionParamBuilder struct {
	parameters []Parameter
}

// NewFunctionParamBuilder returns an empty FunctionParamBuilder.
// Mirrors FunctionParamBuilder(JSTypeRegistry) in Java.
func NewFunctionParamBuilder() *FunctionParamBuilder {
	return &FunctionParamBuilder{}
}

// NewFunctionParamBuilderWithCapacity returns an empty builder with
// the given initial capacity pre-allocated.
// Mirrors FunctionParamBuilder(JSTypeRegistry, int) in Java.
func NewFunctionParamBuilderWithCapacity(cap int) *FunctionParamBuilder {
	return &FunctionParamBuilder{parameters: make([]Parameter, 0, cap)}
}

// AddRequiredParams appends required (non-optional, non-variadic)
// parameters to the list. Returns false and is a no-op if optional or
// variadic params have already been added (callers must add required
// params before optional/variadic ones).
// Mirrors FunctionParamBuilder.addRequiredParams(JSType...) in Java.
func (b *FunctionParamBuilder) AddRequiredParams(types ...JSType) bool {
	if b.hasOptionalOrVarArgs() {
		return false
	}
	for _, t := range types {
		b.newParameter(t, false, false)
	}
	return true
}

// AddOptionalParams appends optional (isOptional=true) parameters.
// The caller is responsible for wrapping the type with |undefined if
// that is the desired JSDoc behaviour (Java's builder wraps via
// registry.createOptionalType). Returns false if variadic params have
// already been added.
// Mirrors FunctionParamBuilder.addOptionalParams(JSType...) in Java.
func (b *FunctionParamBuilder) AddOptionalParams(types ...JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	for _, t := range types {
		b.newParameter(t, true, false)
	}
	return true
}

// AddVarArgs appends a single variadic (rest) parameter. Returns false
// if a variadic parameter has already been added (only one is allowed,
// and it must be last).
// Mirrors FunctionParamBuilder.addVarArgs(JSType) in Java.
func (b *FunctionParamBuilder) AddVarArgs(t JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	b.newParameter(t, false, true)
	return true
}

// NewParameterFrom copies an existing Parameter into the builder.
// Mirrors FunctionParamBuilder.newParameterFrom(Parameter) in Java.
func (b *FunctionParamBuilder) NewParameterFrom(p Parameter) {
	b.parameters = append(b.parameters, p)
}

// NewOptionalParameterFrom copies p but ensures the isOptional bit is
// set (unless it is variadic, in which case the isOptional is left as
// derived from the variadic-or-not semantics).
// Mirrors FunctionParamBuilder.newOptionalParameterFrom(Parameter) in Java.
func (b *FunctionParamBuilder) NewOptionalParameterFrom(p Parameter) {
	isOptional := p.isOptional || !p.isVariadic
	if isOptional != p.isOptional {
		b.newParameter(p.jsType, isOptional, p.isVariadic)
	} else {
		b.NewParameterFrom(p)
	}
}

// newParameter is the internal helper that creates and appends a
// Parameter. Mirrors the private newParameter method in Java.
func (b *FunctionParamBuilder) newParameter(t JSType, isOptional, isVariadic bool) {
	b.parameters = append(b.parameters, NewParameter(t, isOptional, isVariadic))
}

// hasOptionalOrVarArgs returns true when the last parameter is
// optional or variadic. Mirrors the private hasOptionalOrVarArgs()
// in Java.
func (b *FunctionParamBuilder) hasOptionalOrVarArgs() bool {
	if len(b.parameters) == 0 {
		return false
	}
	last := b.parameters[len(b.parameters)-1]
	return last.isOptional || last.isVariadic
}

// HasVarArgs reports whether a variadic parameter has already been
// added. Mirrors FunctionParamBuilder.hasVarArgs() in Java.
func (b *FunctionParamBuilder) HasVarArgs() bool {
	return len(b.parameters) > 0 && b.parameters[len(b.parameters)-1].isVariadic
}

// Build returns the accumulated Parameter slice.
// The returned slice must not be mutated by the caller.
// Mirrors FunctionParamBuilder.build() in Java (which returns an
// ImmutableList).
func (b *FunctionParamBuilder) Build() []Parameter {
	out := make([]Parameter, len(b.parameters))
	copy(out, b.parameters)
	return out
}
