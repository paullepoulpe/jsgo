// Port of com/google/javascript/rhino/jstype/SymbolType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// SymbolType is the symbol type.
//
// Java: class SymbolType extends ValueType (subclassed by KnownSymbolType).
type SymbolType interface {
	ValueType

	// GetSymbolTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetSymbolTypeBase() *SymbolTypeBase
}

// SymbolTypeBase holds SymbolType's methods; plain symbol types are
// *SymbolTypeBase values.
type SymbolTypeBase struct {
	ValueTypeBase
}

// NewSymbolType constructs the plain symbol type.
//
// Java: SymbolType(JSTypeRegistry)
func NewSymbolType(registry *JSTypeRegistry) *SymbolTypeBase {
	t := &SymbolTypeBase{}
	t.InitSymbolTypeBase(t, registry)
	return t
}

// InitSymbolTypeBase is the Go form of the constructor for subclasses.
//
// Java: SymbolType(JSTypeRegistry)
func (t *SymbolTypeBase) InitSymbolTypeBase(self SymbolType, registry *JSTypeRegistry) {
	t.InitValueTypeBase(self, registry)
}

// GetSymbolTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *SymbolTypeBase) GetSymbolTypeBase() *SymbolTypeBase { return t }

// Java: SymbolType.getTypeClass()
func (t *SymbolTypeBase) GetTypeClass() JSTypeClass {
	return JSTypeClass_SYMBOL
}

// Java: SymbolType.testForEquality(JSType)
func (t *SymbolTypeBase) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}

	if that.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
		that.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: SymbolType.isSymbolValueType()
func (t *SymbolTypeBase) IsSymbolValueType() bool {
	return true
}

// Java: SymbolType.matchesNumberContext()
func (t *SymbolTypeBase) MatchesNumberContext() bool {
	return false
}

// Java: SymbolType.matchesStringContext()
func (t *SymbolTypeBase) MatchesStringContext() bool {
	return false
}

// Java: SymbolType.matchesSymbolContext()
func (t *SymbolTypeBase) MatchesSymbolContext() bool {
	return true
}

// Java: SymbolType.matchesObjectContext()
func (t *SymbolTypeBase) MatchesObjectContext() bool {
	return true
}

// Java: SymbolType.getDisplayName()
func (t *SymbolTypeBase) GetDisplayName() *string {
	return &jsType_displayName_symbol
}

var jsType_displayName_symbol = "symbol"

// Java: SymbolType.autoboxesTo()
func (t *SymbolTypeBase) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)
}

// Java: SymbolType.getPossibleToBooleanOutcomes()
func (t *SymbolTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_TRUE
}

// Java: SymbolType.visit(Visitor<T>)
func (t *SymbolTypeBase) Visit(visitor Visitor) any {
	return visitor.CaseSymbolType(t.self.(SymbolType))
}
