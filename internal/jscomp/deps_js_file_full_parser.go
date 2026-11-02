// Port of com/google/javascript/jscomp/deps/JsFileFullParser.java @ closure-compiler v20261028

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsFileFullParser is a parser that extracts dependency information from a
// .js file, including goog.require, goog.provide, goog.module, and
// import/export statements.
//
// Java: class deps.JsFileFullParser
type JsFileFullParser struct {
	// Go-only: language-adapter — distinct identity.
	_ byte
}

// JsFileFullParser_FileInfo is the dependency information contained in a .js
// source file.
//
// Java: static final class JsFileFullParser.FileInfo
type JsFileFullParser_FileInfo struct {
	goog              bool
	testonly          bool
	moduleType        JsFileFullParser_FileInfo_ModuleType
	isLegacyNamespace bool

	// Import order matters!
	importedModules []string
	dynamicRequires []string
	readToggles     []string

	// Note: multiple copies doesn't make much sense, but we report
	// each copy so that calling code can choose how to handle it
	provides      []string
	requires      []string
	typeRequires  []string
	maybeRequires []string

	parsed *parsing.ParserRunner_ParseResult
}

func jsFileFullParser_appendUnique(dst []string, elem string) []string {
	if slices.Contains(dst, elem) {
		return dst
	}
	return append(dst, elem)
}

func (info *JsFileFullParser_FileInfo) sortMultisets() {
	slices.SortFunc(info.provides, jutil.StringCompareTo)
	slices.SortFunc(info.requires, jutil.StringCompareTo)
	slices.SortFunc(info.typeRequires, jutil.StringCompareTo)
	slices.SortFunc(info.maybeRequires, jutil.StringCompareTo)
}

// Java: JsFileFullParser.FileInfo()
func NewJsFileFullParser_FileInfo() *JsFileFullParser_FileInfo {
	return &JsFileFullParser_FileInfo{
		moduleType: JsFileFullParser_FileInfo_ModuleType_UNKNOWN,
	}
}

// JsFileFullParser_FileInfo_ModuleType is the module system declared by the
// file, e.g. goog.provide/goog.module.
//
// Java: enum JsFileFullParser.FileInfo.ModuleType
type JsFileFullParser_FileInfo_ModuleType int32

const (
	JsFileFullParser_FileInfo_ModuleType_UNKNOWN JsFileFullParser_FileInfo_ModuleType = iota
	JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE
	JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE
	JsFileFullParser_FileInfo_ModuleType_ES_MODULE
)

var jsFileFullParser_FileInfo_ModuleType_names = [...]string{"UNKNOWN", "GOOG_PROVIDE", "GOOG_MODULE", "ES_MODULE"}

// String is ModuleType.name().
func (t JsFileFullParser_FileInfo_ModuleType) String() string {
	return jsFileFullParser_FileInfo_ModuleType_names[t]
}

// JsFileFullParser_Parse parses a JavaScript file for dependencies and
// annotations. It returns empty info if a syntax error was encountered.
//
// Java: JsFileFullParser.parse(String, String, Reporter)
func JsFileFullParser_Parse(code string, filename string, reporter JsFileFullParser_Reporter) *JsFileFullParser_FileInfo {
	errorReporter := newJsFileFullParser_DelegatingReporter(reporter)
	em := &jsFileFullParser_ErrorManager{errorReporter: errorReporter}
	em.initReportHooks(em.Println, nil)
	compiler := NewCompilerWithErrorManager(em)
	source := SourceFile_FromCode(filename, code)

	options := NewCompilerOptions()
	options.SetChecksOnly(true)

	compiler.Init([]*SourceFile{}, []*SourceFile{source}, options)

	config :=
		parsing.ParserRunner_CreateConfigWithModes(
			// TODO(sdh): ES8 STRICT, with a non-strict fallback - then give
			// warnings.
			parsing.Config_LanguageMode_ES_NEXT,
			parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE,
			parsing.Config_RunMode_STOP_AFTER_ERROR,
			/* extraAnnotationNames= */ jutil.NewLinkedHashSet[string](),
			/* parseInlineSourceMaps= */ true,
			parsing.Config_StrictMode_SLOPPY)
	info := NewJsFileFullParser_FileInfo()
	parsed := parsing.ParserRunner_Parse(source, code, config, errorReporter)
	if errorReporter.seenError {
		return info
	}
	info.parsed = parsed

	parsed.Ast.SetInputId(rhino.NewInputId(filename))

	gatherModuleMetadata :=
		NewGatherModuleMetadata(
			compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER)
	gatherModuleMetadata.Process(rhino.NewNode(rhino.Token_ROOT), parsed.Ast)
	compiler.GenerateReport()
	byPath := compiler.GetModuleMetadataMap().GetModulesByPath()
	if byPath.Size() != 1 {
		info.sortMultisets()
		return info // Avoid potential crashes due to assumptions of the code below being violated.
	}
	module := byPath.Entries()[0].Value
	switch module.ModuleType() {
	case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
		info.moduleType = JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE
	case ModuleMetadataMap_ModuleType_GOOG_MODULE:
		info.moduleType = JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE
	case ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		info.moduleType = JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE
		info.isLegacyNamespace = true
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		info.moduleType = JsFileFullParser_FileInfo_ModuleType_ES_MODULE
	case ModuleMetadataMap_ModuleType_COMMON_JS, ModuleMetadataMap_ModuleType_SCRIPT:
		// Treat these as unknown for now; we can extend the enum if we care
		// about these.
		info.moduleType = JsFileFullParser_FileInfo_ModuleType_UNKNOWN
	}
	info.goog = module.UsesClosure()
	jsFileFullParser_recordModuleMetadata(info, module)
	info.sortMultisets()
	return info
}

// jsFileFullParser_ErrorManager is the anonymous BasicErrorManager of parse:
// it forwards errors and warnings to the DelegatingReporter.
type jsFileFullParser_ErrorManager struct {
	SortingErrorManager
	errorReporter *jsFileFullParser_DelegatingReporter
}

// Java: BasicErrorManager.println(CheckLevel, JSError)
func (m *jsFileFullParser_ErrorManager) Println(level *CheckLevel, error *JSError) {
	if level != nil && *level == CheckLevel_ERROR {
		m.errorReporter.Error(
			error.Description(),
			jsFileFullParser_sourceName(error),
			error.GetLineNumber(),
			error.Charno())
	} else if level != nil && *level == CheckLevel_WARNING {
		m.errorReporter.Warning(
			error.Description(),
			jsFileFullParser_sourceName(error),
			error.GetLineNumber(),
			error.Charno())
	}
}

// jsFileFullParser_sourceName is error.sourceName() for an ErrorReporter.
//
// Go-only: language-adapter — rhino.ErrorReporter takes a Go string, so a
// null source name is passed as "".
func jsFileFullParser_sourceName(error *JSError) string {
	if s := error.SourceName(); s != nil {
		return *s
	}
	return ""
}

// Java: JsFileFullParser.recordModuleMetadata(FileInfo, ModuleMetadata)
func jsFileFullParser_recordModuleMetadata(info *JsFileFullParser_FileInfo, module *ModuleMetadataMap_ModuleMetadata) {
	for _, e := range module.Es6ImportSpecifiers().Entries() {
		info.importedModules = jsFileFullParser_appendUnique(info.importedModules, e.Key)
	}

	// If something doesn't have an external dependency on Closure, then it
	// does not have any externally required files or symbols to provide.
	// This is needed for bundles that contain base.js as well as other
	// files. These bundles should look like they do not require or provide
	// anything at all.
	if module.UsesClosure() {
		info.provides = slices.AppendSeq(info.provides, module.GoogNamespaces().All())
		info.requires = slices.AppendSeq(info.requires, module.StronglyRequiredGoogNamespaces().All())
		info.maybeRequires = slices.AppendSeq(info.maybeRequires, module.MaybeRequiredGoogNamespaces().All())
		for _, e := range module.DynamicallyRequiredGoogNamespaces().Entries() {
			info.dynamicRequires = jsFileFullParser_appendUnique(info.dynamicRequires, e.Key)
		}
		for _, e := range module.ReadToggles().Entries() {
			info.readToggles = jsFileFullParser_appendUnique(info.readToggles, e.Key)
		}
		info.typeRequires = slices.AppendSeq(info.typeRequires, module.WeaklyRequiredGoogNamespaces().All())
		info.testonly = module.IsTestOnly()
	}
	// Traverse any nested modules (goog.loadModule calls).
	for _, nested := range module.NestedModules() {
		jsFileFullParser_recordModuleMetadata(info, nested)
	}
}

// JsFileFullParser_Reporter is the callback function for reporting errors.
//
// Java: interface JsFileFullParser.Reporter
type JsFileFullParser_Reporter func(fatal bool, message string, sourceName string, line int, lineOffset int)

// Java: static final class JsFileFullParser.DelegatingReporter implements ErrorReporter
type jsFileFullParser_DelegatingReporter struct {
	delegate  JsFileFullParser_Reporter
	seenError bool
}

// Java: DelegatingReporter(Reporter)
func newJsFileFullParser_DelegatingReporter(delegate JsFileFullParser_Reporter) *jsFileFullParser_DelegatingReporter {
	if delegate == nil {
		panic(jutil.NewNullPointerException())
	}
	return &jsFileFullParser_DelegatingReporter{delegate: delegate}
}

// Java: DelegatingReporter.warning(String, String, int, int)
func (r *jsFileFullParser_DelegatingReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.delegate(false, message, sourceName, line, lineOffset)
}

// Java: DelegatingReporter.error(String, String, int, int)
func (r *jsFileFullParser_DelegatingReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.seenError = true
	r.delegate(true, message, sourceName, line, lineOffset)
}
