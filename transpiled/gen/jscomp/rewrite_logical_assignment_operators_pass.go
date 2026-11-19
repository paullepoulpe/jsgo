/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsPass
type RewriteLogicalAssignmentOperatorsPass struct {
	compiler                                AbstractCompiler                         `gson:"compiler"`
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper `gson:"rewriteLogicalAssignmentOperatorsHelper"`
	idHash_                                 atomic.Int32
}

func (this_ *RewriteLogicalAssignmentOperatorsPass) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteLogicalAssignmentOperatorsPass) IsJavaNil_() bool { return this_ == nil }
func (*RewriteLogicalAssignmentOperatorsPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsPass", "RewriteLogicalAssignmentOperatorsPass"
}
func (this_ *RewriteLogicalAssignmentOperatorsPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteLogicalAssignmentOperatorsPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsPass", "com.google.javascript.jscomp.RewriteLogicalAssignmentOperatorsPass", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteLogicalAssignmentOperatorsPass](v_)
	})
}
func (*RewriteLogicalAssignmentOperatorsPass) Implements_jscomp_NodeTraversal_Callback() {}
func (*RewriteLogicalAssignmentOperatorsPass) Implements_jscomp_CompilerPass()           {}

func (this_ *RewriteLogicalAssignmentOperatorsPass) initFields_RewriteLogicalAssignmentOperatorsPass() {
	_ = this_
}

// Java constructor: RewriteLogicalAssignmentOperatorsPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteLogicalAssignmentOperatorsPass(compiler AbstractCompiler) *RewriteLogicalAssignmentOperatorsPass {
	this_ := &RewriteLogicalAssignmentOperatorsPass{}
	InitRewriteLogicalAssignmentOperatorsPass(this_, compiler)
	return this_
}

func InitRewriteLogicalAssignmentOperatorsPass(this_ *RewriteLogicalAssignmentOperatorsPass, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RewriteLogicalAssignmentOperatorsPass()
	this_.compiler = compiler
	this_.rewriteLogicalAssignmentOperatorsHelper = NewRewriteLogicalAssignmentOperatorsHelper(compiler, compiler.CreateAstFactory(), compiler.GetUniqueIdSupplier())
}

// Java method: RewriteLogicalAssignmentOperatorsPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteLogicalAssignmentOperatorsPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(this_.compiler, root, pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT())
}

// Java method: RewriteLogicalAssignmentOperatorsPass.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteLogicalAssignmentOperatorsPass) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return ((scriptFeatures == nil) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT()))
	}
	return true
}

// Java method: RewriteLogicalAssignmentOperatorsPass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node logicalAssignment, com.google.javascript.rhino.Node parent)
func (this_ *RewriteLogicalAssignmentOperatorsPass) Visit_3(t *NodeTraversal, logicalAssignment pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = logicalAssignment
	_ = parent
	switch logicalAssignment.TokenSwitchOrd_() {
	case int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */ :
		this_.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(t, logicalAssignment)
	default:
	}
}
func (*RewriteLogicalAssignmentOperatorsPass) CanReleaseSyntacticScopes_() {}
