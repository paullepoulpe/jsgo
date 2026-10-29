// Port of com/google/javascript/jscomp/DiagnosticType.java @ closure-compiler v20261021

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
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DiagnosticType is the type of a compile or analysis error.
type DiagnosticType struct {
	// The error type. Used as the BugPattern and BugInstance types by
	// BugBot's XML
	key string

	// The default way to format errors. The style of format is
	// java.text.MessageFormat.
	format string

	// The default reporting level for this diagnostic
	level CheckLevel

	msgFormat atomic.Pointer[diagPattern]
}

// Java: DiagnosticType.serialVersionUID
const DiagnosticType_serialVersionUID int64 = 1

// DiagnosticType_Error creates a DiagnosticType at level CheckLevel.ERROR.
//
// Java: DiagnosticType.error(String, String)
func DiagnosticType_Error(name string, descriptionFormat string) *DiagnosticType {
	return DiagnosticType_Make(name, CheckLevel_ERROR, descriptionFormat)
}

// DiagnosticType_Warning creates a DiagnosticType at level
// CheckLevel.WARNING.
//
// Java: DiagnosticType.warning(String, String)
func DiagnosticType_Warning(name string, descriptionFormat string) *DiagnosticType {
	return DiagnosticType_Make(name, CheckLevel_WARNING, descriptionFormat)
}

// DiagnosticType_Disabled creates a DiagnosticType at level CheckLevel.OFF.
//
// Java: DiagnosticType.disabled(String, String)
func DiagnosticType_Disabled(name string, descriptionFormat string) *DiagnosticType {
	return DiagnosticType_Make(name, CheckLevel_OFF, descriptionFormat)
}

// DiagnosticType_Make creates a DiagnosticType at a given CheckLevel.
//
// Java: DiagnosticType.make(String, CheckLevel, String)
func DiagnosticType_Make(name string, level CheckLevel, descriptionFormat string) *DiagnosticType {
	return newDiagnosticType(name, level, descriptionFormat)
}

// Create a DiagnosticType. Private to force use of static factory methods.
func newDiagnosticType(key string, level CheckLevel, format string) *DiagnosticType {
	return &DiagnosticType{key: key, level: level, format: format}
}

// Java: DiagnosticType.format(String...)
func (t *DiagnosticType) Format(arguments ...string) string {
	mf := t.msgFormat.Load()
	if mf == nil {
		mf = parseDiagPattern(t.format)
		t.msgFormat.Store(mf)
	}
	return mf.format(arguments...)
}

type diagElement struct {
	prefix string
	index  int
}

type diagPattern struct {
	elements []diagElement
	tail     string
}

func parseDiagPattern(pattern string) *diagPattern {
	if strings.IndexByte(pattern, '\'') < 0 && strings.IndexByte(pattern, '{') < 0 {
		return &diagPattern{tail: pattern}
	}
	var elements []diagElement
	var lit strings.Builder
	quoted := false
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		if c == '\'' {
			if i+1 < len(pattern) && pattern[i+1] == '\'' {
				lit.WriteByte('\'')
				i++
			} else {
				quoted = !quoted
			}
			continue
		}
		if c == '{' && !quoted {
			closeIdx, level, inQuote := -1, 0, false
			for j := i + 1; j < len(pattern); j++ {
				switch ch := pattern[j]; {
				case ch == '\'':
					inQuote = !inQuote
				case inQuote:
				case ch == '{':
					level++
				case ch == '}' && level > 0:
					level--
				case ch == '}':
					closeIdx = j
				}
				if closeIdx >= 0 {
					break
				}
			}
			if closeIdx < 0 {
				panic(jutil.NewIllegalArgumentExceptionMsg("Unmatched braces in the pattern."))
			}
			idxStr := pattern[i+1 : closeIdx]
			if comma := strings.IndexByte(idxStr, ','); comma >= 0 {
				idxStr = idxStr[:comma]
			}
			idx, err := strconv.Atoi(idxStr)
			if err != nil || strings.HasPrefix(idxStr, "+") {
				panic(jutil.NewIllegalArgumentExceptionMsg("can't parse argument number: " + idxStr))
			}
			if idx < 0 {
				panic(jutil.NewIllegalArgumentExceptionMsg("negative argument number: " + strconv.Itoa(idx)))
			}
			if idx >= 10000 {
				panic(jutil.NewIllegalArgumentExceptionMsg(strconv.Itoa(idx) + " exceeds the ArgumentIndex implementation limit"))
			}
			elements = append(elements, diagElement{prefix: lit.String(), index: idx})
			lit.Reset()
			i = closeIdx
			continue
		}
		lit.WriteByte(c)
	}
	return &diagPattern{elements: elements, tail: lit.String()}
}

func (p *diagPattern) format(arguments ...string) string {
	if len(p.elements) == 0 {
		return p.tail
	}
	n := len(p.tail)
	for i := range p.elements {
		e := &p.elements[i]
		n += len(e.prefix)
		if e.index < len(arguments) {
			n += len(arguments[e.index])
		} else {
			n += 3
		}
	}
	var sb strings.Builder
	sb.Grow(n)
	for i := range p.elements {
		e := &p.elements[i]
		sb.WriteString(e.prefix)
		if e.index < len(arguments) {
			sb.WriteString(arguments[e.index])
		} else {
			sb.WriteByte('{')
			sb.WriteString(strconv.Itoa(e.index))
			sb.WriteByte('}')
		}
	}
	sb.WriteString(p.tail)
	return sb.String()
}

// Java: DiagnosticType.equals(Object)
func (t *DiagnosticType) Equals(typ any) bool {
	diagnosticType, ok := typ.(*DiagnosticType)
	return ok && diagnosticType != nil && diagnosticType.key == t.key
}

// Java: DiagnosticType.hashCode()
func (t *DiagnosticType) HashCode() int32 {
	return jutil.StringHashCode(t.key)
}

// Java: DiagnosticType.compareTo(DiagnosticType)
func (t *DiagnosticType) CompareTo(diagnosticType *DiagnosticType) int {
	return jutil.StringCompareTo(t.key, diagnosticType.key)
}

// Java: DiagnosticType.toString()
func (t *DiagnosticType) String() string {
	return t.key + ": " + t.format
}

// Key returns the diagnostic's unique identifier (Java's DiagnosticType.key).
func (t *DiagnosticType) Key() string {
	return t.key
}

// FormatPattern returns the diagnostic's MessageFormat template (Java's DiagnosticType.format).
func (t *DiagnosticType) FormatPattern() string {
	return t.format
}
