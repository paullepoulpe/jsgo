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
// Port of com/google/javascript/jscomp/AbstractVar.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractVar
type AbstractVar interface {
	Implements_jscomp_AbstractVar()
	IsJavaNil_() bool
	GetAbstractVarBase() *AbstractVarBase
	ScopedName
	pkg_rhino.StaticSlot
	pkg_rhino.StaticRef
	GetInput() *CompilerInput
	GetParentNode() pkg_rhino.Node
	IsBleedingFunction() bool
	GetIndex() int32
	IsGlobal() bool
	IsLocal() bool
	IsExtern() bool
	IsDeclaredOrInferredConst() bool
	IsDefine() bool
	GetInitialValue() pkg_rhino.Node
	GetNameNode() pkg_rhino.Node
	IsVar() bool
	IsCatch() bool
	IsLet() bool
	IsConst() bool
	IsClass() bool
	IsParam() bool
	IsDefaultParam() bool
	IsImport() bool
	IsArguments() bool
	IsGoogModuleExports() bool
	IsThis() bool
	DeclarationType() pkg_rhino.Token
	IsImplicitGoogNamespace() bool
	AddImplicitGoogNamespaceDefinition(pkg_rhino.Node)
	GetImplicitGoogNamespaceStrength() pkg_rhino.StaticSourceFile_SourceKind
}

func ToIface_AbstractVar[T interface {
	AbstractVar
	comparable
}](v T) AbstractVar {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AbstractVar[T AbstractVar](v T) AbstractVar {
	var vi AbstractVar = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.AbstractVar
type AbstractVarBase struct {
	ScopedNameBase
	input                         *CompilerInput                        `gson:"input"`
	scope                         AbstractScope                         `gson:"scope"`
	nameNode                      pkg_rhino.Node                        `gson:"nameNode"`
	index                         int32                                 `gson:"index"`
	implicitGoogNamespaceStrength pkg_rhino.StaticSourceFile_SourceKind `gson:"implicitGoogNamespaceStrength"`
}

func (this_ *AbstractVarBase) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractVarBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractVarBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractVar", "AbstractVar"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractVar", "com.google.javascript.jscomp.AbstractVar", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractVar](v_)
	})
}
func (b *AbstractVarBase) GetAbstractVarBase() *AbstractVarBase { return b }
func (*AbstractVarBase) Implements_jscomp_AbstractVar()         {}
func (*AbstractVarBase) Implements_jscomp_ScopedName()          {}
func (*AbstractVarBase) Implements_rhino_StaticSlot()           {}
func (*AbstractVarBase) Implements_rhino_StaticRef()            {}

var static_AbstractVar_DECLARATION_TYPES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AbstractVar_DECLARATION_TYPES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AbstractVar()
	return static_AbstractVar_DECLARATION_TYPES
}
func StaticSet_AbstractVar_DECLARATION_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AbstractVar()
	static_AbstractVar_DECLARATION_TYPES = v
}

var clinitDone_AbstractVar uint32
var clinitState_AbstractVar j2gort.ClinitState

func clinit_AbstractVar() {
	if atomic.LoadUint32(&clinitDone_AbstractVar) != 0 {
		return
	}
	clinitSlow_AbstractVar()
}

func clinitSlow_AbstractVar() {
	if !j2gort.EnterClinit(&clinitState_AbstractVar, &clinitDone_AbstractVar) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AbstractVar, &clinitDone_AbstractVar)
	static_AbstractVar_DECLARATION_TYPES = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_VAR()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_LET()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CONST()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FUNCTION()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CLASS()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CATCH()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_IMPORT()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_PARAM_LIST())})
}

func (this_ *AbstractVarBase) initFields_AbstractVar() {
	_ = this_
}

func InitAbstractVar(this_ *AbstractVarBase, name string, nameNode pkg_rhino.Node, scope AbstractScope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) {
	_ = name
	_ = nameNode
	_ = scope
	_ = index
	_ = input
	_ = implicitGoogNamespaceDefinition
	InitScopedName(&this_.ScopedNameBase, j2gort.CheckNotNullString(name), func() pkg_rhino.Node {
		if scope != nil && !scope.IsJavaNil_() {
			return scope.GetAbstractScopeBase().GetRootNode_0()
		}
		return pkg_rhino.Node(0)
	}())
	this_.initFields_AbstractVar()
	if !(index >= int32(-1)) {
		j2gort.CheckArgumentMsg(false, j2gort.ValueOfInt(index))
	}
	if implicitGoogNamespaceDefinition != 0 {
		this_.nameNode = pkg_rhino.Node(0)
		this_.implicitGoogNamespaceStrength = AbstractVar_StrengthOf(implicitGoogNamespaceDefinition)
	} else {
		this_.nameNode = nameNode
		this_.implicitGoogNamespaceStrength = pkg_rhino.StaticSourceFile_SourceKind(0)
	}
	this_.scope = scope
	this_.index = index
	this_.input = input
}

// Java method: AbstractVar.getNode()
func (this_ *AbstractVarBase) GetNode() pkg_rhino.Node {
	return this_.nameNode
}

// Java method: AbstractVar.getInput()
func (this_ *AbstractVarBase) GetInput() *CompilerInput {
	return this_.input
}

// Java method: AbstractVar.getSourceFile()
func (this_ *AbstractVarBase) GetSourceFile() pkg_rhino.StaticSourceFile {
	return (func() pkg_rhino.Node {
		if this_.nameNode != 0 {
			return this_.nameNode
		}
		return this_.scope.GetAbstractScopeBase().GetRootNode_0()
	}()).GetStaticSourceFile()
}

// Java method: AbstractVar.getSymbol()
func (this_ *AbstractVarBase) GetSymbol() any {
	panic("abstract or native method invoked")
}

// Java method: AbstractVar.getDeclaration()
func (this_ *AbstractVarBase) GetDeclaration() any {
	panic("abstract or native method invoked")
}

// Java method: AbstractVar.getParentNode()
func (this_ *AbstractVarBase) GetParentNode() pkg_rhino.Node {
	if this_.nameNode == 0 {
		return pkg_rhino.Node(0)
	} else {
		return this_.nameNode.GetParent()
	}
}

// Java method: AbstractVar.isBleedingFunction()
func (this_ *AbstractVarBase) IsBleedingFunction() bool {
	var parent pkg_rhino.Node
	_ = parent
	if this_.nameNode == 0 {
		parent = pkg_rhino.Node(0)
	} else {
		parent = this_.nameNode.GetParent()
	}
	return ((parent != 0) && NodeUtil_IsFunctionExpression(parent))
}

// Java method: AbstractVar.getScope()
func (this_ *AbstractVarBase) GetScope() any {
	return this_.scope
}

// Java method: AbstractVar.getIndex()
func (this_ *AbstractVarBase) GetIndex() int32 {
	return this_.index
}

// Java method: AbstractVar.isGlobal()
func (this_ *AbstractVarBase) IsGlobal() bool {
	return this_.scope.IsGlobal()
}

// Java method: AbstractVar.isLocal()
func (this_ *AbstractVarBase) IsLocal() bool {
	return this_.scope.IsLocal()
}

// Java method: AbstractVar.isExtern()
func (this_ *AbstractVarBase) IsExtern() bool {
	return ((this_.input == nil) || this_.input.IsExtern())
}

// Java method: AbstractVar.isDeclaredOrInferredConst()
func (this_ *AbstractVarBase) IsDeclaredOrInferredConst() bool {
	var declarationNode pkg_rhino.Node = this_.nameNode
	_ = declarationNode
	if declarationNode == 0 {
		return false
	}
	return ((declarationNode.IsDeclaredConstantVar() || declarationNode.IsInferredConstantVar()) || declarationNode.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()))
}

// Java method: AbstractVar.isDefine()
func (this_ *AbstractVarBase) IsDefine() bool {
	var info *pkg_rhino.JSDocInfo
	_ = info
	if this_.nameNode == 0 {
		info = nil
	} else {
		info = NodeUtil_GetBestJSDocInfo(this_.nameNode)
	}
	return ((info != nil) && info.IsDefine())
}

// Java method: AbstractVar.getInitialValue()
func (this_ *AbstractVarBase) GetInitialValue() pkg_rhino.Node {
	if this_.nameNode == 0 {
		return pkg_rhino.Node(0)
	} else {
		return NodeUtil_GetRValueOfLValue(this_.nameNode)
	}
}

// Java method: AbstractVar.getNameNode()
func (this_ *AbstractVarBase) GetNameNode() pkg_rhino.Node {
	return this_.nameNode
}

// Java method: AbstractVar.getJSDocInfo()
func (this_ *AbstractVarBase) GetJSDocInfo() *pkg_rhino.JSDocInfo {
	if this_.nameNode == 0 {
		return nil
	} else {
		return NodeUtil_GetBestJSDocInfo(this_.nameNode)
	}
}

// Java method: AbstractVar.isVar()
func (this_ *AbstractVarBase) IsVar() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_VAR())
}

// Java method: AbstractVar.isCatch()
func (this_ *AbstractVarBase) IsCatch() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_CATCH())
}

// Java method: AbstractVar.isLet()
func (this_ *AbstractVarBase) IsLet() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_LET())
}

// Java method: AbstractVar.isConst()
func (this_ *AbstractVarBase) IsConst() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_CONST())
}

// Java method: AbstractVar.isClass()
func (this_ *AbstractVarBase) IsClass() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_CLASS())
}

// Java method: AbstractVar.isParam()
func (this_ *AbstractVarBase) IsParam() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_PARAM_LIST())
}

// Java method: AbstractVar.isDefaultParam()
func (this_ *AbstractVarBase) IsDefaultParam() bool {
	var parent pkg_rhino.Node = this_.nameNode.GetParent()
	_ = parent
	return ((parent.GetParent().IsParamList() && parent.IsDefaultValue()) && (parent.GetFirstChild() == this_.nameNode))
}

// Java method: AbstractVar.isImport()
func (this_ *AbstractVarBase) IsImport() bool {
	return (this_.priv_AbstractVar_declarationTypeImpl_530() == pkg_rhino.StaticGet_Token_IMPORT())
}

// Java method: AbstractVar.isArguments()
func (this_ *AbstractVarBase) IsArguments() bool {
	var n string = this_.GetScopedNameBase().name
	_ = n
	if !("arguments" == n) {
		return false
	}
	return this_.scope.GetAbstractScopeBase().IsFunctionScope()
}

// Java method: AbstractVar.isGoogModuleExports()
func (this_ *AbstractVarBase) IsGoogModuleExports() bool {
	var n string = this_.GetScopedNameBase().name
	_ = n
	if !("exports" == n) {
		return false
	}
	return (this_.scope.GetAbstractScopeBase().IsModuleScope() && this_.priv_AbstractVar_isImplicit_528())
}

// Java method: AbstractVar.isThis()
func (this_ *AbstractVarBase) IsThis() bool {
	var n string = this_.GetScopedNameBase().name
	_ = n
	if !("this" == n) {
		return false
	}
	return this_.scope.GetAbstractScopeBase().IsFunctionScope()
}

// Java method: AbstractVar.isImplicit()
func (this_ *AbstractVarBase) priv_AbstractVar_isImplicit_528() bool {
	if this_.implicitGoogNamespaceStrength != 0 {
		return true
	}
	var var_ AbstractScope_ImplicitVar = AbstractScope_ImplicitVar_Of(this_.GetScopedNameBase().name)
	_ = var_
	return ((var_ != 0) && var_.IsMadeByScope(this_.scope))
}

// Java method: AbstractVar.declarationType()
func (this_ *AbstractVarBase) DeclarationType() pkg_rhino.Token {
	return this_.priv_AbstractVar_declarationTypeImpl_530()
}

// Java method: AbstractVar.declarationTypeImpl()
func (this_ *AbstractVarBase) priv_AbstractVar_declarationTypeImpl_530() pkg_rhino.Token {
	if this_.implicitGoogNamespaceStrength != 0 {
		return pkg_rhino.Token(0)
	}
	{
		var current pkg_rhino.Node = this_.nameNode
		_ = current
		for ; current != 0; current = current.GetParent() {
			if j2gort.ImmSetContainsTypedU8_(StaticGet_AbstractVar_DECLARATION_TYPES(), uint8(current.GetToken())) {
				return current.GetToken()
			}
		}
	}
	if !(this_.priv_AbstractVar_isImplicit_528()) {
		j2gort.CheckStateFmt(false, "The nameNode for %s must be a descendant of one of: %s", j2gort.BoxString(this_.GetScopedNameBase().name), StaticGet_AbstractVar_DECLARATION_TYPES())
	}
	return pkg_rhino.Token(0)
}

// Java method: AbstractVar.isImplicitGoogNamespace()
func (this_ *AbstractVarBase) IsImplicitGoogNamespace() bool {
	return (this_.implicitGoogNamespaceStrength != 0)
}

// Java method: AbstractVar.addImplicitGoogNamespaceDefinition(com.google.javascript.rhino.Node definition)
func (this_ *AbstractVarBase) AddImplicitGoogNamespaceDefinition(definition pkg_rhino.Node) {
	_ = definition
	if !(this_.implicitGoogNamespaceStrength != 0) {
		j2gort.CheckStateStr(false, this_.GetScopedNameBase().name)
	}
	this_.implicitGoogNamespaceStrength = AbstractVar_StrongerOf(this_.implicitGoogNamespaceStrength, AbstractVar_StrengthOf(definition))
}

// Java method: AbstractVar.getImplicitGoogNamespaceStrength()
func (this_ *AbstractVarBase) GetImplicitGoogNamespaceStrength() pkg_rhino.StaticSourceFile_SourceKind {
	if !(this_.implicitGoogNamespaceStrength != 0) {
		j2gort.CheckStateStr(false, this_.GetScopedNameBase().name)
	}
	return this_.implicitGoogNamespaceStrength
}

// Java static method: AbstractVar.strengthOf(com.google.javascript.rhino.Node n)
func AbstractVar_StrengthOf(n pkg_rhino.Node) pkg_rhino.StaticSourceFile_SourceKind {
	_ = n
	var source pkg_rhino.StaticSourceFile = n.GetStaticSourceFile()
	_ = source
	if source == nil || source.IsJavaNil_() {
		return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()
	}
	return source.GetKind()
}

// Java static method: AbstractVar.strongerOf(com.google.javascript.rhino.StaticSourceFile.SourceKind left, com.google.javascript.rhino.StaticSourceFile.SourceKind right)
func AbstractVar_StrongerOf(left pkg_rhino.StaticSourceFile_SourceKind, right pkg_rhino.StaticSourceFile_SourceKind) pkg_rhino.StaticSourceFile_SourceKind {
	_ = left
	_ = right
	{
		var iifeCond_2 bool
		{

			re_1 := left
			ae_2 := pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
			if re_1 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_2 = re_1 == ae_2
		}
		if iifeCond_2 || func() bool {
			re_3 := right
			ae_4 := pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
			if re_3 == 0 {
				j2gort.ThrowNPE()
			}
			return re_3 == ae_4
		}() {
			return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
		} else {
			var iifeCond_1 bool
			{

				re_5 := left
				ae_6 := pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()
				if re_5 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_5 == ae_6
			}
			if iifeCond_1 || func() bool {
				re_7 := right
				ae_8 := pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()
				if re_7 == 0 {
					j2gort.ThrowNPE()
				}
				return re_7 == ae_8
			}() {
				return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()
			}
		}
	}
	return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_WEAK()
}
