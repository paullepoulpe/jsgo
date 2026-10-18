package jscomp

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func zipCacheWriteArchive(t *testing.T, path, prefix string) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	for i := 0; i < 8; i++ {
		entry, err := writer.Create(fmt.Sprintf("part%d.js", i))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprintf(entry, "%s: entry %d é😀\n", prefix, i); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func zipCacheConcurrentReads(t *testing.T, readers int, open func(worker, iteration int) (io.Reader, string)) {
	t.Helper()
	start := make(chan struct{})
	failures := make(chan error, readers)
	var wg sync.WaitGroup
	for worker := 0; worker < readers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if failure := recover(); failure != nil {
					failures <- fmt.Errorf("worker %d panicked: %v", worker, failure)
				}
			}()
			<-start
			for iteration := 0; iteration < 24; iteration++ {
				reader, want := open(worker, iteration)
				got, err := io.ReadAll(reader)
				if closer, ok := reader.(io.Closer); ok {
					if closeErr := closer.Close(); err == nil {
						err = closeErr
					}
				}
				if err != nil {
					failures <- fmt.Errorf("worker %d read failed: %w", worker, err)
					return
				}
				if string(got) != want {
					failures <- fmt.Errorf("worker %d got %q, want %q", worker, got, want)
					return
				}
			}
		}()
	}
	close(start)
	wg.Wait()
	close(failures)
	for failure := range failures {
		t.Error(failure)
	}
}

func TestZipFileCacheConcurrentSharedCacheReads(t *testing.T) {
	paths := make([]string, 4)
	for i := range paths {
		paths[i] = filepath.Join(t.TempDir(), fmt.Sprintf("input%d.zip", i))
		zipCacheWriteArchive(t, paths[i], fmt.Sprintf("archive %d", i))
	}
	zipCacheConcurrentReads(t, 16, func(worker, iteration int) (io.Reader, string) {
		archive, entry := (worker+iteration)%len(paths), iteration%8
		return JSCompZipFileCache_GetEntryStream(paths[archive], fmt.Sprintf("part%d.js", entry)),
			fmt.Sprintf("archive %d: entry %d é😀\n", archive, entry)
	})
}

func TestZipFileCacheConcurrentRefreshAndEntryOpening(t *testing.T) {
	path := filepath.Join(t.TempDir(), "refresh.zip")
	zipCacheWriteArchive(t, path, "first")
	cached := newJSCompZipFileCache_CachedZipFile(path)
	defer cached.MaybeClose()
	for _, prefix := range []string{"first", "changed"} {
		if prefix == "changed" {
			previous, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			zipCacheWriteArchive(t, path, prefix)
			// A deterministic newer stamp avoids a sleep or filesystem clock
			// granularity dependency. All prior streams have finished reading.
			stamp := previous.ModTime().Add(2 * time.Second)
			if err := os.Chtimes(path, stamp, stamp); err != nil {
				t.Fatal(err)
			}
		}
		zipCacheConcurrentReads(t, 16, func(worker, iteration int) (io.Reader, string) {
			entry := (worker + iteration) % 8
			return cached.GetEntryStream(fmt.Sprintf("part%d.js", entry)),
				fmt.Sprintf("%s: entry %d é😀\n", prefix, entry)
		})
	}
}
