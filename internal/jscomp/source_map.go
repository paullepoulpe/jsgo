// Port of com/google/javascript/jscomp/SourceMap.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/SourceMapInput.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/SourceMapResolver.java @ closure-compiler v20261021

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
	"encoding/base64"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// SourceMap collects information mapping the generated (compiled) source
// back to its original source for debugging purposes.
//
// See CodeConsumer, CodeGenerator, CodePrinter.
type SourceMap struct {
	generator                sourcemap.SourceMapGenerator
	generatorV3              *sourcemap.SourceMapGeneratorV3
	prefixMappings           []SourceMap_LocationMapping
	sourceLocationFixupCache map[string]string
	lastSourceFile           rhino.StaticSourceFile
	lastFixedSourceFileName  string

	// A mapping derived from input source maps. Maps back to input sources
	// that inputs to this compilation job have been generated from, and used
	// to create a source map that maps all the way back to original inputs.
	// nil if no such mapping is wanted.
	compilerMapping   *Compiler
	resolvedSourceMap compiler_ResolvedSourceMap
	lineCache         sourcemap.SourceMapLineCache
}

// SourceMap_Format is an enumeration of available source map formats
type SourceMap_Format int32

const (
	SourceMap_Format_DEFAULT SourceMap_Format = iota
	SourceMap_Format_V3
)

var sourceMap_Format_names = [...]string{"DEFAULT", "V3"}

// SourceMap_Format_values is Format.values().
func SourceMap_Format_values() []SourceMap_Format {
	return []SourceMap_Format{SourceMap_Format_DEFAULT, SourceMap_Format_V3}
}

// String is Enum.name().
func (f SourceMap_Format) String() string { return sourceMap_Format_names[f] }

// Ordinal is Enum.ordinal().
func (f SourceMap_Format) Ordinal() int { return int(f) }

// Java: SourceMap.Format.getInstance()
func (f SourceMap_Format) GetInstance() *SourceMap {
	switch f {
	case SourceMap_Format_DEFAULT:
		return newSourceMap(sourcemap.SourceMapGeneratorFactory_GetInstance(sourcemap.SourceMapFormat_DEFAULT))
	case SourceMap_Format_V3:
		return newSourceMap(sourcemap.SourceMapGeneratorFactory_GetInstance(sourcemap.SourceMapFormat_V3))
	}
	panic(jutil.NewIllegalStateException())
}

// SourceMap_DetailLevel: source maps can be very large different levels of
// detail can be specified.
type SourceMap_DetailLevel int32

const (
	// ALL is best when the fullest details are needed for debugging or for
	// code-origin analysis.
	SourceMap_DetailLevel_ALL SourceMap_DetailLevel = iota
	// SYMBOLS is intended to be used for stack trace deobfuscation when full
	// detail is not needed.
	SourceMap_DetailLevel_SYMBOLS
)

var sourceMap_DetailLevel_names = [...]string{"ALL", "SYMBOLS"}

// SourceMap_DetailLevel_values is DetailLevel.values().
func SourceMap_DetailLevel_values() []SourceMap_DetailLevel {
	return []SourceMap_DetailLevel{SourceMap_DetailLevel_ALL, SourceMap_DetailLevel_SYMBOLS}
}

// String is Enum.name().
func (l SourceMap_DetailLevel) String() string { return sourceMap_DetailLevel_names[l] }

// Ordinal is Enum.ordinal().
func (l SourceMap_DetailLevel) Ordinal() int { return int(l) }

// Apply is Predicate<Node>.apply(Node).
//
// Java: SourceMap.DetailLevel.apply(Node)
func (l SourceMap_DetailLevel) Apply(node rhino.Node) bool {
	switch l {
	case SourceMap_DetailLevel_ALL:
		return true
	case SourceMap_DetailLevel_SYMBOLS:
		return node.IsCall() ||
			node.IsNew() ||
			node.IsFunction() ||
			node.IsName() ||
			NodeUtil_IsNormalOrOptChainGet(node) ||
			NodeUtil_MayBeObjectLitKey(node) ||
			node.IsTaggedTemplateLit()
	}
	panic(jutil.NewIllegalStateException())
}

// SourceMap_LocationMapping is a function that maps a "destination" location
// to use within the source map. Returns ("", false) if the value is not
// mapped.
type SourceMap_LocationMapping interface {
	// Map: location is the location to transform; returns the transformed
	// location and true, or ("", false) if not transformed.
	//
	// Java: SourceMap.LocationMapping.map(String)
	Map(location string) (string, bool)
}

// SourceMap_PrefixLocationMapping is a simple LocationMapping that strips a
// prefix from a location.
type SourceMap_PrefixLocationMapping struct {
	prefix      string
	replacement string
}

var _ SourceMap_LocationMapping = (*SourceMap_PrefixLocationMapping)(nil)

// Java: SourceMap.PrefixLocationMapping(String, String)
func NewSourceMap_PrefixLocationMapping(prefix string, replacement string) *SourceMap_PrefixLocationMapping {
	return &SourceMap_PrefixLocationMapping{prefix: prefix, replacement: replacement}
}

// Java: SourceMap.PrefixLocationMapping.map(String)
func (m *SourceMap_PrefixLocationMapping) Map(location string) (string, bool) {
	if strings.HasPrefix(location, m.prefix) {
		return m.replacement + location[len(m.prefix):], true
	}
	return "", false
}

// Java: SourceMap.PrefixLocationMapping.toString()
func (m *SourceMap_PrefixLocationMapping) String() string {
	return "(" + m.prefix + "|" + m.replacement + ")"
}

// Java: SourceMap.PrefixLocationMapping.equals(Object)
func (m *SourceMap_PrefixLocationMapping) Equals(other any) bool {
	if prefixLocationMapping, ok := other.(*SourceMap_PrefixLocationMapping); ok && prefixLocationMapping != nil {
		return prefixLocationMapping.prefix == m.prefix &&
			prefixLocationMapping.replacement == m.replacement
	} else {
		return false
	}
}

// Java: SourceMap.PrefixLocationMapping.hashCode()
func (m *SourceMap_PrefixLocationMapping) HashCode() int32 {
	return 31*(31+jutil.StringHashCode(m.prefix)) + jutil.StringHashCode(m.replacement)
}

// Java: SourceMap(SourceMapGenerator)
func newSourceMap(generator sourcemap.SourceMapGenerator) *SourceMap {
	v3, _ := generator.(*sourcemap.SourceMapGeneratorV3)
	return &SourceMap{
		generator:   generator,
		generatorV3: v3,
	}
}

// SourceMap_Mapping maintains a mapping from a given node to the position in
// the source code at which its generated form was placed. The positions are
// typically relative to the source file the node is located in, but might be
// adjusted if that source is being concatenated to other sources.
type SourceMap_Mapping struct {
	node  rhino.Node
	start sourcemap.FilePosition
	end   sourcemap.FilePosition
}

// Java: SourceMap.Mapping.toString()
func (m *SourceMap_Mapping) String() string {
	// This toString() representation is used for debugging purposes only.
	return "Mapping: start " + jutil.ToJavaString(&m.start) + ", end " + jutil.ToJavaString(&m.end) +
		", node " + jutil.ToJavaString(m.node)
}

// AddMapping adds the mapping.
//
// Java: SourceMap.addMapping(Mapping)
func (s *SourceMap) AddMapping(mapping *SourceMap_Mapping) {
	s.AddMappingOfNodePositions(mapping.node, mapping.start, mapping.end)
}

// sourceMap_resolvedMapping is a mapping whose node has been resolved to
// its source file name, original name and source position (Go-only).
type sourceMap_resolvedMapping struct {
	sourceFileName string
	originalName   *string
	srcPos         sourcemap.FilePosition
	start, end     sourcemap.FilePosition
}

func (s *SourceMap) AddMappingOfNodePositions(
	node rhino.Node, outputStartPosition sourcemap.FilePosition, outputEndPosition sourcemap.FilePosition) {
	if r, ok := s.resolveMappingOfNodePositions(node, outputStartPosition, outputEndPosition); ok {
		s.addResolvedMapping(&r)
	}
}

// resolveMappingOfNodePositions is the node-dependent half of
// AddMappingOfNodePositions (Go-only split: fast mode runs it on the print
// workers, each with its own SourceMap and caches).
func (s *SourceMap) resolveMappingOfNodePositions(
	node rhino.Node, outputStartPosition sourcemap.FilePosition, outputEndPosition sourcemap.FilePosition,
) (r sourceMap_resolvedMapping, ok bool) {
	// If the node does not have an associated source file or
	// its line number is -1, then the node does not have sufficient
	// information for a mapping to be useful.
	sourceFile := node.GetStaticSourceFile()
	if sourceFile == nil || node.GetLineno() < 0 {
		return r, false
	}

	lineNo := node.GetLineno()
	charNo := node.GetCharno()
	originalName := sourceMap_getOriginalName(node)

	var sourceFileName string
	if s.compilerMapping == nil {
		if sourceFile == s.lastSourceFile {
			sourceFileName = s.lastFixedSourceFileName
		} else {
			sourceFileName = s.fixupSourceLocation(sourceFile.GetName())
			s.lastSourceFile = sourceFile
			s.lastFixedSourceFileName = sourceFileName
		}
	} else {
		sourceFileName = sourceFile.GetName()
		sourceMapping, identifierPtr, hasInputMap, isDefaultEmpty := s.compilerMapping.getSourceMappingValueCached(
			sourceFileName, lineNo, charNo, &s.resolvedSourceMap, &s.lineCache)
		if !hasInputMap {
			// The source file does not have a input map. We consider this to be an
			// original source range and include it in the output map.
		} else if isDefaultEmpty {
			// The source file does have an input map, but it does not map the
			// location. We consider this to be a synthetic code range and do not
			// include it in the output map.
			// TODO b/452676030 - Report the sourceless mapping.
			return r, false
		} else {
			// The source file mapped our code range to its original source location.
			sourceFileName = sourceMapping.GetOriginalFile()
			lineNo = sourceMapping.GetLineNumber()
			charNo = sourceMapping.GetColumnPosition()
			if identifierPtr != nil {
				originalName = identifierPtr
			}
		}
		sourceFileName = s.fixupSourceLocation(sourceFileName)
	}

	// Rhino source lines are one based but for v3 source maps, we make
	// them zero based.
	lineBaseOffset := 1

	srcPos := sourcemap.MakeFilePosition(lineNo-lineBaseOffset, charNo)
	return sourceMap_resolvedMapping{sourceFileName, originalName, srcPos, outputStartPosition, outputEndPosition}, true
}

// addResolvedMapping is the generator half of AddMappingOfNodePositions.
func (s *SourceMap) addResolvedMapping(r *sourceMap_resolvedMapping) {
	if s.generatorV3 != nil {
		s.generatorV3.AddMappingPositions(r.sourceFileName, r.originalName, r.srcPos, r.start, r.end)
	} else {
		s.addMappingFallback(r.sourceFileName, r.originalName, r.srcPos, r.start, r.end)
	}
}

//go:noinline
func (s *SourceMap) addMappingFallback(
	sourceFileName string,
	originalName *string,
	srcPos sourcemap.FilePosition,
	outputStartPosition sourcemap.FilePosition,
	outputEndPosition sourcemap.FilePosition) {
	s.generator.AddMapping(
		&sourceFileName,
		originalName,
		&srcPos,
		&outputStartPosition,
		&outputEndPosition)
}

// Java: SourceMap.addSourceFile(String, String)
func (s *SourceMap) AddSourceFile(name string, code *string) {
	s.generator.AddSourcesContent(s.fixupSourceLocation(name), code)
}

// Java: SourceMap.getOriginalName(Node)
func sourceMap_getOriginalName(node rhino.Node) *string {
	if orig := node.GetOriginalName(); orig != nil {
		return orig
	}
	// If this is the name identifier of a function declaration/expression and the name node itself
	// has no originalName, fall back to the enclosing FUNCTION node's originalName.
	// This covers functions that were originally anonymous (class methods compiled to prototype
	// assignments, `var f = function(){}` patterns): SourceInformationAnnotator skips empty-string
	// NAME nodes but always annotates the FUNCTION node. When the renaming pass later fills in the
	// NAME, the node's originalName remains null. The FUNCTION node's originalName is the correct
	// original name to use.
	if node.IsName() {
		parent := node.GetParent()
		if !parent.IsNil() && parent.IsFunction() && parent.GetFirstChild() == node {
			return parent.GetOriginalName()
		}
		return nil
	}
	if node.IsMemberFunctionDef() {
		return node.GetFirstChild().GetOriginalName()
	}
	return nil
}

// fixupSourceLocation: sourceFile is the source file location to fixup;
// the result is a remapped source file.
//
// Java: SourceMap.fixupSourceLocation(String)
func (s *SourceMap) fixupSourceLocation(sourceFile string) string {
	if len(s.prefixMappings) == 0 {
		return sourceFile
	}

	if fixed, ok := s.sourceLocationFixupCache[sourceFile]; ok {
		return fixed
	}

	// Replace the first prefix found with its replacement.
	// If none of the mappings match then use the original file path.
	res := sourceFile
	for _, mapping := range s.prefixMappings {
		if fixed, ok := mapping.Map(sourceFile); ok {
			res = fixed
			break
		}
	}

	if s.sourceLocationFixupCache == nil {
		s.sourceLocationFixupCache = make(map[string]string)
	}
	s.sourceLocationFixupCache[sourceFile] = res
	return res
}

// AppendTo appends the source map to out. It panics with an IOException if
// out does.
//
// Java: SourceMap.appendTo(Appendable, String)
func (s *SourceMap) AppendTo(out io.StringWriter, name *string) {
	// Go-only: language-adapter — a null name is passed through unchanged
	// when there are no prefix mappings, as fixupSourceLocation(null) does
	// in Java.
	var fixed *string
	if name != nil || len(s.prefixMappings) != 0 {
		f := s.fixupSourceLocation(*name)
		fixed = &f
	}
	s.generator.AppendTo(out, fixed)
}

// Java: SourceMap.reset()
func (s *SourceMap) Reset() {
	s.generator.Reset()
	clear(s.sourceLocationFixupCache)
	s.lastSourceFile = nil
	s.lastFixedSourceFileName = ""
	s.resolvedSourceMap = compiler_ResolvedSourceMap{}
	s.lineCache.Reset()
}

// Java: SourceMap.setStartingPosition(int, int)
func (s *SourceMap) SetStartingPosition(offsetLine int, offsetIndex int) {
	s.generator.SetStartingPosition(offsetLine, offsetIndex)
}

// Java: SourceMap.setWrapperPrefix(String)
func (s *SourceMap) SetWrapperPrefix(prefix string) {
	s.generator.SetWrapperPrefix(prefix)
}

// Java: SourceMap.validate(boolean)
func (s *SourceMap) Validate(validate bool) {
	s.generator.Validate(validate)
}

// Java: SourceMap.setPrefixMappings(List<? extends LocationMapping>)
func (s *SourceMap) SetPrefixMappings(sourceMapLocationMappings []SourceMap_LocationMapping) {
	s.prefixMappings = sourceMapLocationMappings
	clear(s.sourceLocationFixupCache)
	s.lastSourceFile = nil
	s.lastFixedSourceFileName = ""
}

// SetCompilerSourceMapping sets the Compiler used to resolve input source maps.
func (s *SourceMap) SetCompilerSourceMapping(c *Compiler) {
	s.compilerMapping = c
}

// SourceMapInput is a lazy-loaded SourceMapConsumerV3 instance.
type SourceMapInput struct {
	mu              sync.Mutex
	sourceFile      *SourceFile
	parsedSourceMap *sourcemap.SourceMapConsumerV3
	cached          bool
	ready           atomic.Bool
	contentDropped  atomic.Bool
}

var SourceMapInput_SOURCEMAP_RESOLVE_FAILED = DiagnosticType_Warning("SOURCEMAP_RESOLVE_FAILED", "Failed to resolve sourcemap at {0}: {1}")

var SourceMapInput_SOURCEMAP_PARSE_FAILED = DiagnosticType_Warning(
	"SOURCEMAP_PARSE_FAILED", "Failed to parse malformed sourcemap in {0}: {1}")

// Java: SourceMapInput(SourceFile)
func NewSourceMapInput(sourceFile *SourceFile) *SourceMapInput {
	return &SourceMapInput{sourceFile: sourceFile}
}

// GetSourceMap gets the source map, reading from disk and parsing if
// necessary. Returns nil if the sourcemap cannot be resolved or is
// malformed.
//
// Java: SourceMapInput.getSourceMap(ErrorManager)
func (i *SourceMapInput) GetSourceMap(errorManager ErrorManager) *sourcemap.SourceMapConsumerV3 {
	if i.ready.Load() {
		return i.parsedSourceMap
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	if !i.cached {
		// Avoid re-reading or reparsing files.
		i.cached = true
		sourceMapPath := i.sourceFile.GetName()
		sourceMapContents, ioErr := i.sourceFile.getCode()
		if ioErr != nil {
			error := JSError_Make(SourceMapInput_SOURCEMAP_RESOLVE_FAILED, sourceMapPath, ioErr.GetMessage())
			errorManager.Report(error.DefaultLevel().Ptr(), error)
		} else {
			consumer := sourcemap.NewSourceMapConsumerV3()
			if parseErr := consumer.TryParse(sourceMapContents); parseErr != nil {
				error := JSError_Make(SourceMapInput_SOURCEMAP_PARSE_FAILED, sourceMapPath, parseErr.GetMessage())
				errorManager.Report(error.DefaultLevel().Ptr(), error)
			} else {
				// Detect an empty "default" source map provided by Gulp. We do not want
				// to treat this file as having a real source map (implying this is an
				// intermediate source file). Rather, this is the original source file.
				// https://github.com/gulp-sourcemaps/gulp-sourcemaps/blob/91b94954/src/init/index.js#L116-L124
				if file := consumer.GetFile(); file != nil &&
					*file != "" &&
					len(consumer.GetOriginalSources()) == 1 &&
					sourceMapInput_contains(consumer.GetOriginalSources(), *file) &&
					len(consumer.GetOriginalSourcesContent()) == 1 &&
					len(consumer.GetOriginalNames()) == 0 &&
					consumer.GetLineCount() == -1 {
					consumer = nil
				}

				i.parsedSourceMap = consumer
				// Go-only: performance — the parsed map holds all that is
				// needed; do not also keep the map's text (read again on demand).
				i.sourceFile.ClearCachedSource()
			}
		}
		i.ready.Store(true)
	}
	return i.parsedSourceMap
}

func (i *SourceMapInput) dropSourcesContentOnce(consumer *sourcemap.SourceMapConsumerV3) {
	if i.contentDropped.Load() {
		return
	}
	i.mu.Lock()
	if !i.contentDropped.Load() {
		consumer.DropSourcesContent()
		i.contentDropped.Store(true)
	}
	i.mu.Unlock()
}

// sourceMapInput_contains reports whether list contains a non-nil string equal to target.
func sourceMapInput_contains(list []*string, target string) bool {
	for _, s := range list {
		if s != nil && *s == target {
			return true
		}
	}
	return false
}

// GetOriginalPath gets the original location of this sourcemap file on
// disk.
//
// Java: SourceMapInput.getOriginalPath()
func (i *SourceMapInput) GetOriginalPath() string {
	return i.sourceFile.GetName()
}

// GetRawSourceMapContents may return nil.
//
// Java: SourceMapInput.getRawSourceMapContents()
func (i *SourceMapInput) GetRawSourceMapContents() *string {
	code, err := i.sourceFile.getCode()
	if err != nil {
		return nil
	}
	return &code
}

const sourceMapResolver_BASE64_URL_PREFIX = "data:"
const sourceMapResolver_BASE64_START = "base64,"

// For now only accept UTF-8. We could use the actual charset information in
// the future. Default is UTF-8 if unspecified.
func sourceMapResolver_isAcceptedMediaType(mediaType string) bool {
	return mediaType == "application/json;charset=utf-8;" || mediaType == "application/json;"
}

// SourceMapResolver_ExtractSourceMap locates the appropriate sourcemap on
// disk for a given //# sourceMappingUrl. This is use for sourcemap merging
// (--apply_input_source_maps) and for error resolution. parseInlineSourceMaps
// says whether to parse Base64 encoded source maps. It may return nil.
//
// Java: SourceMapResolver.extractSourceMap(SourceFile, String, boolean)
func SourceMapResolver_ExtractSourceMap(jsFile *SourceFile, sourceMapURL string, parseInlineSourceMaps bool) *SourceFile {
	if parseInlineSourceMaps && strings.HasPrefix(sourceMapURL, sourceMapResolver_BASE64_URL_PREFIX) {
		if rawBase64, canonicalUTF8, ok := sourceMapResolver_extractRawBase64(sourceMapURL); ok {
			return newInlineBase64SourceFile(
				jsFile.GetName()+".inline.map", rawBase64, canonicalUTF8)
		}
		return nil
	}
	// TODO(tdeegan): Handle absolute urls here.  The compiler needs to come up with a scheme for
	// properly resolving absolute urls from http:// or the root /some/abs/path/... See b/62544959.
	if sourceMapResolver_isAbsolute(sourceMapURL) {
		return nil
	}
	// If not absolute, its relative.
	// TODO(tdeegan): Handle urls relative to //# sourceURL. See the sourcemap spec.
	return SourceMapResolver_GetRelativePath(jsFile.GetName(), sourceMapURL)
}

func sourceMapResolver_extractRawBase64(url string) (rawBase64 string, canonicalUTF8 bool, ok bool) {
	if strings.HasPrefix(url, sourceMapResolver_BASE64_URL_PREFIX) {
		base64StartIndex := strings.Index(url, sourceMapResolver_BASE64_START)
		if base64StartIndex >= 0 {
			mediaType := url[len(sourceMapResolver_BASE64_URL_PREFIX):base64StartIndex]
			if sourceMapResolver_isAcceptedMediaType(mediaType) {
				raw := url[base64StartIndex+len(sourceMapResolver_BASE64_START):]
				canonical := validateAndCheckCanonicalBase64UTF8(raw)
				return strings.Clone(raw), canonical, true
			}
		}
	}
	return "", false, false
}

func decodeInlineBase64(chars string) []byte {
	trimmed := strings.TrimRight(chars, "=")
	if strings.IndexAny(trimmed, "\r\n") >= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Unrecognized newline in base64"))
	}
	data, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(trimmed)
	if err != nil {
		panic(jutil.NewIllegalArgumentExceptionMsg(err.Error()))
	}
	return data
}

func validateAndCheckCanonicalBase64UTF8(chars string) bool {
	trimmed := strings.TrimRight(chars, "=")
	if strings.IndexAny(trimmed, "\r\n") >= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Unrecognized newline in base64"))
	}
	pad := len(chars) - len(trimmed)
	// "77u/" decodes to the 3-byte UTF-8 BOM (0xEF, 0xBB, 0xBF), which
	// SourceFile.SetCodeAndDoBookkeeping strips on load.
	canonical := len(chars)%4 == 0 && pad <= 2 && !strings.HasPrefix(chars, "77u/")

	var srcBuf [1024]byte
	var dstBuf [772]byte
	carry := 0
	for i := 0; i < len(trimmed); i += 1024 {
		chunk := trimmed[i:min(i+1024, len(trimmed))]
		copy(srcBuf[:len(chunk)], chunk)
		n, err := base64.StdEncoding.WithPadding(base64.NoPadding).Decode(dstBuf[carry:], srcBuf[:len(chunk)])
		if err != nil {
			panic(jutil.NewIllegalArgumentExceptionMsg(err.Error()))
		}
		if canonical {
			if rem := len(chunk) % 4; rem != 0 {
				var tail [3]byte
				if _, strictErr := base64.StdEncoding.WithPadding(base64.NoPadding).Strict().Decode(tail[:], srcBuf[len(chunk)-rem:len(chunk)]); strictErr != nil {
					canonical = false
				}
			}
		}
		if !canonical {
			continue
		}
		total := carry + n
		final := i+1024 >= len(trimmed)
		if final || total == 0 || dstBuf[total-1] < 0x80 {
			if !utf8.Valid(dstBuf[:total]) {
				canonical = false
			}
			carry = 0
			continue
		}
		split := total - 1
		for split > 0 && total-split < 4 && (dstBuf[split]&0xC0) == 0x80 {
			split--
		}
		if !utf8.FullRune(dstBuf[split:total]) {
			if !utf8.Valid(dstBuf[:split]) {
				canonical = false
				carry = 0
				continue
			}
			carry = copy(dstBuf[:], dstBuf[split:total])
			continue
		}
		if !utf8.Valid(dstBuf[:total]) {
			canonical = false
		}
		carry = 0
	}
	if canonical && carry > 0 {
		canonical = false
	}
	return canonical
}

// Java: SourceMapResolver.isAbsolute(String)
func sourceMapResolver_isAbsolute(url string) bool {
	u, err := jutil.ParseURI(url)
	if err != nil {
		cause, _ := err.(jutil.Throwable)
		panic(jutil.NewRuntimeExceptionWithMessageAndCause("Sourcemap url was invalid: "+url, cause))
	}
	return u.IsAbsolute() || strings.HasPrefix(url, "/")
}

// SourceMapResolver_GetRelativePath returns the relative path, resolved
// relative to the base path, where the base path is interpreted as a
// filename rather than a directory. E.g.: getRelativeTo("../foo/bar.js",
// "baz/bam/qux.js") --> "baz/foo/bar.js"
//
// Java: SourceMapResolver.getRelativePath(String, String)
func SourceMapResolver_GetRelativePath(baseFilePath string, relativePath string) *SourceFile {
	return SourceFile_Builder_().
		WithPathFromPath(
			jutil.Path_Of(baseFilePath).ResolveSiblingString(relativePath).Normalize()).
		WithKind(rhino.StaticSourceFile_SourceKind_NON_CODE).
		Build()
}

// SourceMapResolver_AddBase64PrefixToEncodedSourceMap returns the encoded
// source map with the base64 prefix attached. (i.e. prefix with
// "data:application/json;base64,")
//
// E.g. turns "eyJ2ZXJzaW9uI..." into "data:application/json;base64,eyJ2ZXJzaW9uI..."
//
// Java: SourceMapResolver.addBase64PrefixToEncodedSourceMap(String)
func SourceMapResolver_AddBase64PrefixToEncodedSourceMap(encoded string) string {
	return sourceMapResolver_BASE64_URL_PREFIX + "application/json;" + sourceMapResolver_BASE64_START + encoded
}
