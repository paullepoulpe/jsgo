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
// Port of com/google/javascript/jscomp/AstManipulations.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.AstManipulations
type AstManipulations struct {
	idHash_ atomic.Int32
}

func (this_ *AstManipulations) BindSelf_(s any) {
	_ = s
}
func (this_ *AstManipulations) IsJavaNil_() bool { return this_ == nil }
func (*AstManipulations) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AstManipulations", "AstManipulations"
}
func (this_ *AstManipulations) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AstManipulations) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AstManipulations", "com.google.javascript.jscomp.AstManipulations", func(v_ any) bool {
		return j2gort.InstanceOf[*AstManipulations](v_)
	})
}

func (this_ *AstManipulations) initFields_AstManipulations() {
	_ = this_
}

// Java constructor: AstManipulations()
func NewAstManipulations() *AstManipulations {
	this_ := &AstManipulations{}
	InitAstManipulations(this_)
	return this_
}

func InitAstManipulations(this_ *AstManipulations) {
	this_.initFields_AstManipulations()
}

// Java static method: AstManipulations.fuseExpressions(com.google.javascript.rhino.Node exp1, com.google.javascript.rhino.Node exp2)
func AstManipulations_FuseExpressions(exp1 pkg_rhino.Node, exp2 pkg_rhino.Node) pkg_rhino.Node {
	_ = exp1
	_ = exp2
	if !(exp1.GetParent() == 0) {
		j2gort.CheckArgumentFmt(false, "Expected detached node, got %s", j2gort.BoxNode(exp1))
	}
	if !(exp2.GetParent() == 0) {
		j2gort.CheckArgumentFmt(false, "Expected detached node, got %s", j2gort.BoxNode(exp2))
	}
	if exp2.IsEmpty() {
		return exp1
	}
	var comma pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_COMMA(), exp1)
	_ = comma
	comma.SrcrefIfMissing(exp2)
	if exp2.IsComma() {
		var leftMostChild pkg_rhino.Node = exp2
		_ = leftMostChild
		for leftMostChild.IsComma() {
			leftMostChild = leftMostChild.GetFirstChild()
		}
		var parent pkg_rhino.Node = leftMostChild.GetParent()
		_ = parent
		comma.AddChildToBack(leftMostChild.Detach())
		parent.AddChildToFront(comma)
		return exp2
	} else {
		comma.AddChildToBack(exp2)
		return comma
	}
}
