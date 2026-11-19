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
	InitTypedScope_4_TypedScope_Node_slice_string_Module(this_, parent, rootNode, nil, nil)
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
			{
				r_2 := o
				r_2.reservedNames = NewTypedScope_ReservedNames_1_Ext_java_util_Set(reservedNames)
			}
		}
		{
			r_3 := o
			r_3.module = module
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
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_TypedScope()
	this_.CheckChildScope(ToIface_AbstractScope(parent))
	this_.parent = parent
	this_.depth = (parent.depth + int32(1))
	this_.isBottom = false
	var hasReserved bool = ((reservedNamesOwned != nil) && (int32(len(reservedNamesOwned)) != int32(0)))
	_ = hasReserved
	if hasReserved || (module != nil) {
		var val_4 TypedScope_Overflow
		InitTypedScope_Overflow(&val_4)
		var o *TypedScope_Overflow = &val_4
		_ = o
		if hasReserved {
			{
				r_5 := o
				r_5.reservedNames = NewTypedScope_ReservedNames_1_slice_string(reservedNamesOwned)
			}
		}
		{
			r_6 := o
			r_6.module = module
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
func (this_ *TypedScope) priv_TypedScope_ensureOverflow_5688() *TypedScope_Overflow {
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if o == nil {
		o = NewTypedScope_Overflow()
		this_.overflow = o
	}
	return o
}

// Java static method: TypedScope.createGlobalScope(com.google.javascript.rhino.Node rootNode)
func TypedScope_CreateGlobalScope(rootNode pkg_rhino.Node) *TypedScope {
	_ = rootNode
	return NewTypedScope_2_Node_bool(rootNode, false)
}

// Java static method: TypedScope.createLatticeBottom(com.google.javascript.rhino.Node rootNode)
func TypedScope_CreateLatticeBottom(rootNode pkg_rhino.Node) *TypedScope {
	_ = rootNode
	return NewTypedScope_2_Node_bool(rootNode, true)
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
		if !((o.reservedNames == nil) || o.reservedNames.IsEmpty()) {
			j2gort.CheckStateFmt(false, "Expected %s to have no reserved names, found: %s. This probably indicates a bug in TypedScopeCreator where it is failing to declare a variable.", this_, j2gort.PtrToAny(o.reservedNames))
		}
		{
			r_7 := o
			r_7.reservedNames = nil
		}
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
		return j2gort.AsIface[pkg_rhino.JSType](pkg_rhino.ObjectType_Cast(root.GetJSType()))
	} else if NodeUtil_IsNonArrowFunction(root) {
		var nodeType pkg_rhino.JSType = root.GetJSType()
		_ = nodeType
		if (nodeType != nil && !nodeType.IsJavaNil_()) && nodeType.IsFunctionType() {
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
			return j2gort.AsIface[pkg_rhino.JSType](classType.AssertFunctionType().GetInstanceType())
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
	if (o != nil) && (o.reservedNames != nil) {
		if o.reservedNames.RemoveAndCheckEmpty(name) {
			{
				r_8 := o
				r_8.reservedNames = nil
			}
			if o.IsEmpty() {
				this_.overflow = nil
			}
		}
	}
	var var_ *TypedVar = NewTypedVar(inferred, name, nameNode, type_, this_, this_.GetVarCount(), input)
	_ = var_
	this_.DeclareInternal(name, var_)
	var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dot
	if dot >= int32(0) {
		var qo *TypedScope_Overflow = this_.priv_TypedScope_ensureOverflow_5688()
		_ = qo
		if qo.qualifiedRootNames == nil {
			{
				r_9 := qo
				r_9.qualifiedRootNames = j2gort.NewOrderedSetPtrCap[string](int32(4))
			}
		}
		qo.qualifiedRootNames.Add(j2gort.JStrSub(name, int32(0), dot))
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
			r_10 := o
			r_10.qualifiedVarHashes = table
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
						a_11 := next
						i_12 := s
						a_11[i_12] = existing
					}
				}
			}
		}
		table = next
		{
			r_13 := o
			r_13.qualifiedVarHashes = table
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
				a_14 := table
				i_15 := slot
				a_14[i_15] = h
			}
			r_16 := o
			r_16.qualifiedVarHashCount = r_16.qualifiedVarHashCount + 1
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
func (this_ *TypedScope) priv_TypedScope_findOwnDeclaredVar_5699(name string) *TypedVar {
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
	return j2gort.PtrToAny(this_.priv_TypedScope_findOwnDeclaredVar_5699(name))
}

// Java method: TypedScope.putOwnDeclaredVar(java.lang.String name, com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScope) PutOwnDeclaredVar(name string, arg_17 any) {
	_ = name
	var_ := j2gort.CastPtr[TypedVar](arg_17)
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
		o = this_.priv_TypedScope_ensureOverflow_5688()
	}
	if (o.var2 == nil) || (j2gort.CheckNotNullString(o.var2.GetName()) == name) {
		{
			r_18 := o
			r_18.var2 = var_
		}
		return
	}
	if (o.var3 == nil) || (j2gort.CheckNotNullString(o.var3.GetName()) == name) {
		{
			r_19 := o
			r_19.var3 = var_
		}
		return
	}
	if (o.var4 == nil) || (j2gort.CheckNotNullString(o.var4.GetName()) == name) {
		{
			r_20 := o
			r_20.var4 = var_
		}
		return
	}
	if (o.var5 == nil) || (j2gort.CheckNotNullString(o.var5.GetName()) == name) {
		{
			r_21 := o
			r_21.var5 = var_
		}
		return
	}
	{
		r_22 := o
		r_22.vars = j2gort.NewOrderedMapPtrCap[string, *TypedVar](func() int32 {
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
		r_23 := o
		r_23.var2 = nil
	}
	{
		r_24 := o
		r_24.var3 = nil
	}
	{
		r_25 := o
		r_25.var4 = nil
	}
	{
		r_26 := o
		r_26.var5 = nil
	}
}

// Java method: TypedScope.undeclareInteral(com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScope) UndeclareInteral(arg_27 any) {
	var_ := j2gort.CastPtr[TypedVar](arg_27)
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
					r_28 := o
					r_28.var2 = o.var3
				}
				{
					r_29 := o
					r_29.var3 = o.var4
				}
				{
					r_30 := o
					r_30.var4 = o.var5
				}
				{
					r_31 := o
					r_31.var5 = nil
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
						r_32 := o
						r_32.var2 = o.var3
					}
					{
						r_33 := o
						r_33.var3 = o.var4
					}
					{
						r_34 := o
						r_34.var4 = o.var5
					}
					{
						r_35 := o
						r_35.var5 = nil
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
						r_36 := o
						r_36.var2 = o.var3
					}
					{
						r_37 := o
						r_37.var3 = o.var4
					}
					{
						r_38 := o
						r_38.var4 = o.var5
					}
					{
						r_39 := o
						r_39.var5 = nil
					}
					if o.IsEmpty() {
						this_.overflow = nil
					}
				} else if o.var3 != nil {
					if j2gort.CheckNotNullString(o.var3.GetName()) == name {
						{
							r_40 := o
							r_40.var3 = o.var4
						}
						{
							r_41 := o
							r_41.var4 = o.var5
						}
						{
							r_42 := o
							r_42.var5 = nil
						}
					} else if o.var4 != nil {
						if j2gort.CheckNotNullString(o.var4.GetName()) == name {
							{
								r_43 := o
								r_43.var4 = o.var5
							}
							{
								r_44 := o
								r_44.var5 = nil
							}
						} else if (o.var5 != nil) && (j2gort.CheckNotNullString(o.var5.GetName()) == name) {
							{
								r_45 := o
								r_45.var5 = nil
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
			r_46 := o
			r_46.var2 = nil
		}
		{
			r_47 := o
			r_47.var3 = nil
		}
		{
			r_48 := o
			r_48.var4 = nil
		}
		{
			r_49 := o
			r_49.var5 = nil
		}
		{
			r_50 := o
			r_50.qualifiedRootNames = nil
		}
		{
			r_51 := o
			r_51.qualifiedVarHashes = nil
		}
		{
			r_52 := o
			r_52.qualifiedVarHashCount = int32(0)
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
			cur_53 := o.vars.IterValues()
			for cur_53.HasNext() {

				var_ := cur_53.
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
			cur_55 := o.vars.IterValues()
			for cur_55.HasNext() {

				v := cur_55.
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
			cur_57 := o.vars.IterValues()
			for cur_57.HasNext() {

				v := cur_57.
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
			cur_59 := o.vars.IterKeys()
			for cur_59.HasNext() {
				name := cur_59.
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

			re_60 := var_
			ae_61 := StaticGet_AbstractScope_ImplicitVar_EXPORTS()
			if re_60 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_60 == ae_61
		}
		if iifeCond_1 {
			return nil
		}
	}
	return NewTypedVar(false, var_.FieldGet_name(), pkg_rhino.Node(0), this_.priv_TypedScope_getImplicitVarType_5701(var_), this_, int32(-1), nil)
}

// Java method: TypedScope.hasOwnImplicitSlot(com.google.javascript.jscomp.AbstractScope.@org.jspecify.annotations.Nullable ImplicitVar name)
func (this_ *TypedScope) HasOwnImplicitSlot(name AbstractScope_ImplicitVar) bool {
	_ = name
	return (((name != 0) && !(func() bool {
		re_62 := name
		ae_63 := StaticGet_AbstractScope_ImplicitVar_EXPORTS()
		if re_62 == 0 {
			j2gort.ThrowNPE()
		}
		return re_62 == ae_63
	}())) && name.IsMadeByScope(this_))
}

// Java method: TypedScope.getImplicitVarType(com.google.javascript.jscomp.AbstractScope.ImplicitVar var)
func (this_ *TypedScope) priv_TypedScope_getImplicitVarType_5701(var_ AbstractScope_ImplicitVar) pkg_rhino.JSType {
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
				return j2gort.AsIface[pkg_rhino.JSType](superclassCtor.GetInstanceType())
			}
		} else {
			return j2gort.AsIface[pkg_rhino.JSType](receiverType.GetImplicitPrototype())
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
			cur_64 := o.vars.IterValues()
			for cur_64.HasNext() {

				var_ := cur_64.
					Next()
				_ = var_
				if this_.priv_TypedScope_isDeclarativelyUnboundVarWithoutType_5704(var_) {
					if result == nil {
						result = j2gort.NewExt_java_util_ArrayList()
					}
					result.Add_1(j2gort.PtrToAny(var_))
				}
			}
		}
	} else if this_.var0 != nil {
		result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(this_.var0, result)
		if this_.var1 != nil {
			result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(this_.var1, result)
			if (o != nil) && (o.var2 != nil) {
				result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(o.var2, result)
				if o.var3 != nil {
					result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(o.var3, result)
					if o.var4 != nil {
						result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(o.var4, result)
						if o.var5 != nil {
							result = this_.priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(o.var5, result)
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
func (this_ *TypedScope) priv_TypedScope_addIfDeclarativelyUnboundWithoutType_5703(var_ *TypedVar, result j2gort.Ext_java_util_ArrayList) j2gort.Ext_java_util_ArrayList {
	_ = var_
	_ = result
	if this_.priv_TypedScope_isDeclarativelyUnboundVarWithoutType_5704(var_) {
		if result == nil {
			result = j2gort.NewExt_java_util_ArrayList()
		}
		result.Add_1(j2gort.PtrToAny(var_))
	}
	return result
}

// Java method: TypedScope.isDeclarativelyUnboundVarWithoutType(com.google.javascript.jscomp.TypedVar var)
func (this_ *TypedScope) priv_TypedScope_isDeclarativelyUnboundVarWithoutType_5704(var_ *TypedVar) bool {
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
func (this_ *TypedScope) priv_TypedScope_findRootVar_5706(name string) *TypedVar {
	_ = name
	var implicit AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(name)
	_ = implicit
	if implicit == 0 {
		{
			var scope *TypedScope = this_
			_ = scope
			for ; scope != nil; scope = scope.parent {
				var var_ *TypedVar = scope.priv_TypedScope_findOwnDeclaredVar_5699(name)
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
			var var__2 *TypedVar = scope_2.priv_TypedScope_findOwnDeclaredVar_5699(name)
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
				if scope.priv_TypedScope_findOwnDeclaredVar_5699(name) != nil {
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
			if scope_2.priv_TypedScope_findOwnDeclaredVar_5699(name) != nil {
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
		return this_.priv_TypedScope_findRootVar_5706(name)
	}
	var o *TypedScope_Overflow = this_.overflow
	_ = o
	if (o != nil) && (o.qualifiedRootNames != nil) {
		var ownSlot *TypedVar = this_.priv_TypedScope_findOwnDeclaredVar_5699(name)
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
	var rootVar *TypedVar = this_.priv_TypedScope_findRootVar_5706(rootName)
	_ = rootVar
	if rootVar == nil {
		var global *TypedScope = this_.parent.GetGlobalTypedScope()
		_ = global
		if global.HasQualifiedOwnVarWithRoot(rootName) {
			return global.priv_TypedScope_findOwnDeclaredVar_5699(name)
		} else {
			return nil
		}
	}
	var rootScope *TypedScope = j2gort.CastPtr[TypedScope](rootVar.GetScope())
	_ = rootScope
	if rootScope.HasQualifiedOwnVarWithRoot(rootName) {
		return rootScope.priv_TypedScope_findOwnDeclaredVar_5699(name)
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
		if (parentType == nil || parentType.IsJavaNil_()) || j2gort.IfaceEquals[pkg_rhino.ObjectType](parentType.ToMaybeObjectType(), nil) {
			return nil
		}
		return parentType.AssertObjectType().GetPropertyType_1_string(prop)
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
	if (j2gort.CastPtr[TypedVar](this_.GetOwnSlot(name)) != nil) || (((o != nil) && (o.reservedNames != nil)) && o.reservedNames.Contains(name)) {
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

// Java class: com.google.javascript.jscomp.TypedScope.ReservedNames
type TypedScope_ReservedNames struct {
	names     []string                            `gson:"names"`
	largeSet  *j2gort.Ext_java_util_LinkedHashSet `gson:"largeSet"`
	remaining int32                               `gson:"remaining"`
	idHash_   atomic.Int32
}

func (this_ *TypedScope_ReservedNames) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedScope_ReservedNames) IsJavaNil_() bool { return this_ == nil }
func (*TypedScope_ReservedNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedScope$ReservedNames", "ReservedNames"
}
func (this_ *TypedScope_ReservedNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypedScope_ReservedNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedScope.ReservedNames", "com.google.javascript.jscomp.TypedScope$ReservedNames", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedScope_ReservedNames](v_)
	})
}

func (this_ *TypedScope_ReservedNames) initFields_TypedScope_ReservedNames() {
	_ = this_
}

// Java constructor: ReservedNames(java.lang.String[] names)
func NewTypedScope_ReservedNames_1_slice_string(names []string) *TypedScope_ReservedNames {
	this_ := &TypedScope_ReservedNames{}
	InitTypedScope_ReservedNames_1_slice_string(this_, names)
	return this_
}

func InitTypedScope_ReservedNames_1_slice_string(this_ *TypedScope_ReservedNames, names []string) {
	_ = names
	this_.initFields_TypedScope_ReservedNames()
	var len_ int32 = int32(len(names))
	_ = len_
	if len_ <= int32(8) {
		this_.names = names
		this_.remaining = len_
	} else {
		this_.names = nil
		this_.largeSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32((((len_ * int32(4)) / int32(3)) + int32(1)))
		{
			var i int32 = int32(0)
			_ = i
			for ; i < len_; i = i + 1 {
				j2gort.SetAddStr(this_.largeSet, names[i])
			}
		}
		this_.remaining = this_.largeSet.Size()
	}
}

// Java constructor: ReservedNames(java.util.Set<java.lang.String> set)
func NewTypedScope_ReservedNames_1_Ext_java_util_Set(set j2gort.Ext_java_util_Set) *TypedScope_ReservedNames {
	this_ := &TypedScope_ReservedNames{}
	InitTypedScope_ReservedNames_1_Ext_java_util_Set(this_, set)
	return this_
}

func InitTypedScope_ReservedNames_1_Ext_java_util_Set(this_ *TypedScope_ReservedNames, set j2gort.Ext_java_util_Set) {
	_ = set
	this_.initFields_TypedScope_ReservedNames()
	var len_ int32 = set.Size()
	_ = len_
	if len_ <= int32(8) {
		var arr []string = j2gort.NewStringArray(len_)
		_ = arr
		var i int32 = int32(0)
		_ = i
		{
			cur_66 := j2gort.NewIterCursor(j2gort.CheckNotNull(set))
			for cur_66.HasNext() {
				s := cur_66.NextStr()
				_ = s
				{
					a_67 := arr
					i_68 := func() int32 { old_69 := i; nw_70 := old_69 + 1; i = nw_70; return old_69 }()
					a_67[i_68] = s
				}
			}
		}
		this_.names = arr
		this_.remaining = len_
	} else {
		this_.names = nil
		this_.largeSet = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](set))
		this_.remaining = this_.largeSet.Size()
	}
}

// Java method: ReservedNames.contains(java.lang.String name)
func (this_ *TypedScope_ReservedNames) Contains(name string) bool {
	_ = name
	if this_.largeSet != nil {
		return j2gort.SetContainsStr(this_.largeSet, name)
	}
	var arr []string = this_.names
	_ = arr
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(arr))
		_ = len_
		for ; i < len_; i = i + 1 {
			{
				var iifeCond_1 bool
				{

					sr_71 := name
					se_72 := arr[i]
					iifeCond_1 = j2gort.CheckNotNullString(sr_71) == se_72
				}
				if iifeCond_1 {
					return true
				}
			}
		}
	}
	return false
}

// Java method: ReservedNames.removeAndCheckEmpty(java.lang.String name)
func (this_ *TypedScope_ReservedNames) RemoveAndCheckEmpty(name string) bool {
	_ = name
	if this_.largeSet != nil {
		if j2gort.SetRemoveStr(this_.largeSet, name) {
			this_.remaining = this_.remaining - 1
		}
		return (this_.remaining == int32(0))
	}
	var arr []string = this_.names
	_ = arr
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(arr))
		_ = len_
		for ; i < len_; i = i + 1 {
			{
				var iifeCond_1 bool
				{

					sr_73 := name
					se_74 := arr[i]
					iifeCond_1 = j2gort.CheckNotNullString(sr_73) == se_74
				}
				if iifeCond_1 {
					{
						a_75 := arr
						i_76 := i
						a_75[i_76] = j2gort.NullString
					}
					this_.remaining = this_.remaining - 1
					break
				}
			}
		}
	}
	return (this_.remaining == int32(0))
}

// Java method: ReservedNames.isEmpty()
func (this_ *TypedScope_ReservedNames) IsEmpty() bool {
	return (this_.remaining == int32(0))
}

// Java method: ReservedNames.toString()
func (this_ *TypedScope_ReservedNames) String() string {
	if this_.largeSet != nil {
		return this_.largeSet.String()
	}
	var left j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(this_.remaining)
	_ = left
	for _, el_77 := range j2gort.CheckNotNullSlice(this_.names) {
		s := el_77
		_ = s
		if s != j2gort.NullString {
			left.Add_1(j2gort.BoxString(s))
		}
	}
	return left.String()
}

// Java class: com.google.javascript.jscomp.TypedScope.Overflow
type TypedScope_Overflow struct {
	module                *Module                               `gson:"module"`
	reservedNames         *TypedScope_ReservedNames             `gson:"reservedNames"`
	var2                  *TypedVar                             `gson:"var2"`
	var3                  *TypedVar                             `gson:"var3"`
	var4                  *TypedVar                             `gson:"var4"`
	var5                  *TypedVar                             `gson:"var5"`
	vars                  *j2gort.OrderedMap[string, *TypedVar] `gson:"vars"`
	qualifiedRootNames    *j2gort.OrderedSet[string]            `gson:"qualifiedRootNames"`
	qualifiedVarHashes    []int32                               `gson:"qualifiedVarHashes"`
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

// Java method: Overflow.isEmpty()
func (this_ *TypedScope_Overflow) IsEmpty() bool {
	return ((((((this_.module == nil) && (this_.reservedNames == nil)) && (this_.var2 == nil)) && (this_.vars == nil)) && (this_.qualifiedRootNames == nil)) && (this_.qualifiedVarHashes == nil))
}
