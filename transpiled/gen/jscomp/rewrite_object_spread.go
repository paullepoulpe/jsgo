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
// Port of com/google/javascript/jscomp/RewriteObjectSpread.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RewriteObjectSpread
type RewriteObjectSpread struct {
	compiler        AbstractCompiler      `gson:"compiler"`
	astFactory      *AstFactory           `gson:"astFactory"`
	namespace       pkg_rhino.StaticScope `gson:"namespace"`
	idHash_         atomic.Int32
	useJscompAssign bool `gson:"useJscompAssign"`
}

func (this_ *RewriteObjectSpread) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteObjectSpread) IsJavaNil_() bool { return this_ == nil }
func (*RewriteObjectSpread) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteObjectSpread", "RewriteObjectSpread"
}
func (this_ *RewriteObjectSpread) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteObjectSpread) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteObjectSpread", "com.google.javascript.jscomp.RewriteObjectSpread", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteObjectSpread](v_)
	})
}
func (*RewriteObjectSpread) Implements_jscomp_NodeTraversal_Callback() {}
func (*RewriteObjectSpread) Implements_jscomp_CompilerPass()           {}

var static_RewriteObjectSpread_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_RewriteObjectSpread_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_RewriteObjectSpread()
	return static_RewriteObjectSpread_transpiledFeatures
}
func StaticSet_RewriteObjectSpread_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_RewriteObjectSpread()
	static_RewriteObjectSpread_transpiledFeatures = v
}

var clinitDone_RewriteObjectSpread uint32
var clinitState_RewriteObjectSpread j2gort.ClinitState

func clinit_RewriteObjectSpread() {
	if atomic.LoadUint32(&clinitDone_RewriteObjectSpread) != 0 {
		return
	}
	clinitSlow_RewriteObjectSpread()
}

func clinitSlow_RewriteObjectSpread() {
	if !j2gort.EnterClinit(&clinitState_RewriteObjectSpread, &clinitDone_RewriteObjectSpread) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteObjectSpread, &clinitDone_RewriteObjectSpread)
	static_RewriteObjectSpread_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD())
}

func (this_ *RewriteObjectSpread) initFields_RewriteObjectSpread() {
	_ = this_
}

// Java constructor: RewriteObjectSpread(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteObjectSpread(compiler AbstractCompiler) *RewriteObjectSpread {
	clinit_RewriteObjectSpread()
	this_ := &RewriteObjectSpread{}
	InitRewriteObjectSpread(this_, compiler)
	return this_
}

func InitRewriteObjectSpread(this_ *RewriteObjectSpread, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RewriteObjectSpread()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.namespace = this_.compiler.GetTranspilationNamespace()
	this_.useJscompAssign = !(this_.compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015()))
}

// Java method: RewriteObjectSpread.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteObjectSpread) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_RewriteObjectSpread_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_RewriteObjectSpread_TranspiledFeatures())
}

// Java method: RewriteObjectSpread.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteObjectSpread) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}

// Java method: RewriteObjectSpread.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteObjectSpread) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(47) /* Token.OBJECTLIT */ :
		this_.priv_RewriteObjectSpread_visitObject_5000(n)
	default:
	}
}

// Java method: RewriteObjectSpread.visitObject(com.google.javascript.rhino.Node obj)
func (this_ *RewriteObjectSpread) priv_RewriteObjectSpread_visitObject_5000(obj pkg_rhino.Node) {
	_ = obj
	{
		var child pkg_rhino.Node = obj.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsSpread() {
				this_.priv_RewriteObjectSpread_visitObjectWithSpread_5001(obj)
				return
			}
		}
	}
}

// Java method: RewriteObjectSpread.visitObjectWithSpread(com.google.javascript.rhino.Node obj)
func (this_ *RewriteObjectSpread) priv_RewriteObjectSpread_visitObjectWithSpread_5001(obj pkg_rhino.Node) {
	_ = obj
	if !(obj.IsObjectLit()) {
		j2gort.CheckArgument(false)
	}
	var result pkg_rhino.Node
	_ = result
	if this_.useJscompAssign {
		result = this_.astFactory.CreateJscompAssignCall(this_.namespace, AstFactory_S_Type_1_Node(obj), []pkg_rhino.Node{this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})})
	} else {
		result = this_.astFactory.CreateObjectDotAssignCall(this_.namespace, AstFactory_S_Type_1_Node(obj), []pkg_rhino.Node{this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})})
	}
	var trailingObjectLiteral pkg_rhino.Node = pkg_rhino.Node(0)
	_ = trailingObjectLiteral
	{
		var child pkg_rhino.Node = obj.GetFirstChild()
		_ = child
		for child != 0 {
			var next pkg_rhino.Node = child.GetNext()
			_ = next
			if child.IsSpread() {
				var spreaded pkg_rhino.Node = child.RemoveFirstChild()
				_ = spreaded
				result.AddChildToBack(spreaded)
				trailingObjectLiteral = pkg_rhino.Node(0)
			} else {
				if trailingObjectLiteral == 0 {
					trailingObjectLiteral = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
					result.AddChildToBack(trailingObjectLiteral)
				}
				trailingObjectLiteral.AddChildToBack(child.Detach())
			}
			child = next
		}
	}
	result.SrcrefTreeIfMissing(obj)
	obj.ReplaceWith(result)
	this_.compiler.ReportChangeToEnclosingScope(result)
}
func (*RewriteObjectSpread) AlwaysShouldTraverse_()      {}
func (*RewriteObjectSpread) CanReleaseSyntacticScopes_() {}
