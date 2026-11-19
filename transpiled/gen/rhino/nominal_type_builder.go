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
// Port of com/google/javascript/rhino/NominalTypeBuilder.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.NominalTypeBuilder
type NominalTypeBuilder struct {
	constructor *FunctionType `gson:"constructor"`
	instance ObjectType `gson:"instance"`
	prototype ObjectType `gson:"prototype"`
	idHash_ atomic.Int32
}

func (this_ *NominalTypeBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *NominalTypeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*NominalTypeBuilder) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.NominalTypeBuilder", "NominalTypeBuilder" }
func (this_ *NominalTypeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NominalTypeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.NominalTypeBuilder", "com.google.javascript.rhino.NominalTypeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*NominalTypeBuilder](v_)
	})
}

func (this_ *NominalTypeBuilder) initFields_NominalTypeBuilder() {
	_ = this_
}

// Java constructor: NominalTypeBuilder(com.google.javascript.rhino.jstype.FunctionType constructor, com.google.javascript.rhino.jstype.ObjectType instance)
func NewNominalTypeBuilder(constructor *FunctionType, instance ObjectType) *NominalTypeBuilder {
	this_ := &NominalTypeBuilder{}
	InitNominalTypeBuilder(this_, constructor, instance)
	return this_
}

func InitNominalTypeBuilder(this_ *NominalTypeBuilder, constructor *FunctionType, instance ObjectType) {
	_ = constructor
	_ = instance
	this_.initFields_NominalTypeBuilder()
	this_.constructor = constructor
	this_.instance = instance
	this_.prototype = constructor.GetPrototypeProperty()
}

// Java method: NominalTypeBuilder.declarePrototypeProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node defSite)
func (this_ *NominalTypeBuilder) DeclarePrototypeProperty(name string, type_ JSType, defSite Node) {
	_ = name
	_ = type_
	_ = defSite
	this_.prototype.DefineDeclaredProperty_3_string_JSType_Node(name, type_, defSite)
}

// Java method: NominalTypeBuilder.declareInstanceProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node defSite)
func (this_ *NominalTypeBuilder) DeclareInstanceProperty(name string, type_ JSType, defSite Node) {
	_ = name
	_ = type_
	_ = defSite
	this_.instance.DefineDeclaredProperty_3_string_JSType_Node(name, type_, defSite)
}

// Java method: NominalTypeBuilder.declareConstructorProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node defSite)
func (this_ *NominalTypeBuilder) DeclareConstructorProperty(name string, type_ JSType, defSite Node) {
	_ = name
	_ = type_
	_ = defSite
	this_.constructor.DefineDeclaredProperty_3_string_JSType_Node(name, type_, defSite)
}

// Java method: NominalTypeBuilder.superClass()
func (this_ *NominalTypeBuilder) SuperClass() *NominalTypeBuilder {
	var ctor *FunctionType = this_.instance.GetSuperClassConstructor()
	_ = ctor
	if ((ctor == nil)) {
		return nil
	}
	return NewNominalTypeBuilder(ctor, ctor.GetInstanceType())
}

// Java method: NominalTypeBuilder.constructor()
func (this_ *NominalTypeBuilder) Constructor() *FunctionType {
	return this_.constructor
}

// Java method: NominalTypeBuilder.instance()
func (this_ *NominalTypeBuilder) Instance() ObjectType {
	return this_.instance
}

// Java method: NominalTypeBuilder.prototypeOrInstance()
func (this_ *NominalTypeBuilder) PrototypeOrInstance() ObjectType {
	return this_.prototype
}

