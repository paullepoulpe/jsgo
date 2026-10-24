package jscomp

// Go mirror of tools/oracle/java/HashOrderDump.java: hand-written cases whose
// diagnostics come out in Java's hash-iteration order (keys sharing a bucket,
// a table that grows, a table cleared and reused), compiled and checked the
// way the Java tool does, and compared with its output in report order
// (testdata/hash_order/java.txt, regenerated with:
//
//	java -cp .cache/$CLOSURE_JAR:.cache/oracle-build \
//	  com.google.javascript.jscomp.HashOrderDump \
//	  internal/jscomp/testdata/hash_order/cases.txt \
//	  internal/jscomp/testdata/hash_order/java.txt
//
// after tools/oracle/build.sh).

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type hashOrderCase struct {
	name, pass string
	files      []*SourceFile
}

// hashOrderReadCases is HashOrderDump.readCases.
func hashOrderReadCases(lines []string) []hashOrderCase {
	var cases []hashOrderCase
	var cur *hashOrderCase
	fileName := ""
	inFile := false
	var src strings.Builder
	flushFile := func() {
		if inFile {
			cur.files = append(cur.files, SourceFile_FromCode(fileName, src.String()))
			inFile = false
		}
	}
	for _, l := range lines {
		switch {
		case strings.HasPrefix(l, "# case "):
			flushFile()
			if cur != nil {
				cases = append(cases, *cur)
			}
			parts := strings.Split(strings.TrimPrefix(l, "# case "), " ")
			cur = &hashOrderCase{name: parts[0], pass: parts[1]}
		case strings.HasPrefix(l, "## file "):
			flushFile()
			fileName = strings.TrimPrefix(l, "## file ")
			inFile = true
			src.Reset()
		case inFile:
			src.WriteString(l)
			src.WriteByte('\n')
		}
	}
	flushFile()
	if cur != nil {
		cases = append(cases, *cur)
	}
	return cases
}

func hashOrderCreatePass(t testing.TB, pass string, compiler AbstractCompiler) CompilerPass {
	switch pass {
	case "CheckNoMutatedEs6Exports":
		return NewCheckNoMutatedEs6Exports(compiler)
	}
	t.Fatalf("unknown pass %q", pass)
	return nil
}

func hashOrderFormat(e *JSError) string {
	sourceName := "null"
	if s := e.GetSourceName(); s != nil {
		sourceName = *s
	}
	return e.GetType().key +
		" " + sourceName +
		":" + strconv.Itoa(e.GetLineno()) +
		":" + strconv.Itoa(e.GetCharno()) +
		" \"" + e.GetDescription() + "\""
}

// hashOrderCompiler records every JSError reported while recording is on,
// in report order, as HashOrderDump's Compiler.report override does.
type hashOrderCompiler struct {
	*Compiler
	recording bool
	diags     []string
}

func (c *hashOrderCompiler) Report(error *JSError) {
	if c.recording {
		c.diags = append(c.diags, hashOrderFormat(error))
	}
	c.Compiler.Report(error)
}

// hashOrderDumpOne is HashOrderDump.dumpOne.
func hashOrderDumpOne(t testing.TB, sb *strings.Builder, c hashOrderCase) {
	lintChecksDumpLine(sb, "# case "+c.name+" "+c.pass)
	compiler := &hashOrderCompiler{Compiler: NewCompiler()}
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", "")},
		c.files,
		lintChecksDumpOptions(t, "google"))
	var root rhino.Node
	jutil.Try(func() {
		root = compiler.ParseInputs()
	}, func(e jutil.IRuntimeException) {
		root = rhino.Node{}
	})
	if root.IsNil() || compiler.HasErrors() {
		failed := "# parse-failed"
		for _, e := range compiler.GetErrors() {
			failed += " " + e.GetType().key
		}
		lintChecksDumpLine(sb, failed)
		return
	}
	compiler.recording = true
	jutil.Try(func() {
		hashOrderCreatePass(t, c.pass, compiler).Process(root.GetFirstChild(), root.GetSecondChild())
	}, func(e jutil.IRuntimeException) {
		compiler.diags = append(compiler.diags, "# exception "+lintChecksDumpClassName(e))
	})
	compiler.recording = false
	for _, d := range compiler.diags {
		lintChecksDumpLine(sb, d)
	}
}

func TestHashOrderOracle(t *testing.T) {
	casesText, err := os.ReadFile("testdata/hash_order/cases.txt")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/hash_order/java.txt")
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	for _, c := range hashOrderReadCases(strings.Split(string(casesText), "\n")) {
		hashOrderDumpOne(t, &sb, c)
	}
	if got := sb.String(); got != string(want) {
		gotLines, wantLines := strings.Split(got, "\n"), strings.Split(string(want), "\n")
		for i := 0; i < len(gotLines) || i < len(wantLines); i++ {
			var g, w string
			if i < len(gotLines) {
				g = gotLines[i]
			}
			if i < len(wantLines) {
				w = wantLines[i]
			}
			if g != w {
				t.Fatalf("hash order dump differs from Java at line %d:\n got: %s\nwant: %s", i+1, g, w)
			}
		}
	}
}
