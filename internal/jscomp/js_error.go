// Port of com/google/javascript/jscomp/JSError.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/CheckLevel.java @ closure-compiler v20261028

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
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JSError is a compile error description.
//
// type: a type of the error. description: description of the error.
// sourceName: name of the source (nullable). lineno: one-indexed line number
// of the error location. charno: zero-indexed character number of the error
// location. length: length of the error region. node: node where the
// warning occurred (nullable). defaultLevel: the default level, before any
// of the WarningsGuards are applied. requirement: requirement that fails in
// the case of conformance violations (nullable).
type JSError struct {
	type_         *DiagnosticType
	requirement   *Requirement
	description   string
	sourceName    string
	lineno        int
	charno        int
	length        int
	node          rhino.Node
	defaultLevel  CheckLevel
	hasSourceName bool
}

// jsErrorPinsDeferred is non-zero while worker goroutines may report on
// nodes other workers read (fast mode's parallel TypeCheck reports on
// declarations in other scripts): PinForDiagnostic writes the node's
// record, so the errors are pinned when their reports are replayed on the
// main goroutine instead (replayCompilerEvents). One compilation runs at a
// time per process (rhino/node_store.go).
var jsErrorPinsDeferred atomic.Int32

func newJSError(
	type_ *DiagnosticType,
	description string,
	sourceName string,
	hasSourceName bool,
	lineno int,
	charno int,
	length int,
	node rhino.Node,
	defaultLevel CheckLevel,
	requirement *Requirement,
) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	if !node.IsNil() && jsErrorPinsDeferred.Load() == 0 {
		node.PinForDiagnostic()
	}
	return &JSError{
		type_:         type_,
		description:   description,
		sourceName:    sourceName,
		hasSourceName: hasSourceName,
		lineno:        lineno,
		charno:        charno,
		length:        length,
		node:          node,
		defaultLevel:  defaultLevel,
		requirement:   requirement,
	}
}

// Java: JSError(DiagnosticType, String, String, int, int, int, Node, CheckLevel, Requirement)
func NewJSError(
	type_ *DiagnosticType,
	description string,
	sourceName *string,
	lineno int,
	charno int,
	length int,
	node rhino.Node,
	defaultLevel CheckLevel,
	requirement *Requirement) *JSError {
	var src string
	if sourceName != nil {
		src = *sourceName
	}
	return newJSError(type_, description, src, sourceName != nil, lineno, charno, length, node, defaultLevel, requirement)
}

// Java: JSError.type()
func (e *JSError) Type() *DiagnosticType { return e.type_ }

// Java: JSError.description()
func (e *JSError) Description() string { return e.description }

// SourceName returns the source file name and whether one is present.
//
// Java: JSError.sourceName()
func (e *JSError) SourceName() (string, bool) { return e.sourceName, e.hasSourceName }

// Java: JSError.lineno()
func (e *JSError) Lineno() int { return e.lineno }

// Java: JSError.charno()
func (e *JSError) Charno() int { return e.charno }

// Java: JSError.length()
func (e *JSError) Length() int { return e.length }

// Java: JSError.node()
func (e *JSError) Node() rhino.Node { return e.node }

// Java: JSError.defaultLevel()
func (e *JSError) DefaultLevel() CheckLevel { return e.defaultLevel }

// Java: JSError.requirement()
func (e *JSError) Requirement() *Requirement { return e.requirement }

// Java: JSError.getType()
func (e *JSError) GetType() *DiagnosticType { return e.Type() }

// Java: JSError.getDescription()
func (e *JSError) GetDescription() string { return e.Description() }

// Java: JSError.getLineno()
func (e *JSError) GetLineno() int { return e.Lineno() }

// Java: JSError.getCharno()
func (e *JSError) GetCharno() int { return e.Charno() }

// Java: JSError.getLength()
func (e *JSError) GetLength() int { return e.Length() }

// Java: JSError.getNode()
func (e *JSError) GetNode() rhino.Node { return e.Node() }

// Java: JSError.getDefaultLevel()
func (e *JSError) GetDefaultLevel() CheckLevel { return e.DefaultLevel() }

// Java: JSError.getRequirement()
func (e *JSError) GetRequirement() *Requirement { return e.Requirement() }

const (
	jSError_DEFAULT_LINENO = -1
	jSError_DEFAULT_CHARNO = -1
	jSError_DEFAULT_LENGTH = 0
)

// DEFAULT_SOURCENAME, DEFAULT_NODE and DEFAULT_REQUIREMENT are null.

// JSError_Make creates a JSError with no source information.
//
// Java: JSError.make(DiagnosticType, String...)
func JSError_Make(type_ *DiagnosticType, arguments ...string) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	return newJSError(
		type_, type_.Format(arguments...), "", false, jSError_DEFAULT_LINENO, jSError_DEFAULT_CHARNO, jSError_DEFAULT_LENGTH, rhino.NilNode, type_.level, nil)
}

// JSError_MakeWithSourceLocation creates a JSError at a given source
// location: lineno is the line number with source file, or -1 if unknown;
// charno the column number within line, or -1 for whole line.
//
// Java: JSError.make(String, int, int, DiagnosticType, String...)
func JSError_MakeWithSourceLocation(
	sourceName *string, lineno int, charno int, type_ *DiagnosticType, arguments ...string) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	return NewJSError(
		type_, type_.Format(arguments...), sourceName, lineno, charno, jSError_DEFAULT_LENGTH, rhino.NilNode, type_.level, nil)
}

// JSError_MakeAtSourceLocation creates a JSError at a non-null source location
// without boxing sourceName into a *string.
func JSError_MakeAtSourceLocation(
	sourceName string, lineno int, charno int, type_ *DiagnosticType, arguments ...string) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	return newJSError(
		type_, type_.Format(arguments...), sourceName, true, lineno, charno, jSError_DEFAULT_LENGTH, rhino.NilNode, type_.level, nil)
}

// JSError_MakeWithNode creates a JSError from a file and Node position.
//
// Java: JSError.make(Node, DiagnosticType, String...)
func JSError_MakeWithNode(n rhino.Node, type_ *DiagnosticType, arguments ...string) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	return NewJSError(
		type_, type_.Format(arguments...), n.GetSourceFileName(), n.GetLineno(), n.GetCharno(), n.GetLength(), n, type_.level, nil)
}

// JSError_MakeWithNodeRange creates a JSError from a file and Node position:
// start determines the line and char position and source file name, end
// determines the length of the error region.
//
// Java: JSError.make(Node, Node, DiagnosticType, String...)
func JSError_MakeWithNodeRange(start rhino.Node, end rhino.Node, type_ *DiagnosticType, arguments ...string) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	length := jSError_DEFAULT_LENGTH
	endOffset := end.GetSourceOffset()
	startOffset := start.GetSourceOffset()
	if endOffset != -1 && startOffset != -1 {
		length = (endOffset + end.GetLength()) - startOffset
	}
	return NewJSError(
		type_, type_.Format(arguments...), start.GetSourceFileName(), start.GetLineno(), start.GetCharno(), length, start, type_.level, nil)
}

// JSError_MakeWithRequirement creates a JSError from a requirement, file and
// Node position.
//
// Java: JSError.make(Requirement, Node, DiagnosticType, String...)
func JSError_MakeWithRequirement(
	requirement *Requirement, n rhino.Node, type_ *DiagnosticType, arguments ...string) *JSError {
	if type_ == nil {
		panic(jutil.NewNullPointerExceptionMsg("type"))
	}
	jutil.CheckNotNull(requirement)
	return NewJSError(
		type_, type_.Format(arguments...), n.GetSourceFileName(), n.GetLineno(), n.GetCharno(), n.GetLength(), n, type_.level, requirement)
}

// JSErrorBuilder builds a JSError (JSError.Builder).
type JSErrorBuilder struct {
	type_ *DiagnosticType
	args  []string

	sourceName    string
	lineno        int
	charno        int
	length        int
	n             rhino.Node
	level         CheckLevel
	hasSourceName bool
}

// NewJSErrorBuilder creates a new builder (JSError.builder): type_ is the
// associated DiagnosticType; args are the formatting arguments used to
// format the DiagnosticType's description.
//
// Java: JSError.builder(DiagnosticType, String...)
func NewJSErrorBuilder(type_ *DiagnosticType, args ...string) *JSErrorBuilder {
	return &JSErrorBuilder{
		type_:  type_,
		args:   args,
		level:  type_.level, // may be overwritten later
		lineno: jSError_DEFAULT_LINENO,
		charno: jSError_DEFAULT_CHARNO,
		length: jSError_DEFAULT_LENGTH,
	}
}

// SetNode sets the location where this error occurred.
//
// Incompatible with SetSourceLocation.
//
// Java: JSError.Builder.setNode(Node)
func (b *JSErrorBuilder) SetNode(n rhino.Node) *JSErrorBuilder {
	jutil.CheckStateMsg(!b.hasSourceName, "Cannot provide a Node when there's already a source name")
	b.n = n
	if src := n.GetSourceFileName(); src != nil {
		b.sourceName = *src
		b.hasSourceName = true
	}
	b.lineno = n.GetLineno()
	b.charno = n.GetCharno()
	b.length = n.GetLength()
	return b
}

// SetLevel sets the default level, before any WarningGuards, of this
// JSError. Overwrites the default level of the DiagnosticType.
//
// This method is rarely useful: prefer whenever possible to rely on the
// default level of the DiagnosticType.
//
// Java: JSError.Builder.setLevel(CheckLevel)
func (b *JSErrorBuilder) SetLevel(level CheckLevel) *JSErrorBuilder {
	b.level = level
	return b
}

// SetSourceLocation sets the location where this error occurred: lineno is
// the line number with source file, or -1 if unknown; charno the column
// number within line, or -1 for whole line.
//
// Incompatible with SetNode.
//
// Java: JSError.Builder.setSourceLocation(String, int, int)
func (b *JSErrorBuilder) SetSourceLocation(sourceName *string, lineno int, charno int) *JSErrorBuilder {
	jutil.CheckStateMsg(b.n.IsNil(), "Cannot provide a source location when there is already a Node")
	if sourceName != nil {
		b.sourceName = *sourceName
		b.hasSourceName = true
	} else {
		b.sourceName = ""
		b.hasSourceName = false
	}
	b.lineno = lineno
	b.charno = charno
	return b
}

// Java: JSError.Builder.build()
func (b *JSErrorBuilder) Build() *JSError {
	return newJSError(
		b.type_, b.type_.Format(b.args...), b.sourceName, b.hasSourceName, b.lineno, b.charno, b.length, b.n, b.level, nil)
}

// String returns the default rendering of an error as text.
//
// Java: JSError.toString()
func (e *JSError) String() string {
	sourceName := "(unknown source)"
	if e.hasSourceName && e.sourceName != "" {
		sourceName = e.sourceName
	}
	lineno := "(unknown line)"
	if e.Lineno() != jSError_DEFAULT_LINENO {
		lineno = strconv.Itoa(e.Lineno())
	}
	charno := "(unknown column)"
	if e.Charno() != jSError_DEFAULT_CHARNO {
		charno = strconv.Itoa(e.Charno())
	}

	return e.Type().key +
		". " +
		e.Description() +
		" at " +
		sourceName +
		" line " +
		lineno +
		" : " +
		charno
}

// Format formats a message at the given level, returning ("", false) for
// levels other than ERROR and WARNING.
//
// Java: JSError.format(CheckLevel, MessageFormatter)
func (e *JSError) Format(level CheckLevel, formatter MessageFormatter) (string, bool) {
	switch level {
	case CheckLevel_ERROR:
		return formatter.FormatError(e), true
	case CheckLevel_WARNING:
		return formatter.FormatWarning(e), true
	default:
		return "", false
	}
}

// GetNodeSourceOffset returns the offset of the region the Error applies
// to, or -1 if the offset is unknown.
//
// Java: JSError.getNodeSourceOffset()
func (e *JSError) GetNodeSourceOffset() int {
	if !e.Node().IsNil() {
		return e.Node().GetSourceOffset()
	}
	return -1
}

// GetLineNumber is an alias for GetLineno().
//
// Java: JSError.getLineNumber()
func (e *JSError) GetLineNumber() int {
	return e.Lineno()
}

// Java: JSError.equals(Object) (record)
func (e *JSError) Equals(o any) bool {
	that, ok := o.(*JSError)
	if !ok || that == nil {
		return false
	}
	return e.type_.Equals(that.type_) &&
		e.description == that.description &&
		e.hasSourceName == that.hasSourceName &&
		e.sourceName == that.sourceName &&
		e.lineno == that.lineno &&
		e.charno == that.charno &&
		e.length == that.length &&
		e.node == that.node &&
		e.defaultLevel == that.defaultLevel &&
		jutil.ObjectsEquals(e.requirement, that.requirement)
}

// Java: JSError.hashCode() (record)
//
// Go-only: language-adapter — the Node and the CheckLevel hash by identity
// in Java, which varies from run to run.
func (e *JSError) HashCode() int32 {
	var h int32
	h = 31*h + e.type_.HashCode()
	h = 31*h + jutil.StringHashCode(e.description)
	if e.hasSourceName {
		h = 31*h + jutil.StringHashCode(e.sourceName)
	} else {
		h = 31 * h
	}
	h = 31*h + int32(e.lineno)
	h = 31*h + int32(e.charno)
	h = 31*h + int32(e.length)
	h = 31*h + jutil.IdentityHashCode(e.node)
	h = 31*h + int32(e.defaultLevel)
	if e.requirement != nil {
		h = 31*h + e.requirement.HashCode()
	} else {
		h = 31 * h
	}
	return h
}

// CheckLevel controls checking levels of certain options. For all checks
// going forward, this should be used instead of booleans, so teams and
// individuals can control which checks are off, which produce only warnings,
// and which produce errors, without everyone having to agree.
type CheckLevel int32

const (
	CheckLevel_ERROR CheckLevel = iota
	CheckLevel_WARNING
	CheckLevel_OFF
)

var checkLevel_names = [...]string{"ERROR", "WARNING", "OFF"}

// String is Enum.name().
func (l CheckLevel) String() string { return checkLevel_names[l] }

// Java: CheckLevel.isOn()
func (l CheckLevel) IsOn() bool {
	return l != CheckLevel_OFF
}
