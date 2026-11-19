/*
 * Copyright 2010 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/AbstractPeepholeTranspilation.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractPeepholeTranspilation
type AbstractPeepholeTranspilation interface {
	Implements_jscomp_AbstractPeepholeTranspilation()
	IsJavaNil_() bool
	GetAbstractPeepholeTranspilationBase() *AbstractPeepholeTranspilationBase
	GetTranspiledAwayFeatures() *pkg_parser.FeatureSet
	GetAdditionalFeaturesToRunOn() *pkg_parser.FeatureSet
	TranspileSubtree(pkg_rhino.Node) pkg_rhino.Node
}
func ToIface_AbstractPeepholeTranspilation[T interface {
	AbstractPeepholeTranspilation
	comparable
}](v T) AbstractPeepholeTranspilation {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AbstractPeepholeTranspilation[T AbstractPeepholeTranspilation](v T) AbstractPeepholeTranspilation {
	var vi AbstractPeepholeTranspilation = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.AbstractPeepholeTranspilation
type AbstractPeepholeTranspilationBase struct {
	idHash_ atomic.Int32
}

func (this_ *AbstractPeepholeTranspilationBase) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractPeepholeTranspilationBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractPeepholeTranspilationBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.AbstractPeepholeTranspilation", "AbstractPeepholeTranspilation" }
func (this_ *AbstractPeepholeTranspilationBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractPeepholeTranspilation", "com.google.javascript.jscomp.AbstractPeepholeTranspilation", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractPeepholeTranspilation](v_)
	})
}
func (b *AbstractPeepholeTranspilationBase) GetAbstractPeepholeTranspilationBase() *AbstractPeepholeTranspilationBase { return b }
func (*AbstractPeepholeTranspilationBase) Implements_jscomp_AbstractPeepholeTranspilation() {}

func (this_ *AbstractPeepholeTranspilationBase) initFields_AbstractPeepholeTranspilation() {
	_ = this_
}

func InitAbstractPeepholeTranspilation(this_ *AbstractPeepholeTranspilationBase) {
	this_.initFields_AbstractPeepholeTranspilation()
}

// Java method: AbstractPeepholeTranspilation.getTranspiledAwayFeatures()
func (this_ *AbstractPeepholeTranspilationBase) GetTranspiledAwayFeatures() *pkg_parser.FeatureSet {
	panic("abstract or native method invoked")
}

// Java method: AbstractPeepholeTranspilation.getAdditionalFeaturesToRunOn()
func (this_ *AbstractPeepholeTranspilationBase) GetAdditionalFeaturesToRunOn() *pkg_parser.FeatureSet {
	return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
}

// Java method: AbstractPeepholeTranspilation.transpileSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *AbstractPeepholeTranspilationBase) TranspileSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	panic("abstract or native method invoked")
}

