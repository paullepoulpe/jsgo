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
// Port of com/google/javascript/jscomp/CollectFileOverviewVisibility.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CollectFileOverviewVisibility
type CollectFileOverviewVisibility struct {
	builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder `gson:"builder"`
	idHash_ atomic.Int32
}

func (this_ *CollectFileOverviewVisibility) BindSelf_(s any) {
	_ = s
}
func (this_ *CollectFileOverviewVisibility) IsJavaNil_() bool { return this_ == nil }
func (*CollectFileOverviewVisibility) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CollectFileOverviewVisibility", "CollectFileOverviewVisibility"
}
func (this_ *CollectFileOverviewVisibility) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CollectFileOverviewVisibility) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CollectFileOverviewVisibility", "com.google.javascript.jscomp.CollectFileOverviewVisibility", func(v_ any) bool {
		return j2gort.InstanceOf[*CollectFileOverviewVisibility](v_)
	})
}
func (*CollectFileOverviewVisibility) Implements_jscomp_CompilerPass() {}

func (this_ *CollectFileOverviewVisibility) initFields_CollectFileOverviewVisibility() {
	_ = this_
	this_.builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
}

// Java constructor: CollectFileOverviewVisibility(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCollectFileOverviewVisibility(compiler AbstractCompiler) *CollectFileOverviewVisibility {
	this_ := &CollectFileOverviewVisibility{}
	InitCollectFileOverviewVisibility(this_, compiler)
	return this_
}

func InitCollectFileOverviewVisibility(this_ *CollectFileOverviewVisibility, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CollectFileOverviewVisibility()
}

// Java method: CollectFileOverviewVisibility.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CollectFileOverviewVisibility) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; script != 0; script = script.GetNext() {
			if !(script.IsScript()) {
				j2gort.CheckState(false)
			}
			this_.priv_CollectFileOverviewVisibility_visit_1397(script)
		}
	}
}

// Java method: CollectFileOverviewVisibility.visit(com.google.javascript.rhino.Node scriptNode)
func (this_ *CollectFileOverviewVisibility) priv_CollectFileOverviewVisibility_visit_1397(scriptNode pkg_rhino.Node) {
	_ = scriptNode
	var jsDocInfo *pkg_rhino.JSDocInfo = scriptNode.GetJSDocInfo()
	_ = jsDocInfo
	if jsDocInfo == nil {
		return
	}
	var v pkg_rhino.JSDocInfo_Visibility = jsDocInfo.GetVisibility()
	_ = v
	if v == 0 {
		return
	}
	this_.builder.Put_2(scriptNode.GetStaticSourceFile(), j2gort.BoxNullableVal(v))
}

// Java method: CollectFileOverviewVisibility.getFileOverviewVisibilityMap()
func (this_ *CollectFileOverviewVisibility) GetFileOverviewVisibilityMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](this_.builder.BuildOrThrow())
}
