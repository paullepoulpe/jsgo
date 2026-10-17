// Port of com/google/javascript/jscomp/deps/SourceCodeEscapersTest.java @ closure-compiler v20261015

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestSourceCodeEscapers_GoogModuleCode(t *testing.T) {
	input := "goog.module('Foo');\nclass Foo {}"
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(input)
	assertEqual(t, "escaped", escaped, "goog.module(\\x27Foo\\x27);\\nclass Foo {}")
}

func TestSourceCodeEscapers_ExplicitCharacterEscaping(t *testing.T) {
	input := []uint16{'\'', '"', '<', '=', '>', '&', '\b', '\t', '\n', '\f', '\r', '\\'}
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(jutil.StringFromUTF16(input))
	assertEqual(t, "escaped", escaped, "\\x27\\x22\\x3c\\x3d\\x3e\\x26\\b\\t\\n\\f\\r\\\\")
}

func TestSourceCodeEscapers_LowUnicodeValueEscaping(t *testing.T) {
	input := []uint16{0x00, 0x10, 0x1F, 0x20, 0x7E, 0x7F, 0xC0, 0xFF}
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(jutil.StringFromUTF16(input))
	assertEqual(t, "escaped", escaped, "\\x00\\x10\\x1f ~\\x7f\\xc0\\xff")
}

func TestSourceCodeEscapers_HighUnicodeValueEscaping(t *testing.T) {
	input := []uint16{0x100, 0xF00, 0xFFFF}
	escaped := SourceCodeEscapers_JavascriptEscaper().Escape(jutil.StringFromUTF16(input))
	assertEqual(t, "escaped", escaped, "\\u0100\\u0f00\\uffff")
}
