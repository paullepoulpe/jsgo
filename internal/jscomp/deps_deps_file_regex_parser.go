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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DepsFileRegexParser is a parser that can extract dependency information
// from existing deps.js files.
//
// See //javascript/closure/deps.js for an example file.
type DepsFileRegexParser struct {
	JsFileLineParserBase

	// The dependency information extracted from the current file.
	depInfos []DependencyInfo

	// Translates paths in different build systems.
	pathTranslator func(string) string
}

var _ JsFileLineParser = (*DepsFileRegexParser)(nil)

var depsFileRegexParser_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.DepsFileRegexParser")

// Go-only: language-adapter — Java regex. The Java class holds one Matcher
// per pattern (Pattern.compile(...).matcher("")) and resets it per line; the
// patterns are stateless here, so they are compiled once. They are the Java
// patterns rewritten in RE2 syntax with Java's meaning: \s is Java's
// [ \t\n\x0B\f\r] (RE2's \s lacks \x0B), `.` excludes Java's line
// terminators (\n, \r, \x{85}, \x{2028}, \x{2029}), and ^...$ anchor the
// whole input as Matcher.matches() does. They use no back-references or
// lookaround, and RE2's leftmost-first matching selects the same submatches
// as Java's backtracking engine.

// depsFileRegexParser_depMatcher is the pattern for matching JavaScript
// string literals. The group is: goog.addDependency({1});
//
// Java: Pattern.compile("\\s*goog.addDependency\\((.*)\\);?\\s*").matcher("")
var depsFileRegexParser_depMatcher = regexp.MustCompile(
	`^[\t\n\x0B\f\r ]*goog[^\n\r\x{85}\x{2028}\x{2029}]addDependency\(([^\n\r\x{85}\x{2028}\x{2029}]*)\);?[\t\n\x0B\f\r ]*$`)

// depsFileRegexParser_depArgsMatch is the pattern for matching the args of
// a goog.addDependency(). The group is: goog.addDependency({1}, {2}, {3},
// {4?});
//
// Java:
//
//	Pattern.compile(
//	        "\\s*([^,]*), (\\[[^\\]]*\\]), (\\[[^\\]]*\\])"
//	            + "(?:, (true|false|\\{[^{}]*\\}))?\\s*")
//	    .matcher("")
var depsFileRegexParser_depArgsMatch = regexp.MustCompile(
	`^[\t\n\x0B\f\r ]*([^,]*), (\[[^\]]*\]), (\[[^\]]*\])` +
		`(?:, (true|false|\{[^{}]*\}))?[\t\n\x0B\f\r ]*$`)

// depsFileRegexParser_group is Matcher.group(int) after a successful
// match: the group's text, or null if the group did not participate.
//
// Go-only: language-adapter — java.util.regex.Matcher.group(int).
func depsFileRegexParser_group(s string, loc []int, group int) *string {
	if loc[2*group] < 0 {
		return nil
	}
	g := s[loc[2*group]:loc[2*group+1]]
	return &g
}

// NewDepsFileRegexParser is the constructor; errorManager handles parse
// errors.
//
// Java: DepsFileRegexParser(ErrorManager)
func NewDepsFileRegexParser(errorManager ErrorManager) *DepsFileRegexParser {
	return NewDepsFileRegexParserWithPathTranslator(func(s string) string { return s }, errorManager)
}

// NewDepsFileRegexParserWithPathTranslator: pathTranslator translates paths
// in different build systems; errorManager handles parse errors.
//
// Java: DepsFileRegexParser(Function<String, String>, ErrorManager)
func NewDepsFileRegexParserWithPathTranslator(pathTranslator func(string) string, errorManager ErrorManager) *DepsFileRegexParser {
	d := &DepsFileRegexParser{}
	d.InitJsFileLineParserBase(d, errorManager)
	d.pathTranslator = pathTranslator
	return d
}

// ParseFile parses the given file and returns a list of dependency
// information that it contained.
//
// filePath is the path to the file to parse. Returns a list of
// DependencyInfo objects. Panics with an IOException if the file could not
// be read.
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

// ParseFileWithContents parses the given file and returns a list of
// dependency information that it contained. It uses the passed in
// fileContents instead of reading the file.
//
// filePath is the path to the file to parse; fileContents are the contents
// to parse. Returns a list of DependencyInfo objects.
//
// Java: DepsFileRegexParser.parseFile(String, String)
func (d *DepsFileRegexParser) ParseFileWithContents(filePath string, fileContents string) []DependencyInfo {
	return d.ParseFileReader(filePath, jutil.NewStringReader(fileContents))
}

// ParseFileReader parses the file from the given reader and returns a list
// of dependency information that it contained.
//
// filePath is the path to the file to parse; reader is a reader for the
// file. Returns a list of DependencyInfo objects.
//
// Java: DepsFileRegexParser.parseFileReader(String, Reader)
func (d *DepsFileRegexParser) ParseFileReader(filePath string, reader jutil.JavaReader) []DependencyInfo {
	d.depInfos = []DependencyInfo{}
	if depsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
		depsFileRegexParser_logger.Fine("Parsing Dep: " + filePath)
	}
	d.DoParse(filePath, reader)
	return d.depInfos
}

// ParseLine extracts dependency information from lines that look like
// goog.addDependency('pathRelativeToClosure', ['provides'], ['requires']);
// Adds the dependencies to depInfos.
//
// Panics with *JsFileLineParser_ParseException if the given line has a
// malformed goog.addDependency().
//
// Java: DepsFileRegexParser.parseLine(String)
func (d *DepsFileRegexParser) ParseLine(line string) bool {
	hasDependencies := false

	// Quick check that will catch most cases. This is a performance win for teams with a lot of JS.
	if jutil.StringIndexOf(line, "addDependency") >= 0 {
		// See if the line looks like: goog.addDependency(...)
		depLoc := depsFileRegexParser_depMatcher.FindStringSubmatchIndex(line)
		if depLoc != nil {
			hasDependencies = true
			addDependencyParams := *depsFileRegexParser_group(line, depLoc, 1)
			argsLoc := depsFileRegexParser_depArgsMatch.FindStringSubmatchIndex(addDependencyParams)
			// Extract the three parameters.
			if argsLoc == nil {
				// Although we could recover, we mark this as fatal since there should
				// not be problems with generated deps.js files.
				panic(NewJsFileLineParser_ParseException("Invalid arguments to goog.addDependency(). Found: "+
					addDependencyParams, true))
			}
			// Parse the file path.
			relativePath := d.ParseJsString(*depsFileRegexParser_group(addDependencyParams, argsLoc, 1))
			path := d.pathTranslator(relativePath)

			provides := d.ParseJsStringArray(*depsFileRegexParser_group(addDependencyParams, argsLoc, 2))
			loadFlags := d.ParseLoadFlags(depsFileRegexParser_group(addDependencyParams, argsLoc, 4))

			// ES6 modules are require'd by path but do not provide them in the addDependency call.
			if module, ok := loadFlags.Lookup("module"); ok && "es6" == module {
				provides = append(provides, relativePath)
			}

			requireStrings := d.ParseJsStringArray(*depsFileRegexParser_group(addDependencyParams, argsLoc, 3))
			requires := make([]*DependencyInfo_Require, 0, len(requireStrings))
			for _, r := range requireStrings {
				requires = append(requires, DependencyInfo_Require_ParsedFromDeps(r))
			}
			depInfo :=
				SimpleDependencyInfo_Builder_(path, d.filePath).
					SetProvides(provides).
					SetRequires(requires).
					SetLoadFlags(loadFlags).
					Build()

			if depsFileRegexParser_logger.IsLoggable(jutil.Level_FINE) {
				depsFileRegexParser_logger.Fine("Found dep: " + depInfo.String())
			}
			d.depInfos = append(d.depInfos, depInfo)
		}
	}

	return !d.shortcutMode || hasDependencies ||
		jsFileRegexParser_isAllWhitespace(line)
}

// ParseLoadFlags: loadFlags is a nullable Java String.
//
// Java: DepsFileRegexParser.parseLoadFlags(String)
func (d *DepsFileRegexParser) ParseLoadFlags(loadFlags *string) *jutil.LinkedHashMap[string, string] {
	if loadFlags == nil || *loadFlags == "false" {
		return jutil.NewLinkedHashMap[string, string]()
	} else if *loadFlags == "true" {
		m := jutil.NewLinkedHashMap[string, string]()
		m.Put("module", "goog")
		return m
	} else {
		return d.ParseJsStringMap(*loadFlags)
	}
}
