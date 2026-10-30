// Port of com/google/javascript/jscomp/LazyParsedDependencyInfo.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// LazyParsedDependencyInfo is a DependencyInfo class that determines load
// flags by parsing the AST just-in-time.
type LazyParsedDependencyInfo struct {
	delegate DependencyInfo
	ast      *CompilerInput
	compiler AbstractCompiler

	loadFlags    LoadFlags
	hasLoadFlags bool
}

var _ DependencyInfo = (*LazyParsedDependencyInfo)(nil)

// Java: LazyParsedDependencyInfo(DependencyInfo, CompilerInput, AbstractCompiler)
func NewLazyParsedDependencyInfo(
	delegate DependencyInfo, ast *CompilerInput, compiler AbstractCompiler) *LazyParsedDependencyInfo {
	return &LazyParsedDependencyInfo{
		delegate: jutil.CheckNotNull(delegate),
		ast:      jutil.CheckNotNull(ast),
		compiler: jutil.CheckNotNull(compiler),
	}
}

// Java: LazyParsedDependencyInfo.isEs6Module()
func (d *LazyParsedDependencyInfo) IsEs6Module() bool {
	// Instead of doing a full parse to read all load flags, just ask the delegate, which at least
	// has this much info
	return d.delegate.IsEs6Module()
}

// Java: LazyParsedDependencyInfo.isGoogModule()
func (d *LazyParsedDependencyInfo) IsGoogModule() bool {
	// Instead of doing a full parse to read all load flags, just ask the delegate, which at least
	// has this much info
	return d.delegate.IsGoogModule()
}

// Java: LazyParsedDependencyInfo.getLoadFlags()
func (d *LazyParsedDependencyInfo) GetLoadFlags() LoadFlags {
	if !d.hasLoadFlags {
		hadSourceInMemory := d.ast.GetSourceFile().HasSourceInMemory()
		delegateFlags := d.delegate.GetLoadFlags()
		features := d.ast.GetFeatures(d.compiler)
		hasModules := features.Has(parser.FeatureSet_Feature_MODULES)
		version := features.Version()

		var langFlag, modFlag string
		hasLang, hasMod := false, false
		onlyLangAndMod := true
		for _, e := range delegateFlags.Entries() {
			switch e.Key {
			case "lang":
				langFlag, hasLang = e.Value, true
			case "module":
				modFlag, hasMod = e.Value, true
			default:
				onlyLangAndMod = false
			}
		}

		if hasModules {
			if hasMod && modFlag != "es6" {
				name := d.GetName()
				d.compiler.Report(
					JSError_MakeWithSourceLocation(
						&name,
						/* lineno= */ -1,
						/* charno= */ -1,
						ModuleLoader_MODULE_CONFLICT,
						d.GetName()))
			}
			modFlag, hasMod = "es6", true
		}
		if version != "es3" {
			langFlag, hasLang = version, true
		}

		if onlyLangAndMod && (hasLang || langFlag == "") && (hasMod || modFlag == "") {
			// Fast path: reuse canonical LoadFlags in UTF-16 key order ("lang" < "module").
			d.loadFlags = simpleDependencyInfo_flagsForLangModule(langFlag, modFlag)
		} else {
			entries := make([]LoadFlag, 0, delegateFlags.Len()+2)
			for _, e := range delegateFlags.Entries() {
				putLoadFlag(&entries, e.Key, e.Value)
			}
			if hasMod {
				putLoadFlag(&entries, "module", modFlag)
			}
			if hasLang {
				putLoadFlag(&entries, "lang", langFlag)
			}
			slices.SortFunc(entries, func(a, b LoadFlag) int {
				return jutil.StringCompareTo(a.Key, b.Key)
			})
			d.loadFlags, _ = normalizeLoadFlagsEntries(entries, false)
		}
		d.hasLoadFlags = true

		// Don't preserve the source in memory if it was loaded only for getLoadFlags().
		if !hadSourceInMemory {
			d.ast.GetSourceFile().ClearCachedSource()
		}
		// Don't preserve the full AST longer than necessary.  It can consume a lot of memory.
		d.ast = nil
	}
	return d.loadFlags
}

// Java: LazyParsedDependencyInfo.getName()
func (d *LazyParsedDependencyInfo) GetName() string {
	return d.delegate.GetName()
}

// Java: LazyParsedDependencyInfo.getPathRelativeToClosureBase()
func (d *LazyParsedDependencyInfo) GetPathRelativeToClosureBase() string {
	return d.delegate.GetPathRelativeToClosureBase()
}

// Java: LazyParsedDependencyInfo.getRequires()
func (d *LazyParsedDependencyInfo) GetRequires() []*DependencyInfo_Require {
	return d.delegate.GetRequires()
}

// Java: LazyParsedDependencyInfo.getTypeRequires()
func (d *LazyParsedDependencyInfo) GetTypeRequires() []string {
	return d.delegate.GetTypeRequires()
}

// Java: LazyParsedDependencyInfo.getProvides()
func (d *LazyParsedDependencyInfo) GetProvides() []string {
	return d.delegate.GetProvides()
}

// Java: LazyParsedDependencyInfo.getHasExternsAnnotation()
func (d *LazyParsedDependencyInfo) GetHasExternsAnnotation() bool {
	return d.delegate.GetHasExternsAnnotation()
}

// Java: LazyParsedDependencyInfo.getHasNoCompileAnnotation()
func (d *LazyParsedDependencyInfo) GetHasNoCompileAnnotation() bool {
	return d.delegate.GetHasNoCompileAnnotation()
}

// Java: DependencyInfo.getRequiredSymbols() (default method)
func (d *LazyParsedDependencyInfo) GetRequiredSymbols() []string {
	return DependencyInfo_GetRequiredSymbols(d)
}
