package parsing

// Go-only: test-helper — differential check of internal/rhino's type system
// on the type expressions real code writes. The Java side is
// tools/oracle/java/CorpusTypeDump.java; this mirrors it call for call:
// parse with ParserRunner_Parse, create a fresh JSTypeRegistry per file, walk
// the AST in pre-order and evaluate every JSDocInfo type expression with
// EvaluateTypeExpressionInGlobalScope.
//
//	TYPES_ROOT=<corpus root> TYPES_ORACLE_DIR=<CorpusTypeDump output>
//	TYPES_FILES=<file list> [TYPES_DIFF_DIR=<dir>] go test -run TestCorpusTypesOracle

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

type typeLineReporter struct{ line func(string) }

func (r typeLineReporter) Warning(message, sourceName string, lineno, lineOffset int) {
	r.line("WARNING " + strconv.Itoa(lineno) + ":" + strconv.Itoa(lineOffset) + " " + message)
}

func (r typeLineReporter) Error(message, sourceName string, lineno, lineOffset int) {
	r.line("ERROR " + strconv.Itoa(lineno) + ":" + strconv.Itoa(lineOffset) + " " + message)
}

func dumpCorpusTypes(name, src string) string {
	var sb strings.Builder
	line := func(s string) {
		sb.WriteString(dumpEscape(s))
		sb.WriteByte('\n')
	}
	result := ParserRunner_Parse(rhino.NewSimpleSourceFile(name, rhino.StaticSourceFile_SourceKind_STRONG),
		src, defaultDumpConfig(), rhino.ErrorReporter_NULL_INSTANCE)
	if result.Ast.IsNil() {
		line("# no ast")
		return sb.String()
	}
	registry := rhino.NewJSTypeRegistry(typeLineReporter{line})
	number := registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	str := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	object := registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	b := strconv.FormatBool
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		if info := n.GetJSDocInfo(); info != nil {
			for _, expr := range info.GetTypeExpressions() {
				root := expr.GetRoot()
				where := strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno()) +
					" @" + strconv.Itoa(root.GetLineno()) + ":" + strconv.Itoa(root.GetCharno())
				var t rhino.JSType
				var exc jutil.Throwable
				jutil.Try(func() { t = registry.EvaluateTypeExpressionInGlobalScope(expr) },
					func(e jutil.IRuntimeException) { exc = e })
				if exc != nil {
					line(where + " EXCEPTION " + javaClassName(exc) + ": " + exc.GetMessage())
					continue
				}
				line(where + " | " + t.String() +
					" | ann=" + t.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) +
					" | unk=" + b(t.IsUnknownType()) + " null=" + b(t.IsNullable()) + " void=" + b(t.IsVoidable()) +
					" union=" + b(t.IsUnionType()) + " fn=" + b(t.IsFunctionType()) + " rec=" + b(t.IsRecordType()) +
					" tmpl=" + b(t.IsTemplatizedType()) + " enumel=" + b(t.IsEnumElementType()) +
					" noresolved=" + b(t.IsNoResolvedType()) + " obj=" + b(t.IsObject()) +
					" | <:num=" + b(t.IsSubtypeOf(number)) + " <:str=" + b(t.IsSubtypeOf(str)) +
					" <:obj=" + b(t.IsSubtypeOf(object)) +
					" | nonnull=" + t.RestrictByNotNullOrUndefined().String() +
					" | autobox=" + t.Autobox().String() +
					" | lub(num)=" + t.GetLeastSupertype(number).String() +
					" | glb(obj)=" + t.GetGreatestSubtype(object).String())
			}
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
	}
	walk(result.Ast)
	return sb.String()
}

// npeMessage strips NullPointerException messages: the JVM's "helpful" NPE
// text describes bytecode and has no Go equivalent; the exception class is
// what Java code can observe (jutil maps a Go nil dereference to it).
var npeMessage = regexp.MustCompile(`(?m)(EXCEPTION java\.lang\.NullPointerException): .*$`)

func TestCorpusTypesOracle(t *testing.T) {
	root, oracle, list := os.Getenv("TYPES_ROOT"), os.Getenv("TYPES_ORACLE_DIR"), os.Getenv("TYPES_FILES")
	if root == "" || oracle == "" || list == "" {
		t.Skip("set TYPES_ROOT, TYPES_ORACLE_DIR and TYPES_FILES to run the corpus type oracle")
	}
	rels, err := corpus.ReadList(list)
	if err != nil {
		t.Fatal(err)
	}
	diffDir := os.Getenv("TYPES_DIFF_DIR")
	files, same, exprs, differ := 0, 0, 0, 0
	for _, rel := range rels {
		want, err := readDump(filepath.Join(oracle, rel+".types"))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		want = npeMessage.ReplaceAllString(want, "$1")
		got := npeMessage.ReplaceAllString(dumpCorpusTypes(rel, string(src)), "$1")
		files++
		exprs += strings.Count(want, " | ann=")
		if got == want {
			same++
			continue
		}
		differ++
		if differ <= 10 {
			t.Errorf("%s: %s", rel, firstDiff(want, got))
		}
		if diffDir != "" {
			p := filepath.Join(diffDir, rel+".types")
			_ = os.MkdirAll(filepath.Dir(p), 0o755)
			_ = os.WriteFile(p, []byte(got), 0o644)
		}
	}
	t.Logf("compared %d files (%d typed expressions): %d identical, %d differ", files, exprs, same, differ)
}
