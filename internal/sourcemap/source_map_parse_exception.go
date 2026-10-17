// Port of com/google/debugging/sourcemap/SourceMapParseException.java @ closure-compiler v20260415

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

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapParseException is thrown if an invalid or unknown source map is
// encountered. It is a checked exception (java.lang.Exception).
type SourceMapParseException struct{ jutil.JavaException }

func (*SourceMapParseException) sourceMapParseException() {}

// ISourceMapParseException matches SourceMapParseException.
type ISourceMapParseException interface {
	jutil.IException
	sourceMapParseException()
}

// NewSourceMapParseException constructs a SourceMapParseException.
//
// Java: SourceMapParseException(String)
func NewSourceMapParseException(message string) *SourceMapParseException {
	e := &SourceMapParseException{}
	jutil.InitThrowable(e, "com.google.debugging.sourcemap.SourceMapParseException", &message, nil)
	return e
}
