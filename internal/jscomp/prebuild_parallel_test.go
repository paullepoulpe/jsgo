package jscomp

import (
	"fmt"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type prebuildBlockingLoader struct {
	*SourceFile_CodeLoader_Preloaded
	entered chan<- struct{}
	release <-chan struct{}
}

func (l *prebuildBlockingLoader) LoadUncachedCode() string {
	l.entered <- struct{}{}
	<-l.release
	return l.SourceFile_CodeLoader_Preloaded.LoadUncachedCode()
}

// Blocking two independent source reads proves that real parsing/scanning
// work overlaps; the third cannot start until one of the bounded slots opens.
func TestPrebuildWorkersOverlapAndRespectLimit(t *testing.T) {
	for _, dependency := range []bool{false, true} {
		t.Run(fmt.Sprintf("dependency=%v", dependency), func(t *testing.T) {
			compiler := NewCompiler()
			compiler.InitOptions(NewCompilerOptions())
			compiler.SetErrorManager(NewSortingErrorManager(nil))
			compiler.SetPreferRegexParser(dependency)
			compiler.moduleLoader = ModuleLoader_EMPTY
			entered := make(chan struct{}, 4)
			release := make(chan struct{})
			released := false
			defer func() {
				if !released {
					close(release)
				}
			}()
			var inputs []*CompilerInput
			for i := 0; i < 4; i++ {
				loader := &prebuildBlockingLoader{
					SourceFile_CodeLoader_Preloaded: NewSourceFile_CodeLoader_Preloaded(fmt.Sprintf("goog.provide('test.file%d'); var value%d = %d;", i, i, i)),
					entered:                         entered, release: release,
				}
				source := newSourceFile(loader, fmt.Sprintf("file%d.js", i), rhino.StaticSourceFile_SourceKind_STRONG)
				input := NewCompilerInput(source)
				input.SetCompiler(compiler)
				inputs = append(inputs, input)
			}
			done := make(chan any, 1)
			go func() {
				defer func() { done <- recover() }()
				if dependency {
					NewPrebuildDependencyInfo(2).Prebuild(slices.Values(inputs))
				} else {
					NewPrebuildAst(compiler, 2).Prebuild(slices.Values(inputs))
				}
			}()
			for range 2 {
				select {
				case <-entered:
				case failure := <-done:
					t.Fatalf("prebuild ended before overlapping reads: %v", failure)
				case <-time.After(5 * time.Second):
					t.Fatal("source reads did not overlap")
				}
			}
			select {
			case <-entered:
				t.Fatal("more source reads than the requested worker limit")
			default:
			}
			close(release)
			released = true
			select {
			case failure := <-done:
				if failure != nil {
					t.Fatalf("prebuild failed: %v", failure)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("prebuild did not finish")
			}
			for i, input := range inputs {
				if dependency {
					if !slices.Contains(input.GetProvides(), fmt.Sprintf("test.file%d", i)) {
						t.Fatalf("missing dependency info for input %d", i)
					}
				} else if input.ast.root.IsNil() || input.ast.root.GetStaticSourceFile() != rhino.StaticSourceFile(input.sourceFile) {
					t.Fatalf("missing/misbound AST for input %d", i)
				}
			}
		})
	}
}

func TestParallelPrebuildRunsAllTasksAfterFailure(t *testing.T) {
	var completed atomic.Int32
	failures := parallelPrebuild(8, 3, func(i int) {
		defer completed.Add(1)
		if i == 1 || i == 5 {
			panic(jutil.NewIllegalArgumentExceptionMsg(fmt.Sprintf("failure-%d", i)))
		}
	})
	if completed.Load() != 8 {
		t.Fatalf("completed %d of 8 tasks", completed.Load())
	}
	for i, failure := range failures {
		if i == 1 || i == 5 {
			if failure == nil || failure.GetMessage() != fmt.Sprintf("failure-%d", i) {
				t.Fatalf("wrong failure in slot %d: %v", i, failure)
			}
		} else if failure != nil {
			t.Fatalf("unexpected failure in slot %d: %v", i, failure)
		}
	}
}
