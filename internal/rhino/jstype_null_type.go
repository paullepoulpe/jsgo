// Port of com/google/javascript/rhino/jstype/NullType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// NullType is the null type.
type NullType struct {
	ValueTypeBase
}

// Java: NullType(JSTypeRegistry)
func NewNullType(registry *JSTypeRegistry) *NullType {
	t := &NullType{}
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: NullType.getTypeClass()
func (t *NullType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NULL
}

// Java: NullType.isNullType()
func (t *NullType) IsNullType() bool {
	return true
}

// Java: NullType.isNullable()
func (t *NullType) IsNullable() bool {
	return true
}

// Java: NullType.matchesNumberContext()
func (t *NullType) MatchesNumberContext() bool {
	return true
}

// Java: NullType.matchesObjectContext()
func (t *NullType) MatchesObjectContext() bool {
	return false
}

// Java: NullType.matchesStringContext()
func (t *NullType) MatchesStringContext() bool {
	return true
}

// Java: NullType.restrictByNotNullOrUndefined()
func (t *NullType) RestrictByNotNullOrUndefined() JSType {
	return t.registry.GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: NullType.restrictByNotNull()
func (t *NullType) RestrictByNotNull() JSType {
	return t.registry.GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: NullType.testForEquality(JSType)
func (t *NullType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsNullType() || that.IsVoidType() {
		return base.Tri_TRUE
	}
	if that.IsUnknownType() || that.IsNullable() {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: NullType.getDisplayName()
func (t *NullType) GetDisplayName() *string {
	return &jsType_displayName_null
}

// Java: NullType.getPossibleToBooleanOutcomes()
func (t *NullType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_FALSE
}

// Java: NullType.visit(Visitor<T>)
func (t *NullType) Visit(visitor Visitor) any {
	return visitor.CaseNullType()
}

var jsType_displayName_null = "null"
