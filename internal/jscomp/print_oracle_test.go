package jscomp

// Differential oracle for the code printer, Go side of
// tools/oracle/java/PrintDump.java (see its header for the format).
// printDumpOne reproduces PrintDump.dumpOne byte for byte with the ported
// CodePrinter, JSDocInfoPrinter, SourceInformationAnnotator, SourceMap and
// DotFormatter.
//
//   - TestPrintDumpGolden compares against the checked-in Java dumps
//     testdata/print/dump/own/<rel>.print.gz of the printer edge cases in
//     testdata/print/src/<rel>, and testdata/print/dump/parsing/<rel>.print.gz
//     of the parsing layer's golden sources ../parsing/testdata/parse/src/<rel>
//     (always runs). Regenerate with
//       java -Xss1g -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.PrintDump --gzip \
//         testdata/print/dump/own testdata/print/src <rel>...
//     (and likewise for dump/parsing).
//   - TestPrintDumpOracle compares a whole set when PRINT_ROOT (the root of
//     the sources) and PRINT_ORACLE_DIR (the PrintDump output for it, as
//     <rel>.print or <rel>.print.gz) are set; PRINT_FILES optionally names a
//     file listing the relative paths (default: every .print(.gz) under
//     PRINT_ORACLE_DIR), and PRINT_DIFF_DIR, if set, receives the Go dump of
//     every mismatching file. Use absolute paths: `go test` runs in the
//     package directory.
//
// The options a dump was made with (--only, --no-mappings, --dot-limit) are
// read back from its first line.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// printDumpOptions are PrintDump's command-line options.
type printDumpOptions struct {
	only       map[string]bool // nil selects everything
	noMappings bool
	dotLimit   int
}

func (o *printDumpOptions) selected(name string) bool {
	return o.only == nil || o.only[name]
}

func (o *printDumpOptions) header() string {
	only := "*"
	if o.only != nil {
		var names []string
		for n := range o.only {
			names = append(names, n)
		}
		sort.Slice(names, func(i, j int) bool { return jutil.StringCompareTo(names[i], names[j]) < 0 })
		only = strings.Join(names, ",")
	}
	return "# options only=" + only + " no-mappings=" + strconv.FormatBool(o.noMappings) +
		" dot-limit=" + strconv.Itoa(o.dotLimit)
}

// printDumpOptionsFromDump reads the options back from the first line of a
// Java dump.
func printDumpOptionsFromDump(dump string) *printDumpOptions {
	o := &printDumpOptions{dotLimit: 4096}
	first, _, _ := strings.Cut(dump, "\n")
	for _, kv := range strings.Fields(strings.TrimPrefix(first, "# options ")) {
		k, v, _ := strings.Cut(kv, "=")
		switch k {
		case "only":
			if v != "*" {
				o.only = map[string]bool{}
				for _, n := range strings.Split(v, ",") {
					o.only[n] = true
				}
			}
		case "no-mappings":
			o.noMappings = v == "true"
		case "dot-limit":
			o.dotLimit, _ = strconv.Atoi(v)
		}
	}
	return o
}

// printDumpWriter accumulates the escaped lines of a dump.
type printDumpWriter struct {
	sb strings.Builder
}

// line is PrintDump.line: the line escaped as ParsingDump.escape does.
func (w *printDumpWriter) line(s string) {
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			w.sb.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			w.sb.WriteByte(byte(c))
		default:
			fmt.Fprintf(&w.sb, "\\u%04x", c)
		}
	}
	w.sb.WriteByte('\n')
}

// lines is PrintDump.lines: s split at '\n' as "| " lines.
func (w *printDumpWriter) lines(s string) {
	for {
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			w.line("| " + s)
			return
		}
		w.line("| " + s[:i])
		s = s[i+1:]
	}
}

// exception is PrintDump.exception for a recovered panic.
func (w *printDumpWriter) exception(r any) {
	t := jutil.AsThrowable(r)
	full := t.Error()
	msg := t.GetMessage()
	class := full
	if strings.HasSuffix(full, ": "+msg) {
		class = strings.TrimSuffix(full, ": "+msg)
	} else {
		msg = "null"
	}
	w.line("# exception " + class + ": " + msg)
}

// try runs body, writing the exception line if it panics.
func (w *printDumpWriter) try(body func()) {
	defer func() {
		if r := recover(); r != nil {
			w.exception(r)
		}
	}()
	body()
}

// printDumpParserConfig is the parsing.Config of PrintDump.compilerParserConfig.
func printDumpParserConfig(jsdoc parsing.Config_JsDocParsing) *parsing.Config {
	return parsing.ParserRunner_CreateConfigWithModes(parsing.Config_LanguageMode_ES_NEXT, jsdoc,
		parsing.Config_RunMode_STOP_AFTER_ERROR, nil, true, parsing.Config_StrictMode_STRICT)
}

// printDumpReporter is PrintDump's anonymous ErrorReporter: it counts the
// errors.
type printDumpReporter struct {
	errors int
}

func (r *printDumpReporter) Warning(message string, sourceName string, line int, lineOffset int) {}

func (r *printDumpReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.errors++
}

// printDumpParse is PrintDump.parse; w may be nil.
func printDumpParse(w *printDumpWriter, config *parsing.Config, name, src string) *rhino.Node {
	reporter := &printDumpReporter{}
	var result *parsing.ParserRunner_ParseResult
	var exception jutil.Throwable
	jutil.Try(func() {
		result = parsing.ParserRunner_Parse(SourceFile_FromCode(name, src), src, config, reporter)
	}, func(e jutil.IRuntimeException) { exception = e })
	if exception != nil {
		if w != nil {
			full := exception.Error()
			msg := exception.GetMessage()
			class := strings.TrimSuffix(full, ": "+msg)
			w.line("# parse exception " + class + ": " + msg)
		}
		return nil
	}
	if reporter.errors > 0 || result.Ast == nil {
		if w != nil {
			w.line("# parse errors " + strconv.Itoa(reporter.errors))
		}
		return nil
	}
	if w != nil {
		w.line("# parse ok")
	}
	return result.Ast
}

// printDumpOne is PrintDump.dumpOne.
func printDumpOne(o *printDumpOptions, name, src string) string {
	w := &printDumpWriter{}
	w.line(o.header())
	root := printDumpParse(w, printDumpParserConfig(parsing.Config_JsDocParsing_TYPES_ONLY), name, src)
	if root == nil {
		return w.sb.String()
	}

	print := func(name string, root *rhino.Node, configure func(b *CodePrinter_Builder)) {
		if !o.selected(name) {
			return
		}
		w.line("# print " + name)
		w.try(func() {
			b := NewCodePrinter_Builder(root)
			configure(b)
			w.lines(b.Build())
		})
	}
	withOptions := func(set func(o *CompilerOptions)) func(b *CodePrinter_Builder) {
		return func(b *CodePrinter_Builder) {
			opts := NewCompilerOptions()
			set(opts)
			b.SetCompilerOptions(opts)
		}
	}

	// Plain configurations.
	print("compact", root, func(b *CodePrinter_Builder) {})
	print("pretty", root, func(b *CodePrinter_Builder) { b.SetPrettyPrint(true) })
	print("linebreak", root, func(b *CodePrinter_Builder) { b.SetLineBreak(true) })
	print("compact-threshold-60", root, withOptions(func(o *CompilerOptions) { o.SetLineLengthThreshold(60) }))
	print("single-quotes", root, withOptions(func(o *CompilerOptions) { o.SetPreferSingleQuotes(true) }))
	print("es3-untrusted", root, withOptions(func(o *CompilerOptions) {
		o.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
		o.SetTrustedStrings(false)
	}))
	print("es5-utf8", root, withOptions(func(o *CompilerOptions) {
		o.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
		o.SetOutputCharset(jutil.StandardCharsets_UTF_8)
	}))
	print("pretty-types", root, withOptions(func(o *CompilerOptions) {
		o.SetPreserveTypeAnnotations(true)
		o.SetPrettyPrint(true)
	}))
	print("typesummary-strict", root, func(b *CodePrinter_Builder) {
		b.SetPrettyPrint(true).SetTagAsTypeSummary(true).SetTagAsStrict(true)
	})
	print("typed", root, func(b *CodePrinter_Builder) {
		b.SetOutputTypes(true).SetTypeRegistry(rhino.NewJSTypeRegistry(nil))
	})

	// Source maps, after the original names are recorded as the compiler
	// does.
	compiler := newTestCompiler(NewCompilerOptions())
	NodeTraversal_Traverse(compiler, root, SourceInformationAnnotator_Create())
	printDumpWithSourceMap(o, w, "sourcemap-compact", root, false, false, 500, SourceMap_DetailLevel_ALL)
	printDumpWithSourceMap(o, w, "sourcemap-compact-threshold-60", root, false, false, 60, SourceMap_DetailLevel_ALL)
	printDumpWithSourceMap(o, w, "sourcemap-linebreak", root, false, true, 500, SourceMap_DetailLevel_ALL)
	printDumpWithSourceMap(o, w, "sourcemap-pretty", root, true, false, 500, SourceMap_DetailLevel_ALL)
	printDumpWithSourceMap(o, w, "sourcemap-symbols", root, false, false, 500, SourceMap_DetailLevel_SYMBOLS)
	print("original-names", root, withOptions(func(o *CompilerOptions) { o.SetUseOriginalNamesInOutput(true) }))

	// Comments: a parse that keeps every comment.
	var allRoot *rhino.Node
	if o.selected("pretty-comments") || o.selected("compact-comments") || o.selected("jsdoc-all") {
		allRoot = printDumpParse(nil, printDumpParserConfig(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS), name, src)
	}
	if allRoot != nil {
		print("pretty-comments", allRoot, withOptions(func(o *CompilerOptions) {
			o.SetPreserveTypeAnnotations(true)
			o.SetPreserveNonJSDocComments(true)
			o.SetPrettyPrint(true)
		}))
		print("compact-comments", allRoot, withOptions(func(o *CompilerOptions) {
			o.SetPreserveNonJSDocComments(true)
		}))
	}

	// JSDocInfoPrinter on every JSDocInfo.
	if o.selected("jsdoc-types") {
		printDumpJsDoc(w, "types", root, NewJSDocInfoPrinter(false))
	}
	if allRoot != nil && o.selected("jsdoc-all") {
		printDumpJsDoc(w, "all", allRoot, NewJSDocInfoPrinterWithPrintDesc(false, true))
	}

	if jutil.StringLength(src) <= o.dotLimit && o.selected("dot") {
		w.line("# dot")
		w.try(func() { w.lines(DotFormatter_ToDot(root)) })
	}
	return w.sb.String()
}

// printDumpWithSourceMap is PrintDump.printWithSourceMap.
func printDumpWithSourceMap(o *printDumpOptions, w *printDumpWriter, name string, root *rhino.Node,
	pretty bool, lineBreak bool, threshold int, detail SourceMap_DetailLevel) {
	if !o.selected(name) {
		return
	}
	w.line("# print " + name)
	w.try(func() {
		opts := NewCompilerOptions()
		opts.SetSourceMapOutputPath(jutil.StringPtr("out.map"))
		opts.SetSourceMapDetailLevel(detail)
		opts.SetLineLengthThreshold(threshold)
		opts.SetPrettyPrint(pretty)
		opts.SetLineBreak(lineBreak)
		result := NewCodePrinter_Builder(root).SetCompilerOptions(opts).BuildWithSourceMappings()
		w.lines(result.source)
		if !o.noMappings {
			w.line("# mappings")
		}
		sourceMap := SourceMap_Format_V3.GetInstance()
		for _, m := range result.mappings {
			original := "-"
			if m.node.GetOriginalName() != nil {
				original = *m.node.GetOriginalName()
			}
			if !o.noMappings {
				w.line(m.node.GetToken().String() + " " +
					strconv.Itoa(m.node.GetLineno()) + ":" + strconv.Itoa(m.node.GetCharno()) + " " +
					original + " " +
					strconv.Itoa(m.start.GetLine()) + ":" + strconv.Itoa(m.start.GetColumn()) + " " +
					strconv.Itoa(m.end.GetLine()) + ":" + strconv.Itoa(m.end.GetColumn()))
			}
			sourceMap.AddMapping(m)
		}
		w.line("# map")
		var out strings.Builder
		sourceMap.AppendTo(&out, jutil.StringPtr("out.js"))
		w.lines(out.String())
	})
}

// printDumpJsDoc is PrintDump.dumpJsDoc.
func printDumpJsDoc(w *printDumpWriter, parse string, n *rhino.Node, printer *JSDocInfoPrinter) {
	if info := n.GetJSDocInfo(); info != nil {
		w.line("# jsdoc " + parse + " " + n.GetToken().String() + " " +
			strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno()))
		w.try(func() { w.lines(printer.Print(info)) })
	}
	if shadow := n.GetClosureUnawareShadow(); shadow != nil {
		printDumpJsDoc(w, parse, shadow, printer)
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		printDumpJsDoc(w, parse, c, printer)
	}
}

func printDumpRead(path string) (string, error) {
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

// printDumpFirstDiff describes the first differing line of two dumps, with
// the section it is in.
func printDumpFirstDiff(want, got string) string {
	wl := strings.Split(want, "\n")
	gl := strings.Split(got, "\n")
	section := ""
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return fmt.Sprintf("line %d (in %s):\n  java: %s\n  go:   %s", i+1, section, w, g)
		}
		if strings.HasPrefix(w, "# ") {
			section = w
		}
	}
	return "(identical)"
}

func TestPrintDumpGolden(t *testing.T) {
	n := 0
	for _, set := range []struct{ dumps, srcs string }{
		{"testdata/print/dump/own", "testdata/print/src"},
		{"testdata/print/dump/parsing", "../parsing/testdata/parse/src"},
	} {
		err := filepath.Walk(set.dumps, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".print.gz") {
				return err
			}
			rel, _ := filepath.Rel(set.dumps, strings.TrimSuffix(path, ".print.gz"))
			rel = filepath.ToSlash(rel) // the Java dumps were made with this name
			src, err := os.ReadFile(filepath.Join(set.srcs, rel))
			if err != nil {
				return err
			}
			want, err := printDumpRead(strings.TrimSuffix(path, ".gz"))
			if err != nil {
				return err
			}
			n++
			got := printDumpOne(printDumpOptionsFromDump(want), rel, string(src))
			if want != got {
				t.Errorf("%s: dump differs from Java at %s", path, printDumpFirstDiff(want, got))
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

func TestPrintDumpOracle(t *testing.T) {
	root := os.Getenv("PRINT_ROOT")
	oracle := os.Getenv("PRINT_ORACLE_DIR")
	if root == "" || oracle == "" {
		t.Skip("set PRINT_ROOT and PRINT_ORACLE_DIR to run the CodePrinter oracle")
	}
	var rels []string
	if list := os.Getenv("PRINT_FILES"); list != "" {
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
		filepath.Walk(oracle, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				for _, suffix := range []string{".print.gz", ".print"} {
					if strings.HasSuffix(path, suffix) {
						rel, _ := filepath.Rel(oracle, strings.TrimSuffix(path, suffix))
						rels = append(rels, filepath.ToSlash(rel))
						break
					}
				}
			}
			return nil
		})
		sort.Strings(rels)
	}
	diffDir := os.Getenv("PRINT_DIFF_DIR")
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
				want, err := printDumpRead(filepath.Join(oracle, rel+".print"))
				if err != nil {
					results[i] = result{rel: rel, missing: true}
					continue
				}
				src, err := os.ReadFile(filepath.Join(root, rel))
				if err != nil {
					t.Error(err)
					continue
				}
				got := printDumpOne(printDumpOptionsFromDump(want), rel, string(src))
				if want != got {
					results[i] = result{rel: rel, diff: printDumpFirstDiff(want, got), got: got}
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
				out := filepath.Join(diffDir, r.rel+".go.print")
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
