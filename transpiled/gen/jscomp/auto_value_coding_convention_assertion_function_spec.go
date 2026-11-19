// Port of com/google/javascript/jscomp/AutoValue_CodingConvention_AssertionFunctionSpec.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec
type AutoValue_CodingConvention_AssertionFunctionSpec struct {
	CodingConvention_AssertionFunctionSpecBase
	functionName     string                                               `gson:"functionName"`
	paramIndex       int32                                                `gson:"paramIndex"`
	closurePrimitive pkg_rhino.ClosurePrimitive                           `gson:"closurePrimitive"`
	assertionKind    CodingConvention_AssertionFunctionSpec_AssertionKind `gson:"assertionKind"`
}

func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) BindSelf_(s any) {
	_ = s
	this_.CodingConvention_AssertionFunctionSpecBase.BindSelf_(s)
}
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_CodingConvention_AssertionFunctionSpec) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec", "AutoValue_CodingConvention_AssertionFunctionSpec"
}
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodingConvention_AssertionFunctionSpecBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec", "com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_CodingConvention_AssertionFunctionSpec](v_)
	})
}
func (*AutoValue_CodingConvention_AssertionFunctionSpec) Implements_jscomp_CodingConvention_AssertionFunctionSpec() {
}

func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) initFields_AutoValue_CodingConvention_AssertionFunctionSpec() {
	_ = this_
	this_.functionName = j2gort.NullString
}

// Java constructor: AutoValue_CodingConvention_AssertionFunctionSpec(java.lang.@org.jspecify.annotations.Nullable String functionName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable ClosurePrimitive closurePrimitive, com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind assertionKind, int paramIndex)
func NewAutoValue_CodingConvention_AssertionFunctionSpec(functionName string, closurePrimitive pkg_rhino.ClosurePrimitive, assertionKind CodingConvention_AssertionFunctionSpec_AssertionKind, paramIndex int32) *AutoValue_CodingConvention_AssertionFunctionSpec {
	this_ := &AutoValue_CodingConvention_AssertionFunctionSpec{}
	this_.BindSelf_(this_)
	InitAutoValue_CodingConvention_AssertionFunctionSpec(this_, functionName, closurePrimitive, assertionKind, paramIndex)
	return this_
}

func InitAutoValue_CodingConvention_AssertionFunctionSpec(this_ *AutoValue_CodingConvention_AssertionFunctionSpec, functionName string, closurePrimitive pkg_rhino.ClosurePrimitive, assertionKind CodingConvention_AssertionFunctionSpec_AssertionKind, paramIndex int32) {
	_ = functionName
	_ = closurePrimitive
	_ = assertionKind
	_ = paramIndex
	InitCodingConvention_AssertionFunctionSpec(&this_.CodingConvention_AssertionFunctionSpecBase)
	this_.initFields_AutoValue_CodingConvention_AssertionFunctionSpec()
	this_.functionName = functionName
	this_.closurePrimitive = closurePrimitive
	this_.assertionKind = assertionKind
	this_.paramIndex = paramIndex
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.getFunctionName()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) GetFunctionName() string {
	return this_.functionName
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.getClosurePrimitive()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) GetClosurePrimitive() pkg_rhino.ClosurePrimitive {
	return this_.closurePrimitive
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.getAssertionKind()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) GetAssertionKind() CodingConvention_AssertionFunctionSpec_AssertionKind {
	return this_.assertionKind
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.getParamIndex()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) GetParamIndex() int32 {
	return this_.paramIndex
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.toString()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) String() string {
	return ("AssertionFunctionSpec{functionName=" + j2gort.StrValString(this_.functionName) + ", closurePrimitive=" + j2gort.StrVal(this_.closurePrimitive) + ", assertionKind=" + j2gort.StrVal(this_.assertionKind) + ", paramIndex=" + j2gort.StrValInt32(this_.paramIndex) + "}")
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.equals(java.lang.@org.jspecify.annotations.Nullable Object o)
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[CodingConvention_AssertionFunctionSpec](o) {
		var that CodingConvention_AssertionFunctionSpec = j2gort.AsIface[CodingConvention_AssertionFunctionSpec](o)
		_ = that
		return ((((func() bool {
			if this_.functionName == j2gort.NullString {
				return (that.GetFunctionName() == j2gort.NullString)
			}
			return func() bool {
				sr_2 := this_.functionName
				se_3 := that.GetFunctionName()
				return j2gort.CheckNotNullString(sr_2) == se_3
			}()
		}()) && (func() bool {
			if this_.closurePrimitive == 0 {
				return (that.GetClosurePrimitive() == 0)
			}
			return func() bool {
				re_4 := this_.closurePrimitive
				ae_5 := that.GetClosurePrimitive()
				if re_4 == 0 {
					j2gort.ThrowNPE()
				}
				return re_4 == ae_5
			}()
		}())) && func() bool {
			re_6 := this_.assertionKind
			ae_7 := that.GetAssertionKind()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			return re_6 == ae_7
		}()) && (this_.paramIndex == that.GetParamIndex()))
	}
	return false
}

// Java method: AutoValue_CodingConvention_AssertionFunctionSpec.hashCode()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.functionName == j2gort.NullString {
			return int32(0)
		}
		return javalang.StringHashCode(j2gort.CheckNotNullString(this_.functionName))
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.closurePrimitive == 0 {
			return int32(0)
		}
		return j2gort.CheckNotNullEnum(this_.closurePrimitive).HashCode()
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ j2gort.CheckNotNullEnum(this_.assertionKind).HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.paramIndex)
	return h_d
}

// Java class: com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec.Builder
type AutoValue_CodingConvention_AssertionFunctionSpec_Builder struct {
	CodingConvention_AssertionFunctionSpec_BuilderBase
	functionName     string                                               `gson:"functionName"`
	paramIndex       int32                                                `gson:"paramIndex"`
	closurePrimitive pkg_rhino.ClosurePrimitive                           `gson:"closurePrimitive"`
	assertionKind    CodingConvention_AssertionFunctionSpec_AssertionKind `gson:"assertionKind"`
	set_d0           int8                                                 `gson:"set$0"`
}

func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) BindSelf_(s any) {
	_ = s
	this_.CodingConvention_AssertionFunctionSpec_BuilderBase.BindSelf_(s)
}
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) IsJavaNil_() bool {
	return this_ == nil
}
func (*AutoValue_CodingConvention_AssertionFunctionSpec_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec$Builder", "Builder"
}
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodingConvention_AssertionFunctionSpec_BuilderBase.IdentityHash_()
}
func (*AutoValue_CodingConvention_AssertionFunctionSpec_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec.Builder", "com.google.javascript.jscomp.AutoValue_CodingConvention_AssertionFunctionSpec$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_CodingConvention_AssertionFunctionSpec_Builder](v_)
	})
}
func (*AutoValue_CodingConvention_AssertionFunctionSpec_Builder) Implements_jscomp_CodingConvention_AssertionFunctionSpec_Builder() {
}

func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) initFields_AutoValue_CodingConvention_AssertionFunctionSpec_Builder() {
	_ = this_
	this_.functionName = j2gort.NullString
}

// Java constructor: Builder()
func NewAutoValue_CodingConvention_AssertionFunctionSpec_Builder() *AutoValue_CodingConvention_AssertionFunctionSpec_Builder {
	this_ := &AutoValue_CodingConvention_AssertionFunctionSpec_Builder{}
	this_.BindSelf_(this_)
	InitAutoValue_CodingConvention_AssertionFunctionSpec_Builder(this_)
	return this_
}

func InitAutoValue_CodingConvention_AssertionFunctionSpec_Builder(this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) {
	InitCodingConvention_AssertionFunctionSpec_Builder(&this_.CodingConvention_AssertionFunctionSpec_BuilderBase)
	this_.initFields_AutoValue_CodingConvention_AssertionFunctionSpec_Builder()
}

// Java method: Builder.setFunctionName(java.lang.String functionName)
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) SetFunctionName(functionName string) CodingConvention_AssertionFunctionSpec_Builder {
	_ = functionName
	this_.functionName = functionName
	return this_
}

// Java method: Builder.setClosurePrimitive(com.google.javascript.rhino.ClosurePrimitive closurePrimitive)
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) SetClosurePrimitive(closurePrimitive pkg_rhino.ClosurePrimitive) CodingConvention_AssertionFunctionSpec_Builder {
	_ = closurePrimitive
	this_.closurePrimitive = closurePrimitive
	return this_
}

// Java method: Builder.setAssertionKind(com.google.javascript.jscomp.CodingConvention.AssertionFunctionSpec.AssertionKind assertionKind)
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) SetAssertionKind(assertionKind CodingConvention_AssertionFunctionSpec_AssertionKind) CodingConvention_AssertionFunctionSpec_Builder {
	_ = assertionKind
	if assertionKind == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null assertionKind")))
	}
	this_.assertionKind = assertionKind
	return this_
}

// Java method: Builder.setParamIndex(int paramIndex)
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) SetParamIndex(paramIndex int32) CodingConvention_AssertionFunctionSpec_Builder {
	_ = paramIndex
	this_.paramIndex = paramIndex
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.autoBuild()
func (this_ *AutoValue_CodingConvention_AssertionFunctionSpec_Builder) AutoBuild() CodingConvention_AssertionFunctionSpec {
	if (int32(this_.set_d0) != int32(1)) || (this_.assertionKind == 0) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.assertionKind == 0 {
			missing.Append_1_string(" assertionKind")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" paramIndex")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewAutoValue_CodingConvention_AssertionFunctionSpec(this_.functionName, this_.closurePrimitive, this_.assertionKind, this_.paramIndex)
}
