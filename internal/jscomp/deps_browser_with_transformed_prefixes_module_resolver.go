// Port of com/google/javascript/jscomp/deps/BrowserWithTransformedPrefixesModuleResolver.java @ closure-compiler v20261021

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

type browserPrefixReplacement struct {
	prefix      string
	replacement string
}

// BrowserWithTransformedPrefixesModuleResolver is a limited superset of the
// BrowserModuleResolver that allows for replacing some path prefixes before
// resolving.
type BrowserWithTransformedPrefixesModuleResolver struct {
	moduleResolverCommon
	prefixReplacements []browserPrefixReplacement
	expectedPrefixes   string
}

var _ ModuleResolver = (*BrowserWithTransformedPrefixesModuleResolver)(nil)

var BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH = DiagnosticType_Error(
	"JSC_INVALID_AMBIGUOUS_PATH",
	"The path \"{0}\" is invalid. Expected any of the following prefixes for non-relative "+
		"paths: {1}.")

// BrowserWithTransformedPrefixesModuleResolver_Factory is the factory for
// BrowserWithTransformedPrefixesModuleResolver.
type BrowserWithTransformedPrefixesModuleResolver_Factory struct {
	prefixReplacements *jutil.LinkedHashMap[string, string]
}

var _ ModuleLoader_ModuleResolverFactory = (*BrowserWithTransformedPrefixesModuleResolver_Factory)(nil)

// Java: BrowserWithTransformedPrefixesModuleResolver.Factory(ImmutableMap<String, String>)
func NewBrowserWithTransformedPrefixesModuleResolver_Factory(prefixReplacements *jutil.LinkedHashMap[string, string]) *BrowserWithTransformedPrefixesModuleResolver_Factory {
	return &BrowserWithTransformedPrefixesModuleResolver_Factory{prefixReplacements: prefixReplacements}
}

// Java: BrowserWithTransformedPrefixesModuleResolver.Factory.create(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func (f *BrowserWithTransformedPrefixesModuleResolver_Factory) Create(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	return NewBrowserWithTransformedPrefixesModuleResolver(
		modulePaths, moduleRootPaths, errorHandler, pathEscaper, f.prefixReplacements)
}

// Java: BrowserWithTransformedPrefixesModuleResolver(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper, ImmutableMap<String, String>)
func NewBrowserWithTransformedPrefixesModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper,
	prefixReplacements *jutil.LinkedHashMap[string, string]) *BrowserWithTransformedPrefixesModuleResolver {
	r := &BrowserWithTransformedPrefixesModuleResolver{
		moduleResolverCommon: newModuleResolverCommon(modulePaths, moduleRootPaths, errorHandler, pathEscaper),
	}
	n := prefixReplacements.Size()
	replacements := make([]browserPrefixReplacement, 0, n)
	prefixes := make([]string, 0, n)
	for k, v := range prefixReplacements.All() {
		replacements = append(replacements, browserPrefixReplacement{prefix: k, replacement: v})
		prefixes = append(prefixes, k)
	}
	slices.SortFunc(replacements, func(a, b browserPrefixReplacement) int {
		if c := cmp.Compare(jutil.StringLength(b.prefix), jutil.StringLength(a.prefix)); c != 0 {
			return c
		}
		return jutil.StringCompareTo(a.prefix, b.prefix)
	})
	r.prefixReplacements = replacements
	slices.SortFunc(prefixes, jutil.StringCompareTo)
	r.expectedPrefixes = strings.Join(prefixes, ", ")
	return r
}

func (r *BrowserWithTransformedPrefixesModuleResolver) replacePathPrefixes(moduleAddress string) (string, bool) {
	for _, pr := range r.prefixReplacements {
		if strings.HasPrefix(moduleAddress, pr.prefix) {
			return pr.replacement + moduleAddress[len(pr.prefix):], true
		}
	}
	return moduleAddress, false
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {
	transformedAddress, replaced := r.replacePathPrefixes(moduleAddress)
	//Note: in Java, `replaced = !transformedAddress.equals(moduleAddress)` (so replacing "" -> "" or "x" -> "x" has replaced == false!)
	replaced = replaced && transformedAddress != moduleAddress
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(
				&sourcename,
				lineno,
				colno,
				BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH,
				transformedAddress,
				r.expectedPrefixes))
		return nil
	}
	return r.Locate(scriptAddress, transformedAddress)
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveJsModuleSilently(String, String)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	transformedAddress, replaced := r.replacePathPrefixes(moduleAddress)
	replaced = replaced && transformedAddress != moduleAddress
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		return nil
	}
	return r.Locate(scriptAddress, transformedAddress)
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveModuleAsPath(String, String)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		return r.moduleResolverCommon.ResolveModuleAsPath(scriptAddress, moduleAddress)
	}
	transformedAddress, _ := r.replacePathPrefixes(moduleAddress)
	return ModuleLoader_Normalize(transformedAddress, r.moduleRootPaths)
}
