// Port of com/google/javascript/jscomp/AutoBuilder_Es6RewriteClass_ClassProperty_Builder.java (generated mechanically by j2go).

package jscomp

import (
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassProperty_Builder
type AutoBuilder_Es6RewriteClass_ClassProperty_Builder struct {
	Es6RewriteClass_ClassProperty_BuilderBase
	propertyKey  string                                     `gson:"propertyKey"`
	jsDocInfo    *pkg_rhino.JSDocInfo                       `gson:"jsDocInfo"`
	propertyType *pkg_colors.Color                          `gson:"propertyType"`
	kind         Es6RewriteClass_ClassProperty_PropertyKind `gson:"kind"`
}

func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) IsJavaNil_() bool {
	return this_ == nil
}
func (*AutoBuilder_Es6RewriteClass_ClassProperty_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassProperty_Builder", "AutoBuilder_Es6RewriteClass_ClassProperty_Builder"
}
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Es6RewriteClass_ClassProperty_BuilderBase.IdentityHash_()
}
func (*AutoBuilder_Es6RewriteClass_ClassProperty_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassProperty_Builder", "com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassProperty_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_Es6RewriteClass_ClassProperty_Builder](v_)
	})
}
func (*AutoBuilder_Es6RewriteClass_ClassProperty_Builder) Implements_jscomp_Es6RewriteClass_ClassProperty_Builder() {
}

func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) initFields_AutoBuilder_Es6RewriteClass_ClassProperty_Builder() {
	_ = this_
	this_.propertyKey = j2gort.NullString
}

// Java constructor: AutoBuilder_Es6RewriteClass_ClassProperty_Builder()
func NewAutoBuilder_Es6RewriteClass_ClassProperty_Builder() *AutoBuilder_Es6RewriteClass_ClassProperty_Builder {
	this_ := &AutoBuilder_Es6RewriteClass_ClassProperty_Builder{}
	InitAutoBuilder_Es6RewriteClass_ClassProperty_Builder(this_)
	return this_
}

func InitAutoBuilder_Es6RewriteClass_ClassProperty_Builder(this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) {
	InitEs6RewriteClass_ClassProperty_Builder(&this_.Es6RewriteClass_ClassProperty_BuilderBase)
	this_.initFields_AutoBuilder_Es6RewriteClass_ClassProperty_Builder()
}

// Java constructor: AutoBuilder_Es6RewriteClass_ClassProperty_Builder(com.google.javascript.jscomp.Es6RewriteClass.ClassProperty source)
func NewAutoBuilder_Es6RewriteClass_ClassProperty_Builder_1(source *Es6RewriteClass_ClassProperty) *AutoBuilder_Es6RewriteClass_ClassProperty_Builder {
	this_ := &AutoBuilder_Es6RewriteClass_ClassProperty_Builder{}
	InitAutoBuilder_Es6RewriteClass_ClassProperty_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_Es6RewriteClass_ClassProperty_Builder_1(this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder, source *Es6RewriteClass_ClassProperty) {
	_ = source
	InitEs6RewriteClass_ClassProperty_Builder(&this_.Es6RewriteClass_ClassProperty_BuilderBase)
	this_.initFields_AutoBuilder_Es6RewriteClass_ClassProperty_Builder()
	this_.propertyKey = source.PropertyKey()
	this_.kind = source.Kind()
	this_.jsDocInfo = source.JsDocInfo()
	this_.propertyType = source.PropertyType()
}

// Java method: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.propertyKey(java.lang.String propertyKey)
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) PropertyKey(propertyKey string) Es6RewriteClass_ClassProperty_Builder {
	_ = propertyKey
	if propertyKey == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null propertyKey")))
	}
	this_.propertyKey = propertyKey
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.kind(com.google.javascript.jscomp.Es6RewriteClass.ClassProperty.PropertyKind kind)
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) Kind(kind Es6RewriteClass_ClassProperty_PropertyKind) Es6RewriteClass_ClassProperty_Builder {
	_ = kind
	if kind == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null kind")))
	}
	this_.kind = kind
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.jsDocInfo(com.google.javascript.rhino.JSDocInfo jsDocInfo)
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) JsDocInfo(jsDocInfo *pkg_rhino.JSDocInfo) Es6RewriteClass_ClassProperty_Builder {
	_ = jsDocInfo
	if jsDocInfo == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null jsDocInfo")))
	}
	this_.jsDocInfo = jsDocInfo
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.propertyType(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color propertyType)
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) PropertyType(propertyType *pkg_colors.Color) Es6RewriteClass_ClassProperty_Builder {
	_ = propertyType
	this_.propertyType = propertyType
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.build()
func (this_ *AutoBuilder_Es6RewriteClass_ClassProperty_Builder) Build() *Es6RewriteClass_ClassProperty {
	if ((this_.propertyKey == j2gort.NullString) || (this_.kind == 0)) || (this_.jsDocInfo == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.propertyKey == j2gort.NullString {
			missing.Append_1_string(" propertyKey")
		}
		if this_.kind == 0 {
			missing.Append_1_string(" kind")
		}
		if this_.jsDocInfo == nil {
			missing.Append_1_string(" jsDocInfo")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewEs6RewriteClass_ClassProperty(this_.propertyKey, this_.kind, this_.jsDocInfo, this_.propertyType)
}
