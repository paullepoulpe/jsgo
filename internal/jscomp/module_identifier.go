// Port of com/google/javascript/jscomp/ModuleIdentifier.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ModuleIdentifier holds basic information on an entry point module.
//
// Closure entry points are namespace names, while ES and CommonJS entry
// points are file paths which are normalized to a namespace name.
//
// This class allows error messages to be based on the user-provided name
// rather than the normalized name.
//
// name: the user-provided name. closureNamespace: the Closure namespace
// name. moduleName: the module name.
type ModuleIdentifier struct {
	name             string
	closureNamespace string
	moduleName       string
}

// Java: ModuleIdentifier(String, String, String)
func NewModuleIdentifier(name string, closureNamespace string, moduleName string) *ModuleIdentifier {
	return &ModuleIdentifier{name: name, closureNamespace: closureNamespace, moduleName: moduleName}
}

// Java: ModuleIdentifier.name()
func (m *ModuleIdentifier) Name() string { return m.name }

// Java: ModuleIdentifier.closureNamespace()
func (m *ModuleIdentifier) ClosureNamespace() string { return m.closureNamespace }

// Java: ModuleIdentifier.moduleName()
func (m *ModuleIdentifier) ModuleName() string { return m.moduleName }

// Java: ModuleIdentifier.getName()
func (m *ModuleIdentifier) GetName() string { return m.Name() }

// Java: ModuleIdentifier.getClosureNamespace()
func (m *ModuleIdentifier) GetClosureNamespace() string { return m.ClosureNamespace() }

// Java: ModuleIdentifier.getModuleName()
func (m *ModuleIdentifier) GetModuleName() string { return m.ModuleName() }

// Java: ModuleIdentifier.toString()
func (m *ModuleIdentifier) String() string {
	if m.ClosureNamespace() == m.ModuleName() {
		return m.ClosureNamespace()
	}
	return m.ModuleName() + ":" + m.ClosureNamespace()
}

// Java: ModuleIdentifier.equals(Object) (record)
func (m *ModuleIdentifier) Equals(o any) bool {
	that, ok := o.(*ModuleIdentifier)
	return ok && that != nil && m.name == that.name && m.closureNamespace == that.closureNamespace &&
		m.moduleName == that.moduleName
}

// Java: ModuleIdentifier.hashCode() (record)
func (m *ModuleIdentifier) HashCode() int32 {
	var h int32
	h = 31*h + jutil.StringHashCode(m.name)
	h = 31*h + jutil.StringHashCode(m.closureNamespace)
	h = 31*h + jutil.StringHashCode(m.moduleName)
	return h
}

// ModuleIdentifier_ForClosure returns an identifier for a Closure
// namespace.
//
// name is the Closure namespace. It may be in one of the formats
// `name.space`, `goog:name.space` or `goog:moduleName:name.space`, where the
// latter specifies that the module and namespace names are different.
//
// Java: ModuleIdentifier.forClosure(String)
func ModuleIdentifier_ForClosure(name string) *ModuleIdentifier {
	normalizedName := name
	if strings.HasPrefix(normalizedName, "goog:") {
		normalizedName = normalizedName[len("goog:"):]
	}

	namespace := normalizedName
	moduleName := normalizedName
	splitPoint := jutil.StringIndexOfChar(normalizedName, ':')
	if splitPoint != -1 {
		moduleName = jutil.StringSubstring(normalizedName, 0, splitPoint)
		namespace = jutil.StringSubstringFrom(normalizedName, min(splitPoint+1, jutil.StringLength(normalizedName)-1))
	}

	return NewModuleIdentifier(normalizedName, namespace, moduleName)
}

// ModuleIdentifier_ForFile returns an identifier for an ES or CommonJS
// module.
//
// Java: ModuleIdentifier.forFile(String)
func ModuleIdentifier_ForFile(filepath string) *ModuleIdentifier {
	normalizedName := ModuleNames_FileToModuleName(filepath)
	return NewModuleIdentifier(filepath, normalizedName, normalizedName)
}

// ModuleIdentifier_ForFlagValue returns an identifier for an --entry_point
// flag value. If it is in one of the formats `goog:name.space` or
// `goog:moduleName:name.space`, it is interpreted as a Closure namespace.
// Otherwise, it is interpreted as the path to an ES or CommonJS module.
//
// Java: ModuleIdentifier.forFlagValue(String)
func ModuleIdentifier_ForFlagValue(flagValue string) *ModuleIdentifier {
	if strings.HasPrefix(flagValue, "goog:") {
		return ModuleIdentifier_ForClosure(flagValue)
	}
	return ModuleIdentifier_ForFile(flagValue)
}
