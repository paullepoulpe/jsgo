// javaport is a skeleton generator that turns a closure-compiler
// Java source file into a Go stub ready for manual porting. It
// emits:
//
//   - the // Port of <java-path> @ v20260415 (skeleton) header
//     required by ci-guards.sh
//   - the target package declaration
//   - a struct for each top-level class and inner class (prefixed
//     with the outer class name, e.g., OuterInner)
//   - a method stub for every public / package / protected method,
//     with the original javadoc preserved above it
//   - a panic("unimplemented: …") in each untranslated stub body so
//     missing ports surface loudly at runtime
//   - // @ <java-path>:<line> pin comments on each method
//   - Go generics for Java generic type parameters
//   - Guava/JDK collection types mapped to Go equivalents
//   - AutoValue @AutoValue classes → struct + factory
//   - enum-with-methods → interface + per-constant struct
//   - Java import → Go import block via package mapping table
//   - per-class translation cache in .javaport-cache/<hash>.json
//   - Javadoc /** */ → Go // doc comments
//
// Usage:
//
//	go run ./tools/javaport \
//	    -in  .cache/closure-src/com/google/javascript/jscomp/CheckJSDoc.java \
//	    -out internal/jscomp/checkjsdoc.go \
//	    -pkg jscomp

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/paullepoulpe/jsgo/tools/javaport/javaast"
)

// javaportVersion is bumped whenever the generator logic changes to
// invalidate stale cache entries.
const javaportVersion = "2"

func main() {
	inPath := flag.String("in", "", "java source file (required)")
	outPath := flag.String("out", "", "go output file (required)")
	pkg := flag.String("pkg", "jscomp", "go package name")
	buildTag := flag.String("build-tag", "", "add a //go:build <tag> constraint so the file is excluded from the default build")
	cacheDir := flag.String("cache-dir", ".javaport-cache", "directory for translation cache")
	gofmtCheck := flag.Bool("gofmt-check", false, "after writing, run gofmt and report any remaining errors")
	flag.Parse()
	if *inPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: javaport -in <foo.java> -out <foo.go> [-pkg <name>] [-build-tag <tag>] [-gofmt-check]")
		os.Exit(2)
	}
	src, err := os.ReadFile(*inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	javaPath, err := relToSrcRoot(*inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// #171 — check translation cache before doing work.
	cacheKey := computeCacheKey(src)
	if cached := loadCache(*cacheDir, cacheKey); cached != "" {
		if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err := os.WriteFile(*outPath, []byte(cached), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Printf("wrote %s (%d bytes) [cache hit]\n", *outPath, len(cached))
		return
	}

	out := generateSkeleton(string(src), javaPath, *pkg, *buildTag)
	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := os.WriteFile(*outPath, []byte(out), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	saveCache(*cacheDir, cacheKey, out)
	fmt.Printf("wrote %s (%d bytes)\n", *outPath, len(out))
	if *gofmtCheck {
		reportGofmtErrors(*outPath, []byte(out))
	}
}

// reportGofmtErrors runs go/format on src and prints any lines that differ
// from the formatted version, prefixed with their line number.  This is
// equivalent to `gofmt -d` but uses the in-process formatter so no external
// binary is required.  Errors are printed to stderr; the process does not
// exit non-zero so that -gofmt-check is advisory only.
func reportGofmtErrors(path string, src []byte) {
	formatted, err := format.Source(src)
	if err != nil {
		// Syntax errors: just report the parse failure.
		fmt.Fprintf(os.Stderr, "gofmt-check %s: parse error: %v\n", path, err)
		return
	}
	if bytes.Equal(src, formatted) {
		fmt.Printf("gofmt-check %s: OK (0 errors)\n", path)
		return
	}
	// Count differing lines.
	srcLines := bytes.Split(src, []byte("\n"))
	fmtLines := bytes.Split(formatted, []byte("\n"))
	n := len(srcLines)
	if len(fmtLines) > n {
		n = len(fmtLines)
	}
	errCount := 0
	for i := 0; i < n; i++ {
		var sl, fl []byte
		if i < len(srcLines) {
			sl = srcLines[i]
		}
		if i < len(fmtLines) {
			fl = fmtLines[i]
		}
		if !bytes.Equal(sl, fl) {
			errCount++
			fmt.Fprintf(os.Stderr, "gofmt-check %s:%d: want %q got %q\n",
				path, i+1, string(fl), string(sl))
		}
	}
	fmt.Fprintf(os.Stderr, "gofmt-check %s: %d error(s)\n", path, errCount)
}

// ── #171 translation cache ──────────────────────────────────────────────────

type cacheEntry struct {
	Version string `json:"version"`
	Output  string `json:"output"`
}

func computeCacheKey(src []byte) string {
	h := sha256.Sum256(src)
	return fmt.Sprintf("%x", h)
}

func cacheFilePath(dir, key string) string {
	return filepath.Join(dir, key[:2], key+".json")
}

func loadCache(dir, key string) string {
	data, err := os.ReadFile(cacheFilePath(dir, key))
	if err != nil {
		return ""
	}
	var e cacheEntry
	if err := json.Unmarshal(data, &e); err != nil {
		return ""
	}
	if e.Version != javaportVersion {
		return ""
	}
	return e.Output
}

func saveCache(dir, key, output string) {
	path := cacheFilePath(dir, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, _ := json.Marshal(cacheEntry{Version: javaportVersion, Output: output})
	_ = os.WriteFile(path, data, 0o644)
}

// relToSrcRoot returns the Java file's package-relative path — the
// `com/google/javascript/…/FooBar.java` fragment. We strip any
// leading `.cache/closure-src/` so the output's `// Port of` header
// matches the convention ci-guards.sh enforces.
func relToSrcRoot(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	abs = filepath.ToSlash(abs)
	i := strings.Index(abs, "com/google/javascript/")
	if i < 0 {
		return "", fmt.Errorf("path %q missing com/google/javascript/ root", p)
	}
	return abs[i:], nil
}

// generateSkeleton drives the pipeline: strip comments, detect the
// class declaration, emit the Port-of header, translate fields +
// method signatures, write panic("unimplemented") bodies.
func generateSkeleton(src, javaPath, pkg, buildTag string) string {
	var b strings.Builder
	if buildTag != "" {
		b.WriteString("//go:build " + buildTag + "\n\n")
	}
	b.WriteString(headerFor(javaPath))
	b.WriteString("\npackage " + pkg + "\n\n")
	b.WriteString(importBlockFor(src, javaPath))
	cls := parseClass(src)
	if cls == nil {
		b.WriteString("// TODO: manual — no top-level class found in source\n")
		return b.String()
	}
	b.WriteString(emitDiagnosticConsts(cls, src))
	b.WriteString(emitEnum(cls, javaPath))
	b.WriteString(emitAutoValue(cls, javaPath))
	b.WriteString(emitStruct(cls))
	b.WriteString(emitMethods(cls, javaPath))
	// #162 — emit inner classes
	for _, inner := range cls.innerClasses {
		b.WriteString(emitInnerClass(cls.name, inner, javaPath))
	}
	return b.String()
}

// headerFor produces the Port-of traceability header ci-guards.sh
// requires. (skeleton) qualifier flags this file as untranslated
// body logic.
func headerFor(javaPath string) string {
	return "// Port of " + javaPath + "\n" +
		"// @ closure-compiler v20260415 (skeleton; bodies TODO).\n" +
		"//\n" +
		"// Generated by tools/javaport from the upstream Java source.\n" +
		"// Method signatures and field types are mapped mechanically;\n" +
		"// logic, control flow, and Java-idiom conversions (visitors,\n" +
		"// builders, enums-with-methods) are left to a follow-up pass.\n" +
		"// Hand-edit freely — re-running the generator will overwrite.\n"
}

// ── Data structures ──────────────────────────────────────────────────────────

// classDecl captures the parsed shape of a Java top-level class.
type classDecl struct {
	name         string
	javadoc      string
	modifier     string // "final" / "abstract" / ""
	extends      string
	implList     []string
	fields       []fieldDecl
	methods      []methodDecl
	body         string
	isAutoValue  bool
	isEnum       bool
	enumConsts   []enumConst
	innerClasses []*classDecl
	// #170 — source line of class declaration
	sourceLine int
	// #166 — generic type parameters
	typeParams []typeParam
}

// #163 — enum constant (may carry method bodies)
type enumConst struct {
	name    string
	javadoc string
	body    string // non-empty if constant has anonymous-class body
}

// #166 — generic type parameter
type typeParam struct {
	name  string
	bound string // "" means "any"
}

type fieldDecl struct {
	javadoc     string
	modifiers   string // "private final", "public static", …
	typeName    string
	name        string
	initialiser string // raw rhs text, empty when absent
}

type methodDecl struct {
	javadoc    string
	modifiers  string
	returnType string
	name       string
	params     []paramDecl
	throws     string
	body       string // raw body text (including braces)
	sourceLine int    // #170
}

type paramDecl struct {
	typeName string
	name     string
}

// ── Java source parsing ──────────────────────────────────────────────────────

// classHeadRe matches the top-level class declaration line.
var classHeadRe = regexp.MustCompile(
	`(?m)^\s*(?:public\s+|private\s+|protected\s+)?(?:(abstract|final)\s+)?class\s+(\w+)` +
		`(?:\s*<([^{>]+)>)?` +
		`(?:\s+extends\s+([\w.<>,\s]+?))?` +
		`(?:\s+implements\s+([\w.<>,\s]+?))?` +
		`\s*\{`)

// enumHeadRe matches Java enum declarations.
var enumHeadRe = regexp.MustCompile(
	`(?m)^\s*(?:public\s+|private\s+|protected\s+)?enum\s+(\w+)` +
		`(?:\s+implements\s+([\w.<>,\s]+?))?` +
		`\s*\{`)

func parseClass(src string) *classDecl {
	// Check for top-level enum. We only treat the file as an enum-first source
	// when the enum declaration appears before any class declaration; otherwise
	// the enum is an inner type inside a class (e.g. NodeUtil.ValueType).
	em := enumHeadRe.FindStringSubmatchIndex(src)
	cm := classHeadRe.FindStringSubmatchIndex(src)
	if em != nil && (cm == nil || em[0] < cm[0]) {
		cls := &classDecl{
			name:   src[em[2]:em[3]],
			isEnum: true,
		}
		if em[4] != -1 {
			for _, p := range strings.Split(src[em[4]:em[5]], ",") {
				cls.implList = append(cls.implList, strings.TrimSpace(p))
			}
		}
		cls.javadoc = javadocBefore(src, em[0])
		cls.sourceLine = lineNumber(src, em[0])
		body, ok := braceBody(src, em[1]-1)
		if ok {
			cls.body = body
			cls.enumConsts = parseEnumConsts(body)
			_, cls.methods = parseBody(body)
		}
		return cls
	}

	m := cm
	if m == nil {
		return nil
	}
	cls := &classDecl{
		name: src[m[4]:m[5]],
	}
	if m[2] != -1 {
		cls.modifier = strings.TrimSpace(src[m[2]:m[3]])
	}
	// #166 — type parameters
	if m[6] != -1 {
		cls.typeParams = parseTypeParams(src[m[6]:m[7]])
	}
	if m[8] != -1 {
		cls.extends = strings.TrimSpace(src[m[8]:m[9]])
	}
	if m[10] != -1 {
		for _, p := range strings.Split(src[m[10]:m[11]], ",") {
			cls.implList = append(cls.implList, strings.TrimSpace(p))
		}
	}
	cls.javadoc = javadocBefore(src, m[0])
	cls.sourceLine = lineNumber(src, m[0])
	// #168 — detect @AutoValue
	cls.isAutoValue = strings.Contains(javadocBefore(src, m[0]), "@AutoValue") ||
		strings.Contains(src[:m[0]], "@AutoValue")

	body, ok := braceBody(src, m[1]-1)
	if !ok {
		return cls
	}
	cls.body = body
	cls.fields, cls.methods = parseBody(body)
	// #162 — parse inner classes
	cls.innerClasses = parseInnerClasses(body)
	return cls
}

// parseTypeParams parses `T extends Foo, K, V` into []typeParam.
// #166
func parseTypeParams(s string) []typeParam {
	var out []typeParam
	parts := splitTopLevel(s, ',')
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Handle wildcard ?
		if p == "?" {
			out = append(out, typeParam{name: "T", bound: "any"})
			continue
		}
		if strings.HasPrefix(p, "? extends ") {
			bound := javaToGoType(strings.TrimPrefix(p, "? extends "))
			out = append(out, typeParam{name: "T", bound: bound})
			continue
		}
		if strings.HasPrefix(p, "? super ") {
			// Go has no contravariance — drop bound
			out = append(out, typeParam{name: "T", bound: "any"})
			continue
		}
		if idx := strings.Index(p, " extends "); idx >= 0 {
			name := strings.TrimSpace(p[:idx])
			bound := javaToGoType(strings.TrimSpace(p[idx+len(" extends "):]))
			out = append(out, typeParam{name: name, bound: bound})
			continue
		}
		out = append(out, typeParam{name: p, bound: "any"})
	}
	return out
}

// typeParamsGoString formats []typeParam as `[T any, K comparable, V any]`.
// Returns "" when params is empty.
func typeParamsGoString(params []typeParam) string {
	if len(params) == 0 {
		return ""
	}
	var parts []string
	for _, tp := range params {
		if tp.bound == "" || tp.bound == "any" {
			parts = append(parts, tp.name+" any")
		} else {
			parts = append(parts, tp.name+" "+tp.bound)
		}
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// splitTopLevel splits s at every sep that is at nesting depth 0.
func splitTopLevel(s string, sep rune) []string {
	depth := 0
	var out []string
	start := 0
	for i, r := range s {
		switch r {
		case '<', '(', '[':
			depth++
		case '>', ')', ']':
			depth--
		case sep:
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}

// lineNumber returns the 1-based line number of offset pos in src.
func lineNumber(src string, pos int) int {
	return strings.Count(src[:pos], "\n") + 1
}

// #163 — parse enum constants from the body.
// Handles both simple `A, B, C` and method-carrying `A { ... }`.
func parseEnumConsts(body string) []enumConst {
	// Enum body: constants first, then optionally `;` separating methods.
	semiIdx := -1
	depth := 0
	for i := 0; i < len(body); i++ {
		switch body[i] {
		case '{':
			depth++
		case '}':
			depth--
		case ';':
			if depth == 0 {
				semiIdx = i
				break
			}
		}
		if semiIdx >= 0 {
			break
		}
	}
	constPart := body
	if semiIdx >= 0 {
		constPart = body[:semiIdx]
	}

	var consts []enumConst
	i := 0
	for i < len(constPart) {
		i = skipWS(constPart, i)
		if i >= len(constPart) {
			break
		}
		// Skip comments
		if strings.HasPrefix(constPart[i:], "//") {
			nl := strings.Index(constPart[i:], "\n")
			if nl < 0 {
				break
			}
			i += nl + 1
			continue
		}
		if strings.HasPrefix(constPart[i:], "/*") {
			end := strings.Index(constPart[i:], "*/")
			if end < 0 {
				break
			}
			i += end + 2
			continue
		}
		// Read constant name
		doc := javadocBefore(constPart, i)
		j := i
		for j < len(constPart) && isJavaIdentChar(constPart[j]) {
			j++
		}
		if j == i {
			break
		}
		name := constPart[i:j]
		i = j
		i = skipWS(constPart, i)
		// Skip constructor args if present
		if i < len(constPart) && constPart[i] == '(' {
			_, end2 := findMatchingParen(constPart, i)
			if end2 >= 0 {
				i = end2 + 1
			}
		}
		i = skipWS(constPart, i)
		// Check for anonymous class body
		var constBody string
		if i < len(constPart) && constPart[i] == '{' {
			var ok bool
			constBody, ok = braceBody(constPart, i)
			if ok {
				i += 2 + len(constBody)
			}
		}
		consts = append(consts, enumConst{name: name, javadoc: doc, body: constBody})
		i = skipWS(constPart, i)
		if i < len(constPart) && constPart[i] == ',' {
			i++
		}
	}
	return consts
}

func findMatchingParen(s string, open int) (int, int) {
	if open >= len(s) || s[open] != '(' {
		return -1, -1
	}
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return open, i
			}
		}
	}
	return -1, -1
}

// #162 — parse inner (nested) class declarations from a class body.
var innerClassRe = regexp.MustCompile(
	`(?m)(?:public\s+|private\s+|protected\s+)?(?:static\s+)?(?:(abstract|final)\s+)?class\s+(\w+)` +
		`(?:\s*<([^{>]+)>)?` +
		`(?:\s+extends\s+([\w.<>,\s]+?))?` +
		`(?:\s+implements\s+([\w.<>,\s]+?))?` +
		`\s*\{`)

func parseInnerClasses(body string) []*classDecl {
	var result []*classDecl
	matches := innerClassRe.FindAllStringSubmatchIndex(body, -1)
	for _, m := range matches {
		inner := &classDecl{
			name: body[m[4]:m[5]],
		}
		if m[2] != -1 {
			inner.modifier = strings.TrimSpace(body[m[2]:m[3]])
		}
		if m[6] != -1 {
			inner.typeParams = parseTypeParams(body[m[6]:m[7]])
		}
		if m[8] != -1 {
			inner.extends = strings.TrimSpace(body[m[8]:m[9]])
		}
		if m[10] != -1 {
			for _, p := range strings.Split(body[m[10]:m[11]], ",") {
				inner.implList = append(inner.implList, strings.TrimSpace(p))
			}
		}
		inner.javadoc = javadocBefore(body, m[0])
		braceStart := m[1] - 1
		if braceStart >= 0 && braceStart < len(body) && body[braceStart] == '{' {
			if innerBody, ok := braceBody(body, braceStart); ok {
				inner.body = innerBody
				inner.fields, inner.methods = parseBody(innerBody)
			}
		}
		result = append(result, inner)
	}
	return result
}

// javadocBefore returns the /** … */ doc comment that immediately
// precedes the declaration at offset `pos`.
// #172 — convert Javadoc tags to prose.
func javadocBefore(src string, pos int) string {
	end := strings.LastIndex(src[:pos], "*/")
	if end < 0 {
		return ""
	}
	start := strings.LastIndex(src[:end], "/**")
	if start < 0 {
		return ""
	}
	gap := src[end+2 : pos]
	// Allow whitespace and annotation lines (@Override, @Nullable, etc.)
	// between the closing */ and the declaration. Non-annotation, non-
	// whitespace content means this javadoc belongs to another declaration.
	if !isJavadocGapOK(gap) {
		return ""
	}
	raw := src[start+3 : end]
	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, " ")
		// #172 — convert @param, @return, @see to prose
		line = convertJavadocTag(line)
		lines = append(lines, line)
	}
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

// isJavadocGapOK returns true if the text between a `*/` and the next
// declaration contains only whitespace and Java annotation lines (@…).
func isJavadocGapOK(gap string) bool {
	for _, line := range strings.Split(gap, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Allow annotation lines such as @Override, @Nullable, @SuppressWarnings(...)
		if strings.HasPrefix(line, "@") {
			continue
		}
		return false
	}
	return true
}

// convertJavadocTag converts @param, @return, @see, @throws tags to prose.
// #172
func convertJavadocTag(line string) string {
	for _, tag := range []string{"@param", "@return", "@returns", "@see", "@throws", "@deprecated"} {
		if strings.HasPrefix(line, tag) {
			rest := strings.TrimSpace(line[len(tag):])
			switch tag {
			case "@param":
				return "Param: " + rest
			case "@return", "@returns":
				return "Returns: " + rest
			case "@see":
				return "See: " + rest
			case "@throws":
				return "Throws: " + rest
			case "@deprecated":
				return "Deprecated: " + rest
			}
		}
	}
	return line
}

// braceBody returns the text between matching braces starting at
// openIdx (the `{` offset, inclusive).
func braceBody(src string, openIdx int) (string, bool) {
	depth := 0
	for i := openIdx; i < len(src); i++ {
		switch src[i] {
		case '"', '\'':
			i = skipStringLit(src, i)
		case '/':
			if i+1 < len(src) && src[i+1] == '/' {
				i = skipLineComment(src, i)
			} else if i+1 < len(src) && src[i+1] == '*' {
				i = skipBlockComment(src, i)
			}
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return src[openIdx+1 : i], true
			}
		}
	}
	return "", false
}

func skipStringLit(src string, i int) int {
	quote := src[i]
	i++
	for i < len(src) {
		switch src[i] {
		case '\\':
			i += 2
			continue
		case quote:
			return i
		}
		i++
	}
	return len(src) - 1
}

func skipLineComment(src string, i int) int {
	for i < len(src) && src[i] != '\n' {
		i++
	}
	return i
}

func skipBlockComment(src string, i int) int {
	i += 2
	for i+1 < len(src) {
		if src[i] == '*' && src[i+1] == '/' {
			return i + 1
		}
		i++
	}
	return len(src) - 1
}

// parseBody splits a class body into field / method declarations.
func parseBody(body string) ([]fieldDecl, []methodDecl) {
	var fields []fieldDecl
	var methods []methodDecl
	pos := 0
	for pos < len(body) {
		pos = skipWS(body, pos)
		if pos >= len(body) {
			break
		}
		if strings.HasPrefix(body[pos:], "//") {
			nl := strings.Index(body[pos:], "\n")
			if nl < 0 {
				break
			}
			pos += nl + 1
			continue
		}
		if strings.HasPrefix(body[pos:], "/*") {
			end := strings.Index(body[pos:], "*/")
			if end < 0 {
				break
			}
			pos += end + 2
			continue
		}
		if body[pos] == '@' {
			nl := strings.Index(body[pos:], "\n")
			if nl < 0 {
				break
			}
			line := body[pos : pos+nl]
			if !strings.Contains(line, "(") {
				pos += nl + 1
				continue
			}
		}
		doc := javadocBefore(body, pos)
		srcLine := lineNumber(body, pos)
		semi, brace := findDeclTerminator(body[pos:])
		if semi < 0 && brace < 0 {
			break
		}
		var advance int
		if semi >= 0 && (brace < 0 || semi < brace) {
			decl := body[pos : pos+semi]
			advance = semi + 1
			if f, ok := parseField(decl, doc); ok {
				fields = append(fields, f)
			} else if strings.Contains(decl, "(") {
				// Abstract / interface method declaration: `abstract int x();`
				if m, ok := parseMethod(decl, "", doc); ok {
					m.sourceLine = srcLine
					methods = append(methods, m)
				}
			}
		} else {
			sig := body[pos : pos+brace]
			methodBody, ok := braceBody(body, pos+brace)
			if !ok {
				break
			}
			advance = brace + 2 + len(methodBody)
			if m, ok := parseMethod(sig, methodBody, doc); ok {
				m.sourceLine = srcLine
				methods = append(methods, m)
			}
		}
		pos += advance
	}
	return fields, methods
}

// findDeclTerminator scans `s` from the start and returns (semiIdx,
// braceIdx) — whichever comes first at paren-depth zero.
//
// We deliberately do NOT track `<`/`>` as nesting because Java uses them
// both for generic type arguments AND as arithmetic/shift operators (`<<`,
// `>>`, `>>>`, `<=`, `>=`). Tracking `<`/`>` causes the depth counter to
// desync on expressions like `(1L << 53)`, hiding the terminating `;` or `{`.
// Generic angle-brackets never contain `;` or `{`, so omitting them is safe.
func findDeclTerminator(s string) (int, int) {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\'':
			i = skipStringLit(s, i)
		case '/':
			if i+1 < len(s) && s[i+1] == '/' {
				i = skipLineComment(s, i)
			} else if i+1 < len(s) && s[i+1] == '*' {
				i = skipBlockComment(s, i)
			}
		case '(', '[':
			depth++
		case ')', ']':
			if depth > 0 {
				depth--
			}
		case ';':
			if depth == 0 {
				return i, -1
			}
		case '{':
			if depth == 0 {
				return -1, i
			}
		}
	}
	return -1, -1
}

func indexAtDepthZero(s string, b byte) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		case '{':
			if b == '{' && depth == 0 {
				return i
			}
			return -1
		case b:
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func skipWS(s string, i int) int {
	for i < len(s) {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
		}
		break
	}
	return i
}

// parseField recognises `[modifiers] Type name [= init]` pattern.
func parseField(decl, doc string) (fieldDecl, bool) {
	decl = strings.TrimSpace(decl)
	if decl == "" {
		return fieldDecl{}, false
	}
	// Strip initialiser first; if there is none and the declaration ends with
	// `)` it must be a method signature, not a field.
	var init string
	if eq := indexAssignOp(decl); eq >= 0 {
		init = strings.TrimSpace(decl[eq+1:])
		decl = strings.TrimSpace(decl[:eq])
	} else if strings.HasSuffix(decl, ")") {
		// No initialiser and ends with `)` → method signature, not a field.
		return fieldDecl{}, false
	}
	mods, rest := peelModifiers(decl)
	tn := strings.TrimSpace(rest)
	sp := lastTopLevelSpace(tn)
	if sp < 0 {
		return fieldDecl{}, false
	}
	return fieldDecl{
		javadoc:     doc,
		modifiers:   mods,
		typeName:    strings.TrimSpace(tn[:sp]),
		name:        strings.TrimSpace(tn[sp+1:]),
		initialiser: init,
	}, true
}

// parseMethod recognises `[modifiers] [returnType] name(params) [throws …]`.
func parseMethod(sig, body, doc string) (methodDecl, bool) {
	sig = strings.TrimSpace(sig)
	var throws string
	if i := strings.Index(sig, " throws "); i >= 0 {
		throws = strings.TrimSpace(sig[i+len(" throws "):])
		sig = strings.TrimSpace(sig[:i])
	}
	open := strings.Index(sig, "(")
	if open < 0 {
		return methodDecl{}, false
	}
	close := strings.LastIndex(sig, ")")
	if close < 0 || close < open {
		return methodDecl{}, false
	}
	params := parseParams(sig[open+1 : close])
	head := strings.TrimSpace(sig[:open])
	mods, rest := peelModifiers(head)
	// #166 — strip generic type params from method-level `<T> T foo(...)`
	rest = strings.TrimSpace(rest)
	if strings.HasPrefix(rest, "<") {
		end := strings.Index(rest, ">")
		if end >= 0 {
			rest = strings.TrimSpace(rest[end+1:])
		}
	}
	sp := lastTopLevelSpace(rest)
	var ret, name string
	if sp < 0 {
		name = rest
	} else {
		ret = strings.TrimSpace(rest[:sp])
		name = strings.TrimSpace(rest[sp+1:])
	}
	return methodDecl{
		javadoc:    doc,
		modifiers:  mods,
		returnType: ret,
		name:       name,
		params:     params,
		throws:     throws,
		body:       body,
	}, true
}

func parseParams(s string) []paramDecl {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []paramDecl
	depth := 0
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || (s[i] == ',' && depth == 0) {
			piece := strings.TrimSpace(s[start:i])
			if piece != "" {
				// Strip final/annotations from param
				_, piece = peelModifiers(piece)
				// Handle varargs: `Type... name` → treat as array
				piece = strings.ReplaceAll(piece, "...", "[]")
				sp := lastTopLevelSpace(piece)
				if sp > 0 {
					out = append(out, paramDecl{
						typeName: strings.TrimSpace(piece[:sp]),
						name:     strings.TrimSpace(piece[sp+1:]),
					})
				}
			}
			start = i + 1
			continue
		}
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		}
	}
	return out
}

var modifierSet = map[string]bool{
	"public": true, "private": true, "protected": true, "static": true,
	"final": true, "abstract": true, "synchronized": true, "volatile": true,
	"transient": true, "native": true, "default": true,
}

func peelModifiers(s string) (string, string) {
	var mods []string
	for {
		s = strings.TrimSpace(s)
		sp := firstTopLevelSpace(s)
		if sp < 0 {
			break
		}
		head := s[:sp]
		if modifierSet[head] || strings.HasPrefix(head, "@") {
			mods = append(mods, head)
			s = s[sp+1:]
			continue
		}
		break
	}
	return strings.Join(mods, " "), s
}

func firstTopLevelSpace(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		case ' ':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func lastTopLevelSpace(s string) int {
	depth := 0
	last := -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		case ' ':
			if depth == 0 {
				last = i
			}
		}
	}
	return last
}

func indexAssignOp(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', '<', '[':
			depth++
		case ')', '>', ']':
			depth--
		case '=':
			if depth == 0 {
				if i+1 < len(s) && s[i+1] == '=' {
					continue
				}
				return i
			}
		}
	}
	return -1
}

// ── #164 Import block ────────────────────────────────────────────────────────

// javaImportToGoPackage maps Java fully-qualified class prefixes to Go import paths.
var javaImportToGoPackage = map[string]string{
	"com.google.javascript.jscomp.":    "github.com/paullepoulpe/jsgo/internal/jscomp",
	"com.google.javascript.rhino.":     "github.com/paullepoulpe/jsgo/internal/rhino",
	"com.google.javascript.rhino.jstype.": "github.com/paullepoulpe/jsgo/internal/rhino/jstype",
	"com.google.javascript.jscomp.parsing.": "github.com/paullepoulpe/jsgo/internal/parser",
	"com.google.javascript.jscomp.diagnostic.": "github.com/paullepoulpe/jsgo/internal/diag",
}

// parseJavaImports extracts the Java import lines from src.
func parseJavaImports(src string) []string {
	var imports []string
	importRe := regexp.MustCompile(`(?m)^import\s+(static\s+)?([\w.]+);`)
	for _, m := range importRe.FindAllStringSubmatch(src, -1) {
		imports = append(imports, m[2])
	}
	return imports
}

// javaImportsToGoPkgs converts Java import lines to Go package paths.
func javaImportsToGoPkgs(javaImports []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, imp := range javaImports {
		for prefix, goPkg := range javaImportToGoPackage {
			if strings.HasPrefix(imp, prefix) {
				if !seen[goPkg] {
					seen[goPkg] = true
					out = append(out, `"`+goPkg+`"`)
				}
			}
		}
	}
	return out
}

// importBlockFor emits the Go import block. It uses the Java import
// lines to derive correct Go packages (#164) plus heuristic checks
// for well-known types.
func importBlockFor(src, javaPath string) string {
	seen := map[string]bool{}
	var needed []string

	add := func(pkg string) {
		q := `"` + pkg + `"`
		if !seen[q] {
			seen[q] = true
			needed = append(needed, q)
		}
	}

	// #164 — map Java imports
	for _, pkg := range javaImportsToGoPkgs(parseJavaImports(src)) {
		p := strings.Trim(pkg, `"`)
		add(p)
	}

	// Heuristic fallbacks for types referenced without import
	if strings.Contains(src, "Node ") || strings.Contains(src, "Node.") {
		add("github.com/paullepoulpe/jsgo/internal/rhino")
	}
	if strings.Contains(src, "JSType") {
		add("github.com/paullepoulpe/jsgo/internal/rhino/jstype")
	}
	if strings.Contains(src, "JSError") || strings.Contains(src, "DiagnosticType") {
		add("github.com/paullepoulpe/jsgo/internal/diag")
	}

	if len(needed) == 0 {
		return ""
	}
	return "import (\n\t" + strings.Join(needed, "\n\t") + "\n)\n\n"
}

// ── Type mapping ─────────────────────────────────────────────────────────────

// javaTypeMap is the Java→Go type translation table.
var javaTypeMap = map[string]string{
	"":                 "",
	"String":           "string",
	"boolean":          "bool",
	"int":              "int",
	"long":             "int64",
	"short":            "int16",
	"byte":             "int8",
	"char":             "rune",
	"double":           "float64",
	"float":            "float32",
	"void":             "",
	"Object":           "any",
	"Node":             "*rhino.Node",
	"JSDocInfo":        "*rhino.JSDocInfo",
	"JSType":           "jstype.JSType",
	"JSTypeExpression": "*rhino.JSTypeExpression",
	"FunctionType":     "*jstype.FunctionType",
	"ObjectType":       "*jstype.ObjectType",
	"UnionType":        "*jstype.UnionType",
	"TemplatizedType":  "*jstype.TemplatizedType",
	"Token":            "rhino.Token",
	"JSError":          "diag.Diagnostic",
	"DiagnosticType":   "string",
	"AbstractCompiler": "*Compiler",
	"NodeTraversal":    "*NodeTraversal",
	"ControlFlowGraph": "*ControlFlowGraph",
	"CompilerPass":     "CompilerPass",
	"Scope":            "*TypedScope",
	"StaticSourceFile": "string",
}

// javaToGoType maps a Java type expression to its Go equivalent.
// #166 handles generics; #167 handles Guava collections.
func javaToGoType(t string) string {
	t = stripTypeAnnotations(strings.TrimSpace(t))
	for _, prefix := range []string{"final ", "transient "} {
		t = strings.TrimPrefix(t, prefix)
	}
	if v, ok := javaTypeMap[t]; ok {
		return v
	}
	// Array types → Go slice.
	if strings.HasSuffix(t, "[]") {
		return "[]" + javaToGoType(strings.TrimSuffix(t, "[]"))
	}
	// #167 — Guava/JDK collections
	if strings.HasPrefix(t, "List<") || strings.HasPrefix(t, "ImmutableList<") ||
		strings.HasPrefix(t, "Collection<") || strings.HasPrefix(t, "Iterable<") ||
		strings.HasPrefix(t, "ArrayList<") || strings.HasPrefix(t, "LinkedList<") {
		return "[]" + javaToGoType(innerGeneric(t))
	}
	if strings.HasPrefix(t, "Set<") || strings.HasPrefix(t, "ImmutableSet<") ||
		strings.HasPrefix(t, "HashSet<") || strings.HasPrefix(t, "LinkedHashSet<") {
		return "map[" + javaToGoType(innerGeneric(t)) + "]struct{}"
	}
	if strings.HasPrefix(t, "Map<") || strings.HasPrefix(t, "ImmutableMap<") ||
		strings.HasPrefix(t, "HashMap<") || strings.HasPrefix(t, "LinkedHashMap<") {
		inner := innerGeneric(t)
		parts := strings.SplitN(inner, ",", 2)
		if len(parts) == 2 {
			return "map[" + javaToGoType(strings.TrimSpace(parts[0])) + "]" + javaToGoType(strings.TrimSpace(parts[1]))
		}
	}
	// #167 — Optional<T> → *T
	if strings.HasPrefix(t, "Optional<") {
		inner := javaToGoType(innerGeneric(t))
		if strings.HasPrefix(inner, "*") {
			return inner
		}
		return "*" + inner
	}
	// Generic of a mapped base type
	if lt := strings.Index(t, "<"); lt > 0 && strings.HasSuffix(t, ">") {
		base := t[:lt]
		if v, ok := javaTypeMap[base]; ok {
			return v
		}
	}
	return t
}

func stripTypeAnnotations(t string) string {
	for strings.HasPrefix(t, "@") {
		i := 1
		for i < len(t) && (isJavaIdentChar(t[i]) || t[i] == '.') {
			i++
		}
		if i < len(t) && t[i] == '(' {
			depth := 1
			i++
			for i < len(t) && depth > 0 {
				if t[i] == '(' {
					depth++
				} else if t[i] == ')' {
					depth--
				}
				i++
			}
		}
		t = strings.TrimLeft(t[i:], " \t")
	}
	return t
}

func isJavaIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '_' || c == '$'
}

func isMappedType(javaType string) bool {
	javaType = stripTypeAnnotations(strings.TrimSpace(javaType))
	if _, ok := javaTypeMap[javaType]; ok {
		return true
	}
	if strings.HasSuffix(javaType, "[]") {
		return isMappedType(strings.TrimSuffix(javaType, "[]"))
	}
	if strings.HasPrefix(javaType, "List<") ||
		strings.HasPrefix(javaType, "ImmutableList<") ||
		strings.HasPrefix(javaType, "Collection<") ||
		strings.HasPrefix(javaType, "Iterable<") ||
		strings.HasPrefix(javaType, "ArrayList<") ||
		strings.HasPrefix(javaType, "Set<") ||
		strings.HasPrefix(javaType, "ImmutableSet<") ||
		strings.HasPrefix(javaType, "HashSet<") ||
		strings.HasPrefix(javaType, "Map<") ||
		strings.HasPrefix(javaType, "ImmutableMap<") ||
		strings.HasPrefix(javaType, "HashMap<") ||
		strings.HasPrefix(javaType, "Optional<") {
		return true
	}
	if lt := strings.Index(javaType, "<"); lt > 0 && strings.HasSuffix(javaType, ">") {
		if _, ok := javaTypeMap[javaType[:lt]]; ok {
			return true
		}
	}
	return false
}

func innerGeneric(t string) string {
	open := strings.Index(t, "<")
	close := strings.LastIndex(t, ">")
	if open < 0 || close < 0 || close <= open {
		return t
	}
	return strings.TrimSpace(t[open+1 : close])
}

// ── Diagnostic constants ─────────────────────────────────────────────────────

var diagRe = regexp.MustCompile(
	`(?s)DiagnosticType\.(error|warning|disabled)\(\s*"(JSC_[A-Z_0-9]+)"\s*,\s*([^)]+)\)`)

func emitDiagnosticConsts(cls *classDecl, src string) string {
	matches := diagRe.FindAllStringSubmatch(src, -1)
	if len(matches) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("// Diagnostic codes discovered in the Java source.\n")
	b.WriteString("// Add the matching entries to internal/diag/diag.go.\n")
	b.WriteString("//\n")
	for _, m := range matches {
		b.WriteString("//   ")
		b.WriteString(m[2])
		b.WriteString(" (")
		b.WriteString(m[1])
		b.WriteString(")\n")
	}
	b.WriteString("\n")
	return b.String()
}

// ── #168 AutoValue ───────────────────────────────────────────────────────────

// emitAutoValue emits a Go struct + factory for @AutoValue classes.
func emitAutoValue(cls *classDecl, javaPath string) string {
	if !cls.isAutoValue {
		return ""
	}
	var b strings.Builder
	// Doc comment
	if cls.javadoc != "" {
		for _, line := range strings.Split(cls.javadoc, "\n") {
			b.WriteString("// " + line + "\n")
		}
	}
	b.WriteString("// AutoValue: value-based equality; compare with == or reflect.DeepEqual.\n")
	tpStr := typeParamsGoString(cls.typeParams)
	b.WriteString("type " + cls.name + tpStr + " struct {\n")
	for _, f := range cls.fields {
		goType := javaToGoType(f.typeName)
		if f.javadoc != "" {
			for _, l := range strings.Split(f.javadoc, "\n") {
				b.WriteString("\t// " + l + "\n")
			}
		}
		b.WriteString("\t" + toExported(f.name) + " " + goType + "\n")
	}
	b.WriteString("}\n\n")
	// Abstract accessor methods become fields; emit factory
	b.WriteString("// New" + cls.name + " constructs a " + cls.name + " (AutoValue factory).\n")
	b.WriteString("func New" + cls.name + "(")
	first := true
	for _, m := range cls.methods {
		if strings.Contains(m.modifiers, "abstract") && len(m.params) == 0 && m.returnType != "" {
			if !first {
				b.WriteString(", ")
			}
			b.WriteString(m.name + " " + javaToGoType(m.returnType))
			first = false
		}
	}
	b.WriteString(") *" + cls.name + " {\n")
	b.WriteString("\treturn &" + cls.name + "{\n")
	for _, m := range cls.methods {
		if strings.Contains(m.modifiers, "abstract") && len(m.params) == 0 && m.returnType != "" {
			b.WriteString("\t\t" + toExported(m.name) + ": " + m.name + ",\n")
		}
	}
	b.WriteString("\t}\n}\n\n")
	return b.String()
}

// ── #163 Enum-with-methods ───────────────────────────────────────────────────

// emitEnum emits a Go interface + per-constant struct for Java enums.
func emitEnum(cls *classDecl, javaPath string) string {
	if !cls.isEnum {
		return ""
	}
	var b strings.Builder
	if cls.javadoc != "" {
		for _, line := range strings.Split(cls.javadoc, "\n") {
			b.WriteString("// " + line + "\n")
		}
	}

	// Determine if any constant has a body (methods) — if so, emit interface+structs.
	hasMethodBodies := false
	for _, ec := range cls.enumConsts {
		if ec.body != "" {
			hasMethodBodies = true
			break
		}
	}

	if hasMethodBodies || len(cls.methods) > 0 {
		// Rich enum: interface + per-constant struct + factory
		b.WriteString("// " + cls.name + " is a Java enum with methods; represented as an interface.\n")
		b.WriteString("type " + cls.name + " interface {\n")
		for _, m := range cls.methods {
			if strings.Contains(m.modifiers, "private") {
				continue
			}
			goName := toExported(m.name)
			b.WriteString("\t" + goName + "(")
			for i, p := range m.params {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(p.name + " " + javaToGoType(p.typeName))
			}
			b.WriteString(")")
			if ret := javaToGoType(m.returnType); ret != "" {
				b.WriteString(" " + ret)
			}
			b.WriteString("\n")
		}
		b.WriteString("\t// Name returns the enum constant name.\n")
		b.WriteString("\tName() string\n")
		b.WriteString("}\n\n")

		// Per-constant structs
		for _, ec := range cls.enumConsts {
			structName := cls.name + "_" + ec.name
			if ec.javadoc != "" {
				for _, l := range strings.Split(ec.javadoc, "\n") {
					b.WriteString("// " + l + "\n")
				}
			}
			b.WriteString("type " + structName + " struct{}\n\n")
			b.WriteString("func (" + strings.ToLower(string(ec.name[0])) + " *" + structName + ") Name() string {\n")
			b.WriteString("\treturn \"" + ec.name + "\"\n}\n\n")
			// Stub out interface methods per constant
			for _, m := range cls.methods {
				if strings.Contains(m.modifiers, "private") {
					continue
				}
				goName := toExported(m.name)
				recv := strings.ToLower(string(ec.name[0]))
				b.WriteString("// @ " + javaPath + ":" + fmt.Sprint(m.sourceLine) + "\n")
				b.WriteString("func (" + recv + " *" + structName + ") " + goName + "(")
				for i, p := range m.params {
					if i > 0 {
						b.WriteString(", ")
					}
					b.WriteString(p.name + " " + javaToGoType(p.typeName))
				}
				b.WriteString(")")
				goRet := javaToGoType(m.returnType)
				if goRet != "" {
					b.WriteString(" " + goRet)
				}
				b.WriteString(" {\n")
				sig := buildPanicSig(m)
				b.WriteString("\tpanic(\"unimplemented: " + javaPath + ":" + fmt.Sprint(m.sourceLine) + " " + sig + "\")\n")
				b.WriteString("}\n\n")
			}
		}

		// Factory
		b.WriteString("// New" + cls.name + " returns the " + cls.name + " for a given constant name.\n")
		b.WriteString("func New" + cls.name + "(name string) " + cls.name + " {\n")
		b.WriteString("\tswitch name {\n")
		for _, ec := range cls.enumConsts {
			b.WriteString("\tcase \"" + ec.name + "\":\n")
			b.WriteString("\t\treturn &" + cls.name + "_" + ec.name + "{}\n")
		}
		b.WriteString("\t}\n")
		b.WriteString("\treturn nil\n}\n\n")

		// Sentinel vars for simple access
		b.WriteString("var (\n")
		for _, ec := range cls.enumConsts {
			b.WriteString("\t" + cls.name + ec.name + " " + cls.name + " = &" + cls.name + "_" + ec.name + "{}\n")
		}
		b.WriteString(")\n\n")
	} else {
		// Simple value-only enum: iota const block
		b.WriteString("type " + cls.name + " int\n\n")
		b.WriteString("const (\n")
		for i, ec := range cls.enumConsts {
			if ec.javadoc != "" {
				for _, l := range strings.Split(ec.javadoc, "\n") {
					b.WriteString("\t// " + l + "\n")
				}
			}
			if i == 0 {
				b.WriteString("\t" + ec.name + " " + cls.name + " = iota\n")
			} else {
				b.WriteString("\t" + ec.name + "\n")
			}
		}
		b.WriteString(")\n\n")
	}
	return b.String()
}

// ── #162 Inner class emission ────────────────────────────────────────────────

func emitInnerClass(outerName string, inner *classDecl, javaPath string) string {
	// Prefix inner class name with outer name: OuterInner
	prefixedName := outerName + inner.name
	// Save original name and swap for emission
	orig := inner.name
	inner.name = prefixedName
	defer func() { inner.name = orig }()

	var b strings.Builder
	if inner.javadoc != "" {
		for _, line := range strings.Split(inner.javadoc, "\n") {
			b.WriteString("// " + line + "\n")
		}
	}
	b.WriteString("// Inner class " + outerName + "." + orig + " — ported as " + prefixedName + ".\n")
	tpStr := typeParamsGoString(inner.typeParams)
	b.WriteString("type " + prefixedName + tpStr + " struct {\n")
	for _, f := range inner.fields {
		if strings.Contains(f.modifiers, "static") {
			continue
		}
		goType := javaToGoType(f.typeName)
		if f.javadoc != "" {
			for _, l := range strings.Split(f.javadoc, "\n") {
				b.WriteString("\t// " + l + "\n")
			}
		}
		if !isMappedType(f.typeName) {
			b.WriteString("\t// TODO: Java type " + f.typeName + " not mapped\n")
		}
		b.WriteString("\t" + toExported(f.name) + " " + goType + "\n")
	}
	b.WriteString("}\n\n")
	b.WriteString(emitMethods(inner, javaPath))
	return b.String()
}

// ── Struct + methods emission ────────────────────────────────────────────────

func emitStruct(cls *classDecl) string {
	if cls.isEnum || cls.isAutoValue {
		return ""
	}
	var b strings.Builder
	if cls.javadoc != "" {
		for _, line := range strings.Split(cls.javadoc, "\n") {
			b.WriteString("// " + line + "\n")
		}
	}
	tpStr := typeParamsGoString(cls.typeParams)
	b.WriteString("type " + cls.name + tpStr + " struct {\n")
	for _, f := range cls.fields {
		if strings.Contains(f.modifiers, "static") {
			continue
		}
		goType := javaToGoType(f.typeName)
		if f.javadoc != "" {
			for _, l := range strings.Split(f.javadoc, "\n") {
				b.WriteString("\t// " + l + "\n")
			}
		}
		if !isMappedType(f.typeName) {
			b.WriteString("\t// TODO: Java type " + f.typeName + " not mapped\n")
		}
		b.WriteString("\t" + toExported(f.name) + " " + goType + "\n")
	}
	b.WriteString("}\n\n")
	return b.String()
}

func emitMethods(cls *classDecl, javaPath string) string {
	if cls.isEnum {
		return ""
	}
	var b strings.Builder
	for _, m := range cls.methods {
		if strings.Contains(m.modifiers, "static") {
			// Emit static methods as package-level functions (Go has no static
			// methods; Java statics become free functions or methods on a
			// singleton, depending on the class).
			emitStaticFunc(&b, cls, m, javaPath)
			continue
		}
		if m.name == cls.name {
			emitConstructor(&b, cls, m, javaPath)
			continue
		}
		emitMethodStub(&b, cls, m, javaPath)
	}
	return b.String()
}

// emitStaticFunc emits a Java static method as a Go package-level function.
func emitStaticFunc(b *strings.Builder, cls *classDecl, m methodDecl, javaPath string) {
	if m.javadoc != "" {
		for _, l := range strings.Split(m.javadoc, "\n") {
			b.WriteString("// " + l + "\n")
		}
	}
	if m.sourceLine > 0 {
		b.WriteString("// @ " + javaPath + ":" + fmt.Sprint(m.sourceLine) + "\n")
	}
	goName := toExported(m.name)
	if strings.Contains(m.modifiers, "private") {
		goName = strings.ToLower(m.name[:1]) + m.name[1:]
	}
	b.WriteString("func " + goName + "(")
	for i, p := range m.params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(goIdent(p.name) + " " + javaToGoType(p.typeName))
	}
	b.WriteString(")")
	goRet := javaToGoType(m.returnType)
	if goRet != "" {
		b.WriteString(" " + goRet)
	}
	b.WriteString(" {\n")
	if m.throws != "" {
		b.WriteString("\t// Java threw: " + m.throws + "\n")
	}
	recv := strings.ToLower(string(cls.name[0]))
	if bodyText, ok := translateMethodBody(m.body, recv, cls.name, goRet); ok {
		b.WriteString(bodyText)
	} else {
		sig := buildPanicSig(m)
		b.WriteString("\tpanic(\"unimplemented: " + javaPath + ":" + fmt.Sprint(m.sourceLine) + " " + sig + "\")\n")
	}
	b.WriteString("}\n\n")
}

// buildPanicSig builds a human-readable method signature for the panic message.
// #169
func buildPanicSig(m methodDecl) string {
	var sb strings.Builder
	sb.WriteString(m.name + "(")
	for i, p := range m.params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.typeName + " " + p.name)
	}
	sb.WriteString(")")
	if m.returnType != "" {
		sb.WriteString(" " + m.returnType)
	}
	return sb.String()
}

func emitConstructor(b *strings.Builder, cls *classDecl, m methodDecl, javaPath string) {
	if m.javadoc != "" {
		for _, l := range strings.Split(m.javadoc, "\n") {
			b.WriteString("// " + l + "\n")
		}
	}
	// #170 — source-line pin
	if m.sourceLine > 0 {
		b.WriteString("// @ " + javaPath + ":" + fmt.Sprint(m.sourceLine) + "\n")
	}
	b.WriteString("func New" + cls.name + "(")
	for i, p := range m.params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(goIdent(p.name) + " " + javaToGoType(p.typeName))
	}
	b.WriteString(") *" + cls.name + " {\n")
	recv := strings.ToLower(string(cls.name[0]))
	if bodyText, ok := translateMethodBody(m.body, recv, cls.name, ""); ok {
		if !strings.Contains(bodyText, "return") && !strings.Contains(bodyText, "panic(") {
			// Body is field-init assignments (this.X = Y pattern). Declare the
			// receiver first so that all assignments target a valid variable.
			if strings.Contains(bodyText, recv+".") {
				b.WriteString("\t" + recv + " := &" + cls.name + "{}\n")
				b.WriteString(bodyText)
				b.WriteString("\treturn " + recv + "\n")
			} else {
				b.WriteString(bodyText)
				b.WriteString("\treturn &" + cls.name + "{}\n")
			}
		} else {
			b.WriteString(bodyText)
		}
	} else {
		// #169 — panic stub instead of bare return
		sig := buildPanicSig(m)
		b.WriteString("\tpanic(\"unimplemented: " + javaPath + ":" + fmt.Sprint(m.sourceLine) + " " + sig + "\")\n")
	}
	b.WriteString("}\n\n")
}

func emitMethodStub(b *strings.Builder, cls *classDecl, m methodDecl, javaPath string) {
	recv := "c"
	if strings.Contains(m.modifiers, "private") {
		recv = strings.ToLower(string(cls.name[0]))
	}
	if m.javadoc != "" {
		for _, l := range strings.Split(m.javadoc, "\n") {
			b.WriteString("// " + l + "\n")
		}
	}
	// #170 — source-line pin
	if m.sourceLine > 0 {
		b.WriteString("// @ " + javaPath + ":" + fmt.Sprint(m.sourceLine) + "\n")
	}
	goName := toExported(m.name)
	if strings.Contains(m.modifiers, "private") {
		goName = strings.ToLower(m.name[:1]) + m.name[1:]
	}
	b.WriteString("func (" + recv + " *" + cls.name + ") " + goName + "(")
	for i, p := range m.params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(goIdent(p.name) + " " + javaToGoType(p.typeName))
	}
	b.WriteString(")")
	goRet := javaToGoType(m.returnType)
	if goRet != "" {
		b.WriteString(" " + goRet)
	}
	b.WriteString(" {\n")
	if m.throws != "" {
		b.WriteString("\t// Java threw: " + m.throws + "\n")
	}
	if bodyText, ok := translateMethodBody(m.body, recv, cls.name, goRet); ok {
		b.WriteString(bodyText)
	} else {
		// #169 — panic stub
		sig := buildPanicSig(m)
		b.WriteString("\tpanic(\"unimplemented: " + javaPath + ":" + fmt.Sprint(m.sourceLine) + " " + sig + "\")\n")
	}
	b.WriteString("}\n\n")
}

// toExported capitalises the first rune of a name.
func toExported(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

// goIdent escapes Java identifiers that are reserved words in Go.
func goIdent(name string) string {
	switch name {
	case "type":
		return "typ"
	case "len":
		return "length"
	case "cap":
		return "capacity"
	case "range":
		return "rng"
	case "map":
		return "mp"
	case "var":
		return "v"
	case "func":
		return "fn"
	case "package":
		return "pkg"
	case "interface":
		return "iface"
	case "chan":
		return "ch"
	case "select":
		return "sel"
	case "go":
		return "goroutine"
	case "defer":
		return "deferred"
	}
	return name
}

// translateMethodBody attempts to translate the raw Java method body text
// (between the braces) into Go using the javaast package. It tries the
// trivial-pattern recogniser first, then falls back to the full AST walker.
// Returns (goBody, true) on success, or ("", false) when neither succeeded.
//
// goRet is the Go return type (empty string for void); used to synthesise a
// zero-value return when the Java body is empty.
func translateMethodBody(javaBody, recvName, className, goRet string) (string, bool) {
	body := strings.TrimSpace(javaBody)
	if body == "" {
		// Empty Java body — emit a zero-value return for non-void methods.
		switch goRet {
		case "":
			return "", true
		case "bool":
			return "\treturn false\n", true
		case "int", "int64", "float32", "float64":
			return "\treturn 0\n", true
		case "string":
			return "\treturn \"\"\n", true
		}
		return "\treturn nil\n", true
	}

	ctx := javaast.MethodContext{
		RecvName:  recvName,
		ClassName: className,
		TypeMap:   javaToGoType,
	}

	// Try trivial patterns first (getter/setter/delegation/null-check/field-init).
	if s, ok := javaast.TryTrivial(body, recvName, className, javaToGoType); ok {
		return s, true
	}

	// Fall back to the full recursive-descent AST walker.
	translated, ok := javaast.TranslateBody(body, ctx)
	if ok {
		if translated != "" {
			return translated, true
		}
		// Empty translation = body was all comments (e.g. "// Do nothing").
		// Treat as void body — add zero-value return for non-void methods.
		switch goRet {
		case "":
			return "", true
		case "bool":
			return "\treturn false\n", true
		case "int", "int64", "float32", "float64":
			return "\treturn 0\n", true
		case "string":
			return "\treturn \"\"\n", true
		default:
			return "\treturn nil\n", true
		}
	}

	return "", false
}
