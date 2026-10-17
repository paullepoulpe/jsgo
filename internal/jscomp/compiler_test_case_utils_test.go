// Port of com/google/javascript/jscomp/CompilerTestCaseUtils.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CompilerTestCaseUtils.java)

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

import (
	"bytes"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CompilerTestCaseUtils are CompilerTestCase utilities

// Java: CompilerTestCaseUtils.multistageSerializeAndDeserialize(CompilerTestCase, Compiler, List<SourceFile>, List<SourceFile>, CodeChangeHandler)
func CompilerTestCaseUtils_MultistageSerializeAndDeserialize(
	testCase CompilerTestCase,
	compiler *Compiler,
	externs []*SourceFile,
	inputs []*SourceFile,
	changeHandler CodeChangeHandler) *Compiler {
	NewRemoveCastNodes(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	errorManager := compiler.GetErrorManager()
	jutil.Try(func() {
		baos := &bytes.Buffer{}
		compiler.GetChangeTracker().RemoveChangeHandler(changeHandler)
		compiler.DisableThreads()
		compiler.SaveState(baos)

		bais := bytes.NewReader(baos.Bytes())
		compiler = testCase.CreateCompiler()
		compiler.DisableThreads()
		compiler.Init(externs, inputs, testCase.GetOptions())
		compiler.RestoreState(bais)
		compiler.SetErrorManager(errorManager)
		compiler.GetChangeTracker().AddChangeHandler(changeHandler)
	}, func(e jutil.Throwable) {
		// Java: catch (Exception e), which does not see java.lang.Errors.
		if _, isError := e.(jutil.IError); isError {
			panic(e)
		}
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return compiler
}

// Java: CompilerTestCaseUtils.setDebugLogDirectoryOn(CompilerOptions)
func CompilerTestCaseUtils_SetDebugLogDirectoryOn(options *CompilerOptions) {
}
