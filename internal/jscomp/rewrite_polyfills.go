// Port of com/google/javascript/jscomp/RewritePolyfills.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewritePolyfills injects polyfill libraries to ensure that ES6+ library
// functions are available.
//
// Also runs if polyfill isolation is enabled, even if polyfill injection is
// disabled, in order to prevent deletion of a required library function by
// dead code elimination.
//
// TODO(b/120486392): consider merging this pass with InjectRuntimeLibraries
// and InjectTranspilationRuntimeLibraries.
type RewritePolyfills struct {
	compiler            AbstractCompiler
	runtimeJsLibManager *RuntimeJsLibManager
	polyfills           *PolyfillUsageFinder_Polyfills
	injectPolyfills     bool
	isolatePolyfills    bool
	libraries           *jutil.LinkedHashSet[string]
	// Go-only: language-adapter — a nullable LanguageMode: nil is null.
	injectPolyfillsNewerThan *CompilerOptions_LanguageMode
}

// Java: RewritePolyfills.INSUFFICIENT_OUTPUT_VERSION_ERROR
var RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR = DiagnosticType_Disabled(
	"JSC_INSUFFICIENT_OUTPUT_VERSION",
	"Built-in ''{0}'' not supported in output version {1}")

// Java: RewritePolyfills.JSCOMP_POLYFILL
var rewritePolyfills_JSCOMP_POLYFILL = rhino.QualifiedName_Of("$jscomp.polyfill")

// NewRewritePolyfills: if injectPolyfills, injects $jscomp.polyfill
// initializations into the first input; if false, no polyfills are
// injected. If isolatePolyfills, adds externs for library functions used by
// IsolatePolyfills to prevent their deletion.
//
// Java: RewritePolyfills(AbstractCompiler, boolean, boolean, LanguageMode)
func NewRewritePolyfills(
	compiler AbstractCompiler,
	injectPolyfills bool,
	isolatePolyfills bool,
	injectPolyfillsNewerThan *CompilerOptions_LanguageMode) *RewritePolyfills {
	return NewRewritePolyfillsWithRuntimeJsLibManager(
		compiler,
		compiler.GetRuntimeJsLibManager(),
		defaultPolyfillsTable(),
		injectPolyfills,
		isolatePolyfills,
		injectPolyfillsNewerThan)
}

// NewRewritePolyfillsWithRuntimeJsLibManager is for unit testing: allows
// instantiating RewritePolyfills with a different RuntimeJsLibManager than
// AbstractCompiler#getRuntimeJsLibManager.
//
// Java: RewritePolyfills(AbstractCompiler, RuntimeJsLibManager, Polyfills, boolean, boolean, LanguageMode)
func NewRewritePolyfillsWithRuntimeJsLibManager(
	compiler AbstractCompiler,
	runtimeJsLibManager *RuntimeJsLibManager,
	polyfills *PolyfillUsageFinder_Polyfills,
	injectPolyfills bool,
	isolatePolyfills bool,
	injectPolyfillsNewerThan *CompilerOptions_LanguageMode) *RewritePolyfills {
	r := &RewritePolyfills{}
	r.compiler = compiler
	r.runtimeJsLibManager = runtimeJsLibManager
	r.polyfills = polyfills
	r.injectPolyfills = injectPolyfills
	r.isolatePolyfills = isolatePolyfills
	r.injectPolyfillsNewerThan = injectPolyfillsNewerThan
	return r
}

// Java: RewritePolyfills.process(Node, Node)
func (r *RewritePolyfills) Process(externs rhino.Node, root rhino.Node) {
	if r.isolatePolyfills {
		// Polyfill isolation requires a pass to run near the end of optimizations. That pass may call
		// into a library method injected in this pass. Adding an externs declaration of that library
		// method prevents it from being dead-code-elimiated before polyfill isolation runs.
		jscompLookupMethodDecl := rhino.IR_Var(rhino.IR_Name("$jscomp$lookupPolyfilledValue"))
		synthesizedExternsAstRoot :=
			r.compiler.GetSynthesizedExternsInput().GetAstRoot(r.compiler)
		jscompLookupMethodDecl.SrcrefTree(synthesizedExternsAstRoot)
		synthesizedExternsAstRoot.AddChildToBack(jscompLookupMethodDecl)
		r.compiler.ReportChangeToEnclosingScope(jscompLookupMethodDecl)
	}

	if !r.injectPolyfills && r.injectPolyfillsNewerThan == nil {
		// Nothing left to do. Probably this pass only needed to run because --isolate_polyfills is
		// enabled but not --rewrite_polyfills.
		return
	}
	if r.injectPolyfills {
		r.libraries = jutil.NewLinkedHashSet[string]()
		NewPolyfillUsageFinder(r.compiler, r.polyfills).TraverseExcludingGuarded(root, r.Inject)
	}

	var librariesToInject []string
	if r.injectPolyfillsNewerThan != nil {
		polyfillsToInject :=
			r.polyfills.GetPolyfillsNewerThan(*r.injectPolyfillsNewerThan)
		librariesToInject = []string{}
		for _, p := range polyfillsToInject {
			// Skip polyfills that have no associated library. This is true for language
			// features like `Proxy` and `String.raw` that have no associated polyfill, hence
			// there's
			// nothing to inject here.
			if p.library != "" {
				librariesToInject = append(librariesToInject, p.library)
			}
		}
	} else {
		librariesToInject = r.libraries.ToSlice()
	}

	r.InjectAll(librariesToInject /* forceInjection= */, r.injectPolyfillsNewerThan != nil)
}

// Java: RewritePolyfills.injectAll(Iterable<String>, boolean)
func (r *RewritePolyfills) InjectAll(librariesToInject []string, forceInjection bool) {
	var lastNode rhino.Node
	for _, library := range librariesToInject {
		// Go-only: language-adapter — checkNotNull(library): a Go string
		// cannot be null.
		jutil.CheckStateMsg(library != "", "unexpected empty library")
		lastNode = r.runtimeJsLibManager.EnsureLibraryInjected(library, forceInjection)
	}
	if !lastNode.IsNil() {
		parent := lastNode.GetParent()
		r.RemoveUnneededPolyfills(parent, lastNode.GetNext())
		r.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// RemoveUnneededPolyfills removes any $jscomp.polyfill calls whose 3rd
// parameter (the language version that already contains the library) is
// the same or lower than languageOut.
//
// Java: RewritePolyfills.removeUnneededPolyfills(Node, Node)
func (r *RewritePolyfills) RemoveUnneededPolyfills(parent rhino.Node, runtimeEnd rhino.Node) {
	node := parent.GetFirstChild()
	// The target environment is assumed to support the features in this set.
	outputFeatureSet := r.compiler.GetOptions().GetOutputFeatureSet()
	for !node.IsNil() && node != runtimeEnd {
		// look up the next node now, because we may be removing this one.
		next := node.GetNext()
		polyfillSupportedFeatureSet := r.GetPolyfillSupportedFeatureSet(node)
		if polyfillSupportedFeatureSet != nil &&
			outputFeatureSet.Contains(polyfillSupportedFeatureSet) {
			NodeUtil_RemoveChild(parent, node)
			NodeUtil_MarkFunctionsDeleted(node, r.compiler)
		}
		node = next
	}
}

// GetPolyfillSupportedFeatureSet returns, if the given `Node` is a polyfill
// definition, the `FeatureSet` which should be considered to already
// include that polyfill (making it unnecessary).
//
// Otherwise, returns `null`.
//
// Java: RewritePolyfills.getPolyfillSupportedFeatureSet(Node)
func (r *RewritePolyfills) GetPolyfillSupportedFeatureSet(maybePolyfill rhino.Node) *parser.FeatureSet {
	stmt, ok := ast.AsExprStmt(maybePolyfill)
	if !ok {
		return nil
	}
	call, ok := stmt.Expr().AsCall()
	if !ok || !rewritePolyfills_JSCOMP_POLYFILL.Matches(call.Callee().Node()) {
		return nil
	}
	arg, ok := call.AsCallOrNew().ArgAt(2)
	if !ok {
		return nil
	}
	nativeVersionStr := arg.MustStringLit().Value()
	return PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr)
}

// Java: RewritePolyfills.inject(PolyfillUsage)
func (r *RewritePolyfills) Inject(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	polyfill := polyfillUsage.Polyfill()
	outputFeatureSet := r.compiler.GetOptions().GetOutputFeatureSet()
	featuresRequiredByPolyfill := parser.FeatureSet_ValueOf(polyfill.polyfillVersion)
	if polyfill.kind == PolyfillUsageFinder_Polyfill_Kind_STATIC &&
		!outputFeatureSet.Contains(featuresRequiredByPolyfill) {
		r.compiler.Report(
			JSError_MakeWithNode(
				polyfillUsage.Node(),
				RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR,
				polyfillUsage.Name(),
				outputFeatureSet.Version()))
	}

	// The question we want to ask here is:
	// "Does the target platform already have the symbol this polyfill provides?"
	// We approximate it by asking instead:
	// "Does the target platform support all of the features that existed in the language
	// version that introduced this symbol?"
	if !outputFeatureSet.Contains(parser.FeatureSet_ValueOf(polyfill.nativeVersion)) &&
		polyfill.library != "" {
		r.libraries.Add(polyfill.library)
	}
}
