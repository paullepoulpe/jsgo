/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InferConsts.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.InferConsts
type InferConsts struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *InferConsts) BindSelf_(s any) {
	_ = s
}
func (this_ *InferConsts) IsJavaNil_() bool { return this_ == nil }
func (*InferConsts) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InferConsts", "InferConsts"
}
func (this_ *InferConsts) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InferConsts) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InferConsts", "com.google.javascript.jscomp.InferConsts", func(v_ any) bool {
		return j2gort.InstanceOf[*InferConsts](v_)
	})
}
func (*InferConsts) Implements_jscomp_CompilerPass() {}

func (this_ *InferConsts) initFields_InferConsts() {
	_ = this_
}

// Java constructor: InferConsts(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewInferConsts(compiler AbstractCompiler) *InferConsts {
	this_ := &InferConsts{}
	InitInferConsts(this_, compiler)
	return this_
}

func InitInferConsts(this_ *InferConsts, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_InferConsts()
	this_.compiler = compiler
}

// Java method: InferConsts.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node js)
func (this_ *InferConsts) Process_2(externs pkg_rhino.Node, js pkg_rhino.Node) {
	_ = externs
	_ = js
	clinit_ReferenceCollector()
	var val_1 ReferenceCollector
	InitReferenceCollector_3(&val_1, this_.compiler, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(this_.compiler))
	var collector *ReferenceCollector = &val_1
	_ = collector
	collector.ownsScopeCreator_ = true
	defer collector.ReleaseSlabs_()

	collector.Process_1(js)
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](collector.GetAllSymbols())))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			v := j2gort.CastPtr[Var](rawElem_3)
			_ = v
			this_.priv_InferConsts_considerVar_3214(v, j2gort.CastPtr[ReferenceCollection](collector.GetReferences(j2gort.PtrToAny(v))))
		}
	}
	scInline_1 := NewSyntacticScopeCreator_1(this_.compiler)
	defer scInline_1.ReleaseSlabs_()
	var globalExternsScope *Scope = scInline_1.CreateScope_2_Node_Scope(externs, nil)

	_ = globalExternsScope
	{
		cur_4 := (globalExternsScope).iterVars_()
		for cur_4.HasNext() {

			v_2 := cur_4.Next()

			_ = v_2
			this_.priv_InferConsts_considerVar_3214(v_2, nil)
		}
	}
}

// Java method: InferConsts.considerVar(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable ReferenceCollection refCollection)
func (this_ *InferConsts) priv_InferConsts_considerVar_3214(v *Var, refCollection *ReferenceCollection) {
	_ = v
	_ = refCollection
	if v.IsImplicitGoogNamespace() {
		return
	}
	var nameNode pkg_rhino.Node = v.GetNameNode()
	_ = nameNode
	var docInfo *pkg_rhino.JSDocInfo = v.GetJSDocInfo()
	_ = docInfo
	if (docInfo != nil) && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	} else if (nameNode != 0) && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	}
	if InferConsts_IsInferredConst(v, refCollection) {
		nameNode.SetInferredConstantVar(true)
	}
}

// Java static method: InferConsts.isInferredConst(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.ReferenceCollection refCollection)
func InferConsts_IsInferredConst(v *Var, refCollection *ReferenceCollection) bool {
	_ = v
	_ = refCollection
	var nameNode pkg_rhino.Node = v.GetNameNode()
	_ = nameNode
	if ((nameNode == 0) || (refCollection == nil)) || !(refCollection.IsAssignedOnceInLifetime()) {
		return false
	}
	if v.IsImplicitGoogNamespace() {
		return false
	}
	var declarationType pkg_rhino.Token = v.DeclarationType()
	_ = declarationType
	switch (declarationType).Ordinal() {
	case int32(111) /* Token.LET */ :
		return (!(nameNode.GetParent().IsLet()) || nameNode.HasChildren())
	case int32(97) /* Token.CONST */, int32(86) /* Token.CATCH */, int32(105) /* Token.CLASS */, int32(49) /* Token.PARAM_LIST */, int32(73) /* Token.FUNCTION */ :
		return true
	case int32(116) /* Token.IMPORT */ :
		return false
	case int32(84) /* Token.VAR */ :
		return (refCollection.FirstReferenceIsAssigningDeclaration() && refCollection.IsWellDefined())
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unrecognized declaration type "+j2gort.StrVal(declarationType))))
	}
	panic("unreachable switch expression")
}
