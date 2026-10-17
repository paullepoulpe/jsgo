// Port of com/google/javascript/jscomp/CompilerOptionsPreprocessor.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CompilerOptionsPreprocessor checks for combinations of options that are
// incompatible, i.e. will produce incorrect code.
//
// This is run by Compiler#compileInternal, which is not run during unit
// tests. The catch is that it's run after Compiler#initOptions, so if for
// example you want to change the warningsGuard, you can't do it here.
//
// Also, turns off options if the provided options don't make sense
// together.
//
// Don't instantiate (Java: private CompilerOptionsPreprocessor()).
type CompilerOptionsPreprocessor struct{}

// Java: CompilerOptionsPreprocessor.preprocess(CompilerOptions)
func CompilerOptionsPreprocessor_Preprocess(options *CompilerOptions) {
	if options.GetInlineFunctionsLevel() == CompilerOptions_Reach_NONE &&
		options.GetMaxFunctionSizeAfterInlining() != CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING {
		panic(newCompilerOptionsPreprocessor_InvalidOptionsException(
			"max_function_size_after_inlining has no effect if inlining is disabled."))
	}
}

// CompilerOptionsPreprocessor_InvalidOptionsException is an exception to
// indicate incompatible options in the CompilerOptions.
type CompilerOptionsPreprocessor_InvalidOptionsException struct {
	jutil.RuntimeException
}

func (*CompilerOptionsPreprocessor_InvalidOptionsException) invalidOptionsException() {}

// ICompilerOptionsPreprocessor_InvalidOptionsException matches
// CompilerOptionsPreprocessor.InvalidOptionsException and its subclasses.
type ICompilerOptionsPreprocessor_InvalidOptionsException interface {
	jutil.IRuntimeException
	invalidOptionsException()
}

// Java: CompilerOptionsPreprocessor.InvalidOptionsException(String)
func newCompilerOptionsPreprocessor_InvalidOptionsException(message string) *CompilerOptionsPreprocessor_InvalidOptionsException {
	e := &CompilerOptionsPreprocessor_InvalidOptionsException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.CompilerOptionsPreprocessor$InvalidOptionsException", &message, nil)
	return e
}
