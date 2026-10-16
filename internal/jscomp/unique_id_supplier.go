// Port of com/google/javascript/jscomp/UniqueIdSupplier.java @ closure-compiler v20260415

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// UniqueIdSupplier generates unique ids across all compiler input source
// files.
type UniqueIdSupplier struct {
	// Go-only: language-adapter — Multiset<Integer> (a HashMultiset that is
	// only counted, never iterated) is a Go map of counts.
	counter map[int32]int
}

// Java: UniqueIdSupplier()
func NewUniqueIdSupplier() *UniqueIdSupplier {
	return &UniqueIdSupplier{counter: make(map[int32]int)}
}

// GetUniqueId creates and returns a unique Id across all compiler input
// source files.
//
// input is the compiler input for which the unique Id is requested. The
// result is the unique ID as String.
//
// Java: UniqueIdSupplier.getUniqueId(CompilerInput)
func (s *UniqueIdSupplier) GetUniqueId(input *CompilerInput) string {
	filePath := input.GetSourceFile().GetName()
	fileHashCode := jutil.StringHashCode(filePath)
	// Multiset.add(E, int) returns the count before the operation.
	id := s.counter[fileHashCode]
	s.counter[fileHashCode] = id + 1
	var fileHashString string
	if fileHashCode < 0 {
		fileHashString = "m" + strconv.Itoa(int(-fileHashCode))
	} else {
		fileHashString = strconv.Itoa(int(fileHashCode))
	}
	return fileHashString + "$" + strconv.Itoa(id)
}
