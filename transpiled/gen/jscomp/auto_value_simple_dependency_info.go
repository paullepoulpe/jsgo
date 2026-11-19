// Port of com/google/javascript/jscomp/deps/AutoValue_SimpleDependencyInfo.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo
type AutoValue_SimpleDependencyInfo struct {
	SimpleDependencyInfoBase
	name                      string                                             `gson:"name"`
	pathRelativeToClosureBase string                                             `gson:"pathRelativeToClosureBase"`
	provides                  j2gort.Ext_com_google_common_collect_ImmutableList `gson:"provides"`
	requires                  j2gort.Ext_com_google_common_collect_ImmutableList `gson:"requires"`
	typeRequires              j2gort.Ext_com_google_common_collect_ImmutableList `gson:"typeRequires"`
	loadFlags                 j2gort.Ext_com_google_common_collect_ImmutableMap  `gson:"loadFlags"`
	hasExternsAnnotation      bool                                               `gson:"hasExternsAnnotation"`
	hasNoCompileAnnotation    bool                                               `gson:"hasNoCompileAnnotation"`
}

func (this_ *AutoValue_SimpleDependencyInfo) BindSelf_(s any) {
	_ = s
	this_.SimpleDependencyInfoBase.BindSelf_(s)
}
func (this_ *AutoValue_SimpleDependencyInfo) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_SimpleDependencyInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo", "AutoValue_SimpleDependencyInfo"
}
func (this_ *AutoValue_SimpleDependencyInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SimpleDependencyInfoBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo", "com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_SimpleDependencyInfo](v_)
	})
}
func (*AutoValue_SimpleDependencyInfo) Implements_jscomp_SimpleDependencyInfo() {}
func (*AutoValue_SimpleDependencyInfo) Implements_jscomp_DependencyInfo()       {}

func (this_ *AutoValue_SimpleDependencyInfo) initFields_AutoValue_SimpleDependencyInfo() {
	_ = this_
	this_.name = j2gort.NullString
	this_.pathRelativeToClosureBase = j2gort.NullString
}

// Java constructor: AutoValue_SimpleDependencyInfo(java.lang.String name, java.lang.String pathRelativeToClosureBase, com.google.common.collect.ImmutableList<java.lang.String> provides, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.deps.DependencyInfo.Require> requires, com.google.common.collect.ImmutableList<java.lang.String> typeRequires, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> loadFlags, boolean hasExternsAnnotation, boolean hasNoCompileAnnotation)
func NewAutoValue_SimpleDependencyInfo(name string, pathRelativeToClosureBase string, provides j2gort.Ext_com_google_common_collect_ImmutableList, requires j2gort.Ext_com_google_common_collect_ImmutableList, typeRequires j2gort.Ext_com_google_common_collect_ImmutableList, loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap, hasExternsAnnotation bool, hasNoCompileAnnotation bool) *AutoValue_SimpleDependencyInfo {
	this_ := &AutoValue_SimpleDependencyInfo{}
	this_.BindSelf_(this_)
	InitAutoValue_SimpleDependencyInfo(this_, name, pathRelativeToClosureBase, provides, requires, typeRequires, loadFlags, hasExternsAnnotation, hasNoCompileAnnotation)
	return this_
}

func InitAutoValue_SimpleDependencyInfo(this_ *AutoValue_SimpleDependencyInfo, name string, pathRelativeToClosureBase string, provides j2gort.Ext_com_google_common_collect_ImmutableList, requires j2gort.Ext_com_google_common_collect_ImmutableList, typeRequires j2gort.Ext_com_google_common_collect_ImmutableList, loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap, hasExternsAnnotation bool, hasNoCompileAnnotation bool) {
	_ = name
	_ = pathRelativeToClosureBase
	_ = provides
	_ = requires
	_ = typeRequires
	_ = loadFlags
	_ = hasExternsAnnotation
	_ = hasNoCompileAnnotation
	InitSimpleDependencyInfo(&this_.SimpleDependencyInfoBase)
	this_.initFields_AutoValue_SimpleDependencyInfo()
	this_.name = name
	this_.pathRelativeToClosureBase = pathRelativeToClosureBase
	this_.provides = provides
	this_.requires = requires
	this_.typeRequires = typeRequires
	this_.loadFlags = loadFlags
	this_.hasExternsAnnotation = hasExternsAnnotation
	this_.hasNoCompileAnnotation = hasNoCompileAnnotation
}

// Java method: AutoValue_SimpleDependencyInfo.getName()
func (this_ *AutoValue_SimpleDependencyInfo) GetName() string {
	return this_.name
}

// Java method: AutoValue_SimpleDependencyInfo.getPathRelativeToClosureBase()
func (this_ *AutoValue_SimpleDependencyInfo) GetPathRelativeToClosureBase() string {
	return this_.pathRelativeToClosureBase
}

// Java method: AutoValue_SimpleDependencyInfo.getProvides()
func (this_ *AutoValue_SimpleDependencyInfo) GetProvides() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.provides
}

// Java method: AutoValue_SimpleDependencyInfo.getRequires()
func (this_ *AutoValue_SimpleDependencyInfo) GetRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.requires
}

// Java method: AutoValue_SimpleDependencyInfo.getTypeRequires()
func (this_ *AutoValue_SimpleDependencyInfo) GetTypeRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.typeRequires
}

// Java method: AutoValue_SimpleDependencyInfo.getLoadFlags()
func (this_ *AutoValue_SimpleDependencyInfo) GetLoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.loadFlags
}

// Java method: AutoValue_SimpleDependencyInfo.getHasExternsAnnotation()
func (this_ *AutoValue_SimpleDependencyInfo) GetHasExternsAnnotation() bool {
	return this_.hasExternsAnnotation
}

// Java method: AutoValue_SimpleDependencyInfo.getHasNoCompileAnnotation()
func (this_ *AutoValue_SimpleDependencyInfo) GetHasNoCompileAnnotation() bool {
	return this_.hasNoCompileAnnotation
}

// Java method: AutoValue_SimpleDependencyInfo.toString()
func (this_ *AutoValue_SimpleDependencyInfo) String() string {
	return ("SimpleDependencyInfo{name=" + j2gort.StrValString(this_.name) + ", pathRelativeToClosureBase=" + j2gort.StrValString(this_.pathRelativeToClosureBase) + ", provides=" + j2gort.StrVal(this_.provides) + ", requires=" + j2gort.StrVal(this_.requires) + ", typeRequires=" + j2gort.StrVal(this_.typeRequires) + ", loadFlags=" + j2gort.StrVal(this_.loadFlags) + ", hasExternsAnnotation=" + j2gort.StrValBool(this_.hasExternsAnnotation) + ", hasNoCompileAnnotation=" + j2gort.StrValBool(this_.hasNoCompileAnnotation) + "}")
}

// Java method: AutoValue_SimpleDependencyInfo.equals(java.lang.Object o)
func (this_ *AutoValue_SimpleDependencyInfo) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[SimpleDependencyInfo](o) {
		var that SimpleDependencyInfo = j2gort.AsIface[SimpleDependencyInfo](o)
		_ = that
		return (((((((func() bool {
			sr_2 := this_.name
			se_3 := that.GetName()
			return j2gort.CheckNotNullString(sr_2) == se_3
		}() && func() bool {
			sr_4 := this_.pathRelativeToClosureBase
			se_5 := that.GetPathRelativeToClosureBase()
			return j2gort.CheckNotNullString(sr_4) == se_5
		}()) && this_.provides.Equals(that.GetProvides())) && this_.requires.Equals(that.GetRequires())) && this_.typeRequires.Equals(that.GetTypeRequires())) && this_.loadFlags.Equals(that.GetLoadFlags())) && (this_.hasExternsAnnotation == that.GetHasExternsAnnotation())) && (this_.hasNoCompileAnnotation == that.GetHasNoCompileAnnotation()))
	}
	return false
}

// Java method: AutoValue_SimpleDependencyInfo.hashCode()
func (this_ *AutoValue_SimpleDependencyInfo) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.name)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.pathRelativeToClosureBase)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.provides.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.requires.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.typeRequires.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.loadFlags.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.hasExternsAnnotation {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.hasNoCompileAnnotation {
			return int32(1231)
		}
		return int32(1237)
	}())
	return h_d
}

// Java class: com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo.Builder
type AutoValue_SimpleDependencyInfo_Builder struct {
	SimpleDependencyInfo_BuilderBase
	name                      string                                             `gson:"name"`
	pathRelativeToClosureBase string                                             `gson:"pathRelativeToClosureBase"`
	provides                  j2gort.Ext_com_google_common_collect_ImmutableList `gson:"provides"`
	requires                  j2gort.Ext_com_google_common_collect_ImmutableList `gson:"requires"`
	typeRequires              j2gort.Ext_com_google_common_collect_ImmutableList `gson:"typeRequires"`
	loadFlags                 j2gort.Ext_com_google_common_collect_ImmutableMap  `gson:"loadFlags"`
	hasExternsAnnotation      bool                                               `gson:"hasExternsAnnotation"`
	hasNoCompileAnnotation    bool                                               `gson:"hasNoCompileAnnotation"`
	set_d0                    int8                                               `gson:"set$0"`
}

func (this_ *AutoValue_SimpleDependencyInfo_Builder) BindSelf_(s any) {
	_ = s
	this_.SimpleDependencyInfo_BuilderBase.BindSelf_(s)
}
func (this_ *AutoValue_SimpleDependencyInfo_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_SimpleDependencyInfo_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo$Builder", "Builder"
}
func (this_ *AutoValue_SimpleDependencyInfo_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SimpleDependencyInfo_BuilderBase.IdentityHash_()
}
func (*AutoValue_SimpleDependencyInfo_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo.Builder", "com.google.javascript.jscomp.deps.AutoValue_SimpleDependencyInfo$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_SimpleDependencyInfo_Builder](v_)
	})
}
func (*AutoValue_SimpleDependencyInfo_Builder) Implements_jscomp_SimpleDependencyInfo_Builder() {}

func (this_ *AutoValue_SimpleDependencyInfo_Builder) initFields_AutoValue_SimpleDependencyInfo_Builder() {
	_ = this_
	this_.name = j2gort.NullString
	this_.pathRelativeToClosureBase = j2gort.NullString
}

// Java constructor: Builder()
func NewAutoValue_SimpleDependencyInfo_Builder() *AutoValue_SimpleDependencyInfo_Builder {
	this_ := &AutoValue_SimpleDependencyInfo_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_SimpleDependencyInfo_Builder(this_)
	return this_
}

func InitAutoValue_SimpleDependencyInfo_Builder(this_ *AutoValue_SimpleDependencyInfo_Builder) {
	InitSimpleDependencyInfo_Builder(&this_.SimpleDependencyInfo_BuilderBase)
	this_.initFields_AutoValue_SimpleDependencyInfo_Builder()
}

// Java method: Builder.setName(java.lang.String name)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetName(name string) SimpleDependencyInfo_Builder {
	_ = name
	if name == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null name")))
	}
	this_.name = name
	return this_
}

// Java method: Builder.setPathRelativeToClosureBase(java.lang.String pathRelativeToClosureBase)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetPathRelativeToClosureBase(pathRelativeToClosureBase string) SimpleDependencyInfo_Builder {
	_ = pathRelativeToClosureBase
	if pathRelativeToClosureBase == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null pathRelativeToClosureBase")))
	}
	this_.pathRelativeToClosureBase = pathRelativeToClosureBase
	return this_
}

// Java method: Builder.setProvides(java.util.Collection<java.lang.String> provides)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetProvides_1_Ext_java_util_Collection(provides j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder {
	_ = provides
	this_.provides = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(provides)
	return this_
}

// Java method: Builder.setProvides(java.lang.String[] provides)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetProvides_1_slice_string(provides []string) SimpleDependencyInfo_Builder {
	_ = provides
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_6 := provides
		if arr_6 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_7 := make([]any, len(arr_6))
		for i_, el_8 := range arr_6 {
			out_7[i_] = j2gort.BoxString(el_8)
		}
		iifeArg_1 = out_7
		break iifeSwitch_1
	}
	this_.provides = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(iifeArg_1)
	return this_
}

// Java method: Builder.setRequires(java.util.Collection<com.google.javascript.jscomp.deps.DependencyInfo.Require> requires)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetRequires_1_Ext_java_util_Collection(requires j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder {
	_ = requires
	this_.requires = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(requires)
	return this_
}

// Java method: Builder.setRequires(com.google.javascript.jscomp.deps.DependencyInfo.Require[] requires)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetRequires_1_DependencyInfo_Require(requires []*DependencyInfo_Require) SimpleDependencyInfo_Builder {
	_ = requires
	var iifeArg_2 []any
iifeSwitch_2:
	switch {
	default:

		arr_9 := requires
		if arr_9 == nil {
			iifeArg_2 = nil
			break iifeSwitch_2
		}
		out_10 := make([]any, len(arr_9))
		for i_, el_11 := range arr_9 {
			out_10[i_] = j2gort.PtrToAny(el_11)
		}
		iifeArg_2 = out_10
		break iifeSwitch_2
	}
	this_.requires = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(iifeArg_2)
	return this_
}

// Java method: Builder.setTypeRequires(java.util.Collection<java.lang.String> typeRequires)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetTypeRequires_1_Ext_java_util_Collection(typeRequires j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder {
	_ = typeRequires
	this_.typeRequires = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(typeRequires)
	return this_
}

// Java method: Builder.setTypeRequires(java.lang.String[] typeRequires)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetTypeRequires_1_slice_string(typeRequires []string) SimpleDependencyInfo_Builder {
	_ = typeRequires
	var iifeArg_3 []any
iifeSwitch_3:
	switch {
	default:

		arr_12 := typeRequires
		if arr_12 == nil {
			iifeArg_3 = nil
			break iifeSwitch_3
		}
		out_13 := make([]any, len(arr_12))
		for i_, el_14 := range arr_12 {
			out_13[i_] = j2gort.BoxString(el_14)
		}
		iifeArg_3 = out_13
		break iifeSwitch_3
	}
	this_.typeRequires = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(iifeArg_3)
	return this_
}

// Java method: Builder.setLoadFlags(java.util.Map<java.lang.String,java.lang.String> loadFlags)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetLoadFlags(loadFlags j2gort.Ext_java_util_Map) SimpleDependencyInfo_Builder {
	_ = loadFlags
	this_.loadFlags = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(loadFlags)
	return this_
}

// Java method: Builder.setHasExternsAnnotation(boolean hasExternsAnnotation)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetHasExternsAnnotation(hasExternsAnnotation bool) SimpleDependencyInfo_Builder {
	_ = hasExternsAnnotation
	this_.hasExternsAnnotation = hasExternsAnnotation
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.setHasNoCompileAnnotation(boolean hasNoCompileAnnotation)
func (this_ *AutoValue_SimpleDependencyInfo_Builder) SetHasNoCompileAnnotation(hasNoCompileAnnotation bool) SimpleDependencyInfo_Builder {
	_ = hasNoCompileAnnotation
	this_.hasNoCompileAnnotation = hasNoCompileAnnotation
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(2))))
	return this_
}

// Java method: Builder.build()
func (this_ *AutoValue_SimpleDependencyInfo_Builder) Build() SimpleDependencyInfo {
	if ((((((int32(this_.set_d0) != int32(3)) || (this_.name == j2gort.NullString)) || (this_.pathRelativeToClosureBase == j2gort.NullString)) || (this_.provides == nil)) || (this_.requires == nil)) || (this_.typeRequires == nil)) || (this_.loadFlags == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.name == j2gort.NullString {
			missing.Append_1_string(" name")
		}
		if this_.pathRelativeToClosureBase == j2gort.NullString {
			missing.Append_1_string(" pathRelativeToClosureBase")
		}
		if this_.provides == nil {
			missing.Append_1_string(" provides")
		}
		if this_.requires == nil {
			missing.Append_1_string(" requires")
		}
		if this_.typeRequires == nil {
			missing.Append_1_string(" typeRequires")
		}
		if this_.loadFlags == nil {
			missing.Append_1_string(" loadFlags")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" hasExternsAnnotation")
		}
		if (int32(this_.set_d0) & int32(2)) == int32(0) {
			missing.Append_1_string(" hasNoCompileAnnotation")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAutoValue_SimpleDependencyInfo(this_.name, this_.pathRelativeToClosureBase, this_.provides, this_.requires, this_.typeRequires, this_.loadFlags, this_.hasExternsAnnotation, this_.hasNoCompileAnnotation)
}
