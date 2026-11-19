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
// Port of com/google/javascript/jscomp/RewriteGoogJsImports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RewriteGoogJsImports
type RewriteGoogJsImports struct {
	moduleMap          *ModuleMap                                `gson:"moduleMap"`
	compiler           AbstractCompiler                          `gson:"compiler"`
	googModule         *Module                                   `gson:"googModule"`
	moduleReplacements j2gort.CustomOrderedMap[*Module, *Module] `gson:"moduleReplacements"`
	idHash_            atomic.Int32
	mode               RewriteGoogJsImports_Mode `gson:"mode"`
}

func (this_ *RewriteGoogJsImports) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteGoogJsImports) IsJavaNil_() bool { return this_ == nil }
func (*RewriteGoogJsImports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteGoogJsImports", "RewriteGoogJsImports"
}
func (this_ *RewriteGoogJsImports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteGoogJsImports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteGoogJsImports", "com.google.javascript.jscomp.RewriteGoogJsImports", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteGoogJsImports](v_)
	})
}
func (*RewriteGoogJsImports) Implements_jscomp_CompilerPass() {}

var static_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR *DiagnosticType

func StaticGet_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR() *DiagnosticType {
	clinit_RewriteGoogJsImports()
	return static_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR
}
func StaticSet_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR(v *DiagnosticType) {
	clinit_RewriteGoogJsImports()
	static_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR = v
}

var static_RewriteGoogJsImports_GOOG_JS_REEXPORTED *DiagnosticType

func StaticGet_RewriteGoogJsImports_GOOG_JS_REEXPORTED() *DiagnosticType {
	clinit_RewriteGoogJsImports()
	return static_RewriteGoogJsImports_GOOG_JS_REEXPORTED
}
func StaticSet_RewriteGoogJsImports_GOOG_JS_REEXPORTED(v *DiagnosticType) {
	clinit_RewriteGoogJsImports()
	static_RewriteGoogJsImports_GOOG_JS_REEXPORTED = v
}

var static_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG *DiagnosticType

func StaticGet_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG() *DiagnosticType {
	clinit_RewriteGoogJsImports()
	return static_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG
}
func StaticSet_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG(v *DiagnosticType) {
	clinit_RewriteGoogJsImports()
	static_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG = v
}

var static_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG *DiagnosticType

func StaticGet_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG() *DiagnosticType {
	clinit_RewriteGoogJsImports()
	return static_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG
}
func StaticSet_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG(v *DiagnosticType) {
	clinit_RewriteGoogJsImports()
	static_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG = v
}
func StaticGet_RewriteGoogJsImports_EXPECTED_BASE_PROVIDE() string { return "goog" }

var clinitDone_RewriteGoogJsImports uint32
var clinitState_RewriteGoogJsImports j2gort.ClinitState

func clinit_RewriteGoogJsImports() {
	if atomic.LoadUint32(&clinitDone_RewriteGoogJsImports) != 0 {
		return
	}
	clinitSlow_RewriteGoogJsImports()
}

func clinitSlow_RewriteGoogJsImports() {
	if !j2gort.EnterClinit(&clinitState_RewriteGoogJsImports, &clinitDone_RewriteGoogJsImports) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteGoogJsImports, &clinitDone_RewriteGoogJsImports)
	static_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR = DiagnosticType_Error("JSC_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR", "Closure''s goog.js file must be imported as `import * as goog`.")
	static_RewriteGoogJsImports_GOOG_JS_REEXPORTED = DiagnosticType_Error("JSC_GOOG_JS_REEXPORTED", "Do not re-export from goog.js.")
	static_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG = DiagnosticType_Error("JSC_CANNOT_NAME_FILE_GOOG", "Do not name files goog.js, it is reserved for Closure Library.")
	static_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG = DiagnosticType_Error("JSC_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG", "Module scoped variables named ''goog'' must come from importing Closure Library''s goog.js file..")
}

func (this_ *RewriteGoogJsImports) initFields_RewriteGoogJsImports() {
	_ = this_
	this_.moduleReplacements = j2gort.NewCustomOrderedMap[*Module, *Module]()
}

// Java constructor: RewriteGoogJsImports(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.RewriteGoogJsImports.Mode mode, com.google.javascript.jscomp.modules.ModuleMap moduleMap)
func NewRewriteGoogJsImports(compiler AbstractCompiler, mode RewriteGoogJsImports_Mode, moduleMap *ModuleMap) *RewriteGoogJsImports {
	clinit_RewriteGoogJsImports()
	this_ := &RewriteGoogJsImports{}
	InitRewriteGoogJsImports(this_, compiler, mode, moduleMap)
	return this_
}

func InitRewriteGoogJsImports(this_ *RewriteGoogJsImports, compiler AbstractCompiler, mode RewriteGoogJsImports_Mode, moduleMap *ModuleMap) {
	_ = compiler
	_ = mode
	_ = moduleMap
	this_.initFields_RewriteGoogJsImports()
	j2gort.CheckNotNullPtr(moduleMap)
	this_.compiler = compiler
	this_.mode = mode
	this_.moduleMap = moduleMap
}

// Java method: RewriteGoogJsImports.changeModules()
func (this_ *RewriteGoogJsImports) priv_RewriteGoogJsImports_changeModules_4981() {
	var resolvedModules j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = resolvedModules
	var closureModules j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = closureModules
	{
		cur_1 := j2gort.NewMapEntryCursor(this_.moduleMap.GetModulesByPath())
		for cur_1.NextEntry() {
			var newModule *Module = this_.moduleReplacements.GetOrDefault(j2gort.CastPtr[Module](cur_1.GetValue()), j2gort.CastPtr[Module](cur_1.GetValue()))
			_ = newModule
			resolvedModules.Put_2(j2gort.BoxString(cur_1.KeyStr()), j2gort.PtrToAny(newModule))
		}
	}
	{
		cur_2 := j2gort.NewMapEntryCursor(this_.moduleMap.GetModulesByClosureNamespace())
		for cur_2.NextEntry() {
			var newModule_2 *Module = this_.moduleReplacements.GetOrDefault(j2gort.CastPtr[Module](cur_2.GetValue()), j2gort.CastPtr[Module](cur_2.GetValue()))
			_ = newModule_2
			closureModules.Put_2(j2gort.BoxString(cur_2.KeyStr()), j2gort.PtrToAny(newModule_2))
		}
	}
	this_.compiler.SetModuleMap(NewModuleMap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](resolvedModules.BuildOrThrow()), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](closureModules.BuildOrThrow())))
	this_.moduleReplacements.Clear()
}

// Java method: RewriteGoogJsImports.findGoogImportNode(com.google.javascript.rhino.Node scriptRoot)
func (this_ *RewriteGoogJsImports) priv_RewriteGoogJsImports_findGoogImportNode_4982(scriptRoot pkg_rhino.Node) pkg_rhino.Node {
	_ = scriptRoot
	var valid bool = true
	_ = valid
	var googImportNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = googImportNode
	{
		var child pkg_rhino.Node = scriptRoot.GetFirstFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsImport() && javalang.StringEndsWith(j2gort.CheckNotNullString(child.GetLastChild().GetString()), "/goog.js") {
				if (child.GetFirstChild().IsEmpty() && child.GetSecondChild().IsImportStar()) && (j2gort.CheckNotNullString(child.GetSecondChild().GetString()) == "goog") {
					googImportNode = child
				} else {
					valid = false
					this_.compiler.Report_1(JSError_Make_5(scriptRoot.GetSourceFileName(), child.GetLineno(), child.GetCharno(), StaticGet_RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR(), []string{}))
				}
			}
		}
	}
	if !(valid) {
		return pkg_rhino.Node(0)
	} else if googImportNode != 0 {
		return googImportNode
	}
	var moduleScope *Scope = NewSyntacticScopeCreator_1(this_.compiler).CreateScope_2_Node_Scope(scriptRoot.GetFirstChild(), Scope_CreateGlobalScope(scriptRoot))
	_ = moduleScope
	var googVar *Var = j2gort.CastPtr[Var](moduleScope.GetVar("goog"))
	_ = googVar
	if (googVar != nil) && (googVar.GetNameNode() != 0) {
		this_.compiler.Report_1(JSError_Make_5(scriptRoot.GetSourceFileName(), googVar.GetNameNode().GetLineno(), googVar.GetNameNode().GetCharno(), StaticGet_RewriteGoogJsImports_CANNOT_HAVE_MODULE_VAR_NAMED_GOOG(), []string{}))
	}
	return pkg_rhino.Node(0)
}

// Java method: RewriteGoogJsImports.rewriteImports(com.google.javascript.rhino.Node scriptRoot)
func (this_ *RewriteGoogJsImports) priv_RewriteGoogJsImports_rewriteImports_4983(scriptRoot pkg_rhino.Node) {
	_ = scriptRoot
	if Es6RewriteModules_IsEs6ModuleRoot(scriptRoot) {
		var googImportNode pkg_rhino.Node = this_.priv_RewriteGoogJsImports_findGoogImportNode_4982(scriptRoot)
		_ = googImportNode
		NodeTraversal_Traverse(this_.compiler, scriptRoot, NewRewriteGoogJsImports_FindReexports((googImportNode != 0)))
		var module *Module = this_.moduleMap.GetModule_1_ModuleLoader_ModulePath(this_.compiler.GetInput(scriptRoot.GetInputId()).GetPath())
		_ = module
		j2gort.CheckNotNullPtr(module)
		if (googImportNode != 0) && (this_.mode == StaticGet_RewriteGoogJsImports_Mode_LINT_AND_REWRITE()) {
			var unused *RewriteGoogJsImports_ReferenceReplacer = NewRewriteGoogJsImports_ReferenceReplacer(this_, scriptRoot, googImportNode, module, (this_.googModule == nil))
			_ = unused
		}
	}
}

// Java method: RewriteGoogJsImports.findGoogJsScriptNode(com.google.javascript.rhino.Node root)
func (this_ *RewriteGoogJsImports) priv_RewriteGoogJsImports_findGoogJsScriptNode_4984(root pkg_rhino.Node) pkg_rhino.Node {
	_ = root
	var expectedGoogPath *ModuleLoader_ModulePath = nil
	_ = expectedGoogPath
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; script != 0; script = script.GetNext() {
			var provides []string = j2gort.ImmListToSliceNullable[string](this_.compiler.GetInput(script.GetInputId()).GetProvides())
			_ = provides
			if j2gort.SliceContains(provides, "goog") {
				expectedGoogPath = this_.compiler.GetInput(script.GetInputId()).GetPath().ResolveModuleAsPath("./goog.js")
				break
			}
		}
	}
	if expectedGoogPath != nil {
		var googScriptNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = googScriptNode
		{
			var script_2 pkg_rhino.Node = root.GetFirstChild()
			_ = script_2
			for ; script_2 != 0; script_2 = script_2.GetNext() {
				if this_.compiler.GetInput(script_2.GetInputId()).GetPath().EqualsIgnoreLeadingSlash(expectedGoogPath) {
					googScriptNode = script_2
				} else if javalang.StringEndsWith(j2gort.CheckNotNullString(script_2.GetSourceFileName()), "/goog.js") {
					this_.compiler.Report_1(JSError_Make_5(script_2.GetSourceFileName(), int32(-1), int32(-1), StaticGet_RewriteGoogJsImports_CANNOT_NAME_FILE_GOOG(), []string{}))
				}
			}
		}
		return googScriptNode
	}
	return pkg_rhino.Node(0)
}

// Java method: RewriteGoogJsImports.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteGoogJsImports) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.googModule = nil
	var googJsScriptNode pkg_rhino.Node = this_.priv_RewriteGoogJsImports_findGoogJsScriptNode_4984(root)
	_ = googJsScriptNode
	if googJsScriptNode == 0 {
		googJsScriptNode = this_.priv_RewriteGoogJsImports_findGoogJsScriptNode_4984(externs)
	}
	if this_.mode == StaticGet_RewriteGoogJsImports_Mode_LINT_AND_REWRITE() {
		if googJsScriptNode != 0 {
			var googJsPath *ModuleLoader_ModulePath = this_.compiler.GetInput(googJsScriptNode.GetInputId()).GetPath()
			_ = googJsPath
			this_.googModule = this_.moduleMap.GetModule_1_ModuleLoader_ModulePath(googJsPath)
			j2gort.CheckNotNullPtr(this_.googModule)
			var isFromGoog j2gort.Ext_com_google_common_base_Predicate = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(lp_3 any) bool {
				b := j2gort.CastPtr[Binding](lp_3)
				_ = b
				return (b.OriginatingExport().ModulePath() == googJsPath)
			})
			_ = isFromGoog
			if !(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](this_.googModule.BoundNames().Values()).Stream().AllMatch(isFromGoog)) {
				j2gort.CheckStateStr(false, "goog.js should never import anything")
			}
			if !(!(j2gort.MapContainsKeyStr(this_.googModule.Namespace(), "default"))) {
				j2gort.CheckStateStr(false, "goog.js should never have a default export.")
			}
			if !(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](this_.googModule.Namespace().Values()).Stream().AllMatch(isFromGoog)) {
				j2gort.CheckStateStr(false, "goog.js should never export from anything.")
			}
		}
	} else {
		if !(this_.mode == StaticGet_RewriteGoogJsImports_Mode_LINT_ONLY()) {
			j2gort.CheckState(false)
		}
	}
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; script != 0; script = script.GetNext() {
			this_.priv_RewriteGoogJsImports_rewriteImports_4983(script)
		}
	}
	this_.priv_RewriteGoogJsImports_changeModules_4981()
}

// Java enum: com.google.javascript.jscomp.RewriteGoogJsImports.Mode
type RewriteGoogJsImports_Mode uint8

const (
	enum_RewriteGoogJsImports_Mode_LINT_ONLY        RewriteGoogJsImports_Mode = 1
	enum_RewriteGoogJsImports_Mode_LINT_AND_REWRITE RewriteGoogJsImports_Mode = 2
)

var enumBaseTable_RewriteGoogJsImports_Mode = [...]j2gort.EnumBase{
	{Name_: "LINT_ONLY", Ordinal_: 0},
	{Name_: "LINT_AND_REWRITE", Ordinal_: 1},
}
var enumValues_RewriteGoogJsImports_Mode = []RewriteGoogJsImports_Mode{1, 2}

func StaticGet_RewriteGoogJsImports_Mode_LINT_ONLY() RewriteGoogJsImports_Mode {
	return enum_RewriteGoogJsImports_Mode_LINT_ONLY
}
func StaticGet_RewriteGoogJsImports_Mode_LINT_AND_REWRITE() RewriteGoogJsImports_Mode {
	return enum_RewriteGoogJsImports_Mode_LINT_AND_REWRITE
}

func RewriteGoogJsImports_Mode_Values() []RewriteGoogJsImports_Mode {
	return j2gort.CloneSlice(enumValues_RewriteGoogJsImports_Mode)
}

func RewriteGoogJsImports_Mode_ValueOf(name string) RewriteGoogJsImports_Mode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_RewriteGoogJsImports_Mode {
		if enumBaseTable_RewriteGoogJsImports_Mode[i_].Name_ == name {
			return RewriteGoogJsImports_Mode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.RewriteGoogJsImports.Mode." + name))
}

func (this_ RewriteGoogJsImports_Mode) IsJavaNil_() bool                              { return this_ == 0 }
func (RewriteGoogJsImports_Mode) HasIdentityEquals_()                                 {}
func (RewriteGoogJsImports_Mode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (RewriteGoogJsImports_Mode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ RewriteGoogJsImports_Mode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (RewriteGoogJsImports_Mode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteGoogJsImports$Mode", "Mode"
}
func (this_ RewriteGoogJsImports_Mode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_RewriteGoogJsImports_Mode[this_-1]
}
func (this_ RewriteGoogJsImports_Mode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ RewriteGoogJsImports_Mode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ RewriteGoogJsImports_Mode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RewriteGoogJsImports_Mode[this_-1].Name_
}
func (this_ RewriteGoogJsImports_Mode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_RewriteGoogJsImports_Mode[this_-1].Name_
}
func (this_ RewriteGoogJsImports_Mode) Equals(other_ any) bool {
	o_, ok_ := other_.(RewriteGoogJsImports_Mode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ RewriteGoogJsImports_Mode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_RewriteGoogJsImports_Mode[this_-1].HashCode()
}
func (this_ RewriteGoogJsImports_Mode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RewriteGoogJsImports_Mode[this_-1].CompareTo(other_)
}
func (this_ RewriteGoogJsImports_Mode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_RewriteGoogJsImports_Mode[this_-1].Name_)
}
func init() {
	vals_ := []any{RewriteGoogJsImports_Mode(1), RewriteGoogJsImports_Mode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.RewriteGoogJsImports.Mode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.RewriteGoogJsImports$Mode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteGoogJsImports.Mode", "com.google.javascript.jscomp.RewriteGoogJsImports$Mode", func(v_ any) bool {
		return j2gort.InstanceOf[RewriteGoogJsImports_Mode](v_)
	})
}

// Java class: com.google.javascript.jscomp.RewriteGoogJsImports.ReferenceReplacer
type RewriteGoogJsImports_ReferenceReplacer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                 *RewriteGoogJsImports
	googImportNode         pkg_rhino.Node `gson:"googImportNode"`
	hasBadExport           bool           `gson:"hasBadExport"`
	globalizeAllReferences bool           `gson:"globalizeAllReferences"`
}

func (this_ *RewriteGoogJsImports_ReferenceReplacer) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteGoogJsImports_ReferenceReplacer) IsJavaNil_() bool { return this_ == nil }
func (*RewriteGoogJsImports_ReferenceReplacer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteGoogJsImports$ReferenceReplacer", "ReferenceReplacer"
}
func (this_ *RewriteGoogJsImports_ReferenceReplacer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RewriteGoogJsImports_ReferenceReplacer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteGoogJsImports.ReferenceReplacer", "com.google.javascript.jscomp.RewriteGoogJsImports$ReferenceReplacer", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteGoogJsImports_ReferenceReplacer](v_)
	})
}
func (this_ *RewriteGoogJsImports_ReferenceReplacer) SetOuter_(o *RewriteGoogJsImports) {
	this_.outer_ = o
}
func (this_ *RewriteGoogJsImports_ReferenceReplacer) GetOuter_() *RewriteGoogJsImports {
	return this_.outer_
}
func (*RewriteGoogJsImports_ReferenceReplacer) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RewriteGoogJsImports_ReferenceReplacer) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RewriteGoogJsImports_ReferenceReplacer) initFields_RewriteGoogJsImports_ReferenceReplacer() {
	_ = this_
}

// Java constructor: ReferenceReplacer(com.google.javascript.rhino.Node script, com.google.javascript.rhino.Node googImportNode, com.google.javascript.jscomp.modules.Module module, boolean globalizeAllReferences)
func NewRewriteGoogJsImports_ReferenceReplacer(outer_ *RewriteGoogJsImports, script pkg_rhino.Node, googImportNode pkg_rhino.Node, module *Module, globalizeAllReferences bool) *RewriteGoogJsImports_ReferenceReplacer {
	this_ := &RewriteGoogJsImports_ReferenceReplacer{}
	this_.outer_ = outer_
	InitRewriteGoogJsImports_ReferenceReplacer(this_, script, googImportNode, module, globalizeAllReferences)
	return this_
}

func InitRewriteGoogJsImports_ReferenceReplacer(this_ *RewriteGoogJsImports_ReferenceReplacer, script pkg_rhino.Node, googImportNode pkg_rhino.Node, module *Module, globalizeAllReferences bool) {
	_ = script
	_ = googImportNode
	_ = module
	_ = globalizeAllReferences
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RewriteGoogJsImports_ReferenceReplacer()
	this_.googImportNode = googImportNode
	this_.globalizeAllReferences = globalizeAllReferences
	NodeTraversal_Traverse(this_.outer_.compiler, script, ToIface_NodeTraversal_Callback(this_))
	if googImportNode.GetSecondChild().IsImportStar() {
		var newBindings j2gort.OrderedMap[string, *Binding] = j2gort.OrderedMapFromJavaMap[string, *Binding](j2gort.AsIface[j2gort.Ext_java_util_Map](module.BoundNames()))
		_ = newBindings
		newBindings.Remove("goog")
		if this_.hasBadExport {
			googImportNode.GetSecondChild().SetOriginalName("goog")
			googImportNode.GetSecondChild().SetString("$goog")
			newBindings.Put("$goog", j2gort.CastPtr[Binding](j2gort.MapGetStr(module.BoundNames(), "goog")))
		} else {
			googImportNode.GetSecondChild().ReplaceWith(pkg_rhino.IR_Empty())
		}
		var newModule *Module = module.ToBuilder().BoundNames(newBindings.ToImmutableMap()).Build()
		_ = newModule
		this_.outer_.moduleReplacements.Put(module, newModule)
		this_.outer_.compiler.ReportChangeToEnclosingScope(googImportNode)
	}
}

// Java method: ReferenceReplacer.maybeRewriteBadGoogJsImportRef(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node parent)
func (this_ *RewriteGoogJsImports_ReferenceReplacer) priv_ReferenceReplacer_maybeRewriteBadGoogJsImportRef_4985(t *NodeTraversal, nameNode pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = nameNode
	_ = parent
	if !(parent.IsGetProp()) || !(j2gort.CheckNotNullString(nameNode.GetString()) == "goog") {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(nameNode.GetString()))
	_ = var_
	if ((var_ == nil) || (var_.GetNameNode() == 0)) || (var_.GetNameNode().GetParent() != this_.googImportNode) {
		return
	}
	if this_.globalizeAllReferences || j2gort.MapContainsKeyStr(this_.outer_.googModule.Namespace(), parent.GetString()) {
		return
	}
	this_.hasBadExport = true
	nameNode.SetOriginalName("goog")
	nameNode.SetString("$goog")
	t.ReportCodeChange_0()
}

// Java method: ReferenceReplacer.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteGoogJsImports_ReferenceReplacer) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_ReferenceReplacer_maybeRewriteBadGoogJsImportRef_4985(t, n, parent)
	default:
	}
}

// Java class: com.google.javascript.jscomp.RewriteGoogJsImports.FindReexports
type RewriteGoogJsImports_FindReexports struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	hasGoogImport bool `gson:"hasGoogImport"`
}

func (this_ *RewriteGoogJsImports_FindReexports) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteGoogJsImports_FindReexports) IsJavaNil_() bool { return this_ == nil }
func (*RewriteGoogJsImports_FindReexports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteGoogJsImports$FindReexports", "FindReexports"
}
func (this_ *RewriteGoogJsImports_FindReexports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*RewriteGoogJsImports_FindReexports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteGoogJsImports.FindReexports", "com.google.javascript.jscomp.RewriteGoogJsImports$FindReexports", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteGoogJsImports_FindReexports](v_)
	})
}
func (*RewriteGoogJsImports_FindReexports) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*RewriteGoogJsImports_FindReexports) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RewriteGoogJsImports_FindReexports) initFields_RewriteGoogJsImports_FindReexports() {
	_ = this_
}

// Java constructor: FindReexports(boolean hasGoogImport)
func NewRewriteGoogJsImports_FindReexports(hasGoogImport bool) *RewriteGoogJsImports_FindReexports {
	this_ := &RewriteGoogJsImports_FindReexports{}
	InitRewriteGoogJsImports_FindReexports(this_, hasGoogImport)
	return this_
}

func InitRewriteGoogJsImports_FindReexports(this_ *RewriteGoogJsImports_FindReexports, hasGoogImport bool) {
	_ = hasGoogImport
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_RewriteGoogJsImports_FindReexports()
	this_.hasGoogImport = hasGoogImport
}

// Java method: FindReexports.checkIfForwardingExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *RewriteGoogJsImports_FindReexports) priv_FindReexports_checkIfForwardingExport_4986(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	if export.HasTwoChildren() && javalang.StringEndsWith(j2gort.CheckNotNullString(export.GetLastChild().GetString()), "/goog.js") {
		t.Report_3(export, StaticGet_RewriteGoogJsImports_GOOG_JS_REEXPORTED(), []string{})
	}
}

// Java method: FindReexports.checkIfNameFowardedExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node parent)
func (this_ *RewriteGoogJsImports_FindReexports) priv_FindReexports_checkIfNameFowardedExport_4987(t *NodeTraversal, nameNode pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = nameNode
	_ = parent
	if this_.hasGoogImport && (j2gort.CheckNotNullString(nameNode.GetString()) == "goog") {
		if (parent.IsExportSpec() && (parent.GetFirstChild() == nameNode)) || parent.IsExport() {
			t.Report_3(nameNode, StaticGet_RewriteGoogJsImports_GOOG_JS_REEXPORTED(), []string{})
		}
	}
}

// Java method: FindReexports.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteGoogJsImports_FindReexports) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */, int32(123) /* Token.MODULE_BODY */, int32(121) /* Token.EXPORT_SPECS */, int32(122) /* Token.EXPORT_SPEC */ :
		return true
	case int32(120) /* Token.EXPORT */ :
		this_.priv_FindReexports_checkIfForwardingExport_4986(t, n)
		return true
	case int32(32) /* Token.NAME */ :
		this_.priv_FindReexports_checkIfNameFowardedExport_4987(t, n, parent)
		return false
	default:
		return false
	}
	panic("unreachable switch expression")
}
func (*RewriteGoogJsImports_ReferenceReplacer) AlwaysShouldTraverse_()      {}
func (*RewriteGoogJsImports_ReferenceReplacer) CanReleaseSyntacticScopes_() {}
func (*RewriteGoogJsImports_FindReexports) NoOpVisit_()                     {}
func (*RewriteGoogJsImports_FindReexports) CanReleaseSyntacticScopes_()     {}
