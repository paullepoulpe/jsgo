// Port of com/google/javascript/jscomp/FlagUsageException.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// FlagUsageException is an exception thrown when command-line flags are used
// incorrectly.
type FlagUsageException struct{ jutil.RuntimeException }

func (*FlagUsageException) flagUsageException() {}

// IFlagUsageException matches FlagUsageException.
type IFlagUsageException interface {
	jutil.IRuntimeException
	flagUsageException()
}

// Java: FlagUsageException.serialVersionUID
const FlagUsageException_serialVersionUID int64 = 1

// Java: FlagUsageException(String)
func NewFlagUsageException(message string) *FlagUsageException {
	e := &FlagUsageException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.FlagUsageException", &message, nil)
	return e
}
