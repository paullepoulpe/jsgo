// Port of com/google/javascript/jscomp/PolymerBehaviorExtractor.java @ closure-compiler v20261015

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var polymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY = jutil.NewLinkedHashSet(
	"created",
	"attached",
	"detached",
	"attributeChanged",
	"configure",
	"ready",
	"properties",
	"listeners",
	"observers",
	"hostAttributes")

const polymerBehaviorExtractor_GOOG_MODULE_EXPORTS = "exports"

// PolymerBehaviorExtractor finds the Polymer behavior definitions associated
// with Polymer element definitions.
//
// See https://www.polymer-project.org/1.0/docs/devguide/behaviors
//
// Java: final class PolymerBehaviorExtractor
type PolymerBehaviorExtractor struct {
	compiler          AbstractCompiler
	globalNames       *GlobalNamespace
	moduleMetadataMap *ModuleMetadataMap
	moduleMap         *ModuleMap

	// Java HashBasedTable: LinkedHashMap-backed, and only looked up here.
	resolveMemoized       *jutil.HashBasedTable[string, *ModuleMetadataMap_ModuleMetadata, *PolymerBehaviorExtractor_ResolveBehaviorNameResult]
	globalResolveMemoized *jutil.LinkedHashMap[string, *PolymerBehaviorExtractor_ResolveBehaviorNameResult]
}

// Java: PolymerBehaviorExtractor(AbstractCompiler, GlobalNamespace, ModuleMetadataMap, ModuleMap)
func NewPolymerBehaviorExtractor(
	compiler AbstractCompiler,
	globalNames *GlobalNamespace,
	moduleMetadataMap *ModuleMetadataMap,
	moduleMap *ModuleMap) *PolymerBehaviorExtractor {
	e := &PolymerBehaviorExtractor{}
	e.resolveMemoized = jutil.HashBasedTable_Create[string, *ModuleMetadataMap_ModuleMetadata, *PolymerBehaviorExtractor_ResolveBehaviorNameResult]()
	e.globalResolveMemoized = jutil.NewLinkedHashMap[string, *PolymerBehaviorExtractor_ResolveBehaviorNameResult]()
	e.compiler = compiler
	e.globalNames = globalNames
	e.moduleMetadataMap = moduleMetadataMap
	e.moduleMap = moduleMap
	return e
}

// ExtractBehaviors extracts all Behaviors from an array literal, recursively.
// Entries in the array can be object literals or array literals (of other
// behaviors). Behavior names must be global, fully qualified names.
// TODO(rishipal): Make this function better handle case where the same
// behavior transitively gets included in the same Polymer element more than
// once.
//
// See https://github.com/Polymer/polymer/blob/0.8-preview/PRIMER.md#behaviors
//
// moduleMetadata is the module in which these behaviors are being resolved,
// or null if not in a module. Returns a list of all BehaviorDefinitions in
// the array.
//
// Java: PolymerBehaviorExtractor.extractBehaviors(Node, ModuleMetadata)
func (e *PolymerBehaviorExtractor) ExtractBehaviors(
	behaviorArray rhino.Node, moduleMetadata *ModuleMetadataMap_ModuleMetadata) []*PolymerBehaviorExtractor_BehaviorDefinition {
	if behaviorArray.IsNil() {
		return nil
	}

	if !behaviorArray.IsArrayLit() {
		e.compiler.Report(
			JSError_MakeWithNode(behaviorArray, PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY))
		return nil
	}

	var behaviors []*PolymerBehaviorExtractor_BehaviorDefinition
	for behaviorName := behaviorArray.GetFirstChild(); !behaviorName.IsNil(); behaviorName = behaviorName.GetNext() {
		if behaviorName.IsObjectLit() {
			PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(behaviorName, e.compiler)
			PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(behaviorName, e.compiler)
			PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(behaviorName)
			if !NodeUtil_GetFirstPropMatchingKey(behaviorName, "is").IsNil() {
				e.compiler.Report(JSError_MakeWithNode(behaviorName, PolymerPassErrors_POLYMER_INVALID_BEHAVIOR))
			}
			behaviorModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorName)
			features, _ := NodeUtil_GetEnclosingScript(behaviorName).GetProp(rhino.Node_FEATURE_SET).(*parser.FeatureSet)
			behaviors = append(behaviors,
				NewPolymerBehaviorExtractor_BehaviorDefinition(
					PolymerPassStaticUtils_ExtractProperties(
						behaviorName,
						PolymerClassDefinition_DefinitionType_ObjectLiteral,
						e.compiler,
						/* constructor= */ rhino.Node{}),
					polymerBehaviorExtractor_getBehaviorFunctionsToCopy(behaviorName),
					polymerBehaviorExtractor_getNonPropertyMembersToCopy(behaviorName),
					/* isGlobalDeclaration= */ NodeUtil_GetEnclosingScopeRoot(behaviorName).IsRoot(),
					features,
					behaviorModule))
			continue
		}

		var resolveResult *PolymerBehaviorExtractor_ResolveBehaviorNameResult
		if polymerBehaviorExtractor_isGoogModuleGetCall(behaviorName) {
			resolveResult = e.resolveGoogModuleGet(behaviorName.GetSecondChild().GetString())
		} else {
			resolveResult = e.resolveBehaviorName(polymerBehaviorExtractor_getQualifiedNameThroughCast(behaviorName), moduleMetadata)
		}
		if resolveResult == polymerBehaviorExtractor_FAILED_RESOLVE_RESULT {
			e.compiler.Report(JSError_MakeWithNode(behaviorName, PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR))
			continue
		}
		behaviorValue := resolveResult.node

		if behaviorValue.IsArrayLit() {
			// Individual behaviors can also be arrays of behaviors. Parse them recursively.
			behaviors = append(behaviors, e.ExtractBehaviors(behaviorValue, resolveResult.moduleMetadata)...)
		} else if behaviorValue.IsObjectLit() {
			PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(behaviorValue, e.compiler)
			PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(behaviorValue, e.compiler)
			if !NodeUtil_GetFirstPropMatchingKey(behaviorValue, "is").IsNil() {
				e.compiler.Report(JSError_MakeWithNode(behaviorValue, PolymerPassErrors_POLYMER_INVALID_BEHAVIOR))
			}
			behaviorModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorValue)
			features, _ := NodeUtil_GetEnclosingScript(behaviorValue).GetProp(rhino.Node_FEATURE_SET).(*parser.FeatureSet)
			behaviors = append(behaviors,
				NewPolymerBehaviorExtractor_BehaviorDefinition(
					PolymerPassStaticUtils_ExtractProperties(
						behaviorValue,
						PolymerClassDefinition_DefinitionType_ObjectLiteral,
						e.compiler,
						/* constructor= */ rhino.Node{}),
					polymerBehaviorExtractor_getBehaviorFunctionsToCopy(behaviorValue),
					polymerBehaviorExtractor_getNonPropertyMembersToCopy(behaviorValue),
					resolveResult.isGlobalDeclaration,
					features,
					behaviorModule))
		} else {
			e.compiler.Report(JSError_MakeWithNode(behaviorName, PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR))
		}
	}
	return behaviors
}

// Java: private static class PolymerBehaviorExtractor.ResolveBehaviorNameResult
type PolymerBehaviorExtractor_ResolveBehaviorNameResult struct {
	node                rhino.Node
	isGlobalDeclaration bool
	moduleMetadata      *ModuleMetadataMap_ModuleMetadata
}

// Java: PolymerBehaviorExtractor.ResolveBehaviorNameResult(Node, boolean, ModuleMetadata)
func NewPolymerBehaviorExtractor_ResolveBehaviorNameResult(
	node rhino.Node, isGlobalDeclaration bool, moduleMetadata *ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	return &PolymerBehaviorExtractor_ResolveBehaviorNameResult{
		node:                node,
		isGlobalDeclaration: isGlobalDeclaration,
		moduleMetadata:      moduleMetadata,
	}
}

var polymerBehaviorExtractor_FAILED_RESOLVE_RESULT = NewPolymerBehaviorExtractor_ResolveBehaviorNameResult(rhino.Node{}, false, nil)

// resolveBehaviorName resolves an identifier, which is presumed to refer to a
// Polymer Behavior declaration, using the global namespace. Recurses to
// resolve assignment chains of any length.
//
// This method memoizes resolveBehaviorNameInternal(String, ModuleMetadata)
//
// name is the name of the identifier, which may be qualified; moduleMetadata
// the module (ES module or goog.module) this name is resolved in, or null if
// not in a module. Returns the behavior declaration node, or
// FAILED_RESOLVE_RESULT if it couldn't be resolved.
//
// Java: PolymerBehaviorExtractor.resolveBehaviorName(String, ModuleMetadata)
func (e *PolymerBehaviorExtractor) resolveBehaviorName(
	name *string, moduleMetadata *ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	if name == nil {
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}
	var memoized *PolymerBehaviorExtractor_ResolveBehaviorNameResult
	if moduleMetadata != nil {
		memoized, _ = e.resolveMemoized.Get(*name, moduleMetadata)
	} else {
		memoized = e.globalResolveMemoized.Get(*name)
	}
	if memoized == nil {
		memoized = jutil.CheckNotNull(e.resolveBehaviorNameInternal(*name, moduleMetadata))
		if moduleMetadata != nil {
			e.resolveMemoized.Put(*name, moduleMetadata, memoized)
		} else {
			e.globalResolveMemoized.Put(*name, memoized)
		}
	}
	return memoized
}

// resolveBehaviorNameInternal implements behavior resolution. Call
// resolveBehaviorName(String, ModuleMetadata) instead.
//
// Java: PolymerBehaviorExtractor.resolveBehaviorNameInternal(String, ModuleMetadata)
func (e *PolymerBehaviorExtractor) resolveBehaviorNameInternal(
	name string, moduleMetadata *ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	// Check if this name is a module import/require.
	result := e.getNameIfModuleImport(name, moduleMetadata)
	if result != nil {
		return result
	}

	// Check if this name is possibly from a legacy goog.module
	legacyResolve := e.resolveReferenceToLegacyGoogModule(name)
	if legacyResolve != nil {
		return legacyResolve
	}

	// If not, look it up within the current module.
	var moduleLevelName *GlobalNamespace_Name
	if moduleMetadata != nil {
		moduleLevelName = e.globalNames.GetNameFromModule(moduleMetadata, name)
	}
	globalName := moduleLevelName
	if moduleLevelName == nil {
		globalName, _ = e.globalNames.GetSlot(name).(*GlobalNamespace_Name)
	}
	if globalName == nil {
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}

	// Whether the declaration of this node is in the top-level global scope, as opposed to a module
	// or an IIFE.
	isGlobalDeclaration := moduleLevelName == nil

	// Use any set as a backup declaration, even if it's local.
	declarationRef, _ := globalName.GetDeclaration().(*GlobalNamespace_Ref)
	if declarationRef == nil {
		for _, ref := range globalName.GetRefs() {
			if ref.IsSet() {
				isGlobalDeclaration = false
				declarationRef = ref
				break
			}
		}
	}
	if declarationRef == nil {
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}

	declarationNode := declarationRef.GetNode()
	if declarationNode.IsNil() {
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}
	rValue := NodeUtil_GetRValueOfLValue(declarationNode)
	if rValue.IsNil() {
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}

	if rValue.IsQualifiedName() {
		// Another identifier; recurse.
		declarationScope := declarationRef.scope.GetClosestHoistScope()
		m := ModuleImportResolver_GetModuleFromScopeRoot(
			e.compiler.GetModuleMap(), e.compiler, declarationScope.GetRootNode())
		var mMetadata *ModuleMetadataMap_ModuleMetadata
		if m != nil {
			mMetadata = m.Metadata()
		}
		return e.resolveBehaviorName(polymerBehaviorExtractor_getQualifiedNameThroughCast(rValue), mMetadata)
	}

	behaviorInfo := NodeUtil_GetBestJSDocInfo(declarationNode)
	if behaviorInfo == nil || !behaviorInfo.IsPolymerBehavior() {
		e.compiler.Report(
			JSError_MakeWithNode(declarationNode, PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR))
	}

	return NewPolymerBehaviorExtractor_ResolveBehaviorNameResult(rValue, isGlobalDeclaration, moduleMetadata)
}

// resolveReferenceToLegacyGoogModule handles resolving behaviors if they are
// references to legacy modules
//
// Returns null if the name is not from a legacy module, and resolution
// should continue normally.
//
// Java: PolymerBehaviorExtractor.resolveReferenceToLegacyGoogModule(String)
func (e *PolymerBehaviorExtractor) resolveReferenceToLegacyGoogModule(name string) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	// Qualified names are made of identifiers and '.', so UTF-16 and byte
	// indices of '.' agree for the substrings taken here.
	dot := len(name)
	for dot >= 0 {
		subNamespace := name[:dot]
		metadata := e.moduleMetadataMap.GetModulesByGoogNamespace().Get(subNamespace)

		if metadata == nil || !metadata.IsLegacyGoogModule() {
			// name.lastIndexOf('.', dot - 1)
			if dot-1 < 0 {
				dot = -1
			} else {
				dot = strings.LastIndexByte(name[:dot], '.')
			}
			continue
		}

		rest := ""
		if dot != len(name) {
			rest = name[dot:]
		}
		exportsName := polymerBehaviorExtractor_GOOG_MODULE_EXPORTS + rest
		result := e.resolveBehaviorName(&exportsName, metadata)
		// TODO(lharker): Remove this check and just fail to resolve once we have moved module
		// rewriting unconditionally after the PolymerPass.
		if result == polymerBehaviorExtractor_FAILED_RESOLVE_RESULT {
			return nil
		}
		return result
	}
	return nil
}

var polymerBehaviorExtractor_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// Java: PolymerBehaviorExtractor.isGoogModuleGetCall(Node)
func polymerBehaviorExtractor_isGoogModuleGetCall(callNode rhino.Node) bool {
	if !callNode.IsCall() {
		return false
	}
	return polymerBehaviorExtractor_GOOG_MODULE_GET.Matches(callNode.GetFirstChild()) &&
		callNode.HasTwoChildren() &&
		callNode.GetSecondChild().IsStringLit()
}

// Java: PolymerBehaviorExtractor.resolveGoogModuleGet(String)
func (e *PolymerBehaviorExtractor) resolveGoogModuleGet(moduleNamespace string) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	closureModule := e.moduleMetadataMap.GetModulesByGoogNamespace().Get(moduleNamespace)
	if closureModule == nil {
		// Invalid goog.module.get() call.
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	} else if closureModule.IsGoogProvide() {
		return e.resolveBehaviorName(&moduleNamespace, nil)
	}
	jutil.CheckStateMsg1(closureModule.IsGoogModule(), "%s", closureModule)
	exportsName := polymerBehaviorExtractor_GOOG_MODULE_EXPORTS
	return e.resolveBehaviorName(&exportsName, closureModule)
}

// getNameIfModuleImport handles resolving behaviors whose root is imported
// from another module or a provide.
//
// Returns null if the given name is not imported or FAILED_RESOLVE_RESULT if
// it is imported but is not annotated @polymerBehavior.
//
// Java: PolymerBehaviorExtractor.getNameIfModuleImport(String, ModuleMetadata)
func (e *PolymerBehaviorExtractor) getNameIfModuleImport(
	name string, metadata *ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	if metadata == nil || (!metadata.IsEs6Module() && !metadata.IsGoogModule()) {
		return nil
	}
	var module *Module
	if metadata.IsGoogModule() {
		module = e.moduleMap.GetClosureModule(polymerBehaviorExtractor_firstGoogNamespace(metadata))
	} else {
		module = e.moduleMap.GetModuleByPath(metadata.Path())
	}

	jutil.CheckNotNullMsg1(module, "%s", metadata)

	// '.' is ASCII and qualified names are made of identifiers: byte and
	// UTF-16 indices agree for the substrings taken here.
	dot := strings.IndexByte(name, '.')
	root := name
	if dot != -1 {
		root = name[:dot]
	}
	b := module.BoundNames().Get(root)

	if b == nil || !b.IsSomeImport() {
		return nil
	}
	rest := ""
	if dot != -1 {
		rest = name[dot:]
	}

	if b.IsModuleNamespace() {
		// `import * as x from '';` or `const ns = goog.require('...`
		return e.resolveModuleNamespaceBinding(b, rest)
	}
	importMetadata := b.OriginatingExport().ModuleMetadata()
	var originatingName string // The name in the module being imported
	if importMetadata.IsEs6Module() {
		// import {exportName} from './mod';
		originatingName = jutil.ToJavaString(b.OriginatingExport().LocalName()) + rest
	} else if importMetadata.IsGoogModule() {
		// `const {exportName: localName} = goog.require('some.module');`
		originatingName = polymerBehaviorExtractor_GOOG_MODULE_EXPORTS + "." + jutil.ToJavaString(b.OriginatingExport().ExportName()) + rest
	} else {
		// `const {exportName: localName} = goog.require('some.provide');`
		jutil.CheckStateMsg1(importMetadata.IsGoogProvide(), "%s", importMetadata)
		originatingName = jutil.ToJavaString(b.ClosureNamespace()) + "." + jutil.ToJavaString(b.OriginatingExport().ExportName()) + rest
	}
	return e.resolveBehaviorName(&originatingName, importMetadata)
}

// polymerBehaviorExtractor_firstGoogNamespace is
// metadata.googNamespaces().asList().get(0).
//
// Go-only: language-adapter — ImmutableMultiset.asList().get(0).
func polymerBehaviorExtractor_firstGoogNamespace(metadata *ModuleMetadataMap_ModuleMetadata) string {
	namespaces := slices.Collect(metadata.GoogNamespaces().All())
	return namespaces[jutil.CheckElementIndex(0, len(namespaces))]
}

// resolveModuleNamespaceBinding resolves a name that imports the 'namespace'
// of a module or provide.
//
// Java: PolymerBehaviorExtractor.resolveModuleNamespaceBinding(Binding, String)
func (e *PolymerBehaviorExtractor) resolveModuleNamespaceBinding(b *Binding, rest string) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	if b.Metadata().IsGoogModule() {
		exportsName := polymerBehaviorExtractor_GOOG_MODULE_EXPORTS + rest
		return e.resolveBehaviorName(&exportsName, b.Metadata())
	} else if b.Metadata().IsGoogProvide() {
		providedName := jutil.ToJavaString(b.ClosureNamespace()) + rest
		return e.resolveBehaviorName(&providedName, b.Metadata())
	}

	// ES module import *.
	jutil.CheckState(b.Metadata().IsEs6Module())
	if rest == "" {
		// The namespace imported by `import *` is never a @polymerBehavior.
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}
	rest = rest[1:] // Remove leading '.'.
	dot := strings.IndexByte(rest, '.')

	// Given:
	//   `const internalName = 0; export {internalName as exportName};`
	//   `import * as mod from './x'; use(mod.exportName.Behavior);`
	// 1. get the internal name `internalName` from` exportName`.
	// 2. then proceed to resolve `internalName.Behavior` in './x'.
	exportedName := rest
	if dot != -1 {
		exportedName = rest[:dot]
	}
	originalModule := e.moduleMap.GetModuleByPath(b.Metadata().Path())
	exportBinding := originalModule.Namespace().Get(exportedName)
	if exportBinding == nil || !exportBinding.IsCreatedByEsExport() {
		// This is an invalid import, and will cause an error elsewhere.
		return polymerBehaviorExtractor_FAILED_RESOLVE_RESULT
	}
	tail := ""
	if dot != -1 {
		tail = rest[dot:]
	}
	resolvedName := jutil.ToJavaString(exportBinding.OriginatingExport().LocalName()) + tail
	return e.resolveBehaviorName(&resolvedName, b.Metadata())
}

// polymerBehaviorExtractor_getBehaviorFunctionsToCopy returns a list of
// functions from a behavior which should be copied to the element prototype.
//
// Java: PolymerBehaviorExtractor.getBehaviorFunctionsToCopy(Node)
func polymerBehaviorExtractor_getBehaviorFunctionsToCopy(behaviorObjLit rhino.Node) []*PolymerPass_MemberDefinition {
	jutil.CheckState(behaviorObjLit.IsObjectLit())
	var functionsToCopy []*PolymerPass_MemberDefinition
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)

	for keyNode := behaviorObjLit.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		isFunctionDefinition := (keyNode.IsStringKey() && keyNode.GetFirstChild().IsFunction()) ||
			keyNode.IsMemberFunctionDef()
		if isFunctionDefinition && !polymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY.Contains(keyNode.GetString()) {
			functionsToCopy = append(functionsToCopy,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(keyNode),
					keyNode,
					keyNode.GetFirstChild(),
					enclosingModule))
		}
	}

	return functionsToCopy
}

// polymerBehaviorExtractor_getQualifiedNameThroughCast is similar to
// Node#getQualifiedName but also handles CAST nodes. For example, given a
// GETPROP representing "(/** @type {?} */ (x)).y.z" returns "x.y.z". Returns
// null if node is not a NAME, GETPROP, or CAST. See b/64389806 for
// Polymer-specific context.
//
// Java: PolymerBehaviorExtractor.getQualifiedNameThroughCast(Node)
func polymerBehaviorExtractor_getQualifiedNameThroughCast(node rhino.Node) *string {
	if node.IsName() {
		name := node.GetString()
		if name == "" {
			return nil
		}
		return &name
	} else if node.IsGetProp() {
		left := polymerBehaviorExtractor_getQualifiedNameThroughCast(node.GetFirstChild())
		if left == nil {
			return nil
		}
		right := node.GetString()
		qname := *left + "." + right
		return &qname
	} else if node.IsCast() {
		return polymerBehaviorExtractor_getQualifiedNameThroughCast(node.GetFirstChild())
	}
	return nil
}

// polymerBehaviorExtractor_getNonPropertyMembersToCopy returns a list of
// MemberDefinitions in a behavior which are not in the properties block, but
// should still be copied to the element prototype.
//
// Java: PolymerBehaviorExtractor.getNonPropertyMembersToCopy(Node)
func polymerBehaviorExtractor_getNonPropertyMembersToCopy(behaviorObjLit rhino.Node) []*PolymerPass_MemberDefinition {
	jutil.CheckState(behaviorObjLit.IsObjectLit())
	var membersToCopy []*PolymerPass_MemberDefinition
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)

	for keyNode := behaviorObjLit.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		isNonFunctionMember := keyNode.IsGetterDef() || (keyNode.IsStringKey() && !keyNode.GetFirstChild().IsFunction())
		if isNonFunctionMember && !polymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY.Contains(keyNode.GetString()) {
			membersToCopy = append(membersToCopy,
				NewPolymerPass_MemberDefinition(
					NodeUtil_GetBestJSDocInfo(keyNode),
					keyNode,
					keyNode.GetFirstChild(),
					enclosingModule))
		}
	}

	return membersToCopy
}

// PolymerBehaviorExtractor_BehaviorDefinition is a parsed definition of a
// Polymer Behavior. Includes members which should be copied to elements which
// use the behavior.
//
// Java: static final class PolymerBehaviorExtractor.BehaviorDefinition
type PolymerBehaviorExtractor_BehaviorDefinition struct {
	// Properties declared in the behavior 'properties' block.
	props []*PolymerPass_MemberDefinition

	// Functions intended to be copied to elements which use this Behavior.
	functionsToCopy []*PolymerPass_MemberDefinition

	// Other members intended to be copied to elements which use this Behavior.
	nonPropertyMembersToCopy []*PolymerPass_MemberDefinition

	// Whether this Behavior is declared in the global scope.
	isGlobalDeclaration bool

	// Language features to carry over to the extraction destination.
	features *parser.FeatureSet

	// Containing MODULE_BODY if this behavior is defined inside a module, otherwise null
	behaviorModule rhino.Node

	lazyModuleLocalNames *jutil.LinkedHashSet[string]
}

// Java: PolymerBehaviorExtractor.BehaviorDefinition(List<MemberDefinition>, List<MemberDefinition>, List<MemberDefinition>, boolean, FeatureSet, Node)
func NewPolymerBehaviorExtractor_BehaviorDefinition(
	props []*PolymerPass_MemberDefinition,
	functionsToCopy []*PolymerPass_MemberDefinition,
	nonPropertyMembersToCopy []*PolymerPass_MemberDefinition,
	isGlobalDeclaration bool,
	features *parser.FeatureSet,
	behaviorModule rhino.Node) *PolymerBehaviorExtractor_BehaviorDefinition {
	return &PolymerBehaviorExtractor_BehaviorDefinition{
		props:                    props,
		functionsToCopy:          functionsToCopy,
		nonPropertyMembersToCopy: nonPropertyMembersToCopy,
		isGlobalDeclaration:      isGlobalDeclaration,
		features:                 features,
		behaviorModule:           behaviorModule,
	}
}

// Java: PolymerBehaviorExtractor.BehaviorDefinition.getModuleLocalNames(AbstractCompiler)
func (d *PolymerBehaviorExtractor_BehaviorDefinition) GetModuleLocalNames(compiler AbstractCompiler) *jutil.LinkedHashSet[string] {
	if d.lazyModuleLocalNames == nil {
		d.lazyModuleLocalNames = d.accumulateModuleLocalVars(compiler)
	}
	return d.lazyModuleLocalNames
}

// Java: PolymerBehaviorExtractor.BehaviorDefinition.accumulateModuleLocalVars(AbstractCompiler)
func (d *PolymerBehaviorExtractor_BehaviorDefinition) accumulateModuleLocalVars(compiler AbstractCompiler) *jutil.LinkedHashSet[string] {
	scopeCreator := NewSyntacticScopeCreator(compiler)
	globalScope := Scope_CreateGlobalScope(d.behaviorModule.GetParent())
	return NodeUtil_GetAllVarNamesDeclaredInModule(
		d.behaviorModule, compiler, scopeCreator, globalScope)
}
