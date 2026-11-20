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
// Port of com/google/javascript/refactoring/ScriptMetadata.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.ScriptMetadata
type ScriptMetadata struct {
	requireAliases         j2gort.OrderedMap[string, string] `gson:"requireAliases"`
	namesInUse             j2gort.OrderedSet[string]         `gson:"namesInUse"`
	script                 pkg_rhino.Node                    `gson:"script"`
	idHash_                atomic.Int32
	supportsRequireAliases bool `gson:"supportsRequireAliases"`
}

func (this_ *ScriptMetadata) BindSelf_(s any) {
	_ = s
}
func (this_ *ScriptMetadata) IsJavaNil_() bool { return this_ == nil }
func (*ScriptMetadata) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.ScriptMetadata", "ScriptMetadata"
}
func (this_ *ScriptMetadata) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ScriptMetadata) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.ScriptMetadata", "com.google.javascript.refactoring.ScriptMetadata", func(v_ any) bool {
		return j2gort.InstanceOf[*ScriptMetadata](v_)
	})
}

var static_ScriptMetadata_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_ScriptMetadata_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_ScriptMetadata()
	return static_ScriptMetadata_DOT_SPLITTER
}
func StaticSet_ScriptMetadata_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_ScriptMetadata()
	static_ScriptMetadata_DOT_SPLITTER = v
}

var clinitDone_ScriptMetadata uint32
var clinitState_ScriptMetadata j2gort.ClinitState

func clinit_ScriptMetadata() {
	if atomic.LoadUint32(&clinitDone_ScriptMetadata) != 0 {
		return
	}
	clinitSlow_ScriptMetadata()
}

func clinitSlow_ScriptMetadata() {
	if !j2gort.EnterClinit(&clinitState_ScriptMetadata, &clinitDone_ScriptMetadata) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ScriptMetadata, &clinitDone_ScriptMetadata)
	static_ScriptMetadata_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46))
}

func (this_ *ScriptMetadata) initFields_ScriptMetadata() {
	_ = this_
	this_.requireAliases = j2gort.NewOrderedMap[string, string]()
	this_.namesInUse = j2gort.NewOrderedSet[string]()
}

// Java constructor: ScriptMetadata(com.google.javascript.rhino.Node script)
func NewScriptMetadata(script pkg_rhino.Node) *ScriptMetadata {
	clinit_ScriptMetadata()
	this_ := &ScriptMetadata{}
	InitScriptMetadata(this_, script)
	return this_
}

func InitScriptMetadata(this_ *ScriptMetadata, script pkg_rhino.Node) {
	_ = script
	this_.initFields_ScriptMetadata()
	this_.script = script
}

// Java static method: ScriptMetadata.create(com.google.javascript.rhino.Node script, com.google.javascript.refactoring.NodeMetadata metadata)
func ScriptMetadata_Create_2_Node_NodeMetadata(script pkg_rhino.Node, metadata *NodeMetadata) *ScriptMetadata {
	clinit_ScriptMetadata()
	_ = script
	_ = metadata
	return ScriptMetadata_Create_2_Node_AbstractCompiler(script, metadata.GetCompiler())
}

// Java static method: ScriptMetadata.create(com.google.javascript.rhino.Node script, com.google.javascript.jscomp.AbstractCompiler compiler)
func ScriptMetadata_Create_2_Node_AbstractCompiler(script pkg_rhino.Node, compiler AbstractCompiler) *ScriptMetadata {
	clinit_ScriptMetadata()
	_ = script
	_ = compiler
	if !(script.IsScript()) {
		j2gort.CheckArgument(false)
	}
	clinit_ScriptMetadata()
	var val_1 ScriptMetadata
	InitScriptMetadata(&val_1, script)
	var toFill *ScriptMetadata = &val_1
	_ = toFill
	NodeTraversal_Traverse(compiler, script, NewScriptMetadata_FillerCallback(toFill))
	if !(toFill.requireAliases.IsEmpty()) {
		{
			r_2 := toFill
			r_2.supportsRequireAliases = true
		}
	}
	return toFill
}

// Java method: ScriptMetadata.getScript()
func (this_ *ScriptMetadata) GetScript() pkg_rhino.Node {
	return this_.script
}

// Java method: ScriptMetadata.supportsRequireAliases()
func (this_ *ScriptMetadata) SupportsRequireAliases() bool {
	return this_.supportsRequireAliases
}

// Java method: ScriptMetadata.usesName(java.lang.String name)
func (this_ *ScriptMetadata) UsesName(name string) bool {
	_ = name
	return this_.namesInUse.Contains(name)
}

// Java method: ScriptMetadata.getAlias(java.lang.String namespace)
func (this_ *ScriptMetadata) GetAlias(namespace string) string {
	_ = namespace
	return j2gort.OrderedMapGetStr(this_.requireAliases, namespace)
}

// Java method: ScriptMetadata.addAlias(java.lang.String namespace, java.lang.String alias)
func (this_ *ScriptMetadata) AddAlias(namespace string, alias string) {
	_ = namespace
	_ = alias
	if !(this_.supportsRequireAliases) {
		j2gort.CheckState(false)
	}
	this_.namesInUse.Add(alias)
	this_.requireAliases.Put(namespace, alias)
}

// Java class: com.google.javascript.refactoring.ScriptMetadata.FillerCallback
type ScriptMetadata_FillerCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	toFill *ScriptMetadata `gson:"toFill"`
}

func (this_ *ScriptMetadata_FillerCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *ScriptMetadata_FillerCallback) IsJavaNil_() bool { return this_ == nil }
func (*ScriptMetadata_FillerCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.ScriptMetadata$FillerCallback", "FillerCallback"
}
func (this_ *ScriptMetadata_FillerCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ScriptMetadata_FillerCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.ScriptMetadata.FillerCallback", "com.google.javascript.refactoring.ScriptMetadata$FillerCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*ScriptMetadata_FillerCallback](v_)
	})
}
func (*ScriptMetadata_FillerCallback) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ScriptMetadata_FillerCallback) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ScriptMetadata_FillerCallback_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_ScriptMetadata_FillerCallback_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_ScriptMetadata_FillerCallback()
	return static_ScriptMetadata_FillerCallback_GOOG_MODULE
}
func StaticSet_ScriptMetadata_FillerCallback_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_ScriptMetadata_FillerCallback()
	static_ScriptMetadata_FillerCallback_GOOG_MODULE = v
}

var static_ScriptMetadata_FillerCallback_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_ScriptMetadata_FillerCallback_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_ScriptMetadata_FillerCallback()
	return static_ScriptMetadata_FillerCallback_GOOG_REQUIRE
}
func StaticSet_ScriptMetadata_FillerCallback_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_ScriptMetadata_FillerCallback()
	static_ScriptMetadata_FillerCallback_GOOG_REQUIRE = v
}

var static_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE pkg_rhino.QualifiedName

func StaticGet_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE() pkg_rhino.QualifiedName {
	clinit_ScriptMetadata_FillerCallback()
	return static_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE
}
func StaticSet_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE(v pkg_rhino.QualifiedName) {
	clinit_ScriptMetadata_FillerCallback()
	static_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE = v
}

var static_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE pkg_rhino.QualifiedName

func StaticGet_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE() pkg_rhino.QualifiedName {
	clinit_ScriptMetadata_FillerCallback()
	return static_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE
}
func StaticSet_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE(v pkg_rhino.QualifiedName) {
	clinit_ScriptMetadata_FillerCallback()
	static_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE = v
}

var clinitDone_ScriptMetadata_FillerCallback uint32
var clinitState_ScriptMetadata_FillerCallback j2gort.ClinitState

func clinit_ScriptMetadata_FillerCallback() {
	if atomic.LoadUint32(&clinitDone_ScriptMetadata_FillerCallback) != 0 {
		return
	}
	clinitSlow_ScriptMetadata_FillerCallback()
}

func clinitSlow_ScriptMetadata_FillerCallback() {
	if !j2gort.EnterClinit(&clinitState_ScriptMetadata_FillerCallback, &clinitDone_ScriptMetadata_FillerCallback) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ScriptMetadata_FillerCallback, &clinitDone_ScriptMetadata_FillerCallback)
	static_ScriptMetadata_FillerCallback_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
	static_ScriptMetadata_FillerCallback_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
}

func (this_ *ScriptMetadata_FillerCallback) initFields_ScriptMetadata_FillerCallback() {
	_ = this_
}

// Java constructor: FillerCallback(com.google.javascript.refactoring.ScriptMetadata toFill)
func NewScriptMetadata_FillerCallback(toFill *ScriptMetadata) *ScriptMetadata_FillerCallback {
	this_ := &ScriptMetadata_FillerCallback{}
	InitScriptMetadata_FillerCallback(this_, toFill)
	return this_
}

func InitScriptMetadata_FillerCallback(this_ *ScriptMetadata_FillerCallback, toFill *ScriptMetadata) {
	_ = toFill
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ScriptMetadata_FillerCallback()
	this_.toFill = toFill
}

// Java method: FillerCallback.visit(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ScriptMetadata_FillerCallback) Visit_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = nodeTraversal
	_ = n
	_ = parent
	this_.priv_FillerCallback_updateSupportsRequireAlias_7992(n)
	this_.priv_FillerCallback_maybeCollectRequirelikeAlias_7993(n)
	this_.priv_FillerCallback_maybeCollectNames_7995(n)
}

// Java method: FillerCallback.updateSupportsRequireAlias(com.google.javascript.rhino.Node n)
func (this_ *ScriptMetadata_FillerCallback) priv_FillerCallback_updateSupportsRequireAlias_7992(n pkg_rhino.Node) {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(123) /* Token.MODULE_BODY */ :
	case int32(28) /* Token.CALL */ :
		if ((n.GetChildCount() != int32(2)) || !(StaticGet_ScriptMetadata_FillerCallback_GOOG_MODULE().Matches(n.GetFirstChild()))) || !(n.GetSecondChild().IsStringLit()) {
			return
		}
	default:
		return
	}
	{
		r_3 := this_.toFill
		r_3.supportsRequireAliases = true
	}
}

// Java method: FillerCallback.maybeCollectRequirelikeAlias(com.google.javascript.rhino.Node n)
func (this_ *ScriptMetadata_FillerCallback) priv_FillerCallback_maybeCollectRequirelikeAlias_7993(n pkg_rhino.Node) {
	_ = n
	if !(NodeUtil_IsNameDeclaration(n)) {
		return
	}
	var aliasNode pkg_rhino.Node = n.GetFirstChild()
	_ = aliasNode
	if aliasNode.IsDestructuringLhs() {
		return
	}
	var requirelikeCall pkg_rhino.Node = aliasNode.GetFirstChild()
	_ = requirelikeCall
	if !(this_.priv_FillerCallback_isRequirelikeCall_7994(requirelikeCall)) {
		return
	}
	var alias string = aliasNode.GetOriginalName()
	_ = alias
	if alias == j2gort.NullString {
		alias = aliasNode.GetString()
	}
	var namespace string = requirelikeCall.GetSecondChild().GetString()
	_ = namespace
	this_.toFill.requireAliases.Put(namespace, alias)
}

// Java method: FillerCallback.isRequirelikeCall(com.google.javascript.rhino.Node call)
func (this_ *ScriptMetadata_FillerCallback) priv_FillerCallback_isRequirelikeCall_7994(call pkg_rhino.Node) bool {
	_ = call
	if (((call == 0) || !(call.IsCall())) || (call.GetChildCount() != int32(2))) || !(call.GetSecondChild().IsStringLit()) {
		return false
	}
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	return ((StaticGet_ScriptMetadata_FillerCallback_GOOG_REQUIRE().Matches(callee) || StaticGet_ScriptMetadata_FillerCallback_GOOG_REQUIRETYPE().Matches(callee)) || StaticGet_ScriptMetadata_FillerCallback_GOOG_FORWARDDECLARE().Matches(callee))
}

// Java method: FillerCallback.maybeCollectNames(com.google.javascript.rhino.Node n)
func (this_ *ScriptMetadata_FillerCallback) priv_FillerCallback_maybeCollectNames_7995(n pkg_rhino.Node) {
	_ = n
	if n.IsName() {
		this_.toFill.namesInUse.Add(n.GetString())
	}
	var jsdoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsdoc
	if jsdoc != nil {
		{
			cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(jsdoc.GetTypeExpressions()))
			for cur_4.HasNext() {
				rawElem_5 := cur_4.Next()
				expr := j2gort.Cast[pkg_rhino.JSTypeExpression](rawElem_5)
				_ = expr
				{
					cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(expr.GetAllTypeNames()))
					for cur_6.HasNext() {
						type_ := cur_6.NextStr()
						_ = type_
						this_.toFill.namesInUse.Add(j2gort.CastString(j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.SplitterSplitStr(StaticGet_ScriptMetadata_DOT_SPLITTER(), type_).Iterator()).Next()))
					}
				}
			}
		}
	}
}
func (*ScriptMetadata_FillerCallback) AlwaysShouldTraverse_()      {}
func (*ScriptMetadata_FillerCallback) CanReleaseSyntacticScopes_() {}
