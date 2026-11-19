/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/type/ReverseAbstractInterpreter.java (generated mechanically by j2go).

package jscomptype

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java interface: com.google.javascript.jscomp.type.ReverseAbstractInterpreter
type ReverseAbstractInterpreter interface {
	Implements_jscomptype_ReverseAbstractInterpreter()
	IsJavaNil_() bool
	GetPreciserScopeKnowingConditionOutcome(pkg_rhino.Node, FlowScope, *pkg_rhino.Outcome) FlowScope
}
func ToIface_ReverseAbstractInterpreter[T interface {
	ReverseAbstractInterpreter
	comparable
}](v T) ReverseAbstractInterpreter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.ReverseAbstractInterpreter", "com.google.javascript.jscomp.type.ReverseAbstractInterpreter", func(v_ any) bool {
		return j2gort.InstanceOf[ReverseAbstractInterpreter](v_)
	})
}
