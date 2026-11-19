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
// Port of com/google/javascript/jscomp/Var.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Var
type Var struct {
	AbstractVarBase
}

func (this_ *Var) BindSelf_(s any) {
	_ = s
}
func (this_ *Var) IsJavaNil_() bool           { return this_ == nil }
func (*Var) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Var", "Var" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Var", "com.google.javascript.jscomp.Var", func(v_ any) bool {
		return j2gort.InstanceOf[*Var](v_)
	})
}
func (*Var) Implements_jscomp_AbstractVar() {}
func (*Var) Implements_rhino_StaticSlot()   {}
func (*Var) Implements_rhino_StaticRef()    {}
func (*Var) Implements_jscomp_ScopedName()  {}

func StaticGet_Var_ARGUMENTS() string { return "arguments" }
func StaticGet_Var_EXPORTS() string   { return "exports" }
func (this_ *Var) initFields_Var() {
	_ = this_
}

// Java constructor: Var(java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node nameNode, com.google.javascript.jscomp.Scope scope, int index, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable CompilerInput input, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node implicitGoogNamespaceDefinition)
func NewVar(name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) *Var {
	this_ := &Var{}
	InitVar(this_, name, nameNode, scope, index, input, implicitGoogNamespaceDefinition)
	return this_
}
func allocOpt_Var(slab_ *j2gort.PooledSlab[Var], name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) *Var {
	if slab_ != nil && name != j2gort.NullString && scope != nil && index >= int32(-1) && implicitGoogNamespaceDefinition == 0 {
		okNode_ := nameNode == 0
		if !okNode_ {
			switch nameNode.TokenSwitchOrd_() {
			case int32(123), int32(32), int32(119):
				okNode_ = true
			}
		}
		if okNode_ {
			this_ := slab_.Alloc()
			this_.name = name
			this_.scopeRoot = scope.rootNode
			if nameNode != 0 {
				this_.nameNode = nameNode
			}
			this_.scope = scope
			if index != 0 {
				this_.index = index
			}
			if input != nil {
				this_.input = input
			}
			return this_
		}
	}
	var this_ *Var
	if slab_ != nil {
		this_ = slab_.Alloc()
	} else {
		this_ = &Var{}
	}
	InitVar(this_, name, nameNode, scope, index, input, implicitGoogNamespaceDefinition)
	return this_
}

func InitVar(this_ *Var, name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) {
	_ = name
	_ = nameNode
	_ = scope
	_ = index
	_ = input
	_ = implicitGoogNamespaceDefinition
	InitAbstractVar(&this_.AbstractVarBase, name, nameNode, ToIface_AbstractScope(scope), index, input, implicitGoogNamespaceDefinition)
	this_.initFields_Var()
	if nameNode != 0 {
		switch nameNode.TokenSwitchOrd_() {
		case int32(123) /* Token.MODULE_BODY */, int32(32) /* Token.NAME */, int32(119) /* Token.IMPORT_STAR */ :
		default:
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid name node "+j2gort.StrVal(nameNode))))
		}
	}
}

// Java static method: Var.createImplicitGoogNamespace(java.lang.String name, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node definition)
func Var_CreateImplicitGoogNamespace(name string, scope *Scope, definition pkg_rhino.Node) *Var {
	_ = name
	_ = scope
	_ = definition
	return NewVar(name, pkg_rhino.Node(0), scope, int32(-1), nil, definition)
}
func allocOptVia_Var_CreateImplicitGoogNamespace(slab_ *j2gort.PooledSlab[Var], name string, scope *Scope, definition pkg_rhino.Node) *Var {
	_ = name
	_ = scope
	_ = definition
	return allocOpt_Var(slab_, name, pkg_rhino.Node(0), scope, int32(-1), nil, definition)
}

// Java method: Var.getSymbol()
func (this_ *Var) GetSymbol() any {
	return this_
}

// Java method: Var.getDeclaration()
func (this_ *Var) GetDeclaration() any {
	if this_.GetNode() == 0 {
		return nil
	} else {
		return this_
	}
}

// Java method: Var.toString()
func (this_ *Var) String() string {
	return ("Var " + j2gort.StrValString(this_.GetName()) + " @ " + j2gort.StrVal(this_.GetNameNode()))
}
