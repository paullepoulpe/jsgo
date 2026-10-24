// Port of com/google/javascript/jscomp/serialization/MalformedTypedAstException.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
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

// MalformedTypedAstException is emitted when deserialization sees a
// TypedAst with semantic errors.
//
// This excpetion means a TypedAst is flawed in some way. It is distinct from
// other kinds of expections because it describes a user error, rather than a
// bug in deserialization. Most likely, the tool that generated the
// serialized data is faulty.
//
// This exception says nothing about the byte-level encoding or parsing of
// the input data. Instead, it means that TypedAst parsed from those bytes
// can't be translated into the compiler's internal representation.
type MalformedTypedAstException struct{ jutil.RuntimeException }

func (*MalformedTypedAstException) malformedTypedAstException() {}

// IMalformedTypedAstException matches MalformedTypedAstException.
type IMalformedTypedAstException interface {
	jutil.IRuntimeException
	malformedTypedAstException()
}

// Java: MalformedTypedAstException.checkWellFormed(boolean, String)
func MalformedTypedAstException_CheckWellFormed(condition bool, description string) {
	if !condition {
		panic(NewMalformedTypedAstException(description))
	}
}

// Java: MalformedTypedAstException.checkWellFormed(boolean, String, Object)
func MalformedTypedAstException_CheckWellFormedWithParam(condition bool, description string, param any) {
	if !condition {
		message := description + ": " + jutil.ToJavaString(param)
		panic(NewMalformedTypedAstException(message))
	}
}

// Java: MalformedTypedAstException(Object)
func NewMalformedTypedAstException(msg any) *MalformedTypedAstException {
	e := &MalformedTypedAstException{}
	s := jutil.ToJavaString(jutil.CheckNotNull(msg))
	jutil.InitThrowable(e, "com.google.javascript.jscomp.serialization.MalformedTypedAstException", &s, nil)
	return e
}

// Java: MalformedTypedAstException(Object, Throwable)
func NewMalformedTypedAstExceptionWithCause(msg any, cause jutil.Throwable) *MalformedTypedAstException {
	e := &MalformedTypedAstException{}
	s := jutil.ToJavaString(jutil.CheckNotNull(msg))
	jutil.InitThrowable(e, "com.google.javascript.jscomp.serialization.MalformedTypedAstException", &s, cause)
	return e
}
