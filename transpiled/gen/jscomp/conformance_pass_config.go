/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ConformancePassConfig.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.ConformancePassConfig
type ConformancePassConfig struct {
	PassConfig_PassConfigDelegateBase
	f_ConformancePassConfig_delegate PassConfig `gson:"delegate"`
}

func (this_ *ConformancePassConfig) BindSelf_(s any) {
	_ = s
	this_.PassConfig_PassConfigDelegateBase.BindSelf_(s)
}
func (this_ *ConformancePassConfig) IsJavaNil_() bool { return this_ == nil }
func (*ConformancePassConfig) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ConformancePassConfig", "ConformancePassConfig" }
func (this_ *ConformancePassConfig) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassConfig_PassConfigDelegateBase.IdentityHash_()
}
func (*ConformancePassConfig) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConformancePassConfig", "com.google.javascript.jscomp.ConformancePassConfig", func(v_ any) bool {
		return j2gort.InstanceOf[*ConformancePassConfig](v_)
	})
}
func (*ConformancePassConfig) Implements_jscomp_PassConfig_PassConfigDelegate() {}
func (*ConformancePassConfig) Implements_jscomp_PassConfig() {}

func (this_ *ConformancePassConfig) initFields_ConformancePassConfig() {
	_ = this_
}

// Java constructor: ConformancePassConfig(com.google.javascript.jscomp.PassConfig delegate)
func NewConformancePassConfig(delegate PassConfig) *ConformancePassConfig {
	this_ := &ConformancePassConfig{}
	this_.BindSelf_(this_)
	InitConformancePassConfig(this_, delegate)
	return this_
}

func InitConformancePassConfig(this_ *ConformancePassConfig, delegate PassConfig) {
	_ = delegate
	InitPassConfig_PassConfigDelegate(&this_.PassConfig_PassConfigDelegateBase, delegate)
	this_.initFields_ConformancePassConfig()
	this_.f_ConformancePassConfig_delegate = delegate
}

// Java method: ConformancePassConfig.getChecks()
func (this_ *ConformancePassConfig) GetChecks() *PassListBuilder {
	var fromDelegate *PassListBuilder = this_.f_ConformancePassConfig_delegate.GetChecks()
	_ = fromDelegate
	var val_1 PassListBuilder
	InitPassListBuilder(&val_1, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_1
	_ = passes
	passes.AddAllUpTo(fromDelegate, "checkConformance")
	return passes
}

// Java method: ConformancePassConfig.getOptimizations()
func (this_ *ConformancePassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(this_.GetPassConfigBase().options)
}

