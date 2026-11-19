/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6CheckModule.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6CheckModule
type Es6CheckModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *Es6CheckModule) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6CheckModule) IsJavaNil_() bool { return this_ == nil }
func (*Es6CheckModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6CheckModule", "Es6CheckModule"
}
func (this_ *Es6CheckModule) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6CheckModule) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6CheckModule", "com.google.javascript.jscomp.Es6CheckModule", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6CheckModule](v_)
	})
}
func (*Es6CheckModule) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6CheckModule) Implements_jscomp_CompilerPass()                            {}
func (*Es6CheckModule) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_Es6CheckModule_ES6_MODULE_REFERENCES_THIS *DiagnosticType

func StaticGet_Es6CheckModule_ES6_MODULE_REFERENCES_THIS() *DiagnosticType {
	clinit_Es6CheckModule()
	return static_Es6CheckModule_ES6_MODULE_REFERENCES_THIS
}
func StaticSet_Es6CheckModule_ES6_MODULE_REFERENCES_THIS(v *DiagnosticType) {
	clinit_Es6CheckModule()
	static_Es6CheckModule_ES6_MODULE_REFERENCES_THIS = v
}

var static_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED *DiagnosticType

func StaticGet_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED() *DiagnosticType {
	clinit_Es6CheckModule()
	return static_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED
}
func StaticSet_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED(v *DiagnosticType) {
	clinit_Es6CheckModule()
	static_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED = v
}

var clinitDone_Es6CheckModule uint32
var clinitState_Es6CheckModule j2gort.ClinitState

func clinit_Es6CheckModule() {
	if atomic.LoadUint32(&clinitDone_Es6CheckModule) != 0 {
		return
	}
	clinitSlow_Es6CheckModule()
}

func clinitSlow_Es6CheckModule() {
	if !j2gort.EnterClinit(&clinitState_Es6CheckModule, &clinitDone_Es6CheckModule) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6CheckModule, &clinitDone_Es6CheckModule)
	static_Es6CheckModule_ES6_MODULE_REFERENCES_THIS = DiagnosticType_Warning("ES6_MODULE_REFERENCES_THIS", "The body of an ES6 module cannot reference 'this'.")
	static_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED = DiagnosticType_Error("JSC_IMPORT_CANNOT_BE_REASSIGNED", "Assignment to constant variable \"{0}\".")
}

func (this_ *Es6CheckModule) initFields_Es6CheckModule() {
	_ = this_
}

// Java constructor: Es6CheckModule(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6CheckModule(compiler AbstractCompiler) *Es6CheckModule {
	clinit_Es6CheckModule()
	this_ := &Es6CheckModule{}
	InitEs6CheckModule(this_, compiler)
	return this_
}

func InitEs6CheckModule(this_ *Es6CheckModule, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6CheckModule()
	this_.compiler = compiler
}

// Java method: Es6CheckModule.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6CheckModule) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: Es6CheckModule.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6CheckModule) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(37) /* Token.THIS */ :
		if t.InModuleHoistScope() {
			t.Report_3(n, StaticGet_Es6CheckModule_ES6_MODULE_REFERENCES_THIS(), []string{})
		}
	case int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		if (NodeUtil_IsLValue(n) && !(NodeUtil_IsDeclarationLValue(n))) && n.GetFirstChild().IsName() {
			var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetFirstChild().GetString()))
			_ = var_
			if var_ != nil {
				var nameNode pkg_rhino.Node = var_.GetNameNode()
				_ = nameNode
				if (nameNode != 0) && nameNode.IsImportStar() {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED(), []string{nameNode.GetString()}))
				}
			}
		}
	case int32(32) /* Token.NAME */ :
		if NodeUtil_IsLValue(n) && !(NodeUtil_IsDeclarationLValue(n)) {
			var var__2 *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
			_ = var__2
			if var__2 != nil {
				var nameNode_2 pkg_rhino.Node = var__2.GetNameNode()
				_ = nameNode_2
				if (nameNode_2 != 0) && (nameNode_2 != n) {
					if NodeUtil_IsImportedName(nameNode_2) {
						this_.compiler.Report_1(JSError_Make_3(n, StaticGet_Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED(), []string{nameNode_2.GetString()}))
					}
				}
			}
		}
	default:
	}
}
func (*Es6CheckModule) AlwaysShouldTraverse_()      {}
func (*Es6CheckModule) CanReleaseSyntacticScopes_() {}
