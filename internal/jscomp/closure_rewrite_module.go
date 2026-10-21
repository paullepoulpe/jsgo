// Port of com/google/javascript/jscomp/ClosureRewriteModule.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	syntheticExterns        *jutil.LinkedHashMap[string, *rhino.Node]

	globalScope *Scope // non-final because it must be set after process() is called

	// Per script state needed for rewriting.
	//
	// Go-only: language-adapter — the ArrayDeque is an ArrayList whose first
	// element is the deque's first.
	scriptStack   *jutil.ArrayList[*ClosureRewriteModule_ScriptDescription]
	currentScript *ClosureRewriteModule_ScriptDescription

	rewriteState *ClosureRewriteModule_GlobalRewriteState
	// All prefix namespaces from goog.provides and legacy goog.modules.
	legacyScriptNamespacesAndPrefixes *jutil.LinkedHashSet[string]
	unrecognizedRequires              *jutil.ArrayList[*ClosureRewriteModule_UnrecognizedRequire]
	googModuleGetCalls                *jutil.ArrayList[*rhino.Node]
	googRequireDynamicCalls           *jutil.ArrayList[*rhino.Node]

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

// Java: ClosureRewriteModule.AddAt.values()
func ClosureRewriteModule_AddAt_values() []ClosureRewriteModule_AddAt {
	return []ClosureRewriteModule_AddAt{ClosureRewriteModule_AddAt_BEFORE, ClosureRewriteModule_AddAt_AFTER}
}

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

// Java: ClosureRewriteModule.ScopeType.values()
func ClosureRewriteModule_ScopeType_values() []ClosureRewriteModule_ScopeType {
	return []ClosureRewriteModule_ScopeType{ClosureRewriteModule_ScopeType_EXEC_CONTEXT, ClosureRewriteModule_ScopeType_BLOCK}
}

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
	requireNode *rhino.Node
	namespaceId string
}

// Java: ClosureRewriteModule.UnrecognizedRequire(Node, String)
func NewClosureRewriteModule_UnrecognizedRequire(requireNode *rhino.Node, namespaceId string) *ClosureRewriteModule_UnrecognizedRequire {
	return &ClosureRewriteModule_UnrecognizedRequire{requireNode: requireNode, namespaceId: namespaceId}
}

// ClosureRewriteModule_ExportDefinition is
// ClosureRewriteModule.ExportDefinition.
type ClosureRewriteModule_ExportDefinition struct {
	// Null if the export is a default export (exports = expr)
	exportName *string
	// Null if the export is of a @typedef
	rhs *rhino.Node
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
	if e.rhs != nil {
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

var closureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS = jutil.NewLinkedHashSet(
	rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET, rhino.Token_FUNCTION, rhino.Token_CLASS)

// Java: ClosureRewriteModule.ExportDefinition.newDefaultExport(NodeTraversal, Node)
func ClosureRewriteModule_ExportDefinition_NewDefaultExport(t *NodeTraversal, rhs *rhino.Node) *ClosureRewriteModule_ExportDefinition {
	return ClosureRewriteModule_ExportDefinition_NewNamedExport(t, nil, rhs)
}

// Java: ClosureRewriteModule.ExportDefinition.newNamedExport(NodeTraversal, String, Node)
func ClosureRewriteModule_ExportDefinition_NewNamedExport(t *NodeTraversal, name *string, rhs *rhino.Node) *ClosureRewriteModule_ExportDefinition {
	newExport := &ClosureRewriteModule_ExportDefinition{}
	newExport.exportName = name
	newExport.rhs = rhs
	if rhs != nil && (rhs.IsName() || rhs.IsStringKey()) {
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
		!closureRewriteModule_ExportDefinition_INLINABLE_NAME_PARENTS.Contains(e.nameDecl.GetParentNode().GetToken()) ||
		NodeUtil_IsFunctionDeclaration(e.nameDecl.GetParentNode()) {
		return false
	}
	initialValue := e.nameDecl.GetInitialValue()
	if initialValue == nil || !initialValue.IsCall() {
		return true
	}
	method := initialValue.GetFirstChild()
	if !method.IsGetProp() {
		return true
	}
	maybeGoog := method.GetFirstChild()
	if !maybeGoog.IsName() || maybeGoog.GetString() != "goog" {
		return true
	}
	name := method.GetString()
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
	namespaceId            *string                      // "a.b.c"
	contentsPrefix         *string                      // "module$contents$a$b$c_
	topLevelNames          *jutil.LinkedHashSet[string] // For prefixed content renaming.
	// Go-only: language-adapter — the ArrayDeque is an ArrayList.
	childScripts         *jutil.ArrayList[*ClosureRewriteModule_ScriptDescription]
	namesToInlineByAlias *jutil.LinkedHashMap[string, *ClosureRewriteModule_AliasName] // For alias inlining.

	// Transient state.
	willCreateExportsObject bool

	hasCreatedExportObject bool
	defaultExport          *ClosureRewriteModule_ExportDefinition
	defaultExportLocalName *string
	namedExports           *jutil.LinkedHashSet[string]
	exportsToInline        *jutil.LinkedHashMap[*Var, *ClosureRewriteModule_ExportDefinition]

	// The root of the module. The MODULE_BODY node that contains the module
	// contents. For recognizing top level names.
	rootNode *rhino.Node
}

// Java: ClosureRewriteModule.ScriptDescription()
func NewClosureRewriteModule_ScriptDescription() *ClosureRewriteModule_ScriptDescription {
	return &ClosureRewriteModule_ScriptDescription{
		topLevelNames:        jutil.NewLinkedHashSet[string](),
		childScripts:         jutil.NewArrayList[*ClosureRewriteModule_ScriptDescription](),
		namesToInlineByAlias: jutil.NewLinkedHashMap[string, *ClosureRewriteModule_AliasName](),
		namedExports:         jutil.NewLinkedHashSet[string](),
		exportsToInline:      jutil.NewLinkedHashMap[*Var, *ClosureRewriteModule_ExportDefinition](),
	}
}

// Java: ClosureRewriteModule.ScriptDescription.addChildScript(ScriptDescription)
func (s *ClosureRewriteModule_ScriptDescription) AddChildScript(childScript *ClosureRewriteModule_ScriptDescription) {
	s.childScripts.Add(childScript)
}

// Java: ClosureRewriteModule.ScriptDescription.removeFirstChildScript()
func (s *ClosureRewriteModule_ScriptDescription) RemoveFirstChildScript() *ClosureRewriteModule_ScriptDescription {
	if s.childScripts.IsEmpty() {
		panic(jutil.NewNoSuchElementException())
	}
	return s.childScripts.RemoveAt(0)
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
func (p *ClosureRewriteModule_ScriptPreprocessor) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
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
func (r *ClosureRewriteModule_ScriptRecorder) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	m := r.outer
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
		m.RecordModuleBody(n)
	case rhino.Token_CALL:
		method := n.GetFirstChild()
		if !method.IsGetProp() {
			break
		}
		if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE) {
			m.RecordGoogModule(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE) {
			m.RecordGoogDeclareLegacyNamespace()
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_PROVIDE) {
			m.RecordGoogProvide(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRE) {
			m.RecordGoogRequire(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRETYPE) {
			m.RecordGoogRequireType(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIREDYNAMIC) {
			m.RecordGoogRequireDynamic(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_FORWARDDECLARE) && !parent.IsExprResult() {
			m.RecordGoogForwardDeclare(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_GET) {
			m.RecordGoogModuleGet(t, n)
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
		if m.currentScript.isModule && closureRewriteModule_isExportPropertyAssignment(t, n) {
			m.RecordExportsPropertyAssignment(t, n)
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
func (r *ClosureRewriteModule_ScriptRecorder) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsModuleBody() {
		r.outer.PopScript()
	}
}

// ClosureRewriteModule_ScriptUpdater is ClosureRewriteModule.ScriptUpdater
// (an inner class).
type ClosureRewriteModule_ScriptUpdater struct {
	outer              *ClosureRewriteModule
	scriptDescriptions *jutil.ArrayList[*ClosureRewriteModule_ScriptDescription]
}

// Java: ClosureRewriteModule.ScriptUpdater(Deque<ScriptDescription>)
func (m *ClosureRewriteModule) NewScriptUpdater(
	scriptDescriptions *jutil.ArrayList[*ClosureRewriteModule_ScriptDescription]) *ClosureRewriteModule_ScriptUpdater {
	return &ClosureRewriteModule_ScriptUpdater{outer: m, scriptDescriptions: scriptDescriptions}
}

// Java: ClosureRewriteModule.ScriptUpdater.shouldTraverse(NodeTraversal, Node, Node)
func (u *ClosureRewriteModule_ScriptUpdater) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	m := u.outer

	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if u.scriptDescriptions.IsEmpty() {
			panic(jutil.NewNoSuchElementException())
		}
		currentDescription := u.scriptDescriptions.RemoveAt(0)
		jutil.CheckState(currentDescription.rootNode == n)
		if n.IsFromExterns() && !NodeUtil_IsFromTypeSummary(n) {
			return false
		}

		jutil.CheckState(m.scriptStack.IsEmpty())
		m.PushScript(currentDescription)

		// Capture the scope before doing any rewriting to the scope.
		t.GetScope()

		// Capture the global scope for later reference.
		if m.globalScope == nil {
			m.globalScope = t.GetScope().GetGlobalScope()
		}
	case rhino.Token_MODULE_BODY:
		if parent.GetBooleanProp(rhino.Node_GOOG_MODULE) {
			m.UpdateModuleBodyEarly(n)
		} else {
			return false
		}
	case rhino.Token_CALL:
		method := n.GetFirstChild()
		if !method.IsGetProp() {
			break
		}
		if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE) {
			m.UpdateGoogModule(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE) {
			closureRewriteModule_updateGoogDeclareLegacyNamespace(n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRE) ||
			method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRETYPE) {
			m.UpdateGoogRequire(t, n)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_FORWARDDECLARE) && !parent.IsExprResult() {
			m.UpdateGoogForwardDeclare(t, n)
		} else if closureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING.Matches(method) {
			closureRewriteModule_updateGoogPreventModuleExportsSealing(n)
		}
	case rhino.Token_GETPROP:
		if m.currentScript.isModule && closureRewriteModule_isExportPropertyAssignment(t, n) {
			m.UpdateExportsPropertyAssignment(n, t)
		}
	default:
	}

	if info := n.GetJSDocInfo(); info != nil {
		if typeNodes := info.GetTypeNodes(); len(typeNodes) > 0 {
			m.rewriteJsdocWithTraversal(typeNodes, t)
		}
	}

	return true
}

// Java: ClosureRewriteModule.ScriptUpdater.visit(NodeTraversal, Node, Node)
func (u *ClosureRewriteModule_ScriptUpdater) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	m := u.outer
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
		m.UpdateModuleBody(n)
	case rhino.Token_NAME:
		if m.currentScript.isModule {
			m.MaybeUpdateTopLevelName(t, n)
			m.MaybeUpdateExportDeclaration(t, n)
			t.GetScope() // Creating the scope here has load-bearing side-effects.
			m.MaybeUpdateExportNameRef(t, n)
		}
	case rhino.Token_SCRIPT:
		jutil.CheckState(m.currentScript.rootNode == n)
		m.PopScript()
	default:
	}
}

// RewriteJsdoc rewrites JsDoc type references to match AST changes
// resulting from imported alias inlining, module content renaming of top
// level constructor functions and classes, and module renaming from fully
// qualified legacy namespace to its binary name.
//
// Java: ClosureRewriteModule.rewriteJsdoc(JSDocInfo, Scope)
func (m *ClosureRewriteModule) RewriteJsdoc(info *rhino.JSDocInfo, scope *Scope) {
	replacer := m.NewReplaceJsDocRefs(scope)
	for _, typeNode := range info.GetTypeNodes() {
		NodeUtil_VisitPreOrder(typeNode, replacer)
	}
}

func (m *ClosureRewriteModule) rewriteJsdocWithTraversal(typeNodes []*rhino.Node, t *NodeTraversal) {
	replacer := ClosureRewriteModule_ReplaceJsDocRefs{outer: m, t: t}
	for _, typeNode := range typeNodes {
		NodeUtil_VisitPreOrder(typeNode, &replacer)
	}
}

// ClosureRewriteModule_ReplaceJsDocRefs rewrites JsDoc type references to
// match AST changes resulting from imported alias inlining, module content
// renaming of top level constructor functions and classes, and module
// renaming from fully qualified legacy namespace to its binary name.
//
// It is an inner class: outer is the enclosing ClosureRewriteModule.
type ClosureRewriteModule_ReplaceJsDocRefs struct {
	outer *ClosureRewriteModule
	scope *Scope
	t     *NodeTraversal
}

var _ NodeUtil_Visitor = (*ClosureRewriteModule_ReplaceJsDocRefs)(nil)

// Java: ClosureRewriteModule.ReplaceJsDocRefs(Scope)
func (m *ClosureRewriteModule) NewReplaceJsDocRefs(scope *Scope) *ClosureRewriteModule_ReplaceJsDocRefs {
	return &ClosureRewriteModule_ReplaceJsDocRefs{outer: m, scope: scope}
}

func (r *ClosureRewriteModule_ReplaceJsDocRefs) isLocalShadowVar(rootOfType string) bool {
	var rootVar *Var
	if r.t != nil {
		if r.t.InGlobalScope() {
			return false
		}
		rootVar = r.t.LookupVar(rootOfType)
	} else if r.scope != nil {
		rootVar = r.scope.GetVar(rootOfType)
	}
	return rootVar != nil &&
		!rootVar.GetScope().(*Scope).GetClosestHoistScope().IsGlobal() &&
		!rootVar.GetScope().(*Scope).IsModuleScope()
}

// Java: ClosureRewriteModule.ReplaceJsDocRefs.visit(Node)
func (r *ClosureRewriteModule_ReplaceJsDocRefs) Visit(typeRefNode *rhino.Node) {
	m := r.outer
	if !typeRefNode.IsStringLit() {
		return
	}
	// A type name that might be simple like "Foo" or qualified like "foo.Bar".
	typeName := typeRefNode.GetString()
	dot := strings.IndexByte(typeName, '.')
	rootOfType := typeName
	if dot != -1 {
		rootOfType = typeName[:dot]
	}

	// Rewrite the type node if any of the following hold, in priority order:
	//  - the root of the type name is in our set of aliases to inline
	//  - the root of the type name is a name defined in the module scope
	//  - a prefix of the type name matches a Closure namespace
	// AND the following is false:
	//  - the root of the type name is defined in an inner scope

	// If the name is an alias for an imported namespace rewrite from
	// "{Foo}" to
	// "{module$exports$bar$Foo}" or
	// "{bar.Foo}"
	if m.currentScript.namesToInlineByAlias.ContainsKey(rootOfType) {
		if r.isLocalShadowVar(rootOfType) {
			return
		}
		if m.preprocessorSymbolTable != nil {
			// Jsdoc type node is a single STRING node that spans the whole
			// type. For example STRING node "bar.Foo". When rewriting modules
			// potentially replace only "module" part of the type: "bar.Foo" =>
			// "module$exports$bar$Foo". So we need to remember that "bar" as
			// alias. To do that we clone type node and make "bar" node from it.
			moduleOnlyNode := typeRefNode.CloneNode()
			m.SafeSetString(moduleOnlyNode, rootOfType)
			moduleOnlyNode.SetLength(jutil.StringLength(rootOfType))
			m.MaybeAddAliasToSymbolTable(moduleOnlyNode, m.currentScript.namespaceId)
		}

		aliasedNamespace := m.currentScript.namesToInlineByAlias.Get(rootOfType).newName
		remainder := ""
		if dot != -1 {
			remainder = typeName[dot:]
		}
		m.SafeSetString(typeRefNode, aliasedNamespace+remainder)
	} else if m.currentScript.isModule && m.currentScript.topLevelNames.Contains(rootOfType) {
		if r.isLocalShadowVar(rootOfType) {
			return
		}
		// If this is a module and the type name is the name of a top level
		// var/function/class defined in this script then that var will have
		// been previously renamed from Foo to module$contents$Foo_Foo. Update
		// the JsDoc reference to match.
		m.SafeSetString(typeRefNode, jutil.ToJavaString(m.currentScript.contentsPrefix)+typeName)
	} else if m.currentScript.isModule && rootOfType == "exports" {
		if r.isLocalShadowVar(rootOfType) {
			return
		}
		// rewrite a type reference to the implicit "exports" variable
		// e.g. /** @type {exports.Foo} */ -> /** @type {module$exports$my$mod.Foo} */
		namespace := m.currentScript.GetBinaryNamespace()
		remainder := ""
		if dot != -1 {
			remainder = typeName[dot:]
		}
		m.SafeSetString(typeRefNode, jutil.ToJavaString(namespace)+remainder)
	} else if !m.rewriteState.scriptDescriptionsByGoogModuleNamespace.IsEmpty() {
		r.rewriteIfClosureNamespaceRefForRoot(typeName, rootOfType, typeRefNode)
	}
}

func (r *ClosureRewriteModule_ReplaceJsDocRefs) rewriteIfClosureNamespaceRefForRoot(typeName, rootOfType string, typeRefNode *rhino.Node) {
	m := r.outer
	prefixTypeName := typeName
	suffix := ""

	for {
		binaryNamespaceIfModule := m.rewriteState.GetBinaryNamespace(prefixTypeName)
		if m.legacyScriptNamespacesAndPrefixes.Contains(prefixTypeName) &&
			binaryNamespaceIfModule == nil {
			return
		}
		if binaryNamespaceIfModule != nil {
			if r.isLocalShadowVar(rootOfType) {
				return
			}
			m.SafeSetString(typeRefNode, *binaryNamespaceIfModule+suffix)
			return
		}
		if strings.Contains(prefixTypeName, ".") {
			prefixTypeName = prefixTypeName[:strings.LastIndexByte(prefixTypeName, '.')]
			suffix = typeName[len(prefixTypeName):]
		} else {
			break
		}
	}
}

// RewriteIfClosureNamespaceRef tries to match the longest possible prefix of
// this type to a Closure namespace.
//
// If the longest prefix match is a legacy module or provide, this is a
// no-op. If the longest prefix match is a non-legacy module, this method
// rewrites the type node.
//
// Java: ClosureRewriteModule.ReplaceJsDocRefs.rewriteIfClosureNamespaceRef(String, Node)
func (r *ClosureRewriteModule_ReplaceJsDocRefs) RewriteIfClosureNamespaceRef(typeName string, typeRefNode *rhino.Node) {
	m := r.outer
	// Tries to rename progressively shorter type prefixes like "foo.Bar.Baz",
	// then "foo.Bar", then "foo".
	prefixTypeName := typeName
	suffix := ""

	for {
		binaryNamespaceIfModule := m.rewriteState.GetBinaryNamespace(prefixTypeName)
		if m.legacyScriptNamespacesAndPrefixes.Contains(prefixTypeName) &&
			binaryNamespaceIfModule == nil {
			// This thing is definitely coming from a legacy script and so the
			// fully qualified type name will always resolve as is.
			return
		}

		// If the typeName is a reference to a fully qualified legacy namespace
		// like "foo.bar.Baz" of something that is actually a module then
		// rewrite the JsDoc reference to "module$exports$Bar".
		// Note: we may want to ban this pattern in the future. See b/133501660.
		if binaryNamespaceIfModule != nil {
			m.SafeSetString(typeRefNode, *binaryNamespaceIfModule+suffix)
			return
		}

		if strings.Contains(prefixTypeName, ".") {
			prefixTypeName = prefixTypeName[:strings.LastIndexByte(prefixTypeName, '.')]
			suffix = typeName[len(prefixTypeName):]
		} else {
			break
		}
	}
}

// ClosureRewriteModule_GlobalRewriteState is the global state tracking an
// association between the dotted names of goog.module()s and whether the
// goog.module declares itself as a legacy namespace. Allows for detecting
// duplicate goog.module()s and for rewriting fully qualified JsDoc type
// references to goog.module() types in legacy scripts.
type ClosureRewriteModule_GlobalRewriteState struct {
	scriptDescriptionsByGoogModuleNamespace *jutil.LinkedHashMap[string, *ClosureRewriteModule_ScriptDescription]
	// Go-only: language-adapter — a Guava HashMultimap with Node keys that is
	// only written (put), never read, so an insertion-ordered multimap.
	namespaceIdsByScriptNode *jutil.SetMultimap[*rhino.Node, string]
	providedNamespaces       *jutil.LinkedHashSet[string]
}

// Java: ClosureRewriteModule.GlobalRewriteState()
func NewClosureRewriteModule_GlobalRewriteState() *ClosureRewriteModule_GlobalRewriteState {
	return &ClosureRewriteModule_GlobalRewriteState{
		scriptDescriptionsByGoogModuleNamespace: jutil.NewLinkedHashMap[string, *ClosureRewriteModule_ScriptDescription](),
		namespaceIdsByScriptNode:                jutil.NewSetMultimap[*rhino.Node, string](),
		providedNamespaces:                      jutil.NewLinkedHashSet[string](),
	}
}

// Java: ClosureRewriteModule.GlobalRewriteState.containsModule(String)
func (g *ClosureRewriteModule_GlobalRewriteState) ContainsModule(namespaceId string) bool {
	return g.scriptDescriptionsByGoogModuleNamespace.ContainsKey(namespaceId)
}

// Java: ClosureRewriteModule.GlobalRewriteState.isLegacyModule(String)
func (g *ClosureRewriteModule_GlobalRewriteState) IsLegacyModule(namespaceId string) bool {
	jutil.CheckArgument(g.ContainsModule(namespaceId))
	return g.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId).declareLegacyNamespace
}

// Java: ClosureRewriteModule.GlobalRewriteState.getBinaryNamespace(String)
func (g *ClosureRewriteModule_GlobalRewriteState) GetBinaryNamespace(namespaceId string) *string {
	script := g.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId)
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
	googModule := g.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId)
	if googModule == nil {
		return nil
	}
	return googModule.rootNode.GetJSType()
}

// Java: ClosureRewriteModule.GlobalRewriteState.getExportedNamespaceOrScript(String)
func (g *ClosureRewriteModule_GlobalRewriteState) GetExportedNamespaceOrScript(namespaceId string) *string {
	if g.providedNamespaces.Contains(namespaceId) {
		return &namespaceId
	}
	script := g.scriptDescriptionsByGoogModuleNamespace.Get(namespaceId)
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
		syntheticExterns:                  jutil.NewLinkedHashMap[string, *rhino.Node](),
		scriptStack:                       jutil.NewArrayList[*ClosureRewriteModule_ScriptDescription](),
		rewriteState:                      NewClosureRewriteModule_GlobalRewriteState(),
		legacyScriptNamespacesAndPrefixes: jutil.NewLinkedHashSet[string](),
		unrecognizedRequires:              jutil.NewArrayList[*ClosureRewriteModule_UnrecognizedRequire](),
		googModuleGetCalls:                jutil.NewArrayList[*rhino.Node](),
		googRequireDynamicCalls:           jutil.NewArrayList[*rhino.Node](),
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
func (m *ClosureRewriteModule) IsUnrequiredWeakNamespace(call *rhino.Node, namespaceId string) bool {
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
func (u *ClosureRewriteModule_UnwrapGoogLoadModule) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	compiler := u.outer.compiler
	switch n.GetToken() {
	case rhino.Token_ROOT, rhino.Token_SCRIPT:
		return true
	case rhino.Token_EXPR_RESULT:
		call := n.GetFirstChild()
		if NodeUtil_IsCallToNode(call, closureRewriteModule_GOOG_LOADMODULE) && call.GetLastChild().IsFunction() {
			parent.PutBooleanProp(rhino.Node_GOOG_MODULE, true)
			functionNode := call.GetLastChild()
			compiler.ReportFunctionDeleted(functionNode)
			moduleBody := functionNode.GetLastChild().Detach()
			moduleBody.SetToken(rhino.Token_MODULE_BODY)
			exportsParameter := NodeUtil_GetFunctionParameters(functionNode).GetOnlyChild()
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
func (m *ClosureRewriteModule) Process(externs *rhino.Node, root *rhino.Node) {
	// Record all the scripts first so that the googModuleNamespaces global
	// state can be complete before doing any updating also queue up
	// scriptDescriptions for later use in ScriptUpdater runs.

	sc := NewSyntacticScopeCreator(m.compiler)
	sc.EnableArenaPool()
	defer func() {
		m.globalScope = nil
		sc.ReleaseArena()
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

	scriptDescriptions := jutil.NewArrayList[*ClosureRewriteModule_ScriptDescription]()
	for _, parent := range []*rhino.Node{externs, root} {
		for script := parent.GetFirstChild(); script != nil; script = script.GetNext() {
			jutil.CheckStateMsg1(script.IsScript(), "%s", script)
			NodeTraversal_Traverse(m.compiler, script, &unwrap)

			m.PushScript(NewClosureRewriteModule_ScriptDescription()) // sets currentScript

			m.currentScript.rootNode = script
			scriptDescriptions.Add(m.currentScript)
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
	for call := range m.googModuleGetCalls.All() {
		m.UpdateGoogModuleGetCall(call)
	}
	for call := range m.googRequireDynamicCalls.All() {
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
	var vars []*rhino.Node
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
	return m.legacyScriptNamespacesAndPrefixes.Contains(name) || m.globalScope.GetVar(name) != nil
}

// PreprocessExportDeclaration rewrites object literal exports to the
// standard named exports style. i.e. exports = {Foo, Bar} to exports.Foo =
// Foo; exports.Bar = Bar; This makes the module exports into a more standard
// format for later passes.
//
// Java: ClosureRewriteModule.preprocessExportDeclaration(NodeTraversal, Node)
func (m *ClosureRewriteModule) PreprocessExportDeclaration(t *NodeTraversal, n *rhino.Node) {
	if !closureRewriteModule_isGoogModuleExportsRef(t, n) || !closureRewriteModule_isAssignTarget(n) || !n.GetGrandparent().IsExprResult() {
		return
	}

	jutil.CheckState(m.currentScript.defaultExport == nil)
	exportRhs := n.GetNext()
	if NodeUtil_IsNamedExportsLiteral(exportRhs) {
		insertionPoint := n.GetGrandparent()
		for key := exportRhs.GetFirstChild(); key != nil; key = key.GetNext() {
			exportName := key.GetString()
			jsdoc := key.GetJSDocInfo()
			rhs := key.RemoveFirstChild()
			lhs :=
				m.astFactory.
					CreateGetProp(m.astFactory.CreateName("exports", AstFactory_Type_(n)), exportName, AstFactory_Type_(rhs)).
					SrcrefTree(key)
			newExport :=
				rhino.IR_ExprResult(m.astFactory.CreateAssign(lhs, rhs).Srcref(key).SetJSDocInfo(jsdoc)).
					Srcref(key)
			newExport.InsertAfter(insertionPoint)
			insertionPoint = newExport
		}
		n.GetGrandparent().Detach()
	}
}

// Java: ClosureRewriteModule.recordModuleBody(Node)
func (m *ClosureRewriteModule) RecordModuleBody(moduleRoot *rhino.Node) {
	m.PushScript(NewClosureRewriteModule_ScriptDescription())

	m.currentScript.rootNode = moduleRoot
	m.currentScript.isModule = true
}

// Java: ClosureRewriteModule.recordGoogModule(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogModule(t *NodeTraversal, call *rhino.Node) {
	namespaceIdNode := call.GetLastChild()
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_MODULE_ID_ARG)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	m.currentScript.namespaceId = &namespaceId
	contentsPrefix := closureRewriteModule_toModuleContentsPrefix(namespaceId)
	m.currentScript.contentsPrefix = &contentsPrefix

	scriptNode := NodeUtil_GetEnclosingScript(m.currentScript.rootNode)
	m.rewriteState.scriptDescriptionsByGoogModuleNamespace.Put(namespaceId, m.currentScript)
	m.rewriteState.namespaceIdsByScriptNode.Put(scriptNode, namespaceId)
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
	m.legacyScriptNamespacesAndPrefixes.Add(namespace)
	for dot := strings.LastIndexByte(namespace, '.'); dot != -1; dot = strings.LastIndexByte(namespace, '.') {
		namespace = namespace[:dot]
		m.legacyScriptNamespacesAndPrefixes.Add(namespace)
	}
}

// Java: ClosureRewriteModule.recordGoogProvide(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogProvide(t *NodeTraversal, call *rhino.Node) {
	namespaceIdNode := call.GetLastChild()
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_PROVIDE_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	if m.currentScript.isModule {
		t.Report(namespaceIdNode, ClosureRewriteModule_INVALID_PROVIDE_CALL)
	}

	scriptNode := NodeUtil_GetEnclosingScript(call)
	// Log legacy namespaces and prefixes.
	m.rewriteState.providedNamespaces.Add(namespaceId)
	m.rewriteState.namespaceIdsByScriptNode.Put(scriptNode, namespaceId)
	m.UpdateLegacyScriptNamespacesAndPrefixes(namespaceId)
}

// Java: ClosureRewriteModule.recordGoogRequire(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogRequire(t *NodeTraversal, call *rhino.Node) {
	closureRewriteModule_maybeSplitMultiVar(call)

	namespaceIdNode := call.GetLastChild()
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	// Maybe report an error if there is an attempt to import something that
	// is expected to be a goog.module() but no such goog.module() has been
	// defined.
	targetIsAModule := m.rewriteState.ContainsModule(namespaceId)
	targetIsALegacyScript := m.rewriteState.providedNamespaces.Contains(namespaceId)
	if m.currentScript.isModule && !targetIsAModule && !targetIsALegacyScript {
		m.unrecognizedRequires.Add(NewClosureRewriteModule_UnrecognizedRequire(call, namespaceId))
	}
}

// Java: ClosureRewriteModule.recordGoogRequireType(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogRequireType(t *NodeTraversal, call *rhino.Node) {
	namespaceIdNode := call.GetLastChild()
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
func (m *ClosureRewriteModule) RecordGoogForwardDeclare(t *NodeTraversal, call *rhino.Node) {
	namespaceNode := call.GetLastChild()
	if !call.HasTwoChildren() || !namespaceNode.IsStringLit() {
		t.Report(namespaceNode, ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE)
		return
	}

	// For purposes of import collection, goog.forwardDeclare is the same as
	// goog.require.
	m.RecordGoogRequire(t, call)
}

// Java: ClosureRewriteModule.recordGoogRequireDynamic(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogRequireDynamic(t *NodeTraversal, call *rhino.Node) {
	namespaceIdNode := call.GetLastChild()
	if !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_REQUIRE_DYNAMIC)
		return
	}

	namespaceId := namespaceIdNode.GetString()

	if !m.rewriteState.ContainsModule(namespaceId) {
		m.unrecognizedRequires.Add(NewClosureRewriteModule_UnrecognizedRequire(call, namespaceId))
	}
	m.googRequireDynamicCalls.Add(call)
}

// Java: ClosureRewriteModule.recordGoogModuleGet(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordGoogModuleGet(t *NodeTraversal, call *rhino.Node) {
	namespaceIdNode := call.GetLastChild()
	if !call.HasTwoChildren() || !namespaceIdNode.IsStringLit() {
		t.Report(namespaceIdNode, ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
		return
	}
	namespaceId := namespaceIdNode.GetString()

	if !m.rewriteState.ContainsModule(namespaceId) {
		m.unrecognizedRequires.Add(NewClosureRewriteModule_UnrecognizedRequire(call, namespaceId))
	}
	m.googModuleGetCalls.Add(call)

	maybeAssign := call.GetParent()
	isFillingAnAlias :=
		maybeAssign.IsAssign() &&
			maybeAssign.GetFirstChild().IsName() &&
			maybeAssign.GetParent().IsExprResult()
	if !isFillingAnAlias || !m.currentScript.isModule {
		return
	}

	aliasName := call.GetParent().GetFirstChild().GetString()

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
	if aliasVarNodeRhs == nil ||
		!NodeUtil_IsCallToNode(aliasVarNodeRhs, closureRewriteModule_GOOG_FORWARDDECLARE) ||
		namespaceId != aliasVarNodeRhs.GetLastChild().GetString() {
		// Reported in CheckClosureImports
		return
	}

	if m.IsUnrequiredWeakNamespace(call, namespaceId) {
		return
	}

	// Each goog.module.get() calling filling an alias will have the alias
	// importing logic handled at the goog.forwardDeclare call, and the
	// corresponding goog.module.get can simply be removed.
	m.compiler.ReportChangeToEnclosingScope(maybeAssign)
	maybeAssign.GetParent().Detach()
	m.googModuleGetCalls.RemoveAt(m.googModuleGetCalls.Size() - 1)
}

// Java: ClosureRewriteModule.recordTopLevelClassOrFunctionName(Node)
func (m *ClosureRewriteModule) RecordTopLevelClassOrFunctionName(classOrFunctionNode *rhino.Node) {
	nameNode := classOrFunctionNode.GetFirstChild()
	if nameNode.IsName() && nameNode.GetString() != "" {
		name := nameNode.GetString()
		m.currentScript.topLevelNames.Add(name)
	}
}

// Java: ClosureRewriteModule.recordTopLevelVarNames(Node)
func (m *ClosureRewriteModule) RecordTopLevelVarNames(varNode *rhino.Node) {
	NodeUtil_VisitLhsNodesInNode(
		varNode, func(lhs *rhino.Node) { m.currentScript.topLevelNames.Add(lhs.GetString()) })
}

// Java: ClosureRewriteModule.maybeRecordExportDeclaration(NodeTraversal, Node)
func (m *ClosureRewriteModule) MaybeRecordExportDeclaration(t *NodeTraversal, n *rhino.Node) {
	if !m.currentScript.isModule || !closureRewriteModule_isGoogModuleExportsRef(t, n) || !closureRewriteModule_isAssignTarget(n) {
		return
	}

	// ClosureCheckModule reports an error for duplicate 'exports = '
	// assignments, but that error may be suppressed. If so then use the final
	// assignment as the canonical one.
	if m.currentScript.defaultExport != nil {
		previousExport := m.currentScript.defaultExport
		localName := previousExport.GetLocalName()
		if localName != nil && m.currentScript.namesToInlineByAlias.ContainsKey(*localName) {
			m.currentScript.namesToInlineByAlias.Remove(*localName)
		}
		m.currentScript.defaultExportLocalName = nil
	}
	exportRhs := n.GetNext()

	// Exports object should have already been converted in ScriptPreprocess
	// step.
	jutil.CheckStateMsg(
		!NodeUtil_IsNamedExportsLiteral(exportRhs),
		"Exports object should have been converted already")

	m.currentScript.willCreateExportsObject = true
	defaultExport := ClosureRewriteModule_ExportDefinition_NewDefaultExport(t, exportRhs)
	m.currentScript.defaultExport = defaultExport
	if !m.currentScript.declareLegacyNamespace &&
		defaultExport.HasInlinableName(m.currentScript.exportsToInline) {
		localName := defaultExport.GetLocalName()
		m.currentScript.defaultExportLocalName = localName
		m.RecordExportToInline(defaultExport)
	}
}

// Java: ClosureRewriteModule.updateModuleBodyEarly(Node)
func (m *ClosureRewriteModule) UpdateModuleBodyEarly(moduleScopeRoot *rhino.Node) {
	m.PushScript(m.currentScript.RemoveFirstChildScript())
	m.currentScript.rootNode = moduleScopeRoot
}

// Java: ClosureRewriteModule.updateGoogModule(NodeTraversal, Node)
func (m *ClosureRewriteModule) UpdateGoogModule(t *NodeTraversal, call *rhino.Node) {
	if !m.currentScript.isModule {
		m.compiler.ReportChangeToEnclosingScope(call)
		undefined := m.astFactory.CreateVoid(m.astFactory.CreateNumber(0)).SrcrefTree(call)
		call.ReplaceWith(undefined)
		return
	}

	// If it's a goog.module() with a legacy namespace.
	if m.currentScript.declareLegacyNamespace {
		// Rewrite "goog.module('Foo');" as "goog.provide('Foo');".
		call.GetFirstChild().SetString("provide")
		m.compiler.ReportChangeToEnclosingScope(call)
	}

	// If this script file isn't going to eventually create it's own exports
	// object, then we know we'll need to do it ourselves, and so we might as
	// well create it as early as possible to avoid ordering issues with
	// goog.define().
	if !m.currentScript.willCreateExportsObject {
		jutil.CheckStateMsg1(!m.currentScript.hasCreatedExportObject, "%s", m.currentScript)
		m.ExportTheEmptyBinaryNamespaceAt(NodeUtil_GetEnclosingStatement(call), ClosureRewriteModule_AddAt_AFTER, t)
	}

	if !m.currentScript.declareLegacyNamespace && !m.preserveSugar {
		// Otherwise it's a regular module and the goog.module() line can be
		// removed.
		m.compiler.ReportChangeToEnclosingScope(call)
		NodeUtil_GetEnclosingStatement(call).Detach()
	}
}

// Java: ClosureRewriteModule.updateGoogDeclareLegacyNamespace(Node)
func closureRewriteModule_updateGoogDeclareLegacyNamespace(call *rhino.Node) {
	NodeUtil_GetEnclosingStatement(call).Detach()
}

// Java: ClosureRewriteModule.updateGoogPreventModuleExportsSealing(Node)
func closureRewriteModule_updateGoogPreventModuleExportsSealing(call *rhino.Node) {
	NodeUtil_GetEnclosingStatement(call).Detach()
}

// Java: ClosureRewriteModule.updateGoogRequire(NodeTraversal, Node)
func (m *ClosureRewriteModule) UpdateGoogRequire(t *NodeTraversal, call *rhino.Node) {
	namespaceIdNode := call.GetLastChild()
	statementNode := NodeUtil_GetEnclosingStatement(call)
	namespaceId := namespaceIdNode.GetString()

	targetIsNonLegacyGoogModule :=
		m.rewriteState.ContainsModule(namespaceId) && !m.rewriteState.IsLegacyModule(namespaceId)
	importHasAlias := NodeUtil_IsNameDeclaration(statementNode)
	isDestructuring := statementNode.GetFirstChild().IsDestructuringLhs()

	// If the current script is a module or the require statement has a return
	// value that is stored in an alias then the require is goog.module()
	// style.
	currentScriptIsAModule := m.currentScript.isModule
	// "var Foo = goog.require("bar.Foo");" or "const {Foo} = goog.require('bar');" style.
	requireDirectlyStoredInAlias := NodeUtil_IsNameDeclaration(call.GetGrandparent())
	if currentScriptIsAModule &&
		requireDirectlyStoredInAlias &&
		m.IsTopLevel(t, statementNode, ClosureRewriteModule_ScopeType_EXEC_CONTEXT) {
		// Record alias -> exportedNamespace associations for later inlining.
		lhs := call.GetParent()
		exportedNamespace := m.rewriteState.GetExportedNamespaceOrScript(namespaceId)
		if exportedNamespace == nil {
			// There's nothing to inline. The missing provide/module will be
			// reported elsewhere.
		} else if lhs.IsName() {
			// `var Foo` case
			aliasName := statementNode.GetFirstChild().GetString()
			m.RecordNameToInline(aliasName, *exportedNamespace, &namespaceId)
			m.MaybeAddAliasToSymbolTable(statementNode.GetFirstChild(), m.currentScript.namespaceId)
		} else if lhs.IsDestructuringLhs() && lhs.GetFirstChild().IsObjectPattern() {
			// `const {Foo}` case
			m.MaybeWarnForInvalidDestructuring(t, lhs.GetParent(), namespaceId)
			for importSpec := lhs.GetFirstFirstChild(); importSpec != nil; importSpec = importSpec.GetNext() {
				jutil.CheckStateMsg1(importSpec.HasChildren(), "%s", importSpec)
				importedProperty := importSpec.GetString()
				aliasNode := importSpec.GetFirstChild()
				aliasName := aliasNode.GetString()
				fullName := *exportedNamespace + "." + importedProperty
				m.RecordNameToInline(aliasName, fullName /* namespaceId= */, nil)

				// Record alias before we rename node.
				m.MaybeAddAliasToSymbolTable(aliasNode, m.currentScript.namespaceId)
				// Need to rename node otherwise it will stay global and messes up
				// index if there are other files that use the same destructuring
				// alias.
				m.SafeSetString(aliasNode, jutil.ToJavaString(m.currentScript.contentsPrefix)+aliasName)
			}
		} else {
			panic(jutil.NewRuntimeException("Illegal goog.module import: " + lhs.String()))
		}
	}

	if m.currentScript.isModule || targetIsNonLegacyGoogModule {
		if isDestructuring {
			if !m.preserveSugar {
				// Delete the goog.require() because we're going to inline its
				// alias later.
				m.compiler.ReportChangeToEnclosingScope(statementNode)
				statementNode.Detach()
			}
		} else if targetIsNonLegacyGoogModule {
			jutil.CheckStateMsg1(
				m.IsTopLevel(t, statementNode, ClosureRewriteModule_ScopeType_EXEC_CONTEXT),
				"Unexpected non-top-level require at %s",
				call)
			if importHasAlias || !m.rewriteState.IsLegacyModule(namespaceId) {
				if !m.preserveSugar {
					// Delete the goog.require() because we're going to inline its
					// alias later.
					m.compiler.ReportChangeToEnclosingScope(statementNode)
					statementNode.Detach()
				}
			}
		} else {
			// TODO(bangert): make this compatible with preserveSugar. const B =
			// goog.require('b') runs into problems because the type checker
			// cannot handle const.
			// Rewrite
			//   "var B = goog.require('B');" to
			//   "goog.require('B');"
			// because even though we're going to inline the B alias,
			// ProcessClosurePrimitives is going to want to see this legacy
			// require.
			call.Detach()
			statementNode.ReplaceWith(rhino.IR_ExprResult(call))
			m.compiler.ReportChangeToEnclosingScope(call)
		}
	}
}

// MaybeWarnForInvalidDestructuring: these restrictions are in place to make
// it easier to migrate goog.modules to ES6 modules, by structuring the
// imports/exports in a consistent way.
//
// Java: ClosureRewriteModule.maybeWarnForInvalidDestructuring(NodeTraversal, Node, String)
func (m *ClosureRewriteModule) MaybeWarnForInvalidDestructuring(
	t *NodeTraversal, importNode *rhino.Node, importedNamespace string) {
	jutil.CheckArgumentMsg1(importNode.GetFirstChild().IsDestructuringLhs(), "%s", importNode)
	importedModule :=
		m.rewriteState.scriptDescriptionsByGoogModuleNamespace.Get(importedNamespace)
	if importedModule == nil {
		// Don't know enough to give a good warning here.
		return
	}
	if importedModule.defaultExport != nil {
		t.Report(importNode, ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT)
		return
	}
	objPattern := importNode.GetFirstFirstChild()
	for key := objPattern.GetFirstChild(); key != nil; key = key.GetNext() {
		exportName := key.GetString()
		if !importedModule.namedExports.Contains(exportName) {
			t.Report(importNode, ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED, exportName, importedNamespace)
		}
	}
}

// Java: ClosureRewriteModule.updateGoogForwardDeclare(NodeTraversal, Node)
func (m *ClosureRewriteModule) UpdateGoogForwardDeclare(t *NodeTraversal, call *rhino.Node) {
	namespaceId := call.GetLastChild().GetString()
	if m.IsUnrequiredWeakNamespace(call, namespaceId) {
		// For the delayed goog.module.get pattern where an alias variable is initialized by
		// goog.forwardDeclare('b') and later assigned by goog.module.get('b'):
		// When 'b' is in a weak chunk, the goog.module.get call is not inlined and is rewritten
		// to null. Therefore, the alias declaration must be initialized to null rather than
		// registered for inlining into references of the alias.
		m.compiler.ReportChangeToEnclosingScope(call)
		call.ReplaceWith(m.astFactory.CreateNull().Srcref(call))
		return
	}

	// For non-weak imports, goog.forwardDeclare is handled the same as goog.require:
	// the alias is recorded for inlining and the declaration statement is detached.
	m.UpdateGoogRequire(t, call)
}

// Java: ClosureRewriteModule.updateGoogRequireDynamicCall(Node)
func (m *ClosureRewriteModule) UpdateGoogRequireDynamicCall(call *rhino.Node) {
	parent := call.GetParent()
	if parent == nil {
		// This call has been detached from the AST in unrecognized require
		// handling. Nothing else to do here.
		return
	}
	jutil.CheckStateMsg(
		parent.IsAwait() || (parent.IsGetProp() && parent.GetString() == "then"),
		"goog.requireDynamic() in only allowed in await/then expression")

	if parent.IsAwait() {
		m.UpdateGoogRequireDynamicCallInAwait(call)
	} else {
		m.UpdateGoogRequireDynamicCallInThen(call)
	}
}

// UpdateGoogRequireDynamicCallInThen rewrites
//
//	goog.requireDynamic('a.b.c').then( ({Foo}) => { new Foo().render(); });
//
// to
//
//	goog.importHandler_('h4sh').then(() => {
//	     const {Foo} = module$exports$a$b$c;
//	     new Foo().render();
//	});
//
// Note that rewriting works for both destructuring pattern and name, i.e.,
// `{Foo}` or `foo`.
//
// Java: ClosureRewriteModule.updateGoogRequireDynamicCallInThen(Node)
func (m *ClosureRewriteModule) UpdateGoogRequireDynamicCallInThen(call *rhino.Node) {
	namespaceIdNode := call.GetSecondChild()
	namespaceId := namespaceIdNode.GetString()
	exportedNamespace := m.rewriteState.GetExportedNamespaceOrScript(namespaceId)
	jutil.CheckStateMsg(
		exportedNamespace != nil, "Exported namespace for goog.requireDynamic() canot be null")

	// `goog.requireDynamic('a.b.c').then()`
	getProp := call.GetParent()
	thenCallNode := getProp.GetParent()
	jutil.CheckStateMsg(thenCallNode != nil && thenCallNode.IsCall(), "must be a 'then' call expression")

	// Create a string literal containing the ID of the chunk containing the
	// module we want
	argNode :=
		m.astFactory.CreateString(m.NamespaceIdToXid(namespaceId)).Srcref(namespaceIdNode)

	// `goog.requireDynamic('a.b.c')` -> `goog.requireDynamic('chunkId')`
	namespaceIdNode.ReplaceWith(argNode)

	// `goog.requireDynamic` -> `goog.importHandler_`
	calleeNode := call.GetFirstChild()
	jutil.CheckStateMsg1(calleeNode.MatchesQualifiedName(closureRewriteModule_GOOG_REQUIREDYNAMIC_NAME), "%s", calleeNode)
	calleeNode.SetString(closureRewriteModule_IMPORT_HANDLER_NAME)

	// `module$exports$a$b$c`
	exportedNamespaceNameNode :=
		m.astFactory.CreateNameInScope(m.globalTypedStaticScope(), *exportedNamespace).SrcrefTree(call)
	exportedNamespaceNameNode.SetOriginalName(&namespaceId)

	// Get `{Foo}` from `goog.requireDynamic().then(   ({Foo}) => {  }   )`;
	functionNode := thenCallNode.GetSecondChild()
	jutil.CheckStateMsg(functionNode != nil && functionNode.IsFunction(), "must be a function in `then`")

	// Get param list of the callback function
	paramListNode := functionNode.GetSecondChild()
	jutil.CheckStateMsg(paramListNode.GetChildCount() == 1, "function must have only one parameter")

	// Get the callback function body
	functionBody := paramListNode.GetNext()

	// Get parameter, object pattern or name
	objectPatternOrNameNode := paramListNode.GetOnlyChild()
	jutil.CheckStateMsg(
		objectPatternOrNameNode.IsObjectPattern() || objectPatternOrNameNode.IsName(),
		"parameter of callback function must be object pattern or name")
	objectPatternOrNameNode.Detach()

	// `({Foo}) => {}` -> `() => {}`
	paramListNode.RemoveChildren()

	// Dynamic require is not allowed for JS that is too old to have `const`
	// (< ES6), We don't expect to ever see `--language_in=ES5` in combination
	// with dynamic require.
	languageIn := m.compiler.GetOptions().GetLanguageIn()
	jutil.CheckStateMsg2(
		languageIn.ToFeatureSet().ContainsFeature(parser.FeatureSet_Feature_CONST_DECLARATIONS),
		"'%s' does not contain '%s'",
		languageIn,
		parser.FeatureSet_Feature_CONST_DECLARATIONS)

	enclosingScript := NodeUtil_GetEnclosingScript(call)
	NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, m.compiler)

	// Create `const {Foo} = module$exports$a$b$c;`
	// Or `const foo = module$exports$a$b$c;'
	var declarationNode *rhino.Node
	if objectPatternOrNameNode.IsObjectPattern() {
		declarationNode =
			m.astFactory.
				CreateSingleConstObjectPatternDeclaration(
					objectPatternOrNameNode, exportedNamespaceNameNode).
				SrcrefTreeIfMissing(call)
	} else {
		declarationNode =
			m.astFactory.
				CreateSingleConstNameDeclaration(
					objectPatternOrNameNode.GetString(), exportedNamespaceNameNode).
				SrcrefTreeIfMissing(call)
	}

	// If right hand side of the callback arrow function is an expression
	// instead of BLOCK, e.g., `({Foo}) => Foo` instead of `() => { return Foo;
	// }`, create a BLOCK to host the expression.
	if functionBody.IsBlock() {
		functionBody.AddChildToFront(declarationNode)
		m.compiler.ReportChangeToEnclosingScope(declarationNode)
	} else {
		returnStmt := m.astFactory.CreateReturn(functionBody.Detach())
		newBlock := m.astFactory.CreateBlock(returnStmt).SrcrefTree(call)
		newBlock.InsertAfter(paramListNode)
		newBlock.AddChildToFront(declarationNode)
		m.compiler.ReportChangeToEnclosingScope(newBlock)
	}

	m.compiler.ReportChangeToEnclosingScope(call)
}

// UpdateGoogRequireDynamicCallInAwait rewrites
//
//	const {Foo} = await goog.requireDynamic('a.b.c')`
//
// to
//
//	await goog.importHandler_('tJJovc');
//	const {Foo} = module$exports$a$b$c;
//
// Java: ClosureRewriteModule.updateGoogRequireDynamicCallInAwait(Node)
func (m *ClosureRewriteModule) UpdateGoogRequireDynamicCallInAwait(call *rhino.Node) {
	namespaceIdNode := call.GetSecondChild()
	namespaceId := namespaceIdNode.GetString()
	exportedNamespace := m.rewriteState.GetExportedNamespaceOrScript(namespaceId)
	jutil.CheckStateMsg(
		exportedNamespace != nil, "Exported namespace for goog.requireDynamic() cannot be null")

	// `await goog.requireDynamic('a.b.c')`
	existingAwait := call.GetParent()
	jutil.CheckStateMsg(
		existingAwait.IsAwait(), "Only goog.requireDynamic() in await expression is supported now")

	// `{Foo} = await goog.requireDynamic('a.b.c')`
	awaitParent := existingAwait.GetParent()
	if awaitParent == nil || (!awaitParent.IsDestructuringLhs() && !awaitParent.IsName()) {
		m.compiler.Report(JSError_MakeWithNode(call, ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT))
	}

	// `const {Foo} = await goog.requireDynamic('a.b.c')`
	declarationStatement := awaitParent.GetParent()
	// Reject non-declarations or declarations in a for loop.
	if !NodeUtil_IsNameDeclaration(declarationStatement) ||
		!NodeUtil_IsStatement(declarationStatement) {
		m.compiler.Report(JSError_MakeWithNode(call, ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT))
	}

	// `module$exports$a$b$c`
	exportedNamespaceNameNode :=
		m.astFactory.CreateQName(m.globalTypedStaticScope(), *exportedNamespace).SrcrefTree(call)
	exportedNamespaceNameNode.SetJSType(m.rewriteState.GetGoogModuleNamespaceType(namespaceId))
	exportedNamespaceNameNode.SetOriginalName(&namespaceId)

	// Create a string literal containing the ID of the chunk containing the
	// module we want
	argNode :=
		m.astFactory.CreateString(m.NamespaceIdToXid(namespaceId)).Srcref(namespaceIdNode)

	// `goog.requireDynamic('a.b.c')` -> `goog.requireDynamic('chunkId')`
	namespaceIdNode.ReplaceWith(argNode)

	// `goog.requireDynamic` -> `goog.importHandler_`
	calleeNode := call.GetFirstChild()
	jutil.CheckStateMsg1(calleeNode.MatchesQualifiedName(closureRewriteModule_GOOG_REQUIREDYNAMIC_NAME), "%s", calleeNode)
	calleeNode.SetString(closureRewriteModule_IMPORT_HANDLER_NAME)

	// Replace the await with the module object
	// `const {Foo} = module$exports$a$b$c;`
	existingAwait.ReplaceWith(exportedNamespaceNameNode)

	// Insert the await as a statement before the declaration
	awaitStatement := m.astFactory.ExprResult(existingAwait).Srcref(existingAwait)
	awaitStatement.InsertBefore(declarationStatement)
	m.compiler.ReportChangeToEnclosingScope(call)
}

// Java: ClosureRewriteModule.namespaceIdToXid(String)
func (m *ClosureRewriteModule) NamespaceIdToXid(namespaceId string) string {
	hashFunction := m.compiler.GetOptions().GetChunkIdHashFunction()
	var xid *Xid
	if hashFunction == nil {
		xid = NewXid()
	} else {
		xid = NewXidWithHasher(hashFunction)
	}
	return xid.Get(namespaceId)
}

// Java: ClosureRewriteModule.updateGoogModuleGetCall(Node)
func (m *ClosureRewriteModule) UpdateGoogModuleGetCall(call *rhino.Node) {
	namespaceIdNode := call.GetSecondChild()
	namespaceId := namespaceIdNode.GetString()

	if m.IsUnrequiredWeakNamespace(call, namespaceId) {
		m.compiler.Report(JSError_MakeWithNode(call, ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE, namespaceId))
		m.compiler.ReportChangeToEnclosingScope(call)
		nullNode := m.astFactory.CreateNull().Srcref(call)
		call.ReplaceWith(nullNode)
		return
	}

	// Remaining calls to goog.module.get() are not alias updates,
	// and should be replaced by a reference to the proper name.
	// Replace "goog.module.get('pkg.Foo')" with either "pkg.Foo" or
	// "module$exports$pkg$Foo".
	exportedNamespace := m.rewriteState.GetExportedNamespaceOrScript(namespaceId)
	if exportedNamespace != nil {
		m.compiler.ReportChangeToEnclosingScope(call)
		exportedNamespaceName :=
			m.astFactory.
				CreateQNameUsingJSTypeInfo(m.globalTypedScope, *exportedNamespace).
				SrcrefTree(call)
		exportedNamespaceName.SetJSType(m.rewriteState.GetGoogModuleNamespaceType(namespaceId))
		exportedNamespaceName.SetOriginalName(&namespaceId)
		call.ReplaceWith(exportedNamespaceName)
	}
}

// Java: ClosureRewriteModule.recordExportsPropertyAssignment(NodeTraversal, Node)
func (m *ClosureRewriteModule) RecordExportsPropertyAssignment(t *NodeTraversal, getpropNode *rhino.Node) {
	if !m.currentScript.isModule {
		return
	}

	parent := getpropNode.GetParent()
	jutil.CheckStateMsg1(parent.IsAssign() || parent.IsExprResult(), "%s", parent)

	exportsNameNode := getpropNode.GetFirstChild()
	jutil.CheckStateMsg1(exportsNameNode.GetString() == "exports", "%s", exportsNameNode)

	if t.InModuleScope() {
		exportName := getpropNode.GetString()
		m.currentScript.namedExports.Add(exportName)
		exportRhs := getpropNode.GetNext()
		namedExport := ClosureRewriteModule_ExportDefinition_NewNamedExport(t, &exportName, exportRhs)
		if !m.currentScript.declareLegacyNamespace &&
			m.currentScript.defaultExport == nil &&
			namedExport.HasInlinableName(m.currentScript.exportsToInline) {
			m.RecordExportToInline(namedExport)
			parent.GetParent().Detach()
		}
	}
}

// Java: ClosureRewriteModule.updateExportsPropertyAssignment(Node, NodeTraversal)
func (m *ClosureRewriteModule) UpdateExportsPropertyAssignment(getpropNode *rhino.Node, t *NodeTraversal) {
	if !m.currentScript.isModule {
		return
	}

	parent := getpropNode.GetParent()
	jutil.CheckStateMsg1(parent.IsAssign() || parent.IsExprResult(), "%s", parent)

	// Update "exports.foo = Foo" to "module$exports$pkg$Foo.foo = Foo";
	exportsNameNode := getpropNode.GetFirstChild()
	jutil.CheckState(exportsNameNode.GetString() == "exports")
	exportedNamespace := m.currentScript.GetExportedNamespace()
	m.SafeSetMaybeQualifiedString(exportsNameNode, *exportedNamespace /* isModuleNamespace= */, false)

	jsdocNode := getpropNode
	if parent.IsAssign() {
		jsdocNode = parent
	}
	closureRewriteModule_markConstAndCopyJsDoc(jsdocNode, jsdocNode)

	// When seeing the first "exports.foo = ..." line put a "var
	// module$exports$pkg$Foo = {};" before it.
	if !m.currentScript.hasCreatedExportObject {
		m.ExportTheEmptyBinaryNamespaceAt(NodeUtil_GetEnclosingStatement(parent), ClosureRewriteModule_AddAt_BEFORE, t)
	}
}

// MaybeUpdateTopLevelName rewrites top level var names from "var foo;
// console.log(foo);" to "var module$contents$Foo_foo;
// console.log(module$contents$Foo_foo);"
//
// Java: ClosureRewriteModule.maybeUpdateTopLevelName(NodeTraversal, Node)
func (m *ClosureRewriteModule) MaybeUpdateTopLevelName(t *NodeTraversal, nameNode *rhino.Node) {
	name := nameNode.GetString()
	if !m.currentScript.isModule || !m.currentScript.topLevelNames.Contains(name) {
		return
	}
	var_ := t.LookupVar(name)
	// If the name refers to a var that is not from the top level scope.
	if var_ == nil || var_.GetScope().(*Scope).GetRootNode() != m.currentScript.rootNode {
		// Then it shouldn't be renamed.
		return
	}

	// If the name is part of a destructuring import, the import rewriting will
	// take care of it
	if var_.GetNameNode() == nameNode &&
		nameNode.GetParent().IsStringKey() &&
		nameNode.GetGrandparent().IsObjectPattern() {
		destructuringLhsNode := nameNode.GetGrandparent().GetParent()
		if NodeUtil_IsCallToNode(destructuringLhsNode.GetLastChild(), closureRewriteModule_GOOG_REQUIRE) ||
			NodeUtil_IsCallToNode(destructuringLhsNode.GetLastChild(), closureRewriteModule_GOOG_REQUIRETYPE) {
			return
		}
	}

	// If the name is an alias for an imported namespace or an exported local,
	// rewrite from "new Foo;" to "new module$exports$Foo;" or "new Foo" to "new
	// module$contents$bar$Foo".
	nameIsAnAlias := m.currentScript.namesToInlineByAlias.ContainsKey(name)
	if nameIsAnAlias && var_.GetNode() != nameNode {
		m.MaybeAddAliasToSymbolTable(nameNode, m.currentScript.namespaceId)

		inline := m.currentScript.namesToInlineByAlias.Get(name)
		namespaceToInline := inline.newName
		if binaryNamespace := m.currentScript.GetBinaryNamespace(); binaryNamespace != nil && namespaceToInline == *binaryNamespace {
			m.currentScript.hasCreatedExportObject = true
		}
		isModuleNamespace :=
			inline.namespaceId != nil &&
				m.rewriteState.scriptDescriptionsByGoogModuleNamespace.ContainsKey(*inline.namespaceId) &&
				!m.rewriteState.scriptDescriptionsByGoogModuleNamespace.Get(*inline.namespaceId).
					willCreateExportsObject
		m.SafeSetMaybeQualifiedString(nameNode, namespaceToInline, isModuleNamespace)

		// Make sure this action won't shadow a local variable.
		if strings.IndexByte(namespaceToInline, '.') != -1 {
			firstQualifiedName := namespaceToInline[:strings.IndexByte(namespaceToInline, '.')]
			shadowedVar := t.LookupVar(firstQualifiedName)
			if shadowedVar == nil ||
				shadowedVar.IsGlobal() ||
				shadowedVar.GetScope().(*Scope).IsModuleScope() {
				return
			}
			t.Report(
				shadowedVar.GetNode(),
				ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR,
				shadowedVar.GetName(),
				namespaceToInline)
		}
		return
	}

	// For non-import alias names rewrite from
	// "var foo; console.log(foo);" to
	// "var module$contents$Foo_foo; console.log(module$contents$Foo_foo);"
	m.SafeSetString(nameNode, jutil.ToJavaString(m.currentScript.contentsPrefix)+name)
}

// MaybeUpdateExportObjectLiteral: for exports like "exports = {prop:
// value}" update the declarations to enforce @const ness (and typedef
// exports).
//
// TODO(blickly): Remove as much of this functionality as possible, now that
// these style of exports are rewritten in ScriptPreprocess step.
//
// Java: ClosureRewriteModule.maybeUpdateExportObjectLiteral(NodeTraversal, Node)
func (m *ClosureRewriteModule) MaybeUpdateExportObjectLiteral(t *NodeTraversal, n *rhino.Node) {
	if !m.currentScript.isModule {
		return
	}

	parent := n.GetParent()
	rhs := parent.GetLastChild()

	if rhs.IsObjectLit() {
		for c := rhs.GetFirstChild(); c != nil; c = c.GetNext() {
			if c.IsComputedProp() {
				t.Report(c, ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY)
			} else if c.IsStringKey() {
				value := c.GetFirstChild()
				m.MaybeUpdateExportDeclToNode(t, c, value)
			}
		}
	}
}

// Java: ClosureRewriteModule.maybeUpdateExportDeclToNode(NodeTraversal, Node, Node)
func (m *ClosureRewriteModule) MaybeUpdateExportDeclToNode(t *NodeTraversal, target *rhino.Node, value *rhino.Node) {
	if !m.currentScript.isModule {
		return
	}

	// If the RHS is a typedef, clone the declaration.
	// Hack alert: clone the typedef declaration if one exists
	// this is a simple attempt that covers the common case of the
	// exports being in the same scope as the typedef declaration.
	// Otherwise the type name might be invalid.
	if value.IsName() {
		currentScope := t.GetScope()
		v := t.GetScope().GetVar(value.GetString())
		if v != nil {
			varScope := v.GetScope().(*Scope)
			if varScope.GetDepth() == currentScope.GetDepth() {
				info := v.GetJSDocInfo()
				if info != nil && info.HasTypedefType() {
					builder := rhino.JSDocInfo_Builder_CopyFrom(info)
					target.SetJSDocInfo(builder.Build())
					return
				}
			}
		}
	}

	closureRewriteModule_markConstAndCopyJsDoc(target, target)
}

// MaybeUpdateExportDeclaration: in module "foo.Bar", rewrite "exports = Bar"
// to "var module$exports$foo$Bar = Bar".
//
// Java: ClosureRewriteModule.maybeUpdateExportDeclaration(NodeTraversal, Node)
func (m *ClosureRewriteModule) MaybeUpdateExportDeclaration(t *NodeTraversal, n *rhino.Node) {
	if !m.currentScript.isModule || !closureRewriteModule_isGoogModuleExportsRef(t, n) || !closureRewriteModule_isAssignTarget(n) {
		return
	}

	assignNode := n.GetParent()
	rhs := assignNode.GetLastChild()
	if rhs != m.currentScript.defaultExport.rhs {
		// This script has duplicate 'exports = ' assignments. Preserve the rhs
		// as an expression but don't declare it as a global variable.
		assignNode.ReplaceWith(rhs.Detach())
		return
	}

	if !m.currentScript.declareLegacyNamespace && m.currentScript.defaultExportLocalName != nil {
		assignNode.GetParent().Detach()

		binaryNamespaceName := m.astFactory.CreateName(*m.currentScript.GetBinaryNamespace(), AstFactory_Type_(n))
		m.DeclareGlobalVariable(binaryNamespaceName, t)
		return
	}

	// Rewrite "exports = ..." as "var module$exports$foo$Bar = ..."
	var jsdocNode *rhino.Node
	if m.currentScript.declareLegacyNamespace {
		legacyQname :=
			m.astFactory.
				CreateQNameUsingJSTypeInfo(m.globalTypedScope, *m.currentScript.namespaceId).
				SrcrefTree(n)
		legacyQname.SetJSType(n.GetJSType())
		n.ReplaceWith(legacyQname)
		jsdocNode = assignNode
	} else {
		rhs.Detach()
		exprResultNode := assignNode.GetParent()
		binaryNamespaceName := m.astFactory.CreateName(*m.currentScript.GetBinaryNamespace(), AstFactory_Type_(n))
		exportsName := "exports"
		binaryNamespaceName.SetOriginalName(&exportsName)
		m.DeclareGlobalVariable(binaryNamespaceName, t)

		exportsObjectCreationNode := rhino.IR_VarWithValue(binaryNamespaceName, rhs)
		exportsObjectCreationNode.SrcrefTreeIfMissing(exprResultNode)
		exportsObjectCreationNode.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
		exprResultNode.ReplaceWith(exportsObjectCreationNode)
		jsdocNode = exportsObjectCreationNode
		m.currentScript.hasCreatedExportObject = true
	}
	closureRewriteModule_markConstAndCopyJsDoc(assignNode, jsdocNode)
	m.compiler.ReportChangeToEnclosingScope(jsdocNode)

	m.MaybeUpdateExportObjectLiteral(t, rhs)
}

// Java: ClosureRewriteModule.maybeUpdateExportNameRef(NodeTraversal, Node)
func (m *ClosureRewriteModule) MaybeUpdateExportNameRef(t *NodeTraversal, n *rhino.Node) {
	if !m.currentScript.isModule || !closureRewriteModule_isGoogModuleExportsRef(t, n) || n.GetParent() == nil {
		return
	}
	if n.GetParent().IsParamList() {
		return
	}

	if m.currentScript.declareLegacyNamespace {
		legacyQname :=
			m.astFactory.
				CreateQName(m.globalTypedStaticScope(), *m.currentScript.namespaceId).
				SrcrefTree(n)
		legacyQname.SetJSType(n.GetJSType())
		n.ReplaceWith(legacyQname)
		m.compiler.ReportChangeToEnclosingScope(legacyQname)
		return
	}

	m.SafeSetString(n, jutil.ToJavaString(m.currentScript.GetBinaryNamespace()))

	// Either this module is going to create it's own exports object at some
	// point or else if it's going to be defensively created automatically then
	// that should have occurred at the top of the file and been done by now.
	jutil.CheckState(m.currentScript.willCreateExportsObject || m.currentScript.hasCreatedExportObject)
}

// Java: ClosureRewriteModule.updateModuleBody(Node)
func (m *ClosureRewriteModule) UpdateModuleBody(moduleBody *rhino.Node) {
	jutil.CheckArgumentMsg1(
		moduleBody.IsModuleBody() && moduleBody.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE),
		"%s",
		moduleBody)
	moduleBody.SetToken(rhino.Token_BLOCK)
	NodeUtil_TryMergeBlock(moduleBody, true)

	for export := range m.currentScript.exportsToInline.Values() {
		nameNode := export.nameDecl.GetNameNode()
		m.SafeSetMaybeQualifiedString(
			nameNode, jutil.ToJavaString(m.currentScript.GetBinaryNamespace())+export.GetExportPostfix(), false)
	}
	jutil.CheckStateMsg1(m.currentScript.isModule, "%s", m.currentScript)
	jutil.CheckStateMsg1(
		m.currentScript.declareLegacyNamespace || m.currentScript.hasCreatedExportObject,
		"%s",
		m.currentScript)

	m.PopScript()
}

// PushScript records the provided script as the current script at top of the
// script stack and add it as a child of the previous current script if there
// was one.
//
// Keeping track of the current script facilitates aggregation of accurate
// script state so that rewriting can run properly. Handles scripts and
// nested goog.modules.
//
// Java: ClosureRewriteModule.pushScript(ScriptDescription)
func (m *ClosureRewriteModule) PushScript(newCurrentScript *ClosureRewriteModule_ScriptDescription) {
	m.currentScript = newCurrentScript
	if !m.scriptStack.IsEmpty() {
		parentScript := m.scriptStack.Get(0)
		parentScript.AddChildScript(m.currentScript)
	}
	m.scriptStack.AddAt(0, m.currentScript)
}

// Java: ClosureRewriteModule.popScript()
func (m *ClosureRewriteModule) PopScript() {
	if m.scriptStack.IsEmpty() {
		panic(jutil.NewNoSuchElementException())
	}
	m.scriptStack.RemoveAt(0)
	if m.scriptStack.IsEmpty() {
		m.currentScript = nil
	} else {
		m.currentScript = m.scriptStack.Get(0)
	}
}

// ExportTheEmptyBinaryNamespaceAt adds the missing "var
// module$exports$pkg$Foo = {};" line.
//
// Java: ClosureRewriteModule.exportTheEmptyBinaryNamespaceAt(Node, AddAt, NodeTraversal)
func (m *ClosureRewriteModule) ExportTheEmptyBinaryNamespaceAt(atNode *rhino.Node, addAt ClosureRewriteModule_AddAt, t *NodeTraversal) {
	if m.currentScript.declareLegacyNamespace {
		return
	}

	binaryNamespaceString := m.currentScript.GetBinaryNamespace()
	moduleType := AstFactory_Type_(m.currentScript.rootNode)
	binaryNamespaceName := m.astFactory.CreateName(*binaryNamespaceString, moduleType)
	binaryNamespaceName.SetOriginalName(m.currentScript.namespaceId)
	m.DeclareGlobalVariable(binaryNamespaceName, t)

	binaryNamespaceExportNode := rhino.IR_VarWithValue(binaryNamespaceName, m.astFactory.CreateObjectLit())
	if addAt == ClosureRewriteModule_AddAt_BEFORE {
		binaryNamespaceExportNode.InsertBefore(atNode)
	} else if addAt == ClosureRewriteModule_AddAt_AFTER {
		binaryNamespaceExportNode.InsertAfter(atNode)
	}
	binaryNamespaceExportNode.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
	binaryNamespaceExportNode.SrcrefTree(atNode)
	closureRewriteModule_markConst(binaryNamespaceExportNode)
	m.compiler.ReportChangeToEnclosingScope(binaryNamespaceExportNode)
	m.currentScript.hasCreatedExportObject = true
}

// Java: ClosureRewriteModule.checkAndSetStrictModeDirective(NodeTraversal, Node)
func ClosureRewriteModule_CheckAndSetStrictModeDirective(t *NodeTraversal, n *rhino.Node) {
	jutil.CheckStateMsg1(n.IsScript(), "%s", n)

	if n.IsUseStrict() {
		t.Report(n, ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE)
	} else {
		n.SetUseStrict(true)
	}
}

// Java: ClosureRewriteModule.markConst(Node)
func closureRewriteModule_markConst(n *rhino.Node) {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(n.GetJSDocInfo())
	builder.RecordConstancy()
	n.SetJSDocInfo(builder.Build())
}

// Java: ClosureRewriteModule.maybeSplitMultiVar(Node)
func closureRewriteModule_maybeSplitMultiVar(rhsNode *rhino.Node) {
	statementNode := rhsNode.GetGrandparent()
	if !statementNode.IsVar() || !statementNode.HasMoreThanOneChild() {
		return
	}

	nameNode := rhsNode.GetParent()
	nameNode.Detach()
	rhsNode.Detach()
	rhino.IR_VarWithValue(nameNode, rhsNode).InsertBefore(statementNode)
}

// Java: ClosureRewriteModule.markConstAndCopyJsDoc(Node, Node)
func closureRewriteModule_markConstAndCopyJsDoc(from *rhino.Node, target *rhino.Node) {
	info := from.GetJSDocInfo()
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
	builder.RecordConstancy()
	target.SetJSDocInfo(builder.Build())
}

// Java: ClosureRewriteModule.recordExportToInline(ExportDefinition)
func (m *ClosureRewriteModule) RecordExportToInline(exportDefinition *ClosureRewriteModule_ExportDefinition) {
	if !exportDefinition.HasInlinableName(m.currentScript.exportsToInline) {
		jutil.CheckStateMsg2(
			false,
			"exportDefinition: %s\n\nexportsToInline keys: %s",
			exportDefinition,
			jutil.CollectionToString(m.currentScript.exportsToInline.KeySlice()))
	}
	previous := m.currentScript.exportsToInline.Put(exportDefinition.nameDecl, exportDefinition)
	jutil.CheckStateMsg1(
		nil == previous,
		"Already found a mapping for inlining export: %s",
		exportDefinition.nameDecl)
	localName := exportDefinition.GetLocalName()
	fullExportedName :=
		jutil.ToJavaString(m.currentScript.GetBinaryNamespace()) + exportDefinition.GetExportPostfix()
	m.RecordNameToInline(*localName, fullExportedName /* namespaceId= */, nil)
}

// Java: ClosureRewriteModule.recordNameToInline(String, String, String)
func (m *ClosureRewriteModule) RecordNameToInline(aliasName string, newName string, namespaceId *string) {
	// This intentionally overwrites a possibly pre-existing alias of the same
	// name. User code might import the same name twice, with the same variable
	// name. That's an error (duplicate variable definition, reported in
	// TypeValidator), but this code still shouldn't crash on it.
	m.currentScript.namesToInlineByAlias.Put(aliasName, NewClosureRewriteModule_AliasName(newName, namespaceId))
}

// ReportUnrecognizedRequires examines queue'ed unrecognizedRequires to
// categorize and report them as either missing module, missing namespace or
// late provide.
//
// Java: ClosureRewriteModule.reportUnrecognizedRequires()
func (m *ClosureRewriteModule) ReportUnrecognizedRequires() {
	for unrecognizedRequire := range m.unrecognizedRequires.All() {
		namespaceId := unrecognizedRequire.namespaceId

		requireNode := unrecognizedRequire.requireNode
		targetGoogModuleExists := m.rewriteState.ContainsModule(namespaceId)
		targetLegacyScriptExists := m.rewriteState.providedNamespaces.Contains(namespaceId)

		if targetGoogModuleExists || targetLegacyScriptExists {
			// The required thing actually was available somewhere in the program
			// but just wasn't available as early as the require statement would
			// have liked.
			continue
		}

		// Remove the require node so this problem isn't reported again in
		// ProcessClosurePrimitives.
		if m.preserveSugar {
			continue
		}

		if NodeUtil_GetEnclosingScript(requireNode) == nil {
			continue // It's already been removed; nothing to do.
		}

		m.compiler.ReportChangeToEnclosingScope(requireNode)
		enclosingStatement := NodeUtil_GetEnclosingStatement(requireNode)

		// To make compilation with partial source information work for Clutz,
		// delete any name declarations in the enclosing statement completely.
		// For non-declarations, simply replace the invalid require with null.
		if !NodeUtil_IsNameDeclaration(enclosingStatement) {
			requireNode.ReplaceWith(m.astFactory.CreateNull().Srcref(requireNode))
			continue
		}

		enclosingStatement.Detach()
		NodeUtil_VisitLhsNodesInNode(
			enclosingStatement, func(lhs *rhino.Node) { m.syntheticExterns.PutIfAbsent(lhs.GetString(), lhs) })
	}
}

// Java: ClosureRewriteModule.safeSetString(Node, String)
func (m *ClosureRewriteModule) SafeSetString(n *rhino.Node, newString string) {
	if n.GetString() == newString {
		return
	}

	originalName := n.GetString()
	n.SetString(newString)
	if n.GetOriginalName() == nil {
		n.SetOriginalName(&originalName)
	}
	// TODO(blickly): It would be better not to be renaming detached nodes
	changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(n)
	if changeScope != nil {
		m.compiler.ReportChangeToChangeScope(changeScope)
	}
}

// SafeSetMaybeQualifiedString replaces an identifier with a potentially
// qualified name.
//
// Java: ClosureRewriteModule.safeSetMaybeQualifiedString(Node, String, boolean)
func (m *ClosureRewriteModule) SafeSetMaybeQualifiedString(
	nameNode *rhino.Node, newString string, isModuleNamespace bool) {
	if !strings.Contains(newString, ".") {
		m.SafeSetString(nameNode, newString)
		parent := nameNode.GetParent()
		if isModuleNamespace &&
			parent.IsGetProp() &&
			nameNode.GetGrandparent().IsCall() &&
			parent.IsFirstChildOf(nameNode.GetGrandparent()) {
			// In cases where we're calling a function off a module namespace,
			// don't pass the module namespace as `this`.
			nameNode.GetGrandparent().PutBooleanProp(rhino.Node_FREE_CALL, true)
		}
		return
	}

	// When replacing with a dotted fully qualified name it's already better
	// than an original name.
	nameParent := nameNode.GetParent()
	newQualifiedName :=
		m.astFactory.
			CreateQNameUsingJSTypeInfo(m.globalTypedScope, newString).
			SrcrefTree(nameNode)
	// Sometimes the typechecker gave `nameNode` the correct type, but we can't
	// infer the right type for `newQualifiedName`. If so, giving
	// `newQualifiedName` the same type typechecking used for `nameNode` is
	// less confusing.
	newQualifiedName.SetJSType(nameNode.GetJSType())

	replaced := closureRewriteModule_safeSetStringIfDeclaration(nameParent, nameNode, newQualifiedName)
	if replaced {
		return
	}

	nameNode.ReplaceWith(newQualifiedName)
	// Given import "var Bar = goog.require('foo.Bar');" here we replace a
	// usage of Bar with foo.Bar if Bar is goog.provided. 'foo' node is
	// generated and never visible to user. Because of that we should mark all
	// such nodes as non-indexable leaving only Bar indexable. Given that
	// replacement is GETPROP node, prefix is first child. It's also possible
	// that replacement is single-part namespace. Like goog.provide('Bar') in
	// that case replacement won't have children.
	if newQualifiedName.HasChildren() {
		newQualifiedName.GetFirstChild().MakeNonIndexableRecursive()
	}
	m.compiler.ReportChangeToEnclosingScope(newQualifiedName)
}

// closureRewriteModule_safeSetStringIfDeclaration sets the string if given
// a declaration, return whether or not the name was changed.
//
// Java: ClosureRewriteModule.safeSetStringIfDeclaration(Node, Node, Node)
func closureRewriteModule_safeSetStringIfDeclaration(
	nameParent *rhino.Node, nameNode *rhino.Node, newQualifiedName *rhino.Node) bool {

	jsdoc := nameParent.GetJSDocInfo()

	switch nameParent.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if !NodeUtil_IsStatement(nameParent) || nameParent.GetFirstChild() != nameNode {
			return false
		}

		placeholder := rhino.IR_Empty()
		nameParent.ReplaceWith(placeholder)
		newDeclaration :=
			NodeUtil_GetDeclarationFromName(newQualifiedName, nameParent, rhino.Token_VAR, jsdoc)
		if NodeUtil_IsExprAssign(newDeclaration) {
			assign := newDeclaration.GetOnlyChild()
			assign.SetJSType(nameNode.GetJSType())
			closureRewriteModule_updateSourceInfoForExportedTopLevelVariable(assign, nameNode)
		}
		nameParent.SetJSDocInfo(nil)
		newDeclaration.SrcrefTreeIfMissing(nameParent)
		placeholder.ReplaceWith(newDeclaration)
		NodeUtil_RemoveName(nameParent)
		return true
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		var rhs *rhino.Node
		if nameNode.HasChildren() {
			rhs = nameNode.GetLastChild().Detach()
		}
		if jsdoc == nil {
			// Get inline JSDocInfo if there is no JSDoc on the actual
			// declaration.
			jsdoc = nameNode.GetJSDocInfo()
		}
		newStatement :=
			NodeUtil_GetDeclarationFromName(newQualifiedName, rhs, rhino.Token_VAR, jsdoc)
		if NodeUtil_IsExprAssign(newStatement) {
			assign := newStatement.GetOnlyChild()
			assign.SetJSType(nameNode.GetJSType())
			closureRewriteModule_updateSourceInfoForExportedTopLevelVariable(assign, nameNode)
			if nameParent.IsConst() {
				// When replacing `const name = ...;` with `some.prop = ...`,
				// ensure that `some.prop` is annotated @const.
				jsdocBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(jsdoc)
				jsdocBuilder.RecordConstancy()
				jsdoc = jsdocBuilder.Build()
				assign.SetJSDocInfo(jsdoc)
			}
		}
		newStatement.SrcrefTreeIfMissing(nameParent)
		NodeUtil_ReplaceDeclarationChild(nameNode, newStatement)
		return true
	case rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		panic(jutil.NewRuntimeException("Not supported"))
	default:
	}
	return false
}

// closureRewriteModule_updateSourceInfoForExportedTopLevelVariable: if we had
// something like `const FOO = "text"` and we export `FOO`, change the source
// location information for the rewritten FOO. The replacement should be
// something like MOD.FOO = "text", so we look for MOD.FOO and replace the
// source location for FOO to the original location of FOO.
//
// Java: ClosureRewriteModule.updateSourceInfoForExportedTopLevelVariable(Node, Node)
func closureRewriteModule_updateSourceInfoForExportedTopLevelVariable(assign *rhino.Node, sourceName *rhino.Node) {
	jutil.CheckState(assign.IsAssign())
	jutil.CheckState(sourceName.IsName())

	// ASSIGN always has two children.
	getProp := assign.GetFirstChild()
	if !getProp.IsGetProp() {
		return
	}

	var name string
	if originalName := sourceName.GetOriginalName(); originalName != nil {
		name = *originalName
	} else {
		name = sourceName.GetString()
	}

	// The source range of this NAME includes its declared value, which we
	// don't want.
	sourceName = sourceName.CloneNode()
	sourceName.SetLength(jutil.StringLength(name))

	// Receiver and prop string should both use the position of the source
	// name.
	getProp.SrcrefTree(sourceName)
}

// Java: ClosureRewriteModule.isTopLevel(NodeTraversal, Node, ScopeType)
func (m *ClosureRewriteModule) IsTopLevel(t *NodeTraversal, n *rhino.Node, scopeType ClosureRewriteModule_ScopeType) bool {
	if scopeType == ClosureRewriteModule_ScopeType_EXEC_CONTEXT {
		return t.InGlobalScope() || t.GetClosestHoistScopeRoot() == m.currentScript.rootNode
	} else {
		// Must be ScopeType.BLOCK;
		return n.GetParent() == m.currentScript.rootNode
	}
}

// Java: ClosureRewriteModule.toModuleContentsPrefix(String)
func closureRewriteModule_toModuleContentsPrefix(namespaceId string) string {
	return closureRewriteModule_MODULE_CONTENTS_PREFIX + strings.ReplaceAll(namespaceId, ".", "$") + "_"
}

// Java: ClosureRewriteModule.isModuleExport(String)
func ClosureRewriteModule_IsModuleExport(name string) bool {
	return strings.HasPrefix(name, ClosureRewriteModule_MODULE_EXPORTS_PREFIX)
}

// Java: ClosureRewriteModule.isModuleContent(String)
func ClosureRewriteModule_IsModuleContent(name string) bool {
	return strings.HasPrefix(name, closureRewriteModule_MODULE_CONTENTS_PREFIX)
}

// closureRewriteModule_isGoogModuleExportsRef returns whether this is a) a
// reference to the name "exports" and b) based on scoping, actually refers
// to the implicit goog.module exports object.
//
// Java: ClosureRewriteModule.isGoogModuleExportsRef(NodeTraversal, Node)
func closureRewriteModule_isGoogModuleExportsRef(t *NodeTraversal, target *rhino.Node) bool {
	if !target.IsName() || !target.MatchesName("exports") {
		return false
	}
	exportsVar := t.LookupVar("exports")
	// note: exportsVar may be null for third-party code using UMD patterns
	// like
	//   if ("object"==typeof exports) { [...]
	// and if null, is not in a goog.module => not a goog.module exports var
	return exportsVar != nil && exportsVar.IsGoogModuleExports()
}

// closureRewriteModule_isExportPropertyAssignment returns whether the
// getprop is used as an assignment target, and that target represents a
// module export. Note: that "export.name = value" is an export, while
// "export.name.foo = value" is not (it is an assignment to a property of an
// exported value).
//
// Java: ClosureRewriteModule.isExportPropertyAssignment(NodeTraversal, Node)
func closureRewriteModule_isExportPropertyAssignment(t *NodeTraversal, n *rhino.Node) bool {
	target := n.GetFirstChild()
	return (closureRewriteModule_isAssignTarget(n) || closureRewriteModule_isTypedefTarget(n)) &&
		closureRewriteModule_isGoogModuleExportsRef(t, target)
}

// Java: ClosureRewriteModule.isAssignTarget(Node)
func closureRewriteModule_isAssignTarget(n *rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsAssign() && parent.GetFirstChild() == n
}

// Java: ClosureRewriteModule.isTypedefTarget(Node)
func closureRewriteModule_isTypedefTarget(n *rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsExprResult() && parent.GetFirstChild() == n
}

// MaybeAddAliasToSymbolTable adds alias nodes to the symbol table as they
// going to be removed by rewriter. Example aliases:
//
// const Foo = goog.require('my.project.Foo'); const bar =
// goog.require('my.project.baz'); const {baz} =
// goog.require('my.project.utils');
//
// Java: ClosureRewriteModule.maybeAddAliasToSymbolTable(Node, String)
func (m *ClosureRewriteModule) MaybeAddAliasToSymbolTable(n *rhino.Node, module *string) {
	if m.preprocessorSymbolTable != nil {
		n.PutBooleanProp(rhino.Node_MODULE_ALIAS, true)
		// Alias can be used in js types. Types have node type STRING and not
		// NAME so we have to use their name as string.
		var nodeName *string
		if n.IsStringLit() {
			s := n.GetString()
			nodeName = &s
		} else {
			nodeName = m.preprocessorSymbolTable.GetQualifiedName(n)
		}
		// We need to include module as part of the name because aliases are
		// local to current module. Aliases with the same name from different
		// module should be completely different entities.
		name := "alias_" + jutil.ToJavaString(module) + "_" + jutil.ToJavaString(nodeName)
		m.preprocessorSymbolTable.AddReferenceWithName(n, name)
	}
}

// Java: ClosureRewriteModule.declareGlobalVariable(Node, NodeTraversal)
func (m *ClosureRewriteModule) DeclareGlobalVariable(n *rhino.Node, t *NodeTraversal) {
	jutil.CheckState(n.IsName())
	if m.globalTypedScope == nil {
		return
	}

	name := n.GetString()
	if m.globalTypedScope.HasOwnSlot(name) {
		original := m.globalTypedScope.GetOwnSlot(name).(*TypedVar).GetNode()
		originalString := "<unknown>"
		if original != nil {
			originalString = original.String()
		}
		t.Report(
			t.GetCurrentScript(),
			ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT,
			name,
			originalString)
	} else {
		type_ := jutil.CheckNotNull(n.GetJSType())
		m.globalTypedScope.Declare(name, n, type_, t.GetInput(), false)
	}
}
