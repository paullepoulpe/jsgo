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
// Port of com/google/javascript/jscomp/PassConfig.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.PassConfig
type PassConfig interface {
	Implements_jscomp_PassConfig()
	IsJavaNil_() bool
	GetPassConfigBase() *PassConfigBase
	GetWhitespaceOnlyPasses() *PassListBuilder
	GetTranspileOnlyPasses() *PassListBuilder
	GetChecks() *PassListBuilder
	GetOptimizations() *PassListBuilder
	GetFinalizations() *PassListBuilder
	GetPassGraph() pkg_graph.GraphvizGraph
	GetBasePassConfig() PassConfig
}

func ToIface_PassConfig[T interface {
	PassConfig
	comparable
}](v T) PassConfig {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.PassConfig
type PassConfigBase struct {
	self_   PassConfig
	options *CompilerOptions `gson:"options"`
	idHash_ atomic.Int32
}

func (this_ *PassConfigBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(PassConfig)
}
func (this_ *PassConfigBase) IsJavaNil_() bool { return this_ == nil }
func (*PassConfigBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PassConfig", "PassConfig"
}
func (this_ *PassConfigBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassConfig", "com.google.javascript.jscomp.PassConfig", func(v_ any) bool {
		return j2gort.InstanceOf[PassConfig](v_)
	})
}
func (b *PassConfigBase) GetPassConfigBase() *PassConfigBase { return b }
func (*PassConfigBase) Implements_jscomp_PassConfig()        {}

func (this_ *PassConfigBase) initFields_PassConfig() {
	_ = this_
}

func InitPassConfig(this_ *PassConfigBase, options *CompilerOptions) {
	_ = options
	this_.initFields_PassConfig()
	this_.options = options
}

// Java method: PassConfig.getWhitespaceOnlyPasses()
func (this_ *PassConfigBase) GetWhitespaceOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(this_.options)
}

// Java method: PassConfig.getTranspileOnlyPasses()
func (this_ *PassConfigBase) GetTranspileOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(this_.options)
}

// Java method: PassConfig.getChecks()
func (this_ *PassConfigBase) GetChecks() *PassListBuilder {
	panic("abstract or native method invoked")
}

// Java method: PassConfig.getOptimizations()
func (this_ *PassConfigBase) GetOptimizations() *PassListBuilder {
	panic("abstract or native method invoked")
}

// Java method: PassConfig.getFinalizations()
func (this_ *PassConfigBase) GetFinalizations() *PassListBuilder {
	panic("abstract or native method invoked")
}

// Java method: PassConfig.getPassGraph()
func (this_ *PassConfigBase) GetPassGraph() pkg_graph.GraphvizGraph {
	var graph pkg_graph.LinkedDirectedGraph[string, string] = pkg_graph.LinkedDirectedGraph_CreateWithoutAnnotations[string, string]()
	_ = graph
	var allPasses j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Concat_2(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_.GetChecks().Build()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_.GetOptimizations().Build()))
	_ = allPasses
	var lastPass string = j2gort.NullString
	_ = lastPass
	var loopStart string = j2gort.NullString
	_ = loopStart
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(allPasses))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			pass := j2gort.AsIface[PassFactory](rawElem_2)
			_ = pass
			var passName string = pass.GetName()
			_ = passName
			var i int32 = int32(1)
			_ = i
			for graph.HasNode(passName) {
				passName = (j2gort.StrValString(pass.GetName()) + j2gort.StrValInt32(func() int32 { old_3 := i; nw_4 := old_3 + 1; i = nw_4; return old_3 }()))
			}
			graph.CreateNode(passName)
			if (loopStart == j2gort.NullString) && pass.IsRunInFixedPointLoop() {
				loopStart = passName
			} else if (loopStart != j2gort.NullString) && !(pass.IsRunInFixedPointLoop()) {
				graph.Connect_3_any_any_any(lastPass, "loop", loopStart)
				loopStart = j2gort.NullString
			}
			if lastPass != j2gort.NullString {
				graph.Connect_3_any_any_any(lastPass, "", passName)
			}
			lastPass = passName
		}
	}
	return j2gort.AsIface[pkg_graph.GraphvizGraph](graph)
}

// Java method: PassConfig.getBasePassConfig()
func (this_ *PassConfigBase) GetBasePassConfig() PassConfig {
	var current PassConfig = this_.self_
	_ = current
	for j2gort.InstanceOf[PassConfig_PassConfigDelegate](current) {
		current = (j2gort.AsIface[PassConfig_PassConfigDelegate](current)).GetPassConfig_PassConfigDelegateBase().delegate
	}
	return current
}

// Java polymorphic class interface: com.google.javascript.jscomp.PassConfig.PassConfigDelegate
type PassConfig_PassConfigDelegate interface {
	Implements_jscomp_PassConfig_PassConfigDelegate()
	IsJavaNil_() bool
	GetPassConfig_PassConfigDelegateBase() *PassConfig_PassConfigDelegateBase
	PassConfig
}

func ToIface_PassConfig_PassConfigDelegate[T interface {
	PassConfig_PassConfigDelegate
	comparable
}](v T) PassConfig_PassConfigDelegate {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.PassConfig.PassConfigDelegate
type PassConfig_PassConfigDelegateBase struct {
	PassConfigBase
	delegate PassConfig `gson:"delegate"`
}

func (this_ *PassConfig_PassConfigDelegateBase) BindSelf_(s any) {
	_ = s
	this_.PassConfigBase.BindSelf_(s)
}
func (this_ *PassConfig_PassConfigDelegateBase) IsJavaNil_() bool { return this_ == nil }
func (*PassConfig_PassConfigDelegateBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PassConfig$PassConfigDelegate", "PassConfigDelegate"
}
func (this_ *PassConfig_PassConfigDelegateBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassConfigBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassConfig.PassConfigDelegate", "com.google.javascript.jscomp.PassConfig$PassConfigDelegate", func(v_ any) bool {
		return j2gort.InstanceOf[PassConfig_PassConfigDelegate](v_)
	})
}
func (b *PassConfig_PassConfigDelegateBase) GetPassConfig_PassConfigDelegateBase() *PassConfig_PassConfigDelegateBase {
	return b
}
func (*PassConfig_PassConfigDelegateBase) Implements_jscomp_PassConfig_PassConfigDelegate() {}
func (*PassConfig_PassConfigDelegateBase) Implements_jscomp_PassConfig()                    {}

func (this_ *PassConfig_PassConfigDelegateBase) initFields_PassConfig_PassConfigDelegate() {
	_ = this_
}

// Java constructor: PassConfigDelegate(com.google.javascript.jscomp.PassConfig delegate)
func NewPassConfig_PassConfigDelegate(delegate PassConfig) PassConfig_PassConfigDelegate {
	this_ := &PassConfig_PassConfigDelegateBase{}
	this_.BindSelf_(this_)
	InitPassConfig_PassConfigDelegate(this_, delegate)
	return this_
}

func InitPassConfig_PassConfigDelegate(this_ *PassConfig_PassConfigDelegateBase, delegate PassConfig) {
	_ = delegate
	InitPassConfig(&this_.PassConfigBase, delegate.GetPassConfigBase().options)
	this_.initFields_PassConfig_PassConfigDelegate()
	this_.delegate = delegate
}

// Java method: PassConfigDelegate.getWhitespaceOnlyPasses()
func (this_ *PassConfig_PassConfigDelegateBase) GetWhitespaceOnlyPasses() *PassListBuilder {
	return this_.delegate.GetWhitespaceOnlyPasses()
}

// Java method: PassConfigDelegate.getChecks()
func (this_ *PassConfig_PassConfigDelegateBase) GetChecks() *PassListBuilder {
	return this_.delegate.GetChecks()
}

// Java method: PassConfigDelegate.getOptimizations()
func (this_ *PassConfig_PassConfigDelegateBase) GetOptimizations() *PassListBuilder {
	return this_.delegate.GetOptimizations()
}

// Java method: PassConfigDelegate.getFinalizations()
func (this_ *PassConfig_PassConfigDelegateBase) GetFinalizations() *PassListBuilder {
	return this_.delegate.GetFinalizations()
}

// Java method: PassConfigDelegate.getTranspileOnlyPasses()
func (this_ *PassConfig_PassConfigDelegateBase) GetTranspileOnlyPasses() *PassListBuilder {
	return this_.delegate.GetTranspileOnlyPasses()
}
