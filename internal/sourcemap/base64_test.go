// Port of test/com/google/debugging/sourcemap/Base64Test.java @ closure-compiler v20261104

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
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

func TestBase64FromBase64OutOfBoundsAndInvalid(t *testing.T) {
	if got := sourcemap.Base64_FromBase64('!'); got != -1 {
		t.Fatalf("Base64_FromBase64('!') = %d, want -1", got)
	}
	for _, ch := range []uint16{256, 0x1234, 0xffff} {
		func() {
			defer func() {
				r := recover()
				if _, ok := r.(*jutil.ArrayIndexOutOfBoundsException); !ok {
					t.Fatalf("Base64_FromBase64(0x%x) panic = %T (%v), want *jutil.ArrayIndexOutOfBoundsException", ch, r, r)
				}
			}()
			sourcemap.Base64_FromBase64(ch)
		}()
	}

	ci := &base64VLQTest_CharIteratorImpl{}
	ci.set("\u0100")
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(*jutil.ArrayIndexOutOfBoundsException); !ok {
				t.Fatalf("Base64VLQ_Decode(\"\\u0100\") panic = %T (%v), want *jutil.ArrayIndexOutOfBoundsException", r, r)
			}
		}()
		sourcemap.Base64VLQ_Decode(ci)
	}()
}
