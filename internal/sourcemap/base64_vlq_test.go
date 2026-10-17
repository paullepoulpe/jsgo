// Port of test/com/google/debugging/sourcemap/Base64VLQTest.java @ closure-compiler v20260415

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
	"math"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

func TestBase64VLQSelectedValues1(t *testing.T) {
	for i := int32(0); i < 63; i++ {
		base64VLQTestValue(t, i)
	}
}

func TestBase64VLQSelectedValues2(t *testing.T) {
	base := int32(1)
	for i := 0; i < 30; i++ {
		base64VLQTestValue(t, base-1)
		base64VLQTestValue(t, base)
		base *= 2
	}
}

func TestBase64VLQSelectedSignedValues1(t *testing.T) {
	for i := int32(-(64*64 - 1)); i < (64*64 - 1); i++ {
		base64VLQTestValue(t, i)
	}
}

func TestBase64VLQSelectedSignedValues2(t *testing.T) {
	base := int32(1)
	for i := 0; i < 30; i++ {
		base64VLQTestValue(t, base-1)
		base64VLQTestValue(t, base)
		base *= 2
	}
	base = -1
	for i := 0; i < 30; i++ {
		base64VLQTestValue(t, base-1)
		base64VLQTestValue(t, base)
		base *= 2
	}
}

func TestBase64VLQSelectedSignedValues3(t *testing.T) {
	base64VLQTestValue(t, math.MaxInt32)
	base64VLQTestValue(t, math.MaxInt32-1)
	base64VLQTestValue(t, math.MaxInt32-2)
	base64VLQTestValue(t, 0x70000000)
	base64VLQTestValue(t, math.MinInt32)
	base64VLQTestValue(t, math.MinInt32+1)
	base64VLQTestValue(t, math.MinInt32+2)
}

// base64VLQTest_CharIteratorImpl is Base64VLQTest.CharIteratorImpl.
type base64VLQTest_CharIteratorImpl struct {
	current int
	length  int
	cs      []uint16
}

func (c *base64VLQTest_CharIteratorImpl) set(sb string) {
	c.current = 0
	c.cs = jutil.StringToUTF16(sb)
	c.length = len(c.cs)
}

func (c *base64VLQTest_CharIteratorImpl) HasNext() bool {
	return c.current < c.length
}

func (c *base64VLQTest_CharIteratorImpl) Next() uint16 {
	ch := c.cs[c.current]
	c.current++
	return ch
}

func base64VLQTestValue(t *testing.T, value int32) {
	t.Helper()
	var sb strings.Builder
	sourcemap.Base64VLQ_Encode(&sb, value)
	ci := &base64VLQTest_CharIteratorImpl{}
	ci.set(sb.String())
	result := sourcemap.Base64VLQ_Decode(ci)
	if result != value {
		t.Fatalf("failed for value %d: got %d", value, result)
	}
}
