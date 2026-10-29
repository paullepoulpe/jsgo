// Port of com/google/javascript/jscomp/CodeConsumer.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"compress/gzip"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

type codeSinkMode uint8

const (
	sinkCompact codeSinkMode = iota
	sinkPretty
	sinkCost
	sinkSize
)

// CodeConsumer is the unified concrete sink of CodeGenerator output across
// compact printing, pretty printing, inline cost estimation, and performance
// tracker code-size estimation.
//
// Go-only: performance & architecture — flattens Java's 4-level CodeConsumer ->
// MappedCodePrinter -> CompactCodePrinter / PrettyCodePrinter / CostEstimator /
// CodeSizeEstimator hierarchy into one concrete struct with a mode tag so all
// CodeGenerator calls dispatch directly without interface or self-pointer indirection.
type CodeConsumer struct {
	mode                codeSinkMode
	statementNeedsEnded bool
	statementStarted    bool
	sawFunction         bool
	continueProcessing  bool
	lineBreak           bool
	createSrcMap        bool
	trackGzSize         bool
	lastChar            uint16

	// State tracking for template literals. Remember that template literal
	// substitutions can contain additional template literals.
	templateLitDepth    int
	templateLitSubDepth int

	// Printer buffer and line state (sinkCompact, sinkPretty).
	code                []byte
	hasNonASCII         bool
	lineLengthThreshold int
	lineLength          int
	lineIndex           int
	lineStartByte       int // sinkCompact
	preferredBreakByte  int // sinkCompact
	preferredBreakCol   int // sinkCompact
	indent              int // sinkPretty

	// Cost / size counters (sinkCost, sinkSize).
	cost    int
	maxCost int

	// Source map and license tracking (sinkCompact, sinkPretty).
	sourceMapDetailLevel SourceMap_DetailLevel
	licenseTracker       CodePrinter_LicenseTracker
	allMappings          []SourceMap_Mapping
	openMappings         []int
	completeMappings     []int
	cutIdx               int
	cutStackDepth        int

	// Pretty-printer original number lookup cache (sinkPretty).
	prettySrc *prettyNumberSource

	// Gzip size tracking (sinkSize when trackGzSize == true).
	gzStream  *gzip.Writer
	gzScratch []byte
}

// StartSourceMapping starts the source mapping for the given node at the
// current position.
//
// Java: CodeConsumer.startSourceMapping(Node)
func (c *CodeConsumer) StartSourceMapping(node rhino.Node) {
	if c.mode > sinkPretty {
		return
	}
	jutil.CheckState(!node.IsNil())
	if c.createSrcMap &&
		node.GetLineno() > 0 &&
		(c.sourceMapDetailLevel == SourceMap_DetailLevel_ALL || c.sourceMapDetailLevel.Apply(node)) &&
		node.GetStaticSourceFile() != nil {
		line := c.lineIndex
		index := c.lineLength
		jutil.CheckState(line >= 0)
		idx := len(c.allMappings)
		c.allMappings = append(c.allMappings, SourceMap_Mapping{
			node:  node,
			start: sourcemap.MakeFilePosition(line, index),
			end:   sourcemap.MakeFilePosition(-1, -1),
		})
		c.openMappings = append(c.openMappings, idx)
	}
}

// EndSourceMapping finishes the source mapping for the given node at the
// current position.
//
// Java: CodeConsumer.endSourceMapping(Node)
func (c *CodeConsumer) EndSourceMapping(node rhino.Node) {
	if c.createSrcMap && len(c.openMappings) != 0 {
		topIdx := c.openMappings[len(c.openMappings)-1]
		if c.allMappings[topIdx].node == node {
			c.openMappings = c.openMappings[:len(c.openMappings)-1]
			line := c.lineIndex
			index := c.lineLength
			jutil.CheckState(line >= 0)
			c.allMappings[topIdx].end = sourcemap.MakeFilePosition(line, index)
			if len(c.openMappings) < c.cutStackDepth {
				c.cutStackDepth = len(c.openMappings)
				c.completeMappings = append(c.completeMappings, topIdx)
			}
		}
	}
}

// TrackLicenses indicates to the CodeConsumer that this Node might carry
// licensing information, and allows the code consumer to manage licenses
// as it sees fit.
//
// Java: CodeConsumer.trackLicenses(Node)
func (c *CodeConsumer) TrackLicenses(node rhino.Node) {
	if c.licenseTracker != nil {
		c.licenseTracker.TrackLicensesForNode(node)
	}
}

// ContinueProcessing provides a means of interrupting the CodeGenerator.
//
// Java: CodeConsumer.continueProcessing()
func (c *CodeConsumer) ContinueProcessing() bool {
	return c.continueProcessing
}

// Java: CodeConsumer.getLastChar()
func (c *CodeConsumer) GetLastChar() uint16 {
	return c.lastChar
}

// codeConsumer_firstChar is s.charAt(0).
func codeConsumer_firstChar(s string) uint16 {
	if len(s) > 0 && s[0] < 0x80 {
		return uint16(s[0])
	}
	return jutil.StringCharAt(s, 0)
}

func (c *CodeConsumer) maybeIndent() {
	if c.lineLength == 0 && !c.IsInTemplateLiteral() {
		for i := 0; i < c.indent; i++ {
			c.codeAppendASCII(CodePrinter_PrettyCodePrinter_INDENT)
			c.lineLength += len(CodePrinter_PrettyCodePrinter_INDENT)
		}
	}
}

func (c *CodeConsumer) codeAppend(str string) int {
	if len(str) == 0 {
		return 0
	}
	if jutil.IsASCII(str) {
		c.code = append(c.code, str...)
		c.lastChar = uint16(str[len(str)-1])
		return len(str)
	}
	c.hasNonASCII = true
	if c.lastChar >= 0xD800 && c.lastChar <= 0xDBFF &&
		len(c.code) >= 3 &&
		len(str) >= 3 && str[0] == 0xED && str[1] >= 0xB0 && str[1] <= 0xBF && str[2]&0xC0 == 0x80 &&
		c.preferredBreakByte != len(c.code) {
		lo := rune(0xED&0x0F)<<12 | rune(str[1]&0x3F)<<6 | rune(str[2]&0x3F)
		cp := utf16.DecodeRune(rune(c.lastChar), lo)
		c.code = utf8.AppendRune(c.code[:len(c.code)-3], cp)
		c.code = append(c.code, str[3:]...)
	} else {
		c.code = append(c.code, str...)
	}
	units := 0
	var last uint16
	for i := 0; i < len(str); {
		b0 := str[i]
		if b0 < 0x80 {
			last = uint16(b0)
			units++
			i++
			continue
		}
		if b0 == 0xED && i+3 <= len(str) && str[i+1] >= 0xA0 && str[i+1] <= 0xBF && str[i+2]&0xC0 == 0x80 {
			last = uint16(b0&0x0F)<<12 | uint16(str[i+1]&0x3F)<<6 | uint16(str[i+2]&0x3F)
			units++
			i += 3
			continue
		}
		r, sz := utf8.DecodeRuneInString(str[i:])
		i += sz
		if r >= 0x10000 {
			last = uint16(0xDC00 + ((r - 0x10000) & 0x3FF))
			units += 2
		} else {
			last = uint16(r)
			units++
		}
	}
	c.lastChar = last
	return units
}

func (c *CodeConsumer) codeAppendASCII(str string) {
	if len(str) == 0 {
		return
	}
	c.lastChar = uint16(str[len(str)-1])
	c.code = append(c.code, str...)
}

func (c *CodeConsumer) codeAppendByte(b byte) {
	c.lastChar = uint16(b)
	c.code = append(c.code, b)
}

func (c *CodeConsumer) codeInsertByte(offset int, b byte) {
	n := len(c.code)
	if offset < 0 || offset > n {
		panic(jutil.NewStringIndexOutOfBoundsException(offset, n))
	}
	if offset == n {
		c.lastChar = uint16(b)
	}
	c.code = append(c.code, 0)
	copy(c.code[offset+1:], c.code[offset:n])
	c.code[offset] = b
}

func (c *CodeConsumer) recordCost(units int, lastChar uint16) {
	c.lastChar = lastChar
	c.cost += units
	if c.maxCost <= c.cost {
		c.continueProcessing = false
	}
}

func (c *CodeConsumer) flushGzScratch() {
	if _, err := c.gzStream.Write(c.gzScratch); err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
	}
}

func (c *CodeConsumer) recordSizeBytes(units int, lastChar uint16, b []byte) {
	c.lastChar = lastChar
	c.cost += units
	if c.trackGzSize {
		c.gzScratch = append(c.gzScratch[:0], b...)
		c.flushGzScratch()
	}
}

func (c *CodeConsumer) appendRawASCIIBytes(b []byte) {
	last := uint16(b[len(b)-1])
	switch c.mode {
	case sinkCompact, sinkPretty:
		c.lastChar = last
		c.lineLength += len(b)
		c.code = append(c.code, b...)
	case sinkCost:
		c.recordCost(len(b), last)
	case sinkSize:
		c.recordSizeBytes(len(b), last, b)
	}
}

// emitTokenSeparator inserts any mandatory token-separation or optional pretty
// space needed before emitting a token starting with ch.
func (c *CodeConsumer) emitTokenSeparator(ch uint16) {
	if (CodeConsumer_IsWordChar(ch) || ch == '\\') && CodeConsumer_IsWordChar(c.lastChar) {
		// need space to separate. This is not pretty printing.
		// For example: "return foo;"
		c.Append(" ")
	} else if ch == '/' && (c.lastChar == '/' || c.lastChar == '<') {
		// Do not allow a forward slash to appear after a DIV or LT.
		// For example,
		// REGEXP DIV REGEXP
		// is valid and should print like
		// / // / /
		c.Append(" ")
	} else if (ch == '"' || ch == '\'') && CodeConsumer_IsWordChar(c.lastChar) {
		c.MaybeInsertSpace()
	}
}

// Append appends a string to the consumer's output.
//
// Java: CodeConsumer.append(String)
func (c *CodeConsumer) Append(str string) {
	switch c.mode {
	case sinkPretty:
		c.maybeIndent()
		c.lineLength += c.codeAppend(str)
	case sinkCompact:
		c.lineLength += c.codeAppend(str)
	case sinkCost:
		if n := len(str); n > 0 && jutil.IsASCII(str) {
			c.recordCost(n, uint16(str[n-1]))
			return
		}
		length := jutil.StringLength(str)
		c.recordCost(length, jutil.StringCharAt(str, length-1))
	case sinkSize:
		if n := len(str); n > 0 {
			if jutil.IsASCII(str) {
				c.cost += n
				c.lastChar = uint16(str[n-1])
				if c.trackGzSize {
					c.gzScratch = append(c.gzScratch[:0], str...)
					c.flushGzScratch()
				}
				return
			}
			len_ := jutil.StringLength(str)
			c.cost += len_
			c.lastChar = jutil.StringCharAt(str, len_-1)
			if c.trackGzSize {
				c.gzScratch = jutil.AppendStringBytesUTF8(c.gzScratch[:0], str)
				c.flushGzScratch()
			}
		}
	}
}

// appendASCIIBytes appends a non-empty pure-ASCII byte slice directly without
// allocating a temporary string.
func (c *CodeConsumer) appendASCIIBytes(b []byte) {
	if c.mode == sinkPretty {
		c.maybeIndent()
	}
	c.appendRawASCIIBytes(b)
}

// beginToken prepares statement-end, separator-space, and pretty indentation
// before streaming a single-line token whose first ASCII byte is firstByte.
func (c *CodeConsumer) beginToken(firstByte byte) {
	c.MaybeEndStatement()
	c.emitTokenSeparator(uint16(firstByte))
	if c.mode == sinkPretty {
		c.maybeIndent()
	}
}

// appendContinuationChunk appends a non-empty chunk of valid BMP UTF-8 (all
// surrogates already escaped to \uXXXX) after beginToken has run.
func (c *CodeConsumer) appendContinuationChunk(b []byte, isASCII bool) {
	if isASCII {
		c.appendRawASCIIBytes(b)
		return
	}

	// Non-ASCII BMP UTF-8 chunk: every rune is 1 UTF-16 code unit.
	unitsAdded := utf8.RuneCount(b)
	r, _ := utf8.DecodeLastRune(b)
	last := uint16(r)
	switch c.mode {
	case sinkCompact, sinkPretty:
		c.hasNonASCII = true
		c.code = append(c.code, b...)
		c.lastChar = last
		c.lineLength += unitsAdded
	case sinkCost:
		c.recordCost(unitsAdded, last)
	case sinkSize:
		c.recordSizeBytes(unitsAdded, last, b)
	}
}

// addPlainQuotedASCII writes quote + s + quote directly when s is known to be
// plain single-line ASCII.
func (c *CodeConsumer) addPlainQuotedASCII(quote byte, s string) {
	c.beginToken(quote)
	switch c.mode {
	case sinkCompact, sinkPretty:
		c.lastChar = uint16(quote)
		c.lineLength += len(s) + 2
		c.code = append(c.code, quote)
		c.code = append(c.code, s...)
		c.code = append(c.code, quote)
	case sinkCost:
		c.recordCost(len(s)+2, uint16(quote))
	case sinkSize:
		c.lastChar = uint16(quote)
		c.cost += len(s) + 2
		if c.trackGzSize {
			c.gzScratch = append(c.gzScratch[:0], quote)
			c.gzScratch = append(c.gzScratch, s...)
			c.gzScratch = append(c.gzScratch, quote)
			c.flushGzScratch()
		}
	}
}

// addSingleLineASCIIBytes sanitizes and appends a non-empty single-line ASCII slice.
func (c *CodeConsumer) addSingleLineASCIIBytes(b []byte) {
	c.MaybeEndStatement()
	if len(b) == 0 {
		return
	}
	c.emitTokenSeparator(uint16(b[0]))
	c.appendASCIIBytes(b)
}

// Add appends a (possibly multiline) string to the code, keeping track of
// the current line length and count.
//
// Java: CodeConsumer.add(String)
func (c *CodeConsumer) Add(newcode string) {
	c.MaybeEndStatement()

	if newcode == "" {
		return
	}

	c.emitTokenSeparator(codeConsumer_firstChar(newcode))

	endOfLine := strings.IndexByte(newcode, '\n')
	if endOfLine < 0 {
		c.Append(newcode)
		return
	}

	startOfLine := 0
	for endOfLine >= 0 {
		if endOfLine > startOfLine {
			c.Append(newcode[startOfLine:endOfLine])
		}
		// Breaking is non-optional. Newlines added this way must be
		// preserved (e.g. newlines in template literals).
		c.StartNewLine()

		startOfLine = endOfLine + 1
		endOfLine = strings.IndexByte(newcode[startOfLine:], '\n')
		if endOfLine >= 0 {
			endOfLine += startOfLine
		}
	}
	if len(newcode) > startOfLine {
		c.Append(newcode[startOfLine:])
	}
}

// Java: CodeConsumer.addIdentifier(String)
func (c *CodeConsumer) AddIdentifier(identifier string) {
	if c.mode == sinkCost {
		c.Add(InlineCostEstimator_ESTIMATED_IDENTIFIER)
		return
	}
	c.Add(identifier)
}

// Java: CodeConsumer.appendBlockStart()
func (c *CodeConsumer) AppendBlockStart() {
	if c.mode == sinkPretty {
		c.MaybeInsertSpace()
		c.Add("{")
		c.indent++
		return
	}
	c.Append("{")
}

// Java: CodeConsumer.appendBlockEnd()
func (c *CodeConsumer) AppendBlockEnd() {
	if c.mode == sinkPretty {
		c.MaybeEndStatement()
		c.EndLine()
		c.indent--
		c.Add("}")
		return
	}
	c.Append("}")
}

// Java: CodeConsumer.startNewLine()
func (c *CodeConsumer) StartNewLine() {
	if c.mode > sinkPretty {
		return
	}
	if c.lineLength <= 0 && !c.IsInTemplateLiteral() {
		return
	}
	c.codeAppendByte('\n')
	c.lineLength = 0
	c.lineIndex++
	c.lineStartByte = len(c.code)
}

// Java: CodeConsumer.maybeLineBreak()
func (c *CodeConsumer) MaybeLineBreak() {
	switch c.mode {
	case sinkCompact:
		if c.lineBreak && c.sawFunction {
			c.StartNewLine()
			c.sawFunction = false
		}
		len_ := len(c.code)
		if c.preferredBreakByte >= c.lineStartByte && c.preferredBreakByte == len_-1 && c.code[len_-1] == ';' {
			c.preferredBreakByte = len_
			c.preferredBreakCol++
		}
		if c.lineLength > c.lineLengthThreshold {
			c.MaybeCutLine()
		}
	case sinkPretty:
		c.MaybeCutLine()
	}
}

// Java: CodeConsumer.maybeCutLine()
func (c *CodeConsumer) MaybeCutLine() {
	switch c.mode {
	case sinkCompact:
		if c.lineLength <= c.lineLengthThreshold {
			return
		}
		if c.preferredBreakByte > c.lineStartByte &&
			c.preferredBreakCol < c.lineLength &&
			c.preferredBreakByte < len(c.code) {
			c.codeInsertByte(c.preferredBreakByte, '\n')
			c.ReportLineCut(c.lineIndex, c.preferredBreakCol)
			c.lineIndex++
			c.lineLength -= c.preferredBreakCol
			c.lineStartByte = c.preferredBreakByte + 1
		} else {
			c.StartNewLine()
		}
	case sinkPretty:
		if c.lineLength > c.lineLengthThreshold {
			c.StartNewLine()
		}
	}
}

// Java: CodeConsumer.endLine()
func (c *CodeConsumer) EndLine() {
	if c.mode == sinkPretty {
		c.StartNewLine()
	}
}

// Java: CodeConsumer.notePreferredLineBreak()
func (c *CodeConsumer) NotePreferredLineBreak() {
	if c.mode == sinkCompact {
		c.preferredBreakByte = len(c.code)
		c.preferredBreakCol = c.lineLength
	}
}

// Java: CodeConsumer.beginBlock()
func (c *CodeConsumer) BeginBlock() {
	if c.statementNeedsEnded {
		c.Append(";")
		c.MaybeLineBreak()
	}
	c.AppendBlockStart()
	c.EndLine()
	c.statementNeedsEnded = false
}

// Java: CodeConsumer.endBlock()
func (c *CodeConsumer) EndBlock() {
	c.EndBlockWithShouldEndLine(false)
}

// Java: CodeConsumer.endBlock(boolean)
func (c *CodeConsumer) EndBlockWithShouldEndLine(shouldEndLine bool) {
	c.AppendBlockEnd()
	if shouldEndLine {
		c.EndLine()
	}
	c.statementNeedsEnded = false
}

// Java: CodeConsumer.listSeparator()
func (c *CodeConsumer) ListSeparator() {
	if c.mode == sinkPretty {
		c.Add(", ")
	} else {
		c.Add(",")
	}
	c.MaybeLineBreak()
}

// Java: CodeConsumer.optionalListSeparator()
func (c *CodeConsumer) OptionalListSeparator() {
	if c.mode == sinkPretty {
		c.Add(",")
		c.MaybeLineBreak()
	}
}

// Java: CodeConsumer.endStatement(boolean)
func (c *CodeConsumer) EndStatement(hasTrailingCommentOnSameLine bool) {
	c.EndStatementWithSemiColon(false, hasTrailingCommentOnSameLine)
}

// Java: CodeConsumer.endStatement(boolean, boolean)
func (c *CodeConsumer) EndStatementWithSemiColon(needSemiColon bool, hasTrailingCommentOnSameLine bool) {
	if c.mode == sinkPretty {
		c.Add(";")
		if !hasTrailingCommentOnSameLine {
			c.EndLine()
		}
		c.statementNeedsEnded = false
		return
	}
	if needSemiColon {
		c.Append(";")
		if !hasTrailingCommentOnSameLine {
			c.MaybeLineBreak()
		}
		c.statementNeedsEnded = false
	} else if c.statementStarted {
		c.statementNeedsEnded = true
	}
}

// Java: CodeConsumer.maybeEndStatement()
func (c *CodeConsumer) MaybeEndStatement() {
	if c.statementNeedsEnded {
		c.Append(";")
		c.MaybeLineBreak()
		c.EndLine()
		c.statementNeedsEnded = false
	}
	c.statementStarted = true
}

// Java: CodeConsumer.endFunction(boolean)
func (c *CodeConsumer) EndFunction(statementContext bool) {
	c.sawFunction = true
	if statementContext {
		c.EndLine()
		if c.mode == sinkPretty {
			c.StartNewLine()
		}
	}
}

// Java: CodeConsumer.endClass(boolean)
func (c *CodeConsumer) EndClass(statementContext bool) {
	if statementContext {
		c.EndLine()
	}
}

// Java: CodeConsumer.beginCaseBody()
func (c *CodeConsumer) BeginCaseBody() {
	c.Append(":")
	if c.mode == sinkPretty {
		c.indent++
		c.EndLine()
	}
}

// Java: CodeConsumer.endCaseBody()
func (c *CodeConsumer) EndCaseBody() {
	if c.mode == sinkPretty {
		c.indent--
	}
}

// Java: CodeConsumer.beginTemplateLit()
func (c *CodeConsumer) BeginTemplateLit() {
	jutil.CheckState(c.templateLitDepth == c.templateLitSubDepth)

	c.MaybeEndStatement()
	c.Append("`")
	c.templateLitDepth++
}

// Java: CodeConsumer.beginTemplateLitSub()
func (c *CodeConsumer) BeginTemplateLitSub() {
	jutil.CheckState(c.IsInTemplateLiteral())

	c.Append("${")
	c.templateLitSubDepth++
}

// Java: CodeConsumer.endTemplateLitSub()
func (c *CodeConsumer) EndTemplateLitSub() {
	jutil.CheckState(c.templateLitSubDepth > 0)
	jutil.CheckState(c.templateLitDepth == c.templateLitSubDepth)

	c.Append("}")
	c.templateLitSubDepth--
}

// Java: CodeConsumer.endTemplateLit()
func (c *CodeConsumer) EndTemplateLit() {
	jutil.CheckState(c.templateLitDepth > 0)
	jutil.CheckState(c.IsInTemplateLiteral())

	c.Append("`")
	c.templateLitDepth--
}

// Java: CodeConsumer.isInTemplateLiteral()
func (c *CodeConsumer) IsInTemplateLiteral() bool {
	return c.templateLitDepth == c.templateLitSubDepth+1
}

// Java: CodeConsumer.appendOp(String, boolean)
func (c *CodeConsumer) AppendOp(op string, binOp bool) {
	if c.mode == sinkPretty {
		if c.lastChar != ' ' && binOp && codeConsumer_firstChar(op) != ',' {
			c.Add(" ")
		}
		c.Add(op)
		if binOp {
			c.Add(" ")
		}
		return
	}
	c.Append(op)
}

// Java: CodeConsumer.addOp(String, boolean)
func (c *CodeConsumer) AddOp(op string, binOp bool) {
	c.MaybeEndStatement()

	first := codeConsumer_firstChar(op)
	prev := c.lastChar

	if (first == '+' || first == '-') && prev == first {
		// This is not pretty printing. This is to prevent misparsing of
		// things like "x + ++y" or "x++ + ++y"
		c.Append(" ")
	} else if codeConsumer_isLetterChar(first) && CodeConsumer_IsWordChar(prev) {
		// Make sure there is a space after e.g. instanceof , typeof
		c.Append(" ")
	} else if (prev == '-' && first == '>') || (prev == '<' && first == '!') {
		// Make sure that we don't emit "<!--" or "-->"
		c.Append(" ")
	}

	c.AppendOp(op, binOp)

	if binOp {
		c.MaybeCutLine()
	}
}

// Java: CodeConsumer.addNumber(double, Node)
func (c *CodeConsumer) AddNumber(x float64, n rhino.Node) {
	if !base.JSCompDoubles_IsPositive(x) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.DoubleToString(x)))
	}
	if c.mode == sinkCost {
		c.Add("0")
		return
	}
	if c.mode == sinkPretty {
		if numberFromSource, ok := c.GetNumberFromSource(n); ok {
			if d, ok := codePrinter_parseDouble(numberFromSource); ok && x == d {
				c.AddConstant(numberFromSource)
				return
			}
		}
	}
	var buf [32]byte
	if b, ok := codeConsumer_formatIntNumber(buf[:0], x); ok {
		c.addSingleLineASCIIBytes(b)
		return
	}
	c.AddConstant(CodeConsumer_FormatNumber(x))
}

func codeConsumer_formatIntNumber(dst []byte, x float64) ([]byte, bool) {
	value := jutil.D2L(x)
	if float64(value) != x {
		return nil, false
	}

	mantissa := value
	exp := 0
	if x >= 100 {
		for mantissa%10 == 0 {
			mantissa /= 10
			exp++
		}
	}

	if exp > 2 {
		dst = strconv.AppendInt(dst, mantissa, 10)
		dst = append(dst, 'E')
		dst = strconv.AppendInt(dst, int64(exp), 10)
		return dst, true
	}

	if value <= 1000000000000 {
		// Values <1E12 are shorter in decimal
		return strconv.AppendInt(dst, value, 10), true
	}

	var decBuf, hexBuf [24]byte
	dec := strconv.AppendInt(decBuf[:0], value, 10)
	hex := strconv.AppendUint(hexBuf[:0], uint64(value), 16)
	if len(hex)+2 < len(dec) {
		dst = append(dst, '0', 'x')
		dst = append(dst, hex...)
		return dst, true
	}
	return append(dst, dec...), true
}

// Java: CodeConsumer.formatNumber(double)
func CodeConsumer_FormatNumber(x float64) string {
	var buf [32]byte
	if b, ok := codeConsumer_formatIntNumber(buf[:0], x); ok {
		return string(b)
	}
	s := strings.ReplaceAll(jutil.DoubleToString(x), ".0E", "E")
	if strings.HasPrefix(s, "0.") {
		s = "." + s[2:]
	}
	return s
}

// Java: CodeConsumer.addBigInt(BigInteger)
func (c *CodeConsumer) AddBigInt(bi *big.Int) {
	if c.mode == sinkCost {
		c.Add("0")
		return
	}
	var hexBuf, decBuf [64]byte
	hex := append(hexBuf[:0], '0', 'x')
	hex = bi.Append(hex, 16)
	hex = append(hex, 'n')

	dec := bi.Append(decBuf[:0], 10)
	dec = append(dec, 'n')

	if len(hex) < len(dec) {
		c.addSingleLineASCIIBytes(hex)
	} else {
		c.addSingleLineASCIIBytes(dec)
	}
}

// Java: CodeConsumer.addConstant(String)
func (c *CodeConsumer) AddConstant(newcode string) {
	if c.mode == sinkCost {
		c.Add("0")
		return
	}
	c.Add(newcode)
}

var codeConsumer_asciiIsLetter = [128]bool{
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true,
	'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true,
	'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true,
	'Y': true, 'Z': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true,
	'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true,
	'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true,
	'y': true, 'z': true,
}

var codeConsumer_asciiIsWord = [128]bool{
	'$': true,
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true,
	'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true,
	'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true,
	'Y': true, 'Z': true,
	'_': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true,
	'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true,
	'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true,
	'y': true, 'z': true,
}

func codeConsumer_isLetterChar(ch uint16) bool {
	if ch < 0x80 {
		return codeConsumer_asciiIsLetter[ch]
	}
	return jutil.CharacterIsLetter(ch)
}

//go:noinline
func codeConsumer_isWordCharNonASCII(ch uint16) bool {
	return jutil.CharacterIsLetter(ch) || jutil.CharacterIsDigit(ch)
}

// Java: CodeConsumer.isWordChar(char)
func CodeConsumer_IsWordChar(ch uint16) bool {
	if ch < 0x80 {
		return codeConsumer_asciiIsWord[ch]
	}
	return codeConsumer_isWordCharNonASCII(ch)
}

// ShouldPreserveExtras returns whether single-statement bodies of loops/ifs
// should be wrapped in a block.
//
// Java: CodeConsumer.shouldPreserveExtras(Node)
func (c *CodeConsumer) ShouldPreserveExtras(n rhino.Node) bool {
	if c.mode != sinkPretty {
		return false
	}
	if !n.IsBlock() || !n.IsAddedBlock() || !n.HasParent() {
		return true
	}
	parent := n.GetParent()
	isElse := parent.IsIf() && parent.HasXChildren(3) && n == parent.GetLastChild()
	onlyChildIsIf := n.HasOneChild() && n.GetFirstChild().IsIf()
	return !(isElse && onlyChildIsIf)
}

// MaybeInsertSpace inserts optional pretty-printing spaces.
//
// Java: CodeConsumer.maybeInsertSpace()
func (c *CodeConsumer) MaybeInsertSpace() {
	if c.mode == sinkPretty && c.lastChar != ' ' && c.lastChar != '\n' {
		c.Add(" ")
	}
}

// BreakAfterBlockFor returns whether a line break can be added after the
// specified BLOCK.
//
// Java: CodeConsumer.breakAfterBlockFor(Node, boolean)
func (c *CodeConsumer) BreakAfterBlockFor(n rhino.Node, statementContext bool) bool {
	if c.mode != sinkPretty {
		return statementContext
	}
	jutil.CheckStateMsg1(n.IsBlock(), "%s", n)
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_DO:
		return false // Don't break before 'while' in DO-WHILE statements.
	case rhino.Token_FUNCTION:
		return false // FUNCTIONs are handled separately, don't break here.
	case rhino.Token_TRY:
		return n != parent.GetFirstChild() // Don't break before catch
	case rhino.Token_CATCH:
		return !NodeUtil_HasFinally(CodePrinter_PrettyCodePrinter_GetTryForCatch(parent)) // Don't break before finally
	case rhino.Token_IF:
		return n == parent.GetLastChild() // Don't break before else
	default:
		return true
	}
}

// EndFile is called when we're at the end of a file.
//
// Java: CodeConsumer.endFile()
func (c *CodeConsumer) EndFile() {
	if c.mode == sinkPretty {
		c.MaybeEndStatement()
	}
}
