// Port of com/google/javascript/jscomp/deps/PathUtil.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"os"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PathUtil holds utility methods for manipulation of UNIX-like paths.
//
// TODO(tbreisacher): We can probably nuke this entire class and replace it
// with https://docs.oracle.com/javase/8/docs/api/java/nio/file/Path.html
//
// Go-only: language-adapter — the CharMatchers and splits of the Java code
// work on bytes here: '/' and '.' are ASCII, so every split point and index
// is at the same character.
type PathUtil struct{}

// PathUtil_CollapseDots removes all ../ and ./ entries from within the
// given path. If there are extra ..s that move beyond the first directory
// given, they are removed.
//
// Examples:
//
//	"a/b/../c" results in "a/c"
//	"./foo/./../bar" results in "bar"
//	"a/.." results in ""
//	"a/../../foo" results in "foo"
//
// Java: PathUtil.collapseDots(String)
func PathUtil_CollapseDots(path string) string {
	path = PathUtil_RemoveExtraneousSlashes(path)
	// Optimization: Most paths don't contain dots.
	if !strings.Contains(path, ".") {
		return path
	}

	dstFragments := []string{}
	for _, fragment := range strings.Split(path, "/") {
		if fragment == ".." {
			if len(dstFragments) != 0 {
				dstFragments = dstFragments[:len(dstFragments)-1]
			}
		} else if fragment != "." {
			dstFragments = append(dstFragments, fragment)
		}
	}

	// Special case for Join.join([""]); -> "/"
	if len(dstFragments) == 1 && dstFragments[0] == "" {
		return "/"
	}
	return strings.Join(dstFragments, "/")
}

// PathUtil_IsAbsolute determines if a path is absolute or not by testing
// for the presence of "/" at the front of the string.
//
// Java: PathUtil.isAbsolute(String)
func PathUtil_IsAbsolute(path string) bool {
	return strings.HasPrefix(path, "/")
}

// PathUtil_RemoveExtraneousSlashes removes extra slashes from a path.
// Leading slash is preserved, trailing slash is stripped, and any runs of
// more than one slash in the middle is replaced by a single slash.
//
// Java: PathUtil.removeExtraneousSlashes(String)
func PathUtil_RemoveExtraneousSlashes(s string) string {
	// NON_SLASH_MATCHER.lastIndexIn(s)
	lastNonSlash := strings.LastIndexFunc(s, func(r rune) bool { return r != '/' })
	if lastNonSlash != -1 {
		// LastIndexFunc gives the first byte of the last non-slash
		// character; keep all of its bytes.
		end := lastNonSlash + 1
		for end < len(s) && s[end] != '/' {
			end++
		}
		s = s[:end]
	}

	// SLASH_MATCHER.collapseFrom(s, '/')
	var sb strings.Builder
	sb.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			sb.WriteByte('/')
			for i+1 < len(s) && s[i+1] == '/' {
				i++
			}
			continue
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}

// PathUtil_MakeAbsolute converts the given path into an absolute one. This
// prepends the current working directory and removes all .'s from the
// path. If an absolute path is given, it will not be prefixed.
//
// Unlike File.getAbsolutePath(), this function does remove .'s from the
// path and unlike File.getCanonicalPath(), this function does not resolve
// symlinks and does not use filesystem calls.
//
// Java: PathUtil.makeAbsolute(String)
func PathUtil_MakeAbsolute(path string) string {
	// Go-only: language-adapter — StandardSystemProperty.USER_DIR.value()
	// is the process working directory.
	userDir, err := os.Getwd()
	if err != nil {
		panic(jutil.NewRuntimeException(err.Error()))
	}
	return PathUtil_MakeAbsoluteWithRoot(path, userDir)
}

// PathUtil_MakeAbsoluteWithRoot converts the given path into an absolute
// one. This prepends the given rootPath and removes all .'s from the path.
// If an absolute path is given, it will not be prefixed.
//
// Unlike File.getAbsolutePath(), this function does remove .'s from the
// path and unlike File.getCanonicalPath(), this function does not resolve
// symlinks and does not use filesystem calls.
//
// Java: PathUtil.makeAbsolute(String, String)
func PathUtil_MakeAbsoluteWithRoot(path string, rootPath string) string {
	if !PathUtil_IsAbsolute(path) {
		path = rootPath + "/" + path
	}
	return PathUtil_CollapseDots(path)
}

// PathUtil_MakeRelative returns targetPath relative to basePath.
//
// basePath and targetPath must either both be relative, or both be
// absolute paths.
//
// This function is different from makeRelative in that it is able to add
// in ../ components and collapse existing ones as well.
//
// Examples:
//
//	base="some/relative/path" target="some/relative/path/foo" return="foo"
//	base="some/relative/path" target="some/relative" return=".."
//	base="some/relative/path" target="foo/bar" return="../../../foo/bar"
//	base="/some/abs/path" target="/foo/bar" return="../../../foo/bar"
//
// Returns a path relative to targetPath. The returned value will never
// start with a slash.
//
// Java: PathUtil.makeRelative(String, String)
func PathUtil_MakeRelative(basePath string, targetPath string) string {
	// Ensure the paths are both absolute or both relative.
	if PathUtil_IsAbsolute(basePath) != PathUtil_IsAbsolute(targetPath) {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Paths must both be relative or both absolute.\n" +
				"  basePath: " + basePath + "\n" +
				"  targetPath: " + targetPath))
	}

	basePath = PathUtil_CollapseDots(basePath)
	targetPath = PathUtil_CollapseDots(targetPath)
	baseFragments := jutil.StringSplitLiteral(basePath, "/")
	targetFragments := jutil.StringSplitLiteral(targetPath, "/")

	i := -1
	for {
		i += 1
		if i == len(baseFragments) && i == len(targetFragments) {
			// Eg) base: /java/com/google
			//   target: /java/com/google
			//   result: .  <-- . is better than "" since "" + "/path" = "/path"
			return "."
		} else if i == len(baseFragments) {
			// Eg) base: /java/com/google
			//   target: /java/com/google/c/ui
			//   result: c/ui
			return strings.Join(targetFragments[i:], "/")
		} else if i == len(targetFragments) {
			// Eg) base: /java/com/google/c/ui
			//   target: /java/com/google
			//   result: ../..
			return strings.Repeat("../", len(baseFragments)-i-1) + ".."
		}
		if baseFragments[i] != targetFragments[i] {
			break
		}
	}

	// Eg) base: /java/com/google/c
	//   target: /java/com/google/common/base
	//   result: ../common/base
	return strings.Repeat("../", len(baseFragments)-i) +
		strings.Join(targetFragments[i:], "/")
}
