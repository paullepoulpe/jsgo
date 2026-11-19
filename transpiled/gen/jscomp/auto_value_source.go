// Port of com/google/javascript/jscomp/bundle/AutoValue_Source.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.bundle.AutoValue_Source
type AutoValue_Source struct {
	SourceBase
	path                 j2gort.Ext_java_nio_file_Path                     `gson:"path"`
	sourceMap            string                                            `gson:"sourceMap"`
	sourceUrl            string                                            `gson:"sourceUrl"`
	sourceMappingUrl     string                                            `gson:"sourceMappingUrl"`
	runtimes             j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"runtimes"`
	loadFlags            j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"loadFlags"`
	codeSupplier         Source_Lazy                                       `gson:"codeSupplier"`
	originalCodeSupplier Source_Lazy                                       `gson:"originalCodeSupplier"`
	estimatedSize        int32                                             `gson:"estimatedSize"`
}

func (this_ *AutoValue_Source) BindSelf_(s any) {
	_ = s
	this_.SourceBase.BindSelf_(s)
}
func (this_ *AutoValue_Source) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_Source) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.AutoValue_Source", "AutoValue_Source"
}
func (this_ *AutoValue_Source) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SourceBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.AutoValue_Source", "com.google.javascript.jscomp.bundle.AutoValue_Source", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_Source](v_)
	})
}
func (*AutoValue_Source) Implements_jscomp_Source() {}

func (this_ *AutoValue_Source) initFields_AutoValue_Source() {
	_ = this_
	this_.sourceMap = j2gort.NullString
	this_.sourceUrl = j2gort.NullString
	this_.sourceMappingUrl = j2gort.NullString
}

// Java constructor: AutoValue_Source(java.nio.file.Path path, java.lang.String sourceMap, java.lang.String sourceUrl, java.lang.String sourceMappingUrl, com.google.common.collect.ImmutableSet<java.lang.String> runtimes, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> loadFlags, int estimatedSize, com.google.javascript.jscomp.bundle.Source.Lazy<java.lang.String> codeSupplier, com.google.javascript.jscomp.bundle.Source.@org.jspecify.annotations.Nullable Lazy<java.lang.String> originalCodeSupplier)
func NewAutoValue_Source(path j2gort.Ext_java_nio_file_Path, sourceMap string, sourceUrl string, sourceMappingUrl string, runtimes j2gort.Ext_com_google_common_collect_ImmutableSet, loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap, estimatedSize int32, codeSupplier Source_Lazy, originalCodeSupplier Source_Lazy) *AutoValue_Source {
	this_ := &AutoValue_Source{}
	this_.BindSelf_(this_)
	InitAutoValue_Source(this_, path, sourceMap, sourceUrl, sourceMappingUrl, runtimes, loadFlags, estimatedSize, codeSupplier, originalCodeSupplier)
	return this_
}

func InitAutoValue_Source(this_ *AutoValue_Source, path j2gort.Ext_java_nio_file_Path, sourceMap string, sourceUrl string, sourceMappingUrl string, runtimes j2gort.Ext_com_google_common_collect_ImmutableSet, loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap, estimatedSize int32, codeSupplier Source_Lazy, originalCodeSupplier Source_Lazy) {
	_ = path
	_ = sourceMap
	_ = sourceUrl
	_ = sourceMappingUrl
	_ = runtimes
	_ = loadFlags
	_ = estimatedSize
	_ = codeSupplier
	_ = originalCodeSupplier
	InitSource(&this_.SourceBase)
	this_.initFields_AutoValue_Source()
	this_.path = path
	this_.sourceMap = sourceMap
	this_.sourceUrl = sourceUrl
	this_.sourceMappingUrl = sourceMappingUrl
	this_.runtimes = runtimes
	this_.loadFlags = loadFlags
	this_.estimatedSize = estimatedSize
	this_.codeSupplier = codeSupplier
	this_.originalCodeSupplier = originalCodeSupplier
}

// Java method: AutoValue_Source.path()
func (this_ *AutoValue_Source) Path() j2gort.Ext_java_nio_file_Path {
	return this_.path
}

// Java method: AutoValue_Source.sourceMap()
func (this_ *AutoValue_Source) SourceMap() string {
	return this_.sourceMap
}

// Java method: AutoValue_Source.sourceUrl()
func (this_ *AutoValue_Source) SourceUrl() string {
	return this_.sourceUrl
}

// Java method: AutoValue_Source.sourceMappingUrl()
func (this_ *AutoValue_Source) SourceMappingUrl() string {
	return this_.sourceMappingUrl
}

// Java method: AutoValue_Source.runtimes()
func (this_ *AutoValue_Source) Runtimes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.runtimes
}

// Java method: AutoValue_Source.loadFlags()
func (this_ *AutoValue_Source) LoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.loadFlags
}

// Java method: AutoValue_Source.estimatedSize()
func (this_ *AutoValue_Source) EstimatedSize() int32 {
	return this_.estimatedSize
}

// Java method: AutoValue_Source.codeSupplier()
func (this_ *AutoValue_Source) CodeSupplier() Source_Lazy {
	return this_.codeSupplier
}

// Java method: AutoValue_Source.originalCodeSupplier()
func (this_ *AutoValue_Source) OriginalCodeSupplier() Source_Lazy {
	return this_.originalCodeSupplier
}

// Java method: AutoValue_Source.toString()
func (this_ *AutoValue_Source) String() string {
	return ("Source{path=" + j2gort.StrVal(this_.path) + ", sourceMap=" + j2gort.StrValString(this_.sourceMap) + ", sourceUrl=" + j2gort.StrValString(this_.sourceUrl) + ", sourceMappingUrl=" + j2gort.StrValString(this_.sourceMappingUrl) + ", runtimes=" + j2gort.StrVal(this_.runtimes) + ", loadFlags=" + j2gort.StrVal(this_.loadFlags) + ", estimatedSize=" + j2gort.StrValInt32(this_.estimatedSize) + ", codeSupplier=" + j2gort.StrVal(this_.codeSupplier) + ", originalCodeSupplier=" + j2gort.StrVal(this_.originalCodeSupplier) + "}")
}

// Java method: AutoValue_Source.equals(java.lang.@org.jspecify.annotations.Nullable Object o)
func (this_ *AutoValue_Source) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[Source](o) {
		var that Source = j2gort.AsIface[Source](o)
		_ = that
		return ((((((((this_.path.Equals(that.Path()) && func() bool {
			sr_2 := this_.sourceMap
			se_3 := that.SourceMap()
			return j2gort.CheckNotNullString(sr_2) == se_3
		}()) && func() bool {
			sr_4 := this_.sourceUrl
			se_5 := that.SourceUrl()
			return j2gort.CheckNotNullString(sr_4) == se_5
		}()) && func() bool {
			sr_6 := this_.sourceMappingUrl
			se_7 := that.SourceMappingUrl()
			return j2gort.CheckNotNullString(sr_6) == se_7
		}()) && this_.runtimes.Equals(that.Runtimes())) && this_.loadFlags.Equals(that.LoadFlags())) && (this_.estimatedSize == that.EstimatedSize())) && this_.codeSupplier.Equals(that.CodeSupplier())) && (func() bool {
			if this_.originalCodeSupplier == nil || this_.originalCodeSupplier.IsJavaNil_() {
				return j2gort.IfaceEquals[Source_Lazy](that.OriginalCodeSupplier(), nil)
			}
			return this_.originalCodeSupplier.Equals(that.OriginalCodeSupplier())
		}()))
	}
	return false
}

// Java method: AutoValue_Source.hashCode()
func (this_ *AutoValue_Source) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.path.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.sourceMap)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.sourceUrl)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ javalang.StringHashCode(j2gort.CheckNotNullString(this_.sourceMappingUrl)))
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.runtimes.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.loadFlags.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.estimatedSize)
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.codeSupplier.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.originalCodeSupplier == nil || this_.originalCodeSupplier.IsJavaNil_() {
			return int32(0)
		}
		return this_.originalCodeSupplier.HashCode()
	}())
	return h_d
}

// Java method: AutoValue_Source.toBuilder()
func (this_ *AutoValue_Source) ToBuilder() Source_Builder {
	return NewAutoValue_Source_Builder_1(this_)
}

// Java class: com.google.javascript.jscomp.bundle.AutoValue_Source.Builder
type AutoValue_Source_Builder struct {
	Source_BuilderBase
	path                 j2gort.Ext_java_nio_file_Path                     `gson:"path"`
	sourceMap            string                                            `gson:"sourceMap"`
	sourceUrl            string                                            `gson:"sourceUrl"`
	sourceMappingUrl     string                                            `gson:"sourceMappingUrl"`
	runtimes             j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"runtimes"`
	loadFlags            j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"loadFlags"`
	codeSupplier         Source_Lazy                                       `gson:"codeSupplier"`
	originalCodeSupplier Source_Lazy                                       `gson:"originalCodeSupplier"`
	estimatedSize        int32                                             `gson:"estimatedSize"`
	set_d0               int8                                              `gson:"set$0"`
}

func (this_ *AutoValue_Source_Builder) BindSelf_(s any) {
	_ = s
	this_.Source_BuilderBase.BindSelf_(s)
}
func (this_ *AutoValue_Source_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_Source_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.AutoValue_Source$Builder", "Builder"
}
func (this_ *AutoValue_Source_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Source_BuilderBase.IdentityHash_()
}
func (*AutoValue_Source_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.AutoValue_Source.Builder", "com.google.javascript.jscomp.bundle.AutoValue_Source$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_Source_Builder](v_)
	})
}
func (*AutoValue_Source_Builder) Implements_jscomp_Source_Builder() {}

func (this_ *AutoValue_Source_Builder) initFields_AutoValue_Source_Builder() {
	_ = this_
	this_.sourceMap = j2gort.NullString
	this_.sourceUrl = j2gort.NullString
	this_.sourceMappingUrl = j2gort.NullString
}

// Java constructor: Builder()
func NewAutoValue_Source_Builder() *AutoValue_Source_Builder {
	this_ := &AutoValue_Source_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_Source_Builder(this_)
	return this_
}

func InitAutoValue_Source_Builder(this_ *AutoValue_Source_Builder) {
	InitSource_Builder(&this_.Source_BuilderBase)
	this_.initFields_AutoValue_Source_Builder()
}

// Java constructor: Builder(com.google.javascript.jscomp.bundle.Source source)
func NewAutoValue_Source_Builder_1(source Source) *AutoValue_Source_Builder {
	this_ := &AutoValue_Source_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_Source_Builder_1(this_, source)
	return this_
}

func InitAutoValue_Source_Builder_1(this_ *AutoValue_Source_Builder, source Source) {
	_ = source
	InitSource_Builder(&this_.Source_BuilderBase)
	this_.initFields_AutoValue_Source_Builder()
	this_.path = source.Path()
	this_.sourceMap = source.SourceMap()
	this_.sourceUrl = source.SourceUrl()
	this_.sourceMappingUrl = source.SourceMappingUrl()
	this_.runtimes = source.Runtimes()
	this_.loadFlags = source.LoadFlags()
	this_.estimatedSize = source.EstimatedSize()
	this_.codeSupplier = source.CodeSupplier()
	this_.originalCodeSupplier = source.OriginalCodeSupplier()
	this_.set_d0 = int8(1)
}

// Java method: Builder.setPath(java.nio.file.Path path)
func (this_ *AutoValue_Source_Builder) SetPath(path j2gort.Ext_java_nio_file_Path) Source_Builder {
	_ = path
	if path == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null path")))
	}
	this_.path = path
	return this_
}

// Java method: Builder.setSourceMap(java.lang.String sourceMap)
func (this_ *AutoValue_Source_Builder) SetSourceMap(sourceMap string) Source_Builder {
	_ = sourceMap
	if sourceMap == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null sourceMap")))
	}
	this_.sourceMap = sourceMap
	return this_
}

// Java method: Builder.setSourceUrl(java.lang.String sourceUrl)
func (this_ *AutoValue_Source_Builder) SetSourceUrl(sourceUrl string) Source_Builder {
	_ = sourceUrl
	if sourceUrl == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null sourceUrl")))
	}
	this_.sourceUrl = sourceUrl
	return this_
}

// Java method: Builder.setSourceMappingUrl(java.lang.String sourceMappingUrl)
func (this_ *AutoValue_Source_Builder) SetSourceMappingUrl(sourceMappingUrl string) Source_Builder {
	_ = sourceMappingUrl
	if sourceMappingUrl == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null sourceMappingUrl")))
	}
	this_.sourceMappingUrl = sourceMappingUrl
	return this_
}

// Java method: Builder.setRuntimes(com.google.common.collect.ImmutableSet<java.lang.String> runtimes)
func (this_ *AutoValue_Source_Builder) SetRuntimes(runtimes j2gort.Ext_com_google_common_collect_ImmutableSet) Source_Builder {
	_ = runtimes
	if runtimes == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null runtimes")))
	}
	this_.runtimes = runtimes
	return this_
}

// Java method: Builder.runtimes()
func (this_ *AutoValue_Source_Builder) Runtimes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	if this_.runtimes == nil {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Property \"runtimes\" has not been set")))
	}
	return this_.runtimes
}

// Java method: Builder.setLoadFlags(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> loadFlags)
func (this_ *AutoValue_Source_Builder) SetLoadFlags(loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap) Source_Builder {
	_ = loadFlags
	if loadFlags == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null loadFlags")))
	}
	this_.loadFlags = loadFlags
	return this_
}

// Java method: Builder.setEstimatedSize(int estimatedSize)
func (this_ *AutoValue_Source_Builder) SetEstimatedSize(estimatedSize int32) Source_Builder {
	_ = estimatedSize
	this_.estimatedSize = estimatedSize
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.setCodeSupplier(com.google.javascript.jscomp.bundle.Source.Lazy<java.lang.String> codeSupplier)
func (this_ *AutoValue_Source_Builder) SetCodeSupplier(codeSupplier Source_Lazy) Source_Builder {
	_ = codeSupplier
	if codeSupplier == nil || codeSupplier.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null codeSupplier")))
	}
	this_.codeSupplier = codeSupplier
	return this_
}

// Java method: Builder.codeSupplier()
func (this_ *AutoValue_Source_Builder) CodeSupplier() Source_Lazy {
	if this_.codeSupplier == nil || this_.codeSupplier.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Property \"codeSupplier\" has not been set")))
	}
	return this_.codeSupplier
}

// Java method: Builder.setOriginalCodeSupplier(com.google.javascript.jscomp.bundle.Source.@org.jspecify.annotations.Nullable Lazy<java.lang.String> originalCodeSupplier)
func (this_ *AutoValue_Source_Builder) SetOriginalCodeSupplier(originalCodeSupplier Source_Lazy) Source_Builder {
	_ = originalCodeSupplier
	this_.originalCodeSupplier = originalCodeSupplier
	return this_
}

// Java method: Builder.originalCodeSupplier()
func (this_ *AutoValue_Source_Builder) OriginalCodeSupplier() Source_Lazy {
	return this_.originalCodeSupplier
}

// Java method: Builder.autoBuild()
func (this_ *AutoValue_Source_Builder) AutoBuild() Source {
	if (((((((int32(this_.set_d0) != int32(1)) || (this_.path == nil)) || (this_.sourceMap == j2gort.NullString)) || (this_.sourceUrl == j2gort.NullString)) || (this_.sourceMappingUrl == j2gort.NullString)) || (this_.runtimes == nil)) || (this_.loadFlags == nil)) || (this_.codeSupplier == nil || this_.codeSupplier.IsJavaNil_()) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.path == nil {
			missing.Append_1_string(" path")
		}
		if this_.sourceMap == j2gort.NullString {
			missing.Append_1_string(" sourceMap")
		}
		if this_.sourceUrl == j2gort.NullString {
			missing.Append_1_string(" sourceUrl")
		}
		if this_.sourceMappingUrl == j2gort.NullString {
			missing.Append_1_string(" sourceMappingUrl")
		}
		if this_.runtimes == nil {
			missing.Append_1_string(" runtimes")
		}
		if this_.loadFlags == nil {
			missing.Append_1_string(" loadFlags")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" estimatedSize")
		}
		if this_.codeSupplier == nil || this_.codeSupplier.IsJavaNil_() {
			missing.Append_1_string(" codeSupplier")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAutoValue_Source(this_.path, this_.sourceMap, this_.sourceUrl, this_.sourceMappingUrl, this_.runtimes, this_.loadFlags, this_.estimatedSize, this_.codeSupplier, this_.originalCodeSupplier)
}
