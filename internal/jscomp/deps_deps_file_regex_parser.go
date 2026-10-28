// Port of com/google/javascript/jscomp/deps/DepsFileRegexParser.java @ closure-compiler v20261021

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
	"os"
	"regexp"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DepsFileRegexParser is a parser that can extract dependency information
// from existing deps.js files.
//
// See //javascript/closure/deps.js for an example file.
type DepsFileRegexParser struct {
	jsFileLineParser

	// The dependency information extracted from the current file.
	depInfos []DependencyInfo

	// Translates paths in different build systems.
	pathTranslator func(string) string
}

var depsFileRegexParser_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.DepsFileRegexParser")

// Go-only: language-adapter — Java regex. \s is Java's [ \t\n\x0B\f\r] and `.`
// excludes Java's line terminators (\n, \r, \x{85}, \x{2028}, \x{2029}).

// Java: Pattern.compile("\\s*goog.addDependency\\((.*)\\);?\\s*").matcher("")
var depsFileRegexParser_depMatcher = regexp.MustCompile(
	`^[\t\n\x0B\f\r ]*goog[^\n\r\x{85}\x{2028}\x{2029}]addDependency\(([^\n\r\x{85}\x{2028}\x{2029}]*)\);?[\t\n\x0B\f\r ]*$`)

// Java:
//
//	Pattern.compile(
//	        "\\s*([^,]*), (\\[[^\\]]*\\]), (\\[[^\\]]*\\])"
//	            + "(?:, (true|false|\\{[^{}]*\\}))?\\s*")
//	    .matcher("")
var depsFileRegexParser_depArgsMatch = regexp.MustCompile(
	`^[\t\n\x0B\f\r ]*([^,]*), (\[[^\]]*\]), (\[[^\]]*\])` +
		`(?:, (true|false|\{[^{}]*\}))?[\t\n\x0B\f\r ]*$`)

// Java: DepsFileRegexParser(ErrorManager)
func NewDepsFileRegexParser(errorManager ErrorManager) *DepsFileRegexParser {
	return NewDepsFileRegexParserWithPathTranslator(func(s string) string { return s }, errorManager)
}

// Java: DepsFileRegexParser(Function<String, String>, ErrorManager)
func NewDepsFileRegexParserWithPathTranslator(pathTranslator func(string) string, errorManager ErrorManager) *DepsFileRegexParser {
	return &DepsFileRegexParser{
		jsFileLineParser: jsFileLineParser{errorManager: errorManager},
		pathTranslator:   pathTranslator,
	}
}

// ParseFile parses the given file and returns a list of dependency
// information that it contained. Panics with an IOException if the file could
// not be read.
//
// Java: DepsFileRegexParser.parseFile(String)
func (d *DepsFileRegexParser) ParseFile(filePath string) []DependencyInfo {
	// Go-only: language-adapter — Guava Files.newReader(new File(filePath),
	// UTF_8): a BufferedReader over an InputStreamReader of a
	// FileInputStream (malformed input is replaced, not reported).
	f, err := os.Open(filePath)
	if err != nil {
		panic(jutil.IOExceptionFromError(err, filePath))
	}
	return d.ParseFileReader(filePath, jutil.NewInputStreamReader(f, jutil.StandardCharsets_UTF_8))
}

// ParseFileWithContents parses the given file using the passed in fileContents.
//
// Java: DepsFileRegexParser.parseFile(String, String)
func (d *DepsFileRegexParser) ParseFileWithContents(filePath string, fileContents string) []DependencyInfo {
	d.depInfos = []DependencyInfo{}
	if depsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		depsFileRegexParser_logger.Fine("Parsing Dep: " + filePath)
	}
	parseJsString(&d.jsFileLineParser, d, filePath, fileContents)
	return d.depInfos
}

// ParseFileReader parses the file from the given reader and returns a list
// of dependency information that it contained.
//
// Java: DepsFileRegexParser.parseFileReader(String, Reader)
func (d *DepsFileRegexParser) ParseFileReader(filePath string, reader jutil.JavaReader) []DependencyInfo {
	d.depInfos = []DependencyInfo{}
	if depsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		depsFileRegexParser_logger.Fine("Parsing Dep: " + filePath)
	}
	parseJsReader(&d.jsFileLineParser, d, filePath, reader)
	return d.depInfos
}

// Java: DepsFileRegexParser.parseLine(String)
func (d *DepsFileRegexParser) parseLine(line string) (bool, string) {
	hasDependencies := false

	// Quick check that will catch most cases. This is a performance win for teams with a lot of JS.
	if strings.Contains(line, "addDependency") {
		if depLoc := depsFileRegexParser_depMatcher.FindStringSubmatchIndex(line); depLoc != nil {
			hasDependencies = true
			addDependencyParams := line[depLoc[2]:depLoc[3]]
			argsLoc := depsFileRegexParser_depArgsMatch.FindStringSubmatchIndex(addDependencyParams)
			if argsLoc == nil {
				// Although we could recover, we mark this as fatal since there should
				// not be problems with generated deps.js files.
				return false, "Invalid arguments to goog.addDependency(). Found: " + addDependencyParams
			}
			relativePath, errMsg, ok := d.parseJsString(addDependencyParams[argsLoc[2]:argsLoc[3]])
			if !ok {
				return false, errMsg
			}
			path := d.pathTranslator(relativePath)

			provides, errMsg, ok := d.parseJsStringArray(addDependencyParams[argsLoc[4]:argsLoc[5]])
			if !ok {
				return false, errMsg
			}
			var rawFlags string
			hasFlags := argsLoc[8] >= 0
			if hasFlags {
				rawFlags = addDependencyParams[argsLoc[8]:argsLoc[9]]
			}
			loadFlags, errMsg, ok := d.parseLoadFlags(rawFlags, hasFlags)
			if !ok {
				return false, errMsg
			}

			// ES6 modules are require'd by path but do not provide them in the addDependency call.
			if module, ok := loadFlags.Lookup("module"); ok && module == "es6" {
				provides = append(provides, relativePath)
			}

			requireStrings, errMsg, ok := d.parseJsStringArray(addDependencyParams[argsLoc[6]:argsLoc[7]])
			if !ok {
				return false, errMsg
			}
			requires := make([]*DependencyInfo_Require, len(requireStrings))
			for i, r := range requireStrings {
				requires[i] = DependencyInfo_Require_ParsedFromDeps(r)
			}
			depInfo := newSimpleDependencyInfoDirect(
				path,
				d.filePath,
				provides,
				requires,
				nil,
				loadFlags,
				false,
				false,
			)

			if depsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
				depsFileRegexParser_logger.Fine("Found dep: " + depInfo.String())
			}
			d.depInfos = append(d.depInfos, depInfo)
		}
	}

	return !d.shortcutMode || hasDependencies || jsFileRegexParser_isAllWhitespace(line), ""
}

// Java: DepsFileRegexParser.parseLoadFlags(String)
func (d *DepsFileRegexParser) parseLoadFlags(loadFlags string, hasFlags bool) (*jutil.LinkedHashMap[string, string], string, bool) {
	if !hasFlags || loadFlags == "false" {
		return simpleDependencyInfo_EMPTY_FLAGS, "", true
	}
	if loadFlags == "true" {
		return simpleDependencyInfo_Builder_GOOG_MODULE_FLAGS, "", true
	}
	return d.parseJsStringMap(loadFlags)
}
