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
// Port of com/google/javascript/jscomp/LocaleDataPasses.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.LocaleDataPasses
type LocaleDataPasses struct {
	idHash_ atomic.Int32
}

func (this_ *LocaleDataPasses) BindSelf_(s any) {
	_ = s
}
func (this_ *LocaleDataPasses) IsJavaNil_() bool { return this_ == nil }
func (*LocaleDataPasses) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LocaleDataPasses", "LocaleDataPasses"
}
func (this_ *LocaleDataPasses) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LocaleDataPasses) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LocaleDataPasses", "com.google.javascript.jscomp.LocaleDataPasses", func(v_ any) bool {
		return j2gort.InstanceOf[*LocaleDataPasses](v_)
	})
}

func StaticGet_LocaleDataPasses_GOOG_LOCALE_REPLACEMENT() string { return "__JSC_LOCALE__" }

var static_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE pkg_rhino.Node

func StaticGet_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE() pkg_rhino.Node {
	clinit_LocaleDataPasses()
	return static_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE
}
func StaticSet_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE(v pkg_rhino.Node) {
	clinit_LocaleDataPasses()
	static_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE = v
}

var clinitDone_LocaleDataPasses uint32
var clinitState_LocaleDataPasses j2gort.ClinitState

func clinit_LocaleDataPasses() {
	if atomic.LoadUint32(&clinitDone_LocaleDataPasses) != 0 {
		return
	}
	clinitSlow_LocaleDataPasses()
}

func clinitSlow_LocaleDataPasses() {
	if !j2gort.EnterClinit(&clinitState_LocaleDataPasses, &clinitDone_LocaleDataPasses) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LocaleDataPasses, &clinitDone_LocaleDataPasses)
	static_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "LOCALE")
}

func (this_ *LocaleDataPasses) initFields_LocaleDataPasses() {
	_ = this_
}

// Java constructor: LocaleDataPasses()
func NewLocaleDataPasses() *LocaleDataPasses {
	this_ := &LocaleDataPasses{}
	InitLocaleDataPasses(this_)
	return this_
}

func InitLocaleDataPasses(this_ *LocaleDataPasses) {
	this_.initFields_LocaleDataPasses()
}

// Java static method: LocaleDataPasses.isGoogDotLocaleReference(com.google.javascript.rhino.Node n)
func LocaleDataPasses_IsGoogDotLocaleReference(n pkg_rhino.Node) bool {
	_ = n
	return n.MatchesQualifiedName_1_Node(StaticGet_LocaleDataPasses_QNAME_FOR_GOOG_LOCALE())
}

// Java class: com.google.javascript.jscomp.LocaleDataPasses.ProtectGoogLocale
type LocaleDataPasses_ProtectGoogLocale struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *LocaleDataPasses_ProtectGoogLocale) BindSelf_(s any) {
	_ = s
}
func (this_ *LocaleDataPasses_ProtectGoogLocale) IsJavaNil_() bool { return this_ == nil }
func (*LocaleDataPasses_ProtectGoogLocale) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LocaleDataPasses$ProtectGoogLocale", "ProtectGoogLocale"
}
func (this_ *LocaleDataPasses_ProtectGoogLocale) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LocaleDataPasses_ProtectGoogLocale) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LocaleDataPasses.ProtectGoogLocale", "com.google.javascript.jscomp.LocaleDataPasses$ProtectGoogLocale", func(v_ any) bool {
		return j2gort.InstanceOf[*LocaleDataPasses_ProtectGoogLocale](v_)
	})
}
func (*LocaleDataPasses_ProtectGoogLocale) Implements_jscomp_CompilerPass() {}

func (this_ *LocaleDataPasses_ProtectGoogLocale) initFields_LocaleDataPasses_ProtectGoogLocale() {
	_ = this_
}

// Java constructor: ProtectGoogLocale(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewLocaleDataPasses_ProtectGoogLocale(compiler AbstractCompiler) *LocaleDataPasses_ProtectGoogLocale {
	this_ := &LocaleDataPasses_ProtectGoogLocale{}
	InitLocaleDataPasses_ProtectGoogLocale(this_, compiler)
	return this_
}

func InitLocaleDataPasses_ProtectGoogLocale(this_ *LocaleDataPasses_ProtectGoogLocale, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_LocaleDataPasses_ProtectGoogLocale()
	this_.compiler = compiler
}

// Java method: ProtectGoogLocale.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *LocaleDataPasses_ProtectGoogLocale) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeUtil_CreateSynthesizedExternsSymbol(this_.compiler, "__JSC_LOCALE__")
	var protectLocaleCallback *LocaleDataPasses_ProtectCurrentLocale = NewLocaleDataPasses_ProtectCurrentLocale(this_.compiler)
	_ = protectLocaleCallback
	NodeTraversal_Traverse(this_.compiler, root, protectLocaleCallback)
}

// Java class: com.google.javascript.jscomp.LocaleDataPasses.ProtectCurrentLocale
type LocaleDataPasses_ProtectCurrentLocale struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler   AbstractCompiler `gson:"compiler"`
	astFactory *AstFactory      `gson:"astFactory"`
}

func (this_ *LocaleDataPasses_ProtectCurrentLocale) BindSelf_(s any) {
	_ = s
}
func (this_ *LocaleDataPasses_ProtectCurrentLocale) IsJavaNil_() bool { return this_ == nil }
func (*LocaleDataPasses_ProtectCurrentLocale) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LocaleDataPasses$ProtectCurrentLocale", "ProtectCurrentLocale"
}
func (this_ *LocaleDataPasses_ProtectCurrentLocale) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*LocaleDataPasses_ProtectCurrentLocale) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LocaleDataPasses.ProtectCurrentLocale", "com.google.javascript.jscomp.LocaleDataPasses$ProtectCurrentLocale", func(v_ any) bool {
		return j2gort.InstanceOf[*LocaleDataPasses_ProtectCurrentLocale](v_)
	})
}
func (*LocaleDataPasses_ProtectCurrentLocale) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*LocaleDataPasses_ProtectCurrentLocale) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *LocaleDataPasses_ProtectCurrentLocale) initFields_LocaleDataPasses_ProtectCurrentLocale() {
	_ = this_
}

// Java constructor: ProtectCurrentLocale(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewLocaleDataPasses_ProtectCurrentLocale(compiler AbstractCompiler) *LocaleDataPasses_ProtectCurrentLocale {
	this_ := &LocaleDataPasses_ProtectCurrentLocale{}
	InitLocaleDataPasses_ProtectCurrentLocale(this_, compiler)
	return this_
}

func InitLocaleDataPasses_ProtectCurrentLocale(this_ *LocaleDataPasses_ProtectCurrentLocale, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_LocaleDataPasses_ProtectCurrentLocale()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: ProtectCurrentLocale.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *LocaleDataPasses_ProtectCurrentLocale) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if LocaleDataPasses_IsGoogDotLocaleReference(n) {
		var nodeToReplace pkg_rhino.Node
		_ = nodeToReplace
		if NodeUtil_IsLhsOfAssign(n) {
			nodeToReplace = n.GetNext()
		} else {
			nodeToReplace = n
		}
		var replacement pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type("__JSC_LOCALE__", AstFactory_S_Type_1_Node(nodeToReplace))
		_ = replacement
		replacement.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		nodeToReplace.ReplaceWith(replacement)
		this_.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java class: com.google.javascript.jscomp.LocaleDataPasses.LocaleSubstitutions
type LocaleDataPasses_LocaleSubstitutions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler       AbstractCompiler `gson:"compiler"`
	astFactory     *AstFactory      `gson:"astFactory"`
	locale         string           `gson:"locale"`
	qnameForLocale pkg_rhino.Node   `gson:"qnameForLocale"`
}

func (this_ *LocaleDataPasses_LocaleSubstitutions) BindSelf_(s any) {
	_ = s
}
func (this_ *LocaleDataPasses_LocaleSubstitutions) IsJavaNil_() bool { return this_ == nil }
func (*LocaleDataPasses_LocaleSubstitutions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LocaleDataPasses$LocaleSubstitutions", "LocaleSubstitutions"
}
func (this_ *LocaleDataPasses_LocaleSubstitutions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*LocaleDataPasses_LocaleSubstitutions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LocaleDataPasses.LocaleSubstitutions", "com.google.javascript.jscomp.LocaleDataPasses$LocaleSubstitutions", func(v_ any) bool {
		return j2gort.InstanceOf[*LocaleDataPasses_LocaleSubstitutions](v_)
	})
}
func (*LocaleDataPasses_LocaleSubstitutions) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*LocaleDataPasses_LocaleSubstitutions) Implements_jscomp_CompilerPass()           {}
func (*LocaleDataPasses_LocaleSubstitutions) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_LocaleDataPasses_LocaleSubstitutions_DEFAULT_LOCALE() string { return "en" }
func (this_ *LocaleDataPasses_LocaleSubstitutions) initFields_LocaleDataPasses_LocaleSubstitutions() {
	_ = this_
	this_.qnameForLocale = pkg_rhino.IR_Name("__JSC_LOCALE__")
	this_.locale = j2gort.NullString
}

// Java constructor: LocaleSubstitutions(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String locale)
func NewLocaleDataPasses_LocaleSubstitutions(compiler AbstractCompiler, locale string) *LocaleDataPasses_LocaleSubstitutions {
	this_ := &LocaleDataPasses_LocaleSubstitutions{}
	InitLocaleDataPasses_LocaleSubstitutions(this_, compiler, locale)
	return this_
}

func InitLocaleDataPasses_LocaleSubstitutions(this_ *LocaleDataPasses_LocaleSubstitutions, compiler AbstractCompiler, locale string) {
	_ = compiler
	_ = locale
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_LocaleDataPasses_LocaleSubstitutions()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	if locale == j2gort.NullString {
		this_.locale = "en"
	} else {
		this_.locale = locale
	}
}

// Java method: LocaleSubstitutions.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *LocaleDataPasses_LocaleSubstitutions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: LocaleSubstitutions.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *LocaleDataPasses_LocaleSubstitutions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.MatchesName_1_Node(this_.qnameForLocale) {
		var replacement pkg_rhino.Node = this_.astFactory.CreateString(j2gort.CheckNotNullString(this_.locale)).Srcref(n)
		_ = replacement
		n.ReplaceWith(replacement)
		this_.compiler.ReportChangeToEnclosingScope(replacement)
	}
}
func (*LocaleDataPasses_ProtectCurrentLocale) AlwaysShouldTraverse_()      {}
func (*LocaleDataPasses_ProtectCurrentLocale) CanReleaseSyntacticScopes_() {}
func (*LocaleDataPasses_LocaleSubstitutions) AlwaysShouldTraverse_()       {}
func (*LocaleDataPasses_LocaleSubstitutions) CanReleaseSyntacticScopes_()  {}
