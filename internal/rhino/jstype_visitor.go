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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Visitor is a generic JSType visitor.
//
// Java: interface Visitor<T>; dispatched via VisitJSType since Go methods
// cannot declare their own type parameters.
type Visitor[T any] interface {
	// CaseNoType is the bottom type's case.
	CaseNoType(type_ *NoType) T

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(type_ *EnumElementType) T

	// CaseAllType is the all type's case.
	CaseAllType() T

	// CaseBooleanType is the boolean value type's case.
	CaseBooleanType() T

	// CaseNoObjectType is the bottom Object type's case.
	CaseNoObjectType() T

	// CaseFunctionType is the function type's case.
	CaseFunctionType(type_ FunctionType) T

	// CaseObjectType is the object type's case.
	CaseObjectType(type_ ObjectType) T

	// CaseUnknownType is the unknown type's case.
	CaseUnknownType() T

	// CaseNullType is the null type's case.
	CaseNullType() T

	// CaseNamedType is the named type's case.
	CaseNamedType(type_ *NamedType) T

	// CaseProxyObjectType is the proxy type's case.
	CaseProxyObjectType(type_ ProxyObjectType) T

	// CaseNumberType is the number value type's case.
	CaseNumberType() T

	// CaseBigIntType is the bigint value type's case.
	CaseBigIntType() T

	// CaseStringType is the string value type's case.
	CaseStringType() T

	// CaseSymbolType is the symbol value type's case.
	CaseSymbolType(type_ SymbolType) T

	// CaseVoidType is the void type's case.
	CaseVoidType() T

	// CaseUnionType is the union type's case.
	CaseUnionType(type_ *UnionType) T

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(type_ *TemplatizedType) T

	// CaseTemplateType is the template type's case.
	CaseTemplateType(templateType *TemplateType) T
}

// VisitJSType dispatches t to visitor's matching case method.
//
// Java: JSType.visit(Visitor<T>)
func VisitJSType[T any](t JSType, visitor Visitor[T]) T {
	switch x := t.(type) {
	case *NoType:
		return visitor.CaseNoType(x)
	case *EnumElementType:
		return visitor.CaseEnumElementType(x)
	case *AllType:
		return visitor.CaseAllType()
	case *BooleanType:
		return visitor.CaseBooleanType()
	case *NoObjectTypeBase:
		return visitor.CaseNoObjectType()
	case *FunctionTypeBase:
		return visitor.CaseFunctionType(x)
	case *InstanceObjectType:
		return visitor.CaseObjectType(x)
	case *PrototypeObjectTypeBase:
		return visitor.CaseObjectType(x)
	case *RecordType:
		return visitor.CaseObjectType(x)
	case *EnumType:
		return visitor.CaseObjectType(x)
	case *NoResolvedType:
		return visitor.CaseObjectType(x)
	case *UnknownType:
		return visitor.CaseUnknownType()
	case *NullType:
		return visitor.CaseNullType()
	case *NamedType:
		return visitor.CaseNamedType(x)
	case *ProxyObjectTypeBase:
		return visitor.CaseProxyObjectType(x)
	case *NumberType:
		return visitor.CaseNumberType()
	case *BigIntType:
		return visitor.CaseBigIntType()
	case *StringType:
		return visitor.CaseStringType()
	case *SymbolTypeBase:
		return visitor.CaseSymbolType(x)
	case *KnownSymbolType:
		return visitor.CaseSymbolType(x)
	case *VoidType:
		return visitor.CaseVoidType()
	case *UnionType:
		return visitor.CaseUnionType(x)
	case *TemplatizedType:
		return visitor.CaseTemplatizedType(x)
	case *TemplateType:
		return visitor.CaseTemplateType(x)
	default:
		if _, ok := t.(NoObjectType); ok {
			return visitor.CaseNoObjectType()
		}
		if fn, ok := t.(FunctionType); ok {
			return visitor.CaseFunctionType(fn)
		}
		if proxy, ok := t.(ProxyObjectType); ok {
			return visitor.CaseProxyObjectType(proxy)
		}
		if obj, ok := t.(ObjectType); ok {
			return visitor.CaseObjectType(obj)
		}
		if sym, ok := t.(SymbolType); ok {
			return visitor.CaseSymbolType(sym)
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Visitor_WithDefaultCase is a Visitor whose cases all delegate to
// caseDefault, which subclasses implement (Java: an abstract class nested in
// the Visitor interface).
type Visitor_WithDefaultCase[T any] interface {
	Visitor[T]

	// CaseDefault is the case all the others delegate to.
	CaseDefault(type_ JSType) T
}

// Visitor_WithDefaultCaseBase holds the default case methods; self is the
// concrete visitor (PORTING.md §4).
type Visitor_WithDefaultCaseBase[T any] struct {
	self Visitor_WithDefaultCase[T]
}

// InitVisitor_WithDefaultCaseBase is the Go form of the constructor.
//
// Java: Visitor.WithDefaultCase()
func (v *Visitor_WithDefaultCaseBase[T]) InitVisitor_WithDefaultCaseBase(self Visitor_WithDefaultCase[T]) {
	v.self = self
}

// Java: Visitor.WithDefaultCase.caseNoType(NoType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseNoType(type_ *NoType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseEnumElementType(EnumElementType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseEnumElementType(type_ *EnumElementType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseAllType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseAllType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseBooleanType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseBooleanType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseNoObjectType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseNoObjectType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseFunctionType(FunctionType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseFunctionType(type_ FunctionType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseObjectType(ObjectType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseObjectType(type_ ObjectType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseUnknownType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseUnknownType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseNullType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseNullType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseNamedType(NamedType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseNamedType(type_ *NamedType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseProxyObjectType(ProxyObjectType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseProxyObjectType(type_ ProxyObjectType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseNumberType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseNumberType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseBigIntType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseBigIntType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseStringType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseStringType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseSymbolType(SymbolType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseSymbolType(type_ SymbolType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseVoidType()
func (v *Visitor_WithDefaultCaseBase[T]) CaseVoidType() T {
	return v.self.CaseDefault(nil)
}

// Java: Visitor.WithDefaultCase.caseUnionType(UnionType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseUnionType(type_ *UnionType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseTemplatizedType(TemplatizedType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseTemplatizedType(type_ *TemplatizedType) T {
	return v.self.CaseDefault(type_)
}

// Java: Visitor.WithDefaultCase.caseTemplateType(TemplateType)
func (v *Visitor_WithDefaultCaseBase[T]) CaseTemplateType(type_ *TemplateType) T {
	return v.self.CaseDefault(type_)
}
