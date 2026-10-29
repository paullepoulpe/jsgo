// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/OutputCharsetEncoder.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// The number of characters after which we insert a line break in the code
//
// Java: CompilerOptions.DEFAULT_LINE_LENGTH_THRESHOLD
const CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD = 500

// Java: CompilerOptions.POLYMER_PROPERTY_RESERVED_FIRST_CHARS
var compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS = jutil.NewLinkedHashSet(
	jutil.StringToUTF16("ABCDEFGHIJKLMNOPQRSTUVWXYZ$")...)

// Java: CompilerOptions.POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS
var compilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS = jutil.NewLinkedHashSet[uint16]('_', '$')

// Java: CompilerOptions.ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
var compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS = jutil.NewLinkedHashSet[uint16]('$')

// Java: CompilerOptions.getAngularPropertyReservedFirstChars()
func CompilerOptions_GetAngularPropertyReservedFirstChars() *jutil.LinkedHashSet[uint16] {
	return compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
}

// Java: CompilerOptions.getPolymerPropertyReservedFirstChars()
func CompilerOptions_GetPolymerPropertyReservedFirstChars() *jutil.LinkedHashSet[uint16] {
	return compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS
}

// Java: CompilerOptions.setOutputJs(OutputJs)
func (o *CompilerOptions) SetOutputJs(outputJs CompilerOptions_OutputJs) {
	o.outputJs = outputJs
}

// Java: CompilerOptions.getOutputJs()
func (o *CompilerOptions) GetOutputJs() CompilerOptions_OutputJs {
	return o.outputJs
}

// Java: CompilerOptions.setPrintExterns(boolean)
func (o *CompilerOptions) SetPrintExterns(printExterns bool) {
	o.printExterns = printExterns
}

// Java: CompilerOptions.shouldPrintExterns()
func (o *CompilerOptions) ShouldPrintExterns() bool {
	return o.printExterns || o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS
}

// Java: CompilerOptions.setGenerateExports(boolean)
func (o *CompilerOptions) SetGenerateExports(generateExports bool) {
	o.generateExports = generateExports
}

// Java: CompilerOptions.shouldGenerateExports()
func (o *CompilerOptions) ShouldGenerateExports() bool {
	return o.generateExports
}

// Java: CompilerOptions.setExportLocalPropertyDefinitions(boolean)
func (o *CompilerOptions) SetExportLocalPropertyDefinitions(export bool) {
	o.exportLocalPropertyDefinitions = export
}

// Java: CompilerOptions.shouldExportLocalPropertyDefinitions()
func (o *CompilerOptions) ShouldExportLocalPropertyDefinitions() bool {
	return o.exportLocalPropertyDefinitions
}

// Normally, when there are an equal number of single and double quotes in a string, the compiler
// will use double quotes. Set this to true to prefer single quotes.
//
// Java: CompilerOptions.setPreferSingleQuotes(boolean)
func (o *CompilerOptions) SetPreferSingleQuotes(enabled bool) {
	o.preferSingleQuotes = enabled
}

// Java: CompilerOptions.shouldPreferSingleQuotes()
func (o *CompilerOptions) ShouldPreferSingleQuotes() bool {
	return o.preferSingleQuotes
}

// Some people want to put arbitrary user input into strings, which are then run through the
// compiler. These scripts are then put into HTML. By default, we assume strings are untrusted. If
// the compiler is run from the command-line, we assume that strings are trusted.
//
// Java: CompilerOptions.setTrustedStrings(boolean)
func (o *CompilerOptions) SetTrustedStrings(yes bool) {
	o.trustedStrings = yes
}

// Java: CompilerOptions.assumeTrustedStrings()
func (o *CompilerOptions) AssumeTrustedStrings() bool {
	return o.trustedStrings
}

// Sets the output charset.
//
// Java: CompilerOptions.setOutputCharset(Charset)
func (o *CompilerOptions) SetOutputCharset(charset *jutil.Charset) {
	o.outputCharset = charset
}

// Gets the output charset.
//
// Java: CompilerOptions.getOutputCharset()
func (o *CompilerOptions) GetOutputCharset() *jutil.Charset {
	return o.outputCharset
}

// OutputCharsetEncoder is the output charset encoder for CodeGenerator that
// delegates to a CharsetEncoder.
type OutputCharsetEncoder struct {
	// Go-only: language-adapter — a CharsetEncoder is only used for
	// canEncode(char), which depends on the charset alone, so the charset
	// stands in for it (nil is Java's null encoder).
	encoder *jutil.Charset
}

// Java: OutputCharsetEncoder(Charset)
func NewOutputCharsetEncoder(outputCharset *jutil.Charset) *OutputCharsetEncoder {
	e := &OutputCharsetEncoder{}
	if outputCharset == nil || outputCharset == jutil.StandardCharsets_US_ASCII {
		// If we want our default (pretending to be UTF-8, but escaping
		// anything outside of straight ASCII), then don't use the encoder,
		// but just special-case the code. This keeps the normal path
		// through the code identical to how it's been for years.
		e.encoder = nil
	} else {
		e.encoder = outputCharset
	}
	return e
}

// Java: OutputCharsetEncoder.canEncode(char)
func (e *OutputCharsetEncoder) CanEncode(c uint16) bool {
	return e.encoder != nil && e.encoder.CanEncodeChar(c)
}

// Java: CompilerOptions.setInstrumentForCoverageOnly(boolean)
func (o *CompilerOptions) SetInstrumentForCoverageOnly(instrumentForCoverageOnly bool) {
	o.instrumentForCoverageOnly = instrumentForCoverageOnly
}

// Java: CompilerOptions.getInstrumentForCoverageOnly()
func (o *CompilerOptions) GetInstrumentForCoverageOnly() bool {
	return o.instrumentForCoverageOnly
}

// Sets file to output in-progress TypedAST format to. DO NOT USE!
//
// The "TypedAST format" is currently a gzipped TypedAst proto but this is not stable.
//
// Java: CompilerOptions.setTypedAstOutputFile(Path)
func (o *CompilerOptions) SetTypedAstOutputFile(file *jutil.Path) {
	o.typedAstOutputFile = file
}

// Java: CompilerOptions.getTypedAstOutputFile()
func (o *CompilerOptions) GetTypedAstOutputFile() *jutil.Path {
	return o.typedAstOutputFile
}

// Java: CompilerOptions.getSourceMapOutputPath()
func (o *CompilerOptions) GetSourceMapOutputPath() *string {
	return o.sourceMapOutputPath
}

// Java: CompilerOptions.shouldGatherSourceMapInfo()
func (o *CompilerOptions) ShouldGatherSourceMapInfo() bool {
	return o.shouldAlwaysGatherSourceMapInfo || !(o.sourceMapOutputPath == nil || *o.sourceMapOutputPath == "")
}

// Java: CompilerOptions.setAlwaysGatherSourceMapInfo(boolean)
func (o *CompilerOptions) SetAlwaysGatherSourceMapInfo(shouldAlwaysGatherSourceMapInfo bool) {
	o.shouldAlwaysGatherSourceMapInfo = shouldAlwaysGatherSourceMapInfo
}

// Sets the input sourcemap files, indexed by the JS files they refer to.
//
// @param inputSourceMaps the collection of input sourcemap files
//
// Java: CompilerOptions.setInputSourceMaps(ImmutableMap<String, SourceMapInput>)
func (o *CompilerOptions) SetInputSourceMaps(inputSourceMaps *jutil.LinkedHashMap[string, *SourceMapInput]) {
	o.inputSourceMaps = inputSourceMaps
}

// Java: CompilerOptions.getInputSourceMaps()
func (o *CompilerOptions) GetInputSourceMaps() *jutil.LinkedHashMap[string, *SourceMapInput] {
	return o.inputSourceMaps
}

// Sets the variable and property renaming policies for the compiler, in a way that clears
// warnings about the renaming policy being uninitialized from flags.
//
// Java: CompilerOptions.setRenamingPolicy(VariableRenamingPolicy, PropertyRenamingPolicy)
func (o *CompilerOptions) SetRenamingPolicy(newVariablePolicy VariableRenamingPolicy, newPropertyPolicy PropertyRenamingPolicy) {
	o.variableRenaming = newVariablePolicy
	o.propertyRenaming = newPropertyPolicy
}

// Java: CompilerOptions.setNameGenerator(NameGenerator)
func (o *CompilerOptions) SetNameGenerator(nameGenerator NameGenerator) {
	o.nameGenerator = nameGenerator
}

// Java: CompilerOptions.getNameGenerator()
func (o *CompilerOptions) GetNameGenerator() NameGenerator {
	return o.nameGenerator
}

// Useful for tests to avoid having to declare two chunks
//
// Java: CompilerOptions.setRenamePrefixNamespaceAssumeCrossChunkNames(boolean)
func (o *CompilerOptions) SetRenamePrefixNamespaceAssumeCrossChunkNames(assume bool) {
	o.renamePrefixNamespaceAssumeCrossChunkNames = assume
}

// Java: CompilerOptions.assumeCrossChunkNamesForRenamePrefixNamespace()
func (o *CompilerOptions) AssumeCrossChunkNamesForRenamePrefixNamespace() bool {
	return o.renamePrefixNamespaceAssumeCrossChunkNames
}

// SetOptimizeLocalAccessForGlobalSymbolNamespace: optimize access to global
// symbols.
//
// Java: CompilerOptions.setOptimizeLocalAccessForGlobalSymbolNamespace(OptimizeLocalAccess)
func (o *CompilerOptions) SetOptimizeLocalAccessForGlobalSymbolNamespace(optimize CompilerOptions_OptimizeLocalAccess) {
	o.optimizeLocalAccessForGlobalSymbolNamespace = optimize
}

// Java: CompilerOptions.getOptimizeLocalAccessForGlobalSymbolNamespace()
func (o *CompilerOptions) GetOptimizeLocalAccessForGlobalSymbolNamespace() CompilerOptions_OptimizeLocalAccess {
	return o.optimizeLocalAccessForGlobalSymbolNamespace
}

// @param replaceIdGenerators the replaceIdGenerators to set
//
// Java: CompilerOptions.setReplaceIdGenerators(boolean)
func (o *CompilerOptions) SetReplaceIdGenerators(replaceIdGenerators bool) {
	o.replaceIdGenerators = replaceIdGenerators
}

// Java: CompilerOptions.shouldReplaceIdGenerators()
func (o *CompilerOptions) ShouldReplaceIdGenerators() bool {
	return o.replaceIdGenerators
}

// Sets the id generators to replace.
//
// Java: CompilerOptions.setIdGenerators(Set<String>)
func (o *CompilerOptions) SetIdGenerators(idGenerators *jutil.LinkedHashSet[string]) {
	builder := jutil.NewLinkedHashMap[string, RenamingMap]()
	for name := range idGenerators.All() {
		builder.Put(name, RenamingMap(RenamingToken_INCONSISTENT))
	}
	o.idGenerators = builder
}

// Sets the id generators to replace.
//
// Java: CompilerOptions.setIdGenerators(Map<String, RenamingMap>)
func (o *CompilerOptions) SetIdGeneratorsOfMap(idGenerators *jutil.LinkedHashMap[string, RenamingMap]) {
	o.idGenerators = jutil.NewLinkedHashMapFrom(idGenerators)
}

// Java: CompilerOptions.getIdGenerators()
func (o *CompilerOptions) GetIdGenerators() *jutil.LinkedHashMap[string, RenamingMap] {
	return o.idGenerators
}

// A previous map of ids (serialized to a string by a previous compile). This will be used as a
// hint during the ReplaceIdGenerators pass, which will attempt to reuse the same ids.
//
// Java: CompilerOptions.setIdGeneratorsMap(String)
func (o *CompilerOptions) SetIdGeneratorsMap(previousMappings *string) {
	o.idGeneratorsMapSerialized = previousMappings
}

// Java: CompilerOptions.getIdGeneratorsMapSerialized()
func (o *CompilerOptions) GetIdGeneratorsMapSerialized() *string {
	return o.idGeneratorsMapSerialized
}

// Sets the hash function to use for Xid
//
// Java: CompilerOptions.setXidHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetXidHashFunction(xidHashFunction Xid_HashFunction) {
	o.xidHashFunction = xidHashFunction
}

// Java: CompilerOptions.getXidHashFunction()
func (o *CompilerOptions) GetXidHashFunction() Xid_HashFunction {
	return o.xidHashFunction
}

// Sets the hash function to use for chunk ID generation
//
// Java: CompilerOptions.setChunkIdHashFunction(Xid.HashFunction)
func (o *CompilerOptions) SetChunkIdHashFunction(chunkIdHashFunction Xid_HashFunction) {
	o.chunkIdHashFunction = chunkIdHashFunction
}

// Java: CompilerOptions.getChunkIdHashFunction()
func (o *CompilerOptions) GetChunkIdHashFunction() Xid_HashFunction {
	return o.chunkIdHashFunction
}

// Sets the functions whose debug strings to replace.
//
// Java: CompilerOptions.setReplaceStringsConfiguration(String, List<String>)
func (o *CompilerOptions) SetReplaceStringsConfiguration(placeholderToken string, functionDescriptors []string) {
	o.replaceStringsPlaceholderToken = placeholderToken
	o.replaceStringsFunctionDescriptions = append([]string{}, functionDescriptors...)
}

// Java: CompilerOptions.setReplaceMessagesWithChromeI18n(boolean, String)
func (o *CompilerOptions) SetReplaceMessagesWithChromeI18n(replaceMessagesWithChromeI18n bool, tcProjectId *string) {
	if replaceMessagesWithChromeI18n && o.messageBundle != nil {
		if _, isEmpty := o.messageBundle.(*EmptyMessageBundle); !isEmpty {
			panic(jutil.NewRuntimeException(
				"When replacing messages with" +
					" chrome.i18n.getMessage, a message bundle should not be specified."))
		}
	}

	o.replaceMessagesWithChromeI18n = replaceMessagesWithChromeI18n
	o.tcProjectId = tcProjectId
}

// Java: CompilerOptions.getTcProjectId()
func (o *CompilerOptions) GetTcProjectId() *string {
	return o.tcProjectId
}

// Should we run the pass that does replacement of the chrome-specific `chrome.i18n.getMessage()`
// translatable message definitions?
//
// This form of l10n is incompatible with our standard `goog.getMsg()` messages.
//
// Java: CompilerOptions.shouldRunReplaceMessagesForChrome()
func (o *CompilerOptions) ShouldRunReplaceMessagesForChrome() bool {
	if o.replaceMessagesWithChromeI18n {
		_, isEmpty := o.messageBundle.(*EmptyMessageBundle)
		jutil.CheckStateMsg(
			o.messageBundle == nil || isEmpty,
			"When replacing messages with chrome.i18n.getMessage, a message bundle should not be"+
				" specified.")
		jutil.CheckStateMsg(
			!o.doLateLocalization, "Late localization is not supported for chrome.i18n.getMessage")
		return true
	} else {
		return false
	}
}

// Java: CompilerOptions.setAngularPass(boolean)
func (o *CompilerOptions) SetAngularPass(angularPass bool) {
	o.angularPass = angularPass
}

// Java: CompilerOptions.getAngularPass()
func (o *CompilerOptions) GetAngularPass() bool {
	return o.angularPass
}

// Java: CompilerOptions.setPolymerVersion(Integer)
func (o *CompilerOptions) SetPolymerVersion(polymerVersion *int32) {
	jutil.CheckArgumentMsg1(
		polymerVersion == nil || *polymerVersion == 1 || *polymerVersion == 2,
		"Invalid Polymer version: (%s)",
		polymerVersion)
	// Internally we model the polymerVersion as a boolean. if True, run the PolymerPass; otherwise
	// don't. The public API is an integer for historical reasons & in case we need to modify it
	// in the future.
	o.polymerPass = polymerVersion != nil
}

// Java: CompilerOptions.isPolymerPassEnabled()
func (o *CompilerOptions) IsPolymerPassEnabled() bool {
	return o.polymerPass
}

// Java: CompilerOptions.setChromePass(boolean)
func (o *CompilerOptions) SetChromePass(chromePass bool) {
	o.chromePass = chromePass
}

// Java: CompilerOptions.isChromePassEnabled()
func (o *CompilerOptions) IsChromePassEnabled() bool {
	return o.chromePass
}

// Java: CompilerOptions.setJ2clPass(J2clPassMode)
func (o *CompilerOptions) SetJ2clPass(j2clPassMode CompilerOptions_J2clPassMode) {
	o.j2clPassMode = j2clPassMode
}

// Java: CompilerOptions.getJ2clPass()
func (o *CompilerOptions) GetJ2clPass() CompilerOptions_J2clPassMode {
	return o.j2clPassMode
}

// Java: CompilerOptions.setJ2clMinifierEnabled(boolean)
func (o *CompilerOptions) SetJ2clMinifierEnabled(enabled bool) {
	o.j2clMinifierEnabled = enabled
}

// Java: CompilerOptions.isJ2clMinifierEnabled()
func (o *CompilerOptions) IsJ2clMinifierEnabled() bool {
	return o.j2clMinifierEnabled
}

// Java: CompilerOptions.setJ2clMinifierPruningManifest(String)
func (o *CompilerOptions) SetJ2clMinifierPruningManifest(j2clMinifierPruningManifest *string) {
	o.j2clMinifierPruningManifest = j2clMinifierPruningManifest
}

// Java: CompilerOptions.getJ2clMinifierPruningManifest()
func (o *CompilerOptions) GetJ2clMinifierPruningManifest() *string {
	return o.j2clMinifierPruningManifest
}

// Java: CompilerOptions.setPrintSourceAfterEachPass(boolean)
func (o *CompilerOptions) SetPrintSourceAfterEachPass(printSource bool) {
	o.printSourceAfterEachPass = printSource
}

// Java: CompilerOptions.shouldPrintSourceAfterEachPass()
func (o *CompilerOptions) ShouldPrintSourceAfterEachPass() bool {
	return o.printSourceAfterEachPass
}

// Java: CompilerOptions.setFilesToPrintAfterEachPassRegexList(List<String>)
func (o *CompilerOptions) SetFilesToPrintAfterEachPassRegexList(filePathRegexList []string) {
	o.filesToPrintAfterEachPassRegexList = filePathRegexList
}

// Java: CompilerOptions.setChunksToPrintAfterEachPassRegexList(List<String>)
func (o *CompilerOptions) SetChunksToPrintAfterEachPassRegexList(chunkPathRegexList []string) {
	o.chunksToPrintAfterEachPassRegexList = chunkPathRegexList
}

// Java: CompilerOptions.setQnameUsesToPrintAfterEachPassList(List<String>)
func (o *CompilerOptions) SetQnameUsesToPrintAfterEachPassList(qnameRegexList []string) {
	o.qnameUsesToPrintAfterEachPassList = qnameRegexList
}

// Java: CompilerOptions.getFilesToPrintAfterEachPassRegexList()
func (o *CompilerOptions) GetFilesToPrintAfterEachPassRegexList() []string {
	return o.filesToPrintAfterEachPassRegexList
}

// Java: CompilerOptions.getChunksToPrintAfterEachPassRegexList()
func (o *CompilerOptions) GetChunksToPrintAfterEachPassRegexList() []string {
	return o.chunksToPrintAfterEachPassRegexList
}

// Java: CompilerOptions.getQnameUsesToPrintAfterEachPassList()
func (o *CompilerOptions) GetQnameUsesToPrintAfterEachPassList() []string {
	return o.qnameUsesToPrintAfterEachPassList
}

// Java: CompilerOptions.getTracerMode()
func (o *CompilerOptions) GetTracerMode() CompilerOptions_TracerMode {
	return o.tracer
}

// NOTE: Timing information will not be printed if compiler.disableThreads() is called!
//
// Java: CompilerOptions.setTracerMode(TracerMode)
func (o *CompilerOptions) SetTracerMode(mode CompilerOptions_TracerMode) {
	o.tracer = mode
}

// Java: CompilerOptions.getTracerOutput()
func (o *CompilerOptions) GetTracerOutput() *jutil.Path {
	return o.tracerOutput
}

// Java: CompilerOptions.setTracerOutput(Path)
func (o *CompilerOptions) SetTracerOutput(out *jutil.Path) {
	o.tracerOutput = out
}
