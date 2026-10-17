// Port of com/google/debugging/sourcemap/SourceMapObject.java @ closure-compiler v20260415

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

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapObject wraps a JsonObject to provide a V3 source map.
type SourceMapObject struct {
	version        int
	lineCount      int
	sourceRoot     *string
	file           *string
	mappings       *string
	sources        []*string
	sourcesContent []*string
	names          []*string
	sections       []*SourceMapSection
	extensions     *jutil.LinkedHashMap[string, any]
}

// NewSourceMapObject constructs a SourceMapObject (private in Java; use
// the Builder).
//
// Java: SourceMapObject(int, int, String, String, String, String[], String[], String[], List, Map)
func NewSourceMapObject(
	version int,
	lineCount int,
	sourceRoot *string,
	file *string,
	mappings *string,
	sources []*string,
	sourcesContent []*string,
	names []*string,
	sections []*SourceMapSection,
	extensions *jutil.LinkedHashMap[string, any]) *SourceMapObject {
	o := &SourceMapObject{}
	o.version = version
	o.lineCount = lineCount
	o.sourceRoot = sourceRoot
	o.file = file
	o.mappings = mappings
	o.sources = sources
	o.sourcesContent = sourcesContent
	o.names = names
	o.sections = sections
	o.extensions = extensions
	return o
}

// Java: SourceMapObject.getVersion()
func (o *SourceMapObject) GetVersion() int {
	return o.version
}

// Java: SourceMapObject.getLineCount()
func (o *SourceMapObject) GetLineCount() int {
	return o.lineCount
}

// Java: SourceMapObject.getSourceRoot()
func (o *SourceMapObject) GetSourceRoot() *string {
	return o.sourceRoot
}

// Java: SourceMapObject.getFile()
func (o *SourceMapObject) GetFile() *string {
	return o.file
}

// Java: SourceMapObject.getMappings()
func (o *SourceMapObject) GetMappings() *string {
	return o.mappings
}

// GetSources returns the sources (nil for Java null; elements may be nil).
//
// Java: SourceMapObject.getSources()
func (o *SourceMapObject) GetSources() []*string {
	return o.sources
}

// Java: SourceMapObject.getSourcesContent()
func (o *SourceMapObject) GetSourcesContent() []*string {
	return o.sourcesContent
}

// Java: SourceMapObject.getNames()
func (o *SourceMapObject) GetNames() []*string {
	return o.names
}

// GetSections returns the index map sections (nil for Java null).
//
// Java: SourceMapObject.getSections()
func (o *SourceMapObject) GetSections() []*SourceMapSection {
	return o.sections
}

// Java: SourceMapObject.getExtensions()
func (o *SourceMapObject) GetExtensions() *jutil.LinkedHashMap[string, any] {
	return o.extensions
}

// SourceMapObject_Builder_ returns a new builder.
//
// Go-only: language-adapter — named with a trailing underscore because the
// nested type SourceMapObject_Builder has the same Go name.
//
// Java: SourceMapObject.builder()
func SourceMapObject_Builder_() *SourceMapObject_Builder {
	return NewSourceMapObject_Builder()
}

// SourceMapObject_Builder is SourceMapObject.Builder.
type SourceMapObject_Builder struct {
	version        int
	lineCount      int
	sourceRoot     *string
	file           *string
	mappings       *string
	sources        []*string
	sourcesContent []*string
	names          []*string
	sections       []*SourceMapSection
	extensions     *jutil.LinkedHashMap[string, any]
}

// NewSourceMapObject_Builder constructs an empty builder.
//
// Java: SourceMapObject.Builder()
func NewSourceMapObject_Builder() *SourceMapObject_Builder {
	return &SourceMapObject_Builder{}
}

// Java: SourceMapObject.Builder.setVersion(int)
func (b *SourceMapObject_Builder) SetVersion(version int) *SourceMapObject_Builder {
	b.version = version
	return b
}

// Java: SourceMapObject.Builder.setLineCount(int)
func (b *SourceMapObject_Builder) SetLineCount(lineCount int) *SourceMapObject_Builder {
	b.lineCount = lineCount
	return b
}

// Java: SourceMapObject.Builder.setSourceRoot(String)
func (b *SourceMapObject_Builder) SetSourceRoot(sourceRoot *string) *SourceMapObject_Builder {
	b.sourceRoot = sourceRoot
	return b
}

// Java: SourceMapObject.Builder.setFile(String)
func (b *SourceMapObject_Builder) SetFile(file *string) *SourceMapObject_Builder {
	b.file = file
	return b
}

// Java: SourceMapObject.Builder.setMappings(String)
func (b *SourceMapObject_Builder) SetMappings(mappings *string) *SourceMapObject_Builder {
	b.mappings = mappings
	return b
}

// Java: SourceMapObject.Builder.setSources(String[])
func (b *SourceMapObject_Builder) SetSources(sources []*string) *SourceMapObject_Builder {
	b.sources = sources
	return b
}

// Java: SourceMapObject.Builder.setSourcesContent(String[])
func (b *SourceMapObject_Builder) SetSourcesContent(sourcesContent []*string) *SourceMapObject_Builder {
	b.sourcesContent = sourcesContent
	return b
}

// Java: SourceMapObject.Builder.setNames(String[])
func (b *SourceMapObject_Builder) SetNames(names []*string) *SourceMapObject_Builder {
	b.names = names
	return b
}

// Java: SourceMapObject.Builder.setSections(List)
func (b *SourceMapObject_Builder) SetSections(sections []*SourceMapSection) *SourceMapObject_Builder {
	b.sections = sections
	return b
}

// Java: SourceMapObject.Builder.setExtensions(Map)
func (b *SourceMapObject_Builder) SetExtensions(extensions *jutil.LinkedHashMap[string, any]) *SourceMapObject_Builder {
	b.extensions = extensions
	return b
}

// Java: SourceMapObject.Builder.build()
func (b *SourceMapObject_Builder) Build() *SourceMapObject {
	return NewSourceMapObject(
		b.version,
		b.lineCount,
		b.sourceRoot,
		b.file,
		b.mappings,
		b.sources,
		b.sourcesContent,
		b.names,
		b.sections,
		b.extensions)
}
