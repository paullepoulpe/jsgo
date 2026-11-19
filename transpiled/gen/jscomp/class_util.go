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
// Port of com/google/javascript/jscomp/ijs/ClassUtil.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ijs.ClassUtil
type ClassUtil struct {
	idHash_ atomic.Int32
}

func (this_ *ClassUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *ClassUtil) IsJavaNil_() bool { return this_ == nil }
func (*ClassUtil) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ijs.ClassUtil", "ClassUtil" }
func (this_ *ClassUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClassUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ClassUtil", "com.google.javascript.jscomp.ijs.ClassUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*ClassUtil](v_)
	})
}

func (this_ *ClassUtil) initFields_ClassUtil() {
	_ = this_
}

// Java constructor: ClassUtil()
func NewClassUtil() *ClassUtil {
	this_ := &ClassUtil{}
	InitClassUtil(this_)
	return this_
}

func InitClassUtil(this_ *ClassUtil) {
	this_.initFields_ClassUtil()
}

// Java static method: ClassUtil.isThisPropInsideClassWithName(com.google.javascript.rhino.Node maybeGetProp)
func ClassUtil_IsThisPropInsideClassWithName(maybeGetProp pkg_rhino.Node) bool {
	_ = maybeGetProp
	return (ClassUtil_GetClassNameOfThisProp(maybeGetProp) != j2gort.NullString)
}

// Java static method: ClassUtil.getFullyQualifiedNameOfThisProp(com.google.javascript.rhino.Node getProp)
func ClassUtil_GetFullyQualifiedNameOfThisProp(getProp pkg_rhino.Node) string {
	_ = getProp
	if !(ClassUtil_IsThisPropInsideClassWithName(getProp)) { j2gort.CheckArgument(false) }
	var className string = ClassUtil_GetClassNameOfThisProp(getProp)
	_ = className
	return (j2gort.StrValString(className) + ".prototype." + j2gort.StrValString(getProp.GetString()))
}

// Java static method: ClassUtil.getClassNameOfThisProp(com.google.javascript.rhino.Node getprop)
func ClassUtil_GetClassNameOfThisProp(getprop pkg_rhino.Node) string {
	_ = getprop
	if ((!(getprop.IsGetProp()) || !(getprop.GetFirstChild().IsThis()))) {
		return j2gort.NullString
	}
	if (getprop.IsPrivateIdentifier()) {
		return j2gort.NullString
	}
	var function pkg_rhino.Node = NodeUtil_GetEnclosingFunction(getprop)
	_ = function
	if ((function == 0)) {
		return j2gort.NullString
	}
	var className string = ClassUtil_GetMemberFunctionClassName(function)
	_ = className
	if (j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(className)) {
		return j2gort.NullString
	}
	return className
}

// Java static method: ClassUtil.isMemberFieldDefInsideClassWithName(com.google.javascript.rhino.Node fieldNode)
func ClassUtil_IsMemberFieldDefInsideClassWithName(fieldNode pkg_rhino.Node) bool {
	_ = fieldNode
	return (ClassUtil_GetMemberFieldDefClassName(fieldNode) != j2gort.NullString)
}

// Java static method: ClassUtil.isComputedMemberInsideClassWithName(com.google.javascript.rhino.Node fieldNode)
func ClassUtil_IsComputedMemberInsideClassWithName(fieldNode pkg_rhino.Node) bool {
	_ = fieldNode
	return ((fieldNode.GetParent().IsClassMembers() && fieldNode.GetFirstChild().IsGetProp()) && fieldNode.GetFirstFirstChild().MatchesName_1_string("Symbol"))
}

// Java static method: ClassUtil.getFullyQualifiedNameOfMemberFieldDef(com.google.javascript.rhino.Node fieldNode)
func ClassUtil_GetFullyQualifiedNameOfMemberFieldDef(fieldNode pkg_rhino.Node) string {
	_ = fieldNode
	if !(ClassUtil_IsMemberFieldDefInsideClassWithName(fieldNode)) { j2gort.CheckArgument(false) }
	var className string = ClassUtil_GetMemberFieldDefClassName(fieldNode)
	_ = className
	if fieldNode.IsStaticMember() {
		return (j2gort.StrValString(className) + "." + j2gort.StrValString(fieldNode.GetString()))
	} else {
		return (j2gort.StrValString(className) + ".prototype." + j2gort.StrValString(fieldNode.GetString()))
	}
}

// Java static method: ClassUtil.getMemberFieldDefClassName(com.google.javascript.rhino.Node fieldNode)
func ClassUtil_GetMemberFieldDefClassName(fieldNode pkg_rhino.Node) string {
	_ = fieldNode
	if !(fieldNode.IsMemberFieldDef()) { j2gort.CheckArgument(false) }
	var classNode pkg_rhino.Node = fieldNode.GetGrandparent()
	_ = classNode
	if !(classNode.IsClass()) { j2gort.CheckState(false) }
	return NodeUtil_GetName(classNode)
}

// Java static method: ClassUtil.getFullyQualifiedNameOfMethod(com.google.javascript.rhino.Node function)
func ClassUtil_GetFullyQualifiedNameOfMethod(function pkg_rhino.Node) string {
	_ = function
	if !(ClassUtil_IsClassMethod(function)) { j2gort.CheckArgument(false) }
	var className string = ClassUtil_GetMemberFunctionClassName(function)
	_ = className
	if !(((className != j2gort.NullString) && !((j2gort.CheckNotNullString(className) == "")))) { j2gort.CheckState(false) }
	var memberFunctionDef pkg_rhino.Node = function.GetParent()
	_ = memberFunctionDef
	var methodName string = memberFunctionDef.GetString()
	_ = methodName
	if memberFunctionDef.IsStaticMember() {
		return (j2gort.StrValString(className) + "." + j2gort.StrValString(methodName))
	} else {
		return (j2gort.StrValString(className) + ".prototype." + j2gort.StrValString(methodName))
	}
}

// Java static method: ClassUtil.isClassMethod(com.google.javascript.rhino.Node functionNode)
func ClassUtil_IsClassMethod(functionNode pkg_rhino.Node) bool {
	_ = functionNode
	if !(functionNode.IsFunction()) { j2gort.CheckArgument(false) }
	var parent pkg_rhino.Node = functionNode.GetParent()
	_ = parent
	return ((parent.IsMemberFunctionDef() && parent.GetParent().IsClassMembers()))
}

// Java static method: ClassUtil.hasNamedClass(com.google.javascript.rhino.Node functionNode)
func ClassUtil_HasNamedClass(functionNode pkg_rhino.Node) bool {
	_ = functionNode
	if !(functionNode.IsFunction()) { j2gort.CheckArgument(false) }
	return (ClassUtil_GetMemberFunctionClassName(functionNode) != j2gort.NullString)
}

// Java static method: ClassUtil.getMemberFunctionClassName(com.google.javascript.rhino.Node functionNode)
func ClassUtil_GetMemberFunctionClassName(functionNode pkg_rhino.Node) string {
	_ = functionNode
	if !(functionNode.IsFunction()) { j2gort.CheckArgument(false) }
	if (ClassUtil_IsClassMethod(functionNode)) {
		var parent pkg_rhino.Node = functionNode.GetParent()
		_ = parent
		if !(parent.IsMemberFunctionDef()) { j2gort.CheckState(false) }
		var classNode pkg_rhino.Node = functionNode.GetGrandparent().GetParent()
		_ = classNode
		if !(classNode.IsClass()) { j2gort.CheckState(false) }
		return NodeUtil_GetName(classNode)
	}
	return NodeUtil_GetName(functionNode)
}

// Java static method: ClassUtil.isConstructor(com.google.javascript.rhino.Node functionNode)
func ClassUtil_IsConstructor(functionNode pkg_rhino.Node) bool {
	_ = functionNode
	if (ClassUtil_IsClassMethod(functionNode)) {
		return NodeUtil_IsEs6Constructor(functionNode)
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(functionNode)
	_ = jsdoc
	return ((jsdoc != nil) && jsdoc.IsConstructorOrInterface())
}

