// Port of com/google/javascript/jscomp/resources/ResourceLoader.java @ closure-compiler v20261104

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// Package resources provides embedded Closure Compiler runtime libraries,
// polyfill/transpilation tables, builtin externs, and precompiled TypedASTs.
package resources

import (
	"embed"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JsDir is the embedded directory root for Closure JS runtime libraries and tables.
const JsDir = "com/google/javascript/jscomp/js/"

// classpath holds com/google/javascript/jscomp/js/** and externs.zip from the
// pinned closure-compiler jar.
//
//go:embed com externs.zip
var classpath embed.FS

// runtimeLibsTypedAst is /runtime_libs.typedast, embedded directly as []byte so
// it is read in place without copying its 3.5 MB onto the heap.
//
//go:embed runtime_libs.typedast
var runtimeLibsTypedAst []byte

const runtimeLibsTypedAstName = "runtime_libs.typedast"

var (
	overridesMu    sync.RWMutex
	overridesMap   map[string][]byte
	overridesEpoch atomic.Uint64
)

// RegisterOverride registers an in-memory resource override at path (with any
// leading '/' stripped), taking precedence over embedded resources such as
// "runtime_libs.typedast" or files under JsDir.
func RegisterOverride(path string, data []byte) {
	name := strings.TrimPrefix(path, "/")
	cloned := slices.Clone(data)
	if cloned == nil {
		cloned = []byte{}
	}
	overridesMu.Lock()
	if overridesMap == nil {
		overridesMap = make(map[string][]byte)
	}
	overridesMap[name] = cloned
	overridesEpoch.Add(1)
	overridesMu.Unlock()
}

// UnregisterOverride removes any in-memory resource override registered at path.
func UnregisterOverride(path string) {
	name := strings.TrimPrefix(path, "/")
	overridesMu.Lock()
	if len(overridesMap) > 0 {
		delete(overridesMap, name)
		overridesEpoch.Add(1)
	}
	overridesMu.Unlock()
}

// ClearOverrides removes all registered in-memory resource overrides.
func ClearOverrides() {
	overridesMu.Lock()
	if len(overridesMap) > 0 {
		clear(overridesMap)
		overridesEpoch.Add(1)
	}
	overridesMu.Unlock()
}

// OverridesEpoch returns a monotonic counter incremented whenever in-memory
// resource overrides change.
func OverridesEpoch() uint64 {
	return overridesEpoch.Load()
}

// LoadBytes returns the bytes of the embedded resource at path (with an
// optional leading '/' stripped), or an error wrapping fs.ErrNotExist.
func LoadBytes(path string) ([]byte, error) {
	data, ok := getResource(strings.TrimPrefix(path, "/"))
	if !ok {
		return nil, fmt.Errorf("no such resource %q: %w", path, fs.ErrNotExist)
	}
	return data, nil
}

// MustLoadBytes returns the bytes of the embedded resource at path, panicking
// with a RuntimeException if the resource does not exist.
func MustLoadBytes(path string) []byte {
	data, ok := getResource(strings.TrimPrefix(path, "/"))
	if !ok {
		panic(jutil.NewRuntimeException("No such resource: " + path))
	}
	return data
}

// LoadText returns the UTF-8 text of the embedded resource at path, panicking
// with a RuntimeException if the resource does not exist.
func LoadText(path string) string {
	return string(MustLoadBytes(path))
}

// Exists reports whether the embedded resource at path exists.
func Exists(path string) bool {
	_, ok := getResource(strings.TrimPrefix(path, "/"))
	return ok
}

func getResource(name string) ([]byte, bool) {
	if overridesEpoch.Load() != 0 {
		overridesMu.RLock()
		data, ok := overridesMap[name]
		overridesMu.RUnlock()
		if ok {
			return data, true
		}
	}
	if name == runtimeLibsTypedAstName {
		return runtimeLibsTypedAst, true
	}
	if data, err := fs.ReadFile(classpath, name); err == nil {
		return data, true
	}
	dir := strings.TrimSuffix(name, "/")
	if info, err := fs.Stat(classpath, dir); err == nil && info.IsDir() {
		return []byte{}, true
	}
	return nil, false
}
