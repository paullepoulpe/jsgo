/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteNullishCoalesceOperator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteNullishCoalesceOperator
type RewriteNullishCoalesceOperator struct {
	compiler         AbstractCompiler  `gson:"compiler"`
	astFactory       *AstFactory       `gson:"astFactory"`
	uniqueIdSupplier *UniqueIdSupplier `gson:"uniqueIdSupplier"`
	idHash_          atomic.Int32
}

func (this_ *RewriteNullishCoalesceOperator) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteNullishCoalesceOperator) IsJavaNil_() bool { return this_ == nil }
func (*RewriteNullishCoalesceOperator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteNullishCoalesceOperator", "RewriteNullishCoalesceOperator"
}
func (this_ *RewriteNullishCoalesceOperator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteNullishCoalesceOperator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteNullishCoalesceOperator", "com.google.javascript.jscomp.RewriteNullishCoalesceOperator", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteNullishCoalesceOperator](v_)
	})
}
func (*RewriteNullishCoalesceOperator) Implements_jscomp_NodeTraversal_Callback() {}
func (*RewriteNullishCoalesceOperator) Implements_jscomp_CompilerPass()           {}

func StaticGet_RewriteNullishCoalesceOperator_TEMP_VAR_NAME_PREFIX() string {
	return "$jscomp$nullish$tmp"
}
func (this_ *RewriteNullishCoalesceOperator) initFields_RewriteNullishCoalesceOperator() {
	_ = this_
}

// Java constructor: RewriteNullishCoalesceOperator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteNullishCoalesceOperator(compiler AbstractCompiler) *RewriteNullishCoalesceOperator {
	this_ := &RewriteNullishCoalesceOperator{}
	InitRewriteNullishCoalesceOperator(this_, compiler)
	return this_
}

func InitRewriteNullishCoalesceOperator(this_ *RewriteNullishCoalesceOperator, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RewriteNullishCoalesceOperator()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
}

// Java method: RewriteNullishCoalesceOperator.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteNullishCoalesceOperator) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(this_.compiler, root, pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP())
}

// Java method: RewriteNullishCoalesceOperator.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteNullishCoalesceOperator) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return ((scriptFeatures == nil) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP()))
	}
	return true
}

// Java method: RewriteNullishCoalesceOperator.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteNullishCoalesceOperator) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsNullishCoalesce() {
		this_.priv_RewriteNullishCoalesceOperator_visitNullishCoalesce_5018(t, n)
		t.ReportCodeChange_0()
	}
}

// Java method: RewriteNullishCoalesceOperator.visitNullishCoalesce(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *RewriteNullishCoalesceOperator) priv_RewriteNullishCoalesceOperator_visitNullishCoalesce_5018(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var tempVarName string = ("$jscomp$nullish$tmp" + j2gort.StrValString(this_.uniqueIdSupplier.GetUniqueId(t.GetInput())))
	_ = tempVarName
	var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(n)
	_ = enclosingStatement
	var left pkg_rhino.Node = n.RemoveFirstChild()
	_ = left
	var right pkg_rhino.Node = n.GetLastChild().Detach()
	_ = right
	var let pkg_rhino.Node = this_.astFactory.CreateSingleLetNameDeclaration_1(tempVarName)
	_ = let
	var assignName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(tempVarName, AstFactory_S_Type_1_Node(left))
	_ = assignName
	var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(assignName, left)
	_ = assign
	var ne pkg_rhino.Node = this_.astFactory.CreateNe(assign, this_.astFactory.CreateNull())
	_ = ne
	var hookName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(tempVarName, AstFactory_S_Type_1_Node(left))
	_ = hookName
	var hook pkg_rhino.Node = this_.astFactory.CreateHook(ne, hookName, right)
	_ = hook
	let.SrcrefTreeIfMissing(left)
	assignName.SrcrefTreeIfMissing(left)
	assign.SrcrefTreeIfMissing(left)
	ne.SrcrefTreeIfMissing(left)
	hookName.SrcrefTreeIfMissing(left)
	let.InsertBefore(enclosingStatement)
	n.ReplaceWith(hook)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(hook)
}
func (*RewriteNullishCoalesceOperator) CanReleaseSyntacticScopes_() {}
