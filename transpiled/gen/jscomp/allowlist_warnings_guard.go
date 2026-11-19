/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AllowlistWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AllowlistWarningsGuard
type AllowlistWarningsGuard struct {
	WarningsGuardBase
	allowlist j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"allowlist"`
}

func (this_ *AllowlistWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *AllowlistWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*AllowlistWarningsGuard) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AllowlistWarningsGuard", "AllowlistWarningsGuard"
}
func (this_ *AllowlistWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*AllowlistWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AllowlistWarningsGuard", "com.google.javascript.jscomp.AllowlistWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*AllowlistWarningsGuard](v_)
	})
}
func (*AllowlistWarningsGuard) Implements_jscomp_WarningsGuard() {}

var static_AllowlistWarningsGuard_LINE_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_AllowlistWarningsGuard_LINE_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_AllowlistWarningsGuard()
	return static_AllowlistWarningsGuard_LINE_SPLITTER
}
func StaticSet_AllowlistWarningsGuard_LINE_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_AllowlistWarningsGuard()
	static_AllowlistWarningsGuard_LINE_SPLITTER = v
}

var static_AllowlistWarningsGuard_LINE_NUMBER *j2gort.Ext_java_util_regex_Pattern

func StaticGet_AllowlistWarningsGuard_LINE_NUMBER() *j2gort.Ext_java_util_regex_Pattern {
	clinit_AllowlistWarningsGuard()
	return static_AllowlistWarningsGuard_LINE_NUMBER
}
func StaticSet_AllowlistWarningsGuard_LINE_NUMBER(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_AllowlistWarningsGuard()
	static_AllowlistWarningsGuard_LINE_NUMBER = v
}

var clinitDone_AllowlistWarningsGuard uint32
var clinitState_AllowlistWarningsGuard j2gort.ClinitState

func clinit_AllowlistWarningsGuard() {
	if atomic.LoadUint32(&clinitDone_AllowlistWarningsGuard) != 0 {
		return
	}
	clinitSlow_AllowlistWarningsGuard()
}

func clinitSlow_AllowlistWarningsGuard() {
	if !j2gort.EnterClinit(&clinitState_AllowlistWarningsGuard, &clinitDone_AllowlistWarningsGuard) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AllowlistWarningsGuard, &clinitDone_AllowlistWarningsGuard)
	static_AllowlistWarningsGuard_LINE_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(10))
	static_AllowlistWarningsGuard_LINE_NUMBER = j2gort.Ext_java_util_regex_Pattern_Compile_1(":-?\\d+")
}

func (this_ *AllowlistWarningsGuard) initFields_AllowlistWarningsGuard() {
	_ = this_
}

// Java constructor: AllowlistWarningsGuard()
func NewAllowlistWarningsGuard() *AllowlistWarningsGuard {
	clinit_AllowlistWarningsGuard()
	this_ := &AllowlistWarningsGuard{}
	InitAllowlistWarningsGuard(this_)
	return this_
}

func InitAllowlistWarningsGuard(this_ *AllowlistWarningsGuard) {
	InitAllowlistWarningsGuard_1(this_, j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
}

// Java constructor: AllowlistWarningsGuard(java.util.Set<java.lang.String> allowlist)
func NewAllowlistWarningsGuard_1(allowlist j2gort.Ext_java_util_Set) *AllowlistWarningsGuard {
	clinit_AllowlistWarningsGuard()
	this_ := &AllowlistWarningsGuard{}
	InitAllowlistWarningsGuard_1(this_, allowlist)
	return this_
}

func InitAllowlistWarningsGuard_1(this_ *AllowlistWarningsGuard, allowlist j2gort.Ext_java_util_Set) {
	_ = allowlist
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_AllowlistWarningsGuard()
	j2gort.CheckNotNull(allowlist)
	this_.allowlist = AllowlistWarningsGuard_NormalizeAllowlist(allowlist)
}

// Java static method: AllowlistWarningsGuard.normalizeAllowlist(java.util.Set<java.lang.String> allowlist)
func AllowlistWarningsGuard_NormalizeAllowlist(allowlist j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AllowlistWarningsGuard()
	_ = allowlist
	var result j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = result
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(allowlist))
		for cur_1.HasNext() {
			line := cur_1.NextStr()
			_ = line
			var trimmed string = javalang.StringTrim(j2gort.CheckNotNullString(line))
			_ = trimmed
			if (j2gort.CheckNotNullString(trimmed) == "") || (int32(j2gort.JStrCharAt(trimmed, int32(0))) == int32(uint16(35))) {
				continue
			}
			result.Add(StaticGet_AllowlistWarningsGuard_LINE_NUMBER().Matcher(j2gort.BoxString(trimmed)).ReplaceFirst_1_string(":"))
		}
	}
	return result.ToImmutableSet()
}

// Java method: AllowlistWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *AllowlistWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	if j2gort.RefEquals(error_.DefaultLevel(), StaticGet_CheckLevel_ERROR()) {
		return nil
	}
	if !(this_.allowlist.IsEmpty()) && this_.ContainWarning(this_.FormatWarning_1(error_)) {
		return StaticGet_CheckLevel_WARNING()
	}
	return nil
}

// Java method: AllowlistWarningsGuard.containWarning(java.lang.String formattedWarning)
func (this_ *AllowlistWarningsGuard) ContainWarning(formattedWarning string) bool {
	_ = formattedWarning
	return j2gort.SetContainsStr(this_.allowlist, formattedWarning)
}

// Java method: AllowlistWarningsGuard.getPriority()
func (this_ *AllowlistWarningsGuard) GetPriority() int32 {
	return StaticGet_WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST().GetValue()
}

// Java static method: AllowlistWarningsGuard.fromFile(java.io.File file)
func AllowlistWarningsGuard_FromFile(file *j2gort.Ext_java_io_File) *AllowlistWarningsGuard {
	clinit_AllowlistWarningsGuard()
	_ = file
	return NewAllowlistWarningsGuard_1(AllowlistWarningsGuard_LoadAllowlistedJsWarnings_1_Ext_java_io_File(file))
}

// Java static method: AllowlistWarningsGuard.loadAllowlistedJsWarnings(java.io.File file)
func AllowlistWarningsGuard_LoadAllowlistedJsWarnings_1_Ext_java_io_File(file *j2gort.Ext_java_io_File) j2gort.Ext_java_util_Set {
	clinit_AllowlistWarningsGuard()
	_ = file
	return AllowlistWarningsGuard_LoadAllowlistedJsWarnings_1_Ext_com_google_common_io_CharSource(j2gort.Ext_com_google_common_io_Files_AsCharSource(file, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8))
}

// Java static method: AllowlistWarningsGuard.loadAllowlistedJsWarnings(com.google.common.io.CharSource supplier)
func AllowlistWarningsGuard_LoadAllowlistedJsWarnings_1_Ext_com_google_common_io_CharSource(supplier j2gort.Ext_com_google_common_io_CharSource) j2gort.Ext_java_util_Set {
	clinit_AllowlistWarningsGuard()
	_ = supplier
	var didRet_2 bool
	var retVal_3 j2gort.Ext_java_util_Set
	func() {
		defer func() {
			if rec_4 := recover(); rec_4 != nil {
				jex_5 := javalang.RecoverException(rec_4)
				_ = jex_5
				if j2gort.InstanceOf[javaio.IIOException](jex_5) {
					e := j2gort.AsIface[javaio.IIOException](jex_5)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_4)
				}
			}
		}()
		retVal_3 = AllowlistWarningsGuard_LoadAllowlistedJsWarnings_1_Ext_java_io_Reader(supplier.OpenStream())
		didRet_2 = true
		return
	}()
	if didRet_2 {
		return retVal_3
	}
	panic("unreachable")
}

// Java static method: AllowlistWarningsGuard.loadAllowlistedJsWarnings(java.io.Reader reader)
func AllowlistWarningsGuard_LoadAllowlistedJsWarnings_1_Ext_java_io_Reader(reader j2gort.Ext_java_io_Reader) j2gort.Ext_java_util_Set {
	clinit_AllowlistWarningsGuard()
	_ = reader
	j2gort.CheckNotNull(reader)
	var result j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = result
	result.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_io_CharStreams_ReadLines_1(j2gort.AsIface[j2gort.Ext_java_lang_Readable](reader))))
	return result
}

// Java method: AllowlistWarningsGuard.formatWarning(com.google.javascript.jscomp.JSError error)
func (this_ *AllowlistWarningsGuard) FormatWarning_1(error_ *JSError) string {
	_ = error_
	return this_.FormatWarning_2(error_, false)
}

// Java method: AllowlistWarningsGuard.formatWarning(com.google.javascript.jscomp.JSError error, boolean withMetaData)
func (this_ *AllowlistWarningsGuard) FormatWarning_2(error_ *JSError, withMetaData bool) string {
	_ = error_
	_ = withMetaData
	var sb j2gort.LocalStringBuilder

	sb.AppendStr(this_.priv_AllowlistWarningsGuard_normalizeSourceName_556(error_.SourceName())).AppendStr(":")
	if withMetaData {
		sb.AppendInt(error_.GetLineNumber())
	}
	var lines j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(StaticGet_AllowlistWarningsGuard_LINE_SPLITTER(), error_.Description())
	_ = lines
	sb.AppendStr("  ").AppendStr(j2gort.ListGetStr(lines, int32(0)))
	if withMetaData {
		{
			var i int32 = int32(1)
			_ = i
			for ; i < lines.Size(); i = i + 1 {
				sb.AppendStr("\n# ").AppendStr(j2gort.ListGetStr(lines, i))
			}
		}
		sb.AppendStr("\n")
	}
	return sb.String()
}

// Java method: AllowlistWarningsGuard.normalizeSourceName(java.lang.String sourceName)
func (this_ *AllowlistWarningsGuard) priv_AllowlistWarningsGuard_normalizeSourceName_556(sourceName string) string {
	_ = sourceName
	if sourceName != j2gort.NullString {
		return j2gort.JStrReplaceFirst(sourceName, "blaze-out/[^/]*/(bin|genfiles)/", "")
	}
	return sourceName
}

// Java static method: AllowlistWarningsGuard.getFirstLine(java.lang.String warning)
func AllowlistWarningsGuard_GetFirstLine(warning string) string {
	clinit_AllowlistWarningsGuard()
	_ = warning
	var lineLength int32 = j2gort.JStrIndexOfCp(warning, int32(uint16(10)))
	_ = lineLength
	if lineLength > int32(0) {
		warning = j2gort.JStrSub(warning, int32(0), lineLength)
	}
	return warning
}

// Java class: com.google.javascript.jscomp.AllowlistWarningsGuard.AllowlistBuilder
type AllowlistWarningsGuard_AllowlistBuilder struct {
	outer_          *AllowlistWarningsGuard
	warnings        j2gort.CustomOrderedSet[*JSError] `gson:"warnings"`
	productName     string                            `gson:"productName"`
	generatorTarget string                            `gson:"generatorTarget"`
	headerNote      string                            `gson:"headerNote"`
	idHash_         atomic.Int32
}

func (this_ *AllowlistWarningsGuard_AllowlistBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) IsJavaNil_() bool { return this_ == nil }
func (*AllowlistWarningsGuard_AllowlistBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AllowlistWarningsGuard$AllowlistBuilder", "AllowlistBuilder"
}
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AllowlistWarningsGuard_AllowlistBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AllowlistWarningsGuard.AllowlistBuilder", "com.google.javascript.jscomp.AllowlistWarningsGuard$AllowlistBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*AllowlistWarningsGuard_AllowlistBuilder](v_)
	})
}
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) SetOuter_(o *AllowlistWarningsGuard) {
	this_.outer_ = o
}
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) GetOuter_() *AllowlistWarningsGuard {
	return this_.outer_
}
func (*AllowlistWarningsGuard_AllowlistBuilder) Implements_jscomp_ErrorHandler() {}

func (this_ *AllowlistWarningsGuard_AllowlistBuilder) initFields_AllowlistWarningsGuard_AllowlistBuilder() {
	_ = this_
	this_.warnings = j2gort.NewCustomOrderedSet[*JSError]()
	this_.productName = j2gort.NullString
	this_.generatorTarget = j2gort.NullString
	this_.headerNote = j2gort.NullString
}

// Java constructor: AllowlistBuilder()
func NewAllowlistWarningsGuard_AllowlistBuilder(outer_ *AllowlistWarningsGuard) *AllowlistWarningsGuard_AllowlistBuilder {
	this_ := &AllowlistWarningsGuard_AllowlistBuilder{}
	this_.outer_ = outer_
	InitAllowlistWarningsGuard_AllowlistBuilder(this_)
	return this_
}

func InitAllowlistWarningsGuard_AllowlistBuilder(this_ *AllowlistWarningsGuard_AllowlistBuilder) {
	this_.initFields_AllowlistWarningsGuard_AllowlistBuilder()
}

// Java method: AllowlistBuilder.setProductName(java.lang.String name)
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) SetProductName(name string) *AllowlistWarningsGuard_AllowlistBuilder {
	_ = name
	this_.productName = name
	return this_
}

// Java method: AllowlistBuilder.setGeneratorTarget(java.lang.String name)
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) SetGeneratorTarget(name string) *AllowlistWarningsGuard_AllowlistBuilder {
	_ = name
	this_.generatorTarget = name
	return this_
}

// Java method: AllowlistBuilder.setNote(java.lang.String note)
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) SetNote(note string) *AllowlistWarningsGuard_AllowlistBuilder {
	_ = note
	this_.headerNote = note
	return this_
}

// Java method: AllowlistBuilder.report(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) Report_2(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	if j2gort.RefEquals(error_.DefaultLevel(), StaticGet_CheckLevel_ERROR()) {
		return
	}
	this_.warnings.Add(error_)
}

// Java method: AllowlistBuilder.writeAllowlist(java.io.File out)
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) WriteAllowlist(out *j2gort.Ext_java_io_File) {
	_ = out
	func() {
		var stream *j2gort.Ext_java_io_PrintStream = j2gort.NewExt_java_io_PrintStream_1_Ext_java_io_File(out)
		_ = stream
		defer j2gort.CloseResource(stream)
		this_.AppendAllowlist(stream)
	}()
}

// Java method: AllowlistBuilder.appendAllowlist(java.io.PrintStream out)
func (this_ *AllowlistWarningsGuard_AllowlistBuilder) AppendAllowlist(out *j2gort.Ext_java_io_PrintStream) {
	_ = out
	out.Append_1_any("# This is a list of legacy warnings that have yet to be fixed.\n")
	if ((this_.productName != j2gort.NullString) && !(j2gort.CheckNotNullString(this_.productName) == "")) && !(this_.warnings.IsEmpty()) {
		out.Append_1_any("# Please find some time and fix at least one of them and it will be the happiest day for " + j2gort.StrValString(this_.productName) + ".\n")
	}
	if (this_.generatorTarget != j2gort.NullString) && !(j2gort.CheckNotNullString(this_.generatorTarget) == "") {
		out.Append_1_any("# When you fix any of these warnings, run " + j2gort.StrValString(this_.generatorTarget) + " task.\n")
	}
	if this_.headerNote != j2gort.NullString {
		out.Append_1_any("#" + j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n# ").Join_1_Ext_java_lang_Iterable(j2gort.SplitterSplitStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(10)), this_.headerNote))) + "\n")
	}
	var warningsByType *SimpleMultimap[*DiagnosticType, string] = SimpleMultimap_Create_0[*DiagnosticType, string]()
	_ = warningsByType
	{
		cur_6 := this_.warnings.Iter()
		for cur_6.HasNext() {

			warning := cur_6.
				Next()
			_ = warning
			warningsByType.Put(warning.Type(), this_.outer_.FormatWarning_2(warning, true))
		}
	}
	var sortedTypes []*DiagnosticType = j2gort.CopyCollectionToSlice[*DiagnosticType](j2gort.AsIface[j2gort.Ext_java_util_Collection](warningsByType.KeySet()))
	_ = sortedTypes
	j2gort.SliceSortNatural(sortedTypes)
	for _, type_ := range sortedTypes {
		_ = type_
		if StaticGet_DiagnosticGroups_DEPRECATED().Matches_1_DiagnosticType(type_) {
			continue
		}
		j2gort.CastPtr[j2gort.Ext_java_io_PrintStream](j2gort.CastPtr[j2gort.Ext_java_io_PrintStream](j2gort.CastPtr[j2gort.Ext_java_io_PrintStream](out.Append_1_any("\n# Warning ")).Append_1_any(j2gort.BoxString(type_.key))).Append_1_any(": ")).Println_1_string(j2gort.CastString(j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.SplitterSplitStr(StaticGet_AllowlistWarningsGuard_LINE_SPLITTER(), type_.format).Iterator()).Next()))
		{
			cur_10 := j2gort.NewIterCursor(j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](warningsByType.Get_1_any(j2gort.PtrToAny(type_)))))
			for cur_10.HasNext() {
				warning_2 := cur_10.NextStr()
				_ = warning_2
				out.Println_1_string(warning_2)
			}
		}
	}

	out.Flush()
}
