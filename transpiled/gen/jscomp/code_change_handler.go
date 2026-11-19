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
// Port of com/google/javascript/jscomp/CodeChangeHandler.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.CodeChangeHandler
type CodeChangeHandler interface {
	Implements_jscomp_CodeChangeHandler()
	IsJavaNil_() bool
	ReportChange()
}
func ToIface_CodeChangeHandler[T interface {
	CodeChangeHandler
	comparable
}](v T) CodeChangeHandler {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CodeChangeHandler[T CodeChangeHandler](v T) CodeChangeHandler {
	var vi CodeChangeHandler = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodeChangeHandler", "com.google.javascript.jscomp.CodeChangeHandler", func(v_ any) bool {
		return j2gort.InstanceOf[CodeChangeHandler](v_)
	})
}
