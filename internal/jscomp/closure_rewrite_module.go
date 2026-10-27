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
//
// and
//
//	goog.loadModule(function(exports) {
//	  goog.module('foo.Bar');
//	  var Baz = goog.require('foo.Baz');
//	  class Bar extends Baz {}
//	  exports = Bar;
//	  return exports;
//	})
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
	syntheticExterns        *jutil.LinkedHashMap[string, rhino.Node]

	globalScope *Scope // non-final because it must be set after process() is called

	// Per script state needed for rewriting.
	scriptStack   []*ClosureRewriteModule_ScriptDescription
	currentScript *ClosureRewriteModule_ScriptDescription

	rewriteState *ClosureRewriteModule_GlobalRewriteState
	// All prefix namespaces from goog.provides and legacy goog.modules.
	legacyScriptNamespacesAndPrefixes map[string]struct{}
	unrecognizedRequires              []ClosureRewriteModule_UnrecognizedRequire
	googModuleGetCalls                []ast.Call
	googRequireDynamicCalls           []ast.Call

	globalTypedScope *TypedScope
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

// ClosureRewriteModule_AddAt indicates where new nodes should be added in
// relation to some other node.
type ClosureRewriteModule_AddAt int32

const (
	ClosureRewriteModule_AddAt_BEFORE ClosureRewriteModule_AddAt = iota
	ClosureRewriteModule_AddAt_AFTER
)

// Java: ClosureRewriteModule.AddAt.name()
func (a ClosureRewriteModule_AddAt) String() string {
	return [...]string{"BEFORE", "AFTER"}[a]
}

// ClosureRewriteModule_ScopeType is ClosureRewriteModule.ScopeType.
type ClosureRewriteModule_ScopeType int32

const (
	ClosureRewriteModule_ScopeType_EXEC_CONTEXT ClosureRewriteModule_ScopeType = iota
	ClosureRewriteModule_ScopeType_BLOCK
)

// Java: ClosureRewriteModule.ScopeType.name()
func (s ClosureRewriteModule_ScopeType) String() string {
	return [...]string{"EXEC_CONTEXT", "BLOCK"}[s]
}

// ClosureRewriteModule_UnrecognizedRequire describes the context of an
// "unrecognized require" scenario so that it will be possible to categorize
// and report it as either a "not provided yet" or "not provided at all"
// error at the end.
type ClosureRewriteModule_UnrecognizedRequire struct {
	// A goog.require() call, or a goog.module.get() call.
	requireNode rhino.Node
	namespaceId string
}

// ClosureRewriteModule_ExportDefinition is
// ClosureRewriteModule.ExportDefinition.
type ClosureRewriteModule_ExportDefinition struct {
	// Null if the export is a default export (exports = expr)
	exportName *string
	// Null if the export is of a @typedef
	rhs rhino.Node
	// Null if the export is of anything other than a name
	nameDecl *Var
}

// Java: ClosureRewriteModule.ExportDefinition.toString()
func (e *ClosureRewriteModule_ExportDefinition) String() string {
	h := jutil.MoreObjects_ToStringHelper("ExportDefinition")
	if e.exportName != nil {
		h.Add("exportName", *e.exportName)
	} else {
		h.Add("exportName", nil)
	}
	if !e.rhs.IsNil() {
		h.Add("rhs", e.rhs)
	} else {
		h.Add("rhs", nil)
	}
	if e.nameDecl != nil {
		h.Add("nameDecl", e.nameDecl)
	} else {
		h.Add("nameDecl", nil)
	}
	return h.OmitNullValues().String()
}

func closureRewriteModule_isInlinableNameParent(t rhino.Token) bool {
	switch t {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET, rhino.Token_FUNCTION, rhino.Token_CLASS:
		return true
	default:
		return false
	}
}

// Java: ClosureRewriteModule.ExportDefinition.newDefaultExport(NodeTraversal, Node)
func ClosureRewriteModule_ExportDefinition_NewDefaultExport(t *NodeTraversal, rhs rhino.Node) *ClosureRewriteModule_ExportDefinition {
	return ClosureRewriteModule_ExportDefinition_NewNamedExport(t, nil, rhs)
}

// Java: ClosureRewriteModule.ExportDefinition.newNamedExport(NodeTraversal, String, Node)
func ClosureRewriteModule_ExportDefinition_NewNamedExport(t *NodeTraversal, name *string, rhs rhino.Node) *ClosureRewriteModule_ExportDefinition {
	newExport := &ClosureRewriteModule_ExportDefinition{}
	newExport.exportName = name
	newExport.rhs = rhs
	if !rhs.IsNil() && (rhs.IsName() || rhs.IsStringKey()) {
		newExport.nameDecl = t.GetScope().GetVar(rhs.GetString())
	}
	return newExport
}

// Java: ClosureRewriteModule.ExportDefinition.getExportPostfix()
func (e *ClosureRewriteModule_ExportDefinition) GetExportPostfix() string {
	if e.exportName == nil {
		return ""
	}
	return "." + *e.exportName
}

// HasInlinableName: exportedNames is the key set of a Map<Var,
// ExportDefinition>.
//
// Go-only: language-adapter — the Java parameter is the map's keySet()
// view; Go passes the map and reads its keys.
//
// Java: ClosureRewriteModule.ExportDefinition.hasInlinableName(Set<Var>)
func (e *ClosureRewriteModule_ExportDefinition) HasInlinableName(
	exportedNames *jutil.LinkedHashMap[*Var, *ClosureRewriteModule_ExportDefinition]) bool {
	if e.nameDecl == nil ||
		exportedNames.ContainsKey(e.nameDecl) ||
		!closureRewriteModule_isInlinableNameParent(e.nameDecl.GetParentNode().GetToken()) ||
		NodeUtil_IsFunctionDeclaration(e.nameDecl.GetParentNode()) {
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

// Java: ClosureRewriteModule.ExportDefinition.getLocalName()
func (e *ClosureRewriteModule_ExportDefinition) GetLocalName() *string {
	if e.nameDecl != nil {
		name := e.nameDecl.GetName()
		return &name
	}
	return nil
}

// ClosureRewriteModule_AliasName is ClosureRewriteModule.AliasName.
type ClosureRewriteModule_AliasName struct {
	newName     string
	namespaceId *string // non-null only if this is an alias of a module itself
}

// Java: ClosureRewriteModule.AliasName(String, String)
func NewClosureRewriteModule_AliasName(newName string, namespaceId *string) *ClosureRewriteModule_AliasName {
	return &ClosureRewriteModule_AliasName{newName: newName, namespaceId: namespaceId}
}

// ClosureRewriteModule_ScriptDescription is
// ClosureRewriteModule.ScriptDescription.
type ClosureRewriteModule_ScriptDescription struct {
	isModule               bool
	declareLegacyNamespace bool
	namespaceId            *string             // "a.b.c"
	contentsPrefix         *string             // "module$contents$a$b$c_
	topLevelNames          map[string]struct{} // For prefixed content renaming.
	childScripts           []*ClosureRewriteModule_ScriptDescription
	namesToInlineByAlias   map[string]*ClosureRewriteModule_AliasName // For alias inlining.

	// Transient state.
	willCreateExportsObject bool

	hasCreatedExportObject bool
	defaultExport          *ClosureRewriteModule_ExportDefinition
	defaultExportLocalName *string
	namedExports           map[string]struct{}
	exportsToInline        jutil.LinkedHashMap[*Var, *ClosureRewriteModule_ExportDefinition]

	// The root of the module. The MODULE_BODY node that contains the module
	// contents. For recognizing top level names.
	rootNode rhino.Node
}

// Java: ClosureRewriteModule.ScriptDescription()
func NewClosureRewriteModule_ScriptDescription() *ClosureRewriteModule_ScriptDescription {
	return &ClosureRewriteModule_ScriptDescription{}
}

// Java: ClosureRewriteModule.ScriptDescription.addChildScript(ScriptDescription)
func (s *ClosureRewriteModule_ScriptDescription) AddChildScript(childScript *ClosureRewriteModule_ScriptDescription) {
	s.childScripts = append(s.childScripts, childScript)
}

// Java: ClosureRewriteModule.ScriptDescription.removeFirstChildScript()
func (s *ClosureRewriteModule_ScriptDescription) RemoveFirstChildScript() *ClosureRewriteModule_ScriptDescription {
	if len(s.childScripts) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	first := s.childScripts[0]
	s.childScripts[0] = nil
	s.childScripts = s.childScripts[1:]
	return first
}

// GetBinaryNamespace is "module$exports$a$b$c" for non-legacy modules.
//
// Java: ClosureRewriteModule.ScriptDescription.getBinaryNamespace()
func (s *ClosureRewriteModule_ScriptDescription) GetBinaryNamespace() *string {
	if !s.isModule || s.declareLegacyNamespace {
		return nil
	}
	ns := ClosureRewriteModule_GetBinaryModuleNamespace(*s.namespaceId)
	return &ns
}

// Java: ClosureRewriteModule.ScriptDescription.getExportedNamespace()
func (s *ClosureRewriteModule_ScriptDescription) GetExportedNamespace() *string {
	if s.declareLegacyNamespace {
		return s.namespaceId
	}
	return s.GetBinaryNamespace()
}

// Java: ClosureRewriteModule.getBinaryModuleNamespace(String)
func ClosureRewriteModule_GetBinaryModuleNamespace(namespaceId string) string {
	return ClosureRewriteModule_MODULE_EXPORTS_PREFIX + strings.ReplaceAll(namespaceId, ".", "$")
}

// ClosureRewriteModule_ScriptPreprocessor is
// ClosureRewriteModule.ScriptPreprocessor (an inner class).
type ClosureRewriteModule_ScriptPreprocessor struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *ClosureRewriteModule
}

// Java: ClosureRewriteModule.ScriptPreprocessor.shouldTraverse(NodeTraversal, Node, Node)
func (p *ClosureRewriteModule_ScriptPreprocessor) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_MODULE_BODY:
		return true
	case rhino.Token_SCRIPT:
		if NodeUtil_IsGoogModuleFile(n) {
			ClosureRewriteModule_CheckAndSetStrictModeDirective(t, n)
		}
		return true
	case rhino.Token_NAME:
		p.outer.PreprocessExportDeclaration(t, n)
		return true
	default:
		// Don't traverse into non-module scripts.
		return !parent.IsScript()
	}
}

// ClosureRewriteModule_ScriptRecorder is ClosureRewriteModule.ScriptRecorder
// (an inner class).
type ClosureRewriteModule_ScriptRecorder struct {
	outer *ClosureRewriteModule
}

func (*ClosureRewriteModule_ScriptRecorder) noVisitLeaves() {}

// Java: ClosureRewriteModule.ScriptRecorder.shouldTraverse(NodeTraversal, Node, Node)
func (r *ClosureRewriteModule_ScriptRecorder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	m := r.outer
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
		m.RecordModuleBody(n)
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		getProp, ok := call.Callee().AsGetProp()
		if !ok {
			break
		}
		method := getProp.Node()
		if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE) {
			m.RecordGoogModule(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE) {
			m.RecordGoogDeclareLegacyNamespace()
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_PROVIDE) {
			m.RecordGoogProvide(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRE) {
			m.RecordGoogRequire(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRETYPE) {
			m.RecordGoogRequireType(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIREDYNAMIC) {
			m.RecordGoogRequireDynamic(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_FORWARDDECLARE) && !parent.IsExprResult() {
			m.RecordGoogForwardDeclare(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_GET) {
			m.RecordGoogModuleGet(t, call)
		}
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		if m.IsTopLevel(t, n, ClosureRewriteModule_ScopeType_BLOCK) {
			m.RecordTopLevelClassOrFunctionName(n)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		scopeType := ClosureRewriteModule_ScopeType_BLOCK
		if n.IsVar() {
			scopeType = ClosureRewriteModule_ScopeType_EXEC_CONTEXT
		}
		if m.IsTopLevel(t, n, scopeType) {
			m.RecordTopLevelVarNames(n)
		}
	case rhino.Token_GETPROP:
		if getProp := ast.MustGetProp(n); m.currentScript.isModule && closureRewriteModule_isExportPropertyAssignment(t, getProp) {
			m.RecordExportsPropertyAssignment(t, getProp)
		}
	case rhino.Token_NAME:
		if m.currentScript.isModule {
			m.MaybeRecordExportDeclaration(t, n)
		}
	default:
	}

	return true
}

// Java: ClosureRewriteModule.ScriptRecorder.visit(NodeTraversal, Node, Node)
func (r *ClosureRewriteModule_ScriptRecorder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsModuleBody() {
		r.outer.PopScript()
	}
}

// ClosureRewriteModule_GlobalRewriteState is the global state tracking an
// association between the dotted names of goog.module()s and whether the
// goog.module declares itself as a legacy namespace. Allows for detecting
// duplicate goog.module()s and for rewriting fully qualified JsDoc type
// references to goog.module() types in legacy scripts.
type ClosureRewriteModule_GlobalRewriteState struct {
	scriptDescriptionsByGoogModuleNamespace map[string]*ClosureRewriteModule_ScriptDescription
	providedNamespaces                      map[string]struct{}
}

// Java: ClosureRewriteModule.GlobalRewriteState()
func NewClosureRewriteModule_GlobalRewriteState() *ClosureRewriteModule_GlobalRewriteState {
	return &ClosureRewriteModule_GlobalRewriteState{
		scriptDescriptionsByGoogModuleNamespace: make(map[string]*ClosureRewriteModule_ScriptDescription),
		providedNamespaces:                      make(map[string]struct{}),
	}
}

// Java: ClosureRewriteModule.GlobalRewriteState.containsModule(String)
func (g *ClosureRewriteModule_GlobalRewriteState) ContainsModule(namespaceId string) bool {
	_, ok := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	return ok
}

// Java: ClosureRewriteModule.GlobalRewriteState.isLegacyModule(String)
func (g *ClosureRewriteModule_GlobalRewriteState) IsLegacyModule(namespaceId string) bool {
	jutil.CheckArgument(g.ContainsModule(namespaceId))
	return g.scriptDescriptionsByGoogModuleNamespace[namespaceId].declareLegacyNamespace
}

// Java: ClosureRewriteModule.GlobalRewriteState.getBinaryNamespace(String)
func (g *ClosureRewriteModule_GlobalRewriteState) GetBinaryNamespace(namespaceId string) *string {
	script := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	if script == nil {
		return nil
	}
	return script.GetBinaryNamespace()
}

// GetGoogModuleNamespaceType returns the type of a goog.require of the given
// goog.module, or null if not a module.
//
// Java: ClosureRewriteModule.GlobalRewriteState.getGoogModuleNamespaceType(String)
func (g *ClosureRewriteModule_GlobalRewriteState) GetGoogModuleNamespaceType(namespaceId string) rhino.JSType {
	googModule := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	if googModule == nil {
		return nil
	}
	return googModule.rootNode.GetJSType()
}

// Java: ClosureRewriteModule.GlobalRewriteState.getExportedNamespaceOrScript(String)
func (g *ClosureRewriteModule_GlobalRewriteState) GetExportedNamespaceOrScript(namespaceId string) *string {
	if _, ok := g.providedNamespaces[namespaceId]; ok {
		return &namespaceId
	}
	script := g.scriptDescriptionsByGoogModuleNamespace[namespaceId]
	if script == nil {
		return nil
	}
	return script.GetExportedNamespace()
}

// Java: ClosureRewriteModule(AbstractCompiler, PreprocessorSymbolTable, TypedScope)
func NewClosureRewriteModule(
	compiler AbstractCompiler,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	globalTypedScope *TypedScope) *ClosureRewriteModule {
	jutil.CheckArgument(globalTypedScope == nil || globalTypedScope.IsGlobal())

	m := &ClosureRewriteModule{
		syntheticExterns:                  jutil.NewLinkedHashMap[string, rhino.Node](),
		rewriteState:                      NewClosureRewriteModule_GlobalRewriteState(),
		legacyScriptNamespacesAndPrefixes: make(map[string]struct{}),
	}
	m.compiler = compiler
	m.astFactory = compiler.CreateAstFactory()
	m.preprocessorSymbolTable = preprocessorSymbolTable
	m.preserveSugar = compiler.GetOptions().ShouldPreserveGoogModule()
	m.globalTypedScope = globalTypedScope
	return m
}

// closureRewriteModule_staticScope is the nullable globalTypedScope as a
// rhino.StaticScope.
//
// Go-only: language-adapter — a nil *TypedScope must stay a nil interface.
func (m *ClosureRewriteModule) globalTypedStaticScope() rhino.StaticScope {
	if m.globalTypedScope == nil {
		return nil
	}
	return m.globalTypedScope
}

// Java: ClosureRewriteModule.isWeakNamespace(String)
func (m *ClosureRewriteModule) IsWeakNamespace(namespaceId string) bool {
	metadata :=
		m.compiler.GetModuleMetadataMap().GetModulesByGoogNamespace().Get(namespaceId)
	return metadata != nil && metadata.RootNode().GetStaticSourceFile().IsWeak()
}

// Java: ClosureRewriteModule.isUnrequiredWeakNamespace(Node, String)
func (m *ClosureRewriteModule) IsUnrequiredWeakNamespace(call rhino.Node, namespaceId string) bool {
	if call.GetStaticSourceFile().IsWeak() || !m.IsWeakNamespace(namespaceId) {
		return false
	}
	// Go-only: language-adapter — ImmutableMap.get(null) is null.
	var callerMetadata *ModuleMetadataMap_ModuleMetadata
	if sourceFileName := call.GetSourceFileName(); sourceFileName != nil {
		callerMetadata = m.compiler.GetModuleMetadataMap().GetModulesByPath().Get(*sourceFileName)
	}
	return callerMetadata == nil ||
		!callerMetadata.StronglyRequiredGoogNamespaces().Contains(namespaceId)
}

// ClosureRewriteModule_UnwrapGoogLoadModule is
// ClosureRewriteModule.UnwrapGoogLoadModule (an inner class).
type ClosureRewriteModule_UnwrapGoogLoadModule struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *ClosureRewriteModule
}

// Java: ClosureRewriteModule.UnwrapGoogLoadModule.shouldTraverse(NodeTraversal, Node, Node)
func (u *ClosureRewriteModule_UnwrapGoogLoadModule) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	compiler := u.outer.compiler
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT:
		return true
	case rhino.Token_EXPR_RESULT:
		call, ok := ast.MustExprStmt(n).Expr().AsCall()
		if !ok || !NodeUtil_IsCallToNode(call.Node(), closureRewriteModule_GOOG_LOADMODULE) {
			return false
		}
		// goog.loadModule(function(exports) { ... }): the function is the last argument.
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
	// Record all the scripts first so that the googModuleNamespaces global
	// state can be complete before doing any updating also queue up
	// scriptDescriptions for later use in ScriptUpdater runs.

	sc := NewSyntacticScopeCreator(m.compiler)
	defer func() {
		m.globalScope = nil
	}()
	unwrap := ClosureRewriteModule_UnwrapGoogLoadModule{outer: m}
	preprocessor := ClosureRewriteModule_ScriptPreprocessor{outer: m}
	recorder := ClosureRewriteModule_ScriptRecorder{outer: m}
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

	var scriptDescriptions []*ClosureRewriteModule_ScriptDescription
	for _, parent := range []rhino.Node{externs, root} {
		for script := parent.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			jutil.CheckStateMsg1(script.IsScript(), "%s", script)
			NodeTraversal_Traverse(m.compiler, script, &unwrap)

			m.PushScript(NewClosureRewriteModule_ScriptDescription()) // sets currentScript

			m.currentScript.rootNode = script
			scriptDescriptions = append(scriptDescriptions, m.currentScript)
			bPrep.Traverse(script)
			bRec.Traverse(script)
			m.PopScript()
		}
	}

	m.ReportUnrecognizedRequires()
	if m.compiler.HasHaltingErrors() {
		return
	}

	// Update scripts using the now complete googModuleNamespaces global state
	// and unspool the scriptDescriptions that were queued up by all the
	// recording.

	bUpd := NodeTraversal_Builder{
		compiler:     m.compiler,
		callback:     m.NewScriptUpdater(scriptDescriptions),
		scopeCreator: sc,
	}
	bUpd.TraverseRoots(externs, root)
	m.DeclareSyntheticExterns()
	for _, call := range m.googModuleGetCalls {
		m.UpdateGoogModuleGetCall(call)
	}
	for _, call := range m.googRequireDynamicCalls {
		m.UpdateGoogRequireDynamicCall(call)
	}
}

// DeclareSyntheticExterns declares `var foo;` in the externs for all
// synthetic_externs names that aren't already in the global scope.
//
// Only add externs in the error case where there is an unrecognized
// goog.require or goog.module.get. Clutz depends on the compiler deleting
// local name declarations that alias or destructure the invalid call from
// this file. In order to preserve AST validity, we instead declare the
// deleted name in the externs. (which is fine with Clutz, as it just ignores
// the synthetic externs)
//
// Java: ClosureRewriteModule.declareSyntheticExterns()
func (m *ClosureRewriteModule) DeclareSyntheticExterns() {
	var vars []rhino.Node
	for lhs := range m.syntheticExterns.Values() {
		// Skip roots of goog.provide or goog.module.declareLegacyNamespace();
		if m.IsNameInGlobalScope(lhs.GetString()) {
			continue
		}
		vars = append(vars,
			rhino.IR_Var(m.astFactory.CreateNameWithUnknownType(lhs.GetString())).SrcrefTree(lhs))
	}

	if len(vars) == 0 {
		return
	}

	root := m.compiler.GetSynthesizedExternsInput().GetAstRoot(m.compiler)
	for _, v := range vars {
		root.AddChildToBack(v)
	}
}

// IsNameInGlobalScope returns whether the name is declared in the global
// scope, either explicitly with var/const/let or implicitly by a
// goog.provide or legacy goog.module.
//
// Java: ClosureRewriteModule.isNameInGlobalScope(String)
func (m *ClosureRewriteModule) IsNameInGlobalScope(name string) bool {
	_, isLegacy := m.legacyScriptNamespacesAndPrefixes[name]
	return isLegacy || m.globalScope.GetVar(name) != nil
}

// PreprocessExportDeclaration rewrites object literal exports to the
// standard named exports style. i.e. exports = {Foo, Bar} to exports.Foo =
// Foo; exports.Bar = Bar; This makes the module exports into a more standard
// format for later passes.
//
// Java: ClosureRewriteModule.preprocessExportDeclaration(NodeTraversal, Node)
func (m *ClosureRewriteModule) PreprocessExportDeclaration(t *NodeTraversal, n rhino.Node) {
	if !closureRewriteModule_isGoogModuleExportsRef(t, n) {
		return
	}
	assign, ok := closureRewriteModule_assignToTarget(n)
	if !ok || !n.GetGrandparent().IsExprResult() {
		return
	}

	jutil.CheckState(m.currentScript.defaultExport == nil)
	exportRhs := assign.Value().Node()
	if NodeUtil_IsNamedExportsLiteral(exportRhs) {
		insertionPoint := n.GetGrandparent()
		// A named exports literal holds string keys: `exports = {Foo, Bar: Baz}`.
		for member := range ast.MustObjectLit(exportRhs).Members() {
			key := ast.MustStringKey(member.Node())
			exportName := key.Key()
			jsdoc := key.Node().GetJSDocInfo()
			rhs := key.Value().Detach().Node()
			lhs :=
				m.astFactory.
					CreateGetProp(m.astFactory.CreateName("exports", AstFactory_Type_(n)), exportName, AstFactory_Type_(rhs)).
					SrcrefTree(key.Node())
			newExport :=
				rhino.IR_ExprResult(m.astFactory.CreateAssign(lhs, rhs).Srcref(key.Node()).SetJSDocInfo(jsdoc)).
					Srcref(key.Node())
			newExport.InsertAfter(insertionPoint)
			insertionPoint = newExport
		}
		n.GetGrandparent().Detach()
	}
}

// Java: ClosureRewriteModule.recordModuleBody(Node)
func (m *ClosureRewriteModule) RecordModuleBody(moduleRoot rhino.Node) {
	m.PushScript(NewClosureRewriteModule_ScriptDescription())

	m.currentScript.rootNode = moduleRoot
	m.currentScript.isModule = true
}

// Java: ClosureRewriteModule.recordGoogModule(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogModule(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_MODULE_ID_ARG)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	m.currentScript.namespaceId = &namespaceId
	contentsPrefix := closureRewriteModule_toModuleContentsPrefix(namespaceId)
	m.currentScript.contentsPrefix = &contentsPrefix

	m.rewriteState.scriptDescriptionsByGoogModuleNamespace[namespaceId] = m.currentScript
}

// Java: ClosureRewriteModule.recordGoogDeclareLegacyNamespace()
func (m *ClosureRewriteModule) RecordGoogDeclareLegacyNamespace() {
	m.currentScript.declareLegacyNamespace = true
	// Go-only: language-adapter — namespaceId is a nullable String; a nil one
	// panics here as Java's NullPointerException does.
	m.UpdateLegacyScriptNamespacesAndPrefixes(*m.currentScript.namespaceId)
}

// Java: ClosureRewriteModule.updateLegacyScriptNamespacesAndPrefixes(String)
func (m *ClosureRewriteModule) UpdateLegacyScriptNamespacesAndPrefixes(namespace string) {
	m.legacyScriptNamespacesAndPrefixes[namespace] = struct{}{}
	for dot := strings.LastIndexByte(namespace, '.'); dot != -1; dot = strings.LastIndexByte(namespace, '.') {
		namespace = namespace[:dot]
		m.legacyScriptNamespacesAndPrefixes[namespace] = struct{}{}
	}
}

// Java: ClosureRewriteModule.recordGoogProvide(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogProvide(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	if m.currentScript.isModule {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_PROVIDE_CALL)
	}

	// Log legacy namespaces and prefixes.
	m.rewriteState.providedNamespaces[namespaceId] = struct{}{}
	m.UpdateLegacyScriptNamespacesAndPrefixes(namespaceId)
}

// Java: ClosureRewriteModule.recordGoogRequire(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogRequire(t *NodeTraversal, call ast.Call) {
	closureRewriteModule_maybeSplitMultiVar(call.Node())

	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	// Maybe report an error if there is an attempt to import something that
	// is expected to be a goog.module() but no such goog.module() has been
	// defined.
	targetIsAModule := m.rewriteState.ContainsModule(namespaceId)
	_, targetIsALegacyScript := m.rewriteState.providedNamespaces[namespaceId]
	if m.currentScript.isModule && !targetIsAModule && !targetIsALegacyScript {
		m.unrecognizedRequires = append(m.unrecognizedRequires, ClosureRewriteModule_UnrecognizedRequire{requireNode: call.Node(), namespaceId: namespaceId})
	}
}

// Java: ClosureRewriteModule.recordGoogRequireType(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogRequireType(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_TYPE_NAMESPACE)
		return
	}

	// For purposes of import collection, goog.requireType is the same as
	// goog.require but a goog.requireType call is not required to appear
	// after the corresponding namespace definition.
	m.RecordGoogRequire(t, call)
}

// Java: ClosureRewriteModule.recordGoogForwardDeclare(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogForwardDeclare(t *NodeTraversal, call ast.Call) {
	namespaceNode := closureRewriteModule_namespaceArg(call)
	if call.NumArgs() != 1 || !namespaceNode.IsStringLit() {
		t.Report(namespaceNode, ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE)
		return
	}

	// For purposes of import collection, goog.forwardDeclare is the same as
	// goog.require.
	m.RecordGoogRequire(t, call)
}

// Java: ClosureRewriteModule.recordGoogRequireDynamic(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogRequireDynamic(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC)
		return
	}

	namespaceId := namespaceIdNode.GetString()

	if !m.rewriteState.ContainsModule(namespaceId) {
		m.unrecognizedRequires = append(m.unrecognizedRequires, ClosureRewriteModule_UnrecognizedRequire{requireNode: call.Node(), namespaceId: namespaceId})
	}
	m.googRequireDynamicCalls = append(m.googRequireDynamicCalls, call)
}

// Java: ClosureRewriteModule.recordGoogModuleGet(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogModuleGet(t *NodeTraversal, call ast.Call) {
	namespaceIdNode := closureRewriteModule_namespaceArg(call)
	if call.NumArgs() != 1 || !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	if !m.rewriteState.ContainsModule(namespaceId) {
		m.unrecognizedRequires = append(m.unrecognizedRequires, ClosureRewriteModule_UnrecognizedRequire{requireNode: call.Node(), namespaceId: namespaceId})
	}
	m.googModuleGetCalls = append(m.googModuleGetCalls, call)

	// Is the call filling an alias: `alias = goog.module.get('a.namespace');`?
	maybeAssign := call.Node().GetParent()
	assign, ok := ast.AsAssign(maybeAssign)
	if !ok || !maybeAssign.GetParent().IsExprResult() || !m.currentScript.isModule {
		return
	}
	alias, ok := assign.Target().AsName()
	if !ok {
		return
	}

	aliasName := alias.Ident()

	// If the assignment isn't into a var in our scope then it's not ok.
	aliasVar := t.GetScope().GetVar(aliasName)
	if aliasVar == nil {
		// Reported in CheckClosureImports
		return
	}

	// Even if it was to a var in our scope it should still only rewrite if
	// the var looked like:
	//   let x = goog.forwardDeclare('a.namespace');
	aliasVarNodeRhs := NodeUtil_GetRValueOfLValue(aliasVar.GetNode())
	if aliasVarNodeRhs.IsNil() ||
		!NodeUtil_IsCallToNode(aliasVarNodeRhs, closureRewriteModule_GOOG_FORWARDDECLARE) ||
		namespaceId != closureRewriteModule_namespaceArg(ast.MustCall(aliasVarNodeRhs)).GetString() {
		// Reported in CheckClosureImports
		return
	}

	if m.IsUnrequiredWeakNamespace(call.Node(), namespaceId) {
		return
	}

	// Each goog.module.get() calling filling an alias will have the alias
	// importing logic handled at the goog.forwardDeclare call, and the
	// corresponding goog.module.get can simply be removed.
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

// Java: ClosureRewriteModule.recordTopLevelClassOrFunctionName(Node)
func (m *ClosureRewriteModule) RecordTopLevelClassOrFunctionName(classOrFunctionNode rhino.Node) {
	var nameNode ast.Name
	var ok bool
	if fn, isFunction := ast.AsFunction(classOrFunctionNode); isFunction {
		nameNode, ok = fn.Name()
	} else {
		nameNode, ok = ast.MustClass(classOrFunctionNode).Name()
	}
	if ok {
		name := nameNode.Ident()
		if m.currentScript.topLevelNames == nil {
			m.currentScript.topLevelNames = make(map[string]struct{})
		}
		m.currentScript.topLevelNames[name] = struct{}{}
	}
}

// Java: ClosureRewriteModule.recordTopLevelVarNames(Node)
func (m *ClosureRewriteModule) RecordTopLevelVarNames(varNode rhino.Node) {
	NodeUtil_VisitLhsNodesInNode(
		varNode, func(lhs rhino.Node) {
			if m.currentScript.topLevelNames == nil {
				m.currentScript.topLevelNames = make(map[string]struct{})
			}
			m.currentScript.topLevelNames[lhs.GetString()] = struct{}{}
		})
}

// Java: ClosureRewriteModule.maybeRecordExportDeclaration(NodeTraversal, Node)
func (m *ClosureRewriteModule) MaybeRecordExportDeclaration(t *NodeTraversal, n rhino.Node) {
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
	if m.currentScript.defaultExport != nil {
		previousExport := m.currentScript.defaultExport
		localName := previousExport.GetLocalName()
		if localName != nil {
			delete(m.currentScript.namesToInlineByAlias, *localName)
		}
		m.currentScript.defaultExportLocalName = nil
	}
	exportRhs := assign.Value().Node()

	// Exports object should have already been converted in ScriptPreprocess
	// step.
	jutil.CheckStateMsg(
		!NodeUtil_IsNamedExportsLiteral(exportRhs),
		"Exports object should have been converted already")

	m.currentScript.willCreateExportsObject = true
	defaultExport := ClosureRewriteModule_ExportDefinition_NewDefaultExport(t, exportRhs)
	m.currentScript.defaultExport = defaultExport
	if !m.currentScript.declareLegacyNamespace &&
		defaultExport.HasInlinableName(&m.currentScript.exportsToInline) {
		localName := defaultExport.GetLocalName()
		m.currentScript.defaultExportLocalName = localName
		m.RecordExportToInline(defaultExport)
	}
}
