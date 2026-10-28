// Port of com/google/javascript/jscomp/JSCompZipFileCache.java @ closure-compiler v20261021

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
	"sync"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSCompZipFileCache abstracts entries from zip files via timestamp-checked
// bounded LRU caching so updated zip archives are refreshed automatically
// without re-opening the central directory on every entry read.
type JSCompZipFileCache struct{}

const jSCompZipFileCacheMaxSize = 1000

type zipFileCacheNode struct {
	key        string
	val        *JSCompZipFileCache_CachedZipFile
	prev, next *zipFileCacheNode
}

type zipFileLRUCache struct {
	mu      sync.Mutex
	entries map[string]*zipFileCacheNode
	head    *zipFileCacheNode
	tail    *zipFileCacheNode
}

var jSCompZipFileCache_zipFileCache = zipFileLRUCache{
	entries: make(map[string]*zipFileCacheNode),
}

func (c *zipFileLRUCache) get(zipName string) *JSCompZipFileCache_CachedZipFile {
	c.mu.Lock()
	if node, ok := c.entries[zipName]; ok {
		c.moveToFrontLocked(node)
		val := node.val
		c.mu.Unlock()
		return val
	}
	var evicted *JSCompZipFileCache_CachedZipFile
	if len(c.entries) >= jSCompZipFileCacheMaxSize && c.tail != nil {
		old := c.tail
		c.unlinkLocked(old)
		delete(c.entries, old.key)
		evicted = old.val
	}
	val := newJSCompZipFileCache_CachedZipFile(zipName)
	node := &zipFileCacheNode{key: zipName, val: val}
	c.pushFrontLocked(node)
	c.entries[zipName] = node
	c.mu.Unlock()

	if evicted != nil {
		if err := evicted.maybeCloseErr(); err != nil {
			panic(jutil.NewRuntimeExceptionWithCause(err))
		}
	}
	return val
}

func (c *zipFileLRUCache) pushFrontLocked(n *zipFileCacheNode) {
	n.prev = nil
	n.next = c.head
	if c.head != nil {
		c.head.prev = n
	} else {
		c.tail = n
	}
	c.head = n
}

func (c *zipFileLRUCache) unlinkLocked(n *zipFileCacheNode) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
	n.prev = nil
	n.next = nil
}

func (c *zipFileLRUCache) moveToFrontLocked(n *zipFileCacheNode) {
	if c.head == n {
		return
	}
	c.unlinkLocked(n)
	c.pushFrontLocked(n)
}

// JSCompZipFileCache_GetEntryStream returns a stream for an entry.
//
// Java: JSCompZipFileCache.getEntryStream(String, String)
func JSCompZipFileCache_GetEntryStream(zipName string, entryName string) io.Reader {
	rc, err := jSCompZipFileCache_getEntryStreamErr(zipName, entryName)
	if err != nil {
		panic(err)
	}
	return rc
}

func jSCompZipFileCache_getEntryStreamErr(zipName string, entryName string) (io.ReadCloser, jutil.IIOException) {
	return jSCompZipFileCache_zipFileCache.get(zipName).getEntryStreamErr(entryName)
}

// JSCompZipFileCache_CachedZipFile is JSCompZipFileCache.CachedZipFile.
type JSCompZipFileCache_CachedZipFile struct {
	path         *jutil.Path
	zipFile      *zip.ReadCloser
	entries      map[string]*zip.File
	lastModified *time.Time
	mu           sync.Mutex
}

func newJSCompZipFileCache_CachedZipFile(zipName string) *JSCompZipFileCache_CachedZipFile {
	return &JSCompZipFileCache_CachedZipFile{path: jutil.Path_Of(zipName)}
}

// Java: JSCompZipFileCache.CachedZipFile.getEntryStream(String)
func (z *JSCompZipFileCache_CachedZipFile) GetEntryStream(entryName string) io.Reader {
	rc, err := z.getEntryStreamErr(entryName)
	if err != nil {
		panic(err)
	}
	return rc
}

func (z *JSCompZipFileCache_CachedZipFile) getEntryStreamErr(entryName string) (io.ReadCloser, jutil.IIOException) {
	z.mu.Lock()
	defer z.mu.Unlock()
	if err := z.refreshIfNeededLockedErr(); err != nil {
		return nil, err
	}
	entry := z.getEntryLocked(entryName)
	if entry == nil {
		panic(jutil.NewNullPointerExceptionMsg(jutil.LenientFormat("%s!/%s", z.path, entryName)))
	}
	rc, err := entry.Open()
	if err != nil {
		return nil, jSCompZipFileCache_ioException(err, z.path.String())
	}
	return rc, nil
}

func (z *JSCompZipFileCache_CachedZipFile) getEntryLocked(name string) *zip.File {
	if f, ok := z.entries[name]; ok {
		return f
	}
	if f, ok := z.entries[name+"/"]; ok {
		return f
	}
	return nil
}

func jSCompZipFileCache_ioException(err error, path string) jutil.IIOException {
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
	z.mu.Lock()
	defer z.mu.Unlock()
	if err := z.refreshIfNeededLockedErr(); err != nil {
		panic(err)
	}
}

func (z *JSCompZipFileCache_CachedZipFile) refreshIfNeededLockedErr() jutil.IIOException {
	info, err := os.Stat(z.path.String())
	if err != nil {
		return jSCompZipFileCache_ioException(err, z.path.String())
	}
	newLastModified := info.ModTime()
	if z.lastModified != nil && newLastModified.Equal(*z.lastModified) {
		return nil
	}

	if err := z.maybeCloseLockedErr(); err != nil {
		return err
	}
	zf, err := zip.OpenReader(z.path.String())
	if err != nil {
		return jSCompZipFileCache_ioException(err, z.path.String())
	}
	entries := make(map[string]*zip.File, len(zf.File))
	for _, f := range zf.File {
		if _, exists := entries[f.Name]; !exists {
			entries[f.Name] = f
		}
	}
	z.zipFile = zf
	z.entries = entries
	z.lastModified = &newLastModified
	return nil
}

// Java: JSCompZipFileCache.CachedZipFile.maybeClose()
func (z *JSCompZipFileCache_CachedZipFile) MaybeClose() {
	if err := z.maybeCloseErr(); err != nil {
		panic(err)
	}
}

func (z *JSCompZipFileCache_CachedZipFile) maybeCloseErr() jutil.IIOException {
	z.mu.Lock()
	defer z.mu.Unlock()
	return z.maybeCloseLockedErr()
}

func (z *JSCompZipFileCache_CachedZipFile) maybeCloseLockedErr() jutil.IIOException {
	if z.zipFile != nil {
		if err := z.zipFile.Close(); err != nil {
			return jSCompZipFileCache_ioException(err, z.path.String())
		}
		z.zipFile = nil
		z.entries = nil
	}
	return nil
}
