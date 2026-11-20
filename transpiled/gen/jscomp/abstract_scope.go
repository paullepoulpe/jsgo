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
// Port of com/google/javascript/jscomp/AbstractScope.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractScope
type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetAbstractScopeBase() *AbstractScopeBase
	pkg_rhino.StaticScope
	GetDepth() int32
	GetParent() any
	String() string
	Untyped() *Scope
	Typed() *TypedScope
	Contains(AbstractScope) bool
	GetGlobalScope() AbstractScope
	MakeImplicitVar(AbstractScope_ImplicitVar) any
	GetOwnDeclaredVar(string) any
	PutOwnDeclaredVar(string, any)
	Undeclare(AbstractVar)
	UndeclareInteral(any)
	DeclareInternal(string, AbstractVar)
	ClearVarsInternal()
	HasOwnImplicitSlot(AbstractScope_ImplicitVar) bool
	HasOwnSlot(string) bool
	HasSlot(string) bool
	GetOwnImplicitSlot(AbstractScope_ImplicitVar) AbstractVar
	GetVar(string) any
	GetArgumentsVar() AbstractVar
	CanDeclare(string) bool
	GetVarIterable() j2gort.Ext_java_lang_Iterable
	PutOwnVarsIfAbsent(*j2gort.Ext_java_util_LinkedHashMap)
	GetAllAccessibleVariables() j2gort.Ext_java_lang_Iterable
	GetAllSymbols() j2gort.Ext_java_lang_Iterable
	AddVarsTo(j2gort.Ext_java_util_ArrayList)
	AddVarNamesTo(*j2gort.Ext_java_util_LinkedHashSet)
	GetVarCount() int32
	IsGlobal() bool
	IsLocal() bool
	IsBlockScope() bool
	IsStaticBlockScope() bool
	IsFunctionBlockScope() bool
	IsFunctionScope() bool
	IsModuleScope() bool
	IsMemberFieldDefScope() bool
	IsComputedFieldDefRhsScope() bool
	IsCatchScope() bool
	IsCfgRootScope() bool
	IsHoistScope() bool
	GetClosestHoistScope() AbstractScope
	GetClosestCfgRootScope() AbstractScope
	GetClosestContainerScope() AbstractScope
	CheckChildScope(AbstractScope)
	CheckRootScope()
	GetCommonParent(AbstractScope) AbstractScope
	HasSameContainerScope(AbstractScope) bool
	GetScopeOfThis() AbstractScope
}

func ToIface_AbstractScope[T interface {
	AbstractScope
	comparable
}](v T) AbstractScope {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AbstractScope[T AbstractScope](v T) AbstractScope {
	var vi AbstractScope = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.AbstractScope
type AbstractScopeBase struct {
	self_        AbstractScope
	implicitVars any            `gson:"implicitVars"`
	rootNode     pkg_rhino.Node `gson:"rootNode"`
	idHash_      atomic.Int32
}

func (this_ *AbstractScopeBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(AbstractScope)
}
func (this_ *AbstractScopeBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractScopeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractScope", "AbstractScope"
}
func (this_ *AbstractScopeBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractScope", "com.google.javascript.jscomp.AbstractScope", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractScope](v_)
	})
}
func (b *AbstractScopeBase) GetAbstractScopeBase() *AbstractScopeBase { return b }
func (*AbstractScopeBase) Implements_jscomp_AbstractScope()           {}
func (*AbstractScopeBase) Implements_rhino_StaticScope()              {}

func (this_ *AbstractScopeBase) initFields_AbstractScope() {
	_ = this_
	this_.implicitVars = nil
}

func InitAbstractScope(this_ *AbstractScopeBase, rootNode pkg_rhino.Node) {
	_ = rootNode
	this_.initFields_AbstractScope()
	this_.rootNode = j2gort.CheckNotNullNode(rootNode)
}

// Java method: AbstractScope.getDepth()
func (this_ *AbstractScopeBase) GetDepth() int32 {
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.getParent()
func (this_ *AbstractScopeBase) GetParent() any {
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.toString()
func (this_ *AbstractScopeBase) String() string {
	return ("Scope@" + j2gort.StrVal(this_.rootNode))
}

// Java method: AbstractScope.untyped()
func (this_ *AbstractScopeBase) Untyped() *Scope {
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "untyped() called, but not an untyped scope.")))
}

// Java method: AbstractScope.typed()
func (this_ *AbstractScopeBase) Typed() *TypedScope {
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "typed() called, but not a typed scope.")))
}

// Java method: AbstractScope.contains(S other)
func (this_ *AbstractScopeBase) Contains(other AbstractScope) bool {
	_ = other
	var s AbstractScope = j2gort.CheckNotNullIface(other)
	_ = s
	for s != nil && !s.IsJavaNil_() {
		if s == this_.self_ || ((s == nil || s.IsJavaNil_()) && (this_.self_ == nil || this_.self_.IsJavaNil_())) {
			return true
		}
		s = j2gort.AsIface[AbstractScope](s.GetParent())
	}
	return false
}

// Java method: AbstractScope.getRootNode()
func (this_ *AbstractScopeBase) GetRootNode_0() pkg_rhino.Node {
	return this_.rootNode
}

// Java method: AbstractScope.getGlobalScope()
func (this_ *AbstractScopeBase) GetGlobalScope() AbstractScope {
	var result AbstractScope = this_.priv_AbstractScope_thisScope_491()
	_ = result
	for !(j2gort.IfaceEquals[AbstractScope](j2gort.AsIface[AbstractScope](result.GetParent()), nil)) {
		result = j2gort.AsIface[AbstractScope](result.GetParent())
	}
	return result
}

// Java method: AbstractScope.getParentScope()
func (this_ *AbstractScopeBase) GetParentScope() any {
	return j2gort.AsIface[AbstractScope](this_.self_.GetParent())
}

// Java method: AbstractScope.makeImplicitVar(com.google.javascript.jscomp.AbstractScope.ImplicitVar type)
func (this_ *AbstractScopeBase) MakeImplicitVar(type_ AbstractScope_ImplicitVar) any {
	_ = type_
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.getOwnDeclaredVar(java.lang.String name)
func (this_ *AbstractScopeBase) GetOwnDeclaredVar(name string) any {
	_ = name
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.putOwnDeclaredVar(java.lang.String name, V var)
func (this_ *AbstractScopeBase) PutOwnDeclaredVar(name string, arg_1 any) {
	_ = name
	var_ := j2gort.AsIface[AbstractVar](arg_1)
	_ = var_
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.undeclare(V var)
func (this_ *AbstractScopeBase) Undeclare(var_ AbstractVar) {
	_ = var_
	if !(j2gort.IfaceEquals[AbstractScope](j2gort.AsIface[AbstractScope](var_.GetAbstractVarBase().GetScope()), this_.self_)) {
		j2gort.CheckState(false)
	}
	var existing AbstractVar = j2gort.AsIface[AbstractVar](this_.self_.GetOwnDeclaredVar(var_.GetAbstractVarBase().GetName()))
	_ = existing
	if !((existing != nil && !existing.IsJavaNil_()) && existing.GetAbstractVarBase().Equals(var_)) {
		j2gort.CheckState(false)
	}
	this_.self_.UndeclareInteral(var_)
}

// Java method: AbstractScope.undeclareInteral(V var)
func (this_ *AbstractScopeBase) UndeclareInteral(arg_2 any) {
	var_ := j2gort.AsIface[AbstractVar](arg_2)
	_ = var_
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.declareInternal(java.lang.String name, V var)
func (this_ *AbstractScopeBase) DeclareInternal(name string, var_ AbstractVar) {
	_ = name
	_ = var_
	if !(this_.HasOwnSlot(name) || this_.CanDeclare(name)) {
		j2gort.CheckStateFmt(false, "Illegal shadow: %s", j2gort.BoxNode(var_.GetAbstractVarBase().GetNode()))
	}
	{
		var iifeCond_1 bool
		{

			sr_3 := name
			se_4 := var_.GetAbstractVarBase().GetName()
			iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
		}
		if !(iifeCond_1) {
			j2gort.CheckArgument(false)
		}
	}
	this_.self_.PutOwnDeclaredVar(name, var_)
}

// Java method: AbstractScope.clearVarsInternal()
func (this_ *AbstractScopeBase) ClearVarsInternal() {
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.hasOwnImplicitSlot(com.google.javascript.jscomp.AbstractScope.@org.jspecify.annotations.Nullable ImplicitVar name)
func (this_ *AbstractScopeBase) HasOwnImplicitSlot(name AbstractScope_ImplicitVar) bool {
	_ = name
	return ((name != 0) && name.IsMadeByScope(this_.self_))
}

// Java method: AbstractScope.hasOwnSlot(java.lang.String name)
func (this_ *AbstractScopeBase) HasOwnSlot(name string) bool {
	_ = name
	return (!(j2gort.IfaceEquals[AbstractVar](j2gort.AsIface[AbstractVar](this_.self_.GetOwnDeclaredVar(name)), nil)) || this_.self_.HasOwnImplicitSlot(AbstractScope_ImplicitVar_Of(name)))
}

// Java method: AbstractScope.hasSlot(java.lang.String name)
func (this_ *AbstractScopeBase) HasSlot(name string) bool {
	_ = name
	var implicit AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(name)
	_ = implicit
	if implicit == 0 {
		{
			var scope AbstractScope = this_.self_
			_ = scope
			for ; scope != nil && !scope.IsJavaNil_(); scope = j2gort.AsIface[AbstractScope](scope.GetParent()) {
				if !(j2gort.IfaceEquals[AbstractVar](j2gort.AsIface[AbstractVar](scope.GetOwnDeclaredVar(name)), nil)) {
					return true
				}
			}
		}
		return false
	}
	{
		var scope_2 AbstractScope = this_.self_
		_ = scope_2
		for ; scope_2 != nil && !scope_2.IsJavaNil_(); scope_2 = j2gort.AsIface[AbstractScope](scope_2.GetParent()) {
			if !(j2gort.IfaceEquals[AbstractVar](j2gort.AsIface[AbstractVar](scope_2.GetOwnDeclaredVar(name)), nil)) || scope_2.HasOwnImplicitSlot(implicit) {
				return true
			}
		}
	}
	return false
}

// Java method: AbstractScope.getOwnImplicitSlot(com.google.javascript.jscomp.AbstractScope.@org.jspecify.annotations.Nullable ImplicitVar name)
func (this_ *AbstractScopeBase) GetOwnImplicitSlot(name AbstractScope_ImplicitVar) AbstractVar {
	_ = name
	if !(this_.self_.HasOwnImplicitSlot(name)) {
		return nil
	}
	var slotOrHolder any = this_.implicitVars
	_ = slotOrHolder
	if slotOrHolder == nil {
		var var_ AbstractVar = j2gort.AsIface[AbstractVar](this_.self_.MakeImplicitVar(name))
		_ = var_
		if var_ != nil && !var_.IsJavaNil_() {
			this_.implicitVars = var_
		}
		return var_
	}
	var holder *AbstractScope_ImplicitVars
	_ = holder
	var ok_5 bool
	holder, ok_5 = j2gort.TryCast[*AbstractScope_ImplicitVars](slotOrHolder)
	if ok_5 {
		var var__2 AbstractVar = j2gort.AsIface[AbstractVar](holder.Get(name))
		_ = var__2
		if var__2 == nil || var__2.IsJavaNil_() {
			var__2 = j2gort.AsIface[AbstractVar](this_.self_.MakeImplicitVar(name))
			if var__2 != nil && !var__2.IsJavaNil_() {
				holder.Put(name, var__2)
			}
		}
		return var__2
	}
	var singleVar AbstractVar = j2gort.AsIface[AbstractVar](slotOrHolder)
	_ = singleVar
	{
		var iifeCond_1 bool
		{

			sr_6 := singleVar.GetAbstractVarBase().GetName()
			se_7 := name.FieldGet_name()
			iifeCond_1 = j2gort.CheckNotNullString(sr_6) == se_7
		}
		if iifeCond_1 {
			return singleVar
		}
	}
	var var__3 AbstractVar = j2gort.AsIface[AbstractVar](this_.self_.MakeImplicitVar(name))
	_ = var__3
	if var__3 != nil && !var__3.IsJavaNil_() {
		var val_8 AbstractScope_ImplicitVars
		InitAbstractScope_ImplicitVars(&val_8)
		var holder_2 *AbstractScope_ImplicitVars = &val_8
		_ = holder_2
		holder_2.Put(j2gort.CheckNotNullEnum(AbstractScope_ImplicitVar_Of(singleVar.GetAbstractVarBase().GetName())), singleVar)
		holder_2.Put(name, var__3)
		this_.implicitVars = holder_2
	}
	return var__3
}

// Java method: AbstractScope.getOwnSlot(java.lang.String name)
func (this_ *AbstractScopeBase) GetOwnSlot(name string) any {
	_ = name
	var var_ AbstractVar = j2gort.AsIface[AbstractVar](this_.self_.GetOwnDeclaredVar(name))
	_ = var_
	if var_ != nil && !var_.IsJavaNil_() {
		return var_
	}
	return this_.GetOwnImplicitSlot(AbstractScope_ImplicitVar_Of(name))
}

// Java method: AbstractScope.getSlot(java.lang.String name)
func (this_ *AbstractScopeBase) GetSlot(name string) any {
	_ = name
	return j2gort.AsIface[AbstractVar](this_.self_.GetVar(name))
}

// Java method: AbstractScope.getVar(java.lang.String name)
func (this_ *AbstractScopeBase) GetVar(name string) any {
	_ = name
	var implicit AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(name)
	_ = implicit
	{
		var scope AbstractScope = this_.priv_AbstractScope_thisScope_491()
		_ = scope
		for ; scope != nil && !scope.IsJavaNil_(); scope = j2gort.AsIface[AbstractScope](scope.GetParent()) {
			var var_ AbstractVar = j2gort.AsIface[AbstractVar](scope.GetOwnDeclaredVar(name))
			_ = var_
			if var_ != nil && !var_.IsJavaNil_() {
				return var_
			}
			if implicit != 0 {
				var_ = scope.GetAbstractScopeBase().GetOwnImplicitSlot(implicit)
				if var_ != nil && !var_.IsJavaNil_() {
					return var_
				}
			}
		}
	}
	return nil
}

// Java method: AbstractScope.getArgumentsVar()
func (this_ *AbstractScopeBase) GetArgumentsVar() AbstractVar {
	{
		var scope AbstractScope = this_.priv_AbstractScope_thisScope_491()
		_ = scope
		for ; scope != nil && !scope.IsJavaNil_(); scope = j2gort.AsIface[AbstractScope](scope.GetParent()) {
			var arguments AbstractVar = scope.GetAbstractScopeBase().GetOwnImplicitSlot(StaticGet_AbstractScope_ImplicitVar_ARGUMENTS())
			_ = arguments
			if arguments != nil && !arguments.IsJavaNil_() {
				return arguments
			}
		}
	}
	return nil
}

// Java method: AbstractScope.canDeclare(java.lang.String name)
func (this_ *AbstractScopeBase) CanDeclare(name string) bool {
	_ = name
	return (!(this_.HasOwnSlot(name)) && ((!(this_.IsFunctionBlockScope()) || !(j2gort.AsIface[AbstractScope](this_.self_.GetParent()).GetAbstractScopeBase().HasOwnSlot(name))) || this_.priv_AbstractScope_isBleedingFunctionName_468(name)))
}

// Java method: AbstractScope.isBleedingFunctionName(java.lang.String name)
func (this_ *AbstractScopeBase) priv_AbstractScope_isBleedingFunctionName_468(name string) bool {
	_ = name
	var var_ AbstractVar = j2gort.AsIface[AbstractVar](this_.self_.GetVar(name))
	_ = var_
	return (((var_ != nil && !var_.IsJavaNil_()) && (var_.GetAbstractVarBase().GetNode() != 0)) && var_.GetAbstractVarBase().GetNode().GetParent().IsFunction())
}

// Java method: AbstractScope.getVarIterable()
func (this_ *AbstractScopeBase) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.putOwnVarsIfAbsent(java.util.LinkedHashMap<java.lang.String,V> dst)
func (this_ *AbstractScopeBase) PutOwnVarsIfAbsent(dst *j2gort.Ext_java_util_LinkedHashMap) {
	_ = dst
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.getAllAccessibleVariables()
func (this_ *AbstractScopeBase) GetAllAccessibleVariables() j2gort.Ext_java_lang_Iterable {
	var accessibleVars *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap_1_int32(int32(16))
	_ = accessibleVars
	var s AbstractScope = this_.self_
	_ = s
	for s != nil && !s.IsJavaNil_() {
		s.PutOwnVarsIfAbsent(accessibleVars)
		s = j2gort.AsIface[AbstractScope](s.GetParent())
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](accessibleVars.Values())
}

// Java method: AbstractScope.getAllSymbols()
func (this_ *AbstractScopeBase) GetAllSymbols() j2gort.Ext_java_lang_Iterable {
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.addVarsTo(java.util.ArrayList<V> dst)
func (this_ *AbstractScopeBase) AddVarsTo(dst j2gort.Ext_java_util_ArrayList) {
	_ = dst
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.addVarNamesTo(java.util.LinkedHashSet<java.lang.String> dst)
func (this_ *AbstractScopeBase) AddVarNamesTo(dst *j2gort.Ext_java_util_LinkedHashSet) {
	_ = dst
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.getVarCount()
func (this_ *AbstractScopeBase) GetVarCount() int32 {
	panic("abstract or native method invoked")
}

// Java method: AbstractScope.isGlobal()
func (this_ *AbstractScopeBase) IsGlobal() bool {
	return j2gort.IfaceEquals[AbstractScope](j2gort.AsIface[AbstractScope](this_.self_.GetParent()), nil)
}

// Java method: AbstractScope.isLocal()
func (this_ *AbstractScopeBase) IsLocal() bool {
	return !(j2gort.IfaceEquals[AbstractScope](j2gort.AsIface[AbstractScope](this_.self_.GetParent()), nil))
}

// Java method: AbstractScope.isBlockScope()
func (this_ *AbstractScopeBase) IsBlockScope() bool {
	return NodeUtil_CreatesBlockScope(this_.rootNode)
}

// Java method: AbstractScope.isStaticBlockScope()
func (this_ *AbstractScopeBase) IsStaticBlockScope() bool {
	return NodeUtil_IsClassStaticBlock(this_.GetRootNode_0())
}

// Java method: AbstractScope.isFunctionBlockScope()
func (this_ *AbstractScopeBase) IsFunctionBlockScope() bool {
	return NodeUtil_IsFunctionBlock(this_.GetRootNode_0())
}

// Java method: AbstractScope.isFunctionScope()
func (this_ *AbstractScopeBase) IsFunctionScope() bool {
	return this_.GetRootNode_0().IsFunction()
}

// Java method: AbstractScope.isModuleScope()
func (this_ *AbstractScopeBase) IsModuleScope() bool {
	return this_.GetRootNode_0().IsModuleBody()
}

// Java method: AbstractScope.isMemberFieldDefScope()
func (this_ *AbstractScopeBase) IsMemberFieldDefScope() bool {
	return this_.GetRootNode_0().IsMemberFieldDef()
}

// Java method: AbstractScope.isComputedFieldDefRhsScope()
func (this_ *AbstractScopeBase) IsComputedFieldDefRhsScope() bool {
	return this_.GetRootNode_0().IsComputedFieldDef()
}

// Java method: AbstractScope.isCatchScope()
func (this_ *AbstractScopeBase) IsCatchScope() bool {
	return ((this_.GetRootNode_0().IsBlock() && this_.GetRootNode_0().HasOneChild()) && this_.GetRootNode_0().GetFirstChild().IsCatch())
}

// Java method: AbstractScope.isCfgRootScope()
func (this_ *AbstractScopeBase) IsCfgRootScope() bool {
	return NodeUtil_IsValidCfgRoot(this_.rootNode)
}

// Java method: AbstractScope.isHoistScope()
func (this_ *AbstractScopeBase) IsHoistScope() bool {
	return ((((this_.IsFunctionScope() || this_.IsFunctionBlockScope()) || this_.self_.IsGlobal()) || this_.IsModuleScope()) || this_.IsStaticBlockScope())
}

// Java method: AbstractScope.getClosestHoistScope()
func (this_ *AbstractScopeBase) GetClosestHoistScope() AbstractScope {
	var current AbstractScope = this_.priv_AbstractScope_thisScope_491()
	_ = current
	for current != nil && !current.IsJavaNil_() {
		if current.GetAbstractScopeBase().IsHoistScope() {
			return current
		}
		current = j2gort.AsIface[AbstractScope](current.GetParent())
	}
	return nil
}

// Java method: AbstractScope.getClosestCfgRootScope()
func (this_ *AbstractScopeBase) GetClosestCfgRootScope() AbstractScope {
	var current AbstractScope = this_.priv_AbstractScope_thisScope_491()
	_ = current
	for !(current.GetAbstractScopeBase().IsCfgRootScope()) {
		current = j2gort.AsIface[AbstractScope](current.GetParent())
	}
	return current
}

// Java method: AbstractScope.getClosestContainerScope()
func (this_ *AbstractScopeBase) GetClosestContainerScope() AbstractScope {
	var scope AbstractScope = this_.GetClosestHoistScope()
	_ = scope
	if scope.GetAbstractScopeBase().IsFunctionBlockScope() {
		scope = j2gort.AsIface[AbstractScope](scope.GetParent())
		if !(!(scope.GetAbstractScopeBase().IsBlockScope())) {
			j2gort.CheckStateMsg(false, scope)
		}
	}
	return scope
}

// Java method: AbstractScope.thisScope()
func (this_ *AbstractScopeBase) priv_AbstractScope_thisScope_491() AbstractScope {
	return this_.self_
}

// Java method: AbstractScope.checkChildScope(S parent)
func (this_ *AbstractScopeBase) CheckChildScope(parent AbstractScope) {
	_ = parent
	j2gort.CheckNotNullIface(parent)
	if !(NodeUtil_CreatesScope(this_.rootNode)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(this_.rootNode))
	}
	if !(this_.rootNode != parent.GetAbstractScopeBase().GetRootNode_0()) {
		j2gort.CheckArgumentFmt(false, "rootNode should not be the parent's root node: %s", j2gort.BoxNode(this_.rootNode))
	}
}

// Java method: AbstractScope.checkRootScope()
func (this_ *AbstractScopeBase) CheckRootScope() {
	if !((NodeUtil_CreatesScope(this_.rootNode) || this_.rootNode.IsScript()) || this_.rootNode.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(this_.rootNode))
	}
}

// Java method: AbstractScope.getCommonParent(S other)
func (this_ *AbstractScopeBase) GetCommonParent(other AbstractScope) AbstractScope {
	_ = other
	var left AbstractScope = this_.priv_AbstractScope_thisScope_491()
	_ = left
	var right AbstractScope = other
	_ = right
	for ((left != nil && !left.IsJavaNil_()) && (right != nil && !right.IsJavaNil_())) && !(left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_()))) {
		var leftDepth int32 = left.GetDepth()
		_ = leftDepth
		var rightDepth int32 = right.GetDepth()
		_ = rightDepth
		if leftDepth >= rightDepth {
			left = j2gort.AsIface[AbstractScope](left.GetParent())
		}
		if leftDepth <= rightDepth {
			right = j2gort.AsIface[AbstractScope](right.GetParent())
		}
	}
	if !((left != nil && !left.IsJavaNil_()) && (left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())))) {
		j2gort.CheckState(false)
	}
	return left
}

// Java method: AbstractScope.hasSameContainerScope(S other)
func (this_ *AbstractScopeBase) HasSameContainerScope(other AbstractScope) bool {
	_ = other
	return ((this_.self_ == other || ((this_.self_ == nil || this_.self_.IsJavaNil_()) && (other == nil || other.IsJavaNil_()))) || j2gort.IfaceEquals[AbstractScope](this_.GetClosestContainerScope(), other.GetAbstractScopeBase().GetClosestContainerScope()))
}

// Java method: AbstractScope.getScopeOfThis()
func (this_ *AbstractScopeBase) GetScopeOfThis() AbstractScope {
	var scope AbstractScope = this_.GetClosestContainerScope()
	_ = scope
	for !(scope.IsGlobal()) && !(NodeUtil_IsNonArrowFunction(scope.GetAbstractScopeBase().GetRootNode_0())) {
		scope = j2gort.AsIface[AbstractScope](scope.GetParent()).GetAbstractScopeBase().GetClosestContainerScope()
	}
	return scope
}

func (this_ *AbstractScopeBase) GetTopmostScopeOfEventualDeclaration(p0_ string) pkg_rhino.StaticScope {
	return pkg_rhino.IfaceDefault_StaticScope_GetTopmostScopeOfEventualDeclaration(this_.self_, p0_)
}

// Java class: com.google.javascript.jscomp.AbstractScope.ImplicitVars
type AbstractScope_ImplicitVars struct {
	argumentsVar any `gson:"argumentsVar"`
	exportsVar   any `gson:"exportsVar"`
	superVar     any `gson:"superVar"`
	thisVar      any `gson:"thisVar"`
	idHash_      atomic.Int32
}

func (this_ *AbstractScope_ImplicitVars) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractScope_ImplicitVars) IsJavaNil_() bool { return this_ == nil }
func (*AbstractScope_ImplicitVars) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractScope$ImplicitVars", "ImplicitVars"
}
func (this_ *AbstractScope_ImplicitVars) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AbstractScope_ImplicitVars) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractScope.ImplicitVars", "com.google.javascript.jscomp.AbstractScope$ImplicitVars", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractScope_ImplicitVars](v_)
	})
}

func (this_ *AbstractScope_ImplicitVars) initFields_AbstractScope_ImplicitVars() {
	_ = this_
}

// Java constructor: ImplicitVars()
func NewAbstractScope_ImplicitVars() *AbstractScope_ImplicitVars {
	this_ := &AbstractScope_ImplicitVars{}
	InitAbstractScope_ImplicitVars(this_)
	return this_
}

func InitAbstractScope_ImplicitVars(this_ *AbstractScope_ImplicitVars) {
	this_.initFields_AbstractScope_ImplicitVars()
}

// Java method: ImplicitVars.get(com.google.javascript.jscomp.AbstractScope.ImplicitVar kind)
func (this_ *AbstractScope_ImplicitVars) Get(kind AbstractScope_ImplicitVar) any {
	_ = kind
	switch (kind).Ordinal() {
	case int32(0) /* ImplicitVar.ARGUMENTS */ :
		return this_.argumentsVar
	case int32(1) /* ImplicitVar.EXPORTS */ :
		return this_.exportsVar
	case int32(2) /* ImplicitVar.SUPER */ :
		return this_.superVar
	case int32(3) /* ImplicitVar.THIS */ :
		return this_.thisVar
	}
	panic("unreachable switch expression")
}

// Java method: ImplicitVars.put(com.google.javascript.jscomp.AbstractScope.ImplicitVar kind, java.lang.Object var)
func (this_ *AbstractScope_ImplicitVars) Put(kind AbstractScope_ImplicitVar, var_ any) {
	_ = kind
	_ = var_
	switch (kind).Ordinal() {
	case int32(0) /* ImplicitVar.ARGUMENTS */ :
		this_.argumentsVar = var_
	case int32(1) /* ImplicitVar.EXPORTS */ :
		this_.exportsVar = var_
	case int32(2) /* ImplicitVar.SUPER */ :
		this_.superVar = var_
	case int32(3) /* ImplicitVar.THIS */ :
		this_.thisVar = var_
	}
}

// Java enum: com.google.javascript.jscomp.AbstractScope.ImplicitVar
type AbstractScope_ImplicitVar uint8

const (
	enum_AbstractScope_ImplicitVar_ARGUMENTS AbstractScope_ImplicitVar = 1
	enum_AbstractScope_ImplicitVar_EXPORTS   AbstractScope_ImplicitVar = 2
	enum_AbstractScope_ImplicitVar_SUPER     AbstractScope_ImplicitVar = 3
	enum_AbstractScope_ImplicitVar_THIS      AbstractScope_ImplicitVar = 4
)

var enumBaseTable_AbstractScope_ImplicitVar = [...]j2gort.EnumBase{
	{Name_: "ARGUMENTS", Ordinal_: 0},
	{Name_: "EXPORTS", Ordinal_: 1},
	{Name_: "SUPER", Ordinal_: 2},
	{Name_: "THIS", Ordinal_: 3},
}
var enumValues_AbstractScope_ImplicitVar = []AbstractScope_ImplicitVar{1, 2, 3, 4}

var enumField_AbstractScope_ImplicitVar_name = [...]string{
	"arguments",
	"exports",
	"super",
	"this",
}

func (this_ AbstractScope_ImplicitVar) FieldGet_name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_AbstractScope_ImplicitVar_name[this_-1]
}

func StaticGet_AbstractScope_ImplicitVar_ARGUMENTS() AbstractScope_ImplicitVar {
	return enum_AbstractScope_ImplicitVar_ARGUMENTS
}
func StaticGet_AbstractScope_ImplicitVar_EXPORTS() AbstractScope_ImplicitVar {
	return enum_AbstractScope_ImplicitVar_EXPORTS
}
func StaticGet_AbstractScope_ImplicitVar_SUPER() AbstractScope_ImplicitVar {
	return enum_AbstractScope_ImplicitVar_SUPER
}
func StaticGet_AbstractScope_ImplicitVar_THIS() AbstractScope_ImplicitVar {
	return enum_AbstractScope_ImplicitVar_THIS
}

func AbstractScope_ImplicitVar_Values() []AbstractScope_ImplicitVar {
	return j2gort.CloneSlice(enumValues_AbstractScope_ImplicitVar)
}

func AbstractScope_ImplicitVar_ValueOf(name string) AbstractScope_ImplicitVar {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AbstractScope_ImplicitVar {
		if enumBaseTable_AbstractScope_ImplicitVar[i_].Name_ == name {
			return AbstractScope_ImplicitVar(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AbstractScope.ImplicitVar." + name))
}

// Java method: ImplicitVar.isMadeByScope(com.google.javascript.jscomp.AbstractScope<?,?> scope)
func (this_ AbstractScope_ImplicitVar) IsMadeByScope(scope AbstractScope) bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = scope
	switch (this_).Ordinal() {
	case int32(1) /* ImplicitVar.EXPORTS */ :
		return (scope.GetAbstractScopeBase().IsModuleScope() && scope.GetAbstractScopeBase().GetRootNode_0().GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE()))
	case int32(2) /* ImplicitVar.SUPER */, int32(3) /* ImplicitVar.THIS */ :
		return (((scope.GetAbstractScopeBase().IsStaticBlockScope() || NodeUtil_IsNonArrowFunction(scope.GetAbstractScopeBase().GetRootNode_0())) || scope.GetAbstractScopeBase().IsMemberFieldDefScope()) || scope.GetAbstractScopeBase().IsComputedFieldDefRhsScope())
	case int32(0) /* ImplicitVar.ARGUMENTS */ :
		return NodeUtil_IsNonArrowFunction(scope.GetAbstractScopeBase().GetRootNode_0())
	}
	panic("unreachable switch expression")
}

// Java static method: ImplicitVar.of(java.lang.String name)
func AbstractScope_ImplicitVar_Of(name string) AbstractScope_ImplicitVar {
	_ = name
	switch j2gort.CheckNotNullString(name) {
	case "arguments":
		return StaticGet_AbstractScope_ImplicitVar_ARGUMENTS()
	case "super":
		return StaticGet_AbstractScope_ImplicitVar_SUPER()
	case "this":
		return StaticGet_AbstractScope_ImplicitVar_THIS()
	case "exports":
		return StaticGet_AbstractScope_ImplicitVar_EXPORTS()
	default:
		return AbstractScope_ImplicitVar(0)
	}
	panic("unreachable switch expression")
}

func (this_ AbstractScope_ImplicitVar) IsJavaNil_() bool                              { return this_ == 0 }
func (AbstractScope_ImplicitVar) HasIdentityEquals_()                                 {}
func (AbstractScope_ImplicitVar) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AbstractScope_ImplicitVar) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AbstractScope_ImplicitVar) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AbstractScope_ImplicitVar) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractScope$ImplicitVar", "ImplicitVar"
}
func (this_ AbstractScope_ImplicitVar) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AbstractScope_ImplicitVar[this_-1]
}
func (this_ AbstractScope_ImplicitVar) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AbstractScope_ImplicitVar) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AbstractScope_ImplicitVar) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractScope_ImplicitVar[this_-1].Name_
}
func (this_ AbstractScope_ImplicitVar) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AbstractScope_ImplicitVar[this_-1].Name_
}
func (this_ AbstractScope_ImplicitVar) Equals(other_ any) bool {
	o_, ok_ := other_.(AbstractScope_ImplicitVar)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AbstractScope_ImplicitVar) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AbstractScope_ImplicitVar[this_-1].HashCode()
}
func (this_ AbstractScope_ImplicitVar) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractScope_ImplicitVar[this_-1].CompareTo(other_)
}
func (this_ AbstractScope_ImplicitVar) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AbstractScope_ImplicitVar[this_-1].Name_)
}
func init() {
	vals_ := []any{AbstractScope_ImplicitVar(1), AbstractScope_ImplicitVar(2), AbstractScope_ImplicitVar(3), AbstractScope_ImplicitVar(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractScope.ImplicitVar", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractScope$ImplicitVar", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractScope.ImplicitVar", "com.google.javascript.jscomp.AbstractScope$ImplicitVar", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractScope_ImplicitVar](v_)
	})
}
