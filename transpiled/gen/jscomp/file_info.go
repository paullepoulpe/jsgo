/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ijs/FileInfo.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ijs.FileInfo
type FileInfo struct {
	providedNamespaces j2gort.Ext_java_util_Set                          `gson:"providedNamespaces"`
	requiredLocalNames j2gort.Ext_java_util_Set                          `gson:"requiredLocalNames"`
	declarations       j2gort.Ext_com_google_common_collect_ListMultimap `gson:"declarations"`
	idHash_            atomic.Int32
	isFromTypeScript   bool `gson:"isFromTypeScript"`
}

func (this_ *FileInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *FileInfo) IsJavaNil_() bool { return this_ == nil }
func (*FileInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.FileInfo", "FileInfo"
}
func (this_ *FileInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FileInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.FileInfo", "com.google.javascript.jscomp.ijs.FileInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*FileInfo](v_)
	})
}

func (this_ *FileInfo) initFields_FileInfo() {
	_ = this_
	this_.providedNamespaces = j2gort.NewExt_java_util_LinkedHashSet()
	this_.requiredLocalNames = j2gort.NewExt_java_util_LinkedHashSet()
	this_.declarations = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_MultimapBuilder_LinkedHashKeys_0().ArrayListValues_0().Build_0())
}

// Java constructor: FileInfo(java.lang.String sourceFileName)
func NewFileInfo(sourceFileName string) *FileInfo {
	this_ := &FileInfo{}
	InitFileInfo(this_, sourceFileName)
	return this_
}

func InitFileInfo(this_ *FileInfo, sourceFileName string) {
	_ = sourceFileName
	this_.initFields_FileInfo()
	this_.isFromTypeScript = javalang.StringEndsWith(j2gort.CheckNotNullString(sourceFileName), ".closure.js")
}

// Java method: FileInfo.isFromTypeScript()
func (this_ *FileInfo) IsFromTypeScript() bool {
	return this_.isFromTypeScript
}

// Java method: FileInfo.recordNameDeclaration(com.google.javascript.rhino.Node qualifiedNameNode)
func (this_ *FileInfo) RecordNameDeclaration(qualifiedNameNode pkg_rhino.Node) {
	_ = qualifiedNameNode
	this_.RecordDeclaration(PotentialDeclaration_FromName(qualifiedNameNode))
}

// Java method: FileInfo.recordMemberFieldDef(com.google.javascript.rhino.Node fieldNode)
func (this_ *FileInfo) RecordMemberFieldDef(fieldNode pkg_rhino.Node) {
	_ = fieldNode
	this_.RecordDeclaration(PotentialDeclaration_FromMemberFieldDef(fieldNode))
}

// Java method: FileInfo.recordMethod(com.google.javascript.rhino.Node functionNode)
func (this_ *FileInfo) RecordMethod(functionNode pkg_rhino.Node) {
	_ = functionNode
	this_.RecordDeclaration(PotentialDeclaration_FromMethod(functionNode))
}

// Java method: FileInfo.recordStringKeyDeclaration(com.google.javascript.rhino.Node stringKeyNode)
func (this_ *FileInfo) RecordStringKeyDeclaration(stringKeyNode pkg_rhino.Node) {
	_ = stringKeyNode
	this_.RecordDeclaration(PotentialDeclaration_FromStringKey(stringKeyNode))
}

// Java method: FileInfo.recordDefine(com.google.javascript.rhino.Node callNode)
func (this_ *FileInfo) RecordDefine(callNode pkg_rhino.Node) {
	_ = callNode
	this_.RecordDeclaration(PotentialDeclaration_FromDefine(callNode))
}

// Java method: FileInfo.recordAliasDeclaration(com.google.javascript.rhino.Node nameNode)
func (this_ *FileInfo) RecordAliasDeclaration(nameNode pkg_rhino.Node) {
	_ = nameNode
	this_.RecordDeclaration(PotentialDeclaration_FromAlias(nameNode))
}

// Java method: FileInfo.getDeclarations()
func (this_ *FileInfo) GetDeclarations() j2gort.Ext_com_google_common_collect_ListMultimap {
	return this_.declarations
}

// Java method: FileInfo.recordDeclaration(com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func (this_ *FileInfo) RecordDeclaration(decl PotentialDeclaration) {
	_ = decl
	j2gort.MultimapPutStr(this_.declarations, decl.GetFullyQualifiedName(), decl)
}

// Java method: FileInfo.recordImport(java.lang.String localName)
func (this_ *FileInfo) RecordImport(localName string) {
	_ = localName
	j2gort.SetAddStr(this_.requiredLocalNames, localName)
}

// Java method: FileInfo.isNameDeclared(java.lang.String fullyQualifiedName)
func (this_ *FileInfo) IsNameDeclared(fullyQualifiedName string) bool {
	_ = fullyQualifiedName
	return j2gort.MultimapContainsKeyStr(this_.declarations, fullyQualifiedName)
}

// Java static method: FileInfo.containsPrefix(java.lang.String fullyQualifiedName, java.lang.Iterable<java.lang.String> prefixNamespaces)
func FileInfo_ContainsPrefix(fullyQualifiedName string, prefixNamespaces j2gort.Ext_java_lang_Iterable) bool {
	_ = fullyQualifiedName
	_ = prefixNamespaces
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(prefixNamespaces))
		for cur_1.HasNext() {
			prefix := cur_1.NextStr()
			_ = prefix
			if (j2gort.CheckNotNullString(fullyQualifiedName) == prefix) || javalang.StringStartsWith(j2gort.CheckNotNullString(fullyQualifiedName), j2gort.StrValString(prefix)+".") {
				return true
			}
		}
	}
	return false
}

// Java method: FileInfo.isPrefixProvided(java.lang.String fullyQualifiedName)
func (this_ *FileInfo) IsPrefixProvided(fullyQualifiedName string) bool {
	_ = fullyQualifiedName
	return FileInfo_ContainsPrefix(fullyQualifiedName, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.providedNamespaces))
}

// Java method: FileInfo.isPrefixRequired(java.lang.String fullyQualifiedName)
func (this_ *FileInfo) IsPrefixRequired(fullyQualifiedName string) bool {
	_ = fullyQualifiedName
	return FileInfo_ContainsPrefix(fullyQualifiedName, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.requiredLocalNames))
}

// Java method: FileInfo.isStrictPrefixDeclared(java.lang.String fullyQualifiedName)
func (this_ *FileInfo) IsStrictPrefixDeclared(fullyQualifiedName string) bool {
	_ = fullyQualifiedName
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.declarations.KeySet())))
		for cur_2.HasNext() {
			prefix := cur_2.NextStr()
			_ = prefix
			if javalang.StringStartsWith(j2gort.CheckNotNullString(fullyQualifiedName), j2gort.StrValString(prefix)+".") {
				return true
			}
		}
	}
	return false
}

// Java method: FileInfo.markProvided(java.lang.String providedName)
func (this_ *FileInfo) MarkProvided(providedName string) {
	_ = providedName
	j2gort.CheckNotNullString(providedName)
	j2gort.SetAddStr(this_.providedNamespaces, providedName)
}
