// Port of com/google/javascript/jscomp/bundle/CompilerBasedTransformer.java @ closure-compiler v20260415

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

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CompilerBasedTransformer is a source transformer base class. May also
// include a runtime that needs to be shipped with the final bundle.
//
// Java: abstract class bundle.CompilerBasedTransformer implements Source.Transformer
type CompilerBasedTransformer interface {
	Source_Transformer

	// GetRuntime returns nil for Optional.empty().
	//
	// Java: CompilerBasedTransformer.getRuntime()
	GetRuntime() *string

	// Java: CompilerBasedTransformer.getTranformationName()
	GetTranformationName() string

	GetCompilerBasedTransformerBase() *CompilerBasedTransformerBase
}

// CompilerBasedTransformerBase holds the fields of CompilerBasedTransformer
// and its implementation; subclasses embed it.
type CompilerBasedTransformerBase struct {
	self             CompilerBasedTransformer
	compilerSupplier CompilerBasedTransformer_CompilerSupplier
}

// InitCompilerBasedTransformerBase is the constructor as subclasses call
// it; self is the subclass instance.
//
// Java: CompilerBasedTransformer(CompilerSupplier)
func (t *CompilerBasedTransformerBase) InitCompilerBasedTransformerBase(self CompilerBasedTransformer, compilerSupplier CompilerBasedTransformer_CompilerSupplier) {
	t.self = self
	t.compilerSupplier = jutil.CheckNotNull(compilerSupplier)
}

// GetCompilerBasedTransformerBase returns the CompilerBasedTransformer
// fields.
func (t *CompilerBasedTransformerBase) GetCompilerBasedTransformerBase() *CompilerBasedTransformerBase {
	return t
}

// Java: CompilerBasedTransformer.transform(Source)
func (t *CompilerBasedTransformerBase) Transform(input *Source) *Source {
	result := t.compilerSupplier.Compile(input.Path(), input.Code())
	if len(result.errors) != 0 {
		// TODO(sdh): how to handle this?  Currently we throw an ISE with the
		// message, but this may not be the most appropriate option.  It
		// might make sense to add console.log() statements to any JS that
		// comes out, particularly for warnings.
		formatter := ErrorFormat_SOURCELESS.ToFormatter(nil, false)
		message :=
			jutil.NewStringBuilder().Append(t.self.GetTranformationName()).Append(" failed.\n")
		for _, error := range result.errors {
			message.Append(formatter.FormatError(error))
		}
		panic(jutil.NewIllegalStateExceptionMsg(message.String()))
	}
	if !result.transformed {
		return input
	}
	builder := input.ToBuilder().
		SetCode(result.source).
		SetSourceMap(result.sourceMap)
	if runtime := t.self.GetRuntime(); runtime != nil {
		builder.AddRuntime(t.compilerSupplier.Runtime(*runtime))
	}
	return builder.Build()
}

// CompilerBasedTransformer_CompilerSupplier wraps the Compiler into a more
// relevant interface, making it easy to test the CompilerBasedTransformer
// without depending on implementation details of the Compiler itself. Also
// works around the fact that the Compiler is not thread-safe (since we may
// do multiple transpiles concurrently), so we supply a fresh instance each
// time when we're in single-file mode.
//
// Java: static class CompilerBasedTransformer.CompilerSupplier
type CompilerBasedTransformer_CompilerSupplier interface {
	// Java: CompilerSupplier.compile(Path, String)
	Compile(path *jutil.Path, code string) *CompilerBasedTransformer_CompileResult
	// Java: CompilerSupplier.transformed(Result)
	Transformed(result *Result) bool
	// Java: CompilerSupplier.runtime(String)
	Runtime(library string) string
	// Java: CompilerSupplier.compiler()
	Compiler() *Compiler
	// Java: CompilerSupplier.options()
	Options() *CompilerOptions
	// Java: CompilerSupplier.setOptions(CompilerOptions)
	SetOptions(options *CompilerOptions)

	GetCompilerBasedTransformer_CompilerSupplierBase() *CompilerBasedTransformer_CompilerSupplierBase
}

// CompilerBasedTransformer_CompilerSupplierBase holds the implementation of
// CompilerSupplier; subclasses embed it.
type CompilerBasedTransformer_CompilerSupplierBase struct {
	self CompilerBasedTransformer_CompilerSupplier
}

// Java: CompilerSupplier()
func NewCompilerBasedTransformer_CompilerSupplier() *CompilerBasedTransformer_CompilerSupplierBase {
	s := &CompilerBasedTransformer_CompilerSupplierBase{}
	s.InitCompilerBasedTransformer_CompilerSupplierBase(s)
	return s
}

// InitCompilerBasedTransformer_CompilerSupplierBase is the constructor as
// subclasses call it; self is the subclass instance.
//
// Java: CompilerSupplier()
func (s *CompilerBasedTransformer_CompilerSupplierBase) InitCompilerBasedTransformer_CompilerSupplierBase(self CompilerBasedTransformer_CompilerSupplier) {
	s.self = self
}

// GetCompilerBasedTransformer_CompilerSupplierBase returns the
// CompilerSupplier fields.
func (s *CompilerBasedTransformer_CompilerSupplierBase) GetCompilerBasedTransformer_CompilerSupplierBase() *CompilerBasedTransformer_CompilerSupplierBase {
	return s
}

// Java: CompilerSupplier.compile(Path, String)
func (s *CompilerBasedTransformer_CompilerSupplierBase) Compile(path *jutil.Path, code string) *CompilerBasedTransformer_CompileResult {
	compiler := s.self.Compiler()
	result :=
		compiler.CompileOfSourceFiles(CompilerBasedTransformer_CompilerSupplier_EXTERNS, SourceFile_FromCode(path.String(), code), s.self.Options())
	source := compiler.ToSource()
	var sourceMap strings.Builder
	if result.sourceMap != nil {
		jutil.Try(func() {
			name := path.String()
			result.sourceMap.AppendTo(&sourceMap, &name)
		}, func(e jutil.IIOException) {
			// impossible, and not a big deal even if it did happen.
		})
	}
	transformed := s.self.Transformed(result)
	sm := ""
	if transformed {
		sm = sourceMap.String()
	}
	return NewCompilerBasedTransformer_CompileResult(
		source, result.errors, transformed, sm)
}

// Java: CompilerSupplier.transformed(Result)
func (s *CompilerBasedTransformer_CompilerSupplierBase) Transformed(result *Result) bool {
	return result.transpiledFiles
}

// Java: CompilerSupplier.runtime(String)
func (s *CompilerBasedTransformer_CompilerSupplierBase) Runtime(library string) string {
	compiler := s.self.Compiler()
	options := s.self.Options()
	options.SetForceLibraryInjection([]string{library})
	compiler.CompileOfSourceFiles(CompilerBasedTransformer_CompilerSupplier_EXTERNS, CompilerBasedTransformer_CompilerSupplier_EMPTY, options)
	return compiler.ToSource()
}

// Java: CompilerSupplier.compiler()
func (s *CompilerBasedTransformer_CompilerSupplierBase) Compiler() *Compiler {
	return NewCompiler()
}

// Java: CompilerSupplier.options()
func (s *CompilerBasedTransformer_CompilerSupplierBase) Options() *CompilerOptions {
	options := NewCompilerOptions()
	s.self.SetOptions(options)
	return options
}

// Java: CompilerSupplier.setOptions(CompilerOptions)
func (s *CompilerBasedTransformer_CompilerSupplierBase) SetOptions(options *CompilerOptions) {
	// TODO(sdh): It would be nice to allow people to output code in
	// strict mode.  But currently we swallow all the input language
	// strictness checks, and there are various tests that are never
	// compiled and so are broken when we output 'use strict'.  We
	// could consider adding some sort of logging/warning/error in
	// cases where the input was not strict, though there could still
	// be semantic differences even if syntax is strict.  Possibly
	// the first step would be to allow the option of outputting strict
	// and then change the default and see what breaks.  b/33005948
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	options.SetQuoteKeywordProperties(true)
	options.SetSkipNonTranspilationPasses(true)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	options.SetPropertyRenaming(PropertyRenamingPolicy_OFF)
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	options.SetPrettyPrint(true)
	devNull := "/dev/null"
	options.SetSourceMapOutputPath(&devNull)
	options.SetSourceMapIncludeSourcesContent(true)
	options.SetWarningLevel(CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS, CheckLevel_OFF)
}

// Java: CompilerSupplier.EXTERNS
var CompilerBasedTransformer_CompilerSupplier_EXTERNS = SourceFile_FromCode("externs.js", "function Symbol() {}")

// Java: CompilerSupplier.EMPTY
var CompilerBasedTransformer_CompilerSupplier_EMPTY = SourceFile_FromCode("empty.js", "")

// Java: CompilerSupplier.ES5_WARNINGS
var CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS = NewDiagnosticGroup(DiagnosticType_Error("JSC_CANNOT_CONVERT", ""))

// CompilerBasedTransformer_CompileResult is the source together with the
// additional compilation results.
//
// Java: static class CompilerBasedTransformer.CompileResult
type CompilerBasedTransformer_CompileResult struct {
	source      string
	errors      []*JSError
	transformed bool
	sourceMap   string
}

// Java: CompileResult(String, ImmutableList<JSError>, boolean, String)
func NewCompilerBasedTransformer_CompileResult(
	source string, errors []*JSError, transformed bool, sourceMap string) *CompilerBasedTransformer_CompileResult {
	r := &CompilerBasedTransformer_CompileResult{}
	r.source = source
	r.errors = errors // checkNotNull: a Go slice is never null (nil is an empty list)
	r.transformed = transformed
	r.sourceMap = sourceMap
	return r
}
