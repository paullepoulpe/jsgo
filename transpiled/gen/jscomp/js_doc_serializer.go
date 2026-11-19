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
// Port of com/google/javascript/jscomp/serialization/JSDocSerializer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.JSDocSerializer
type JSDocSerializer struct {
	idHash_ atomic.Int32
}

func (this_ *JSDocSerializer) BindSelf_(s any) {
	_ = s
}
func (this_ *JSDocSerializer) IsJavaNil_() bool { return this_ == nil }
func (*JSDocSerializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.JSDocSerializer", "JSDocSerializer"
}
func (this_ *JSDocSerializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSDocSerializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.JSDocSerializer", "com.google.javascript.jscomp.serialization.JSDocSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[*JSDocSerializer](v_)
	})
}

var static_JSDocSerializer_SYNTHETIC_SOURCE *SourceFile

func StaticGet_JSDocSerializer_SYNTHETIC_SOURCE() *SourceFile {
	clinit_JSDocSerializer()
	return static_JSDocSerializer_SYNTHETIC_SOURCE
}
func StaticSet_JSDocSerializer_SYNTHETIC_SOURCE(v *SourceFile) {
	clinit_JSDocSerializer()
	static_JSDocSerializer_SYNTHETIC_SOURCE = v
}
func StaticGet_JSDocSerializer_PLACEHOLDER_TYPE_NAME() string {
	return "JSDocSerializer_placeholder_type"
}

var static_JSDocSerializer_placeholderType pkg_rhino.JSTypeExpression

func StaticGet_JSDocSerializer_PlaceholderType() pkg_rhino.JSTypeExpression {
	clinit_JSDocSerializer()
	return static_JSDocSerializer_placeholderType
}
func StaticSet_JSDocSerializer_PlaceholderType(v pkg_rhino.JSTypeExpression) {
	clinit_JSDocSerializer()
	static_JSDocSerializer_placeholderType = v
}

var clinitDone_JSDocSerializer uint32
var clinitState_JSDocSerializer j2gort.ClinitState

func clinit_JSDocSerializer() {
	if atomic.LoadUint32(&clinitDone_JSDocSerializer) != 0 {
		return
	}
	clinitSlow_JSDocSerializer()
}

func clinitSlow_JSDocSerializer() {
	if !j2gort.EnterClinit(&clinitState_JSDocSerializer, &clinitDone_JSDocSerializer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSDocSerializer, &clinitDone_JSDocSerializer)
	static_JSDocSerializer_SYNTHETIC_SOURCE = SourceFile_FromCode_2("JSDocSerializer_placeholder_source", "")
	static_JSDocSerializer_placeholderType = JSDocSerializer_CreatePlaceholderType()
}

func (this_ *JSDocSerializer) initFields_JSDocSerializer() {
	_ = this_
}

// Java constructor: JSDocSerializer()
func NewJSDocSerializer() *JSDocSerializer {
	clinit_JSDocSerializer()
	this_ := &JSDocSerializer{}
	InitJSDocSerializer(this_)
	return this_
}

func InitJSDocSerializer(this_ *JSDocSerializer) {
	this_.initFields_JSDocSerializer()
}

// Java static method: JSDocSerializer.convertJSDocInfoForOptimizations(com.google.javascript.rhino.JSDocInfo jsdoc)
func JSDocSerializer_ConvertJSDocInfoForOptimizations(jsdoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JSDocSerializer()
	_ = jsdoc
	var stringPool *StringPool_Builder = StringPool_S_Builder()
	_ = stringPool
	return JSDocSerializer_DeserializeJsdoc(JSDocSerializer_SerializeJsdoc(jsdoc, stringPool), stringPool.Build())
}

// Java static method: JSDocSerializer.addKind(com.google.javascript.jscomp.serialization.OptimizationJsdoc.@org.jspecify.annotations.Nullable Builder b, com.google.javascript.jscomp.serialization.JsdocTag tag)
func JSDocSerializer_AddKind(b *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder, tag *j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag) *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder {
	clinit_JSDocSerializer()
	_ = b
	_ = tag
	if b == nil {
		b = j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_NewBuilder_0()
	}
	b.AddKind(tag)
	return b
}

// Java static method: JSDocSerializer.ensureBuilder(com.google.javascript.jscomp.serialization.OptimizationJsdoc.@org.jspecify.annotations.Nullable Builder b)
func JSDocSerializer_EnsureBuilder(b *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder {
	clinit_JSDocSerializer()
	_ = b
	if b != nil {
		return b
	} else {
		return j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_NewBuilder_0()
	}
}

// Java static method: JSDocSerializer.serializeJsdoc(com.google.javascript.rhino.JSDocInfo jsdoc, com.google.javascript.jscomp.serialization.StringPool.Builder stringPool)
func JSDocSerializer_SerializeJsdoc(jsdoc *pkg_rhino.JSDocInfo, stringPool *StringPool_Builder) *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc {
	clinit_JSDocSerializer()
	_ = jsdoc
	_ = stringPool
	if jsdoc == nil {
		return nil
	}
	var builder *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder = nil
	_ = builder
	if jsdoc.HasFileOverview() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_FILEOVERVIEW)
	}
	if jsdoc.GetLicense() != j2gort.NullString {
		builder = JSDocSerializer_EnsureBuilder(builder)
		builder.SetLicenseTextPointer(stringPool.Put(jsdoc.GetLicense()))
	}
	if jsdoc.IsSassGeneratedCssTs() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_SASS_GENERATED_CSS_TS)
	}
	if jsdoc.IsNoCoverage() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_NO_COVERAGE)
	}
	if jsdoc.IsNoInline() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_NO_INLINE)
	}
	if jsdoc.IsEncourageInlining() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ENCOURAGE_INLINING)
	}
	if jsdoc.IsRequireInlining() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_REQUIRE_INLINING)
	}
	if jsdoc.IsNoCollapse() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_NO_COLLAPSE)
	}
	if jsdoc.IsProvideGoog() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_PROVIDE_GOOG)
	}
	if jsdoc.IsProvideAlreadyProvided() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_PROVIDE_ALREADY_PROVIDED)
	}
	if jsdoc.IsTypeSummary() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_TYPE_SUMMARY_FILE)
	}
	if jsdoc.IsPureOrBreakMyCode() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_PURE_OR_BREAK_MY_CODE)
	}
	if jsdoc.IsCollapsibleOrBreakMyCode() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE)
	}
	if jsdoc.HasThisType() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_THIS)
	}
	if jsdoc.HasEnumParameterType() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ENUM)
	}
	if jsdoc.IsDefine() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_DEFINE)
	}
	if jsdoc.HasConstAnnotation() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_CONST)
	}
	if jsdoc.IsAnyIdGenerator() {
		builder = JSDocSerializer_AddKind(builder, JSDocSerializer_SerializeIdGenerator(jsdoc))
	}
	if jsdoc.IsNoSideEffects() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_NO_SIDE_EFFECTS)
	}
	if jsdoc.HasModifies() {
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(jsdoc.GetModifies()))
			for cur_1.HasNext() {
				modifies := cur_1.NextStr()
				_ = modifies
				switch j2gort.CheckNotNullString(modifies) {
				case "this":
					builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_MODIFIES_THIS)
				case "arguments":
					builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_MODIFIES_ARGUMENTS)
				default:
					builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_MODIFIES_ARGUMENTS)
				}
			}
		}
	}
	if !(jsdoc.GetThrowsAnnotations().IsEmpty()) {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_THROWS)
	}
	if jsdoc.IsConstructor() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_CONSTRUCTOR)
	}
	if jsdoc.IsInterface() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_INTERFACE)
	}
	var suppressions j2gort.Ext_java_util_Set = jsdoc.GetSuppressions()
	_ = suppressions
	if !(suppressions.IsEmpty()) {
		if j2gort.SetContainsStr(suppressions, "partialAlias") {
			builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_SUPPRESS_PARTIAL_ALIAS)
		}
	}
	if jsdoc.IsAbstract() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ABSTRACT)
	}
	if jsdoc.GetDescription() != j2gort.NullString {
		builder = JSDocSerializer_EnsureBuilder(builder)
		builder.SetDescriptionPointer(stringPool.Put(jsdoc.GetDescription()))
	}
	if jsdoc.GetAlternateMessageId() != j2gort.NullString {
		builder = JSDocSerializer_EnsureBuilder(builder)
		builder.SetAlternateMessageIdPointer(stringPool.Put(jsdoc.GetAlternateMessageId()))
	}
	if jsdoc.GetMeaning() != j2gort.NullString {
		builder = JSDocSerializer_EnsureBuilder(builder)
		builder.SetMeaningPointer(stringPool.Put(jsdoc.GetMeaning()))
	}
	if !(suppressions.IsEmpty()) {
		if j2gort.SetContainsStr(suppressions, "messageConventions") {
			builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_SUPPRESS_MESSAGE_CONVENTION)
		}
		if j2gort.SetContainsStr(suppressions, "untranspilableFeatures") {
			builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES)
		}
	}
	if jsdoc.IsUsedViaDotConstructor() {
		builder = JSDocSerializer_AddKind(builder, j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_USED_VIA_DOT_CONSTRUCTOR)
	}
	if jsdoc.GetPerFileClosureUnawareMode() != 0 {
		builder = JSDocSerializer_EnsureBuilder(builder)
		builder.SetClosureUnawareModePointer(stringPool.Put((jsdoc.GetPerFileClosureUnawareMode()).Name()))
	}
	if builder == nil {
		return nil
	}
	var result *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc](builder.Build())
	_ = result
	if j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_GetDefaultInstance().Equals(j2gort.PtrToAny(result)) {
		return nil
	}
	return j2gort.CheckNotNullPtr(result)
}

// Java static method: JSDocSerializer.serializeIdGenerator(com.google.javascript.rhino.JSDocInfo doc)
func JSDocSerializer_SerializeIdGenerator(doc *pkg_rhino.JSDocInfo) *j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag {
	clinit_JSDocSerializer()
	_ = doc
	if doc.IsConsistentIdGenerator() {
		return j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ID_GENERATOR_CONSISTENT
	} else if doc.IsStableIdGenerator() {
		return j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ID_GENERATOR_STABLE
	} else if doc.IsXidGenerator() {
		return j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ID_GENERATOR_XID
	} else if doc.IsMappedIdGenerator() {
		return j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ID_GENERATOR_MAPPED
	} else if doc.IsIdGenerator() {
		return j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag_F_JSDOC_ID_GENERATOR_INCONSISTENT
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Failed to identify idGenerator inside JSDoc: "+j2gort.StrVal(doc))))
}

// Java static method: JSDocSerializer.createPlaceholderType()
func JSDocSerializer_CreatePlaceholderType() pkg_rhino.JSTypeExpression {
	clinit_JSDocSerializer()
	var name pkg_rhino.Node = pkg_rhino.IR_String("JSDocSerializer_placeholder_type")
	_ = name
	var bang pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), name)
	_ = bang
	name.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(StaticGet_JSDocSerializer_SYNTHETIC_SOURCE()))
	bang.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(StaticGet_JSDocSerializer_SYNTHETIC_SOURCE()))
	return pkg_rhino.NewJSTypeExpression(bang, StaticGet_JSDocSerializer_SYNTHETIC_SOURCE().GetName())
}

// Java static method: JSDocSerializer.deserializeJsdoc(com.google.javascript.jscomp.serialization.OptimizationJsdoc serializedJsdoc, com.google.javascript.jscomp.serialization.StringPool stringPool)
func JSDocSerializer_DeserializeJsdoc(serializedJsdoc *j2gort.Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc, stringPool *StringPool) *pkg_rhino.JSDocInfo {
	clinit_JSDocSerializer()
	_ = serializedJsdoc
	_ = stringPool
	if serializedJsdoc == nil {
		return nil
	}
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = builder
	if serializedJsdoc.GetClosureUnawareModePointer() != int32(0) {
		builder.RecordClosureUnawareCode_1(pkg_rhino.JSDocInfo_PerFileClosureUnawareMode_ValueOf(stringPool.Get(serializedJsdoc.GetClosureUnawareModePointer())))
	}
	if serializedJsdoc.GetLicenseTextPointer() != int32(0) {
		builder.AddLicense(stringPool.Get(serializedJsdoc.GetLicenseTextPointer()))
	}
	if serializedJsdoc.GetMeaningPointer() != int32(0) {
		builder.RecordMeaning(stringPool.Get(serializedJsdoc.GetMeaningPointer()))
	}
	if serializedJsdoc.GetDescriptionPointer() != int32(0) {
		builder.RecordDescription(stringPool.Get(serializedJsdoc.GetDescriptionPointer()))
	}
	if serializedJsdoc.GetAlternateMessageIdPointer() != int32(0) {
		builder.RecordAlternateMessageId(stringPool.Get(serializedJsdoc.GetAlternateMessageIdPointer()))
	}
	var modifies *j2gort.Ext_java_util_TreeSet = nil
	_ = modifies
	var suppressions *j2gort.Ext_java_util_TreeSet = nil
	_ = suppressions
	{
		var i int32 = int32(0)
		_ = i
		for ; i < serializedJsdoc.GetKindCount(); i = i + 1 {
			var tag *j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_JsdocTag](serializedJsdoc.GetKindList().Get(i))
			_ = tag
			switch j2gort.CheckNotNullPtr(tag).Ordinal() {
			case int32(12) /* JsdocTag.JSDOC_CONST */ :
				builder.RecordConstancy()
				continue
			case int32(9) /* JsdocTag.JSDOC_ENUM */ :
				builder.RecordEnumParameterType(StaticGet_JSDocSerializer_PlaceholderType())
				continue
			case int32(11) /* JsdocTag.JSDOC_THIS */ :
				builder.RecordThisType(StaticGet_JSDocSerializer_PlaceholderType())
				continue
			case int32(10) /* JsdocTag.JSDOC_NO_COLLAPSE */ :
				builder.RecordNoCollapse()
				continue
			case int32(1) /* JsdocTag.JSDOC_NO_INLINE */ :
				builder.RecordNoInline()
				continue
			case int32(30) /* JsdocTag.JSDOC_REQUIRE_INLINING */ :
				builder.RecordRequireInlining()
				continue
			case int32(32) /* JsdocTag.JSDOC_ENCOURAGE_INLINING */ :
				builder.RecordEncourageInlining()
				continue
			case int32(24) /* JsdocTag.JSDOC_PROVIDE_GOOG */ :
				builder.RecordProvideGoog()
				continue
			case int32(26) /* JsdocTag.JSDOC_PROVIDE_ALREADY_PROVIDED */ :
				builder.RecordProvideAlreadyProvided()
				continue
			case int32(25) /* JsdocTag.JSDOC_TYPE_SUMMARY_FILE */ :
				builder.RecordTypeSummary()
				continue
			case int32(21) /* JsdocTag.JSDOC_PURE_OR_BREAK_MY_CODE */ :
				builder.RecordPureOrBreakMyCode()
				continue
			case int32(22) /* JsdocTag.JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE */ :
				builder.RecordCollapsibleOrBreakMyCode()
				continue
			case int32(2) /* JsdocTag.JSDOC_DEFINE */ :
				builder.RecordDefineType(StaticGet_JSDocSerializer_PlaceholderType())
				continue
			case int32(6) /* JsdocTag.JSDOC_NO_SIDE_EFFECTS */ :
				builder.RecordNoSideEffects()
				continue
			case int32(3) /* JsdocTag.JSDOC_MODIFIES_THIS */ :
				if modifies != nil {
					modifies = modifies
				} else {
					modifies = j2gort.NewExt_java_util_TreeSet()
				}
				j2gort.SetAddStr(modifies, "this")
				continue
			case int32(4) /* JsdocTag.JSDOC_MODIFIES_ARGUMENTS */ :
				if modifies != nil {
					modifies = modifies
				} else {
					modifies = j2gort.NewExt_java_util_TreeSet()
				}
				j2gort.SetAddStr(modifies, "arguments")
				continue
			case int32(5) /* JsdocTag.JSDOC_THROWS */ :
				builder.RecordThrowsAnnotation(("{!" + "JSDocSerializer_placeholder_type" + "}"))
				continue
			case int32(7) /* JsdocTag.JSDOC_CONSTRUCTOR */ :
				builder.RecordConstructor()
				continue
			case int32(8) /* JsdocTag.JSDOC_INTERFACE */ :
				builder.RecordInterface()
				continue
			case int32(20) /* JsdocTag.JSDOC_SUPPRESS_PARTIAL_ALIAS */ :
				if suppressions != nil {
					suppressions = suppressions
				} else {
					suppressions = j2gort.NewExt_java_util_TreeSet()
				}
				j2gort.SetAddStr(suppressions, "partialAlias")
				continue
			case int32(13) /* JsdocTag.JSDOC_ID_GENERATOR_CONSISTENT */ :
				builder.RecordConsistentIdGenerator()
				continue
			case int32(15) /* JsdocTag.JSDOC_ID_GENERATOR_STABLE */ :
				builder.RecordStableIdGenerator()
				continue
			case int32(16) /* JsdocTag.JSDOC_ID_GENERATOR_MAPPED */ :
				builder.RecordMappedIdGenerator()
				continue
			case int32(17) /* JsdocTag.JSDOC_ID_GENERATOR_XID */ :
				builder.RecordXidGenerator()
				continue
			case int32(14) /* JsdocTag.JSDOC_ID_GENERATOR_INCONSISTENT */ :
				builder.RecordIdGenerator()
				continue
			case int32(18) /* JsdocTag.JSDOC_ABSTRACT */ :
				builder.RecordAbstract()
				continue
			case int32(19) /* JsdocTag.JSDOC_SUPPRESS_MESSAGE_CONVENTION */ :
				if suppressions != nil {
					suppressions = suppressions
				} else {
					suppressions = j2gort.NewExt_java_util_TreeSet()
				}
				j2gort.SetAddStr(suppressions, "messageConventions")
				continue
			case int32(28) /* JsdocTag.JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES */ :
				if suppressions != nil {
					suppressions = suppressions
				} else {
					suppressions = j2gort.NewExt_java_util_TreeSet()
				}
				j2gort.SetAddStr(suppressions, "untranspilableFeatures")
				continue
			case int32(23) /* JsdocTag.JSDOC_FILEOVERVIEW */ :
				builder.RecordFileOverview("")
				continue
			case int32(29) /* JsdocTag.JSDOC_NO_COVERAGE */ :
				builder.RecordNoCoverage()
				continue
			case int32(27) /* JsdocTag.JSDOC_SASS_GENERATED_CSS_TS */ :
				builder.RecordSassGeneratedCssTs()
				continue
			case int32(31) /* JsdocTag.JSDOC_USED_VIA_DOT_CONSTRUCTOR */ :
				var unused bool = builder.RecordUsedViaDotConstructor()
				_ = unused
				continue
			case int32(0) /* JsdocTag.JSDOC_UNSPECIFIED */, int32(33) /* JsdocTag.UNRECOGNIZED */ :
				panic(NewMalformedTypedAstException_1("Unsupported JSDoc tag can't be deserialized: " + j2gort.StrVal(tag)))
			}
		}
	}
	if modifies != nil {
		builder.RecordModifies(j2gort.AsIface[j2gort.Ext_java_util_Set](modifies))
	}
	if suppressions != nil {
		builder.RecordSuppressions_1(j2gort.AsIface[j2gort.Ext_java_util_Set](suppressions))
	}
	return builder.Build_0()
}
