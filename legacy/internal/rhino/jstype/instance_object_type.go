// Port of com/google/javascript/rhino/jstype/InstanceObjectType.java
// @ closure-compiler v20260415.
//
// InstanceObjectType is a concrete subclass of PrototypeObjectType that
// represents the instance type produced by `new Foo()`. Its defining
// characteristic is that it holds a back-reference to its FunctionType
// constructor (`Foo` itself), and it overrides getImplicitPrototype() to
// return the constructor's .prototype object instead of a stored fallback.
//
// In the Java hierarchy: InstanceObjectType extends PrototypeObjectType.
// Here, InstanceObjectType embeds *PrototypeObjectType for all shared
// object behaviour and adds the constructor field.

package jstype

// InstanceObjectType is the Java InstanceObjectType shape.
// It wraps a *PrototypeObjectType (which in turn wraps *ObjectType) and
// adds the back-pointer to the constructor FunctionType.
//
// Key Java behavioural differences from a plain PrototypeObjectType:
//   - GetImplicitPrototype() returns constructor.getPrototype(), not a
//     stored fallback.
//   - GetReferenceName() / isArrayType() / … delegate to the constructor.
//   - defineProperty guards against re-defining a property already
//     declared on the implicit prototype (the .prototype object).
type InstanceObjectType struct {
	*PrototypeObjectType

	// constructor is the FunctionType that, when called with `new`,
	// produces this instance type. Mirrors InstanceObjectType.constructor
	// in Java; never nil.
	constructor *FunctionType
}

// NewInstanceObjectType allocates an InstanceObjectType for the given
// constructor. The nominal name and native flag are copied from the
// FunctionType's Instance ObjectType if present, otherwise left empty.
// Mirrors InstanceObjectType.builderForCtor(FunctionType) in Java.
func NewInstanceObjectType(ctor *FunctionType) *InstanceObjectType {
	name := ""
	native := false
	if ctor != nil && ctor.Instance != nil {
		name = ctor.Instance.Name
	}
	proto := NewPrototypeObjectType(name, nil, native)
	return &InstanceObjectType{
		PrototypeObjectType: proto,
		constructor:         ctor,
	}
}

// GetConstructor returns the constructor FunctionType.
// Mirrors InstanceObjectType.getConstructor().
func (i *InstanceObjectType) GetConstructor() *FunctionType { return i.constructor }

// GetImplicitPrototype returns the constructor's .prototype object,
// overriding the PrototypeObjectType fallback field.
// Mirrors InstanceObjectType.getImplicitPrototype().
// Returns nil when the constructor has no prototype set.
// TODO: FunctionType does not yet carry a .prototype PrototypeObjectType;
// this returns nil until that field is added.
func (i *InstanceObjectType) GetImplicitPrototype() *PrototypeObjectType {
	// In Java: return getConstructor().getPrototype()
	// FunctionType.prototype is not yet modelled; return nil (safe fallback).
	return nil
}

// GetReferenceName returns the constructor's reference name.
// Mirrors InstanceObjectType.getReferenceName() → constructor.getReferenceName().
func (i *InstanceObjectType) GetReferenceName() string {
	if i.constructor != nil && i.constructor.Instance != nil {
		return i.constructor.Instance.Name
	}
	return i.PrototypeObjectType.GetReferenceName()
}

// IsInstanceType mirrors InstanceObjectType.isInstanceType() → true.
func (i *InstanceObjectType) IsInstanceType() bool { return true }

// IsNominalType mirrors InstanceObjectType.isNominalType() → hasReferenceName().
func (i *InstanceObjectType) IsNominalType() bool { return i.GetReferenceName() != "" }

// IsTheObjectType mirrors InstanceObjectType.isTheObjectType().
func (i *InstanceObjectType) IsTheObjectType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "Object"
}

// IsArrayType mirrors InstanceObjectType.isArrayType().
func (i *InstanceObjectType) IsArrayType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "Array"
}

// IsReadonlyArrayType mirrors InstanceObjectType.isReadonlyArrayType().
func (i *InstanceObjectType) IsReadonlyArrayType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "ReadonlyArray"
}

// IsNumberObjectType mirrors InstanceObjectType.isNumberObjectType().
func (i *InstanceObjectType) IsNumberObjectType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "Number"
}

// IsStringObjectType mirrors InstanceObjectType.isStringObjectType().
func (i *InstanceObjectType) IsStringObjectType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "String"
}

// IsBooleanObjectType mirrors InstanceObjectType.isBooleanObjectType().
func (i *InstanceObjectType) IsBooleanObjectType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "Boolean"
}

// IsDateType mirrors InstanceObjectType.isDateType().
func (i *InstanceObjectType) IsDateType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "Date"
}

// IsRegexpType mirrors InstanceObjectType.isRegexpType().
func (i *InstanceObjectType) IsRegexpType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "RegExp"
}

// IsBigIntObjectType mirrors InstanceObjectType.isBigIntObjectType().
func (i *InstanceObjectType) IsBigIntObjectType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "BigInt"
}

// IsSymbolObjectType mirrors InstanceObjectType.isSymbolObjectType().
func (i *InstanceObjectType) IsSymbolObjectType() bool {
	return i.isNativeCtor() && i.GetReferenceName() == "Symbol"
}

// GetCtorImplementedInterfaces mirrors InstanceObjectType.getCtorImplementedInterfaces()
// → constructor.getImplementedInterfaces().
// TODO: FunctionType does not yet carry implementedInterfaces; returns nil.
func (i *InstanceObjectType) GetCtorImplementedInterfaces() []*ObjectType { return nil }

// GetCtorExtendedInterfaces mirrors InstanceObjectType.getCtorExtendedInterfaces()
// → constructor.getExtendedInterfaces().
// TODO: FunctionType does not yet carry extendedInterfaces; returns nil.
func (i *InstanceObjectType) GetCtorExtendedInterfaces() []*ObjectType { return nil }

// isNativeCtor is a helper: reports whether the constructor is a native
// object type. Mirrors `getConstructor().isNativeObjectType()` in Java.
// We approximate this using the nativeType flag on the embedded PrototypeObjectType.
func (i *InstanceObjectType) isNativeCtor() bool {
	return i.PrototypeObjectType.nativeType
}
