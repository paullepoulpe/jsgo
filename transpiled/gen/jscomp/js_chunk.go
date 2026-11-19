/*
 * Copyright 2005 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/JSChunk.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.JSChunk
type JSChunk struct {
	name    string                                    `gson:"name"`
	inputs  j2gort.OrderedMap[string, *CompilerInput] `gson:"-"`
	deps    j2gort.Ext_java_util_List                 `gson:"deps"`
	depth   int32                                     `gson:"depth"`
	index   int32                                     `gson:"index"`
	idHash_ atomic.Int32
}

func (this_ *JSChunk) BindSelf_(s any) {
	_ = s
}
func (this_ *JSChunk) IsJavaNil_() bool { return this_ == nil }
func (*JSChunk) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JSChunk", "JSChunk"
}
func (this_ *JSChunk) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSChunk) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JSChunk", "com.google.javascript.jscomp.JSChunk", func(v_ any) bool {
		return j2gort.InstanceOf[*JSChunk](v_)
	})
}
func (*JSChunk) Implements_jscomp_DependencyInfo() {}

func StaticGet_JSChunk_STRONG_CHUNK_NAME() string { return "$strong$" }
func StaticGet_JSChunk_WEAK_CHUNK_NAME() string   { return "$weak$" }
func StaticGet_JSChunk_SerialVersionUID() int64   { return int64(1) }
func (this_ *JSChunk) initFields_JSChunk() {
	_ = this_
	this_.name = j2gort.NullString
	this_.inputs = j2gort.NewOrderedMap[string, *CompilerInput]()
	this_.deps = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: JSChunk(java.lang.String name)
func NewJSChunk(name string) *JSChunk {
	this_ := &JSChunk{}
	InitJSChunk(this_, name)
	return this_
}

func InitJSChunk(this_ *JSChunk, name string) {
	_ = name
	this_.initFields_JSChunk()
	this_.name = name
	this_.depth = int32(-1)
	this_.index = int32(-1)
}

// Java method: JSChunk.getName()
func (this_ *JSChunk) GetName() string {
	return this_.name
}

// Java method: JSChunk.getProvides()
func (this_ *JSChunk) GetProvides() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListOfStr1(this_.name)
}

// Java method: JSChunk.getHasExternsAnnotation()
func (this_ *JSChunk) GetHasExternsAnnotation() bool {
	return false
}

// Java method: JSChunk.getHasNoCompileAnnotation()
func (this_ *JSChunk) GetHasNoCompileAnnotation() bool {
	return false
}

// Java method: JSChunk.getRequires()
func (this_ *JSChunk) GetRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	var builder []*DependencyInfo_Require = []*DependencyInfo_Require{}
	_ = builder
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.deps))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_2)
			_ = chunk
			builder = append(builder, j2gort.CheckNotNullElem(DependencyInfo_Require_CompilerChunk(chunk.GetName())))
		}
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java method: JSChunk.getTypeRequires()
func (this_ *JSChunk) GetTypeRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

// Java method: JSChunk.getPathRelativeToClosureBase()
func (this_ *JSChunk) GetPathRelativeToClosureBase() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: JSChunk.getLoadFlags()
func (this_ *JSChunk) GetLoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: JSChunk.isGoogModule()
func (this_ *JSChunk) IsGoogModule() bool {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: JSChunk.isEs6Module()
func (this_ *JSChunk) IsEs6Module() bool {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: JSChunk.add(com.google.javascript.jscomp.SourceFile file)
func (this_ *JSChunk) Add_1_SourceFile(file *SourceFile) {
	_ = file
	this_.Add_1_CompilerInput(NewCompilerInput_1(file))
}

// Java method: JSChunk.add(com.google.javascript.jscomp.CompilerInput input)
func (this_ *JSChunk) Add_1_CompilerInput(input *CompilerInput) {
	_ = input
	var inputName string = input.GetName()
	_ = inputName
	var previous *CompilerInput = this_.inputs.Put(inputName, input)
	_ = previous
	if !(previous == nil) {
		j2gort.CheckStateFmt(false, "%s already exist in chunk %s", j2gort.BoxString(inputName), j2gort.BoxString(this_.GetName()))
	}
	input.SetChunk(this_)
}

// Java method: JSChunk.addDependency(com.google.javascript.jscomp.JSChunk dep)
func (this_ *JSChunk) AddDependency(dep *JSChunk) {
	_ = dep
	j2gort.CheckNotNullPtr(dep)
	if !(dep != this_) {
		j2gort.CheckStateFmt(false, "Cannot add dependency on self (%s)", this_)
	}
	this_.deps.Add_1(j2gort.PtrToAny(dep))
}

// Java method: JSChunk.remove(com.google.javascript.jscomp.CompilerInput input)
func (this_ *JSChunk) Remove(input *CompilerInput) {
	_ = input
	input.SetChunk(nil)
	this_.inputs.Remove(input.GetName())
}

// Java method: JSChunk.removeAll()
func (this_ *JSChunk) RemoveAll() {
	{
		cur_3 := this_.inputs.IterValues()
		for cur_3.HasNext() {

			input := cur_3.
				Next()
			_ = input
			input.SetChunk(nil)
		}
	}
	this_.inputs.Clear()
}

// Java method: JSChunk.getDependencies()
func (this_ *JSChunk) GetDependencies() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.deps))
}

// Java method: JSChunk.getSortedDependencyNames()
func (this_ *JSChunk) GetSortedDependencyNames() j2gort.Ext_java_util_List {
	var names []string = []string{}
	_ = names
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetDependencies()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_6)
			_ = chunk
			names = append(names, chunk.GetName())
		}
	}
	j2gort.SliceSortNatural(names)
	return j2gort.WrapTypedArrayList(names)
}

// Java method: JSChunk.getAllDependencies()
func (this_ *JSChunk) GetAllDependencies() j2gort.Ext_java_util_Set {
	var allDeps *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = allDeps
	allDeps.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.deps))
	var stack j2gort.Deque[*JSChunk] = j2gort.DequeFromCol[*JSChunk](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.deps))
	_ = stack
	for !(stack.IsEmpty()) {
		var chunk *JSChunk = stack.Pop()
		_ = chunk
		var chunkDeps j2gort.Ext_java_util_List = chunk.deps
		_ = chunkDeps
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkDeps))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				dep := j2gort.CastPtr[JSChunk](rawElem_8)
				_ = dep
				if allDeps.Add_1(j2gort.PtrToAny(dep)) {
					stack.Push(dep)
				}
			}
		}
	}
	return allDeps
}

// Java method: JSChunk.getThisAndAllDependencies()
func (this_ *JSChunk) GetThisAndAllDependencies() j2gort.Ext_java_util_Set {
	var deps j2gort.Ext_java_util_Set = this_.GetAllDependencies()
	_ = deps
	deps.Add_1(this_)
	return deps
}

// Java method: JSChunk.getInputCount()
func (this_ *JSChunk) GetInputCount() int32 {
	return this_.inputs.Size()
}

// Java method: JSChunk.getFirst()
func (this_ *JSChunk) GetFirst() *CompilerInput {
	return this_.inputs.FirstValue()
}

// Java method: JSChunk.getInputs()
func (this_ *JSChunk) GetInputs() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.inputs.ValuesToImmutableList()
}

// Java method: JSChunk.getInputsIterable()
func (this_ *JSChunk) GetInputsIterable() j2gort.Ext_java_lang_Iterable {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.inputs.LiveValues())
}

// Java method: JSChunk.getByName(java.lang.String name)
func (this_ *JSChunk) GetByName(name string) *CompilerInput {
	_ = name
	return this_.inputs.Get(name)
}

// Java method: JSChunk.removeByName(java.lang.String name)
func (this_ *JSChunk) RemoveByName(name string) bool {
	_ = name
	var value *CompilerInput = this_.inputs.Remove(name)
	_ = value
	return (value != nil)
}

// Java method: JSChunk.isSynthetic()
func (this_ *JSChunk) IsSynthetic() bool {
	return (func() bool { sr_9 := this_.name; se_10 := "$strong$"; return j2gort.CheckNotNullString(sr_9) == se_10 }() || func() bool { sr_11 := this_.name; se_12 := "$weak$"; return j2gort.CheckNotNullString(sr_11) == se_12 }())
}

// Java method: JSChunk.isWeak()
func (this_ *JSChunk) IsWeak() bool {
	return func() bool { sr_13 := this_.name; se_14 := "$weak$"; return j2gort.CheckNotNullString(sr_13) == se_14 }()
}

// Java method: JSChunk.toString()
func (this_ *JSChunk) String() string {
	return this_.name
}

// Java method: JSChunk.sortInputsByDeps(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *JSChunk) SortInputsByDeps(compiler AbstractCompiler) {
	_ = compiler
	{
		cur_15 := this_.inputs.IterValues()
		for cur_15.HasNext() {

			input := cur_15.
				Next()
			_ = input
			input.SetCompiler(compiler)
		}
	}
	var sortedList j2gort.Ext_java_util_List = NewSortedDependencies(j2gort.AsIface[j2gort.Ext_java_util_List](this_.GetInputs())).GetSortedList()
	_ = sortedList
	this_.inputs.Clear()
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(sortedList))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			input_2 := j2gort.CastPtr[CompilerInput](rawElem_18)
			_ = input_2
			this_.inputs.Put(input_2.GetName(), input_2)
		}
	}
}

// Java method: JSChunk.setDepth(int dep)
func (this_ *JSChunk) SetDepth(dep int32) {
	_ = dep
	if !(dep >= int32(0)) {
		j2gort.CheckArgumentFmt(false, "invalid depth: %s", j2gort.ValueOfInt(dep))
	}
	this_.depth = dep
}

// Java method: JSChunk.getDepth()
func (this_ *JSChunk) GetDepth() int32 {
	return this_.depth
}

// Java method: JSChunk.setIndex(int index)
func (this_ *JSChunk) SetIndex(index int32) {
	_ = index
	if !(index >= int32(0)) {
		j2gort.CheckArgumentFmt(false, "Invalid chunk index: %s", j2gort.ValueOfInt(index))
	}
	this_.index = index
}

// Java method: JSChunk.getIndex()
func (this_ *JSChunk) GetIndex() int32 {
	return this_.index
}

// Java method: JSChunk.readObject(java.io.ObjectInputStream in)
func (this_ *JSChunk) priv_JSChunk_readObject_3521(in *j2gort.Ext_java_io_ObjectInputStream) {
	_ = in
	in.DefaultReadObject()
	this_.inputs = j2gort.NewOrderedMap[string, *CompilerInput]()
}

func (this_ *JSChunk) GetRequiredSymbols() j2gort.Ext_com_google_common_collect_ImmutableList {
	return IfaceDefault_DependencyInfo_GetRequiredSymbols(this_)
}
