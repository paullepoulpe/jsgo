package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type sourceFileReentryThenIOExceptionLoader struct {
	SourceFile_CodeLoaderBase
	file  *SourceFile
	loads int
}

func (l *sourceFileReentryThenIOExceptionLoader) LoadUncachedCode() string {
	l.loads++
	l.file.SetCodeDeprecated("callback\nsource")
	// Reenter position initialization too, so the outer catch must replace
	// both offset arrays without discarding the successfully published text.
	if got := l.file.GetLineOffset(2); got != 9 {
		panic(jutil.NewIllegalStateExceptionMsg("callback line offset changed"))
	}
	panic(jutil.NewIOException("loader failed after callback"))
}

func TestSourceFileReentrantIOExceptionPreservesPublishedCode(t *testing.T) {
	loader := &sourceFileReentryThenIOExceptionLoader{}
	file := newSourceFile(loader, "callback-error.js", rhino.StaticSourceFile_SourceKind_STRONG)
	loader.file = file
	// Offset lookup catches the loader's IOException and caches its original
	// one-line fallback. The callback's source and metrics must remain cached.
	if got := file.GetLineOffset(1); got != 0 {
		t.Fatalf("IOException offset fallback = %d", got)
	}
	if !file.HasSourceInMemory() || file.GetCode() != "callback\nsource" || loader.loads != 1 {
		t.Fatalf("callback source discarded or reloaded: inMemory=%v loads=%d", file.HasSourceInMemory(), loader.loads)
	}
	if file.GetNumLines() != 2 || file.GetNumBytes() != 15 {
		t.Fatal("callback metrics discarded by IOException fallback")
	}
	state := file.getLineOffsets()
	if len(state.lineOffsets) != 1 || len(state.lineByteOffsets) != 1 || state.lineOffsets[0] != 0 || state.lineByteOffsets[0] != 0 {
		t.Fatal("IOException did not replace both position arrays together")
	}
}
