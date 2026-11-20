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
// Port of com/google/javascript/jscomp/OptimizeCalls.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.OptimizeCalls
type OptimizeCalls struct {
	compiler        AbstractCompiler                      `gson:"compiler"`
	passes          []OptimizeCalls_CallGraphCompilerPass `gson:"passes"`
	idHash_         atomic.Int32
	considerExterns bool `gson:"considerExterns"`
}

func (this_ *OptimizeCalls) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeCalls) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeCalls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeCalls", "OptimizeCalls"
}
func (this_ *OptimizeCalls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeCalls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeCalls", "com.google.javascript.jscomp.OptimizeCalls", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeCalls](v_)
	})
}
func (*OptimizeCalls) Implements_jscomp_CompilerPass() {}

func (this_ *OptimizeCalls) initFields_OptimizeCalls() {
	_ = this_
}

// Java constructor: OptimizeCalls(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.OptimizeCalls.CallGraphCompilerPass> passes, boolean considerExterns)
func NewOptimizeCalls(compiler AbstractCompiler, passes j2gort.Ext_com_google_common_collect_ImmutableList, considerExterns bool) *OptimizeCalls {
	this_ := &OptimizeCalls{}
	InitOptimizeCalls(this_, compiler, passes, considerExterns)
	return this_
}

func InitOptimizeCalls(this_ *OptimizeCalls, compiler AbstractCompiler, passes j2gort.Ext_com_google_common_collect_ImmutableList, considerExterns bool) {
	_ = compiler
	_ = passes
	_ = considerExterns
	this_.initFields_OptimizeCalls()
	this_.compiler = compiler
	this_.passes = j2gort.ImmListToSliceNullable[OptimizeCalls_CallGraphCompilerPass](passes)
	this_.considerExterns = considerExterns
}

// Java static method: OptimizeCalls.builder()
func OptimizeCalls_S_Builder() *OptimizeCalls_Builder {
	return NewOptimizeCalls_Builder()
}

// Java method: OptimizeCalls.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *OptimizeCalls) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	if len(this_.passes) == 0 {
		return
	}
	var val_1 OptimizeCalls_ReferenceMap
	InitOptimizeCalls_ReferenceMap(&val_1)
	var references *OptimizeCalls_ReferenceMap = &val_1
	_ = references
	defer references.ReleaseSlabs_()
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		references.globalScope = nil
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewOptimizeCalls_ReferenceMapBuildingCallback(this_, references))
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).TraverseRoots(externs, root)
	}

	this_.priv_OptimizeCalls_eliminateAccessorsFrom_4075(references)
	{

		for _, pass := range this_.passes {

			_ = pass
			pass.Process_3(externs, root, references)
		}
	}
}

// Java method: OptimizeCalls.eliminateAccessorsFrom(com.google.javascript.jscomp.OptimizeCalls.ReferenceMap references)
func (this_ *OptimizeCalls) priv_OptimizeCalls_eliminateAccessorsFrom_4075(references *OptimizeCalls_ReferenceMap) {
	_ = references
	{
		cur_4 := j2gort.NewMapKeyStrCursor(this_.compiler.GetAccessorSummary().GetAccessors())
		for cur_4.HasNext() {
			accessorName := cur_4.NextStr()
			_ = accessorName
			references.RemovePropReferences(accessorName)
		}
	}
}

// Java static method: OptimizeCalls.safeSet(java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> set)
func OptimizeCalls_SafeSet(set j2gort.Ext_java_util_Set) j2gort.Ext_java_util_Set {
	_ = set
	if set != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](set)))
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
}

// Java static method: OptimizeCalls.mayBeOptimizableName(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String name)
func OptimizeCalls_MayBeOptimizableName(compiler AbstractCompiler, name string) bool {
	_ = compiler
	_ = name
	if compiler.GetCodingConvention().IsExported_1(name) {
		return false
	}
	{
		var iifeCond_1 bool
		{

			sr_5 := name
			se_6 := "JSCompiler_renameProperty"
			iifeCond_1 = j2gort.CheckNotNullString(sr_5) == se_6
		}
		if ((iifeCond_1 || (j2gort.CheckNotNullString(name) == "inherits")) || (j2gort.CheckNotNullString(name) == "$jscomp$inherits")) || (j2gort.CheckNotNullString(name) == "goog$inherits") {
			return false
		}
	}
	return true
}

// Java static method: OptimizeCalls.isAllowedReference(com.google.javascript.rhino.Node n)
func OptimizeCalls_IsAllowedReference(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		return (parent.GetSecondChild() == n)
	case int32(45) /* Token.INSTANCEOF */, int32(25) /* Token.TYPEOF */, int32(44) /* Token.IN */ :
		return true
	case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if ((n == parent.GetFirstChild()) && (grandparent != 0)) && ((grandparent.IsCall() || grandparent.IsOptChainCall()) || grandparent.IsTaggedTemplateLit()) {
			return false
		}
		return true
	case int32(105) /* Token.CLASS */ :
		if n.IsFirstChildOf(parent) {
			return false
		} else {
			if !(n.IsSecondChildOf(parent)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
			}
			if NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent) == 0 {
				return false
			} else {
				return true
			}
		}
	default:
		if NodeUtil_IsNameDeclaration(parent) && !(n.HasChildren()) {
			return true
		}
	}
	return false
}

// Java interface: com.google.javascript.jscomp.OptimizeCalls.CallGraphCompilerPass
type OptimizeCalls_CallGraphCompilerPass interface {
	Implements_jscomp_OptimizeCalls_CallGraphCompilerPass()
	IsJavaNil_() bool
	Process_3(pkg_rhino.Node, pkg_rhino.Node, *OptimizeCalls_ReferenceMap)
}

func ToIface_OptimizeCalls_CallGraphCompilerPass[T interface {
	OptimizeCalls_CallGraphCompilerPass
	comparable
}](v T) OptimizeCalls_CallGraphCompilerPass {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_OptimizeCalls_CallGraphCompilerPass[T OptimizeCalls_CallGraphCompilerPass](v T) OptimizeCalls_CallGraphCompilerPass {
	var vi OptimizeCalls_CallGraphCompilerPass = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeCalls.CallGraphCompilerPass", "com.google.javascript.jscomp.OptimizeCalls$CallGraphCompilerPass", func(v_ any) bool {
		return j2gort.InstanceOf[OptimizeCalls_CallGraphCompilerPass](v_)
	})
}

// Java class: com.google.javascript.jscomp.OptimizeCalls.Builder
type OptimizeCalls_Builder struct {
	compiler        AbstractCompiler                                            `gson:"compiler"`
	passes          *j2gort.Ext_com_google_common_collect_ImmutableList_Builder `gson:"passes"`
	considerExterns *j2gort.BoxBoolean                                          `gson:"considerExterns"`
	idHash_         atomic.Int32
}

func (this_ *OptimizeCalls_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeCalls_Builder) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeCalls_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeCalls$Builder", "Builder"
}
func (this_ *OptimizeCalls_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeCalls_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeCalls.Builder", "com.google.javascript.jscomp.OptimizeCalls$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeCalls_Builder](v_)
	})
}

func (this_ *OptimizeCalls_Builder) initFields_OptimizeCalls_Builder() {
	_ = this_
	this_.passes = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
}

// Java constructor: Builder()
func NewOptimizeCalls_Builder() *OptimizeCalls_Builder {
	this_ := &OptimizeCalls_Builder{}
	InitOptimizeCalls_Builder(this_)
	return this_
}

func InitOptimizeCalls_Builder(this_ *OptimizeCalls_Builder) {
	this_.initFields_OptimizeCalls_Builder()
}

// Java method: Builder.setCompiler(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *OptimizeCalls_Builder) SetCompiler(compiler AbstractCompiler) *OptimizeCalls_Builder {
	_ = compiler
	this_.compiler = compiler
	return this_
}

// Java method: Builder.setConsiderExterns(boolean b)
func (this_ *OptimizeCalls_Builder) SetConsiderExterns(b bool) *OptimizeCalls_Builder {
	_ = b
	this_.considerExterns = j2gort.ValueOfBoolean(b)
	return this_
}

// Java method: Builder.addPass(com.google.javascript.jscomp.OptimizeCalls.CallGraphCompilerPass pass)
func (this_ *OptimizeCalls_Builder) AddPass(pass OptimizeCalls_CallGraphCompilerPass) *OptimizeCalls_Builder {
	_ = pass
	this_.passes.AddVoid_(pass)
	return this_
}

// Java method: Builder.build()
func (this_ *OptimizeCalls_Builder) Build() *OptimizeCalls {
	j2gort.CheckNotNullIface(this_.compiler)
	j2gort.CheckNotNullPtr(this_.considerExterns)
	return NewOptimizeCalls(this_.compiler, this_.passes.BuildImmList_(), j2gort.UnboxBoolean(this_.considerExterns))
}

// Java class: com.google.javascript.jscomp.OptimizeCalls.NodeList
type OptimizeCalls_NodeList struct {
	nodes     []pkg_rhino.Node `gson:"nodes"`
	hashTable []pkg_rhino.Node `gson:"hashTable"`
	n0        pkg_rhino.Node   `gson:"n0"`
	n1        pkg_rhino.Node   `gson:"n1"`
	n2        pkg_rhino.Node   `gson:"n2"`
	n3        pkg_rhino.Node   `gson:"n3"`
	hashCount int32            `gson:"hashCount"`
	size      int32            `gson:"size"`
	idHash_   atomic.Int32
}

func (this_ *OptimizeCalls_NodeList) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeCalls_NodeList) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeCalls_NodeList) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeCalls$NodeList", "NodeList"
}
func (this_ *OptimizeCalls_NodeList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeCalls_NodeList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeCalls.NodeList", "com.google.javascript.jscomp.OptimizeCalls$NodeList", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeCalls_NodeList](v_)
	})
}

var static_OptimizeCalls_NodeList_EMPTY *OptimizeCalls_NodeList

func StaticGet_OptimizeCalls_NodeList_EMPTY() *OptimizeCalls_NodeList {
	clinit_OptimizeCalls_NodeList()
	return static_OptimizeCalls_NodeList_EMPTY
}
func StaticSet_OptimizeCalls_NodeList_EMPTY(v *OptimizeCalls_NodeList) {
	clinit_OptimizeCalls_NodeList()
	static_OptimizeCalls_NodeList_EMPTY = v
}

var clinitDone_OptimizeCalls_NodeList uint32
var clinitState_OptimizeCalls_NodeList j2gort.ClinitState

func clinit_OptimizeCalls_NodeList() {
	if atomic.LoadUint32(&clinitDone_OptimizeCalls_NodeList) != 0 {
		return
	}
	clinitSlow_OptimizeCalls_NodeList()
}

func clinitSlow_OptimizeCalls_NodeList() {
	if !j2gort.EnterClinit(&clinitState_OptimizeCalls_NodeList, &clinitDone_OptimizeCalls_NodeList) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_OptimizeCalls_NodeList, &clinitDone_OptimizeCalls_NodeList)
	static_OptimizeCalls_NodeList_EMPTY = NewOptimizeCalls_NodeList()
}

func (this_ *OptimizeCalls_NodeList) initFields_OptimizeCalls_NodeList() {
	_ = this_
}

// Java constructor: NodeList()
func NewOptimizeCalls_NodeList() *OptimizeCalls_NodeList {
	clinit_OptimizeCalls_NodeList()
	this_ := &OptimizeCalls_NodeList{}
	InitOptimizeCalls_NodeList(this_)
	return this_
}
func allocPooled_OptimizeCalls_NodeList(slab_ *j2gort.PooledSlab[OptimizeCalls_NodeList]) *OptimizeCalls_NodeList {
	clinit_OptimizeCalls_NodeList()
	this_ := slab_.Alloc()
	InitOptimizeCalls_NodeList(this_)
	return this_
}

func InitOptimizeCalls_NodeList(this_ *OptimizeCalls_NodeList) {
	this_.initFields_OptimizeCalls_NodeList()
	this_.n0 = pkg_rhino.Node(0)
	this_.n1 = pkg_rhino.Node(0)
	this_.n2 = pkg_rhino.Node(0)
	this_.n3 = pkg_rhino.Node(0)
	this_.nodes = nil
	this_.size = int32(0)
}

// Java constructor: NodeList(int initialCapacity)
func NewOptimizeCalls_NodeList_1(initialCapacity int32) *OptimizeCalls_NodeList {
	clinit_OptimizeCalls_NodeList()
	this_ := &OptimizeCalls_NodeList{}
	InitOptimizeCalls_NodeList_1(this_, initialCapacity)
	return this_
}

func InitOptimizeCalls_NodeList_1(this_ *OptimizeCalls_NodeList, initialCapacity int32) {
	_ = initialCapacity
	this_.initFields_OptimizeCalls_NodeList()
	this_.n0 = pkg_rhino.Node(0)
	this_.n1 = pkg_rhino.Node(0)
	this_.n2 = pkg_rhino.Node(0)
	this_.n3 = pkg_rhino.Node(0)
	if initialCapacity > int32(4) {
		this_.nodes = make([]pkg_rhino.Node, initialCapacity)
	} else {
		this_.nodes = nil
	}
	this_.size = int32(0)
}

// Java method: NodeList.add(com.google.javascript.rhino.Node n)
func (this_ *OptimizeCalls_NodeList) Add(n pkg_rhino.Node) {
	_ = n
	var s int32 = this_.size
	_ = s
	var arr []pkg_rhino.Node = this_.nodes
	_ = arr
	if arr == nil {
		switch s {
		case int32(0):
			this_.n0 = n
			this_.size = int32(1)
			return
		case int32(1):
			this_.n1 = n
			this_.size = int32(2)
			return
		case int32(2):
			this_.n2 = n
			this_.size = int32(3)
			return
		case int32(3):
			this_.n3 = n
			this_.size = int32(4)
			return
		default:
			arr = make([]pkg_rhino.Node, int32(8))
			{
				a_7 := arr
				i_8 := int32(0)
				a_7[i_8] = this_.n0
			}
			{
				a_9 := arr
				i_10 := int32(1)
				a_9[i_10] = this_.n1
			}
			{
				a_11 := arr
				i_12 := int32(2)
				a_11[i_12] = this_.n2
			}
			{
				a_13 := arr
				i_14 := int32(3)
				a_13[i_14] = this_.n3
			}
			{
				a_15 := arr
				i_16 := int32(4)
				a_15[i_16] = n
			}
			this_.n0 = pkg_rhino.Node(0)
			this_.n1 = pkg_rhino.Node(0)
			this_.n2 = pkg_rhino.Node(0)
			this_.n3 = pkg_rhino.Node(0)
			this_.nodes = arr
			this_.size = int32(5)
			return
		}
	}
	if s == int32(len(arr)) {
		{

			arrSrc_17 := arr
			arrLen_18 := max(int32(8), (s * int32(2)))
			j2gort.CheckNotNullSlice(arrSrc_17)
			if arrLen_18 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_19 := make([]pkg_rhino.Node, arrLen_18)
			copy(arrOut_19, arrSrc_17)
			arr = arrOut_19
		}

		this_.nodes = arr
	}
	{
		a_20 := arr
		i_21 := s
		a_20[i_21] = n
	}
	this_.size = (s + int32(1))
	if (this_.hashTable != nil) && (n != 0) {
		this_.priv_NodeList_insertIntoHashTable_4081(n)
	}
}

// Java static method: NodeList.mixNodeHash(com.google.javascript.rhino.Node n)
func OptimizeCalls_NodeList_MixNodeHash(n pkg_rhino.Node) int32 {
	clinit_OptimizeCalls_NodeList()
	_ = n
	var h int32 = (int32(j2gort.CheckNotNullNode(n)) * int32(-1640531527))
	_ = h
	return (h ^ (javalang.Ushr32(h, int(int32(16)))))
}

// Java method: NodeList.insertIntoHashTable(com.google.javascript.rhino.Node n)
func (this_ *OptimizeCalls_NodeList) priv_NodeList_insertIntoHashTable_4081(n pkg_rhino.Node) {
	_ = n
	var table []pkg_rhino.Node = this_.hashTable
	_ = table
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (OptimizeCalls_NodeList_MixNodeHash(n) & mask)
	_ = slot
	for {
		var existing pkg_rhino.Node = table[slot]
		_ = existing
		if existing == 0 {
			if ((this_.hashCount + int32(1)) * int32(2)) > int32(len(table)) {
				this_.priv_NodeList_growHashTable_4082()
				table = this_.hashTable
				mask = (int32(len(table)) - int32(1))
				slot = (OptimizeCalls_NodeList_MixNodeHash(n) & mask)
				for table[slot] != 0 {
					slot = ((slot + int32(1)) & mask)
				}
			}
			{
				a_22 := table
				i_23 := slot
				a_22[i_23] = n
			}
			this_.hashCount = this_.hashCount + 1
			return
		}
		if existing == n {
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: NodeList.growHashTable()
func (this_ *OptimizeCalls_NodeList) priv_NodeList_growHashTable_4082() {
	var oldTable []pkg_rhino.Node = this_.hashTable
	_ = oldTable
	var newCap int32 = (int32(len(oldTable)) * int32(2))
	_ = newCap
	var newTable []pkg_rhino.Node = make([]pkg_rhino.Node, newCap)
	_ = newTable
	var mask int32 = (newCap - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(oldTable))
		_ = len_
		for ; i < len_; i = i + 1 {
			var elem pkg_rhino.Node = oldTable[i]
			_ = elem
			if elem != 0 {
				var slot int32 = (OptimizeCalls_NodeList_MixNodeHash(elem) & mask)
				_ = slot
				for newTable[slot] != 0 {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_24 := newTable
					i_25 := slot
					a_24[i_25] = elem
				}
			}
		}
	}
	this_.hashTable = newTable
}

// Java method: NodeList.buildHashTable()
func (this_ *OptimizeCalls_NodeList) priv_NodeList_buildHashTable_4083() []pkg_rhino.Node {
	var s int32 = this_.size
	_ = s
	var cap_ int32 = int32(64)
	_ = cap_
	for cap_ < (s * int32(2)) {
		cap_ = javalang.Shl32(cap_, int(int32(1)))
	}
	var table []pkg_rhino.Node = make([]pkg_rhino.Node, cap_)
	_ = table
	var mask int32 = (cap_ - int32(1))
	_ = mask
	var unique int32 = int32(0)
	_ = unique
	var arr []pkg_rhino.Node = this_.nodes
	_ = arr
	{
		var i int32 = int32(0)
		_ = i
		for ; i < s; i = i + 1 {
			var elem pkg_rhino.Node = arr[i]
			_ = elem
			if elem != 0 {
				var slot int32 = (OptimizeCalls_NodeList_MixNodeHash(elem) & mask)
				_ = slot
				for {
					var existing pkg_rhino.Node = table[slot]
					_ = existing
					if existing == 0 {
						{
							a_26 := table
							i_27 := slot
							a_26[i_27] = elem
						}
						unique = unique + 1
						break
					}
					if existing == elem {
						break
					}
					slot = ((slot + int32(1)) & mask)
				}
			}
		}
	}
	this_.hashTable = table
	this_.hashCount = unique
	return table
}

// Java method: NodeList.addAll(com.google.javascript.jscomp.OptimizeCalls.NodeList other)
func (this_ *OptimizeCalls_NodeList) AddAll(other *OptimizeCalls_NodeList) {
	_ = other
	var count int32 = other.size
	_ = count
	if count == int32(0) {
		return
	}
	var curSize int32 = this_.size
	_ = curSize
	var newSize int32 = (curSize + count)
	_ = newSize
	if (this_.nodes == nil) && (newSize <= int32(4)) {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < count; i = i + 1 {
				this_.Add(other.Get(i))
			}
		}
		return
	}
	var arr []pkg_rhino.Node = this_.nodes
	_ = arr
	if arr == nil {
		arr = make([]pkg_rhino.Node, max(int32(8), newSize))
		if curSize > int32(0) {
			{
				a_28 := arr
				i_29 := int32(0)
				a_28[i_29] = this_.n0
			}
			if curSize > int32(1) {
				{
					a_30 := arr
					i_31 := int32(1)
					a_30[i_31] = this_.n1
				}
				if curSize > int32(2) {
					{
						a_32 := arr
						i_33 := int32(2)
						a_32[i_33] = this_.n2
					}
					if curSize > int32(3) {
						{
							a_34 := arr
							i_35 := int32(3)
							a_34[i_35] = this_.n3
						}
					}
				}
			}
		}
		this_.n0 = pkg_rhino.Node(0)
		this_.n1 = pkg_rhino.Node(0)
		this_.n2 = pkg_rhino.Node(0)
		this_.n3 = pkg_rhino.Node(0)
		this_.nodes = arr
	} else if newSize > int32(len(arr)) {
		{

			arrSrc_36 := arr
			arrLen_37 := max((int32(len(arr)) * int32(2)), newSize)
			j2gort.CheckNotNullSlice(arrSrc_36)
			if arrLen_37 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_38 := make([]pkg_rhino.Node, arrLen_37)
			copy(arrOut_38, arrSrc_36)
			arr = arrOut_38
		}

		this_.nodes = arr
	}
	var otherArr []pkg_rhino.Node = other.nodes
	_ = otherArr
	if otherArr != nil {
		func() {
			acSrc_39 := otherArr
			acSP_40 := int32(0)
			acDst_41 := arr
			acDP_42 := curSize
			acLen_43 := count
			j2gort.CheckNotNullSlice(acSrc_39)
			j2gort.CheckNotNullSlice(acDst_41)
			if acSP_40 < 0 || acDP_42 < 0 || acLen_43 < 0 || int64(acSP_40)+int64(acLen_43) > int64(len(acSrc_39)) || int64(acDP_42)+int64(acLen_43) > int64(len(acDst_41)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_41[int(acDP_42):int(acDP_42)+int(acLen_43)], acSrc_39[int(acSP_40):int(acSP_40)+int(acLen_43)])
		}()
	} else {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < count; i_2 = i_2 + 1 {
				{
					a_44 := arr
					i_45 := (curSize + i_2)
					a_44[i_45] = other.Get(i_2)
				}
			}
		}
	}
	this_.size = newSize
	if this_.hashTable != nil {
		if ((this_.hashCount + count) * int32(2)) > int32(len(this_.hashTable)) {
			this_.hashTable = nil
			this_.hashCount = int32(0)
		} else {
			{
				var i_3 int32 = curSize
				_ = i_3
				for ; i_3 < newSize; i_3 = i_3 + 1 {
					var elem pkg_rhino.Node = arr[i_3]
					_ = elem
					if elem != 0 {
						this_.priv_NodeList_insertIntoHashTable_4081(elem)
					}
				}
			}
		}
	}
}

// Java method: NodeList.clear()
func (this_ *OptimizeCalls_NodeList) Clear() {
	if this_.nodes != nil {
		j2gort.ArraysFillRange(this_.nodes, int32(0), this_.size, pkg_rhino.Node(0))
	} else {
		this_.n0 = pkg_rhino.Node(0)
		this_.n1 = pkg_rhino.Node(0)
		this_.n2 = pkg_rhino.Node(0)
		this_.n3 = pkg_rhino.Node(0)
	}
	this_.hashTable = nil
	this_.hashCount = int32(0)
	this_.size = int32(0)
}

// Java method: NodeList.size()
func (this_ *OptimizeCalls_NodeList) Size() int32 {
	return this_.size
}

// Java method: NodeList.isEmpty()
func (this_ *OptimizeCalls_NodeList) IsEmpty() bool {
	return (this_.size == int32(0))
}

// Java method: NodeList.contains(com.google.javascript.rhino.Node target)
func (this_ *OptimizeCalls_NodeList) Contains(target pkg_rhino.Node) bool {
	_ = target
	var s int32 = this_.size
	_ = s
	if s == int32(0) {
		return false
	}
	var arr []pkg_rhino.Node = this_.nodes
	_ = arr
	if arr != nil {
		if (target != 0) && (s > int32(16)) {
			var table []pkg_rhino.Node = this_.hashTable
			_ = table
			if table == nil {
				table = this_.priv_NodeList_buildHashTable_4083()
			}
			var mask int32 = (int32(len(table)) - int32(1))
			_ = mask
			var slot int32 = (OptimizeCalls_NodeList_MixNodeHash(target) & mask)
			_ = slot
			for {
				var existing pkg_rhino.Node = table[slot]
				_ = existing
				if existing == 0 {
					return false
				}
				if existing == target {
					return true
				}
				slot = ((slot + int32(1)) & mask)
			}
		}
		{
			var i int32 = int32(0)
			_ = i
			for ; i < s; i = i + 1 {
				if arr[i] == target {
					return true
				}
			}
		}
		return false
	}
	if this_.n0 == target {
		return true
	}
	if (s > int32(1)) && (this_.n1 == target) {
		return true
	}
	if (s > int32(2)) && (this_.n2 == target) {
		return true
	}
	return ((s > int32(3)) && (this_.n3 == target))
}

// Java method: NodeList.get(int index)
func (this_ *OptimizeCalls_NodeList) Get(index int32) pkg_rhino.Node {
	_ = index
	var arr []pkg_rhino.Node = this_.nodes
	_ = arr
	if arr != nil {
		return arr[index]
	}
	switch index {
	case int32(0):
		return this_.n0
	case int32(1):
		return this_.n1
	case int32(2):
		return this_.n2
	case int32(3):
		return this_.n3
	default:
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	panic("unreachable switch expression")
}

// Java class: com.google.javascript.jscomp.OptimizeCalls.ReferenceMap
type OptimizeCalls_ReferenceMap struct {
	globalScope                 *Scope                                             `gson:"globalScope"`
	names                       j2gort.OrderedMap[string, *OptimizeCalls_NodeList] `gson:"names"`
	props                       j2gort.OrderedMap[string, *OptimizeCalls_NodeList] `gson:"props"`
	idHash_                     atomic.Int32
	slab_OptimizeCalls_NodeList j2gort.PooledSlab[OptimizeCalls_NodeList]
}

func (this_ *OptimizeCalls_ReferenceMap) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeCalls_ReferenceMap) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeCalls_ReferenceMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeCalls$ReferenceMap", "ReferenceMap"
}
func (this_ *OptimizeCalls_ReferenceMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeCalls_ReferenceMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeCalls.ReferenceMap", "com.google.javascript.jscomp.OptimizeCalls$ReferenceMap", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeCalls_ReferenceMap](v_)
	})
}

func (this_ *OptimizeCalls_ReferenceMap) initFields_OptimizeCalls_ReferenceMap() {
	_ = this_
	this_.names = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList]()
	this_.props = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList](

	// Java constructor: ReferenceMap()
	)
}

func NewOptimizeCalls_ReferenceMap() *OptimizeCalls_ReferenceMap {
	this_ := &OptimizeCalls_ReferenceMap{}
	InitOptimizeCalls_ReferenceMap(this_)
	return this_
}

func InitOptimizeCalls_ReferenceMap(this_ *OptimizeCalls_ReferenceMap) {
	this_.initFields_OptimizeCalls_ReferenceMap()
}

// Java method: ReferenceMap.addNameReference(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *OptimizeCalls_ReferenceMap) AddNameReference(name string, n pkg_rhino.Node) {
	_ = name
	_ = n
	var refs *OptimizeCalls_NodeList = this_.names.Get(name)
	_ = refs
	if refs == nil {
		refs = allocPooled_OptimizeCalls_NodeList(&this_.slab_OptimizeCalls_NodeList)
		this_.names.Put(name, refs)
	}
	refs.Add(n)
}

// Java method: ReferenceMap.addPropReference(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *OptimizeCalls_ReferenceMap) AddPropReference(name string, n pkg_rhino.Node) {
	_ = name
	_ = n
	var refs *OptimizeCalls_NodeList = this_.props.Get(name)
	_ = refs
	if refs == nil {
		refs = allocPooled_OptimizeCalls_NodeList(&this_.slab_OptimizeCalls_NodeList)
		this_.props.Put(name, refs)
	}
	refs.Add(n)
}

// Java method: ReferenceMap.removePropReferences(java.lang.String name)
func (this_ *OptimizeCalls_ReferenceMap) RemovePropReferences(name string) {
	_ = name
	this_.props.Remove(name)
}

// Java method: ReferenceMap.getGlobalScope()
func (this_ *OptimizeCalls_ReferenceMap) GetGlobalScope() *Scope {
	return this_.globalScope
}

// Java static method: ReferenceMap.getFunctionNodeList(com.google.javascript.jscomp.OptimizeCalls.NodeList definitionSites)
func OptimizeCalls_ReferenceMap_GetFunctionNodeList(definitionSites *OptimizeCalls_NodeList) *OptimizeCalls_NodeList {
	_ = definitionSites
	var result *OptimizeCalls_NodeList = nil
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = definitionSites.Size()
		_ = count
		for ; i < count; i = i + 1 {
			result = OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(definitionSites.Get(i), result)
		}
	}
	if result != nil {
		return result
	} else {
		return StaticGet_OptimizeCalls_NodeList_EMPTY()
	}
}
func allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList(slab_ *j2gort.PooledSlab[OptimizeCalls_NodeList], definitionSites *OptimizeCalls_NodeList) *OptimizeCalls_NodeList {
	_ = definitionSites
	var result *OptimizeCalls_NodeList = nil
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = definitionSites.Size()
		_ = count
		for ; i < count; i = i + 1 {
			result = allocPooledVia_OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(slab_, definitionSites.Get(i), result)
		}
	}
	if result != nil {
		return result
	} else {
		return StaticGet_OptimizeCalls_NodeList_EMPTY()
	}
}

// Java static method: ReferenceMap.definitionFunctionNodeListFor(com.google.javascript.rhino.Node definitionSite)
func OptimizeCalls_ReferenceMap_DefinitionFunctionNodeListFor(definitionSite pkg_rhino.Node) *OptimizeCalls_NodeList {
	_ = definitionSite
	var result *OptimizeCalls_NodeList = OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(definitionSite, nil)
	_ = result
	if result != nil {
		return result
	} else {
		return StaticGet_OptimizeCalls_NodeList_EMPTY()
	}
}
func allocPooledVia_OptimizeCalls_ReferenceMap_DefinitionFunctionNodeListFor(slab_ *j2gort.PooledSlab[OptimizeCalls_NodeList], definitionSite pkg_rhino.Node) *OptimizeCalls_NodeList {
	_ = definitionSite
	var result *OptimizeCalls_NodeList = allocPooledVia_OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(slab_, definitionSite, nil)
	_ = result
	if result != nil {
		return result
	} else {
		return StaticGet_OptimizeCalls_NodeList_EMPTY()
	}
}

// Java static method: ReferenceMap.appendDefinitionFunctionNodes(com.google.javascript.rhino.Node definitionSite, com.google.javascript.jscomp.OptimizeCalls.@org.jspecify.annotations.Nullable NodeList fns)
func OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(definitionSite pkg_rhino.Node, fns *OptimizeCalls_NodeList) *OptimizeCalls_NodeList {
	_ = definitionSite
	_ = fns
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		return fns
	}
	var parent pkg_rhino.Node = definitionSite.GetParent()
	_ = parent
	if parent == 0 {
		return fns
	}
	switch parent.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		if definitionSite.IsFirstChildOf(parent) {
			var constructorFnDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent)
			_ = constructorFnDef
			if constructorFnDef != 0 {
				if fns == nil {
					fns = NewOptimizeCalls_NodeList()
				}
				fns.Add(constructorFnDef.GetOnlyChild())
			}
		}
	case int32(73) /* Token.FUNCTION */ :
		if fns == nil {
			fns = NewOptimizeCalls_NodeList()
		}
		fns.Add(parent)
	case int32(106) /* Token.CLASS_MEMBERS */ :
		if definitionSite.IsMemberFunctionDef() {
			if fns == nil {
				fns = NewOptimizeCalls_NodeList()
			}
			fns.Add(definitionSite.GetLastChild())
		} else {
			if !(definitionSite.IsMemberFieldDef()) {
				j2gort.CheckArgumentMsg(false, j2gort.BoxNode(definitionSite))
			}
			var value pkg_rhino.Node = definitionSite.GetFirstChild()
			_ = value
			if value != 0 {
				fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, value)
			}
		}
	case int32(47) /* Token.OBJECTLIT */ :
		if !(definitionSite.IsStringKey() || definitionSite.IsMemberFunctionDef()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(definitionSite))
		}
		fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, definitionSite.GetLastChild())
	case int32(51) /* Token.ASSIGN */ :
		var target pkg_rhino.Node = parent.GetFirstChild()
		_ = target
		var value_2 pkg_rhino.Node = parent.GetLastChild()
		_ = value_2
		if definitionSite == target {
			fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, value_2)
		}
	case int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(84) /* Token.VAR */ :
		if definitionSite.IsName() && definitionSite.HasChildren() {
			fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, definitionSite.GetFirstChild())
		}
	default:
	}
	return fns
}
func allocPooledVia_OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(slab_ *j2gort.PooledSlab[OptimizeCalls_NodeList], definitionSite pkg_rhino.Node, fns *OptimizeCalls_NodeList) *OptimizeCalls_NodeList {
	_ = definitionSite
	_ = fns
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		return fns
	}
	var parent pkg_rhino.Node = definitionSite.GetParent()
	_ = parent
	if parent == 0 {
		return fns
	}
	switch parent.TokenSwitchOrd_() {
	case int32(105):
		if definitionSite.IsFirstChildOf(parent) {
			var constructorFnDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent)
			_ = constructorFnDef
			if constructorFnDef != 0 {
				if fns == nil {
					fns = allocPooled_OptimizeCalls_NodeList(slab_)
				}
				fns.Add(constructorFnDef.GetOnlyChild())
			}
		}
	case int32(73):
		if fns == nil {
			fns = allocPooled_OptimizeCalls_NodeList(slab_)
		}
		fns.Add(parent)
	case int32(106):
		if definitionSite.IsMemberFunctionDef() {
			if fns == nil {
				fns = allocPooled_OptimizeCalls_NodeList(slab_)
			}
			fns.Add(definitionSite.GetLastChild())
		} else {
			if !(definitionSite.IsMemberFieldDef()) {
				j2gort.CheckArgumentMsg(false, j2gort.BoxNode(definitionSite))
			}
			var value pkg_rhino.Node = definitionSite.GetFirstChild()
			_ = value
			if value != 0 {
				fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, value)
			}
		}
	case int32(47):
		if !(definitionSite.IsStringKey() || definitionSite.IsMemberFunctionDef()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(definitionSite))
		}
		fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, definitionSite.GetLastChild())
	case int32(51):
		var target pkg_rhino.Node = parent.GetFirstChild()
		_ = target
		var value_2 pkg_rhino.Node = parent.GetLastChild()
		_ = value_2
		if definitionSite == target {
			fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, value_2)
		}
	case int32(97), int32(111), int32(84):
		if definitionSite.IsName() && definitionSite.HasChildren() {
			fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, definitionSite.GetFirstChild())
		}
	default:
	}
	return fns
}

// Java static method: ReferenceMap.addValueFunctionNodes(com.google.javascript.jscomp.OptimizeCalls.@org.jspecify.annotations.Nullable NodeList fns, com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns *OptimizeCalls_NodeList, n pkg_rhino.Node) *OptimizeCalls_NodeList {
	_ = fns
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		var constructorFnDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		_ = constructorFnDef
		if constructorFnDef != 0 {
			if fns == nil {
				fns = NewOptimizeCalls_NodeList()
			}
			fns.Add(constructorFnDef.GetOnlyChild())
		}
	case int32(73) /* Token.FUNCTION */ :
		if fns == nil {
			fns = NewOptimizeCalls_NodeList()
		}
		fns.Add(n)
	case int32(67) /* Token.HOOK */ :
		fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetSecondChild())
		fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetLastChild())
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */ :
		fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetFirstChild())
		fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetLastChild())
	case int32(101) /* Token.CAST */, int32(50) /* Token.COMMA */ :
		fns = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetLastChild())
	default:
	}
	return fns
}
func allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_ *j2gort.PooledSlab[OptimizeCalls_NodeList], fns *OptimizeCalls_NodeList, n pkg_rhino.Node) *OptimizeCalls_NodeList {
	_ = fns
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(105):
		var constructorFnDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		_ = constructorFnDef
		if constructorFnDef != 0 {
			if fns == nil {
				fns = allocPooled_OptimizeCalls_NodeList(slab_)
			}
			fns.Add(constructorFnDef.GetOnlyChild())
		}
	case int32(73):
		if fns == nil {
			fns = allocPooled_OptimizeCalls_NodeList(slab_)
		}
		fns.Add(n)
	case int32(67):
		fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, n.GetSecondChild())
		fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, n.GetLastChild())
	case int32(68), int32(69), int32(70):
		fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, n.GetFirstChild())
		fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, n.GetLastChild())
	case int32(101), int32(50):
		fns = allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes(slab_, fns, n.GetLastChild())
	default:
	}
	return fns
}

// Java static method: ReferenceMap.isNormalOrOptChainCallOrNewTarget(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n pkg_rhino.Node) bool {
	_ = n
	return ((OptimizeCalls_ReferenceMap_IsCallTarget(n) || OptimizeCalls_ReferenceMap_IsNewTarget(n)) || OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n))
}

// Java static method: ReferenceMap.isCallTarget(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_IsCallTarget(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsCall() && n.IsFirstChildOf(parent) {
		return true
	}
	var grandParent pkg_rhino.Node = parent.GetParent()
	_ = grandParent
	return (((parent.IsGetProp() && grandParent.IsCall()) && parent.IsFirstChildOf(grandParent)) && (j2gort.CheckNotNullString(parent.GetString()) == "call"))
}

// Java static method: ReferenceMap.isOptChainCallTarget(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsOptChainCall() && n.IsFirstChildOf(parent) {
		return true
	}
	var grandParent pkg_rhino.Node = parent.GetParent()
	_ = grandParent
	return (((parent.IsOptChainGetProp() && grandParent.IsOptChainCall()) && parent.IsFirstChildOf(grandParent)) && (j2gort.CheckNotNullString(parent.GetString()) == "call"))
}

// Java static method: ReferenceMap.isNewTarget(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_IsNewTarget(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (parent.IsNew() && (parent.GetFirstChild() == n))
}

// Java static method: ReferenceMap.getCallOrNewNodeForTarget(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var maybeCall pkg_rhino.Node = n.GetParent()
	_ = maybeCall
	if !(n.IsFirstChildOf(maybeCall)) {
		j2gort.CheckStateFmt(false, "%s\n\n%s", j2gort.BoxNode(maybeCall), j2gort.BoxNode(n))
	}
	if NodeUtil_IsCallOrNew(maybeCall) {
		return maybeCall
	} else {
		var child pkg_rhino.Node = maybeCall
		_ = child
		maybeCall = child.GetParent()
		if !(NodeUtil_IsNormalOrOptChainGetProp(child)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(child))
		}
		if !(NodeUtil_IsNormalOrOptChainCall(maybeCall)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(maybeCall))
		}
		if !(child.IsFirstChildOf(maybeCall)) {
			j2gort.CheckStateFmt(false, "%s\n\n%s", j2gort.BoxNode(maybeCall), j2gort.BoxNode(child))
		}
		return maybeCall
	}
}

// Java static method: ReferenceMap.getFirstArgumentForCallOrNewOrDotCall(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	return OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(n, int32(0))
}

// Java static method: ReferenceMap.getArgumentForCallOrNewOrDotCall(com.google.javascript.rhino.Node n, int index)
func OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(n pkg_rhino.Node, index int32) pkg_rhino.Node {
	_ = n
	_ = index
	var adjustedIndex int32 = index
	_ = adjustedIndex
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if !((parent.IsCall() || parent.IsOptChainCall()) || parent.IsNew()) {
		parent = parent.GetParent()
		if NodeUtil_IsFunctionObjectCall(parent) {
			adjustedIndex = adjustedIndex + 1
		}
	}
	return NodeUtil_GetArgumentForCallOrNew(parent, adjustedIndex)
}

// Java static method: ReferenceMap.isSimpleAssignmentTarget(com.google.javascript.rhino.Node n)
func OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((parent.IsAssign() && n.IsFirstChildOf(parent)) && parent.GetParent().IsExprResult())
}

// Java class: com.google.javascript.jscomp.OptimizeCalls.ReferenceMapBuildingCallback
type OptimizeCalls_ReferenceMapBuildingCallback struct {
	outer_      *OptimizeCalls
	externProps j2gort.Ext_java_util_Set    `gson:"externProps"`
	references  *OptimizeCalls_ReferenceMap `gson:"references"`
	globalScope *Scope                      `gson:"globalScope"`
	idHash_     atomic.Int32
}

func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeCalls_ReferenceMapBuildingCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeCalls$ReferenceMapBuildingCallback", "ReferenceMapBuildingCallback"
}
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeCalls_ReferenceMapBuildingCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeCalls.ReferenceMapBuildingCallback", "com.google.javascript.jscomp.OptimizeCalls$ReferenceMapBuildingCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeCalls_ReferenceMapBuildingCallback](v_)
	})
}
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) SetOuter_(o *OptimizeCalls) {
	this_.outer_ = o
}
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) GetOuter_() *OptimizeCalls {
	return this_.outer_
}
func (*OptimizeCalls_ReferenceMapBuildingCallback) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*OptimizeCalls_ReferenceMapBuildingCallback) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) initFields_OptimizeCalls_ReferenceMapBuildingCallback() {
	_ = this_
}

// Java constructor: ReferenceMapBuildingCallback(com.google.javascript.jscomp.OptimizeCalls.ReferenceMap references)
func NewOptimizeCalls_ReferenceMapBuildingCallback(outer_ *OptimizeCalls, references *OptimizeCalls_ReferenceMap) *OptimizeCalls_ReferenceMapBuildingCallback {
	this_ := &OptimizeCalls_ReferenceMapBuildingCallback{}
	this_.outer_ = outer_
	InitOptimizeCalls_ReferenceMapBuildingCallback(this_, references)
	return this_
}

func InitOptimizeCalls_ReferenceMapBuildingCallback(this_ *OptimizeCalls_ReferenceMapBuildingCallback, references *OptimizeCalls_ReferenceMap) {
	_ = references
	this_.initFields_OptimizeCalls_ReferenceMapBuildingCallback()
	this_.externProps = OptimizeCalls_SafeSet(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.outer_.compiler.GetExternProperties()))
	this_.references = references
}

// Java method: ReferenceMapBuildingCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node unused)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, unused pkg_rhino.Node) {
	_ = t
	_ = n
	_ = unused
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_ReferenceMapBuildingCallback_maybeAddNameReference_4095(n.GetString(), n)
	case int32(29) /* Token.OPTCHAIN_GETPROP */, int32(26) /* Token.GETPROP */ :
		this_.priv_ReferenceMapBuildingCallback_maybeAddPropReference_4096(n.GetString(), n)
	case int32(28) /* Token.CALL */ :
		var fnName pkg_rhino.Node = n.GetFirstChild()
		_ = fnName
		if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(fnName) {
			var propName pkg_rhino.Node = NodeUtil_GetArgumentForCallOrNew(n, int32(0))
			_ = propName
			if propName != 0 {
				this_.priv_ReferenceMapBuildingCallback_maybeAddPropReference_4096(propName.GetString(), n)
			}
		}
	case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if !(n.IsQuotedStringKey()) {
			this_.priv_ReferenceMapBuildingCallback_maybeAddPropReference_4096(n.GetString(), n)
		}
	case int32(110) /* Token.SUPER */ :
		this_.priv_ReferenceMapBuildingCallback_visitSuper_4094(n)
	case int32(129) /* Token.COMPUTED_PROP */, int32(109) /* Token.COMPUTED_FIELD_DEF */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(27) /* Token.GETELEM */, int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */ :
	default:
	}
}

// Java method: ReferenceMapBuildingCallback.visitSuper(com.google.javascript.rhino.Node superNode)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) priv_ReferenceMapBuildingCallback_visitSuper_4094(superNode pkg_rhino.Node) {
	_ = superNode
	var parent pkg_rhino.Node = superNode.GetParent()
	_ = parent
	if parent.IsCall() && superNode.IsFirstChildOf(parent) {
		var enclosingClass pkg_rhino.Node = j2gort.CheckNotNullNode(NodeUtil_GetEnclosingClass(parent))
		_ = enclosingClass
		var extendsNode pkg_rhino.Node = enclosingClass.GetSecondChild()
		_ = extendsNode
		if !(!(extendsNode.IsEmpty())) {
			j2gort.CheckStateStr(false, "super call appears in class without extends clause")
		}
		if extendsNode.IsName() {
			this_.priv_ReferenceMapBuildingCallback_maybeAddNameReference_4095(extendsNode.GetString(), superNode)
		} else if extendsNode.IsGetProp() {
			this_.priv_ReferenceMapBuildingCallback_maybeAddPropReference_4096(extendsNode.GetString(), superNode)
		}
	}
}

// Java method: ReferenceMapBuildingCallback.maybeAddNameReference(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) priv_ReferenceMapBuildingCallback_maybeAddNameReference_4095(name string, n pkg_rhino.Node) {
	_ = name
	_ = n
	var var_ *Var = j2gort.CastPtr[Var](this_.globalScope.GetSlot(name))
	_ = var_
	if (var_ != nil) && (this_.outer_.considerExterns || !(var_.IsExtern())) {
		this_.references.AddNameReference(name, n)
	}
}

// Java method: ReferenceMapBuildingCallback.maybeAddPropReference(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) priv_ReferenceMapBuildingCallback_maybeAddPropReference_4096(name string, n pkg_rhino.Node) {
	_ = name
	_ = n
	if this_.outer_.considerExterns || !(j2gort.SetContainsStr(this_.externProps, name)) {
		this_.references.AddPropReference(name, n)
	}
}

// Java method: ReferenceMapBuildingCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		return ((this_.outer_.considerExterns && t.InGlobalScope()) || !(n.IsFromExterns()))
	} else {
		return true
	}
}

// Java method: ReferenceMapBuildingCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) EnterScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalScope() {
		this_.globalScope = t.GetScope()
		{
			r_46 := this_.references
			r_46.globalScope = this_.globalScope
		}
	}
}

// Java method: ReferenceMapBuildingCallback.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *OptimizeCalls_ReferenceMapBuildingCallback) ExitScope(t *NodeTraversal) {
	_ = t
}
func (this_ *OptimizeCalls_ReferenceMap) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	j2gort.ReleaseOrderedMap(&this_.names)
	j2gort.ReleaseOrderedMap(&this_.props)
	this_.slab_OptimizeCalls_NodeList.Release()
}
