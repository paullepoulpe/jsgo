// Port of com/google/javascript/jscomp/CompilerInput.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CompilerInput is a class for the internal representation of an input to
// the compiler. Wraps a JsAst and maintain state such as the JSChunk for the
// input and whether the input is an extern. Also calculates provided and
// required types.
type CompilerInput struct {
	// Info about where the file lives.
	chunk *JSChunk
	// Go-only: performance — the input's index in compilerInputTable
	// (abstract_var.go), 0 until a variable refers to it.
	varTableIndex atomic.Uint32
	// Go-only: the InputId is its own allocation, never a field of the
	// CompilerInput: the SCRIPT node keeps it as a property, and an interior
	// pointer would keep the CompilerInput (and through it the Compiler)
	// alive from rhino's node tables for as long as the node's cold data
	// is registered there (see rhino/node_tables.go).
	id         *rhino.InputId
	sourceFile *SourceFile

	// The lazily constructed AST.
	ast       *CompilerInput_JsAst
	inlineAst CompilerInput_JsAst

	// DependencyInfo to delegate to.
	dependencyInfo  DependencyInfo
	extraRequires   []*DependencyInfo_Require
	extraProvides   []string
	orderedRequires []*DependencyInfo_Require
	dynamicRequires []string
	// Modules imported by goog.requireDynamic()
	requireDynamicImports      []string
	hasFullParseDependencyInfo bool
	jsModuleType               CompilerInput_ModuleType

	// If this input file has a MODULE_BODY (i.e. it is a goog.module() or
	// ES module), TypedScopeCreator will store the scope created for that
	// Node here.
	//
	// DefaultPassConfig is responsible for erasing this value when later use
	// of it is not needed.
	typedScope *TypedScope

	// An AbstractCompiler for doing parsing.
	compiler   AbstractCompiler
	modulePath *ModuleLoader_ModulePath
}

var _ DependencyInfo = (*CompilerInput)(nil)

// Java: CompilerInput.serialVersionUID
const compilerInput_serialVersionUID int64 = 2

// Java: CompilerInput.setTypedScope(TypedScope)
func (c *CompilerInput) SetTypedScope(typedScope *TypedScope) {
	c.typedScope = typedScope
}

// Java: CompilerInput.getTypedScope()
func (c *CompilerInput) GetTypedScope() *TypedScope {
	return c.typedScope
}

// TODO(tjgq): Whether a CompilerInput is an externs file is determined by
// the `isExtern` constructor argument and the `setIsExtern` method. Both are
// necessary because, while externs files passed under the --externs flag are
// not required to contain an @externs annotation, externs files passed under
// any other flag must have one, and the presence of the latter can only be
// known once the file has been parsed. To add to the confusion, note that
// CompilerInput doesn't actually store the extern bit itself, but instead
// mutates the SourceFile associated with the AST node. Once (when?) we
// enforce that extern files always contain an @externs annotation, we can
// store the extern bit in the AST node, and make SourceFile immutable.

// Java: CompilerInput(SourceFile, InputId)
func NewCompilerInputWithInputId(file *SourceFile, inputId *rhino.InputId) *CompilerInput {
	return NewCompilerInputWithInputIdAndIsExtern(file, inputId, false)
}

// Deprecated: the inputId is read from the SourceFile. Use
// NewCompilerInputWithIsExtern(file, isExtern)
//
// Java: CompilerInput(SourceFile, String, boolean)
func NewCompilerInputWithStringIdAndIsExtern(sourceFile *SourceFile, inputId string, isExtern bool) *CompilerInput {
	c := &CompilerInput{}
	return initCompilerInputWithInputId(c, sourceFile, rhino.NewInputId(inputId), isExtern)
}

var (
	compilerInputEmptyRequires = []*DependencyInfo_Require{}
	compilerInputEmptyStrings  = []string{}
)

func initCompilerInputWithInputId(c *CompilerInput, sourceFile *SourceFile, inputId *rhino.InputId, isExtern bool) *CompilerInput {
	c.extraRequires = compilerInputEmptyRequires
	c.extraProvides = compilerInputEmptyStrings
	c.orderedRequires = compilerInputEmptyRequires
	c.dynamicRequires = compilerInputEmptyStrings
	c.requireDynamicImports = compilerInputEmptyStrings
	c.jsModuleType = CompilerInput_ModuleType_NONE
	c.inlineAst.outer = c
	c.ast = &c.inlineAst
	c.sourceFile = jutil.CheckNotNull(sourceFile)
	c.id = jutil.CheckNotNull(inputId)
	if isExtern {
		c.SetIsExtern()
	}
	return c
}

// Deprecated: the inputId is read from the SourceFile. Use
// NewCompilerInputWithIsExtern(file, isExtern)
//
// Java: CompilerInput(SourceFile, InputId, boolean)
func NewCompilerInputWithInputIdAndIsExtern(sourceFile *SourceFile, inputId *rhino.InputId, isExtern bool) *CompilerInput {
	return initCompilerInputWithInputId(&CompilerInput{}, sourceFile, inputId, isExtern)
}

// Java: CompilerInput(SourceFile)
func NewCompilerInput(file *SourceFile) *CompilerInput {
	return NewCompilerInputWithIsExtern(file, false)
}

// Java: CompilerInput(SourceFile, boolean)
func NewCompilerInputWithIsExtern(file *SourceFile, isExtern bool) *CompilerInput {
	jutil.CheckNotNull(file)
	return initCompilerInputWithInputId(&CompilerInput{}, file, rhino.NewInputId(file.GetName()), isExtern)
}

// GetInputId returns a name for this input. Must be unique across all
// inputs.
//
// Java: CompilerInput.getInputId()
func (c *CompilerInput) GetInputId() *rhino.InputId {
	return c.id
}

// GetName returns a name for this input. Must be unique across all inputs.
//
// Java: CompilerInput.getName()
func (c *CompilerInput) GetName() string {
	return c.id.GetIdName()
}

// GetPathRelativeToClosureBase gets the path relative to closure-base, if
// one is available.
//
// Java: CompilerInput.getPathRelativeToClosureBase()
func (c *CompilerInput) GetPathRelativeToClosureBase() string {
	// TODO(nicksantos): Implement me.
	panic(jutil.NewUnsupportedOperationException())
}

// Java: CompilerInput.getAstRoot(AbstractCompiler)
func (c *CompilerInput) GetAstRoot(compiler AbstractCompiler) rhino.Node {
	root := rhino.CheckNotNull(c.ast.GetAstRoot(compiler))
	jutil.CheckState(root.IsScript())
	jutil.CheckNotNull(root.GetInputId())
	return root
}

// Java: CompilerInput.clearAst()
func (c *CompilerInput) ClearAst() {
	c.ast.ClearAst()
}

// Java: CompilerInput.initShadowAst(Node)
func (c *CompilerInput) InitShadowAst(shadowScript rhino.Node) {
	jutil.CheckArgumentMsg1(shadowScript.IsScript(), "%s", shadowScript)
	jutil.CheckStateMsg1(shadowScript.GetIsInClosureUnawareSubtree(), "%s", shadowScript)
	c.ast.root = shadowScript
}

// Java: CompilerInput.getSourceFile()
func (c *CompilerInput) GetSourceFile() *SourceFile {
	return c.sourceFile
}

// SetCompiler sets an abstract compiler for doing parsing.
//
// Java: CompilerInput.setCompiler(AbstractCompiler)
func (c *CompilerInput) SetCompiler(compiler AbstractCompiler) {
	c.compiler = compiler
}

// GetRequires gets a list of types depended on by this input.
//
// Java: CompilerInput.getRequires()
func (c *CompilerInput) GetRequires() []*DependencyInfo_Require {
	if c.hasFullParseDependencyInfo {
		return append([]*DependencyInfo_Require{}, c.orderedRequires...)
	}

	return c.GetDependencyInfo().GetRequires()
}

// Java: CompilerInput.getTypeRequires()
func (c *CompilerInput) GetTypeRequires() []string {
	return c.GetDependencyInfo().GetTypeRequires()
}

// GetKnownRequires gets a list of namespaces and paths depended on by this
// input, but does not attempt to regenerate the dependency information.
// Typically this occurs from module rewriting.
//
// Java: CompilerInput.getKnownRequires()
func (c *CompilerInput) GetKnownRequires() []*DependencyInfo_Require {
	var first []*DependencyInfo_Require
	if c.dependencyInfo != nil {
		first = c.dependencyInfo.GetRequires()
	} else {
		first = []*DependencyInfo_Require{}
	}
	return compilerInput_concat(first, c.extraRequires)
}

// Java: CompilerInput.getKnownRequiredSymbols()
func (c *CompilerInput) GetKnownRequiredSymbols() []string {
	return DependencyInfo_Require_AsSymbolList(c.GetKnownRequires())
}

// GetProvides gets a list of types provided by this input.
//
// Java: CompilerInput.getProvides()
func (c *CompilerInput) GetProvides() []string {
	return c.GetDependencyInfo().GetProvides()
}

// Java: CompilerInput.getHasExternsAnnotation()
func (c *CompilerInput) GetHasExternsAnnotation() bool {
	return c.GetDependencyInfo().GetHasExternsAnnotation()
}

// Java: CompilerInput.getHasNoCompileAnnotation()
func (c *CompilerInput) GetHasNoCompileAnnotation() bool {
	return c.GetDependencyInfo().GetHasNoCompileAnnotation()
}

// GetKnownProvides gets a list of types provided, but does not attempt to
// regenerate the dependency information. Typically this occurs from module
// rewriting.
//
// Java: CompilerInput.getKnownProvides()
func (c *CompilerInput) GetKnownProvides() []string {
	var first []string
	if c.dependencyInfo != nil {
		first = c.dependencyInfo.GetProvides()
	} else {
		first = []string{}
	}
	return compilerInput_concat(first, c.extraProvides)
}

// AddProvide registers a type that this input defines. Includes both
// explicitly declared namespaces via goog.provide and goog.module calls as
// well as implicit namespaces provided by module rewriting.
//
// Java: CompilerInput.addProvide(String)
func (c *CompilerInput) AddProvide(provide string) {
	c.extraProvides = append(c.extraProvides, provide)
}

// AddOrderedRequire registers a type that this input depends on in the
// order seen in the file.
//
// Java: CompilerInput.addOrderedRequire(Require)
func (c *CompilerInput) AddOrderedRequire(require *DependencyInfo_Require) bool {
	for _, r := range c.orderedRequires {
		if jutil.ObjectsEquals(r, require) {
			return false
		}
	}
	c.orderedRequires = append(c.orderedRequires, require)
	return true
}

// GetDynamicRequires returns the types that this input dynamically depends
// on in the order seen in the file. The returned types were loaded
// dynamically so while they are part of the dependency graph, they do not
// need sorted before this input.
//
// Java: CompilerInput.getDynamicRequires()
func (c *CompilerInput) GetDynamicRequires() []string {
	return append([]string{}, c.dynamicRequires...)
}

// AddDynamicRequire registers a type that this input depends on in the
// order seen in the file. The type was loaded dynamically so while it is
// part of the dependency graph, it does not need sorted before this input.
//
// Java: CompilerInput.addDynamicRequire(String)
func (c *CompilerInput) AddDynamicRequire(require string) bool {
	for _, r := range c.dynamicRequires {
		if r == require {
			return false
		}
	}
	c.dynamicRequires = append(c.dynamicRequires, require)
	return true
}

// GetRequireDynamicImports returns the types that this input dynamically
// imports by goog.requireDynamic() in the order seen in the file. The
// returned types were loaded dynamically so while they are part of the
// dependency graph, they do not need sorted before this input.
//
// Java: CompilerInput.getRequireDynamicImports()
func (c *CompilerInput) GetRequireDynamicImports() []string {
	return append([]string{}, c.requireDynamicImports...)
}

// AddRequireDynamicImports registers a type that this input depends on in
// the order seen in the file. The type was loaded by goog.requireDynamic() so
// while it is part of the dependency graph, it does not need sorted before
// this input.
//
// Java: CompilerInput.addRequireDynamicImports(String)
func (c *CompilerInput) AddRequireDynamicImports(require string) {
	for _, r := range c.requireDynamicImports {
		if r == require {
			return
		}
	}
	c.requireDynamicImports = append(c.requireDynamicImports, require)
}

// Java: CompilerInput.setHasFullParseDependencyInfo(boolean)
func (c *CompilerInput) SetHasFullParseDependencyInfo(hasFullParseDependencyInfo bool) {
	c.hasFullParseDependencyInfo = hasFullParseDependencyInfo
}

// Java: CompilerInput.getJsModuleType()
func (c *CompilerInput) GetJsModuleType() CompilerInput_ModuleType {
	return c.jsModuleType
}

// Java: CompilerInput.setJsModuleType(ModuleType)
func (c *CompilerInput) SetJsModuleType(moduleType CompilerInput_ModuleType) {
	c.jsModuleType = moduleType
}

// AddRequire registers a type that this input depends on.
//
// Java: CompilerInput.addRequire(Require)
func (c *CompilerInput) AddRequire(require *DependencyInfo_Require) {
	c.extraRequires = append(c.extraRequires, require)
}

// GetDependencyInfo returns the DependencyInfo object, generating it lazily
// if necessary.
//
// Java: CompilerInput.getDependencyInfo()
func (c *CompilerInput) GetDependencyInfo() DependencyInfo {
	if c.dependencyInfo == nil {
		c.dependencyInfo = c.generateDependencyInfo()
	}
	if len(c.extraRequires) != 0 || len(c.extraProvides) != 0 {
		c.dependencyInfo = SimpleDependencyInfo_Builder_(c.GetName(), c.GetName()).
			SetProvides(compilerInput_concat(c.dependencyInfo.GetProvides(), c.extraProvides)).
			SetRequires(compilerInput_concat(c.dependencyInfo.GetRequires(), c.extraRequires)).
			SetTypeRequires(c.dependencyInfo.GetTypeRequires()).
			SetLoadFlags(c.dependencyInfo.GetLoadFlags()).
			SetHasExternsAnnotation(c.dependencyInfo.GetHasExternsAnnotation()).
			SetHasNoCompileAnnotation(c.dependencyInfo.GetHasNoCompileAnnotation()).
			Build()
		c.extraRequires = c.extraRequires[:0]
		c.extraProvides = c.extraProvides[:0]
	}
	return c.dependencyInfo
}

// generateDependencyInfo generates the DependencyInfo by scanning and/or
// parsing the file. This is called lazily by getDependencyInfo, and does not
// take into account any extra requires/provides added by AddRequire or
// AddProvide.
//
// Java: CompilerInput.generateDependencyInfo()
func (c *CompilerInput) generateDependencyInfo() DependencyInfo {
	jutil.CheckNotNullMsg1(c.compiler, "Expected setCompiler to be called first: %s", c)
	jutil.CheckNotNullMsg1(
		c.compiler.GetErrorManager(), "Expected compiler to call an error manager: %s", c)

	// If the code is a JsAst, then it was originally JS code, and is compatible with the
	// regex-based parsing of JsFileRegexParser.
	if JsFileRegexParser_IsSupported() && c.compiler.PreferRegexParser() {
		// Look at the source code.
		// Note: it's OK to use getName() instead of
		// getPathRelativeToClosureBase() here because we're not using
		// this to generate deps files. (We're only using it for
		// symbol dependencies.)
		var result DependencyInfo
		jutil.Try(func() {
			info := NewJsFileRegexParser(c.compiler.GetErrorManager()).
				SetModuleLoader(c.compiler.GetModuleLoader()).
				SetIncludeGoogBase(true).
				ParseFile(c.GetName(), c.GetName(), c.GetCode())
			result = NewLazyParsedDependencyInfo(info, c, c.compiler)
		}, func(e jutil.IIOException) {
			c.compiler.GetErrorManager().Report(
				CheckLevel_ERROR.Ptr(),
				JSError_Make(AbstractCompiler_READ_ERROR, c.GetName(), e.GetMessage()))
			result = SimpleDependencyInfo_EMPTY
		})
		return result
	} else {
		// Otherwise, just look at the AST.
		root := c.GetAstRoot(c.compiler)
		if root.IsNil() {
			return SimpleDependencyInfo_EMPTY
		}

		var providesBuf [4]string
		var requiresBuf [16]*DependencyInfo_Require
		finder := CompilerInput_DepsFinder{
			provides:   providesBuf[:0],
			requires:   requiresBuf[:0],
			modulePath: c.GetPath(),
		}

		finder.VisitTree(root)
		info := root.GetJSDocInfo()

		// TODO(nicksantos|user): This caching behavior is a bit
		// odd, and only works if you assume the exact call flow that
		// clients are currently using.  In that flow, they call
		// getProvides(), then remove the goog.provide calls from the
		// AST, and then call getProvides() again.
		//
		// This won't work for any other call flow, or any sort of incremental
		// compilation scheme. The API needs to be fixed so callers aren't
		// doing weird things like this, and then we should get rid of the
		// multiple-scan strategy.
		return newSimpleDependencyInfoDirect(
			"",
			"",
			finder.provides,
			finder.requires,
			finder.typeRequires,
			simpleDependencyInfo_flagsForLangModule(finder.langFlag, finder.moduleFlag),
			info != nil && info.IsExterns(),
			info != nil && info.IsNoCompile(),
		)
	}
}

// CompilerInput_DepsFinder is CompilerInput.DepsFinder.
type CompilerInput_DepsFinder struct {
	langFlag     string
	moduleFlag   string
	provides     []string
	requires     []*DependencyInfo_Require
	typeRequires []string
	modulePath   *ModuleLoader_ModulePath
}

// Java: CompilerInput.DepsFinder(ModulePath)
func newCompilerInput_DepsFinder(modulePath *ModuleLoader_ModulePath) *CompilerInput_DepsFinder {
	return &CompilerInput_DepsFinder{
		provides:     []string{},
		requires:     []*DependencyInfo_Require{},
		typeRequires: []string{},
		modulePath:   modulePath,
	}
}

// Java: CompilerInput.DepsFinder.visitTree(Node)
func (f *CompilerInput_DepsFinder) VisitTree(n rhino.Node) {
	f.VisitSubtree(n, rhino.NilNode)
	jutil.CheckArgument(n.IsScript())
	var features *parser.FeatureSet
	if p := n.GetProp(rhino.Node_FEATURE_SET); p != nil {
		features = p.(*parser.FeatureSet)
	}
	if features != nil {
		// Only add the "lang" load flag if it's not the default (es3), so that
		// legacy deps files will remain unchanged (i.e. load flags omitted).
		version := features.Version()
		if version != "es3" {
			f.langFlag = version
		}
	}
}

func (f *CompilerInput_DepsFinder) containsRequire(r *DependencyInfo_Require) bool {
	for _, x := range f.requires {
		if jutil.ObjectsEquals(x, r) {
			return true
		}
	}
	return false
}

// Java: CompilerInput.DepsFinder.visitSubtree(Node, Node)
func (f *CompilerInput_DepsFinder) VisitSubtree(n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		if n.HasTwoChildren() &&
			n.GetFirstChild().IsGetProp() &&
			n.GetFirstFirstChild().MatchesName("goog") {

			if !f.containsRequire(DependencyInfo_Require_BASE) {
				f.requires = append(f.requires, DependencyInfo_Require_BASE)
			}

			callee := n.GetFirstChild()
			argument := n.GetLastChild()
			switch callee.GetString() {
			case "module", "provide":
				if callee.GetString() == "module" {
					// only mark as a goog.module if this is not bundled, e.g.
					// goog.loadModule(function(exports) {"use strict";goog.module('Foo');
					if parent.IsExprResult() && parent.GetParent().IsModuleBody() {
						f.moduleFlag = "goog"
					}
				}
				// Fall-through
				if !argument.IsStringLit() {
					return
				}
				f.provides = append(f.provides, argument.GetString())
				return

			case "require":
				if !argument.IsStringLit() {
					return
				}
				f.requires = append(f.requires, DependencyInfo_Require_GoogRequireSymbol(argument.GetString()))
				return

			case "requireType":
				if !argument.IsStringLit() {
					return
				}
				f.typeRequires = append(f.typeRequires, argument.GetString())
				return

			case "loadModule":
				// Process the block of the loadModule argument if it's a function, as opposed to
				// a string.
				if argument.IsStringLit() {
					panic(jutil.NewIllegalArgumentExceptionMsg(
						"Unsupported parse of goog.loadModule with string literal"))
				}
				n = argument.GetLastChild()

			case "declareModuleId":
				if !argument.IsStringLit() {
					return
				}
				f.provides = append(f.provides, argument.GetString())

			default:
				return
			}
		}
	case rhino.Token_MODULE_BODY:
		if !parent.GetBooleanProp(rhino.Node_GOOG_MODULE) {
			f.provides = append(f.provides, f.modulePath.ToModuleName())
			f.moduleFlag = "es6"
		}
	case rhino.Token_IMPORT:
		f.VisitEs6ModuleName(n.GetLastChild(), n)
		return
	case rhino.Token_EXPORT:
		if NodeUtil_IsExportFrom(n) {
			f.VisitEs6ModuleName(n.GetLastChild(), n)
		}
		return
	case rhino.Token_EXPR_RESULT, rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR,
		rhino.Token_BLOCK, rhino.Token_NAME, rhino.Token_DESTRUCTURING_LHS:
	case rhino.Token_SCRIPT:
		jsdoc := n.GetJSDocInfo()
		if jsdoc != nil && jsdoc.IsProvideGoog() {
			f.provides = append(f.provides, "goog")
		}
	default:
		return
	}

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		f.VisitSubtree(child, n)
	}
}

// Java: CompilerInput.DepsFinder.visitEs6ModuleName(Node, Node)
func (f *CompilerInput_DepsFinder) VisitEs6ModuleName(n rhino.Node, parent rhino.Node) {
	jutil.CheckArgument(n.IsStringLit())
	jutil.CheckArgument(parent.IsExport() || parent.IsImport())

	// TODO(blickly): Move this (and the duplicated logic in JsFileRegexParser/Es6RewriteModules)
	// into ModuleLoader.
	moduleName := n.GetString()
	if strings.HasPrefix(moduleName, "goog:") {
		// cut off the "goog:" prefix
		f.requires = append(f.requires,
			DependencyInfo_Require_GoogRequireSymbol(jutil.StringSubstringFrom(moduleName, 5)))
		return
	}
	importedModule := f.modulePath.ResolveJsModule(
		moduleName, f.modulePath.String(), n.GetLineno(), n.GetCharno())

	if importedModule == nil {
		importedModule = f.modulePath.ResolveModuleAsPath(moduleName)
	}

	f.requires = append(f.requires, DependencyInfo_Require_Es6Import(importedModule.ToModuleName(), n.GetString()))
}

// GetCode returns the source code. It panics with an IOException if the
// code cannot be loaded.
//
// Java: CompilerInput.getCode()
func (c *CompilerInput) GetCode() string {
	return c.GetSourceFile().GetCode()
}

// GetChunk returns the chunk to which the input belongs.
//
// Java: CompilerInput.getChunk()
func (c *CompilerInput) GetChunk() *JSChunk {
	return c.chunk
}

// SetChunk sets the chunk to which the input belongs.
//
// Java: CompilerInput.setChunk(JSChunk)
func (c *CompilerInput) SetChunk(chunk *JSChunk) {
	// An input may only belong to one chunk.
	jutil.CheckArgument(chunk == nil || c.chunk == nil || c.chunk == chunk)
	c.chunk = chunk
}

// OverrideChunk overrides the chunk to which the input belongs.
//
// Java: CompilerInput.overrideChunk(JSChunk)
func (c *CompilerInput) OverrideChunk(chunk *JSChunk) {
	c.chunk = chunk
}

// Java: CompilerInput.isExtern()
func (c *CompilerInput) IsExtern() bool {
	return c.sourceFile.IsExtern()
}

// Java: CompilerInput.setIsExtern()
func (c *CompilerInput) SetIsExtern() {
	c.sourceFile.SetKind(rhino.StaticSourceFile_SourceKind_EXTERN)
}

// Java: CompilerInput.getLineOffset(int)
func (c *CompilerInput) GetLineOffset(lineno int) int {
	return c.sourceFile.GetLineOffset(lineno)
}

// Java: CompilerInput.toString()
func (c *CompilerInput) String() string {
	return c.GetName()
}

// Java: CompilerInput.isEs6Module()
func (c *CompilerInput) IsEs6Module() bool {
	// Instead of doing a full parse to read all load flags, just ask the delegate, which at least
	// has this much info
	return c.GetDependencyInfo().IsEs6Module()
}

// Java: CompilerInput.isGoogModule()
func (c *CompilerInput) IsGoogModule() bool {
	// Instead of doing a full parse to read all load flags, just ask the delegate, which at least
	// has this much info
	return c.GetDependencyInfo().IsGoogModule()
}

// Java: CompilerInput.getLoadFlags()
func (c *CompilerInput) GetLoadFlags() *jutil.LinkedHashMap[string, string] {
	return c.GetDependencyInfo().GetLoadFlags()
}

// Java: DependencyInfo.getRequiredSymbols() (default method)
func (c *CompilerInput) GetRequiredSymbols() []string {
	return DependencyInfo_GetRequiredSymbols(c)
}

// compilerInput_concat is ImmutableSet.builder().addAll(first).addAll(second).build():
// the elements of both, in order, without duplicates (by Java equals).
//
// Java: CompilerInput.concat(Iterable<T>, Iterable<T>)
func compilerInput_concat[T any](first []T, second []T) []T {
	set := jutil.NewObjectLinkedHashSet[T]()
	for _, e := range first {
		set.Add(e)
	}
	for _, e := range second {
		set.Add(e)
	}
	return set.ToSlice()
}

// Java: CompilerInput.getPath()
func (c *CompilerInput) GetPath() *ModuleLoader_ModulePath {
	if c.modulePath == nil {
		moduleLoader := c.compiler.GetModuleLoader()
		c.modulePath = moduleLoader.Resolve(c.GetName())
	}
	return c.modulePath
}

// CompilerInput_ModuleType is the JavaScript module type.
type CompilerInput_ModuleType int32

const (
	CompilerInput_ModuleType_NONE CompilerInput_ModuleType = iota
	CompilerInput_ModuleType_GOOG
	CompilerInput_ModuleType_ES6
	CompilerInput_ModuleType_COMMONJS
	CompilerInput_ModuleType_JSON
	CompilerInput_ModuleType_IMPORTED_SCRIPT
)

var compilerInput_ModuleType_names = [...]string{
	"NONE", "GOOG", "ES6", "COMMONJS", "JSON", "IMPORTED_SCRIPT",
}

// CompilerInput_ModuleType_values is ModuleType.values().
func CompilerInput_ModuleType_values() []CompilerInput_ModuleType {
	return []CompilerInput_ModuleType{
		CompilerInput_ModuleType_NONE, CompilerInput_ModuleType_GOOG, CompilerInput_ModuleType_ES6,
		CompilerInput_ModuleType_COMMONJS, CompilerInput_ModuleType_JSON,
		CompilerInput_ModuleType_IMPORTED_SCRIPT,
	}
}

// String is Enum.name().
func (t CompilerInput_ModuleType) String() string { return compilerInput_ModuleType_names[t] }

// Ordinal is Enum.ordinal().
func (t CompilerInput_ModuleType) Ordinal() int { return int(t) }

var compilerInput_ModuleType_ptrs = CompilerInput_ModuleType_values()

// Ptr returns the canonical pointer to t.
//
// Go-only: language-adapter — a nullable Java ModuleType reference is a
// *CompilerInput_ModuleType; nil is null.
func (t CompilerInput_ModuleType) Ptr() *CompilerInput_ModuleType {
	return &compilerInput_ModuleType_ptrs[t]
}

// Java: CompilerInput.getFeatures(AbstractCompiler)
func (c *CompilerInput) GetFeatures(compiler AbstractCompiler) *parser.FeatureSet {
	_ = c.ast.GetAstRoot(compiler) // parse if required
	return c.ast.features
}

// CompilerInput_JsAst wraps the Rhino AST for this CompilerInput.
//
// Use this wrapper class instead of directly storing a Node because the AST
// is lazily constructed. CompilerInputs are created before the parsing step,
// and in some cases the AST for a given CompilerInput may never be parsed if
// dependency pruning is enabled.
//
// Additionally if using precompiled libraries, this class handles
// deserializing the AST from the TypedAST format, instead of parsing it from
// the source.
type CompilerInput_JsAst struct {
	// Go-only: language-adapter — the enclosing CompilerInput of this
	// (non-static) inner class.
	outer    *CompilerInput
	root     rhino.Node
	features *parser.FeatureSet
}

func newCompilerInput_JsAst(outer *CompilerInput) *CompilerInput_JsAst {
	return &CompilerInput_JsAst{outer: outer}
}

// Java: CompilerInput.JsAst.getAstRoot(AbstractCompiler)
func (a *CompilerInput_JsAst) GetAstRoot(compiler AbstractCompiler) rhino.Node {
	if a.isParsed() {
		return a.root
	}

	sourceFile := a.outer.sourceFile
	astRootSource := compiler.GetTypedAstDeserializer(sourceFile)
	if astRootSource != nil {
		a.installDeserializedRoot(astRootSource())
	} else {
		a.parse(compiler)
	}
	jutil.CheckState(a.root.GetStaticSourceFile() == rhino.StaticSourceFile(sourceFile))
	a.root.SetInputId(a.outer.id)
	// Clear the cached source after parsing.  It will be re-read for snippet generation if
	// needed.
	sourceFile.ClearCachedSource()
	return a.root
}

func (a *CompilerInput_JsAst) installDeserializedRoot(root rhino.Node) {
	a.root = root
	if p := a.root.GetProp(rhino.Node_FEATURE_SET); p != nil {
		a.features = p.(*parser.FeatureSet)
	} else {
		a.features = nil
	}
}

// Java: CompilerInput.JsAst.clearAst()
func (a *CompilerInput_JsAst) ClearAst() {
	a.root = rhino.NilNode
	// While we're at it, clear out any saved text in the source file on
	// the assumption that if we're dumping the parse tree, then we probably
	// assume regenerating everything else is a smart idea also.
	a.outer.sourceFile.ClearCachedSource()
}

// Java: CompilerInput.JsAst.isParsed()
func (a *CompilerInput_JsAst) isParsed() bool {
	return !a.root.IsNil()
}

// Java: CompilerInput.JsAst.parse(AbstractCompiler)
func (a *CompilerInput_JsAst) parse(compiler AbstractCompiler) {
	sourceFile := a.outer.sourceFile
	jutil.Try(func() {
		context := AbstractCompiler_ConfigContext_DEFAULT
		if sourceFile.IsExtern() {
			context = AbstractCompiler_ConfigContext_EXTERNS
		}
		result := parsing.ParserRunner_Parse(
			sourceFile,
			sourceFile.GetCode(),
			compiler.GetParserConfig(context),
			compiler.GetDefaultErrorReporter())
		a.installParseResult(compiler, result)
	}, func(e jutil.IIOException) {
		compiler.Report(
			JSError_Make(AbstractCompiler_READ_ERROR, sourceFile.GetName(), e.GetMessage()))
	})

	a.ensureScriptRoot()
}

// Go-only: language-adapter — serial publication after independent parsing.
func (a *CompilerInput_JsAst) installParseResult(compiler AbstractCompiler, result *parsing.ParserRunner_ParseResult) {
	sourceFile := a.outer.sourceFile
	a.root = result.Ast
	a.features = result.Features

	if compiler.GetOptions().PreservesDetailedSourceInfo() {
		compiler.AddComments(sourceFile.GetName(), result.Comments)
	}
	if result.SourceMapURL != nil && compiler.GetOptions().GetResolveSourceMapAnnotations() {
		parseInline := compiler.GetOptions().GetParseInlineSourceMaps()
		sourceMapSourceFile :=
			SourceMapResolver_ExtractSourceMap(sourceFile, *result.SourceMapURL, parseInline)
		if sourceMapSourceFile != nil {
			compiler.AddInputSourceMap(
				sourceFile.GetName(), NewSourceMapInput(sourceMapSourceFile))
		}
	}
}

func (a *CompilerInput_JsAst) ensureScriptRoot() {
	sourceFile := a.outer.sourceFile
	if a.root.IsNil() {
		a.root = rhino.IR_Script()
	}

	// Set the source name so that the compiler passes can track
	// the source file and module.
	a.root.SetStaticSourceFile(sourceFile)
}
