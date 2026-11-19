/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/modules/ClosureRequireProcessor.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.ClosureRequireProcessor
type ClosureRequireProcessor struct {
	nameDeclaration pkg_rhino.Node `gson:"nameDeclaration"`
	idHash_         atomic.Int32
	requireKind     Binding_CreatedBy `gson:"requireKind"`
}

func (this_ *ClosureRequireProcessor) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRequireProcessor) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRequireProcessor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ClosureRequireProcessor", "ClosureRequireProcessor"
}
func (this_ *ClosureRequireProcessor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureRequireProcessor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ClosureRequireProcessor", "com.google.javascript.jscomp.modules.ClosureRequireProcessor", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRequireProcessor](v_)
	})
}

var static_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ClosureRequireProcessor()
	return static_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS
}
func StaticSet_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_ClosureRequireProcessor()
	static_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS = v
}

var clinitDone_ClosureRequireProcessor uint32
var clinitState_ClosureRequireProcessor j2gort.ClinitState

func clinit_ClosureRequireProcessor() {
	if atomic.LoadUint32(&clinitDone_ClosureRequireProcessor) != 0 {
		return
	}
	clinitSlow_ClosureRequireProcessor()
}

func clinitSlow_ClosureRequireProcessor() {
	if !j2gort.EnterClinit(&clinitState_ClosureRequireProcessor, &clinitDone_ClosureRequireProcessor) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureRequireProcessor, &clinitDone_ClosureRequireProcessor)
	static_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_6("require", j2gort.BoxNullableVal(StaticGet_Binding_CreatedBy_GOOG_REQUIRE()), "requireType", j2gort.BoxNullableVal(StaticGet_Binding_CreatedBy_GOOG_REQUIRE_TYPE()), "forwardDeclare", j2gort.BoxNullableVal(StaticGet_Binding_CreatedBy_GOOG_FORWARD_DECLARE()))
}

func (this_ *ClosureRequireProcessor) initFields_ClosureRequireProcessor() {
	_ = this_
}

// Java constructor: ClosureRequireProcessor(com.google.javascript.rhino.Node nameDeclaration, com.google.javascript.jscomp.modules.Binding.CreatedBy requireKind)
func NewClosureRequireProcessor(nameDeclaration pkg_rhino.Node, requireKind Binding_CreatedBy) *ClosureRequireProcessor {
	clinit_ClosureRequireProcessor()
	this_ := &ClosureRequireProcessor{}
	InitClosureRequireProcessor(this_, nameDeclaration, requireKind)
	return this_
}

func InitClosureRequireProcessor(this_ *ClosureRequireProcessor, nameDeclaration pkg_rhino.Node, requireKind Binding_CreatedBy) {
	_ = nameDeclaration
	_ = requireKind
	this_.initFields_ClosureRequireProcessor()
	if !(NodeUtil_IsNameDeclaration(nameDeclaration)) {
		j2gort.CheckArgument(false)
	}
	this_.nameDeclaration = nameDeclaration
	this_.requireKind = requireKind
}

// Java static method: ClosureRequireProcessor.getAllRequires(com.google.javascript.rhino.Node nameDeclaration)
func ClosureRequireProcessor_GetAllRequires(nameDeclaration pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_ClosureRequireProcessor()
	_ = nameDeclaration
	var rhs pkg_rhino.Node
	_ = rhs
	if nameDeclaration.GetFirstChild().IsDestructuringLhs() {
		rhs = nameDeclaration.GetFirstChild().GetSecondChild()
	} else {
		rhs = nameDeclaration.GetFirstFirstChild()
	}
	var requireKind Binding_CreatedBy = ClosureRequireProcessor_GetModuleDependencyTypeFromRhs(rhs)
	_ = requireKind
	if requireKind == 0 {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	return NewClosureRequireProcessor(nameDeclaration, requireKind).priv_ClosureRequireProcessor_getAllRequiresInDeclaration_6654()
}

// Java static method: ClosureRequireProcessor.getModuleDependencyTypeFromRhs(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value)
func ClosureRequireProcessor_GetModuleDependencyTypeFromRhs(value pkg_rhino.Node) Binding_CreatedBy {
	clinit_ClosureRequireProcessor()
	_ = value
	if (((value == 0) || !(value.IsCall())) || !(value.HasTwoChildren())) || !(value.GetSecondChild().IsStringLit()) {
		return Binding_CreatedBy(0)
	}
	var callee pkg_rhino.Node = value.GetFirstChild()
	_ = callee
	if !(callee.IsGetProp()) {
		return Binding_CreatedBy(0)
	}
	var owner pkg_rhino.Node = callee.GetFirstChild()
	_ = owner
	if !(owner.IsName()) || !(j2gort.CheckNotNullString(owner.GetString()) == "goog") {
		return Binding_CreatedBy(0)
	}
	return j2gort.Cast[Binding_CreatedBy](j2gort.MapGetStr(StaticGet_ClosureRequireProcessor_GOOG_DEPENDENCY_CALLS(), callee.GetString()))
}

// Java method: ClosureRequireProcessor.getAllRequiresInDeclaration()
func (this_ *ClosureRequireProcessor) priv_ClosureRequireProcessor_getAllRequiresInDeclaration_6654() j2gort.Ext_com_google_common_collect_ImmutableList {
	var rhs pkg_rhino.Node
	_ = rhs
	if this_.nameDeclaration.GetFirstChild().IsDestructuringLhs() {
		rhs = this_.nameDeclaration.GetFirstChild().GetSecondChild()
	} else {
		rhs = this_.nameDeclaration.GetFirstFirstChild()
	}
	var namespace string = rhs.GetSecondChild().GetString()
	_ = namespace
	if this_.nameDeclaration.GetFirstChild().IsName() {
		var lhs pkg_rhino.Node = this_.nameDeclaration.GetFirstChild()
		_ = lhs
		return j2gort.ImmListOfPtr(ClosureRequireProcessor_Require_Create(lhs.GetString(), Import_S_Builder().ModuleRequest(namespace).LocalName(lhs.GetString()).ImportName("*exports*").ImportNode(this_.nameDeclaration).NameNode(lhs).Build(), this_.requireKind))
	} else {
		var objectPattern pkg_rhino.Node = this_.nameDeclaration.GetFirstFirstChild()
		_ = objectPattern
		if !(objectPattern.IsObjectPattern()) {
			return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
		}
		return this_.priv_ClosureRequireProcessor_getAllRequiresFromDestructuring_6655(objectPattern, namespace)
	}
}

// Java method: ClosureRequireProcessor.getAllRequiresFromDestructuring(com.google.javascript.rhino.Node objectPattern, java.lang.String namespace)
func (this_ *ClosureRequireProcessor) priv_ClosureRequireProcessor_getAllRequiresFromDestructuring_6655(objectPattern pkg_rhino.Node, namespace string) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = objectPattern
	_ = namespace
	var requireBuilder []*ClosureRequireProcessor_Require = []*ClosureRequireProcessor_Require{}
	_ = requireBuilder
	{
		var key pkg_rhino.Node = objectPattern.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if !(key.IsStringKey()) {
				continue
			}
			var lhs pkg_rhino.Node = key.GetOnlyChild()
			_ = lhs
			if !(lhs.IsName()) {
				continue
			}
			requireBuilder = append(requireBuilder, j2gort.CheckNotNullElem(ClosureRequireProcessor_Require_Create(lhs.GetString(), Import_S_Builder().ModuleRequest(namespace).LocalName(lhs.GetString()).ImportName(key.GetString()).ImportNode(this_.nameDeclaration).NameNode(lhs).Build(), this_.requireKind)))
		}
	}
	return j2gort.ImmListFromTypedSlice(requireBuilder)
}

// Java record: com.google.javascript.jscomp.modules.ClosureRequireProcessor.Require
type ClosureRequireProcessor_Require struct {
	localName    string            `gson:"localName"`
	importRecord *Import           `gson:"importRecord"`
	createdBy    Binding_CreatedBy `gson:"createdBy"`
}

func (this_ *ClosureRequireProcessor_Require) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureRequireProcessor_Require) IsJavaNil_() bool { return this_ == nil }
func (*ClosureRequireProcessor_Require) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.ClosureRequireProcessor$Require", "Require"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ClosureRequireProcessor.Require", "com.google.javascript.jscomp.modules.ClosureRequireProcessor$Require", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureRequireProcessor_Require](v_)
	})
}
func (*ClosureRequireProcessor_Require) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ClosureRequireProcessor_Require) initFields_ClosureRequireProcessor_Require() {
	_ = this_
}

// Java constructor: Require(java.lang.String localName, com.google.javascript.jscomp.modules.Import importRecord, com.google.javascript.jscomp.modules.Binding.CreatedBy createdBy)
func NewClosureRequireProcessor_Require(localName string, importRecord *Import, createdBy Binding_CreatedBy) *ClosureRequireProcessor_Require {
	this_ := &ClosureRequireProcessor_Require{}
	InitClosureRequireProcessor_Require(this_, localName, importRecord, createdBy)
	return this_
}

func InitClosureRequireProcessor_Require(this_ *ClosureRequireProcessor_Require, localName string, importRecord *Import, createdBy Binding_CreatedBy) {
	_ = localName
	_ = importRecord
	_ = createdBy
	this_.initFields_ClosureRequireProcessor_Require()
	j2gort.CheckNotNullString(localName)
	j2gort.CheckNotNullPtr(importRecord)
	j2gort.CheckNotNullEnum(createdBy)
	if !(createdBy.IsClosureImport()) {
		j2gort.CheckArgument(false)
	}
	this_.localName = localName
	this_.importRecord = importRecord
	this_.createdBy = createdBy
}

// Java static method: Require.create(java.lang.String localName, com.google.javascript.jscomp.modules.Import importRecord, com.google.javascript.jscomp.modules.Binding.CreatedBy createdBy)
func ClosureRequireProcessor_Require_Create(localName string, importRecord *Import, createdBy Binding_CreatedBy) *ClosureRequireProcessor_Require {
	_ = localName
	_ = importRecord
	_ = createdBy
	return NewClosureRequireProcessor_Require(localName, importRecord, createdBy)
}

func (this_ *ClosureRequireProcessor_Require) LocalName() string            { return this_.localName }
func (this_ *ClosureRequireProcessor_Require) ImportRecord() *Import        { return this_.importRecord }
func (this_ *ClosureRequireProcessor_Require) CreatedBy() Binding_CreatedBy { return this_.createdBy }
func (this_ *ClosureRequireProcessor_Require) Equals(other_ any) bool {
	o_, ok_ := other_.(*ClosureRequireProcessor_Require)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.localName != o_.localName {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.importRecord), j2gort.PtrToAny(o_.importRecord)) {
		return false
	}
	if this_.createdBy != o_.createdBy {
		return false
	}
	return true
}
func (this_ *ClosureRequireProcessor_Require) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_1 := this_.localName
		h_ = 31 * h_
		if hs_1 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_1)
		}
	}
	{
		hp_2 := this_.importRecord
		h_ = 31 * h_
		if hp_2 != nil {
			h_ += hp_2.HashCode()
		}
	}
	h_ = 31*h_ + (this_.createdBy).HashCode()
	return h_
}
func (this_ *ClosureRequireProcessor_Require) String() string {
	_ = fmt.Sprint
	return "Require[" + "localName=" + j2gort.StrVal(j2gort.BoxString(this_.localName)) + ", " + "importRecord=" + j2gort.StrVal(j2gort.PtrToAny(this_.importRecord)) + ", " + "createdBy=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.createdBy)) + "]"
}
