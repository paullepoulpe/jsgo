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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/Visitor.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.rhino.jstype.Visitor
type Visitor interface {
	Implements_rhino_Visitor()
	IsJavaNil_() bool
	CaseNoType(*FunctionType) any
	CaseEnumElementType(*EnumElementType) any
	CaseAllType() any
	CaseBooleanType() any
	CaseNoObjectType() any
	CaseFunctionType(*FunctionType) any
	CaseObjectType(ObjectType) any
	CaseUnknownType() any
	CaseNullType() any
	CaseNamedType(*NamedType) any
	CaseProxyObjectType(ProxyObjectType) any
	CaseNumberType() any
	CaseBigIntType() any
	CaseStringType() any
	CaseSymbolType(SymbolType) any
	CaseVoidType() any
	CaseUnionType(*UnionType) any
	CaseTemplatizedType(*TemplatizedType) any
	CaseTemplateType(*TemplateType) any
}
func ToIface_Visitor[T interface {
	Visitor
	comparable
}](v T) Visitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Visitor[T Visitor](v T) Visitor {
	var vi Visitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Visitor", "com.google.javascript.rhino.jstype.Visitor", func(v_ any) bool {
		return j2gort.InstanceOf[Visitor](v_)
	})
}
// Java polymorphic class interface: com.google.javascript.rhino.jstype.Visitor.WithDefaultCase
type Visitor_WithDefaultCase interface {
	Implements_rhino_Visitor_WithDefaultCase()
	IsJavaNil_() bool
	GetVisitor_WithDefaultCaseBase() *Visitor_WithDefaultCaseBase
	Visitor
	CaseDefault(JSType) any
}
func ToIface_Visitor_WithDefaultCase[T interface {
	Visitor_WithDefaultCase
	comparable
}](v T) Visitor_WithDefaultCase {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Visitor_WithDefaultCase[T Visitor_WithDefaultCase](v T) Visitor_WithDefaultCase {
	var vi Visitor_WithDefaultCase = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.Visitor.WithDefaultCase
type Visitor_WithDefaultCaseBase struct {
	self_ Visitor_WithDefaultCase
	idHash_ atomic.Int32
}

func (this_ *Visitor_WithDefaultCaseBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Visitor_WithDefaultCase)
}
func (this_ *Visitor_WithDefaultCaseBase) IsJavaNil_() bool { return this_ == nil }
func (*Visitor_WithDefaultCaseBase) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.Visitor$WithDefaultCase", "WithDefaultCase" }
func (this_ *Visitor_WithDefaultCaseBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Visitor.WithDefaultCase", "com.google.javascript.rhino.jstype.Visitor$WithDefaultCase", func(v_ any) bool {
		return j2gort.InstanceOf[Visitor_WithDefaultCase](v_)
	})
}
func (b *Visitor_WithDefaultCaseBase) GetVisitor_WithDefaultCaseBase() *Visitor_WithDefaultCaseBase { return b }
func (*Visitor_WithDefaultCaseBase) Implements_rhino_Visitor_WithDefaultCase() {}
func (*Visitor_WithDefaultCaseBase) Implements_rhino_Visitor() {}

func (this_ *Visitor_WithDefaultCaseBase) initFields_Visitor_WithDefaultCase() {
	_ = this_
}

func InitVisitor_WithDefaultCase(this_ *Visitor_WithDefaultCaseBase) {
	this_.initFields_Visitor_WithDefaultCase()
}

// Java method: WithDefaultCase.caseDefault(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseDefault(type_ JSType) any {
	_ = type_
	panic("abstract or native method invoked")
}

// Java method: WithDefaultCase.caseNoType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseNoType(type_ *FunctionType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

// Java method: WithDefaultCase.caseEnumElementType(com.google.javascript.rhino.jstype.EnumElementType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseEnumElementType(type_ *EnumElementType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

// Java method: WithDefaultCase.caseAllType()
func (this_ *Visitor_WithDefaultCaseBase) CaseAllType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseBooleanType()
func (this_ *Visitor_WithDefaultCaseBase) CaseBooleanType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseNoObjectType()
func (this_ *Visitor_WithDefaultCaseBase) CaseNoObjectType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseFunctionType(type_ *FunctionType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

// Java method: WithDefaultCase.caseObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseObjectType(type_ ObjectType) any {
	_ = type_
	return this_.self_.CaseDefault(Upcast_JSType(type_))
}

// Java method: WithDefaultCase.caseUnknownType()
func (this_ *Visitor_WithDefaultCaseBase) CaseUnknownType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseNullType()
func (this_ *Visitor_WithDefaultCaseBase) CaseNullType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseNamedType(com.google.javascript.rhino.jstype.NamedType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseNamedType(type_ *NamedType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

// Java method: WithDefaultCase.caseProxyObjectType(com.google.javascript.rhino.jstype.ProxyObjectType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseProxyObjectType(type_ ProxyObjectType) any {
	_ = type_
	return this_.self_.CaseDefault(Upcast_JSType(type_))
}

// Java method: WithDefaultCase.caseNumberType()
func (this_ *Visitor_WithDefaultCaseBase) CaseNumberType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseBigIntType()
func (this_ *Visitor_WithDefaultCaseBase) CaseBigIntType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseStringType()
func (this_ *Visitor_WithDefaultCaseBase) CaseStringType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseSymbolType(com.google.javascript.rhino.jstype.SymbolType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseSymbolType(type_ SymbolType) any {
	_ = type_
	return this_.self_.CaseDefault(Upcast_JSType(type_))
}

// Java method: WithDefaultCase.caseVoidType()
func (this_ *Visitor_WithDefaultCaseBase) CaseVoidType() any {
	return this_.self_.CaseDefault(nil)
}

// Java method: WithDefaultCase.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseUnionType(type_ *UnionType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

// Java method: WithDefaultCase.caseTemplatizedType(com.google.javascript.rhino.jstype.TemplatizedType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseTemplatizedType(type_ *TemplatizedType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

// Java method: WithDefaultCase.caseTemplateType(com.google.javascript.rhino.jstype.TemplateType type)
func (this_ *Visitor_WithDefaultCaseBase) CaseTemplateType(type_ *TemplateType) any {
	_ = type_
	return this_.self_.CaseDefault(ToIface_JSType(type_))
}

