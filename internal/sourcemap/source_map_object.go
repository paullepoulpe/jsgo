// Port of com/google/debugging/sourcemap/SourceMapObject.java @ closure-compiler v20261028

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

import "slices"

// SourceMapExtension is a custom "x_..." extension property on a V3 source map.
type SourceMapExtension struct {
	Key   string
	Value any
}

func putExtension(exts []SourceMapExtension, key string, val any) []SourceMapExtension {
	for i := range exts {
		if exts[i].Key == key {
			exts[i].Value = val
			return exts
		}
	}
	return append(exts, SourceMapExtension{Key: key, Value: val})
}

func lookupExtension(exts []SourceMapExtension, key string) (any, bool) {
	for i := range exts {
		if exts[i].Key == key {
			return exts[i].Value, true
		}
	}
	return nil, false
}

func removeExtension(exts []SourceMapExtension, key string) []SourceMapExtension {
	for i := range exts {
		if exts[i].Key == key {
			return slices.Delete(exts, i, i+1)
		}
	}
	return exts
}

// SourceMapObject wraps a JsonObject to provide a V3 source map.
type SourceMapObject struct {
	version        int
	lineCount      int
	sourceRoot     OptionalString
	file           OptionalString
	mappings       OptionalString
	sources        []OptionalString
	sourcesContent []OptionalString
	names          []OptionalString
	sections       []*SourceMapSection
	extensions     []SourceMapExtension
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
func (o *SourceMapObject) GetSourceRoot() (string, bool) {
	return o.sourceRoot.Get()
}

// Java: SourceMapObject.getFile()
func (o *SourceMapObject) GetFile() (string, bool) {
	return o.file.Get()
}

// Java: SourceMapObject.getMappings()
func (o *SourceMapObject) GetMappings() (string, bool) {
	return o.mappings.Get()
}

// GetSources returns the sources (nil for Java null; elements may be absent).
//
// Java: SourceMapObject.getSources()
func (o *SourceMapObject) GetSources() []OptionalString {
	return o.sources
}

// Java: SourceMapObject.getSourcesContent()
func (o *SourceMapObject) GetSourcesContent() []OptionalString {
	return o.sourcesContent
}

// Java: SourceMapObject.getNames()
func (o *SourceMapObject) GetNames() []OptionalString {
	return o.names
}

// GetSections returns the index map sections (nil for Java null).
//
// Java: SourceMapObject.getSections()
func (o *SourceMapObject) GetSections() []*SourceMapSection {
	return o.sections
}

// Java: SourceMapObject.getExtensions()
func (o *SourceMapObject) GetExtensions() []SourceMapExtension {
	return o.extensions
}
