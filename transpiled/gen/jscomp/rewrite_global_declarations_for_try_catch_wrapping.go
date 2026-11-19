/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteGlobalDeclarationsForTryCatchWrapping.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteGlobalDeclarationsForTryCatchWrapping
type RewriteGlobalDeclarationsForTryCatchWrapping struct {
	compiler  AbstractCompiler                                  `gson:"compiler"`
	functions j2gort.Ext_com_google_common_collect_ListMultimap `gson:"functions"`
	classes   []pkg_rhino.Node                                  `gson:"classes"`
	idHash_   atomic.Int32
}

func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) IsJavaNil_() bool { return this_ == nil }
func (*RewriteGlobalDeclarationsForTryCatchWrapping) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteGlobalDeclarationsForTryCatchWrapping", "RewriteGlobalDeclarationsForTryCatchWrapping"
}
func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteGlobalDeclarationsForTryCatchWrapping) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteGlobalDeclarationsForTryCatchWrapping", "com.google.javascript.jscomp.RewriteGlobalDeclarationsForTryCatchWrapping", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteGlobalDeclarationsForTryCatchWrapping](v_)
	})
}
func (*RewriteGlobalDeclarationsForTryCatchWrapping) Implements_jscomp_NodeTraversal_Callback() {}
func (*RewriteGlobalDeclarationsForTryCatchWrapping) Implements_jscomp_CompilerPass()           {}

func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) initFields_RewriteGlobalDeclarationsForTryCatchWrapping() {
	_ = this_
	this_.functions = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
	this_.classes = []pkg_rhino.Node{}
}

// Java constructor: RewriteGlobalDeclarationsForTryCatchWrapping(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteGlobalDeclarationsForTryCatchWrapping(compiler AbstractCompiler) *RewriteGlobalDeclarationsForTryCatchWrapping {
	this_ := &RewriteGlobalDeclarationsForTryCatchWrapping{}
	InitRewriteGlobalDeclarationsForTryCatchWrapping(this_, compiler)
	return this_
}

func InitRewriteGlobalDeclarationsForTryCatchWrapping(this_ *RewriteGlobalDeclarationsForTryCatchWrapping, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RewriteGlobalDeclarationsForTryCatchWrapping()
	this_.compiler = compiler
}

// Java method: RewriteGlobalDeclarationsForTryCatchWrapping.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	{
		cur_1 := j2gort.NewMapEntryCursor(j2gort.Ext_com_google_common_collect_Multimaps_S_AsMap_1_Ext_com_google_common_collect_ListMultimap(this_.functions))
		for cur_1.NextEntry() {
			var addingRoot pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(j2gort.CastPtr[JSChunk](cur_1.GetKey()))
			_ = addingRoot
			var fnNodes j2gort.Ext_java_util_List = j2gort.Ext_com_google_common_collect_Lists_Reverse(j2gort.AsIface[j2gort.Ext_java_util_List](cur_1.GetValue()))
			_ = fnNodes
			if !(fnNodes.IsEmpty()) {
				{
					cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(fnNodes))
					for cur_2.HasNext() {
						n := pkg_rhino.Node(cur_2.NextNode())
						_ = n
						var parent pkg_rhino.Node = n.GetParent()
						_ = parent
						n.Detach()
						this_.compiler.ReportChangeToEnclosingScope(parent)
						var nameNode pkg_rhino.Node = n.GetFirstChild()
						_ = nameNode
						var name string = nameNode.GetString()
						_ = name
						nameNode.SetString("")
						addingRoot.AddChildToFront(pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(name), n).SrcrefTreeIfMissing(n))
						this_.compiler.ReportChangeToEnclosingScope(nameNode)
					}
				}
				this_.compiler.ReportChangeToEnclosingScope(addingRoot)
			}
		}
	}
	for _, n_2 := range this_.classes {
		_ = n_2
		var originalNameNode pkg_rhino.Node = n_2.GetFirstChild()
		_ = originalNameNode
		originalNameNode.ReplaceWith(pkg_rhino.IR_Empty())
		var var_ pkg_rhino.Node = pkg_rhino.IR_Var_1(originalNameNode)
		_ = var_
		n_2.ReplaceWith(var_)
		originalNameNode.AddChildToFront(n_2)
		this_.compiler.ReportChangeToEnclosingScope(n_2)
	}

}

// Java method: RewriteGlobalDeclarationsForTryCatchWrapping.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	var grandparent pkg_rhino.Node = n.GetAncestor(int32(2))
	_ = grandparent
	return ((grandparent == 0) || !(grandparent.IsScript()))
}

// Java method: RewriteGlobalDeclarationsForTryCatchWrapping.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteGlobalDeclarationsForTryCatchWrapping) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (parent == 0) || !(parent.IsScript()) {
		return
	}
	if NodeUtil_IsFunctionDeclaration(n) {
		this_.functions.Put(j2gort.PtrToAny(t.GetChunk()), j2gort.BoxNode(n))
	}
	if NodeUtil_IsClassDeclaration(n) {
		this_.classes = append(this_.classes, n)
	}
	if n.IsConst() || n.IsLet() {
		n.SetToken(pkg_rhino.StaticGet_Token_VAR())
		this_.compiler.ReportChangeToEnclosingScope(n)
	}
}
func (*RewriteGlobalDeclarationsForTryCatchWrapping) CanReleaseSyntacticScopes_() {}
