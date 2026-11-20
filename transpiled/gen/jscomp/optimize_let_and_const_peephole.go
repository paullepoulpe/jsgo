/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/OptimizeLetAndConstPeephole.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.OptimizeLetAndConstPeephole
type OptimizeLetAndConstPeephole struct {
	AbstractPeepholeOptimizationBase
	assumeOutputIsWrapped bool `gson:"assumeOutputIsWrapped"`
}

func (this_ *OptimizeLetAndConstPeephole) BindSelf_(s any) {
	_ = s
	this_.AbstractPeepholeOptimizationBase.BindSelf_(s)
}
func (this_ *OptimizeLetAndConstPeephole) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeLetAndConstPeephole) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.OptimizeLetAndConstPeephole", "OptimizeLetAndConstPeephole" }
func (this_ *OptimizeLetAndConstPeephole) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeOptimizationBase.IdentityHash_()
}
func (*OptimizeLetAndConstPeephole) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeLetAndConstPeephole", "com.google.javascript.jscomp.OptimizeLetAndConstPeephole", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeLetAndConstPeephole](v_)
	})
}
func (*OptimizeLetAndConstPeephole) Implements_jscomp_AbstractPeepholeOptimization() {}

func (this_ *OptimizeLetAndConstPeephole) initFields_OptimizeLetAndConstPeephole() {
	_ = this_
}

// Java constructor: OptimizeLetAndConstPeephole(boolean assumeOutputIsWrapped)
func NewOptimizeLetAndConstPeephole(assumeOutputIsWrapped bool) *OptimizeLetAndConstPeephole {
	this_ := &OptimizeLetAndConstPeephole{}
	this_.BindSelf_(this_)
	InitOptimizeLetAndConstPeephole(this_, assumeOutputIsWrapped)
	return this_
}

func InitOptimizeLetAndConstPeephole(this_ *OptimizeLetAndConstPeephole, assumeOutputIsWrapped bool) {
	_ = assumeOutputIsWrapped
	InitAbstractPeepholeOptimization(&this_.AbstractPeepholeOptimizationBase)
	this_.initFields_OptimizeLetAndConstPeephole()
	this_.assumeOutputIsWrapped = assumeOutputIsWrapped
}

// Java method: OptimizeLetAndConstPeephole.optimizeSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *OptimizeLetAndConstPeephole) OptimizeSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	switch subtree.TokenSwitchOrd_() {
	case int32(111) /* Token.LET */:
		return this_.priv_OptimizeLetAndConstPeephole_processLet_4094(subtree)
	case int32(97) /* Token.CONST */:
		return this_.priv_OptimizeLetAndConstPeephole_processConst_4095(subtree)
	default:
		return subtree
	}
	panic("unreachable switch expression")
}

// Java method: OptimizeLetAndConstPeephole.processLet(com.google.javascript.rhino.Node n)
func (this_ *OptimizeLetAndConstPeephole) priv_OptimizeLetAndConstPeephole_processLet_4094(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if ((!(this_.assumeOutputIsWrapped) && n.GetParent().IsScript())) {
		return n
	}
	if (this_.priv_OptimizeLetAndConstPeephole_isDirectlyInHoistScope_4096(n)) {
		n.SetToken(pkg_rhino.StaticGet_Token_VAR())
		this_.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java method: OptimizeLetAndConstPeephole.processConst(com.google.javascript.rhino.Node n)
func (this_ *OptimizeLetAndConstPeephole) priv_OptimizeLetAndConstPeephole_processConst_4095(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if ((!(this_.assumeOutputIsWrapped) && n.GetParent().IsScript())) {
		return n
	}
	if (this_.priv_OptimizeLetAndConstPeephole_isDirectlyInHoistScope_4096(n)) {
		n.SetToken(pkg_rhino.StaticGet_Token_VAR())
	} else {
		n.SetToken(pkg_rhino.StaticGet_Token_LET())
		this_.AddFeatureToEnclosingScript(pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS())
	}
	this_.ReportChangeToEnclosingScope(n)
	return n
}

// Java method: OptimizeLetAndConstPeephole.isDirectlyInHoistScope(com.google.javascript.rhino.Node n)
func (this_ *OptimizeLetAndConstPeephole) priv_OptimizeLetAndConstPeephole_isDirectlyInHoistScope_4096(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	for (((parent != 0) && parent.IsLabel())) {
		parent = parent.GetParent()
	}
	if ((parent == 0)) {
		return false
	}
	if ((parent.IsScript() || parent.IsModuleBody())) {
		return true
	}
	if (parent.IsBlock()) {
		var grandparent pkg_rhino.Node = parent.GetParent()
		_ = grandparent
		if (((grandparent != 0) && ((grandparent.IsFunction() || grandparent.IsClassMembers())))) {
			return true
		}
	}
	return false
}

