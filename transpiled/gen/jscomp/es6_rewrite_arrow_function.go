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
// Port of com/google/javascript/jscomp/Es6RewriteArrowFunction.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RewriteArrowFunction
type Es6RewriteArrowFunction struct {
	compiler     AbstractCompiler                                                        `gson:"compiler"`
	astFactory   *AstFactory                                                             `gson:"astFactory"`
	contextStack j2gort.Deque[*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext] `gson:"contextStack"`
	idHash_      atomic.Int32
}

func (this_ *Es6RewriteArrowFunction) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteArrowFunction) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteArrowFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteArrowFunction", "Es6RewriteArrowFunction"
}
func (this_ *Es6RewriteArrowFunction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteArrowFunction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteArrowFunction", "com.google.javascript.jscomp.Es6RewriteArrowFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteArrowFunction](v_)
	})
}
func (*Es6RewriteArrowFunction) Implements_jscomp_NodeTraversal_Callback() {}
func (*Es6RewriteArrowFunction) Implements_jscomp_CompilerPass()           {}

var static_Es6RewriteArrowFunction_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_Es6RewriteArrowFunction_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_Es6RewriteArrowFunction()
	return static_Es6RewriteArrowFunction_transpiledFeatures
}
func StaticSet_Es6RewriteArrowFunction_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteArrowFunction()
	static_Es6RewriteArrowFunction_transpiledFeatures = v
}

var clinitDone_Es6RewriteArrowFunction uint32
var clinitState_Es6RewriteArrowFunction j2gort.ClinitState

func clinit_Es6RewriteArrowFunction() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteArrowFunction) != 0 {
		return
	}
	clinitSlow_Es6RewriteArrowFunction()
}

func clinitSlow_Es6RewriteArrowFunction() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteArrowFunction, &clinitDone_Es6RewriteArrowFunction) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteArrowFunction, &clinitDone_Es6RewriteArrowFunction)
	static_Es6RewriteArrowFunction_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS())
}

func (this_ *Es6RewriteArrowFunction) initFields_Es6RewriteArrowFunction() {
	_ = this_
}

// Java constructor: Es6RewriteArrowFunction(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteArrowFunction(compiler AbstractCompiler) *Es6RewriteArrowFunction {
	clinit_Es6RewriteArrowFunction()
	this_ := &Es6RewriteArrowFunction{}
	InitEs6RewriteArrowFunction(this_, compiler)
	return this_
}

func InitEs6RewriteArrowFunction(this_ *Es6RewriteArrowFunction, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6RewriteArrowFunction()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.contextStack = j2gort.Deque[*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext]{}
}

// Java method: Es6RewriteArrowFunction.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteArrowFunction) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6RewriteArrowFunction_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteArrowFunction_TranspiledFeatures())
}

// Java method: Es6RewriteArrowFunction.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteArrowFunction) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		this_.contextStack.Push(this_.priv_Es6RewriteArrowFunction_contextForScript_2581(n, t.GetInput()))
	case int32(73) /* Token.FUNCTION */ :
		if !(n.IsArrowFunction()) {
			this_.contextStack.Push(this_.priv_Es6RewriteArrowFunction_contextForFunction_2580(n, t.GetInput()))
		}
	case int32(110) /* Token.SUPER */ :
		var context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext = j2gort.CheckNotNullPtr(this_.contextStack.Peek())
		_ = context
		if (context.isConstructor && parent.IsCall()) && (parent.GetFirstChild() == n) {
			{
				r_1 := context
				r_1.lastSuperStatement = this_.priv_Es6RewriteArrowFunction_getEnclosingStatement_2578(parent, context.scopeBody)
			}
		}
	default:
	}
	return true
}

// Java method: Es6RewriteArrowFunction.getEnclosingStatement(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node block)
func (this_ *Es6RewriteArrowFunction) priv_Es6RewriteArrowFunction_getEnclosingStatement_2578(n pkg_rhino.Node, block pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	_ = block
	for j2gort.CheckNotNullNode(n.GetParent()) != block {
		n = j2gort.CheckNotNullNode(NodeUtil_GetEnclosingStatement(n.GetParent()))
	}
	return n
}

// Java method: Es6RewriteArrowFunction.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteArrowFunction) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext = this_.contextStack.Peek()
	_ = context
	if n.IsArrowFunction() {
		this_.priv_Es6RewriteArrowFunction_visitArrowFunction_2579(t, n, j2gort.CheckNotNullPtr(context))
	} else if (context != nil) && (context.scopeBody == n) {
		this_.contextStack.Pop()
		context.AddVarDeclarations(this_.compiler, this_.astFactory, t.GetCurrentScript())
	}
}

// Java method: Es6RewriteArrowFunction.visitArrowFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater.ThisAndArgumentsContext context)
func (this_ *Es6RewriteArrowFunction) priv_Es6RewriteArrowFunction_visitArrowFunction_2579(t *NodeTraversal, n pkg_rhino.Node, context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) {
	_ = t
	_ = n
	_ = context
	n.SetIsArrowFunction(false)
	var body pkg_rhino.Node = n.GetLastChild()
	_ = body
	if !(body.IsBlock()) {
		j2gort.CheckStateStr(false, "Arrow function body must be a block after normalization")
	}
	var val_2 ThisAndArgumentsReferenceUpdater
	InitThisAndArgumentsReferenceUpdater(&val_2, this_.compiler, context, this_.astFactory)
	var updater *ThisAndArgumentsReferenceUpdater = &val_2
	_ = updater
	NodeTraversal_Traverse(this_.compiler, body, updater)
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteArrowFunction.contextForFunction(com.google.javascript.rhino.Node functionNode, com.google.javascript.jscomp.CompilerInput input)
func (this_ *Es6RewriteArrowFunction) priv_Es6RewriteArrowFunction_contextForFunction_2580(functionNode pkg_rhino.Node, input *CompilerInput) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	_ = functionNode
	_ = input
	var scopeBody pkg_rhino.Node = functionNode.GetLastChild()
	_ = scopeBody
	return NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(scopeBody, NodeUtil_IsEs6Constructor(functionNode), this_.compiler.GetUniqueIdSupplier().GetUniqueId(input))
}

// Java method: Es6RewriteArrowFunction.contextForScript(com.google.javascript.rhino.Node scriptNode, com.google.javascript.jscomp.CompilerInput input)
func (this_ *Es6RewriteArrowFunction) priv_Es6RewriteArrowFunction_contextForScript_2581(scriptNode pkg_rhino.Node, input *CompilerInput) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	_ = scriptNode
	_ = input
	return NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(scriptNode, false, this_.compiler.GetUniqueIdSupplier().GetUniqueId(input))
}
func (*Es6RewriteArrowFunction) CanReleaseSyntacticScopes_() {}
