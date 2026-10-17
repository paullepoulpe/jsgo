// Port of com/google/javascript/jscomp/PolymerPass.java @ closure-compiler v20261015

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const polymerPass_VIRTUAL_FILE = "<PolymerPass.java>"

// PolymerPass rewrites "Polymer({})" calls into a form that is suitable for
// type checking and dead code elimination. Also ensures proper format and
// types.
//
// Only works with Polymer version 0.8 and above.
//
// Design and examples: https://github.com/google/closure-compiler/wiki/Polymer-Pass
//
// Java: final class PolymerPass extends ExternsSkippingCallback implements CompilerPass
type PolymerPass struct {
	NodeTraversal_ExternsSkippingCallbackBase

	compiler   AbstractCompiler
	tagNameMap *jutil.LinkedHashMap[string, string]

	polymerElementExterns        *rhino.Node
	nativeExternsAdded           *jutil.LinkedHashSet[string]
	polymerElementProps          []*rhino.Node
	globalNames                  *GlobalNamespace
	behaviorExtractor            *PolymerBehaviorExtractor
	warnedPolymer1ExternsMissing bool
	propertySinkExternInjected   bool
}

// Java: PolymerPass(AbstractCompiler)
func NewPolymerPass(compiler AbstractCompiler) *PolymerPass {
	p := &PolymerPass{}
	p.warnedPolymer1ExternsMissing = false
	p.propertySinkExternInjected = false
	p.compiler = compiler
	p.tagNameMap = TagNameToType_GetMap()
	p.nativeExternsAdded = jutil.NewLinkedHashSet[string]()
	return p
}

// Java: PolymerPass.process(Node, Node)
func (p *PolymerPass) Process(externs *rhino.Node, root *rhino.Node) {
	externsCallback := NewPolymerPassFindExterns()
	NodeTraversal_Traverse(p.compiler, externs, externsCallback)
	p.polymerElementExterns = externsCallback.GetPolymerElementExterns()
	p.polymerElementProps = externsCallback.GetPolymerElementProps()

	p.globalNames = NewGlobalNamespaceWithRoots(p.compiler, externs, root)
	p.behaviorExtractor = NewPolymerBehaviorExtractor(
		p.compiler, p.globalNames, p.compiler.GetModuleMetadataMap(), p.compiler.GetModuleMap())

	externsAndJsRoot := root.GetParent()
	NodeTraversal_Traverse(p.compiler, externsAndJsRoot, p)
	suppressBehaviorsCallback := NewPolymerPassSuppressBehaviorsAndProtectKeys(p.compiler)
	NodeTraversal_Traverse(p.compiler, externsAndJsRoot, suppressBehaviorsCallback)
}

// Java: PolymerPass.visit(NodeTraversal, Node, Node)
func (p *PolymerPass) Visit(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	if PolymerPassStaticUtils_IsPolymerCall(node) {
		if p.polymerElementExterns != nil {
			p.rewritePolymer1ClassDefinition(node, parent, traversal)
		} else if !p.warnedPolymer1ExternsMissing {
			p.compiler.Report(JSError_MakeWithNode(node, PolymerPassErrors_POLYMER_MISSING_EXTERNS))
			p.warnedPolymer1ExternsMissing = true
		}
	} else if PolymerPassStaticUtils_IsPolymerClass(node) {
		p.rewritePolymer2ClassDefinition(node, traversal)
	}
}

// Java: PolymerPass.getModuleMetadata(NodeTraversal)
func (p *PolymerPass) getModuleMetadata(traversal *NodeTraversal) *ModuleMetadataMap_ModuleMetadata {
	script := traversal.GetCurrentScript()
	if script != nil && script.GetFirstChild().IsModuleBody() {
		return ModuleImportResolver_GetModuleFromScopeRoot(
			p.compiler.GetModuleMap(), p.compiler, script.GetFirstChild()).
			Metadata()
	}
	// Check for bundled goog.loadModule calls
	scopeRoot := traversal.GetScopeRoot()
	if NodeUtil_IsBundledGoogModuleScopeRoot(scopeRoot) {
		return ModuleImportResolver_GetModuleFromScopeRoot(
			p.compiler.GetModuleMap(), p.compiler, scopeRoot).
			Metadata()
	}
	return nil
}

// rewritePolymer1ClassDefinition handles Polymer 1.x and Polymer 2 Legacy
// Element Definitions
//
// Java: PolymerPass.rewritePolymer1ClassDefinition(Node, Node, NodeTraversal)
func (p *PolymerPass) rewritePolymer1ClassDefinition(node *rhino.Node, parent *rhino.Node, traversal *NodeTraversal) {
	grandparent := parent.GetParent()
	if grandparent.IsConst() {
		grandparent.SetToken(rhino.Token_LET)
		scriptNode := traversal.GetCurrentScript()
		if scriptNode != nil {
			NodeUtil_AddFeatureToScript(scriptNode, parser.FeatureSet_Feature_LET_DECLARATIONS, p.compiler)
		}
		traversal.ReportCodeChange()
	}
	if NodeUtil_IsNameDeclaration(grandparent) && grandparent.GetParent().IsExport() {
		p.normalizePolymerExport(grandparent, grandparent.GetParent())
		traversal.ReportCodeChange()
	}
	def := PolymerClassDefinition_ExtractFromCallNode(
		node, p.compiler, p.getModuleMetadata(traversal), p.behaviorExtractor)
	if def != nil {
		if def.nativeBaseElement != nil {
			p.appendPolymerElementExterns(def)
		}
		rewriter := NewPolymerClassRewriter(p.compiler)
		rewriter.RewritePolymerCall(def, traversal)
	}
}

// rewritePolymer2ClassDefinition handles Polymer 2.x Class Nodes
//
// Java: PolymerPass.rewritePolymer2ClassDefinition(Node, NodeTraversal)
func (p *PolymerPass) rewritePolymer2ClassDefinition(node *rhino.Node, traversal *NodeTraversal) {
	def := PolymerClassDefinition_ExtractFromClassNode(node, p.compiler)
	if def != nil {
		rewriter := NewPolymerClassRewriter(p.compiler)
		rewriter.propertySinkExternInjected = p.propertySinkExternInjected
		rewriter.RewritePolymerClassDeclaration(node, traversal, def)
		p.propertySinkExternInjected = rewriter.propertySinkExternInjected

		// Support Polymer implicit property edge cases by just @exporting all non-private properties.
		polymerPass_exportTsickleProps(def.definition, *def.target.GetQualifiedName())
	}
}

// normalizePolymerExport replaces `export let Element = ...` with
// `let Element = ...; export {Element};`
//
// Java: PolymerPass.normalizePolymerExport(Node, Node)
func (p *PolymerPass) normalizePolymerExport(nameDecl *rhino.Node, export *rhino.Node) {
	name := nameDecl.GetFirstChild()
	nameDecl.Detach().InsertBefore(export)

	exportSpec := rhino.NewNode(rhino.Token_EXPORT_SPEC)
	exportSpec.AddChildToFront(name.CloneNode())
	exportSpec.AddChildToFront(name.CloneNode())
	export.AddChildToFront(rhino.NewNodeWithChild(rhino.Token_EXPORT_SPECS, exportSpec).SrcrefTree(export))
}

// appendPolymerElementExterns duplicates the PolymerElement externs with a
// different element base class if needed. For example, if the base class is
// HTMLInputElement, then a class PolymerInputElement will be added. If the
// element does not extend a native HTML element, this method is a no-op.
//
// Java: PolymerPass.appendPolymerElementExterns(PolymerClassDefinition)
func (p *PolymerPass) appendPolymerElementExterns(def *PolymerClassDefinition) {
	if !p.nativeExternsAdded.Add(*def.nativeBaseElement) {
		return
	}

	block := rhino.IR_Block()

	baseExterns := p.polymerElementExterns.CloneTree()
	polymerElementType := PolymerPassStaticUtils_GetPolymerElementType(def)
	baseExterns.GetFirstChild().SetString(polymerElementType)

	elementType, ok := p.tagNameMap.Lookup(*def.nativeBaseElement)
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(def.descriptor, PolymerPassErrors_POLYMER_INVALID_EXTENDS, *def.nativeBaseElement))
		return
	}
	elementBaseType := rhino.NewJSTypeExpression(
		rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String(elementType).SrcrefTree(p.polymerElementExterns)),
		polymerPass_VIRTUAL_FILE)
	baseDocs := rhino.JSDocInfo_Builder_CopyFrom(baseExterns.GetJSDocInfo())
	baseDocs.ChangeBaseType(elementBaseType)
	baseExterns.SetJSDocInfo(baseDocs.Build())
	block.AddChildToBack(baseExterns)

	for _, baseProp := range p.polymerElementProps {
		newProp := baseProp.CloneTree()
		newPropRootName := NodeUtil_GetRootOfQualifiedName(newProp.GetFirstFirstChild())
		newPropRootName.SetString(polymerElementType)
		block.AddChildToBack(newProp)
	}

	block.SrcrefTreeIfMissing(p.polymerElementExterns)

	parent := p.polymerElementExterns.GetParent()
	stmts := block.RemoveChildren()
	parent.AddChildrenAfter(stmts, p.polymerElementExterns)

	p.compiler.ReportChangeToEnclosingScope(stmts)
}

// polymerPass_exportTsickleProps finds the tsickle-generated `if (false)`
// block and exports all non-private properties.
//
// Java: PolymerPass.exportTsickleProps(Node, String)
func polymerPass_exportTsickleProps(classNode *rhino.Node, classQualifiedName string) {
	classStmt := jutil.CheckNotNull(NodeUtil_GetEnclosingStatement(classNode))

	for sibling := classStmt.GetNext(); sibling != nil; sibling = sibling.GetNext() {
		ifBlock := polymerPass_maybeGetTsicklePropTypesIfBlock(sibling)
		if ifBlock == nil {
			continue
		}

		for stmt := ifBlock.GetFirstChild(); stmt != nil; stmt = stmt.GetNext() {
			prop := polymerPass_getTsickleClassPrototypeProp(stmt, classQualifiedName)
			if prop == nil {
				continue
			}

			jsdoc := jutil.CheckNotNullMsg1(prop.GetJSDocInfo(), "No JSDoc found for %s", prop)
			if jsdoc.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE {
				continue
			}

			builder := rhino.JSDocInfo_Builder_CopyFrom(jsdoc)
			builder.RecordExport()
			// @export implicitly overwrites visibility.
			builder.OverwriteVisibility(rhino.JSDocInfo_Visibility_INHERITED)
			prop.SetJSDocInfo(builder.Build())
		}
	}
}

// polymerPass_maybeGetTsicklePropTypesIfBlock: if the given node is an
// `if (false) {}`, returns its block node. Otherwise returns null.
//
// Java: PolymerPass.maybeGetTsicklePropTypesIfBlock(Node)
func polymerPass_maybeGetTsicklePropTypesIfBlock(node *rhino.Node) *rhino.Node {
	if !node.IsIf() {
		return nil
	}
	condition := node.GetFirstChild()
	if !condition.IsFalse() {
		return nil
	}
	block := condition.GetNext()
	jutil.CheckStateMsg1(block != nil && block.IsBlock(), "Expected block but got %s", block)
	return block
}

// polymerPass_getTsickleClassPrototypeProp looks for
// `ClassName.prototype.foo` expressions and returns the property node (if
// found).
//
// Java: PolymerPass.getTsickleClassPrototypeProp(Node, String)
func polymerPass_getTsickleClassPrototypeProp(stmt *rhino.Node, classQualifiedName string) *rhino.Node {
	if !stmt.IsExprResult() {
		return nil
	}
	expr := stmt.GetFirstChild()
	if !expr.IsGetProp() {
		return nil
	}
	owner := expr.GetFirstChild()
	if !owner.IsGetProp() || owner.GetString() != "prototype" {
		return nil
	}
	classRef := owner.GetFirstChild()
	if !classRef.MatchesQualifiedName(classQualifiedName) {
		return nil
	}
	return expr
}

// PolymerPass_MemberDefinition is any member of a Polymer element or
// Behavior. These can be functions, properties, etc.
//
// Java: static class PolymerPass.MemberDefinition
type PolymerPass_MemberDefinition struct {
	// Any JSDocInfo tied to this member.
	info *rhino.JSDocInfo

	// Name Node for the definition of this member.
	name *rhino.Node

	// Value Node (RHS) for the definition of this member.
	value *rhino.Node

	enclosingModule *rhino.Node
}

// Java: PolymerPass.MemberDefinition(JSDocInfo, Node, Node, Node)
func NewPolymerPass_MemberDefinition(info *rhino.JSDocInfo, name *rhino.Node, value *rhino.Node, enclosingModule *rhino.Node) *PolymerPass_MemberDefinition {
	return &PolymerPass_MemberDefinition{
		info:            info,
		name:            name,
		value:           value,
		enclosingModule: enclosingModule,
	}
}

// Java: PolymerPass.MemberDefinition.toString()
func (m *PolymerPass_MemberDefinition) String() string {
	return jutil.MoreObjects_ToStringHelper("MemberDefinition").Add("name", m.name).Add("value", m.value).String()
}
