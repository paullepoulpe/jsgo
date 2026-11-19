/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ClosureExportTypeClassifier.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_closure_export_type_classifier_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_closure_export_type_classifier_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_19 any) any {
		unused := j2gort.CastString(lp_19)
		_ = unused
		return NewClosureExportTypeClassifier_NameInfo()
	})
}

// Java class: com.google.javascript.jscomp.ClosureExportTypeClassifier
type ClosureExportTypeClassifier struct {
	topLevelNames     j2gort.OrderedMap[string, *ClosureExportTypeClassifier_NameInfo] `gson:"topLevelNames"`
	forcedNamedExport j2gort.OrderedSet[string]                                        `gson:"forcedNamedExport"`
	idHash_           atomic.Int32
}

func (this_ *ClosureExportTypeClassifier) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureExportTypeClassifier) IsJavaNil_() bool { return this_ == nil }
func (*ClosureExportTypeClassifier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureExportTypeClassifier", "ClosureExportTypeClassifier"
}
func (this_ *ClosureExportTypeClassifier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureExportTypeClassifier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureExportTypeClassifier", "com.google.javascript.jscomp.ClosureExportTypeClassifier", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureExportTypeClassifier](v_)
	})
}

func StaticGet_ClosureExportTypeClassifier_EXPORTS() string        { return "exports" }
func StaticGet_ClosureExportTypeClassifier_MAX_ALIAS_DEPTH() int32 { return int32(4) }
func (this_ *ClosureExportTypeClassifier) initFields_ClosureExportTypeClassifier() {
	_ = this_
	this_.topLevelNames = j2gort.NewOrderedMap[string, *ClosureExportTypeClassifier_NameInfo]()
	this_.forcedNamedExport = j2gort.NewOrderedSet[string]()
}

// Java constructor: ClosureExportTypeClassifier()
func NewClosureExportTypeClassifier() *ClosureExportTypeClassifier {
	this_ := &ClosureExportTypeClassifier{}
	InitClosureExportTypeClassifier(this_)
	return this_
}

func InitClosureExportTypeClassifier(this_ *ClosureExportTypeClassifier) {
	this_.initFields_ClosureExportTypeClassifier()
}

// Java static method: ClosureExportTypeClassifier.classify(com.google.javascript.rhino.Node moduleRoot, java.lang.String googNamespace)
func ClosureExportTypeClassifier_Classify(moduleRoot pkg_rhino.Node, googNamespace string) ClosureExportTypeClassifier_TsExportType {
	_ = moduleRoot
	_ = googNamespace
	var val_1 ClosureExportTypeClassifier
	InitClosureExportTypeClassifier(&val_1)
	var classifier *ClosureExportTypeClassifier = &val_1
	_ = classifier
	classifier.priv_ClosureExportTypeClassifier_addGoogModuleExports_1165()
	ClosureExportTypeClassifier_ForEachTopLevelStatement(moduleRoot, func() j2gort.Ext_java_util_function_Consumer {
		bound_2 := classifier
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_2.priv_ClosureExportTypeClassifier_visitStatement_1166(j2gort.CastNode[pkg_rhino.Node](mrp0_))
		})
	}())
	var info *ClosureExportTypeClassifier_NameInfo = classifier.topLevelNames.Get("exports")
	_ = info
	if classifier.forcedNamedExport.Contains(googNamespace) || (classifier.priv_ClosureExportTypeClassifier_kindOfName_1173(info, int32(0)) == StaticGet_ClosureExportTypeClassifier_Kind_NAMESPACE()) {
		return StaticGet_ClosureExportTypeClassifier_TsExportType_NAMED()
	}
	return StaticGet_ClosureExportTypeClassifier_TsExportType_DEFAULT()
}

// Java method: ClosureExportTypeClassifier.addGoogModuleExports()
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_addGoogModuleExports_1165() {
	var val_3 ClosureExportTypeClassifier_NameInfo
	InitClosureExportTypeClassifier_NameInfo(&val_3)
	var info *ClosureExportTypeClassifier_NameInfo = &val_3
	_ = info
	{
		r_4 := info
		r_4.isGoogModuleExports = true
	}
	this_.topLevelNames.Put("exports", info)
}

// Java static method: ClosureExportTypeClassifier.forEachTopLevelStatement(com.google.javascript.rhino.Node moduleRoot, java.util.function.Consumer<com.google.javascript.rhino.Node> callback)
func ClosureExportTypeClassifier_ForEachTopLevelStatement(moduleRoot pkg_rhino.Node, callback j2gort.Ext_java_util_function_Consumer) {
	_ = moduleRoot
	_ = callback
	if !(moduleRoot.IsScript() && moduleRoot.GetFirstChild().IsModuleBody()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(moduleRoot))
	}
	var body pkg_rhino.Node = moduleRoot.GetFirstChild()
	_ = body
	{
		var statement pkg_rhino.Node = body.GetFirstChild()
		_ = statement
		for ; statement != 0; statement = statement.GetNext() {
			callback.Accept(j2gort.BoxNode(statement))
		}
	}
}

// Java method: ClosureExportTypeClassifier.visitStatement(com.google.javascript.rhino.Node statement)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_visitStatement_1166(statement pkg_rhino.Node) {
	_ = statement
	switch statement.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.priv_ClosureExportTypeClassifier_visitNameDeclaration_1168(statement)
	case int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
		var name pkg_rhino.Node = statement.GetFirstChild()
		_ = name
		var info *ClosureExportTypeClassifier_NameInfo = this_.priv_ClosureExportTypeClassifier_getOrCreate_1172(name.GetString())
		_ = info
		{
			r_5 := info
			r_5.isNeverNamespace = true
		}
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.priv_ClosureExportTypeClassifier_maybeRecordLegacyNamedTsExport_1167(statement)
		this_.priv_ClosureExportTypeClassifier_visitExpression_1170(statement.GetFirstChild())
	default:
	}
}

// Java method: ClosureExportTypeClassifier.maybeRecordLegacyNamedTsExport(com.google.javascript.rhino.Node statement)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_maybeRecordLegacyNamedTsExport_1167(statement pkg_rhino.Node) {
	_ = statement
	if !(NodeUtil_IsGoogModuleCall(statement)) {
		return
	}
	var call pkg_rhino.Node = statement.GetFirstChild()
	_ = call
	var jsdoc *pkg_rhino.JSDocInfo = call.GetJSDocInfo()
	_ = jsdoc
	if (jsdoc == nil) || !(jsdoc.IsLegacyNamedTsExport()) {
		return
	}
	this_.forcedNamedExport.Add(call.GetSecondChild().GetString())
}

// Java method: ClosureExportTypeClassifier.visitNameDeclaration(com.google.javascript.rhino.Node declaration)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_visitNameDeclaration_1168(declaration pkg_rhino.Node) {
	_ = declaration
	{
		var lhs pkg_rhino.Node = declaration.GetFirstChild()
		_ = lhs
		for ; lhs != 0; lhs = lhs.GetNext() {
			if lhs.IsName() {
				this_.priv_ClosureExportTypeClassifier_visitDeclaredName_1169(lhs)
			}
		}
	}
}

// Java method: ClosureExportTypeClassifier.visitDeclaredName(com.google.javascript.rhino.Node lhs)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_visitDeclaredName_1169(lhs pkg_rhino.Node) {
	_ = lhs
	var initializer pkg_rhino.Node = lhs.GetFirstChild()
	_ = initializer
	var info *ClosureExportTypeClassifier_NameInfo = this_.priv_ClosureExportTypeClassifier_getOrCreate_1172(lhs.GetString())
	_ = info
	if (initializer != 0) && initializer.IsAssign() {
		this_.priv_ClosureExportTypeClassifier_visitExpression_1170(initializer)
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(lhs)
	_ = jsdoc
	if initializer != 0 {
		{
			r_6 := info
			r_6.value = initializer
		}
		{
			r_7 := info
			r_7.jsdoc = jsdoc
		}
	} else if (info.value == 0) && (jsdoc != nil) {
		{
			r_8 := info
			r_8.jsdoc = jsdoc
		}
	}
}

// Java method: ClosureExportTypeClassifier.visitExpression(com.google.javascript.rhino.Node expr)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_visitExpression_1170(expr pkg_rhino.Node) {
	_ = expr
	switch expr.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_ClosureExportTypeClassifier_visitAssign_1171(expr)
	case int32(26) /* Token.GETPROP */ :
		if (expr.GetJSDocInfo() != nil) && expr.GetFirstChild().IsName() {
			var variableName string = expr.GetFirstChild().GetString()
			_ = variableName
			r_9 := this_.priv_ClosureExportTypeClassifier_getOrCreate_1172(variableName)
			r_9.propertyCount = r_9.propertyCount + 1
		}
	default:
	}
}

// Java method: ClosureExportTypeClassifier.visitAssign(com.google.javascript.rhino.Node expr)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_visitAssign_1171(expr pkg_rhino.Node) {
	_ = expr
	var lhs pkg_rhino.Node = expr.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = expr.GetSecondChild()
	_ = rhs
	var jsdoc *pkg_rhino.JSDocInfo = expr.GetJSDocInfo()
	_ = jsdoc
	if lhs.IsName() {
		var assigned *ClosureExportTypeClassifier_NameInfo = this_.priv_ClosureExportTypeClassifier_getOrCreate_1172(lhs.GetString())
		_ = assigned
		{
			r_10 := assigned
			r_10.value = rhs
		}
		{
			r_11 := assigned
			r_11.jsdoc = jsdoc
		}
		{
			r_12 := assigned
			r_12.isNeverNamespace = false
		}
		{
			r_13 := assigned
			r_13.isNamedExportsLiteral = (assigned.isGoogModuleExports && NodeUtil_IsNamedExportsLiteral(rhs))
		}
		return
	}
	if lhs.IsGetProp() && lhs.GetFirstChild().IsName() {
		r_14 := this_.priv_ClosureExportTypeClassifier_getOrCreate_1172(lhs.GetFirstChild().GetString())
		r_14.propertyCount = r_14.propertyCount + 1
	}
}

// Java method: ClosureExportTypeClassifier.getOrCreate(java.lang.String name)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_getOrCreate_1172(name string) *ClosureExportTypeClassifier_NameInfo {
	_ = name
	return j2gort.OrderedMapComputeIfAbsentFn(&this_.topLevelNames, name, nil, lambdaSingleton_closure_export_type_classifier_1)

}

// Java method: ClosureExportTypeClassifier.kindOfName(com.google.javascript.jscomp.ClosureExportTypeClassifier.NameInfo info, int depth)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_kindOfName_1173(info *ClosureExportTypeClassifier_NameInfo, depth int32) ClosureExportTypeClassifier_Kind {
	_ = info
	_ = depth
	if info.isNeverNamespace {
		return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
	}
	if info.isNamedExportsLiteral {
		return StaticGet_ClosureExportTypeClassifier_Kind_NAMESPACE()
	}
	var kind ClosureExportTypeClassifier_Kind = ClosureExportTypeClassifier_Kind(0)
	_ = kind
	if (info.value != 0) || (info.jsdoc != nil) {
		kind = this_.priv_ClosureExportTypeClassifier_kindOfValue_1174(info.value, info.jsdoc, depth)
	} else {
		if info.isGoogModuleExports {
			kind = StaticGet_ClosureExportTypeClassifier_Kind_EMPTY_OBJECT()
		} else {
			kind = StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
		}
	}
	if (kind == StaticGet_ClosureExportTypeClassifier_Kind_EMPTY_OBJECT()) && (info.propertyCount > int32(0)) {
		return StaticGet_ClosureExportTypeClassifier_Kind_NAMESPACE()
	}
	return kind
}

// Java method: ClosureExportTypeClassifier.kindOfValue(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo jsdoc, int depth)
func (this_ *ClosureExportTypeClassifier) priv_ClosureExportTypeClassifier_kindOfValue_1174(value pkg_rhino.Node, jsdoc *pkg_rhino.JSDocInfo, depth int32) ClosureExportTypeClassifier_Kind {
	_ = value
	_ = jsdoc
	_ = depth
	if jsdoc != nil {
		if (jsdoc.HasEnumParameterType() || jsdoc.HasTypedefType()) || jsdoc.IsConstructorOrInterface() {
			return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
		}
		if jsdoc.HasType() {
			return ClosureExportTypeClassifier_KindOfDeclaredType(jsdoc.GetType())
		}
	}
	if value == 0 {
		return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
	}
	switch value.TokenSwitchOrd_() {
	case int32(47) /* Token.OBJECTLIT */ :
		{
			var member pkg_rhino.Node = value.GetFirstChild()
			_ = member
			for ; member != 0; member = member.GetNext() {
				if !(member.IsComputedProp()) && !(member.IsSpread()) {
					return StaticGet_ClosureExportTypeClassifier_Kind_NAMESPACE()
				}
			}
		}
		return StaticGet_ClosureExportTypeClassifier_Kind_EMPTY_OBJECT()
	case int32(32) /* Token.NAME */ :
		if depth >= int32(4) {
			return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
		}
		var aliased *ClosureExportTypeClassifier_NameInfo = this_.topLevelNames.Get(value.GetString())
		_ = aliased
		return func() ClosureExportTypeClassifier_Kind {
			if aliased == nil {
				return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
			}
			return this_.priv_ClosureExportTypeClassifier_kindOfName_1173(aliased, (depth + int32(1)))
		}()
	case int32(101) /* Token.CAST */ :
		return this_.priv_ClosureExportTypeClassifier_kindOfValue_1174(value.GetFirstChild(), value.GetJSDocInfo(), depth)
	default:
		return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
	}
	panic("unreachable switch expression")
}

// Java static method: ClosureExportTypeClassifier.kindOfDeclaredType(com.google.javascript.rhino.JSTypeExpression type)
func ClosureExportTypeClassifier_KindOfDeclaredType(type_ pkg_rhino.JSTypeExpression) ClosureExportTypeClassifier_Kind {
	_ = type_
	var root pkg_rhino.Node = type_.GetRoot()
	_ = root
	if root.RawToken_() == uint8(162) {
		return StaticGet_ClosureExportTypeClassifier_Kind_NAMESPACE()
	} else {
		return StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE()
	}
}

// Java enum: com.google.javascript.jscomp.ClosureExportTypeClassifier.Kind
type ClosureExportTypeClassifier_Kind uint8

const (
	enum_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE ClosureExportTypeClassifier_Kind = 1
	enum_ClosureExportTypeClassifier_Kind_EMPTY_OBJECT  ClosureExportTypeClassifier_Kind = 2
	enum_ClosureExportTypeClassifier_Kind_NAMESPACE     ClosureExportTypeClassifier_Kind = 3
)

var enumBaseTable_ClosureExportTypeClassifier_Kind = [...]j2gort.EnumBase{
	{Name_: "NOT_NAMESPACE", Ordinal_: 0},
	{Name_: "EMPTY_OBJECT", Ordinal_: 1},
	{Name_: "NAMESPACE", Ordinal_: 2},
}
var enumValues_ClosureExportTypeClassifier_Kind = []ClosureExportTypeClassifier_Kind{1, 2, 3}

func StaticGet_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE() ClosureExportTypeClassifier_Kind {
	return enum_ClosureExportTypeClassifier_Kind_NOT_NAMESPACE
}
func StaticGet_ClosureExportTypeClassifier_Kind_EMPTY_OBJECT() ClosureExportTypeClassifier_Kind {
	return enum_ClosureExportTypeClassifier_Kind_EMPTY_OBJECT
}
func StaticGet_ClosureExportTypeClassifier_Kind_NAMESPACE() ClosureExportTypeClassifier_Kind {
	return enum_ClosureExportTypeClassifier_Kind_NAMESPACE
}

func ClosureExportTypeClassifier_Kind_Values() []ClosureExportTypeClassifier_Kind {
	return j2gort.CloneSlice(enumValues_ClosureExportTypeClassifier_Kind)
}

func ClosureExportTypeClassifier_Kind_ValueOf(name string) ClosureExportTypeClassifier_Kind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ClosureExportTypeClassifier_Kind {
		if enumBaseTable_ClosureExportTypeClassifier_Kind[i_].Name_ == name {
			return ClosureExportTypeClassifier_Kind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ClosureExportTypeClassifier.Kind." + name))
}

func (this_ ClosureExportTypeClassifier_Kind) IsJavaNil_() bool                              { return this_ == 0 }
func (ClosureExportTypeClassifier_Kind) HasIdentityEquals_()                                 {}
func (ClosureExportTypeClassifier_Kind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ClosureExportTypeClassifier_Kind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ClosureExportTypeClassifier_Kind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ClosureExportTypeClassifier_Kind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureExportTypeClassifier$Kind", "Kind"
}
func (this_ ClosureExportTypeClassifier_Kind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ClosureExportTypeClassifier_Kind[this_-1]
}
func (this_ ClosureExportTypeClassifier_Kind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ClosureExportTypeClassifier_Kind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ClosureExportTypeClassifier_Kind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureExportTypeClassifier_Kind[this_-1].Name_
}
func (this_ ClosureExportTypeClassifier_Kind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ClosureExportTypeClassifier_Kind[this_-1].Name_
}
func (this_ ClosureExportTypeClassifier_Kind) Equals(other_ any) bool {
	o_, ok_ := other_.(ClosureExportTypeClassifier_Kind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ClosureExportTypeClassifier_Kind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ClosureExportTypeClassifier_Kind[this_-1].HashCode()
}
func (this_ ClosureExportTypeClassifier_Kind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureExportTypeClassifier_Kind[this_-1].CompareTo(other_)
}
func (this_ ClosureExportTypeClassifier_Kind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ClosureExportTypeClassifier_Kind[this_-1].Name_)
}
func init() {
	vals_ := []any{ClosureExportTypeClassifier_Kind(1), ClosureExportTypeClassifier_Kind(2), ClosureExportTypeClassifier_Kind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureExportTypeClassifier.Kind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureExportTypeClassifier$Kind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureExportTypeClassifier.Kind", "com.google.javascript.jscomp.ClosureExportTypeClassifier$Kind", func(v_ any) bool {
		return j2gort.InstanceOf[ClosureExportTypeClassifier_Kind](v_)
	})
}

// Java class: com.google.javascript.jscomp.ClosureExportTypeClassifier.NameInfo
type ClosureExportTypeClassifier_NameInfo struct {
	jsdoc                 *pkg_rhino.JSDocInfo `gson:"jsdoc"`
	value                 pkg_rhino.Node       `gson:"value"`
	propertyCount         int32                `gson:"propertyCount"`
	idHash_               atomic.Int32
	isNeverNamespace      bool `gson:"isNeverNamespace"`
	isGoogModuleExports   bool `gson:"isGoogModuleExports"`
	isNamedExportsLiteral bool `gson:"isNamedExportsLiteral"`
}

func (this_ *ClosureExportTypeClassifier_NameInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureExportTypeClassifier_NameInfo) IsJavaNil_() bool { return this_ == nil }
func (*ClosureExportTypeClassifier_NameInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureExportTypeClassifier$NameInfo", "NameInfo"
}
func (this_ *ClosureExportTypeClassifier_NameInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureExportTypeClassifier_NameInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureExportTypeClassifier.NameInfo", "com.google.javascript.jscomp.ClosureExportTypeClassifier$NameInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureExportTypeClassifier_NameInfo](v_)
	})
}

func (this_ *ClosureExportTypeClassifier_NameInfo) initFields_ClosureExportTypeClassifier_NameInfo() {
	_ = this_
}

// Java constructor: NameInfo()
func NewClosureExportTypeClassifier_NameInfo() *ClosureExportTypeClassifier_NameInfo {
	this_ := &ClosureExportTypeClassifier_NameInfo{}
	InitClosureExportTypeClassifier_NameInfo(this_)
	return this_
}

func InitClosureExportTypeClassifier_NameInfo(this_ *ClosureExportTypeClassifier_NameInfo) {
	this_.initFields_ClosureExportTypeClassifier_NameInfo()
}

// Java enum: com.google.javascript.jscomp.ClosureExportTypeClassifier.TsExportType
type ClosureExportTypeClassifier_TsExportType uint8

const (
	enum_ClosureExportTypeClassifier_TsExportType_DEFAULT ClosureExportTypeClassifier_TsExportType = 1
	enum_ClosureExportTypeClassifier_TsExportType_NAMED   ClosureExportTypeClassifier_TsExportType = 2
)

var enumBaseTable_ClosureExportTypeClassifier_TsExportType = [...]j2gort.EnumBase{
	{Name_: "DEFAULT", Ordinal_: 0},
	{Name_: "NAMED", Ordinal_: 1},
}
var enumValues_ClosureExportTypeClassifier_TsExportType = []ClosureExportTypeClassifier_TsExportType{1, 2}

func StaticGet_ClosureExportTypeClassifier_TsExportType_DEFAULT() ClosureExportTypeClassifier_TsExportType {
	return enum_ClosureExportTypeClassifier_TsExportType_DEFAULT
}
func StaticGet_ClosureExportTypeClassifier_TsExportType_NAMED() ClosureExportTypeClassifier_TsExportType {
	return enum_ClosureExportTypeClassifier_TsExportType_NAMED
}

func ClosureExportTypeClassifier_TsExportType_Values() []ClosureExportTypeClassifier_TsExportType {
	return j2gort.CloneSlice(enumValues_ClosureExportTypeClassifier_TsExportType)
}

func ClosureExportTypeClassifier_TsExportType_ValueOf(name string) ClosureExportTypeClassifier_TsExportType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ClosureExportTypeClassifier_TsExportType {
		if enumBaseTable_ClosureExportTypeClassifier_TsExportType[i_].Name_ == name {
			return ClosureExportTypeClassifier_TsExportType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ClosureExportTypeClassifier.TsExportType." + name))
}

func (this_ ClosureExportTypeClassifier_TsExportType) IsJavaNil_() bool                { return this_ == 0 }
func (ClosureExportTypeClassifier_TsExportType) HasIdentityEquals_()                   {}
func (ClosureExportTypeClassifier_TsExportType) Implements_j2gort_Ext_java_lang_Enum() {}
func (ClosureExportTypeClassifier_TsExportType) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ ClosureExportTypeClassifier_TsExportType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ClosureExportTypeClassifier_TsExportType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ClosureExportTypeClassifier$TsExportType", "TsExportType"
}
func (this_ ClosureExportTypeClassifier_TsExportType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ClosureExportTypeClassifier_TsExportType[this_-1]
}
func (this_ ClosureExportTypeClassifier_TsExportType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ClosureExportTypeClassifier_TsExportType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ClosureExportTypeClassifier_TsExportType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureExportTypeClassifier_TsExportType[this_-1].Name_
}
func (this_ ClosureExportTypeClassifier_TsExportType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ClosureExportTypeClassifier_TsExportType[this_-1].Name_
}
func (this_ ClosureExportTypeClassifier_TsExportType) Equals(other_ any) bool {
	o_, ok_ := other_.(ClosureExportTypeClassifier_TsExportType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ClosureExportTypeClassifier_TsExportType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ClosureExportTypeClassifier_TsExportType[this_-1].HashCode()
}
func (this_ ClosureExportTypeClassifier_TsExportType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ClosureExportTypeClassifier_TsExportType[this_-1].CompareTo(other_)
}
func (this_ ClosureExportTypeClassifier_TsExportType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ClosureExportTypeClassifier_TsExportType[this_-1].Name_)
}
func init() {
	vals_ := []any{ClosureExportTypeClassifier_TsExportType(1), ClosureExportTypeClassifier_TsExportType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureExportTypeClassifier.TsExportType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ClosureExportTypeClassifier$TsExportType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureExportTypeClassifier.TsExportType", "com.google.javascript.jscomp.ClosureExportTypeClassifier$TsExportType", func(v_ any) bool {
		return j2gort.InstanceOf[ClosureExportTypeClassifier_TsExportType](v_)
	})
}
