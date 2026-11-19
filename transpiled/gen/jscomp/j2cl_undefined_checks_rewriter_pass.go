/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clUndefinedChecksRewriterPass.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clUndefinedChecksRewriterPass
type J2clUndefinedChecksRewriterPass struct {
	AbstractPeepholeOptimizationBase
	shouldRunJ2clPasses bool `gson:"shouldRunJ2clPasses"`
}

func (this_ *J2clUndefinedChecksRewriterPass) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *J2clUndefinedChecksRewriterPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clUndefinedChecksRewriterPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clUndefinedChecksRewriterPass", "J2clUndefinedChecksRewriterPass"
}
func (this_ *J2clUndefinedChecksRewriterPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*J2clUndefinedChecksRewriterPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clUndefinedChecksRewriterPass", "com.google.javascript.jscomp.J2clUndefinedChecksRewriterPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clUndefinedChecksRewriterPass](v_)
	})
}
func (*J2clUndefinedChecksRewriterPass) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *J2clUndefinedChecksRewriterPass) initFields_J2clUndefinedChecksRewriterPass() {
	_ = this_

}

// Java constructor: J2clUndefinedChecksRewriterPass()
func NewJ2clUndefinedChecksRewriterPass() *J2clUndefinedChecksRewriterPass {
	this_ := &J2clUndefinedChecksRewriterPass{}
	this_.BindSelf_(this_)
	InitJ2clUndefinedChecksRewriterPass(this_)
	return this_
}

func InitJ2clUndefinedChecksRewriterPass(this_ *J2clUndefinedChecksRewriterPass) {
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_J2clUndefinedChecksRewriterPass()
}

// Java method: J2clUndefinedChecksRewriterPass.beginTraversal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *J2clUndefinedChecksRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	_ = compiler
	this_.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	this_.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java method: J2clUndefinedChecksRewriterPass.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *J2clUndefinedChecksRewriterPass) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	if !(this_.shouldRunJ2clPasses) {
		return subtree
	}
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	if J2clUndefinedChecksRewriterPass_IsUtilCall(subtree, "$isUndefined") && subtree.HasXChildren(int32(2)) {
		replacement = J2clUndefinedChecksRewriterPass_MaybeOptimizeIsUndefinedCall(subtree)
	} else if J2clUndefinedChecksRewriterPass_IsUtilCall(subtree, "$coerceToNull") && subtree.HasXChildren(int32(2)) {
		replacement = J2clUndefinedChecksRewriterPass_MaybeOptimizeCoerceToNullCall(subtree)
	}
	if replacement != 0 {
		if !(replacement != subtree) {
			j2gort.CheckState(false)
		}
		replacement.SrcrefTreeIfMissing(subtree)
		subtree.ReplaceWith(replacement)
		this_.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return subtree
}

// Java static method: J2clUndefinedChecksRewriterPass.maybeOptimizeIsUndefinedCall(com.google.javascript.rhino.Node node)
func J2clUndefinedChecksRewriterPass_MaybeOptimizeIsUndefinedCall(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	var valueExpression pkg_rhino.Node = j2gort.CheckNotNullNode(node.GetSecondChild())
	_ = valueExpression
	if NodeUtil_IsUndefined(valueExpression) {
		return pkg_rhino.IR_TrueNode()
	} else if NodeUtil_IsDefinedValue(valueExpression) {
		return pkg_rhino.IR_Comma(valueExpression.Detach(), pkg_rhino.IR_FalseNode())
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java static method: J2clUndefinedChecksRewriterPass.maybeOptimizeCoerceToNullCall(com.google.javascript.rhino.Node node)
func J2clUndefinedChecksRewriterPass_MaybeOptimizeCoerceToNullCall(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	var valueExpression pkg_rhino.Node = j2gort.CheckNotNullNode(node.GetSecondChild())
	_ = valueExpression
	if NodeUtil_IsUndefined(valueExpression) {
		return pkg_rhino.IR_NullNode()
	}
	if NodeUtil_IsDefinedValue(valueExpression) {
		return valueExpression.Detach()
	}
	return pkg_rhino.Node(0)
}

// Java static method: J2clUndefinedChecksRewriterPass.isUtilCall(com.google.javascript.rhino.Node node, java.lang.String methodName)
func J2clUndefinedChecksRewriterPass_IsUtilCall(node pkg_rhino.Node, methodName string) bool {
	_ = node
	_ = methodName
	if !(node.IsCall()) || !(node.GetFirstChild().IsQualifiedName()) {
		return false
	}
	var qualifiedName string = node.GetFirstChild().GetOriginalQualifiedName()
	_ = qualifiedName
	return ((j2gort.JStrIndexOf(qualifiedName, "nativebootstrap$Util") >= 0) && javalang.StringEndsWith(j2gort.CheckNotNullString(qualifiedName), j2gort.CheckNotNullString(methodName)))
}
