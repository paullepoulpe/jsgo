// Port of com/google/javascript/jscomp/UniqueIdSupplier.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/IdGenerator.java @ closure-compiler v20261021

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
	counts map[int32]int
	order  []int32
}

// Java: UniqueIdSupplier()
func NewUniqueIdSupplier() *UniqueIdSupplier {
	return &UniqueIdSupplier{}
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
	if len(s.order) == 0 {
		return nil
	}
	result := make([]*UniqueIdProto, len(s.order))
	for i, element := range s.order {
		result[i] = UniqueIdProto_NewBuilder().
			SetHash(element).
			SetCounter(int32(s.counts[element])).
			Build()
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

func (s *UniqueIdSupplier) addToCounter(element int32, occurrences int) int {
	if occurrences < 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"occurrences cannot be negative: " + strconv.Itoa(occurrences)))
	}
	oldCount, exists := s.counts[element]
	if occurrences > 0 {
		if s.counts == nil {
			s.counts = make(map[int32]int)
		}
		if !exists {
			s.order = append(s.order, element)
		}
		s.counts[element] = oldCount + occurrences
	}
	return oldCount
}

// IdGenerator is the id generator used in CrossChunkMethodMotion.
type IdGenerator struct {
	// Ids for cross-chunk method stubbing, so that each method has a unique
	// id.
	currentId int
}

// NewIdGenerator is new IdGenerator().
func NewIdGenerator() *IdGenerator { return &IdGenerator{} }

// HasGeneratedAnyIds returns whether we've generated any new ids.
//
// Java: IdGenerator.hasGeneratedAnyIds()
func (g *IdGenerator) HasGeneratedAnyIds() bool {
	return g.currentId != 0
}

// NewId creates a new id for stubbing a method.
//
// Java: IdGenerator.newId()
func (g *IdGenerator) NewId() int {
	id := g.currentId
	g.currentId++
	return id
}

// Java: IdGenerator.getCurrentId()
func (g *IdGenerator) GetCurrentId() int {
	return g.currentId
}

// Java: IdGenerator.setCurrentId(int)
func (g *IdGenerator) SetCurrentId(id int) {
	g.currentId = id
}
