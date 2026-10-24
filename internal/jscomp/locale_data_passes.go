// Port of com/google/javascript/jscomp/LocaleDataPasses.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// LocaleDataPasses contains compiler passes to protect `goog.LOCALE` from
// optimization during the main optimizations phase of compilation, then
// replace it with the specific destination locale near the end of
// compilation.
//
// Go-only: language-adapter — the Java class has only a private
// constructor and static members; it has no Go type.

// LocaleDataPasses_GOOG_LOCALE_REPLACEMENT is the replacement for values
// that needed to be protected from optimizations.
const LocaleDataPasses_GOOG_LOCALE_REPLACEMENT = "__JSC_LOCALE__"

// LocaleDataPasses_ProtectGoogLocale is the pass that protects goog.LOCALE.
//
// Java: static class LocaleDataPasses.ProtectGoogLocale implements CompilerPass
type LocaleDataPasses_ProtectGoogLocale struct {
	compiler AbstractCompiler
}

// Java: LocaleDataPasses.ProtectGoogLocale(AbstractCompiler)
func NewLocaleDataPasses_ProtectGoogLocale(compiler AbstractCompiler) *LocaleDataPasses_ProtectGoogLocale {
	return &LocaleDataPasses_ProtectGoogLocale{compiler: compiler}
}

// Java: LocaleDataPasses.ProtectGoogLocale.process(Node, Node)
func (p *LocaleDataPasses_ProtectGoogLocale) Process(externs *rhino.Node, root *rhino.Node) {
	// Create the extern symbols
	NodeUtil_CreateSynthesizedExternsSymbol(p.compiler, LocaleDataPasses_GOOG_LOCALE_REPLACEMENT)
	protectLocaleCallback := newLocaleDataPasses_ProtectCurrentLocale(p.compiler)
	NodeTraversal_Traverse(p.compiler, root, protectLocaleCallback)
}

// localeDataPasses_ProtectCurrentLocale protects `goog.LOCALE` by replacing
// it with an extern-defined name.
//
// Java: private static class LocaleDataPasses.ProtectCurrentLocale extends AbstractPostOrderCallback
type localeDataPasses_ProtectCurrentLocale struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler   AbstractCompiler
	astFactory *AstFactory
}

// Java: LocaleDataPasses.ProtectCurrentLocale(AbstractCompiler)
func newLocaleDataPasses_ProtectCurrentLocale(compiler AbstractCompiler) *localeDataPasses_ProtectCurrentLocale {
	return &localeDataPasses_ProtectCurrentLocale{
		compiler:   compiler,
		astFactory: compiler.CreateAstFactory(),
	}
}

// Java: LocaleDataPasses.ProtectCurrentLocale.visit(NodeTraversal, Node, Node)
func (c *localeDataPasses_ProtectCurrentLocale) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if localeDataPasses_isGoogDotLocaleReference(n) {
		// We will replace the RHS of `goog.LOCALE = goog.define(...);`, but everywhere
		// else we will replace `goog.LOCALE` itself, so we don't have to waste time
		// inlining it later.
		nodeToReplace := n
		if NodeUtil_IsLhsOfAssign(n) {
			nodeToReplace = n.GetNext()
		}
		replacement := c.astFactory.CreateName(LocaleDataPasses_GOOG_LOCALE_REPLACEMENT, AstFactory_Type_(nodeToReplace))
		replacement.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		nodeToReplace.ReplaceWith(replacement)
		c.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// localeDataPasses_QNAME_FOR_GOOG_LOCALE: matching against an actual Node is
// faster than matching against the string "goog.LOCALE"
var localeDataPasses_QNAME_FOR_GOOG_LOCALE = rhino.IR_Getprop(rhino.IR_Name("goog"), "LOCALE")

// Java: LocaleDataPasses.isGoogDotLocaleReference(Node)
func localeDataPasses_isGoogDotLocaleReference(n *rhino.Node) bool {
	// NOTE: Theoretically there could be a local variable named `goog`, but it's not worth checking
	// for that.
	return n.MatchesQualifiedNameNode(localeDataPasses_QNAME_FOR_GOOG_LOCALE)
}

const localeDataPasses_LocaleSubstitutions_DEFAULT_LOCALE = "en"

// LocaleDataPasses_LocaleSubstitutions replaces `__JSC_LOCALE__` with the
// actual locale string.
//
// Java: static class LocaleDataPasses.LocaleSubstitutions extends AbstractPostOrderCallback implements CompilerPass
type LocaleDataPasses_LocaleSubstitutions struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	qnameForLocale *rhino.Node

	compiler   AbstractCompiler
	astFactory *AstFactory
	locale     *string
}

// Java: LocaleDataPasses.LocaleSubstitutions(AbstractCompiler, String)
func NewLocaleDataPasses_LocaleSubstitutions(compiler AbstractCompiler, locale *string) *LocaleDataPasses_LocaleSubstitutions {
	s := &LocaleDataPasses_LocaleSubstitutions{}
	s.qnameForLocale = rhino.IR_Name(LocaleDataPasses_GOOG_LOCALE_REPLACEMENT)
	s.compiler = compiler
	s.astFactory = compiler.CreateAstFactory()
	// Use the "default" locale if not otherwise set.
	if locale == nil {
		s.locale = jutil.StringPtr(localeDataPasses_LocaleSubstitutions_DEFAULT_LOCALE)
	} else {
		s.locale = locale
	}
	return s
}

// Java: LocaleDataPasses.LocaleSubstitutions.process(Node, Node)
func (s *LocaleDataPasses_LocaleSubstitutions) Process(externs *rhino.Node, root *rhino.Node) {
	// Create the extern symbol
	NodeTraversal_Traverse(s.compiler, root, s)
}

// Java: LocaleDataPasses.LocaleSubstitutions.visit(NodeTraversal, Node, Node)
func (s *LocaleDataPasses_LocaleSubstitutions) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.MatchesNameNode(s.qnameForLocale) {
		replacement := s.astFactory.CreateString(*jutil.CheckNotNull(s.locale)).Srcref(n)
		n.ReplaceWith(replacement)
		s.compiler.ReportChangeToEnclosingScope(replacement)
	}
}
