// Port of com/google/javascript/jscomp/ShowByPathWarningsGuard.java @ closure-compiler v20260415

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ShowByPathWarningsGuard controls whether warnings should be restricted or
// suppressed for specified paths.
type ShowByPathWarningsGuard struct {
	WarningsGuardBase
	warningsGuard *ByPathWarningsGuard
}

var _ WarningsGuard = (*ShowByPathWarningsGuard)(nil)

const showByPathWarningsGuard_serialVersionUID int64 = 1

// ShowByPathWarningsGuard_ShowType controls whether warnings should be
// restricted to a specified path or suppressed within the specified path.
//
// Java: ShowByPathWarningsGuard.ShowType
type ShowByPathWarningsGuard_ShowType int32

const (
	ShowByPathWarningsGuard_ShowType_INCLUDE ShowByPathWarningsGuard_ShowType = iota // Suppress warnings outside the path.
	ShowByPathWarningsGuard_ShowType_EXCLUDE                                         // Suppress warnings within the path.
)

var showByPathWarningsGuard_ShowType_names = [...]string{"INCLUDE", "EXCLUDE"}

// Java: ShowByPathWarningsGuard.ShowType.values()
func ShowByPathWarningsGuard_ShowType_values() []ShowByPathWarningsGuard_ShowType {
	return []ShowByPathWarningsGuard_ShowType{ShowByPathWarningsGuard_ShowType_INCLUDE, ShowByPathWarningsGuard_ShowType_EXCLUDE}
}

// String is Enum.name().
func (t ShowByPathWarningsGuard_ShowType) String() string {
	return showByPathWarningsGuard_ShowType_names[t]
}

// Java: ShowByPathWarningsGuard(String)
func NewShowByPathWarningsGuard(checkWarningsOnlyForPath string) *ShowByPathWarningsGuard {
	return NewShowByPathWarningsGuardWithShowType(checkWarningsOnlyForPath, ShowByPathWarningsGuard_ShowType_INCLUDE)
}

// Java: ShowByPathWarningsGuard(String[])
func NewShowByPathWarningsGuardWithPaths(checkWarningsOnlyForPath []string) *ShowByPathWarningsGuard {
	return NewShowByPathWarningsGuardWithPathsAndShowType(checkWarningsOnlyForPath, ShowByPathWarningsGuard_ShowType_INCLUDE)
}

// Java: ShowByPathWarningsGuard(String, ShowType)
func NewShowByPathWarningsGuardWithShowType(path string, showType ShowByPathWarningsGuard_ShowType) *ShowByPathWarningsGuard {
	return NewShowByPathWarningsGuardWithPathsAndShowType([]string{path}, showType)
}

// Java: ShowByPathWarningsGuard(String[], ShowType)
func NewShowByPathWarningsGuardWithPathsAndShowType(paths []string, showType ShowByPathWarningsGuard_ShowType) *ShowByPathWarningsGuard {
	jutil.CheckArgument(paths != nil)
	g := &ShowByPathWarningsGuard{}
	pathList := append([]string{}, paths...)
	if showType == ShowByPathWarningsGuard_ShowType_INCLUDE {
		g.warningsGuard = ByPathWarningsGuard_ExceptPath(pathList, CheckLevel_OFF)
	} else {
		g.warningsGuard = ByPathWarningsGuard_ForPath(pathList, CheckLevel_OFF)
	}
	return g
}

// Java: ShowByPathWarningsGuard.level(JSError)
func (g *ShowByPathWarningsGuard) Level(error *JSError) *CheckLevel {
	return g.warningsGuard.Level(error)
}

// Java: ShowByPathWarningsGuard.getPriority()
func (g *ShowByPathWarningsGuard) GetPriority() int {
	return g.warningsGuard.GetPriority()
}

// String is Object.toString().
func (g *ShowByPathWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ShowByPathWarningsGuard", g)
}
