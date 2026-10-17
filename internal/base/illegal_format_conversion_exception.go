// Port of com/google/javascript/jscomp/base/format/IllegalFormatConversionException.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IllegalFormatConversionException is com.google.javascript.jscomp.base.format.IllegalFormatConversionException.
type IllegalFormatConversionException struct{ jutil.RuntimeException }

func (*IllegalFormatConversionException) isIllegalFormatConversionException() {}

// IIllegalFormatConversionException matches IllegalFormatConversionException and its subclasses in jutil.Try.
type IIllegalFormatConversionException interface {
	jutil.IRuntimeException
	isIllegalFormatConversionException()
}

// NewIllegalFormatConversionException constructs the exception.
//
// Java: IllegalFormatConversionException(char, Class<? extends Object>)
//
// Go-only: language-adapter — the Class argument is passed as its
// Class.toString() form (e.g. "class java.lang.String").
func NewIllegalFormatConversionException(conversionType uint16, class1 string) *IllegalFormatConversionException {
	e := &IllegalFormatConversionException{}
	m := "ConversionType = " + jutil.StringFromChar(conversionType) + " class = " + class1
	jutil.InitThrowable(e, "com.google.javascript.jscomp.base.format.IllegalFormatConversionException", &m, nil)
	return e
}
