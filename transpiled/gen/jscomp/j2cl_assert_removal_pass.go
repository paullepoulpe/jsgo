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
// Port of com/google/javascript/jscomp/J2clAssertRemovalPass.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.J2clAssertRemovalPass
type J2clAssertRemovalPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *J2clAssertRemovalPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clAssertRemovalPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clAssertRemovalPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clAssertRemovalPass", "J2clAssertRemovalPass"
}
func (this_ *J2clAssertRemovalPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clAssertRemovalPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clAssertRemovalPass", "com.google.javascript.jscomp.J2clAssertRemovalPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clAssertRemovalPass](v_)
	})
}
func (*J2clAssertRemovalPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*J2clAssertRemovalPass) Implements_jscomp_CompilerPass()                            {}
func (*J2clAssertRemovalPass) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *J2clAssertRemovalPass) initFields_J2clAssertRemovalPass() {
	_ = this_
}

// Java constructor: J2clAssertRemovalPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clAssertRemovalPass(compiler AbstractCompiler) *J2clAssertRemovalPass {
	this_ := &J2clAssertRemovalPass{}
	InitJ2clAssertRemovalPass(this_, compiler)
	return this_
}

func InitJ2clAssertRemovalPass(this_ *J2clAssertRemovalPass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clAssertRemovalPass()
	this_.compiler = compiler
}

// Java method: J2clAssertRemovalPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clAssertRemovalPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler)) {
		return
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: J2clAssertRemovalPass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *J2clAssertRemovalPass) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	if J2clAssertRemovalPass_IsAssertionCall(node) {
		NodeUtil_DeleteFunctionCall(node, this_.compiler)
	}
}

// Java static method: J2clAssertRemovalPass.isAssertionCall(com.google.javascript.rhino.Node node)
func J2clAssertRemovalPass_IsAssertionCall(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsCall() && J2clAssertRemovalPass_IsAssertionMethodName(node.GetFirstChild()))
}

// Java static method: J2clAssertRemovalPass.isAssertionMethodName(com.google.javascript.rhino.Node fnName)
func J2clAssertRemovalPass_IsAssertionMethodName(fnName pkg_rhino.Node) bool {
	_ = fnName
	if !(fnName.IsQualifiedName()) {
		return false
	}
	var originalQname string = fnName.GetOriginalQualifiedName()
	_ = originalQname
	return ((j2gort.JStrIndexOf(originalQname, "Asserts") >= 0) && (j2gort.JStrIndexOf(originalQname, ".$assert") >= 0))
}
func (*J2clAssertRemovalPass) AlwaysShouldTraverse_()      {}
func (*J2clAssertRemovalPass) CanReleaseSyntacticScopes_() {}
