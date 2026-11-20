/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/NodeTraversal.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.jscomp.NodeTraversal
type NodeTraversal struct {
	compiler                                       AbstractCompiler             `gson:"compiler"`
	callback                                       NodeTraversal_Callback       `gson:"callback"`
	scopeCallback                                  NodeTraversal_ScopedCallback `gson:"scopeCallback"`
	scopeCreator                                   ScopeCreator                 `gson:"scopeCreator"`
	scopeNodes                                     []pkg_rhino.Node             `gson:"scopeNodes"`
	scopeObjs                                      []AbstractScope              `gson:"scopeObjs"`
	sourceName                                     string                       `gson:"sourceName"`
	inputId                                        *pkg_rhino.InputId           `gson:"inputId"`
	compilerInput                                  *CompilerInput               `gson:"compilerInput"`
	currentNode                                    pkg_rhino.Node               `gson:"currentNode"`
	currentHoistScopeRoot                          pkg_rhino.Node               `gson:"currentHoistScopeRoot"`
	currentFunction                                pkg_rhino.Node               `gson:"currentFunction"`
	currentScript                                  pkg_rhino.Node               `gson:"currentScript"`
	currentChangeScope                             pkg_rhino.Node               `gson:"currentChangeScope"`
	scopeDepth                                     int32                        `gson:"scopeDepth"`
	idHash_                                        atomic.Int32
	obeyDestructuringAndDefaultValueExecutionOrder bool `gson:"obeyDestructuringAndDefaultValueExecutionOrder"`
	mayContainSyntheticBlocks                      bool `gson:"mayContainSyntheticBlocks"`
	scopeNodesBuf_                                 [12]pkg_rhino.Node
	scopeObjsBuf_                                  [12]AbstractScope
	defaultScopeCreator_                           SyntacticScopeCreator
	alwaysTraverse_                                bool
	noOpVisit_                                     bool
	canReleaseScopes_                              bool
}

func (this_ *NodeTraversal) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal", "NodeTraversal"
}
func (this_ *NodeTraversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeTraversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal", "com.google.javascript.jscomp.NodeTraversal", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeTraversal](v_)
	})
}

var static_NodeTraversal_EMPTY_SCOPE_NODES []pkg_rhino.Node

func StaticGet_NodeTraversal_EMPTY_SCOPE_NODES() []pkg_rhino.Node {
	clinit_NodeTraversal()
	return static_NodeTraversal_EMPTY_SCOPE_NODES
}
func StaticSet_NodeTraversal_EMPTY_SCOPE_NODES(v []pkg_rhino.Node) {
	clinit_NodeTraversal()
	static_NodeTraversal_EMPTY_SCOPE_NODES = v
}

var static_NodeTraversal_EMPTY_SCOPE_OBJS []AbstractScope

func StaticGet_NodeTraversal_EMPTY_SCOPE_OBJS() []AbstractScope {
	clinit_NodeTraversal()
	return static_NodeTraversal_EMPTY_SCOPE_OBJS
}
func StaticSet_NodeTraversal_EMPTY_SCOPE_OBJS(v []AbstractScope) {
	clinit_NodeTraversal()
	static_NodeTraversal_EMPTY_SCOPE_OBJS = v
}
func StaticGet_NodeTraversal_MISSING_SOURCE() string { return "[source unknown]" }

var clinitDone_NodeTraversal uint32
var clinitState_NodeTraversal j2gort.ClinitState

func clinit_NodeTraversal() {
	if atomic.LoadUint32(&clinitDone_NodeTraversal) != 0 {
		return
	}
	clinitSlow_NodeTraversal()
}

func clinitSlow_NodeTraversal() {
	if !j2gort.EnterClinit(&clinitState_NodeTraversal, &clinitDone_NodeTraversal) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_NodeTraversal, &clinitDone_NodeTraversal)
	static_NodeTraversal_EMPTY_SCOPE_NODES = make([]pkg_rhino.Node, int32(0))
	static_NodeTraversal_EMPTY_SCOPE_OBJS = make([]AbstractScope, int32(0))
}

func (this_ *NodeTraversal) initFields_NodeTraversal() {
	_ = this_
	this_.scopeNodes = StaticGet_NodeTraversal_EMPTY_SCOPE_NODES()
	this_.scopeObjs = StaticGet_NodeTraversal_EMPTY_SCOPE_OBJS()
	this_.scopeDepth = int32(0)
	this_.sourceName = j2gort.NullString
}

// Java constructor: NodeTraversal(com.google.javascript.jscomp.NodeTraversal.Builder builder)
func NewNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	clinit_NodeTraversal()
	this_ := &NodeTraversal{}
	InitNodeTraversal(this_, builder)
	return this_
}

func InitNodeTraversal(this_ *NodeTraversal, builder *NodeTraversal_Builder) {
	_ = builder
	this_.initFields_NodeTraversal()
	this_.compiler = j2gort.CheckNotNullIface(builder.compiler)
	this_.callback = j2gort.CheckNotNullIface(builder.callback)
	var scopedCallback NodeTraversal_ScopedCallback
	_ = scopedCallback
	var ok_1 bool
	scopedCallback, ok_1 = j2gort.TryCast[NodeTraversal_ScopedCallback](this_.callback)
	if ok_1 {
		this_.scopeCallback = scopedCallback
	} else {
		this_.scopeCallback = nil
	}
	if builder.scopeCreator == nil || builder.scopeCreator.IsJavaNil_() {
		if this_.defaultScopeCreator_.reusableScanner == nil {
			InitSyntacticScopeCreator_1(&this_.defaultScopeCreator_, this_.compiler)
		} else {
			handler_ := StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER()
			this_.defaultScopeCreator_.compiler = this_.compiler
			this_.defaultScopeCreator_.redeclarationHandler = handler_
			this_.defaultScopeCreator_.treatProvidesAsRedeclarations = false
			this_.defaultScopeCreator_.scannerInUse = false
			this_.defaultScopeCreator_.reusableScanner.compiler = this_.compiler
			this_.defaultScopeCreator_.reusableScanner.redeclarationHandler = handler_
			this_.defaultScopeCreator_.reusableScanner.treatProvidesAsRedeclarations = false
			this_.defaultScopeCreator_.reusableScanner.Reset()
		}
		this_.scopeCreator = &this_.defaultScopeCreator_
	} else {
		this_.scopeCreator = builder.scopeCreator
	}
	this_.obeyDestructuringAndDefaultValueExecutionOrder = builder.obeyDestructuringAndDefaultValueExecutionOrder
	var options *CompilerOptions = this_.compiler.GetOptions()
	_ = options
	this_.mayContainSyntheticBlocks = ((options == nil) || (options.GetSyntheticBlockStartMarker() != j2gort.NullString))
	cbTarget_ := unwrapNodeTraversalCallback_(this_.callback)
	_, this_.alwaysTraverse_ = any(this_.callback).(interface{ AlwaysShouldTraverse_() })
	if !this_.alwaysTraverse_ {
		_, this_.alwaysTraverse_ = cbTarget_.(interface{ AlwaysShouldTraverse_() })
	}
	_, this_.noOpVisit_ = any(this_.callback).(interface{ NoOpVisit_() })
	if !this_.noOpVisit_ {
		_, this_.noOpVisit_ = cbTarget_.(interface{ NoOpVisit_() })
	}
	this_.canReleaseScopes_ = canReleaseSyntacticScopesForCallback_(this_.callback)

}

// Java static method: NodeTraversal.builder()
func NodeTraversal_S_Builder() *NodeTraversal_Builder {
	clinit_NodeTraversal()
	return NewNodeTraversal_Builder()
}

// Java method: NodeTraversal.throwUnexpectedException(java.lang.Throwable unexpectedException)
func (this_ *NodeTraversal) priv_NodeTraversal_throwUnexpectedException_3975(unexpectedException javalang.Throwable) {
	_ = unexpectedException
	var message string = j2gort.ThrowableGetMessage(unexpectedException)
	_ = message
	if this_.currentScript != 0 {
		message = (j2gort.StrValString(j2gort.ThrowableGetMessage(unexpectedException)) + "\n" + j2gort.StrValString(this_.priv_NodeTraversal_formatNodeContext_3976("Node", this_.currentNode)) + j2gort.StrValString(func() string {
			if this_.currentNode == 0 {
				return ""
			}
			return this_.priv_NodeTraversal_formatNodeContext_3976("Parent", this_.currentNode.GetParent())
		}()))
	}
	this_.compiler.ThrowInternalError(message, unexpectedException)
}

// Java method: NodeTraversal.formatNodeContext(java.lang.String label, com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_formatNodeContext_3976(label string, n pkg_rhino.Node) string {
	_ = label
	_ = n
	if n == 0 {
		return ("  " + j2gort.StrValString(label) + ": NULL")
	}
	return ("  " + j2gort.StrValString(label) + "(" + j2gort.StrValString(n.ToString_3(false, false, false)) + "): " + j2gort.StrValString(this_.priv_NodeTraversal_formatNodePosition_3979(n)))
}

// Java method: NodeTraversal.traverse(com.google.javascript.rhino.Node root)
func (this_ *NodeTraversal) priv_NodeTraversal_traverse_3977(root pkg_rhino.Node) {
	_ = root
	func() {
		defer func() {
			if rec_2 := recover(); rec_2 != nil {
				jex_3 := javalang.RecoverException(rec_2)
				_ = jex_3
				if j2gort.InstanceOf[javalang.IError](jex_3) || j2gort.InstanceOf[javalang.IException](jex_3) {
					unexpectedException := jex_3
					_ = unexpectedException
					this_.priv_NodeTraversal_throwUnexpectedException_3975(unexpectedException)
				} else {
					panic(rec_2)
				}
			}
		}()
		this_.priv_NodeTraversal_initTraversal_4031(root)
		this_.currentNode = root
		this_.priv_NodeTraversal_pushScope_4001(root)
		this_.priv_NodeTraversal_traverseBranch_3993(root, pkg_rhino.Node(0))
		this_.priv_NodeTraversal_popScope_4004()
	}()
}

// Java static method: NodeTraversal.traverse(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.NodeTraversal.Callback cb)
func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	t_ := acquireNodeTraversal_(compiler, cb, nil, false)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverse_3977(root)
}

// Java method: NodeTraversal.traverseRoots(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseRoots_3978(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	func() {
		defer func() {
			if rec_4 := recover(); rec_4 != nil {
				jex_5 := javalang.RecoverException(rec_4)
				_ = jex_5
				if j2gort.InstanceOf[javalang.IError](jex_5) || j2gort.InstanceOf[javalang.IException](jex_5) {
					unexpectedException := jex_5
					_ = unexpectedException
					this_.priv_NodeTraversal_throwUnexpectedException_3975(unexpectedException)
				} else {
					panic(rec_4)
				}
			}
		}()
		var scopeRoot pkg_rhino.Node = externs.GetParent()
		_ = scopeRoot
		j2gort.CheckNotNullNode(scopeRoot)
		this_.priv_NodeTraversal_initTraversal_4031(scopeRoot)
		this_.currentNode = scopeRoot
		this_.priv_NodeTraversal_pushScope_4001(scopeRoot)
		this_.priv_NodeTraversal_traverseBranch_3993(externs, scopeRoot)
		if !(root.GetParent() == scopeRoot) {
			j2gort.CheckState(false)
		}
		this_.priv_NodeTraversal_traverseBranch_3993(root, scopeRoot)
		this_.priv_NodeTraversal_popScope_4004()
	}()
}

// Java static method: NodeTraversal.traverseRoots(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal.Callback cb, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func NodeTraversal_TraverseRoots(compiler AbstractCompiler, cb NodeTraversal_Callback, externs pkg_rhino.Node, root pkg_rhino.Node) {
	t_ := acquireNodeTraversal_(compiler, cb, nil, false)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverseRoots_3978(externs, root)
}

// Java method: NodeTraversal.formatNodePosition(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_formatNodePosition_3979(n pkg_rhino.Node) string {
	_ = n
	var sourceFileName string = this_.priv_NodeTraversal_getBestSourceFileName_4036(n)
	_ = sourceFileName
	if sourceFileName == j2gort.NullString {
		return ("[source unknown]" + "\n")
	}
	var lineNumber int32 = n.GetLineno()
	_ = lineNumber
	var columnNumber int32 = n.GetCharno()
	_ = columnNumber
	var src string = this_.compiler.GetSourceLine(sourceFileName, lineNumber)
	_ = src
	if src == j2gort.NullString {
		src = "[source unknown]"
	}
	return (j2gort.StrValString(sourceFileName) + ":" + j2gort.StrValInt32(lineNumber) + ":" + j2gort.StrValInt32(columnNumber) + "\n" + j2gort.StrValString(src) + "\n")
}

// Java method: NodeTraversal.traverseWithScope(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.AbstractScope<?,?> s)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseWithScope_3980(root pkg_rhino.Node, s AbstractScope) {
	_ = root
	_ = s
	if !(s.IsGlobal() || s.GetAbstractScopeBase().IsModuleScope()) {
		j2gort.CheckStateMsg(false, s)
	}
	func() {
		defer func() {
			if rec_6 := recover(); rec_6 != nil {
				jex_7 := javalang.RecoverException(rec_6)
				_ = jex_7
				if j2gort.InstanceOf[javalang.IError](jex_7) || j2gort.InstanceOf[javalang.IException](jex_7) {
					unexpectedException := jex_7
					_ = unexpectedException
					this_.priv_NodeTraversal_throwUnexpectedException_3975(unexpectedException)
				} else {
					panic(rec_6)
				}
			}
		}()
		this_.priv_NodeTraversal_initTraversal_4031(root)
		this_.currentNode = root
		this_.priv_NodeTraversal_pushScope_4002(s)
		this_.priv_NodeTraversal_traverseBranch_3993(root, pkg_rhino.Node(0))
		this_.priv_NodeTraversal_popScope_4004()
	}()
}

// Java method: NodeTraversal.traverseAtScope(com.google.javascript.jscomp.AbstractScope<?,?> s)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseAtScope_3981(s AbstractScope) {
	_ = s
	base_s_ := s.GetAbstractScopeBase()
	var n pkg_rhino.Node = base_s_.GetRootNode_0()
	_ = n
	this_.priv_NodeTraversal_initTraversal_4031(n)
	this_.currentNode = n
	var parentScopes j2gort.Deque[AbstractScope]
	_ = parentScopes
	var temp AbstractScope = j2gort.AsIface[AbstractScope](s.GetParent())
	_ = temp
	for temp != nil && !temp.IsJavaNil_() {
		parentScopes.Push(temp)
		temp = j2gort.AsIface[AbstractScope](temp.GetParent())
	}
	for !(parentScopes.IsEmpty()) {
		this_.priv_NodeTraversal_pushScope_4003(parentScopes.Pop(), true)
	}
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		if this_.callback.ShouldTraverse_3(this_, n, pkg_rhino.Node(0)) {
			this_.priv_NodeTraversal_pushScope_4002(s)
			var fnName pkg_rhino.Node = n.GetFirstChild()
			_ = fnName
			var args pkg_rhino.Node = fnName.GetNext()
			_ = args
			var body pkg_rhino.Node = args.GetNext()
			_ = body
			if !(NodeUtil_IsFunctionDeclaration(n)) {
				this_.priv_NodeTraversal_traverseBranch_3993(fnName, n)
			}
			this_.priv_NodeTraversal_traverseBranch_3993(args, n)
			this_.priv_NodeTraversal_traverseBranch_3993(body, n)
			this_.priv_NodeTraversal_popScope_4004()
			this_.callback.Visit_3(this_, n, pkg_rhino.Node(0))
		}
	case int32(105) /* Token.CLASS */ :
		if this_.callback.ShouldTraverse_3(this_, n, pkg_rhino.Node(0)) {
			this_.priv_NodeTraversal_pushScope_4002(s)
			var className pkg_rhino.Node = n.GetFirstChild()
			_ = className
			var body_2 pkg_rhino.Node = n.GetLastChild()
			_ = body_2
			if NodeUtil_IsClassExpression(n) {
				this_.priv_NodeTraversal_traverseBranch_3993(className, n)
			}
			this_.priv_NodeTraversal_traverseBranch_3993(body_2, n)
			this_.priv_NodeTraversal_popScope_4004()
			this_.callback.Visit_3(this_, n, pkg_rhino.Node(0))
		}
	case int32(90) /* Token.BLOCK */, int32(91) /* Token.SWITCH_BODY */ :
		if this_.callback.ShouldTraverse_3(this_, n, pkg_rhino.Node(0)) {
			this_.priv_NodeTraversal_pushScope_4002(s)
			this_.priv_NodeTraversal_traverseChildren_3998(n)
			this_.priv_NodeTraversal_popScope_4004()
			this_.callback.Visit_3(this_, n, pkg_rhino.Node(0))
		}
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		this_.priv_NodeTraversal_pushScope_4002(s)
		if this_.callback.ShouldTraverse_3(this_, n, pkg_rhino.Node(0)) {
			this_.priv_NodeTraversal_traverseChildren_3998(n)
			this_.callback.Visit_3(this_, n, pkg_rhino.Node(0))
		}
		this_.priv_NodeTraversal_popScope_4004()
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		this_.priv_NodeTraversal_pushScope_4002(s)
		if this_.callback.ShouldTraverse_3(this_, n, pkg_rhino.Node(0)) {
			this_.priv_NodeTraversal_traverseBranch_3993(n.GetLastChild(), n)
			this_.callback.Visit_3(this_, n, pkg_rhino.Node(0))
		}
		this_.priv_NodeTraversal_popScope_4004()
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(80) /* Token.FOR */ :
		if this_.callback.ShouldTraverse_3(this_, n, pkg_rhino.Node(0)) {
			this_.priv_NodeTraversal_pushScope_4002(s)
			var forAssignmentParam pkg_rhino.Node = n.GetFirstChild()
			_ = forAssignmentParam
			var forIterableParam pkg_rhino.Node = forAssignmentParam.GetNext()
			_ = forIterableParam
			var forBodyScope pkg_rhino.Node = forIterableParam.GetNext()
			_ = forBodyScope
			this_.priv_NodeTraversal_traverseBranch_3993(forAssignmentParam, n)
			this_.priv_NodeTraversal_traverseBranch_3993(forIterableParam, n)
			this_.priv_NodeTraversal_traverseBranch_3993(forBodyScope, n)
			this_.priv_NodeTraversal_popScope_4004()
			this_.callback.Visit_3(this_, n, pkg_rhino.Node(0))
		}
	default:
		if !(s.IsGlobal() || base_s_.IsModuleScope()) {
			j2gort.CheckStateFmt(false, "Expected global or module scope. Got: (%s)", s)
		}
		this_.priv_NodeTraversal_traverseWithScope_3980(n, s)
	}
}

// Java static method: NodeTraversal.traverseAtScope(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal.Callback cb, com.google.javascript.jscomp.AbstractScope<?,?> scope)
func NodeTraversal_TraverseAtScope(compiler AbstractCompiler, cb NodeTraversal_Callback, scope AbstractScope) {
	t_ := acquireNodeTraversal_(compiler, cb, nil, false)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverseAtScope_3981(scope)
}

// Java static method: NodeTraversal.traverseScopeRoots(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.rhino.Node> scopeNodes, com.google.javascript.jscomp.NodeTraversal.Callback cb, boolean traverseNested)
func NodeTraversal_TraverseScopeRoots(compiler AbstractCompiler, scopeNodes j2gort.Ext_java_util_List, cb NodeTraversal_Callback, traverseNested bool) {
	clinit_NodeTraversal()
	_ = compiler
	_ = scopeNodes
	_ = cb
	_ = traverseNested
	var scb *NodeTraversal__Local_65_TraverseScopeRootsCallback = NewNodeTraversal__Local_65_TraverseScopeRootsCallback(traverseNested, cb)
	_ = scb
	var val_8 MemoizedScopeCreator
	InitMemoizedScopeCreator_2(&val_8, NewSyntacticScopeCreator_1(compiler), max(int32(8), (scopeNodes.Size()*int32(2))))
	var scopeCreator *MemoizedScopeCreator = &val_8
	_ = scopeCreator
	canRel_ := canReleaseSyntacticScopesForCallback_(cb)
	if canRel_ {
		defer scopeCreator.ReleaseSlabs_()
	}

	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(scopeNodes))
		for cur_9.HasNext() {
			scopeNode := pkg_rhino.Node(cur_9.NextNode())
			_ = scopeNode
			{
				r_10 := scb
				r_10.scopeNode = scopeNode
			}
			traverseScopeRootPooled_(compiler, scb, scopeCreator, scopeNode, canRel_)
		}
	}
}

// Java method: NodeTraversal.traverseScopeRoot(com.google.javascript.rhino.Node scopeRoot)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseScopeRoot_3982(scopeRoot pkg_rhino.Node) {
	_ = scopeRoot
	func() {
		defer func() {
			if rec_11 := recover(); rec_11 != nil {
				jex_12 := javalang.RecoverException(rec_11)
				_ = jex_12
				if j2gort.InstanceOf[javalang.IError](jex_12) || j2gort.InstanceOf[javalang.IException](jex_12) {
					unexpectedException := jex_12
					_ = unexpectedException
					this_.priv_NodeTraversal_throwUnexpectedException_3975(unexpectedException)
				} else {
					panic(rec_11)
				}
			}
		}()
		this_.priv_NodeTraversal_initTraversal_4031(scopeRoot)
		this_.currentNode = scopeRoot
		this_.priv_NodeTraversal_initScopeRoots_4032(scopeRoot.GetParent())
		this_.priv_NodeTraversal_traverseBranch_3993(scopeRoot, scopeRoot.GetParent())
	}()
}

// Java method: NodeTraversal.getCompiler()
func (this_ *NodeTraversal) GetCompiler() AbstractCompiler {
	return this_.compiler
}

// Java method: NodeTraversal.getSourceName()
func (this_ *NodeTraversal) GetSourceName() string {
	if this_.sourceName == j2gort.NullString {
		if this_.currentScript != 0 {
			this_.sourceName = this_.currentScript.GetSourceFileName()
		} else {
			this_.sourceName = ""
		}
	}
	return this_.sourceName
}

// Java method: NodeTraversal.getInput()
func (this_ *NodeTraversal) GetInput() *CompilerInput {
	var inputId *pkg_rhino.InputId = this_.GetInputId()
	_ = inputId
	if (this_.compilerInput == nil) && (inputId != nil) {
		this_.compilerInput = this_.compiler.GetInput(inputId)
	}
	return this_.compilerInput
}

// Java method: NodeTraversal.getChunk()
func (this_ *NodeTraversal) GetChunk() *JSChunk {
	var input *CompilerInput = this_.GetInput()
	_ = input
	if input == nil {
		return nil
	} else {
		return input.GetChunk()
	}
}

// Java method: NodeTraversal.getCurrentNode()
func (this_ *NodeTraversal) GetCurrentNode() pkg_rhino.Node {
	return this_.currentNode
}

// Java method: NodeTraversal.handleScript(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_handleScript_3988(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if Platform_IsThreadInterrupted() {
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", j2gort.NewExt_java_lang_InterruptedException())))
	}
	this_.priv_NodeTraversal_setChangeScope_4029(n)
	this_.currentNode = n
	this_.currentScript = n
	this_.priv_NodeTraversal_clearScriptState_4034()
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
	this_.priv_NodeTraversal_setChangeScope_4029(pkg_rhino.Node(0))
}

// Java method: NodeTraversal.handleFunction(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_handleFunction_3989(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	var changeScope pkg_rhino.Node = this_.currentChangeScope
	_ = changeScope
	this_.priv_NodeTraversal_setChangeScope_4029(n)
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseFunction_3994(n, parent)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
	this_.priv_NodeTraversal_setChangeScope_4029(changeScope)
}

// Java method: NodeTraversal.handleModule(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_handleModule_3990(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	this_.currentHoistScopeRoot = n
	this_.priv_NodeTraversal_pushScope_4001(n)
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
	this_.priv_NodeTraversal_popScope_4004()
	this_.currentHoistScopeRoot = pkg_rhino.Node(0)
}

// Java method: NodeTraversal.handleDestructuringOrDefaultValue(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_handleDestructuringOrDefaultValue_3991(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		var first pkg_rhino.Node = n.GetFirstChild()
		_ = first
		var second pkg_rhino.Node = first.GetNext()
		_ = second
		if second != 0 {
			if !(second.GetNext() == 0) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(second))
			}
			this_.priv_NodeTraversal_traverseBranch_3993(second, n)
		}
		this_.priv_NodeTraversal_traverseBranch_3993(first, n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

// Java method: NodeTraversal.createsBlockScope(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_createsBlockScope_3992(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */ :
		if this_.mayContainSyntheticBlocks && n.IsSyntheticBlock() {
			return false
		}
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		return (((parent != 0) && !(NodeUtil_IsSwitchCase(parent))) && !(parent.IsCatch()))
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(91) /* Token.SWITCH_BODY */, int32(105) /* Token.CLASS */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: NodeTraversal.traverseBranch(com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_3993(n pkg_rhino.Node, parent pkg_rhino.Node) {
	curBlkIdx_ := uint32(n) >> nodestore.BlockShift
	curBlk_ := nodestore.BlockOf(nodestore.Node(n))
	f_, _, tok_ := curBlk_.FirstNextTok(nodestore.Node(n))
	this_.traverseBranchBlk_(n, parent, pkg_rhino.Node(f_), tok_, curBlkIdx_, curBlk_)

}

// Java method: NodeTraversal.traverseFunction(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseFunction_3994(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	var fnName pkg_rhino.Node = n.GetFirstChild()
	_ = fnName
	var isFunctionDeclaration bool = ((parent != 0) && NodeUtil_IsFunctionDeclaration(n))
	_ = isFunctionDeclaration
	if isFunctionDeclaration {
		this_.priv_NodeTraversal_traverseBranch_3993(fnName, n)
	}
	this_.currentNode = n
	var previousHoistScopeRoot pkg_rhino.Node = this_.currentHoistScopeRoot
	_ = previousHoistScopeRoot
	this_.currentHoistScopeRoot = n
	var previousFunction pkg_rhino.Node = this_.currentFunction
	_ = previousFunction
	this_.currentFunction = n
	this_.priv_NodeTraversal_pushScope_4001(n)
	if !(isFunctionDeclaration) {
		this_.priv_NodeTraversal_traverseBranch_3993(fnName, n)
	}
	var args pkg_rhino.Node = fnName.GetNext()
	_ = args
	var body pkg_rhino.Node = args.GetNext()
	_ = body
	this_.priv_NodeTraversal_traverseBranch_3993(args, n)
	this_.priv_NodeTraversal_traverseBranch_3993(body, n)
	this_.priv_NodeTraversal_popScope_4004()
	this_.currentFunction = previousFunction
	this_.currentHoistScopeRoot = previousHoistScopeRoot
}

// Java method: NodeTraversal.handleClass(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_handleClass_3995(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	this_.currentNode = n
	if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
		return
	}
	var className pkg_rhino.Node = n.GetFirstChild()
	_ = className
	var extendsClause pkg_rhino.Node = className.GetNext()
	_ = extendsClause
	var body pkg_rhino.Node = extendsClause.GetNext()
	_ = body
	var isClassExpression bool = NodeUtil_IsClassExpression(n)
	_ = isClassExpression
	this_.priv_NodeTraversal_traverseBranch_3993(extendsClause, n)
	if !(isClassExpression) {
		this_.priv_NodeTraversal_traverseBranch_3993(className, n)
	}
	this_.currentNode = n
	this_.priv_NodeTraversal_pushScope_4001(n)
	if isClassExpression {
		this_.priv_NodeTraversal_traverseBranch_3993(className, n)
	}
	this_.priv_NodeTraversal_traverseBranch_3993(body, n)
	this_.priv_NodeTraversal_popScope_4004()
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

// Java method: NodeTraversal.handleClassMembers(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal) priv_NodeTraversal_handleClassMembers_3996(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	this_.currentNode = n
	if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
		return
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for child != 0 {
			var next pkg_rhino.Node = child.GetNext()
			_ = next
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				this_.priv_NodeTraversal_traverseBranch_3993(child.GetFirstChild(), child)
			}
			child = next
		}
	}
	{
		var child_2 pkg_rhino.Node = n.GetFirstChild()
		_ = child_2
		for child_2 != 0 {
			var next_2 pkg_rhino.Node = child_2.GetNext()
			_ = next_2
			switch child_2.TokenSwitchOrd_() {
			case int32(129) /* Token.COMPUTED_PROP */ :
				this_.currentNode = n
				if this_.callback.ShouldTraverse_3(this_, child_2, n) {
					this_.priv_NodeTraversal_traverseBranch_3993(child_2.GetLastChild(), child_2)
					this_.currentNode = n
					this_.callback.Visit_3(this_, child_2, n)
				}
			case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
				this_.currentNode = n
				var previousHoistScopeRoot pkg_rhino.Node = this_.currentHoistScopeRoot
				_ = previousHoistScopeRoot
				this_.currentHoistScopeRoot = n
				this_.priv_NodeTraversal_pushScope_4001(child_2)
				if this_.callback.ShouldTraverse_3(this_, child_2, n) {
					if child_2.HasTwoChildren() {
						this_.priv_NodeTraversal_traverseBranch_3993(child_2.GetLastChild(), child_2)
					}
					this_.currentNode = n
					this_.callback.Visit_3(this_, child_2, n)
				}
				this_.priv_NodeTraversal_popScope_4004()
				this_.currentHoistScopeRoot = previousHoistScopeRoot
			case int32(108) /* Token.MEMBER_FIELD_DEF */ :
				this_.priv_NodeTraversal_handleMemberFieldDef_3997(n, child_2)
			case int32(90) /* Token.BLOCK */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(172) /* Token.MEMBER_VARIABLE_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
				this_.priv_NodeTraversal_traverseBranch_3993(child_2, n)
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid class member: "+j2gort.StrVal(child_2.GetToken()))))
			}
			child_2 = next_2
		}
	}
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

// Java method: NodeTraversal.handleMemberFieldDef(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node child)
func (this_ *NodeTraversal) priv_NodeTraversal_handleMemberFieldDef_3997(n pkg_rhino.Node, child pkg_rhino.Node) {
	_ = n
	_ = child
	var previousHoistScopeRoot pkg_rhino.Node = this_.currentHoistScopeRoot
	_ = previousHoistScopeRoot
	this_.currentHoistScopeRoot = n
	this_.priv_NodeTraversal_pushScope_4001(child)
	this_.priv_NodeTraversal_traverseBranch_3993(child, n)
	this_.priv_NodeTraversal_popScope_4004()
	this_.currentHoistScopeRoot = previousHoistScopeRoot
}

// Java method: NodeTraversal.traverseChildren(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_traverseChildren_3998(n pkg_rhino.Node) {
	_ = n
	curBlkIdx_ := uint32(n) >> nodestore.BlockShift
	curBlk_ := nodestore.BlockOf(nodestore.Node(n))
	nFirst_, _, _ := curBlk_.FirstNextTok(nodestore.Node(n))
	for child_ := pkg_rhino.Node(nFirst_); child_ != 0; {
		var cFirst_, next_ pkg_rhino.Node
		var cTok_ uint8
		if idx_ := uint32(child_) >> nodestore.BlockShift; idx_ == curBlkIdx_ {
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		} else {
			curBlk_ = nodestore.BlockOf(nodestore.Node(child_))
			curBlkIdx_ = idx_
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		}
		if cFirst_ == 0 && cTok_ != 0xFF && !nodeTraversalSpecialOrBlockTok_[cTok_] {
			if this_.alwaysTraverse_ {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			this_.currentNode = child_
			if !this_.callback.ShouldTraverse_3(this_, child_, n) {
				child_ = next_
				continue
			}
			var gc_ pkg_rhino.Node
			if cIdx_ := uint32(child_) >> nodestore.BlockShift; cIdx_ == curBlkIdx_ {
				gf_, _, _ := curBlk_.FirstNextTok(nodestore.Node(child_))
				gc_ = pkg_rhino.Node(gf_)
			} else {
				curBlk_ = nodestore.BlockOf(nodestore.Node(child_))
				curBlkIdx_ = cIdx_
				gf_, _, _ := curBlk_.FirstNextTok(nodestore.Node(child_))
				gc_ = pkg_rhino.Node(gf_)
			}
			if gc_ == 0 {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			for gc_ != 0 {
				var gFirst_, gNext_ pkg_rhino.Node
				var gTok_ uint8
				if gIdx_ := uint32(gc_) >> nodestore.BlockShift; gIdx_ == curBlkIdx_ {
					gf_, gnx_, gt_ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gFirst_, gNext_, gTok_ = pkg_rhino.Node(gf_), pkg_rhino.Node(gnx_), gt_
				} else {
					curBlk_ = nodestore.BlockOf(nodestore.Node(gc_))
					curBlkIdx_ = gIdx_
					gf_, gnx_, gt_ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gFirst_, gNext_, gTok_ = pkg_rhino.Node(gf_), pkg_rhino.Node(gnx_), gt_
				}
				this_.traverseBranchBlk_(gc_, child_, gFirst_, gTok_, curBlkIdx_, curBlk_)
				gc_ = gNext_
			}
			this_.currentNode = child_
			if !this_.noOpVisit_ {
				this_.callback.Visit_3(this_, child_, n)
			}
			child_ = next_
			continue
		}
		this_.traverseBranchBlk_(child_, n, cFirst_, cTok_, curBlkIdx_, curBlk_)
		child_ = next_
	}

}

// Java method: NodeTraversal.getEnclosingFunction()
func (this_ *NodeTraversal) GetEnclosingFunction() pkg_rhino.Node {
	return this_.currentFunction
}

// Java method: NodeTraversal.ensureScopeCapacity()
func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_4000() {
	if this_.scopeDepth == int32(len(this_.scopeNodes)) {
		if this_.scopeDepth == 0 {
			this_.scopeNodes = this_.scopeNodesBuf_[:12]
			this_.scopeObjs = this_.scopeObjsBuf_[:12]
			return
		}
		newLen := int(this_.scopeDepth) << 1
		newScopeNodes := make([]pkg_rhino.Node, newLen)
		copy(newScopeNodes, this_.scopeNodes)
		this_.scopeNodes = newScopeNodes
		newScopeObjs := make([]AbstractScope, newLen)
		copy(newScopeObjs, this_.scopeObjs)
		this_.scopeObjs = newScopeObjs
	}
}

// Java method: NodeTraversal.pushScope(com.google.javascript.rhino.Node node)
func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_4001(node pkg_rhino.Node) {
	_ = node
	j2gort.CheckNotNullNode(this_.currentNode)
	j2gort.CheckNotNullNode(node)
	this_.priv_NodeTraversal_ensureScopeCapacity_4000()
	{
		a_19 := this_.scopeNodes
		i_20 := this_.scopeDepth
		a_19[i_20] = node
	}
	{
		a_21 := this_.scopeObjs
		i_22 := this_.scopeDepth
		a_21[i_22] = nil
	}
	this_.scopeDepth = this_.scopeDepth + 1
	if this_.scopeCallback != nil && !this_.scopeCallback.IsJavaNil_() {
		this_.scopeCallback.EnterScope(this_)
	}
}

// Java method: NodeTraversal.pushScope(com.google.javascript.jscomp.AbstractScope<?,?> s)
func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_4002(s AbstractScope) {
	_ = s
	this_.priv_NodeTraversal_pushScope_4003(s, false)
}

// Java method: NodeTraversal.pushScope(com.google.javascript.jscomp.AbstractScope<?,?> s, boolean quietly)
func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_4003(s AbstractScope, quietly bool) {
	_ = s
	_ = quietly
	j2gort.CheckNotNullNode(this_.currentNode)
	this_.priv_NodeTraversal_ensureScopeCapacity_4000()
	{
		a_23 := this_.scopeNodes
		i_24 := this_.scopeDepth
		a_23[i_24] = s.GetAbstractScopeBase().GetRootNode_0()
	}
	{
		a_25 := this_.scopeObjs
		i_26 := this_.scopeDepth
		a_25[i_26] = s
	}
	this_.scopeDepth = this_.scopeDepth + 1
	if !(quietly) && (this_.scopeCallback != nil && !this_.scopeCallback.IsJavaNil_()) {
		this_.scopeCallback.EnterScope(this_)
	}
}

// Java method: NodeTraversal.popScope()
func (this_ *NodeTraversal) priv_NodeTraversal_popScope_4004() {
	this_.priv_NodeTraversal_popScope_4005(false)
}

// Java method: NodeTraversal.popScope(boolean quietly)
func (this_ *NodeTraversal) priv_NodeTraversal_popScope_4005(quietly bool) {
	_ = quietly
	if !(quietly) && (this_.scopeCallback != nil && !this_.scopeCallback.IsJavaNil_()) {
		this_.scopeCallback.ExitScope(this_)
	}
	this_.scopeDepth = this_.scopeDepth - 1
	{
		a_27 := this_.scopeNodes
		i_28 := this_.scopeDepth
		a_27[i_28] = pkg_rhino.Node(0)
	}
	{
		a_29 := this_.scopeObjs
		i_30 := this_.scopeDepth
		a_29[i_30] = nil
	}
}

// Java method: NodeTraversal.getScopeRoot()
func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node {
	if this_.scopeDepth > int32(0) {
		return this_.scopeNodes[(this_.scopeDepth - int32(1))]
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: NodeTraversal.getScopeDepth()
func (this_ *NodeTraversal) GetScopeDepth() int32 {
	if !(this_.scopeDepth > int32(0)) {
		j2gort.CheckState(false)
	}
	return (this_.scopeDepth - int32(1))
}

// Java method: NodeTraversal.getScope()
func (this_ *NodeTraversal) GetScope() *Scope {
	return this_.GetAbstractScope_0().Untyped()
}

// Java method: NodeTraversal.getTypedScope()
func (this_ *NodeTraversal) GetTypedScope() *TypedScope {
	return this_.GetAbstractScope_0().Typed()
}

// Java method: NodeTraversal.getAbstractScope()
func (this_ *NodeTraversal) GetAbstractScope_0() AbstractScope {
	return this_.priv_NodeTraversal_getAbstractScope_4011((this_.scopeDepth - int32(1)))
}

// Java method: NodeTraversal.getAbstractScope(int rootDepth)
func (this_ *NodeTraversal) priv_NodeTraversal_getAbstractScope_4011(rootDepth int32) AbstractScope {
	_ = rootDepth
	var s AbstractScope = this_.scopeObjs[rootDepth]
	_ = s
	if s != nil && !s.IsJavaNil_() {
		return s
	}
	var parentScope AbstractScope
	_ = parentScope
	if rootDepth > int32(0) {
		parentScope = this_.priv_NodeTraversal_getAbstractScope_4011((rootDepth - int32(1)))
	} else {
		parentScope = nil
	}
	var scope AbstractScope = j2gort.AsIface[AbstractScope](this_.scopeCreator.CreateScope_2_Node_AbstractScope(this_.scopeNodes[rootDepth], parentScope))
	_ = scope
	{
		a_31 := this_.scopeObjs
		i_32 := rootDepth
		a_31[i_32] = scope
	}
	return scope
}

// Java method: NodeTraversal.isHoistScope()
func (this_ *NodeTraversal) IsHoistScope() bool {
	return NodeTraversal_IsHoistScopeRootNode(this_.GetScopeRoot())
}

// Java method: NodeTraversal.getClosestHoistScopeRoot()
func (this_ *NodeTraversal) GetClosestHoistScopeRoot() pkg_rhino.Node {
	{
		var i int32 = (this_.scopeDepth - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			var rootNode pkg_rhino.Node = this_.scopeNodes[i]
			_ = rootNode
			if NodeTraversal_IsHoistScopeRootNode(rootNode) {
				return rootNode
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: NodeTraversal.getClosestContainerScope()
func (this_ *NodeTraversal) GetClosestContainerScope() AbstractScope {
	{
		var i int32 = (this_.scopeDepth - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			if !(NodeUtil_CreatesBlockScope(this_.scopeNodes[i])) {
				return this_.priv_NodeTraversal_getAbstractScope_4011(i)
			}
		}
	}
	return nil
}

// Java method: NodeTraversal.getClosestHoistScope()
func (this_ *NodeTraversal) GetClosestHoistScope() AbstractScope {
	{
		var i int32 = (this_.scopeDepth - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			if NodeTraversal_IsHoistScopeRootNode(this_.scopeNodes[i]) {
				return this_.priv_NodeTraversal_getAbstractScope_4011(i)
			}
		}
	}
	return nil
}

// Java static method: NodeTraversal.isHoistScopeRootNode(com.google.javascript.rhino.Node n)
func NodeTraversal_IsHoistScopeRootNode(n pkg_rhino.Node) bool {
	clinit_NodeTraversal()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(123) /* Token.MODULE_BODY */, int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */ :
		return true
	default:
		return NodeUtil_IsFunctionBlock(n)
	}
	panic("unreachable switch expression")
}

// Java method: NodeTraversal.getClosestScopeRootNodeBindingThisOrSuper()
func (this_ *NodeTraversal) GetClosestScopeRootNodeBindingThisOrSuper() pkg_rhino.Node {
	{
		var i int32 = (this_.scopeDepth - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			var rootNode pkg_rhino.Node = this_.scopeNodes[i]
			_ = rootNode
			switch rootNode.TokenSwitchOrd_() {
			case int32(73) /* Token.FUNCTION */ :
				if rootNode.IsArrowFunction() {
					continue
				}
				return rootNode
			case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */, int32(105) /* Token.CLASS */, int32(123) /* Token.MODULE_BODY */, int32(89) /* Token.ROOT */ :
				return rootNode
			case int32(90) /* Token.BLOCK */ :
				if NodeUtil_IsClassStaticBlock(rootNode) {
					return rootNode
				}
				continue
			default:
				continue
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: NodeTraversal.getScopeCreator()
func (this_ *NodeTraversal) GetScopeCreator() ScopeCreator {
	return this_.scopeCreator
}

// Java method: NodeTraversal.inGlobalScope()
func (this_ *NodeTraversal) InGlobalScope() bool {
	return (this_.GetScopeDepth() == int32(0))
}

// Java method: NodeTraversal.inModuleScope()
func (this_ *NodeTraversal) InModuleScope() bool {
	return NodeUtil_IsModuleScopeRoot(this_.GetScopeRoot())
}

// Java method: NodeTraversal.inGlobalOrModuleScope()
func (this_ *NodeTraversal) InGlobalOrModuleScope() bool {
	return (this_.InGlobalScope() || this_.InModuleScope())
}

// Java method: NodeTraversal.inFunctionBlockScope()
func (this_ *NodeTraversal) InFunctionBlockScope() bool {
	return NodeUtil_IsFunctionBlock(this_.GetScopeRoot())
}

// Java method: NodeTraversal.inGlobalHoistScope()
func (this_ *NodeTraversal) InGlobalHoistScope() bool {
	return (this_.currentHoistScopeRoot == 0)
}

// Java method: NodeTraversal.inModuleHoistScope()
func (this_ *NodeTraversal) InModuleHoistScope() bool {
	var moduleRoot pkg_rhino.Node = this_.currentHoistScopeRoot
	_ = moduleRoot
	if moduleRoot == 0 {
		return false
	}
	if moduleRoot.IsFunction() {
		moduleRoot = moduleRoot.GetLastChild()
	}
	return NodeUtil_IsModuleScopeRoot(moduleRoot)
}

// Java method: NodeTraversal.report(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType diagnosticType, java.lang.String[] arguments)
func (this_ *NodeTraversal) Report_3(n pkg_rhino.Node, diagnosticType *DiagnosticType, arguments []string) {
	_ = n
	_ = diagnosticType
	_ = arguments
	var error_ *JSError = JSError_Make_3(n, diagnosticType, arguments)
	_ = error_
	this_.compiler.Report_1(error_)
}

// Java method: NodeTraversal.report(com.google.javascript.rhino.Node start, com.google.javascript.rhino.Node end, com.google.javascript.jscomp.DiagnosticType diagnosticType, java.lang.String[] arguments)
func (this_ *NodeTraversal) Report_4(start pkg_rhino.Node, end pkg_rhino.Node, diagnosticType *DiagnosticType, arguments []string) {
	_ = start
	_ = end
	_ = diagnosticType
	_ = arguments
	var error_ *JSError = JSError_Make_4_Node_Node_DiagnosticType_slice_string(start, end, diagnosticType, arguments)
	_ = error_
	this_.compiler.Report_1(error_)
}

// Java method: NodeTraversal.reportCodeChange()
func (this_ *NodeTraversal) ReportCodeChange_0() {
	var changeScope pkg_rhino.Node = this_.currentChangeScope
	_ = changeScope
	j2gort.CheckNotNullNode(changeScope)
	if !(ChangeTracker_IsChangeScopeRoot(changeScope)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(changeScope))
	}
	this_.compiler.ReportChangeToChangeScope(changeScope)
}

// Java method: NodeTraversal.reportCodeChange(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) ReportCodeChange_1(n pkg_rhino.Node) {
	_ = n
	this_.compiler.ReportChangeToEnclosingScope(n)
}

// Java method: NodeTraversal.getCurrentScript()
func (this_ *NodeTraversal) GetCurrentScript() pkg_rhino.Node {
	if this_.currentScript == 0 {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "getCurrentScript not supported")))
	}
	return this_.currentScript
}

// Java method: NodeTraversal.setChangeScope(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_setChangeScope_4029(n pkg_rhino.Node) {
	_ = n
	this_.currentChangeScope = n
}

// Java method: NodeTraversal.getEnclosingScript(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_getEnclosingScript_4030(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	for (n != 0) && !(n.IsScript()) {
		n = n.GetParent()
	}
	return n
}

// Java method: NodeTraversal.initTraversal(com.google.javascript.rhino.Node traversalRoot)
func (this_ *NodeTraversal) priv_NodeTraversal_initTraversal_4031(traversalRoot pkg_rhino.Node) {
	_ = traversalRoot
	if Platform_IsThreadInterrupted() {
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", j2gort.NewExt_java_lang_InterruptedException())))
	}
	var hoistScopeRoot pkg_rhino.Node = NodeUtil_GetEnclosingHoistScopeRoot(traversalRoot)
	_ = hoistScopeRoot
	this_.currentHoistScopeRoot = hoistScopeRoot
	var iifeArg_1 pkg_rhino.
		Node
iifeSwitch_1:
	switch {
	default:

		if hoistScopeRoot != 0 {
			iifeArg_1 = hoistScopeRoot
			break iifeSwitch_1
		}
		iifeArg_1 = traversalRoot
		break iifeSwitch_1
	}
	var changeScope pkg_rhino.Node = ChangeTracker_GetEnclosingChangeScopeRoot(iifeArg_1)
	_ = changeScope
	this_.priv_NodeTraversal_setChangeScope_4029(changeScope)
	var enclosingFunction pkg_rhino.Node
	_ = enclosingFunction
	if hoistScopeRoot != 0 {
		enclosingFunction = NodeUtil_GetEnclosingFunction(hoistScopeRoot)
	} else {
		enclosingFunction = pkg_rhino.Node(0)
	}
	this_.currentFunction = enclosingFunction
	var script pkg_rhino.Node = this_.priv_NodeTraversal_getEnclosingScript_4030(changeScope)
	_ = script
	this_.currentScript = script
	this_.priv_NodeTraversal_clearScriptState_4034()
}

// Java method: NodeTraversal.initScopeRoots(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_initScopeRoots_4032(n pkg_rhino.Node) {
	_ = n
	var queuedScopeRoots j2gort.Deque[pkg_rhino.Node]
	_ = queuedScopeRoots
	for n != 0 {
		if this_.priv_NodeTraversal_isScopeRoot_4033(n) {
			queuedScopeRoots.AddFirst(j2gort.CheckNotNullNode(n))
		}
		n = n.GetParent()
	}
	{
		cur_33 := queuedScopeRoots.Iter()
		for cur_33.HasNext() {
			queuedScopeRoot := cur_33.Next()
			_ = queuedScopeRoot
			this_.priv_NodeTraversal_pushScope_4001(queuedScopeRoot)
		}
	}
}

// Java method: NodeTraversal.isScopeRoot(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_isScopeRoot_4033(n pkg_rhino.Node) bool {
	_ = n
	if n.IsRoot() && (n.GetParent() == 0) {
		return true
	} else if n.IsFunction() {
		return true
	} else if NodeUtil_CreatesBlockScope(n) {
		return true
	}
	return false
}

// Java method: NodeTraversal.clearScriptState()
func (this_ *NodeTraversal) priv_NodeTraversal_clearScriptState_4034() {
	this_.inputId = nil
	this_.sourceName = j2gort.NullString
	this_.compilerInput = nil
}

// Java method: NodeTraversal.getInputId()
func (this_ *NodeTraversal) GetInputId() *pkg_rhino.InputId {
	if (this_.currentScript != 0) && (this_.inputId == nil) {
		this_.inputId = this_.currentScript.GetInputId()
	}
	return this_.inputId
}

// Java method: NodeTraversal.getBestSourceFileName(com.google.javascript.rhino.Node n)
func (this_ *NodeTraversal) priv_NodeTraversal_getBestSourceFileName_4036(n pkg_rhino.Node) string {
	_ = n
	if n == 0 {
		return this_.GetSourceName()
	} else {
		return n.GetSourceFileName()
	}
}

// Java interface: com.google.javascript.jscomp.NodeTraversal.Callback
type NodeTraversal_Callback interface {
	Implements_jscomp_NodeTraversal_Callback()
	IsJavaNil_() bool
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func ToIface_NodeTraversal_Callback[T interface {
	NodeTraversal_Callback
	comparable
}](v T) NodeTraversal_Callback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_Callback[T NodeTraversal_Callback](v T) NodeTraversal_Callback {
	var vi NodeTraversal_Callback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.Callback", "com.google.javascript.jscomp.NodeTraversal$Callback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_Callback](v_)
	})
}

// Java interface: com.google.javascript.jscomp.NodeTraversal.ScopedCallback
type NodeTraversal_ScopedCallback interface {
	Implements_jscomp_NodeTraversal_ScopedCallback()
	IsJavaNil_() bool
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

func ToIface_NodeTraversal_ScopedCallback[T interface {
	NodeTraversal_ScopedCallback
	comparable
}](v T) NodeTraversal_ScopedCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_ScopedCallback[T NodeTraversal_ScopedCallback](v T) NodeTraversal_ScopedCallback {
	var vi NodeTraversal_ScopedCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.ScopedCallback", "com.google.javascript.jscomp.NodeTraversal$ScopedCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_ScopedCallback](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractPostOrderCallback
type NodeTraversal_AbstractPostOrderCallback interface {
	Implements_jscomp_NodeTraversal_AbstractPostOrderCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractPostOrderCallbackBase() *NodeTraversal_AbstractPostOrderCallbackBase
	NodeTraversal_Callback
}

func ToIface_NodeTraversal_AbstractPostOrderCallback[T interface {
	NodeTraversal_AbstractPostOrderCallback
	comparable
}](v T) NodeTraversal_AbstractPostOrderCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractPostOrderCallback[T NodeTraversal_AbstractPostOrderCallback](v T) NodeTraversal_AbstractPostOrderCallback {
	var vi NodeTraversal_AbstractPostOrderCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractPostOrderCallback
type NodeTraversal_AbstractPostOrderCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *NodeTraversal_AbstractPostOrderCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_AbstractPostOrderCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractPostOrderCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractPostOrderCallback", "AbstractPostOrderCallback"
}
func (this_ *NodeTraversal_AbstractPostOrderCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractPostOrderCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractPostOrderCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractPostOrderCallback](v_)
	})
}
func (b *NodeTraversal_AbstractPostOrderCallbackBase) GetNodeTraversal_AbstractPostOrderCallbackBase() *NodeTraversal_AbstractPostOrderCallbackBase {
	return b
}
func (*NodeTraversal_AbstractPostOrderCallbackBase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*NodeTraversal_AbstractPostOrderCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *NodeTraversal_AbstractPostOrderCallbackBase) initFields_NodeTraversal_AbstractPostOrderCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractPostOrderCallback(this_ *NodeTraversal_AbstractPostOrderCallbackBase) {
	this_.initFields_NodeTraversal_AbstractPostOrderCallback()
}

// Java method: AbstractPostOrderCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return true
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.ExternsSkippingCallback
type NodeTraversal_ExternsSkippingCallback interface {
	Implements_jscomp_NodeTraversal_ExternsSkippingCallback()
	IsJavaNil_() bool
	GetNodeTraversal_ExternsSkippingCallbackBase() *NodeTraversal_ExternsSkippingCallbackBase
	NodeTraversal_Callback
}

func ToIface_NodeTraversal_ExternsSkippingCallback[T interface {
	NodeTraversal_ExternsSkippingCallback
	comparable
}](v T) NodeTraversal_ExternsSkippingCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_ExternsSkippingCallback[T NodeTraversal_ExternsSkippingCallback](v T) NodeTraversal_ExternsSkippingCallback {
	var vi NodeTraversal_ExternsSkippingCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.ExternsSkippingCallback
type NodeTraversal_ExternsSkippingCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *NodeTraversal_ExternsSkippingCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_ExternsSkippingCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_ExternsSkippingCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$ExternsSkippingCallback", "ExternsSkippingCallback"
}
func (this_ *NodeTraversal_ExternsSkippingCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.ExternsSkippingCallback", "com.google.javascript.jscomp.NodeTraversal$ExternsSkippingCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_ExternsSkippingCallback](v_)
	})
}
func (b *NodeTraversal_ExternsSkippingCallbackBase) GetNodeTraversal_ExternsSkippingCallbackBase() *NodeTraversal_ExternsSkippingCallbackBase {
	return b
}
func (*NodeTraversal_ExternsSkippingCallbackBase) Implements_jscomp_NodeTraversal_ExternsSkippingCallback() {
}
func (*NodeTraversal_ExternsSkippingCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *NodeTraversal_ExternsSkippingCallbackBase) initFields_NodeTraversal_ExternsSkippingCallback() {
	_ = this_
}

func InitNodeTraversal_ExternsSkippingCallback(this_ *NodeTraversal_ExternsSkippingCallbackBase) {
	this_.initFields_NodeTraversal_ExternsSkippingCallback()
}

// Java method: ExternsSkippingCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_ExternsSkippingCallbackBase) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return ((!(n.IsScript()) || !(n.IsFromExterns())) || NodeUtil_IsFromTypeSummary(n))
}

// Java interface: com.google.javascript.jscomp.NodeTraversal.AbstractPostOrderCallbackInterface
type NodeTraversal_AbstractPostOrderCallbackInterface interface {
	Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface()
	IsJavaNil_() bool
	Visit(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func ToIface_NodeTraversal_AbstractPostOrderCallbackInterface[T interface {
	NodeTraversal_AbstractPostOrderCallbackInterface
	comparable
}](v T) NodeTraversal_AbstractPostOrderCallbackInterface {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractPostOrderCallbackInterface[T NodeTraversal_AbstractPostOrderCallbackInterface](v T) NodeTraversal_AbstractPostOrderCallbackInterface {
	var vi NodeTraversal_AbstractPostOrderCallbackInterface = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractPostOrderCallbackInterface", "com.google.javascript.jscomp.NodeTraversal$AbstractPostOrderCallbackInterface", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractPostOrderCallbackInterface](v_)
	})
}

type lambda_NodeTraversal_AbstractPostOrderCallbackInterface struct {
	fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func (l *lambda_NodeTraversal_AbstractPostOrderCallbackInterface) IsJavaNil_() bool { return l == nil }
func NewLambda_NodeTraversal_AbstractPostOrderCallbackInterface(fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)) NodeTraversal_AbstractPostOrderCallbackInterface {
	if fn == nil {
		return nil
	}
	return &lambda_NodeTraversal_AbstractPostOrderCallbackInterface{fn: fn}
}

func (*lambda_NodeTraversal_AbstractPostOrderCallbackInterface) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {
}
func (l *lambda_NodeTraversal_AbstractPostOrderCallbackInterface) Visit(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	l.fn(p0_, p1_, p2_)
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractPreOrderCallback
type NodeTraversal_AbstractPreOrderCallback interface {
	Implements_jscomp_NodeTraversal_AbstractPreOrderCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractPreOrderCallbackBase() *NodeTraversal_AbstractPreOrderCallbackBase
	NodeTraversal_Callback
}

func ToIface_NodeTraversal_AbstractPreOrderCallback[T interface {
	NodeTraversal_AbstractPreOrderCallback
	comparable
}](v T) NodeTraversal_AbstractPreOrderCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractPreOrderCallback[T NodeTraversal_AbstractPreOrderCallback](v T) NodeTraversal_AbstractPreOrderCallback {
	var vi NodeTraversal_AbstractPreOrderCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractPreOrderCallback
type NodeTraversal_AbstractPreOrderCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *NodeTraversal_AbstractPreOrderCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_AbstractPreOrderCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractPreOrderCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractPreOrderCallback", "AbstractPreOrderCallback"
}
func (this_ *NodeTraversal_AbstractPreOrderCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractPreOrderCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractPreOrderCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractPreOrderCallback](v_)
	})
}
func (b *NodeTraversal_AbstractPreOrderCallbackBase) GetNodeTraversal_AbstractPreOrderCallbackBase() *NodeTraversal_AbstractPreOrderCallbackBase {
	return b
}
func (*NodeTraversal_AbstractPreOrderCallbackBase) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*NodeTraversal_AbstractPreOrderCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *NodeTraversal_AbstractPreOrderCallbackBase) initFields_NodeTraversal_AbstractPreOrderCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractPreOrderCallback(this_ *NodeTraversal_AbstractPreOrderCallbackBase) {
	this_.initFields_NodeTraversal_AbstractPreOrderCallback()
}

// Java method: AbstractPreOrderCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractPreOrderCallbackBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractScopedCallback
type NodeTraversal_AbstractScopedCallback interface {
	Implements_jscomp_NodeTraversal_AbstractScopedCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractScopedCallbackBase() *NodeTraversal_AbstractScopedCallbackBase
	NodeTraversal_ScopedCallback
}

func ToIface_NodeTraversal_AbstractScopedCallback[T interface {
	NodeTraversal_AbstractScopedCallback
	comparable
}](v T) NodeTraversal_AbstractScopedCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractScopedCallback[T NodeTraversal_AbstractScopedCallback](v T) NodeTraversal_AbstractScopedCallback {
	var vi NodeTraversal_AbstractScopedCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractScopedCallback
type NodeTraversal_AbstractScopedCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *NodeTraversal_AbstractScopedCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_AbstractScopedCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractScopedCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractScopedCallback", "AbstractScopedCallback"
}
func (this_ *NodeTraversal_AbstractScopedCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractScopedCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractScopedCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractScopedCallback](v_)
	})
}
func (b *NodeTraversal_AbstractScopedCallbackBase) GetNodeTraversal_AbstractScopedCallbackBase() *NodeTraversal_AbstractScopedCallbackBase {
	return b
}
func (*NodeTraversal_AbstractScopedCallbackBase) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {
}
func (*NodeTraversal_AbstractScopedCallbackBase) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*NodeTraversal_AbstractScopedCallbackBase) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *NodeTraversal_AbstractScopedCallbackBase) initFields_NodeTraversal_AbstractScopedCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractScopedCallback(this_ *NodeTraversal_AbstractScopedCallbackBase) {
	this_.initFields_NodeTraversal_AbstractScopedCallback()
}

// Java method: AbstractScopedCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractScopedCallbackBase) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return true
}

// Java method: AbstractScopedCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal_AbstractScopedCallbackBase) EnterScope(t *NodeTraversal) {
	_ = t
}

// Java method: AbstractScopedCallback.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal_AbstractScopedCallbackBase) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractCfgCallback
type NodeTraversal_AbstractCfgCallback interface {
	Implements_jscomp_NodeTraversal_AbstractCfgCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase
	NodeTraversal_ScopedCallback
	GetControlFlowGraph(AbstractCompiler) ControlFlowGraph[pkg_rhino.Node]
	EnterScopeWithCfg(*NodeTraversal)
	ExitScopeWithCfg(*NodeTraversal)
}

func ToIface_NodeTraversal_AbstractCfgCallback[T interface {
	NodeTraversal_AbstractCfgCallback
	comparable
}](v T) NodeTraversal_AbstractCfgCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractCfgCallback[T NodeTraversal_AbstractCfgCallback](v T) NodeTraversal_AbstractCfgCallback {
	var vi NodeTraversal_AbstractCfgCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractCfgCallback
type NodeTraversal_AbstractCfgCallbackBase struct {
	self_      NodeTraversal_AbstractCfgCallback
	cfgs       *j2gort.Ext_java_util_ArrayDeque `gson:"cfgs"`
	currentCfg ControlFlowGraph[pkg_rhino.Node] `gson:"currentCfg"`
	idHash_    atomic.Int32
}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(NodeTraversal_AbstractCfgCallback)
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractCfgCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractCfgCallback", "AbstractCfgCallback"
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractCfgCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractCfgCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractCfgCallback](v_)
	})
}
func (b *NodeTraversal_AbstractCfgCallbackBase) GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase {
	return b
}
func (*NodeTraversal_AbstractCfgCallbackBase) Implements_jscomp_NodeTraversal_AbstractCfgCallback() {}
func (*NodeTraversal_AbstractCfgCallbackBase) Implements_jscomp_NodeTraversal_ScopedCallback()      {}
func (*NodeTraversal_AbstractCfgCallbackBase) Implements_jscomp_NodeTraversal_Callback()            {}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) initFields_NodeTraversal_AbstractCfgCallback() {
	_ = this_
	this_.cfgs = j2gort.NewExt_java_util_ArrayDeque()
}

func InitNodeTraversal_AbstractCfgCallback(this_ *NodeTraversal_AbstractCfgCallbackBase) {
	this_.initFields_NodeTraversal_AbstractCfgCallback()
}

// Java method: AbstractCfgCallback.getControlFlowGraph(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) GetControlFlowGraph(compiler AbstractCompiler) ControlFlowGraph[pkg_rhino.Node] {
	_ = compiler
	var cached ControlFlowGraph[pkg_rhino.Node] = this_.currentCfg
	_ = cached
	if cached != nil && !cached.IsJavaNil_() {
		return cached
	}
	var o any = this_.cfgs.Peek()
	_ = o
	if !(o != nil) {
		j2gort.CheckState(false)
	}
	var result ControlFlowGraph[pkg_rhino.Node] = nil
	_ = result
	var cfgRoot pkg_rhino.Node
	_ = cfgRoot
	var ok_34 bool
	cfgRoot, ok_34 = j2gort.TryCast[pkg_rhino.Node](o)
	if ok_34 {
		result = ControlFlowAnalysis_ComputeCfgDirect(compiler, cfgRoot, false, true)
		this_.cfgs.Pop()
		this_.cfgs.Push(result)
	} else {
		result = j2gort.AsIface[ControlFlowGraph[pkg_rhino.Node]](o)
	}
	this_.currentCfg = result
	return result
}

// Java method: AbstractCfgCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	_ = t
	var currentScopeRoot pkg_rhino.Node = t.GetScopeRoot()
	_ = currentScopeRoot
	if NodeUtil_IsValidCfgRoot(currentScopeRoot) {
		j2gort.DequePushNode(this_.cfgs, uint32(currentScopeRoot))
		this_.currentCfg = nil
	}
	this_.self_.EnterScopeWithCfg(t)
}

// Java method: AbstractCfgCallback.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	_ = t
	defer func() {
		var currentScopeRoot pkg_rhino.Node = t.GetScopeRoot()
		_ = currentScopeRoot
		if NodeUtil_IsValidCfgRoot(currentScopeRoot) {
			ControlFlowAnalysis_ReleaseCfg(this_.cfgs.Pop())
			var top any = this_.cfgs.Peek()
			_ = top
			{
				var iifeCond_1 bool
				{

					v_, ok_ := top.(interface {
						Implements_jscomp_ControlFlowGraph()
						IsJavaNil_() bool
					})
					iifeCond_1 = ok_ && !v_.IsJavaNil_()
				}
				if iifeCond_1 {
					this_.currentCfg = j2gort.AsIface[ControlFlowGraph[pkg_rhino.Node]](top)
				} else {
					this_.currentCfg = nil
				}
			}
		}
	}()
	{
		this_.self_.ExitScopeWithCfg(t)
	}
}

// Java method: AbstractCfgCallback.enterScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) EnterScopeWithCfg(t *NodeTraversal) {
	_ = t
}

// Java method: AbstractCfgCallback.exitScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) ExitScopeWithCfg(t *NodeTraversal) {
	_ = t
}

// Java method: AbstractCfgCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}

// Java method: AbstractCfgCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractCfgCallbackBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractShallowCallback
type NodeTraversal_AbstractShallowCallback interface {
	Implements_jscomp_NodeTraversal_AbstractShallowCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractShallowCallbackBase() *NodeTraversal_AbstractShallowCallbackBase
	NodeTraversal_Callback
}

func ToIface_NodeTraversal_AbstractShallowCallback[T interface {
	NodeTraversal_AbstractShallowCallback
	comparable
}](v T) NodeTraversal_AbstractShallowCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractShallowCallback[T NodeTraversal_AbstractShallowCallback](v T) NodeTraversal_AbstractShallowCallback {
	var vi NodeTraversal_AbstractShallowCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractShallowCallback
type NodeTraversal_AbstractShallowCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *NodeTraversal_AbstractShallowCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_AbstractShallowCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractShallowCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractShallowCallback", "AbstractShallowCallback"
}
func (this_ *NodeTraversal_AbstractShallowCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractShallowCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractShallowCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractShallowCallback](v_)
	})
}
func (b *NodeTraversal_AbstractShallowCallbackBase) GetNodeTraversal_AbstractShallowCallbackBase() *NodeTraversal_AbstractShallowCallbackBase {
	return b
}
func (*NodeTraversal_AbstractShallowCallbackBase) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {
}
func (*NodeTraversal_AbstractShallowCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *NodeTraversal_AbstractShallowCallbackBase) initFields_NodeTraversal_AbstractShallowCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractShallowCallback(this_ *NodeTraversal_AbstractShallowCallbackBase) {
	this_.initFields_NodeTraversal_AbstractShallowCallback()
}

// Java method: AbstractShallowCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractShallowCallbackBase) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return (((parent == 0) || !(parent.IsFunction())) || (n == parent.GetFirstChild()))
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractShallowStatementCallback
type NodeTraversal_AbstractShallowStatementCallback interface {
	Implements_jscomp_NodeTraversal_AbstractShallowStatementCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractShallowStatementCallbackBase() *NodeTraversal_AbstractShallowStatementCallbackBase
	NodeTraversal_Callback
}

func ToIface_NodeTraversal_AbstractShallowStatementCallback[T interface {
	NodeTraversal_AbstractShallowStatementCallback
	comparable
}](v T) NodeTraversal_AbstractShallowStatementCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractShallowStatementCallback[T NodeTraversal_AbstractShallowStatementCallback](v T) NodeTraversal_AbstractShallowStatementCallback {
	var vi NodeTraversal_AbstractShallowStatementCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractShallowStatementCallback
type NodeTraversal_AbstractShallowStatementCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *NodeTraversal_AbstractShallowStatementCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_AbstractShallowStatementCallbackBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*NodeTraversal_AbstractShallowStatementCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractShallowStatementCallback", "AbstractShallowStatementCallback"
}
func (this_ *NodeTraversal_AbstractShallowStatementCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractShallowStatementCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractShallowStatementCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractShallowStatementCallback](v_)
	})
}
func (b *NodeTraversal_AbstractShallowStatementCallbackBase) GetNodeTraversal_AbstractShallowStatementCallbackBase() *NodeTraversal_AbstractShallowStatementCallbackBase {
	return b
}
func (*NodeTraversal_AbstractShallowStatementCallbackBase) Implements_jscomp_NodeTraversal_AbstractShallowStatementCallback() {
}
func (*NodeTraversal_AbstractShallowStatementCallbackBase) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *NodeTraversal_AbstractShallowStatementCallbackBase) initFields_NodeTraversal_AbstractShallowStatementCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractShallowStatementCallback(this_ *NodeTraversal_AbstractShallowStatementCallbackBase) {
	this_.initFields_NodeTraversal_AbstractShallowStatementCallback()
}

// Java method: AbstractShallowStatementCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractShallowStatementCallbackBase) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return NodeUtil_IsShallowStatementTree(parent)
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractModuleCallback
type NodeTraversal_AbstractModuleCallback interface {
	Implements_jscomp_NodeTraversal_AbstractModuleCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractModuleCallbackBase() *NodeTraversal_AbstractModuleCallbackBase
	NodeTraversal_Callback
	EnterModule(ModuleMetadataMap_ModuleMetadata, pkg_rhino.Node)
	ExitModule(ModuleMetadataMap_ModuleMetadata, pkg_rhino.Node)
	ShouldTraverse_4(*NodeTraversal, pkg_rhino.Node, ModuleMetadataMap_ModuleMetadata, pkg_rhino.Node) bool
	Visit_4(*NodeTraversal, pkg_rhino.Node, ModuleMetadataMap_ModuleMetadata, pkg_rhino.Node)
}

func ToIface_NodeTraversal_AbstractModuleCallback[T interface {
	NodeTraversal_AbstractModuleCallback
	comparable
}](v T) NodeTraversal_AbstractModuleCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractModuleCallback[T NodeTraversal_AbstractModuleCallback](v T) NodeTraversal_AbstractModuleCallback {
	var vi NodeTraversal_AbstractModuleCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractModuleCallback
type NodeTraversal_AbstractModuleCallbackBase struct {
	self_             NodeTraversal_AbstractModuleCallback
	compiler          AbstractCompiler                 `gson:"compiler"`
	moduleMetadataMap *ModuleMetadataMap               `gson:"moduleMetadataMap"`
	currentModule     ModuleMetadataMap_ModuleMetadata `gson:"currentModule"`
	scopeRoot         pkg_rhino.Node                   `gson:"scopeRoot"`
	idHash_           atomic.Int32
	inLoadModule      bool `gson:"inLoadModule"`
}

func (this_ *NodeTraversal_AbstractModuleCallbackBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(NodeTraversal_AbstractModuleCallback)
}
func (this_ *NodeTraversal_AbstractModuleCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractModuleCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractModuleCallback", "AbstractModuleCallback"
}
func (this_ *NodeTraversal_AbstractModuleCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractModuleCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractModuleCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractModuleCallback](v_)
	})
}
func (b *NodeTraversal_AbstractModuleCallbackBase) GetNodeTraversal_AbstractModuleCallbackBase() *NodeTraversal_AbstractModuleCallbackBase {
	return b
}
func (*NodeTraversal_AbstractModuleCallbackBase) Implements_jscomp_NodeTraversal_AbstractModuleCallback() {
}
func (*NodeTraversal_AbstractModuleCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

var static_NodeTraversal_AbstractModuleCallback_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_NodeTraversal_AbstractModuleCallback_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_NodeTraversal_AbstractModuleCallback()
	return static_NodeTraversal_AbstractModuleCallback_GOOG_MODULE
}
func StaticSet_NodeTraversal_AbstractModuleCallback_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_NodeTraversal_AbstractModuleCallback()
	static_NodeTraversal_AbstractModuleCallback_GOOG_MODULE = v
}

var clinitDone_NodeTraversal_AbstractModuleCallback uint32
var clinitState_NodeTraversal_AbstractModuleCallback j2gort.ClinitState

func clinit_NodeTraversal_AbstractModuleCallback() {
	if atomic.LoadUint32(&clinitDone_NodeTraversal_AbstractModuleCallback) != 0 {
		return
	}
	clinitSlow_NodeTraversal_AbstractModuleCallback()
}

func clinitSlow_NodeTraversal_AbstractModuleCallback() {
	if !j2gort.EnterClinit(&clinitState_NodeTraversal_AbstractModuleCallback, &clinitDone_NodeTraversal_AbstractModuleCallback) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_NodeTraversal_AbstractModuleCallback, &clinitDone_NodeTraversal_AbstractModuleCallback)
	static_NodeTraversal_AbstractModuleCallback_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
}

func (this_ *NodeTraversal_AbstractModuleCallbackBase) initFields_NodeTraversal_AbstractModuleCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractModuleCallback(this_ *NodeTraversal_AbstractModuleCallbackBase, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	_ = compiler
	_ = moduleMetadataMap
	this_.initFields_NodeTraversal_AbstractModuleCallback()
	this_.compiler = compiler
	this_.moduleMetadataMap = moduleMetadataMap
}

// Java method: AbstractModuleCallback.enterModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata currentModule, com.google.javascript.rhino.Node moduleScopeRoot)
func (this_ *NodeTraversal_AbstractModuleCallbackBase) EnterModule(currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = currentModule
	_ = moduleScopeRoot
}

// Java method: AbstractModuleCallback.exitModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata oldModule, com.google.javascript.rhino.Node moduleScopeRoot)
func (this_ *NodeTraversal_AbstractModuleCallbackBase) ExitModule(oldModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = oldModule
	_ = moduleScopeRoot
}

// Java method: AbstractModuleCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractModuleCallbackBase) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		this_.currentModule = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByPath(), t.GetInput().GetPath().String()))
		j2gort.CheckNotNullIface(this_.currentModule)
		if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
			this_.scopeRoot = n.GetFirstChild()
		} else {
			this_.scopeRoot = n
		}
		this_.self_.EnterModule(this_.currentModule, this_.scopeRoot)
	case int32(90) /* Token.BLOCK */ :
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			this_.scopeRoot = n
			this_.inLoadModule = true
		}
	case int32(28) /* Token.CALL */ :
		if this_.inLoadModule && StaticGet_NodeTraversal_AbstractModuleCallback_GOOG_MODULE().Matches(n.GetFirstChild()) {
			var newModule ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByGoogNamespace(), n.GetLastChild().GetString()))
			_ = newModule
			j2gort.CheckNotNullIface(newModule)
			if !(newModule == this_.currentModule || ((newModule == nil || newModule.IsJavaNil_()) && (this_.currentModule == nil || this_.currentModule.IsJavaNil_()))) {
				this_.currentModule = newModule
				this_.self_.EnterModule(this_.currentModule, this_.scopeRoot)
			}
		}
	default:
	}
	return this_.self_.ShouldTraverse_4(t, n, this_.currentModule, this_.scopeRoot)
}

// Java method: AbstractModuleCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata currentModule, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node moduleScopeRoot)
func (this_ *NodeTraversal_AbstractModuleCallbackBase) ShouldTraverse_4(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = currentModule
	_ = moduleScopeRoot
	return true
}

// Java method: AbstractModuleCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractModuleCallbackBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		j2gort.CheckNotNullIface(this_.currentModule)
		this_.self_.ExitModule(this_.currentModule, this_.scopeRoot)
		this_.currentModule = nil
		this_.scopeRoot = pkg_rhino.Node(0)
	case int32(90) /* Token.BLOCK */ :
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			j2gort.CheckNotNullIface(this_.currentModule)
			this_.self_.ExitModule(this_.currentModule, this_.scopeRoot)
			this_.scopeRoot = n.GetGrandparent().GetGrandparent()
			this_.inLoadModule = false
			this_.currentModule = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByPath(), t.GetInput().GetPath().String()))
			j2gort.CheckNotNullIface(this_.currentModule)
		}
	default:
	}
	this_.self_.Visit_4(t, n, this_.currentModule, this_.scopeRoot)
}

// Java method: AbstractModuleCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata currentModule, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node moduleScopeRoot)
func (this_ *NodeTraversal_AbstractModuleCallbackBase) Visit_4(t *NodeTraversal, n pkg_rhino.Node, currentModule ModuleMetadataMap_ModuleMetadata, moduleScopeRoot pkg_rhino.Node) {
	_ = t
	_ = n
	_ = currentModule
	_ = moduleScopeRoot
}

// Java polymorphic class interface: com.google.javascript.jscomp.NodeTraversal.AbstractChangedScopeCallback
type NodeTraversal_AbstractChangedScopeCallback interface {
	Implements_jscomp_NodeTraversal_AbstractChangedScopeCallback()
	IsJavaNil_() bool
	GetNodeTraversal_AbstractChangedScopeCallbackBase() *NodeTraversal_AbstractChangedScopeCallbackBase
	NodeTraversal_AbstractPreOrderCallback
	EnterChangedScopeRoot(AbstractCompiler, pkg_rhino.Node)
}

func ToIface_NodeTraversal_AbstractChangedScopeCallback[T interface {
	NodeTraversal_AbstractChangedScopeCallback
	comparable
}](v T) NodeTraversal_AbstractChangedScopeCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeTraversal_AbstractChangedScopeCallback[T NodeTraversal_AbstractChangedScopeCallback](v T) NodeTraversal_AbstractChangedScopeCallback {
	var vi NodeTraversal_AbstractChangedScopeCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.NodeTraversal.AbstractChangedScopeCallback
type NodeTraversal_AbstractChangedScopeCallbackBase struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	self_ NodeTraversal_AbstractChangedScopeCallback
}

func (this_ *NodeTraversal_AbstractChangedScopeCallbackBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(NodeTraversal_AbstractChangedScopeCallback)
}
func (this_ *NodeTraversal_AbstractChangedScopeCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_AbstractChangedScopeCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$AbstractChangedScopeCallback", "AbstractChangedScopeCallback"
}
func (this_ *NodeTraversal_AbstractChangedScopeCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.AbstractChangedScopeCallback", "com.google.javascript.jscomp.NodeTraversal$AbstractChangedScopeCallback", func(v_ any) bool {
		return j2gort.InstanceOf[NodeTraversal_AbstractChangedScopeCallback](v_)
	})
}
func (b *NodeTraversal_AbstractChangedScopeCallbackBase) GetNodeTraversal_AbstractChangedScopeCallbackBase() *NodeTraversal_AbstractChangedScopeCallbackBase {
	return b
}
func (*NodeTraversal_AbstractChangedScopeCallbackBase) Implements_jscomp_NodeTraversal_AbstractChangedScopeCallback() {
}
func (*NodeTraversal_AbstractChangedScopeCallbackBase) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*NodeTraversal_AbstractChangedScopeCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *NodeTraversal_AbstractChangedScopeCallbackBase) initFields_NodeTraversal_AbstractChangedScopeCallback() {
	_ = this_
}

func InitNodeTraversal_AbstractChangedScopeCallback(this_ *NodeTraversal_AbstractChangedScopeCallbackBase) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_NodeTraversal_AbstractChangedScopeCallback()
}

// Java method: AbstractChangedScopeCallback.enterChangedScopeRoot(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root)
func (this_ *NodeTraversal_AbstractChangedScopeCallbackBase) EnterChangedScopeRoot(compiler AbstractCompiler, root pkg_rhino.Node) {
	_ = compiler
	_ = root
	panic("abstract or native method invoked")
}

// Java method: AbstractChangedScopeCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal_AbstractChangedScopeCallbackBase) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if ChangeTracker_IsChangeScopeRoot(n) && t.GetCompiler().HasScopeChanged(n) {
		this_.self_.EnterChangedScopeRoot(t.GetCompiler(), n)
	}
	return true
}

// Java class: com.google.javascript.jscomp.NodeTraversal.Builder
type NodeTraversal_Builder struct {
	callback                                       NodeTraversal_Callback `gson:"callback"`
	compiler                                       AbstractCompiler       `gson:"compiler"`
	scopeCreator                                   ScopeCreator           `gson:"scopeCreator"`
	idHash_                                        atomic.Int32
	obeyDestructuringAndDefaultValueExecutionOrder bool `gson:"obeyDestructuringAndDefaultValueExecutionOrder"`
}

func (this_ *NodeTraversal_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal_Builder) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$Builder", "Builder"
}
func (this_ *NodeTraversal_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeTraversal_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.NodeTraversal.Builder", "com.google.javascript.jscomp.NodeTraversal$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeTraversal_Builder](v_)
	})
}

func (this_ *NodeTraversal_Builder) initFields_NodeTraversal_Builder() {
	_ = this_

}

// Java constructor: Builder()
func NewNodeTraversal_Builder() *NodeTraversal_Builder {
	this_ := &NodeTraversal_Builder{}
	InitNodeTraversal_Builder(this_)
	return this_
}

func InitNodeTraversal_Builder(this_ *NodeTraversal_Builder) {
	this_.initFields_NodeTraversal_Builder()
}

// Java method: Builder.setCallback(com.google.javascript.jscomp.NodeTraversal.Callback x)
func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(x NodeTraversal_Callback) *NodeTraversal_Builder {
	_ = x
	this_.callback = x
	return this_
}

// Java method: Builder.setCallback(com.google.javascript.jscomp.NodeTraversal.AbstractPostOrderCallbackInterface x)
func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal_Builder {
	_ = x
	this_.callback = NewNodeTraversal__Anon_64(this_, x)
	return this_
}

// Java method: Builder.setCompiler(com.google.javascript.jscomp.AbstractCompiler x)
func (this_ *NodeTraversal_Builder) SetCompiler(x AbstractCompiler) *NodeTraversal_Builder {
	_ = x
	this_.compiler = x
	return this_
}

// Java method: Builder.setScopeCreator(com.google.javascript.jscomp.ScopeCreator x)
func (this_ *NodeTraversal_Builder) SetScopeCreator(x ScopeCreator) *NodeTraversal_Builder {
	_ = x
	this_.scopeCreator = x
	return this_
}

// Java method: Builder.setObeyDestructuringAndDefaultValueExecutionOrder(boolean x)
func (this_ *NodeTraversal_Builder) SetObeyDestructuringAndDefaultValueExecutionOrder(x bool) *NodeTraversal_Builder {
	_ = x
	this_.obeyDestructuringAndDefaultValueExecutionOrder = x
	return this_
}

// Java method: Builder.build()
func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	return NewNodeTraversal(this_)
}

// Java method: Builder.traverse(com.google.javascript.rhino.Node root)
func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	t_ := acquireNodeTraversal_(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverse_3977(root)
}

// Java method: Builder.traverseAtScope(com.google.javascript.jscomp.AbstractScope<?,?> scope)
func (this_ *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	t_ := acquireNodeTraversal_(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverseAtScope_3981(scope)
}

// Java method: Builder.traverseRoots(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	t_ := acquireNodeTraversal_(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverseRoots_3978(externs, root)
}

// Java method: Builder.traverseWithScope(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.AbstractScope<?,?> s)
func (this_ *NodeTraversal_Builder) TraverseWithScope(root pkg_rhino.Node, s AbstractScope) {
	t_ := acquireNodeTraversal_(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder)
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverseWithScope_3980(root, s)
}

// Java class: TraverseScopeRootsCallback
type NodeTraversal__Local_65_TraverseScopeRootsCallback struct {
	cap_1_cb             NodeTraversal_Callback
	scopeNode            pkg_rhino.Node `gson:"scopeNode"`
	idHash_              atomic.Int32
	cap_0_traverseNested bool
	insideScopeNode      bool `gson:"insideScopeNode"`
}

func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) IsJavaNil_() bool {
	return this_ == nil
}
func (*NodeTraversal__Local_65_TraverseScopeRootsCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$1TraverseScopeRootsCallback", "TraverseScopeRootsCallback"
}
func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeTraversal__Local_65_TraverseScopeRootsCallback) HasIdentityEquals_() {}
func (*NodeTraversal__Local_65_TraverseScopeRootsCallback) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*NodeTraversal__Local_65_TraverseScopeRootsCallback) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) initFields_NodeTraversal__Local_65_TraverseScopeRootsCallback() {
	_ = this_
	this_.insideScopeNode = false
	this_.scopeNode = pkg_rhino.Node(0)
}

// Java constructor: TraverseScopeRootsCallback()
func NewNodeTraversal__Local_65_TraverseScopeRootsCallback(cap_p_traverseNested bool, cap_p_cb NodeTraversal_Callback) *NodeTraversal__Local_65_TraverseScopeRootsCallback {
	this_ := &NodeTraversal__Local_65_TraverseScopeRootsCallback{}
	this_.cap_0_traverseNested = cap_p_traverseNested
	this_.cap_1_cb = cap_p_cb
	InitNodeTraversal__Local_65_TraverseScopeRootsCallback(this_)
	return this_
}

func InitNodeTraversal__Local_65_TraverseScopeRootsCallback(this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) {
	this_.initFields_NodeTraversal__Local_65_TraverseScopeRootsCallback()
}

// Java method: TraverseScopeRootsCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if this_.scopeNode == n {
		this_.insideScopeNode = true
	}
	return (((this_.cap_0_traverseNested || (this_.scopeNode == n)) || !(ChangeTracker_IsChangeScopeRoot(n))) && this_.cap_1_cb.ShouldTraverse_3(t, n, parent))
}

// Java method: TraverseScopeRootsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.scopeNode == n {
		this_.insideScopeNode = false
	}
	this_.cap_1_cb.Visit_3(t, n, parent)
}

// Java method: TraverseScopeRootsCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) EnterScope(t *NodeTraversal) {
	_ = t
	var scopedCallback NodeTraversal_ScopedCallback
	_ = scopedCallback
	ok_35 := this_.insideScopeNode
	if ok_35 {
		scopedCallback, ok_35 = j2gort.TryCast[NodeTraversal_ScopedCallback](this_.cap_1_cb)
	}
	if ok_35 {
		scopedCallback.EnterScope(t)
	}
}

// Java method: TraverseScopeRootsCallback.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *NodeTraversal__Local_65_TraverseScopeRootsCallback) ExitScope(t *NodeTraversal) {
	_ = t
	var scopedCallback NodeTraversal_ScopedCallback
	_ = scopedCallback
	ok_36 := this_.insideScopeNode
	if ok_36 {
		scopedCallback, ok_36 = j2gort.TryCast[NodeTraversal_ScopedCallback](this_.cap_1_cb)
	}
	if ok_36 {
		scopedCallback.ExitScope(t)
	}
}

// Java class:
type NodeTraversal__Anon_64 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_  *NodeTraversal_Builder
	cap_0_x NodeTraversal_AbstractPostOrderCallbackInterface
}

func (this_ *NodeTraversal__Anon_64) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeTraversal__Anon_64) IsJavaNil_() bool { return this_ == nil }
func (*NodeTraversal__Anon_64) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.NodeTraversal$Builder$1", ""
}
func (this_ *NodeTraversal__Anon_64) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*NodeTraversal__Anon_64) HasIdentityEquals_()                                        {}
func (this_ *NodeTraversal__Anon_64) SetOuter_(o *NodeTraversal_Builder)                   { this_.outer_ = o }
func (this_ *NodeTraversal__Anon_64) GetOuter_() *NodeTraversal_Builder                    { return this_.outer_ }
func (*NodeTraversal__Anon_64) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*NodeTraversal__Anon_64) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *NodeTraversal__Anon_64) initFields_NodeTraversal__Anon_64() {
	_ = this_
}

func NewNodeTraversal__Anon_64(outer_ *NodeTraversal_Builder, cap_p_x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal__Anon_64 {
	this_ := &NodeTraversal__Anon_64{}
	this_.outer_ = outer_
	this_.cap_0_x = cap_p_x
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_NodeTraversal__Anon_64()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *NodeTraversal__Anon_64) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.cap_0_x.Visit(t, n, parent)
}
func (*NodeTraversal__Anon_64) AlwaysShouldTraverse_() {}
func (this_ *NodeTraversal__Anon_64) UnwrapCallback_() any {
	if this_ == nil {
		return nil
	}
	return this_.cap_0_x
}

type releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface struct {
	fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func (l *releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) IsJavaNil_() bool {
	return l == nil
}
func (*releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {
}
func (*releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) CanReleaseSyntacticScopes_() {
}
func (l *releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface) Visit(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	l.fn(p0_, p1_, p2_)
}
func newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface(fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)) NodeTraversal_AbstractPostOrderCallbackInterface {
	if fn == nil {
		return nil
	}
	return &releasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface{fn: fn}
}
func unwrapNodeTraversalCallback_(cb NodeTraversal_Callback) any {
	if uw_, ok_ := any(cb).(interface{ UnwrapCallback_() any }); ok_ {
		return uw_.UnwrapCallback_()
	}
	return cb
}
func canReleaseSyntacticScopesForCallback_(cb NodeTraversal_Callback) bool {
	u_ := unwrapNodeTraversalCallback_(cb)
	if _, ok_ := u_.(interface{ CanReleaseSyntacticScopes_() }); ok_ {
		return true
	}
	if d_, ok_ := u_.(interface{ CanReleaseSyntacticScopesDynamic_() bool }); ok_ {
		return d_.CanReleaseSyntacticScopesDynamic_()
	}
	return false
}

var nodeTraversalPool_ = sync.Pool{
	New: func() any { return &NodeTraversal{} },
}

func acquireNodeTraversal_(compiler AbstractCompiler, callback NodeTraversal_Callback, scopeCreator ScopeCreator, obey bool) *NodeTraversal {
	clinit_NodeTraversal()
	t_ := nodeTraversalPool_.Get().(*NodeTraversal)
	savedNodes_ := t_.scopeNodes
	savedObjs_ := t_.scopeObjs
	b_ := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     callback,
		scopeCreator: scopeCreator,
		obeyDestructuringAndDefaultValueExecutionOrder: obey,
	}
	InitNodeTraversal(t_, &b_)
	if cap(savedNodes_) > len(t_.scopeNodesBuf_) {
		t_.scopeNodes = savedNodes_[:cap(savedNodes_)]
	} else {
		t_.scopeNodes = t_.scopeNodesBuf_[:len(t_.scopeNodesBuf_)]
	}
	if cap(savedObjs_) > len(t_.scopeObjsBuf_) {
		t_.scopeObjs = savedObjs_[:cap(savedObjs_)]
	} else {
		t_.scopeObjs = t_.scopeObjsBuf_[:len(t_.scopeObjsBuf_)]
	}
	return t_
}

func releaseNodeTraversal_(t_ *NodeTraversal) {
	savedNodes_ := t_.scopeNodes
	savedObjs_ := t_.scopeObjs
	if cap(savedNodes_) > len(t_.scopeNodesBuf_) {
		clear(savedNodes_)
	}
	if cap(savedObjs_) > len(t_.scopeObjsBuf_) {
		clear(savedObjs_)
	}
	savedScanner_ := t_.defaultScopeCreator_.reusableScanner
	if savedScanner_ != nil {
		savedScanner_.compiler = nil
		savedScanner_.redeclarationHandler = nil
		savedScanner_.treatProvidesAsRedeclarations = false
		savedScanner_.Reset()
	}
	if t_.canReleaseScopes_ {
		t_.defaultScopeCreator_.ReleaseSlabs_()
	}
	*t_ = NodeTraversal{}
	t_.defaultScopeCreator_.reusableScanner = savedScanner_
	if cap(savedNodes_) > len(t_.scopeNodesBuf_) {
		t_.scopeNodes = savedNodes_[:cap(savedNodes_)]
	}
	if cap(savedObjs_) > len(t_.scopeObjsBuf_) {
		t_.scopeObjs = savedObjs_[:cap(savedObjs_)]
	}
	nodeTraversalPool_.Put(t_)
}
func traverseScopeRootPooled_(compiler AbstractCompiler, callback NodeTraversal_Callback, scopeCreator ScopeCreator, scopeRoot pkg_rhino.Node, canRel_ bool) {
	t_ := acquireNodeTraversal_(compiler, callback, scopeCreator, false)
	if canRel_ {
		t_.canReleaseScopes_ = true
	}
	defer releaseNodeTraversal_(t_)
	t_.priv_NodeTraversal_traverseScopeRoot_3982(scopeRoot)
}

var nodeTraversalBlockScopeTok_ = [256]bool{80: true, 81: true, 90: true, 91: true, 105: true, 112: true, 113: true}

var nodeTraversalSpecialOrBlockTok_ = [256]bool{73: true, 80: true, 81: true, 90: true, 91: true, 94: true, 104: true, 105: true, 106: true, 112: true, 113: true, 123: true, 134: true, 255: true}

func (this_ *NodeTraversal) traverseBranchBlk_(n, parent, nFirst pkg_rhino.Node, nTok uint8, curBlkIdx_ uint32, curBlk_ *nodestore.NodeBlock) {

	_ = n
	_ = parent
	if nTok == 0xFF {
		_ = n.TokenSwitchOrd_()
	}

	switch int32(nTok) {
	case int32(94):
		this_.priv_NodeTraversal_handleScript_3988(n, parent)
		return
	case int32(73):
		this_.priv_NodeTraversal_handleFunction_3989(n, parent)
		return
	case int32(123):
		this_.priv_NodeTraversal_handleModule_3990(n, parent)
		return
	case int32(105):
		this_.priv_NodeTraversal_handleClass_3995(n, parent)
		return
	case int32(106):
		this_.priv_NodeTraversal_handleClassMembers_3996(n, parent)
		return
	case int32(134), int32(104):
		if this_.obeyDestructuringAndDefaultValueExecutionOrder {
			this_.priv_NodeTraversal_handleDestructuringOrDefaultValue_3991(n, parent)
			return
		}
	default:
	}
	this_.currentNode = n
	if !this_.
		alwaysTraverse_ {

		if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
			return
		}
	}
	var createsBlockScope bool
	if nodeTraversalBlockScopeTok_[nTok] {
		createsBlockScope = this_.priv_NodeTraversal_createsBlockScope_3992(n)
	}

	_ = createsBlockScope
	var previousHoistScopeRoot pkg_rhino.Node = this_.currentHoistScopeRoot
	_ = previousHoistScopeRoot
	if createsBlockScope {
		this_.priv_NodeTraversal_pushScope_4001(n)
		if NodeUtil_IsClassStaticBlock(n) {
			this_.currentHoistScopeRoot = n
		}
	}
	var child_ pkg_rhino.Node
	if this_.alwaysTraverse_ && !createsBlockScope {
		child_ = nFirst
	} else if nIdx_ := uint32(n) >> nodestore.BlockShift; nIdx_ == curBlkIdx_ {
		f_, _, _ := curBlk_.FirstNextTok(nodestore.Node(n))
		child_ = pkg_rhino.Node(f_)
	} else {
		curBlk_ = nodestore.BlockOf(nodestore.Node(n))
		curBlkIdx_ = nIdx_
		f_, _, _ := curBlk_.FirstNextTok(nodestore.Node(n))
		child_ = pkg_rhino.Node(f_)
	}
	for child_ != 0 {
		var cFirst_, next_ pkg_rhino.Node
		var cTok_ uint8
		if idx_ := uint32(child_) >> nodestore.BlockShift; idx_ == curBlkIdx_ {
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		} else {
			curBlk_ = nodestore.BlockOf(nodestore.Node(child_))
			curBlkIdx_ = idx_
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		}
		if cFirst_ == 0 && cTok_ != 0xFF && !nodeTraversalSpecialOrBlockTok_[cTok_] {
			if this_.alwaysTraverse_ {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			this_.currentNode = child_
			if !this_.callback.ShouldTraverse_3(this_, child_, n) {
				child_ = next_
				continue
			}
			var gc_ pkg_rhino.Node
			if cIdx_ := uint32(child_) >> nodestore.BlockShift; cIdx_ == curBlkIdx_ {
				gf_, _, _ := curBlk_.FirstNextTok(nodestore.Node(child_))
				gc_ = pkg_rhino.Node(gf_)
			} else {
				curBlk_ = nodestore.BlockOf(nodestore.Node(child_))
				curBlkIdx_ = cIdx_
				gf_, _, _ := curBlk_.FirstNextTok(nodestore.Node(child_))
				gc_ = pkg_rhino.Node(gf_)
			}
			if gc_ == 0 {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			for gc_ != 0 {
				var gFirst_, gNext_ pkg_rhino.Node
				var gTok_ uint8
				if gIdx_ := uint32(gc_) >> nodestore.BlockShift; gIdx_ == curBlkIdx_ {
					gf_, gnx_, gt_ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gFirst_, gNext_, gTok_ = pkg_rhino.Node(gf_), pkg_rhino.Node(gnx_), gt_
				} else {
					curBlk_ = nodestore.BlockOf(nodestore.Node(gc_))
					curBlkIdx_ = gIdx_
					gf_, gnx_, gt_ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gFirst_, gNext_, gTok_ = pkg_rhino.Node(gf_), pkg_rhino.Node(gnx_), gt_
				}
				this_.traverseBranchBlk_(gc_, child_, gFirst_, gTok_, curBlkIdx_, curBlk_)
				gc_ = gNext_
			}
			this_.currentNode = child_
			if !this_.noOpVisit_ {
				this_.callback.Visit_3(this_, child_, n)
			}
			child_ = next_
			continue
		}
		this_.traverseBranchBlk_(child_, n, cFirst_, cTok_, curBlkIdx_, curBlk_)
		child_ = next_
	}

	if createsBlockScope {
		this_.priv_NodeTraversal_popScope_4004()
		this_.currentHoistScopeRoot = previousHoistScopeRoot
	}
	this_.currentNode = n
	if !this_.
		noOpVisit_ {

		this_.callback.Visit_3(this_, n, parent)
	}
}
