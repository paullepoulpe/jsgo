// Port of com/google/javascript/jscomp/CodingConventions.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodingConventions holds helper classes for dealing with coding
// conventions.
type CodingConventions struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

// Java: CodingConventions()
func newCodingConventions() *CodingConventions { return &CodingConventions{} }

// CodingConventions_GetDefault gets the default coding convention.
//
// Java: CodingConventions.getDefault()
func CodingConventions_GetDefault() CodingConvention {
	return newCodingConventions_DefaultCodingConvention()
}

// CodingConventions_DefaultIsFunctionCallThatAlwaysThrows: n is the last
// statement of a block to check for an always throws function call. Used by
// CheckMissingReturn. alwaysThrowsFunctionName is the name of a function that
// always throws. Returns true if n is call to alwaysThrowsFunctionName,
// otherwise false.
//
// Java: CodingConventions.defaultIsFunctionCallThatAlwaysThrows(Node, String)
func CodingConventions_DefaultIsFunctionCallThatAlwaysThrows(
	n rhino.Node, alwaysThrowsFunctionName string) bool {
	if n.IsExprResult() {
		if !n.GetFirstChild().IsCall() {
			return false
		}
	} else if !n.IsCall() {
		return false
	}
	if n.IsExprResult() {
		n = n.GetFirstChild()
	}
	// n is a call
	return n.GetFirstChild().MatchesQualifiedName(alwaysThrowsFunctionName)
}

// CodingConventions_Proxy is a convention that wraps another.
//
// When you want to support a new library, you should subclass this
// delegate, and override the methods that you want to customize.
//
// This way, a person using jQuery and Closure Library can create a new
// coding convention by creating a jQueryCodingConvention that delegates
// to a ClosureCodingConvention that delegates to a DefaultCodingConvention.
type CodingConventions_Proxy interface {
	CodingConvention

	GetCodingConventions_ProxyBase() *CodingConventions_ProxyBase
}

// CodingConventions_ProxyBase holds the fields and base implementations of
// CodingConventions.Proxy. Subclasses embed it, set self, then call
// InitCodingConventions_ProxyBase.
type CodingConventions_ProxyBase struct {
	self CodingConvention

	nextConvention CodingConvention
}

// GetCodingConventions_ProxyBase gives access to the Proxy fields through
// the interface.
func (p *CodingConventions_ProxyBase) GetCodingConventions_ProxyBase() *CodingConventions_ProxyBase {
	return p
}

// Java: CodingConventions.Proxy(CodingConvention)
func (p *CodingConventions_ProxyBase) InitCodingConventions_ProxyBase(convention CodingConvention) {
	p.nextConvention = convention
}

// Java: CodingConventions.Proxy.isConstant(String)
func (p *CodingConventions_ProxyBase) IsConstant(variableName string) bool {
	return p.nextConvention.IsConstant(variableName)
}

// Java: CodingConventions.Proxy.isConstantKey(String)
func (p *CodingConventions_ProxyBase) IsConstantKey(keyName string) bool {
	return p.nextConvention.IsConstantKey(keyName)
}

// Java: CodingConventions.Proxy.isValidEnumKey(String)
func (p *CodingConventions_ProxyBase) IsValidEnumKey(key string) bool {
	return p.nextConvention.IsValidEnumKey(key)
}

// Java: CodingConventions.Proxy.isOptionalParameter(Node)
func (p *CodingConventions_ProxyBase) IsOptionalParameter(parameter rhino.Node) bool {
	return p.nextConvention.IsOptionalParameter(parameter)
}

// Java: CodingConventions.Proxy.isVarArgsParameter(Node)
func (p *CodingConventions_ProxyBase) IsVarArgsParameter(parameter rhino.Node) bool {
	return p.nextConvention.IsVarArgsParameter(parameter)
}

// Java: CodingConventions.Proxy.isFunctionCallThatAlwaysThrows(Node)
func (p *CodingConventions_ProxyBase) IsFunctionCallThatAlwaysThrows(n rhino.Node) bool {
	return p.nextConvention.IsFunctionCallThatAlwaysThrows(n)
}

// Java: CodingConventions.Proxy.isExported(String, boolean)
func (p *CodingConventions_ProxyBase) IsExportedWithLocal(name string, local bool) bool {
	return p.nextConvention.IsExportedWithLocal(name, local)
}

// IsExported is final in Java: it runs CodingConvention's default body on
// this convention (through self, so that subclass overrides of
// isExported(String, boolean) are seen).
//
// Java: CodingConventions.Proxy.isExported(String)
func (p *CodingConventions_ProxyBase) IsExported(name string) bool {
	return CodingConvention_IsExported(p.self, name)
}

// BlockRenamingForProperty is CodingConvention's default method, which
// Proxy does not override.
//
// Java: CodingConvention.blockRenamingForProperty(String) (default method)
func (p *CodingConventions_ProxyBase) BlockRenamingForProperty(name string) bool {
	return CodingConvention_BlockRenamingForProperty(p.self, name)
}

// Java: CodingConventions.Proxy.getPackageName(StaticSourceFile)
func (p *CodingConventions_ProxyBase) GetPackageName(source rhino.StaticSourceFile) *string {
	return p.nextConvention.GetPackageName(source)
}

// Java: CodingConventions.Proxy.getClassesDefinedByCall(Node)
func (p *CodingConventions_ProxyBase) GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship {
	return p.nextConvention.GetClassesDefinedByCall(callNode)
}

// Java: CodingConventions.Proxy.isClassFactoryCall(Node)
func (p *CodingConventions_ProxyBase) IsClassFactoryCall(callNode rhino.Node) bool {
	return p.nextConvention.IsClassFactoryCall(callNode)
}

// Java: CodingConventions.Proxy.isSuperClassReference(String)
func (p *CodingConventions_ProxyBase) IsSuperClassReference(propertyName string) bool {
	return p.nextConvention.IsSuperClassReference(propertyName)
}

// Java: CodingConventions.Proxy.extractIsModuleFile(Node, Node)
func (p *CodingConventions_ProxyBase) ExtractIsModuleFile(node rhino.Node, parent rhino.Node) bool {
	return p.nextConvention.ExtractIsModuleFile(node, parent)
}

// Java: CodingConventions.Proxy.extractClassNameIfProvide(Node, Node)
func (p *CodingConventions_ProxyBase) ExtractClassNameIfProvide(node rhino.Node, parent rhino.Node) *string {
	return p.nextConvention.ExtractClassNameIfProvide(node, parent)
}

// Java: CodingConventions.Proxy.extractClassNameIfRequire(Node, Node)
func (p *CodingConventions_ProxyBase) ExtractClassNameIfRequire(node rhino.Node, parent rhino.Node) *string {
	return p.nextConvention.ExtractClassNameIfRequire(node, parent)
}

// Java: CodingConventions.Proxy.getExportPropertyFunction()
func (p *CodingConventions_ProxyBase) GetExportPropertyFunction() *string {
	return p.nextConvention.GetExportPropertyFunction()
}

// Java: CodingConventions.Proxy.getExportSymbolFunction()
func (p *CodingConventions_ProxyBase) GetExportSymbolFunction() *string {
	return p.nextConvention.GetExportSymbolFunction()
}

// Java: CodingConventions.Proxy.identifyTypeDeclarationCall(Node)
func (p *CodingConventions_ProxyBase) IdentifyTypeDeclarationCall(n rhino.Node) []string {
	return p.nextConvention.IdentifyTypeDeclarationCall(n)
}

// Java: CodingConventions.Proxy.applySubclassRelationship(NominalTypeBuilder, NominalTypeBuilder, SubclassType)
func (p *CodingConventions_ProxyBase) ApplySubclassRelationship(
	parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	p.nextConvention.ApplySubclassRelationship(parent, child, type_)
}

// Java: CodingConventions.Proxy.getAbstractMethodName()
func (p *CodingConventions_ProxyBase) GetAbstractMethodName() *string {
	return p.nextConvention.GetAbstractMethodName()
}

// Java: CodingConventions.Proxy.getSingletonGetterClassName(Node)
func (p *CodingConventions_ProxyBase) GetSingletonGetterClassName(callNode rhino.Node) *string {
	return p.nextConvention.GetSingletonGetterClassName(callNode)
}

// Java: CodingConventions.Proxy.applySingletonGetter(NominalTypeBuilder, FunctionType)
func (p *CodingConventions_ProxyBase) ApplySingletonGetter(
	classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	p.nextConvention.ApplySingletonGetter(classType, getterType)
}

// Java: CodingConventions.Proxy.getAssertionFunctions()
func (p *CodingConventions_ProxyBase) GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec {
	return p.nextConvention.GetAssertionFunctions()
}

// Java: CodingConventions.Proxy.describeFunctionBind(Node, boolean)
func (p *CodingConventions_ProxyBase) DescribeFunctionBind(n rhino.Node, checkTypes bool) *CodingConvention_Bind {
	return p.nextConvention.DescribeFunctionBind(n, checkTypes)
}

// Java: CodingConventions.Proxy.describeCachingCall(Node)
func (p *CodingConventions_ProxyBase) DescribeCachingCall(node rhino.Node) *CodingConvention_Cache {
	return p.nextConvention.DescribeCachingCall(node)
}

// Java: CodingConventions.Proxy.isPropertyTestFunction(Node)
func (p *CodingConventions_ProxyBase) IsPropertyTestFunction(call rhino.Node) bool {
	return p.nextConvention.IsPropertyTestFunction(call)
}

// Java: CodingConventions.Proxy.isPropertyRenameFunction(Node)
func (p *CodingConventions_ProxyBase) IsPropertyRenameFunction(nameNode rhino.Node) bool {
	return p.nextConvention.IsPropertyRenameFunction(nameNode)
}

// Java: CodingConventions.Proxy.isPrototypeAlias(Node)
func (p *CodingConventions_ProxyBase) IsPrototypeAlias(getProp rhino.Node) bool {
	return false
}

// Java: CodingConventions.Proxy.getObjectLiteralCast(Node)
func (p *CodingConventions_ProxyBase) GetObjectLiteralCast(callNode rhino.Node) *CodingConvention_ObjectLiteralCast {
	return p.nextConvention.GetObjectLiteralCast(callNode)
}

// CodingConventions_DefaultCodingConvention is the default coding
// convention. Should be at the bottom of all proxy chains.
type CodingConventions_DefaultCodingConvention struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

const codingConventions_DefaultCodingConvention_serialVersionUID int64 = 1

// Java: CodingConventions.DefaultCodingConvention()
func newCodingConventions_DefaultCodingConvention() *CodingConventions_DefaultCodingConvention {
	return &CodingConventions_DefaultCodingConvention{}
}

// Java: CodingConventions.DefaultCodingConvention.isConstant(String)
func (d *CodingConventions_DefaultCodingConvention) IsConstant(variableName string) bool {
	return false
}

// Java: CodingConventions.DefaultCodingConvention.isConstantKey(String)
func (d *CodingConventions_DefaultCodingConvention) IsConstantKey(variableName string) bool {
	return false
}

// IsValidEnumKey: key != null && key.length() > 0 (a Go string is never
// null).
//
// Java: CodingConventions.DefaultCodingConvention.isValidEnumKey(String)
func (d *CodingConventions_DefaultCodingConvention) IsValidEnumKey(key string) bool {
	return jutil.StringLength(key) > 0
}

// Java: CodingConventions.DefaultCodingConvention.isOptionalParameter(Node)
func (d *CodingConventions_DefaultCodingConvention) IsOptionalParameter(parameter rhino.Node) bool {
	return false
}

// Java: CodingConventions.DefaultCodingConvention.isVarArgsParameter(Node)
func (d *CodingConventions_DefaultCodingConvention) IsVarArgsParameter(parameter rhino.Node) bool {
	// be as lax as possible
	return parameter.IsRest()
}

// Java: CodingConventions.DefaultCodingConvention.isFunctionCallThatAlwaysThrows(Node)
func (d *CodingConventions_DefaultCodingConvention) IsFunctionCallThatAlwaysThrows(n rhino.Node) bool {
	if NodeUtil_IsExprCall(n) {
		fnType := rhino.JSType_ToMaybeFunctionType(n.GetFirstFirstChild().GetJSType())
		return fnType != nil && fnType.GetClosurePrimitive() != nil &&
			rhino.ClosurePrimitive_ASSERTS_FAIL == *fnType.GetClosurePrimitive()
	}
	return false
}

// Java: CodingConventions.DefaultCodingConvention.getPackageName(StaticSourceFile)
func (d *CodingConventions_DefaultCodingConvention) GetPackageName(source rhino.StaticSourceFile) *string {
	// The package name of a source file is its file path.
	name := source.GetName()
	lastSlash := jutil.StringLastIndexOf(name, "/")
	var result string
	if lastSlash == -1 {
		result = ""
	} else {
		result = jutil.StringSubstring(name, 0, lastSlash)
	}
	return &result
}

// Java: CodingConventions.DefaultCodingConvention.isExported(String, boolean)
func (d *CodingConventions_DefaultCodingConvention) IsExportedWithLocal(name string, local bool) bool {
	return local && jutil.StringStartsWith(name, "$super")
}

// Java: CodingConventions.DefaultCodingConvention.isExported(String)
func (d *CodingConventions_DefaultCodingConvention) IsExported(name string) bool {
	return CodingConvention_IsExported(d, name)
}

// BlockRenamingForProperty is CodingConvention's default method, which
// DefaultCodingConvention does not override.
//
// Java: CodingConvention.blockRenamingForProperty(String) (default method)
func (d *CodingConventions_DefaultCodingConvention) BlockRenamingForProperty(name string) bool {
	return CodingConvention_BlockRenamingForProperty(d, name)
}

var codingConventions_DefaultCodingConvention_JSCOMP_INHERITS = rhino.QualifiedName_Of("$jscomp.inherits")

// Java: CodingConventions.DefaultCodingConvention.getClassesDefinedByCall(Node)
func (d *CodingConventions_DefaultCodingConvention) GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship {
	callName := callNode.GetFirstChild()
	if (codingConventions_DefaultCodingConvention_JSCOMP_INHERITS.Matches(callName) || callName.MatchesName("$jscomp$inherits")) &&
		callNode.HasXChildren(3) {
		subclass := callName.GetNext()
		superclass := subclass.GetNext()
		// The StripCode pass may create $jscomp.inherits calls with NULL arguments.
		if subclass.IsQualifiedName() && superclass.IsQualifiedName() {
			return NewCodingConvention_SubclassRelationship(CodingConvention_SubclassType_INHERITS, subclass, superclass)
		}
	}
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.isClassFactoryCall(Node)
func (d *CodingConventions_DefaultCodingConvention) IsClassFactoryCall(callNode rhino.Node) bool {
	return false
}

// Java: CodingConventions.DefaultCodingConvention.isSuperClassReference(String)
func (d *CodingConventions_DefaultCodingConvention) IsSuperClassReference(propertyName string) bool {
	return false
}

// Java: CodingConventions.DefaultCodingConvention.extractIsModuleFile(Node, Node)
func (d *CodingConventions_DefaultCodingConvention) ExtractIsModuleFile(node rhino.Node, parent rhino.Node) bool {
	message := "only implemented in ClosureCodingConvention"
	panic(jutil.NewUnsupportedOperationExceptionMsg(message))
}

// Java: CodingConventions.DefaultCodingConvention.extractClassNameIfProvide(Node, Node)
func (d *CodingConventions_DefaultCodingConvention) ExtractClassNameIfProvide(node rhino.Node, parent rhino.Node) *string {
	message := "only implemented in ClosureCodingConvention"
	panic(jutil.NewUnsupportedOperationExceptionMsg(message))
}

// Java: CodingConventions.DefaultCodingConvention.extractClassNameIfRequire(Node, Node)
func (d *CodingConventions_DefaultCodingConvention) ExtractClassNameIfRequire(node rhino.Node, parent rhino.Node) *string {
	message := "only implemented in ClosureCodingConvention"
	panic(jutil.NewUnsupportedOperationExceptionMsg(message))
}

// Java: CodingConventions.DefaultCodingConvention.getExportPropertyFunction()
func (d *CodingConventions_DefaultCodingConvention) GetExportPropertyFunction() *string {
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.getExportSymbolFunction()
func (d *CodingConventions_DefaultCodingConvention) GetExportSymbolFunction() *string {
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.identifyTypeDeclarationCall(Node)
func (d *CodingConventions_DefaultCodingConvention) IdentifyTypeDeclarationCall(n rhino.Node) []string {
	return []string{}
}

// Java: CodingConventions.DefaultCodingConvention.applySubclassRelationship(NominalTypeBuilder, NominalTypeBuilder, SubclassType)
func (d *CodingConventions_DefaultCodingConvention) ApplySubclassRelationship(
	parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	// do nothing
}

// Java: CodingConventions.DefaultCodingConvention.getAbstractMethodName()
func (d *CodingConventions_DefaultCodingConvention) GetAbstractMethodName() *string {
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.getSingletonGetterClassName(Node)
func (d *CodingConventions_DefaultCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) *string {
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.applySingletonGetter(NominalTypeBuilder, FunctionType)
func (d *CodingConventions_DefaultCodingConvention) ApplySingletonGetter(
	classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	// do nothing.
}

// Java: CodingConventions.DefaultCodingConvention.isPropertyTestFunction(Node)
func (d *CodingConventions_DefaultCodingConvention) IsPropertyTestFunction(call rhino.Node) bool {
	// Avoid building the qualified name and check for
	// "goog.isArray"
	target := call.GetFirstChild()
	if target.IsGetProp() {
		src := target.GetFirstChild()
		prop := target.GetString()
		if src.IsName() && src.GetString() == "Array" && prop == "isArray" {
			return true
		}
	}
	return false
}

// Java: CodingConventions.DefaultCodingConvention.isPropertyRenameFunction(Node)
func (d *CodingConventions_DefaultCodingConvention) IsPropertyRenameFunction(nameNode rhino.Node) bool {
	return nameNode.MatchesName(NodeUtil_JSC_PROPERTY_NAME_FN)
}

// Java: CodingConventions.DefaultCodingConvention.isPrototypeAlias(Node)
func (d *CodingConventions_DefaultCodingConvention) IsPrototypeAlias(getProp rhino.Node) bool {
	return false
}

// Java: CodingConventions.DefaultCodingConvention.getObjectLiteralCast(Node)
func (d *CodingConventions_DefaultCodingConvention) GetObjectLiteralCast(callNode rhino.Node) *CodingConvention_ObjectLiteralCast {
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.getAssertionFunctions()
func (d *CodingConventions_DefaultCodingConvention) GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec {
	return []*CodingConvention_AssertionFunctionSpec{
		CodingConvention_AssertionFunctionSpec_ForTruthy().
			SetClosurePrimitive(rhino.ClosurePrimitive_ASSERTS_TRUTHY).
			Build(),
		CodingConvention_AssertionFunctionSpec_ForMatchesReturn().
			SetClosurePrimitive(rhino.ClosurePrimitive_ASSERTS_MATCHES_RETURN).
			Build(),
	}
}

var codingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL = rhino.QualifiedName_Of("Function.prototype.bind.call")

// Java: CodingConventions.DefaultCodingConvention.describeFunctionBind(Node, boolean)
func (d *CodingConventions_DefaultCodingConvention) DescribeFunctionBind(n rhino.Node, checkTypes bool) *CodingConvention_Bind {
	if !n.IsCall() {
		return nil
	}

	callTarget := n.GetFirstChild()
	if callTarget.IsQualifiedName() {
		if codingConventions_DefaultCodingConvention_FUNCTION_PROTOTYPE_BIND_CALL.Matches(callTarget) {
			// goog.bind(fn, self, args...);
			fn := callTarget.GetNext()
			if fn.IsNil() {
				return nil
			}
			thisValue := CodingConventions_DefaultCodingConvention_SafeNext(fn)
			parameters := CodingConventions_DefaultCodingConvention_SafeNext(thisValue)
			return NewCodingConvention_Bind(fn, thisValue, parameters)
		}
	}

	if callTarget.IsGetProp() && callTarget.GetString() == "bind" {
		maybeFn := callTarget.GetFirstChild()
		maybeFnType := maybeFn.GetJSType()
		var fnType rhino.FunctionType
		if checkTypes && maybeFnType != nil {
			fnType = maybeFnType.RestrictByNotNullOrUndefined().
				ToMaybeFunctionType()
		}

		if fnType != nil || maybeFn.IsFunction() {
			// (function(){}).bind(self, args...);
			thisValue := callTarget.GetNext()
			parameters := CodingConventions_DefaultCodingConvention_SafeNext(thisValue)
			return NewCodingConvention_Bind(maybeFn, thisValue, parameters)
		}
	}

	return nil
}

// Java: CodingConventions.DefaultCodingConvention.describeCachingCall(Node)
func (d *CodingConventions_DefaultCodingConvention) DescribeCachingCall(node rhino.Node) *CodingConvention_Cache {
	return nil
}

// Java: CodingConventions.DefaultCodingConvention.safeNext(Node)
func CodingConventions_DefaultCodingConvention_SafeNext(n rhino.Node) rhino.Node {
	if !n.IsNil() {
		return n.GetNext()
	}
	return rhino.NilNode
}
