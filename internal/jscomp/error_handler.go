// Port of com/google/javascript/jscomp/ErrorHandler.java @ closure-compiler v20261021

/*
 * Copyright 2012 The Closure Compiler Authors.
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

// ErrorHandler is any generic sink for warnings and errors, after they've
// passed through any filtering WarningsGuards.
type ErrorHandler interface {
	// Report reports error at the reporting level. level may be nil (Java
	// callers pass null).
	//
	// Java: ErrorHandler.report(CheckLevel, JSError)
	Report(level *CheckLevel, error *JSError)
}
