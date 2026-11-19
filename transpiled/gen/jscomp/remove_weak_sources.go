/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RemoveWeakSources.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RemoveWeakSources
type RemoveWeakSources struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_ atomic.Int32
}

func (this_ *RemoveWeakSources) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveWeakSources) IsJavaNil_() bool { return this_ == nil }
func (*RemoveWeakSources) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.RemoveWeakSources", "RemoveWeakSources" }
func (this_ *RemoveWeakSources) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveWeakSources) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveWeakSources", "com.google.javascript.jscomp.RemoveWeakSources", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveWeakSources](v_)
	})
}
func (*RemoveWeakSources) Implements_jscomp_CompilerPass() {}

func (this_ *RemoveWeakSources) initFields_RemoveWeakSources() {
	_ = this_
}

// Java constructor: RemoveWeakSources(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRemoveWeakSources(compiler AbstractCompiler) *RemoveWeakSources {
	this_ := &RemoveWeakSources{}
	InitRemoveWeakSources(this_, compiler)
	return this_
}

func InitRemoveWeakSources(this_ *RemoveWeakSources, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RemoveWeakSources()
	this_.compiler = compiler
}

// Java method: RemoveWeakSources.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RemoveWeakSources) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; (script != 0); script = script.GetNext() {
			if !(script.IsScript()) { j2gort.CheckState(false) }
			if (script.GetStaticSourceFile().IsWeak()) {
				NodeUtil_DeleteChildren(script, this_.compiler)
			}
		}
	}
}

