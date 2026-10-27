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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DependencyInfo is a data structure for JS dependency information for a
// single .js file.
//
// Java default methods are declared in the interface; their default bodies
// are the DependencyInfo_X functions, which implementations call.
type DependencyInfo interface {
	// GetName gets the unique name / path of this file.
	GetName() string

	// GetPathRelativeToClosureBase gets the path of this file relative to
	// Closure's base.js file.
	GetPathRelativeToClosureBase() string

	// GetProvides gets the symbols provided by this file.
	GetProvides() []string

	// GetRequires gets the symbols required by this file.
	GetRequires() []*DependencyInfo_Require

	GetRequiredSymbols() []string

	// GetTypeRequires gets the symbols type-required by this file (i.e. for
	// typechecking only).
	GetTypeRequires() []string

	// GetLoadFlags gets the loading information for this file.
	GetLoadFlags() *jutil.LinkedHashMap[string, string]

	// IsEs6Module reports whether the symbol is provided by an ES6 module.
	IsEs6Module() bool

	// IsGoogModule reports whether the symbol is provided by a goog module.
	IsGoogModule() bool

	// GetHasExternsAnnotation reports whether the file '@externs' annotation.
	GetHasExternsAnnotation() bool

	// GetHasNoCompileAnnotation reports whether the file has the '@nocompile'
	// annotation.
	GetHasNoCompileAnnotation() bool
}

// Java: DependencyInfo.getRequiredSymbols() (default method)
func DependencyInfo_GetRequiredSymbols(d DependencyInfo) []string {
	return DependencyInfo_Require_AsSymbolList(d.GetRequires())
}

// Java: DependencyInfo.isEs6Module() (default method)
func DependencyInfo_IsEs6Module(d DependencyInfo) bool {
	v, ok := d.GetLoadFlags().Lookup("module")
	return ok && v == "es6"
}

// Java: DependencyInfo.isGoogModule() (default method)
func DependencyInfo_IsGoogModule(d DependencyInfo) bool {
	v, ok := d.GetLoadFlags().Lookup("module")
	return ok && v == "goog"
}

// DependencyInfo_Require is a dependency link between two files, e.g.
// goog.require('namespace'), import 'file';
//
// Go-only: language-adapter — Require is an @AutoValue: the abstract class
// and the generated AutoValue_DependencyInfo_Require (fields, equals,
// hashCode, toString, Builder) are one struct here (PORTING.md §4).
type DependencyInfo_Require struct {
	symbol  string
	rawText string
	type_   DependencyInfo_Require_Type
}

// DependencyInfo_Require_Type is DependencyInfo.Require.Type.
type DependencyInfo_Require_Type int32

const (
	// Standard goog.require call for a symbol from a goog.provide or
	// goog.module.
	DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL DependencyInfo_Require_Type = iota
	// ES6 import statement.
	DependencyInfo_Require_Type_ES6_IMPORT
	// Parsed from an existing Closure dependency file.
	DependencyInfo_Require_Type_PARSED_FROM_DEPS
	// CommonJS require() call.
	DependencyInfo_Require_Type_COMMON_JS
	// Compiler chunk graph dependencies.
	DependencyInfo_Require_Type_COMPILER_CHUNK
)

var dependencyInfo_Require_Type_names = [...]string{
	"GOOG_REQUIRE_SYMBOL", "ES6_IMPORT", "PARSED_FROM_DEPS", "COMMON_JS", "COMPILER_CHUNK",
}

// String is Enum.name().
func (t DependencyInfo_Require_Type) String() string { return dependencyInfo_Require_Type_names[t] }

// Java: DependencyInfo.Require.BASE
var DependencyInfo_Require_BASE = DependencyInfo_Require_GoogRequireSymbol("goog")

// Java: DependencyInfo.Require.asSymbolList(Iterable<Require>)
func DependencyInfo_Require_AsSymbolList(requires []*DependencyInfo_Require) []string {
	out := make([]string, 0, len(requires))
	for _, r := range requires {
		out = append(out, r.GetSymbol())
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

// DependencyInfo_Require_Builder_ is DependencyInfo.Require.builder().
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name DependencyInfo_Require_Builder for the nested Builder type.
func DependencyInfo_Require_Builder_() *DependencyInfo_Require_Builder {
	return NewDependencyInfo_Require_Builder()
}

// Java: DependencyInfo.Require.toBuilder()
func (r *DependencyInfo_Require) ToBuilder() *DependencyInfo_Require_Builder {
	return &DependencyInfo_Require_Builder{symbol: &r.symbol, rawText: &r.rawText, type_: r.type_.ptr()}
}

// Java: DependencyInfo.Require.withSymbol(String)
func (r *DependencyInfo_Require) WithSymbol(symbol string) *DependencyInfo_Require {
	return &DependencyInfo_Require{symbol: symbol, rawText: r.rawText, type_: r.type_}
}

// GetSymbol returns the symbol provided by another DependencyInfo's
// GetProvides().
//
// Java: DependencyInfo.Require.getSymbol()
func (r *DependencyInfo_Require) GetSymbol() string { return r.symbol }

// GetRawText returns the raw text of the import string as it appears in the
// file. Used mostly for error reporting.
//
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
	if !ok || that == nil {
		return false
	}
	if that == r {
		return true
	}
	return r.symbol == that.GetSymbol() &&
		r.rawText == that.GetRawText() &&
		r.type_ == that.GetType()
}

// Java: AutoValue_DependencyInfo_Require.hashCode()
//
// Go-only: language-adapter — Java's Enum.hashCode() is the identity hash,
// which varies from run to run; the ordinal stands in for it.
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

func (t DependencyInfo_Require_Type) ptr() *DependencyInfo_Require_Type { return &t }

// DependencyInfo_Require_Builder is DependencyInfo.Require.Builder
// (AutoValue_DependencyInfo_Require.Builder).
type DependencyInfo_Require_Builder struct {
	symbol  *string
	rawText *string
	type_   *DependencyInfo_Require_Type
}

// Java: AutoValue_DependencyInfo_Require.Builder()
func NewDependencyInfo_Require_Builder() *DependencyInfo_Require_Builder {
	return &DependencyInfo_Require_Builder{}
}

// Java: DependencyInfo.Require.Builder.setType(Type)
func (b *DependencyInfo_Require_Builder) SetType(value DependencyInfo_Require_Type) *DependencyInfo_Require_Builder {
	b.type_ = value.ptr()
	return b
}

// Java: DependencyInfo.Require.Builder.setRawText(String)
func (b *DependencyInfo_Require_Builder) SetRawText(rawText string) *DependencyInfo_Require_Builder {
	b.rawText = &rawText
	return b
}

// Java: DependencyInfo.Require.Builder.setSymbol(String)
func (b *DependencyInfo_Require_Builder) SetSymbol(value string) *DependencyInfo_Require_Builder {
	b.symbol = &value
	return b
}

// Java: DependencyInfo.Require.Builder.build()
func (b *DependencyInfo_Require_Builder) Build() *DependencyInfo_Require {
	if b.symbol == nil || b.rawText == nil || b.type_ == nil {
		var missing strings.Builder
		if b.symbol == nil {
			missing.WriteString(" symbol")
		}
		if b.rawText == nil {
			missing.WriteString(" rawText")
		}
		if b.type_ == nil {
			missing.WriteString(" type")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return &DependencyInfo_Require{symbol: *b.symbol, rawText: *b.rawText, type_: *b.type_}
}

// DependencyInfo_Util holds utility methods.
type DependencyInfo_Util struct{}

// DependencyInfo_Util_WriteAddDependency prints a goog.addDependency call
// for a single DependencyInfo. It panics with an IOException if out does.
//
// Java: DependencyInfo.Util.writeAddDependency(Appendable, DependencyInfo)
func DependencyInfo_Util_WriteAddDependency(out jutil.Appendable, info DependencyInfo) {
	out.AppendCharSequence("goog.addDependency('")
	out.AppendCharSequence(info.GetPathRelativeToClosureBase())
	out.AppendCharSequence("', ")
	dependencyInfo_Util_writeJsArray(out, info.GetProvides())
	out.AppendCharSequence(", ")
	dependencyInfo_Util_writeJsArray(out, DependencyInfo_Require_AsSymbolList(info.GetRequires()))
	loadFlags := info.GetLoadFlags()
	if !loadFlags.IsEmpty() {
		out.AppendCharSequence(", ")
		dependencyInfo_Util_writeJsObject(out, loadFlags)
	}
	out.AppendCharSequence(");\n")
}

// Prints a map as a JS object literal.
func dependencyInfo_Util_writeJsObject(out jutil.Appendable, m *jutil.LinkedHashMap[string, string]) {
	var entries []string
	for k, v := range m.All() {
		key := strings.ReplaceAll(k, "'", "\\'")
		value := strings.ReplaceAll(v, "'", "\\'")
		entries = append(entries, "'"+key+"': '"+value+"'")
	}
	out.AppendCharSequence("{")
	out.AppendCharSequence(strings.Join(entries, ", "))
	out.AppendCharSequence("}")
}

// Prints a list of strings formatted as a JavaScript array of string
// literals.
func dependencyInfo_Util_writeJsArray(out jutil.Appendable, values []string) {
	quoted := make([]string, len(values))
	for i, arg := range values {
		quoted[i] = "'" + strings.ReplaceAll(arg, "'", "\\'") + "'"
	}
	out.AppendCharSequence("[")
	out.AppendCharSequence(strings.Join(quoted, ", "))
	out.AppendCharSequence("]")
}

// DependencyResolver is the interface for DependencyResolver to keep in
// line with com.google.javascript.jscomp.deps.DefaultDependencyResolver,
// which doesn't provide an interface. This is so that later on, we can
// merge back with it.
//
// Every method throws ServiceException (as a panic).
type DependencyResolver interface {
	// GetDependencies gets a list of dependencies for the provided code.
	//
	// Java: DependencyResolver.getDependencies(String)
	GetDependencies(code string) []string

	// GetDependenciesOfSymbols gets a list of dependencies for the provided
	// list of symbols.
	//
	// Java: DependencyResolver.getDependencies(Collection<String>)
	GetDependenciesOfSymbols(symbols []string) []string

	// GetDependenciesWithSeen: code is the raw code to be parsed for
	// requires; seen is the set of already seen symbols;
	// addClosureBaseFile indicates whether the closure base file should be
	// added to the dependency list. Returns a list of filenames for each of
	// the dependencies for the provided code.
	//
	// Java: DependencyResolver.getDependencies(String, Set<String>, boolean)
	GetDependenciesWithSeen(code string, seen *jutil.LinkedHashSet[string], addClosureBaseFile bool) []string

	// GetDependenciesOfSymbolsWithSeen: symbols is a list of required
	// symbols; seen is the set of already seen symbols. Returns a list of
	// filenames for each of the required symbols.
	//
	// Java: DependencyResolver.getDependencies(Collection<String>, Set<String>)
	GetDependenciesOfSymbolsWithSeen(symbols []string, seen *jutil.LinkedHashSet[string]) []string
}

// Deps_SourceFile is an interface for representing source files for the
// Compilation Service.
//
// Go-only: language-adapter — com.google.javascript.jscomp.deps.SourceFile
// and com.google.javascript.jscomp.SourceFile are both in the Go package
// jscomp, so this one carries its subpackage name.
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
	jutil.InitThrowable(e, "com.google.javascript.jscomp.deps.ServiceException", jutil.StringPtr(message), nil)
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
	if vf, ok := o.(*VirtualFile); ok && vf != nil {
		return f.name == vf.name && f.code == vf.code
	}
	return false
}
