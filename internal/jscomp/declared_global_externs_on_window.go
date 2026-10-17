// Port of com/google/javascript/jscomp/DeclaredGlobalExternsOnWindow.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DeclaredGlobalExternsOnWindow is a compiler pass to normalize externs by
// declaring global names on the "window" object, if it is declared in
// externs. The new declarations are added to the window instance, not to
// Window.prototype.
type DeclaredGlobalExternsOnWindow struct {
	NodeTraversal_AbstractShallowStatementCallbackBase

	compiler AbstractCompiler
	nodes    *jutil.LinkedHashSet[*rhino.Node]

	// Whether there is a "var window" declaration in the externs.
	windowInExterns bool
}

// Java: DeclaredGlobalExternsOnWindow.WINDOW_NAME
const DeclaredGlobalExternsOnWindow_WINDOW_NAME = "window"

// Java: DeclaredGlobalExternsOnWindow(AbstractCompiler)
func NewDeclaredGlobalExternsOnWindow(compiler AbstractCompiler) *DeclaredGlobalExternsOnWindow {
	d := &DeclaredGlobalExternsOnWindow{
		nodes:           jutil.NewLinkedHashSet[*rhino.Node](),
		windowInExterns: false,
	}
	d.compiler = compiler
	return d
}

// Java: DeclaredGlobalExternsOnWindow.process(Node, Node)
func (d *DeclaredGlobalExternsOnWindow) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(d.compiler, externs, d)
	d.AddWindowProperties()
}

// Java: DeclaredGlobalExternsOnWindow.addWindowProperties()
func (d *DeclaredGlobalExternsOnWindow) AddWindowProperties() {
	if !d.nodes.IsEmpty() {
		for node := range d.nodes.All() {
			DeclaredGlobalExternsOnWindow_AddExtern(node, d.windowInExterns)
			d.compiler.ReportChangeToEnclosingScope(node)
		}
	}
}

// Java: DeclaredGlobalExternsOnWindow.addExtern(Node, boolean)
func DeclaredGlobalExternsOnWindow_AddExtern(node *rhino.Node, defineOnWindow bool) {
	name := node.GetString()
	oldJSDocInfo := NodeUtil_GetBestJSDocInfo(node)

	var globalRef *rhino.Node
	if defineOnWindow {
		globalRef = rhino.IR_Name(DeclaredGlobalExternsOnWindow_WINDOW_NAME)
	} else {
		globalRef = rhino.IR_ThisNode()
	}
	getprop := rhino.IR_Getprop(globalRef, name)
	newNode := getprop

	if oldJSDocInfo != nil {
		var builder *rhino.JSDocInfo_Builder

		if oldJSDocInfo.IsConstructorOrInterface() ||
			oldJSDocInfo.HasEnumParameterType() {
			nameNode := rhino.IR_Name(name)
			newNode = rhino.IR_Assign(getprop, nameNode)

			builder = rhino.JSDocInfo_Builder_()
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
				newNode = rhino.IR_Assign(getprop, rhino.IR_Name(name))
			} else {
				rhs := NodeUtil_GetRValueOfLValue(node)
				// Type-aliasing definition
				if oldJSDocInfo.HasConstAnnotation() && rhs != nil && rhs.IsQualifiedName() {
					newNode = rhino.IR_Assign(getprop, rhs.CloneTree())
				}
			}
			builder = rhino.JSDocInfo_Builder_CopyFrom(oldJSDocInfo)
		}

		// TODO(blickly): Remove these suppressions when all externs declarations on window are gone.
		builder.RecordSuppressions(jutil.NewLinkedHashSet("const", "duplicate"))
		jsDocInfo := builder.Build()
		newNode.SetJSDocInfo(jsDocInfo)
	}

	newNode.SrcrefTree(node)
	newNode.SetOriginalName(&name)
	newNode.MakeNonIndexableRecursive()
	node.GetGrandparent().AddChildToBack(rhino.IR_ExprResult(newNode))
}

// Java: DeclaredGlobalExternsOnWindow.visit(NodeTraversal, Node, Node)
func (d *DeclaredGlobalExternsOnWindow) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsFunction() {
		d.nodes.Add(n.GetFirstChild())
	} else if n.IsVar() {
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			if c.GetString() == DeclaredGlobalExternsOnWindow_WINDOW_NAME {
				d.windowInExterns = true
				continue
			}
			// Skip 'location' since there is an existing definition
			// for window.location which conflicts with the "var location" one.
			if c.GetString() != "location" {
				d.nodes.Add(c)
			}
		}
	}
}
