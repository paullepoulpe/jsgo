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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureCodingConvention describes the Closure-specific JavaScript coding
// conventions.
type ClosureCodingConvention struct {
	CodingConventions_ProxyBase
}

const closureCodingConvention_serialVersionUID int64 = 1

var ClosureCodingConvention_OBJECTLIT_EXPECTED = DiagnosticType_Warning(
	"JSC_REFLECT_OBJECTLIT_EXPECTED",
	"Object literal expected as second argument")

// Java: ClosureCodingConvention()
func NewClosureCodingConvention() *ClosureCodingConvention {
	return NewClosureCodingConventionWithConvention(CodingConventions_GetDefault())
}

// Java: ClosureCodingConvention(CodingConvention)
func NewClosureCodingConventionWithConvention(wrapped CodingConvention) *ClosureCodingConvention {
	c := &ClosureCodingConvention{}
	c.self = c
	c.InitCodingConventions_ProxyBase(wrapped)
	return c
}

// ApplySubclassRelationship: Closure's goog.inherits adds a superClass_
// property to the subclass, and a constructor property.
//
// Java: ClosureCodingConvention.applySubclassRelationship(NominalTypeBuilder, NominalTypeBuilder, SubclassType)
func (c *ClosureCodingConvention) ApplySubclassRelationship(
	parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType) {
	c.CodingConventions_ProxyBase.ApplySubclassRelationship(parent, child, type_)
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

// GetClassesDefinedByCall understands several different inheritance
// patterns that occur in Google code (various uses of inherits and mixin).
//
// Java: ClosureCodingConvention.getClassesDefinedByCall(Node)
func (c *ClosureCodingConvention) GetClassesDefinedByCall(callNode *rhino.Node) *CodingConvention_SubclassRelationship {
	relationship :=
		c.CodingConventions_ProxyBase.GetClassesDefinedByCall(callNode)
	if relationship != nil {
		return relationship
	}

	callName := callNode.GetFirstChild()
	type_ := ClosureCodingConvention_TypeofClassDefiningName(callName)
	if type_ != nil {
		// Possible formats for a class-defining method call:
		// goog.inherits(SubClass, SuperClass)
		// goog$inherits(SubClass, SuperClass)
		// ValueType.mixin(SubClass, SuperClass, ...) // used by J2CL.
		// ValueType$mixin(SubClass, SuperClass, ...)
		if callNode.GetChildCount() < 3 {
			return nil
		}

		// goog.inherits(SubClass, SuperClass)
		subclass := callName.GetNext()
		superclass := subclass.GetNext()

		if *type_ == CodingConvention_SubclassType_MIXIN {
			// Strip off the prototype from the name.
			if ClosureCodingConvention_EndsWithPrototype(superclass) {
				superclass = superclass.GetFirstChild()
			}
			if ClosureCodingConvention_EndsWithPrototype(subclass) {
				subclass = subclass.GetFirstChild()
			}
		}

		// bail out if either of the side of the "inherits"
		// isn't a real class name. This prevents us from
		// doing something weird in cases like:
		// goog.inherits(MySubClass, cond ? SuperClass1 : BaseClass2)
		if subclass != nil &&
			subclass.IsUnscopedQualifiedName() &&
			superclass.IsUnscopedQualifiedName() {
			return NewCodingConvention_SubclassRelationship(*type_, subclass, superclass)
		}
	}

	return nil
}

// closureCodingConvention_subclassTypes holds the SubclassType constants
// that typeofClassDefiningName returns.
//
// Go-only: language-adapter — a Java enum reference may be null, so the
// result is a pointer; these are the shared non-null values.
var closureCodingConvention_subclassTypes = [...]CodingConvention_SubclassType{
	CodingConvention_SubclassType_INHERITS, CodingConvention_SubclassType_MIXIN,
}

// ClosureCodingConvention_TypeofClassDefiningName determines whether the
// given node is a class-defining name, like "inherits" or "mixin."
//
// Returns the type of class-defining name, or nil.
//
// Java: ClosureCodingConvention.typeofClassDefiningName(Node)
func ClosureCodingConvention_TypeofClassDefiningName(callName *rhino.Node) *CodingConvention_SubclassType {
	// Check if the method name matches one of the class-defining methods.
	var methodName *string
	if callName.IsGetProp() {
		s := callName.GetString()
		methodName = &s
	} else if callName.IsName() {
		name := callName.GetString()
		dollarIndex := jutil.StringLastIndexOf(name, "$")
		if dollarIndex != -1 {
			s := jutil.StringSubstringFrom(name, dollarIndex+1)
			methodName = &s
		}
	}

	if methodName != nil {
		if *methodName == "inherits" {
			return &closureCodingConvention_subclassTypes[CodingConvention_SubclassType_INHERITS]
		} else if *methodName == "mixin" {
			return &closureCodingConvention_subclassTypes[CodingConvention_SubclassType_MIXIN]
		}
	}
	return nil
}

// Java: ClosureCodingConvention.isSuperClassReference(String)
func (c *ClosureCodingConvention) IsSuperClassReference(propertyName string) bool {
	return "superClass_" == propertyName ||
		c.CodingConventions_ProxyBase.IsSuperClassReference(propertyName)
}

// ClosureCodingConvention_EndsWithPrototype: given a qualified name node,
// returns whether "prototype" is at the end. For example: a.b.c => false
// a.b.c.prototype => true
//
// Java: ClosureCodingConvention.endsWithPrototype(Node)
func ClosureCodingConvention_EndsWithPrototype(qualifiedName *rhino.Node) bool {
	return qualifiedName.IsGetProp() && qualifiedName.GetString() == "prototype"
}

// ExtractIsModuleFile returns whether the node indicates that the file
// represents a "module", a file whose top level declarations are not in
// global scope.
//
// Java: ClosureCodingConvention.extractIsModuleFile(Node, Node)
func (c *ClosureCodingConvention) ExtractIsModuleFile(node *rhino.Node, parent *rhino.Node) bool {
	namespace := ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.module")
	return namespace != nil
}

// ExtractClassNameIfProvide extracts X from goog.provide('X'), if the
// applied Node is goog.
//
// Returns the extracted class name, or nil.
//
// Java: ClosureCodingConvention.extractClassNameIfProvide(Node, Node)
func (c *ClosureCodingConvention) ExtractClassNameIfProvide(node *rhino.Node, parent *rhino.Node) *string {
	namespace := ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.provide")
	if namespace == nil {
		namespace = ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.module")
	}
	return namespace
}

// ExtractClassNameIfRequire extracts X from goog.require('X'), if the
// applied Node is goog.
//
// Returns the extracted class name, or nil.
//
// Java: ClosureCodingConvention.extractClassNameIfRequire(Node, Node)
func (c *ClosureCodingConvention) ExtractClassNameIfRequire(node *rhino.Node, parent *rhino.Node) *string {
	return ClosureCodingConvention_ExtractClassNameIfGoog(node, parent, "goog.require")
}

// Java: ClosureCodingConvention.extractClassNameIfGoog(Node, Node, String)
func ClosureCodingConvention_ExtractClassNameIfGoog(node *rhino.Node, parent *rhino.Node,
	functionName string) *string {
	var className *string
	if NodeUtil_IsExprCall(parent) {
		callee := node.GetFirstChild()
		if callee != nil && callee.IsGetProp() && callee.MatchesQualifiedName(functionName) {
			target := callee.GetNext()
			if target != nil && target.IsStringLit() {
				s := target.GetString()
				className = &s
			}
		}
	}
	return className
}

// GetExportPropertyFunction uses closure's implementation.
//
// Returns closure's function name for exporting properties.
//
var (
	closureCodingConvention_exportPropertyFunction = "goog.exportProperty"
	closureCodingConvention_exportSymbolFunction   = "goog.exportSymbol"
	closureCodingConvention_abstractMethodName     = "goog.abstractMethod"
)

// Java: ClosureCodingConvention.getExportPropertyFunction()
func (c *ClosureCodingConvention) GetExportPropertyFunction() *string {
	return &closureCodingConvention_exportPropertyFunction
}

// GetExportSymbolFunction uses closure's implementation.
//
// Returns closure's function name for exporting symbols.
//
// Java: ClosureCodingConvention.getExportSymbolFunction()
func (c *ClosureCodingConvention) GetExportSymbolFunction() *string {
	return &closureCodingConvention_exportSymbolFunction
}

var ClosureCodingConvention_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// Java: ClosureCodingConvention.identifyTypeDeclarationCall(Node)
func (c *ClosureCodingConvention) IdentifyTypeDeclarationCall(n *rhino.Node) []string {
	callName := n.GetFirstChild()
	// Identify forward declaration of form goog.forwardDeclare('foo.bar')
	if ClosureCodingConvention_GOOG_FORWARDDECLARE.Matches(callName) && n.HasTwoChildren() {
		typeDeclaration := n.GetSecondChild()
		if typeDeclaration.IsStringLit() {
			return []string{typeDeclaration.GetString()}
		}
	}

	return c.CodingConventions_ProxyBase.IdentifyTypeDeclarationCall(n)
}

// Java: ClosureCodingConvention.getAbstractMethodName()
func (c *ClosureCodingConvention) GetAbstractMethodName() *string {
	return &closureCodingConvention_abstractMethodName
}

var ClosureCodingConvention_GOOG_ADDSINGLETONGETTER = rhino.QualifiedName_Of("goog.addSingletonGetter")

var ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED = rhino.QualifiedName_Of("goog$addSingletonGetter")

// Java: ClosureCodingConvention.getSingletonGetterClassName(Node)
func (c *ClosureCodingConvention) GetSingletonGetterClassName(callNode *rhino.Node) *string {
	callArg := callNode.GetFirstChild()
	// Use both the original name and the post-CollapseProperties name.
	if callNode.HasTwoChildren() &&
		(ClosureCodingConvention_GOOG_ADDSINGLETONGETTER.Matches(callArg) ||
			ClosureCodingConvention_GOOG_ADDSINGLETONGETTER_MANGLED.Matches(callArg)) {
		return callArg.GetNext().GetQualifiedName()
	}
	return c.CodingConventions_ProxyBase.GetSingletonGetterClassName(callNode)
}

// Java: ClosureCodingConvention.applySingletonGetter(NominalTypeBuilder, FunctionType)
func (c *ClosureCodingConvention) ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType) {
	defSite := classType.Constructor().GetSource()
	classType.DeclareConstructorProperty("getInstance", getterType, defSite)
	classType.DeclareConstructorProperty("instance_", classType.Instance(), defSite)
}

// Java: ClosureCodingConvention.isPropertyTestFunction(Node)
func (c *ClosureCodingConvention) IsPropertyTestFunction(call *rhino.Node) bool {
	jutil.CheckArgument(call.IsCall())
	// Avoid building the qualified name and check for
	// "goog.isArrayLike", "goog.isObject"
	target := call.GetFirstChild()
	if target.IsGetProp() {
		src := target.GetFirstChild()
		prop := target.GetString()
		if src.IsName() &&
			src.GetString() == "goog" &&
			(prop == "isArrayLike" || prop == "isObject") {
			return true
		}
	}

	return c.CodingConventions_ProxyBase.IsPropertyTestFunction(call)
}

var ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY = rhino.QualifiedName_Of("goog.reflect.objectProperty")

var ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED = rhino.QualifiedName_Of("goog$reflect$objectProperty")

// Java: ClosureCodingConvention.isPropertyRenameFunction(Node)
func (c *ClosureCodingConvention) IsPropertyRenameFunction(nameNode *rhino.Node) bool {
	if c.CodingConventions_ProxyBase.IsPropertyRenameFunction(nameNode) {
		return true
	}
	return ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY.Matches(nameNode) ||
		ClosureCodingConvention_GOOG_REFLECT_OBJECTPROPERTY_MANGLED.Matches(nameNode)
}

// Java: ClosureCodingConvention.isFunctionCallThatAlwaysThrows(Node)
func (c *ClosureCodingConvention) IsFunctionCallThatAlwaysThrows(n *rhino.Node) bool {
	return c.CodingConventions_ProxyBase.IsFunctionCallThatAlwaysThrows(n) ||
		CodingConventions_DefaultIsFunctionCallThatAlwaysThrows(n, "goog.asserts.fail")
}

var ClosureCodingConvention_GOOG_REFLECT_OBJECT = rhino.QualifiedName_Of("goog.reflect.object")

var ClosureCodingConvention_JSCOMP_REFLECTOBJECT = rhino.QualifiedName_Of("$jscomp.reflectObject")

// Java: ClosureCodingConvention.getObjectLiteralCast(Node)
func (c *ClosureCodingConvention) GetObjectLiteralCast(callNode *rhino.Node) *CodingConvention_ObjectLiteralCast {
	jutil.CheckArgumentMsg1(callNode.IsCall(), "Expected call node but found %s", callNode)
	proxyCast := c.CodingConventions_ProxyBase.GetObjectLiteralCast(callNode)
	if proxyCast != nil {
		return proxyCast
	}

	callName := callNode.GetFirstChild()
	if !(ClosureCodingConvention_GOOG_REFLECT_OBJECT.Matches(callName) || ClosureCodingConvention_JSCOMP_REFLECTOBJECT.Matches(callName)) ||
		!callNode.HasXChildren(3) {
		return nil
	}

	typeNode := callName.GetNext()
	if !typeNode.IsQualifiedName() {
		return nil
	}

	objectNode := typeNode.GetNext()
	if !objectNode.IsObjectLit() {
		return NewCodingConvention_ObjectLiteralCast(nil, nil, ClosureCodingConvention_OBJECTLIT_EXPECTED)
	}

	return NewCodingConvention_ObjectLiteralCast(typeNode.GetQualifiedName(), typeNode.GetNext(), nil)
}

// Java: ClosureCodingConvention.getAssertionFunctions()
func (c *ClosureCodingConvention) GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec {
	// ImmutableSet.builder(): duplicates (by equals) are dropped, first
	// occurrence order is kept.
	builder := jutil.NewObjectLinkedHashSet[*CodingConvention_AssertionFunctionSpec]()
	for _, spec := range c.CodingConventions_ProxyBase.GetAssertionFunctions() {
		builder.Add(spec)
	}
	for _, spec := range []*CodingConvention_AssertionFunctionSpec{
		CodingConvention_AssertionFunctionSpec_ForTruthy().SetFunctionName("goog.asserts.assert").Build(),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Array"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Boolean"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Element"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Function"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Instanceof"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Number"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("Object"),
		ClosureCodingConvention_CreateGoogAssertOnReturn("String"),
	} {
		builder.Add(spec)
	}
	return builder.ToSlice()
}

// ClosureCodingConvention_CreateGoogAssertOnReturn returns a new assertion
// function goog.asserts.assert[assertedTypeName]
//
// Java: ClosureCodingConvention.createGoogAssertOnReturn(String)
func ClosureCodingConvention_CreateGoogAssertOnReturn(assertedTypeName string) *CodingConvention_AssertionFunctionSpec {
	return CodingConvention_AssertionFunctionSpec_ForMatchesReturn().
		SetFunctionName("goog.asserts.assert" + assertedTypeName).
		Build()
}

var ClosureCodingConvention_GOOG_BIND = rhino.QualifiedName_Of("goog.bind")

var ClosureCodingConvention_GOOG_PARTIAL = rhino.QualifiedName_Of("goog.partial")

// Java: ClosureCodingConvention.describeFunctionBind(Node, boolean)
func (c *ClosureCodingConvention) DescribeFunctionBind(n *rhino.Node, checkTypes bool) *CodingConvention_Bind {
	if !n.IsCall() {
		return nil
	}
	callTarget := n.GetFirstChild()
	if callTarget.IsQualifiedName() {
		if ClosureCodingConvention_GOOG_BIND.Matches(callTarget) || callTarget.MatchesName("goog$bind") {
			// goog.bind(fn, self, args...);
			fn := callTarget.GetNext()
			if fn == nil {
				return nil
			}
			thisValue := ClosureCodingConvention_SafeNext(fn)
			parameters := ClosureCodingConvention_SafeNext(thisValue)
			return NewCodingConvention_Bind(fn, thisValue, parameters)
		}

		if ClosureCodingConvention_GOOG_PARTIAL.Matches(callTarget) || callTarget.MatchesName("goog$partial") {
			// goog.partial(fn, args...);
			fn := callTarget.GetNext()
			if fn == nil {
				return nil
			}
			var thisValue *rhino.Node
			parameters := ClosureCodingConvention_SafeNext(fn)
			return NewCodingConvention_Bind(fn, thisValue, parameters)
		}
	}
	return c.CodingConventions_ProxyBase.DescribeFunctionBind(n, checkTypes)
}

// Java: ClosureCodingConvention.describeCachingCall(Node)
func (c *ClosureCodingConvention) DescribeCachingCall(node *rhino.Node) *CodingConvention_Cache {
	if !node.IsCall() {
		return nil
	}

	callTarget := node.GetFirstChild()
	if c.MatchesCacheMethodName(callTarget) {
		paramCount := node.GetChildCount() - 1
		if 3 <= paramCount && paramCount <= 4 {
			cacheObj := callTarget.GetNext()
			keyNode := cacheObj.GetNext()
			valueFn := keyNode.GetNext()
			keyFn := valueFn.GetNext()

			return NewCodingConvention_Cache(cacheObj, keyNode, valueFn, keyFn)
		}
	}

	return c.CodingConventions_ProxyBase.DescribeCachingCall(node)
}

var ClosureCodingConvention_googCacheReflect = rhino.IR_GetpropWithMoreProps(rhino.IR_Name("goog"), "reflect", "cache")

// Java: ClosureCodingConvention.matchesCacheMethodName(Node)
func (c *ClosureCodingConvention) MatchesCacheMethodName(target *rhino.Node) bool {
	if target.IsGetProp() {
		return target.MatchesQualifiedNameNode(ClosureCodingConvention_googCacheReflect)
	} else if target.IsName() {
		return target.GetString() == "goog$reflect$cache"
	}
	return false
}

// Java: ClosureCodingConvention.safeNext(Node)
func ClosureCodingConvention_SafeNext(n *rhino.Node) *rhino.Node {
	if n != nil {
		return n.GetNext()
	}
	return nil
}
