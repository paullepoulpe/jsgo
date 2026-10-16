// Port of com/google/javascript/rhino/NominalTypeBuilder.java @ closure-compiler v20260415

package rhino

// NominalTypeBuilder is the builder interface for declaring properties on
// class-like (nominal) types. Nominal types consist primarily of three
// separate object types, which may each have their own properties declared
// on them separately: (1) the constructor function, (2) the instance type,
// and (3) the prototype object.
//
// This builder is used during the first part of type checking, while the
// type checker builds up the set of known types. Thus, this interface serves
// as a placeholder to allow operating on (i.e. both assigning properties to,
// and referencing the yet-to-exist type) these not-yet-available types.
type NominalTypeBuilder struct {
	constructor FunctionType
	instance    ObjectType
	prototype   ObjectType
}

// Java: NominalTypeBuilder(FunctionType, ObjectType)
func NewNominalTypeBuilder(constructor FunctionType, instance ObjectType) *NominalTypeBuilder {
	return &NominalTypeBuilder{
		constructor: constructor,
		instance:    instance,
		prototype:   constructor.GetPrototypeProperty(),
	}
}

// DeclarePrototypeProperty declares a property on the nominal type's
// prototype.
//
// Java: NominalTypeBuilder.declarePrototypeProperty(String, JSType, Node)
func (b *NominalTypeBuilder) DeclarePrototypeProperty(name string, type_ JSType, defSite *Node) {
	b.prototype.DefineDeclaredProperty(name, type_, defSite)
}

// DeclareInstanceProperty declares an instance property on the nominal
// type.
//
// Java: NominalTypeBuilder.declareInstanceProperty(String, JSType, Node)
func (b *NominalTypeBuilder) DeclareInstanceProperty(name string, type_ JSType, defSite *Node) {
	b.instance.DefineDeclaredProperty(name, type_, defSite)
}

// DeclareConstructorProperty declares a static property on the nominal
// type's constructor.
//
// Java: NominalTypeBuilder.declareConstructorProperty(String, JSType, Node)
func (b *NominalTypeBuilder) DeclareConstructorProperty(name string, type_ JSType, defSite *Node) {
	b.constructor.DefineDeclaredProperty(name, type_, defSite)
}

// SuperClass returns a NominalTypeBuilder for this type's superclass.
//
// Java: NominalTypeBuilder.superClass()
func (b *NominalTypeBuilder) SuperClass() *NominalTypeBuilder {
	ctor := b.instance.GetSuperClassConstructor()
	if ctor == nil {
		return nil
	}
	return NewNominalTypeBuilder(ctor, ctor.GetInstanceType())
}

// Constructor returns the constructor as a JSType.
//
// Java: NominalTypeBuilder.constructor()
func (b *NominalTypeBuilder) Constructor() FunctionType {
	return b.constructor
}

// Instance returns the instance type as a JSType.
//
// Java: NominalTypeBuilder.instance()
func (b *NominalTypeBuilder) Instance() ObjectType {
	return b.instance
}

// PrototypeOrInstance returns the type of the prototype object (OTI).
//
// TODO(sdh): See if we can just delete this entirely and use instance()
// instead?
//
// Java: NominalTypeBuilder.prototypeOrInstance()
func (b *NominalTypeBuilder) PrototypeOrInstance() ObjectType {
	return b.prototype
}
