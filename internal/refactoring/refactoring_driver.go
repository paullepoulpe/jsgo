// Port of com/google/javascript/refactoring/RefactoringDriver.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/Scanner.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/JsFlumeCallback.java @ closure-compiler v20261021
// Port of com/google/javascript/refactoring/RefactoringUtils.java @ closure-compiler v20261021

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

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// Package refactoring implements AST pattern matching, refactoring drivers,
// and automated source-code suggested fixes (SuggestedFix, ErrorToFixMapper,
// FixingErrorManager, and ApplySuggestedFixes).
package refactoring

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Scanner is the class that developers should implement to perform a JsFlume
// refactoring.
//
// Java: public abstract class Scanner implements Serializable
type Scanner interface {
	// Matches returns true if the given node and node traversal should match
	// for this particular scanner. Typically this function uses the Matcher
	// class or predefined matchers from Matchers to match against the Node and
	// NodeMetadata.
	//
	// If this function returns true, a Match for this node will be passed to
	// ProcessMatch and all matches will be passed to ProcessAllMatches at the
	// end of the traversal.
	//
	// Java: Scanner.matches(Node, NodeMetadata)
	Matches(node rhino.Node, t *NodeMetadata) bool

	// Java: Scanner.processMatch(Match)
	ProcessMatch(match *Match) []*SuggestedFix

	// Java: Scanner.processAllMatches(Collection<Match>)
	ProcessAllMatches(matches []*Match) []*SuggestedFix
}

// ScannerBase holds Scanner's method bodies; subclasses embed it and define
// Matches.
type ScannerBase struct{}

// ProcessMatch processes one Match at a time. There is no order guaranteed for
// when this function will be called with the Match.
//
// Java: Scanner.processMatch(Match)
func (s *ScannerBase) ProcessMatch(match *Match) []*SuggestedFix {
	return []*SuggestedFix{}
}

// ProcessAllMatches processes every given match at one time. This function
// can be used when the refactoring needs the information from the entire run
// to perform the refactoring, such as moving functions around.
//
// Java: Scanner.processAllMatches(Collection<Match>)
func (s *ScannerBase) ProcessAllMatches(matches []*Match) []*SuggestedFix {
	return []*SuggestedFix{}
}

// JsFlumeCallback is a compiler node traversal callback that invokes matchers
// against every node and keeps track of any suggested fixes from the
// refactoring.
//
// Java: final class JsFlumeCallback implements NodeTraversal.Callback
type JsFlumeCallback struct {
	scanner            Scanner
	includeFilePattern *jutil.Pattern // may be nil
	matches            []*Match
	fixes              []*SuggestedFix
}

// Java: JsFlumeCallback(Scanner, Pattern)
func NewJsFlumeCallback(scanner Scanner, includeFilePattern *jutil.Pattern) *JsFlumeCallback {
	return &JsFlumeCallback{
		scanner:            scanner,
		includeFilePattern: includeFilePattern,
	}
}

// Java: JsFlumeCallback.getMatches()
func (c *JsFlumeCallback) GetMatches() []*Match {
	return c.matches
}

// Java: JsFlumeCallback.getFixes()
func (c *JsFlumeCallback) GetFixes() []*SuggestedFix {
	return c.fixes
}

// Java: JsFlumeCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *JsFlumeCallback) ShouldTraverse(nodeTraversal *jscomp.NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsFromExterns() {
		return false
	}
	if n.IsScript() {
		filename := n.GetSourceFileName()
		if c.includeFilePattern != nil &&
			c.includeFilePattern.Pattern() != "" &&
			filename != nil && *filename != "" {
			return c.includeFilePattern.Matcher(*filename).Find()
		}
	}
	return true
}

// Java: JsFlumeCallback.visit(NodeTraversal, Node, Node)
func (c *JsFlumeCallback) Visit(t *jscomp.NodeTraversal, n rhino.Node, parent rhino.Node) {
	metadata := NodeMetadata_FromTraversal(t)
	if c.scanner.Matches(n, metadata) {
		match := NewMatch(n, metadata)
		c.fixes = append(c.fixes, c.scanner.ProcessMatch(match)...)
		c.matches = append(c.matches, match)
	}
}

// RefactoringUtils holds utility methods for refactoring code.
type RefactoringUtils struct{}

// Java: RefactoringUtils()
func NewRefactoringUtils() *RefactoringUtils {
	return &RefactoringUtils{}
}

// RefactoringUtils_IsInClosurizedFile looks for a goog.require(),
// goog.provide() or goog.module() call in the node's file.
//
// Java: RefactoringUtils.isInClosurizedFile(Node, NodeMetadata)
func RefactoringUtils_IsInClosurizedFile(node rhino.Node, metadata *NodeMetadata) bool {
	script := jscomp.NodeUtil_GetEnclosingScript(node)
	if script.IsNil() {
		return false
	}

	reqMatcher := Matchers_GoogRequire()
	modOrProvideMatcher := Matchers_GoogModuleOrProvide()
	for child := script.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if exprStmt, ok := ast.AsExprStmt(child); ok {
			expr := exprStmt.Expr().Node()
			if expr.IsCall() && (reqMatcher.Matches(expr, metadata) || modOrProvideMatcher.Matches(expr, metadata)) {
				return true
			}
		} else if child.IsVar() && child.GetBooleanProp(rhino.Node_IS_NAMESPACE) {
			return true
		}
	}
	return false
}

// RefactoringDriver is the primary driver of a refactoring. This class
// collects the inputs, runs the refactoring over the compiled input, and
// then collects the suggested fixes based on the refactoring.
//
// Java: final class refactoring.RefactoringDriver
type RefactoringDriver struct {
	compiler *jscomp.Compiler
	rootNode rhino.Node
}

// Java: RefactoringDriver(List<SourceFile>, List<SourceFile>, CompilerOptions)
func newRefactoringDriver(
	inputs []*jscomp.SourceFile,
	externs []*jscomp.SourceFile,
	compilerOptions *jscomp.CompilerOptions) *RefactoringDriver {
	compiler := refactoringDriver_createCompiler(inputs, externs, compilerOptions)
	return &RefactoringDriver{
		compiler: compiler,
		rootNode: compiler.GetRoot(),
	}
}

// DriveWithIncludeFilePattern runs a refactoring and returns any suggested
// fixes as a result. includeFilePattern may be nil.
//
// Java: RefactoringDriver.drive(Scanner, Pattern)
func (d *RefactoringDriver) DriveWithIncludeFilePattern(scanner Scanner, includeFilePattern *jutil.Pattern) []*SuggestedFix {
	callback := NewJsFlumeCallback(scanner, includeFilePattern)
	jscomp.NodeTraversal_Traverse(d.compiler, d.rootNode, callback)
	return append(callback.GetFixes(), scanner.ProcessAllMatches(callback.GetMatches())...)
}

// Drive runs a refactoring and returns any suggested fixes as a result.
//
// Java: RefactoringDriver.drive(Scanner)
func (d *RefactoringDriver) Drive(scanner Scanner) []*SuggestedFix {
	return d.DriveWithIncludeFilePattern(scanner, nil)
}

// Java: RefactoringDriver.getCompiler()
func (d *RefactoringDriver) GetCompiler() *jscomp.Compiler {
	return d.compiler
}

// Java: RefactoringDriver.createCompiler(List<SourceFile>, List<SourceFile>, CompilerOptions)
func refactoringDriver_createCompiler(
	inputs []*jscomp.SourceFile, externs []*jscomp.SourceFile, compilerOptions *jscomp.CompilerOptions) *jscomp.Compiler {
	compiler := jscomp.NewCompilerWithErrorManager(jscomp.NewBlackHoleErrorManager())
	compiler.DisableThreads()
	compiler.Compile(externs, inputs, compilerOptions)
	return compiler
}

// TODO(tbreisacher): Make this package-private by refactoring tests so they
// don't need to call it directly.

// Java: RefactoringDriver.getCompilerOptions()
func RefactoringDriver_GetCompilerOptions() *jscomp.CompilerOptions {
	options := jscomp.NewCompilerOptions()
	options.SetLanguageIn(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5)
	options.SetSummaryDetailLevel(0)

	options.SetDependencyOptions(jscomp.DependencyOptions_SortOnly())

	options.SetChecksOnly(true)
	options.SetContinueAfterErrors(true)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	options.SetCheckSuspiciousCode(true)
	options.SetCheckSymbols(true)
	options.SetCheckTypes(true)
	// TODO(bangert): Remove this -- we want to rewrite code before closure
	// syntax is removed. Unfortunately, setClosurePass is required, or code
	// doesn't type check.
	options.SetClosurePass(true)
	options.SetGenerateExports(true)
	options.SetPreserveClosurePrimitives(true)

	options.SetWarningLevel(jscomp.DiagnosticGroups_EXTRA_REQUIRE, jscomp.CheckLevel_WARNING)

	return options
}

// Java: static class RefactoringDriver.Builder
type RefactoringDriver_Builder struct {
	inputs          []*jscomp.SourceFile
	externs         []*jscomp.SourceFile
	compilerOptions *jscomp.CompilerOptions
}

// Java: RefactoringDriver.Builder.TO_SOURCE_FILE_FN
func refactoringDriver_Builder_TO_SOURCE_FILE_FN(file string) *jscomp.SourceFile {
	return jscomp.SourceFile_Builder_().WithPath(file).Build()
}

// Java: RefactoringDriver.Builder()
func NewRefactoringDriver_Builder() *RefactoringDriver_Builder {
	return &RefactoringDriver_Builder{
		compilerOptions: RefactoringDriver_GetCompilerOptions(),
	}
}

// Java: RefactoringDriver.Builder.addExternsFromFile(String)
func (b *RefactoringDriver_Builder) AddExternsFromFile(filename string) *RefactoringDriver_Builder {
	b.externs = append(b.externs, jscomp.SourceFile_FromFile(filename))
	return b
}

// Java: RefactoringDriver.Builder.addExternsFromFile(Iterable<String>)
func (b *RefactoringDriver_Builder) AddExternsFromFileOfList(externs []string) *RefactoringDriver_Builder {
	for _, f := range externs {
		b.externs = append(b.externs, refactoringDriver_Builder_TO_SOURCE_FILE_FN(f))
	}
	return b
}

// Java: RefactoringDriver.Builder.addExternsFromCode(String)
func (b *RefactoringDriver_Builder) AddExternsFromCode(code string) *RefactoringDriver_Builder {
	b.externs = append(b.externs, jscomp.SourceFile_FromCode("externs", code))
	return b
}

// Java: RefactoringDriver.Builder.addExterns(Iterable<SourceFile>)
func (b *RefactoringDriver_Builder) AddExterns(externs []*jscomp.SourceFile) *RefactoringDriver_Builder {
	b.externs = append(b.externs, externs...)
	return b
}

// Java: RefactoringDriver.Builder.addInputsFromFile(String)
func (b *RefactoringDriver_Builder) AddInputsFromFile(filename string) *RefactoringDriver_Builder {
	b.inputs = append(b.inputs, jscomp.SourceFile_FromFile(filename))
	return b
}

// Java: RefactoringDriver.Builder.addInputsFromFile(Iterable<String>)
func (b *RefactoringDriver_Builder) AddInputsFromFileOfList(inputs []string) *RefactoringDriver_Builder {
	for _, f := range inputs {
		b.inputs = append(b.inputs, refactoringDriver_Builder_TO_SOURCE_FILE_FN(f))
	}
	return b
}

// Java: RefactoringDriver.Builder.addInputsFromCode(String)
func (b *RefactoringDriver_Builder) AddInputsFromCode(code string) *RefactoringDriver_Builder {
	return b.AddInputsFromCodeWithFilename(code, "input")
}

// Java: RefactoringDriver.Builder.addInputsFromCode(String, String)
func (b *RefactoringDriver_Builder) AddInputsFromCodeWithFilename(code string, filename string) *RefactoringDriver_Builder {
	b.inputs = append(b.inputs, jscomp.SourceFile_FromCode(filename, code))
	return b
}

// Java: RefactoringDriver.Builder.addInputs(Iterable<SourceFile>)
func (b *RefactoringDriver_Builder) AddInputs(inputs []*jscomp.SourceFile) *RefactoringDriver_Builder {
	b.inputs = append(b.inputs, inputs...)
	return b
}

// Java: RefactoringDriver.Builder.withCompilerOptions(CompilerOptions)
func (b *RefactoringDriver_Builder) WithCompilerOptions(compilerOptions *jscomp.CompilerOptions) *RefactoringDriver_Builder {
	b.compilerOptions = jutil.CheckNotNull(compilerOptions)
	return b
}

// Java: RefactoringDriver.Builder.build()
func (b *RefactoringDriver_Builder) Build() *RefactoringDriver {
	return newRefactoringDriver(slices.Clone(b.inputs), slices.Clone(b.externs), b.compilerOptions)
}
