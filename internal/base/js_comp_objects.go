// Port of com/google/javascript/jscomp/base/JSCompObjects.java @ closure-compiler v20260415

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// JSCompObjects holds basic Object functions used by JSComp.
type JSCompObjects struct{}

// JSCompObjects_Identical is reference equality (`a == b`).
//
// Go-only: language-adapter — Java `==` on references. Go's == on interface
// values compares pointers for reference types and values for Go value types
// (which stand in for Java enums and interned strings, where Java code
// relies on identity being value equality).
//
// Java: JSCompObjects.identical(Object, Object)
func JSCompObjects_Identical(a, b any) bool {
	return a == b
}

// NewJSCompObjects is the private constructor, which always throws.
//
// Java: JSCompObjects()
func NewJSCompObjects() *JSCompObjects {
	panic(jutil.NewAssertionErrorNoMessage())
}
