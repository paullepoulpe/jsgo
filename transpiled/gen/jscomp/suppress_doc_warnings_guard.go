/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SuppressDocWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.SuppressDocWarningsGuard
type SuppressDocWarningsGuard struct {
	WarningsGuardBase
	suppressors j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"suppressors"`
	compiler    AbstractCompiler                                  `gson:"compiler"`
}

func (this_ *SuppressDocWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *SuppressDocWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*SuppressDocWarningsGuard) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SuppressDocWarningsGuard", "SuppressDocWarningsGuard"
}
func (this_ *SuppressDocWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*SuppressDocWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SuppressDocWarningsGuard", "com.google.javascript.jscomp.SuppressDocWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*SuppressDocWarningsGuard](v_)
	})
}
func (*SuppressDocWarningsGuard) Implements_jscomp_WarningsGuard() {}

func StaticGet_SuppressDocWarningsGuard_SerialVersionUID() int64 { return int64(1) }
func (this_ *SuppressDocWarningsGuard) initFields_SuppressDocWarningsGuard() {
	_ = this_
}

// Java constructor: SuppressDocWarningsGuard(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.Map<java.lang.String,com.google.javascript.jscomp.DiagnosticGroup> suppressors)
func NewSuppressDocWarningsGuard(compiler AbstractCompiler, suppressors j2gort.Ext_java_util_Map) *SuppressDocWarningsGuard {
	this_ := &SuppressDocWarningsGuard{}
	InitSuppressDocWarningsGuard(this_, compiler, suppressors)
	return this_
}

func InitSuppressDocWarningsGuard(this_ *SuppressDocWarningsGuard, compiler AbstractCompiler, suppressors j2gort.Ext_java_util_Map) {
	_ = compiler
	_ = suppressors
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_SuppressDocWarningsGuard()
	this_.compiler = compiler
	this_.suppressors = SuppressDocWarningsGuard_CreateSuppressors(suppressors)
}

// Java static method: SuppressDocWarningsGuard.createSuppressors(java.util.Map<java.lang.String,com.google.javascript.jscomp.DiagnosticGroup> suppressors)
func SuppressDocWarningsGuard_CreateSuppressors(suppressors j2gort.Ext_java_util_Map) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = suppressors
	var builder j2gort.OrderedMap[string, *DiagnosticGroup] = j2gort.OrderedMapFromJavaMap[string, *DiagnosticGroup](suppressors)
	_ = builder
	builder.Put("missingProperties", NewDiagnosticGroup_1_DiagnosticGroup([]*DiagnosticGroup{StaticGet_DiagnosticGroups_MISSING_PROPERTIES(), StaticGet_DiagnosticGroups_STRICT_MISSING_PROPERTIES()}))
	builder.Put("checkTypes", NewDiagnosticGroup_1_DiagnosticGroup([]*DiagnosticGroup{StaticGet_DiagnosticGroups_CHECK_TYPES(), StaticGet_DiagnosticGroups_STRICT_CHECK_TYPES()}))
	return builder.ToImmutableMap()
}

// Java method: SuppressDocWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *SuppressDocWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	var node pkg_rhino.Node = error_.Node()
	_ = node
	if node == 0 {
		node = this_.priv_SuppressDocWarningsGuard_getScriptNodeBySourceName_5221(error_)
	}
	if node == 0 {
		return nil
	}
	var level *CheckLevel = this_.priv_SuppressDocWarningsGuard_getCheckLevelFromAncestors_5219(error_, node)
	_ = level
	if level != nil {
		return level
	}
	var scriptNode pkg_rhino.Node = this_.priv_SuppressDocWarningsGuard_getScriptNodeBySourceName_5221(error_)
	_ = scriptNode
	if scriptNode != 0 {
		var info *pkg_rhino.JSDocInfo = scriptNode.GetJSDocInfo()
		_ = info
		if info != nil {
			return this_.priv_SuppressDocWarningsGuard_getCheckLevelFromInfo_5220(error_, info)
		}
	}
	return nil
}

// Java method: SuppressDocWarningsGuard.getCheckLevelFromAncestors(com.google.javascript.jscomp.JSError error, com.google.javascript.rhino.Node node)
func (this_ *SuppressDocWarningsGuard) priv_SuppressDocWarningsGuard_getCheckLevelFromAncestors_5219(error_ *JSError, node pkg_rhino.Node) *CheckLevel {
	_ = error_
	_ = node
	{
		var current pkg_rhino.Node = node
		_ = current
		for ; current != 0; current = current.GetParent() {
			var info *pkg_rhino.JSDocInfo = nil
			_ = info
			if current.IsFunction() || current.IsClass() {
				info = NodeUtil_GetBestJSDocInfo(current)
			} else if current.IsScript() {
				info = current.GetJSDocInfo()
			} else if ((((NodeUtil_IsNameDeclaration(current) || NodeUtil_MayBeObjectLitKey(current)) || current.IsComputedProp()) || current.IsMemberFieldDef()) || current.IsComputedFieldDef()) || (current.HasParent() && current.GetParent().IsExprResult()) {
				info = NodeUtil_GetBestJSDocInfo(current)
			}
			if info != nil {
				var level *CheckLevel = this_.priv_SuppressDocWarningsGuard_getCheckLevelFromInfo_5220(error_, info)
				_ = level
				if level != nil {
					return level
				}
			}
		}
	}
	return nil
}

// Java method: SuppressDocWarningsGuard.getCheckLevelFromInfo(com.google.javascript.jscomp.JSError error, com.google.javascript.rhino.JSDocInfo info)
func (this_ *SuppressDocWarningsGuard) priv_SuppressDocWarningsGuard_getCheckLevelFromInfo_5220(error_ *JSError, info *pkg_rhino.JSDocInfo) *CheckLevel {
	_ = error_
	_ = info
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetSuppressions()))
		for cur_1.HasNext() {
			suppressor := cur_1.NextStr()
			_ = suppressor
			var group *DiagnosticGroup = j2gort.CastPtr[DiagnosticGroup](j2gort.MapGetStr(this_.suppressors, suppressor))
			_ = group
			if group == nil {
				continue
			}
			if group.Matches_1_JSError(error_) {
				return StaticGet_CheckLevel_OFF()
			}
		}
	}
	return nil
}

// Java method: SuppressDocWarningsGuard.getScriptNodeBySourceName(com.google.javascript.jscomp.JSError error)
func (this_ *SuppressDocWarningsGuard) priv_SuppressDocWarningsGuard_getScriptNodeBySourceName_5221(error_ *JSError) pkg_rhino.Node {
	_ = error_
	if error_.SourceName() == j2gort.NullString {
		return pkg_rhino.Node(0)
	}
	var scriptNode pkg_rhino.Node = this_.compiler.GetScriptNode(error_.SourceName())
	_ = scriptNode
	if scriptNode == 0 {
		return pkg_rhino.Node(0)
	}
	if !(scriptNode.IsScript()) {
		j2gort.CheckState(false)
	}
	return scriptNode
}

// Java method: SuppressDocWarningsGuard.getPriority()
func (this_ *SuppressDocWarningsGuard) GetPriority() int32 {
	return (StaticGet_WarningsGuard_Priority_SUPPRESS_DOC()).FieldGet_value()
}
