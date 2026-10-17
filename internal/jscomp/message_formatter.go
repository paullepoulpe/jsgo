// Port of com/google/javascript/jscomp/MessageFormatter.java @ closure-compiler v20260415

/*
 * Copyright 2007 The Closure Compiler Authors.
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

// MessageFormatter formats warnings and errors. This interface may be used
// by implementations of ErrorManager to request message formatting
// capabilities.
type MessageFormatter interface {
	// FormatError formats an error.
	FormatError(error *JSError) string

	// FormatWarning formats a warning.
	FormatWarning(warning *JSError) string
}
