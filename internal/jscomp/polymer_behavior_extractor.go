// Port of com/google/javascript/jscomp/PolymerBehaviorExtractor.java @ closure-compiler v20261104

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func isBehaviorNameNotToCopy(name string) bool {
	switch name {
	case "created",
		"attached",
		"detached",
		"attributeChanged",
		"configure",
		"ready",
		"properties",
		"listeners",
		"observers",
		"hostAttributes":
		return true
	default:
		return false
	}
}

const polymerGoogModuleExports = "exports"

type polymerResolveKey struct {
	name           string
	moduleMetadata *ModuleMetadataMap_ModuleMetadata
}

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

	resolveMemoized map[polymerResolveKey]*polymerResolveBehaviorResult
}

// Java: PolymerBehaviorExtractor(AbstractCompiler, GlobalNamespace, ModuleMetadataMap, ModuleMap)
func NewPolymerBehaviorExtractor(
	compiler AbstractCompiler,
	globalNames *GlobalNamespace,
	moduleMetadataMap *ModuleMetadataMap,
	moduleMap *ModuleMap) *PolymerBehaviorExtractor {
	return &PolymerBehaviorExtractor{
		compiler:          compiler,
		globalNames:       globalNames,
		moduleMetadataMap: moduleMetadataMap,
		moduleMap:         moduleMap,
		resolveMemoized:   make(map[polymerResolveKey]*polymerResolveBehaviorResult),
	}
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

	arr, ok := ast.AsArrayLit(behaviorArray)
	if !ok {
		e.compiler.Report(
			JSError_MakeWithNode(behaviorArray, PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY))
		return nil
	}

	var behaviors []*PolymerBehaviorExtractor_BehaviorDefinition
	for elem := range arr.Elements() {
		behaviorName := elem.Node()
		if behaviorName.IsObjectLit() {
			behaviors = append(behaviors,
				e.extractBehaviorObjectLit(
					behaviorName,
					/* isGlobalDeclaration= */ NodeUtil_GetEnclosingScopeRoot(behaviorName).IsRoot(),
					/* protectFunctionKeys= */ true))
			continue
		}

		var resolveResult *polymerResolveBehaviorResult
		if ns, ok := asPolymerGoogModuleGetNamespace(behaviorName); ok {
			resolveResult = e.resolveGoogModuleGet(ns)
		} else if qname, ok := qualifiedNameThroughCast(behaviorName); ok {
			resolveResult = e.resolveBehaviorName(qname, moduleMetadata)
		} else {
			resolveResult = failedBehaviorResolveResult
		}
		if resolveResult == failedBehaviorResolveResult {
			e.compiler.Report(JSError_MakeWithNode(behaviorName, PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR))
			continue
		}
		behaviorValue := resolveResult.node

		if behaviorValue.IsArrayLit() {
			// Individual behaviors can also be arrays of behaviors. Parse them recursively.
			behaviors = append(behaviors, e.ExtractBehaviors(behaviorValue, resolveResult.moduleMetadata)...)
		} else if behaviorValue.IsObjectLit() {
			behaviors = append(behaviors,
				e.extractBehaviorObjectLit(behaviorValue, resolveResult.isGlobalDeclaration /* protectFunctionKeys= */, false))
		} else {
			e.compiler.Report(JSError_MakeWithNode(behaviorName, PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR))
		}
	}
	return behaviors
}

func (e *PolymerBehaviorExtractor) extractBehaviorObjectLit(
	behaviorObjLit rhino.Node, isGlobalDeclaration bool, protectFunctionKeys bool) *PolymerBehaviorExtractor_BehaviorDefinition {
	PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(behaviorObjLit, e.compiler)
	PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(behaviorObjLit, e.compiler)
	if protectFunctionKeys {
		PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(behaviorObjLit)
	}
	if !NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "is").IsNil() {
		e.compiler.Report(JSError_MakeWithNode(behaviorObjLit, PolymerPassErrors_POLYMER_INVALID_BEHAVIOR))
	}
	behaviorModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)
	features, _ := NodeUtil_GetEnclosingScript(behaviorObjLit).GetProp(rhino.Node_FEATURE_SET).(parser.FeatureSet)
	return NewPolymerBehaviorExtractor_BehaviorDefinition(
		PolymerPassStaticUtils_ExtractProperties(
			behaviorObjLit,
			PolymerClassDefinition_DefinitionType_ObjectLiteral,
			e.compiler,
			/* constructor= */ rhino.NilNode),
		behaviorFunctionsToCopy(behaviorObjLit),
		nonPropertyBehaviorMembersToCopy(behaviorObjLit),
		isGlobalDeclaration,
		features,
		behaviorModule)
}

// Java: private static class PolymerBehaviorExtractor.ResolveBehaviorNameResult
type polymerResolveBehaviorResult struct {
	node                rhino.Node
	isGlobalDeclaration bool
	moduleMetadata      *ModuleMetadataMap_ModuleMetadata
}

var failedBehaviorResolveResult = &polymerResolveBehaviorResult{node: rhino.NilNode}

// resolveBehaviorName resolves an identifier, which is presumed to refer to a
// Polymer Behavior declaration, using the global namespace. Recurses to
// resolve assignment chains of any length.
//
// This method memoizes resolveBehaviorNameInternal(String, ModuleMetadata)
//
// name is the name of the identifier, which may be qualified; moduleMetadata
// the module (ES module or goog.module) this name is resolved in, or null if
// not in a module. Returns the behavior declaration node, or
// failedBehaviorResolveResult if it couldn't be resolved.
//
// Java: PolymerBehaviorExtractor.resolveBehaviorName(String, ModuleMetadata)
func (e *PolymerBehaviorExtractor) resolveBehaviorName(
	name string, moduleMetadata *ModuleMetadataMap_ModuleMetadata) *polymerResolveBehaviorResult {
	key := polymerResolveKey{name: name, moduleMetadata: moduleMetadata}
	memoized := e.resolveMemoized[key]
	if memoized == nil {
		memoized = jutil.CheckNotNull(e.resolveBehaviorNameInternal(name, moduleMetadata))
		e.resolveMemoized[key] = memoized
	}
	return memoized
}

// resolveBehaviorNameInternal implements behavior resolution. Call
// resolveBehaviorName(String, ModuleMetadata) instead.
//
// Java: PolymerBehaviorExtractor.resolveBehaviorNameInternal(String, ModuleMetadata)
func (e *PolymerBehaviorExtractor) resolveBehaviorNameInternal(
	name string, moduleMetadata *ModuleMetadataMap_ModuleMetadata) *polymerResolveBehaviorResult {
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
		return failedBehaviorResolveResult
	}

	// Whether the declaration of this node is in the top-level global scope, as opposed to a module
	// or an IIFE.
	isGlobalDeclaration := moduleLevelName == nil

	// Use any set as a backup declaration, even if it's local.
	declarationRef, _ := globalName.GetDeclaration().(*GlobalNamespace_Ref)
	if declarationRef == nil {
		for ref := range globalName.Refs() {
			if ref.IsSet() {
				isGlobalDeclaration = false
				declarationRef = ref
				break
			}
		}
	}
	if declarationRef == nil {
		return failedBehaviorResolveResult
	}

	declarationNode := declarationRef.GetNode()
	if declarationNode.IsNil() {
		return failedBehaviorResolveResult
	}
	rValue := NodeUtil_GetRValueOfLValue(declarationNode)
	if rValue.IsNil() {
		return failedBehaviorResolveResult
	}

	if rValue.IsQualifiedName() {
		// Another identifier; recurse.
		qname, ok := qualifiedNameThroughCast(rValue)
		if !ok {
			return failedBehaviorResolveResult
		}
		declarationScope := declarationRef.scope.GetClosestHoistScope()
		m := ModuleImportResolver_GetModuleFromScopeRoot(
			e.compiler.GetModuleMap(), e.compiler, declarationScope.GetRootNode())
		var mMetadata *ModuleMetadataMap_ModuleMetadata
		if m != nil {
			mMetadata = m.Metadata()
		}
		return e.resolveBehaviorName(qname, mMetadata)
	}

	behaviorInfo := NodeUtil_GetBestJSDocInfo(declarationNode)
	if behaviorInfo == nil || !behaviorInfo.IsPolymerBehavior() {
		e.compiler.Report(
			JSError_MakeWithNode(declarationNode, PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR))
	}

	return &polymerResolveBehaviorResult{
		node:                rValue,
		isGlobalDeclaration: isGlobalDeclaration,
		moduleMetadata:      moduleMetadata,
	}
}

// resolveReferenceToLegacyGoogModule handles resolving behaviors if they are
// references to legacy modules
//
// Returns null if the name is not from a legacy module, and resolution
// should continue normally.
//
// Java: PolymerBehaviorExtractor.resolveReferenceToLegacyGoogModule(String)
func (e *PolymerBehaviorExtractor) resolveReferenceToLegacyGoogModule(name string) *polymerResolveBehaviorResult {
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
		result := e.resolveBehaviorName(polymerGoogModuleExports+rest, metadata)
		// TODO(lharker): Remove this check and just fail to resolve once we have moved module
		// rewriting unconditionally after the PolymerPass.
		if result == failedBehaviorResolveResult {
			return nil
		}
		return result
	}
	return nil
}

var polymerGoogModuleGetQName = rhino.QualifiedName_Of("goog.module.get")

// Java: PolymerBehaviorExtractor.isGoogModuleGetCall(Node)
func asPolymerGoogModuleGetNamespace(callNode rhino.Node) (string, bool) {
	call, ok := ast.AsCall(callNode)
	if !ok || !polymerGoogModuleGetQName.Matches(call.Callee().Node()) {
		return "", false
	}
	arg, ok := call.OnlyArg()
	if !ok {
		return "", false
	}
	str, ok := arg.AsStringLit()
	if !ok {
		return "", false
	}
	return str.Value(), true
}

// Java: PolymerBehaviorExtractor.resolveGoogModuleGet(String)
func (e *PolymerBehaviorExtractor) resolveGoogModuleGet(moduleNamespace string) *polymerResolveBehaviorResult {
	closureModule := e.moduleMetadataMap.GetModulesByGoogNamespace().Get(moduleNamespace)
	if closureModule == nil {
		// Invalid goog.module.get() call.
		return failedBehaviorResolveResult
	} else if closureModule.IsGoogProvide() {
		return e.resolveBehaviorName(moduleNamespace, nil)
	}
	jutil.CheckStateMsg1(closureModule.IsGoogModule(), "%s", closureModule)
	return e.resolveBehaviorName(polymerGoogModuleExports, closureModule)
}

// getNameIfModuleImport handles resolving behaviors whose root is imported
// from another module or a provide.
//
// Returns null if the given name is not imported or failedBehaviorResolveResult
// if it is imported but is not annotated @polymerBehavior.
//
// Java: PolymerBehaviorExtractor.getNameIfModuleImport(String, ModuleMetadata)
func (e *PolymerBehaviorExtractor) getNameIfModuleImport(
	name string, metadata *ModuleMetadataMap_ModuleMetadata) *polymerResolveBehaviorResult {
	if metadata == nil || (!metadata.IsEs6Module() && !metadata.IsGoogModule()) {
		return nil
	}
	var module *Module
	if metadata.IsGoogModule() {
		module = e.moduleMap.GetClosureModule(firstPolymerGoogNamespace(metadata))
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
	b := module.BoundNamesMap().Get(root)

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
	origExport := b.OriginatingExport()
	importMetadata := origExport.ModuleMetadata()
	var originatingName string // The name in the module being imported
	if importMetadata.IsEs6Module() {
		// import {exportName} from './mod';
		originatingName = modules_optString(origExport.LocalName(), origExport.HasLocalName()) + rest
	} else if importMetadata.IsGoogModule() {
		// `const {exportName: localName} = goog.require('some.module');`
		originatingName = polymerGoogModuleExports + "." +
			modules_optString(origExport.ExportName(), origExport.HasExportName()) + rest
	} else {
		// `const {exportName: localName} = goog.require('some.provide');`
		jutil.CheckStateMsg1(importMetadata.IsGoogProvide(), "%s", importMetadata)
		originatingName = modules_optString(b.ClosureNamespace(), b.HasClosureNamespace()) + "." +
			modules_optString(origExport.ExportName(), origExport.HasExportName()) + rest
	}
	return e.resolveBehaviorName(originatingName, importMetadata)
}

// firstPolymerGoogNamespace is metadata.googNamespaces().asList().get(0).
func firstPolymerGoogNamespace(metadata *ModuleMetadataMap_ModuleMetadata) string {
	namespaces := slices.Collect(metadata.GoogNamespaces().All())
	return namespaces[jutil.CheckElementIndex(0, len(namespaces))]
}

// resolveModuleNamespaceBinding resolves a name that imports the 'namespace'
// of a module or provide.
//
// Java: PolymerBehaviorExtractor.resolveModuleNamespaceBinding(Binding, String)
func (e *PolymerBehaviorExtractor) resolveModuleNamespaceBinding(b *Binding, rest string) *polymerResolveBehaviorResult {
	if b.Metadata().IsGoogModule() {
		return e.resolveBehaviorName(polymerGoogModuleExports+rest, b.Metadata())
	} else if b.Metadata().IsGoogProvide() {
		providedName := modules_optString(b.ClosureNamespace(), b.HasClosureNamespace()) + rest
		return e.resolveBehaviorName(providedName, b.Metadata())
	}

	// ES module import *.
	jutil.CheckState(b.Metadata().IsEs6Module())
	if rest == "" {
		// The namespace imported by `import *` is never a @polymerBehavior.
		return failedBehaviorResolveResult
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
	exportBinding := originalModule.NamespaceMap().Get(exportedName)
	if exportBinding == nil || !exportBinding.IsCreatedByEsExport() {
		// This is an invalid import, and will cause an error elsewhere.
		return failedBehaviorResolveResult
	}
	tail := ""
	if dot != -1 {
		tail = rest[dot:]
	}
	origExport := exportBinding.OriginatingExport()
	resolvedName := modules_optString(origExport.LocalName(), origExport.HasLocalName()) + tail
	return e.resolveBehaviorName(resolvedName, b.Metadata())
}

// behaviorFunctionsToCopy returns a list of functions from a behavior which
// should be copied to the element prototype.
//
// Java: PolymerBehaviorExtractor.getBehaviorFunctionsToCopy(Node)
func behaviorFunctionsToCopy(behaviorObjLit rhino.Node) []*PolymerPass_MemberDefinition {
	obj := ast.MustObjectLit(behaviorObjLit)
	var functionsToCopy []*PolymerPass_MemberDefinition
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)

	for member := range obj.Members() {
		if fnDef, ok := member.AsMemberFunctionDef(); ok {
			if !isBehaviorNameNotToCopy(fnDef.Key()) {
				keyNode := fnDef.Node()
				functionsToCopy = append(functionsToCopy,
					NewPolymerPass_MemberDefinition(
						NodeUtil_GetBestJSDocInfo(keyNode), keyNode, fnDef.Function().Node(), enclosingModule))
			}
		} else if sk, ok := member.AsStringKey(); ok && sk.Value().Node().IsFunction() {
			if !isBehaviorNameNotToCopy(sk.Key()) {
				keyNode := sk.Node()
				functionsToCopy = append(functionsToCopy,
					NewPolymerPass_MemberDefinition(
						NodeUtil_GetBestJSDocInfo(keyNode), keyNode, sk.Value().Node(), enclosingModule))
			}
		}
	}

	return functionsToCopy
}

// qualifiedNameThroughCast is similar to Node#getQualifiedName but also handles
// CAST nodes. For example, given a GETPROP representing
// "(/** @type {?} */ (x)).y.z" returns "x.y.z". Returns false if node is not a
// NAME, GETPROP, or CAST. See b/64389806 for Polymer-specific context.
//
// Java: PolymerBehaviorExtractor.getQualifiedNameThroughCast(Node)
func qualifiedNameThroughCast(node rhino.Node) (string, bool) {
	if nameNode, ok := ast.AsName(node); ok {
		name := nameNode.Ident()
		if name == "" {
			return "", false
		}
		return name, true
	} else if gp, ok := ast.AsGetProp(node); ok {
		left, ok := qualifiedNameThroughCast(gp.Object().Node())
		if !ok {
			return "", false
		}
		return left + "." + gp.Property(), true
	} else if cast, ok := ast.AsCast(node); ok {
		return qualifiedNameThroughCast(cast.Expr().Node())
	}
	return "", false
}

// nonPropertyBehaviorMembersToCopy returns a list of MemberDefinitions in a
// behavior which are not in the properties block, but should still be copied to
// the element prototype.
//
// Java: PolymerBehaviorExtractor.getNonPropertyMembersToCopy(Node)
func nonPropertyBehaviorMembersToCopy(behaviorObjLit rhino.Node) []*PolymerPass_MemberDefinition {
	obj := ast.MustObjectLit(behaviorObjLit)
	var membersToCopy []*PolymerPass_MemberDefinition
	enclosingModule := NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)

	for member := range obj.Members() {
		if getter, ok := member.AsGetterDef(); ok {
			if !isBehaviorNameNotToCopy(getter.Key()) {
				keyNode := getter.Node()
				membersToCopy = append(membersToCopy,
					NewPolymerPass_MemberDefinition(
						NodeUtil_GetBestJSDocInfo(keyNode), keyNode, getter.Function().Node(), enclosingModule))
			}
		} else if sk, ok := member.AsStringKey(); ok && !sk.Value().Node().IsFunction() {
			if !isBehaviorNameNotToCopy(sk.Key()) {
				keyNode := sk.Node()
				membersToCopy = append(membersToCopy,
					NewPolymerPass_MemberDefinition(
						NodeUtil_GetBestJSDocInfo(keyNode), keyNode, sk.Value().Node(), enclosingModule))
			}
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
	features parser.FeatureSet

	// Containing MODULE_BODY if this behavior is defined inside a module, otherwise null
	behaviorModule rhino.Node

	lazyModuleLocalNames *OrderedSet[string]
}

// Java: PolymerBehaviorExtractor.BehaviorDefinition(List<MemberDefinition>, List<MemberDefinition>, List<MemberDefinition>, boolean, FeatureSet, Node)
func NewPolymerBehaviorExtractor_BehaviorDefinition(
	props []*PolymerPass_MemberDefinition,
	functionsToCopy []*PolymerPass_MemberDefinition,
	nonPropertyMembersToCopy []*PolymerPass_MemberDefinition,
	isGlobalDeclaration bool,
	features parser.FeatureSet,
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
func (d *PolymerBehaviorExtractor_BehaviorDefinition) GetModuleLocalNames(compiler AbstractCompiler) *OrderedSet[string] {
	if d.lazyModuleLocalNames == nil {
		d.lazyModuleLocalNames = d.accumulateModuleLocalVars(compiler)
	}
	return d.lazyModuleLocalNames
}

// Java: PolymerBehaviorExtractor.BehaviorDefinition.accumulateModuleLocalVars(AbstractCompiler)
func (d *PolymerBehaviorExtractor_BehaviorDefinition) accumulateModuleLocalVars(compiler AbstractCompiler) *OrderedSet[string] {
	scopeCreator := NewSyntacticScopeCreator(compiler)
	globalScope := Scope_CreateGlobalScope(d.behaviorModule.GetParent())
	return NodeUtil_GetAllVarNamesDeclaredInModule(
		d.behaviorModule, compiler, scopeCreator, globalScope)
}
