package linter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestLinterConcurrentAndNodeStoreIsolation(t *testing.T) {
	rhino.ResetNodeStore()
	defer rhino.ResetNodeStore()

	savedOut := jutil.System_Out()
	jutil.System_SetOut(jutil.NewPrintStream(io.Discard))
	defer jutil.System_SetOut(savedOut)

	dir := t.TempDir()
	const workers = 6
	lintFiles := make([]string, workers)
	fixFiles := make([]string, workers)
	for i := 0; i < workers; i++ {
		lintFiles[i] = filepath.Join(dir, fmt.Sprintf("lint_%d.js", i))
		if err := os.WriteFile(lintFiles[i], []byte(fmt.Sprintf("goog.module('test.lint%d');\nconst x = %d;\nexports = x;\n", i, i)), 0o644); err != nil {
			t.Fatal(err)
		}
		fixFiles[i] = filepath.Join(dir, fmt.Sprintf("fix_%d.js", i))
		if err := os.WriteFile(fixFiles[i], []byte("goog.module('test.fix');\ngoog.require('b.Mod');\ngoog.require('a.Mod');\nconst u = [a, b]\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	l := NewLinter_Builder().Build()
	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			l.Lint(lintFiles[id])
			if err := l.FixRepeatedly(fixFiles[id]); err != nil {
				errCh <- err
			}
		}(i)
	}
	close(start)
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	if !rhino.CurrentNodeStore().IsStatic() {
		t.Fatalf("expected rhino.CurrentNodeStore().IsStatic() == true after standalone Linter calls")
	}

	fixedBytes, err := os.ReadFile(fixFiles[0])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(fixedBytes), "goog.require('a.Mod')") {
		t.Fatalf("expected FixRepeatedly to remove extra goog.require, got:\n%s", string(fixedBytes))
	}

	// Verify an already-active outer NodeStore is preserved.
	outer := rhino.NewNodeStore()
	prev := rhino.SetCurrentNodeStore(outer)
	defer rhino.SetCurrentNodeStore(prev)
	l.Lint(lintFiles[0])
	if rhino.CurrentNodeStore() != outer {
		t.Fatalf("Linter.Lint clobbered active outer NodeStore")
	}
}
