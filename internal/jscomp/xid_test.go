// Port of com/google/javascript/jscomp/XidTest.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"math"
	"regexp"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TestXid_Uniqueness verifies that Xid#toString generates unique strings
// for the integers close to zero, Integer.MIN_VALUE, and
// Integer.MAX_VALUE.
func TestXid_Uniqueness(t *testing.T) {
	m := map[string]int32{}
	xidHelpTestUniqueness(t, m, -1000, 1000)
	xidHelpTestUniqueness(t, m, math.MinInt32, math.MinInt32+1000)
	xidHelpTestUniqueness(t, m, math.MaxInt32, math.MaxInt32-1000)
}

func xidHelpTestUniqueness(t *testing.T, m map[string]int32, lo, hi int32) {
	for i := int64(lo); i <= int64(hi); i++ {
		key := Xid_ToString(int32(i))
		if dup, ok := m[key]; ok {
			t.Errorf("Both %d and %d map to: %s", dup, i, key)
		}
		m[key] = int32(i)
	}
}

// TestXid_Length verifies that Xid#toString generates strings with lengths
// between 1 and 6.
func TestXid_Length(t *testing.T) {
	xidHelpTestLength(t, 0)
	xidHelpTestLength(t, 1)
	xidHelpTestLength(t, -1)
	xidHelpTestLength(t, math.MinInt32)
	xidHelpTestLength(t, math.MaxInt32)

	r := jutil.NewRandom()
	for i := 0; i <= 10000; i++ {
		n := r.NextInt()
		xidHelpTestLength(t, n)
	}
}

var xidLengthRE = regexp.MustCompile(`^.{1,6}$`)

func xidHelpTestLength(t *testing.T, i int32) {
	if s := Xid_ToString(i); !xidLengthRE.MatchString(s) {
		t.Errorf("Xid.toString(%d) = %q", i, s)
	}
}

func TestXid_ToString(t *testing.T) {
	assertEqual(t, "1", Xid_ToString(1), "z6ArXc")
	assertEqual(t, "2", Xid_ToString(2), "A6ArXc")
	assertEqual(t, "-1", Xid_ToString(-1), "x6ArXc")
	assertEqual(t, "10000", Xid_ToString(10000), "OcErXc")
	assertEqual(t, "-1951591049", Xid_ToString(-1951591049), "tTaYp")
}

func TestXid_Get(t *testing.T) {
	dummyMap := NewXid()
	assertEqual(t, "today", dummyMap.Get("today"), "nZzm6c")
	assertEqual(t, "tomorrow", dummyMap.Get("tomorrow"), "fkPKBb")
	assertEqual(t, "value", dummyMap.Get("value"), "b6Lt6c")

	assertEqual(t, "foo", dummyMap.Get("foo"), "QB6rXc")
	assertEqual(t, "foo.Bar", dummyMap.Get("foo.Bar"), "RW4o4b")

	assertEqual(t, "prop1", dummyMap.Get("prop1"), "MiB45c")
	assertEqual(t, "prop2", dummyMap.Get("prop2"), "NiB45c")
}
