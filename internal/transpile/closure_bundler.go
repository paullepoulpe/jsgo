// Port of com/google/javascript/jscomp/deps/ClosureBundler.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package transpile

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TODO(user): Convert this class to a builder/autovalue.

// ClosureBundler is a utility class to assist in creating JS bundle files.
//
// Java: final class deps.ClosureBundler
type ClosureBundler struct {
	transpiler          Transpiler
	es6ModuleTranspiler Transpiler

	mode           closureBundler_EvalMode
	sourceUrl      string
	hasSourceUrl   bool
	path           string
	embedSourcemap bool

	// TODO(sdh): This cache should be moved out into a higher level, but is
	// currently required due to the API that source maps must be accessible
	// via just a path (and not the file contents).
	sourceMapCache *closureBundler_SourceMapCache
}

type closureBundler_SourceMapCache struct {
	mu sync.Mutex
	m  map[string]string
}

// Java: ClosureBundler()
func NewClosureBundler() *ClosureBundler {
	return NewClosureBundlerWithTranspiler(Transpiler_NULL)
}

// Java: ClosureBundler(Transpiler)
func NewClosureBundlerWithTranspiler(transpiler Transpiler) *ClosureBundler {
	return NewClosureBundlerWithEs6ModuleTranspiler(transpiler, BaseTranspiler_LATEST_TRANSPILER)
}

// Java: ClosureBundler(Transpiler, Transpiler)
func NewClosureBundlerWithEs6ModuleTranspiler(transpiler Transpiler, es6ModuleTranspiler Transpiler) *ClosureBundler {
	return &ClosureBundler{
		transpiler:          transpiler,
		es6ModuleTranspiler: es6ModuleTranspiler,
		mode:                closureBundler_EvalMode_NORMAL,
		path:                "unknown_source",
		sourceMapCache:      &closureBundler_SourceMapCache{m: make(map[string]string)},
	}
}

// Java: ClosureBundler.withTranspilers(Transpiler, Transpiler)
func (b *ClosureBundler) WithTranspilers(
	newTranspiler Transpiler, newEs6ModuleTranspiler Transpiler) *ClosureBundler {
	cp := *b
	cp.transpiler = newTranspiler
	cp.es6ModuleTranspiler = newEs6ModuleTranspiler
	return &cp
}

// Java: ClosureBundler.withTranspiler(Transpiler)
func (b *ClosureBundler) WithTranspiler(newTranspiler Transpiler) *ClosureBundler {
	return b.WithTranspilers(newTranspiler, b.es6ModuleTranspiler)
}

// Java: ClosureBundler.withEs6ModuleTranspiler(Transpiler)
func (b *ClosureBundler) WithEs6ModuleTranspiler(newEs6ModuleTranspiler Transpiler) *ClosureBundler {
	return b.WithTranspilers(b.transpiler, newEs6ModuleTranspiler)
}

// Java: ClosureBundler.disableJ2clMinifier()
func (b *ClosureBundler) DisableJ2clMinifier() *ClosureBundler {
	cp := *b
	return &cp
}

// Java: ClosureBundler.useEval(boolean)
func (b *ClosureBundler) UseEval(useEval bool) *ClosureBundler {
	cp := *b
	cp.mode = closureBundler_EvalMode_NORMAL
	if useEval {
		cp.mode = closureBundler_EvalMode_EVAL
	}
	return &cp
}

// Java: ClosureBundler.withSourceUrl(String)
func (b *ClosureBundler) WithSourceUrl(newSourceUrl string) *ClosureBundler {
	cp := *b
	cp.sourceUrl = newSourceUrl
	cp.hasSourceUrl = true
	return &cp
}

// Java: ClosureBundler.withPath(String)
func (b *ClosureBundler) WithPath(newPath string) *ClosureBundler {
	cp := *b
	cp.path = newPath
	return &cp
}

// Java: ClosureBundler.embedSourcemap()
func (b *ClosureBundler) EmbedSourcemap() *ClosureBundler {
	cp := *b
	cp.embedSourcemap = true
	return &cp
}

// ClosureBundler_AppendInput appends the contents of the string to the
// supplied appendable. It throws IOException.
//
// Java: ClosureBundler.appendInput(Appendable, DependencyInfo, String)
func ClosureBundler_AppendInput(out jutil.Appendable, info jscomp.DependencyInfo, contents string) {
	NewClosureBundler().AppendTo(out, info, contents)
}

// AppendTo appends the contents of the string to the supplied appendable.
// It throws IOException.
//
// Java: ClosureBundler.appendTo(Appendable, DependencyInfo, String)
func (b *ClosureBundler) AppendTo(out jutil.Appendable, info jscomp.DependencyInfo, code string) {
	if info.IsGoogModule() {
		b.mode.appendGoogModule(b.transpile(code), out, b.sourceUrl, b.hasSourceUrl)
	} else if info.IsEs6Module() && b.transpiler == Transpiler_NULL {
		// TODO(johnplaisted): Make the default transpiler the
		// ES_MODULE_TO_CJS_TRANSPILER. Currently some code is passing in
		// unicode identifiers in non-ES6 modules the compiler fails to parse.
		// Once this compiler bug is fixed we can always transpile.
		b.mode.appendTraditional(b.transpileEs6Module(code), out, b.sourceUrl, b.hasSourceUrl)
	} else {
		b.mode.appendTraditional(b.transpile(code), out, b.sourceUrl, b.hasSourceUrl)
	}
}

// AppendToFromFile appends the contents of the file to the supplied
// appendable. It throws IOException.
//
// Java: ClosureBundler.appendTo(Appendable, DependencyInfo, File, Charset)
func (b *ClosureBundler) AppendToFromFile(out jutil.Appendable, info jscomp.DependencyInfo, content *jutil.File, contentCharset *jutil.Charset) {
	b.AppendToFromCharSource(out, info, jutil.Files_AsCharSource(content, contentCharset))
}

// AppendToFromCharSource appends the contents of the CharSource to the
// supplied appendable. It throws IOException.
//
// Java: ClosureBundler.appendTo(Appendable, DependencyInfo, CharSource)
func (b *ClosureBundler) AppendToFromCharSource(out jutil.Appendable, info jscomp.DependencyInfo, content *jutil.CharSource) {
	b.AppendTo(out, info, content.Read())
}

// AppendRuntimeTo throws IOException.
//
// Java: ClosureBundler.appendRuntimeTo(Appendable)
func (b *ClosureBundler) AppendRuntimeTo(out jutil.Appendable) {
	runtime := b.transpiler.Runtime()
	if runtime != "" {
		b.mode.appendTraditional(runtime, out, "", false)
	}
	if b.transpiler == Transpiler_NULL {
		b.mode.appendTraditional(b.es6ModuleTranspiler.Runtime(), out, "", false)
	}
	b.mode.appendTraditional("this.CLOSURE_EVAL_PREFILTER = function(s) { return s; };", out, "", false)
	b.mode.appendTraditional("(function(thisValue){", out, "", false)
	// Check for Chrome <87 which does not eval properly in workers.
	b.mode.appendTraditional(
		"var isChrome87 = false; try {isChrome87 =  eval(trustedTypes.emptyScript) !=="+
			" trustedTypes.emptyScript } catch (e) {} if (typeof trustedTypes !=="+
			" 'undefined' && trustedTypes.createPolicy &&isChrome87 ) {",
		out,
		"", false)
	b.mode.appendTraditional(
		"  var policy = trustedTypes.createPolicy('goog#devserver',{ createScript: function(s){"+
			" return s; }});",
		out,
		"", false)
	b.mode.appendTraditional(
		"  thisValue.CLOSURE_EVAL_PREFILTER = policy.createScript.bind(policy);", out, "", false)
	b.mode.appendTraditional("}", out, "", false)
	b.mode.appendTraditional("})(this);", out, "", false)
}

// GetSourceMap: subclasses that need to provide a source map for any
// transformed input can return it with this method.
//
// Java: ClosureBundler.getSourceMap(String)
func (b *ClosureBundler) GetSourceMap(path string) string {
	b.sourceMapCache.mu.Lock()
	defer b.sourceMapCache.mu.Unlock()
	return b.sourceMapCache.m[path]
}

// Java: ClosureBundler.transpile(String, Transpiler)
func (b *ClosureBundler) transpileWithTranspiler(s string, t Transpiler) string {
	result := t.Transpile(b.path, s)
	b.sourceMapCache.mu.Lock()
	b.sourceMapCache.m[b.path] = result.SourceMap()
	b.sourceMapCache.mu.Unlock()
	if b.embedSourcemap {
		return result.EmbedSourcemapBase64().Transpiled()
	}
	return result.Transpiled()
}

// Java: ClosureBundler.transpile(String)
func (b *ClosureBundler) transpile(s string) string {
	return b.transpileWithTranspiler(s, b.transpiler)
}

// Java: ClosureBundler.transpileEs6Module(String)
func (b *ClosureBundler) transpileEs6Module(s string) string {
	return b.transpileWithTranspiler(s, b.es6ModuleTranspiler)
}

// Java: enum ClosureBundler.EvalMode
type closureBundler_EvalMode int32

const (
	closureBundler_EvalMode_EVAL closureBundler_EvalMode = iota
	closureBundler_EvalMode_NORMAL
)

// String is EvalMode.name().
func (m closureBundler_EvalMode) String() string {
	return [...]string{"EVAL", "NORMAL"}[m]
}

// Java: EvalMode.appendTraditional(String, Appendable, String)
func (m closureBundler_EvalMode) appendTraditional(s string, out jutil.Appendable, sourceUrl string, hasSourceUrl bool) {
	switch m {
	case closureBundler_EvalMode_EVAL:
		out.AppendCharSequence("eval(this.CLOSURE_EVAL_PREFILTER(\"")
		closureBundler_EscapeMode_ESCAPED.append(s, out)
		closureBundler_appendSourceUrl(out, closureBundler_EscapeMode_ESCAPED, sourceUrl, hasSourceUrl)
		out.AppendCharSequence("\"));\n")
	default: // NORMAL
		closureBundler_EscapeMode_NORMAL.append(s, out)
		closureBundler_appendSourceUrl(out, closureBundler_EscapeMode_NORMAL, sourceUrl, hasSourceUrl)
	}
}

// Java: EvalMode.appendGoogModule(String, Appendable, String)
func (m closureBundler_EvalMode) appendGoogModule(s string, out jutil.Appendable, sourceUrl string, hasSourceUrl bool) {
	switch m {
	case closureBundler_EvalMode_EVAL:
		out.AppendCharSequence("goog.loadModule(\"")
		closureBundler_EscapeMode_ESCAPED.append(s, out)
		closureBundler_appendSourceUrl(out, closureBundler_EscapeMode_ESCAPED, sourceUrl, hasSourceUrl)
		out.AppendCharSequence("\");\n")
	default: // NORMAL
		// add the prefix on the first line so the line numbers aren't
		// affected.
		out.AppendCharSequence("goog.loadModule(function(exports) {" + "'use strict';")
		closureBundler_EscapeMode_NORMAL.append(s, out)
		out.AppendCharSequence(
			"\n" + // terminate any trailing single line comment.
				";" + // terminate any trailing expression.
				"return exports;});\n")
		closureBundler_appendSourceUrl(out, closureBundler_EscapeMode_NORMAL, sourceUrl, hasSourceUrl)
	}
}

// Java: enum ClosureBundler.EscapeMode
type closureBundler_EscapeMode int32

const (
	closureBundler_EscapeMode_ESCAPED closureBundler_EscapeMode = iota
	closureBundler_EscapeMode_NORMAL
)

// String is EscapeMode.name().
func (m closureBundler_EscapeMode) String() string {
	return [...]string{"ESCAPED", "NORMAL"}[m]
}

// Java: EscapeMode.append(String, Appendable)
func (m closureBundler_EscapeMode) append(s string, out jutil.Appendable) {
	switch m {
	case closureBundler_EscapeMode_ESCAPED:
		jscomp.SourceCodeEscapers_AppendWithJavascriptEscaper(s, out)
	default: // NORMAL
		out.AppendCharSequence(s)
	}
}

// Java: ClosureBundler.appendSourceUrl(Appendable, EscapeMode, String)
func closureBundler_appendSourceUrl(out jutil.Appendable, mode closureBundler_EscapeMode, sourceUrl string, hasSourceUrl bool) {
	if !hasSourceUrl {
		return
	}
	toAppend := "\n//# sourceURL=" + sourceUrl + "\n"
	// Don't go through #append. That method relies on #transformInput,
	// but source URLs generally aren't valid JS inputs.
	mode.append(toAppend, out)
}
