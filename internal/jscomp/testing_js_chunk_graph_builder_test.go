// Port of com/google/javascript/jscomp/testing/JSChunkGraphBuilder.java @ closure-compiler v20260415

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
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSChunkGraphBuilder is a utility to create various input JSChunk graphs
type JSChunkGraphBuilder struct {
	graphType jsChunkGraphBuilder_GraphType
	chunks    []string
	// correspondances between chunk indices in `chunks` and specified names. if no name is specified
	// defaults to "m{i}"
	chunkNames     *jutil.LinkedHashMap[int, string]
	filenameFormat string
}

// Java: JSChunkGraphBuilder.GraphType
type jsChunkGraphBuilder_GraphType int32

const (
	// each chunk depends on the chunk immediately before
	jsChunkGraphBuilder_GraphType_CHAIN jsChunkGraphBuilder_GraphType = iota
	// each chunk depends on the fist chunk
	jsChunkGraphBuilder_GraphType_STAR
	// no dependencies between chunks
	jsChunkGraphBuilder_GraphType_DISJOINT
	// chunk 2 depends on chunk 1, and all others depend on chunk 2
	jsChunkGraphBuilder_GraphType_BUSH
	// binary tree
	jsChunkGraphBuilder_GraphType_TREE
)

// Java: JSChunkGraphBuilder.GraphType.addDependencyEdges(JSChunk[])
func (g jsChunkGraphBuilder_GraphType) addDependencyEdges(chunks []*JSChunk) {
	switch g {
	case jsChunkGraphBuilder_GraphType_CHAIN:
		for i := 1; i < len(chunks); i++ {
			chunks[i].AddDependency(chunks[i-1])
		}
	case jsChunkGraphBuilder_GraphType_STAR:
		for i := 1; i < len(chunks); i++ {
			chunks[i].AddDependency(chunks[0])
		}
	case jsChunkGraphBuilder_GraphType_BUSH:
		jutil.CheckStateMsg(len(chunks) > 2, "BUSHes need at least three graph nodes")
		for i := 1; i < len(chunks); i++ {
			if i == 1 {
				chunks[i].AddDependency(chunks[0])
			} else {
				chunks[i].AddDependency(chunks[1])
			}
		}
	case jsChunkGraphBuilder_GraphType_TREE:
		for i := 1; i < len(chunks); i++ {
			chunks[i].AddDependency(chunks[(i-1)/2])
		}
	default:
	}
}

// Java: JSChunkGraphBuilder(GraphType)
func newJSChunkGraphBuilder(graphType jsChunkGraphBuilder_GraphType) *JSChunkGraphBuilder {
	return &JSChunkGraphBuilder{
		graphType:      graphType,
		chunks:         []string{},
		chunkNames:     jutil.NewLinkedHashMap[int, string](),
		filenameFormat: "i%s.js",
	}
}

// JSChunkGraphBuilder_ForChain generates a list of chunks from a list of
// inputs, such that each chunk depends on the chunk before it.
//
// Java: JSChunkGraphBuilder.forChain()
func JSChunkGraphBuilder_ForChain() *JSChunkGraphBuilder {
	return newJSChunkGraphBuilder(jsChunkGraphBuilder_GraphType_CHAIN)
}

// JSChunkGraphBuilder_ForStar generates a list of chunks from a list of
// inputs, such that each chunk depends on the first chunk.
//
// Java: JSChunkGraphBuilder.forStar()
func JSChunkGraphBuilder_ForStar() *JSChunkGraphBuilder {
	return newJSChunkGraphBuilder(jsChunkGraphBuilder_GraphType_STAR)
}

// JSChunkGraphBuilder_ForBush generates a list of chunks from a list of
// inputs, such that chunks form a bush formation. In a bush formation, chunk
// 2 depends on chunk 1, and all other chunks depend on chunk 2.
//
// Java: JSChunkGraphBuilder.forBush()
func JSChunkGraphBuilder_ForBush() *JSChunkGraphBuilder {
	return newJSChunkGraphBuilder(jsChunkGraphBuilder_GraphType_BUSH)
}

// JSChunkGraphBuilder_ForTree generates a list of chunks from a list of
// inputs, such that chunks form a tree formation. In a tree formation, chunk
// N depends on chunk `floor(N/2)`, So the chunks form a balanced binary
// tree.
//
// Java: JSChunkGraphBuilder.forTree()
func JSChunkGraphBuilder_ForTree() *JSChunkGraphBuilder {
	return newJSChunkGraphBuilder(jsChunkGraphBuilder_GraphType_TREE)
}

// JSChunkGraphBuilder_ForUnordered generates a list of chunks from a list
// of inputs with no dependency edges between the chunks
//
// Java: JSChunkGraphBuilder.forUnordered()
func JSChunkGraphBuilder_ForUnordered() *JSChunkGraphBuilder {
	return newJSChunkGraphBuilder(jsChunkGraphBuilder_GraphType_DISJOINT)
}

// Java: JSChunkGraphBuilder.addChunk(String)
func (b *JSChunkGraphBuilder) AddChunk(input string) *JSChunkGraphBuilder {
	b.chunks = append(b.chunks, input)
	return b
}

// Java: JSChunkGraphBuilder.addChunkWithName(String, String)
func (b *JSChunkGraphBuilder) AddChunkWithName(input string, chunkName string) *JSChunkGraphBuilder {
	b.chunkNames.Put(len(b.chunks), chunkName)
	b.chunks = append(b.chunks, input)
	return b
}

// AddChunks takes Java's List<String> and String[] overloads alike.
//
// Java: JSChunkGraphBuilder.addChunks(List<String>)
// Java: JSChunkGraphBuilder.addChunks(String[])
func (b *JSChunkGraphBuilder) AddChunks(inputs []string) *JSChunkGraphBuilder {
	b.chunks = append(b.chunks, inputs...)
	return b
}

// SetFilenameFormat sets a format string for the generated file names for
// each chunk.
//
// String must contain one substitution, the index of the file. e.g.
// "i%s.js" will produce i0.js, i1.js, etc.
//
// Java: JSChunkGraphBuilder.setFilenameFormat(String)
func (b *JSChunkGraphBuilder) SetFilenameFormat(filenameFormat string) *JSChunkGraphBuilder {
	b.filenameFormat = filenameFormat
	return b
}

// Java: JSChunkGraphBuilder.build()
func (b *JSChunkGraphBuilder) Build() []*JSChunk {
	chunks := make([]*JSChunk, len(b.chunks))
	for i := 0; i < len(b.chunks); i++ {
		chunkName := b.chunkNames.GetOrDefault(i, "m"+fmt.Sprint(i))
		chunks[i] = NewJSChunk(chunkName)
		chunk := chunks[i]

		chunk.AddSourceFile(
			SourceFile_FromCode(jutil.LenientFormat(b.filenameFormat, i), b.chunks[i]))
	}
	b.graphType.addDependencyEdges(chunks)
	return chunks
}
