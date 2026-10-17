// Port of com/google/javascript/jscomp/base/format/IllegalFormatCodePointException.java @ closure-compiler v20260415

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package base

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatCodePointException is com.google.javascript.jscomp.base.format.IllegalFormatCodePointException.
type IllegalFormatCodePointException struct{ jutil.RuntimeException }

func (*IllegalFormatCodePointException) isIllegalFormatCodePointException() {}

// IIllegalFormatCodePointException matches IllegalFormatCodePointException and its subclasses in jutil.Try.
type IIllegalFormatCodePointException interface {
	jutil.IRuntimeException
	isIllegalFormatCodePointException()
}

// NewIllegalFormatCodePointException constructs the exception.
//
// Java: IllegalFormatCodePointException(int)
func NewIllegalFormatCodePointException(codePoint int) *IllegalFormatCodePointException {
	e := &IllegalFormatCodePointException{}
	m := "codePoint:" + strconv.Itoa(codePoint)
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatCodePointException", &m, nil)
	return e
}
