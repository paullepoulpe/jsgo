// Port of com/google/javascript/rhino/jstype/NumberType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// NumberType is the number type.
type NumberType struct {
	ValueTypeBase
}

// Java: NumberType(JSTypeRegistry)
func NewNumberType(registry *JSTypeRegistry) *NumberType {
	t := &NumberType{}
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: NumberType.getTypeClass()
func (t *NumberType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NUMBER
}

// Java: NumberType.testForEquality(JSType)
func (t *NumberType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BOOLEAN_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: NumberType.isNumberValueType()
func (t *NumberType) IsNumberValueType() bool {
	return true
}

// Java: NumberType.matchesNumberContext()
func (t *NumberType) MatchesNumberContext() bool {
	return true
}

// Java: NumberType.matchesStringContext()
func (t *NumberType) MatchesStringContext() bool {
	return true
}

// Java: NumberType.matchesObjectContext()
func (t *NumberType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: NumberType.getDisplayName()
func (t *NumberType) GetDisplayName() *string {
	return &jsType_displayName_number
}

// Java: NumberType.getPossibleToBooleanOutcomes()
func (t *NumberType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NumberType.visit(Visitor<T>)
func (t *NumberType) Visit(visitor Visitor) any {
	return visitor.CaseNumberType()
}

// Java: NumberType.autoboxesTo()
func (t *NumberType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_NUMBER_OBJECT_TYPE)
}

var jsType_displayName_number = "number"
