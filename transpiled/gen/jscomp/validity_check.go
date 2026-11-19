/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ValidityCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ValidityCheck
type ValidityCheck struct {
	compiler      AbstractCompiler `gson:"compiler"`
	changeTracker *ChangeTracker   `gson:"changeTracker"`
	astValidator  *AstValidator    `gson:"astValidator"`
	idHash_       atomic.Int32
}

func (this_ *ValidityCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *ValidityCheck) IsJavaNil_() bool { return this_ == nil }
func (*ValidityCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ValidityCheck", "ValidityCheck"
}
func (this_ *ValidityCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ValidityCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ValidityCheck", "com.google.javascript.jscomp.ValidityCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*ValidityCheck](v_)
	})
}
func (*ValidityCheck) Implements_jscomp_CompilerPass() {}

var static_ValidityCheck_EXTERN_PROPERTIES_CHANGED *DiagnosticType

func StaticGet_ValidityCheck_EXTERN_PROPERTIES_CHANGED() *DiagnosticType {
	clinit_ValidityCheck()
	return static_ValidityCheck_EXTERN_PROPERTIES_CHANGED
}
func StaticSet_ValidityCheck_EXTERN_PROPERTIES_CHANGED(v *DiagnosticType) {
	clinit_ValidityCheck()
	static_ValidityCheck_EXTERN_PROPERTIES_CHANGED = v
}

var clinitDone_ValidityCheck uint32
var clinitState_ValidityCheck j2gort.ClinitState

func clinit_ValidityCheck() {
	if atomic.LoadUint32(&clinitDone_ValidityCheck) != 0 {
		return
	}
	clinitSlow_ValidityCheck()
}

func clinitSlow_ValidityCheck() {
	if !j2gort.EnterClinit(&clinitState_ValidityCheck, &clinitDone_ValidityCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ValidityCheck, &clinitDone_ValidityCheck)
	static_ValidityCheck_EXTERN_PROPERTIES_CHANGED = DiagnosticType_Error("JSC_EXTERN_PROPERTIES_CHANGED", "Internal compiler error. Extern properties modified from:\n{0}\nto:\n{1}")
}

func (this_ *ValidityCheck) initFields_ValidityCheck() {
	_ = this_
}

// Java constructor: ValidityCheck(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewValidityCheck(compiler AbstractCompiler) *ValidityCheck {
	clinit_ValidityCheck()
	this_ := &ValidityCheck{}
	InitValidityCheck(this_, compiler)
	return this_
}

func InitValidityCheck(this_ *ValidityCheck, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ValidityCheck()
	this_.compiler = compiler
	this_.changeTracker = compiler.GetChangeTracker()
	this_.astValidator = NewAstValidator_1(compiler)
}

// Java method: ValidityCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ValidityCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.priv_ValidityCheck_checkAst_5873(externs, root)
	this_.priv_ValidityCheck_checkNormalization_5875(externs, root)
	this_.priv_ValidityCheck_checkVars_5874(externs, root)
	this_.priv_ValidityCheck_checkExternProperties_5876(externs)
}

// Java method: ValidityCheck.checkAst(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ValidityCheck) priv_ValidityCheck_checkAst_5873(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.astValidator.ValidateCodeRoot(externs)
	this_.astValidator.ValidateCodeRoot(root)
}

// Java method: ValidityCheck.checkVars(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ValidityCheck) priv_ValidityCheck_checkVars_5874(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.compiler.GetLifeCycleStage().IsNormalized() {
		{
			clinit_VarCheck()
			var val_1 VarCheck
			InitVarCheck_2(&val_1, this_.compiler, true)
			(&val_1).Process_2(externs, root)
		}
	}
}

// Java method: ValidityCheck.checkNormalization(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ValidityCheck) priv_ValidityCheck_checkNormalization_5875(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var handler CodeChangeHandler = NewForbiddenChange()
	_ = handler
	this_.changeTracker.AddChangeHandler(handler)
	if this_.compiler.GetLifeCycleStage().IsNormalized() {
		Normalize_S_Builder(this_.compiler).AssertOnChange(true).Build().Process_2(externs, root)
		if this_.compiler.GetLifeCycleStage().IsNormalizedUnobfuscated() {
			var checkUserDeclarations bool = true
			_ = checkUserDeclarations
			var pass CompilerPass = NewValidityCheck_VerifyConstants(this_.compiler, checkUserDeclarations)
			_ = pass
			pass.Process_2(externs, root)
		}
	}
	this_.changeTracker.RemoveChangeHandler(handler)
}

// Java method: ValidityCheck.checkExternProperties(com.google.javascript.rhino.Node externs)
func (this_ *ValidityCheck) priv_ValidityCheck_checkExternProperties_5876(externs pkg_rhino.Node) {
	_ = externs
	var externProperties j2gort.Ext_com_google_common_collect_ImmutableSet = this_.compiler.GetExternProperties()
	_ = externProperties
	if externProperties == nil {
		return
	}
	{
		var val_2 GatherExternProperties
		InitGatherExternProperties(&val_2, this_.compiler, StaticGet_GatherExternProperties_Mode_CHECK())
		(&val_2).Process_2(externs, pkg_rhino.Node(0))
	}
	if !(this_.compiler.GetExternProperties().Equals(externProperties)) {
		this_.compiler.Report_1(JSError_Make_2(StaticGet_ValidityCheck_EXTERN_PROPERTIES_CHANGED(), []string{externProperties.String(), this_.compiler.GetExternProperties().String()}))
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Validity Check failed: Extern properties changed from:\n" + j2gort.StrVal(externProperties) + "\nto:\n" + j2gort.StrVal(this_.compiler.GetExternProperties()))))))
	}
}

// Java class: com.google.javascript.jscomp.ValidityCheck.VerifyConstants
type ValidityCheck_VerifyConstants struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler              AbstractCompiler                              `gson:"compiler"`
	constantMap           j2gort.OrderedMap[string, *j2gort.BoxBoolean] `gson:"constantMap"`
	checkUserDeclarations bool                                          `gson:"checkUserDeclarations"`
}

func (this_ *ValidityCheck_VerifyConstants) BindSelf_(s any) {
	_ = s
}
func (this_ *ValidityCheck_VerifyConstants) IsJavaNil_() bool { return this_ == nil }
func (*ValidityCheck_VerifyConstants) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ValidityCheck$VerifyConstants", "VerifyConstants"
}
func (this_ *ValidityCheck_VerifyConstants) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ValidityCheck_VerifyConstants) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ValidityCheck.VerifyConstants", "com.google.javascript.jscomp.ValidityCheck$VerifyConstants", func(v_ any) bool {
		return j2gort.InstanceOf[*ValidityCheck_VerifyConstants](v_)
	})
}
func (*ValidityCheck_VerifyConstants) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ValidityCheck_VerifyConstants) Implements_jscomp_CompilerPass()                            {}
func (*ValidityCheck_VerifyConstants) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ValidityCheck_VerifyConstants) initFields_ValidityCheck_VerifyConstants() {
	_ = this_
	this_.constantMap = j2gort.NewOrderedMap[string, *j2gort.BoxBoolean]()
}

// Java constructor: VerifyConstants(com.google.javascript.jscomp.AbstractCompiler compiler, boolean checkUserDeclarations)
func NewValidityCheck_VerifyConstants(compiler AbstractCompiler, checkUserDeclarations bool) *ValidityCheck_VerifyConstants {
	this_ := &ValidityCheck_VerifyConstants{}
	InitValidityCheck_VerifyConstants(this_, compiler, checkUserDeclarations)
	return this_
}

func InitValidityCheck_VerifyConstants(this_ *ValidityCheck_VerifyConstants, compiler AbstractCompiler, checkUserDeclarations bool) {
	_ = compiler
	_ = checkUserDeclarations
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ValidityCheck_VerifyConstants()
	this_.compiler = compiler
	this_.checkUserDeclarations = checkUserDeclarations
}

// Java method: VerifyConstants.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ValidityCheck_VerifyConstants) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var externsAndJs pkg_rhino.Node = root.GetParent()
	_ = externsAndJs
	if !(externsAndJs != 0) {
		j2gort.CheckState(false)
	}
	if !(externsAndJs.HasChild(externs)) {
		j2gort.CheckState(false)
	}
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}

// Java method: VerifyConstants.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ValidityCheck_VerifyConstants) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if j2gort.CheckNotNullString(n.GetString()) == "" {
			return
		}
		var isConst bool = n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())
		_ = isConst
		if this_.checkUserDeclarations {
			var expectedConst bool = false
			_ = expectedConst
			var convention CodingConvention = this_.compiler.GetCodingConvention()
			_ = convention
			if NodeUtil_IsConstantName(n) || NodeUtil_IsConstantByConvention(convention, n) {
				expectedConst = true
			} else {
				expectedConst = false
				var info *pkg_rhino.JSDocInfo = nil
				_ = info
				var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
				_ = var_
				if var_ != nil {
					info = var_.GetJSDocInfo()
				}
				if (info != nil) && info.IsConstant() {
					expectedConst = true
				} else {
					expectedConst = false
				}
			}
			if expectedConst {
				if !(expectedConst == isConst) {
					j2gort.CheckStateFmt(false, "The name %s is not annotated as constant.", j2gort.BoxString(name))
				}
			} else {
				if !(expectedConst == isConst) {
					j2gort.CheckStateFmt(false, "The name %s should not be annotated as constant.", j2gort.BoxString(name))
				}
			}
		}
		var value *j2gort.BoxBoolean = this_.constantMap.Get(name)
		_ = value
		if value == nil {
			this_.constantMap.Put(name, j2gort.ValueOfBoolean(isConst))
		} else if value.BooleanValue() != isConst {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("The name " + j2gort.StrValString(name) + " is not consistently annotated as constant. Current constness: " + j2gort.StrValBool(isConst) + ", previous constness: " + j2gort.StrValBool(value.BooleanValue()) + "\n\nAt node: " + j2gort.StrVal(n))))))
		}
	}
}
func (*ValidityCheck_VerifyConstants) AlwaysShouldTraverse_()      {}
func (*ValidityCheck_VerifyConstants) CanReleaseSyntacticScopes_() {}
