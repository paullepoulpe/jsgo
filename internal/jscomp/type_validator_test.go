// Port of test/com/google/javascript/jscomp/TypeValidatorTest.java @ closure-compiler v20260415

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

package jscomp

import "testing"

// TypeValidatorTest extends CompilerTestCase with enableTypeCheck(): every
// case type-checks source code (TypedScopeCreator, TypeInference, TypeCheck)
// and inspects the reported warnings and Compiler.getTypeMismatches(). They
// are ported when CompilerTestCase and the type-checking passes are; the
// TypeValidator API itself is covered by type_validator_oracle_test.go.

func TestTypeValidator_BasicMismatch(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_FunctionMismatch(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_FunctionMismatch2(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_FunctionMismatchMediumLengthTypes(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_MissingPropertiesPrintedWithMismatch_defaultParam(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_MissingPropertiesPrintedWithMismatch_nullableParam(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_FunctionMismatchLongTypes(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_FunctionMismatchTypedef(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_Bug_testMismatchRecursively_throughFields(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_Bug_testMismatchRecursively_throughTemplates(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_NullUndefined(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_Subclass(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_UnionsMismatch(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef1(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef2(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef3(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef4(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef5(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef6(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef7(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef8(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef9(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef10(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef11(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef12(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_ModuloNullUndef13(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_InheritanceModuloNullUndef1(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_InheritanceModuloNullUndef2(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_InheritanceModuloNullUndef3(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_OptionalProperties_dontNeedInitializer_inConstructor(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_DuplicateSuppression(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_DuplicateSuppression_class(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_DuplicateSuppression_typeMismatch(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_DuplicateSuppression_stubs(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_DuplicateSuppression_topLevelVariables(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}

func TestTypeValidator_DuplicateSuppression_topLevelFunctions(t *testing.T) {
	t.Skip("needs CompilerTestCase, TypeCheck")
}
