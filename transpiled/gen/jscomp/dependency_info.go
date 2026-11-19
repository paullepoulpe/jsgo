/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/DependencyInfo.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_dependency_info_1 j2gort.Ext_com_google_common_base_Function

func init() {
	lambdaSingleton_dependency_info_1 = j2gort.NewLambda_Ext_com_google_common_base_Function(func(lp_11 any) any {
		arg := j2gort.CastString(lp_11)
		_ = arg
		return "'" + j2gort.StrValString(javalang.StringReplace(arg, "'", "\\'")) + "'"
	})
}

// Java interface: com.google.javascript.jscomp.deps.DependencyInfo
type DependencyInfo interface {
	Implements_jscomp_DependencyInfo()
	IsJavaNil_() bool
	GetName() string
	GetPathRelativeToClosureBase() string
	GetProvides() j2gort.Ext_com_google_common_collect_ImmutableList
	GetRequires() j2gort.Ext_com_google_common_collect_ImmutableList
	GetRequiredSymbols() j2gort.Ext_com_google_common_collect_ImmutableList
	GetTypeRequires() j2gort.Ext_com_google_common_collect_ImmutableList
	GetLoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap
	IsEs6Module() bool
	IsGoogModule() bool
	GetHasExternsAnnotation() bool
	GetHasNoCompileAnnotation() bool
}

func ToIface_DependencyInfo[T interface {
	DependencyInfo
	comparable
}](v T) DependencyInfo {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DependencyInfo", "com.google.javascript.jscomp.deps.DependencyInfo", func(v_ any) bool {
		return j2gort.InstanceOf[DependencyInfo](v_)
	})
}

// Java default interface method: DependencyInfo.getRequiredSymbols()
func IfaceDefault_DependencyInfo_GetRequiredSymbols(this_ DependencyInfo) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = this_
	return DependencyInfo_Require_AsSymbolList(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetRequires()))
}

// Java default interface method: DependencyInfo.isEs6Module()
func IfaceDefault_DependencyInfo_IsEs6Module(this_ DependencyInfo) bool {
	_ = this_
	return (j2gort.CastString(j2gort.MapGetStr(this_.GetLoadFlags(), "module")) == "es6")
}

// Java default interface method: DependencyInfo.isGoogModule()
func IfaceDefault_DependencyInfo_IsGoogModule(this_ DependencyInfo) bool {
	_ = this_
	return (j2gort.CastString(j2gort.MapGetStr(this_.GetLoadFlags(), "module")) == "goog")
}

// Java class: com.google.javascript.jscomp.deps.DependencyInfo.Require
type DependencyInfo_Require struct {
	symbol  string                      `gson:"symbol"`
	rawText string                      `gson:"rawText"`
	type_   DependencyInfo_Require_Type `gson:"type"`
}

func (this_ *DependencyInfo_Require) BindSelf_(s any) {
	_ = s
}
func (this_ *DependencyInfo_Require) IsJavaNil_() bool { return this_ == nil }
func (*DependencyInfo_Require) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DependencyInfo$Require", "Require"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DependencyInfo.Require", "com.google.javascript.jscomp.deps.DependencyInfo$Require", func(v_ any) bool {
		return j2gort.InstanceOf[*DependencyInfo_Require](v_)
	})
}

var static_DependencyInfo_Require_BASE *DependencyInfo_Require

func StaticGet_DependencyInfo_Require_BASE() *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	return static_DependencyInfo_Require_BASE
}
func StaticSet_DependencyInfo_Require_BASE(v *DependencyInfo_Require) {
	clinit_DependencyInfo_Require()
	static_DependencyInfo_Require_BASE = v
}

var clinitDone_DependencyInfo_Require uint32
var clinitState_DependencyInfo_Require j2gort.ClinitState

func clinit_DependencyInfo_Require() {
	if atomic.LoadUint32(&clinitDone_DependencyInfo_Require) != 0 {
		return
	}
	clinitSlow_DependencyInfo_Require()
}

func clinitSlow_DependencyInfo_Require() {
	if !j2gort.EnterClinit(&clinitState_DependencyInfo_Require, &clinitDone_DependencyInfo_Require) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DependencyInfo_Require, &clinitDone_DependencyInfo_Require)
	static_DependencyInfo_Require_BASE = DependencyInfo_Require_GoogRequireSymbol("goog")
}

func (this_ *DependencyInfo_Require) initFields_DependencyInfo_Require() {
	_ = this_
	this_.symbol = j2gort.NullString
	this_.rawText = j2gort.NullString
}

// Java constructor: Require(java.lang.String symbol, java.lang.String rawText, com.google.javascript.jscomp.deps.DependencyInfo.Require.Type type)
func NewDependencyInfo_Require(symbol string, rawText string, type_ DependencyInfo_Require_Type) *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	this_ := &DependencyInfo_Require{}
	InitDependencyInfo_Require(this_, symbol, rawText, type_)
	return this_
}

func InitDependencyInfo_Require(this_ *DependencyInfo_Require, symbol string, rawText string, type_ DependencyInfo_Require_Type) {
	_ = symbol
	_ = rawText
	_ = type_
	this_.initFields_DependencyInfo_Require()
	if symbol == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null symbol")))
	}
	if rawText == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null rawText")))
	}
	if type_ == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null type")))
	}
	this_.symbol = symbol
	this_.rawText = rawText
	this_.type_ = type_
}

// Java static method: Require.asSymbolList(java.lang.Iterable<com.google.javascript.jscomp.deps.DependencyInfo.Require> requires)
func DependencyInfo_Require_AsSymbolList(requires j2gort.Ext_java_lang_Iterable) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_DependencyInfo_Require()
	_ = requires
	var list j2gort.Ext_java_util_List
	_ = list
	var ok_1 bool
	list, ok_1 = j2gort.TryCast[j2gort.Ext_java_util_List](requires)
	if ok_1 {
		var size int32 = list.Size()
		_ = size
		if size == int32(0) {
			return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
		}
		if size == int32(1) {
			return j2gort.ImmListOfStr1(j2gort.CastPtr[DependencyInfo_Require](list.Get(int32(0))).symbol)
		}
		if size == int32(2) {
			return j2gort.ImmListOfStr2(j2gort.CastPtr[DependencyInfo_Require](list.Get(int32(0))).symbol, j2gort.CastPtr[DependencyInfo_Require](list.Get(int32(1))).symbol)
		}
		if size == int32(3) {
			return j2gort.ImmListOfStr3(j2gort.CastPtr[DependencyInfo_Require](list.Get(int32(0))).symbol, j2gort.CastPtr[DependencyInfo_Require](list.Get(int32(1))).symbol, j2gort.CastPtr[DependencyInfo_Require](list.Get(int32(2))).symbol)
		}
		var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, size)
		_ = builder
		{
			var i int32 = int32(0)
			_ = i
			for ; i < size; i = i + 1 {
				builder.AddStrVoid_(j2gort.CastPtr[DependencyInfo_Require](list.Get(i)).symbol)
			}
		}
		return builder.BuildImmList_()
	}
	var builder_2 []string = []string{}
	_ = builder_2
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(requires))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			require := j2gort.CastPtr[DependencyInfo_Require](rawElem_3)
			_ = require
			builder_2 = append(builder_2, j2gort.CheckNotNullString(require.symbol))
		}
	}
	return j2gort.ImmListFromTypedSlice(builder_2)
}

// Java static method: Require.googRequireSymbol(java.lang.String symbol)
func DependencyInfo_Require_GoogRequireSymbol(symbol string) *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	_ = symbol
	return NewDependencyInfo_Require(symbol, symbol, StaticGet_DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL())
}

// Java static method: Require.es6Import(java.lang.String symbol, java.lang.String rawPath)
func DependencyInfo_Require_Es6Import(symbol string, rawPath string) *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	_ = symbol
	_ = rawPath
	return NewDependencyInfo_Require(symbol, rawPath, StaticGet_DependencyInfo_Require_Type_ES6_IMPORT())
}

// Java static method: Require.commonJs(java.lang.String symbol, java.lang.String rawPath)
func DependencyInfo_Require_CommonJs(symbol string, rawPath string) *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	_ = symbol
	_ = rawPath
	return NewDependencyInfo_Require(symbol, rawPath, StaticGet_DependencyInfo_Require_Type_COMMON_JS())
}

// Java static method: Require.compilerChunk(java.lang.String symbol)
func DependencyInfo_Require_CompilerChunk(symbol string) *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	_ = symbol
	return NewDependencyInfo_Require(symbol, symbol, StaticGet_DependencyInfo_Require_Type_COMPILER_CHUNK())
}

// Java static method: Require.parsedFromDeps(java.lang.String symbol)
func DependencyInfo_Require_ParsedFromDeps(symbol string) *DependencyInfo_Require {
	clinit_DependencyInfo_Require()
	_ = symbol
	return NewDependencyInfo_Require(symbol, symbol, StaticGet_DependencyInfo_Require_Type_PARSED_FROM_DEPS())
}

// Java method: Require.withSymbol(java.lang.String symbol)
func (this_ *DependencyInfo_Require) WithSymbol(symbol string) *DependencyInfo_Require {
	_ = symbol
	return NewDependencyInfo_Require(symbol, this_.rawText, this_.type_)
}

// Java method: Require.getSymbol()
func (this_ *DependencyInfo_Require) GetSymbol() string {
	return this_.symbol
}

// Java method: Require.getRawText()
func (this_ *DependencyInfo_Require) GetRawText() string {
	return this_.rawText
}

// Java method: Require.getType()
func (this_ *DependencyInfo_Require) GetType() DependencyInfo_Require_Type {
	return this_.type_
}

// Java method: Require.toString()
func (this_ *DependencyInfo_Require) String() string {
	return ("Require{symbol=" + j2gort.StrValString(this_.symbol) + ", rawText=" + j2gort.StrValString(this_.rawText) + ", type=" + j2gort.StrVal(this_.type_) + "}")
}

// Java method: Require.equals(java.lang.Object o)
func (this_ *DependencyInfo_Require) Equals(arg_4 any) bool {
	o := arg_4
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	var that *DependencyInfo_Require
	_ = that
	var ok_5 bool
	that, ok_5 = j2gort.TryCast[*DependencyInfo_Require](o)
	if ok_5 {
		return ((func() bool { sr_6 := this_.symbol; se_7 := that.symbol; return j2gort.CheckNotNullString(sr_6) == se_7 }() && func() bool {
			sr_8 := this_.rawText
			se_9 := that.rawText
			return j2gort.CheckNotNullString(sr_8) == se_9
		}()) && (this_.type_ == that.type_))
	}
	return false
}

// Java method: Require.hashCode()
func (this_ *DependencyInfo_Require) HashCode() int32 {
	var h int32 = int32(1)
	_ = h
	h = (h * int32(1000003))
	h = (h ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.symbol)))
	h = (h * int32(1000003))
	h = (h ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.rawText)))
	h = (h * int32(1000003))
	h = (h ^ j2gort.CheckNotNullEnum(this_.type_).HashCode())
	return h
}

// Java enum: com.google.javascript.jscomp.deps.DependencyInfo.Require.Type
type DependencyInfo_Require_Type uint8

const (
	enum_DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL DependencyInfo_Require_Type = 1
	enum_DependencyInfo_Require_Type_ES6_IMPORT          DependencyInfo_Require_Type = 2
	enum_DependencyInfo_Require_Type_PARSED_FROM_DEPS    DependencyInfo_Require_Type = 3
	enum_DependencyInfo_Require_Type_COMMON_JS           DependencyInfo_Require_Type = 4
	enum_DependencyInfo_Require_Type_COMPILER_CHUNK      DependencyInfo_Require_Type = 5
)

var enumBaseTable_DependencyInfo_Require_Type = [...]j2gort.EnumBase{
	{Name_: "GOOG_REQUIRE_SYMBOL", Ordinal_: 0},
	{Name_: "ES6_IMPORT", Ordinal_: 1},
	{Name_: "PARSED_FROM_DEPS", Ordinal_: 2},
	{Name_: "COMMON_JS", Ordinal_: 3},
	{Name_: "COMPILER_CHUNK", Ordinal_: 4},
}
var enumValues_DependencyInfo_Require_Type = []DependencyInfo_Require_Type{1, 2, 3, 4, 5}

func StaticGet_DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL() DependencyInfo_Require_Type {
	return enum_DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL
}
func StaticGet_DependencyInfo_Require_Type_ES6_IMPORT() DependencyInfo_Require_Type {
	return enum_DependencyInfo_Require_Type_ES6_IMPORT
}
func StaticGet_DependencyInfo_Require_Type_PARSED_FROM_DEPS() DependencyInfo_Require_Type {
	return enum_DependencyInfo_Require_Type_PARSED_FROM_DEPS
}
func StaticGet_DependencyInfo_Require_Type_COMMON_JS() DependencyInfo_Require_Type {
	return enum_DependencyInfo_Require_Type_COMMON_JS
}
func StaticGet_DependencyInfo_Require_Type_COMPILER_CHUNK() DependencyInfo_Require_Type {
	return enum_DependencyInfo_Require_Type_COMPILER_CHUNK
}

func DependencyInfo_Require_Type_Values() []DependencyInfo_Require_Type {
	return j2gort.CloneSlice(enumValues_DependencyInfo_Require_Type)
}

func DependencyInfo_Require_Type_ValueOf(name string) DependencyInfo_Require_Type {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_DependencyInfo_Require_Type {
		if enumBaseTable_DependencyInfo_Require_Type[i_].Name_ == name {
			return DependencyInfo_Require_Type(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.DependencyInfo.Require.Type." + name))
}

func (this_ DependencyInfo_Require_Type) IsJavaNil_() bool                              { return this_ == 0 }
func (DependencyInfo_Require_Type) HasIdentityEquals_()                                 {}
func (DependencyInfo_Require_Type) Implements_j2gort_Ext_java_lang_Enum()               {}
func (DependencyInfo_Require_Type) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ DependencyInfo_Require_Type) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (DependencyInfo_Require_Type) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DependencyInfo$Require$Type", "Type"
}
func (this_ DependencyInfo_Require_Type) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_DependencyInfo_Require_Type[this_-1]
}
func (this_ DependencyInfo_Require_Type) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ DependencyInfo_Require_Type) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ DependencyInfo_Require_Type) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DependencyInfo_Require_Type[this_-1].Name_
}
func (this_ DependencyInfo_Require_Type) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_DependencyInfo_Require_Type[this_-1].Name_
}
func (this_ DependencyInfo_Require_Type) Equals(other_ any) bool {
	o_, ok_ := other_.(DependencyInfo_Require_Type)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ DependencyInfo_Require_Type) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_DependencyInfo_Require_Type[this_-1].HashCode()
}
func (this_ DependencyInfo_Require_Type) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DependencyInfo_Require_Type[this_-1].CompareTo(other_)
}
func (this_ DependencyInfo_Require_Type) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_DependencyInfo_Require_Type[this_-1].Name_)
}
func init() {
	vals_ := []any{DependencyInfo_Require_Type(1), DependencyInfo_Require_Type(2), DependencyInfo_Require_Type(3), DependencyInfo_Require_Type(4), DependencyInfo_Require_Type(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.DependencyInfo.Require.Type", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.DependencyInfo$Require$Type", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DependencyInfo.Require.Type", "com.google.javascript.jscomp.deps.DependencyInfo$Require$Type", func(v_ any) bool {
		return j2gort.InstanceOf[DependencyInfo_Require_Type](v_)
	})
}

// Java class: com.google.javascript.jscomp.deps.DependencyInfo.Util
type DependencyInfo_Util struct {
	idHash_ atomic.Int32
}

func (this_ *DependencyInfo_Util) BindSelf_(s any) {
	_ = s
}
func (this_ *DependencyInfo_Util) IsJavaNil_() bool { return this_ == nil }
func (*DependencyInfo_Util) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DependencyInfo$Util", "Util"
}
func (this_ *DependencyInfo_Util) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DependencyInfo_Util) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DependencyInfo.Util", "com.google.javascript.jscomp.deps.DependencyInfo$Util", func(v_ any) bool {
		return j2gort.InstanceOf[*DependencyInfo_Util](v_)
	})
}

func (this_ *DependencyInfo_Util) initFields_DependencyInfo_Util() {
	_ = this_
}

// Java constructor: Util()
func NewDependencyInfo_Util() *DependencyInfo_Util {
	this_ := &DependencyInfo_Util{}
	InitDependencyInfo_Util(this_)
	return this_
}

func InitDependencyInfo_Util(this_ *DependencyInfo_Util) {
	this_.initFields_DependencyInfo_Util()
}

// Java static method: Util.writeAddDependency(java.lang.Appendable out, com.google.javascript.jscomp.deps.DependencyInfo info)
func DependencyInfo_Util_WriteAddDependency(out j2gort.Ext_java_lang_Appendable, info DependencyInfo) {
	_ = out
	_ = info
	j2gort.AsIface[j2gort.Ext_java_lang_Appendable](j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out.Append_1_any("goog.addDependency('")).Append_1_any(j2gort.BoxString(info.GetPathRelativeToClosureBase()))).Append_1_any("', ")
	DependencyInfo_Util_WriteJsArray(out, j2gort.AsIface[j2gort.Ext_java_util_Collection](info.GetProvides()))
	out.Append_1_any(", ")
	DependencyInfo_Util_WriteJsArray(out, j2gort.AsIface[j2gort.Ext_java_util_Collection](DependencyInfo_Require_AsSymbolList(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](info.GetRequires()))))
	var loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap = info.GetLoadFlags()
	_ = loadFlags
	if !(loadFlags.IsEmpty()) {
		out.Append_1_any(", ")
		DependencyInfo_Util_WriteJsObject(out, j2gort.AsIface[j2gort.Ext_java_util_Map](loadFlags))
	}
	out.Append_1_any(");\n")
}

// Java static method: Util.writeJsObject(java.lang.Appendable out, java.util.Map<java.lang.String,java.lang.String> map)
func DependencyInfo_Util_WriteJsObject(out j2gort.Ext_java_lang_Appendable, map_ j2gort.Ext_java_util_Map) {
	_ = out
	_ = map_
	var entries []string = []string{}
	_ = entries
	{
		cur_10 := j2gort.NewMapEntryCursor(map_)
		for cur_10.NextEntry() {
			var key string = javalang.StringReplace(cur_10.KeyStr(), "'", "\\'")
			_ = key
			var value string = javalang.StringReplace(j2gort.CastString(cur_10.GetValue()), "'", "\\'")
			_ = value
			entries = append(entries, "'"+j2gort.StrValString(key)+"': '"+j2gort.StrValString(value)+"'")
		}
	}
	out.Append_1_any("{")
	out.Append_1_any(j2gort.BoxString(j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(entries))))
	out.Append_1_any("}")
}

// Java static method: Util.writeJsArray(java.lang.Appendable out, java.util.Collection<java.lang.String> values)
func DependencyInfo_Util_WriteJsArray(out j2gort.Ext_java_lang_Appendable, values j2gort.Ext_java_util_Collection) {
	_ = out
	_ = values
	var quoted j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Transform(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](values), lambdaSingleton_dependency_info_1)
	_ = quoted
	out.Append_1_any("[")
	out.Append_1_any(j2gort.BoxString(j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_Ext_java_lang_Iterable(quoted)))
	out.Append_1_any("]")
}
