/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CompilerPass.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java interface: com.google.javascript.jscomp.CompilerPass
type CompilerPass interface {
	Implements_jscomp_CompilerPass()
	IsJavaNil_() bool
	Process_2(pkg_rhino.Node, pkg_rhino.Node)
}
func ToIface_CompilerPass[T interface {
	CompilerPass
	comparable
}](v T) CompilerPass {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CompilerPass[T CompilerPass](v T) CompilerPass {
	var vi CompilerPass = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerPass", "com.google.javascript.jscomp.CompilerPass", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerPass](v_)
	})
}
type lambda_CompilerPass struct {
	fn func(pkg_rhino.Node, pkg_rhino.Node)
}

func (l *lambda_CompilerPass) IsJavaNil_() bool { return l == nil }
func NewLambda_CompilerPass(fn func(pkg_rhino.Node, pkg_rhino.Node)) CompilerPass {
	if fn == nil {
		return nil
	}
	return &lambda_CompilerPass{fn: fn}
}

func (*lambda_CompilerPass) Implements_jscomp_CompilerPass() {}
func (l *lambda_CompilerPass) Process_2(p0_ pkg_rhino.Node, p1_ pkg_rhino.Node) {
	l.fn(p0_, p1_)
}

