// Port of com/google/javascript/jscomp/CodingConventions.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/ChromeCodingConvention.java @ closure-compiler v20261028

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

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	defaultCodingConvention CodingConvention = &CodingConventions_DefaultCodingConvention{}

	defaultAssertionFunctions = []*CodingConvention_AssertionFunctionSpec{
		newAssertionSpecPrimitive(
			rhino.ClosurePrimitive_ASSERTS_TRUTHY,
			CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY,
		),
		newAssertionSpecPrimitive(
			rhino.ClosurePrimitive_ASSERTS_MATCHES_RETURN,
			CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE,
		),
	}

	jscompInheritsQName        = rhino.QualifiedName_Of("$jscomp.inherits")
	functionPrototypeBindCallQ = rhino.QualifiedName_Of("Function.prototype.bind.call")
	chromeAddSingletonGetterQ  = rhino.QualifiedName_Of("cr.addSingletonGetter")

	codingConvArrayID   = strpool.Intern("Array")
	codingConvIsArrayID = strpool.Intern("isArray")
	codingConvBindID    = strpool.Intern("bind")
)

// CodingConventions_GetDefault gets the default coding convention.
func CodingConventions_GetDefault() CodingConvention {
	return defaultCodingConvention
}

// CodingConventions_DefaultIsFunctionCallThatAlwaysThrows reports whether n
// (an EXPR_RESULT(CALL) or CALL node) calls alwaysThrowsFunctionName.
func CodingConventions_DefaultIsFunctionCallThatAlwaysThrows(
	n rhino.Node, alwaysThrowsFunctionName string) bool {
	if stmt, ok := ast.AsExprStmt(n); ok {
		n = stmt.Expr().Node()
	}
	call, ok := ast.AsCall(n)
	return ok && call.Callee().Node().MatchesQualifiedName(alwaysThrowsFunctionName)
}

func callAlwaysThrowsQName(n rhino.Node, qname rhino.QualifiedName) bool {
	if stmt, ok := ast.AsExprStmt(n); ok {
		n = stmt.Expr().Node()
	}
	call, ok := ast.AsCall(n)
	return ok && qname.Matches(call.Callee().Node())
}

// CodingConventions_DefaultCodingConvention is the default coding convention
// at the bottom of convention chains.
type CodingConventions_DefaultCodingConvention struct{}

func (d *CodingConventions_DefaultCodingConvention) IsConstant(variableName string) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) IsConstantKey(variableName string) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) IsValidEnumKey(key string) bool {
	return len(key) > 0
}

func (d *CodingConventions_DefaultCodingConvention) IsOptionalParameter(parameter rhino.Node) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) IsVarArgsParameter(parameter rhino.Node) bool {
	return parameter.IsRest()
}

func (d *CodingConventions_DefaultCodingConvention) IsFunctionCallThatAlwaysThrows(n rhino.Node) bool {
	stmt, ok := ast.AsExprStmt(n)
	if !ok {
		return false
	}
	call, ok := stmt.Expr().AsCall()
	if !ok {
		return false
	}
	callee := call.Callee().Node()
	if callee.IsNil() {
		return false
	}
	if fnType := rhino.JSType_ToMaybeFunctionType(callee.GetJSType()); fnType != nil {
		prim, ok := fnType.GetClosurePrimitive()
		return ok && prim == rhino.ClosurePrimitive_ASSERTS_FAIL
	}
	return false
}

func (d *CodingConventions_DefaultCodingConvention) GetPackageName(source rhino.StaticSourceFile) (string, bool) {
	name := source.GetName()
	lastSlash := strings.LastIndexByte(name, '/')
	if lastSlash < 0 {
		return "", true
	}
	return name[:lastSlash], true
}

func (d *CodingConventions_DefaultCodingConvention) IsExportedWithLocal(name string, local bool) bool {
	return local && strings.HasPrefix(name, "$super")
}

func (d *CodingConventions_DefaultCodingConvention) IsExported(name string) bool {
	return CodingConvention_IsExported(d, name)
}

func (d *CodingConventions_DefaultCodingConvention) GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship {
	call, ok := ast.AsCallOrNew(callNode)
	if !ok {
		return nil
	}
	callName := call.Callee().Node()
	if !jscompInheritsQName.Matches(callName) && !callName.MatchesName("$jscomp$inherits") {
		return nil
	}
	subclassArg, ok := call.FirstArg()
	if !ok {
		return nil
	}
	superclassArg, ok := subclassArg.Next()
	if !ok || !superclassArg.IsLast() {
		return nil
	}
	subclass := subclassArg.Node()
	superclass := superclassArg.Node()
	// The StripCode pass may create $jscomp.inherits calls with NULL arguments.
	if subclass.IsQualifiedName() && superclass.IsQualifiedName() {
		return NewCodingConvention_SubclassRelationship(CodingConvention_SubclassType_INHERITS, subclass, superclass)
	}
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) IsClassFactoryCall(callNode rhino.Node) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) IsSuperClassReference(propertyName string) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) GetExportPropertyFunction() (string, bool) {
	return "", false
}

func (d *CodingConventions_DefaultCodingConvention) GetExportSymbolFunction() (string, bool) {
	return "", false
}

func (d *CodingConventions_DefaultCodingConvention) IdentifyTypeDeclarationCall(n rhino.Node) (string, bool) {
	return "", false
}

func (d *CodingConventions_DefaultCodingConvention) ApplySubclassRelationship(
	parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
}

func (d *CodingConventions_DefaultCodingConvention) GetAbstractMethodName() string {
	return ""
}

func (d *CodingConventions_DefaultCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) (string, bool) {
	return "", false
}

func (d *CodingConventions_DefaultCodingConvention) ApplySingletonGetter(
	classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
}

func (d *CodingConventions_DefaultCodingConvention) IsPropertyTestFunction(call rhino.Node) bool {
	callView, ok := ast.AsCall(call)
	if !ok {
		return false
	}
	callee, ok := callView.Callee().AsGetProp()
	if !ok || nodeStringID(callee.Node()) != codingConvIsArrayID {
		return false
	}
	recv := callee.Object().Node()
	return recv.IsName() && nodeStringID(recv) == codingConvArrayID
}

func (d *CodingConventions_DefaultCodingConvention) IsPropertyRenameFunction(nameNode rhino.Node) bool {
	return nameNode.MatchesName(NodeUtil_JSC_PROPERTY_NAME_FN)
}

func (d *CodingConventions_DefaultCodingConvention) GetObjectLiteralCast(callNode rhino.Node) (CodingConvention_ObjectLiteralCast, bool) {
	return CodingConvention_ObjectLiteralCast{}, false
}

func (d *CodingConventions_DefaultCodingConvention) GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec {
	return defaultAssertionFunctions
}

func (d *CodingConventions_DefaultCodingConvention) DescribeFunctionBind(n rhino.Node, checkTypes bool) (CodingConvention_Bind, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return CodingConvention_Bind{}, false
	}

	callOrNew := call.AsCallOrNew()
	callTarget := call.Callee().Node()
	if callTarget.IsQualifiedName() && functionPrototypeBindCallQ.Matches(callTarget) {
		fn, ok := call.FirstArg()
		if !ok {
			return CodingConvention_Bind{}, false
		}
		thisValue, _ := callOrNew.ArgAt(1)
		parameters, _ := callOrNew.ArgAt(2)
		return NewCodingConvention_Bind(fn.Node(), thisValue.Node(), parameters.Node()), true
	}

	if gp, ok := call.Callee().AsGetProp(); ok && nodeStringID(gp.Node()) == codingConvBindID {
		maybeFn := gp.Object().Node()
		if maybeFn.IsNil() {
			return CodingConvention_Bind{}, false
		}
		maybeFnType := maybeFn.GetJSType()
		var fnType rhino.FunctionType
		if checkTypes && maybeFnType != nil {
			fnType = maybeFnType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
		}
		if fnType != nil || maybeFn.IsFunction() {
			thisValue, _ := callOrNew.ArgAt(0)
			parameters, _ := callOrNew.ArgAt(1)
			return NewCodingConvention_Bind(maybeFn, thisValue.Node(), parameters.Node()), true
		}
	}

	return CodingConvention_Bind{}, false
}

func (d *CodingConventions_DefaultCodingConvention) DescribeCachingCall(node rhino.Node) (CodingConvention_Cache, bool) {
	return CodingConvention_Cache{}, false
}

// ChromeCodingConvention is the coding convention used by the Chrome team to
// compile Chrome's JS.
type ChromeCodingConvention struct {
	CodingConvention
}

func NewChromeCodingConvention() *ChromeCodingConvention {
	return NewChromeCodingConventionWithConvention(CodingConventions_GetDefault())
}

func NewChromeCodingConventionWithConvention(wrapped CodingConvention) *ChromeCodingConvention {
	return &ChromeCodingConvention{CodingConvention: wrapped}
}

func (c *ChromeCodingConvention) IsExported(name string) bool {
	return CodingConvention_IsExported(c, name)
}

func (c *ChromeCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) (string, bool) {
	if call, ok := ast.AsCall(callNode); ok {
		if arg, ok := call.OnlyArg(); ok && chromeAddSingletonGetterQ.Matches(call.Callee().Node()) {
			return arg.Node().GetQualifiedNameString()
		}
	}
	return c.CodingConvention.GetSingletonGetterClassName(callNode)
}

func (c *ChromeCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}
