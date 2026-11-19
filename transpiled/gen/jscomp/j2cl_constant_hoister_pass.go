/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clConstantHoisterPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clConstantHoisterPass
type J2clConstantHoisterPass struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *J2clConstantHoisterPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clConstantHoisterPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clConstantHoisterPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clConstantHoisterPass", "J2clConstantHoisterPass"
}
func (this_ *J2clConstantHoisterPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clConstantHoisterPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clConstantHoisterPass", "com.google.javascript.jscomp.J2clConstantHoisterPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clConstantHoisterPass](v_)
	})
}
func (*J2clConstantHoisterPass) Implements_jscomp_CompilerPass() {}

func (this_ *J2clConstantHoisterPass) initFields_J2clConstantHoisterPass() {
	_ = this_
}

// Java constructor: J2clConstantHoisterPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clConstantHoisterPass(compiler AbstractCompiler) *J2clConstantHoisterPass {
	this_ := &J2clConstantHoisterPass{}
	InitJ2clConstantHoisterPass(this_, compiler)
	return this_
}

func InitJ2clConstantHoisterPass(this_ *J2clConstantHoisterPass, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_J2clConstantHoisterPass()
	this_.compiler = compiler
}

// Java method: J2clConstantHoisterPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clConstantHoisterPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler)) {
		return
	}
	var fieldAssignments j2gort.Ext_com_google_common_collect_Multimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
	_ = fieldAssignments
	var hoistableFunctions j2gort.OrderedSet[pkg_rhino.Node] = j2gort.NewOrderedSet[pkg_rhino.Node]()
	_ = hoistableFunctions
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface(func(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
		_ = t
		_ = node
		_ = parent
		if (parent != 0) && NodeUtil_IsLValue(node) {
			j2gort.MultimapPutStr(fieldAssignments, node.GetQualifiedName(), j2gort.BoxNode(parent))
		}
		if J2clConstantHoisterPass_IsHoistableFunction(t, node) {
			hoistableFunctions.Add(node)
		}
	})).Traverse(root)
	{
		cur_1 := j2gort.NewMapValCursor(j2gort.AsIface[j2gort.Ext_java_util_Map](fieldAssignments.AsMap()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			assignments := j2gort.AsIface[j2gort.Ext_java_util_Collection](rawElem_2)
			_ = assignments
			this_.priv_J2clConstantHoisterPass_maybeHoistClassField_3476(assignments, &hoistableFunctions)
		}
	}
}

// Java static method: J2clConstantHoisterPass.isHoistableFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node)
func J2clConstantHoisterPass_IsHoistableFunction(t *NodeTraversal, node pkg_rhino.Node) bool {
	_ = t
	_ = node
	return (node.IsFunction() && (t.GetScope().GetVarCount() == int32(0)))
}

// Java method: J2clConstantHoisterPass.maybeHoistClassField(java.util.Collection<com.google.javascript.rhino.Node> assignments, java.util.Collection<com.google.javascript.rhino.Node> hoistableFunctions)
func (this_ *J2clConstantHoisterPass) priv_J2clConstantHoisterPass_maybeHoistClassField_3476(assignments j2gort.Ext_java_util_Collection, hoistableFunctions *j2gort.OrderedSet[pkg_rhino.Node]) {
	_ = assignments
	_ = hoistableFunctions
	if assignments.Size() != int32(2) {
		return
	}
	var first pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_collect_Iterables_Get_2(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](assignments), int32(0)))
	_ = first
	var second pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.Ext_com_google_common_collect_Iterables_Get_2(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](assignments), int32(1)))
	_ = second
	var topLevelDeclaration pkg_rhino.Node
	_ = topLevelDeclaration
	if J2clConstantHoisterPass_IsClassFieldDeclaration(first) {
		topLevelDeclaration = first
	} else {
		topLevelDeclaration = second
	}
	var clinitAssignment pkg_rhino.Node
	_ = clinitAssignment
	if J2clConstantHoisterPass_IsClinitFieldAssignment(first) {
		clinitAssignment = first
	} else {
		clinitAssignment = second
	}
	if !(J2clConstantHoisterPass_IsClassFieldDeclaration(topLevelDeclaration)) || !(J2clConstantHoisterPass_IsClinitFieldAssignment(clinitAssignment)) {
		return
	}
	var assignmentRhs pkg_rhino.Node = clinitAssignment.GetSecondChild()
	_ = assignmentRhs
	if !(NodeUtil_IsLiteralValue(assignmentRhs, true)) || (assignmentRhs.IsFunction() && !(j2gort.CheckNotNullPtr(hoistableFunctions).Contains(assignmentRhs))) {
		return
	}
	if NodeUtil_GetEnclosingScript(clinitAssignment) != NodeUtil_GetEnclosingScript(topLevelDeclaration) {
		return
	}
	this_.priv_J2clConstantHoisterPass_hoistConstantLikeField_3477(clinitAssignment, topLevelDeclaration)
}

// Java method: J2clConstantHoisterPass.hoistConstantLikeField(com.google.javascript.rhino.Node clinitAssignment, com.google.javascript.rhino.Node topLevelDeclaration)
func (this_ *J2clConstantHoisterPass) priv_J2clConstantHoisterPass_hoistConstantLikeField_3477(clinitAssignment pkg_rhino.Node, topLevelDeclaration pkg_rhino.Node) {
	_ = clinitAssignment
	_ = topLevelDeclaration
	var clinitAssignedValue pkg_rhino.Node = clinitAssignment.GetSecondChild()
	_ = clinitAssignedValue
	var declarationInClass pkg_rhino.Node = topLevelDeclaration.GetFirstChild()
	_ = declarationInClass
	var declarationAssignedValue pkg_rhino.Node = declarationInClass.GetFirstChild()
	_ = declarationAssignedValue
	var clinitChangeScope pkg_rhino.Node = ChangeTracker_GetEnclosingChangeScopeRoot(clinitAssignment)
	_ = clinitChangeScope
	NodeUtil_RemoveChild(clinitAssignment.GetParent(), clinitAssignment)
	clinitAssignedValue.Detach()
	this_.compiler.ReportChangeToChangeScope(clinitChangeScope)
	if declarationAssignedValue == 0 {
		declarationInClass.AddChildToFront(clinitAssignedValue)
		this_.compiler.ReportChangeToEnclosingScope(topLevelDeclaration)
	} else if !(j2gort.CheckNotNullNode(declarationAssignedValue) == clinitAssignedValue) {
		if !(NodeUtil_IsLiteralValue(declarationAssignedValue, false)) {
			j2gort.CheckState(false)
		}
		declarationAssignedValue.ReplaceWith(clinitAssignedValue)
		this_.compiler.ReportChangeToEnclosingScope(topLevelDeclaration)
	}
	declarationInClass.SetDeclaredConstantVar(true)
}

// Java static method: J2clConstantHoisterPass.isClassFieldDeclaration(com.google.javascript.rhino.Node node)
func J2clConstantHoisterPass_IsClassFieldDeclaration(node pkg_rhino.Node) bool {
	_ = node
	return ((node.GetParent().IsScript() && node.IsVar()) && (!(node.GetFirstChild().HasChildren()) || ((node.GetFirstFirstChild() != 0) && NodeUtil_IsLiteralValue(node.GetFirstFirstChild(), false))))
}

// Java static method: J2clConstantHoisterPass.isClinitFieldAssignment(com.google.javascript.rhino.Node node)
func J2clConstantHoisterPass_IsClinitFieldAssignment(node pkg_rhino.Node) bool {
	_ = node
	return ((node.GetParent().IsExprResult() && node.GetGrandparent().IsBlock()) && J2clConstantHoisterPass_IsClinitMethod(node.GetGrandparent().GetParent()))
}

// Java static method: J2clConstantHoisterPass.isClinitMethod(com.google.javascript.rhino.Node fnNode)
func J2clConstantHoisterPass_IsClinitMethod(fnNode pkg_rhino.Node) bool {
	_ = fnNode
	if !(fnNode.IsFunction()) {
		return false
	}
	var fnName string = NodeUtil_GetName(fnNode)
	_ = fnName
	return ((fnName != j2gort.NullString) && J2clConstantHoisterPass_IsClinitMethodName(fnName))
}

// Java static method: J2clConstantHoisterPass.isClinitMethodName(java.lang.String methodName)
func J2clConstantHoisterPass_IsClinitMethodName(methodName string) bool {
	_ = methodName
	return ((methodName != j2gort.NullString) && (javalang.StringEndsWith(j2gort.CheckNotNullString(methodName), "$$0clinit") || javalang.StringEndsWith(j2gort.CheckNotNullString(methodName), ".$clinit")))
}
