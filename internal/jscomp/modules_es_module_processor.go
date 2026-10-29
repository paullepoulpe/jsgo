// Port of com/google/javascript/jscomp/modules/EsModuleProcessor.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// EsModuleProcessor collects information related to and resolves ES
// imports and exports. Also performs several ES module related checks.
type EsModuleProcessor struct {
	compiler             AbstractCompiler
	currentModuleBuilder *esUnresolvedModuleBuilder
	metadata             *ModuleMetadataMap_ModuleMetadata
}

var _ NodeTraversal_Callback = (*EsModuleProcessor)(nil)
var _ moduleProcessor = (*EsModuleProcessor)(nil)

var (
	// EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION occurs when multiple
	// `export * from` statements export the same name without a local override.
	EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION = DiagnosticType_Warning(
		"JSC_AMBIGUOUS_EXPORT_DEFINITION", "The export \"{0}\" is ambiguous.")

	// Java: EsModuleProcessor.CYCLIC_EXPORT_DEFINITION
	EsModuleProcessor_CYCLIC_EXPORT_DEFINITION = DiagnosticType_Error(
		"JSC_CYCLIC_EXPORT_DEFINITION", "Cyclic export detected while resolving name \"{0}\".")

	// Java: EsModuleProcessor.DUPLICATE_EXPORT
	EsModuleProcessor_DUPLICATE_EXPORT = DiagnosticType_Error("JSC_DUPLICATE_EXPORT", "Duplicate export of \"{0}\".")

	// Java: EsModuleProcessor.IMPORTED_AMBIGUOUS_EXPORT
	EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT = DiagnosticType_Error(
		"JSC_IMPORTED_AMBIGUOUS_EXPORT", "The requested name \"{0}\" is ambiguous.")

	// Java: EsModuleProcessor.NAMESPACE_IMPORT_CANNOT_USE_STAR
	EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR = DiagnosticType_Error(
		"JSC_NAMESPACE_IMPORT_CANNOT_USE_STAR",
		"Namespace imports ('goog:some.Namespace') cannot use import * as. "+
			"Did you mean to import {0} from ''{1}'';?")

	// Java: EsModuleProcessor.CANNOT_PATH_IMPORT_CLOSURE_FILE
	EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE = DiagnosticType_Error(
		"JSC_CANNOT_PATH_IMPORT_CLOSURE_FILE",
		"Cannot import Closure files by path. Use either import 'goog:namespace' or"+
			" goog.require('namespace')")
)

// esFindMutableExports marks all exports that are mutated in an inner scope as mutable.
type esFindMutableExports struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	exports          []*Export
	exportsByLocalID map[uint32][]*Export
}

func newEsFindMutableExports(exports []*Export) *esFindMutableExports {
	f := &esFindMutableExports{exports: exports}
	if len(exports) > 0 {
		f.exportsByLocalID = make(map[uint32][]*Export, len(exports))
		for _, e := range exports {
			id := rhino.RhinoStringPool_InternID(e.LocalName())
			f.exportsByLocalID[id] = append(f.exportsByLocalID[id], e)
		}
	}
	return f
}

func (f *esFindMutableExports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsName() || len(f.exportsByLocalID) == 0 {
		return
	}
	exports := f.exportsByLocalID[n.GetStringID()]
	if len(exports) == 0 {
		return
	}
	scope := t.GetScope()
	if !NodeUtil_IsLValue(n) || scope.GetClosestHoistScope().IsModuleScope() {
		return
	}
	var_ := scope.GetVar(n.GetString())
	if var_ == nil || !var_.GetScope().(*Scope).IsModuleScope() {
		return
	}
	for idx, e := range exports {
		i := slices.Index(f.exports, e)
		mutated := e.MutatedCopy()
		f.exports[i] = mutated
		exports[idx] = mutated
	}
}

type esUnresolvedModuleBuilder struct {
	outer              *EsModuleProcessor
	path               *ModuleLoader_ModulePath
	root               rhino.Node
	importNames        []string
	importsByLocalName map[string]*Import
	exports            []*Export
	exportedNames      map[string]struct{}
}

func newEsUnresolvedModuleBuilder(outer *EsModuleProcessor, path *ModuleLoader_ModulePath, root rhino.Node) *esUnresolvedModuleBuilder {
	return &esUnresolvedModuleBuilder{
		outer:              outer,
		path:               path,
		root:               root,
		importsByLocalName: make(map[string]*Import),
		exportedNames:      make(map[string]struct{}),
	}
}

func (b *esUnresolvedModuleBuilder) addImport(i *Import) {
	name := i.LocalName()
	if _, exists := b.importsByLocalName[name]; !exists {
		b.importNames = append(b.importNames, name)
	}
	b.importsByLocalName[name] = i
}

func (b *esUnresolvedModuleBuilder) addExport(e *Export) bool {
	if !e.HasExportName() {
		b.exports = append(b.exports, e)
		return true
	}
	name := e.ExportName()
	if _, ok := b.exportedNames[name]; ok {
		return false
	}
	b.exports = append(b.exports, e)
	b.exportedNames[name] = struct{}{}
	return true
}

func (b *esUnresolvedModuleBuilder) build() UnresolvedModule {
	var localExports []*Export
	var indirectExports []*Export
	var starExports []*Export

	for _, ee := range b.exports {
		if !ee.HasModuleRequest() {
			if ee.HasLocalName() {
				if _, isImport := b.importsByLocalName[ee.LocalName()]; isImport {
					indirectExports = append(indirectExports, ee)
					continue
				}
			}
			localExports = append(localExports, ee)
		} else if ee.HasImportName() && ee.ImportName() == "*" {
			starExports = append(starExports, ee)
		} else {
			indirectExports = append(indirectExports, ee)
		}
	}

	NodeTraversal_Traverse(b.outer.compiler, b.root, newEsFindMutableExports(localExports))

	return &unresolvedEsModule{
		outer:              b.outer,
		metadata:           b.outer.metadata,
		path:               b.path,
		importNames:        b.importNames,
		importsByLocalName: b.importsByLocalName,
		localExports:       localExports,
		indirectExports:    indirectExports,
		starExports:        starExports,
		resolvedImports:    make(map[string]*ResolveExportResult, len(b.importsByLocalName)),
		resolvedExports:    make(map[string]*ResolveExportResult, len(b.exports)),
	}
}

type unresolvedEsModule struct {
	outer              *EsModuleProcessor
	metadata           *ModuleMetadataMap_ModuleMetadata
	path               *ModuleLoader_ModulePath
	importNames        []string
	importsByLocalName map[string]*Import
	localExports       []*Export
	indirectExports    []*Export
	starExports        []*Export
	exportedNames      []string
	hasExportedNames   bool
	resolvedImports    map[string]*ResolveExportResult
	resolvedExports    map[string]*ResolveExportResult
	resolved           *Module
}

var _ UnresolvedModule = (*unresolvedEsModule)(nil)

func (m *unresolvedEsModule) Reset() {
	m.resolved = nil
	m.exportedNames = nil
	m.hasExportedNames = false
	clear(m.resolvedImports)
	clear(m.resolvedExports)
}

func (m *unresolvedEsModule) Resolve(moduleRequestResolver ModuleRequestResolver) *Module {
	return m.ResolveWithModuleSpecifier(moduleRequestResolver, "", false)
}

func (m *unresolvedEsModule) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	return m.ResolveExportWithTraces(moduleRequestResolver, "", false, exportName, nil, nil)
}

func (m *unresolvedEsModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, hasModuleSpecifier bool) *Module {
	if m.resolved == nil {
		boundNames := m.getAllResolvedImports(moduleRequestResolver)
		localNameToLocalExport := nonEsEmptyExportMap
		if len(m.localExports) > 0 {
			localNameToLocalExport = NewOrderedStringMapWithCapacity[*Export](len(m.localExports))
			for _, e := range m.localExports {
				localName := e.LocalName()
				localNameToLocalExport.Put(localName, e)
				if localName != Export_DEFAULT_EXPORT_NAME {
					b := m.ResolveExport(moduleRequestResolver, e.ExportName())
					jutil.CheckStateMsg(b.Resolved(), "Cannot have invalid missing own export!")
					if !b.IsAmbiguous() {
						if boundNames == nonEsEmptyBindingMap {
							boundNames = NewOrderedStringMapWithCapacity[*Binding](len(m.localExports))
						}
						boundNames.Put(localName, jutil.CheckNotNull(b.GetBinding()))
					}
				}
			}
		}

		namespace := m.getAllResolvedExports(moduleRequestResolver)
		m.resolved = NewModule(m.metadata, m.path, namespace, boundNames, localNameToLocalExport, "", false)
	}
	return m.resolved
}

func (m *unresolvedEsModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

func (m *unresolvedEsModule) getAllResolvedImports(
	moduleRequestResolver ModuleRequestResolver) *OrderedStringMap[*Binding] {
	if len(m.importNames) == 0 {
		return nonEsEmptyBindingMap
	}
	var imports *OrderedStringMap[*Binding]
	for _, name := range m.importNames {
		b := m.resolveImport(moduleRequestResolver, name)
		if b.Resolved() {
			if imports == nil {
				imports = NewOrderedStringMapWithCapacity[*Binding](len(m.importNames))
			}
			imports.Put(name, jutil.CheckNotNull(b.GetBinding()))
		}
	}
	if imports == nil {
		return nonEsEmptyBindingMap
	}
	return imports
}

func (m *unresolvedEsModule) resolveImportWithSets(
	moduleRequestResolver ModuleRequestResolver,
	name string,
	resolveSet map[ExportTrace]struct{},
	exportStarSet map[UnresolvedModule]struct{}) *ResolveExportResult {
	if res, ok := m.resolvedImports[name]; ok {
		return res
	}
	res := m.resolveImportImpl(moduleRequestResolver, name, resolveSet, exportStarSet)
	m.resolvedImports[name] = res
	return res
}

func (m *unresolvedEsModule) resolveImport(
	moduleRequestResolver ModuleRequestResolver, name string) *ResolveExportResult {
	return m.resolveImportWithSets(moduleRequestResolver, name, nil, nil)
}

func (m *unresolvedEsModule) resolveImportImpl(
	moduleRequestResolver ModuleRequestResolver,
	name string,
	resolveSet map[ExportTrace]struct{},
	exportStarSet map[UnresolvedModule]struct{}) *ResolveExportResult {
	i := m.importsByLocalName[name]
	requested := moduleRequestResolver.Resolve(i)
	if requested == nil {
		return ResolveExportResult_ERROR
	}

	importStar := i.ImportName() == "*"
	reqMeta := requested.Metadata()
	if importStar || (i.ImportName() == Export_DEFAULT && (reqMeta.IsGoogProvide() || reqMeta.IsGoogModule())) {
		isGoogSpec := GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest())
		if !isGoogSpec && (reqMeta.IsGoogModule() || reqMeta.IsGoogProvide()) {
			srcName := m.path.String()
			m.outer.compiler.Report(
				JSError_MakeWithSourceLocation(
					&srcName,
					i.ImportNode().GetLineno(),
					i.ImportNode().GetCharno(),
					EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE,
					i.LocalName(),
					i.ModuleRequest()))
			return ResolveExportResult_ERROR
		}
		if importStar && isGoogSpec {
			srcName := m.path.String()
			m.outer.compiler.Report(
				JSError_MakeWithSourceLocation(
					&srcName,
					i.ImportNode().GetLineno(),
					i.ImportNode().GetCharno(),
					EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR,
					i.LocalName(),
					i.ModuleRequest()))
			return ResolveExportResult_ERROR
		}
		var closureNamespace string
		if isGoogSpec {
			closureNamespace = GoogEsImports_GetClosureIdFromGoogImportSpecifier(i.ModuleRequest())
		}
		return ResolveExportResult_Of(Binding_FromImportStar(reqMeta, closureNamespace, isGoogSpec, i.NameNode()))
	}

	result := requested.ResolveExportWithTraces(
		moduleRequestResolver,
		i.ModuleRequest(),
		true,
		i.ImportName(),
		resolveSet,
		exportStarSet)
	if !result.Found() && !result.HadError() {
		srcName := m.path.String()
		m.outer.compiler.Report(
			JSError_MakeWithSourceLocation(
				&srcName,
				i.ImportNode().GetLineno(),
				i.ImportNode().GetCharno(),
				ModuleMapCreator_DOES_NOT_HAVE_EXPORT,
				i.ImportName()))
		return ResolveExportResult_ERROR
	} else if result.IsAmbiguous() {
		srcName := m.path.String()
		m.outer.compiler.Report(
			JSError_MakeWithSourceLocation(
				&srcName,
				i.ImportNode().GetLineno(),
				i.ImportNode().GetCharno(),
				EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT,
				i.ImportName()))
		return ResolveExportResult_ERROR
	}
	forSourceInfo := i.NameNode()
	if forSourceInfo.IsNil() {
		forSourceInfo = i.ImportNode()
	}
	return result.Copy(forSourceInfo, Binding_CreatedBy_IMPORT)
}

func (m *unresolvedEsModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) []string {
	if !m.hasExportedNames {
		var visited map[UnresolvedModule]struct{}
		if len(m.starExports) > 0 {
			visited = make(map[UnresolvedModule]struct{})
		}
		m.exportedNames = m.GetExportedNamesWithVisited(moduleRequestResolver, visited)
		m.hasExportedNames = true
	}
	return m.exportedNames
}

func (m *unresolvedEsModule) hasExportedName(moduleRequestResolver ModuleRequestResolver, name string) bool {
	_, ok := slices.BinarySearchFunc(m.GetExportedNames(moduleRequestResolver), name, jutil.StringCompareTo)
	return ok
}

func (m *unresolvedEsModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited map[UnresolvedModule]struct{}) []string {
	if visited != nil {
		if _, seen := visited[m]; seen {
			return nil
		}
		visited[m] = struct{}{}
	} else if len(m.starExports) > 0 {
		visited = map[UnresolvedModule]struct{}{m: {}}
	}

	capHint := len(m.localExports) + len(m.indirectExports)
	if capHint == 0 && len(m.starExports) == 0 {
		return nil
	}
	names := make([]string, 0, capHint)
	for _, e := range m.localExports {
		names = append(names, e.ExportName())
	}
	for _, e := range m.indirectExports {
		names = append(names, e.ExportName())
	}
	for _, e := range m.starExports {
		requested := moduleRequestResolver.ResolveWithExport(e)
		if requested != nil {
			if requested.Metadata().IsEs6Module() {
				for _, n := range requested.GetExportedNamesWithVisited(moduleRequestResolver, visited) {
					if n != Export_DEFAULT {
						names = append(names, n)
					}
				}
			} else {
				m.outer.compiler.Report(
					JSError_MakeWithNode(
						e.ExportNode(),
						TranspilationUtil_CANNOT_CONVERT_YET,
						"Wildcard export for non-ES module"))
			}
		}
	}

	slices.SortFunc(names, jutil.StringCompareTo)
	return slices.Compact(names)
}

func (m *unresolvedEsModule) getAllResolvedExports(
	moduleRequestResolver ModuleRequestResolver) *OrderedStringMap[*Binding] {
	exportedNames := m.GetExportedNames(moduleRequestResolver)
	if len(exportedNames) == 0 {
		return nonEsEmptyBindingMap
	}
	var exports *OrderedStringMap[*Binding]
	for _, name := range exportedNames {
		b := m.ResolveExport(moduleRequestResolver, name)
		jutil.CheckStateMsg(b.Found(), "Cannot have invalid own export.")
		if b.Resolved() {
			if exports == nil {
				exports = NewOrderedStringMapWithCapacity[*Binding](len(exportedNames))
			}
			exports.Put(name, jutil.CheckNotNull(b.GetBinding()))
		} else if b.IsAmbiguous() {
			srcName := m.path.String()
			m.outer.compiler.Report(JSError_MakeWithSourceLocation(
				&srcName, -1, -1, EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION, name))
		}
	}
	if exports == nil {
		return nonEsEmptyBindingMap
	}
	return exports
}

func (m *unresolvedEsModule) resolveExportWithSets(
	moduleRequestResolver ModuleRequestResolver,
	exportName string,
	resolveSet map[ExportTrace]struct{},
	exportStarSet map[UnresolvedModule]struct{}) *ResolveExportResult {
	if !m.hasExportedName(moduleRequestResolver, exportName) {
		return ResolveExportResult_NOT_FOUND
	}

	if resolveSet != nil {
		trace := ExportTrace_Create(m, exportName)
		if _, seen := resolveSet[trace]; seen {
			srcName := m.path.String()
			m.outer.compiler.Report(JSError_MakeWithSourceLocation(
				&srcName, 0, 0, EsModuleProcessor_CYCLIC_EXPORT_DEFINITION, exportName))
			return ResolveExportResult_ERROR
		}
		resolveSet[trace] = struct{}{}
	}

	for _, e := range m.localExports {
		if e.HasExportName() && exportName == e.ExportName() {
			forSourceInfo := e.ExportNode()
			if !e.NameNode().IsNil() {
				forSourceInfo = e.NameNode()
			}
			return ResolveExportResult_Of(Binding_From(e, forSourceInfo))
		}
	}

	if resolveSet == nil {
		resolveSet = map[ExportTrace]struct{}{ExportTrace_Create(m, exportName): {}}
	}

	for _, e := range m.indirectExports {
		if e.HasExportName() && exportName == e.ExportName() {
			if e.HasLocalName() {
				if _, isImport := m.importsByLocalName[e.LocalName()]; isImport {
					return m.resolveImportWithSets(moduleRequestResolver, e.LocalName(), resolveSet, exportStarSet).
						Copy(e.NameNode(), Binding_CreatedBy_EXPORT)
				}
			}
			requested := moduleRequestResolver.ResolveWithExport(e)
			if requested == nil {
				return ResolveExportResult_ERROR
			}
			result := requested.ResolveExportWithTraces(
				moduleRequestResolver,
				e.ModuleRequest(),
				e.HasModuleRequest(),
				e.ImportName(),
				resolveSet,
				exportStarSet)
			if !result.Found() && !result.HadError() {
				srcName := m.path.String()
				m.outer.compiler.Report(
					JSError_MakeWithSourceLocation(
						&srcName,
						e.ExportNode().GetLineno(),
						e.ExportNode().GetCharno(),
						ModuleMapCreator_DOES_NOT_HAVE_EXPORT,
						e.ImportName()))
				return ResolveExportResult_ERROR
			} else if result.IsAmbiguous() {
				srcName := m.path.String()
				m.outer.compiler.Report(
					JSError_MakeWithSourceLocation(
						&srcName,
						e.ExportNode().GetLineno(),
						e.ExportNode().GetCharno(),
						EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT,
						e.ImportName()))
			}
			return result.Copy(e.NameNode(), Binding_CreatedBy_EXPORT)
		}
	}

	jutil.CheckStateMsg(exportName != Export_DEFAULT, "Default export cannot come from export *.")

	if exportStarSet == nil {
		exportStarSet = map[UnresolvedModule]struct{}{m: {}}
	} else {
		if _, seen := exportStarSet[m]; seen {
			srcName := m.path.String()
			m.outer.compiler.Report(JSError_MakeWithSourceLocation(
				&srcName, -1, -1, EsModuleProcessor_CYCLIC_EXPORT_DEFINITION, exportName))
			return ResolveExportResult_ERROR
		}
		exportStarSet[m] = struct{}{}
	}

	var starResolution *ResolveExportResult
	for _, e := range m.starExports {
		requested := moduleRequestResolver.ResolveWithExport(e)
		if requested == nil {
			return ResolveExportResult_ERROR
		}
		if _, has := slices.BinarySearchFunc(requested.GetExportedNames(moduleRequestResolver), exportName, jutil.StringCompareTo); has {
			resolution := requested.ResolveExportWithTraces(
				moduleRequestResolver, e.ModuleRequest(), e.HasModuleRequest(), exportName, resolveSet, exportStarSet)
			if resolution.HadError() || resolution.IsAmbiguous() {
				return resolution
			}
			if starResolution == nil {
				starResolution = resolution.Copy(e.ExportNode(), Binding_CreatedBy_EXPORT)
			} else if starResolution != resolution {
				return ResolveExportResult_AMBIGUOUS
			}
		}
	}

	if starResolution == nil {
		return ResolveExportResult_ERROR
	}
	return starResolution
}

func (m *unresolvedEsModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier string,
	hasModuleSpecifier bool,
	exportName string,
	resolveSet map[ExportTrace]struct{},
	exportStarSet map[UnresolvedModule]struct{}) *ResolveExportResult {
	if res, ok := m.resolvedExports[exportName]; ok {
		return res
	}
	res := m.resolveExportWithSets(moduleRequestResolver, exportName, resolveSet, exportStarSet)
	m.resolvedExports[exportName] = res
	return res
}

// Java: EsModuleProcessor(AbstractCompiler)
func NewEsModuleProcessor(compiler AbstractCompiler) *EsModuleProcessor {
	return &EsModuleProcessor{compiler: compiler}
}

// Java: EsModuleProcessor.process(ModuleMetadata, ModulePath, Node)
func (p *EsModuleProcessor) Process(
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script rhino.Node) UnresolvedModule {
	p.metadata = metadata
	p.currentModuleBuilder = newEsUnresolvedModuleBuilder(p, path, script)
	NodeTraversal_Traverse(p.compiler, script, p)
	m := p.currentModuleBuilder.build()
	p.metadata = nil
	p.currentModuleBuilder = nil
	return m
}

// Java: EsModuleProcessor.shouldTraverse(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT, rhino.Token_IMPORT:
		return true
	default:
		return false
	}
}

// Java: EsModuleProcessor.visit(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_EXPORT:
		p.visitExport(t, n)
	case rhino.Token_IMPORT:
		p.visitImport(t, n)
	}
}

func (p *EsModuleProcessor) visitExportAllFrom(t *NodeTraversal, export ast.Export) {
	source, ok := export.Source()
	jutil.CheckState(ok)
	p.currentModuleBuilder.addExport(
		Export_Builder_().
			ModuleRequest(source.Value()).
			ImportName("*").
			ModulePath(t.GetInput().GetPath()).
			ExportNode(export.Node()).
			ModuleMetadata(p.metadata).
			Build())
}

func (p *EsModuleProcessor) visitExportDefault(t *NodeTraversal, export ast.Export) {
	body := export.Body().Node()
	name := Export_DEFAULT_EXPORT_NAME
	if fn, ok := ast.AsFunction(body); ok {
		if fnName, hasName := fn.Name(); hasName {
			name = fnName.Ident()
		}
	} else if cls, ok := ast.AsClass(body); ok {
		if clsName, hasName := cls.Name(); hasName {
			name = clsName.Ident()
		}
	}

	if !p.currentModuleBuilder.addExport(
		Export_Builder_().
			ExportName(Export_DEFAULT).
			LocalName(name).
			ModulePath(t.GetInput().GetPath()).
			ExportNode(export.Node()).
			ModuleMetadata(p.metadata).
			Build()) {
		t.Report(export.Node(), EsModuleProcessor_DUPLICATE_EXPORT, Export_DEFAULT)
	}
}

func (p *EsModuleProcessor) visitExportFrom(t *NodeTraversal, export ast.Export, specs ast.ExportSpecs, moduleRequest string) {
	for spec := range specs.Specs() {
		exportedName := spec.Exported().Ident()
		if !p.currentModuleBuilder.addExport(
			Export_Builder_().
				ExportName(exportedName).
				ModuleRequest(moduleRequest).
				ImportName(spec.Local().Ident()).
				ModulePath(t.GetInput().GetPath()).
				ExportNode(export.Node()).
				NameNode(spec.Local().Node()).
				ModuleMetadata(p.metadata).
				Build()) {
			t.Report(export.Node(), EsModuleProcessor_DUPLICATE_EXPORT, exportedName)
		}
	}
}

func (p *EsModuleProcessor) visitExportSpecs(t *NodeTraversal, export ast.Export, specs ast.ExportSpecs) {
	for spec := range specs.Specs() {
		exportedName := spec.Exported().Ident()
		if !p.currentModuleBuilder.addExport(
			Export_Builder_().
				ExportName(exportedName).
				LocalName(spec.Local().Ident()).
				ModulePath(t.GetInput().GetPath()).
				ExportNode(export.Node()).
				NameNode(spec.Local().Node()).
				ModuleMetadata(p.metadata).
				Build()) {
			t.Report(export.Node(), EsModuleProcessor_DUPLICATE_EXPORT, exportedName)
		}
	}
}

func (p *EsModuleProcessor) addLocalDeclExport(t *NodeTraversal, name string, nameNode rhino.Node, exportNode rhino.Node) {
	if !p.currentModuleBuilder.addExport(
		Export_Builder_().
			ExportName(name).
			LocalName(name).
			ModulePath(t.GetInput().GetPath()).
			ExportNode(exportNode).
			NameNode(nameNode).
			ModuleMetadata(p.metadata).
			Build()) {
		t.Report(exportNode, EsModuleProcessor_DUPLICATE_EXPORT, name)
	}
}

func (p *EsModuleProcessor) visitExport(t *NodeTraversal, exportNode rhino.Node) {
	export := ast.MustExport(exportNode)
	if export.IsAllFrom() {
		p.visitExportAllFrom(t, export)
	} else if export.IsDefault() {
		p.visitExportDefault(t, export)
	} else if source, hasSource := export.Source(); hasSource {
		p.visitExportFrom(t, export, export.Body().MustExportSpecs(), source.Value())
	} else if specs, ok := export.Body().AsExportSpecs(); ok {
		p.visitExportSpecs(t, export, specs)
	} else {
		declaration := export.Body().Node()
		if nameDecl, ok := ast.AsNameDeclaration(declaration); ok {
			NodeUtil_VisitLhsNodesInNode(nameDecl.Node(), func(lhs rhino.Node) {
				lhsName := ast.MustName(lhs)
				p.addLocalDeclExport(t, lhsName.Ident(), lhsName.Node(), exportNode)
			})
		} else {
			var nameNode ast.Name
			if fn, ok := ast.AsFunction(declaration); ok {
				nameNode = fn.NameNode()
			} else {
				nameNode = ast.MustName(ast.MustClass(declaration).NameNode())
			}
			p.addLocalDeclExport(t, nameNode.Ident(), nameNode.Node(), exportNode)
		}
	}
}

func (p *EsModuleProcessor) visitImport(t *NodeTraversal, importNode rhino.Node) {
	imp := ast.MustImport(importNode)
	moduleRequest := imp.Source().Value()
	modulePath := t.GetInput().GetPath()

	if defaultName, ok := imp.Default(); ok {
		p.currentModuleBuilder.addImport(
			NewImport(moduleRequest, Export_DEFAULT, defaultName.Ident(), modulePath, importNode, defaultName.Node()))
	}

	if bindings, ok := imp.Bindings(); ok {
		if specs, ok := bindings.AsImportSpecs(); ok {
			for spec := range specs.Specs() {
				p.currentModuleBuilder.addImport(
					NewImport(moduleRequest, spec.Imported().Ident(), spec.Local().Ident(), modulePath, importNode, spec.Local().Node()))
			}
		} else if star, ok := bindings.AsImportStar(); ok {
			p.currentModuleBuilder.addImport(
				NewImport(moduleRequest, "*", star.Ident(), modulePath, importNode, star.Node()))
		}
	}
}
