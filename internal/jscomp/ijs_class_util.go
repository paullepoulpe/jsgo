// Port of com/google/javascript/jscomp/ijs/ClassUtil.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// ClassUtil holds static utility methods for dealing with classes. The
// primary benefit is for papering over the differences between ES6 class
// syntax.
//
// Go-only: language-adapter — the Java class has only a private
// constructor and static methods; it has no Go type.

// ClassUtil_IsThisPropInsideClassWithName returns whether the given node
// represents a GETPROP with a first child THIS inside a named class.
//
// Java: ClassUtil.isThisPropInsideClassWithName(Node)
func ClassUtil_IsThisPropInsideClassWithName(maybeGetProp *rhino.Node) bool {
	return classUtil_getClassNameOfThisProp(maybeGetProp) != nil
}

// ClassUtil_GetFullyQualifiedNameOfThisProp returns the fully qualified name
// of a this.property inside a constructor. This method called should only
// be called if `isThisPropInsideClassWithName` returns true.
//
// Java: ClassUtil.getFullyQualifiedNameOfThisProp(Node)
func ClassUtil_GetFullyQualifiedNameOfThisProp(getProp *rhino.Node) string {
	jutil.CheckArgument(ClassUtil_IsThisPropInsideClassWithName(getProp))
	className := classUtil_getClassNameOfThisProp(getProp)
	return *className + ".prototype." + getProp.GetString()
}

// Java: ClassUtil.getClassNameOfThisProp(Node)
func classUtil_getClassNameOfThisProp(getprop *rhino.Node) *string {
	if !getprop.IsGetProp() || !getprop.GetFirstChild().IsThis() {
		return nil
	}
	function := NodeUtil_GetEnclosingFunction(getprop)
	if function == nil {
		return nil
	}
	className := classUtil_getMemberFunctionClassName(function)
	if className == nil || *className == "" {
		return nil
	}
	return className
}

// ClassUtil_IsMemberFieldDefInsideClassWithName returns whether the given
// node represents a MEMBER_FIELD_DEF that is inside a class with a name.
//
// Java: ClassUtil.isMemberFieldDefInsideClassWithName(Node)
func ClassUtil_IsMemberFieldDefInsideClassWithName(fieldNode *rhino.Node) bool {
	return classUtil_getMemberFieldDefClassName(fieldNode) != nil
}

// ClassUtil_IsComputedMemberInsideClassWithName returns whether the given
// node represents a MEMBER_FIELD_DEF that is inside a class with a name.
//
// Java: ClassUtil.isComputedMemberInsideClassWithName(Node)
func ClassUtil_IsComputedMemberInsideClassWithName(fieldNode *rhino.Node) bool {
	return fieldNode.GetParent().IsClassMembers() &&
		fieldNode.GetFirstChild().IsGetProp() &&
		fieldNode.GetFirstFirstChild().MatchesName("Symbol")
}

// ClassUtil_GetFullyQualifiedNameOfMemberFieldDef returns the fully
// qualified name of a MEMBER_FIELD_DEF. It is invalid to call this method
// for a field that belongs to a nameless class.
//
// Java: ClassUtil.getFullyQualifiedNameOfMemberFieldDef(Node)
func ClassUtil_GetFullyQualifiedNameOfMemberFieldDef(fieldNode *rhino.Node) string {
	jutil.CheckArgument(ClassUtil_IsMemberFieldDefInsideClassWithName(fieldNode))
	className := classUtil_getMemberFieldDefClassName(fieldNode)
	if fieldNode.IsStaticMember() {
		return *className + "." + fieldNode.GetString()
	}
	return *className + ".prototype." + fieldNode.GetString()
}

// Java: ClassUtil.getMemberFieldDefClassName(Node)
func classUtil_getMemberFieldDefClassName(fieldNode *rhino.Node) *string {
	jutil.CheckArgument(fieldNode.IsMemberFieldDef())
	classNode := fieldNode.GetGrandparent()
	jutil.CheckState(classNode.IsClass())
	return NodeUtil_GetName(classNode)
}

// Java: ClassUtil.getFullyQualifiedNameOfMethod(Node)
func ClassUtil_GetFullyQualifiedNameOfMethod(function *rhino.Node) string {
	jutil.CheckArgument(ClassUtil_IsClassMethod(function))
	className := classUtil_getMemberFunctionClassName(function)
	jutil.CheckState(className != nil && *className != "")
	memberFunctionDef := function.GetParent()
	methodName := memberFunctionDef.GetString()
	if memberFunctionDef.IsStaticMember() {
		return *className + "." + methodName
	}
	return *className + ".prototype." + methodName
}

// Java: ClassUtil.isClassMethod(Node)
func ClassUtil_IsClassMethod(functionNode *rhino.Node) bool {
	jutil.CheckArgument(functionNode.IsFunction())
	parent := functionNode.GetParent()
	return parent.IsMemberFunctionDef() && parent.GetParent().IsClassMembers()
}

// ClassUtil_HasNamedClass checks whether the given constructor/member
// function belongs to a named class, as opposed to an anonymous class.
//
// Java: ClassUtil.hasNamedClass(Node)
func ClassUtil_HasNamedClass(functionNode *rhino.Node) bool {
	jutil.CheckArgument(functionNode.IsFunction())
	return classUtil_getMemberFunctionClassName(functionNode) != nil
}

// Java: ClassUtil.getMemberFunctionClassName(Node)
func classUtil_getMemberFunctionClassName(functionNode *rhino.Node) *string {
	jutil.CheckArgument(functionNode.IsFunction())
	if ClassUtil_IsClassMethod(functionNode) {
		parent := functionNode.GetParent()
		jutil.CheckState(parent.IsMemberFunctionDef())
		// ES6 class
		classNode := functionNode.GetGrandparent().GetParent()
		jutil.CheckState(classNode.IsClass())
		return NodeUtil_GetName(classNode)
	}
	return NodeUtil_GetName(functionNode)
}

// Java: ClassUtil.isConstructor(Node)
func ClassUtil_IsConstructor(functionNode *rhino.Node) bool {
	if ClassUtil_IsClassMethod(functionNode) {
		return NodeUtil_IsEs6Constructor(functionNode)
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(functionNode)
	return jsdoc != nil && jsdoc.IsConstructorOrInterface()
}
