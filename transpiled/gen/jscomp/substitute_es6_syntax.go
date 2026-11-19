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
// Port of com/google/javascript/jscomp/SubstituteEs6Syntax.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.SubstituteEs6Syntax
type SubstituteEs6Syntax struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                       AbstractCompiler `gson:"compiler"`
	objectLiteralShorthandWasAdded bool             `gson:"objectLiteralShorthandWasAdded"`
}

func (this_ *SubstituteEs6Syntax) BindSelf_(s any) {
	_ = s
}
func (this_ *SubstituteEs6Syntax) IsJavaNil_() bool { return this_ == nil }
func (*SubstituteEs6Syntax) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SubstituteEs6Syntax", "SubstituteEs6Syntax"
}
func (this_ *SubstituteEs6Syntax) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SubstituteEs6Syntax) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SubstituteEs6Syntax", "com.google.javascript.jscomp.SubstituteEs6Syntax", func(v_ any) bool {
		return j2gort.InstanceOf[*SubstituteEs6Syntax](v_)
	})
}
func (*SubstituteEs6Syntax) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SubstituteEs6Syntax) Implements_jscomp_CompilerPass()                            {}
func (*SubstituteEs6Syntax) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SubstituteEs6Syntax) initFields_SubstituteEs6Syntax() {
	_ = this_

}

// Java constructor: SubstituteEs6Syntax(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewSubstituteEs6Syntax(compiler AbstractCompiler) *SubstituteEs6Syntax {
	this_ := &SubstituteEs6Syntax{}
	InitSubstituteEs6Syntax(this_, compiler)
	return this_
}

func InitSubstituteEs6Syntax(this_ *SubstituteEs6Syntax, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SubstituteEs6Syntax()
	this_.compiler = compiler
}

// Java method: SubstituteEs6Syntax.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SubstituteEs6Syntax) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: SubstituteEs6Syntax.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SubstituteEs6Syntax) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		if n.IsArrowFunction() {
			this_.priv_SubstituteEs6Syntax_maybeSimplifyArrowFunctionBody_5194(n, n.GetLastChild())
		}
	case int32(100) /* Token.STRING_KEY */ :
		if n.GetFirstChild().IsName() && func() bool {
			sr_1 := n.GetFirstChild().GetString()
			se_2 := n.GetString()
			return j2gort.CheckNotNullString(sr_1) == se_2
		}() {
			n.SetShorthandProperty(true)
			this_.objectLiteralShorthandWasAdded = true
		}
	case int32(94) /* Token.SCRIPT */ :
		if this_.objectLiteralShorthandWasAdded {
			NodeUtil_AddFeatureToScript(n, pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES(), this_.compiler)
			this_.objectLiteralShorthandWasAdded = false
		}
	default:
	}
}

// Java method: SubstituteEs6Syntax.maybeSimplifyArrowFunctionBody(com.google.javascript.rhino.Node arrowFunction, com.google.javascript.rhino.Node body)
func (this_ *SubstituteEs6Syntax) priv_SubstituteEs6Syntax_maybeSimplifyArrowFunctionBody_5194(arrowFunction pkg_rhino.Node, body pkg_rhino.Node) {
	_ = arrowFunction
	_ = body
	if !(arrowFunction.IsArrowFunction()) {
		j2gort.CheckArgument(false)
	}
	if (!(body.IsBlock()) || !(body.HasOneChild())) || !(body.GetFirstChild().IsReturn()) {
		return
	}
	var returnValue pkg_rhino.Node = body.GetFirstChild().RemoveFirstChild()
	_ = returnValue
	var replacement pkg_rhino.Node
	_ = replacement
	if returnValue != 0 {
		replacement = returnValue
	} else {
		replacement = pkg_rhino.IR_Name("undefined")
	}
	body.ReplaceWith(replacement)
	this_.compiler.ReportChangeToEnclosingScope(replacement)
}
func (*SubstituteEs6Syntax) AlwaysShouldTraverse_()      {}
func (*SubstituteEs6Syntax) CanReleaseSyntacticScopes_() {}
