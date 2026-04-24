// tools/portcheck/main.go — measures per-file fidelity of the Go port against the Java source.
//
// For each Go file in internal/ with a "// Port of <java-path>" header,
// it computes a fidelity score (0–100) and emits docs/port-fidelity.md and
// docs/port-fidelity.json.
//
// Score formula:
//   weighted(methodNameOverlap*25 + methodCountRatio*15 + sigSimilarity*10 + bodyLocRatio*20) / 70
//   + 30 baseline
//   - goOnlyPenalty (up to 20)
//   - stubPenalty   (up to 30)
//   - versionPenalty (5 if version mismatch)
//   clamped to [0,100].

package main

import (
	"encoding/json"
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
)

const (
	javaSourceRoot = "/Users/damien/code/jsgo/.cache/closure-src"
	expectedVersion = "v20260415"
)

// -------- data structures --------

// FileResult holds computed metrics for one Go file.
type FileResult struct {
	GoFile          string  `json:"go_file"`
	JavaFile        string  `json:"java_file"`
	Score           float64 `json:"score"`
	MethodOverlapPct float64 `json:"method_overlap_pct"`
	GoMethodCount   int     `json:"go_method_count"`
	JavaMethodCount int     `json:"java_method_count"`
	StubCount       int     `json:"stub_count"`
	GoBodyLoC       int     `json:"go_body_loc"`
	JavaBodyLoC     int     `json:"java_body_loc"`
	LocRatio        float64 `json:"loc_ratio"`
	GoOnlyCount     int     `json:"go_only_count"`
	Notes           string  `json:"notes"`
	// Internal breakdown
	WeightedScore   float64 `json:"weighted_score"`
	GoOnlyPenalty   float64 `json:"go_only_penalty"`
	StubPenalty     float64 `json:"stub_penalty"`
	VersionPenalty  float64 `json:"version_penalty"`
}

// -------- Java parsing --------

// javaMethodRe matches Java method declarations after annotation stripping.
// Captures: (modifiers) ReturnType methodName ( params
var javaMethodRe = regexp.MustCompile(
	`(?m)^[ \t]*(?:(?:public|private|protected|static|final|synchronized|abstract|native|default|override|@\w+(?:\([^)]*\))?)\s+)*` +
		`(?:<[^>]+>\s*)?` + // optional leading generic type params (e.g. <T>)
		`(?:[\w\$][\w<>\[\],\s]*?\s+)` + // return type
		`(\w+)\s*\(`)

// annotationLineRe matches Java lines that are solely an annotation.
var annotationLineRe = regexp.MustCompile(`(?m)^[ \t]*@\w+(?:\([^)]*\))?[ \t]*$`)

// stripJavaAnnotationLines removes standalone annotation lines (e.g. @Override, @Nullable)
// from Java source so that the method regex can match multi-line annotation blocks.
func stripJavaAnnotationLines(src string) string {
	return annotationLineRe.ReplaceAllString(src, "")
}

// javaParentRe extracts the first extends clause from a Java class declaration.
var javaParentRe = regexp.MustCompile(`(?m)\bclass\s+\w+\s+extends\s+(\w+)`)

// extractJavaParentClass returns the immediate parent class name (lowercase), or "".
func extractJavaParentClass(src string) string {
	m := javaParentRe.FindStringSubmatch(src)
	if m == nil {
		return ""
	}
	return strings.ToLower(m[1])
}

// extractJavaMethods returns a map of lowercase method name → arity (param count approximation).
func extractJavaMethods(src string) map[string]int {
	stripped := stripJavaAnnotationLines(src)
	methods := make(map[string]int)
	matches := javaMethodRe.FindAllStringSubmatch(stripped, -1)
	for _, m := range matches {
		name := strings.ToLower(m[1])
		// Skip constructors, keywords, obvious non-method tokens
		if isJavaKeyword(name) {
			continue
		}
		methods[name]++
	}
	return methods
}

var javaKeywords = map[string]bool{
	"if": true, "for": true, "while": true, "switch": true, "catch": true,
	"return": true, "new": true, "throw": true, "else": true, "do": true,
	"try": true, "finally": true, "assert": true, "synchronized": true,
	"instanceof": true, "class": true, "interface": true, "enum": true,
}

func isJavaKeyword(s string) bool {
	return javaKeywords[s]
}

// javaBodyLoC counts non-blank non-comment lines in Java source.
func javaBodyLoC(src string) int {
	count := 0
	inBlockComment := false
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if inBlockComment {
			if strings.Contains(trimmed, "*/") {
				inBlockComment = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "/*") {
			inBlockComment = true
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
			continue
		}
		count++
	}
	return count
}

// -------- Go parsing --------

// GoFunc holds info about a parsed Go function/method.
type GoFunc struct {
	Name     string
	Params   int
	IsStub   bool
	IsGoOnly bool // declared intentional via "// Go-only:" comment
}

// extractGoFunctions parses a Go file and extracts all function and method declarations.
// It recognises functions preceded by a "// Go-only:" comment as intentional divergences
// that should not be penalised.
func extractGoFunctions(filename string) ([]GoFunc, int, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, 0, err
	}

	// Build a set of lines that carry "// Go-only:" comments so we can
	// match them to the function declaration that immediately follows.
	goOnlyLines := make(map[int]bool)
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, "Go-only:") {
				goOnlyLines[fset.Position(c.Pos()).Line] = true
			}
		}
	}

	var funcs []GoFunc
	bodyLoC := 0

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		name := fd.Name.Name
		params := 0
		if fd.Type.Params != nil {
			for _, p := range fd.Type.Params.List {
				if len(p.Names) == 0 {
					params++
				} else {
					params += len(p.Names)
				}
			}
		}
		isStub := false
		if fd.Body != nil {
			start := fset.Position(fd.Body.Lbrace).Line
			end := fset.Position(fd.Body.Rbrace).Line
			loc := end - start - 1
			if loc < 0 {
				loc = 0
			}
			bodyLoC += loc
			isStub = isStubBody(fd, fset)
		}

		// A function is "Go-only" if any comment on the line immediately preceding
		// the function declaration contains "Go-only:".
		funcLine := fset.Position(fd.Pos()).Line
		isGoOnly := false
		// Check the two lines before the func keyword for a Go-only comment.
		for delta := 1; delta <= 3; delta++ {
			if goOnlyLines[funcLine-delta] {
				isGoOnly = true
				break
			}
		}
		// Also check the function's own doc comment if it exists.
		if fd.Doc != nil {
			for _, c := range fd.Doc.List {
				if strings.Contains(c.Text, "Go-only:") {
					isGoOnly = true
				}
			}
		}

		funcs = append(funcs, GoFunc{Name: name, Params: params, IsStub: isStub, IsGoOnly: isGoOnly})
	}

	return funcs, bodyLoC, nil
}

// isStubBody returns true if the function body is a stub.
// Criteria:
// (a) empty body
// (b) single return with zero/nil/false
// (c) panic("not implemented") / panic("unimplemented")
// (d) only TODO comment(s)
func isStubBody(fd *ast.FuncDecl, fset *token.FileSet) bool {
	if fd.Body == nil {
		return true
	}
	stmts := fd.Body.List
	if len(stmts) == 0 {
		return true
	}
	if len(stmts) == 1 {
		stmt := stmts[0]
		// Check for return zero/nil/false
		if ret, ok := stmt.(*ast.ReturnStmt); ok {
			if len(ret.Results) == 0 {
				return true
			}
			if len(ret.Results) == 1 {
				switch v := ret.Results[0].(type) {
				case *ast.Ident:
					if v.Name == "nil" || v.Name == "false" {
						return true
					}
				case *ast.BasicLit:
					if v.Value == "0" || v.Value == `""` {
						return true
					}
				}
			}
			// Multiple returns, all zero-ish
			allZero := true
			for _, r := range ret.Results {
				switch v := r.(type) {
				case *ast.Ident:
					if v.Name != "nil" && v.Name != "false" {
						allZero = false
					}
				case *ast.BasicLit:
					if v.Value != "0" && v.Value != `""` {
						allZero = false
					}
				default:
					allZero = false
				}
			}
			if allZero && len(ret.Results) > 0 {
				return true
			}
		}
		// Check for panic("not implemented") / panic("unimplemented")
		if expr, ok := stmt.(*ast.ExprStmt); ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
					if len(call.Args) == 1 {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							low := strings.ToLower(lit.Value)
							if strings.Contains(low, "not implemented") || strings.Contains(low, "unimplemented") {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

// goBodyLoC counts actual body lines (already computed during AST walk).
// This function is kept for clarity; we pass in the precomputed value.

// -------- header parsing --------

type headerInfo struct {
	javaPath string   // e.g. "com/google/javascript/jscomp/TypeCheck.java"
	version  string   // e.g. "v20260415"
	valid    bool
	multi    bool     // true if header references multiple Java files
	allPaths []string // all referenced java paths
}

// portHeaderRe matches lines like: // Port of com/google/javascript/...java @ closure-compiler v20260415
var portHeaderRe = regexp.MustCompile(`// Port of (com/google/javascript/\S+\.java)`)
var versionRe = regexp.MustCompile(`@ closure-compiler (v\d+)`)

func parseHeader(content string) headerInfo {
	lines := strings.Split(content, "\n")
	// Look at first 10 lines
	headerLines := lines
	if len(headerLines) > 10 {
		headerLines = lines[:10]
	}
	header := strings.Join(headerLines, "\n")

	matches := portHeaderRe.FindAllStringSubmatch(header, -1)
	if len(matches) == 0 {
		return headerInfo{valid: false}
	}

	var paths []string
	for _, m := range matches {
		paths = append(paths, m[1])
	}

	version := ""
	vm := versionRe.FindStringSubmatch(header)
	if vm != nil {
		version = vm[1]
	}

	return headerInfo{
		javaPath: paths[0],
		version:  version,
		valid:    true,
		multi:    len(paths) > 1,
		allPaths: paths,
	}
}

// -------- scoring --------

func computeScore(goFile string, content string, javaRoot string) FileResult {
	result := FileResult{GoFile: goFile}

	hdr := parseHeader(content)
	if !hdr.valid {
		result.Score = 0
		result.Notes = "missing Port of header"
		return result
	}

	// Version check
	versionPenalty := 0.0
	if hdr.version == "" {
		versionPenalty = 5
		result.Notes += "version missing; "
	} else if hdr.version != expectedVersion {
		versionPenalty = 5
		result.Notes += fmt.Sprintf("version mismatch (%s != %s); ", hdr.version, expectedVersion)
	}
	result.VersionPenalty = versionPenalty

	// Find Java source
	javaPath := filepath.Join(javaRoot, hdr.javaPath)
	result.JavaFile = hdr.javaPath

	javaContent, err := os.ReadFile(javaPath)
	if err != nil {
		// Try other paths if multi
		found := false
		for _, p := range hdr.allPaths[1:] {
			alt := filepath.Join(javaRoot, p)
			if jc, e := os.ReadFile(alt); e == nil {
				javaContent = jc
				result.JavaFile = p
				found = true
				break
			}
		}
		if !found {
			result.Score = 0
			result.Notes += fmt.Sprintf("java source not found: %s", hdr.javaPath)
			return result
		}
	}

	javaSrc := string(javaContent)

	// Extract Java methods from the primary source.
	javaMethods := extractJavaMethods(javaSrc)
	javaMethodCount := len(javaMethods)
	result.JavaMethodCount = javaMethodCount
	javaBodyLoc := javaBodyLoC(javaSrc)
	result.JavaBodyLoC = javaBodyLoc

	// For tiny type wrappers (≤50 LoC Go source), also load the full ancestor
	// class hierarchy so we can credit @Override coverage against the inherited
	// surface, not just the (often empty) override list in the subclass file.
	// We walk up to 4 levels of the extends chain.
	goLineCount := strings.Count(content, "\n")
	parentMethods := make(map[string]bool)
	if goLineCount <= 50 {
		javaDir := filepath.Dir(filepath.Join(javaRoot, hdr.javaPath))
		currentSrc := javaSrc
		for depth := 0; depth < 4; depth++ {
			parentClass := extractJavaParentClass(currentSrc)
			if parentClass == "" {
				break
			}
			parentFile := filepath.Join(javaDir, capitalise(parentClass)+".java")
			parentSrc, e := os.ReadFile(parentFile)
			if e != nil {
				break
			}
			for name := range extractJavaMethods(string(parentSrc)) {
				parentMethods[name] = true
			}
			currentSrc = string(parentSrc)
		}
	}

	// Extract Go functions.
	goFuncs, goBodyLoc, err := extractGoFunctions(goFile)
	if err != nil {
		result.Score = 0
		result.Notes += fmt.Sprintf("parse error: %v", err)
		return result
	}
	result.GoBodyLoC = goBodyLoc

	// Build Go method name set.  Functions explicitly tagged "// Go-only:" are
	// intentional divergences and are excluded from the penalty count but still
	// included in the method-count denominator.
	goMethodNames := make(map[string]bool)
	goMethodCount := 0
	stubCount := 0
	goOnlyCount := 0

	javaSrcLower := strings.ToLower(javaSrc)
	for _, fn := range goFuncs {
		goMethodCount++
		nameLower := strings.ToLower(fn.Name)
		goMethodNames[nameLower] = true
		if fn.IsStub {
			stubCount++
		}
		// A function is a "Go-only" divergence when:
		//   (a) it carries an explicit "// Go-only:" comment, or
		//   (b) its lowercase name does not appear anywhere in the Java source.
		// Explicitly tagged functions are still counted but do NOT incur a penalty.
		if fn.IsGoOnly {
			// Intentional divergence — don't penalise.
		} else if !strings.Contains(javaSrcLower, nameLower) {
			// Name missing from Java source entirely.
			goOnlyCount++
		}
	}

	result.GoMethodCount = goMethodCount
	result.StubCount = stubCount
	result.GoOnlyCount = goOnlyCount

	// Metric 1: Method name overlap (Jaccard).
	// For tiny wrappers we union the primary Java method set with the parent's
	// set so that correct @Override coverage is rewarded.
	javaSet := make(map[string]bool)
	for name := range javaMethods {
		javaSet[name] = true
	}
	// Augment with parent methods for tiny wrappers.
	augmentedJavaSet := make(map[string]bool)
	for name := range javaSet {
		augmentedJavaSet[name] = true
	}
	for name := range parentMethods {
		augmentedJavaSet[name] = true
	}

	intersection := 0
	for name := range goMethodNames {
		if augmentedJavaSet[name] {
			intersection++
		}
	}
	union := len(goMethodNames) + len(augmentedJavaSet) - intersection
	jaccard := 0.0
	if union > 0 {
		jaccard = float64(intersection) / float64(union)
	}
	result.MethodOverlapPct = jaccard * 100

	// Tiny-wrapper bonus: when the Go file is very small (≤50 LoC) and correctly
	// names a parent class that exists in our Java cache, reward the structural
	// fidelity even when method names diverge (Go uses interface methods instead
	// of @Override names).  Cap the bonus at 30 points so it cannot manufacture
	// a perfect score on its own.
	tinyWrapperBonus := 0.0
	if goLineCount <= 50 && len(parentMethods) > 0 {
		// Check how many Go method names appear in the combined parent+primary set.
		parentHits := 0
		for name := range goMethodNames {
			if parentMethods[name] || augmentedJavaSet[name] {
				parentHits++
			}
		}
		// Give a bonus proportional to how much of the Go surface is covered.
		if goMethodCount > 0 {
			tinyWrapperBonus = math.Min(30, float64(parentHits)/float64(goMethodCount)*30)
		} else {
			// File has no Go functions but inherits from the right parent — modest bonus.
			tinyWrapperBonus = 15
		}
		result.Notes += fmt.Sprintf("tiny-wrapper(parent bonus %.0f); ", tinyWrapperBonus)
	}

	// Metric 2: Method count ratio.
	countRatio := 0.0
	if goMethodCount > 0 && javaMethodCount > 0 {
		mn := math.Min(float64(goMethodCount), float64(javaMethodCount))
		mx := math.Max(float64(goMethodCount), float64(javaMethodCount))
		countRatio = mn / mx
	} else if goMethodCount == 0 && javaMethodCount == 0 {
		countRatio = 1.0
	}

	// Metric 3: Signature similarity (arity comparison for shared methods).
	sigSimilarity := 0.0
	{
		sharedCount := 0
		totalArityRatio := 0.0
		// Build map of Go func name -> param count.
		goParamMap := make(map[string]int)
		for _, fn := range goFuncs {
			goParamMap[strings.ToLower(fn.Name)] = fn.Params
		}
		// Build approximate Java arity from method regex.
		javaArityMap := extractJavaArities(javaSrc)

		for name := range goMethodNames {
			if augmentedJavaSet[name] {
				sharedCount++
				goArity := goParamMap[name]
				javaArity := javaArityMap[name]
				// Compute arity ratio.
				if goArity == 0 && javaArity == 0 {
					totalArityRatio += 1.0
				} else {
					mn := math.Min(float64(goArity), float64(javaArity))
					mx := math.Max(float64(goArity), float64(javaArity))
					if mx > 0 {
						totalArityRatio += mn / mx
					}
				}
			}
		}
		if sharedCount > 0 {
			sigSimilarity = totalArityRatio / float64(sharedCount)
		} else if goMethodCount == 0 {
			sigSimilarity = 1.0
		}
	}

	// Metric 6: Body LoC ratio.
	locRatio := 0.0
	locScore := 0.0
	if javaBodyLoc > 0 {
		locRatio = float64(goBodyLoc) / float64(javaBodyLoc)
		if locRatio > 1.5 {
			locRatio = 1.5
		}
		// Score = 100 - |1 - ratio| * 50
		locScore = 100 - math.Abs(1.0-locRatio)*50
		if locScore < 0 {
			locScore = 0
		}
	} else if goBodyLoc == 0 {
		locScore = 100
		locRatio = 1.0
	} else {
		// Java has no body LoC but Go does.
		locRatio = 1.5
		locScore = 50 // half penalty
	}
	result.LocRatio = locRatio

	// Weighted score: (jaccard*25 + countRatio*15 + sigSimilarity*10 + locScore*20) / 70
	// Normalise each metric to [0,100] scale.
	weightedSum := jaccard*100*25 + countRatio*100*15 + sigSimilarity*100*10 + locScore*20
	weightedScore := weightedSum / 70.0
	result.WeightedScore = weightedScore

	// Penalties.
	goOnlyPenalty := math.Min(20, float64(2*goOnlyCount))
	stubPenalty := math.Min(30, float64(3*stubCount))
	result.GoOnlyPenalty = goOnlyPenalty
	result.StubPenalty = stubPenalty

	// Final score = weightedScore * 0.7 + 30 + tinyWrapperBonus - penalties - versionPenalty
	score := weightedScore*0.7 + 30 + tinyWrapperBonus - goOnlyPenalty - stubPenalty - versionPenalty

	// Clamp
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	result.Score = score

	// Build notes
	var noteParts []string
	if result.Notes != "" {
		noteParts = append(noteParts, strings.TrimRight(result.Notes, "; "))
	}
	if stubCount > 0 {
		noteParts = append(noteParts, fmt.Sprintf("%d stubs", stubCount))
	}
	if goOnlyCount > 0 {
		noteParts = append(noteParts, fmt.Sprintf("%d go-only fns", goOnlyCount))
	}
	if hdr.multi {
		noteParts = append(noteParts, "multi-java-source")
	}
	result.Notes = strings.Join(noteParts, "; ")

	return result
}

// capitalise upper-cases the first rune of s.  Used to derive a Java
// class filename from a lowercase parent class name.
func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// extractJavaArities returns method name → param count for Java methods.
var javaMethodArityRe = regexp.MustCompile(
	`(?m)^[ \t]*(?:(?:public|private|protected|static|final|synchronized|abstract|native|default|override|@\w+(?:\([^)]*\))?)\s+)*` +
		`(?:<[^>]+>\s*)?` +
		`(?:[\w<>\[\],\s]+?\s+)` +
		`(\w+)\s*\(([^)]*)\)`)

func extractJavaArities(src string) map[string]int {
	stripped := stripJavaAnnotationLines(src)
	result := make(map[string]int)
	matches := javaMethodArityRe.FindAllStringSubmatch(stripped, -1)
	for _, m := range matches {
		name := strings.ToLower(m[1])
		if isJavaKeyword(name) {
			continue
		}
		params := m[2]
		params = strings.TrimSpace(params)
		if params == "" {
			result[name] = 0
		} else {
			// Count comma-separated params
			result[name] = strings.Count(params, ",")+1
		}
	}
	return result
}

// -------- walker --------

func walkInternal(internalDir string, javaRoot string) ([]FileResult, error) {
	var results []FileResult

	err := filepath.Walk(internalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip
		}
		if info.IsDir() {
			// Skip skel/ directories — these contain Java-like skeleton templates
			// that intentionally don't compile under normal build constraints.
			if info.Name() == "skel" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Quick check: does it have a Port of header?
		src := string(content)
		if !strings.Contains(src, "// Port of com/google/javascript") {
			return nil
		}

		// Skip files with //go:build skeletons tag — they're intentional skeletons.
		if strings.Contains(src, "//go:build skeletons") {
			return nil
		}

		result := computeScore(path, src, javaRoot)
		results = append(results, result)
		return nil
	})

	return results, err
}

// -------- report generation --------

const mdTemplate = `# Port Fidelity Report

Generated: auto-computed by tools/portcheck.

## Summary

| Metric | Value |
|--------|-------|
| Total files scored | {{.TotalFiles}} |
| Average score | {{printf "%.1f" .AvgScore}} |
| Std deviation | {{printf "%.1f" .StdDev}} |
| Min score | {{printf "%.1f" .MinScore}} |
| Median score | {{printf "%.1f" .MedianScore}} |
| Max score | {{printf "%.1f" .MaxScore}} |

## Files by Fidelity Score (worst first)

| Rank | Score | Go File | Java Source | Method Overlap % | Stub Count | LoC Ratio | Notes |
|------|-------|---------|-------------|-----------------|------------|-----------|-------|
{{range .Rows}}| {{.Rank}} | {{printf "%.0f" .Score}} | ` + "`" + `{{.GoFileShort}}` + "`" + ` | ` + "`" + `{{.JavaFileShort}}` + "`" + ` | {{printf "%.0f" .MethodOverlapPct}}% | {{.StubCount}} | {{printf "%.2f" .LocRatio}} | {{.Notes}} |
{{end}}
`

type mdRow struct {
	Rank            int
	Score           float64
	GoFileShort     string
	JavaFileShort   string
	MethodOverlapPct float64
	StubCount       int
	LocRatio        float64
	Notes           string
}

type mdData struct {
	TotalFiles  int
	AvgScore    float64
	StdDev      float64
	MinScore    float64
	MedianScore float64
	MaxScore    float64
	Rows        []mdRow
}

func computeStats(results []FileResult) mdData {
	if len(results) == 0 {
		return mdData{}
	}

	scores := make([]float64, len(results))
	for i, r := range results {
		scores[i] = r.Score
	}

	// Sort scores for median
	sorted := make([]float64, len(scores))
	copy(sorted, scores)
	sort.Float64s(sorted)

	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	avg := sum / float64(len(scores))

	variance := 0.0
	for _, s := range scores {
		d := s - avg
		variance += d * d
	}
	variance /= float64(len(scores))

	median := sorted[len(sorted)/2]
	if len(sorted)%2 == 0 {
		median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	}

	return mdData{
		TotalFiles:  len(results),
		AvgScore:    avg,
		StdDev:      math.Sqrt(variance),
		MinScore:    sorted[0],
		MedianScore: median,
		MaxScore:    sorted[len(sorted)-1],
	}
}

func generateReport(results []FileResult, internalDir string, outputMD string, outputJSON string) error {
	// Sort ascending by score (worst first)
	sorted := make([]FileResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Score != sorted[j].Score {
			return sorted[i].Score < sorted[j].Score
		}
		return sorted[i].GoFile < sorted[j].GoFile
	})

	stats := computeStats(results)

	// Build rows
	rows := make([]mdRow, len(sorted))
	for i, r := range sorted {
		goShort := strings.TrimPrefix(r.GoFile, internalDir+"/")
		if goShort == r.GoFile {
			// try relative from parent
			goShort = filepath.Base(r.GoFile)
		}
		javaShort := r.JavaFile
		// Trim com/google/javascript/ prefix for brevity
		javaShort = strings.TrimPrefix(javaShort, "com/google/javascript/")

		rows[i] = mdRow{
			Rank:            i + 1,
			Score:           r.Score,
			GoFileShort:     goShort,
			JavaFileShort:   javaShort,
			MethodOverlapPct: r.MethodOverlapPct,
			StubCount:       r.StubCount,
			LocRatio:        r.LocRatio,
			Notes:           r.Notes,
		}
	}
	stats.Rows = rows

	// Write Markdown
	tmpl, err := template.New("report").Parse(mdTemplate)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	mdFile, err := os.Create(outputMD)
	if err != nil {
		return fmt.Errorf("create md: %w", err)
	}
	defer mdFile.Close()

	if err := tmpl.Execute(mdFile, stats); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	// Write JSON
	jsonData, err := json.MarshalIndent(sorted, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	if err := os.WriteFile(outputJSON, jsonData, 0644); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}

// -------- main --------

func main() {
	// Locate repo root relative to this binary's path.
	// When run as `go run ./tools/portcheck/...` from repo root, use CWD.
	repoRoot := findRepoRoot()
	if repoRoot == "" {
		fmt.Fprintf(os.Stderr, "error: could not find repo root (no go.mod found)\n")
		os.Exit(1)
	}

	internalDir := filepath.Join(repoRoot, "internal")
	javaRoot := javaSourceRoot

	fmt.Fprintf(os.Stderr, "Scanning %s...\n", internalDir)
	fmt.Fprintf(os.Stderr, "Java source root: %s\n", javaRoot)

	results, err := walkInternal(internalDir, javaRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking internal/: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Scored %d files.\n", len(results))

	docsDir := filepath.Join(repoRoot, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating docs/: %v\n", err)
		os.Exit(1)
	}

	outputMD := filepath.Join(docsDir, "port-fidelity.md")
	outputJSON := filepath.Join(docsDir, "port-fidelity.json")

	if err := generateReport(results, internalDir, outputMD, outputJSON); err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report written to %s\n", outputMD)
	fmt.Printf("JSON written to %s\n", outputJSON)

	// Print summary
	stats := computeStats(results)
	fmt.Printf("Total files: %d, avg score: %.1f, min: %.1f, max: %.1f\n",
		stats.TotalFiles, stats.AvgScore, stats.MinScore, stats.MaxScore)
}

// findRepoRoot walks up from the current directory looking for go.mod.
func findRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
