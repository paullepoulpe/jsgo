/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DestructuringGlobalNameExtractor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.DestructuringGlobalNameExtractor
type DestructuringGlobalNameExtractor struct {
	idHash_ atomic.Int32
}

func (this_ *DestructuringGlobalNameExtractor) BindSelf_(s any) {
	_ = s
}
func (this_ *DestructuringGlobalNameExtractor) IsJavaNil_() bool { return this_ == nil }
func (*DestructuringGlobalNameExtractor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DestructuringGlobalNameExtractor", "DestructuringGlobalNameExtractor"
}
func (this_ *DestructuringGlobalNameExtractor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DestructuringGlobalNameExtractor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DestructuringGlobalNameExtractor", "com.google.javascript.jscomp.DestructuringGlobalNameExtractor", func(v_ any) bool {
		return j2gort.InstanceOf[*DestructuringGlobalNameExtractor](v_)
	})
}

func (this_ *DestructuringGlobalNameExtractor) initFields_DestructuringGlobalNameExtractor() {
	_ = this_
}

// Java constructor: DestructuringGlobalNameExtractor()
func NewDestructuringGlobalNameExtractor() *DestructuringGlobalNameExtractor {
	this_ := &DestructuringGlobalNameExtractor{}
	InitDestructuringGlobalNameExtractor(this_)
	return this_
}

func InitDestructuringGlobalNameExtractor(this_ *DestructuringGlobalNameExtractor) {
	this_.initFields_DestructuringGlobalNameExtractor()
}

// Java static method: DestructuringGlobalNameExtractor.reassignDestructringLvalue(com.google.javascript.rhino.Node stringKey, com.google.javascript.rhino.Node newName, java.util.@org.jspecify.annotations.Nullable Set<com.google.javascript.jscomp.GlobalNamespace.AstChange> newNodes, com.google.javascript.jscomp.GlobalNamespace.Ref ref, com.google.javascript.jscomp.AbstractCompiler compiler)
func DestructuringGlobalNameExtractor_ReassignDestructringLvalue(stringKey pkg_rhino.Node, newName pkg_rhino.Node, newNodes j2gort.Ext_java_util_Set, ref *GlobalNamespace_Ref, compiler AbstractCompiler) {
	_ = stringKey
	_ = newName
	_ = newNodes
	_ = ref
	_ = compiler
	var pattern pkg_rhino.Node = stringKey.GetParent()
	_ = pattern
	var assignmentType pkg_rhino.Node = pattern.GetParent()
	_ = assignmentType
	if !(assignmentType.IsAssign() || assignmentType.IsDestructuringLhs()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(assignmentType))
	}
	var originalRvalue pkg_rhino.Node = pattern.GetNext()
	_ = originalRvalue
	if !(originalRvalue.IsQualifiedName()) {
		j2gort.CheckState(false)
	}
	var lvalueToReassign pkg_rhino.Node
	_ = lvalueToReassign
	if stringKey.GetOnlyChild().IsDefaultValue() {
		lvalueToReassign = stringKey.GetOnlyChild().GetFirstChild()
	} else {
		lvalueToReassign = stringKey.GetOnlyChild()
	}
	if newNodes != nil {
		newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(ref, newName))
	}
	var rvalue pkg_rhino.Node = DestructuringGlobalNameExtractor_MakeNewRvalueForDestructuringKey(stringKey, newName, newNodes, ref)
	_ = rvalue
	if (stringKey.GetPrevious() == 0) && (stringKey.GetNext() != 0) {
		DestructuringGlobalNameExtractor_ReplaceDestructuringAssignment(pattern, lvalueToReassign.Detach(), rvalue)
	} else {
		DestructuringGlobalNameExtractor_AddAfter(pattern, lvalueToReassign.Detach(), rvalue)
	}
	if stringKey.GetNext() != 0 {
		var newPattern pkg_rhino.Node = DestructuringGlobalNameExtractor_CreateNewObjectPatternFromSuccessiveKeys(stringKey).Srcref(pattern)
		_ = newPattern
		var newRvalue pkg_rhino.Node = pkg_rhino.Node(0)
		_ = newRvalue
		if stringKey.GetPrevious() == 0 {
			newRvalue = originalRvalue.Detach()
		} else {
			newRvalue = originalRvalue.CloneTree_0()
			if newNodes != nil {
				newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(ref, newRvalue))
			}
		}
		DestructuringGlobalNameExtractor_AddAfter(lvalueToReassign, newPattern, newRvalue)
	}
	stringKey.Detach()
	compiler.ReportChangeToEnclosingScope(lvalueToReassign)
}

// Java static method: DestructuringGlobalNameExtractor.addAfter(com.google.javascript.rhino.Node originalLvalue, com.google.javascript.rhino.Node newLvalue, com.google.javascript.rhino.Node newRvalue)
func DestructuringGlobalNameExtractor_AddAfter(originalLvalue pkg_rhino.Node, newLvalue pkg_rhino.Node, newRvalue pkg_rhino.Node) {
	_ = originalLvalue
	_ = newLvalue
	_ = newRvalue
	var parent pkg_rhino.Node = originalLvalue.GetParent()
	_ = parent
	if parent.IsAssign() {
		var newAssign pkg_rhino.Node = pkg_rhino.IR_Assign(newLvalue, newRvalue).Srcref(parent)
		_ = newAssign
		var newComma pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_COMMA(), newAssign)
		_ = newComma
		parent.ReplaceWith(newComma)
		newComma.AddChildToFront(parent)
		return
	}
	if newLvalue.IsDestructuringPattern() {
		newLvalue = pkg_rhino.NewNode_3(pkg_rhino.StaticGet_Token_DESTRUCTURING_LHS(), newLvalue, newRvalue).Srcref(parent)
	} else {
		newLvalue.AddChildToBack(newRvalue)
	}
	var declaration pkg_rhino.Node
	_ = declaration
	if parent.IsDestructuringLhs() {
		declaration = originalLvalue.GetGrandparent()
	} else {
		declaration = parent
	}
	if !(NodeUtil_IsNameDeclaration(declaration)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(declaration))
	}
	if NodeUtil_IsStatementParent(declaration.GetParent()) {
		var newDeclaration pkg_rhino.Node = pkg_rhino.NewNode_1(declaration.GetToken()).Srcref(declaration)
		_ = newDeclaration
		newDeclaration.AddChildToBack(newLvalue)
		newDeclaration.InsertAfter(declaration)
	} else {
		declaration.AddChildToBack(newLvalue)
	}
}

// Java static method: DestructuringGlobalNameExtractor.replaceDestructuringAssignment(com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.Node newLvalue, com.google.javascript.rhino.Node newRvalue)
func DestructuringGlobalNameExtractor_ReplaceDestructuringAssignment(pattern pkg_rhino.Node, newLvalue pkg_rhino.Node, newRvalue pkg_rhino.Node) {
	_ = pattern
	_ = newLvalue
	_ = newRvalue
	var parent pkg_rhino.Node = pattern.GetParent()
	_ = parent
	if parent.IsAssign() {
		var newAssign pkg_rhino.Node = pkg_rhino.IR_Assign(newLvalue, newRvalue).Srcref(parent)
		_ = newAssign
		parent.ReplaceWith(newAssign)
	} else if newLvalue.IsName() {
		if !(parent.IsDestructuringLhs()) {
			j2gort.CheckState(false)
		}
		parent.ReplaceWith(newLvalue)
		newLvalue.AddChildToBack(newRvalue)
	} else {
		pattern.GetNext().Detach()
		pattern.Detach()
		parent.AddChildToBack(newLvalue)
		parent.AddChildToBack(newRvalue)
	}
}

// Java static method: DestructuringGlobalNameExtractor.makeNewRvalueForDestructuringKey(com.google.javascript.rhino.Node stringKey, com.google.javascript.rhino.Node rvalue, java.util.Set<com.google.javascript.jscomp.GlobalNamespace.AstChange> newNodes, com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func DestructuringGlobalNameExtractor_MakeNewRvalueForDestructuringKey(stringKey pkg_rhino.Node, rvalue pkg_rhino.Node, newNodes j2gort.Ext_java_util_Set, ref *GlobalNamespace_Ref) pkg_rhino.Node {
	_ = stringKey
	_ = rvalue
	_ = newNodes
	_ = ref
	if stringKey.GetOnlyChild().IsDefaultValue() {
		var defaultValue pkg_rhino.Node = stringKey.GetFirstChild().GetSecondChild().Detach()
		_ = defaultValue
		var rvalueForSheq pkg_rhino.Node = rvalue.CloneTree_0()
		_ = rvalueForSheq
		if newNodes != nil {
			newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(ref, rvalueForSheq))
		}
		rvalue = pkg_rhino.IR_Hook(pkg_rhino.IR_Sheq(NodeUtil_NewUndefinedNode(rvalue), rvalueForSheq), defaultValue, rvalue).SrcrefTree(defaultValue)
	}
	return rvalue
}

// Java static method: DestructuringGlobalNameExtractor.createNewObjectPatternFromSuccessiveKeys(com.google.javascript.rhino.Node stringKey)
func DestructuringGlobalNameExtractor_CreateNewObjectPatternFromSuccessiveKeys(stringKey pkg_rhino.Node) pkg_rhino.Node {
	_ = stringKey
	var newPattern pkg_rhino.Node = stringKey.GetParent().CloneNode_0()
	_ = newPattern
	{
		var next pkg_rhino.Node = stringKey.GetNext()
		_ = next
		for next != 0 {
			var newKey pkg_rhino.Node = next
			_ = newKey
			next = newKey.GetNext()
			newPattern.AddChildToBack(newKey.Detach())
		}
	}
	return newPattern
}
