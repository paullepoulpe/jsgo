// Port of com/google/javascript/jscomp/Es6RewriteModules.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/Es6RelativizeImportPaths.java @ closure-compiler v20261104

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteModules rewrites an ES6 module into a form that can be safely
// concatenated. Also rewrites any goog.{require,requireType,forwardDeclare,module.get}
// calls that are either in an ES module or target an ES module using
// goog.declareModuleId.
type Es6RewriteModules struct {
	compiler    AbstractCompiler
	astFactory  *AstFactory
	unknownType rhino.JSType

	preprocessorSymbolTable *PreprocessorSymbolTable

	namesToInlineByAlias map[string]ModuleRenaming_GlobalizedModuleName
	typedefs             map[string]struct{}

	moduleMetadataMap *ModuleMetadataMap
	moduleMap         *ModuleMap
	globalTypedScope  *TypedScope
	chunkOutputType   CompilerOptions_ChunkOutputType

	requireRewriter       es6RewriteRequiresForEs6Modules
	moduleRequireRewriter es6RewriteModuleRequiresCallback
	inlineAliasesCallback es6RewriteInlineAliasesCallback
	globalVarRenamer      es6RenameGlobalVars
	hasEs6GoogModules     bool
}

var (
	es6RewriteCallVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_CALL)
		return s
	}()

	es6RewriteModulesVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_IMPORT)
		s.Add(rhino.Token_EXPORT)
		s.Add(rhino.Token_SCRIPT)
		s.Add(rhino.Token_CALL)
		s.Add(rhino.Token_IMPORT_META)
		return s
	}()
)

var Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST = DiagnosticType_Error(
	"JSC_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST",
	"The left side of a goog.require() or goog.requireType() "+
		"must use ''const'' (not ''let'' or ''var'')")

var Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST = DiagnosticType_Error(
	"JSC_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST",
	"goog.requireType alias for ES6 module should be const.")

var Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST = DiagnosticType_Error(
	"JSC_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST",
	"goog.forwardDeclare alias for ES6 module should be const.")

var Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE = DiagnosticType_Warning(
	"JSC_SHOULD_IMPORT_ES6_MODULE",
	"ES6 modules should import other ES6 modules rather than goog.require them.")

// NewEs6RewriteModules creates a new Es6RewriteModules instance which can be
// used to rewrite ES6 modules to a concatenable form.
//
// Java: Es6RewriteModules(AbstractCompiler, ModuleMetadataMap, ModuleMap, PreprocessorSymbolTable, TypedScope)
func NewEs6RewriteModules(
	compiler AbstractCompiler,
	moduleMetadataMap *ModuleMetadataMap,
	moduleMap *ModuleMap,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	globalTypedScope *TypedScope) *Es6RewriteModules {
	return NewEs6RewriteModulesWithChunkOutputType(
		compiler,
		moduleMetadataMap,
		moduleMap,
		preprocessorSymbolTable,
		globalTypedScope,
		CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE)
}

// NewEs6RewriteModulesWithChunkOutputType creates a new Es6RewriteModules
// instance which can be used to rewrite ES6 modules to a concatenable form.
//
// Java: Es6RewriteModules(AbstractCompiler, ModuleMetadataMap, ModuleMap, PreprocessorSymbolTable, TypedScope, ChunkOutputType)
func NewEs6RewriteModulesWithChunkOutputType(
	compiler AbstractCompiler,
	moduleMetadataMap *ModuleMetadataMap,
	moduleMap *ModuleMap,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	globalTypedScope *TypedScope,
	chunkOutputType CompilerOptions_ChunkOutputType) *Es6RewriteModules {
	jutil.CheckNotNull(moduleMetadataMap)
	e := &Es6RewriteModules{
		compiler:                compiler,
		astFactory:              compiler.CreateAstFactory(),
		moduleMetadataMap:       moduleMetadataMap,
		moduleMap:               moduleMap,
		preprocessorSymbolTable: preprocessorSymbolTable,
		globalTypedScope:        globalTypedScope,
		unknownType:             compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE),
		chunkOutputType:         chunkOutputType,
	}
	e.requireRewriter.pass = e
	e.moduleRequireRewriter.pass = e
	e.inlineAliasesCallback.pass = e
	e.globalVarRenamer.pass = e
	return e
}

// Es6RewriteModules_IsEs6ModuleRoot returns whether the given script node
// represents an ES6 module file.
//
// Java: Es6RewriteModules.isEs6ModuleRoot(Node)
func Es6RewriteModules_IsEs6ModuleRoot(scriptNode rhino.Node) bool {
	script, ok := ast.AsScript(scriptNode)
	jutil.CheckArgumentMsg1(ok, "%s", scriptNode)
	if scriptNode.GetBooleanProp(rhino.Node_GOOG_MODULE) {
		return false
	}
	firstItem, hasFirst := script.FirstItem()
	return hasFirst && firstItem.Node().IsModuleBody()
}

func rootHasEs6ModuleScript(root rhino.Node) bool {
	for s := root.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		if !s.IsScript() || Es6RewriteModules_IsEs6ModuleRoot(s) {
			return true
		}
	}
	return false
}

// Java: Es6RewriteModules.process(Node, Node)
func (e *Es6RewriteModules) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckArgumentMsg1(externs.IsRoot(), "%s", externs)
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	hasEs6Goog := false
	for _, entry := range e.moduleMetadataMap.GetModulesByGoogNamespace().Entries() {
		if entry.Value != nil && entry.Value.IsEs6Module() {
			hasEs6Goog = true
			break
		}
	}
	e.hasEs6GoogModules = hasEs6Goog
	hasEs6Scripts := hasEs6Goog || rootHasEs6ModuleScript(externs) || rootHasEs6ModuleScript(root)
	if hasEs6Scripts {
		NodeTraversal_TraverseRoots(e.compiler, e, externs, root)
	}
	NodeUtil_RemoveFeatureFromAllScripts(root, parser.FeatureSet_Feature_MODULES, e.compiler)
	if hasEs6Scripts || e.compiler.GetAccessorSummary() == nil {
		GatherGetterAndSetterProperties_Update(e.compiler, externs, root)
	}
}

func (e *Es6RewriteModules) clearPerFileState() {
	clear(e.typedefs)
	clear(e.namesToInlineByAlias)
}

func (e *Es6RewriteModules) addTypedef(name string) {
	if e.typedefs == nil {
		e.typedefs = make(map[string]struct{})
	}
	e.typedefs[name] = struct{}{}
}

func (e *Es6RewriteModules) putInlineAlias(alias string, rep ModuleRenaming_GlobalizedModuleName) {
	if e.namesToInlineByAlias == nil {
		e.namesToInlineByAlias = make(map[string]ModuleRenaming_GlobalizedModuleName)
	}
	e.namesToInlineByAlias[alias] = rep
}

var (
	es6RewriteModules_GOOG_REQUIRE         = rhino.QualifiedName_Of("goog.require")
	es6RewriteModules_GOOG_REQUIRETYPE     = rhino.QualifiedName_Of("goog.requireType")
	es6RewriteModules_GOOG_MODULE_GET      = rhino.QualifiedName_Of("goog.module.get")
	es6RewriteModules_GOOG_FORWARDDECLARE  = rhino.QualifiedName_Of("goog.forwardDeclare")
	es6RewriteModules_GOOG_DECLAREMODULEID = rhino.QualifiedName_Of("goog.declareModuleId")
)

// es6RewriteRequiresForEs6Modules rewrites goog.require, goog.requireType,
// goog.module.get, and goog.forwardDeclare calls that import ES6 modules.
type es6RewriteRequiresForEs6Modules struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	pass        *Es6RewriteModules
	renameTable map[Es6ScopeRenameKey]string
}

func (*es6RewriteRequiresForEs6Modules) noVisitLeaves() {}
func (*es6RewriteRequiresForEs6Modules) visitTokens() *ast.TokenSet {
	return &es6RewriteCallVisitTokens
}

func (r *es6RewriteRequiresForEs6Modules) rewrite(scriptNode rhino.Node) {
	clear(r.renameTable)
	NodeTraversal_Traverse(r.pass.compiler, scriptNode, r)

	if len(r.renameTable) > 0 {
		NodeTraversal_Traverse(r.pass.compiler, scriptNode, NewEs6RenameTypeReferences(r.renameTable))
	}
}

func (r *es6RewriteRequiresForEs6Modules) putRename(root rhino.Node, oldName, newName string) {
	if r.renameTable == nil {
		r.renameTable = make(map[Es6ScopeRenameKey]string)
	}
	r.renameTable[Es6ScopeRenameKey{Root: root, Name: oldName}] = newName
}

func (r *es6RewriteRequiresForEs6Modules) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	e := r.pass
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}

	callTarget := call.Callee().Node()
	if _, isGetProp := ast.AsGetProp(callTarget); !isGetProp {
		return
	}

	isRequire := es6RewriteModules_GOOG_REQUIRE.Matches(callTarget)
	isRequireType := es6RewriteModules_GOOG_REQUIRETYPE.Matches(callTarget)
	isGet := es6RewriteModules_GOOG_MODULE_GET.Matches(callTarget)
	isForwardDeclare := es6RewriteModules_GOOG_FORWARDDECLARE.Matches(callTarget)

	if !isRequire && !isRequireType && !isGet && !isForwardDeclare {
		return
	}

	arg, hasOneArg := call.OnlyArg()
	if !hasOneArg {
		return
	}
	argStr, isStr := arg.AsStringLit()
	if !isStr {
		// Reported in CheckClosureImports, don't report here.
		return
	}

	name := argStr.Value()
	moduleMetadata := e.moduleMetadataMap.ModuleByGoogNamespace(name)
	if moduleMetadata == nil || !moduleMetadata.IsEs6Module() {
		return
	}

	if isRequire {
		currentModuleMetadata := e.moduleMetadataMap.ModuleByPath(t.GetInput().GetPath().String())
		if currentModuleMetadata != nil && currentModuleMetadata.IsEs6Module() {
			t.Report(n, Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE)
		}
	}

	if isGet && t.InGlobalHoistScope() {
		return
	}

	statementNode := NodeUtil_GetEnclosingStatement(n)
	decl, isNameDecl := ast.AsNameDeclaration(statementNode)
	var firstDecl ast.Declarator
	if isNameDecl {
		var hasOne bool
		firstDecl, hasOne = decl.OnlyDeclarator()
		isNameDecl = hasOne && firstDecl.Node() == n.GetParent()
		if isNameDecl {
			if destructuringLhs, isDestruct := firstDecl.AsDestructuringLhs(); isDestruct && !isForwardDeclare {
				isNameDecl = ClosureCheckModule_IsValidDestructuringImport(destructuringLhs)
			}
		}
	}
	if isNameDecl {
		if destructuringLhs, isDestruct := firstDecl.AsDestructuringLhs(); isDestruct {
			if isForwardDeclare {
				t.Report(n, ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE)
				return
			}
			objPat := ast.MustObjectPattern(destructuringLhs.Pattern().Node())
			if isRequireType {
				if !statementNode.IsConst() {
					t.Report(statementNode, Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST)
					return
				}
				baseGlobal := ModuleRenaming_GetGlobalNameWithGoogNamespace(moduleMetadata, name)
				scopeRoot := t.GetScopeRoot()
				for prop := range objPat.Props() {
					sk := ast.MustStringKey(prop.Node())
					local := ast.MustName(sk.Value().Node())
					r.putRename(scopeRoot, local.Ident(), baseGlobal.Getprop(sk.Key()).Join())
				}
			} else {
				baseGlobal := e.getGlobalNameAndType(moduleMetadata, name /* isFromMissingModuleOrProvide= */, false)
				for prop := range objPat.Props() {
					sk := ast.MustStringKey(prop.Node())
					local := ast.MustName(sk.Value().Node())
					globalName := baseGlobal.Getprop(sk.Key())
					constNode := e.astFactory.CreateSingleConstNameDeclaration(
						local.Ident(), globalName.ToQname(e.astFactory))
					constNode.SrcrefTree(sk.Node())
					constNode.InsertBefore(statementNode)
				}
			}
			statementNode.Detach()
			t.ReportCodeChange()
		} else {
			if isForwardDeclare || isRequireType {
				if !statementNode.IsConst() {
					diagnostic := Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST
					if isForwardDeclare {
						diagnostic = Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST
					}
					t.Report(statementNode, diagnostic)
					return
				}
				lhsName := firstDecl.MustName()
				r.putRename(
					t.GetScopeRoot(),
					lhsName.Ident(),
					ModuleRenaming_GetGlobalNameWithGoogNamespace(moduleMetadata, name).Join())
				statementNode.Detach()
				t.ReportCodeChange()
			} else {
				n.ReplaceWith(
					e.astFactory.
						CreateName(
							ModuleRenaming_GetGlobalNameWithGoogNamespace(moduleMetadata, name).GetRoot(), AstFactory_Type_(n)).
						SrcrefTree(n))
				t.ReportCodeChange()
			}
		}
	} else {
		if isForwardDeclare || isRequireType {
			r.putRename(
				t.GetScopeRoot(), name, ModuleRenaming_GetGlobalNameWithGoogNamespace(moduleMetadata, name).Join())
			statementNode.Detach()
		} else if exprStmt, ok := ast.AsExprStmt(statementNode); ok && exprStmt.Expr().Node() == n {
			statementNode.Detach()
		} else {
			n.ReplaceWith(
				e.astFactory.
					CreateName(
						ModuleRenaming_GetGlobalNameWithGoogNamespace(moduleMetadata, name).GetRoot(), AstFactory_Type_(n)).
					SrcrefTree(n))
		}
		t.ReportCodeChange()
	}
}

func (*Es6RewriteModules) visitTokens() *ast.TokenSet { return &es6RewriteModulesVisitTokens }

// Java: Es6RewriteModules.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if _, ok := ast.AsScript(n); ok {
		isEs6Root := Es6RewriteModules_IsEs6ModuleRoot(n)
		if !e.hasEs6GoogModules && !isEs6Root {
			return false
		}
		// Trigger creation of the global scope before inserting any synthetic code.
		nodeTraversal.GetScope()
		if e.hasEs6GoogModules {
			e.requireRewriter.rewrite(n)
		}
		if isEs6Root {
			e.clearPerFileState()
		} else {
			return false
		}
	}
	return true
}

// Java: Es6RewriteModules.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteModules) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IMPORT:
		e.maybeWarnExternModule(t, n, parent)
		e.visitImport(t, ast.MustImport(n), parent)
	case rhino.Token_EXPORT:
		e.maybeWarnExternModule(t, n, parent)
		e.visitExport(t, ast.MustExport(n), parent)
	case rhino.Token_SCRIPT:
		e.visitScript(t, n)
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if es6RewriteModules_GOOG_DECLAREMODULEID.Matches(call.Callee().Node()) {
			n.GetParent().Detach()
		}
	case rhino.Token_IMPORT_META:
		if e.chunkOutputType != CompilerOptions_ChunkOutputType_ES_MODULES {
			t.Report(
				n,
				TranspilationUtil_CANNOT_CONVERT,
				"import.meta. Use --chunk_output_type=ES_MODULES to allow passthrough support.")
		}
	}
}

func (e *Es6RewriteModules) maybeWarnExternModule(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	jutil.CheckState(parent.IsModuleBody())
	if parent.IsFromExterns() && !NodeUtil_IsFromTypeSummary(parent.GetParent()) {
		t.Report(n, TranspilationUtil_CANNOT_CONVERT_YET, "ES6 modules in externs")
	}
}

func (e *Es6RewriteModules) visitImport(t *NodeTraversal, importDecl ast.Import, parent rhino.Node) {
	jutil.CheckArgumentMsg1(parent.IsModuleBody(), "%s", parent)
	importDeclNode := importDecl.Node()
	sourceLit := importDecl.Source()
	importName := sourceLit.Value()
	if namespace, ok := strings.CutPrefix(importName, "goog:"); ok {
		m := e.moduleMetadataMap.ModuleByGoogNamespace(namespace)
		if m == nil {
			t.Report(importDeclNode, ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, namespace)
		} else {
			jutil.CheckState(m.IsEs6Module() || m.IsGoogModule() || m.IsGoogProvide())
		}
	} else {
		inputPath := t.GetInput().GetPath()
		modulePath := inputPath.ResolveJsModule(
			importName,
			processCommonJSModules_sourceFileName(importDeclNode),
			importDeclNode.GetLineno(),
			importDeclNode.GetCharno())
		if modulePath == nil {
			modulePath = inputPath.ResolveModuleAsPath(importName)
		}
		e.maybeAddImportedFileReferenceToSymbolTable(sourceLit.Node(), modulePath.String())
	}

	if e.preprocessorSymbolTable != nil {
		if bindings, hasBindings := importDecl.Bindings(); hasBindings {
			sourceName, hasSourceName := t.SourceNameOk()
			if specs, isSpecs := bindings.AsImportSpecs(); isSpecs {
				for spec := range specs.Specs() {
					e.maybeAddAliasToSymbolTable(spec.Imported().Node(), sourceName, hasSourceName)
				}
			} else if star, isStar := bindings.AsImportStar(); isStar {
				e.maybeAddAliasToSymbolTable(star.Node(), sourceName, hasSourceName)
			}
		}
	}

	importDeclNode.Detach()
	t.ReportCodeChange()
}

func (e *Es6RewriteModules) visitExport(t *NodeTraversal, export ast.Export, parent rhino.Node) {
	jutil.CheckArgumentMsg1(parent.IsModuleBody(), "%s", parent)
	exportNode := export.Node()
	if export.IsDefault() {
		child := export.Body().Node()
		hasDeclName := false
		if fn, ok := ast.AsFunction(child); ok {
			_, hasDeclName = fn.Name()
		} else if cls, ok := ast.AsClass(child); ok {
			_, hasDeclName = cls.Name()
		}

		if hasDeclName {
			exportNode.ReplaceWith(child.Detach())
		} else {
			varDecl := e.astFactory.CreateSingleVarNameDeclarationWithValue(
				ModuleRenaming_DEFAULT_EXPORT_VAR_PREFIX, child.Detach())
			varDecl.SetJSDocInfo(child.GetJSDocInfo())
			child.SetJSDocInfo(nil)
			varDecl.SrcrefTreeIfMissing(exportNode)
			exportNode.ReplaceWith(varDecl)
		}
		t.ReportCodeChange()
	} else if _, hasSource := export.Source(); export.IsAllFrom() || hasSource || export.Body().Node().IsExportSpecs() {
		exportNode.Detach()
		t.ReportCodeChange()
	} else {
		e.visitExportDeclaration(t, export)
	}
}

func (e *Es6RewriteModules) visitExportNameDeclaration(declaration rhino.Node) {
	if info := declaration.GetJSDocInfo(); info != nil && info.HasTypedefType() {
		NodeUtil_VisitLhsNodesInNode(declaration, func(lhs rhino.Node) { e.addTypedef(lhs.GetString()) })
	}
}

func (e *Es6RewriteModules) visitExportDeclaration(t *NodeTraversal, export ast.Export) {
	declaration := export.Body().Node()
	if _, ok := ast.AsNameDeclaration(declaration); ok {
		e.visitExportNameDeclaration(declaration)
	}
	export.Node().ReplaceWith(declaration.Detach())
	t.ReportCodeChange()
}

func (e *Es6RewriteModules) inlineModuleToGlobalScope(moduleBody ast.ModuleBody) {
	moduleNode := moduleBody.Node()
	scriptNode := moduleNode.GetParent()
	moduleNode.Detach()
	scriptNode.AddChildrenToFront(moduleNode.RemoveChildren())
}

func (e *Es6RewriteModules) visitScript(t *NodeTraversal, script rhino.Node) {
	moduleBody := ast.MustScript(script).MustFirstItem().MustModuleBody()
	// TypedScopeCreator sets the module object type on the MODULE_BODY during type checking.
	moduleObjectType := AstFactory_Type_(moduleBody.Node())
	e.inlineModuleToGlobalScope(moduleBody)

	ClosureRewriteModule_CheckAndSetStrictModeDirective(t, script)

	thisModule := e.moduleMap.GetModuleByPath(t.GetInput().GetPath())
	qualifiedName := ModuleRenaming_GetGlobalNameWithoutGoogNamespace(thisModule.Metadata())
	jutil.CheckStateMsg1(qualifiedName.IsSimple(), "Unexpected qualified name %s", qualifiedName)
	moduleName := qualifiedName.GetRoot()

	moduleVar := e.createExportsObject(moduleName, t, script, thisModule, moduleObjectType)

	// Rename vars to not conflict in global scope.
	e.globalVarRenamer.thisModule = thisModule
	NodeTraversal_Traverse(e.compiler, script, &e.globalVarRenamer)
	e.globalVarRenamer.thisModule = nil

	// Rename the exports object to something we can reference later.
	moduleVarName := ast.MustNameDeclaration(moduleVar).MustOnlyDeclarator().MustName()
	moduleVarName.SetIdent(moduleName)
	moduleVar.MakeNonIndexableRecursive()
	e.declareGlobalVariable(moduleVarName.Node(), t)

	e.rewriteRequires(script)

	t.ReportCodeChange()
}

func (e *Es6RewriteModules) createExportsObject(
	moduleName string,
	t *NodeTraversal,
	script rhino.Node,
	thisModule *Module,
	moduleObjectType AstFactory_Type) rhino.Node {
	moduleObject := e.astFactory.CreateObjectLitWithType(moduleObjectType)
	moduleVar := e.astFactory.CreateSingleVarNameDeclarationWithValue("$jscomp$tmp$exports$module$name", moduleObject)
	moduleVarLhs := ast.MustNameDeclaration(moduleVar).MustOnlyDeclarator().Node()
	moduleVarLhs.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	infoBuilder := rhino.JSDocInfo_Builder_()
	infoBuilder.RecordConstancy()
	moduleVar.SetJSDocInfo(infoBuilder.Build())
	moduleVarLhs.SetDeclaredConstantVar(true)
	script.AddChildToBack(moduleVar.SrcrefTreeIfMissing(script))

	for _, entry := range thisModule.NamespaceMap().Entries() {
		exportedName, binding := entry.Key, entry.Value
		nodeForSourceInfo := binding.SourceNode()
		mutated := binding.IsMutated()
		boundVariableQualifiedName := ModuleRenaming_GetGlobalNameWithBinding(binding)
		jutil.CheckStateMsg1(
			boundVariableQualifiedName.IsSimple(),
			"unexpected qualified name: %s",
			boundVariableQualifiedName)
		boundVariableName := boundVariableQualifiedName.GetRoot()

		getProp := e.astFactory.CreateGetPropWithoutColor(
			e.astFactory.CreateName(moduleName, moduleObjectType), exportedName)
		getProp.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)

		if _, isTypedef := e.typedefs[exportedName]; isTypedef {
			builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
			typeExpr := rhino.NewJSTypeExpression(
				e.astFactory.CreateString(exportedName).Srcref(nodeForSourceInfo),
				processCommonJSModules_sourceFileName(script))
			builder.RecordTypedef(typeExpr)
			getProp.SetJSDocInfo(builder.Build())
			exprResult := e.astFactory.ExprResult(getProp).SrcrefTreeIfMissing(nodeForSourceInfo)
			script.AddChildToBack(exprResult)
		} else if mutated {
			globalExportName := e.astFactory.CreateName(boundVariableName, AstFactory_Type_(getProp))
			e.addGetterExport(script, nodeForSourceInfo, moduleObject, exportedName, globalExportName)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_GETTER, e.compiler)
		} else {
			assign := e.astFactory.CreateAssign(
				getProp, e.astFactory.CreateName(boundVariableName, AstFactory_Type_(getProp)))
			builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
			builder.RecordConstancy()
			assign.SetJSDocInfo(builder.Build())
			script.AddChildToBack(e.astFactory.ExprResult(assign).SrcrefTreeIfMissing(nodeForSourceInfo))
		}
	}

	return moduleVar
}

func (e *Es6RewriteModules) addGetterExport(
	script rhino.Node, forSourceInfo rhino.Node, objLit rhino.Node, exportedName string, value rhino.Node) {
	getter := e.astFactory.CreateGetterDef(exportedName, value)
	getter.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	objLit.AddChildToBack(getter)

	if !e.astFactory.IsAddingTypes() {
		builder := rhino.JSDocInfo_Builder_().ParseDocumentation()
		builder.RecordReturnType(
			rhino.NewJSTypeExpression(
				rhino.NewNode(rhino.Token_QMARK).Srcref(forSourceInfo), processCommonJSModules_sourceFileName(script)))
		getter.SetJSDocInfo(builder.Build())
	} else {
		getter.SetJSType(e.compiler.GetTypeRegistry().CreateFunctionType(value.GetJSType()))
	}

	getter.SrcrefTreeIfMissing(forSourceInfo)
	e.compiler.ReportChangeToEnclosingScope(ast.MustGetterDef(getter).Function().Body().Node())
	e.compiler.ReportChangeToEnclosingScope(getter)
}

type es6RewriteModuleRequiresCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	pass *Es6RewriteModules
}

func (*es6RewriteModuleRequiresCallback) noVisitLeaves() {}
func (*es6RewriteModuleRequiresCallback) visitTokens() *ast.TokenSet {
	return &es6RewriteCallVisitTokens
}

func (r *es6RewriteModuleRequiresCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if call, ok := ast.AsCall(n); ok {
		e := r.pass
		fn := call.Callee().Node()
		if es6RewriteModules_GOOG_REQUIRE.Matches(fn) || es6RewriteModules_GOOG_REQUIRETYPE.Matches(fn) {
			e.visitRequireOrGet(t, n, parent /* isRequire= */, true)
		} else if es6RewriteModules_GOOG_MODULE_GET.Matches(fn) {
			e.visitGoogModuleGet(t, n, parent)
		}
	}
}

type es6RewriteInlineAliasesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	pass *Es6RewriteModules
}

func (r *es6RewriteInlineAliasesCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	e := r.pass
	if info := n.GetJSDocInfo(); info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			e.inlineAliasedTypes(t, it.Node())
		}
	}

	if nm, ok := ast.AsName(n); ok {
		ident := nm.Ident()
		if replacementName, found := e.namesToInlineByAlias[ident]; found {
			v := t.GetScope().GetVarByNode(n)
			if v == nil || v.GetNameNode() != n {
				replacement := replacementName.ToQname(e.astFactory).SrcrefTree(n)
				n.ReplaceWith(replacement)
			}
		}
	}
}

func (e *Es6RewriteModules) rewriteRequires(script rhino.Node) {
	NodeTraversal_Traverse(e.compiler, script, &e.moduleRequireRewriter)
	if len(e.namesToInlineByAlias) == 0 {
		return
	}
	NodeTraversal_Traverse(e.compiler, script, &e.inlineAliasesCallback)
}

func (e *Es6RewriteModules) inlineAliasedTypes(t *NodeTraversal, typeNode rhino.Node) {
	if typeNode.IsStringLit() {
		name := typeNode.GetString()
		base, rest, hasDot := strings.Cut(name, ".")
		if replacement, ok := e.namesToInlineByAlias[base]; ok && t.GetScope().GetVar(base) == nil {
			newTypeName := replacement.AliasName().Join()
			if hasDot {
				newTypeName += "." + rest
			}
			typeNode.SetOriginalNameString(name)
			typeNode.SetString(newTypeName)
			t.ReportCodeChange()
		}
	}
	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		e.inlineAliasedTypes(t, child)
	}
}

func (e *Es6RewriteModules) visitGoogModuleGet(t *NodeTraversal, getCall rhino.Node, parent rhino.Node) {
	arg, hasOneArg := ast.MustCall(getCall).OnlyArg()
	if !hasOneArg || !arg.Node().IsStringLit() {
		t.Report(getCall, ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
		return
	}

	if t.InGlobalHoistScope() {
		t.Report(getCall, ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
		return
	}

	e.visitRequireOrGet(t, getCall, parent /* isRequire= */, false)
}

func (e *Es6RewriteModules) fallbackMetadataForNamespace(namespace string) *ModuleMetadataMap_ModuleMetadata {
	builder := NewModuleMetadataBuilder().
		ModuleType(ModuleMetadataMap_ModuleType_GOOG_PROVIDE).
		UsesClosure(true).
		IsTestOnly(false)
	builder.GoogNamespacesBuilder().Add(namespace)
	return builder.Build()
}

func (e *Es6RewriteModules) visitRequireOrGet(
	t *NodeTraversal, requireCall rhino.Node, parent rhino.Node, isRequire bool) {
	arg, hasOneArg := ast.MustCall(requireCall).OnlyArg()
	argStr, isStr := arg.AsStringLit()
	if !hasOneArg || !isStr {
		t.Report(requireCall, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
		return
	}

	if isRequire && !t.GetScope().IsGlobal() {
		t.Report(requireCall, ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
		return
	}

	namespace := argStr.Value()
	isStoredInDeclaration := NodeUtil_IsDeclaration(parent.GetParent())

	if isStoredInDeclaration && !parent.GetParent().IsConst() {
		e.compiler.Report(JSError_MakeWithNode(parent.GetParent(), Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST))
	}

	m := e.moduleMetadataMap.ModuleByGoogNamespace(namespace)
	isFromFallbackMetadata := m == nil
	if isFromFallbackMetadata {
		t.Report(requireCall, ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, namespace)
		m = e.fallbackMetadataForNamespace(namespace)
	}

	if isStoredInDeclaration {
		if isRequire {
			var toDetach rhino.Node
			if destructLhs, ok := ast.AsDestructuringLhs(parent); ok {
				objPat := ast.MustObjectPattern(destructLhs.Pattern().Node())
				toDetach = parent.GetParent()
				baseGlobal := e.getGlobalNameAndType(m, namespace, isFromFallbackMetadata)
				for prop := range objPat.Props() {
					sk := ast.MustStringKey(prop.Node())
					local := ast.MustName(sk.Value().Node())
					e.putInlineAlias(local.Ident(), baseGlobal.Getprop(sk.Key()))
				}
			} else if nm, ok := ast.AsName(parent); ok {
				alias := e.getGlobalNameAndType(m, namespace, isFromFallbackMetadata)
				e.putInlineAlias(nm.Ident(), alias)
				toDetach = parent.GetParent()
			} else {
				jutil.CheckState(parent.IsExprResult())
				toDetach = parent
			}
			toDetach.Detach()
		} else {
			name := e.getGlobalNameAndType(m, namespace, isFromFallbackMetadata)
			replacement := name.ToQname(e.astFactory).SrcrefTree(requireCall)
			requireCall.ReplaceWith(replacement)
		}
	} else {
		jutil.CheckState(requireCall.GetParent().IsExprResult())
		requireCall.GetParent().Detach()
	}
}

func (e *Es6RewriteModules) getGlobalNameAndType(
	metadata *ModuleMetadataMap_ModuleMetadata,
	googNamespace string,
	isFromMissingModuleOrProvide bool) ModuleRenaming_GlobalizedModuleName {
	if isFromMissingModuleOrProvide {
		globalName := ModuleRenaming_GetGlobalNameWithGoogNamespace(metadata, googNamespace)
		return ModuleRenaming_GlobalizedModuleName_Create(globalName, e.unknownType)
	}
	return ModuleRenaming_GlobalizedModuleName_CreateWithGoogNamespace(metadata, googNamespace, e.globalTypedScope)
}

// es6RenameGlobalVars appends a module suffix to module-scoped variables and
// rewrites references to imported bindings.
type es6RenameGlobalVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	pass       *Es6RewriteModules
	thisModule *Module
}

func (r *es6RenameGlobalVars) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if info := n.GetJSDocInfo(); info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			r.fixTypeNode(t, it.Node())
		}
	}

	nm, ok := ast.AsName(n)
	if !ok {
		return
	}
	v := t.GetScope().GetVarByNode(n)
	if v != nil && v.IsGlobal() {
		name := nm.Ident()
		newName := ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(r.thisModule.Metadata(), name).Join()
		nm.SetIdent(newName)
		n.SetOriginalNameString(name)
		t.ReportCodeChangeWithNode(n)
		if NodeUtil_IsDeclarationLValue(n) {
			r.pass.declareGlobalVariable(n, t)
		}
	} else if v == nil {
		name := nm.Ident()
		if binding := r.thisModule.BoundNamesMap().Get(name); binding != nil {
			if r.pass.preprocessorSymbolTable != nil {
				sourceName, hasSourceName := t.SourceNameOk()
				r.pass.maybeAddAliasToSymbolTable(n, sourceName, hasSourceName)
			}

			replacement := r.replaceBindingRef(n, binding)

			if _, isName := ast.AsName(replacement); isName {
				if call, isCall := ast.AsCall(parent); isCall && call.Callee().Node() == n && call.IsFreeCall() {
					call.SetFreeCall(true)
				}
			}

			if NodeUtil_IsDeclarationLValue(n) {
				r.pass.declareGlobalVariable(n, t)
			}
			t.ReportCodeChange()
		}
	}
}

func (r *es6RenameGlobalVars) replaceBindingRef(n rhino.Node, binding *Binding) rhino.Node {
	e := r.pass
	jutil.CheckState(n.IsName())

	for binding.IsModuleNamespace() && binding.Metadata().IsEs6Module() {
		gp, ok := ast.AsGetProp(n.GetParent())
		if !ok {
			break
		}
		propertyName := gp.Property()
		m := e.moduleMap.GetModuleByPath(binding.Metadata().Path())
		if propBinding := m.NamespaceMap().Get(propertyName); propBinding != nil {
			binding = propBinding
			n = gp.Node()
		} else {
			break
		}
	}

	globalName := ModuleRenaming_GetGlobalNameWithBinding(binding)
	var newNode rhino.Node
	if !globalName.IsSimple() {
		root := globalName.GetRoot()
		if e.globalTypedScope != nil && e.globalTypedScope.HasSlot(root) {
			newNode = e.astFactory.CreateQNameUsingJSTypeInfo(e.globalTypedScope, globalName.Join())
		} else {
			newNode = e.astFactory.CreateQNameWithUnknownType(globalName.Join())
		}
	} else {
		newNode = e.astFactory.CreateName(globalName.GetRoot(), AstFactory_Type_(n))
	}

	origName := n.GetString()
	n.ReplaceWith(newNode)
	newNode.SrcrefTree(n)
	newNode.SetOriginalNameString(origName)
	return newNode
}

func (r *es6RenameGlobalVars) fixTypeNode(t *NodeTraversal, typeNode rhino.Node) {
	e := r.pass
	if typeNode.IsStringLit() {
		thisModule := r.thisModule
		name := typeNode.GetString()
		baseName, rest, hasDot := strings.Cut(name, ".")
		v := t.GetScope().GetVar(baseName)
		if v != nil && v.IsGlobal() {
			newName := ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(thisModule.Metadata(), baseName).Join()
			if hasDot {
				newName += "." + rest
			}
			r.maybeSetNewName(t, typeNode, name, newName)
		} else if v == nil {
			if binding := thisModule.BoundNamesMap().Get(baseName); binding != nil {
				var propChain []string
				if hasDot {
					propChain = strings.Split(rest, ".")
				}
				globalName := ModuleRenaming_GetGlobalNameForJsDoc(e.moduleMap, binding, propChain)
				r.maybeSetNewName(t, typeNode, name, globalName)

				if e.preprocessorSymbolTable != nil {
					onlyBaseName := rhino.Node_NewString(baseName).Srcref(typeNode)
					onlyBaseName.SetLength(jutil.StringLength(baseName))
					sourceName, hasSourceName := t.SourceNameOk()
					e.maybeAddAliasToSymbolTable(onlyBaseName, sourceName, hasSourceName)
				}
			}
		}

		typeNode.SetOriginalNameString(name)
	}

	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		r.fixTypeNode(t, child)
	}
}

func (r *es6RenameGlobalVars) maybeSetNewName(t *NodeTraversal, node rhino.Node, name string, newName string) {
	if name != newName {
		node.SetString(newName)
		node.SetOriginalNameString(name)
		t.ReportCodeChange()
	}
}

func (e *Es6RewriteModules) maybeAddAliasToSymbolTable(n rhino.Node, module string, hasModule bool) {
	if e.preprocessorSymbolTable == nil {
		return
	}
	n.PutBooleanProp(rhino.Node_MODULE_ALIAS, true)
	moduleName := "null"
	if hasModule {
		moduleName = module
	}
	nodeName := "null"
	if n.IsStringLit() || n.IsImportStar() {
		nodeName = n.GetString()
	} else if qn, ok := e.preprocessorSymbolTable.GetQualifiedName(n); ok {
		nodeName = qn
	}
	e.preprocessorSymbolTable.AddReferenceWithName(n, "alias_"+moduleName+"_"+nodeName)
}

func (e *Es6RewriteModules) maybeAddImportedFileReferenceToSymbolTable(
	importStringNode rhino.Node, importedFilePath string) {
	if e.preprocessorSymbolTable == nil {
		return
	}
	if e.preprocessorSymbolTable.GetSlot(importedFilePath) == nil {
		scriptNode := e.compiler.GetScriptNode(importedFilePath)
		if !scriptNode.IsNil() {
			e.preprocessorSymbolTable.AddReferenceWithName(scriptNode, importedFilePath)
		}
	}
	e.preprocessorSymbolTable.AddReferenceWithName(importStringNode, importedFilePath)
}

func (e *Es6RewriteModules) declareGlobalVariable(n rhino.Node, t *NodeTraversal) {
	nm := ast.MustName(n)
	if !e.astFactory.IsAddingTypes() {
		return
	}
	jutil.CheckNotNull(e.globalTypedScope)

	name := nm.Ident()
	if e.globalTypedScope.HasOwnSlot(name) {
		t.Report(t.GetCurrentScript(), ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT, name)
	} else {
		type_ := jutil.CheckNotNull(n.GetJSType())
		e.globalTypedScope.Declare(name, n, type_, t.GetInput(), false)
	}
}

// Es6RelativizeImportPaths rewrites ES6 import paths to be relative after
// resolving according to the compiler's module resolver.
//
// Useful for servers that wish to preserve ES6 modules, meaning their paths
// need to be valid in the browser.
type Es6RelativizeImportPaths struct {
	compiler AbstractCompiler
}

// Java: Es6RelativizeImportPaths(AbstractCompiler)
func NewEs6RelativizeImportPaths(compiler AbstractCompiler) *Es6RelativizeImportPaths {
	return &Es6RelativizeImportPaths{compiler: compiler}
}

// Java: Es6RelativizeImportPaths.process(Node, Node)
func (e *Es6RelativizeImportPaths) Process(externs rhino.Node, root rhino.Node) {
	var rewriter Es6RelativizeImportPaths_Rewriter
	for item := range ast.MustRoot(root).Items() {
		script := item.Node()
		if Es6RewriteModules_IsEs6ModuleRoot(script) {
			NodeTraversal_Traverse(e.compiler, script, &rewriter)
		}
	}
}

// Es6RelativizeImportPaths_Rewriter is the private static nested class
// Es6RelativizeImportPaths.Rewriter.
type Es6RelativizeImportPaths_Rewriter struct {
	NodeTraversal_AbstractPreOrderCallbackBase
}

// Java: Es6RelativizeImportPaths.Rewriter()
func NewEs6RelativizeImportPaths_Rewriter() *Es6RelativizeImportPaths_Rewriter {
	return &Es6RelativizeImportPaths_Rewriter{}
}

// Java: Es6RelativizeImportPaths.Rewriter.shouldTraverse(NodeTraversal, Node, Node)
func (r *Es6RelativizeImportPaths_Rewriter) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_MODULE_BODY, rhino.Token_SCRIPT:
		return true
	case rhino.Token_IMPORT:
		r.VisitImport(nodeTraversal, n)
		return false
	case rhino.Token_EXPORT:
		r.VisitExport(nodeTraversal, n)
		return false
	default:
		return false
	}
}

// Java: Es6RelativizeImportPaths.Rewriter.visitImport(NodeTraversal, Node)
func (r *Es6RelativizeImportPaths_Rewriter) VisitImport(t *NodeTraversal, importDecl rhino.Node) {
	r.relativizeModuleSource(t, ast.MustImport(importDecl).Source())
}

func (r *Es6RelativizeImportPaths_Rewriter) relativizeModuleSource(t *NodeTraversal, specifierLit ast.StringLit) {
	specifier := specifierLit.Value()

	// Leave relative and truly absolute paths (those with a scheme) alone. Only transform
	// absolute paths without a scheme (starting with "/") and ambiguous paths.
	if ModuleLoader_IsRelativeIdentifier(specifier) {
		return
	}
	specifierURI, err := jutil.ParseURI(specifier)
	if err != nil || specifierURI.IsAbsolute() {
		return
	}

	scriptPath := t.GetInput().GetPath().String()

	// If the script path has a scheme / host / port then just use the path part.
	scriptURI, err := jutil.ParseURI(scriptPath)
	if err != nil {
		return
	}
	uriPath, hasPath := scriptURI.Path()
	jutil.CheckState(hasPath)
	scriptPath = uriPath

	newSpecifier := t.GetInput().GetPath().ResolveModuleAsPath(specifier).String()

	// If a module root is stripped then this won't start with "/" when it probably should.
	if !strings.HasPrefix(newSpecifier, "/") {
		newSpecifier = "/" + newSpecifier
	}

	newSpecifier = jutil.Path_Of(scriptPath).GetParent().Relativize(jutil.Path_Of(newSpecifier)).String()

	// Relativizing two paths with the same directory yields an ambiguous path rather than one
	// starting with "./".
	if ModuleLoader_IsAmbiguousIdentifier(newSpecifier) {
		newSpecifier = "./" + newSpecifier
	}

	if newSpecifier != specifier {
		specifierLit.SetValue(newSpecifier)
		t.ReportCodeChangeWithNode(specifierLit.Node())
	}
}

// Java: Es6RelativizeImportPaths.visitExport(NodeTraversal, Node)
func (r *Es6RelativizeImportPaths_Rewriter) VisitExport(t *NodeTraversal, export rhino.Node) {
	if src, ok := ast.MustExport(export).Source(); ok {
		r.relativizeModuleSource(t, src)
	}
}
