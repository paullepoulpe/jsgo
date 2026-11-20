/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Scope.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Scope
type Scope struct {
	AbstractScopeBase
	parent   *Scope `gson:"parent"`
	varSlab_ *j2gort.PooledSlab[Var]
	var0     *Var                             `gson:"var0"`
	var1     *Var                             `gson:"var1"`
	var2     *Var                             `gson:"var2"`
	var3     *Var                             `gson:"var3"`
	var4     *Var                             `gson:"var4"`
	var5     *Var                             `gson:"var5"`
	vars     *j2gort.OrderedMap[string, *Var] `gson:"vars"`
	depth    int32                            `gson:"depth"`
}

func (this_ *Scope) BindSelf_(s any) {
	_ = s
	this_.AbstractScopeBase.BindSelf_(s)
}
func (this_ *Scope) IsJavaNil_() bool           { return this_ == nil }
func (*Scope) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Scope", "Scope" }
func (this_ *Scope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractScopeBase.IdentityHash_()
}
func (*Scope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Scope", "com.google.javascript.jscomp.Scope", func(v_ any) bool {
		return j2gort.InstanceOf[*Scope](v_)
	})
}
func (*Scope) Implements_jscomp_AbstractScope() {}
func (*Scope) Implements_rhino_StaticScope()    {}

func (this_ *Scope) initFields_Scope() {
	_ = this_
}

// Java constructor: Scope(com.google.javascript.jscomp.Scope parent, com.google.javascript.rhino.Node rootNode)
func NewScope_2(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, rootNode)
	return this_
}
func allocPooled_Scope_2(slab_ *j2gort.PooledSlab[Scope], parent *Scope, rootNode pkg_rhino.Node) *Scope {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, rootNode)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, rootNode pkg_rhino.Node) {
	_ = parent
	_ = rootNode
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_Scope()
	this_.CheckChildScope(ToIface_AbstractScope(parent))
	this_.parent = parent
	this_.varSlab_ = parent.varSlab_
	this_.depth = (parent.GetDepth() + int32(1))
}

// Java constructor: Scope(com.google.javascript.rhino.Node rootNode)
func NewScope_1(rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_1(this_, rootNode)
	return this_
}
func allocPooled_Scope_1(slab_ *j2gort.PooledSlab[Scope], rootNode pkg_rhino.Node) *Scope {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitScope_1(this_, rootNode)
	return this_
}

func InitScope_1(this_ *Scope, rootNode pkg_rhino.Node) {
	_ = rootNode
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_Scope()
	this_.CheckRootScope()
	this_.parent = nil
	this_.depth = int32(0)
}

// Java static method: Scope.createGlobalScope(com.google.javascript.rhino.Node rootNode)
func Scope_CreateGlobalScope(rootNode pkg_rhino.Node) *Scope {
	_ = rootNode
	return NewScope_1(rootNode)
}
func allocPooledVia_Scope_CreateGlobalScope(slab_ *j2gort.PooledSlab[Scope], rootNode pkg_rhino.Node) *Scope {
	_ = rootNode
	return allocPooled_Scope_1(slab_, rootNode)
}

// Java static method: Scope.createChildScope(com.google.javascript.jscomp.Scope parent, com.google.javascript.rhino.Node rootNode)
func Scope_CreateChildScope(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	_ = parent
	_ = rootNode
	return NewScope_2(parent, rootNode)
}
func allocPooledVia_Scope_CreateChildScope(slab_ *j2gort.PooledSlab[Scope], parent *Scope, rootNode pkg_rhino.Node) *Scope {
	_ = parent
	_ = rootNode
	return allocPooled_Scope_2(slab_, parent, rootNode)
}
func allocFastScopeGlobal_(slab_ *j2gort.PooledSlab[Scope], varSlab_ *j2gort.PooledSlab[Var], rootNode pkg_rhino.Node) *Scope {
	if rootNode == 0 || !(rootNode.IsRoot() || rootNode.IsScript() || NodeUtil_CreatesScope(rootNode)) {
		s_ := allocPooledVia_Scope_CreateGlobalScope(slab_, rootNode)
		s_.varSlab_ = varSlab_
		return s_
	}
	this_ := slab_.Alloc()
	this_.self_ = this_
	this_.rootNode = rootNode
	this_.varSlab_ = varSlab_
	return this_
}

func allocFastScopeChild_(slab_ *j2gort.PooledSlab[Scope], varSlab_ *j2gort.PooledSlab[Var], parent *Scope, rootNode pkg_rhino.Node) *Scope {
	if rootNode == 0 || parent == nil || !NodeUtil_CreatesScope(rootNode) || rootNode == parent.rootNode {
		s_ := allocPooledVia_Scope_CreateChildScope(slab_, parent, rootNode)
		s_.varSlab_ = varSlab_
		return s_
	}
	this_ := slab_.Alloc()
	this_.self_ = this_
	this_.rootNode = rootNode
	this_.parent = parent
	this_.varSlab_ = varSlab_
	this_.depth = parent.depth + 1
	return this_
}

// Java method: Scope.untyped()
func (this_ *Scope) Untyped() *Scope {
	return this_
}

// Java method: Scope.getDepth()
func (this_ *Scope) GetDepth() int32 {
	return this_.depth
}

// Java method: Scope.getParent()
func (this_ *Scope) GetParent() any {
	return j2gort.PtrToAny(this_.parent)
}

// Java method: Scope.declare(java.lang.String name, com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.CompilerInput input)
func (this_ *Scope) Declare(name string, nameNode pkg_rhino.Node, input *CompilerInput) *Var {
	_ = name
	_ = nameNode
	_ = input
	if !(!(j2gort.CheckNotNullString(name) == "")) {
		j2gort.CheckArgument(false)
	}
	if !(j2gort.CastPtr[Var](this_.GetOwnSlot(name)) == nil) {
		j2gort.CheckState(false)
	}
	var var_ *Var = allocOpt_Var(this_.varSlab_, name, nameNode, this_, this_.GetVarCount(), input, pkg_rhino.Node(0))
	_ = var_
	this_.DeclareInternal(name, ToIface_AbstractVar(var_))
	return var_
}

// Java method: Scope.declareImplicitGoogNamespaceIfAbsent(java.lang.String name, com.google.javascript.rhino.Node definition)
func (this_ *Scope) DeclareImplicitGoogNamespaceIfAbsent(name string, definition pkg_rhino.Node) *Var {
	_ = name
	_ = definition
	if !(!(j2gort.CheckNotNullString(name) == "")) {
		j2gort.CheckArgument(false)
	}
	if !(this_.IsGlobal()) {
		j2gort.CheckStateFmt(false, "Cannot declare implicit goog namespace in local scope %s", this_)
	}
	var var_ *Var = j2gort.CastPtr[Var](this_.GetOwnSlot(name))
	_ = var_
	if var_ == nil {
		var_ = allocOptVia_Var_CreateImplicitGoogNamespace(this_.varSlab_, name, this_, definition)
		this_.DeclareInternal(name, ToIface_AbstractVar(var_))
	} else if var_.IsImplicitGoogNamespace() {
		var_.AddImplicitGoogNamespaceDefinition(definition)
	}
	return var_
}

// Java method: Scope.makeImplicitVar(com.google.javascript.jscomp.AbstractScope.ImplicitVar var)
func (this_ *Scope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) any {
	_ = var_
	return j2gort.PtrToAny(allocOpt_Var(this_.varSlab_, var_.FieldGet_name(), pkg_rhino.Node(0), this_, int32(-1), nil, pkg_rhino.Node(0)))
}

// Java method: Scope.findOwnDeclaredVar(java.lang.String name)
func (this_ *Scope) priv_Scope_findOwnDeclaredVar_5046(name string) *Var {
	_ = name
	if this_.vars != nil {
		return this_.vars.Get(name)
	}
	if this_.var0 != nil {
		if j2gort.CheckNotNullString(this_.var0.GetName()) == name {
			return this_.var0
		}
		if this_.var1 != nil {
			if j2gort.CheckNotNullString(this_.var1.GetName()) == name {
				return this_.var1
			}
			if this_.var2 != nil {
				if j2gort.CheckNotNullString(this_.var2.GetName()) == name {
					return this_.var2
				}
				if this_.var3 != nil {
					if j2gort.CheckNotNullString(this_.var3.GetName()) == name {
						return this_.var3
					}
					if this_.var4 != nil {
						if j2gort.CheckNotNullString(this_.var4.GetName()) == name {
							return this_.var4
						}
						if (this_.var5 != nil) && (j2gort.CheckNotNullString(this_.var5.GetName()) == name) {
							return this_.var5
						}
					}
				}
			}
		}
	}
	return nil
}

// Java method: Scope.getOwnDeclaredVar(java.lang.String name)
func (this_ *Scope) GetOwnDeclaredVar(name string) any {
	_ = name
	return j2gort.PtrToAny(this_.priv_Scope_findOwnDeclaredVar_5046(name))
}

// Java method: Scope.putOwnDeclaredVar(java.lang.String name, com.google.javascript.jscomp.Var var)
func (this_ *Scope) PutOwnDeclaredVar(name string, arg_1 any) {
	_ = name
	var_ := j2gort.CastPtr[Var](arg_1)
	_ = var_
	if this_.vars != nil {
		this_.vars.Put(name, var_)
		return
	}
	if (this_.var0 == nil) || (j2gort.CheckNotNullString(this_.var0.GetName()) == name) {
		this_.var0 = var_
		return
	}
	if (this_.var1 == nil) || (j2gort.CheckNotNullString(this_.var1.GetName()) == name) {
		this_.var1 = var_
		return
	}
	if (this_.var2 == nil) || (j2gort.CheckNotNullString(this_.var2.GetName()) == name) {
		this_.var2 = var_
		return
	}
	if (this_.var3 == nil) || (j2gort.CheckNotNullString(this_.var3.GetName()) == name) {
		this_.var3 = var_
		return
	}
	if (this_.var4 == nil) || (j2gort.CheckNotNullString(this_.var4.GetName()) == name) {
		this_.var4 = var_
		return
	}
	if (this_.var5 == nil) || (j2gort.CheckNotNullString(this_.var5.GetName()) == name) {
		this_.var5 = var_
		return
	}
	var iifeArg_1 int32
iifeSwitch_1:
	switch {
	default:

		if this_.parent == nil {
			iifeArg_1 = int32(1024)
			break iifeSwitch_1
		}
		iifeArg_1 = int32(16)
		break iifeSwitch_1
	}
	this_.vars = j2gort.AcquireOrderedMapPtrCap[string, *Var](iifeArg_1)
	this_.vars.Put(this_.var0.GetName(), this_.var0)
	this_.vars.Put(this_.var1.GetName(), this_.var1)
	this_.vars.Put(this_.var2.GetName(), this_.var2)
	this_.vars.Put(this_.var3.GetName(), this_.var3)
	this_.vars.Put(this_.var4.GetName(), this_.var4)
	this_.vars.Put(this_.var5.GetName(), this_.var5)
	this_.vars.Put(name, var_)
	this_.var0 = nil
	this_.var1 = nil
	this_.var2 = nil
	this_.var3 = nil
	this_.var4 = nil
	this_.var5 = nil
}

// Java method: Scope.undeclareInteral(com.google.javascript.jscomp.Var var)
func (this_ *Scope) UndeclareInteral(arg_2 any) {
	var_ := j2gort.CastPtr[Var](arg_2)
	_ = var_
	if this_.vars != nil {
		this_.vars.Remove(var_.GetName())
		return
	}
	var name string = var_.GetName()
	_ = name
	if this_.var0 != nil {
		if j2gort.CheckNotNullString(this_.var0.GetName()) == name {
			this_.var0 = this_.var1
			this_.var1 = this_.var2
			this_.var2 = this_.var3
			this_.var3 = this_.var4
			this_.var4 = this_.var5
			this_.var5 = nil
		} else if this_.var1 != nil {
			if j2gort.CheckNotNullString(this_.var1.GetName()) == name {
				this_.var1 = this_.var2
				this_.var2 = this_.var3
				this_.var3 = this_.var4
				this_.var4 = this_.var5
				this_.var5 = nil
			} else if this_.var2 != nil {
				if j2gort.CheckNotNullString(this_.var2.GetName()) == name {
					this_.var2 = this_.var3
					this_.var3 = this_.var4
					this_.var4 = this_.var5
					this_.var5 = nil
				} else if this_.var3 != nil {
					if j2gort.CheckNotNullString(this_.var3.GetName()) == name {
						this_.var3 = this_.var4
						this_.var4 = this_.var5
						this_.var5 = nil
					} else if this_.var4 != nil {
						if j2gort.CheckNotNullString(this_.var4.GetName()) == name {
							this_.var4 = this_.var5
							this_.var5 = nil
						} else if (this_.var5 != nil) && (j2gort.CheckNotNullString(this_.var5.GetName()) == name) {
							this_.var5 = nil
						}
					}
				}
			}
		}
	}
}

// Java method: Scope.clearVarsInternal()
func (this_ *Scope) ClearVarsInternal() {
	this_.var0 = nil
	this_.var1 = nil
	this_.var2 = nil
	this_.var3 = nil
	this_.var4 = nil
	this_.var5 = nil
	if (this_.vars != nil) && !(this_.vars.IsEmpty()) {
		this_.vars.Clear()
	}
}

// Java method: Scope.lookupVar(java.lang.String name)
func (this_ *Scope) LookupVar(name string) *Var {
	_ = name
	var implicit AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(name)
	_ = implicit
	if implicit == 0 {
		{
			var scope *Scope = this_
			_ = scope
			for ; scope != nil; scope = scope.parent {
				var var_ *Var = scope.priv_Scope_findOwnDeclaredVar_5046(name)
				_ = var_
				if var_ != nil {
					return var_
				}
			}
		}
		return nil
	}
	{
		var scope_2 *Scope = this_
		_ = scope_2
		for ; scope_2 != nil; scope_2 = scope_2.parent {
			var var__2 *Var = scope_2.priv_Scope_findOwnDeclaredVar_5046(name)
			_ = var__2
			if var__2 != nil {
				return var__2
			}
			var__2 = j2gort.CastPtr[Var](scope_2.GetOwnImplicitSlot(implicit))
			if var__2 != nil {
				return var__2
			}
		}
	}
	return nil
}

// Java method: Scope.getVar(java.lang.String name)
func (this_ *Scope) GetVar(name string) any {
	_ = name
	return j2gort.PtrToAny(this_.LookupVar(name))
}

// Java method: Scope.getVarIterable()
func (this_ *Scope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	if this_.vars != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.vars.LiveValues())
	}
	if this_.var0 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_EmptyList())
	}
	if this_.var1 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.SingletonListPtr(this_.var0))
	}
	if this_.var2 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1)))
	}
	if this_.var3 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(this_.var2)))
	}
	if this_.var4 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_4(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(this_.var2), j2gort.PtrToAny(this_.var3)))
	}
	if this_.var5 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_5(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(this_.var2), j2gort.PtrToAny(this_.var3), j2gort.PtrToAny(this_.var4)))
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_6(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(this_.var2), j2gort.PtrToAny(this_.var3), j2gort.PtrToAny(this_.var4), j2gort.PtrToAny(this_.var5)))
}

// Java method: Scope.putOwnVarsIfAbsent(java.util.LinkedHashMap<java.lang.String,com.google.javascript.jscomp.Var> dst)
func (this_ *Scope) PutOwnVarsIfAbsent(dst *j2gort.Ext_java_util_LinkedHashMap) {
	_ = dst
	if this_.vars != nil {
		{
			cur_3 := this_.vars.IterValues()
			for cur_3.HasNext() {

				v := cur_3.
					Next()
				_ = v
				dst.PutIfAbsent(j2gort.BoxString(v.GetName()), j2gort.PtrToAny(v))
			}
		}
	} else if this_.var0 != nil {
		dst.PutIfAbsent(j2gort.BoxString(this_.var0.GetName()), j2gort.PtrToAny(this_.var0))
		if this_.var1 != nil {
			dst.PutIfAbsent(j2gort.BoxString(this_.var1.GetName()), j2gort.PtrToAny(this_.var1))
			if this_.var2 != nil {
				dst.PutIfAbsent(j2gort.BoxString(this_.var2.GetName()), j2gort.PtrToAny(this_.var2))
				if this_.var3 != nil {
					dst.PutIfAbsent(j2gort.BoxString(this_.var3.GetName()), j2gort.PtrToAny(this_.var3))
					if this_.var4 != nil {
						dst.PutIfAbsent(j2gort.BoxString(this_.var4.GetName()), j2gort.PtrToAny(this_.var4))
						if this_.var5 != nil {
							dst.PutIfAbsent(j2gort.BoxString(this_.var5.GetName()), j2gort.PtrToAny(this_.var5))
						}
					}
				}
			}
		}
	}
}

// Java method: Scope.getAllSymbols()
func (this_ *Scope) GetAllSymbols() j2gort.Ext_java_lang_Iterable {
	if this_.vars != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.vars.LiveValues())))
	}
	return this_.GetVarIterable()
}

// Java method: Scope.addVarsTo(java.util.ArrayList<com.google.javascript.jscomp.Var> dst)
func (this_ *Scope) AddVarsTo(dst j2gort.Ext_java_util_ArrayList) {
	_ = dst
	if this_.vars != nil {
		{
			cur_5 := this_.vars.IterValues()
			for cur_5.HasNext() {

				v := cur_5.
					Next()
				_ = v
				dst.Add_1(j2gort.PtrToAny(v))
			}
		}
	} else if this_.var0 != nil {
		dst.Add_1(j2gort.PtrToAny(this_.var0))
		if this_.var1 != nil {
			dst.Add_1(j2gort.PtrToAny(this_.var1))
			if this_.var2 != nil {
				dst.Add_1(j2gort.PtrToAny(this_.var2))
				if this_.var3 != nil {
					dst.Add_1(j2gort.PtrToAny(this_.var3))
					if this_.var4 != nil {
						dst.Add_1(j2gort.PtrToAny(this_.var4))
						if this_.var5 != nil {
							dst.Add_1(j2gort.PtrToAny(this_.var5))
						}
					}
				}
			}
		}
	}
}

// Java method: Scope.addVarNamesTo(java.util.LinkedHashSet<java.lang.String> dst)
func (this_ *Scope) AddVarNamesTo(dst *j2gort.Ext_java_util_LinkedHashSet) {
	_ = dst
	if this_.vars != nil {
		{
			cur_7 := this_.vars.IterKeys()
			for cur_7.HasNext() {
				name := cur_7.
					Next()
				_ = name
				j2gort.SetAddStr(dst, name)
			}
		}
	} else if this_.var0 != nil {
		j2gort.SetAddStr(dst, this_.var0.GetName())
		if this_.var1 != nil {
			j2gort.SetAddStr(dst, this_.var1.GetName())
			if this_.var2 != nil {
				j2gort.SetAddStr(dst, this_.var2.GetName())
				if this_.var3 != nil {
					j2gort.SetAddStr(dst, this_.var3.GetName())
					if this_.var4 != nil {
						j2gort.SetAddStr(dst, this_.var4.GetName())
						if this_.var5 != nil {
							j2gort.SetAddStr(dst, this_.var5.GetName())
						}
					}
				}
			}
		}
	}
}

// Java method: Scope.getOwnVarNamesArray()
func (this_ *Scope) GetOwnVarNamesArray() []string {
	if this_.vars != nil {
		var n int32 = this_.vars.Size()
		_ = n
		if n == int32(0) {
			return TypedScopeCreator_EmptyVarNames()
		}
		var out []string = j2gort.NewStringArray(n)
		_ = out
		var i int32 = int32(0)
		_ = i
		{
			cur_8 := this_.vars.IterKeys()
			for cur_8.HasNext() {
				name := cur_8.
					Next()
				_ = name
				{
					a_9 := out
					var i_10 int32
					{

						old_11 := i
						nw_12 := old_11 + 1
						i = nw_12
						i_10 = old_11
					}
					a_9[i_10] = name
				}
			}
		}
		return out
	}
	var count int32 = this_.GetVarCount()
	_ = count
	if count == int32(0) {
		return TypedScopeCreator_EmptyVarNames()
	}
	if count == int32(1) {
		return TypedScopeCreator_SingletonVarNameArray(this_.var0.GetName())
	}
	if count == int32(2) {
		return TypedScopeCreator_PairVarNameArray(this_.var0.GetName(), this_.var1.GetName())
	}
	var out_2 []string = j2gort.NewStringArray(count)
	_ = out_2
	{
		a_13 := out_2
		i_14 := int32(0)
		a_13[i_14] = this_.var0.GetName()
	}
	{
		a_15 := out_2
		i_16 := int32(1)
		a_15[i_16] = this_.var1.GetName()
	}
	{
		a_17 := out_2
		i_18 := int32(2)
		a_17[i_18] = this_.var2.GetName()
	}
	if count > int32(3) {
		{
			a_19 := out_2
			i_20 := int32(3)
			a_19[i_20] = this_.var3.GetName()
		}
		if count > int32(4) {
			{
				a_21 := out_2
				i_22 := int32(4)
				a_21[i_22] = this_.var4.GetName()
			}
			if count > int32(5) {
				{
					a_23 := out_2
					i_24 := int32(5)
					a_23[i_24] = this_.var5.GetName()
				}
			}
		}
	}
	return out_2
}

// Java method: Scope.mergeVarNamesArray(java.lang.String[] existing)
func (this_ *Scope) MergeVarNamesArray(existing []string) []string {
	_ = existing
	var merged *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(existing)) + this_.GetVarCount()))
	_ = merged
	var iifeArg_2 []any
iifeSwitch_2:
	switch {
	default:

		arr_25 := existing
		if arr_25 == nil {
			iifeArg_2 = nil
			break iifeSwitch_2
		}
		out_26 := make([]any, len(arr_25))
		for i_, el_27 := range arr_25 {
			out_26[i_] = j2gort.BoxString(el_27)
		}
		iifeArg_2 = out_26
		break iifeSwitch_2
	}
	j2gort.Ext_java_util_Collections_AddAll(merged, iifeArg_2)
	this_.AddVarNamesTo(merged)
	return j2gort.ColToNewStringArray(merged)
}

// Java method: Scope.getVarCount()
func (this_ *Scope) GetVarCount() int32 {
	if this_.vars != nil {
		return this_.vars.Size()
	}
	if this_.var0 == nil {
		return int32(0)
	}
	if this_.var1 == nil {
		return int32(1)
	}
	if this_.var2 == nil {
		return int32(2)
	}
	if this_.var3 == nil {
		return int32(3)
	}
	if this_.var4 == nil {
		return int32(4)
	}
	if this_.var5 == nil {
		return int32(5)
	} else {
		return int32(6)
	}
}

type scopeVarCursor_[V any] struct {
	vars   [6]V
	idx    int
	n      int
	useMap bool
	mapCur j2gort.OrderedMapValCursor[string, V]
}

func (c *scopeVarCursor_[V]) HasNext() bool {
	if c.useMap {
		return c.mapCur.HasNext()
	}
	return c.idx < c.n
}

func (c *scopeVarCursor_[V]) Next() V {
	if c.useMap {
		return c.mapCur.Next()
	}
	v := c.vars[c.idx]
	c.idx++
	return v
}
func (this_ *Scope) iterVars_() scopeVarCursor_[*Var] {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if this_.vars != nil {
		return scopeVarCursor_[*Var]{useMap: true, mapCur: this_.vars.IterValues()}
	}
	if this_.var0 == nil {
		return scopeVarCursor_[*Var]{}
	}
	var c scopeVarCursor_[*Var]
	c.vars[0] = this_.var0
	if this_.var1 == nil {
		c.n = 1
		return c
	}
	c.vars[1] = this_.var1
	if this_.var2 == nil {
		c.n = 2
		return c
	}
	c.vars[2] = this_.var2
	if this_.var3 == nil {
		c.n = 3
		return c
	}
	c.vars[3] = this_.var3
	if this_.var4 == nil {
		c.n = 4
		return c
	}
	c.vars[4] = this_.var4
	if this_.var5 == nil {
		c.n = 5
		return c
	}
	c.vars[5] = this_.var5
	c.n = 6
	return c
}
