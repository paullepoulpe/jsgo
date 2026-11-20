/*
 * Copyright 2024 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ExtraRequireRemover.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ExtraRequireRemover
type ExtraRequireRemover struct {
	compiler AbstractCompiler                          `gson:"compiler"`
	requires j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"requires"`
	usages   j2gort.OrderedSet[string]                 `gson:"usages"`
	idHash_  atomic.Int32
}

func (this_ *ExtraRequireRemover) BindSelf_(s any) {
	_ = s
}
func (this_ *ExtraRequireRemover) IsJavaNil_() bool { return this_ == nil }
func (*ExtraRequireRemover) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExtraRequireRemover", "ExtraRequireRemover"
}
func (this_ *ExtraRequireRemover) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExtraRequireRemover) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExtraRequireRemover", "com.google.javascript.jscomp.ExtraRequireRemover", func(v_ any) bool {
		return j2gort.InstanceOf[*ExtraRequireRemover](v_)
	})
}
func (*ExtraRequireRemover) Implements_jscomp_CompilerPass()           {}
func (*ExtraRequireRemover) Implements_jscomp_NodeTraversal_Callback() {}

var static_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ExtraRequireRemover()
	return static_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES
}
func StaticSet_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ExtraRequireRemover()
	static_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES = v
}

var static_ExtraRequireRemover_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_ExtraRequireRemover_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_ExtraRequireRemover()
	return static_ExtraRequireRemover_GOOG_REQUIRE
}
func StaticSet_ExtraRequireRemover_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_ExtraRequireRemover()
	static_ExtraRequireRemover_GOOG_REQUIRE = v
}

var static_ExtraRequireRemover_GOOG_REQUIRE_TYPE pkg_rhino.QualifiedName

func StaticGet_ExtraRequireRemover_GOOG_REQUIRE_TYPE() pkg_rhino.QualifiedName {
	clinit_ExtraRequireRemover()
	return static_ExtraRequireRemover_GOOG_REQUIRE_TYPE
}
func StaticSet_ExtraRequireRemover_GOOG_REQUIRE_TYPE(v pkg_rhino.QualifiedName) {
	clinit_ExtraRequireRemover()
	static_ExtraRequireRemover_GOOG_REQUIRE_TYPE = v
}

var static_ExtraRequireRemover_GOOG_FORWARD_DECLARE pkg_rhino.QualifiedName

func StaticGet_ExtraRequireRemover_GOOG_FORWARD_DECLARE() pkg_rhino.QualifiedName {
	clinit_ExtraRequireRemover()
	return static_ExtraRequireRemover_GOOG_FORWARD_DECLARE
}
func StaticSet_ExtraRequireRemover_GOOG_FORWARD_DECLARE(v pkg_rhino.QualifiedName) {
	clinit_ExtraRequireRemover()
	static_ExtraRequireRemover_GOOG_FORWARD_DECLARE = v
}

var static_ExtraRequireRemover_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_ExtraRequireRemover_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_ExtraRequireRemover()
	return static_ExtraRequireRemover_GOOG_MODULE_GET
}
func StaticSet_ExtraRequireRemover_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_ExtraRequireRemover()
	static_ExtraRequireRemover_GOOG_MODULE_GET = v
}

var clinitDone_ExtraRequireRemover uint32
var clinitState_ExtraRequireRemover j2gort.ClinitState

func clinit_ExtraRequireRemover() {
	if atomic.LoadUint32(&clinitDone_ExtraRequireRemover) != 0 {
		return
	}
	clinitSlow_ExtraRequireRemover()
}

func clinitSlow_ExtraRequireRemover() {
	if !j2gort.EnterClinit(&clinitState_ExtraRequireRemover, &clinitDone_ExtraRequireRemover) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ExtraRequireRemover, &clinitDone_ExtraRequireRemover)
	static_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("goog.testing.asserts", "goog.testing.jsunit", "goog.labs.testing.Environment")
	static_ExtraRequireRemover_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_ExtraRequireRemover_GOOG_REQUIRE_TYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_ExtraRequireRemover_GOOG_FORWARD_DECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_ExtraRequireRemover_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
}

func (this_ *ExtraRequireRemover) initFields_ExtraRequireRemover() {
	_ = this_
	this_.requires = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.usages = j2gort.NewOrderedSet[string]()
}

// Java constructor: ExtraRequireRemover(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewExtraRequireRemover(compiler AbstractCompiler) *ExtraRequireRemover {
	this_ := &ExtraRequireRemover{}
	InitExtraRequireRemover(this_, compiler)
	return this_
}

func InitExtraRequireRemover(this_ *ExtraRequireRemover, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ExtraRequireRemover()
	this_.compiler = compiler
}

// Java method: ExtraRequireRemover.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ExtraRequireRemover) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.priv_ExtraRequireRemover_reset_2861()
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: ExtraRequireRemover.extractNamespace(com.google.javascript.rhino.Node call, com.google.javascript.rhino.QualifiedName[] primitiveNames)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_extractNamespace_2859(call pkg_rhino.Node, primitiveNames []pkg_rhino.QualifiedName) string {
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
		if primitiveName.Matches(callee) {
			var target pkg_rhino.Node = callee.GetNext()
			_ = target
			if (target != 0) && target.IsStringLit() {
				return target.GetString()
			}
		}
	}
	return j2gort.NullString
}

// Java method: ExtraRequireRemover.extractNamespaceIfRequire(com.google.javascript.rhino.Node call)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_extractNamespaceIfRequire_2860(call pkg_rhino.Node) string {
	_ = call
	return this_.priv_ExtraRequireRemover_extractNamespace_2859(call, []pkg_rhino.QualifiedName{StaticGet_ExtraRequireRemover_GOOG_REQUIRE(), StaticGet_ExtraRequireRemover_GOOG_REQUIRE_TYPE(), StaticGet_ExtraRequireRemover_GOOG_FORWARD_DECLARE()})
}

// Java method: ExtraRequireRemover.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ExtraRequireRemover) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsModuleBody()) || n.GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE()))
}

// Java method: ExtraRequireRemover.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ExtraRequireRemover) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_ExtraRequireRemover_maybeAddJsDocUsages_2869(n)
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		if !(NodeUtil_IsLValue(n)) && !(parent.IsGetProp()) {
			this_.priv_ExtraRequireRemover_visitQualifiedName_2868(n)
		}
	case int32(26) /* Token.GETPROP */ :
		if !(parent.IsGetProp()) && n.IsQualifiedName() {
			this_.priv_ExtraRequireRemover_visitQualifiedName_2868(n)
		}
	case int32(28) /* Token.CALL */ :
		this_.priv_ExtraRequireRemover_visitCallNode_2866(n, parent)
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_ExtraRequireRemover_removeExtraRequiresInScript_2862()
		this_.priv_ExtraRequireRemover_reset_2861()
	default:
	}
}

// Java method: ExtraRequireRemover.reset()
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_reset_2861() {
	this_.usages.Clear()
	this_.requires.Clear()
}

// Java method: ExtraRequireRemover.removeExtraRequiresInScript()
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_removeExtraRequiresInScript_2862() {
	{
		cur_2 := this_.requires.IterEntries()
		for cur_2.NextEntry() {
			var require string = cur_2.Key()
			_ = require
			var call pkg_rhino.Node = cur_2.Value()
			_ = call
			this_.priv_ExtraRequireRemover_removeExtraRequire_2863(call, require)
		}
	}
}

// Java method: ExtraRequireRemover.removeExtraRequire(com.google.javascript.rhino.Node call, java.lang.String require)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_removeExtraRequire_2863(call pkg_rhino.Node, require string) {
	_ = call
	_ = require
	if this_.usages.Contains(require) {
		return
	}
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(call)
	_ = jsDoc
	if (jsDoc != nil) && j2gort.SetContainsStr(jsDoc.GetSuppressions(), "extraRequire") {
		return
	}
	if ExtraRequireRemover_IsTsickleProtoCall(call) {
		return
	}
	if call.GetParent().IsExprResult() {
		if !(call.GetGrandparent().IsModuleBody()) || J2clSourceUtils_IsJ2clHeaderSource(call.GetSourceFileName()) {
			this_.compiler.ReportChangeToEnclosingScope(call.GetParent())
			call.GetParent().Detach()
		}
		return
	}
	if call.GetParent().IsName() {
		var nameDeclaration pkg_rhino.Node = call.GetGrandparent()
		_ = nameDeclaration
		if !(NodeUtil_IsNameDeclaration(call.GetGrandparent())) {
			j2gort.CheckState(false)
		}
		this_.compiler.ReportChangeToEnclosingScope(nameDeclaration.GetParent())
		nameDeclaration.Detach()
		return
	}
	if !(call.GetParent().IsStringKey()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(call.GetParent()))
	}
	var objectPattern pkg_rhino.Node = call.GetGrandparent()
	_ = objectPattern
	if objectPattern.HasOneChild() {
		var nameDeclaration_2 pkg_rhino.Node = objectPattern.GetGrandparent()
		_ = nameDeclaration_2
		if !(NodeUtil_IsNameDeclaration(nameDeclaration_2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(nameDeclaration_2))
		}
		this_.compiler.ReportChangeToEnclosingScope(nameDeclaration_2.GetParent())
		nameDeclaration_2.Detach()
		return
	}
	this_.compiler.ReportChangeToEnclosingScope(call.GetGrandparent())
	call.GetParent().Detach()
}

// Java method: ExtraRequireRemover.visitRequire(java.lang.String localName, java.lang.String namespace, com.google.javascript.rhino.Node node)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_visitRequire_2864(localName string, namespace string, node pkg_rhino.Node) {
	_ = localName
	_ = namespace
	_ = node
	if j2gort.SetContainsStr(StaticGet_ExtraRequireRemover_DEFAULT_EXTRA_NAMESPACES(), namespace) {
		return
	}
	this_.requires.PutIfAbsent(localName, node)
}

// Java method: ExtraRequireRemover.visitGoogRequire(java.lang.String namespace, com.google.javascript.rhino.Node googRequireCall, com.google.javascript.rhino.Node parent)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_visitGoogRequire_2865(namespace string, googRequireCall pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = namespace
	_ = googRequireCall
	_ = parent
	if parent.IsName() {
		this_.priv_ExtraRequireRemover_visitRequire_2864(parent.GetString(), namespace, googRequireCall)
	} else if parent.IsDestructuringLhs() && parent.GetFirstChild().IsObjectPattern() {
		if parent.GetFirstChild().HasChildren() {
			{
				var stringKey pkg_rhino.Node = parent.GetFirstFirstChild()
				_ = stringKey
				for ; stringKey != 0; stringKey = stringKey.GetNext() {
					var importName pkg_rhino.Node = stringKey.GetFirstChild()
					_ = importName
					this_.priv_ExtraRequireRemover_visitRequire_2864(importName.GetString(), namespace, importName)
				}
			}
		} else {
			this_.priv_ExtraRequireRemover_visitRequire_2864(namespace, namespace, googRequireCall)
		}
	} else {
		this_.priv_ExtraRequireRemover_visitRequire_2864(namespace, namespace, googRequireCall)
	}
}

// Java method: ExtraRequireRemover.visitCallNode(com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node parent)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_visitCallNode_2866(call pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = call
	_ = parent
	var required string = this_.priv_ExtraRequireRemover_extractNamespaceIfRequire_2860(call)
	_ = required
	if required != j2gort.NullString {
		this_.priv_ExtraRequireRemover_visitGoogRequire_2865(required, call, parent)
		return
	}
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	if StaticGet_ExtraRequireRemover_GOOG_MODULE_GET().Matches(callee) && call.GetSecondChild().IsStringLit() {
		this_.usages.Add(call.GetSecondChild().GetString())
	}
}

// Java method: ExtraRequireRemover.addUsagesOfAllPrefixes(com.google.javascript.rhino.QualifiedName qualifiedName)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_addUsagesOfAllPrefixes_2867(qualifiedName pkg_rhino.QualifiedName) {
	_ = qualifiedName
	this_.usages.Add(qualifiedName.Join_0())
	for !(qualifiedName.IsSimple()) {
		qualifiedName = qualifiedName.GetOwner()
		this_.usages.Add(qualifiedName.Join_0())
	}
}

// Java method: ExtraRequireRemover.visitQualifiedName(com.google.javascript.rhino.Node n)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_visitQualifiedName_2868(n pkg_rhino.Node) {
	_ = n
	if !(n.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	this_.priv_ExtraRequireRemover_addUsagesOfAllPrefixes_2867(n.GetQualifiedNameObject())
}

// Java method: ExtraRequireRemover.maybeAddJsDocUsages(com.google.javascript.rhino.Node n)
func (this_ *ExtraRequireRemover) priv_ExtraRequireRemover_maybeAddJsDocUsages_2869(n pkg_rhino.Node) {
	_ = n
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info == nil {
		return
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetTypeExpressions()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			e := j2gort.Cast[pkg_rhino.JSTypeExpression](rawElem_4)
			_ = e
			{
				cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(e.GetAllTypeNames()))
				for cur_5.HasNext() {
					typeName := cur_5.NextStr()
					_ = typeName
					this_.priv_ExtraRequireRemover_addUsagesOfAllPrefixes_2867(pkg_rhino.QualifiedName_Of(typeName))
				}
			}
		}
	}
}

// Java static method: ExtraRequireRemover.isTsickleProtoCall(com.google.javascript.rhino.Node call)
func ExtraRequireRemover_IsTsickleProtoCall(call pkg_rhino.Node) bool {
	_ = call
	if !(call.IsCall()) || !(javalang.StringEndsWith(j2gort.CheckNotNullString(call.GetSourceFileName()), ".closure.js")) {
		return false
	}
	var moduleId pkg_rhino.Node = call.GetSecondChild()
	_ = moduleId
	return (((moduleId != 0) && moduleId.IsStringLit()) && javalang.StringStartsWith(j2gort.CheckNotNullString(moduleId.GetString()), "proto."))
}
func (*ExtraRequireRemover) CanReleaseSyntacticScopes_() {}
