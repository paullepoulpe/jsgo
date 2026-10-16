// Port of com/google/javascript/rhino/jstype/AllType.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// AllType is the all type, representing all values.
type AllType struct {
	JSTypeBase
}

// Java: AllType(JSTypeRegistry)
func NewAllType(registry *JSTypeRegistry) *AllType {
	t := &AllType{}
	t.InitJSTypeBase(t, registry)
	t.EagerlyResolveToSelf()
	return t
}

// Java: AllType.getTypeClass()
func (t *AllType) GetTypeClass() JSTypeClass {
	return JSTypeClass_ALL
}

// Java: AllType.isAllType()
func (t *AllType) IsAllType() bool {
	return true
}

// Java: AllType.matchesStringContext()
func (t *AllType) MatchesStringContext() bool {
	// Be lenient.
	return true
}

// Java: AllType.matchesObjectContext()
func (t *AllType) MatchesObjectContext() bool {
	// Be lenient.
	return true
}

// Java: AllType.testForEquality(JSType)
func (t *AllType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: AllType.appendTo(TypeStringBuilder)
func (t *AllType) AppendTo(sb *TypeStringBuilder) {
	sb.Append("*")
}

// Java: AllType.getDisplayName()
func (t *AllType) GetDisplayName() *string {
	return &allType_displayName
}

var allType_displayName = "<Any Type>"

// Java: AllType.hasDisplayName()
func (t *AllType) HasDisplayName() bool {
	return true
}

// Java: AllType.visit(Visitor<T>)
func (t *AllType) Visit(visitor Visitor) any {
	return visitor.CaseAllType()
}

// Java: AllType.visit(RelationshipVisitor<T>, JSType)
func (t *AllType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseAllType(that)
}

// Java: AllType.getPossibleToBooleanOutcomes()
func (t *AllType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: AllType.resolveInternal(ErrorReporter)
func (t *AllType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: AllType.recursionUnsafeHashCode()
func (t *AllType) RecursionUnsafeHashCode() int32 {
	return jutil.IdentityHashCode(t)
}

// Java: AllType.isNullable()
func (t *AllType) IsNullable() bool {
	return true
}

// Java: AllType.isVoidable()
func (t *AllType) IsVoidable() bool {
	return true
}
