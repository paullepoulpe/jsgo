/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReferenceCollection.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ReferenceCollection
type ReferenceCollection struct {
	ref0    *Reference   `gson:"ref0"`
	ref1    *Reference   `gson:"ref1"`
	ref2    *Reference   `gson:"ref2"`
	ref3    *Reference   `gson:"ref3"`
	rest    []*Reference `gson:"rest"`
	idHash_ atomic.Int32
}

func (this_ *ReferenceCollection) BindSelf_(s any) {
	_ = s
}
func (this_ *ReferenceCollection) IsJavaNil_() bool { return this_ == nil }
func (*ReferenceCollection) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReferenceCollection", "ReferenceCollection"
}
func (this_ *ReferenceCollection) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReferenceCollection) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReferenceCollection", "com.google.javascript.jscomp.ReferenceCollection", func(v_ any) bool {
		return j2gort.InstanceOf[*ReferenceCollection](v_)
	})
}
func (*ReferenceCollection) Implements_j2gort_Ext_java_lang_Iterable() {}

func (this_ *ReferenceCollection) initFields_ReferenceCollection() {
	_ = this_
}

// Java constructor: ReferenceCollection()
func NewReferenceCollection() *ReferenceCollection {
	this_ := &ReferenceCollection{}
	InitReferenceCollection(this_)
	return this_
}
func allocPooled_ReferenceCollection(slab_ *j2gort.PooledSlab[ReferenceCollection]) *ReferenceCollection {
	this_ := slab_.Alloc()
	InitReferenceCollection(this_)
	return this_
}

func InitReferenceCollection(this_ *ReferenceCollection) {
	this_.initFields_ReferenceCollection()
}

// Java method: ReferenceCollection.iterator()
func (this_ *ReferenceCollection) Iterator() any {
	if this_.rest != nil {
		return j2gort.NewTypedSliceIter(this_.rest)
	}
	var n int32 = this_.Size()
	_ = n
	if n == int32(0) {
		return j2gort.Ext_java_util_Collections_EmptyIterator()
	}
	return NewReferenceCollection_InlineIterator(this_.ref0, this_.ref1, this_.ref2, this_.ref3, n)
}

// Java method: ReferenceCollection.size()
func (this_ *ReferenceCollection) Size() int32 {
	if this_.rest != nil {
		return int32(len(this_.rest))
	}
	if this_.ref0 == nil {
		return int32(0)
	}
	if this_.ref1 == nil {
		return int32(1)
	}
	if this_.ref2 == nil {
		return int32(2)
	}
	if this_.ref3 == nil {
		return int32(3)
	} else {
		return int32(4)
	}
}

// Java method: ReferenceCollection.isEmpty()
func (this_ *ReferenceCollection) IsEmpty() bool {
	return ((this_.ref0 == nil) && (this_.rest == nil))
}

// Java method: ReferenceCollection.get(int index)
func (this_ *ReferenceCollection) Get(index int32) *Reference {
	_ = index
	if this_.rest != nil {
		return this_.rest[index]
	}
	var r *Reference = nil
	_ = r
	switch index {
	case int32(0):
		r = this_.ref0
	case int32(1):
		r = this_.ref1
	case int32(2):
		r = this_.ref2
	case int32(3):
		r = this_.ref3
	default:
		r = nil
	}
	if r != nil {
		return r
	}
	panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
}

// Java method: ReferenceCollection.add(com.google.javascript.jscomp.Reference reference)
func (this_ *ReferenceCollection) Add(reference *Reference) {
	_ = reference
	if this_.rest != nil {
		j2gort.AppendPooledSlice(&this_.rest, reference)
	} else if this_.ref0 == nil {
		this_.ref0 = reference
	} else if this_.ref1 == nil {
		this_.ref1 = reference
	} else if this_.ref2 == nil {
		this_.ref2 = reference
	} else if this_.ref3 == nil {
		this_.ref3 = reference
	} else {
		var list []*Reference = j2gort.AcquirePooledSlice[*Reference](8)
		_ = list
		list = append(list, this_.ref0)
		list = append(list, this_.ref1)
		list = append(list, this_.ref2)
		list = append(list, this_.ref3)
		list = append(list, reference)
		this_.rest = list
		this_.ref0 = nil
		this_.ref1 = nil
		this_.ref2 = nil
		this_.ref3 = nil
	}
}

// Java method: ReferenceCollection.isWellDefined()
func (this_ *ReferenceCollection) IsWellDefined() bool {
	var size int32 = this_.Size()
	_ = size
	if size == int32(0) {
		return false
	}
	var init *Reference = this_.GetInitializingReference()
	_ = init
	if init == nil {
		return false
	}
	if !(this_.Get(int32(0)).IsDeclaration()) {
		j2gort.CheckState(false)
	}
	var initBlock *BasicBlock = init.GetBasicBlock()
	_ = initBlock
	{
		var i int32 = int32(1)
		_ = i
		for ; i < size; i = i + 1 {
			if !(initBlock.ProvablyExecutesBefore(this_.Get(i).GetBasicBlock())) {
				return false
			}
		}
	}
	return true
}

// Java method: ReferenceCollection.isEscaped()
func (this_ *ReferenceCollection) IsEscaped() bool {
	var closureScope *Scope = nil
	_ = closureScope
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = this_.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = this_.Get(i)
			_ = ref
			var refClosureScope *Scope = ReferenceCollection_GetEnclosingClosureScope(ref.GetScope())
			_ = refClosureScope
			if closureScope == nil {
				closureScope = refClosureScope
			} else if closureScope != refClosureScope {
				return true
			}
		}
	}
	return false
}

// Java static method: ReferenceCollection.getEnclosingClosureScope(com.google.javascript.jscomp.Scope scope)
func ReferenceCollection_GetEnclosingClosureScope(scope *Scope) *Scope {
	_ = scope
	{
		var s *Scope = scope
		_ = s
		for ; s != nil; s = j2gort.CastPtr[Scope](s.GetParent()) {
			if s.IsHoistScope() || ((s.IsMemberFieldDefScope() || s.IsComputedFieldDefRhsScope()) && !(s.GetRootNode_0().IsStaticMember())) {
				return s
			}
		}
	}
	return nil
}

// Java method: ReferenceCollection.isThisRebound()
func (this_ *ReferenceCollection) IsThisRebound() bool {
	var hoistScope *Scope = nil
	_ = hoistScope
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = this_.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = this_.Get(i)
			_ = ref
			var refHoistScope *Scope = j2gort.CastPtr[Scope](ref.GetScope().GetClosestHoistScope())
			_ = refHoistScope
			if hoistScope == nil {
				hoistScope = refHoistScope
			} else if hoistScope != refHoistScope {
				return true
			}
			{
				var s *Scope = ref.GetScope()
				_ = s
				for ; (s != hoistScope) && (s != nil); s = j2gort.CastPtr[Scope](s.GetParent()) {
					if StaticGet_AbstractScope_ImplicitVar_THIS().IsMadeByScope(ToIface_AbstractScope(s)) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Java method: ReferenceCollection.isInitializingDeclarationAt(int index)
func (this_ *ReferenceCollection) priv_ReferenceCollection_isInitializingDeclarationAt_4632(index int32) bool {
	_ = index
	var maybeInit *Reference = this_.Get(index)
	_ = maybeInit
	if maybeInit.IsInitializingDeclaration() {
		return true
	}
	return false
}

// Java method: ReferenceCollection.isInitializingAssignmentAt(int index)
func (this_ *ReferenceCollection) priv_ReferenceCollection_isInitializingAssignmentAt_4633(index int32) bool {
	_ = index
	if (index < this_.Size()) && (index > int32(0)) {
		var maybeDecl *Reference = this_.Get((index - int32(1)))
		_ = maybeDecl
		if maybeDecl.IsVarDeclaration() || maybeDecl.IsLetDeclaration() {
			if !(!(maybeDecl.IsInitializingDeclaration())) {
				j2gort.CheckState(false)
			}
			var maybeInit *Reference = this_.Get(index)
			_ = maybeInit
			if maybeInit.IsSimpleAssignmentToName() {
				return true
			}
		}
	}
	return false
}

// Java method: ReferenceCollection.getInitializingReference()
func (this_ *ReferenceCollection) GetInitializingReference() *Reference {
	if this_.priv_ReferenceCollection_isInitializingDeclarationAt_4632(int32(0)) {
		return this_.Get(int32(0))
	} else if this_.priv_ReferenceCollection_isInitializingAssignmentAt_4633(int32(1)) && this_.priv_ReferenceCollection_areInSameFile_4635(int32(0), int32(1)) {
		return this_.Get(int32(1))
	}
	return nil
}

// Java method: ReferenceCollection.areInSameFile(int idx0, int idx1)
func (this_ *ReferenceCollection) priv_ReferenceCollection_areInSameFile_4635(idx0 int32, idx1 int32) bool {
	_ = idx0
	_ = idx1
	return j2gort.ObjEqCustomPtr(this_.Get(idx0).GetInputId(), this_.Get(idx1).GetInputId())
}

// Java method: ReferenceCollection.getInitializingReferenceForConstants()
func (this_ *ReferenceCollection) GetInitializingReferenceForConstants() *Reference {
	var size int32 = this_.Size()
	_ = size
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			if this_.priv_ReferenceCollection_isInitializingDeclarationAt_4632(i) || this_.priv_ReferenceCollection_isInitializingAssignmentAt_4633(i) {
				return this_.Get(i)
			}
		}
	}
	return nil
}

// Java method: ReferenceCollection.isAssignedOnceInLifetime()
func (this_ *ReferenceCollection) IsAssignedOnceInLifetime() bool {
	var ref *Reference = this_.GetOneAndOnlyAssignment()
	_ = ref
	if ref == nil {
		return false
	}
	{
		var block *BasicBlock = ref.GetBasicBlock()
		_ = block
		for ; block != nil; block = block.GetParent() {
			if block.IsFunction() {
				if j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](j2gort.CastPtr[Var](ref.GetSymbol()).GetScope()).GetClosestHoistScope()) != j2gort.CastPtr[Scope](ref.GetScope().GetClosestHoistScope()) {
					return false
				}
				break
			}
			if block.GetRoot() == j2gort.CastPtr[Scope](j2gort.CastPtr[Var](ref.GetSymbol()).GetScope()).GetRootNode_0() {
				break
			}
			if block.IsLoop() {
				return false
			}
		}
	}
	return true
}

// Java method: ReferenceCollection.isReferencedWeakly()
func (this_ *ReferenceCollection) IsReferencedWeakly() bool {
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = this_.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = this_.Get(i)
			_ = ref
			if NodeUtil_IsGoogWeakUsageCall(ref.GetParent()) {
				return true
			}
		}
	}
	return false
}

// Java method: ReferenceCollection.getOneAndOnlyAssignment()
func (this_ *ReferenceCollection) GetOneAndOnlyAssignment() *Reference {
	var assignment *Reference = nil
	_ = assignment
	var size int32 = this_.Size()
	_ = size
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			var ref *Reference = this_.Get(i)
			_ = ref
			if ref.IsLvalue() || ref.IsInitializingDeclaration() {
				if assignment == nil {
					assignment = ref
				} else {
					return nil
				}
			}
		}
	}
	return assignment
}

// Java method: ReferenceCollection.isNeverAssigned()
func (this_ *ReferenceCollection) IsNeverAssigned() bool {
	var size int32 = this_.Size()
	_ = size
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			var ref *Reference = this_.Get(i)
			_ = ref
			if ref.IsLvalue() || ref.IsInitializingDeclaration() {
				return false
			}
		}
	}
	return true
}

// Java method: ReferenceCollection.firstReferenceIsAssigningDeclaration()
func (this_ *ReferenceCollection) FirstReferenceIsAssigningDeclaration() bool {
	return (!(this_.IsEmpty()) && this_.Get(int32(0)).IsInitializingDeclaration())
}

// Java method: ReferenceCollection.toString()
func (this_ *ReferenceCollection) String() string {
	var refsView any = nil
	_ = refsView
	if this_.rest != nil {
		refsView = j2gort.SliceToJavaList(this_.rest)
	} else {
		var n int32 = this_.Size()
		_ = n
		if n == int32(0) {
			refsView = j2gort.Ext_java_util_Collections_EmptyList()
		} else if n == int32(1) {
			refsView = j2gort.SingletonListPtr(this_.ref0)
		} else {
			var list []*Reference = j2gort.MakeNonNilSlice[*Reference](n)
			_ = list
			{
				var i int32 = int32(0)
				_ = i
				for ; i < n; i = i + 1 {
					list = append(list, this_.Get(i))
				}
			}
			refsView = j2gort.SliceToJavaList(list)
		}
	}
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("initRef", j2gort.PtrToAny(this_.GetInitializingReference())).Add_2_string_any("references", refsView).Add_2_string_bool("wellDefined", this_.IsWellDefined()).Add_2_string_bool("assignedOnce", this_.IsAssignedOnceInLifetime()).String()
}

func (this_ *ReferenceCollection) ForEach(p0_ j2gort.Ext_java_util_function_Consumer) {
	j2gort.ExtDefault_Iterable_ForEach(this_, p0_)
}

// Java class: com.google.javascript.jscomp.ReferenceCollection.InlineIterator
type ReferenceCollection_InlineIterator struct {
	r0      *Reference `gson:"r0"`
	r1      *Reference `gson:"r1"`
	r2      *Reference `gson:"r2"`
	r3      *Reference `gson:"r3"`
	size    int32      `gson:"size"`
	index   int32      `gson:"index"`
	idHash_ atomic.Int32
}

func (this_ *ReferenceCollection_InlineIterator) BindSelf_(s any) {
	_ = s
}
func (this_ *ReferenceCollection_InlineIterator) IsJavaNil_() bool { return this_ == nil }
func (*ReferenceCollection_InlineIterator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReferenceCollection$InlineIterator", "InlineIterator"
}
func (this_ *ReferenceCollection_InlineIterator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReferenceCollection_InlineIterator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReferenceCollection.InlineIterator", "com.google.javascript.jscomp.ReferenceCollection$InlineIterator", func(v_ any) bool {
		return j2gort.InstanceOf[*ReferenceCollection_InlineIterator](v_)
	})
}
func (*ReferenceCollection_InlineIterator) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *ReferenceCollection_InlineIterator) initFields_ReferenceCollection_InlineIterator() {
	_ = this_
}

// Java constructor: InlineIterator(com.google.javascript.jscomp.Reference r0, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference r1, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference r2, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference r3, int size)
func NewReferenceCollection_InlineIterator(r0 *Reference, r1 *Reference, r2 *Reference, r3 *Reference, size int32) *ReferenceCollection_InlineIterator {
	this_ := &ReferenceCollection_InlineIterator{}
	InitReferenceCollection_InlineIterator(this_, r0, r1, r2, r3, size)
	return this_
}

func InitReferenceCollection_InlineIterator(this_ *ReferenceCollection_InlineIterator, r0 *Reference, r1 *Reference, r2 *Reference, r3 *Reference, size int32) {
	_ = r0
	_ = r1
	_ = r2
	_ = r3
	_ = size
	this_.initFields_ReferenceCollection_InlineIterator()
	this_.r0 = r0
	this_.r1 = r1
	this_.r2 = r2
	this_.r3 = r3
	this_.size = size
}

// Java method: InlineIterator.hasNext()
func (this_ *ReferenceCollection_InlineIterator) HasNext() bool {
	return (this_.index < this_.size)
}

// Java method: InlineIterator.next()
func (this_ *ReferenceCollection_InlineIterator) Next() any {
	if this_.index >= this_.size {
		panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
	}
	switch func() int32 { old_1 := this_.index; nw_2 := old_1 + 1; this_.index = nw_2; return old_1 }() {
	case int32(0):
		return j2gort.PtrToAny(this_.r0)
	case int32(1):
		return j2gort.PtrToAny(this_.r1)
	case int32(2):
		return j2gort.PtrToAny(this_.r2)
	default:
		return j2gort.PtrToAny(this_.r3)
	}
	panic("unreachable switch expression")
}

func (this_ *ReferenceCollection_InlineIterator) Remove() {
	panic(javalang.NewUnsupportedOperationException("remove"))
}
