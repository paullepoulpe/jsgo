package diagnostic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// testdata/log_file.txt is tools/oracle/java/LogFileVectors.java run in a
// temporary directory (shown as <tmp>): the result of every LogFile call
// and the bytes of each log file. Closure has no LogFile tests.
//
// Go-only: language-adapter — the java.nio.file exceptions behind a failed
// createOrReopen (FileAlreadyExistsException, FileSystemException) are
// IOExceptions with the Go error's message in jutil.Files, so for those two
// steps only the RuntimeException wrapping is compared.

type logFileScenario struct {
	tmp string
	sb  strings.Builder
}

func (s *logFileScenario) step(name string, f func() any) {
	var result string
	jutil.Try(func() {
		r := f()
		if _, ok := r.(LogFile); ok {
			result = "ok"
		} else {
			result = jutil.ToJavaString(r)
		}
	}, func(e jutil.Throwable) {
		result = exceptionString(e)
		if c := e.GetCause(); c != nil {
			result += " / cause " + exceptionString(c)
		}
	})
	s.sb.WriteString(escape(strings.ReplaceAll(name+" -> "+result, s.tmp, "<tmp>")))
	s.sb.WriteByte('\n')
}

func (s *logFileScenario) bytes(file string) {
	b, err := os.ReadFile(file)
	if err != nil {
		s.sb.WriteString("bytes <missing>\n")
		return
	}
	fmt.Fprintf(&s.sb, "bytes %x\n", b)
}

type producerFunc func(w *jutil.JsonWriter)

func (f producerFunc) WriteJson(w *jutil.JsonWriter) { f(w) }

func ptr(s string) *string { return &s }

func TestLogFileMatchesJava(t *testing.T) {
	want, err := os.ReadFile("testdata/log_file.txt")
	if err != nil {
		t.Fatal(err)
	}
	s := &logFileScenario{tmp: t.TempDir()}
	var log LogFile

	s.sb.WriteString("=== writing\n")
	f1 := filepath.Join(s.tmp, "a/b/one.log")
	s.step("createOrReopen", func() any { log = LogFile_CreateOrReopen(jutil.Path_Of(f1)); return log })
	s.step("isLogging", func() any { return log.IsLogging() })
	s.step("log(String)", func() any { return log.LogString("plain <&> \u00e9 \U0001F600") })
	s.step("log(Object)", func() any { return log.Log(jutil.NewArrayListFrom([]any{int32(1), "x"})) })
	s.step("log(Object null)", func() any { return log.Log(nil) })
	s.step("log(Supplier)", func() any { return log.LogSupplier(func() string { return "supplied" }) })
	s.step("log(template)", func() any { return log.LogFormat("%s-%d-%s", "x", int32(3), nil) })
	s.step("logJson(Object)", func() any {
		m := jutil.NewLinkedHashMap[string, any]()
		m.Put("k", []string{"<v>"})
		return log.LogJson(m)
	})
	s.step("logJson(Supplier)", func() any { return log.LogJsonSupplier(func() any { return "s\u2028" }) })
	s.step("log high", func() any { return log.LogString(jutil.StringFromUTF16([]uint16{'h', 'i', 0xd83d})) })
	s.step("log low", func() any { return log.LogString(jutil.StringFromUTF16([]uint16{0xde00, 'l', 'o'})) })
	s.step("log lone", func() any { return log.LogString(jutil.StringFromUTF16([]uint16{'l', 'o', 'n', 'e', 0xd800, '!'})) })
	s.step("log lone low", func() any { return log.LogString(jutil.StringFromUTF16([]uint16{0xdc00})) })
	s.step("close", func() any { log.Close(); return "void" })
	s.step("close again", func() any { log.Close(); return "void" })
	s.step("log after close", func() any { return log.LogString("x") })
	s.bytes(f1)

	s.sb.WriteString("=== reopen appends\n")
	s.step("createOrReopen", func() any { log = LogFile_CreateOrReopen(jutil.Path_Of(f1)); return log })
	s.step("log(String)", func() any { return log.LogString("second") })
	s.step("close", func() any { log.Close(); return "void" })
	s.bytes(f1)

	s.sb.WriteString("=== streamed json\n")
	f2 := filepath.Join(s.tmp, "streamed.log")
	s.step("createOrReopen", func() any { log = LogFile_CreateOrReopen(jutil.Path_Of(f2)); return log })
	s.step("log(String)", func() any { return log.LogString("before") })
	s.step("logJson(producer)", func() any {
		return log.LogJsonProducer(producerFunc(func(w *jutil.JsonWriter) {
			w.BeginObject()
			w.Name("a<").BeginArray().ValueNumber(jutil.NewLazilyParsedNumber("1")).ValueString(ptr("<x>")).NullValue().EndArray()
			w.Name("n").NullValue()
			w.EndObject()
		}))
	})
	s.step("log after producer", func() any { return log.LogString("after") })
	s.step("close", func() any { log.Close(); return "void" })
	s.bytes(f2)

	s.sb.WriteString("=== incomplete streamed json\n")
	f3 := filepath.Join(s.tmp, "incomplete.log")
	s.step("createOrReopen", func() any { log = LogFile_CreateOrReopen(jutil.Path_Of(f3)); return log })
	s.step("logJson(producer)", func() any {
		return log.LogJsonProducer(producerFunc(func(w *jutil.JsonWriter) { w.BeginArray().ValueBoolean(true) }))
	})
	s.bytes(f3)

	s.sb.WriteString("=== throwing producer\n")
	f4 := filepath.Join(s.tmp, "throwing.log")
	s.step("createOrReopen", func() any { log = LogFile_CreateOrReopen(jutil.Path_Of(f4)); return log })
	s.step("logJson(producer)", func() any {
		return log.LogJsonProducer(producerFunc(func(w *jutil.JsonWriter) {
			w.BeginArray()
			panic(jutil.NewIllegalStateExceptionMsg("boom"))
		}))
	})
	s.bytes(f4)

	s.sb.WriteString("=== bad paths\n")
	s.step("relative file without parent", func() any { return LogFile_CreateOrReopen(jutil.Path_Of("rel-no-parent.log")) })
	if err := os.WriteFile(filepath.Join(s.tmp, "afile"), []byte("x"), 0o666); err != nil {
		t.Fatal(err)
	}
	s.step("parent is a file", func() any { return LogFile_CreateOrReopen(jutil.Path_Of(filepath.Join(s.tmp, "afile/x.log"))) })
	if err := os.MkdirAll(filepath.Join(s.tmp, "adir"), 0o777); err != nil {
		t.Fatal(err)
	}
	s.step("file is a directory", func() any { return LogFile_CreateOrReopen(jutil.Path_Of(filepath.Join(s.tmp, "adir"))) })

	s.sb.WriteString("=== no-op\n")
	s.step("createNoOp", func() any { log = LogFile_CreateNoOp(); return log })
	s.step("isLogging", func() any { return log.IsLogging() })
	s.step("log(Object null)", func() any { return log.Log(nil) })
	s.step("logJson(producer)", func() any {
		return log.LogJsonProducer(producerFunc(func(w *jutil.JsonWriter) { panic(jutil.NewAssertionError("not called")) }))
	})
	s.step("close", func() any { log.Close(); return "void" })

	gotLines := strings.Split(s.sb.String(), "\n")
	wantLines := strings.Split(string(want), "\n")
	if len(gotLines) != len(wantLines) {
		t.Fatalf("%d lines, Java %d:\n%s", len(gotLines), len(wantLines), s.sb.String())
	}
	for i := range wantLines {
		g, w := gotLines[i], wantLines[i]
		if strings.HasPrefix(w, "parent is a file -> ") || strings.HasPrefix(w, "file is a directory -> ") {
			prefix := w[:strings.Index(w, " -> ")] + " -> java.lang.RuntimeException: "
			if !strings.HasPrefix(g, prefix) || !strings.Contains(g, " / cause java.io.IOException: ") {
				t.Errorf("line %d: got %s, Java %s", i+1, g, w)
			}
			continue
		}
		if g != w {
			t.Errorf("line %d: got %s, Java %s", i+1, g, w)
		}
	}
}
