// Port of com/google/javascript/jscomp/base/JSCompStrings.java @ closure-compiler v20261015

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

package base

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSCompStrings holds basic String functions used by JSComp.
type JSCompStrings struct{}

// JSCompStrings_Lines joins the lines with "\n".
//
// Java: JSCompStrings.lines(String...)
func JSCompStrings_Lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

// JSCompStrings_LinesOfIterable joins the lines with "\n".
//
// Java: JSCompStrings.lines(Iterable<String>)
func JSCompStrings_LinesOfIterable(lines []string) string {
	return strings.Join(lines, "\n")
}

// NewJSCompStrings is the private constructor, which always throws.
//
// Java: JSCompStrings()
func NewJSCompStrings() *JSCompStrings {
	panic(jutil.NewAssertionErrorNoMessage())
}
