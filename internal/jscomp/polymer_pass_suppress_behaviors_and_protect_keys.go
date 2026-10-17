// Port of com/google/javascript/jscomp/PolymerPassSuppressBehaviorsAndProtectKeys.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// PolymerPassSuppressBehaviorsAndProtectKeys: for every Polymer Behavior,
// strip property type annotations and add suppress checktypes on functions.
//
// Also find listener and hostAttribute keys in the behavior, and mark them as
// quoted. This protects the keys from being renamed.
//
// Java: final class PolymerPassSuppressBehaviorsAndProtectKeys extends ExternsSkippingCallback
type PolymerPassSuppressBehaviorsAndProtectKeys struct {
	NodeTraversal_ExternsSkippingCallbackBase
	compiler AbstractCompiler
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeys(AbstractCompiler)
func NewPolymerPassSuppressBehaviorsAndProtectKeys(compiler AbstractCompiler) *PolymerPassSuppressBehaviorsAndProtectKeys {
	return &PolymerPassSuppressBehaviorsAndProtectKeys{compiler: compiler}
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeys.visit(NodeTraversal, Node, Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if polymerPassSuppressBehaviorsAndProtectKeys_isBehavior(n) {
		if !NodeUtil_IsNameDeclaration(n) && !n.IsAssign() {
			p.compiler.Report(JSError_MakeWithNode(n, PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR))
			return
		}

		// Find listener and hostAttribute keys in the behavior, and mark them as quoted.
		// This ensures that the keys are not renamed.
		p.traverseBehaviorsAndProtectKeys(n)

		// Add @nocollapse.
		newDocs := rhino.JSDocInfo_Builder_MaybeCopyFrom(n.GetJSDocInfo())
		newDocs.RecordNoCollapse()
		n.SetJSDocInfo(newDocs.Build())

		behaviorValue := n.GetSecondChild()
		if NodeUtil_IsNameDeclaration(n) {
			behaviorValue = n.GetFirstFirstChild()
		}
		p.suppressBehavior(behaviorValue, n)
	}
}

// traverseBehaviorsAndProtectKeys traverses the behavior, finds listener and
// hostAttribute keys, and marks them as quoted.
//
// Java: PolymerPassSuppressBehaviorsAndProtectKeys.traverseBehaviorsAndProtectKeys(Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) traverseBehaviorsAndProtectKeys(n *rhino.Node) {
	if n.IsObjectLit() {
		PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(n, p.compiler)
		PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(n)
	}

	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		p.traverseBehaviorsAndProtectKeys(child)
	}
}

// suppressBehavior strips property type annotations and adds suppressions on
// functions.
//
// Java: PolymerPassSuppressBehaviorsAndProtectKeys.suppressBehavior(Node, Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) suppressBehavior(behaviorValue *rhino.Node, reportNode *rhino.Node) {
	if behaviorValue == nil {
		p.compiler.Report(JSError_MakeWithNode(reportNode, PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR))
		return
	}

	if behaviorValue.IsArrayLit() {
		for child := behaviorValue.GetFirstChild(); child != nil; child = child.GetNext() {
			p.suppressBehavior(child, behaviorValue)
		}
	} else if behaviorValue.IsObjectLit() {
		p.stripPropertyTypes(behaviorValue)
		p.addBehaviorSuppressions(behaviorValue)
	}
}

// polymerPassSuppressBehaviorsAndProtectKeys_isBehavior returns whether the
// node is the declaration of a Behavior.
//
// Java: PolymerPassSuppressBehaviorsAndProtectKeys.isBehavior(Node)
func polymerPassSuppressBehaviorsAndProtectKeys_isBehavior(value *rhino.Node) bool {
	return value.GetJSDocInfo() != nil && value.GetJSDocInfo().IsPolymerBehavior()
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeys.stripPropertyTypes(Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) stripPropertyTypes(behaviorValue *rhino.Node) {
	properties := PolymerPassStaticUtils_ExtractProperties(
		behaviorValue,
		PolymerClassDefinition_DefinitionType_ObjectLiteral,
		p.compiler,
		/* constructor= */ nil)
	for _, property := range properties {
		property.name.SetJSDocInfo(nil)
	}
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeys.suppressDefaultValues(Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) suppressDefaultValues(behaviorValue *rhino.Node) {
	for _, property := range PolymerPassStaticUtils_ExtractProperties(
		behaviorValue,
		PolymerClassDefinition_DefinitionType_ObjectLiteral,
		p.compiler,
		/* constructor= */ nil) {
		if !property.value.IsObjectLit() {
			continue
		}

		defaultValue := NodeUtil_GetFirstPropMatchingKey(property.value, "value")
		if defaultValue == nil || !defaultValue.IsFunction() {
			continue
		}
		defaultValueKey := defaultValue.GetParent()
		p.addSuppressionJsDocAnnotations(defaultValueKey)
	}
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeys.addBehaviorSuppressions(Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) addBehaviorSuppressions(behaviorValue *rhino.Node) {
	for keyNode := behaviorValue.GetFirstChild(); keyNode != nil; keyNode = keyNode.GetNext() {
		if keyNode.GetFirstChild().IsFunction() {
			p.addSuppressionJsDocAnnotations(keyNode)
		}
	}
	p.suppressDefaultValues(behaviorValue)
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeys.addSuppressionJsDocAnnotations(Node)
func (p *PolymerPassSuppressBehaviorsAndProtectKeys) addSuppressionJsDocAnnotations(node *rhino.Node) {
	suppressDoc := rhino.JSDocInfo_Builder_MaybeCopyFrom(node.GetJSDocInfo())
	suppressDoc.RecordSuppressions(jutil.NewLinkedHashSet("checkTypes", "globalThis", "visibility"))
	node.SetJSDocInfo(suppressDoc.Build())
}
