// Port of com/google/javascript/jscomp/modules/ModuleMap.java @ closure-compiler v20260415

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ModuleMap is a map containing information about all modules in the
// compilation.
//
// This is currently used for ES modules and other types of module are not
// processed in detail.
type ModuleMap struct {
	resolvedModules        *jutil.LinkedHashMap[string, *Module]
	resolvedClosureModules *jutil.LinkedHashMap[string, *Module]
}

// Java: ModuleMap(ImmutableMap<String, Module>, ImmutableMap<String, Module>)
func NewModuleMap(
	resolvedModules *jutil.LinkedHashMap[string, *Module],
	resolvedClosureModules *jutil.LinkedHashMap[string, *Module]) *ModuleMap {
	return &ModuleMap{resolvedModules: resolvedModules, resolvedClosureModules: resolvedClosureModules}
}

// GetModule may return nil.
//
// Java: ModuleMap.getModule(String)
func (m *ModuleMap) GetModule(moduleName string) *Module {
	return m.resolvedModules.Get(moduleName)
}

// GetModuleByPath may return nil.
//
// Java: ModuleMap.getModule(ModuleLoader.ModulePath)
func (m *ModuleMap) GetModuleByPath(path *ModuleLoader_ModulePath) *Module {
	return m.GetModule(path.ToModuleName())
}

// Java: ModuleMap.getModulesByPath()
func (m *ModuleMap) GetModulesByPath() *jutil.LinkedHashMap[string, *Module] {
	return m.resolvedModules
}

// Java: ModuleMap.getModulesByClosureNamespace()
func (m *ModuleMap) GetModulesByClosureNamespace() *jutil.LinkedHashMap[string, *Module] {
	return m.resolvedClosureModules
}

// GetClosureModule may return nil.
//
// Java: ModuleMap.getClosureModule(String)
func (m *ModuleMap) GetClosureModule(namespace string) *Module {
	return m.resolvedClosureModules.Get(namespace)
}

// Java: ModuleMap.emptyForTesting()
func ModuleMap_EmptyForTesting() *ModuleMap {
	return NewModuleMap(jutil.NewLinkedHashMap[string, *Module](), jutil.NewLinkedHashMap[string, *Module]())
}
