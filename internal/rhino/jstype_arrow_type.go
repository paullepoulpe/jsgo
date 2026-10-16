// Port of com/google/javascript/rhino/jstype/ArrowType.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ArrowType models a "bare" function type: from some parameter types to a
// return type.
//
// JavaScript functions include more things like properties, type of THIS,
// etc, and are modeled by FunctionType.
type ArrowType struct {
	JSTypeBase

	parameterList []*FunctionType_Parameter
	returnType    JSType

	// Whether the return type is inferred.
	returnTypeInferred bool
}

// Java: ArrowType(JSTypeRegistry, List<Parameter>, JSType, boolean)
func NewArrowType(
	registry *JSTypeRegistry,
	parameters []*FunctionType_Parameter,
	returnType JSType,
	returnTypeInferred bool) *ArrowType {
	t := &ArrowType{}
	t.InitJSTypeBase(t, registry)

	if parameters == nil {
		t.parameterList = registry.CreateParametersWithVarArgs(t.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
	} else {
		t.parameterList = arrowTypeCopyParameters(parameters)
	}
	if returnType == nil {
		t.returnType = t.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	} else {
		t.returnType = returnType
	}
	t.returnTypeInferred = returnTypeInferred

	registry.GetResolver().ResolveIfClosed(t, JSTypeClass_ARROW)
	return t
}

// arrowTypeCopyParameters is ImmutableList.copyOf(parameters), which
// rejects null elements.
func arrowTypeCopyParameters(parameters []*FunctionType_Parameter) []*FunctionType_Parameter {
	out := make([]*FunctionType_Parameter, len(parameters))
	for i, p := range parameters {
		out[i] = jutil.CheckNotNull(p)
	}
	return out
}

// Java: ArrowType.getTypeClass()
func (t *ArrowType) GetTypeClass() JSTypeClass {
	return JSTypeClass_ARROW
}

// Java: ArrowType.recursionUnsafeHashCode()
func (t *ArrowType) RecursionUnsafeHashCode() int32 {
	hashCode := jutil.ObjectHashCode(t.returnType)
	for _, param := range t.parameterList {
		hashCode = hashCode*31 + jutil.ObjectHashCode(param.GetJSType())
	}
	return hashCode
}

// Java: ArrowType.getReturnType()
func (t *ArrowType) GetReturnType() JSType {
	return t.returnType
}

// Java: ArrowType.getParameterList()
func (t *ArrowType) GetParameterList() []*FunctionType_Parameter {
	return t.parameterList
}

// IsReturnTypeInferred reads the package-private field returnTypeInferred.
//
// Go-only: language-adapter — field read by FunctionType.
func (t *ArrowType) IsReturnTypeInferred() bool {
	return t.returnTypeInferred
}

// Java: ArrowType.getLeastSupertype(JSType)
func (t *ArrowType) GetLeastSupertype(that JSType) JSType {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ArrowType.getGreatestSubtype(JSType)
func (t *ArrowType) GetGreatestSubtype(that JSType) JSType {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ArrowType.testForEquality(JSType)
func (t *ArrowType) TestForEquality(that JSType) base.Tri {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ArrowType.visit(Visitor<T>)
func (t *ArrowType) Visit(visitor Visitor) any {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ArrowType.visit(RelationshipVisitor<T>, JSType)
func (t *ArrowType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ArrowType.getPossibleToBooleanOutcomes()
func (t *ArrowType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_TRUE
}

// Java: ArrowType.resolveInternal(ErrorReporter)
func (t *ArrowType) ResolveInternal(reporter ErrorReporter) JSType {
	t.returnType = JSType_SafeResolve(t.returnType, reporter)
	for _, param := range t.parameterList {
		param.GetJSType().Resolve(reporter)
	}

	return t
}

// Java: ArrowType.hasUnknownParamsOrReturn()
func (t *ArrowType) HasUnknownParamsOrReturn() bool {
	for _, param := range t.parameterList {
		type_ := param.GetJSType()
		if type_.IsUnknownType() {
			return true
		}
	}

	return t.returnType == nil || t.returnType.IsUnknownType()
}

// Java: ArrowType.appendTo(TypeStringBuilder)
func (t *ArrowType) AppendTo(sb *TypeStringBuilder) {
	sb.Append("[ArrowType]")
}

// Java: ArrowType.hasAnyTemplateTypesInternal()
func (t *ArrowType) HasAnyTemplateTypesInternal() bool {
	return t.returnType.HasAnyTemplateTypes() || t.hasTemplatedParameterType()
}

// Java: ArrowType.hasTemplatedParameterType()
func (t *ArrowType) hasTemplatedParameterType() bool {
	for _, param := range t.parameterList {
		type_ := param.GetJSType()
		if type_.HasAnyTemplateTypes() {
			return true
		}
	}
	return false
}

// String is a string representation of this type, suitable for printing in
// warnings.
//
// Java: ArrowType.toString()
func (t *ArrowType) String() string {
	typeStringBuilder := NewTypeStringBuilder(false)
	typeStringBuilder.Append("(")
	elements := make([]any, len(t.parameterList))
	for i, p := range t.parameterList {
		elements[i] = p
	}
	typeStringBuilder.AppendAll(elements, ",")
	typeStringBuilder.Append(") -> ")
	typeStringBuilder.AppendType(t.returnType)
	return typeStringBuilder.Build()
}
