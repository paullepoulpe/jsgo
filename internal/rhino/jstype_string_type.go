// Port of com/google/javascript/rhino/jstype/StringType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// StringType is the string type.
type StringType struct {
	ValueTypeBase
}

// Java: StringType(JSTypeRegistry)
func NewStringType(registry *JSTypeRegistry) *StringType {
	t := &StringType{}
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: StringType.getTypeClass()
func (t *StringType) GetTypeClass() JSTypeClass {
	return JSTypeClass_STRING
}

// Java: StringType.testForEquality(JSType)
func (t *StringType) TestForEquality(that JSType) base.Tri {
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

// Java: StringType.isStringValueType()
func (t *StringType) IsStringValueType() bool {
	return true
}

// Java: StringType.matchesNumberContext()
func (t *StringType) MatchesNumberContext() bool {
	return true
}

// Java: StringType.matchesStringContext()
func (t *StringType) MatchesStringContext() bool {
	return true
}

// Java: StringType.matchesObjectContext()
func (t *StringType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: StringType.getDisplayName()
func (t *StringType) GetDisplayName() *string {
	return &jsType_displayName_string
}

// Java: StringType.autoboxesTo()
func (t *StringType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: StringType.getPossibleToBooleanOutcomes()
func (t *StringType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: StringType.visit(Visitor<T>)
func (t *StringType) Visit(visitor Visitor) any {
	return visitor.CaseStringType()
}

var jsType_displayName_string = "string"
