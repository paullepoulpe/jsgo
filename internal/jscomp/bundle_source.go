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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Source is an abstract representation of a source file.
//
// Go-only: language-adapter — Source is an @AutoValue: the abstract class
// and the generated AutoValue_Source (fields, equals, hashCode, toString,
// Builder) are one struct here (PORTING.md §4).
type Source struct {
	path                 *jutil.Path
	sourceMap            string
	sourceUrl            string
	sourceMappingUrl     string
	runtimes             *jutil.LinkedHashSet[string]
	loadFlags            *jutil.LinkedHashMap[string, string]
	estimatedSize        int
	codeSupplier         *Source_Lazy[string]
	originalCodeSupplier *Source_Lazy[string]
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
func (s *Source) Runtimes() *jutil.LinkedHashSet[string] { return s.runtimes }

// LoadFlags are the load flags, specifying module type and language level.
//
// Java: Source.loadFlags()
func (s *Source) LoadFlags() *jutil.LinkedHashMap[string, string] { return s.loadFlags }

// EstimatedSize is a best estimate of the size of this source (in case the
// source itself is not yet loaded.
//
// Java: Source.estimatedSize()
func (s *Source) EstimatedSize() int { return s.estimatedSize }

// Code is the actual code in this source file.
//
// Java: Source.code()
func (s *Source) Code() string { return s.CodeSupplier().Get() }

// OriginalCode is the untransformed code from the original source file.
//
// Java: Source.originalCode()
func (s *Source) OriginalCode() string { return s.OriginalCodeSupplier().Get() }

// ToBuilder copies the data from this source to a new builder.
//
// Java: Source.toBuilder()
func (s *Source) ToBuilder() *Source_Builder {
	return &Source_Builder{
		path:                 s.path,
		sourceMap:            &s.sourceMap,
		sourceUrl:            &s.sourceUrl,
		sourceMappingUrl:     &s.sourceMappingUrl,
		runtimes:             s.runtimes,
		loadFlags:            s.loadFlags,
		estimatedSize:        s.estimatedSize,
		codeSupplier:         s.codeSupplier,
		originalCodeSupplier: s.originalCodeSupplier,
		set0:                 1,
	}
}

// Source_Builder_ makes a new empty builder (Source.builder()).
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name Source_Builder for the nested Builder type.
func Source_Builder_() *Source_Builder {
	return NewSource_Builder().
		SetPath(source_DEV_NULL).
		SetCode("").
		SetOriginalCodeSupplier(nil).
		SetSourceMap("").
		SetSourceUrl("").
		SetSourceMappingUrl("").
		SetRuntimes(jutil.NewLinkedHashSet[string]()).
		SetLoadFlags(jutil.NewLinkedHashMap[string, string]()).
		SetEstimatedSize(0)
}

var source_DEV_NULL = jutil.Path_Of("/dev/null")

// Internal-only properties: the code suppliers are necessary for lazy
// bundling, but we cannot use an ordinary supplier since we need guarantees
// about equals and hash code. Thus, we use an internal-only Supplier
// subtype.

// Java: Source.codeSupplier()
func (s *Source) CodeSupplier() *Source_Lazy[string] { return s.codeSupplier }

// Java: Source.originalCodeSupplier()
func (s *Source) OriginalCodeSupplier() *Source_Lazy[string] { return s.originalCodeSupplier }

// Java: AutoValue_Source.toString()
func (s *Source) String() string {
	return "Source{" +
		"path=" + s.path.String() + ", " +
		"sourceMap=" + s.sourceMap + ", " +
		"sourceUrl=" + s.sourceUrl + ", " +
		"sourceMappingUrl=" + s.sourceMappingUrl + ", " +
		"runtimes=" + s.runtimes.String() + ", " +
		"loadFlags=" + s.loadFlags.String() + ", " +
		"estimatedSize=" + jutil.ToJavaString(s.estimatedSize) + ", " +
		"codeSupplier=" + s.codeSupplier.String() + ", " +
		"originalCodeSupplier=" + jutil.ToJavaString(s.originalCodeSupplier) +
		"}"
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
	return s.path.Equals(that.Path()) &&
		s.sourceMap == that.SourceMap() &&
		s.sourceUrl == that.SourceUrl() &&
		s.sourceMappingUrl == that.SourceMappingUrl() &&
		s.runtimes.Equals(that.Runtimes()) &&
		s.loadFlags.Equals(that.LoadFlags()) &&
		s.estimatedSize == that.EstimatedSize() &&
		s.codeSupplier.Equals(that.CodeSupplier()) &&
		(func() bool {
			if s.originalCodeSupplier == nil {
				return that.OriginalCodeSupplier() == nil
			}
			return s.originalCodeSupplier.Equals(that.OriginalCodeSupplier())
		})()
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
	h ^= s.runtimes.HashCode()
	h *= 1000003
	h ^= s.loadFlags.HashCode()
	h *= 1000003
	h ^= int32(s.estimatedSize)
	h *= 1000003
	h ^= s.codeSupplier.HashCode()
	h *= 1000003
	if s.originalCodeSupplier != nil {
		h ^= s.originalCodeSupplier.HashCode()
	}
	return h
}

// Source_Builder is the builder for Source instances
// (AutoValue_Source.Builder).
type Source_Builder struct {
	path                 *jutil.Path
	sourceMap            *string
	sourceUrl            *string
	sourceMappingUrl     *string
	runtimes             *jutil.LinkedHashSet[string]
	loadFlags            *jutil.LinkedHashMap[string, string]
	estimatedSize        int
	codeSupplier         *Source_Lazy[string]
	originalCodeSupplier *Source_Lazy[string]
	set0                 byte
}

// Java: AutoValue_Source.Builder()
func NewSource_Builder() *Source_Builder { return &Source_Builder{} }

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
	b.sourceMap = &sourceMap
	return b
}

// Java: Source.Builder.setSourceUrl(String)
func (b *Source_Builder) SetSourceUrl(sourceUrl string) *Source_Builder {
	b.sourceUrl = &sourceUrl
	return b
}

// Java: Source.Builder.setSourceMappingUrl(String)
func (b *Source_Builder) SetSourceMappingUrl(sourceMappingUrl string) *Source_Builder {
	b.sourceMappingUrl = &sourceMappingUrl
	return b
}

// Java: Source.Builder.setRuntimes(ImmutableSet<String>)
func (b *Source_Builder) SetRuntimes(runtimes *jutil.LinkedHashSet[string]) *Source_Builder {
	if runtimes == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null runtimes"))
	}
	b.runtimes = runtimes
	return b
}

// Java: Source.Builder.setLoadFlags(ImmutableMap<String, String>)
func (b *Source_Builder) SetLoadFlags(flags *jutil.LinkedHashMap[string, string]) *Source_Builder {
	if flags == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null loadFlags"))
	}
	b.loadFlags = flags
	return b
}

// Java: Source.Builder.setEstimatedSize(int)
func (b *Source_Builder) SetEstimatedSize(estimatedSize int) *Source_Builder {
	b.estimatedSize = estimatedSize
	b.set0 |= 1
	return b
}

// Java: Source.Builder.setCode(Supplier<String>)
func (b *Source_Builder) SetCodeFromSupplier(code func() string) *Source_Builder {
	return b.SetCodeSupplier(Source_Lazy_Memoize(code))
}

// Java: Source.Builder.setCode(String)
func (b *Source_Builder) SetCode(code string) *Source_Builder {
	return b.SetCodeSupplier(Source_Lazy_OfInstance(code))
}

// Java: Source.Builder.setOriginalCode(String)
func (b *Source_Builder) SetOriginalCode(code string) *Source_Builder {
	return b.SetOriginalCodeSupplier(Source_Lazy_OfInstance(code))
}

// Java: Source.Builder.addRuntime(String...)
func (b *Source_Builder) AddRuntime(runtimes ...string) *Source_Builder {
	set := jutil.NewLinkedHashSet[string]()
	for r := range b.Runtimes().All() {
		set.Add(r)
	}
	set.AddAll(runtimes...)
	return b.SetRuntimes(set)
}

// Java: Source.Builder.setDependencyInfo(DependencyInfo)
func (b *Source_Builder) SetDependencyInfo(info DependencyInfo) *Source_Builder {
	// TODO(sdh): consider whether to set path.
	return b.SetLoadFlags(info.GetLoadFlags())
}

// Java: Source.Builder.build()
func (b *Source_Builder) Build() *Source {
	if b.OriginalCodeSupplier() == nil {
		b.SetOriginalCodeSupplier(b.CodeSupplier())
	}
	return b.AutoBuild()
}

// Java: Source.Builder.setCodeSupplier(Lazy<String>)
func (b *Source_Builder) SetCodeSupplier(codeSupplier *Source_Lazy[string]) *Source_Builder {
	if codeSupplier == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null codeSupplier"))
	}
	b.codeSupplier = codeSupplier
	return b
}

// Java: Source.Builder.setOriginalCodeSupplier(Lazy<String>)
func (b *Source_Builder) SetOriginalCodeSupplier(originalCodeSupplier *Source_Lazy[string]) *Source_Builder {
	b.originalCodeSupplier = originalCodeSupplier
	return b
}

// Java: Source.Builder.runtimes()
func (b *Source_Builder) Runtimes() *jutil.LinkedHashSet[string] {
	if b.runtimes == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Property \"runtimes\" has not been set"))
	}
	return b.runtimes
}

// Java: Source.Builder.codeSupplier()
func (b *Source_Builder) CodeSupplier() *Source_Lazy[string] {
	if b.codeSupplier == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Property \"codeSupplier\" has not been set"))
	}
	return b.codeSupplier
}

// Java: Source.Builder.originalCodeSupplier()
func (b *Source_Builder) OriginalCodeSupplier() *Source_Lazy[string] { return b.originalCodeSupplier }

// Java: Source.Builder.autoBuild()
func (b *Source_Builder) AutoBuild() *Source {
	if b.set0 != 1 || b.path == nil || b.sourceMap == nil || b.sourceUrl == nil || b.sourceMappingUrl == nil ||
		b.runtimes == nil || b.loadFlags == nil || b.codeSupplier == nil {
		missing := jutil.NewStringBuilder()
		if b.path == nil {
			missing.Append(" path")
		}
		if b.sourceMap == nil {
			missing.Append(" sourceMap")
		}
		if b.sourceUrl == nil {
			missing.Append(" sourceUrl")
		}
		if b.sourceMappingUrl == nil {
			missing.Append(" sourceMappingUrl")
		}
		if b.runtimes == nil {
			missing.Append(" runtimes")
		}
		if b.loadFlags == nil {
			missing.Append(" loadFlags")
		}
		if b.set0&1 == 0 {
			missing.Append(" estimatedSize")
		}
		if b.codeSupplier == nil {
			missing.Append(" codeSupplier")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return &Source{
		path:                 b.path,
		sourceMap:            *b.sourceMap,
		sourceUrl:            *b.sourceUrl,
		sourceMappingUrl:     *b.sourceMappingUrl,
		runtimes:             b.runtimes,
		loadFlags:            b.loadFlags,
		estimatedSize:        b.estimatedSize,
		codeSupplier:         b.codeSupplier,
		originalCodeSupplier: b.originalCodeSupplier,
	}
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

// Source_Transformer_Identity returns an identity transformer.
//
// Java: Source.Transformer.identity()
func Source_Transformer_Identity() Source_Transformer {
	return Source_TransformerFunc(func(x *Source) *Source { return x })
}

// Source_Transformer_AsFunction converts this Transformer to a Function.
//
// Java: Source.Transformer.asFunction() (default method)
func Source_Transformer_AsFunction(t Source_Transformer) func(*Source) *Source {
	return t.Transform
}

// Source_Transformer_AndThen concatenates two Transformers.
//
// Java: Source.Transformer.andThen(Transformer) (default method)
func Source_Transformer_AndThen(t Source_Transformer, after Source_Transformer) Source_Transformer {
	before := t
	return Source_TransformerFunc(func(x *Source) *Source { return after.Transform(before.Transform(x)) })
}

// Source_Transformer_Compose concatenates two Transformers.
//
// Java: Source.Transformer.compose(Transformer) (default method)
func Source_Transformer_Compose(t Source_Transformer, before Source_Transformer) Source_Transformer {
	after := t
	return Source_TransformerFunc(func(x *Source) *Source { return after.Transform(before.Transform(x)) })
}

// Source_Lazy is essentially the same as Supplier, but wraps equals and
// hashCode.
type Source_Lazy[T any] struct {
	get       func() T
	className string
}

// Get is Supplier.get().
func (l *Source_Lazy[T]) Get() T { return l.get() }

// Java: Source.Lazy.equals(Object)
func (l *Source_Lazy[T]) Equals(other any) bool {
	o, ok := other.(*Source_Lazy[T])
	return ok && o != nil && jutil.ObjectsEquals(l.Get(), o.Get())
}

// Java: Source.Lazy.hashCode()
func (l *Source_Lazy[T]) HashCode() int32 {
	return jutil.ObjectHashCode(l.Get())
}

// String is Object.toString() of the anonymous Lazy subclass.
func (l *Source_Lazy[T]) String() string {
	return jutil.ObjectToString(l.className, l)
}

// Source_Lazy_OfInstance returns a Lazy that always returns the same
// instance.
//
// Java: Source.Lazy.ofInstance(T)
func Source_Lazy_OfInstance[T any](instance T) *Source_Lazy[T] {
	return &Source_Lazy[T]{
		get:       func() T { return instance },
		className: "com.google.javascript.jscomp.bundle.Source$Lazy$1",
	}
}

// Source_Lazy_Memoize returns a Lazy from a memoized supplier.
//
// Java: Source.Lazy.memoize(Supplier<T>)
func Source_Lazy_Memoize[T any](supplier func() T) *Source_Lazy[T] {
	memoized := source_Suppliers_Memoize(supplier)
	return &Source_Lazy[T]{
		get:       memoized,
		className: "com.google.javascript.jscomp.bundle.Source$Lazy$2",
	}
}

// source_Suppliers_Memoize is Guava's Suppliers.memoize: the delegate is
// called until it returns normally once; then its value is cached.
func source_Suppliers_Memoize[T any](delegate func() T) func() T {
	initialized := false
	var value T
	return func() T {
		if !initialized {
			t := delegate()
			value = t
			initialized = true
			delegate = nil
			return t
		}
		return value
	}
}
