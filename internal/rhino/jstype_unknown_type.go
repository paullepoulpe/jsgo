// Port of com/google/javascript/rhino/jstype/UnknownType.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// UnknownType is the unknown type.
type UnknownType struct {
	ObjectTypeBase

	// See the explanation of checked unknown types in JSTypeNative.
	isChecked bool
}

// Java: UnknownType(JSTypeRegistry, boolean)
func NewUnknownType(registry *JSTypeRegistry, isChecked bool) *UnknownType {
	t := &UnknownType{}
	t.InitObjectTypeBase(t, registry)
	t.EagerlyResolveToSelf()
	t.isChecked = isChecked
	return t
}

// Java: UnknownType.getTypeClass()
func (t *UnknownType) GetTypeClass() JSTypeClass {
	return JSTypeClass_UNKNOWN
}

// Java: UnknownType.isUnknownType()
func (t *UnknownType) IsUnknownType() bool {
	return true
}

// Java: UnknownType.isCheckedUnknownType()
func (t *UnknownType) IsCheckedUnknownType() bool {
	return t.isChecked
}

// Java: UnknownType.canBeCalled()
func (t *UnknownType) CanBeCalled() bool {
	return true
}

// Java: UnknownType.matchesNumberContext()
func (t *UnknownType) MatchesNumberContext() bool {
	return true
}

// Java: UnknownType.matchesObjectContext()
func (t *UnknownType) MatchesObjectContext() bool {
	return true
}

// Java: UnknownType.matchesStringContext()
func (t *UnknownType) MatchesStringContext() bool {
	return true
}

// Java: UnknownType.matchesSymbolContext()
func (t *UnknownType) MatchesSymbolContext() bool {
	return true
}

// Java: UnknownType.testForEquality(JSType)
func (t *UnknownType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: UnknownType.isNullable()
func (t *UnknownType) IsNullable() bool {
	return true
}

// Java: UnknownType.isVoidable()
func (t *UnknownType) IsVoidable() bool {
	return true
}

// Java: UnknownType.visit(Visitor<T>)
func (t *UnknownType) Visit(visitor Visitor) any {
	return visitor.CaseUnknownType()
}

// Java: UnknownType.visit(RelationshipVisitor<T>, JSType)
func (t *UnknownType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseUnknownType(t, that)
}

// Java: UnknownType.appendTo(TypeStringBuilder)
func (t *UnknownType) AppendTo(sb *TypeStringBuilder) {
	sb.Append(*t.GetReferenceName())
}

// Java: UnknownType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *UnknownType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool {
	// nothing to define
	return true
}

// Java: UnknownType.getImplicitPrototype()
func (t *UnknownType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java: UnknownType.getConstructor()
func (t *UnknownType) GetConstructor() FunctionType {
	return nil
}

// Java: UnknownType.getReferenceName()
func (t *UnknownType) GetReferenceName() *string {
	if t.isChecked {
		return &unknownType_checkedName
	}
	return &unknownType_name
}

var (
	unknownType_checkedName = "??"
	unknownType_name        = "?"
	unknownType_displayName = "Unknown"
)

// Java: UnknownType.getDisplayName()
func (t *UnknownType) GetDisplayName() *string {
	return &unknownType_displayName
}

// Java: UnknownType.hasDisplayName()
func (t *UnknownType) HasDisplayName() bool {
	return true
}

// Java: UnknownType.getPossibleToBooleanOutcomes()
func (t *UnknownType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: UnknownType.resolveInternal(ErrorReporter)
func (t *UnknownType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: UnknownType.recursionUnsafeHashCode()
func (t *UnknownType) RecursionUnsafeHashCode() int32 {
	return jutil.IdentityHashCode(t)
}
