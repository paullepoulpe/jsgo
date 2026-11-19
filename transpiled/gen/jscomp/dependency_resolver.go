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
// Port of com/google/javascript/jscomp/deps/DependencyResolver.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.deps.DependencyResolver
type DependencyResolver interface {
	Implements_jscomp_DependencyResolver()
	IsJavaNil_() bool
	GetDependencies_1_string(string) any
	GetDependencies_1_Ext_java_util_Collection(j2gort.Ext_java_util_Collection) any
	GetDependencies_3(string, j2gort.Ext_java_util_Set, bool) any
	GetDependencies_2(j2gort.Ext_java_util_Collection, j2gort.Ext_java_util_Set) any
}
func ToIface_DependencyResolver[T interface {
	DependencyResolver
	comparable
}](v T) DependencyResolver {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_DependencyResolver[T DependencyResolver](v T) DependencyResolver {
	var vi DependencyResolver = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DependencyResolver", "com.google.javascript.jscomp.deps.DependencyResolver", func(v_ any) bool {
		return j2gort.InstanceOf[DependencyResolver](v_)
	})
}
