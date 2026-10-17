// Port of com/google/javascript/jscomp/SourceMap.java @ closure-compiler v20261015

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

package jscomp

import (
	"io"
	"strings"

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
	prefixMappings           []SourceMap_LocationMapping
	sourceLocationFixupCache *jutil.LinkedHashMap[string, string]

	// A mapping derived from input source maps. Maps back to input sources
	// that inputs to this compilation job have been generated from, and used
	// to create a source map that maps all the way back to original inputs.
	// nil if no such mapping is wanted.
	mapping SourceFileMapping
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
func (l SourceMap_DetailLevel) Apply(node *rhino.Node) bool {
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

// SourceMap_LocationMapping is a function that mape a "destination" location
// to use within the source map. Should return nil if the value is not
// mapped.
type SourceMap_LocationMapping interface {
	// Map: location is the location to transform; the result is the
	// transformed location or nil if not transformed
	//
	// Java: SourceMap.LocationMapping.map(String)
	Map(location string) *string
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
func (m *SourceMap_PrefixLocationMapping) Map(location string) *string {
	if strings.HasPrefix(location, m.prefix) {
		s := m.replacement + location[len(m.prefix):]
		return &s
	}
	return nil
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
	return jutil.ObjectsHash(m.prefix, m.replacement)
}

// Java: SourceMap(SourceMapGenerator)
func newSourceMap(generator sourcemap.SourceMapGenerator) *SourceMap {
	return &SourceMap{
		generator:                generator,
		prefixMappings:           []SourceMap_LocationMapping{},
		sourceLocationFixupCache: jutil.NewLinkedHashMap[string, string](),
	}
}

// SourceMap_Mapping maintains a mapping from a given node to the position in
// the source code at which its generated form was placed. The positions are
// typically relative to the source file the node is located in, but might be
// adjusted if that source is being concatenated to other sources.
type SourceMap_Mapping struct {
	node  *rhino.Node
	start *sourcemap.FilePosition
	end   *sourcemap.FilePosition
}

// Java: SourceMap.Mapping()
func NewSourceMap_Mapping() *SourceMap_Mapping {
	return &SourceMap_Mapping{}
}

// Java: SourceMap.Mapping.toString()
func (m *SourceMap_Mapping) String() string {
	// This toString() representation is used for debugging purposes only.
	return "Mapping: start " + jutil.ToJavaString(m.start) + ", end " + jutil.ToJavaString(m.end) +
		", node " + jutil.ToJavaString(m.node)
}

// AddMapping adds the mapping.
//
// Java: SourceMap.addMapping(Mapping)
func (s *SourceMap) AddMapping(mapping *SourceMap_Mapping) {
	s.AddMappingOfNode(mapping.node, mapping.start, mapping.end)
}

// Java: SourceMap.addMapping(Node, FilePosition, FilePosition)
func (s *SourceMap) AddMappingOfNode(
	node *rhino.Node, outputStartPosition *sourcemap.FilePosition, outputEndPosition *sourcemap.FilePosition) {
	// If the node does not have an associated source file or
	// its line number is -1, then the node does not have sufficient
	// information for a mapping to be useful.
	sourceFile := node.GetStaticSourceFile()
	if sourceFile == nil || node.GetLineno() < 0 {
		return
	}

	sourceFileName := sourceFile.GetName()
	lineNo := node.GetLineno()
	charNo := node.GetCharno()
	originalName := sourceMap_getOriginalName(node)

	if s.mapping != nil {
		sourceMapping := s.mapping.GetSourceMapping(sourceFileName, lineNo, charNo)
		if sourceMapping == nil {
			// The source file does not have a input map. We consider this to be an
			// original source range and include it in the output map.
		} else if sourceMapping.Equals(sourcemap.Mapping_OriginalMapping_GetDefaultInstance()) {
			// The source file does have an input map, but it does not map the
			// location. We consider this to be a synthetic code range and do not
			// include it in the output map.
			// TODO b/452676030 - Report the sourceless mapping.
			return
		} else {
			// The source file mapped our code range to its original source location.
			sourceFileName = sourceMapping.GetOriginalFile()
			lineNo = sourceMapping.GetLineNumber()
			charNo = sourceMapping.GetColumnPosition()
			identifier := sourceMapping.GetIdentifier()
			if sourceMapping.HasIdentifier() && identifier != "" {
				originalName = &identifier
			}
		}
	}

	sourceFileName = s.fixupSourceLocation(sourceFileName)

	// Rhino source lines are one based but for v3 source maps, we make
	// them zero based.
	lineBaseOffset := 1

	s.generator.AddMapping(
		&sourceFileName,
		originalName,
		sourcemap.NewFilePosition(lineNo-lineBaseOffset, charNo),
		outputStartPosition,
		outputEndPosition)
}

// Java: SourceMap.addSourceFile(String, String)
func (s *SourceMap) AddSourceFile(name string, code *string) {
	s.generator.AddSourcesContent(s.fixupSourceLocation(name), code)
}

// Java: SourceMap.getOriginalName(Node)
func sourceMap_getOriginalName(node *rhino.Node) *string {
	if node.GetOriginalName() != nil {
		return node.GetOriginalName()
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
		if parent != nil && parent.IsFunction() && parent.GetFirstChild() == node {
			return parent.GetOriginalName()
		}
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

	if fixed, ok := s.sourceLocationFixupCache.Lookup(sourceFile); ok {
		return fixed
	}

	// Replace the first prefix found with its replacement
	var fixed *string
	for _, mapping := range s.prefixMappings {
		fixed = mapping.Map(sourceFile)
		if fixed != nil {
			break
		}
	}

	// If none of the mappings match then use the original file path.
	if fixed == nil {
		fixed = &sourceFile
	}

	s.sourceLocationFixupCache.Put(sourceFile, *fixed)
	return *fixed
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
	s.sourceLocationFixupCache.Clear()
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
}

// Java: SourceMap.setSourceFileMapping(SourceFileMapping)
func (s *SourceMap) SetSourceFileMapping(mapping SourceFileMapping) {
	s.mapping = mapping
}
