// Port of com/google/javascript/jscomp/modules/NonEsModuleProcessor.java @ closure-compiler v20261015

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
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	nonEsEmptyBindingMap = jutil.NewLinkedHashMap[string, *Binding]()
	nonEsEmptyExportMap  = jutil.NewLinkedHashMap[string, *Export]()
)

// NonEsModuleProcessor is the catch all module processor for non-ES and
// non-goog modules that doesn't do any scanning of exports but instead will
// always treat every name as `munged.name`.
//
// This includes goog.provide'd names and CommonJS modules.
//
// TODO(johnplaisted): It may not be too hard to actually scan these modules
// and figure out their exported keys. At least, assuming we're rigid with
// what we allow (no dynamic exports).
type NonEsModuleProcessor struct {
	nonEsSlab  []NonEsModuleProcessor_NonEsModule
	moduleSlab []Module
}

var _ ModuleMapCreator_ModuleProcessor = (*NonEsModuleProcessor)(nil)

// Java: NonEsModuleProcessor()
func NewNonEsModuleProcessor() *NonEsModuleProcessor {
	return &NonEsModuleProcessor{}
}

// NonEsModuleProcessor_NonEsModule is NonEsModuleProcessor.NonEsModule.
type NonEsModuleProcessor_NonEsModule struct {
	UnresolvedModuleBase

	processor  *NonEsModuleProcessor
	metadata   *ModuleMetadataMap_ModuleMetadata
	path       *ModuleLoader_ModulePath
	scriptNode *rhino.Node
}

// Java: NonEsModuleProcessor.NonEsModule(ModuleMetadata, ModulePath, Node)
func NewNonEsModuleProcessor_NonEsModule(
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, scriptNode *rhino.Node) *NonEsModuleProcessor_NonEsModule {
	m := &NonEsModuleProcessor_NonEsModule{}
	m.InitUnresolvedModuleBase(m)
	m.metadata = metadata
	m.path = path
	m.scriptNode = scriptNode
	return m
}

// Java: NonEsModuleProcessor.NonEsModule.resolveExport(ModuleRequestResolver, String)
func (m *NonEsModuleProcessor_NonEsModule) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: NonEsModuleProcessor.NonEsModule.resolveExport(ModuleRequestResolver, String, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *NonEsModuleProcessor_NonEsModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier *string,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	var namespace *string
	if moduleSpecifier != nil && GoogEsImports_IsGoogImportSpecifier(*moduleSpecifier) {
		ns := GoogEsImports_GetClosureIdFromGoogImportSpecifier(*moduleSpecifier)
		namespace = &ns
	} else if m.metadata.IsGoogProvide() {
		namespace = moduleSpecifier
	}
	if m.metadata.IsCommonJs() {
		// Currently we don't scan require()s, so this only gets hit if an ES
		// module imports a CJS module. If so then this module should only have
		// a default export.
		if exportName != "default" {
			return ResolveExportResult_NOT_FOUND
		}
	}
	return ResolveExportResult_Of(
		Binding_From(
			Export_Builder_().
				ExportName(&exportName).
				ModuleMetadata(m.metadata).
				ModulePath(m.path).
				ClosureNamespace(namespace).
				Build(),
			m.scriptNode))
}

// Java: NonEsModuleProcessor.NonEsModule.reset()
func (m *NonEsModuleProcessor_NonEsModule) Reset() {}

// Java: NonEsModuleProcessor.NonEsModule.resolve(ModuleRequestResolver, String)
func (m *NonEsModuleProcessor_NonEsModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module {
	var namespace *string
	if moduleSpecifier != nil && GoogEsImports_IsGoogImportSpecifier(*moduleSpecifier) {
		ns := GoogEsImports_GetClosureIdFromGoogImportSpecifier(*moduleSpecifier)
		namespace = &ns
	}
	jutil.CheckNotNullMsg(m.metadata, "Null metadata")
	if m.processor != nil {
		if len(m.processor.moduleSlab) == cap(m.processor.moduleSlab) {
			newCap := max(16, min(max(cap(m.processor.moduleSlab)*2, 64), jscompMaxSmallAlloc/int(unsafe.Sizeof(Module{}))))
			m.processor.moduleSlab = make([]Module, 0, newCap)
		}
		idx := len(m.processor.moduleSlab)
		m.processor.moduleSlab = m.processor.moduleSlab[:idx+1]
		mod := &m.processor.moduleSlab[idx]
		mod.metadata = m.metadata
		mod.path = m.path
		mod.namespace = nonEsEmptyBindingMap
		mod.boundNames = nonEsEmptyBindingMap
		mod.localNameToLocalExport = nonEsEmptyExportMap
		mod.closureNamespace = namespace
		return mod
	}
	return &Module{
		metadata:               m.metadata,
		path:                   m.path,
		namespace:              nonEsEmptyBindingMap,
		boundNames:             nonEsEmptyBindingMap,
		localNameToLocalExport: nonEsEmptyExportMap,
		closureNamespace:       namespace,
	}
}

// Java: NonEsModuleProcessor.NonEsModule.metadata()
func (m *NonEsModuleProcessor_NonEsModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

// Java: NonEsModuleProcessor.NonEsModule.getExportedNames(ModuleRequestResolver)
func (m *NonEsModuleProcessor_NonEsModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: NonEsModuleProcessor.NonEsModule.getExportedNames(ModuleRequestResolver, Set<UnresolvedModule>)
func (m *NonEsModuleProcessor_NonEsModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: NonEsModuleProcessor.process(ModuleMetadata, ModulePath, Node)
func (p *NonEsModuleProcessor) Process(
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script *rhino.Node) UnresolvedModule {
	if len(p.nonEsSlab) == cap(p.nonEsSlab) {
		newCap := max(16, min(max(cap(p.nonEsSlab)*2, 64), jscompMaxSmallAlloc/int(unsafe.Sizeof(NonEsModuleProcessor_NonEsModule{}))))
		p.nonEsSlab = make([]NonEsModuleProcessor_NonEsModule, 0, newCap)
	}
	idx := len(p.nonEsSlab)
	p.nonEsSlab = p.nonEsSlab[:idx+1]
	m := &p.nonEsSlab[idx]
	m.InitUnresolvedModuleBase(m)
	m.processor = p
	m.metadata = metadata
	m.path = path
	m.scriptNode = script
	return m
}
