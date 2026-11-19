/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ModuleIdentifier.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.ModuleIdentifier
type ModuleIdentifier struct {
	name             string `gson:"name"`
	closureNamespace string `gson:"closureNamespace"`
	moduleName       string `gson:"moduleName"`
}

func (this_ *ModuleIdentifier) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleIdentifier) IsJavaNil_() bool { return this_ == nil }
func (*ModuleIdentifier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ModuleIdentifier", "ModuleIdentifier"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ModuleIdentifier", "com.google.javascript.jscomp.ModuleIdentifier", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleIdentifier](v_)
	})
}
func (*ModuleIdentifier) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ModuleIdentifier) initFields_ModuleIdentifier() {
	_ = this_
}

// Java constructor: ModuleIdentifier(java.lang.String name, java.lang.String closureNamespace, java.lang.String moduleName)
func NewModuleIdentifier(name string, closureNamespace string, moduleName string) *ModuleIdentifier {
	this_ := &ModuleIdentifier{}
	InitModuleIdentifier(this_, name, closureNamespace, moduleName)
	return this_
}

func InitModuleIdentifier(this_ *ModuleIdentifier, name string, closureNamespace string, moduleName string) {
	_ = name
	_ = closureNamespace
	_ = moduleName
	this_.initFields_ModuleIdentifier()
	j2gort.CheckNotNullString(name)
	j2gort.CheckNotNullString(closureNamespace)
	j2gort.CheckNotNullString(moduleName)
	this_.name = name
	this_.closureNamespace = closureNamespace
	this_.moduleName = moduleName
}

// Java method: ModuleIdentifier.getName()
func (this_ *ModuleIdentifier) GetName() string {
	return this_.Name()
}

// Java method: ModuleIdentifier.getClosureNamespace()
func (this_ *ModuleIdentifier) GetClosureNamespace() string {
	return this_.ClosureNamespace()
}

// Java method: ModuleIdentifier.getModuleName()
func (this_ *ModuleIdentifier) GetModuleName() string {
	return this_.ModuleName()
}

// Java method: ModuleIdentifier.toString()
func (this_ *ModuleIdentifier) String() string {
	{
		var iifeCond_1 bool
		{

			sr_1 := this_.ClosureNamespace()
			se_2 := this_.ModuleName()
			iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
		}
		if iifeCond_1 {
			return this_.ClosureNamespace()
		}
	}
	return (j2gort.StrValString(this_.ModuleName()) + ":" + j2gort.StrValString(this_.ClosureNamespace()))
}

// Java static method: ModuleIdentifier.forClosure(java.lang.String name)
func ModuleIdentifier_ForClosure(name string) *ModuleIdentifier {
	_ = name
	var normalizedName string = name
	_ = normalizedName
	if javalang.StringStartsWith(j2gort.CheckNotNullString(normalizedName), "goog:") {
		normalizedName = j2gort.JStrSubFrom(normalizedName, j2gort.JStrLen("goog:"))
	}
	var namespace string = normalizedName
	_ = namespace
	var moduleName string = normalizedName
	_ = moduleName
	var splitPoint int32 = j2gort.JStrIndexOfCp(normalizedName, int32(uint16(58)))
	_ = splitPoint
	if splitPoint != int32(-1) {
		moduleName = j2gort.JStrSub(normalizedName, int32(0), splitPoint)
		namespace = j2gort.JStrSubFrom(normalizedName, min((splitPoint+int32(1)), (j2gort.JStrLen(normalizedName)-int32(1))))
	}
	return NewModuleIdentifier(normalizedName, namespace, moduleName)
}

// Java static method: ModuleIdentifier.forFile(java.lang.String filepath)
func ModuleIdentifier_ForFile(filepath string) *ModuleIdentifier {
	_ = filepath
	var normalizedName string = ModuleNames_FileToModuleName(filepath)
	_ = normalizedName
	return NewModuleIdentifier(filepath, normalizedName, normalizedName)
}

// Java static method: ModuleIdentifier.forFlagValue(java.lang.String flagValue)
func ModuleIdentifier_ForFlagValue(flagValue string) *ModuleIdentifier {
	_ = flagValue
	if javalang.StringStartsWith(j2gort.CheckNotNullString(flagValue), "goog:") {
		return ModuleIdentifier_ForClosure(flagValue)
	} else {
		return ModuleIdentifier_ForFile(flagValue)
	}
}

func (this_ *ModuleIdentifier) Name() string             { return this_.name }
func (this_ *ModuleIdentifier) ClosureNamespace() string { return this_.closureNamespace }
func (this_ *ModuleIdentifier) ModuleName() string       { return this_.moduleName }
func (this_ *ModuleIdentifier) Equals(other_ any) bool {
	o_, ok_ := other_.(*ModuleIdentifier)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.name != o_.name {
		return false
	}
	if this_.closureNamespace != o_.closureNamespace {
		return false
	}
	if this_.moduleName != o_.moduleName {
		return false
	}
	return true
}
func (this_ *ModuleIdentifier) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_3 := this_.name
		h_ = 31 * h_
		if hs_3 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_3)
		}
	}
	{
		hs_4 := this_.closureNamespace
		h_ = 31 * h_
		if hs_4 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_4)
		}
	}
	{
		hs_5 := this_.moduleName
		h_ = 31 * h_
		if hs_5 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_5)
		}
	}
	return h_
}
