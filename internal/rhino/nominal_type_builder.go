// Port of com/google/javascript/rhino/NominalTypeBuilder.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Ben Lickly
 *   Dimitris Vardoulakis
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
