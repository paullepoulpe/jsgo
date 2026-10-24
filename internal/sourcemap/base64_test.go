// Port of test/com/google/debugging/sourcemap/Base64Test.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

func TestBase64(t *testing.T) {
	for i := int32(0); i < 64; i++ {
		base64TestValue(t, i)
	}
}

func TestBase64EncodeInt(t *testing.T) {
	for _, tc := range []struct {
		v    int32
		want string
	}{
		{0, "AAAAAA"},
		{1, "AAAAAQ"},
		{42, "AAAAKg"},
		{-100, "////nA"},
		{-1, "/////w"}, // 0xffffffff
	} {
		if got := sourcemap.Base64_Base64EncodeInt(tc.v); got != tc.want {
			t.Errorf("base64EncodeInt(%d) = %q, want %q", tc.v, got, tc.want)
		}
	}
}

func base64TestValue(t *testing.T, value int32) {
	t.Helper()
	if got := sourcemap.Base64_FromBase64(sourcemap.Base64_ToBase64(value)); got != value {
		t.Errorf("fromBase64(toBase64(%d)) = %d", value, got)
	}
}
