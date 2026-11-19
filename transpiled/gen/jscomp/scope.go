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
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Scope
type Scope struct {
	AbstractScopeBase
	parent   *Scope                           `gson:"parent"`
	var0     *Var                             `gson:"var0"`
	var1     *Var                             `gson:"var1"`
	var2     *Var                             `gson:"var2"`
	var3     *Var                             `gson:"var3"`
	var4     *Var                             `gson:"var4"`
	var5     *Var                             `gson:"var5"`
	vars     *j2gort.OrderedMap[string, *Var] `gson:"vars"`
	depth    int32                            `gson:"depth"`
	varSlab_ *j2gort.PooledSlab[Var]
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

var static_Scope_EMPTY_VAR_NAMES []string

func StaticGet_Scope_EMPTY_VAR_NAMES() []string {
	clinit_Scope()
	return static_Scope_EMPTY_VAR_NAMES
}
func StaticSet_Scope_EMPTY_VAR_NAMES(v []string) {
	clinit_Scope()
	static_Scope_EMPTY_VAR_NAMES = v
}

var static_Scope_SINGLE_ASCII_VAR_NAMES [][]string

func StaticGet_Scope_SINGLE_ASCII_VAR_NAMES() [][]string {
	clinit_Scope()
	return static_Scope_SINGLE_ASCII_VAR_NAMES
}
func StaticSet_Scope_SINGLE_ASCII_VAR_NAMES(v [][]string) {
	clinit_Scope()
	static_Scope_SINGLE_ASCII_VAR_NAMES = v
}

var static_Scope_SINGLE_EXPORTS []string

func StaticGet_Scope_SINGLE_EXPORTS() []string {
	clinit_Scope()
	return static_Scope_SINGLE_EXPORTS
}
func StaticSet_Scope_SINGLE_EXPORTS(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_EXPORTS = v
}

var static_Scope_SINGLE_THIS []string

func StaticGet_Scope_SINGLE_THIS() []string {
	clinit_Scope()
	return static_Scope_SINGLE_THIS
}
func StaticSet_Scope_SINGLE_THIS(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_THIS = v
}

var static_Scope_SINGLE_ARGUMENTS []string

func StaticGet_Scope_SINGLE_ARGUMENTS() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ARGUMENTS
}
func StaticSet_Scope_SINGLE_ARGUMENTS(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ARGUMENTS = v
}

var static_Scope_SINGLE_VALUE []string

func StaticGet_Scope_SINGLE_VALUE() []string {
	clinit_Scope()
	return static_Scope_SINGLE_VALUE
}
func StaticSet_Scope_SINGLE_VALUE(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_VALUE = v
}

var static_Scope_SINGLE_EVENT []string

func StaticGet_Scope_SINGLE_EVENT() []string {
	clinit_Scope()
	return static_Scope_SINGLE_EVENT
}
func StaticSet_Scope_SINGLE_EVENT(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_EVENT = v
}

var static_Scope_SINGLE_ERROR []string

func StaticGet_Scope_SINGLE_ERROR() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ERROR
}
func StaticSet_Scope_SINGLE_ERROR(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ERROR = v
}

var static_Scope_SINGLE_INDEX []string

func StaticGet_Scope_SINGLE_INDEX() []string {
	clinit_Scope()
	return static_Scope_SINGLE_INDEX
}
func StaticSet_Scope_SINGLE_INDEX(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_INDEX = v
}

var static_Scope_SINGLE_OPTIONS []string

func StaticGet_Scope_SINGLE_OPTIONS() []string {
	clinit_Scope()
	return static_Scope_SINGLE_OPTIONS
}
func StaticSet_Scope_SINGLE_OPTIONS(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_OPTIONS = v
}

var static_Scope_SINGLE_ELEMENT []string

func StaticGet_Scope_SINGLE_ELEMENT() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ELEMENT
}
func StaticSet_Scope_SINGLE_ELEMENT(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ELEMENT = v
}

var static_Scope_SINGLE_RESULT []string

func StaticGet_Scope_SINGLE_RESULT() []string {
	clinit_Scope()
	return static_Scope_SINGLE_RESULT
}
func StaticSet_Scope_SINGLE_RESULT(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_RESULT = v
}

var static_Scope_SINGLE_TARGET []string

func StaticGet_Scope_SINGLE_TARGET() []string {
	clinit_Scope()
	return static_Scope_SINGLE_TARGET
}
func StaticSet_Scope_SINGLE_TARGET(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_TARGET = v
}

var static_Scope_SINGLE_STATE []string

func StaticGet_Scope_SINGLE_STATE() []string {
	clinit_Scope()
	return static_Scope_SINGLE_STATE
}
func StaticSet_Scope_SINGLE_STATE(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_STATE = v
}

var static_Scope_SINGLE_PROPS []string

func StaticGet_Scope_SINGLE_PROPS() []string {
	clinit_Scope()
	return static_Scope_SINGLE_PROPS
}
func StaticSet_Scope_SINGLE_PROPS(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_PROPS = v
}

var static_Scope_SINGLE_DATA []string

func StaticGet_Scope_SINGLE_DATA() []string {
	clinit_Scope()
	return static_Scope_SINGLE_DATA
}
func StaticSet_Scope_SINGLE_DATA(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_DATA = v
}

var static_Scope_SINGLE_ITEM []string

func StaticGet_Scope_SINGLE_ITEM() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ITEM
}
func StaticSet_Scope_SINGLE_ITEM(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ITEM = v
}

var static_Scope_SINGLE_NAME []string

func StaticGet_Scope_SINGLE_NAME() []string {
	clinit_Scope()
	return static_Scope_SINGLE_NAME
}
func StaticSet_Scope_SINGLE_NAME(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_NAME = v
}

var static_Scope_SINGLE_TYPE []string

func StaticGet_Scope_SINGLE_TYPE() []string {
	clinit_Scope()
	return static_Scope_SINGLE_TYPE
}
func StaticSet_Scope_SINGLE_TYPE(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_TYPE = v
}

var static_Scope_SINGLE_NODE []string

func StaticGet_Scope_SINGLE_NODE() []string {
	clinit_Scope()
	return static_Scope_SINGLE_NODE
}
func StaticSet_Scope_SINGLE_NODE(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_NODE = v
}

var static_Scope_SINGLE_KEY []string

func StaticGet_Scope_SINGLE_KEY() []string {
	clinit_Scope()
	return static_Scope_SINGLE_KEY
}
func StaticSet_Scope_SINGLE_KEY(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_KEY = v
}

var static_Scope_SINGLE_VAL []string

func StaticGet_Scope_SINGLE_VAL() []string {
	clinit_Scope()
	return static_Scope_SINGLE_VAL
}
func StaticSet_Scope_SINGLE_VAL(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_VAL = v
}

var static_Scope_SINGLE_ERR []string

func StaticGet_Scope_SINGLE_ERR() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ERR
}
func StaticSet_Scope_SINGLE_ERR(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ERR = v
}

var static_Scope_SINGLE_RES []string

func StaticGet_Scope_SINGLE_RES() []string {
	clinit_Scope()
	return static_Scope_SINGLE_RES
}
func StaticSet_Scope_SINGLE_RES(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_RES = v
}

var static_Scope_SINGLE_REQ []string

func StaticGet_Scope_SINGLE_REQ() []string {
	clinit_Scope()
	return static_Scope_SINGLE_REQ
}
func StaticSet_Scope_SINGLE_REQ(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_REQ = v
}

var static_Scope_SINGLE_CB []string

func StaticGet_Scope_SINGLE_CB() []string {
	clinit_Scope()
	return static_Scope_SINGLE_CB
}
func StaticSet_Scope_SINGLE_CB(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_CB = v
}

var static_Scope_SINGLE_FN []string

func StaticGet_Scope_SINGLE_FN() []string {
	clinit_Scope()
	return static_Scope_SINGLE_FN
}
func StaticSet_Scope_SINGLE_FN(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_FN = v
}

var static_Scope_SINGLE_EL []string

func StaticGet_Scope_SINGLE_EL() []string {
	clinit_Scope()
	return static_Scope_SINGLE_EL
}
func StaticSet_Scope_SINGLE_EL(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_EL = v
}

var static_Scope_SINGLE_ID []string

func StaticGet_Scope_SINGLE_ID() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ID
}
func StaticSet_Scope_SINGLE_ID(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ID = v
}

var static_Scope_SINGLE_SELF []string

func StaticGet_Scope_SINGLE_SELF() []string {
	clinit_Scope()
	return static_Scope_SINGLE_SELF
}
func StaticSet_Scope_SINGLE_SELF(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_SELF = v
}

var static_Scope_SINGLE_THAT []string

func StaticGet_Scope_SINGLE_THAT() []string {
	clinit_Scope()
	return static_Scope_SINGLE_THAT
}
func StaticSet_Scope_SINGLE_THAT(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_THAT = v
}

var static_Scope_SINGLE_ARGS []string

func StaticGet_Scope_SINGLE_ARGS() []string {
	clinit_Scope()
	return static_Scope_SINGLE_ARGS
}
func StaticSet_Scope_SINGLE_ARGS(v []string) {
	clinit_Scope()
	static_Scope_SINGLE_ARGS = v
}

var static_Scope_PAIR_A_B []string

func StaticGet_Scope_PAIR_A_B() []string {
	clinit_Scope()
	return static_Scope_PAIR_A_B
}
func StaticSet_Scope_PAIR_A_B(v []string) {
	clinit_Scope()
	static_Scope_PAIR_A_B = v
}

var static_Scope_PAIR_X_Y []string

func StaticGet_Scope_PAIR_X_Y() []string {
	clinit_Scope()
	return static_Scope_PAIR_X_Y
}
func StaticSet_Scope_PAIR_X_Y(v []string) {
	clinit_Scope()
	static_Scope_PAIR_X_Y = v
}

var static_Scope_PAIR_I_N []string

func StaticGet_Scope_PAIR_I_N() []string {
	clinit_Scope()
	return static_Scope_PAIR_I_N
}
func StaticSet_Scope_PAIR_I_N(v []string) {
	clinit_Scope()
	static_Scope_PAIR_I_N = v
}

var static_Scope_PAIR_K_V []string

func StaticGet_Scope_PAIR_K_V() []string {
	clinit_Scope()
	return static_Scope_PAIR_K_V
}
func StaticSet_Scope_PAIR_K_V(v []string) {
	clinit_Scope()
	static_Scope_PAIR_K_V = v
}

var static_Scope_PAIR_E_T []string

func StaticGet_Scope_PAIR_E_T() []string {
	clinit_Scope()
	return static_Scope_PAIR_E_T
}
func StaticSet_Scope_PAIR_E_T(v []string) {
	clinit_Scope()
	static_Scope_PAIR_E_T = v
}

var static_Scope_PAIR_T_E []string

func StaticGet_Scope_PAIR_T_E() []string {
	clinit_Scope()
	return static_Scope_PAIR_T_E
}
func StaticSet_Scope_PAIR_T_E(v []string) {
	clinit_Scope()
	static_Scope_PAIR_T_E = v
}

var static_Scope_PAIR_KEY_VALUE []string

func StaticGet_Scope_PAIR_KEY_VALUE() []string {
	clinit_Scope()
	return static_Scope_PAIR_KEY_VALUE
}
func StaticSet_Scope_PAIR_KEY_VALUE(v []string) {
	clinit_Scope()
	static_Scope_PAIR_KEY_VALUE = v
}

var static_Scope_PAIR_RESOLVE_REJECT []string

func StaticGet_Scope_PAIR_RESOLVE_REJECT() []string {
	clinit_Scope()
	return static_Scope_PAIR_RESOLVE_REJECT
}
func StaticSet_Scope_PAIR_RESOLVE_REJECT(v []string) {
	clinit_Scope()
	static_Scope_PAIR_RESOLVE_REJECT = v
}

var clinitDone_Scope uint32
var clinitState_Scope j2gort.ClinitState

func clinit_Scope() {
	if atomic.LoadUint32(&clinitDone_Scope) != 0 {
		return
	}
	clinitSlow_Scope()
}

func clinitSlow_Scope() {
	if !j2gort.EnterClinit(&clinitState_Scope, &clinitDone_Scope) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Scope, &clinitDone_Scope)
	static_Scope_EMPTY_VAR_NAMES = j2gort.NewStringArray(int32(0))
	static_Scope_SINGLE_ASCII_VAR_NAMES = make([][]string, int32(128))
	static_Scope_SINGLE_EXPORTS = []string{"exports"}
	static_Scope_SINGLE_THIS = []string{"this"}
	static_Scope_SINGLE_ARGUMENTS = []string{"arguments"}
	static_Scope_SINGLE_VALUE = []string{"value"}
	static_Scope_SINGLE_EVENT = []string{"event"}
	static_Scope_SINGLE_ERROR = []string{"error"}
	static_Scope_SINGLE_INDEX = []string{"index"}
	static_Scope_SINGLE_OPTIONS = []string{"options"}
	static_Scope_SINGLE_ELEMENT = []string{"element"}
	static_Scope_SINGLE_RESULT = []string{"result"}
	static_Scope_SINGLE_TARGET = []string{"target"}
	static_Scope_SINGLE_STATE = []string{"state"}
	static_Scope_SINGLE_PROPS = []string{"props"}
	static_Scope_SINGLE_DATA = []string{"data"}
	static_Scope_SINGLE_ITEM = []string{"item"}
	static_Scope_SINGLE_NAME = []string{"name"}
	static_Scope_SINGLE_TYPE = []string{"type"}
	static_Scope_SINGLE_NODE = []string{"node"}
	static_Scope_SINGLE_KEY = []string{"key"}
	static_Scope_SINGLE_VAL = []string{"val"}
	static_Scope_SINGLE_ERR = []string{"err"}
	static_Scope_SINGLE_RES = []string{"res"}
	static_Scope_SINGLE_REQ = []string{"req"}
	static_Scope_SINGLE_CB = []string{"cb"}
	static_Scope_SINGLE_FN = []string{"fn"}
	static_Scope_SINGLE_EL = []string{"el"}
	static_Scope_SINGLE_ID = []string{"id"}
	static_Scope_SINGLE_SELF = []string{"self"}
	static_Scope_SINGLE_THAT = []string{"that"}
	static_Scope_SINGLE_ARGS = []string{"args"}
	static_Scope_PAIR_A_B = []string{"a", "b"}
	static_Scope_PAIR_X_Y = []string{"x", "y"}
	static_Scope_PAIR_I_N = []string{"i", "n"}
	static_Scope_PAIR_K_V = []string{"k", "v"}
	static_Scope_PAIR_E_T = []string{"e", "t"}
	static_Scope_PAIR_T_E = []string{"t", "e"}
	static_Scope_PAIR_KEY_VALUE = []string{"key", "value"}
	static_Scope_PAIR_RESOLVE_REJECT = []string{"resolve", "reject"}
	{
		var c int32 = int32(uint16(48))
		_ = c
		for ; c <= int32(uint16(57)); c = c + 1 {
			{
				a_1 := StaticGet_Scope_SINGLE_ASCII_VAR_NAMES()
				i_2 := c
				a_1[i_2] = []string{javalang.StringFromChar(uint16(c))}
			}
		}
	}
	{
		var c_2 int32 = int32(uint16(65))
		_ = c_2
		for ; c_2 <= int32(uint16(90)); c_2 = c_2 + 1 {
			{
				a_3 := StaticGet_Scope_SINGLE_ASCII_VAR_NAMES()
				i_4 := c_2
				a_3[i_4] = []string{javalang.StringFromChar(uint16(c_2))}
			}
		}
	}
	{
		var c_3 int32 = int32(uint16(97))
		_ = c_3
		for ; c_3 <= int32(uint16(122)); c_3 = c_3 + 1 {
			{
				a_5 := StaticGet_Scope_SINGLE_ASCII_VAR_NAMES()
				i_6 := c_3
				a_5[i_6] = []string{javalang.StringFromChar(uint16(c_3))}
			}
		}
	}
	{
		a_7 := StaticGet_Scope_SINGLE_ASCII_VAR_NAMES()
		i_8 := int32(uint16(95))
		a_7[i_8] = []string{"_"}
	}
	{
		a_9 := StaticGet_Scope_SINGLE_ASCII_VAR_NAMES()
		i_10 := int32(uint16(36))
		a_9[i_10] = []string{"$"}
	}
}

func (this_ *Scope) initFields_Scope() {
	_ = this_
}

// Java constructor: Scope(com.google.javascript.jscomp.Scope parent, com.google.javascript.rhino.Node rootNode)
func NewScope_2(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	clinit_Scope()
	this_ := &Scope{}
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
	clinit_Scope()
	this_ := &Scope{}
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
	clinit_Scope()
	_ = rootNode
	return NewScope_1(rootNode)
}

// Java static method: Scope.createChildScope(com.google.javascript.jscomp.Scope parent, com.google.javascript.rhino.Node rootNode)
func Scope_CreateChildScope(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	clinit_Scope()
	_ = parent
	_ = rootNode
	return NewScope_2(parent, rootNode)
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
func (this_ *Scope) priv_Scope_findOwnDeclaredVar_5009(name string) *Var {
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
	return j2gort.PtrToAny(this_.priv_Scope_findOwnDeclaredVar_5009(name))
}

// Java method: Scope.putOwnDeclaredVar(java.lang.String name, com.google.javascript.jscomp.Var var)
func (this_ *Scope) PutOwnDeclaredVar(name string, arg_11 any) {
	_ = name
	var_ := j2gort.CastPtr[Var](arg_11)
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
	this_.vars = j2gort.NewOrderedMapPtrCap[string, *Var](func() int32 {
		if this_.parent == nil {
			return int32(1024)
		}
		return int32(16)
	}())
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
func (this_ *Scope) UndeclareInteral(arg_12 any) {
	var_ := j2gort.CastPtr[Var](arg_12)
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
				var var_ *Var = scope.priv_Scope_findOwnDeclaredVar_5009(name)
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
			var var__2 *Var = scope_2.priv_Scope_findOwnDeclaredVar_5009(name)
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
			cur_13 := this_.vars.IterValues()
			for cur_13.HasNext() {

				v := cur_13.
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
			cur_15 := this_.vars.IterValues()
			for cur_15.HasNext() {

				v := cur_15.
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
			cur_17 := this_.vars.IterKeys()
			for cur_17.HasNext() {
				name := cur_17.
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

// Java static method: Scope.singletonVarNameArray(java.lang.String name)
func Scope_SingletonVarNameArray(name string) []string {
	clinit_Scope()
	_ = name
	var len_ int32 = j2gort.JStrLen(name)
	_ = len_
	if len_ == int32(1) {
		var c uint16 = j2gort.JStrCharAt(name, int32(0))
		_ = c
		if int32(c) < int32(128) {
			var cached []string = StaticGet_Scope_SINGLE_ASCII_VAR_NAMES()[int32(c)]
			_ = cached
			if cached != nil {
				return cached
			}
		}
		return []string{name}
	}
	switch j2gort.CheckNotNullString(name) {
	case "exports":
		return StaticGet_Scope_SINGLE_EXPORTS()
	case "this":
		return StaticGet_Scope_SINGLE_THIS()
	case "arguments":
		return StaticGet_Scope_SINGLE_ARGUMENTS()
	case "value":
		return StaticGet_Scope_SINGLE_VALUE()
	case "event":
		return StaticGet_Scope_SINGLE_EVENT()
	case "error":
		return StaticGet_Scope_SINGLE_ERROR()
	case "index":
		return StaticGet_Scope_SINGLE_INDEX()
	case "options":
		return StaticGet_Scope_SINGLE_OPTIONS()
	case "element":
		return StaticGet_Scope_SINGLE_ELEMENT()
	case "result":
		return StaticGet_Scope_SINGLE_RESULT()
	case "target":
		return StaticGet_Scope_SINGLE_TARGET()
	case "state":
		return StaticGet_Scope_SINGLE_STATE()
	case "props":
		return StaticGet_Scope_SINGLE_PROPS()
	case "data":
		return StaticGet_Scope_SINGLE_DATA()
	case "item":
		return StaticGet_Scope_SINGLE_ITEM()
	case "name":
		return StaticGet_Scope_SINGLE_NAME()
	case "type":
		return StaticGet_Scope_SINGLE_TYPE()
	case "node":
		return StaticGet_Scope_SINGLE_NODE()
	case "key":
		return StaticGet_Scope_SINGLE_KEY()
	case "val":
		return StaticGet_Scope_SINGLE_VAL()
	case "err":
		return StaticGet_Scope_SINGLE_ERR()
	case "res":
		return StaticGet_Scope_SINGLE_RES()
	case "req":
		return StaticGet_Scope_SINGLE_REQ()
	case "cb":
		return StaticGet_Scope_SINGLE_CB()
	case "fn":
		return StaticGet_Scope_SINGLE_FN()
	case "el":
		return StaticGet_Scope_SINGLE_EL()
	case "id":
		return StaticGet_Scope_SINGLE_ID()
	case "self":
		return StaticGet_Scope_SINGLE_SELF()
	case "that":
		return StaticGet_Scope_SINGLE_THAT()
	case "args":
		return StaticGet_Scope_SINGLE_ARGS()
	default:
		return []string{name}
	}
	panic("unreachable switch expression")
}

// Java static method: Scope.pairVarNameArray(java.lang.String n0, java.lang.String n1)
func Scope_PairVarNameArray(n0 string, n1 string) []string {
	clinit_Scope()
	_ = n0
	_ = n1
	if ("a" == n0) && ("b" == n1) {
		return StaticGet_Scope_PAIR_A_B()
	}
	if ("x" == n0) && ("y" == n1) {
		return StaticGet_Scope_PAIR_X_Y()
	}
	if ("i" == n0) && ("n" == n1) {
		return StaticGet_Scope_PAIR_I_N()
	}
	if ("k" == n0) && ("v" == n1) {
		return StaticGet_Scope_PAIR_K_V()
	}
	if ("e" == n0) && ("t" == n1) {
		return StaticGet_Scope_PAIR_E_T()
	}
	if ("t" == n0) && ("e" == n1) {
		return StaticGet_Scope_PAIR_T_E()
	}
	if ("key" == n0) && ("value" == n1) {
		return StaticGet_Scope_PAIR_KEY_VALUE()
	}
	if ("resolve" == n0) && ("reject" == n1) {
		return StaticGet_Scope_PAIR_RESOLVE_REJECT()
	}
	return []string{n0, n1}
}

// Java method: Scope.getOwnVarNamesArray()
func (this_ *Scope) GetOwnVarNamesArray() []string {
	if this_.vars != nil {
		var n int32 = this_.vars.Size()
		_ = n
		if n == int32(0) {
			return StaticGet_Scope_EMPTY_VAR_NAMES()
		}
		var out []string = j2gort.NewStringArray(n)
		_ = out
		var i int32 = int32(0)
		_ = i
		{
			cur_18 := this_.vars.IterKeys()
			for cur_18.HasNext() {
				name := cur_18.
					Next()
				_ = name
				{
					a_19 := out
					i_20 := func() int32 { old_21 := i; nw_22 := old_21 + 1; i = nw_22; return old_21 }()
					a_19[i_20] = name
				}
			}
		}
		return out
	}
	var count int32 = this_.GetVarCount()
	_ = count
	if count == int32(0) {
		return StaticGet_Scope_EMPTY_VAR_NAMES()
	}
	if count == int32(1) {
		return Scope_SingletonVarNameArray(this_.var0.GetName())
	}
	if count == int32(2) {
		return Scope_PairVarNameArray(this_.var0.GetName(), this_.var1.GetName())
	}
	var out_2 []string = j2gort.NewStringArray(count)
	_ = out_2
	{
		a_23 := out_2
		i_24 := int32(0)
		a_23[i_24] = this_.var0.GetName()
	}
	{
		a_25 := out_2
		i_26 := int32(1)
		a_25[i_26] = this_.var1.GetName()
	}
	{
		a_27 := out_2
		i_28 := int32(2)
		a_27[i_28] = this_.var2.GetName()
	}
	if count > int32(3) {
		{
			a_29 := out_2
			i_30 := int32(3)
			a_29[i_30] = this_.var3.GetName()
		}
		if count > int32(4) {
			{
				a_31 := out_2
				i_32 := int32(4)
				a_31[i_32] = this_.var4.GetName()
			}
			if count > int32(5) {
				{
					a_33 := out_2
					i_34 := int32(5)
					a_33[i_34] = this_.var5.GetName()
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
	j2gort.Ext_java_util_Collections_AddAll(merged, func() []any {
		arr_35 := existing
		if arr_35 == nil {
			return nil
		}
		out_36 := make([]any, len(arr_35))
		for i_, el_37 := range arr_35 {
			out_36[i_] = j2gort.BoxString(el_37)
		}
		return out_36
	}())
	this_.AddVarNamesTo(merged)
	return func() []string {
		arr_41 := merged.ToArray_1_slice_any(func() []any {
			arr_38 := j2gort.NewStringArray(int32(0))
			if arr_38 == nil {
				return nil
			}
			out_39 := make([]any, len(arr_38))
			for i_, el_40 := range arr_38 {
				out_39[i_] = j2gort.BoxString(el_40)
			}
			return out_39
		}())
		if arr_41 == nil {
			return nil
		}
		out_42 := make([]string, len(arr_41))
		for i_, el_43 := range arr_41 {
			out_42[i_] = j2gort.CastString(el_43)
		}
		return out_42
	}()
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
