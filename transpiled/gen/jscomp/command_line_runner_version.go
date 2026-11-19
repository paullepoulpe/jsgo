/*
 * Copyright 2022 The Closure Compiler Authors.
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
// Port of CommandLineRunnerVersion.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.CommandLineRunnerVersion
type CommandLineRunnerVersion struct {
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunnerVersion) BindSelf_(s any) {
	_ = s
}
func (this_ *CommandLineRunnerVersion) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunnerVersion) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CommandLineRunnerVersion", "CommandLineRunnerVersion" }
func (this_ *CommandLineRunnerVersion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CommandLineRunnerVersion) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunnerVersion", "com.google.javascript.jscomp.CommandLineRunnerVersion", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunnerVersion](v_)
	})
}

func StaticGet_CommandLineRunnerVersion_COMPILER_VERSION() string { return "v20261104" }
func (this_ *CommandLineRunnerVersion) initFields_CommandLineRunnerVersion() {
	_ = this_
}

// Java constructor: CommandLineRunnerVersion()
func NewCommandLineRunnerVersion() *CommandLineRunnerVersion {
	this_ := &CommandLineRunnerVersion{}
	InitCommandLineRunnerVersion(this_)
	return this_
}

func InitCommandLineRunnerVersion(this_ *CommandLineRunnerVersion) {
	this_.initFields_CommandLineRunnerVersion()
}

