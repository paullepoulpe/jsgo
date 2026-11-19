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
// Port of com/google/javascript/jscomp/lint/CheckProvidesSorted.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_check_provides_sorted_1 j2gort.Ext_com_google_common_base_Function

func init() {
	lambdaSingleton_check_provides_sorted_1 = j2gort.NewLambda_Ext_com_google_common_base_Function(func(mrp0_ any) any {
		return j2gort.BoxString(CheckProvidesSorted_FormatProvide(j2gort.CastString(mrp0_)))
	})
}

// Java class: com.google.javascript.jscomp.lint.CheckProvidesSorted
type CheckProvidesSorted struct {
	originalProvides j2gort.Ext_java_util_List `gson:"originalProvides"`
	replacement      string                    `gson:"replacement"`
	firstNode        pkg_rhino.Node            `gson:"firstNode"`
	lastNode         pkg_rhino.Node            `gson:"lastNode"`
	idHash_          atomic.Int32
	mode             CheckProvidesSorted_Mode `gson:"mode"`
	finished         bool                     `gson:"finished"`
	needsFix         bool                     `gson:"needsFix"`
}

func (this_ *CheckProvidesSorted) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckProvidesSorted) IsJavaNil_() bool { return this_ == nil }
func (*CheckProvidesSorted) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckProvidesSorted", "CheckProvidesSorted"
}
func (this_ *CheckProvidesSorted) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckProvidesSorted) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckProvidesSorted", "com.google.javascript.jscomp.lint.CheckProvidesSorted", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckProvidesSorted](v_)
	})
}
func (*CheckProvidesSorted) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckProvidesSorted_PROVIDES_NOT_SORTED *DiagnosticType

func StaticGet_CheckProvidesSorted_PROVIDES_NOT_SORTED() *DiagnosticType {
	clinit_CheckProvidesSorted()
	return static_CheckProvidesSorted_PROVIDES_NOT_SORTED
}
func StaticSet_CheckProvidesSorted_PROVIDES_NOT_SORTED(v *DiagnosticType) {
	clinit_CheckProvidesSorted()
	static_CheckProvidesSorted_PROVIDES_NOT_SORTED = v
}

var clinitDone_CheckProvidesSorted uint32
var clinitState_CheckProvidesSorted j2gort.ClinitState

func clinit_CheckProvidesSorted() {
	if atomic.LoadUint32(&clinitDone_CheckProvidesSorted) != 0 {
		return
	}
	clinitSlow_CheckProvidesSorted()
}

func clinitSlow_CheckProvidesSorted() {
	if !j2gort.EnterClinit(&clinitState_CheckProvidesSorted, &clinitDone_CheckProvidesSorted) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckProvidesSorted, &clinitDone_CheckProvidesSorted)
	static_CheckProvidesSorted_PROVIDES_NOT_SORTED = DiagnosticType_Warning("JSC_PROVIDES_NOT_SORTED", "goog.provide() statements are not sorted. The correct order is:\n\n{0}\n")
}

func (this_ *CheckProvidesSorted) initFields_CheckProvidesSorted() {
	_ = this_
	this_.originalProvides = j2gort.NewExt_java_util_ArrayList()

	this_.replacement = j2gort.NullString

}

// Java constructor: CheckProvidesSorted(com.google.javascript.jscomp.lint.CheckProvidesSorted.Mode mode)
func NewCheckProvidesSorted(mode CheckProvidesSorted_Mode) *CheckProvidesSorted {
	clinit_CheckProvidesSorted()
	this_ := &CheckProvidesSorted{}
	InitCheckProvidesSorted(this_, mode)
	return this_
}

func InitCheckProvidesSorted(this_ *CheckProvidesSorted, mode CheckProvidesSorted_Mode) {
	_ = mode
	this_.initFields_CheckProvidesSorted()
	this_.mode = mode
}

// Java method: CheckProvidesSorted.getFirstNode()
func (this_ *CheckProvidesSorted) GetFirstNode() pkg_rhino.Node {
	return this_.firstNode
}

// Java method: CheckProvidesSorted.getLastNode()
func (this_ *CheckProvidesSorted) GetLastNode() pkg_rhino.Node {
	return this_.lastNode
}

// Java method: CheckProvidesSorted.getReplacement()
func (this_ *CheckProvidesSorted) GetReplacement() string {
	return this_.replacement
}

// Java method: CheckProvidesSorted.needsFix()
func (this_ *CheckProvidesSorted) NeedsFix() bool {
	return this_.needsFix
}

// Java method: CheckProvidesSorted.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckProvidesSorted) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(this_.finished) && ((((parent == 0) || parent.IsRoot()) || parent.IsScript()) || parent.IsModuleBody()))
}

// Java method: CheckProvidesSorted.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckProvidesSorted) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.priv_CheckProvidesSorted_checkCanonical_6571(t)
		return
	}
	if NodeUtil_IsGoogProvideCall(n) && CheckProvidesSorted_AreProvideArgumentsValid(n.GetFirstChild()) {
		this_.originalProvides.Add_1(j2gort.BoxString(CheckProvidesSorted_GetNamespace(n)))
		if this_.firstNode == 0 {
			this_.firstNode = func() pkg_rhino.Node { var asgn_1 pkg_rhino.Node = n; this_.lastNode = asgn_1; return asgn_1 }()
		} else {
			this_.lastNode = n
		}
	} else if !(this_.originalProvides.IsEmpty()) {
		this_.finished = true
	}
}

// Java static method: CheckProvidesSorted.areProvideArgumentsValid(com.google.javascript.rhino.Node call)
func CheckProvidesSorted_AreProvideArgumentsValid(call pkg_rhino.Node) bool {
	clinit_CheckProvidesSorted()
	_ = call
	return (call.HasTwoChildren() && call.GetSecondChild().IsStringLit())
}

// Java static method: CheckProvidesSorted.getNamespace(com.google.javascript.rhino.Node n)
func CheckProvidesSorted_GetNamespace(n pkg_rhino.Node) string {
	clinit_CheckProvidesSorted()
	_ = n
	return n.GetFirstChild().GetSecondChild().GetString()
}

// Java static method: CheckProvidesSorted.formatProvide(java.lang.String namespace)
func CheckProvidesSorted_FormatProvide(namespace string) string {
	clinit_CheckProvidesSorted()
	_ = namespace
	var sb j2gort.LocalStringBuilder

	sb.AppendStr("goog.provide('")
	sb.AppendStr(namespace)
	sb.AppendStr("');")
	return sb.String()
}

// Java method: CheckProvidesSorted.checkCanonical(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckProvidesSorted) priv_CheckProvidesSorted_checkCanonical_6571(t *NodeTraversal) {
	_ = t
	var canonicalProvides j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.originalProvides)))
	_ = canonicalProvides
	if !(this_.originalProvides.Equals(canonicalProvides)) {
		this_.needsFix = true
		this_.replacement = j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable("\n", j2gort.Ext_com_google_common_collect_Iterables_Transform(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](canonicalProvides), lambdaSingleton_check_provides_sorted_1))
		if this_.mode == StaticGet_CheckProvidesSorted_Mode_COLLECT_AND_REPORT() {
			t.Report_3(this_.firstNode, StaticGet_CheckProvidesSorted_PROVIDES_NOT_SORTED(), []string{this_.replacement})
		}
	}
}

// Java enum: com.google.javascript.jscomp.lint.CheckProvidesSorted.Mode
type CheckProvidesSorted_Mode uint8

const (
	enum_CheckProvidesSorted_Mode_COLLECT_ONLY       CheckProvidesSorted_Mode = 1
	enum_CheckProvidesSorted_Mode_COLLECT_AND_REPORT CheckProvidesSorted_Mode = 2
)

var enumBaseTable_CheckProvidesSorted_Mode = [...]j2gort.EnumBase{
	{Name_: "COLLECT_ONLY", Ordinal_: 0},
	{Name_: "COLLECT_AND_REPORT", Ordinal_: 1},
}
var enumValues_CheckProvidesSorted_Mode = []CheckProvidesSorted_Mode{1, 2}

func StaticGet_CheckProvidesSorted_Mode_COLLECT_ONLY() CheckProvidesSorted_Mode {
	return enum_CheckProvidesSorted_Mode_COLLECT_ONLY
}
func StaticGet_CheckProvidesSorted_Mode_COLLECT_AND_REPORT() CheckProvidesSorted_Mode {
	return enum_CheckProvidesSorted_Mode_COLLECT_AND_REPORT
}

func CheckProvidesSorted_Mode_Values() []CheckProvidesSorted_Mode {
	return j2gort.CloneSlice(enumValues_CheckProvidesSorted_Mode)
}

func CheckProvidesSorted_Mode_ValueOf(name string) CheckProvidesSorted_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckProvidesSorted_Mode {
		if enumBaseTable_CheckProvidesSorted_Mode[i_].Name_ == name {
			return CheckProvidesSorted_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.lint.CheckProvidesSorted.Mode." + name))
}

func (this_ CheckProvidesSorted_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckProvidesSorted_Mode) HasIdentityEquals_()                                 {}
func (CheckProvidesSorted_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckProvidesSorted_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckProvidesSorted_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckProvidesSorted_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckProvidesSorted$Mode", "Mode"
}
func (this_ CheckProvidesSorted_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckProvidesSorted_Mode[this_-1]
}
func (this_ CheckProvidesSorted_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckProvidesSorted_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckProvidesSorted_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckProvidesSorted_Mode[this_-1].Name_
}
func (this_ CheckProvidesSorted_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CheckProvidesSorted_Mode[this_-1].Name_
}
func (this_ CheckProvidesSorted_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckProvidesSorted_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckProvidesSorted_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckProvidesSorted_Mode[this_-1].HashCode()
}
func (this_ CheckProvidesSorted_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckProvidesSorted_Mode[this_-1].CompareTo(other_)
}
func (this_ CheckProvidesSorted_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckProvidesSorted_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckProvidesSorted_Mode(1), CheckProvidesSorted_Mode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckProvidesSorted.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckProvidesSorted$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckProvidesSorted.Mode", "com.google.javascript.jscomp.lint.CheckProvidesSorted$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[CheckProvidesSorted_Mode](v_)
	})
}
func (*CheckProvidesSorted) CanReleaseSyntacticScopes_() {}
