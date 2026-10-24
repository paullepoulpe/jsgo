// Port of com/google/javascript/jscomp/CheckClosureImports.java @ closure-compiler v20261021

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

// CheckClosureImports checks all goog.requires, goog.module.gets,
// goog.forwardDeclares, and goog.requireTypes in all files.
//
// Checks that these dependency calls both contain a valid Closure namespace
// and are in an acceptable location (e.g. a goog.require cannot be within a
// function).
type CheckClosureImports struct {
	compiler       AbstractCompiler
	checker        *CheckClosureImports_Checker
	namespacesSeen *jutil.LinkedHashSet[string]
	chunkGraph     *JSChunkGraph
}

var _ CompilerPass = (*CheckClosureImports)(nil)

// CheckClosureImports_ClosureImport is CheckClosureImports.ClosureImport.
type CheckClosureImports_ClosureImport int32

const (
	CheckClosureImports_ClosureImport_REQUIRE CheckClosureImports_ClosureImport = iota
	CheckClosureImports_ClosureImport_FORWARD_DECLARE
	CheckClosureImports_ClosureImport_REQUIRE_TYPE
	CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC
)

// Java: CheckClosureImports.ClosureImport.values()
func CheckClosureImports_ClosureImport_values() []CheckClosureImports_ClosureImport {
	return []CheckClosureImports_ClosureImport{
		CheckClosureImports_ClosureImport_REQUIRE,
		CheckClosureImports_ClosureImport_FORWARD_DECLARE,
		CheckClosureImports_ClosureImport_REQUIRE_TYPE,
		CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC,
	}
}

// Java: CheckClosureImports.ClosureImport.name()
func (c CheckClosureImports_ClosureImport) String() string {
	return [...]string{"REQUIRE", "FORWARD_DECLARE", "REQUIRE_TYPE", "REQUIRE_DYNAMIC"}[c]
}

// AllowDestructuring is whether or not a destructuring alias is allowed,
// e.g. `const {Class} = goog.require('namespace');`.
//
// Java: CheckClosureImports.ClosureImport.allowDestructuring()
func (c CheckClosureImports_ClosureImport) AllowDestructuring() bool {
	switch c {
	case CheckClosureImports_ClosureImport_REQUIRE:
		return true
	case CheckClosureImports_ClosureImport_FORWARD_DECLARE:
		return false
	case CheckClosureImports_ClosureImport_REQUIRE_TYPE:
		return true
	default: // REQUIRE_DYNAMIC
		return true
	}
}

// AliasMustBeConstant is true if the alias of the call must be constant,
// false if it can be `let`.
//
// Note: since not every goog.module is using ES6, `const` is only checked in
// ES6 modules. `let` will always be disallowed if this is false. But `var`
// will be allowed in goog.module files even if this returns true (since it
// assumes the file is non-ES6).
//
// Java: CheckClosureImports.ClosureImport.aliasMustBeConstant()
func (c CheckClosureImports_ClosureImport) AliasMustBeConstant() bool {
	switch c {
	case CheckClosureImports_ClosureImport_REQUIRE:
		return true
	case CheckClosureImports_ClosureImport_FORWARD_DECLARE:
		return false
	case CheckClosureImports_ClosureImport_REQUIRE_TYPE:
		return true
	default: // REQUIRE_DYNAMIC
		return true
	}
}

// MustBeOrdered is true if the providing symbol must appear before this call
// in the program.
//
// A `goog.require` for something can only appear after it is defined (e.g.
// `goog.provide`, `goog.module`, or `goog.declareModuleId` call). But the
// same is not true for `goog.requireType` and `goog.forwardDeclare`, which
// can appear before the definition.
//
// Java: CheckClosureImports.ClosureImport.mustBeOrdered()
func (c CheckClosureImports_ClosureImport) MustBeOrdered() bool {
	switch c {
	case CheckClosureImports_ClosureImport_REQUIRE:
		return true
	case CheckClosureImports_ClosureImport_FORWARD_DECLARE:
		return false
	case CheckClosureImports_ClosureImport_REQUIRE_TYPE:
		return false
	default: // REQUIRE_DYNAMIC
		return false
	}
}

// CallName is the human readable name of this call, used in error
// reporting.
//
// Java: CheckClosureImports.ClosureImport.callName()
func (c CheckClosureImports_ClosureImport) CallName() string {
	switch c {
	case CheckClosureImports_ClosureImport_REQUIRE:
		return "goog.require"
	case CheckClosureImports_ClosureImport_FORWARD_DECLARE:
		return "goog.forwardDeclare"
	case CheckClosureImports_ClosureImport_REQUIRE_TYPE:
		return "goog.requireType"
	default: // REQUIRE_DYNAMIC
		return "goog.requireDynamic"
	}
}

// Java: CheckClosureImports.INVALID_CLOSURE_IMPORT_DESTRUCTURING
var CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING = DiagnosticType_Error(
	"JSC_INVALID_CLOSURE_IMPORT_DESTRUCTURING",
	"Destructuring {0} must be a simple object pattern.")

// Java: CheckClosureImports.ONE_CLOSURE_IMPORT_PER_DECLARATION
var CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION = DiagnosticType_Error(
	"JSC_ONE_CLOSURE_IMPORT_PER_DECLARATION",
	"There may only be one {0} per var/let/const declaration.")

// Java: CheckClosureImports.INVALID_CLOSURE_IMPORT_CALL
var CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL = DiagnosticType_Error(
	"JSC_INVALID_CLOSURE_IMPORT_CALL", "{0} parameter must be a string literal.")

// Java: CheckClosureImports.LATE_PROVIDE_ERROR
var CheckClosureImports_LATE_PROVIDE_ERROR = DiagnosticType_Error(
	"JSC_LATE_PROVIDE_ERROR", "Required namespace \"{0}\" not provided yet.")

// Java: CheckClosureImports.LET_CLOSURE_IMPORT
var CheckClosureImports_LET_CLOSURE_IMPORT = DiagnosticType_Disabled(
	"JSC_LET_CLOSURE_IMPORT",
	"Module imports must be constant. Please use ''const'' instead of ''let''.")

// Java: CheckClosureImports.NO_CLOSURE_IMPORT_DESTRUCTURING
var CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING = DiagnosticType_Error(
	"JSC_NO_CLOSURE_IMPORT_DESTRUCTURING", "Cannot destructure the return value of {0}")

// Java: CheckClosureImports.LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE
var CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE = DiagnosticType_Error(
	"JSC_LHS_OF_CLOUSRE_IMPORT_MUST_BE_CONST_IN_ES_MODULE",
	"The left side of a {0} must use ''const'' (not ''let'' or ''var'') in an ES module.")

// Java: CheckClosureImports.CROSS_CHUNK_REQUIRE_ERROR
var CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR = DiagnosticType_Warning(
	"JSC_XMODULE_REQUIRE_ERROR",
	"namespace \"{0}\" is required in chunk {2} but provided in chunk {1}."+
		" Is chunk {2} missing a dependency on chunk {1}?")

var (
	checkClosureImports_GOOG_REQUIRE         = rhino.IR_Getprop(rhino.IR_Name("goog"), "require")
	checkClosureImports_GOOG_MODULE_GET      = rhino.IR_GetpropWithMoreProps(rhino.IR_Name("goog"), "module", "get")
	checkClosureImports_GOOG_FORWARD_DECLARE = rhino.IR_Getprop(rhino.IR_Name("goog"), "forwardDeclare")
	checkClosureImports_GOOG_REQUIRE_TYPE    = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireType")
	checkClosureImports_GOOG_REQUIRE_DYNAMIC = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireDynamic")
)

// Java: CheckClosureImports(AbstractCompiler, ModuleMetadataMap)
func NewCheckClosureImports(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *CheckClosureImports {
	c := &CheckClosureImports{}
	c.compiler = compiler
	c.checker = c.NewChecker(compiler, moduleMetadataMap)
	c.namespacesSeen = jutil.NewLinkedHashSet[string]()
	c.chunkGraph = compiler.GetChunkGraph()
	return c
}

// Java: CheckClosureImports.process(Node, Node)
func (c *CheckClosureImports) Process(externs rhino.Node, root rhino.Node) {
	sc := NewSyntacticScopeCreator(c.compiler)
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c.checker,
		scopeCreator: sc,
	}
	b.Traverse(externs)
	b.Traverse(root)
}

// CheckClosureImports_Checker is CheckClosureImports.Checker (an inner
// class: outer is the enclosing CheckClosureImports).
type CheckClosureImports_Checker struct {
	NodeTraversal_AbstractModuleCallbackBase
	outer             *CheckClosureImports
	moduleMetadataMap *ModuleMetadataMap
}

func (*CheckClosureImports_Checker) noVisitLeaves() {}

// Java: CheckClosureImports.Checker(AbstractCompiler, ModuleMetadataMap)
func (c *CheckClosureImports) NewChecker(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *CheckClosureImports_Checker {
	k := &CheckClosureImports_Checker{outer: c}
	k.InitNodeTraversal_AbstractModuleCallbackBase(k, compiler, moduleMetadataMap)
	k.moduleMetadataMap = moduleMetadataMap
	return k
}

// Java: CheckClosureImports.Checker.enterModule(ModuleMetadata, Node)
func (k *CheckClosureImports_Checker) EnterModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node) {
	for ns := range currentModule.GoogNamespaces().ElementSet().All() {
		k.outer.namespacesSeen.Add(ns)
	}
}

// Java: CheckClosureImports.Checker.visit(NodeTraversal, Node, ModuleMetadata, Node)
func (k *CheckClosureImports_Checker) VisitWithModule(
	t *NodeTraversal,
	n rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot rhino.Node) {
	// currentModule is null on ROOT nodes.
	if currentModule == nil {
		return
	}

	parent := n.GetParent()

	if n.IsCall() {
		if kindOfCall, ok := k.kindOfCallValue(n); ok {
			k.CheckImport(t, n, parent, currentModule, kindOfCall)
		} else if n.GetFirstChild().MatchesQualifiedNameNode(checkClosureImports_GOOG_MODULE_GET) {
			k.CheckGoogModuleGet(t, n, currentModule)
		}
	} else if n.IsName() {
		k.CheckValidImportCodeReference(t, n)
	}
}

// KindOfCall, given a call node, determines what kind of ClosureImport it
// is. Returns null if not a ClosureImport.
//
// Java: CheckClosureImports.Checker.kindOfCall(Node)
func (k *CheckClosureImports_Checker) KindOfCall(callNode rhino.Node) *CheckClosureImports_ClosureImport {
	kind, ok := k.kindOfCallValue(callNode)
	if !ok {
		return nil
	}
	return &kind
}

func (k *CheckClosureImports_Checker) kindOfCallValue(callNode rhino.Node) (CheckClosureImports_ClosureImport, bool) {
	jutil.CheckState(callNode.IsCall())
	callee := callNode.GetFirstChild()
	if callee.MatchesQualifiedNameNode(checkClosureImports_GOOG_REQUIRE) {
		return CheckClosureImports_ClosureImport_REQUIRE, true
	} else if callee.MatchesQualifiedNameNode(checkClosureImports_GOOG_FORWARD_DECLARE) {
		return CheckClosureImports_ClosureImport_FORWARD_DECLARE, true
	} else if callee.MatchesQualifiedNameNode(checkClosureImports_GOOG_REQUIRE_TYPE) {
		return CheckClosureImports_ClosureImport_REQUIRE_TYPE, true
	} else if callee.MatchesQualifiedNameNode(checkClosureImports_GOOG_REQUIRE_DYNAMIC) {
		return CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC, true
	}
	return 0, false
}

// Java: CheckClosureImports.Checker.checkValidImportCodeReference(NodeTraversal, Node)
func (k *CheckClosureImports_Checker) CheckValidImportCodeReference(t *NodeTraversal, nameNode rhino.Node) {
}

// Java: CheckClosureImports.Checker.checkGoogModuleGet(NodeTraversal, Node, ModuleMetadata)
func (k *CheckClosureImports_Checker) CheckGoogModuleGet(t *NodeTraversal, call rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata) {
	if t.InModuleHoistScope() {
		t.Report(call, ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
		return
	}

	if !currentModule.IsEs6Module() &&
		!currentModule.IsGoogModule() &&
		k.outer.compiler.GetOptions().GetModuleResolutionMode() != ModuleLoader_ResolutionMode_WEBPACK {
		// Here we are making a heuristic check of the use of goog.module.get.
		// There are many different ways to deliberately subvert these checks.
		// What we are trying to avoid is accidential use of a `goog.provide`
		// file as module scoped. So we want to avoid:
		//   `const Foo = goog.module.get('a.b.Foo');
		// and similar patterns, but we won't try to make this a perfect check.

		// Anything not in global scope is ok.
		if t.InGlobalScope() {
			// Find the root of expressions like `goog.module.get().export`
			getSubExprRoot := call
			for NodeUtil_IsNormalOrOptChainGetProp(getSubExprRoot.GetParent()) {
				getSubExprRoot = getSubExprRoot.GetParent()
			}

			invalid := false
			parent := getSubExprRoot.GetParent()
			if parent.IsName() || parent.IsDestructuringLhs() {
				invalid = true
			} else if parent.IsAssign() {
				target := parent.GetFirstChild()
				if target.IsName() || target.IsObjectPattern() {
					invalid = true
				}
			}

			if invalid {
				t.Report(call, ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE)
			}
		}
	}

	if !call.HasTwoChildren() || !call.GetSecondChild().IsStringLit() {
		t.Report(call, CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL, "goog.module.get")
		return
	}

	namespace := call.GetSecondChild().GetString()
	requiredModule := k.moduleMetadataMap.GetModulesByGoogNamespace().Get(namespace)

	if requiredModule == nil {
		t.Report(call, ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, namespace)
		return
	}

	maybeAssign := call.GetParent()
	isFillingAnAlias := maybeAssign.IsAssign() && maybeAssign.GetFirstChild().IsName()
	isModule := currentModule.IsGoogModule() || currentModule.IsEs6Module()
	if isFillingAnAlias && isModule {
		aliasName := call.GetParent().GetFirstChild().GetString()

		// If the assignment isn't into a var in our scope then it's not ok.
		aliasVar := t.GetScope().GetVar(aliasName)
		if aliasVar == nil {
			t.Report(call, ClosureRewriteModule_INVALID_GET_ALIAS)
			return
		}

		// Even if it was to a var in our scope it should still only rewrite if
		// the var looked like:
		//   let x = goog.forwardDeclare('a.namespace');
		aliasVarNodeRhs := NodeUtil_GetRValueOfLValue(aliasVar.GetNode())

		if aliasVarNodeRhs.IsNil() || !NodeUtil_IsCallToNode(aliasVarNodeRhs, checkClosureImports_GOOG_FORWARD_DECLARE) {
			t.Report(call, ClosureRewriteModule_INVALID_GET_ALIAS)
			return
		}

		if namespace != aliasVarNodeRhs.GetLastChild().GetString() {
			t.Report(call, ClosureRewriteModule_INVALID_GET_ALIAS)
			return
		}
	}
}

// Java: CheckClosureImports.Checker.isValidDestructuringImport(Node)
func (k *CheckClosureImports_Checker) IsValidDestructuringImport(destructuringLhs rhino.Node) bool {
	jutil.CheckArgument(destructuringLhs.IsDestructuringLhs())
	objectPattern := destructuringLhs.GetFirstChild()
	if !objectPattern.IsObjectPattern() {
		return false
	}
	for stringKey := objectPattern.GetFirstChild(); !stringKey.IsNil(); stringKey = stringKey.GetNext() {
		if !stringKey.IsStringKey() {
			return false
		}
		if !stringKey.HasChildren() || !stringKey.GetFirstChild().IsName() {
			return false
		}
	}
	return true
}

// Java: CheckClosureImports.Checker.checkShortName(NodeTraversal, Node, String)
func (k *CheckClosureImports_Checker) CheckShortName(t *NodeTraversal, shortNameNode rhino.Node, namespace string) {
	shortName := shortNameNode.GetString()
	lastSegment := jutil.StringSubstringFrom(namespace, jutil.StringLastIndexOf(namespace, ".")+1)
	if shortName == lastSegment || lastSegment == "" {
		return
	}

	if closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(shortName, 0)) != closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(lastSegment, 0)) {
		var newStartChar uint16
		if closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(shortName, 0)) {
			newStartChar = closureCheckModule_asciiToLowerCase(jutil.StringCharAt(shortName, 0))
		} else {
			newStartChar = closureCheckModule_asciiToUpperCase(jutil.StringCharAt(shortName, 0))
		}
		correctedName := jutil.StringFromChar(newStartChar) + jutil.StringSubstringFrom(shortName, 1)
		t.Report(shortNameNode, ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION, shortName, correctedName)
	}
}

// Java: CheckClosureImports.Checker.checkImport(NodeTraversal, Node, Node, ModuleMetadata, ClosureImport)
func (k *CheckClosureImports_Checker) CheckImport(
	t *NodeTraversal,
	call rhino.Node,
	parent rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	importType CheckClosureImports_ClosureImport) {
	atTopLevelScope := t.InGlobalHoistScope() || t.InModuleScope()
	isModule := currentModule.IsModule()
	validAssignment := isModule || parent.IsExprResult()
	isAliased := NodeUtil_IsNameDeclaration(parent.GetParent())

	// goog.requireDynamic() can be in non-top scope.
	if (!atTopLevelScope || !validAssignment) && importType != CheckClosureImports_ClosureImport_REQUIRE_DYNAMIC {
		t.Report(call, ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
		return
	}

	if !call.HasTwoChildren() || !call.GetSecondChild().IsStringLit() {
		t.Report(call, CheckClosureImports_INVALID_CLOSURE_IMPORT_CALL, importType.CallName())
		return
	}

	if isModule && isAliased {
		declaration := parent.GetParent()

		if !declaration.HasOneChild() {
			t.Report(call, CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION, importType.CallName())
			return
		}

		if importType.AliasMustBeConstant() {
			if declaration.IsLet() {
				t.Report(call, CheckClosureImports_LET_CLOSURE_IMPORT, importType.CallName())
			} else if !declaration.IsConst() && currentModule.IsEs6Module() {
				t.Report(call, CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE, importType.CallName())
			}
		}

		lhs := declaration.GetFirstChild()
		if lhs.IsDestructuringLhs() {
			if importType.AllowDestructuring() {
				if !k.IsValidDestructuringImport(lhs) {
					t.Report(declaration, CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING, importType.CallName())
				}
			} else {
				t.Report(declaration, CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING, importType.CallName())
			}
		} else if lhs.IsName() {
			k.CheckShortName(t, lhs, call.GetLastChild().GetString())
		}
	}

	namespace := call.GetSecondChild().GetString()
	requiredModule := k.moduleMetadataMap.GetModulesByGoogNamespace().Get(namespace)

	if requiredModule == nil {
		if importType == CheckClosureImports_ClosureImport_FORWARD_DECLARE {
			// Ok to forwardDeclare any global, sadly. This is some bad legacy
			// behavior and people ought to use externs.
			if isAliased && isModule {
				// Special case `const Foo = goog.forwardDeclare('a.Foo');`. For
				// legacy reasons, we allow js_library to be less strict about
				// missing forwardDeclares vs. other missing Closure imports.
				t.Report(
					call,
					ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE,
					namespace)
			}
			return
		}
		t.Report(call, ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, namespace)
		return
	}

	if importType.MustBeOrdered() {
		k.outer.VerifyRequireOrder(namespace, call, t, requiredModule)
	}

	if importType == CheckClosureImports_ClosureImport_REQUIRE &&
		currentModule.IsEs6Module() &&
		requiredModule.IsEs6Module() {
		t.Report(call, Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE)
	}
}

// Java: CheckClosureImports.verifyRequireOrder(String, Node, NodeTraversal, ModuleMetadata)
func (c *CheckClosureImports) VerifyRequireOrder(
	namespace string, call rhino.Node, t *NodeTraversal, requiredModule *ModuleMetadataMap_ModuleMetadata) {
	if !c.namespacesSeen.Contains(namespace) {
		t.Report(call, CheckClosureImports_LATE_PROVIDE_ERROR, namespace)
		return
	}

	if requiredModule.RootNode().IsNil() || requiredModule.RootNode().IsFromExterns() {
		return // synthetic metadata for tests may have no root node.
	}

	requiredInputId := NodeUtil_GetInputId(requiredModule.RootNode())
	requiredInput :=
		jutil.CheckNotNullMsg1(
			c.compiler.GetInput(requiredInputId), "Cannot find CompilerInput for %s", requiredModule)

	requiredChunk := requiredInput.GetChunk()
	currentChunk := t.GetChunk()
	if currentChunk != requiredChunk && !c.chunkGraph.DependsOn(currentChunk, requiredChunk) {
		c.compiler.Report(
			JSError_MakeWithNode(
				call,
				CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR,
				namespace,
				requiredChunk.GetName(),
				currentChunk.GetName()))
	}
}
