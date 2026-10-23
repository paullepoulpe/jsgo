// Port of com/google/javascript/jscomp/deps/JsFileRegexParser.java @ closure-compiler v20261015

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
	JsFileLineParserBase

	// Previous lines stored for multi-line googMatcher application.
	googMatcherBuffer string

	// The info for the file we are currently parsing.
	provides []string

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

var _ JsFileLineParser = (*JsFileRegexParser)(nil)

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
// Java:
//
//	"(?:^|;)(?:[\\p{L}\\p{Nl}\\p{Nd}$_,:{}\\s]+=)?\\s*"
//	    + "goog\\.(?<func>provide|module|require|requireType|addDependency|declareModuleId)"
//	    + "\\s*\\(\\s*(?<args>.*?)\\s*\\)"
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

// ES6_MODULE_PATTERN is the pattern for matching import ... from
// './path/to/file'.
//
// Unlike the goog.require() pattern above, this pattern does not allow
// multiple statements per line. The import/export must be at the beginning
// of the line to match.
//
// Java:
//
//	// Require the import/export to be at the beginning of the line
//	"^"
//	    // Either an import or export, but we don't care which, followed by at least one space
//	    + "(?:import|export)\\b\\s*"
//	    // Skip any identifier chars, as well as star, comma, braces, and spaces
//	    // This should match, e.g., "* as foo from ", or "Foo, {Bar as Baz} from ".
//	    // The 'from' keyword is required except in the case of "import '...';",
//	    // where there's nothing between 'import' and the module key string literal.
//	    + "(?:[a-zA-Z0-9$_*,{}\\s]+\\bfrom\\s*|)"
//	    // Imports require a string literal at the end; it's optional for exports
//	    // (e.g. "export * from './other';", which is effectively also an import).
//	    // This optionally captures group #1, which is the imported module name.
//	    + "(?:['\"]([^'\"]+)['\"])?"
//	    // Finally, this should be the entire statement, so ensure there's a semicolon.
//	    + "\\s*;"
//
// Here: the part after "^(?:import|export)\\b", anchored at offset 6.
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

// ES6_EXPORT_PATTERN is the pattern for 'export' keyword, e.g. "export
// default class ..." or "export {blah}". The '\b' ensures we don't also match
// "exports = ...", which is not an ES6 module.
//
// Java: Pattern.compile("^export\\b"), used with lookingAt().
func jsFileRegexParser_ES6_EXPORT_PATTERN_LookingAt(line string) bool {
	return strings.HasPrefix(line, "export") && jsFileRegexParser_isWordBoundaryAfterKeyword(line)
}

// Line in comment indicating that the file is Closure's base.js.
const jsFileRegexParser_PROVIDES_GOOG_COMMENT = "@provideGoog"

// Line in comment indicating that the file is an extern.
const jsFileRegexParser_EXTERNS_COMMENT = "@externs"

// Line in comment indicating that the file should not be touched by JSCompiler.
const jsFileRegexParser_NOCOMPILE_COMMENT = "@nocompile"

// The start of a bundled goog.module, i.e. one that is wrapped in a goog.loadModule call
const jsFileRegexParser_BUNDLED_GOOG_MODULE_START = "goog.loadModule(function("

// jsFileRegexParser_ModuleType is JsFileRegexParser.ModuleType.
type jsFileRegexParser_ModuleType int32

const (
	jsFileRegexParser_ModuleType_NON_MODULE jsFileRegexParser_ModuleType = iota
	jsFileRegexParser_ModuleType_GOOG_MODULE
	jsFileRegexParser_ModuleType_GOOG_PROVIDE
	jsFileRegexParser_ModuleType_ES6_MODULE
)

// NewJsFileRegexParser is the constructor; errorManager handles parse
// errors.
//
// Java: JsFileRegexParser(ErrorManager)
func NewJsFileRegexParser(errorManager ErrorManager) *JsFileRegexParser {
	p := &JsFileRegexParser{
		loader: ModuleLoader_EMPTY,
	}
	p.InitJsFileLineParserBase(p, errorManager)
	return p
}

// SetIncludeGoogBase sets whether we should create implicit provides and
// requires of the root namespace.
//
// When generating deps files, you do not want this behavior. Deps files
// need base.js to run anyway, so they don't need information about it.
//
// When generating abstract build graphs, you probably do want this
// behavior. It will create an implicit dependency of all files with
// provides/requires on base.js.
//
// It returns this for easy chaining.
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
// It returns this for easy chaining.
//
// Java: JsFileRegexParser.setModuleLoader(ModuleLoader)
func (p *JsFileRegexParser) SetModuleLoader(loader *ModuleLoader) *JsFileRegexParser {
	p.loader = loader
	return p
}

// ParseFile parses the given file and returns the dependency information
// that it contained.
//
// filePath is the path to the file to parse; closureRelativePath the path
// of the file relative to closure; fileContents the contents to parse. The
// result is a DependencyInfo containing all provides/requires found in the
// file.
//
// Java: JsFileRegexParser.parseFile(String, String, String)
func (p *JsFileRegexParser) ParseFile(filePath string, closureRelativePath string, fileContents string) DependencyInfo {
	return p.parseReader(filePath, closureRelativePath, jutil.NewStringReader(fileContents))
}

// Java: JsFileRegexParser.parseReader(String, String, Reader)
func (p *JsFileRegexParser) parseReader(filePath string, closureRelativePath string, fileContents jutil.JavaReader) DependencyInfo {
	p.provides = []string{}
	p.requires = []*DependencyInfo_Require{}
	p.typeRequires = []string{}
	p.fileHasProvidesOrRequires = false
	p.hasExternsAnnotation = false
	p.hasNoCompileAnnotation = false
	p.file = p.loader.Resolve(filePath)
	p.moduleType = jsFileRegexParser_ModuleType_NON_MODULE
	p.seenLoadModule = false
	p.googMatcherBuffer = ""

	if jsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		jsFileRegexParser_logger.Log(jutil.Level_FINE, "Parsing Source: "+filePath)
	}
	p.DoParse(filePath, fileContents)

	if p.moduleType == jsFileRegexParser_ModuleType_ES6_MODULE {
		p.provides = append(p.provides, p.file.ToModuleName())
	}

	loadFlags := jutil.NewLinkedHashMap[string, string]()
	switch p.moduleType {
	case jsFileRegexParser_ModuleType_GOOG_MODULE:
		loadFlags.Put("module", "goog")
	case jsFileRegexParser_ModuleType_ES6_MODULE:
		loadFlags.Put("module", "es6")
	default:
		// Nothing to do here.
	}

	dependencyInfo := SimpleDependencyInfo_Builder_(closureRelativePath, filePath).
		SetProvides(p.provides).
		SetRequires(p.requires).
		SetTypeRequires(p.typeRequires).
		SetLoadFlags(loadFlags).
		SetHasExternsAnnotation(p.hasExternsAnnotation).
		SetHasNoCompileAnnotation(p.hasNoCompileAnnotation).
		Build()
	if jsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		jsFileRegexParser_logger.Log(jutil.Level_FINE, "DepInfo: "+dependencyInfo.String())
	}
	return dependencyInfo
}

// setModuleType returns true if the moduleType is successfully set and
// false otherwise (e.g. goog.provide, goog.module conflict).
//
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
				p.file.String()))

		return false
	}

	p.moduleType = type_
	return true
}

// Java: JsFileRegexParser.parseJsDocCommentLine(String)
func (p *JsFileRegexParser) ParseJsDocCommentLine(line string) bool {
	// Since these checks are all mutually exclusive, bail out after we see the first one.
	// @providesGoog should only ever be in one file (namely base.js),
	// and @nocompile means the file is thrown out entirely,
	// so it should never exist with @externs.
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

// applyGoogMatcher panics with *JsFileLineParser_ParseException.
//
// Java: JsFileRegexParser.applyGoogMatcher(String)
func (p *JsFileRegexParser) applyGoogMatcher(line string) bool {
	lineHasProvidesOrRequires := false
	// Iterate over the provides/requires.
	for _, m := range jsFileRegexParser_GOOG_PROVIDE_REQUIRE_PATTERN.FindAllStringSubmatchIndex(line, -1) {
		lineHasProvidesOrRequires = true

		if p.includeGoogBase && !p.fileHasProvidesOrRequires {
			p.fileHasProvidesOrRequires = true
			p.requires = append(p.requires, DependencyInfo_Require_BASE)
		}

		// See if it's a require or provide.
		methodName := line[m[2*jsFileRegexParser_googFuncGroup]:m[2*jsFileRegexParser_googFuncGroup+1]]
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
			// Parse the param.
			//
			// Go-only: performance — clone the name: a substring would keep
			// the whole file text alive as long as the DependencyInfo.
			arg := strings.Clone(p.ParseJsString(line[m[2*jsFileRegexParser_googArgsGroup]:m[2*jsFileRegexParser_googArgsGroup+1]]))
			// Add the dependency.
			if isRequire {
				if methodName == "requireType" {
					p.typeRequires = append(p.typeRequires, arg)
				} else if arg != "goog" {
					// goog is always implicit.
					require := DependencyInfo_Require_GoogRequireSymbol(arg)
					p.requires = append(p.requires, require)
				}
			} else {
				p.provides = append(p.provides, arg)
			}
		}
	}
	return lineHasProvidesOrRequires
}

// ParseLine parses a line of JavaScript, extracting goog.provide and
// goog.require information.
//
// Java: JsFileRegexParser.parseLine(String)
func (p *JsFileRegexParser) ParseLine(line string) bool {
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
		// Retry with the previous line.
		lineHasProvidesOrRequires = p.applyGoogMatcher(p.googMatcherBuffer + line)
	}
	if !lineHasProvidesOrRequires && lineHasProvidesOrRequiresWords {
		// Try without the previous line.
		lineHasProvidesOrRequires = p.applyGoogMatcher(line)
	}
	if !lineHasProvidesOrRequires && lineHasProvidesOrRequiresWords {
		// We haven't found any matches in this line and will retry on the next line.
		// We support only 2-lines retries for efficiency; otherwise if there were no matches, we
		// could end up applying the regexp to the whole file. It covers majority of multi-line
		// cases (according to the JS formatting rules).
		p.googMatcherBuffer = line
	} else {
		p.googMatcherBuffer = ""
	}

	if strings.HasPrefix(line, "import") || strings.HasPrefix(line, "export") {
		// "^(?:import|export)\\b" can only match at the start of the line, so
		// the find() loop runs at most once and es6Matcher.start() is 0.
		if jsFileRegexParser_isWordBoundaryAfterKeyword(line) {
			if m := jsFileRegexParser_ES6_MODULE_PATTERN_REST.FindStringSubmatchIndex(line[6:]); m != nil {
				p.setModuleType(jsFileRegexParser_ModuleType_ES6_MODULE)
				lineHasProvidesOrRequires = true

				if m[2] >= 0 {
					// Go-only: performance — cloned, as in applyGoogMatcher.
					arg := strings.Clone(line[6+m[2] : 6+m[3]])
					if strings.HasPrefix(arg, "goog:") {
						// cut off the "goog:" prefix
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

	return !p.shortcutMode ||
		lineHasProvidesOrRequires ||
		p.googMatcherBuffer != "" ||
		jsFileRegexParser_isAllWhitespace(line) ||
		!strings.Contains(line, ";") ||
		strings.Contains(line, "goog.setTestOnly") ||
		strings.Contains(line, "goog.module.declareLegacyNamespace")
}

// jsFileRegexParser_isAllWhitespace is CharMatcher.whitespace().matchesAllOf(line).
func jsFileRegexParser_isAllWhitespace(line string) bool {
	for _, c := range jutil.StringToUTF16(line) {
		if !jutil.CharMatcher_Whitespace_Matches(c) {
			return false
		}
	}
	return true
}

// Java: JsFileRegexParser.isSupported()
func JsFileRegexParser_IsSupported() bool {
	return true
}
