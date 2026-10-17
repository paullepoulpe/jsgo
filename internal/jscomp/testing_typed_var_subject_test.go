// Port of com/google/javascript/jscomp/testing/TypedVarSubject.java @ closure-compiler v20261015
// (src/com/google/javascript/jscomp/testing/TypedVarSubject.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// Go-only: test-helper — jscomp.testing classes are ported as test files of
// package jscomp (see testing_js_error_subject_test.go).

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// TypedVarSubject is a Truth Subject for the TypedVar class. Usage:
//
//	TypedVarSubject_AssertThat(var_).IsInferred()
//	TypedVarSubject_AssertThat(var_).HasJSTypeThat().IsString()
type TypedVarSubject struct {
	*Subject
	actual *TypedVar
}

// TypedVarSubject_AssertThat is TypedVarSubject.assertThat(TypedVar).
//
// Java: TypedVarSubject.assertThat(TypedVar)
func TypedVarSubject_AssertThat(var_ *TypedVar) *TypedVarSubject {
	// NB: Eclipse's Java compiler bails on just passing TypedVarSubject::new below, so wrap it in a
	// Closure.
	return newTypedVarSubject(&Subject{actual: var_}, var_)
}

// Java: TypedVarSubject(FailureMetadata, TypedVar)
func newTypedVarSubject(failureMetadata *Subject, var_ *TypedVar) *TypedVarSubject {
	return &TypedVarSubject{Subject: failureMetadata, actual: var_}
}

// Java: TypedVarSubject.actualNonNull()
func (s *TypedVarSubject) actualNonNull() *TypedVar {
	s.IsNotNull()
	return s.actual
}

// check is Subject.check(format, args): a subject for a property of the
// var.
func (s *TypedVarSubject) check(format string, args ...any) *StandardSubjectBuilder {
	msg := "value of: var." + jutil.LenientFormat(format, args...)
	if s.message != "" {
		msg = s.message + "\n" + msg
	}
	return &StandardSubjectBuilder{message: msg}
}

// Java: TypedVarSubject.hasJSTypeThat()
func (s *TypedVarSubject) HasJSTypeThat() *rhinotesting.TypeSubject {
	b := s.check("getJSType()")
	return rhinotesting.NewTruth(assertionTB{}).AssertWithMessage("%s", b.message).AboutTypes().That(s.actualNonNull().GetType())
}

// Java: TypedVarSubject.isInferred()
func (s *TypedVarSubject) IsInferred() {
	s.check("isTypeInferred()").That(s.actualNonNull().IsTypeInferred()).IsTrue()
}

// Java: TypedVarSubject.isNotInferred()
func (s *TypedVarSubject) IsNotInferred() {
	s.check("isTypeInferred()").That(s.actualNonNull().IsTypeInferred()).IsFalse()
}
