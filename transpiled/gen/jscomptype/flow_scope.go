/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/type/FlowScope.java (generated mechanically by j2go).

package jscomptype

import (
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java interface: com.google.javascript.jscomp.type.FlowScope
type FlowScope interface {
	Implements_jscomptype_FlowScope()
	IsJavaNil_() bool
	pkg_rhino.StaticTypedScope
	pkg_graph.LatticeElement
	WithSyntacticScope(pkg_rhino.StaticTypedScope) FlowScope
	InferSlotType(string, pkg_rhino.JSType) FlowScope
	InferQualifiedSlot(pkg_rhino.Node, string, pkg_rhino.JSType, pkg_rhino.JSType, bool) FlowScope
	GetDeclarationScope() pkg_rhino.StaticTypedScope
	EqualsFlowScope(FlowScope) bool
}

func ToIface_FlowScope[T interface {
	FlowScope
	comparable
}](v T) FlowScope {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_FlowScope[T FlowScope](v T) FlowScope {
	var vi FlowScope = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.type.FlowScope", "com.google.javascript.jscomp.type.FlowScope", func(v_ any) bool {
		return j2gort.InstanceOf[FlowScope](v_)
	})
}

// Java default interface method: FlowScope.equalsFlowScope(com.google.javascript.jscomp.type.FlowScope other)
func IfaceDefault_FlowScope_EqualsFlowScope(this_ FlowScope, other FlowScope) bool {
	_ = this_
	_ = other
	return ((this_ == other || ((this_ == nil || this_.IsJavaNil_()) && (other == nil || other.IsJavaNil_()))) || javalang.ObjectsEquals(j2gort.CheckNotNull(this_), other))
}
