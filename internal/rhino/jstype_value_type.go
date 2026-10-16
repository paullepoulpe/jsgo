// Port of com/google/javascript/rhino/jstype/ValueType.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ValueType is the class of value types (null, void, number, boolean,
// string).
//
// Java: abstract class ValueType extends JSType.
type ValueType interface {
	JSType

	// GetValueTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetValueTypeBase() *ValueTypeBase
}

// ValueTypeBase holds ValueType's base method implementations; the virtual
// calls go through JSTypeBase.self.
type ValueTypeBase struct {
	JSTypeBase
}

// InitValueTypeBase is the Go form of the ValueType constructor.
//
// Java: ValueType(JSTypeRegistry)
func (v *ValueTypeBase) InitValueTypeBase(self ValueType, registry *JSTypeRegistry) {
	v.InitJSTypeBase(self, registry)
	v.EagerlyResolveToSelf()
}

// GetValueTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (v *ValueTypeBase) GetValueTypeBase() *ValueTypeBase { return v }

// Java: ValueType.resolveInternal(ErrorReporter)
func (v *ValueTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: ValueType.appendTo(TypeStringBuilder)
func (v *ValueTypeBase) AppendTo(sb *TypeStringBuilder) {
	sb.Append(jutil.ToJavaString(v.self.GetDisplayName()))
}

// Java: ValueType.hasDisplayName()
func (v *ValueTypeBase) HasDisplayName() bool {
	return true
}

// Java: ValueType.visit(RelationshipVisitor<T>, JSType)
func (v *ValueTypeBase) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseValueType(v.self.(ValueType), that)
}

// Java: ValueType.getPropertyKind(Property.Key, boolean)
func (v *ValueTypeBase) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	if autobox && v.self.IsBoxableScalar() {
		return v.self.AutoboxesTo().GetPropertyKindByKeyWithAutobox(propertyName, autobox)
	} else {
		return v.JSTypeBase.GetPropertyKindByKeyWithAutobox(propertyName, autobox)
	}
}

// Java: ValueType.recursionUnsafeHashCode()
func (v *ValueTypeBase) RecursionUnsafeHashCode() int32 {
	// Subclasses of this type are unique within a JSTypeRegistry.
	return jutil.IdentityHashCode(v.self)
}
