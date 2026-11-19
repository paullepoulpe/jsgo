// Port of com/google/javascript/jscomp/deps/AutoValue_DependencyInfo_Require.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.deps.AutoValue_DependencyInfo_Require
type AutoValue_DependencyInfo_Require struct {
	idHash_ atomic.Int32
}

func (this_ *AutoValue_DependencyInfo_Require) BindSelf_(s any) {
	_ = s
}
func (this_ *AutoValue_DependencyInfo_Require) IsJavaNil_() bool { return this_ == nil }
func (*AutoValue_DependencyInfo_Require) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.deps.AutoValue_DependencyInfo_Require", "AutoValue_DependencyInfo_Require" }
func (this_ *AutoValue_DependencyInfo_Require) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AutoValue_DependencyInfo_Require) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.AutoValue_DependencyInfo_Require", "com.google.javascript.jscomp.deps.AutoValue_DependencyInfo_Require", func(v_ any) bool {
		return j2gort.InstanceOf[*AutoValue_DependencyInfo_Require](v_)
	})
}

func (this_ *AutoValue_DependencyInfo_Require) initFields_AutoValue_DependencyInfo_Require() {
	_ = this_
}

// Java constructor: AutoValue_DependencyInfo_Require()
func NewAutoValue_DependencyInfo_Require() *AutoValue_DependencyInfo_Require {
	this_ := &AutoValue_DependencyInfo_Require{}
	InitAutoValue_DependencyInfo_Require(this_)
	return this_
}

func InitAutoValue_DependencyInfo_Require(this_ *AutoValue_DependencyInfo_Require) {
	this_.initFields_AutoValue_DependencyInfo_Require()
}

