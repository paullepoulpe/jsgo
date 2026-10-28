// Port of com/google/javascript/jscomp/modules/ClosureModuleProcessor.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureModuleProcessor is the processor for goog.module.
type ClosureModuleProcessor struct {
	compiler AbstractCompiler
}

var _ moduleProcessor = (*ClosureModuleProcessor)(nil)

type unresolvedGoogModule struct {
	UnresolvedModuleBase

	metadata            *ModuleMetadataMap_ModuleMetadata
	srcFileName         *string
	path                *ModuleLoader_ModulePath
	namespace           *jutil.LinkedHashMap[string, *Binding]
	requireNames        []string
	requiresByLocalName map[string]*ClosureRequireProcessor_Require
	compiler            AbstractCompiler
	resolved            *Module
}

func (m *unresolvedGoogModule) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	if b, ok := m.namespace.Lookup(exportName); ok {
		return ResolveExportResult_Of(b)
	}
	return ResolveExportResult_NOT_FOUND
}

func (m *unresolvedGoogModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier *string,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	return m.ResolveExport(moduleRequestResolver, exportName)
}

func closureModuleProcessor_onlyNamespace(metadata *ModuleMetadataMap_ModuleMetadata) string {
	var out string
	found := false
	for ns := range metadata.GoogNamespaces().All() {
		if found {
			panic(jutil.NewIllegalArgumentExceptionMsg("expected one element"))
		}
		out = ns
		found = true
	}
	if !found {
		panic(jutil.NewNoSuchElementException())
	}
	return out
}

func (m *unresolvedGoogModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module {
	if m.resolved == nil {
		boundNames := m.getAllResolvedImports(moduleRequestResolver)
		closureNamespace := closureModuleProcessor_onlyNamespace(m.metadata)
		m.resolved = NewModule(
			m.metadata,
			m.path,
			m.namespace,
			boundNames,
			nonEsEmptyExportMap,
			&closureNamespace,
		)
	}
	return m.resolved
}

func (m *unresolvedGoogModule) getAllResolvedImports(
	moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashMap[string, *Binding] {
	imports := jutil.NewLinkedHashMapWithCapacity[string, *Binding](len(m.requireNames))
	for _, name := range m.requireNames {
		b := m.resolveImport(moduleRequestResolver, name)
		if b.Resolved() {
			imports.Put(name, jutil.CheckNotNull(b.GetBinding()))
		}
	}
	return imports
}

// closureModuleProcessor_immutableMapCopyOf is ImmutableMap.copyOf(Map),
// which rejects null values.
func closureModuleProcessor_immutableMapCopyOf[V any](m *jutil.LinkedHashMap[string, V]) *jutil.LinkedHashMap[string, V] {
	out := jutil.NewLinkedHashMapWithCapacity[string, V](m.Size())
	for k, v := range m.All() {
		out.Put(k, jutil.CheckNotNull(v))
	}
	return out
}

func (m *unresolvedGoogModule) resolveImport(
	moduleRequestResolver ModuleRequestResolver, name string) *ResolveExportResult {
	require := m.requiresByLocalName[name]
	importRecord := require.ImportRecord()
	requested := moduleRequestResolver.Resolve(importRecord)

	if requested == nil {
		return ResolveExportResult_ERROR
	}
	moduleRequest := importRecord.ModuleRequest()
	if importRecord.ImportName() == Export_NAMESPACE {
		return ResolveExportResult_Of(
			Binding_FromClosureImport(
				requested.Metadata(),
				importRecord.NameNode(),
				&moduleRequest,
				require.CreatedBy()))
	}

	result := requested.ResolveExportWithTraces(
		moduleRequestResolver,
		&moduleRequest,
		importRecord.ImportName(),
		jutil.NewLinkedHashSet[ExportTrace](),
		jutil.NewLinkedHashSet[UnresolvedModule]())
	if !result.Found() && !result.HadError() {
		m.reportInvalidDestructuringRequire(requested, importRecord)
		return ResolveExportResult_ERROR
	}
	forSourceInfo := importRecord.NameNode()
	if forSourceInfo.IsNil() {
		forSourceInfo = importRecord.ImportNode()
	}
	return result.Copy(forSourceInfo, require.CreatedBy())
}

func (m *unresolvedGoogModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

func (m *unresolvedGoogModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

func (m *unresolvedGoogModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

func (m *unresolvedGoogModule) Reset() {
	m.resolved = nil
}

func (m *unresolvedGoogModule) reportInvalidDestructuringRequire(
	requested UnresolvedModule, importRecord *Import) {
	additionalInfo := ""
	if otherGoog, ok := requested.(*unresolvedGoogModule); ok {
		if defaultBinding, hasDefault := otherGoog.namespace.Lookup(Export_NAMESPACE); hasDefault {
			additionalInfo =
				jutil.LenientFormat(
					"\n"+
						"The goog.module \"%s\" cannot be destructured as it contains a default"+
						" export, not named exports. See %s.",
					importRecord.ModuleRequest(),
					"https://github.com/google/closure-library/wiki/goog.module%3A-an-ES6-module-like-alternative-to-goog.provide#destructuring-imports")

			if closureModuleProcessor_mayBeAccidentalDefaultExport(importRecord.ImportName(), defaultBinding) {
				additionalInfo +=
					jutil.LenientFormat(
						"\n"+
							"Either use a non-destructuring require or rewrite the goog.module"+
							" \"%s\" to support destructuring requires. For example, consider"+
							" replacing\n"+
							"  exports = {%s: <value>[, ...]};\n"+
							"with individual named export assignments like\n"+
							"  exports.%s = <value>;\n",
						importRecord.ModuleRequest(),
						importRecord.ImportName(),
						importRecord.ImportName())
			}
		}
	}
	m.compiler.Report(
		JSError_MakeWithSourceLocation(
			m.srcFileName,
			importRecord.ImportNode().GetLineno(),
			importRecord.ImportNode().GetCharno(),
			ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS,
			importRecord.ImportName(),
			additionalInfo))
}

func closureModuleProcessor_mayBeAccidentalDefaultExport(importName string, defaultBinding *Binding) bool {
	defaultExport := defaultBinding.OriginatingExport().ExportNode()
	assign, isAssign := ast.AsAssign(defaultExport.GetParent())
	jutil.CheckStateMsg1(defaultExport.MatchesName("exports") && isAssign, "%s", defaultExport)
	objLit, ok := assign.Value().AsObjectLit()
	if !ok {
		return false
	}
	for member := range objLit.Members() {
		if key, ok := member.AsStringKey(); ok && key.Key() == importName {
			return true
		}
	}
	return false
}

// Java: ClosureModuleProcessor(AbstractCompiler)
func NewClosureModuleProcessor(compiler AbstractCompiler) *ClosureModuleProcessor {
	return &ClosureModuleProcessor{compiler: compiler}
}

// Java: ClosureModuleProcessor.process(ModuleMetadata, ModulePath, Node)
func (p *ClosureModuleProcessor) Process(
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script rhino.Node) UnresolvedModule {
	jutil.CheckArgumentMsg1(
		script.IsScript() || script.IsCall(), "Unexpected module root %s", script)
	jutil.CheckArgumentMsg(
		script.IsCall() || path != nil, "Non goog.loadModules must have a path")

	cb := &closureModuleProcessingCallback{
		metadata:            metadata,
		closureNamespace:    closureModuleProcessor_onlyNamespace(metadata),
		namespace:           jutil.NewLinkedHashMap[string, *Binding](),
		requiresByLocalName: make(map[string]*ClosureRequireProcessor_Require),
	}
	NodeTraversal_Traverse(p.compiler, script, cb)

	m := &unresolvedGoogModule{
		metadata:            metadata,
		srcFileName:         script.GetSourceFileName(),
		path:                path,
		namespace:           cb.namespace,
		requireNames:        cb.requireNames,
		requiresByLocalName: cb.requiresByLocalName,
		compiler:            p.compiler,
	}
	m.InitUnresolvedModuleBase(m)
	return m
}

type closureModuleProcessingCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	metadata              *ModuleMetadataMap_ModuleMetadata
	closureNamespace      string
	namespace             *jutil.LinkedHashMap[string, *Binding]
	requireNames          []string
	requiresByLocalName   map[string]*ClosureRequireProcessor_Require
	seenExportsAssignment bool
}

func (c *closureModuleProcessingCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY, rhino.Token_SCRIPT, rhino.Token_CALL, rhino.Token_BLOCK:
		return true
	case rhino.Token_FUNCTION:
		return parent.IsCall() && parent == c.metadata.RootNode()
	case rhino.Token_EXPR_RESULT:
		expr := ast.MustExprStmt(n).Expr()
		if assign, ok := expr.AsAssign(); ok {
			c.maybeInitializeExports(assign)
		} else if getProp, ok := expr.AsGetProp(); ok {
			c.maybeInitializeExportsStub(getProp)
		}
		return false
	case rhino.Token_CONST, rhino.Token_VAR, rhino.Token_LET:
		c.maybeInitializeRequire(n)
		return false
	default:
		return false
	}
}

func (c *closureModuleProcessingCallback) maybeInitializeExports(assign ast.Assign) {
	lhs := assign.Target()
	if lhsName, ok := lhs.AsName(); ok && lhsName.Ident() == "exports" {
		rhs := assign.Value().Node()
		if NodeUtil_IsNamedExportsLiteral(rhs) {
			for member := range ast.MustObjectLit(rhs).Members() {
				c.addPropertyExport(member.Node().GetString(), member.Node())
			}
		} else {
			c.seenExportsAssignment = true
			c.markExportsAssignmentInNamespace(lhs.Node())
		}
	} else if getProp, ok := lhs.AsGetProp(); ok {
		if owner, ok := getProp.Object().AsName(); ok && owner.Ident() == "exports" {
			c.addPropertyExport(getProp.Property(), getProp.Node())
		}
	}
}

func (c *closureModuleProcessingCallback) maybeInitializeExportsStub(getProp ast.GetProp) {
	if owner, ok := getProp.Object().AsName(); ok && owner.Ident() == "exports" {
		c.addPropertyExport(getProp.Property(), getProp.Node())
	}
}

func (c *closureModuleProcessingCallback) markExportsAssignmentInNamespace(exportsNode rhino.Node) {
	exportName := Export_NAMESPACE
	closureNamespace := c.closureNamespace
	c.namespace.Put(
		Export_NAMESPACE,
		Binding_From(
			Export_Builder_().
				ExportName(&exportName).
				ExportNode(exportsNode).
				ModuleMetadata(c.metadata).
				ClosureNamespace(&closureNamespace).
				ModulePath(c.metadata.Path()).
				Build(),
			exportsNode))
}

func (c *closureModuleProcessingCallback) addPropertyExport(exportedId string, propNode rhino.Node) {
	if c.seenExportsAssignment || c.namespace.ContainsKey(exportedId) {
		return
	}
	closureNamespace := c.closureNamespace
	c.namespace.Put(
		exportedId,
		Binding_From(
			Export_Builder_().
				ExportName(&exportedId).
				ExportNode(propNode).
				ModuleMetadata(c.metadata).
				ClosureNamespace(&closureNamespace).
				ModulePath(c.metadata.Path()).
				Build(),
			propNode))
}

func (c *closureModuleProcessingCallback) maybeInitializeRequire(nameDeclaration rhino.Node) {
	for _, require := range ClosureRequireProcessor_GetAllRequires(nameDeclaration) {
		localName := require.LocalName()
		if _, exists := c.requiresByLocalName[localName]; !exists {
			c.requireNames = append(c.requireNames, localName)
			c.requiresByLocalName[localName] = require
		}
	}
}
