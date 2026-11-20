/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DeclaredGlobalExternsOnWindow.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.DeclaredGlobalExternsOnWindow
type DeclaredGlobalExternsOnWindow struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler        AbstractCompiler                  `gson:"compiler"`
	nodes           j2gort.OrderedSet[pkg_rhino.Node] `gson:"nodes"`
	windowInExterns bool                              `gson:"windowInExterns"`
}

func (this_ *DeclaredGlobalExternsOnWindow) BindSelf_(s any) {
	_ = s
}
func (this_ *DeclaredGlobalExternsOnWindow) IsJavaNil_() bool { return this_ == nil }
func (*DeclaredGlobalExternsOnWindow) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeclaredGlobalExternsOnWindow", "DeclaredGlobalExternsOnWindow"
}
func (this_ *DeclaredGlobalExternsOnWindow) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowStatementCallbackBase.IdentityHash_()
}
func (*DeclaredGlobalExternsOnWindow) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeclaredGlobalExternsOnWindow", "com.google.javascript.jscomp.DeclaredGlobalExternsOnWindow", func(v_ any) bool {
		return j2gort.InstanceOf[*DeclaredGlobalExternsOnWindow](v_)
	})
}
func (*DeclaredGlobalExternsOnWindow) Implements_jscomp_NodeTraversal_AbstractShallowStatementCallback() {
}
func (*DeclaredGlobalExternsOnWindow) Implements_jscomp_CompilerPass()           {}
func (*DeclaredGlobalExternsOnWindow) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_DeclaredGlobalExternsOnWindow_WINDOW_NAME() string { return "window" }
func (this_ *DeclaredGlobalExternsOnWindow) initFields_DeclaredGlobalExternsOnWindow() {
	_ = this_
	this_.nodes = j2gort.NewOrderedSet[pkg_rhino.Node]()

}

// Java constructor: DeclaredGlobalExternsOnWindow(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewDeclaredGlobalExternsOnWindow(compiler AbstractCompiler) *DeclaredGlobalExternsOnWindow {
	this_ := &DeclaredGlobalExternsOnWindow{}
	InitDeclaredGlobalExternsOnWindow(this_, compiler)
	return this_
}

func InitDeclaredGlobalExternsOnWindow(this_ *DeclaredGlobalExternsOnWindow, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractShallowStatementCallback(&this_.NodeTraversal_AbstractShallowStatementCallbackBase)
	this_.initFields_DeclaredGlobalExternsOnWindow()
	this_.compiler = compiler
}

// Java method: DeclaredGlobalExternsOnWindow.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DeclaredGlobalExternsOnWindow) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, externs, this_)
	this_.priv_DeclaredGlobalExternsOnWindow_addWindowProperties_2371()
}

// Java method: DeclaredGlobalExternsOnWindow.addWindowProperties()
func (this_ *DeclaredGlobalExternsOnWindow) priv_DeclaredGlobalExternsOnWindow_addWindowProperties_2371() {
	if !(this_.nodes.IsEmpty()) {
		{
			cur_1 := this_.nodes.Iter()
			for cur_1.HasNext() {
				node := cur_1.
					Next()
				_ = node
				DeclaredGlobalExternsOnWindow_AddExtern(node, this_.windowInExterns)
				this_.compiler.ReportChangeToEnclosingScope(node)
			}
		}
	}
}

// Java static method: DeclaredGlobalExternsOnWindow.addExtern(com.google.javascript.rhino.Node node, boolean defineOnWindow)
func DeclaredGlobalExternsOnWindow_AddExtern(node pkg_rhino.Node, defineOnWindow bool) {
	_ = node
	_ = defineOnWindow
	var name string = node.GetString()
	_ = name
	var oldJSDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(node)
	_ = oldJSDocInfo
	var globalRef pkg_rhino.Node
	_ = globalRef
	if defineOnWindow {
		globalRef = pkg_rhino.IR_Name("window")
	} else {
		globalRef = pkg_rhino.IR_ThisNode()
	}
	var getprop pkg_rhino.Node = pkg_rhino.IR_Getprop_2(globalRef, name)
	_ = getprop
	var newNode pkg_rhino.Node = getprop
	_ = newNode
	if oldJSDocInfo != nil {
		var builder *pkg_rhino.JSDocInfo_Builder = nil
		_ = builder
		if oldJSDocInfo.IsConstructorOrInterface() || oldJSDocInfo.HasEnumParameterType() {
			var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(name)
			_ = nameNode
			newNode = pkg_rhino.IR_Assign(getprop, nameNode)
			builder = pkg_rhino.JSDocInfo_S_Builder()
			if oldJSDocInfo.IsConstructor() {
				builder.RecordConstructor()
			}
			if oldJSDocInfo.IsInterface() {
				builder.RecordInterface()
			}
			if oldJSDocInfo.UsesImplicitMatch() {
				builder.RecordImplicitMatch()
			}
			if oldJSDocInfo.HasEnumParameterType() {
				builder.RecordEnumParameterType(oldJSDocInfo.GetEnumParameterType())
			}
		} else {
			if NodeUtil_IsNamespaceDecl(node) {
				newNode = pkg_rhino.IR_Assign(getprop, pkg_rhino.IR_Name(name))
			} else {
				var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(node)
				_ = rhs
				if (oldJSDocInfo.HasConstAnnotation() && (rhs != 0)) && rhs.IsQualifiedName() {
					newNode = pkg_rhino.IR_Assign(getprop, rhs.CloneTree_0())
				}
			}
			builder = pkg_rhino.JSDocInfo_Builder_CopyFrom(oldJSDocInfo)
		}
		builder.RecordSuppressions_1(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.ImmSetOfStr2("const", "duplicate")))
		var jsDocInfo *pkg_rhino.JSDocInfo = builder.Build_0()
		_ = jsDocInfo
		newNode.SetJSDocInfo(jsDocInfo)
	}
	newNode.SrcrefTree(node)
	newNode.SetOriginalName(name)
	newNode.MakeNonIndexableRecursive()
	node.GetGrandparent().AddChildToBack(pkg_rhino.IR_ExprResult(newNode))
}

// Java method: DeclaredGlobalExternsOnWindow.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *DeclaredGlobalExternsOnWindow) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsFunction() {
		this_.nodes.Add(n.GetFirstChild())
	} else if n.IsVar() {
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				{
					var iifeCond_1 bool
					{

						sr_2 := c.GetString()
						se_3 := "window"
						iifeCond_1 = j2gort.CheckNotNullString(sr_2) == se_3
					}
					if iifeCond_1 {
						this_.windowInExterns = true
						continue
					}
				}
				if !(j2gort.CheckNotNullString(c.GetString()) == "location") {
					this_.nodes.Add(c)
				}
			}
		}
	}
}
func (*DeclaredGlobalExternsOnWindow) CanReleaseSyntacticScopes_() {}
