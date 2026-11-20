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
// Port of com/google/javascript/jscomp/lint/CheckRequiresSorted.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_check_requires_sorted_1 j2gort.Ext_com_google_common_base_Function
var lambdaSingleton_check_requires_sorted_2 j2gort.Ext_java_util_function_BinaryOperator
var lambdaSingleton_check_requires_sorted_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_requires_sorted_4 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_check_requires_sorted_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_requires_sorted_6 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_check_requires_sorted_7 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_requires_sorted_8 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_requires_sorted_9 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_check_requires_sorted_10 j2gort.Ext_java_util_function_Function
var lambdaSingleton_check_requires_sorted_11 j2gort.Ext_java_util_function_Predicate

func init() {
	lambdaSingleton_check_requires_sorted_1 = j2gort.NewLambda_Ext_com_google_common_base_Function(func(mrp0_ any) any {
		return j2gort.BoxString(j2gort.CastPtr[CheckRequiresSorted_ImportStatement](mrp0_).Format())
	})
	lambdaSingleton_check_requires_sorted_2 = j2gort.NewLambda_Ext_java_util_function_BinaryOperator(func(mrp0_ any, mrp1_ any) any {
		return j2gort.BoxNullableVal(CheckRequiresSorted_ImportPrimitive_Stronger(j2gort.Cast[CheckRequiresSorted_ImportPrimitive](mrp0_), j2gort.Cast[CheckRequiresSorted_ImportPrimitive](mrp1_)))
	})
	lambdaSingleton_check_requires_sorted_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxNullableVal(j2gort.CastPtr[CheckRequiresSorted_ImportStatement](mrp0_).Primitive())
	})
	lambdaSingleton_check_requires_sorted_4 = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(mrp0_ any) bool { return j2gort.CastPtr[CheckRequiresSorted_ImportStatement](mrp0_).IsAliasing() })
	lambdaSingleton_check_requires_sorted_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		i := j2gort.CastPtr[CheckRequiresSorted_ImportStatement](lp_5)
		_ = i
		return i.Nodes().Stream()
	})
	lambdaSingleton_check_requires_sorted_6 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(mrp0_ any) bool {
		return j2gort.CastPtr[CheckRequiresSorted_ImportStatement](mrp0_).IsDestructuring()
	})
	lambdaSingleton_check_requires_sorted_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(j2gort.CastPtr[CheckRequiresSorted_DestructuringBinding](mrp0_).CanonicalizeShorthandProperties())
	})
	lambdaSingleton_check_requires_sorted_8 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_6 any) any {
		i := j2gort.CastPtr[CheckRequiresSorted_ImportStatement](lp_6)
		_ = i
		return i.Destructures().Stream()
	})
	lambdaSingleton_check_requires_sorted_9 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(mrp0_ any) bool {
		return j2gort.CastPtr[CheckRequiresSorted_ImportStatement](mrp0_).IsDestructuring()
	})
	lambdaSingleton_check_requires_sorted_10 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
		i := j2gort.CastPtr[CheckRequiresSorted_ImportStatement](lp_7)
		_ = i
		return i.Nodes().Stream()
	})
	lambdaSingleton_check_requires_sorted_11 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(mrp0_ any) bool { return j2gort.CastPtr[CheckRequiresSorted_ImportStatement](mrp0_).IsStandalone() })
}

// Java class: com.google.javascript.jscomp.lint.CheckRequiresSorted
type CheckRequiresSorted struct {
	importsByNamespace j2gort.Ext_com_google_common_collect_Multimap `gson:"importsByNamespace"`
	originalImports    j2gort.Ext_java_util_List                     `gson:"originalImports"`
	replacement        string                                        `gson:"replacement"`
	firstNode          pkg_rhino.Node                                `gson:"firstNode"`
	lastNode           pkg_rhino.Node                                `gson:"lastNode"`
	idHash_            atomic.Int32
	mode               CheckRequiresSorted_Mode `gson:"mode"`
	finished           bool                     `gson:"finished"`
	needsFix           bool                     `gson:"needsFix"`
}

func (this_ *CheckRequiresSorted) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckRequiresSorted) IsJavaNil_() bool { return this_ == nil }
func (*CheckRequiresSorted) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckRequiresSorted", "CheckRequiresSorted"
}
func (this_ *CheckRequiresSorted) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckRequiresSorted) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckRequiresSorted", "com.google.javascript.jscomp.lint.CheckRequiresSorted", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckRequiresSorted](v_)
	})
}
func (*CheckRequiresSorted) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckRequiresSorted_REQUIRES_NOT_SORTED *DiagnosticType

func StaticGet_CheckRequiresSorted_REQUIRES_NOT_SORTED() *DiagnosticType {
	clinit_CheckRequiresSorted()
	return static_CheckRequiresSorted_REQUIRES_NOT_SORTED
}
func StaticSet_CheckRequiresSorted_REQUIRES_NOT_SORTED(v *DiagnosticType) {
	clinit_CheckRequiresSorted()
	static_CheckRequiresSorted_REQUIRES_NOT_SORTED = v
}

var static_CheckRequiresSorted_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_CheckRequiresSorted_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_CheckRequiresSorted()
	return static_CheckRequiresSorted_GOOG_REQUIRE
}
func StaticSet_CheckRequiresSorted_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_CheckRequiresSorted()
	static_CheckRequiresSorted_GOOG_REQUIRE = v
}

var static_CheckRequiresSorted_GOOG_REQUIRETYPE pkg_rhino.QualifiedName

func StaticGet_CheckRequiresSorted_GOOG_REQUIRETYPE() pkg_rhino.QualifiedName {
	clinit_CheckRequiresSorted()
	return static_CheckRequiresSorted_GOOG_REQUIRETYPE
}
func StaticSet_CheckRequiresSorted_GOOG_REQUIRETYPE(v pkg_rhino.QualifiedName) {
	clinit_CheckRequiresSorted()
	static_CheckRequiresSorted_GOOG_REQUIRETYPE = v
}

var static_CheckRequiresSorted_GOOG_FORWARDDECLARE pkg_rhino.QualifiedName

func StaticGet_CheckRequiresSorted_GOOG_FORWARDDECLARE() pkg_rhino.QualifiedName {
	clinit_CheckRequiresSorted()
	return static_CheckRequiresSorted_GOOG_FORWARDDECLARE
}
func StaticSet_CheckRequiresSorted_GOOG_FORWARDDECLARE(v pkg_rhino.QualifiedName) {
	clinit_CheckRequiresSorted()
	static_CheckRequiresSorted_GOOG_FORWARDDECLARE = v
}

var clinitDone_CheckRequiresSorted uint32
var clinitState_CheckRequiresSorted j2gort.ClinitState

func clinit_CheckRequiresSorted() {
	if atomic.LoadUint32(&clinitDone_CheckRequiresSorted) != 0 {
		return
	}
	clinitSlow_CheckRequiresSorted()
}

func clinitSlow_CheckRequiresSorted() {
	if !j2gort.EnterClinit(&clinitState_CheckRequiresSorted, &clinitDone_CheckRequiresSorted) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckRequiresSorted, &clinitDone_CheckRequiresSorted)
	static_CheckRequiresSorted_REQUIRES_NOT_SORTED = DiagnosticType_Warning("JSC_REQUIRES_NOT_SORTED", "goog.require() and goog.requireType() statements are not in recommended format. The correct order is:\n\n{0}\n")
	static_CheckRequiresSorted_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_CheckRequiresSorted_GOOG_REQUIRETYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_CheckRequiresSorted_GOOG_FORWARDDECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
}

func (this_ *CheckRequiresSorted) initFields_CheckRequiresSorted() {
	_ = this_
	this_.importsByNamespace = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
	this_.originalImports = j2gort.NewExt_java_util_ArrayList()

	this_.replacement = j2gort.NullString
}

// Java constructor: CheckRequiresSorted(com.google.javascript.jscomp.lint.CheckRequiresSorted.Mode mode)
func NewCheckRequiresSorted(mode CheckRequiresSorted_Mode) *CheckRequiresSorted {
	clinit_CheckRequiresSorted()
	this_ := &CheckRequiresSorted{}
	InitCheckRequiresSorted(this_, mode)
	return this_
}

func InitCheckRequiresSorted(this_ *CheckRequiresSorted, mode CheckRequiresSorted_Mode) {
	_ = mode
	this_.initFields_CheckRequiresSorted()
	this_.mode = mode
}

// Java method: CheckRequiresSorted.getFirstNode()
func (this_ *CheckRequiresSorted) GetFirstNode() pkg_rhino.Node {
	return this_.firstNode
}

// Java method: CheckRequiresSorted.getLastNode()
func (this_ *CheckRequiresSorted) GetLastNode() pkg_rhino.Node {
	return this_.lastNode
}

// Java method: CheckRequiresSorted.getReplacement()
func (this_ *CheckRequiresSorted) GetReplacement() string {
	return this_.replacement
}

// Java method: CheckRequiresSorted.needsFix()
func (this_ *CheckRequiresSorted) NeedsFix() bool {
	return this_.needsFix
}

// Java method: CheckRequiresSorted.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckRequiresSorted) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(this_.finished) && ((((parent == 0) || parent.IsRoot()) || parent.IsScript()) || parent.IsModuleBody()))
}

// Java method: CheckRequiresSorted.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckRequiresSorted) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.priv_CheckRequiresSorted_checkCanonical_6627(t)
		return
	}
	var callNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = callNode
	if n.IsExprResult() {
		callNode = n.GetFirstChild()
	} else if NodeUtil_IsNameDeclaration(n) {
		callNode = n.GetFirstChild().GetLastChild()
	}
	if (callNode != 0) && CheckRequiresSorted_IsValidImportCall(callNode) {
		var stmt *CheckRequiresSorted_ImportStatement = CheckRequiresSorted_ParseImport(callNode)
		_ = stmt
		this_.originalImports.Add_1(j2gort.PtrToAny(stmt))
		j2gort.MultimapPutStr(this_.importsByNamespace, stmt.Namespace(), j2gort.PtrToAny(stmt))
		if this_.firstNode == 0 {
			this_.firstNode = func() pkg_rhino.Node { var asgn_1 pkg_rhino.Node = n; this_.lastNode = asgn_1; return asgn_1 }()
		} else {
			this_.lastNode = n
		}
	} else if !(this_.importsByNamespace.IsEmpty()) {
		this_.finished = true
	}
}

// Java static method: CheckRequiresSorted.isValidImportCall(com.google.javascript.rhino.Node n)
func CheckRequiresSorted_IsValidImportCall(n pkg_rhino.Node) bool {
	clinit_CheckRequiresSorted()
	_ = n
	return (((n.IsCall() && n.HasTwoChildren()) && ((StaticGet_CheckRequiresSorted_GOOG_REQUIRE().Matches(n.GetFirstChild()) || StaticGet_CheckRequiresSorted_GOOG_REQUIRETYPE().Matches(n.GetFirstChild())) || StaticGet_CheckRequiresSorted_GOOG_FORWARDDECLARE().Matches(n.GetFirstChild()))) && n.GetSecondChild().IsStringLit())
}

// Java static method: CheckRequiresSorted.parseImport(com.google.javascript.rhino.Node callNode)
func CheckRequiresSorted_ParseImport(callNode pkg_rhino.Node) *CheckRequiresSorted_ImportStatement {
	clinit_CheckRequiresSorted()
	_ = callNode
	var primitive CheckRequiresSorted_ImportPrimitive = CheckRequiresSorted_ImportPrimitive_FromName(callNode.GetFirstChild().GetQualifiedName())
	_ = primitive
	var namespace string = callNode.GetSecondChild().GetString()
	_ = namespace
	var parent pkg_rhino.Node = callNode.GetParent()
	_ = parent
	if parent.IsExprResult() {
		return CheckRequiresSorted_ImportStatement_Of(j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(parent)), primitive, namespace, j2gort.NullString, nil)
	}
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if parent.IsName() {
		var alias string = parent.GetString()
		_ = alias
		return CheckRequiresSorted_ImportStatement_Of(j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(grandparent)), primitive, namespace, alias, nil)
	}
	var destructures []*CheckRequiresSorted_DestructuringBinding = []*CheckRequiresSorted_DestructuringBinding{}
	_ = destructures
	{
		var name pkg_rhino.Node = parent.GetFirstFirstChild()
		_ = name
		for ; name != 0; name = name.GetNext() {
			var exportedName string = name.GetString()
			_ = exportedName
			var localName string = name.GetFirstChild().GetString()
			_ = localName
			var isShorthandProperty bool = name.IsShorthandProperty()
			_ = isShorthandProperty
			destructures = append(destructures, j2gort.CheckNotNullElem(CheckRequiresSorted_DestructuringBinding_Of(exportedName, localName, isShorthandProperty)))
		}
	}
	return CheckRequiresSorted_ImportStatement_Of(j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(grandparent)), primitive, namespace, j2gort.NullString, j2gort.ImmListFromTypedSlice(destructures))
}

// Java method: CheckRequiresSorted.checkCanonical(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckRequiresSorted) priv_CheckRequiresSorted_checkCanonical_6627(t *NodeTraversal) {
	_ = t
	var canonicalImports j2gort.Ext_java_util_List = CheckRequiresSorted_CanonicalizeImports(this_.importsByNamespace)
	_ = canonicalImports
	if !(this_.originalImports.Equals(canonicalImports)) {
		this_.needsFix = true
		this_.replacement = j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable("\n", j2gort.Ext_com_google_common_collect_Iterables_Transform(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](canonicalImports), lambdaSingleton_check_requires_sorted_1))
		if this_.mode == StaticGet_CheckRequiresSorted_Mode_COLLECT_AND_REPORT() {
			t.Report_3(this_.firstNode, StaticGet_CheckRequiresSorted_REQUIRES_NOT_SORTED(), []string{this_.replacement})
		}
	}
}

// Java static method: CheckRequiresSorted.canonicalizeImports(com.google.common.collect.Multimap<java.lang.String,com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportStatement> importsByNamespace)
func CheckRequiresSorted_CanonicalizeImports(importsByNamespace j2gort.Ext_com_google_common_collect_Multimap) j2gort.Ext_java_util_List {
	clinit_CheckRequiresSorted()
	_ = importsByNamespace
	var canonicalImports []*CheckRequiresSorted_ImportStatement = []*CheckRequiresSorted_ImportStatement{}
	_ = canonicalImports
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](importsByNamespace.KeySet())))
		for cur_2.HasNext() {
			namespace := cur_2.NextStr()
			_ = namespace
			var allImports j2gort.Ext_java_util_Collection = j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetStr(importsByNamespace, namespace))
			_ = allImports
			var strongestPrimitive CheckRequiresSorted_ImportPrimitive = j2gort.Cast[CheckRequiresSorted_ImportPrimitive](allImports.Stream().Map(lambdaSingleton_check_requires_sorted_3).Reduce_2(j2gort.BoxNullableVal(StaticGet_CheckRequiresSorted_ImportPrimitive_WEAKEST()), lambdaSingleton_check_requires_sorted_2))
			_ = strongestPrimitive
			var hasAliasing bool = false
			_ = hasAliasing
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Iterables_Filter_2_Ext_java_lang_Iterable_Ext_com_google_common_base_Predicate(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](allImports), lambdaSingleton_check_requires_sorted_4)))
				for cur_3.HasNext() {
					rawElem_4 := cur_3.Next()
					stmt := j2gort.CastPtr[CheckRequiresSorted_ImportStatement](rawElem_4)
					_ = stmt
					canonicalImports = append(canonicalImports, stmt.Upgrade(strongestPrimitive))
					hasAliasing = true
				}
			}
			var hasDestructuring bool = false
			_ = hasDestructuring
			var destructuringNodes j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](allImports.Stream().Filter(lambdaSingleton_check_requires_sorted_6).FlatMap(lambdaSingleton_check_requires_sorted_5).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
			_ = destructuringNodes
			var destructures j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](allImports.Stream().Filter(lambdaSingleton_check_requires_sorted_9).FlatMap(lambdaSingleton_check_requires_sorted_8).Map(lambdaSingleton_check_requires_sorted_7).Sorted_0().Distinct().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
			_ = destructures
			if !(destructures.IsEmpty()) {
				canonicalImports = append(canonicalImports, CheckRequiresSorted_ImportStatement_Of(destructuringNodes, strongestPrimitive, namespace, j2gort.NullString, destructures))
				hasDestructuring = true
			}
			if !(hasAliasing) && !(hasDestructuring) {
				var standaloneNodes j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](allImports.Stream().Filter(lambdaSingleton_check_requires_sorted_11).FlatMap(lambdaSingleton_check_requires_sorted_10).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
				_ = standaloneNodes
				canonicalImports = append(canonicalImports, CheckRequiresSorted_ImportStatement_Of(standaloneNodes, strongestPrimitive, namespace, j2gort.NullString, nil))
			}
		}
	}
	j2gort.SliceSortNatural(canonicalImports)
	return j2gort.WrapTypedArrayList(canonicalImports)
}

// Java enum: com.google.javascript.jscomp.lint.CheckRequiresSorted.Mode
type CheckRequiresSorted_Mode uint8

const (
	enum_CheckRequiresSorted_Mode_COLLECT_ONLY       CheckRequiresSorted_Mode = 1
	enum_CheckRequiresSorted_Mode_COLLECT_AND_REPORT CheckRequiresSorted_Mode = 2
)

var enumBaseTable_CheckRequiresSorted_Mode = [...]j2gort.EnumBase{
	{Name_: "COLLECT_ONLY", Ordinal_: 0},
	{Name_: "COLLECT_AND_REPORT", Ordinal_: 1},
}
var enumValues_CheckRequiresSorted_Mode = []CheckRequiresSorted_Mode{1, 2}

func StaticGet_CheckRequiresSorted_Mode_COLLECT_ONLY() CheckRequiresSorted_Mode {
	return enum_CheckRequiresSorted_Mode_COLLECT_ONLY
}
func StaticGet_CheckRequiresSorted_Mode_COLLECT_AND_REPORT() CheckRequiresSorted_Mode {
	return enum_CheckRequiresSorted_Mode_COLLECT_AND_REPORT
}

func CheckRequiresSorted_Mode_Values() []CheckRequiresSorted_Mode {
	return j2gort.CloneSlice(enumValues_CheckRequiresSorted_Mode)
}

func CheckRequiresSorted_Mode_ValueOf(name string) CheckRequiresSorted_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckRequiresSorted_Mode {
		if enumBaseTable_CheckRequiresSorted_Mode[i_].Name_ == name {
			return CheckRequiresSorted_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.lint.CheckRequiresSorted.Mode." + name))
}

func (this_ CheckRequiresSorted_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckRequiresSorted_Mode) HasIdentityEquals_()                                 {}
func (CheckRequiresSorted_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckRequiresSorted_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckRequiresSorted_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckRequiresSorted_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckRequiresSorted$Mode", "Mode"
}
func (this_ CheckRequiresSorted_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckRequiresSorted_Mode[this_-1]
}
func (this_ CheckRequiresSorted_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckRequiresSorted_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckRequiresSorted_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckRequiresSorted_Mode[this_-1].Name_
}
func (this_ CheckRequiresSorted_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CheckRequiresSorted_Mode[this_-1].Name_
}
func (this_ CheckRequiresSorted_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckRequiresSorted_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckRequiresSorted_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckRequiresSorted_Mode[this_-1].HashCode()
}
func (this_ CheckRequiresSorted_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckRequiresSorted_Mode[this_-1].CompareTo(other_)
}
func (this_ CheckRequiresSorted_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckRequiresSorted_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckRequiresSorted_Mode(1), CheckRequiresSorted_Mode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckRequiresSorted.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckRequiresSorted$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckRequiresSorted.Mode", "com.google.javascript.jscomp.lint.CheckRequiresSorted$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[CheckRequiresSorted_Mode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive
type CheckRequiresSorted_ImportPrimitive uint8

const (
	enum_CheckRequiresSorted_ImportPrimitive_REQUIRE         CheckRequiresSorted_ImportPrimitive = 1
	enum_CheckRequiresSorted_ImportPrimitive_REQUIRE_TYPE    CheckRequiresSorted_ImportPrimitive = 2
	enum_CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE CheckRequiresSorted_ImportPrimitive = 3
)

var enumBaseTable_CheckRequiresSorted_ImportPrimitive = [...]j2gort.EnumBase{
	{Name_: "REQUIRE", Ordinal_: 0},
	{Name_: "REQUIRE_TYPE", Ordinal_: 1},
	{Name_: "FORWARD_DECLARE", Ordinal_: 2},
}
var enumValues_CheckRequiresSorted_ImportPrimitive = []CheckRequiresSorted_ImportPrimitive{1, 2, 3}

var enumField_CheckRequiresSorted_ImportPrimitive_name = [...]string{
	"goog.require",
	"goog.requireType",
	"goog.forwardDeclare",
}

func (this_ CheckRequiresSorted_ImportPrimitive) FieldGet_name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_CheckRequiresSorted_ImportPrimitive_name[this_-1]
}

func StaticGet_CheckRequiresSorted_ImportPrimitive_REQUIRE() CheckRequiresSorted_ImportPrimitive {
	return enum_CheckRequiresSorted_ImportPrimitive_REQUIRE
}
func StaticGet_CheckRequiresSorted_ImportPrimitive_REQUIRE_TYPE() CheckRequiresSorted_ImportPrimitive {
	return enum_CheckRequiresSorted_ImportPrimitive_REQUIRE_TYPE
}
func StaticGet_CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE() CheckRequiresSorted_ImportPrimitive {
	return enum_CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE
}

var static_CheckRequiresSorted_ImportPrimitive_WEAKEST CheckRequiresSorted_ImportPrimitive

func StaticGet_CheckRequiresSorted_ImportPrimitive_WEAKEST() CheckRequiresSorted_ImportPrimitive {
	clinit_CheckRequiresSorted_ImportPrimitive()
	return static_CheckRequiresSorted_ImportPrimitive_WEAKEST
}
func StaticSet_CheckRequiresSorted_ImportPrimitive_WEAKEST(v CheckRequiresSorted_ImportPrimitive) {
	clinit_CheckRequiresSorted_ImportPrimitive()
	static_CheckRequiresSorted_ImportPrimitive_WEAKEST = v
}

var clinitDone_CheckRequiresSorted_ImportPrimitive uint32
var clinitState_CheckRequiresSorted_ImportPrimitive j2gort.ClinitState

func clinit_CheckRequiresSorted_ImportPrimitive() {
	if atomic.LoadUint32(&clinitDone_CheckRequiresSorted_ImportPrimitive) != 0 {
		return
	}
	clinitSlow_CheckRequiresSorted_ImportPrimitive()
}

func clinitSlow_CheckRequiresSorted_ImportPrimitive() {
	if !j2gort.EnterClinit(&clinitState_CheckRequiresSorted_ImportPrimitive, &clinitDone_CheckRequiresSorted_ImportPrimitive) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckRequiresSorted_ImportPrimitive, &clinitDone_CheckRequiresSorted_ImportPrimitive)
	static_CheckRequiresSorted_ImportPrimitive_WEAKEST = StaticGet_CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE()
}

func CheckRequiresSorted_ImportPrimitive_Values() []CheckRequiresSorted_ImportPrimitive {
	clinit_CheckRequiresSorted_ImportPrimitive()
	return j2gort.CloneSlice(enumValues_CheckRequiresSorted_ImportPrimitive)
}

func CheckRequiresSorted_ImportPrimitive_ValueOf(name string) CheckRequiresSorted_ImportPrimitive {
	clinit_CheckRequiresSorted_ImportPrimitive()
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckRequiresSorted_ImportPrimitive {
		if enumBaseTable_CheckRequiresSorted_ImportPrimitive[i_].Name_ == name {
			return CheckRequiresSorted_ImportPrimitive(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive." + name))
}

// Java static method: ImportPrimitive.fromName(java.lang.String name)
func CheckRequiresSorted_ImportPrimitive_FromName(name string) CheckRequiresSorted_ImportPrimitive {
	clinit_CheckRequiresSorted_ImportPrimitive()
	_ = name
	for _, el_8 := range j2gort.CheckNotNullSlice(CheckRequiresSorted_ImportPrimitive_Values()) {
		primitive := el_8
		_ = primitive
		if j2gort.CheckNotNullString(primitive.FieldGet_name()) == name {
			return primitive
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid primitive name "+j2gort.StrValString(name))))
}

// Java static method: ImportPrimitive.stronger(com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive p1, com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive p2)
func CheckRequiresSorted_ImportPrimitive_Stronger(p1 CheckRequiresSorted_ImportPrimitive, p2 CheckRequiresSorted_ImportPrimitive) CheckRequiresSorted_ImportPrimitive {
	clinit_CheckRequiresSorted_ImportPrimitive()
	_ = p1
	_ = p2
	if (p1).Ordinal() < (p2).Ordinal() {
		return p1
	} else {
		return p2
	}
}

// Java method: ImportPrimitive.toString()
func (this_ CheckRequiresSorted_ImportPrimitive) String() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_name()
}

func (this_ CheckRequiresSorted_ImportPrimitive) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckRequiresSorted_ImportPrimitive) HasIdentityEquals_()                                 {}
func (CheckRequiresSorted_ImportPrimitive) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckRequiresSorted_ImportPrimitive) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckRequiresSorted_ImportPrimitive) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckRequiresSorted_ImportPrimitive) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckRequiresSorted$ImportPrimitive", "ImportPrimitive"
}
func (this_ CheckRequiresSorted_ImportPrimitive) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckRequiresSorted_ImportPrimitive[this_-1]
}
func (this_ CheckRequiresSorted_ImportPrimitive) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckRequiresSorted_ImportPrimitive) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckRequiresSorted_ImportPrimitive) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckRequiresSorted_ImportPrimitive[this_-1].Name_
}
func (this_ CheckRequiresSorted_ImportPrimitive) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckRequiresSorted_ImportPrimitive)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckRequiresSorted_ImportPrimitive) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckRequiresSorted_ImportPrimitive[this_-1].HashCode()
}
func (this_ CheckRequiresSorted_ImportPrimitive) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckRequiresSorted_ImportPrimitive[this_-1].CompareTo(other_)
}
func (this_ CheckRequiresSorted_ImportPrimitive) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckRequiresSorted_ImportPrimitive[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckRequiresSorted_ImportPrimitive(1), CheckRequiresSorted_ImportPrimitive(2), CheckRequiresSorted_ImportPrimitive(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckRequiresSorted$ImportPrimitive", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive", "com.google.javascript.jscomp.lint.CheckRequiresSorted$ImportPrimitive", func(v_ any) bool {
		return j2gort.InstanceOf[CheckRequiresSorted_ImportPrimitive](v_)
	})
}

// Java record: com.google.javascript.jscomp.lint.CheckRequiresSorted.DestructuringBinding
type CheckRequiresSorted_DestructuringBinding struct {
	exportedName        string `gson:"exportedName"`
	localName           string `gson:"localName"`
	isShorthandProperty bool   `gson:"isShorthandProperty"`
}

func (this_ *CheckRequiresSorted_DestructuringBinding) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckRequiresSorted_DestructuringBinding) IsJavaNil_() bool { return this_ == nil }
func (*CheckRequiresSorted_DestructuringBinding) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckRequiresSorted$DestructuringBinding", "DestructuringBinding"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckRequiresSorted.DestructuringBinding", "com.google.javascript.jscomp.lint.CheckRequiresSorted$DestructuringBinding", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckRequiresSorted_DestructuringBinding](v_)
	})
}
func (*CheckRequiresSorted_DestructuringBinding) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *CheckRequiresSorted_DestructuringBinding) initFields_CheckRequiresSorted_DestructuringBinding() {
	_ = this_
}

// Java constructor: DestructuringBinding(java.lang.String exportedName, java.lang.String localName, boolean isShorthandProperty)
func NewCheckRequiresSorted_DestructuringBinding(exportedName string, localName string, isShorthandProperty bool) *CheckRequiresSorted_DestructuringBinding {
	this_ := &CheckRequiresSorted_DestructuringBinding{}
	InitCheckRequiresSorted_DestructuringBinding(this_, exportedName, localName, isShorthandProperty)
	return this_
}

func InitCheckRequiresSorted_DestructuringBinding(this_ *CheckRequiresSorted_DestructuringBinding, exportedName string, localName string, isShorthandProperty bool) {
	_ = exportedName
	_ = localName
	_ = isShorthandProperty
	this_.initFields_CheckRequiresSorted_DestructuringBinding()
	j2gort.CheckNotNullString(exportedName)
	j2gort.CheckNotNullString(localName)
	if !(!(isShorthandProperty) || (j2gort.CheckNotNullString(exportedName) == localName)) {
		j2gort.CheckArgument(false)
	}
	this_.exportedName = exportedName
	this_.localName = localName
	this_.isShorthandProperty = isShorthandProperty
}

// Java static method: DestructuringBinding.of(java.lang.String exportedName, java.lang.String localName, boolean isShorthandProperty)
func CheckRequiresSorted_DestructuringBinding_Of(exportedName string, localName string, isShorthandProperty bool) *CheckRequiresSorted_DestructuringBinding {
	_ = exportedName
	_ = localName
	_ = isShorthandProperty
	return NewCheckRequiresSorted_DestructuringBinding(exportedName, localName, isShorthandProperty)
}

// Java method: DestructuringBinding.compareTo(com.google.javascript.jscomp.lint.CheckRequiresSorted.DestructuringBinding other)
func (this_ *CheckRequiresSorted_DestructuringBinding) CompareTo(arg_9 any) int32 {
	other := j2gort.CastPtr[CheckRequiresSorted_DestructuringBinding](arg_9)
	_ = other
	return j2gort.Ext_com_google_common_collect_ComparisonChain_Start().Compare_2_any_any(j2gort.BoxString(this_.ExportedName()), j2gort.BoxString(other.ExportedName())).Compare_2_any_any(j2gort.BoxString(this_.LocalName()), j2gort.BoxString(other.LocalName())).Result()
}

// Java method: DestructuringBinding.isCanonical()
func (this_ *CheckRequiresSorted_DestructuringBinding) priv_DestructuringBinding_isCanonical_6628() bool {
	return (!(func() bool {
		sr_10 := this_.ExportedName()
		se_11 := this_.LocalName()
		return j2gort.CheckNotNullString(sr_10) == se_11
	}()) || this_.IsShorthandProperty())
}

// Java method: DestructuringBinding.canonicalizeShorthandProperties()
func (this_ *CheckRequiresSorted_DestructuringBinding) CanonicalizeShorthandProperties() *CheckRequiresSorted_DestructuringBinding {
	if this_.priv_DestructuringBinding_isCanonical_6628() {
		return this_
	} else {
		return CheckRequiresSorted_DestructuringBinding_Of(this_.ExportedName(), this_.LocalName(), true)
	}
}

func (this_ *CheckRequiresSorted_DestructuringBinding) ExportedName() string {
	return this_.exportedName
}
func (this_ *CheckRequiresSorted_DestructuringBinding) LocalName() string { return this_.localName }
func (this_ *CheckRequiresSorted_DestructuringBinding) IsShorthandProperty() bool {
	return this_.isShorthandProperty
}
func (this_ *CheckRequiresSorted_DestructuringBinding) Equals(other_ any) bool {
	o_, ok_ := other_.(*CheckRequiresSorted_DestructuringBinding)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.exportedName != o_.exportedName {
		return false
	}
	if this_.localName != o_.localName {
		return false
	}
	if this_.isShorthandProperty != o_.isShorthandProperty {
		return false
	}
	return true
}
func (this_ *CheckRequiresSorted_DestructuringBinding) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_12 := this_.exportedName
		h_ = 31 * h_
		if hs_12 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_12)
		}
	}
	{
		hs_13 := this_.localName
		h_ = 31 * h_
		if hs_13 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_13)
		}
	}
	if this_.isShorthandProperty {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *CheckRequiresSorted_DestructuringBinding) String() string {
	_ = fmt.Sprint
	return "DestructuringBinding[" + "exportedName=" + j2gort.StrVal(j2gort.BoxString(this_.exportedName)) + ", " + "localName=" + j2gort.StrVal(j2gort.BoxString(this_.localName)) + ", " + "isShorthandProperty=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isShorthandProperty)) + "]"
}

// Java record: com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportStatement
type CheckRequiresSorted_ImportStatement struct {
	nodes        []pkg_rhino.Node                            `gson:"nodes"`
	namespace    string                                      `gson:"namespace"`
	alias        string                                      `gson:"alias"`
	destructures []*CheckRequiresSorted_DestructuringBinding `gson:"destructures"`
	primitive    CheckRequiresSorted_ImportPrimitive         `gson:"primitive"`
}

func (this_ *CheckRequiresSorted_ImportStatement) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckRequiresSorted_ImportStatement) IsJavaNil_() bool { return this_ == nil }
func (*CheckRequiresSorted_ImportStatement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckRequiresSorted$ImportStatement", "ImportStatement"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportStatement", "com.google.javascript.jscomp.lint.CheckRequiresSorted$ImportStatement", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckRequiresSorted_ImportStatement](v_)
	})
}
func (*CheckRequiresSorted_ImportStatement) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *CheckRequiresSorted_ImportStatement) initFields_CheckRequiresSorted_ImportStatement() {
	_ = this_
}

// Java constructor: ImportStatement(com.google.common.collect.ImmutableList<com.google.javascript.rhino.Node> nodes, com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive primitive, java.lang.String namespace, java.lang.@org.jspecify.annotations.Nullable String alias, com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableList<com.google.javascript.jscomp.lint.CheckRequiresSorted.DestructuringBinding> destructures)
func NewCheckRequiresSorted_ImportStatement(nodes j2gort.Ext_com_google_common_collect_ImmutableList, primitive CheckRequiresSorted_ImportPrimitive, namespace string, alias string, destructures j2gort.Ext_com_google_common_collect_ImmutableList) *CheckRequiresSorted_ImportStatement {
	this_ := &CheckRequiresSorted_ImportStatement{}
	InitCheckRequiresSorted_ImportStatement(this_, nodes, primitive, namespace, alias, destructures)
	return this_
}

func InitCheckRequiresSorted_ImportStatement(this_ *CheckRequiresSorted_ImportStatement, nodes j2gort.Ext_com_google_common_collect_ImmutableList, primitive CheckRequiresSorted_ImportPrimitive, namespace string, alias string, destructures j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = nodes
	_ = primitive
	_ = namespace
	_ = alias
	_ = destructures
	this_.initFields_CheckRequiresSorted_ImportStatement()
	j2gort.CheckNotNull(nodes)
	j2gort.CheckNotNullEnum(primitive)
	j2gort.CheckNotNullString(namespace)
	if !((alias == j2gort.NullString) || (destructures == nil)) {
		j2gort.CheckArgumentStr(false, "Import statement cannot be simultaneously aliasing and destructuring")
	}
	this_.nodes = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](nodes)
	this_.primitive = primitive
	this_.namespace = namespace
	this_.alias = alias
	this_.destructures = j2gort.ImmListToSliceNullable[

	// Java static method: ImportStatement.of(com.google.common.collect.ImmutableList<com.google.javascript.rhino.Node> nodes, com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive primitive, java.lang.String namespace, java.lang.@org.jspecify.annotations.Nullable String alias, com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableList<com.google.javascript.jscomp.lint.CheckRequiresSorted.DestructuringBinding> destructures)
	*CheckRequiresSorted_DestructuringBinding](destructures)
}

func CheckRequiresSorted_ImportStatement_Of(nodes j2gort.Ext_com_google_common_collect_ImmutableList, primitive CheckRequiresSorted_ImportPrimitive, namespace string, alias string, destructures j2gort.Ext_com_google_common_collect_ImmutableList) *CheckRequiresSorted_ImportStatement {
	_ = nodes
	_ = primitive
	_ = namespace
	_ = alias
	_ = destructures
	return NewCheckRequiresSorted_ImportStatement(nodes, primitive, namespace, alias, destructures)
}

// Java method: ImportStatement.isStandalone()
func (this_ *CheckRequiresSorted_ImportStatement) IsStandalone() bool {
	return (!(this_.IsAliasing()) && !(this_.IsDestructuring()))
}

// Java method: ImportStatement.isAliasing()
func (this_ *CheckRequiresSorted_ImportStatement) IsAliasing() bool {
	return (this_.Alias() != j2gort.NullString)
}

// Java method: ImportStatement.isDestructuring()
func (this_ *CheckRequiresSorted_ImportStatement) IsDestructuring() bool {
	return (this_.Destructures() != nil)
}

// Java method: ImportStatement.upgrade(com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportPrimitive otherPrimitive)
func (this_ *CheckRequiresSorted_ImportStatement) Upgrade(otherPrimitive CheckRequiresSorted_ImportPrimitive) *CheckRequiresSorted_ImportStatement {
	_ = otherPrimitive
	if CheckRequiresSorted_ImportPrimitive_Stronger(this_.Primitive(), otherPrimitive) != this_.Primitive() {
		return NewCheckRequiresSorted_ImportStatement(this_.Nodes(), otherPrimitive, this_.Namespace(), this_.Alias(), this_.Destructures())
	}
	return this_
}

// Java method: ImportStatement.formatWithoutDoc()
func (this_ *CheckRequiresSorted_ImportStatement) priv_ImportStatement_formatWithoutDoc_6637() string {
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	if !(this_.IsStandalone()) {
		sb.Append_1_string("const ")
	}
	if this_.Alias() != j2gort.NullString {
		sb.Append_1_string(this_.Alias())
	}
	if this_.Destructures() != nil {
		sb.Append_1_string("{")
		var first bool = true
		_ = first
		{

			for _, binding := range j2gort.CheckNotNullSlice(this_.
				Slice_Destructures()) {

				_ = binding
				var exportedName string = binding.ExportedName()
				_ = exportedName
				var localName string = binding.LocalName()
				_ = localName
				if first {
					first = false
				} else {
					sb.Append_1_string(", ")
				}
				sb.Append_1_string(exportedName)
				if !(j2gort.CheckNotNullString(exportedName) == localName) {
					sb.Append_1_string(": ")
					sb.Append_1_string(localName)
				}
			}
		}
		sb.Append_1_string("}")
	}
	if !(this_.IsStandalone()) {
		sb.Append_1_string(" = ")
	}
	sb.Append_1_any_9926(j2gort.BoxNullableVal(this_.Primitive()))
	sb.Append_1_string("('")
	sb.Append_1_string(this_.Namespace())
	sb.Append_1_string("');")
	return sb.String()
}

// Java method: ImportStatement.format()
func (this_ *CheckRequiresSorted_ImportStatement) Format() string {
	var sb j2gort.LocalStringBuilder

	{

		for _, node := range j2gort.CheckNotNullSlice(this_.
			Slice_Nodes()) {

			_ = node
			var comment string = node.GetNonJSDocCommentString()
			_ = comment
			if !(j2gort.CheckNotNullString(comment) == "") {
				sb.AppendStr(node.GetNonJSDocCommentString()).AppendStr("\n")
			}
			var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(node)
			_ = jsDoc
			if jsDoc != nil {
				sb.AppendStr(jsDoc.GetOriginalCommentString()).AppendStr("\n")
			}
		}
	}
	sb.AppendStr(this_.priv_ImportStatement_formatWithoutDoc_6637())
	return sb.String()
}

// Java method: ImportStatement.compareTo(com.google.javascript.jscomp.lint.CheckRequiresSorted.ImportStatement other)
func (this_ *CheckRequiresSorted_ImportStatement) CompareTo(arg_15 any) int32 {
	other := j2gort.CastPtr[CheckRequiresSorted_ImportStatement](arg_15)
	_ = other
	return j2gort.JStrCompareTo(this_.priv_ImportStatement_formatWithoutDoc_6637(), other.priv_ImportStatement_formatWithoutDoc_6637())
}

func (this_ *CheckRequiresSorted_ImportStatement) Nodes() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.nodes)
}
func (this_ *CheckRequiresSorted_ImportStatement) Primitive() CheckRequiresSorted_ImportPrimitive {
	return this_.primitive
}
func (this_ *CheckRequiresSorted_ImportStatement) Namespace() string { return this_.namespace }
func (this_ *CheckRequiresSorted_ImportStatement) Alias() string     { return this_.alias }
func (this_ *CheckRequiresSorted_ImportStatement) Destructures() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.destructures)
}
func (this_ *CheckRequiresSorted_ImportStatement) Equals(other_ any) bool {
	o_, ok_ := other_.(*CheckRequiresSorted_ImportStatement)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(j2gort.ImmListFromTypedSlice(this_.nodes), j2gort.ImmListFromTypedSlice(o_.nodes)) {
		return false
	}
	if this_.primitive != o_.primitive {
		return false
	}
	if this_.namespace != o_.namespace {
		return false
	}
	if this_.alias != o_.alias {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.ImmListFromTypedSlice(this_.destructures), j2gort.ImmListFromTypedSlice(o_.destructures)) {
		return false
	}
	return true
}
func (this_ *CheckRequiresSorted_ImportStatement) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + j2gort.ObjectHash(j2gort.ImmListFromTypedSlice(this_.nodes))
	h_ = 31*h_ + (this_.primitive).HashCode()
	{
		hs_16 := this_.namespace
		h_ = 31 * h_
		if hs_16 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_16)
		}
	}
	{
		hs_17 := this_.alias
		h_ = 31 * h_
		if hs_17 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_17)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(j2gort.ImmListFromTypedSlice(this_.destructures))
	return h_
}
func (this_ *CheckRequiresSorted_ImportStatement) String() string {
	_ = fmt.Sprint
	return "ImportStatement[" + "nodes=" + j2gort.StrVal(j2gort.ImmListFromTypedSlice(this_.nodes)) + ", " + "primitive=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.primitive)) + ", " + "namespace=" + j2gort.StrVal(j2gort.BoxString(this_.namespace)) + ", " + "alias=" + j2gort.StrVal(j2gort.BoxString(this_.alias)) + ", " + "destructures=" + j2gort.StrVal(j2gort.ImmListFromTypedSlice(this_.destructures)) + "]"
}
func (this_ *CheckRequiresSorted_ImportStatement) Slice_Nodes() []pkg_rhino.Node {
	return this_.nodes
}
func (this_ *CheckRequiresSorted_ImportStatement) Slice_Destructures() []*CheckRequiresSorted_DestructuringBinding {
	return this_.destructures
}
func (*CheckRequiresSorted) CanReleaseSyntacticScopes_() {}
