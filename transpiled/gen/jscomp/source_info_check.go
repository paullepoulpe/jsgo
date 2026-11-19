/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SourceInfoCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.SourceInfoCheck
type SourceInfoCheck struct {
	compiler            AbstractCompiler `gson:"compiler"`
	idHash_             atomic.Int32
	requiresLineNumbers bool `gson:"requiresLineNumbers"`
}

func (this_ *SourceInfoCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceInfoCheck) IsJavaNil_() bool { return this_ == nil }
func (*SourceInfoCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SourceInfoCheck", "SourceInfoCheck"
}
func (this_ *SourceInfoCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceInfoCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SourceInfoCheck", "com.google.javascript.jscomp.SourceInfoCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceInfoCheck](v_)
	})
}
func (*SourceInfoCheck) Implements_jscomp_NodeTraversal_Callback() {}
func (*SourceInfoCheck) Implements_jscomp_CompilerPass()           {}

var static_SourceInfoCheck_MISSING_LINE_INFO *DiagnosticType

func StaticGet_SourceInfoCheck_MISSING_LINE_INFO() *DiagnosticType {
	clinit_SourceInfoCheck()
	return static_SourceInfoCheck_MISSING_LINE_INFO
}
func StaticSet_SourceInfoCheck_MISSING_LINE_INFO(v *DiagnosticType) {
	clinit_SourceInfoCheck()
	static_SourceInfoCheck_MISSING_LINE_INFO = v
}

var static_SourceInfoCheck_MISSING_LENGTH *DiagnosticType

func StaticGet_SourceInfoCheck_MISSING_LENGTH() *DiagnosticType {
	clinit_SourceInfoCheck()
	return static_SourceInfoCheck_MISSING_LENGTH
}
func StaticSet_SourceInfoCheck_MISSING_LENGTH(v *DiagnosticType) {
	clinit_SourceInfoCheck()
	static_SourceInfoCheck_MISSING_LENGTH = v
}

var static_SourceInfoCheck_MISSING_SOURCE_NAME *DiagnosticType

func StaticGet_SourceInfoCheck_MISSING_SOURCE_NAME() *DiagnosticType {
	clinit_SourceInfoCheck()
	return static_SourceInfoCheck_MISSING_SOURCE_NAME
}
func StaticSet_SourceInfoCheck_MISSING_SOURCE_NAME(v *DiagnosticType) {
	clinit_SourceInfoCheck()
	static_SourceInfoCheck_MISSING_SOURCE_NAME = v
}

var clinitDone_SourceInfoCheck uint32
var clinitState_SourceInfoCheck j2gort.ClinitState

func clinit_SourceInfoCheck() {
	if atomic.LoadUint32(&clinitDone_SourceInfoCheck) != 0 {
		return
	}
	clinitSlow_SourceInfoCheck()
}

func clinitSlow_SourceInfoCheck() {
	if !j2gort.EnterClinit(&clinitState_SourceInfoCheck, &clinitDone_SourceInfoCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceInfoCheck, &clinitDone_SourceInfoCheck)
	static_SourceInfoCheck_MISSING_LINE_INFO = DiagnosticType_Error("JSC_MISSING_LINE_INFO", "No source location information associated with {0}.\nMost likely a Node has been created without setting the source file and line/column location.  Usually this is done using Node.srcrefIfMissing and supplying a Node from the source AST.")
	static_SourceInfoCheck_MISSING_LENGTH = DiagnosticType_Error("JSC_MISSING_LENGTH", "Negative length associated with {0}.\nMost likely a Node's source information was set incorrectly at parse time.")
	static_SourceInfoCheck_MISSING_SOURCE_NAME = DiagnosticType_Error("JSC_MISSING_SOURCE_NAME", "No source name associated with {0}.\nMost likely a new type was created without setting the source name.")
}

func (this_ *SourceInfoCheck) initFields_SourceInfoCheck() {
	_ = this_

}

// Java constructor: SourceInfoCheck(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewSourceInfoCheck(compiler AbstractCompiler) *SourceInfoCheck {
	clinit_SourceInfoCheck()
	this_ := &SourceInfoCheck{}
	InitSourceInfoCheck(this_, compiler)
	return this_
}

func InitSourceInfoCheck(this_ *SourceInfoCheck, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_SourceInfoCheck()
	this_.compiler = compiler
}

// Java method: SourceInfoCheck.setCheckSubTree(com.google.javascript.rhino.Node root)
func (this_ *SourceInfoCheck) SetCheckSubTree(root pkg_rhino.Node) {
	_ = root
	this_.requiresLineNumbers = true
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: SourceInfoCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SourceInfoCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.requiresLineNumbers = false
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: SourceInfoCheck.shouldTraverse(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SourceInfoCheck) ShouldTraverse_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = unused
	_ = n
	_ = parent
	if n.IsScript() {
		this_.requiresLineNumbers = true
	}
	return true
}

// Java method: SourceInfoCheck.visit(com.google.javascript.jscomp.NodeTraversal unused, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SourceInfoCheck) Visit_3(unused *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = unused
	_ = n
	_ = parent
	if n.IsScript() {
		this_.requiresLineNumbers = false
	} else if this_.requiresLineNumbers {
		if n.GetLineno() == int32(-1) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_SourceInfoCheck_MISSING_LINE_INFO(), []string{n.ToStringTree()}))
		} else if n.GetLength() < int32(0) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_SourceInfoCheck_MISSING_LENGTH(), []string{n.ToStringTree()}))
		}
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		if info != nil {
			this_.priv_SourceInfoCheck_checkJSDoc_5142(info)
		}
	}
}

// Java method: SourceInfoCheck.checkJSDoc(com.google.javascript.rhino.JSDocInfo info)
func (this_ *SourceInfoCheck) priv_SourceInfoCheck_checkJSDoc_5142(info *pkg_rhino.JSDocInfo) {
	_ = info
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetTypeExpressions()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			expression := j2gort.Cast[pkg_rhino.JSTypeExpression](rawElem_2)
			_ = expression
			if expression.GetSourceName() == j2gort.NullString {
				this_.compiler.Report_1(JSError_Make_3(expression.GetRoot(), StaticGet_SourceInfoCheck_MISSING_SOURCE_NAME(), []string{("JSTypeExpression " + j2gort.StrVal(expression))}))
			} else {
				NodeUtil_VisitPreOrder_2(expression.GetRoot(), func() NodeUtil_Visitor {
					bound_3 := this_
					return NewLambda_NodeUtil_Visitor(func(mrp0_ pkg_rhino.Node) { bound_3.priv_SourceInfoCheck_checkSourceFile_5143(mrp0_) })
				}())
			}
		}
	}
}

// Java method: SourceInfoCheck.checkSourceFile(com.google.javascript.rhino.Node n)
func (this_ *SourceInfoCheck) priv_SourceInfoCheck_checkSourceFile_5143(n pkg_rhino.Node) {
	_ = n
	var sourceName pkg_rhino.StaticSourceFile = n.GetStaticSourceFile()
	_ = sourceName
	if sourceName == nil || sourceName.IsJavaNil_() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_SourceInfoCheck_MISSING_SOURCE_NAME(), []string{("type node " + j2gort.StrValString(n.ToStringTree()))}))
	}
}
func (*SourceInfoCheck) CanReleaseSyntacticScopes_() {}
