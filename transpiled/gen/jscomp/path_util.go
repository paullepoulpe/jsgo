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
// Port of com/google/javascript/jscomp/deps/PathUtil.java (generated mechanically by j2go).

package jscomp

import (
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.PathUtil
type PathUtil struct {
	idHash_ atomic.Int32
}

func (this_ *PathUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *PathUtil) IsJavaNil_() bool { return this_ == nil }
func (*PathUtil) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.PathUtil", "PathUtil"
}
func (this_ *PathUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PathUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.PathUtil", "com.google.javascript.jscomp.deps.PathUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*PathUtil](v_)
	})
}

var static_PathUtil_SLASH_MATCHER j2gort.Ext_com_google_common_base_CharMatcher

func StaticGet_PathUtil_SLASH_MATCHER() j2gort.Ext_com_google_common_base_CharMatcher {
	clinit_PathUtil()
	return static_PathUtil_SLASH_MATCHER
}
func StaticSet_PathUtil_SLASH_MATCHER(v j2gort.Ext_com_google_common_base_CharMatcher) {
	clinit_PathUtil()
	static_PathUtil_SLASH_MATCHER = v
}

var static_PathUtil_NON_SLASH_MATCHER j2gort.Ext_com_google_common_base_CharMatcher

func StaticGet_PathUtil_NON_SLASH_MATCHER() j2gort.Ext_com_google_common_base_CharMatcher {
	clinit_PathUtil()
	return static_PathUtil_NON_SLASH_MATCHER
}
func StaticSet_PathUtil_NON_SLASH_MATCHER(v j2gort.Ext_com_google_common_base_CharMatcher) {
	clinit_PathUtil()
	static_PathUtil_NON_SLASH_MATCHER = v
}

var clinitDone_PathUtil uint32
var clinitState_PathUtil j2gort.ClinitState

func clinit_PathUtil() {
	if atomic.LoadUint32(&clinitDone_PathUtil) != 0 {
		return
	}
	clinitSlow_PathUtil()
}

func clinitSlow_PathUtil() {
	if !j2gort.EnterClinit(&clinitState_PathUtil, &clinitDone_PathUtil) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PathUtil, &clinitDone_PathUtil)
	static_PathUtil_SLASH_MATCHER = j2gort.Ext_com_google_common_base_CharMatcher_S_Is(uint16(47))
	static_PathUtil_NON_SLASH_MATCHER = j2gort.Ext_com_google_common_base_CharMatcher_S_IsNot(uint16(47))
}

func (this_ *PathUtil) initFields_PathUtil() {
	_ = this_
}

// Java constructor: PathUtil()
func NewPathUtil() *PathUtil {
	clinit_PathUtil()
	this_ := &PathUtil{}
	InitPathUtil(this_)
	return this_
}

func InitPathUtil(this_ *PathUtil) {
	this_.initFields_PathUtil()
}

// Java static method: PathUtil.collapseDots(java.lang.String path)
func PathUtil_CollapseDots(path string) string {
	clinit_PathUtil()
	_ = path
	path = PathUtil_RemoveExtraneousSlashes(path)
	if !(j2gort.JStrIndexOf(path, ".") >= 0) {
		return path
	}
	var dstFragments []string = []string{}
	_ = dstFragments
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(47)).Split(j2gort.BoxString(path))))
		for cur_1.HasNext() {
			fragment := cur_1.NextStr()
			_ = fragment
			if j2gort.CheckNotNullString(fragment) == ".." {
				if len(dstFragments) != 0 {
					dstFragments = j2gort.SliceRemoveAt(dstFragments, (int32(len(dstFragments)) - int32(1)))
				}
			} else if !(j2gort.CheckNotNullString(fragment) == ".") {
				dstFragments = append(dstFragments, fragment)
			}
		}
	}
	if (int32(len(dstFragments)) == int32(1)) && (j2gort.CheckNotNullString(dstFragments[int32(0)]) == "") {
		return "/"
	}
	return j2gort.Ext_com_google_common_base_Joiner_On_1_string("/").Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(dstFragments))
}

// Java static method: PathUtil.isAbsolute(java.lang.String path)
func PathUtil_IsAbsolute(path string) bool {
	clinit_PathUtil()
	_ = path
	return javalang.StringStartsWith(j2gort.CheckNotNullString(path), "/")
}

// Java static method: PathUtil.removeExtraneousSlashes(java.lang.String s)
func PathUtil_RemoveExtraneousSlashes(s string) string {
	clinit_PathUtil()
	_ = s
	var lastNonSlash int32 = StaticGet_PathUtil_NON_SLASH_MATCHER().LastIndexIn(j2gort.BoxString(s))
	_ = lastNonSlash
	if lastNonSlash != int32(-1) {
		s = j2gort.JStrSub(s, int32(0), (lastNonSlash + int32(1)))
	}
	return StaticGet_PathUtil_SLASH_MATCHER().CollapseFrom(j2gort.BoxString(s), uint16(47))
}

// Java static method: PathUtil.makeAbsolute(java.lang.String path)
func PathUtil_MakeAbsolute_1(path string) string {
	clinit_PathUtil()
	_ = path
	return PathUtil_MakeAbsolute_2(path, j2gort.Ext_com_google_common_base_StandardSystemProperty_F_USER_DIR.Value())
}

// Java static method: PathUtil.makeAbsolute(java.lang.String path, java.lang.String rootPath)
func PathUtil_MakeAbsolute_2(path string, rootPath string) string {
	clinit_PathUtil()
	_ = path
	_ = rootPath
	if !(PathUtil_IsAbsolute(path)) {
		path = (j2gort.StrValString(rootPath) + "/" + j2gort.StrValString(path))
	}
	return PathUtil_CollapseDots(path)
}

// Java static method: PathUtil.makeRelative(java.lang.String basePath, java.lang.String targetPath)
func PathUtil_MakeRelative(basePath string, targetPath string) string {
	clinit_PathUtil()
	_ = basePath
	_ = targetPath
	if PathUtil_IsAbsolute(basePath) != PathUtil_IsAbsolute(targetPath) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(("Paths must both be relative or both absolute.\n  basePath: " + j2gort.StrValString(basePath) + "\n  targetPath: " + j2gort.StrValString(targetPath))))))
	}
	basePath = PathUtil_CollapseDots(basePath)
	targetPath = PathUtil_CollapseDots(targetPath)
	var baseFragments []string = j2gort.JStrSplit(basePath, "/")
	_ = baseFragments
	var targetFragments []string = j2gort.JStrSplit(targetPath, "/")
	_ = targetFragments
	var i int32 = int32(-1)
	_ = i
	for doFirst_2 := true; doFirst_2 || (func() bool {
		sr_3 := baseFragments[i]
		se_4 := targetFragments[i]
		return j2gort.CheckNotNullString(sr_3) == se_4
	}()); doFirst_2 = false {
		i = (i + int32(1))
		if (i == int32(len(baseFragments))) && (i == int32(len(targetFragments))) {
			return "."
		} else if i == int32(len(baseFragments)) {
			return j2gort.Ext_com_google_common_base_Joiner_On_1_string("/").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_Lists_NewArrayList_1_Ext_java_util_Iterator(j2gort.AsIface[j2gort.Ext_java_util_ListIterator](j2gort.Ext_java_util_Arrays_AsList(func() []any {
				arr_5 := targetFragments
				if arr_5 == nil {
					return nil
				}
				out_6 := make([]any, len(arr_5))
				for i_, el_7 := range arr_5 {
					out_6[i_] = j2gort.BoxString(el_7)
				}
				return out_6
			}()).ListIterator_1(i)))))
		} else if i == int32(len(targetFragments)) {
			return (j2gort.StrValString(strings.Repeat("../", int(((int32(len(baseFragments))-i)-int32(1))))) + "..")
		}
	}
	return (j2gort.StrValString(strings.Repeat("../", int((int32(len(baseFragments))-i)))) + j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string("/").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_Lists_NewArrayList_1_Ext_java_util_Iterator(j2gort.AsIface[j2gort.Ext_java_util_ListIterator](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_8 := targetFragments
		if arr_8 == nil {
			return nil
		}
		out_9 := make([]any, len(arr_8))
		for i_, el_10 := range arr_8 {
			out_9[i_] = j2gort.BoxString(el_10)
		}
		return out_9
	}()).ListIterator_1(i)))))))
}
