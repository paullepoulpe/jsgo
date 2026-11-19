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
// Port of com/google/javascript/jscomp/Es6RewriteScriptsToModules.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RewriteScriptsToModules
type Es6RewriteScriptsToModules struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *Es6RewriteScriptsToModules) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteScriptsToModules) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteScriptsToModules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteScriptsToModules", "Es6RewriteScriptsToModules"
}
func (this_ *Es6RewriteScriptsToModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteScriptsToModules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteScriptsToModules", "com.google.javascript.jscomp.Es6RewriteScriptsToModules", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteScriptsToModules](v_)
	})
}
func (*Es6RewriteScriptsToModules) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {}
func (*Es6RewriteScriptsToModules) Implements_jscomp_CompilerPass()                           {}
func (*Es6RewriteScriptsToModules) Implements_jscomp_NodeTraversal_Callback()                 {}

func (this_ *Es6RewriteScriptsToModules) initFields_Es6RewriteScriptsToModules() {
	_ = this_
}

// Java constructor: Es6RewriteScriptsToModules(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteScriptsToModules(compiler AbstractCompiler) *Es6RewriteScriptsToModules {
	this_ := &Es6RewriteScriptsToModules{}
	InitEs6RewriteScriptsToModules(this_, compiler)
	return this_
}

func InitEs6RewriteScriptsToModules(this_ *Es6RewriteScriptsToModules, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_Es6RewriteScriptsToModules()
	this_.compiler = compiler
}

// Java method: Es6RewriteScriptsToModules.forceToEs6Module(com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteScriptsToModules) ForceToEs6Module(root pkg_rhino.Node) {
	_ = root
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		return
	}
	var moduleNode pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_MODULE_BODY()).Srcref(root)
	_ = moduleNode
	moduleNode.AddChildrenToBack(root.RemoveChildren())
	root.AddChildToBack(moduleNode)
	this_.compiler.ReportChangeToChangeScope(root)
}

// Java method: Es6RewriteScriptsToModules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteScriptsToModules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var file pkg_rhino.Node = root.GetFirstChild()
		_ = file
		for ; file != 0; file = file.GetNext() {
			if !(file.IsScript()) {
				j2gort.CheckArgument(false)
			}
			NodeTraversal_Traverse(this_.compiler, file, this_)
		}
	}
}

// Java method: Es6RewriteScriptsToModules.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteScriptsToModules) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsScript() {
		var moduleType CompilerInput_ModuleType = this_.compiler.GetInput(n.GetInputId()).GetJsModuleType()
		_ = moduleType
		if moduleType == StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT() {
			this_.ForceToEs6Module(n)
		}
	}
	return false
}
func (*Es6RewriteScriptsToModules) NoOpVisit_()                 {}
func (*Es6RewriteScriptsToModules) CanReleaseSyntacticScopes_() {}
