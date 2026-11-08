// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261104

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

package cli

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

// SetPrintVersion sets whether --version output is requested.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintVersion(printVersion bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printVersion = printVersion
	return c
}

// SetPrintTree sets whether to print the parse tree.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintTree(printTree bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printTree = printTree
	return c
}

// SetPrintTreeJson sets whether to print the parse tree as JSON.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintTreeJson(printTreeJson bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printTreeJson = printTreeJson
	return c
}

// SetPrintAst sets whether to print the DOT AST graph.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintAst(printAst bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printAst = printAst
	return c
}

// SetJscompDevMode sets the compiler developer self-check mode.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJscompDevMode(devMode jscomp.CompilerOptions_DevMode) *AbstractCommandLineRunner_CommandLineConfig {
	c.jscompDevMode = devMode
	return c
}

// SetLoggingLevel sets the java.util.logging level name.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetLoggingLevel(loggingLevel string) *AbstractCommandLineRunner_CommandLineConfig {
	c.loggingLevel = loggingLevel
	return c
}

// SetDefaultToStdin sets whether compilation defaults to reading stdin when no inputs are specified.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetDefaultToStdin(defaultToStdin bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.defaultToStdin = defaultToStdin
	return c
}

// SetJsOutputFile sets the output JS file path.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setJsOutputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsOutputFile(jsOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.jsOutputFile = jsOutputFile
	return c
}

// GetJsOutputFile returns the configured output JS file path.
func (c *AbstractCommandLineRunner_CommandLineConfig) GetJsOutputFile() string {
	return c.jsOutputFile
}

// SetContinueSavedCompilationOpt sets or leaves unset the saved compilation
// state input filename and restored stage.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetContinueSavedCompilationOpt(
	fileName base.OptString, restoredCompilationStage int,
) *AbstractCommandLineRunner_CommandLineConfig {
	if !fileName.IsPresent() {
		c.continueSavedCompilationFileName = base.NoString()
		c.restoredCompilationStage = -1
		return c
	}
	return c.setContinueSavedCompilationOpt(fileName, restoredCompilationStage)
}

// SetSaveCompilationStateOpt sets or clears the saved compilation state output
// filename and stage.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSaveCompilationStateOpt(
	fileName base.OptString, saveAfterCompilationStage int,
) *AbstractCommandLineRunner_CommandLineConfig {
	return c.setSaveCompilationStateOpt(fileName, saveAfterCompilationStage)
}

// SetTypedAstListInputFilename sets the TypedAST list input file path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetTypedAstListInputFilename(
	fileName string,
) *AbstractCommandLineRunner_CommandLineConfig {
	c.typedAstListInputFilename = base.SomeString(fileName)
	return c
}

// SetTypedAstListInputFilenameOpt sets or clears the TypedAST list input file path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetTypedAstListInputFilenameOpt(
	fileName base.OptString,
) *AbstractCommandLineRunner_CommandLineConfig {
	c.typedAstListInputFilename = fileName
	return c
}

// GetTypedAstListInputFilename returns the configured TypedAST list input filename, if present.
func (c *AbstractCommandLineRunner_CommandLineConfig) GetTypedAstListInputFilename() (string, bool) {
	return c.typedAstListInputFilename.Get()
}

// SetParseInlineSourceMaps sets whether inline source maps are parsed.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetParseInlineSourceMaps(parseInlineSourceMaps bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.parseInlineSourceMaps = parseInlineSourceMaps
	return c
}

// SetApplyInputSourceMaps sets whether input source maps are applied.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetApplyInputSourceMaps(applyInputSourceMaps bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.applyInputSourceMaps = applyInputSourceMaps
	return c
}

// SetVariableMapInputFile sets the input variable renaming map path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetVariableMapInputFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.variableMapInputFile = file
	return c
}

// SetPropertyMapInputFile sets the input property renaming map path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPropertyMapInputFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.propertyMapInputFile = file
	return c
}

// SetVariableMapOutputFile sets the output variable renaming map path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetVariableMapOutputFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.variableMapOutputFile = file
	return c
}

// SetCreateNameMapFiles sets whether default variable/property map files are written.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCreateNameMapFiles(create bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.createNameMapFiles = create
	return c
}

// SetPropertyMapOutputFile sets the output property renaming map path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPropertyMapOutputFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.propertyMapOutputFile = file
	return c
}

// SetStringMapOutputPath sets the ReplaceStrings output map path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetStringMapOutputPath(path string) *AbstractCommandLineRunner_CommandLineConfig {
	c.stringMapOutputPath = path
	return c
}

// SetInstrumentationMappingFile sets the production instrumentation mapping output path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetInstrumentationMappingFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.instrumentationMappingFile = file
	return c
}

// SetCodingConvention sets the CodingConvention used during compilation.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCodingConvention(cc jscomp.CodingConvention) *AbstractCommandLineRunner_CommandLineConfig {
	c.codingConvention = cc
	return c
}

// SetSummaryDetailLevel sets the diagnostic summary detail level.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSummaryDetailLevel(level int) *AbstractCommandLineRunner_CommandLineConfig {
	c.summaryDetailLevel = level
	return c
}

// SetOutputWrapper sets the single-binary output wrapper template.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputWrapper(outputWrapper string) *AbstractCommandLineRunner_CommandLineConfig {
	c.outputWrapper = outputWrapper
	return c
}

// SetChunkOutputPathPrefix sets the prefix prepended to chunk output file names.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkOutputPathPrefix(prefix string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkOutputPathPrefix = prefix
	return c
}

// SetChunkOutputFiles sets explicit <chunk>:<output_file> mappings.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkOutputFiles(chunkOutputFiles []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkOutputFiles = slices.Clone(chunkOutputFiles)
	return c
}

// SetChunkConformanceFiles sets explicit <chunk>:<conformance_file> mappings.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkConformanceFiles(chunkConformanceFiles []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkConformanceFiles = slices.Clone(chunkConformanceFiles)
	return c
}

// SetCreateSourceMap sets the source map output path or pattern.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCreateSourceMap(createSourceMap string) *AbstractCommandLineRunner_CommandLineConfig {
	c.createSourceMap = createSourceMap
	return c
}

// SetSourceMapDetailLevel sets the source map detail level.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapDetailLevel(level jscomp.SourceMap_DetailLevel) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapDetailLevel = level
	return c
}

// SetSourceMapFormat sets the source map format.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapFormat(format jscomp.SourceMap_Format) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapFormat = format
	return c
}

// SetBrowserFeaturesetYear sets the browser featureset year.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetBrowserFeaturesetYear(year int) *AbstractCommandLineRunner_CommandLineConfig {
	c.browserFeaturesetYear = year
	return c
}

// SetTweakProcessing sets the tweak processing mode.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetTweakProcessing(tp jscomp.CompilerOptions_TweakProcessing) *AbstractCommandLineRunner_CommandLineConfig {
	c.tweakProcessing = tp
	return c
}

// SetCharset sets the input/output charset name.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCharset(charset string) *AbstractCommandLineRunner_CommandLineConfig {
	c.charset = charset
	return c
}

// SetDependencyOptions sets the dependency sorting/pruning options.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetDependencyOptions(opts *jscomp.DependencyOptions) *AbstractCommandLineRunner_CommandLineConfig {
	c.dependencyOptions = opts
	return c
}

// SetOutputChunkDependencies sets the chunk dependencies JSON output path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputChunkDependencies(path string) *AbstractCommandLineRunner_CommandLineConfig {
	c.outputChunkDependencies = path
	return c
}

// SetOutputBundles sets the output bundle file paths (filtering empty names).
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputBundles(bundles []string) *AbstractCommandLineRunner_CommandLineConfig {
	out := make([]string, 0, len(bundles))
	for _, b := range bundles {
		if b != "" {
			out = append(out, b)
		}
	}
	c.outputBundles = out
	return c
}

// SetSkipNormalOutputs sets whether normal compilation and JS outputs are skipped.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSkipNormalOutputs(skip bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.skipNormalOutputs = skip
	return c
}

// SetManifestMaps sets the chunk manifest map output paths.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetManifestMaps(manifestMaps []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.manifestMaps = slices.Clone(manifestMaps)
	return c
}

// SetProcessCommonJSModules sets whether CommonJS modules are rewritten.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetProcessCommonJSModules(process bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.processCommonJSModules = process
	return c
}

// SetModuleRoots sets the module root prefixes.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetModuleRoots(roots []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.moduleRoots = slices.Clone(roots)
	return c
}

// SetWarningsAllowFile sets the warnings allowlist file path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetWarningsAllowFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.warningsAllowFile = file
	return c
}

// SetHideWarningsFor sets path prefixes for which warnings are hidden.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetHideWarningsFor(paths []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.hideWarningsFor = slices.Clone(paths)
	return c
}

// SetAngularPass sets whether Angular @ngInject processing is enabled.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetAngularPass(angular bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.angularPass = angular
	return c
}

// SetJsonStreamMode sets the JSON stream input/output mode.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsonStreamMode(mode jscomp.CompilerOptions_JsonStreamMode) *AbstractCommandLineRunner_CommandLineConfig {
	c.jsonStreamMode = mode
	return c
}

// SetErrorFormat sets the CLI error output format option.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetErrorFormat(format AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) *AbstractCommandLineRunner_CommandLineConfig {
	c.errorFormat = format
	return c
}

// SetJsonWarningsFile sets the JSON warnings report file path.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsonWarningsFile(file string) *AbstractCommandLineRunner_CommandLineConfig {
	c.jsonWarningsFile = file
	return c
}

// SetIncludeTrailingNewline sets whether a trailing newline is appended to JS outputs.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetIncludeTrailingNewline(include bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.includeTrailingNewline = include
	return c
}
