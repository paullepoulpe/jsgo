package jutil

import (
	"regexp"
	"strings"
)

// Go-only: language-adapter — FileSystems.getDefault().getPathMatcher(
// "glob:...") of the Unix file system, written from the Javadoc of
// FileSystem.getPathMatcher: the glob is matched against the whole
// Path.toString(), case-sensitively, and
//   - `*` matches zero or more characters of a name component (not '/');
//   - `**` matches zero or more characters across directory boundaries;
//   - `?` matches exactly one character of a name component;
//   - `\` escapes the following character;
//   - `[...]` matches one character of a name component out of a set: `!`
//     first negates it, `-` forms a range (and matches itself first or
//     last), and `*`, `?` and `\` match themselves;
//   - `{a,b,...}` matches any of the comma-separated subpatterns; groups
//     cannot be nested;
//   - a leading period is not special.
// The glob is compiled to an RE2 regular expression, which RE2 and Java
// match alike for these constructs.

// PathMatcher is java.nio.file.PathMatcher.
type PathMatcher struct {
	re *regexp.Regexp
}

// FileSystem_GetPathMatcher is FileSystem.getPathMatcher(String) for the
// "glob:" syntax (the only one Closure uses). A missing ':' throws
// IllegalArgumentException, another syntax UnsupportedOperationException,
// and a malformed glob PatternSyntaxException.
func FileSystem_GetPathMatcher(syntaxAndPattern string) *PathMatcher {
	pos := strings.IndexByte(syntaxAndPattern, ':')
	if pos <= 0 {
		panic(NewIllegalArgumentExceptionMsg(""))
	}
	syntax, input := syntaxAndPattern[:pos], syntaxAndPattern[pos+1:]
	if !strings.EqualFold(syntax, "glob") {
		panic(NewUnsupportedOperationExceptionMsg("Syntax '" + syntax + "' not recognized"))
	}
	return &PathMatcher{re: regexp.MustCompile(pathMatcher_globToRegex(input))}
}

// Matches is PathMatcher.matches(Path).
func (m *PathMatcher) Matches(path *Path) bool {
	return m.re.MatchString(path.String())
}

func pathMatcher_globToRegex(glob string) string {
	var sb strings.Builder
	sb.WriteString(`(?s)^`)
	groupOpen := false
	n := len(glob)
	for i := 0; i < n; i++ {
		c := glob[i]
		switch c {
		case '\\':
			if i+1 == n {
				panic(NewPatternSyntaxException("No character to escape", glob, i))
			}
			i++
			sb.WriteString(regexp.QuoteMeta(pathMatcher_charAt(glob, &i)))
		case '*':
			if i+1 < n && glob[i+1] == '*' {
				sb.WriteString(`.*`)
				i++
			} else {
				sb.WriteString(`[^/]*`)
			}
		case '?':
			sb.WriteString(`[^/]`)
		case '[':
			i = pathMatcher_bracket(glob, i, &sb)
		case '{':
			if groupOpen {
				panic(NewPatternSyntaxException("Cannot nest groups", glob, i))
			}
			groupOpen = true
			sb.WriteString(`(?:`)
		case '}':
			if groupOpen {
				groupOpen = false
				sb.WriteString(`)`)
			} else {
				sb.WriteString(`\}`)
			}
		case ',':
			if groupOpen {
				sb.WriteString(`|`)
			} else {
				sb.WriteString(`,`)
			}
		default:
			sb.WriteString(regexp.QuoteMeta(pathMatcher_charAt(glob, &i)))
		}
	}
	if groupOpen {
		panic(NewPatternSyntaxException("Missing '}", glob, n-1))
	}
	sb.WriteString(`$`)
	return sb.String()
}

// pathMatcher_charAt returns the (UTF-8 encoded) character at *i, leaving
// *i on its last byte.
func pathMatcher_charAt(s string, i *int) string {
	j := *i + 1
	for j < len(s) && s[j]&0xC0 == 0x80 {
		j++
	}
	r := s[*i:j]
	*i = j - 1
	return r
}

// pathMatcher_bracket translates the bracket expression starting at
// glob[start] == '[' and returns the index of its closing ']'. The class
// never matches '/'.
func pathMatcher_bracket(glob string, start int, sb *strings.Builder) int {
	n := len(glob)
	i := start + 1
	negate := false
	if i < n && glob[i] == '!' {
		negate = true
		i++
	}
	type rng struct{ lo, hi rune }
	var ranges []rng
	for ; i < n && glob[i] != ']'; i++ {
		loIndex := i
		lo := []rune(pathMatcher_charAt(glob, &i))[0]
		if lo == '/' {
			panic(NewPatternSyntaxException("Explicit 'name separator' in class", glob, i))
		}
		hi := lo
		if i+2 < n && glob[i+1] == '-' && glob[i+2] != ']' {
			i += 2
			hi = []rune(pathMatcher_charAt(glob, &i))[0]
			if hi < lo {
				panic(NewPatternSyntaxException("Invalid range", glob, loIndex))
			}
		}
		ranges = append(ranges, rng{lo, hi})
	}
	if i >= n {
		panic(NewPatternSyntaxException("Missing ']", glob, n-1))
	}
	sb.WriteString(`[`)
	if negate {
		sb.WriteString(`^/`)
	}
	for _, r := range ranges {
		for _, part := range []rng{{r.lo, min(r.hi, '/'-1)}, {max(r.lo, '/'+1), r.hi}} {
			if part.lo > part.hi {
				continue
			}
			sb.WriteString(pathMatcher_classEscape(part.lo))
			if part.hi != part.lo {
				sb.WriteString(`-`)
				sb.WriteString(pathMatcher_classEscape(part.hi))
			}
		}
	}
	if len(ranges) == 0 {
		// Go-only: language-adapter — Java's translation of an empty class
		// ("[]" or "[!]") does not compile, and its message quotes Java's
		// internal regex; the index here is the glob's ']'.
		panic(NewPatternSyntaxException("Unclosed character class", glob, i))
	}
	sb.WriteString(`]`)
	return i
}

// pathMatcher_classEscape writes r for use inside an RE2 character class.
func pathMatcher_classEscape(r rune) string {
	switch r {
	case '\\', ']', '[', '^', '-':
		return `\` + string(r)
	}
	return string(r)
}
