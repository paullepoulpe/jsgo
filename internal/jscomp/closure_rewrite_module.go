// Port of com/google/javascript/jscomp/ClosureRewriteModule.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// the shapes of Closure calls and module statements: the callee and
// arguments of goog.module/require/forwardDeclare/module.get/requireDynamic
// calls (each record and update method takes the ast.Call), the export
// assignments (`exports = value`, `exports.name = value`, read through
// ast.Assign), the members of an exported object literal, the alias or
// destructuring pattern of a require, the `then` callback and the `await`
// of a dynamic require, the function of goog.loadModule, and the name, init
// and value of rewritten declarations. The port no longer follows the Java
// line by line.
//
// Java quirk kept: the namespace argument of the goog.* calls is the call's
// last child, which is the callee when the call has no arguments, so the
// "must be a string literal" errors of `goog.module()` are reported at the
// callee (namespaceArg says so where it is read).
//
// Left on the rhino API: the per-node token switches of the preprocessor,
// recorder and updater (they hand typed calls on), the JSDoc type
// rewriting (JSDoc type expressions are not in the grammar), the
// qualified-name matching against the prebuilt goog.* nodes, and code that
// looks at a node's parents or siblings of any kind (the statement around
// a require, the declaration a name belongs to).

// ClosureRewriteModule processes aliases in goog.modules.
//
//	goog.module('foo.Bar');
//	var Baz = goog.require('foo.Baz');
//	class Bar extends Baz {}
//	exports = Bar;
//
// becomes
//
//	class module$contents$foo$Bar_Bar extends module$exports$foo$Baz {}
//	var module$exports$foo$Bar = module$contents$foo$Bar_Bar;
type ClosureRewriteModule struct {
	compiler                AbstractCompiler
	astFactory              *AstFactory
	preprocessorSymbolTable *PreprocessorSymbolTable
	preserveSugar           bool

	syntheticExternNodes []rhino.Node
	syntheticExternSeen  map[string]struct{}

	globalScope *Scope // set after Process() starts updating scripts

	// Per script state needed for rewriting.
	scriptStack   []*closureScriptDescription
	currentScript *closureScriptDescription

	rewriteState closureGlobalRewriteState
	// All prefix namespaces from goog.provides and legacy goog.modules.
	legacyScriptNamespacesAndPrefixes map[string]struct{}
	unrecognizedRequires              []closureUnrecognizedRequire
	googModuleGetCalls                []ast.Call
	googRequireDynamicCalls           []ast.Call

	globalTypedScope *TypedScope
	xid              *Xid
}

var _ CompilerPass = (*ClosureRewriteModule)(nil)

// Java: ClosureRewriteModule.INVALID_MODULE_ID_ARG
var ClosureRewriteModule_INVALID_MODULE_ID_ARG = DiagnosticType_Error(
	"JSC_GOOG_MODULE_INVALID_MODULE_ID_ARG",
	"goog.module parameter must be a string literal")

// Java: ClosureRewriteModule.INVALID_PROVIDE_NAMESPACE
var ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE = DiagnosticType_Error(
	"JSC_GOOG_MODULE_INVALID_PROVIDE_NAMESPACE",
	"goog.provide parameter must be a string literal.")

// Java: ClosureRewriteModule.INVALID_PROVIDE_CALL
var ClosureRewriteModule_INVALID_PROVIDE_CALL = DiagnosticType_Error(
	"JSC_GOOG_MODULE_INVALID_PROVIDE_CALL", "goog.provide can not be called in goog.module.")

// Java: ClosureRewriteModule.INVALID_GET_ALIAS
var ClosureRewriteModule_INVALID_GET_ALIAS = DiagnosticType_Error(
	"JSC_GOOG_MODULE_INVALID_GET_ALIAS", "goog.module.get should not be aliased.")

// Java: ClosureRewriteModule.INVALID_EXPORT_COMPUTED_PROPERTY
var ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY = DiagnosticType_Error(
	"JSC_GOOG_MODULE_INVALID_EXPORT_COMPUTED_PROPERTY",
	"Computed properties are not yet supported in goog.module exports.")

// Java: ClosureRewriteModule.USELESS_USE_STRICT_DIRECTIVE
var ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE = DiagnosticType_Disabled(
	"JSC_USELESS_USE_STRICT_DIRECTIVE", "'use strict' is unnecessary in goog.module files.")

// Java: ClosureRewriteModule.IMPORT_INLINING_SHADOWS_VAR
var ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR = DiagnosticType_Error(
	"JSC_IMPORT_INLINING_SHADOWS_VAR",
	"Inlining of reference to import \"{1}\" shadows var \"{0}\".")

// Java: ClosureRewriteModule.ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT
var ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT = DiagnosticType_Error(
	"JSC_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT",
	"Destructuring import only allowed for importing module with named exports.\n"+
		"See https://github.com/google/closure-compiler/wiki/goog.module-style")

// Java: ClosureRewriteModule.ILLEGAL_DESTRUCTURING_NOT_EXPORTED
var ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED = DiagnosticType_Error(
	"JSC_ILLEGAL_DESTRUCTURING_NOT_EXPORTED",
	"Destructuring import reference to name \"{0}\" was not exported in module {1}")

// Java: ClosureRewriteModule.LOAD_MODULE_FN_MISSING_RETURN
var ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN = DiagnosticType_Error(
	"JSC_LOAD_MODULE_FN_MISSING_RETURN",
	"goog.loadModule function should end with 'return exports;'")

// Java: ClosureRewriteModule.ILLEGAL_MODULE_RENAMING_CONFLICT
var ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT = DiagnosticType_Error(
	"JSC_ILLEGAL_MODULE_RENAMING_CONFLICT",
	"Internal compiler error: rewritten module global name {0} is already in use.\n"+
		"Original definition: {1}")

// Java: ClosureRewriteModule.ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT
var ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT = DiagnosticType_Error(
	"ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT",
	"Illegal use of dynamic import: LHS of await goog.requireDynamic() must be a destructing"+
		" LHS or name, and it must be in a declaration statement.")

// Java: ClosureRewriteModule.MODULE_EXPORTS_PREFIX
const ClosureRewriteModule_MODULE_EXPORTS_PREFIX = "module$exports$"

const closureRewriteModule_MODULE_CONTENTS_PREFIX = "module$contents$"

// Prebuilt Nodes to speed up Node.matchesQualifiedName() calls
var (
	closureRewriteModule_GOOG_FORWARDDECLARE                    = rhino.IR_Getprop(rhino.IR_Name("goog"), "forwardDeclare")
	closureRewriteModule_GOOG_LOADMODULE                        = rhino.IR_Getprop(rhino.IR_Name("goog"), "loadModule")
	closureRewriteModule_GOOG_MODULE                            = rhino.IR_Getprop(rhino.IR_Name("goog"), "module")
	closureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE     = rhino.IR_Getprop(closureRewriteModule_GOOG_MODULE, "declareLegacyNamespace")
	closureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING = rhino.QualifiedName_Of("goog.module.preventModuleExportSealing")
	closureRewriteModule_GOOG_MODULE_GET                        = rhino.IR_Getprop(closureRewriteModule_GOOG_MODULE.CloneTree(), "get")
	closureRewriteModule_GOOG_PROVIDE                           = rhino.IR_Getprop(rhino.IR_Name("goog"), "provide")
	closureRewriteModule_GOOG_REQUIRE                           = rhino.IR_Getprop(rhino.IR_Name("goog"), "require")
	closureRewriteModule_GOOG_REQUIRETYPE                       = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireType")
	closureRewriteModule_GOOG_REQUIREDYNAMIC                    = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireDynamic")
)

const closureRewriteModule_GOOG_REQUIREDYNAMIC_NAME = "goog.requireDynamic"
const closureRewriteModule_IMPORT_HANDLER_NAME = "importHandler_"

type closureAddAt uint8

const (
	closureAddBefore closureAddAt = iota
	closureAddAfter
)

type closureScopeType uint8

const (
	closureScopeExecContext closureScopeType = iota
	closureScopeBlock
)

type closureUnrecognizedRequire struct {
	requireNode rhino.Node
	namespaceId string
}

type closureExportDefinition struct {
	exportName string
	isNamed    bool
	rhs        rhino.Node
	nameDecl   *Var
}

func isInlinableExportNameParent(t rhino.Token) bool {
	switch t {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET, rhino.Token_FUNCTION, rhino.Token_CLASS:
		return true
	default:
		return false
	}
}

func newDefaultExportDefinition(t *NodeTraversal, rhs rhino.Node) closureExportDefinition {
	return newExportDefinition(t, "", false, rhs)
}

func newNamedExportDefinition(t *NodeTraversal, name string, rhs rhino.Node) closureExportDefinition {
	return newExportDefinition(t, name, true, rhs)
}

func newExportDefinition(t *NodeTraversal, name string, isNamed bool, rhs rhino.Node) closureExportDefinition {
	def := closureExportDefinition{
		exportName: name,
		isNamed:    isNamed,
		rhs:        rhs,
	}
	if !rhs.IsNil() && (rhs.IsName() || rhs.IsStringKey()) {
		def.nameDecl = t.GetScope().GetVar(rhs.GetString())
	}
	return def
}

func (e *closureExportDefinition) exportPostfix() string {
	if !e.isNamed {
		return ""
	}
	return "." + e.exportName
}

func (e *closureExportDefinition) hasInlinableName(exportedVars map[*Var]struct{}) bool {
	if e.nameDecl == nil {
		return false
	}
	if _, alreadyExported := exportedVars[e.nameDecl]; alreadyExported {
		return false
	}
	parent := e.nameDecl.GetParentNode()
	if !isInlinableExportNameParent(parent.GetToken()) || NodeUtil_IsFunctionDeclaration(parent) {
		return false
	}
	initialValue, ok := ast.AsCall(e.nameDecl.GetInitialValue())
	if !ok {
		return true
	}
	method, ok := initialValue.Callee().AsGetProp()
	if !ok {
		return true
	}
	maybeGoog, ok := method.Object().AsName()
	if !ok || maybeGoog.Ident() != "goog" {
		return true
	}
	name := method.Property()
	return name != "require" && name != "forwardDeclare" && name != "getMsg"
}

func (e *closureExportDefinition) localName() (string, bool) {
	if e.nameDecl != nil {
		return e.nameDecl.GetName(), true
	}
	return "", false
}

type closureAliasName struct {
	newName      string
	namespaceId  string
	hasNamespace bool
}

type closureScriptDescription struct {
	isModule               bool
	declareLegacyNamespace bool
	hasNamespaceId         bool
	namespaceId            string // "a.b.c"
	binaryNamespace        string // "module$exports$a$b$c"
	contentsPrefix         string // "module$contents$a$b$c_"
	topLevelNames          map[string]struct{}
	childScripts           []*closureScriptDescription
	namesToInlineByAlias   map[string]closureAliasName

	// Transient state.
	willCreateExportsObject   bool
	hasCreatedExportObject    bool
	hasDefaultExport          bool
	defaultExport             closureExportDefinition
	hasDefaultExportLocalName bool
	defaultExportLocalName    string
	namedExports              map[string]struct{}
	exportsToInline           []closureExportDefinition
	exportsToInlineVars       map[*Var]struct{}

	// The MODULE_BODY node (or SCRIPT node at outer level) that contains the contents.
	rootNode rhino.Node
}

func (s *closureScriptDescription) addChildScript(childScript *closureScriptDescription) {
	s.childScripts = append(s.childScripts, childScript)
}

func (s *closureScriptDescription) removeFirstChildScript() *closureScriptDescription {
	if len(s.childScripts) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	first := s.childScripts[0]
	s.childScripts[0] = nil
	s.childScripts = s.childScripts[1:]
	return first
}

func (s *closureScriptDescription) getBinaryNamespace() (string, bool) {
	if !s.isModule || s.declareLegacyNamespace {
		return "", false
	}
	if !s.hasNamespaceId {
		panic(jutil.NewNullPointerException())
	}
	return s.binaryNamespace, true
}

func (s *closureScriptDescription) binaryNamespaceOrNull() string {
	if ns, ok := s.getBinaryNamespace(); ok {
		return ns
	}
	return "null"
}

func (s *closureScriptDescription) contentsPrefixOrNull() string {
	if !s.hasNamespaceId {
		return "null"
	}
	return s.contentsPrefix
}

func (s *closureScriptDescription) namespaceIdOrNull() string {
	if !s.hasNamespaceId {
		return "null"
	}
	return s.namespaceId
}

func (s *closureScriptDescription) getExportedNamespace() (string, bool) {
	if s.declareLegacyNamespace {
		return s.namespaceId, s.hasNamespaceId
	}
	return s.getBinaryNamespace()
}

// Java: ClosureRewriteModule.getBinaryModuleNamespace(String)
func ClosureRewriteModule_GetBinaryModuleNamespace(namespaceId string) string {
	return ClosureRewriteModule_MODULE_EXPORTS_PREFIX + strings.ReplaceAll(namespaceId, ".", "$")
}

type closureScriptPreprocessor struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *ClosureRewriteModule
}

func (p *closureScriptPreprocessor) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_MODULE_BODY:
		return true
	case rhino.Token_SCRIPT:
		if NodeUtil_IsGoogModuleFile(n) {
			ClosureRewriteModule_CheckAndSetStrictModeDirective(t, n)
		}
		return true
	case rhino.Token_NAME:
		p.outer.preprocessExportDeclaration(t, n)
		return true
	default:
		// Don't traverse into non-module scripts.
		return !parent.IsScript()
	}
}

type closureScriptRecorder struct {
	outer *ClosureRewriteModule
}

func (*closureScriptRecorder) noVisitLeaves() {}

func (r *closureScriptRecorder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	m := r.outer
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
		m.recordModuleBody(n)
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		getProp, ok := call.Callee().AsGetProp()
		if !ok {
			break
		}
		method := getProp.Node()
		if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE) {
			m.recordGoogModule(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE) {
			m.recordGoogDeclareLegacyNamespace()
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_PROVIDE) {
			m.recordGoogProvide(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRE) {
			m.recordGoogRequire(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRETYPE) {
			m.recordGoogRequireType(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIREDYNAMIC) {
			m.recordGoogRequireDynamic(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_FORWARDDECLARE) && !parent.IsExprResult() {
			m.recordGoogForwardDeclare(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_GET) {
			m.recordGoogModuleGet(t, call)
		}
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		if m.isTopLevel(t, n, closureScopeBlock) {
			m.recordTopLevelClassOrFunctionName(n)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		scopeType := closureScopeBlock
		if n.IsVar() {
			scopeType = closureScopeExecContext
		}
		if m.isTopLevel(t, n, scopeType) {
			m.recordTopLevelVarNames(n)
		}
	case rhino.Token_GETPROP:
		if getProp := ast.MustGetProp(n); m.currentScript.isModule && closureRewriteModule_isExportPropertyAssignment(t, getProp) {
			m.recordExportsPropertyAssignment(t, getProp)
		}
	case rhino.Token_NAME:
		if m.currentScript.isModule {
			m.maybeRecordExportDeclaration(t, n)
		}
	}

	return true
}

func (r *closureScriptRecorder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsModuleBody() {
		r.outer.popScript()
	}
}

type closureGlobalRewriteState struct {
	scriptDescriptionsByGoogModuleNamespace map[string]*closureScriptDescription
	providedNamespaces                      map[string]struct{}
}

func (g *closureGlobalRewriteState) containsModule(namespaceId string) bool {
	_, ok := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	return ok
}

func (g *closureGlobalRewriteState) isLegacyModule(namespaceId string) bool {
	s, ok := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	jutil.CheckArgument(ok)
	return s.declareLegacyNamespace
}

func (g *closureGlobalRewriteState) getBinaryNamespace(namespaceId string) (string, bool) {
	script := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	if script == nil {
		return "", false
	}
	return script.getBinaryNamespace()
}

func (g *closureGlobalRewriteState) getGoogModuleNamespaceType(namespaceId string) rhino.JSType {
	googModule := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	if googModule == nil {
		return nil
	}
	return googModule.rootNode.GetJSType()
}

func (g *closureGlobalRewriteState) getExportedNamespaceOrScript(namespaceId string) (string, bool) {
	if _, ok := g.providedNamespaces[namespaceId]; ok {
		return namespaceId, true
	}
	script := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	if script == nil {
		return "", false
	}
	return script.getExportedNamespace()
}

// Java: ClosureRewriteModule(AbstractCompiler, PreprocessorSymbolTable, TypedScope)
func NewClosureRewriteModule(
	compiler AbstractCompiler,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	globalTypedScope *TypedScope) *ClosureRewriteModule {
	jutil.CheckArgument(globalTypedScope == nil || globalTypedScope.IsGlobal())

	return &ClosureRewriteModule{
		compiler:                compiler,
		astFactory:              compiler.CreateAstFactory(),
		preprocessorSymbolTable: preprocessorSymbolTable,
		preserveSugar:           compiler.GetOptions().ShouldPreserveGoogModule(),
		rewriteState: closureGlobalRewriteState{
			scriptDescriptionsByGoogModuleNamespace: make(map[string]*closureScriptDescription),
			providedNamespaces:                      make(map[string]struct{}),
		},
		legacyScriptNamespacesAndPrefixes: make(map[string]struct{}),
		globalTypedScope:                  globalTypedScope,
	}
}

func (m *ClosureRewriteModule) globalTypedStaticScope() rhino.StaticScope {
	if m.globalTypedScope == nil {
		return nil
	}
	return m.globalTypedScope
}

func (m *ClosureRewriteModule) isWeakNamespace(namespaceId string) bool {
	metadata := m.compiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get(namespaceId)
	return metadata != nil && metadata.RootNode().GetStaticSourceFile().IsWeak()
}

func (m *ClosureRewriteModule) isUnrequiredWeakNamespace(call rhino.Node, namespaceId string) bool {
	if call.GetStaticSourceFile().IsWeak() || !m.isWeakNamespace(namespaceId) {
		return false
	}
	var callerMetadata *ModuleMetadataMap_ModuleMetadata
	if sourceFileName := call.GetSourceFileName(); sourceFileName != nil {
		callerMetadata = m.compiler.GetModuleMetadataMap().GetModulesByPath().Get(*sourceFileName)
	}
	return callerMetadata == nil ||
		!callerMetadata.StronglyRequiredGoogNamespaces().Contains(namespaceId)
}

type closureUnwrapGoogLoadModule struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *ClosureRewriteModule
}

func (u *closureUnwrapGoogLoadModule) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	compiler := u.outer.compiler
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT:
		return true
	case rhino.Token_EXPR_RESULT:
		call, ok := ast.MustExprStmt(n).Expr().AsCall()
		if !ok || !NodeUtil_IsCallToNode(call.Node(), closureRewriteModule_GOOG_LOADMODULE) {
			return false
		}
		lastArg, _ := call.LastArg()
		if function, ok := ast.AsFunction(lastArg.Node()); ok {
			parent.PutBooleanProp(rhino.Node_GOOG_MODULE, true)
			functionNode := function.Node()
			compiler.ReportFunctionDeleted(functionNode)
			moduleBody := function.Body().Detach().Node()
			moduleBody.SetToken(rhino.Token_MODULE_BODY)
			exportsParameter := function.Params().Node().GetOnlyChild()
			moduleBody.SetJSType(exportsParameter.GetJSType())
			n.ReplaceWith(moduleBody)
			returnNode := moduleBody.GetLastChild()
			if !returnNode.IsReturn() {
				compiler.Report(JSError_MakeWithNode(moduleBody, ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN))
			} else {
				returnNode.Detach()
			}
			t.ReportCodeChange()
		}
		return false
	default:
		return false
	}
}

// Java: ClosureRewriteModule.process(Node, Node)
func (m *ClosureRewriteModule) Process(externs rhino.Node, root rhino.Node) {
	sc := NewSyntacticScopeCreator(m.compiler)
	defer func() {
		m.globalScope = nil
	}()
	unwrap := closureUnwrapGoogLoadModule{outer: m}
	preprocessor := closureScriptPreprocessor{outer: m}
	recorder := closureScriptRecorder{outer: m}
	bPrep := NodeTraversal_Builder{
		compiler:     m.compiler,
		callback:     &preprocessor,
		scopeCreator: sc,
	}
	bRec := NodeTraversal_Builder{
		compiler:     m.compiler,
		callback:     &recorder,
		scopeCreator: sc,
	}

	scriptCount := externs.GetChildCount() + root.GetChildCount()
	scriptSlab := make([]closureScriptDescription, scriptCount)
	scriptDescriptions := make([]*closureScriptDescription, 0, scriptCount)
	slabIdx := 0
	for _, parent := range [2]rhino.Node{externs, root} {
		for script := parent.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			jutil.CheckStateMsg1(script.IsScript(), "%s", script)
			NodeTraversal_Traverse(m.compiler, script, &unwrap)

			var desc *closureScriptDescription
			if slabIdx < len(scriptSlab) {
				desc = &scriptSlab[slabIdx]
				slabIdx++
			} else {
				desc = &closureScriptDescription{}
			}
			m.pushScript(desc)

			m.currentScript.rootNode = script
			scriptDescriptions = append(scriptDescriptions, m.currentScript)
			bPrep.Traverse(script)
			bRec.Traverse(script)
			m.popScript()
		}
	}

	m.reportUnrecognizedRequires()
	if m.compiler.HasHaltingErrors() {
		return
	}

	updater := closureScriptUpdater{outer: m, scriptDescriptions: scriptDescriptions}
	bUpd := NodeTraversal_Builder{
		compiler:     m.compiler,
		callback:     &updater,
		scopeCreator: sc,
	}
	bUpd.TraverseRoots(externs, root)
	m.declareSyntheticExterns()
	for _, call := range m.googModuleGetCalls {
		m.updateGoogModuleGetCall(call)
	}
	for _, call := range m.googRequireDynamicCalls {
		m.updateGoogRequireDynamicCall(call)
	}
}

func (m *ClosureRewriteModule) addSyntheticExternIfAbsent(lhs rhino.Node) {
	name := lhs.GetString()
	if m.syntheticExternSeen == nil {
		m.syntheticExternSeen = make(map[string]struct{})
	}
	if _, seen := m.syntheticExternSeen[name]; !seen {
		m.syntheticExternSeen[name] = struct{}{}
		m.syntheticExternNodes = append(m.syntheticExternNodes, lhs)
	}
}

func (m *ClosureRewriteModule) declareSyntheticExterns() {
	if len(m.syntheticExternNodes) == 0 {
		return
	}
	var vars []rhino.Node
	for _, lhs := range m.syntheticExternNodes {
		name := lhs.GetString()
		if m.isNameInGlobalScope(name) {
			continue
		}
		vars = append(vars,
			rhino.IR_Var(m.astFactory.CreateNameWithUnknownType(name)).SrcrefTree(lhs))
	}
	if len(vars) == 0 {
		return
	}
	root := m.compiler.GetSynthesizedExternsInput().GetAstRoot(m.compiler)
	for _, v := range vars {
		root.AddChildToBack(v)
	}
}

func (m *ClosureRewriteModule) isNameInGlobalScope(name string) bool {
	_, isLegacy := m.legacyScriptNamespacesAndPrefixes[name]
	return isLegacy || m.globalScope.GetVar(name) != nil
}

func (m *ClosureRewriteModule) preprocessExportDeclaration(t *NodeTraversal, n rhino.Node) {
	if !closureRewriteModule_isGoogModuleExportsRef(t, n) {
		return
	}
	assign, ok := closureRewriteModule_assignToTarget(n)
	if !ok || !n.GetGrandparent().IsExprResult() {
		return
	}

	jutil.CheckState(!m.currentScript.hasDefaultExport)
	exportRhs := assign.Value().Node()
	if NodeUtil_IsNamedExportsLiteral(exportRhs) {
		insertionPoint := n.GetGrandparent()
		for member := range ast.MustObjectLit(exportRhs).Members() {
			key := ast.MustStringKey(member.Node())
			exportName := key.Key()
			jsdoc := key.Node().GetJSDocInfo()
			rhs := key.Value().Detach().Node()
			lhs := m.astFactory.
				CreateGetProp(m.astFactory.CreateName("exports", AstFactory_Type_(n)), exportName, AstFactory_Type_(rhs)).
				SrcrefTree(key.Node())
			newExport := rhino.IR_ExprResult(m.astFactory.CreateAssign(lhs, rhs).Srcref(key.Node()).SetJSDocInfo(jsdoc)).
				Srcref(key.Node())
			newExport.InsertAfter(insertionPoint)
			insertionPoint = newExport
		}
		n.GetGrandparent().Detach()
	}
}

func (m *ClosureRewriteModule) recordModuleBody(moduleRoot rhino.Node) {
	m.pushScript(&closureScriptDescription{
		rootNode: moduleRoot,
		isModule: true,
	})
}

func (m *ClosureRewriteModule) recordGoogModule(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_MODULE_ID_ARG)
		return
	}
	namespaceId := namespaceIdNode.GetString()
	escaped := strings.ReplaceAll(namespaceId, ".", "$")

	m.currentScript.hasNamespaceId = true
	m.currentScript.namespaceId = namespaceId
	m.currentScript.binaryNamespace = ClosureRewriteModule_MODULE_EXPORTS_PREFIX + escaped
	m.currentScript.contentsPrefix = closureRewriteModule_MODULE_CONTENTS_PREFIX + escaped + "_"

	m.rewriteState.scriptDescriptionsByGoogModuleNamespace[namespaceId] = m.currentScript
}

func (m *ClosureRewriteModule) recordGoogDeclareLegacyNamespace() {
	m.currentScript.declareLegacyNamespace = true
	if !m.currentScript.hasNamespaceId {
		panic(jutil.NewNullPointerException())
	}
	m.updateLegacyScriptNamespacesAndPrefixes(m.currentScript.namespaceId)
}

func (m *ClosureRewriteModule) updateLegacyScriptNamespacesAndPrefixes(namespace string) {
	m.legacyScriptNamespacesAndPrefixes[namespace] = struct{}{}
	for dot := strings.LastIndexByte(namespace, '.'); dot != -1; dot = strings.LastIndexByte(namespace, '.') {
		namespace = namespace[:dot]
		m.legacyScriptNamespacesAndPrefixes[namespace] = struct{}{}
	}
}

func (m *ClosureRewriteModule) recordGoogProvide(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	if m.currentScript.isModule {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_PROVIDE_CALL)
	}

	m.rewriteState.providedNamespaces[namespaceId] = struct{}{}
	m.updateLegacyScriptNamespacesAndPrefixes(namespaceId)
}

func (m *ClosureRewriteModule) recordGoogRequire(t *NodeTraversal, call ast.Call) {
	closureRewriteModule_maybeSplitMultiVar(call.Node())

	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	targetIsAModule := m.rewriteState.containsModule(namespaceId)
	_, targetIsALegacyScript := m.rewriteState.providedNamespaces[namespaceId]
	if m.currentScript.isModule && !targetIsAModule && !targetIsALegacyScript {
		m.unrecognizedRequires = append(m.unrecognizedRequires, closureUnrecognizedRequire{requireNode: call.Node(), namespaceId: namespaceId})
	}
}

func (m *ClosureRewriteModule) recordGoogRequireType(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE)
		return
	}
	m.recordGoogRequire(t, call)
}

func (m *ClosureRewriteModule) recordGoogForwardDeclare(t *NodeTraversal, call ast.Call) {
	namespaceNode := closureRewriteModule_namespaceArg(call)
	if call.NumArgs() != 1 || !namespaceNode.IsStringLit() {
		t.Report(namespaceNode, ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE)
		return
	}
	m.recordGoogRequire(t, call)
}

func (m *ClosureRewriteModule) recordGoogRequireDynamic(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC)
		return
	}

	namespaceId := namespaceIdNode.GetString()
	if !m.rewriteState.containsModule(namespaceId) {
		m.unrecognizedRequires = append(m.unrecognizedRequires, closureUnrecognizedRequire{requireNode: call.Node(), namespaceId: namespaceId})
	}
	m.googRequireDynamicCalls = append(m.googRequireDynamicCalls, call)
}

func (m *ClosureRewriteModule) recordGoogModuleGet(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if call.NumArgs() != 1 || !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	if !m.rewriteState.containsModule(namespaceId) {
		m.unrecognizedRequires = append(m.unrecognizedRequires, closureUnrecognizedRequire{requireNode: call.Node(), namespaceId: namespaceId})
	}
	m.googModuleGetCalls = append(m.googModuleGetCalls, call)

	maybeAssign := call.Node().GetParent()
	assign, ok := ast.AsAssign(maybeAssign)
	if !ok || !maybeAssign.GetParent().IsExprResult() || !m.currentScript.isModule {
		return
	}
	alias, ok := assign.Target().AsName()
	if !ok {
		return
	}

	aliasVar := t.GetScope().GetVar(alias.Ident())
	if aliasVar == nil {
		return
	}

	aliasVarNodeRhs := NodeUtil_GetRValueOfLValue(aliasVar.GetNode())
	if aliasVarNodeRhs.IsNil() ||
		!NodeUtil_IsCallToNode(aliasVarNodeRhs, closureRewriteModule_GOOG_FORWARDDECLARE) ||
		namespaceId != closureRewriteModule_namespaceArg(ast.MustCall(aliasVarNodeRhs)).GetString() {
		return
	}

	if m.isUnrequiredWeakNamespace(call.Node(), namespaceId) {
		return
	}

	m.compiler.ReportChangeToEnclosingScope(maybeAssign)
	maybeAssign.GetParent().Detach()
	m.googModuleGetCalls = m.googModuleGetCalls[:len(m.googModuleGetCalls)-1]
}

// closureRewriteModule_namespaceArg returns the namespace argument of a
// goog.module, goog.provide, goog.require... call: its last argument.
//
// Go-only: Java reads the call's last child, which is the callee when the
// call has no arguments; this keeps that (an error about the argument of
// `goog.module()` is reported at `goog.module`).
func closureRewriteModule_namespaceArg(call ast.Call) rhino.Node {
	if last, ok := call.LastArg(); ok {
		return last.Node()
	}
	return call.Callee().Node()
}

func (m *ClosureRewriteModule) recordTopLevelClassOrFunctionName(classOrFunctionNode rhino.Node) {
	var nameNode ast.Name
	var ok bool
	if fn, isFunction := ast.AsFunction(classOrFunctionNode); isFunction {
		nameNode, ok = fn.Name()
	} else {
		nameNode, ok = ast.MustClass(classOrFunctionNode).Name()
	}
	if ok {
		if m.currentScript.topLevelNames == nil {
			m.currentScript.topLevelNames = make(map[string]struct{})
		}
		m.currentScript.topLevelNames[nameNode.Ident()] = struct{}{}
	}
}

func (m *ClosureRewriteModule) recordTopLevelVarNames(varNode rhino.Node) {
	NodeUtil_VisitLhsNodesInNode(
		varNode, func(lhs rhino.Node) {
			if m.currentScript.topLevelNames == nil {
				m.currentScript.topLevelNames = make(map[string]struct{})
			}
			m.currentScript.topLevelNames[lhs.GetString()] = struct{}{}
		})
}

func (m *ClosureRewriteModule) maybeRecordExportDeclaration(t *NodeTraversal, n rhino.Node) {
	if !m.currentScript.isModule || !closureRewriteModule_isGoogModuleExportsRef(t, n) {
		return
	}
	assign, ok := closureRewriteModule_assignToTarget(n)
	if !ok {
		return
	}

	// ClosureCheckModule reports an error for duplicate 'exports = '
	// assignments, but that error may be suppressed. If so then use the final
	// assignment as the canonical one.
	if m.currentScript.hasDefaultExport {
		if localName, hasLocal := m.currentScript.defaultExport.localName(); hasLocal {
			delete(m.currentScript.namesToInlineByAlias, localName)
		}
		m.currentScript.hasDefaultExportLocalName = false
		m.currentScript.defaultExportLocalName = ""
	}
	exportRhs := assign.Value().Node()

	jutil.CheckStateMsg(
		!NodeUtil_IsNamedExportsLiteral(exportRhs),
		"Exports object should have been converted already")

	m.currentScript.willCreateExportsObject = true
	defaultExport := newDefaultExportDefinition(t, exportRhs)
	m.currentScript.hasDefaultExport = true
	m.currentScript.defaultExport = defaultExport
	if !m.currentScript.declareLegacyNamespace &&
		defaultExport.hasInlinableName(m.currentScript.exportsToInlineVars) {
		localName, _ := defaultExport.localName()
		m.currentScript.hasDefaultExportLocalName = true
		m.currentScript.defaultExportLocalName = localName
		m.recordExportToInline(defaultExport)
	}
}
