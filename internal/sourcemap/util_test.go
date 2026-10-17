// Port of test/com/google/debugging/sourcemap/UtilTest.java @ closure-compiler v20260415

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

package sourcemap_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

func utilTestAppendHexJavaScriptRepresentation(t *testing.T, ch uint16, expectedOut string) {
	t.Helper()
	var sb strings.Builder
	sourcemap.Util_AppendHexJavaScriptRepresentation(&sb, ch)
	if sb.String() != expectedOut {
		t.Errorf("appendHexJavaScriptRepresentation(%q) = %q, want %q", rune(ch), sb.String(), expectedOut)
	}
}

func TestAppendHexJavaScriptRepresentation(t *testing.T) {
	utilTestAppendHexJavaScriptRepresentation(t, 'a', "\\u0061")
	utilTestAppendHexJavaScriptRepresentation(t, 'z', "\\u007a")
	utilTestAppendHexJavaScriptRepresentation(t, 0, "\\u0000")
	utilTestAppendHexJavaScriptRepresentation(t, 0xa1, "\\u00a1")
}
