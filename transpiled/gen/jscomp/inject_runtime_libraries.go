/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InjectRuntimeLibraries.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.InjectRuntimeLibraries
type InjectRuntimeLibraries struct {
	runtimeLibs *RuntimeJsLibManager `gson:"runtimeLibs"`
	forceInjectedLibraries j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"forceInjectedLibraries"`
	idHash_ atomic.Int32
}

func (this_ *InjectRuntimeLibraries) BindSelf_(s any) {
	_ = s
}
func (this_ *InjectRuntimeLibraries) IsJavaNil_() bool { return this_ == nil }
func (*InjectRuntimeLibraries) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.InjectRuntimeLibraries", "InjectRuntimeLibraries" }
func (this_ *InjectRuntimeLibraries) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InjectRuntimeLibraries) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InjectRuntimeLibraries", "com.google.javascript.jscomp.InjectRuntimeLibraries", func(v_ any) bool {
		return j2gort.InstanceOf[*InjectRuntimeLibraries](v_)
	})
}
func (*InjectRuntimeLibraries) Implements_jscomp_CompilerPass() {}

func (this_ *InjectRuntimeLibraries) initFields_InjectRuntimeLibraries() {
	_ = this_
}

// Java constructor: InjectRuntimeLibraries(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableSet<java.lang.String> forceInjectedLibraries)
func NewInjectRuntimeLibraries(compiler AbstractCompiler, forceInjectedLibraries j2gort.Ext_com_google_common_collect_ImmutableSet) *InjectRuntimeLibraries {
	this_ := &InjectRuntimeLibraries{}
	InitInjectRuntimeLibraries(this_, compiler, forceInjectedLibraries)
	return this_
}

func InitInjectRuntimeLibraries(this_ *InjectRuntimeLibraries, compiler AbstractCompiler, forceInjectedLibraries j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = compiler
	_ = forceInjectedLibraries
	this_.initFields_InjectRuntimeLibraries()
	this_.runtimeLibs = compiler.GetRuntimeJsLibManager()
	this_.forceInjectedLibraries = forceInjectedLibraries
}

// Java method: InjectRuntimeLibraries.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InjectRuntimeLibraries) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.priv_InjectRuntimeLibraries_injectLibraries_3206()
}

// Java method: InjectRuntimeLibraries.injectLibraries()
func (this_ *InjectRuntimeLibraries) priv_InjectRuntimeLibraries_injectLibraries_3206() {
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.forceInjectedLibraries))
		for cur_1.HasNext() {
			forced := cur_1.NextStr()
			_ = forced
			this_.runtimeLibs.EnsureLibraryInjected(forced, true)
		}
	}
}

