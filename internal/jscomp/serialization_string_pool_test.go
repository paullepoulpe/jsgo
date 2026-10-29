// Port of test/com/google/javascript/jscomp/serialization/StringPoolTest.java @ closure-compiler v20261021

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

package jscomp

import (
	"bytes"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestStringPool_EmptyPool_fromBuilderOrConstant(t *testing.T) {
	jsaAssertIsEmptyPool(t, StringPool_Empty())
	jsaAssertIsEmptyPool(t, StringPool_Builder_().Build())
}

func TestStringPool_EmptyPool_serializesToDefaultProto(t *testing.T) {
	got := StringPool_Empty().ToProto()
	jsaTrue(t, bytes.Equal(got.ToByteArray(), StringPoolProto_GetDefaultInstance().ToByteArray()))
	jsaEq(t, got.String(), StringPoolProto_GetDefaultInstance().String())
}

func TestStringPool_EmptyPool_parsesFromDefaultProto(t *testing.T) {
	jsaAssertIsEmptyPool(t, StringPool_FromProto(StringPoolProto_GetDefaultInstance()))
}

func jsaAssertIsEmptyPool(t *testing.T, pool *StringPool) {
	t.Helper()
	jsaEq(t, pool.Get(0), "")
	jsaThrows[jutil.Throwable](t, func() { pool.Get(1) })
}

// Go-only: extra coverage — strings round-trip through the proto.
func TestStringPool_RoundTrip(t *testing.T) {
	b := StringPool_Builder_()
	a := b.Put("a")
	x := b.Put(jutil.StringFromUTF16([]uint16{'x', 0xd800}))
	jsaEq(t, b.Put("a"), a)
	jsaEq(t, b.Put(""), 0)
	pool := StringPool_FromProto(jsaReparseStringPool(b.Build().ToProto()))
	jsaEq(t, pool.Get(a), "a")
	jsaEq(t, pool.Get(x), jutil.StringFromUTF16([]uint16{'x', 0xd800}))
}

func jsaReparseStringPool(p *StringPoolProto) *StringPoolProto {
	return StringPoolProto_ParseFrom(p.ToByteArray())
}
