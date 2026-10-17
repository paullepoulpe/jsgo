package jscomp

// Go-only: test-helper — differential test against
// tools/oracle/java/PeepholeDump.java: every peephole pass run alone (the way
// its unit test's getProcessor runs it), and DefaultPassConfig's peephole
// sequences, on the normalized AST of each input, printed in the same format
// and compared block by block.
//
//   - TestPeepholeGolden replays the checked-in golden dump
//     testdata/peephole/golden.txt.gz (--full blocks): the hand-written
//     inputs testdata/peephole/js/*.js, then the corpus files of
//     testdata/peephole/golden-files.txt (skipped when the corpus is not
//     present). Regenerate with
//       CP=.cache/closure-compiler.jar:<PeepholeDump build dir>
//       J="java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp $CP com.google.javascript.jscomp.PeepholeDump"
//       $J --full a.gz internal/jscomp/testdata/peephole/js peephole-fold.js peephole-control.js peephole-methods.js
//       $J --full b.gz .cache/corpus @internal/jscomp/testdata/peephole/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/peephole/golden.txt.gz
//   - TestPeepholeOracle compares any other dump (hash lines or --full
//     blocks): PEEPHOLE_ORACLE=<dump.gz> [PEEPHOLE_ROOT=<root>] (default: the
//     corpus under /home/user/jsgo/.cache/corpus) [PEEPHOLE_EVERY=N: only every
//     Nth file of the dump] [PEEPHOLE_DIFF_DIR=<dir>: write the Go block of
//     every differing (file, config) there] [PEEPHOLE_MAX_DIFFS=N]. The
//     held-out corpus (.cache/corpus-holdout, corpus-holdout-files.txt) is for
//     a final check only.

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// peepholeOracleConfig is one configuration of PeepholeDump.CONFIGS.
type peepholeOracleConfig struct {
	name string
	pass func(c AbstractCompiler) CompilerPass // nil result: no pass
}

func peepholeOracleSingle(c AbstractCompiler, name string, optimization AbstractPeepholeOptimization) CompilerPass {
	return NewPeepholeOptimizationsPass(c, name, optimization)
}

// peepholeOracleConfigs is PeepholeDump.CONFIGS, in the same order.
var peepholeOracleConfigs = []peepholeOracleConfig{
	{"normalize", func(c AbstractCompiler) CompilerPass { return nil }},
	{"fold", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "fold", NewPeepholeFoldConstants(false, false))
	}},
	{"fold-late", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "fold-late", NewPeepholeFoldConstants(true, false))
	}},
	{"minimize", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "minimize", NewPeepholeMinimizeConditions(false))
	}},
	{"minimize-late", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "minimize-late", NewPeepholeMinimizeConditions(true))
	}},
	{"substitute", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "substitute", NewPeepholeSubstituteAlternateSyntax(false))
	}},
	{"substitute-late", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "substitute-late", NewPeepholeSubstituteAlternateSyntax(true))
	}},
	{"known-methods", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "known-methods", NewPeepholeReplaceKnownMethods(false, false))
	}},
	{"known-methods-late", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "known-methods-late", NewPeepholeReplaceKnownMethods(true, false))
	}},
	{"dead-code", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "dead-code", NewPeepholeRemoveDeadCode())
	}},
	{"collect-props", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "collect-props", NewPeepholeCollectPropertyAssignments())
	}},
	{"exit-points", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "exit-points", NewMinimizeExitPoints())
	}},
	{"fusion", func(c AbstractCompiler) CompilerPass {
		return peepholeOracleSingle(c, "fusion", NewStatementFusion())
	}},
	{"exploit-assigns", func(c AbstractCompiler) CompilerPass {
		return NewPeepholeOptimizationsPass(c, PassNames_POST_NORMALIZE_PEEPHOLE, NewExploitAssigns())
	}},
	{"dotted", func(c AbstractCompiler) CompilerPass { return NewConvertToDottedProperties(c) }},
	{"peephole", func(c AbstractCompiler) CompilerPass {
		// DefaultPassConfig.createPeepholeOptimizationsPass(compiler,
		// PEEPHOLE_OPTIMIZATIONS, true).
		late := false
		useTypes := c.GetOptions().ShouldUseTypesForLocalOptimization()
		var optimizations []AbstractPeepholeOptimization
		optimizations = append(optimizations, NewMinimizeExitPoints())
		optimizations = append(optimizations, NewPeepholeMinimizeConditions(late))
		optimizations = append(optimizations, NewPeepholeSubstituteAlternateSyntax(late))
		optimizations = append(optimizations, NewPeepholeReplaceKnownMethods(late, useTypes))
		optimizations = append(optimizations, NewPeepholeRemoveDeadCode())
		optimizations = append(optimizations, NewPeepholeFoldConstants(late, useTypes))
		optimizations = append(optimizations, NewPeepholeCollectPropertyAssignments())
		return NewPeepholeOptimizationsPassOfList(c, PassNames_PEEPHOLE_OPTIMIZATIONS, optimizations)
	}},
	{"peephole-late", func(c AbstractCompiler) CompilerPass {
		// DefaultPassConfig.latePeepholeOptimizations.
		late := true
		useTypes := c.GetOptions().ShouldUseTypesForLocalOptimization()
		return NewPeepholeOptimizationsPass(
			c,
			"latePeepholeOptimizations",
			NewStatementFusion(),
			NewPeepholeRemoveDeadCode(),
			NewPeepholeMinimizeConditions(late),
			NewPeepholeSubstituteAlternateSyntax(late),
			NewPeepholeReplaceKnownMethods(late, useTypes),
			NewPeepholeFoldConstants(late, useTypes))
	}},
}

func peepholeOracleConfigByName(name string) *peepholeOracleConfig {
	for i := range peepholeOracleConfigs {
		if peepholeOracleConfigs[i].name == name {
			return &peepholeOracleConfigs[i]
		}
	}
	return nil
}

// peepholeOracleEscape is PeepholeDump.escape (ParsingDump's escaping).
func peepholeOracleEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			const hexDigits = "0123456789abcdef"
			b.WriteString("\\u")
			b.WriteByte(hexDigits[c>>12&0xf])
			b.WriteByte(hexDigits[c>>8&0xf])
			b.WriteByte(hexDigits[c>>4&0xf])
			b.WriteByte(hexDigits[c&0xf])
		}
	}
	return b.String()
}

type peepholeOracleDumper struct {
	sb strings.Builder
}

func (d *peepholeOracleDumper) line(s string) {
	d.sb.WriteString(peepholeOracleEscape(s))
	d.sb.WriteByte('\n')
}

// peepholeOracleError is PeepholeDump.error.
func peepholeOracleError(level string, e *JSError) string {
	return level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " " + e.GetDescription()
}

func (d *peepholeOracleDumper) diagnostics(compiler *Compiler, skipWarnings int) {
	d.line("# diagnostics")
	for _, e := range compiler.GetErrors() {
		d.line(peepholeOracleError("ERROR", e))
	}
	warnings := compiler.GetWarnings()
	for i := skipWarnings; i < len(warnings); i++ {
		d.line(peepholeOracleError("WARNING", warnings[i]))
	}
}

func (d *peepholeOracleDumper) bigints(n *rhino.Node) {
	if n.IsBigInt() {
		d.line(strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno()) + " " + n.GetBigInt().String())
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		d.bigints(c)
	}
}

// peepholeOracleException is PeepholeDump's exception line:
// getClass().getName() + ": " + the first line of String.valueOf(getMessage()).
func peepholeOracleException(r any) string {
	t := jutil.AsThrowable(r)
	s := t.Error()
	if jutil.GetMessageOrNull(t) == nil {
		s += ": null"
	}
	if nl := strings.IndexByte(s, '\n'); nl >= 0 {
		s = s[:nl]
	}
	return s
}

// peepholeOracleDumpOne is PeepholeDump.dumpOne.
func peepholeOracleDumpOne(rel string, src string, config string) string {
	d := &peepholeOracleDumper{}
	d.line("=== " + rel + " " + config)
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", "")},
		[]*SourceFile{SourceFile_FromCode(rel, src)},
		options)
	root := compiler.ParseInputs()
	if root == nil || compiler.GetErrorCount() > 0 {
		d.line("# parse errors")
		for _, e := range compiler.GetErrors() {
			d.line(peepholeOracleError("ERROR", e))
		}
		return d.sb.String()
	}
	parseWarnings := compiler.GetWarningCount()
	externsRoot := root.GetFirstChild()
	mainRoot := root.GetLastChild()
	var exception any
	func() {
		defer func() {
			if r := recover(); r != nil {
				exception = r
			}
		}()
		// CompilerTestCase.updateAccessorSummary
		compiler.SetAccessorSummary(
			AccessorSummary_Create(GatherGetterAndSetterProperties_Gather(compiler, root)))
		// CompilerTestCase.normalizeActualCode
		Normalize_CreateNormalizeForOptimizations(compiler).Process(externsRoot, mainRoot)
		pass := peepholeOracleConfigByName(config).pass(compiler)
		if pass != nil {
			pass.Process(externsRoot, mainRoot)
		}
	}()
	if exception != nil {
		d.diagnostics(compiler, parseWarnings)
		d.line("# exception " + peepholeOracleException(exception))
		return d.sb.String()
	}
	d.diagnostics(compiler, parseWarnings)
	d.line("# tree")
	for _, l := range strings.Split(mainRoot.ToStringTree(), "\n") {
		d.line(l)
	}
	// Node.toString does not print BIGINT values.
	d.line("# bigints")
	d.bigints(mainRoot)
	return d.sb.String()
}

func peepholeOracleSha256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// peepholeOracleEntry is one block of a dump: its file, configuration and
// either the full text (--full) or its hash.
type peepholeOracleEntry struct {
	rel, config string
	full        string // the --full block, or ""
	hash        string
}

// peepholeOracleRead reads a PeepholeDump output (hash lines, or --full
// blocks starting with "=== <rel> <config>").
func peepholeOracleRead(t *testing.T, path string) []peepholeOracleEntry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var r io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		gz.Multistream(true)
		r = gz
	}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<28)
	var out []peepholeOracleEntry
	var cur *strings.Builder
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "=== ") {
			head := strings.TrimPrefix(l, "=== ")
			sp := strings.LastIndexByte(head, ' ')
			out = append(out, peepholeOracleEntry{rel: head[:sp], config: head[sp+1:]})
			cur = &strings.Builder{}
			cur.WriteString(l + "\n")
			out[len(out)-1].full = cur.String()
			continue
		}
		if cur != nil {
			cur.WriteString(l + "\n")
			out[len(out)-1].full = cur.String()
			continue
		}
		parts := strings.Split(l, "\t")
		if len(parts) != 3 {
			t.Fatalf("bad oracle line %q", l)
		}
		out = append(out, peepholeOracleEntry{rel: parts[0], config: parts[1], hash: parts[2]})
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	for i := range out {
		if out[i].full != "" {
			out[i].hash = peepholeOracleSha256(out[i].full)
		}
	}
	return out
}

// peepholeOracleCompare runs every entry whose file exists under one of the
// roots and reports the differences.
func peepholeOracleCompare(t *testing.T, entries []peepholeOracleEntry, roots []string, every int, diffDir string, maxDiffs int) {
	t.Helper()
	diffs, checked, missing := 0, 0, 0
	files := map[string]int{}
	srcCache := map[string]string{}
	for _, e := range entries {
		idx, seen := files[e.rel]
		if !seen {
			idx = len(files)
			files[e.rel] = idx
		}
		if every > 1 && idx%every != 0 {
			continue
		}
		src, ok := srcCache[e.rel]
		if !ok {
			found := false
			for _, root := range roots {
				b, err := os.ReadFile(filepath.Join(root, e.rel))
				if err == nil {
					src, found = string(b), true
					break
				}
			}
			if !found {
				missing++
				continue
			}
			srcCache = map[string]string{e.rel: src} // keep only the current file
		}
		checked++
		got := peepholeOracleDumpOne(e.rel, src, e.config)
		if peepholeOracleSha256(got) == e.hash {
			continue
		}
		diffs++
		if diffs <= maxDiffs {
			if e.full != "" {
				t.Errorf("%s %s differs:\n%s", e.rel, e.config, peepholeOracleFirstDiff(e.full, got))
			} else {
				t.Errorf("%s %s differs (hash)", e.rel, e.config)
			}
		}
		if diffDir != "" {
			name := strings.ReplaceAll(e.rel, "/", "__") + "." + e.config
			os.WriteFile(filepath.Join(diffDir, name+".go.txt"), []byte(got), 0o644)
			if e.full != "" {
				os.WriteFile(filepath.Join(diffDir, name+".java.txt"), []byte(e.full), 0o644)
			}
		}
	}
	t.Logf("%d blocks checked, %d differ, %d skipped (input not found)", checked, diffs, missing)
}

// peepholeOracleFirstDiff shows the first differing line with some context.
func peepholeOracleFirstDiff(want, got string) string {
	w := strings.Split(want, "\n")
	g := strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			lo := max(0, i-3)
			var sb strings.Builder
			for j := lo; j < i; j++ {
				sb.WriteString("    " + w[j] + "\n")
			}
			sb.WriteString("Java: " + wl + "\n")
			sb.WriteString("Go:   " + gl + "\n")
			return sb.String()
		}
	}
	return "(identical lines)"
}

func TestPeepholeGolden(t *testing.T) {
	entries := peepholeOracleRead(t, "testdata/peephole/golden.txt.gz")
	if len(entries) < 1000 {
		t.Fatalf("only %d golden blocks", len(entries))
	}
	corpus := "/home/user/jsgo/.cache/corpus"
	roots := []string{"testdata/peephole/js"}
	if _, err := os.Stat(corpus); err == nil {
		roots = append(roots, corpus)
	}
	peepholeOracleCompare(t, entries, roots, 1, os.Getenv("PEEPHOLE_DIFF_DIR"), 20)
}

func TestPeepholeOracle(t *testing.T) {
	path := os.Getenv("PEEPHOLE_ORACLE")
	if path == "" {
		t.Skip("set PEEPHOLE_ORACLE=<PeepholeDump output> to run")
	}
	root := os.Getenv("PEEPHOLE_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	every := 1
	if s := os.Getenv("PEEPHOLE_EVERY"); s != "" {
		every, _ = strconv.Atoi(s)
	}
	maxDiffs := 50
	if s := os.Getenv("PEEPHOLE_MAX_DIFFS"); s != "" {
		maxDiffs, _ = strconv.Atoi(s)
	}
	peepholeOracleCompare(t, peepholeOracleRead(t, path), []string{root}, every, os.Getenv("PEEPHOLE_DIFF_DIR"), maxDiffs)
}
