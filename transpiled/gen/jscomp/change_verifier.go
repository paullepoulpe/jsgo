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
// Port of com/google/javascript/jscomp/ChangeVerifier.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_change_verifier_1 j2gort.Ext_java_util_function_Supplier

func init() {
	lambdaSingleton_change_verifier_1 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return nil
	})
}

// Java class: com.google.javascript.jscomp.ChangeVerifier
type ChangeVerifier struct {
	changeTracker   *ChangeTracker                             `gson:"changeTracker"`
	clonesByCurrent j2gort.Ext_com_google_common_collect_BiMap `gson:"clonesByCurrent"`
	snapshotChange  int32                                      `gson:"snapshotChange"`
	idHash_         atomic.Int32
}

func (this_ *ChangeVerifier) BindSelf_(s any) {
	_ = s
}
func (this_ *ChangeVerifier) IsJavaNil_() bool { return this_ == nil }
func (*ChangeVerifier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChangeVerifier", "ChangeVerifier"
}
func (this_ *ChangeVerifier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ChangeVerifier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ChangeVerifier", "com.google.javascript.jscomp.ChangeVerifier", func(v_ any) bool {
		return j2gort.InstanceOf[*ChangeVerifier](v_)
	})
}

func (this_ *ChangeVerifier) initFields_ChangeVerifier() {
	_ = this_
	this_.clonesByCurrent = j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](j2gort.Ext_com_google_common_collect_HashBiMap_Create_0())
}

// Java constructor: ChangeVerifier(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewChangeVerifier(compiler AbstractCompiler) *ChangeVerifier {
	this_ := &ChangeVerifier{}
	InitChangeVerifier(this_, compiler)
	return this_
}

func InitChangeVerifier(this_ *ChangeVerifier, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ChangeVerifier()
	this_.changeTracker = compiler.GetChangeTracker()
}

// Java method: ChangeVerifier.snapshot(com.google.javascript.rhino.Node root)
func (this_ *ChangeVerifier) Snapshot(root pkg_rhino.Node) *ChangeVerifier {
	_ = root
	this_.clonesByCurrent.Clear()
	this_.snapshotChange = this_.changeTracker.GetChangeStamp()
	var snapshot pkg_rhino.Node = root.CloneTree_0()
	_ = snapshot
	this_.priv_ChangeVerifier_associateClones_978(root, snapshot)
	return this_
}

// Java method: ChangeVerifier.checkRecordedChanges(com.google.javascript.rhino.Node current)
func (this_ *ChangeVerifier) CheckRecordedChanges_1(current pkg_rhino.Node) {
	_ = current
	this_.CheckRecordedChanges_2("", current)
}

// Java method: ChangeVerifier.checkRecordedChanges(java.lang.String passName, com.google.javascript.rhino.Node root)
func (this_ *ChangeVerifier) CheckRecordedChanges_2(passName string, root pkg_rhino.Node) {
	_ = passName
	_ = root
	this_.priv_ChangeVerifier_verifyScopeChangesHaveBeenRecorded_979(passName, root)
}

// Java method: ChangeVerifier.associateClones(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node snapshot)
func (this_ *ChangeVerifier) priv_ChangeVerifier_associateClones_978(n pkg_rhino.Node, snapshot pkg_rhino.Node) {
	_ = n
	_ = snapshot
	if n.IsRoot() || ChangeTracker_IsChangeScopeRoot(n) {
		j2gort.MapPutNode(this_.clonesByCurrent, n, j2gort.BoxNode(snapshot))
	}
	var child pkg_rhino.Node = n.GetFirstChild()
	_ = child
	var snapshotChild pkg_rhino.Node = snapshot.GetFirstChild()
	_ = snapshotChild
	for child != 0 {
		this_.priv_ChangeVerifier_associateClones_978(child, snapshotChild)
		child = child.GetNext()
		snapshotChild = snapshotChild.GetNext()
	}
}

// Java method: ChangeVerifier.verifyScopeChangesHaveBeenRecorded(java.lang.String passName, com.google.javascript.rhino.Node root)
func (this_ *ChangeVerifier) priv_ChangeVerifier_verifyScopeChangesHaveBeenRecorded_979(passName string, root pkg_rhino.Node) {
	_ = passName
	_ = root
	var passNameMsg string
	_ = passNameMsg
	if j2gort.CheckNotNullString(passName) == "" {
		passNameMsg = ""
	} else {
		passNameMsg = (j2gort.StrValString(passName) + ": ")
	}
	var snapshotScopeNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = snapshotScopeNodes
	NodeUtil_VisitPreOrder_2(j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetNode(this_.clonesByCurrent, root)), NewChangeVerifier__Anon_3(this_, snapshotScopeNodes))
	NodeUtil_VisitPreOrder_2(root, NewChangeVerifier__Anon_4(this_, snapshotScopeNodes, passNameMsg))
	this_.priv_ChangeVerifier_verifyDeletedScopeNodes_980(passNameMsg, snapshotScopeNodes)
}

// Java method: ChangeVerifier.verifyDeletedScopeNodes(java.lang.String passNameMsg, java.util.Set<com.google.javascript.rhino.Node> deletedScopeNodes)
func (this_ *ChangeVerifier) priv_ChangeVerifier_verifyDeletedScopeNodes_980(passNameMsg string, deletedScopeNodes j2gort.Ext_java_util_Set) {
	_ = passNameMsg
	_ = deletedScopeNodes
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(deletedScopeNodes))
		for cur_1.HasNext() {
			snapshotScopeNode := pkg_rhino.Node(cur_1.NextNode())
			_ = snapshotScopeNode
			var currentNode pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetNode(j2gort.AsIface[j2gort.Ext_com_google_common_collect_BiMap](this_.clonesByCurrent.Inverse()), snapshotScopeNode))
			_ = currentNode
			if currentNode.IsDeleted() {
				continue
			}
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString((j2gort.StrValString(passNameMsg) + "deleted scope was not reported:\n" + j2gort.StrValString(currentNode.ToStringTree()))))))
		}
	}
}

// Java method: ChangeVerifier.verifyNode(java.lang.String passNameMsg, com.google.javascript.rhino.Node n)
func (this_ *ChangeVerifier) priv_ChangeVerifier_verifyNode_981(passNameMsg string, n pkg_rhino.Node) {
	_ = passNameMsg
	_ = n
	if n.IsDeleted() {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString((j2gort.StrValString(passNameMsg) + "existing scope is improperly marked as deleted:\n" + j2gort.StrValString(n.ToStringTree()))))))
	}
}

// Java method: ChangeVerifier.verifyNewNode(java.lang.String passNameMsg, com.google.javascript.rhino.Node n)
func (this_ *ChangeVerifier) priv_ChangeVerifier_verifyNewNode_982(passNameMsg string, n pkg_rhino.Node) {
	_ = passNameMsg
	_ = n
	var changeTime int32 = n.GetChangeTime()
	_ = changeTime
	if (changeTime == int32(0)) || (changeTime < this_.snapshotChange) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString((j2gort.StrValString(passNameMsg) + "new scope not explicitly marked as changed:\n" + j2gort.StrValString(n.ToStringTree()))))))
	}
}

// Java method: ChangeVerifier.verifyRoot(com.google.javascript.rhino.Node root)
func (this_ *ChangeVerifier) priv_ChangeVerifier_verifyRoot_983(root pkg_rhino.Node) {
	_ = root
	if !(root.IsRoot()) {
		j2gort.CheckState(false)
	}
	if root.GetChangeTime() != int32(0) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Root nodes should never be marked as changed.")))
	}
}

// Java method: ChangeVerifier.verifyNodeChange(java.lang.String passNameMsg, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node snapshot)
func (this_ *ChangeVerifier) priv_ChangeVerifier_verifyNodeChange_984(passNameMsg string, n pkg_rhino.Node, snapshot pkg_rhino.Node) {
	_ = passNameMsg
	_ = n
	_ = snapshot
	if n.IsRoot() {
		return
	}
	var result *ChangeVerifier_EqualsResult = ChangeVerifier_GetInequivalenceReasonExcludingFunctions(n, snapshot)
	_ = result
	if n.GetChangeTime() > snapshot.GetChangeTime() {
		if result.M_Equals_0() {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString((j2gort.StrValString(passNameMsg) + "unchanged scope marked as changed: " + j2gort.StrValString(this_.GetNameForNode(n)))))))
		}
	} else {
		if !(result.M_Equals_0()) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("\"%schanged scope not marked as changed: %s.\n%s\nAncestor nodes:\n%s\n", []any{j2gort.BoxString(passNameMsg), j2gort.BoxString(this_.GetNameForNode(n)), j2gort.BoxString(j2gort.CastString(result.errorMessage.Get())), j2gort.BoxString(this_.priv_ChangeVerifier_path_986(n, result.ErrorNode()))})))))
		}
	}
}

// Java method: ChangeVerifier.getNameForNode(com.google.javascript.rhino.Node n)
func (this_ *ChangeVerifier) GetNameForNode(n pkg_rhino.Node) string {
	_ = n
	var sourceName string = NodeUtil_GetSourceName(n)
	_ = sourceName
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		return ("SCRIPT: " + j2gort.StrValString(sourceName))
	case int32(73) /* Token.FUNCTION */ :
		var fnName string = NodeUtil_GetNearestFunctionName(n)
		_ = fnName
		if fnName == j2gort.NullString {
			fnName = ("anonymous@" + j2gort.StrValInt32(n.GetLineno()) + ":" + j2gort.StrValInt32(n.GetCharno()))
		}
		return ("FUNCTION: " + j2gort.StrValString(fnName) + " in " + j2gort.StrValString(sourceName))
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected Node type")))
	}
	panic("unreachable switch expression")
}

// Java method: ChangeVerifier.path(com.google.javascript.rhino.Node ancestor, com.google.javascript.rhino.Node child)
func (this_ *ChangeVerifier) priv_ChangeVerifier_path_986(ancestor pkg_rhino.Node, child pkg_rhino.Node) string {
	_ = ancestor
	_ = child
	var childToAncestor j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = childToAncestor
	{
		var current pkg_rhino.Node = child
		_ = current
		for ; current != ancestor; current = current.GetParent() {
			childToAncestor.Add_1(j2gort.BoxString(current.String()))
		}
	}
	childToAncestor.Add_1(j2gort.BoxString(ancestor.String()))
	j2gort.Ext_java_util_Collections_Reverse(childToAncestor)
	var result *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		for ; i < childToAncestor.Size(); i = i + 1 {
			result.Repeat_2_int32_int32(int32(uint16(32)), (i * int32(2)))
			result.Append_1_string(j2gort.ListGetStr(childToAncestor, i))
			result.Append_1_string("\n")
		}
	}
	return result.String()
}

// Java static method: ChangeVerifier.getInequivalenceReasonExcludingFunctions(com.google.javascript.rhino.Node thisNode, com.google.javascript.rhino.Node thatNode)
func ChangeVerifier_GetInequivalenceReasonExcludingFunctions(thisNode pkg_rhino.Node, thatNode pkg_rhino.Node) *ChangeVerifier_EqualsResult {
	_ = thisNode
	_ = thatNode
	j2gort.CheckNotNullNode(thisNode)
	j2gort.CheckNotNullNode(thatNode)
	if thisNode.GetChildCount() != thatNode.GetChildCount() {
		return ChangeVerifier_EqualsResult_NotEqual(thisNode, "differing child count", j2gort.BoxString((j2gort.StrVal(thisNode) + ": " + j2gort.StrValInt32(thisNode.GetChildCount()))), j2gort.BoxString((j2gort.StrVal(thatNode) + ": " + j2gort.StrValInt32(thatNode.GetChildCount()))))
	}
	if !(thisNode.IsEquivalentWithSideEffectsToShallow(thatNode)) {
		return ChangeVerifier_EqualsResult_NotEqual(thisNode, "shallow inequivalence", j2gort.BoxNode(thisNode), j2gort.BoxNode(thatNode))
	}
	if thisNode.IsFunction() && thatNode.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(thisNode) != NodeUtil_IsFunctionDeclaration(thatNode) {
			return ChangeVerifier_EqualsResult_NotEqual(thisNode, "mismatched isFunctionDeclaration", j2gort.BoxNode(thisNode), j2gort.BoxNode(thatNode))
		}
	}
	var thisChild pkg_rhino.Node = thisNode.GetFirstChild()
	_ = thisChild
	var thatChild pkg_rhino.Node = thatNode.GetFirstChild()
	_ = thatChild
	for (thisChild != 0) && (thatChild != 0) {
		if thisChild.IsFunction() || thisChild.IsScript() {
			if thatChild.RawToken_() != thisChild.RawToken_() {
				return ChangeVerifier_EqualsResult_NotEqual(thisNode, "different tokens", j2gort.BoxNode(thisChild), j2gort.BoxNode(thatChild))
			}
			if thisChild.IsFunction() && NodeUtil_IsFunctionDeclaration(thisChild) {
				var thisName string = thisChild.GetFirstChild().GetString()
				_ = thisName
				var thatName string = thatChild.GetFirstChild().GetString()
				_ = thatName
				if !(j2gort.CheckNotNullString(thisName) == thatName) {
					return ChangeVerifier_EqualsResult_NotEqual(thisNode, "function name changed", j2gort.BoxString(thisName), j2gort.BoxString(thatName))
				}
			}
		} else {
			var result *ChangeVerifier_EqualsResult = ChangeVerifier_GetInequivalenceReasonExcludingFunctions(thisChild, thatChild)
			_ = result
			if !(result.M_Equals_0()) {
				return result
			}
		}
		thisChild = thisChild.GetNext()
		thatChild = thatChild.GetNext()
	}
	return ChangeVerifier_EqualsResult_Equal()
}

// Java record: com.google.javascript.jscomp.ChangeVerifier.EqualsResult
type ChangeVerifier_EqualsResult struct {
	errorMessage j2gort.Ext_java_util_function_Supplier `gson:"errorMessage"`
	errorNode    pkg_rhino.Node                         `gson:"errorNode"`
	equals       bool                                   `gson:"equals"`
}

func (this_ *ChangeVerifier_EqualsResult) BindSelf_(s any) {
	_ = s
}
func (this_ *ChangeVerifier_EqualsResult) IsJavaNil_() bool { return this_ == nil }
func (*ChangeVerifier_EqualsResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChangeVerifier$EqualsResult", "EqualsResult"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ChangeVerifier.EqualsResult", "com.google.javascript.jscomp.ChangeVerifier$EqualsResult", func(v_ any) bool {
		return j2gort.InstanceOf[*ChangeVerifier_EqualsResult](v_)
	})
}
func (*ChangeVerifier_EqualsResult) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ChangeVerifier_EqualsResult) initFields_ChangeVerifier_EqualsResult() {
	_ = this_
}

// Java constructor: EqualsResult(boolean equals, com.google.javascript.rhino.Node errorNode, java.util.function.Supplier<java.lang.String> errorMessage)
func NewChangeVerifier_EqualsResult(equals bool, errorNode pkg_rhino.Node, errorMessage j2gort.Ext_java_util_function_Supplier) *ChangeVerifier_EqualsResult {
	this_ := &ChangeVerifier_EqualsResult{}
	InitChangeVerifier_EqualsResult(this_, equals, errorNode, errorMessage)
	return this_
}

func InitChangeVerifier_EqualsResult(this_ *ChangeVerifier_EqualsResult, equals bool, errorNode pkg_rhino.Node, errorMessage j2gort.Ext_java_util_function_Supplier) {
	_ = equals
	_ = errorNode
	_ = errorMessage
	this_.initFields_ChangeVerifier_EqualsResult()
	this_.equals = equals
	this_.errorNode = errorNode
	this_.errorMessage = errorMessage
}

// Java static method: EqualsResult.equal()
func ChangeVerifier_EqualsResult_Equal() *ChangeVerifier_EqualsResult {
	return NewChangeVerifier_EqualsResult(true, pkg_rhino.Node(0), lambdaSingleton_change_verifier_1)
}

// Java static method: EqualsResult.notEqual(com.google.javascript.rhino.Node errorNode, java.lang.String error, java.lang.Object after, java.lang.Object before)
func ChangeVerifier_EqualsResult_NotEqual(errorNode pkg_rhino.Node, error_ string, after any, before any) *ChangeVerifier_EqualsResult {
	_ = errorNode
	_ = error_
	_ = after
	_ = before
	return NewChangeVerifier_EqualsResult(false, errorNode, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s\nBefore: %s\nAfter:  %s\n", []any{j2gort.BoxString(error_), before, after}))
	}))
}

func (this_ *ChangeVerifier_EqualsResult) M_Equals_0() bool          { return this_.equals }
func (this_ *ChangeVerifier_EqualsResult) ErrorNode() pkg_rhino.Node { return this_.errorNode }
func (this_ *ChangeVerifier_EqualsResult) ErrorMessage() j2gort.Ext_java_util_function_Supplier {
	return this_.errorMessage
}
func (this_ *ChangeVerifier_EqualsResult) Equals(other_ any) bool {
	o_, ok_ := other_.(*ChangeVerifier_EqualsResult)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.equals != o_.equals {
		return false
	}
	if this_.errorNode != o_.errorNode {
		return false
	}
	if !javalang.ObjectsEquals(this_.errorMessage, o_.errorMessage) {
		return false
	}
	return true
}
func (this_ *ChangeVerifier_EqualsResult) HashCode() int32 {
	var h_ int32 = 0
	if this_.equals {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	h_ = 31*h_ + int32(this_.errorNode)
	h_ = 31*h_ + j2gort.ObjectHash(this_.errorMessage)
	return h_
}
func (this_ *ChangeVerifier_EqualsResult) String() string {
	_ = fmt.Sprint
	return "EqualsResult[" + "equals=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.equals)) + ", " + "errorNode=" + j2gort.StrVal(j2gort.BoxNode(this_.errorNode)) + ", " + "errorMessage=" + j2gort.StrVal(this_.errorMessage) + "]"
}

// Java class:
type ChangeVerifier__Anon_3 struct {
	outer_                   *ChangeVerifier
	cap_0_snapshotScopeNodes j2gort.Ext_java_util_Set
	idHash_                  atomic.Int32
}

func (this_ *ChangeVerifier__Anon_3) BindSelf_(s any) {
	_ = s
}
func (this_ *ChangeVerifier__Anon_3) IsJavaNil_() bool { return this_ == nil }
func (*ChangeVerifier__Anon_3) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChangeVerifier$1", ""
}
func (this_ *ChangeVerifier__Anon_3) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ChangeVerifier__Anon_3) HasIdentityEquals_()                 {}
func (this_ *ChangeVerifier__Anon_3) SetOuter_(o *ChangeVerifier)   { this_.outer_ = o }
func (this_ *ChangeVerifier__Anon_3) GetOuter_() *ChangeVerifier    { return this_.outer_ }
func (*ChangeVerifier__Anon_3) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *ChangeVerifier__Anon_3) initFields_ChangeVerifier__Anon_3() {
	_ = this_
}

func NewChangeVerifier__Anon_3(outer_ *ChangeVerifier, cap_p_snapshotScopeNodes j2gort.Ext_java_util_Set) *ChangeVerifier__Anon_3 {
	this_ := &ChangeVerifier__Anon_3{}
	this_.outer_ = outer_
	this_.cap_0_snapshotScopeNodes = cap_p_snapshotScopeNodes
	this_.initFields_ChangeVerifier__Anon_3()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node oldNode)
func (this_ *ChangeVerifier__Anon_3) Visit(oldNode pkg_rhino.Node) {
	_ = oldNode
	if ChangeTracker_IsChangeScopeRoot(oldNode) {
		j2gort.SetAddNode(this_.cap_0_snapshotScopeNodes, oldNode)
	}
}

// Java class:
type ChangeVerifier__Anon_4 struct {
	outer_                   *ChangeVerifier
	cap_0_snapshotScopeNodes j2gort.Ext_java_util_Set
	cap_1_passNameMsg        string
	idHash_                  atomic.Int32
}

func (this_ *ChangeVerifier__Anon_4) BindSelf_(s any) {
	_ = s
}
func (this_ *ChangeVerifier__Anon_4) IsJavaNil_() bool { return this_ == nil }
func (*ChangeVerifier__Anon_4) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ChangeVerifier$2", ""
}
func (this_ *ChangeVerifier__Anon_4) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ChangeVerifier__Anon_4) HasIdentityEquals_()                 {}
func (this_ *ChangeVerifier__Anon_4) SetOuter_(o *ChangeVerifier)   { this_.outer_ = o }
func (this_ *ChangeVerifier__Anon_4) GetOuter_() *ChangeVerifier    { return this_.outer_ }
func (*ChangeVerifier__Anon_4) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *ChangeVerifier__Anon_4) initFields_ChangeVerifier__Anon_4() {
	_ = this_
}

func NewChangeVerifier__Anon_4(outer_ *ChangeVerifier, cap_p_snapshotScopeNodes j2gort.Ext_java_util_Set, cap_p_passNameMsg string) *ChangeVerifier__Anon_4 {
	this_ := &ChangeVerifier__Anon_4{}
	this_.outer_ = outer_
	this_.cap_0_snapshotScopeNodes = cap_p_snapshotScopeNodes
	this_.cap_1_passNameMsg = cap_p_passNameMsg
	this_.initFields_ChangeVerifier__Anon_4()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node n)
func (this_ *ChangeVerifier__Anon_4) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsRoot() {
		this_.outer_.priv_ChangeVerifier_verifyRoot_983(n)
	} else if ChangeTracker_IsChangeScopeRoot(n) {
		var clone pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetNode(this_.outer_.clonesByCurrent, n))
		_ = clone
		j2gort.SetRemoveNode(this_.cap_0_snapshotScopeNodes, clone)
		this_.outer_.priv_ChangeVerifier_verifyNode_981(this_.cap_1_passNameMsg, n)
		if clone == 0 {
			this_.outer_.priv_ChangeVerifier_verifyNewNode_982(this_.cap_1_passNameMsg, n)
		} else {
			this_.outer_.priv_ChangeVerifier_verifyNodeChange_984(this_.cap_1_passNameMsg, n, clone)
		}
	}
}
