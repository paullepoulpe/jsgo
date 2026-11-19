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
// Port of com/google/javascript/jscomp/CompilerInputProvider.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java interface: com.google.javascript.jscomp.CompilerInputProvider
type CompilerInputProvider interface {
	Implements_jscomp_CompilerInputProvider()
	IsJavaNil_() bool
	GetInput(*pkg_rhino.InputId) *CompilerInput
}
func ToIface_CompilerInputProvider[T interface {
	CompilerInputProvider
	comparable
}](v T) CompilerInputProvider {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CompilerInputProvider[T CompilerInputProvider](v T) CompilerInputProvider {
	var vi CompilerInputProvider = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerInputProvider", "com.google.javascript.jscomp.CompilerInputProvider", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerInputProvider](v_)
	})
}
type lambda_CompilerInputProvider struct {
	fn func(*pkg_rhino.InputId) *CompilerInput
}

func (l *lambda_CompilerInputProvider) IsJavaNil_() bool { return l == nil }
func NewLambda_CompilerInputProvider(fn func(*pkg_rhino.InputId) *CompilerInput) CompilerInputProvider {
	if fn == nil {
		return nil
	}
	return &lambda_CompilerInputProvider{fn: fn}
}

func (*lambda_CompilerInputProvider) Implements_jscomp_CompilerInputProvider() {}
func (l *lambda_CompilerInputProvider) GetInput(p0_ *pkg_rhino.InputId) *CompilerInput {
	return l.fn(p0_)
}

