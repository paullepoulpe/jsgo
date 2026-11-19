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
// Port of com/google/javascript/jscomp/MessageFormatter.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.MessageFormatter
type MessageFormatter interface {
	Implements_jscomp_MessageFormatter()
	IsJavaNil_() bool
	FormatError(*JSError) string
	FormatWarning(*JSError) string
}
func ToIface_MessageFormatter[T interface {
	MessageFormatter
	comparable
}](v T) MessageFormatter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_MessageFormatter[T MessageFormatter](v T) MessageFormatter {
	var vi MessageFormatter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MessageFormatter", "com.google.javascript.jscomp.MessageFormatter", func(v_ any) bool {
		return j2gort.InstanceOf[MessageFormatter](v_)
	})
}
