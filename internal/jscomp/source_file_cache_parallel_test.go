package jscomp

import (
	"fmt"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestParallelSourceSharedNumberOrigin(t *testing.T) {
	// Moved/inlined nodes in different output scripts can share an original
	// source. Its line cache is independent of the parser's position scanner.
	origin := SourceFile_FromCode("original.js", "'😀';\r\n0x2a;")
	origin.GetCode()
	if origin.cached.Load().lineOffsets != nil {
		t.Fatal("test requires an uninitialized line-position cache")
	}
	options := NewCompilerOptions()
	options.SetPrettyPrint(true)
	options.SetNumParallelThreads(4)
	compiler := newTestCompiler(options)
	roots := make([]*rhino.Node, 8)
	for i := range roots {
		number := rhino.IR_Number(42)
		number.SetStaticSourceFile(origin).SetLinenoCharno(2, 0).SetLength(4)
		roots[i] = rhino.IR_Script()
		roots[i].AddChildToBack(rhino.IR_ExprResult(number))
	}
	// Use an independent source for the serial oracle, keeping the shared
	// parallel source's lazy offsets cold until the worker printers run.
	serialOrigin := SourceFile_FromCode("original.js", "'😀';\r\n0x2a;")
	serialOptions := NewCompilerOptions()
	serialOptions.SetPrettyPrint(true)
	serialCompiler := newTestCompiler(serialOptions)
	want := make([]string, len(roots))
	for i := range want {
		number := rhino.IR_Number(42)
		number.SetStaticSourceFile(serialOrigin).SetLinenoCharno(2, 0).SetLength(4)
		root := rhino.IR_Script()
		root.AddChildToBack(rhino.IR_ExprResult(number))
		want[i] = serialCompiler.toSourceAndMappings(root, i == 0,
			newCompiler_SingleBinaryLicenseTracker(serialCompiler)).source
	}
	if origin.cached.Load().lineOffsets != nil {
		t.Fatal("serial oracle warmed the shared parallel source")
	}
	count := 0
	compiler.toSourceRoots(roots, newCompiler_SingleBinaryLicenseTracker(compiler), func(i int, root *rhino.Node, result *parallelSource) {
		if result == nil {
			t.Fatal("test requires the parallel printer path")
		}
		if got := result.code.source; got != want[i] {
			t.Errorf("root %d: parallel pretty source %q, serial source %q", i, got, want[i])
		}
		count++
	})
	if count != len(roots) {
		t.Fatalf("printed %d roots, want %d", count, len(roots))
	}
	if origin.cached.Load().lineOffsets == nil {
		t.Fatal("parallel printers did not exercise the shared lazy line cache")
	}
	if got := origin.GetLineOffset(2); got != jutil.StringLength("'😀';\r\n") {
		t.Errorf("line offset %d does not use UTF-16 positions", got)
	}
}

func TestSourceFileConcurrentCacheInvalidation(t *testing.T) {
	const a = "😀\r\n0x2a;\u2028last"
	const b = "x\n0x2a;\nend"
	file := SourceFile_FromCode("changing.js", a)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				file.SetCodeDeprecated(b)
			} else {
				file.SetCodeDeprecated(a)
			}
			if i%5 == 0 {
				file.ClearCachedSource()
			}
		}
	}()
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for range 100 {
				if got := file.GetLine(2); got == nil || *got != "0x2a;" {
					t.Errorf("line/code cache generations were mixed: %v", got)
					return
				}
				if got := file.GetLines(2, 4); got == nil || got.GetSourceExcerpt() != "0x2a;" {
					t.Errorf("region/code cache generations were mixed: %v", got)
					return
				}
				if got := file.GetLineOffset(2); got != 2 && got != 4 {
					t.Errorf("invalid line offset %d", got)
					return
				}
				if file.GetLineOfOffset(0) != 1 || file.GetColumnOfOffset(0) != 0 || file.GetNumLines() != 3 {
					t.Error("invalid source positions/counts")
					return
				}
				if got := file.GetNumBytes(); got != jutil.StringLength(a) && got != jutil.StringLength(b) {
					t.Errorf("invalid source length %d", got)
					return
				}
				file.HasSourceInMemory()
				file.GetProto()
			}
		}()
	}
	close(start)
	wg.Wait()
}

type sourceFileCacheFailingLoader struct {
	SourceFile_CodeLoaderBase
	failed bool
	loads  int
}

func (l *sourceFileCacheFailingLoader) ToProtoLocationBuilder(name string) *SourceFileProto_Builder {
	return NewSourceFile_CodeLoader_Preloaded("").ToProtoLocationBuilder(name)
}

func (l *sourceFileCacheFailingLoader) LoadUncachedCode() string {
	l.loads++
	if l.failed {
		panic(jutil.NewIOException("not available"))
	}
	return "first\nsecond"
}

func TestSourceFilePositionCacheIOExceptionAndRetry(t *testing.T) {
	loader := &sourceFileCacheFailingLoader{failed: true}
	file := newSourceFile(loader, "unavailable.js", rhino.StaticSourceFile_SourceKind_STRONG)
	if got := file.GetLineOffset(1); got != 0 || loader.loads != 1 {
		t.Fatalf("offset IOException fallback: offset=%d loads=%d", got, loader.loads)
	}
	if got := file.GetLineOfOffset(10); got != 1 || loader.loads != 1 {
		t.Fatalf("fallback should remain cached: line=%d loads=%d", got, loader.loads)
	}
	if got := file.GetLine(1); got != nil || loader.loads != 2 {
		t.Fatalf("line API must retry code loading: line=%v loads=%d", got, loader.loads)
	}
	loader.failed = false
	if got := file.GetLine(2); got == nil || *got != "second" || loader.loads != 3 {
		t.Fatalf("successful retry must replace fallback: line=%v loads=%d", got, loader.loads)
	}
	old := file.getLineOffsets()
	file.SetCodeDeprecated("new\r\n😀")
	if got := file.GetLineOffset(2); got != 5 {
		t.Errorf("replacement did not invalidate positions: %d", got)
	}
	if got := fmt.Sprint(old.lineOffsets); got != "[0 6]" || *old.code != "first\nsecond" {
		t.Errorf("published snapshot was mutated: %s %q", got, *old.code)
	}
}
