// Port of com/google/javascript/jscomp/base/format/IllegalFormatWidthException.java @ closure-compiler v20261015

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

// IllegalFormatWidthException is com.google.javascript.jscomp.base.format.IllegalFormatWidthException.
type IllegalFormatWidthException struct{ jutil.RuntimeException }

func (*IllegalFormatWidthException) isIllegalFormatWidthException() {}

// IIllegalFormatWidthException matches IllegalFormatWidthException and its subclasses in jutil.Try.
type IIllegalFormatWidthException interface {
	jutil.IRuntimeException
	isIllegalFormatWidthException()
}

// NewIllegalFormatWidthException constructs the exception.
//
// Java: IllegalFormatWidthException(int)
func NewIllegalFormatWidthException(width int) *IllegalFormatWidthException {
	e := &IllegalFormatWidthException{}
	m := "Width:" + strconv.Itoa(width)
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatWidthException", &m, nil)
	return e
}
