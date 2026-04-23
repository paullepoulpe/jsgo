// Port of com/google/javascript/rhino/jstype/ArrowType.java @ closure-compiler v20260415
//
// ArrowType models the [[Call]] slot of a function: an ordered parameter
// list plus a return type. In Java this is a package-private JSType
// subclass; here it is a standalone struct implementing the JSType
// interface.
//
// FunctionType currently carries these fields inline (Params/Return). A
// later migration step will wire FunctionType.call to use *ArrowType.
// Nothing in this file modifies function_type.go.

package jstype

import "strings"

// Parameter and NewParameter are defined in function_type.go (the
// earlier Parameter widening landed those there). ArrowType uses the
// same struct to avoid two different `Parameter` definitions in the
// jstype package.

// ArrowType models a "bare" function type: a parameter list and a return
// type. JavaScript functions carry more (this-type, prototype, own
// properties); those reside on FunctionType. Mirrors ArrowType in Java.
type ArrowType struct {
	// parameterList holds the ordered formal parameters. Never nil after
	// construction; defaults to a single var-args Unknown parameter when
	// the caller passes nil (mirrors the Java constructor behaviour).
	parameterList []Parameter

	// returnType is the declared return type. Defaults to Unknown when
	// the Java constructor receives nil.
	returnType JSType

	// ReturnTypeInferred mirrors ArrowType.returnTypeInferred — true when
	// the return type was inferred rather than explicitly declared via
	// @return. Exported so FunctionType (and later callers) can read it.
	ReturnTypeInferred bool
}

// NewArrowType constructs an ArrowType. A nil parameters slice is
// replaced with a single var-args Unknown parameter (matching Java's
// registry.createParametersWithVarArgs(getNativeType(UNKNOWN_TYPE))).
// A nil returnType is replaced with Unknown. Mirrors the Java constructor.
func NewArrowType(parameters []Parameter, returnType JSType, returnTypeInferred bool) *ArrowType {
	if parameters == nil {
		parameters = []Parameter{NewParameter(Unknown, false, true)}
	}
	if returnType == nil {
		returnType = Unknown
	}
	return &ArrowType{
		parameterList:      parameters,
		returnType:         returnType,
		ReturnTypeInferred: returnTypeInferred,
	}
}

// GetReturnType mirrors ArrowType.getReturnType().
func (a *ArrowType) GetReturnType() JSType { return a.returnType }

// GetParameterList mirrors ArrowType.getParameterList().
func (a *ArrowType) GetParameterList() []Parameter { return a.parameterList }

// HasUnknownParamsOrReturn mirrors ArrowType.hasUnknownParamsOrReturn().
// Returns true when any parameter type or the return type is Unknown.
func (a *ArrowType) HasUnknownParamsOrReturn() bool {
	for _, p := range a.parameterList {
		if isUnknown(p.GetJSType()) {
			return true
		}
	}
	return a.returnType == nil || isUnknown(a.returnType)
}

// HasAnyTemplateTypes mirrors ArrowType.hasAnyTemplateTypesInternal().
// Returns true when the return type or any parameter type contains an
// unresolved TemplateType.
func (a *ArrowType) HasAnyTemplateTypes() bool {
	return hasAnyTemplateTypes(a.returnType) || a.hasTemplatedParameterType()
}

// hasTemplatedParameterType mirrors ArrowType.hasTemplatedParameterType().
func (a *ArrowType) hasTemplatedParameterType() bool {
	for _, p := range a.parameterList {
		if hasAnyTemplateTypes(p.GetJSType()) {
			return true
		}
	}
	return false
}

// hasAnyTemplateTypes reports whether t is or transitively contains an
// unresolved *TemplateType. Mirrors JSType.hasAnyTemplateTypes() for the
// types that ArrowType uses: TemplateType, UnionType, TemplatizedType,
// FunctionType, and ArrowType itself.
func hasAnyTemplateTypes(t JSType) bool {
	if t == nil {
		return false
	}
	switch v := t.(type) {
	case *TemplateType:
		return true
	case *UnionType:
		for _, m := range v.Members() {
			if hasAnyTemplateTypes(m) {
				return true
			}
		}
		return false
	case *TemplatizedType:
		for _, a := range v.Args {
			if hasAnyTemplateTypes(a) {
				return true
			}
		}
		return false
	case *FunctionType:
		for _, p := range v.Params {
			if hasAnyTemplateTypes(p) {
				return true
			}
		}
		return hasAnyTemplateTypes(v.Return)
	case *ArrowType:
		return v.HasAnyTemplateTypes()
	}
	return false
}

// RecursionUnsafeHashCode mirrors ArrowType.recursionUnsafeHashCode().
// Uses TypeName() as a proxy for Objects.hashCode since Go has no
// built-in object identity hash for interface values.
func (a *ArrowType) RecursionUnsafeHashCode() int {
	h := typePtrHash(a.returnType)
	for _, p := range a.parameterList {
		h = h*31 + typePtrHash(p.GetJSType())
	}
	return h
}

// typePtrHash returns a simple hash for a JSType value, analogous to
// Objects.hashCode(t) in Java (which returns 0 for nil, identity hash
// otherwise). We use the TypeName string as a stable stand-in.
func typePtrHash(t JSType) int {
	if t == nil {
		return 0
	}
	s := t.TypeName()
	h := 0
	for _, c := range s {
		h = h*31 + int(c)
	}
	return h
}

// TypeName mirrors ArrowType.toString() and satisfies the JSType
// interface. Renders as `(T1,T2) -> R`.
func (a *ArrowType) TypeName() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, p := range a.parameterList {
		if i > 0 {
			sb.WriteString(",")
		}
		if p.GetJSType() != nil {
			sb.WriteString(p.GetJSType().TypeName())
		} else {
			sb.WriteString("?")
		}
	}
	sb.WriteString(") -> ")
	if a.returnType != nil {
		sb.WriteString(a.returnType.TypeName())
	} else {
		sb.WriteString("?")
	}
	return sb.String()
}

// IsSubtypeOf satisfies the JSType interface. ArrowType is not a
// first-class JSType in Java (getLeastSupertype / getGreatestSubtype /
// testForEquality all throw UnsupportedOperationException there), so
// IsSubtypeOf follows the conservative rule: only AllType and Unknown
// are supertypes.
func (a *ArrowType) IsSubtypeOf(other JSType) bool {
	return isAll(other) || isUnknown(other)
}
