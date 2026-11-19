/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ForbidDynamicImportUsage.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ForbidDynamicImportUsage
type ForbidDynamicImportUsage struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *ForbidDynamicImportUsage) BindSelf_(s any) {
	_ = s
}
func (this_ *ForbidDynamicImportUsage) IsJavaNil_() bool { return this_ == nil }
func (*ForbidDynamicImportUsage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ForbidDynamicImportUsage", "ForbidDynamicImportUsage"
}
func (this_ *ForbidDynamicImportUsage) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ForbidDynamicImportUsage) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ForbidDynamicImportUsage", "com.google.javascript.jscomp.ForbidDynamicImportUsage", func(v_ any) bool {
		return j2gort.InstanceOf[*ForbidDynamicImportUsage](v_)
	})
}
func (*ForbidDynamicImportUsage) Implements_jscomp_CompilerPass()           {}
func (*ForbidDynamicImportUsage) Implements_jscomp_NodeTraversal_Callback() {}

var static_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE *DiagnosticType

func StaticGet_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE() *DiagnosticType {
	clinit_ForbidDynamicImportUsage()
	return static_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE
}
func StaticSet_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE(v *DiagnosticType) {
	clinit_ForbidDynamicImportUsage()
	static_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE = v
}

var clinitDone_ForbidDynamicImportUsage uint32
var clinitState_ForbidDynamicImportUsage j2gort.ClinitState

func clinit_ForbidDynamicImportUsage() {
	if atomic.LoadUint32(&clinitDone_ForbidDynamicImportUsage) != 0 {
		return
	}
	clinitSlow_ForbidDynamicImportUsage()
}

func clinitSlow_ForbidDynamicImportUsage() {
	if !j2gort.EnterClinit(&clinitState_ForbidDynamicImportUsage, &clinitDone_ForbidDynamicImportUsage) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ForbidDynamicImportUsage, &clinitDone_ForbidDynamicImportUsage)
	static_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE = DiagnosticType_Error("JSC_DYNAMIC_IMPORT_USAGE", "Dynamic import expressions cannot be transpiled.")
}

func (this_ *ForbidDynamicImportUsage) initFields_ForbidDynamicImportUsage() {
	_ = this_
}

// Java constructor: ForbidDynamicImportUsage(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewForbidDynamicImportUsage(compiler AbstractCompiler) *ForbidDynamicImportUsage {
	clinit_ForbidDynamicImportUsage()
	this_ := &ForbidDynamicImportUsage{}
	InitForbidDynamicImportUsage(this_, compiler)
	return this_
}

func InitForbidDynamicImportUsage(this_ *ForbidDynamicImportUsage, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ForbidDynamicImportUsage()
	this_.compiler = compiler
}

// Java method: ForbidDynamicImportUsage.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ForbidDynamicImportUsage) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	this_.compiler.MarkFeatureNotAllowed(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
}

// Java method: ForbidDynamicImportUsage.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ForbidDynamicImportUsage) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		return this_.priv_ForbidDynamicImportUsage_getFeatureSetOfScript_2875(n).Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
	} else {
		return true
	}
}

// Java method: ForbidDynamicImportUsage.getFeatureSetOfScript(com.google.javascript.rhino.Node n)
func (this_ *ForbidDynamicImportUsage) priv_ForbidDynamicImportUsage_getFeatureSetOfScript_2875(n pkg_rhino.Node) *pkg_parser.FeatureSet {
	_ = n
	var featureSet *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
	_ = featureSet
	if featureSet == nil {
		return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
	} else {
		return featureSet
	}
}

// Java method: ForbidDynamicImportUsage.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ForbidDynamicImportUsage) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		t.Report_3(n, StaticGet_ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE(), []string{})
	default:
	}
}
func (*ForbidDynamicImportUsage) CanReleaseSyntacticScopes_() {}
