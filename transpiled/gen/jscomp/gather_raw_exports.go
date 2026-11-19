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
// Port of com/google/javascript/jscomp/GatherRawExports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.GatherRawExports
type GatherRawExports struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler          AbstractCompiler         `gson:"compiler"`
	exportedVariables j2gort.Ext_java_util_Set `gson:"exportedVariables"`
}

func (this_ *GatherRawExports) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherRawExports) IsJavaNil_() bool { return this_ == nil }
func (*GatherRawExports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherRawExports", "GatherRawExports"
}
func (this_ *GatherRawExports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*GatherRawExports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherRawExports", "com.google.javascript.jscomp.GatherRawExports", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherRawExports](v_)
	})
}
func (*GatherRawExports) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*GatherRawExports) Implements_jscomp_CompilerPass()                            {}
func (*GatherRawExports) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_GatherRawExports_GLOBAL_THIS_NAMES []string

func StaticGet_GatherRawExports_GLOBAL_THIS_NAMES() []string {
	clinit_GatherRawExports()
	return static_GatherRawExports_GLOBAL_THIS_NAMES
}
func StaticSet_GatherRawExports_GLOBAL_THIS_NAMES(v []string) {
	clinit_GatherRawExports()
	static_GatherRawExports_GLOBAL_THIS_NAMES = v
}

var clinitDone_GatherRawExports uint32
var clinitState_GatherRawExports j2gort.ClinitState

func clinit_GatherRawExports() {
	if atomic.LoadUint32(&clinitDone_GatherRawExports) != 0 {
		return
	}
	clinitSlow_GatherRawExports()
}

func clinitSlow_GatherRawExports() {
	if !j2gort.EnterClinit(&clinitState_GatherRawExports, &clinitDone_GatherRawExports) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GatherRawExports, &clinitDone_GatherRawExports)
	static_GatherRawExports_GLOBAL_THIS_NAMES = []string{"window", "globalThis", "top", "self", "goog$global", "goog.global", "$jscomp.global", "$jscomp$global"}
}

func (this_ *GatherRawExports) initFields_GatherRawExports() {
	_ = this_
	this_.exportedVariables = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: GatherRawExports(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewGatherRawExports(compiler AbstractCompiler) *GatherRawExports {
	clinit_GatherRawExports()
	this_ := &GatherRawExports{}
	InitGatherRawExports(this_, compiler)
	return this_
}

func InitGatherRawExports(this_ *GatherRawExports, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_GatherRawExports()
	this_.compiler = compiler
}

// Java method: GatherRawExports.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *GatherRawExports) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: GatherRawExports.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GatherRawExports) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsNormalOrOptChainGet(n) && GatherRawExports_IsGlobalThisObject(t, n.GetFirstChild()) {
		if NodeUtil_IsNormalOrOptChainGetProp(n) {
			j2gort.SetAddStr(this_.exportedVariables, n.GetString())
		} else if NodeUtil_IsNormalOrOptChainGet(n) && n.GetSecondChild().IsStringLit() {
			j2gort.SetAddStr(this_.exportedVariables, n.GetSecondChild().GetString())
		}
	}
}

// Java static method: GatherRawExports.isGlobalThisObject(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func GatherRawExports_IsGlobalThisObject(t *NodeTraversal, n pkg_rhino.Node) bool {
	clinit_GatherRawExports()
	_ = t
	_ = n
	if n.IsThis() {
		return t.InGlobalHoistScope()
	} else if n.IsQualifiedName() {
		var items int32 = int32(len(StaticGet_GatherRawExports_GLOBAL_THIS_NAMES()))
		_ = items
		{
			var i int32 = int32(0)
			_ = i
			for ; i < items; i = i + 1 {
				if n.MatchesQualifiedName_1_string(StaticGet_GatherRawExports_GLOBAL_THIS_NAMES()[i]) {
					return true
				}
			}
		}
	}
	return false
}

// Java method: GatherRawExports.getExportedVariableNames()
func (this_ *GatherRawExports) GetExportedVariableNames() j2gort.Ext_java_util_Set {
	return this_.exportedVariables
}
func (*GatherRawExports) AlwaysShouldTraverse_()      {}
func (*GatherRawExports) CanReleaseSyntacticScopes_() {}
