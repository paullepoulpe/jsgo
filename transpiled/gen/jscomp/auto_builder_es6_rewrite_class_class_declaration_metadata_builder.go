// Port of com/google/javascript/jscomp/AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder
type AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder struct {
	Es6RewriteClass_ClassDeclarationMetadata_BuilderBase
	insertionPoint                  *Es6RewriteClass_InsertionPoint `gson:"insertionPoint"`
	classMembersToDeclare           j2gort.Ext_java_util_Map        `gson:"classMembersToDeclare"`
	definePropertiesObjForPrototype pkg_rhino.Node                  `gson:"definePropertiesObjForPrototype"`
	definePropertiesObjForClass     pkg_rhino.Node                  `gson:"definePropertiesObjForClass"`
	fullClassNameNode               pkg_rhino.Node                  `gson:"fullClassNameNode"`
	classPrototypeNode              pkg_rhino.Node                  `gson:"classPrototypeNode"`
	classNameNode                   pkg_rhino.Node                  `gson:"classNameNode"`
	superClassNameNode              pkg_rhino.Node                  `gson:"superClassNameNode"`
	anonymous                       bool                            `gson:"anonymous"`
	set_d0                          int8                            `gson:"set$0"`
}

func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) IsJavaNil_() bool {
	return this_ == nil
}
func (*AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder", "AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder"
}
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Es6RewriteClass_ClassDeclarationMetadata_BuilderBase.IdentityHash_()
}
func (*AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder", "com.google.javascript.jscomp.AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder](v_)
	})
}
func (*AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) Implements_jscomp_Es6RewriteClass_ClassDeclarationMetadata_Builder() {
}

func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) initFields_AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder() {
	_ = this_
}

// Java constructor: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder()
func NewAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder() *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder {
	this_ := &AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder{}
	InitAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder(this_)
	return this_
}

func InitAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder(this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) {
	InitEs6RewriteClass_ClassDeclarationMetadata_Builder(&this_.Es6RewriteClass_ClassDeclarationMetadata_BuilderBase)
	this_.initFields_AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder()
}

// Java constructor: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder(com.google.javascript.jscomp.Es6RewriteClass.ClassDeclarationMetadata source)
func NewAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder_1(source *Es6RewriteClass_ClassDeclarationMetadata) *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder {
	this_ := &AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder{}
	InitAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder_1(this_, source)
	return this_
}

func InitAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder_1(this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder, source *Es6RewriteClass_ClassDeclarationMetadata) {
	_ = source
	InitEs6RewriteClass_ClassDeclarationMetadata_Builder(&this_.Es6RewriteClass_ClassDeclarationMetadata_BuilderBase)
	this_.initFields_AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder()
	this_.insertionPoint = source.InsertionPoint()
	this_.definePropertiesObjForPrototype = source.DefinePropertiesObjForPrototype()
	this_.definePropertiesObjForClass = source.DefinePropertiesObjForClass()
	this_.classMembersToDeclare = source.ClassMembersToDeclare()
	this_.fullClassNameNode = source.FullClassNameNode()
	this_.classPrototypeNode = source.ClassPrototypeNode()
	this_.anonymous = source.Anonymous()
	this_.classNameNode = source.ClassNameNode()
	this_.superClassNameNode = source.SuperClassNameNode()
	this_.set_d0 = int8(1)
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setInsertionPoint(com.google.javascript.jscomp.Es6RewriteClass.InsertionPoint insertionPoint)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetInsertionPoint(insertionPoint *Es6RewriteClass_InsertionPoint) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = insertionPoint
	if insertionPoint == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null insertionPoint")))
	}
	this_.insertionPoint = insertionPoint
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setDefinePropertiesObjForPrototype(com.google.javascript.rhino.Node definePropertiesObjForPrototype)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetDefinePropertiesObjForPrototype(definePropertiesObjForPrototype pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = definePropertiesObjForPrototype
	if definePropertiesObjForPrototype == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null definePropertiesObjForPrototype")))
	}
	this_.definePropertiesObjForPrototype = definePropertiesObjForPrototype
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setDefinePropertiesObjForClass(com.google.javascript.rhino.Node definePropertiesObjForClass)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetDefinePropertiesObjForClass(definePropertiesObjForClass pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = definePropertiesObjForClass
	if definePropertiesObjForClass == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null definePropertiesObjForClass")))
	}
	this_.definePropertiesObjForClass = definePropertiesObjForClass
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setClassMembersToDeclare(java.util.Map<java.lang.String,com.google.javascript.jscomp.Es6RewriteClass.ClassProperty> classMembersToDeclare)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetClassMembersToDeclare(classMembersToDeclare j2gort.Ext_java_util_Map) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = classMembersToDeclare
	if classMembersToDeclare == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null classMembersToDeclare")))
	}
	this_.classMembersToDeclare = classMembersToDeclare
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setFullClassNameNode(com.google.javascript.rhino.Node fullClassNameNode)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetFullClassNameNode(fullClassNameNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = fullClassNameNode
	if fullClassNameNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null fullClassNameNode")))
	}
	this_.fullClassNameNode = fullClassNameNode
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.fullClassNameNode()
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) FullClassNameNode() pkg_rhino.Node {
	if this_.fullClassNameNode == 0 {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Property \"fullClassNameNode\" has not been set")))
	}
	return this_.fullClassNameNode
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setClassPrototypeNode(com.google.javascript.rhino.Node classPrototypeNode)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetClassPrototypeNode(classPrototypeNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = classPrototypeNode
	if classPrototypeNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null classPrototypeNode")))
	}
	this_.classPrototypeNode = classPrototypeNode
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setAnonymous(boolean anonymous)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetAnonymous(anonymous bool) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = anonymous
	this_.anonymous = anonymous
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setClassNameNode(com.google.javascript.rhino.Node classNameNode)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetClassNameNode(classNameNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = classNameNode
	if classNameNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null classNameNode")))
	}
	this_.classNameNode = classNameNode
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setSuperClassNameNode(com.google.javascript.rhino.Node superClassNameNode)
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) SetSuperClassNameNode(superClassNameNode pkg_rhino.Node) Es6RewriteClass_ClassDeclarationMetadata_Builder {
	_ = superClassNameNode
	if superClassNameNode == 0 {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null superClassNameNode")))
	}
	this_.superClassNameNode = superClassNameNode
	return this_
}

// Java method: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.build()
func (this_ *AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder) Build() *Es6RewriteClass_ClassDeclarationMetadata {
	if ((((((((int32(this_.set_d0) != int32(1)) || (this_.insertionPoint == nil)) || (this_.definePropertiesObjForPrototype == 0)) || (this_.definePropertiesObjForClass == 0)) || (this_.classMembersToDeclare == nil)) || (this_.fullClassNameNode == 0)) || (this_.classPrototypeNode == 0)) || (this_.classNameNode == 0)) || (this_.superClassNameNode == 0) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.insertionPoint == nil {
			missing.Append_1_string(" insertionPoint")
		}
		if this_.definePropertiesObjForPrototype == 0 {
			missing.Append_1_string(" definePropertiesObjForPrototype")
		}
		if this_.definePropertiesObjForClass == 0 {
			missing.Append_1_string(" definePropertiesObjForClass")
		}
		if this_.classMembersToDeclare == nil {
			missing.Append_1_string(" classMembersToDeclare")
		}
		if this_.fullClassNameNode == 0 {
			missing.Append_1_string(" fullClassNameNode")
		}
		if this_.classPrototypeNode == 0 {
			missing.Append_1_string(" classPrototypeNode")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" anonymous")
		}
		if this_.classNameNode == 0 {
			missing.Append_1_string(" classNameNode")
		}
		if this_.superClassNameNode == 0 {
			missing.Append_1_string(" superClassNameNode")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewEs6RewriteClass_ClassDeclarationMetadata(this_.insertionPoint, this_.definePropertiesObjForPrototype, this_.definePropertiesObjForClass, this_.classMembersToDeclare, this_.fullClassNameNode, this_.classPrototypeNode, this_.anonymous, this_.classNameNode, this_.superClassNameNode)
}
