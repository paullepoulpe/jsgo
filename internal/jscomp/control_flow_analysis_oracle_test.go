package jscomp

// Differential test against tools/oracle/java/ControlFlowAnalysisDump.java:
// the same sweep of the ControlFlowAnalysis / ControlFlowGraph API over
// parsed files, printed in the same format, compared section by section.
//
//   - TestControlFlowAnalysisGolden replays the checked-in golden dump
//     testdata/cfa/golden.txt.gz: the hand-written inputs testdata/cfa/js/cfa-*.js
//     followed by the corpus files of testdata/cfa/golden-files.txt (every
//     120th file of corpus-files.txt, plus files whose finallyMap holds
//     several values for a key). Regenerate with
//       CP=.cache/closure-compiler.jar:.cache/oracle-build
//       java -Xss1g -cp $CP com.google.javascript.jscomp.ControlFlowAnalysisDump a.gz \
//         internal/jscomp/testdata/cfa/js --list=<the cfa-*.js names>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.ControlFlowAnalysisDump b.gz \
//         .cache/corpus --list=<the corpus names>
//       cat a.gz b.gz > internal/jscomp/testdata/cfa/golden.txt.gz
//   - TestControlFlowAnalysisOracle compares any other dump:
//     CFA_ORACLE=<dump.gz> [CFA_ROOT=<root>] (default root: the
//     corpus under /home/user/jsgo/.cache/corpus). The held-out corpus
//     (.cache/corpus-holdout, corpus-holdout-files.txt) is for a final check
//     only.

import (
	"bufio"
	"compress/gzip"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// cfaEscape is ControlFlowAnalysisDump.esc (ParsingDump's escaping).
func cfaEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			const hex = "0123456789abcdef"
			b.WriteString("\\u")
			b.WriteByte(hex[c>>12&0xf])
			b.WriteByte(hex[c>>8&0xf])
			b.WriteByte(hex[c>>4&0xf])
			b.WriteByte(hex[c&0xf])
		}
	}
	return b.String()
}

// cfaJavaClass is Throwable.getClass().getName().
func cfaJavaClass(t jutil.Throwable) string {
	full := t.Error()
	if msg := t.GetMessage(); strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg)
	}
	return full
}

// cfaJavaMessage is String.valueOf(Throwable.getMessage()).
func cfaJavaMessage(t jutil.Throwable) string {
	if t.Error() == cfaJavaClass(t) && t.GetMessage() == "" {
		return "null"
	}
	return t.GetMessage()
}

// cfaJavaSimpleClass is Throwable.getClass().getSimpleName().
func cfaJavaSimpleClass(t jutil.Throwable) string {
	c := cfaJavaClass(t)
	if i := strings.LastIndexAny(c, ".$"); i >= 0 {
		return c[i+1:]
	}
	return c
}

type cfaErrorReporter struct{ hasError bool }

func (r *cfaErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {}
func (r *cfaErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.hasError = true
}

type cfaDumper struct {
	compiler AbstractCompiler
	config   *parsing.Config
	sb       strings.Builder
	ids      map[*rhino.Node]int
}

func (d *cfaDumper) line(s string) {
	d.sb.WriteString(cfaEscape(s))
	d.sb.WriteByte('\n')
}

func (d *cfaDumper) id(n *rhino.Node) string {
	if n == nil {
		return "R"
	}
	if i, ok := d.ids[n]; ok {
		return strconv.Itoa(i)
	}
	return "?"
}

// newCfaDumper mirrors ControlFlowAnalysisDump.main's setup.
func newCfaDumper() *cfaDumper {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	return &cfaDumper{
		compiler: newTestCompiler(options),
		// Compiler.getParserConfig(ConfigContext.DEFAULT) for these options
		// (see internal/parsing/parsingdump_test.go).
		config: parsing.ParserRunner_CreateConfigWithModes(
			parsing.Config_LanguageMode_ES_NEXT, parsing.Config_JsDocParsing_TYPES_ONLY,
			parsing.Config_RunMode_STOP_AFTER_ERROR, nil, true, parsing.Config_StrictMode_STRICT),
	}
}

// dumpOne mirrors ControlFlowAnalysisDump.dumpOne.
func (d *cfaDumper) dumpOne(name, src string) string {
	d.sb.Reset()
	d.line("## file " + name)
	reporter := &cfaErrorReporter{}
	var result *parsing.ParserRunner_ParseResult
	var exception jutil.Throwable
	jutil.Try(func() {
		result = parsing.ParserRunner_Parse(SourceFile_FromCode(name, src), src, d.config, reporter)
	}, func(e jutil.IRuntimeException) { exception = e })
	if exception != nil {
		d.line("# parse-exception")
		return d.sb.String()
	}
	if reporter.hasError || result.Ast == nil {
		d.line("# parse-error")
		return d.sb.String()
	}
	script := result.Ast
	externsRoot := rhino.IR_Root()
	jsRoot := rhino.IR_Root(script)
	globalRoot := rhino.IR_Root(externsRoot, jsRoot)

	d.ids = map[*rhino.Node]int{}
	var preorder []*rhino.Node
	var number func(n *rhino.Node)
	number = func(n *rhino.Node) {
		d.ids[n] = len(preorder)
		preorder = append(preorder, n)
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			number(c)
		}
	}
	number(globalRoot)
	d.line("# nodes " + strconv.Itoa(len(preorder)))

	// The static API.
	for i := 3; i < len(preorder); i++ {
		n := preorder[i]
		var follow, handler string
		jutil.Try(func() {
			follow = d.id(ControlFlowAnalysis_ComputeFollowNode(n))
		}, func(e jutil.IRuntimeException) { follow = "X" + cfaJavaSimpleClass(e) })
		jutil.Try(func() {
			handler = d.id(ControlFlowAnalysis_GetExceptionHandler(n))
		}, func(e jutil.IRuntimeException) { handler = "X" + cfaJavaSimpleClass(e) })
		mayThrow, entering := "0", "0"
		if ControlFlowAnalysis_MayThrowException(n) {
			mayThrow = "1"
		}
		if ControlFlowGraph_IsEnteringNewCfgNode(n) {
			entering = "1"
		}
		d.line("S " + strconv.Itoa(i) + " " + mayThrow + " " + entering + " " + follow + " " + handler)
	}

	// The CFGs.
	for _, n := range preorder {
		if NodeUtil_IsValidCfgRoot(n) {
			d.dumpCfg(n, "A", false, true)
		}
	}
	d.dumpCfg(jsRoot, "F", true, false)
	return d.sb.String()
}

// dumpCfg mirrors ControlFlowAnalysisDump.dumpCfg: it runs exactly what
// Builder.computeCfg runs, keeping the analysis to look at its finallyMap.
func (d *cfaDumper) dumpCfg(root *rhino.Node, mode string, traverseFunctions, edgeAnnotations bool) {
	d.line("cfg " + d.id(root) + " " + mode)
	cfa := newControlFlowAnalysis(d.compiler, traverseFunctions, edgeAnnotations)
	var exception jutil.Throwable
	jutil.Try(func() { cfa.computeCfg(root) }, func(e jutil.IRuntimeException) { exception = e })
	if exception != nil {
		d.line("exception " + cfaJavaClass(exception) + ": " + cfaJavaMessage(exception))
		if cause := exception.GetCause(); cause != nil {
			d.line("cause " + cfaJavaClass(cause) + ": " + cfaJavaMessage(cause))
		}
		return
	}
	cfg := cfa.GetCfg()
	for _, gn := range cfg.GetDiGraphNodes() {
		v := gn.GetValue()
		var desc string
		if v == nil {
			desc = "N R - -"
		} else {
			desc = "N " + d.id(v) + " " + v.GetToken().String() + " " +
				strconv.Itoa(v.GetLineno()) + ":" + strconv.Itoa(v.GetCharno())
		}
		d.line(desc + " p=" + strconv.Itoa(gn.GetPriority()))
		for _, e := range gn.GetOutEdges() {
			d.line(" > " + e.GetValue().String() + " " + d.id(e.GetDestination().GetValue()))
		}
		for _, e := range gn.GetInEdges() {
			d.line(" < " + e.GetValue().String() + " " + d.id(e.GetSource().GetValue()))
		}
	}
	for k := range cfa.finallyMap.KeySet() {
		if cfa.finallyMap.Get(k).Size() > 1 {
			d.line("multi")
			break
		}
	}
}

// cfaReadSections splits a gzipped dump into its "## file" sections.
func cfaReadSections(t *testing.T, path string) (names []string, sections map[string]string) {
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
	sections = map[string]string{}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<28)
	var cur string
	var sb strings.Builder
	flush := func() {
		if cur != "" {
			sections[cur] = sb.String()
		}
		sb.Reset()
	}
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "## file ") {
			flush()
			cur = l[len("## file "):]
			names = append(names, cur)
		}
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	return names, sections
}

func cfaFirstDiff(want, got string) string {
	wl, gl := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return "line " + strconv.Itoa(i+1) + ":\n  java: " + w + "\n  go:   " + g
		}
	}
	return ""
}

// cfaCompare dumps every listed file and compares it with the Java dump.
func cfaCompare(t *testing.T, dumpPath string, rootFor func(name string) string) {
	names, sections := cfaReadSections(t, dumpPath)
	if len(names) == 0 {
		t.Fatalf("no sections in %s", dumpPath)
	}
	d := newCfaDumper()
	failures := 0
	for _, name := range names {
		src, err := os.ReadFile(filepath.Join(rootFor(name), name))
		if err != nil {
			t.Fatal(err)
		}
		got := d.dumpOne(name, string(src))
		if want := sections[name]; got != want {
			failures++
			if failures <= 20 {
				t.Errorf("%s: differs from Java at %s", name, cfaFirstDiff(want, got))
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, len(names))
	}
}

func TestControlFlowAnalysisGolden(t *testing.T) {
	cfaCompare(t, "testdata/cfa/golden.txt.gz", func(name string) string {
		if strings.HasPrefix(name, "cfa-") {
			return "testdata/cfa/js"
		}
		root := os.Getenv("CFA_CORPUS")
		if root == "" {
			root = "/home/user/jsgo/.cache/corpus"
		}
		if _, err := os.Stat(root); err != nil {
			t.Skipf("corpus not available (%v)", err)
		}
		return root
	})
}

func TestControlFlowAnalysisOracle(t *testing.T) {
	dump := os.Getenv("CFA_ORACLE")
	if dump == "" {
		t.Skip("set CFA_ORACLE=<ControlFlowAnalysisDump output .gz> (and CFA_ROOT)")
	}
	root := os.Getenv("CFA_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	cfaCompare(t, dump, func(string) string { return root })
}
