// Port of com/google/javascript/jscomp/CodingConventions.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ChromeCodingConvention.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	defaultCodingConvention CodingConvention = &CodingConventions_DefaultCodingConvention{}
	emptyPackageName                         = ""

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

	codingConvArrayID   = rhino.RhinoStringPool_InternID("Array")
	codingConvIsArrayID = rhino.RhinoStringPool_InternID("isArray")
	codingConvBindID    = rhino.RhinoStringPool_InternID("bind")
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
	if NodeUtil_IsExprCall(n) {
		fnType := rhino.JSType_ToMaybeFunctionType(n.GetFirstFirstChild().GetJSType())
		return fnType != nil && fnType.GetClosurePrimitive() != nil &&
			*fnType.GetClosurePrimitive() == rhino.ClosurePrimitive_ASSERTS_FAIL
	}
	return false
}

func (d *CodingConventions_DefaultCodingConvention) GetPackageName(source rhino.StaticSourceFile) *string {
	name := source.GetName()
	lastSlash := strings.LastIndexByte(name, '/')
	if lastSlash < 0 {
		return &emptyPackageName
	}
	pkg := name[:lastSlash]
	return &pkg
}

func (d *CodingConventions_DefaultCodingConvention) IsExportedWithLocal(name string, local bool) bool {
	return local && strings.HasPrefix(name, "$super")
}

func (d *CodingConventions_DefaultCodingConvention) IsExported(name string) bool {
	return CodingConvention_IsExported(d, name)
}

func (d *CodingConventions_DefaultCodingConvention) BlockRenamingForProperty(name string) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship {
	callName := callNode.GetFirstChild()
	if (jscompInheritsQName.Matches(callName) || callName.MatchesName("$jscomp$inherits")) &&
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

func (d *CodingConventions_DefaultCodingConvention) IsClassFactoryCall(callNode rhino.Node) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) IsSuperClassReference(propertyName string) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) ExtractIsModuleFile(node rhino.Node, parent rhino.Node) bool {
	panic(jutil.NewUnsupportedOperationExceptionMsg("only implemented in ClosureCodingConvention"))
}

func (d *CodingConventions_DefaultCodingConvention) ExtractClassNameIfProvide(node rhino.Node, parent rhino.Node) *string {
	panic(jutil.NewUnsupportedOperationExceptionMsg("only implemented in ClosureCodingConvention"))
}

func (d *CodingConventions_DefaultCodingConvention) ExtractClassNameIfRequire(node rhino.Node, parent rhino.Node) *string {
	panic(jutil.NewUnsupportedOperationExceptionMsg("only implemented in ClosureCodingConvention"))
}

func (d *CodingConventions_DefaultCodingConvention) GetExportPropertyFunction() *string {
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) GetExportSymbolFunction() *string {
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) IdentifyTypeDeclarationCall(n rhino.Node) []string {
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) ApplySubclassRelationship(
	parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
}

func (d *CodingConventions_DefaultCodingConvention) GetAbstractMethodName() *string {
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) *string {
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) ApplySingletonGetter(
	classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
}

func (d *CodingConventions_DefaultCodingConvention) IsPropertyTestFunction(call rhino.Node) bool {
	callee := call.GetFirstChild()
	if callee.IsGetProp() && callee.GetStringID() == codingConvIsArrayID {
		recv := callee.GetFirstChild()
		return recv.IsName() && recv.GetStringID() == codingConvArrayID
	}
	return false
}

func (d *CodingConventions_DefaultCodingConvention) IsPropertyRenameFunction(nameNode rhino.Node) bool {
	return nameNode.MatchesName(NodeUtil_JSC_PROPERTY_NAME_FN)
}

func (d *CodingConventions_DefaultCodingConvention) IsPrototypeAlias(getProp rhino.Node) bool {
	return false
}

func (d *CodingConventions_DefaultCodingConvention) GetObjectLiteralCast(callNode rhino.Node) *CodingConvention_ObjectLiteralCast {
	return nil
}

func (d *CodingConventions_DefaultCodingConvention) GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec {
	return defaultAssertionFunctions
}

func (d *CodingConventions_DefaultCodingConvention) DescribeFunctionBind(n rhino.Node, checkTypes bool) *CodingConvention_Bind {
	call, ok := ast.AsCall(n)
	if !ok {
		return nil
	}

	callTarget := call.Callee().Node()
	if callTarget.IsQualifiedName() && functionPrototypeBindCallQ.Matches(callTarget) {
		fn := callTarget.GetNext()
		if fn.IsNil() {
			return nil
		}
		thisValue := safeNextNode(fn)
		parameters := safeNextNode(thisValue)
		return NewCodingConvention_Bind(fn, thisValue, parameters)
	}

	if callTarget.IsGetProp() && callTarget.GetStringID() == codingConvBindID {
		maybeFn := callTarget.GetFirstChild()
		maybeFnType := maybeFn.GetJSType()
		var fnType rhino.FunctionType
		if checkTypes && maybeFnType != nil {
			fnType = maybeFnType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
		}
		if fnType != nil || maybeFn.IsFunction() {
			thisValue := callTarget.GetNext()
			parameters := safeNextNode(thisValue)
			return NewCodingConvention_Bind(maybeFn, thisValue, parameters)
		}
	}

	return nil
}

func (d *CodingConventions_DefaultCodingConvention) DescribeCachingCall(node rhino.Node) *CodingConvention_Cache {
	return nil
}

func safeNextNode(n rhino.Node) rhino.Node {
	if !n.IsNil() {
		return n.GetNext()
	}
	return rhino.NilNode
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

func (c *ChromeCodingConvention) BlockRenamingForProperty(name string) bool {
	return false
}

func (c *ChromeCodingConvention) IsPrototypeAlias(getProp rhino.Node) bool {
	return false
}

func (c *ChromeCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) *string {
	callArg := callNode.GetFirstChild()
	if !chromeAddSingletonGetterQ.Matches(callArg) || !callNode.HasTwoChildren() {
		return c.CodingConvention.GetSingletonGetterClassName(callNode)
	}
	return callArg.GetNext().GetQualifiedName()
}

func (c *ChromeCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}
