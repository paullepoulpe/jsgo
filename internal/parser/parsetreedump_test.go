package parser

// Differential parse-tree oracle, Go side. dumpParseTree reproduces, byte for
// byte, the format of tools/oracle/java/ParseTreeDump.java, using reflection
// over the ported trees and tokens (Go field names with the first letter
// lower-cased must equal the Java field names, so a naming mistake shows up
// as a diff).
//
//   - TestParseTreeGolden compares against the checked-in Java dumps in
//     testdata/ptree/<lang>/ (always runs); es_next/ uses the default
//     configuration, es5/ and es3/ were dumped with --lang=ECMASCRIPT5 /
//     ECMASCRIPT3 (sloppy mode, older Parser.Config modes).
//   - TestParseTreeOracle compares a whole corpus when PTREE_CORPUS (the
//     corpus root) and PTREE_ORACLE_DIR (the ParseTreeDump output for it) are
//     set; PTREE_FILES optionally names a file listing the relative paths
//     (default: every .js/.mjs/.cjs under the root), and PTREE_DIFF_DIR, if
//     set, receives the Go dump of every mismatching file.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// defaultParserConfig is the Parser.Config that ParserRunner.parse builds for
// the compiler's default configuration with --language_in=ECMASCRIPT_NEXT:
// ES_NEXT maps to Mode.ES8_OR_GREATER, and ECMASCRIPT_NEXT is default-strict
// (CompilerOptions.expectStrictModeInput()). The "# config" line of every
// dump checks that it agrees with what Java derives.
func defaultParserConfig() *Parser_Config {
	return NewParser_ConfigWithMode(Parser_Config_Mode_ES8_OR_GREATER /* isStrictMode= */, true)
}

type recordingErrorReporter struct {
	ErrorReporterBase
	lines []string
}

func newRecordingErrorReporter() *recordingErrorReporter {
	r := &recordingErrorReporter{}
	r.InitErrorReporterBase(r)
	return r
}

func (r *recordingErrorReporter) ReportError(location LinePosition, message string) {
	r.lines = append(r.lines, "ERROR "+dumpLinePosition(location)+" "+dumpQuote(message))
}

func (r *recordingErrorReporter) ReportWarning(location LinePosition, message string) {
	r.lines = append(r.lines, "WARNING "+dumpLinePosition(location)+" "+dumpQuote(message))
}

// configFromDump returns the Parser.Config whose fields the "# config" line
// of a Java dump shows (the Java side derives it from a --lang through
// ParserRunner; the Go dump prints the fields of the config built here, so a
// wrong reconstruction shows up as a diff on that line).
func configFromDump(dump string) *Parser_Config {
	first, _, _ := strings.Cut(dump, "\n")
	if !strings.HasPrefix(first, "# config ") {
		return defaultParserConfig()
	}
	fields := map[string]bool{}
	for _, kv := range strings.Fields(strings.TrimPrefix(first, "# config ")) {
		k, v, _ := strings.Cut(kv, "=")
		fields[k] = v == "true"
	}
	mode := Parser_Config_Mode_ES5
	switch {
	case fields["atLeast8"]:
		mode = Parser_Config_Mode_ES8_OR_GREATER
	case fields["atLeast6"]:
		mode = Parser_Config_Mode_ES6_OR_ES7
	case fields["warnTrailingCommas"]:
		mode = Parser_Config_Mode_ES3
	}
	return NewParser_ConfigWithMode(mode, fields["isStrictMode"])
}

// dumpParseTree parses src (named name) with config and returns the dump.
func dumpParseTree(name, src string, config *Parser_Config) string {
	var sb strings.Builder
	file := NewSourceFile(name, src)
	reporter := newRecordingErrorReporter()
	sb.WriteString("# config")
	cv := reflect.ValueOf(config).Elem()
	for _, f := range sortedStructFields(cv.Type()) {
		sb.WriteString(" " + f.name + "=" + fmt.Sprint(cv.FieldByIndex(f.index).Bool()))
	}
	sb.WriteByte('\n')
	p := NewParser(config, reporter, file)
	var tree *ProgramTree
	var exception string
	func() {
		defer func() {
			if r := recover(); r != nil {
				exception = describePanic(r)
			}
		}()
		tree = p.ParseProgram()
	}()
	sb.WriteString("# diagnostics\n")
	for _, l := range reporter.lines {
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	if exception != "" {
		sb.WriteString("# exception " + exception + "\n")
		return sb.String()
	}
	sb.WriteString("# features " + dumpQuote(p.GetFeatures().String()) + "\n")
	if url := p.GetSourceMapURL(); url == nil {
		sb.WriteString("# sourceMapURL null\n")
	} else {
		sb.WriteString("# sourceMapURL " + dumpQuote(*url) + "\n")
	}
	sb.WriteString("# comments\n")
	d := &dumper{sb: &sb, file: file}
	d.render(reflect.ValueOf(p.GetComments()), 0, "")
	sb.WriteString("# tree\n")
	if tree == nil {
		d.line(0, "", "null")
	} else {
		d.render(reflect.ValueOf(tree), 0, "")
	}
	return sb.String()
}

// describePanic renders a recovered panic like the Java side renders a
// thrown exception: "<class>: <quoted message>".
func describePanic(r any) string {
	if t, ok := r.(jutil.Throwable); ok {
		full := t.Error()
		msg := t.GetMessage()
		if !strings.Contains(full, ": ") {
			return full + ": " + dumpQuote("null")
		}
		return full[:len(full)-len(msg)-2] + ": " + dumpQuote(msg)
	}
	buf := make([]byte, 1<<16)
	buf = buf[:runtime.Stack(buf, false)]
	return "go-panic: " + dumpQuote(fmt.Sprint(r)) + "\n" + string(buf)
}

type dumper struct {
	sb   *strings.Builder
	file *SourceFile // positions are offsets: lines and columns come from here
}

func (d *dumper) line(depth int, label, text string) {
	for i := 0; i < depth; i++ {
		d.sb.WriteByte(' ')
	}
	d.sb.WriteString(label)
	d.sb.WriteString(text)
	d.sb.WriteByte('\n')
}

var (
	parseTreeType     = reflect.TypeOf((*ParseTree)(nil)).Elem()
	tokenType         = reflect.TypeOf((*Token)(nil)).Elem()
	commentPtrType    = reflect.TypeOf((*Comment)(nil))
	positionType      = reflect.TypeOf(SourcePosition{})
	rangeType         = reflect.TypeOf(SourceRange{})
	lineRangeType     = reflect.TypeOf(LineRange{})
	parseTreeBaseType = reflect.TypeOf(ParseTreeBase{})
	tokenBaseType     = reflect.TypeOf(TokenBase{})
)

type namer interface{ Name() string }

// render prints v (a field value, list element or root) at depth.
func (d *dumper) render(v reflect.Value, depth int, label string) {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			d.line(depth, label, "null")
			return
		}
		v = v.Elem()
		if v.Kind() == reflect.Pointer && v.IsNil() {
			// A typed nil pointer in an interface is non-nil in Go but would be
			// Java null: always a porting bug.
			d.line(depth, label, "TYPED-NIL "+v.Type().String())
			return
		}
	}
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			d.line(depth, label, "null")
			return
		}
		if v.Len() == 0 {
			d.line(depth, label, "[]")
			return
		}
		d.line(depth, label, fmt.Sprintf("[%d]", v.Len()))
		for i := 0; i < v.Len(); i++ {
			d.render(v.Index(i), depth+1, "")
		}
		return
	case reflect.Pointer:
		if v.IsNil() {
			d.line(depth, label, "null")
			return
		}
	}
	t := v.Type()
	switch {
	case t.Implements(parseTreeType):
		pt := v.Interface().(ParseTree).GetParseTreeBase()
		d.line(depth, label, javaSimpleName(t)+" "+pt.Type.Name()+" "+d.dumpRange(pt.Location))
		d.fields(v.Elem(), depth+1)
	case t.Implements(tokenType):
		tb := v.Interface().(Token).GetTokenBase()
		d.line(depth, label, javaSimpleName(t)+" "+tb.Type.Name()+" "+d.dumpRange(tb.Location))
		d.fields(v.Elem(), depth+1)
	case t == commentPtrType:
		d.line(depth, label, "Comment")
		d.fields(v.Elem(), depth+1)
	default:
		d.line(depth, label, d.scalar(v))
	}
}

func (d *dumper) fields(sv reflect.Value, depth int) {
	for _, f := range sortedStructFields(sv.Type()) {
		fv := sv.FieldByIndex(f.index)
		if !f.exported {
			// Unexported Go fields are Java private fields of basic types.
			switch fv.Kind() {
			case reflect.Bool:
				d.line(depth, f.name+"=", fmt.Sprint(fv.Bool()))
			case reflect.String:
				d.line(depth, f.name+"=", dumpQuote(fv.String()))
			case reflect.Uint32:
				if fv.Type() != reflect.TypeOf(strpool.ID(0)) {
					panic("unsupported unexported field " + sv.Type().String() + "." + f.name)
				}
				// Go-only: an interned string (IdentifierToken.value).
				d.line(depth, f.name+"=", dumpQuote(strpool.String(strpool.ID(fv.Uint()))))
			default:
				panic("unsupported unexported field " + sv.Type().String() + "." + f.name)
			}
			continue
		}
		d.render(fv, depth, f.name+"=")
	}
}

func (d *dumper) scalar(v reflect.Value) string {
	switch {
	case v.Type() == positionType:
		return d.dumpPosition(v.Interface().(SourcePosition))
	case v.Type() == rangeType:
		return d.dumpRange(v.Interface().(SourceRange))
	case v.Type() == lineRangeType:
		r := v.Interface().(LineRange)
		return dumpLinePosition(r.Start) + "-" + dumpLinePosition(r.End)
	}
	if v.Kind() == reflect.Pointer {
		// Go-only nullable scalars (*string, *enum) are non-nil here.
		v = v.Elem()
	}
	if n, ok := v.Interface().(namer); ok && v.Kind() == reflect.Int32 {
		return n.Name()
	}
	switch v.Kind() {
	case reflect.String:
		return dumpQuote(v.String())
	case reflect.Bool:
		return fmt.Sprint(v.Bool())
	case reflect.Int, reflect.Int32, reflect.Int64:
		return fmt.Sprint(v.Int())
	case reflect.Float64:
		return jutil.DoubleToString(v.Float())
	}
	return "?" + v.Type().String() + " " + dumpQuote(fmt.Sprint(v.Interface()))
}

type structField struct {
	name     string // Java name: Go name with the first letter lower-cased
	index    []int
	exported bool
}

var fieldCache sync.Map // reflect.Type -> []structField

// sortedStructFields lists the fields of struct type t below the
// ParseTreeBase/TokenBase level (their type and location are in the
// header), flattening other embedded bases (e.g. LiteralTokenBase), sorted
// by Java name. The Go-only `self` back-pointer lives in the excluded bases.
func sortedStructFields(t reflect.Type) []structField {
	if c, ok := fieldCache.Load(t); ok {
		return c.([]structField)
	}
	var out []structField
	if t == parseTreeBaseType || t == tokenBaseType {
		// A plain Token: everything is in the header.
		fieldCache.Store(t, out)
		return out
	}
	var walk func(t reflect.Type, prefix []int)
	walk = func(t reflect.Type, prefix []int) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			index := append(append([]int{}, prefix...), i)
			if f.Anonymous {
				if f.Type == parseTreeBaseType || f.Type == tokenBaseType {
					continue
				}
				walk(f.Type, index)
				continue
			}
			out = append(out, structField{
				name:     strings.ToLower(f.Name[:1]) + f.Name[1:],
				index:    index,
				exported: f.IsExported(),
			})
		}
	}
	walk(t, nil)
	sort.SliceStable(out, func(i, j int) bool { return jutil.StringCompareTo(out[i].name, out[j].name) < 0 })
	fieldCache.Store(t, out)
	return out
}

// javaSimpleName maps a Go type to the Java class simple name: base structs
// of non-abstract Java classes (TokenBase, LiteralTokenBase) are the Java
// class itself.
func javaSimpleName(t reflect.Type) string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return strings.TrimSuffix(t.Name(), "Base")
}

func (d *dumper) dumpRange(r SourceRange) string {
	return d.dumpPosition(r.Start) + "-" + d.dumpPosition(r.End)
}

func (d *dumper) dumpPosition(p SourcePosition) string {
	return dumpLinePosition(d.file.LinePosition(p))
}

func dumpLinePosition(p LinePosition) string {
	return fmt.Sprintf("%d:%d@%d", p.Line, p.Column, p.Offset)
}

// dumpQuote quotes s like ParseTreeDump.quote: \" \\ and \uXXXX for every
// UTF-16 code unit outside printable ASCII.
func dumpQuote(s string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '"' || c == '\\':
			sb.WriteByte('\\')
			sb.WriteByte(byte(c))
		case c >= 0x20 && c < 0x7f:
			sb.WriteByte(byte(c))
		default:
			fmt.Fprintf(&sb, "\\u%04X", c)
		}
	}
	sb.WriteByte('"')
	return sb.String()
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

func TestParseTreeGolden(t *testing.T) {
	files, err := filepath.Glob("testdata/ptree/*/*.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden files")
	}
	for _, f := range files {
		name := filepath.Base(f) // the Java dumps were made with the file's base name
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		want, err := readDump(f + ".ptree")
		if err != nil {
			t.Fatal(err)
		}
		got := dumpParseTree(name, string(src), configFromDump(want))
		if got != want {
			t.Errorf("%s: dump differs from Java at %s", name, firstDiff(want, got))
		}
	}
}

func TestParseTreeOracle(t *testing.T) {
	root := os.Getenv("PTREE_CORPUS")
	oracle := os.Getenv("PTREE_ORACLE_DIR")
	if root == "" || oracle == "" {
		t.Skip("set PTREE_CORPUS and PTREE_ORACLE_DIR to run the corpus oracle")
	}
	var rels []string
	if list := os.Getenv("PTREE_FILES"); list != "" {
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
	diffDir := os.Getenv("PTREE_DIFF_DIR")
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
				want, err := readDump(filepath.Join(oracle, rel+".ptree"))
				if err != nil {
					results[i] = result{rel: rel, missing: true}
					continue
				}
				src, err := os.ReadFile(filepath.Join(root, rel))
				if err != nil {
					t.Error(err)
					continue
				}
				got := dumpParseTree(rel, string(src), configFromDump(want))
				if got != want {
					results[i] = result{rel: rel, diff: firstDiff(want, got), got: got}
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
				out := filepath.Join(diffDir, r.rel+".go.ptree")
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
