/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/NameGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.NameGenerator
type NameGenerator interface {
	Implements_jscomp_NameGenerator()
	IsJavaNil_() bool
	Reset_3(j2gort.Ext_java_util_Set, string, j2gort.Ext_java_util_Set)
	Reset_4(j2gort.Ext_java_util_Set, string, j2gort.Ext_java_util_Set, j2gort.Ext_java_util_Set)
	Clone(j2gort.Ext_java_util_Set, string, j2gort.Ext_java_util_Set) NameGenerator
	GenerateNextName() string
}
func ToIface_NameGenerator[T interface {
	NameGenerator
	comparable
}](v T) NameGenerator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NameGenerator", "com.google.javascript.jscomp.NameGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[NameGenerator](v_)
	})
}
