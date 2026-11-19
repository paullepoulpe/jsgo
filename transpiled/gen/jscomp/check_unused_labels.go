/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckUnusedLabels.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckUnusedLabels
type CheckUnusedLabels struct {
	compiler       AbstractCompiler                `gson:"compiler"`
	currentContext *CheckUnusedLabels_LabelContext `gson:"currentContext"`
	idHash_        atomic.Int32
}

func (this_ *CheckUnusedLabels) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckUnusedLabels) IsJavaNil_() bool { return this_ == nil }
func (*CheckUnusedLabels) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckUnusedLabels", "CheckUnusedLabels"
}
func (this_ *CheckUnusedLabels) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckUnusedLabels) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckUnusedLabels", "com.google.javascript.jscomp.lint.CheckUnusedLabels", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckUnusedLabels](v_)
	})
}
func (*CheckUnusedLabels) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckUnusedLabels) Implements_jscomp_CompilerPass()           {}

var static_CheckUnusedLabels_UNUSED_LABEL *DiagnosticType

func StaticGet_CheckUnusedLabels_UNUSED_LABEL() *DiagnosticType {
	clinit_CheckUnusedLabels()
	return static_CheckUnusedLabels_UNUSED_LABEL
}
func StaticSet_CheckUnusedLabels_UNUSED_LABEL(v *DiagnosticType) {
	clinit_CheckUnusedLabels()
	static_CheckUnusedLabels_UNUSED_LABEL = v
}

var clinitDone_CheckUnusedLabels uint32
var clinitState_CheckUnusedLabels j2gort.ClinitState

func clinit_CheckUnusedLabels() {
	if atomic.LoadUint32(&clinitDone_CheckUnusedLabels) != 0 {
		return
	}
	clinitSlow_CheckUnusedLabels()
}

func clinitSlow_CheckUnusedLabels() {
	if !j2gort.EnterClinit(&clinitState_CheckUnusedLabels, &clinitDone_CheckUnusedLabels) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckUnusedLabels, &clinitDone_CheckUnusedLabels)
	static_CheckUnusedLabels_UNUSED_LABEL = DiagnosticType_Disabled("JSC_UNUSED_LABEL", "Unused label {0}.")
}

func (this_ *CheckUnusedLabels) initFields_CheckUnusedLabels() {
	_ = this_
}

// Java constructor: CheckUnusedLabels(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckUnusedLabels(compiler AbstractCompiler) *CheckUnusedLabels {
	clinit_CheckUnusedLabels()
	this_ := &CheckUnusedLabels{}
	InitCheckUnusedLabels(this_, compiler)
	return this_
}

func InitCheckUnusedLabels(this_ *CheckUnusedLabels, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckUnusedLabels()
	this_.compiler = compiler
}

// Java method: CheckUnusedLabels.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckUnusedLabels) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckUnusedLabels.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUnusedLabels) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(82) /* Token.BREAK */, int32(83) /* Token.CONTINUE */ :
		if n.HasChildren() {
			var temp *CheckUnusedLabels_LabelContext = this_.currentContext
			_ = temp
			for temp != nil {
				{
					var iifeCond_1 bool
					{

						sr_1 := temp.name
						se_2 := n.GetFirstChild().GetString()
						iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
					}
					if iifeCond_1 {
						{
							r_3 := temp
							r_3.used = true
						}
						break
					}
				}
				temp = temp.parent
			}
		}
		return false
	case int32(92) /* Token.LABEL */ :
		this_.currentContext = NewCheckUnusedLabels_LabelContext(n.GetFirstChild().GetString(), this_.currentContext)
	default:
	}
	return true
}

// Java method: CheckUnusedLabels.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUnusedLabels) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsLabel() && (this_.currentContext != nil) {
		if !(this_.currentContext.used) {
			t.Report_3(n, StaticGet_CheckUnusedLabels_UNUSED_LABEL(), []string{n.GetFirstChild().GetString()})
		}
		this_.currentContext = this_.currentContext.parent
	}
}

// Java class: com.google.javascript.jscomp.lint.CheckUnusedLabels.LabelContext
type CheckUnusedLabels_LabelContext struct {
	name    string                          `gson:"name"`
	parent  *CheckUnusedLabels_LabelContext `gson:"parent"`
	idHash_ atomic.Int32
	used    bool `gson:"used"`
}

func (this_ *CheckUnusedLabels_LabelContext) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckUnusedLabels_LabelContext) IsJavaNil_() bool { return this_ == nil }
func (*CheckUnusedLabels_LabelContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckUnusedLabels$LabelContext", "LabelContext"
}
func (this_ *CheckUnusedLabels_LabelContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckUnusedLabels_LabelContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckUnusedLabels.LabelContext", "com.google.javascript.jscomp.lint.CheckUnusedLabels$LabelContext", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckUnusedLabels_LabelContext](v_)
	})
}

func (this_ *CheckUnusedLabels_LabelContext) initFields_CheckUnusedLabels_LabelContext() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: LabelContext(java.lang.String name, com.google.javascript.jscomp.lint.CheckUnusedLabels.LabelContext parent)
func NewCheckUnusedLabels_LabelContext(name string, parent *CheckUnusedLabels_LabelContext) *CheckUnusedLabels_LabelContext {
	this_ := &CheckUnusedLabels_LabelContext{}
	InitCheckUnusedLabels_LabelContext(this_, name, parent)
	return this_
}

func InitCheckUnusedLabels_LabelContext(this_ *CheckUnusedLabels_LabelContext, name string, parent *CheckUnusedLabels_LabelContext) {
	_ = name
	_ = parent
	this_.initFields_CheckUnusedLabels_LabelContext()
	this_.name = name
	this_.parent = parent
}
func (*CheckUnusedLabels) CanReleaseSyntacticScopes_() {}
