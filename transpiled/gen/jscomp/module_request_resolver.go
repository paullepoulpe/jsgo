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
// Port of com/google/javascript/jscomp/modules/ModuleRequestResolver.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.modules.ModuleRequestResolver
type ModuleRequestResolver interface {
	Implements_jscomp_ModuleRequestResolver()
	IsJavaNil_() bool
	Resolve_1_Import(*Import) UnresolvedModule
	Resolve_1_Export(*Export) UnresolvedModule
}
func ToIface_ModuleRequestResolver[T interface {
	ModuleRequestResolver
	comparable
}](v T) ModuleRequestResolver {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ModuleRequestResolver[T ModuleRequestResolver](v T) ModuleRequestResolver {
	var vi ModuleRequestResolver = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleRequestResolver", "com.google.javascript.jscomp.modules.ModuleRequestResolver", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleRequestResolver](v_)
	})
}
