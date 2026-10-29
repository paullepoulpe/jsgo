// Port of com/google/javascript/jscomp/deps/JsFileFullParser.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsFileFullParser is a parser that extracts dependency information from a
// .js file, including goog.require, goog.provide, goog.module, import/export
// statements, and JSDoc annotations related to dependency management.
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
	isConfig          bool
	isExterns         bool
	provideGoog       bool
	testonly          bool
	moduleType        JsFileFullParser_FileInfo_ModuleType
	isLegacyNamespace bool

	delcalls     []string
	deltemplates []string
	// Import order matters!
	importedModules []string
	dynamicRequires []string
	readToggles     []string
	modName         []string
	mods            []string

	// Note: multiple copies doesn't make much sense, but we report
	// each copy so that calling code can choose how to handle it
	provides      []string
	requires      []string
	typeRequires  []string
	maybeRequires []string
	requiresCss   []string
	visibility    []string
}

func jsFileFullParser_appendUnique(dst []string, elem string) []string {
	if slices.Contains(dst, elem) {
		return dst
	}
	return append(dst, elem)
}

func (info *JsFileFullParser_FileInfo) sortMultisets() {
	slices.SortFunc(info.delcalls, jutil.StringCompareTo)
	info.delcalls = slices.Compact(info.delcalls)
	slices.SortFunc(info.deltemplates, jutil.StringCompareTo)
	info.deltemplates = slices.Compact(info.deltemplates)
	slices.SortFunc(info.provides, jutil.StringCompareTo)
	slices.SortFunc(info.requires, jutil.StringCompareTo)
	slices.SortFunc(info.typeRequires, jutil.StringCompareTo)
	slices.SortFunc(info.maybeRequires, jutil.StringCompareTo)
	slices.SortFunc(info.requiresCss, jutil.StringCompareTo)
	slices.SortFunc(info.visibility, jutil.StringCompareTo)
}

// Java: JsFileFullParser.FileInfo()
func NewJsFileFullParser_FileInfo() *JsFileFullParser_FileInfo {
	return &JsFileFullParser_FileInfo{
		moduleType: JsFileFullParser_FileInfo_ModuleType_UNKNOWN,
		modName:    []string{},
		mods:       []string{},
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

// jsFileFullParser_CommentAnnotation represents a single JSDoc annotation,
// with an optional argument.
//
// Java: static class JsFileFullParser.CommentAnnotation
type jsFileFullParser_CommentAnnotation struct {
	// Annotation name, e.g. "@fileoverview" or "@externs".
	name string
	// Annotation value: either the bare identifier immediately after the
	// annotation, or else string in braces.
	value string
}

// Java: CommentAnnotation(String, String)
func newJsFileFullParser_CommentAnnotation(name string, value string) *jsFileFullParser_CommentAnnotation {
	a := &jsFileFullParser_CommentAnnotation{}
	a.name = name
	a.value = value
	return a
}

// jsFileFullParser_CommentAnnotation_parse returns all the annotations in a
// given comment string.
//
// Java: CommentAnnotation.parse(String)
func jsFileFullParser_CommentAnnotation_parse(comment string) []*jsFileFullParser_CommentAnnotation {
	// TODO(b/139159612): This is reinventing a large part of JSDocInfoParser.
	// We should try to consolidate as much as possible.  This requires
	// several steps:
	//   1. Make all the annotations we look for first-class in JSDocInfo
	//   2. Support custom annotations (may already be done?)
	//   3. Fix up existing code so that all these annotations are in @fileoverview
	//   4. Change this code to simply inspect the script's JSDocInfo instead of re-parsing
	out := []*jsFileFullParser_CommentAnnotation{}
	matcher := jsFileFullParser_CommentAnnotation_ANNOTATION_RE.Matcher(comment)
	for matcher.Find() {
		name := *matcher.Group(jsFileFullParser_CommentAnnotation_ANNOTATION_NAME_GROUP)
		value := ""
		if v := matcher.Group(jsFileFullParser_CommentAnnotation_ANNOTATION_VALUE_GROUP); v != nil {
			value = *v
		}
		out = append(out, newJsFileFullParser_CommentAnnotation(name, value))
	}
	return out
}

// Regex for a JSDoc annotation with an `@name` and an optional
// brace-delimited `{value}`. The `@` should not match the middle of a word.
var jsFileFullParser_CommentAnnotation_ANNOTATION_RE = jutil.Pattern_Compile(
	"(?:[^a-zA-Z0-9_$]|^)(@[a-zA-Z]+)(?:\\s*\\{\\s*([^}\\t\\n\\v\\f\\r ]+)\\s*\\})?")

const jsFileFullParser_CommentAnnotation_ANNOTATION_NAME_GROUP = 1
const jsFileFullParser_CommentAnnotation_ANNOTATION_VALUE_GROUP = 2

// JsFileFullParser_Parse parses a JavaScript file for dependencies and
// annotations. It returns empty info if a syntax error was encountered.
//
// Java: JsFileFullParser.parse(String, String, Reporter)
func JsFileFullParser_Parse(code string, filename string, reporter JsFileFullParser_Reporter) *JsFileFullParser_FileInfo {
	errorReporter := newJsFileFullParser_DelegatingReporter(reporter)
	em := &jsFileFullParser_ErrorManager{errorReporter: errorReporter}
	em.InitBasicErrorManagerBase(em)
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

	parsed.Ast.SetInputId(rhino.NewInputId(filename))

	jsFileFullParser_parseTopLevelJsDocComments(parsed.TopLevelStatementRanges, parsed.Comments, info)
	gatherModuleMetadata :=
		NewGatherModuleMetadata(
			compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER)
	gatherModuleMetadata.Process(rhino.NewNode(rhino.Token_ROOT), parsed.Ast)
	compiler.GenerateReport()
	byPath := compiler.GetModuleMetadataMap().ModulesByPathMap()
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
	BasicErrorManagerBase
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

// Java: BasicErrorManager.printSummary()
func (m *jsFileFullParser_ErrorManager) PrintSummary() {}

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

// jsFileFullParser_parseTopLevelJsDocComments parses all top-level JSDoc
// comments (comments before or after top-level statements, or all comments
// if there are no statements), ignoring comments nested inside top-level
// statement bodies (functions, classes, methods, blocks).
//
// Uses the parser's statement ranges rather than Node.getSourceOffset, which
// is lossy: Node wraps line numbers past 2^20, clamps columns past
// Node.MAX_COLUMN_NUMBER, and resolves lines via SourceFile, which disagrees
// with the parser on line terminators other than \n.
//
// Java: JsFileFullParser.parseTopLevelJsDocComments(List<SourceRange>, List<Comment>, FileInfo)
func jsFileFullParser_parseTopLevelJsDocComments(
	statementRanges []parser.SourceRange, comments []parser.Comment, info *JsFileFullParser_FileInfo) {
	i := 0
	for idx := range comments {
		comment := &comments[idx]
		if comment.Type != parser.Comment_Type_JSDOC {
			continue
		}
		commentStart := comment.Location.Start.Offset
		commentEnd := comment.Location.End.Offset
		for i < len(statementRanges) && statementRanges[i].End.Offset <= commentStart {
			i++
		}
		if i < len(statementRanges) &&
			commentStart >= statementRanges[i].Start.Offset &&
			commentEnd <= statementRanges[i].End.Offset {
			continue
		}
		jsFileFullParser_parseComment(comment, info)
	}
}

// jsFileFullParser_parseComment mutates info with information from the given
// comment.
//
// Java: JsFileFullParser.parseComment(Comment, FileInfo)
func jsFileFullParser_parseComment(comment *parser.Comment, info *JsFileFullParser_FileInfo) {
	for _, annotation := range jsFileFullParser_CommentAnnotation_parse(comment.Value) {
		switch annotation.name {
		case "@fileoverview", "@author", "@see", "@link":
		case "@mods":
			if annotation.value != "" {
				info.mods = append(info.mods, annotation.value)
			}
		case "@visibility":
			if annotation.value != "" {
				info.visibility = append(info.visibility, annotation.value)
			}
		case "@modName":
			if annotation.value != "" {
				info.modName = append(info.modName, annotation.value)
			}
		case "@config":
			info.isConfig = true
		case "@provideGoog":
			info.provideGoog = true
		case "@requirecss":
			if annotation.value != "" {
				info.requiresCss = append(info.requiresCss, annotation.value)
			}
		case "@deltemplate", "@hassoydeltemplate":
			// TODO(b/210468818): Remove legacy @hassoydeltemplate annotation
			// once Soy gencode has been updated to use the new one.
			if annotation.value != "" {
				info.deltemplates = append(info.deltemplates, annotation.value)
			}
		case "@delcall", "@hassoydelcall":
			// TODO(b/210468818): Remove legacy @hassoydelcall annotation once
			// Soy gencode has been updated to use the new one.
			if annotation.value != "" {
				info.delcalls = append(info.delcalls, annotation.value)
			}
		case "@externs":
			info.isExterns = true
		}
	}
}

// JsFileFullParser_Reporter is the interface for reporting errors.
//
// Java: interface JsFileFullParser.Reporter
type JsFileFullParser_Reporter interface {
	// Java: Reporter.report(boolean, String, String, int, int)
	Report(fatal bool, message string, sourceName string, line int, lineOffset int)
}

// Java: static final class JsFileFullParser.DelegatingReporter implements ErrorReporter
type jsFileFullParser_DelegatingReporter struct {
	delegate  JsFileFullParser_Reporter
	seenError bool
}

// Java: DelegatingReporter(Reporter)
func newJsFileFullParser_DelegatingReporter(delegate JsFileFullParser_Reporter) *jsFileFullParser_DelegatingReporter {
	r := &jsFileFullParser_DelegatingReporter{}
	r.delegate = jutil.CheckNotNull(delegate)
	return r
}

// Java: DelegatingReporter.warning(String, String, int, int)
func (r *jsFileFullParser_DelegatingReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.delegate.Report(false, message, sourceName, line, lineOffset)
}

// Java: DelegatingReporter.error(String, String, int, int)
func (r *jsFileFullParser_DelegatingReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.seenError = true
	r.delegate.Report(true, message, sourceName, line, lineOffset)
}
