// tools/bodyfidelity/main.go — per-method body-level fidelity checker.
//
// Walks all internal/**/*.go files. For each Go method whose enclosing file
// has a "// Port of <java-relative-path> @ closure-compiler v20260415" header,
// it extracts the Go body and the matching Java method body, computes a
// body-similarity score, and emits:
//   - docs/body-fidelity.md  — Markdown table sorted worst-first
//   - docs/body-fidelity.json — same data as JSON
//
// Score formula (0-100):
//   50 pts: LoC ratio between 0.66 and 1.50 → full; otherwise scaled
//   30 pts: control-flow keyword count within ±20% → full; otherwise scaled
//   20 pts: call count within ±25% → full; otherwise scaled
//
// Usage:
//
//	go run ./tools/bodyfidelity           # produce output files + console summary
//	go run ./tools/bodyfidelity -fix      # also print suggestions for top-N divergent methods

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

const (
	javaSourceRoot   = "/Users/damien/code/jsgo/.cache/closure-src"
	internalRoot     = "internal"
	expectedVersion  = "v20260415"
	mdOutputPath     = "docs/body-fidelity.md"
	jsonOutputPath   = "docs/body-fidelity.json"
	fixSuggestionTop = 20
)

// -------- regexps --------

// portOfFileRe matches a file-level "// Port of <path> @ closure-compiler <version>" header.
// Group 1 = java relative path, Group 2 = version.
var portOfFileRe = regexp.MustCompile(`^//\s*Port of\s+(\S+\.java[^\s@]*)\s*@\s*closure-compiler\s+(v\d+)`)

// portOfMethodRe matches a per-method "// Port of <path> @ closure-compiler <version>" comment.
// Same format as file-level but appearing right above a func declaration.
var portOfMethodRe = portOfFileRe

// javaKeywords maps Java keywords that should not be counted as method identifiers.
var javaKeywords = map[string]bool{
	"if": true, "for": true, "while": true, "switch": true, "catch": true,
	"return": true, "new": true, "throw": true, "else": true, "do": true,
	"try": true, "finally": true, "assert": true, "synchronized": true,
	"instanceof": true, "class": true, "interface": true, "enum": true,
	"static": true, "final": true, "public": true, "private": true, "protected": true,
	"void": true, "int": true, "boolean": true, "long": true, "double": true,
	"float": true, "char": true, "byte": true, "short": true,
}

// cfKeywordsGo and cfKeywordsJava are the control-flow keywords to count.
var cfKeywordsGo = []string{"if", "for", "switch", "return", "break", "continue"}

// cfKeywordsJava: same set for Java
var cfKeywordsJava = []string{"if", "for", "switch", "return", "break", "continue", "while"}

// javaMethodRe matches Java method declarations.
// Captures group 1 = method name.
var javaMethodRe = regexp.MustCompile(
	`(?m)^[ \t]*(?:(?:public|private|protected|static|final|synchronized|abstract|native|default|override|@\w+(?:\([^)]*\))?)\s+)*` +
		`(?:<[^>]+>\s*)?` +
		`(?:[\w$][\w<>\[\],\s]*?\s+)` +
		`(\w+)\s*\(`)

// annotationLineRe matches lines that are solely a Java annotation.
var annotationLineRe = regexp.MustCompile(`(?m)^[ \t]*@\w+(?:\([^)]*\))?[ \t]*$`)

// -------- data structures --------

// MethodResult holds computed metrics for one matched method pair.
type MethodResult struct {
	GoFile      string  `json:"go_file"`
	Method      string  `json:"method"`
	JavaFile    string  `json:"java_file"`
	JavaLine    int     `json:"java_line"`
	GoLoC       int     `json:"go_loc"`
	JavaLoC     int     `json:"java_loc"`
	GoCF        int     `json:"go_cf"`
	JavaCF      int     `json:"java_cf"`
	GoCalls     int     `json:"go_calls"`
	JavaCalls   int     `json:"java_calls"`
	Score       float64 `json:"score"`
	Notes       string  `json:"notes"`
}

// GoMethodInfo holds parsed info about a Go method.
type GoMethodInfo struct {
	Name     string
	Body     string // raw source text of the method body (between braces)
	Line     int    // line of the func keyword
	JavaPath string // file-level java relative path
}

// -------- Go parsing --------

// filePortPath extracts the java relative path from a file-level Port of header.
// Returns "" if not found or version doesn't match.
func filePortPath(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum > 10 {
			break
		}
		line := scanner.Text()
		m := portOfFileRe.FindStringSubmatch(line)
		if m != nil {
			if m[2] == expectedVersion {
				return m[1]
			}
		}
	}
	return ""
}

// goToJavaName converts a Go method name to the expected Java camelCase name.
// e.g. GetName → getName, IsSubtype → isSubtype, Visit → visit.
func goToJavaName(name string) string {
	if name == "" {
		return name
	}
	runes := []rune(name)
	// Lower-case the first run of consecutive uppercase letters, but keep
	// the rest intact (e.g. "GetURL" → "getURL", "IsOK" → "isOK").
	i := 0
	for i < len(runes) && unicode.IsUpper(runes[i]) {
		i++
	}
	if i == 0 {
		return name
	}
	// If whole string is uppercase, lower all; otherwise lower up to i-1.
	lower := i
	if i < len(runes) && i > 1 {
		lower = i - 1
	}
	result := make([]rune, len(runes))
	copy(result, runes)
	for j := 0; j < lower; j++ {
		result[j] = unicode.ToLower(runes[j])
	}
	return string(result)
}

// extractGoMethods parses filename and returns all exported+unexported method/function
// declarations, with their body source text.
func extractGoMethods(filename, javaPath string, src []byte, fset *token.FileSet) []GoMethodInfo {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil
	}

	// Build map: funcLine → port path, for per-method Port of headers.
	// A comment on line N applies to a func on line N+1 or N+2.
	type perMethodPort struct {
		javaPath string
		line     int
	}
	perMethod := map[int]perMethodPort{}
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			m := portOfMethodRe.FindStringSubmatch(c.Text)
			if m != nil && m[2] == expectedVersion {
				commentLine := fset.Position(c.Pos()).Line
				perMethod[commentLine] = perMethodPort{javaPath: m[1], line: commentLine}
			}
		}
	}

	var methods []GoMethodInfo
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		funcLine := fset.Position(fn.Pos()).Line

		// Determine which java path applies to this method.
		// Per-method header wins; fall back to file-level.
		jp := javaPath
		for cl, pm := range perMethod {
			if funcLine-cl >= 1 && funcLine-cl <= 3 {
				jp = pm.javaPath
				break
			}
		}
		if jp == "" {
			continue
		}

		// Extract body source.
		bodyStart := fset.Position(fn.Body.Lbrace).Offset
		bodyEnd := fset.Position(fn.Body.Rbrace).Offset + 1
		if bodyEnd > len(src) {
			bodyEnd = len(src)
		}
		body := string(src[bodyStart:bodyEnd])

		methods = append(methods, GoMethodInfo{
			Name:     fn.Name.Name,
			Body:     body,
			Line:     funcLine,
			JavaPath: jp,
		})
	}
	return methods
}

// -------- body metric helpers --------

// countLoC counts non-blank, non-comment lines in Go body text.
func countGoLoC(body string) int {
	count := 0
	inBlock := false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if inBlock {
			if strings.Contains(trimmed, "*/") {
				inBlock = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "/*") {
			inBlock = true
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || trimmed == "{" || trimmed == "}" {
			continue
		}
		count++
	}
	return count
}

// countJavaLoC counts non-blank, non-comment lines in Java body text.
func countJavaLoC(body string) int {
	count := 0
	inBlock := false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if inBlock {
			if strings.Contains(trimmed, "*/") {
				inBlock = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "/*") {
			inBlock = true
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") || trimmed == "{" || trimmed == "}" {
			continue
		}
		count++
	}
	return count
}

// countCFKeywords counts occurrences of control-flow keywords in source.
// Uses word-boundary matching to avoid false positives.
var cfWordRe = regexp.MustCompile(`\b(if|for|switch|return|break|continue|while)\b`)

func countCFKeywords(body string, keywords []string) int {
	// Strip strings and comments first to avoid false hits inside literals.
	stripped := stripLiterals(body)
	count := 0
	for _, kw := range keywords {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(kw) + `\b`)
		count += len(re.FindAllString(stripped, -1))
	}
	return count
}

// countCalls counts rough function call occurrences: `identifier(` where identifier
// is not a keyword or starts with uppercase (type assertion or conversion allowed).
var callRe = regexp.MustCompile(`\b([a-zA-Z_]\w*)\s*\(`)

var goKeywordSet = map[string]bool{
	"if": true, "for": true, "switch": true, "return": true, "break": true,
	"continue": true, "select": true, "go": true, "defer": true, "func": true,
	"make": true, "new": true, "append": true, "copy": true, "delete": true,
	"len": true, "cap": true, "close": true, "panic": true, "recover": true,
	"print": true, "println": true, "range": true, "type": true, "map": true,
	"chan": true, "var": true, "const": true, "package": true, "import": true,
}

var javaKeywordSet = map[string]bool{
	"if": true, "for": true, "while": true, "switch": true, "return": true,
	"break": true, "continue": true, "new": true, "throw": true, "catch": true,
	"try": true, "finally": true, "assert": true, "synchronized": true,
	"instanceof": true, "class": true, "interface": true, "enum": true,
	"super": true, "this": true,
}

func countGoCallsInBody(body string) int {
	stripped := stripLiterals(body)
	matches := callRe.FindAllStringSubmatch(stripped, -1)
	count := 0
	for _, m := range matches {
		name := m[1]
		if !goKeywordSet[name] {
			count++
		}
	}
	return count
}

func countJavaCallsInBody(body string) int {
	stripped := stripLiterals(body)
	matches := callRe.FindAllStringSubmatch(stripped, -1)
	count := 0
	for _, m := range matches {
		name := m[1]
		if !javaKeywordSet[name] {
			count++
		}
	}
	return count
}

// stripLiterals replaces string literals and comments with spaces to avoid
// counting keywords or parens inside them.
func stripLiterals(src string) string {
	var b strings.Builder
	b.Grow(len(src))
	i := 0
	for i < len(src) {
		// Single-line comment
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '/' {
			for i < len(src) && src[i] != '\n' {
				b.WriteByte(' ')
				i++
			}
			continue
		}
		// Block comment
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '*' {
			b.WriteString("  ")
			i += 2
			for i < len(src) {
				if i+1 < len(src) && src[i] == '*' && src[i+1] == '/' {
					b.WriteString("  ")
					i += 2
					break
				}
				b.WriteByte(' ')
				i++
			}
			continue
		}
		// Double-quoted string
		if src[i] == '"' {
			b.WriteByte(' ')
			i++
			for i < len(src) && src[i] != '"' {
				if src[i] == '\\' {
					b.WriteString("  ")
					i += 2
					continue
				}
				b.WriteByte(' ')
				i++
			}
			if i < len(src) {
				b.WriteByte(' ')
				i++
			}
			continue
		}
		// Single-quoted char (Java)
		if src[i] == '\'' {
			b.WriteByte(' ')
			i++
			for i < len(src) && src[i] != '\'' {
				if src[i] == '\\' {
					b.WriteString("  ")
					i += 2
					continue
				}
				b.WriteByte(' ')
				i++
			}
			if i < len(src) {
				b.WriteByte(' ')
				i++
			}
			continue
		}
		// Backtick string (Go)
		if src[i] == '`' {
			b.WriteByte(' ')
			i++
			for i < len(src) && src[i] != '`' {
				b.WriteByte(' ')
				i++
			}
			if i < len(src) {
				b.WriteByte(' ')
				i++
			}
			continue
		}
		b.WriteByte(src[i])
		i++
	}
	return b.String()
}

// -------- Java parsing --------

// javaSourcePath converts a java relative path from the Port of header to an absolute path.
// It handles paths that start with "com/google/..." directly or just "jscomp/..." etc.
func javaSourcePath(javaRelPath string) string {
	// Clean up trailing suffixes like "(narrowed)" etc.
	javaRelPath = strings.Fields(javaRelPath)[0]
	javaRelPath = strings.TrimSuffix(javaRelPath, ")")

	// Try direct join
	candidate := filepath.Join(javaSourceRoot, javaRelPath)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	// Try prepending com/google/javascript/
	candidate2 := filepath.Join(javaSourceRoot, "com/google/javascript", javaRelPath)
	if _, err := os.Stat(candidate2); err == nil {
		return candidate2
	}
	return ""
}

// extractJavaMethodBody finds the first Java method named methodName in src,
// returns (body, lineNumber) where body is the text between the matching braces.
// If multiple overloads exist, returns the one with most lines (heuristic).
func extractJavaMethodBody(src, methodName string) (string, int) {
	stripped := annotationLineRe.ReplaceAllString(src, "")
	lines := strings.Split(src, "\n")

	// Find all positions where the method appears.
	pattern := regexp.MustCompile(`(?m)^[ \t]*(?:(?:public|private|protected|static|final|synchronized|abstract|native|default|override)\s+)*` +
		`(?:<[^>]+>\s*)?` +
		`(?:[\w$][\w<>\[\],\s]*?\s+)` +
		`\b` + regexp.QuoteMeta(methodName) + `\s*\(`)

	_ = stripped // used pattern on stripped above; use src for offset tracking
	allMatches := pattern.FindAllStringIndex(src, -1)
	if len(allMatches) == 0 {
		return "", 0
	}

	type candidate struct {
		body    string
		lineNum int
		bodyLen int
	}
	var candidates []candidate

	for _, match := range allMatches {
		matchStart := match[0]
		// Find the opening brace of this method.
		braceStart := strings.Index(src[matchStart:], "{")
		if braceStart < 0 {
			continue
		}
		braceStart += matchStart

		// Check it's not an abstract method (no body).
		betweenSig := src[matchStart:braceStart]
		if strings.Contains(betweenSig, ";") {
			continue
		}

		// Extract balanced braces.
		body, endPos := extractBalancedBraces(src, braceStart)
		if endPos < 0 {
			continue
		}

		// Compute line number.
		lineNum := 1 + strings.Count(src[:matchStart], "\n")
		// Verify lineNum is within range.
		if lineNum > len(lines) {
			lineNum = len(lines)
		}

		candidates = append(candidates, candidate{body: body, lineNum: lineNum, bodyLen: len(body)})
	}

	if len(candidates) == 0 {
		return "", 0
	}
	// Return longest body (covers most overloads heuristic).
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.bodyLen > best.bodyLen {
			best = c
		}
	}
	return best.body, best.lineNum
}

// extractBalancedBraces extracts the text from src[start] (which must be '{')
// to the matching '}', returning (text, endIndex+1) or ("", -1) on failure.
func extractBalancedBraces(src string, start int) (string, int) {
	if start >= len(src) || src[start] != '{' {
		return "", -1
	}
	depth := 0
	i := start
	for i < len(src) {
		ch := src[i]
		if ch == '"' || ch == '\'' || ch == '`' {
			// Skip string literal.
			quote := ch
			i++
			for i < len(src) {
				if src[i] == '\\' {
					i += 2
					continue
				}
				if src[i] == quote {
					i++
					break
				}
				i++
			}
			continue
		}
		if i+1 < len(src) && ch == '/' && src[i+1] == '/' {
			for i < len(src) && src[i] != '\n' {
				i++
			}
			continue
		}
		if i+1 < len(src) && ch == '/' && src[i+1] == '*' {
			i += 2
			for i < len(src) {
				if i+1 < len(src) && src[i] == '*' && src[i+1] == '/' {
					i += 2
					break
				}
				i++
			}
			continue
		}
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				return src[start : i+1], i + 1
			}
		}
		i++
	}
	return "", -1
}

// -------- scoring --------

// computeScore returns a body-similarity score 0-100.
func computeScore(goLoC, javaLoC, goCF, javaCF, goCalls, javaCalls int) float64 {
	// LoC component (50 pts).
	locScore := 0.0
	if javaLoC == 0 && goLoC == 0 {
		locScore = 50
	} else if javaLoC == 0 {
		locScore = 0
	} else {
		ratio := float64(goLoC) / float64(javaLoC)
		if ratio >= 0.66 && ratio <= 1.50 {
			locScore = 50
		} else {
			// Scale: clamp ratio to [0.01, 10], then 50 * (1 - |log(ratio)|/log(thresholdDistance))
			// Use simpler: distance from nearest threshold.
			dist := 0.0
			if ratio < 0.66 {
				dist = 0.66 - ratio
			} else {
				dist = ratio - 1.50
			}
			// Max possible distance in a useful range: ~3 (ratio=4.5 → dist=3).
			locScore = math.Max(0, 50*(1-dist/2.0))
		}
	}

	// CF component (30 pts).
	cfScore := 0.0
	if javaCF == 0 && goCF == 0 {
		cfScore = 30
	} else if javaCF == 0 {
		cfScore = 0
	} else {
		ratio := float64(goCF) / float64(javaCF)
		if ratio >= 0.80 && ratio <= 1.20 {
			cfScore = 30
		} else {
			dist := 0.0
			if ratio < 0.80 {
				dist = 0.80 - ratio
			} else {
				dist = ratio - 1.20
			}
			cfScore = math.Max(0, 30*(1-dist/1.5))
		}
	}

	// Call count component (20 pts).
	callScore := 0.0
	if javaCalls == 0 && goCalls == 0 {
		callScore = 20
	} else if javaCalls == 0 {
		callScore = 0
	} else {
		ratio := float64(goCalls) / float64(javaCalls)
		if ratio >= 0.75 && ratio <= 1.25 {
			callScore = 20
		} else {
			dist := 0.0
			if ratio < 0.75 {
				dist = 0.75 - ratio
			} else {
				dist = ratio - 1.25
			}
			callScore = math.Max(0, 20*(1-dist/1.5))
		}
	}

	total := locScore + cfScore + callScore
	if total > 100 {
		total = 100
	}
	return math.Round(total)
}

// -------- main logic --------

func run(fixMode bool) error {
	// Find the repo root (where internal/ lives).
	repoRoot, err := findRepoRoot()
	if err != nil {
		return fmt.Errorf("could not find repo root: %w", err)
	}

	internalDir := filepath.Join(repoRoot, internalRoot)
	var results []MethodResult
	javaCache := map[string]string{} // java abs path → source

	err = filepath.Walk(internalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Check file-level port header.
		javaPath := filePortPath(path)
		if javaPath == "" {
			return nil
		}

		// Read Go source.
		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		fset := token.NewFileSet()
		methods := extractGoMethods(path, javaPath, src, fset)
		if len(methods) == 0 {
			return nil
		}

		// Load Java source (cached).
		for _, m := range methods {
			javaAbsPath, ok := javaCache["path:"+m.JavaPath]
			if !ok {
				javaAbsPath = javaSourcePath(m.JavaPath)
				javaCache["path:"+m.JavaPath] = javaAbsPath
			}
			if javaAbsPath == "" {
				continue
			}

			javaSrc, ok := javaCache["src:"+javaAbsPath]
			if !ok {
				data, err := os.ReadFile(javaAbsPath)
				if err != nil {
					javaCache["src:"+javaAbsPath] = ""
					continue
				}
				javaSrc = string(data)
				javaCache["src:"+javaAbsPath] = javaSrc
			}
			if javaSrc == "" {
				continue
			}

			// Try Go name and Java camelCase name.
			javaName := goToJavaName(m.Name)
			javaBody, javaLine := extractJavaMethodBody(javaSrc, javaName)
			if javaBody == "" && javaName != m.Name {
				javaBody, javaLine = extractJavaMethodBody(javaSrc, m.Name)
			}
			if javaBody == "" {
				// Try lower-case first letter only.
				alt := strings.ToLower(m.Name[:1]) + m.Name[1:]
				javaBody, javaLine = extractJavaMethodBody(javaSrc, alt)
			}
			if javaBody == "" {
				continue
			}

			goLoC := countGoLoC(m.Body)
			javaLoC := countJavaLoC(javaBody)
			goCF := countCFKeywords(m.Body, cfKeywordsGo)
			javaCF := countCFKeywords(javaBody, cfKeywordsJava)
			goCalls := countGoCallsInBody(m.Body)
			javaCalls := countJavaCallsInBody(javaBody)

			score := computeScore(goLoC, javaLoC, goCF, javaCF, goCalls, javaCalls)

			// Build notes.
			var notes []string
			if javaLoC > 0 && goLoC > 0 {
				ratio := float64(goLoC) / float64(javaLoC)
				if ratio < 0.3 {
					notes = append(notes, fmt.Sprintf("likely stub (%d LoC vs Java's %d)", goLoC, javaLoC))
				} else if ratio > 3.0 {
					notes = append(notes, fmt.Sprintf("much larger (%d LoC vs Java's %d)", goLoC, javaLoC))
				}
			}

			// Relative Go file path for display.
			relGoFile, _ := filepath.Rel(repoRoot, path)
			relJavaFile, _ := filepath.Rel(javaSourceRoot, javaAbsPath)

			results = append(results, MethodResult{
				GoFile:   relGoFile,
				Method:   m.Name,
				JavaFile: relJavaFile,
				JavaLine: javaLine,
				GoLoC:    goLoC,
				JavaLoC:  javaLoC,
				GoCF:     goCF,
				JavaCF:   javaCF,
				GoCalls:  goCalls,
				JavaCalls: javaCalls,
				Score:    score,
				Notes:    strings.Join(notes, "; "),
			})
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Sort worst-first.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	// Compute summary stats.
	total := len(results)
	var sumScore float64
	below50 := 0
	for _, r := range results {
		sumScore += r.Score
		if r.Score < 50 {
			below50++
		}
	}
	avgScore := 0.0
	if total > 0 {
		avgScore = sumScore / float64(total)
	}

	// Write JSON.
	docsDir := filepath.Join(repoRoot, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(repoRoot, jsonOutputPath), jsonData, 0644); err != nil {
		return err
	}

	// Write Markdown.
	mdPath := filepath.Join(repoRoot, mdOutputPath)
	if err := writeMD(mdPath, results, total, avgScore, below50); err != nil {
		return err
	}

	// Console summary.
	fmt.Printf("Body fidelity check complete.\n")
	fmt.Printf("  Total matched methods : %d\n", total)
	fmt.Printf("  Average score         : %.1f\n", avgScore)
	fmt.Printf("  Methods below 50      : %d\n", below50)
	fmt.Printf("  Output: %s\n", mdOutputPath)
	fmt.Printf("  Output: %s\n", jsonOutputPath)

	// -fix suggestions.
	if fixMode {
		fmt.Printf("\n--- Top-%d divergent method suggestions ---\n", fixSuggestionTop)
		n := fixSuggestionTop
		if n > len(results) {
			n = len(results)
		}
		for i := 0; i < n; i++ {
			r := results[i]
			fmt.Printf("  [score=%.0f] method %s in %s — Go %d LoC vs Java %d LoC (CF Go=%d Java=%d, calls Go=%d Java=%d) — %s\n",
				r.Score, r.Method, r.GoFile, r.GoLoC, r.JavaLoC, r.GoCF, r.JavaCF, r.GoCalls, r.JavaCalls,
				suggestion(r))
		}
	}

	return nil
}

func suggestion(r MethodResult) string {
	if r.JavaLoC > 0 && r.GoLoC > 0 {
		ratio := float64(r.GoLoC) / float64(r.JavaLoC)
		if ratio < 0.3 {
			return fmt.Sprintf("likely a stub; expand body to ~%d lines", r.JavaLoC)
		}
		if ratio > 3.0 {
			return fmt.Sprintf("Go body is much larger; review for over-expansion")
		}
	}
	if r.GoLoC == 0 {
		return "empty Go body; likely a stub"
	}
	if r.GoCF == 0 && r.JavaCF > 3 {
		return fmt.Sprintf("missing control flow; Java has %d CF keywords", r.JavaCF)
	}
	if r.GoCalls < r.JavaCalls/2 {
		return fmt.Sprintf("fewer calls (%d) than Java (%d); may be missing logic", r.GoCalls, r.JavaCalls)
	}
	return "review body divergence"
}

// findRepoRoot walks up from the current directory to find the repo root
// (the directory containing go.mod).
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found")
}

// -------- Markdown output --------

const mdTemplate = `# Body Fidelity Report

Generated: auto-computed by tools/bodyfidelity.

## Summary

| Metric | Value |
|--------|-------|
| Total matched methods | {{.Total}} |
| Average score | {{printf "%.1f" .AvgScore}} |
| Methods below 50 | {{.Below50}} |

## Methods by Body-Similarity Score (worst first)

| Rank | Score | Go File | Method | Java File | Java Line | Go LoC | Java LoC | Notes |
|------|-------|---------|--------|-----------|-----------|--------|----------|-------|
{{range $i, $r := .Results}}| {{inc $i}} | {{printf "%.0f" $r.Score}} | ` + "`{{$r.GoFile}}`" + ` | ` + "`{{$r.Method}}`" + ` | ` + "`{{$r.JavaFile}}`" + ` | {{$r.JavaLine}} | {{$r.GoLoC}} | {{$r.JavaLoC}} | {{$r.Notes}} |
{{end}}`

func writeMD(path string, results []MethodResult, total int, avgScore float64, below50 int) error {
	funcMap := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}
	tmpl, err := template.New("md").Funcs(funcMap).Parse(mdTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, struct {
		Total   int
		AvgScore float64
		Below50 int
		Results []MethodResult
	}{total, avgScore, below50, results})
}

func main() {
	fixFlag := flag.Bool("fix", false, "print suggestions for top divergent methods")
	flag.Parse()

	if err := run(*fixFlag); err != nil {
		fmt.Fprintf(os.Stderr, "bodyfidelity: %v\n", err)
		os.Exit(1)
	}
}
