// Port of com/google/javascript/jscomp/ijs/FileInfo.java @ closure-compiler v20261015

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FileInfo is the class to keep track of what has been seen so far in a
// given file.
type FileInfo struct {
	providedNamespaces *jutil.LinkedHashSet[string]
	requiredLocalNames *jutil.LinkedHashSet[string]
	// MultimapBuilder.linkedHashKeys().arrayListValues().build()
	declarations     *jutil.ListMultimap[string, PotentialDeclaration]
	isFromTypeScript bool
}

// Java: FileInfo(String)
func NewFileInfo(sourceFileName string) *FileInfo {
	f := &FileInfo{
		providedNamespaces: jutil.NewLinkedHashSet[string](),
		requiredLocalNames: jutil.NewLinkedHashSet[string](),
		declarations:       jutil.ArrayListMultimap_Create[string, PotentialDeclaration](),
	}
	f.isFromTypeScript = strings.HasSuffix(sourceFileName, ".closure.js")
	return f
}

// Java: FileInfo.isFromTypeScript()
func (f *FileInfo) IsFromTypeScript() bool {
	return f.isFromTypeScript
}

// Java: FileInfo.recordNameDeclaration(Node)
func (f *FileInfo) RecordNameDeclaration(qualifiedNameNode *rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromName(qualifiedNameNode))
}

// Java: FileInfo.recordMemberFieldDef(Node)
func (f *FileInfo) RecordMemberFieldDef(fieldNode *rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromMemberFieldDef(fieldNode))
}

// Java: FileInfo.recordMethod(Node)
func (f *FileInfo) RecordMethod(functionNode *rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromMethod(functionNode))
}

// Java: FileInfo.recordStringKeyDeclaration(Node)
func (f *FileInfo) RecordStringKeyDeclaration(stringKeyNode *rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromStringKey(stringKeyNode))
}

// Java: FileInfo.recordDefine(Node)
func (f *FileInfo) RecordDefine(callNode *rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromDefine(callNode))
}

// Java: FileInfo.recordAliasDeclaration(Node)
func (f *FileInfo) RecordAliasDeclaration(nameNode *rhino.Node) {
	f.RecordDeclaration(PotentialDeclaration_FromAlias(nameNode))
}

// Java: FileInfo.getDeclarations()
func (f *FileInfo) GetDeclarations() *jutil.ListMultimap[string, PotentialDeclaration] {
	return f.declarations
}

// Java: FileInfo.recordDeclaration(PotentialDeclaration)
func (f *FileInfo) RecordDeclaration(decl PotentialDeclaration) {
	f.declarations.Put(decl.GetFullyQualifiedName(), decl)
}

// Java: FileInfo.recordImport(String)
func (f *FileInfo) RecordImport(localName string) {
	f.requiredLocalNames.Add(localName)
}

// Java: FileInfo.isNameDeclared(String)
func (f *FileInfo) IsNameDeclared(fullyQualifiedName string) bool {
	return f.declarations.ContainsKey(fullyQualifiedName)
}

// Java: FileInfo.containsPrefix(String, Iterable<String>)
func FileInfo_ContainsPrefix(fullyQualifiedName string, prefixNamespaces *jutil.LinkedHashSet[string]) bool {
	for prefix := range prefixNamespaces.All() {
		if fullyQualifiedName == prefix || strings.HasPrefix(fullyQualifiedName, prefix+".") {
			return true
		}
	}
	return false
}

// Java: FileInfo.isPrefixProvided(String)
func (f *FileInfo) IsPrefixProvided(fullyQualifiedName string) bool {
	return FileInfo_ContainsPrefix(fullyQualifiedName, f.providedNamespaces)
}

// Java: FileInfo.isPrefixRequired(String)
func (f *FileInfo) IsPrefixRequired(fullyQualifiedName string) bool {
	return FileInfo_ContainsPrefix(fullyQualifiedName, f.requiredLocalNames)
}

// Java: FileInfo.isStrictPrefixDeclared(String)
func (f *FileInfo) IsStrictPrefixDeclared(fullyQualifiedName string) bool {
	for prefix := range f.declarations.KeySet() {
		if strings.HasPrefix(fullyQualifiedName, prefix+".") {
			return true
		}
	}
	return false
}

// Java: FileInfo.markProvided(String)
func (f *FileInfo) MarkProvided(providedName string) {
	// checkNotNull(providedName): a Go string cannot be null.
	f.providedNamespaces.Add(providedName)
}
