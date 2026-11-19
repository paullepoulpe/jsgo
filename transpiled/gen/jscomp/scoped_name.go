/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ScopedName.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.ScopedName
type ScopedName interface {
	Implements_jscomp_ScopedName()
	IsJavaNil_() bool
	GetScopedNameBase() *ScopedNameBase
	GetName() string
	GetScopeRoot() pkg_rhino.Node
	Equals(any) bool
	HashCode() int32
}
func ToIface_ScopedName[T interface {
	ScopedName
	comparable
}](v T) ScopedName {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ScopedName
type ScopedNameBase struct {
	name string `gson:"name"`
	cachedHashCode int32 `gson:"cachedHashCode"`
	scopeRoot pkg_rhino.Node `gson:"scopeRoot"`
}

func (this_ *ScopedNameBase) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedNameBase) IsJavaNil_() bool { return this_ == nil }
func (*ScopedNameBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ScopedName", "ScopedName" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedName", "com.google.javascript.jscomp.ScopedName", func(v_ any) bool {
		return j2gort.InstanceOf[ScopedName](v_)
	})
}
func (b *ScopedNameBase) GetScopedNameBase() *ScopedNameBase { return b }
func (*ScopedNameBase) Implements_jscomp_ScopedName() {}

func (this_ *ScopedNameBase) initFields_ScopedName() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: ScopedName(java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node scopeRoot)
func NewScopedName(name string, scopeRoot pkg_rhino.Node) ScopedName {
	this_ := &ScopedNameBase{}
	InitScopedName(this_, name, scopeRoot)
	return this_
}

func InitScopedName(this_ *ScopedNameBase, name string, scopeRoot pkg_rhino.Node) {
	_ = name
	_ = scopeRoot
	this_.initFields_ScopedName()
	this_.name = name
	this_.scopeRoot = scopeRoot
}

// Java method: ScopedName.getName()
func (this_ *ScopedNameBase) GetName() string {
	return this_.name
}

// Java method: ScopedName.getScopeRoot()
func (this_ *ScopedNameBase) GetScopeRoot() pkg_rhino.Node {
	return this_.scopeRoot
}

// Java static method: ScopedName.of(java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node scopeRoot)
func ScopedName_Of(name string, scopeRoot pkg_rhino.Node) ScopedName {
	_ = name
	_ = scopeRoot
	return NewScopedName(name, scopeRoot)
}

// Java method: ScopedName.equals(java.lang.Object other)
func (this_ *ScopedNameBase) Equals(arg_1 any) bool {
	other := arg_1
	_ = other
	var scopedName ScopedName
	_ = scopedName
	var ok_2 bool
	scopedName, ok_2 = j2gort.TryCast[ScopedName](other)
	if !ok_2 {
		return false
	}
	return ((this_.scopeRoot == scopedName.GetScopedNameBase().scopeRoot) && func() bool { sr_3 := this_.name; se_4 := scopedName.GetScopedNameBase().name; return j2gort.CheckNotNullString(sr_3) == se_4 }())
}

// Java method: ScopedName.hashCode()
func (this_ *ScopedNameBase) HashCode() int32 {
	var h int32 = this_.cachedHashCode
	_ = h
	if ((h == int32(0))) {
		return this_.priv_ScopedName_computeHashCode_5030()
	}
	return h
}

// Java method: ScopedName.computeHashCode()
func (this_ *ScopedNameBase) priv_ScopedName_computeHashCode_5030() int32 {
	var result int32 = (int32(31) + javalang.StringHashCode(j2gort.CheckNotNullString(this_.name)))
	_ = result
	var rootHash int32 = int32(0)
	_ = rootHash
	if ((this_.scopeRoot != 0)) {
		rootHash = int32(j2gort.CheckNotNullNode(this_.scopeRoot))
	}
	result = ((int32(31) * result) + rootHash)
	this_.cachedHashCode = result
	return result
}

