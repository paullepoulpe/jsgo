/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/GuardedCallback.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java polymorphic class interface: com.google.javascript.jscomp.GuardedCallback
type GuardedCallback interface {
	Implements_jscomp_GuardedCallback()
	IsJavaNil_() bool
	GetGuardedCallbackBase() *GuardedCallbackBase
	NodeTraversal_Callback
	VisitGuarded(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
	IsGuarded(any) bool
}

func ToIface_GuardedCallback[T interface {
	GuardedCallback
	comparable
}](v T) GuardedCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_GuardedCallback[T GuardedCallback](v T) GuardedCallback {
	var vi GuardedCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.GuardedCallback
type GuardedCallbackBase struct {
	self_                        GuardedCallback
	compiler                     AbstractCompiler                                  `gson:"compiler"`
	registeredGuards             j2gort.Ext_com_google_common_collect_SetMultimap  `gson:"registeredGuards"`
	guarded                      j2gort.Ext_com_google_common_collect_Multiset     `gson:"guarded"`
	installedGuards              j2gort.Ext_com_google_common_collect_ListMultimap `gson:"installedGuards"`
	contextStack                 j2gort.Deque[*GuardedCallback_Context]            `gson:"contextStack"`
	idHash_                      atomic.Int32
	slab_GuardedCallback_Context j2gort.PooledSlab[GuardedCallback_Context]
	slabCtxCount_                int32
}

func (this_ *GuardedCallbackBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(GuardedCallback)
}
func (this_ *GuardedCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*GuardedCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GuardedCallback", "GuardedCallback"
}
func (this_ *GuardedCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GuardedCallback", "com.google.javascript.jscomp.GuardedCallback", func(v_ any) bool {
		return j2gort.InstanceOf[GuardedCallback](v_)
	})
}
func (b *GuardedCallbackBase) GetGuardedCallbackBase() *GuardedCallbackBase { return b }
func (*GuardedCallbackBase) Implements_jscomp_GuardedCallback()             {}
func (*GuardedCallbackBase) Implements_jscomp_NodeTraversal_Callback()      {}

var static_GuardedCallback_CAN_HAVE_GUARDS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_GuardedCallback_CAN_HAVE_GUARDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_GuardedCallback()
	return static_GuardedCallback_CAN_HAVE_GUARDS
}
func StaticSet_GuardedCallback_CAN_HAVE_GUARDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_GuardedCallback()
	static_GuardedCallback_CAN_HAVE_GUARDS = v
}

var clinitDone_GuardedCallback uint32
var clinitState_GuardedCallback j2gort.ClinitState

func clinit_GuardedCallback() {
	if atomic.LoadUint32(&clinitDone_GuardedCallback) != 0 {
		return
	}
	clinitSlow_GuardedCallback()
}

func clinitSlow_GuardedCallback() {
	if !j2gort.EnterClinit(&clinitState_GuardedCallback, &clinitDone_GuardedCallback) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GuardedCallback, &clinitDone_GuardedCallback)
	static_GuardedCallback_CAN_HAVE_GUARDS = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_AND()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_OR()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_COALESCE()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_HOOK()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_IF()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_BLOCK()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_SCRIPT()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_OPTCHAIN_CALL()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_OPTCHAIN_GETELEM()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_OPTCHAIN_GETPROP())})
}

func (this_ *GuardedCallbackBase) initFields_GuardedCallback() {
	_ = this_
	this_.registeredGuards = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_MultimapBuilder_HashKeys_0().HashSetValues_0().Build_0())
	this_.guarded = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_HashMultiset_Create_0())
	this_.installedGuards = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_MultimapBuilder_HashKeys_0().ArrayListValues_0().Build_0())
	this_.contextStack = j2gort.Deque[*GuardedCallback_Context]{}
}

func InitGuardedCallback(this_ *GuardedCallbackBase, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_GuardedCallback()
	this_.compiler = compiler
}

// Java method: GuardedCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GuardedCallbackBase) ShouldTraverse_3(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = traversal
	_ = n
	_ = parent
	if this_.contextStack.IsEmpty() {
		this_.slabCtxCount_ = 0
	}

	if parent == 0 {
		this_.contextStack.Push(StaticGet_GuardedCallback_Context_EMPTY())
	} else {
		ctx_ := this_.contextStack.Peek().descendWithSlab_(&this_.slab_GuardedCallback_Context, this_.compiler, parent, n)
		if ctx_ != StaticGet_GuardedCallback_Context_EMPTY() {
			this_.slabCtxCount_++
		}
		this_.contextStack.Push(ctx_)

		if ((parent != 0) && j2gort.ImmSetContainsTypedU8_(StaticGet_GuardedCallback_CAN_HAVE_GUARDS(), uint8(parent.GetToken()))) && j2gort.MultimapContainsKeyNode(this_.registeredGuards, parent) {
			{
				cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapRemoveAllNode(this_.registeredGuards, parent))))
				for cur_1.HasNext() {
					rawElem_2 := cur_1.Next()
					resource := rawElem_2
					_ = resource
					this_.guarded.Add_1(resource)
					j2gort.MultimapPutNode(this_.installedGuards, parent, resource)
				}
			}
		}
	}
	return true
}

// Java method: GuardedCallback.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GuardedCallbackBase) Visit_3(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = n
	_ = parent
	if ((parent != 0) && j2gort.ImmSetContainsTypedU8_(StaticGet_GuardedCallback_CAN_HAVE_GUARDS(), uint8(n.GetToken()))) && j2gort.MultimapContainsKeyNode(this_.installedGuards, n) {
		this_.guarded.RemoveAll(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MultimapRemoveAllNode(this_.installedGuards, n)))
	}
	if GuardedCallback_IsAbrupt(n) {
		this_.priv_GuardedCallback_promoteAbruptReturns_3190(parent)
	}
	this_.self_.VisitGuarded(traversal, n, parent)
	poppedCtx_ := this_.contextStack.Pop()
	if this_.contextStack.IsEmpty() {
		this_.slabCtxCount_ = 0
		this_.slab_GuardedCallback_Context.Release()
	} else if poppedCtx_ != StaticGet_GuardedCallback_Context_EMPTY() {
		this_.slabCtxCount_--
		if this_.slabCtxCount_ == 0 {
			this_.slab_GuardedCallback_Context.Reset()
		}
	}

}

// Java method: GuardedCallback.promoteAbruptReturns(com.google.javascript.rhino.Node parent)
func (this_ *GuardedCallbackBase) priv_GuardedCallback_promoteAbruptReturns_3190(parent pkg_rhino.Node) {
	_ = parent
	if parent.IsBlock() {
		parent = parent.GetParent()
	}
	if parent.IsIf() && j2gort.MultimapContainsKeyNode(this_.installedGuards, parent) {
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if grandparent.IsBlock() || grandparent.IsScript() {
			j2gort.MultimapPutAllNode(this_.registeredGuards, grandparent, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MultimapGetNode(this_.installedGuards, parent)))
		}
	}
}

// Java method: GuardedCallback.visitGuarded(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GuardedCallbackBase) VisitGuarded(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = n
	_ = parent
	panic("abstract or native method invoked")
}

// Java method: GuardedCallback.isGuarded(T resource)
func (this_ *GuardedCallbackBase) IsGuarded(resource any) bool {
	_ = resource
	if this_.guarded.Contains(resource) {
		return true
	}
	var context *GuardedCallback_Context = this_.contextStack.Peek()
	_ = context
	if !(context.safe) {
		return false
	}
	for (context != nil) && (context.conditional != 0) {
		j2gort.MultimapPutNode(this_.registeredGuards, context.conditional, resource)
		context = context.linked
	}
	return true
}

// Java static method: GuardedCallback.isAbrupt(com.google.javascript.rhino.Node n)
func GuardedCallback_IsAbrupt(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsReturn() || n.IsThrow())
}

// Java static method: GuardedCallback.isPropertyTestFunction(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node n)
func GuardedCallback_IsPropertyTestFunction(compiler AbstractCompiler, n pkg_rhino.Node) bool {
	_ = compiler
	_ = n
	if compiler.GetCodingConvention().IsPropertyTestFunction(n) {
		return true
	}
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	if target.IsName() {
		var name string = target.GetString()
		_ = name
		return ((j2gort.CheckNotNullString(name) == "String") || (j2gort.CheckNotNullString(name) == "Boolean"))
	}
	return false
}

// Java class: com.google.javascript.jscomp.GuardedCallback.Context
type GuardedCallback_Context struct {
	linked      *GuardedCallback_Context `gson:"linked"`
	conditional pkg_rhino.Node           `gson:"conditional"`
	idHash_     atomic.Int32
	safe        bool `gson:"safe"`
}

func (this_ *GuardedCallback_Context) BindSelf_(s any) {
	_ = s
}
func (this_ *GuardedCallback_Context) IsJavaNil_() bool { return this_ == nil }
func (*GuardedCallback_Context) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GuardedCallback$Context", "Context"
}
func (this_ *GuardedCallback_Context) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GuardedCallback_Context) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GuardedCallback.Context", "com.google.javascript.jscomp.GuardedCallback$Context", func(v_ any) bool {
		return j2gort.InstanceOf[*GuardedCallback_Context](v_)
	})
}

var static_GuardedCallback_Context_EMPTY *GuardedCallback_Context

func StaticGet_GuardedCallback_Context_EMPTY() *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	return static_GuardedCallback_Context_EMPTY
}
func StaticSet_GuardedCallback_Context_EMPTY(v *GuardedCallback_Context) {
	clinit_GuardedCallback_Context()
	static_GuardedCallback_Context_EMPTY = v
}

var clinitDone_GuardedCallback_Context uint32
var clinitState_GuardedCallback_Context j2gort.ClinitState

func clinit_GuardedCallback_Context() {
	if atomic.LoadUint32(&clinitDone_GuardedCallback_Context) != 0 {
		return
	}
	clinitSlow_GuardedCallback_Context()
}

func clinitSlow_GuardedCallback_Context() {
	if !j2gort.EnterClinit(&clinitState_GuardedCallback_Context, &clinitDone_GuardedCallback_Context) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GuardedCallback_Context, &clinitDone_GuardedCallback_Context)
	static_GuardedCallback_Context_EMPTY = NewGuardedCallback_Context(pkg_rhino.Node(0), false, nil)
}

func (this_ *GuardedCallback_Context) initFields_GuardedCallback_Context() {
	_ = this_
}

// Java constructor: Context(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node conditional, boolean safe, com.google.javascript.jscomp.GuardedCallback.@org.jspecify.annotations.Nullable Context linked)
func NewGuardedCallback_Context(conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	this_ := &GuardedCallback_Context{}
	InitGuardedCallback_Context(this_, conditional, safe, linked)
	return this_
}
func allocPooled_GuardedCallback_Context(slab_ *j2gort.PooledSlab[GuardedCallback_Context], conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	this_ := slab_.Alloc()
	InitGuardedCallback_Context(this_, conditional, safe, linked)
	return this_
}

func InitGuardedCallback_Context(this_ *GuardedCallback_Context, conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) {
	_ = conditional
	_ = safe
	_ = linked
	this_.initFields_GuardedCallback_Context()
	this_.conditional = conditional
	this_.safe = safe
	this_.linked = linked
}

// Java method: Context.link(com.google.javascript.rhino.Node newConditional, boolean newSafe)
func (this_ *GuardedCallback_Context) Link(newConditional pkg_rhino.Node, newSafe bool) *GuardedCallback_Context {
	_ = newConditional
	_ = newSafe
	var iifeArg_1 *GuardedCallback_Context
iifeSwitch_1:
	switch {
	default:

		if this_.conditional != 0 {
			iifeArg_1 = this_
			break iifeSwitch_1
		}
		iifeArg_1 = nil
		break iifeSwitch_1
	}
	return NewGuardedCallback_Context(newConditional, newSafe, iifeArg_1)
}
func (this_ *GuardedCallback_Context) linkWithSlab_(slab_ *j2gort.PooledSlab[GuardedCallback_Context], newConditional pkg_rhino.Node, newSafe bool) *GuardedCallback_Context {
	_ = newConditional
	_ = newSafe
	var iifeArg_1 *GuardedCallback_Context
iifeSwitch_1:
	switch {
	default:
		if this_.conditional != 0 {
			iifeArg_1 = this_
			break iifeSwitch_1
		}
		iifeArg_1 = nil
		break iifeSwitch_1
	}
	return allocPooled_GuardedCallback_Context(slab_, newConditional, newSafe, iifeArg_1)
}

// Java method: Context.propagate(boolean newSafe)
func (this_ *GuardedCallback_Context) Propagate(newSafe bool) *GuardedCallback_Context {
	_ = newSafe
	if newSafe == this_.safe {
		return this_
	} else {
		return NewGuardedCallback_Context(this_.conditional, newSafe, this_.linked)
	}
}
func (this_ *GuardedCallback_Context) propagateWithSlab_(slab_ *j2gort.PooledSlab[GuardedCallback_Context], newSafe bool) *GuardedCallback_Context {
	_ = newSafe
	if newSafe == this_.safe {
		return this_
	} else {
		return allocPooled_GuardedCallback_Context(slab_, this_.conditional, newSafe, this_.linked)
	}
}

// Java method: Context.descend(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node child)
func (this_ *GuardedCallback_Context) Descend(compiler AbstractCompiler, parent pkg_rhino.Node, child pkg_rhino.Node) *GuardedCallback_Context {
	_ = compiler
	_ = parent
	_ = child
	var first bool = (child == parent.GetFirstChild())
	_ = first
	switch parent.TokenSwitchOrd_() {
	case int32(101) /* Token.CAST */ :
		return this_
	case int32(50) /* Token.COMMA */ :
		if child == parent.GetLastChild() {
			return this_
		} else {
			return this_.Propagate(true)
		}
	case int32(69) /* Token.AND */ :
		if first {
			return this_.Link(parent, true)
		} else {
			return this_
		}
	case int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		if first {
			return this_.Link(parent, false)
		} else {
			return this_
		}
	case int32(67) /* Token.HOOK */ :
		if first {
			return this_.Link(parent, true)
		} else {
			return this_
		}
	case int32(74) /* Token.IF */ :
		if first {
			return this_.Link(parent, true)
		} else {
			return StaticGet_GuardedCallback_Context_EMPTY()
		}
	case int32(45) /* Token.INSTANCEOF */, int32(51) /* Token.ASSIGN */ :
		return this_.Propagate(first)
	case int32(25) /* Token.TYPEOF */, int32(19) /* Token.NOT */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */ :
		return this_.Propagate(true)
	case int32(28) /* Token.CALL */ :
		return this_.Propagate((!(first) && GuardedCallback_IsPropertyTestFunction(compiler, parent)))
	case int32(89) /* Token.ROOT */ :
		return StaticGet_GuardedCallback_Context_EMPTY()
	case int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		if first {
			return this_.Link(parent, parent.IsOptionalChainStart())
		} else {
			return this_.Propagate(false)
		}
	default:
		if NodeUtil_IsStatement(parent) {
			return StaticGet_GuardedCallback_Context_EMPTY()
		} else {
			return this_.Propagate(false)
		}
	}
	panic("unreachable")
}
func (this_ *GuardedCallback_Context) descendWithSlab_(slab_ *j2gort.PooledSlab[GuardedCallback_Context], compiler AbstractCompiler, parent pkg_rhino.Node, child pkg_rhino.Node) *GuardedCallback_Context {
	_ = compiler
	_ = parent
	_ = child
	var first bool = (child == parent.GetFirstChild())
	_ = first
	switch parent.TokenSwitchOrd_() {
	case int32(101):
		return this_
	case int32(50):
		if child == parent.GetLastChild() {
			return this_
		} else {
			return this_.propagateWithSlab_(slab_, true)
		}
	case int32(69):
		if first {
			return this_.linkWithSlab_(slab_, parent, true)
		} else {
			return this_
		}
	case int32(68), int32(70):
		if first {
			return this_.linkWithSlab_(slab_, parent, false)
		} else {
			return this_
		}
	case int32(67):
		if first {
			return this_.linkWithSlab_(slab_, parent, true)
		} else {
			return this_
		}
	case int32(74):
		if first {
			return this_.linkWithSlab_(slab_, parent, true)
		} else {
			return StaticGet_GuardedCallback_Context_EMPTY()
		}
	case int32(45), int32(51):
		return this_.propagateWithSlab_(slab_, first)
	case int32(25), int32(19), int32(4), int32(5), int32(40), int32(41):
		return this_.propagateWithSlab_(slab_, true)
	case int32(28):
		return this_.propagateWithSlab_(slab_, (!(first) && GuardedCallback_IsPropertyTestFunction(compiler, parent)))
	case int32(89):
		return StaticGet_GuardedCallback_Context_EMPTY()
	case int32(31), int32(30), int32(29):
		if first {
			return this_.linkWithSlab_(slab_, parent, parent.IsOptionalChainStart())
		} else {
			return this_.propagateWithSlab_(slab_, false)
		}
	default:
		if NodeUtil_IsStatement(parent) {
			return StaticGet_GuardedCallback_Context_EMPTY()
		} else {
			return this_.propagateWithSlab_(slab_, false)
		}
	}
	panic("unreachable")
}
func (this_ *GuardedCallbackBase) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slabCtxCount_ = 0
	this_.slab_GuardedCallback_Context.Release()
}
