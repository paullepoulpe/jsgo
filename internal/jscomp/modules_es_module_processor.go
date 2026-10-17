// Port of com/google/javascript/jscomp/modules/EsModuleProcessor.java @ closure-compiler v20261015

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

// EsModuleProcessor collects information related to and resolves ES
// imports and exports. Also performs several ES module related checks.
//
// This information is stored outside of any Scope because it should never
// be recalculated because it is expensive.
//
// Suggested reading:
// https://www.ecma-international.org/ecma-262/9.0/index.html#sec-source-text-module-records
type EsModuleProcessor struct {
	compiler             AbstractCompiler
	currentModuleBuilder *EsModuleProcessor_UnresolvedModuleBuilder
	metadata             *ModuleMetadataMap_ModuleMetadata
}

var _ NodeTraversal_Callback = (*EsModuleProcessor)(nil)
var _ ModuleMapCreator_ModuleProcessor = (*EsModuleProcessor)(nil)

// EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION: error occurs when there is
// an ambiguous export, which can happen if there are multiple `export *
// from` statements. If two modules that were `export * from`'d have the
// same key as an export that export is now ambiguous. It can be resolved by
// exporting that key explicitly locally.
//
// So if modules "first" and "second" both export "x" then this is
// ambiguous:
//
//	export * from 'first';
//	export * from 'second';
//
// However we can resolve it by explicitly exporting an "x" locally:
//
//	// Note: any one of these is a solution. Using 2 or more together causes a duplicate key error.
//
//	// specifically export x from first
//	export {x} from 'first';
//	// specifically export x from second
//	export {x} from 'second';
//	// export our own x instead
//	export let x;
//
// Note: This is purposefully a warning. The spec does not treat this as an
// error. It is only an error if an import attempts to use an ambiguous
// name. But just having ambiguous names is not itself an error.
var EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION = DiagnosticType_Warning(
	"JSC_AMBIGUOUS_EXPORT_DEFINITION", "The export \"{0}\" is ambiguous.")

// Java: EsModuleProcessor.CYCLIC_EXPORT_DEFINITION
var EsModuleProcessor_CYCLIC_EXPORT_DEFINITION = DiagnosticType_Error(
	"JSC_CYCLIC_EXPORT_DEFINITION", "Cyclic export detected while resolving name \"{0}\".")

// EsModuleProcessor_DUPLICATE_EXPORT. Note: We only check for duplicate
// exports, not imports. Imports cannot be shadowed (by any other binding,
// including other imports) and so we check that in VariableReferenceCheck.
var EsModuleProcessor_DUPLICATE_EXPORT = DiagnosticType_Error("JSC_DUPLICATE_EXPORT", "Duplicate export of \"{0}\".")

// Java: EsModuleProcessor.IMPORTED_AMBIGUOUS_EXPORT
var EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT = DiagnosticType_Error(
	"JSC_IMPORTED_AMBIGUOUS_EXPORT", "The requested name \"{0}\" is ambiguous.")

// Java: EsModuleProcessor.NAMESPACE_IMPORT_CANNOT_USE_STAR
var EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR = DiagnosticType_Error(
	"JSC_NAMESPACE_IMPORT_CANNOT_USE_STAR",
	"Namespace imports ('goog:some.Namespace') cannot use import * as. "+
		"Did you mean to import {0} from ''{1}'';?")

// Java: EsModuleProcessor.CANNOT_PATH_IMPORT_CLOSURE_FILE
var EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE = DiagnosticType_Error(
	"JSC_CANNOT_PATH_IMPORT_CLOSURE_FILE",
	"Cannot import Closure files by path. Use either import 'goog:namespace' or"+
		" goog.require('namespace')")

// EsModuleProcessor_FindMutableExports marks all exports that are mutated in
// an inner scope as mutable.
//
// Exports mutated at the module scope are not marked as mutable as they are
// effectively constant after module evaluation.
type EsModuleProcessor_FindMutableExports struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	exports            *jutil.ArrayList[*Export]
	exportsByLocalName *jutil.ListMultimap[string, *Export]
}

// Java: EsModuleProcessor.FindMutableExports(List<Export>)
func NewEsModuleProcessor_FindMutableExports(exports *jutil.ArrayList[*Export]) *EsModuleProcessor_FindMutableExports {
	f := &EsModuleProcessor_FindMutableExports{}
	// There may be multiple exports with the same local name because you can
	// export a local variable with an alias many different times. Example:
	//
	// let x;
	// export {x as y, x as z};
	f.exports = exports
	f.exportsByLocalName = jutil.ArrayListMultimap_Create[string, *Export]()
	for e := range exports.All() {
		f.exportsByLocalName.Put(esModuleProcessor_localNameKey(e), e)
	}
	return f
}

// esModuleProcessor_localNameKey is e.localName() as a multimap key.
//
// Go-only: language-adapter — Export.localName() is a nullable String; the
// local exports this multimap holds always have one.
func esModuleProcessor_localNameKey(e *Export) string {
	return *e.LocalName()
}

// Java: EsModuleProcessor.FindMutableExports.visit(NodeTraversal, Node, Node)
func (f *EsModuleProcessor_FindMutableExports) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsName() {
		return
	}

	scope := t.GetScope()
	if !NodeUtil_IsLValue(n) || scope.GetClosestHoistScope().IsModuleScope() {
		return
	}

	exports := f.exportsByLocalName.Get(n.GetString())
	if len(exports) == 0 {
		return
	}

	var_ := scope.GetVar(n.GetString())
	// A var declared in the module scope with the same name as an export must
	// be the export. And we know we're setting it in a function scope, so
	// this cannot be the declaration itself. We must be mutating.
	if var_ == nil || !var_.GetScope().(*Scope).IsModuleScope() {
		return
	}

	for _, e := range append([]*Export(nil), exports...) {
		i := esModuleProcessor_indexOf(f.exports, e)
		mutated := e.MutatedCopy()
		f.exports.Set(i, mutated)
		f.exportsByLocalName.Remove(esModuleProcessor_localNameKey(e), e)
		f.exportsByLocalName.Put(esModuleProcessor_localNameKey(e), mutated)
	}
}

// esModuleProcessor_indexOf is List.indexOf(Object) with Export.equals.
//
// Go-only: language-adapter — jutil.ArrayList.IndexOf compares pointers.
func esModuleProcessor_indexOf(list *jutil.ArrayList[*Export], e *Export) int {
	for i, x := range list.Slice() {
		if e.Equals(x) {
			return i
		}
	}
	return -1
}

// EsModuleProcessor_UnresolvedModuleBuilder collects all imports and exports
// of a module. Builds a UnresolvedEsModule by slotting the exports into one
// of three categories:
//
//   - Local exports - exports whose original definition is in this file.
//   - Indirect exports - exports that come from another file. Either via
//     `export {name} from` or `import {name}; export {name};`
//   - Star exports - `export * from` exports
//
// It is an inner class: outer is the enclosing EsModuleProcessor.
type EsModuleProcessor_UnresolvedModuleBuilder struct {
	outer              *EsModuleProcessor
	path               *ModuleLoader_ModulePath
	root               *rhino.Node
	importsByLocalName *jutil.LinkedHashMap[string, *Import]
	exports            *jutil.ArrayList[*Export]
	exportedNames      *jutil.LinkedHashSet[string]
}

// Java: EsModuleProcessor.UnresolvedModuleBuilder(ModulePath, Node)
func (p *EsModuleProcessor) NewUnresolvedModuleBuilder(path *ModuleLoader_ModulePath, root *rhino.Node) *EsModuleProcessor_UnresolvedModuleBuilder {
	b := &EsModuleProcessor_UnresolvedModuleBuilder{outer: p}
	b.path = path
	b.root = root
	b.importsByLocalName = jutil.NewLinkedHashMap[string, *Import]()
	b.exports = jutil.NewArrayList[*Export]()
	b.exportedNames = jutil.NewLinkedHashSet[string]()
	return b
}

// Java: EsModuleProcessor.UnresolvedModuleBuilder.add(Import)
func (b *EsModuleProcessor_UnresolvedModuleBuilder) Add(i *Import) {
	b.importsByLocalName.Put(i.LocalName(), i)
}

// AddExport is true if the export was successfully added or false if the
// exported name already exists.
//
// Java: EsModuleProcessor.UnresolvedModuleBuilder.add(Export)
func (b *EsModuleProcessor_UnresolvedModuleBuilder) AddExport(e *Export) bool {
	if e.ExportName() == nil {
		b.exports.Add(e)
		return true
	} else {
		if b.exportedNames.Contains(*e.ExportName()) {
			return false
		}

		b.exports.Add(e)
		b.exportedNames.Add(*e.ExportName())
		return true
	}
}

// Java: EsModuleProcessor.UnresolvedModuleBuilder.build()
func (b *EsModuleProcessor_UnresolvedModuleBuilder) Build() UnresolvedModule {
	localExports := jutil.NewArrayList[*Export]()
	var indirectExportsBuilder []*Export
	var starExportsBuilder []*Export

	for ee := range b.exports.All() {
		if ee.ModuleRequest() == nil {
			if ee.LocalName() != nil && b.importsByLocalName.ContainsKey(*ee.LocalName()) {
				// import A from '';
				// export { A };

				// import * as ns from '';
				// export { ns };
				indirectExportsBuilder = append(indirectExportsBuilder, ee)
			} else {
				// var y;
				// export { y };
				// export var x;
				localExports.Add(ee)
			}
		} else if ee.ImportName() != nil && "*" == *ee.ImportName() {
			// export * from '';
			starExportsBuilder = append(starExportsBuilder, ee)
		} else {
			// export { A } from '';
			indirectExportsBuilder = append(indirectExportsBuilder, ee)
		}
	}

	NodeTraversal_Traverse(b.outer.compiler, b.root, NewEsModuleProcessor_FindMutableExports(localExports))

	return b.outer.NewUnresolvedEsModule(
		b.outer.metadata,
		b.path,
		jutil.NewLinkedHashMapFrom(b.importsByLocalName),
		localExports.ToSlice(),
		indirectExportsBuilder,
		starExportsBuilder)
}

// EsModuleProcessor_UnresolvedEsModule is a module that has had its exports
// and imports parsed and categorized but has yet to transitively resolve
// them.
//
// It is an inner class: outer is the enclosing EsModuleProcessor.
type EsModuleProcessor_UnresolvedEsModule struct {
	UnresolvedModuleBase

	outer              *EsModuleProcessor
	metadata           *ModuleMetadataMap_ModuleMetadata
	path               *ModuleLoader_ModulePath
	importsByLocalName *jutil.LinkedHashMap[string, *Import]
	localExports       []*Export
	indirectExports    []*Export
	starExports        []*Export
	exportedNames      *jutil.LinkedHashSet[string]
	resolvedImports    *jutil.LinkedHashMap[string, *ResolveExportResult]
	resolvedExports    *jutil.LinkedHashMap[string, *ResolveExportResult]
	resolved           *Module
}

// Java: EsModuleProcessor.UnresolvedEsModule(ModuleMetadata, ModulePath, ImmutableMap<String, Import>, ImmutableList<Export>, ImmutableList<Export>, ImmutableList<Export>)
func (p *EsModuleProcessor) NewUnresolvedEsModule(
	metadata *ModuleMetadataMap_ModuleMetadata,
	path *ModuleLoader_ModulePath,
	importsByLocalName *jutil.LinkedHashMap[string, *Import],
	localExports []*Export,
	indirectExports []*Export,
	starExports []*Export) *EsModuleProcessor_UnresolvedEsModule {
	m := &EsModuleProcessor_UnresolvedEsModule{outer: p}
	m.InitUnresolvedModuleBase(m)
	m.metadata = metadata
	m.path = path
	m.importsByLocalName = importsByLocalName
	m.localExports = localExports
	m.indirectExports = indirectExports
	m.starExports = starExports
	m.exportedNames = nil
	m.resolvedImports = jutil.NewLinkedHashMap[string, *ResolveExportResult]()
	m.resolvedExports = jutil.NewLinkedHashMap[string, *ResolveExportResult]()
	m.resolved = nil
	return m
}

// Java: EsModuleProcessor.UnresolvedEsModule.reset()
func (m *EsModuleProcessor_UnresolvedEsModule) Reset() {
	m.resolved = nil
	m.exportedNames = nil
	m.resolvedImports.Clear()
	m.resolvedExports.Clear()
}

// Java: EsModuleProcessor.UnresolvedEsModule.resolve(ModuleRequestResolver, String)
func (m *EsModuleProcessor_UnresolvedEsModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module {
	if m.resolved == nil {
		// Every import creates a locally bound name.
		boundNames := jutil.NewLinkedHashMapFrom(m.GetAllResolvedImports(moduleRequestResolver))

		localNameToLocalExport := jutil.NewLinkedHashMap[string, *Export]()

		// Only local exports that are not an anonymous default export create
		// local bindings.
		for _, e := range m.localExports {
			localNameToLocalExport.Put(*e.LocalName(), e)

			if Export_DEFAULT_EXPORT_NAME != *e.LocalName() {
				b := m.ResolveExport(moduleRequestResolver, *e.ExportName())
				jutil.CheckStateMsg(b.Resolved(), "Cannot have invalid missing own export!")
				if !b.IsAmbiguous() {
					boundNames.Put(*e.LocalName(), b.GetBinding())
				}
			}
		}

		// getAllResolvedExports is required to make a module but also
		// performs cycle checking.
		m.resolved =
			Module_Builder_().
				BoundNames(closureModuleProcessor_immutableMapCopyOf(boundNames)).
				Namespace(closureModuleProcessor_immutableMapCopyOf(m.GetAllResolvedExports(moduleRequestResolver))).
				Metadata(m.metadata).
				Path(m.path).
				LocalNameToLocalExport(closureModuleProcessor_immutableMapCopyOf(localNameToLocalExport)).
				Build()
	}

	return m.resolved
}

// Java: EsModuleProcessor.UnresolvedEsModule.metadata()
func (m *EsModuleProcessor_UnresolvedEsModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

// GetAllResolvedImports is a map from import bound name to binding.
//
// Java: EsModuleProcessor.UnresolvedEsModule.getAllResolvedImports(ModuleRequestResolver)
func (m *EsModuleProcessor_UnresolvedEsModule) GetAllResolvedImports(
	moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashMap[string, *Binding] {
	imports := jutil.NewLinkedHashMap[string, *Binding]()

	for name := range m.importsByLocalName.Keys() {
		b := m.ResolveImport(moduleRequestResolver, name)
		if b.Resolved() {
			imports.Put(name, b.GetBinding())
		}
	}

	return imports
}

// Java: EsModuleProcessor.UnresolvedEsModule.resolveImport(ModuleRequestResolver, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *EsModuleProcessor_UnresolvedEsModule) ResolveImportWithSets(
	moduleRequestResolver ModuleRequestResolver,
	name string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	if m.resolvedImports.ContainsKey(name) {
		return m.resolvedImports.Get(name)
	}
	b := m.ResolveImportImpl(moduleRequestResolver, name, resolveSet, exportStarSet)
	m.resolvedImports.Put(name, b)
	return b
}

// Java: EsModuleProcessor.UnresolvedEsModule.resolveImport(ModuleRequestResolver, String)
func (m *EsModuleProcessor_UnresolvedEsModule) ResolveImport(
	moduleRequestResolver ModuleRequestResolver, name string) *ResolveExportResult {
	return m.ResolveImportWithSets(
		moduleRequestResolver, name, jutil.NewLinkedHashSet[ExportTrace](), jutil.NewLinkedHashSet[UnresolvedModule]())
}

// Java: EsModuleProcessor.UnresolvedEsModule.resolveImportImpl(ModuleRequestResolver, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *EsModuleProcessor_UnresolvedEsModule) ResolveImportImpl(
	moduleRequestResolver ModuleRequestResolver,
	name string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	i := m.importsByLocalName.Get(name)

	requested := moduleRequestResolver.Resolve(i)
	if requested == nil {
		return ResolveExportResult_ERROR
	} else {
		importStar := i.ImportName() == "*"
		if importStar ||
			(i.ImportName() == Export_DEFAULT &&
				(requested.Metadata().IsGoogProvide() || requested.Metadata().IsGoogModule())) {
			if !GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest()) &&
				(requested.Metadata().IsGoogModule() || requested.Metadata().IsGoogProvide()) {
				m.outer.compiler.Report(
					JSError_MakeWithSourceLocation(
						jutil.StringPtr(m.path.String()),
						i.ImportNode().GetLineno(),
						i.ImportNode().GetCharno(),
						EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE,
						i.LocalName(),
						i.ModuleRequest()))
				return ResolveExportResult_ERROR
			}
			if importStar && GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest()) {
				m.outer.compiler.Report(
					JSError_MakeWithSourceLocation(
						jutil.StringPtr(m.path.String()),
						i.ImportNode().GetLineno(),
						i.ImportNode().GetCharno(),
						EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR,
						i.LocalName(),
						i.ModuleRequest()))
				return ResolveExportResult_ERROR
			}
			var closureNamespace *string
			if GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest()) {
				ns := GoogEsImports_GetClosureIdFromGoogImportSpecifier(i.ModuleRequest())
				closureNamespace = &ns
			}
			return ResolveExportResult_Of(
				Binding_FromImportStar(requested.Metadata(), closureNamespace, i.NameNode()))
		} else {
			moduleRequest := i.ModuleRequest()
			result :=
				requested.ResolveExportWithTraces(
					moduleRequestResolver,
					&moduleRequest,
					i.ImportName(),
					resolveSet,
					exportStarSet)
			if !result.Found() && !result.HadError() {
				m.outer.compiler.Report(
					JSError_MakeWithSourceLocation(
						jutil.StringPtr(m.path.String()),
						i.ImportNode().GetLineno(),
						i.ImportNode().GetCharno(),
						ModuleMapCreator_DOES_NOT_HAVE_EXPORT,
						i.ImportName()))
				return ResolveExportResult_ERROR
			} else if result.IsAmbiguous() {
				m.outer.compiler.Report(
					JSError_MakeWithSourceLocation(
						jutil.StringPtr(m.path.String()),
						i.ImportNode().GetLineno(),
						i.ImportNode().GetCharno(),
						EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT,
						i.ImportName()))
				return ResolveExportResult_ERROR
			}
			forSourceInfo := i.NameNode()
			if forSourceInfo == nil {
				forSourceInfo = i.ImportNode()
			}
			return result.Copy(forSourceInfo, Binding_CreatedBy_IMPORT)
		}
	}
}

// Java: EsModuleProcessor.UnresolvedEsModule.getExportedNames(ModuleRequestResolver)
func (m *EsModuleProcessor_UnresolvedEsModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	if m.exportedNames == nil {
		m.exportedNames = m.GetExportedNamesWithVisited(moduleRequestResolver, jutil.NewLinkedHashSet[UnresolvedModule]())
	}
	return m.exportedNames
}

// Java: EsModuleProcessor.UnresolvedEsModule.getExportedNames(ModuleRequestResolver, Set<UnresolvedModule>)
func (m *EsModuleProcessor_UnresolvedEsModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	if visited.Contains(m) {
		// import * cycle
		return jutil.NewLinkedHashSet[string]()
	}

	visited.Add(m)
	// exports are in Array.prototype.sort() order.
	exportedNames := jutil.NewTreeSet[string](jutil.StringCompareTo)

	for _, e := range m.localExports {
		exportedNames.Add(*e.ExportName())
	}

	for _, e := range m.indirectExports {
		exportedNames.Add(*e.ExportName())
	}

	for _, e := range m.starExports {
		requested := moduleRequestResolver.ResolveWithExport(e)

		if requested != nil {
			if requested.Metadata().IsEs6Module() {
				for n := range requested.GetExportedNamesWithVisited(moduleRequestResolver, visited).All() {
					// Default exports are not exported with export *.
					if Export_DEFAULT != n && !exportedNames.Contains(n) {
						exportedNames.Add(n)
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

	return jutil.NewLinkedHashSet(exportedNames.ToSlice()...)
}

// GetAllResolvedExports is a map of exported name to binding.
//
// Java: EsModuleProcessor.UnresolvedEsModule.getAllResolvedExports(ModuleRequestResolver)
func (m *EsModuleProcessor_UnresolvedEsModule) GetAllResolvedExports(
	moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashMap[string, *Binding] {
	exports := jutil.NewLinkedHashMap[string, *Binding]()

	for name := range m.GetExportedNames(moduleRequestResolver).All() {
		b := m.ResolveExport(moduleRequestResolver, name)
		jutil.CheckStateMsg(b.Found(), "Cannot have invalid own export.")
		if b.Resolved() {
			exports.Put(name, b.GetBinding())
		} else if b.IsAmbiguous() {
			m.outer.compiler.Report(JSError_MakeWithSourceLocation(
				jutil.StringPtr(m.path.String()), -1, -1, EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION, name))
		}
	}

	return exports
}

// Java: EsModuleProcessor.UnresolvedEsModule.resolveExport(ModuleRequestResolver, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *EsModuleProcessor_UnresolvedEsModule) ResolveExportWithSets(
	moduleRequestResolver ModuleRequestResolver,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	if !m.GetExportedNames(moduleRequestResolver).Contains(exportName) {
		return ResolveExportResult_NOT_FOUND
	}

	if !resolveSet.Add(ExportTrace_Create(m, exportName)) {
		// Cycle!
		m.outer.compiler.Report(JSError_MakeWithSourceLocation(
			jutil.StringPtr(m.path.String()), 0, 0, EsModuleProcessor_CYCLIC_EXPORT_DEFINITION, exportName))
		return ResolveExportResult_ERROR
	}

	for _, e := range m.localExports {
		if e.ExportName() != nil && exportName == *e.ExportName() {
			forSourceInfo := e.ExportNode()
			if e.NameNode() != nil {
				forSourceInfo = e.NameNode()
			}
			return ResolveExportResult_Of(Binding_From(e, forSourceInfo))
		}
	}

	for _, e := range m.indirectExports {
		if e.ExportName() != nil && exportName == *e.ExportName() {
			if e.LocalName() != nil && m.importsByLocalName.ContainsKey(*e.LocalName()) {
				// import whatever from 'mod';
				// export { whatever };
				return m.ResolveImportWithSets(moduleRequestResolver, *e.LocalName(), resolveSet, exportStarSet).
					Copy(e.NameNode(), Binding_CreatedBy_EXPORT)
			} else {
				requested := moduleRequestResolver.ResolveWithExport(e)

				if requested == nil {
					return ResolveExportResult_ERROR
				} else {
					// export { whatever } from 'mod';
					result :=
						requested.ResolveExportWithTraces(
							moduleRequestResolver,
							e.ModuleRequest(),
							*e.ImportName(),
							resolveSet,
							exportStarSet)
					if !result.Found() && !result.HadError() {
						m.outer.compiler.Report(
							JSError_MakeWithSourceLocation(
								jutil.StringPtr(m.path.String()),
								e.ExportNode().GetLineno(),
								e.ExportNode().GetCharno(),
								ModuleMapCreator_DOES_NOT_HAVE_EXPORT,
								*e.ImportName()))
						return ResolveExportResult_ERROR
					} else if result.IsAmbiguous() {
						m.outer.compiler.Report(
							JSError_MakeWithSourceLocation(
								jutil.StringPtr(m.path.String()),
								e.ExportNode().GetLineno(),
								e.ExportNode().GetCharno(),
								EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT,
								*e.ImportName()))
					}
					return result.Copy(e.NameNode(), Binding_CreatedBy_EXPORT)
				}
			}
		}
	}

	jutil.CheckStateMsg(Export_DEFAULT != exportName, "Default export cannot come from export *.")

	if exportStarSet.Contains(m) {
		// Cycle!
		m.outer.compiler.Report(JSError_MakeWithSourceLocation(
			jutil.StringPtr(m.path.String()), -1, -1, EsModuleProcessor_CYCLIC_EXPORT_DEFINITION, exportName))
		return ResolveExportResult_ERROR
	}

	exportStarSet.Add(m)

	var starResolution *ResolveExportResult

	for _, e := range m.starExports {
		requested := moduleRequestResolver.ResolveWithExport(e)

		if requested == nil {
			return ResolveExportResult_ERROR
		} else if requested.GetExportedNames(moduleRequestResolver).Contains(exportName) {
			resolution :=
				requested.ResolveExportWithTraces(
					moduleRequestResolver, e.ModuleRequest(), exportName, resolveSet, exportStarSet)

			if resolution.HadError() {
				// Recursive case; error was reported on base case.
				return resolution
			} else if resolution.IsAmbiguous() {
				return resolution
			} else {
				if starResolution == nil {
					// First time finding something, not ambiguous.
					starResolution = resolution.Copy(e.ExportNode(), Binding_CreatedBy_EXPORT)
				} else {
					// Second time finding something, might be ambiguous!
					// Not ambiguous if it is the same export (same module and
					// export name).
					if starResolution != resolution {
						return ResolveExportResult_AMBIGUOUS
					}
				}
			}
		}
	}

	if starResolution == nil {
		return ResolveExportResult_ERROR
	}

	return starResolution
}

// Java: EsModuleProcessor.UnresolvedEsModule.resolveExport(ModuleRequestResolver, String, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *EsModuleProcessor_UnresolvedEsModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier *string,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	// Explicit containsKey check here since values can be null!
	if m.resolvedExports.ContainsKey(exportName) {
		return m.resolvedExports.Get(exportName)
	}

	b := m.ResolveExportWithSets(moduleRequestResolver, exportName, resolveSet, exportStarSet)
	m.resolvedExports.Put(exportName, b)
	return b
}

// Java: EsModuleProcessor(AbstractCompiler)
func NewEsModuleProcessor(compiler AbstractCompiler) *EsModuleProcessor {
	return &EsModuleProcessor{compiler: compiler}
}

// Java: EsModuleProcessor.process(ModuleMetadata, ModulePath, Node)
func (p *EsModuleProcessor) Process(
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script *rhino.Node) UnresolvedModule {
	p.metadata = metadata
	p.currentModuleBuilder = p.NewUnresolvedModuleBuilder(path, script)

	NodeTraversal_Traverse(p.compiler, script, p)

	m := p.currentModuleBuilder.Build()
	p.metadata = nil
	p.currentModuleBuilder = nil
	return m
}

// Java: EsModuleProcessor.shouldTraverse(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT, rhino.Token_IMPORT:
		return true
	default:
		return false
	}
}

// Java: EsModuleProcessor.visit(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_EXPORT:
		p.VisitExport(t, n)
	case rhino.Token_IMPORT:
		p.VisitImport(t, n)
	default:
	}
}

// Java: EsModuleProcessor.visitExportAllFrom(NodeTraversal, Node)
func (p *EsModuleProcessor) VisitExportAllFrom(t *NodeTraversal, export *rhino.Node) {
	// export * from '';
	moduleRequest := export.GetSecondChild().GetString()
	importName := "*"
	p.currentModuleBuilder.AddExport(
		Export_Builder_().
			ExportName(nil).
			ModuleRequest(&moduleRequest).
			ImportName(&importName).
			LocalName(nil).
			ModulePath(t.GetInput().GetPath()).
			ExportNode(export).
			ModuleMetadata(p.metadata).
			Build())
}

// Java: EsModuleProcessor.visitExportDefault(NodeTraversal, Node)
func (p *EsModuleProcessor) VisitExportDefault(t *NodeTraversal, export *rhino.Node) {
	// export default <expression>;
	child := export.GetFirstChild()
	name := Export_DEFAULT_EXPORT_NAME

	if child.IsFunction() || child.IsClass() {
		maybeName := NodeUtil_GetName(child)

		if maybeName != nil && *maybeName != "" {
			name = *maybeName
		}
	}

	exportName := Export_DEFAULT
	if !p.currentModuleBuilder.AddExport(
		Export_Builder_().
			ExportName(&exportName).
			ModuleRequest(nil).
			ImportName(nil).
			LocalName(&name).
			ModulePath(t.GetInput().GetPath()).
			ExportNode(export).
			ModuleMetadata(p.metadata).
			Build()) {
		t.Report(export, EsModuleProcessor_DUPLICATE_EXPORT, Export_DEFAULT)
	}
}

// Java: EsModuleProcessor.visitExportFrom(NodeTraversal, Node)
func (p *EsModuleProcessor) VisitExportFrom(t *NodeTraversal, export *rhino.Node) {
	// export { Foo, Bar as Rab } from '';
	for child := export.GetFirstFirstChild(); child != nil; child = child.GetNext() {
		importName := child.GetFirstChild().GetString()
		exportedName := child.GetLastChild().GetString()
		moduleRequest := export.GetSecondChild().GetString()
		if !p.currentModuleBuilder.AddExport(
			Export_Builder_().
				ExportName(&exportedName).
				ModuleRequest(&moduleRequest).
				ImportName(&importName).
				LocalName(nil).
				ModulePath(t.GetInput().GetPath()).
				ExportNode(export).
				NameNode(child.GetFirstChild()).
				ModuleMetadata(p.metadata).
				Build()) {
			t.Report(export, EsModuleProcessor_DUPLICATE_EXPORT, exportedName)
		}
	}
}

// Java: EsModuleProcessor.visitExportSpecs(NodeTraversal, Node)
func (p *EsModuleProcessor) VisitExportSpecs(t *NodeTraversal, export *rhino.Node) {
	// export { Foo, Bar as Rab };
	for child := export.GetFirstFirstChild(); child != nil; child = child.GetNext() {
		localName := child.GetFirstChild().GetString()
		exportedName := child.GetLastChild().GetString()
		if !p.currentModuleBuilder.AddExport(
			Export_Builder_().
				ExportName(&exportedName).
				ModuleRequest(nil).
				ImportName(nil).
				LocalName(&localName).
				ModulePath(t.GetInput().GetPath()).
				ExportNode(export).
				NameNode(child.GetFirstChild()).
				ModuleMetadata(p.metadata).
				Build()) {
			t.Report(export, EsModuleProcessor_DUPLICATE_EXPORT, exportedName)
		}
	}
}

// Java: EsModuleProcessor.visitExportNameDeclaration(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) VisitExportNameDeclaration(t *NodeTraversal, export *rhino.Node, declaration *rhino.Node) {
	//    export var Foo;
	//    export let {a, b:[c,d]} = {};
	NodeUtil_VisitLhsNodesInNode(
		declaration, func(lhs *rhino.Node) { p.AddExportNameDeclaration(t, lhs, export) })
}

// Java: EsModuleProcessor.addExportNameDeclaration(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) AddExportNameDeclaration(t *NodeTraversal, lhs *rhino.Node, export *rhino.Node) {
	jutil.CheckState(lhs.IsName())
	name := lhs.GetString()
	if !p.currentModuleBuilder.AddExport(
		Export_Builder_().
			ExportName(&name).
			ModuleRequest(nil).
			ImportName(nil).
			LocalName(&name).
			ModulePath(t.GetInput().GetPath()).
			ExportNode(export).
			NameNode(lhs).
			ModuleMetadata(p.metadata).
			Build()) {
		t.Report(export, EsModuleProcessor_DUPLICATE_EXPORT, name)
	}
}

// Java: EsModuleProcessor.visitExportFunctionOrClass(NodeTraversal, Node, Node)
func (p *EsModuleProcessor) VisitExportFunctionOrClass(t *NodeTraversal, export *rhino.Node, declaration *rhino.Node) {
	// export function foo() {}
	// export class Foo {}
	jutil.CheckState(declaration.IsFunction() || declaration.IsClass())
	nameNode := declaration.GetFirstChild()
	name := nameNode.GetString()
	if !p.currentModuleBuilder.AddExport(
		Export_Builder_().
			ExportName(&name).
			ModuleRequest(nil).
			ImportName(nil).
			LocalName(&name).
			ModulePath(t.GetInput().GetPath()).
			ExportNode(export).
			NameNode(nameNode).
			ModuleMetadata(p.metadata).
			Build()) {
		t.Report(export, EsModuleProcessor_DUPLICATE_EXPORT, name)
	}
}

// Java: EsModuleProcessor.visitExport(NodeTraversal, Node)
func (p *EsModuleProcessor) VisitExport(t *NodeTraversal, export *rhino.Node) {
	if export.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) {
		p.VisitExportAllFrom(t, export)
	} else if export.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
		p.VisitExportDefault(t, export)
	} else if export.HasTwoChildren() {
		p.VisitExportFrom(t, export)
	} else if export.GetFirstChild().IsExportSpecs() {
		p.VisitExportSpecs(t, export)
	} else {
		declaration := export.GetFirstChild()
		if NodeUtil_IsNameDeclaration(declaration) {
			p.VisitExportNameDeclaration(t, export, declaration)
		} else {
			p.VisitExportFunctionOrClass(t, export, declaration)
		}
	}
}

// Java: EsModuleProcessor.visitImportDefault(NodeTraversal, Node, String)
func (p *EsModuleProcessor) VisitImportDefault(t *NodeTraversal, importNode *rhino.Node, moduleRequest string) {
	// import Default from '';
	p.currentModuleBuilder.Add(
		Import_Builder_().
			ModuleRequest(moduleRequest).
			ImportName(Export_DEFAULT).
			LocalName(importNode.GetFirstChild().GetString()).
			ModulePath(t.GetInput().GetPath()).
			ImportNode(importNode).
			NameNode(importNode.GetFirstChild()).
			Build())
}

// Java: EsModuleProcessor.visitImportSpecs(NodeTraversal, Node, String)
func (p *EsModuleProcessor) VisitImportSpecs(t *NodeTraversal, importNode *rhino.Node, moduleRequest string) {
	// import { A, b as B } from '';
	for child := importNode.GetSecondChild().GetFirstChild(); child != nil; child = child.GetNext() {
		importName := child.GetFirstChild().GetString()
		localName := child.GetLastChild().GetString()
		p.currentModuleBuilder.Add(
			Import_Builder_().
				ModuleRequest(moduleRequest).
				ImportName(importName).
				LocalName(localName).
				ModulePath(t.GetInput().GetPath()).
				ImportNode(importNode).
				NameNode(child.GetSecondChild()).
				Build())
	}
}

// Java: EsModuleProcessor.visitImportStar(NodeTraversal, Node, String)
func (p *EsModuleProcessor) VisitImportStar(t *NodeTraversal, importNode *rhino.Node, moduleRequest string) {
	// import * as ns from '';
	p.currentModuleBuilder.Add(
		Import_Builder_().
			ModuleRequest(moduleRequest).
			ImportName("*").
			LocalName(importNode.GetSecondChild().GetString()).
			ImportNode(importNode).
			ModulePath(t.GetInput().GetPath()).
			NameNode(importNode.GetSecondChild()).
			Build())
}

// Java: EsModuleProcessor.visitImport(NodeTraversal, Node)
func (p *EsModuleProcessor) VisitImport(t *NodeTraversal, importNode *rhino.Node) {
	moduleRequest := importNode.GetLastChild().GetString()

	if importNode.GetFirstChild().IsName() {
		p.VisitImportDefault(t, importNode, moduleRequest)
	}

	if importNode.GetSecondChild().IsImportSpecs() {
		p.VisitImportSpecs(t, importNode, moduleRequest)
	} else if importNode.GetSecondChild().IsImportStar() {
		p.VisitImportStar(t, importNode, moduleRequest)
	}

	// no entry for import '';
}
