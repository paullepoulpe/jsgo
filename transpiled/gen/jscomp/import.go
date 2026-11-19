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
// Port of com/google/javascript/jscomp/modules/Import.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.modules.Import
type Import struct {
	moduleRequest string                   `gson:"moduleRequest"`
	importName    string                   `gson:"importName"`
	localName     string                   `gson:"localName"`
	modulePath    *ModuleLoader_ModulePath `gson:"modulePath"`
	importNode    pkg_rhino.Node           `gson:"importNode"`
	nameNode      pkg_rhino.Node           `gson:"nameNode"`
}

func (this_ *Import) BindSelf_(s any) {
	_ = s
}
func (this_ *Import) IsJavaNil_() bool { return this_ == nil }
func (*Import) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Import", "Import"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Import", "com.google.javascript.jscomp.modules.Import", func(v_ any) bool {
		return j2gort.InstanceOf[*Import](v_)
	})
}
func (*Import) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Import) initFields_Import() {
	_ = this_
}

// Java constructor: Import(java.lang.String moduleRequest, java.lang.String importName, java.lang.String localName, com.google.javascript.jscomp.deps.ModuleLoader.@org.jspecify.annotations.Nullable ModulePath modulePath, com.google.javascript.rhino.Node importNode, com.google.javascript.rhino.Node nameNode)
func NewImport(moduleRequest string, importName string, localName string, modulePath *ModuleLoader_ModulePath, importNode pkg_rhino.Node, nameNode pkg_rhino.Node) *Import {
	this_ := &Import{}
	InitImport(this_, moduleRequest, importName, localName, modulePath, importNode, nameNode)
	return this_
}

func InitImport(this_ *Import, moduleRequest string, importName string, localName string, modulePath *ModuleLoader_ModulePath, importNode pkg_rhino.Node, nameNode pkg_rhino.Node) {
	_ = moduleRequest
	_ = importName
	_ = localName
	_ = modulePath
	_ = importNode
	_ = nameNode
	this_.initFields_Import()
	j2gort.CheckNotNullString(moduleRequest)
	j2gort.CheckNotNullString(importName)
	j2gort.CheckNotNullString(localName)
	j2gort.CheckNotNullNode(importNode)
	j2gort.CheckNotNullNode(nameNode)
	this_.moduleRequest = moduleRequest
	this_.importName = importName
	this_.localName = localName
	this_.modulePath = modulePath
	this_.importNode = importNode
	this_.nameNode = nameNode
}

// Java static method: Import.builder()
func Import_S_Builder() Import_Builder {
	return NewAutoBuilder_Import_Builder()
}

func (this_ *Import) ModuleRequest() string                { return this_.moduleRequest }
func (this_ *Import) ImportName() string                   { return this_.importName }
func (this_ *Import) LocalName() string                    { return this_.localName }
func (this_ *Import) ModulePath() *ModuleLoader_ModulePath { return this_.modulePath }
func (this_ *Import) ImportNode() pkg_rhino.Node           { return this_.importNode }
func (this_ *Import) NameNode() pkg_rhino.Node             { return this_.nameNode }
func (this_ *Import) Equals(other_ any) bool {
	o_, ok_ := other_.(*Import)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.moduleRequest != o_.moduleRequest {
		return false
	}
	if this_.importName != o_.importName {
		return false
	}
	if this_.localName != o_.localName {
		return false
	}
	if this_.modulePath != o_.modulePath {
		return false
	}
	if this_.importNode != o_.importNode {
		return false
	}
	if this_.nameNode != o_.nameNode {
		return false
	}
	return true
}
func (this_ *Import) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_1 := this_.moduleRequest
		h_ = 31 * h_
		if hs_1 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_1)
		}
	}
	{
		hs_2 := this_.importName
		h_ = 31 * h_
		if hs_2 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_2)
		}
	}
	{
		hs_3 := this_.localName
		h_ = 31 * h_
		if hs_3 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_3)
		}
	}
	{
		hp_4 := this_.modulePath
		h_ = 31 * h_
		if hp_4 != nil {
			h_ += j2gort.IdentityHash(hp_4)
		}
	}
	h_ = 31*h_ + int32(this_.importNode)
	h_ = 31*h_ + int32(this_.nameNode)
	return h_
}
func (this_ *Import) String() string {
	_ = fmt.Sprint
	return "Import[" + "moduleRequest=" + j2gort.StrVal(j2gort.BoxString(this_.moduleRequest)) + ", " + "importName=" + j2gort.StrVal(j2gort.BoxString(this_.importName)) + ", " + "localName=" + j2gort.StrVal(j2gort.BoxString(this_.localName)) + ", " + "modulePath=" + j2gort.StrVal(j2gort.PtrToAny(this_.modulePath)) + ", " + "importNode=" + j2gort.StrVal(j2gort.BoxNode(this_.importNode)) + ", " + "nameNode=" + j2gort.StrVal(j2gort.BoxNode(this_.nameNode)) + "]"
}

// Java polymorphic class interface: com.google.javascript.jscomp.modules.Import.Builder
type Import_Builder interface {
	Implements_jscomp_Import_Builder()
	IsJavaNil_() bool
	GetImport_BuilderBase() *Import_BuilderBase
	ModuleRequest(string) Import_Builder
	ImportName(string) Import_Builder
	LocalName(string) Import_Builder
	ModulePath(*ModuleLoader_ModulePath) Import_Builder
	ImportNode(pkg_rhino.Node) Import_Builder
	NameNode(pkg_rhino.Node) Import_Builder
	Build() *Import
}

func ToIface_Import_Builder[T interface {
	Import_Builder
	comparable
}](v T) Import_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.modules.Import.Builder
type Import_BuilderBase struct {
	idHash_ atomic.Int32
}

func (this_ *Import_BuilderBase) BindSelf_(s any) {
	_ = s
}
func (this_ *Import_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*Import_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.Import$Builder", "Builder"
}
func (this_ *Import_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.Import.Builder", "com.google.javascript.jscomp.modules.Import$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[Import_Builder](v_)
	})
}
func (b *Import_BuilderBase) GetImport_BuilderBase() *Import_BuilderBase { return b }
func (*Import_BuilderBase) Implements_jscomp_Import_Builder()            {}

func (this_ *Import_BuilderBase) initFields_Import_Builder() {
	_ = this_
}

func InitImport_Builder(this_ *Import_BuilderBase) {
	this_.initFields_Import_Builder()
}

// Java method: Builder.moduleRequest(java.lang.String value)
func (this_ *Import_BuilderBase) ModuleRequest(value string) Import_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.importName(java.lang.String value)
func (this_ *Import_BuilderBase) ImportName(value string) Import_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.localName(java.lang.String value)
func (this_ *Import_BuilderBase) LocalName(value string) Import_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.modulePath(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath value)
func (this_ *Import_BuilderBase) ModulePath(value *ModuleLoader_ModulePath) Import_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.importNode(com.google.javascript.rhino.Node value)
func (this_ *Import_BuilderBase) ImportNode(value pkg_rhino.Node) Import_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.nameNode(com.google.javascript.rhino.Node value)
func (this_ *Import_BuilderBase) NameNode(value pkg_rhino.Node) Import_Builder {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *Import_BuilderBase) Build() *Import {
	panic("abstract or native method invoked")
}
