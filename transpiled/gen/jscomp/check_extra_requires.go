/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckExtraRequires.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_check_extra_requires_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_extra_requires_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_check_extra_requires_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return pkg_rhino.QualifiedName_Of(j2gort.CastString(mrp0_)) })
	lambdaSingleton_check_extra_requires_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		e := j2gort.Cast[pkg_rhino.JSTypeExpression](lp_4)
		_ = e
		return e.GetAllTypeNames().Stream()
	})
}

// Java class: com.google.javascript.jscomp.lint.CheckExtraRequires
type CheckExtraRequires struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler         AbstractCompiler                                  `gson:"compiler"`
	requires         j2gort.OrderedMap[string, pkg_rhino.Node]         `gson:"requires"`
	usages           j2gort.OrderedSet[string]                         `gson:"usages"`
	requiresToRemove j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"requiresToRemove"`
}

func (this_ *CheckExtraRequires) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckExtraRequires) IsJavaNil_() bool { return this_ == nil }
func (*CheckExtraRequires) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckExtraRequires", "CheckExtraRequires"
}
func (this_ *CheckExtraRequires) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckExtraRequires) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckExtraRequires", "com.google.javascript.jscomp.lint.CheckExtraRequires", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckExtraRequires](v_)
	})
}
func (*CheckExtraRequires) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckExtraRequires) Implements_jscomp_CompilerPass()                            {}
func (*CheckExtraRequires) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckExtraRequires_EXTRA_REQUIRE_WARNING *DiagnosticType

func StaticGet_CheckExtraRequires_EXTRA_REQUIRE_WARNING() *DiagnosticType {
	clinit_CheckExtraRequires()
	return static_CheckExtraRequires_EXTRA_REQUIRE_WARNING
}
func StaticSet_CheckExtraRequires_EXTRA_REQUIRE_WARNING(v *DiagnosticType) {
	clinit_CheckExtraRequires()
	static_CheckExtraRequires_EXTRA_REQUIRE_WARNING = v
}

var static_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CheckExtraRequires()
	return static_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES
}
func StaticSet_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CheckExtraRequires()
	static_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES = v
}

var static_CheckExtraRequires_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_CheckExtraRequires_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_CheckExtraRequires()
	return static_CheckExtraRequires_GOOG_MODULE_GET
}
func StaticSet_CheckExtraRequires_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_CheckExtraRequires()
	static_CheckExtraRequires_GOOG_MODULE_GET = v
}

var clinitDone_CheckExtraRequires uint32
var clinitState_CheckExtraRequires j2gort.ClinitState

func clinit_CheckExtraRequires() {
	if atomic.LoadUint32(&clinitDone_CheckExtraRequires) != 0 {
		return
	}
	clinitSlow_CheckExtraRequires()
}

func clinitSlow_CheckExtraRequires() {
	if !j2gort.EnterClinit(&clinitState_CheckExtraRequires, &clinitDone_CheckExtraRequires) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckExtraRequires, &clinitDone_CheckExtraRequires)
	static_CheckExtraRequires_EXTRA_REQUIRE_WARNING = DiagnosticType_Disabled("JSC_EXTRA_REQUIRE_WARNING", "extra require: ''{0}'' is never referenced in this file")
	static_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES = j2gort.ImmSetOfStr2("goog.testing.asserts", "goog.testing.jsunit")
	static_CheckExtraRequires_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
}

func (this_ *CheckExtraRequires) initFields_CheckExtraRequires() {
	_ = this_
	this_.requires = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.usages = j2gort.NewOrderedSet[string]()
}

// Java constructor: CheckExtraRequires(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableSet<java.lang.String> requiresToRemove)
func NewCheckExtraRequires(compiler AbstractCompiler, requiresToRemove j2gort.Ext_com_google_common_collect_ImmutableSet) *CheckExtraRequires {
	clinit_CheckExtraRequires()
	this_ := &CheckExtraRequires{}
	InitCheckExtraRequires(this_, compiler, requiresToRemove)
	return this_
}

func InitCheckExtraRequires(this_ *CheckExtraRequires, compiler AbstractCompiler, requiresToRemove j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = compiler
	_ = requiresToRemove
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckExtraRequires()
	this_.compiler = compiler
	this_.requiresToRemove = requiresToRemove
}

// Java method: CheckExtraRequires.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckExtraRequires) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.priv_CheckExtraRequires_reset_6495()
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckExtraRequires.extractNamespace(com.google.javascript.rhino.Node call, java.lang.String[] primitiveNames)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_extractNamespace_6492(call pkg_rhino.Node, primitiveNames []string) string {
	_ = call
	_ = primitiveNames
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	if !(callee.IsGetProp()) {
		return j2gort.NullString
	}
	for _, el_1 := range j2gort.CheckNotNullSlice(primitiveNames) {
		primitiveName := el_1
		_ = primitiveName
		if callee.MatchesQualifiedName_1_string(primitiveName) {
			var target pkg_rhino.Node = callee.GetNext()
			_ = target
			if (target != 0) && target.IsStringLit() {
				return target.GetString()
			}
		}
	}
	return j2gort.NullString
}

// Java method: CheckExtraRequires.extractNamespaceIfRequire(com.google.javascript.rhino.Node call)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_extractNamespaceIfRequire_6493(call pkg_rhino.Node) string {
	_ = call
	return this_.priv_CheckExtraRequires_extractNamespace_6492(call, []string{"goog.require", "goog.requireType"})
}

// Java method: CheckExtraRequires.extractNamespaceIfForwardDeclare(com.google.javascript.rhino.Node call)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_extractNamespaceIfForwardDeclare_6494(call pkg_rhino.Node) string {
	_ = call
	return this_.priv_CheckExtraRequires_extractNamespace_6492(call, []string{"goog.forwardDeclare"})
}

// Java method: CheckExtraRequires.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckExtraRequires) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_CheckExtraRequires_maybeAddJsDocUsages_6505(n)
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		if (!(NodeUtil_IsLValue(n)) && !(parent.IsGetProp())) && !(parent.IsImportSpec()) {
			this_.priv_CheckExtraRequires_visitQualifiedName_6504(n)
		}
	case int32(26) /* Token.GETPROP */ :
		if !(parent.IsGetProp()) && n.IsQualifiedName() {
			this_.priv_CheckExtraRequires_visitQualifiedName_6504(n)
		}
	case int32(28) /* Token.CALL */ :
		this_.priv_CheckExtraRequires_visitCallNode_6502(n, parent)
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_CheckExtraRequires_visitScriptNode_6496()
		this_.priv_CheckExtraRequires_reset_6495()
	case int32(116) /* Token.IMPORT */ :
		this_.priv_CheckExtraRequires_visitImportNode_6499(n)
	default:
	}
}

// Java method: CheckExtraRequires.reset()
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_reset_6495() {
	this_.usages.Clear()
	this_.requires.Clear()
}

// Java method: CheckExtraRequires.visitScriptNode()
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitScriptNode_6496() {
	{
		cur_2 := this_.requires.IterEntries()
		for cur_2.NextEntry() {
			var require string = cur_2.Key()
			_ = require
			var call pkg_rhino.Node = cur_2.Value()
			_ = call
			if !(this_.usages.Contains(require)) && ((this_.requiresToRemove == nil) || j2gort.SetContainsStr(this_.requiresToRemove, require)) {
				this_.priv_CheckExtraRequires_reportExtraRequireWarning_6497(call, require)
			}
		}
	}
}

// Java method: CheckExtraRequires.reportExtraRequireWarning(com.google.javascript.rhino.Node call, java.lang.String require)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_reportExtraRequireWarning_6497(call pkg_rhino.Node, require string) {
	_ = call
	_ = require
	if j2gort.SetContainsStr(StaticGet_CheckExtraRequires_DEFAULT_EXTRA_NAMESPACES(), require) {
		return
	}
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(call)
	_ = jsDoc
	if (jsDoc != nil) && j2gort.SetContainsStr(jsDoc.GetSuppressions(), "extraRequire") {
		return
	}
	this_.compiler.Report_1(JSError_Make_3(call, StaticGet_CheckExtraRequires_EXTRA_REQUIRE_WARNING(), []string{require}))
}

// Java method: CheckExtraRequires.visitRequire(java.lang.String localName, com.google.javascript.rhino.Node node)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitRequire_6498(localName string, node pkg_rhino.Node) {
	_ = localName
	_ = node
	this_.requires.PutIfAbsent(localName, node)
}

// Java method: CheckExtraRequires.visitImportNode(com.google.javascript.rhino.Node importNode)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitImportNode_6499(importNode pkg_rhino.Node) {
	_ = importNode
	var defaultImport pkg_rhino.Node = importNode.GetFirstChild()
	_ = defaultImport
	if defaultImport.IsName() {
		this_.priv_CheckExtraRequires_visitRequire_6498(defaultImport.GetString(), importNode)
	}
	var namedImports pkg_rhino.Node = defaultImport.GetNext()
	_ = namedImports
	if namedImports.IsImportSpecs() {
		{
			var importSpec pkg_rhino.Node = namedImports.GetFirstChild()
			_ = importSpec
			for ; importSpec != 0; importSpec = importSpec.GetNext() {
				this_.priv_CheckExtraRequires_visitRequire_6498(importSpec.GetLastChild().GetString(), importNode)
			}
		}
	}
}

// Java method: CheckExtraRequires.visitForwardDeclare(java.lang.String namespace, com.google.javascript.rhino.Node forwardDeclareCall, com.google.javascript.rhino.Node parent)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitForwardDeclare_6500(namespace string, forwardDeclareCall pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = namespace
	_ = forwardDeclareCall
	_ = parent
	this_.priv_CheckExtraRequires_visitGoogRequire_6501(namespace, forwardDeclareCall, parent)
}

// Java method: CheckExtraRequires.visitGoogRequire(java.lang.String namespace, com.google.javascript.rhino.Node googRequireCall, com.google.javascript.rhino.Node parent)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitGoogRequire_6501(namespace string, googRequireCall pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = namespace
	_ = googRequireCall
	_ = parent
	if parent.IsName() {
		this_.priv_CheckExtraRequires_visitRequire_6498(parent.GetString(), googRequireCall)
	} else if parent.IsDestructuringLhs() && parent.GetFirstChild().IsObjectPattern() {
		if parent.GetFirstChild().HasChildren() {
			{
				var stringKey pkg_rhino.Node = parent.GetFirstFirstChild()
				_ = stringKey
				for ; stringKey != 0; stringKey = stringKey.GetNext() {
					var importName pkg_rhino.Node = stringKey.GetFirstChild()
					_ = importName
					if !(importName.IsName()) {
						continue
					}
					this_.priv_CheckExtraRequires_visitRequire_6498(importName.GetString(), importName)
				}
			}
		} else {
			this_.priv_CheckExtraRequires_visitRequire_6498(namespace, googRequireCall)
		}
	} else {
		this_.priv_CheckExtraRequires_visitRequire_6498(namespace, googRequireCall)
	}
}

// Java method: CheckExtraRequires.visitCallNode(com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node parent)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitCallNode_6502(call pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = call
	_ = parent
	var required string = this_.priv_CheckExtraRequires_extractNamespaceIfRequire_6493(call)
	_ = required
	if required != j2gort.NullString {
		this_.priv_CheckExtraRequires_visitGoogRequire_6501(required, call, parent)
		return
	}
	var declare string = this_.priv_CheckExtraRequires_extractNamespaceIfForwardDeclare_6494(call)
	_ = declare
	if declare != j2gort.NullString {
		this_.priv_CheckExtraRequires_visitForwardDeclare_6500(declare, call, parent)
		return
	}
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	if StaticGet_CheckExtraRequires_GOOG_MODULE_GET().Matches(callee) && call.GetSecondChild().IsStringLit() {
		this_.usages.Add(call.GetSecondChild().GetString())
	}
}

// Java method: CheckExtraRequires.addUsagesOfAllPrefixes(com.google.javascript.rhino.QualifiedName qualifiedName)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_addUsagesOfAllPrefixes_6503(qualifiedName pkg_rhino.QualifiedName) {
	_ = qualifiedName
	this_.usages.Add(qualifiedName.Join_0())
	for !(qualifiedName.IsSimple()) {
		qualifiedName = qualifiedName.GetOwner()
		this_.usages.Add(qualifiedName.Join_0())
	}
}

// Java method: CheckExtraRequires.visitQualifiedName(com.google.javascript.rhino.Node n)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_visitQualifiedName_6504(n pkg_rhino.Node) {
	_ = n
	if !(n.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	this_.priv_CheckExtraRequires_addUsagesOfAllPrefixes_6503(n.GetQualifiedNameObject())
}

// Java method: CheckExtraRequires.maybeAddJsDocUsages(com.google.javascript.rhino.Node n)
func (this_ *CheckExtraRequires) priv_CheckExtraRequires_maybeAddJsDocUsages_6505(n pkg_rhino.Node) {
	_ = n
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info == nil {
		return
	}
	info.GetTypeExpressions().Stream().FlatMap(lambdaSingleton_check_extra_requires_2).Map(lambdaSingleton_check_extra_requires_1).ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_3 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_3.priv_CheckExtraRequires_addUsagesOfAllPrefixes_6503(j2gort.AsIface[pkg_rhino.QualifiedName](mrp0_))
		})
	}())
}
func (*CheckExtraRequires) AlwaysShouldTraverse_()      {}
func (*CheckExtraRequires) CanReleaseSyntacticScopes_() {}
