// Port of com/google/javascript/jscomp/modules/ModuleMapCreator.java @ closure-compiler v20260415

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// ModuleMapCreator creates a ModuleMap.
type ModuleMapCreator struct {
	compiler                            AbstractCompiler
	esModuleProcessor                   *EsModuleProcessor
	closureModuleProcessor              *ClosureModuleProcessor
	nonEsModuleProcessor                *NonEsModuleProcessor
	unresolvedModules                   *jutil.LinkedHashMap[string, UnresolvedModule]
	unresolvedModulesByClosureNamespace *jutil.LinkedHashMap[string, UnresolvedModule]
	moduleMetadataMap                   *ModuleMetadataMap
}

var _ CompilerPass = (*ModuleMapCreator)(nil)

// Java: ModuleMapCreator.DOES_NOT_HAVE_EXPORT
var ModuleMapCreator_DOES_NOT_HAVE_EXPORT = DiagnosticType_Error(
	"JSC_DOES_NOT_HAVE_EXPORT", "Requested module does not have an export \"{0}\".")

// Java: ModuleMapCreator.DOES_NOT_HAVE_EXPORT_WITH_DETAILS
var ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS = DiagnosticType_Error(
	"JSC_DOES_NOT_HAVE_EXPORT_WITH_DETAILS",
	"Requested module does not have an export \"{0}\".{1}")

// ModuleMapCreator_ModuleRequestResolverImpl is
// ModuleMapCreator.ModuleRequestResolverImpl (an inner class: outer is the
// enclosing ModuleMapCreator).
type ModuleMapCreator_ModuleRequestResolverImpl struct {
	outer *ModuleMapCreator
}

var _ ModuleRequestResolver = (*ModuleMapCreator_ModuleRequestResolverImpl)(nil)

// moduleMapCreator_fallbackForMissingNonClosureModule is the anonymous
// UnresolvedModule of
// ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule.
type moduleMapCreator_fallbackForMissingNonClosureModule struct {
	UnresolvedModuleBase

	path     *ModuleLoader_ModulePath
	metadata *ModuleMetadataMap_ModuleMetadata
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath)
func (r *ModuleMapCreator_ModuleRequestResolverImpl) GetFallbackForMissingNonClosureModule(path *ModuleLoader_ModulePath) UnresolvedModule {
	metadata :=
		ModuleMetadataMap_ModuleMetadata_Builder_().
			RootNode(nil).
			Path(path).
			ModuleType(ModuleMetadataMap_ModuleType_ES6_MODULE).
			IsTestOnly(false).
			UsesClosure(false).
			Build()

	m := &moduleMapCreator_fallbackForMissingNonClosureModule{path: path, metadata: metadata}
	m.InitUnresolvedModuleBase(m)
	return m
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath).new UnresolvedModule() {...}.reset()
func (m *moduleMapCreator_fallbackForMissingNonClosureModule) Reset() {}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath).new UnresolvedModule() {...}.resolve(ModuleRequestResolver, String)
func (m *moduleMapCreator_fallbackForMissingNonClosureModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module {
	return Module_Builder_().
		BoundNames(jutil.NewLinkedHashMap[string, *Binding]()).
		Namespace(jutil.NewLinkedHashMap[string, *Binding]()).
		LocalNameToLocalExport(jutil.NewLinkedHashMap[string, *Export]()).
		Path(m.path).
		Metadata(m.metadata).
		Build()
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath).new UnresolvedModule() {...}.metadata()
func (m *moduleMapCreator_fallbackForMissingNonClosureModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath).new UnresolvedModule() {...}.getExportedNames(ModuleRequestResolver)
func (m *moduleMapCreator_fallbackForMissingNonClosureModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet[string]()
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath).new UnresolvedModule() {...}.getExportedNames(ModuleRequestResolver, Set<UnresolvedModule>)
func (m *moduleMapCreator_fallbackForMissingNonClosureModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet[string]()
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingNonClosureModule(ModulePath).new UnresolvedModule() {...}.resolveExport(ModuleRequestResolver, String, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *moduleMapCreator_fallbackForMissingNonClosureModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier *string,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	return ResolveExportResult_Of(
		Binding_From(
			Export_Builder_().
				LocalName(&exportName).
				ModuleMetadata(m.metadata).
				ModulePath(m.path).
				ClosureNamespace(nil).
				Build(),
			/* sourceNode= */ nil))
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.getFallbackForMissingClosureModule(String)
func (r *ModuleMapCreator_ModuleRequestResolverImpl) GetFallbackForMissingClosureModule(namespace string) UnresolvedModule {
	return r.outer.nonEsModuleProcessor.Process(
		ModuleMetadataMap_ModuleMetadata_Builder_().
			AddGoogNamespace(namespace).
			IsTestOnly(false).
			ModuleType(ModuleMetadataMap_ModuleType_GOOG_PROVIDE).
			Path(nil).
			RootNode(nil).
			UsesClosure(true).
			Build(),
		/* path= */ nil,
		/* script= */ nil)
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.resolve(Import)
func (r *ModuleMapCreator_ModuleRequestResolverImpl) Resolve(i *Import) UnresolvedModule {
	if i.ModulePath() == nil {
		return r.ResolveForClosure(i.ModuleRequest())
	}
	return r.ResolveWithModulePath(i.ModuleRequest(), i.ModulePath(), i.ImportNode())
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.resolve(Export)
func (r *ModuleMapCreator_ModuleRequestResolverImpl) ResolveWithExport(e *Export) UnresolvedModule {
	// Go-only: language-adapter — Export.moduleRequest() is a nullable
	// String; the callers only resolve exports that have one.
	return r.ResolveWithModulePath(*e.ModuleRequest(), e.ModulePath(), e.ExportNode())
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.resolveForClosure(String)
func (r *ModuleMapCreator_ModuleRequestResolverImpl) ResolveForClosure(namespace string) UnresolvedModule {
	return r.outer.unresolvedModulesByClosureNamespace.ComputeIfAbsent(
		namespace, r.GetFallbackForMissingClosureModule)
}

// Java: ModuleMapCreator.ModuleRequestResolverImpl.resolve(String, ModulePath, Node)
func (r *ModuleMapCreator_ModuleRequestResolverImpl) ResolveWithModulePath(
	moduleRequest string, modulePath *ModuleLoader_ModulePath, forLineInfo *rhino.Node) UnresolvedModule {

	if GoogEsImports_IsGoogImportSpecifier(moduleRequest) {
		namespace := GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleRequest)
		return r.ResolveForClosure(namespace)
	}

	// Go-only: language-adapter — Node.getSourceFileName() is a nullable
	// String; ModulePath.resolveJsModule takes a Go string.
	sourceFileName := ""
	if s := forLineInfo.GetSourceFileName(); s != nil {
		sourceFileName = *s
	}
	requestedPath :=
		modulePath.ResolveJsModule(
			moduleRequest,
			sourceFileName,
			forLineInfo.GetLineno(),
			forLineInfo.GetCharno())

	if requestedPath == nil {
		requestedPath = modulePath.ResolveModuleAsPath(moduleRequest)

		if !r.outer.unresolvedModules.ContainsKey(requestedPath.ToModuleName()) {
			module := r.GetFallbackForMissingNonClosureModule(requestedPath)
			r.outer.unresolvedModules.Put(requestedPath.ToModuleName(), module)
			return module
		}
	}

	return r.outer.unresolvedModules.Get(requestedPath.ToModuleName())
}

// ModuleMapCreator_ModuleProcessor is a basic interface that can scan and
// return information about a module.
type ModuleMapCreator_ModuleProcessor interface {
	// Java: ModuleMapCreator.ModuleProcessor.process(ModuleMetadata, ModulePath, Node)
	Process(metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script *rhino.Node) UnresolvedModule
}

// Java: ModuleMapCreator(AbstractCompiler, ModuleMetadataMap)
func NewModuleMapCreator(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *ModuleMapCreator {
	c := &ModuleMapCreator{}
	c.compiler = compiler
	c.moduleMetadataMap = moduleMetadataMap
	c.esModuleProcessor = NewEsModuleProcessor(compiler)
	c.closureModuleProcessor = NewClosureModuleProcessor(compiler)
	c.nonEsModuleProcessor = NewNonEsModuleProcessor()
	c.unresolvedModules = jutil.NewLinkedHashMap[string, UnresolvedModule]()
	c.unresolvedModulesByClosureNamespace = jutil.NewLinkedHashMap[string, UnresolvedModule]()
	return c
}

// Java: ModuleMapCreator.create()
func (c *ModuleMapCreator) Create() *ModuleMap {
	c.unresolvedModules.Clear()
	c.unresolvedModulesByClosureNamespace.Clear()

	// There are modules that aren't associated with scripts - nested
	// goog.modules in goog.loadModule calls.
	for moduleMetadata := range c.moduleMetadataMap.GetAllModuleMetadata().All() {
		c.ProcessModuleMetadata(moduleMetadata)
	}

	return c.Resolve()
}

// Java: ModuleMapCreator.resolve()
func (c *ModuleMapCreator) Resolve() *ModuleMap {
	var requestResolver ModuleRequestResolver = &ModuleMapCreator_ModuleRequestResolverImpl{outer: c}
	resolvedModules := jutil.NewLinkedHashMap[string, *Module]()
	resolvedClosureModules := jutil.NewLinkedHashMap[string, *Module]()

	// We need to resolve in a loop as any missing reference will add a fake to
	// the unresolvedModules map (see getFallback* methods above). This would
	// cause a concurrent modification exception if we just iterated over
	// unresolvedModules. So the first loop through should resolve any "known"
	// modules, and the second any "unrecognized" modules.
	for {
		toResolve := moduleMapCreator_difference(c.unresolvedModules, resolvedModules)

		for _, key := range toResolve {
			resolved := c.unresolvedModules.Get(key).Resolve(requestResolver)
			resolvedModules.Put(key, resolved)

			for namespace := range resolved.Metadata().GoogNamespaces().All() {
				resolvedClosureModules.Put(namespace, resolved)
			}
		}

		if moduleMapCreator_containsAllKeys(resolvedModules, c.unresolvedModules) {
			break
		}
	}

	for {
		toResolve := moduleMapCreator_difference(c.unresolvedModulesByClosureNamespace, resolvedClosureModules)

		for _, namespace := range toResolve {
			resolvedClosureModules.Put(
				namespace, c.unresolvedModulesByClosureNamespace.Get(namespace).Resolve(requestResolver))
		}

		if moduleMapCreator_containsAllKeys(resolvedClosureModules, c.unresolvedModulesByClosureNamespace) {
			break
		}
	}

	c.unresolvedModules.Clear()
	c.unresolvedModulesByClosureNamespace.Clear()

	return NewModuleMap(
		jutil.NewLinkedHashMapFrom(resolvedModules), jutil.NewLinkedHashMapFrom(resolvedClosureModules))
}

// moduleMapCreator_difference is
// Sets.difference(a.keySet(), b.keySet()).immutableCopy(): the keys of a
// that are not keys of b, in a's iteration order.
//
// Go-only: language-adapter — Guava Sets.difference.
func moduleMapCreator_difference(
	a *jutil.LinkedHashMap[string, UnresolvedModule], b *jutil.LinkedHashMap[string, *Module]) []string {
	var out []string
	for k := range a.Keys() {
		if !b.ContainsKey(k) {
			out = append(out, k)
		}
	}
	return out
}

// moduleMapCreator_containsAllKeys is a.keySet().containsAll(b.keySet()).
//
// Go-only: language-adapter — Set.containsAll.
func moduleMapCreator_containsAllKeys(
	a *jutil.LinkedHashMap[string, *Module], b *jutil.LinkedHashMap[string, UnresolvedModule]) bool {
	for k := range b.Keys() {
		if !a.ContainsKey(k) {
			return false
		}
	}
	return true
}

// Java: ModuleMapCreator.process(ModuleMetadata)
func (c *ModuleMapCreator) ProcessModuleMetadata(moduleMetadata *ModuleMetadataMap_ModuleMetadata) {
	var processor ModuleMapCreator_ModuleProcessor
	switch moduleMetadata.ModuleType() {
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		processor = c.esModuleProcessor
	case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		if moduleMetadata.GoogNamespaces().Size() == 1 {
			processor = c.closureModuleProcessor
		} else {
			// this indicates some malformed Closure module. We should already
			// have reported an error.
			processor = c.nonEsModuleProcessor
		}
	default:
		processor = c.nonEsModuleProcessor
	}
	module :=
		processor.Process(moduleMetadata, moduleMetadata.Path(), moduleMetadata.RootNode())

	// Have to use module names as keys because path "resolution"
	// (ModuleLoader) "respects" leading slashes. Meaning that if you look up
	// "file.js" and "/file.js" you'll get different paths back. But they'll
	// have the same module name.
	if moduleMetadata.Path() != nil {
		c.unresolvedModules.Put(moduleMetadata.Path().ToModuleName(), module)
	}

	for namespace := range moduleMetadata.GoogNamespaces().All() {
		c.unresolvedModulesByClosureNamespace.Put(namespace, module)
	}
}

// Java: ModuleMapCreator.process(Node, Node)
func (c *ModuleMapCreator) Process(externs *rhino.Node, root *rhino.Node) {
	c.compiler.SetModuleMap(c.Create())
}
