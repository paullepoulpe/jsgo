// Port of com/google/javascript/jscomp/JSChunk.java @ closure-compiler v20261021

/*
 * Copyright 2005 The Closure Compiler Authors.
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
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSChunk is a JavaScript chunk: it has a unique name, consists of a list of
// compiler inputs, and can depend on other chunks.
//
// Go-only: language-adapter — readObject(ObjectInputStream), the
// java.io.Serializable hook that re-creates the transient input map after
// deserialization, is ReadObject(). JSChunk still declares it, but nothing
// deserializes a JSChunk any more: since v20261015 Compiler saves the chunk
// graph of its state as ChunkProto messages and rebuilds the chunks.
type JSChunk struct {
	// Chunk name
	name string

	// Source code inputs
	// non-final for deserialization
	// CompilerInputs must be explicitly added to the JSChunk again after deserialization
	// A map keyed by the CompilerInput.getName() to speed up getByName and removeByName.
	inputs       *jutil.LinkedHashMap[string, *CompilerInput]
	cachedInputs []*CompilerInput

	// Chunks that this chunk depends on
	deps []*JSChunk

	// The length of the longest path starting from this chunk
	depth int
	// The position of this chunk relative to all others in the AST.
	index int
}

var _ DependencyInfo = (*JSChunk)(nil)

// JSChunk_STRONG_CHUNK_NAME is the name of the artificial chunk containing
// all strong sources when there is no chunk spec. If there is a chunk spec,
// strong sources go in their respective chunks, and this chunk does not
// exist.
//
// Java: JSChunk.STRONG_CHUNK_NAME
const JSChunk_STRONG_CHUNK_NAME = "$strong$"

// JSChunk_WEAK_CHUNK_NAME is the name of the artificial chunk containing all
// weak sources. Regardless of the chunk spec, weak sources are moved into
// this chunk, which is made to depend on every other chunk. This is
// necessary so that removing weak sources (as an optimization) does not
// accidentally remove namespace declarations whose existence strong sources
// rely upon.
//
// Java: JSChunk.WEAK_CHUNK_NAME
const JSChunk_WEAK_CHUNK_NAME = "$weak$"

// Java: JSChunk.serialVersionUID
const jsChunk_serialVersionUID int64 = 1

// NewJSChunk creates an instance.
//
// name is a unique name for the chunk
//
// Java: JSChunk(String)
func NewJSChunk(name string) *JSChunk {
	return newJSChunkFromSlab(nil, name)
}

func newJSChunkFromSlab(slab *jutil.Slab[JSChunk], name string) *JSChunk {
	var c *JSChunk
	if slab != nil {
		c = slab.New()
	} else {
		c = &JSChunk{}
	}
	c.name = name
	c.inputs = jutil.NewLinkedHashMap[string, *CompilerInput]()
	// Depth and index will be set to their correct values by the JSChunkGraph into which they
	// are placed.
	c.depth = -1
	c.index = -1
	return c
}

// GetName gets the chunk name.
//
// Java: JSChunk.getName()
func (c *JSChunk) GetName() string {
	return c.name
}

// Java: JSChunk.getProvides()
func (c *JSChunk) GetProvides() []string {
	return []string{c.name}
}

// Java: JSChunk.getHasExternsAnnotation()
func (c *JSChunk) GetHasExternsAnnotation() bool {
	return false
}

// Java: JSChunk.getHasNoCompileAnnotation()
func (c *JSChunk) GetHasNoCompileAnnotation() bool {
	return false
}

// Java: JSChunk.getRequires()
func (c *JSChunk) GetRequires() []*DependencyInfo_Require {
	builder := make([]*DependencyInfo_Require, 0, len(c.deps))
	for _, chunk := range c.deps {
		builder = append(builder, DependencyInfo_Require_CompilerChunk(chunk.GetName()))
	}
	return builder
}

// Java: JSChunk.getTypeRequires()
func (c *JSChunk) GetTypeRequires() []string {
	// TODO(blickly): Actually allow weak chunk deps
	return []string{}
}

// Java: JSChunk.getPathRelativeToClosureBase()
func (c *JSChunk) GetPathRelativeToClosureBase() string {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: JSChunk.getLoadFlags()
func (c *JSChunk) GetLoadFlags() *jutil.LinkedHashMap[string, string] {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: JSChunk.isGoogModule()
func (c *JSChunk) IsGoogModule() bool {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: JSChunk.isEs6Module()
func (c *JSChunk) IsEs6Module() bool {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: DependencyInfo.getRequiredSymbols() (default method)
func (c *JSChunk) GetRequiredSymbols() []string {
	return DependencyInfo_GetRequiredSymbols(c)
}

// AddSourceFile adds a source file input to this chunk.
//
// Java: JSChunk.add(SourceFile)
func (c *JSChunk) AddSourceFile(file *SourceFile) {
	c.Add(NewCompilerInput(file))
}

// Add adds a source code input to this chunk.
//
// Java: JSChunk.add(CompilerInput)
func (c *JSChunk) Add(input *CompilerInput) {
	inputName := input.GetName()
	previous := c.inputs.Put(inputName, input)
	jutil.CheckStateMsg2(previous == nil, "%s already exist in chunk %s", inputName, c.GetName())
	c.cachedInputs = nil
	input.SetChunk(c)
}

// AddDependency adds a dependency on another chunk.
//
// Java: JSChunk.addDependency(JSChunk)
func (c *JSChunk) AddDependency(dep *JSChunk) {
	jutil.CheckNotNull(dep)
	jutil.CheckStateMsg1(dep != c, "Cannot add dependency on self (%s)", c)
	c.deps = append(c.deps, dep)
}

// Remove removes an input from this chunk.
//
// Java: JSChunk.remove(CompilerInput)
func (c *JSChunk) Remove(input *CompilerInput) {
	input.SetChunk(nil)
	c.inputs.Remove(input.GetName())
	c.cachedInputs = nil
}

// RemoveAll removes all of the inputs from this chunk.
//
// Java: JSChunk.removeAll()
func (c *JSChunk) RemoveAll() {
	for input := range c.inputs.Values() {
		input.SetChunk(nil)
	}
	c.inputs.Clear()
	c.cachedInputs = nil
}

// GetDependencies gets the list of chunks that this chunk depends on.
//
// # The result is a list that may be empty but not nil
//
// Java: JSChunk.getDependencies()
func (c *JSChunk) GetDependencies() []*JSChunk {
	return append([]*JSChunk{}, c.deps...)
}

// GetSortedDependencyNames gets the names of the chunks that this chunk
// depends on, sorted alphabetically.
//
// Java: JSChunk.getSortedDependencyNames()
func (c *JSChunk) GetSortedDependencyNames() []string {
	names := []string{}
	for _, chunk := range c.GetDependencies() {
		names = append(names, chunk.GetName())
	}
	slices.SortStableFunc(names, jutil.StringCompareTo)
	return names
}

// GetAllDependencies returns the transitive closure of dependencies
// starting from the dependencies of this chunk.
//
// Java: JSChunk.getAllDependencies()
func (c *JSChunk) GetAllDependencies() *jutil.LinkedHashSet[*JSChunk] {
	allDeps := jutil.NewLinkedHashSet[*JSChunk]()
	allDeps.AddAll(c.deps...)
	// Go-only: language-adapter — the ArrayDeque used as a stack (push and
	// pop at its head) is a slice whose top is its last element, so it
	// starts as deps reversed.
	stack := make([]*JSChunk, 0, len(c.deps))
	for i := len(c.deps) - 1; i >= 0; i-- {
		stack = append(stack, c.deps[i])
	}
	for len(stack) != 0 {
		chunk := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		chunkDeps := chunk.deps
		for _, dep := range chunkDeps {
			if allDeps.Add(dep) {
				stack = append(stack, dep)
			}
		}
	}
	return allDeps
}

// GetThisAndAllDependencies returns this chunk and all of its dependencies
// in one list.
//
// Java: JSChunk.getThisAndAllDependencies()
func (c *JSChunk) GetThisAndAllDependencies() *jutil.LinkedHashSet[*JSChunk] {
	deps := c.GetAllDependencies()
	deps.Add(c)
	return deps
}

// GetInputCount returns the number of source code inputs.
//
// Java: JSChunk.getInputCount()
func (c *JSChunk) GetInputCount() int {
	return c.inputs.Size()
}

var emptyCompilerInputs = []*CompilerInput{}

// GetFirst returns the first source code input.
//
// Java: JSChunk.getFirst()
func (c *JSChunk) GetFirst() *CompilerInput {
	if input, ok := c.inputs.ValueAt(0); ok {
		return input
	}
	panic(jutil.NewNoSuchElementException())
}

// GetInputs gets this chunk's list of source code inputs.
//
// # The result is a list that may be empty but not nil
//
// Java: JSChunk.getInputs()
func (c *JSChunk) GetInputs() []*CompilerInput {
	if c.cachedInputs != nil {
		return c.cachedInputs[:len(c.cachedInputs):len(c.cachedInputs)]
	}
	n := c.inputs.Size()
	if n == 0 {
		c.cachedInputs = emptyCompilerInputs
		return emptyCompilerInputs[:0:0]
	}
	out := make([]*CompilerInput, 0, n)
	slot := 0
	for {
		input, nextSlot, ok := c.inputs.NextSlotValue(slot)
		if !ok {
			break
		}
		out = append(out, input)
		slot = nextSlot
	}
	c.cachedInputs = out
	return out[:len(out):len(out)]
}

// GetInputsIterable is a live view of the inputs.
//
// Java: JSChunk.getInputsIterable()
func (c *JSChunk) GetInputsIterable() iter.Seq[*CompilerInput] {
	return c.inputs.Values()
}

// GetByName returns the input with the given name or nil if none.
//
// Java: JSChunk.getByName(String)
func (c *JSChunk) GetByName(name string) *CompilerInput {
	return c.inputs.Get(name)
}

// RemoveByName removes any input with the given name. Returns whether any
// were removed.
//
// Java: JSChunk.removeByName(String)
func (c *JSChunk) RemoveByName(name string) bool {
	value := c.inputs.Remove(name)
	if value != nil {
		c.cachedInputs = nil
		return true
	}
	return false
}

// IsSynthetic returns whether this chunk is synthetic (i.e. one of the
// special strong or weak chunks created by the compiler.
//
// Java: JSChunk.isSynthetic()
func (c *JSChunk) IsSynthetic() bool {
	return c.name == JSChunk_STRONG_CHUNK_NAME || c.name == JSChunk_WEAK_CHUNK_NAME
}

// Java: JSChunk.isWeak()
func (c *JSChunk) IsWeak() bool {
	return c.name == JSChunk_WEAK_CHUNK_NAME
}

// String returns the chunk name (primarily for debugging).
//
// Java: JSChunk.toString()
func (c *JSChunk) String() string {
	return c.name
}

// SortInputsByDeps puts the JS files into a topologically sorted order by
// their dependencies.
//
// Java: JSChunk.sortInputsByDeps(AbstractCompiler)
func (c *JSChunk) SortInputsByDeps(compiler AbstractCompiler) {
	// Set the compiler, so that we can parse requires/provides and report
	// errors properly.
	for input := range c.inputs.Values() {
		input.SetCompiler(compiler)
	}

	// Sort the JSChunk in this order.
	sortedList := NewSortedDependencies[*CompilerInput](c.GetInputs()).GetSortedList()
	c.inputs.Clear()
	for _, input := range sortedList {
		c.inputs.Put(input.GetName(), input)
	}
	c.cachedInputs = sortedList
}

// SetDepth sets the depth.
//
// Java: JSChunk.setDepth(int)
func (c *JSChunk) SetDepth(dep int) {
	jutil.CheckArgumentMsg1(dep >= 0, "invalid depth: %s", dep)
	c.depth = dep
}

// GetDepth returns the depth.
//
// Java: JSChunk.getDepth()
func (c *JSChunk) GetDepth() int {
	return c.depth
}

// Java: JSChunk.setIndex(int)
func (c *JSChunk) SetIndex(index int) {
	jutil.CheckArgumentMsg1(index >= 0, "Invalid chunk index: %s", index)
	c.index = index
}

// Java: JSChunk.getIndex()
func (c *JSChunk) GetIndex() int {
	return c.index
}

// ReadObject is the deserialization hook: after defaultReadObject, the
// transient input map is a new empty one (CompilerInputs must be explicitly
// added to the JSChunk again after deserialization).
//
// Go-only: language-adapter — there is no ObjectInputStream; see the type
// comment.
//
// Java: JSChunk.readObject(ObjectInputStream)
func (c *JSChunk) ReadObject() {
	c.inputs = jutil.NewLinkedHashMap[string, *CompilerInput]()
	c.cachedInputs = nil
}
