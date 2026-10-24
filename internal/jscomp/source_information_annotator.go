// Port of com/google/javascript/jscomp/SourceInformationAnnotator.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// SourceInformationAnnotator annotates nodes with information from their
// original input file before the compiler performs work that changes this
// information (such as its original location, its original name, etc).
//
// Information saved:
//
//   - Annotates all NAME nodes with an ORIGINALNAME_PROP indicating its
//     original name.
//   - Annotates all GETPROP and OPTCHAIN_GETPROP nodes with an original name.
//   - Annotates all OBJECT_LITERAL unquoted string key nodes with an original
//     name.
//   - Annotates all FUNCTION nodes with an original name indicating its
//     nearest original name.
type SourceInformationAnnotator struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	sourceFileToCheck *string // nil is null
}

// Java: SourceInformationAnnotator(String)
func NewSourceInformationAnnotator(sourceFileToCheck *string) *SourceInformationAnnotator {
	return &SourceInformationAnnotator{sourceFileToCheck: sourceFileToCheck}
}

// SourceInformationAnnotator_Create returns an annotator that sets the
// original name of nodes.
//
// Java: SourceInformationAnnotator.create()
func SourceInformationAnnotator_Create() *SourceInformationAnnotator {
	return NewSourceInformationAnnotator(nil)
}

// SourceInformationAnnotator_CreateWithAnnotationChecks returns an annotator
// that both sets original name and verifies Node#getSourceFileName is
// accurate.
//
// Java: SourceInformationAnnotator.createWithAnnotationChecks(String)
func SourceInformationAnnotator_CreateWithAnnotationChecks(sourceFile string) *SourceInformationAnnotator {
	return NewSourceInformationAnnotator(&sourceFile)
}

// Java: SourceInformationAnnotator.visit(NodeTraversal, Node, Node)
func (a *SourceInformationAnnotator) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Verify the source file is annotated.
	if a.sourceFileToCheck != nil {
		sourceFileName := n.GetSourceFileName()
		jutil.CheckState(sourceFileName != nil && *a.sourceFileToCheck == *sourceFileName)
	}

	// Annotate the original name.
	if SourceInformationAnnotator_IsStringNodeRequiringOriginalName(n) {
		SourceInformationAnnotator_SetOriginalNameFromNode(n, n)
		return
	}

	if n.IsFunction() {
		annotateFunctionOriginalName(n)
	}
}

// AnnotateTree performs the post-order Visit walk directly without allocating
// a NodeTraversal or Deque stack when dev-mode sourceFileToCheck is disabled.
func (a *SourceInformationAnnotator) AnnotateTree(n rhino.Node) {
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		a.AnnotateTree(child)
	}
	if SourceInformationAnnotator_IsStringNodeRequiringOriginalName(n) {
		SourceInformationAnnotator_SetOriginalNameFromNode(n, n)
		return
	}
	if n.IsFunction() {
		annotateFunctionOriginalName(n)
	}
}

func annotateFunctionOriginalName(n rhino.Node) {
	if n.GetOriginalName() != nil {
		return
	}
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_NAME:
		SourceInformationAnnotator_SetOriginalNameFromNode(n, parent)
		return
	case rhino.Token_ASSIGN:
		firstChild := parent.GetFirstChild()
		if firstChild.IsName() {
			SourceInformationAnnotator_SetOriginalNameFromNode(n, firstChild)
			return
		}
		if qname, ok := firstChild.GetQualifiedNameString(); ok && qname != "" {
			n.SetOriginalNameString(qname)
		}
		return
	default:
		funNameNode := n.GetFirstChild()
		if !funNameNode.IsEmpty() && funNameNode.GetString() != "" {
			SourceInformationAnnotator_SetOriginalNameFromNode(n, funNameNode)
			return
		}
	}
	switch parent.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF, rhino.Token_STRING_KEY:
		SourceInformationAnnotator_SetOriginalNameFromNode(n, parent)
	case rhino.Token_NUMBER:
		if s := NodeUtil_GetStringValue(parent); s != nil && *s != "" {
			n.SetOriginalNameString(*s)
		}
	}
}

// SourceInformationAnnotator_IsStringNodeRequiringOriginalName: whether
// JSCompiler attempts to preserve the original string attached to this node
// on the AST post-mangling and in source maps.
//
// Java: SourceInformationAnnotator.isStringNodeRequiringOriginalName(Node)
func SourceInformationAnnotator_IsStringNodeRequiringOriginalName(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME:
		return true
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		return node.GetParent().IsObjectLit() && !node.IsQuotedStringKey()
	default:
		return false
	}
}

// Java: SourceInformationAnnotator.setOriginalName(Node, String)
func SourceInformationAnnotator_SetOriginalName(n rhino.Node, name string) {
	if name != "" && n.GetOriginalName() == nil {
		n.SetOriginalName(&name)
	}
}

// SourceInformationAnnotator_SetOriginalNameFromNode sets the original name
// from the "getString" of the supplied node. This avoids having to check
// for if the string is interned.
//
// Java: SourceInformationAnnotator.setOriginalName(Node, Node)
func SourceInformationAnnotator_SetOriginalNameFromNode(n rhino.Node, name rhino.Node) {
	if name.GetString() != "" && n.GetOriginalName() == nil {
		n.SetOriginalNameFromName(name)
	}
}
