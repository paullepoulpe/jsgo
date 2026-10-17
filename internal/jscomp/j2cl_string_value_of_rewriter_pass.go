// Port of com/google/javascript/jscomp/J2clStringValueOfRewriterPass.java @ closure-compiler v20260415

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const j2clStringValueOfRewriterPass_METHOD_NAME = "module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String"

// J2clStringValueOfRewriterPass is an optimization pass to rewrite J2CL
// String.valueOf.
//
// Java: final class J2clStringValueOfRewriterPass extends AbstractPeepholeOptimization
type J2clStringValueOfRewriterPass struct {
	AbstractPeepholeOptimizationBase
	shouldRunJ2clPasses bool
}

// Java: J2clStringValueOfRewriterPass()
func NewJ2clStringValueOfRewriterPass() *J2clStringValueOfRewriterPass {
	p := &J2clStringValueOfRewriterPass{}
	p.InitAbstractPeepholeOptimizationBase(p)
	p.shouldRunJ2clPasses = false
	return p
}

// Java: J2clStringValueOfRewriterPass.beginTraversal(AbstractCompiler)
func (p *J2clStringValueOfRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	p.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java: J2clStringValueOfRewriterPass.optimizeSubtree(Node)
func (p *J2clStringValueOfRewriterPass) OptimizeSubtree(node *rhino.Node) *rhino.Node {
	if !p.shouldRunJ2clPasses || !j2clStringValueOfRewriterPass_isStringValueOfCall(node) {
		return node
	}

	replacement := p.tryRewriteStringValueOfCall(node)
	if replacement != node {
		replacement = replacement.SrcrefIfMissing(node)
		node.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
	}
	return replacement
}

// Cases to rewrite:
//
//	String.valueOf("foo") -> String("foo")
//	String.valueOf(null) -> String("null")
//	String.valueOf(undefined) -> "null"
//	String.valueOf("foo" + bar) -> String("foo" + bar)
//	String.valueOf(bar + "foo") -> String(bar + "foo")
//	String.valueOf("foo" + String.valueOf(bar)) -> String("foo" + String.valueOf(bar))
//
// Java: J2clStringValueOfRewriterPass.tryRewriteStringValueOfCall(Node)
func (p *J2clStringValueOfRewriterPass) tryRewriteStringValueOfCall(n *rhino.Node) *rhino.Node {
	param := n.GetSecondChild()
	if NodeUtil_IsUndefined(param) {
		return rhino.IR_String("null")
	} else if NodeUtil_IsDefinedValue(param) && !param.IsArrayLit() {
		// Generate String(param), let other peephole optimizations handle the rest when safe
		return NodeUtil_NewCallNode(rhino.IR_Name("String").Srcref(n), param.Detach())
	}
	return n
}

// Java: J2clStringValueOfRewriterPass.isStringValueOfCall(Node)
func j2clStringValueOfRewriterPass_isStringValueOfCall(node *rhino.Node) bool {
	return node.IsCall() &&
		// Do not optimize if the parameter was removed
		node.HasXChildren(2) &&
		j2clStringValueOfRewriterPass_isStringValueOfMethodName(node.GetFirstChild())
}

// Java: J2clStringValueOfRewriterPass.isStringValueOfMethodName(Node)
func j2clStringValueOfRewriterPass_isStringValueOfMethodName(fnName *rhino.Node) bool {
	return fnName.IsQualifiedName() && *fnName.GetOriginalQualifiedName() == j2clStringValueOfRewriterPass_METHOD_NAME
}
