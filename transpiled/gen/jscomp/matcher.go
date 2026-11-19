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
// Port of com/google/javascript/refactoring/Matcher.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java interface: com.google.javascript.refactoring.Matcher
type Matcher interface {
	Implements_jscomp_Matcher()
	IsJavaNil_() bool
	Matches(pkg_rhino.Node, *NodeMetadata) bool
}
func ToIface_Matcher[T interface {
	Matcher
	comparable
}](v T) Matcher {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.Matcher", "com.google.javascript.refactoring.Matcher", func(v_ any) bool {
		return j2gort.InstanceOf[Matcher](v_)
	})
}
type lambda_Matcher struct {
	fn func(pkg_rhino.Node, *NodeMetadata) bool
}

func (l *lambda_Matcher) IsJavaNil_() bool { return l == nil }
func NewLambda_Matcher(fn func(pkg_rhino.Node, *NodeMetadata) bool) Matcher {
	if fn == nil {
		return nil
	}
	return &lambda_Matcher{fn: fn}
}

func (*lambda_Matcher) Implements_jscomp_Matcher() {}
func (l *lambda_Matcher) Matches(p0_ pkg_rhino.Node, p1_ *NodeMetadata) bool {
	return l.fn(p0_, p1_)
}

