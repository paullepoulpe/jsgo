// Port of com/google/javascript/jscomp/modules/ClosureModuleProcessor.java @ closure-compiler v20260415

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureModuleProcessor is the processor for goog.module.
//
// The namespace of a goog.module contains all named exports, e.g.
// `exports.x = 0`, and any 'default export's that assign directly to the
// `exports` object, e.g. `exports = class {}`).
//
// The bound names include any names imported through a
// goog.require(Type)/forwardDeclare.
type ClosureModuleProcessor struct {
	compiler AbstractCompiler
}

var _ ModuleMapCreator_ModuleProcessor = (*ClosureModuleProcessor)(nil)

// ClosureModuleProcessor_UnresolvedGoogModule is
// ClosureModuleProcessor.UnresolvedGoogModule.
type ClosureModuleProcessor_UnresolvedGoogModule struct {
	UnresolvedModuleBase

	metadata            *ModuleMetadataMap_ModuleMetadata
	srcFileName         *string
	path                *ModuleLoader_ModulePath
	namespace           *jutil.LinkedHashMap[string, *Binding]
	requiresByLocalName *jutil.LinkedHashMap[string, *ClosureRequireProcessor_Require]
	compiler            AbstractCompiler
	resolved            *Module
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule(ModuleMetadata, String, ModulePath, ImmutableMap<String, Binding>, ImmutableMap<String, Require>, AbstractCompiler)
func NewClosureModuleProcessor_UnresolvedGoogModule(
	metadata *ModuleMetadataMap_ModuleMetadata,
	srcFileName *string,
	path *ModuleLoader_ModulePath,
	namespace *jutil.LinkedHashMap[string, *Binding],
	requires *jutil.LinkedHashMap[string, *ClosureRequireProcessor_Require],
	compiler AbstractCompiler) *ClosureModuleProcessor_UnresolvedGoogModule {
	m := &ClosureModuleProcessor_UnresolvedGoogModule{}
	m.InitUnresolvedModuleBase(m)
	m.metadata = metadata
	m.srcFileName = srcFileName
	m.path = path
	m.namespace = namespace
	m.requiresByLocalName = requires
	m.compiler = compiler
	return m
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.resolveExport(ModuleRequestResolver, String)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) ResolveExport(
	moduleRequestResolver ModuleRequestResolver, exportName string) *ResolveExportResult {
	if m.namespace.ContainsKey(exportName) {
		return ResolveExportResult_Of(m.namespace.Get(exportName))
	}
	return ResolveExportResult_NOT_FOUND
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.resolveExport(ModuleRequestResolver, String, String, Set<ExportTrace>, Set<UnresolvedModule>)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) ResolveExportWithTraces(
	moduleRequestResolver ModuleRequestResolver,
	moduleSpecifier *string,
	exportName string,
	resolveSet *jutil.LinkedHashSet[ExportTrace],
	exportStarSet *jutil.LinkedHashSet[UnresolvedModule]) *ResolveExportResult {
	return m.ResolveExport(moduleRequestResolver, exportName)
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.resolve(ModuleRequestResolver, String)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) ResolveWithModuleSpecifier(
	moduleRequestResolver ModuleRequestResolver, moduleSpecifier *string) *Module {
	if m.resolved == nil {
		// Every import creates a locally bound name.
		boundNames := jutil.NewLinkedHashMapFrom(m.GetAllResolvedImports(moduleRequestResolver))

		closureNamespace := colorPool_getOnlyElement(slices.Collect(m.metadata.GoogNamespaces().All()))
		m.resolved =
			Module_Builder_().
				Path(m.path).
				Metadata(m.metadata).
				Namespace(m.namespace).
				BoundNames(closureModuleProcessor_immutableMapCopyOf(boundNames)).
				LocalNameToLocalExport(jutil.NewLinkedHashMap[string, *Export]()).
				ClosureNamespace(&closureNamespace).
				Build()
	}
	return m.resolved
}

// closureModuleProcessor_immutableMapCopyOf is ImmutableMap.copyOf(Map),
// which rejects null values.
//
// Go-only: language-adapter — Guava ImmutableMap.copyOf.
func closureModuleProcessor_immutableMapCopyOf[V any](m *jutil.LinkedHashMap[string, V]) *jutil.LinkedHashMap[string, V] {
	out := jutil.NewLinkedHashMap[string, V]()
	for k, v := range m.All() {
		out.Put(k, jutil.CheckNotNull(v))
	}
	return out
}

// GetAllResolvedImports is a map from import bound name to binding.
//
// Java: ClosureModuleProcessor.UnresolvedGoogModule.getAllResolvedImports(ModuleRequestResolver)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) GetAllResolvedImports(
	moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashMap[string, *Binding] {
	imports := jutil.NewLinkedHashMap[string, *Binding]()

	for name := range m.requiresByLocalName.Keys() {
		b := m.ResolveImport(moduleRequestResolver, name)
		if b.Resolved() {
			imports.Put(name, b.GetBinding())
		}
	}

	return imports
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.resolveImport(ModuleRequestResolver, String)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) ResolveImport(
	moduleRequestResolver ModuleRequestResolver, name string) *ResolveExportResult {
	require := m.requiresByLocalName.Get(name)
	importRecord := require.ImportRecord()

	requested := moduleRequestResolver.Resolve(importRecord)

	if requested == nil {
		return ResolveExportResult_ERROR
	} else if importRecord.ImportName() == Export_NAMESPACE {
		// Return a binding based on the other module's metadata.
		moduleRequest := importRecord.ModuleRequest()
		return ResolveExportResult_Of(
			Binding_FromClosureImport(
				requested.Metadata(),
				importRecord.NameNode(),
				&moduleRequest,
				require.CreatedBy()))

	} else {
		moduleRequest := importRecord.ModuleRequest()
		result :=
			requested.ResolveExportWithTraces(
				moduleRequestResolver,
				&moduleRequest,
				importRecord.ImportName(),
				jutil.NewLinkedHashSet[ExportTrace](),
				jutil.NewLinkedHashSet[UnresolvedModule]())
		if !result.Found() && !result.HadError() {
			m.ReportInvalidDestructuringRequire(requested, importRecord)
			return ResolveExportResult_ERROR
		}
		forSourceInfo := importRecord.NameNode()
		if forSourceInfo == nil {
			forSourceInfo = importRecord.ImportNode()
		}
		return result.Copy(forSourceInfo, require.CreatedBy())
	}
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.metadata()
func (m *ClosureModuleProcessor_UnresolvedGoogModule) Metadata() *ModuleMetadataMap_ModuleMetadata {
	return m.metadata
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.getExportedNames(ModuleRequestResolver)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) GetExportedNames(moduleRequestResolver ModuleRequestResolver) *jutil.LinkedHashSet[string] {
	// Unsupported until such time as it becomes useful
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.getExportedNames(ModuleRequestResolver, Set<UnresolvedModule>)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) GetExportedNamesWithVisited(
	moduleRequestResolver ModuleRequestResolver, visited *jutil.LinkedHashSet[UnresolvedModule]) *jutil.LinkedHashSet[string] {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: ClosureModuleProcessor.UnresolvedGoogModule.reset()
func (m *ClosureModuleProcessor_UnresolvedGoogModule) Reset() {
	m.resolved = nil
}

// ReportInvalidDestructuringRequire reports an error given an invalid
// destructuring require.
//
// Java: ClosureModuleProcessor.UnresolvedGoogModule.reportInvalidDestructuringRequire(UnresolvedModule, Import)
func (m *ClosureModuleProcessor_UnresolvedGoogModule) ReportInvalidDestructuringRequire(
	requested UnresolvedModule, importRecord *Import) {
	additionalInfo := ""
	if unresolvedGoogModule, ok := requested.(*ClosureModuleProcessor_UnresolvedGoogModule); ok {
		// Detect some edge cases and given more helpful error messages.
		exports := unresolvedGoogModule.namespace
		if exports.ContainsKey(Export_NAMESPACE) {
			// Can't use destructuring imports on a goog.module with a default
			// export like
			//   exports = class {
			// (even if there is an assignment like `exports.Bar = 0;` later)
			additionalInfo =
				jutil.LenientFormat( // Use Strings.lenientFormat for GWT/J2CL compatability
					"\n"+
						"The goog.module \"%s\" cannot be destructured as it contains a default"+
						" export, not named exports. See %s.",
					importRecord.ModuleRequest(),
					"https://github.com/google/closure-library/wiki/goog.module%3A-an-ES6-module-like-alternative-to-goog.provide#destructuring-imports")

			if closureModuleProcessor_mayBeAccidentalDefaultExport(importRecord.ImportName(), exports) {
				// Give the user a more detailed error message, since this is a
				// tricky edge case.
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

// closureModuleProcessor_mayBeAccidentalDefaultExport returns whether the
// user appears to have confused a default export of an object literal with
// the named export object literal shorthand.
//
// Basically, `exports = {foo: 0};`, does /not/ create a 'named export' of
// 'foo' because 0 is not a name. So users cannot destructuring-require
// `const {foo} = goog.require('the.module');`. However, if instead of `0`
// the user exported a name like `bar`, the user could use a destructuring
// require.
//
// Java: ClosureModuleProcessor.mayBeAccidentalDefaultExport(String, Map<String, Binding>)
func closureModuleProcessor_mayBeAccidentalDefaultExport(
	importName string, exports *jutil.LinkedHashMap[string, *Binding]) bool {
	defaultExport := exports.Get(Export_NAMESPACE).OriginatingExport().ExportNode()
	jutil.CheckStateMsg1(
		defaultExport.MatchesName("exports") && defaultExport.GetParent().IsAssign(),
		"%s",
		defaultExport)
	exportedValue := defaultExport.GetNext()
	if !exportedValue.IsObjectLit() {
		return false
	}
	// Look for `importName` in the exported object literal.
	for key := exportedValue.GetFirstChild(); key != nil; key = key.GetNext() {
		if key.IsStringKey() && key.GetString() == importName {
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
	metadata *ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script *rhino.Node) UnresolvedModule {
	jutil.CheckArgumentMsg1(
		script.IsScript() || script.IsCall(), "Unexpected module root %s", script)
	jutil.CheckArgumentMsg(
		script.IsCall() || path != nil, "Non goog.loadModules must have a path")

	moduleProcessingCallback := NewClosureModuleProcessor_ModuleProcessingCallback(metadata)
	NodeTraversal_Traverse(p.compiler, script, moduleProcessingCallback)
	return NewClosureModuleProcessor_UnresolvedGoogModule(
		metadata,
		script.GetSourceFileName(),
		path,
		closureModuleProcessor_immutableMapCopyOf(moduleProcessingCallback.namespace),
		closureModuleProcessor_immutableMapCopyOf(moduleProcessingCallback.requiresByLocalName),
		p.compiler)
}

// ClosureModuleProcessor_ModuleProcessingCallback traverses a subtree rooted
// at a module, gathering all exports and requires.
type ClosureModuleProcessor_ModuleProcessingCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	metadata *ModuleMetadataMap_ModuleMetadata

	// The Closure namespace 'a.b.c' from the `goog.module('a.b.c');` statement
	closureNamespace string

	// Note: the following two maps are mutable because in some cases, we need
	// to check if a key has already been added before trying to add a second.

	// All named exports and explicit assignments of the `exports` object
	namespace *jutil.LinkedHashMap[string, *Binding]

	// All required/forwardDeclared local names
	requiresByLocalName *jutil.LinkedHashMap[string, *ClosureRequireProcessor_Require]

	// Whether we've come across an "exports = ..." assignment
	seenExportsAssignment bool
}

// Java: ClosureModuleProcessor.ModuleProcessingCallback(ModuleMetadata)
func NewClosureModuleProcessor_ModuleProcessingCallback(
	metadata *ModuleMetadataMap_ModuleMetadata) *ClosureModuleProcessor_ModuleProcessingCallback {
	c := &ClosureModuleProcessor_ModuleProcessingCallback{}
	c.metadata = metadata
	c.namespace = jutil.NewLinkedHashMap[string, *Binding]()
	c.requiresByLocalName = jutil.NewLinkedHashMap[string, *ClosureRequireProcessor_Require]()
	c.closureNamespace = colorPool_getOnlyElement(slices.Collect(metadata.GoogNamespaces().All()))
	c.seenExportsAssignment = false
	return c
}

// Java: ClosureModuleProcessor.ModuleProcessingCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY,
		rhino.Token_SCRIPT,
		rhino.Token_CALL, // Traverse into goog.loadModule calls.
		rhino.Token_BLOCK:
		return true
	case rhino.Token_FUNCTION:
		// Only traverse into functions that are the argument of a
		// goog.loadModule call, which is the module root. Avoid traversing
		// function declarations like:
		//     goog.module('a.b'); function (exports) { exports.x = 0; }
		return parent.IsCall() && parent == c.metadata.RootNode()
	case rhino.Token_EXPR_RESULT:
		expr := n.GetFirstChild()
		if expr.IsAssign() {
			c.MaybeInitializeExports(expr)
		} else if expr.IsGetProp() {
			c.MaybeInitializeExportsStub(expr)
		}
		return false
	case rhino.Token_CONST, rhino.Token_VAR, rhino.Token_LET:
		// Note that `let` is valid only for `goog.forwardDeclare`.
		c.MaybeInitializeRequire(n)
		return false
	default:
		return false
	}
}

// MaybeInitializeExports: if an assignment is to 'exports', adds it to the
// list of Exports.
//
// Java: ClosureModuleProcessor.ModuleProcessingCallback.maybeInitializeExports(Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) MaybeInitializeExports(assignment *rhino.Node) {
	lhs := assignment.GetFirstChild()
	rhs := assignment.GetSecondChild()
	if lhs.IsName() && lhs.GetString() == "exports" {
		// This may be a 'named exports' or may be a default export.
		// It is a 'named export' if and only if it is assigned an object
		// literal w/ string keys, whose values are all names.
		if NodeUtil_IsNamedExportsLiteral(rhs) {
			c.InitializeNamedExportsLiteral(rhs)
		} else {
			c.seenExportsAssignment = true
			c.MarkExportsAssignmentInNamespace(lhs)
		}
	} else if lhs.IsGetProp() &&
		lhs.GetFirstChild().IsName() &&
		lhs.GetFirstChild().GetString() == "exports" {
		exportedId := lhs.GetString()
		c.AddPropertyExport(exportedId, lhs)
	}
}

// MaybeInitializeExportsStub adds stub export declarations `exports.Foo;`
// to the list of Exports.
//
// Java: ClosureModuleProcessor.ModuleProcessingCallback.maybeInitializeExportsStub(Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) MaybeInitializeExportsStub(qname *rhino.Node) {
	owner := qname.GetFirstChild()
	if owner.IsName() && owner.GetString() == "exports" {
		c.AddPropertyExport(qname.GetString(), qname)
	}
}

// MarkExportsAssignmentInNamespace adds an explicit namespace export.
//
// Note that all goog.modules create an 'exports' object, but this object
// is only added to the Module namespace if there is an explicit' exports =
// ...' assignment
//
// Java: ClosureModuleProcessor.ModuleProcessingCallback.markExportsAssignmentInNamespace(Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) MarkExportsAssignmentInNamespace(exportsNode *rhino.Node) {
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

// Java: ClosureModuleProcessor.ModuleProcessingCallback.initializeNamedExportsLiteral(Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) InitializeNamedExportsLiteral(objectLit *rhino.Node) {
	for key := objectLit.GetFirstChild(); key != nil; key = key.GetNext() {
		c.AddPropertyExport(key.GetString(), key)
	}
}

// AddPropertyExport adds a named export to the list of Exports.
//
// Java: ClosureModuleProcessor.ModuleProcessingCallback.addPropertyExport(String, Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) AddPropertyExport(exportedId string, propNode *rhino.Node) {
	if c.seenExportsAssignment {
		// We've seen an assignment "exports = ...", so this is not a named
		// export.
		return
	} else if c.namespace.ContainsKey(exportedId) {
		// Ignore duplicate exports - this is an error but checked elsewhere.
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

// MaybeInitializeRequire adds a goog.require(Type) or forwardDeclare to the
// list of requiresByLocalName.
//
// Java: ClosureModuleProcessor.ModuleProcessingCallback.maybeInitializeRequire(Node)
func (c *ClosureModuleProcessor_ModuleProcessingCallback) MaybeInitializeRequire(nameDeclaration *rhino.Node) {
	for _, require := range ClosureRequireProcessor_GetAllRequires(nameDeclaration) {
		c.requiresByLocalName.PutIfAbsent(require.LocalName(), require)
	}
}
