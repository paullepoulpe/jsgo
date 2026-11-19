// Port of com/google/javascript/rhino/jstype/AutoValue_FunctionType_Parameter.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.AutoValue_FunctionType_Parameter
type AutoValue_FunctionType_Parameter struct {
	idHash_ atomic.Int32
}

func (this_ *AutoValue_FunctionType_Parameter) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoValue_FunctionType_Parameter) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_FunctionType_Parameter) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.AutoValue_FunctionType_Parameter", "AutoValue_FunctionType_Parameter" }
func (this_ *AutoValue_FunctionType_Parameter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AutoValue_FunctionType_Parameter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.AutoValue_FunctionType_Parameter", "com.google.javascript.rhino.jstype.AutoValue_FunctionType_Parameter", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_FunctionType_Parameter](v_)
	})
}

func (this_ *AutoValue_FunctionType_Parameter) initFields_AutoValue_FunctionType_Parameter() {
	_ = this_
}

// Java constructor: AutoValue_FunctionType_Parameter()
func NewAutoValue_FunctionType_Parameter() *AutoValue_FunctionType_Parameter {
	this_ := &AutoValue_FunctionType_Parameter{}
	InitAutoValue_FunctionType_Parameter(this_)
	return this_
}

func InitAutoValue_FunctionType_Parameter(this_ *AutoValue_FunctionType_Parameter) {
	this_.initFields_AutoValue_FunctionType_Parameter()
}

