// Port of com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// ContainsUpperBoundSuperTypeVisitor is a type visitor that traverse through
// the referenced types of "forwaring types" to search for a target type.
//
// The result of the visitation indicates whether the target type was
// present, absent, or undcidable due to a cycle.
//
// This is used by TemplateType as a helper in cycle detection and
// subtyping. Cycles are checked for once after TemplateTypes are constructed
// in FunctionTypeBuilder and an error is emitted if they are found.
//
// Java: extends Visitor.WithDefaultCase<Result>; the case methods return
// Results boxed in `any` (see Visitor).
type ContainsUpperBoundSuperTypeVisitor struct {
	Visitor_WithDefaultCaseBase

	target JSType
	seen   *base.LinkedIdentityHashSet[JSType]
}

// Java: ContainsUpperBoundSuperTypeVisitor(JSType)
func NewContainsUpperBoundSuperTypeVisitor(target JSType) *ContainsUpperBoundSuperTypeVisitor {
	v := &ContainsUpperBoundSuperTypeVisitor{}
	v.InitVisitor_WithDefaultCaseBase(v)
	v.seen = base.NewLinkedIdentityHashSet[JSType]()
	v.target = target
	return v
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseDefault(JSType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseDefault(type_ JSType) any {
	if v.target == nil {
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	}

	if type_ == v.target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	}
	return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseTemplateType(TemplateType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseTemplateType(type_ *TemplateType) any {
	return v.CaseForwardingType(type_, type_.GetBound())
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseNamedType(NamedType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseNamedType(type_ *NamedType) any {
	return v.CaseForwardingType(type_, type_.GetReferencedType())
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseTemplatizedType(TemplatizedType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseTemplatizedType(type_ *TemplatizedType) any {
	return v.CaseForwardingType(type_, type_.GetReferencedType())
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseUnionType(UnionType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseUnionType(type_ *UnionType) any {
	if JSType(type_) == v.target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	}

	for _, alt := range type_.GetAlternates() {
		foundInAlt := alt.Visit(v).(ContainsUpperBoundSuperTypeVisitor_Result)
		if foundInAlt != ContainsUpperBoundSuperTypeVisitor_Result_ABSENT {
			return foundInAlt
		}
	}

	return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseForwardingType(JSType, JSType)
func (v *ContainsUpperBoundSuperTypeVisitor) CaseForwardingType(type_ JSType, reference JSType) any {
	if type_ == v.target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	} else if v.seen.Contains(type_) {
		return ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
	} else {
		v.seen.Add(type_)
		return reference.Visit(v)
	}
}

// ContainsUpperBoundSuperTypeVisitor_Result represents the outcome of a
// visitation of the ContainsUpperBoundSuperTypeVisitor.
type ContainsUpperBoundSuperTypeVisitor_Result int32

const (
	ContainsUpperBoundSuperTypeVisitor_Result_PRESENT ContainsUpperBoundSuperTypeVisitor_Result = iota
	ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	// Containment is undecidable due to a reference cycle.
	ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
)

var containsUpperBoundSuperTypeVisitor_Result_names = [...]string{"PRESENT", "ABSENT", "CYCLE"}

// String is Enum.name().
func (r ContainsUpperBoundSuperTypeVisitor_Result) String() string {
	return containsUpperBoundSuperTypeVisitor_Result_names[r]
}
