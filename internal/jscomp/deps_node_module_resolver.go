// Port of com/google/javascript/jscomp/deps/NodeModuleResolver.java @ closure-compiler v20260415

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// NodeModuleResolver is the resolution algorithm for NodeJS. See
// https://nodejs.org/api/modules.html#modules_all_together
//
// Unambiguous paths are file paths resolved from the current script.
// Ambiguous paths are located within the nearest node_modules folder
// ancestor.
type NodeModuleResolver struct {
	ModuleResolverBase

	// Named modules found in node_modules folders
	packageJsonMainEntries *jutil.LinkedHashMap[string, string]

	// Named modules found in node_modules folders
	nodeModulesFolders []string
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

// nodeModuleResolver_buildNodeModulesFoldersRegistry builds a list of node
// module paths. Given the following path:
//
// /foo/node_modules/bar/node_modules/baz/foo_bar_baz.js
//
// Return a set containing:
//
// /foo/ /foo/node_modules/bar/
//
// modulePaths is the set of all module paths where the key is the module
// path normalized to have a leading slash; moduleRootPaths is a possibly
// empty list of root paths that should be ignored when processing module
// paths. It returns a sorted set with the longest paths first where each
// entry is the folder containing a node_modules sub-folder.
//
// Java: NodeModuleResolver.buildNodeModulesFoldersRegistry(Iterable<String>, Iterable<String>)
func nodeModuleResolver_buildNodeModulesFoldersRegistry(
	modulePaths *jutil.LinkedHashSet[string], moduleRootPaths []string) []string {
	registry :=
		jutil.NewTreeSet(
			// TODO(b/28382956): Take better advantage of Java8 comparing() to simplify this
			func(a, b string) int {
				// Order longest path first
				comparison := cmp.Compare(jutil.StringLength(b), jutil.StringLength(a)) // Integer.compare
				if comparison != 0 {
					return comparison
				}

				return jutil.StringCompareTo(a, b)
			})

	// For each modulePath, find all the node_modules folders
	// There might be more than one:
	//    /foo/node_modules/bar/node_modules/baz/foo_bar_baz.js
	// Should add:
	//   /foo/ -> bar/node_modules/baz/foo_bar_baz.js
	//   /foo/node_modules/bar/ -> baz/foo_bar_baz.js
	//
	// If there are root paths provided - they should be ignored from paths. For example if
	// root paths contains "/generated" root then the following module path should produce exact
	// same result as above:
	//    /generated/foo/node_modules/bar/node_modules/baz/foo_bar_baz.js
	for modulePath := range modulePaths.All() {
		// Strip root path from the beginning if it matches any.
		for _, moduleRootPath := range moduleRootPaths {
			if strings.HasPrefix(modulePath, moduleRootPath) {
				modulePath = modulePath[len(moduleRootPath):]
				break
			}
		}
		nodeModulesDirs := jutil.StringSplitLiteral(modulePath, "/node_modules/")
		parentPath := ""

		for i := 0; i < len(nodeModulesDirs)-1; i++ {
			if i+1 < len(nodeModulesDirs) {
				parentPath += nodeModulesDirs[i] + "/"
			}

			registry.Add(parentPath)
			parentPath += "node_modules/"
		}
	}

	return registry.ToSlice()
}

// NodeModuleResolver_Factory is the factory for NodeModuleResolver.
type NodeModuleResolver_Factory struct {
	packageJsonMainEntries *jutil.LinkedHashMap[string, string]
}

var _ ModuleLoader_ModuleResolverFactory = (*NodeModuleResolver_Factory)(nil)

// Java: NodeModuleResolver.Factory()
func NewNodeModuleResolver_Factory() *NodeModuleResolver_Factory {
	return NewNodeModuleResolver_FactoryWithEntries( /* packageJsonMainEntries= */ nil)
}

// NewNodeModuleResolver_FactoryWithEntries takes a possibly-nil map.
//
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

// NewNodeModuleResolver takes a possibly-nil packageJsonMainEntries.
//
// Java: NodeModuleResolver(ImmutableSet<String>, ImmutableList<String>, Map<String, String>, ErrorHandler, PathEscaper)
func NewNodeModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	packageJsonMainEntries *jutil.LinkedHashMap[string, string],
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) *NodeModuleResolver {
	r := &NodeModuleResolver{}
	r.InitModuleResolverBase(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	r.nodeModulesFolders = nodeModuleResolver_buildNodeModulesFoldersRegistry(modulePaths, moduleRootPaths)

	if packageJsonMainEntries == nil {
		r.packageJsonMainEntries = jutil.NewLinkedHashMap[string, string]()
	} else {
		r.packageJsonMainEntries = r.BuildPackageJsonMainEntries(packageJsonMainEntries)
	}
	return r
}

// BuildPackageJsonMainEntries takes packageJsonMainEntries, a map with keys
// that are package.json file paths and values which are the "main" entry
// from the package.json. "main" entries are absolute paths rooted from the
// folder containing the package.json file.
//
// Java: NodeModuleResolver.buildPackageJsonMainEntries(Map<String, String>)
func (r *NodeModuleResolver) BuildPackageJsonMainEntries(
	packageJsonMainEntries *jutil.LinkedHashMap[string, string]) *jutil.LinkedHashMap[string, string] {
	builder := jutil.NewLinkedHashMap[string, string]()
	for entryKey, value := range packageJsonMainEntries.All() {
		if ModuleLoader_IsAmbiguousIdentifier(entryKey) {
			entryKey = ModuleLoader_MODULE_SLASH + entryKey
		}

		// ImmutableMap.Builder.put + buildOrThrow(): duplicate keys throw.
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

// ResolveJsModule may return nil.
//
// Java: NodeModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *NodeModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {
	loadAddress := r.InternalResolveJsModule(scriptAddress, moduleAddress)

	if loadAddress == nil {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(&sourcename, lineno, colno, ModuleLoader_LOAD_WARNING, moduleAddress))
	}
	return loadAddress
}

// Java: NodeModuleResolver.internalResolveJsModule(String, String)
func (r *NodeModuleResolver) InternalResolveJsModule(scriptAddress string, moduleAddress string) *string {
	var loadAddress *string
	if ModuleLoader_IsAbsoluteIdentifier(moduleAddress) ||
		ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		loadAddress = r.ResolveJsModuleNodeFileOrDirectory(scriptAddress, moduleAddress)
	} else {
		loadAddress = r.ResolveJsModuleFromRegistry(scriptAddress, moduleAddress)
	}

	return loadAddress
}

// ResolveJsModuleSilently may return nil.
//
// Java: NodeModuleResolver.resolveJsModuleSilently(String, String)
func (r *NodeModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	return r.InternalResolveJsModule(scriptAddress, moduleAddress)
}

// ResolveJsModuleFile may return nil.
//
// Java: NodeModuleResolver.resolveJsModuleFile(String, String)
func (r *NodeModuleResolver) ResolveJsModuleFile(scriptAddress string, moduleAddress string) *string {
	for _, extension := range nodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH {
		moduleAddressCandidate := moduleAddress + extension
		canonicalizedCandidatePath := r.CanonicalizePath(scriptAddress, moduleAddressCandidate)

		// Also look for mappings in packageJsonMainEntries because browser field
		// advanced usage allows to override / skip specific files, including
		// the main entry.
		if r.packageJsonMainEntries.ContainsKey(canonicalizedCandidatePath) {
			moduleAddressCandidate = r.packageJsonMainEntries.Get(canonicalizedCandidatePath)

			if ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER == moduleAddressCandidate {
				return nil
			}
		}

		loadAddress := r.Locate(scriptAddress, moduleAddressCandidate)
		if loadAddress != nil {
			return loadAddress
		}
	}

	return nil
}

// Java: NodeModuleResolver.resolveJsModuleNodeFileOrDirectory(String, String)
func (r *NodeModuleResolver) ResolveJsModuleNodeFileOrDirectory(
	scriptAddress string, moduleAddress string) *string {
	var loadAddress *string
	loadAddress = r.ResolveJsModuleFile(scriptAddress, moduleAddress)
	if loadAddress == nil {
		loadAddress = r.ResolveJsModuleNodeDirectory(scriptAddress, moduleAddress)
	}
	return loadAddress
}

// Java: NodeModuleResolver.resolveJsModuleNodeDirectory(String, String)
func (r *NodeModuleResolver) ResolveJsModuleNodeDirectory(
	scriptAddress string, moduleAddress string) *string {
	if strings.HasSuffix(moduleAddress, ModuleLoader_MODULE_SLASH) {
		moduleAddress = moduleAddress[0 : len(moduleAddress)-1]
	}

	// Load as a file
	for i := 0; i < len(nodeModuleResolver_FILES_TO_SEARCH); i++ {
		loadAddress := r.Locate(scriptAddress, moduleAddress+nodeModuleResolver_FILES_TO_SEARCH[i])
		if loadAddress != nil {
			if nodeModuleResolver_FILES_TO_SEARCH[i] == ModuleLoader_MODULE_SLASH+"package.json" {
				if r.packageJsonMainEntries.ContainsKey(*loadAddress) {
					return r.ResolveJsModuleFile(scriptAddress, r.packageJsonMainEntries.Get(*loadAddress))
				}
			} else {
				return loadAddress
			}
		}
	}

	return nil
}

// Java: NodeModuleResolver.resolveJsModuleFromRegistry(String, String)
func (r *NodeModuleResolver) ResolveJsModuleFromRegistry(scriptAddress string, moduleAddress string) *string {
	prefix := ""
	if ModuleLoader_IsAmbiguousIdentifier(scriptAddress) {
		prefix = ModuleLoader_MODULE_SLASH
	}
	normalizedScriptAddress := prefix + scriptAddress
	for _, nodeModulesFolder := range r.nodeModulesFolders {

		if !strings.HasPrefix(normalizedScriptAddress, nodeModulesFolder) {
			continue
		}

		// Load as a file
		fullModulePath := nodeModulesFolder + "node_modules/" + moduleAddress
		loadAddress := r.ResolveJsModuleFile(scriptAddress, fullModulePath)
		if loadAddress == nil {
			// Load as a directory
			loadAddress = r.ResolveJsModuleNodeDirectory(scriptAddress, fullModulePath)
		}

		if loadAddress != nil {
			return loadAddress
		}
	}

	return nil
}
