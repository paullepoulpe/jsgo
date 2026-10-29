// Port of com/google/javascript/jscomp/bundle/Source.java @ closure-compiler v20261021
// (with the AutoValue-generated bundle/AutoValue_Source.java)

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

package bundle

import (
	"slices"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type sourceLoadFlag struct {
	key   string
	value string
}

type lazyString struct {
	val string
	fn  func() string
}

func (l lazyString) get() string {
	if l.fn != nil {
		return l.fn()
	}
	return l.val
}

// Source is an abstract representation of a source file.
//
// Go-only: language-adapter — Source is an @AutoValue: the abstract class
// and the generated AutoValue_Source (fields, equals, hashCode, toString,
// Builder) are flattened into plain value fields here (PORTING.md §4).
type Source struct {
	path             *jutil.Path
	sourceMap        string
	sourceUrl        string
	sourceMappingUrl string
	runtimes         []string
	loadFlags        []sourceLoadFlag
	estimatedSize    int
	code             lazyString
	originalCode     lazyString
}

// Path is the path of this source. This may refer to a path on disk or a
// path on the HTTP server.
//
// Java: Source.path()
func (s *Source) Path() *jutil.Path { return s.path }

// SourceMap is the text of any source map applicable to this file.
//
// Java: Source.sourceMap()
func (s *Source) SourceMap() string { return s.sourceMap }

// SourceUrl is the source URL associated with this file.
//
// Java: Source.sourceUrl()
func (s *Source) SourceUrl() string { return s.sourceUrl }

// SourceMappingUrl is the URL for a source map associated with this file.
//
// Java: Source.sourceMappingUrl()
func (s *Source) SourceMappingUrl() string { return s.sourceMappingUrl }

// Runtimes are any runtime libraries necessary for this source. Any
// transformation that adds a runtime library to any sources must be
// responsible to never add the same library as a substring to a different
// source (so that the "no duplicates" invariant of Set will work
// correctly).
//
// Java: Source.runtimes()
func (s *Source) Runtimes() []string {
	if len(s.runtimes) == 0 {
		return []string{}
	}
	return slices.Clone(s.runtimes)
}

// LoadFlags are the load flags, specifying module type and language level.
//
// Java: Source.loadFlags()
func (s *Source) LoadFlags() map[string]string {
	m := make(map[string]string, len(s.loadFlags))
	for _, f := range s.loadFlags {
		m[f.key] = f.value
	}
	return m
}

// EstimatedSize is a best estimate of the size of this source (in case the
// source itself is not yet loaded.
//
// Java: Source.estimatedSize()
func (s *Source) EstimatedSize() int { return s.estimatedSize }

// Code is the actual code in this source file.
//
// Java: Source.code()
func (s *Source) Code() string { return s.code.get() }

// OriginalCode is the untransformed code from the original source file.
//
// Java: Source.originalCode()
func (s *Source) OriginalCode() string { return s.originalCode.get() }

// ToBuilder copies the data from this source to a new builder.
//
// Java: Source.toBuilder()
func (s *Source) ToBuilder() *Source_Builder {
	return &Source_Builder{
		path:             s.path,
		sourceMap:        s.sourceMap,
		sourceUrl:        s.sourceUrl,
		sourceMappingUrl: s.sourceMappingUrl,
		runtimes:         s.runtimes,
		loadFlags:        s.loadFlags,
		estimatedSize:    s.estimatedSize,
		code:             s.code,
		hasCode:          true,
		originalCode:     s.originalCode,
		hasOriginalCode:  true,
	}
}

// Source_Builder_ makes a new empty builder (Source.builder()).
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name Source_Builder for the nested Builder type.
func Source_Builder_() *Source_Builder {
	return NewSource_Builder()
}

var source_DEV_NULL = jutil.Path_Of("/dev/null")

// Java: AutoValue_Source.toString()
func (s *Source) String() string {
	var sb strings.Builder
	sb.WriteString("Source{path=")
	sb.WriteString(s.path.String())
	sb.WriteString(", sourceMap=")
	sb.WriteString(s.sourceMap)
	sb.WriteString(", sourceUrl=")
	sb.WriteString(s.sourceUrl)
	sb.WriteString(", sourceMappingUrl=")
	sb.WriteString(s.sourceMappingUrl)
	sb.WriteString(", runtimes=[")
	for i, r := range s.runtimes {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(r)
	}
	sb.WriteString("], loadFlags={")
	for i, f := range s.loadFlags {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(f.key)
		sb.WriteByte('=')
		sb.WriteString(f.value)
	}
	sb.WriteString("}, estimatedSize=")
	sb.WriteString(jutil.ToJavaString(s.estimatedSize))
	sb.WriteString(", code=")
	sb.WriteString(s.code.get())
	sb.WriteString(", originalCode=")
	sb.WriteString(s.originalCode.get())
	sb.WriteByte('}')
	return sb.String()
}

// Java: AutoValue_Source.equals(Object)
func (s *Source) Equals(o any) bool {
	that, ok := o.(*Source)
	if !ok || that == nil {
		return false
	}
	if that == s {
		return true
	}
	return s.estimatedSize == that.estimatedSize &&
		s.path.Equals(that.path) &&
		s.sourceMap == that.sourceMap &&
		s.sourceUrl == that.sourceUrl &&
		s.sourceMappingUrl == that.sourceMappingUrl &&
		slices.Equal(s.runtimes, that.runtimes) &&
		sourceLoadFlagsEqual(s.loadFlags, that.loadFlags) &&
		s.code.get() == that.code.get() &&
		s.originalCode.get() == that.originalCode.get()
}

func sourceLoadFlagsEqual(a, b []sourceLoadFlag) bool {
	if len(a) != len(b) {
		return false
	}
	for _, fa := range a {
		found := false
		for _, fb := range b {
			if fa.key == fb.key {
				found = fa.value == fb.value
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func sourceRuntimesHashCode(runtimes []string) int32 {
	var h int32
	for _, r := range runtimes {
		h += jutil.StringHashCode(r)
	}
	return h
}

func sourceLoadFlagsHashCode(flags []sourceLoadFlag) int32 {
	var h int32
	for _, f := range flags {
		h += jutil.StringHashCode(f.key) ^ jutil.StringHashCode(f.value)
	}
	return h
}

// Java: AutoValue_Source.hashCode()
func (s *Source) HashCode() int32 {
	h := int32(1)
	h *= 1000003
	h ^= s.path.HashCode()
	h *= 1000003
	h ^= jutil.StringHashCode(s.sourceMap)
	h *= 1000003
	h ^= jutil.StringHashCode(s.sourceUrl)
	h *= 1000003
	h ^= jutil.StringHashCode(s.sourceMappingUrl)
	h *= 1000003
	h ^= sourceRuntimesHashCode(s.runtimes)
	h *= 1000003
	h ^= sourceLoadFlagsHashCode(s.loadFlags)
	h *= 1000003
	h ^= int32(s.estimatedSize)
	h *= 1000003
	h ^= jutil.StringHashCode(s.code.get())
	h *= 1000003
	h ^= jutil.StringHashCode(s.originalCode.get())
	return h
}

// Source_Builder is the builder for Source instances.
type Source_Builder struct {
	path             *jutil.Path
	sourceMap        string
	sourceUrl        string
	sourceMappingUrl string
	runtimes         []string
	loadFlags        []sourceLoadFlag
	estimatedSize    int
	code             lazyString
	hasCode          bool
	originalCode     lazyString
	hasOriginalCode  bool
}

// Java: AutoValue_Source.Builder()
func NewSource_Builder() *Source_Builder {
	return &Source_Builder{
		path:    source_DEV_NULL,
		hasCode: true,
	}
}

// Java: Source.Builder.setPath(Path)
func (b *Source_Builder) SetPath(path *jutil.Path) *Source_Builder {
	if path == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null path"))
	}
	b.path = path
	return b
}

// Java: Source.Builder.setSourceMap(String)
func (b *Source_Builder) SetSourceMap(sourceMap string) *Source_Builder {
	b.sourceMap = sourceMap
	return b
}

// Java: Source.Builder.setSourceUrl(String)
func (b *Source_Builder) SetSourceUrl(sourceUrl string) *Source_Builder {
	b.sourceUrl = sourceUrl
	return b
}

// Java: Source.Builder.setSourceMappingUrl(String)
func (b *Source_Builder) SetSourceMappingUrl(sourceMappingUrl string) *Source_Builder {
	b.sourceMappingUrl = sourceMappingUrl
	return b
}

func appendUniqueRuntimes(dst []string, runtimes ...string) []string {
	for _, r := range runtimes {
		if !slices.Contains(dst, r) {
			dst = append(dst, r)
		}
	}
	return dst
}

// Java: Source.Builder.setRuntimes(ImmutableSet<String>)
func (b *Source_Builder) SetRuntimes(runtimes []string) *Source_Builder {
	if runtimes == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null runtimes"))
	}
	if len(runtimes) == 0 {
		b.runtimes = nil
		return b
	}
	b.runtimes = appendUniqueRuntimes(make([]string, 0, len(runtimes)), runtimes...)
	return b
}

// Java: Source.Builder.setLoadFlags(ImmutableMap<String, String>)
func (b *Source_Builder) SetLoadFlags(flags map[string]string) *Source_Builder {
	if flags == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null loadFlags"))
	}
	if len(flags) == 0 {
		b.loadFlags = nil
		return b
	}
	entries := make([]sourceLoadFlag, 0, len(flags))
	for k, v := range flags { // order-independent: entries are sorted by key immediately below
		entries = append(entries, sourceLoadFlag{key: k, value: v})
	}
	slices.SortFunc(entries, func(a, c sourceLoadFlag) int {
		return strings.Compare(a.key, c.key)
	})
	b.loadFlags = entries
	return b
}

// Java: Source.Builder.setEstimatedSize(int)
func (b *Source_Builder) SetEstimatedSize(estimatedSize int) *Source_Builder {
	b.estimatedSize = estimatedSize
	return b
}

// Java: Source.Builder.setCode(Supplier<String>)
func (b *Source_Builder) SetCodeFromSupplier(code func() string) *Source_Builder {
	jutil.CheckNotNull(code)
	b.code = lazyString{fn: sync.OnceValue(code)}
	b.hasCode = true
	return b
}

// Java: Source.Builder.setCode(String)
func (b *Source_Builder) SetCode(code string) *Source_Builder {
	b.code = lazyString{val: code}
	b.hasCode = true
	return b
}

// Java: Source.Builder.setOriginalCode(String)
func (b *Source_Builder) SetOriginalCode(code string) *Source_Builder {
	b.originalCode = lazyString{val: code}
	b.hasOriginalCode = true
	return b
}

// Java: Source.Builder.addRuntime(String...)
func (b *Source_Builder) AddRuntime(runtimes ...string) *Source_Builder {
	if len(runtimes) == 0 {
		return b
	}
	next := make([]string, len(b.runtimes), len(b.runtimes)+len(runtimes))
	copy(next, b.runtimes)
	b.runtimes = appendUniqueRuntimes(next, runtimes...)
	return b
}

// Java: Source.Builder.build()
func (b *Source_Builder) Build() *Source {
	if !b.hasCode {
		panic(jutil.NewIllegalStateExceptionMsg("Property \"codeSupplier\" has not been set"))
	}
	if b.path == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties: path"))
	}
	orig := b.originalCode
	if !b.hasOriginalCode {
		orig = b.code
	}
	return &Source{
		path:             b.path,
		sourceMap:        b.sourceMap,
		sourceUrl:        b.sourceUrl,
		sourceMappingUrl: b.sourceMappingUrl,
		runtimes:         b.runtimes,
		loadFlags:        b.loadFlags,
		estimatedSize:    b.estimatedSize,
		code:             b.code,
		originalCode:     orig,
	}
}

// Java: Source.Builder.runtimes()
func (b *Source_Builder) Runtimes() []string {
	if len(b.runtimes) == 0 {
		return []string{}
	}
	return slices.Clone(b.runtimes)
}

// Source_Transformer is an automorphic transformation on sources.
type Source_Transformer interface {
	// Transform is the main transformation method.
	Transform(input *Source) *Source
}

// Source_TransformerFunc adapts a function to Source_Transformer.
//
// Go-only: language-adapter — Java lambdas implementing the functional
// interface.
type Source_TransformerFunc func(input *Source) *Source

// Transform calls f.
func (f Source_TransformerFunc) Transform(input *Source) *Source { return f(input) }

// Java: Source.Transformer.of(Function<Source, Source>)
func Source_Transformer_Of(function func(*Source) *Source) Source_Transformer {
	return Source_TransformerFunc(function)
}
