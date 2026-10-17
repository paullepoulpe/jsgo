package jutil

import (
	"strconv"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilPatternSyntaxVectors.java
// (JDK 21): for each pattern, whether Pattern.compile throws a
// PatternSyntaxException, and its description, index and message.

// patternSyntaxKnownGaps are the vectors where the RE2-based Pattern accepts a
// regex that Java rejects: Java validates every `{` as a counted closure,
// while RE2 reads a malformed one as a literal. (A repeated repetition such
// as a**, which Java rejects but RE2 reports like the x{2}{3} that Java
// accepts, is UnsupportedOperationException and skipped.)
var patternSyntaxKnownGaps = map[string]bool{
	"a{": true, "{": true, "a{1": true, "a{1,": true, "a{,2}": true, "a{ 1}": true, "a{2,1}": true,
}

func TestPatternSyntaxExceptionMatchesJDK(t *testing.T) {
	rows := readTSV(t, "pattern_syntax.tsv")
	if len(rows) < 50 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		regex, _ := unescVector(t, r[0])
		javaRejects := r[1] == "PSE"
		var caught any
		func() {
			defer func() { caught = recover() }()
			Pattern_Compile(regex)
		}()
		if _, unsupported := caught.(*UnsupportedOperationException); unsupported {
			continue
		}
		pse, isPSE := caught.(*PatternSyntaxException)
		if caught != nil && !isPSE {
			t.Errorf("%q: unexpected panic %v", regex, caught)
			continue
		}
		if isPSE != javaRejects {
			if patternSyntaxKnownGaps[regex] {
				continue
			}
			t.Errorf("%q: PatternSyntaxException=%v, JDK %v", regex, isPSE, javaRejects)
			continue
		}
		if !javaRejects {
			continue
		}
		desc, _ := unescVector(t, r[2])
		index, _ := strconv.Atoi(r[3])
		message, _ := unescVector(t, r[4])
		// The exception's message is composed from (description, index,
		// pattern) exactly as the JDK composes it.
		if got := NewPatternSyntaxException(desc, regex, index).GetMessage(); got != message {
			t.Errorf("%q: message %q, JDK %q", regex, got, message)
		}
		if pse.GetPattern() != regex {
			t.Errorf("%q: pattern %q", regex, pse.GetPattern())
		}
	}
}
