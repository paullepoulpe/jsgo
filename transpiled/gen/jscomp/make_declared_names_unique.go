/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/MakeDeclaredNamesUnique.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique
type MakeDeclaredNamesUnique struct {
	NodeTraversal_AbstractScopedCallbackBase
	renamerStack           j2gort.Deque[MakeDeclaredNamesUnique_Renamer] `gson:"renamerStack"`
	rootRenamer            MakeDeclaredNamesUnique_Renamer               `gson:"rootRenamer"`
	markChanges            bool                                          `gson:"markChanges"`
	assertOnChange         bool                                          `gson:"assertOnChange"`
	slab_ContextualRenamer j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenamerBase]
}

func (this_ *MakeDeclaredNamesUnique) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique", "MakeDeclaredNamesUnique"
}
func (this_ *MakeDeclaredNamesUnique) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*MakeDeclaredNamesUnique) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique", "com.google.javascript.jscomp.MakeDeclaredNamesUnique", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique](v_)
	})
}
func (*MakeDeclaredNamesUnique) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {}
func (*MakeDeclaredNamesUnique) Implements_jscomp_NodeTraversal_ScopedCallback()         {}
func (*MakeDeclaredNamesUnique) Implements_jscomp_NodeTraversal_Callback()               {}

func StaticGet_MakeDeclaredNamesUnique_ARGUMENTS() string { return "arguments" }
func (this_ *MakeDeclaredNamesUnique) initFields_MakeDeclaredNamesUnique() {
	_ = this_
	this_.renamerStack = j2gort.Deque[MakeDeclaredNamesUnique_Renamer]{}
}

// Java constructor: MakeDeclaredNamesUnique(com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer renamer, boolean markChanges, boolean assertOnChange)
func NewMakeDeclaredNamesUnique(renamer MakeDeclaredNamesUnique_Renamer, markChanges bool, assertOnChange bool) *MakeDeclaredNamesUnique {
	this_ := &MakeDeclaredNamesUnique{}
	InitMakeDeclaredNamesUnique(this_, renamer, markChanges, assertOnChange)
	return this_
}

func InitMakeDeclaredNamesUnique(this_ *MakeDeclaredNamesUnique, renamer MakeDeclaredNamesUnique_Renamer, markChanges bool, assertOnChange bool) {
	_ = renamer
	_ = markChanges
	_ = assertOnChange
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_MakeDeclaredNamesUnique()
	this_.rootRenamer = renamer
	this_.markChanges = markChanges
	this_.assertOnChange = assertOnChange
}

// Java static method: MakeDeclaredNamesUnique.builder()
func MakeDeclaredNamesUnique_S_Builder() *MakeDeclaredNamesUnique_Builder {
	return NewMakeDeclaredNamesUnique_Builder()
}

// Java static method: MakeDeclaredNamesUnique.getContextualRenameInverter(com.google.javascript.jscomp.AbstractCompiler compiler)
func MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler AbstractCompiler) CompilerPass {
	_ = compiler
	return NewMakeDeclaredNamesUnique_ContextualRenameInverter(compiler)
}

// Java method: MakeDeclaredNamesUnique.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	_ = t
	var declarationRoot pkg_rhino.Node = t.GetScopeRoot()
	_ = declarationRoot
	var renamer MakeDeclaredNamesUnique_Renamer = nil
	_ = renamer
	if this_.renamerStack.IsEmpty() {
		if !(!(declarationRoot.IsFunction()) || !(j2gort.InstanceOf[MakeDeclaredNamesUnique_ContextualRenamer](this_.rootRenamer))) {
			j2gort.CheckState(false)
		}
		renamer = this_.rootRenamer
	} else {
		var hoist bool = (!(declarationRoot.IsFunction()) && !(NodeUtil_CreatesBlockScope(declarationRoot)))
		_ = hoist
		topRenamer_ := this_.renamerStack.Peek()
		if crBase_, okCr_ := topRenamer_.(*MakeDeclaredNamesUnique_ContextualRenamerBase); okCr_ {
			renamer = crBase_.createForChildScopeWithSlab_(&this_.slab_ContextualRenamer, t.GetScopeRoot(), hoist)
		} else {
			renamer = topRenamer_.CreateForChildScope(t.GetScopeRoot(), hoist)
		}

	}
	this_.renamerStack.Push(renamer)
	this_.priv_MakeDeclaredNamesUnique_findDeclaredNames_3807(t, declarationRoot)
}

// Java method: MakeDeclaredNamesUnique.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	_ = t
	if !(t.InGlobalScope()) {
		popped_ := this_.renamerStack.Pop()
		if crPopped_, okCr_ := popped_.(*MakeDeclaredNamesUnique_ContextualRenamerBase); okCr_ && !crPopped_.global {
			if crPopped_.declarations != nil {
				j2gort.ReleaseOrderedMapPtr(crPopped_.declarations)
				crPopped_.declarations = nil
			}
			if this_.renamerStack.Size() <= 1 {
				this_.slab_ContextualRenamer.Reset()
			}
		}

	}
}

// Java method: MakeDeclaredNamesUnique.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *MakeDeclaredNamesUnique) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */, int32(119) /* Token.IMPORT_STAR */ :
		this_.priv_MakeDeclaredNamesUnique_visitNameOrImportStar_3805(t, n, parent)
	default:
	}
}

// Java method: MakeDeclaredNamesUnique.visitNameOrImportStar(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *MakeDeclaredNamesUnique) priv_MakeDeclaredNamesUnique_visitNameOrImportStar_3805(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() && NodeUtil_IsNonlocalModuleExportName(n) {
		return
	}
	var newName string = this_.priv_MakeDeclaredNamesUnique_getReplacementName_3806(n.GetString())
	_ = newName
	if newName == j2gort.NullString {
		return
	}
	if this_.assertOnChange {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected renaming in MakeDeclaredNamesUnique. new name: " + j2gort.StrValString(newName) + " for " + j2gort.StrVal(n))))))
	}
	var renamer MakeDeclaredNamesUnique_Renamer = this_.renamerStack.Peek()
	_ = renamer
	if renamer.StripConstIfReplaced() {
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), false)
		var jsDocInfoNode pkg_rhino.Node = NodeUtil_GetBestJSDocInfoNode(n)
		_ = jsDocInfoNode
		if (jsDocInfoNode != 0) && (jsDocInfoNode.GetJSDocInfo() != nil) {
			var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_CopyFrom(jsDocInfoNode.GetJSDocInfo())
			_ = builder
			builder.RecordMutable()
			jsDocInfoNode.SetJSDocInfo(builder.Build_0())
		}
	}
	n.SetString(newName)
	if this_.markChanges {
		t.ReportCodeChange_0()
		if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
			t.GetCompiler().ReportChangeToEnclosingScope(parent)
		}
	}
}

// Java method: MakeDeclaredNamesUnique.getReplacementName(java.lang.String oldName)
func (this_ *MakeDeclaredNamesUnique) priv_MakeDeclaredNamesUnique_getReplacementName_3806(oldName string) string {
	_ = oldName
	{
		cur_1 := this_.renamerStack.Iter()
		for cur_1.HasNext() {

			renamer := cur_1.Next()
			_ = renamer
			var newName string = renamer.GetReplacementName(oldName)
			_ = newName
			if newName != j2gort.NullString {
				return newName
			}
		}
	}
	return j2gort.NullString
}

// Java method: MakeDeclaredNamesUnique.findDeclaredNames(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *MakeDeclaredNamesUnique) priv_MakeDeclaredNamesUnique_findDeclaredNames_3807(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(NodeUtil_CreatesScope(n) || n.IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	{
		cur_3 := (t.GetScope()).iterVars_()
		for cur_3.HasNext() {

			v := cur_3.Next()

			_ = v
			this_.renamerStack.Peek().AddDeclaredName(v.GetName(), false)
		}
	}
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.Builder
type MakeDeclaredNamesUnique_Builder struct {
	renamer        MakeDeclaredNamesUnique_Renamer `gson:"renamer"`
	idHash_        atomic.Int32
	markChanges    bool `gson:"markChanges"`
	assertOnChange bool `gson:"assertOnChange"`
}

func (this_ *MakeDeclaredNamesUnique_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_Builder) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$Builder", "Builder"
}
func (this_ *MakeDeclaredNamesUnique_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.Builder", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_Builder](v_)
	})
}

func (this_ *MakeDeclaredNamesUnique_Builder) initFields_MakeDeclaredNamesUnique_Builder() {
	_ = this_
	this_.markChanges = true

}

// Java constructor: Builder()
func NewMakeDeclaredNamesUnique_Builder() *MakeDeclaredNamesUnique_Builder {
	this_ := &MakeDeclaredNamesUnique_Builder{}
	InitMakeDeclaredNamesUnique_Builder(this_)
	return this_
}

func InitMakeDeclaredNamesUnique_Builder(this_ *MakeDeclaredNamesUnique_Builder) {
	this_.initFields_MakeDeclaredNamesUnique_Builder()
}

// Java method: Builder.withRenamer(com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer renamer)
func (this_ *MakeDeclaredNamesUnique_Builder) WithRenamer(renamer MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_Builder {
	_ = renamer
	this_.renamer = renamer
	return this_
}

// Java method: Builder.withMarkChanges(boolean markChanges)
func (this_ *MakeDeclaredNamesUnique_Builder) WithMarkChanges(markChanges bool) *MakeDeclaredNamesUnique_Builder {
	_ = markChanges
	this_.markChanges = markChanges
	return this_
}

// Java method: Builder.withAssertOnChange(boolean assertOnChange)
func (this_ *MakeDeclaredNamesUnique_Builder) WithAssertOnChange(assertOnChange bool) *MakeDeclaredNamesUnique_Builder {
	_ = assertOnChange
	this_.assertOnChange = assertOnChange
	return this_
}

// Java method: Builder.build()
func (this_ *MakeDeclaredNamesUnique_Builder) Build() *MakeDeclaredNamesUnique {
	if this_.renamer == nil || this_.renamer.IsJavaNil_() {
		this_.renamer = Upcast_MakeDeclaredNamesUnique_Renamer(NewMakeDeclaredNamesUnique_ContextualRenamer())
	}
	return NewMakeDeclaredNamesUnique(this_.renamer, this_.markChanges, this_.assertOnChange)
}

// Java interface: com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer
type MakeDeclaredNamesUnique_Renamer interface {
	Implements_jscomp_MakeDeclaredNamesUnique_Renamer()
	IsJavaNil_() bool
	AddDeclaredName(string, bool)
	GetReplacementName(string) string
	StripConstIfReplaced() bool
	CreateForChildScope(pkg_rhino.Node, bool) MakeDeclaredNamesUnique_Renamer
	GetHoistRenamer() MakeDeclaredNamesUnique_Renamer
}

func ToIface_MakeDeclaredNamesUnique_Renamer[T interface {
	MakeDeclaredNamesUnique_Renamer
	comparable
}](v T) MakeDeclaredNamesUnique_Renamer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_MakeDeclaredNamesUnique_Renamer[T MakeDeclaredNamesUnique_Renamer](v T) MakeDeclaredNamesUnique_Renamer {
	var vi MakeDeclaredNamesUnique_Renamer = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$Renamer", func(v_ any) bool {
		return j2gort.InstanceOf[MakeDeclaredNamesUnique_Renamer](v_)
	})
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter
type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	compiler                    AbstractCompiler                                                             `gson:"compiler"`
	scopeContextStack           j2gort.Deque[*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext] `gson:"scopeContextStack"`
	variablesStackedByName      *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap       `gson:"variablesStackedByName"`
	idHash_                     atomic.Int32
	slab_ScopeContext           j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext]
	slab_RenamingInfo           j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo]
	slab_DeclaredVariableInfo   j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo]
	slab_UndeclaredVariableInfo j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo]
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter", "ContextualRenameInverter"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenameInverter](v_)
	})
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) Implements_jscomp_CompilerPass()           {}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) initFields_MakeDeclaredNamesUnique_ContextualRenameInverter() {
	_ = this_
	this_.scopeContextStack = j2gort.Deque[*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext]{}
	this_.variablesStackedByName = NewMakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap(this_)
}

// Java constructor: ContextualRenameInverter(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter(compiler AbstractCompiler) *MakeDeclaredNamesUnique_ContextualRenameInverter {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter(this_, compiler)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenameInverter(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenameInverter()
	this_.compiler = compiler
}

// Java method: ContextualRenameInverter.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node js)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) Process_2(externs pkg_rhino.Node, js pkg_rhino.Node) {
	_ = externs
	_ = js
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		this_.scopeContextStack.Clear()
		this_.variablesStackedByName.Clear()
		sc_.ReleaseSlabs_()
		this_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(js)
	}

}

// Java static method: ContextualRenameInverter.getOriginalName(java.lang.String name)
func MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name string) string {
	_ = name
	var index int32 = MakeDeclaredNamesUnique_ContextualRenameInverter_IndexOfSeparator(name)
	_ = index
	if index <= int32(0) {
		return name
	} else {
		return j2gort.JStrSub(name, int32(0), index)
	}
}

// Java static method: ContextualRenameInverter.indexOfSeparator(java.lang.String name)
func MakeDeclaredNamesUnique_ContextualRenameInverter_IndexOfSeparator(name string) int32 {
	_ = name
	return j2gort.JStrLastIndexOf(name, "$jscomp$")
}

// Java method: ContextualRenameInverter.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalScope() {
		return
	}
	var scope *Scope = t.GetScope()
	_ = scope
	var declaredVariablesBuilder []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo = []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo{}
	_ = declaredVariablesBuilder
	{
		cur_5 := (scope).iterVars_()
		for cur_5.HasNext() {

			var_ := cur_5.Next()

			_ = var_
			var variableInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo = this_.priv_ContextualRenameInverter_createDeclaredVariableInfo_3817(var_)
			_ = variableInfo
			declaredVariablesBuilder = append(declaredVariablesBuilder, j2gort.CheckNotNullElem(variableInfo))
			this_.variablesStackedByName.PushVariableInfo(ToIface_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo(variableInfo))
		}
	}
	var parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext = this_.scopeContextStack.Peek()
	_ = parent
	this_.scopeContextStack.Push(allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(&this_.slab_ScopeContext, scope, j2gort.ImmListFromTypedSlice(declaredVariablesBuilder), parent))
}

// Java method: ContextualRenameInverter.createDeclaredVariableInfo(com.google.javascript.jscomp.Var var)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) priv_ContextualRenameInverter_createDeclaredVariableInfo_3817(var_ *Var) *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo {
	_ = var_
	var currentName string = var_.GetName()
	_ = currentName
	var indexOfSeparator int32 = MakeDeclaredNamesUnique_ContextualRenameInverter_IndexOfSeparator(currentName)
	_ = indexOfSeparator
	var renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo = nil
	_ = renamingInfo
	if indexOfSeparator > int32(0) {
		var invertedName string = j2gort.JStrSub(currentName, int32(0), indexOfSeparator)
		_ = invertedName
		renamingInfo = allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(&this_.slab_RenamingInfo, this_, currentName, invertedName)
	} else {
		renamingInfo = nil
	}
	return allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(&this_.slab_DeclaredVariableInfo, var_, renamingInfo)
}

// Java method: ContextualRenameInverter.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalScope() {
		return
	}
	var scopeContext *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext = j2gort.CheckNotNullPtr(this_.scopeContextStack.Pop())
	_ = scopeContext
	{
		cur_7 := scopeContext.referencedVariables.Iter()
		for cur_7.HasNext() {

			referencedVariable := cur_7.
				Next()
			_ = referencedVariable
			scopeContext.RecordPotentialConflicts(referencedVariable)
		}
	}
	{

		for _, declaredVariableInfo := range scopeContext.declaredVariableInfos {

			_ = declaredVariableInfo
			var variableName string = declaredVariableInfo.GetName()
			_ = variableName
			var oldStackTop MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo = this_.variablesStackedByName.PopVariableInfo(variableName)
			_ = oldStackTop
			if !(oldStackTop == declaredVariableInfo || ((oldStackTop == nil || oldStackTop.IsJavaNil_()) && (declaredVariableInfo == nil || declaredVariableInfo.IsJavaNil_()))) {
				j2gort.CheckStateFmt(false, "Declared variable \"%s\" was not the top of the name stack", j2gort.BoxString(variableName))
			}
			declaredVariableInfo.priv_DeclaredVariableInfo_tryToInvertName_3825()
		}
	}
	if this_.scopeContextStack.IsEmpty() {
		this_.variablesStackedByName.Clear()
	}
}

// Java method: ContextualRenameInverter.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}

// Java method: ContextualRenameInverter.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) Visit_3(t *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = node
	_ = parent
	var scopeContext *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext = this_.scopeContextStack.Peek()
	_ = scopeContext
	if scopeContext != nil {
		if NodeUtil_IsReferenceName(node) || node.IsImportStar() {
			var name string = node.GetString()
			_ = name
			var variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo = this_.variablesStackedByName.GetOrCreateCurrentVariableInfo(name)
			_ = variableInfo
			variableInfo.AddReferenceNode(node)
			scopeContext.AddReferencedVariableInfo(variableInfo)
		}
	}
}

// Java method: ContextualRenameInverter.createUndeclaredVariableInfo(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) priv_ContextualRenameInverter_createUndeclaredVariableInfo_3818(name string) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	_ = name
	return ToIface_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo(allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(&this_.slab_UndeclaredVariableInfo, name))
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo
type MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo struct {
	outer_                   *MakeDeclaredNamesUnique_ContextualRenameInverter
	referenceNodes           []pkg_rhino.Node                                                                 `gson:"referenceNodes"`
	potentialShadowVariables j2gort.OrderedSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo] `gson:"potentialShadowVariables"`
	currentName              string                                                                           `gson:"currentName"`
	preferredName            string                                                                           `gson:"preferredName"`
	idHash_                  atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) IsJavaNil_() bool {
	return this_ == nil
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$RenamingInfo", "RenamingInfo"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$RenamingInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo](v_)
	})
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) SetOuter_(o *MakeDeclaredNamesUnique_ContextualRenameInverter) {
	this_.outer_ = o
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) GetOuter_() *MakeDeclaredNamesUnique_ContextualRenameInverter {
	return this_.outer_
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo() {
	_ = this_
	this_.referenceNodes = []pkg_rhino.Node{}
	this_.potentialShadowVariables = j2gort.NewOrderedSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]()
	this_.currentName = j2gort.NullString
	this_.preferredName = j2gort.NullString
}

// Java constructor: RenamingInfo(java.lang.String currentName, java.lang.String preferredName)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(outer_ *MakeDeclaredNamesUnique_ContextualRenameInverter, currentName string, preferredName string) *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo{}
	this_.outer_ = outer_
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(this_, currentName, preferredName)
	return this_
}
func allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(slab_ *j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo], outer_ *MakeDeclaredNamesUnique_ContextualRenameInverter, currentName string, preferredName string) *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(this_, currentName, preferredName)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo, currentName string, preferredName string) {
	_ = currentName
	_ = preferredName
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo()
	this_.currentName = currentName
	this_.preferredName = preferredName
}

// Java method: RenamingInfo.addPotentialShadowVariable(com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo variableInfo)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) AddPotentialShadowVariable(variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	_ = variableInfo
	this_.potentialShadowVariables.Add(variableInfo)
}

// Java method: RenamingInfo.attemptRename()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) AttemptRename() string {
	var disallowedNames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = disallowedNames
	disallowedNames.Add("arguments")
	{
		cur_11 := this_.potentialShadowVariables.Iter()
		for cur_11.HasNext() {

			potentialShadowVariable := cur_11.
				Next()
			_ = potentialShadowVariable
			var potentialShadowName string = potentialShadowVariable.GetName()
			_ = potentialShadowName
			disallowedNames.Add(potentialShadowName)
		}
	}
	var newName string = this_.preferredName
	_ = newName
	if disallowedNames.Contains(this_.preferredName) {
		var baseName string = (j2gort.StrValString(this_.preferredName) + "$jscomp$")
		_ = baseName
		var i int32 = int32(0)
		_ = i
		for doFirst_13 := true; doFirst_13 || (disallowedNames.Contains(newName)); doFirst_13 = false {
			newName = (j2gort.StrValString(baseName) + j2gort.StrValInt32(func() int32 { old_14 := i; nw_15 := old_14 + 1; i = nw_15; return old_14 }()))
		}
	}
	{
		var iifeCond_1 bool
		{

			sr_16 := newName
			se_17 := this_.currentName
			iifeCond_1 = j2gort.CheckNotNullString(sr_16) == se_17
		}
		if iifeCond_1 {
			return this_.currentName
		}
	}
	for _, referenceNode := range this_.referenceNodes {
		_ = referenceNode
		referenceNode.SetString(newName)
		this_.outer_.compiler.ReportChangeToEnclosingScope(referenceNode)
		var parent pkg_rhino.Node = referenceNode.GetParent()
		_ = parent
		if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
			this_.outer_.compiler.ReportChangeToEnclosingScope(parent)
		}
	}

	return newName
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo
type MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo struct {
	var_         *Var                                                           `gson:"var"`
	name         string                                                         `gson:"name"`
	renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo `gson:"renamingInfo"`
	idHash_      atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) IsJavaNil_() bool {
	return this_ == nil
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$DeclaredVariableInfo", "DeclaredVariableInfo"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$DeclaredVariableInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo](v_)
	})
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo() {
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: DeclaredVariableInfo(com.google.javascript.jscomp.Var var, com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.@org.jspecify.annotations.Nullable RenamingInfo renamingInfo)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(var_ *Var, renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(this_, var_, renamingInfo)
	return this_
}
func allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(slab_ *j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo], var_ *Var, renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo {
	this_ := slab_.Alloc()
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(this_, var_, renamingInfo)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo, var_ *Var, renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) {
	_ = var_
	_ = renamingInfo
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo()
	this_.var_ = var_
	this_.name = var_.GetName()
	this_.renamingInfo = renamingInfo
}

// Java method: DeclaredVariableInfo.getName()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) GetName() string {
	return this_.name
}

// Java method: DeclaredVariableInfo.getScopeDepth()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) GetScopeDepth() int32 {
	return j2gort.CastPtr[Scope](this_.var_.GetScope()).GetDepth()
}

// Java method: DeclaredVariableInfo.addPotentialConflict(com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo variableInfo)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) AddPotentialConflict(variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	_ = variableInfo
	if (this_.renamingInfo != nil) && !(variableInfo == this_ || ((variableInfo == nil || variableInfo.IsJavaNil_()) && (this_ == nil || this_.IsJavaNil_()))) {
		this_.renamingInfo.AddPotentialShadowVariable(variableInfo)
	}
}

// Java method: DeclaredVariableInfo.addReferenceNode(com.google.javascript.rhino.Node nameNode)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) AddReferenceNode(nameNode pkg_rhino.Node) {
	_ = nameNode
	if this_.renamingInfo != nil {
		this_.renamingInfo.referenceNodes = append(this_.renamingInfo.referenceNodes, nameNode)
	}
}

// Java method: DeclaredVariableInfo.tryToInvertName()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) priv_DeclaredVariableInfo_tryToInvertName_3825() {
	if this_.renamingInfo != nil {
		this_.name = this_.renamingInfo.AttemptRename()
		this_.renamingInfo = nil
	}
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext
type MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext struct {
	parent                *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext                   `gson:"parent"`
	scope                 *Scope                                                                           `gson:"scope"`
	declaredVariableInfos []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo         `gson:"declaredVariableInfos"`
	referencedVariables   j2gort.OrderedSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo] `gson:"referencedVariables"`
	idHash_               atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) IsJavaNil_() bool {
	return this_ == nil
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$ScopeContext", "ScopeContext"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$ScopeContext", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext](v_)
	})
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext() {
	_ = this_
	this_.referencedVariables = j2gort.NewOrderedSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo](

	// Java constructor: ScopeContext(com.google.javascript.jscomp.Scope scope, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo> declaredVariableInfos, com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.@org.jspecify.annotations.Nullable ScopeContext parent)
	)
}

func NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(scope *Scope, declaredVariableInfos j2gort.Ext_com_google_common_collect_ImmutableList, parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(this_, scope, declaredVariableInfos, parent)
	return this_
}
func allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(slab_ *j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext], scope *Scope, declaredVariableInfos j2gort.Ext_com_google_common_collect_ImmutableList, parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext {
	this_ := slab_.Alloc()
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(this_, scope, declaredVariableInfos, parent)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext, scope *Scope, declaredVariableInfos j2gort.Ext_com_google_common_collect_ImmutableList, parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) {
	_ = scope
	_ = declaredVariableInfos
	_ = parent
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext()
	this_.scope = scope
	this_.declaredVariableInfos = j2gort.ImmListToSliceNullable[*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo](declaredVariableInfos)
	this_.parent = parent
}

// Java method: ScopeContext.addReferencedVariableInfo(com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo variableInfo)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) AddReferencedVariableInfo(variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	_ = variableInfo
	this_.referencedVariables.Add(variableInfo)
}

// Java method: ScopeContext.recordPotentialConflicts(com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo referencedVariable)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) RecordPotentialConflicts(referencedVariable MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	_ = referencedVariable
	var thisScopeDepth int32 = this_.scope.GetDepth()
	_ = thisScopeDepth
	var referencedVariableDeclarationDepth int32 = referencedVariable.GetScopeDepth()
	_ = referencedVariableDeclarationDepth
	if thisScopeDepth >= referencedVariableDeclarationDepth {
		{

			for _, declaredVariableInfo := range this_.declaredVariableInfos {

				_ = declaredVariableInfo
				referencedVariable.AddPotentialConflict(ToIface_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo(declaredVariableInfo))
				declaredVariableInfo.AddPotentialConflict(referencedVariable)
			}
		}
		if this_.parent != nil {
			this_.parent.RecordPotentialConflicts(referencedVariable)
		}
	}
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap struct {
	outer_      *MakeDeclaredNamesUnique_ContextualRenameInverter
	mapOfStacks j2gort.OrderedMap[string, j2gort.Ext_java_util_Deque] `gson:"mapOfStacks"`
	idHash_     atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) IsJavaNil_() bool {
	return this_ == nil
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$VariableInfoStackMap", "VariableInfoStackMap"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$VariableInfoStackMap", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap](v_)
	})
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) SetOuter_(o *MakeDeclaredNamesUnique_ContextualRenameInverter) {
	this_.outer_ = o
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) GetOuter_() *MakeDeclaredNamesUnique_ContextualRenameInverter {
	return this_.outer_
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap() {
	_ = this_
	this_.mapOfStacks = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_Deque](

	// Java constructor: VariableInfoStackMap()
	)
}

func NewMakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap(outer_ *MakeDeclaredNamesUnique_ContextualRenameInverter) *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap{}
	this_.outer_ = outer_
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap(this_)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) {
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap()
}

// Java method: VariableInfoStackMap.pushVariableInfo(com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo info)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) PushVariableInfo(info MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	_ = info
	var name string = info.GetName()
	_ = name
	var variableInfoStack j2gort.Ext_java_util_Deque = this_.priv_VariableInfoStackMap_getVariableInfoStack_3829(name)
	_ = variableInfoStack
	variableInfoStack.Push(info)
}

// Java method: VariableInfoStackMap.getVariableInfoStack(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) priv_VariableInfoStackMap_getVariableInfoStack_3829(name string) j2gort.Ext_java_util_Deque {
	_ = name
	var stack j2gort.Ext_java_util_Deque = this_.mapOfStacks.Get(name)
	_ = stack
	if stack == nil {
		stack = j2gort.NewExt_java_util_ArrayDeque()
		this_.mapOfStacks.Put(name, stack)
	}
	return stack
}

// Java method: VariableInfoStackMap.getOrCreateCurrentVariableInfo(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) GetOrCreateCurrentVariableInfo(name string) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	_ = name
	var variableInfoStack j2gort.Ext_java_util_Deque = this_.priv_VariableInfoStackMap_getVariableInfoStack_3829(name)
	_ = variableInfoStack
	var variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo = nil
	_ = variableInfo
	if variableInfoStack.IsEmpty() {
		variableInfo = this_.outer_.priv_ContextualRenameInverter_createUndeclaredVariableInfo_3818(name)
		variableInfoStack.Push(variableInfo)
	} else {
		variableInfo = j2gort.AsIface[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo](variableInfoStack.Peek())
	}
	return variableInfo
}

// Java method: VariableInfoStackMap.popVariableInfo(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) PopVariableInfo(name string) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	_ = name
	var variableInfos j2gort.Ext_java_util_Deque = this_.mapOfStacks.Get(name)
	_ = variableInfos
	j2gort.CheckNotNull(variableInfos)
	var info MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo = j2gort.AsIface[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo](variableInfos.Pop())
	_ = info
	if variableInfos.IsEmpty() {
		this_.mapOfStacks.Remove(name)
	}
	return info
}

// Java method: VariableInfoStackMap.clear()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) Clear() {
	this_.mapOfStacks.Clear()
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo
type MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo struct {
	name    string `gson:"name"`
	idHash_ atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) IsJavaNil_() bool {
	return this_ == nil
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$UndeclaredVariableInfo", "UndeclaredVariableInfo"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$UndeclaredVariableInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo](v_)
	})
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo() {
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: UndeclaredVariableInfo(java.lang.String name)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(name string) *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(this_, name)
	return this_
}
func allocPooled_MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(slab_ *j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo], name string) *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo {
	this_ := slab_.Alloc()
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(this_, name)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo, name string) {
	_ = name
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo()
	this_.name = name
}

// Java method: UndeclaredVariableInfo.addReferenceNode(com.google.javascript.rhino.Node nameNode)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) AddReferenceNode(nameNode pkg_rhino.Node) {
	_ = nameNode
}

// Java method: UndeclaredVariableInfo.getName()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) GetName() string {
	return this_.name
}

// Java method: UndeclaredVariableInfo.getScopeDepth()
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) GetScopeDepth() int32 {
	return int32(0)
}

// Java method: UndeclaredVariableInfo.addPotentialConflict(com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo variableInfo)
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) AddPotentialConflict(variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	_ = variableInfo
}

// Java interface: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo interface {
	Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo()
	IsJavaNil_() bool
	AddReferenceNode(pkg_rhino.Node)
	GetName() string
	GetScopeDepth() int32
	AddPotentialConflict(MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo)
}

func ToIface_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo[T interface {
	MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo
	comparable
}](v T) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo[T MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo](v T) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	var vi MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenameInverter$VariableInfo", func(v_ any) bool {
		return j2gort.InstanceOf[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenamer
type MakeDeclaredNamesUnique_ContextualRenamer interface {
	Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenamer()
	IsJavaNil_() bool
	GetMakeDeclaredNamesUnique_ContextualRenamerBase() *MakeDeclaredNamesUnique_ContextualRenamerBase
	MakeDeclaredNamesUnique_Renamer
	String() string
}

func ToIface_MakeDeclaredNamesUnique_ContextualRenamer[T interface {
	MakeDeclaredNamesUnique_ContextualRenamer
	comparable
}](v T) MakeDeclaredNamesUnique_ContextualRenamer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_MakeDeclaredNamesUnique_ContextualRenamer[T MakeDeclaredNamesUnique_ContextualRenamer](v T) MakeDeclaredNamesUnique_ContextualRenamer {
	var vi MakeDeclaredNamesUnique_ContextualRenamer = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenamer
type MakeDeclaredNamesUnique_ContextualRenamerBase struct {
	self_        MakeDeclaredNamesUnique_ContextualRenamer
	nameUsage    *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage `gson:"nameUsage"`
	declName0    string                                               `gson:"declName0"`
	declRep0     string                                               `gson:"declRep0"`
	declName1    string                                               `gson:"declName1"`
	declRep1     string                                               `gson:"declRep1"`
	declName2    string                                               `gson:"declName2"`
	declRep2     string                                               `gson:"declRep2"`
	declName3    string                                               `gson:"declName3"`
	declRep3     string                                               `gson:"declRep3"`
	declarations *j2gort.OrderedMap[string, string]                   `gson:"declarations"`
	hoistRenamer MakeDeclaredNamesUnique_Renamer                      `gson:"hoistRenamer"`
	scopeRoot    pkg_rhino.Node                                       `gson:"scopeRoot"`
	declCount    int32                                                `gson:"declCount"`
	idHash_      atomic.Int32
	global       bool `gson:"global"`
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(MakeDeclaredNamesUnique_ContextualRenamer)
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique_ContextualRenamerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenamer", "ContextualRenamer"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenamer", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenamer", func(v_ any) bool {
		return j2gort.InstanceOf[MakeDeclaredNamesUnique_ContextualRenamer](v_)
	})
}
func (b *MakeDeclaredNamesUnique_ContextualRenamerBase) GetMakeDeclaredNamesUnique_ContextualRenamerBase() *MakeDeclaredNamesUnique_ContextualRenamerBase {
	return b
}
func (*MakeDeclaredNamesUnique_ContextualRenamerBase) Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenamer() {
}
func (*MakeDeclaredNamesUnique_ContextualRenamerBase) Implements_jscomp_MakeDeclaredNamesUnique_Renamer() {
}

func StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR() string {
	return "$jscomp$"
}
func StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_SUFFIX_CACHE_SIZE() int32 {
	return int32(16384)
}

var static_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES []string

func StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES() []string {
	clinit_MakeDeclaredNamesUnique_ContextualRenamer()
	return static_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES
}
func StaticSet_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES(v []string) {
	clinit_MakeDeclaredNamesUnique_ContextualRenamer()
	static_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES = v
}

var clinitDone_MakeDeclaredNamesUnique_ContextualRenamer uint32
var clinitState_MakeDeclaredNamesUnique_ContextualRenamer j2gort.ClinitState

func clinit_MakeDeclaredNamesUnique_ContextualRenamer() {
	if atomic.LoadUint32(&clinitDone_MakeDeclaredNamesUnique_ContextualRenamer) != 0 {
		return
	}
	clinitSlow_MakeDeclaredNamesUnique_ContextualRenamer()
}

func clinitSlow_MakeDeclaredNamesUnique_ContextualRenamer() {
	if !j2gort.EnterClinit(&clinitState_MakeDeclaredNamesUnique_ContextualRenamer, &clinitDone_MakeDeclaredNamesUnique_ContextualRenamer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_MakeDeclaredNamesUnique_ContextualRenamer, &clinitDone_MakeDeclaredNamesUnique_ContextualRenamer)
	static_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES = j2gort.NewStringArray(int32(16384))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(16384); i = i + 1 {
			{
				a_21 := StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES()
				i_22 := i
				a_21[i_22] = ("$jscomp$" + j2gort.StrValInt32(i))
			}
		}
	}
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) initFields_MakeDeclaredNamesUnique_ContextualRenamer() {
	_ = this_
	this_.declName0 = j2gort.NullString
	this_.declRep0 = j2gort.NullString
	this_.declName1 = j2gort.NullString
	this_.declRep1 = j2gort.NullString
	this_.declName2 = j2gort.NullString
	this_.declRep2 = j2gort.NullString
	this_.declName3 = j2gort.NullString
	this_.declRep3 = j2gort.NullString

}

// Java constructor: ContextualRenamer()
func NewMakeDeclaredNamesUnique_ContextualRenamer() MakeDeclaredNamesUnique_ContextualRenamer {
	clinit_MakeDeclaredNamesUnique_ContextualRenamer()
	this_ := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	this_.BindSelf_(this_)
	InitMakeDeclaredNamesUnique_ContextualRenamer(this_)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenamer(this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) {
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenamer()
	this_.scopeRoot = pkg_rhino.Node(0)
	this_.global = true
	this_.nameUsage = NewMakeDeclaredNamesUnique_ContextualRenamer_NameUsage()
	this_.declCount = int32(0)
	this_.declName0 = j2gort.NullString
	this_.declRep0 = j2gort.NullString
	this_.declName1 = j2gort.NullString
	this_.declRep1 = j2gort.NullString
	this_.declName2 = j2gort.NullString
	this_.declRep2 = j2gort.NullString
	this_.declName3 = j2gort.NullString
	this_.declRep3 = j2gort.NullString
	this_.hoistRenamer = Upcast_MakeDeclaredNamesUnique_Renamer(this_.self_)
}

// Java constructor: ContextualRenamer(com.google.javascript.rhino.Node scopeRoot, com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenamer.NameUsage nameUsage, boolean hoistingTargetScope, com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer parent)
func NewMakeDeclaredNamesUnique_ContextualRenamer_4(scopeRoot pkg_rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) MakeDeclaredNamesUnique_ContextualRenamer {
	clinit_MakeDeclaredNamesUnique_ContextualRenamer()
	this_ := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	this_.BindSelf_(this_)
	InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_, scopeRoot, nameUsage, hoistingTargetScope, parent)
	return this_
}
func allocPooled_MakeDeclaredNamesUnique_ContextualRenamer_4(slab_ *j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenamerBase], scopeRoot pkg_rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) MakeDeclaredNamesUnique_ContextualRenamer {
	clinit_MakeDeclaredNamesUnique_ContextualRenamer()
	this_ := slab_.Alloc()

	this_.BindSelf_(this_)
	InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_, scopeRoot, nameUsage, hoistingTargetScope, parent)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_ *MakeDeclaredNamesUnique_ContextualRenamerBase, scopeRoot pkg_rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) {
	_ = scopeRoot
	_ = nameUsage
	_ = hoistingTargetScope
	_ = parent
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenamer()
	if !(NodeUtil_CreatesScope(scopeRoot)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(scopeRoot))
	}
	if scopeRoot.IsFunction() {
		if !(!(hoistingTargetScope)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(scopeRoot))
		}
	}
	this_.scopeRoot = scopeRoot
	this_.global = false
	this_.nameUsage = nameUsage
	this_.declCount = int32(0)
	this_.declName0 = j2gort.NullString
	this_.declRep0 = j2gort.NullString
	this_.declName1 = j2gort.NullString
	this_.declRep1 = j2gort.NullString
	this_.declName2 = j2gort.NullString
	this_.declRep2 = j2gort.NullString
	this_.declName3 = j2gort.NullString
	this_.declRep3 = j2gort.NullString
	if hoistingTargetScope {
		if !(!(NodeUtil_CreatesBlockScope(scopeRoot))) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(scopeRoot))
		}
		this_.hoistRenamer = Upcast_MakeDeclaredNamesUnique_Renamer(this_.self_)
	} else {
		if !(NodeUtil_CreatesBlockScope(scopeRoot) || scopeRoot.IsFunction()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(scopeRoot))
		}
		this_.hoistRenamer = parent.GetHoistRenamer()
	}
}

// Java method: ContextualRenamer.toString()
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) String() string {
	if this_.declarations != nil {
		return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_.self_).Add_2_string_any("scopeRoot", j2gort.BoxNode(this_.scopeRoot)).Add_2_string_any("nameUsage", j2gort.PtrToAny(this_.nameUsage)).Add_2_string_any("declarations", j2gort.PtrToAny(this_.declarations)).Add_2_string_bool("global", this_.global).String()
	}
	var declView j2gort.OrderedMap[string, string] = j2gort.NewOrderedMapCap[string, string](this_.declCount)
	_ = declView
	if this_.declCount > int32(0) {
		declView.Put(this_.declName0, this_.declRep0)
		if this_.declCount > int32(1) {
			declView.Put(this_.declName1, this_.declRep1)
			if this_.declCount > int32(2) {
				declView.Put(this_.declName2, this_.declRep2)
				if this_.declCount > int32(3) {
					declView.Put(this_.declName3, this_.declRep3)
				}
			}
		}
	}
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_.self_).Add_2_string_any("scopeRoot", j2gort.BoxNode(this_.scopeRoot)).Add_2_string_any("nameUsage", j2gort.PtrToAny(this_.nameUsage)).Add_2_string_any("declarations", &declView).Add_2_string_bool("global", this_.global).String()
}

// Java method: ContextualRenamer.createForChildScope(com.google.javascript.rhino.Node scopeRoot, boolean hoistingTargetScope)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) CreateForChildScope(scopeRoot pkg_rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	_ = scopeRoot
	_ = hoistingTargetScope
	return Upcast_MakeDeclaredNamesUnique_Renamer(NewMakeDeclaredNamesUnique_ContextualRenamer_4(scopeRoot, this_.nameUsage, hoistingTargetScope, Upcast_MakeDeclaredNamesUnique_Renamer(this_.self_)))
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) createForChildScopeWithSlab_(slab_ *j2gort.PooledSlab[MakeDeclaredNamesUnique_ContextualRenamerBase], scopeRoot pkg_rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	_ = scopeRoot
	_ = hoistingTargetScope
	return Upcast_MakeDeclaredNamesUnique_Renamer(allocPooled_MakeDeclaredNamesUnique_ContextualRenamer_4(slab_, scopeRoot, this_.nameUsage, hoistingTargetScope, Upcast_MakeDeclaredNamesUnique_Renamer(this_.self_)))
}

// Java method: ContextualRenamer.hasLocalDeclaration(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) priv_ContextualRenamer_hasLocalDeclaration_3834(name string) bool {
	_ = name
	var c int32 = this_.declCount
	_ = c
	if c == int32(0) {
		return false
	}
	if c <= int32(4) {
		{
			var iifeCond_1 bool
			{

				sr_23 := name
				se_24 := this_.declName0
				iifeCond_1 = j2gort.CheckNotNullString(sr_23) == se_24
			}
			if iifeCond_1 {
				return true
			}
		}
		if (c > int32(1)) && func() bool { sr_25 := name; se_26 := this_.declName1; return j2gort.CheckNotNullString(sr_25) == se_26 }() {
			return true
		}
		if (c > int32(2)) && func() bool { sr_27 := name; se_28 := this_.declName2; return j2gort.CheckNotNullString(sr_27) == se_28 }() {
			return true
		}
		return ((c > int32(3)) && func() bool { sr_29 := name; se_30 := this_.declName3; return j2gort.CheckNotNullString(sr_29) == se_30 }())
	}
	return this_.declarations.ContainsKey(name)
}

// Java method: ContextualRenamer.putLocalDeclaration(java.lang.String name, java.lang.@org.jspecify.annotations.Nullable String newName)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) priv_ContextualRenamer_putLocalDeclaration_3835(name string, newName string) {
	_ = name
	_ = newName
	var c int32 = this_.declCount
	_ = c
	switch c {
	case int32(0):
		this_.declName0 = name
		this_.declRep0 = newName
		this_.declCount = int32(1)
	case int32(1):
		this_.declName1 = name
		this_.declRep1 = newName
		this_.declCount = int32(2)
	case int32(2):
		this_.declName2 = name
		this_.declRep2 = newName
		this_.declCount = int32(3)
	case int32(3):
		this_.declName3 = name
		this_.declRep3 = newName
		this_.declCount = int32(4)
	case int32(4):
		if this_.declarations == nil {
			this_.declarations = j2gort.AcquireOrderedMapPtrCap[string, string](8)
		}

		var n0 string = this_.declName0
		_ = n0
		var r0 string = this_.declRep0
		_ = r0
		var n1 string = this_.declName1
		_ = n1
		var r1 string = this_.declRep1
		_ = r1
		var n2 string = this_.declName2
		_ = n2
		var r2 string = this_.declRep2
		_ = r2
		var n3 string = this_.declName3
		_ = n3
		var r3 string = this_.declRep3
		_ = r3
		this_.declName0 = j2gort.NullString
		this_.declRep0 = j2gort.NullString
		this_.declName1 = j2gort.NullString
		this_.declRep1 = j2gort.NullString
		this_.declName2 = j2gort.NullString
		this_.declRep2 = j2gort.NullString
		this_.declName3 = j2gort.NullString
		this_.declRep3 = j2gort.NullString
		this_.declarations = j2gort.NewOrderedMapPtrCap[string, string](int32(8))
		this_.declarations.Put(n0, r0)
		this_.declarations.Put(n1, r1)
		this_.declarations.Put(n2, r2)
		this_.declarations.Put(n3, r3)
		this_.declarations.Put(name, newName)
		this_.declCount = int32(5)
	default:
		this_.declarations.Put(name, newName)
	}
}

// Java method: ContextualRenamer.addDeclaredName(java.lang.String name, boolean hoisted)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) AddDeclaredName(name string, hoisted bool) {
	_ = name
	_ = hoisted
	if hoisted && !(this_.hoistRenamer == this_.self_ || ((this_.hoistRenamer == nil || this_.hoistRenamer.IsJavaNil_()) && (this_.self_ == nil || this_.self_.IsJavaNil_()))) {
		this_.hoistRenamer.AddDeclaredName(name, true)
	} else {
		{
			var iifeCond_1 bool
			{

				sr_31 := name
				se_32 := "arguments"
				iifeCond_1 = j2gort.CheckNotNullString(sr_31) == se_32
			}
			if !(iifeCond_1) {
				if this_.global {
					this_.priv_ContextualRenamer_reserveName_3836(name)
				} else {
					if !(this_.priv_ContextualRenamer_hasLocalDeclaration_3834(name)) {
						var id int32 = this_.priv_ContextualRenamer_incrementNameCount_3837(name)
						_ = id
						var newName string = j2gort.NullString
						_ = newName
						if id != int32(0) {
							newName = this_.nameUsage.GetUniqueName(name, id)
							var nextId int32 = int32(0)
							_ = nextId
							for (func() int32 {
								var asgn_33 int32 = this_.nameUsage.AddOneIfPresent(newName)
								nextId = asgn_33
								return asgn_33
							}()) >= int32(0) {
								newName = this_.nameUsage.GetUniqueName(name, nextId)
							}
							this_.priv_ContextualRenamer_reserveName_3836(newName)
						}
						this_.priv_ContextualRenamer_putLocalDeclaration_3835(name, newName)
					}
				}
			}
		}
	}
}

// Java method: ContextualRenamer.getReplacementName(java.lang.String oldName)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) GetReplacementName(oldName string) string {
	_ = oldName
	var c int32 = this_.declCount
	_ = c
	if c == int32(0) {
		return j2gort.NullString
	}
	if c <= int32(4) {
		{
			var iifeCond_1 bool
			{

				sr_34 := oldName
				se_35 := this_.declName0
				iifeCond_1 = j2gort.CheckNotNullString(sr_34) == se_35
			}
			if iifeCond_1 {
				return this_.declRep0
			}
		}
		if (c > int32(1)) && func() bool {
			sr_36 := oldName
			se_37 := this_.declName1
			return j2gort.CheckNotNullString(sr_36) == se_37
		}() {
			return this_.declRep1
		}
		if (c > int32(2)) && func() bool {
			sr_38 := oldName
			se_39 := this_.declName2
			return j2gort.CheckNotNullString(sr_38) == se_39
		}() {
			return this_.declRep2
		}
		if (c > int32(3)) && func() bool {
			sr_40 := oldName
			se_41 := this_.declName3
			return j2gort.CheckNotNullString(sr_40) == se_41
		}() {
			return this_.declRep3
		}
		return j2gort.NullString
	}
	return j2gort.OrderedMapPtrGetStr(this_.declarations, oldName)
}

// Java static method: ContextualRenamer.formatUniqueName(java.lang.String name, int id)
func MakeDeclaredNamesUnique_ContextualRenamer_FormatUniqueName(name string, id int32) string {
	clinit_MakeDeclaredNamesUnique_ContextualRenamer()
	_ = name
	_ = id
	if (id >= int32(0)) && (id < int32(16384)) {
		return (j2gort.StrValString(name) + j2gort.StrValString(StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SUFFIXES()[id]))
	}
	return (j2gort.StrValString(name) + "$jscomp$" + j2gort.StrValInt32(id))
}

// Java method: ContextualRenamer.reserveName(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) priv_ContextualRenamer_reserveName_3836(name string) {
	_ = name
	this_.nameUsage.Reserve(name)
}

// Java method: ContextualRenamer.incrementNameCount(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) priv_ContextualRenamer_incrementNameCount_3837(name string) int32 {
	_ = name
	return this_.nameUsage.AddOne(name)
}

// Java method: ContextualRenamer.stripConstIfReplaced()
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) StripConstIfReplaced() bool {
	return false
}

// Java method: ContextualRenamer.getHoistRenamer()
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return this_.hoistRenamer
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenamer.NameUsage
type MakeDeclaredNamesUnique_ContextualRenamer_NameUsage struct {
	counts        j2gort.PlainMap[string, int32] `gson:"counts"`
	lcUniqueCache [][]string                     `gson:"lcUniqueCache"`
	memoNames     []string                       `gson:"memoNames"`
	memoIds       []int32                        `gson:"memoIds"`
	memoResults   []string                       `gson:"memoResults"`
	idHash_       atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) IsJavaNil_() bool {
	return this_ == nil
}
func (*MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenamer$NameUsage", "NameUsage"
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.ContextualRenamer.NameUsage", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$ContextualRenamer$NameUsage", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_ContextualRenamer_NameUsage](v_)
	})
}

func StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_NameUsage_LC_CACHE_IDS() int32 {
	return int32(1024)
}
func StaticGet_MakeDeclaredNamesUnique_ContextualRenamer_NameUsage_MEMO_MASK() int32 {
	return int32(2047)
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) initFields_MakeDeclaredNamesUnique_ContextualRenamer_NameUsage() {
	_ = this_
	this_.counts = j2gort.NewPlainMapVal[string, int32](

	// Java constructor: NameUsage()
	)
}

func NewMakeDeclaredNamesUnique_ContextualRenamer_NameUsage() *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage {
	this_ := &MakeDeclaredNamesUnique_ContextualRenamer_NameUsage{}
	InitMakeDeclaredNamesUnique_ContextualRenamer_NameUsage(this_)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenamer_NameUsage(this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) {
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenamer_NameUsage()
}

// Java method: NameUsage.contains(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) Contains(name string) bool {
	_ = name
	return this_.counts.ContainsKey(name)
}

// Java method: NameUsage.reserve(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) Reserve(name string) {
	_ = name
	this_.counts.PutIfAbsent(name, int32(1))
}

// Java method: NameUsage.addOne(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) AddOne(name string) int32 {
	_ = name
	prev, prevOk_ := this_.counts.GetOk(name)
	_ = prevOk_
	_ = prev
	if !prevOk_ {
		this_.counts.PutVoid(name, int32(1))
		return int32(0)
	}
	var old int32
	{
		if !prevOk_ {
			j2gort.ThrowNPE()
		}
		old = prev
	}

	_ = old
	this_.counts.PutVoid(name, (old + int32(1)))
	return old
}

// Java method: NameUsage.addOneIfPresent(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) AddOneIfPresent(name string) int32 {
	_ = name
	prev, prevOk_ := this_.counts.GetOk(name)
	_ = prevOk_
	_ = prev
	if !prevOk_ {
		return int32(-1)
	}
	var old int32
	{
		if !prevOk_ {
			j2gort.ThrowNPE()
		}
		old = prev
	}

	_ = old
	this_.counts.PutVoid(name, (old + int32(1)))
	return old
}

// Java method: NameUsage.getUniqueName(java.lang.String name, int id)
func (this_ *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage) GetUniqueName(name string, id int32) string {
	_ = name
	_ = id
	if ((j2gort.JStrLen(name) == int32(1)) && (id >= int32(0))) && (id < int32(1024)) {
		var c uint16 = j2gort.JStrCharAt(name, int32(0))
		_ = c
		if (int32(c) >= int32(uint16(97))) && (int32(c) <= int32(uint16(122))) {
			var table [][]string = this_.lcUniqueCache
			_ = table
			if table == nil {
				table = make([][]string, int32(26))
				this_.lcUniqueCache = table
			}
			var idx int32 = (int32(c) - int32(uint16(97)))
			_ = idx
			var row []string = table[idx]
			_ = row
			if row == nil {
				row = j2gort.NewStringArray(int32(1024))
				{
					a_42 := table
					i_43 := idx
					a_42[i_43] = row
				}
			}
			var cached string = row[id]
			_ = cached
			if cached == j2gort.NullString {
				cached = MakeDeclaredNamesUnique_ContextualRenamer_FormatUniqueName(name, id)
				{
					a_44 := row
					i_45 := id
					a_44[i_45] = cached
				}
			}
			return cached
		}
	}
	if id >= int32(2) {
		var mNames []string = this_.memoNames
		_ = mNames
		if mNames == nil {
			mNames = j2gort.NewStringArray((int32(2047) + int32(1)))
			this_.memoNames = mNames
			this_.memoIds = make([]int32, (int32(2047) + int32(1)))
			this_.memoResults = j2gort.NewStringArray((int32(2047) + int32(1)))
		}
		var slot int32 = (((javalang.StringHashCode(j2gort.CheckNotNullString(name)) * int32(31)) + id) & int32(2047))
		_ = slot
		if (this_.memoIds[slot] == id) && func() bool { sr_46 := name; se_47 := mNames[slot]; return j2gort.CheckNotNullString(sr_46) == se_47 }() {
			return this_.memoResults[slot]
		}
		var res string = MakeDeclaredNamesUnique_ContextualRenamer_FormatUniqueName(name, id)
		_ = res
		{
			a_48 := mNames
			i_49 := slot
			a_48[i_49] = name
		}
		{
			a_50 := this_.memoIds
			i_51 := slot
			a_50[i_51] = id
		}
		{
			a_52 := this_.memoResults
			i_53 := slot
			a_52[i_53] = res
		}
		return res
	}
	return MakeDeclaredNamesUnique_ContextualRenamer_FormatUniqueName(name, id)
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.InlineRenamer
type MakeDeclaredNamesUnique_InlineRenamer struct {
	declName0          string                                 `gson:"declName0"`
	declRep0           string                                 `gson:"declRep0"`
	declName1          string                                 `gson:"declName1"`
	declRep1           string                                 `gson:"declRep1"`
	declarations       *j2gort.OrderedMap[string, string]     `gson:"declarations"`
	uniqueIdSupplier   j2gort.Ext_java_util_function_Supplier `gson:"uniqueIdSupplier"`
	idPrefix           string                                 `gson:"idPrefix"`
	separatorAndPrefix string                                 `gson:"separatorAndPrefix"`
	convention         CodingConvention                       `gson:"convention"`
	hoistRenamer       MakeDeclaredNamesUnique_Renamer        `gson:"hoistRenamer"`
	declCount          int32                                  `gson:"declCount"`
	idHash_            atomic.Int32
	removeConstness    bool `gson:"removeConstness"`
}

func (this_ *MakeDeclaredNamesUnique_InlineRenamer) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique_InlineRenamer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$InlineRenamer", "InlineRenamer"
}
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_InlineRenamer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.InlineRenamer", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$InlineRenamer", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_InlineRenamer](v_)
	})
}
func (*MakeDeclaredNamesUnique_InlineRenamer) Implements_jscomp_MakeDeclaredNamesUnique_Renamer() {}

func (this_ *MakeDeclaredNamesUnique_InlineRenamer) initFields_MakeDeclaredNamesUnique_InlineRenamer() {
	_ = this_
	this_.declName0 = j2gort.NullString
	this_.declRep0 = j2gort.NullString
	this_.declName1 = j2gort.NullString
	this_.declRep1 = j2gort.NullString

	this_.idPrefix = j2gort.NullString
	this_.separatorAndPrefix = j2gort.NullString
}

// Java constructor: InlineRenamer(com.google.javascript.jscomp.CodingConvention convention, java.util.function.Supplier<java.lang.String> uniqueIdSupplier, java.lang.String idPrefix, boolean removeConstness, boolean hoistingTargetScope, com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer parent)
func NewMakeDeclaredNamesUnique_InlineRenamer(convention CodingConvention, uniqueIdSupplier j2gort.Ext_java_util_function_Supplier, idPrefix string, removeConstness bool, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_InlineRenamer {
	this_ := &MakeDeclaredNamesUnique_InlineRenamer{}
	InitMakeDeclaredNamesUnique_InlineRenamer(this_, convention, uniqueIdSupplier, idPrefix, removeConstness, hoistingTargetScope, parent)
	return this_
}

func InitMakeDeclaredNamesUnique_InlineRenamer(this_ *MakeDeclaredNamesUnique_InlineRenamer, convention CodingConvention, uniqueIdSupplier j2gort.Ext_java_util_function_Supplier, idPrefix string, removeConstness bool, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) {
	_ = convention
	_ = uniqueIdSupplier
	_ = idPrefix
	_ = removeConstness
	_ = hoistingTargetScope
	_ = parent
	this_.initFields_MakeDeclaredNamesUnique_InlineRenamer()
	this_.convention = convention
	this_.uniqueIdSupplier = uniqueIdSupplier
	if !(!(j2gort.CheckNotNullString(idPrefix) == "")) {
		j2gort.CheckArgument(false)
	}
	this_.idPrefix = idPrefix
	this_.separatorAndPrefix = ("$jscomp$" + j2gort.StrValString(idPrefix))
	this_.removeConstness = removeConstness
	if hoistingTargetScope {
		this_.hoistRenamer = ToIface_MakeDeclaredNamesUnique_Renamer(this_)
	} else {
		this_.hoistRenamer = parent.GetHoistRenamer()
	}
}

// Java method: InlineRenamer.addDeclaredName(java.lang.String name, boolean hoisted)
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) AddDeclaredName(name string, hoisted bool) {
	_ = name
	_ = hoisted
	{
		var iifeCond_1 bool
		{

			sr_54 := name
			se_55 := "arguments"
			iifeCond_1 = j2gort.CheckNotNullString(sr_54) == se_55
		}
		if !(!(iifeCond_1)) {
			j2gort.CheckState(false)
		}
	}
	if hoisted && !(this_.hoistRenamer == this_ || ((this_.hoistRenamer == nil || this_.hoistRenamer.IsJavaNil_()) && (this_ == nil || this_.IsJavaNil_()))) {
		this_.hoistRenamer.AddDeclaredName(name, hoisted)
		return
	}
	var c int32 = this_.declCount
	_ = c
	if c == int32(0) {
		this_.declName0 = name
		this_.declRep0 = this_.priv_InlineRenamer_getUniqueName_3843(name)
		this_.declCount = int32(1)
	} else if c == int32(1) {
		{
			var iifeCond_1 bool
			{

				sr_56 := name
				se_57 := this_.declName0
				iifeCond_1 = j2gort.CheckNotNullString(sr_56) == se_57
			}
			if !(iifeCond_1) {
				this_.declName1 = name
				this_.declRep1 = this_.priv_InlineRenamer_getUniqueName_3843(name)
				this_.declCount = int32(2)
			}
		}
	} else if c == int32(2) {
		{
			var iifeCond_1 bool
			{

				sr_58 := name
				se_59 := this_.declName0
				iifeCond_1 = j2gort.CheckNotNullString(sr_58) == se_59
			}
			if !(iifeCond_1) && !(func() bool { sr_60 := name; se_61 := this_.declName1; return j2gort.CheckNotNullString(sr_60) == se_61 }()) {
				var n0 string = this_.declName0
				_ = n0
				var r0 string = this_.declRep0
				_ = r0
				var n1 string = this_.declName1
				_ = n1
				var r1 string = this_.declRep1
				_ = r1
				var newName string = this_.priv_InlineRenamer_getUniqueName_3843(name)
				_ = newName
				this_.declName0 = j2gort.NullString
				this_.declRep0 = j2gort.NullString
				this_.declName1 = j2gort.NullString
				this_.declRep1 = j2gort.NullString
				this_.declarations = j2gort.NewOrderedMapPtrCap[string, string](int32(8))
				this_.declarations.Put(n0, r0)
				this_.declarations.Put(n1, r1)
				this_.declarations.Put(name, newName)
				this_.declCount = int32(3)
			}
		}
	} else if !(this_.declarations.ContainsKey(name)) {
		this_.declarations.Put(name, this_.priv_InlineRenamer_getUniqueName_3843(name))
	}
}

// Java method: InlineRenamer.getUniqueName(java.lang.String name)
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) priv_InlineRenamer_getUniqueName_3843(name string) string {
	_ = name
	if j2gort.CheckNotNullString(name) == "" {
		return name
	}
	var sepIdx int32 = j2gort.JStrLastIndexOf(name, "$jscomp$")
	_ = sepIdx
	if sepIdx != int32(-1) {
		name = j2gort.JStrSub(name, int32(0), sepIdx)
	}
	if this_.convention.IsExported_1(name) {
		name = ("JSCompiler_" + j2gort.StrValString(name))
	}
	return (j2gort.StrValString(name) + j2gort.StrValString(this_.separatorAndPrefix) + j2gort.StrValString(j2gort.CastString(this_.uniqueIdSupplier.Get())))
}

// Java method: InlineRenamer.getReplacementName(java.lang.String oldName)
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) GetReplacementName(oldName string) string {
	_ = oldName
	var c int32 = this_.declCount
	_ = c
	if c == int32(0) {
		return j2gort.NullString
	}
	if c <= int32(2) {
		{
			var iifeCond_1 bool
			{

				sr_62 := oldName
				se_63 := this_.declName0
				iifeCond_1 = j2gort.CheckNotNullString(sr_62) == se_63
			}
			if iifeCond_1 {
				return this_.declRep0
			}
		}
		if (c > int32(1)) && func() bool {
			sr_64 := oldName
			se_65 := this_.declName1
			return j2gort.CheckNotNullString(sr_64) == se_65
		}() {
			return this_.declRep1
		} else {
			return j2gort.NullString
		}
	}
	return j2gort.OrderedMapPtrGetStr(this_.declarations, oldName)
}

// Java method: InlineRenamer.createForChildScope(com.google.javascript.rhino.Node scopeRoot, boolean hoistingTargetScope)
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) CreateForChildScope(scopeRoot pkg_rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	_ = scopeRoot
	_ = hoistingTargetScope
	return NewMakeDeclaredNamesUnique_InlineRenamer(this_.convention, this_.uniqueIdSupplier, this_.idPrefix, this_.removeConstness, hoistingTargetScope, this_)
}

// Java method: InlineRenamer.stripConstIfReplaced()
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) StripConstIfReplaced() bool {
	return this_.removeConstness
}

// Java method: InlineRenamer.getHoistRenamer()
func (this_ *MakeDeclaredNamesUnique_InlineRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return this_.hoistRenamer
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.BoilerplateRenamer
type MakeDeclaredNamesUnique_BoilerplateRenamer struct {
	MakeDeclaredNamesUnique_ContextualRenamerBase
	uniqueIdSupplier j2gort.Ext_java_util_function_Supplier `gson:"uniqueIdSupplier"`
	idPrefix         string                                 `gson:"idPrefix"`
	convention       CodingConvention                       `gson:"convention"`
}

func (this_ *MakeDeclaredNamesUnique_BoilerplateRenamer) BindSelf_(s any) {
	_ = s
	this_.MakeDeclaredNamesUnique_ContextualRenamerBase.BindSelf_(s)
}
func (this_ *MakeDeclaredNamesUnique_BoilerplateRenamer) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique_BoilerplateRenamer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$BoilerplateRenamer", "BoilerplateRenamer"
}
func (this_ *MakeDeclaredNamesUnique_BoilerplateRenamer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.MakeDeclaredNamesUnique_ContextualRenamerBase.IdentityHash_()
}
func (*MakeDeclaredNamesUnique_BoilerplateRenamer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.BoilerplateRenamer", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$BoilerplateRenamer", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_BoilerplateRenamer](v_)
	})
}
func (*MakeDeclaredNamesUnique_BoilerplateRenamer) Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenamer() {
}
func (*MakeDeclaredNamesUnique_BoilerplateRenamer) Implements_jscomp_MakeDeclaredNamesUnique_Renamer() {
}

func (this_ *MakeDeclaredNamesUnique_BoilerplateRenamer) initFields_MakeDeclaredNamesUnique_BoilerplateRenamer() {
	_ = this_
	this_.idPrefix = j2gort.NullString
}

// Java constructor: BoilerplateRenamer(com.google.javascript.jscomp.CodingConvention convention, java.util.function.Supplier<java.lang.String> uniqueIdSupplier, java.lang.String idPrefix)
func NewMakeDeclaredNamesUnique_BoilerplateRenamer(convention CodingConvention, uniqueIdSupplier j2gort.Ext_java_util_function_Supplier, idPrefix string) *MakeDeclaredNamesUnique_BoilerplateRenamer {
	this_ := &MakeDeclaredNamesUnique_BoilerplateRenamer{}
	this_.BindSelf_(this_)
	InitMakeDeclaredNamesUnique_BoilerplateRenamer(this_, convention, uniqueIdSupplier, idPrefix)
	return this_
}

func InitMakeDeclaredNamesUnique_BoilerplateRenamer(this_ *MakeDeclaredNamesUnique_BoilerplateRenamer, convention CodingConvention, uniqueIdSupplier j2gort.Ext_java_util_function_Supplier, idPrefix string) {
	_ = convention
	_ = uniqueIdSupplier
	_ = idPrefix
	InitMakeDeclaredNamesUnique_ContextualRenamer(&this_.MakeDeclaredNamesUnique_ContextualRenamerBase)
	this_.initFields_MakeDeclaredNamesUnique_BoilerplateRenamer()
	this_.convention = convention
	this_.uniqueIdSupplier = uniqueIdSupplier
	this_.idPrefix = idPrefix
}

// Java method: BoilerplateRenamer.createForChildScope(com.google.javascript.rhino.Node scopeRoot, boolean hoisted)
func (this_ *MakeDeclaredNamesUnique_BoilerplateRenamer) CreateForChildScope(scopeRoot pkg_rhino.Node, hoisted bool) MakeDeclaredNamesUnique_Renamer {
	_ = scopeRoot
	_ = hoisted
	return NewMakeDeclaredNamesUnique_InlineRenamer(this_.convention, this_.uniqueIdSupplier, this_.idPrefix, false, hoisted, this_)
}

// Java class: com.google.javascript.jscomp.MakeDeclaredNamesUnique.TargettedRenamer
type MakeDeclaredNamesUnique_TargettedRenamer struct {
	delegate MakeDeclaredNamesUnique_Renamer `gson:"delegate"`
	targets  j2gort.Ext_java_util_Set        `gson:"targets"`
	idHash_  atomic.Int32
}

func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) BindSelf_(s any) {
	_ = s
}
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) IsJavaNil_() bool { return this_ == nil }
func (*MakeDeclaredNamesUnique_TargettedRenamer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MakeDeclaredNamesUnique$TargettedRenamer", "TargettedRenamer"
}
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MakeDeclaredNamesUnique_TargettedRenamer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MakeDeclaredNamesUnique.TargettedRenamer", "com.google.javascript.jscomp.MakeDeclaredNamesUnique$TargettedRenamer", func(v_ any) bool {
		return j2gort.InstanceOf[*MakeDeclaredNamesUnique_TargettedRenamer](v_)
	})
}
func (*MakeDeclaredNamesUnique_TargettedRenamer) Implements_jscomp_MakeDeclaredNamesUnique_Renamer() {
}

func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) initFields_MakeDeclaredNamesUnique_TargettedRenamer() {
	_ = this_
}

// Java constructor: TargettedRenamer(com.google.javascript.jscomp.MakeDeclaredNamesUnique.Renamer delegate, java.util.Set<java.lang.String> targets)
func NewMakeDeclaredNamesUnique_TargettedRenamer(delegate MakeDeclaredNamesUnique_Renamer, targets j2gort.Ext_java_util_Set) *MakeDeclaredNamesUnique_TargettedRenamer {
	this_ := &MakeDeclaredNamesUnique_TargettedRenamer{}
	InitMakeDeclaredNamesUnique_TargettedRenamer(this_, delegate, targets)
	return this_
}

func InitMakeDeclaredNamesUnique_TargettedRenamer(this_ *MakeDeclaredNamesUnique_TargettedRenamer, delegate MakeDeclaredNamesUnique_Renamer, targets j2gort.Ext_java_util_Set) {
	_ = delegate
	_ = targets
	this_.initFields_MakeDeclaredNamesUnique_TargettedRenamer()
	this_.delegate = delegate
	this_.targets = targets
}

// Java method: TargettedRenamer.addDeclaredName(java.lang.String name, boolean hoisted)
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) AddDeclaredName(name string, hoisted bool) {
	_ = name
	_ = hoisted
	if j2gort.SetContainsStr(this_.targets, name) {
		this_.delegate.AddDeclaredName(name, hoisted)
	}
}

// Java method: TargettedRenamer.getReplacementName(java.lang.String oldName)
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) GetReplacementName(oldName string) string {
	_ = oldName
	if j2gort.SetContainsStr(this_.targets, oldName) {
		return this_.delegate.GetReplacementName(oldName)
	} else {
		return j2gort.NullString
	}
}

// Java method: TargettedRenamer.stripConstIfReplaced()
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) StripConstIfReplaced() bool {
	return this_.delegate.StripConstIfReplaced()
}

// Java method: TargettedRenamer.createForChildScope(com.google.javascript.rhino.Node scopeRoot, boolean hoistingTargetScope)
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) CreateForChildScope(scopeRoot pkg_rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	_ = scopeRoot
	_ = hoistingTargetScope
	return NewMakeDeclaredNamesUnique_TargettedRenamer(this_.delegate.CreateForChildScope(scopeRoot, hoistingTargetScope), this_.targets)
}

// Java method: TargettedRenamer.getHoistRenamer()
func (this_ *MakeDeclaredNamesUnique_TargettedRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return this_.delegate.GetHoistRenamer()
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_ScopeContext.Release()
	this_.slab_DeclaredVariableInfo.Release()
	this_.slab_RenamingInfo.Release()
	this_.slab_UndeclaredVariableInfo.Release()
}
func (*MakeDeclaredNamesUnique) AlwaysShouldTraverse_()                          {}
func (*MakeDeclaredNamesUnique) CanReleaseSyntacticScopes_()                     {}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) AlwaysShouldTraverse_() {}
func (this_ *MakeDeclaredNamesUnique) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_ContextualRenamer.ReleaseWith(func(cr_ *MakeDeclaredNamesUnique_ContextualRenamerBase) {
		if cr_.declarations != nil {
			j2gort.ReleaseOrderedMapPtr(cr_.declarations)
			cr_.declarations = nil
		}
	})
}
