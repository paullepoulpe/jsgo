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
// Port of com/google/javascript/jscomp/ByPathWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.ByPathWarningsGuard
type ByPathWarningsGuard struct {
	WarningsGuardBase
	paths    j2gort.Ext_java_util_List `gson:"paths"`
	level    *CheckLevel               `gson:"level"`
	priority int32                     `gson:"priority"`
	include  bool                      `gson:"include"`
}

func (this_ *ByPathWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *ByPathWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*ByPathWarningsGuard) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ByPathWarningsGuard", "ByPathWarningsGuard"
}
func (this_ *ByPathWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*ByPathWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ByPathWarningsGuard", "com.google.javascript.jscomp.ByPathWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*ByPathWarningsGuard](v_)
	})
}
func (*ByPathWarningsGuard) Implements_jscomp_WarningsGuard() {}

func StaticGet_ByPathWarningsGuard_SerialVersionUID() int64 { return int64(1) }
func (this_ *ByPathWarningsGuard) initFields_ByPathWarningsGuard() {
	_ = this_
}

// Java constructor: ByPathWarningsGuard(java.util.List<java.lang.String> paths, boolean include, com.google.javascript.jscomp.CheckLevel level)
func NewByPathWarningsGuard(paths j2gort.Ext_java_util_List, include bool, level *CheckLevel) *ByPathWarningsGuard {
	this_ := &ByPathWarningsGuard{}
	InitByPathWarningsGuard(this_, paths, include, level)
	return this_
}

func InitByPathWarningsGuard(this_ *ByPathWarningsGuard, paths j2gort.Ext_java_util_List, include bool, level *CheckLevel) {
	_ = paths
	_ = include
	_ = level
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_ByPathWarningsGuard()
	if !(paths != nil) {
		j2gort.CheckArgument(false)
	}
	if !((level == StaticGet_CheckLevel_OFF()) || (level == StaticGet_CheckLevel_ERROR())) {
		j2gort.CheckArgument(false)
	}
	this_.paths = paths
	this_.include = include
	this_.level = level
	if level == StaticGet_CheckLevel_ERROR() {
		this_.priority = (StaticGet_WarningsGuard_Priority_STRICT()).FieldGet_value()
	} else {
		this_.priority = (StaticGet_WarningsGuard_Priority_FILTER_BY_PATH()).FieldGet_value()
	}
}

// Java static method: ByPathWarningsGuard.forPath(java.util.List<java.lang.String> paths, com.google.javascript.jscomp.CheckLevel level)
func ByPathWarningsGuard_ForPath(paths j2gort.Ext_java_util_List, level *CheckLevel) *ByPathWarningsGuard {
	_ = paths
	_ = level
	return NewByPathWarningsGuard(paths, true, level)
}

// Java static method: ByPathWarningsGuard.exceptPath(java.util.List<java.lang.String> paths, com.google.javascript.jscomp.CheckLevel level)
func ByPathWarningsGuard_ExceptPath(paths j2gort.Ext_java_util_List, level *CheckLevel) *ByPathWarningsGuard {
	_ = paths
	_ = level
	return NewByPathWarningsGuard(paths, false, level)
}

// Java method: ByPathWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *ByPathWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	var errorPath string = error_.SourceName()
	_ = errorPath
	var defaultLevel *CheckLevel = error_.DefaultLevel()
	_ = defaultLevel
	if (defaultLevel != StaticGet_CheckLevel_ERROR()) && (errorPath != j2gort.NullString) {
		var inPath bool = false
		_ = inPath
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.paths))
			for cur_1.HasNext() {
				path := cur_1.NextStr()
				_ = path
				{
					ba_2 := inPath
					bb_3 := (j2gort.JStrIndexOf(errorPath, j2gort.CheckNotNullString(path)) >= 0)
					inPath = ba_2 || bb_3
				}
			}
		}
		if inPath == this_.include {
			return this_.level
		}
	}
	return nil
}

// Java method: ByPathWarningsGuard.getPriority()
func (this_ *ByPathWarningsGuard) GetPriority() int32 {
	return this_.priority
}
