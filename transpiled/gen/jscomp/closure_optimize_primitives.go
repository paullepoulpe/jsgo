/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ClosureOptimizePrimitives.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_dtoa "github.com/paullepoulpe/jsgo/transpiled/gen/dtoa"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ClosureOptimizePrimitives
type ClosureOptimizePrimitives struct {
	compiler        AbstractCompiler `gson:"compiler"`
	idHash_         atomic.Int32
	canUseEs6Syntax bool `gson:"canUseEs6Syntax"`
}

func (this_ *ClosureOptimizePrimitives) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureOptimizePrimitives) IsJavaNil_() bool { return this_ == nil }
func (*ClosureOptimizePrimitives) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureOptimizePrimitives", "ClosureOptimizePrimitives"
}
func (this_ *ClosureOptimizePrimitives) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureOptimizePrimitives) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureOptimizePrimitives", "com.google.javascript.jscomp.ClosureOptimizePrimitives", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureOptimizePrimitives](v_)
	})
}
func (*ClosureOptimizePrimitives) Implements_jscomp_CompilerPass() {}

var static_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER *DiagnosticType

func StaticGet_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER() *DiagnosticType {
	clinit_ClosureOptimizePrimitives()
	return static_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER
}
func StaticSet_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER(v *DiagnosticType) {
	clinit_ClosureOptimizePrimitives()
	static_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER = v
}

var static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE pkg_rhino.QualifiedName

func StaticGet_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE() pkg_rhino.QualifiedName {
	clinit_ClosureOptimizePrimitives()
	return static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE
}
func StaticSet_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE(v pkg_rhino.QualifiedName) {
	clinit_ClosureOptimizePrimitives()
	static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE = v
}

var static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET pkg_rhino.QualifiedName

func StaticGet_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET() pkg_rhino.QualifiedName {
	clinit_ClosureOptimizePrimitives()
	return static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET
}
func StaticSet_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET(v pkg_rhino.QualifiedName) {
	clinit_ClosureOptimizePrimitives()
	static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET = v
}

var clinitDone_ClosureOptimizePrimitives uint32
var clinitState_ClosureOptimizePrimitives j2gort.ClinitState

func clinit_ClosureOptimizePrimitives() {
	if atomic.LoadUint32(&clinitDone_ClosureOptimizePrimitives) != 0 {
		return
	}
	clinitSlow_ClosureOptimizePrimitives()
}

func clinitSlow_ClosureOptimizePrimitives() {
	if !j2gort.EnterClinit(&clinitState_ClosureOptimizePrimitives, &clinitDone_ClosureOptimizePrimitives) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureOptimizePrimitives, &clinitDone_ClosureOptimizePrimitives)
	static_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER = DiagnosticType_Warning("JSC_DUPLICATE_SET_MEMBER", "Found duplicate value ''{0}'' in set")
	static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE = pkg_rhino.QualifiedName_Of("goog.object.create")
	static_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET = pkg_rhino.QualifiedName_Of("goog.object.createSet")
}

func (this_ *ClosureOptimizePrimitives) initFields_ClosureOptimizePrimitives() {
	_ = this_
}

// Java constructor: ClosureOptimizePrimitives(com.google.javascript.jscomp.AbstractCompiler compiler, boolean canUseEs6Syntax)
func NewClosureOptimizePrimitives(compiler AbstractCompiler, canUseEs6Syntax bool) *ClosureOptimizePrimitives {
	clinit_ClosureOptimizePrimitives()
	this_ := &ClosureOptimizePrimitives{}
	InitClosureOptimizePrimitives(this_, compiler, canUseEs6Syntax)
	return this_
}

func InitClosureOptimizePrimitives(this_ *ClosureOptimizePrimitives, compiler AbstractCompiler, canUseEs6Syntax bool) {
	_ = compiler
	_ = canUseEs6Syntax
	this_.initFields_ClosureOptimizePrimitives()
	this_.compiler = compiler
	this_.canUseEs6Syntax = canUseEs6Syntax
}

// Java method: ClosureOptimizePrimitives.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ClosureOptimizePrimitives) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var pass *ClosureOptimizePrimitives_FindPrimitives = NewClosureOptimizePrimitives_FindPrimitives(this_)
	_ = pass
	NodeTraversal_Traverse(this_.compiler, root, pass)
}

// Java method: ClosureOptimizePrimitives.processObjectCreateCall(com.google.javascript.rhino.Node callNode)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_processObjectCreateCall_1174(callNode pkg_rhino.Node) {
	_ = callNode
	var curParam pkg_rhino.Node = callNode.GetSecondChild()
	_ = curParam
	if this_.priv_ClosureOptimizePrimitives_canOptimizeObjectCreate_1176(curParam) {
		var objNode pkg_rhino.Node = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}).Srcref(callNode)
		_ = objNode
		for curParam != 0 {
			var keyNode pkg_rhino.Node = curParam
			_ = keyNode
			var valueNode pkg_rhino.Node = curParam.GetNext()
			_ = valueNode
			curParam = valueNode.GetNext()
			keyNode.Detach()
			valueNode.Detach()
			this_.priv_ClosureOptimizePrimitives_addKeyValueToObjLit_1179(objNode, keyNode, valueNode, NodeUtil_GetEnclosingScript(callNode))
		}
		callNode.ReplaceWith(objNode)
		this_.compiler.ReportChangeToEnclosingScope(objNode)
	}
}

// Java method: ClosureOptimizePrimitives.processRenamePropertyCall(com.google.javascript.rhino.Node callNode)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_processRenamePropertyCall_1175(callNode pkg_rhino.Node) {
	_ = callNode
	var nameNode pkg_rhino.Node = callNode.GetFirstChild()
	_ = nameNode
	if nameNode.MatchesQualifiedName_1_string("JSCompiler_renameProperty") {
		return
	}
	var newTarget pkg_rhino.Node = pkg_rhino.IR_Name("JSCompiler_renameProperty").Srcref(nameNode)
	_ = newTarget
	newTarget.SetOriginalName(nameNode.GetOriginalQualifiedName())
	nameNode.ReplaceWith(newTarget)
	callNode.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	this_.compiler.ReportChangeToEnclosingScope(callNode)
}

// Java method: ClosureOptimizePrimitives.canOptimizeObjectCreate(com.google.javascript.rhino.Node firstParam)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_canOptimizeObjectCreate_1176(firstParam pkg_rhino.Node) bool {
	_ = firstParam
	var curParam pkg_rhino.Node = firstParam
	_ = curParam
	for curParam != 0 {
		if !(this_.priv_ClosureOptimizePrimitives_isOptimizableKey_1180(curParam)) {
			return false
		}
		curParam = curParam.GetNext()
		if curParam == 0 {
			return false
		}
		curParam = curParam.GetNext()
	}
	return true
}

// Java method: ClosureOptimizePrimitives.processObjectCreateSetCall(com.google.javascript.rhino.Node callNode)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_processObjectCreateSetCall_1177(callNode pkg_rhino.Node) {
	_ = callNode
	var curParam pkg_rhino.Node = callNode.GetSecondChild()
	_ = curParam
	if this_.priv_ClosureOptimizePrimitives_canOptimizeObjectCreateSet_1178(curParam) {
		var objNode pkg_rhino.Node = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}).Srcref(callNode)
		_ = objNode
		for curParam != 0 {
			var keyNode pkg_rhino.Node = curParam
			_ = keyNode
			var valueNode pkg_rhino.Node = pkg_rhino.IR_TrueNode().Srcref(keyNode)
			_ = valueNode
			curParam = curParam.GetNext()
			keyNode.Detach()
			this_.priv_ClosureOptimizePrimitives_addKeyValueToObjLit_1179(objNode, keyNode, valueNode, NodeUtil_GetEnclosingScript(callNode))
		}
		callNode.ReplaceWith(objNode)
		this_.compiler.ReportChangeToEnclosingScope(objNode)
	}
}

// Java method: ClosureOptimizePrimitives.canOptimizeObjectCreateSet(com.google.javascript.rhino.Node firstParam)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_canOptimizeObjectCreateSet_1178(firstParam pkg_rhino.Node) bool {
	_ = firstParam
	if ((firstParam != 0) && (firstParam.GetNext() == 0)) && !(firstParam.IsNumber() || firstParam.IsStringLit()) {
		return false
	}
	var curParam pkg_rhino.Node = firstParam
	_ = curParam
	var keys j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = keys
	for curParam != 0 {
		if !(this_.priv_ClosureOptimizePrimitives_isOptimizableKey_1180(curParam)) {
			return false
		}
		if curParam.IsStringLit() || curParam.IsNumber() {
			var key string
			_ = key
			if curParam.IsStringLit() {
				key = curParam.GetString()
			} else {
				key = pkg_dtoa.DToA_NumberToString(curParam.GetDouble())
			}
			if !(keys.Add(key)) {
				this_.compiler.Report_1(JSError_Make_3(firstParam.GetPrevious(), StaticGet_ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER(), []string{key}))
				return false
			}
		}
		curParam = curParam.GetNext()
	}
	return true
}

// Java method: ClosureOptimizePrimitives.addKeyValueToObjLit(com.google.javascript.rhino.Node objNode, com.google.javascript.rhino.Node keyNode, com.google.javascript.rhino.Node valueNode, com.google.javascript.rhino.Node scriptNode)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_addKeyValueToObjLit_1179(objNode pkg_rhino.Node, keyNode pkg_rhino.Node, valueNode pkg_rhino.Node, scriptNode pkg_rhino.Node) {
	_ = objNode
	_ = keyNode
	_ = valueNode
	_ = scriptNode
	if keyNode.IsNumber() || keyNode.IsStringLit() {
		if keyNode.IsNumber() {
			keyNode = pkg_rhino.IR_String(pkg_dtoa.DToA_NumberToString(keyNode.GetDouble())).Srcref(keyNode)
		}
		keyNode.SetIsParenthesized(false)
		keyNode.SetToken(pkg_rhino.StaticGet_Token_STRING_KEY())
		keyNode.SetQuotedStringKey()
		objNode.AddChildToBack(pkg_rhino.IR_Propdef(keyNode, valueNode))
	} else {
		objNode.AddChildToBack(pkg_rhino.IR_ComputedProp(keyNode, valueNode).Srcref(keyNode))
		NodeUtil_AddFeatureToScript(scriptNode, pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), this_.compiler)
	}
}

// Java method: ClosureOptimizePrimitives.isOptimizableKey(com.google.javascript.rhino.Node curParam)
func (this_ *ClosureOptimizePrimitives) priv_ClosureOptimizePrimitives_isOptimizableKey_1180(curParam pkg_rhino.Node) bool {
	_ = curParam
	if this_.canUseEs6Syntax {
		return !(NodeUtil_IsStatement(curParam))
	} else {
		return (curParam.IsStringLit() || curParam.IsNumber())
	}
}

// Java class: com.google.javascript.jscomp.ClosureOptimizePrimitives.FindPrimitives
type ClosureOptimizePrimitives_FindPrimitives struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ClosureOptimizePrimitives
}

func (this_ *ClosureOptimizePrimitives_FindPrimitives) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureOptimizePrimitives_FindPrimitives) IsJavaNil_() bool { return this_ == nil }
func (*ClosureOptimizePrimitives_FindPrimitives) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureOptimizePrimitives$FindPrimitives", "FindPrimitives"
}
func (this_ *ClosureOptimizePrimitives_FindPrimitives) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ClosureOptimizePrimitives_FindPrimitives) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureOptimizePrimitives.FindPrimitives", "com.google.javascript.jscomp.ClosureOptimizePrimitives$FindPrimitives", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureOptimizePrimitives_FindPrimitives](v_)
	})
}
func (this_ *ClosureOptimizePrimitives_FindPrimitives) SetOuter_(o *ClosureOptimizePrimitives) {
	this_.outer_ = o
}
func (this_ *ClosureOptimizePrimitives_FindPrimitives) GetOuter_() *ClosureOptimizePrimitives {
	return this_.outer_
}
func (*ClosureOptimizePrimitives_FindPrimitives) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ClosureOptimizePrimitives_FindPrimitives) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ClosureOptimizePrimitives_FindPrimitives) initFields_ClosureOptimizePrimitives_FindPrimitives() {
	_ = this_
}

// Java constructor: FindPrimitives()
func NewClosureOptimizePrimitives_FindPrimitives(outer_ *ClosureOptimizePrimitives) *ClosureOptimizePrimitives_FindPrimitives {
	this_ := &ClosureOptimizePrimitives_FindPrimitives{}
	this_.outer_ = outer_
	InitClosureOptimizePrimitives_FindPrimitives(this_)
	return this_
}

func InitClosureOptimizePrimitives_FindPrimitives(this_ *ClosureOptimizePrimitives_FindPrimitives) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ClosureOptimizePrimitives_FindPrimitives()
}

// Java method: FindPrimitives.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ClosureOptimizePrimitives_FindPrimitives) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsCall() {
		var fn pkg_rhino.Node = n.GetFirstChild()
		_ = fn
		if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(fn) {
			this_.outer_.priv_ClosureOptimizePrimitives_processRenamePropertyCall_1175(n)
		} else if (fn.MatchesName_1_string("goog$object$create") || fn.MatchesName_1_string("module$contents$goog$object_create")) || StaticGet_ClosureOptimizePrimitives_GOOG_OBJECT_CREATE().Matches(fn) {
			this_.outer_.priv_ClosureOptimizePrimitives_processObjectCreateCall_1174(n)
		} else if (fn.MatchesName_1_string("goog$object$createSet") || fn.MatchesName_1_string("module$contents$goog$object_createSet")) || StaticGet_ClosureOptimizePrimitives_GOOG_OBJECT_CREATESET().Matches(fn) {
			this_.outer_.priv_ClosureOptimizePrimitives_processObjectCreateSetCall_1177(n)
		}
	}
}
func (*ClosureOptimizePrimitives_FindPrimitives) AlwaysShouldTraverse_()      {}
func (*ClosureOptimizePrimitives_FindPrimitives) CanReleaseSyntacticScopes_() {}
