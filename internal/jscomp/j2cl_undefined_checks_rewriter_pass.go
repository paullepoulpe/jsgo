// Port of com/google/javascript/jscomp/J2clUndefinedChecksRewriterPass.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: final class J2clUndefinedChecksRewriterPass extends AbstractPeepholeOptimization
type J2clUndefinedChecksRewriterPass struct {
	AbstractPeepholeOptimizationBase
	shouldRunJ2clPasses bool
}

// Java: J2clUndefinedChecksRewriterPass()
func NewJ2clUndefinedChecksRewriterPass() *J2clUndefinedChecksRewriterPass {
	p := &J2clUndefinedChecksRewriterPass{}
	p.InitAbstractPeepholeOptimizationBase(p)
	p.shouldRunJ2clPasses = false
	return p
}

// Java: J2clUndefinedChecksRewriterPass.beginTraversal(AbstractCompiler)
func (p *J2clUndefinedChecksRewriterPass) BeginTraversal(compiler AbstractCompiler) {
	p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
	p.shouldRunJ2clPasses = J2clSourceFileChecker_ShouldRunJ2clPasses(compiler)
}

// Java: J2clUndefinedChecksRewriterPass.optimizeSubtree(Node)
func (p *J2clUndefinedChecksRewriterPass) OptimizeSubtree(subtree *rhino.Node) *rhino.Node {
	if !p.shouldRunJ2clPasses {
		return subtree
	}

	var replacement *rhino.Node
	if j2clUndefinedChecksRewriterPass_isUtilCall(subtree, "$isUndefined") && subtree.HasXChildren(2) {
		replacement = j2clUndefinedChecksRewriterPass_maybeOptimizeIsUndefinedCall(subtree)
	} else if j2clUndefinedChecksRewriterPass_isUtilCall(subtree, "$coerceToNull") && subtree.HasXChildren(2) {
		replacement = j2clUndefinedChecksRewriterPass_maybeOptimizeCoerceToNullCall(subtree)
	}

	if replacement != nil {
		jutil.CheckState(replacement != subtree)
		replacement.SrcrefTreeIfMissing(subtree)
		subtree.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}

	return subtree
}

// Java: J2clUndefinedChecksRewriterPass.maybeOptimizeIsUndefinedCall(Node)
func j2clUndefinedChecksRewriterPass_maybeOptimizeIsUndefinedCall(node *rhino.Node) *rhino.Node {
	valueExpression := jutil.CheckNotNull(node.GetSecondChild())
	if NodeUtil_IsUndefined(valueExpression) {
		return rhino.IR_TrueNode()
	} else if NodeUtil_IsDefinedValue(valueExpression) {
		return rhino.IR_Comma(valueExpression.Detach(), rhino.IR_FalseNode())
	} else {
		return nil
	}
}

// Java: J2clUndefinedChecksRewriterPass.maybeOptimizeCoerceToNullCall(Node)
func j2clUndefinedChecksRewriterPass_maybeOptimizeCoerceToNullCall(node *rhino.Node) *rhino.Node {
	valueExpression := jutil.CheckNotNull(node.GetSecondChild())
	if NodeUtil_IsUndefined(valueExpression) {
		return rhino.IR_NullNode()
	}

	if NodeUtil_IsDefinedValue(valueExpression) {
		// detach the valueExpression as we'll be using it as the replacement node.
		return valueExpression.Detach()
	}

	return nil
}

// Java: J2clUndefinedChecksRewriterPass.isUtilCall(Node, String)
func j2clUndefinedChecksRewriterPass_isUtilCall(node *rhino.Node, methodName string) bool {
	if !node.IsCall() || !node.GetFirstChild().IsQualifiedName() {
		return false
	}
	qualifiedName := *node.GetFirstChild().GetOriginalQualifiedName()
	return strings.Contains(qualifiedName, "nativebootstrap$Util") && strings.HasSuffix(qualifiedName, methodName)
}
