package jutil

import (
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.util.regex.PatternSyntaxException,
// written from its Javadoc ("Unchecked exception thrown to indicate a syntax
// error in a regular-expression pattern"; constructor (desc, regex, index),
// with index -1 if the index is not known; getDescription, getIndex,
// getPattern, and getMessage: "a multi-line string containing the
// description of the syntax error and its index, the erroneous
// regular-expression pattern, and a visual indication of the error index
// within the pattern"). The exact layout of that message was checked
// black-box against JDK 21 (testdata/pattern_syntax.tsv, from
// tools/oracle/java/JutilPatternSyntaxVectors.java).

// PatternSyntaxException is java.util.regex.PatternSyntaxException (an
// IllegalArgumentException).
type PatternSyntaxException struct {
	IllegalArgumentException
	desc    string
	pattern string
	index   int
}

func (*PatternSyntaxException) patternSyntaxException() {}

// IPatternSyntaxException matches PatternSyntaxException.
type IPatternSyntaxException interface {
	IIllegalArgumentException
	patternSyntaxException()
}

// NewPatternSyntaxException is new PatternSyntaxException(desc, regex,
// index); index is -1 if it is not known.
func NewPatternSyntaxException(desc string, regex string, index int) *PatternSyntaxException {
	e := &PatternSyntaxException{desc: desc, pattern: regex, index: index}
	var sb strings.Builder
	sb.WriteString(desc)
	if index >= 0 {
		sb.WriteString(" near index ")
		sb.WriteString(strconv.Itoa(index))
	}
	sb.WriteString("\n")
	sb.WriteString(regex)
	if index >= 0 && index < StringLength(regex) {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(" ", index))
		sb.WriteString("^")
	}
	InitThrowable(e, "java.util.regex.PatternSyntaxException", msgPtr(sb.String()), nil)
	return e
}

// GetDescription is PatternSyntaxException.getDescription().
func (e *PatternSyntaxException) GetDescription() string { return e.desc }

// GetIndex is PatternSyntaxException.getIndex() (-1 if not known).
func (e *PatternSyntaxException) GetIndex() int { return e.index }

// GetPattern is PatternSyntaxException.getPattern().
func (e *PatternSyntaxException) GetPattern() string { return e.pattern }

// compileTranslatedJavaRegex compiles translated, the RE2 form of the Java
// regex regex. A syntax error that RE2 reports is a Java syntax error too
// (the translation keeps the shared constructs' syntax), so it becomes a
// PatternSyntaxException, except for RE2's own limits and its stricter
// rules on repeated repetitions (Java accepts x{2}{3}, a{1001}): those are
// constructs this adapter cannot express (UnsupportedOperationException).
// RE2 does not report Java's index or description, so the exception has
// RE2's description and index -1.
func compileTranslatedJavaRegex(regex string, translated string) *regexp.Regexp {
	re, err := regexp.Compile(translated)
	if err == nil {
		return re
	}
	if se, ok := err.(*syntax.Error); ok {
		switch se.Code {
		case syntax.ErrInvalidRepeatSize, syntax.ErrInvalidRepeatOp, syntax.ErrLarge,
			syntax.ErrNestingDepth, syntax.ErrInternalError:
			unsupportedJavaRegex(regex, "RE2: "+se.Code.String())
		}
		panic(NewPatternSyntaxException(se.Code.String(), regex, -1))
	}
	panic(err)
}
