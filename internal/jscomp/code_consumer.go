// Port of com/google/javascript/jscomp/CodeConsumer.java @ closure-compiler v20260415

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
	"math/big"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodeConsumer is the abstracted consumer of the CodeGenerator output.
//
// See CodeGenerator, CodePrinter, InlineCostEstimator.
//
// Go-only: language-adapter — an abstract class with subclasses: this
// interface declares every method (the final ones too, so that callers
// holding a CodeConsumer reach them), and CodeConsumerBase holds the fields
// and the base implementations. Subclasses embed CodeConsumerBase and call
// InitCodeConsumerBase with themselves.
type CodeConsumer interface {
	StartSourceMapping(node *rhino.Node)
	EndSourceMapping(node *rhino.Node)
	TrackLicenses(node *rhino.Node)
	ContinueProcessing() bool
	GetLastChar() uint16
	Add(newcode string)
	Append(str string)
	AddIdentifier(identifier string)
	AppendBlockStart()
	AppendBlockEnd()
	StartNewLine()
	MaybeLineBreak()
	MaybeCutLine()
	EndLine()
	NotePreferredLineBreak()
	BeginBlock()
	EndBlock()
	EndBlockWithShouldEndLine(shouldEndLine bool)
	ListSeparator()
	OptionalListSeparator()
	EndStatement(hasTrailingCommentOnSameLine bool)
	EndStatementWithSemiColon(needSemiColon bool, hasTrailingCommentOnSameLine bool)
	MaybeEndStatement()
	EndFunction(statementContext bool)
	EndClass(statementContext bool)
	BeginCaseBody()
	EndCaseBody()
	BeginTemplateLit()
	BeginTemplateLitSub()
	EndTemplateLitSub()
	EndTemplateLit()
	IsInTemplateLiteral() bool
	AppendOp(op string, binOp bool)
	AddOp(op string, binOp bool)
	AddNumber(x float64, n *rhino.Node)
	AddBigInt(bi *big.Int)
	AddConstant(newcode string)
	ShouldPreserveExtras(n *rhino.Node) bool
	MaybeInsertSpace()
	BreakAfterBlockFor(n *rhino.Node, statementContext bool) bool
	EndFile()

	GetCodeConsumerBase() *CodeConsumerBase
}

// CodeConsumerBase holds CodeConsumer's fields and base implementations.
type CodeConsumerBase struct {
	self CodeConsumer

	statementNeedsEnded bool
	statementStarted    bool
	sawFunction         bool

	// State tracking for template literals. Remember that template literal
	// substitutions can contain additional template literals.
	templateLitDepth    int
	templateLitSubDepth int
}

// InitCodeConsumerBase is the CodeConsumer() constructor; self is the
// object being constructed.
func (c *CodeConsumerBase) InitCodeConsumerBase(self CodeConsumer) {
	c.self = self
	c.statementNeedsEnded = false
	c.statementStarted = false
	c.sawFunction = false
	c.templateLitDepth = 0
	c.templateLitSubDepth = 0
}

// GetCodeConsumerBase gives access to the base fields through the
// interface.
func (c *CodeConsumerBase) GetCodeConsumerBase() *CodeConsumerBase { return c }

// StartSourceMapping starts the source mapping for the given node at the
// current position.
//
// Java: CodeConsumer.startSourceMapping(Node)
func (c *CodeConsumerBase) StartSourceMapping(node *rhino.Node) {
}

// EndSourceMapping finishes the source mapping for the given node at the
// current position.
//
// Java: CodeConsumer.endSourceMapping(Node)
func (c *CodeConsumerBase) EndSourceMapping(node *rhino.Node) {
}

// TrackLicenses indicates to the CodeConsumer that this Node might carry
// licensing information, and allows the code consumer to manage licenses
// as it sees fit.
//
// Java: CodeConsumer.trackLicenses(Node)
func (c *CodeConsumerBase) TrackLicenses(node *rhino.Node) {}

// ContinueProcessing provides a means of interrupting the CodeGenerator.
// Derived classes should return false to stop further processing.
//
// Java: CodeConsumer.continueProcessing()
func (c *CodeConsumerBase) ContinueProcessing() bool {
	return true
}

// codeConsumer_firstChar is s.charAt(0).
func codeConsumer_firstChar(s string) uint16 {
	if len(s) > 0 && s[0] < 0x80 {
		return uint16(s[0])
	}
	return jutil.StringCharAt(s, 0)
}

// Add appends a (possibly multiline) string to the code, keeping track of
// the current line length and count.
//
// Clients should use this method rather than Append. It sanitizes the input
// to Append.
//
// Final in Java.
//
// Java: CodeConsumer.add(String)
func (c *CodeConsumerBase) Add(newcode string) {
	c.self.MaybeEndStatement()

	if newcode == "" {
		return
	}

	ch := codeConsumer_firstChar(newcode)
	if (CodeConsumer_IsWordChar(ch) || ch == '\\') && CodeConsumer_IsWordChar(c.self.GetLastChar()) {
		// need space to separate. This is not pretty printing.
		// For example: "return foo;"
		c.self.Append(" ")
	} else if ch == '/' && c.self.GetLastChar() == '/' {
		// Do not allow a forward slash to appear after a DIV.
		// For example,
		// REGEXP DIV REGEXP
		// is valid and should print like
		// / // / /
		c.self.Append(" ")
	} else if (ch == '"' || ch == '\'') && CodeConsumer_IsWordChar(c.self.GetLastChar()) {
		c.self.MaybeInsertSpace()
	}

	// Iterate through the new code and add each contained line, followed by
	// a break. Remember that the string may start and end in the middle of
	// a line. We do this rather primitively because this method is called
	// frequently and most invocations have only one line.
	// TODO(nickreid): There are other possible newline characters
	// recognized by the JS spec. We should also be considering them.
	//
	// Go-only: language-adapter — byte offsets: every split is at an ASCII
	// '\n', so the lines are the same as with Java's UTF-16 offsets.
	startOfLine := 0
	endOfLine := strings.IndexByte(newcode, '\n')
	for endOfLine >= 0 {
		if endOfLine > startOfLine {
			// Append line only if it is non-empty.
			c.self.Append(newcode[startOfLine:endOfLine])
		}
		// Breaking is non-optional. Newlines added this way must be
		// preserved (e.g. newlines in template literals).
		c.self.StartNewLine()

		startOfLine = endOfLine + 1 // Jump over the newline char.
		endOfLine = strings.IndexByte(newcode[startOfLine:], '\n')
		if endOfLine >= 0 {
			endOfLine += startOfLine
		}
	}
	if len(newcode) > startOfLine {
		// Append line only if it is non-empty.
		c.self.Append(newcode[startOfLine:]) // Append the last or only line without breaking.
	}
}

// Java: CodeConsumer.addIdentifier(String)
func (c *CodeConsumerBase) AddIdentifier(identifier string) {
	c.self.Add(identifier)
}

// Java: CodeConsumer.appendBlockStart()
func (c *CodeConsumerBase) AppendBlockStart() {
	c.self.Append("{")
}

// Java: CodeConsumer.appendBlockEnd()
func (c *CodeConsumerBase) AppendBlockEnd() {
	c.self.Append("}")
}

// Java: CodeConsumer.startNewLine()
func (c *CodeConsumerBase) StartNewLine() {
}

// Java: CodeConsumer.maybeLineBreak()
func (c *CodeConsumerBase) MaybeLineBreak() {
	c.self.MaybeCutLine()
}

// Java: CodeConsumer.maybeCutLine()
func (c *CodeConsumerBase) MaybeCutLine() {
}

// Java: CodeConsumer.endLine()
func (c *CodeConsumerBase) EndLine() {
}

// Java: CodeConsumer.notePreferredLineBreak()
func (c *CodeConsumerBase) NotePreferredLineBreak() {
}

// Java: CodeConsumer.beginBlock()
func (c *CodeConsumerBase) BeginBlock() {
	if c.statementNeedsEnded {
		c.self.Append(";")
		c.self.MaybeLineBreak()
	}
	c.self.AppendBlockStart()

	c.self.EndLine()
	c.statementNeedsEnded = false
}

// Java: CodeConsumer.endBlock()
func (c *CodeConsumerBase) EndBlock() {
	c.self.EndBlockWithShouldEndLine(false)
}

// Java: CodeConsumer.endBlock(boolean)
func (c *CodeConsumerBase) EndBlockWithShouldEndLine(shouldEndLine bool) {
	c.self.AppendBlockEnd()
	if shouldEndLine {
		c.self.EndLine()
	}
	c.statementNeedsEnded = false
}

// Java: CodeConsumer.listSeparator()
func (c *CodeConsumerBase) ListSeparator() {
	c.self.Add(",")
	c.self.MaybeLineBreak()
}

// Java: CodeConsumer.optionalListSeparator()
func (c *CodeConsumerBase) OptionalListSeparator() {}

// EndStatement indicates the end of a statement and a ';' may need to be
// added. But we don't add it now, in case we're at the end of a block (in
// which case we don't have to add the ';'). See maybeEndStatement()
//
// Java: CodeConsumer.endStatement(boolean)
func (c *CodeConsumerBase) EndStatement(hasTrailingCommentOnSameLine bool) {
	c.self.EndStatementWithSemiColon(false, hasTrailingCommentOnSameLine)
}

// Java: CodeConsumer.endStatement(boolean, boolean)
func (c *CodeConsumerBase) EndStatementWithSemiColon(needSemiColon bool, hasTrailingCommentOnSameLine bool) {
	if needSemiColon {
		c.self.Append(";")
		if !hasTrailingCommentOnSameLine {
			c.self.MaybeLineBreak()
		}
		c.statementNeedsEnded = false
	} else if c.statementStarted {
		c.statementNeedsEnded = true
	}
}

// MaybeEndStatement is to be called when we're in a statement. If the prev
// statement needs to be ended, add a ';'.
//
// Java: CodeConsumer.maybeEndStatement()
func (c *CodeConsumerBase) MaybeEndStatement() {
	// Add a ';' if we need to.
	if c.statementNeedsEnded {
		c.self.Append(";")
		c.self.MaybeLineBreak()
		c.self.EndLine()
		c.statementNeedsEnded = false
	}
	c.statementStarted = true
}

// Java: CodeConsumer.endFunction(boolean)
func (c *CodeConsumerBase) EndFunction(statementContext bool) {
	c.sawFunction = true
	if statementContext {
		c.self.EndLine()
	}
}

// Java: CodeConsumer.endClass(boolean)
func (c *CodeConsumerBase) EndClass(statementContext bool) {
	if statementContext {
		c.self.EndLine()
	}
}

// Java: CodeConsumer.beginCaseBody()
func (c *CodeConsumerBase) BeginCaseBody() {
	c.self.Append(":")
}

// Java: CodeConsumer.endCaseBody()
func (c *CodeConsumerBase) EndCaseBody() {
}

// BeginTemplateLit is final in Java.
//
// Java: CodeConsumer.beginTemplateLit()
func (c *CodeConsumerBase) BeginTemplateLit() {
	jutil.CheckState(c.templateLitDepth == c.templateLitSubDepth)

	c.self.MaybeEndStatement()
	c.self.Append("`")
	c.templateLitDepth++
}

// BeginTemplateLitSub is final in Java.
//
// Java: CodeConsumer.beginTemplateLitSub()
func (c *CodeConsumerBase) BeginTemplateLitSub() {
	// This method pair exists because '$' behaves differently inside
	// template literals. We want to append it without sanitizing in the
	// generic way add() would.
	jutil.CheckState(c.self.IsInTemplateLiteral())

	c.self.Append("${")
	c.templateLitSubDepth++
}

// EndTemplateLitSub is final in Java.
//
// Java: CodeConsumer.endTemplateLitSub()
func (c *CodeConsumerBase) EndTemplateLitSub() {
	// This method pair exists because '$' behaves differently inside
	// template literals. We want to append it without sanitizing in the
	// generic way add() would.
	jutil.CheckState(c.templateLitSubDepth > 0)
	jutil.CheckState(c.templateLitDepth == c.templateLitSubDepth)

	c.self.Append("}")
	c.templateLitSubDepth--
}

// EndTemplateLit is final in Java.
//
// Java: CodeConsumer.endTemplateLit()
func (c *CodeConsumerBase) EndTemplateLit() {
	jutil.CheckState(c.templateLitDepth > 0)
	jutil.CheckState(c.self.IsInTemplateLiteral())

	c.self.Append("`")
	c.templateLitDepth--
}

// IsInTemplateLiteral is final in Java.
//
// Java: CodeConsumer.isInTemplateLiteral()
func (c *CodeConsumerBase) IsInTemplateLiteral() bool {
	// We're inside a template literal but not a substitution within that
	// literal.
	return c.templateLitDepth == c.templateLitSubDepth+1
}

// Java: CodeConsumer.appendOp(String, boolean)
func (c *CodeConsumerBase) AppendOp(op string, binOp bool) {
	c.self.Append(op)
}

// Java: CodeConsumer.addOp(String, boolean)
func (c *CodeConsumerBase) AddOp(op string, binOp bool) {
	c.self.MaybeEndStatement()

	first := codeConsumer_firstChar(op)
	prev := c.self.GetLastChar()

	if (first == '+' || first == '-') && prev == first {
		// This is not pretty printing. This is to prevent misparsing of
		// things like "x + ++y" or "x++ + ++y"
		c.self.Append(" ")
	} else if jutil.CharacterIsLetter(first) && CodeConsumer_IsWordChar(prev) {
		// Make sure there is a space after e.g. instanceof , typeof
		c.self.Append(" ")
	} else if (prev == '-' && first == '>') || (prev == '<' && first == '!') {
		// Make sure that we don't emit "<!--" or "-->"
		c.self.Append(" ")
	}

	// Allow formatting around the operator.
	c.self.AppendOp(op, binOp)

	// Line breaking after an operator is always safe. Line breaking before
	// an operator on the other hand is not. We only line break after a bin
	// op because it looks strange.
	if binOp {
		c.self.MaybeCutLine()
	}
}

// Java: CodeConsumer.addNumber(double, Node)
func (c *CodeConsumerBase) AddNumber(x float64, n *rhino.Node) {
	if !base.JSCompDoubles_IsPositive(x) {
		// checkState(isPositive(x), x)
		panic(jutil.NewIllegalStateExceptionMsg(jutil.DoubleToString(x)))
	}

	if float64(jutil.D2L(x)) != x {
		s := strings.ReplaceAll(jutil.DoubleToString(x), ".0E", "E")
		if strings.HasPrefix(s, "0.") {
			s = "." + s[2:]
		}
		c.self.AddConstant(s)
		return
	}

	value := jutil.D2L(x)
	mantissa := value
	exp := 0
	if x >= 100 {
		for mantissa%10 == 0 {
			mantissa /= 10
			exp++
		}
	}

	if exp > 2 {
		c.self.AddConstant(strconv.FormatInt(mantissa, 10) + "E" + strconv.Itoa(exp))
		return
	}

	decValueString := strconv.FormatInt(value, 10)
	if value <= 1000000000000 {
		// Values <1E12 are shorter in decimal
		c.self.AddConstant(decValueString)
		return
	}

	hexValueString := jutil.LongToHexString(value)
	if len(hexValueString)+2 < len(decValueString) {
		c.self.AddConstant("0x" + hexValueString)
	} else {
		c.self.AddConstant(decValueString)
	}
}

// Java: CodeConsumer.addBigInt(BigInteger)
func (c *CodeConsumerBase) AddBigInt(bi *big.Int) {
	hexEncoded := "0x" + bi.Text(16) + "n"
	decimalEncoded := bi.String() + "n"
	if len(hexEncoded) < len(decimalEncoded) {
		c.self.AddConstant(hexEncoded)
	} else {
		c.self.AddConstant(decimalEncoded)
	}
}

// Java: CodeConsumer.addConstant(String)
func (c *CodeConsumerBase) AddConstant(newcode string) {
	c.self.Add(newcode)
}

// Java: CodeConsumer.isWordChar(char)
func CodeConsumer_IsWordChar(ch uint16) bool {
	return ch == '_' ||
		ch == '$' ||
		jutil.CharacterIsLetter(ch) || jutil.CharacterIsDigit(ch)
}

// ShouldPreserveExtras: if the body of a for loop or the then clause of an
// if statement has a single statement, should it be wrapped in a block?
// Doing so can help when pretty-printing the code, and permits putting a
// debugging breakpoint on the statement inside the condition.
//
// Returns true if such expressions should be wrapped.
//
// Java: CodeConsumer.shouldPreserveExtras(Node)
func (c *CodeConsumerBase) ShouldPreserveExtras(n *rhino.Node) bool {
	return false
}

// MaybeInsertSpace allows a consumer to insert spaces in locations where
// it is unnecessary but may improve the readability of the code. This will
// be called in such places as after a statement and before opening
// parentheses, or after the end of a if block before the start of an else
// block.
//
// Java: CodeConsumer.maybeInsertSpace()
func (c *CodeConsumerBase) MaybeInsertSpace() {}

// BreakAfterBlockFor returns whether the a line break can be added after
// the specified BLOCK.
//
// Java: CodeConsumer.breakAfterBlockFor(Node, boolean)
func (c *CodeConsumerBase) BreakAfterBlockFor(n *rhino.Node, statementContext bool) bool {
	return statementContext
}

// EndFile is called when we're at the end of a file.
//
// Java: CodeConsumer.endFile()
func (c *CodeConsumerBase) EndFile() {}
