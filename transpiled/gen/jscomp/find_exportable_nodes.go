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
// Port of com/google/javascript/jscomp/FindExportableNodes.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FindExportableNodes
type FindExportableNodes struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler          AbstractCompiler                    `gson:"compiler"`
	exports           *j2gort.Ext_java_util_LinkedHashMap `gson:"exports"`
	es6ClassExports   *j2gort.Ext_java_util_LinkedHashMap `gson:"es6ClassExports"`
	localExports      *j2gort.Ext_java_util_LinkedHashSet `gson:"localExports"`
	allowLocalExports bool                                `gson:"allowLocalExports"`
}

func (this_ *FindExportableNodes) BindSelf_(s any) {
	_ = s
}
func (this_ *FindExportableNodes) IsJavaNil_() bool { return this_ == nil }
func (*FindExportableNodes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FindExportableNodes", "FindExportableNodes"
}
func (this_ *FindExportableNodes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*FindExportableNodes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FindExportableNodes", "com.google.javascript.jscomp.FindExportableNodes", func(v_ any) bool {
		return j2gort.InstanceOf[*FindExportableNodes](v_)
	})
}
func (*FindExportableNodes) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*FindExportableNodes) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_FindExportableNodes_NON_GLOBAL_ERROR *DiagnosticType

func StaticGet_FindExportableNodes_NON_GLOBAL_ERROR() *DiagnosticType {
	clinit_FindExportableNodes()
	return static_FindExportableNodes_NON_GLOBAL_ERROR
}
func StaticSet_FindExportableNodes_NON_GLOBAL_ERROR(v *DiagnosticType) {
	clinit_FindExportableNodes()
	static_FindExportableNodes_NON_GLOBAL_ERROR = v
}

var static_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED *DiagnosticType

func StaticGet_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED() *DiagnosticType {
	clinit_FindExportableNodes()
	return static_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED
}
func StaticSet_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED(v *DiagnosticType) {
	clinit_FindExportableNodes()
	static_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED = v
}

var clinitDone_FindExportableNodes uint32
var clinitState_FindExportableNodes j2gort.ClinitState

func clinit_FindExportableNodes() {
	if atomic.LoadUint32(&clinitDone_FindExportableNodes) != 0 {
		return
	}
	clinitSlow_FindExportableNodes()
}

func clinitSlow_FindExportableNodes() {
	if !j2gort.EnterClinit(&clinitState_FindExportableNodes, &clinitDone_FindExportableNodes) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FindExportableNodes, &clinitDone_FindExportableNodes)
	static_FindExportableNodes_NON_GLOBAL_ERROR = DiagnosticType_Error("JSC_NON_GLOBAL_ERROR", "@export only allowed on symbols/properties defined in the global scope.\nConvert to a global definition or enable --export_local_property_definitions/ options.setExportLocalPropertyDefinitions(true)")
	static_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED = DiagnosticType_Error("JSC_EXPORT_ANNOTATION_NOT_ALLOWED", "@export is not supported on this expression.")
}

func (this_ *FindExportableNodes) initFields_FindExportableNodes() {
	_ = this_
	this_.exports = j2gort.NewExt_java_util_LinkedHashMap()
	this_.es6ClassExports = j2gort.NewExt_java_util_LinkedHashMap()
	this_.localExports = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: FindExportableNodes(com.google.javascript.jscomp.AbstractCompiler compiler, boolean allowLocalExports)
func NewFindExportableNodes(compiler AbstractCompiler, allowLocalExports bool) *FindExportableNodes {
	clinit_FindExportableNodes()
	this_ := &FindExportableNodes{}
	InitFindExportableNodes(this_, compiler, allowLocalExports)
	return this_
}

func InitFindExportableNodes(this_ *FindExportableNodes, compiler AbstractCompiler, allowLocalExports bool) {
	_ = compiler
	_ = allowLocalExports
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_FindExportableNodes()
	this_.compiler = compiler
	this_.allowLocalExports = allowLocalExports
}

// Java method: FindExportableNodes.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FindExportableNodes) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var docInfo *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = docInfo
	if (docInfo == nil) || !(docInfo.IsExport()) {
		return
	}
	if parent.IsAssign() && (n.IsFunction() || n.IsClass()) {
		var parentInfo *pkg_rhino.JSDocInfo = parent.GetJSDocInfo()
		_ = parentInfo
		if (parentInfo != nil) && parentInfo.IsExport() {
			return
		}
	}
	var mode FindExportableNodes_Mode = FindExportableNodes_Mode(0)
	_ = mode
	var export string = j2gort.NullString
	_ = export
	var context pkg_rhino.Node = pkg_rhino.Node(0)
	_ = context
	{
		var grandparent pkg_rhino.Node = pkg_rhino.Node(0)
		_ = grandparent
		var child pkg_rhino.Node = pkg_rhino.Node(0)
		_ = child
		switch n.TokenSwitchOrd_() {
		case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */ :
			if parent.IsScript() {
				export = NodeUtil_GetName(n)
				context = n
				mode = StaticGet_FindExportableNodes_Mode_EXPORT()
			}
		case int32(51) /* Token.ASSIGN */ :
			grandparent = parent.GetParent()
			child = n.GetFirstChild()
			if parent.IsExprResult() && !(n.GetLastChild().IsAssign()) {
				if ((grandparent != 0) && grandparent.IsScript()) && child.IsQualifiedName() {
					export = child.GetQualifiedName()
					context = n
					mode = StaticGet_FindExportableNodes_Mode_EXPORT()
				} else if this_.allowLocalExports && child.IsGetProp() {
					export = child.GetString()
					mode = StaticGet_FindExportableNodes_Mode_EXTERN()
				}
			}
		case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
			if parent.IsScript() {
				if n.GetFirstChild().HasChildren() && !(n.GetFirstFirstChild().IsAssign()) {
					export = n.GetFirstChild().GetString()
					context = n
					mode = StaticGet_FindExportableNodes_Mode_EXPORT()
				}
			}
		case int32(26) /* Token.GETPROP */ :
			{
				if this_.allowLocalExports && parent.IsExprResult() {
					export = n.GetString()
					mode = StaticGet_FindExportableNodes_Mode_EXTERN()
				}
				break
			}
		case int32(108) /* Token.MEMBER_FIELD_DEF */ :
			if n.IsStaticMember() {
				var classNode pkg_rhino.Node = parent.GetParent()
				_ = classNode
				var className string = NodeUtil_GetName(classNode)
				_ = className
				if className == j2gort.NullString {
					t.Report_3(n, StaticGet_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED(), []string{})
					return
				}
				export = (j2gort.StrValString(className) + "." + j2gort.StrValString(n.GetString()))
				mode = StaticGet_FindExportableNodes_Mode_EXPORT()
				context = n
			} else if this_.allowLocalExports {
				export = n.GetString()
				mode = StaticGet_FindExportableNodes_Mode_EXTERN()
			}
		case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
			if parent.GetParent().IsClass() {
				var methodOwnerName string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(parent.GetParent()))
				_ = methodOwnerName
				if methodOwnerName == j2gort.NullString {
					t.Report_3(n, StaticGet_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED(), []string{})
					return
				}
				j2gort.MapPutNode(this_.es6ClassExports, n, j2gort.BoxString((j2gort.StrValString(methodOwnerName) + j2gort.StrValString(func() string {
					if n.IsStaticMember() {
						return ""
					}
					return ".prototype"
				}()))))
				return
			}
			fallthrough
		case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
			if this_.allowLocalExports {
				export = n.GetString()
				mode = StaticGet_FindExportableNodes_Mode_EXTERN()
			}
		default:
		}
	}
	if export != j2gort.NullString {
		if mode == StaticGet_FindExportableNodes_Mode_EXPORT() {
			j2gort.CheckNotNullNode(context)
			j2gort.MapPutStr(this_.exports, export, j2gort.BoxNode(context))
		} else {
			if !(context == 0) {
				j2gort.CheckState(false)
			}
			if !(mode == StaticGet_FindExportableNodes_Mode_EXTERN()) {
				j2gort.CheckState(false)
			}
			if !(!(j2gort.CheckNotNullString(export) == "")) {
				j2gort.CheckState(false)
			}
			j2gort.SetAddStr(this_.localExports, export)
		}
	} else if !(n.IsGetProp() && parent.IsExprResult()) {
		if !(n.IsFunction()) || !(NodeUtil_MayBeObjectLitKey(parent)) {
			if this_.allowLocalExports {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED(), []string{}))
			} else {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_FindExportableNodes_NON_GLOBAL_ERROR(), []string{}))
			}
		}
	}
}

// Java method: FindExportableNodes.getExports()
func (this_ *FindExportableNodes) GetExports() *j2gort.Ext_java_util_LinkedHashMap {
	return this_.exports
}

// Java method: FindExportableNodes.getLocalExports()
func (this_ *FindExportableNodes) GetLocalExports() *j2gort.Ext_java_util_LinkedHashSet {
	return this_.localExports
}

// Java method: FindExportableNodes.getEs6ClassExports()
func (this_ *FindExportableNodes) GetEs6ClassExports() *j2gort.Ext_java_util_LinkedHashMap {
	return this_.es6ClassExports
}

// Java enum: com.google.javascript.jscomp.FindExportableNodes.Mode
type FindExportableNodes_Mode uint8

const (
	enum_FindExportableNodes_Mode_EXPORT FindExportableNodes_Mode = 1
	enum_FindExportableNodes_Mode_EXTERN FindExportableNodes_Mode = 2
)

var enumBaseTable_FindExportableNodes_Mode = [...]j2gort.EnumBase{
	{Name_: "EXPORT", Ordinal_: 0},
	{Name_: "EXTERN", Ordinal_: 1},
}
var enumValues_FindExportableNodes_Mode = []FindExportableNodes_Mode{1, 2}

func StaticGet_FindExportableNodes_Mode_EXPORT() FindExportableNodes_Mode {
	return enum_FindExportableNodes_Mode_EXPORT
}
func StaticGet_FindExportableNodes_Mode_EXTERN() FindExportableNodes_Mode {
	return enum_FindExportableNodes_Mode_EXTERN
}

func FindExportableNodes_Mode_Values() []FindExportableNodes_Mode {
	return j2gort.CloneSlice(enumValues_FindExportableNodes_Mode)
}

func FindExportableNodes_Mode_ValueOf(name string) FindExportableNodes_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FindExportableNodes_Mode {
		if enumBaseTable_FindExportableNodes_Mode[i_].Name_ == name {
			return FindExportableNodes_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.FindExportableNodes.Mode." + name))
}

func (this_ FindExportableNodes_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (FindExportableNodes_Mode) HasIdentityEquals_()                                 {}
func (FindExportableNodes_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FindExportableNodes_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FindExportableNodes_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FindExportableNodes_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FindExportableNodes$Mode", "Mode"
}
func (this_ FindExportableNodes_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FindExportableNodes_Mode[this_-1]
}
func (this_ FindExportableNodes_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FindExportableNodes_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FindExportableNodes_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FindExportableNodes_Mode[this_-1].Name_
}
func (this_ FindExportableNodes_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FindExportableNodes_Mode[this_-1].Name_
}
func (this_ FindExportableNodes_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(FindExportableNodes_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FindExportableNodes_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FindExportableNodes_Mode[this_-1].HashCode()
}
func (this_ FindExportableNodes_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FindExportableNodes_Mode[this_-1].CompareTo(other_)
}
func (this_ FindExportableNodes_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FindExportableNodes_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{FindExportableNodes_Mode(1), FindExportableNodes_Mode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FindExportableNodes.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FindExportableNodes$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FindExportableNodes.Mode", "com.google.javascript.jscomp.FindExportableNodes$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[FindExportableNodes_Mode](v_)
	})
}
func (*FindExportableNodes) AlwaysShouldTraverse_()      {}
func (*FindExportableNodes) CanReleaseSyntacticScopes_() {}
