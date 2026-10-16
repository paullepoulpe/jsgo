// Port of com/google/javascript/rhino/jstype/BigIntType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// BigIntType is the bigint type.
type BigIntType struct {
	ValueTypeBase
}

// Java: BigIntType(JSTypeRegistry)
func NewBigIntType(registry *JSTypeRegistry) *BigIntType {
	t := &BigIntType{}
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: BigIntType.getTypeClass()
func (t *BigIntType) GetTypeClass() JSTypeClass {
	return JSTypeClass_BIGINT
}

// Java: BigIntType.testForEquality(JSType)
func (t *BigIntType) TestForEquality(that JSType) base.Tri {
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

// Java: BigIntType.isBigIntValueType()
func (t *BigIntType) IsBigIntValueType() bool {
	return true
}

// Java: BigIntType.matchesNumberContext()
func (t *BigIntType) MatchesNumberContext() bool {
	// BigInt cannot be coerced to a Number
	return false
}

// Java: BigIntType.matchesStringContext()
func (t *BigIntType) MatchesStringContext() bool {
	return true
}

// Java: BigIntType.matchesObjectContext()
func (t *BigIntType) MatchesObjectContext() bool {
	return true
}

// Java: BigIntType.getDisplayName()
func (t *BigIntType) GetDisplayName() *string {
	return &jsType_displayName_bigint
}

// Java: BigIntType.getPossibleToBooleanOutcomes()
func (t *BigIntType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: BigIntType.visit(Visitor<T>)
func (t *BigIntType) Visit(visitor Visitor) any {
	return visitor.CaseBigIntType()
}

// Java: BigIntType.autoboxesTo()
func (t *BigIntType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_BIGINT_OBJECT_TYPE)
}

var jsType_displayName_bigint = "bigint"
