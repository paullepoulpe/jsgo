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
// Port of com/google/javascript/jscomp/PolymerPassFindExterns.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.PolymerPassFindExterns
type PolymerPassFindExterns struct {
	polymerElementProps   *j2gort.Ext_com_google_common_collect_ImmutableList_Builder `gson:"polymerElementProps"`
	polymerElementExterns pkg_rhino.Node                                              `gson:"polymerElementExterns"`
	idHash_               atomic.Int32
}

func (this_ *PolymerPassFindExterns) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPassFindExterns) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPassFindExterns) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerPassFindExterns", "PolymerPassFindExterns"
}
func (this_ *PolymerPassFindExterns) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerPassFindExterns) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerPassFindExterns", "com.google.javascript.jscomp.PolymerPassFindExterns", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerPassFindExterns](v_)
	})
}
func (*PolymerPassFindExterns) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_PolymerPassFindExterns_POLYMER_ELEMENT_NAME() string { return "PolymerElement" }
func (this_ *PolymerPassFindExterns) initFields_PolymerPassFindExterns() {
	_ = this_
}

// Java constructor: PolymerPassFindExterns()
func NewPolymerPassFindExterns() *PolymerPassFindExterns {
	this_ := &PolymerPassFindExterns{}
	InitPolymerPassFindExterns(this_)
	return this_
}

func InitPolymerPassFindExterns(this_ *PolymerPassFindExterns) {
	this_.initFields_PolymerPassFindExterns()
	this_.polymerElementProps = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
}

// Java method: PolymerPassFindExterns.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *PolymerPassFindExterns) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return !(n.IsModuleBody())
}

// Java method: PolymerPassFindExterns.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *PolymerPassFindExterns) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if PolymerPassFindExterns_IsPolymerElementExterns(n) {
		this_.polymerElementExterns = n
	} else if PolymerPassFindExterns_IsPolymerElementPropExpr(n) {
		this_.polymerElementProps.AddVoid_(j2gort.BoxNode(n))
	}
}

// Java method: PolymerPassFindExterns.getPolymerElementProps()
func (this_ *PolymerPassFindExterns) GetPolymerElementProps() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.polymerElementProps.BuildImmList_()
}

// Java method: PolymerPassFindExterns.getPolymerElementExterns()
func (this_ *PolymerPassFindExterns) GetPolymerElementExterns() pkg_rhino.Node {
	return this_.polymerElementExterns
}

// Java static method: PolymerPassFindExterns.isPolymerElementExterns(com.google.javascript.rhino.Node value)
func PolymerPassFindExterns_IsPolymerElementExterns(value pkg_rhino.Node) bool {
	_ = value
	return (NodeUtil_IsNameDeclaration(value) && value.GetFirstChild().MatchesQualifiedName_1_string("PolymerElement"))
}

// Java static method: PolymerPassFindExterns.isPolymerElementPropExpr(com.google.javascript.rhino.Node value)
func PolymerPassFindExterns_IsPolymerElementPropExpr(value pkg_rhino.Node) bool {
	_ = value
	return ((((((value != 0) && value.IsExprResult()) && (value.GetFirstFirstChild() != 0)) && value.GetFirstFirstChild().IsGetProp()) && value.GetFirstFirstChild().IsQualifiedName()) && NodeUtil_GetRootOfQualifiedName_1_Node(value.GetFirstFirstChild()).MatchesQualifiedName_1_string("PolymerElement"))
}
func (*PolymerPassFindExterns) CanReleaseSyntacticScopes_() {}
