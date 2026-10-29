// Port of com/google/javascript/jscomp/modules/ModuleMapCreator.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/ModuleMap.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/ModuleRequestResolver.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/GoogEsImports.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/ExportTrace.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/NonEsModuleProcessor.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/ResolveExportResult.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/modules/UnresolvedModule.java @ closure-compiler v20261021

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

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type orderedUnresolvedMap struct {
	keys  []string
	byKey map[string]UnresolvedModule
}

func (m *orderedUnresolvedMap) reset() {
	m.keys = m.keys[:0]
	if m.byKey == nil {
		m.byKey = make(map[string]UnresolvedModule)
	} else {
		clear(m.byKey)
	}
}

func (m *orderedUnresolvedMap) put(k string, v UnresolvedModule) {
	if _, exists := m.byKey[k]; !exists {
		m.keys = append(m.keys, k)
	}
	m.byKey[k] = v
}

// ModuleMapCreator creates a ModuleMap.
type ModuleMapCreator struct {
	compiler                            AbstractCompiler
	esModuleProcessor                   *EsModuleProcessor
	closureModuleProcessor              *ClosureModuleProcessor
	nonEsModuleProcessor                *NonEsModuleProcessor
	unresolvedModules                   orderedUnresolvedMap
	unresolvedModulesByClosureNamespace orderedUnresolvedMap
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

type moduleRequestResolverImpl struct {
	outer *ModuleMapCreator
}

var _ ModuleRequestResolver = (*moduleRequestResolverImpl)(nil)

type fallbackMissingNonClosureModule struct {
	UnresolvedModuleBase

	path     *ModuleLoader_ModulePath
	metadata *ModuleMetadataMap_ModuleMetadata
}

func newFallbackMissingNonClosureModule(path *ModuleLoader_ModulePath) UnresolvedModule {
	metadata := &ModuleMetadataMap_ModuleMetadata{
		path:       path,
		moduleType: ModuleMetadataMap_ModuleType_ES6_MODULE,
	}
	m := &fallbackMissingNonClosureModule{path: path, metadata: metadata}
	m.InitUnresolvedModuleBase(m)
	return m
}

func (m *fallbackMissingNonClosureModule) Reset() {}

func (m *fallbackMissingNonClosureModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, hasModuleSpecifier bool) *Module {
	return NewModule(m.metadata, m.path, nonEsEmptyBindingMap, nonEsEmptyBindingMap, nonEsEmptyExportMap, "", false)
}

func (m *fallbackMissingNonClosureModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

func (m *fallbackMissingNonClosureModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	return nonEsEmptyExportedSet
}

func (m *fallbackMissingNonClosureModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	return nonEsEmptyExportedSet
}

func (m *fallbackMissingNonClosureModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier string,
	hasModuleSpecifier bool,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	return ResolveExportResult_Of(
		Binding_From(
			Export_Builder_().
				LocalName(exportName).
				ModuleMetadata(m.metadata).
				ModulePath(m.path).
				Build(),
			rhino.NilNode))
}

func (r *moduleRequestResolverImpl) fallbackForMissingClosureModule(namespace string) UnresolvedModule {
	md := &ModuleMetadataMap_ModuleMetadata{
		moduleType:  ModuleMetadataMap_ModuleType_GOOG_PROVIDE,
		usesClosure: true,
	}
	md.googNamespaces.Add(namespace)
	return r.outer.nonEsModuleProcessor.Process(md, nil, rhino.NilNode)
}

func (r *moduleRequestResolverImpl) Resolve(i *Import) UnresolvedModule {
	if i.ModulePath() == nil {
		return r.resolveForClosure(i.ModuleRequest())
	}
	return r.resolveWithModulePath(i.ModuleRequest(), i.ModulePath(), i.ImportNode())
}

func (r *moduleRequestResolverImpl) ResolveWithExport(e *Export) UnresolvedModule {
	return r.resolveWithModulePath(e.ModuleRequest(), e.ModulePath(), e.ExportNode())
}

func (r *moduleRequestResolverImpl) resolveForClosure(namespace string) UnresolvedModule {
	if m, ok := r.outer.unresolvedModulesByClosureNamespace.byKey[namespace]; ok {
		return m
	}
	m := r.fallbackForMissingClosureModule(namespace)
	r.outer.unresolvedModulesByClosureNamespace.put(namespace, m)
	return m
}

func (r *moduleRequestResolverImpl) resolveWithModulePath(
	moduleRequest string, modulePath *ModuleLoader_ModulePath, forLineInfo rhino.Node) UnresolvedModule {
	if GoogEsImports_IsGoogImportSpecifier(moduleRequest) {
		return r.resolveForClosure(GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleRequest))
	}

	sourceFileName := ""
	if s := forLineInfo.GetSourceFileName(); s != nil {
		sourceFileName = *s
	}
	requestedPath := modulePath.ResolveJsModule(
		moduleRequest,
		sourceFileName,
		forLineInfo.GetLineno(),
		forLineInfo.GetCharno())

	if requestedPath == nil {
		requestedPath = modulePath.ResolveModuleAsPath(moduleRequest)
		modName := requestedPath.ToModuleName()
		if existing, ok := r.outer.unresolvedModules.byKey[modName]; ok {
			return existing
		}
		module := newFallbackMissingNonClosureModule(requestedPath)
		r.outer.unresolvedModules.put(modName, module)
		return module
	}

	return r.outer.unresolvedModules.byKey[requestedPath.ToModuleName()]
}

type moduleProcessor interface {
	Process(metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script rhino.Node) UnresolvedModule
}

// Java: ModuleMapCreator(AbstractCompiler, ModuleMetadataMap)
func NewModuleMapCreator(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *ModuleMapCreator {
	return &ModuleMapCreator{
		compiler:               compiler,
		moduleMetadataMap:      moduleMetadataMap,
		esModuleProcessor:      NewEsModuleProcessor(compiler),
		closureModuleProcessor: NewClosureModuleProcessor(compiler),
		nonEsModuleProcessor:   NewNonEsModuleProcessor(),
	}
}

// Java: ModuleMapCreator.create()
func (c *ModuleMapCreator) Create() *ModuleMap {
	c.unresolvedModules.reset()
	c.unresolvedModulesByClosureNamespace.reset()

	for moduleMetadata := range c.moduleMetadataMap.GetAllModuleMetadata().All() {
		c.processModuleMetadata(moduleMetadata)
	}
	return c.resolve()
}

func (c *ModuleMapCreator) resolve() *ModuleMap {
	requestResolver := &moduleRequestResolverImpl{outer: c}
	resolvedModules := jutil.NewLinkedHashMap[string, *Module]()
	resolvedClosureModules := jutil.NewLinkedHashMap[string, *Module]()

	// Missing references append fallback entries to unresolvedModules.keys during
	// resolution, so advancing an index cursor resolves known and fallback modules
	// in deterministic insertion order in a single linear pass.
	for idx := 0; idx < len(c.unresolvedModules.keys); idx++ {
		key := c.unresolvedModules.keys[idx]
		resolved := c.unresolvedModules.byKey[key].Resolve(requestResolver)
		resolvedModules.Put(key, resolved)
		for namespace := range resolved.Metadata().GoogNamespaces().All() {
			resolvedClosureModules.Put(namespace, resolved)
		}
	}

	for idx := 0; idx < len(c.unresolvedModulesByClosureNamespace.keys); idx++ {
		namespace := c.unresolvedModulesByClosureNamespace.keys[idx]
		if !resolvedClosureModules.ContainsKey(namespace) {
			resolvedClosureModules.Put(
				namespace,
				c.unresolvedModulesByClosureNamespace.byKey[namespace].Resolve(requestResolver))
		}
	}

	c.unresolvedModules.reset()
	c.unresolvedModulesByClosureNamespace.reset()
	return NewModuleMap(resolvedModules, resolvedClosureModules)
}

func (c *ModuleMapCreator) processModuleMetadata(moduleMetadata *ModuleMetadataMap_ModuleMetadata) {
	var processor moduleProcessor
	switch moduleMetadata.ModuleType() {
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		processor = c.esModuleProcessor
	case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		if moduleMetadata.GoogNamespaces().Size() == 1 {
			processor = c.closureModuleProcessor
		} else {
			processor = c.nonEsModuleProcessor
		}
	default:
		processor = c.nonEsModuleProcessor
	}
	module := processor.Process(moduleMetadata, moduleMetadata.Path(), moduleMetadata.RootNode())

	if moduleMetadata.Path() != nil {
		c.unresolvedModules.put(moduleMetadata.Path().ToModuleName(), module)
	}
	for namespace := range moduleMetadata.GoogNamespaces().All() {
		c.unresolvedModulesByClosureNamespace.put(namespace, module)
	}
}

// Java: ModuleMapCreator.process(Node, Node)
func (c *ModuleMapCreator) Process(externs rhino.Node, root rhino.Node) {
	c.compiler.SetModuleMap(c.Create())
}

// ModuleMap holds resolved module information for the compilation.
type ModuleMap struct {
	resolvedModules        *jutil.LinkedHashMap[string, *Module]
	resolvedClosureModules *jutil.LinkedHashMap[string, *Module]
}

// Java: ModuleMap(ImmutableMap<String, Module>, ImmutableMap<String, Module>)
func NewModuleMap(
	resolvedModules *jutil.LinkedHashMap[string, *Module],
	resolvedClosureModules *jutil.LinkedHashMap[string, *Module]) *ModuleMap {
	return &ModuleMap{resolvedModules: resolvedModules, resolvedClosureModules: resolvedClosureModules}
}

// Java: ModuleMap.getModule(String)
func (m *ModuleMap) GetModule(moduleName string) *Module {
	return m.resolvedModules.Get(moduleName)
}

// Java: ModuleMap.getModule(ModuleLoader.ModulePath)
func (m *ModuleMap) GetModuleByPath(path *ModuleLoader_ModulePath) *Module {
	return m.resolvedModules.Get(path.ToModuleName())
}

// Java: ModuleMap.getModulesByPath()
func (m *ModuleMap) GetModulesByPath() *jutil.LinkedHashMap[string, *Module] {
	return m.resolvedModules
}

// Java: ModuleMap.getModulesByClosureNamespace()
func (m *ModuleMap) GetModulesByClosureNamespace() *jutil.LinkedHashMap[string, *Module] {
	return m.resolvedClosureModules
}

// Java: ModuleMap.getClosureModule(String)
func (m *ModuleMap) GetClosureModule(namespace string) *Module {
	return m.resolvedClosureModules.Get(namespace)
}

// Java: ModuleMap.emptyForTesting()
func ModuleMap_EmptyForTesting() *ModuleMap {
	return NewModuleMap(jutil.NewLinkedHashMap[string, *Module](), jutil.NewLinkedHashMap[string, *Module]())
}

// ModuleRequestResolver resolves requests for other modules.
type ModuleRequestResolver interface {
	Resolve(i *Import) UnresolvedModule
	ResolveWithExport(e *Export) UnresolvedModule
}

const googEsImports_GOOG_IMPORT_PREFIX = "goog:"

// Java: GoogEsImports.isGoogImportSpecifier(String)
func GoogEsImports_IsGoogImportSpecifier(moduleSpecifier string) bool {
	return strings.HasPrefix(moduleSpecifier, googEsImports_GOOG_IMPORT_PREFIX)
}

// Java: GoogEsImports.getClosureIdFromGoogImportSpecifier(String)
func GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleSpecifier string) string {
	return moduleSpecifier[len(googEsImports_GOOG_IMPORT_PREFIX):]
}

// ExportTrace keeps track of visited (module, exportName) pairs for cycle detection.
type ExportTrace struct {
	module     UnresolvedModule
	exportName string
}

// Java: ExportTrace(UnresolvedModule, String)
func NewExportTrace(module UnresolvedModule, exportName string) ExportTrace {
	if module == nil {
		panic(jutil.NewNullPointerExceptionMsg("module"))
	}
	return ExportTrace{module: module, exportName: exportName}
}

// Java: ExportTrace.module()
func (t ExportTrace) Module() UnresolvedModule { return t.module }

// Java: ExportTrace.exportName()
func (t ExportTrace) ExportName() string { return t.exportName }

// Java: ExportTrace.create(UnresolvedModule, String)
func ExportTrace_Create(module UnresolvedModule, exportName string) ExportTrace {
	return NewExportTrace(module, exportName)
}

// UnresolvedModule is a module which has had some of its imports and
// exports statements scanned but has yet to resolve anything transitively.
//
// Reference equality is expected in ExportTrace, so implementations must be
// pointer types (Go interface equality is then identity).
type UnresolvedModule interface {
	// Reset clears any caching so that Resolve(ModuleRequestResolver) will
	// create a new Module.
	//
	// Java: UnresolvedModule.reset()
	Reset()

	// Java: UnresolvedModule.resolve(ModuleRequestResolver)
	Resolve(moduleRequestResolver ModuleRequestResolver) *Module

	// ResolveWithModuleSpecifier resolves all imports and exports and
	// returns a resolved module. moduleSpecifier is the module specifier
	// that was used to import this module, if resolving an import.
	//
	// Java: UnresolvedModule.resolve(ModuleRequestResolver, String)
	ResolveWithModuleSpecifier(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, hasModuleSpecifier bool) *Module

	// Metadata returns the metadata corresponding to this module.
	//
	// Java: UnresolvedModule.metadata()
	Metadata() *ModuleMetadataMap_ModuleMetadata

	// GetExportedNames returns all names in this module's namespace, sorted
	// per Java's string ordering.
	//
	// Java: UnresolvedModule.getExportedNames(ModuleRequestResolver)
	GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string]

	// GetExportedNamesWithVisited returns all names in this module's
	// namespace, sorted per Java's string ordering. visited is the set used
	// to detect export * cycles.
	//
	// Java: UnresolvedModule.getExportedNames(ModuleRequestResolver, Set<UnresolvedModule>)
	GetExportedNamesWithVisited(moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string]

	// Java: UnresolvedModule.resolveExport(ModuleRequestResolver, String)
	ResolveExport(moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult

	// ResolveExportWithTraces resolves exportName (the name of the export to
	// resolve).
	//
	// Java: UnresolvedModule.resolveExport(ModuleRequestResolver, String, String, Set<ExportTrace>, Set<UnresolvedModule>)
	ResolveExportWithTraces(
		moduleRequestResolver ModuleRequestResolver,
		moduleSpecifier string,
		hasModuleSpecifier bool,
		exportName string,
		resolveSet *jutil.LinkedHashSet[ExportTrace],
		exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult
}

// UnresolvedModuleBase holds the base implementations of UnresolvedModule;
// subclasses embed it and call InitUnresolvedModuleBase.
type UnresolvedModuleBase struct {
	self UnresolvedModule
}

// InitUnresolvedModuleBase is the Go form of super(); self is the object
// being constructed.
func (m *UnresolvedModuleBase) InitUnresolvedModuleBase(self UnresolvedModule) {
	m.self = self
}

// Java: UnresolvedModule.resolve(ModuleRequestResolver)
func (m *UnresolvedModuleBase) Resolve(moduleRequestResolver ModuleRequestResolver) *Module {
	return m.self.ResolveWithModuleSpecifier(moduleRequestResolver, "", false)
}

// ResolveExport resolves exportName; see ResolveExportWithTraces.
//
// Java: UnresolvedModule.resolveExport(ModuleRequestResolver, String)
func (m *UnresolvedModuleBase) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	return m.self.ResolveExportWithTraces(
		moduleRequestResolver,
		"",
		false,
		exportName,
		jutil.NewLinkedHashSet[ExportTrace](),
		jutil.NewLinkedHashSet[UnresolvedModule]())
}

// ResolveExportResult is the result of resolving an export, which can be a
// valid binding, ambiguous, not found, or an error.
type ResolveExportResult struct {
	binding *Binding
	state   resolveExportResult_State
}

// Java: ResolveExportResult.State
type resolveExportResult_State int32

const (
	resolveExportResult_State_RESOLVED resolveExportResult_State = iota
	resolveExportResult_State_AMBIGUOUS
	resolveExportResult_State_NOT_FOUND
	resolveExportResult_State_ERROR
)

// Java: ResolveExportResult(Binding, State)
func newResolveExportResult(binding *Binding, state resolveExportResult_State) *ResolveExportResult {
	return &ResolveExportResult{binding: binding, state: state}
}

// Copy creates a new result that has the given node for the source of the
// binding and given type of binding.
//
// Java: ResolveExportResult.copy(Node, Binding.CreatedBy)
func (r *ResolveExportResult) Copy(sourceNode rhino.Node, createdBy Binding_CreatedBy) *ResolveExportResult {
	rhino.CheckNotNull(sourceNode)
	if r.binding == nil {
		return r
	}
	return newResolveExportResult(r.binding.Copy(sourceNode, createdBy), r.state)
}

// HadError is true if there was an error resolving the export, false
// otherwise.
//
// Java: ResolveExportResult.hadError()
func (r *ResolveExportResult) HadError() bool {
	return r.state == resolveExportResult_State_ERROR
}

// IsAmbiguous is true if the export is ambiguous, false otherwise.
//
// Java: ResolveExportResult.isAmbiguous()
func (r *ResolveExportResult) IsAmbiguous() bool {
	return r.state == resolveExportResult_State_AMBIGUOUS
}

// Resolved is true if the export was successfully resolved, false
// otherwise.
//
// Java: ResolveExportResult.resolved()
func (r *ResolveExportResult) Resolved() bool {
	return r.state == resolveExportResult_State_RESOLVED
}

// Found is true if the export key exists on the given module, even if it is
// ambiguous or had an error.
//
// Java: ResolveExportResult.found()
func (r *ResolveExportResult) Found() bool {
	return r.state != resolveExportResult_State_NOT_FOUND
}

// GetBinding may return nil.
//
// Java: ResolveExportResult.getBinding()
func (r *ResolveExportResult) GetBinding() *Binding {
	return r.binding
}

// ResolveExportResult_AMBIGUOUS: the result of resolving the export was
// ambiguous.
var ResolveExportResult_AMBIGUOUS = newResolveExportResult(nil, resolveExportResult_State_AMBIGUOUS)

// ResolveExportResult_NOT_FOUND: the export was not found because the
// module never exported the key.
var ResolveExportResult_NOT_FOUND = newResolveExportResult(nil, resolveExportResult_State_NOT_FOUND)

// ResolveExportResult_ERROR: there was an error resolving the export.
var ResolveExportResult_ERROR = newResolveExportResult(nil, resolveExportResult_State_ERROR)

// Java: ResolveExportResult.of(Binding)
func ResolveExportResult_Of(binding *Binding) *ResolveExportResult {
	jutil.CheckNotNull(binding)
	return newResolveExportResult(binding, resolveExportResult_State_RESOLVED)
}

var (
	nonEsEmptyBindingMap  = jutil.NewLinkedHashMap[string, *Binding]()
	nonEsEmptyExportMap   = jutil.NewLinkedHashMap[string, *Export]()
	nonEsEmptyExportedSet = jutil.NewLinkedHashSet[string]()
)

// NonEsModuleProcessor is the catch-all module processor for non-ES and
// non-goog modules (including goog.provide'd names and CommonJS modules).
type NonEsModuleProcessor struct{}

var _ moduleProcessor = (*NonEsModuleProcessor)(nil)

// Java: NonEsModuleProcessor()
func NewNonEsModuleProcessor() *NonEsModuleProcessor {
	return &NonEsModuleProcessor{}
}

type nonEsModule struct {
	UnresolvedModuleBase

	metadata   *ModuleMetadataMap_ModuleMetadata
	path       *ModuleLoader_ModulePath
	scriptNode rhino.Node
}

func (m *nonEsModule) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	panic(jutil.NewUnsupportedOperationException())
}

func (m *nonEsModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier string,
	hasModuleSpecifier bool,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	var namespace string
	hasNamespace := false
	if hasModuleSpecifier && GoogEsImports_IsGoogImportSpecifier(moduleSpecifier) {
		namespace = GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleSpecifier)
		hasNamespace = true
	} else if m.metadata.IsGoogProvide() && hasModuleSpecifier {
		namespace = moduleSpecifier
		hasNamespace = true
	}
	if m.metadata.IsCommonJs() && exportName != "default" {
		return ResolveExportResult_NOT_FOUND
	}
	b := Export_Builder_().
		ExportName(exportName).
		ModuleMetadata(m.metadata).
		ModulePath(m.path)
	if hasNamespace {
		b = b.ClosureNamespace(namespace)
	}
	return ResolveExportResult_Of(Binding_From(b.Build(), m.scriptNode))
}

func (m *nonEsModule) Reset() {}

func (m *nonEsModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, hasModuleSpecifier bool) *Module {
	var namespace string
	hasNamespace := false
	if hasModuleSpecifier && GoogEsImports_IsGoogImportSpecifier(moduleSpecifier) {
		namespace = GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleSpecifier)
		hasNamespace = true
	}
	jutil.CheckNotNullMsg(m.metadata, "Null metadata")
	return NewModule(
		m.metadata,
		m.path,
		nonEsEmptyBindingMap,
		nonEsEmptyBindingMap,
		nonEsEmptyExportMap,
		namespace,
		hasNamespace)
}

func (m *nonEsModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

func (m *nonEsModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

func (m *nonEsModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: NonEsModuleProcessor.process(ModuleMetadata, ModulePath, Node)
func (p *NonEsModuleProcessor) Process(
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script rhino.Node) UnresolvedModule {
	m := &nonEsModule{
		metadata:   metadata,
		path:       path,
		scriptNode: script,
	}
	m.InitUnresolvedModuleBase(m)
	return m
}
