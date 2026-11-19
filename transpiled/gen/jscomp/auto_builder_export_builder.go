// Port of com/google/javascript/jscomp/modules/AutoBuilder_Export_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.AutoBuilder_Export_Builder
type AutoBuilder_Export_Builder struct {
	Export_BuilderBase
	exportName       string                           `gson:"exportName"`
	moduleRequest    string                           `gson:"moduleRequest"`
	importName       string                           `gson:"importName"`
	localName        string                           `gson:"localName"`
	modulePath       *ModuleLoader_ModulePath         `gson:"modulePath"`
	moduleMetadata   ModuleMetadataMap_ModuleMetadata `gson:"moduleMetadata"`
	closureNamespace string                           `gson:"closureNamespace"`
	exportNode       pkg_rhino.Node                   `gson:"exportNode"`
	nameNode         pkg_rhino.Node                   `gson:"nameNode"`
	mutated          bool                             `gson:"mutated"`
	set_d0           int8                             `gson:"set$0"`
}

func (this_ *AutoBuilder_Export_Builder) BindSelf_(s any) {
	_ = s
	this_.Export_BuilderBase.BindSelf_(s)
}
func (this_ *AutoBuilder_Export_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoBuilder_Export_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.AutoBuilder_Export_Builder", "AutoBuilder_Export_Builder"
}
func (this_ *AutoBuilder_Export_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Export_BuilderBase.IdentityHash_()
}
func (*AutoBuilder_Export_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.AutoBuilder_Export_Builder", "com.google.javascript.jscomp.modules.AutoBuilder_Export_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_Export_Builder](v_)
	})
}
func (*AutoBuilder_Export_Builder) Implements_jscomp_Export_Builder() {}

func (this_ *AutoBuilder_Export_Builder) initFields_AutoBuilder_Export_Builder() {
	_ = this_
	this_.exportName = j2gort.NullString
	this_.moduleRequest = j2gort.NullString
	this_.importName = j2gort.NullString
	this_.localName = j2gort.NullString
	this_.closureNamespace = j2gort.NullString
}

// Java constructor: AutoBuilder_Export_Builder()
func NewAutoBuilder_Export_Builder() *AutoBuilder_Export_Builder {
	this_ := &AutoBuilder_Export_Builder{}
	this_.BindSelf_(this_)
	InitAutoBuilder_Export_Builder(this_)
	return this_
}

func InitAutoBuilder_Export_Builder(this_ *AutoBuilder_Export_Builder) {
	InitExport_Builder(&this_.Export_BuilderBase)
	this_.initFields_AutoBuilder_Export_Builder()
}

// Java constructor: AutoBuilder_Export_Builder(com.google.javascript.jscomp.modules.Export source)
func NewAutoBuilder_Export_Builder_1(source *Export) *AutoBuilder_Export_Builder {
	this_ := &AutoBuilder_Export_Builder{}
	this_.BindSelf_(this_)
	InitAutoBuilder_Export_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_Export_Builder_1(this_ *AutoBuilder_Export_Builder, source *Export) {
	_ = source
	InitExport_Builder(&this_.Export_BuilderBase)
	this_.initFields_AutoBuilder_Export_Builder()
	this_.exportName = source.ExportName()
	this_.moduleRequest = source.ModuleRequest()
	this_.importName = source.ImportName()
	this_.localName = source.LocalName()
	this_.modulePath = source.ModulePath()
	this_.exportNode = source.ExportNode()
	this_.nameNode = source.NameNode()
	this_.moduleMetadata = source.ModuleMetadata()
	this_.closureNamespace = source.ClosureNamespace()
	this_.mutated = source.Mutated()
	this_.set_d0 = int8(1)
}

// Java method: AutoBuilder_Export_Builder.exportName(java.lang.@org.jspecify.annotations.Nullable String exportName)
func (this_ *AutoBuilder_Export_Builder) ExportName(exportName string) Export_Builder {
	_ = exportName
	this_.exportName = exportName
	return this_
}

// Java method: AutoBuilder_Export_Builder.moduleRequest(java.lang.@org.jspecify.annotations.Nullable String moduleRequest)
func (this_ *AutoBuilder_Export_Builder) ModuleRequest(moduleRequest string) Export_Builder {
	_ = moduleRequest
	this_.moduleRequest = moduleRequest
	return this_
}

// Java method: AutoBuilder_Export_Builder.importName(java.lang.@org.jspecify.annotations.Nullable String importName)
func (this_ *AutoBuilder_Export_Builder) ImportName(importName string) Export_Builder {
	_ = importName
	this_.importName = importName
	return this_
}

// Java method: AutoBuilder_Export_Builder.localName(java.lang.@org.jspecify.annotations.Nullable String localName)
func (this_ *AutoBuilder_Export_Builder) LocalName(localName string) Export_Builder {
	_ = localName
	this_.localName = localName
	return this_
}

// Java method: AutoBuilder_Export_Builder.modulePath(com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath modulePath)
func (this_ *AutoBuilder_Export_Builder) ModulePath(modulePath *ModuleLoader_ModulePath) Export_Builder {
	_ = modulePath
	this_.modulePath = modulePath
	return this_
}

// Java method: AutoBuilder_Export_Builder.exportNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node exportNode)
func (this_ *AutoBuilder_Export_Builder) ExportNode(exportNode pkg_rhino.Node) Export_Builder {
	_ = exportNode
	this_.exportNode = exportNode
	return this_
}

// Java method: AutoBuilder_Export_Builder.nameNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node nameNode)
func (this_ *AutoBuilder_Export_Builder) NameNode(nameNode pkg_rhino.Node) Export_Builder {
	_ = nameNode
	this_.nameNode = nameNode
	return this_
}

// Java method: AutoBuilder_Export_Builder.moduleMetadata(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata)
func (this_ *AutoBuilder_Export_Builder) ModuleMetadata(moduleMetadata ModuleMetadataMap_ModuleMetadata) Export_Builder {
	_ = moduleMetadata
	if moduleMetadata == nil || moduleMetadata.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null moduleMetadata")))
	}
	this_.moduleMetadata = moduleMetadata
	return this_
}

// Java method: AutoBuilder_Export_Builder.closureNamespace(java.lang.@org.jspecify.annotations.Nullable String closureNamespace)
func (this_ *AutoBuilder_Export_Builder) ClosureNamespace(closureNamespace string) Export_Builder {
	_ = closureNamespace
	this_.closureNamespace = closureNamespace
	return this_
}

// Java method: AutoBuilder_Export_Builder.mutated(boolean mutated)
func (this_ *AutoBuilder_Export_Builder) Mutated(mutated bool) Export_Builder {
	_ = mutated
	this_.mutated = mutated
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: AutoBuilder_Export_Builder.autoBuild()
func (this_ *AutoBuilder_Export_Builder) AutoBuild() *Export {
	if (int32(this_.set_d0) != int32(1)) || (this_.moduleMetadata == nil || this_.moduleMetadata.IsJavaNil_()) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.moduleMetadata == nil || this_.moduleMetadata.IsJavaNil_() {
			missing.Append_1_string(" moduleMetadata")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" mutated")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewExport(this_.exportName, this_.moduleRequest, this_.importName, this_.localName, this_.modulePath, this_.exportNode, this_.nameNode, this_.moduleMetadata, this_.closureNamespace, this_.mutated)
}
