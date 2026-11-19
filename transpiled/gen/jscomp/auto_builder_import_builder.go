// Port of com/google/javascript/jscomp/modules/AutoBuilder_Import_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.AutoBuilder_Import_Builder
type AutoBuilder_Import_Builder struct {
	Import_BuilderBase
	moduleRequest string                   `gson:"moduleRequest"`
	importName    string                   `gson:"importName"`
	localName     string                   `gson:"localName"`
	modulePath    *ModuleLoader_ModulePath `gson:"modulePath"`
	importNode    pkg_rhino.Node           `gson:"importNode"`
	nameNode      pkg_rhino.Node           `gson:"nameNode"`
}

func (this_ *AutoBuilder_Import_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_Import_Builder) IsJavaNil_() bool { return this_ == nil }
func (*AutoBuilder_Import_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.AutoBuilder_Import_Builder", "AutoBuilder_Import_Builder"
}
func (this_ *AutoBuilder_Import_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Import_BuilderBase.IdentityHash_()
}
func (*AutoBuilder_Import_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.AutoBuilder_Import_Builder", "com.google.javascript.jscomp.modules.AutoBuilder_Import_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_Import_Builder](v_)
	})
}
func (*AutoBuilder_Import_Builder) Implements_jscomp_Import_Builder() {}

func (this_ *AutoBuilder_Import_Builder) initFields_AutoBuilder_Import_Builder() {
	_ = this_
	this_.moduleRequest = j2gort.NullString
	this_.importName = j2gort.NullString
	this_.localName = j2gort.NullString
}

// Java constructor: AutoBuilder_Import_Builder()
func NewAutoBuilder_Import_Builder() *AutoBuilder_Import_Builder {
	this_ := &AutoBuilder_Import_Builder{}
	InitAutoBuilder_Import_Builder(this_)
	return this_
}

func InitAutoBuilder_Import_Builder(this_ *AutoBuilder_Import_Builder) {
	InitImport_Builder(&this_.Import_BuilderBase)
	this_.initFields_AutoBuilder_Import_Builder()
}

// Java constructor: AutoBuilder_Import_Builder(com.google.javascript.jscomp.modules.Import source)
func NewAutoBuilder_Import_Builder_1(source *Import) *AutoBuilder_Import_Builder {
	this_ := &AutoBuilder_Import_Builder{}
	InitAutoBuilder_Import_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_Import_Builder_1(this_ *AutoBuilder_Import_Builder, source *Import) {
	_ = source
	InitImport_Builder(&this_.Import_BuilderBase)
	this_.initFields_AutoBuilder_Import_Builder()
	this_.moduleRequest = source.ModuleRequest()
	this_.importName = source.ImportName()
	this_.localName = source.LocalName()
	this_.modulePath = source.ModulePath()
	this_.importNode = source.ImportNode()
	this_.nameNode = source.NameNode()
}

// Java method: AutoBuilder_Import_Builder.moduleRequest(java.lang.String moduleRequest)
func (this_ *AutoBuilder_Import_Builder) ModuleRequest(moduleRequest string) Import_Builder {
	_ = moduleRequest
	if moduleRequest == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null moduleRequest")))
	}
	this_.moduleRequest = moduleRequest
	return this_
}

// Java method: AutoBuilder_Import_Builder.importName(java.lang.String importName)
func (this_ *AutoBuilder_Import_Builder) ImportName(importName string) Import_Builder {
	_ = importName
	if importName == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null importName")))
	}
	this_.importName = importName
	return this_
}

// Java method: AutoBuilder_Import_Builder.localName(java.lang.String localName)
func (this_ *AutoBuilder_Import_Builder) LocalName(localName string) Import_Builder {
	_ = localName
	if localName == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null localName")))
	}
	this_.localName = localName
	return this_
}

// Java method: AutoBuilder_Import_Builder.modulePath(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath modulePath)
func (this_ *AutoBuilder_Import_Builder) ModulePath(modulePath *ModuleLoader_ModulePath) Import_Builder {
	_ = modulePath
	this_.modulePath = modulePath
	return this_
}

// Java method: AutoBuilder_Import_Builder.importNode(com.google.javascript.rhino.Node importNode)
func (this_ *AutoBuilder_Import_Builder) ImportNode(importNode pkg_rhino.Node) Import_Builder {
	_ = importNode
	if importNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null importNode")))
	}
	this_.importNode = importNode
	return this_
}

// Java method: AutoBuilder_Import_Builder.nameNode(com.google.javascript.rhino.Node nameNode)
func (this_ *AutoBuilder_Import_Builder) NameNode(nameNode pkg_rhino.Node) Import_Builder {
	_ = nameNode
	if nameNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null nameNode")))
	}
	this_.nameNode = nameNode
	return this_
}

// Java method: AutoBuilder_Import_Builder.build()
func (this_ *AutoBuilder_Import_Builder) Build() *Import {
	if ((((this_.moduleRequest == j2gort.NullString) || (this_.importName == j2gort.NullString)) || (this_.localName == j2gort.NullString)) || (this_.importNode == 0)) || (this_.nameNode == 0) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.moduleRequest == j2gort.NullString {
			missing.Append_1_string(" moduleRequest")
		}
		if this_.importName == j2gort.NullString {
			missing.Append_1_string(" importName")
		}
		if this_.localName == j2gort.NullString {
			missing.Append_1_string(" localName")
		}
		if this_.importNode == 0 {
			missing.Append_1_string(" importNode")
		}
		if this_.nameNode == 0 {
			missing.Append_1_string(" nameNode")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewImport(this_.moduleRequest, this_.importName, this_.localName, this_.modulePath, this_.importNode, this_.nameNode)
}
