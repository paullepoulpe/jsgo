package jscomp

// Go-only: test-helper — the Go side of tools/oracle/java/MessagesTypesMiscDump.java (see its
// header for the modes and the exact procedure); each dump below mirrors the Java one line for
// line, with the same escaping.
//
// Checked-in goldens (always run):
//   - testdata/xtb/expected.txt: XtbMessageBundle on testdata/xtb/*.xtb. The pinned jar's
//     XtbMessageBundle cannot run (its jarjar rules relocated javax.xml.parsers without bundling
//     it: NoClassDefFoundError), so tools/oracle/xtb_oracle.sh compiles XtbMessageBundle.java from
//     the pinned sources against the jar and runs the dump with it.
//   - testdata/messages_types_misc/ttl.txt.gz: TypeTransformation on ttl_expressions.txt.
//   - testdata/messages_types_misc/{msgs,angular,chrome}.txt.gz: the corpus dumps restricted to
//     the corpus files that report something (list in *.files).
// Full corpus (env-gated): L11MSG_CORPUS_ROOT=<corpus dir> L11MSG_CORPUS_FILES=<list>
// L11MSG_ORACLE_DIR=<dir with msgs.txt.gz, angular.txt.gz, chrome.txt.gz from the Java tool>.

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

func l11msgEscape(s string) string {
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

func l11msgEscapeNullable(s *string) string {
	if s == nil {
		return "<null>"
	}
	return l11msgEscape(*s)
}

func l11msgBool(b bool) string { return strconv.FormatBool(b) }

func l11msgDiagnostics(sb *strings.Builder, compiler *Compiler) {
	l11msgDiagnosticsFrom(sb, compiler, 0, 0)
}

// l11msgDiagnosticsFrom prints the errors and warnings reported after the
// first fromErrors / fromWarnings ones.
func l11msgDiagnosticsFrom(sb *strings.Builder, compiler *Compiler, fromErrors int, fromWarnings int) {
	for _, e := range compiler.GetErrors()[fromErrors:] {
		sb.WriteString("error " + l11msgEscape(e.String()) + "\n")
	}
	for _, e := range compiler.GetWarnings()[fromWarnings:] {
		sb.WriteString("warning " + l11msgEscape(e.String()) + "\n")
	}
}

func l11msgParts(parts []JsMessage_Part) string {
	var b strings.Builder
	b.WriteString("[")
	for i, p := range parts {
		if i > 0 {
			b.WriteString(", ")
		}
		if p.IsPlaceholder() {
			b.WriteString("PH(" + p.GetJsPlaceholderName() + "|" + p.GetCanonicalPlaceholderName() + ")")
		} else {
			b.WriteString("\"" + p.GetString() + "\"")
		}
	}
	b.WriteString("]")
	return b.String()
}

func l11msgMap(m *jutil.LinkedHashMap[string, string]) string {
	var b strings.Builder
	b.WriteString("{")
	first := true
	for k, v := range m.All() {
		if !first {
			b.WriteString(", ")
		}
		first = false
		b.WriteString(k + "=" + v)
	}
	b.WriteString("}")
	return b.String()
}

func l11msgMessage(sb *strings.Builder, m *JsMessage) {
	sb.WriteString("  key=" + l11msgEscape(m.GetKey()) + "\n")
	sb.WriteString("  id=" + l11msgEscape(m.GetId()) + "\n")
	sb.WriteString("  anonymous=" + l11msgBool(m.IsAnonymous()) + " external=" + l11msgBool(m.IsExternal()) +
		" empty=" + l11msgBool(m.IsEmpty()) + "\n")
	sb.WriteString("  source=" + l11msgEscapeNullable(m.GetSourceName()) + "\n")
	sb.WriteString("  desc=" + l11msgEscapeNullable(m.GetDesc()) + "\n")
	sb.WriteString("  meaning=" + l11msgEscapeNullable(m.GetMeaning()) + "\n")
	sb.WriteString("  alternateId=" + l11msgEscapeNullable(m.GetAlternateId()) + "\n")
	sb.WriteString("  parts=" + l11msgEscape(l11msgParts(m.GetParts())) + "\n")
	sb.WriteString("  jsPlaceholders=" + l11msgEscape(m.JsPlaceholderNames().String()) + "\n")
	sb.WriteString("  examples=" + l11msgEscape(l11msgMap(m.GetPlaceholderNameToExampleMap())) + "\n")
	sb.WriteString("  originalCode=" + l11msgEscape(l11msgMap(m.GetPlaceholderNameToOriginalCodeMap())) + "\n")
}

func l11msgNode(n rhino.Node) string {
	if n.IsNil() {
		return "<null>"
	}
	return n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

// l11msgCollectMessages is the anonymous JsMessageVisitor of dumpMessages.
type l11msgCollectMessages struct {
	JsMessageVisitorBase
	found *strings.Builder
}

func (v *l11msgCollectMessages) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	v.found.WriteString("getMsg " + l11msgNode(definition.GetMessageNode()) + " text=" +
		l11msgNode(definition.GetTemplateTextNode()) + " values=" +
		l11msgNode(definition.GetPlaceholderValuesNode()) + " valueMap=" +
		l11msgEscape(jutil.CollectionToString(definition.GetPlaceholderValueMap().KeySlice())) +
		" html=" + l11msgBool(definition.ShouldEscapeLessThan()) +
		" unescape=" + l11msgBool(definition.ShouldUnescapeHtmlEntities()) + "\n")
	l11msgMessage(v.found, definition.GetMessage())
}

func (v *l11msgCollectMessages) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	v.found.WriteString("icu " + l11msgNode(definition.GetMessageNode()) + " text=" +
		l11msgNode(definition.GetTemplateTextNode()) + "\n")
	l11msgMessage(v.found, definition.GetMessage())
}

func (v *l11msgCollectMessages) ProcessMessageFallback(callNode rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	v.found.WriteString("fallback " + l11msgNode(callNode) + " " + l11msgEscape(message1.GetKey()) + " " +
		l11msgEscape(message2.GetKey()) + "\n")
}

func l11msgDumpMessages(sb *strings.Builder, code string) {
	compiler := NewCompiler()
	root := compiler.ParseTestCode(code)
	parseErrors := len(compiler.GetErrors())
	parseWarnings := len(compiler.GetWarnings())
	found := &strings.Builder{}
	visitor := &l11msgCollectMessages{found: found}
	visitor.InitJsMessageVisitorBase(visitor, compiler, NewGoogleJsMessageIdGenerator(nil))
	visitor.Process(rhino.NilNode, root)
	sb.WriteString(found.String())
	l11msgDiagnosticsFrom(sb, compiler, parseErrors, parseWarnings)
}

func l11msgDumpXtb(sb *strings.Builder, xtb []byte) {
	var bundle *XtbMessageBundle
	failed := false
	jutil.Try(func() {
		bundle = NewXtbMessageBundle(bytes.NewReader(xtb), jutil.StringPtr("TestProject"))
	}, func(e jutil.IRuntimeException) {
		failed = true
		if re, ok := e.(*jutil.RuntimeException); ok && re.GetCause() != nil {
			sb.WriteString("exception java.lang.RuntimeException cause " + l11msgJavaClassName(re.GetCause()) + "\n")
		} else {
			sb.WriteString("exception " + l11msgJavaClassName(e) + " " + l11msgThrowableMessage(e) + "\n")
		}
	})
	if failed {
		return
	}
	for m := range bundle.GetAllMessages() {
		sb.WriteString("message " + l11msgEscape(m.GetId()) + "\n")
		l11msgMessage(sb, m)
		for _, g := range m.GetGenderedMessageVariants() {
			sb.WriteString("  gender " + g.String() + " " + l11msgEscape(l11msgParts(m.GetGenderedMessageParts(g))) + "\n")
		}
	}
}

// l11msgJavaClassName is the Java class name of an exception: the start of
// Throwable.toString(), up to ": " (or "; " for SAXParseException, whose
// toString() goes on with "; lineNumber: ..."). Java class names contain
// neither ':' nor ';'.
func l11msgJavaClassName(e jutil.Throwable) string {
	s := e.Error()
	if i := strings.IndexAny(s, ":;"); i >= 0 {
		return s[:i]
	}
	return s
}

// l11msgThrowableMessage is escape(e.getMessage()).
func l11msgThrowableMessage(e jutil.Throwable) string {
	if !strings.Contains(e.Error(), ": ") {
		return "<null>"
	}
	return l11msgEscape(e.GetMessage())
}

const l11msgExtraTypeDefs = "/** @type {number} */ Array.prototype.length;\n" +
	"/** @typedef {!Array<?>} */ var ArrayAlias;\n" +
	"\n" +
	"/** @constructor */\n" +
	"function Bar() {}\n" +
	"\n" +
	"/** @type {number} */\n" +
	"var n = 10;\n"

func l11msgRecord(registry *rhino.JSTypeRegistry, propsAndTypes ...any) rhino.JSType {
	builder := rhino.NewRecordTypeBuilder(registry)
	for i := 0; i < len(propsAndTypes); i += 2 {
		builder.AddProperty(propsAndTypes[i].(string), propsAndTypes[i+1].(rhino.JSType), rhino.NilNode)
	}
	return builder.Build()
}

// l11msgTtlErrorReporter records the TypeTransformationParser's reports.
type l11msgTtlErrorReporter struct {
	reported *strings.Builder
}

func (r *l11msgTtlErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.reported.WriteString("parser warning " + l11msgEscape(message) + "\n")
}

func (r *l11msgTtlErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.reported.WriteString("parser error " + l11msgEscape(message) + "\n")
}

func l11msgDumpTtl(sb *strings.Builder, ttlExp string) {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetCheckTypes(true)
	compiler.InitOptions(options)
	registry := compiler.GetTypeRegistry()
	reported := &strings.Builder{}
	errorReporter := &l11msgTtlErrorReporter{reported: reported}
	func() {
		closer := registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		str := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
		num := registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
		bool_ := registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
		arr := registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
		typeVars := jutil.NewImmutableMapBuilder[string, rhino.JSType]().
			Put("S", str).
			Put("N", num).
			Put("B", bool_).
			Put("BOT", registry.GetNativeType(rhino.JSTypeNative_NO_TYPE)).
			Put("TOP", registry.GetNativeType(rhino.JSTypeNative_ALL_TYPE)).
			Put("UNK", registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)).
			Put("CHKUNK", registry.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)).
			Put("SO", registry.GetNativeType(rhino.JSTypeNative_STRING_OBJECT_TYPE)).
			Put("NO", registry.GetNativeType(rhino.JSTypeNative_NUMBER_OBJECT_TYPE)).
			Put("BO", registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE)).
			Put("NULL", registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE)).
			Put("OBJ", registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)).
			Put("UNDEF", registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE)).
			Put("ARR", arr).
			Put("ARRNUM", registry.CreateTemplatizedType(arr, num)).
			Put("REC", l11msgRecord(registry, "n", num, "s", str, "b", bool_)).
			Put("NESTEDREC", l11msgRecord(registry, "s", str, "r", l11msgRecord(registry, "s", str, "b", bool_))).
			Put("ASYNCH", l11msgRecord(registry, "service", rhino.JSType(registry.CreateTemplatizedType(arr, num)))).
			BuildOrThrow()
		nameVars := jutil.NewImmutableMapBuilder[string, string]().
			Put("s", "string").
			Put("n", "number").
			Put("b", "boolean").
			BuildOrThrow()
		ttlParser := parsing.NewTypeTransformationParser(
			ttlExp, SourceFile_FromCode("[testcode]", ttlExp), errorReporter, 0, 0)
		if ttlParser.ParseTypeTransformation() {
			ast := ttlParser.GetTypeTransformationAst()
			extraTypeDefs := compiler.ParseTestCode(l11msgExtraTypeDefs)
			// TypedScopeCreator.createScope returns a TypedScope (asserted through any so that
			// this compiles whatever the declared Go result type is).
			var scopeAny any = NewTypedScopeCreator(compiler).
				CreateScope(rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(extraTypeDefs)), nil)
			scope := scopeAny.(*TypedScope)
			typeTransformation := NewTypeTransformation(compiler, scope)
			// printType() writes to System.out: capture it.
			savedOut := jutil.System_Out()
			captured := &strings.Builder{}
			jutil.System_SetOut(jutil.NewPrintStream(captured))
			func() {
				defer jutil.System_SetOut(savedOut)
				jutil.Try(func() {
					resultType := typeTransformation.EvalWithNameVars(ast, typeVars, nameVars)
					reported.WriteString("result " + l11msgEscape(resultType.String()) + "\n")
				}, func(e jutil.IRuntimeException) {
					reported.WriteString("exception " + l11msgJavaClassName(e) + " " + l11msgThrowableMessage(e) + "\n")
				})
			}()
			if captured.Len() > 0 {
				reported.WriteString("stdout " + l11msgEscape(captured.String()) + "\n")
			}
		} else {
			reported.WriteString("unparsed\n")
		}
	}()
	sb.WriteString(reported.String())
	l11msgDiagnostics(sb, compiler)
}

func l11msgDumpPass(sb *strings.Builder, mode string, code string) {
	compiler := NewCompiler()
	root := compiler.ParseTestCode(code)
	parseErrors := len(compiler.GetErrors())
	parseWarnings := len(compiler.GetWarnings())
	before := root.CloneTree()
	var pass CompilerPass
	if mode == "angular" {
		pass = NewAngularPass(compiler)
	} else {
		pass = NewChromePass(compiler)
	}
	jutil.Try(func() {
		pass.Process(rhino.NilNode, root)
	}, func(e jutil.IRuntimeException) {
		message := "<null>"
		if strings.Contains(e.Error(), ": ") {
			first, _, _ := strings.Cut(e.GetMessage(), "\n")
			message = l11msgEscape(first)
		}
		sb.WriteString("exception " + l11msgJavaClassName(e) + " " + message)
		if cause := e.GetCause(); cause != nil {
			sb.WriteString(" cause " + l11msgJavaClassName(cause))
		}
		sb.WriteString("\n")
	})
	l11msgDiagnosticsFrom(sb, compiler, parseErrors, parseWarnings)
	l11msgDiffTree(sb, before, root, "0")
}

// l11msgDiffTree is MessagesTypesMiscDump.diffTree.
func l11msgDiffTree(sb *strings.Builder, before rhino.Node, after rhino.Node, path string) {
	afterLine := after.String()
	if before.String() != afterLine {
		sb.WriteString("changed " + path + " " + l11msgEscape(afterLine) + "\n")
	}
	if before.GetChildCount() == after.GetChildCount() {
		i := 0
		for b, a := before.GetFirstChild(), after.GetFirstChild(); !b.IsNil(); b, a = b.GetNext(), a.GetNext() {
			l11msgDiffTree(sb, b, a, path+"/"+strconv.Itoa(i))
			i++
		}
		return
	}
	unmatched := map[string]int{}
	for b := before.GetFirstChild(); !b.IsNil(); b = b.GetNext() {
		unmatched[b.ToStringTree()]++
	}
	left := before.GetChildCount()
	i := 0
	for a := after.GetFirstChild(); !a.IsNil(); a = a.GetNext() {
		tree := a.ToStringTree()
		if unmatched[tree] > 0 {
			unmatched[tree]--
			left--
		} else {
			sb.WriteString("inserted " + path + "/" + strconv.Itoa(i) + "\n")
			for _, line := range strings.Split(tree, "\n") {
				sb.WriteString("| " + l11msgEscape(line) + "\n")
			}
		}
		i++
	}
	sb.WriteString("unmatched-before " + path + " " + strconv.Itoa(left) + "\n")
}

// l11msgSplitDump splits a dump into its "== name" sections.
func l11msgSplitDump(dump string) map[string]string {
	sections := map[string]string{}
	var name string
	var body strings.Builder
	flush := func() {
		if name != "" {
			sections[name] = body.String()
		}
		body.Reset()
	}
	sc := bufio.NewScanner(strings.NewReader(dump))
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "== ") {
			flush()
			name = line[3:]
			continue
		}
		body.WriteString(line + "\n")
	}
	flush()
	return sections
}

func l11msgReadMaybeGzip(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(path, ".gz") {
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		data, err = io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
	}
	return string(data)
}

func l11msgCompareSections(t *testing.T, what string, want map[string]string, got map[string]string, maxDiffs int) {
	t.Helper()
	var names []string
	for n := range want {
		names = append(names, n)
	}
	for n := range got {
		if _, ok := want[n]; !ok {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	diffs := 0
	for _, n := range names {
		if want[n] != got[n] {
			diffs++
			if diffs <= maxDiffs {
				t.Errorf("%s %s:\n got:\n%s\n Java:\n%s", what, n, got[n], want[n])
			}
		}
	}
	if diffs > 0 {
		t.Errorf("%s: %d of %d sections differ", what, diffs, len(names))
	}
}

// TestXtbMessageBundleOracle compares XtbMessageBundle on testdata/xtb with
// the pinned Java class (testdata/xtb/expected.txt).
func TestXtbMessageBundleOracle(t *testing.T) {
	files, _ := filepath.Glob("testdata/xtb/*.xtb")
	sort.Strings(files)
	got := map[string]string{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		sb := &strings.Builder{}
		l11msgDumpXtb(sb, data)
		got[filepath.Base(f)] = sb.String()
	}
	l11msgCompareSections(t, "xtb", l11msgSplitDump(l11msgReadMaybeGzip(t, "testdata/xtb/expected.txt")), got, 10)
}

// TestTypeTransformationOracle compares TypeTransformation on the TTL
// expressions of testdata/messages_types_misc/ttl_expressions.txt.
func TestTypeTransformationOracle(t *testing.T) {
	data, err := os.ReadFile("testdata/messages_types_misc/ttl_expressions.txt")
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		sb := &strings.Builder{}
		l11msgDumpTtl(sb, line)
		got[l11msgEscape(line)] = sb.String()
	}
	l11msgCompareSections(t, "ttl", l11msgSplitDump(l11msgReadMaybeGzip(t, "testdata/messages_types_misc/ttl.txt.gz")), got, 10)
}

func l11msgRunCorpusMode(t *testing.T, mode string, root string, files []string, want map[string]string, maxDiffs int) {
	t.Helper()
	got := map[string]string{}
	for _, rel := range files {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		sb := &strings.Builder{}
		jutil.Try(func() {
			if mode == "msgs" {
				l11msgDumpMessages(sb, string(data))
			} else {
				l11msgDumpPass(sb, mode, string(data))
			}
		}, func(e jutil.Throwable) {
			sb.WriteString("go panic " + l11msgEscape(e.Error()) + "\n")
		})
		if sb.Len() > 0 {
			got[l11msgEscape(rel)] = sb.String()
		}
	}
	l11msgCompareSections(t, mode, want, got, maxDiffs)
}

// l11msgReadList reads a file list, without the files of
// tools/oracle/gpl-excluded-corpus.txt.
func l11msgReadList(t *testing.T, path string) []string {
	t.Helper()
	files, err := corpus.ReadList(path)
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// TestMessagesTypesMiscSeedsGolden runs the corpus modes on the inputs of
// JsMessageVisitorTest, AngularPassTest and ChromePassTest
// (testdata/messages_types_misc/seeds/<mode>_*.js) against the Java dumps
// testdata/messages_types_misc/seeds_<mode>.txt.gz.
func TestMessagesTypesMiscSeedsGolden(t *testing.T) {
	for _, mode := range []string{"msgs", "angular", "chrome"} {
		t.Run(mode, func(t *testing.T) {
			paths, _ := filepath.Glob("testdata/messages_types_misc/seeds/" + mode + "_*.js")
			sort.Strings(paths)
			var files []string
			for _, p := range paths {
				files = append(files, filepath.Base(p))
			}
			want := l11msgSplitDump(l11msgReadMaybeGzip(t, "testdata/messages_types_misc/seeds_"+mode+".txt.gz"))
			l11msgRunCorpusMode(t, mode, "testdata/messages_types_misc/seeds", files, want, 10)
		})
	}
}

// TestMessagesTypesMiscCorpusGolden runs the msgs and chrome modes on the
// corpus files of the checked-in goldens, when the corpus is present:
// testdata/messages_types_misc/<mode>.files lists the corpus files the Java
// dump reports something for (for chrome, those of at most 20000 bytes),
// <mode>.txt.gz holds their Java sections. No corpus file reports anything in
// the angular mode (the seeds cover AngularPass).
func TestMessagesTypesMiscCorpusGolden(t *testing.T) {
	root := os.Getenv("L11MSG_CORPUS_ROOT")
	if root == "" {
		root = "../../.cache/corpus"
	}
	if _, err := os.Stat(root); err != nil {
		t.Skip("corpus not present: set L11MSG_CORPUS_ROOT")
	}
	for _, mode := range []string{"msgs", "chrome"} {
		t.Run(mode, func(t *testing.T) {
			want := l11msgSplitDump(l11msgReadMaybeGzip(t, "testdata/messages_types_misc/"+mode+".txt.gz"))
			l11msgRunCorpusMode(t, mode, root, l11msgReadList(t, "testdata/messages_types_misc/"+mode+".files"), want, 10)
		})
	}
}

// TestMessagesTypesMiscCorpusOracle runs the corpus modes on a whole corpus
// list against full Java dumps (env-gated).
func TestMessagesTypesMiscCorpusOracle(t *testing.T) {
	root, list, oracleDir := os.Getenv("L11MSG_CORPUS_ROOT"), os.Getenv("L11MSG_CORPUS_FILES"), os.Getenv("L11MSG_ORACLE_DIR")
	if root == "" || list == "" || oracleDir == "" {
		t.Skip("set L11MSG_CORPUS_ROOT, L11MSG_CORPUS_FILES and L11MSG_ORACLE_DIR")
	}
	files := l11msgReadList(t, list)
	for _, mode := range []string{"msgs", "angular", "chrome"} {
		t.Run(mode, func(t *testing.T) {
			want := l11msgSplitDump(l11msgReadMaybeGzip(t, filepath.Join(oracleDir, mode+".txt.gz")))
			l11msgRunCorpusMode(t, mode, root, files, want, 20)
		})
	}
}
