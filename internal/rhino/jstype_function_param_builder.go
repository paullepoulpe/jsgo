// Port of com/google/javascript/rhino/jstype/FunctionParamBuilder.java @ closure-compiler v20260415

package rhino

// FunctionParamBuilder is a builder for the list representing FunctionType
// Parameters.
type FunctionParamBuilder struct {
	registry   *JSTypeRegistry
	parameters []*FunctionType_Parameter
}

// Java: FunctionParamBuilder(JSTypeRegistry)
func NewFunctionParamBuilder(registry *JSTypeRegistry) *FunctionParamBuilder {
	return &FunctionParamBuilder{registry: registry, parameters: []*FunctionType_Parameter{}}
}

// Java: FunctionParamBuilder(JSTypeRegistry, int)
func NewFunctionParamBuilderWithCapacity(registry *JSTypeRegistry, initialParameterCapacity int) *FunctionParamBuilder {
	return &FunctionParamBuilder{
		registry:   registry,
		parameters: make([]*FunctionType_Parameter, 0, initialParameterCapacity),
	}
}

// AddRequiredParams adds parameters of the given type to the end of the
// param list. Returns false if this is called after optional params are
// added.
//
// Java: FunctionParamBuilder.addRequiredParams(JSType...)
func (b *FunctionParamBuilder) AddRequiredParams(types ...JSType) bool {
	if b.hasOptionalOrVarArgs() {
		return false
	}
	for _, type_ := range types {
		b.newParameter(type_, /* isOptional= */ false, /* isVariadic= */ false)
	}
	return true
}

// AddOptionalParams adds optional parameters of the given type to the end
// of the param list. The builder will make them undefine-able. Returns false
// if this is called after var args are added.
//
// Java: FunctionParamBuilder.addOptionalParams(JSType...)
func (b *FunctionParamBuilder) AddOptionalParams(types ...JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	for _, type_ := range types {
		b.newParameter(
			b.registry.CreateOptionalType(type_), /* isOptional= */ true, /* isVariadic= */ false)
	}
	return true
}

// AddVarArgs adds variable arguments to the end of the parameter list.
// Returns false if this is called after var args are added.
//
// Java: FunctionParamBuilder.addVarArgs(JSType)
func (b *FunctionParamBuilder) AddVarArgs(type_ JSType) bool {
	if b.HasVarArgs() {
		return false
	}
	b.newParameter(type_, /* isOptional= */ false, /* isVariadic= */ true)
	return true
}

// NewParameterFrom copies the existing parameter into this builder.
//
// Java: FunctionParamBuilder.newParameterFrom(Parameter)
func (b *FunctionParamBuilder) NewParameterFrom(n *FunctionType_Parameter) {
	b.parameters = append(b.parameters, n)
}

// NewOptionalParameterFrom copies the parameter specification from the given
// parameter, but makes sure it's optional.
//
// Java: FunctionParamBuilder.newOptionalParameterFrom(Parameter)
func (b *FunctionParamBuilder) NewOptionalParameterFrom(p *FunctionType_Parameter) {
	isOptional := p.IsOptional() || !p.IsVariadic()
	if isOptional != p.IsOptional() {
		b.newParameter(p.GetJSType(), isOptional, p.IsVariadic())
	} else {
		b.NewParameterFrom(p)
	}
}

// newParameter adds a parameter with the given type.
//
// Java: FunctionParamBuilder.newParameter(JSType, boolean, boolean)
func (b *FunctionParamBuilder) newParameter(type_ JSType, isOptional bool, isVariadic bool) {
	parameter := FunctionType_Parameter_Create(type_, isOptional, isVariadic)
	b.parameters = append(b.parameters, parameter)
}

// Java: FunctionParamBuilder.hasOptionalOrVarArgs()
func (b *FunctionParamBuilder) hasOptionalOrVarArgs() bool {
	if len(b.parameters) == 0 {
		return false
	}
	lastParam := b.parameters[len(b.parameters)-1]
	return lastParam.IsOptional() || lastParam.IsVariadic()
}

// Java: FunctionParamBuilder.hasVarArgs()
func (b *FunctionParamBuilder) HasVarArgs() bool {
	return len(b.parameters) != 0 && b.parameters[len(b.parameters)-1].IsVariadic()
}

// Java: FunctionParamBuilder.build()
func (b *FunctionParamBuilder) Build() []*FunctionType_Parameter {
	out := make([]*FunctionType_Parameter, len(b.parameters))
	copy(out, b.parameters)
	return out
}
