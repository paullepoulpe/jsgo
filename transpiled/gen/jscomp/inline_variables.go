/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InlineVariables.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_inline_variables_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_inline_variables_2 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_inline_variables_1 = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(mrp0_ any) bool {
		return j2gort.AsIface[AbstractVar](mrp0_).GetAbstractVarBase().IsDeclaredOrInferredConst()
	})
	lambdaSingleton_inline_variables_2 = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(mrp0_ any) bool {
		return j2gort.AsIface[AbstractVar](mrp0_).GetAbstractVarBase(

		// Java class: com.google.javascript.jscomp.InlineVariables
		).IsLocal()
	})
}

type InlineVariables struct {
	compiler                                                      AbstractCompiler      `gson:"compiler"`
	mode                                                          *InlineVariables_Mode `gson:"mode"`
	idHash_                                                       atomic.Int32
	slab_InlineVariables_VarExpert                                j2gort.PooledSlab[InlineVariables_VarExpert]
	slab_InlineVariables_InlineVarAnalysis                        j2gort.PooledSlab[InlineVariables_InlineVarAnalysis]
	slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler j2gort.PooledSlab[InlineVariables_InliningBehavior_AliasInlineRetryHandler]
}

func (this_ *InlineVariables) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineVariables) IsJavaNil_() bool { return this_ == nil }
func (*InlineVariables) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables", "InlineVariables"
}
func (this_ *InlineVariables) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineVariables) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables", "com.google.javascript.jscomp.InlineVariables", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables](v_)
	})
}
func (*InlineVariables) Implements_jscomp_CompilerPass() {}

func StaticGet_InlineVariables_UNKNOWN() int8 { return int8(-1) }
func StaticGet_InlineVariables_FALSE() int8   { return int8(0) }
func StaticGet_InlineVariables_TRUE() int8    { return int8(1) }

var static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS *InlineVariables_InlineVarAnalysis

func StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS() *InlineVariables_InlineVarAnalysis {
	clinit_InlineVariables()
	return static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS
}
func StaticSet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS(v *InlineVariables_InlineVarAnalysis) {
	clinit_InlineVariables()
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS = v
}

var static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS *InlineVariables_InlineVarAnalysis

func StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS() *InlineVariables_InlineVarAnalysis {
	clinit_InlineVariables()
	return static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS
}
func StaticSet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS(v *InlineVariables_InlineVarAnalysis) {
	clinit_InlineVariables()
	static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS = v
}

var static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT *InlineVariables_VarExpert

func StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT() *InlineVariables_VarExpert {
	clinit_InlineVariables()
	return static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT
}
func StaticSet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT(v *InlineVariables_VarExpert) {
	clinit_InlineVariables()
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT = v
}

var static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT *InlineVariables_VarExpert

func StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT() *InlineVariables_VarExpert {
	clinit_InlineVariables()
	return static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT
}
func StaticSet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT(v *InlineVariables_VarExpert) {
	clinit_InlineVariables()
	static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT = v
}

var clinitDone_InlineVariables uint32
var clinitState_InlineVariables j2gort.ClinitState

func clinit_InlineVariables() {
	if atomic.LoadUint32(&clinitDone_InlineVariables) != 0 {
		return
	}
	clinitSlow_InlineVariables()
}

func clinitSlow_InlineVariables() {
	if !j2gort.EnterClinit(&clinitState_InlineVariables, &clinitDone_InlineVariables) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InlineVariables, &clinitDone_InlineVariables)
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS = NewInlineVariables_InlineVarAnalysis(int8(0), nil, nil, pkg_rhino.Node(0), nil, nil)
	static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS = NewInlineVariables_InlineVarAnalysis(int8(1), nil, nil, pkg_rhino.Node(0), nil, nil)
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT = NewInlineVariables_VarExpert_1(StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS())
	static_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT = NewInlineVariables_VarExpert_1(StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS())
}

func (this_ *InlineVariables) initFields_InlineVariables() {
	_ = this_
}

// Java constructor: InlineVariables(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.InlineVariables.Mode mode)
func NewInlineVariables(compiler AbstractCompiler, mode *InlineVariables_Mode) *InlineVariables {
	clinit_InlineVariables()
	this_ := &InlineVariables{}
	InitInlineVariables(this_, compiler, mode)
	return this_
}

func InitInlineVariables(this_ *InlineVariables, compiler AbstractCompiler, mode *InlineVariables_Mode) {
	_ = compiler
	_ = mode
	this_.initFields_InlineVariables()
	this_.compiler = compiler
	this_.mode = mode
}

// Java method: InlineVariables.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineVariables) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.ReleaseSlabs_()
	clinit_ReferenceCollector()
	var val_1 ReferenceCollector
	InitReferenceCollector_4(&val_1, this_.compiler, NewInlineVariables_InliningBehavior(this_), NewSyntacticScopeCreator_1(this_.compiler), this_.mode.varPredicate)
	var callback *ReferenceCollector = &val_1
	_ = callback
	callback.ownsScopeCreator_ = true
	defer callback.ReleaseSlabs_()

	callback.Process_2(externs, root)
}

// Java static method: InlineVariables.hasNoInlineAnnotation(com.google.javascript.jscomp.Var var)
func InlineVariables_HasNoInlineAnnotation(var_ *Var) bool {
	clinit_InlineVariables()
	_ = var_
	var jsDocInfo *pkg_rhino.JSDocInfo = var_.GetJSDocInfo()
	_ = jsDocInfo
	return ((jsDocInfo != nil) && jsDocInfo.IsNoInline())
}

// Java enum: com.google.javascript.jscomp.InlineVariables.Mode
type InlineVariables_Mode struct {
	j2gort.EnumBase
	varPredicate j2gort.Ext_com_google_common_base_Predicate
}

func (this_ *InlineVariables_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *InlineVariables_Mode) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*InlineVariables_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*InlineVariables_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *InlineVariables_Mode) initFields_InlineVariables_Mode() {
	_ = this_
}

func initEnumInlineVariables_Mode(this_ *InlineVariables_Mode, enumName_ string, enumOrd_ int32, varPredicate j2gort.Ext_com_google_common_base_Predicate) *InlineVariables_Mode {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = varPredicate
	this_.initFields_InlineVariables_Mode()
	this_.varPredicate = varPredicate
	return this_
}

var enum_InlineVariables_Mode_CONSTANTS_ONLY *InlineVariables_Mode

func StaticGet_InlineVariables_Mode_CONSTANTS_ONLY() *InlineVariables_Mode {
	clinit_InlineVariables_Mode()
	return enum_InlineVariables_Mode_CONSTANTS_ONLY
}

var enum_InlineVariables_Mode_LOCALS_ONLY *InlineVariables_Mode

func StaticGet_InlineVariables_Mode_LOCALS_ONLY() *InlineVariables_Mode {
	clinit_InlineVariables_Mode()
	return enum_InlineVariables_Mode_LOCALS_ONLY
}

var enum_InlineVariables_Mode_ALL *InlineVariables_Mode

func StaticGet_InlineVariables_Mode_ALL() *InlineVariables_Mode {
	clinit_InlineVariables_Mode()
	return enum_InlineVariables_Mode_ALL
}

var clinitDone_InlineVariables_Mode uint32
var clinitState_InlineVariables_Mode j2gort.ClinitState

func clinit_InlineVariables_Mode() {
	if atomic.LoadUint32(&clinitDone_InlineVariables_Mode) != 0 {
		return
	}
	clinitSlow_InlineVariables_Mode()
}

func clinitSlow_InlineVariables_Mode() {
	if !j2gort.EnterClinit(&clinitState_InlineVariables_Mode, &clinitDone_InlineVariables_Mode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InlineVariables_Mode, &clinitDone_InlineVariables_Mode)
	enum_InlineVariables_Mode_CONSTANTS_ONLY = &InlineVariables_Mode{}
	initEnumInlineVariables_Mode(enum_InlineVariables_Mode_CONSTANTS_ONLY, "CONSTANTS_ONLY", int32(0), lambdaSingleton_inline_variables_1)
	enum_InlineVariables_Mode_LOCALS_ONLY = &InlineVariables_Mode{}
	initEnumInlineVariables_Mode(enum_InlineVariables_Mode_LOCALS_ONLY, "LOCALS_ONLY", int32(1), lambdaSingleton_inline_variables_2)
	enum_InlineVariables_Mode_ALL = &InlineVariables_Mode{}
	initEnumInlineVariables_Mode(enum_InlineVariables_Mode_ALL, "ALL", int32(2), j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
	enumValues_InlineVariables_Mode = []*InlineVariables_Mode{enum_InlineVariables_Mode_CONSTANTS_ONLY, enum_InlineVariables_Mode_LOCALS_ONLY, enum_InlineVariables_Mode_ALL}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineVariables.Mode", []any{enum_InlineVariables_Mode_CONSTANTS_ONLY, enum_InlineVariables_Mode_LOCALS_ONLY, enum_InlineVariables_Mode_ALL})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineVariables$Mode", []any{enum_InlineVariables_Mode_CONSTANTS_ONLY, enum_InlineVariables_Mode_LOCALS_ONLY, enum_InlineVariables_Mode_ALL})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.InlineVariables.Mode", "com.google.javascript.jscomp.InlineVariables$Mode", func() {
		clinit_InlineVariables_Mode()
	})
}

var enumValues_InlineVariables_Mode []*InlineVariables_Mode

func InlineVariables_Mode_Values() []*InlineVariables_Mode {
	clinit_InlineVariables_Mode()
	return j2gort.CloneSlice(enumValues_InlineVariables_Mode)
}

func InlineVariables_Mode_ValueOf(name string) *InlineVariables_Mode {
	clinit_InlineVariables_Mode()
	for _, c := range enumValues_InlineVariables_Mode {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.InlineVariables.Mode." + name))
}

func (this_ *InlineVariables_Mode) IsJavaNil_() bool { return this_ == nil }
func (*InlineVariables_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables$Mode", "Mode"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables.Mode", "com.google.javascript.jscomp.InlineVariables$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables_Mode](v_)
	})
}

// Java class: com.google.javascript.jscomp.InlineVariables.InlineVarAnalysis
type InlineVariables_InlineVarAnalysis struct {
	aliasedVar          *Var                       `gson:"aliasedVar"`
	expert              *InlineVariables_VarExpert `gson:"expert"`
	initialization      *Reference                 `gson:"initialization"`
	singleReadReference *Reference                 `gson:"singleReadReference"`
	initValue           pkg_rhino.Node             `gson:"initValue"`
	idHash_             atomic.Int32
	kind                int8 `gson:"kind"`
	lazyRefInfo_        *ReferenceCollection
	lazyWellDefined_    int8
	lazyAssignedOnce_   int8
}

func (this_ *InlineVariables_InlineVarAnalysis) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineVariables_InlineVarAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*InlineVariables_InlineVarAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables$InlineVarAnalysis", "InlineVarAnalysis"
}
func (this_ *InlineVariables_InlineVarAnalysis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineVariables_InlineVarAnalysis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables.InlineVarAnalysis", "com.google.javascript.jscomp.InlineVariables$InlineVarAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables_InlineVarAnalysis](v_)
	})
}

func StaticGet_InlineVariables_InlineVarAnalysis_KIND_NO_INLINE_UNSAFE_ALIASES() int8 { return int8(0) }
func StaticGet_InlineVariables_InlineVarAnalysis_KIND_NO_INLINE_SAFE_ALIASES() int8   { return int8(1) }
func StaticGet_InlineVariables_InlineVarAnalysis_KIND_NO_INLINE_LAZY_ALIASES() int8   { return int8(2) }
func StaticGet_InlineVariables_InlineVarAnalysis_KIND_WAIT_FOR_ALIAS() int8           { return int8(3) }
func StaticGet_InlineVariables_InlineVarAnalysis_KIND_INLINE_UNDEFINED() int8         { return int8(4) }
func StaticGet_InlineVariables_InlineVarAnalysis_KIND_INLINE_WELL_DEFINED() int8      { return int8(5) }
func StaticGet_InlineVariables_InlineVarAnalysis_KIND_INLINE_SINGLE_READ() int8       { return int8(6) }
func (this_ *InlineVariables_InlineVarAnalysis) initFields_InlineVariables_InlineVarAnalysis() {
	_ = this_
}

// Java constructor: InlineVarAnalysis(byte kind, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Var aliasedVar, com.google.javascript.jscomp.InlineVariables.@org.jspecify.annotations.Nullable VarExpert expert, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node initValue, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference initialization, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference singleReadReference)
func NewInlineVariables_InlineVarAnalysis(kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	this_ := &InlineVariables_InlineVarAnalysis{}
	InitInlineVariables_InlineVarAnalysis(this_, kind, aliasedVar, expert, initValue, initialization, singleReadReference)
	return this_
}
func allocIvPooled_InlineVariables_InlineVarAnalysis(beh_ *InlineVariables_InliningBehavior, kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	if beh_ == nil || beh_.outer_ == nil {
		return NewInlineVariables_InlineVarAnalysis(kind, aliasedVar, expert, initValue, initialization, singleReadReference)
	}
	this_ := beh_.outer_.slab_InlineVariables_InlineVarAnalysis.Alloc()
	InitInlineVariables_InlineVarAnalysis(this_, kind, aliasedVar, expert, initValue, initialization, singleReadReference)
	return this_
}

func InitInlineVariables_InlineVarAnalysis(this_ *InlineVariables_InlineVarAnalysis, kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) {
	_ = kind
	_ = aliasedVar
	_ = expert
	_ = initValue
	_ = initialization
	_ = singleReadReference
	this_.initFields_InlineVariables_InlineVarAnalysis()
	this_.kind = kind
	this_.aliasedVar = aliasedVar
	this_.expert = expert
	this_.initValue = initValue
	this_.initialization = initialization
	this_.singleReadReference = singleReadReference
	this_.

		// Java static method: InlineVarAnalysis.createLazyNegative(com.google.javascript.jscomp.InlineVariables.VarExpert expert)
		lazyRefInfo_ = nil
}

func InlineVariables_InlineVarAnalysis_CreateLazyNegative(expert *InlineVariables_VarExpert) *InlineVariables_InlineVarAnalysis {
	_ = expert
	return NewInlineVariables_InlineVarAnalysis(int8(2), nil, expert, pkg_rhino.Node(0), nil, nil)
}
func allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateLazyNegative(beh_ *InlineVariables_InliningBehavior, expert *InlineVariables_VarExpert) *InlineVariables_InlineVarAnalysis {
	_ = expert
	if beh_ != nil && beh_.outer_ != nil && expert != nil && expert.fixedAnalysis == nil && expert.referenceInfo != nil {
		if expert.isWellDefined == 0 || expert.isAssignedOnceInLifetime == 0 {
			return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS()
		}
		if expert.isWellDefined == 1 && expert.isAssignedOnceInLifetime == 1 {
			return StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS()
		}
		if beh_.activeRefMap_ != nil {
			if expert.isWellDefined == 1 {
				return &staticIvLazyNegWellDef_
			}
			if expert.isAssignedOnceInLifetime == 1 {
				return &staticIvLazyNegAssignedOnce_
			}
			return &staticIvLazyNegUnknown_
		}
		a_ := beh_.outer_.slab_InlineVariables_InlineVarAnalysis.Alloc()
		a_.kind = int8(2)
		a_.lazyRefInfo_ = expert.referenceInfo
		a_.lazyWellDefined_ = expert.isWellDefined
		a_.lazyAssignedOnce_ = expert.isAssignedOnceInLifetime
		return a_
	}
	expert = promoteIvScratchVarExpert_(beh_, expert)

	return allocIvPooled_InlineVariables_InlineVarAnalysis(beh_, int8(2), nil, expert, pkg_rhino.Node(0), nil, nil)
}

// Java static method: InlineVarAnalysis.createWaitForAlias(com.google.javascript.jscomp.Var aliasedVar)
func InlineVariables_InlineVarAnalysis_CreateWaitForAlias(aliasedVar *Var) *InlineVariables_InlineVarAnalysis {
	_ = aliasedVar
	return NewInlineVariables_InlineVarAnalysis(int8(3), aliasedVar, nil, pkg_rhino.Node(0), nil, nil)
}
func allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateWaitForAlias(beh_ *InlineVariables_InliningBehavior, aliasedVar *Var) *InlineVariables_InlineVarAnalysis {
	_ = aliasedVar
	return allocIvPooled_InlineVariables_InlineVarAnalysis(beh_, int8(3), aliasedVar, nil, pkg_rhino.Node(0), nil, nil)
}

// Java static method: InlineVarAnalysis.createInlineUndefined(com.google.javascript.jscomp.InlineVariables.VarExpert expert)
func InlineVariables_InlineVarAnalysis_CreateInlineUndefined(expert *InlineVariables_VarExpert) *InlineVariables_InlineVarAnalysis {
	_ = expert
	return NewInlineVariables_InlineVarAnalysis(int8(4), nil, expert, pkg_rhino.Node(0), nil, nil)
}
func allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateInlineUndefined(beh_ *InlineVariables_InliningBehavior, expert *InlineVariables_VarExpert) *InlineVariables_InlineVarAnalysis {
	_ = expert
	return allocIvPooled_InlineVariables_InlineVarAnalysis(beh_, int8(4), nil, expert, pkg_rhino.Node(0), nil, nil)
}

// Java static method: InlineVarAnalysis.createInlineWellDefined(com.google.javascript.jscomp.InlineVariables.VarExpert expert, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node initValue)
func InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(expert *InlineVariables_VarExpert, initValue pkg_rhino.Node) *InlineVariables_InlineVarAnalysis {
	_ = expert
	_ = initValue
	return NewInlineVariables_InlineVarAnalysis(int8(5), nil, expert, initValue, nil, nil)
}
func allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(beh_ *InlineVariables_InliningBehavior, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node) *InlineVariables_InlineVarAnalysis {
	_ = expert
	_ = initValue
	return allocIvPooled_InlineVariables_InlineVarAnalysis(beh_, int8(5), nil, expert, initValue, nil, nil)
}

// Java static method: InlineVarAnalysis.createInlineSingleRead(com.google.javascript.jscomp.InlineVariables.VarExpert expert, com.google.javascript.jscomp.Reference initialization, com.google.javascript.jscomp.Reference singleReadReference)
func InlineVariables_InlineVarAnalysis_CreateInlineSingleRead(expert *InlineVariables_VarExpert, initialization *Reference, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	_ = expert
	_ = initialization
	_ = singleReadReference
	return NewInlineVariables_InlineVarAnalysis(int8(6), nil, expert, pkg_rhino.Node(0), initialization, singleReadReference)
}
func allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateInlineSingleRead(beh_ *InlineVariables_InliningBehavior, expert *InlineVariables_VarExpert, initialization *Reference, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	_ = expert
	_ = initialization
	_ = singleReadReference
	return allocIvPooled_InlineVariables_InlineVarAnalysis(beh_, int8(6), nil, expert, pkg_rhino.Node(0), initialization, singleReadReference)
}

// Java method: InlineVarAnalysis.shouldInline()
func (this_ *InlineVariables_InlineVarAnalysis) ShouldInline() bool {
	return (int32(this_.kind) >= int32(int8(4)))
}

// Java method: InlineVarAnalysis.shouldWaitForAliasedVar()
func (this_ *InlineVariables_InlineVarAnalysis) ShouldWaitForAliasedVar() bool {
	return (int32(this_.kind) == int32(int8(3)))
}

// Java method: InlineVarAnalysis.getAliasedVar()
func (this_ *InlineVariables_InlineVarAnalysis) GetAliasedVar() *Var {
	if int32(this_.kind) == int32(int8(3)) {
		return this_.aliasedVar
	}
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "no aliased Var")))
}

// Java method: InlineVarAnalysis.isSafeToInlineAliases()
func (this_ *InlineVariables_InlineVarAnalysis) IsSafeToInlineAliases() bool {
	switch this_.kind {
	case int8(0):
		return false
	case int8(1), int8(4), int8(5), int8(6):
		return true
	case int8(2):
		if this_.lazyRefInfo_ != nil {
			if this_.lazyWellDefined_ == 0 || this_.lazyAssignedOnce_ == 0 {
				return false
			}
			if this_.lazyWellDefined_ < 0 {
				if this_.lazyRefInfo_.IsWellDefined() {
					this_.lazyWellDefined_ = 1
				} else {
					this_.lazyWellDefined_ = 0
					return false
				}
			}
			if this_.lazyAssignedOnce_ < 0 {
				if this_.lazyRefInfo_.IsAssignedOnceInLifetime() {
					this_.lazyAssignedOnce_ = 1
				} else {
					this_.lazyAssignedOnce_ = 0
					return false
				}
			}
			return true
		}

		return this_.expert.IsWellDefinedAssignedOnce()
	case int8(3):
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "analysis is incomplete")))
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	panic("unreachable switch expression")
}

// Java method: InlineVarAnalysis.performInline()
func (this_ *InlineVariables_InlineVarAnalysis) PerformInline() {
	switch this_.kind {
	case int8(4):
		this_.expert.InlineUndefined()
	case int8(5):
		this_.expert.InlineWellDefined(this_.initValue)
	case int8(6):
		this_.expert.InlineSingleRead(this_.initialization, this_.singleReadReference)
	default:
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "cannot inline")))
	}
}

// Java class: com.google.javascript.jscomp.InlineVariables.VarExpert
type InlineVariables_VarExpert struct {
	fixedAnalysis                       *InlineVariables_InlineVarAnalysis `gson:"fixedAnalysis"`
	behavior                            *InlineVariables_InliningBehavior  `gson:"behavior"`
	v                                   *Var                               `gson:"v"`
	referenceInfo                       *ReferenceCollection               `gson:"referenceInfo"`
	declaration                         *Reference                         `gson:"declaration"`
	initializationReference             *Reference                         `gson:"initializationReference"`
	idHash_                             atomic.Int32
	isDeclaredOrInferredConstant        bool `gson:"isDeclaredOrInferredConstant"`
	mayBeAParameterModifiedViaArguments bool `gson:"mayBeAParameterModifiedViaArguments"`
	isNeverAssigned                     int8 `gson:"isNeverAssigned"`
	isWellDefined                       int8 `gson:"isWellDefined"`
	isAssignedOnceInLifetime            int8 `gson:"isAssignedOnceInLifetime"`
	isReferencedWeakly                  int8 `gson:"isReferencedWeakly"`
	isWellDefinedAssignedOnce           int8 `gson:"isWellDefinedAssignedOnce"`
	initializationReferenceKnown        bool `gson:"initializationReferenceKnown"`
	allReferencesAreValid               int8 `gson:"allReferencesAreValid"`
}

func (this_ *InlineVariables_VarExpert) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineVariables_VarExpert) IsJavaNil_() bool { return this_ == nil }
func (*InlineVariables_VarExpert) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables$VarExpert", "VarExpert"
}
func (this_ *InlineVariables_VarExpert) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineVariables_VarExpert) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables.VarExpert", "com.google.javascript.jscomp.InlineVariables$VarExpert", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables_VarExpert](v_)
	})
}

func (this_ *InlineVariables_VarExpert) initFields_InlineVariables_VarExpert() {
	_ = this_
	this_.isNeverAssigned = int8(-1)
	this_.isWellDefined = int8(-1)
	this_.isAssignedOnceInLifetime = int8(-1)
	this_.isReferencedWeakly = int8(-1)
	this_.isWellDefinedAssignedOnce = int8(-1)

	this_.allReferencesAreValid = int8(-1)
}

// Java constructor: VarExpert(com.google.javascript.jscomp.InlineVariables.InlineVarAnalysis fixedAnalysis)
func NewInlineVariables_VarExpert_1(fixedAnalysis *InlineVariables_InlineVarAnalysis) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_1(this_, fixedAnalysis)
	return this_
}

func InitInlineVariables_VarExpert_1(this_ *InlineVariables_VarExpert, fixedAnalysis *InlineVariables_InlineVarAnalysis) {
	_ = fixedAnalysis
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = fixedAnalysis
	this_.behavior = nil
	this_.v = nil
	this_.referenceInfo = nil
	this_.declaration = nil
	this_.isDeclaredOrInferredConstant = false
	this_.mayBeAParameterModifiedViaArguments = false
}

// Java constructor: VarExpert(com.google.javascript.jscomp.InlineVariables.InliningBehavior behavior, com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.ReferenceCollection referenceInfo, boolean isDeclaredOrInferredConstant, boolean mayBeAParameterModifiedViaArguments)
func NewInlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isDeclaredOrInferredConstant bool, mayBeAParameterModifiedViaArguments bool) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_5(this_, behavior, v, referenceInfo, isDeclaredOrInferredConstant, mayBeAParameterModifiedViaArguments)
	return this_
}
func allocIvPooled_InlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isDeclaredOrInferredConstant bool, mayBeAParameterModifiedViaArguments bool) *InlineVariables_VarExpert {
	if behavior == nil || behavior.outer_ == nil {
		return NewInlineVariables_VarExpert_5(behavior, v, referenceInfo, isDeclaredOrInferredConstant, mayBeAParameterModifiedViaArguments)
	}
	this_ := &behavior.scratchVarExpert_
	this_.isNeverAssigned = int8(-1)
	this_.isWellDefined = int8(-1)
	this_.isAssignedOnceInLifetime = int8(-1)
	this_.isReferencedWeakly = int8(-1)
	this_.isWellDefinedAssignedOnce = int8(-1)
	this_.initializationReferenceKnown = false
	this_.initializationReference = nil
	this_.allReferencesAreValid = int8(-1)
	this_.fixedAnalysis = nil
	this_.behavior = behavior
	this_.v = v
	this_.referenceInfo = referenceInfo
	this_.declaration = referenceInfo.
		Get(int32(0))
	this_.isDeclaredOrInferredConstant = isDeclaredOrInferredConstant
	this_.mayBeAParameterModifiedViaArguments = mayBeAParameterModifiedViaArguments
	return this_
}
func promoteIvScratchVarExpert_(beh_ *InlineVariables_InliningBehavior, exp_ *InlineVariables_VarExpert) *InlineVariables_VarExpert {
	if beh_ == nil || beh_.outer_ == nil || exp_ == nil || exp_ != &beh_.scratchVarExpert_ {
		return exp_
	}
	dst_ := beh_.outer_.slab_InlineVariables_VarExpert.Alloc()
	dst_.fixedAnalysis = exp_.fixedAnalysis
	dst_.behavior = exp_.behavior
	dst_.v = exp_.v
	dst_.referenceInfo = exp_.referenceInfo
	dst_.declaration = exp_.declaration
	dst_.initializationReference = exp_.initializationReference
	dst_.isDeclaredOrInferredConstant = exp_.isDeclaredOrInferredConstant
	dst_.mayBeAParameterModifiedViaArguments = exp_.mayBeAParameterModifiedViaArguments
	dst_.isNeverAssigned = exp_.isNeverAssigned
	dst_.isWellDefined = exp_.isWellDefined
	dst_.isAssignedOnceInLifetime = exp_.isAssignedOnceInLifetime
	dst_.isReferencedWeakly = exp_.isReferencedWeakly
	dst_.isWellDefinedAssignedOnce = exp_.isWellDefinedAssignedOnce
	dst_.initializationReferenceKnown = exp_.initializationReferenceKnown
	dst_.allReferencesAreValid = exp_.allReferencesAreValid
	return dst_
}

func InitInlineVariables_VarExpert_5(this_ *InlineVariables_VarExpert, behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isDeclaredOrInferredConstant bool, mayBeAParameterModifiedViaArguments bool) {
	_ = behavior
	_ = v
	_ = referenceInfo
	_ = isDeclaredOrInferredConstant
	_ = mayBeAParameterModifiedViaArguments
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = nil
	this_.behavior = behavior
	this_.v = v
	this_.referenceInfo = referenceInfo
	this_.declaration = referenceInfo.Get(int32(0))
	this_.isDeclaredOrInferredConstant = isDeclaredOrInferredConstant
	this_.mayBeAParameterModifiedViaArguments = mayBeAParameterModifiedViaArguments
}

// Java method: VarExpert.isNeverAssigned()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_isNeverAssigned_3372() bool {
	if int32(this_.isNeverAssigned) != int32(int8(-1)) {
		return (int32(this_.isNeverAssigned) == int32(int8(1)))
	} else if (this_.initializationReferenceKnown && (this_.initializationReference != nil)) || (int32(this_.isAssignedOnceInLifetime) == int32(int8(1))) {
		this_.isNeverAssigned = int8(0)
		return false
	} else {
		var val bool = this_.referenceInfo.IsNeverAssigned()
		_ = val
		if !(int32(this_.isNeverAssigned) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		if val {
			this_.isNeverAssigned = int8(1)
		} else {
			this_.isNeverAssigned = int8(0)
		}
		return val
	}
}

// Java method: VarExpert.isWellDefined()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_isWellDefined_3373() bool {
	if int32(this_.isWellDefined) != int32(int8(-1)) {
		return (int32(this_.isWellDefined) == int32(int8(1)))
	} else {
		var val bool = this_.referenceInfo.IsWellDefined()
		_ = val
		if !(int32(this_.isWellDefined) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		if val {
			this_.isWellDefined = int8(1)
		} else {
			this_.isWellDefined = int8(0)
		}
		return val
	}
}

// Java method: VarExpert.isAssignedOnceInLifetime()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_isAssignedOnceInLifetime_3374() bool {
	if int32(this_.isAssignedOnceInLifetime) != int32(int8(-1)) {
		return (int32(this_.isAssignedOnceInLifetime) == int32(int8(1)))
	} else {
		var val bool = this_.referenceInfo.IsAssignedOnceInLifetime()
		_ = val
		if !(int32(this_.isAssignedOnceInLifetime) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		if val {
			this_.isAssignedOnceInLifetime = int8(1)
		} else {
			this_.isAssignedOnceInLifetime = int8(0)
		}
		return val
	}
}

// Java method: VarExpert.isReferencedWeakly()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_isReferencedWeakly_3375() bool {
	if int32(this_.isReferencedWeakly) != int32(int8(-1)) {
		return (int32(this_.isReferencedWeakly) == int32(int8(1)))
	} else {
		var val bool = this_.referenceInfo.IsReferencedWeakly()
		_ = val
		if !(int32(this_.isReferencedWeakly) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		if val {
			this_.isReferencedWeakly = int8(1)
		} else {
			this_.isReferencedWeakly = int8(0)
		}
		return val
	}
}

// Java method: VarExpert.isWellDefinedAssignedOnce()
func (this_ *InlineVariables_VarExpert) IsWellDefinedAssignedOnce() bool {
	if int32(this_.isWellDefinedAssignedOnce) != int32(int8(-1)) {
		return (int32(this_.isWellDefinedAssignedOnce) == int32(int8(1)))
	} else {
		var val bool = ((((int32(this_.isWellDefined) != int32(int8(0))) && (int32(this_.isAssignedOnceInLifetime) != int32(int8(0)))) && this_.priv_VarExpert_isWellDefined_3373()) && this_.priv_VarExpert_isAssignedOnceInLifetime_3374())
		_ = val
		if !(int32(this_.isWellDefinedAssignedOnce) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		if val {
			this_.isWellDefinedAssignedOnce = int8(1)
		} else {
			this_.isWellDefinedAssignedOnce = int8(0)
		}
		return val
	}
}

// Java method: VarExpert.getInitialization()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_getInitialization_3377() *Reference {
	if this_.initializationReferenceKnown {
		return this_.initializationReference
	} else {
		var initRef *Reference
		_ = initRef
		if this_.isDeclaredOrInferredConstant {
			initRef = this_.referenceInfo.GetInitializingReferenceForConstants()
		} else {
			initRef = this_.referenceInfo.GetInitializingReference()
		}
		if !(!(this_.initializationReferenceKnown)) {
			j2gort.CheckStateStr(false, "already known")
		}
		this_.initializationReference = initRef
		this_.initializationReferenceKnown = true
		return initRef
	}
}

// Java method: VarExpert.hasValidDeclaration()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_hasValidDeclaration_3378() bool {
	return this_.behavior.IsValidDeclaration(this_.declaration)
}

// Java method: VarExpert.hasValidInitialization()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_hasValidInitialization_3379() bool {
	return this_.behavior.IsValidInitialization(this_.priv_VarExpert_getInitialization_3377())
}

// Java method: VarExpert.allReferencesAreValid()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_allReferencesAreValid_3380() bool {
	if int32(this_.allReferencesAreValid) != int32(int8(-1)) {
		return (int32(this_.allReferencesAreValid) == int32(int8(1)))
	} else {
		var hasValidDeclaration bool = this_.priv_VarExpert_hasValidDeclaration_3378()
		_ = hasValidDeclaration
		var hasValidInitialization bool = this_.priv_VarExpert_hasValidInitialization_3379()
		_ = hasValidInitialization
		var neverAssigned bool = this_.priv_VarExpert_isNeverAssigned_3372()
		_ = neverAssigned
		if hasValidDeclaration && (hasValidInitialization || neverAssigned) {
			var initialization *Reference = this_.priv_VarExpert_getInitialization_3377()
			_ = initialization
			{
				var i int32 = int32(0)
				_ = i
				var size int32 = this_.referenceInfo.Size()
				_ = size
				for ; i < size; i = i + 1 {
					var ref *Reference = this_.referenceInfo.Get(i)
					_ = ref
					if ((ref != this_.declaration) && (ref != initialization)) && !(this_.behavior.IsValidReference(ref)) {
						if !(int32(this_.allReferencesAreValid) == int32(int8(-1))) {
							j2gort.CheckStateStr(false, "already known")
						}
						this_.allReferencesAreValid = int8(0)
						return false
					}
				}
			}
			if !(int32(this_.allReferencesAreValid) == int32(int8(-1))) {
				j2gort.CheckStateStr(false, "already known")
			}
			this_.allReferencesAreValid = int8(1)
			return true
		}
		if !(int32(this_.allReferencesAreValid) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		this_.allReferencesAreValid = int8(0)
		return false
	}
}

// Java method: VarExpert.inlineUndefined()
func (this_ *InlineVariables_VarExpert) InlineUndefined() {
	var srcLocation pkg_rhino.Node = this_.declaration.GetNode()
	_ = srcLocation
	var undefinedNode pkg_rhino.Node = NodeUtil_NewUndefinedNode(srcLocation)
	_ = undefinedNode
	this_.behavior.InlineWellDefinedVariable(this_.v, undefinedNode, this_.referenceInfo)
}

// Java method: VarExpert.inlineWellDefined(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node initValue)
func (this_ *InlineVariables_VarExpert) InlineWellDefined(initValue pkg_rhino.Node) {
	_ = initValue
	this_.behavior.InlineWellDefinedVariable(this_.v, initValue, this_.referenceInfo)
}

// Java method: VarExpert.inlineSingleRead(com.google.javascript.jscomp.Reference initialization, com.google.javascript.jscomp.Reference singleReadReference)
func (this_ *InlineVariables_VarExpert) InlineSingleRead(initialization *Reference, singleReadReference *Reference) {
	_ = initialization
	_ = singleReadReference
	this_.behavior.Inline(this_.declaration, initialization, singleReadReference)
}

// Java method: VarExpert.analyze()
func (this_ *InlineVariables_VarExpert) Analyze() *InlineVariables_InlineVarAnalysis {
	if this_.fixedAnalysis != nil {
		return this_.fixedAnalysis
	}
	if InlineVariables_HasNoInlineAnnotation(this_.v) || this_.priv_VarExpert_isReferencedWeakly_3375() {
		return this_.priv_VarExpert_getNegativeInlineVarAnalysis_3386()
	}
	var initialization *Reference = this_.priv_VarExpert_getInitialization_3377()
	_ = initialization
	var initValue pkg_rhino.Node
	_ = initValue
	if initialization == nil {
		initValue = pkg_rhino.Node(0)
	} else {
		initValue = initialization.GetAssignedValue()
	}
	var val_2 InlineVariables_VarExpert_InitialValueAnalysis
	InitInlineVariables_VarExpert_InitialValueAnalysis(&val_2, this_.v, initValue)
	var initialValueAnalysis *InlineVariables_VarExpert_InitialValueAnalysis = &val_2
	_ = initialValueAnalysis
	if (this_.isDeclaredOrInferredConstant && initialValueAnalysis.IsImmutableValue()) && this_.priv_VarExpert_isAssignedOnceInLifetime_3374() {
		return this_.priv_VarExpert_createInlineWellDefinedVarAnalysis_3387(initValue)
	}
	if this_.behavior.Mode() == StaticGet_InlineVariables_Mode_CONSTANTS_ONLY() {
		return this_.priv_VarExpert_getNegativeInlineVarAnalysis_3386()
	}
	if initialValueAnalysis.IsAlias() {
		return allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateWaitForAlias(this_.behavior, initialValueAnalysis.GetAliasedVar())
	}
	return this_.priv_VarExpert_analyzeWithInitialValue_3385(initialization, initValue, initialValueAnalysis)
}

// Java method: VarExpert.analyzeWithInitialValue(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference initialization, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node initValue, com.google.javascript.jscomp.InlineVariables.VarExpert.InitialValueAnalysis initialValueAnalysis)
func (this_ *InlineVariables_VarExpert) priv_VarExpert_analyzeWithInitialValue_3385(initialization *Reference, initValue pkg_rhino.Node, initialValueAnalysis *InlineVariables_VarExpert_InitialValueAnalysis) *InlineVariables_InlineVarAnalysis {
	_ = initialization
	_ = initValue
	_ = initialValueAnalysis
	var refCount int32 = this_.referenceInfo.Size()
	_ = refCount
	if (refCount > int32(1)) && this_.priv_VarExpert_allReferencesAreValid_3380() {
		if this_.referenceInfo.IsNeverAssigned() {
			return allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateInlineUndefined(this_.behavior, this_)
		}
		if this_.priv_VarExpert_isWellDefined_3373() && (initialValueAnalysis.IsImmutableValue() || (initValue.IsThis() && !(this_.referenceInfo.IsThisRebound()))) {
			return this_.priv_VarExpert_createInlineWellDefinedVarAnalysis_3387(initValue)
		}
		var firstReadRefIndex int32
		_ = firstReadRefIndex
		if this_.declaration == initialization {
			firstReadRefIndex = int32(1)
		} else {
			firstReadRefIndex = int32(2)
		}
		var numReadRefs int32 = (refCount - firstReadRefIndex)
		_ = numReadRefs
		if numReadRefs == int32(0) {
			return this_.priv_VarExpert_createInlineWellDefinedVarAnalysis_3387(initValue)
		}
		if numReadRefs == int32(1) {
			var singleReadReference *Reference = this_.referenceInfo.Get(firstReadRefIndex)
			_ = singleReadReference
			if this_.behavior.CanInline(this_.declaration, initialization, singleReadReference, initValue) {
				return allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateInlineSingleRead(this_.behavior, this_, initialization, singleReadReference)
			}
		}
	}
	return this_.priv_VarExpert_getNegativeInlineVarAnalysis_3386()
}

// Java method: VarExpert.getNegativeInlineVarAnalysis()
func (this_ *InlineVariables_VarExpert) priv_VarExpert_getNegativeInlineVarAnalysis_3386() *InlineVariables_InlineVarAnalysis {
	if (this_.mayBeAParameterModifiedViaArguments || (this_.behavior.Mode() == StaticGet_InlineVariables_Mode_CONSTANTS_ONLY())) || (int32(this_.isWellDefinedAssignedOnce) == int32(int8(0))) {
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS()
	}
	if int32(this_.isWellDefinedAssignedOnce) == int32(int8(1)) {
		return StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS()
	}
	return allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateLazyNegative(this_.behavior, this_)
}

// Java method: VarExpert.createInlineWellDefinedVarAnalysis(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node initValue)
func (this_ *InlineVariables_VarExpert) priv_VarExpert_createInlineWellDefinedVarAnalysis_3387(initValue pkg_rhino.Node) *InlineVariables_InlineVarAnalysis {
	_ = initValue
	return allocIvPooledVia_InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_.behavior, this_, initValue)
}

// Java method: VarExpert.reanalyzeAfterAliasedVar(com.google.javascript.jscomp.Var aliasedVar, com.google.javascript.jscomp.InlineVariables.InlineVarAnalysis aliasedVarAnalysis)
func (this_ *InlineVariables_VarExpert) ReanalyzeAfterAliasedVar(aliasedVar *Var, aliasedVarAnalysis *InlineVariables_InlineVarAnalysis) *InlineVariables_InlineVarAnalysis {
	_ = aliasedVar
	_ = aliasedVarAnalysis
	if this_.fixedAnalysis != nil {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "not waiting for an aliased variable")))
	}
	var initialization *Reference = this_.priv_VarExpert_getInitialization_3377()
	_ = initialization
	var initValue pkg_rhino.Node
	_ = initValue
	if initialization == nil {
		initValue = pkg_rhino.Node(0)
	} else {
		initValue = initialization.GetAssignedValue()
	}
	var val_3 InlineVariables_VarExpert_InitialValueAnalysis
	InitInlineVariables_VarExpert_InitialValueAnalysis(&val_3, this_.v, initValue)
	var initialValueAnalysis *InlineVariables_VarExpert_InitialValueAnalysis = &val_3
	_ = initialValueAnalysis
	if (this_.isDeclaredOrInferredConstant && initialValueAnalysis.IsImmutableValue()) && this_.priv_VarExpert_isAssignedOnceInLifetime_3374() {
		return this_.priv_VarExpert_createInlineWellDefinedVarAnalysis_3387(initValue)
	}
	if ((initialValueAnalysis.IsAlias() && j2gort.ObjEqCustomPtr(aliasedVar, initialValueAnalysis.GetAliasedVar())) && isSafeToInlineAliasesIv_(this_.behavior, aliasedVar, aliasedVarAnalysis)) && this_.IsWellDefinedAssignedOnce() {
		return this_.priv_VarExpert_createInlineWellDefinedVarAnalysis_3387(initValue)
	}
	return this_.priv_VarExpert_analyzeWithInitialValue_3385(initialization, initValue, initialValueAnalysis)
}

// Java class: com.google.javascript.jscomp.InlineVariables.VarExpert.InitialValueAnalysis
type InlineVariables_VarExpert_InitialValueAnalysis struct {
	v                *Var           `gson:"v"`
	aliasedVar       *Var           `gson:"aliasedVar"`
	value            pkg_rhino.Node `gson:"value"`
	idHash_          atomic.Int32
	isImmutableValue int8 `gson:"isImmutableValue"`
	aliasedVarKnown  bool `gson:"aliasedVarKnown"`
}

func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*InlineVariables_VarExpert_InitialValueAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables$VarExpert$InitialValueAnalysis", "InitialValueAnalysis"
}
func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineVariables_VarExpert_InitialValueAnalysis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables.VarExpert.InitialValueAnalysis", "com.google.javascript.jscomp.InlineVariables$VarExpert$InitialValueAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables_VarExpert_InitialValueAnalysis](v_)
	})
}

func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) initFields_InlineVariables_VarExpert_InitialValueAnalysis() {
	_ = this_
	this_.isImmutableValue = int8(-1)

}

// Java constructor: InitialValueAnalysis(com.google.javascript.jscomp.Var v, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value)
func NewInlineVariables_VarExpert_InitialValueAnalysis(v *Var, value pkg_rhino.Node) *InlineVariables_VarExpert_InitialValueAnalysis {
	this_ := &InlineVariables_VarExpert_InitialValueAnalysis{}
	InitInlineVariables_VarExpert_InitialValueAnalysis(this_, v, value)
	return this_
}

func InitInlineVariables_VarExpert_InitialValueAnalysis(this_ *InlineVariables_VarExpert_InitialValueAnalysis, v *Var, value pkg_rhino.Node) {
	_ = v
	_ = value
	this_.initFields_InlineVariables_VarExpert_InitialValueAnalysis()
	this_.v = v
	this_.value = value
}

// Java method: InitialValueAnalysis.isImmutableValue()
func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) IsImmutableValue() bool {
	if int32(this_.isImmutableValue) != int32(int8(-1)) {
		return (int32(this_.isImmutableValue) == int32(int8(1)))
	} else {
		var val bool = ((this_.value != 0) && NodeUtil_IsImmutableValue(this_.value))
		_ = val
		if !(int32(this_.isImmutableValue) == int32(int8(-1))) {
			j2gort.CheckStateStr(false, "already known")
		}
		if val {
			this_.isImmutableValue = int8(1)
		} else {
			this_.isImmutableValue = int8(0)
		}
		return val
	}
}

// Java method: InitialValueAnalysis.isAlias()
func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) IsAlias() bool {
	return (this_.GetAliasedVar() != nil)
}

// Java method: InitialValueAnalysis.getAliasedVar()
func (this_ *InlineVariables_VarExpert_InitialValueAnalysis) GetAliasedVar() *Var {
	if this_.aliasedVarKnown {
		return this_.aliasedVar
	} else {
		var resolved *Var = nil
		_ = resolved
		if (this_.value != 0) && this_.value.IsName() {
			var aliasedName string = this_.value.GetString()
			_ = aliasedName
			{
				var iifeCond_1 bool
				{

					sr_4 := aliasedName
					se_5 := this_.v.GetName()
					iifeCond_1 = j2gort.CheckNotNullString(sr_4) == se_5
				}
				if !(iifeCond_1) {
					resolved = j2gort.CastPtr[Scope](this_.v.GetScope()).LookupVar(aliasedName)
				}
			}
		}
		if !(!(this_.aliasedVarKnown)) {
			j2gort.CheckStateStr(false, "already known")
		}
		this_.aliasedVar = resolved
		this_.aliasedVarKnown = true
		return resolved
	}
}

// Java class: com.google.javascript.jscomp.InlineVariables.InliningBehavior
type InlineVariables_InliningBehavior struct {
	outer_                            *InlineVariables
	currentScopeHandledVarAnalysesMap j2gort.OrderedMap[*Var, *InlineVariables_InlineVarAnalysis]                      `gson:"currentScopeHandledVarAnalysesMap"`
	varToAliasRetryHandlersMap        *SimpleMultimap[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler] `gson:"varToAliasRetryHandlersMap"`
	idHash_                           atomic.Int32
	scratchVarExpert_                 InlineVariables_VarExpert
	activeRefMap_                     ReferenceMap
	scratchLocalScope_                NodeIterators_FunctionlessLocalScope
	scratchLocalVarMotion_            NodeIterators_LocalVarMotion
}

func (this_ *InlineVariables_InliningBehavior) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineVariables_InliningBehavior) IsJavaNil_() bool { return this_ == nil }
func (*InlineVariables_InliningBehavior) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables$InliningBehavior", "InliningBehavior"
}
func (this_ *InlineVariables_InliningBehavior) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineVariables_InliningBehavior) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables.InliningBehavior", "com.google.javascript.jscomp.InlineVariables$InliningBehavior", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables_InliningBehavior](v_)
	})
}
func (this_ *InlineVariables_InliningBehavior) SetOuter_(o *InlineVariables)             { this_.outer_ = o }
func (this_ *InlineVariables_InliningBehavior) GetOuter_() *InlineVariables              { return this_.outer_ }
func (*InlineVariables_InliningBehavior) Implements_jscomp_ReferenceCollector_Behavior() {}

func StaticGet_InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS() int32 { return int32(5000) }
func (this_ *InlineVariables_InliningBehavior) initFields_InlineVariables_InliningBehavior() {
	_ = this_
	this_.currentScopeHandledVarAnalysesMap = j2gort.NewOrderedMap[*Var, *InlineVariables_InlineVarAnalysis]()
	this_.varToAliasRetryHandlersMap = SimpleMultimap_Create_0[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler]()
}

// Java constructor: InliningBehavior()
func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	this_ := &InlineVariables_InliningBehavior{}
	this_.outer_ = outer_
	InitInlineVariables_InliningBehavior(this_)
	return this_
}

func InitInlineVariables_InliningBehavior(this_ *InlineVariables_InliningBehavior) {
	this_.initFields_InlineVariables_InliningBehavior()
}

// Java method: InliningBehavior.mode()
func (this_ *InlineVariables_InliningBehavior) Mode() *InlineVariables_Mode {
	return this_.outer_.mode
}

// Java method: InliningBehavior.afterExitScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *InlineVariables_InliningBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	this_.priv_InliningBehavior_doInlinesForScope_3393(t, referenceMap)
}

// Java method: InliningBehavior.doInlinesForScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_doInlinesForScope_3393(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	this_.activeRefMap_ = referenceMap

	this_.currentScopeHandledVarAnalysesMap.Clear()
	var mayBeAParameterModifiedViaArguments bool = this_.priv_InliningBehavior_varsInThisScopeMayBeModifiedUsingArguments_3394(t.GetScope(), referenceMap)
	_ = mayBeAParameterModifiedViaArguments
	{
		cur_6 := (t.GetScope()).iterVars_()
		for cur_6.HasNext() {

			v := cur_6.Next()

			_ = v
			var referenceInfo *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(v)))
			_ = referenceInfo
			var expert *InlineVariables_VarExpert = this_.priv_InliningBehavior_createVarExpert_3396(v, referenceInfo, mayBeAParameterModifiedViaArguments)
			_ = expert
			var analysis *InlineVariables_InlineVarAnalysis = expert.Analyze()
			_ = analysis
			if analysis.ShouldWaitForAliasedVar() {
				var retryHandler *InlineVariables_InliningBehavior_AliasInlineRetryHandler = allocIvPooled_InlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
				_ = retryHandler
				var aliasedVar *Var = analysis.GetAliasedVar()
				_ = aliasedVar
				var aliasedVarAnalysis *InlineVariables_InlineVarAnalysis = this_.currentScopeHandledVarAnalysesMap.Get(aliasedVar)
				_ = aliasedVarAnalysis
				if aliasedVarAnalysis != nil {
					retryHandler.HandleAliasedVarCompletion(aliasedVar, aliasedVarAnalysis)
				} else {
					this_.varToAliasRetryHandlersMap.Put(aliasedVar, retryHandler)
				}
			} else {
				if analysis.ShouldInline() {
					analysis.PerformInline()
				}
				this_.currentScopeHandledVarAnalysesMap.Put(v, analysis)
				this_.priv_InliningBehavior_retryAliases_3395(v, analysis)
			}
		}
	}
	this_.currentScopeHandledVarAnalysesMap.Clear()
	if this_.outer_ != nil {
		this_.outer_.slab_InlineVariables_InlineVarAnalysis.Reset()
	}
	if this_.varToAliasRetryHandlersMap.IsEmpty() {

		if this_.outer_ != nil {
			this_.outer_.slab_InlineVariables_VarExpert.Reset()

			this_.outer_.slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler.Reset()
		}
	}
	this_.activeRefMap_ = nil

}

// Java method: InliningBehavior.varsInThisScopeMayBeModifiedUsingArguments(com.google.javascript.jscomp.Scope scope, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_varsInThisScopeMayBeModifiedUsingArguments_3394(scope *Scope, referenceMap ReferenceMap) bool {
	_ = scope
	_ = referenceMap
	if scope.IsFunctionScope() && !(scope.GetRootNode_0().IsArrowFunction()) {
		var arguments *Var = j2gort.CastPtr[Var](scope.GetArgumentsVar())
		_ = arguments
		var refs *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(arguments)))
		_ = refs
		if (refs != nil) && !(refs.IsEmpty()) {
			{
				var i int32 = int32(0)
				_ = i
				var size int32 = refs.Size()
				_ = size
				for ; i < size; i = i + 1 {
					var ref *Reference = refs.Get(i)
					_ = ref
					if !(this_.IsSafeUseOfArguments(ref.GetNode())) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Java method: InliningBehavior.retryAliases(com.google.javascript.jscomp.Var aliasedVar, com.google.javascript.jscomp.InlineVariables.InlineVarAnalysis aliasedVarAnalysis)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_retryAliases_3395(aliasedVar *Var, aliasedVarAnalysis *InlineVariables_InlineVarAnalysis) {
	_ = aliasedVar
	_ = aliasedVarAnalysis
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.varToAliasRetryHandlersMap.RemoveAll_1_any(j2gort.PtrToAny(aliasedVar))))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			aliasInlineRetryHandler := j2gort.CastPtr[InlineVariables_InliningBehavior_AliasInlineRetryHandler](rawElem_9)
			_ = aliasInlineRetryHandler
			aliasInlineRetryHandler.HandleAliasedVarCompletion(aliasedVar, aliasedVarAnalysis)
		}
	}
}

// Java method: InliningBehavior.createVarExpert(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.ReferenceCollection referenceInfo, boolean mayBeAParameterModifiedViaArguments)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_createVarExpert_3396(v *Var, referenceInfo *ReferenceCollection, mayBeAParameterModifiedViaArguments bool) *InlineVariables_VarExpert {
	_ = v
	_ = referenceInfo
	_ = mayBeAParameterModifiedViaArguments
	if referenceInfo == nil {
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT()
	}
	var isDeclaredOrInferredConstant bool = v.IsDeclaredOrInferredConst()
	_ = isDeclaredOrInferredConstant
	if !(isDeclaredOrInferredConstant) && (this_.outer_.mode == StaticGet_InlineVariables_Mode_CONSTANTS_ONLY()) {
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT()
	}
	if v.IsExtern() {
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT()
	}
	if this_.outer_.compiler.GetCodingConvention().IsExported_2(v.GetName(), v.IsLocal()) {
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT()
	}
	if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(v.GetNameNode()) {
		return StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT()
	}
	return allocIvPooled_InlineVariables_VarExpert_5(this_, v, referenceInfo, isDeclaredOrInferredConstant, mayBeAParameterModifiedViaArguments)
}

// Java method: InliningBehavior.isSafeUseOfArguments(com.google.javascript.rhino.Node argumentsNode)
func (this_ *InlineVariables_InliningBehavior) IsSafeUseOfArguments(argumentsNode pkg_rhino.Node) bool {
	_ = argumentsNode
	if !(argumentsNode.MatchesName_1_string("arguments")) {
		j2gort.CheckArgument(false)
	}
	return (this_.IsTargetOfPropertyRead(argumentsNode) || this_.IsSecondArgumentToDotApplyMethod(argumentsNode))
}

// Java method: InliningBehavior.isTargetOfPropertyRead(com.google.javascript.rhino.Node n)
func (this_ *InlineVariables_InliningBehavior) IsTargetOfPropertyRead(n pkg_rhino.Node) bool {
	_ = n
	var getNode pkg_rhino.Node = n.GetParent()
	_ = getNode
	return ((n.IsFirstChildOf(getNode) && NodeUtil_IsNormalOrOptChainGet(getNode)) && !(NodeUtil_IsLValue(getNode)))
}

// Java method: InliningBehavior.isSecondArgumentToDotApplyMethod(com.google.javascript.rhino.Node n)
func (this_ *InlineVariables_InliningBehavior) IsSecondArgumentToDotApplyMethod(n pkg_rhino.Node) bool {
	_ = n
	var callNode pkg_rhino.Node = n.GetParent()
	_ = callNode
	if NodeUtil_IsNormalOrOptChainCall(callNode) {
		var calleeNode pkg_rhino.Node = callNode.GetFirstChild()
		_ = calleeNode
		if NodeUtil_IsNormalOrOptChainGetProp(calleeNode) {
			if j2gort.CheckNotNullString(calleeNode.GetString()) == "apply" {
				var thisArgNode pkg_rhino.Node = calleeNode.GetNext()
				_ = thisArgNode
				if thisArgNode != 0 {
					return (thisArgNode.GetNext() == n)
				}
			}
		}
	}
	return false
}

// Java method: InliningBehavior.inline(com.google.javascript.jscomp.Reference decl, com.google.javascript.jscomp.Reference init, com.google.javascript.jscomp.Reference ref)
func (this_ *InlineVariables_InliningBehavior) Inline(decl *Reference, init *Reference, ref *Reference) {
	_ = decl
	_ = init
	_ = ref
	var value pkg_rhino.Node = init.GetAssignedValue()
	_ = value
	if !(value != 0) {
		j2gort.CheckState(false)
	}
	var isFunctionDeclaration bool = NodeUtil_IsFunctionDeclaration(value)
	_ = isFunctionDeclaration
	if isFunctionDeclaration {
		var funcName pkg_rhino.Node = value.GetFirstChild()
		_ = funcName
		j2gort.CheckNotNullNode(funcName)
		if !(funcName.IsName()) {
			j2gort.CheckState(false)
		}
		funcName.SetString("")
		this_.outer_.compiler.ReportChangeToChangeScope(value)
		this_.outer_.compiler.ReportChangeToEnclosingScope(value.GetParent())
	}
	this_.priv_InliningBehavior_inlineValue_3403(ref.GetNode(), value.Detach())
	if decl != init {
		var expressRoot pkg_rhino.Node = init.GetGrandparent()
		_ = expressRoot
		if !(expressRoot.IsExprResult()) {
			j2gort.CheckState(false)
		}
		NodeUtil_RemoveChild(expressRoot.GetParent(), expressRoot)
	}
	if !(isFunctionDeclaration) {
		this_.priv_InliningBehavior_removeDeclaration_3402(decl)
	}
}

// Java method: InliningBehavior.inlineWellDefinedVariable(com.google.javascript.jscomp.Var v, com.google.javascript.rhino.Node value, com.google.javascript.jscomp.ReferenceCollection refSet)
func (this_ *InlineVariables_InliningBehavior) InlineWellDefinedVariable(v *Var, value pkg_rhino.Node, refSet *ReferenceCollection) {
	_ = v
	_ = value
	_ = refSet
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = refSet.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var r *Reference = refSet.Get(i)
			_ = r
			if r.GetNode() == v.GetNameNode() {
				this_.priv_InliningBehavior_removeDeclaration_3402(r)
			} else if r.IsSimpleAssignmentToName() {
				this_.priv_InliningBehavior_inlineValue_3403(r.GetParent(), value.Detach())
			} else {
				var clonedValue pkg_rhino.Node = value.CloneTree_0()
				_ = clonedValue
				NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(clonedValue, this_.outer_.compiler)
				this_.priv_InliningBehavior_inlineValue_3403(r.GetNode(), clonedValue)
			}
		}
	}
}

// Java method: InliningBehavior.removeDeclaration(com.google.javascript.jscomp.Reference decl)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_removeDeclaration_3402(decl *Reference) {
	_ = decl
	var varNode pkg_rhino.Node = decl.GetParent()
	_ = varNode
	if !(NodeUtil_IsNameDeclaration(varNode)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(varNode))
	}
	var grandparent pkg_rhino.Node = decl.GetGrandparent()
	_ = grandparent
	this_.outer_.compiler.ReportChangeToEnclosingScope(decl.GetNode())
	decl.GetNode().Detach()
	if !(varNode.HasChildren()) {
		NodeUtil_RemoveChild(grandparent, varNode)
	}
}

// Java method: InliningBehavior.inlineValue(com.google.javascript.rhino.Node toRemove, com.google.javascript.rhino.Node toInsert)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_inlineValue_3403(toRemove pkg_rhino.Node, toInsert pkg_rhino.Node) {
	_ = toRemove
	_ = toInsert
	this_.outer_.compiler.ReportChangeToEnclosingScope(toRemove)
	if (toRemove.GetColor() != nil) && toRemove.IsColorFromTypeCast() {
		toInsert.SetColor(toRemove.GetColor())
		toInsert.SetColorFromTypeCast()
	}
	toRemove.ReplaceWith(toInsert)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(toRemove, this_.outer_.compiler)
}

// Java method: InliningBehavior.canInline(com.google.javascript.jscomp.Reference declaration, com.google.javascript.jscomp.Reference initialization, com.google.javascript.jscomp.Reference reference, com.google.javascript.rhino.Node initValue)
func (this_ *InlineVariables_InliningBehavior) CanInline(declaration *Reference, initialization *Reference, reference *Reference, initValue pkg_rhino.Node) bool {
	_ = declaration
	_ = initialization
	_ = reference
	_ = initValue
	if (declaration != initialization) && !(initialization.GetGrandparent().IsExprResult()) {
		return false
	}
	if (declaration.GetBasicBlock() != initialization.GetBasicBlock()) || (declaration.GetBasicBlock() != reference.GetBasicBlock()) {
		return false
	}
	if !(initValue != 0) {
		j2gort.CheckState(false)
	}
	if (initValue.IsGetProp() && reference.GetParent().IsCall()) && (reference.GetParent().GetFirstChild() == reference.GetNode()) {
		return false
	}
	if initValue.IsFunction() {
		var callNode pkg_rhino.Node = reference.GetParent()
		_ = callNode
		if reference.GetParent().IsCall() {
			var convention CodingConvention = this_.outer_.compiler.GetCodingConvention()
			_ = convention
			var relationship *CodingConvention_SubclassRelationship = convention.GetClassesDefinedByCall(callNode)
			_ = relationship
			if relationship != nil {
				return false
			}
			if convention.GetSingletonGetterClassName(callNode) != j2gort.NullString {
				return false
			}
		}
	}
	if (initialization.GetScope() != declaration.GetScope()) || !(initialization.GetScope().Contains(ToIface_AbstractScope(reference.GetScope()))) {
		return false
	}
	return (this_.priv_InliningBehavior_canMoveAggressively_3405(initValue) || this_.priv_InliningBehavior_canMoveModerately_3406(initialization, reference))
}

// Java method: InliningBehavior.canMoveAggressively(com.google.javascript.rhino.Node value)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_canMoveAggressively_3405(value pkg_rhino.Node) bool {
	_ = value
	return (NodeUtil_IsLiteralValue(value, true) || value.IsFunction())
}

// Java method: InliningBehavior.canMoveModerately(com.google.javascript.jscomp.Reference initialization, com.google.javascript.jscomp.Reference reference)
func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_canMoveModerately_3406(initialization *Reference, reference *Reference) bool {
	_ = initialization
	_ = reference
	var it *NodeIterators_LocalVarMotion = nil
	_ = it
	if NodeUtil_IsNameDeclaration(initialization.GetParent()) {
		it = allocReuseVia_NodeIterators_LocalVarMotion_ForVar(&this_.scratchLocalVarMotion_, &this_.scratchLocalScope_, this_.outer_.compiler, initialization.GetNode(), initialization.GetParent(), initialization.GetGrandparent())
	} else if initialization.GetParent().IsAssign() {
		if !(initialization.GetGrandparent().IsExprResult()) {
			j2gort.CheckState(false)
		}
		it = allocReuseVia_NodeIterators_LocalVarMotion_ForAssign(&this_.scratchLocalVarMotion_, &this_.scratchLocalScope_, this_.outer_.compiler, initialization.GetNode(), initialization.GetParent(), initialization.GetGrandparent(), initialization.GetGrandparent().GetParent())
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected initialization parent\n"+j2gort.StrValString(initialization.GetParent().ToStringTree()))))
	}
	var targetName pkg_rhino.Node = reference.GetNode()
	_ = targetName
	var steps int32 = int32(0)
	_ = steps
	for it.HasNext() {
		var curNode pkg_rhino.Node = it.NextNode()
		_ = curNode
		if curNode == targetName {
			return true
		}
		{
			var iifeCond_1 int32
			{

				old_10 := steps
				nw_11 := old_10 + 1
				steps = nw_11
				iifeCond_1 = nw_11
			}
			if iifeCond_1 > int32(5000) {
				return false
			}
		}
	}
	return false
}

// Java method: InliningBehavior.isValidDeclaration(com.google.javascript.jscomp.Reference declaration)
func (this_ *InlineVariables_InliningBehavior) IsValidDeclaration(declaration *Reference) bool {
	_ = declaration
	return ((NodeUtil_IsNameDeclaration(declaration.GetParent()) && !(NodeUtil_IsLoopStructure(declaration.GetGrandparent()))) || NodeUtil_IsFunctionDeclaration(declaration.GetParent()))
}

// Java method: InliningBehavior.isValidInitialization(com.google.javascript.jscomp.Reference initialization)
func (this_ *InlineVariables_InliningBehavior) IsValidInitialization(initialization *Reference) bool {
	_ = initialization
	if initialization == nil {
		return false
	} else if initialization.IsDeclaration() {
		if !(NodeUtil_IsFunctionDeclaration(initialization.GetParent())) && !(initialization.GetNode().HasChildren()) {
			return false
		}
	} else {
		var parent pkg_rhino.Node = initialization.GetParent()
		_ = parent
		if !(parent.IsAssign() && (parent.GetFirstChild() == initialization.GetNode())) {
			j2gort.CheckState(false)
		}
	}
	return true
}

// Java method: InliningBehavior.isValidReference(com.google.javascript.jscomp.Reference reference)
func (this_ *InlineVariables_InliningBehavior) IsValidReference(reference *Reference) bool {
	_ = reference
	return (!(reference.IsDeclaration()) && !(reference.IsLvalue()))
}

func (this_ *InlineVariables_InliningBehavior) BeforeTraverseNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	IfaceDefault_ReferenceCollector_Behavior_BeforeTraverseNode(this_, p0_, p1_, p2_)
}

func (this_ *InlineVariables_InliningBehavior) AfterVisitNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node, p3_ *Reference) {
	IfaceDefault_ReferenceCollector_Behavior_AfterVisitNode(this_, p0_, p1_, p2_, p3_)
}

// Java class: com.google.javascript.jscomp.InlineVariables.InliningBehavior.AliasInlineRetryHandler
type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer_  *InlineVariables_InliningBehavior
	v       *Var                       `gson:"v"`
	expert  *InlineVariables_VarExpert `gson:"expert"`
	idHash_ atomic.Int32
}

func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) IsJavaNil_() bool {
	return this_ == nil
}
func (*InlineVariables_InliningBehavior_AliasInlineRetryHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineVariables$InliningBehavior$AliasInlineRetryHandler", "AliasInlineRetryHandler"
}
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineVariables_InliningBehavior_AliasInlineRetryHandler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineVariables.InliningBehavior.AliasInlineRetryHandler", "com.google.javascript.jscomp.InlineVariables$InliningBehavior$AliasInlineRetryHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineVariables_InliningBehavior_AliasInlineRetryHandler](v_)
	})
}
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) SetOuter_(o *InlineVariables_InliningBehavior) {
	this_.outer_ = o
}
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) GetOuter_() *InlineVariables_InliningBehavior {
	return this_.outer_
}

func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) initFields_InlineVariables_InliningBehavior_AliasInlineRetryHandler() {
	_ = this_
}

// Java constructor: AliasInlineRetryHandler(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.InlineVariables.VarExpert expert)
func NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(outer_ *InlineVariables_InliningBehavior, v *Var, expert *InlineVariables_VarExpert) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	this_ := &InlineVariables_InliningBehavior_AliasInlineRetryHandler{}
	this_.outer_ = outer_
	InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
	return this_
}
func allocIvPooled_InlineVariables_InliningBehavior_AliasInlineRetryHandler(outer_ *InlineVariables_InliningBehavior, v *Var, expert *InlineVariables_VarExpert) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	expert = promoteIvScratchVarExpert_(outer_, expert)
	if outer_ == nil || outer_.outer_ == nil {
		return NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(outer_, v, expert)
	}
	this_ := outer_.outer_.slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler.Alloc()
	this_.outer_ = outer_
	InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
	return this_
}

func InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler, v *Var, expert *InlineVariables_VarExpert) {
	_ = v
	_ = expert
	this_.initFields_InlineVariables_InliningBehavior_AliasInlineRetryHandler()
	this_.v = v
	this_.expert = expert
}

// Java method: AliasInlineRetryHandler.handleAliasedVarCompletion(com.google.javascript.jscomp.Var aliasedVar, com.google.javascript.jscomp.InlineVariables.InlineVarAnalysis aliasedVarAnalyis)
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) HandleAliasedVarCompletion(aliasedVar *Var, aliasedVarAnalyis *InlineVariables_InlineVarAnalysis) {
	_ = aliasedVar
	_ = aliasedVarAnalyis
	var newAnalysis *InlineVariables_InlineVarAnalysis = this_.expert.ReanalyzeAfterAliasedVar(aliasedVar, aliasedVarAnalyis)
	_ = newAnalysis
	if !(!(newAnalysis.ShouldWaitForAliasedVar())) {
		j2gort.CheckStateFmt(false, "expert for %s asked to wait a second time", j2gort.PtrToAny(this_.v))
	}
	if newAnalysis.ShouldInline() {
		newAnalysis.PerformInline()
	}
	this_.outer_.currentScopeHandledVarAnalysesMap.Put(this_.v, newAnalysis)
	this_.outer_.priv_InliningBehavior_retryAliases_3395(this_.v, newAnalysis)
}
func (this_ *InlineVariables) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_InlineVariables_VarExpert.Release()
	this_.slab_InlineVariables_InlineVarAnalysis.Release()
	this_.slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler.Release()
}

var (
	staticIvLazyNegUnknown_      = InlineVariables_InlineVarAnalysis{kind: int8(2), lazyWellDefined_: int8(-1), lazyAssignedOnce_: int8(-1)}
	staticIvLazyNegWellDef_      = InlineVariables_InlineVarAnalysis{kind: int8(2), lazyWellDefined_: int8(1), lazyAssignedOnce_: int8(-1)}
	staticIvLazyNegAssignedOnce_ = InlineVariables_InlineVarAnalysis{kind: int8(2), lazyWellDefined_: int8(-1), lazyAssignedOnce_: int8(1)}
)

func isSafeToInlineAliasesIv_(beh_ *InlineVariables_InliningBehavior, aliasedVar *Var, analysis *InlineVariables_InlineVarAnalysis) bool {
	if analysis == &staticIvLazyNegUnknown_ || analysis == &staticIvLazyNegWellDef_ || analysis == &staticIvLazyNegAssignedOnce_ {
		if beh_ == nil || beh_.activeRefMap_ == nil || aliasedVar == nil {
			return false
		}
		var rc *ReferenceCollection
		if rcCol, ok := beh_.activeRefMap_.(*ReferenceCollector); ok && rcCol != nil {
			rc = rcCol.referenceMap.Get(aliasedVar)
		} else {
			rc, _ = beh_.activeRefMap_.GetReferences(j2gort.PtrToAny(aliasedVar)).(*ReferenceCollection)
		}
		if rc == nil {
			return false
		}
		if rc.ivSafeAliases_ != 0 {
			return rc.ivSafeAliases_ == 2
		}
		ok := (analysis == &staticIvLazyNegWellDef_ || rc.IsWellDefined()) &&
			(analysis == &staticIvLazyNegAssignedOnce_ || rc.IsAssignedOnceInLifetime())
		if ok {
			rc.ivSafeAliases_ = 2
		} else {
			rc.ivSafeAliases_ = 1
		}
		return ok
	}
	return analysis.IsSafeToInlineAliases()
}
