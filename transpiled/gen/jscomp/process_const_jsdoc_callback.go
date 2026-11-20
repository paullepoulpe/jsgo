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
// Port of com/google/javascript/jscomp/ijs/ProcessConstJsdocCallback.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.ijs.ProcessConstJsdocCallback
type ProcessConstJsdocCallback interface {
	Implements_jscomp_ProcessConstJsdocCallback()
	IsJavaNil_() bool
	GetProcessConstJsdocCallbackBase() *ProcessConstJsdocCallbackBase
	NodeTraversal_AbstractPostOrderCallback
	ProcessConstWithRhs(*NodeTraversal, pkg_rhino.Node)
}

func ToIface_ProcessConstJsdocCallback[T interface {
	ProcessConstJsdocCallback
	comparable
}](v T) ProcessConstJsdocCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ProcessConstJsdocCallback[T ProcessConstJsdocCallback](v T) ProcessConstJsdocCallback {
	var vi ProcessConstJsdocCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.ijs.ProcessConstJsdocCallback
type ProcessConstJsdocCallbackBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	self_       ProcessConstJsdocCallback
	currentFile *FileInfo `gson:"currentFile"`
}

func (this_ *ProcessConstJsdocCallbackBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ProcessConstJsdocCallback)
}
func (this_ *ProcessConstJsdocCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*ProcessConstJsdocCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.ProcessConstJsdocCallback", "ProcessConstJsdocCallback"
}
func (this_ *ProcessConstJsdocCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ProcessConstJsdocCallback", "com.google.javascript.jscomp.ijs.ProcessConstJsdocCallback", func(v_ any) bool {
		return j2gort.InstanceOf[ProcessConstJsdocCallback](v_)
	})
}
func (b *ProcessConstJsdocCallbackBase) GetProcessConstJsdocCallbackBase() *ProcessConstJsdocCallbackBase {
	return b
}
func (*ProcessConstJsdocCallbackBase) Implements_jscomp_ProcessConstJsdocCallback()               {}
func (*ProcessConstJsdocCallbackBase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ProcessConstJsdocCallbackBase) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ProcessConstJsdocCallback_GOOG_DEFINE pkg_rhino.QualifiedName

func StaticGet_ProcessConstJsdocCallback_GOOG_DEFINE() pkg_rhino.QualifiedName {
	clinit_ProcessConstJsdocCallback()
	return static_ProcessConstJsdocCallback_GOOG_DEFINE
}
func StaticSet_ProcessConstJsdocCallback_GOOG_DEFINE(v pkg_rhino.QualifiedName) {
	clinit_ProcessConstJsdocCallback()
	static_ProcessConstJsdocCallback_GOOG_DEFINE = v
}

var static_ProcessConstJsdocCallback_GOOG_PROVIDE pkg_rhino.QualifiedName

func StaticGet_ProcessConstJsdocCallback_GOOG_PROVIDE() pkg_rhino.QualifiedName {
	clinit_ProcessConstJsdocCallback()
	return static_ProcessConstJsdocCallback_GOOG_PROVIDE
}
func StaticSet_ProcessConstJsdocCallback_GOOG_PROVIDE(v pkg_rhino.QualifiedName) {
	clinit_ProcessConstJsdocCallback()
	static_ProcessConstJsdocCallback_GOOG_PROVIDE = v
}

var static_ProcessConstJsdocCallback_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_ProcessConstJsdocCallback_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_ProcessConstJsdocCallback()
	return static_ProcessConstJsdocCallback_GOOG_REQUIRE
}
func StaticSet_ProcessConstJsdocCallback_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_ProcessConstJsdocCallback()
	static_ProcessConstJsdocCallback_GOOG_REQUIRE = v
}

var static_ProcessConstJsdocCallback_CJS_REQUIRE pkg_rhino.QualifiedName

func StaticGet_ProcessConstJsdocCallback_CJS_REQUIRE() pkg_rhino.QualifiedName {
	clinit_ProcessConstJsdocCallback()
	return static_ProcessConstJsdocCallback_CJS_REQUIRE
}
func StaticSet_ProcessConstJsdocCallback_CJS_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_ProcessConstJsdocCallback()
	static_ProcessConstJsdocCallback_CJS_REQUIRE = v
}

var clinitDone_ProcessConstJsdocCallback uint32
var clinitState_ProcessConstJsdocCallback j2gort.ClinitState

func clinit_ProcessConstJsdocCallback() {
	if atomic.LoadUint32(&clinitDone_ProcessConstJsdocCallback) != 0 {
		return
	}
	clinitSlow_ProcessConstJsdocCallback()
}

func clinitSlow_ProcessConstJsdocCallback() {
	if !j2gort.EnterClinit(&clinitState_ProcessConstJsdocCallback, &clinitDone_ProcessConstJsdocCallback) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessConstJsdocCallback, &clinitDone_ProcessConstJsdocCallback)
	static_ProcessConstJsdocCallback_GOOG_DEFINE = pkg_rhino.QualifiedName_Of("goog.define")
	static_ProcessConstJsdocCallback_GOOG_PROVIDE = pkg_rhino.QualifiedName_Of("goog.provide")
	static_ProcessConstJsdocCallback_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_ProcessConstJsdocCallback_CJS_REQUIRE = pkg_rhino.QualifiedName_Of("require")
}

func (this_ *ProcessConstJsdocCallbackBase) initFields_ProcessConstJsdocCallback() {
	_ = this_
}

func InitProcessConstJsdocCallback(this_ *ProcessConstJsdocCallbackBase, currentFile *FileInfo) {
	_ = currentFile
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ProcessConstJsdocCallback()
	this_.currentFile = currentFile
}

// Java method: ProcessConstJsdocCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessConstJsdocCallbackBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		if NodeUtil_IsStatementParent(parent) {
			this_.currentFile.RecordNameDeclaration(n.GetFirstChild())
		}
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if NodeUtil_GetRValueOfLValue(n) != 0 {
			this_.priv_ProcessConstJsdocCallback_processDeclarationWithRhs_6491(t, n)
		}
		this_.currentFile.RecordMemberFieldDef(n)
	case int32(73) /* Token.FUNCTION */ :
		if NodeUtil_IsStatementParent(parent) {
			this_.currentFile.RecordNameDeclaration(n.GetFirstChild())
		} else if ClassUtil_IsClassMethod(n) && ClassUtil_HasNamedClass(n) {
			this_.currentFile.RecordMethod(n)
		}
	case int32(93) /* Token.EXPR_RESULT */ :
		var expr pkg_rhino.Node = n.GetFirstChild()
		_ = expr
		{
			var callee pkg_rhino.Node = pkg_rhino.Node(0)
			_ = callee
			switch expr.TokenSwitchOrd_() {
			case int32(28) /* Token.CALL */ :
				callee = expr.GetFirstChild()
				if StaticGet_ProcessConstJsdocCallback_GOOG_PROVIDE().Matches(callee) {
					this_.currentFile.MarkProvided(expr.GetLastChild().GetString())
				} else if StaticGet_ProcessConstJsdocCallback_GOOG_REQUIRE().Matches(callee) || StaticGet_ProcessConstJsdocCallback_CJS_REQUIRE().Matches(callee) {
					this_.currentFile.RecordImport(expr.GetLastChild().GetString())
				} else if StaticGet_ProcessConstJsdocCallback_GOOG_DEFINE().Matches(callee) {
					this_.currentFile.RecordDefine(expr)
				}
			case int32(51) /* Token.ASSIGN */ :
				this_.priv_ProcessConstJsdocCallback_recordDeclaration_6489(t, expr.GetFirstChild(), expr.GetLastChild())
			case int32(26) /* Token.GETPROP */ :
				this_.currentFile.RecordNameDeclaration(expr)
			case int32(27) /* Token.GETELEM */ :
			default:
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected declaration: "+j2gort.StrVal(expr))))
			}
		}
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
		if !(n.HasOneChild()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		this_.priv_ProcessConstJsdocCallback_recordDeclaration_6489(t, n.GetFirstChild(), n.GetFirstChild().GetLastChild())
	case int32(100) /* Token.STRING_KEY */ :
		if parent.IsObjectLit() && n.HasOneChild() {
			this_.priv_ProcessConstJsdocCallback_processDeclarationWithRhs_6491(t, n)
			this_.currentFile.RecordStringKeyDeclaration(n)
		}
	default:
	}
}

// Java method: ProcessConstJsdocCallback.recordDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhs)
func (this_ *ProcessConstJsdocCallbackBase) priv_ProcessConstJsdocCallback_recordDeclaration_6489(t *NodeTraversal, lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = t
	_ = lhs
	_ = rhs
	if (((rhs != 0) && rhs.IsCall()) && StaticGet_ProcessConstJsdocCallback_GOOG_DEFINE().Matches(rhs.GetFirstChild())) && lhs.IsQualifiedName() {
		this_.currentFile.RecordDefine(rhs)
	} else if lhs.IsGetElem() {
		return
	} else {
		this_.priv_ProcessConstJsdocCallback_recordNameDeclaration_6490(lhs, rhs)
		if !(lhs.IsDestructuringLhs()) && (rhs != 0) {
			this_.priv_ProcessConstJsdocCallback_processDeclarationWithRhs_6491(t, lhs)
		}
	}
}

// Java method: ProcessConstJsdocCallback.recordNameDeclaration(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhs)
func (this_ *ProcessConstJsdocCallbackBase) priv_ProcessConstJsdocCallback_recordNameDeclaration_6490(lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = lhs
	_ = rhs
	if !(NodeUtil_IsNameDeclaration(lhs.GetParent()) || lhs.GetParent().IsAssign()) {
		j2gort.CheckArgument(false)
	}
	var isImport bool = PotentialDeclaration_IsImportRhs(rhs)
	_ = isImport
	var isAlias bool = PotentialDeclaration_IsAliasDeclaration(lhs, rhs)
	_ = isAlias
	NodeUtil_VisitLhsNodesInNode(lhs.GetParent(), j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
		name := j2gort.CastNode[pkg_rhino.Node](lp_1)
		_ = name
		if isAlias || isImport {
			this_.currentFile.RecordAliasDeclaration(name)
		} else {
			this_.currentFile.RecordNameDeclaration(name)
		}
	}))
}

// Java method: ProcessConstJsdocCallback.processDeclarationWithRhs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node lhs)
func (this_ *ProcessConstJsdocCallbackBase) priv_ProcessConstJsdocCallback_processDeclarationWithRhs_6491(t *NodeTraversal, lhs pkg_rhino.Node) {
	_ = t
	_ = lhs
	if !((lhs.IsQualifiedName() || lhs.IsStringKey()) || lhs.IsMemberFieldDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(lhs))
	}
	if !(NodeUtil_GetRValueOfLValue(lhs) != 0) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
	}
	if PotentialDeclaration_IsConstToBeInferred(lhs) {
		this_.self_.ProcessConstWithRhs(t, lhs)
	}
}

// Java method: ProcessConstJsdocCallback.processConstWithRhs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node lhs)
func (this_ *ProcessConstJsdocCallbackBase) ProcessConstWithRhs(t *NodeTraversal, lhs pkg_rhino.Node) {
	_ = t
	_ = lhs
	panic("abstract or native method invoked")
}
