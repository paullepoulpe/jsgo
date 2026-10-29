// Port of com/google/javascript/jscomp/PolymerPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/PolymerPassErrors.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/PolymerPassFindExterns.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/TagNameToType.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const polymerPass_VIRTUAL_FILE = "<PolymerPass.java>"

var polymerPassVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	s.Add(rhino.Token_CLASS)
	return s
}()

var _ nodeTraversal_visitTokens = (*PolymerPass)(nil)

func (*PolymerPass) visitTokens() *ast.TokenSet {
	return &polymerPassVisitTokens
}

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
	tagNameMap map[string]string

	polymerElementExterns        rhino.Node
	nativeExternsAdded           map[string]bool
	polymerElementProps          []rhino.Node
	globalNames                  *GlobalNamespace
	behaviorExtractor            *PolymerBehaviorExtractor
	warnedPolymer1ExternsMissing bool
	propertySinkExternInjected   bool
}

// Java: PolymerPass(AbstractCompiler)
func NewPolymerPass(compiler AbstractCompiler) *PolymerPass {
	return &PolymerPass{
		compiler:           compiler,
		tagNameMap:         TagNameToType_GetMap(),
		nativeExternsAdded: make(map[string]bool),
	}
}

// Java: PolymerPass.process(Node, Node)
func (p *PolymerPass) Process(externs rhino.Node, root rhino.Node) {
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
func (p *PolymerPass) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if PolymerPassStaticUtils_IsPolymerCall(node) {
		if !p.polymerElementExterns.IsNil() {
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
	if script, ok := ast.AsScript(traversal.GetCurrentScript()); ok {
		if first, hasFirst := script.FirstItem(); hasFirst && first.Node().IsModuleBody() {
			return ModuleImportResolver_GetModuleFromScopeRoot(
				p.compiler.GetModuleMap(), p.compiler, first.Node()).
				Metadata()
		}
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
func (p *PolymerPass) rewritePolymer1ClassDefinition(node rhino.Node, parent rhino.Node, traversal *NodeTraversal) {
	grandparent := parent.GetParent()
	if grandparent.IsConst() {
		grandparent.SetToken(rhino.Token_LET)
		scriptNode := traversal.GetCurrentScript()
		if !scriptNode.IsNil() {
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
func (p *PolymerPass) rewritePolymer2ClassDefinition(node rhino.Node, traversal *NodeTraversal) {
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
func (p *PolymerPass) normalizePolymerExport(nameDecl rhino.Node, export rhino.Node) {
	name := ast.MustNameDeclaration(nameDecl).MustFirstDeclarator().Node()
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
	if p.nativeExternsAdded[*def.nativeBaseElement] {
		return
	}
	p.nativeExternsAdded[*def.nativeBaseElement] = true

	block := rhino.IR_Block()

	baseExterns := p.polymerElementExterns.CloneTree()
	polymerElementType := PolymerPassStaticUtils_GetPolymerElementType(def)
	ast.MustNameDeclaration(baseExterns).MustFirstDeclarator().Node().SetString(polymerElementType)

	elementType, ok := p.tagNameMap[*def.nativeBaseElement]
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
		newPropRootName := NodeUtil_GetRootOfQualifiedName(
			ast.MustExprStmt(newProp).Expr().Node().GetFirstChild())
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
func polymerPass_exportTsickleProps(classNode rhino.Node, classQualifiedName string) {
	classStmt := rhino.CheckNotNull(NodeUtil_GetEnclosingStatement(classNode))

	for sibling := classStmt.GetNext(); !sibling.IsNil(); sibling = sibling.GetNext() {
		ifBlock := polymerPass_maybeGetTsicklePropTypesIfBlock(sibling)
		if ifBlock.IsNil() {
			continue
		}

		for stmt := ifBlock.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
			prop := polymerPass_getTsickleClassPrototypeProp(stmt, classQualifiedName)
			if prop.IsNil() {
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
func polymerPass_maybeGetTsicklePropTypesIfBlock(node rhino.Node) rhino.Node {
	ifNode, ok := ast.AsIf(node)
	if !ok || !ifNode.Cond().Node().IsFalse() {
		return rhino.NilNode
	}
	return ifNode.Then().Node()
}

// polymerPass_getTsickleClassPrototypeProp looks for
// `ClassName.prototype.foo` expressions and returns the property node (if
// found).
//
// Java: PolymerPass.getTsickleClassPrototypeProp(Node, String)
func polymerPass_getTsickleClassPrototypeProp(stmt rhino.Node, classQualifiedName string) rhino.Node {
	exprStmt, ok := ast.AsExprStmt(stmt)
	if !ok {
		return rhino.NilNode
	}
	expr, ok := exprStmt.Expr().AsGetProp()
	if !ok {
		return rhino.NilNode
	}
	owner, ok := expr.Object().AsGetProp()
	if !ok || owner.Property() != "prototype" {
		return rhino.NilNode
	}
	if !owner.Object().Node().MatchesQualifiedName(classQualifiedName) {
		return rhino.NilNode
	}
	return expr.Node()
}

// PolymerPass_MemberDefinition is any member of a Polymer element or
// Behavior. These can be functions, properties, etc.
//
// Java: static class PolymerPass.MemberDefinition
type PolymerPass_MemberDefinition struct {
	// Any JSDocInfo tied to this member.
	info *rhino.JSDocInfo

	// Name Node for the definition of this member.
	name rhino.Node

	// Value Node (RHS) for the definition of this member.
	value rhino.Node

	enclosingModule rhino.Node
}

// Java: PolymerPass.MemberDefinition(JSDocInfo, Node, Node, Node)
func NewPolymerPass_MemberDefinition(info *rhino.JSDocInfo, name rhino.Node, value rhino.Node, enclosingModule rhino.Node) *PolymerPass_MemberDefinition {
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

var (
	// TODO(jlklein): Switch back to an error when everyone is upgraded to
	// Polymer 1.0
	PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID = DiagnosticType_Warning(
		"JSC_POLYMER_DESCRIPTOR_NOT_VALID",
		"The argument to Polymer() is not an obj lit or the Polymer 2 class does not have a"+
			" static getter named 'config'. Ignoring this definition.")

	PolymerPassErrors_POLYMER_INVALID_BEHAVIOR = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_BEHAVIOR", "A Polymer behavior may not include an ''is'' property.")

	PolymerPassErrors_POLYMER_MISSING_IS = DiagnosticType_Error("JSC_POLYMER_MISSING_IS",
		"The class descriptor must include an ''is'' property.")

	PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS = DiagnosticType_Error(
		"JSC_POLYMER_UNEXPECTED_PARAMS", "The class definition has too many arguments.")

	PolymerPassErrors_POLYMER_MISSING_EXTERNS = DiagnosticType_Error(
		"JSC_POLYMER_MISSING_EXTERNS", "Missing Polymer externs.")

	PolymerPassErrors_POLYMER_INVALID_PROPERTY = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_PROPERTY", "Polymer property has an invalid or missing type.")

	PolymerPassErrors_POLYMER_INVALID_EXTENDS = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_EXTENDS",
		"Cannot extend HTML element ''{0}''. The element is probably either misspelled,"+
			" or needs to be added to the list of known elements.")

	PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY = DiagnosticType_Error(
		"JSC_POLYMER_INVALID_BEHAVIOR_ARRAY", "The behaviors property must be an array literal.")

	PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR = DiagnosticType_Error(
		"JSC_POLYMER_UNQUALIFIED_BEHAVIOR",
		"Behaviors must be global names or qualified names that are declared as object literals or "+
			"array literals of other valid Behaviors.")

	PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR = DiagnosticType_Error(
		"JSC_POLYMER_UNANNOTATED_BEHAVIOR",
		"Behavior declarations must be annotated with @polymerBehavior.")

	PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID = DiagnosticType_Error(
		"JSC_POLYMER_CLASS_PROPERTIES_INVALID",
		"The Polymer element class 'properties' getter does not return an object literal. "+
			"Ignoring this definition.")

	PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC = DiagnosticType_Error(
		"JSC_POLYMER_CLASS_PROPERTIES_NOT_STATIC",
		"The Polymer element class 'properties' getter is not declared static. "+
			"Ignoring this definition.")

	PolymerPassErrors_POLYMER_CLASS_UNNAMED = DiagnosticType_Warning(
		"JSC_POLYMER2_UNNAMED",
		"Unable to locate a valid name for the Polymer element class."+
			"Ignoring this definition.")

	PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC = DiagnosticType_Warning(
		"JSC_POLYMER_MISPLACED_PROPERTY_JSDOC",
		"When a Polymer property is declared in the constructor, its JSDoc "+
			"should only be in the constructor, not on the Polymer properties configuration.")

	PolymerPassErrors_POLYMER_UNPARSABLE_STRING = DiagnosticType_Error(
		"JSC_POLYMER_UNPARSABLE_STRING",
		"The Polymer computed property or complex observer string could not be parsed.")
)

const polymerPassFindExterns_POLYMER_ELEMENT_NAME = "PolymerElement"

var polymerPassFindExternsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_CONST)
	s.Add(rhino.Token_EXPR_RESULT)
	return s
}()

var _ nodeTraversal_visitTokens = (*PolymerPassFindExterns)(nil)

func (*PolymerPassFindExterns) visitTokens() *ast.TokenSet {
	return &polymerPassFindExternsVisitTokens
}

// PolymerPassFindExterns finds the externs for the PolymerElement base class
// and all of its properties in the externs.
//
// Java: final class PolymerPassFindExterns implements NodeTraversal.Callback
type PolymerPassFindExterns struct {
	polymerElementExterns rhino.Node
	polymerElementProps   []rhino.Node
}

// Java: PolymerPassFindExterns()
func NewPolymerPassFindExterns() *PolymerPassFindExterns {
	return &PolymerPassFindExterns{}
}

// Java: PolymerPassFindExterns.shouldTraverse(NodeTraversal, Node, Node)
func (f *PolymerPassFindExterns) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	/*
	 * The externs subtree also includes @typeSummary files, which are slightly different from
	 * normal externs. The Polymer externs may exist in either file type. However,
	 * some @typeSummaries are goog.modules, which definitely don't contain Polymer externs, and
	 * which we don't want to search.
	 */
	return !n.IsModuleBody()
}

// Java: PolymerPassFindExterns.visit(NodeTraversal, Node, Node)
func (f *PolymerPassFindExterns) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if polymerPassFindExterns_isPolymerElementExterns(n) {
		f.polymerElementExterns = n
	} else if polymerPassFindExterns_isPolymerElementPropExpr(n) {
		f.polymerElementProps = append(f.polymerElementProps, n)
	}
}

// GetPolymerElementProps returns all properties of the PolymerElement class
// in externs.
//
// Java: PolymerPassFindExterns.getPolymerElementProps()
func (f *PolymerPassFindExterns) GetPolymerElementProps() []rhino.Node {
	// ImmutableList.Builder.build() snapshots the elements added so far.
	return append([]rhino.Node(nil), f.polymerElementProps...)
}

// GetPolymerElementExterns returns the root Node PolymerElement class in
// externs.
//
// Java: PolymerPassFindExterns.getPolymerElementExterns()
func (f *PolymerPassFindExterns) GetPolymerElementExterns() rhino.Node {
	return f.polymerElementExterns
}

// polymerPassFindExterns_isPolymerElementExterns returns whether the node is
// the declaration of PolymerElement.
//
// Java: PolymerPassFindExterns.isPolymerElementExterns(Node)
func polymerPassFindExterns_isPolymerElementExterns(value rhino.Node) bool {
	// Explicitly allow var, let or const.
	decl, ok := ast.AsNameDeclaration(value)
	if !ok {
		return false
	}
	first, ok := decl.FirstDeclarator()
	return ok && first.Node().MatchesQualifiedName(polymerPassFindExterns_POLYMER_ELEMENT_NAME)
}

// polymerPassFindExterns_isPolymerElementPropExpr returns whether the node is
// an expression result of an assignment to a property of PolymerElement.
//
// Java: PolymerPassFindExterns.isPolymerElementPropExpr(Node)
func polymerPassFindExterns_isPolymerElementPropExpr(value rhino.Node) bool {
	stmt, ok := ast.AsExprStmt(value)
	if !ok {
		return false
	}
	first := stmt.Expr().Node().GetFirstChild()
	return first.IsGetProp() &&
		first.IsQualifiedName() &&
		NodeUtil_GetRootOfQualifiedName(first).MatchesQualifiedName(polymerPassFindExterns_POLYMER_ELEMENT_NAME)
}

// tagNameToTypeMap contains a mapping from HTML Element tag name to the
// javascript type of that element at runtime.
var tagNameToTypeMap = map[string]string{
	"a":        "HTMLAnchorElement",
	"area":     "HTMLAreaElement",
	"audio":    "HTMLAudioElement",
	"base":     "HTMLBaseElement",
	"body":     "HTMLBodyElement",
	"br":       "HTMLBRElement",
	"button":   "HTMLButtonElement",
	"canvas":   "HTMLCanvasElement",
	"caption":  "HTMLTableCaptionElement",
	"col":      "HTMLTableColElement",
	"content":  "HTMLContentElement",
	"data":     "HTMLDataElement",
	"datalist": "HTMLDataListElement",
	"del":      "HTMLModElement",
	"dialog":   "HTMLDialogElement",
	"dir":      "HTMLDirectoryElement",
	"div":      "HTMLDivElement",
	"dl":       "HTMLDListElement",
	"embed":    "HTMLEmbedElement",
	"fieldset": "HTMLFieldSetElement",
	"font":     "HTMLFontElement",
	"form":     "HTMLFormElement",
	"frame":    "HTMLFrameElement",
	"frameset": "HTMLFrameSetElement",
	"h1":       "HTMLHeadingElement",
	"head":     "HTMLHeadElement",
	"hr":       "HTMLHRElement",
	"html":     "HTMLHtmlElement",
	"iframe":   "HTMLIFrameElement",
	"img":      "HTMLImageElement",
	"input":    "HTMLInputElement",
	"keygen":   "HTMLKeygenElement",
	"label":    "HTMLLabelElement",
	"legend":   "HTMLLegendElement",
	"li":       "HTMLLIElement",
	"link":     "HTMLLinkElement",
	"map":      "HTMLMapElement",
	"marquee":  "HTMLMarqueeElement",
	"menu":     "HTMLMenuElement",
	"menuitem": "HTMLMenuItemElement",
	"meta":     "HTMLMetaElement",
	"meter":    "HTMLMeterElement",
	"object":   "HTMLObjectElement",
	"ol":       "HTMLOListElement",
	"optgroup": "HTMLOptGroupElement",
	"option":   "HTMLOptionElement",
	"output":   "HTMLOutputElement",
	"p":        "HTMLParagraphElement",
	"param":    "HTMLParamElement",
	"pre":      "HTMLPreElement",
	"progress": "HTMLProgressElement",
	"q":        "HTMLQuoteElement",
	"script":   "HTMLScriptElement",
	"select":   "HTMLSelectElement",
	"shadow":   "HTMLShadowElement",
	"source":   "HTMLSourceElement",
	"span":     "HTMLSpanElement",
	"style":    "HTMLStyleElement",
	"table":    "HTMLTableElement",
	"tbody":    "HTMLTableSectionElement",
	"tfoot":    "HTMLTableSectionElement",
	"template": "HTMLTemplateElement",
	"td":       "HTMLTableSectionElement",
	"textarea": "HTMLTextAreaElement",
	"th":       "HTMLTableCellElement",
	"thead":    "HTMLTableCellElement",
	"time":     "HTMLTimeElement",
	"title":    "HTMLTitleElement",
	"tr":       "HTMLTableRowElement",
	"track":    "HTMLTrackElement",
	"ul":       "HTMLUListElement",
	"video":    "HTMLVideoElement",
}

// Java: TagNameToType.getMap()
func TagNameToType_GetMap() map[string]string {
	return tagNameToTypeMap
}
