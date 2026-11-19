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
// Port of com/google/javascript/jscomp/J2clClinitPrunerPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clClinitPrunerPass
type J2clClinitPrunerPass struct {
	emptiedClinitMethods     j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"emptiedClinitMethods"`
	compiler                 AbstractCompiler                          `gson:"compiler"`
	initialChangedScopeNodes j2gort.Ext_java_util_List                 `gson:"initialChangedScopeNodes"`
	idHash_                  atomic.Int32
}

func (this_ *J2clClinitPrunerPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clClinitPrunerPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clClinitPrunerPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clClinitPrunerPass", "J2clClinitPrunerPass"
}
func (this_ *J2clClinitPrunerPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clClinitPrunerPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clClinitPrunerPass", "com.google.javascript.jscomp.J2clClinitPrunerPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clClinitPrunerPass](v_)
	})
}
func (*J2clClinitPrunerPass) Implements_jscomp_CompilerPass() {}

func (this_ *J2clClinitPrunerPass) initFields_J2clClinitPrunerPass() {
	_ = this_
	this_.emptiedClinitMethods = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
}

// Java constructor: J2clClinitPrunerPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.rhino.Node> initialChangedScopeNodes)
func NewJ2clClinitPrunerPass(compiler AbstractCompiler, initialChangedScopeNodes j2gort.Ext_java_util_List) *J2clClinitPrunerPass {
	this_ := &J2clClinitPrunerPass{}
	InitJ2clClinitPrunerPass(this_, compiler, initialChangedScopeNodes)
	return this_
}

func InitJ2clClinitPrunerPass(this_ *J2clClinitPrunerPass, compiler AbstractCompiler, initialChangedScopeNodes j2gort.Ext_java_util_List) {
	_ = compiler
	_ = initialChangedScopeNodes
	this_.initFields_J2clClinitPrunerPass()
	this_.compiler = compiler
	this_.initialChangedScopeNodes = initialChangedScopeNodes
}

// Java method: J2clClinitPrunerPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clClinitPrunerPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler)) {
		return
	}
	this_.priv_J2clClinitPrunerPass_removeRedundantClinits_3450(root, this_.initialChangedScopeNodes)
	this_.priv_J2clClinitPrunerPass_pruneEmptyClinits_3451(root, this_.initialChangedScopeNodes)
	if this_.emptiedClinitMethods.IsEmpty() {
		return
	}
	var clinitReferences j2gort.Ext_com_google_common_collect_Multimap = this_.priv_J2clClinitPrunerPass_collectClinitReferences_3454(root)
	_ = clinitReferences
	for doFirst_1 := true; doFirst_1 || (!(this_.emptiedClinitMethods.IsEmpty())); doFirst_1 = false {
		var newChangedScopes j2gort.Ext_java_util_List = this_.priv_J2clClinitPrunerPass_cleanEmptyClinitReferences_3455(clinitReferences)
		_ = newChangedScopes
		this_.priv_J2clClinitPrunerPass_pruneEmptyClinits_3451(root, newChangedScopes)
	}
	{
		var val_2 PureFunctionIdentifier_Driver
		InitPureFunctionIdentifier_Driver(&val_2, this_.compiler)
		(&val_2).Process_2(externs, root)
	}
}

// Java method: J2clClinitPrunerPass.removeRedundantClinits(com.google.javascript.rhino.Node root, java.util.@org.jspecify.annotations.Nullable List<com.google.javascript.rhino.Node> changedScopeNodes)
func (this_ *J2clClinitPrunerPass) priv_J2clClinitPrunerPass_removeRedundantClinits_3450(root pkg_rhino.Node, changedScopeNodes j2gort.Ext_java_util_List) {
	_ = root
	_ = changedScopeNodes
	var changedRoots j2gort.Ext_java_util_List = J2clClinitPrunerPass_GetNonNestedParentScopeNodes(changedScopeNodes)
	_ = changedRoots
	this_.priv_J2clClinitPrunerPass_traverseScopeRootsWithFallback_3452(root, changedRoots, NewJ2clClinitPrunerPass_RedundantClinitPruner(this_, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](changedRoots)), true)
	this_.priv_J2clClinitPrunerPass_traverseScopeRootsWithFallback_3452(root, changedScopeNodes, NewJ2clClinitPrunerPass_LookAheadRedundantClinitPruner(this_), false)
}

// Java method: J2clClinitPrunerPass.pruneEmptyClinits(com.google.javascript.rhino.Node root, java.util.@org.jspecify.annotations.Nullable List<com.google.javascript.rhino.Node> changedScopes)
func (this_ *J2clClinitPrunerPass) priv_J2clClinitPrunerPass_pruneEmptyClinits_3451(root pkg_rhino.Node, changedScopes j2gort.Ext_java_util_List) {
	_ = root
	_ = changedScopes
	this_.emptiedClinitMethods.Clear()
	this_.priv_J2clClinitPrunerPass_traverseScopeRootsWithFallback_3452(root, changedScopes, NewJ2clClinitPrunerPass_EmptyClinitPruner(this_), false)
	{
		cur_3 := this_.emptiedClinitMethods.IterEntries()
		for cur_3.NextEntry() {
			cur_3.SetValue(this_.priv_J2clClinitPrunerPass_resolveReplacement_3453(cur_3.Value()))
		}
	}
}

// Java method: J2clClinitPrunerPass.traverseScopeRootsWithFallback(com.google.javascript.rhino.Node root, java.util.@org.jspecify.annotations.Nullable List<com.google.javascript.rhino.Node> scopeRoots, com.google.javascript.jscomp.NodeTraversal.Callback callback, boolean traverseNested)
func (this_ *J2clClinitPrunerPass) priv_J2clClinitPrunerPass_traverseScopeRootsWithFallback_3452(root pkg_rhino.Node, scopeRoots j2gort.Ext_java_util_List, callback NodeTraversal_Callback, traverseNested bool) {
	_ = root
	_ = scopeRoots
	_ = callback
	_ = traverseNested
	if scopeRoots == nil {
		NodeTraversal_Traverse(this_.compiler, root, callback)
	} else {
		NodeTraversal_TraverseScopeRoots(this_.compiler, scopeRoots, callback, traverseNested)
	}
}

// Java method: J2clClinitPrunerPass.resolveReplacement(com.google.javascript.rhino.Node node)
func (this_ *J2clClinitPrunerPass) priv_J2clClinitPrunerPass_resolveReplacement_3453(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if node == 0 {
		return pkg_rhino.Node(0)
	}
	var clinitName string = J2clClinitPrunerPass_GetClinitMethodName(node)
	_ = clinitName
	if this_.emptiedClinitMethods.ContainsKey(clinitName) {
		return this_.priv_J2clClinitPrunerPass_resolveReplacement_3453(this_.emptiedClinitMethods.Get(clinitName))
	} else {
		return node
	}
}

// Java method: J2clClinitPrunerPass.collectClinitReferences(com.google.javascript.rhino.Node root)
func (this_ *J2clClinitPrunerPass) priv_J2clClinitPrunerPass_collectClinitReferences_3454(root pkg_rhino.Node) j2gort.Ext_com_google_common_collect_Multimap {
	_ = root
	var clinitReferences j2gort.Ext_com_google_common_collect_Multimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_HashMultimap_Create_0())
	_ = clinitReferences
	NodeTraversal_Traverse(this_.compiler, root, NewJ2clClinitPrunerPass__Anon_54(this_, clinitReferences))
	return clinitReferences
}

// Java method: J2clClinitPrunerPass.cleanEmptyClinitReferences(com.google.common.collect.Multimap<java.lang.String,com.google.javascript.rhino.Node> clinitReferences)
func (this_ *J2clClinitPrunerPass) priv_J2clClinitPrunerPass_cleanEmptyClinitReferences_3455(clinitReferences j2gort.Ext_com_google_common_collect_Multimap) j2gort.Ext_java_util_List {
	_ = clinitReferences
	var newChangedScopes []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = newChangedScopes
	{
		cur_4 := this_.emptiedClinitMethods.IterEntries()
		for cur_4.NextEntry() {
			var clinitName string = cur_4.Key()
			_ = clinitName
			var replacement pkg_rhino.Node = cur_4.Value()
			_ = replacement
			var references j2gort.Ext_java_util_Collection = j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapRemoveAllStr(clinitReferences, clinitName))
			_ = references
			{
				cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(references))
				for cur_5.HasNext() {
					reference := pkg_rhino.Node(cur_5.NextNode())
					_ = reference
					var changedScope pkg_rhino.Node = ChangeTracker_GetEnclosingChangeScopeRoot(reference.GetParent())
					_ = changedScope
					if replacement == 0 {
						NodeUtil_DeleteFunctionCall(reference, this_.compiler)
					} else {
						replacement = replacement.CloneTree_0()
						reference.ReplaceWith(replacement)
						this_.compiler.ReportChangeToChangeScope(changedScope)
						j2gort.MultimapPutStr(clinitReferences, J2clClinitPrunerPass_GetClinitMethodName(replacement), j2gort.BoxNode(replacement))
					}
					newChangedScopes = append(newChangedScopes, changedScope)
				}
			}
		}
	}
	return j2gort.WrapTypedArrayList(newChangedScopes)
}

// Java static method: J2clClinitPrunerPass.getNonNestedParentScopeNodes(java.util.List<com.google.javascript.rhino.Node> changedScopeNodes)
func J2clClinitPrunerPass_GetNonNestedParentScopeNodes(changedScopeNodes j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = changedScopeNodes
	if changedScopeNodes == nil {
		return nil
	} else {
		return NodeUtil_RemoveNestedChangeScopeNodes(NodeUtil_GetParentChangeScopeNodes(changedScopeNodes))
	}
}

// Java static method: J2clClinitPrunerPass.isClinitMethod(com.google.javascript.rhino.Node node)
func J2clClinitPrunerPass_IsClinitMethod(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsFunction() && J2clClinitPrunerPass_IsClinitMethodName(J2clClinitPrunerPass_GetQualifiedNameOfFunction(node)))
}

// Java static method: J2clClinitPrunerPass.getClinitMethodName(com.google.javascript.rhino.Node node)
func J2clClinitPrunerPass_GetClinitMethodName(node pkg_rhino.Node) string {
	_ = node
	if node.IsCall() {
		var fnName string = NodeUtil_GetBestLValueName(node.GetFirstChild())
		_ = fnName
		if J2clClinitPrunerPass_IsClinitMethodName(fnName) {
			return fnName
		} else {
			return j2gort.NullString
		}
	}
	return j2gort.NullString
}

// Java static method: J2clClinitPrunerPass.isClinitMethodName(java.lang.String fnName)
func J2clClinitPrunerPass_IsClinitMethodName(fnName string) bool {
	_ = fnName
	return ((fnName != j2gort.NullString) && (javalang.StringEndsWith(j2gort.CheckNotNullString(fnName), "$$0clinit") || javalang.StringEndsWith(j2gort.CheckNotNullString(fnName), ".$clinit")))
}

// Java static method: J2clClinitPrunerPass.getQualifiedNameOfFunction(com.google.javascript.rhino.Node function)
func J2clClinitPrunerPass_GetQualifiedNameOfFunction(function pkg_rhino.Node) string {
	_ = function
	if !(function.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(function))
	}
	var lValue pkg_rhino.Node = NodeUtil_GetBestLValue(function)
	_ = lValue
	return NodeUtil_GetBestLValueName(lValue)
}

// Java class: com.google.javascript.jscomp.J2clClinitPrunerPass.RedundantClinitPruner
type J2clClinitPrunerPass_RedundantClinitPruner struct {
	outer_                *J2clClinitPrunerPass
	roots                 j2gort.Ext_com_google_common_collect_ImmutableSet   `gson:"roots"`
	stateStack            j2gort.Deque[*J2clClinitPrunerPass_HierarchicalSet] `gson:"stateStack"`
	clinitsCalledAtBranch *J2clClinitPrunerPass_HierarchicalSet               `gson:"clinitsCalledAtBranch"`
	idHash_               atomic.Int32
}

func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) IsJavaNil_() bool { return this_ == nil }
func (*J2clClinitPrunerPass_RedundantClinitPruner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clClinitPrunerPass$RedundantClinitPruner", "RedundantClinitPruner"
}
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clClinitPrunerPass_RedundantClinitPruner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clClinitPrunerPass.RedundantClinitPruner", "com.google.javascript.jscomp.J2clClinitPrunerPass$RedundantClinitPruner", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clClinitPrunerPass_RedundantClinitPruner](v_)
	})
}
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) SetOuter_(o *J2clClinitPrunerPass) {
	this_.outer_ = o
}
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) GetOuter_() *J2clClinitPrunerPass {
	return this_.outer_
}
func (*J2clClinitPrunerPass_RedundantClinitPruner) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) initFields_J2clClinitPrunerPass_RedundantClinitPruner() {
	_ = this_
	this_.stateStack = j2gort.Deque[*J2clClinitPrunerPass_HierarchicalSet]{}
	this_.clinitsCalledAtBranch = NewJ2clClinitPrunerPass_HierarchicalSet(nil)
}

// Java constructor: RedundantClinitPruner(java.lang.Iterable<com.google.javascript.rhino.Node> roots)
func NewJ2clClinitPrunerPass_RedundantClinitPruner(outer_ *J2clClinitPrunerPass, roots j2gort.Ext_java_lang_Iterable) *J2clClinitPrunerPass_RedundantClinitPruner {
	this_ := &J2clClinitPrunerPass_RedundantClinitPruner{}
	this_.outer_ = outer_
	InitJ2clClinitPrunerPass_RedundantClinitPruner(this_, roots)
	return this_
}

func InitJ2clClinitPrunerPass_RedundantClinitPruner(this_ *J2clClinitPrunerPass_RedundantClinitPruner, roots j2gort.Ext_java_lang_Iterable) {
	_ = roots
	this_.initFields_J2clClinitPrunerPass_RedundantClinitPruner()
	if roots == nil {
		this_.roots = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	} else {
		this_.roots = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable(roots)
	}
}

// Java method: RedundantClinitPruner.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) ShouldTraverse_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = node
	_ = parent
	if j2gort.SetContainsNode(this_.roots, node) {
		this_.clinitsCalledAtBranch = NewJ2clClinitPrunerPass_HierarchicalSet(nil)
		this_.stateStack.Clear()
	}
	if NodeUtil_IsFunctionDeclaration(node) || node.IsScript() {
		this_.stateStack.AddLast(this_.clinitsCalledAtBranch)
		this_.clinitsCalledAtBranch = NewJ2clClinitPrunerPass_HierarchicalSet(nil)
	}
	if this_.priv_RedundantClinitPruner_isNewControlBranch_3457(parent) {
		this_.clinitsCalledAtBranch = NewJ2clClinitPrunerPass_HierarchicalSet(this_.clinitsCalledAtBranch)
		if J2clClinitPrunerPass_IsClinitMethod(parent) {
			this_.clinitsCalledAtBranch.Add(j2gort.BoxString(J2clClinitPrunerPass_GetQualifiedNameOfFunction(parent)))
		}
	}
	return true
}

// Java method: RedundantClinitPruner.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	this_.priv_RedundantClinitPruner_tryRemovingClinit_3456(node)
	if this_.priv_RedundantClinitPruner_isNewControlBranch_3457(parent) {
		this_.clinitsCalledAtBranch = this_.clinitsCalledAtBranch.parent
	}
	if NodeUtil_IsFunctionDeclaration(node) || node.IsScript() {
		this_.clinitsCalledAtBranch = this_.stateStack.RemoveLast()
	}
}

// Java method: RedundantClinitPruner.tryRemovingClinit(com.google.javascript.rhino.Node node)
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) priv_RedundantClinitPruner_tryRemovingClinit_3456(node pkg_rhino.Node) {
	_ = node
	var clinitName string = J2clClinitPrunerPass_GetClinitMethodName(node)
	_ = clinitName
	if clinitName == j2gort.NullString {
		return
	}
	if this_.clinitsCalledAtBranch.Add(j2gort.BoxString(clinitName)) {
		return
	}
	NodeUtil_DeleteFunctionCall(node, this_.outer_.compiler)
}

// Java method: RedundantClinitPruner.isNewControlBranch(com.google.javascript.rhino.Node n)
func (this_ *J2clClinitPrunerPass_RedundantClinitPruner) priv_RedundantClinitPruner_isNewControlBranch_3457(n pkg_rhino.Node) bool {
	_ = n
	return ((n != 0) && (((((NodeUtil_IsControlStructure(n) || n.IsDefaultValue()) || n.IsHook()) || n.IsAnd()) || n.IsOr()) || n.IsFunction()))
}

// Java class: com.google.javascript.jscomp.J2clClinitPrunerPass.LookAheadRedundantClinitPruner
type J2clClinitPrunerPass_LookAheadRedundantClinitPruner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *J2clClinitPrunerPass
}

func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clClinitPrunerPass_LookAheadRedundantClinitPruner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clClinitPrunerPass$LookAheadRedundantClinitPruner", "LookAheadRedundantClinitPruner"
}
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clClinitPrunerPass_LookAheadRedundantClinitPruner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clClinitPrunerPass.LookAheadRedundantClinitPruner", "com.google.javascript.jscomp.J2clClinitPrunerPass$LookAheadRedundantClinitPruner", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clClinitPrunerPass_LookAheadRedundantClinitPruner](v_)
	})
}
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) SetOuter_(o *J2clClinitPrunerPass) {
	this_.outer_ = o
}
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) GetOuter_() *J2clClinitPrunerPass {
	return this_.outer_
}
func (*J2clClinitPrunerPass_LookAheadRedundantClinitPruner) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clClinitPrunerPass_LookAheadRedundantClinitPruner) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) initFields_J2clClinitPrunerPass_LookAheadRedundantClinitPruner() {
	_ = this_
}

// Java constructor: LookAheadRedundantClinitPruner()
func NewJ2clClinitPrunerPass_LookAheadRedundantClinitPruner(outer_ *J2clClinitPrunerPass) *J2clClinitPrunerPass_LookAheadRedundantClinitPruner {
	this_ := &J2clClinitPrunerPass_LookAheadRedundantClinitPruner{}
	this_.outer_ = outer_
	InitJ2clClinitPrunerPass_LookAheadRedundantClinitPruner(this_)
	return this_
}

func InitJ2clClinitPrunerPass_LookAheadRedundantClinitPruner(this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clClinitPrunerPass_LookAheadRedundantClinitPruner()
}

// Java method: LookAheadRedundantClinitPruner.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	if !(node.IsExprResult()) {
		return
	}
	var clinitName string = J2clClinitPrunerPass_GetClinitMethodName(node.GetFirstChild())
	_ = clinitName
	if clinitName == j2gort.NullString {
		return
	}
	var callOrNewNode pkg_rhino.Node = this_.priv_LookAheadRedundantClinitPruner_getCallOrNewNode_3460(node.GetNext())
	_ = callOrNewNode
	if (callOrNewNode == 0) || !(callOrNewNode.GetFirstChild().IsName()) {
		return
	}
	var enclosingFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(node)
	_ = enclosingFunction
	if (enclosingFunction == 0) || func() bool {
		sr_6 := callOrNewNode.GetFirstChild().GetString()
		se_7 := NodeUtil_GetNearestFunctionName(enclosingFunction)
		return j2gort.CheckNotNullString(sr_6) == se_7
	}() {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(callOrNewNode.GetFirstChild().GetString()))
	_ = var_
	if ((var_ == nil) || (var_.GetInitialValue() == 0)) || !(var_.GetInitialValue().IsFunction()) {
		return
	}
	var staticFnNode pkg_rhino.Node = var_.GetInitialValue()
	_ = staticFnNode
	if this_.priv_LookAheadRedundantClinitPruner_callsClinit_3461(staticFnNode, clinitName) && this_.priv_LookAheadRedundantClinitPruner_hasSafeArguments_3458(t, callOrNewNode) {
		node.Detach()
		this_.outer_.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java method: LookAheadRedundantClinitPruner.hasSafeArguments(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callOrNewNode)
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) priv_LookAheadRedundantClinitPruner_hasSafeArguments_3458(t *NodeTraversal, callOrNewNode pkg_rhino.Node) bool {
	_ = t
	_ = callOrNewNode
	var child pkg_rhino.Node = callOrNewNode.GetSecondChild()
	_ = child
	for child != 0 {
		if !(NodeUtil_IsLiteralValue(child, false)) && !(this_.priv_LookAheadRedundantClinitPruner_isParameter_3459(t, child)) {
			return false
		}
		child = child.GetNext()
	}
	return true
}

// Java method: LookAheadRedundantClinitPruner.isParameter(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) priv_LookAheadRedundantClinitPruner_isParameter_3459(t *NodeTraversal, n pkg_rhino.Node) bool {
	_ = t
	_ = n
	if !(n.IsName()) {
		return false
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
	_ = var_
	return var_.GetParentNode().IsParamList()
}

// Java method: LookAheadRedundantClinitPruner.getCallOrNewNode(com.google.javascript.rhino.Node n)
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) priv_LookAheadRedundantClinitPruner_getCallOrNewNode_3460(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n == 0 {
		return pkg_rhino.Node(0)
	}
	switch n.TokenSwitchOrd_() {
	case int32(93) /* Token.EXPR_RESULT */, int32(0) /* Token.RETURN */ :
		return this_.priv_LookAheadRedundantClinitPruner_getCallOrNewNode_3460(n.GetFirstChild())
	case int32(28) /* Token.CALL */, int32(23) /* Token.NEW */ :
		return n
	case int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(84) /* Token.VAR */ :
		return func() pkg_rhino.Node {
			if n.HasOneChild() {
				return this_.priv_LookAheadRedundantClinitPruner_getCallOrNewNode_3460(n.GetFirstFirstChild())
			}
			return pkg_rhino.Node(0)
		}()
	default:
		return pkg_rhino.Node(0)
	}
	panic("unreachable switch expression")
}

// Java method: LookAheadRedundantClinitPruner.callsClinit(com.google.javascript.rhino.Node fnNode, java.lang.String clinitName)
func (this_ *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) priv_LookAheadRedundantClinitPruner_callsClinit_3461(fnNode pkg_rhino.Node, clinitName string) bool {
	_ = fnNode
	_ = clinitName
	j2gort.CheckNotNullString(clinitName)
	var child pkg_rhino.Node = fnNode.GetLastChild().GetFirstChild()
	_ = child
	return (((child != 0) && child.IsExprResult()) && func() bool {
		sr_8 := clinitName
		se_9 := J2clClinitPrunerPass_GetClinitMethodName(child.GetFirstChild())
		return j2gort.CheckNotNullString(sr_8) == se_9
	}())
}

// Java class: com.google.javascript.jscomp.J2clClinitPrunerPass.EmptyClinitPruner
type J2clClinitPrunerPass_EmptyClinitPruner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *J2clClinitPrunerPass
}

func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) IsJavaNil_() bool { return this_ == nil }
func (*J2clClinitPrunerPass_EmptyClinitPruner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clClinitPrunerPass$EmptyClinitPruner", "EmptyClinitPruner"
}
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clClinitPrunerPass_EmptyClinitPruner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clClinitPrunerPass.EmptyClinitPruner", "com.google.javascript.jscomp.J2clClinitPrunerPass$EmptyClinitPruner", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clClinitPrunerPass_EmptyClinitPruner](v_)
	})
}
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) SetOuter_(o *J2clClinitPrunerPass) {
	this_.outer_ = o
}
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) GetOuter_() *J2clClinitPrunerPass {
	return this_.outer_
}
func (*J2clClinitPrunerPass_EmptyClinitPruner) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clClinitPrunerPass_EmptyClinitPruner) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) initFields_J2clClinitPrunerPass_EmptyClinitPruner() {
	_ = this_
}

// Java constructor: EmptyClinitPruner()
func NewJ2clClinitPrunerPass_EmptyClinitPruner(outer_ *J2clClinitPrunerPass) *J2clClinitPrunerPass_EmptyClinitPruner {
	this_ := &J2clClinitPrunerPass_EmptyClinitPruner{}
	this_.outer_ = outer_
	InitJ2clClinitPrunerPass_EmptyClinitPruner(this_)
	return this_
}

func InitJ2clClinitPrunerPass_EmptyClinitPruner(this_ *J2clClinitPrunerPass_EmptyClinitPruner) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clClinitPrunerPass_EmptyClinitPruner()
}

// Java method: EmptyClinitPruner.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	if !(J2clClinitPrunerPass_IsClinitMethod(node)) {
		return
	}
	this_.priv_EmptyClinitPruner_trySubstituteEmptyFunction_3462(node)
}

// Java method: EmptyClinitPruner.trySubstituteEmptyFunction(com.google.javascript.rhino.Node fnNode)
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) priv_EmptyClinitPruner_trySubstituteEmptyFunction_3462(fnNode pkg_rhino.Node) {
	_ = fnNode
	var fnQualifiedName string = J2clClinitPrunerPass_GetQualifiedNameOfFunction(fnNode)
	_ = fnQualifiedName
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(fnQualifiedName) {
		return
	}
	var body pkg_rhino.Node = fnNode.GetLastChild()
	_ = body
	if !(body.HasChildren()) {
		return
	}
	var firstExpr pkg_rhino.Node = body.GetFirstChild()
	_ = firstExpr
	if !(this_.priv_EmptyClinitPruner_isAssignToEmptyFn_3463(firstExpr, fnQualifiedName)) {
		return
	}
	var secondExpr pkg_rhino.Node = firstExpr.GetNext()
	_ = secondExpr
	var replaceReferencesWith pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replaceReferencesWith
	if secondExpr == 0 {
		replaceReferencesWith = pkg_rhino.Node(0)
	} else if ((secondExpr.GetNext() == 0) && secondExpr.IsExprResult()) && (J2clClinitPrunerPass_GetClinitMethodName(secondExpr.GetFirstChild()) != j2gort.NullString) {
		replaceReferencesWith = secondExpr.GetFirstChild()
	} else {
		return
	}
	this_.outer_.emptiedClinitMethods.Put(fnQualifiedName, replaceReferencesWith)
	NodeUtil_DeleteNode(firstExpr, this_.outer_.compiler)
}

// Java method: EmptyClinitPruner.isAssignToEmptyFn(com.google.javascript.rhino.Node node, java.lang.String enclosingFnName)
func (this_ *J2clClinitPrunerPass_EmptyClinitPruner) priv_EmptyClinitPruner_isAssignToEmptyFn_3463(node pkg_rhino.Node, enclosingFnName string) bool {
	_ = node
	_ = enclosingFnName
	if !(NodeUtil_IsExprAssign(node)) {
		return false
	}
	var lhs pkg_rhino.Node = node.GetFirstFirstChild()
	_ = lhs
	if lhs.IsGetElem() {
		return false
	}
	if !(lhs.IsName() || lhs.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
	}
	var rhs pkg_rhino.Node = node.GetFirstChild().GetLastChild()
	_ = rhs
	return (NodeUtil_IsEmptyFunctionExpression(rhs) && (NodeUtil_GetBestLValueName(lhs) == enclosingFnName))
}

// Java class: com.google.javascript.jscomp.J2clClinitPrunerPass.HierarchicalSet
type J2clClinitPrunerPass_HierarchicalSet struct {
	currentSet j2gort.Ext_java_util_Set              `gson:"currentSet"`
	parent     *J2clClinitPrunerPass_HierarchicalSet `gson:"parent"`
	idHash_    atomic.Int32
}

func (this_ *J2clClinitPrunerPass_HierarchicalSet) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clClinitPrunerPass_HierarchicalSet) IsJavaNil_() bool { return this_ == nil }
func (*J2clClinitPrunerPass_HierarchicalSet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clClinitPrunerPass$HierarchicalSet", "HierarchicalSet"
}
func (this_ *J2clClinitPrunerPass_HierarchicalSet) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clClinitPrunerPass_HierarchicalSet) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clClinitPrunerPass.HierarchicalSet", "com.google.javascript.jscomp.J2clClinitPrunerPass$HierarchicalSet", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clClinitPrunerPass_HierarchicalSet](v_)
	})
}

func (this_ *J2clClinitPrunerPass_HierarchicalSet) initFields_J2clClinitPrunerPass_HierarchicalSet() {
	_ = this_
	this_.currentSet = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: HierarchicalSet(com.google.javascript.jscomp.J2clClinitPrunerPass.@org.jspecify.annotations.Nullable HierarchicalSet<T> parent)
func NewJ2clClinitPrunerPass_HierarchicalSet(parent *J2clClinitPrunerPass_HierarchicalSet) *J2clClinitPrunerPass_HierarchicalSet {
	this_ := &J2clClinitPrunerPass_HierarchicalSet{}
	InitJ2clClinitPrunerPass_HierarchicalSet(this_, parent)
	return this_
}

func InitJ2clClinitPrunerPass_HierarchicalSet(this_ *J2clClinitPrunerPass_HierarchicalSet, parent *J2clClinitPrunerPass_HierarchicalSet) {
	_ = parent
	this_.initFields_J2clClinitPrunerPass_HierarchicalSet()
	this_.parent = parent
}

// Java method: HierarchicalSet.add(T o)
func (this_ *J2clClinitPrunerPass_HierarchicalSet) Add(o any) bool {
	_ = o
	return (!(this_.priv_HierarchicalSet_parentsContains_3465(o)) && this_.currentSet.Add_1(o))
}

// Java method: HierarchicalSet.parentsContains(T o)
func (this_ *J2clClinitPrunerPass_HierarchicalSet) priv_HierarchicalSet_parentsContains_3465(o any) bool {
	_ = o
	return ((this_.parent != nil) && (this_.parent.currentSet.Contains(o) || this_.parent.priv_HierarchicalSet_parentsContains_3465(o)))
}

// Java class:
type J2clClinitPrunerPass__Anon_54 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                 *J2clClinitPrunerPass
	cap_0_clinitReferences j2gort.Ext_com_google_common_collect_Multimap
}

func (this_ *J2clClinitPrunerPass__Anon_54) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clClinitPrunerPass__Anon_54) IsJavaNil_() bool { return this_ == nil }
func (*J2clClinitPrunerPass__Anon_54) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clClinitPrunerPass$1", ""
}
func (this_ *J2clClinitPrunerPass__Anon_54) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clClinitPrunerPass__Anon_54) HasIdentityEquals_()                                        {}
func (this_ *J2clClinitPrunerPass__Anon_54) SetOuter_(o *J2clClinitPrunerPass)                    { this_.outer_ = o }
func (this_ *J2clClinitPrunerPass__Anon_54) GetOuter_() *J2clClinitPrunerPass                     { return this_.outer_ }
func (*J2clClinitPrunerPass__Anon_54) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*J2clClinitPrunerPass__Anon_54) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *J2clClinitPrunerPass__Anon_54) initFields_J2clClinitPrunerPass__Anon_54() {
	_ = this_
}

func NewJ2clClinitPrunerPass__Anon_54(outer_ *J2clClinitPrunerPass, cap_p_clinitReferences j2gort.Ext_com_google_common_collect_Multimap) *J2clClinitPrunerPass__Anon_54 {
	this_ := &J2clClinitPrunerPass__Anon_54{}
	this_.outer_ = outer_
	this_.cap_0_clinitReferences = cap_p_clinitReferences
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clClinitPrunerPass__Anon_54()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *J2clClinitPrunerPass__Anon_54) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	var clinitName string = J2clClinitPrunerPass_GetClinitMethodName(node)
	_ = clinitName
	if clinitName != j2gort.NullString {
		j2gort.MultimapPutStr(this_.cap_0_clinitReferences, clinitName, j2gort.BoxNode(node))
	}
}
func (*J2clClinitPrunerPass_RedundantClinitPruner) CanReleaseSyntacticScopes_()          {}
func (*J2clClinitPrunerPass_LookAheadRedundantClinitPruner) AlwaysShouldTraverse_()      {}
func (*J2clClinitPrunerPass_LookAheadRedundantClinitPruner) CanReleaseSyntacticScopes_() {}
func (*J2clClinitPrunerPass_EmptyClinitPruner) AlwaysShouldTraverse_()                   {}
func (*J2clClinitPrunerPass_EmptyClinitPruner) CanReleaseSyntacticScopes_()              {}
func (*J2clClinitPrunerPass__Anon_54) AlwaysShouldTraverse_()                            {}
func (*J2clClinitPrunerPass__Anon_54) CanReleaseSyntacticScopes_()                       {}
