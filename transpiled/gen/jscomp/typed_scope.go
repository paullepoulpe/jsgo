/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypedScope.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.TypedScope
type TypedScope struct {
	AbstractScopeBase
	parent   *TypedScope          `gson:"parent"`
	var0     *TypedVar            `gson:"var0"`
	var1     *TypedVar            `gson:"var1"`
	overflow *TypedScope_Overflow `gson:"overflow"`
	depth    int32                `gson:"depth"`
	isBottom bool                 `gson:"isBottom"`
	tsArena_ *typedScopeArena_
}

func (this_ *TypedScope) BindSelf_(s any) {
	_ = s
	this_.AbstractScopeBase.BindSelf_(s)
}
func (this_ *TypedScope) IsJavaNil_() bool { return this_ == nil }
func (*TypedScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScope", "TypedScope"
}
func (this_ *TypedScope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractScopeBase.IdentityHash_()
}
func (*TypedScope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScope", "com.google.javascript.jscomp.TypedScope", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScope](v_)
	})
}
func (*TypedScope) Implements_jscomp_AbstractScope()   {}
func (*TypedScope) Implements_rhino_StaticTypedScope() {}
func (*TypedScope) Implements_rhino_StaticScope()      {}

func (this_ *TypedScope) initFields_TypedScope() {
	_ = this_
}

// Java constructor: TypedScope(com.google.javascript.jscomp.TypedScope parent, com.google.javascript.rhino.Node rootNode)
func NewTypedScope_2_TypedScope_Node(parent *TypedScope, rootNode pkg_rhino.Node) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_2_TypedScope_Node(this_, parent, rootNode)
	return this_
}

func InitTypedScope_2_TypedScope_Node(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node) {
	_ = parent
	_ = rootNode
	if this_.tsArena_ == nil && parent != nil {
		this_.tsArena_ = parent.tsArena_
	}

	InitTypedScope_5(this_, parent, rootNode, nil, nil, nil)
}

// Java constructor: TypedScope(com.google.javascript.jscomp.TypedScope parent, com.google.javascript.rhino.Node rootNode, java.util.Set<java.lang.String> reservedNames, com.google.javascript.jscomp.modules.@org.jspecify.annotations.Nullable Module module)
func NewTypedScope_4_TypedScope_Node_Ext_java_util_Set_Module(parent *TypedScope, rootNode pkg_rhino.Node, reservedNames j2gort.Ext_java_util_Set, module *Module) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_4_TypedScope_Node_Ext_java_util_Set_Module(this_, parent, rootNode, reservedNames, module)
	return this_
}

func InitTypedScope_4_TypedScope_Node_Ext_java_util_Set_Module(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node, reservedNames j2gort.Ext_java_util_Set, module *Module) {
	_ = parent
	_ = rootNode
	_ = reservedNames
	_ = module
	if this_.tsArena_ == nil && parent != nil {
		this_.tsArena_ = parent.tsArena_
	}

	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_TypedScope()
	this_.CheckChildScope(ToIface_AbstractScope(parent))
	this_.parent = parent
	this_.depth = (parent.depth + int32(1))
	this_.isBottom = false
	var hasReserved bool = !(reservedNames.IsEmpty())
	_ = hasReserved
	if hasReserved || (module != nil) {
		var val_1 TypedScope_Overflow
		InitTypedScope_Overflow(&val_1)
		var o *TypedScope_Overflow = &val_1
		_ = o
		if hasReserved {
			o.InitReservedSet(reservedNames)
		}
		{
			r_2 := o
			r_2.module = module
		}
		this_.overflow = o
	}
}

// Java constructor: TypedScope(com.google.javascript.jscomp.TypedScope parent, com.google.javascript.rhino.Node rootNode, java.lang.String @org.jspecify.annotations.Nullable [] reservedNamesOwned, com.google.javascript.jscomp.modules.@org.jspecify.annotations.Nullable Module module)
func NewTypedScope_4_TypedScope_Node_slice_string_Module(parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_4_TypedScope_Node_slice_string_Module(this_, parent, rootNode, reservedNamesOwned, module)
	return this_
}

func InitTypedScope_4_TypedScope_Node_slice_string_Module(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module) {
	_ = parent
	_ = rootNode
	_ = reservedNamesOwned
	_ = module
	if this_.tsArena_ == nil && parent != nil {
		this_.tsArena_ = parent.tsArena_
	}

	InitTypedScope_5(this_, parent, rootNode, reservedNamesOwned, module, nil)
}

// Java constructor: TypedScope(com.google.javascript.jscomp.TypedScope parent, com.google.javascript.rhino.Node rootNode, java.lang.String @org.jspecify.annotations.Nullable [] reservedNamesOwned, com.google.javascript.jscomp.modules.@org.jspecify.annotations.Nullable Module module, com.google.javascript.jscomp.TypedScope.@org.jspecify.annotations.Nullable Overflow reusableOverflow)
func NewTypedScope_5(parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module, reusableOverflow *TypedScope_Overflow) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_5(this_, parent, rootNode, reservedNamesOwned, module, reusableOverflow)
	return this_
}
func allocPooled_TypedScope_5(slab_ *j2gort.PooledSlab[TypedScope], parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module, reusableOverflow *TypedScope_Overflow) *TypedScope {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTypedScope_5(this_, parent, rootNode, reservedNamesOwned, module, reusableOverflow)
	return this_
}

func InitTypedScope_5(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module, reusableOverflow *TypedScope_Overflow) {
	_ = parent
	_ = rootNode
	_ = reservedNamesOwned
	_ = module
	_ = reusableOverflow
	if this_.tsArena_ == nil && parent != nil {
		this_.tsArena_ = parent.tsArena_
	}

	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_TypedScope()
	this_.CheckChildScope(ToIface_AbstractScope(parent))
	this_.parent = parent
	this_.depth = (parent.depth + int32(1))
	this_.isBottom = false
	var hasReserved bool = ((reservedNamesOwned != nil) && (int32(len(reservedNamesOwned)) != int32(0)))
	_ = hasReserved
	if hasReserved || (module != nil) {
		var o *TypedScope_Overflow
		_ = o
		if reusableOverflow != nil {
			o = reusableOverflow
		} else {
			o = NewTypedScope_Overflow()
		}
		if hasReserved {
			o.InitReservedArray(reservedNamesOwned)
		}
		{
			r_3 := o
			r_3.module = module
		}
		this_.overflow = o
	}
}

// Java constructor: TypedScope(com.google.javascript.rhino.Node rootNode, boolean isBottom)
func NewTypedScope_2_Node_bool(rootNode pkg_rhino.Node, isBottom bool) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_2_Node_bool(this_, rootNode, isBottom)
	return this_
}
func allocPooled_TypedScope_2_Node_bool(slab_ *j2gort.PooledSlab[TypedScope], rootNode pkg_rhino.Node, isBottom bool) *TypedScope {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTypedScope_2_Node_bool(this_, rootNode, isBottom)
	return this_
}

func InitTypedScope_2_Node_bool(this_ *TypedScope, rootNode pkg_rhino.Node, isBottom bool) {
	_ = rootNode
	_ = isBottom
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_TypedScope()
	this_.CheckRootScope()
	this_.parent = nil
	this_.depth = int32(0)
	this_.isBottom = isBottom
}

// Java method: TypedScope.ensureOverflow()
func (this_ *TypedScope) priv_TypedScope_ensureOverflow_5714() *TypedScope_Overflow {
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if o == nil {
		o = NewTypedScope_Overflow()
		this_.overflow = o
	}
	return o
}

// Java method: TypedScope.isUsingOverflow(com.google.javascript.jscomp.TypedScope.Overflow candidate)
func (this_ *TypedScope) IsUsingOverflow(candidate *TypedScope_Overflow) bool {
	_ = candidate
	return (this_.overflow == candidate)
}

// Java static method: TypedScope.createGlobalScope(com.google.javascript.rhino.Node rootNode)
func TypedScope_CreateGlobalScope(rootNode pkg_rhino.Node) *TypedScope {
	_ = rootNode
	return NewTypedScope_2_Node_bool(rootNode, false)
}
func allocPooledVia_TypedScope_CreateGlobalScope(slab_ *j2gort.PooledSlab[TypedScope], rootNode pkg_rhino.Node) *TypedScope {
	_ = rootNode
	return allocPooled_TypedScope_2_Node_bool(slab_, rootNode, false)
}

// Java static method: TypedScope.createLatticeBottom(com.google.javascript.rhino.Node rootNode)
func TypedScope_CreateLatticeBottom(rootNode pkg_rhino.Node) *TypedScope {
	_ = rootNode
	return NewTypedScope_2_Node_bool(rootNode, true)
}
func allocPooledVia_TypedScope_CreateLatticeBottom(slab_ *j2gort.PooledSlab[TypedScope], rootNode pkg_rhino.Node) *TypedScope {
	_ = rootNode
	return allocPooled_TypedScope_2_Node_bool(slab_, rootNode, true)
}

// Java method: TypedScope.typed()
func (this_ *TypedScope) Typed() *TypedScope {
	return this_
}

// Java method: TypedScope.validateCompletelyBuilt()
func (this_ *TypedScope) ValidateCompletelyBuilt() {
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if o != nil {
		if o.HasReservedNames() {
			j2gort.CheckStateFmt(false, "Expected %s to have no reserved names, found: %s. This probably indicates a bug in TypedScopeCreator where it is failing to declare a variable.", this_, j2gort.BoxString(o.ReservedToString()))
		}
		o.ClearReserved()
		if o.IsEmpty() {
			this_.overflow = nil
		}
	}
}

// Java method: TypedScope.isBottom()
func (this_ *TypedScope) IsBottom() bool {
	return this_.isBottom
}

// Java method: TypedScope.getModule()
func (this_ *TypedScope) GetModule() *Module {
	if this_.overflow != nil {
		return this_.overflow.module
	} else {
		return nil
	}
}

// Java method: TypedScope.getDepth()
func (this_ *TypedScope) GetDepth() int32 {
	return this_.depth
}

// Java method: TypedScope.getParent()
func (this_ *TypedScope) GetParent() any {
	return j2gort.PtrToAny(this_.parent)
}

// Java method: TypedScope.getTypeOfThis()
func (this_ *TypedScope) GetTypeOfThis() pkg_rhino.JSType {
	var root pkg_rhino.Node = this_.GetRootNode_0()
	_ = root
	if this_.IsGlobal() {
		return pkg_rhino.Upcast_JSType(pkg_rhino.ObjectType_Cast(root.GetJSType()))
	} else if NodeUtil_IsNonArrowFunction(root) {
		var nodeType pkg_rhino.JSType = root.GetJSType()
		_ = nodeType
		if (nodeType != nil && !nodeType.IsJavaNil_()) && nodeType.GetJSTypeBase().IsFunctionType() {
			return nodeType.ToMaybeFunctionType().GetTypeOfThis()
		} else {
			return nil
		}
	} else if this_.IsStaticBlockScope() {
		return j2gort.CastPtr[TypedScope](this_.GetParent()).GetRootNode_0().GetJSType()
	} else if this_.IsMemberFieldDefScope() || this_.IsComputedFieldDefRhsScope() {
		var classType pkg_rhino.JSType = j2gort.CastPtr[TypedScope](this_.GetParent()).GetRootNode_0().GetJSType()
		_ = classType
		if root.IsStaticMember() {
			return classType
		} else {
			return pkg_rhino.Upcast_JSType(classType.GetJSTypeBase().AssertFunctionType().GetInstanceType())
		}
	} else {
		return j2gort.CastPtr[TypedScope](this_.GetParent()).GetTypeOfThis()
	}
}

// Java method: TypedScope.declare(java.lang.String name, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.jscomp.CompilerInput input, boolean inferred)
func (this_ *TypedScope) Declare(name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	_ = name
	_ = nameNode
	_ = type_
	_ = input
	_ = inferred
	if !((name != j2gort.NullString) && !(j2gort.CheckNotNullString(name) == "")) {
		j2gort.CheckState(false)
	}
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.reservedMask != int32(0)) {
		if o.RemoveReservedAndCheckEmpty(name) && o.IsEmpty() {
			this_.overflow = nil
		}
	}
	var var_ *TypedVar
	if this_.tsArena_ != nil {
		var_ = allocPooled_TypedVar(&this_.tsArena_.slab_TypedVar, inferred, name, nameNode, type_, this_,
			this_.
				GetVarCount(), input)
	} else {
		var_ = NewTypedVar(inferred, name, nameNode, type_, this_, this_.GetVarCount(), input)
	}

	_ = var_
	this_.DeclareInternal(name, ToIface_AbstractVar(var_))
	var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dot
	if dot >= int32(0) {
		var qo *TypedScope_Overflow = this_.priv_TypedScope_ensureOverflow_5714()
		_ = qo
		if qo.qualifiedRootNames == nil {
			{
				r_4 := qo
				r_4.qualifiedRootNames = j2gort.NewPlainSetPtrCap[string](int32(4))
			}
		}
		qo.qualifiedRootNames.AddVoid(j2gort.JStrSub(name, int32(0), dot))
		TypedScope_AddQualifiedVarHash(qo, pkg_rhino.Node_HashQualifiedNameString(name))
	}
	return var_
}

// Java static method: TypedScope.addQualifiedVarHash(com.google.javascript.jscomp.TypedScope.Overflow o, int h)
func TypedScope_AddQualifiedVarHash(o *TypedScope_Overflow, h int32) {
	_ = o
	_ = h
	if h == int32(0) {
		return
	}
	var table []int32 = o.qualifiedVarHashes
	_ = table
	if table == nil {
		table = make([]int32, int32(8))
		{
			r_5 := o
			r_5.qualifiedVarHashes = table
		}
	} else if (o.qualifiedVarHashCount * int32(4)) >= (int32(len(table)) * int32(3)) {
		var next []int32 = make([]int32, (int32(len(table)) * int32(2)))
		_ = next
		var nextMask int32 = (int32(len(next)) - int32(1))
		_ = nextMask
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = int32(len(table))
			_ = len_
			for ; i < len_; i = i + 1 {
				var existing int32 = table[i]
				_ = existing
				if existing != int32(0) {
					var s int32 = (existing & nextMask)
					_ = s
					for next[s] != int32(0) {
						s = ((s + int32(1)) & nextMask)
					}
					{
						a_6 := next
						i_7 := s
						a_6[i_7] = existing
					}
				}
			}
		}
		table = next
		{
			r_8 := o
			r_8.qualifiedVarHashes = table
		}
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing_2 int32 = table[slot]
		_ = existing_2
		if existing_2 == int32(0) {
			{
				a_9 := table
				i_10 := slot
				a_9[i_10] = h
			}
			r_11 := o
			r_11.qualifiedVarHashCount = r_11.qualifiedVarHashCount + 1
			return
		}
		if existing_2 == h {
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: TypedScope.containsQualifiedVarHash(int h)
func (this_ *TypedScope) ContainsQualifiedVarHash(h int32) bool {
	_ = h
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o == nil) || (h == int32(0)) {
		return false
	}
	var table []int32 = o.qualifiedVarHashes
	_ = table
	if table == nil {
		return false
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing int32 = table[slot]
		_ = existing
		if existing == int32(0) {
			return false
		}
		if existing == h {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: TypedScope.mayHaveQualifiedVarWithHash(int h)
func (this_ *TypedScope) MayHaveQualifiedVarWithHash(h int32) bool {
	_ = h
	{
		var s *TypedScope = this_
		_ = s
		for ; s != nil; s = s.parent {
			if s.ContainsQualifiedVarHash(h) {
				return true
			}
		}
	}
	return false
}

// Java method: TypedScope.mayHaveQualifiedVarForNode(com.google.javascript.rhino.Node n)
func (this_ *TypedScope) MayHaveQualifiedVarForNode(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsGetProp()) {
		return n.IsQualifiedName()
	}
	var h int32 = n.HashQualifiedName()
	_ = h
	return ((h != int32(0)) && this_.MayHaveQualifiedVarWithHash(h))
}

// Java method: TypedScope.hasQualifiedOwnVars()
func (this_ *TypedScope) HasQualifiedOwnVars() bool {
	return ((this_.overflow != nil) && (this_.overflow.qualifiedRootNames != nil))
}

// Java method: TypedScope.hasQualifiedOwnVarWithRoot(java.lang.String rootName)
func (this_ *TypedScope) HasQualifiedOwnVarWithRoot(rootName string) bool {
	_ = rootName
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	return (((o != nil) && (o.qualifiedRootNames != nil)) && o.qualifiedRootNames.Contains(rootName))
}

// Java method: TypedScope.mayHaveQualifiedVarForRoot(java.lang.String rootName)
func (this_ *TypedScope) MayHaveQualifiedVarForRoot(rootName string) bool {
	_ = rootName
	{
		var s *TypedScope = this_
		_ = s
		for ; s != nil; s = s.parent {
			if s.HasQualifiedOwnVarWithRoot(rootName) {
				return true
			}
		}
	}
	return false
}

// Java method: TypedScope.findOwnDeclaredVar(java.lang.String name)
func (this_ *TypedScope) priv_TypedScope_findOwnDeclaredVar_5726(name string) *TypedVar {
	_ = name
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		return o.vars.Get(name)
	}
	if this_.var0 != nil {
		if j2gort.CheckNotNullString(this_.var0.GetName()) == name {
			return this_.var0
		}
		if this_.var1 != nil {
			if j2gort.CheckNotNullString(this_.var1.GetName()) == name {
				return this_.var1
			}
			if (o != nil) && (o.var2 != nil) {
				if j2gort.CheckNotNullString(o.var2.GetName()) == name {
					return o.var2
				}
				if o.var3 != nil {
					if j2gort.CheckNotNullString(o.var3.GetName()) == name {
						return o.var3
					}
					if o.var4 != nil {
						if j2gort.CheckNotNullString(o.var4.GetName()) == name {
							return o.var4
						}
						if (o.var5 != nil) && (j2gort.CheckNotNullString(o.var5.GetName()) == name) {
							return o.var5
						}
					}
				}
			}
		}
	}
	return nil
}

// Java method: TypedScope.getOwnDeclaredVar(java.lang.String name)
func (this_ *TypedScope) GetOwnDeclaredVar(name string) any {
	_ = name
	return j2gort.PtrToAny(this_.priv_TypedScope_findOwnDeclaredVar_5726(name))
}

// Java method: TypedScope.putOwnDeclaredVar(java.lang.String name, com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScope) PutOwnDeclaredVar(name string, arg_12 any) {
	_ = name
	var_ := j2gort.CastPtr[TypedVar](arg_12)
	_ = var_
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		o.vars.Put(name, var_)
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
	if o == nil {
		o = this_.priv_TypedScope_ensureOverflow_5714()
	}
	if (o.var2 == nil) || (j2gort.CheckNotNullString(o.var2.GetName()) == name) {
		{
			r_13 := o
			r_13.var2 = var_
		}
		return
	}
	if (o.var3 == nil) || (j2gort.CheckNotNullString(o.var3.GetName()) == name) {
		{
			r_14 := o
			r_14.var3 = var_
		}
		return
	}
	if (o.var4 == nil) || (j2gort.CheckNotNullString(o.var4.GetName()) == name) {
		{
			r_15 := o
			r_15.var4 = var_
		}
		return
	}
	if (o.var5 == nil) || (j2gort.CheckNotNullString(o.var5.GetName()) == name) {
		{
			r_16 := o
			r_16.var5 = var_
		}
		return
	}
	{
		r_17 := o
		r_17.vars = j2gort.NewOrderedMapPtrCap[string, *TypedVar](func() int32 {
			if this_.parent == nil {
				return int32(1024)
			}
			return int32(16)
		}())
	}
	o.vars.Put(this_.var0.GetName(), this_.var0)
	o.vars.Put(this_.var1.GetName(), this_.var1)
	o.vars.Put(o.var2.GetName(), o.var2)
	o.vars.Put(o.var3.GetName(), o.var3)
	o.vars.Put(o.var4.GetName(), o.var4)
	o.vars.Put(o.var5.GetName(), o.var5)
	o.vars.Put(name, var_)
	this_.var0 = nil
	this_.var1 = nil
	{
		r_18 := o
		r_18.var2 = nil
	}
	{
		r_19 := o
		r_19.var3 = nil
	}
	{
		r_20 := o
		r_20.var4 = nil
	}
	{
		r_21 := o
		r_21.var5 = nil
	}
}

// Java method: TypedScope.undeclareInteral(com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScope) UndeclareInteral(arg_22 any) {
	var_ := j2gort.CastPtr[TypedVar](arg_22)
	_ = var_
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		o.vars.Remove(var_.GetName())
		return
	}
	var name string = var_.GetName()
	_ = name
	if this_.var0 != nil {
		if j2gort.CheckNotNullString(this_.var0.GetName()) == name {
			this_.var0 = this_.var1
			if o != nil {
				this_.var1 = o.var2
				{
					r_23 := o
					r_23.var2 = o.var3
				}
				{
					r_24 := o
					r_24.var3 = o.var4
				}
				{
					r_25 := o
					r_25.var4 = o.var5
				}
				{
					r_26 := o
					r_26.var5 = nil
				}
				if o.IsEmpty() {
					this_.overflow = nil
				}
			} else {
				this_.var1 = nil
			}
		} else if this_.var1 != nil {
			if j2gort.CheckNotNullString(this_.var1.GetName()) == name {
				if o != nil {
					this_.var1 = o.var2
					{
						r_27 := o
						r_27.var2 = o.var3
					}
					{
						r_28 := o
						r_28.var3 = o.var4
					}
					{
						r_29 := o
						r_29.var4 = o.var5
					}
					{
						r_30 := o
						r_30.var5 = nil
					}
					if o.IsEmpty() {
						this_.overflow = nil
					}
				} else {
					this_.var1 = nil
				}
			} else if (o != nil) && (o.var2 != nil) {
				if j2gort.CheckNotNullString(o.var2.GetName()) == name {
					{
						r_31 := o
						r_31.var2 = o.var3
					}
					{
						r_32 := o
						r_32.var3 = o.var4
					}
					{
						r_33 := o
						r_33.var4 = o.var5
					}
					{
						r_34 := o
						r_34.var5 = nil
					}
					if o.IsEmpty() {
						this_.overflow = nil
					}
				} else if o.var3 != nil {
					if j2gort.CheckNotNullString(o.var3.GetName()) == name {
						{
							r_35 := o
							r_35.var3 = o.var4
						}
						{
							r_36 := o
							r_36.var4 = o.var5
						}
						{
							r_37 := o
							r_37.var5 = nil
						}
					} else if o.var4 != nil {
						if j2gort.CheckNotNullString(o.var4.GetName()) == name {
							{
								r_38 := o
								r_38.var4 = o.var5
							}
							{
								r_39 := o
								r_39.var5 = nil
							}
						} else if (o.var5 != nil) && (j2gort.CheckNotNullString(o.var5.GetName()) == name) {
							{
								r_40 := o
								r_40.var5 = nil
							}
						}
					}
				}
			}
		}
	}
}

// Java method: TypedScope.clearVarsInternal()
func (this_ *TypedScope) ClearVarsInternal() {
	this_.var0 = nil
	this_.var1 = nil
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if o != nil {
		{
			r_41 := o
			r_41.var2 = nil
		}
		{
			r_42 := o
			r_42.var3 = nil
		}
		{
			r_43 := o
			r_43.var4 = nil
		}
		{
			r_44 := o
			r_44.var5 = nil
		}
		{
			r_45 := o
			r_45.qualifiedRootNames = nil
		}
		{
			r_46 := o
			r_46.qualifiedVarHashes = nil
		}
		{
			r_47 := o
			r_47.qualifiedVarHashCount = int32(0)
		}
		if (o.vars != nil) && !(o.vars.IsEmpty()) {
			o.vars.Clear()
		}
		if o.IsEmpty() {
			this_.overflow = nil
		}
	}
}

// Java method: TypedScope.resolveOwnVarTypes(com.google.javascript.rhino.ErrorReporter errorReporter)
func (this_ *TypedScope) ResolveOwnVarTypes(errorReporter pkg_rhino.ErrorReporter) {
	_ = errorReporter
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		{
			cur_48 := o.vars.IterValues()
			for cur_48.HasNext() {

				var_ := cur_48.
					Next()
				_ = var_
				var_.ResolveType(errorReporter)
			}
		}
		return
	}
	if this_.var0 != nil {
		this_.var0.ResolveType(errorReporter)
		if this_.var1 != nil {
			this_.var1.ResolveType(errorReporter)
			if (o != nil) && (o.var2 != nil) {
				o.var2.ResolveType(errorReporter)
				if o.var3 != nil {
					o.var3.ResolveType(errorReporter)
					if o.var4 != nil {
						o.var4.ResolveType(errorReporter)
						if o.var5 != nil {
							o.var5.ResolveType(errorReporter)
						}
					}
				}
			}
		}
	}
}

// Java method: TypedScope.getVarIterable()
func (this_ *TypedScope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](o.vars.LiveValues())
	}
	if this_.var0 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_EmptyList())
	}
	if this_.var1 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.SingletonListPtr(this_.var0))
	}
	if (o == nil) || (o.var2 == nil) {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1)))
	}
	if o.var3 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(o.var2)))
	}
	if o.var4 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_4(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(o.var2), j2gort.PtrToAny(o.var3)))
	}
	if o.var5 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_5(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(o.var2), j2gort.PtrToAny(o.var3), j2gort.PtrToAny(o.var4)))
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_6(j2gort.PtrToAny(this_.var0), j2gort.PtrToAny(this_.var1), j2gort.PtrToAny(o.var2), j2gort.PtrToAny(o.var3), j2gort.PtrToAny(o.var4), j2gort.PtrToAny(o.var5)))
}

// Java method: TypedScope.putOwnVarsIfAbsent(java.util.LinkedHashMap<java.lang.String,com.google.javascript.jscomp.TypedVar> dst)
func (this_ *TypedScope) PutOwnVarsIfAbsent(dst *j2gort.Ext_java_util_LinkedHashMap) {
	_ = dst
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		{
			cur_50 := o.vars.IterValues()
			for cur_50.HasNext() {

				v := cur_50.
					Next()
				_ = v
				dst.PutIfAbsent(j2gort.BoxString(v.GetName()), j2gort.PtrToAny(v))
			}
		}
	} else if this_.var0 != nil {
		dst.PutIfAbsent(j2gort.BoxString(this_.var0.GetName()), j2gort.PtrToAny(this_.var0))
		if this_.var1 != nil {
			dst.PutIfAbsent(j2gort.BoxString(this_.var1.GetName()), j2gort.PtrToAny(this_.var1))
			if (o != nil) && (o.var2 != nil) {
				dst.PutIfAbsent(j2gort.BoxString(o.var2.GetName()), j2gort.PtrToAny(o.var2))
				if o.var3 != nil {
					dst.PutIfAbsent(j2gort.BoxString(o.var3.GetName()), j2gort.PtrToAny(o.var3))
					if o.var4 != nil {
						dst.PutIfAbsent(j2gort.BoxString(o.var4.GetName()), j2gort.PtrToAny(o.var4))
						if o.var5 != nil {
							dst.PutIfAbsent(j2gort.BoxString(o.var5.GetName()), j2gort.PtrToAny(o.var5))
						}
					}
				}
			}
		}
	}
}

// Java method: TypedScope.getAllSymbols()
func (this_ *TypedScope) GetAllSymbols() j2gort.Ext_java_lang_Iterable {
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](o.vars.LiveValues())))
	}
	return this_.GetVarIterable()
}

// Java method: TypedScope.addVarsTo(java.util.ArrayList<com.google.javascript.jscomp.TypedVar> dst)
func (this_ *TypedScope) AddVarsTo(dst j2gort.Ext_java_util_ArrayList) {
	_ = dst
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		{
			cur_52 := o.vars.IterValues()
			for cur_52.HasNext() {

				v := cur_52.
					Next()
				_ = v
				dst.Add_1(j2gort.PtrToAny(v))
			}
		}
	} else if this_.var0 != nil {
		dst.Add_1(j2gort.PtrToAny(this_.var0))
		if this_.var1 != nil {
			dst.Add_1(j2gort.PtrToAny(this_.var1))
			if (o != nil) && (o.var2 != nil) {
				dst.Add_1(j2gort.PtrToAny(o.var2))
				if o.var3 != nil {
					dst.Add_1(j2gort.PtrToAny(o.var3))
					if o.var4 != nil {
						dst.Add_1(j2gort.PtrToAny(o.var4))
						if o.var5 != nil {
							dst.Add_1(j2gort.PtrToAny(o.var5))
						}
					}
				}
			}
		}
	}
}

// Java method: TypedScope.addVarNamesTo(java.util.LinkedHashSet<java.lang.String> dst)
func (this_ *TypedScope) AddVarNamesTo(dst *j2gort.Ext_java_util_LinkedHashSet) {
	_ = dst
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		{
			cur_54 := o.vars.IterKeys()
			for cur_54.HasNext() {
				name := cur_54.
					Next()
				_ = name
				j2gort.SetAddStr(dst, name)
			}
		}
	} else if this_.var0 != nil {
		j2gort.SetAddStr(dst, this_.var0.GetName())
		if this_.var1 != nil {
			j2gort.SetAddStr(dst, this_.var1.GetName())
			if (o != nil) && (o.var2 != nil) {
				j2gort.SetAddStr(dst, o.var2.GetName())
				if o.var3 != nil {
					j2gort.SetAddStr(dst, o.var3.GetName())
					if o.var4 != nil {
						j2gort.SetAddStr(dst, o.var4.GetName())
						if o.var5 != nil {
							j2gort.SetAddStr(dst, o.var5.GetName())
						}
					}
				}
			}
		}
	}
}

// Java method: TypedScope.getVarCount()
func (this_ *TypedScope) GetVarCount() int32 {
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		return o.vars.Size()
	}
	if this_.var0 == nil {
		return int32(0)
	}
	if this_.var1 == nil {
		return int32(1)
	}
	if (o == nil) || (o.var2 == nil) {
		return int32(2)
	}
	if o.var3 == nil {
		return int32(3)
	}
	if o.var4 == nil {
		return int32(4)
	}
	if o.var5 == nil {
		return int32(5)
	} else {
		return int32(6)
	}
}

// Java method: TypedScope.makeImplicitVar(com.google.javascript.jscomp.AbstractScope.ImplicitVar var)
func (this_ *TypedScope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) any {
	_ = var_
	if this_.IsGlobal() {
		return nil
	} else {
		var iifeCond_1 bool
		{

			re_55 := var_
			ae_56 := StaticGet_AbstractScope_ImplicitVar_EXPORTS()
			if re_55 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_55 == ae_56
		}
		if iifeCond_1 {
			return nil
		}
	}
	if this_.tsArena_ != nil {
		return allocPooled_TypedVar(&this_.tsArena_.slab_TypedVar, false, var_.FieldGet_name(), pkg_rhino.
			Node(0), this_.priv_TypedScope_getImplicitVarType_5728(var_), this_, int32(-1), nil)
	}
	return NewTypedVar(false, var_.FieldGet_name(), pkg_rhino.Node(0), this_.priv_TypedScope_getImplicitVarType_5728(var_), this_, int32(-1), nil)

}

// Java method: TypedScope.hasOwnImplicitSlot(com.google.javascript.jscomp.AbstractScope.@org.jspecify.annotations.Nullable ImplicitVar name)
func (this_ *TypedScope) HasOwnImplicitSlot(name AbstractScope_ImplicitVar) bool {
	_ = name
	return (((name != 0) && !(func() bool {
		re_57 := name
		ae_58 := StaticGet_AbstractScope_ImplicitVar_EXPORTS()
		if re_57 == 0 {
			j2gort.ThrowNPE()
		}
		return re_57 == ae_58
	}())) && name.IsMadeByScope(this_))
}

// Java method: TypedScope.getImplicitVarType(com.google.javascript.jscomp.AbstractScope.ImplicitVar var)
func (this_ *TypedScope) priv_TypedScope_getImplicitVarType_5728(var_ AbstractScope_ImplicitVar) pkg_rhino.JSType {
	_ = var_
	switch (var_).Ordinal() {
	case int32(0) /* ImplicitVar.ARGUMENTS */ :
		var globalArgs *TypedVar = this_.GetGlobalTypedScope().LookupTypedVar("arguments")
		_ = globalArgs
		if (globalArgs != nil) && globalArgs.IsExtern() {
			return globalArgs.GetType()
		} else {
			return nil
		}
	case int32(3) /* ImplicitVar.THIS */ :
		return this_.GetTypeOfThis()
	case int32(2) /* ImplicitVar.SUPER */ :
		var receiverType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(this_.GetTypeOfThis())
		_ = receiverType
		if receiverType == nil || receiverType.IsJavaNil_() {
			return nil
		} else if receiverType.IsInstanceType() {
			var superclassCtor *pkg_rhino.FunctionType = receiverType.GetSuperClassConstructor()
			_ = superclassCtor
			if superclassCtor == nil {
				return nil
			} else {
				return pkg_rhino.Upcast_JSType(superclassCtor.GetInstanceType())
			}
		} else {
			return pkg_rhino.Upcast_JSType(receiverType.GetImplicitPrototype())
		}
	case int32(1) /* ImplicitVar.EXPORTS */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "TypedScopes should not contain an implicit 'exports'")))
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: TypedScope.getDeclarativelyUnboundVarsWithoutTypes()
func (this_ *TypedScope) GetDeclarativelyUnboundVarsWithoutTypes() j2gort.Ext_java_lang_Iterable {
	var result j2gort.Ext_java_util_ArrayList = nil
	_ = result
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.vars != nil) {
		{
			cur_59 := o.vars.IterValues()
			for cur_59.HasNext() {

				var_ := cur_59.
					Next()
				_ = var_
				if this_.priv_TypedScope_isDeclarativelyUnboundVarWithoutType_5731(var_) {
					if result == nil {
						result = j2gort.NewExt_java_util_ArrayList()
					}
					result.Add_1(j2gort.PtrToAny(var_))
				}
			}
		}
	} else if this_.var0 != nil {
		result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(this_.var0, result)
		if this_.var1 != nil {
			result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(this_.var1, result)
			if (o != nil) && (o.var2 != nil) {
				result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(o.var2, result)
				if o.var3 != nil {
					result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(o.var3, result)
					if o.var4 != nil {
						result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(o.var4, result)
						if o.var5 != nil {
							result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(o.var5, result)
						}
					}
				}
			}
		}
	}
	if result != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](result)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
}

// Java method: TypedScope.addIfDeclarativelyUnboundWithoutType(com.google.javascript.jscomp.TypedVar var, java.util.@org.jspecify.annotations.Nullable ArrayList<com.google.javascript.jscomp.TypedVar> result)
func (this_ *TypedScope) priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5730(var_ *TypedVar, result j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList {
	_ = var_
	_ = result
	if this_.priv_TypedScope_isDeclarativelyUnboundVarWithoutType_5731(var_) {
		if result == nil {
			result = j2gort.NewExt_java_util_ArrayList()
		}
		result.Add_1(j2gort.PtrToAny(var_))
	}
	return result
}

// Java method: TypedScope.isDeclarativelyUnboundVarWithoutType(com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScope) priv_TypedScope_isDeclarativelyUnboundVarWithoutType_5731(var_ *TypedVar) bool {
	_ = var_
	return ((((var_.GetParentNode() != 0) && j2gort.IfaceEquals[pkg_rhino.JSType](var_.GetType(), nil)) && var_.GetParentNode().IsVar()) && !(var_.IsExtern()))
}

// Java method: TypedScope.getGlobalTypedScope()
func (this_ *TypedScope) GetGlobalTypedScope() *TypedScope {
	var s *TypedScope = this_
	_ = s
	for s.parent != nil {
		s = s.parent
	}
	return s
}

// Java method: TypedScope.findRootVar(java.lang.String name)
func (this_ *TypedScope) priv_TypedScope_findRootVar_5733(name string) *TypedVar {
	_ = name
	var implicit AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(name)
	_ = implicit
	if implicit == 0 {
		{
			var scope *TypedScope = this_
			_ = scope
			for ; scope != nil; scope = scope.parent {
				var var_ *TypedVar = scope.priv_TypedScope_findOwnDeclaredVar_5726(name)
				_ = var_
				if var_ != nil {
					return var_
				}
			}
		}
		return nil
	}
	{
		var scope_2 *TypedScope = this_
		_ = scope_2
		for ; scope_2 != nil; scope_2 = scope_2.parent {
			var var__2 *TypedVar = scope_2.priv_TypedScope_findOwnDeclaredVar_5726(name)
			_ = var__2
			if var__2 != nil {
				return var__2
			}
			var__2 = j2gort.CastPtr[TypedVar](scope_2.GetOwnImplicitSlot(implicit))
			if var__2 != nil {
				return var__2
			}
		}
	}
	return nil
}

// Java method: TypedScope.findScopeForRootName(java.lang.String name)
func (this_ *TypedScope) FindScopeForRootName(name string) *TypedScope {
	_ = name
	var implicit AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(name)
	_ = implicit
	if implicit == 0 {
		{
			var scope *TypedScope = this_
			_ = scope
			for ; scope != nil; scope = scope.parent {
				if scope.priv_TypedScope_findOwnDeclaredVar_5726(name) != nil {
					return scope
				}
			}
		}
		return nil
	}
	{
		var scope_2 *TypedScope = this_
		_ = scope_2
		for ; scope_2 != nil; scope_2 = scope_2.parent {
			if scope_2.priv_TypedScope_findOwnDeclaredVar_5726(name) != nil {
				return scope_2
			}
			if !(scope_2.IsGlobal()) && scope_2.HasOwnImplicitSlot(implicit) {
				return scope_2
			}
		}
	}
	return nil
}

// Java method: TypedScope.lookupTypedVar(java.lang.String name)
func (this_ *TypedScope) LookupTypedVar(name string) *TypedVar {
	_ = name
	var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dot
	if dot < int32(0) {
		return this_.priv_TypedScope_findRootVar_5733(name)
	}
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.qualifiedRootNames != nil) {
		var ownSlot *TypedVar = this_.priv_TypedScope_findOwnDeclaredVar_5726(name)
		_ = ownSlot
		if ownSlot != nil {
			return ownSlot
		}
	}
	if this_.parent == nil {
		return nil
	}
	var rootName string = j2gort.JStrSub(name, int32(0), dot)
	_ = rootName
	var rootVar *TypedVar = this_.priv_TypedScope_findRootVar_5733(rootName)
	_ = rootVar
	if rootVar == nil {
		var global *TypedScope = this_.parent.GetGlobalTypedScope()
		_ = global
		if global.HasQualifiedOwnVarWithRoot(rootName) {
			return global.priv_TypedScope_findOwnDeclaredVar_5726(name)
		} else {
			return nil
		}
	}
	var rootScope *TypedScope = j2gort.CastPtr[TypedScope](rootVar.GetScope())
	_ = rootScope
	if rootScope.HasQualifiedOwnVarWithRoot(rootName) {
		return rootScope.priv_TypedScope_findOwnDeclaredVar_5726(name)
	} else {
		return nil
	}
}

// Java method: TypedScope.getVar(java.lang.String name)
func (this_ *TypedScope) GetVar(name string) any {
	_ = name
	return j2gort.PtrToAny(this_.LookupTypedVar(name))
}

// Java method: TypedScope.getTypeThroughNamespace(java.lang.String moduleId)
func (this_ *TypedScope) GetTypeThroughNamespace(moduleId string) pkg_rhino.JSType {
	_ = moduleId
	var split int32 = j2gort.JStrLastIndexOfCp(moduleId, int32(uint16(46)))
	_ = split
	if split >= int32(0) {
		var parentName string = j2gort.JStrSub(moduleId, int32(0), split)
		_ = parentName
		var prop string = j2gort.JStrSubFrom(moduleId, (split + int32(1)))
		_ = prop
		var parentType pkg_rhino.JSType = this_.GetTypeThroughNamespace(parentName)
		_ = parentType
		if (parentType == nil || parentType.IsJavaNil_()) || j2gort.IfaceEquals[pkg_rhino.ObjectType](parentType.GetJSTypeBase().ToMaybeObjectType(), nil) {
			return nil
		}
		return parentType.GetJSTypeBase().AssertObjectType().GetObjectTypeBase().GetPropertyType_1_string(prop)
	} else {
		var var_ *TypedVar = this_.LookupTypedVar(moduleId)
		_ = var_
		if var_ != nil {
			return var_.GetType()
		} else {
			return nil
		}
	}
}

// Java method: TypedScope.getTopmostScopeOfEventualDeclaration(java.lang.String name)
func (this_ *TypedScope) GetTopmostScopeOfEventualDeclaration(name string) pkg_rhino.StaticScope {
	_ = name
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (j2gort.CastPtr[TypedVar](this_.GetOwnSlot(name)) != nil) || (((o != nil) && (o.reservedMask != int32(0))) && o.ContainsReservedName(name)) {
		return this_
	} else if j2gort.CastPtr[TypedScope](this_.GetParent()) == nil {
		return nil
	} else {
		return j2gort.CastPtr[TypedScope](this_.GetParent()).GetTopmostScopeOfEventualDeclaration(name)
	}
}

func (this_ *TypedScope) LookupQualifiedName_1_QualifiedName(p0_ pkg_rhino.QualifiedName) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_QualifiedName(this_, p0_)
}

func (this_ *TypedScope) LookupQualifiedName_1_string(p0_ string) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_string(this_, p0_)
}

func (this_ *TypedScope) LookupQualifiedNameNode(p0_ pkg_rhino.Node) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedNameNode(this_, p0_)
}

// Java class: com.google.javascript.jscomp.TypedScope.Overflow
type TypedScope_Overflow struct {
	module                *Module                               `gson:"module"`
	reservedNames         []string                              `gson:"reservedNames"`
	reservedLargeSet      *j2gort.Ext_java_util_LinkedHashSet   `gson:"reservedLargeSet"`
	var2                  *TypedVar                             `gson:"var2"`
	var3                  *TypedVar                             `gson:"var3"`
	var4                  *TypedVar                             `gson:"var4"`
	var5                  *TypedVar                             `gson:"var5"`
	vars                  *j2gort.OrderedMap[string, *TypedVar] `gson:"vars"`
	qualifiedRootNames    *j2gort.PlainSet[string]              `gson:"qualifiedRootNames"`
	qualifiedVarHashes    []int32                               `gson:"qualifiedVarHashes"`
	reservedMask          int32                                 `gson:"reservedMask"`
	qualifiedVarHashCount int32                                 `gson:"qualifiedVarHashCount"`
	idHash_               atomic.Int32
}

func (this_ *TypedScope_Overflow) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScope_Overflow) IsJavaNil_() bool { return this_ == nil }
func (*TypedScope_Overflow) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScope$Overflow", "Overflow"
}
func (this_ *TypedScope_Overflow) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScope_Overflow) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScope.Overflow", "com.google.javascript.jscomp.TypedScope$Overflow", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScope_Overflow](v_)
	})
}

func (this_ *TypedScope_Overflow) initFields_TypedScope_Overflow() {
	_ = this_
}

// Java constructor: Overflow()
func NewTypedScope_Overflow() *TypedScope_Overflow {
	this_ := &TypedScope_Overflow{}
	InitTypedScope_Overflow(this_)
	return this_
}

func InitTypedScope_Overflow(this_ *TypedScope_Overflow) {
	this_.initFields_TypedScope_Overflow()
}

// Java method: Overflow.initReservedArray(java.lang.String[] names)
func (this_ *TypedScope_Overflow) InitReservedArray(names []string) {
	_ = names
	var len_ int32 = int32(len(names))
	_ = len_
	if len_ == int32(0) {
		return
	}
	if len_ <= int32(32) {
		this_.reservedNames = names
		if len_ == int32(32) {
			this_.reservedMask = int32(-1)
		} else {
			this_.reservedMask = ((javalang.Shl32(int32(1), int(len_))) - int32(1))
		}
	} else {
		var set *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32((((len_ * int32(4)) / int32(3)) + int32(1)))
		_ = set
		{
			var i int32 = int32(0)
			_ = i
			for ; i < len_; i = i + 1 {
				j2gort.SetAddStr(set, names[i])
			}
		}
		this_.reservedLargeSet = set
		this_.reservedMask = set.Size()
	}
}

// Java method: Overflow.initReservedSet(java.util.Set<java.lang.String> set)
func (this_ *TypedScope_Overflow) InitReservedSet(set j2gort.Ext_java_util_Set) {
	_ = set
	var len_ int32 = set.Size()
	_ = len_
	if len_ == int32(0) {
		return
	}
	if len_ <= int32(32) {
		var arr []string = j2gort.NewStringArray(len_)
		_ = arr
		var i int32 = int32(0)
		_ = i
		{
			cur_61 := j2gort.NewIterCursor(j2gort.CheckNotNull(set))
			for cur_61.HasNext() {
				s := cur_61.NextStr()
				_ = s
				{
					a_62 := arr
					var i_63 int32
					{

						old_64 := i
						nw_65 := old_64 + 1
						i = nw_65
						i_63 = old_64
					}
					a_62[i_63] = s
				}
			}
		}
		this_.reservedNames = arr
		if len_ == int32(32) {
			this_.reservedMask = int32(-1)
		} else {
			this_.reservedMask = ((javalang.Shl32(int32(1), int(len_))) - int32(1))
		}
	} else {
		var copy_ *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](set))
		_ = copy_
		this_.reservedLargeSet = copy_
		this_.reservedMask = copy_.Size()
	}
}

// Java method: Overflow.hasReservedNames()
func (this_ *TypedScope_Overflow) HasReservedNames() bool {
	return (this_.reservedMask != int32(0))
}

// Java method: Overflow.containsReservedName(java.lang.String name)
func (this_ *TypedScope_Overflow) ContainsReservedName(name string) bool {
	_ = name
	if this_.reservedLargeSet != nil {
		return j2gort.SetContainsStr(this_.reservedLargeSet, name)
	}
	var m int32 = this_.reservedMask
	_ = m
	var arr []string = this_.reservedNames
	_ = arr
	for m != int32(0) {
		var idx int32 = j2gort.IntTrailingZeros(m)
		_ = idx
		{
			var iifeCond_1 bool
			{

				sr_66 := name
				se_67 := arr[idx]
				iifeCond_1 = j2gort.CheckNotNullString(sr_66) == se_67
			}
			if iifeCond_1 {
				return true
			}
		}
		m = (m & (m - int32(1)))
	}
	return false
}

// Java method: Overflow.removeReservedAndCheckEmpty(java.lang.String name)
func (this_ *TypedScope_Overflow) RemoveReservedAndCheckEmpty(name string) bool {
	_ = name
	if this_.reservedLargeSet != nil {
		if j2gort.SetRemoveStr(this_.reservedLargeSet, name) {
			var next int32 = (this_.reservedMask - int32(1))
			_ = next
			this_.reservedMask = next
			if next == int32(0) {
				this_.reservedLargeSet = nil
				return true
			}
		}
		return (this_.reservedMask == int32(0))
	}
	var m int32 = this_.reservedMask
	_ = m
	var arr []string = this_.reservedNames
	_ = arr
	for m != int32(0) {
		var bit int32 = (m & -(m))
		_ = bit
		var idx int32 = j2gort.IntTrailingZeros(bit)
		_ = idx
		{
			var iifeCond_1 bool
			{

				sr_68 := name
				se_69 := arr[idx]
				iifeCond_1 = j2gort.CheckNotNullString(sr_68) == se_69
			}
			if iifeCond_1 {
				var next_2 int32 = (this_.reservedMask & ^(bit))
				_ = next_2
				this_.reservedMask = next_2
				if next_2 == int32(0) {
					this_.reservedNames = nil
					return true
				}
				return false
			}
		}
		m = (m & ^(bit))
	}
	return (this_.reservedMask == int32(0))
}

// Java method: Overflow.clearReserved()
func (this_ *TypedScope_Overflow) ClearReserved() {
	this_.reservedNames = nil
	this_.reservedMask = int32(0)
	this_.reservedLargeSet = nil
}

// Java method: Overflow.reservedToString()
func (this_ *TypedScope_Overflow) ReservedToString() string {
	if this_.reservedLargeSet != nil {
		return this_.reservedLargeSet.String()
	}
	var left j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(j2gort.IntBitCount(this_.reservedMask))
	_ = left
	var m int32 = this_.reservedMask
	_ = m
	var arr []string = this_.reservedNames
	_ = arr
	for m != int32(0) {
		var idx int32 = j2gort.IntTrailingZeros(m)
		_ = idx
		left.Add_1(j2gort.BoxString(arr[idx]))
		m = (m & (m - int32(1)))
	}
	return left.String()
}

// Java method: Overflow.isEmpty()
func (this_ *TypedScope_Overflow) IsEmpty() bool {
	return ((((((this_.module == nil) && (this_.reservedMask == int32(0))) && (this_.var2 == nil)) && (this_.vars == nil)) && (this_.qualifiedRootNames == nil)) && (this_.qualifiedVarHashes == nil))
}

// Java method: Overflow.resetForPool()
func (this_ *TypedScope_Overflow) ResetForPool() {
	this_.module = nil
	this_.reservedNames = nil
	this_.reservedMask = int32(0)
	this_.reservedLargeSet = nil
	this_.var2 = nil
	this_.var3 = nil
	this_.var4 = nil
	this_.var5 = nil
	this_.vars = nil
	this_.qualifiedRootNames = nil
	this_.qualifiedVarHashes = nil
	this_.qualifiedVarHashCount = int32(0)
}

type typedScopeArena_ struct {
	slab_TypedScope j2gort.PooledSlab[TypedScope]
	slab_TypedVar   j2gort.PooledSlab[TypedVar]
}

func (a_ *typedScopeArena_) release_() {
	if a_ == nil {
		return
	}
	a_.slab_TypedScope.ReleaseWith(func(s_ *TypedScope) {
		if s_.overflow != nil {
			if s_.overflow.vars != nil {
				j2gort.ReleaseOrderedMap(s_.overflow.vars)
				s_.overflow.vars = nil
			}
			s_.overflow.ResetForPool()
			s_.overflow = nil
		}
		*s_ = TypedScope{}
	})
	a_.slab_TypedVar.ReleaseWith(func(v_ *TypedVar) { *v_ = TypedVar{} })
}
