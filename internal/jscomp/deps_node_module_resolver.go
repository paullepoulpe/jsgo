// Port of com/google/javascript/jscomp/deps/NodeModuleResolver.java @ closure-compiler v20261021

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
	"cmp"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// NodeModuleResolver is the resolution algorithm for NodeJS. See
// https://nodejs.org/api/modules.html#modules_all_together
type NodeModuleResolver struct {
	moduleResolverCommon

	packageJsonMainEntries *jutil.LinkedHashMap[string, string]
	nodeModulesFolders     []string
}

var _ ModuleResolver = (*NodeModuleResolver)(nil)

var nodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH = []string{
	"",
	".js",
	".json",
	// .i.js is typed interfaces. See ConvertToTypedInterface.java
	".i.js",
	".js.i.js",
}

var nodeModuleResolver_FILES_TO_SEARCH = []string{
	ModuleLoader_MODULE_SLASH + "package.json",
	ModuleLoader_MODULE_SLASH + "index.js",
	ModuleLoader_MODULE_SLASH + "index.json",
}

// Java: NodeModuleResolver.buildNodeModulesFoldersRegistry(Iterable<String>, Iterable<String>)
func nodeModuleResolver_buildNodeModulesFoldersRegistry(
	modulePaths *jutil.LinkedHashSet[string], moduleRootPaths []string) []string {
	seen := make(map[string]struct{})
	var registry []string

	for modulePath := range modulePaths.All() {
		for _, moduleRootPath := range moduleRootPaths {
			if strings.HasPrefix(modulePath, moduleRootPath) {
				modulePath = modulePath[len(moduleRootPath):]
				break
			}
		}
		nodeModulesDirs := jutil.StringSplitLiteral(modulePath, "/node_modules/")
		parentPath := ""
		for i := 0; i < len(nodeModulesDirs)-1; i++ {
			parentPath += nodeModulesDirs[i] + "/"
			if _, ok := seen[parentPath]; !ok {
				seen[parentPath] = struct{}{}
				registry = append(registry, parentPath)
			}
			parentPath += "node_modules/"
		}
	}

	slices.SortFunc(registry, func(a, b string) int {
		if c := cmp.Compare(jutil.StringLength(b), jutil.StringLength(a)); c != 0 {
			return c
		}
		return jutil.StringCompareTo(a, b)
	})
	return registry
}

// NodeModuleResolver_Factory is the factory for NodeModuleResolver.
type NodeModuleResolver_Factory struct {
	packageJsonMainEntries *jutil.LinkedHashMap[string, string]
}

var _ ModuleLoader_ModuleResolverFactory = (*NodeModuleResolver_Factory)(nil)

// Java: NodeModuleResolver.Factory()
func NewNodeModuleResolver_Factory() *NodeModuleResolver_Factory {
	return NewNodeModuleResolver_FactoryWithEntries(nil)
}

// Java: NodeModuleResolver.Factory(Map<String, String>)
func NewNodeModuleResolver_FactoryWithEntries(packageJsonMainEntries *jutil.LinkedHashMap[string, string]) *NodeModuleResolver_Factory {
	return &NodeModuleResolver_Factory{packageJsonMainEntries: packageJsonMainEntries}
}

// Java: NodeModuleResolver.Factory.create(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func (f *NodeModuleResolver_Factory) Create(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	return NewNodeModuleResolver(
		modulePaths, moduleRootPaths, f.packageJsonMainEntries, errorHandler, pathEscaper)
}

// Java: NodeModuleResolver(ImmutableSet<String>, ImmutableList<String>, Map<String, String>, ErrorHandler, PathEscaper)
func NewNodeModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	packageJsonMainEntries *jutil.LinkedHashMap[string, string],
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) *NodeModuleResolver {
	r := &NodeModuleResolver{
		moduleResolverCommon: newModuleResolverCommon(modulePaths, moduleRootPaths, errorHandler, pathEscaper),
		nodeModulesFolders:   nodeModuleResolver_buildNodeModulesFoldersRegistry(modulePaths, moduleRootPaths),
	}

	if packageJsonMainEntries == nil {
		r.packageJsonMainEntries = emptyPackageJsonMainEntries
	} else {
		r.packageJsonMainEntries = nodeModuleResolver_buildPackageJsonMainEntries(packageJsonMainEntries)
	}
	return r
}

func nodeModuleResolver_buildPackageJsonMainEntries(
	packageJsonMainEntries *jutil.LinkedHashMap[string, string]) *jutil.LinkedHashMap[string, string] {
	builder := jutil.NewLinkedHashMap[string, string]()
	for entryKey, value := range packageJsonMainEntries.All() {
		if ModuleLoader_IsAmbiguousIdentifier(entryKey) {
			entryKey = ModuleLoader_MODULE_SLASH + entryKey
		}
		if old, ok := builder.Lookup(entryKey); ok {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + entryKey + "=" + value + " and " + entryKey + "=" + old))
		}
		builder.Put(entryKey, value)
	}
	return builder
}

// Java: NodeModuleResolver.getPackageJsonMainEntries()
func (r *NodeModuleResolver) GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string] {
	return r.packageJsonMainEntries
}

// Java: NodeModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *NodeModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {
	loadAddress := r.internalResolveJsModule(scriptAddress, moduleAddress)
	if loadAddress == nil {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(&sourcename, lineno, colno, ModuleLoader_LOAD_WARNING, moduleAddress))
	}
	return loadAddress
}

func (r *NodeModuleResolver) internalResolveJsModule(scriptAddress string, moduleAddress string) *string {
	if ModuleLoader_IsAbsoluteIdentifier(moduleAddress) ||
		ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		return r.resolveJsModuleNodeFileOrDirectory(scriptAddress, moduleAddress)
	}
	return r.resolveJsModuleFromRegistry(scriptAddress, moduleAddress)
}

// Java: NodeModuleResolver.resolveJsModuleSilently(String, String)
func (r *NodeModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	return r.internalResolveJsModule(scriptAddress, moduleAddress)
}

func (r *NodeModuleResolver) resolveJsModuleFile(scriptAddress string, moduleAddress string) *string {
	for _, extension := range nodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH {
		moduleAddressCandidate := moduleAddress + extension
		canonicalizedCandidatePath := r.CanonicalizePath(scriptAddress, moduleAddressCandidate)

		if mapped, ok := r.packageJsonMainEntries.Lookup(canonicalizedCandidatePath); ok {
			if mapped == ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER {
				return nil
			}
			moduleAddressCandidate = mapped
		}

		if loadAddress := r.Locate(scriptAddress, moduleAddressCandidate); loadAddress != nil {
			return loadAddress
		}
	}
	return nil
}

func (r *NodeModuleResolver) resolveJsModuleNodeFileOrDirectory(scriptAddress string, moduleAddress string) *string {
	if loadAddress := r.resolveJsModuleFile(scriptAddress, moduleAddress); loadAddress != nil {
		return loadAddress
	}
	return r.resolveJsModuleNodeDirectory(scriptAddress, moduleAddress)
}

func (r *NodeModuleResolver) resolveJsModuleNodeDirectory(scriptAddress string, moduleAddress string) *string {
	moduleAddress = strings.TrimSuffix(moduleAddress, ModuleLoader_MODULE_SLASH)
	for _, suffix := range nodeModuleResolver_FILES_TO_SEARCH {
		if loadAddress := r.Locate(scriptAddress, moduleAddress+suffix); loadAddress != nil {
			if suffix == ModuleLoader_MODULE_SLASH+"package.json" {
				if mainEntry, ok := r.packageJsonMainEntries.Lookup(*loadAddress); ok {
					return r.resolveJsModuleFile(scriptAddress, mainEntry)
				}
			} else {
				return loadAddress
			}
		}
	}
	return nil
}

func (r *NodeModuleResolver) resolveJsModuleFromRegistry(scriptAddress string, moduleAddress string) *string {
	normalizedScriptAddress := scriptAddress
	if ModuleLoader_IsAmbiguousIdentifier(scriptAddress) {
		normalizedScriptAddress = ModuleLoader_MODULE_SLASH + scriptAddress
	}
	for _, nodeModulesFolder := range r.nodeModulesFolders {
		if !strings.HasPrefix(normalizedScriptAddress, nodeModulesFolder) {
			continue
		}
		if loadAddress := r.resolveJsModuleNodeFileOrDirectory(
			scriptAddress, nodeModulesFolder+"node_modules/"+moduleAddress); loadAddress != nil {
			return loadAddress
		}
	}
	return nil
}
