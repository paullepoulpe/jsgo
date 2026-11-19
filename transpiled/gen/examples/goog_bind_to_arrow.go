/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/refactoring/examples/GoogBindToArrow.java (generated mechanically by j2go).

package examples

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomp "github.com/paullepoulpe/jsgo/transpiled/gen/jscomp"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.refactoring.examples.GoogBindToArrow
type GoogBindToArrow struct {
	pkg_jscomp.ScannerBase
}

func (this_ *GoogBindToArrow) BindSelf_(s any) {
	_ = s
}
func (this_ *GoogBindToArrow) IsJavaNil_() bool { return this_ == nil }
func (*GoogBindToArrow) JavaClassInfo_() (string, string) { return "com.google.javascript.refactoring.examples.GoogBindToArrow", "GoogBindToArrow" }
func (this_ *GoogBindToArrow) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ScannerBase.IdentityHash_()
}
func (*GoogBindToArrow) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.examples.GoogBindToArrow", "com.google.javascript.refactoring.examples.GoogBindToArrow", func(v_ any) bool {
		return j2gort.InstanceOf[*GoogBindToArrow](v_)
	})
}
func (*GoogBindToArrow) Implements_jscomp_Scanner() {}

var static_GoogBindToArrow_GOOG_BIND_NAME pkg_rhino.QualifiedName
func StaticGet_GoogBindToArrow_GOOG_BIND_NAME() pkg_rhino.QualifiedName {
	clinit_GoogBindToArrow()
	return static_GoogBindToArrow_GOOG_BIND_NAME
}
func StaticSet_GoogBindToArrow_GOOG_BIND_NAME(v pkg_rhino.QualifiedName) {
	clinit_GoogBindToArrow()
	static_GoogBindToArrow_GOOG_BIND_NAME = v
}
var clinitDone_GoogBindToArrow uint32
var clinitState_GoogBindToArrow j2gort.ClinitState
func clinit_GoogBindToArrow() {
	if atomic.LoadUint32(&clinitDone_GoogBindToArrow) != 0 {
		return
	}
	clinitSlow_GoogBindToArrow()
}

func clinitSlow_GoogBindToArrow() {
	if !j2gort.EnterClinit(&clinitState_GoogBindToArrow, &clinitDone_GoogBindToArrow) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GoogBindToArrow, &clinitDone_GoogBindToArrow)
	static_GoogBindToArrow_GOOG_BIND_NAME = pkg_rhino.QualifiedName_Of("goog.bind")
}

func (this_ *GoogBindToArrow) initFields_GoogBindToArrow() {
	_ = this_
}

// Java constructor: GoogBindToArrow()
func NewGoogBindToArrow() *GoogBindToArrow {
	this_ := &GoogBindToArrow{}
	InitGoogBindToArrow(this_)
	return this_
}

func InitGoogBindToArrow(this_ *GoogBindToArrow) {
	pkg_jscomp.InitScanner(&this_.ScannerBase)
	this_.initFields_GoogBindToArrow()
}

// Java static method: GoogBindToArrow.isGoogBind(com.google.javascript.rhino.Node n)
func GoogBindToArrow_IsGoogBind(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsGetProp() && StaticGet_GoogBindToArrow_GOOG_BIND_NAME().Matches(n))
}

// Java method: GoogBindToArrow.matches(com.google.javascript.rhino.Node node, com.google.javascript.refactoring.NodeMetadata metadata)
func (this_ *GoogBindToArrow) Matches(node pkg_rhino.Node, metadata *pkg_jscomp.NodeMetadata) bool {
	_ = node
	_ = metadata
	if (!(node.IsCall())) {
		return false
	}
	var callee pkg_rhino.Node = node.GetFirstChild()
	_ = callee
	if (!(GoogBindToArrow_IsGoogBind(callee))) {
		return false
	}
	var firstArg pkg_rhino.Node = callee.GetNext()
	_ = firstArg
	if (((firstArg == 0) || !(firstArg.IsFunction()))) {
		return false
	}
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	return (((secondArg != 0) && secondArg.IsThis()) && (secondArg.GetNext() == 0))
}

// Java method: GoogBindToArrow.processMatch(com.google.javascript.refactoring.Match match)
func (this_ *GoogBindToArrow) ProcessMatch(match *pkg_jscomp.Match) any {
	_ = match
	var compiler pkg_jscomp.AbstractCompiler = match.GetMetadata().GetCompiler()
	_ = compiler
	var googBindCall pkg_rhino.Node = match.GetNode()
	_ = googBindCall
	var function pkg_rhino.Node = googBindCall.GetFirstChild().GetNext()
	_ = function
	var arrowFunction pkg_rhino.Node = function.CloneTree_0()
	_ = arrowFunction
	arrowFunction.SetIsArrowFunction(true)
	var body pkg_rhino.Node = function.GetLastChild()
	_ = body
	if (body.HasOneChild()) {
		var returnNode pkg_rhino.Node = body.GetFirstChild()
		_ = returnNode
		if (returnNode.IsReturn()) {
			arrowFunction.GetLastChild().ReplaceWith(returnNode.RemoveFirstChild())
		}
	}
	var fix *pkg_jscomp.SuggestedFix_Builder = pkg_jscomp.NewSuggestedFix_Builder()
	_ = fix
	fix.Replace(googBindCall, arrowFunction, compiler)
	return j2gort.ImmListOfPtr(fix.Build())
}

