// Port of com/google/javascript/rhino/jstype/AutoValue_JSTypeRegistry_ClosureNamespace.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.AutoValue_JSTypeRegistry_ClosureNamespace
type AutoValue_JSTypeRegistry_ClosureNamespace struct {
	JSTypeRegistry_ClosureNamespaceBase
	type_          JSType `gson:"type"`
	definitionNode Node   `gson:"definitionNode"`
	isLegacy       bool   `gson:"isLegacy"`
}

func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_JSTypeRegistry_ClosureNamespace) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.AutoValue_JSTypeRegistry_ClosureNamespace", "AutoValue_JSTypeRegistry_ClosureNamespace"
}
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JSTypeRegistry_ClosureNamespaceBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.AutoValue_JSTypeRegistry_ClosureNamespace", "com.google.javascript.rhino.jstype.AutoValue_JSTypeRegistry_ClosureNamespace", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_JSTypeRegistry_ClosureNamespace](v_)
	})
}
func (*AutoValue_JSTypeRegistry_ClosureNamespace) Implements_rhino_JSTypeRegistry_ClosureNamespace() {
}

func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) initFields_AutoValue_JSTypeRegistry_ClosureNamespace() {
	_ = this_
}

// Java constructor: AutoValue_JSTypeRegistry_ClosureNamespace(boolean isLegacy, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node definitionNode, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func NewAutoValue_JSTypeRegistry_ClosureNamespace(isLegacy bool, definitionNode Node, type_ JSType) *AutoValue_JSTypeRegistry_ClosureNamespace {
	this_ := &AutoValue_JSTypeRegistry_ClosureNamespace{}
	InitAutoValue_JSTypeRegistry_ClosureNamespace(this_, isLegacy, definitionNode, type_)
	return this_
}

func InitAutoValue_JSTypeRegistry_ClosureNamespace(this_ *AutoValue_JSTypeRegistry_ClosureNamespace, isLegacy bool, definitionNode Node, type_ JSType) {
	_ = isLegacy
	_ = definitionNode
	_ = type_
	InitJSTypeRegistry_ClosureNamespace(&this_.JSTypeRegistry_ClosureNamespaceBase)
	this_.initFields_AutoValue_JSTypeRegistry_ClosureNamespace()
	this_.isLegacy = isLegacy
	this_.definitionNode = definitionNode
	this_.type_ = type_
}

// Java method: AutoValue_JSTypeRegistry_ClosureNamespace.isLegacy()
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) IsLegacy() bool {
	return this_.isLegacy
}

// Java method: AutoValue_JSTypeRegistry_ClosureNamespace.definitionNode()
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) DefinitionNode() Node {
	return this_.definitionNode
}

// Java method: AutoValue_JSTypeRegistry_ClosureNamespace.type()
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) Type() JSType {
	return this_.type_
}

// Java method: AutoValue_JSTypeRegistry_ClosureNamespace.toString()
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) String() string {
	return ("ClosureNamespace{isLegacy=" + j2gort.StrValBool(this_.isLegacy) + ", definitionNode=" + j2gort.StrVal(this_.definitionNode) + ", type=" + j2gort.StrVal(this_.type_) + "}")
}

// Java method: AutoValue_JSTypeRegistry_ClosureNamespace.equals(java.lang.@org.jspecify.annotations.Nullable Object o)
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[JSTypeRegistry_ClosureNamespace](o) {
		var that JSTypeRegistry_ClosureNamespace = j2gort.AsIface[JSTypeRegistry_ClosureNamespace](o)
		_ = that
		return (((this_.isLegacy == that.IsLegacy()) && (func() bool {
			if this_.definitionNode == 0 {
				return (that.DefinitionNode() == 0)
			}
			return func() bool {
				rn_2 := this_.definitionNode
				an_3 := that.DefinitionNode()
				return j2gort.CheckNotNullNode(rn_2) == an_3
			}()
		}())) && (func() bool {
			if this_.type_ == nil || this_.type_.IsJavaNil_() {
				return j2gort.IfaceEquals[JSType](that.Type(), nil)
			}
			return this_.type_.Equals(that.Type())
		}()))
	}
	return false
}

// Java method: AutoValue_JSTypeRegistry_ClosureNamespace.hashCode()
func (this_ *AutoValue_JSTypeRegistry_ClosureNamespace) HashCode() int32 {
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.isLegacy {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.definitionNode == 0 {
			return int32(0)
		}
		return int32(j2gort.CheckNotNullNode(this_.definitionNode))
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.type_ == nil || this_.type_.IsJavaNil_() {
			return int32(0)
		}
		return this_.type_.HashCode()
	}())
	return h_d
}
