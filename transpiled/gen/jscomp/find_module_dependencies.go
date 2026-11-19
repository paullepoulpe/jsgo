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
// Port of com/google/javascript/jscomp/FindModuleDependencies.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FindModuleDependencies
type FindModuleDependencies struct {
	compiler                AbstractCompiler                                  `gson:"compiler"`
	dynamicImportScope      *Scope                                            `gson:"dynamicImportScope"`
	inputPathByWebpackId    j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"inputPathByWebpackId"`
	idHash_                 atomic.Int32
	supportsEs6Modules      bool                     `gson:"supportsEs6Modules"`
	supportsCommonJsModules bool                     `gson:"supportsCommonJsModules"`
	moduleType              CompilerInput_ModuleType `gson:"moduleType"`
}

func (this_ *FindModuleDependencies) BindSelf_(s any) {
	_ = s
}
func (this_ *FindModuleDependencies) IsJavaNil_() bool { return this_ == nil }
func (*FindModuleDependencies) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FindModuleDependencies", "FindModuleDependencies"
}
func (this_ *FindModuleDependencies) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FindModuleDependencies) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FindModuleDependencies", "com.google.javascript.jscomp.FindModuleDependencies", func(v_ any) bool {
		return j2gort.InstanceOf[*FindModuleDependencies](v_)
	})
}
func (*FindModuleDependencies) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*FindModuleDependencies) Implements_jscomp_NodeTraversal_Callback()       {}

var static_FindModuleDependencies_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_FindModuleDependencies_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_FindModuleDependencies()
	return static_FindModuleDependencies_GOOG_MODULE
}
func StaticSet_FindModuleDependencies_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_FindModuleDependencies()
	static_FindModuleDependencies_GOOG_MODULE = v
}

var static_FindModuleDependencies_GOOG_PROVIDE pkg_rhino.QualifiedName

func StaticGet_FindModuleDependencies_GOOG_PROVIDE() pkg_rhino.QualifiedName {
	clinit_FindModuleDependencies()
	return static_FindModuleDependencies_GOOG_PROVIDE
}
func StaticSet_FindModuleDependencies_GOOG_PROVIDE(v pkg_rhino.QualifiedName) {
	clinit_FindModuleDependencies()
	static_FindModuleDependencies_GOOG_PROVIDE = v
}

var clinitDone_FindModuleDependencies uint32
var clinitState_FindModuleDependencies j2gort.ClinitState

func clinit_FindModuleDependencies() {
	if atomic.LoadUint32(&clinitDone_FindModuleDependencies) != 0 {
		return
	}
	clinitSlow_FindModuleDependencies()
}

func clinitSlow_FindModuleDependencies() {
	if !j2gort.EnterClinit(&clinitState_FindModuleDependencies, &clinitDone_FindModuleDependencies) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FindModuleDependencies, &clinitDone_FindModuleDependencies)
	static_FindModuleDependencies_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
	static_FindModuleDependencies_GOOG_PROVIDE = pkg_rhino.QualifiedName_Of("goog.provide")
}

func (this_ *FindModuleDependencies) initFields_FindModuleDependencies() {
	_ = this_
	this_.moduleType = StaticGet_CompilerInput_ModuleType_NONE()

}

// Java constructor: FindModuleDependencies(com.google.javascript.jscomp.AbstractCompiler compiler, boolean supportsEs6Modules, boolean supportsCommonJsModules, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> inputPathByWebpackId)
func NewFindModuleDependencies(compiler AbstractCompiler, supportsEs6Modules bool, supportsCommonJsModules bool, inputPathByWebpackId j2gort.Ext_com_google_common_collect_ImmutableMap) *FindModuleDependencies {
	this_ := &FindModuleDependencies{}
	InitFindModuleDependencies(this_, compiler, supportsEs6Modules, supportsCommonJsModules, inputPathByWebpackId)
	return this_
}

func InitFindModuleDependencies(this_ *FindModuleDependencies, compiler AbstractCompiler, supportsEs6Modules bool, supportsCommonJsModules bool, inputPathByWebpackId j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = compiler
	_ = supportsEs6Modules
	_ = supportsCommonJsModules
	_ = inputPathByWebpackId
	this_.initFields_FindModuleDependencies()
	this_.compiler = compiler
	this_.supportsEs6Modules = supportsEs6Modules
	this_.supportsCommonJsModules = supportsCommonJsModules
	this_.inputPathByWebpackId = inputPathByWebpackId
}

// Java method: FindModuleDependencies.process(com.google.javascript.rhino.Node root)
func (this_ *FindModuleDependencies) Process(root pkg_rhino.Node) {
	_ = root
	if !(root.IsScript()) {
		j2gort.CheckArgument(false)
	}
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		this_.moduleType = StaticGet_CompilerInput_ModuleType_ES6()
	}
	var input *CompilerInput = this_.compiler.GetInput(root.GetInputId())
	_ = input
	if input.GetRequires().Contains(j2gort.PtrToAny(StaticGet_DependencyInfo_Require_BASE())) {
		input.AddOrderedRequire(StaticGet_DependencyInfo_Require_BASE())
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
	if this_.moduleType == StaticGet_CompilerInput_ModuleType_ES6() {
		this_.priv_FindModuleDependencies_convertToEs6Module_2873(root, true)
	} else if ((this_.moduleType == StaticGet_CompilerInput_ModuleType_NONE()) && (this_.inputPathByWebpackId != nil)) && this_.inputPathByWebpackId.ContainsValue(j2gort.BoxString(input.GetPath().String())) {
		this_.moduleType = StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT()
	}
	input.AddProvide(input.GetPath().ToModuleName())
	input.SetJsModuleType(this_.moduleType)
	input.SetHasFullParseDependencyInfo(true)
}

// Java method: FindModuleDependencies.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FindModuleDependencies) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if (this_.supportsCommonJsModules && n.IsFunction()) && ProcessCommonJSModules_IsCommonJsDynamicImportCallback(n, this_.compiler.GetOptions().GetModuleResolutionMode()) {
		if this_.dynamicImportScope == nil {
			this_.dynamicImportScope = t.GetScope()
		}
	}
	return true
}

// Java method: FindModuleDependencies.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FindModuleDependencies) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var resolutionMode ModuleLoader_ResolutionMode = this_.compiler.GetOptions().GetModuleResolutionMode()
	_ = resolutionMode
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsExprResult() {
			var maybeGetProp pkg_rhino.Node = n.GetFirstFirstChild()
			_ = maybeGetProp
			if (maybeGetProp != 0) && (StaticGet_FindModuleDependencies_GOOG_PROVIDE().Matches(maybeGetProp) || StaticGet_FindModuleDependencies_GOOG_MODULE().Matches(maybeGetProp)) {
				this_.moduleType = StaticGet_CompilerInput_ModuleType_GOOG()
				return
			}
		}
	}
	if NodeUtil_IsGoogRequireDynamicCall(n) {
		t.GetInput().AddRequireDynamicImports(n.GetLastChild().GetString())
	}
	if this_.supportsEs6Modules && n.IsExport() {
		this_.moduleType = StaticGet_CompilerInput_ModuleType_ES6()
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		} else if n.HasTwoChildren() {
			this_.priv_FindModuleDependencies_addEs6ModuleImportToGraph_2871(t, n)
		}
	} else if this_.supportsEs6Modules && n.IsImport() {
		this_.moduleType = StaticGet_CompilerInput_ModuleType_ES6()
		this_.priv_FindModuleDependencies_addEs6ModuleImportToGraph_2871(t, n)
	} else if (this_.supportsEs6Modules && (n.RawToken_() == uint8(124))) && n.GetFirstChild().IsString() {
		var path string = n.GetFirstChild().GetString()
		_ = path
		var modulePath *ModuleLoader_ModulePath = t.GetInput().GetPath().ResolveJsModuleSilently(path)
		_ = modulePath
		if modulePath != nil {
			t.GetInput().AddDynamicRequire(modulePath.ToModuleName())
		}
	} else if this_.supportsCommonJsModules {
		if (this_.moduleType != StaticGet_CompilerInput_ModuleType_GOOG()) && ProcessCommonJSModules_IsCommonJsExport(t, n, resolutionMode) {
			this_.moduleType = StaticGet_CompilerInput_ModuleType_COMMONJS()
		} else if ProcessCommonJSModules_IsCommonJsImport(n, resolutionMode) {
			var path_2 string = ProcessCommonJSModules_GetCommonJsImportPath(n, resolutionMode)
			_ = path_2
			var modulePath_2 *ModuleLoader_ModulePath = t.GetInput().GetPath().ResolveJsModule(path_2, n.GetSourceFileName(), n.GetLineno(), n.GetCharno())
			_ = modulePath_2
			if modulePath_2 != nil {
				if (this_.dynamicImportScope != nil) || ProcessCommonJSModules_IsCommonJsDynamicImportCallback(NodeUtil_GetEnclosingFunction(n), resolutionMode) {
					t.GetInput().AddDynamicRequire(modulePath_2.ToModuleName())
				} else {
					t.GetInput().AddOrderedRequire(DependencyInfo_Require_CommonJs(modulePath_2.ToModuleName(), path_2))
				}
			}
		}
	}
	if ((((parent != 0) && (parent.IsExprResult() || !(t.InGlobalHoistScope()))) && NodeUtil_IsGoogRequireCall(n)) && (n.GetSecondChild() != 0)) && n.GetSecondChild().IsStringLit() {
		var namespace string = n.GetSecondChild().GetString()
		_ = namespace
		if javalang.StringStartsWith(j2gort.CheckNotNullString(namespace), "goog.") {
			t.GetInput().AddOrderedRequire(StaticGet_DependencyInfo_Require_BASE())
		}
		t.GetInput().AddOrderedRequire(DependencyInfo_Require_GoogRequireSymbol(namespace))
	}
}

// Java method: FindModuleDependencies.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *FindModuleDependencies) EnterScope(t *NodeTraversal) {
	_ = t
}

// Java method: FindModuleDependencies.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *FindModuleDependencies) ExitScope(t *NodeTraversal) {
	_ = t
	if t.GetScope() == this_.dynamicImportScope {
		this_.dynamicImportScope = nil
	}
}

// Java method: FindModuleDependencies.addEs6ModuleImportToGraph(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *FindModuleDependencies) priv_FindModuleDependencies_addEs6ModuleImportToGraph_2871(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var moduleName string = this_.priv_FindModuleDependencies_getEs6ModuleNameFromImportNode_2872(t, n)
	_ = moduleName
	if javalang.StringStartsWith(j2gort.CheckNotNullString(moduleName), "goog.") {
		t.GetInput().AddOrderedRequire(StaticGet_DependencyInfo_Require_BASE())
	}
	t.GetInput().AddOrderedRequire(DependencyInfo_Require_Es6Import(moduleName, n.GetLastChild().GetString()))
}

// Java method: FindModuleDependencies.getEs6ModuleNameFromImportNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *FindModuleDependencies) priv_FindModuleDependencies_getEs6ModuleNameFromImportNode_2872(t *NodeTraversal, n pkg_rhino.Node) string {
	_ = t
	_ = n
	var importName string = n.GetLastChild().GetString()
	_ = importName
	var isNamespaceImport bool = javalang.StringStartsWith(j2gort.CheckNotNullString(importName), "goog:")
	_ = isNamespaceImport
	if isNamespaceImport {
		return j2gort.JStrSubFrom(importName, j2gort.JStrLen("goog:"))
	} else {
		var modulePath *ModuleLoader_ModulePath = t.GetInput().GetPath().ResolveJsModule(importName, n.GetSourceFileName(), n.GetLineno(), n.GetCharno())
		_ = modulePath
		if modulePath == nil {
			modulePath = t.GetInput().GetPath().ResolveModuleAsPath(importName)
		}
		return modulePath.ToModuleName()
	}
}

// Java method: FindModuleDependencies.convertToEs6Module(com.google.javascript.rhino.Node root, boolean skipGoogProvideModuleCheck)
func (this_ *FindModuleDependencies) priv_FindModuleDependencies_convertToEs6Module_2873(root pkg_rhino.Node, skipGoogProvideModuleCheck bool) bool {
	_ = root
	_ = skipGoogProvideModuleCheck
	if Es6RewriteModules_IsEs6ModuleRoot(root) {
		return true
	}
	if !(skipGoogProvideModuleCheck) {
		var finder *FindModuleDependencies_FindGoogProvideOrGoogModule = NewFindModuleDependencies_FindGoogProvideOrGoogModule()
		_ = finder
		NodeTraversal_Traverse(this_.compiler, root, finder)
		if finder.IsFound() {
			return false
		}
	}
	var moduleNode pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_MODULE_BODY()).Srcref(root)
	_ = moduleNode
	moduleNode.AddChildrenToBack(root.RemoveChildren())
	root.AddChildToBack(moduleNode)
	return true
}

// Java class: com.google.javascript.jscomp.FindModuleDependencies.FindGoogProvideOrGoogModule
type FindModuleDependencies_FindGoogProvideOrGoogModule struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	found bool `gson:"found"`
}

func (this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) BindSelf_(s any) {
	_ = s
}
func (this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) IsJavaNil_() bool {
	return this_ == nil
}
func (*FindModuleDependencies_FindGoogProvideOrGoogModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FindModuleDependencies$FindGoogProvideOrGoogModule", "FindGoogProvideOrGoogModule"
}
func (this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*FindModuleDependencies_FindGoogProvideOrGoogModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FindModuleDependencies.FindGoogProvideOrGoogModule", "com.google.javascript.jscomp.FindModuleDependencies$FindGoogProvideOrGoogModule", func(v_ any) bool {
		return j2gort.InstanceOf[*FindModuleDependencies_FindGoogProvideOrGoogModule](v_)
	})
}
func (*FindModuleDependencies_FindGoogProvideOrGoogModule) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*FindModuleDependencies_FindGoogProvideOrGoogModule) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) initFields_FindModuleDependencies_FindGoogProvideOrGoogModule() {
	_ = this_
}

// Java constructor: FindGoogProvideOrGoogModule()
func NewFindModuleDependencies_FindGoogProvideOrGoogModule() *FindModuleDependencies_FindGoogProvideOrGoogModule {
	this_ := &FindModuleDependencies_FindGoogProvideOrGoogModule{}
	InitFindModuleDependencies_FindGoogProvideOrGoogModule(this_)
	return this_
}

func InitFindModuleDependencies_FindGoogProvideOrGoogModule(this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_FindModuleDependencies_FindGoogProvideOrGoogModule()
}

// Java method: FindGoogProvideOrGoogModule.isFound()
func (this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) IsFound() bool {
	return this_.found
}

// Java method: FindGoogProvideOrGoogModule.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FindModuleDependencies_FindGoogProvideOrGoogModule) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if this_.found {
		return false
	}
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsExprResult() {
			var maybeGetProp pkg_rhino.Node = n.GetFirstFirstChild()
			_ = maybeGetProp
			if (maybeGetProp != 0) && (StaticGet_FindModuleDependencies_GOOG_PROVIDE().Matches(maybeGetProp) || StaticGet_FindModuleDependencies_GOOG_MODULE().Matches(maybeGetProp)) {
				this_.found = true
				return false
			}
		}
		return true
	}
	return false
}
func (*FindModuleDependencies_FindGoogProvideOrGoogModule) NoOpVisit_()                 {}
func (*FindModuleDependencies_FindGoogProvideOrGoogModule) CanReleaseSyntacticScopes_() {}
