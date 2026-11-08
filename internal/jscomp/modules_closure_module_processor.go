// Port of com/google/javascript/jscomp/modules/ClosureModuleProcessor.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// ClosureModuleProcessor is the processor for goog.module.
type ClosureModuleProcessor struct {
	compiler AbstractCompiler
}

var _ moduleProcessor = (*ClosureModuleProcessor)(nil)

var closureModuleExportsID = uint32(strpool.Intern("exports"))

type unresolvedGoogModule struct {
	metadata            *ModuleMetadataMap_ModuleMetadata
	srcFileName         base.OptString
	path                *ModuleLoader_ModulePath
	namespace           *OrderedStringMap[*Binding]
	requireNames        []string
	requiresByLocalName map[string]*ClosureRequireProcessor_Require
	compiler            AbstractCompiler
	resolved            *Module
}

var _ UnresolvedModule = (*unresolvedGoogModule)(nil)

func (m *unresolvedGoogModule) Resolve(moduleRequestResolver ModuleRequestResolver) *Module {
	return m.ResolveWithModuleSpecifier(moduleRequestResolver, "", false)
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
	moduleSpecifier string,
	hasModuleSpecifier bool,
	exportName string,
	resolveSet map[ExportTrace]struct{},
	exportStarSet map[UnresolvedModule]struct{}) *ResolveExportResult {
	return m.ResolveExport(moduleRequestResolver, exportName)
}

func onlyClosureNamespace(metadata *ModuleMetadataMap_ModuleMetadata) string {
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
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, hasModuleSpecifier bool) *Module {
	if m.resolved == nil {
		boundNames := m.getAllResolvedImports(moduleRequestResolver)
		closureNamespace := onlyClosureNamespace(m.metadata)
		m.resolved = NewModule(
			m.metadata,
			m.path,
			m.namespace,
			boundNames,
			nonEsEmptyExportMap,
			closureNamespace,
			true,
		)
	}
	return m.resolved
}

func (m *unresolvedGoogModule) getAllResolvedImports(
	moduleRequestResolver ModuleRequestResolver) *OrderedStringMap[*Binding] {
	if len(m.requireNames) == 0 {
		return nonEsEmptyBindingMap
	}
	var imports *OrderedStringMap[*Binding]
	for _, name := range m.requireNames {
		b := m.resolveImport(moduleRequestResolver, name)
		if b.Resolved() {
			if imports == nil {
				imports = NewOrderedStringMapWithCapacity[*Binding](len(m.requireNames))
			}
			imports.Put(name, jutil.CheckNotNull(b.GetBinding()))
		}
	}
	if imports == nil {
		return nonEsEmptyBindingMap
	}
	return imports
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
				moduleRequest,
				require.CreatedBy()))
	}

	result := requested.ResolveExportWithTraces(
		moduleRequestResolver,
		moduleRequest,
		true,
		importRecord.ImportName(),
		nil,
		nil)
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

func (m *unresolvedGoogModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) []string {
	panic(jutil.NewUnsupportedOperationException())
}

func (m *unresolvedGoogModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited map[UnresolvedModule]struct{}) []string {
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

			if mayBeAccidentalDefaultExport(importRecord.ImportName(), defaultBinding) {
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
		JSError_MakeWithOptSourceLocation(
			m.srcFileName,
			importRecord.ImportNode().GetLineno(),
			importRecord.ImportNode().GetCharno(),
			ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS,
			importRecord.ImportName(),
			additionalInfo))
}

func mayBeAccidentalDefaultExport(importName string, defaultBinding *Binding) bool {
	defaultExport := defaultBinding.OriginatingExport().ExportNode()
	assign, isAssign := ast.AsAssign(defaultExport.GetParent())
	jutil.CheckStateMsg1(isClosureExportsIdent(defaultExport) && isAssign, "%s", defaultExport)
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
		closureNamespace:    onlyClosureNamespace(metadata),
		namespace:           &OrderedStringMap[*Binding]{},
		requiresByLocalName: make(map[string]*ClosureRequireProcessor_Require),
	}
	NodeTraversal_Traverse(p.compiler, script, cb)
	namespace := cb.namespace
	if namespace.IsEmpty() {
		namespace = nonEsEmptyBindingMap
	}

	var srcFileName base.OptString
	if src, ok := script.SourceFileName(); ok {
		srcFileName = base.SomeString(src)
	}
	return &unresolvedGoogModule{
		metadata:            metadata,
		srcFileName:         srcFileName,
		path:                path,
		namespace:           namespace,
		requireNames:        cb.requireNames,
		requiresByLocalName: cb.requiresByLocalName,
		compiler:            p.compiler,
	}
}

type closureModuleProcessingCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	metadata              *ModuleMetadataMap_ModuleMetadata
	closureNamespace      string
	namespace             *OrderedStringMap[*Binding]
	requireNames          []string
	requiresByLocalName   map[string]*ClosureRequireProcessor_Require
	seenExportsAssignment bool
}

var (
	_ NodeTraversal_Callback = (*closureModuleProcessingCallback)(nil)
	_ nodeTraversal_noVisit  = (*closureModuleProcessingCallback)(nil)
)

func isClosureExportsIdent(n rhino.Node) bool {
	return n.IsName() && n.GetStringID() == closureModuleExportsID
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
	if isClosureExportsIdent(lhs.Node()) {
		rhs := assign.Value().Node()
		if NodeUtil_IsNamedExportsLiteral(rhs) {
			for member := range ast.MustObjectLit(rhs).Members() {
				c.addPropertyExport(member.MustStringKey().Key(), member.Node())
			}
		} else {
			c.seenExportsAssignment = true
			c.markExportsAssignmentInNamespace(lhs.Node())
		}
	} else if getProp, ok := lhs.AsGetProp(); ok {
		c.maybeInitializeExportsStub(getProp)
	}
}

func (c *closureModuleProcessingCallback) maybeInitializeExportsStub(getProp ast.GetProp) {
	if isClosureExportsIdent(getProp.Object().Node()) {
		c.addPropertyExport(getProp.Property(), getProp.Node())
	}
}

func (c *closureModuleProcessingCallback) markExportsAssignmentInNamespace(exportsNode rhino.Node) {
	c.namespace.Put(
		Export_NAMESPACE,
		Binding_From(
			NewExportBuilder().
				ExportName(Export_NAMESPACE).
				ExportNode(exportsNode).
				ModuleMetadata(c.metadata).
				ClosureNamespace(c.closureNamespace).
				ModulePath(c.metadata.Path()).
				Build(),
			exportsNode))
}

func (c *closureModuleProcessingCallback) addPropertyExport(exportedId string, propNode rhino.Node) {
	if c.seenExportsAssignment || c.namespace.ContainsKey(exportedId) {
		return
	}
	c.namespace.Put(
		exportedId,
		Binding_From(
			NewExportBuilder().
				ExportName(exportedId).
				ExportNode(propNode).
				ModuleMetadata(c.metadata).
				ClosureNamespace(c.closureNamespace).
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
