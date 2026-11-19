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
// Port of com/google/javascript/jscomp/J2clSourceFileChecker.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.J2clSourceFileChecker
type J2clSourceFileChecker struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_ atomic.Int32
}

func (this_ *J2clSourceFileChecker) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clSourceFileChecker) IsJavaNil_() bool { return this_ == nil }
func (*J2clSourceFileChecker) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.J2clSourceFileChecker", "J2clSourceFileChecker" }
func (this_ *J2clSourceFileChecker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*J2clSourceFileChecker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clSourceFileChecker", "com.google.javascript.jscomp.J2clSourceFileChecker", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clSourceFileChecker](v_)
	})
}
func (*J2clSourceFileChecker) Implements_jscomp_CompilerPass() {}

func (this_ *J2clSourceFileChecker) initFields_J2clSourceFileChecker() {
	_ = this_
}

// Java constructor: J2clSourceFileChecker(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clSourceFileChecker(compiler AbstractCompiler) *J2clSourceFileChecker {
	this_ := &J2clSourceFileChecker{}
	InitJ2clSourceFileChecker(this_, compiler)
	return this_
}

func InitJ2clSourceFileChecker(this_ *J2clSourceFileChecker, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_J2clSourceFileChecker()
	this_.compiler = compiler
}

// Java static method: J2clSourceFileChecker.hasJ2cl(com.google.javascript.rhino.Node root)
func J2clSourceFileChecker_HasJ2cl(root pkg_rhino.Node) bool {
	_ = root
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; (script != 0); script = script.GetNext() {
			if !(script.IsScript()) { j2gort.CheckState(false) }
			if (J2clSourceUtils_IsJ2clSource_1_Node(script)) {
				return true
			}
		}
	}
	return false
}

// Java method: J2clSourceFileChecker.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clSourceFileChecker) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if (J2clSourceFileChecker_HasJ2cl(root)) {
		J2clSourceFileChecker_MarkToRunJ2clPasses(this_.compiler)
	}
}

// Java static method: J2clSourceFileChecker.markToRunJ2clPasses(com.google.javascript.jscomp.AbstractCompiler compiler)
func J2clSourceFileChecker_MarkToRunJ2clPasses(compiler AbstractCompiler) {
	_ = compiler
	compiler.SetRunJ2clPasses(true)
}

// Java static method: J2clSourceFileChecker.shouldRunJ2clPasses(com.google.javascript.jscomp.AbstractCompiler compiler)
func J2clSourceFileChecker_ShouldRunJ2clPasses(compiler AbstractCompiler) bool {
	_ = compiler
	return compiler.RunJ2clPasses()
}

