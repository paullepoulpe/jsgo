// Port of com/google/javascript/jscomp/ByPathWarningsGuard.java @ closure-compiler v20261015

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ByPathWarningsGuard is an implementation of a WarningsGuard that can
// modify the CheckLevel based on the file that caused the warning, and
// whether this file matches a set of paths (specified either as include or
// exclude of path name parts).
//
// For example:
//
//	List<String> paths = new ArrayList<String>();
//	paths.add("foo");
//	WarningsGuard guard =
//	    ByPathWarningsGuard.forPath(paths, CheckLevel.ERROR, 1);
//
// This guard will convert any warning that came from a file that contains
// "foo" in its path to an error.
type ByPathWarningsGuard struct {
	WarningsGuardBase
	paths    []string
	include  bool
	priority int
	level    CheckLevel
}

var _ WarningsGuard = (*ByPathWarningsGuard)(nil)

const byPathWarningsGuard_serialVersionUID int64 = 1

// newByPathWarningsGuard constructs a new instance. The priority is
// determined by the CheckLevel: ERROR have Priority.STRICT, and OFF have
// priority FILTER_BY_PATH.
//
// Use ByPathWarningsGuard_ForPath or ByPathWarningsGuard_ExceptPath to
// actually create a new instance.
//
// Java: ByPathWarningsGuard(List<String>, boolean, CheckLevel)
func newByPathWarningsGuard(
	paths []string, include bool, level CheckLevel) *ByPathWarningsGuard {
	jutil.CheckArgument(paths != nil)
	jutil.CheckArgument(level == CheckLevel_OFF || level == CheckLevel_ERROR)
	g := &ByPathWarningsGuard{paths: paths, include: include, level: level}
	if level == CheckLevel_ERROR {
		g.priority = WarningsGuard_Priority_STRICT.value()
	} else {
		g.priority = WarningsGuard_Priority_FILTER_BY_PATH.value()
	}
	return g
}

// ByPathWarningsGuard_ForPath returns a new ByPathWarningsGuard that would
// affect any file in the given set of paths (paths for matching), applying
// level to the affected files.
//
// Java: ByPathWarningsGuard.forPath(List<String>, CheckLevel)
func ByPathWarningsGuard_ForPath(
	paths []string, level CheckLevel) *ByPathWarningsGuard {
	return newByPathWarningsGuard(paths, true, level)
}

// ByPathWarningsGuard_ExceptPath returns a new ByPathWarningsGuard that
// would affect any file not in the given set of paths (paths for matching),
// applying level to the affected files.
//
// Java: ByPathWarningsGuard.exceptPath(List<String>, CheckLevel)
func ByPathWarningsGuard_ExceptPath(
	paths []string, level CheckLevel) *ByPathWarningsGuard {
	return newByPathWarningsGuard(paths, false, level)
}

// Level may return nil.
//
// Java: ByPathWarningsGuard.level(JSError)
func (g *ByPathWarningsGuard) Level(error *JSError) *CheckLevel {
	errorPath := error.SourceName()
	defaultLevel := error.DefaultLevel()
	if defaultLevel != CheckLevel_ERROR && errorPath != nil {
		inPath := false
		for _, path := range g.paths {
			inPath = inPath || strings.Contains(*errorPath, path)
		}
		if inPath == g.include {
			return g.level.Ptr()
		}
	}
	return nil
}

// Java: ByPathWarningsGuard.getPriority()
func (g *ByPathWarningsGuard) GetPriority() int {
	return g.priority
}

// String is Object.toString().
func (g *ByPathWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ByPathWarningsGuard", g)
}
