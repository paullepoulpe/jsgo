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
// Port of com/google/javascript/jscomp/JSDocInfoPrinter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.JSDocInfoPrinter
type JSDocInfoPrinter struct {
	idHash_         atomic.Int32
	useOriginalName bool `gson:"useOriginalName"`
	printDesc       bool `gson:"printDesc"`
}

func (this_ *JSDocInfoPrinter) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocInfoPrinter) IsJavaNil_() bool { return this_ == nil }
func (*JSDocInfoPrinter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSDocInfoPrinter", "JSDocInfoPrinter"
}
func (this_ *JSDocInfoPrinter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocInfoPrinter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSDocInfoPrinter", "com.google.javascript.jscomp.JSDocInfoPrinter", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocInfoPrinter](v_)
	})
}

func (this_ *JSDocInfoPrinter) initFields_JSDocInfoPrinter() {
	_ = this_
}

// Java constructor: JSDocInfoPrinter(boolean useOriginalName)
func NewJSDocInfoPrinter_1(useOriginalName bool) *JSDocInfoPrinter {
	this_ := &JSDocInfoPrinter{}
	InitJSDocInfoPrinter_1(this_, useOriginalName)
	return this_
}

func InitJSDocInfoPrinter_1(this_ *JSDocInfoPrinter, useOriginalName bool) {
	_ = useOriginalName
	InitJSDocInfoPrinter_2(this_, useOriginalName, false)
}

// Java constructor: JSDocInfoPrinter(boolean useOriginalName, boolean printDesc)
func NewJSDocInfoPrinter_2(useOriginalName bool, printDesc bool) *JSDocInfoPrinter {
	this_ := &JSDocInfoPrinter{}
	InitJSDocInfoPrinter_2(this_, useOriginalName, printDesc)
	return this_
}

func InitJSDocInfoPrinter_2(this_ *JSDocInfoPrinter, useOriginalName bool, printDesc bool) {
	_ = useOriginalName
	_ = printDesc
	this_.initFields_JSDocInfoPrinter()
	this_.useOriginalName = useOriginalName
	this_.printDesc = printDesc
}

// Java method: JSDocInfoPrinter.print(com.google.javascript.rhino.JSDocInfo info)
func (this_ *JSDocInfoPrinter) Print(info *pkg_rhino.JSDocInfo) string {
	_ = info
	var multiline bool = false
	_ = multiline
	var parts []string = []string{}
	_ = parts
	parts = append(parts, "/**")
	if info.IsExterns() {
		parts = append(parts, "@externs")
	}
	if info.IsTypeSummary() {
		parts = append(parts, "@typeSummary")
	}
	if info.IsExport() {
		parts = append(parts, "@export")
	} else if (info.GetVisibility() != 0) && (info.GetVisibility() != pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) {
		parts = append(parts, "@"+j2gort.StrValString(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.CheckNotNullEnum(info.GetVisibility()).String())))
	}
	if info.GetAuthors() != nil {
		multiline = true
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetAuthors()))
			for cur_1.HasNext() {
				name := cur_1.NextStr()
				_ = name
				parts = append(parts, "@author "+j2gort.StrValString(name))
			}
		}
	}
	if info.IsAbstract() {
		parts = append(parts, "@abstract")
	}
	if info.HasLendsName() {
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3562("lends", info.GetLendsName().GetRoot()))
	}
	if info.HasConstAnnotation() && !(info.IsDefine()) {
		parts = append(parts, "@const")
	}
	if info.IsFinal() {
		parts = append(parts, "@final")
	}
	var description string = info.GetDescription()
	_ = description
	if description != j2gort.NullString {
		multiline = true
		parts = append(parts, "@desc "+j2gort.StrValString(description))
	}
	if info.GetReferences() != nil {
		multiline = true
		{
			cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetReferences()))
			for cur_2.HasNext() {
				desc := cur_2.NextStr()
				_ = desc
				parts = append(parts, "@see "+j2gort.StrValString(desc))
			}
		}
	}
	if info.IsWizaction() {
		parts = append(parts, "@wizaction")
	}
	if info.IsNoSideEffects() {
		parts = append(parts, "@nosideeffects")
	}
	if info.IsNoCompile() {
		parts = append(parts, "@nocompile")
	}
	if info.IsNoInline() {
		parts = append(parts, "@noinline")
	}
	if info.IsRequireInlining() {
		parts = append(parts, "@requireInlining")
	}
	if info.IsEncourageInlining() {
		parts = append(parts, "@encourageInlining")
	}
	if info.IsProvideAlreadyProvided() {
		parts = append(parts, "@provideAlreadyProvided")
	}
	if info.IsIdGenerator() {
		parts = append(parts, "@idGenerator {unique}")
	}
	if info.IsConsistentIdGenerator() {
		parts = append(parts, "@idGenerator {consistent}")
	}
	if info.IsStableIdGenerator() {
		parts = append(parts, "@idGenerator {stable}")
	}
	if info.IsXidGenerator() {
		parts = append(parts, "@idGenerator {xid}")
	}
	if info.IsMappedIdGenerator() {
		parts = append(parts, "@idGenerator {mapped}")
	}
	if info.MakesDicts() {
		parts = append(parts, "@dict")
	}
	if info.MakesStructs() {
		parts = append(parts, "@struct")
	}
	if info.MakesUnrestricted() {
		parts = append(parts, "@unrestricted ")
	}
	if info.IsConstructor() {
		parts = append(parts, "@constructor")
	}
	if info.IsInterface() && !(info.UsesImplicitMatch()) {
		parts = append(parts, "@interface")
	}
	if info.IsInterface() && info.UsesImplicitMatch() {
		parts = append(parts, "@record")
	}
	if info.IsClosureUnawareCode() {
		parts = append(parts, "@closureUnaware")
	}
	if info.HasBaseType() {
		multiline = true
		var typeNode pkg_rhino.Node = this_.priv_JSDocInfoPrinter_stripBang_3559(info.GetBaseType().GetRoot())
		_ = typeNode
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3562("extends", typeNode))
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetExtendedInterfaces()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			type_ := j2gort.Cast[pkg_rhino.JSTypeExpression](rawElem_4)
			_ = type_
			multiline = true
			var typeNode_2 pkg_rhino.Node = this_.priv_JSDocInfoPrinter_stripBang_3559(type_.GetRoot())
			_ = typeNode_2
			parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3562("extends", typeNode_2))
		}
	}
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetImplementedInterfaces()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			type__2 := j2gort.Cast[pkg_rhino.JSTypeExpression](rawElem_6)
			_ = type__2
			multiline = true
			var typeNode_3 pkg_rhino.Node = this_.priv_JSDocInfoPrinter_stripBang_3559(type__2.GetRoot())
			_ = typeNode_3
			parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3562("implements", typeNode_3))
		}
	}
	if info.HasThisType() {
		multiline = true
		var typeNode_4 pkg_rhino.Node = this_.priv_JSDocInfoPrinter_stripBang_3559(info.GetThisType().GetRoot())
		_ = typeNode_4
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3562("this", typeNode_4))
	}
	var paramCount int32 = info.GetParameterCount()
	_ = paramCount
	if paramCount > int32(0) {
		multiline = true
		{
			var i int32 = int32(0)
			_ = i
			for ; i < paramCount; i = i + 1 {
				parts = append(parts, "@param "+j2gort.StrValString(this_.priv_JSDocInfoPrinter_buildParamType_3564(info, info.GetParameterNameAt(i))))
			}
		}
	}
	if info.HasReturnType() {
		multiline = true
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3561("return", info.GetReturnType(), info.GetReturnDescription()))
	}
	if !(info.GetThrowsAnnotations().IsEmpty()) && !(j2gort.CheckNotNullString(j2gort.ListGetStr(info.GetThrowsAnnotations(), int32(0))) == "") {
		multiline = true
		parts = append(parts, "@throws "+j2gort.StrValString(j2gort.ListGetStr(info.GetThrowsAnnotations(), int32(0))))
	}
	var templates j2gort.Ext_com_google_common_collect_ImmutableMap = info.GetTemplateTypes()
	_ = templates
	if !(templates.IsEmpty()) {
		multiline = true
		templates.ForEach_1(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_7 any, lp_8 any) {
			name_2 := j2gort.CastString(lp_7)
			_ = name_2
			boundExpr := j2gort.Cast[pkg_rhino.JSTypeExpression](lp_8)
			_ = boundExpr
			var boundRoot pkg_rhino.Node = boundExpr.GetRoot()
			_ = boundRoot
			if (boundRoot.RawToken_() == uint8(158)) && !(boundRoot.HasChildren()) {
				parts = append(parts, "@template "+j2gort.StrValString(name_2))
			} else {
				parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3561("template", boundExpr, name_2))
			}
		}))
	}
	var typeTransformations j2gort.Ext_com_google_common_collect_ImmutableMap = info.GetTypeTransformations()
	_ = typeTransformations
	if !(typeTransformations.IsEmpty()) {
		multiline = true
		{
			cur_9 := j2gort.NewMapEntryCursor(typeTransformations)
			for cur_9.NextEntry() {
				var name_2 string = cur_9.KeyStr()
				_ = name_2
				var tranformationDefinition string = NewCodePrinter_Builder(j2gort.CastNode[pkg_rhino.Node](cur_9.GetValue())).Build()
				_ = tranformationDefinition
				parts = append(parts, "@template "+j2gort.StrValString(name_2)+" := "+j2gort.StrValString(tranformationDefinition)+" =:")
			}
		}
	}
	if info.IsOverride() {
		parts = append(parts, "@override")
	}
	if info.HasType() && !(info.IsDefine()) {
		if info.IsInlineType() {
			parts = append(parts, this_.priv_JSDocInfoPrinter_typeNode_3565(info.GetType().GetRoot()))
		} else {
			parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3560("type", info.GetType()))
		}
	}
	if info.IsDefine() {
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3560("define", info.GetType()))
	}
	if info.HasTypedefType() {
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3560("typedef", info.GetTypedefType()))
	}
	if info.HasEnumParameterType() {
		parts = append(parts, this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3560("enum", info.GetEnumParameterType()))
	}
	if info.IsImplicitCast() {
		parts = append(parts, "@implicitCast")
	}
	if info.IsNoCollapse() {
		parts = append(parts, "@nocollapse")
	}
	var suppressions j2gort.Ext_com_google_common_collect_ImmutableMap = info.GetSuppressionsAndTheirDescription()
	_ = suppressions
	if !(suppressions.IsEmpty()) {
		{
			cur_10 := j2gort.NewMapEntryCursor(suppressions)
			for cur_10.NextEntry() {
				var warnings []string = func() []string {
					arr_14 := j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](cur_10.GetKey()).ToArray_1_slice_any(func() []any {
						arr_11 := j2gort.NewStringArray(int32(0))
						if arr_11 == nil {
							return nil
						}
						out_12 := make([]any, len(arr_11))
						for i_, el_13 := range arr_11 {
							out_12[i_] = j2gort.BoxString(el_13)
						}
						return out_12
					}())
					if arr_14 == nil {
						return nil
					}
					out_15 := make([]string, len(arr_14))
					for i_, el_16 := range arr_14 {
						out_15[i_] = j2gort.CastString(el_16)
					}
					return out_15
				}()
				_ = warnings
				func() {
					arrMut_17 := warnings
					var arrBox_18 []any
					if arrMut_17 != nil {
						arrBox_18 = make([]any, len(arrMut_17))
						for i_, el_19 := range arrMut_17 {
							arrBox_18[i_] = j2gort.BoxString(el_19)
						}
					}
					j2gort.Ext_java_util_Arrays_Sort_2(arrBox_18, j2gort.Ext_java_util_Comparator_NaturalOrder())
					for i_, el_19 := range arrBox_18 {
						arrMut_17[i_] = j2gort.CastString(el_19)
					}
				}()
				var text string = j2gort.CastString(cur_10.GetValue())
				_ = text
				var sb j2gort.LocalStringBuilder

				sb.AppendStr("@suppress {").AppendStr(j2gort.Ext_com_google_common_base_Joiner_On_1_uint16(uint16(44)).Join_1_slice_any(func() []any {
					arr_20 := warnings
					if arr_20 == nil {
						return nil
					}
					out_21 := make([]any, len(arr_20))
					for i_, el_22 := range arr_20 {
						out_21[i_] = j2gort.BoxString(el_22)
					}
					return out_21
				}())).AppendStr("}")
				if !(j2gort.CheckNotNullString(text) == "") {
					sb.AppendStr(" ").AppendStr(text)
				}
				parts = append(parts, sb.String())
			}
		}
		multiline = true
	}
	if info.IsDeprecated() {
		var deprecationReason string = info.GetDeprecationReason()
		_ = deprecationReason
		parts = append(parts, "@deprecated"+j2gort.StrValString(func() string {
			if deprecationReason != j2gort.NullString {
				return (" " + j2gort.StrValString(deprecationReason))
			}
			return ""
		}()))
		multiline = true
	}
	if info.IsPolymer() {
		multiline = true
		parts = append(parts, "@polymer")
	}
	if info.IsPolymerBehavior() {
		multiline = true
		parts = append(parts, "@polymerBehavior")
	}
	if info.IsMixinFunction() {
		multiline = true
		parts = append(parts, "@mixinFunction")
	}
	if info.IsMixinClass() {
		multiline = true
		parts = append(parts, "@mixinClass")
	}
	if info.IsCustomElement() {
		multiline = true
		parts = append(parts, "@customElement")
	}
	if info.GetClosurePrimitiveId() != j2gort.NullString {
		parts = append(parts, "@closurePrimitive {"+j2gort.StrValString(info.GetClosurePrimitiveId())+"}")
	}
	if info.IsNgInject() {
		parts = append(parts, "@ngInject")
	}
	if info.IsLegacyNamedTsExport() {
		parts = append(parts, "@legacyNamedTsExport")
	}
	{
		cur_23 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetTsTypes()))
		for cur_23.HasNext() {
			tsType := cur_23.NextStr()
			_ = tsType
			parts = append(parts, "@tsType "+j2gort.StrValString(tsType))
		}
	}
	if this_.printDesc && (info.GetBlockDescription() != j2gort.NullString) {
		var cleaned string = j2gort.JStrReplaceAll(info.GetBlockDescription(), "\n\\s*\\*\\s*", "\n")
		_ = cleaned
		if !(j2gort.CheckNotNullString(cleaned) == "") {
			multiline = true
			cleaned = javalang.StringTrim(j2gort.CheckNotNullString(cleaned))
			if int32(len(parts)) > int32(1) {
				cleaned = (j2gort.StrValString(cleaned) + "\n")
			}
			parts = j2gort.SliceInsertAt(parts, int32(1), cleaned)
		}
	}
	var sb_2 *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb_2
	if multiline {
		j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n").AppendTo_2_Ext_java_lang_StringBuilder_Ext_java_lang_Iterable(sb_2, j2gort.SliceReadOnlyIterable(parts))
	} else {
		j2gort.Ext_com_google_common_base_Joiner_On_1_string(" ").AppendTo_2_Ext_java_lang_StringBuilder_Ext_java_lang_Iterable(sb_2, j2gort.SliceReadOnlyIterable(parts))
		sb_2.Append_1_string(" */")
	}
	var s string = j2gort.JStrReplaceAll(javalang.StringReplace(sb_2.String(), "\n", "\n *"), "\n \\*([^ \n])", "\n * $1")
	_ = s
	if multiline {
		s = (j2gort.StrValString(s) + "\n */\n")
	} else {
		s = (j2gort.StrValString(s) + " ")
	}
	return s
}

// Java method: JSDocInfoPrinter.stripBang(com.google.javascript.rhino.Node typeNode)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_stripBang_3559(typeNode pkg_rhino.Node) pkg_rhino.Node {
	_ = typeNode
	if typeNode.RawToken_() == uint8(159) {
		typeNode = typeNode.GetFirstChild()
	}
	return typeNode
}

// Java method: JSDocInfoPrinter.buildAnnotationWithType(java.lang.String annotation, com.google.javascript.rhino.JSTypeExpression type)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_buildAnnotationWithType_3560(annotation string, type_ pkg_rhino.JSTypeExpression) string {
	_ = annotation
	_ = type_
	return this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3561(annotation, type_, j2gort.NullString)
}

// Java method: JSDocInfoPrinter.buildAnnotationWithType(java.lang.String annotation, com.google.javascript.rhino.JSTypeExpression type, java.lang.@org.jspecify.annotations.Nullable String description)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_buildAnnotationWithType_3561(annotation string, type_ pkg_rhino.JSTypeExpression, description string) string {
	_ = annotation
	_ = type_
	_ = description
	return this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3563(annotation, type_.GetRoot(), description)
}

// Java method: JSDocInfoPrinter.buildAnnotationWithType(java.lang.String annotation, com.google.javascript.rhino.Node type)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_buildAnnotationWithType_3562(annotation string, type_ pkg_rhino.Node) string {
	_ = annotation
	_ = type_
	return this_.priv_JSDocInfoPrinter_buildAnnotationWithType_3563(annotation, type_, j2gort.NullString)
}

// Java method: JSDocInfoPrinter.buildAnnotationWithType(java.lang.String annotation, com.google.javascript.rhino.Node type, java.lang.@org.jspecify.annotations.Nullable String description)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_buildAnnotationWithType_3563(annotation string, type_ pkg_rhino.Node, description string) string {
	_ = annotation
	_ = type_
	_ = description
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	sb.Append_1_string("@")
	sb.Append_1_string(annotation)
	sb.Append_1_string(" {")
	this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, type_)
	sb.Append_1_string("}")
	if description != j2gort.NullString {
		sb.Append_1_string(" ")
		sb.Append_1_string(description)
	}
	return sb.String()
}

// Java method: JSDocInfoPrinter.buildParamType(com.google.javascript.rhino.JSDocInfo info, java.lang.String name)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_buildParamType_3564(info *pkg_rhino.JSDocInfo, name string) string {
	_ = info
	_ = name
	var type_ pkg_rhino.JSTypeExpression = info.GetParameterType(name)
	_ = type_
	if type_ != (pkg_rhino.JSTypeExpression{}) {
		var p string = ("{" + j2gort.StrValString(this_.priv_JSDocInfoPrinter_typeNode_3565(type_.GetRoot())) + "} " + j2gort.StrValString(name) + j2gort.StrValString(func() string {
			if this_.printDesc && (info.GetDescriptionForParameter(name) != j2gort.NullString) {
				return info.GetDescriptionForParameter(name)
			}
			return ""
		}()))
		_ = p
		return javalang.StringTrim(j2gort.CheckNotNullString(p))
	} else {
		return name
	}
}

// Java method: JSDocInfoPrinter.typeNode(com.google.javascript.rhino.Node typeNode)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_typeNode_3565(typeNode pkg_rhino.Node) string {
	_ = typeNode
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode)
	return sb.String()
}

// Java method: JSDocInfoPrinter.appendTypeNode(java.lang.StringBuilder sb, com.google.javascript.rhino.Node typeNode)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_appendTypeNode_3566(sb *j2gort.Ext_java_lang_StringBuilder, typeNode pkg_rhino.Node) {
	_ = sb
	_ = typeNode
	if this_.useOriginalName && (typeNode.GetOriginalName() != j2gort.NullString) {
		sb.Append_1_string(typeNode.GetOriginalName())
		return
	}
	switch typeNode.TokenSwitchOrd_() {
	case int32(159) /* Token.BANG */ :
		sb.Append_1_string("!")
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode.GetFirstChild())
	case int32(160) /* Token.EQUALS */ :
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode.GetFirstChild())
		sb.Append_1_string("=")
	case int32(155) /* Token.PIPE */ :
		sb.Append_1_string("(")
		var lastChild pkg_rhino.Node = typeNode.GetLastChild()
		_ = lastChild
		{
			var child pkg_rhino.Node = typeNode.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, child)
				if child != lastChild {
					sb.Append_1_string("|")
				}
			}
		}
		sb.Append_1_string(")")
	case int32(125) /* Token.ITER_REST */ :
		sb.Append_1_string("...")
		if typeNode.HasChildren() && !(typeNode.GetFirstChild().IsEmpty()) {
			this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode.GetFirstChild())
		}
	case int32(156) /* Token.STAR */ :
		sb.Append_1_string("*")
	case int32(158) /* Token.QMARK */ :
		sb.Append_1_string("?")
		if typeNode.HasChildren() {
			this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode.GetFirstChild())
		}
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_JSDocInfoPrinter_appendFunctionNode_3567(sb, typeNode)
	case int32(162) /* Token.LC */ :
		sb.Append_1_string("{")
		var lb pkg_rhino.Node = typeNode.GetFirstChild()
		_ = lb
		var lastColon pkg_rhino.Node = lb.GetLastChild()
		_ = lastColon
		{
			var colon pkg_rhino.Node = lb.GetFirstChild()
			_ = colon
			for ; colon != 0; colon = colon.GetNext() {
				if colon.HasChildren() {
					j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(colon.GetFirstChild().GetString())).Append_1_string(":")
					this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, colon.GetLastChild())
				} else {
					sb.Append_1_string(colon.GetString())
				}
				if colon != lastColon {
					sb.Append_1_string(",")
				}
			}
		}
		sb.Append_1_string("}")
	case int32(87) /* Token.VOID */ :
		sb.Append_1_string("void")
	case int32(25) /* Token.TYPEOF */ :
		sb.Append_1_string("typeof ")
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode.GetFirstChild())
	case int32(90) /* Token.BLOCK */ :
		sb.Append_1_string("<")
		var last pkg_rhino.Node = typeNode.GetLastChild()
		_ = last
		{
			var type_ pkg_rhino.Node = typeNode.GetFirstChild()
			_ = type_
			for ; type_ != 0; type_ = type_.GetNext() {
				this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, type_)
				if type_ != last {
					sb.Append_1_string(",")
				}
			}
		}
		sb.Append_1_string(">")
	case int32(35) /* Token.STRINGLIT */ :
		sb.Append_1_string(typeNode.GetString())
		if typeNode.HasChildren() {
			this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, typeNode.GetOnlyChild())
		}
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected typeNode: "+j2gort.StrVal(typeNode))))
	}
}

// Java method: JSDocInfoPrinter.appendFunctionNode(java.lang.StringBuilder sb, com.google.javascript.rhino.Node function)
func (this_ *JSDocInfoPrinter) priv_JSDocInfoPrinter_appendFunctionNode_3567(sb *j2gort.Ext_java_lang_StringBuilder, function pkg_rhino.Node) {
	_ = sb
	_ = function
	var hasNewOrThis bool = false
	_ = hasNewOrThis
	sb.Append_1_string("function(")
	var first pkg_rhino.Node = function.GetFirstChild()
	_ = first
	if first.IsNew() {
		sb.Append_1_string("new:")
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, first.GetFirstChild())
		hasNewOrThis = true
	} else if first.IsThis() {
		sb.Append_1_string("this:")
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, first.GetFirstChild())
		hasNewOrThis = true
	} else if first.IsEmpty() {
		sb.Append_1_string(")")
		return
	} else if !(first.IsParamList()) {
		sb.Append_1_string("):")
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, first)
		return
	}
	var paramList pkg_rhino.Node = pkg_rhino.Node(0)
	_ = paramList
	if first.IsParamList() {
		paramList = first
	} else if first.GetNext().IsParamList() {
		paramList = first.GetNext()
	}
	if paramList != 0 {
		var firstParam bool = true
		_ = firstParam
		{
			var param pkg_rhino.Node = paramList.GetFirstChild()
			_ = param
			for ; param != 0; param = param.GetNext() {
				if !(firstParam) || hasNewOrThis {
					sb.Append_1_string(",")
				}
				this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, param)
				firstParam = false
			}
		}
	}
	sb.Append_1_string(")")
	var returnType pkg_rhino.Node = function.GetLastChild()
	_ = returnType
	if !(returnType.IsEmpty()) {
		sb.Append_1_string(":")
		this_.priv_JSDocInfoPrinter_appendTypeNode_3566(sb, returnType)
	}
}
