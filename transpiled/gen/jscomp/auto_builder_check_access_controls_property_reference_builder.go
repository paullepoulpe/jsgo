// Port of com/google/javascript/jscomp/AutoBuilder_CheckAccessControls_PropertyReference_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AutoBuilder_CheckAccessControls_PropertyReference_Builder
type AutoBuilder_CheckAccessControls_PropertyReference_Builder struct {
	name             string                                 `gson:"name"`
	receiverType     pkg_rhino.ObjectType                   `gson:"receiverType"`
	readableTypeName j2gort.Ext_java_util_function_Supplier `gson:"readableTypeName"`
	sourceNode       pkg_rhino.Node                         `gson:"sourceNode"`
	idHash_          atomic.Int32
	mutation         bool `gson:"mutation"`
	declaration      bool `gson:"declaration"`
	override         bool `gson:"override"`
	set_d0           int8 `gson:"set$0"`
}

func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) IsJavaNil_() bool {
	return this_ == nil
}
func (*AutoBuilder_CheckAccessControls_PropertyReference_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoBuilder_CheckAccessControls_PropertyReference_Builder", "AutoBuilder_CheckAccessControls_PropertyReference_Builder"
}
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AutoBuilder_CheckAccessControls_PropertyReference_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoBuilder_CheckAccessControls_PropertyReference_Builder", "com.google.javascript.jscomp.AutoBuilder_CheckAccessControls_PropertyReference_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_CheckAccessControls_PropertyReference_Builder](v_)
	})
}
func (*AutoBuilder_CheckAccessControls_PropertyReference_Builder) Implements_jscomp_CheckAccessControls_PropertyReference_Builder() {
}

func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) initFields_AutoBuilder_CheckAccessControls_PropertyReference_Builder() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: AutoBuilder_CheckAccessControls_PropertyReference_Builder()
func NewAutoBuilder_CheckAccessControls_PropertyReference_Builder() *AutoBuilder_CheckAccessControls_PropertyReference_Builder {
	this_ := &AutoBuilder_CheckAccessControls_PropertyReference_Builder{}
	InitAutoBuilder_CheckAccessControls_PropertyReference_Builder(this_)
	return this_
}

func InitAutoBuilder_CheckAccessControls_PropertyReference_Builder(this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) {
	this_.initFields_AutoBuilder_CheckAccessControls_PropertyReference_Builder()
}

// Java constructor: AutoBuilder_CheckAccessControls_PropertyReference_Builder(com.google.javascript.jscomp.CheckAccessControls.PropertyReference source)
func NewAutoBuilder_CheckAccessControls_PropertyReference_Builder_1(source *CheckAccessControls_PropertyReference) *AutoBuilder_CheckAccessControls_PropertyReference_Builder {
	this_ := &AutoBuilder_CheckAccessControls_PropertyReference_Builder{}
	InitAutoBuilder_CheckAccessControls_PropertyReference_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_CheckAccessControls_PropertyReference_Builder_1(this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder, source *CheckAccessControls_PropertyReference) {
	_ = source
	this_.initFields_AutoBuilder_CheckAccessControls_PropertyReference_Builder()
	this_.sourceNode = source.SourceNode()
	this_.name = source.Name()
	this_.receiverType = source.ReceiverType()
	this_.mutation = source.Mutation()
	this_.declaration = source.Declaration()
	this_.override = source.Override()
	this_.readableTypeName = source.ReadableTypeName()
	this_.set_d0 = int8(7)
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setSourceNode(com.google.javascript.rhino.Node sourceNode)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetSourceNode(sourceNode pkg_rhino.Node) CheckAccessControls_PropertyReference_Builder {
	_ = sourceNode
	if sourceNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null sourceNode")))
	}
	this_.sourceNode = sourceNode
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setName(java.lang.String name)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetName(name string) CheckAccessControls_PropertyReference_Builder {
	_ = name
	if name == j2gort.NullString {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null name")))
	}
	this_.name = name
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setReceiverType(com.google.javascript.rhino.jstype.ObjectType receiverType)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetReceiverType(receiverType pkg_rhino.ObjectType) CheckAccessControls_PropertyReference_Builder {
	_ = receiverType
	if receiverType == nil || receiverType.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null receiverType")))
	}
	this_.receiverType = receiverType
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setMutation(boolean mutation)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetMutation(mutation bool) CheckAccessControls_PropertyReference_Builder {
	_ = mutation
	this_.mutation = mutation
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setDeclaration(boolean declaration)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetDeclaration(declaration bool) CheckAccessControls_PropertyReference_Builder {
	_ = declaration
	this_.declaration = declaration
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(2))))
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setOverride(boolean override)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetOverride(override bool) CheckAccessControls_PropertyReference_Builder {
	_ = override
	this_.override = override
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(4))))
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.setReadableTypeName(java.util.function.Supplier<java.lang.String> readableTypeName)
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) SetReadableTypeName(readableTypeName j2gort.Ext_java_util_function_Supplier) CheckAccessControls_PropertyReference_Builder {
	_ = readableTypeName
	if readableTypeName == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null readableTypeName")))
	}
	this_.readableTypeName = readableTypeName
	return this_
}

// Java method: AutoBuilder_CheckAccessControls_PropertyReference_Builder.build()
func (this_ *AutoBuilder_CheckAccessControls_PropertyReference_Builder) Build() *CheckAccessControls_PropertyReference {
	if ((((int32(this_.set_d0) != int32(7)) || (this_.sourceNode == 0)) || (this_.name == j2gort.NullString)) || (this_.receiverType == nil || this_.receiverType.IsJavaNil_())) || (this_.readableTypeName == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.sourceNode == 0 {
			missing.Append_1_string(" sourceNode")
		}
		if this_.name == j2gort.NullString {
			missing.Append_1_string(" name")
		}
		if this_.receiverType == nil || this_.receiverType.IsJavaNil_() {
			missing.Append_1_string(" receiverType")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" mutation")
		}
		if (int32(this_.set_d0) & int32(2)) == int32(0) {
			missing.Append_1_string(" declaration")
		}
		if (int32(this_.set_d0) & int32(4)) == int32(0) {
			missing.Append_1_string(" override")
		}
		if this_.readableTypeName == nil {
			missing.Append_1_string(" readableTypeName")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewCheckAccessControls_PropertyReference(this_.sourceNode, this_.name, this_.receiverType, this_.mutation, this_.declaration, this_.override, this_.readableTypeName)
}
