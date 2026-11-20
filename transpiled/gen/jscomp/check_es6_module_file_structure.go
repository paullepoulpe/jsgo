/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckEs6ModuleFileStructure.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure
type CheckEs6ModuleFileStructure struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	compiler          AbstractCompiler              `gson:"compiler"`
	orderedStatements *j2gort.Ext_java_util_TreeSet `gson:"orderedStatements"`
}

func (this_ *CheckEs6ModuleFileStructure) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckEs6ModuleFileStructure) IsJavaNil_() bool { return this_ == nil }
func (*CheckEs6ModuleFileStructure) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure", "CheckEs6ModuleFileStructure"
}
func (this_ *CheckEs6ModuleFileStructure) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*CheckEs6ModuleFileStructure) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure", "com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckEs6ModuleFileStructure](v_)
	})
}
func (*CheckEs6ModuleFileStructure) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {}
func (*CheckEs6ModuleFileStructure) Implements_jscomp_CompilerPass()                           {}
func (*CheckEs6ModuleFileStructure) Implements_jscomp_NodeTraversal_Callback()                 {}

var static_CheckEs6ModuleFileStructure_MUST_COME_BEFORE *DiagnosticType

func StaticGet_CheckEs6ModuleFileStructure_MUST_COME_BEFORE() *DiagnosticType {
	clinit_CheckEs6ModuleFileStructure()
	return static_CheckEs6ModuleFileStructure_MUST_COME_BEFORE
}
func StaticSet_CheckEs6ModuleFileStructure_MUST_COME_BEFORE(v *DiagnosticType) {
	clinit_CheckEs6ModuleFileStructure()
	static_CheckEs6ModuleFileStructure_MUST_COME_BEFORE = v
}

var static_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID pkg_rhino.QualifiedName

func StaticGet_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID() pkg_rhino.QualifiedName {
	clinit_CheckEs6ModuleFileStructure()
	return static_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID
}
func StaticSet_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID(v pkg_rhino.QualifiedName) {
	clinit_CheckEs6ModuleFileStructure()
	static_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID = v
}

var static_CheckEs6ModuleFileStructure_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_CheckEs6ModuleFileStructure_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_CheckEs6ModuleFileStructure()
	return static_CheckEs6ModuleFileStructure_GOOG_REQUIRE
}
func StaticSet_CheckEs6ModuleFileStructure_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_CheckEs6ModuleFileStructure()
	static_CheckEs6ModuleFileStructure_GOOG_REQUIRE = v
}

var clinitDone_CheckEs6ModuleFileStructure uint32
var clinitState_CheckEs6ModuleFileStructure j2gort.ClinitState

func clinit_CheckEs6ModuleFileStructure() {
	if atomic.LoadUint32(&clinitDone_CheckEs6ModuleFileStructure) != 0 {
		return
	}
	clinitSlow_CheckEs6ModuleFileStructure()
}

func clinitSlow_CheckEs6ModuleFileStructure() {
	if !j2gort.EnterClinit(&clinitState_CheckEs6ModuleFileStructure, &clinitDone_CheckEs6ModuleFileStructure) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckEs6ModuleFileStructure, &clinitDone_CheckEs6ModuleFileStructure)
	static_CheckEs6ModuleFileStructure_MUST_COME_BEFORE = DiagnosticType_Warning("JSC_MUST_COME_BEFORE_IN_ES6_MODULE", "In ES6 modules, {0} should come before {1}.")
	static_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID = pkg_rhino.QualifiedName_Of("goog.declareModuleId")
	static_CheckEs6ModuleFileStructure_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
}

func (this_ *CheckEs6ModuleFileStructure) initFields_CheckEs6ModuleFileStructure() {
	_ = this_
	this_.orderedStatements = j2gort.NewExt_java_util_TreeSet()
}

// Java constructor: CheckEs6ModuleFileStructure(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckEs6ModuleFileStructure(compiler AbstractCompiler) *CheckEs6ModuleFileStructure {
	clinit_CheckEs6ModuleFileStructure()
	this_ := &CheckEs6ModuleFileStructure{}
	InitCheckEs6ModuleFileStructure(this_, compiler)
	return this_
}

func InitCheckEs6ModuleFileStructure(this_ *CheckEs6ModuleFileStructure, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_CheckEs6ModuleFileStructure()
	this_.compiler = compiler
}

// Java method: CheckEs6ModuleFileStructure.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckEs6ModuleFileStructure) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */ :
		return true
	case int32(94) /* Token.SCRIPT */ :
		return n.GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())
	case int32(123) /* Token.MODULE_BODY */ :
		this_.orderedStatements.Clear()
		return true
	case int32(116) /* Token.IMPORT */ :
		this_.priv_CheckEs6ModuleFileStructure_visitImport_6570(t, n)
		return false
	case int32(93) /* Token.EXPR_RESULT */ :
		return this_.priv_CheckEs6ModuleFileStructure_visitExprResult_6568(t, n, parent)
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		return this_.priv_CheckEs6ModuleFileStructure_visitDeclaration_6569(t, n, parent)
	default:
		if (parent != 0) && parent.IsModuleBody() {
			this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, n, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_OTHER())
		}
		return false
	}
	panic("unreachable switch expression")
}

// Java method: CheckEs6ModuleFileStructure.checkOrder(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure.OrderedStatement statement)
func (this_ *CheckEs6ModuleFileStructure) priv_CheckEs6ModuleFileStructure_checkOrder_6567(t *NodeTraversal, n pkg_rhino.Node, statement CheckEs6ModuleFileStructure_OrderedStatement) {
	_ = t
	_ = n
	_ = statement
	this_.orderedStatements.Add_1(j2gort.BoxNullableVal(statement))
	var outOfOrder j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_NavigableSet](this_.orderedStatements.TailSet_2(j2gort.BoxNullableVal(statement), false))
	_ = outOfOrder
	if !(outOfOrder.IsEmpty()) {
		t.Report_3(n, StaticGet_CheckEs6ModuleFileStructure_MUST_COME_BEFORE(), []string{statement.String(), j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](outOfOrder))})
	}
}

// Java method: CheckEs6ModuleFileStructure.visitExprResult(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node exprResult, com.google.javascript.rhino.Node parent)
func (this_ *CheckEs6ModuleFileStructure) priv_CheckEs6ModuleFileStructure_visitExprResult_6568(t *NodeTraversal, exprResult pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = exprResult
	_ = parent
	if parent.IsModuleBody() && exprResult.GetFirstChild().IsCall() {
		var call pkg_rhino.Node = exprResult.GetFirstChild()
		_ = call
		if StaticGet_CheckEs6ModuleFileStructure_GOOG_DECLAREMODULEID().Matches(call.GetFirstChild()) {
			this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, call, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID())
			return false
		} else if StaticGet_CheckEs6ModuleFileStructure_GOOG_REQUIRE().Matches(call.GetFirstChild()) {
			this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, call, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE())
			return false
		}
	}
	this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, exprResult, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_OTHER())
	return true
}

// Java method: CheckEs6ModuleFileStructure.visitDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node declaration, com.google.javascript.rhino.Node parent)
func (this_ *CheckEs6ModuleFileStructure) priv_CheckEs6ModuleFileStructure_visitDeclaration_6569(t *NodeTraversal, declaration pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = declaration
	_ = parent
	if ((parent.IsModuleBody() && declaration.HasOneChild()) && declaration.GetFirstChild().HasOneChild()) && declaration.GetFirstFirstChild().IsCall() {
		var call pkg_rhino.Node = declaration.GetFirstFirstChild()
		_ = call
		if StaticGet_CheckEs6ModuleFileStructure_GOOG_REQUIRE().Matches(call.GetFirstChild()) {
			this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, call, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE())
			return false
		}
	}
	this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, declaration, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_OTHER())
	return true
}

// Java method: CheckEs6ModuleFileStructure.visitImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode)
func (this_ *CheckEs6ModuleFileStructure) priv_CheckEs6ModuleFileStructure_visitImport_6570(t *NodeTraversal, importNode pkg_rhino.Node) {
	_ = t
	_ = importNode
	this_.priv_CheckEs6ModuleFileStructure_checkOrder_6567(t, importNode, StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_IMPORT())
}

// Java method: CheckEs6ModuleFileStructure.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckEs6ModuleFileStructure) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java enum: com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure.OrderedStatement
type CheckEs6ModuleFileStructure_OrderedStatement uint8

const (
	enum_CheckEs6ModuleFileStructure_OrderedStatement_IMPORT            CheckEs6ModuleFileStructure_OrderedStatement = 1
	enum_CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID CheckEs6ModuleFileStructure_OrderedStatement = 2
	enum_CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE      CheckEs6ModuleFileStructure_OrderedStatement = 3
	enum_CheckEs6ModuleFileStructure_OrderedStatement_OTHER             CheckEs6ModuleFileStructure_OrderedStatement = 4
)

var enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement = [...]j2gort.EnumBase{
	{Name_: "IMPORT", Ordinal_: 0},
	{Name_: "DECLARE_MODULE_ID", Ordinal_: 1},
	{Name_: "GOOG_REQUIRE", Ordinal_: 2},
	{Name_: "OTHER", Ordinal_: 3},
}
var enumValues_CheckEs6ModuleFileStructure_OrderedStatement = []CheckEs6ModuleFileStructure_OrderedStatement{1, 2, 3, 4}

var enumField_CheckEs6ModuleFileStructure_OrderedStatement_description = [...]string{
	"import statements",
	"a call to goog.declareModuleId()",
	"calls to goog.require()",
	"other statements",
}

func (this_ CheckEs6ModuleFileStructure_OrderedStatement) FieldGet_description() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_CheckEs6ModuleFileStructure_OrderedStatement_description[this_-1]
}

func StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_IMPORT() CheckEs6ModuleFileStructure_OrderedStatement {
	return enum_CheckEs6ModuleFileStructure_OrderedStatement_IMPORT
}
func StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID() CheckEs6ModuleFileStructure_OrderedStatement {
	return enum_CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID
}
func StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE() CheckEs6ModuleFileStructure_OrderedStatement {
	return enum_CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE
}
func StaticGet_CheckEs6ModuleFileStructure_OrderedStatement_OTHER() CheckEs6ModuleFileStructure_OrderedStatement {
	return enum_CheckEs6ModuleFileStructure_OrderedStatement_OTHER
}

func CheckEs6ModuleFileStructure_OrderedStatement_Values() []CheckEs6ModuleFileStructure_OrderedStatement {
	return j2gort.CloneSlice(enumValues_CheckEs6ModuleFileStructure_OrderedStatement)
}

func CheckEs6ModuleFileStructure_OrderedStatement_ValueOf(name string) CheckEs6ModuleFileStructure_OrderedStatement {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement {
		if enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement[i_].Name_ == name {
			return CheckEs6ModuleFileStructure_OrderedStatement(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure.OrderedStatement." + name))
}

// Java method: OrderedStatement.toString()
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) String() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_description()
}

func (this_ CheckEs6ModuleFileStructure_OrderedStatement) IsJavaNil_() bool                { return this_ == 0 }
func (CheckEs6ModuleFileStructure_OrderedStatement) HasIdentityEquals_()                   {}
func (CheckEs6ModuleFileStructure_OrderedStatement) Implements_j2gort_Ext_java_lang_Enum() {}
func (CheckEs6ModuleFileStructure_OrderedStatement) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CheckEs6ModuleFileStructure_OrderedStatement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure$OrderedStatement", "OrderedStatement"
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement[this_-1]
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement[this_-1].Name_
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) Equals(other_ any) bool {
	o_, ok_ := other_.(CheckEs6ModuleFileStructure_OrderedStatement)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement[this_-1].HashCode()
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement[this_-1].CompareTo(other_)
}
func (this_ CheckEs6ModuleFileStructure_OrderedStatement) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CheckEs6ModuleFileStructure_OrderedStatement[this_-1].Name_)
}
func init() {
	vals_ := []any{CheckEs6ModuleFileStructure_OrderedStatement(1), CheckEs6ModuleFileStructure_OrderedStatement(2), CheckEs6ModuleFileStructure_OrderedStatement(3), CheckEs6ModuleFileStructure_OrderedStatement(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure.OrderedStatement", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure$OrderedStatement", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure.OrderedStatement", "com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure$OrderedStatement", func(v_ any) bool {
		return j2gort.InstanceOf[CheckEs6ModuleFileStructure_OrderedStatement](v_)
	})
}
func (*CheckEs6ModuleFileStructure) NoOpVisit_()                 {}
func (*CheckEs6ModuleFileStructure) CanReleaseSyntacticScopes_() {}
