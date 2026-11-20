/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckNestedNames.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckNestedNames
type CheckNestedNames struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *CheckNestedNames) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckNestedNames) IsJavaNil_() bool { return this_ == nil }
func (*CheckNestedNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckNestedNames", "CheckNestedNames"
}
func (this_ *CheckNestedNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckNestedNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckNestedNames", "com.google.javascript.jscomp.lint.CheckNestedNames", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckNestedNames](v_)
	})
}
func (*CheckNestedNames) Implements_jscomp_CompilerPass()           {}
func (*CheckNestedNames) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE *DiagnosticType

func StaticGet_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE() *DiagnosticType {
	clinit_CheckNestedNames()
	return static_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE
}
func StaticSet_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE(v *DiagnosticType) {
	clinit_CheckNestedNames()
	static_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE = v
}

var static_CheckNestedNames_NESTED_KIND_TO_REPORT j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_CheckNestedNames_NESTED_KIND_TO_REPORT() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_CheckNestedNames()
	return static_CheckNestedNames_NESTED_KIND_TO_REPORT
}
func StaticSet_CheckNestedNames_NESTED_KIND_TO_REPORT(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_CheckNestedNames()
	static_CheckNestedNames_NESTED_KIND_TO_REPORT = v
}

var clinitDone_CheckNestedNames uint32
var clinitState_CheckNestedNames j2gort.ClinitState

func clinit_CheckNestedNames() {
	if atomic.LoadUint32(&clinitDone_CheckNestedNames) != 0 {
		return
	}
	clinitSlow_CheckNestedNames()
}

func clinitSlow_CheckNestedNames() {
	if !j2gort.EnterClinit(&clinitState_CheckNestedNames, &clinitDone_CheckNestedNames) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckNestedNames, &clinitDone_CheckNestedNames)
	static_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE = DiagnosticType_Disabled("JSC_NESTED_NAME_IN_GOOG_MODULE", "A nested {0} is created on the name `{1}`. Fix this linter finding by converting the module-level static property assignment on `{1}` into a module-level flat name (i.e. change `{1}.prop = ...` into `{1}_prop = ...`. You can (if required) export this flat name using named exports (`exports.{1}_prop = {1}_prop`).")
	static_CheckNestedNames_NESTED_KIND_TO_REPORT = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_8(j2gort.BoxNullableVal(StaticGet_CheckNestedNames_DeclarationKind_CLASS()), "class", j2gort.BoxNullableVal(StaticGet_CheckNestedNames_DeclarationKind_ENUM()), "enum", j2gort.BoxNullableVal(StaticGet_CheckNestedNames_DeclarationKind_INTERFACE()), "interface", j2gort.BoxNullableVal(StaticGet_CheckNestedNames_DeclarationKind_TYPEDEF()), "typedef")
}

func (this_ *CheckNestedNames) initFields_CheckNestedNames() {
	_ = this_
}

// Java constructor: CheckNestedNames(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckNestedNames(compiler AbstractCompiler) *CheckNestedNames {
	clinit_CheckNestedNames()
	this_ := &CheckNestedNames{}
	InitCheckNestedNames(this_, compiler)
	return this_
}

func InitCheckNestedNames(this_ *CheckNestedNames, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckNestedNames()
	this_.compiler = compiler
}

// Java method: CheckNestedNames.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckNestedNames) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckNestedNames.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckNestedNames) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsScript() {
		return n.GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE())
	}
	return nodeTraversal.InGlobalOrModuleScope()
}

// Java method: CheckNestedNames.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckNestedNames) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(t.InGlobalOrModuleScope()) {
		j2gort.CheckArgument(false)
	}
	if !(n.IsGetProp()) || !(n.IsQualifiedName()) {
		return
	}
	var targetGetProp pkg_rhino.Node = CheckNestedNames_GetTargetGetProp(n)
	_ = targetGetProp
	if ((parent.IsAssign() || parent.IsExprResult()) && n.IsFirstChildOf(parent)) && !(CheckNestedNames_IsDotPrototype(targetGetProp)) {
		var owner pkg_rhino.Node = targetGetProp.GetFirstChild()
		_ = owner
		if !(owner.IsName()) || (j2gort.CheckNotNullString(owner.GetString()) == "exports") {
			return
		}
		var ownerName string = owner.GetString()
		_ = ownerName
		var declarationKind CheckNestedNames_DeclarationKind = this_.priv_CheckNestedNames_getNestedDeclarationKind_6605(n, parent)
		_ = declarationKind
		if StaticGet_CheckNestedNames_NESTED_KIND_TO_REPORT().ContainsKey(j2gort.BoxNullableVal(declarationKind)) {
			var nestedKind string = j2gort.CastString(StaticGet_CheckNestedNames_NESTED_KIND_TO_REPORT().Get(j2gort.BoxNullableVal(declarationKind)))
			_ = nestedKind
			t.Report_3(targetGetProp, StaticGet_CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE(), []string{nestedKind, ownerName})
		}
	}
}

// Java static method: CheckNestedNames.getTargetGetProp(com.google.javascript.rhino.Node n)
func CheckNestedNames_GetTargetGetProp(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_CheckNestedNames()
	_ = n
	if !(n.IsGetProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	if n.GetFirstChild().IsGetProp() {
		return CheckNestedNames_GetTargetGetProp(n.GetFirstChild())
	}
	return n
}

// Java static method: CheckNestedNames.isDotPrototype(com.google.javascript.rhino.Node getProp)
func CheckNestedNames_IsDotPrototype(getProp pkg_rhino.Node) bool {
	clinit_CheckNestedNames()
	_ = getProp
	return (getProp.IsGetProp() && (j2gort.CheckNotNullString(getProp.GetString()) == "prototype"))
}

// Java method: CheckNestedNames.getNestedDeclarationKind(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node parent)
func (this_ *CheckNestedNames) priv_CheckNestedNames_getNestedDeclarationKind_6605(lhs pkg_rhino.Node, parent pkg_rhino.Node) CheckNestedNames_DeclarationKind {
	_ = lhs
	_ = parent
	if !(lhs.IsFirstChildOf(parent)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(lhs))
	}
	if !(j2gort.IfaceEquals[pkg_rhino.JSType](lhs.GetTypedefTypeProp(), nil)) {
		return StaticGet_CheckNestedNames_DeclarationKind_TYPEDEF()
	}
	var type_ pkg_rhino.JSType = lhs.GetJSType()
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		if type_.GetJSTypeBase().IsEnumType() {
			return StaticGet_CheckNestedNames_DeclarationKind_ENUM()
		} else if type_.GetJSTypeBase().IsInterface() {
			return StaticGet_CheckNestedNames_DeclarationKind_INTERFACE()
		} else if type_.GetJSTypeBase().IsConstructor() {
			return StaticGet_CheckNestedNames_DeclarationKind_CLASS()
		}
	}
	return CheckNestedNames_DeclarationKind(0)
}

// Java enum: com.google.javascript.jscomp.lint.CheckNestedNames.DeclarationKind
type CheckNestedNames_DeclarationKind uint8

const (
	enum_CheckNestedNames_DeclarationKind_CLASS     CheckNestedNames_DeclarationKind = 1
	enum_CheckNestedNames_DeclarationKind_ENUM      CheckNestedNames_DeclarationKind = 2
	enum_CheckNestedNames_DeclarationKind_TYPEDEF   CheckNestedNames_DeclarationKind = 3
	enum_CheckNestedNames_DeclarationKind_INTERFACE CheckNestedNames_DeclarationKind = 4
)

var enumBaseTable_CheckNestedNames_DeclarationKind = [...]j2gort.EnumBase{
	{Name_: "CLASS", Ordinal_: 0},
	{Name_: "ENUM", Ordinal_: 1},
	{Name_: "TYPEDEF", Ordinal_: 2},
	{Name_: "INTERFACE", Ordinal_: 3},
}
var enumValues_CheckNestedNames_DeclarationKind = []CheckNestedNames_DeclarationKind{1, 2, 3, 4}

func StaticGet_CheckNestedNames_DeclarationKind_CLASS() CheckNestedNames_DeclarationKind {
	return enum_CheckNestedNames_DeclarationKind_CLASS
}
func StaticGet_CheckNestedNames_DeclarationKind_ENUM() CheckNestedNames_DeclarationKind {
	return enum_CheckNestedNames_DeclarationKind_ENUM
}
func StaticGet_CheckNestedNames_DeclarationKind_TYPEDEF() CheckNestedNames_DeclarationKind {
	return enum_CheckNestedNames_DeclarationKind_TYPEDEF
}
func StaticGet_CheckNestedNames_DeclarationKind_INTERFACE() CheckNestedNames_DeclarationKind {
	return enum_CheckNestedNames_DeclarationKind_INTERFACE
}

func CheckNestedNames_DeclarationKind_Values() []CheckNestedNames_DeclarationKind {
	return j2gort.CloneSlice(enumValues_CheckNestedNames_DeclarationKind)
}

func CheckNestedNames_DeclarationKind_ValueOf(name string) CheckNestedNames_DeclarationKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckNestedNames_DeclarationKind {
		if enumBaseTable_CheckNestedNames_DeclarationKind[i_].Name_ == name {
			return CheckNestedNames_DeclarationKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.lint.CheckNestedNames.DeclarationKind." + name))
}

func (this_ CheckNestedNames_DeclarationKind) IsJavaNil_() bool                              { return this_ == 0 }
func (CheckNestedNames_DeclarationKind) HasIdentityEquals_()                                 {}
func (CheckNestedNames_DeclarationKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CheckNestedNames_DeclarationKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CheckNestedNames_DeclarationKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckNestedNames_DeclarationKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckNestedNames$DeclarationKind", "DeclarationKind"
}
func (this_ CheckNestedNames_DeclarationKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckNestedNames_DeclarationKind[this_-1]
}
func (this_ CheckNestedNames_DeclarationKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckNestedNames_DeclarationKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CheckNestedNames_DeclarationKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckNestedNames_DeclarationKind[this_-1].Name_
}
func (this_ CheckNestedNames_DeclarationKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CheckNestedNames_DeclarationKind[this_-1].Name_
}
func (this_ CheckNestedNames_DeclarationKind) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckNestedNames_DeclarationKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckNestedNames_DeclarationKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckNestedNames_DeclarationKind[this_-1].HashCode()
}
func (this_ CheckNestedNames_DeclarationKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckNestedNames_DeclarationKind[this_-1].CompareTo(other_)
}
func (this_ CheckNestedNames_DeclarationKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckNestedNames_DeclarationKind[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckNestedNames_DeclarationKind(1), CheckNestedNames_DeclarationKind(2), CheckNestedNames_DeclarationKind(3), CheckNestedNames_DeclarationKind(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckNestedNames.DeclarationKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckNestedNames$DeclarationKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckNestedNames.DeclarationKind", "com.google.javascript.jscomp.lint.CheckNestedNames$DeclarationKind", func(v_ any) bool {
		return j2gort.InstanceOf[CheckNestedNames_DeclarationKind](v_)
	})
}
func (*CheckNestedNames) CanReleaseSyntacticScopes_() {}
