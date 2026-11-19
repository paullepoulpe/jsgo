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
// Port of com/google/javascript/jscomp/TypedVar.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.TypedVar
type TypedVar struct {
	AbstractVarBase
	type_                     pkg_rhino.JSType `gson:"type"`
	markedEscaped             bool             `gson:"markedEscaped"`
	markedAssignedExactlyOnce bool             `gson:"markedAssignedExactlyOnce"`
	typeInferred              bool             `gson:"typeInferred"`
}

func (this_ *TypedVar) BindSelf_(s any) {
	_ = s
}
func (this_ *TypedVar) IsJavaNil_() bool { return this_ == nil }
func (*TypedVar) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypedVar", "TypedVar"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypedVar", "com.google.javascript.jscomp.TypedVar", func(v_ any) bool {
		return j2gort.InstanceOf[*TypedVar](v_)
	})
}
func (*TypedVar) Implements_jscomp_AbstractVar()    {}
func (*TypedVar) Implements_rhino_StaticTypedSlot() {}
func (*TypedVar) Implements_rhino_StaticTypedRef()  {}
func (*TypedVar) Implements_jscomp_ScopedName()     {}
func (*TypedVar) Implements_rhino_StaticSlot()      {}
func (*TypedVar) Implements_rhino_StaticRef()       {}

var static_TypedVar_NAME_NODE_TYPES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_TypedVar_NAME_NODE_TYPES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_TypedVar()
	return static_TypedVar_NAME_NODE_TYPES
}
func StaticSet_TypedVar_NAME_NODE_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_TypedVar()
	static_TypedVar_NAME_NODE_TYPES = v
}

var clinitDone_TypedVar uint32
var clinitState_TypedVar j2gort.ClinitState

func clinit_TypedVar() {
	if atomic.LoadUint32(&clinitDone_TypedVar) != 0 {
		return
	}
	clinitSlow_TypedVar()
}

func clinitSlow_TypedVar() {
	if !j2gort.EnterClinit(&clinitState_TypedVar, &clinitDone_TypedVar) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypedVar, &clinitDone_TypedVar)
	static_TypedVar_NAME_NODE_TYPES = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_NAME()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_THIS()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_IMPORT_STAR()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_EXPR_RESULT()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_EXPORT()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_GETPROP()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FUNCTION()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_STRING_KEY()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_SETTER_DEF()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_GETTER_DEF()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_MEMBER_FUNCTION_DEF()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_MODULE_BODY())})
}

func (this_ *TypedVar) initFields_TypedVar() {
	_ = this_

}

// Java constructor: TypedVar(boolean inferred, java.lang.String name, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.jscomp.TypedScope scope, int index, com.google.javascript.jscomp.CompilerInput input)
func NewTypedVar(inferred bool, name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, scope *TypedScope, index int32, input *CompilerInput) *TypedVar {
	this_ := &TypedVar{}
	InitTypedVar(this_, inferred, name, nameNode, type_, scope, index, input)
	return this_
}
func allocPooled_TypedVar(slab_ *j2gort.PooledSlab[TypedVar], inferred bool, name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, scope *TypedScope, index int32, input *CompilerInput) *TypedVar {
	clinit_TypedVar()
	this_ := slab_.Alloc()
	InitTypedVar(this_, inferred, name, nameNode, type_, scope, index, input)
	return this_
}

func InitTypedVar(this_ *TypedVar, inferred bool, name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, scope *TypedScope, index int32, input *CompilerInput) {
	_ = inferred
	_ = name
	_ = nameNode
	_ = type_
	_ = scope
	_ = index
	_ = input
	InitAbstractVar(&this_.AbstractVarBase, name, nameNode, ToIface_AbstractScope(scope), index, input, pkg_rhino.Node(0))
	this_.initFields_TypedVar()
	if nameNode != 0 {
		if !(j2gort.ImmSetContainsTypedU8_(StaticGet_TypedVar_NAME_NODE_TYPES(), uint8(nameNode.GetToken()))) {
			j2gort.CheckArgumentFmt(false, "Invalid name node token %s", j2gort.BoxNullableVal(nameNode.GetToken()))
		}
	}
	this_.type_ = type_
	this_.typeInferred = inferred
}

// Java method: TypedVar.getType()
func (this_ *TypedVar) GetType() pkg_rhino.JSType {
	return this_.type_
}

// Java method: TypedVar.setType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypedVar) SetType(type_ pkg_rhino.JSType) {
	_ = type_
	this_.type_ = type_
}

// Java method: TypedVar.resolveType(com.google.javascript.rhino.ErrorReporter errorReporter)
func (this_ *TypedVar) ResolveType(errorReporter pkg_rhino.ErrorReporter) {
	_ = errorReporter
	if this_.type_ != nil && !this_.type_.IsJavaNil_() {
		this_.type_ = this_.type_.GetJSTypeBase().Resolve(errorReporter)
	}
}

// Java method: TypedVar.isTypeInferred()
func (this_ *TypedVar) IsTypeInferred() bool {
	return this_.typeInferred
}

// Java method: TypedVar.getInputName()
func (this_ *TypedVar) GetInputName() string {
	if this_.GetInput() == nil {
		return "<non-file>"
	}
	return this_.GetInput().GetName()
}

// Java method: TypedVar.getSymbol()
func (this_ *TypedVar) GetSymbol() any {
	return this_
}

// Java method: TypedVar.getDeclaration()
func (this_ *TypedVar) GetDeclaration() any {
	if this_.GetNode() == 0 {
		return nil
	} else {
		return this_
	}
}

// Java method: TypedVar.toString()
func (this_ *TypedVar) String() string {
	return ("Var " + j2gort.StrValString(this_.GetName()) + "{" + j2gort.StrVal(this_.type_) + "}")
}

// Java method: TypedVar.markEscaped()
func (this_ *TypedVar) MarkEscaped() {
	this_.markedEscaped = true
}

// Java method: TypedVar.isMarkedEscaped()
func (this_ *TypedVar) IsMarkedEscaped() bool {
	return this_.markedEscaped
}

// Java method: TypedVar.markAssignedExactlyOnce()
func (this_ *TypedVar) MarkAssignedExactlyOnce() {
	this_.markedAssignedExactlyOnce = true
}

// Java method: TypedVar.isMarkedAssignedExactlyOnce()
func (this_ *TypedVar) IsMarkedAssignedExactlyOnce() bool {
	return this_.markedAssignedExactlyOnce
}
