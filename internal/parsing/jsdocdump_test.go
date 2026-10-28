package parsing

// Differential JSDoc oracle, Go side. dumpJsDoc reproduces, byte for byte,
// the format of tools/oracle/java/JsDocDump.java: every JSDoc and important
// comment that internal/parser collects is fed to JsDocInfoParser the way
// IRFactory does (createJsDocInfoParser, handlePossibleFileOverviewJsDoc,
// parseJSDocInfoFrom, parseInlineTypeDoc, setFileOverviewJsDoc — mirrored
// below, since IRFactory is ported separately), and the resulting JSDocInfo,
// markers, license text, file-overview state and diagnostics are printed.
//
//   - TestJsDocDumpGolden compares against the checked-in Java dumps in
//     testdata/jsdoc/<mode>/ for the inputs in testdata/jsdoc/src/ (always
//     runs; the dumps are gzipped). <mode> is the JsDocParsing mode the Java side was run with:
//     types (the compiler default, TYPES_ONLY), docs (--jsdoc,
//     INCLUDE_DESCRIPTIONS_NO_WHITESPACE), whitespace
//     (INCLUDE_DESCRIPTIONS_WITH_WHITESPACE), all (INCLUDE_ALL_COMMENTS) and
//     license (LICENSE_COMMENTS_ONLY). Regenerate a mode with
//     (cd testdata/jsdoc/src && java -Xss1g -cp $CP
//     com.google.javascript.jscomp.parsing.JsDocDump --gzip --jsdoc=<MODE>
//     ../<mode> . *.js). The inputs are hand-written edge cases (every
//     annotation, malformed types, unicode, line endings, unterminated
//     comments, file overviews, TTL) plus every comment string of
//     JsDocInfoParserTest.java (jsdocinfoparsertest_*.js).
//   - TestJsDocDumpOracle compares a whole corpus when JSDOC_CORPUS (the
//     corpus root) and JSDOC_ORACLE_DIR (the JsDocDump output for it, .jsdoc
//     or .jsdoc.gz) are set; JSDOC_FILES optionally names a file listing the
//     relative paths (default: every .js/.mjs/.cjs under the root), and
//     JSDOC_DIFF_DIR, if set, receives the Go dump of every mismatching file.
//
// The configuration is read back from the "# config", "# annotations",
// "# suppressions" and "# primitives" header lines of the Java dump and
// rebuilt with Config_Builder_; the Go dump prints the rebuilt Config's own
// values, so a wrong reconstruction (or a wrong Config.buildAnnotations)
// shows up as a diff on those lines.

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

// silentParserReporter drops the parser's own diagnostics (JsDocDump does).
type silentParserReporter struct {
	parser.ErrorReporterBase
}

func newSilentParserReporter() *silentParserReporter {
	r := &silentParserReporter{}
	r.InitErrorReporterBase(r)
	return r
}

func (r *silentParserReporter) ReportError(location parser.LinePosition, message string)   {}
func (r *silentParserReporter) ReportWarning(location parser.LinePosition, message string) {}

// jsdocRecordingReporter is JsDocDump.RecordingReporter.
type jsdocRecordingReporter struct {
	lines []string
}

func (r *jsdocRecordingReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.lines = append(r.lines, fmt.Sprintf("WARNING %d:%d %s", line, lineOffset, jsdocQuote(message)))
}

func (r *jsdocRecordingReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.lines = append(r.lines, fmt.Sprintf("ERROR %d:%d %s", line, lineOffset, jsdocQuote(message)))
}

// jsdocFileState is JsDocDump.FileState: the per-file state IRFactory keeps
// for file-level JSDoc.
type jsdocFileState struct {
	licenseBuilder    *jutil.LinkedHashSet[string]
	firstFileoverview *rhino.JSDocInfo
}

// testNewParserConfig mirrors ParserRunner.newParserConfig(Config) (ported
// separately, with ParserRunner).
func testNewParserConfig(config *Config) *parser.Parser_Config {
	isStrictMode := config.StrictMode().IsStrict()
	var mode parser.Parser_Config_Mode
	switch config.LanguageMode() {
	case Config_LanguageMode_ECMASCRIPT3:
		mode = parser.Parser_Config_Mode_ES3
	case Config_LanguageMode_ECMASCRIPT5:
		mode = parser.Parser_Config_Mode_ES5
	case Config_LanguageMode_ECMASCRIPT_2015, Config_LanguageMode_ECMASCRIPT_2016:
		mode = parser.Parser_Config_Mode_ES6_OR_ES7
	default:
		mode = parser.Parser_Config_Mode_ES8_OR_GREATER
	}
	return parser.NewParser_ConfigWithMode(mode, isStrictMode)
}

// jsdocConfigFromDump rebuilds the parsing.Config described by the header of
// a Java dump.
func jsdocConfigFromDump(dump string) (*Config, error) {
	header := map[string]string{}
	for _, l := range strings.SplitN(dump, "\n", 5)[:4] {
		k, v, ok := strings.Cut(l, " ")
		if !ok || k != "#" {
			return nil, fmt.Errorf("bad header line %q", l)
		}
		k, v, _ = strings.Cut(v, " ")
		header[k] = jsdocUnescape(v)
	}
	f := strings.Fields(header["config"])
	if len(f) != 5 {
		return nil, fmt.Errorf("bad config line %q", header["config"])
	}
	b := Config_Builder_()
	for _, m := range Config_LanguageMode_values() {
		if m.Name() == f[0] {
			b.SetLanguageMode(m)
		}
	}
	for _, m := range Config_StrictMode_values() {
		if m.Name() == f[1] {
			b.SetStrictMode(m)
		}
	}
	for _, m := range Config_JsDocParsing_values() {
		if m.Name() == f[2] {
			b.SetJsDocParsingMode(m)
		}
	}
	for _, m := range Config_RunMode_values() {
		if m.Name() == f[3] {
			b.SetRunMode(m)
		}
	}
	b.SetParseInlineSourceMaps(f[4] == "parseInlineSourceMaps=true")
	var annotationNames []string
	for _, kv := range jsdocSplitList(header["annotations"], "{", "}") {
		k, _, _ := strings.Cut(kv, "=")
		annotationNames = append(annotationNames, k)
	}
	b.SetExtraAnnotationNames(jutil.NewLinkedHashSet(annotationNames...).All())
	b.SetSuppressionNames(jutil.NewLinkedHashSet(jsdocSplitList(header["suppressions"], "[", "]")...).All())
	b.SetClosurePrimitiveNames(jutil.NewLinkedHashSet(jsdocSplitList(header["primitives"], "[", "]")...).All())
	return b.Build(), nil
}

// jsdocSplitList splits a Java collection rendering "[a, b]" / "{a, b}".
func jsdocSplitList(s, open, close string) []string {
	s = strings.TrimSuffix(strings.TrimPrefix(s, open), close)
	if s == "" {
		return nil
	}
	return strings.Split(s, ", ")
}

// jsdocUnescape undoes jsdocEscape for the header lines (which are ASCII in
// practice).
func jsdocUnescape(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var units []uint16
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == '\\' {
			units = append(units, '\\')
			i++
		} else if s[i] == '\\' && i+5 < len(s) && s[i+1] == 'u' {
			var v uint16
			fmt.Sscanf(s[i+2:i+6], "%04x", &v)
			units = append(units, v)
			i += 5
		} else {
			units = append(units, uint16(s[i]))
		}
	}
	return jutil.StringFromUTF16(units)
}

// dumpJsDoc is JsDocDump.dumpOne.
func dumpJsDoc(name, src string, config *Config) string {
	var sb strings.Builder
	line := func(s string) {
		sb.WriteString(jsdocEscape(s))
		sb.WriteByte('\n')
	}
	line("# config " + config.LanguageMode().String() + " " + config.StrictMode().String() + " " +
		config.JsDocParsingMode().String() + " " + config.RunMode().String() +
		" parseInlineSourceMaps=" + jutil.ToJavaString(config.ParseInlineSourceMaps()))
	line("# annotations " + config.Annotations().String())
	line("# suppressions " + config.SuppressionNames().String())
	line("# primitives " + config.ClosurePrimitiveNames().String())

	file := parser.NewSourceFile(name, src)
	p := parser.NewParser(testNewParserConfig(config), newSilentParserReporter(), file)
	if exception := jsdocCatch(func() { p.ParseProgram() }); exception != "" {
		line("# parse-exception " + exception)
	}
	comments := p.GetComments()

	sourceFile := rhino.NewSimpleSourceFile(name, rhino.StaticSourceFile_SourceKind_STRONG)
	templateNode := rhino.NewNode(rhino.Token_SCRIPT)
	templateNode.SetStaticSourceFile(sourceFile)
	state := &jsdocFileState{licenseBuilder: jutil.NewLinkedHashSet[string]()}

	for i := range comments {
		c := &comments[i]
		if c.Type != parser.Comment_Type_JSDOC && c.Type != parser.Comment_Type_IMPORTANT {
			continue
		}
		line("# comment " + c.Type.String() + " " + jsdocRange(c.Location))
		reporter := &jsdocRecordingReporter{}
		flush := func(r *jsdocRecordingReporter) {
			for _, l := range r.lines {
				line(l)
			}
			r.lines = nil
		}
		exception := jsdocCatch(func() {
			jsdocParser := testCreateJsDocInfoParser(c, templateNode, config, reporter, state)
			isFileOverview := testHandlePossibleFileOverviewJsDoc(jsdocParser, state)
			var info *rhino.JSDocInfo
			if !isFileOverview {
				info = jsdocParser.RetrieveAndResetParsedJSDocInfo()
			}
			flush(reporter)
			line("info " + jsdocRender(info))
			jsdocMarkers(line, info)
			license := jsdocParser.GetLicenseText()
			if license == nil {
				line("license null")
			} else {
				line("license " + jsdocQuote(*license))
			}
			if isFileOverview {
				line("fileoverview " + jsdocRender(state.firstFileoverview))
				jsdocMarkers(line, state.firstFileoverview)
			} else {
				line("fileoverview -")
			}
		})
		if exception != "" {
			flush(reporter)
			line("exception " + exception)
		}
		if c.Type == parser.Comment_Type_JSDOC && !strings.Contains(c.Value, "@") {
			inlineReporter := &jsdocRecordingReporter{}
			exception := jsdocCatch(func() {
				inline := testParseInlineTypeDoc(c, templateNode, config, inlineReporter)
				line("inline")
				flush(inlineReporter)
				line("inline " + jsdocRender(inline))
				jsdocMarkers(line, inline)
			})
			if exception != "" {
				line("inline")
				flush(inlineReporter)
				line("exception " + exception)
			}
		}
	}

	// IRFactory.setFileOverviewJsDoc
	var fileoverview *rhino.JSDocInfo_Builder
	if state.firstFileoverview == nil {
		fileoverview = rhino.JSDocInfo_Builder_()
	} else {
		fileoverview = state.firstFileoverview.ToBuilder()
	}
	if !state.licenseBuilder.IsEmpty() {
		fileoverview.RecordLicense(strings.Join(state.licenseBuilder.ToSlice(), ""))
	}
	fileInfo := fileoverview.BuildWithAlways(false)
	line("# file " + jsdocRender(fileInfo))
	jsdocMarkers(line, fileInfo)
	return sb.String()
}

// testCreateJsDocInfoParser mirrors IRFactory.createJsDocInfoParser(Comment,
// boolean) with useLicensesOnlyConfig = false.
func testCreateJsDocInfoParser(
	node *parser.Comment,
	templateNode rhino.Node,
	config *Config,
	errorReporter rhino.ErrorReporter,
	state *jsdocFileState,
) *JsDocInfoParser {
	comment := node.Value
	lineno := int(node.Location.Start.Line) + 1
	charno := int(node.Location.Start.Column)
	position := int(node.Location.Start.Offset)

	// The JsDocInfoParser expects the comment without the initial '/**'.
	numOpeningChars := 3
	jsdocParser := NewJsDocInfoParser(
		NewJsDocTokenStreamWithLinenoCharno(
			jutil.StringSubstringFrom(comment, numOpeningChars), lineno, charno+numOpeningChars),
		comment,
		position,
		templateNode,
		config,
		JsDocInfoParser_JsDocSourceKind_NORMAL,
		errorReporter)
	jsdocParser.SetFileOverviewJSDocInfo(state.firstFileoverview)
	if node.Type == parser.Comment_Type_IMPORTANT && len(node.Value) > 0 {
		jsdocParser.ParseImportantComment()
	} else {
		jsdocParser.Parse()
	}

	return jsdocParser
}

// testHandlePossibleFileOverviewJsDoc mirrors
// IRFactory.handlePossibleFileOverviewJsDoc(JsDocInfoParser).
func testHandlePossibleFileOverviewJsDoc(jsDocParser *JsDocInfoParser, state *jsdocFileState) bool {
	if jsDocParser.GetLicenseText() != nil {
		state.licenseBuilder.Add(*jsDocParser.GetLicenseText())
	}

	newFileoverview := jsDocParser.GetFileOverviewJSDocInfo()
	if newFileoverview == state.firstFileoverview {
		return false
	}

	if state.firstFileoverview == nil {
		state.firstFileoverview = newFileoverview
	} else {
		merged := state.firstFileoverview.ToBuilder()
		if !newFileoverview.GetSuppressions().IsEmpty() {
			merged.RecordSuppressions(newFileoverview.GetSuppressions())
		}
		if newFileoverview.IsExterns() {
			merged.RecordExterns()
		}
		if newFileoverview.IsNoCoverage() {
			merged.RecordNoCoverage()
		}
		state.firstFileoverview = merged.Build()
	}

	return true
}

// testParseInlineTypeDoc mirrors IRFactory.parseInlineTypeDoc(Comment).
func testParseInlineTypeDoc(
	node *parser.Comment, templateNode rhino.Node, config *Config, errorReporter rhino.ErrorReporter,
) *rhino.JSDocInfo {
	comment := node.Value
	lineno := int(node.Location.Start.Line) + 1
	charno := int(node.Location.Start.Column)

	// The JsDocInfoParser expects the comment without the initial '/**'.
	numOpeningChars := 3
	p := NewJsDocInfoParser(
		NewJsDocTokenStreamWithLinenoCharno(
			jutil.StringSubstringFrom(comment, numOpeningChars), lineno, charno+numOpeningChars),
		comment,
		int(node.Location.Start.Offset),
		templateNode,
		config,
		JsDocInfoParser_JsDocSourceKind_NORMAL,
		errorReporter)
	return p.ParseInlineTypeDoc()
}

var jsdocMarkerHash = regexp.MustCompile(`com\.google\.javascript\.rhino\.JSDocInfo\$Marker@[0-9a-f]+`)

func jsdocRender(info *rhino.JSDocInfo) string {
	if info == nil {
		return "null"
	}
	return jsdocMarkerHash.ReplaceAllLiteralString(info.ToStringVerbose(), "JSDocInfo$Marker")
}

func jsdocMarkers(line func(string), info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	markers := info.GetMarkers()
	for i := range markers {
		m := &markers[i]
		type_ := "null"
		if m.GetType() != nil {
			type_ = jsdocPos(m.GetType(), m.GetType().GetItem().ToStringTree())
			if m.GetType().HasBrackets() {
				type_ += "[brackets]"
			}
		}
		annotation := "null"
		if a := m.GetAnnotation(); a != nil {
			annotation = jsdocPos(a, jsdocQuote(a.GetItem()))
		}
		name := "null"
		if n := m.GetNameNode(); n != nil {
			name = jsdocPos(n, n.GetItem().ToStringTree())
		}
		description := "null"
		if d := m.GetDescription(); d != nil {
			description = jsdocPos(d, jsdocQuote(d.GetItem()))
		}
		line("marker annotation=" + annotation + " name=" + name + " description=" + description + " type=" + type_)
	}
}

type jsdocSourcePosition interface {
	GetStartLine() int
	GetPositionOnStartLine() int
	GetEndLine() int
	GetPositionOnEndLine() int
}

func jsdocPos(p jsdocSourcePosition, item string) string {
	return fmt.Sprintf("%s@%d:%d-%d:%d", item, p.GetStartLine(), p.GetPositionOnStartLine(), p.GetEndLine(), p.GetPositionOnEndLine())
}

func jsdocRange(r parser.LineRange) string {
	return jsdocPosition(r.Start) + "-" + jsdocPosition(r.End)
}

func jsdocPosition(p parser.LinePosition) string {
	return fmt.Sprintf("%d:%d@%d", p.Line, p.Column, p.Offset)
}

// jsdocQuote is JsDocDump.quote: escaping happens once, per line.
func jsdocQuote(s string) string { return `"` + s + `"` }

// jsdocEscape is JsDocDump.escape: backslash and every UTF-16 char outside
// printable ASCII become escapes.
func jsdocEscape(s string) string {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' || s[i] < 0x20 || s[i] >= 0x7f {
			ascii = false
			break
		}
	}
	if ascii {
		return s
	}
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString(`\\`)
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, `\u%04x`, c)
		}
	}
	return b.String()
}

// jsdocCatch runs f and renders a panic like the Java side renders a caught
// RuntimeException: "<class>: <message>".
func jsdocCatch(f func()) (exception string) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				full := t.Error()
				if !strings.Contains(full, ": ") {
					full += ": null"
				}
				exception = full
				return
			}
			buf := make([]byte, 1<<16)
			buf = buf[:runtime.Stack(buf, false)]
			exception = "go-panic: " + fmt.Sprint(r) + "\n" + string(buf)
		}
	}()
	f()
	return ""
}

func readJsDocDump(path string) (string, error) {
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

// jsdocSections splits a dump into its header, its "# comment" sections and
// its "# file" trailer.
func jsdocSections(dump string) []string {
	var out []string
	start := 0
	for i := 0; i < len(dump); {
		j := strings.IndexByte(dump[i:], '\n')
		if j < 0 {
			break
		}
		l := dump[i : i+j]
		if (strings.HasPrefix(l, "# comment ") || strings.HasPrefix(l, "# file ")) && i > start {
			out = append(out, dump[start:i])
			start = i
		}
		i += j + 1
	}
	return append(out, dump[start:])
}

// jsdocCompare compares a Go dump with a Java dump section by section.
// It returns the number of comment sections in the Java dump, how many of
// them are identical, and the first differing section pair.
func jsdocCompare(want, got string) (comments, same int, firstDiff string) {
	ws, gs := jsdocSections(want), jsdocSections(got)
	for i, w := range ws {
		isComment := strings.HasPrefix(w, "# comment ")
		if isComment {
			comments++
		}
		var g string
		if i < len(gs) {
			g = gs[i]
		}
		if w == g {
			if isComment {
				same++
			}
			continue
		}
		if firstDiff == "" {
			firstDiff = "want:\n" + w + "got:\n" + g
		}
	}
	if len(gs) != len(ws) && firstDiff == "" {
		firstDiff = fmt.Sprintf("want %d sections, got %d", len(ws), len(gs))
	}
	return
}

func TestJsDocDumpGolden(t *testing.T) {
	srcDir := filepath.Join("testdata", "jsdoc", "src")
	modes, err := os.ReadDir(filepath.Join("testdata", "jsdoc"))
	if err != nil {
		t.Fatal(err)
	}
	total, totalSame := 0, 0
	for _, mode := range modes {
		if !mode.IsDir() || mode.Name() == "src" {
			continue
		}
		dir := filepath.Join("testdata", "jsdoc", mode.Name())
		dumps, err := filepath.Glob(filepath.Join(dir, "*.jsdoc.gz"))
		if err != nil {
			t.Fatal(err)
		}
		if len(dumps) == 0 {
			t.Fatalf("no goldens in %s", dir)
		}
		for _, d := range dumps {
			rel := strings.TrimSuffix(filepath.Base(d), ".jsdoc.gz")
			want, err := readJsDocDump(strings.TrimSuffix(d, ".gz"))
			if err != nil {
				t.Fatal(err)
			}
			src, err := os.ReadFile(filepath.Join(srcDir, rel))
			if err != nil {
				t.Fatal(err)
			}
			config, err := jsdocConfigFromDump(string(want))
			if err != nil {
				t.Fatalf("%s: %v", d, err)
			}
			got := dumpJsDoc(rel, string(src), config)
			comments, same, diff := jsdocCompare(string(want), got)
			total += comments
			totalSame += same
			if diff != "" {
				t.Errorf("%s/%s: %s", mode.Name(), rel, diff)
			}
		}
	}
	t.Logf("%d comments, %d identical", total, totalSame)
}

func TestJsDocDumpOracle(t *testing.T) {
	root := os.Getenv("JSDOC_CORPUS")
	oracle := os.Getenv("JSDOC_ORACLE_DIR")
	if root == "" || oracle == "" {
		t.Skip("set JSDOC_CORPUS and JSDOC_ORACLE_DIR to run the corpus oracle")
	}
	var rels []string
	if list := os.Getenv("JSDOC_FILES"); list != "" {
		f, err := os.Open(list)
		if err != nil {
			t.Fatal(err)
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
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
	diffDir := os.Getenv("JSDOC_DIFF_DIR")
	type result struct {
		rel, diff, got string
		missing        bool
		comments, same int
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
				want, err := readJsDocDump(filepath.Join(oracle, rel+".jsdoc"))
				if err != nil {
					results[i] = result{rel: rel, missing: true}
					continue
				}
				src, err := os.ReadFile(filepath.Join(root, rel))
				if err != nil {
					t.Error(err)
					continue
				}
				config, err := jsdocConfigFromDump(want)
				if err != nil {
					t.Errorf("%s: %v", rel, err)
					continue
				}
				got := dumpJsDoc(rel, string(src), config)
				comments, same, diff := jsdocCompare(want, got)
				results[i] = result{rel: rel, diff: diff, got: got, comments: comments, same: same}
			}
		}()
	}
	for i := range rels {
		work <- i
	}
	close(work)
	wg.Wait()
	files, missing, differ, comments, same := 0, 0, 0, 0, 0
	for _, r := range results {
		if r.missing {
			missing++
			continue
		}
		files++
		comments += r.comments
		same += r.same
		if r.diff != "" {
			differ++
			if differ <= 20 {
				t.Errorf("%s:\n%s", r.rel, r.diff)
			}
			if diffDir != "" {
				out := filepath.Join(diffDir, r.rel+".go.jsdoc")
				os.MkdirAll(filepath.Dir(out), 0o755)
				os.WriteFile(out, []byte(r.got), 0o644)
			}
		}
	}
	t.Logf("%d files (%d differ, %d without oracle output); %d comments: %d identical, %d differ",
		files, differ, missing, comments, same, comments-same)
}
