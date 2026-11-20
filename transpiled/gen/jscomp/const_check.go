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
// Port of com/google/javascript/jscomp/ConstCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ConstCheck
type ConstCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler             AbstractCompiler              `gson:"compiler"`
	initializedConstants j2gort.CustomOrderedSet[*Var] `gson:"initializedConstants"`
}

func (this_ *ConstCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *ConstCheck) IsJavaNil_() bool { return this_ == nil }
func (*ConstCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConstCheck", "ConstCheck"
}
func (this_ *ConstCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ConstCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConstCheck", "com.google.javascript.jscomp.ConstCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*ConstCheck](v_)
	})
}
func (*ConstCheck) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ConstCheck) Implements_jscomp_CompilerPass()                            {}
func (*ConstCheck) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ConstCheck_CONST_REASSIGNED_VALUE_ERROR *DiagnosticType

func StaticGet_ConstCheck_CONST_REASSIGNED_VALUE_ERROR() *DiagnosticType {
	clinit_ConstCheck()
	return static_ConstCheck_CONST_REASSIGNED_VALUE_ERROR
}
func StaticSet_ConstCheck_CONST_REASSIGNED_VALUE_ERROR(v *DiagnosticType) {
	clinit_ConstCheck()
	static_ConstCheck_CONST_REASSIGNED_VALUE_ERROR = v
}

var clinitDone_ConstCheck uint32
var clinitState_ConstCheck j2gort.ClinitState

func clinit_ConstCheck() {
	if atomic.LoadUint32(&clinitDone_ConstCheck) != 0 {
		return
	}
	clinitSlow_ConstCheck()
}

func clinitSlow_ConstCheck() {
	if !j2gort.EnterClinit(&clinitState_ConstCheck, &clinitDone_ConstCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConstCheck, &clinitDone_ConstCheck)
	static_ConstCheck_CONST_REASSIGNED_VALUE_ERROR = DiagnosticType_Warning("JSC_CONSTANT_REASSIGNED_VALUE_ERROR", "constant {0} assigned a value more than once.\nOriginal definition at {1}")
}

func (this_ *ConstCheck) initFields_ConstCheck() {
	_ = this_
}

// Java constructor: ConstCheck(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewConstCheck(compiler AbstractCompiler) *ConstCheck {
	clinit_ConstCheck()
	this_ := &ConstCheck{}
	InitConstCheck(this_, compiler)
	return this_
}

func InitConstCheck(this_ *ConstCheck, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ConstCheck()
	this_.compiler = compiler
	this_.initializedConstants = j2gort.NewCustomOrderedSet[*Var]()
}

// Java method: ConstCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ConstCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		this_.initializedConstants.Clear()
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).TraverseRoots(externs, root)
	}

}

// Java method: ConstCheck.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConstCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		if NodeUtil_IsNameDeclaration(parent) {
			var name string = n.GetString()
			_ = name
			var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
			_ = var_
			if this_.priv_ConstCheck_isConstant_2178(var_, n) {
				if n.IsFromExterns() {
					this_.initializedConstants.Add(var_)
				} else if n.HasChildren() {
					if !(this_.initializedConstants.Add(var_)) {
						this_.priv_ConstCheck_reportError_2179(n, var_, name)
					}
				}
			}
		}
	case int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(63) /* Token.ASSIGN_EXPONENT */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		if lhs.IsName() {
			var name_2 string = lhs.GetString()
			_ = name_2
			var var__2 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name_2))
			_ = var__2
			if this_.priv_ConstCheck_isConstant_2178(var__2, lhs) && !(this_.initializedConstants.Add(var__2)) {
				this_.priv_ConstCheck_reportError_2179(n, var__2, name_2)
			} else if ((var__2 != nil) && var__2.IsGoogModuleExports()) && !(this_.initializedConstants.Add(var__2)) {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ConstCheck_CONST_REASSIGNED_VALUE_ERROR(), []string{"exports", n.GetSourceFileName()}))
			}
		}
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		var lhs_2 pkg_rhino.Node = n.GetFirstChild()
		_ = lhs_2
		if lhs_2.IsName() {
			var name_3 string = lhs_2.GetString()
			_ = name_3
			var var__3 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name_3))
			_ = var__3
			if this_.priv_ConstCheck_isConstant_2178(var__3, lhs_2) {
				this_.priv_ConstCheck_reportError_2179(n, var__3, name_3)
			}
		}
	default:
	}
}

// Java method: ConstCheck.isConstant(com.google.javascript.jscomp.Var var, com.google.javascript.rhino.Node nameNode)
func (this_ *ConstCheck) priv_ConstCheck_isConstant_2178(var_ *Var, nameNode pkg_rhino.Node) bool {
	_ = var_
	_ = nameNode
	j2gort.CheckNotNullPtr(var_)
	return (var_.IsConst() || var_.IsDeclaredOrInferredConst())
}

// Java method: ConstCheck.reportError(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var var, java.lang.String name)
func (this_ *ConstCheck) priv_ConstCheck_reportError_2179(n pkg_rhino.Node, var_ *Var, name string) {
	_ = n
	_ = var_
	_ = name
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	if (info == nil) || !(j2gort.SetContainsStr(info.GetSuppressions(), "const")) {
		var declNode pkg_rhino.Node = var_.GetNode()
		_ = declNode
		var declaredPosition string = (j2gort.StrValString(declNode.GetSourceFileName()) + ":" + j2gort.StrValInt32(declNode.GetLineno()))
		_ = declaredPosition
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ConstCheck_CONST_REASSIGNED_VALUE_ERROR(), []string{name, declaredPosition}))
	}
}
func (*ConstCheck) AlwaysShouldTraverse_() {}
