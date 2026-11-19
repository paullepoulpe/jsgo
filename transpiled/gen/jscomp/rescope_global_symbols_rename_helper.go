/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsRenameHelper.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper
type RescopeGlobalSymbolsRenameHelper struct {
	compiler                           AbstractCompiler                                  `gson:"compiler"`
	globalSymbolNamespace              string                                            `gson:"globalSymbolNamespace"`
	unwrappedLocalNamesInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"unwrappedLocalNamesInDefiningChunk"`
	idHash_                            atomic.Int32
}

func (this_ *RescopeGlobalSymbolsRenameHelper) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbolsRenameHelper) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbolsRenameHelper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper", "RescopeGlobalSymbolsRenameHelper"
}
func (this_ *RescopeGlobalSymbolsRenameHelper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RescopeGlobalSymbolsRenameHelper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper", "com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbolsRenameHelper](v_)
	})
}

func StaticGet_RescopeGlobalSymbolsRenameHelper_DISAMBIGUATION_SUFFIX() string       { return "$" }
func StaticGet_RescopeGlobalSymbolsRenameHelper_UNWRAPPED_LOCAL_NAME_PREFIX() string { return "_" }
func (this_ *RescopeGlobalSymbolsRenameHelper) initFields_RescopeGlobalSymbolsRenameHelper() {
	_ = this_
	this_.globalSymbolNamespace = j2gort.NullString
}

// Java constructor: RescopeGlobalSymbolsRenameHelper(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String globalSymbolNamespace, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> unwrappedLocalNamesInDefiningChunk)
func NewRescopeGlobalSymbolsRenameHelper(compiler AbstractCompiler, globalSymbolNamespace string, unwrappedLocalNamesInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableMap) *RescopeGlobalSymbolsRenameHelper {
	this_ := &RescopeGlobalSymbolsRenameHelper{}
	InitRescopeGlobalSymbolsRenameHelper(this_, compiler, globalSymbolNamespace, unwrappedLocalNamesInDefiningChunk)
	return this_
}

func InitRescopeGlobalSymbolsRenameHelper(this_ *RescopeGlobalSymbolsRenameHelper, compiler AbstractCompiler, globalSymbolNamespace string, unwrappedLocalNamesInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = compiler
	_ = globalSymbolNamespace
	_ = unwrappedLocalNamesInDefiningChunk
	this_.initFields_RescopeGlobalSymbolsRenameHelper()
	this_.compiler = compiler
	this_.globalSymbolNamespace = globalSymbolNamespace
	this_.unwrappedLocalNamesInDefiningChunk = unwrappedLocalNamesInDefiningChunk
}

// Java static method: RescopeGlobalSymbolsRenameHelper.conflictsWithGlobalSymbolNamespace(java.lang.String name, java.lang.String globalSymbolNamespace)
func RescopeGlobalSymbolsRenameHelper_ConflictsWithGlobalSymbolNamespace(name string, globalSymbolNamespace string) bool {
	_ = name
	_ = globalSymbolNamespace
	return ((j2gort.CheckNotNullString(name) == globalSymbolNamespace) || javalang.StringStartsWith(j2gort.CheckNotNullString(name), j2gort.StrValString(globalSymbolNamespace)+"$"))
}

// Java static method: RescopeGlobalSymbolsRenameHelper.getDisambiguatedLocalName(java.lang.String name)
func RescopeGlobalSymbolsRenameHelper_GetDisambiguatedLocalName(name string) string {
	_ = name
	return (j2gort.StrValString(name) + "$")
}

// Java static method: RescopeGlobalSymbolsRenameHelper.collectUnwrappedLocalNamesInDefiningChunk(com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess optimizeLocalAccess, com.google.javascript.jscomp.RescopeGlobalSymbolsInfo symbolInfo, java.util.Set<java.lang.String> namesStartingWithUnderscore)
func RescopeGlobalSymbolsRenameHelper_CollectUnwrappedLocalNamesInDefiningChunk(optimizeLocalAccess CompilerOptions_OptimizeLocalAccess, symbolInfo *RescopeGlobalSymbolsInfo, namesStartingWithUnderscore j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = optimizeLocalAccess
	_ = symbolInfo
	_ = namesStartingWithUnderscore
	if !(optimizeLocalAccess.ShouldKeepLocalAliasForWrappedReassignables()) {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	}
	var usedNames j2gort.LocalOrderedSet[string] = j2gort.LocalOrderedSetFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](namesStartingWithUnderscore))
	_ = usedNames
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = builder
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(symbolInfo.WrappedCrossChunkNames()))
		for cur_1.HasNext() {
			name := cur_1.NextStr()
			_ = name
			if (!(j2gort.SetContainsStr(symbolInfo.CrossChunkNamesWithWriteFromOtherChunk(), name)) && !(j2gort.SetContainsStr(symbolInfo.GlobalNamesWithUnsafeWriteForLocalAliasInDefiningChunk(), name))) && j2gort.SetContainsStr(symbolInfo.GlobalNamesWithReadInDefiningChunk(), name) {
				var unwrappedName string = ("_" + j2gort.StrValString(name))
				_ = unwrappedName
				for (usedNames.Contains(unwrappedName) || j2gort.SetContainsStr(symbolInfo.ExternNames(), unwrappedName)) || func() bool {
					sr_2 := unwrappedName
					se_3 := symbolInfo.GlobalSymbolNamespace()
					return j2gort.CheckNotNullString(sr_2) == se_3
				}() {
					unwrappedName = ("_" + j2gort.StrValString(unwrappedName))
				}
				usedNames.Add(unwrappedName)
				builder.Put_2(j2gort.BoxString(name), j2gort.BoxString(unwrappedName))
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java method: RescopeGlobalSymbolsRenameHelper.hasUnwrappedLocalNameInDefiningChunk(java.lang.String name)
func (this_ *RescopeGlobalSymbolsRenameHelper) HasUnwrappedLocalNameInDefiningChunk(name string) bool {
	_ = name
	return j2gort.MapContainsKeyStr(this_.unwrappedLocalNamesInDefiningChunk, name)
}

// Java method: RescopeGlobalSymbolsRenameHelper.getLocalNameInDefiningChunk(java.lang.String name, boolean isGlobalDeclaration)
func (this_ *RescopeGlobalSymbolsRenameHelper) GetLocalNameInDefiningChunk(name string, isGlobalDeclaration bool) string {
	_ = name
	_ = isGlobalDeclaration
	if isGlobalDeclaration {
		return j2gort.CastString(j2gort.MapGetOrDefaultStr(this_.unwrappedLocalNamesInDefiningChunk, name, j2gort.BoxString(name)))
	} else {
		return name
	}
}

// Java method: RescopeGlobalSymbolsRenameHelper.maybeDisambiguateLocalVar(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var var)
func (this_ *RescopeGlobalSymbolsRenameHelper) MaybeDisambiguateLocalVar(n pkg_rhino.Node, var_ *Var) {
	_ = n
	_ = var_
	var name string = n.GetString()
	_ = name
	if !(var_.IsGlobal()) && RescopeGlobalSymbolsRenameHelper_ConflictsWithGlobalSymbolNamespace(name, this_.globalSymbolNamespace) {
		n.SetString(RescopeGlobalSymbolsRenameHelper_GetDisambiguatedLocalName(name))
		this_.compiler.ReportChangeToEnclosingScope(n)
	}
}

// Java method: RescopeGlobalSymbolsRenameHelper.maybeRenameGlobalInDefiningChunk(com.google.javascript.rhino.Node n)
func (this_ *RescopeGlobalSymbolsRenameHelper) MaybeRenameGlobalInDefiningChunk(n pkg_rhino.Node) string {
	_ = n
	var name string = n.GetString()
	_ = name
	var localName string = j2gort.CastString(j2gort.MapGetOrDefaultStr(this_.unwrappedLocalNamesInDefiningChunk, name, j2gort.BoxString(name)))
	_ = localName
	if !(j2gort.CheckNotNullString(localName) == name) {
		n.SetString(localName)
		this_.compiler.ReportChangeToEnclosingScope(n)
	}
	return localName
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper.HelperBuilder
type RescopeGlobalSymbolsRenameHelper_HelperBuilder struct {
	compiler                    AbstractCompiler         `gson:"compiler"`
	globalSymbolNamespace       string                   `gson:"globalSymbolNamespace"`
	namesStartingWithUnderscore j2gort.Ext_java_util_Set `gson:"namesStartingWithUnderscore"`
	idHash_                     atomic.Int32
	optimizeLocalAccess         CompilerOptions_OptimizeLocalAccess `gson:"optimizeLocalAccess"`
}

func (this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbolsRenameHelper_HelperBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper$HelperBuilder", "HelperBuilder"
}
func (this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RescopeGlobalSymbolsRenameHelper_HelperBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper.HelperBuilder", "com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper$HelperBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbolsRenameHelper_HelperBuilder](v_)
	})
}

func (this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder) initFields_RescopeGlobalSymbolsRenameHelper_HelperBuilder() {
	_ = this_
	this_.globalSymbolNamespace = j2gort.NullString
	this_.namesStartingWithUnderscore = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: HelperBuilder(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String globalSymbolNamespace, com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess optimizeLocalAccess)
func NewRescopeGlobalSymbolsRenameHelper_HelperBuilder(compiler AbstractCompiler, globalSymbolNamespace string, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) *RescopeGlobalSymbolsRenameHelper_HelperBuilder {
	this_ := &RescopeGlobalSymbolsRenameHelper_HelperBuilder{}
	InitRescopeGlobalSymbolsRenameHelper_HelperBuilder(this_, compiler, globalSymbolNamespace, optimizeLocalAccess)
	return this_
}

func InitRescopeGlobalSymbolsRenameHelper_HelperBuilder(this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder, compiler AbstractCompiler, globalSymbolNamespace string, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) {
	_ = compiler
	_ = globalSymbolNamespace
	_ = optimizeLocalAccess
	this_.initFields_RescopeGlobalSymbolsRenameHelper_HelperBuilder()
	this_.compiler = compiler
	this_.globalSymbolNamespace = globalSymbolNamespace
	this_.optimizeLocalAccess = optimizeLocalAccess
}

// Java method: HelperBuilder.recordExistingName(java.lang.String name)
func (this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder) RecordExistingName(name string) {
	_ = name
	if !(this_.optimizeLocalAccess.ShouldKeepLocalAliasForWrappedReassignables()) {
		return
	}
	if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "_") {
		j2gort.SetAddStr(this_.namesStartingWithUnderscore, name)
	}
	if RescopeGlobalSymbolsRenameHelper_ConflictsWithGlobalSymbolNamespace(name, this_.globalSymbolNamespace) {
		var disambiguatedName string = RescopeGlobalSymbolsRenameHelper_GetDisambiguatedLocalName(name)
		_ = disambiguatedName
		if javalang.StringStartsWith(j2gort.CheckNotNullString(disambiguatedName), "_") {
			j2gort.SetAddStr(this_.namesStartingWithUnderscore, disambiguatedName)
		}
	}
}

// Java method: HelperBuilder.build(com.google.javascript.jscomp.RescopeGlobalSymbolsInfo symbolInfo)
func (this_ *RescopeGlobalSymbolsRenameHelper_HelperBuilder) Build(symbolInfo *RescopeGlobalSymbolsInfo) *RescopeGlobalSymbolsRenameHelper {
	_ = symbolInfo
	var unwrappedLocalNamesInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableMap = RescopeGlobalSymbolsRenameHelper_CollectUnwrappedLocalNamesInDefiningChunk(this_.optimizeLocalAccess, symbolInfo, this_.namesStartingWithUnderscore)
	_ = unwrappedLocalNamesInDefiningChunk
	return NewRescopeGlobalSymbolsRenameHelper(this_.compiler, this_.globalSymbolNamespace, unwrappedLocalNamesInDefiningChunk)
}
