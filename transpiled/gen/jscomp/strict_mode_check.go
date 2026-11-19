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
// Port of com/google/javascript/jscomp/StrictModeCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.StrictModeCheck
type StrictModeCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler     AbstractCompiler `gson:"compiler"`
	defaultLevel *CheckLevel      `gson:"defaultLevel"`
}

func (this_ *StrictModeCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *StrictModeCheck) IsJavaNil_() bool { return this_ == nil }
func (*StrictModeCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.StrictModeCheck", "StrictModeCheck"
}
func (this_ *StrictModeCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*StrictModeCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StrictModeCheck", "com.google.javascript.jscomp.StrictModeCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*StrictModeCheck](v_)
	})
}
func (*StrictModeCheck) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*StrictModeCheck) Implements_jscomp_CompilerPass()                            {}
func (*StrictModeCheck) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_StrictModeCheck_USE_OF_WITH *DiagnosticType

func StaticGet_StrictModeCheck_USE_OF_WITH() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_USE_OF_WITH
}
func StaticSet_StrictModeCheck_USE_OF_WITH(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_USE_OF_WITH = v
}

var static_StrictModeCheck_EVAL_DECLARATION *DiagnosticType

func StaticGet_StrictModeCheck_EVAL_DECLARATION() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_EVAL_DECLARATION
}
func StaticSet_StrictModeCheck_EVAL_DECLARATION(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_EVAL_DECLARATION = v
}

var static_StrictModeCheck_EVAL_ASSIGNMENT *DiagnosticType

func StaticGet_StrictModeCheck_EVAL_ASSIGNMENT() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_EVAL_ASSIGNMENT
}
func StaticSet_StrictModeCheck_EVAL_ASSIGNMENT(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_EVAL_ASSIGNMENT = v
}

var static_StrictModeCheck_ARGUMENTS_DECLARATION *DiagnosticType

func StaticGet_StrictModeCheck_ARGUMENTS_DECLARATION() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_ARGUMENTS_DECLARATION
}
func StaticSet_StrictModeCheck_ARGUMENTS_DECLARATION(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_ARGUMENTS_DECLARATION = v
}

var static_StrictModeCheck_ARGUMENTS_ASSIGNMENT *DiagnosticType

func StaticGet_StrictModeCheck_ARGUMENTS_ASSIGNMENT() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_ARGUMENTS_ASSIGNMENT
}
func StaticSet_StrictModeCheck_ARGUMENTS_ASSIGNMENT(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_ARGUMENTS_ASSIGNMENT = v
}

var static_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN *DiagnosticType

func StaticGet_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN
}
func StaticSet_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN = v
}

var static_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN *DiagnosticType

func StaticGet_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN
}
func StaticSet_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN = v
}

var static_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN *DiagnosticType

func StaticGet_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN
}
func StaticSet_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN = v
}

var static_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN *DiagnosticType

func StaticGet_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN
}
func StaticSet_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN = v
}

var static_StrictModeCheck_DELETE_VARIABLE *DiagnosticType

func StaticGet_StrictModeCheck_DELETE_VARIABLE() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_DELETE_VARIABLE
}
func StaticSet_StrictModeCheck_DELETE_VARIABLE(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_DELETE_VARIABLE = v
}

var static_StrictModeCheck_DUPLICATE_MEMBER *DiagnosticType

func StaticGet_StrictModeCheck_DUPLICATE_MEMBER() *DiagnosticType {
	clinit_StrictModeCheck()
	return static_StrictModeCheck_DUPLICATE_MEMBER
}
func StaticSet_StrictModeCheck_DUPLICATE_MEMBER(v *DiagnosticType) {
	clinit_StrictModeCheck()
	static_StrictModeCheck_DUPLICATE_MEMBER = v
}

var clinitDone_StrictModeCheck uint32
var clinitState_StrictModeCheck j2gort.ClinitState

func clinit_StrictModeCheck() {
	if atomic.LoadUint32(&clinitDone_StrictModeCheck) != 0 {
		return
	}
	clinitSlow_StrictModeCheck()
}

func clinitSlow_StrictModeCheck() {
	if !j2gort.EnterClinit(&clinitState_StrictModeCheck, &clinitDone_StrictModeCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_StrictModeCheck, &clinitDone_StrictModeCheck)
	static_StrictModeCheck_USE_OF_WITH = DiagnosticType_Error("JSC_USE_OF_WITH", "The 'with' statement cannot be used in strict mode.")
	static_StrictModeCheck_EVAL_DECLARATION = DiagnosticType_Error("JSC_EVAL_DECLARATION", "\"eval\" cannot be redeclared in strict mode")
	static_StrictModeCheck_EVAL_ASSIGNMENT = DiagnosticType_Error("JSC_EVAL_ASSIGNMENT", "the \"eval\" object cannot be reassigned in strict mode")
	static_StrictModeCheck_ARGUMENTS_DECLARATION = DiagnosticType_Error("JSC_ARGUMENTS_DECLARATION", "\"arguments\" cannot be redeclared in strict mode")
	static_StrictModeCheck_ARGUMENTS_ASSIGNMENT = DiagnosticType_Error("JSC_ARGUMENTS_ASSIGNMENT", "the \"arguments\" object cannot be reassigned in strict mode")
	static_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN = DiagnosticType_Error("JSC_ARGUMENTS_CALLEE_FORBIDDEN", "\"arguments.callee\" cannot be used in strict mode")
	static_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN = DiagnosticType_Error("JSC_ARGUMENTS_CALLER_FORBIDDEN", "\"arguments.caller\" cannot be used in strict mode")
	static_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN = DiagnosticType_Error("JSC_FUNCTION_CALLER_FORBIDDEN", "A function''s \"caller\" property cannot be used in strict mode")
	static_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN = DiagnosticType_Error("JSC_FUNCTION_ARGUMENTS_PROP_FORBIDDEN", "A function''s \"arguments\" property cannot be used in strict mode")
	static_StrictModeCheck_DELETE_VARIABLE = DiagnosticType_Error("JSC_DELETE_VARIABLE", "variables, functions, and arguments cannot be deleted in strict mode")
	static_StrictModeCheck_DUPLICATE_MEMBER = DiagnosticType_Warning("JSC_DUPLICATE_MEMBER", "Class or object literal contains duplicate member \"{0}\". In non-strict code, the last duplicate will overwrite the others.")
}

func (this_ *StrictModeCheck) initFields_StrictModeCheck() {
	_ = this_
}

// Java constructor: StrictModeCheck(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.CheckLevel defaultLevel)
func NewStrictModeCheck(compiler AbstractCompiler, defaultLevel *CheckLevel) *StrictModeCheck {
	clinit_StrictModeCheck()
	this_ := &StrictModeCheck{}
	InitStrictModeCheck(this_, compiler, defaultLevel)
	return this_
}

func InitStrictModeCheck(this_ *StrictModeCheck, compiler AbstractCompiler, defaultLevel *CheckLevel) {
	_ = compiler
	_ = defaultLevel
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_StrictModeCheck()
	this_.compiler = compiler
	this_.defaultLevel = defaultLevel
}

// Java method: StrictModeCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *StrictModeCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
	NodeTraversal_Traverse(this_.compiler, root, NewStrictModeCheck_NonExternChecks(this_))
}

// Java method: StrictModeCheck.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StrictModeCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsAssign() {
		this_.priv_StrictModeCheck_checkAssignment_5187(n)
	} else if n.IsDelProp() {
		this_.priv_StrictModeCheck_checkDelete_5188(t, n)
	} else if n.IsObjectLit() {
		this_.priv_StrictModeCheck_checkObjectLiteralOrClass_5189(n)
	} else if n.IsClass() {
		this_.priv_StrictModeCheck_checkObjectLiteralOrClass_5189(n.GetLastChild())
	} else if n.IsWith() {
		this_.priv_StrictModeCheck_checkWith_5186(n)
	}
}

// Java method: StrictModeCheck.checkWith(com.google.javascript.rhino.Node n)
func (this_ *StrictModeCheck) priv_StrictModeCheck_checkWith_5186(n pkg_rhino.Node) {
	_ = n
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	var allowWith bool = ((info != nil) && j2gort.SetContainsStr(info.GetSuppressions(), "with"))
	_ = allowWith
	if !(allowWith) {
		this_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_USE_OF_WITH(), []string{})
	}
}

// Java static method: StrictModeCheck.isDeclaration(com.google.javascript.rhino.Node n)
func StrictModeCheck_IsDeclaration(n pkg_rhino.Node) bool {
	clinit_StrictModeCheck()
	_ = n
	switch n.GetParent().TokenSwitchOrd_() {
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */, int32(86) /* Token.CATCH */ :
		return true
	case int32(73) /* Token.FUNCTION */ :
		return (n == n.GetParent().GetFirstChild())
	case int32(49) /* Token.PARAM_LIST */ :
		return n.GetGrandparent().IsFunction()
	default:
		return false
	}
	panic("unreachable")
}

// Java method: StrictModeCheck.checkAssignment(com.google.javascript.rhino.Node n)
func (this_ *StrictModeCheck) priv_StrictModeCheck_checkAssignment_5187(n pkg_rhino.Node) {
	_ = n
	if n.GetFirstChild().IsName() {
		{
			var iifeCond_2 bool
			{

				sr_1 := "arguments"
				se_2 := n.GetFirstChild().GetString()
				iifeCond_2 = j2gort.CheckNotNullString(sr_1) == se_2
			}
			if iifeCond_2 {
				this_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_ARGUMENTS_ASSIGNMENT(), []string{})
			} else {
				var iifeCond_1 bool
				{

					sr_3 := "eval"
					se_4 := n.GetFirstChild().GetString()
					iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
				}
				if iifeCond_1 {
					this_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_EVAL_ASSIGNMENT(), []string{})
				}
			}
		}
	}
}

// Java method: StrictModeCheck.checkDelete(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *StrictModeCheck) priv_StrictModeCheck_checkDelete_5188(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.GetFirstChild().IsName() {
		var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetFirstChild().GetString()))
		_ = v
		if v != nil {
			this_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_DELETE_VARIABLE(), []string{})
		}
	}
}

// Java method: StrictModeCheck.checkObjectLiteralOrClass(com.google.javascript.rhino.Node n)
func (this_ *StrictModeCheck) priv_StrictModeCheck_checkObjectLiteralOrClass_5189(n pkg_rhino.Node) {
	_ = n
	var getters *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = getters
	var setters *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = setters
	var staticGetters *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = staticGetters
	var staticSetters *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = staticSetters
	{
		var key pkg_rhino.Node = n.GetLastChild()
		_ = key
		for ; key != 0; key = key.GetPrevious() {
			if (((key.IsEmpty() || key.IsComputedProp()) || key.IsSpread()) || key.IsComputedFieldDef()) || NodeUtil_IsClassStaticBlock(key) {
				continue
			}
			var keyName string = key.GetString()
			_ = keyName
			if !(key.IsSetterDef()) {
				var set *j2gort.Ext_java_util_LinkedHashSet
				_ = set
				if key.IsStaticMember() {
					set = staticGetters
				} else {
					set = getters
				}
				if !(j2gort.SetAddStr(set, keyName)) {
					this_.priv_StrictModeCheck_report_5190(key, StaticGet_StrictModeCheck_DUPLICATE_MEMBER(), []string{keyName})
				}
			}
			if !(key.IsGetterDef()) {
				var set_2 *j2gort.Ext_java_util_LinkedHashSet
				_ = set_2
				if key.IsStaticMember() {
					set_2 = staticSetters
				} else {
					set_2 = setters
				}
				if !(j2gort.SetAddStr(set_2, keyName)) {
					this_.priv_StrictModeCheck_report_5190(key, StaticGet_StrictModeCheck_DUPLICATE_MEMBER(), []string{keyName})
				}
			}
		}
	}
}

// Java static method: StrictModeCheck.isFunctionType(com.google.javascript.rhino.Node n)
func StrictModeCheck_IsFunctionType(n pkg_rhino.Node) bool {
	clinit_StrictModeCheck()
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	return ((type_ != nil && !type_.IsJavaNil_()) && type_.GetJSTypeBase().IsFunctionType())
}

// Java method: StrictModeCheck.report(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType diagnostic, java.lang.String[] args)
func (this_ *StrictModeCheck) priv_StrictModeCheck_report_5190(n pkg_rhino.Node, diagnostic *DiagnosticType, args []string) {
	_ = n
	_ = diagnostic
	_ = args
	this_.compiler.Report_1(JSError_S_Builder(diagnostic, args).SetLevel(this_.defaultLevel).SetNode(n).Build())
}

// Java class: com.google.javascript.jscomp.StrictModeCheck.NonExternChecks
type StrictModeCheck_NonExternChecks struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *StrictModeCheck
}

func (this_ *StrictModeCheck_NonExternChecks) BindSelf_(s any) {
	_ = s
}
func (this_ *StrictModeCheck_NonExternChecks) IsJavaNil_() bool { return this_ == nil }
func (*StrictModeCheck_NonExternChecks) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.StrictModeCheck$NonExternChecks", "NonExternChecks"
}
func (this_ *StrictModeCheck_NonExternChecks) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*StrictModeCheck_NonExternChecks) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.StrictModeCheck.NonExternChecks", "com.google.javascript.jscomp.StrictModeCheck$NonExternChecks", func(v_ any) bool {
		return j2gort.InstanceOf[*StrictModeCheck_NonExternChecks](v_)
	})
}
func (this_ *StrictModeCheck_NonExternChecks) SetOuter_(o *StrictModeCheck)                         { this_.outer_ = o }
func (this_ *StrictModeCheck_NonExternChecks) GetOuter_() *StrictModeCheck                          { return this_.outer_ }
func (*StrictModeCheck_NonExternChecks) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*StrictModeCheck_NonExternChecks) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *StrictModeCheck_NonExternChecks) initFields_StrictModeCheck_NonExternChecks() {
	_ = this_
}

// Java constructor: NonExternChecks()
func NewStrictModeCheck_NonExternChecks(outer_ *StrictModeCheck) *StrictModeCheck_NonExternChecks {
	this_ := &StrictModeCheck_NonExternChecks{}
	this_.outer_ = outer_
	InitStrictModeCheck_NonExternChecks(this_)
	return this_
}

func InitStrictModeCheck_NonExternChecks(this_ *StrictModeCheck_NonExternChecks) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_StrictModeCheck_NonExternChecks()
}

// Java method: NonExternChecks.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *StrictModeCheck_NonExternChecks) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() && StrictModeCheck_IsDeclaration(n) {
		this_.priv_NonExternChecks_checkDeclaration_5191(n)
	} else if n.IsGetProp() {
		this_.priv_NonExternChecks_checkGetProp_5192(n)
	}
}

// Java method: NonExternChecks.checkDeclaration(com.google.javascript.rhino.Node n)
func (this_ *StrictModeCheck_NonExternChecks) priv_NonExternChecks_checkDeclaration_5191(n pkg_rhino.Node) {
	_ = n
	{
		var iifeCond_2 bool
		{

			sr_5 := "eval"
			se_6 := n.GetString()
			iifeCond_2 = j2gort.CheckNotNullString(sr_5) == se_6
		}
		if iifeCond_2 {
			this_.outer_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_EVAL_DECLARATION(), []string{})
		} else {
			var iifeCond_1 bool
			{

				sr_7 := "arguments"
				se_8 := n.GetString()
				iifeCond_1 = j2gort.CheckNotNullString(sr_7) == se_8
			}
			if iifeCond_1 {
				this_.outer_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_ARGUMENTS_DECLARATION(), []string{})
			}
		}
	}
}

// Java method: NonExternChecks.checkGetProp(com.google.javascript.rhino.Node n)
func (this_ *StrictModeCheck_NonExternChecks) priv_NonExternChecks_checkGetProp_5192(n pkg_rhino.Node) {
	_ = n
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	var name string = n.GetString()
	_ = name
	if j2gort.CheckNotNullString(name) == "callee" {
		if target.IsName() && (j2gort.CheckNotNullString(target.GetString()) == "arguments") {
			this_.outer_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN(), []string{})
		}
	} else if j2gort.CheckNotNullString(name) == "caller" {
		if target.IsName() && (j2gort.CheckNotNullString(target.GetString()) == "arguments") {
			this_.outer_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN(), []string{})
		} else if StrictModeCheck_IsFunctionType(target) {
			this_.outer_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_FUNCTION_CALLER_FORBIDDEN(), []string{})
		}
	} else if (j2gort.CheckNotNullString(name) == "arguments") && StrictModeCheck_IsFunctionType(target) {
		this_.outer_.priv_StrictModeCheck_report_5190(n, StaticGet_StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN(), []string{})
	}
}
func (*StrictModeCheck) AlwaysShouldTraverse_()                      {}
func (*StrictModeCheck) CanReleaseSyntacticScopes_()                 {}
func (*StrictModeCheck_NonExternChecks) AlwaysShouldTraverse_()      {}
func (*StrictModeCheck_NonExternChecks) CanReleaseSyntacticScopes_() {}
