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
// Port of com/google/javascript/jscomp/graph/UnionFind.java (generated mechanically by j2go).

package graph

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.graph.UnionFind
type UnionFind[E comparable] interface {
	Implements_graph_UnionFind()
	IsJavaNil_() bool
	Add(E)
	Union(E, E) E
	Find(E) E
	AreEquivalent(E, E) bool
	Elements() j2gort.Ext_java_util_Set
	AllEquivalenceClasses() j2gort.Ext_com_google_common_collect_ImmutableList
	FindAll(E) j2gort.Ext_java_util_Set
}
func ToIface_UnionFind[E comparable, T_iface_ interface {
	UnionFind[E]
	comparable
}](v T_iface_) UnionFind[E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.UnionFind", "com.google.javascript.jscomp.graph.UnionFind", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_UnionFind()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
