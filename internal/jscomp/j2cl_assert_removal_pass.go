// Port of com/google/javascript/jscomp/J2clAssertRemovalPass.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clAssertRemovalPass is an optimization pass to remove J2CL
// Asserts.$assert.
//
// Java: public class J2clAssertRemovalPass extends AbstractPostOrderCallback implements CompilerPass
type J2clAssertRemovalPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: J2clAssertRemovalPass(AbstractCompiler)
func NewJ2clAssertRemovalPass(compiler AbstractCompiler) *J2clAssertRemovalPass {
	return &J2clAssertRemovalPass{compiler: compiler}
}

// Java: J2clAssertRemovalPass.process(Node, Node)
func (p *J2clAssertRemovalPass) Process(externs *rhino.Node, root *rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: J2clAssertRemovalPass.visit(NodeTraversal, Node, Node)
func (p *J2clAssertRemovalPass) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	if j2clAssertRemovalPass_isAssertionCall(node) {
		NodeUtil_DeleteFunctionCall(node, p.compiler)
	}
}

// Java: J2clAssertRemovalPass.isAssertionCall(Node)
func j2clAssertRemovalPass_isAssertionCall(node *rhino.Node) bool {
	return node.IsCall() && j2clAssertRemovalPass_isAssertionMethodName(node.GetFirstChild())
}

// Java: J2clAssertRemovalPass.isAssertionMethodName(Node)
func j2clAssertRemovalPass_isAssertionMethodName(fnName *rhino.Node) bool {
	if !fnName.IsQualifiedName() {
		return false
	}
	originalQname := *fnName.GetOriginalQualifiedName()
	return strings.Contains(originalQname, "Asserts") && strings.Contains(originalQname, ".$assert")
}
