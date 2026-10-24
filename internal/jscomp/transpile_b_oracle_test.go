package jscomp

// Differential test against tools/oracle/java/TranspileBDump.java: the
// w5-transpile-b passes (Es6NormalizeClasses, Es6ConvertSuper,
// Es6RewriteClass + Es6ConvertSuperConstructorCalls, the block-scoping
// rewrites, RewriteNewDotTarget / ReportUntranspilableFeatures,
// RewritePolyfills + PolyfillUsageFinder, IsolatePolyfills), run in DefaultPassConfig's order around the rest of
// the transpilation passes, printed in the same format and compared file by
// file.
//
//   - TestTranspileBGolden replays the checked-in golden dump
//     testdata/transpileb/golden.txt.gz (with --full trees): the hand-written
//     inputs testdata/transpileb/js/*.js, followed by the corpus files of
//     testdata/transpileb/golden-files.txt (skipped when the corpus is not
//     present). Regenerate with
//       CP=.cache/$CLOSURE_JAR:<TranspileBDump build dir>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.TranspileBDump --full a.gz \
//         internal/jscomp/testdata/transpileb/js <the *.js names>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.TranspileBDump --full b.gz \
//         .cache/corpus --list=internal/jscomp/testdata/transpileb/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/transpileb/golden.txt.gz
//   - TestTranspileBOracle compares any other dump (hashes, or --full):
//     TRANSPILEB_ORACLE=<dump.gz> [TRANSPILEB_ROOT=<root>] (default root: the
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
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// transpileBOracleFirst is TranspileBDump.FIRST.
const transpileBOracleFirst = "transpileb:first.js"

// transpileBOracleVariants is TranspileBDump.VARIANTS.
var transpileBOracleVariants = []string{"es5", "es2015", "es5-isolate"}

// transpileBOracleEscape is TranspileBDump.escape (ParsingDump's escaping).
func transpileBOracleEscape(s string) string {
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

// transpileBOracleException is TranspileBDump.exception:
// getClass().getName() + ": " + getMessage() (the Go port's nil
// dereferences carry no message, as the dump normalizes the JVM's).
func transpileBOracleException(r any) string {
	t := jutil.AsThrowable(r)
	full := t.Error()
	msg := t.GetMessage()
	if strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg) + ": " + msg
	}
	// No message: Java prints "null".
	return full + ": null"
}

func transpileBOracleSha16(lines []string) string {
	h := sha256.New()
	for _, l := range lines {
		io.WriteString(h, l+"\n")
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// transpileBOracleTreeLines is TranspileBDump.treeLines.
func transpileBOracleTreeLines(tree string) []string {
	var lines []string
	for _, l := range strings.Split(tree, "\n") {
		lines = append(lines, transpileBOracleEscape(l))
	}
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

type transpileBOracleDumper struct {
	full bool
	sb   strings.Builder
}

func (d *transpileBOracleDumper) line(s string) {
	d.sb.WriteString(transpileBOracleEscape(s))
	d.sb.WriteByte('\n')
}

// newCompiler mirrors TranspileBDump.newCompiler; variant "" is the parse
// check (Java: null).
func (d *transpileBOracleDumper) newCompiler(rel string, src string, variant string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	switch variant {
	case "":
		options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	case "es2015":
		options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	default:
		options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	}
	if variant != "" {
		options.SetRewritePolyfills(true)
		options.SetIsolatePolyfills(variant == "es5-isolate")
	}
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	compiler.Init(
		[]*SourceFile{},
		[]*SourceFile{SourceFile_FromCode(transpileBOracleFirst, ""), SourceFile_FromCode(rel, src)},
		options)
	compiler.Parse()
	return compiler
}

// transpileBOracleNamedPass is TranspileBDump.NamedPass.
type transpileBOracleNamedPass struct {
	name        string
	factory     func(*Compiler) CompilerPass
	passFactory *PassFactory
}

func (p *transpileBOracleNamedPass) create(compiler *Compiler) CompilerPass {
	if p.passFactory != nil {
		p.passFactory.ValidatePreconditions(compiler.GetOptions())
		return p.passFactory.Create(compiler)
	}
	return p.factory(compiler)
}

// transpileBOraclePasses mirrors TranspileBDump.passes.
func transpileBOraclePasses(compiler *Compiler, variant string) []*transpileBOracleNamedPass {
	options := compiler.GetOptions()
	var passes []*transpileBOracleNamedPass
	// As TranspilationPasses.injectTranspilationRuntimeLibraries and .rewritePolyfills create them.
	passes = append(passes, &transpileBOracleNamedPass{name: "es6InjectRuntimeLibraries", factory: func(c *Compiler) CompilerPass {
		return NewInjectTranspilationRuntimeLibraries(c)
	}})
	if options.GetRewritePolyfills() ||
		options.GetIsolatePolyfills() ||
		options.GetInjectPolyfillsNewerThan() != nil {
		passes = append(passes, &transpileBOracleNamedPass{name: "RewritePolyfills", factory: func(c *Compiler) CompilerPass {
			return NewRewritePolyfills(
				c,
				c.GetOptions().GetRewritePolyfills(),
				c.GetOptions().GetIsolatePolyfills(),
				c.GetOptions().GetInjectPolyfillsNewerThan())
		}})
	}
	passes = append(passes, &transpileBOracleNamedPass{name: "InjectRuntimeLibraries", factory: func(c *Compiler) CompilerPass {
		return NewInjectRuntimeLibraries(
			c, jutil.NewLinkedHashSet(c.GetOptions().GetForceLibraryInjectionList()...))
	}})
	passes = append(passes, &transpileBOracleNamedPass{name: PassNames_NORMALIZE, factory: func(c *Compiler) CompilerPass {
		return Normalize_CreateNormalizeForOptimizations(c)
	}})
	passes = append(passes, &transpileBOracleNamedPass{name: PassNames_GATHER_GETTERS_AND_SETTERS, factory: func(c *Compiler) CompilerPass {
		return NewGatherGetterAndSetterProperties(c)
	}})
	transpilation := NewPassListBuilder(options)
	TranspilationPasses_AddTranspilationPasses(transpilation, options)
	for _, f := range transpilation.Build() {
		passes = append(passes, &transpileBOracleNamedPass{name: f.GetName(), passFactory: f})
	}
	if variant == "es5-isolate" {
		passes = append(passes, &transpileBOracleNamedPass{name: "IsolatePolyfills", factory: func(c *Compiler) CompilerPass {
			return NewIsolatePolyfills(c)
		}})
	}
	return passes
}

func transpileBOracleDiagnostics(compiler *Compiler) []string {
	var out []string
	for _, e := range compiler.GetErrors() {
		out = append(out, transpileBOracleDiag("ERROR", e))
	}
	for _, e := range compiler.GetWarnings() {
		out = append(out, transpileBOracleDiag("WARNING", e))
	}
	return out
}

func transpileBOracleDiag(level string, e *JSError) string {
	return "diag " + level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " \"" + e.GetDescription() + "\""
}

// transpileBOracleRemoveFirst is List.remove(Object).
func transpileBOracleRemoveFirst(list []string, s string) []string {
	for i, x := range list {
		if x == s {
			return append(list[:i:i], list[i+1:]...)
		}
	}
	return list
}

// runPass runs one pass; it returns the "exception ..." line if the pass
// panicked with a Java exception.
func (d *transpileBOracleDumper) runPass(c *Compiler, pass *transpileBOracleNamedPass) (exception string) {
	defer func() {
		if r := recover(); r != nil {
			exception = "exception " + transpileBOracleException(r)
		}
	}()
	pass.create(c).Process(c.GetExternsRoot(), c.GetJsRoot())
	return ""
}

// dumpFile mirrors TranspileBDump.dumpFile.
func (d *transpileBOracleDumper) dumpFile(rel string, src string) {
	d.line("## file " + rel)
	compiler := d.newCompiler(rel, src, "")
	parseDiags := transpileBOracleDiagnostics(compiler)
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
	for _, variant := range transpileBOracleVariants {
		d.line("## variant " + variant)
		c := d.newCompiler(rel, src, variant)
		seen := transpileBOracleDiagnostics(c)
		tree := c.GetRoot().ToStringTree()
		ended := false
		for _, pass := range transpileBOraclePasses(c, variant) {
			exception := d.runPass(c, pass)
			after := c.GetRoot().ToStringTree()
			if after == tree {
				d.line("pass " + pass.name + " same")
			} else {
				d.line("pass " + pass.name + " " + transpileBOracleSha16(transpileBOracleTreeLines(after)))
			}
			tree = after
			diags := transpileBOracleDiagnostics(c)
			added := append([]string(nil), diags...)
			for _, diag := range seen {
				added = transpileBOracleRemoveFirst(added, diag)
			}
			for _, diag := range added {
				d.line(diag)
			}
			seen = diags
			if exception != "" {
				d.line(exception)
				ended = true
				break
			}
			if c.HasHaltingErrors() {
				d.line("halt")
				ended = true
				break
			}
		}
		if ended {
			continue
		}
		libs := c.GetRuntimeJsLibManager().GetInjectedLibraries()
		d.line("libs " + strconv.Itoa(len(libs)) + " " + strings.Join(libs, ","))
		d.line("stage " + c.GetLifeCycleStage().String())
		d.line("tree " + transpileBOracleSha16(transpileBOracleTreeLines(tree)))
		for script := c.GetJsRoot().GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			lines := transpileBOracleTreeLines(script.ToStringTree())
			name := jutil.ToJavaString(script.GetSourceFileName())
			d.line("script " + name + " " + transpileBOracleSha16(lines))
			if d.full && name != transpileBOracleFirst {
				for _, l := range lines {
					d.sb.WriteString(l)
					d.sb.WriteByte('\n')
				}
			}
		}
	}
}

// transpileBOracleSplit reads a dump into per-file sections; full reports
// whether it holds --full trees.
func transpileBOracleSplit(t *testing.T, path string) (names []string, sections map[string]string, full bool) {
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
		if strings.HasPrefix(l, "    ") || strings.HasPrefix(l, "SCRIPT") {
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

// transpileBOracleUnescapeName undoes the escaping of a file name.
func transpileBOracleUnescapeName(s string) string {
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

// transpileBOracleCompare dumps every file of the dump at path (sources
// under root) and compares; it returns the number of files compared.
func transpileBOracleCompare(t *testing.T, path string, root string, onlyNames func(string) bool) int {
	t.Helper()
	// The dump replaces the JVM's helpful NullPointerException messages
	// by "null", as -XX:-ShowCodeDetailsInExceptionMessages would.
	defer jutil.SetShowCodeDetailsInExceptionMessages(jutil.SetShowCodeDetailsInExceptionMessages(false))
	names, sections, full := transpileBOracleSplit(t, path)
	n, failures := 0, 0
	for _, name := range names {
		rel := transpileBOracleUnescapeName(name)
		if onlyNames != nil && !onlyNames(rel) {
			continue
		}
		if corpus.IsExcluded(rel) { // tools/oracle/gpl-excluded-corpus.txt
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Logf("skipping %s: %v", rel, err)
			continue
		}
		d := &transpileBOracleDumper{full: full}
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.line("go panic " + transpileBOracleException(r))
				}
			}()
			d.dumpFile(rel, jutil.NewStringFromUTF8Bytes(src))
		}()
		n++
		if got, want := d.sb.String(), sections[name]; got != want {
			failures++
			if failures <= 10 {
				t.Errorf("%s: dump differs:\n%s", rel, transpileBOracleFirstDiff(got, want))
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, n)
	}
	return n
}

func transpileBOracleFirstDiff(got, want string) string {
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

func TestTranspileBGolden(t *testing.T) {
	golden := "testdata/transpileb/golden.txt.gz"
	jsDir := "testdata/transpileb/js"
	hand := map[string]bool{}
	entries, err := os.ReadDir(jsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		hand[e.Name()] = true
	}
	if n := transpileBOracleCompare(t, golden, jsDir, func(rel string) bool { return hand[rel] }); n != len(hand) {
		t.Errorf("compared %d hand-written files, want %d", n, len(hand))
	}
	corpusRoot := "/home/user/jsgo/.cache/corpus"
	if _, err := os.Stat(corpusRoot); err != nil {
		t.Logf("corpus not present; golden corpus files skipped")
		return
	}
	transpileBOracleCompare(t, golden, corpusRoot, func(rel string) bool { return !hand[rel] })
}

func TestTranspileBOracle(t *testing.T) {
	dump := os.Getenv("TRANSPILEB_ORACLE")
	if dump == "" {
		t.Skip("set TRANSPILEB_ORACLE=<TranspileBDump output .gz> (and TRANSPILEB_ROOT) to run")
	}
	root := os.Getenv("TRANSPILEB_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	n := transpileBOracleCompare(t, dump, root, nil)
	t.Logf("compared %d files", n)
}
