package jscomp

// Differential test against tools/oracle/java/TranspileDump.java: the ES
// transpilation passes (TranspilationPasses.addTranspilationPasses, after the
// runtime libraries and Normalize, in DefaultPassConfig's transpile-only
// order) run on parsed files for --language_out ECMASCRIPT5, ECMASCRIPT_2015
// and ECMASCRIPT_2017, printed in the same format and compared file by file.
//
//   - TestTranspileGolden replays the checked-in golden dump
//     testdata/transpile/golden.txt.gz: the hand-written inputs
//     testdata/transpile/js/*.js with --full trees, followed by the hashes of
//     the corpus files of testdata/transpile/golden-files.txt (skipped when
//     the corpus is not present). Regenerate with
//       CP=.cache/closure-compiler.jar:<TranspileDump build dir>
//       J="java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp $CP com.google.javascript.jscomp.TranspileDump"
//       $J --full a.gz internal/jscomp/testdata/transpile/js <the *.js names>
//       $J b.gz .cache/corpus --list=internal/jscomp/testdata/transpile/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/transpile/golden.txt.gz
//   - TestTranspileOracle compares any other dump (hashes, or --full):
//     TRANSPILE_ORACLE=<dump.gz> [TRANSPILE_ROOT=<root>] (default root: the
//     corpus under /home/user/jsgo/.cache/corpus), optionally
//     TRANSPILE_DIFF_DIR=<dir> to write the differing sections there. The
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

// transpileOracleVariants is TranspileDump.variants (its default).
var transpileOracleVariants = []string{"es5", "es2015", "es2017"}

// transpileOracleEscape is TranspileDump.escape (ParsingDump's escaping).
func transpileOracleEscape(s string) string {
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

// transpileOracleException is TranspileDump.exception:
// getClass().getName() + ": " + getMessage().
func transpileOracleException(r any) string {
	t := jutil.AsThrowable(r)
	full := t.Error()
	msg := t.GetMessage()
	if strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg) + ": " + msg
	}
	// No message: Java prints "null".
	return full + ": null"
}

func transpileOracleSha16(lines []string) string {
	h := sha256.New()
	for _, l := range lines {
		io.WriteString(h, l+"\n")
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

type transpileOracleDumper struct {
	full     bool
	variants []string
	sb       strings.Builder
}

func (d *transpileOracleDumper) line(s string) {
	d.sb.WriteString(transpileOracleEscape(s))
	d.sb.WriteByte('\n')
}

// newCompiler mirrors TranspileDump.newCompiler.
func (d *transpileOracleDumper) newCompiler(rel string, src string, variant string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	switch variant {
	case "es5":
		options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	case "es2015":
		options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	case "es2017":
		options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2017)
	case "parse":
		options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("unknown variant " + variant))
	}
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	compiler.Init(nil, []*SourceFile{SourceFile_FromCode(rel, src)}, options)
	compiler.Parse()
	return compiler
}

func transpileOracleDiagnostics(compiler *Compiler) []string {
	var out []string
	for _, e := range compiler.GetErrors() {
		out = append(out, transpileOracleDiag("ERROR", e))
	}
	for _, e := range compiler.GetWarnings() {
		out = append(out, transpileOracleDiag("WARNING", e))
	}
	return out
}

func transpileOracleDiag(level string, e *JSError) string {
	return "diag " + level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " \"" + e.GetDescription() + "\""
}

func transpileOracleCollectChanges(n *rhino.Node, out *[]string) {
	if (n.IsScript() || n.IsFunction()) && n.GetChangeTime() != 0 {
		*out = append(*out, n.GetToken().String()+"@"+strconv.Itoa(n.GetLineno())+":"+
			strconv.Itoa(n.GetCharno())+" "+strconv.Itoa(int(n.GetChangeTime())))
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		transpileOracleCollectChanges(c, out)
	}
}

// transpileOracleRemoveFirst is List.remove(Object).
func transpileOracleRemoveFirst(list []string, s string) []string {
	for i, x := range list {
		if x == s {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// passes mirrors TranspileDump.passes.
func transpileOraclePasses(compiler *Compiler) []*PassFactory {
	options := compiler.GetOptions()
	passes := NewPassListBuilder(options)
	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName(PassNames_NORMALIZE).
			SetInternalFactory(func(c AbstractCompiler) CompilerPass { return Normalize_CreateNormalizeForOptimizations(c) }).
			Build())
	TranspilationPasses_AddTranspilationPasses(passes, options)
	return passes.Build()
}

// dumpFile mirrors TranspileDump.dumpFile.
func (d *transpileOracleDumper) dumpFile(rel string, src string) {
	d.line("## file " + rel)
	compiler := d.newCompiler(rel, src, "parse")
	parseDiags := transpileOracleDiagnostics(compiler)
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
	for _, variant := range d.variants {
		d.line("## variant " + variant)
		c := d.newCompiler(rel, src, variant)
		baseDiags := transpileOracleDiagnostics(c)
		before := c.GetRoot().ToStringTree()
		exception := ""
		count := 0
		func() {
			defer func() {
				if r := recover(); r != nil {
					exception = "exception " + transpileOracleException(r)
				}
			}()
			for _, factory := range transpileOraclePasses(c) {
				count++
				factory.Create(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			}
		}()
		diags := transpileOracleDiagnostics(c)
		for _, diag := range baseDiags {
			diags = transpileOracleRemoveFirst(diags, diag)
		}
		for _, diag := range diags {
			d.line(diag)
		}
		if exception != "" {
			d.line(exception)
			d.line("passes " + strconv.Itoa(count))
			continue
		}
		d.line("passes " + strconv.Itoa(count))
		after := c.GetRoot().ToStringTree()
		if after == before {
			d.line("tree unchanged")
		} else {
			var lines []string
			for _, l := range strings.Split(after, "\n") {
				lines = append(lines, transpileOracleEscape(l))
			}
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			d.line("tree changed " + transpileOracleSha16(lines))
			if d.full {
				for _, l := range lines {
					d.sb.WriteString(l + "\n")
				}
			}
		}
		var changes []string
		transpileOracleCollectChanges(c.GetRoot(), &changes)
		var escapedChanges []string
		for _, l := range changes {
			escapedChanges = append(escapedChanges, transpileOracleEscape(l))
		}
		d.line("changes " + strconv.Itoa(len(changes)) + " " + transpileOracleSha16(escapedChanges))
		if d.full {
			for _, l := range escapedChanges {
				d.sb.WriteString(l + "\n")
			}
		}
		d.line("libs " + strings.Join(c.GetRuntimeJsLibManager().GetInjectedLibraries(), ","))
		for script := c.GetJsRoot().GetFirstChild(); script != nil; script = script.GetNext() {
			d.line("features " + jutil.ToJavaString(script.GetSourceFileName()) + " " +
				jutil.ToJavaString(NodeUtil_GetFeatureSetOfScript(script)))
		}
		d.line("allowable " + c.GetAllowableFeatures().String())
		d.line("stage " + c.GetLifeCycleStage().String())
	}
}

// transpileOracleSplit splits a dump into its "## file" sections, in order,
// and reports which sections have --full trees and which variants it ran.
func transpileOracleSplit(t *testing.T, path string) (names []string, sections map[string]string, full map[string]bool, variants []string) {
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
	full = map[string]bool{}
	seen := map[string]bool{}
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
		if v, ok := strings.CutPrefix(l, "## variant "); ok && !seen[v] {
			seen[v] = true
			variants = append(variants, v)
		}
		if strings.HasPrefix(l, "    ") || strings.HasPrefix(l, "ROOT") {
			full[cur] = true
		}
		sb.WriteString(l + "\n")
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	if len(variants) == 0 {
		variants = transpileOracleVariants
	}
	return names, sections, full, variants
}

// transpileOracleUnescapeName undoes the escaping of a file name.
func transpileOracleUnescapeName(s string) string {
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

// transpileOracleCompare dumps every file of the dump at path (sources under
// root) and compares; it returns the number of files compared.
func transpileOracleCompare(t *testing.T, path string, root string, onlyNames func(string) bool) int {
	t.Helper()
	// The dump was made with -XX:-ShowCodeDetailsInExceptionMessages: a
	// NullPointerException that the JVM throws has no message.
	defer jutil.SetShowCodeDetailsInExceptionMessages(jutil.SetShowCodeDetailsInExceptionMessages(false))
	names, sections, full, variants := transpileOracleSplit(t, path)
	diffDir := os.Getenv("TRANSPILE_DIFF_DIR")
	n, failures := 0, 0
	for _, name := range names {
		rel := transpileOracleUnescapeName(name)
		if onlyNames != nil && !onlyNames(rel) {
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Logf("skipping %s: %v", rel, err)
			continue
		}
		d := &transpileOracleDumper{full: full[name], variants: variants}
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.line("go panic " + transpileOracleException(r))
				}
			}()
			d.dumpFile(rel, jutil.NewStringFromUTF8Bytes(src))
		}()
		n++
		if got, want := d.sb.String(), sections[name]; got != want {
			failures++
			if failures <= 10 {
				t.Errorf("%s: dump differs:\n%s", rel, transpileOracleFirstDiff(got, want))
			}
			if diffDir != "" {
				base := strings.ReplaceAll(rel, "/", "__")
				_ = os.WriteFile(filepath.Join(diffDir, base+".go.txt"), []byte(got), 0o644)
				_ = os.WriteFile(filepath.Join(diffDir, base+".java.txt"), []byte(want), 0o644)
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, n)
	}
	return n
}

func transpileOracleFirstDiff(got, want string) string {
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

func TestTranspileGolden(t *testing.T) {
	golden := "testdata/transpile/golden.txt.gz"
	jsDir := "testdata/transpile/js"
	hand := map[string]bool{}
	entries, err := os.ReadDir(jsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		hand[e.Name()] = true
	}
	if n := transpileOracleCompare(t, golden, jsDir, func(rel string) bool { return hand[rel] }); n != len(hand) {
		t.Errorf("compared %d hand-written files, want %d", n, len(hand))
	}
	corpus := "/home/user/jsgo/.cache/corpus"
	if _, err := os.Stat(corpus); err != nil {
		t.Logf("corpus not present; golden corpus files skipped")
		return
	}
	transpileOracleCompare(t, golden, corpus, func(rel string) bool { return !hand[rel] })
}

func TestTranspileOracle(t *testing.T) {
	dump := os.Getenv("TRANSPILE_ORACLE")
	if dump == "" {
		t.Skip("set TRANSPILE_ORACLE=<TranspileDump output .gz> (and TRANSPILE_ROOT) to run")
	}
	root := os.Getenv("TRANSPILE_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	n := transpileOracleCompare(t, dump, root, nil)
	t.Logf("compared %d files", n)
}
