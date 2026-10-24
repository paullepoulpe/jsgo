// Port of com/google/javascript/rhino/jstype/Visitor.java @ closure-compiler v20261021

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

package rhino

// Visitor is a type visitor.
//
// Go-only: language-adapter — Java's Visitor<T> is generic and
// JSType.visit is a generic method; Go methods cannot have type parameters,
// so T is erased to `any` (as the JVM does). Implementations return their T
// boxed in `any`, and callers convert the result back (see
// JSType_VisitForJSType for the common Visitor<JSType> case).
type Visitor interface {
	// CaseNoType is the bottom type's case.
	CaseNoType(type_ *NoType) any

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(type_ *EnumElementType) any

	// CaseAllType is the all type's case.
	CaseAllType() any

	// CaseBooleanType is the boolean value type's case.
	CaseBooleanType() any

	// CaseNoObjectType is the bottom Object type's case.
	CaseNoObjectType() any

	// CaseFunctionType is the function type's case.
	CaseFunctionType(type_ FunctionType) any

	// CaseObjectType is the object type's case.
	CaseObjectType(type_ ObjectType) any

	// CaseUnknownType is the unknown type's case.
	CaseUnknownType() any

	// CaseNullType is the null type's case.
	CaseNullType() any

	// CaseNamedType is the named type's case.
	CaseNamedType(type_ *NamedType) any

	// CaseProxyObjectType is the proxy type's case.
	CaseProxyObjectType(type_ ProxyObjectType) any

	// CaseNumberType is the number value type's case.
	CaseNumberType() any

	// CaseBigIntType is the bigint value type's case.
	CaseBigIntType() any

	// CaseStringType is the string value type's case.
	CaseStringType() any

	// CaseSymbolType is the symbol value type's case.
	CaseSymbolType(type_ SymbolType) any

	// CaseVoidType is the void type's case.
	CaseVoidType() any

	// CaseUnionType is the union type's case.
	CaseUnionType(type_ *UnionType) any

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(type_ *TemplatizedType) any

	// CaseTemplateType is the template type's case.
	CaseTemplateType(templateType *TemplateType) any
}

// Visitor_WithDefaultCase is a Visitor whose cases all delegate to
// caseDefault, which subclasses implement (Java: an abstract class nested in
// the Visitor interface).
type Visitor_WithDefaultCase interface {
	Visitor

	// CaseDefault is the case all the others delegate to.
	CaseDefault(type_ JSType) any
}

// Visitor_WithDefaultCaseBase holds the default case methods; self is the
// concrete visitor (PORTING.md §4).
type Visitor_WithDefaultCaseBase struct {
	self Visitor_WithDefaultCase
}

// InitVisitor_WithDefaultCaseBase is the Go form of the constructor.
//
// Java: Visitor.WithDefaultCase()
func (v *Visitor_WithDefaultCaseBase) InitVisitor_WithDefaultCaseBase(self Visitor_WithDefaultCase) {
	v.self = self
}

// Java: Visitor.WithDefaultCase.caseNoType(NoType)
func (v *Visitor_WithDefaultCaseBase) CaseNoType(type_ *NoType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseEnumElementType(EnumElementType)
func (v *Visitor_WithDefaultCaseBase) CaseEnumElementType(type_ *EnumElementType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseAllType()
func (v *Visitor_WithDefaultCaseBase) CaseAllType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseBooleanType()
func (v *Visitor_WithDefaultCaseBase) CaseBooleanType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseNoObjectType()
func (v *Visitor_WithDefaultCaseBase) CaseNoObjectType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseFunctionType(FunctionType)
func (v *Visitor_WithDefaultCaseBase) CaseFunctionType(type_ FunctionType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseObjectType(ObjectType)
func (v *Visitor_WithDefaultCaseBase) CaseObjectType(type_ ObjectType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseUnknownType()
func (v *Visitor_WithDefaultCaseBase) CaseUnknownType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseNullType()
func (v *Visitor_WithDefaultCaseBase) CaseNullType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseNamedType(NamedType)
func (v *Visitor_WithDefaultCaseBase) CaseNamedType(type_ *NamedType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseProxyObjectType(ProxyObjectType)
func (v *Visitor_WithDefaultCaseBase) CaseProxyObjectType(type_ ProxyObjectType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseNumberType()
func (v *Visitor_WithDefaultCaseBase) CaseNumberType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseBigIntType()
func (v *Visitor_WithDefaultCaseBase) CaseBigIntType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseStringType()
func (v *Visitor_WithDefaultCaseBase) CaseStringType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseSymbolType(SymbolType)
func (v *Visitor_WithDefaultCaseBase) CaseSymbolType(type_ SymbolType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseVoidType()
func (v *Visitor_WithDefaultCaseBase) CaseVoidType() any {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseUnionType(UnionType)
func (v *Visitor_WithDefaultCaseBase) CaseUnionType(type_ *UnionType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseTemplatizedType(TemplatizedType)
func (v *Visitor_WithDefaultCaseBase) CaseTemplatizedType(type_ *TemplatizedType) any {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseTemplateType(TemplateType)
func (v *Visitor_WithDefaultCaseBase) CaseTemplateType(type_ *TemplateType) any {
	return v.self.CaseDefault(type_)
}
