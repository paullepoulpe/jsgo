/*
 * Copyright 2012 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ErrorHandler.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.ErrorHandler
type ErrorHandler interface {
	Implements_jscomp_ErrorHandler()
	IsJavaNil_() bool
	Report_2(*CheckLevel, *JSError)
}
func ToIface_ErrorHandler[T interface {
	ErrorHandler
	comparable
}](v T) ErrorHandler {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ErrorHandler", "com.google.javascript.jscomp.ErrorHandler", func(v_ any) bool {
		return j2gort.InstanceOf[ErrorHandler](v_)
	})
}
type lambda_ErrorHandler struct {
	fn func(*CheckLevel, *JSError)
}

func (l *lambda_ErrorHandler) IsJavaNil_() bool { return l == nil }
func NewLambda_ErrorHandler(fn func(*CheckLevel, *JSError)) ErrorHandler {
	if fn == nil {
		return nil
	}
	return &lambda_ErrorHandler{fn: fn}
}

func (*lambda_ErrorHandler) Implements_jscomp_ErrorHandler() {}
func (l *lambda_ErrorHandler) Report_2(p0_ *CheckLevel, p1_ *JSError) {
	l.fn(p0_, p1_)
}

