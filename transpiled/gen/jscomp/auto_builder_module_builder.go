// Port of com/google/javascript/jscomp/modules/AutoBuilder_Module_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.AutoBuilder_Module_Builder
type AutoBuilder_Module_Builder struct {
	Module_BuilderBase
	metadata               ModuleMetadataMap_ModuleMetadata                  `gson:"metadata"`
	path                   *ModuleLoader_ModulePath                          `gson:"path"`
	namespace              j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"namespace"`
	boundNames             j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"boundNames"`
	localNameToLocalExport j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"localNameToLocalExport"`
	closureNamespace       string                                            `gson:"closureNamespace"`
}

func (this_ *AutoBuilder_Module_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_Module_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoBuilder_Module_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.AutoBuilder_Module_Builder", "AutoBuilder_Module_Builder"
}
func (this_ *AutoBuilder_Module_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Module_BuilderBase.IdentityHash_()
}
func (*AutoBuilder_Module_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.AutoBuilder_Module_Builder", "com.google.javascript.jscomp.modules.AutoBuilder_Module_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_Module_Builder](v_)
	})
}
func (*AutoBuilder_Module_Builder) Implements_jscomp_Module_Builder() {}

func (this_ *AutoBuilder_Module_Builder) initFields_AutoBuilder_Module_Builder() {
	_ = this_
	this_.closureNamespace = j2gort.NullString
}

// Java constructor: AutoBuilder_Module_Builder()
func NewAutoBuilder_Module_Builder() *AutoBuilder_Module_Builder {
	this_ := &AutoBuilder_Module_Builder{}
	InitAutoBuilder_Module_Builder(this_)
	return this_
}

func InitAutoBuilder_Module_Builder(this_ *AutoBuilder_Module_Builder) {
	InitModule_Builder(&this_.Module_BuilderBase)
	this_.initFields_AutoBuilder_Module_Builder()
}

// Java constructor: AutoBuilder_Module_Builder(com.google.javascript.jscomp.modules.Module source)
func NewAutoBuilder_Module_Builder_1(source *Module) *AutoBuilder_Module_Builder {
	this_ := &AutoBuilder_Module_Builder{}
	InitAutoBuilder_Module_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_Module_Builder_1(this_ *AutoBuilder_Module_Builder, source *Module) {
	_ = source
	InitModule_Builder(&this_.Module_BuilderBase)
	this_.initFields_AutoBuilder_Module_Builder()
	this_.metadata = source.Metadata()
	this_.path = source.Path()
	this_.namespace = source.Namespace()
	this_.boundNames = source.BoundNames()
	this_.localNameToLocalExport = source.LocalNameToLocalExport()
	this_.closureNamespace = source.ClosureNamespace()
}

// Java method: AutoBuilder_Module_Builder.metadata(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func (this_ *AutoBuilder_Module_Builder) Metadata(metadata ModuleMetadataMap_ModuleMetadata) Module_Builder {
	_ = metadata
	if metadata == nil || metadata.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null metadata")))
	}
	this_.metadata = metadata
	return this_
}

// Java method: AutoBuilder_Module_Builder.path(com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath path)
func (this_ *AutoBuilder_Module_Builder) Path(path *ModuleLoader_ModulePath) Module_Builder {
	_ = path
	this_.path = path
	return this_
}

// Java method: AutoBuilder_Module_Builder.namespace(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> namespace)
func (this_ *AutoBuilder_Module_Builder) Namespace(namespace j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder {
	_ = namespace
	if namespace == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null namespace")))
	}
	this_.namespace = namespace
	return this_
}

// Java method: AutoBuilder_Module_Builder.boundNames(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Binding> boundNames)
func (this_ *AutoBuilder_Module_Builder) BoundNames(boundNames j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder {
	_ = boundNames
	if boundNames == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null boundNames")))
	}
	this_.boundNames = boundNames
	return this_
}

// Java method: AutoBuilder_Module_Builder.localNameToLocalExport(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Export> localNameToLocalExport)
func (this_ *AutoBuilder_Module_Builder) LocalNameToLocalExport(localNameToLocalExport j2gort.Ext_com_google_common_collect_ImmutableMap) Module_Builder {
	_ = localNameToLocalExport
	if localNameToLocalExport == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null localNameToLocalExport")))
	}
	this_.localNameToLocalExport = localNameToLocalExport
	return this_
}

// Java method: AutoBuilder_Module_Builder.closureNamespace(java.lang.@org.jspecify.annotations.Nullable String closureNamespace)
func (this_ *AutoBuilder_Module_Builder) ClosureNamespace(closureNamespace string) Module_Builder {
	_ = closureNamespace
	this_.closureNamespace = closureNamespace
	return this_
}

// Java method: AutoBuilder_Module_Builder.build()
func (this_ *AutoBuilder_Module_Builder) Build() *Module {
	if (((this_.metadata == nil || this_.metadata.IsJavaNil_()) || (this_.namespace == nil)) || (this_.boundNames == nil)) || (this_.localNameToLocalExport == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.metadata == nil || this_.metadata.IsJavaNil_() {
			missing.Append_1_string(" metadata")
		}
		if this_.namespace == nil {
			missing.Append_1_string(" namespace")
		}
		if this_.boundNames == nil {
			missing.Append_1_string(" boundNames")
		}
		if this_.localNameToLocalExport == nil {
			missing.Append_1_string(" localNameToLocalExport")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewModule(this_.metadata, this_.path, this_.namespace, this_.boundNames, this_.localNameToLocalExport, this_.closureNamespace)
}
