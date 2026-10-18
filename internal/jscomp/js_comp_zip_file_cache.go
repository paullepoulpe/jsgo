// Port of com/google/javascript/jscomp/JSCompZipFileCache.java @ closure-compiler v20261015

/*
 * Copyright 2018 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSCompZipFileCache is a class that abstract entries from zip files via
// managed caching.
//
// Normally zip files created via URL connection cache indefinite by JDK.
// That's not good since a zip file contents might change over time (e.g.
// compiler running as a worker). This class provides a timestamp controlled
// caching which ensures we always read up-to-date zip while avoiding wasting
// time by re-reading the zip for each entry.
type JSCompZipFileCache struct{}

// JSCompZipFileCache_GetEntryStream returns a stream for an entry.
//
// We expose a stream rather than something like a ZipEntry to avoid
// possible cache-freshness issues. Callers must read the stream "quickly"
// and may only do so once. It panics with an IOException.
//
// Java: JSCompZipFileCache.getEntryStream(String, String)
func JSCompZipFileCache_GetEntryStream(zipName string, entryName string) io.Reader {
	// Go-only: language-adapter — LoadingCache's Go representation mutates
	// its lookup table and LRU links on hits as well as misses. Synchronize
	// cache access and entry opening; reading the returned stream and parsing
	// its source remain independent, parallel work.
	jSCompZipFileCache_monitor.Lock()
	defer jSCompZipFileCache_monitor.Unlock()
	return jSCompZipFileCache_zipFileCache.GetUnchecked(zipName).GetEntryStream(entryName)
}

// Go-only: language-adapter — protects the shared emulated cache, including
// eviction callbacks that close a cached archive.
var jSCompZipFileCache_monitor jutil.Monitor

// Go-only: language-adapter — there are no Java system properties; the
// "jscomp.zipfile.cachesize" property takes its default.
var jSCompZipFileCache_ZIP_CACHE_SIZE = int64(jutil.IntegerParseInt("1000", 10))

var jSCompZipFileCache_zipFileCache = jutil.CacheBuilder_Build(
	jutil.CacheBuilder_NewBuilder().
		MaximumSize(jSCompZipFileCache_ZIP_CACHE_SIZE).
		RemovalListener(func(key, value any) {
			jutil.Try(func() {
				value.(*JSCompZipFileCache_CachedZipFile).MaybeClose()
			}, func(e jutil.IIOException) {
				panic(jutil.NewRuntimeExceptionWithCause(e))
			})
		}),
	func(key string) *JSCompZipFileCache_CachedZipFile {
		return newJSCompZipFileCache_CachedZipFile(key)
	})

// JSCompZipFileCache_CachedZipFile is JSCompZipFileCache.CachedZipFile.
type JSCompZipFileCache_CachedZipFile struct {
	path         *jutil.Path
	zipFile      *zip.ReadCloser
	lastModified *time.Time

	// Go-only: language-adapter — the object's monitor, which the
	// synchronized block of refreshIfNeeded holds.
	monitor jutil.Monitor
}

func newJSCompZipFileCache_CachedZipFile(zipName string) *JSCompZipFileCache_CachedZipFile {
	return &JSCompZipFileCache_CachedZipFile{path: jutil.Path_Of(zipName)}
}

// Java: JSCompZipFileCache.CachedZipFile.getEntryStream(String)
func (z *JSCompZipFileCache_CachedZipFile) GetEntryStream(entryName string) io.Reader {
	// Go-only: language-adapter — acquire metadata and the entry from the
	// same archive generation. The returned reader owns its current stream;
	// no monitor is held while its caller reads or parses it.
	z.monitor.Lock()
	defer z.monitor.Unlock()
	z.refreshIfNeededLocked()
	entry := jSCompZipFileCache_getEntry(z.zipFile, entryName)
	if entry == nil {
		panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat("%s!/%s", z.path, entryName)))
	}
	rc, err := entry.Open()
	if err != nil {
		panic(jSCompZipFileCache_ioException(err, z.path.String()))
	}
	return rc
}

// jSCompZipFileCache_getEntry is ZipFile.getEntry(String): the entry with
// the given name, or else the directory entry name + "/".
//
// Go-only: language-adapter — archive/zip has no lookup by name.
func jSCompZipFileCache_getEntry(z *zip.ReadCloser, name string) *zip.File {
	for _, f := range z.File {
		if f.Name == name {
			return f
		}
	}
	for _, f := range z.File {
		if f.Name == name+"/" {
			return f
		}
	}
	return nil
}

// jSCompZipFileCache_ioException maps a Go I/O error to the IOException the
// JDK would throw.
//
// Go-only: language-adapter — Go reports I/O failures as errors.
func jSCompZipFileCache_ioException(err error, path string) jutil.Throwable {
	e := &jutil.IOException{}
	if errors.Is(err, fs.ErrNotExist) {
		jutil.InitThrowable(e, "java.nio.file.NoSuchFileException", &path, nil)
	} else {
		msg := err.Error()
		jutil.InitThrowable(e, "java.io.IOException", &msg, nil)
	}
	return e
}

// Java: JSCompZipFileCache.CachedZipFile.refreshIfNeeded()
func (z *JSCompZipFileCache_CachedZipFile) RefreshIfNeeded() {
	// Go-only: language-adapter — Java's unsynchronized fast-path reads do
	// not establish Go's required happens-before relation for lastModified
	// and zipFile. Protect the entire check and metadata publication.
	z.monitor.Lock()
	defer z.monitor.Unlock()
	z.refreshIfNeededLocked()
}

// Go-only: language-adapter — caller holds z.monitor, avoiding a second
// re-entrant monitor acquisition in GetEntryStream.
func (z *JSCompZipFileCache_CachedZipFile) refreshIfNeededLocked() {
	info, err := os.Stat(z.path.String())
	if err != nil {
		panic(jSCompZipFileCache_ioException(err, z.path.String()))
	}
	newLastModified := info.ModTime()
	if z.lastModified != nil && newLastModified.Equal(*z.lastModified) {
		return
	}

	z.maybeCloseLocked()
	zf, err := zip.OpenReader(z.path.String())
	if err != nil {
		panic(jSCompZipFileCache_ioException(err, z.path.String()))
	}
	z.zipFile = zf
	z.lastModified = &newLastModified
}

// Java: JSCompZipFileCache.CachedZipFile.maybeClose()
func (z *JSCompZipFileCache_CachedZipFile) MaybeClose() {
	z.monitor.Lock()
	defer z.monitor.Unlock()
	z.maybeCloseLocked()
}

// Go-only: language-adapter — caller holds z.monitor.
func (z *JSCompZipFileCache_CachedZipFile) maybeCloseLocked() {
	if z.zipFile != nil {
		if err := z.zipFile.Close(); err != nil {
			panic(jSCompZipFileCache_ioException(err, z.path.String()))
		}
	}
}

// Java: JSCompZipFileCache()
func newJSCompZipFileCache() *JSCompZipFileCache {
	panic(jutil.NewAssertionErrorNoMessage())
}
