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
// Port of com/google/javascript/jscomp/lint/CheckUnusedPrivateProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckUnusedPrivateProperties
type CheckUnusedPrivateProperties struct {
	compiler                  AbstractCompiler          `gson:"compiler"`
	used                      j2gort.OrderedSet[string] `gson:"used"`
	candidates                []pkg_rhino.Node          `gson:"candidates"`
	constructorsAndInterfaces j2gort.OrderedSet[string] `gson:"constructorsAndInterfaces"`
	idHash_                   atomic.Int32
}

func (this_ *CheckUnusedPrivateProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckUnusedPrivateProperties) IsJavaNil_() bool { return this_ == nil }
func (*CheckUnusedPrivateProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckUnusedPrivateProperties", "CheckUnusedPrivateProperties"
}
func (this_ *CheckUnusedPrivateProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckUnusedPrivateProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckUnusedPrivateProperties", "com.google.javascript.jscomp.lint.CheckUnusedPrivateProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckUnusedPrivateProperties](v_)
	})
}
func (*CheckUnusedPrivateProperties) Implements_jscomp_CompilerPass()           {}
func (*CheckUnusedPrivateProperties) Implements_jscomp_NodeTraversal_Callback() {}

var static_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY *DiagnosticType

func StaticGet_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY() *DiagnosticType {
	clinit_CheckUnusedPrivateProperties()
	return static_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY
}
func StaticSet_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY(v *DiagnosticType) {
	clinit_CheckUnusedPrivateProperties()
	static_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY = v
}

var clinitDone_CheckUnusedPrivateProperties uint32
var clinitState_CheckUnusedPrivateProperties j2gort.ClinitState

func clinit_CheckUnusedPrivateProperties() {
	if atomic.LoadUint32(&clinitDone_CheckUnusedPrivateProperties) != 0 {
		return
	}
	clinitSlow_CheckUnusedPrivateProperties()
}

func clinitSlow_CheckUnusedPrivateProperties() {
	if !j2gort.EnterClinit(&clinitState_CheckUnusedPrivateProperties, &clinitDone_CheckUnusedPrivateProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckUnusedPrivateProperties, &clinitDone_CheckUnusedPrivateProperties)
	static_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY = DiagnosticType_Disabled("JSC_UNUSED_PRIVATE_PROPERTY", "Private property {0} is never read")
}

func (this_ *CheckUnusedPrivateProperties) initFields_CheckUnusedPrivateProperties() {
	_ = this_
	this_.used = j2gort.NewOrderedSet[string]()
	this_.candidates = []pkg_rhino.Node{}
	this_.constructorsAndInterfaces = j2gort.NewOrderedSet[string]()
}

// Java constructor: CheckUnusedPrivateProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckUnusedPrivateProperties(compiler AbstractCompiler) *CheckUnusedPrivateProperties {
	clinit_CheckUnusedPrivateProperties()
	this_ := &CheckUnusedPrivateProperties{}
	InitCheckUnusedPrivateProperties(this_, compiler)
	return this_
}

func InitCheckUnusedPrivateProperties(this_ *CheckUnusedPrivateProperties, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckUnusedPrivateProperties()
	this_.compiler = compiler
}

// Java method: CheckUnusedPrivateProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckUnusedPrivateProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckUnusedPrivateProperties.reportUnused(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckUnusedPrivateProperties) priv_CheckUnusedPrivateProperties_reportUnused_6630(t *NodeTraversal) {
	_ = t
	for _, n := range this_.candidates {
		_ = n
		var propName string = this_.priv_CheckUnusedPrivateProperties_getPropName_6631(n)
		_ = propName
		if !(this_.used.Contains(propName)) {
			t.Report_3(n, StaticGet_CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY(), []string{propName})
		}
	}

}

// Java method: CheckUnusedPrivateProperties.getPropName(com.google.javascript.rhino.Node n)
func (this_ *CheckUnusedPrivateProperties) priv_CheckUnusedPrivateProperties_getPropName_6631(n pkg_rhino.Node) string {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		return n.GetString()
	default:
	}
	panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Unexpected node type: "+j2gort.StrVal(n))))
}

// Java method: CheckUnusedPrivateProperties.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUnusedPrivateProperties) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.used.Clear()
		this_.used.Add("constructor")
		this_.candidates = this_.candidates[:0]
	}
	return true
}

// Java method: CheckUnusedPrivateProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUnusedPrivateProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_CheckUnusedPrivateProperties_reportUnused_6630(t)
	case int32(26) /* Token.GETPROP */ :
		var propName string = n.GetString()
		_ = propName
		if CheckUnusedPrivateProperties_IsPinningPropertyUse(n) || !(this_.priv_CheckUnusedPrivateProperties_isCandidatePropertyDefinition_6635(n)) {
			this_.used.Add(propName)
		} else {
			if this_.priv_CheckUnusedPrivateProperties_isCheckablePrivatePropDecl_6634(n) {
				this_.candidates = append(this_.candidates, n)
			}
		}
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		if this_.priv_CheckUnusedPrivateProperties_isCheckablePrivatePropDecl_6634(n) {
			this_.candidates = append(this_.candidates, n)
		}
	case int32(47) /* Token.OBJECTLIT */ :
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				if ((c.IsStringKey() || c.IsGetterDef()) || c.IsSetterDef()) || c.IsMemberFunctionDef() {
					this_.used.Add(c.GetString())
				}
			}
		}
	case int32(28) /* Token.CALL */ :
		var target pkg_rhino.Node = n.GetFirstChild()
		_ = target
		if n.HasMoreThanOneChild() && this_.compiler.GetCodingConvention().IsPropertyRenameFunction(target) {
			var propName_2 pkg_rhino.Node = target.GetNext()
			_ = propName_2
			if propName_2.IsStringLit() {
				this_.used.Add(propName_2.GetString())
			}
		}
	case int32(73) /* Token.FUNCTION */ :
		var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = info
		if (info != nil) && (info.IsConstructor() || info.IsInterface()) {
			this_.priv_CheckUnusedPrivateProperties_recordConstructorOrInterface_6632(n)
		}
	case int32(105) /* Token.CLASS */ :
		this_.priv_CheckUnusedPrivateProperties_recordConstructorOrInterface_6632(n)
	default:
	}
}

// Java method: CheckUnusedPrivateProperties.recordConstructorOrInterface(com.google.javascript.rhino.Node n)
func (this_ *CheckUnusedPrivateProperties) priv_CheckUnusedPrivateProperties_recordConstructorOrInterface_6632(n pkg_rhino.Node) {
	_ = n
	var className string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	_ = className
	if className != j2gort.NullString {
		this_.constructorsAndInterfaces.Add(className)
	}
}

// Java method: CheckUnusedPrivateProperties.isPrivatePropDecl(com.google.javascript.rhino.Node n)
func (this_ *CheckUnusedPrivateProperties) priv_CheckUnusedPrivateProperties_isPrivatePropDecl_6633(n pkg_rhino.Node) bool {
	_ = n
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	return ((info != nil) && (info.GetVisibility() == pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE()))
}

// Java method: CheckUnusedPrivateProperties.isCheckablePrivatePropDecl(com.google.javascript.rhino.Node n)
func (this_ *CheckUnusedPrivateProperties) priv_CheckUnusedPrivateProperties_isCheckablePrivatePropDecl_6634(n pkg_rhino.Node) bool {
	_ = n
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	return ((this_.priv_CheckUnusedPrivateProperties_isPrivatePropDecl_6633(n) && !(info.HasTypedefType())) && !(info.IsInterface()))
}

// Java method: CheckUnusedPrivateProperties.isCandidatePropertyDefinition(com.google.javascript.rhino.Node n)
func (this_ *CheckUnusedPrivateProperties) priv_CheckUnusedPrivateProperties_isCandidatePropertyDefinition_6635(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	return ((target.IsThis() || this_.constructorsAndInterfaces.Contains(target.GetQualifiedName())) || (target.IsGetProp() && (j2gort.CheckNotNullString(target.GetString()) == "prototype")))
}

// Java static method: CheckUnusedPrivateProperties.isPinningPropertyUse(com.google.javascript.rhino.Node n)
func CheckUnusedPrivateProperties_IsPinningPropertyUse(n pkg_rhino.Node) bool {
	clinit_CheckUnusedPrivateProperties()
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if n == parent.GetFirstChild() {
		if parent.IsExprResult() {
			return false
		} else if parent.IsAssign() {
			return false
		} else if (NodeUtil_IsAssignmentOp(parent) || parent.IsInc()) || parent.IsDec() {
			return NodeUtil_IsExpressionResultUsed(parent)
		}
	}
	return true
}
func (*CheckUnusedPrivateProperties) CanReleaseSyntacticScopes_() {}
