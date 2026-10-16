// Port of com/google/javascript/rhino/jstype/SimpleSlot.java @ closure-compiler v20260415

package rhino

// SimpleSlot is the minimum implementation of StaticTypedSlot.
//
// Java: class SimpleSlot (subclassed by jscomp's SymbolTable.Symbol).
type SimpleSlot interface {
	StaticTypedSlot
	GetSimpleSlotBase() *SimpleSlotBase
}

// SimpleSlotBase holds SimpleSlot's fields and methods; plain slots are
// *SimpleSlotBase values. The base methods make no virtual calls, so no self
// field is needed.
type SimpleSlotBase struct {
	name     string
	type_    JSType
	inferred bool
}

// Java: SimpleSlot.serialVersionUID
const SimpleSlot_serialVersionUID int64 = 1

// NewSimpleSlot constructs a plain SimpleSlot.
//
// Java: SimpleSlot(String, JSType, boolean)
func NewSimpleSlot(name string, type_ JSType, inferred bool) *SimpleSlotBase {
	s := &SimpleSlotBase{}
	s.InitSimpleSlotBase(name, type_, inferred)
	return s
}

// InitSimpleSlotBase is the Go form of the constructor for subclasses.
//
// Java: SimpleSlot(String, JSType, boolean)
func (s *SimpleSlotBase) InitSimpleSlotBase(name string, type_ JSType, inferred bool) {
	s.name = name
	s.type_ = type_
	s.inferred = inferred
}

// GetSimpleSlotBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (s *SimpleSlotBase) GetSimpleSlotBase() *SimpleSlotBase { return s }

// Java: SimpleSlot.getName()
func (s *SimpleSlotBase) GetName() string {
	return s.name
}

// Java: SimpleSlot.getType()
func (s *SimpleSlotBase) GetType() JSType {
	return s.type_
}

// Java: SimpleSlot.isTypeInferred()
func (s *SimpleSlotBase) IsTypeInferred() bool {
	return s.inferred
}

// Java: SimpleSlot.getDeclaration()
func (s *SimpleSlotBase) GetDeclaration() StaticRef {
	return nil
}

// Java: SimpleSlot.getJSDocInfo()
func (s *SimpleSlotBase) GetJSDocInfo() *JSDocInfo {
	return nil
}

// Java: SimpleSlot.getScope()
func (s *SimpleSlotBase) GetScope() StaticScope {
	return nil
}
