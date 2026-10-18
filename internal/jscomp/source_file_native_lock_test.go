package jscomp

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Embedding a native loader must not accidentally select the native lock: its
// overrides can call back into the same SourceFile while loading or opening.
type sourceFileReentrantTestLoader struct {
	*SourceFile_CodeLoader_Preloaded
	file *SourceFile
}

func (l *sourceFileReentrantTestLoader) LoadUncachedCode() string {
	l.file.ClearCachedSource()
	l.file.SetCodeDeprecated("first\ncallback")
	line := l.file.GetLine(2)
	if line == nil || *line != "callback" {
		panic("reentrant cache and offset publication failed")
	}
	return "outer\nresult"
}

func (l *sourceFileReentrantTestLoader) OpenUncachedReader() jutil.JavaReader {
	l.file.ClearCachedSource()
	l.file.SetCodeDeprecated("reader\ncallback")
	return nil // GetCodeReader subsequently uses the reentrantly published code.
}

func TestNativeSourceLockKeepsCustomLoaderReentry(t *testing.T) {
	for _, reader := range []bool{false, true} {
		loader := &sourceFileReentrantTestLoader{SourceFile_CodeLoader_Preloaded: NewSourceFile_CodeLoader_Preloaded("unused")}
		file := newSourceFile(loader, "reentrant.js", rhino.StaticSourceFile_SourceKind_STRONG)
		loader.file = file
		if file.monitor.reentrant == nil {
			t.Fatal("embedding loader lost Java monitor reentrancy")
		}
		done := make(chan any, 1)
		go func() {
			defer func() { done <- recover() }()
			if reader {
				if got := jutil.CharStreams_ToString(file.GetCodeReader()); got != "reader\ncallback" {
					panic("reader callback changed source")
				}
			} else if got := file.GetCode(); got != "outer\nresult" {
				panic("load callback changed source")
			}
		}()
		select {
		case failure := <-done:
			if failure != nil {
				t.Fatalf("reader=%v callback failed: %v", reader, failure)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("reader=%v loader callback deadlocked", reader)
		}
	}
}

func TestNativeSourceLockSelectsOnlyExactNativeLoaders(t *testing.T) {
	loaders := []SourceFile_CodeLoader{
		NewSourceFile_CodeLoader_Preloaded("code"),
		NewSourceFile_CodeLoader_StubSourceFileCodeLoader(),
		NewSourceFile_CodeLoader_OnDisk(jutil.Path_Of("not-opened.js"), jutil.StandardCharsets_UTF_8),
		NewSourceFile_CodeLoader_AtZip("not-opened.zip", "input.js", jutil.StandardCharsets_UTF_8),
	}
	for _, loader := range loaders {
		file := newSourceFile(loader, "native.js", rhino.StaticSourceFile_SourceKind_STRONG)
		if file.monitor.reentrant != nil {
			t.Fatalf("native loader %T still uses a goroutine-owned monitor", loader)
		}
	}
}

func TestNativeSourceLockLoadFailureAndInvalidation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "changes.js")
	file := SourceFile_FromFile(path)
	failed := false
	jutil.Try(func() { file.GetCode() }, func(jutil.IIOException) { failed = true })
	if !failed {
		t.Fatal("missing file must preserve its IOException")
	}
	if err := os.WriteFile(path, []byte("old\nline"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := file.GetLine(2); got == nil || *got != "line" {
		t.Fatalf("load failure left lock/cache unusable: %v", got)
	}
	old := file.getCodeAndLineOffsets()
	if err := os.WriteFile(path, []byte("new\r\n😀"), 0600); err != nil {
		t.Fatal(err)
	}
	file.ClearCachedSource()
	if got := file.GetLine(2); got == nil || *got != "😀" || file.GetLineOffset(2) != 5 {
		t.Fatalf("native reload mixed cache generations: %v", got)
	}
	if *old.code != "old\nline" || old.lineOffsets[1] != 4 {
		t.Fatal("native reload mutated a previously published snapshot")
	}
}
