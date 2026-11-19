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
// Port of com/google/debugging/sourcemap/SourceMapObject.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.debugging.sourcemap.SourceMapObject
type SourceMapObject struct {
	sourceRoot string `gson:"sourceRoot"`
	file string `gson:"file"`
	mappings string `gson:"mappings"`
	sources []string `gson:"sources"`
	sourcesContent []string `gson:"sourcesContent"`
	names []string `gson:"names"`
	sections j2gort.Ext_java_util_List `gson:"sections"`
	extensions j2gort.Ext_java_util_Map `gson:"extensions"`
	version int32 `gson:"version"`
	lineCount int32 `gson:"lineCount"`
	idHash_ atomic.Int32
}

func (this_ *SourceMapObject) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapObject) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapObject) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.SourceMapObject", "SourceMapObject" }
func (this_ *SourceMapObject) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapObject) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapObject", "com.google.debugging.sourcemap.SourceMapObject", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapObject](v_)
	})
}

func (this_ *SourceMapObject) initFields_SourceMapObject() {
	_ = this_
	this_.sourceRoot = j2gort.NullString
	this_.file = j2gort.NullString
	this_.mappings = j2gort.NullString
}

// Java constructor: SourceMapObject(int version, int lineCount, java.lang.String sourceRoot, java.lang.String file, java.lang.String mappings, java.lang.String[] sources, java.lang.String[] sourcesContent, java.lang.String[] names, java.util.List<com.google.debugging.sourcemap.SourceMapSection> sections, java.util.Map<java.lang.String,java.lang.Object> extensions)
func NewSourceMapObject(version int32, lineCount int32, sourceRoot string, file string, mappings string, sources []string, sourcesContent []string, names []string, sections j2gort.Ext_java_util_List, extensions j2gort.Ext_java_util_Map) *SourceMapObject {
	this_ := &SourceMapObject{}
	InitSourceMapObject(this_, version, lineCount, sourceRoot, file, mappings, sources, sourcesContent, names, sections, extensions)
	return this_
}

func InitSourceMapObject(this_ *SourceMapObject, version int32, lineCount int32, sourceRoot string, file string, mappings string, sources []string, sourcesContent []string, names []string, sections j2gort.Ext_java_util_List, extensions j2gort.Ext_java_util_Map) {
	_ = version
	_ = lineCount
	_ = sourceRoot
	_ = file
	_ = mappings
	_ = sources
	_ = sourcesContent
	_ = names
	_ = sections
	_ = extensions
	this_.initFields_SourceMapObject()
	this_.version = version
	this_.lineCount = lineCount
	this_.sourceRoot = sourceRoot
	this_.file = file
	this_.mappings = mappings
	this_.sources = sources
	this_.sourcesContent = sourcesContent
	this_.names = names
	this_.sections = sections
	this_.extensions = extensions
}

// Java method: SourceMapObject.getVersion()
func (this_ *SourceMapObject) GetVersion() int32 {
	return this_.version
}

// Java method: SourceMapObject.getLineCount()
func (this_ *SourceMapObject) GetLineCount() int32 {
	return this_.lineCount
}

// Java method: SourceMapObject.getSourceRoot()
func (this_ *SourceMapObject) GetSourceRoot() string {
	return this_.sourceRoot
}

// Java method: SourceMapObject.getFile()
func (this_ *SourceMapObject) GetFile() string {
	return this_.file
}

// Java method: SourceMapObject.getMappings()
func (this_ *SourceMapObject) GetMappings() string {
	return this_.mappings
}

// Java method: SourceMapObject.getSources()
func (this_ *SourceMapObject) GetSources() []string {
	return this_.sources
}

// Java method: SourceMapObject.getSourcesContent()
func (this_ *SourceMapObject) GetSourcesContent() []string {
	return this_.sourcesContent
}

// Java method: SourceMapObject.getNames()
func (this_ *SourceMapObject) GetNames() []string {
	return this_.names
}

// Java method: SourceMapObject.getSections()
func (this_ *SourceMapObject) GetSections() j2gort.Ext_java_util_List {
	return this_.sections
}

// Java method: SourceMapObject.getExtensions()
func (this_ *SourceMapObject) GetExtensions() j2gort.Ext_java_util_Map {
	return this_.extensions
}

// Java static method: SourceMapObject.builder()
func SourceMapObject_S_Builder() *SourceMapObject_Builder {
	return NewSourceMapObject_Builder()
}

// Java class: com.google.debugging.sourcemap.SourceMapObject.Builder
type SourceMapObject_Builder struct {
	sourceRoot string `gson:"sourceRoot"`
	file string `gson:"file"`
	mappings string `gson:"mappings"`
	sources []string `gson:"sources"`
	sourcesContent []string `gson:"sourcesContent"`
	names []string `gson:"names"`
	sections j2gort.Ext_java_util_List `gson:"sections"`
	extensions j2gort.Ext_java_util_Map `gson:"extensions"`
	version int32 `gson:"version"`
	lineCount int32 `gson:"lineCount"`
	idHash_ atomic.Int32
}

func (this_ *SourceMapObject_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapObject_Builder) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapObject_Builder) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.SourceMapObject$Builder", "Builder" }
func (this_ *SourceMapObject_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceMapObject_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapObject.Builder", "com.google.debugging.sourcemap.SourceMapObject$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapObject_Builder](v_)
	})
}

func (this_ *SourceMapObject_Builder) initFields_SourceMapObject_Builder() {
	_ = this_
	this_.sourceRoot = j2gort.NullString
	this_.file = j2gort.NullString
	this_.mappings = j2gort.NullString
}

// Java constructor: Builder()
func NewSourceMapObject_Builder() *SourceMapObject_Builder {
	this_ := &SourceMapObject_Builder{}
	InitSourceMapObject_Builder(this_)
	return this_
}

func InitSourceMapObject_Builder(this_ *SourceMapObject_Builder) {
	this_.initFields_SourceMapObject_Builder()
}

// Java method: Builder.setVersion(int version)
func (this_ *SourceMapObject_Builder) SetVersion(version int32) *SourceMapObject_Builder {
	_ = version
	this_.version = version
	return this_
}

// Java method: Builder.setLineCount(int lineCount)
func (this_ *SourceMapObject_Builder) SetLineCount(lineCount int32) *SourceMapObject_Builder {
	_ = lineCount
	this_.lineCount = lineCount
	return this_
}

// Java method: Builder.setSourceRoot(java.lang.String sourceRoot)
func (this_ *SourceMapObject_Builder) SetSourceRoot(sourceRoot string) *SourceMapObject_Builder {
	_ = sourceRoot
	this_.sourceRoot = sourceRoot
	return this_
}

// Java method: Builder.setFile(java.lang.String file)
func (this_ *SourceMapObject_Builder) SetFile(file string) *SourceMapObject_Builder {
	_ = file
	this_.file = file
	return this_
}

// Java method: Builder.setMappings(java.lang.String mappings)
func (this_ *SourceMapObject_Builder) SetMappings(mappings string) *SourceMapObject_Builder {
	_ = mappings
	this_.mappings = mappings
	return this_
}

// Java method: Builder.setSources(java.lang.String[] sources)
func (this_ *SourceMapObject_Builder) SetSources(sources []string) *SourceMapObject_Builder {
	_ = sources
	this_.sources = sources
	return this_
}

// Java method: Builder.setSourcesContent(java.lang.String[] sourcesContent)
func (this_ *SourceMapObject_Builder) SetSourcesContent(sourcesContent []string) *SourceMapObject_Builder {
	_ = sourcesContent
	this_.sourcesContent = sourcesContent
	return this_
}

// Java method: Builder.setNames(java.lang.String[] names)
func (this_ *SourceMapObject_Builder) SetNames(names []string) *SourceMapObject_Builder {
	_ = names
	this_.names = names
	return this_
}

// Java method: Builder.setSections(java.util.List<com.google.debugging.sourcemap.SourceMapSection> sections)
func (this_ *SourceMapObject_Builder) SetSections(sections j2gort.Ext_java_util_List) *SourceMapObject_Builder {
	_ = sections
	this_.sections = sections
	return this_
}

// Java method: Builder.setExtensions(java.util.Map<java.lang.String,java.lang.Object> extensions)
func (this_ *SourceMapObject_Builder) SetExtensions(extensions j2gort.Ext_java_util_Map) *SourceMapObject_Builder {
	_ = extensions
	this_.extensions = extensions
	return this_
}

// Java method: Builder.build()
func (this_ *SourceMapObject_Builder) Build() *SourceMapObject {
	return NewSourceMapObject(this_.version, this_.lineCount, this_.sourceRoot, this_.file, this_.mappings, this_.sources, this_.sourcesContent, this_.names, this_.sections, this_.extensions)
}

