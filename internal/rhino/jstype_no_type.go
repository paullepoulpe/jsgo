// Port of com/google/javascript/rhino/jstype/NoType.java @ closure-compiler v20260415

package rhino

// NoType is the bottom type, representing the subclass of any value or
// object.
//
// Although JavaScript programmers can't explicitly denote the bottom type,
// it comes up in static analysis. For example, if we have:
//
//	var x = null;
//	if (x) {
//	  f(x);
//	}
//
// We need to be able to assign x a type within the f(x) call. Since it has
// no possible type, we assign x the NoType, so that f(x) is legal no matter
// what the type of f's first argument is.
type NoType struct {
	NoObjectTypeBase
}

// Java: NoType(JSTypeRegistry)
func NewNoType(registry *JSTypeRegistry) *NoType {
	t := &NoType{}
	t.InitNoObjectTypeBase(t, registry)
	return t
}

// Java: NoType.getTypeClass()
func (t *NoType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NO
}

// Java: NoType.isNoObjectType()
func (t *NoType) IsNoObjectType() bool {
	return false
}

// Java: NoType.isNoType()
func (t *NoType) IsNoType() bool {
	return true
}

// Java: NoType.isNullable()
func (t *NoType) IsNullable() bool {
	return true
}

// Java: NoType.isVoidable()
func (t *NoType) IsVoidable() bool {
	return true
}

// Java: NoType.getPossibleToBooleanOutcomes()
func (t *NoType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NoType.visit(Visitor<T>)
func (t *NoType) Visit(visitor Visitor) any {
	return visitor.CaseNoType(t)
}

// Java: NoType.visit(RelationshipVisitor<T>, JSType)
func (t *NoType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseNoType(that)
}

// Java: NoType.appendTo(TypeStringBuilder)
func (t *NoType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("?")
	} else {
		sb.Append("None")
	}
}
