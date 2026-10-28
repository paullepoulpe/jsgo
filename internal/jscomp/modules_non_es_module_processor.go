// Port of com/google/javascript/jscomp/modules/NonEsModuleProcessor.java @ closure-compiler v20261021

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
	if m.metadata.IsCommonJs() && exportName != "default" {
		return ResolveExportResult_NOT_FOUND
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

func (m *nonEsModule) Reset() {}

func (m *nonEsModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module {
	var namespace *string
	if moduleSpecifier != nil && GoogEsImports_IsGoogImportSpecifier(*moduleSpecifier) {
		ns := GoogEsImports_GetClosureIdFromGoogImportSpecifier(*moduleSpecifier)
		namespace = &ns
	}
	jutil.CheckNotNullMsg(m.metadata, "Null metadata")
	return &Module{
		metadata:               m.metadata,
		path:                   m.path,
		namespace:              nonEsEmptyBindingMap,
		boundNames:             nonEsEmptyBindingMap,
		localNameToLocalExport: nonEsEmptyExportMap,
		closureNamespace:       namespace,
	}
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
