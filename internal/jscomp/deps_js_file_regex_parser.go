// Port of com/google/javascript/jscomp/deps/JsFileRegexParser.java @ closure-compiler v20261021

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
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JsFileRegexParser is a parser that can extract dependency information
// from a .js file, including goog.require, goog.provide, goog.module, import
// statements, and export statements.
type JsFileRegexParser struct {
	jsFileLineParser

	// Previous lines stored for multi-line googMatcher application.
	googMatcherBuffer string

	// The info for the file we are currently parsing.
	provides                  []string
	requires                  []*DependencyInfo_Require
	typeRequires              []string
	fileHasProvidesOrRequires bool
	hasExternsAnnotation      bool
	hasNoCompileAnnotation    bool
	loader                    *ModuleLoader
	file                      *ModuleLoader_ModulePath

	moduleType     jsFileRegexParser_ModuleType
	seenLoadModule bool

	// Whether to provide/require the root namespace.
	includeGoogBase bool
}

// Java: JsFileRegexParser.logger
var jsFileRegexParser_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.JsFileRegexParser")

// Go-only: language-adapter — Java regex. The three patterns below are the
// Java patterns rewritten in RE2 syntax with Java's meaning: \s is Java's
// [ \t\n\x0B\f\r] (RE2's \s lacks \x0B), `.` excludes Java's line
// terminators (\n, \r, \x{85}, \x{2028}, \x{2029}), and the named groups use
// RE2's (?P<name>...). None of them uses back-references or lookaround, and
// RE2's leftmost-first matching selects the same match and submatches as
// Java's backtracking engine. Java's \b (ASCII word characters, plus
// non-spacing marks that follow a letter or digit) differs from RE2's, so
// the boundary after "import"/"export" is checked by
// jsFileRegexParser_isWordBoundaryAfterKeyword; the \b before "from" is
// RE2's (the character before it is always an ASCII class member).

// GOOG_PROVIDE_REQUIRE_PATTERN is the pattern for matching goog.provide(*)
// and goog.require(*).
//
// TODO(sdh): this handles goog.loadModule(function(){"use strict";goog.module
// but fails to match without "use strict"; since we look for semicolon, not open brace.
// TODO(johnplaisted): Remove declareNamespace.
var jsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN = regexp.MustCompile(
	`(?:^|;)(?:[\p{L}\p{Nl}\p{Nd}$_,:{}\t\n\x0B\f\r ]+=)?[\t\n\x0B\f\r ]*` +
		`goog\.(?P<func>provide|module|require|requireType|addDependency|declareModuleId)` +
		`[\t\n\x0B\f\r ]*\([\t\n\x0B\f\r ]*(?P<args>[^\n\r\x{85}\x{2028}\x{2029}]*?)[\t\n\x0B\f\r ]*\)`)

var (
	jsFileRegexParser_googFuncGroup = jsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN.SubexpIndex("func")
	jsFileRegexParser_googArgsGroup = jsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN.SubexpIndex("args")
)

// ES6_MODULE_PATTERN_REST matches the part of ES6_MODULE_PATTERN after
// "^(?:import|export)\\b", anchored at offset 6.
var jsFileRegexParser_ES6_MODULE_PATTERN_REST = regexp.MustCompile(
	`^[\t\n\x0B\f\r ]*` +
		`(?:[a-zA-Z0-9$_*,{}\t\n\x0B\f\r ]+\bfrom[\t\n\x0B\f\r ]*|)` +
		`(?:['"]([^'"]+)['"])?` +
		`[\t\n\x0B\f\r ]*;`)

// jsFileRegexParser_isWordBoundaryAfterKeyword is Java's \b at offset 6 of
// a line starting with "import" or "export": the preceding 't' is a word
// character, so there is a boundary unless the next character is a word
// character ([A-Za-z0-9_]) or a non-spacing mark (which Java attaches to the
// preceding letter).
func jsFileRegexParser_isWordBoundaryAfterKeyword(line string) bool {
	if len(line) <= 6 {
		return true
	}
	c := line[6]
	if c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
		return false
	}
	if c >= utf8.RuneSelf {
		r, _ := utf8.DecodeRuneInString(line[6:])
		if r != utf8.RuneError &&
			jutil.CharacterGetTypeOfCodePoint(int(r)) == jutil.Character_NON_SPACING_MARK {
			return false
		}
	}
	return true
}

// Java: Pattern.compile("^export\\b"), used with lookingAt().
func jsFileRegexParser_ES6_EXPORT_PATTERN_LookingAt(line string) bool {
	return strings.HasPrefix(line, "export") && jsFileRegexParser_isWordBoundaryAfterKeyword(line)
}

const (
	// Line in comment indicating that the file is Closure's base.js.
	jsFileRegexParser_PROVIDES_GOOG_COMMENT = "@provideGoog"
	// Line in comment indicating that the file is an extern.
	jsFileRegexParser_EXTERNS_COMMENT = "@externs"
	// Line in comment indicating that the file should not be touched by JSCompiler.
	jsFileRegexParser_NOCOMPILE_COMMENT = "@nocompile"
	// The start of a bundled goog.module, i.e. one that is wrapped in a goog.loadModule call
	jsFileRegexParser_BUNDLED_GOOG_MODULE_START = "goog.loadModule(function("
)

type jsFileRegexParser_ModuleType int32

const (
	jsFileRegexParser_ModuleType_NON_MODULE jsFileRegexParser_ModuleType = iota
	jsFileRegexParser_ModuleType_GOOG_MODULE
	jsFileRegexParser_ModuleType_GOOG_PROVIDE
	jsFileRegexParser_ModuleType_ES6_MODULE
)

// Java: JsFileRegexParser(ErrorManager)
func NewJsFileRegexParser(errorManager ErrorManager) *JsFileRegexParser {
	return &JsFileRegexParser{
		jsFileLineParser: jsFileLineParser{errorManager: errorManager},
		loader:           ModuleLoader_EMPTY,
	}
}

// SetIncludeGoogBase sets whether we should create implicit provides and
// requires of the root namespace.
//
// Java: JsFileRegexParser.setIncludeGoogBase(boolean)
func (p *JsFileRegexParser) SetIncludeGoogBase(include bool) *JsFileRegexParser {
	jutil.CheckState(JsFileRegexParser_IsSupported())
	p.includeGoogBase = include
	return p
}

// SetModuleLoader sets a list of "module root" URIs, which allow
// relativizing filenames for modules.
//
// Java: JsFileRegexParser.setModuleLoader(ModuleLoader)
func (p *JsFileRegexParser) SetModuleLoader(loader *ModuleLoader) *JsFileRegexParser {
	p.loader = loader
	return p
}

// ParseFile parses the given file and returns the dependency information
// that it contained.
//
// Java: JsFileRegexParser.parseFile(String, String, String)
func (p *JsFileRegexParser) ParseFile(filePath string, closureRelativePath string, fileContents string) DependencyInfo {
	p.resetForFile(filePath)
	if jsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		jsFileRegexParser_logger.Log(jutil.Level_FINE, "Parsing Source: "+filePath)
	}
	parseJsString(&p.jsFileLineParser, p, filePath, fileContents)
	return p.finishDependencyInfo(filePath, closureRelativePath)
}

func (p *JsFileRegexParser) resetForFile(filePath string) {
	p.provides = p.provides[:0]
	p.requires = p.requires[:0]
	p.typeRequires = p.typeRequires[:0]
	p.fileHasProvidesOrRequires = false
	p.hasExternsAnnotation = false
	p.hasNoCompileAnnotation = false
	p.file = p.loader.Resolve(filePath)
	p.moduleType = jsFileRegexParser_ModuleType_NON_MODULE
	p.seenLoadModule = false
	p.googMatcherBuffer = ""
}

func (p *JsFileRegexParser) finishDependencyInfo(filePath string, closureRelativePath string) DependencyInfo {
	if p.moduleType == jsFileRegexParser_ModuleType_ES6_MODULE {
		p.provides = append(p.provides, p.file.ToModuleName())
	}

	var loadFlags *jutil.LinkedHashMap[string, string]
	switch p.moduleType {
	case jsFileRegexParser_ModuleType_GOOG_MODULE:
		loadFlags = simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS
	case jsFileRegexParser_ModuleType_ES6_MODULE:
		loadFlags = simpleDependencyInfo_ES6_FLAGS
	default:
		loadFlags = simpleDependencyInfo_EMPTY_FLAGS
	}

	dependencyInfo := newSimpleDependencyInfoDirect(
		closureRelativePath,
		filePath,
		p.provides,
		p.requires,
		p.typeRequires,
		loadFlags,
		p.hasExternsAnnotation,
		p.hasNoCompileAnnotation,
	)
	if jsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		jsFileRegexParser_logger.Log(jutil.Level_FINE, "DepInfo: "+dependencyInfo.String())
	}
	return dependencyInfo
}

// Java: JsFileRegexParser.setModuleType(ModuleType)
func (p *JsFileRegexParser) setModuleType(type_ jsFileRegexParser_ModuleType) bool {
	provide := type_ == jsFileRegexParser_ModuleType_GOOG_PROVIDE || p.moduleType == jsFileRegexParser_ModuleType_GOOG_PROVIDE
	es6Module := type_ == jsFileRegexParser_ModuleType_ES6_MODULE || p.moduleType == jsFileRegexParser_ModuleType_ES6_MODULE
	googModule := type_ == jsFileRegexParser_ModuleType_GOOG_MODULE || p.moduleType == jsFileRegexParser_ModuleType_GOOG_MODULE

	if googModule && provide && p.seenLoadModule {
		// We have to assume this is a top level goog.provide and a wrapped goog.loadModule. We can't
		// correctly validate this with just regular expressions.
		p.moduleType = jsFileRegexParser_ModuleType_GOOG_PROVIDE
		return true
	}

	provideGoogModuleConflict := googModule && provide && !p.seenLoadModule
	provideEs6ModuleConflict := es6Module && provide
	// Don't allow nested goog modules in ES6 modules.
	googEs6ModuleConflict := (googModule || p.seenLoadModule) && es6Module

	if provideGoogModuleConflict || provideEs6ModuleConflict || googEs6ModuleConflict {
		// TODO(sdh): should this be an error?
		fileName := p.file.String()
		p.errorManager.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(
				&fileName,
				/* lineno= */ -1,
				/* charno= */ -1,
				ModuleLoader_MODULE_CONFLICT,
				fileName))

		return false
	}

	p.moduleType = type_
	return true
}

// Java: JsFileRegexParser.parseJsDocCommentLine(String)
func (p *JsFileRegexParser) parseJsDocCommentLine(line string) bool {
	// Since these checks are all mutually exclusive, bail out after we see the first one.
	if p.includeGoogBase && strings.Contains(line, jsFileRegexParser_PROVIDES_GOOG_COMMENT) {
		p.provides = append(p.provides, "goog")
		return false
	} else if strings.Contains(line, jsFileRegexParser_EXTERNS_COMMENT) {
		p.hasExternsAnnotation = true
		return false
	} else if strings.Contains(line, jsFileRegexParser_NOCOMPILE_COMMENT) {
		p.hasNoCompileAnnotation = true
		return false
	}
	return true
}

func jsFileRegexParser_tryFastGoogCall(line string) (methodName, args string, ok bool) {
	i := 0
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	if i+5 >= len(line) || line[i:i+5] != "goog." {
		return "", "", false
	}
	rest := line[i+5:]
	openIdx := strings.IndexByte(rest, '(')
	if openIdx <= 0 {
		return "", "", false
	}
	fnEnd := openIdx
	for fnEnd > 0 && (rest[fnEnd-1] == ' ' || rest[fnEnd-1] == '\t') {
		fnEnd--
	}
	fn := rest[:fnEnd]
	switch fn {
	case "provide", "module", "require", "requireType", "addDependency", "declareModuleId":
	default:
		return "", "", false
	}
	end := len(rest)
	for end > openIdx+1 && (rest[end-1] == ' ' || rest[end-1] == '\t') {
		end--
	}
	if end > openIdx+1 && rest[end-1] == ';' {
		end--
		for end > openIdx+1 && (rest[end-1] == ' ' || rest[end-1] == '\t') {
			end--
		}
	}
	if end <= openIdx || rest[end-1] != ')' {
		return "", "", false
	}
	aStart := openIdx + 1
	aEnd := end - 1
	for aStart < aEnd && jsFileLineParser_isJavaSpace(rest[aStart]) {
		aStart++
	}
	for aEnd > aStart && jsFileLineParser_isJavaSpace(rest[aEnd-1]) {
		aEnd--
	}
	rawArgs := rest[aStart:aEnd]
	for k := 0; k < len(rawArgs); k++ {
		c := rawArgs[k]
		if c == ')' || c == ';' || c == '\n' || c == '\r' || c >= 0x80 {
			return "", "", false
		}
	}
	return fn, rawArgs, true
}

func (p *JsFileRegexParser) recordGoogMatch(methodName, rawArgs string) string {
	if p.includeGoogBase && !p.fileHasProvidesOrRequires {
		p.fileHasProvidesOrRequires = true
		p.requires = append(p.requires, DependencyInfo_Require_BASE)
	}

	firstChar := methodName[0]
	isDeclareModuleNamespace := firstChar == 'd'
	isModule := !isDeclareModuleNamespace && firstChar == 'm'
	isProvide := firstChar == 'p'
	providesNamespace := isProvide || isModule || isDeclareModuleNamespace
	isRequire := firstChar == 'r'

	if isModule && !p.seenLoadModule {
		providesNamespace = p.setModuleType(jsFileRegexParser_ModuleType_GOOG_MODULE)
	}

	if isProvide {
		providesNamespace = p.setModuleType(jsFileRegexParser_ModuleType_GOOG_PROVIDE)
	}

	if providesNamespace || isRequire {
		parsed, errMsg, ok := p.parseJsString(rawArgs)
		if !ok {
			return errMsg
		}
		// Go-only: performance — clone the name: a substring would keep
		// the whole file text alive as long as the DependencyInfo.
		arg := strings.Clone(parsed)
		if isRequire {
			if methodName == "requireType" {
				p.typeRequires = append(p.typeRequires, arg)
			} else if arg != "goog" {
				p.requires = append(p.requires, DependencyInfo_Require_GoogRequireSymbol(arg))
			}
		} else {
			p.provides = append(p.provides, arg)
		}
	}
	return ""
}

// Java: JsFileRegexParser.applyGoogMatcher(String)
func (p *JsFileRegexParser) applyGoogMatcher(line string) (bool, string) {
	if fn, rawArgs, ok := jsFileRegexParser_tryFastGoogCall(line); ok {
		if errMsg := p.recordGoogMatch(fn, rawArgs); errMsg != "" {
			return false, errMsg
		}
		return true, ""
	}
	if !strings.Contains(line, "goog.") {
		return false, ""
	}
	lineHasProvidesOrRequires := false
	for _, m := range jsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN.FindAllStringSubmatchIndex(line, -1) {
		lineHasProvidesOrRequires = true
		methodName := line[m[2*jsFileRegexParser_googFuncGroup]:m[2*jsFileRegexParser_googFuncGroup+1]]
		rawArgs := line[m[2*jsFileRegexParser_googArgsGroup]:m[2*jsFileRegexParser_googArgsGroup+1]]
		if errMsg := p.recordGoogMatch(methodName, rawArgs); errMsg != "" {
			return false, errMsg
		}
	}
	return lineHasProvidesOrRequires, ""
}

// Java: JsFileRegexParser.parseLine(String)
func (p *JsFileRegexParser) parseLine(line string) (bool, string) {
	lineHasProvidesOrRequires := false

	if strings.HasPrefix(line, jsFileRegexParser_BUNDLED_GOOG_MODULE_START) {
		p.seenLoadModule = true
	}

	// Quick check that will catch most cases. This is a performance win for teams with a lot of JS.
	lineHasProvidesOrRequiresWords :=
		strings.Contains(line, "provide") ||
			strings.Contains(line, "require") ||
			strings.Contains(line, "module") ||
			strings.Contains(line, "addDependency") ||
			strings.Contains(line, "declareModuleId")

	if p.googMatcherBuffer != "" {
		var errMsg string
		lineHasProvidesOrRequires, errMsg = p.applyGoogMatcher(p.googMatcherBuffer + line)
		if errMsg != "" {
			return false, errMsg
		}
	}
	if !lineHasProvidesOrRequires && lineHasProvidesOrRequiresWords {
		var errMsg string
		lineHasProvidesOrRequires, errMsg = p.applyGoogMatcher(line)
		if errMsg != "" {
			return false, errMsg
		}
	}
	if !lineHasProvidesOrRequires && lineHasProvidesOrRequiresWords {
		// Retry on the next line (at most 2-line retries).
		p.googMatcherBuffer = line
	} else {
		p.googMatcherBuffer = ""
	}

	if strings.HasPrefix(line, "import") || strings.HasPrefix(line, "export") {
		if jsFileRegexParser_isWordBoundaryAfterKeyword(line) {
			if m := jsFileRegexParser_ES6_MODULE_PATTERN_REST.FindStringSubmatchIndex(line[6:]); m != nil {
				p.setModuleType(jsFileRegexParser_ModuleType_ES6_MODULE)
				lineHasProvidesOrRequires = true

				if m[2] >= 0 {
					// Go-only: performance — cloned, as in applyGoogMatcher.
					arg := strings.Clone(line[6+m[2] : 6+m[3]])
					if strings.HasPrefix(arg, "goog:") {
						p.requires = append(p.requires, DependencyInfo_Require_GoogRequireSymbol(arg[5:]))
					} else {
						path := p.file.ResolveJsModule(arg, p.filePath, p.lineNum, 0 /* es6Matcher.start() */)
						if path == nil {
							path = p.file.ResolveModuleAsPath(arg)
						}
						p.requires = append(p.requires, DependencyInfo_Require_Es6Import(path.ToModuleName(), arg))
					}
				}
			}
		}

		// This check is only relevant for modules that don't import anything.
		if p.moduleType != jsFileRegexParser_ModuleType_ES6_MODULE && jsFileRegexParser_ES6_EXPORT_PATTERN_LookingAt(line) {
			p.setModuleType(jsFileRegexParser_ModuleType_ES6_MODULE)
		}
	}

	keepGoing := !p.shortcutMode ||
		lineHasProvidesOrRequires ||
		p.googMatcherBuffer != "" ||
		jsFileRegexParser_isAllWhitespace(line) ||
		!strings.Contains(line, ";") ||
		strings.Contains(line, "goog.setTestOnly") ||
		strings.Contains(line, "goog.module.declareLegacyNamespace")
	return keepGoing, ""
}

// jsFileRegexParser_isAllWhitespace is CharMatcher.whitespace().matchesAllOf(line).
func jsFileRegexParser_isAllWhitespace(line string) bool {
	for i := 0; i < len(line); {
		c := line[i]
		if c < utf8.RuneSelf {
			if c != ' ' && (c < 0x09 || c > 0x0D) {
				return false
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(line[i:])
		if !jsFileLineParser_isCharMatcherSpace(r, size) {
			return false
		}
		i += size
	}
	return true
}

// Java: JsFileRegexParser.isSupported()
func JsFileRegexParser_IsSupported() bool {
	return true
}
