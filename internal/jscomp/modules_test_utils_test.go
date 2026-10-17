// Port of com/google/javascript/jscomp/ModulesTestUtils.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ModulesTestUtils.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// ModulesTestUtils holds test utilities for testing modules, used by
// Es6RewriteModulesTest and ProcessCommonJSModulesTest. (A final class with a
// private constructor and only static members: its members are package
// functions.)

// Java: ModulesTestUtils.testModules(CompilerTestCase, String, String, String)
func ModulesTestUtils_TestModules(test CompilerTestCase, fileName string, input string, expected string) {
	// Shared with ProcessCommonJSModulesTest.
	inputs := []*SourceFile{
		SourceFile_FromCode("other.js", "goog.provide('module$other');"),
		SourceFile_FromCode("yet_another.js", "goog.provide('module$yet_another');"),
		SourceFile_FromCode(fileName, input),
	}
	expecteds := []*SourceFile{
		SourceFile_FromCode("other.js", "goog.provide('module$other');"),
		SourceFile_FromCode("yet_another.js", "goog.provide('module$yet_another');"),
		SourceFile_FromCode(fileName, expected),
	}
	test.GetCompilerTestCaseBase().Test(CompilerTestCase_Srcs(inputs), CompilerTestCase_Expected_(expecteds))
}

// Java: ModulesTestUtils.testModulesError(CompilerTestCase, String, DiagnosticType)
func ModulesTestUtils_TestModulesError(test CompilerTestCase, input string, error_ *DiagnosticType) {
	inputs := []*SourceFile{
		SourceFile_FromCode("other.js", ""), SourceFile_FromCode("testcode.js", input),
	}
	test.GetCompilerTestCaseBase().TestError(CompilerTestCase_Srcs(inputs), CompilerTestCase_Error(error_))
}
