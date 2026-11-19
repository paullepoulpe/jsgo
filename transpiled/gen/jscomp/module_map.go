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
// Port of com/google/javascript/jscomp/modules/ModuleMap.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.modules.ModuleMap
type ModuleMap struct {
	resolvedModules j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"resolvedModules"`
	resolvedClosureModules j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"resolvedClosureModules"`
	idHash_ atomic.Int32
}

func (this_ *ModuleMap) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleMap) IsJavaNil_() bool { return this_ == nil }
func (*ModuleMap) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.modules.ModuleMap", "ModuleMap" }
func (this_ *ModuleMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ModuleMap", "com.google.javascript.jscomp.modules.ModuleMap", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleMap](v_)
	})
}

func (this_ *ModuleMap) initFields_ModuleMap() {
	_ = this_
}

// Java constructor: ModuleMap(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Module> resolvedModules, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Module> resolvedClosureModules)
func NewModuleMap(resolvedModules j2gort.Ext_com_google_common_collect_ImmutableMap, resolvedClosureModules j2gort.Ext_com_google_common_collect_ImmutableMap) *ModuleMap {
	this_ := &ModuleMap{}
	InitModuleMap(this_, resolvedModules, resolvedClosureModules)
	return this_
}

func InitModuleMap(this_ *ModuleMap, resolvedModules j2gort.Ext_com_google_common_collect_ImmutableMap, resolvedClosureModules j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = resolvedModules
	_ = resolvedClosureModules
	this_.initFields_ModuleMap()
	this_.resolvedModules = resolvedModules
	this_.resolvedClosureModules = resolvedClosureModules
}

// Java method: ModuleMap.getModule(java.lang.String moduleName)
func (this_ *ModuleMap) GetModule_1_string(moduleName string) *Module {
	_ = moduleName
	return j2gort.CastPtr[Module](j2gort.MapGetStr(this_.resolvedModules, moduleName))
}

// Java method: ModuleMap.getModule(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path)
func (this_ *ModuleMap) GetModule_1_ModuleLoader_ModulePath(path *ModuleLoader_ModulePath) *Module {
	_ = path
	return this_.GetModule_1_string(path.ToModuleName())
}

// Java method: ModuleMap.getModulesByPath()
func (this_ *ModuleMap) GetModulesByPath() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.resolvedModules
}

// Java method: ModuleMap.getModulesByClosureNamespace()
func (this_ *ModuleMap) GetModulesByClosureNamespace() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.resolvedClosureModules
}

// Java method: ModuleMap.getClosureModule(java.lang.String namespace)
func (this_ *ModuleMap) GetClosureModule(namespace string) *Module {
	_ = namespace
	return j2gort.CastPtr[Module](j2gort.MapGetStr(this_.resolvedClosureModules, namespace))
}

// Java static method: ModuleMap.emptyForTesting()
func ModuleMap_EmptyForTesting() *ModuleMap {
	return NewModuleMap(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0(), j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
}

