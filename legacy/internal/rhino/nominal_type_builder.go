// Port of com/google/javascript/rhino/NominalTypeBuilder.java @ closure-compiler v20260415
//
// NominalTypeBuilder depends on FunctionType and ObjectType from
// com.google.javascript.rhino.jstype. Because the Go jstype package imports
// the rhino package (for Node references), importing jstype here would create
// an import cycle. The constructor and instance fields are therefore typed as
// any; callers in the jscomp layer perform type assertions to the concrete
// jstype types at the use-site.

package rhino

// NominalTypeBuilder is a builder for declaring properties on class-like
// (nominal) types. Nominal types consist primarily of three separate object
// types which may each have their own properties declared on them separately:
// (1) the constructor function, (2) the instance type, and (3) the prototype
// object.
//
// This builder is used during the first part of type checking, while the type
// checker builds up the set of known types. It serves as a placeholder to
// allow operating on (assigning properties to, and referencing the
// yet-to-exist type) these not-yet-available types.
//
// Mirrors com.google.javascript.rhino.NominalTypeBuilder.
type NominalTypeBuilder struct {
	// constructor holds a *jstype.FunctionType (typed as any to avoid import cycle).
	constructor any
	// instance holds a *jstype.ObjectType (typed as any to avoid import cycle).
	instance any
	// prototype holds a *jstype.ObjectType (typed as any to avoid import cycle).
	prototype any
}

// NominalTypePrototypeGetter is a function that retrieves the prototype
// ObjectType from a FunctionType. Callers in jscomp inject this to break the
// import cycle.
var NominalTypePrototypeGetter func(constructor any) any

// NewNominalTypeBuilder creates a NominalTypeBuilder for the given constructor
// and instance types. The prototype is derived from the constructor via
// NominalTypePrototypeGetter.
// Mirrors NominalTypeBuilder(FunctionType, ObjectType).
func NewNominalTypeBuilder(constructor, instance any) *NominalTypeBuilder {
	var proto any
	if NominalTypePrototypeGetter != nil {
		proto = NominalTypePrototypeGetter(constructor)
	}
	return &NominalTypeBuilder{
		constructor: constructor,
		instance:    instance,
		prototype:   proto,
	}
}

// RawConstructor returns the raw constructor value (a *jstype.FunctionType).
func (b *NominalTypeBuilder) RawConstructor() any { return b.constructor }

// RawInstance returns the raw instance value (a *jstype.ObjectType).
func (b *NominalTypeBuilder) RawInstance() any { return b.instance }

// RawPrototype returns the raw prototype value (a *jstype.ObjectType).
func (b *NominalTypeBuilder) RawPrototype() any { return b.prototype }

// NominalTypeSuperClassGetter is a function that retrieves the superclass
// constructor from an instance ObjectType. Callers in jscomp inject this.
var NominalTypeSuperClassGetter func(instance any) any

// NominalTypeInstanceGetter is a function that retrieves the instance type
// from a FunctionType constructor. Callers in jscomp inject this.
var NominalTypeInstanceGetter func(constructor any) any

// SuperClass returns a NominalTypeBuilder for this type's superclass, or nil
// if there is no superclass. Mirrors NominalTypeBuilder.superClass().
func (b *NominalTypeBuilder) SuperClass() *NominalTypeBuilder {
	if NominalTypeSuperClassGetter == nil || NominalTypeInstanceGetter == nil {
		return nil
	}
	ctor := NominalTypeSuperClassGetter(b.instance)
	if ctor == nil {
		return nil
	}
	inst := NominalTypeInstanceGetter(ctor)
	return NewNominalTypeBuilder(ctor, inst)
}
