// Port of com/google/debugging/sourcemap/SourceMapObject.java @ closure-compiler v20261021

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
	sourceRoot     optionalString
	file           optionalString
	mappings       optionalString
	sources        []optionalString
	sourcesContent []optionalString
	names          []optionalString
	sections       []*SourceMapSection
	extensions     []SourceMapExtension
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
	extensions []SourceMapExtension) *SourceMapObject {
	return &SourceMapObject{
		version:        version,
		lineCount:      lineCount,
		sourceRoot:     makeOptionalString(sourceRoot),
		file:           makeOptionalString(file),
		mappings:       makeOptionalString(mappings),
		sources:        optionalStringsFromPtrs(sources),
		sourcesContent: optionalStringsFromPtrs(sourcesContent),
		names:          optionalStringsFromPtrs(names),
		sections:       sections,
		extensions:     extensions,
	}
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
	return o.sourceRoot.ptr()
}

// Java: SourceMapObject.getFile()
func (o *SourceMapObject) GetFile() *string {
	return o.file.ptr()
}

// Java: SourceMapObject.getMappings()
func (o *SourceMapObject) GetMappings() *string {
	return o.mappings.ptr()
}

// GetSources returns the sources (nil for Java null; elements may be nil).
//
// Java: SourceMapObject.getSources()
func (o *SourceMapObject) GetSources() []*string {
	return optionalStringsToPtrs(o.sources)
}

// Java: SourceMapObject.getSourcesContent()
func (o *SourceMapObject) GetSourcesContent() []*string {
	return optionalStringsToPtrs(o.sourcesContent)
}

// Java: SourceMapObject.getNames()
func (o *SourceMapObject) GetNames() []*string {
	return optionalStringsToPtrs(o.names)
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
