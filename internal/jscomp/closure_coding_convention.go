// Port of com/google/javascript/jscomp/ClosureCodingConvention.java @ closure-compiler v20261021

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

var (
	closureExportPropertyFunction = "goog.exportProperty"
	closureExportSymbolFunction   = "goog.exportSymbol"
	closureAbstractMethodName     = "goog.abstractMethod"

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

	codingConvPrototypeID               = rhino.RhinoStringPool_InternID("prototype")
	codingConvIsArrayLikeID             = rhino.RhinoStringPool_InternID("isArrayLike")
	codingConvIsObjectID                = rhino.RhinoStringPool_InternID("isObject")
	codingConvGoogReflectCacheMangledID = rhino.RhinoStringPool_InternID("goog$reflect$cache")

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

	callName := callNode.GetFirstChild()
	type_, ok := typeofClassDefiningName(callName)
	if !ok || !callNode.HasXChildren(3) && callNode.GetChildCount() < 3 {
		return nil
	}

	// goog.inherits(SubClass, SuperClass) / ValueType.mixin(SubClass, SuperClass, ...)
	subclass := callName.GetNext()
	superclass := subclass.GetNext()

	if type_ == CodingConvention_SubclassType_MIXIN {
		if endsWithPrototype(superclass) {
			superclass = superclass.GetFirstChild()
		}
		if endsWithPrototype(subclass) {
			subclass = subclass.GetFirstChild()
		}
	}

	// Bail out if either side isn't a real class name (e.g., cond ? A : B).
	if !subclass.IsNil() &&
		subclass.IsUnscopedQualifiedName() &&
		superclass.IsUnscopedQualifiedName() {
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

func endsWithPrototype(qualifiedName rhino.Node) bool {
	return qualifiedName.IsGetProp() && qualifiedName.GetStringID() == codingConvPrototypeID
}

func (c *ClosureCodingConvention) ExtractClassNameIfRequire(node rhino.Node, parent rhino.Node) *string {
	prop, arg := extractGoogCallPropAndArg(node, parent)
	if arg != nil && prop == "require" {
		return arg
	}
	return nil
}

func extractGoogCallPropAndArg(node rhino.Node, parent rhino.Node) (string, *string) {
	if !NodeUtil_IsExprCall(parent) {
		return "", nil
	}
	callee := node.GetFirstChild()
	if callee.IsNil() || !callee.IsGetProp() {
		return "", nil
	}
	recv := callee.GetFirstChild()
	if recv.IsNil() || !recv.IsName() || recv.GetStringID() != rewriteGoogJsGoogID {
		return "", nil
	}
	target := callee.GetNext()
	if target.IsNil() || !target.IsStringLit() {
		return "", nil
	}
	return callee.GetString(), strpool.StringPtr(strpool.ID(target.GetStringID()))
}

func (c *ClosureCodingConvention) GetExportPropertyFunction() *string {
	return &closureExportPropertyFunction
}

func (c *ClosureCodingConvention) GetExportSymbolFunction() *string {
	return &closureExportSymbolFunction
}

func (c *ClosureCodingConvention) IdentifyTypeDeclarationCall(n rhino.Node) []string {
	if googForwardDeclareQ.Matches(n.GetFirstChild()) && n.HasTwoChildren() {
		if lit, ok := ast.AsStringLit(n.GetSecondChild()); ok {
			return []string{lit.Value()}
		}
	}
	return c.CodingConvention.IdentifyTypeDeclarationCall(n)
}

func (c *ClosureCodingConvention) GetAbstractMethodName() *string {
	return &closureAbstractMethodName
}

func (c *ClosureCodingConvention) GetSingletonGetterClassName(callNode rhino.Node) *string {
	callArg := callNode.GetFirstChild()
	if callNode.HasTwoChildren() &&
		(googAddSingletonGetterQ.Matches(callArg) || googAddSingletonGetterMangledQ.Matches(callArg)) {
		return callArg.GetNext().GetQualifiedName()
	}
	return c.CodingConvention.GetSingletonGetterClassName(callNode)
}

func (c *ClosureCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}

func (c *ClosureCodingConvention) IsPropertyTestFunction(call rhino.Node) bool {
	jutil.CheckArgument(call.IsCall())
	callee := call.GetFirstChild()
	if callee.IsGetProp() {
		propID := callee.GetStringID()
		if propID == codingConvIsArrayLikeID || propID == codingConvIsObjectID {
			recv := callee.GetFirstChild()
			if recv.IsName() && recv.GetStringID() == rewriteGoogJsGoogID {
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

func (c *ClosureCodingConvention) GetObjectLiteralCast(callNode rhino.Node) *CodingConvention_ObjectLiteralCast {
	jutil.CheckArgumentMsg1(callNode.IsCall(), "Expected call node but found %s", callNode)
	if proxyCast := c.CodingConvention.GetObjectLiteralCast(callNode); proxyCast != nil {
		return proxyCast
	}

	callName := callNode.GetFirstChild()
	if !(googReflectObjectQ.Matches(callName) || jscompReflectObjectQ.Matches(callName)) ||
		!callNode.HasXChildren(3) {
		return nil
	}

	typeNode := callName.GetNext()
	if !typeNode.IsQualifiedName() {
		return nil
	}

	objectNode := typeNode.GetNext()
	if !objectNode.IsObjectLit() {
		return NewCodingConvention_ObjectLiteralCast(nil, rhino.NilNode, ClosureCodingConvention_OBJECTLIT_EXPECTED)
	}

	return NewCodingConvention_ObjectLiteralCast(typeNode.GetQualifiedName(), objectNode, nil)
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

func (c *ClosureCodingConvention) DescribeFunctionBind(n rhino.Node, checkTypes bool) *CodingConvention_Bind {
	call, ok := ast.AsCall(n)
	if !ok {
		return nil
	}
	callTarget := call.Callee().Node()
	if callTarget.IsQualifiedName() {
		if googBindQ.Matches(callTarget) || callTarget.MatchesName("goog$bind") {
			fn := callTarget.GetNext()
			if fn.IsNil() {
				return nil
			}
			thisValue := safeNextNode(fn)
			parameters := safeNextNode(thisValue)
			return NewCodingConvention_Bind(fn, thisValue, parameters)
		}
		if googPartialQ.Matches(callTarget) || callTarget.MatchesName("goog$partial") {
			fn := callTarget.GetNext()
			if fn.IsNil() {
				return nil
			}
			parameters := safeNextNode(fn)
			return NewCodingConvention_Bind(fn, rhino.NilNode, parameters)
		}
	}
	return c.CodingConvention.DescribeFunctionBind(n, checkTypes)
}

func (c *ClosureCodingConvention) DescribeCachingCall(node rhino.Node) *CodingConvention_Cache {
	call, ok := ast.AsCall(node)
	if !ok {
		return nil
	}
	callTarget := call.Callee().Node()
	if matchesCacheMethodName(callTarget) {
		paramCount := node.GetChildCount() - 1
		if 3 <= paramCount && paramCount <= 4 {
			cacheObj := callTarget.GetNext()
			keyNode := cacheObj.GetNext()
			valueFn := keyNode.GetNext()
			keyFn := valueFn.GetNext()
			return NewCodingConvention_Cache(cacheObj, keyNode, valueFn, keyFn)
		}
	}
	return c.CodingConvention.DescribeCachingCall(node)
}

func matchesCacheMethodName(target rhino.Node) bool {
	if target.IsGetProp() {
		return googReflectCacheQ.Matches(target)
	}
	return target.IsName() && target.GetStringID() == codingConvGoogReflectCacheMangledID
}
