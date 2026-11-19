/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteOptionalChainingOperator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteOptionalChainingOperator
type RewriteOptionalChainingOperator struct {
	compiler                     AbstractCompiler                       `gson:"compiler"`
	getTmpVarNameCreatorForInput j2gort.Ext_java_util_function_Function `gson:"getTmpVarNameCreatorForInput"`
	idHash_                      atomic.Int32
}

func (this_ *RewriteOptionalChainingOperator) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteOptionalChainingOperator) IsJavaNil_() bool { return this_ == nil }
func (*RewriteOptionalChainingOperator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteOptionalChainingOperator", "RewriteOptionalChainingOperator"
}
func (this_ *RewriteOptionalChainingOperator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteOptionalChainingOperator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteOptionalChainingOperator", "com.google.javascript.jscomp.RewriteOptionalChainingOperator", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteOptionalChainingOperator](v_)
	})
}
func (*RewriteOptionalChainingOperator) Implements_jscomp_CompilerPass() {}

func (this_ *RewriteOptionalChainingOperator) initFields_RewriteOptionalChainingOperator() {
	_ = this_
}

// Java constructor: RewriteOptionalChainingOperator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteOptionalChainingOperator_1(compiler AbstractCompiler) *RewriteOptionalChainingOperator {
	this_ := &RewriteOptionalChainingOperator{}
	InitRewriteOptionalChainingOperator_1(this_, compiler)
	return this_
}

func InitRewriteOptionalChainingOperator_1(this_ *RewriteOptionalChainingOperator, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RewriteOptionalChainingOperator()
	this_.compiler = compiler
	var uniqueIdSupplier *UniqueIdSupplier = compiler.GetUniqueIdSupplier()
	_ = uniqueIdSupplier
	this_.getTmpVarNameCreatorForInput = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		input := j2gort.CastPtr[CompilerInput](lp_1)
		_ = input
		return NewLambda_OptionalChainRewriter_TmpVarNameCreator(func() string {
			return ("$jscomp$optchain$tmp" + j2gort.StrValString(uniqueIdSupplier.GetUniqueId(input)))
		})
	})
}

// Java constructor: RewriteOptionalChainingOperator(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.OptionalChainRewriter.TmpVarNameCreator tmpVarNameCreator)
func NewRewriteOptionalChainingOperator_2(compiler AbstractCompiler, tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) *RewriteOptionalChainingOperator {
	this_ := &RewriteOptionalChainingOperator{}
	InitRewriteOptionalChainingOperator_2(this_, compiler, tmpVarNameCreator)
	return this_
}

func InitRewriteOptionalChainingOperator_2(this_ *RewriteOptionalChainingOperator, compiler AbstractCompiler, tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) {
	_ = compiler
	_ = tmpVarNameCreator
	this_.initFields_RewriteOptionalChainingOperator()
	this_.compiler = compiler
	this_.getTmpVarNameCreatorForInput = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		input := j2gort.CastPtr[CompilerInput](lp_2)
		_ = input
		return tmpVarNameCreator
	})
}

// Java method: RewriteOptionalChainingOperator.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteOptionalChainingOperator) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, NewRewriteOptionalChainingOperator_TranspilationCallback(this_), externs, root)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(this_.compiler, root, pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING())
}

// Java class: com.google.javascript.jscomp.RewriteOptionalChainingOperator.TranspilationCallback
type RewriteOptionalChainingOperator_TranspilationCallback struct {
	outer_          *RewriteOptionalChainingOperator
	rewriterBuilder *OptionalChainRewriter_Builder `gson:"rewriterBuilder"`
	optChains       []*OptionalChainRewriter       `gson:"optChains"`
	idHash_         atomic.Int32
}

func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) IsJavaNil_() bool {
	return this_ == nil
}
func (*RewriteOptionalChainingOperator_TranspilationCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteOptionalChainingOperator$TranspilationCallback", "TranspilationCallback"
}
func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteOptionalChainingOperator_TranspilationCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteOptionalChainingOperator.TranspilationCallback", "com.google.javascript.jscomp.RewriteOptionalChainingOperator$TranspilationCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteOptionalChainingOperator_TranspilationCallback](v_)
	})
}
func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) SetOuter_(o *RewriteOptionalChainingOperator) {
	this_.outer_ = o
}
func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) GetOuter_() *RewriteOptionalChainingOperator {
	return this_.outer_
}
func (*RewriteOptionalChainingOperator_TranspilationCallback) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) initFields_RewriteOptionalChainingOperator_TranspilationCallback() {
	_ = this_
	this_.rewriterBuilder = OptionalChainRewriter_S_Builder(this_.outer_.compiler)
	this_.optChains = []*OptionalChainRewriter{}
}

// Java constructor: TranspilationCallback()
func NewRewriteOptionalChainingOperator_TranspilationCallback(outer_ *RewriteOptionalChainingOperator) *RewriteOptionalChainingOperator_TranspilationCallback {
	this_ := &RewriteOptionalChainingOperator_TranspilationCallback{}
	this_.outer_ = outer_
	InitRewriteOptionalChainingOperator_TranspilationCallback(this_)
	return this_
}

func InitRewriteOptionalChainingOperator_TranspilationCallback(this_ *RewriteOptionalChainingOperator_TranspilationCallback) {
	this_.initFields_RewriteOptionalChainingOperator_TranspilationCallback()
}

// Java method: TranspilationCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.rewriterBuilder.SetTmpVarNameCreator(j2gort.AsIface[OptionalChainRewriter_TmpVarNameCreator](this_.outer_.getTmpVarNameCreatorForInput.Apply(j2gort.PtrToAny(t.GetInput()))))
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return ((scriptFeatures == nil) || scriptFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING()))
	}
	return true
}

// Java method: TranspilationCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteOptionalChainingOperator_TranspilationCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsEndOfFullOptChain(n) {
		this_.optChains = append(this_.optChains, this_.rewriterBuilder.Build(n))
	} else if n.IsScript() {
		if len(this_.optChains) != 0 {
			for _, optChain := range this_.optChains {
				_ = optChain
				optChain.Rewrite()
			}

			clear(this_.optChains)
			this_.optChains = this_.optChains[:0]
		}
	}
}
func (*RewriteOptionalChainingOperator_TranspilationCallback) CanReleaseSyntacticScopes_() {}
