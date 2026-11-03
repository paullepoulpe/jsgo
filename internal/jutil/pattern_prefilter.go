package jutil

import (
	"cmp"
	"regexp"
	"regexp/syntax"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	matchesCoreNone uint8 = iota
	matchesCoreDotAll
	matchesCoreJavaDot
)

const (
	maxPrefilterAllOf = 4
	maxPrefilterAnyOf = 8
)

// literalPrefilter holds mandatory exact UTF-8 substrings extracted from a
// compiled RE2 syntax tree. Every match of the regex must contain all strings
// in allOf and (when anyOf is non-empty) at least one string in anyOf.
type literalPrefilter struct {
	allOf []string
	anyOf []string
}

func (f literalPrefilter) match(s string) bool {
	for _, lit := range f.allOf {
		if !strings.Contains(s, lit) {
			return false
		}
	}
	if len(f.anyOf) > 0 {
		for _, lit := range f.anyOf {
			if strings.Contains(s, lit) {
				return true
			}
		}
		return false
	}
	return true
}

var javaDotClassRunes = mustParseJavaDotRunes()

func mustParseJavaDotRunes() []rune {
	re, err := syntax.Parse(javaDotClass, syntax.Perl)
	if err != nil || re.Op != syntax.OpCharClass {
		panic("jutil: failed to parse javaDotClass")
	}
	return re.Rune
}

func populatePatternFastPaths(p *Pattern, whole, search00 string) {
	if wholeAST, err := syntax.Parse(whole, syntax.Perl); err == nil {
		p.fullPrefilter = buildLiteralPrefilter(wholeAST)
		p.initMatchesCore(wholeAST)
	}
	if search00 == whole {
		p.findPrefilter = p.fullPrefilter
	} else if findAST, err := syntax.Parse(search00, syntax.Perl); err == nil {
		p.findPrefilter = buildLiteralPrefilter(findAST)
	}
}

func unwrapCaptures(re *syntax.Regexp) *syntax.Regexp {
	for re.Op == syntax.OpCapture && len(re.Sub) == 1 {
		re = re.Sub[0]
	}
	return re
}

func isDotAllStar(re *syntax.Regexp) bool {
	re = unwrapCaptures(re)
	return re.Op == syntax.OpStar && len(re.Sub) == 1 && unwrapCaptures(re.Sub[0]).Op == syntax.OpAnyChar
}

func isJavaDotStar(re *syntax.Regexp) bool {
	re = unwrapCaptures(re)
	if re.Op != syntax.OpStar || len(re.Sub) != 1 {
		return false
	}
	sub := unwrapCaptures(re.Sub[0])
	return sub.Op == syntax.OpCharClass && slices.Equal(sub.Rune, javaDotClassRunes)
}

func (p *Pattern) initMatchesCore(wholeAST *syntax.Regexp) {
	ast := unwrapCaptures(wholeAST)
	if isDotAllStar(ast) {
		p.matchesCoreMode = matchesCoreDotAll
		return
	}
	if isJavaDotStar(ast) {
		p.matchesCoreMode = matchesCoreJavaDot
		return
	}
	if ast.Op != syntax.OpConcat || len(ast.Sub) < 2 {
		return
	}
	subs := ast.Sub

	leadDotAll := isDotAllStar(subs[0])
	trailDotAll := isDotAllStar(subs[len(subs)-1])
	if leadDotAll || trailDotAll {
		start, end := 0, len(subs)
		for start < end && isDotAllStar(subs[start]) {
			start++
		}
		for end > start && isDotAllStar(subs[end-1]) {
			end--
		}
		if start == end {
			p.matchesCoreMode = matchesCoreDotAll
			return
		}
		if re, ok := compileTrimmedCore(subs[start:end], start > 0, end < len(subs)); ok {
			p.matchesCore = re
			p.matchesCoreMode = matchesCoreDotAll
		}
		return
	}

	leadJavaDot := isJavaDotStar(subs[0])
	trailJavaDot := isJavaDotStar(subs[len(subs)-1])
	if leadJavaDot || trailJavaDot {
		start, end := 0, len(subs)
		for start < end && isJavaDotStar(subs[start]) {
			start++
		}
		for end > start && isJavaDotStar(subs[end-1]) {
			end--
		}
		if start == end {
			p.matchesCoreMode = matchesCoreJavaDot
			return
		}
		mid := subs[start:end]
		if re, ok := compileTrimmedCore(mid, start > 0, end < len(subs)); ok {
			coreNode := &syntax.Regexp{Op: syntax.OpConcat, Sub: mid}
			p.matchesCore = re
			p.matchesCoreMode = matchesCoreJavaDot
			p.matchesCoreCanMatchTerm = syntaxCanMatchJavaTerminator(coreNode)
		}
	}
}

func compileTrimmedCore(mid []*syntax.Regexp, trimmedLead, trimmedTrail bool) (*regexp.Regexp, bool) {
	coreNode := &syntax.Regexp{Op: syntax.OpConcat, Sub: mid}
	expr := `(?:` + coreNode.String() + `)`
	if !trimmedLead {
		expr = `\A` + expr
	}
	if !trimmedTrail {
		expr = expr + `\z`
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, false
	}
	return re, true
}

func syntaxCanMatchJavaTerminator(re *syntax.Regexp) bool {
	switch re.Op {
	case syntax.OpNoMatch, syntax.OpEmptyMatch,
		syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText,
		syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		return false
	case syntax.OpAnyChar, syntax.OpAnyCharNotNL:
		return true
	case syntax.OpLiteral:
		for _, r := range re.Rune {
			if isJavaLineTerminatorRune(r) {
				return true
			}
		}
		return false
	case syntax.OpCharClass:
		for i := 0; i+1 < len(re.Rune); i += 2 {
			lo, hi := re.Rune[i], re.Rune[i+1]
			if charRangeContainsJavaTerminator(lo, hi) {
				return true
			}
		}
		return false
	case syntax.OpCapture, syntax.OpStar, syntax.OpPlus, syntax.OpQuest:
		return len(re.Sub) > 0 && syntaxCanMatchJavaTerminator(re.Sub[0])
	case syntax.OpRepeat:
		if re.Max == 0 || len(re.Sub) == 0 {
			return false
		}
		return syntaxCanMatchJavaTerminator(re.Sub[0])
	case syntax.OpConcat, syntax.OpAlternate:
		for _, sub := range re.Sub {
			if syntaxCanMatchJavaTerminator(sub) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

func charRangeContainsJavaTerminator(lo, hi rune) bool {
	return (lo <= '\n' && '\n' <= hi) ||
		(lo <= '\r' && '\r' <= hi) ||
		(lo <= '\u0085' && '\u0085' <= hi) ||
		(lo <= '\u2028' && '\u2028' <= hi) ||
		(lo <= '\u2029' && '\u2029' <= hi)
}

// containsJavaLineTerminator reports whether s contains any Java line
// terminator (\n, \r, \u0085, \u2028, or \u2029).
func containsJavaLineTerminator(s string) bool {
	if strings.IndexByte(s, '\n') >= 0 || strings.IndexByte(s, '\r') >= 0 {
		return true
	}
	// \u0085 is C2 85 in UTF-8; \u2028 and \u2029 are E2 80 A8 and E2 80 A9.
	if i := strings.IndexByte(s, 0xC2); i >= 0 {
		if strings.Contains(s[i:], "\u0085") {
			return true
		}
	}
	if i := strings.IndexByte(s, 0xE2); i >= 0 {
		sub := s[i:]
		if strings.Contains(sub, "\u2028") || strings.Contains(sub, "\u2029") {
			return true
		}
	}
	return false
}

func exactLiteralString(re *syntax.Regexp) (string, bool) {
	re = unwrapCaptures(re)
	if re.Op != syntax.OpLiteral || re.Flags&syntax.FoldCase != 0 || len(re.Rune) == 0 {
		return "", false
	}
	for _, r := range re.Rune {
		// U+FFFD in RE2 matches arbitrary invalid UTF-8 bytes, not just EF BF BD.
		if r == utf8.RuneError {
			return "", false
		}
	}
	return string(re.Rune), true
}

func buildLiteralPrefilter(re *syntax.Regexp) literalPrefilter {
	f := extractPrefilter(re)
	if len(f.allOf) > 1 {
		slices.SortFunc(f.allOf, func(a, b string) int {
			return cmp.Compare(len(b), len(a))
		})
		dedup := f.allOf[:0]
		for _, lit := range f.allOf {
			redundant := false
			for _, kept := range dedup {
				if strings.Contains(kept, lit) {
					redundant = true
					break
				}
			}
			if !redundant {
				dedup = append(dedup, lit)
				if len(dedup) == maxPrefilterAllOf {
					break
				}
			}
		}
		f.allOf = dedup
	}
	if len(f.anyOf) > 0 && len(f.allOf) > 0 {
		// Drop anyOf if one of its alternatives is already implied by allOf.
		for _, alt := range f.anyOf {
			for _, req := range f.allOf {
				if strings.Contains(req, alt) {
					f.anyOf = nil
					break
				}
			}
			if f.anyOf == nil {
				break
			}
		}
	}
	return f
}

func extractPrefilter(re *syntax.Regexp) literalPrefilter {
	re = unwrapCaptures(re)
	switch re.Op {
	case syntax.OpLiteral:
		if s, ok := exactLiteralString(re); ok {
			return literalPrefilter{allOf: []string{s}}
		}
		return literalPrefilter{}
	case syntax.OpPlus:
		if len(re.Sub) == 1 {
			return extractPrefilter(re.Sub[0])
		}
		return literalPrefilter{}
	case syntax.OpRepeat:
		if re.Min >= 1 && len(re.Sub) == 1 {
			return extractPrefilter(re.Sub[0])
		}
		return literalPrefilter{}
	case syntax.OpConcat:
		var allOf []string
		var bestAnyOf []string
		var runBuilder strings.Builder
		flushRun := func() {
			if runBuilder.Len() > 0 {
				allOf = append(allOf, runBuilder.String())
				runBuilder.Reset()
			}
		}
		for _, sub := range re.Sub {
			if lit, ok := exactLiteralString(sub); ok {
				runBuilder.WriteString(lit)
				continue
			}
			flushRun()
			sf := extractPrefilter(sub)
			allOf = append(allOf, sf.allOf...)
			if len(sf.anyOf) > 0 && (len(bestAnyOf) == 0 || len(sf.anyOf) < len(bestAnyOf)) {
				bestAnyOf = sf.anyOf
			}
		}
		flushRun()
		return literalPrefilter{allOf: allOf, anyOf: bestAnyOf}
	case syntax.OpAlternate:
		if len(re.Sub) == 0 || len(re.Sub) > maxPrefilterAnyOf {
			return literalPrefilter{}
		}
		var anyOf []string
		for _, branch := range re.Sub {
			bf := buildLiteralPrefilter(branch)
			switch {
			case len(bf.allOf) > 0:
				anyOf = append(anyOf, bf.allOf[0])
			case len(bf.anyOf) > 0:
				anyOf = append(anyOf, bf.anyOf...)
			default:
				return literalPrefilter{}
			}
			if len(anyOf) > maxPrefilterAnyOf {
				return literalPrefilter{}
			}
		}
		// Deduplicate exact entries in anyOf.
		dedup := anyOf[:0]
		for _, alt := range anyOf {
			if !slices.Contains(dedup, alt) {
				dedup = append(dedup, alt)
			}
		}
		if len(dedup) == 1 {
			return literalPrefilter{allOf: dedup}
		}
		return literalPrefilter{anyOf: dedup}
	default:
		return literalPrefilter{}
	}
}
