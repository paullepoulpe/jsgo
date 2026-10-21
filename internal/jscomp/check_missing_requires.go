// Port of com/google/javascript/jscomp/CheckMissingRequires.java @ closure-compiler v20261015

/*
 * Copyright 2020 The Closure Compiler Authors.
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

// CheckMissingRequires is a pass to detect references to fully qualified
// Closure namespaces.
type CheckMissingRequires struct {
	NodeTraversal_AbstractModuleCallbackBase

	// The set of template parameter names found so far in the file currently being checked.
	templateParamNames *jutil.LinkedHashSet[string]

	// The mapping from Closure namespace into the module that provides it.
	moduleByNamespace *jutil.LinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]

	// Tracks how many "control flow scopes" we've entered, starting at 0
	//
	// Where a "control flow scope" is here defined as any scope /except/ for either the global
	// control flow scope, or an IIFE or goog.scope body within the global control flow scope.
	controlFlowScopeDepth int
}

// NOTE: "references a .*namespace" is recognized by tooling.
//
// Java: CheckMissingRequires.MISSING_REQUIRE
var CheckMissingRequires_MISSING_REQUIRE = DiagnosticType_Warning(
	"JSC_MISSING_REQUIRE",
	"''{0}'' references a fully qualified namespace, which is disallowed by the style"+
		" guide.\nPlease add a goog.require, assign or destructure it into an alias, and "+
		"use the alias instead.\nNamespace defined in {1}.")

// Java: CheckMissingRequires.MISSING_REQUIRE_TYPE
var CheckMissingRequires_MISSING_REQUIRE_TYPE = DiagnosticType_Disabled(
	"JSC_MISSING_REQUIRE_TYPE",
	"''{0}'' references a fully qualified namespace, which is disallowed by the style"+
		" guide.\nPlease add a goog.requireType, assign or destructure it into an alias, "+
		"and use the alias instead.\nNamespace defined in {1}.")

// Java: CheckMissingRequires.INCORRECT_NAMESPACE_ALIAS_REQUIRE
var CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE = DiagnosticType_Disabled(
	"JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE",
	"''{0}'' is its own namespace.\nPlease add a separate goog.require and "+
		"use that alias instead.")

// Java: CheckMissingRequires.INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE
var CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE = DiagnosticType_Disabled(
	"JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE",
	"''{0}'' is its own namespace.\nPlease add a separate goog.requireType and "+
		"use that alias instead.")

// Java: CheckMissingRequires.INDIRECT_NAMESPACE_REF_REQUIRE
var CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE = DiagnosticType_Disabled(
	"JSC_INDIRECT_NAMESPACE_REF_REQUIRE",
	"''{0}'' should have its own goog.require.\nPlease add a separate goog.require and "+
		"use that alias instead.")

// Java: CheckMissingRequires.INDIRECT_NAMESPACE_REF_REQUIRE_TYPE
var CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE = DiagnosticType_Disabled(
	"JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE",
	"''{0}'' should have its own goog.requireType.\n"+
		"Please add a separate goog.requireType and use that alias instead.")

// Java: CheckMissingRequires.MISSING_REQUIRE_IN_PROVIDES_FILE
var CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE = DiagnosticType_Warning(
	"JSC_MISSING_REQUIRE_IN_PROVIDES_FILE",
	"''{0}'' references a namespace which was not required by this file.\n"+
		"Please add a goog.require.\nNamespace defined in {1}.")

// Java: CheckMissingRequires.MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE
var CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE = DiagnosticType_Disabled(
	"JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE",
	"''{0}'' references a namespace which was not required by this file.\n"+
		"Please add a goog.requireType.\nNamespace defined in {1}.")

// Java: CheckMissingRequires.NON_LEGACY_GOOG_MODULE_REFERENCE
var CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE = DiagnosticType_Error(
	"JSC_NON_LEGACY_GOOG_MODULE_REFERENCE",
	"''{0}'' references the name of a module without goog.declareLegacyNamespace(), which "+
		"is not actually defined. Use goog.module.get() instead.")

// Java: CheckMissingRequires.MISSING_REQUIRE_IN_GOOG_SCOPE
var CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE = DiagnosticType_Warning(
	"JSC_MISSING_REQUIRE_IN_GOOG_SCOPE",
	"''{0}'' is its own namespace. Please add a goog.require and reference that alias.")

// Java: CheckMissingRequires.MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE
var CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE = DiagnosticType_Warning(
	"JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE ",
	"''{0}'' is its own namespace. Please add a goog.requireType and reference that alias.")

// Java: CheckMissingRequires.MISSING_REQUIRE_FOR_GOOG_MODULE_GET
var CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET = DiagnosticType_Warning(
	"JSC_MISSING_REQUIRE_FOR_GOOG_MODULE_GET",
	"''{0}'' references a namespace which was not required by this file.\n"+
		"Please add a goog.require.")

// Java: CheckMissingRequires(AbstractCompiler, ModuleMetadataMap)
func NewCheckMissingRequires(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *CheckMissingRequires {
	c := &CheckMissingRequires{}
	c.InitNodeTraversal_AbstractModuleCallbackBase(c, compiler, moduleMetadataMap)
	c.templateParamNames = jutil.NewLinkedHashSet[string]()
	c.controlFlowScopeDepth = 0
	c.moduleByNamespace = moduleMetadataMap.GetModulesByGoogNamespace()
	return c
}

// Java: CheckMissingRequires.process(Node, Node)
func (c *CheckMissingRequires) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckMissingRequires.shouldTraverse(NodeTraversal, Node, ModuleMetadata, Node)
func (c *CheckMissingRequires) ShouldTraverseWithModule(
	t *NodeTraversal, n *rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata, scopeRoot *rhino.Node) bool {
	if NodeUtil_IsValidCfgRoot(n) && !c.IsInGlobalControlFlowScope(n) {
		c.controlFlowScopeDepth++
	}
	if currentModule == nil {
		return true
	}

	// Traverse nodes in preorder to collect `@template` parameter names before their use.
	c.VisitNode(t, n, currentModule)
	return true
}

// IsInGlobalControlFlowScope takes a valid "control flow root", as defined by
// NodeUtil.isValidCfgRoot, and returns whether it creates a new "control flow
// scope" from the parent scope (the global scope).
//
// Java: CheckMissingRequires.isInGlobalControlFlowScope(Node)
func (c *CheckMissingRequires) IsInGlobalControlFlowScope(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_ROOT:
		return true
	case rhino.Token_MODULE_BODY, rhino.Token_BLOCK:
		// class static blocks need to create a new cfg root
		return false
	case rhino.Token_FUNCTION:
		if c.controlFlowScopeDepth > 0 {
			return false
		}
		// Check for functions invoked immediately within the global scope - so IIFEs or
		// goog.scope callees.
		if CheckMissingRequires_IsGoogScopeBody(NodeUtil_GetFunctionBody(n)) {
			return true
		}
		return CheckMissingRequires_IsIIFE(n)
	default:
		panic(jutil.NewAssertionError("Unexpected control flow root: " + n.String()))
	}
}

// Java: CheckMissingRequires.isIIFE(Node)
func CheckMissingRequires_IsIIFE(n *rhino.Node) bool {
	return n.IsFunction() && n.GetParent().IsCall() && n.IsFirstChildOf(n.GetParent())
}

// Java: CheckMissingRequires.visit(NodeTraversal, Node, ModuleMetadata, Node)
func (c *CheckMissingRequires) VisitWithModule(
	t *NodeTraversal, n *rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata, scopeRoot *rhino.Node) {
	if NodeUtil_IsValidCfgRoot(n) && !c.IsInGlobalControlFlowScope(n) {
		c.controlFlowScopeDepth--
	}
	if currentModule != nil && n == currentModule.RootNode() {
		// For this pass, template parameter names are only meaningful inside the file defining them.
		c.templateParamNames.Clear()
	}
}

// Java: CheckMissingRequires.visitNode(NodeTraversal, Node, ModuleMetadata)
func (c *CheckMissingRequires) VisitNode(t *NodeTraversal, n *rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata) {
	info := n.GetJSDocInfo()
	if info != nil {
		c.VisitJsDocInfo(t, currentModule, info)
	}
	if !n.GetParent().IsGetProp() && n.IsQualifiedName() {
		qualifiedName := n.GetQualifiedNameObject()
		root := qualifiedName.GetRoot()
		if root == "this" || root == "super" {
			return
		}
		c.VisitQualifiedName(t, n, currentModule, qualifiedName, CheckMissingRequires_Strength_CODE)
	}

	if n.IsName() && n.GetString() != "" {
		c.VisitMaybeDeclaration(t, n, currentModule)
	}

	if n.IsCall() && c.controlFlowScopeDepth == 0 {
		c.VisitMaybeGoogModuleGet(t, n, currentModule)
	}
}

// Java: CheckMissingRequires.visitJsDocInfo(NodeTraversal, ModuleMetadata, JSDocInfo)
func (c *CheckMissingRequires) VisitJsDocInfo(t *NodeTraversal, currentModule *ModuleMetadataMap_ModuleMetadata, info *rhino.JSDocInfo) {
	// Collect template parameter names before checking, so that annotations on the same node that
	// reference the name are excluded from the check.
	c.templateParamNames.AddAll(info.GetTemplateTypeNames()...)
	for name := range info.GetTypeTransformations().Keys() {
		c.templateParamNames.Add(name)
	}
	if info.HasType() {
		c.VisitJsDocExpr(t, currentModule, info.GetType(), CheckMissingRequires_Strength_WEAK_TYPE)
	}
	for paramType := range info.ParameterTypes() {
		if paramType != nil {
			c.VisitJsDocExpr(t, currentModule, paramType, CheckMissingRequires_Strength_WEAK_TYPE)
		}
	}
	if info.HasReturnType() {
		c.VisitJsDocExpr(t, currentModule, info.GetReturnType(), CheckMissingRequires_Strength_WEAK_TYPE)
	}
	if info.HasEnumParameterType() {
		c.VisitJsDocExpr(t, currentModule, info.GetEnumParameterType(), CheckMissingRequires_Strength_WEAK_TYPE)
	}
	if info.HasTypedefType() {
		c.VisitJsDocExpr(t, currentModule, info.GetTypedefType(), CheckMissingRequires_Strength_WEAK_TYPE)
	}
	if info.HasThisType() {
		c.VisitJsDocExpr(t, currentModule, info.GetThisType(), CheckMissingRequires_Strength_WEAK_TYPE)
	}
	if info.HasBaseType() {
		// Note that `@extends` requires a goog.require, not a goog.requireType.
		c.VisitJsDocExpr(t, currentModule, info.GetBaseType(), CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS)
	}
	for _, expr := range info.GetExtendedInterfaces() {
		// Note that `@extends` requires a goog.require, not a goog.requireType.
		c.VisitJsDocExpr(t, currentModule, expr, CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS)
	}
	for _, expr := range info.GetImplementedInterfaces() {
		// Note that `@implements` requires a goog.require, not a goog.requireType.
		c.VisitJsDocExpr(t, currentModule, expr, CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS)
	}
}

// Java: CheckMissingRequires.visitJsDocExpr(NodeTraversal, ModuleMetadata, JSTypeExpression, Strength)
func (c *CheckMissingRequires) VisitJsDocExpr(
	t *NodeTraversal,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	expr *rhino.JSTypeExpression,
	referenceStrength CheckMissingRequires_Strength) {
	for _, typeNode := range expr.GetAllTypeNodes() {
		c.VisitQualifiedName(
			t, typeNode, currentModule, rhino.QualifiedName_Of(typeNode.GetString()), referenceStrength)
	}
}

// Java: CheckMissingRequires.GOOG_SCOPE
var checkMissingRequires_GOOG_SCOPE = rhino.QualifiedName_Of("goog.scope")

// Java: CheckMissingRequires.isGoogScopeBody(Node)
func CheckMissingRequires_IsGoogScopeBody(hoistScopeRoot *rhino.Node) bool {
	return hoistScopeRoot.IsBlock() &&
		hoistScopeRoot.GetParent().IsFunction() &&
		hoistScopeRoot.GetGrandparent().IsCall() &&
		checkMissingRequires_GOOG_SCOPE.Matches(hoistScopeRoot.GetGrandparent().GetFirstChild())
}

// VisitMaybeGoogModuleGet checks for invalid goog.module.get calls executed on
// script load.
//
// We don't check for goog.module.get within function bodies though (except
// for goog.scope since that's basically an IIFE) - it's possible that the
// module will have been loaded by the time the goog.module.get is called even
// though there's no strong require. It's up to the caller to ensure that it's
// loaded.
//
// Java: CheckMissingRequires.visitMaybeGoogModuleGet(NodeTraversal, Node, ModuleMetadata)
func (c *CheckMissingRequires) VisitMaybeGoogModuleGet(
	t *NodeTraversal, googModuleGet *rhino.Node, currentFile *ModuleMetadataMap_ModuleMetadata) {
	jutil.CheckState(googModuleGet.IsCall())
	if !NodeUtil_IsGoogModuleGetCall(googModuleGet) {
		return
	}
	importedNamespace := googModuleGet.GetSecondChild().GetString()
	requiredFile := c.moduleByNamespace.Get(importedNamespace)
	if requiredFile == nil {
		// goog.module.get of a non-existing namespace is an error in another pass.
		return
	}
	if !CheckMissingRequires_HasAcceptableRequire(
		currentFile, rhino.QualifiedName_Of(importedNamespace), requiredFile, CheckMissingRequires_Strength_WEAK_TYPE) {
		t.Report(googModuleGet, CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET, importedNamespace)
	}
}

// Java: CheckMissingRequires.visitMaybeDeclaration(NodeTraversal, Node, ModuleMetadata)
func (c *CheckMissingRequires) VisitMaybeDeclaration(t *NodeTraversal, n *rhino.Node, currentFile *ModuleMetadataMap_ModuleMetadata) {
	if !currentFile.IsModule() {
		// This check only makes sense in goog.module files
		return
	}

	// Currently, tsickle can introduce these
	// TODO(b/333952917): Remove this once tsickle is fixed.
	if c.IsTypeScriptSource(n) {
		return
	}

	jutil.CheckState(n.IsName())
	rootName := n.GetString()
	var_ := t.GetScope().GetVar(rootName)
	if var_ == nil {
		return
	}
	declaration := var_.GetDeclaration()
	if declaration == nil || declaration.GetNode() != n {
		// Not a declaration reference
		return
	}

	if var_.GetScope().(*Scope).IsLocal() {
		require := NodeUtil_GetGoogRequireInfo(var_)
		if require == nil {
			// If it refers to a local variable, it's legal.
			return
		}

		// Not destructured.
		if require.Property() == nil {
			return
		}

		requiredFile := c.moduleByNamespace.Get(require.Namespace())
		if requiredFile == nil || !requiredFile.HasLegacyGoogNamespaces() {
			// If we are importing from a non-legacy namespace than any property reference is valid,
			// as it is assumed that the property reference is exported intentionally.

			// NOTE There are corner cases if a legacy namespace object is explicitly reexported
			// but (1) it isn't clear what the "correct" behavior is and (2) we don't have the
			// information here in order to validate it as we would need to the object shapes.
			return
		}

		reference := require.Namespace() + "." + *require.Property()

		// Do not report references to a namespace provided in the same file, and do not recurse
		// into parent namespaces either.
		if !currentFile.GoogNamespaces().Contains(reference) {
			// Try to find an alternate file to import from if, if we construct a namespace using
			// the destructured property name.
			alternateFile := c.moduleByNamespace.Get(reference)

			if alternateFile != nil && alternateFile != requiredFile {
				// TODO: report on the node that needs to be removed, include the namespace that needs
				// to be added.
				toReport := CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE
				if require.IsStrongRequire() {
					toReport = CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE
				}
				t.Report(n, toReport, reference)
			}
		}
	}
}

// Java: CheckMissingRequires.isTypeScriptSource(Node)
func (c *CheckMissingRequires) IsTypeScriptSource(n *rhino.Node) bool {
	return n.GetStaticSourceFile().IsTypeScriptSource()
}

// Java: CheckMissingRequires.visitQualifiedName(NodeTraversal, Node, ModuleMetadata, QualifiedName, Strength)
func (c *CheckMissingRequires) VisitQualifiedName(
	t *NodeTraversal,
	n *rhino.Node,
	currentFile *ModuleMetadataMap_ModuleMetadata,
	qualifiedName rhino.QualifiedName,
	referenceStrength CheckMissingRequires_Strength) {

	rootName := qualifiedName.GetRoot()
	if qualifiedName.IsSimple() {
		if c.templateParamNames.Contains(rootName) {
			// This will produce a false negative when the same name is used in both template and
			// non-template capacity in the same file, and a false positive when the `@template` does
			// not precede the reference within the same source file (e.g. an ES5 ctor in a different
			// file).
			return
		}
	}

	if NodeUtil_IsDeclarationLValue(n) {
		// NOTE: maybe do something here.
		return
	}

	var_ := t.GetScope().GetVar(rootName)
	if var_ != nil && var_.GetScope().(*Scope).IsLocal() {
		c.CheckMissingRequireThroughShortName(t, n, var_, currentFile, qualifiedName, referenceStrength)
	} else {
		c.CheckMissingRequireThroughFullyQualifiedName(
			t, n, currentFile, qualifiedName, referenceStrength)
	}
}

// CheckMissingRequireThroughShortName checks if a reference to a local
// variable should actually trigger a missing require warning.
//
// This only warns within goog.scope and module bodies, in cases where there's
// a local variable that's actually an alias of some import.
//
// n is the node representing the fully qualified name being checked,
// starting with localVar; localVar the root of the qualified name;
// qualifiedName some fully qualified name to check, starting with localVar.
//
// Java: CheckMissingRequires.checkMissingRequireThroughShortName(NodeTraversal, Node, Var, ModuleMetadata, QualifiedName, Strength)
func (c *CheckMissingRequires) CheckMissingRequireThroughShortName(
	t *NodeTraversal,
	n *rhino.Node,
	localVar *Var,
	currentFile *ModuleMetadataMap_ModuleMetadata,
	qualifiedName rhino.QualifiedName,
	referenceStrength CheckMissingRequires_Strength) {
	// TODO(b/333952917): Remove this once tsickle is fixed.
	if c.IsTypeScriptSource(n) {
		return
	}

	inModule := currentFile.IsModule()
	inGoogScope := !inModule && CheckMissingRequires_IsGoogScopeBody(localVar.GetScopeRoot())

	if !inModule && !inGoogScope {
		// Don't worry about aliases outside of module files and goog.scopes for now.
		return
	}

	// The qualified name *is* the root name if it is "simple"
	if qualifiedName.IsSimple() {
		// It is explicitly imported and we have already validated the import in
		// `visitMaybeDeclaration`
		return
	}

	if inGoogScope {
		c.CheckMissingRequireThroughShortNameInGoogScope(
			t, n, localVar, currentFile, qualifiedName, referenceStrength)
	} else {
		c.CheckMissingRequireThroughShortNameInModule(
			t,
			n,
			NodeUtil_GetGoogRequireInfo(localVar),
			currentFile,
			qualifiedName,
			referenceStrength)
	}
}

// Java: CheckMissingRequires.checkMissingRequireThroughShortNameInGoogScope(NodeTraversal, Node, Var, ModuleMetadata, QualifiedName, Strength)
func (c *CheckMissingRequires) CheckMissingRequireThroughShortNameInGoogScope(
	t *NodeTraversal,
	n *rhino.Node,
	localVar *Var,
	currentFile *ModuleMetadataMap_ModuleMetadata,
	qualifiedName rhino.QualifiedName,
	referenceStrength CheckMissingRequires_Strength) {
	initialValue := localVar.GetInitialValue()
	if initialValue == nil {
		return
	}
	// There are two cases we have to handle (assume 'foo.bar' is a module or provide, below);
	//  1) a legacy namespace, possibly with extra properties:
	//        const alias = foo.bar;
	//        const alias = foo.bar.baz;
	//  2) a goog.module.get, possibly with extra properties:
	//      const alias = goog.module.get('foo.bar');
	//      const alias = goog.module.get('foo.bar').baz;
	// Ignore any other initial value. This pass isn't designed to check arbitrary chains of import
	// aliases, and we just look at direct aliases of Closure namespaces. For example, we will not
	// find aliases such as:
	//    alias.something = foo.bar;
	// This aligns with the behavior of the ScopedAliases pass goog.scope handling.

	// Case (1) - the initial value is a legacy namespace + potential extra properties
	if initialValue.IsQualifiedName() {
		aliasedName := initialValue.GetQualifiedNameObject()
		normalizedQualifiedName := CheckMissingRequires_SwapRootOfQualifiedName(qualifiedName, aliasedName)

		// Find longest prefix match that is actually a provide/legacy namespace, if any.
		var originalImport *ModuleMetadataMap_ModuleMetadata = nil
		for subName := aliasedName; subName != nil; subName = subName.GetOwner() {
			imported := c.moduleByNamespace.Get(subName.Join())
			if imported != nil {
				originalImport = imported
				break
			}
		}
		c.CheckNormalizedShortNameImport(
			t,
			n,
			currentFile,
			originalImport,
			qualifiedName,
			normalizedQualifiedName,
			referenceStrength)
		return
	}

	// Case (2) - the initial value is a goog.module.get + potential extra properties
	for initialValue.IsGetProp() {
		initialValue = initialValue.GetFirstChild()
	}
	if !NodeUtil_IsGoogModuleGetCall(initialValue) {
		return
	}
	importedNamespace := initialValue.GetSecondChild().GetString()
	originalImport := c.moduleByNamespace.Get(importedNamespace)
	if originalImport != nil && !originalImport.HasLegacyGoogNamespaces() {
		// We trust the import from a non-legacy module, as they should not overlap.
		return
	}
	aliasedName := rhino.QualifiedName_Of(importedNamespace)
	for prop := initialValue.GetParent(); prop.IsGetProp(); prop = prop.GetParent() {
		aliasedName = aliasedName.Getprop(prop.GetString())
	}
	normalizedQualifiedName := CheckMissingRequires_SwapRootOfQualifiedName(qualifiedName, aliasedName)
	c.CheckNormalizedShortNameImport(
		t,
		n,
		currentFile,
		originalImport,
		qualifiedName,
		normalizedQualifiedName,
		referenceStrength)
}

// Java: CheckMissingRequires.checkMissingRequireThroughShortNameInModule(NodeTraversal, Node, NodeUtil.GoogRequire, ModuleMetadata, QualifiedName, Strength)
func (c *CheckMissingRequires) CheckMissingRequireThroughShortNameInModule(
	t *NodeTraversal,
	n *rhino.Node,
	require *NodeUtil_GoogRequire,
	currentFile *ModuleMetadataMap_ModuleMetadata,
	qualifiedName rhino.QualifiedName,
	referenceStrength CheckMissingRequires_Strength) {
	if require == nil {
		// It is a local name, not an import.
		// NOTE: this *could be a local alias* of an imported name
		//   but at some point the only real fix is to not use `goog.provide`
		//   or `goog.module.declareLegacyNamespace` because it is always
		//   possible to obfuscate the use of the namespace.
		return
	}

	originalRequiredFile := c.moduleByNamespace.Get(require.Namespace())
	if originalRequiredFile == nil || !originalRequiredFile.HasLegacyGoogNamespaces() {
		// We trust the import from a non-legacy module, as they should not overlap.
		return
	}

	// Verify namespace usage
	normalizeQualifiedName := c.NormalizeQualifiedNamePlusImport(qualifiedName, require)

	c.CheckNormalizedShortNameImport(
		t,
		n,
		currentFile,
		originalRequiredFile,
		qualifiedName,
		normalizeQualifiedName,
		referenceStrength)
}

// Java: CheckMissingRequires.checkNormalizedShortNameImport(NodeTraversal, Node, ModuleMetadata, ModuleMetadata, QualifiedName, QualifiedName, Strength)
func (c *CheckMissingRequires) CheckNormalizedShortNameImport(
	t *NodeTraversal,
	n *rhino.Node,
	currentFile *ModuleMetadataMap_ModuleMetadata,
	originalImport *ModuleMetadataMap_ModuleMetadata,
	originalQualifiedName rhino.QualifiedName,
	normalizedQualifiedName rhino.QualifiedName,
	referenceStrength CheckMissingRequires_Strength) {
	// TESTCASES, where A.B should not be used through A:
	// checked here
	// X  const A = require('A'); ref(A.B);  // error, missing require A.B
	// X  const A = require('A'); const B = require('B');   ref(A.B);  // error, incorrect ref
	// checked in declaration
	// X  const {B} = require('A');  // error bad import
	//

	// Look for the longest prefix match against a provided namespace.
	for subName := normalizedQualifiedName; subName != nil; subName = subName.GetOwner() {
		namespace := subName.Join()
		if c.IsAllowedNamespace(currentFile, namespace) {
			return
		}

		requiredFile := c.moduleByNamespace.Get(namespace)
		if requiredFile == nil {
			// Not a known namespace check the parent
			continue
		}

		// TODO: report on the node that needs rewritten, include the namespace that needs
		// to be added:
		// Autofix: add the require if necessary, rewrite to be the full namespace, and let fixjs
		// fix
		// it up. Write a different message if the namespace is already imported vs missing.

		if !CheckMissingRequires_HasAcceptableRequire(currentFile, subName, requiredFile, referenceStrength) ||
			originalImport == nil ||
			originalImport != requiredFile {
			// Why check `originalImport == null|| originalImport != requiredFile`? in case the imported
			// file does have an acceptable goog.require in this file, but this particular reference is
			// still through the wrong namespace.
			// For example:
			//   const foo = goog.require('my.foo');
			//   const bar = goog.require('my.foo.bar');
			// (...)
			//    use(foo.bar); // should be use(bar);
			var toReport *DiagnosticType
			if currentFile.IsModule() {
				if referenceStrength.IsStrong() {
					toReport = CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE
				} else {
					toReport = CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE
				}
			} else {
				if referenceStrength.IsStrong() {
					toReport = CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE
				} else {
					toReport = CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE
				}
			}

			errorNode := n
			// Get the right location if it is a code reference.
			diff := normalizedQualifiedName.GetComponentCount() - subName.GetComponentCount()
			if n.IsQualifiedName() {
				for i := 0; i < diff && errorNode.IsGetProp(); i++ {
					errorNode = errorNode.GetFirstChild()
				}
				t.ReportWithNodeRange(NodeUtil_GetRootOfQualifiedName(errorNode), errorNode, toReport, namespace)
			} else {
				// JSDoc reference case
				// Trim the original qualified name from the source text by the same difference.
				errorQName := originalQualifiedName
				for i := 0; i < diff; i++ {
					errorQName = errorQName.GetOwner()
				}

				correctName := errorQName.Join()
				// For JSDoc, the entire qualified name is a single STRINGLIT node.
				// To report the error on the correct sub-expression, we create a new
				// temporary node with the correct string and length.
				newErrorNode := n.CloneNode()
				newErrorNode.SetString(correctName)
				newErrorNode.SetLength(jutil.StringLength(correctName))

				t.Report(newErrorNode, toReport, namespace)
			}
		}

		// We found the imported namespace: done
		return
	}
}

// CheckMissingRequireThroughFullyQualifiedName checks if a reference to a
// global qualified name should trigger a missing require warning.
//
// Here, "global qualified name" means that the root object for the name is not
// defined in some local module or function scope, but is global.
//
// n is the node representing the fully qualified name being checked;
// qualifiedName some global fully qualified name to check; referenceStrength
// whether this is a non-type-only reference.
//
// Java: CheckMissingRequires.checkMissingRequireThroughFullyQualifiedName(NodeTraversal, Node, ModuleMetadata, QualifiedName, Strength)
func (c *CheckMissingRequires) CheckMissingRequireThroughFullyQualifiedName(
	t *NodeTraversal,
	n *rhino.Node,
	currentFile *ModuleMetadataMap_ModuleMetadata,
	qualifiedName rhino.QualifiedName,
	referenceStrength CheckMissingRequires_Strength) {

	// Look for the longest prefix match against a provided namespace.
	for subName := qualifiedName; subName != nil; subName = subName.GetOwner() {
		namespace := subName.Join()
		if c.IsAllowedNamespace(currentFile, namespace) {
			return
		}

		requiredFile := c.moduleByNamespace.Get(namespace)
		if requiredFile == nil {
			// Not a known namespace check the parent
			continue
		}
		fileName := requiredFile.RootNode().GetSourceFileName()
		var toReport *DiagnosticType
		if currentFile.IsModule() {
			/*
			 * In files that represent modules, report a require without an alias the same as a totally
			 * missing require.
			 */
			if referenceStrength.IsStrong() {
				toReport = CheckMissingRequires_MISSING_REQUIRE
			} else {
				toReport = CheckMissingRequires_MISSING_REQUIRE_TYPE
			}
		} else if !CheckMissingRequires_HasAcceptableRequire(currentFile, subName, requiredFile, referenceStrength) {
			/*
			 * In files that aren't modules, report a qualified name reference only if there's no
			 * require to satisfy it.
			 */
			if referenceStrength.IsStrong() {
				toReport = CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE
			} else {
				toReport = CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE
			}
		} else if !referenceStrength.IsTypeOnly() && requiredFile.IsNonLegacyGoogModule() {
			// The referenced file has a valid goog.require for it, but the reference cannot be in a
			// regular qualified name: that won't be defined at runtime.
			// However, this only applies to references in actual code. Type-only references may still
			// reference the namespace in a script. The Closure type system will still resolve the
			// namespace - it just doesn't have any runtime effect.
			toReport = CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE
		} else {
			return
		}

		t.Report(n, toReport, namespace, jutil.ToJavaString(fileName))
		return
	}
}

// NormalizeQualifiedNamePlusImport: given a qualified name with local root
// name and an import that defines that name, transform that name so that it
// is a fully qualified name.
//
// Java: CheckMissingRequires.normalizeQualifiedNamePlusImport(QualifiedName, NodeUtil.GoogRequire)
func (c *CheckMissingRequires) NormalizeQualifiedNamePlusImport(
	qualifiedName rhino.QualifiedName, imported *NodeUtil_GoogRequire) rhino.QualifiedName {
	newQName := rhino.QualifiedName_Of(imported.Namespace())

	if imported.Property() != nil {
		newQName = newQName.Getprop(*imported.Property())
	}
	return CheckMissingRequires_SwapRootOfQualifiedName(qualifiedName, newQName)
}

// CheckMissingRequires_SwapRootOfQualifiedName replaces the root name of the
// original qname with the newRoot (which may be a simple or complex qualified
// name).
//
// For example: if original is `oldRoot.a.b.c`, and newRoot is
// `newRoot.withProperty`, returns `newRoot.withProperty.a.b.c`.
//
// Java: CheckMissingRequires.swapRootOfQualifiedName(QualifiedName, QualifiedName)
func CheckMissingRequires_SwapRootOfQualifiedName(
	original rhino.QualifiedName, newRoot rhino.QualifiedName) rhino.QualifiedName {
	newQName := newRoot
	components := original.Components()
	if len(components) == 0 {
		panic(jutil.NewNoSuchElementException()) // Iterator.next() skipping the root
	}
	for _, component := range components[1:] { // skip the root
		newQName = newQName.Getprop(component)
	}
	// Verify namespace usage
	return newQName
}

// Java: CheckMissingRequires.isAllowedNamespace(ModuleMetadata, String)
func (c *CheckMissingRequires) IsAllowedNamespace(currentFile *ModuleMetadataMap_ModuleMetadata, namespace string) bool {
	if namespace == "goog.module" {
		// We must special case `goog.module` because Closure Library provides a namespace with that
		// name, but it's (confusingly) unrelated to the `goog.module` primitive.
		return true
	}

	// Do not report references to a namespace provided in the same file, and do not recurse
	// into parent namespaces either.
	// TODO(tjgq): Also check for these references.
	if currentFile.GoogNamespaces().Contains(namespace) {
		return true
	}

	return false
}

// CheckMissingRequires_HasAcceptableRequire: does `rdep` contain an
// acceptable require for `namespace` from `dep`?
//
// Any require for a parent of `namespace` from `rdep` onto `dep`, which
// declares that parent, is sufficient. This constraint still ensures correct
// dependency ordering.
//
// We loosen the check in this way because we want to make it easy to migrate
// multi-provide files into modules. That includes deleting obsolete provides
// and splitting provides into different files.
//
// Java: CheckMissingRequires.hasAcceptableRequire(ModuleMetadata, QualifiedName, ModuleMetadata, Strength)
func CheckMissingRequires_HasAcceptableRequire(
	rdep *ModuleMetadataMap_ModuleMetadata,
	namespace rhino.QualifiedName,
	dep *ModuleMetadataMap_ModuleMetadata,
	referenceStrength CheckMissingRequires_Strength) bool {
	// Go-only: language-adapter — the Java code builds Guava Sets.union and
	// Sets.intersection views of the multisets' element sets and only ever
	// queries them with contains(); acceptableRequires is that membership test.
	acceptableRequires := func(e string) bool {
		inRequires := rdep.StronglyRequiredGoogNamespaces().Contains(e) ||
			(!referenceStrength.IsStrong() && rdep.WeaklyRequiredGoogNamespaces().Contains(e))
		return inRequires && dep.GoogNamespaces().Contains(e)
	}

	for parent := namespace; parent != nil; parent = parent.GetOwner() {
		if acceptableRequires(parent.Join()) {
			return true
		}
	}

	return false
}

// CheckMissingRequires_Strength represents the strength of references of an
// require'd name, i.e. whether references are in code only (strong,
// not-typeOnly), JsDoc annotations for @implements or @extends (strong as well
// as typeOnly) or in other JsDoc annotations only (weak, typeOnly).
//
// Java: CheckMissingRequires.Strength
type CheckMissingRequires_Strength int32

const (
	CheckMissingRequires_Strength_CODE CheckMissingRequires_Strength = iota
	CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS
	CheckMissingRequires_Strength_WEAK_TYPE
)

var checkMissingRequires_Strength_names = [...]string{"CODE", "IMPLEMENTS_EXTENDS", "WEAK_TYPE"}

var checkMissingRequires_Strength_isStrong = [...]bool{true, true, false}

var checkMissingRequires_Strength_isTypeOnly = [...]bool{false, true, true}

// Java: CheckMissingRequires.Strength.values()
func CheckMissingRequires_Strength_values() []CheckMissingRequires_Strength {
	return []CheckMissingRequires_Strength{
		CheckMissingRequires_Strength_CODE,
		CheckMissingRequires_Strength_IMPLEMENTS_EXTENDS,
		CheckMissingRequires_Strength_WEAK_TYPE,
	}
}

// String is Enum.name().
func (s CheckMissingRequires_Strength) String() string { return checkMissingRequires_Strength_names[s] }

// Ordinal is Enum.ordinal().
func (s CheckMissingRequires_Strength) Ordinal() int { return int(s) }

// IsStrong is the field Strength.isStrong.
func (s CheckMissingRequires_Strength) IsStrong() bool {
	return checkMissingRequires_Strength_isStrong[s]
}

// IsTypeOnly is the field Strength.isTypeOnly.
func (s CheckMissingRequires_Strength) IsTypeOnly() bool {
	return checkMissingRequires_Strength_isTypeOnly[s]
}
