/*
 * Copyright 2022 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PassListBuilder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PassListBuilder
type PassListBuilder struct {
	options *CompilerOptions          `gson:"options"`
	passes  j2gort.Ext_java_util_List `gson:"passes"`
	idHash_ atomic.Int32
}

func (this_ *PassListBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *PassListBuilder) IsJavaNil_() bool { return this_ == nil }
func (*PassListBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PassListBuilder", "PassListBuilder"
}
func (this_ *PassListBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PassListBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassListBuilder", "com.google.javascript.jscomp.PassListBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*PassListBuilder](v_)
	})
}

func (this_ *PassListBuilder) initFields_PassListBuilder() {
	_ = this_
	this_.passes = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: PassListBuilder(com.google.javascript.jscomp.CompilerOptions options)
func NewPassListBuilder(options *CompilerOptions) *PassListBuilder {
	this_ := &PassListBuilder{}
	InitPassListBuilder(this_, options)
	return this_
}

func InitPassListBuilder(this_ *PassListBuilder, options *CompilerOptions) {
	_ = options
	this_.initFields_PassListBuilder()
	this_.options = options
}

// Java method: PassListBuilder.build()
func (this_ *PassListBuilder) Build() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.passes))
}

// Java method: PassListBuilder.addAll(com.google.javascript.jscomp.PassListBuilder other)
func (this_ *PassListBuilder) AddAll(other *PassListBuilder) {
	_ = other
	this_.passes.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](other.Build()))
}

// Java method: PassListBuilder.maybeAdd(com.google.javascript.jscomp.PassFactory factory)
func (this_ *PassListBuilder) MaybeAdd(factory PassFactory) {
	_ = factory
	if j2gort.UnboxBoolean((j2gort.CastPtr[j2gort.BoxBoolean](factory.GetCondition().Apply(j2gort.PtrToAny(this_.options))))) {
		this_.passes.Add_1(factory)
	}
}

// Java method: PassListBuilder.addBefore(com.google.javascript.jscomp.PassFactory factory, java.lang.String passName)
func (this_ *PassListBuilder) AddBefore(factory PassFactory, passName string) {
	_ = factory
	_ = passName
	this_.passes.Add_2(this_.priv_PassListBuilder_findIndexByName_4145(passName), factory)
}

// Java method: PassListBuilder.addAfter(com.google.javascript.jscomp.PassFactory factory, java.lang.String passName)
func (this_ *PassListBuilder) AddAfter(factory PassFactory, passName string) {
	_ = factory
	_ = passName
	this_.passes.Add_2((this_.priv_PassListBuilder_findIndexByName_4145(passName) + int32(1)), factory)
}

// Java method: PassListBuilder.addAllUpTo(com.google.javascript.jscomp.PassListBuilder other, java.lang.String passName)
func (this_ *PassListBuilder) AddAllUpTo(other *PassListBuilder, passName string) {
	_ = other
	_ = passName
	var otherPasses []PassFactory = j2gort.ImmListToSliceNullable[PassFactory](other.Build())
	_ = otherPasses
	var index int32 = int32(-1)
	_ = index
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(otherPasses)); i = i + 1 {
			{
				var iifeCond_1 bool
				{

					sr_1 := passName
					se_2 := otherPasses[i].GetName()
					iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
				}
				if iifeCond_1 {
					index = i
					break
				}
			}
		}
	}
	this_.passes.AddAll_1(j2gort.ImmListFromTypedSlice(j2gort.ImmSliceSubList(otherPasses, int32(0), (index + int32(1)))))
}

// Java method: PassListBuilder.findByName(java.lang.String name)
func (this_ *PassListBuilder) FindByName(name string) PassFactory {
	_ = name
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.passes))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			pass := j2gort.AsIface[PassFactory](rawElem_4)
			_ = pass
			if j2gort.CheckNotNullString(pass.GetName()) == name {
				return pass
			}
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "No factory named '"+j2gort.StrValString(name)+"' in the factory list")))
}

// Java method: PassListBuilder.removeByName(java.lang.String name)
func (this_ *PassListBuilder) RemoveByName(name string) {
	_ = name
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.passes.Size(); i = i + 1 {
			if j2gort.CheckNotNullString(j2gort.AsIface[PassFactory](this_.passes.Get(i)).GetName()) == name {
				this_.passes.Remove_1_int32(i)
				return
			}
		}
	}
}

// Java method: PassListBuilder.findIndexByName(java.lang.String name)
func (this_ *PassListBuilder) priv_PassListBuilder_findIndexByName_4145(name string) int32 {
	_ = name
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.passes.Size(); i = i + 1 {
			if j2gort.CheckNotNullString(j2gort.AsIface[PassFactory](this_.passes.Get(i)).GetName()) == name {
				return i
			}
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "No factory named '"+j2gort.StrValString(name)+"' in the factory list")))
}

// Java method: PassListBuilder.contains(com.google.javascript.jscomp.PassFactory factory)
func (this_ *PassListBuilder) Contains(factory PassFactory) bool {
	_ = factory
	return this_.passes.Contains(factory)
}

// Java method: PassListBuilder.assertAllOneTimePasses()
func (this_ *PassListBuilder) AssertAllOneTimePasses() {
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.passes))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			pass := j2gort.AsIface[PassFactory](rawElem_6)
			_ = pass
			if !(!(pass.IsRunInFixedPointLoop())) {
				j2gort.CheckState(false)
			}
		}
	}
}

// Java method: PassListBuilder.assertAllLoopablePasses()
func (this_ *PassListBuilder) AssertAllLoopablePasses() {
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.passes))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			pass := j2gort.AsIface[PassFactory](rawElem_8)
			_ = pass
			if !(pass.IsRunInFixedPointLoop()) {
				j2gort.CheckState(false)
			}
		}
	}
}

// Java method: PassListBuilder.assertPassOrder(com.google.javascript.jscomp.PassFactory pass1, com.google.javascript.jscomp.PassFactory pass2, java.lang.String msg)
func (this_ *PassListBuilder) AssertPassOrder(pass1 PassFactory, pass2 PassFactory, msg string) {
	_ = pass1
	_ = pass2
	_ = msg
	var pass1Index int32 = this_.passes.IndexOf_1(pass1)
	_ = pass1Index
	var pass2Index int32 = this_.passes.IndexOf_1(pass2)
	_ = pass2Index
	if (pass1Index != int32(-1)) && (pass2Index != int32(-1)) {
		if !(pass1Index < pass2Index) {
			j2gort.CheckStateStr(false, msg)
		}
	}
}
