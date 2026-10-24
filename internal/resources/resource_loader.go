// Port of com/google/javascript/jscomp/resources/ResourceLoader.java @ closure-compiler v20261021

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

package resources

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// resourceLoader_classpath holds the classpath resources that
// ResourceLoader's callers read: com/google/javascript/jscomp/js/** of the
// pinned closure-compiler jar (polyfills.txt, transpilation_libs.txt and the
// runtime library .js files), extracted with their original paths, and the
// jar's root-level externs.zip (the builtin externs that
// AbstractCommandLineRunner.getBuiltinExterns reads). The jar's
// /runtime_libs.typedast is resourceLoader_runtimeLibsTypedAst.
//
// Go-only: language-adapter — classpath resources via go:embed.
//
//go:embed com externs.zip
var resourceLoader_classpath embed.FS

// resourceLoader_runtimeLibsTypedAst is /runtime_libs.typedast, the
// precompiled runtime libraries that Compiler.initRuntimeLibraryTypedAsts
// reads through Class.getResourceAsStream.
//
// Go-only: performance — embedded directly as []byte so that it is read in
// place: fs.ReadFile would copy its 3.5 MB onto the heap, where the runtime
// library TypedASTs keep it alive for the whole compilation.
//
//go:embed runtime_libs.typedast
var resourceLoader_runtimeLibsTypedAst []byte

const resourceLoader_runtimeLibsTypedAstName = "runtime_libs.typedast"

// ResourceLoader loads text resources from the classpath.
//
// Java: public final class ResourceLoader (only static members).
//
// Go-only: language-adapter — a Java Class<?> argument is the class's binary
// name (e.g. "com.google.javascript.jscomp.RemoveUnusedCode"); like
// Class.getResource, a relative path is resolved against its package.
type ResourceLoader struct{}

// ResourceLoader_LoadTextResource returns the UTF-8 text of the resource
// path relative to clazz, throwing RuntimeException("No such resource: ...")
// if there is none.
//
// Java: ResourceLoader.loadTextResource(Class<?>, String)
func ResourceLoader_LoadTextResource(clazz string, path string) string {
	var result string
	jutil.Try(func() {
		result = resourceLoader_charStreamsToString(resourceLoader_getResourceAsStream(clazz, path))
	}, func(e jutil.INullPointerException) {
		if !ResourceLoader_ResourceExists(clazz, path) {
			panic(jutil.NewRuntimeException("No such resource: " + path))
		}
		panic(e)
	})
	return result
}

// ResourceLoader_LoadPropertiesMap parses the resource resourceName relative
// to clazz with PropertiesParser.
//
// Java: ResourceLoader.loadPropertiesMap(Class<?>, String)
func ResourceLoader_LoadPropertiesMap(clazz string, resourceName string) *jutil.LinkedHashMap[string, string] {
	return PropertiesParser_Parse(ResourceLoader_LoadTextResource(clazz, resourceName))
}

// ResourceLoader_ResourceExists reports whether the resource path relative
// to clazz exists.
//
// Java: ResourceLoader.resourceExists(Class<?>, String)
func ResourceLoader_ResourceExists(clazz string, path string) bool {
	_, ok := resourceLoader_getResource(resourceLoader_resolveName(clazz, path))
	return ok
}

// NewResourceLoader is the private constructor.
//
// Java: ResourceLoader()
func NewResourceLoader() *ResourceLoader {
	return &ResourceLoader{}
}

// resourceLoader_resolveName is Class.resolveName(String): a relative name is
// prefixed with the package of clazz (dots turned into slashes); a leading
// '/' is removed from an absolute one.
//
// Go-only: language-adapter — java.lang.Class.getResource on a class name.
func resourceLoader_resolveName(clazz string, name string) string {
	if !strings.HasPrefix(name, "/") {
		baseName := "" // Java: getPackageName()
		if i := strings.LastIndexByte(clazz, '.'); i >= 0 {
			baseName = clazz[:i]
		}
		if baseName != "" {
			name = strings.ReplaceAll(baseName, ".", "/") + "/" + name
		}
	} else {
		name = name[1:]
	}
	return name
}

// resourceLoader_getResource looks name up like the class loader looks up a
// jar entry (JarFile.getEntry): the file name, or the directory name + "/"
// (whose stream is empty).
//
// Go-only: language-adapter — ClassLoader.getResource over the embedded
// resources.
func resourceLoader_getResource(name string) ([]byte, bool) {
	if name == resourceLoader_runtimeLibsTypedAstName {
		// Read-only memory: the callers only read it.
		return resourceLoader_runtimeLibsTypedAst, true
	}
	if data, err := fs.ReadFile(resourceLoader_classpath, name); err == nil {
		return data, true
	}
	dir := strings.TrimSuffix(name, "/")
	if info, err := fs.Stat(resourceLoader_classpath, dir); err == nil && info.IsDir() {
		return []byte{}, true
	}
	return nil, false
}

// resourceLoader_getResourceAsStream is clazz.getResourceAsStream(path):
// nil (Java null) when there is no such resource.
//
// Go-only: language-adapter — Class.getResourceAsStream over the embedded
// resources.
func resourceLoader_getResourceAsStream(clazz string, path string) []byte {
	data, ok := resourceLoader_getResource(resourceLoader_resolveName(clazz, path))
	if !ok {
		return nil
	}
	return data
}

// resourceLoader_charStreamsToString is CharStreams.toString(new
// InputStreamReader(stream, UTF_8)); like new InputStreamReader(null, ...)
// it throws NullPointerException for a null stream. The embedded resources
// are valid UTF-8 (resource_loader_test.go checks every one), on which
// Java's decoding with replacement is the identity.
//
// Go-only: language-adapter — java.io stream decoding of embedded bytes.
func resourceLoader_charStreamsToString(stream []byte) string {
	if stream == nil {
		panic(jutil.NewNullPointerException())
	}
	return string(stream)
}

// ResourceLoader_GetResourceAsStream is clazz.getResourceAsStream(path) for
// the callers of Class.getResourceAsStream outside this class (e.g.
// Compiler): the resource's bytes, or nil (Java null) when there is none.
//
// Go-only: language-adapter — Class.getResourceAsStream over the embedded
// resources.
func ResourceLoader_GetResourceAsStream(clazz string, path string) []byte {
	return resourceLoader_getResourceAsStream(clazz, path)
}
