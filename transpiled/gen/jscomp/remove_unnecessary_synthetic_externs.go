/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RemoveUnnecessarySyntheticExterns.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns
type RemoveUnnecessarySyntheticExterns struct {
	compiler      AbstractCompiler                  `gson:"compiler"`
	nodesToDetach j2gort.OrderedSet[pkg_rhino.Node] `gson:"nodesToDetach"`
	idHash_       atomic.Int32
}

func (this_ *RemoveUnnecessarySyntheticExterns) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnnecessarySyntheticExterns) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnnecessarySyntheticExterns) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns", "RemoveUnnecessarySyntheticExterns"
}
func (this_ *RemoveUnnecessarySyntheticExterns) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnnecessarySyntheticExterns) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns", "com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnnecessarySyntheticExterns](v_)
	})
}
func (*RemoveUnnecessarySyntheticExterns) Implements_jscomp_CompilerPass() {}

func (this_ *RemoveUnnecessarySyntheticExterns) initFields_RemoveUnnecessarySyntheticExterns() {
	_ = this_
	this_.nodesToDetach = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: RemoveUnnecessarySyntheticExterns(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRemoveUnnecessarySyntheticExterns(compiler AbstractCompiler) *RemoveUnnecessarySyntheticExterns {
	this_ := &RemoveUnnecessarySyntheticExterns{}
	InitRemoveUnnecessarySyntheticExterns(this_, compiler)
	return this_
}

func InitRemoveUnnecessarySyntheticExterns(this_ *RemoveUnnecessarySyntheticExterns, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RemoveUnnecessarySyntheticExterns()
	this_.compiler = compiler
}

// Java method: RemoveUnnecessarySyntheticExterns.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RemoveUnnecessarySyntheticExterns) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var globalRoot pkg_rhino.Node = externs.GetParent()
	_ = globalRoot
	clinit_SyntacticScopeCreator()
	var val_1 SyntacticScopeCreator
	InitSyntacticScopeCreator_3(&val_1, this_.compiler, NewRemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler(this_), true)
	var scopeCreator *SyntacticScopeCreator = &val_1
	_ = scopeCreator
	defer scopeCreator.ReleaseSlabs_()

	scopeCreator.CreateScope_2_Node_Scope(globalRoot, nil)
	if this_.nodesToDetach.IsEmpty() {
		return
	}
	{
		cur_2 := this_.nodesToDetach.Iter()
		for cur_2.HasNext() {
			n := cur_2.
				Next()
			_ = n
			n.Detach()
		}
	}
	var syntheticExternsScript pkg_rhino.Node = this_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.compiler)
	_ = syntheticExternsScript
	this_.compiler.ReportChangeToEnclosingScope(syntheticExternsScript)
}

// Java class: com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns.RedeclarationCheckHandler
type RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler struct {
	outer_  *RemoveUnnecessarySyntheticExterns
	idHash_ atomic.Int32
}

func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) IsJavaNil_() bool {
	return this_ == nil
}
func (*RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns$RedeclarationCheckHandler", "RedeclarationCheckHandler"
}
func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns.RedeclarationCheckHandler", "com.google.javascript.jscomp.RemoveUnnecessarySyntheticExterns$RedeclarationCheckHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler](v_)
	})
}
func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) SetOuter_(o *RemoveUnnecessarySyntheticExterns) {
	this_.outer_ = o
}
func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) GetOuter_() *RemoveUnnecessarySyntheticExterns {
	return this_.outer_
}
func (*RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) Implements_jscomp_SyntacticScopeCreator_RedeclarationHandler() {
}

func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) initFields_RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler() {
	_ = this_
}

// Java constructor: RedeclarationCheckHandler()
func NewRemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler(outer_ *RemoveUnnecessarySyntheticExterns) *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler {
	this_ := &RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler{}
	this_.outer_ = outer_
	InitRemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler(this_)
	return this_
}

func InitRemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler(this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) {
	this_.initFields_RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler()
}

// Java method: RedeclarationCheckHandler.onRedeclaration(com.google.javascript.jscomp.Scope s, java.lang.String name, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CompilerInput input)
func (this_ *RemoveUnnecessarySyntheticExterns_RedeclarationCheckHandler) OnRedeclaration(s *Scope, name string, n pkg_rhino.Node, input *CompilerInput) {
	_ = s
	_ = name
	_ = n
	_ = input
	if !(s.IsGlobal()) {
		j2gort.CheckArgumentFmt(false, "Unexpected non-global redeclaration %s in scope %s", j2gort.BoxNode(n), j2gort.PtrToAny(s))
	}
	var isSynthetic bool = func() bool {
		rp_3 := input
		ap_4 := this_.outer_.compiler.GetSynthesizedExternsInput()
		return j2gort.CheckNotNullPtr(rp_3) == ap_4
	}()
	_ = isSynthetic
	var origVar *Var = j2gort.CastPtr[Var](s.GetVar(name))
	_ = origVar
	var origNode pkg_rhino.Node = origVar.GetNode()
	_ = origNode
	if isSynthetic {
		if n.GetParent().IsSynthesizedUnfulfilledNameDeclaration() {
			this_.outer_.nodesToDetach.Add(n.GetParent())
		}
	} else if (origNode != 0) && origNode.GetParent().IsSynthesizedUnfulfilledNameDeclaration() {
		this_.outer_.nodesToDetach.Add(origNode.GetParent())
	}
}
