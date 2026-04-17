// extract_fixtures reads a closure-compiler Java test file and emits
// one .js fixture per test method whose single addSource(...) argument
// is a plain string or Java text block — i.e. the cases jsgo can
// plausibly handle.
//
// Intentionally conservative: tests that use addExterns, multiple
// addSource calls, string concatenation, or the lines(...) helper are
// skipped (logged to stderr). The resulting corpus is what we run
// through `make oracle` to get authoritative Java diagnostics.
//
// Usage:
//     go run ./tools/extract_fixtures <TypeCheckTest.java> <out-dir>

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: extract_fixtures <Java test file> <out dir>")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	outDir := os.Args[2]
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tests := extractTests(string(data))
	kept, skipped := 0, 0
	for _, t := range tests {
		if hasForbidden(t.body) {
			skipped++
			continue
		}
		js, ok := extractAddSource(t.body)
		if !ok {
			skipped++
			continue
		}
		name := sanitizeName(t.name)
		path := filepath.Join(outDir, name+".js")
		if err := os.WriteFile(path, []byte(js), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		kept++
	}
	fmt.Fprintf(os.Stderr, "extracted %d fixtures, skipped %d\n", kept, skipped)
}

type test struct {
	name string
	body string
}

// extractTests walks the source and returns every top-level method
// matching `public void test<Name>() { ... }`, via balanced-brace scan.
func extractTests(src string) []test {
	methodHead := regexp.MustCompile(`public\s+void\s+(test[A-Za-z0-9_]+)\s*\(\s*\)\s*\{`)
	idx := methodHead.FindAllStringSubmatchIndex(src, -1)
	var out []test
	for _, m := range idx {
		name := src[m[2]:m[3]]
		// m[1] is the offset of '{'. Scan to find matching '}'.
		open := m[1] - 1
		end := matchBrace(src, open)
		if end < 0 {
			continue
		}
		out = append(out, test{name: name, body: src[open+1 : end]})
	}
	return out
}

func matchBrace(src string, open int) int {
	depth := 0
	i := open
	for i < len(src) {
		c := src[i]
		switch c {
		case '{':
			depth++
			i++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
			i++
		case '/':
			if i+1 < len(src) && src[i+1] == '/' {
				for i < len(src) && src[i] != '\n' {
					i++
				}
			} else if i+1 < len(src) && src[i+1] == '*' {
				i += 2
				for i+1 < len(src) && !(src[i] == '*' && src[i+1] == '/') {
					i++
				}
				i += 2
			} else {
				i++
			}
		case '"':
			if i+2 < len(src) && src[i+1] == '"' && src[i+2] == '"' {
				i += 3
				for i+2 < len(src) && !(src[i] == '"' && src[i+1] == '"' && src[i+2] == '"') {
					i++
				}
				i += 3
			} else {
				i++
				for i < len(src) && src[i] != '"' {
					if src[i] == '\\' && i+1 < len(src) {
						i += 2
						continue
					}
					if src[i] == '\n' {
						break
					}
					i++
				}
				if i < len(src) {
					i++
				}
			}
		case '\'':
			i++
			for i < len(src) && src[i] != '\'' {
				if src[i] == '\\' && i+1 < len(src) {
					i += 2
					continue
				}
				i++
			}
			if i < len(src) {
				i++
			}
		default:
			i++
		}
	}
	return -1
}

// hasForbidden skips tests that rely on features we don't mirror:
// externs and source accumulation (multiple addSource calls).
func hasForbidden(body string) bool {
	if strings.Contains(body, "addExterns") {
		return true
	}
	if strings.Count(body, "addSource(") > 1 {
		return true
	}
	return false
}

// extractAddSource returns the JS source passed to the single
// addSource(...) call in the method body, decoding either an inline
// string literal or a Java text block.
func extractAddSource(body string) (string, bool) {
	i := strings.Index(body, "addSource(")
	if i < 0 {
		return "", false
	}
	i += len("addSource(")
	// skip whitespace
	for i < len(body) && (body[i] == ' ' || body[i] == '\n' || body[i] == '\t' || body[i] == '\r') {
		i++
	}
	// lines("line1", "line2", ...) is a test helper that joins with '\n'.
	if strings.HasPrefix(body[i:], "lines(") {
		return decodeLinesCall(body, i+len("lines("))
	}
	if i+2 < len(body) && body[i] == '"' && body[i+1] == '"' && body[i+2] == '"' {
		return decodeTextBlock(body, i)
	}
	if i < len(body) && body[i] == '"' {
		return decodeStringLiteral(body, i)
	}
	return "", false
}

// decodeLinesCall reads `"l1", "l2", ...)` after `lines(` and joins the
// strings with '\n' — mirroring the closure-compiler test helper of the
// same name.
func decodeLinesCall(body string, i int) (string, bool) {
	var parts []string
	for {
		for i < len(body) && (body[i] == ' ' || body[i] == '\n' || body[i] == '\t' || body[i] == '\r' || body[i] == ',') {
			i++
		}
		if i >= len(body) {
			return "", false
		}
		if body[i] == ')' {
			return strings.Join(parts, "\n"), true
		}
		if i+2 < len(body) && body[i] == '"' && body[i+1] == '"' && body[i+2] == '"' {
			s, ok := decodeTextBlock(body, i)
			if !ok {
				return "", false
			}
			parts = append(parts, s)
			// advance i past the closing """
			j := strings.Index(body[i+3:], `"""`)
			if j < 0 {
				return "", false
			}
			i = i + 3 + j + 3
			continue
		}
		if body[i] == '"' {
			// Reuse decodeStringLiteral but without the trailing `)` check.
			s, ok, next := decodeStringLiteralList(body, i)
			if !ok {
				return "", false
			}
			parts = append(parts, s)
			i = next
			continue
		}
		return "", false
	}
}

// decodeStringLiteralList reads ONE "..." at position i, returns the
// decoded content and the index just past the closing quote. Unlike
// decodeStringLiteral, it does not expect the next non-ws byte to be ')'.
func decodeStringLiteralList(body string, i int) (string, bool, int) {
	if body[i] != '"' {
		return "", false, i
	}
	i++
	var b strings.Builder
	for i < len(body) && body[i] != '"' {
		if body[i] == '\\' && i+1 < len(body) {
			b.WriteByte(unescape(body[i+1]))
			i += 2
			continue
		}
		b.WriteByte(body[i])
		i++
	}
	if i >= len(body) {
		return "", false, i
	}
	return b.String(), true, i + 1
}

func decodeStringLiteral(body string, i int) (string, bool) {
	// Must be followed by `)` after an optional chain of `+ "..."`.
	var b strings.Builder
	for {
		if i >= len(body) || body[i] != '"' {
			return "", false
		}
		i++ // opening quote
		for i < len(body) && body[i] != '"' {
			if body[i] == '\\' && i+1 < len(body) {
				b.WriteByte(unescape(body[i+1]))
				i += 2
				continue
			}
			b.WriteByte(body[i])
			i++
		}
		if i >= len(body) {
			return "", false
		}
		i++ // closing quote
		// Skip whitespace.
		for i < len(body) && (body[i] == ' ' || body[i] == '\n' || body[i] == '\t' || body[i] == '\r') {
			i++
		}
		if i >= len(body) {
			return "", false
		}
		if body[i] == ')' {
			return b.String(), true
		}
		if body[i] == '+' {
			i++
			for i < len(body) && (body[i] == ' ' || body[i] == '\n' || body[i] == '\t' || body[i] == '\r') {
				i++
			}
			continue
		}
		return "", false
	}
}

func unescape(b byte) byte {
	switch b {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case '"':
		return '"'
	case '\'':
		return '\''
	case '\\':
		return '\\'
	}
	return b
}

// decodeTextBlock mirrors the Java 17 text-block rules: drop the leading
// `"""` line, drop the trailing `"""`, and remove the minimum common
// leading-whitespace indent from each line of content.
func decodeTextBlock(body string, i int) (string, bool) {
	if !(i+2 < len(body) && body[i] == '"' && body[i+1] == '"' && body[i+2] == '"') {
		return "", false
	}
	i += 3
	// Expect newline.
	for i < len(body) && body[i] != '\n' {
		i++
	}
	if i >= len(body) {
		return "", false
	}
	i++ // past newline
	start := i
	for i+2 < len(body) && !(body[i] == '"' && body[i+1] == '"' && body[i+2] == '"') {
		i++
	}
	if i+2 >= len(body) {
		return "", false
	}
	raw := body[start:i]
	return stripCommonIndent(raw), true
}

func stripCommonIndent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		n := 0
		for n < len(l) && (l[n] == ' ' || l[n] == '\t') {
			n++
		}
		if minIndent < 0 || n < minIndent {
			minIndent = n
		}
	}
	if minIndent <= 0 {
		return s
	}
	var b strings.Builder
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if len(l) >= minIndent {
			b.WriteString(l[minIndent:])
		} else {
			b.WriteString(l)
		}
	}
	return b.String()
}

var nameClean = regexp.MustCompile(`[^A-Za-z0-9_]+`)

func sanitizeName(n string) string {
	n = nameClean.ReplaceAllString(n, "_")
	if len(n) > 120 {
		n = n[:120]
	}
	return strings.ToLower(n)
}
