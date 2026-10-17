// Port of com/google/javascript/jscomp/testing/JSErrorSubject.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// JSErrorSubject is a Truth Subject for the JSError class. Usage:
//
//	JSErrorSubject_AssertError(error).HasType(TypeValidator_TYPE_MISMATCH)
//
// TODO(tbreisacher): Add assertions on the message text, line number, etc.
type JSErrorSubject struct {
	*Subject
	actual *JSError
}

// JSErrorSubject_AssertError is JSErrorSubject.assertError(JSError).
//
// Java: JSErrorSubject.assertError(JSError)
func JSErrorSubject_AssertError(error_ *JSError) *JSErrorSubject {
	return NewJSErrorSubject(&Subject{actual: error_}, error_)
}

// AssertError is JSErrorSubject_AssertError for Java's
// `import static ...JSErrorSubject.assertError` (Go-only: language-adapter).
func AssertError(error_ *JSError) *JSErrorSubject {
	return JSErrorSubject_AssertError(error_)
}

// Java: JSErrorSubject(FailureMetadata, JSError)
func NewJSErrorSubject(s *Subject, error_ *JSError) *JSErrorSubject {
	return &JSErrorSubject{Subject: s, actual: error_}
}

func (s *JSErrorSubject) check(what string) *StandardSubjectBuilder {
	msg := "value of: jSError." + what + "\njSError was: " + show(s.actual)
	if s.message != "" {
		msg = s.message + "\n" + msg
	}
	return &StandardSubjectBuilder{message: msg}
}

// Java: JSErrorSubject.hasType(DiagnosticType)
func (s *JSErrorSubject) HasType(type_ *DiagnosticType) {
	s.check("getType()").That(s.actual.Type()).IsEqualTo(type_)
}

// Java: JSErrorSubject.hasMessage(String)
func (s *JSErrorSubject) HasMessage(msg string) {
	s.check("description").That(s.actual.Description()).IsEqualTo(msg)
}
