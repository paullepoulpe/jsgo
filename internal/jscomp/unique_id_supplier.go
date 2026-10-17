// Port of com/google/javascript/jscomp/UniqueIdSupplier.java @ closure-compiler v20261015

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

// UniqueIdSupplier generates unique String Ids when requested via a compiler
// instance.
//
// This supplier provides Ids that are deterministic and unique across all
// input files given to the compiler. The generated ID format is: uniqueId =
// "fileHashCode$counterForThisFile"
//
// Java: UniqueIdSupplier
type UniqueIdSupplier struct {
	// Go-only: language-adapter — Multiset<Integer> (a HashMultiset, whose
	// entrySet toProto iterates in the JDK HashMap order of its Integer
	// elements) is a JavaHashMap from element to count.
	counter *jutil.JavaHashMap[int32, int]
}

// Java: UniqueIdSupplier()
func NewUniqueIdSupplier() *UniqueIdSupplier {
	return &UniqueIdSupplier{counter: jutil.NewJavaHashMap[int32, int]()}
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
	id := s.addToCounter(fileHashCode, 1)
	var fileHashString string
	if fileHashCode < 0 {
		fileHashString = "m" + strconv.Itoa(int(-fileHashCode))
	} else {
		fileHashString = strconv.Itoa(int(fileHashCode))
	}
	return fileHashString + "$" + strconv.Itoa(id)
}

// Java: UniqueIdSupplier.toProto()
func (s *UniqueIdSupplier) ToProto() []*UniqueIdProto {
	var result []*UniqueIdProto
	for element, count := range s.counter.All() {
		result = append(result,
			UniqueIdProto_NewBuilder().
				SetHash(element).
				SetCounter(int32(count)).
				Build())
	}
	return result
}

// Java: UniqueIdSupplier.fromProto(List<UniqueIdProto>)
func UniqueIdSupplier_FromProto(protos []*UniqueIdProto) *UniqueIdSupplier {
	supplier := NewUniqueIdSupplier()
	for _, p := range protos {
		supplier.addToCounter(p.GetHash(), int(p.GetCounter()))
	}
	return supplier
}

// addToCounter is Multiset.add(E, int) on counter: adds occurrences
// occurrences of element and returns the count before the operation.
//
// Go-only: language-adapter — Guava's Multiset.add(E, int).
func (s *UniqueIdSupplier) addToCounter(element int32, occurrences int) int {
	jutil.CheckArgumentMsg1(occurrences >= 0, "occurrences cannot be negative: %s", occurrences)
	oldCount := s.counter.Get(element)
	if occurrences > 0 {
		s.counter.Put(element, oldCount+occurrences)
	}
	return oldCount
}
