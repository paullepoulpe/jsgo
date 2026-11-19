/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ImplicitNullabilityCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ImplicitNullabilityCheck
type ImplicitNullabilityCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *ImplicitNullabilityCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *ImplicitNullabilityCheck) IsJavaNil_() bool { return this_ == nil }
func (*ImplicitNullabilityCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ImplicitNullabilityCheck", "ImplicitNullabilityCheck"
}
func (this_ *ImplicitNullabilityCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ImplicitNullabilityCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ImplicitNullabilityCheck", "com.google.javascript.jscomp.ImplicitNullabilityCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*ImplicitNullabilityCheck](v_)
	})
}
func (*ImplicitNullabilityCheck) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ImplicitNullabilityCheck) Implements_jscomp_CompilerPass()                            {}
func (*ImplicitNullabilityCheck) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC *DiagnosticType

func StaticGet_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC() *DiagnosticType {
	clinit_ImplicitNullabilityCheck()
	return static_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC
}
func StaticSet_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC(v *DiagnosticType) {
	clinit_ImplicitNullabilityCheck()
	static_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC = v
}

var static_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC *DiagnosticType

func StaticGet_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC() *DiagnosticType {
	clinit_ImplicitNullabilityCheck()
	return static_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC
}
func StaticSet_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC(v *DiagnosticType) {
	clinit_ImplicitNullabilityCheck()
	static_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC = v
}

var static_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ImplicitNullabilityCheck()
	return static_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES
}
func StaticSet_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ImplicitNullabilityCheck()
	static_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES = v
}

var clinitDone_ImplicitNullabilityCheck uint32
var clinitState_ImplicitNullabilityCheck j2gort.ClinitState

func clinit_ImplicitNullabilityCheck() {
	if atomic.LoadUint32(&clinitDone_ImplicitNullabilityCheck) != 0 {
		return
	}
	clinitSlow_ImplicitNullabilityCheck()
}

func clinitSlow_ImplicitNullabilityCheck() {
	if !j2gort.EnterClinit(&clinitState_ImplicitNullabilityCheck, &clinitDone_ImplicitNullabilityCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ImplicitNullabilityCheck, &clinitDone_ImplicitNullabilityCheck)
	static_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC = DiagnosticType_Disabled("JSC_IMPLICITLY_NULLABLE_JSDOC", "Name {0} in JSDoc is implicitly nullable, and is discouraged by the style guide.\nPlease add a '!' to make it non-nullable, or a '?' to make it explicitly nullable.")
	static_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC = DiagnosticType_Disabled("JSC_IMPLICITLY_NONNULL_JSDOC", "Name {0} in JSDoc is implicitly non-null, and is discouraged by the style guide.\nPlease add a '!' to make it explicit.")
	static_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("*", "?", "bigint", "boolean", "null", "number", []any{"string", "symbol", "undefined", "void"})
}

func (this_ *ImplicitNullabilityCheck) initFields_ImplicitNullabilityCheck() {
	_ = this_
}

// Java constructor: ImplicitNullabilityCheck(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewImplicitNullabilityCheck(compiler AbstractCompiler) *ImplicitNullabilityCheck {
	clinit_ImplicitNullabilityCheck()
	this_ := &ImplicitNullabilityCheck{}
	InitImplicitNullabilityCheck(this_, compiler)
	return this_
}

func InitImplicitNullabilityCheck(this_ *ImplicitNullabilityCheck, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ImplicitNullabilityCheck()
	this_.compiler = compiler
}

// Java method: ImplicitNullabilityCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ImplicitNullabilityCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}

// Java static method: ImplicitNullabilityCheck.findImplicitNullabilityResults(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.jscomp.NodeTraversal t)
func ImplicitNullabilityCheck_FindImplicitNullabilityResults(info *pkg_rhino.JSDocInfo, t *NodeTraversal) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_ImplicitNullabilityCheck()
	_ = info
	_ = t
	if info == nil {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var typeNodeCount int32 = info.GetTypeNodeCount()
	_ = typeNodeCount
	if typeNodeCount == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < typeNodeCount; i = i + 1 {
			var typeRoot pkg_rhino.Node = info.GetTypeNodeAt(i)
			_ = typeRoot
			NodeUtil_VisitPreOrder_3(typeRoot, NewImplicitNullabilityCheck__Anon_52(t, builder), j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
		}
	}
	return builder.BuildImmList_()
}

// Java method: ImplicitNullabilityCheck.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node p)
func (this_ *ImplicitNullabilityCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = t
	_ = n
	_ = p
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(ImplicitNullabilityCheck_FindImplicitNullabilityResults(info, t)))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			r := j2gort.CastPtr[ImplicitNullabilityCheck_Result](rawElem_2)
			_ = r
			var stringNode pkg_rhino.Node = r.GetNode()
			_ = stringNode
			var dt *DiagnosticType
			_ = dt
			if r.GetNullability().IsNullable() {
				dt = StaticGet_ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC()
			} else {
				dt = StaticGet_ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC()
			}
			this_.compiler.Report_1(JSError_Make_3(stringNode, dt, []string{stringNode.GetString()}))
		}
	}
}

// Java enum: com.google.javascript.jscomp.ImplicitNullabilityCheck.Nullability
type ImplicitNullabilityCheck_Nullability uint8

const (
	enum_ImplicitNullabilityCheck_Nullability_NONNULL  ImplicitNullabilityCheck_Nullability = 1
	enum_ImplicitNullabilityCheck_Nullability_NULLABLE ImplicitNullabilityCheck_Nullability = 2
)

var enumBaseTable_ImplicitNullabilityCheck_Nullability = [...]j2gort.EnumBase{
	{Name_: "NONNULL", Ordinal_: 0},
	{Name_: "NULLABLE", Ordinal_: 1},
}
var enumValues_ImplicitNullabilityCheck_Nullability = []ImplicitNullabilityCheck_Nullability{1, 2}

func StaticGet_ImplicitNullabilityCheck_Nullability_NONNULL() ImplicitNullabilityCheck_Nullability {
	return enum_ImplicitNullabilityCheck_Nullability_NONNULL
}
func StaticGet_ImplicitNullabilityCheck_Nullability_NULLABLE() ImplicitNullabilityCheck_Nullability {
	return enum_ImplicitNullabilityCheck_Nullability_NULLABLE
}

func ImplicitNullabilityCheck_Nullability_Values() []ImplicitNullabilityCheck_Nullability {
	return j2gort.CloneSlice(enumValues_ImplicitNullabilityCheck_Nullability)
}

func ImplicitNullabilityCheck_Nullability_ValueOf(name string) ImplicitNullabilityCheck_Nullability {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ImplicitNullabilityCheck_Nullability {
		if enumBaseTable_ImplicitNullabilityCheck_Nullability[i_].Name_ == name {
			return ImplicitNullabilityCheck_Nullability(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ImplicitNullabilityCheck.Nullability." + name))
}

// Java method: Nullability.isNullable()
func (this_ ImplicitNullabilityCheck_Nullability) IsNullable() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_ImplicitNullabilityCheck_Nullability_NULLABLE())
}

func (this_ ImplicitNullabilityCheck_Nullability) IsJavaNil_() bool                              { return this_ == 0 }
func (ImplicitNullabilityCheck_Nullability) HasIdentityEquals_()                                 {}
func (ImplicitNullabilityCheck_Nullability) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ImplicitNullabilityCheck_Nullability) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ImplicitNullabilityCheck_Nullability) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ImplicitNullabilityCheck_Nullability) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ImplicitNullabilityCheck$Nullability", "Nullability"
}
func (this_ ImplicitNullabilityCheck_Nullability) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ImplicitNullabilityCheck_Nullability[this_-1]
}
func (this_ ImplicitNullabilityCheck_Nullability) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ImplicitNullabilityCheck_Nullability) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ImplicitNullabilityCheck_Nullability) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ImplicitNullabilityCheck_Nullability[this_-1].Name_
}
func (this_ ImplicitNullabilityCheck_Nullability) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ImplicitNullabilityCheck_Nullability[this_-1].Name_
}
func (this_ ImplicitNullabilityCheck_Nullability) Equals(other_ any) bool {
	o_, ok_ := other_.(ImplicitNullabilityCheck_Nullability)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ImplicitNullabilityCheck_Nullability) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ImplicitNullabilityCheck_Nullability[this_-1].HashCode()
}
func (this_ ImplicitNullabilityCheck_Nullability) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ImplicitNullabilityCheck_Nullability[this_-1].CompareTo(other_)
}
func (this_ ImplicitNullabilityCheck_Nullability) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ImplicitNullabilityCheck_Nullability[this_-1].Name_)
}
func init() {
	vals_ := []any{ImplicitNullabilityCheck_Nullability(1), ImplicitNullabilityCheck_Nullability(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ImplicitNullabilityCheck.Nullability", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ImplicitNullabilityCheck$Nullability", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ImplicitNullabilityCheck.Nullability", "com.google.javascript.jscomp.ImplicitNullabilityCheck$Nullability", func(v_ any) bool {
		return j2gort.InstanceOf[ImplicitNullabilityCheck_Nullability](v_)
	})
}

// Java class: com.google.javascript.jscomp.ImplicitNullabilityCheck.Result
type ImplicitNullabilityCheck_Result struct {
	node        pkg_rhino.Node `gson:"node"`
	idHash_     atomic.Int32
	nullability ImplicitNullabilityCheck_Nullability `gson:"nullability"`
}

func (this_ *ImplicitNullabilityCheck_Result) BindSelf_(s any) {
	_ = s
}
func (this_ *ImplicitNullabilityCheck_Result) IsJavaNil_() bool { return this_ == nil }
func (*ImplicitNullabilityCheck_Result) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ImplicitNullabilityCheck$Result", "Result"
}
func (this_ *ImplicitNullabilityCheck_Result) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ImplicitNullabilityCheck_Result) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ImplicitNullabilityCheck.Result", "com.google.javascript.jscomp.ImplicitNullabilityCheck$Result", func(v_ any) bool {
		return j2gort.InstanceOf[*ImplicitNullabilityCheck_Result](v_)
	})
}

func (this_ *ImplicitNullabilityCheck_Result) initFields_ImplicitNullabilityCheck_Result() {
	_ = this_
}

// Java constructor: Result(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.ImplicitNullabilityCheck.Nullability nullability)
func NewImplicitNullabilityCheck_Result(node pkg_rhino.Node, nullability ImplicitNullabilityCheck_Nullability) *ImplicitNullabilityCheck_Result {
	this_ := &ImplicitNullabilityCheck_Result{}
	InitImplicitNullabilityCheck_Result(this_, node, nullability)
	return this_
}

func InitImplicitNullabilityCheck_Result(this_ *ImplicitNullabilityCheck_Result, node pkg_rhino.Node, nullability ImplicitNullabilityCheck_Nullability) {
	_ = node
	_ = nullability
	this_.initFields_ImplicitNullabilityCheck_Result()
	if !(node.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	this_.node = node
	this_.nullability = nullability
}

// Java static method: Result.create(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.ImplicitNullabilityCheck.Nullability nullability)
func ImplicitNullabilityCheck_Result_Create(node pkg_rhino.Node, nullability ImplicitNullabilityCheck_Nullability) *ImplicitNullabilityCheck_Result {
	_ = node
	_ = nullability
	return NewImplicitNullabilityCheck_Result(node, nullability)
}

// Java method: Result.getNullability()
func (this_ *ImplicitNullabilityCheck_Result) GetNullability() ImplicitNullabilityCheck_Nullability {
	return this_.nullability
}

// Java method: Result.getNode()
func (this_ *ImplicitNullabilityCheck_Result) GetNode() pkg_rhino.Node {
	return this_.node
}

// Java class:
type ImplicitNullabilityCheck__Anon_52 struct {
	cap_0_t       *NodeTraversal
	cap_1_builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder
	idHash_       atomic.Int32
}

func (this_ *ImplicitNullabilityCheck__Anon_52) BindSelf_(s any) {
	_ = s
}
func (this_ *ImplicitNullabilityCheck__Anon_52) IsJavaNil_() bool { return this_ == nil }
func (*ImplicitNullabilityCheck__Anon_52) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ImplicitNullabilityCheck$1", ""
}
func (this_ *ImplicitNullabilityCheck__Anon_52) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ImplicitNullabilityCheck__Anon_52) HasIdentityEquals_()                 {}
func (*ImplicitNullabilityCheck__Anon_52) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *ImplicitNullabilityCheck__Anon_52) initFields_ImplicitNullabilityCheck__Anon_52() {
	_ = this_
}

func NewImplicitNullabilityCheck__Anon_52(cap_p_t *NodeTraversal, cap_p_builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder) *ImplicitNullabilityCheck__Anon_52 {
	this_ := &ImplicitNullabilityCheck__Anon_52{}
	this_.cap_0_t = cap_p_t
	this_.cap_1_builder = cap_p_builder
	this_.initFields_ImplicitNullabilityCheck__Anon_52()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node node)
func (this_ *ImplicitNullabilityCheck__Anon_52) Visit(node pkg_rhino.Node) {
	_ = node
	if !(node.IsStringLit()) {
		return
	}
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	if parent != 0 {
		switch parent.TokenSwitchOrd_() {
		case int32(159) /* Token.BANG */, int32(158) /* Token.QMARK */, int32(37) /* Token.THIS */, int32(23) /* Token.NEW */, int32(25) /* Token.TYPEOF */ :
			return
		case int32(155) /* Token.PIPE */ :
			var gp pkg_rhino.Node = parent.GetParent()
			_ = gp
			if (gp != 0) && (gp.RawToken_() == uint8(158)) {
				return
			}
			{
				var child pkg_rhino.Node = parent.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					if (child.IsStringLit() && (j2gort.CheckNotNullString(child.GetString()) == "null")) || (child.RawToken_() == uint8(158)) {
						return
					}
				}
			}
		default:
		}
	}
	var typeName string = node.GetString()
	_ = typeName
	if j2gort.SetContainsStr(StaticGet_ImplicitNullabilityCheck_NULLABILITY_OMITTED_TYPES(), typeName) {
		return
	}
	var registry *pkg_rhino.JSTypeRegistry = this_.cap_0_t.GetCompiler().GetTypeRegistry()
	_ = registry
	var type_ pkg_rhino.JSType = registry.GetType_2(pkg_rhino.ToIface_StaticScope(this_.cap_0_t.GetScope()), typeName)
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return
	}
	var isNonNullableName bool = (registry.IsNonNullableName(pkg_rhino.ToIface_StaticScope(this_.cap_0_t.GetScope()), typeName) && !(type_.IsNullable()))
	_ = isNonNullableName
	var nullability ImplicitNullabilityCheck_Nullability
	_ = nullability
	if isNonNullableName {
		nullability = StaticGet_ImplicitNullabilityCheck_Nullability_NONNULL()
	} else {
		nullability = StaticGet_ImplicitNullabilityCheck_Nullability_NULLABLE()
	}
	this_.cap_1_builder.AddVoid_(j2gort.PtrToAny(ImplicitNullabilityCheck_Result_Create(node, nullability)))
}
func (*ImplicitNullabilityCheck) AlwaysShouldTraverse_()      {}
func (*ImplicitNullabilityCheck) CanReleaseSyntacticScopes_() {}
