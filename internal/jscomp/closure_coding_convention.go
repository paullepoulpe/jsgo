// Port of com/google/javascript/jscomp/ClosureCodingConvention.java @ closure-compiler v20261028

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// ClosureCodingConvention describes the Closure-specific JavaScript coding
// conventions.
type ClosureCodingConvention struct {
	CodingConvention
}

var ClosureCodingConvention_OBJECTLIT_EXPECTED = DiagnosticType_Warning(
	"JSC_REFLECT_OBJECTLIT_EXPECTED",
	"Object literal expected as second argument")

const (
	closureAbstractMethodName     = "goog.abstractMethod"
	closureExportPropertyFunction = "goog.exportProperty"
	closureExportSymbolFunction   = "goog.exportSymbol"
)

var (
	googForwardDeclareQ               = rhino.QualifiedName_Of("goog.forwardDeclare")
	googAddSingletonGetterQ           = rhino.QualifiedName_Of("goog.addSingletonGetter")
	googAddSingletonGetterMangledQ    = rhino.QualifiedName_Of("goog$addSingletonGetter")
	googReflectObjectPropertyQ        = rhino.QualifiedName_Of("goog.reflect.objectProperty")
	googReflectObjectPropertyMangledQ = rhino.QualifiedName_Of("goog$reflect$objectProperty")
	googAssertsFailQ                  = rhino.QualifiedName_Of("goog.asserts.fail")
	googReflectObjectQ                = rhino.QualifiedName_Of("goog.reflect.object")
	jscompReflectObjectQ              = rhino.QualifiedName_Of("$jscomp.reflectObject")
	googBindQ                         = rhino.QualifiedName_Of("goog.bind")
	googPartialQ                      = rhino.QualifiedName_Of("goog.partial")
	googReflectCacheQ                 = rhino.QualifiedName_Of("goog.reflect.cache")

	codingConvPrototypeID               = strpool.Intern("prototype")
	codingConvIsArrayLikeID             = strpool.Intern("isArrayLike")
	codingConvIsObjectID                = strpool.Intern("isObject")
	codingConvGoogReflectCacheMangledID = strpool.Intern("goog$reflect$cache")

	closureOwnAssertionFunctions = []*CodingConvention_AssertionFunctionSpec{
		newAssertionSpecFunc("goog.asserts.assert", CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY),
		newAssertionSpecFunc("goog.asserts.assertArray", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertBoolean", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertElement", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertFunction", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertInstanceof", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertNumber", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertObject", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
		newAssertionSpecFunc("goog.asserts.assertString", CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE),
	}

	closureDefaultAssertionFunctions = slices.Concat(defaultAssertionFunctions, closureOwnAssertionFunctions)
)

func NewClosureCodingConvention() *ClosureCodingConvention {
	return NewClosureCodingConventionWithConvention(CodingConventions_GetDefault())
}

func NewClosureCodingConventionWithConvention(wrapped CodingConvention) *ClosureCodingConvention {
	return &ClosureCodingConvention{CodingConvention: wrapped}
}

func (c *ClosureCodingConvention) IsExported(name string) bool {
	return CodingConvention_IsExported(c, name)
}

// ApplySubclassRelationship adds Closure's goog.inherits superClass_ and
// prototype.constructor properties.
func (c *ClosureCodingConvention) ApplySubclassRelationship(
	parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	c.CodingConvention.ApplySubclassRelationship(parent, child, type_)
	if type_ == CodingConvention_SubclassType_INHERITS {
		childCtor := child.Constructor()
		child.DeclareConstructorProperty(
			"superClass_", parent.PrototypeOrInstance(), childCtor.GetSource())
		// Notice that constructor functions do not need to be covariant on the superclass.
		// So if G extends F, new G() and new F() can accept completely different argument
		// types, but G.prototype.constructor needs to be covariant on F.prototype.constructor.
		// To get around this, we just turn off type-checking on arguments and return types
		// of G.prototype.constructor.
		qmarkCtor := childCtor.ForgetParameterAndReturnTypes()
		child.DeclarePrototypeProperty("constructor", qmarkCtor, childCtor.GetSource())
	}
}

// GetClassesDefinedByCall recognizes Google/Closure inherits and mixin calls.
func (c *ClosureCodingConvention) GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship {
	if relationship := c.CodingConvention.GetClassesDefinedByCall(callNode); relationship != nil {
		return relationship
	}

	call, ok := ast.AsCallOrNew(callNode)
	if !ok {
		return nil
	}
	type_, ok := typeofClassDefiningName(call.Callee().Node())
	if !ok {
		return nil
	}
	subclassArg, ok := call.ArgAt(0)
	if !ok {
		return nil
	}
	superclassArg, ok := call.ArgAt(1)
	if !ok {
		return nil
	}

	// goog.inherits(SubClass, SuperClass) / ValueType.mixin(SubClass, SuperClass, ...)
	subclass := subclassArg.Node()
	superclass := superclassArg.Node()

	if type_ == CodingConvention_SubclassType_MIXIN {
		superclass = stripPrototypeProp(superclass)
		subclass = stripPrototypeProp(subclass)
	}

	// Bail out if either side isn't a real class name (e.g., cond ? A : B).
	if subclass.IsUnscopedQualifiedName() && superclass.IsUnscopedQualifiedName() {
		return NewCodingConvention_SubclassRelationship(type_, subclass, superclass)
	}
	return nil
}

func typeofClassDefiningName(callName rhino.Node) (CodingConvention_SubclassType, bool) {
	var methodName string
	if gp, ok := ast.AsGetProp(callName); ok {
		methodName = gp.Property()
	} else if nm, ok := ast.AsName(callName); ok {
		name := nm.Ident()
		idx := strings.LastIndexByte(name, '$')
		if idx < 0 {
			return 0, false
		}
		methodName = name[idx+1:]
	} else {
		return 0, false
	}

	switch methodName {
	case "inherits":
		return CodingConvention_SubclassType_INHERITS, true
	case "mixin":
		return CodingConvention_SubclassType_MIXIN, true
	default:
		return 0, false
	}
}

func (c *ClosureCodingConvention) IsSuperClassReference(propertyName string) bool {
	return propertyName == "superClass_" ||
		c.CodingConvention.IsSuperClassReference(propertyName)
}

func stripPrototypeProp(qualifiedName rhino.Node) rhino.Node {
	if gp, ok := ast.AsGetProp(qualifiedName); ok && nodeStringID(qualifiedName) == codingConvPrototypeID {
		return gp.Object().Node()
	}
	return qualifiedName
}

func (c *ClosureCodingConvention) ExtractClassNameIfRequire(node rhino.Node, parent rhino.Node) (string, bool) {
	prop, arg, ok := extractGoogCallPropAndArg(node, parent)
	if ok && prop == "require" {
		return arg, true
	}
	return "", false
}

func extractGoogCallPropAndArg(node rhino.Node, parent rhino.Node) (string, string, bool) {
	if !NodeUtil_IsExprCall(parent) {
		return "", "", false
	}
	call, ok := ast.AsCall(node)
	if !ok {
		return "", "", false
	}
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return "", "", false
	}
	recv, ok := callee.Object().AsName()
	if !ok || nodeStringID(recv.Node()) != rewriteGoogJsGoogID {
		return "", "", false
	}
	firstArg, ok := call.FirstArg()
	if !ok {
		return "", "", false
	}
	target, ok := firstArg.AsStringLit()
	if !ok {
		return "", "", false
	}
	return callee.Property(), target.Value(), true
}

func (c *ClosureCodingConvention) GetExportPropertyFunction() (string, bool) {
	return closureExportPropertyFunction, true
}

func (c *ClosureCodingConvention) GetExportSymbolFunction() (string, bool) {
	return closureExportSymbolFunction, true
}

func (c *ClosureCodingConvention) IdentifyTypeDeclarationCall(n rhino.Node) (string, bool) {
	if call, ok := ast.AsCall(n); ok && googForwardDeclareQ.Matches(call.Callee().Node()) {
		if arg, ok := call.OnlyArg(); ok {
			if lit, ok := arg.AsStringLit(); ok {
				return lit.Value(), true
			}
		}
	}
	return c.CodingConvention.IdentifyTypeDeclarationCall(n)
}

func (c *ClosureCodingConvention) GetAbstractMethodName() string {
	return closureAbstractMethodName
}

func (c *ClosureCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) (string, bool) {
	if call, ok := ast.AsCall(callNode); ok {
		callee := call.Callee().Node()
		if arg, ok := call.OnlyArg(); ok &&
			(googAddSingletonGetterQ.Matches(callee) || googAddSingletonGetterMangledQ.Matches(callee)) {
			return arg.Node().GetQualifiedNameString()
		}
	}
	return c.CodingConvention.GetSingletonGetterClassName(callNode)
}

func (c *ClosureCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}

func (c *ClosureCodingConvention) IsPropertyTestFunction(call rhino.Node) bool {
	callView, ok := ast.AsCall(call)
	jutil.CheckArgument(ok)
	if callee, ok := callView.Callee().AsGetProp(); ok {
		propID := nodeStringID(callee.Node())
		if propID == codingConvIsArrayLikeID || propID == codingConvIsObjectID {
			recv := callee.Object().Node()
			if recv.IsName() && nodeStringID(recv) == rewriteGoogJsGoogID {
				return true
			}
		}
	}
	return c.CodingConvention.IsPropertyTestFunction(call)
}

func (c *ClosureCodingConvention) IsPropertyRenameFunction(nameNode rhino.Node) bool {
	return c.CodingConvention.IsPropertyRenameFunction(nameNode) ||
		googReflectObjectPropertyQ.Matches(nameNode) ||
		googReflectObjectPropertyMangledQ.Matches(nameNode)
}

func (c *ClosureCodingConvention) IsFunctionCallThatAlwaysThrows(n rhino.Node) bool {
	return c.CodingConvention.IsFunctionCallThatAlwaysThrows(n) ||
		callAlwaysThrowsQName(n, googAssertsFailQ)
}

func (c *ClosureCodingConvention) GetObjectLiteralCast(callNode rhino.Node) (CodingConvention_ObjectLiteralCast, bool) {
	call, ok := ast.AsCall(callNode)
	jutil.CheckArgumentMsg1(ok, "Expected call node but found %s", callNode)
	if proxyCast, ok := c.CodingConvention.GetObjectLiteralCast(callNode); ok {
		return proxyCast, true
	}

	callName := call.Callee().Node()
	if !(googReflectObjectQ.Matches(callName) || jscompReflectObjectQ.Matches(callName)) ||
		call.NumArgs() != 2 {
		return CodingConvention_ObjectLiteralCast{}, false
	}

	callOrNew := call.AsCallOrNew()
	typeArg, _ := callOrNew.ArgAt(0)
	typeNode := typeArg.Node()
	if !typeNode.IsQualifiedName() {
		return CodingConvention_ObjectLiteralCast{}, false
	}

	objectArg, _ := callOrNew.ArgAt(1)
	objectNode := objectArg.Node()
	if !objectNode.IsObjectLit() {
		return NewCodingConvention_ObjectLiteralCast("", rhino.NilNode, ClosureCodingConvention_OBJECTLIT_EXPECTED), true
	}

	typeName, _ := typeNode.GetQualifiedNameString()
	return NewCodingConvention_ObjectLiteralCast(typeName, objectNode, nil), true
}

func (c *ClosureCodingConvention) GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec {
	if c.CodingConvention == defaultCodingConvention {
		return closureDefaultAssertionFunctions
	}
	base := c.CodingConvention.GetAssertionFunctions()
	out := make([]*CodingConvention_AssertionFunctionSpec, 0, len(base)+len(closureOwnAssertionFunctions))
	addUnique := func(spec *CodingConvention_AssertionFunctionSpec) {
		for _, existing := range out {
			if existing.equals(spec) {
				return
			}
		}
		out = append(out, spec)
	}
	for _, spec := range base {
		addUnique(spec)
	}
	for _, spec := range closureOwnAssertionFunctions {
		addUnique(spec)
	}
	return out
}

func (c *ClosureCodingConvention) DescribeFunctionBind(n rhino.Node, checkTypes bool) (CodingConvention_Bind, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return CodingConvention_Bind{}, false
	}
	callTarget := call.Callee().Node()
	if callTarget.IsQualifiedName() {
		callOrNew := call.AsCallOrNew()
		if googBindQ.Matches(callTarget) || callTarget.MatchesName("goog$bind") {
			fn, ok := call.FirstArg()
			if !ok {
				return CodingConvention_Bind{}, false
			}
			thisValue, _ := callOrNew.ArgAt(1)
			parameters, _ := callOrNew.ArgAt(2)
			return NewCodingConvention_Bind(fn.Node(), thisValue.Node(), parameters.Node()), true
		}
		if googPartialQ.Matches(callTarget) || callTarget.MatchesName("goog$partial") {
			fn, ok := call.FirstArg()
			if !ok {
				return CodingConvention_Bind{}, false
			}
			parameters, _ := callOrNew.ArgAt(1)
			return NewCodingConvention_Bind(fn.Node(), rhino.NilNode, parameters.Node()), true
		}
	}
	return c.CodingConvention.DescribeFunctionBind(n, checkTypes)
}

func (c *ClosureCodingConvention) DescribeCachingCall(node rhino.Node) (CodingConvention_Cache, bool) {
	call, ok := ast.AsCall(node)
	if !ok {
		return CodingConvention_Cache{}, false
	}
	callTarget := call.Callee().Node()
	if matchesCacheMethodName(callTarget) {
		paramCount := call.NumArgs()
		if 3 <= paramCount && paramCount <= 4 {
			callOrNew := call.AsCallOrNew()
			cacheObj, _ := callOrNew.ArgAt(0)
			keyNode, _ := callOrNew.ArgAt(1)
			valueFn, _ := callOrNew.ArgAt(2)
			keyFn, _ := callOrNew.ArgAt(3)
			return NewCodingConvention_Cache(cacheObj.Node(), keyNode.Node(), valueFn.Node(), keyFn.Node()), true
		}
	}
	return c.CodingConvention.DescribeCachingCall(node)
}

func matchesCacheMethodName(target rhino.Node) bool {
	if target.IsGetProp() {
		return googReflectCacheQ.Matches(target)
	}
	return target.IsName() && nodeStringID(target) == codingConvGoogReflectCacheMangledID
}
