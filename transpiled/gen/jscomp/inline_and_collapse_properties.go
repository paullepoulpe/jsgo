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
// Port of com/google/javascript/jscomp/InlineAndCollapseProperties.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties
type InlineAndCollapseProperties struct {
	compiler                      AbstractCompiler               `gson:"compiler"`
	optionalGlobalNamespaceTester *j2gort.Ext_java_util_Optional `gson:"optionalGlobalNamespaceTester"`
	decisionsLog                  pkg_diagnostic.LogFile         `gson:"decisionsLog"`
	namespace                     *GlobalNamespace               `gson:"namespace"`
	idHash_                       atomic.Int32
	propertyCollapseLevel         CompilerOptions_PropertyCollapseLevel `gson:"propertyCollapseLevel"`
	chunkOutputType               CompilerOptions_ChunkOutputType       `gson:"chunkOutputType"`
	haveModulesBeenRewritten      bool                                  `gson:"haveModulesBeenRewritten"`
	moduleResolutionMode          ModuleLoader_ResolutionMode           `gson:"moduleResolutionMode"`
	testAggressiveInliningOnly    bool                                  `gson:"testAggressiveInliningOnly"`
}

func (this_ *InlineAndCollapseProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties) IsJavaNil_() bool { return this_ == nil }
func (*InlineAndCollapseProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties", "InlineAndCollapseProperties"
}
func (this_ *InlineAndCollapseProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineAndCollapseProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties", "com.google.javascript.jscomp.InlineAndCollapseProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties](v_)
	})
}
func (*InlineAndCollapseProperties) Implements_jscomp_CompilerPass() {}

var static_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING *DiagnosticType

func StaticGet_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING() *DiagnosticType {
	clinit_InlineAndCollapseProperties()
	return static_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING
}
func StaticSet_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING(v *DiagnosticType) {
	clinit_InlineAndCollapseProperties()
	static_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING = v
}

var static_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING *DiagnosticType

func StaticGet_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING() *DiagnosticType {
	clinit_InlineAndCollapseProperties()
	return static_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING
}
func StaticSet_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING(v *DiagnosticType) {
	clinit_InlineAndCollapseProperties()
	static_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING = v
}

var static_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE *DiagnosticType

func StaticGet_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE() *DiagnosticType {
	clinit_InlineAndCollapseProperties()
	return static_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE
}
func StaticSet_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE(v *DiagnosticType) {
	clinit_InlineAndCollapseProperties()
	static_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE = v
}

var static_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING *DiagnosticType

func StaticGet_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING() *DiagnosticType {
	clinit_InlineAndCollapseProperties()
	return static_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING
}
func StaticSet_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING(v *DiagnosticType) {
	clinit_InlineAndCollapseProperties()
	static_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING = v
}

var static_InlineAndCollapseProperties_ALIAS_CYCLE *DiagnosticType

func StaticGet_InlineAndCollapseProperties_ALIAS_CYCLE() *DiagnosticType {
	clinit_InlineAndCollapseProperties()
	return static_InlineAndCollapseProperties_ALIAS_CYCLE
}
func StaticSet_InlineAndCollapseProperties_ALIAS_CYCLE(v *DiagnosticType) {
	clinit_InlineAndCollapseProperties()
	static_InlineAndCollapseProperties_ALIAS_CYCLE = v
}

var clinitDone_InlineAndCollapseProperties uint32
var clinitState_InlineAndCollapseProperties j2gort.ClinitState

func clinit_InlineAndCollapseProperties() {
	if atomic.LoadUint32(&clinitDone_InlineAndCollapseProperties) != 0 {
		return
	}
	clinitSlow_InlineAndCollapseProperties()
}

func clinitSlow_InlineAndCollapseProperties() {
	if !j2gort.EnterClinit(&clinitState_InlineAndCollapseProperties, &clinitDone_InlineAndCollapseProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InlineAndCollapseProperties, &clinitDone_InlineAndCollapseProperties)
	static_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING = DiagnosticType_Warning("JSC_PARTIAL_NAMESPACE", "Partial alias created for namespace {0}, possibly due to await/yield transpilation.\nThis may prevent optimization of anything nested under this namespace.\nSee https://github.com/google/closure-compiler/wiki/FAQ#i-got-a-partial-alias-created-for-namespace-error--what-do-i-do for more details.")
	static_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING = DiagnosticType_Warning("JSC_NAMESPACE_REDEFINED", "namespace {0} should not be redefined")
	static_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE = DiagnosticType_Warning("JSC_RECEIVER_AFFECTED_BY_COLLAPSE", "Receiver reference in function {0} changes meaning when namespace is collapsed.\n Consider annotating @nocollapse; however, other properties on the receiver may still be collapsed.")
	static_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING = DiagnosticType_Warning("JSC_UNSAFE_CTOR_ALIASING", "Variable {0} aliases a constructor, so it cannot be assigned multiple times")
	static_InlineAndCollapseProperties_ALIAS_CYCLE = DiagnosticType_Error("JSC_ALIAS_CYCLE", "Alias path contains a cycle: {0} to {1}")
}

func (this_ *InlineAndCollapseProperties) initFields_InlineAndCollapseProperties() {
	_ = this_

}

// Java constructor: InlineAndCollapseProperties(com.google.javascript.jscomp.InlineAndCollapseProperties.Builder builder)
func NewInlineAndCollapseProperties(builder *InlineAndCollapseProperties_Builder) *InlineAndCollapseProperties {
	clinit_InlineAndCollapseProperties()
	this_ := &InlineAndCollapseProperties{}
	InitInlineAndCollapseProperties(this_, builder)
	return this_
}

func InitInlineAndCollapseProperties(this_ *InlineAndCollapseProperties, builder *InlineAndCollapseProperties_Builder) {
	_ = builder
	this_.initFields_InlineAndCollapseProperties()
	this_.compiler = builder.compiler
	this_.propertyCollapseLevel = builder.propertyCollapseLevel
	this_.chunkOutputType = builder.chunkOutputType
	this_.haveModulesBeenRewritten = builder.haveModulesBeenRewritten
	this_.moduleResolutionMode = builder.moduleResolutionMode
	this_.testAggressiveInliningOnly = builder.testAggressiveInliningOnly
	this_.optionalGlobalNamespaceTester = builder.optionalGlobalNamespaceTester
}

// Java static method: InlineAndCollapseProperties.builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func InlineAndCollapseProperties_S_Builder(compiler AbstractCompiler) *InlineAndCollapseProperties_Builder {
	clinit_InlineAndCollapseProperties()
	_ = compiler
	return NewInlineAndCollapseProperties_Builder(compiler)
}

// Java method: InlineAndCollapseProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer func() {
		if !this_.testAggressiveInliningOnly && this_.namespace != nil {
			this_.namespace.ReleaseSlabs_()
			this_.namespace = nil
		}
	}()

	if !(!(this_.testAggressiveInliningOnly) || (this_.propertyCollapseLevel == StaticGet_CompilerOptions_PropertyCollapseLevel_ALL())) {
		j2gort.CheckStateFmt(false, "testAggressiveInlining is invalid for: %s", j2gort.BoxNullableVal(this_.propertyCollapseLevel))
	}
	func() {
		defer func() {
			this_.decisionsLog = nil
		}()
		var logFile pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "decisions.log", []string{})
		_ = logFile
		defer j2gort.CloseResource(logFile)
		this_.decisionsLog = logFile
		switch (this_.propertyCollapseLevel).Ordinal() {
		case int32(1) /* PropertyCollapseLevel.NONE */ :
			this_.priv_InlineAndCollapseProperties_performMinimalInliningAndNoCollapsing_3208(externs, root)
		case int32(2) /* PropertyCollapseLevel.MODULE_EXPORT */ :
			this_.priv_InlineAndCollapseProperties_performMinimalInliningAndModuleExportCollapsing_3209(externs, root)
		case int32(0) /* PropertyCollapseLevel.ALL */ :
			if this_.testAggressiveInliningOnly {
				this_.priv_InlineAndCollapseProperties_performAggressiveInliningForTest_3211(externs, root)
			} else {
				this_.priv_InlineAndCollapseProperties_performAggressiveInliningAndCollapsing_3210(externs, root)
			}
		}
	}()
}

// Java method: InlineAndCollapseProperties.performMinimalInliningAndNoCollapsing(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties) priv_InlineAndCollapseProperties_performMinimalInliningAndNoCollapsing_3208(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var val_1 InlineAndCollapseProperties_InlineAliases
		val_1.outer_ = this_
		InitInlineAndCollapseProperties_InlineAliases(&val_1)
		(&val_1).Process_2(externs, root)
	}
}

// Java method: InlineAndCollapseProperties.performMinimalInliningAndModuleExportCollapsing(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties) priv_InlineAndCollapseProperties_performMinimalInliningAndModuleExportCollapsing_3209(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var val_2 InlineAndCollapseProperties_InlineAliases
		val_2.outer_ = this_
		InitInlineAndCollapseProperties_InlineAliases(&val_2)
		(&val_2).Process_2(externs, root)
	}
	this_.namespace = NewGlobalNamespace_3_LogFile_AbstractCompiler_Node(this_.decisionsLog, this_.compiler, root)
	{
		var val_3 InlineAndCollapseProperties_CollapseProperties
		val_3.outer_ = this_
		InitInlineAndCollapseProperties_CollapseProperties(&val_3)
		(&val_3).Process_2(externs, root)
	}
}

// Java method: InlineAndCollapseProperties.performAggressiveInliningAndCollapsing(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties) priv_InlineAndCollapseProperties_performAggressiveInliningAndCollapsing_3210(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		clinit_ConcretizeStaticInheritanceForInlining()
		var val_4 ConcretizeStaticInheritanceForInlining
		InitConcretizeStaticInheritanceForInlining(&val_4, this_.compiler)
		(&val_4).Process_2(externs, root)
	}
	{
		var val_5 InlineAndCollapseProperties_AggressiveInlineAliases
		val_5.outer_ = this_
		InitInlineAndCollapseProperties_AggressiveInlineAliases(&val_5)
		(&val_5).Process_2(externs, root)
	}
	{
		var val_6 InlineAndCollapseProperties_CollapseProperties
		val_6.outer_ = this_
		InitInlineAndCollapseProperties_CollapseProperties(&val_6)
		(&val_6).Process_2(externs, root)
	}
}

// Java method: InlineAndCollapseProperties.performAggressiveInliningForTest(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties) priv_InlineAndCollapseProperties_performAggressiveInliningForTest_3211(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var val_7 InlineAndCollapseProperties_AggressiveInlineAliases
	val_7.outer_ = this_
	InitInlineAndCollapseProperties_AggressiveInlineAliases(&val_7)
	var aggressiveInlineAliases *InlineAndCollapseProperties_AggressiveInlineAliases = &val_7
	_ = aggressiveInlineAliases
	aggressiveInlineAliases.Process_2(externs, root)
	j2gort.AsIface[j2gort.Ext_java_util_function_Consumer](this_.optionalGlobalNamespaceTester.Get()).Accept(j2gort.PtrToAny(aggressiveInlineAliases.GetLastUsedGlobalNamespace()))
}

// Java static method: InlineAndCollapseProperties.maybeGetInnerNameNode(com.google.javascript.rhino.Node maybeFunctionOrClassNode)
func InlineAndCollapseProperties_MaybeGetInnerNameNode(maybeFunctionOrClassNode pkg_rhino.Node) pkg_rhino.Node {
	clinit_InlineAndCollapseProperties()
	_ = maybeFunctionOrClassNode
	if NodeUtil_IsFunctionExpression(maybeFunctionOrClassNode) {
		var nameNode pkg_rhino.Node = maybeFunctionOrClassNode.GetFirstChild()
		_ = nameNode
		if !(nameNode.IsName()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(nameNode))
		}
		if j2gort.CheckNotNullString(nameNode.GetString()) == "" {
			return pkg_rhino.Node(0)
		} else {
			return nameNode
		}
	} else if NodeUtil_IsClassExpression(maybeFunctionOrClassNode) {
		var nameNode_2 pkg_rhino.Node = maybeFunctionOrClassNode.GetFirstChild()
		_ = nameNode_2
		if nameNode_2.IsName() {
			return nameNode_2
		} else {
			return pkg_rhino.Node(0)
		}
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java static method: InlineAndCollapseProperties.mayBeGlobalAlias(com.google.javascript.jscomp.GlobalNamespace.Ref alias)
func InlineAndCollapseProperties_MayBeGlobalAlias(alias *GlobalNamespace_Ref) bool {
	clinit_InlineAndCollapseProperties()
	_ = alias
	if alias.scope.IsGlobal() {
		return true
	}
	var aliasParent pkg_rhino.Node = alias.GetNode().GetParent()
	_ = aliasParent
	if !(aliasParent.IsAssign()) && !(aliasParent.IsName()) {
		return true
	}
	var aliasLhsNode pkg_rhino.Node
	_ = aliasLhsNode
	if aliasParent.IsName() {
		aliasLhsNode = aliasParent
	} else {
		aliasLhsNode = aliasParent.GetFirstChild()
	}
	if !(aliasLhsNode.IsName()) {
		return true
	}
	var aliasVarName string = aliasLhsNode.GetString()
	_ = aliasVarName
	var aliasVar *Var = j2gort.CastPtr[Var](alias.scope.GetVar(aliasVarName))
	_ = aliasVar
	if aliasVar != nil {
		return aliasVar.IsGlobal()
	}
	return true
}

// Java static method: InlineAndCollapseProperties.referencesCollapsibleProperty(com.google.javascript.jscomp.ReferenceCollection aliasRefs, com.google.javascript.jscomp.GlobalNamespace.Name aliasedName, com.google.javascript.jscomp.GlobalNamespace namespace)
func InlineAndCollapseProperties_ReferencesCollapsibleProperty(aliasRefs *ReferenceCollection, aliasedName *GlobalNamespace_Name, namespace *GlobalNamespace) bool {
	clinit_InlineAndCollapseProperties()
	_ = aliasRefs
	_ = aliasedName
	_ = namespace
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = aliasRefs.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = aliasRefs.Get(i)
			_ = ref
			if ref.GetParent() == 0 {
				continue
			}
			var parent pkg_rhino.Node = ref.GetParent()
			_ = parent
			if parent.IsAssign() {
				parent = parent.GetParent()
			}
			if NodeUtil_IsNormalOrOptChainGetProp(parent) {
				var propertyName string = parent.GetString()
				_ = propertyName
				var originalPropertyName string = (j2gort.StrValString(aliasedName.GetName()) + "." + j2gort.StrValString(propertyName))
				_ = originalPropertyName
				var originalProperty *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](namespace.GetOwnSlot(originalPropertyName))
				_ = originalProperty
				if (originalProperty == nil) || !(originalProperty.CanCollapse()) {
					continue
				}
				return true
			}
		}
	}
	return false
}

// Java static method: InlineAndCollapseProperties.isUnsafelyReassigned(com.google.javascript.jscomp.GlobalNamespace.Name name)
func InlineAndCollapseProperties_IsUnsafelyReassigned(name *GlobalNamespace_Name) bool {
	clinit_InlineAndCollapseProperties()
	_ = name
	if name.GetTotalSets() <= int32(1) {
		return false
	}
	var foundOriginalDefinition bool = false
	_ = foundOriginalDefinition
	if name.HasInlineRefs() {
		var r0 *GlobalNamespace_Ref = name.InlineRef0()
		_ = r0
		if ((r0 != nil) && r0.IsSet()) && !(InlineAndCollapseProperties_IsSafeNamespaceReinit(r0)) {
			foundOriginalDefinition = true
		}
		var r1 *GlobalNamespace_Ref = name.InlineRef1()
		_ = r1
		if ((r1 != nil) && r1.IsSet()) && !(InlineAndCollapseProperties_IsSafeNamespaceReinit(r1)) {
			if foundOriginalDefinition {
				return true
			}
			foundOriginalDefinition = true
		}
		var r2 *GlobalNamespace_Ref = name.InlineRef2()
		_ = r2
		if ((r2 != nil) && r2.IsSet()) && !(InlineAndCollapseProperties_IsSafeNamespaceReinit(r2)) {
			if foundOriginalDefinition {
				return true
			}
			foundOriginalDefinition = true
		}
		var r3 *GlobalNamespace_Ref = name.InlineRef3()
		_ = r3
		if ((r3 != nil) && r3.IsSet()) && !(InlineAndCollapseProperties_IsSafeNamespaceReinit(r3)) {
			if foundOriginalDefinition {
				return true
			}
		}
		return false
	}
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(name.GetRefs()))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_9)
			_ = ref
			if !(ref.IsSet()) {
				continue
			}
			if InlineAndCollapseProperties_IsSafeNamespaceReinit(ref) {
				continue
			}
			if !(foundOriginalDefinition) {
				foundOriginalDefinition = true
			} else {
				return true
			}
		}
	}
	return false
}

// Java static method: InlineAndCollapseProperties.getSubclassForEs6Superclass(com.google.javascript.rhino.Node superclass)
func InlineAndCollapseProperties_GetSubclassForEs6Superclass(superclass pkg_rhino.Node) pkg_rhino.Node {
	clinit_InlineAndCollapseProperties()
	_ = superclass
	var classNode pkg_rhino.Node = superclass.GetParent()
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(classNode))
	}
	if NodeUtil_IsNameDeclaration(classNode.GetGrandparent()) {
		return classNode.GetParent()
	} else if superclass.GetGrandparent().IsAssign() {
		return classNode.GetPrevious()
	} else if NodeUtil_IsClassDeclaration(classNode) {
		return classNode.GetFirstChild()
	}
	return pkg_rhino.Node(0)
}

// Java static method: InlineAndCollapseProperties.isSafeNamespaceReinit(com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func InlineAndCollapseProperties_IsSafeNamespaceReinit(ref *GlobalNamespace_Ref) bool {
	clinit_InlineAndCollapseProperties()
	_ = ref
	var valParent pkg_rhino.Node = InlineAndCollapseProperties_GetValueParent(ref)
	_ = valParent
	var val pkg_rhino.Node = valParent.GetLastChild()
	_ = val
	if (val != 0) && val.IsOr() {
		var maybeName pkg_rhino.Node = val.GetFirstChild()
		_ = maybeName
		if ref.GetNode().MatchesQualifiedName_1_Node(maybeName) {
			return true
		}
	}
	return false
}

// Java static method: InlineAndCollapseProperties.getValueParent(com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func InlineAndCollapseProperties_GetValueParent(ref *GlobalNamespace_Ref) pkg_rhino.Node {
	clinit_InlineAndCollapseProperties()
	_ = ref
	var n pkg_rhino.Node = ref.GetNode().GetParent()
	_ = n
	if (n != 0) && NodeUtil_IsNameDeclaration(n) {
		return ref.GetNode()
	} else {
		return ref.GetNode().GetParent()
	}
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.Builder
type InlineAndCollapseProperties_Builder struct {
	compiler                      AbstractCompiler               `gson:"compiler"`
	optionalGlobalNamespaceTester *j2gort.Ext_java_util_Optional `gson:"optionalGlobalNamespaceTester"`
	idHash_                       atomic.Int32
	propertyCollapseLevel         CompilerOptions_PropertyCollapseLevel `gson:"propertyCollapseLevel"`
	chunkOutputType               CompilerOptions_ChunkOutputType       `gson:"chunkOutputType"`
	haveModulesBeenRewritten      bool                                  `gson:"haveModulesBeenRewritten"`
	moduleResolutionMode          ModuleLoader_ResolutionMode           `gson:"moduleResolutionMode"`
	testAggressiveInliningOnly    bool                                  `gson:"testAggressiveInliningOnly"`
}

func (this_ *InlineAndCollapseProperties_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_Builder) IsJavaNil_() bool { return this_ == nil }
func (*InlineAndCollapseProperties_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$Builder", "Builder"
}
func (this_ *InlineAndCollapseProperties_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineAndCollapseProperties_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.Builder", "com.google.javascript.jscomp.InlineAndCollapseProperties$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_Builder](v_)
	})
}

func (this_ *InlineAndCollapseProperties_Builder) initFields_InlineAndCollapseProperties_Builder() {
	_ = this_

	this_.optionalGlobalNamespaceTester = j2gort.Ext_java_util_Optional_Empty()
}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewInlineAndCollapseProperties_Builder(compiler AbstractCompiler) *InlineAndCollapseProperties_Builder {
	this_ := &InlineAndCollapseProperties_Builder{}
	InitInlineAndCollapseProperties_Builder(this_, compiler)
	return this_
}

func InitInlineAndCollapseProperties_Builder(this_ *InlineAndCollapseProperties_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_InlineAndCollapseProperties_Builder()
	this_.compiler = compiler
}

// Java method: Builder.setPropertyCollapseLevel(com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel propertyCollapseLevel)
func (this_ *InlineAndCollapseProperties_Builder) SetPropertyCollapseLevel(propertyCollapseLevel CompilerOptions_PropertyCollapseLevel) *InlineAndCollapseProperties_Builder {
	_ = propertyCollapseLevel
	this_.propertyCollapseLevel = propertyCollapseLevel
	return this_
}

// Java method: Builder.setChunkOutputType(com.google.javascript.jscomp.CompilerOptions.ChunkOutputType chunkOutputType)
func (this_ *InlineAndCollapseProperties_Builder) SetChunkOutputType(chunkOutputType CompilerOptions_ChunkOutputType) *InlineAndCollapseProperties_Builder {
	_ = chunkOutputType
	this_.chunkOutputType = chunkOutputType
	return this_
}

// Java method: Builder.setHaveModulesBeenRewritten(boolean haveModulesBeenRewritten)
func (this_ *InlineAndCollapseProperties_Builder) SetHaveModulesBeenRewritten(haveModulesBeenRewritten bool) *InlineAndCollapseProperties_Builder {
	_ = haveModulesBeenRewritten
	this_.haveModulesBeenRewritten = haveModulesBeenRewritten
	return this_
}

// Java method: Builder.setModuleResolutionMode(com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode moduleResolutionMode)
func (this_ *InlineAndCollapseProperties_Builder) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) *InlineAndCollapseProperties_Builder {
	_ = moduleResolutionMode
	this_.moduleResolutionMode = moduleResolutionMode
	return this_
}

// Java method: Builder.testAggressiveInliningOnly(java.util.function.Consumer<com.google.javascript.jscomp.GlobalNamespace> globalNamespaceTester)
func (this_ *InlineAndCollapseProperties_Builder) TestAggressiveInliningOnly(globalNamespaceTester j2gort.Ext_java_util_function_Consumer) *InlineAndCollapseProperties_Builder {
	_ = globalNamespaceTester
	this_.testAggressiveInliningOnly = true
	this_.optionalGlobalNamespaceTester = j2gort.Ext_java_util_Optional_Of(globalNamespaceTester)
	return this_
}

// Java method: Builder.build()
func (this_ *InlineAndCollapseProperties_Builder) Build() *InlineAndCollapseProperties {
	return NewInlineAndCollapseProperties(this_)
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.AggressiveInlineAliases
type InlineAndCollapseProperties_AggressiveInlineAliases struct {
	outer_      *InlineAndCollapseProperties
	idHash_     atomic.Int32
	codeChanged bool `gson:"codeChanged"`
}

func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) IsJavaNil_() bool {
	return this_ == nil
}
func (*InlineAndCollapseProperties_AggressiveInlineAliases) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$AggressiveInlineAliases", "AggressiveInlineAliases"
}
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineAndCollapseProperties_AggressiveInlineAliases) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.AggressiveInlineAliases", "com.google.javascript.jscomp.InlineAndCollapseProperties$AggressiveInlineAliases", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_AggressiveInlineAliases](v_)
	})
}
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) SetOuter_(o *InlineAndCollapseProperties) {
	this_.outer_ = o
}
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) GetOuter_() *InlineAndCollapseProperties {
	return this_.outer_
}
func (*InlineAndCollapseProperties_AggressiveInlineAliases) Implements_jscomp_CompilerPass() {}

func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) initFields_InlineAndCollapseProperties_AggressiveInlineAliases() {
	_ = this_
}

// Java constructor: AggressiveInlineAliases()
func NewInlineAndCollapseProperties_AggressiveInlineAliases(outer_ *InlineAndCollapseProperties) *InlineAndCollapseProperties_AggressiveInlineAliases {
	this_ := &InlineAndCollapseProperties_AggressiveInlineAliases{}
	this_.outer_ = outer_
	InitInlineAndCollapseProperties_AggressiveInlineAliases(this_)
	return this_
}

func InitInlineAndCollapseProperties_AggressiveInlineAliases(this_ *InlineAndCollapseProperties_AggressiveInlineAliases) {
	this_.initFields_InlineAndCollapseProperties_AggressiveInlineAliases()
	this_.codeChanged = true
}

// Java method: AggressiveInlineAliases.getLastUsedGlobalNamespace()
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) GetLastUsedGlobalNamespace() *GlobalNamespace {
	return this_.outer_.namespace
}

// Java method: AggressiveInlineAliases.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.outer_.compiler, root, NewInlineAndCollapseProperties_RewriteSimpleDestructuringAliases())
	this_.outer_.namespace = NewGlobalNamespace_3_LogFile_AbstractCompiler_Node(this_.outer_.decisionsLog, this_.outer_.compiler, root)
	for this_.codeChanged {
		this_.codeChanged = false
		this_.priv_AggressiveInlineAliases_inlineAliases_3219(this_.outer_.namespace)
	}
}

// Java method: AggressiveInlineAliases.inlineAliases(com.google.javascript.jscomp.GlobalNamespace namespace)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_inlineAliases_3219(namespace *GlobalNamespace) {
	_ = namespace
	namespace.EnsureGenerated()
	var workList j2gort.Deque[*GlobalNamespace_Name]
	_ = workList
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(namespace.globalNames))
		_ = size
		for ; i < size; i = i + 1 {
			workList.AddLast(namespace.globalNames[i])
		}
	}
	for !(workList.IsEmpty()) {
		var name *GlobalNamespace_Name = workList.Pop()
		_ = name
		if name.IsGetOrSetDefinition() {
			continue
		}
		if (!(name.InExterns()) && (name.GetGlobalSets() == int32(1))) && (name.GetLocalSets() == int32(0)) {
			this_.priv_AggressiveInlineAliases_maybeInlineInnerName_3221(name)
			if (name.GetAliasingGets() > int32(0)) || (name.GetSubclassingGets() > int32(0)) {
				this_.priv_AggressiveInlineAliases_inlineAliasesForName_3220(name, namespace)
			}
		}
		if (((name.IsObjectLiteral() || name.IsFunction()) || name.IsClass()) && !(InlineAndCollapseProperties_IsUnsafelyReassigned(name))) && name.HasProps() {
			if name.GetAliasingGets() == int32(0) {
				{
					var i_2 int32 = int32(0)
					_ = i_2
					var size_2 int32 = name.GetPropCount()
					_ = size_2
					for ; i_2 < size_2; i_2 = i_2 + 1 {
						workList.AddLast(name.GetPropAt(i_2))
					}
				}
			} else {
				{
					var i_3 int32 = int32(0)
					_ = i_3
					var size_3 int32 = name.GetPropCount()
					_ = size_3
					for ; i_3 < size_3; i_3 = i_3 + 1 {
						var property *GlobalNamespace_Name = name.GetPropAt(i_3)
						_ = property
						if property.CanCollapse() {
							workList.AddLast(property)
						}
					}
				}
			}
		}
	}
}

// Java method: AggressiveInlineAliases.inlineAliasesForName(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.jscomp.GlobalNamespace namespace)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_inlineAliasesForName_3220(name *GlobalNamespace_Name, namespace *GlobalNamespace) {
	_ = name
	_ = namespace
	var refs []*GlobalNamespace_Ref = j2gort.CopyCollectionToSlice[*GlobalNamespace_Ref](name.GetRefs())
	_ = refs
	for _, ref := range refs {
		_ = ref
		if (ref.IsAliasingGet() && !(InlineAndCollapseProperties_MayBeGlobalAlias(ref))) && !(ref.IsTwin()) {
			this_.priv_AggressiveInlineAliases_inlineAliasIfPossible_3224(name, ref, namespace)
		} else if (ref.IsAliasingGet() && j2gort.CastPtr[Scope](ref.scope.GetClosestHoistScope()).IsGlobal()) && !(ref.IsTwin()) {
			this_.priv_AggressiveInlineAliases_inlineGlobalAliasIfPossible_3228(name, ref, namespace)
		} else if (name.IsClass() && ref.IsSubclassingGet()) && name.HasProps() {
			{
				var i int32 = int32(0)
				_ = i
				var size int32 = name.GetPropCount()
				_ = size
				for ; i < size; i = i + 1 {
					this_.priv_AggressiveInlineAliases_rewriteAllSubclassInheritedAccesses_3222(name, ref, name.GetPropAt(i), namespace)
				}
			}
		}
	}

}

// Java method: AggressiveInlineAliases.maybeInlineInnerName(com.google.javascript.jscomp.GlobalNamespace.Name globalName)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_maybeInlineInnerName_3221(globalName *GlobalNamespace_Name) {
	_ = globalName
	var globalNameDeclaration *GlobalNamespace_Ref = j2gort.CheckNotNullPtr(j2gort.CastPtr[GlobalNamespace_Ref](globalName.GetDeclaration()))
	_ = globalNameDeclaration
	var globalDeclarationNode pkg_rhino.Node = j2gort.CheckNotNullNode(globalNameDeclaration.GetNode())
	_ = globalDeclarationNode
	var valueNode pkg_rhino.Node = NodeUtil_GetRValueOfLValue(globalDeclarationNode)
	_ = valueNode
	if valueNode == 0 {
		return
	}
	var innerNameNode pkg_rhino.Node = InlineAndCollapseProperties_MaybeGetInnerNameNode(valueNode)
	_ = innerNameNode
	if innerNameNode == 0 {
		return
	}
	var innerName string = innerNameNode.GetString()
	_ = innerName
	clinit_SyntacticScopeCreator()
	var val_12 SyntacticScopeCreator
	InitSyntacticScopeCreator_1(&val_12, this_.outer_.compiler)
	var syntacticScopeCreator *SyntacticScopeCreator = &val_12
	_ = syntacticScopeCreator
	defer syntacticScopeCreator.ReleaseSlabs_()

	var innerScope *Scope = syntacticScopeCreator.CreateScope_2_Node_Scope(valueNode, globalNameDeclaration.scope)
	_ = innerScope
	var innerNameVar *Var = j2gort.CheckNotNullPtr(j2gort.CastPtr[Var](innerScope.GetVar(innerName)))
	_ = innerNameVar
	clinit_ReferenceCollector()
	var val_13 ReferenceCollector
	InitReferenceCollector_4(&val_13, this_.outer_.compiler, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), syntacticScopeCreator, j2gort.Ext_com_google_common_base_Predicates_EqualTo(innerNameVar))
	var collector *ReferenceCollector = &val_13
	_ = collector
	defer collector.ReleaseSlabs_()

	collector.ProcessScope(innerScope)
	var innerNameRefs *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](collector.GetReferences(innerNameVar))
	_ = innerNameRefs
	var newNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = newNodes
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(innerNameRefs))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			innerNameRef := j2gort.CastPtr[Reference](rawElem_15)
			_ = innerNameRef
			var innerNameRefNode pkg_rhino.Node = innerNameRef.GetNode()
			_ = innerNameRefNode
			if NodeUtil_IsNormalGet(innerNameRefNode.GetParent()) {
				newNodes.Add_1(this_.priv_AggressiveInlineAliases_replaceAliasReference_3227(globalNameDeclaration, innerNameRef))
			}
		}
	}
	this_.outer_.namespace.ScanNewNodes(newNodes)
}

// Java method: AggressiveInlineAliases.rewriteAllSubclassInheritedAccesses(com.google.javascript.jscomp.GlobalNamespace.Name superclassNameObj, com.google.javascript.jscomp.GlobalNamespace.Ref superclassRef, com.google.javascript.jscomp.GlobalNamespace.Name prop, com.google.javascript.jscomp.GlobalNamespace namespace)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_rewriteAllSubclassInheritedAccesses_3222(superclassNameObj *GlobalNamespace_Name, superclassRef *GlobalNamespace_Ref, prop *GlobalNamespace_Name, namespace *GlobalNamespace) bool {
	_ = superclassNameObj
	_ = superclassRef
	_ = prop
	_ = namespace
	if !(prop.CanCollapse()) {
		return false
	}
	var subclass pkg_rhino.Node = InlineAndCollapseProperties_GetSubclassForEs6Superclass(superclassRef.GetNode())
	_ = subclass
	if (subclass == 0) || !(subclass.IsQualifiedName()) {
		return false
	}
	var subclassName string = subclass.GetQualifiedName()
	_ = subclassName
	var subclassQualifiedPropName string = (j2gort.StrValString(subclassName) + "." + j2gort.StrValString(prop.GetBaseName()))
	_ = subclassQualifiedPropName
	var subclassPropNameObj *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](namespace.GetOwnSlot(subclassQualifiedPropName))
	_ = subclassPropNameObj
	if (subclassPropNameObj != nil) && ((subclassPropNameObj.GetLocalSets() > int32(0)) || (subclassPropNameObj.GetGlobalSets() > int32(0))) {
		return false
	}
	var subclassNameObj *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](namespace.GetOwnSlot(subclassName))
	_ = subclassNameObj
	if (subclassNameObj != nil) && (subclassNameObj.SubclassingGetCount() > int32(0)) {
		{
			cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(subclassNameObj.GetRefs()))
			for cur_16.HasNext() {
				rawElem_17 := cur_16.Next()
				ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_17)
				_ = ref
				if ref.IsSubclassingGet() {
					this_.priv_AggressiveInlineAliases_rewriteAllSubclassInheritedAccesses_3222(superclassNameObj, ref, prop, namespace)
				}
			}
		}
	}
	if subclassPropNameObj != nil {
		var newNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
		_ = newNodes
		var superclassNameNode pkg_rhino.Node = j2gort.CastPtr[GlobalNamespace_Ref](superclassNameObj.GetDeclaration()).GetNode()
		_ = superclassNameNode
		if superclassNameNode.IsName() {
			superclassNameNode = superclassNameNode.CloneNode_0()
		} else if superclassNameNode.IsGetProp() {
			superclassNameNode = superclassNameNode.CloneTree_0()
		} else {
			return false
		}
		this_.priv_AggressiveInlineAliases_rewriteNestedAliasReference_3231(superclassNameNode, int32(0), newNodes, subclassPropNameObj)
		namespace.ScanNewNodes(newNodes)
	}
	return true
}

// Java method: AggressiveInlineAliases.isInlineableModuleExportsAlias(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.jscomp.ReferenceCollection aliasRefs)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_isInlineableModuleExportsAlias_3223(name *GlobalNamespace_Name, aliasRefs *ReferenceCollection) bool {
	_ = name
	_ = aliasRefs
	var aliasedNameStr string = name.GetFullName()
	_ = aliasedNameStr
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(aliasedNameStr), "module$exports$")) {
		return false
	}
	var firstRhs pkg_rhino.Node = pkg_rhino.Node(0)
	_ = firstRhs
	var size int32 = aliasRefs.Size()
	_ = size
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			var ref *Reference = aliasRefs.Get(i)
			_ = ref
			if ref.IsVarDeclaration() || ref.GetParent().IsAssign() {
				var rhs pkg_rhino.Node = pkg_rhino.Node(0)
				_ = rhs
				if ref.GetParent().IsAssign() {
					rhs = ref.GetParent().GetSecondChild()
				} else {
					rhs = ref.GetNode().GetFirstChild()
				}
				if rhs != 0 {
					if firstRhs == 0 {
						if !(rhs.MatchesQualifiedName_1_string(aliasedNameStr)) {
							return false
						}
						firstRhs = rhs
					} else {
						if !(rhs.MatchesQualifiedName_1_Node(firstRhs)) {
							return false
						}
					}
				}
				continue
			}
			if !(ref.IsDotPropertyAccess()) && !(ref.IsAssignedToObjectDestructuringPattern()) {
				return false
			}
		}
	}
	return true
}

// Java method: AggressiveInlineAliases.inlineAliasIfPossible(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.jscomp.GlobalNamespace.Ref alias, com.google.javascript.jscomp.GlobalNamespace namespace)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_inlineAliasIfPossible_3224(name *GlobalNamespace_Name, alias *GlobalNamespace_Ref, namespace *GlobalNamespace) {
	_ = name
	_ = alias
	_ = namespace
	var aliasParent pkg_rhino.Node = alias.GetNode().GetParent()
	_ = aliasParent
	if aliasParent.IsName() || aliasParent.IsAssign() {
		var aliasLhsNode pkg_rhino.Node
		_ = aliasLhsNode
		if aliasParent.IsName() {
			aliasLhsNode = aliasParent
		} else {
			aliasLhsNode = aliasParent.GetFirstChild()
		}
		var aliasVarName string = aliasLhsNode.GetString()
		_ = aliasVarName
		var aliasVar *Var = j2gort.CastPtr[Var](alias.scope.GetVar(aliasVarName))
		_ = aliasVar
		if !(aliasVar != nil) {
			j2gort.CheckStateFmt(false, "Expected variable to be defined in scope (%s)", j2gort.BoxString(aliasVarName))
		}
		clinit_ReferenceCollector()
		var val_18 ReferenceCollector
		InitReferenceCollector_4(&val_18, this_.outer_.compiler, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(this_.outer_.compiler), j2gort.Ext_com_google_common_base_Predicates_EqualTo(j2gort.PtrToAny(aliasVar)))
		var collector *ReferenceCollector = &val_18
		_ = collector
		collector.ownsScopeCreator_ = true
		defer collector.ReleaseSlabs_()

		var aliasScope *Scope = j2gort.CastPtr[Scope](aliasVar.GetScope())
		_ = aliasScope
		collector.ProcessScope(aliasScope)
		var aliasRefs *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](collector.GetReferences(j2gort.PtrToAny(aliasVar)))
		_ = aliasRefs
		var newNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
		_ = newNodes
		if aliasRefs.IsWellDefined() && (aliasRefs.IsAssignedOnceInLifetime() || this_.priv_AggressiveInlineAliases_isInlineableModuleExportsAlias_3223(name, aliasRefs)) {
			var readNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
			_ = readNodes
			var size int32 = aliasRefs.Size()
			_ = size
			var firstRead int32
			_ = firstRead
			if aliasRefs.Get(int32(0)).IsInitializingDeclaration() {
				firstRead = int32(1)
			} else {
				firstRead = int32(2)
			}
			{
				var i int32 = firstRead
				_ = i
				for ; i < size; i = i + 1 {
					var aliasRef *Reference = aliasRefs.Get(i)
					_ = aliasRef
					readNodes.Add_1(this_.priv_AggressiveInlineAliases_replaceAliasReference_3227(alias, aliasRef))
				}
			}
			this_.priv_AggressiveInlineAliases_replaceAliasingAssignment_3226(alias, name, aliasLhsNode, newNodes)
			newNodes.AddAll_1(readNodes)
			namespace.ScanNewNodes(newNodes)
			return
		}
		if name.IsConstructor() {
			if !(this_.priv_AggressiveInlineAliases_partiallyInlineAlias_3225(alias, name, namespace, aliasRefs, aliasLhsNode)) {
				if InlineAndCollapseProperties_ReferencesCollapsibleProperty(aliasRefs, name, namespace) {
					this_.outer_.compiler.Report_1(JSError_Make_3(aliasParent, StaticGet_InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING(), []string{aliasVarName}))
				}
			}
		}
	}
}

// Java method: AggressiveInlineAliases.partiallyInlineAlias(com.google.javascript.jscomp.GlobalNamespace.Ref alias, com.google.javascript.jscomp.GlobalNamespace.Name aliasingName, com.google.javascript.jscomp.GlobalNamespace namespace, com.google.javascript.jscomp.ReferenceCollection aliasRefs, com.google.javascript.rhino.Node aliasLhsNode)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_partiallyInlineAlias_3225(alias *GlobalNamespace_Ref, aliasingName *GlobalNamespace_Name, namespace *GlobalNamespace, aliasRefs *ReferenceCollection, aliasLhsNode pkg_rhino.Node) bool {
	_ = alias
	_ = aliasingName
	_ = namespace
	_ = aliasRefs
	_ = aliasLhsNode
	var aliasBlock *BasicBlock = nil
	_ = aliasBlock
	{
		cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(aliasRefs))
		for cur_19.HasNext() {
			rawElem_20 := cur_19.Next()
			aliasRef := j2gort.CastPtr[Reference](rawElem_20)
			_ = aliasRef
			var aliasRefNode pkg_rhino.Node = aliasRef.GetNode()
			_ = aliasRefNode
			if aliasRefNode == aliasLhsNode {
				aliasBlock = aliasRef.GetBasicBlock()
				continue
			} else if aliasRef.IsLvalue() {
				return false
			}
		}
	}
	var readNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = readNodes
	var alreadySeenInitialAlias bool = false
	_ = alreadySeenInitialAlias
	var foundNonReplaceableAlias bool = false
	_ = foundNonReplaceableAlias
	{
		cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(aliasRefs))
		for cur_21.HasNext() {
			rawElem_22 := cur_21.Next()
			aliasRef_2 := j2gort.CastPtr[Reference](rawElem_22)
			_ = aliasRef_2
			var aliasRefNode_2 pkg_rhino.Node = aliasRef_2.GetNode()
			_ = aliasRefNode_2
			if aliasRefNode_2 == aliasLhsNode {
				alreadySeenInitialAlias = true
				continue
			} else if aliasRef_2.IsDeclaration() {
				continue
			}
			var refBlock *BasicBlock = aliasRef_2.GetBasicBlock()
			_ = refBlock
			if ((refBlock != aliasBlock) && aliasBlock.ProvablyExecutesBefore(refBlock)) || ((refBlock == aliasBlock) && alreadySeenInitialAlias) {
				this_.codeChanged = true
				readNodes.Add_1(this_.priv_AggressiveInlineAliases_replaceAliasReference_3227(alias, aliasRef_2))
			} else {
				foundNonReplaceableAlias = true
			}
		}
	}
	var newNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = newNodes
	if !(foundNonReplaceableAlias) {
		this_.priv_AggressiveInlineAliases_replaceAliasingAssignment_3226(alias, aliasingName, aliasLhsNode, newNodes)
	}
	newNodes.AddAll_1(readNodes)
	if this_.codeChanged {
		namespace.ScanNewNodes(newNodes)
	}
	return !(foundNonReplaceableAlias)
}

// Java method: AggressiveInlineAliases.replaceAliasingAssignment(com.google.javascript.jscomp.GlobalNamespace.Ref alias, com.google.javascript.jscomp.GlobalNamespace.Name aliasName, com.google.javascript.rhino.Node aliasLhsNode, java.util.Set<com.google.javascript.jscomp.GlobalNamespace.AstChange> newNodes)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_replaceAliasingAssignment_3226(alias *GlobalNamespace_Ref, aliasName *GlobalNamespace_Name, aliasLhsNode pkg_rhino.Node, newNodes j2gort.Ext_java_util_Set) {
	_ = alias
	_ = aliasName
	_ = aliasLhsNode
	_ = newNodes
	var assignment pkg_rhino.Node = aliasLhsNode.GetParent()
	_ = assignment
	if !(NodeUtil_IsNameDeclaration(assignment)) && NodeUtil_IsExpressionResultUsed(assignment) {
		if !(assignment.IsAssign()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(assignment))
		}
		var assignmentParent pkg_rhino.Node = assignment.GetParent()
		_ = assignmentParent
		var rhs pkg_rhino.Node = alias.GetNode().Detach()
		_ = rhs
		assignment.ReplaceWith(rhs)
		newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(alias, rhs))
		this_.outer_.compiler.ReportChangeToEnclosingScope(assignmentParent)
	} else {
		var aliasParent pkg_rhino.Node = alias.GetNode().GetParent()
		_ = aliasParent
		alias.GetNode().ReplaceWith(pkg_rhino.IR_NullNode())
		this_.outer_.compiler.ReportChangeToEnclosingScope(aliasParent)
	}
	aliasName.RemoveRef(alias)
	this_.codeChanged = true
}

// Java method: AggressiveInlineAliases.replaceAliasReference(com.google.javascript.jscomp.GlobalNamespace.Ref alias, com.google.javascript.jscomp.Reference aliasRef)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_replaceAliasReference_3227(alias *GlobalNamespace_Ref, aliasRef *Reference) GlobalNamespace_AstChange {
	_ = alias
	_ = aliasRef
	var originalRefNode pkg_rhino.Node = alias.GetNode()
	_ = originalRefNode
	var nodeToReplace pkg_rhino.Node = aliasRef.GetNode()
	_ = nodeToReplace
	if !(nodeToReplace.IsQualifiedName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(nodeToReplace))
	}
	var newNode pkg_rhino.Node
	_ = newNode
	if originalRefNode.IsName() {
		newNode = originalRefNode.CloneNode_0()
	} else {
		newNode = originalRefNode.CloneTree_0()
	}
	newNode.SrcrefTree(nodeToReplace)
	nodeToReplace.ReplaceWith(newNode)
	this_.outer_.compiler.ReportChangeToEnclosingScope(newNode)
	return NewGlobalNamespace_SimpleAstChange(newNode, alias.GetChunk(), alias.scope)
}

// Java method: AggressiveInlineAliases.inlineGlobalAliasIfPossible(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.jscomp.GlobalNamespace.Ref alias, com.google.javascript.jscomp.GlobalNamespace namespace)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_inlineGlobalAliasIfPossible_3228(name *GlobalNamespace_Name, alias *GlobalNamespace_Ref, namespace *GlobalNamespace) {
	_ = name
	_ = alias
	_ = namespace
	var aliasParent pkg_rhino.Node = alias.GetNode().GetParent()
	_ = aliasParent
	if ((aliasParent.IsAssign() || aliasParent.IsName()) && NodeUtil_IsExecutedExactlyOnce(aliasParent)) || (aliasParent.IsName() && name.IsConstructor()) {
		var lvalue pkg_rhino.Node
		_ = lvalue
		if aliasParent.IsName() {
			lvalue = aliasParent
		} else {
			lvalue = aliasParent.GetFirstChild()
		}
		if !(lvalue.IsQualifiedName()) {
			return
		}
		if lvalue.IsName() && this_.outer_.compiler.GetCodingConvention().IsExported_2(lvalue.GetString(), false) {
			return
		}
		var aliasingName *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](namespace.GetSlot(lvalue.GetQualifiedName()))
		_ = aliasingName
		if aliasingName == nil {
			return
		}
		if (j2gort.CheckNotNullPtr(name) == aliasingName) && aliasParent.IsAssign() {
			return
		}
		var aliasInlinability GlobalNamespace_Inlinability = aliasingName.CalculateInlinability()
		_ = aliasInlinability
		if !(aliasInlinability.ShouldInlineUsages()) {
			return
		}
		var readNodes j2gort.Ext_java_util_Set = this_.priv_AggressiveInlineAliases_rewriteAliasReferences_3229(aliasingName, alias, aliasInlinability)
		_ = readNodes
		this_.priv_AggressiveInlineAliases_rewriteAliasProps_3230(aliasingName, alias.GetNode(), int32(0), readNodes)
		var newNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
		_ = newNodes
		if aliasInlinability.ShouldRemoveDeclaration() {
			var aliasInitialization *GlobalNamespace_Ref = aliasingName.GetInitialization()
			_ = aliasInitialization
			if aliasInitialization.IsTwin() {
				if !(aliasParent.IsAssign()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(aliasParent))
				}
				var aliasGrandparent pkg_rhino.Node = aliasParent.GetParent()
				_ = aliasGrandparent
				aliasParent.ReplaceWith(alias.GetNode().Detach())
				aliasingName.RemoveRef(aliasInitialization)
				newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(alias, alias.GetNode()))
				this_.outer_.compiler.ReportChangeToEnclosingScope(aliasGrandparent)
			} else {
				alias.GetNode().ReplaceWith(pkg_rhino.IR_NullNode())
				this_.outer_.compiler.ReportChangeToEnclosingScope(aliasParent)
			}
			this_.codeChanged = true
			name.RemoveRef(alias)
		}
		newNodes.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](readNodes))
		namespace.ScanNewNodes(newNodes)
	}
}

// Java method: AggressiveInlineAliases.rewriteAliasReferences(com.google.javascript.jscomp.GlobalNamespace.Name aliasingName, com.google.javascript.jscomp.GlobalNamespace.Ref aliasingRef, com.google.javascript.jscomp.GlobalNamespace.Inlinability aliasInlinability)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_rewriteAliasReferences_3229(aliasingName *GlobalNamespace_Name, aliasingRef *GlobalNamespace_Ref, aliasInlinability GlobalNamespace_Inlinability) j2gort.Ext_java_util_Set {
	_ = aliasingName
	_ = aliasingRef
	_ = aliasInlinability
	var newNodes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = newNodes
	var refs []*GlobalNamespace_Ref = j2gort.CopyCollectionToSlice[*GlobalNamespace_Ref](aliasingName.GetRefs())
	_ = refs
	var initializationChunk *JSChunk
	_ = initializationChunk
	if aliasingName.GetInitialization() != nil {
		initializationChunk = aliasingName.GetInitialization().GetChunk()
	} else {
		initializationChunk = j2gort.CastPtr[GlobalNamespace_Ref](aliasingName.GetDeclaration()).GetChunk()
	}
	for _, ref := range refs {
		_ = ref
		if ref.IsSetFromGlobal() {
			if !((j2gort.CheckNotNullPtr(ref.GetChunk()) == initializationChunk) || ref.IsUninitializedDeclaration()) {
				j2gort.CheckStateMsg(false, j2gort.PtrToAny(ref))
			}
			continue
		}
		if !(ref.IsGet()) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(ref))
		}
		if !(!(ref.IsTwin())) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(ref))
		}
		if ((aliasInlinability == StaticGet_GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY()) && (ref.GetChunk() != initializationChunk)) && !(this_.outer_.compiler.GetChunkGraph().DependsOn(ref.GetChunk(), initializationChunk)) {
			continue
		}
		if ref.GetNode().IsStringKey() {
			DestructuringGlobalNameExtractor_ReassignDestructringLvalue(ref.GetNode(), aliasingRef.GetNode().CloneTree_0(), newNodes, ref, this_.outer_.compiler)
		} else {
			if !(ref.GetNode().IsGetProp() || ref.GetNode().IsName()) {
				j2gort.CheckState(false)
			}
			var newNode pkg_rhino.Node = aliasingRef.GetNode().CloneTree_0()
			_ = newNode
			var node pkg_rhino.Node = ref.GetNode()
			_ = node
			newNode.Srcref(node)
			node.ReplaceWith(newNode)
			this_.outer_.compiler.ReportChangeToEnclosingScope(newNode)
			newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(ref, newNode))
		}
		aliasingName.RemoveRef(ref)
	}

	return newNodes
}

// Java method: AggressiveInlineAliases.rewriteAliasProps(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.rhino.Node value, int depth, java.util.Set<com.google.javascript.jscomp.GlobalNamespace.AstChange> newNodes)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_rewriteAliasProps_3230(name *GlobalNamespace_Name, value pkg_rhino.Node, depth int32, newNodes j2gort.Ext_java_util_Set) {
	_ = name
	_ = value
	_ = depth
	_ = newNodes
	if !(name.HasProps()) {
		return
	}
	if !(!(value.MatchesQualifiedName_1_string(name.GetFullName()))) {
		j2gort.CheckStateFmt(false, "%s should not match name %s", j2gort.BoxNode(value), j2gort.BoxString(name.GetFullName()))
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = name.GetPropCount()
		_ = size
		for ; i < size; i = i + 1 {
			this_.priv_AggressiveInlineAliases_rewriteNestedAliasReference_3231(value, depth, newNodes, name.GetPropAt(i))
		}
	}
}

// Java method: AggressiveInlineAliases.rewriteNestedAliasReference(com.google.javascript.rhino.Node value, int depth, java.util.Set<com.google.javascript.jscomp.GlobalNamespace.AstChange> newNodes, com.google.javascript.jscomp.GlobalNamespace.Name prop)
func (this_ *InlineAndCollapseProperties_AggressiveInlineAliases) priv_AggressiveInlineAliases_rewriteNestedAliasReference_3231(value pkg_rhino.Node, depth int32, newNodes j2gort.Ext_java_util_Set, prop *GlobalNamespace_Name) {
	_ = value
	_ = depth
	_ = newNodes
	_ = prop
	this_.priv_AggressiveInlineAliases_rewriteAliasProps_3230(prop, value, (depth + int32(1)), newNodes)
	var refs []*GlobalNamespace_Ref = j2gort.CopyCollectionToSlice[*GlobalNamespace_Ref](prop.GetRefs())
	_ = refs
	for _, ref := range refs {
		_ = ref
		var target pkg_rhino.Node = ref.GetNode()
		_ = target
		if target.IsStringKey() && target.GetParent().IsDestructuringPattern() {
			if !(target.GetGrandparent().IsAssign() || target.GetGrandparent().IsDestructuringLhs()) {
				j2gort.CheckStateFmt(false, "Did not expect GlobalNamespace to create Ref for key in nested object pattern %s", j2gort.BoxNode(target))
			}
			continue
		}
		if target.IsMemberFunctionDef() {
			continue
		}
		{
			var i int32 = int32(0)
			_ = i
			for ; i <= depth; i = i + 1 {
				if target.IsGetProp() {
					target = target.GetFirstChild()
				} else if NodeUtil_IsObjectLitKey(target) {
					var gparent pkg_rhino.Node = target.GetGrandparent()
					_ = gparent
					if gparent.IsAssign() {
						target = gparent.GetFirstChild()
					} else {
						if !(NodeUtil_IsObjectLitKey(gparent)) {
							j2gort.CheckState(false)
						}
						target = gparent
					}
				} else {
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected node: "+j2gort.StrVal(target))))
				}
			}
		}
		if !(target.IsGetProp() || target.IsName()) {
			j2gort.CheckState(false)
		}
		var newValue pkg_rhino.Node = value.CloneTree_0()
		_ = newValue
		target.ReplaceWith(newValue)
		this_.outer_.compiler.ReportChangeToEnclosingScope(newValue)
		prop.RemoveRef(ref)
		newNodes.Add_1(NewGlobalNamespace_RefBasedAstChange(ref, ref.GetNode()))
		this_.codeChanged = true
	}

}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.RewriteSimpleDestructuringAliases
type InlineAndCollapseProperties_RewriteSimpleDestructuringAliases struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func (this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) IsJavaNil_() bool {
	return this_ == nil
}
func (*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$RewriteSimpleDestructuringAliases", "RewriteSimpleDestructuringAliases"
}
func (this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.RewriteSimpleDestructuringAliases", "com.google.javascript.jscomp.InlineAndCollapseProperties$RewriteSimpleDestructuringAliases", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases](v_)
	})
}
func (*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) initFields_InlineAndCollapseProperties_RewriteSimpleDestructuringAliases() {
	_ = this_
}

// Java constructor: RewriteSimpleDestructuringAliases()
func NewInlineAndCollapseProperties_RewriteSimpleDestructuringAliases() *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases {
	this_ := &InlineAndCollapseProperties_RewriteSimpleDestructuringAliases{}
	InitInlineAndCollapseProperties_RewriteSimpleDestructuringAliases(this_)
	return this_
}

func InitInlineAndCollapseProperties_RewriteSimpleDestructuringAliases(this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineAndCollapseProperties_RewriteSimpleDestructuringAliases()
}

// Java static method: RewriteSimpleDestructuringAliases.isSimpleDestructuringAlias(com.google.javascript.rhino.Node n)
func InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringAlias(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsStatement(n)) || !(n.IsConst()) {
		return false
	}
	if !(n.HasOneChild()) {
		j2gort.CheckState(false)
	}
	var destructuringLhs pkg_rhino.Node = n.GetFirstChild()
	_ = destructuringLhs
	if !(destructuringLhs.IsDestructuringLhs()) {
		return false
	}
	var objectPattern pkg_rhino.Node = destructuringLhs.GetFirstChild()
	_ = objectPattern
	if !(objectPattern.IsObjectPattern()) {
		return false
	}
	var rhs pkg_rhino.Node = destructuringLhs.GetLastChild()
	_ = rhs
	if !(rhs.IsQualifiedName()) {
		return false
	}
	return InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringPattern(objectPattern)
}

// Java static method: RewriteSimpleDestructuringAliases.isSimpleDestructuringPattern(com.google.javascript.rhino.Node objectPattern)
func InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringPattern(objectPattern pkg_rhino.Node) bool {
	_ = objectPattern
	if !(objectPattern.IsObjectPattern()) {
		j2gort.CheckArgument(false)
	}
	{
		var key pkg_rhino.Node = objectPattern.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if !(key.IsStringKey()) || key.IsQuotedStringKey() {
				return false
			}
			if !(key.HasOneChild()) {
				j2gort.CheckState(false)
			}
			var rhs pkg_rhino.Node = key.GetFirstChild()
			_ = rhs
			if !(rhs.IsObjectPattern()) && !(rhs.IsName()) {
				return false
			}
			if rhs.IsObjectPattern() && !(InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringPattern(rhs)) {
				return false
			}
		}
	}
	return true
}

// Java method: RewriteSimpleDestructuringAliases.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringAlias(n)) {
		return
	}
	var insertionPoint pkg_rhino.Node = n
	_ = insertionPoint
	var destructuringLhs pkg_rhino.Node = n.GetFirstChild()
	_ = destructuringLhs
	var objectPattern pkg_rhino.Node = destructuringLhs.GetFirstChild()
	_ = objectPattern
	var rhs pkg_rhino.Node = destructuringLhs.GetLastChild()
	_ = rhs
	if !(rhs.IsName()) && objectPattern.HasMoreThanOneChild() {
		var uniqueId string = t.GetCompiler().GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		_ = uniqueId
		var tempName pkg_rhino.Node = pkg_rhino.IR_Name(("destructuring$" + j2gort.StrValString(uniqueId))).Srcref(rhs)
		_ = tempName
		tempName.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		var tempConstNode pkg_rhino.Node = pkg_rhino.IR_ConstNode(tempName, rhs.Detach()).Srcref(objectPattern)
		_ = tempConstNode
		tempConstNode.InsertAfter(insertionPoint)
		insertionPoint = tempConstNode
		rhs = tempName.CloneNode_0()
	}
	var unusedNewInsertionPoint pkg_rhino.Node = this_.priv_RewriteSimpleDestructuringAliases_expandObjectPattern_3232(t, insertionPoint, objectPattern, rhs)
	_ = unusedNewInsertionPoint
	n.Detach()
	t.ReportCodeChange_0()
}

// Java method: RewriteSimpleDestructuringAliases.expandObjectPattern(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node insertionPoint, com.google.javascript.rhino.Node objectPattern, com.google.javascript.rhino.Node rhs)
func (this_ *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) priv_RewriteSimpleDestructuringAliases_expandObjectPattern_3232(t *NodeTraversal, insertionPoint pkg_rhino.Node, objectPattern pkg_rhino.Node, rhs pkg_rhino.Node) pkg_rhino.Node {
	_ = t
	_ = insertionPoint
	_ = objectPattern
	_ = rhs
	{
		var key pkg_rhino.Node = objectPattern.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			var keyChild pkg_rhino.Node = key.GetFirstChild()
			_ = keyChild
			var nameNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = nameNode
			if keyChild.IsName() {
				nameNode = keyChild.Detach()
			} else {
				if !(keyChild.IsObjectPattern()) {
					j2gort.CheckState(false)
				}
				var uniqueId string = t.GetCompiler().GetUniqueIdSupplier().GetUniqueId(t.GetInput())
				_ = uniqueId
				nameNode = pkg_rhino.IR_Name(("destructuring$" + j2gort.StrValString(uniqueId))).Srcref(keyChild)
				nameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
			}
			var newRhs pkg_rhino.Node = pkg_rhino.IR_Getprop_2(rhs.CloneTree_0(), key.GetString()).Srcref(keyChild)
			_ = newRhs
			var newConstNode pkg_rhino.Node = pkg_rhino.IR_ConstNode(nameNode, newRhs).Srcref(objectPattern)
			_ = newConstNode
			newConstNode.InsertAfter(insertionPoint)
			insertionPoint = newConstNode
			if keyChild.IsObjectPattern() {
				insertionPoint = this_.priv_RewriteSimpleDestructuringAliases_expandObjectPattern_3232(t, insertionPoint, keyChild, nameNode.CloneNode_0())
			}
		}
	}
	return insertionPoint
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.CollapseProperties
type InlineAndCollapseProperties_CollapseProperties struct {
	outer_                     *InlineAndCollapseProperties
	dynamicallyImportedModules j2gort.OrderedSet[string]           `gson:"dynamicallyImportedModules"`
	prefixOriginalNameBuf      *j2gort.Ext_java_lang_StringBuilder `gson:"prefixOriginalNameBuf"`
	idHash_                    atomic.Int32
}

func (this_ *InlineAndCollapseProperties_CollapseProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_CollapseProperties) IsJavaNil_() bool { return this_ == nil }
func (*InlineAndCollapseProperties_CollapseProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$CollapseProperties", "CollapseProperties"
}
func (this_ *InlineAndCollapseProperties_CollapseProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineAndCollapseProperties_CollapseProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.CollapseProperties", "com.google.javascript.jscomp.InlineAndCollapseProperties$CollapseProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_CollapseProperties](v_)
	})
}
func (this_ *InlineAndCollapseProperties_CollapseProperties) SetOuter_(o *InlineAndCollapseProperties) {
	this_.outer_ = o
}
func (this_ *InlineAndCollapseProperties_CollapseProperties) GetOuter_() *InlineAndCollapseProperties {
	return this_.outer_
}
func (*InlineAndCollapseProperties_CollapseProperties) Implements_jscomp_CompilerPass() {}

func (this_ *InlineAndCollapseProperties_CollapseProperties) initFields_InlineAndCollapseProperties_CollapseProperties() {
	_ = this_
	this_.dynamicallyImportedModules = j2gort.NewOrderedSet[string]()
	this_.prefixOriginalNameBuf = j2gort.NewExt_java_lang_StringBuilder_1_int32(int32(64))
}

// Java constructor: CollapseProperties()
func NewInlineAndCollapseProperties_CollapseProperties(outer_ *InlineAndCollapseProperties) *InlineAndCollapseProperties_CollapseProperties {
	this_ := &InlineAndCollapseProperties_CollapseProperties{}
	this_.outer_ = outer_
	InitInlineAndCollapseProperties_CollapseProperties(this_)
	return this_
}

func InitInlineAndCollapseProperties_CollapseProperties(this_ *InlineAndCollapseProperties_CollapseProperties) {
	this_.initFields_InlineAndCollapseProperties_CollapseProperties()
}

// Java method: CollapseProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties_CollapseProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if (this_.outer_.propertyCollapseLevel == StaticGet_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT()) || (this_.outer_.chunkOutputType == StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES()) {
		NodeTraversal_Traverse(this_.outer_.compiler, root, NewInlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules(this_, this_.outer_.haveModulesBeenRewritten, this_.outer_.moduleResolutionMode))
	}
	j2gort.CheckNotNullPtr(this_.outer_.namespace).EnsureGenerated()
	var escaped j2gort.OrderedSet[*GlobalNamespace_Name] = this_.priv_CollapseProperties_checkNamespaces_3235()
	_ = escaped
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(this_.outer_.namespace.globalNames))
		_ = size
		for ; i < size; i = i + 1 {
			var name *GlobalNamespace_Name = this_.outer_.namespace.globalNames[i]
			_ = name
			this_.priv_CollapseProperties_flattenReferencesToCollapsibleDescendantNames_3238(name, name.GetBaseName(), &escaped)
			this_.priv_CollapseProperties_collapseDeclarationOfNameAndDescendants_3249(name, name.GetBaseName(), &escaped)
		}
	}
}

// Java method: CollapseProperties.canCollapse(com.google.javascript.jscomp.GlobalNamespace.Name name)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_canCollapse_3233(name *GlobalNamespace_Name) bool {
	_ = name
	var inlinability GlobalNamespace_Inlinability = name.CanCollapseOrInline()
	_ = inlinability
	if !(inlinability.CanCollapse()) {
		this_.priv_CollapseProperties_logDecisionForName_3239(name, inlinability, "canCollapse() returns false")
		return false
	}
	if this_.outer_.propertyCollapseLevel == StaticGet_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT() {
		if !(name.IsModuleExport()) {
			this_.priv_CollapseProperties_logDecisionForName_3239(name, inlinability, "module export: canCollapse() returns false")
			return false
		} else if this_.dynamicallyImportedModules.Contains(name.GetBaseName()) {
			this_.priv_CollapseProperties_logDecisionForName_3239(name, inlinability, "dynamic module export: canCollapse() returns false")
			return false
		}
	}
	this_.priv_CollapseProperties_logDecisionForName_3239(name, inlinability, "canCollapse() returns true")
	return true
}

// Java method: CollapseProperties.canEliminate(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.rhino.Node rvalue)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_canEliminate_3234(name *GlobalNamespace_Name, rvalue pkg_rhino.Node) bool {
	_ = name
	_ = rvalue
	if !(name.CanEliminate()) {
		return false
	}
	if rvalue.IsObjectLit() {
		{
			var child pkg_rhino.Node = rvalue.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				switch child.TokenSwitchOrd_() {
				case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(128) /* Token.OBJECT_SPREAD */ :
					return false
				case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
				default:
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected child of OBJECTLIT: "+j2gort.StrValString(child.ToStringTree()))))
				}
			}
		}
	}
	if !(name.HasProps()) || (this_.outer_.propertyCollapseLevel != StaticGet_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT()) {
		return true
	}
	return false
}

// Java method: CollapseProperties.checkNamespaces()
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_checkNamespaces_3235() j2gort.OrderedSet[*GlobalNamespace_Name] {
	var escaped j2gort.OrderedSet[*GlobalNamespace_Name] = j2gort.NewOrderedSet[*GlobalNamespace_Name]()
	_ = escaped
	var dynamicallyImportedModuleRefs *j2gort.OrderedSet[string]
	_ = dynamicallyImportedModuleRefs
	if this_.dynamicallyImportedModules.IsEmpty() {
		dynamicallyImportedModuleRefs = nil
	} else {
		dynamicallyImportedModuleRefs = j2gort.OrderedSetPtrFromSet(&this_.dynamicallyImportedModules)
	}
	if dynamicallyImportedModuleRefs != nil {
		{
			cur_27 := this_.outer_.namespace.nameMap.IterValues()
			for cur_27.HasNext() {

				name := cur_27.
					Next()
				_ = name
				if this_.dynamicallyImportedModules.Contains(name.GetFullName()) {
					this_.priv_CollapseProperties_logDecisionForName_3240(name, "escapes - dynamically imported module namespace")
					escaped.Add(name)
					if !(name.HasProps()) {
						continue
					}
					{
						var i int32 = int32(0)
						_ = i
						var size int32 = name.GetPropCount()
						_ = size
						for ; i < size; i = i + 1 {
							var prop *GlobalNamespace_Name = name.GetPropAt(i)
							_ = prop
							var propDeclaration *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](prop.GetDeclaration())
							_ = propDeclaration
							if propDeclaration == nil {
								continue
							}
							if propDeclaration.GetNode() != 0 {
								var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(propDeclaration.GetNode())
								_ = rValue
								if rValue.IsName() {
									this_.priv_CollapseProperties_logDecisionForName_3240(name, "escapes - dynamically imported module namespace property alias")
									dynamicallyImportedModuleRefs.Add(rValue.GetQualifiedName())
								}
							}
						}
					}
				}
			}
		}
	}
	{
		cur_29 := this_.outer_.namespace.nameMap.IterValues()
		for cur_29.HasNext() {

			name_2 := cur_29.
				Next()
			_ = name_2
			if (dynamicallyImportedModuleRefs != nil) && dynamicallyImportedModuleRefs.Contains(name_2.GetFullName()) {
				escaped.Add(name_2)
			}
			if !(name_2.IsNamespaceObjectLit()) {
				continue
			}
			if ((name_2.GetAliasingGets() == int32(0)) && ((name_2.GetLocalSets() + name_2.GetGlobalSets()) <= int32(1))) && (name_2.GetDeleteProps() == int32(0)) {
				continue
			}
			var initialized bool = (j2gort.CastPtr[GlobalNamespace_Ref](name_2.GetDeclaration()) != nil)
			_ = initialized
			{
				cur_31 := j2gort.NewIterCursor(j2gort.CheckNotNull(name_2.GetRefs()))
				for cur_31.HasNext() {
					rawElem_32 := cur_31.Next()
					ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_32)
					_ = ref
					if ref.IsDeleteProp() {
						if initialized {
							this_.priv_CollapseProperties_warnAboutNamespaceRedefinition_3237(name_2, ref)
						}
					} else if ref.IsSet() && (ref != j2gort.CastPtr[GlobalNamespace_Ref](name_2.GetDeclaration())) {
						if initialized && !(InlineAndCollapseProperties_IsSafeNamespaceReinit(ref)) {
							this_.priv_CollapseProperties_warnAboutNamespaceRedefinition_3237(name_2, ref)
						}
						initialized = true
					} else if ref.IsAliasingGet() {
						this_.priv_CollapseProperties_warnAboutNamespaceAliasing_3236(name_2, ref)
						this_.priv_CollapseProperties_logDecisionForName_3240(name_2, "escapes")
						escaped.Add(name_2)
						break
					}
				}
			}
		}
	}
	return escaped
}

// Java method: CollapseProperties.warnAboutNamespaceAliasing(com.google.javascript.jscomp.GlobalNamespace.Name nameObj, com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_warnAboutNamespaceAliasing_3236(nameObj *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	_ = nameObj
	_ = ref
	this_.outer_.compiler.Report_1(JSError_Make_3(ref.GetNode(), StaticGet_InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING(), []string{nameObj.GetFullName()}))
}

// Java method: CollapseProperties.warnAboutNamespaceRedefinition(com.google.javascript.jscomp.GlobalNamespace.Name nameObj, com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_warnAboutNamespaceRedefinition_3237(nameObj *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	_ = nameObj
	_ = ref
	this_.outer_.compiler.Report_1(JSError_Make_3(ref.GetNode(), StaticGet_InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING(), []string{nameObj.GetFullName()}))
}

// Java method: CollapseProperties.flattenReferencesToCollapsibleDescendantNames(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias, java.util.LinkedHashSet<com.google.javascript.jscomp.GlobalNamespace.Name> escaped)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenReferencesToCollapsibleDescendantNames_3238(n *GlobalNamespace_Name, alias string, escaped *j2gort.OrderedSet[*GlobalNamespace_Name]) {
	_ = n
	_ = alias
	_ = escaped
	if !(n.HasProps()) {
		return
	}
	if n.IsCollapsingExplicitlyDenied() {
		this_.priv_CollapseProperties_logDecisionForName_3240(n, "@nocollapse: will not flatten descendant name references")
		return
	}
	if !(j2gort.CheckNotNullPtr(escaped).IsEmpty()) && j2gort.CheckNotNullPtr(escaped).Contains(n) {
		this_.priv_CollapseProperties_logDecisionForName_3240(n, "escapes: will not flatten descendant name references")
		return
	}
	var isLogging bool = ((this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging())
	_ = isLogging
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = n.GetPropCount()
		_ = size
		for ; i < size; i = i + 1 {
			var p *GlobalNamespace_Name = n.GetPropAt(i)
			_ = p
			var inlinability GlobalNamespace_Inlinability = p.CanCollapseOrInline()
			_ = inlinability
			var isAllowedToCollapse bool = ((this_.outer_.propertyCollapseLevel != StaticGet_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT()) || p.IsModuleExport())
			_ = isAllowedToCollapse
			var propAlias string = j2gort.NullString
			_ = propAlias
			if isAllowedToCollapse {
				if inlinability.CanCollapse() {
					this_.priv_CollapseProperties_logDecisionForName_3239(p, inlinability, "will flatten references")
					propAlias = this_.priv_CollapseProperties_appendPropForAlias_3262(alias, p.GetBaseName())
					this_.priv_CollapseProperties_flattenReferencesTo_3245(p, propAlias)
				} else if p.IsCollapsingExplicitlyDenied() {
					this_.priv_CollapseProperties_logDecisionForName_3240(p, "@nocollapse: will not flatten references")
				} else if p.IsSimpleStubDeclaration() {
					this_.priv_CollapseProperties_logDecisionForName_3240(p, "simple stub declaration: will flatten references")
					propAlias = this_.priv_CollapseProperties_appendPropForAlias_3262(alias, p.GetBaseName())
					this_.priv_CollapseProperties_flattenSimpleStubDeclaration_3242(p, propAlias)
				} else {
					this_.priv_CollapseProperties_logDecisionForName_3239(p, inlinability, "will not flatten references")
				}
			}
			if (p.HasProps() && !(p.IsCollapsingExplicitlyDenied())) && (j2gort.CheckNotNullPtr(escaped).IsEmpty() || !(j2gort.CheckNotNullPtr(escaped).Contains(p))) {
				if propAlias == j2gort.NullString {
					propAlias = this_.priv_CollapseProperties_appendPropForAlias_3262(alias, p.GetBaseName())
				}
				this_.priv_CollapseProperties_flattenReferencesToCollapsibleDescendantNames_3238(p, propAlias, escaped)
			} else if isLogging && p.HasProps() {
				if p.IsCollapsingExplicitlyDenied() {
					this_.priv_CollapseProperties_logDecisionForName_3240(p, "@nocollapse: will not flatten descendant name references")
				} else if !(j2gort.CheckNotNullPtr(escaped).IsEmpty()) && j2gort.CheckNotNullPtr(escaped).Contains(p) {
					this_.priv_CollapseProperties_logDecisionForName_3240(p, "escapes: will not flatten descendant name references")
				}
			}
		}
	}
}

// Java method: CollapseProperties.logDecisionForName(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.jscomp.GlobalNamespace.Inlinability inlinability, java.lang.String message)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_logDecisionForName_3239(name *GlobalNamespace_Name, inlinability GlobalNamespace_Inlinability, message string) {
	_ = name
	_ = inlinability
	_ = message
	if (this_.outer_.decisionsLog == nil || this_.outer_.decisionsLog.IsJavaNil_()) || !(this_.outer_.decisionsLog.IsLogging()) {
		return
	}
	this_.priv_CollapseProperties_logDecisionForName_3241(name, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("inlinability %s: %s", []any{j2gort.BoxNullableVal(inlinability), j2gort.BoxString(message)}))
	}))
}

// Java method: CollapseProperties.logDecisionForName(com.google.javascript.jscomp.GlobalNamespace.Name name, java.lang.String message)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_logDecisionForName_3240(name *GlobalNamespace_Name, message string) {
	_ = name
	_ = message
	if (this_.outer_.decisionsLog == nil || this_.outer_.decisionsLog.IsJavaNil_()) || !(this_.outer_.decisionsLog.IsLogging()) {
		return
	}
	if this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s: %s", []any{j2gort.BoxString(name.GetFullName()), j2gort.BoxString(message)}))
		}))
	}
}

// Java method: CollapseProperties.logDecisionForName(com.google.javascript.jscomp.GlobalNamespace.Name name, java.util.function.Supplier<java.lang.String> messageSupplier)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_logDecisionForName_3241(name *GlobalNamespace_Name, messageSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = name
	_ = messageSupplier
	if (this_.outer_.decisionsLog == nil || this_.outer_.decisionsLog.IsJavaNil_()) || !(this_.outer_.decisionsLog.IsLogging()) {
		return
	}
	if this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s: %s", []any{j2gort.BoxString(name.GetFullName()), j2gort.BoxString(j2gort.CastString(messageSupplier.Get()))}))
		}))
	}
}

// Java method: CollapseProperties.flattenSimpleStubDeclaration(com.google.javascript.jscomp.GlobalNamespace.Name name, java.lang.String alias)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenSimpleStubDeclaration_3242(name *GlobalNamespace_Name, alias string) {
	_ = name
	_ = alias
	var ref *GlobalNamespace_Ref = name.GetOnlyRef()
	_ = ref
	var nameNode pkg_rhino.Node = NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, alias, ref.GetNode(), name.GetFullName())
	_ = nameNode
	var varNode pkg_rhino.Node = pkg_rhino.IR_Var_1(nameNode).SrcrefIfMissing(nameNode)
	_ = varNode
	if !(ref.GetNode().GetParent().IsExprResult()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = ref.GetNode().GetParent()
	_ = parent
	parent.ReplaceWith(varNode)
	this_.outer_.compiler.ReportChangeToEnclosingScope(varNode)
}

// Java method: CollapseProperties.shouldAddConstantName(com.google.javascript.rhino.Node declarationNode, java.lang.String collapsedName)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_shouldAddConstantName_3243(declarationNode pkg_rhino.Node, collapsedName string) bool {
	_ = declarationNode
	_ = collapsedName
	if declarationNode == 0 {
		return false
	}
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(declarationNode.GetParent())
	_ = info
	return ((((info != nil) && info.IsConstant()) || declarationNode.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())) || this_.outer_.compiler.GetCodingConvention().IsConstant(collapsedName))
}

// Java method: CollapseProperties.flattenSingleDirectRef(com.google.javascript.jscomp.GlobalNamespace.Name n, com.google.javascript.jscomp.GlobalNamespace.Ref r, com.google.javascript.jscomp.GlobalNamespace.@org.jspecify.annotations.Nullable Ref decl, java.lang.String alias, java.lang.String originalName, boolean isConstantName)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenSingleDirectRef_3244(n *GlobalNamespace_Name, r *GlobalNamespace_Ref, decl *GlobalNamespace_Ref, alias string, originalName string, isConstantName bool) {
	_ = n
	_ = r
	_ = decl
	_ = alias
	_ = originalName
	_ = isConstantName
	if r == decl {
		return
	}
	var rNode pkg_rhino.Node = r.GetNode()
	_ = rNode
	var rParent pkg_rhino.Node = rNode.GetParent()
	_ = rParent
	if !(NodeUtil_MayBeObjectLitKey(rNode)) {
		this_.priv_CollapseProperties_flattenNameRef_3248(alias, rNode, rParent, originalName, isConstantName)
	} else if rNode.IsStringKey() && rParent.IsObjectPattern() {
		var newNode pkg_rhino.Node = pkg_rhino.IR_Name(alias).Srcref(rNode)
		_ = newNode
		NodeUtil_CopyNameAnnotations(rNode, newNode)
		DestructuringGlobalNameExtractor_ReassignDestructringLvalue(rNode, newNode, nil, r, this_.outer_.compiler)
	}
}

// Java method: CollapseProperties.flattenReferencesTo(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenReferencesTo_3245(n *GlobalNamespace_Name, alias string) {
	_ = n
	_ = alias
	var originalName string = n.GetFullName()
	_ = originalName
	var decl *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = decl
	var isConstantName bool = this_.priv_CollapseProperties_shouldAddConstantName_3243(func() pkg_rhino.Node {
		if decl != nil {
			return decl.GetNode()
		}
		return pkg_rhino.Node(0)
	}(), alias)
	_ = isConstantName
	if n.HasInlineRefs() {
		var r0 *GlobalNamespace_Ref = n.InlineRef0()
		_ = r0
		var r1 *GlobalNamespace_Ref = n.InlineRef1()
		_ = r1
		var r2 *GlobalNamespace_Ref = n.InlineRef2()
		_ = r2
		var r3 *GlobalNamespace_Ref = n.InlineRef3()
		_ = r3
		if r0 != nil {
			this_.priv_CollapseProperties_flattenSingleDirectRef_3244(n, r0, decl, alias, originalName, isConstantName)
		}
		if r1 != nil {
			this_.priv_CollapseProperties_flattenSingleDirectRef_3244(n, r1, decl, alias, originalName, isConstantName)
		}
		if r2 != nil {
			this_.priv_CollapseProperties_flattenSingleDirectRef_3244(n, r2, decl, alias, originalName, isConstantName)
		}
		if r3 != nil {
			this_.priv_CollapseProperties_flattenSingleDirectRef_3244(n, r3, decl, alias, originalName, isConstantName)
		}
	} else {
		{
			cur_33 := j2gort.NewIterCursor(j2gort.CheckNotNull(n.GetRefs()))
			for cur_33.HasNext() {
				rawElem_34 := cur_33.Next()
				r := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_34)
				_ = r
				this_.priv_CollapseProperties_flattenSingleDirectRef_3244(n, r, decl, alias, originalName, isConstantName)
			}
		}
	}
	if n.HasProps() {
		var prefixBuf *j2gort.Ext_java_lang_StringBuilder = this_.prefixOriginalNameBuf
		_ = prefixBuf
		prefixBuf.SetLength(int32(0))
		prefixBuf.Append_1_string(originalName)
		var baseLen int32 = prefixBuf.Length()
		_ = baseLen
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = n.GetPropCount()
			_ = size
			for ; i < size; i = i + 1 {
				var p *GlobalNamespace_Name = n.GetPropAt(i)
				_ = p
				prefixBuf.Append_1_string(p.GetBaseName())
				this_.priv_CollapseProperties_flattenPrefixes_3246(alias, prefixBuf, p, int32(1), isConstantName)
				prefixBuf.SetLength(baseLen)
			}
		}
	}
}

// Java method: CollapseProperties.flattenPrefixes(java.lang.String alias, java.lang.StringBuilder prefixBuf, com.google.javascript.jscomp.GlobalNamespace.Name n, int depth, boolean isConstantName)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenPrefixes_3246(alias string, prefixBuf *j2gort.Ext_java_lang_StringBuilder, n *GlobalNamespace_Name, depth int32, isConstantName bool) {
	_ = alias
	_ = prefixBuf
	_ = n
	_ = depth
	_ = isConstantName
	var materializedOriginalName string = j2gort.NullString
	_ = materializedOriginalName
	var decl *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = decl
	if ((decl != nil) && (decl.GetNode() != 0)) && decl.GetNode().IsGetProp() {
		materializedOriginalName = this_.priv_CollapseProperties_flattenNameRefAtDepth_3247(alias, decl.GetNode(), depth, prefixBuf, materializedOriginalName, isConstantName)
	}
	if n.HasInlineRefs() {
		var r0 *GlobalNamespace_Ref = n.InlineRef0()
		_ = r0
		if r0 != nil {
			if r0 != decl {
				materializedOriginalName = this_.priv_CollapseProperties_flattenNameRefAtDepth_3247(alias, r0.GetNode(), depth, prefixBuf, materializedOriginalName, isConstantName)
			}
			var r1 *GlobalNamespace_Ref = n.InlineRef1()
			_ = r1
			if r1 != nil {
				if r1 != decl {
					materializedOriginalName = this_.priv_CollapseProperties_flattenNameRefAtDepth_3247(alias, r1.GetNode(), depth, prefixBuf, materializedOriginalName, isConstantName)
				}
				var r2 *GlobalNamespace_Ref = n.InlineRef2()
				_ = r2
				if r2 != nil {
					if r2 != decl {
						materializedOriginalName = this_.priv_CollapseProperties_flattenNameRefAtDepth_3247(alias, r2.GetNode(), depth, prefixBuf, materializedOriginalName, isConstantName)
					}
					var r3 *GlobalNamespace_Ref = n.InlineRef3()
					_ = r3
					if (r3 != nil) && (r3 != decl) {
						materializedOriginalName = this_.priv_CollapseProperties_flattenNameRefAtDepth_3247(alias, r3.GetNode(), depth, prefixBuf, materializedOriginalName, isConstantName)
					}
				}
			}
		}
	} else {
		{
			cur_35 := j2gort.NewIterCursor(j2gort.CheckNotNull(n.GetRefs()))
			for cur_35.HasNext() {
				rawElem_36 := cur_35.Next()
				r := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_36)
				_ = r
				if r == decl {
					continue
				}
				materializedOriginalName = this_.priv_CollapseProperties_flattenNameRefAtDepth_3247(alias, r.GetNode(), depth, prefixBuf, materializedOriginalName, isConstantName)
			}
		}
	}
	if n.HasProps() {
		var prevLen int32 = prefixBuf.Length()
		_ = prevLen
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = n.GetPropCount()
			_ = size
			for ; i < size; i = i + 1 {
				var p *GlobalNamespace_Name = n.GetPropAt(i)
				_ = p
				prefixBuf.Append_1_string(p.GetBaseName())
				this_.priv_CollapseProperties_flattenPrefixes_3246(alias, prefixBuf, p, (depth + int32(1)), isConstantName)
				prefixBuf.SetLength(prevLen)
			}
		}
	}
}

// Java method: CollapseProperties.flattenNameRefAtDepth(java.lang.String alias, com.google.javascript.rhino.Node n, int depth, java.lang.StringBuilder prefixBuf, java.lang.@org.jspecify.annotations.Nullable String materializedOriginalName, boolean isConstantName)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenNameRefAtDepth_3247(alias string, n pkg_rhino.Node, depth int32, prefixBuf *j2gort.Ext_java_lang_StringBuilder, materializedOriginalName string, isConstantName bool) string {
	_ = alias
	_ = n
	_ = depth
	_ = prefixBuf
	_ = materializedOriginalName
	_ = isConstantName
	var nType pkg_rhino.Token = n.GetToken()
	_ = nType
	var isQName bool = ((nType == pkg_rhino.StaticGet_Token_NAME()) || (nType == pkg_rhino.StaticGet_Token_GETPROP()))
	_ = isQName
	var isObjKey bool = NodeUtil_MayBeObjectLitKey(n)
	_ = isObjKey
	if !(isObjKey || isQName) {
		j2gort.CheckState(false)
	}
	if isQName {
		{
			var i int32 = int32(1)
			_ = i
			for ; (i < depth) && n.HasChildren(); i = i + 1 {
				n = n.GetFirstChild()
			}
		}
		if n.IsGetProp() && n.GetFirstChild().IsGetProp() {
			if materializedOriginalName == j2gort.NullString {
				materializedOriginalName = prefixBuf.String()
			}
			this_.priv_CollapseProperties_flattenNameRef_3248(alias, n.GetFirstChild(), n, materializedOriginalName, isConstantName)
		}
	}
	return materializedOriginalName
}

// Java method: CollapseProperties.flattenNameRef(java.lang.String alias, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent, java.lang.String originalName, boolean isConstantName)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_flattenNameRef_3248(alias string, n pkg_rhino.Node, parent pkg_rhino.Node, originalName string, isConstantName bool) {
	_ = alias
	_ = n
	_ = parent
	_ = originalName
	_ = isConstantName
	if !(n.IsGetProp()) {
		j2gort.CheckArgumentFmt(false, "Expected GETPROP, found %s. Node: %s", j2gort.BoxNullableVal(n.GetToken()), j2gort.BoxNode(n))
	}
	var ref pkg_rhino.Node = NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, alias, n, originalName).CopyTypeFrom(n)
	_ = ref
	NodeUtil_CopyNameAnnotations(n, ref)
	if isConstantName {
		ref.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), isConstantName)
	}
	if NodeUtil_IsNormalOrOptChainCall(parent) && n.IsFirstChildOf(parent) {
		parent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	}
	n.ReplaceWith(ref)
	this_.outer_.compiler.ReportChangeToEnclosingScope(ref)
}

// Java method: CollapseProperties.collapseDeclarationOfNameAndDescendants(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias, java.util.LinkedHashSet<com.google.javascript.jscomp.GlobalNamespace.Name> escaped)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_collapseDeclarationOfNameAndDescendants_3249(n *GlobalNamespace_Name, alias string, escaped *j2gort.OrderedSet[*GlobalNamespace_Name]) {
	_ = n
	_ = alias
	_ = escaped
	var childNameInlinability GlobalNamespace_Inlinability = n.CanCollapseOrInlineChildNames()
	_ = childNameInlinability
	var canCollapseChildNames bool = false
	_ = canCollapseChildNames
	if !(childNameInlinability.CanCollapse()) {
		if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
			this_.priv_CollapseProperties_logDecisionForName_3241(n, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("child name inlinability: %s: will not collapse child names", []any{j2gort.BoxNullableVal(childNameInlinability)}))
			}))
		}
		canCollapseChildNames = false
	} else if !(j2gort.CheckNotNullPtr(escaped).IsEmpty()) && j2gort.CheckNotNullPtr(escaped).Contains(n) {
		this_.priv_CollapseProperties_logDecisionForName_3240(n, "escapes: will not collapse child names")
		canCollapseChildNames = false
	} else {
		canCollapseChildNames = true
	}
	if this_.priv_CollapseProperties_canCollapse_3233(n) {
		this_.priv_CollapseProperties_logDecisionForName_3240(n, "collapsing")
		this_.priv_CollapseProperties_updateGlobalNameDeclaration_3251(n, alias, canCollapseChildNames)
	}
	if !(n.HasProps()) {
		return
	}
	if !(j2gort.CheckNotNullPtr(escaped).IsEmpty()) && j2gort.CheckNotNullPtr(escaped).Contains(n) {
		return
	}
	this_.priv_CollapseProperties_logDecisionForName_3240(n, "collapsing descendants")
	var isLogging bool = ((this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging())
	_ = isLogging
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = n.GetPropCount()
		_ = size
		for ; i < size; i = i + 1 {
			var p *GlobalNamespace_Name = n.GetPropAt(i)
			_ = p
			if (!(isLogging) && !(p.HasProps())) && !(this_.priv_CollapseProperties_canCollapse_3233(p)) {
				continue
			}
			this_.priv_CollapseProperties_collapseDeclarationOfNameAndDescendants_3249(p, this_.priv_CollapseProperties_appendPropForAlias_3262(alias, p.GetBaseName()), escaped)
		}
	}
}

// Java method: CollapseProperties.updateTwinnedDeclaration(java.lang.String alias, com.google.javascript.jscomp.GlobalNamespace.Name refName, com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateTwinnedDeclaration_3250(alias string, refName *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	_ = alias
	_ = refName
	_ = ref
	if !(ref.IsTwin()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(ref))
	}
	if !(ref.GetNode().IsGetProp()) {
		return
	}
	var rvalue pkg_rhino.Node = ref.GetNode().GetNext()
	_ = rvalue
	var parent pkg_rhino.Node = ref.GetNode().GetParent()
	_ = parent
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if (rvalue != 0) && rvalue.IsFunction() {
		this_.priv_CollapseProperties_checkForReceiverAffectedByCollapse_3255(rvalue, refName.GetJSDocInfo(), refName)
	}
	var nameNode pkg_rhino.Node = NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, alias, grandparent.GetFirstChild(), refName.GetFullName())
	_ = nameNode
	NodeUtil_CopyNameAnnotations(ref.GetNode(), nameNode)
	var current pkg_rhino.Node = grandparent
	_ = current
	var currentParent pkg_rhino.Node = grandparent.GetParent()
	_ = currentParent
	for ; !(currentParent.IsScript()) && !(currentParent.IsBlock()); func() { current = currentParent; currentParent = currentParent.GetParent() }() {
	}
	var stubVar pkg_rhino.Node = pkg_rhino.IR_Var_1(nameNode.CloneTree_0()).SrcrefIfMissing(nameNode)
	_ = stubVar
	stubVar.InsertBefore(current)
	ref.GetNode().ReplaceWith(nameNode)
	this_.outer_.compiler.ReportChangeToEnclosingScope(nameNode)
}

// Java method: CollapseProperties.updateGlobalNameDeclaration(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias, boolean canCollapseChildNames)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateGlobalNameDeclaration_3251(n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	_ = n
	_ = alias
	_ = canCollapseChildNames
	var decl *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = decl
	if decl == nil {
		this_.priv_CollapseProperties_logDecisionForName_3240(n, "no global declaration found")
		return
	}
	var isSimpleName bool = n.IsSimpleName()
	_ = isSimpleName
	var declNode pkg_rhino.Node = decl.GetNode()
	_ = declNode
	switch declNode.GetParent().TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_CollapseProperties_logDeclarationAction_3252(n, declNode, "updating assignment")
		this_.priv_CollapseProperties_updateGlobalNameDeclarationAtAssignNode_3254(n, alias, canCollapseChildNames)
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.priv_CollapseProperties_logDeclarationAction_3252(n, declNode, "updating variable declaration")
		this_.priv_CollapseProperties_updateGlobalNameDeclarationAtVariableNode_3256(n, canCollapseChildNames)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_CollapseProperties_logDeclarationAction_3252(n, declNode, "updating function declaration")
		this_.priv_CollapseProperties_updateGlobalNameDeclarationAtFunctionNode_3257(n, canCollapseChildNames)
	case int32(105) /* Token.CLASS */ :
		this_.priv_CollapseProperties_logDeclarationAction_3252(n, declNode, "updating class declaration")
		this_.priv_CollapseProperties_updateGlobalNameDeclarationAtClassNode_3258(n, canCollapseChildNames)
	case int32(106) /* Token.CLASS_MEMBERS */ :
		this_.priv_CollapseProperties_logDeclarationAction_3252(n, declNode, "updating static member declaration")
		this_.priv_CollapseProperties_updateGlobalNameDeclarationAtStaticMemberNode_3259(n, alias, canCollapseChildNames)
	default:
		this_.priv_CollapseProperties_logDeclarationAction_3252(n, declNode, "not updating an unsupported type of declaration node")
	}
	if !(isSimpleName) {
		var commonAncestor *JSChunk = n.GetDeepestCommonAncestorChunk(this_.outer_.compiler.GetChunkGraph())
		_ = commonAncestor
		if commonAncestor != decl.GetChunk() {
			this_.priv_CollapseProperties_moveDeclarationIntoHigherChunk_3253(decl, commonAncestor)
		}
	}
}

// Java method: CollapseProperties.logDeclarationAction(com.google.javascript.jscomp.GlobalNamespace.Name name, com.google.javascript.rhino.Node declarationNode, java.lang.String message)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_logDeclarationAction_3252(name *GlobalNamespace_Name, declarationNode pkg_rhino.Node, message string) {
	_ = name
	_ = declarationNode
	_ = message
	if (this_.outer_.decisionsLog == nil || this_.outer_.decisionsLog.IsJavaNil_()) || !(this_.outer_.decisionsLog.IsLogging()) {
		return
	}
	this_.priv_CollapseProperties_logDecisionForName_3241(name, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s: %s", []any{j2gort.BoxNode(declarationNode), j2gort.BoxString(message)}))
	}))
}

// Java method: CollapseProperties.moveDeclarationIntoHigherChunk(com.google.javascript.jscomp.GlobalNamespace.Ref decl, com.google.javascript.jscomp.JSChunk commonAncestor)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_moveDeclarationIntoHigherChunk_3253(decl *GlobalNamespace_Ref, commonAncestor *JSChunk) {
	_ = decl
	_ = commonAncestor
	var newNode pkg_rhino.Node = decl.GetNode()
	_ = newNode
	var varNode pkg_rhino.Node = newNode.GetParent()
	_ = varNode
	var initializer pkg_rhino.Node = newNode.GetFirstChild()
	_ = initializer
	if !(varNode.IsVar()) {
		j2gort.CheckStateFmt(false, "Unexpected declaration %s", j2gort.BoxNode(varNode))
	}
	var newRoot pkg_rhino.Node = j2gort.CastPtr[CompilerInput](commonAncestor.GetInputs().Get(int32(0))).GetAstRoot(this_.outer_.compiler)
	_ = newRoot
	if initializer == 0 {
		newRoot.AddChildToFront(varNode.Detach())
	} else {
		var newVarNode pkg_rhino.Node = pkg_rhino.IR_Var_1(newNode.CloneNode_0()).Srcref(newNode)
		_ = newVarNode
		newRoot.AddChildToFront(newVarNode)
		var newAssign pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(newNode.Detach(), initializer.Detach())).SrcrefTreeIfMissing(newNode)
		_ = newAssign
		varNode.ReplaceWith(newAssign)
	}
	this_.outer_.compiler.ReportChangeToChangeScope(newRoot)
}

// Java method: CollapseProperties.updateGlobalNameDeclarationAtAssignNode(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias, boolean canCollapseChildNames)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateGlobalNameDeclarationAtAssignNode_3254(n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	_ = n
	_ = alias
	_ = canCollapseChildNames
	var ref *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = ref
	var rvalue pkg_rhino.Node = ref.GetNode().GetNext()
	_ = rvalue
	if ref.IsTwin() {
		this_.priv_CollapseProperties_updateTwinnedDeclaration_3250(alias, n, ref)
		return
	}
	var varNode pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_VAR())
	_ = varNode
	var varParent pkg_rhino.Node = ref.GetNode().GetAncestor(int32(3))
	_ = varParent
	var grandparent pkg_rhino.Node = ref.GetNode().GetAncestor(int32(2))
	_ = grandparent
	var isObjLit bool = rvalue.IsObjectLit()
	_ = isObjLit
	var insertedVarNode bool = false
	_ = insertedVarNode
	if isObjLit && this_.priv_CollapseProperties_canEliminate_3234(n, rvalue) {
		grandparent.ReplaceWith(varNode)
		n.UpdateRefNode(ref, pkg_rhino.Node(0))
		insertedVarNode = true
		this_.outer_.compiler.ReportChangeToEnclosingScope(varNode)
	} else if !(n.IsSimpleName()) {
		if rvalue.IsFunction() {
			this_.priv_CollapseProperties_checkForReceiverAffectedByCollapse_3255(rvalue, n.GetJSDocInfo(), n)
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(rvalue)
		rvalue.Detach()
		var nameNode pkg_rhino.Node = NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, alias, ref.GetNode().GetAncestor(int32(2)), n.GetFullName())
		_ = nameNode
		nameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), this_.priv_CollapseProperties_shouldAddConstantName_3243(ref.GetNode(), alias))
		var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(ref.GetNode().GetParent())
		_ = info
		if info != nil {
			varNode.SetJSDocInfo(info)
		}
		varNode.AddChildToBack(nameNode)
		nameNode.AddChildToFront(rvalue)
		grandparent.ReplaceWith(varNode)
		n.UpdateRefNode(ref, nameNode)
		insertedVarNode = true
		this_.outer_.compiler.ReportChangeToEnclosingScope(varNode)
	}
	if canCollapseChildNames {
		if isObjLit {
			this_.priv_CollapseProperties_declareVariablesForObjLitValues_3260(n, alias, rvalue, varNode, varNode.GetPrevious())
		}
		this_.priv_CollapseProperties_addStubsForUndeclaredProperties_3261(n, alias, varParent, varNode)
	}
	if insertedVarNode {
		if !(varNode.HasChildren()) {
			varNode.Detach()
		}
	}
}

// Java method: CollapseProperties.checkForReceiverAffectedByCollapse(com.google.javascript.rhino.Node function, com.google.javascript.rhino.JSDocInfo docInfo, com.google.javascript.jscomp.GlobalNamespace.Name name)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_checkForReceiverAffectedByCollapse_3255(function pkg_rhino.Node, docInfo *pkg_rhino.JSDocInfo, name *GlobalNamespace_Name) {
	_ = function
	_ = docInfo
	_ = name
	if !(function.IsFunction()) {
		j2gort.CheckState(false)
	}
	if docInfo != nil {
		if docInfo.IsConstructorOrInterface() {
			return
		}
		if docInfo.HasThisType() {
			return
		}
	}
	if NodeUtil_ReferencesOwnReceiver(function) {
		this_.outer_.compiler.Report_1(JSError_Make_3(function, StaticGet_InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE(), []string{name.GetFullName()}))
	}
}

// Java method: CollapseProperties.updateGlobalNameDeclarationAtVariableNode(com.google.javascript.jscomp.GlobalNamespace.Name n, boolean canCollapseChildNames)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateGlobalNameDeclarationAtVariableNode_3256(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	_ = n
	_ = canCollapseChildNames
	if !(canCollapseChildNames) {
		this_.priv_CollapseProperties_logDecisionForName_3240(n, "cannot collapse child names: skipping")
		return
	}
	var ref *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = ref
	var name string = ref.GetNode().GetString()
	_ = name
	var rvalue pkg_rhino.Node = ref.GetNode().GetFirstChild()
	_ = rvalue
	var variableNode pkg_rhino.Node = ref.GetNode().GetParent()
	_ = variableNode
	var grandparent pkg_rhino.Node = variableNode.GetParent()
	_ = grandparent
	var isObjLit bool = rvalue.IsObjectLit()
	_ = isObjLit
	if isObjLit {
		this_.priv_CollapseProperties_declareVariablesForObjLitValues_3260(n, name, rvalue, variableNode, variableNode.GetPrevious())
	}
	this_.priv_CollapseProperties_addStubsForUndeclaredProperties_3261(n, name, grandparent, variableNode)
	if isObjLit && this_.priv_CollapseProperties_canEliminate_3234(n, rvalue) {
		ref.GetNode().Detach()
		this_.outer_.compiler.ReportChangeToEnclosingScope(variableNode)
		if !(variableNode.HasChildren()) {
			variableNode.Detach()
		}
		n.UpdateRefNode(ref, pkg_rhino.Node(0))
	}
}

// Java method: CollapseProperties.updateGlobalNameDeclarationAtFunctionNode(com.google.javascript.jscomp.GlobalNamespace.Name n, boolean canCollapseChildNames)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateGlobalNameDeclarationAtFunctionNode_3257(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	_ = n
	_ = canCollapseChildNames
	if !(canCollapseChildNames) || !(this_.priv_CollapseProperties_canCollapse_3233(n)) {
		return
	}
	var ref *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = ref
	var fnName string = ref.GetNode().GetString()
	_ = fnName
	this_.priv_CollapseProperties_addStubsForUndeclaredProperties_3261(n, fnName, ref.GetNode().GetAncestor(int32(2)), ref.GetNode().GetParent())
}

// Java method: CollapseProperties.updateGlobalNameDeclarationAtClassNode(com.google.javascript.jscomp.GlobalNamespace.Name n, boolean canCollapseChildNames)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateGlobalNameDeclarationAtClassNode_3258(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	_ = n
	_ = canCollapseChildNames
	if !(canCollapseChildNames) || !(this_.priv_CollapseProperties_canCollapse_3233(n)) {
		return
	}
	var ref *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = ref
	var className string = ref.GetNode().GetString()
	_ = className
	this_.priv_CollapseProperties_addStubsForUndeclaredProperties_3261(n, className, ref.GetNode().GetAncestor(int32(2)), ref.GetNode().GetParent())
}

// Java method: CollapseProperties.updateGlobalNameDeclarationAtStaticMemberNode(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias, boolean canCollapseChildNames)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_updateGlobalNameDeclarationAtStaticMemberNode_3259(n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	_ = n
	_ = alias
	_ = canCollapseChildNames
	var declaration *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](n.GetDeclaration())
	_ = declaration
	var classNode pkg_rhino.Node = declaration.GetNode().GetGrandparent()
	_ = classNode
	if !(classNode.IsClass()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classNode))
	}
	var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(classNode)
	_ = enclosingStatement
	if canCollapseChildNames {
		this_.priv_CollapseProperties_addStubsForUndeclaredProperties_3261(n, alias, enclosingStatement.GetParent(), classNode)
	}
	var member pkg_rhino.Node = declaration.GetNode().Detach()
	_ = member
	var varDecl pkg_rhino.Node = pkg_rhino.Node(0)
	_ = varDecl
	if member.HasChildren() {
		var memberVal pkg_rhino.Node = member.GetOnlyChild()
		_ = memberVal
		this_.priv_CollapseProperties_checkForReceiverAffectedByCollapse_3255(memberVal, member.GetJSDocInfo(), n)
		memberVal = memberVal.Detach()
		varDecl = pkg_rhino.IR_Var_2(NodeUtil_NewName_3_AbstractCompiler_string_Node(this_.outer_.compiler, alias, member), memberVal).Srcref(member)
	} else {
		varDecl = pkg_rhino.IR_Var_1(NodeUtil_NewName_3_AbstractCompiler_string_Node(this_.outer_.compiler, alias, member)).Srcref(member)
	}
	varDecl.InsertBefore(enclosingStatement)
	varDecl.SetJSDocInfo(member.GetJSDocInfo())
	this_.outer_.compiler.ReportChangeToEnclosingScope(varDecl)
	n.UpdateRefNode(declaration, varDecl.GetFirstChild())
}

// Java method: CollapseProperties.declareVariablesForObjLitValues(com.google.javascript.jscomp.GlobalNamespace.Name objlitName, java.lang.String alias, com.google.javascript.rhino.Node objlit, com.google.javascript.rhino.Node varNode, com.google.javascript.rhino.Node nameToAddAfter)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_declareVariablesForObjLitValues_3260(objlitName *GlobalNamespace_Name, alias string, objlit pkg_rhino.Node, varNode pkg_rhino.Node, nameToAddAfter pkg_rhino.Node) {
	_ = objlitName
	_ = alias
	_ = objlit
	_ = varNode
	_ = nameToAddAfter
	var arbitraryNameCounter int32 = int32(0)
	_ = arbitraryNameCounter
	var discardKeys bool = !(objlitName.ShouldKeepKeys())
	_ = discardKeys
	{
		var key pkg_rhino.Node = objlit.GetFirstChild()
		_ = key
		var nextKey pkg_rhino.Node = pkg_rhino.Node(0)
		_ = nextKey
		for ; key != 0; key = nextKey {
			var value pkg_rhino.Node = key.GetFirstChild()
			_ = value
			nextKey = key.GetNext()
			switch key.TokenSwitchOrd_() {
			case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(128) /* Token.OBJECT_SPREAD */ :
				continue
			case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected child of OBJECTLIT: "+j2gort.StrValString(key.ToStringTree()))))
			}
			var isJsIdentifier bool = (!(key.IsNumber()) && pkg_rhino.TokenStream_IsJSIdentifier(key.GetString()))
			_ = isJsIdentifier
			var propName string
			_ = propName
			if isJsIdentifier {
				propName = key.GetString()
			} else {
				propName = strconv.FormatInt(int64(func() int32 {
					old_37 := arbitraryNameCounter
					nw_38 := old_37 + 1
					arbitraryNameCounter = nw_38
					return nw_38
				}()), 10)
			}
			var qName string = (j2gort.StrValString(objlitName.GetFullName()) + "." + j2gort.StrValString(propName))
			_ = qName
			var p *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](this_.outer_.namespace.GetOwnSlot(qName))
			_ = p
			if (p != nil) && !(this_.priv_CollapseProperties_canCollapse_3233(p)) {
				continue
			}
			var propAlias string = this_.priv_CollapseProperties_appendPropForAlias_3262(alias, propName)
			_ = propAlias
			var refNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = refNode
			var isConstantName bool = this_.priv_CollapseProperties_shouldAddConstantName_3243(key, propAlias)
			_ = isConstantName
			if discardKeys {
				key.Detach()
				value.Detach()
			} else {
				refNode = pkg_rhino.IR_Name(propAlias)
				if isConstantName {
					refNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
				}
				value.ReplaceWith(refNode)
				this_.outer_.compiler.ReportChangeToEnclosingScope(refNode)
			}
			var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(propAlias)
			_ = nameNode
			nameNode.AddChildToFront(value)
			if isConstantName {
				nameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
			}
			var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(nameNode).SrcrefTreeIfMissing(key)
			_ = newVar
			if nameToAddAfter != 0 {
				newVar.InsertAfter(nameToAddAfter)
			} else {
				newVar.InsertBefore(varNode)
			}
			this_.outer_.compiler.ReportChangeToEnclosingScope(newVar)
			nameToAddAfter = newVar
			if isJsIdentifier && (p != nil) {
				if !(discardKeys) {
					p.AddAliasingGetClonedFromDeclaration(refNode)
				}
				p.UpdateRefNode(j2gort.CastPtr[GlobalNamespace_Ref](p.GetDeclaration()), nameNode)
				if value.IsFunction() {
					this_.priv_CollapseProperties_checkForReceiverAffectedByCollapse_3255(value, key.GetJSDocInfo(), p)
				}
			}
		}
	}
}

// Java method: CollapseProperties.addStubsForUndeclaredProperties(com.google.javascript.jscomp.GlobalNamespace.Name n, java.lang.String alias, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node addAfter)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_addStubsForUndeclaredProperties_3261(n *GlobalNamespace_Name, alias string, parent pkg_rhino.Node, addAfter pkg_rhino.Node) {
	_ = n
	_ = alias
	_ = parent
	_ = addAfter
	if !(n.CanCollapseUnannotatedChildNames()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(n))
	}
	if !(NodeUtil_IsStatementBlock(parent)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(parent))
	}
	j2gort.CheckNotNullNode(addAfter)
	if !(n.HasProps()) {
		return
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = n.GetPropCount()
		_ = size
		for ; i < size; i = i + 1 {
			var p *GlobalNamespace_Name = n.GetPropAt(i)
			_ = p
			if !(p.NeedsToBeStubbed()) {
				continue
			}
			var propAlias string = this_.priv_CollapseProperties_appendPropForAlias_3262(alias, p.GetBaseName())
			_ = propAlias
			var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(propAlias)
			_ = nameNode
			var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(nameNode).SrcrefTreeIfMissing(addAfter)
			_ = newVar
			newVar.InsertAfter(addAfter)
			var constPropNode pkg_rhino.Node = p.GetFirstRef().GetNode()
			_ = constPropNode
			var isConstantName bool = (constPropNode.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) || this_.outer_.compiler.GetCodingConvention().IsConstant(propAlias))
			_ = isConstantName
			if isConstantName {
				nameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
			}
			this_.outer_.compiler.ReportChangeToEnclosingScope(newVar)
			addAfter = newVar
		}
	}
}

// Java method: CollapseProperties.appendPropForAlias(java.lang.String root, java.lang.String prop)
func (this_ *InlineAndCollapseProperties_CollapseProperties) priv_CollapseProperties_appendPropForAlias_3262(root string, prop string) string {
	_ = root
	_ = prop
	if j2gort.JStrIndexOfCp(prop, int32(uint16(36))) != int32(-1) {
		prop = javalang.StringReplace(prop, "$", "$0")
	}
	var result string = (j2gort.StrValString(root) + "$" + j2gort.StrValString(prop))
	_ = result
	var id int32 = int32(1)
	_ = id
	for this_.outer_.namespace.ContainsName(result) {
		result = (j2gort.StrValString(root) + "$" + j2gort.StrValString(prop) + "$" + j2gort.StrValInt32(id))
		id = id + 1
	}
	return result
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.CollapseProperties.FindDynamicallyImportedModules
type InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                 *InlineAndCollapseProperties_CollapseProperties
	processCommonJSModules bool                        `gson:"processCommonJSModules"`
	moduleResolutionMode   ModuleLoader_ResolutionMode `gson:"moduleResolutionMode"`
}

func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) IsJavaNil_() bool {
	return this_ == nil
}
func (*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$CollapseProperties$FindDynamicallyImportedModules", "FindDynamicallyImportedModules"
}
func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.CollapseProperties.FindDynamicallyImportedModules", "com.google.javascript.jscomp.InlineAndCollapseProperties$CollapseProperties$FindDynamicallyImportedModules", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules](v_)
	})
}
func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) SetOuter_(o *InlineAndCollapseProperties_CollapseProperties) {
	this_.outer_ = o
}
func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) GetOuter_() *InlineAndCollapseProperties_CollapseProperties {
	return this_.outer_
}
func (*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) initFields_InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules() {
	_ = this_
}

// Java constructor: FindDynamicallyImportedModules(boolean processCommonJSModules, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode resolutionMode)
func NewInlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules(outer_ *InlineAndCollapseProperties_CollapseProperties, processCommonJSModules bool, resolutionMode ModuleLoader_ResolutionMode) *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules {
	this_ := &InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules{}
	this_.outer_ = outer_
	InitInlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules(this_, processCommonJSModules, resolutionMode)
	return this_
}

func InitInlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules(this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules, processCommonJSModules bool, resolutionMode ModuleLoader_ResolutionMode) {
	_ = processCommonJSModules
	_ = resolutionMode
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules()
	this_.processCommonJSModules = processCommonJSModules
	this_.moduleResolutionMode = resolutionMode
}

// Java method: FindDynamicallyImportedModules.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if ((((((this_.processCommonJSModules && n.IsGetProp()) && n.IsQualifiedName()) && (n.GetParent() != 0)) && n.GetParent().IsReturn()) && n.GetGrandparent().IsBlock()) && n.GetGrandparent().HasOneChild()) && n.GetGrandparent().GetParent().IsFunction() {
		var potentialCallback pkg_rhino.Node = NodeUtil_GetEnclosingFunction(n)
		_ = potentialCallback
		if (potentialCallback != 0) && ProcessCommonJSModules_IsCommonJsDynamicImportCallback(NodeUtil_GetEnclosingFunction(potentialCallback), this_.moduleResolutionMode) {
			this_.outer_.dynamicallyImportedModules.Add(NodeUtil_GetRootOfQualifiedName_1_string(n.GetQualifiedName()))
		}
	} else if ConvertChunksToESModules_IsDynamicImportCallback(n) {
		var moduleNamespace pkg_rhino.Node = ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(this_.outer_.outer_.compiler, n)
		_ = moduleNamespace
		if moduleNamespace != 0 {
			this_.outer_.dynamicallyImportedModules.Add(moduleNamespace.GetQualifiedName())
		}
	}
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.InlineAliases
type InlineAndCollapseProperties_InlineAliases struct {
	outer_     *InlineAndCollapseProperties
	aliases    j2gort.OrderedMap[string, string] `gson:"aliases"`
	namespace  *GlobalNamespace                  `gson:"namespace"`
	astFactory *AstFactory                       `gson:"astFactory"`
	idHash_    atomic.Int32
}

func (this_ *InlineAndCollapseProperties_InlineAliases) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_InlineAliases) IsJavaNil_() bool { return this_ == nil }
func (*InlineAndCollapseProperties_InlineAliases) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$InlineAliases", "InlineAliases"
}
func (this_ *InlineAndCollapseProperties_InlineAliases) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineAndCollapseProperties_InlineAliases) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.InlineAliases", "com.google.javascript.jscomp.InlineAndCollapseProperties$InlineAliases", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_InlineAliases](v_)
	})
}
func (this_ *InlineAndCollapseProperties_InlineAliases) SetOuter_(o *InlineAndCollapseProperties) {
	this_.outer_ = o
}
func (this_ *InlineAndCollapseProperties_InlineAliases) GetOuter_() *InlineAndCollapseProperties {
	return this_.outer_
}
func (*InlineAndCollapseProperties_InlineAliases) Implements_jscomp_CompilerPass() {}

func (this_ *InlineAndCollapseProperties_InlineAliases) initFields_InlineAndCollapseProperties_InlineAliases() {
	_ = this_
	this_.aliases = j2gort.NewOrderedMap[string, string]()
}

// Java constructor: InlineAliases()
func NewInlineAndCollapseProperties_InlineAliases(outer_ *InlineAndCollapseProperties) *InlineAndCollapseProperties_InlineAliases {
	this_ := &InlineAndCollapseProperties_InlineAliases{}
	this_.outer_ = outer_
	InitInlineAndCollapseProperties_InlineAliases(this_)
	return this_
}

func InitInlineAndCollapseProperties_InlineAliases(this_ *InlineAndCollapseProperties_InlineAliases) {
	this_.initFields_InlineAndCollapseProperties_InlineAliases()
	this_.astFactory = this_.outer_.compiler.CreateAstFactory()
}

// Java method: InlineAliases.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineAndCollapseProperties_InlineAliases) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer func() {
		if this_.namespace != nil {
			this_.namespace.ReleaseSlabs_()
			this_.namespace = nil
		}
	}()

	this_.namespace = NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_.outer_.compiler, externs, root)
	NodeTraversal_TraverseRoots(this_.outer_.compiler, NewInlineAndCollapseProperties_InlineAliases_AliasesCollector(this_), externs, root)
	NodeTraversal_TraverseRoots(this_.outer_.compiler, NewInlineAndCollapseProperties_InlineAliases_AliasesInliner(this_), externs, root)
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.InlineAliases.AliasesCollector
type InlineAndCollapseProperties_InlineAliases_AliasesCollector struct {
	NodeTraversal_ExternsSkippingCallbackBase
	outer_ *InlineAndCollapseProperties_InlineAliases
}

func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) IsJavaNil_() bool {
	return this_ == nil
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$InlineAliases$AliasesCollector", "AliasesCollector"
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_ExternsSkippingCallbackBase.IdentityHash_()
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.InlineAliases.AliasesCollector", "com.google.javascript.jscomp.InlineAndCollapseProperties$InlineAliases$AliasesCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_InlineAliases_AliasesCollector](v_)
	})
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) SetOuter_(o *InlineAndCollapseProperties_InlineAliases) {
	this_.outer_ = o
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) GetOuter_() *InlineAndCollapseProperties_InlineAliases {
	return this_.outer_
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) Implements_jscomp_NodeTraversal_ExternsSkippingCallback() {
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) initFields_InlineAndCollapseProperties_InlineAliases_AliasesCollector() {
	_ = this_
}

// Java constructor: AliasesCollector()
func NewInlineAndCollapseProperties_InlineAliases_AliasesCollector(outer_ *InlineAndCollapseProperties_InlineAliases) *InlineAndCollapseProperties_InlineAliases_AliasesCollector {
	this_ := &InlineAndCollapseProperties_InlineAliases_AliasesCollector{}
	this_.outer_ = outer_
	InitInlineAndCollapseProperties_InlineAliases_AliasesCollector(this_)
	return this_
}

func InitInlineAndCollapseProperties_InlineAliases_AliasesCollector(this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) {
	InitNodeTraversal_ExternsSkippingCallback(&this_.NodeTraversal_ExternsSkippingCallbackBase)
	this_.initFields_InlineAndCollapseProperties_InlineAliases_AliasesCollector()
}

// Java method: AliasesCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
		if n.HasOneChild() && t.InGlobalScope() {
			this_.priv_AliasesCollector_visitAliasDefinition_3263(n.GetFirstChild(), NodeUtil_GetBestJSDocInfo(n.GetFirstChild()))
		}
	case int32(51) /* Token.ASSIGN */ :
		if ((parent != 0) && parent.IsExprResult()) && t.InGlobalScope() {
			this_.priv_AliasesCollector_visitAliasDefinition_3263(n.GetFirstChild(), n.GetJSDocInfo())
		}
	default:
	}
}

// Java method: AliasesCollector.visitAliasDefinition(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.JSDocInfo info)
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) priv_AliasesCollector_visitAliasDefinition_3263(lhs pkg_rhino.Node, info *pkg_rhino.JSDocInfo) {
	_ = lhs
	_ = info
	if (this_.priv_AliasesCollector_isDeclaredConst_3264(lhs, info) && ((info == nil) || !(info.HasTypeInformation()))) && lhs.IsQualifiedName() {
		var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(lhs)
		_ = rhs
		if (rhs != 0) && rhs.IsQualifiedName() {
			var lhsName *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](this_.outer_.namespace.GetOwnSlot(lhs.GetQualifiedName()))
			_ = lhsName
			var rhsName *GlobalNamespace_Name = j2gort.CastPtr[GlobalNamespace_Name](this_.outer_.namespace.GetOwnSlot(rhs.GetQualifiedName()))
			_ = rhsName
			if (((lhsName != nil) && lhsName.CalculateInlinability().ShouldInlineUsages()) && (rhsName != nil)) && rhsName.CalculateInlinability().ShouldInlineUsages() {
				this_.outer_.aliases.Put(lhs.GetQualifiedName(), rhs.GetQualifiedName())
			}
		}
	}
}

// Java method: AliasesCollector.isDeclaredConst(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.JSDocInfo info)
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesCollector) priv_AliasesCollector_isDeclaredConst_3264(lhs pkg_rhino.Node, info *pkg_rhino.JSDocInfo) bool {
	_ = lhs
	_ = info
	if (info != nil) && info.HasConstAnnotation() {
		return true
	}
	return lhs.GetParent().IsConst()
}

// Java class: com.google.javascript.jscomp.InlineAndCollapseProperties.InlineAliases.AliasesInliner
type InlineAndCollapseProperties_InlineAliases_AliasesInliner struct {
	NodeTraversal_ExternsSkippingCallbackBase
	outer_ *InlineAndCollapseProperties_InlineAliases
}

func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) IsJavaNil_() bool {
	return this_ == nil
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineAndCollapseProperties$InlineAliases$AliasesInliner", "AliasesInliner"
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_ExternsSkippingCallbackBase.IdentityHash_()
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineAndCollapseProperties.InlineAliases.AliasesInliner", "com.google.javascript.jscomp.InlineAndCollapseProperties$InlineAliases$AliasesInliner", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineAndCollapseProperties_InlineAliases_AliasesInliner](v_)
	})
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) SetOuter_(o *InlineAndCollapseProperties_InlineAliases) {
	this_.outer_ = o
}
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) GetOuter_() *InlineAndCollapseProperties_InlineAliases {
	return this_.outer_
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) Implements_jscomp_NodeTraversal_ExternsSkippingCallback() {
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) initFields_InlineAndCollapseProperties_InlineAliases_AliasesInliner() {
	_ = this_
}

// Java constructor: AliasesInliner()
func NewInlineAndCollapseProperties_InlineAliases_AliasesInliner(outer_ *InlineAndCollapseProperties_InlineAliases) *InlineAndCollapseProperties_InlineAliases_AliasesInliner {
	this_ := &InlineAndCollapseProperties_InlineAliases_AliasesInliner{}
	this_.outer_ = outer_
	InitInlineAndCollapseProperties_InlineAliases_AliasesInliner(this_)
	return this_
}

func InitInlineAndCollapseProperties_InlineAliases_AliasesInliner(this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) {
	InitNodeTraversal_ExternsSkippingCallback(&this_.NodeTraversal_ExternsSkippingCallbackBase)
	this_.initFields_InlineAndCollapseProperties_InlineAliases_AliasesInliner()
}

// Java method: AliasesInliner.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */ :
		if n.IsQualifiedName() && this_.outer_.aliases.ContainsKey(n.GetQualifiedName()) {
			if this_.priv_AliasesInliner_isLeftmostNameLocal_3265(t, n) {
				return
			}
			if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
				return
			}
			var newNode pkg_rhino.Node = this_.outer_.astFactory.CreateQName_2_StaticScope_string(pkg_rhino.ToIface_StaticScope(this_.outer_.namespace), this_.priv_AliasesInliner_resolveAlias_3266(n.GetQualifiedName(), n))
			_ = newNode
			if this_.priv_AliasesInliner_isLeftmostNameLocal_3265(t, newNode) {
				return
			}
			newNode.SrcrefTree(n)
			if newNode.IsGetProp() {
				newNode.GetFirstChild().MakeNonIndexableRecursive()
			}
			n.ReplaceWith(newNode)
			t.ReportCodeChange_0()
		}
	default:
	}
}

// Java method: AliasesInliner.isLeftmostNameLocal(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) priv_AliasesInliner_isLeftmostNameLocal_3265(t *NodeTraversal, n pkg_rhino.Node) bool {
	_ = t
	_ = n
	if !(n.IsQualifiedName()) {
		j2gort.CheckState(false)
	}
	var leftmostName string = NodeUtil_GetRootOfQualifiedName_1_Node(n).GetString()
	_ = leftmostName
	var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(leftmostName))
	_ = v
	return ((v != nil) && v.IsLocal())
}

// Java method: AliasesInliner.resolveAlias(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *InlineAndCollapseProperties_InlineAliases_AliasesInliner) priv_AliasesInliner_resolveAlias_3266(name string, n pkg_rhino.Node) string {
	_ = name
	_ = n
	var aliasPath j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = aliasPath
	for this_.outer_.aliases.ContainsKey(name) {
		if !(aliasPath.Add(name)) {
			this_.outer_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_InlineAndCollapseProperties_ALIAS_CYCLE(), []string{(&aliasPath).String(), name}))
			this_.outer_.aliases.Remove(name)
			break
		}
		name = j2gort.OrderedMapGetStr(this_.outer_.aliases, name)
	}
	return name
}
func (*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) AlwaysShouldTraverse_()      {}
func (*InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) CanReleaseSyntacticScopes_() {}
func (*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) AlwaysShouldTraverse_() {
}
func (*InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) CanReleaseSyntacticScopes_() {
}
func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) CanReleaseSyntacticScopes_() {}
func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) CanReleaseSyntacticScopes_()   {}
