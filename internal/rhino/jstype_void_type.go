// Port of com/google/javascript/rhino/jstype/VoidType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// VoidType is the void type.
type VoidType struct {
	ValueTypeBase
}

// Java: VoidType(JSTypeRegistry)
func NewVoidType(registry *JSTypeRegistry) *VoidType {
	t := &VoidType{}
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: VoidType.getTypeClass()
func (t *VoidType) GetTypeClass() JSTypeClass {
	return JSTypeClass_VOID
}

// Java: VoidType.restrictByNotNullOrUndefined()
func (t *VoidType) RestrictByNotNullOrUndefined() JSType {
	return t.registry.GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: VoidType.restrictByNotUndefined()
func (t *VoidType) RestrictByNotUndefined() JSType {
	return t.registry.GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: VoidType.testForEquality(JSType)
func (t *VoidType) TestForEquality(that JSType) base.Tri {
	if base.Tri_UNKNOWN == t.ValueTypeBase.TestForEquality(that) {
		return base.Tri_UNKNOWN
	}
	if that.IsSubtypeOf(t) || that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NULL_TYPE)) {
		return base.Tri_TRUE
	}
	return base.Tri_FALSE
}

// Java: VoidType.matchesNumberContext()
func (t *VoidType) MatchesNumberContext() bool {
	return false
}

// Java: VoidType.matchesObjectContext()
func (t *VoidType) MatchesObjectContext() bool {
	return false
}

// Java: VoidType.matchesStringContext()
func (t *VoidType) MatchesStringContext() bool {
	return true
}

// Java: VoidType.isVoidType()
func (t *VoidType) IsVoidType() bool {
	return true
}

// Java: VoidType.isVoidable()
func (t *VoidType) IsVoidable() bool {
	return true
}

// Java: VoidType.isExplicitlyVoidable()
func (t *VoidType) IsExplicitlyVoidable() bool {
	return true
}

// Java: VoidType.getDisplayName()
func (t *VoidType) GetDisplayName() *string {
	return &jsType_displayName_undefined
}

// Java: VoidType.getPossibleToBooleanOutcomes()
func (t *VoidType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_FALSE
}

// Java: VoidType.visit(Visitor<T>)
func (t *VoidType) Visit(visitor Visitor) any {
	return visitor.CaseVoidType()
}

var jsType_displayName_undefined = "undefined"
