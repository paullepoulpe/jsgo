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
// Port of com/google/javascript/jscomp/CheckSuper.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CheckSuper
type CheckSuper struct {
	compiler     AbstractCompiler                 `gson:"compiler"`
	contextStack *j2gort.Ext_java_util_ArrayDeque `gson:"contextStack"`
	idHash_      atomic.Int32
}

func (this_ *CheckSuper) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper", "CheckSuper"
}
func (this_ *CheckSuper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckSuper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper", "com.google.javascript.jscomp.CheckSuper", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuper](v_)
	})
}
func (*CheckSuper) Implements_jscomp_CompilerPass()           {}
func (*CheckSuper) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckSuper_MISSING_CALL_TO_SUPER *DiagnosticType

func StaticGet_CheckSuper_MISSING_CALL_TO_SUPER() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_MISSING_CALL_TO_SUPER
}
func StaticSet_CheckSuper_MISSING_CALL_TO_SUPER(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_MISSING_CALL_TO_SUPER = v
}

var static_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER *DiagnosticType

func StaticGet_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER
}
func StaticSet_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER = v
}

var static_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER *DiagnosticType

func StaticGet_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER
}
func StaticSet_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER = v
}

var static_CheckSuper_THIS_BEFORE_SUPER *DiagnosticType

func StaticGet_CheckSuper_THIS_BEFORE_SUPER() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_THIS_BEFORE_SUPER
}
func StaticSet_CheckSuper_THIS_BEFORE_SUPER(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_THIS_BEFORE_SUPER = v
}

var static_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR *DiagnosticType

func StaticGet_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR
}
func StaticSet_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR = v
}

var static_CheckSuper_INVALID_SUPER_CALL *DiagnosticType

func StaticGet_CheckSuper_INVALID_SUPER_CALL() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_INVALID_SUPER_CALL
}
func StaticSet_CheckSuper_INVALID_SUPER_CALL(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_INVALID_SUPER_CALL = v
}

var static_CheckSuper_SUPER_CALL_IN_ARROW *DiagnosticType

func StaticGet_CheckSuper_SUPER_CALL_IN_ARROW() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_SUPER_CALL_IN_ARROW
}
func StaticSet_CheckSuper_SUPER_CALL_IN_ARROW(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_SUPER_CALL_IN_ARROW = v
}

var static_CheckSuper_INVALID_SUPER_USAGE *DiagnosticType

func StaticGet_CheckSuper_INVALID_SUPER_USAGE() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_INVALID_SUPER_USAGE
}
func StaticSet_CheckSuper_INVALID_SUPER_USAGE(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_INVALID_SUPER_USAGE = v
}

var static_CheckSuper_INVALID_SUPER_ACCESS *DiagnosticType

func StaticGet_CheckSuper_INVALID_SUPER_ACCESS() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_INVALID_SUPER_ACCESS
}
func StaticSet_CheckSuper_INVALID_SUPER_ACCESS(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_INVALID_SUPER_ACCESS = v
}

var static_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION *DiagnosticType

func StaticGet_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION() *DiagnosticType {
	clinit_CheckSuper()
	return static_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION
}
func StaticSet_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION(v *DiagnosticType) {
	clinit_CheckSuper()
	static_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION = v
}

var clinitDone_CheckSuper uint32
var clinitState_CheckSuper j2gort.ClinitState

func clinit_CheckSuper() {
	if atomic.LoadUint32(&clinitDone_CheckSuper) != 0 {
		return
	}
	clinitSlow_CheckSuper()
}

func clinitSlow_CheckSuper() {
	if !j2gort.EnterClinit(&clinitState_CheckSuper, &clinitDone_CheckSuper) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckSuper, &clinitDone_CheckSuper)
	static_CheckSuper_MISSING_CALL_TO_SUPER = DiagnosticType_Error("JSC_MISSING_CALL_TO_SUPER", "constructor is missing a call to super()")
	static_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER = DiagnosticType_Error("JSC_RETURN_WITHOUT_CALL_TO_SUPER", "constructor returns a value without calling super()")
	static_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER = DiagnosticType_Error("JSC_MISSING_ROOT_LEVEL_CALL_TO_SUPER", "constructor of a class with instance fields or private members must call super() unconditionally once in a root-level statement before any return")
	static_CheckSuper_THIS_BEFORE_SUPER = DiagnosticType_Error("JSC_THIS_BEFORE_SUPER", "cannot access this before calling super()")
	static_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR = DiagnosticType_Error("JSC_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR", "cannot access super properties before calling super()")
	static_CheckSuper_INVALID_SUPER_CALL = DiagnosticType_Error("JSC_INVALID_SUPER_CALL", "super() not allowed except in the constructor of a subclass")
	static_CheckSuper_SUPER_CALL_IN_ARROW = DiagnosticType_Error("JSC_SUPER_CALL_IN_ARROW", "closure-compiler does not allow calls to `super()` in arrow functions")
	static_CheckSuper_INVALID_SUPER_USAGE = DiagnosticType_Error("JSC_INVALID_SUPER_USAGE", "''super'' may only be used in a call or property access")
	static_CheckSuper_INVALID_SUPER_ACCESS = DiagnosticType_Error("JSC_INVALID_SUPER_ACCESS", "''super'' may only be accessed within a method")
	static_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION = DiagnosticType_Error("JSC_INVALID_SUPER_CALL_WITH_SUGGESTION", "super() not allowed here. Did you mean super.{0}?")
}

func (this_ *CheckSuper) initFields_CheckSuper() {
	_ = this_
	this_.contextStack = j2gort.NewExt_java_util_ArrayDeque()
}

// Java constructor: CheckSuper(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckSuper(compiler AbstractCompiler) *CheckSuper {
	clinit_CheckSuper()
	this_ := &CheckSuper{}
	InitCheckSuper(this_, compiler)
	return this_
}

func InitCheckSuper(this_ *CheckSuper, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckSuper()
	this_.compiler = compiler
}

// Java method: CheckSuper.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckSuper) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	j2gort.CheckStateMsg(this_.contextStack.IsEmpty(), j2gort.ImmListOfPtr(this_.contextStack))
}

// Java method: CheckSuper.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckSuper) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	{
		var fieldContext CheckSuper_Context = nil
		_ = fieldContext
		switch n.TokenSwitchOrd_() {
		case int32(89) /* Token.ROOT */ :
			if !(this_.contextStack.IsEmpty()) {
				j2gort.CheckState(false)
			}
			this_.contextStack.Push(NewCheckSuper_SuperNotAllowedContext(n))
		case int32(73) /* Token.FUNCTION */ :
			{
				var currentContext CheckSuper_Context = j2gort.AsIface[CheckSuper_Context](this_.contextStack.Peek())
				_ = currentContext
				var newContext CheckSuper_Context = this_.priv_CheckSuper_getContextForFunctionNode_1108(currentContext, n)
				_ = newContext
				if !(newContext == currentContext || ((newContext == nil || newContext.IsJavaNil_()) && (currentContext == nil || currentContext.IsJavaNil_()))) {
					this_.contextStack.Push(newContext)
				}
			}
		case int32(108) /* Token.MEMBER_FIELD_DEF */ :
			fieldContext = this_.priv_CheckSuper_createMemberFieldDefContext_1111(n)
			this_.contextStack.Push(fieldContext)
		case int32(90) /* Token.BLOCK */ :
			if NodeUtil_IsClassStaticBlock(n) {
				var newContext_2 CheckSuper_Context = this_.priv_CheckSuper_createStaticBlockContext_1112(n)
				_ = newContext_2
				this_.contextStack.Push(newContext_2)
			}
		case int32(105) /* Token.CLASS */ :
			fallthrough
		default:
		}
	}
	return true
}

// Java method: CheckSuper.getContextForFunctionNode(com.google.javascript.jscomp.CheckSuper.Context currentContext, com.google.javascript.rhino.Node fn)
func (this_ *CheckSuper) priv_CheckSuper_getContextForFunctionNode_1108(currentContext CheckSuper_Context, fn pkg_rhino.Node) CheckSuper_Context {
	_ = currentContext
	_ = fn
	if NodeUtil_IsMethodDeclaration(fn) {
		if NodeUtil_IsEs6Constructor(fn) {
			return NewCheckSuper_ConstructorContext(fn)
		} else {
			return NewCheckSuper_MethodContext(fn)
		}
	} else {
		if fn.IsArrowFunction() {
			return currentContext.GetContextForArrowFunctionNode(fn)
		} else {
			return NewCheckSuper_SuperNotAllowedContext(fn)
		}
	}
}

// Java method: CheckSuper.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckSuper) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var currentContext CheckSuper_Context = j2gort.CheckNotNullIface(j2gort.AsIface[CheckSuper_Context](this_.contextStack.Peek()))
	_ = currentContext
	switch n.TokenSwitchOrd_() {
	case int32(110) /* Token.SUPER */ :
		if this_.priv_CheckSuper_isSuperConstructorCall_1109(n) {
		} else if this_.priv_CheckSuper_isSuperPropertyAccess_1110(n) {
			currentContext.VisitSuperPropertyAccess(t, n)
		} else {
			t.Report_3(n, StaticGet_CheckSuper_INVALID_SUPER_USAGE(), []string{})
		}
	case int32(28) /* Token.CALL */ :
		var callee pkg_rhino.Node = n.GetFirstChild()
		_ = callee
		if callee.IsSuper() {
			currentContext.VisitSuperConstructorCall(t, callee)
		}
	case int32(37) /* Token.THIS */ :
		currentContext.VisitThis(t, n)
	case int32(0) /* Token.RETURN */ :
		currentContext.VisitReturn(t, n)
	case int32(89) /* Token.ROOT */ :
		if !(this_.contextStack.Size() == int32(1)) {
			j2gort.CheckState(false)
		}
	default:
	}
	if n == currentContext.GetCheckSuper_ContextBase().GetContextNode() {
		currentContext.VisitContextNode(t)
		this_.contextStack.Pop()
	}
}

// Java method: CheckSuper.isSuperConstructorCall(com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper) priv_CheckSuper_isSuperConstructorCall_1109(superNode pkg_rhino.Node) bool {
	_ = superNode
	if !(superNode.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(superNode))
	}
	var parent pkg_rhino.Node = superNode.GetParent()
	_ = parent
	return (parent.IsCall() && superNode.IsFirstChildOf(parent))
}

// Java method: CheckSuper.isSuperPropertyAccess(com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper) priv_CheckSuper_isSuperPropertyAccess_1110(superNode pkg_rhino.Node) bool {
	_ = superNode
	if !(superNode.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(superNode))
	}
	var parent pkg_rhino.Node = superNode.GetParent()
	_ = parent
	return (NodeUtil_IsNormalGet(parent) && superNode.IsFirstChildOf(parent))
}

// Java method: CheckSuper.createMemberFieldDefContext(com.google.javascript.rhino.Node fieldNode)
func (this_ *CheckSuper) priv_CheckSuper_createMemberFieldDefContext_1111(fieldNode pkg_rhino.Node) CheckSuper_Context {
	_ = fieldNode
	return NewCheckSuper_SuperPropertyAccessAllowedContext(fieldNode)
}

// Java method: CheckSuper.createStaticBlockContext(com.google.javascript.rhino.Node staticBlock)
func (this_ *CheckSuper) priv_CheckSuper_createStaticBlockContext_1112(staticBlock pkg_rhino.Node) CheckSuper_Context {
	_ = staticBlock
	return NewCheckSuper_SuperPropertyAccessAllowedContext(staticBlock)
}

// Java polymorphic class interface: com.google.javascript.jscomp.CheckSuper.Context
type CheckSuper_Context interface {
	Implements_jscomp_CheckSuper_Context()
	IsJavaNil_() bool
	GetCheckSuper_ContextBase() *CheckSuper_ContextBase
	GetContextNode() pkg_rhino.Node
	GetContextForArrowFunctionNode(pkg_rhino.Node) CheckSuper_Context
	VisitSuperConstructorCall(*NodeTraversal, pkg_rhino.Node)
	VisitSuperPropertyAccess(*NodeTraversal, pkg_rhino.Node)
	VisitThis(*NodeTraversal, pkg_rhino.Node)
	VisitReturn(*NodeTraversal, pkg_rhino.Node)
	VisitContextNode(*NodeTraversal)
}

func ToIface_CheckSuper_Context[T interface {
	CheckSuper_Context
	comparable
}](v T) CheckSuper_Context {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CheckSuper_Context[T CheckSuper_Context](v T) CheckSuper_Context {
	var vi CheckSuper_Context = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.CheckSuper.Context
type CheckSuper_ContextBase struct {
	contextNode pkg_rhino.Node `gson:"contextNode"`
	idHash_     atomic.Int32
}

func (this_ *CheckSuper_ContextBase) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper_ContextBase) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper_ContextBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper$Context", "Context"
}
func (this_ *CheckSuper_ContextBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper.Context", "com.google.javascript.jscomp.CheckSuper$Context", func(v_ any) bool {
		return j2gort.InstanceOf[CheckSuper_Context](v_)
	})
}
func (b *CheckSuper_ContextBase) GetCheckSuper_ContextBase() *CheckSuper_ContextBase { return b }
func (*CheckSuper_ContextBase) Implements_jscomp_CheckSuper_Context()                {}

func (this_ *CheckSuper_ContextBase) initFields_CheckSuper_Context() {
	_ = this_
}

func InitCheckSuper_Context(this_ *CheckSuper_ContextBase, contextNode pkg_rhino.Node) {
	_ = contextNode
	this_.initFields_CheckSuper_Context()
	this_.contextNode = j2gort.CheckNotNullNode(contextNode)
}

// Java method: Context.getContextNode()
func (this_ *CheckSuper_ContextBase) GetContextNode() pkg_rhino.Node {
	return this_.contextNode
}

// Java method: Context.getContextForArrowFunctionNode(com.google.javascript.rhino.Node arrowFn)
func (this_ *CheckSuper_ContextBase) GetContextForArrowFunctionNode(arrowFn pkg_rhino.Node) CheckSuper_Context {
	_ = arrowFn
	panic("abstract or native method invoked")
}

// Java method: Context.visitSuperConstructorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_ContextBase) VisitSuperConstructorCall(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	panic("abstract or native method invoked")
}

// Java method: Context.visitSuperPropertyAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_ContextBase) VisitSuperPropertyAccess(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	panic("abstract or native method invoked")
}

// Java method: Context.visitThis(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node thisNode)
func (this_ *CheckSuper_ContextBase) VisitThis(t *NodeTraversal, thisNode pkg_rhino.Node) {
	_ = t
	_ = thisNode
}

// Java method: Context.visitReturn(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node returnNode)
func (this_ *CheckSuper_ContextBase) VisitReturn(t *NodeTraversal, returnNode pkg_rhino.Node) {
	_ = t
	_ = returnNode
}

// Java method: Context.visitContextNode(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckSuper_ContextBase) VisitContextNode(t *NodeTraversal) {
	_ = t
}

// Java class: com.google.javascript.jscomp.CheckSuper.SuperPropertyAccessAllowedContext
type CheckSuper_SuperPropertyAccessAllowedContext struct {
	CheckSuper_ContextBase
}

func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper_SuperPropertyAccessAllowedContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper$SuperPropertyAccessAllowedContext", "SuperPropertyAccessAllowedContext"
}
func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CheckSuper_ContextBase.IdentityHash_()
}
func (*CheckSuper_SuperPropertyAccessAllowedContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper.SuperPropertyAccessAllowedContext", "com.google.javascript.jscomp.CheckSuper$SuperPropertyAccessAllowedContext", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuper_SuperPropertyAccessAllowedContext](v_)
	})
}
func (*CheckSuper_SuperPropertyAccessAllowedContext) Implements_jscomp_CheckSuper_Context() {}

func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) initFields_CheckSuper_SuperPropertyAccessAllowedContext() {
	_ = this_
}

// Java constructor: SuperPropertyAccessAllowedContext(com.google.javascript.rhino.Node contextNode)
func NewCheckSuper_SuperPropertyAccessAllowedContext(contextNode pkg_rhino.Node) *CheckSuper_SuperPropertyAccessAllowedContext {
	this_ := &CheckSuper_SuperPropertyAccessAllowedContext{}
	InitCheckSuper_SuperPropertyAccessAllowedContext(this_, contextNode)
	return this_
}

func InitCheckSuper_SuperPropertyAccessAllowedContext(this_ *CheckSuper_SuperPropertyAccessAllowedContext, contextNode pkg_rhino.Node) {
	_ = contextNode
	InitCheckSuper_Context(&this_.CheckSuper_ContextBase, contextNode)
	this_.initFields_CheckSuper_SuperPropertyAccessAllowedContext()
}

// Java method: SuperPropertyAccessAllowedContext.visitSuperPropertyAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
}

// Java method: SuperPropertyAccessAllowedContext.visitSuperConstructorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) VisitSuperConstructorCall(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	t.Report_3(superNode, StaticGet_CheckSuper_INVALID_SUPER_CALL(), []string{})
}

// Java method: SuperPropertyAccessAllowedContext.getContextForArrowFunctionNode(com.google.javascript.rhino.Node arrowFn)
func (this_ *CheckSuper_SuperPropertyAccessAllowedContext) GetContextForArrowFunctionNode(arrowFn pkg_rhino.Node) CheckSuper_Context {
	_ = arrowFn
	return this_
}

// Java class: com.google.javascript.jscomp.CheckSuper.SuperNotAllowedContext
type CheckSuper_SuperNotAllowedContext struct {
	CheckSuper_ContextBase
}

func (this_ *CheckSuper_SuperNotAllowedContext) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper_SuperNotAllowedContext) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper_SuperNotAllowedContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper$SuperNotAllowedContext", "SuperNotAllowedContext"
}
func (this_ *CheckSuper_SuperNotAllowedContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CheckSuper_ContextBase.IdentityHash_()
}
func (*CheckSuper_SuperNotAllowedContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper.SuperNotAllowedContext", "com.google.javascript.jscomp.CheckSuper$SuperNotAllowedContext", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuper_SuperNotAllowedContext](v_)
	})
}
func (*CheckSuper_SuperNotAllowedContext) Implements_jscomp_CheckSuper_Context() {}

func (this_ *CheckSuper_SuperNotAllowedContext) initFields_CheckSuper_SuperNotAllowedContext() {
	_ = this_
}

// Java constructor: SuperNotAllowedContext(com.google.javascript.rhino.Node contextNode)
func NewCheckSuper_SuperNotAllowedContext(contextNode pkg_rhino.Node) *CheckSuper_SuperNotAllowedContext {
	this_ := &CheckSuper_SuperNotAllowedContext{}
	InitCheckSuper_SuperNotAllowedContext(this_, contextNode)
	return this_
}

func InitCheckSuper_SuperNotAllowedContext(this_ *CheckSuper_SuperNotAllowedContext, contextNode pkg_rhino.Node) {
	_ = contextNode
	InitCheckSuper_Context(&this_.CheckSuper_ContextBase, contextNode)
	this_.initFields_CheckSuper_SuperNotAllowedContext()
}

// Java method: SuperNotAllowedContext.getContextForArrowFunctionNode(com.google.javascript.rhino.Node arrowFn)
func (this_ *CheckSuper_SuperNotAllowedContext) GetContextForArrowFunctionNode(arrowFn pkg_rhino.Node) CheckSuper_Context {
	_ = arrowFn
	return this_
}

// Java method: SuperNotAllowedContext.visitSuperConstructorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_SuperNotAllowedContext) VisitSuperConstructorCall(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	t.Report_3(superNode, StaticGet_CheckSuper_INVALID_SUPER_CALL(), []string{})
}

// Java method: SuperNotAllowedContext.visitSuperPropertyAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_SuperNotAllowedContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	t.Report_3(superNode, StaticGet_CheckSuper_INVALID_SUPER_ACCESS(), []string{})
}

// Java class: com.google.javascript.jscomp.CheckSuper.MethodContext
type CheckSuper_MethodContext struct {
	CheckSuper_ContextBase
	isClassMethod bool `gson:"isClassMethod"`
}

func (this_ *CheckSuper_MethodContext) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper_MethodContext) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper_MethodContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper$MethodContext", "MethodContext"
}
func (this_ *CheckSuper_MethodContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CheckSuper_ContextBase.IdentityHash_()
}
func (*CheckSuper_MethodContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper.MethodContext", "com.google.javascript.jscomp.CheckSuper$MethodContext", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuper_MethodContext](v_)
	})
}
func (*CheckSuper_MethodContext) Implements_jscomp_CheckSuper_Context() {}

func (this_ *CheckSuper_MethodContext) initFields_CheckSuper_MethodContext() {
	_ = this_
}

// Java constructor: MethodContext(com.google.javascript.rhino.Node functionNode)
func NewCheckSuper_MethodContext(functionNode pkg_rhino.Node) *CheckSuper_MethodContext {
	this_ := &CheckSuper_MethodContext{}
	InitCheckSuper_MethodContext(this_, functionNode)
	return this_
}

func InitCheckSuper_MethodContext(this_ *CheckSuper_MethodContext, functionNode pkg_rhino.Node) {
	_ = functionNode
	InitCheckSuper_Context(&this_.CheckSuper_ContextBase, functionNode)
	this_.initFields_CheckSuper_MethodContext()
	if !(functionNode.IsFunction() && NodeUtil_IsMethodDeclaration(functionNode)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(functionNode))
	}
	var objLitOrClassMembers pkg_rhino.Node = j2gort.CheckNotNullNode(functionNode.GetGrandparent())
	_ = objLitOrClassMembers
	if objLitOrClassMembers.IsObjectLit() {
		this_.isClassMethod = false
	} else {
		if !(objLitOrClassMembers.IsClassMembers()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(objLitOrClassMembers))
		}
		this_.isClassMethod = true
	}
}

// Java method: MethodContext.getContextForArrowFunctionNode(com.google.javascript.rhino.Node arrowFn)
func (this_ *CheckSuper_MethodContext) GetContextForArrowFunctionNode(arrowFn pkg_rhino.Node) CheckSuper_Context {
	_ = arrowFn
	return this_
}

// Java method: MethodContext.visitSuperPropertyAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_MethodContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
}

// Java method: MethodContext.visitSuperConstructorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_MethodContext) VisitSuperConstructorCall(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	if this_.isClassMethod {
		var propName string = this_.GetPropertyName()
		_ = propName
		if propName != j2gort.NullString {
			t.Report_3(superNode, StaticGet_CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION(), []string{propName})
		} else {
			t.Report_3(superNode, StaticGet_CheckSuper_INVALID_SUPER_CALL(), []string{})
		}
	} else {
		t.Report_3(superNode, StaticGet_CheckSuper_INVALID_SUPER_CALL(), []string{})
	}
}

// Java method: MethodContext.getPropertyName()
func (this_ *CheckSuper_MethodContext) GetPropertyName() string {
	var parent pkg_rhino.Node = j2gort.CheckNotNullNode(this_.GetContextNode().GetParent())
	_ = parent
	if (parent.IsMemberFunctionDef() || parent.IsGetterDef()) || parent.IsSetterDef() {
		return parent.GetString()
	} else {
		return j2gort.NullString
	}
}

// Java class: com.google.javascript.jscomp.CheckSuper.ConstructorContext
type CheckSuper_ConstructorContext struct {
	CheckSuper_ContextBase
	firstSuperCall                       pkg_rhino.Node `gson:"firstSuperCall"`
	rootLevelSuperCall                   pkg_rhino.Node `gson:"rootLevelSuperCall"`
	thisAccessedBeforeSuper              pkg_rhino.Node `gson:"thisAccessedBeforeSuper"`
	superPropertyAccessedBeforeSuperCall pkg_rhino.Node `gson:"superPropertyAccessedBeforeSuperCall"`
	hasParentClass                       bool           `gson:"hasParentClass"`
	hasInstanceElements                  bool           `gson:"hasInstanceElements"`
	hasMultipleSuperCalls                bool           `gson:"hasMultipleSuperCalls"`
	returnsAValue                        bool           `gson:"returnsAValue"`
	returnsBeforeRootLevelSuperCall      bool           `gson:"returnsBeforeRootLevelSuperCall"`
}

func (this_ *CheckSuper_ConstructorContext) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper_ConstructorContext) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper_ConstructorContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper$ConstructorContext", "ConstructorContext"
}
func (this_ *CheckSuper_ConstructorContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CheckSuper_ContextBase.IdentityHash_()
}
func (*CheckSuper_ConstructorContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper.ConstructorContext", "com.google.javascript.jscomp.CheckSuper$ConstructorContext", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuper_ConstructorContext](v_)
	})
}
func (*CheckSuper_ConstructorContext) Implements_jscomp_CheckSuper_Context() {}

func (this_ *CheckSuper_ConstructorContext) initFields_CheckSuper_ConstructorContext() {
	_ = this_

}

// Java constructor: ConstructorContext(com.google.javascript.rhino.Node contextNode)
func NewCheckSuper_ConstructorContext(contextNode pkg_rhino.Node) *CheckSuper_ConstructorContext {
	this_ := &CheckSuper_ConstructorContext{}
	InitCheckSuper_ConstructorContext(this_, contextNode)
	return this_
}

func InitCheckSuper_ConstructorContext(this_ *CheckSuper_ConstructorContext, contextNode pkg_rhino.Node) {
	_ = contextNode
	InitCheckSuper_Context(&this_.CheckSuper_ContextBase, contextNode)
	this_.initFields_CheckSuper_ConstructorContext()
	if !(NodeUtil_IsEs6Constructor(contextNode)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(contextNode))
	}
	var classNode pkg_rhino.Node = NodeUtil_GetEnclosingClass(contextNode)
	_ = classNode
	this_.hasParentClass = !(classNode.GetSecondChild().IsEmpty())
	this_.hasInstanceElements = NodeUtil_HasInstanceElements(classNode)
}

// Java method: ConstructorContext.getContextForArrowFunctionNode(com.google.javascript.rhino.Node arrowFn)
func (this_ *CheckSuper_ConstructorContext) GetContextForArrowFunctionNode(arrowFn pkg_rhino.Node) CheckSuper_Context {
	_ = arrowFn
	return NewCheckSuper_ConstructorArrowContext(this_, arrowFn)
}

// Java method: ConstructorContext.visitSuperConstructorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_ConstructorContext) VisitSuperConstructorCall(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	if this_.firstSuperCall == 0 {
		this_.firstSuperCall = superNode
	} else {
		this_.hasMultipleSuperCalls = true
	}
	var call pkg_rhino.Node = superNode.GetParent()
	_ = call
	if ((this_.rootLevelSuperCall == 0) && func() bool {
		rn_1 := NodeUtil_GetEnclosingStatement(call).GetParent()
		an_2 := NodeUtil_GetFunctionBody(this_.GetContextNode())
		return j2gort.CheckNotNullNode(rn_1) == an_2
	}()) && NodeUtil_IsExecutedExactlyOnce(call) {
		this_.rootLevelSuperCall = superNode
	}
}

// Java method: ConstructorContext.visitSuperPropertyAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_ConstructorContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	if this_.firstSuperCall == 0 {
		this_.superPropertyAccessedBeforeSuperCall = superNode
	}
}

// Java method: ConstructorContext.visitThis(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node thisNode)
func (this_ *CheckSuper_ConstructorContext) VisitThis(t *NodeTraversal, thisNode pkg_rhino.Node) {
	_ = t
	_ = thisNode
	if this_.firstSuperCall == 0 {
		this_.thisAccessedBeforeSuper = thisNode
	}
}

// Java method: ConstructorContext.visitReturn(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node returnNode)
func (this_ *CheckSuper_ConstructorContext) VisitReturn(t *NodeTraversal, returnNode pkg_rhino.Node) {
	_ = t
	_ = returnNode
	if returnNode.HasChildren() {
		this_.returnsAValue = true
	}
	if this_.rootLevelSuperCall == 0 {
		this_.returnsBeforeRootLevelSuperCall = true
	}
}

// Java method: ConstructorContext.visitContextNode(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckSuper_ConstructorContext) VisitContextNode(t *NodeTraversal) {
	_ = t
	if !(this_.hasParentClass) {
		if this_.firstSuperCall != 0 {
			t.Report_3(this_.firstSuperCall, StaticGet_CheckSuper_INVALID_SUPER_CALL(), []string{})
		}
	} else {
		if this_.firstSuperCall == 0 {
			t.Report_3(this_.GetContextNode().GetParent(), func() *DiagnosticType {
				if this_.returnsAValue {
					return StaticGet_CheckSuper_RETURN_WITHOUT_CALL_TO_SUPER()
				}
				return StaticGet_CheckSuper_MISSING_CALL_TO_SUPER()
			}(), []string{})
		} else if this_.hasInstanceElements && (((this_.rootLevelSuperCall == 0) || this_.returnsBeforeRootLevelSuperCall) || this_.hasMultipleSuperCalls) {
			t.Report_3(this_.GetContextNode().GetParent(), StaticGet_CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER(), []string{})
		}
		if this_.thisAccessedBeforeSuper != 0 {
			t.Report_3(this_.thisAccessedBeforeSuper, StaticGet_CheckSuper_THIS_BEFORE_SUPER(), []string{})
		}
		if this_.superPropertyAccessedBeforeSuperCall != 0 {
			t.Report_3(this_.superPropertyAccessedBeforeSuperCall, StaticGet_CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR(), []string{})
		}
	}
}

// Java class: com.google.javascript.jscomp.CheckSuper.ConstructorArrowContext
type CheckSuper_ConstructorArrowContext struct {
	CheckSuper_ContextBase
	constructorContext *CheckSuper_ConstructorContext `gson:"constructorContext"`
}

func (this_ *CheckSuper_ConstructorArrowContext) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSuper_ConstructorArrowContext) IsJavaNil_() bool { return this_ == nil }
func (*CheckSuper_ConstructorArrowContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSuper$ConstructorArrowContext", "ConstructorArrowContext"
}
func (this_ *CheckSuper_ConstructorArrowContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CheckSuper_ContextBase.IdentityHash_()
}
func (*CheckSuper_ConstructorArrowContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSuper.ConstructorArrowContext", "com.google.javascript.jscomp.CheckSuper$ConstructorArrowContext", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSuper_ConstructorArrowContext](v_)
	})
}
func (*CheckSuper_ConstructorArrowContext) Implements_jscomp_CheckSuper_Context() {}

func (this_ *CheckSuper_ConstructorArrowContext) initFields_CheckSuper_ConstructorArrowContext() {
	_ = this_
}

// Java constructor: ConstructorArrowContext(com.google.javascript.jscomp.CheckSuper.ConstructorContext constructorContext, com.google.javascript.rhino.Node arrowFn)
func NewCheckSuper_ConstructorArrowContext(constructorContext *CheckSuper_ConstructorContext, arrowFn pkg_rhino.Node) *CheckSuper_ConstructorArrowContext {
	this_ := &CheckSuper_ConstructorArrowContext{}
	InitCheckSuper_ConstructorArrowContext(this_, constructorContext, arrowFn)
	return this_
}

func InitCheckSuper_ConstructorArrowContext(this_ *CheckSuper_ConstructorArrowContext, constructorContext *CheckSuper_ConstructorContext, arrowFn pkg_rhino.Node) {
	_ = constructorContext
	_ = arrowFn
	InitCheckSuper_Context(&this_.CheckSuper_ContextBase, arrowFn)
	this_.initFields_CheckSuper_ConstructorArrowContext()
	this_.constructorContext = j2gort.CheckNotNullPtr(constructorContext)
}

// Java method: ConstructorArrowContext.getContextForArrowFunctionNode(com.google.javascript.rhino.Node arrowFn)
func (this_ *CheckSuper_ConstructorArrowContext) GetContextForArrowFunctionNode(arrowFn pkg_rhino.Node) CheckSuper_Context {
	_ = arrowFn
	return this_
}

// Java method: ConstructorArrowContext.visitSuperConstructorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_ConstructorArrowContext) VisitSuperConstructorCall(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	t.Report_3(superNode, StaticGet_CheckSuper_SUPER_CALL_IN_ARROW(), []string{})
	if this_.constructorContext.firstSuperCall == 0 {
		{
			r_3 := this_.constructorContext
			r_3.firstSuperCall = superNode
		}
	}
	if this_.constructorContext.rootLevelSuperCall == 0 {
		{
			r_4 := this_.constructorContext
			r_4.rootLevelSuperCall = superNode
		}
	}
}

// Java method: ConstructorArrowContext.visitSuperPropertyAccess(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node superNode)
func (this_ *CheckSuper_ConstructorArrowContext) VisitSuperPropertyAccess(t *NodeTraversal, superNode pkg_rhino.Node) {
	_ = t
	_ = superNode
	this_.constructorContext.VisitSuperPropertyAccess(t, superNode)
}
func (*CheckSuper) CanReleaseSyntacticScopes_() {
}
