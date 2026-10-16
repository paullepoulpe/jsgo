// Port of com/google/javascript/rhino/jstype/KnownSymbolType.java @ closure-compiler v20260415

package rhino

// KnownSymbolType is the type of a well-known (or unique) symbol, such as
// Symbol.iterator.
type KnownSymbolType struct {
	SymbolTypeBase

	name string
}

// Java: KnownSymbolType.TYPE_CLASS
const KnownSymbolType_TYPE_CLASS = JSTypeClass_WELL_KNOWN_SYMBOL

// Java: KnownSymbolType(JSTypeRegistry, String)
func NewKnownSymbolType(registry *JSTypeRegistry, name string) *KnownSymbolType {
	t := &KnownSymbolType{}
	t.InitSymbolTypeBase(t, registry)
	t.name = name
	return t
}

// Java: KnownSymbolType.getTypeClass()
func (t *KnownSymbolType) GetTypeClass() JSTypeClass {
	return KnownSymbolType_TYPE_CLASS
}

// Java: KnownSymbolType.getDisplayName()
func (t *KnownSymbolType) GetDisplayName() *string {
	return &t.name
}

// Java: KnownSymbolType.isKnownSymbolValueType()
func (t *KnownSymbolType) IsKnownSymbolValueType() bool {
	return true
}

// Java: KnownSymbolType.toMaybeKnownSymbolType()
func (t *KnownSymbolType) ToMaybeKnownSymbolType() *KnownSymbolType {
	return t
}
