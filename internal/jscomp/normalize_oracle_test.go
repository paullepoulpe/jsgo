package jscomp

// Differential test against tools/oracle/java/NormalizeDump.java: Normalize
// run the way CompilerTestCase.enableNormalize() runs it, on parsed files,
// printed in the same format and compared file by file.
//
//   - TestNormalizeGolden replays the checked-in golden dump
//     testdata/normalize/golden.txt.gz (with --full trees): the hand-written
//     inputs testdata/normalize/js/*.js, followed by the corpus files of
//     testdata/normalize/golden-files.txt (skipped when the corpus is not
//     present). Regenerate with
//       CP=.cache/$CLOSURE_JAR:<NormalizeDump build dir>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.NormalizeDump --full a.gz \
//         internal/jscomp/testdata/normalize/js <the *.js names>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.NormalizeDump --full b.gz \
//         .cache/corpus --list=internal/jscomp/testdata/normalize/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/normalize/golden.txt.gz
//   - TestNormalizeOracle compares any other dump (hashes, or --full):
//     NORMALIZE_ORACLE=<dump.gz> [NORMALIZE_ROOT=<root>] (default root: the
//     corpus under /home/user/jsgo/.cache/corpus). The held-out corpus
//     (.cache/corpus-holdout, corpus-holdout-files.txt) is for a final check
//     only.

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

const normalizeOracleExterns = "var window; var Arguments;"

// normalizeOracleEscape is NormalizeDump.escape (ParsingDump's escaping).
func normalizeOracleEscape(s string) string {
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

// normalizeOracleException is NormalizeDump.exception:
// getClass().getName() + ": " + getMessage().
func normalizeOracleException(r any) string {
	t := jutil.AsThrowable(r)
	full := t.Error()
	msg := t.GetMessage()
	if strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg) + ": " + msg
	}
	// No message: Java prints "null".
	return full + ": null"
}

func normalizeOracleSha16(lines []string) string {
	h := sha256.New()
	for _, l := range lines {
		io.WriteString(h, l+"\n")
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

type normalizeOracleDumper struct {
	full bool
	sb   strings.Builder
}

func (d *normalizeOracleDumper) line(s string) {
	d.sb.WriteString(normalizeOracleEscape(s))
	d.sb.WriteByte('\n')
}

// newCompiler mirrors NormalizeDump.newCompiler.
func (d *normalizeOracleDumper) newCompiler(rel string, src string, google bool) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	// As NormalizeTest.getOptions does.
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	if google {
		options.SetCodingConvention(NewGoogleCodingConvention())
	}
	var externs []*SourceFile
	if google {
		externs = []*SourceFile{SourceFile_FromCode("externs.js", normalizeOracleExterns)}
	}
	compiler.Init(externs, []*SourceFile{SourceFile_FromCode(rel, src)}, options)
	compiler.Parse()
	return compiler
}

func normalizeOracleDiagnostics(compiler *Compiler) []string {
	var out []string
	for _, e := range compiler.GetErrors() {
		out = append(out, normalizeOracleDiag("ERROR", e))
	}
	for _, e := range compiler.GetWarnings() {
		out = append(out, normalizeOracleDiag("WARNING", e))
	}
	return out
}

func normalizeOracleDiag(level string, e *JSError) string {
	return "diag " + level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " \"" + e.GetDescription() + "\""
}

func normalizeOracleCollectChanges(n *rhino.Node, out *[]string) {
	if (n.IsScript() || n.IsFunction()) && n.GetChangeTime() != 0 {
		*out = append(*out, n.GetToken().String()+"@"+strconv.Itoa(n.GetLineno())+":"+
			strconv.Itoa(n.GetCharno())+" "+strconv.Itoa(int(n.GetChangeTime())))
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		normalizeOracleCollectChanges(c, out)
	}
}

// removeFirst is List.remove(Object).
func normalizeOracleRemoveFirst(list []string, s string) []string {
	for i, x := range list {
		if x == s {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// dumpFile mirrors NormalizeDump.dumpFile.
func (d *normalizeOracleDumper) dumpFile(rel string, src string) {
	d.line("## file " + rel)
	compiler := d.newCompiler(rel, src, false)
	parseDiags := normalizeOracleDiagnostics(compiler)
	if compiler.HasErrors() {
		d.line("# parse errors")
	} else {
		d.line("# parse ok")
	}
	for _, diag := range parseDiags {
		d.line(diag)
	}
	if compiler.HasErrors() {
		return
	}
	for _, google := range []bool{true, false} {
		if google {
			d.line("## variant google")
		} else {
			d.line("## variant closure")
		}
		c := d.newCompiler(rel, src, google)
		baseDiags := normalizeOracleDiagnostics(c)
		before := c.GetRoot().ToStringTree()
		exception := ""
		func() {
			defer func() {
				if r := recover(); r != nil {
					exception = "exception " + normalizeOracleException(r)
				}
			}()
			Normalize_CreateNormalizeForOptimizations(c).Process(c.GetExternsRoot(), c.GetJsRoot())
		}()
		diags := normalizeOracleDiagnostics(c)
		for _, diag := range baseDiags {
			diags = normalizeOracleRemoveFirst(diags, diag)
		}
		for _, diag := range diags {
			d.line(diag)
		}
		if exception != "" {
			d.line(exception)
			continue
		}
		after := c.GetRoot().ToStringTree()
		if after == before {
			d.line("tree unchanged")
		} else {
			var lines []string
			for _, l := range strings.Split(after, "\n") {
				lines = append(lines, normalizeOracleEscape(l))
			}
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			d.line("tree changed " + normalizeOracleSha16(lines))
			if d.full {
				for _, l := range lines {
					d.sb.WriteString(l + "\n")
				}
			}
		}
		var changes []string
		normalizeOracleCollectChanges(c.GetRoot(), &changes)
		var escapedChanges []string
		for _, l := range changes {
			escapedChanges = append(escapedChanges, normalizeOracleEscape(l))
		}
		d.line("changes " + strconv.Itoa(len(changes)) + " " + normalizeOracleSha16(escapedChanges))
		if d.full {
			for _, l := range escapedChanges {
				d.sb.WriteString(l + "\n")
			}
		}
		d.line("stage " + c.GetLifeCycleStage().String())
		recheck := "recheck ok"
		func() {
			defer func() {
				if r := recover(); r != nil {
					recheck = "recheck exception " + normalizeOracleException(r)
				}
			}()
			Normalize_Builder_(c).AssertOnChange(true).Build().Process(c.GetExternsRoot(), c.GetJsRoot())
		}()
		d.line(recheck)
	}
}

// normalizeOracleSplit splits a dump into its "## file" sections, in order.
func normalizeOracleSplit(t *testing.T, path string) (names []string, sections map[string]string, full bool) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	zr.Multistream(true)
	sections = map[string]string{}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	var cur string
	var sb strings.Builder
	flush := func() {
		if cur != "" {
			sections[cur] = sb.String()
			names = append(names, cur)
		}
		sb.Reset()
	}
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "## file ") {
			flush()
			cur = strings.TrimPrefix(l, "## file ")
		}
		if strings.HasPrefix(l, "    ") || strings.HasPrefix(l, "ROOT") {
			full = true
		}
		sb.WriteString(l + "\n")
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	return names, sections, full
}

// normalizeOracleUnescapeName undoes the escaping of a file name.
func normalizeOracleUnescapeName(s string) string {
	if !strings.Contains(s, "\\") {
		return s
	}
	var units []uint16
	for i := 0; i < len(s); {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == '\\' {
			units = append(units, '\\')
			i += 2
		} else if s[i] == '\\' && i+5 < len(s) && s[i+1] == 'u' {
			v, _ := strconv.ParseUint(s[i+2:i+6], 16, 16)
			units = append(units, uint16(v))
			i += 6
		} else {
			units = append(units, uint16(s[i]))
			i++
		}
	}
	return jutil.StringFromUTF16(units)
}

// normalizeOracleCompare dumps every file of the dump at path (sources under
// root) and compares; it returns the number of files compared.
func normalizeOracleCompare(t *testing.T, path string, root string, onlyNames func(string) bool) int {
	t.Helper()
	names, sections, full := normalizeOracleSplit(t, path)
	n, failures := 0, 0
	for _, name := range names {
		rel := normalizeOracleUnescapeName(name)
		if onlyNames != nil && !onlyNames(rel) {
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Logf("skipping %s: %v", rel, err)
			continue
		}
		d := &normalizeOracleDumper{full: full}
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.line("go panic " + normalizeOracleException(r))
				}
			}()
			d.dumpFile(rel, jutil.NewStringFromUTF8Bytes(src))
		}()
		n++
		if got, want := d.sb.String(), sections[name]; got != want {
			failures++
			if failures <= 10 {
				t.Errorf("%s: dump differs:\n%s", rel, normalizeOracleFirstDiff(got, want))
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, n)
	}
	return n
}

func normalizeOracleFirstDiff(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")
	for i := 0; i < len(g) || i < len(w); i++ {
		var gl, wl string
		if i < len(g) {
			gl = g[i]
		}
		if i < len(w) {
			wl = w[i]
		}
		if gl != wl {
			return "line " + strconv.Itoa(i+1) + ":\n  go:   " + gl + "\n  java: " + wl
		}
	}
	return "(no line differs)"
}

func TestNormalizeGolden(t *testing.T) {
	golden := "testdata/normalize/golden.txt.gz"
	jsDir := "testdata/normalize/js"
	hand := map[string]bool{}
	entries, err := os.ReadDir(jsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		hand[e.Name()] = true
	}
	if n := normalizeOracleCompare(t, golden, jsDir, func(rel string) bool { return hand[rel] }); n != len(hand) {
		t.Errorf("compared %d hand-written files, want %d", n, len(hand))
	}
	corpus := "/home/user/jsgo/.cache/corpus"
	if _, err := os.Stat(corpus); err != nil {
		t.Logf("corpus not present; golden corpus files skipped")
		return
	}
	normalizeOracleCompare(t, golden, corpus, func(rel string) bool { return !hand[rel] })
}

func TestNormalizeOracle(t *testing.T) {
	dump := os.Getenv("NORMALIZE_ORACLE")
	if dump == "" {
		t.Skip("set NORMALIZE_ORACLE=<NormalizeDump output .gz> (and NORMALIZE_ROOT) to run")
	}
	root := os.Getenv("NORMALIZE_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	n := normalizeOracleCompare(t, dump, root, nil)
	t.Logf("compared %d files", n)
}
