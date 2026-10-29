// Port of com/google/javascript/jscomp/RewriteDynamicImports.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Es6RewriteScriptsToModules.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ForbidDynamicImportUsage.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
 * Copyright 2017 The Closure Compiler Authors.
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
 * Copyright 2020 The Closure Compiler Authors.
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

var (
	rewriteDynamicImportsVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_DYNAMIC_IMPORT)
		return s
	}()
	rewriteDynamicImportsOriginalName = "import"
)

// RewriteDynamicImports rewrites dynamic import expressions to account for
// bundling and module rewriting. Since dynamic imports cannot be fully
// polyfilled, optionally support replacing the expression with a function
// call (alias function) which indicates an external polyfill is utilized.
//
// If the import specifier is a string literal and the module resolver
// recognizes the target, the pass retargets the specifier to the correct
// output chunk.
//
// Example from:
//
//	import('./input-module1.js');
//
// To:
//
//	imprt_('./output-chunk0.js').then(function() { return module$output$chunk0; });
//
// Unlike other "rewrite" passes for which the DefaultPassConfig checks the
// output language level using options.needsTranspilationOf(feature) before
// adding a pass, it does not check the language level before adding this
// particular pass. Even when supported by the output level, we still invoke
// this pass if options.shouldAllowDynamicImport() is set. The only
// difference is that for ES2020, the dynamic imports syntax can pass
// through if aliasing is not requested, while for < ES2020 they must get
// rewritten using the alias call.
//
//	options.shouldAllowDynamicImport() === false?
//	    ForbidDynamicImportUsage :
//	    (output >= ES2020) ?
//	        don't *really* rewrite it unless aliasing is requested :
//	        rewrite it using the alias call
//
// TODO: b/291319705 For the above reason, the pass should be better named as
// MaybeRewriteDynamicImports.
type RewriteDynamicImports struct {
	compiler                         AbstractCompiler
	astFactory                       *AstFactory
	alias                            *string
	requiresAliasing                 bool
	shouldWrapDynamicImportCallbacks bool
	dynamicImportsRemoved            bool
	wrappedDynamicImportCallback     bool
}

func (*RewriteDynamicImports) noVisitLeaves() {}

func (*RewriteDynamicImports) visitTokens() *ast.TokenSet {
	return &rewriteDynamicImportsVisitTokens
}

var RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED = DiagnosticType_Warning(
	"JSC_DYNAMIC_IMPORT_ALIASING_REQUIRED",
	"Dynamic import expressions should be aliased for language level. "+
		"Use the --dynamic_import_alias flag.")

var RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS = DiagnosticType_Error(
	"JSC_DYNAMIC_IMPORT_INVALID_ALIAS", "Dynamic import alias is not a valid name")

// NewRewriteDynamicImports: compiler is the compiler.
//
// Java: RewriteDynamicImports(AbstractCompiler, String, ChunkOutputType)
func NewRewriteDynamicImports(
	compiler AbstractCompiler, alias *string, chunkOutputType CompilerOptions_ChunkOutputType) *RewriteDynamicImports {
	r := &RewriteDynamicImports{}
	r.dynamicImportsRemoved = false
	r.wrappedDynamicImportCallback = false
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.alias = alias
	r.requiresAliasing =
		!compiler.GetOptions().GetOutputFeatureSet().ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	r.shouldWrapDynamicImportCallbacks = chunkOutputType == CompilerOptions_ChunkOutputType_ES_MODULES
	return r
}

// Java: RewriteDynamicImports.aliasIsValid()
func (r *RewriteDynamicImports) AliasIsValid() bool {
	return r.alias != nil &&
		(*r.alias == "import" ||
			NodeUtil_IsValidQualifiedNameFeatureSet(r.compiler.GetAllowableFeatures(), *r.alias))
}

// Java: RewriteDynamicImports.process(Node, Node)
func (r *RewriteDynamicImports) Process(externs rhino.Node, root rhino.Node) {
	r.dynamicImportsRemoved = false
	jutil.CheckArgumentMsg1(externs.IsRoot(), "%s", externs)
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)

	NodeTraversal_Traverse(r.compiler, root, r)
	if r.wrappedDynamicImportCallback {
		r.InjectWrappingFunctionExtern()
	}

	if r.dynamicImportsRemoved {
		r.compiler.MarkFeatureNotAllowed(parser.FeatureSet_Feature_DYNAMIC_IMPORT)

		// This pass removes dynamic import, but adds arrow functions.
		// It is unusual to call this NodeUtil method instead of the TranspilationPasses one. This
		// pass runs as a stage1 checks pass and does not have access to {@code TranspilationPasses}.
		// Adding that dep produces a cycle in the BUILD dep graph.
		// Regular transpiler passes must use the {@code
		// TranspilationPassses.maybeMarkFeaturesAsTranspiledAway}.
		NodeUtil_RemoveFeatureFromAllScripts(root, parser.FeatureSet_Feature_DYNAMIC_IMPORT, r.compiler)
		if r.requiresAliasing && r.AliasIsValid() {
			NodeTraversal_Traverse(r.compiler, externs, NewRewriteDynamicImports_AliasInjectingTraversal(r))
		}
	}
}

// ShouldTraverse skips scripts that do not contain dynamic imports.
func (r *RewriteDynamicImports) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		fs := NodeUtil_GetFeatureSetOfScript(n)
		return fs == nil || fs.ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	}
	return true
}

// Java: RewriteDynamicImports.visit(NodeTraversal, Node, Node)
func (r *RewriteDynamicImports) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.GetToken() != rhino.Token_DYNAMIC_IMPORT {
		return
	}

	// If the module specifier is a string, attempt to resolve the module
	moduleMap := r.compiler.GetModuleMap()
	importSpecifier := n.GetFirstChild()
	if importSpecifier.IsStringLit() && moduleMap != nil {
		targetPath :=
			t.GetInput().
				GetPath().
				ResolveJsModule(
					importSpecifier.GetString(), processCommonJSModules_sourceFileName(n), n.GetLineno(), n.GetCharno())
		var module *Module
		if targetPath != nil {
			module = moduleMap.GetModuleByPath(targetPath)
		}

		var targetModuleNS *Var
		if module != nil {
			targetModuleVarName := ModuleRenaming_GlobalizedModuleName_CreateWithModuleMetadata(module.Metadata(), nil, nil).AliasName().Join()
			targetModuleNS = t.GetScope().GetVar(targetModuleVarName)
		}

		if targetModuleNS != nil {
			targetChunk := targetModuleNS.GetInput().GetChunk()
			// If the target module is bundled into the same output chunk, replace the import statement
			// with a promise that resolves to the module namespace.
			// No further rewriting occurs for this case.
			if t.GetChunk() == targetChunk {
				r.ReplaceDynamicImportWithPromise(t, n, targetModuleNS)
				return
			} else {
				// The target output chunk is recognized and different from the current chunk.
				// Retarget the import specifier path to the output chunk path and rewrite
				// the import to reference the rewritten global module namespace variable.
				r.RetargetImportSpecifier(t, n, targetChunk)
				if NodeUtil_IsExpressionResultUsed(n) {
					r.AddChainedThen(t, n, targetModuleNS)
				}
			}
		}
	}
	if r.AliasIsValid() {
		r.AliasDynamicImport(t, n)
	} else if r.alias != nil {
		t.Report(n, RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS)
	} else if r.requiresAliasing {
		t.Report(n, RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED)
	}
}

// ReplaceDynamicImportWithPromise replaces a dynamic import expression with
// a promise resolving to the rewritten module namespace.
//
// Before
//
//	import('./foo.js');
//
// After
//
//	Promise.resolve(module$foo);
//
// Java: RewriteDynamicImports.replaceDynamicImportWithPromise(NodeTraversal, Node, Var)
func (r *RewriteDynamicImports) ReplaceDynamicImportWithPromise(
	t *NodeTraversal, dynamicImport rhino.Node, targetModuleNS *Var) {
	registry := r.compiler.GetTypeRegistry()
	promiseDotResolve :=
		r.astFactory.CreateQName(r.compiler.GetTranspilationNamespace(), "Promise.resolve")
	promiseResolveCall := r.astFactory.CreateCall(promiseDotResolve, AstFactory_Type_(dynamicImport))
	isExpressionUsed := NodeUtil_IsExpressionResultUsed(dynamicImport)
	if isExpressionUsed {
		var moduleNamespaceType rhino.JSType
		if dynamicImport.GetJSType() != nil && dynamicImport.GetJSType().IsTemplatizedType() {
			templatizedType := dynamicImport.GetJSType().(*rhino.TemplatizedType)
			moduleNamespaceType = templatizedType.GetTemplateTypes()[0]
		} else {
			moduleNamespaceType = registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		promiseDotResolve.SetJSType(
			registry.CreateFunctionType(
				registry.CreateTemplatizedType(
					registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE), moduleNamespaceType),
				registry.GetNativeType(rhino.JSTypeNative_ALL_TYPE)))
		promiseResolveCall.AddChildToBack(r.CreateModuleNamespaceNode(targetModuleNS))
		promiseResolveCall.CopyTypeFrom(dynamicImport)
	} else {
		promiseDotResolve.SetJSType(
			registry.CreateFunctionType(
				registry.CreateTemplatizedType(
					registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE),
					registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE)),
				registry.GetNativeType(rhino.JSTypeNative_ALL_TYPE)))
		promiseResolveCall.SetJSType(
			registry.CreateTemplatizedType(
				registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE),
				registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE)))
	}
	promiseResolveCall.SrcrefTree(dynamicImport)
	parent := dynamicImport.GetParent()
	dynamicImport.ReplaceWith(promiseResolveCall)
	t.ReportCodeChangeWithNode(parent)
	r.dynamicImportsRemoved = true
}

// Java: RewriteDynamicImports.getChunkFileName(JSChunk)
func RewriteDynamicImports_GetChunkFileName(chunk *JSChunk) string {
	return chunk.GetName() + ".js"
}

// RetargetImportSpecifier rewrites the dynamic import specifier to the
// output chunk path
//
// Before
//
//	import('./foo.js');
//
// After
//
//	import('./chunk0.js');
//
// Java: RewriteDynamicImports.retargetImportSpecifier(NodeTraversal, Node, JSChunk)
func (r *RewriteDynamicImports) RetargetImportSpecifier(t *NodeTraversal, dynamicImport rhino.Node, targetChunk *JSChunk) {
	importingChunkFilename := RewriteDynamicImports_GetChunkFileName(t.GetInput().GetChunk())
	targetChunkFilename := RewriteDynamicImports_GetChunkFileName(targetChunk)
	retargetedSpecifier, ok := tryRelativePathFrom(importingChunkFilename, targetChunkFilename)
	if !ok {
		r.compiler.Report(
			JSError_MakeWithNode(
				dynamicImport,
				ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH,
				importingChunkFilename,
				targetChunkFilename))
		return
	}
	originalSpecifierNode := dynamicImport.GetFirstChild()
	newSpecifier := r.astFactory.CreateString(retargetedSpecifier).Srcref(originalSpecifierNode)
	originalSpecifierNode.ReplaceWith(newSpecifier)
	t.ReportCodeChangeWithNode(newSpecifier)
}

func tryRelativePathFrom(fromUriPath, toUriPath string) (string, bool) {
	fromPath := jutil.Path_Of(fromUriPath)
	toPath := jutil.Path_Of(toUriPath)
	fromFolder := fromPath.GetParent()

	if fromFolder == nil && toPath.GetParent() == nil {
		return "./" + toUriPath, true
	} else if fromFolder == nil &&
		(strings.HasPrefix(toUriPath, ".") ||
			strings.HasPrefix(toPath.String(), "/") ||
			strings.HasPrefix(toPath.String(), "\\")) {
		return toUriPath, true
	} else if fromFolder == nil {
		return "", false
	}

	if !canRelativizePath(fromFolder, toPath) {
		return "", false
	}
	calculatedPath := fromFolder.Relativize(toPath).String()
	if strings.HasPrefix(calculatedPath, ".") || strings.HasPrefix(calculatedPath, "/") {
		return calculatedPath, true
	}
	return "./" + calculatedPath, true
}

func canRelativizePath(base, target *jutil.Path) bool {
	if target.IsAbsolute() != base.IsAbsolute() {
		return false
	}
	if base.Equals(target) || base.String() == "" {
		return true
	}
	nb, nt := base.Normalize(), target.Normalize()
	bCount, tCount := nb.GetNameCount(), nt.GetNameCount()
	common := 0
	for common < bCount && common < tCount && nb.GetName(common).Equals(nt.GetName(common)) {
		common++
	}
	for i := common; i < bCount; i++ {
		if nb.GetName(i).String() == ".." {
			return false
		}
	}
	return true
}

// AddChainedThen: since a dynamic import expression is a promise resolving
// to the namespace export type, add a ".then()" call after it and resolve
// to the rewritten module namespace.
//
// Before
//
//	import('./foo.js');
//
// After
//
//	import('./foo.js').then(() => module$foo);
//
// Java: RewriteDynamicImports.addChainedThen(NodeTraversal, Node, Var)
func (r *RewriteDynamicImports) AddChainedThen(t *NodeTraversal, dynamicImport rhino.Node, targetModuleNs *Var) {
	registry := r.compiler.GetTypeRegistry()
	importParent := dynamicImport.GetParent()
	placeholder := rhino.IR_Empty()
	dynamicImport.ReplaceWith(placeholder)
	moduleNamespaceNode := r.CreateModuleNamespaceNode(targetModuleNs)
	callbackFn := r.astFactory.CreateZeroArgArrowFunctionForExpression(moduleNamespaceNode)
	callbackFn.SetJSType(
		registry.CreateFunctionType(
			registry.CreateTemplatizedType(
				registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE),
				registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)),
			registry.CreateFunctionType(registry.GetNativeType(rhino.JSTypeNative_ALL_TYPE))))
	thenArgument := callbackFn
	if r.shouldWrapDynamicImportCallbacks {
		wrappingFunction :=
			r.astFactory.CreateName(
				ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN,
				AstFactory_Type_OfJSType(registry.CreateFunctionType(callbackFn.GetJSType(), callbackFn.GetJSType())))
		thenArgument = r.astFactory.CreateCall(wrappingFunction, AstFactory_Type_(callbackFn), callbackFn)
		r.wrappedDynamicImportCallback = true
	}
	importThenCall :=
		r.astFactory.CreateCall(
			r.astFactory.CreateGetPropWithUnknownType(dynamicImport, "then"),
			AstFactory_Type_(dynamicImport),
			thenArgument)
	importThenCall.SrcrefTreeIfMissing(dynamicImport)
	if dynamicImport.GetJSType() != nil {
		importThenCall.CopyTypeFrom(dynamicImport)
	}
	placeholder.ReplaceWith(importThenCall)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_ARROW_FUNCTIONS, r.compiler)
	r.compiler.ReportChangeToChangeScope(callbackFn)
	r.compiler.ReportChangeToEnclosingScope(importParent)
}

// AliasDynamicImport replaces a dynamic import expression with a function
// call to the specified alias.
//
// Before
//
//	import('./foo.js');
//
// After
//
//	aliasedName('./foo.js');
//
// Java: RewriteDynamicImports.aliasDynamicImport(NodeTraversal, Node)
func (r *RewriteDynamicImports) AliasDynamicImport(t *NodeTraversal, dynamicImport rhino.Node) {
	jutil.CheckNotNull(r.alias)
	aliasNode := r.astFactory.CreateQNameWithUnknownType(*r.alias)
	aliasNode.SetOriginalName(&rewriteDynamicImportsOriginalName)
	moduleSpecifier := dynamicImport.RemoveFirstChild()
	importAliasCall :=
		r.astFactory.
			CreateCall(aliasNode, AstFactory_Type_(dynamicImport), moduleSpecifier).
			SrcrefTreeIfMissing(dynamicImport)
	if dynamicImport.GetJSType() != nil {
		importAliasCall.CopyTypeFrom(dynamicImport)
	}
	dynamicImport.ReplaceWith(importAliasCall)
	t.ReportCodeChangeWithNode(importAliasCall)
	r.dynamicImportsRemoved = true
}

// CreateModuleNamespaceNode returns, for a given module, a reference to the
// module namespace export
//
// Java: RewriteDynamicImports.createModuleNamespaceNode(Var)
func (r *RewriteDynamicImports) CreateModuleNamespaceNode(moduleVar *Var) rhino.Node {
	moduleVarNode := moduleVar.GetNode()
	moduleNamespace := moduleVarNode.CloneNode()
	moduleNamespace.CopyTypeFrom(moduleVarNode)
	return moduleNamespace
}

// InjectWrappingFunctionExtern returns, for a given module, a reference to
// the module namespace export
//
// Java: RewriteDynamicImports.injectWrappingFunctionExtern()
func (r *RewriteDynamicImports) InjectWrappingFunctionExtern() {
	registry := r.compiler.GetTypeRegistry()
	templateT := registry.CreateTemplateType("T")
	wrappingFunctionDefinition :=
		r.astFactory.CreateFunction(
			ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN,
			r.astFactory.CreateParamList("importCallback"),
			r.astFactory.CreateBlock(),
			AstFactory_Type_OfJSType(registry.CreateFunctionType(templateT, templateT)))

	externsRoot := r.compiler.GetSynthesizedExternsInput().GetAstRoot(r.compiler)
	wrappingFunctionDefinition.SrcrefTree(externsRoot)
	externsRoot.AddChildToBack(wrappingFunctionDefinition)
	r.compiler.ReportChangeToEnclosingScope(wrappingFunctionDefinition)
}

// RewriteDynamicImports_AliasInjectingTraversal is a shallow traversal class
// for the externs to inject the alias.
//
// Later passes require the names to be defined. For simple name aliases
// first check for an existing definition. Qualified name aliases check for
// an existing definition of the root name only.
type RewriteDynamicImports_AliasInjectingTraversal struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	outer *RewriteDynamicImports // the enclosing instance
}

// Java: RewriteDynamicImports.AliasInjectingTraversal()
func NewRewriteDynamicImports_AliasInjectingTraversal(outer *RewriteDynamicImports) *RewriteDynamicImports_AliasInjectingTraversal {
	return &RewriteDynamicImports_AliasInjectingTraversal{outer: outer}
}

// Java: RewriteDynamicImports.AliasInjectingTraversal.shouldTraverse(NodeTraversal, Node, Node)
func (a *RewriteDynamicImports_AliasInjectingTraversal) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	r := a.outer
	jutil.CheckState(r.AliasIsValid())
	alias := *r.alias

	// If the extern is a simple name and already declared, there is nothing to do
	if NodeUtil_IsValidSimpleName(alias) && t.GetScope().GetVar(alias) != nil {
		return false
	}

	// Create the type of the import function
	registry := r.compiler.GetTypeRegistry()
	unknownType := registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	promiseType := registry.GetNativeType(rhino.JSTypeNative_PROMISE_TYPE)
	promiseTemplatizedType :=
		registry.CreateTemplatizedType(promiseType.ToObjectType(), unknownType)
	stringType := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)

	externsRoot := r.compiler.GetSynthesizedExternsInput().GetAstRoot(r.compiler)
	// "import" is not a valid JS name and will not parse. However we can manually create it
	// and inject it into the externs.
	if alias == "import" || NodeUtil_IsValidSimpleName(alias) {
		aliasNode :=
			r.astFactory.CreateFunction(
				alias,
				r.astFactory.CreateParamList("specifier"),
				r.astFactory.CreateBlock(),
				AstFactory_Type_OfJSType(registry.CreateFunctionType(promiseTemplatizedType, stringType)))
		aliasNode.SrcrefTree(externsRoot)
		externsRoot.AddChildToBack(aliasNode)
	} else {
		aliasRootName := NodeUtil_GetRootOfQualifiedNameString(alias)
		aliasVar := t.GetScope().GetVar(aliasRootName)
		// If the namespace root is defined, just assume the whole name is properly defined.
		// This may need revisited in the future, but checking if the full qualified name
		// is properly defined is difficult here. As long as the root of the qualified name is
		// defined, other passes seem be ok even with missing properties.
		if aliasVar != nil {
			return false
		}
		aliasRootNode :=
			r.astFactory.CreateSingleVarNameDeclarationWithValue(aliasRootName, r.astFactory.CreateObjectLit())
		aliasRootNode.SrcrefTree(externsRoot)
		externsRoot.AddChildToBack(aliasRootNode)

		// Define the rest of the parts of the name
		aliasNameParts := strings.Split(alias, ".") // Splitter.on(".").split(alias).iterator()
		qName := aliasRootNode.GetFirstChild().CloneNode()
		next := 1 // skip over root name
		for next < len(aliasNameParts) {
			qName = r.astFactory.CreateGetPropWithUnknownType(qName.CloneTree(), aliasNameParts[next])
			next++
			var assignedValue rhino.Node
			if next >= len(aliasNameParts) {
				assignedValue =
					r.astFactory.CreateFunction(
						"",
						r.astFactory.CreateParamList("specifier"),
						r.astFactory.CreateBlock(),
						AstFactory_Type_OfJSType(registry.CreateFunctionType(promiseTemplatizedType, stringType)))
			} else {
				assignedValue = r.astFactory.CreateObjectLit()
			}
			expr := r.astFactory.ExprResult(r.astFactory.CreateAssign(qName, assignedValue))
			expr.SrcrefTree(externsRoot)
			externsRoot.AddChildToBack(expr)
		}
	}
	r.compiler.ReportChangeToChangeScope(externsRoot)
	return false
}

// Es6RewriteScriptsToModules rewrites a script which was imported as a
// module into an ES6 module.
type Es6RewriteScriptsToModules struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	compiler AbstractCompiler
}

// NewEs6RewriteScriptsToModules creates a new Es6RewriteModules instance
// which can be used to rewrite ES6 modules to a concatenable form.
//
// Java: Es6RewriteScriptsToModules(AbstractCompiler)
func NewEs6RewriteScriptsToModules(compiler AbstractCompiler) *Es6RewriteScriptsToModules {
	e := &Es6RewriteScriptsToModules{}
	e.compiler = compiler
	return e
}

// ForceToEs6Module forces rewriting of a script into an ES6 module, such as
// for imported files that contain no "import" or "export" statements.
//
// Java: Es6RewriteScriptsToModules.forceToEs6Module(Node)
func (e *Es6RewriteScriptsToModules) ForceToEs6Module(root rhino.Node) {
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		return
	}
	moduleNode := rhino.NewNode(rhino.Token_MODULE_BODY).Srcref(root)
	moduleNode.AddChildrenToBack(root.RemoveChildren())
	root.AddChildToBack(moduleNode)
	e.compiler.ReportChangeToChangeScope(root)
}

// Java: Es6RewriteScriptsToModules.process(Node, Node)
func (e *Es6RewriteScriptsToModules) Process(externs rhino.Node, root rhino.Node) {
	for file := root.GetFirstChild(); !file.IsNil(); file = file.GetNext() {
		jutil.CheckArgument(file.IsScript())
		NodeTraversal_Traverse(e.compiler, file, e)
	}
}

// Java: Es6RewriteScriptsToModules.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteScriptsToModules) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		moduleType := e.compiler.GetInput(n.GetInputId()).GetJsModuleType()

		if moduleType == CompilerInput_ModuleType_IMPORTED_SCRIPT {
			e.ForceToEs6Module(n)
		}
	}

	return false
}

// ForbidDynamicImportUsage warns at any usage of Dynamic Import expressions
// that they are unable to be transpiled.
type ForbidDynamicImportUsage struct {
	compiler AbstractCompiler
}

func (*ForbidDynamicImportUsage) noVisitLeaves() {}

func (*ForbidDynamicImportUsage) visitTokens() *ast.TokenSet {
	return &rewriteDynamicImportsVisitTokens
}

// Java: ForbidDynamicImportUsage.DYNAMIC_IMPORT_USAGE
var ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE = DiagnosticType_Error(
	"JSC_DYNAMIC_IMPORT_USAGE", "Dynamic import expressions cannot be transpiled.")

// Java: ForbidDynamicImportUsage(AbstractCompiler)
func NewForbidDynamicImportUsage(compiler AbstractCompiler) *ForbidDynamicImportUsage {
	f := &ForbidDynamicImportUsage{}
	f.compiler = compiler
	return f
}

// Java: ForbidDynamicImportUsage.process(Node, Node)
func (f *ForbidDynamicImportUsage) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(f.compiler, root, f)
	f.compiler.MarkFeatureNotAllowed(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
}

// Java: ForbidDynamicImportUsage.shouldTraverse(NodeTraversal, Node, Node)
func (f *ForbidDynamicImportUsage) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		// We only want to traverse scripts that contain dynamic imports, so we can generate error
		// messages that point to them.
		return f.GetFeatureSetOfScript(n).ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	} else {
		return true
	}
}

// Java: ForbidDynamicImportUsage.getFeatureSetOfScript(Node)
func (f *ForbidDynamicImportUsage) GetFeatureSetOfScript(n rhino.Node) *parser.FeatureSet {
	featureSet := NodeUtil_GetFeatureSetOfScript(n)
	if featureSet == nil {
		return parser.FeatureSet_BARE_MINIMUM
	}
	return featureSet
}

// Java: ForbidDynamicImportUsage.visit(NodeTraversal, Node, Node)
func (f *ForbidDynamicImportUsage) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_DYNAMIC_IMPORT:
		t.Report(n, ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE)
	default:
	}
}
