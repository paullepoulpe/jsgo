// Port of com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java @ closure-compiler v20260415
//
// ContainsUpperBoundSuperTypeVisitor traverses the referenced types of
// "forwarding types" (NamedType, TemplatizedType, TemplateType, UnionType)
// to search for a target JSType.  It is used by TemplateType as a helper
// in cycle detection and subtyping.
//
// The Result enum indicates whether the target was PRESENT, ABSENT, or
// undecidable due to a reference CYCLE.
//
// Dispatch follows the Go port convention: the package-level VisitJSType
// function replaces Java's JSType.visit(Visitor), since the JSType
// interface does not carry an accept/visit method (additive-only phase).
//
// NOTE: TemplateType.getBound() is not yet ported (Tier-4+); the
// caseForwardingType for TemplateType conservatively forwards to Unknown
// when no bound is available.

package jstype

// UpperBoundResult is the outcome of a ContainsUpperBoundSuperTypeVisitor
// visitation.
// Port of ContainsUpperBoundSuperTypeVisitor.Result.
type UpperBoundResult int

const (
	// UpperBoundAbsent — target type was not found.
	// Port of Result.ABSENT.
	UpperBoundAbsent UpperBoundResult = iota

	// UpperBoundPresent — target type was found.
	// Port of Result.PRESENT.
	UpperBoundPresent

	// UpperBoundCycle — containment is undecidable due to a reference cycle.
	// Port of Result.CYCLE.
	UpperBoundCycle
)

// ContainsUpperBoundSuperTypeVisitor searches for a target JSType by
// traversing forwarding types.
// Port of ContainsUpperBoundSuperTypeVisitor extends Visitor.WithDefaultCase<Result>.
type ContainsUpperBoundSuperTypeVisitor struct {
	WithDefaultCase
	target JSType
	seen   map[JSType]bool // mirrors LinkedIdentityHashSet<JSType>
}

// NewContainsUpperBoundSuperTypeVisitor constructs a visitor for the given
// target.  Mirrors ContainsUpperBoundSuperTypeVisitor(JSType target).
func NewContainsUpperBoundSuperTypeVisitor(target JSType) *ContainsUpperBoundSuperTypeVisitor {
	v := &ContainsUpperBoundSuperTypeVisitor{
		target: target,
		seen:   make(map[JSType]bool),
	}
	v.WithDefaultCase.DefaultCase = v.caseDefault
	return v
}

// caseDefault is the fallback for types not explicitly overridden.
// Port of Visitor.WithDefaultCase.caseDefault — returns ABSENT unless
// type is pointer-identical to target.
func (v *ContainsUpperBoundSuperTypeVisitor) caseDefault(typ JSType) any {
	if v.target == nil {
		return UpperBoundAbsent
	}
	if typ == v.target {
		return UpperBoundPresent
	}
	return UpperBoundAbsent
}

// CaseTemplateType overrides the default to follow the template's bound.
// Port of caseTemplateType — calls caseForwardingType(type, type.getBound()).
// TODO(Tier-4): replace Unknown with type.GetBound() once bounds are ported.
func (v *ContainsUpperBoundSuperTypeVisitor) CaseTemplateType(typ *TemplateType) any {
	// getBound() is not yet ported; conservatively forward to Unknown so
	// an unbound template is never mistaken for a cycle.
	return v.caseForwardingType(typ, Unknown)
}

// CaseNamedType overrides the default to follow the referenced type.
// Port of caseNamedType — calls caseForwardingType(type, type.getReferencedType()).
func (v *ContainsUpperBoundSuperTypeVisitor) CaseNamedType(typ JSType) any {
	nt, ok := typ.(*NamedType)
	if !ok {
		return UpperBoundAbsent
	}
	return v.caseForwardingType(nt, nt.GetReferencedType())
}

// CaseTemplatizedType overrides the default to follow the base type.
// Port of caseTemplatizedType — calls caseForwardingType(type, type.getReferencedType()).
// In Go, TemplatizedType.getReferencedType() corresponds to the Base field.
func (v *ContainsUpperBoundSuperTypeVisitor) CaseTemplatizedType(typ *TemplatizedType) any {
	return v.caseForwardingType(typ, typ.Base)
}

// CaseUnionType overrides the default to search each alternate.
// Port of caseUnionType.
func (v *ContainsUpperBoundSuperTypeVisitor) CaseUnionType(typ *UnionType) any {
	if typ == v.target {
		return UpperBoundPresent
	}
	for _, alt := range typ.Members() {
		result := VisitJSType(v, alt).(UpperBoundResult)
		if result != UpperBoundAbsent {
			return result
		}
	}
	return UpperBoundAbsent
}

// caseForwardingType is the common implementation for forwarding types.
// Port of ContainsUpperBoundSuperTypeVisitor.caseForwardingType.
func (v *ContainsUpperBoundSuperTypeVisitor) caseForwardingType(typ JSType, reference JSType) UpperBoundResult {
	if typ == v.target {
		return UpperBoundPresent
	}
	if v.seen[typ] {
		return UpperBoundCycle
	}
	v.seen[typ] = true
	return VisitJSType(v, reference).(UpperBoundResult)
}

// VisitJSType dispatches typ through the Visitor interface.  This is the
// Go equivalent of the Java pattern JSType.visit(Visitor).
//
// Type cases handled:
//   - *TemplateType     → CaseTemplateType
//   - *TemplatizedType  → CaseTemplatizedType
//   - *UnionType        → CaseUnionType
//   - *FunctionType     → CaseFunctionType
//   - *ObjectType       → CaseObjectType
//   - *EnumElementType  → CaseEnumElementType
//   - *NamedType        → CaseNamedType
//   - unknownType       → CaseUnknownType
//   - allType           → CaseAllType
//   - noType / nil      → CaseNoType
//   - all other primitives → CaseBooleanType / CaseNumberType / etc.
//
// NOTE: The JSType.visit / JSType.accept wiring is intentionally absent
// from the JSType interface (additive-only phase). Wire when actively consumed.
func VisitJSType(v Visitor, typ JSType) any {
	if typ == nil {
		return v.CaseNoType(nil)
	}
	switch t := typ.(type) {
	case *TemplateType:
		return v.CaseTemplateType(t)
	case *TemplatizedType:
		return v.CaseTemplatizedType(t)
	case *UnionType:
		return v.CaseUnionType(t)
	case *FunctionType:
		return v.CaseFunctionType(t)
	case *ObjectType:
		return v.CaseObjectType(t)
	case *EnumElementType:
		return v.CaseEnumElementType(t)
	case *NamedType:
		return v.CaseNamedType(t)
	case unknownType:
		return v.CaseUnknownType()
	case allType:
		return v.CaseAllType()
	case numberType:
		return v.CaseNumberType()
	case stringType:
		return v.CaseStringType()
	case booleanType:
		return v.CaseBooleanType()
	case nullType:
		return v.CaseNullType()
	case voidType:
		return v.CaseVoidType()
	default:
		// noType and other not-yet-ported concrete types.
		return v.CaseNoType(typ)
	}
}
