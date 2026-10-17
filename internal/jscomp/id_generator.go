// Port of com/google/javascript/jscomp/IdGenerator.java @ closure-compiler v20260415

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
