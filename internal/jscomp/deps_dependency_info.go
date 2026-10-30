// Port of com/google/javascript/jscomp/deps/DependencyInfo.java @ closure-compiler v20261021
// (with the AutoValue-generated deps/AutoValue_DependencyInfo_Require.java)
// Port of com/google/javascript/jscomp/deps/DependencyResolver.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/deps/SourceFile.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/deps/ServiceException.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/deps/VirtualFile.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
 * Copyright 2008 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// LoadFlag is a single (Key, Value) entry in LoadFlags.
type LoadFlag struct {
	Key   string
	Value string
}

// LoadFlags is an immutable insertion-ordered map of JS file load flags (such
// as "lang" and "module") backed by a compact value slice.
type LoadFlags struct {
	entries []LoadFlag
}

var (
	_ jutil.MapSource[string, string] = LoadFlags{}

	EmptyLoadFlags      = LoadFlags{}
	loadFlagsGoogArray  = [1]LoadFlag{{Key: "module", Value: "goog"}}
	loadFlagsES6Array   = [1]LoadFlag{{Key: "module", Value: "es6"}}
	GoogModuleLoadFlags = LoadFlags{entries: loadFlagsGoogArray[:]}
	ES6ModuleLoadFlags  = LoadFlags{entries: loadFlagsES6Array[:]}
)

// NewLoadFlags constructs LoadFlags from alternating key, value pairs.
func NewLoadFlags(kv ...string) LoadFlags {
	if len(kv) == 0 {
		return EmptyLoadFlags
	}
	var buf [4]LoadFlag
	entries := buf[:0]
	for i := 0; i+1 < len(kv); i += 2 {
		putLoadFlag(&entries, kv[i], kv[i+1])
	}
	flags, _ := normalizeLoadFlagsEntries(entries, true)
	return flags
}

func putLoadFlag(entries *[]LoadFlag, k, v string) {
	s := *entries
	for i := range s {
		if s[i].Key == k {
			s[i].Value = v
			return
		}
	}
	*entries = append(s, LoadFlag{Key: k, Value: v})
}

// Len returns the number of entries in f.
func (f LoadFlags) Len() int { return len(f.entries) }

// Size implements jutil.MapSource.
func (f LoadFlags) Size() int { return len(f.entries) }

// IsEmpty reports whether f has zero entries.
func (f LoadFlags) IsEmpty() bool { return len(f.entries) == 0 }

// Entries returns the underlying entries in insertion order.
func (f LoadFlags) Entries() []LoadFlag { return f.entries }

// Get returns (value, true) if key is present, or ("", false) otherwise.
func (f LoadFlags) Get(key string) (string, bool) {
	for _, e := range f.entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

// Lookup is an alias for Get.
func (f LoadFlags) Lookup(key string) (string, bool) { return f.Get(key) }

// GetOrDefault returns the value for key, or def if absent.
func (f LoadFlags) GetOrDefault(key, def string) string {
	if v, ok := f.Get(key); ok {
		return v
	}
	return def
}

// All yields (key, value) pairs in insertion order.
func (f LoadFlags) All() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for _, e := range f.entries {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	}
}

// Equals reports whether f and other contain the same (key, value) mappings.
func (f LoadFlags) Equals(other LoadFlags) bool {
	if len(f.entries) != len(other.entries) {
		return false
	}
	if slices.Equal(f.entries, other.entries) {
		return true
	}
	for _, e := range f.entries {
		if ov, ok := other.Get(e.Key); !ok || ov != e.Value {
			return false
		}
	}
	return true
}

// HashCode computes Java Map.hashCode(): sum of (key.hashCode() ^ value.hashCode()).
func (f LoadFlags) HashCode() int32 {
	var sum int32
	for _, e := range f.entries {
		sum += jutil.StringHashCode(e.Key) ^ jutil.StringHashCode(e.Value)
	}
	return sum
}

// String formats f like Java AbstractMap.toString(): "{k1=v1, k2=v2}".
func (f LoadFlags) String() string {
	if len(f.entries) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range f.entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e.Key)
		sb.WriteByte('=')
		sb.WriteString(e.Value)
	}
	sb.WriteByte('}')
	return sb.String()
}

// DependencyInfo is a data structure for JS dependency information for a
// single .js file.
type DependencyInfo interface {
	GetName() string
	GetPathRelativeToClosureBase() string
	GetProvides() []string
	GetRequires() []*DependencyInfo_Require
	GetRequiredSymbols() []string
	GetTypeRequires() []string
	GetLoadFlags() LoadFlags
	IsEs6Module() bool
	IsGoogModule() bool
	GetHasExternsAnnotation() bool
	GetHasNoCompileAnnotation() bool
}

// Java: DependencyInfo.getRequiredSymbols() (default method)
func DependencyInfo_GetRequiredSymbols(d DependencyInfo) []string {
	return DependencyInfo_Require_AsSymbolList(d.GetRequires())
}

// Java: DependencyInfo.isEs6Module() (default method)
func DependencyInfo_IsEs6Module(d DependencyInfo) bool {
	v, ok := d.GetLoadFlags().Get("module")
	return ok && v == "es6"
}

// Java: DependencyInfo.isGoogModule() (default method)
func DependencyInfo_IsGoogModule(d DependencyInfo) bool {
	v, ok := d.GetLoadFlags().Get("module")
	return ok && v == "goog"
}

// DependencyInfo_Require is a dependency link between two files, e.g.
// goog.require('namespace'), import 'file';
type DependencyInfo_Require struct {
	symbol  string
	rawText string
	type_   DependencyInfo_Require_Type
}

// DependencyInfo_Require_Type is DependencyInfo.Require.Type.
type DependencyInfo_Require_Type int32

const (
	DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL DependencyInfo_Require_Type = iota
	DependencyInfo_Require_Type_ES6_IMPORT
	DependencyInfo_Require_Type_PARSED_FROM_DEPS
	DependencyInfo_Require_Type_COMMON_JS
	DependencyInfo_Require_Type_COMPILER_CHUNK
)

var dependencyInfo_Require_Type_names = [...]string{
	"GOOG_REQUIRE_SYMBOL", "ES6_IMPORT", "PARSED_FROM_DEPS", "COMMON_JS", "COMPILER_CHUNK",
}

func (t DependencyInfo_Require_Type) String() string { return dependencyInfo_Require_Type_names[t] }

// Java: DependencyInfo.Require.BASE
var DependencyInfo_Require_BASE = DependencyInfo_Require_GoogRequireSymbol("goog")

// Java: DependencyInfo.Require.asSymbolList(Iterable<Require>)
func DependencyInfo_Require_AsSymbolList(requires []*DependencyInfo_Require) []string {
	out := make([]string, len(requires))
	for i, r := range requires {
		out[i] = r.symbol
	}
	return out
}

// Java: DependencyInfo.Require.googRequireSymbol(String)
func DependencyInfo_Require_GoogRequireSymbol(symbol string) *DependencyInfo_Require {
	interned := rhino.RhinoStringPool_AddOrGet(symbol)
	return &DependencyInfo_Require{
		symbol:  interned,
		rawText: interned,
		type_:   DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL,
	}
}

// Java: DependencyInfo.Require.es6Import(String, String)
func DependencyInfo_Require_Es6Import(symbol string, rawPath string) *DependencyInfo_Require {
	return &DependencyInfo_Require{symbol: symbol, rawText: rawPath, type_: DependencyInfo_Require_Type_ES6_IMPORT}
}

// Java: DependencyInfo.Require.commonJs(String, String)
func DependencyInfo_Require_CommonJs(symbol string, rawPath string) *DependencyInfo_Require {
	return &DependencyInfo_Require{symbol: symbol, rawText: rawPath, type_: DependencyInfo_Require_Type_COMMON_JS}
}

// Java: DependencyInfo.Require.compilerChunk(String)
func DependencyInfo_Require_CompilerChunk(symbol string) *DependencyInfo_Require {
	return &DependencyInfo_Require{symbol: symbol, rawText: symbol, type_: DependencyInfo_Require_Type_COMPILER_CHUNK}
}

// Java: DependencyInfo.Require.parsedFromDeps(String)
func DependencyInfo_Require_ParsedFromDeps(symbol string) *DependencyInfo_Require {
	return &DependencyInfo_Require{symbol: symbol, rawText: symbol, type_: DependencyInfo_Require_Type_PARSED_FROM_DEPS}
}

// Java: DependencyInfo.Require.withSymbol(String)
func (r *DependencyInfo_Require) WithSymbol(symbol string) *DependencyInfo_Require {
	return &DependencyInfo_Require{symbol: symbol, rawText: r.rawText, type_: r.type_}
}

// Java: DependencyInfo.Require.getSymbol()
func (r *DependencyInfo_Require) GetSymbol() string { return r.symbol }

// Java: DependencyInfo.Require.getRawText()
func (r *DependencyInfo_Require) GetRawText() string { return r.rawText }

// Java: DependencyInfo.Require.getType()
func (r *DependencyInfo_Require) GetType() DependencyInfo_Require_Type { return r.type_ }

// Java: AutoValue_DependencyInfo_Require.toString()
func (r *DependencyInfo_Require) String() string {
	return "Require{" +
		"symbol=" + r.symbol + ", " +
		"rawText=" + r.rawText + ", " +
		"type=" + r.type_.String() +
		"}"
}

// Java: AutoValue_DependencyInfo_Require.equals(Object)
func (r *DependencyInfo_Require) Equals(o any) bool {
	that, ok := o.(*DependencyInfo_Require)
	return ok && that != nil &&
		(that == r || (r.symbol == that.symbol && r.rawText == that.rawText && r.type_ == that.type_))
}

// Java: AutoValue_DependencyInfo_Require.hashCode()
func (r *DependencyInfo_Require) HashCode() int32 {
	h := int32(1)
	h *= 1000003
	h ^= jutil.StringHashCode(r.symbol)
	h *= 1000003
	h ^= jutil.StringHashCode(r.rawText)
	h *= 1000003
	h ^= int32(r.type_)
	return h
}

// DependencyInfo_Util_WriteAddDependency prints a goog.addDependency call
// for a single DependencyInfo. It panics with an IOException if out does.
//
// Java: DependencyInfo.Util.writeAddDependency(Appendable, DependencyInfo)
func DependencyInfo_Util_WriteAddDependency(out jutil.Appendable, info DependencyInfo) {
	out.AppendCharSequence("goog.addDependency('")
	out.AppendCharSequence(info.GetPathRelativeToClosureBase())
	out.AppendCharSequence("', ")
	dependencyInfo_Util_writeJsArray(out, info.GetProvides())
	out.AppendCharSequence(", [")
	for i, r := range info.GetRequires() {
		if i > 0 {
			out.AppendCharSequence(", ")
		}
		out.AppendCharSequence("'")
		out.AppendCharSequence(strings.ReplaceAll(r.GetSymbol(), "'", "\\'"))
		out.AppendCharSequence("'")
	}
	out.AppendCharSequence("]")
	loadFlags := info.GetLoadFlags()
	if !loadFlags.IsEmpty() {
		out.AppendCharSequence(", ")
		dependencyInfo_Util_writeJsObject(out, loadFlags)
	}
	out.AppendCharSequence(");\n")
}

func dependencyInfo_Util_writeJsObject(out jutil.Appendable, m LoadFlags) {
	out.AppendCharSequence("{")
	for i, e := range m.Entries() {
		if i > 0 {
			out.AppendCharSequence(", ")
		}
		out.AppendCharSequence("'")
		out.AppendCharSequence(strings.ReplaceAll(e.Key, "'", "\\'"))
		out.AppendCharSequence("': '")
		out.AppendCharSequence(strings.ReplaceAll(e.Value, "'", "\\'"))
		out.AppendCharSequence("'")
	}
	out.AppendCharSequence("}")
}

func dependencyInfo_Util_writeJsArray(out jutil.Appendable, values []string) {
	out.AppendCharSequence("[")
	for i, arg := range values {
		if i > 0 {
			out.AppendCharSequence(", ")
		}
		out.AppendCharSequence("'")
		out.AppendCharSequence(strings.ReplaceAll(arg, "'", "\\'"))
		out.AppendCharSequence("'")
	}
	out.AppendCharSequence("]")
}

// DependencyResolver gets dependency lists for code or symbols.
// Every method throws ServiceException (as a panic).
type DependencyResolver interface {
	GetDependencies(code string) []string
	GetDependenciesOfSymbols(symbols []string) []string
	GetDependenciesWithSeen(code string, seen *jutil.LinkedHashSet[string], addClosureBaseFile bool) []string
	GetDependenciesOfSymbolsWithSeen(symbols []string, seen *jutil.LinkedHashSet[string]) []string
}

// Deps_SourceFile represents source files for the Compilation Service.
type Deps_SourceFile interface {
	GetName() string
	GetContent() string
	WasModified() bool
}

// ServiceException is an exception thrown by the compilation service.
type ServiceException struct {
	jutil.JavaException
}

func (*ServiceException) serviceException() {}

// IServiceException matches ServiceException.
type IServiceException interface {
	jutil.IException
	serviceException()
}

// Java: ServiceException(String)
func NewServiceException(message string) *ServiceException {
	e := &ServiceException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.deps.ServiceException", &message, nil)
	return e
}

// VirtualFile is used to represent a "virtual" file.
type VirtualFile struct {
	name string
	code string
}

var _ Deps_SourceFile = (*VirtualFile)(nil)

// Java: VirtualFile(String, String)
func NewVirtualFile(name string, code string) *VirtualFile {
	return &VirtualFile{name: name, code: code}
}

// Java: VirtualFile.getName()
func (f *VirtualFile) GetName() string { return f.name }

// Java: VirtualFile.getContent()
func (f *VirtualFile) GetContent() string { return f.code }

// Java: VirtualFile.wasModified()
func (f *VirtualFile) WasModified() bool { return false }

// Java: VirtualFile.hashCode() (Objects.hash(name, code))
func (f *VirtualFile) HashCode() int32 {
	h := int32(1)
	h = 31*h + jutil.StringHashCode(f.name)
	h = 31*h + jutil.StringHashCode(f.code)
	return h
}

// Java: VirtualFile.equals(Object)
func (f *VirtualFile) Equals(o any) bool {
	vf, ok := o.(*VirtualFile)
	return ok && vf != nil && f.name == vf.name && f.code == vf.code
}
