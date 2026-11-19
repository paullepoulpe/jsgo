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
// Port of com/google/javascript/jscomp/lint/CheckMissingSemicolon.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckMissingSemicolon
type CheckMissingSemicolon struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckMissingSemicolon) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckMissingSemicolon) IsJavaNil_() bool { return this_ == nil }
func (*CheckMissingSemicolon) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckMissingSemicolon", "CheckMissingSemicolon"
}
func (this_ *CheckMissingSemicolon) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckMissingSemicolon) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckMissingSemicolon", "com.google.javascript.jscomp.lint.CheckMissingSemicolon", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckMissingSemicolon](v_)
	})
}
func (*CheckMissingSemicolon) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckMissingSemicolon) Implements_jscomp_CompilerPass()                            {}
func (*CheckMissingSemicolon) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckMissingSemicolon_MISSING_SEMICOLON *DiagnosticType

func StaticGet_CheckMissingSemicolon_MISSING_SEMICOLON() *DiagnosticType {
	clinit_CheckMissingSemicolon()
	return static_CheckMissingSemicolon_MISSING_SEMICOLON
}
func StaticSet_CheckMissingSemicolon_MISSING_SEMICOLON(v *DiagnosticType) {
	clinit_CheckMissingSemicolon()
	static_CheckMissingSemicolon_MISSING_SEMICOLON = v
}

var clinitDone_CheckMissingSemicolon uint32
var clinitState_CheckMissingSemicolon j2gort.ClinitState

func clinit_CheckMissingSemicolon() {
	if atomic.LoadUint32(&clinitDone_CheckMissingSemicolon) != 0 {
		return
	}
	clinitSlow_CheckMissingSemicolon()
}

func clinitSlow_CheckMissingSemicolon() {
	if !j2gort.EnterClinit(&clinitState_CheckMissingSemicolon, &clinitDone_CheckMissingSemicolon) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckMissingSemicolon, &clinitDone_CheckMissingSemicolon)
	static_CheckMissingSemicolon_MISSING_SEMICOLON = DiagnosticType_Disabled("JSC_MISSING_SEMICOLON", "Missing semicolon")
}

func (this_ *CheckMissingSemicolon) initFields_CheckMissingSemicolon() {
	_ = this_
}

// Java constructor: CheckMissingSemicolon(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckMissingSemicolon(compiler AbstractCompiler) *CheckMissingSemicolon {
	clinit_CheckMissingSemicolon()
	this_ := &CheckMissingSemicolon{}
	InitCheckMissingSemicolon(this_, compiler)
	return this_
}

func InitCheckMissingSemicolon(this_ *CheckMissingSemicolon, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckMissingSemicolon()
	this_.compiler = compiler
}

// Java method: CheckMissingSemicolon.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckMissingSemicolon) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckMissingSemicolon.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckMissingSemicolon) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (!(n.IsScript()) && NodeUtil_IsStatement(n)) && this_.priv_CheckMissingSemicolon_shouldHaveSemicolon_6545(n) {
		this_.priv_CheckMissingSemicolon_checkSemicolon_6546(t, n)
	}
}

// Java method: CheckMissingSemicolon.shouldHaveSemicolon(com.google.javascript.rhino.Node statement)
func (this_ *CheckMissingSemicolon) priv_CheckMissingSemicolon_shouldHaveSemicolon_6545(statement pkg_rhino.Node) bool {
	_ = statement
	if ((((statement.IsFunction() || statement.IsClass()) || statement.IsBlock()) || statement.IsLabelName()) || statement.IsModuleBody()) || (NodeUtil_IsControlStructure(statement) && !(statement.IsDo())) {
		return false
	}
	if statement.IsExport() {
		return this_.priv_CheckMissingSemicolon_shouldHaveSemicolon_6545(statement.GetFirstChild())
	}
	return true
}

// Java method: CheckMissingSemicolon.checkSemicolon(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckMissingSemicolon) priv_CheckMissingSemicolon_checkSemicolon_6546(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var staticSourceFile pkg_rhino.StaticSourceFile = n.GetStaticSourceFile()
	_ = staticSourceFile
	var sourceFile *SourceFile
	_ = sourceFile
	var ok_1 bool
	sourceFile, ok_1 = j2gort.TryCast[*SourceFile](staticSourceFile)
	if ok_1 {
		var code string = j2gort.NullString
		_ = code
		var didRet_2 bool
		func() {
			defer func() {
				if rec_3 := recover(); rec_3 != nil {
					jex_4 := javalang.RecoverException(rec_3)
					_ = jex_4
					if j2gort.InstanceOf[javaio.IIOException](jex_4) {
						e := j2gort.AsIface[javaio.IIOException](jex_4)
						_ = e
						didRet_2 = true
						return
					} else {
						panic(rec_3)
					}
				}
			}()
			code = sourceFile.GetCode()
		}()
		if didRet_2 {
			return
		}
		var length int32 = n.GetLength()
		_ = length
		if length == int32(0) {
			return
		}
		var position int32 = ((n.GetSourceOffset() + length) - int32(1))
		_ = position
		var endsWithSemicolon bool = (int32(j2gort.JStrCharAt(code, position)) == int32(uint16(59)))
		_ = endsWithSemicolon
		if !(endsWithSemicolon) {
			t.Report_3(n, StaticGet_CheckMissingSemicolon_MISSING_SEMICOLON(), []string{})
		}
	}
}
func (*CheckMissingSemicolon) AlwaysShouldTraverse_()      {}
func (*CheckMissingSemicolon) CanReleaseSyntacticScopes_() {}
