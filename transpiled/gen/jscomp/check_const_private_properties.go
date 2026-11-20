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
// Port of com/google/javascript/jscomp/lint/CheckConstPrivateProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckConstPrivateProperties
type CheckConstPrivateProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                  AbstractCompiler          `gson:"compiler"`
	candidates                []pkg_rhino.Node          `gson:"candidates"`
	modified                  j2gort.OrderedSet[string] `gson:"modified"`
	constructorsAndInterfaces j2gort.OrderedSet[string] `gson:"constructorsAndInterfaces"`
}

func (this_ *CheckConstPrivateProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckConstPrivateProperties) IsJavaNil_() bool { return this_ == nil }
func (*CheckConstPrivateProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckConstPrivateProperties", "CheckConstPrivateProperties"
}
func (this_ *CheckConstPrivateProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckConstPrivateProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckConstPrivateProperties", "com.google.javascript.jscomp.lint.CheckConstPrivateProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckConstPrivateProperties](v_)
	})
}
func (*CheckConstPrivateProperties) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckConstPrivateProperties) Implements_jscomp_CompilerPass()                            {}
func (*CheckConstPrivateProperties) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckConstPrivateProperties_MISSING_CONST_PROPERTY *DiagnosticType

func StaticGet_CheckConstPrivateProperties_MISSING_CONST_PROPERTY() *DiagnosticType {
	clinit_CheckConstPrivateProperties()
	return static_CheckConstPrivateProperties_MISSING_CONST_PROPERTY
}
func StaticSet_CheckConstPrivateProperties_MISSING_CONST_PROPERTY(v *DiagnosticType) {
	clinit_CheckConstPrivateProperties()
	static_CheckConstPrivateProperties_MISSING_CONST_PROPERTY = v
}

var clinitDone_CheckConstPrivateProperties uint32
var clinitState_CheckConstPrivateProperties j2gort.ClinitState

func clinit_CheckConstPrivateProperties() {
	if atomic.LoadUint32(&clinitDone_CheckConstPrivateProperties) != 0 {
		return
	}
	clinitSlow_CheckConstPrivateProperties()
}

func clinitSlow_CheckConstPrivateProperties() {
	if !j2gort.EnterClinit(&clinitState_CheckConstPrivateProperties, &clinitDone_CheckConstPrivateProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckConstPrivateProperties, &clinitDone_CheckConstPrivateProperties)
	static_CheckConstPrivateProperties_MISSING_CONST_PROPERTY = DiagnosticType_Disabled("JSC_MISSING_CONST_PROPERTY", "Private property {0} is never modified, use the @const annotation")
}

func (this_ *CheckConstPrivateProperties) initFields_CheckConstPrivateProperties() {
	_ = this_
	this_.candidates = []pkg_rhino.Node{}
	this_.modified = j2gort.NewOrderedSet[string]()
	this_.constructorsAndInterfaces = j2gort.NewOrderedSet[string]()
}

// Java constructor: CheckConstPrivateProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckConstPrivateProperties(compiler AbstractCompiler) *CheckConstPrivateProperties {
	clinit_CheckConstPrivateProperties()
	this_ := &CheckConstPrivateProperties{}
	InitCheckConstPrivateProperties(this_, compiler)
	return this_
}

func InitCheckConstPrivateProperties(this_ *CheckConstPrivateProperties, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckConstPrivateProperties()
	this_.compiler = compiler
}

// Java method: CheckConstPrivateProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckConstPrivateProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckConstPrivateProperties.reportMissingConst(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckConstPrivateProperties) priv_CheckConstPrivateProperties_reportMissingConst_6560(t *NodeTraversal) {
	_ = t
	for _, n := range this_.candidates {
		_ = n
		var propName string = n.GetString()
		_ = propName
		if !(this_.modified.Contains(propName)) {
			t.Report_3(n, StaticGet_CheckConstPrivateProperties_MISSING_CONST_PROPERTY(), []string{propName})
		}
	}

	this_.candidates = this_.candidates[:0]
	this_.modified.Clear()
}

// Java method: CheckConstPrivateProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckConstPrivateProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_CheckConstPrivateProperties_reportMissingConst_6560(t)
	case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
		var propNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = propNode
		if n.IsGetProp() {
			propNode = n
		} else if n.GetLastChild().IsStringLit() {
			propNode = n.GetLastChild()
		} else {
			return
		}
		if this_.priv_CheckConstPrivateProperties_isCandidatePropertyDefinition_6562(n) {
			this_.candidates = append(this_.candidates, propNode)
		} else if this_.priv_CheckConstPrivateProperties_isModificationOp_6564(n) {
			this_.modified.Add(propNode.GetString())
		}
	case int32(73) /* Token.FUNCTION */ :
		var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = info
		if (info != nil) && (info.IsConstructor() || info.IsInterface()) {
			this_.priv_CheckConstPrivateProperties_recordConstructorOrInterface_6561(n)
		}
	case int32(105) /* Token.CLASS */ :
		this_.priv_CheckConstPrivateProperties_recordConstructorOrInterface_6561(n)
	default:
	}
}

// Java method: CheckConstPrivateProperties.recordConstructorOrInterface(com.google.javascript.rhino.Node n)
func (this_ *CheckConstPrivateProperties) priv_CheckConstPrivateProperties_recordConstructorOrInterface_6561(n pkg_rhino.Node) {
	_ = n
	var className string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	_ = className
	if className != j2gort.NullString {
		this_.constructorsAndInterfaces.Add(className)
	}
}

// Java method: CheckConstPrivateProperties.isCandidatePropertyDefinition(com.google.javascript.rhino.Node n)
func (this_ *CheckConstPrivateProperties) priv_CheckConstPrivateProperties_isCandidatePropertyDefinition_6562(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsLhsOfAssign(n)) {
		return false
	}
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	if !(target.IsThis() || this_.constructorsAndInterfaces.Contains(target.GetQualifiedName())) {
		return false
	}
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = info
	return (((((((info != nil) && (info.GetVisibility() == pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE())) && !(info.IsConstant())) && !(info.HasTypedefType())) && !(info.HasEnumParameterType())) && !(info.IsInterface())) && !(this_.priv_CheckConstPrivateProperties_isFunctionProperty_6563(n)))
}

// Java method: CheckConstPrivateProperties.isFunctionProperty(com.google.javascript.rhino.Node n)
func (this_ *CheckConstPrivateProperties) priv_CheckConstPrivateProperties_isFunctionProperty_6563(n pkg_rhino.Node) bool {
	_ = n
	if n.IsGetElem() {
		return false
	}
	var assignedValue pkg_rhino.Node = NodeUtil_GetAssignedValue(n)
	_ = assignedValue
	return ((assignedValue != 0) && assignedValue.IsFunction())
}

// Java method: CheckConstPrivateProperties.isModificationOp(com.google.javascript.rhino.Node n)
func (this_ *CheckConstPrivateProperties) priv_CheckConstPrivateProperties_isModificationOp_6564(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsLValue(n) || n.GetParent().IsDelProp())
}
func (*CheckConstPrivateProperties) AlwaysShouldTraverse_()      {}
func (*CheckConstPrivateProperties) CanReleaseSyntacticScopes_() {}
