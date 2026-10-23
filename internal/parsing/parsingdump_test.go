package parsing

// Differential oracle at the ParserRunner.parse boundary, Go side.
// dumpParsing reproduces, byte for byte, the format of
// tools/oracle/java/ParsingDump.java (see its header): the rhino.Node AST
// with JSDoc (through Node.toString(true, true, false) and
// JSDocInfo.toStringVerbose()), the features, the source map URL, the
// comments, and every message sent to the rhino ErrorReporter.
//
//   - TestParsingDumpGolden compares against the checked-in Java dumps in
//     testdata/parse/<config>/<rel>.parse.gz of the sources in
//     testdata/parse/src/<rel> (always runs): IRFactory-specific inputs
//     (closure-unaware code, comment attachment, validations, literals) under
//     every JsDocParsing mode, ES3/ES5/ES2015 and KEEP_GOING, and broken
//     inputs under KEEP_GOING. Regenerate a directory with
//       java -Xss1g -cp closure-compiler.jar:<build> ParsingDump --gzip <options> \
//         testdata/parse/<config> testdata/parse/src <rel>...
//     where <options> is "" (default), --keep-going, --keep-going
//     --jsdoc-mode=INCLUDE_ALL_COMMENTS, --jsdoc-mode=<MODE> (include-*,
//     license-comments-only), --lang=<MODE> (ecmascript*), or --lang=ECMASCRIPT3
//     --keep-going --jsdoc-mode=INCLUDE_ALL_COMMENTS.
//   - TestParsingDumpOracle compares a whole set when PARSE_ROOT (the root of
//     the sources) and PARSE_ORACLE_DIR (the ParsingDump output for it, as
//     <rel>.parse or <rel>.parse.gz) are set; PARSE_FILES optionally names a
//     file listing the relative paths (default: every .js/.mjs/.cjs under the
//     root), and PARSE_DIFF_DIR, if set, receives the Go dump of every
//     mismatching file. Use absolute paths: `go test` runs in the package
//     directory.
//
// The Java side passes a jscomp.SourceFile (SourceFile.fromCode, kind STRONG)
// as the StaticSourceFile; it renders as its name inside [source_file: ...],
// as rhino.SimpleSourceFile does, and nothing on the parse path reads
// anything else from it but the name.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// recordingReporter is ParsingDump's anonymous ErrorReporter.
type recordingReporter struct {
	diags []string
}

func (r *recordingReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.diags = append(r.diags, fmt.Sprintf("WARNING %d:%d %s", line, lineOffset, dumpQuote(message)))
}

func (r *recordingReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.diags = append(r.diags, fmt.Sprintf("ERROR %d:%d %s", line, lineOffset, dumpQuote(message)))
}

// compilerParserConfig is the parsing.Config that Compiler.getParserConfig
// builds for the given options: ParserRunner.createConfig with the
// CompilerOptions defaults (no extra annotation names, inline source maps
// parsed).
func compilerParserConfig(languageMode Config_LanguageMode, jsDocParsing Config_JsDocParsing,
	runMode Config_RunMode, strictMode Config_StrictMode) *Config {
	return ParserRunner_CreateConfigWithModes(languageMode, jsDocParsing, runMode, nil, true, strictMode)
}

// defaultDumpConfig is the configuration of the Java dumps (`# config`
// languageMode=ES_NEXT jsDocParsingMode=TYPES_ONLY runMode=STOP_AFTER_ERROR
// strictMode=STRICT).
func defaultDumpConfig() *Config {
	return compilerParserConfig(Config_LanguageMode_ES_NEXT, Config_JsDocParsing_TYPES_ONLY,
		Config_RunMode_STOP_AFTER_ERROR, Config_StrictMode_STRICT)
}

func enumByName[E fmt.Stringer](values []E, name string) (E, bool) {
	for _, v := range values {
		if v.String() == name {
			return v, true
		}
	}
	var zero E
	return zero, false
}

// configFromDump rebuilds the Config named by the "# config" line of a Java
// dump (the Go dump prints the fields of the config built here, so a wrong
// reconstruction shows up as a diff on that line).
func configFromDump(dump string) *Config {
	first, _, _ := strings.Cut(dump, "\n")
	if !strings.HasPrefix(first, "# config ") {
		return defaultDumpConfig()
	}
	lang, jsdoc, run, strict := Config_LanguageMode_ES_NEXT, Config_JsDocParsing_TYPES_ONLY,
		Config_RunMode_STOP_AFTER_ERROR, Config_StrictMode_STRICT
	for _, kv := range strings.Fields(strings.TrimPrefix(first, "# config ")) {
		k, v, _ := strings.Cut(kv, "=")
		switch k {
		case "languageMode":
			lang, _ = enumByName(Config_LanguageMode_values(), v)
		case "jsDocParsingMode":
			jsdoc, _ = enumByName(Config_JsDocParsing_values(), v)
		case "runMode":
			run, _ = enumByName(Config_RunMode_values(), v)
		case "strictMode":
			strict, _ = enumByName(Config_StrictMode_values(), v)
		}
	}
	return compilerParserConfig(lang, jsdoc, run, strict)
}

// dumpParsing is ParsingDump.dumpOne.
func dumpParsing(name, src string, config *Config) string {
	var sb strings.Builder
	line := func(s string) {
		sb.WriteString(dumpEscape(s))
		sb.WriteByte('\n')
	}
	line("# config languageMode=" + config.LanguageMode().String() +
		" jsDocParsingMode=" + config.JsDocParsingMode().String() +
		" runMode=" + config.RunMode().String() +
		" strictMode=" + config.StrictMode().String())
	reporter := &recordingReporter{}
	var result *ParserRunner_ParseResult
	var exception jutil.Throwable
	jutil.Try(func() {
		result = ParserRunner_Parse(
			rhino.NewSimpleSourceFile(name, rhino.StaticSourceFile_SourceKind_STRONG), src, config, reporter)
	}, func(e jutil.IRuntimeException) { exception = e })
	line("# diagnostics")
	for _, d := range reporter.diags {
		line(d)
	}
	if exception != nil {
		line("# exception " + javaClassName(exception) + ": " + exception.GetMessage())
		return sb.String()
	}
	line("# features " + result.Features.String())
	if result.SourceMapURL == nil {
		line("# sourceMapURL null")
	} else {
		line("# sourceMapURL " + dumpQuote(*result.SourceMapURL))
	}
	line("# comments")
	for _, c := range result.Comments {
		line(c.Type.String() + " " + dumpRange(c.Location) + " " + dumpQuote(c.Value))
	}
	line("# tree")
	if result.Ast == nil {
		line("null")
	} else {
		dumpTree(line, result.Ast, 0)
	}
	return sb.String()
}

func dumpTree(line func(string), n *rhino.Node, level int) {
	indent := strings.Repeat("    ", level)
	line(indent + n.ToStringWithFlags(true, true, false))
	if info := n.GetJSDocInfo(); info != nil {
		line(indent + "  @jsdoc " + info.ToStringVerbose())
	}
	if shadow := n.GetClosureUnawareShadow(); shadow != nil {
		line(indent + "  @shadow")
		dumpTree(line, shadow, level+1)
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		dumpTree(line, c, level+1)
	}
}

// javaClassName is Throwable.getClass().getName().
func javaClassName(t jutil.Throwable) string {
	full := t.Error()
	if msg := t.GetMessage(); strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg)
	}
	return full
}

func dumpRange(r parser.SourceRange) string {
	return dumpPos(r.Start) + "-" + dumpPos(r.End)
}

func dumpPos(p parser.SourcePosition) string {
	return fmt.Sprintf("%d:%d@%d", p.Line, p.Column, p.Offset)
}

// dumpQuote is ParsingDump.quote: double quotes only (escaping happens once,
// per line).
func dumpQuote(s string) string { return "\"" + s + "\"" }

// dumpEscape is ParsingDump.escape: backslash and every UTF-16 code unit
// outside printable ASCII become escapes.
func dumpEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, "\\u%04x", c)
		}
	}
	return b.String()
}

func readDump(path string) (string, error) {
	if b, err := os.ReadFile(path); err == nil {
		return string(b), nil
	}
	f, err := os.Open(path + ".gz")
	if err != nil {
		return "", err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(zr)
	return string(b), err
}

// dumpIdentityHash matches the identity hash that Object.toString() prints
// for a JSDocInfo.Marker (JSDocInfo.toStringVerbose lists the markers when
// descriptions are parsed). It is arbitrary on both sides, so sameDump masks
// it; the number and position of the markers are still compared.
var dumpIdentityHash = regexp.MustCompile(`com\.google\.javascript\.rhino\.JSDocInfo\$Marker@[0-9a-f]+`)

// sameDump compares two dumps, ignoring only Marker identity hashes; on a
// difference it returns the masked dumps.
func sameDump(want, got string) (bool, string, string) {
	if want == got {
		return true, want, got
	}
	w := dumpIdentityHash.ReplaceAllLiteralString(want, "com.google.javascript.rhino.JSDocInfo$Marker@")
	g := dumpIdentityHash.ReplaceAllLiteralString(got, "com.google.javascript.rhino.JSDocInfo$Marker@")
	return w == g, w, g
}

// firstDiff describes the first differing line of two dumps.
func firstDiff(want, got string) string {
	wl := strings.Split(want, "\n")
	gl := strings.Split(got, "\n")
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return fmt.Sprintf("line %d:\n  java: %s\n  go:   %s", i+1, w, g)
		}
	}
	return "(identical)"
}

func TestParsingDumpGolden(t *testing.T) {
	const root = "testdata/parse"
	modes, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, mode := range modes {
		if !mode.IsDir() || mode.Name() == "src" {
			continue
		}
		dir := filepath.Join(root, mode.Name())
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".parse.gz") {
				return err
			}
			rel, _ := filepath.Rel(dir, strings.TrimSuffix(path, ".parse.gz"))
			rel = filepath.ToSlash(rel) // the Java dumps were made with this name
			src, err := os.ReadFile(filepath.Join(root, "src", rel))
			if err != nil {
				return err
			}
			want, err := readDump(strings.TrimSuffix(path, ".gz"))
			if err != nil {
				return err
			}
			n++
			got := dumpParsing(rel, string(src), configFromDump(want))
			if same, w, g := sameDump(want, got); !same {
				t.Errorf("%s/%s: dump differs from Java at %s", mode.Name(), rel, firstDiff(w, g))
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if n < 100 {
		t.Fatalf("only %d golden dumps", n)
	}
}

func TestParsingDumpOracle(t *testing.T) {
	root := os.Getenv("PARSE_ROOT")
	oracle := os.Getenv("PARSE_ORACLE_DIR")
	if root == "" || oracle == "" {
		t.Skip("set PARSE_ROOT and PARSE_ORACLE_DIR to run the ParserRunner.parse oracle")
	}
	var rels []string
	if list := os.Getenv("PARSE_FILES"); list != "" {
		f, err := os.Open(list)
		if err != nil {
			t.Fatal(err)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if l := strings.TrimSpace(sc.Text()); l != "" {
				rels = append(rels, l)
			}
		}
		f.Close()
	} else {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				switch filepath.Ext(path) {
				case ".js", ".mjs", ".cjs":
					rel, _ := filepath.Rel(root, path)
					rels = append(rels, rel)
				}
			}
			return nil
		})
		sort.Strings(rels)
	}
	rels = corpus.Filter(rels) // tools/oracle/gpl-excluded-corpus.txt
	diffDir := os.Getenv("PARSE_DIFF_DIR")
	type result struct {
		rel, diff, got string
		missing        bool
	}
	results := make([]result, len(rels))
	var wg sync.WaitGroup
	work := make(chan int)
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range work {
				rel := rels[i]
				want, err := readDump(filepath.Join(oracle, rel+".parse"))
				if err != nil {
					results[i] = result{rel: rel, missing: true}
					continue
				}
				src, err := os.ReadFile(filepath.Join(root, rel))
				if err != nil {
					t.Error(err)
					continue
				}
				got := dumpParsing(rel, string(src), configFromDump(want))
				if same, w, g := sameDump(want, got); !same {
					results[i] = result{rel: rel, diff: firstDiff(w, g), got: got}
				}
			}
		}()
	}
	for i := range rels {
		work <- i
	}
	close(work)
	wg.Wait()
	same, missing, differ := 0, 0, 0
	for _, r := range results {
		switch {
		case r.missing:
			missing++
		case r.diff != "":
			differ++
			if differ <= 50 {
				t.Errorf("%s: %s", r.rel, r.diff)
			}
			if diffDir != "" {
				out := filepath.Join(diffDir, r.rel+".go.parse")
				os.MkdirAll(filepath.Dir(out), 0o755)
				os.WriteFile(out, []byte(r.got), 0o644)
			}
		default:
			same++
		}
	}
	t.Logf("compared %d files: %d identical, %d differ, %d without oracle dump", len(rels)-missing, same, differ, missing)
	if missing > 0 {
		t.Errorf("%d files have no oracle dump", missing)
	}
}
