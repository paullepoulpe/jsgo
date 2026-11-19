/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypeInference.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomptype "github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.TypeInference
type TypeInference struct {
	DataFlowAnalysisBase[pkg_rhino.Node, pkg_jscomptype.FlowScope]
	compiler                AbstractCompiler                          `gson:"compiler"`
	registry                *pkg_rhino.JSTypeRegistry                 `gson:"registry"`
	reverseInterpreter      pkg_jscomptype.ReverseAbstractInterpreter `gson:"reverseInterpreter"`
	bottomScope             pkg_jscomptype.FlowScope                  `gson:"bottomScope"`
	containerScope          *TypedScope                               `gson:"containerScope"`
	scopeCreator            *TypedScopeCreator                        `gson:"scopeCreator"`
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup `gson:"assertionFunctionLookup"`
	moduleImportResolver    *ModuleImportResolver                     `gson:"moduleImportResolver"`
	optChainArrayDeque      j2gort.Deque[*TypeInference_OptChainInfo] `gson:"optChainArrayDeque"`
	inferredUnboundVars     j2gort.OrderedSet[*TypedScope]            `gson:"inferredUnboundVars"`
	flowBrancher            *TypeInference_ReusableFlowBrancher       `gson:"flowBrancher"`
	flowScopePool           *LinkedFlowScope_Pool                     `gson:"flowScopePool"`
	reusableJoinOp          *LinkedFlowScope_FlowScopeJoinOp          `gson:"reusableJoinOp"`
	unknownType             pkg_rhino.ObjectType                      `gson:"unknownType"`
	numberAdditionSupertype pkg_rhino.JSType                          `gson:"numberAdditionSupertype"`
}

func (this_ *TypeInference) BindSelf_(s any) {
	_ = s
	this_.DataFlowAnalysisBase.BindSelf_(s)
}
func (this_ *TypeInference) IsJavaNil_() bool { return this_ == nil }
func (*TypeInference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInference", "TypeInference"
}
func (this_ *TypeInference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.DataFlowAnalysisBase.IdentityHash_()
}
func (*TypeInference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference", "com.google.javascript.jscomp.TypeInference", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInference](v_)
	})
}
func (*TypeInference) Implements_jscomp_DataFlowAnalysis() {}

var static_TypeInference_REASSIGN_CLASS_PROTOTYPE *DiagnosticType

func StaticGet_TypeInference_REASSIGN_CLASS_PROTOTYPE() *DiagnosticType {
	clinit_TypeInference()
	return static_TypeInference_REASSIGN_CLASS_PROTOTYPE
}
func StaticSet_TypeInference_REASSIGN_CLASS_PROTOTYPE(v *DiagnosticType) {
	clinit_TypeInference()
	static_TypeInference_REASSIGN_CLASS_PROTOTYPE = v
}
func StaticGet_TypeInference_GOOG_REQUIREDYNAMIC_NAME() string { return "goog.requireDynamic" }

var static_TypeInference_TOKENS_ALLOWING_NULL_TYPES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_TypeInference_TOKENS_ALLOWING_NULL_TYPES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_TypeInference()
	return static_TypeInference_TOKENS_ALLOWING_NULL_TYPES
}
func StaticSet_TypeInference_TOKENS_ALLOWING_NULL_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_TypeInference()
	static_TypeInference_TOKENS_ALLOWING_NULL_TYPES = v
}

var clinitDone_TypeInference uint32
var clinitState_TypeInference j2gort.ClinitState

func clinit_TypeInference() {
	if atomic.LoadUint32(&clinitDone_TypeInference) != 0 {
		return
	}
	clinitSlow_TypeInference()
}

func clinitSlow_TypeInference() {
	if !j2gort.EnterClinit(&clinitState_TypeInference, &clinitDone_TypeInference) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeInference, &clinitDone_TypeInference)
	static_TypeInference_REASSIGN_CLASS_PROTOTYPE = DiagnosticType_Error("JSC_REASSIGN_CLASS_PROTOTYPE", "Reassigning a class prototype is not allowed")
	static_TypeInference_TOKENS_ALLOWING_NULL_TYPES = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_NAME()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_CALL()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_NEW())})
}

func (this_ *TypeInference) initFields_TypeInference() {
	_ = this_
	this_.optChainArrayDeque = j2gort.Deque[*TypeInference_OptChainInfo]{}
	this_.inferredUnboundVars = j2gort.NewOrderedSet[*TypedScope]()
	this_.flowBrancher = NewTypeInference_ReusableFlowBrancher(this_)
	this_.flowScopePool = NewLinkedFlowScope_Pool()
}

// Java constructor: TypeInference(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, com.google.javascript.jscomp.type.ReverseAbstractInterpreter reverseInterpreter, com.google.javascript.jscomp.TypedScope syntacticScope, com.google.javascript.jscomp.TypedScopeCreator scopeCreator, com.google.javascript.jscomp.CodingConvention.AssertionFunctionLookup assertionFunctionLookup)
func NewTypeInference(compiler AbstractCompiler, cfg ControlFlowGraph[pkg_rhino.Node], reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, syntacticScope *TypedScope, scopeCreator *TypedScopeCreator, assertionFunctionLookup *CodingConvention_AssertionFunctionLookup) *TypeInference {
	clinit_TypeInference()
	this_ := &TypeInference{}
	this_.BindSelf_(this_)
	InitTypeInference(this_, compiler, cfg, reverseInterpreter, syntacticScope, scopeCreator, assertionFunctionLookup)
	return this_
}

func InitTypeInference(this_ *TypeInference, compiler AbstractCompiler, cfg ControlFlowGraph[pkg_rhino.Node], reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, syntacticScope *TypedScope, scopeCreator *TypedScopeCreator, assertionFunctionLookup *CodingConvention_AssertionFunctionLookup) {
	_ = compiler
	_ = cfg
	_ = reverseInterpreter
	_ = syntacticScope
	_ = scopeCreator
	_ = assertionFunctionLookup
	InitDataFlowAnalysis[pkg_rhino.Node, pkg_jscomptype.FlowScope](&this_.DataFlowAnalysisBase, cfg)
	this_.initFields_TypeInference()
	if !(syntacticScope.GetRootNode_0() == cfg.GetEntry().GetValue()) {
		j2gort.CheckArgumentFmt(false, "Expected syntactic scope to be rooted at CFG root node of %s, but instead got syntactic scope root node of %s", cfg.GetEntry(), j2gort.BoxNode(syntacticScope.GetRootNode_0()))
	}
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.registry = j2gort.CheckNotNullPtr(compiler.GetTypeRegistry())
	this_.reverseInterpreter = j2gort.CheckNotNullIface(reverseInterpreter)
	this_.moduleImportResolver = NewModuleImportResolver(compiler.GetModuleMap(), scopeCreator.GetNodeToScopeMapper(), this_.registry)
	this_.containerScope = syntacticScope
	this_.scopeCreator = j2gort.CheckNotNullPtr(scopeCreator)
	this_.assertionFunctionLookup = j2gort.CheckNotNullPtr(assertionFunctionLookup)
	this_.reusableJoinOp = NewLinkedFlowScope_FlowScopeJoinOp_2(Upcast_CompilerInputProvider(this_.compiler), this_.flowScopePool)
	this_.unknownType = this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	this_.numberAdditionSupertype = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NUMBER_ADDITION_SUPERTYPE())
	this_.bottomScope = pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateEntryLattice_3(Upcast_CompilerInputProvider(compiler), allocPooledVia_TypedScope_CreateLatticeBottom(&this_.scopeCreator.tsArena_.slab_TypedScope, syntacticScope.GetRootNode_0()), this_.flowScopePool))
}

// Java method: TypeInference.resetForScope(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, com.google.javascript.jscomp.TypedScope syntacticScope)
func (this_ *TypeInference) ResetForScope(cfg ControlFlowGraph[pkg_rhino.Node], syntacticScope *TypedScope) {
	_ = cfg
	_ = syntacticScope
	if !(syntacticScope.GetRootNode_0() == cfg.GetEntry().GetValue()) {
		j2gort.CheckArgumentFmt(false, "Expected syntactic scope to be rooted at CFG root node of %s, but instead got syntactic scope root node of %s", cfg.GetEntry(), j2gort.BoxNode(syntacticScope.GetRootNode_0()))
	}
	this_.flowScopePool.Reset()
	this_.InitCfg(cfg)
	this_.containerScope = syntacticScope
	this_.bottomScope = pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateEntryLattice_3(Upcast_CompilerInputProvider(this_.compiler), allocPooledVia_TypedScope_CreateLatticeBottom(&this_.scopeCreator.tsArena_.slab_TypedScope, syntacticScope.GetRootNode_0()), this_.flowScopePool))
	if !(this_.optChainArrayDeque.IsEmpty()) {
		this_.optChainArrayDeque.Clear()
	}
	if !(this_.inferredUnboundVars.IsEmpty()) {
		this_.inferredUnboundVars.Clear()
	}
	this_.flowBrancher.Clear()
}

// Java method: TypeInference.finishScope()
func (this_ *TypeInference) FinishScope() {
	this_.ReleaseCfgReference()
	this_.bottomScope = nil
	this_.containerScope = nil
	if !(this_.optChainArrayDeque.IsEmpty()) {
		this_.optChainArrayDeque.Clear()
	}
	if !(this_.inferredUnboundVars.IsEmpty()) {
		this_.inferredUnboundVars.Clear()
	}
	this_.flowBrancher.Clear()
	this_.flowScopePool.Reset()
}

// Java method: TypeInference.needsNodeAnnotations()
func (this_ *TypeInference) NeedsNodeAnnotations() bool {
	return this_.compiler.IsDebugLoggingEnabled()
}

// Java method: TypeInference.isLatticeEqual(com.google.javascript.jscomp.type.FlowScope a, com.google.javascript.jscomp.type.FlowScope b)
func (this_ *TypeInference) IsLatticeEqual(a pkg_jscomptype.FlowScope, b pkg_jscomptype.FlowScope) bool {
	_ = a
	_ = b
	return a.EqualsFlowScope(b)
}

// Java method: TypeInference.inferDeclarativelyUnboundVarsWithoutTypes(com.google.javascript.jscomp.type.FlowScope flow)
func (this_ *TypeInference) priv_TypeInference_inferDeclarativelyUnboundVarsWithoutTypes_5495(flow pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = flow
	var scope *TypedScope = j2gort.CastPtr[TypedScope](j2gort.AsIface[pkg_rhino.StaticTypedScope](flow.GetDeclarationScope()))
	_ = scope
	if !(this_.inferredUnboundVars.Add(scope)) {
		return flow
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(scope.GetDeclarativelyUnboundVarsWithoutTypes()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			var_ := j2gort.CastPtr[TypedVar](rawElem_2)
			_ = var_
			if this_.priv_TypeInference_isUnflowable_5570(var_) {
				continue
			}
			flow = j2gort.AsIface[pkg_jscomptype.FlowScope](flow.InferSlotType(var_.GetName(), this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())))
		}
	}
	return flow
}

// Java method: TypeInference.inferParameters(com.google.javascript.jscomp.type.FlowScope entryFlowScope)
func (this_ *TypeInference) priv_TypeInference_inferParameters_5496(entryFlowScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = entryFlowScope
	var functionNode pkg_rhino.Node = this_.containerScope.GetRootNode_0()
	_ = functionNode
	if !(functionNode.IsFunction()) {
		return entryFlowScope
	} else if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		return entryFlowScope
	}
	var astParameters pkg_rhino.Node = functionNode.GetSecondChild()
	_ = astParameters
	var iifeArgumentNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = iifeArgumentNode
	if NodeUtil_IsInvocationTarget(functionNode) {
		iifeArgumentNode = functionNode.GetNext()
	}
	var functionType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(functionNode.GetJSType())
	_ = functionType
	var parameterTypes []pkg_rhino.FunctionType_Parameter = functionType.GetParameterArray()
	_ = parameterTypes
	var paramLen int32 = int32(len(parameterTypes))
	_ = paramLen
	var paramIdx int32 = int32(0)
	_ = paramIdx
	{
		var astParamItr pkg_rhino.Node = astParameters.GetFirstChild()
		_ = astParamItr
		for ; astParamItr != 0; astParamItr = astParamItr.GetNext() {
			var astParam pkg_rhino.Node = astParamItr
			_ = astParam
			if (iifeArgumentNode != 0) && iifeArgumentNode.IsSpread() {
				iifeArgumentNode = pkg_rhino.Node(0)
			}
			var inferredType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(astParam)
			_ = inferredType
			var parameter pkg_rhino.FunctionType_Parameter
			_ = parameter
			if paramIdx < paramLen {
				parameter = parameterTypes[func() int32 { old_3 := paramIdx; nw_4 := old_3 + 1; paramIdx = nw_4; return old_3 }()]
			} else {
				parameter = (pkg_rhino.FunctionType_Parameter{})
			}
			if iifeArgumentNode != 0 {
				if !(j2gort.IfaceEquals[pkg_rhino.JSType](iifeArgumentNode.GetJSType(), nil)) {
					inferredType = iifeArgumentNode.GetJSType()
				}
			} else if parameter != (pkg_rhino.FunctionType_Parameter{}) {
				if !(j2gort.IfaceEquals[pkg_rhino.JSType](parameter.GetJSType(), nil)) {
					inferredType = parameter.GetJSType()
				}
			}
			var defaultValue pkg_rhino.Node = pkg_rhino.Node(0)
			_ = defaultValue
			if astParam.IsDefaultValue() {
				defaultValue = astParam.GetSecondChild()
				entryFlowScope = this_.priv_TypeInference_traverse_5500(defaultValue, entryFlowScope)
				astParam = astParam.GetFirstChild()
			} else if astParam.IsRest() {
				astParam = astParam.GetOnlyChild()
				inferredType = pkg_rhino.ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), []pkg_rhino.JSType{inferredType}))
			}
			if defaultValue != 0 {
				inferredType = this_.registry.CreateUnionType_2(inferredType.RestrictByNotUndefined(), this_.priv_TypeInference_getJSType_5571(defaultValue))
			}
			if astParam.IsDestructuringPattern() {
				entryFlowScope = this_.priv_TypeInference_updateDestructuringParameter_5497(astParam, inferredType, entryFlowScope)
			} else {
				entryFlowScope = this_.priv_TypeInference_updateNamedParameter_5498(astParam, (defaultValue != 0), inferredType, entryFlowScope)
			}
			if iifeArgumentNode != 0 {
				iifeArgumentNode = iifeArgumentNode.GetNext()
			} else {
				iifeArgumentNode = pkg_rhino.Node(0)
			}
		}
	}
	return entryFlowScope
}

// Java method: TypeInference.updateDestructuringParameter(com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.jstype.JSType inferredType, com.google.javascript.jscomp.type.FlowScope entryFlowScope)
func (this_ *TypeInference) priv_TypeInference_updateDestructuringParameter_5497(pattern pkg_rhino.Node, inferredType pkg_rhino.JSType, entryFlowScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = pattern
	_ = inferredType
	_ = entryFlowScope
	entryFlowScope = this_.priv_TypeInference_traverseDestructuringPatternHelper_5521(pattern, entryFlowScope, inferredType, NewLambda_TypeInference_TypeDeclaringCallback(func(scope pkg_jscomptype.FlowScope, lvalue pkg_rhino.Node, type_ pkg_rhino.JSType) pkg_jscomptype.FlowScope {
		_ = scope
		_ = lvalue
		_ = type_
		var var_ *TypedVar = this_.containerScope.LookupTypedVar(lvalue.GetString())
		_ = var_
		j2gort.CheckNotNullPtr(var_)
		if var_.IsTypeInferred() {
			var_.SetType(type_)
			lvalue.SetJSType(type_)
		}
		if lvalue.GetParent().IsDefaultValue() {
			scope = this_.priv_TypeInference_updateScopeForAssignment_5513(scope, lvalue, type_, StaticGet_TypeInference_AssignmentType_ASSIGN())
		}
		return scope
	}))
	return entryFlowScope
}

// Java method: TypeInference.updateNamedParameter(com.google.javascript.rhino.Node paramName, boolean hasDefaultValue, com.google.javascript.rhino.jstype.JSType inferredType, com.google.javascript.jscomp.type.FlowScope entryFlowScope)
func (this_ *TypeInference) priv_TypeInference_updateNamedParameter_5498(paramName pkg_rhino.Node, hasDefaultValue bool, inferredType pkg_rhino.JSType, entryFlowScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = paramName
	_ = hasDefaultValue
	_ = inferredType
	_ = entryFlowScope
	var var_ *TypedVar = this_.containerScope.LookupTypedVar(paramName.GetString())
	_ = var_
	j2gort.CheckNotNullPtr(var_)
	paramName.SetJSType(inferredType)
	if var_.IsTypeInferred() {
		var_.SetType(inferredType)
	} else if hasDefaultValue {
		entryFlowScope = this_.priv_TypeInference_redeclareSimpleVar_5569(entryFlowScope, paramName, inferredType)
	}
	return entryFlowScope
}

// Java method: TypeInference.createInitialEstimateLattice()
func (this_ *TypeInference) CreateInitialEstimateLattice() pkg_jscomptype.FlowScope {
	return this_.bottomScope
}

// Java method: TypeInference.createEntryLattice()
func (this_ *TypeInference) CreateEntryLattice() pkg_jscomptype.FlowScope {
	var entryScope pkg_jscomptype.FlowScope = this_.priv_TypeInference_inferDeclarativelyUnboundVarsWithoutTypes_5495(pkg_jscomptype.ToIface_FlowScope(LinkedFlowScope_CreateEntryLattice_3(Upcast_CompilerInputProvider(this_.compiler), this_.containerScope, this_.flowScopePool)))
	_ = entryScope
	return this_.priv_TypeInference_inferParameters_5496(entryScope)
}

// Java method: TypeInference.flowThrough(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope input)
func (this_ *TypeInference) FlowThrough(n pkg_rhino.Node, input pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = input
	if input == this_.bottomScope || ((input == nil || input.IsJavaNil_()) && (this_.bottomScope == nil || this_.bottomScope.IsJavaNil_())) {
		return input
	}
	var root pkg_rhino.Node = NodeUtil_GetEnclosingScopeRoot(n)
	_ = root
	var module *Module = ModuleImportResolver_GetModuleFromScopeRoot(this_.compiler.GetModuleMap(), Upcast_CompilerInputProvider(this_.compiler), root)
	_ = module
	var syntacticBlockScope *TypedScope = this_.scopeCreator.CreateScope_1(root)
	_ = syntacticBlockScope
	if (module != nil) && module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module() {
		this_.moduleImportResolver.DeclareEsModuleImports(module, syntacticBlockScope, this_.compiler.GetInput(NodeUtil_GetInputId(n)))
	}
	var output pkg_jscomptype.FlowScope = input.WithSyntacticScope(pkg_rhino.ToIface_StaticTypedScope(syntacticBlockScope))
	_ = output
	output = this_.priv_TypeInference_inferDeclarativelyUnboundVarsWithoutTypes_5495(output)
	output = this_.priv_TypeInference_traverse_5500(n, output)
	this_.priv_TypeInference_updateModuleScope_5499(module, syntacticBlockScope)
	return output
}

// Java method: TypeInference.updateModuleScope(com.google.javascript.jscomp.modules.Module module, com.google.javascript.jscomp.TypedScope syntacticBlockScope)
func (this_ *TypeInference) priv_TypeInference_updateModuleScope_5499(module *Module, syntacticBlockScope *TypedScope) {
	_ = module
	_ = syntacticBlockScope
	if module == nil {
		return
	}
	switch (module.Metadata().ModuleType()).Ordinal() {
	case int32(0) /* ModuleType.ES6_MODULE */ :
		this_.moduleImportResolver.UpdateEsModuleNamespaceType(j2gort.CastPtr[TypedVar](syntacticBlockScope.GetVar("*exports*")).GetType().GetJSTypeBase().ToObjectType(), module, syntacticBlockScope)
		return
	case int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		var exportsVar *TypedVar = j2gort.CheckNotNullPtr(j2gort.CastPtr[TypedVar](syntacticBlockScope.GetVar("exports")))
		_ = exportsVar
		var exportsType pkg_rhino.JSType
		_ = exportsType
		if !(j2gort.IfaceEquals[pkg_rhino.JSType](exportsVar.GetType(), nil)) {
			exportsType = exportsVar.GetType()
		} else {
			exportsType = pkg_rhino.Upcast_JSType(this_.unknownType)
		}
		var rootNode pkg_rhino.Node = syntacticBlockScope.GetRootNode_0()
		_ = rootNode
		if rootNode.IsModuleBody() {
			rootNode.SetJSType(exportsType)
		} else {
			if !(rootNode.IsBlock()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(rootNode))
			}
			var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(rootNode.GetParent())
			_ = paramList
			paramList.GetOnlyChild().SetJSType(exportsType)
		}
		if !(module.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsLegacyGoogModule()) || j2gort.IfaceEquals[pkg_rhino.JSType](exportsVar.GetType(), nil) {
			break
		}
		var moduleId string = module.ClosureNamespace()
		_ = moduleId
		var globalScope *TypedScope = syntacticBlockScope.GetGlobalTypedScope()
		_ = globalScope
		var globalVar *TypedVar = globalScope.LookupTypedVar(moduleId)
		_ = globalVar
		if globalVar.IsTypeInferred() {
			globalVar.SetType(exportsType)
		}
		var lastDot int32 = j2gort.JStrLastIndexOfCp(moduleId, int32(uint16(46)))
		_ = lastDot
		if lastDot >= int32(0) {
			var propName string = j2gort.JStrSubFrom(moduleId, (lastDot + int32(1)))
			_ = propName
			var parentType pkg_rhino.JSType = globalScope.LookupQualifiedName_1_string(j2gort.JStrSub(moduleId, int32(0), lastDot))
			_ = parentType
			var parentObjectType pkg_rhino.ObjectType
			_ = parentObjectType
			if parentType != nil && !parentType.IsJavaNil_() {
				parentObjectType = parentType.GetJSTypeBase().ToMaybeObjectType()
			} else {
				parentObjectType = nil
			}
			if (parentObjectType != nil && !parentObjectType.IsJavaNil_()) && !(parentObjectType.GetObjectTypeBase().IsPropertyTypeDeclared_1_string(propName)) {
				parentObjectType.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(propName, exportsType, exportsVar.GetNode())
			}
		}
		return
	default:
	}
}

// Java method: TypeInference.isForward()
func (this_ *TypeInference) IsForward() bool {
	return true
}

// Java method: TypeInference.isBranched()
func (this_ *TypeInference) IsBranched() bool {
	return true
}

// Java method: TypeInference.createFlowJoiner()
func (this_ *TypeInference) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[pkg_jscomptype.FlowScope] {
	if this_.reusableJoinOp.TryAcquire() {
		return ToIface_DataFlowAnalysis_FlowJoiner[pkg_jscomptype.FlowScope](this_.reusableJoinOp)
	} else {
		return ToIface_DataFlowAnalysis_FlowJoiner[pkg_jscomptype.FlowScope](NewLinkedFlowScope_FlowScopeJoinOp_2(Upcast_CompilerInputProvider(this_.compiler), this_.flowScopePool))
	}
}

// Java method: TypeInference.createFlowBrancher(com.google.javascript.rhino.Node source, com.google.javascript.jscomp.type.FlowScope output)
func (this_ *TypeInference) CreateFlowBrancher(source pkg_rhino.Node, output pkg_jscomptype.FlowScope) DataFlowAnalysis_FlowBrancher[pkg_jscomptype.FlowScope] {
	_ = source
	_ = output
	return this_.flowBrancher.Reset(source, output)
}

// Java method: TypeInference.traverse(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverse_5500(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var isTypeable bool = true
	_ = isTypeable
	switch n.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		scope = this_.priv_TypeInference_traverseAssign_5511(n, scope)
	case int32(32) /* Token.NAME */ :
		scope = this_.priv_TypeInference_traverseName_5522(n, scope)
	case int32(29) /* Token.OPTCHAIN_GETPROP */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		scope = this_.priv_TypeInference_traverseOptChain_5558(n, scope)
	case int32(26) /* Token.GETPROP */ :
		scope = this_.priv_TypeInference_traverseGetProp_5555(n, scope)
	case int32(105) /* Token.CLASS */ :
		scope = this_.priv_TypeInference_traverseClass_5527(n, scope)
	case int32(65) /* Token.ASSIGN_AND */, int32(64) /* Token.ASSIGN_OR */ :
		scope = this_.priv_TypeInference_traverseShortCircuitingBinOpAssignment_5563(n, scope)
	case int32(69) /* Token.AND */ :
		scope = this_.priv_TypeInference_traverseAnd_5548(n, scope).GetJoinedFlowScope()
	case int32(68) /* Token.OR */ :
		scope = this_.priv_TypeInference_traverseOr_5562(n, scope).GetJoinedFlowScope()
	case int32(66) /* Token.ASSIGN_COALESCE */, int32(70) /* Token.COALESCE */ :
		scope = this_.priv_TypeInference_traverseNullishCoalesce_5534(n, scope)
	case int32(67) /* Token.HOOK */ :
		scope = this_.priv_TypeInference_traverseHook_5533(n, scope)
	case int32(47) /* Token.OBJECTLIT */ :
		scope = this_.priv_TypeInference_traverseObjectLiteral_5530(n, scope)
	case int32(28) /* Token.CALL */ :
		scope = this_.priv_TypeInference_traverseCall_5504(n, scope)
	case int32(23) /* Token.NEW */ :
		scope = this_.priv_TypeInference_traverseNew_5546(n, scope)
	case int32(135) /* Token.NEW_TARGET */ :
		this_.priv_TypeInference_traverseNewTarget_5507(n)
	case int32(58) /* Token.ASSIGN_ADD */, int32(13) /* Token.ADD */ :
		scope = this_.priv_TypeInference_traverseAdd_5531(n, scope)
	case int32(21) /* Token.POS */ :
		scope = this_.priv_TypeInference_traverseUnaryPlus_5524(n, scope)
	case int32(22) /* Token.NEG */, int32(20) /* Token.BITNOT */, int32(72) /* Token.DEC */, int32(71) /* Token.INC */ :
		scope = this_.priv_TypeInference_traverseBigIntCompatibleUnaryOperator_5525(n, scope)
	case int32(46) /* Token.ARRAYLIT */ :
		scope = this_.priv_TypeInference_traverseArrayLiteral_5529(n, scope)
	case int32(37) /* Token.THIS */ :
		n.SetJSType(scope.GetTypeOfThis())
	case int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(54) /* Token.ASSIGN_BITAND */, int32(53) /* Token.ASSIGN_BITXOR */, int32(52) /* Token.ASSIGN_BITOR */, int32(60) /* Token.ASSIGN_MUL */, int32(59) /* Token.ASSIGN_SUB */, int32(63) /* Token.ASSIGN_EXPONENT */ :
		scope = this_.priv_TypeInference_traverseAssignOp_5512(n, scope)
	case int32(57) /* Token.ASSIGN_URSH */ :
		scope = this_.priv_TypeInference_traverseAssignUnsignedRightShift_5503(n, scope)
	case int32(3) /* Token.BITAND */, int32(2) /* Token.BITXOR */, int32(1) /* Token.BITOR */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(14) /* Token.SUB */, int32(15) /* Token.MUL */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */, int32(18) /* Token.EXPONENT */ :
		scope = this_.priv_TypeInference_traverseBigIntCompatibleBinaryOperator_5526(n, scope)
	case int32(12) /* Token.URSH */ :
		scope = this_.priv_TypeInference_traverseUnsignedRightShift_5502(n, scope)
	case int32(50) /* Token.COMMA */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		n.SetJSType(this_.priv_TypeInference_getJSType_5571(n.GetLastChild()))
	case int32(131) /* Token.TEMPLATELIT */, int32(25) /* Token.TYPEOF */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()))
	case int32(132) /* Token.TEMPLATELIT_SUB */, int32(43) /* Token.THROW */, int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */, int32(116) /* Token.IMPORT */, int32(117) /* Token.IMPORT_SPECS */, int32(119) /* Token.IMPORT_STAR */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		isTypeable = false
	case int32(118) /* Token.IMPORT_SPEC */ :
		this_.priv_TypeInference_traverseImportSpec_5523(scope, n)
		isTypeable = false
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		scope = this_.priv_TypeInference_traverseTaggedTemplateLit_5505(n, scope)
	case int32(24) /* Token.DELPROP */, int32(6) /* Token.LT */, int32(7) /* Token.LE */, int32(8) /* Token.GT */, int32(9) /* Token.GE */, int32(19) /* Token.NOT */, int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(45) /* Token.INSTANCEOF */, int32(44) /* Token.IN */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE()))
	case int32(27) /* Token.GETELEM */ :
		scope = this_.priv_TypeInference_traverseGetElem_5550(n, scope)
	case int32(93) /* Token.EXPR_RESULT */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		if n.GetFirstChild().IsGetProp() {
			var getprop pkg_rhino.Node = n.GetFirstChild()
			_ = getprop
			var ownerType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(this_.priv_TypeInference_getJSType_5571(getprop.GetFirstChild()).RestrictByNotNullOrUndefined())
			_ = ownerType
			if ownerType != nil && !ownerType.IsJavaNil_() {
				this_.priv_TypeInference_ensurePropertyDeclaredHelper_5517(getprop, ownerType, scope)
			}
		}
		isTypeable = false
	case int32(75) /* Token.SWITCH */ :
		scope = this_.priv_TypeInference_traverse_5500(n.GetFirstChild(), scope)
		isTypeable = false
	case int32(0) /* Token.RETURN */ :
		scope = this_.priv_TypeInference_traverseReturn_5508(n, scope)
		isTypeable = false
	case int32(114) /* Token.YIELD */ :
		if n.IsYieldAll() {
			scope = this_.priv_TypeInference_traverseYieldAll_5509(n, scope)
		} else {
			scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
			n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()))
		}
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		scope = this_.priv_TypeInference_traverseDeclaration_5518(n, scope)
		isTypeable = false
	case int32(86) /* Token.CATCH */ :
		scope = this_.priv_TypeInference_traverseCatch_5510(n, scope)
		isTypeable = false
	case int32(101) /* Token.CAST */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		j2gort.CheckNotNullPtr(info)
		if info.HasType() {
			n.SetJSType(info.GetType().Evaluate(j2gort.AsIface[pkg_rhino.StaticTypedScope](scope.GetDeclarationScope()), this_.registry).GetJSTypeBase().Resolve(this_.registry.GetErrorReporter()))
		} else {
			n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
		}
	case int32(110) /* Token.SUPER */ :
		this_.priv_TypeInference_traverseSuper_5506(n, scope)
	case int32(115) /* Token.AWAIT */ :
		scope = this_.priv_TypeInference_traverseAwait_5566(n, scope)
	case int32(87) /* Token.VOID */ :
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()))
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	case int32(120) /* Token.EXPORT */ :
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
			var defaultExport *TypedVar = TypeInference_GetDeclaredVar(scope, "*default*")
			_ = defaultExport
			if defaultExport.IsTypeInferred() {
				defaultExport.SetType(this_.priv_TypeInference_getJSType_5571(n.GetOnlyChild()))
			}
		}
		isTypeable = false
	case int32(136) /* Token.IMPORT_META */ :
		n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
	case int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */, int32(123) /* Token.MODULE_BODY */, int32(91) /* Token.SWITCH_BODY */, int32(73) /* Token.FUNCTION */, int32(49) /* Token.PARAM_LIST */, int32(90) /* Token.BLOCK */, int32(88) /* Token.EMPTY */, int32(74) /* Token.IF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(82) /* Token.BREAK */, int32(83) /* Token.CONTINUE */, int32(48) /* Token.TRY */, int32(76) /* Token.CASE */, int32(77) /* Token.DEFAULT_CASE */, int32(85) /* Token.WITH */, int32(98) /* Token.DEBUGGER */, int32(121) /* Token.EXPORT_SPECS */, int32(92) /* Token.LABEL */ :
		isTypeable = false
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.priv_TypeInference_traverseDynamicImport_5567(n, scope)
	case int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(35) /* Token.STRINGLIT */, int32(33) /* Token.NUMBER */, int32(34) /* Token.BIGINT */, int32(36) /* Token.NULL */, int32(42) /* Token.REGEXP */, int32(133) /* Token.TEMPLATELIT_STRING */ :
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Type inference doesn't know to handle token "+j2gort.StrVal(n.GetToken()))))
	}
	if (isTypeable && j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil)) && !(j2gort.ImmSetContainsTypedU8_(StaticGet_TypeInference_TOKENS_ALLOWING_NULL_TYPES(), uint8(n.GetToken()))) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Failed to infer JSType for "+j2gort.StrVal(n))))
	}
	return scope
}

// Java method: TypeInference.initializeEnhancedForScope(com.google.javascript.rhino.Node source, com.google.javascript.jscomp.type.FlowScope output)
func (this_ *TypeInference) priv_TypeInference_initializeEnhancedForScope_5501(source pkg_rhino.Node, output pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = source
	_ = output
	var item pkg_rhino.Node = source.GetFirstChild()
	_ = item
	var obj pkg_rhino.Node = item.GetNext()
	_ = obj
	var informed pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverse_5500(obj, output)
	_ = informed
	var assignmentType TypeInference_AssignmentType = TypeInference_AssignmentType(0)
	_ = assignmentType
	if NodeUtil_IsNameDeclaration(item) {
		item = item.GetFirstChild()
		assignmentType = StaticGet_TypeInference_AssignmentType_DECLARATION()
	} else {
		assignmentType = StaticGet_TypeInference_AssignmentType_ASSIGN()
	}
	if item.IsDestructuringLhs() {
		item = item.GetFirstChild()
	}
	var newType pkg_rhino.JSType = nil
	_ = newType
	switch source.TokenSwitchOrd_() {
	case int32(81) /* Token.FOR_IN */ :
		var iterKeyType pkg_rhino.JSType = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
		_ = iterKeyType
		var objType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(obj).Autobox()
		_ = objType
		var objIndexType pkg_rhino.JSType = objType.GetTemplateTypeMap().GetResolvedTemplateType(this_.registry.GetObjectIndexKey())
		_ = objIndexType
		if (objIndexType != nil && !objIndexType.IsJavaNil_()) && !(objIndexType.GetJSTypeBase().IsUnknownType()) {
			var narrowedKeyType pkg_rhino.JSType = iterKeyType.GetJSTypeBase().GetGreatestSubtype(objIndexType)
			_ = narrowedKeyType
			if !(narrowedKeyType.GetJSTypeBase().IsEmptyType()) {
				iterKeyType = narrowedKeyType
			}
		}
		newType = iterKeyType
	case int32(112) /* Token.FOR_OF */ :
		newType = JsIterables_MaybeBoxIterableOrAsyncIterable(this_.priv_TypeInference_getJSType_5571(obj), this_.registry).OrElse(pkg_rhino.Upcast_JSType(this_.unknownType))
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		var iterableType pkg_rhino.JSType = JsIterables_MaybeBoxIterableOrAsyncIterable(this_.priv_TypeInference_getJSType_5571(obj), this_.registry).OrElse(pkg_rhino.Upcast_JSType(this_.unknownType))
		_ = iterableType
		newType = Promises_GetResolvedType(this_.registry, iterableType)
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected source node "+j2gort.StrVal(source))))
	}
	if item.IsDestructuringPattern() {
		informed = this_.priv_TypeInference_traverseDestructuringPattern_5520(item, informed, newType, assignmentType)
	} else {
		informed = this_.priv_TypeInference_traverse_5500(item, informed)
		informed = this_.priv_TypeInference_updateScopeForAssignment_5513(informed, item, newType, assignmentType)
	}
	return informed
}

// Java method: TypeInference.traverseUnsignedRightShift(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseUnsignedRightShift_5502(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	if (TypeInference_GetBigIntPresence(this_.priv_TypeInference_getJSType_5571(n.GetFirstChild())) != StaticGet_TypeInference_BigIntPresence_NO_BIGINT()) || (TypeInference_GetBigIntPresence(this_.priv_TypeInference_getJSType_5571(n.GetLastChild())) != StaticGet_TypeInference_BigIntPresence_NO_BIGINT()) {
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE()))
	} else {
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))
	}
	return scope
}

// Java method: TypeInference.traverseAssignUnsignedRightShift(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseAssignUnsignedRightShift_5503(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	scope = this_.priv_TypeInference_traverseUnsignedRightShift_5502(n, scope)
	return this_.priv_TypeInference_updateScopeForAssignment_5514(scope, left, this_.priv_TypeInference_getJSType_5571(n), pkg_rhino.Node(0), StaticGet_TypeInference_AssignmentType_ASSIGN())
}

// Java method: TypeInference.traverseCall(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.type.FlowScope originalScope)
func (this_ *TypeInference) priv_TypeInference_traverseCall_5504(callNode pkg_rhino.Node, originalScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = callNode
	_ = originalScope
	if !(callNode.IsCall() || callNode.IsOptChainCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(callNode))
	}
	var scopeAfterChildren pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverseChildren_5549(callNode, originalScope)
	_ = scopeAfterChildren
	var scopeAfterExecution pkg_jscomptype.FlowScope = this_.priv_TypeInference_setCallNodeTypeAfterChildrenTraversed_5535(callNode, scopeAfterChildren)
	_ = scopeAfterExecution
	return this_.priv_TypeInference_tightenTypesAfterAssertions_5537(scopeAfterExecution, callNode)
}

// Java method: TypeInference.traverseTaggedTemplateLit(com.google.javascript.rhino.Node tagTempLitNode, com.google.javascript.jscomp.type.FlowScope originalScope)
func (this_ *TypeInference) priv_TypeInference_traverseTaggedTemplateLit_5505(tagTempLitNode pkg_rhino.Node, originalScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = tagTempLitNode
	_ = originalScope
	if !(tagTempLitNode.IsTaggedTemplateLit()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(tagTempLitNode))
	}
	var scopeAfterChildren pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverseChildren_5549(tagTempLitNode, originalScope)
	_ = scopeAfterChildren
	var scopeAfterExecution pkg_jscomptype.FlowScope = this_.priv_TypeInference_setCallNodeTypeAfterChildrenTraversed_5535(tagTempLitNode, scopeAfterChildren)
	_ = scopeAfterExecution
	if j2gort.IfaceEquals[pkg_rhino.JSType](tagTempLitNode.GetJSType(), nil) {
		tagTempLitNode.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
	}
	return scopeAfterExecution
}

// Java method: TypeInference.traverseSuper(com.google.javascript.rhino.Node superNode, com.google.javascript.jscomp.type.FlowScope currentScope)
func (this_ *TypeInference) priv_TypeInference_traverseSuper_5506(superNode pkg_rhino.Node, currentScope pkg_jscomptype.FlowScope) {
	_ = superNode
	_ = currentScope
	var superNodeType pkg_rhino.ObjectType = nil
	_ = superNodeType
	switch superNode.GetParent().TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		var scope *TypedScope = this_.containerScope
		_ = scope
		for (scope != nil) && !(NodeUtil_IsNonArrowFunction(scope.GetRootNode_0())) {
			scope = j2gort.CastPtr[TypedScope](scope.GetParent())
		}
		if scope == nil {
			superNode.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
			return
		}
		var enclosingFunctionType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(scope.GetRootNode_0().GetJSType())
		_ = enclosingFunctionType
		if (enclosingFunctionType != nil) && enclosingFunctionType.IsConstructor() {
			superNodeType = pkg_rhino.ToIface_ObjectType(enclosingFunctionType.GetSuperClassConstructor())
		}
	case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
		var currentSyntacticScope pkg_rhino.StaticTypedScope = j2gort.AsIface[pkg_rhino.StaticTypedScope](currentScope.GetDeclarationScope())
		_ = currentSyntacticScope
		superNodeType = pkg_rhino.ObjectType_Cast(j2gort.AsIface[pkg_rhino.StaticTypedSlot](currentSyntacticScope.GetSlot("super")).GetType())
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent of SUPER: "+j2gort.StrValString(superNode.GetParent().ToStringTree()))))
	}
	superNode.SetJSType(func() pkg_rhino.JSType {
		if superNodeType != nil && !superNodeType.IsJavaNil_() {
			return pkg_rhino.Upcast_JSType(superNodeType)
		}
		return pkg_rhino.Upcast_JSType(this_.unknownType)
	}())
}

// Java method: TypeInference.traverseNewTarget(com.google.javascript.rhino.Node newTargetNode)
func (this_ *TypeInference) priv_TypeInference_traverseNewTarget_5507(newTargetNode pkg_rhino.Node) {
	_ = newTargetNode
	var scope *TypedScope = this_.containerScope
	_ = scope
	for (scope != nil) && !(NodeUtil_IsNonArrowFunction(scope.GetRootNode_0())) {
		scope = j2gort.CastPtr[TypedScope](scope.GetParent())
	}
	if scope == nil {
		newTargetNode.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
		return
	}
	var root pkg_rhino.Node = scope.GetRootNode_0()
	_ = root
	var parent pkg_rhino.Node = root.GetParent()
	_ = parent
	if parent.GetGrandparent().IsClass() {
		var type_ pkg_rhino.JSTypeNative
		_ = type_
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(parent) {
			type_ = pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE()
		} else {
			type_ = pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()
		}
		newTargetNode.SetJSType(this_.registry.GetNativeType(type_))
	} else {
		newTargetNode.SetJSType(this_.registry.CreateUnionType_2(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE()), this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())))
	}
}

// Java method: TypeInference.traverseReturn(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseReturn_5508(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	var retValue pkg_rhino.Node = n.GetFirstChild()
	_ = retValue
	if retValue != 0 {
		var type_ pkg_rhino.JSType = this_.containerScope.GetRootNode_0().GetJSType()
		_ = type_
		if type_ != nil && !type_.IsJavaNil_() {
			var fnType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
			_ = fnType
			if fnType != nil {
				TypeInference_InferPropertyTypesToMatchConstraint(retValue.GetJSType(), fnType.GetReturnType())
			}
		}
	}
	return scope
}

// Java method: TypeInference.traverseYieldAll(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseYieldAll_5509(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	var innerType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(n.GetFirstChild())
	_ = innerType
	var yieldAllResult pkg_rhino.JSType = JsIterables_GetReturnElementType(innerType, this_.registry)
	_ = yieldAllResult
	n.SetJSType(yieldAllResult)
	return scope
}

// Java method: TypeInference.traverseCatch(com.google.javascript.rhino.Node catchNode, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseCatch_5510(catchNode pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = catchNode
	_ = scope
	var catchTarget pkg_rhino.Node = catchNode.GetFirstChild()
	_ = catchTarget
	if catchTarget.IsName() {
		var name pkg_rhino.Node = catchNode.GetFirstChild()
		_ = name
		var type_ pkg_rhino.JSType = name.GetJSType()
		_ = type_
		if type_ == nil || type_.IsJavaNil_() {
			type_ = pkg_rhino.Upcast_JSType(this_.unknownType)
			name.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
		}
		return this_.priv_TypeInference_redeclareSimpleVar_5569(scope, name, type_)
	} else if catchTarget.IsDestructuringPattern() {
		var pattern pkg_rhino.Node = catchNode.GetFirstChild()
		_ = pattern
		return this_.priv_TypeInference_traverseDestructuringPattern_5520(pattern, scope, pkg_rhino.Upcast_JSType(this_.unknownType), StaticGet_TypeInference_AssignmentType_DECLARATION())
	} else {
		if !(catchTarget.IsEmpty()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(catchTarget))
		}
		return scope
	}
}

// Java method: TypeInference.traverseAssign(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseAssign_5511(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var target pkg_rhino.Node = n.GetFirstChild()
	_ = target
	var value pkg_rhino.Node = n.GetLastChild()
	_ = value
	if target.IsDestructuringPattern() {
		scope = this_.priv_TypeInference_traverse_5500(value, scope)
		var valueType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(value)
		_ = valueType
		n.SetJSType(valueType)
		return this_.priv_TypeInference_traverseDestructuringPattern_5520(target, scope, valueType, StaticGet_TypeInference_AssignmentType_ASSIGN())
	} else {
		scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
		var valueType_2 pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(value)
		_ = valueType_2
		n.SetJSType(valueType_2)
		return this_.priv_TypeInference_updateScopeForAssignment_5513(scope, target, valueType_2, StaticGet_TypeInference_AssignmentType_ASSIGN())
	}
}

// Java method: TypeInference.traverseAssignOp(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseAssignOp_5512(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	scope = this_.priv_TypeInference_traverseBigIntCompatibleBinaryOperator_5526(n, scope)
	return this_.priv_TypeInference_updateScopeForAssignment_5514(scope, left, this_.priv_TypeInference_getJSType_5571(n), pkg_rhino.Node(0), StaticGet_TypeInference_AssignmentType_ASSIGN())
}

// Java static method: TypeInference.isInExternFile(com.google.javascript.rhino.Node n)
func TypeInference_IsInExternFile(n pkg_rhino.Node) bool {
	clinit_TypeInference()
	_ = n
	return NodeUtil_GetSourceFile(n).IsExtern()
}

// Java static method: TypeInference.isPossibleMixinApplication(com.google.javascript.rhino.Node lvalue, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rvalue)
func TypeInference_IsPossibleMixinApplication(lvalue pkg_rhino.Node, rvalue pkg_rhino.Node) bool {
	clinit_TypeInference()
	_ = lvalue
	_ = rvalue
	if TypeInference_IsInExternFile(lvalue) {
		return true
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(lvalue)
	_ = jsdoc
	return (((((jsdoc != nil) && jsdoc.IsConstructor()) && (jsdoc.GetImplementedInterfaceCount() > int32(0))) && lvalue.IsQualifiedName()) && ((rvalue != 0) && rvalue.IsCall()))
}

// Java static method: TypeInference.addMissingInterfaceProperties(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType constructor)
func TypeInference_AddMissingInterfaceProperties(constructor pkg_rhino.JSType) {
	clinit_TypeInference()
	_ = constructor
	if (constructor == nil || constructor.IsJavaNil_()) || !(constructor.GetJSTypeBase().IsConstructor()) {
		return
	}
	var f *pkg_rhino.FunctionType = constructor.ToMaybeFunctionType()
	_ = f
	var proto pkg_rhino.ObjectType = f.GetPrototype()
	_ = proto
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(f.GetImplementedInterfaces()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			interf := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_6)
			_ = interf
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(interf.GetObjectTypeBase().GetPropertyNames()))
				for cur_7.HasNext() {
					pname := cur_7.NextStr()
					_ = pname
					if !(proto.GetObjectTypeBase().HasProperty_1_string(pname)) {
						proto.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(pname, interf.GetObjectTypeBase().GetPropertyType_1_string(pname), pkg_rhino.Node(0))
					}
				}
			}
		}
	}
}

// Java method: TypeInference.updateScopeForAssignment(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node target, com.google.javascript.rhino.jstype.JSType resultType, com.google.javascript.jscomp.TypeInference.AssignmentType type)
func (this_ *TypeInference) priv_TypeInference_updateScopeForAssignment_5513(scope pkg_jscomptype.FlowScope, target pkg_rhino.Node, resultType pkg_rhino.JSType, type_ TypeInference_AssignmentType) pkg_jscomptype.FlowScope {
	_ = scope
	_ = target
	_ = resultType
	_ = type_
	return this_.priv_TypeInference_updateScopeForAssignment_5514(scope, target, resultType, target, type_)
}

// Java method: TypeInference.updateScopeForAssignment(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node target, com.google.javascript.rhino.jstype.JSType resultType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node updateNode, com.google.javascript.jscomp.TypeInference.AssignmentType type)
func (this_ *TypeInference) priv_TypeInference_updateScopeForAssignment_5514(scope pkg_jscomptype.FlowScope, target pkg_rhino.Node, resultType pkg_rhino.JSType, updateNode pkg_rhino.Node, type_ TypeInference_AssignmentType) pkg_jscomptype.FlowScope {
	_ = scope
	_ = target
	_ = resultType
	_ = updateNode
	_ = type_
	j2gort.CheckNotNullIface(resultType)
	if !((updateNode == 0) || (updateNode == target)) {
		j2gort.CheckState(false)
	}
	var targetType pkg_rhino.JSType = target.GetJSType()
	_ = targetType
	var right pkg_rhino.Node = NodeUtil_GetRValueOfLValue(target)
	_ = right
	if TypeInference_IsPossibleMixinApplication(target, right) {
		TypeInference_AddMissingInterfaceProperties(targetType)
	}
	switch target.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var varName string = target.GetString()
		_ = varName
		var var_ *TypedVar = TypeInference_GetDeclaredVar(scope, varName)
		_ = var_
		var varType pkg_rhino.JSType
		_ = varType
		if var_ == nil {
			varType = nil
		} else {
			varType = var_.GetType()
		}
		var isVarDeclaration bool = ((((type_ == StaticGet_TypeInference_AssignmentType_DECLARATION()) && (var_ != nil)) && !(var_.IsTypeInferred())) && (var_.GetNameNode() != 0))
		_ = isVarDeclaration
		var isTypelessConstDecl bool = ((((isVarDeclaration && var_.GetNameNode().IsName()) && NodeUtil_IsConstantDeclaration(var_.GetJSDocInfo(), var_.GetNameNode())) && !((var_.GetJSDocInfo() != nil) && var_.GetJSDocInfo().ContainsDeclarationExcludingTypelessConst())) && !(!(j2gort.IfaceEquals[pkg_rhino.JSType](var_.GetType(), nil)) && var_.GetType().GetJSTypeBase().IsNoResolvedType()))
		_ = isTypelessConstDecl
		var isVarTypeBetter bool = (((isVarDeclaration && !(resultType.GetJSTypeBase().IsNullType())) && !(resultType.GetJSTypeBase().IsVoidType())) && !(isTypelessConstDecl))
		_ = isVarTypeBetter
		if isVarTypeBetter {
			scope = this_.priv_TypeInference_redeclareSimpleVar_5569(scope, target, varType)
		} else {
			scope = this_.priv_TypeInference_redeclareSimpleVar_5569(scope, target, resultType)
		}
		if updateNode != 0 {
			updateNode.SetJSType(resultType)
		}
		if ((var_ != nil) && var_.IsTypeInferred()) && !(target.GetParent().IsLet() && !(target.HasChildren())) {
			var oldType pkg_rhino.JSType = var_.GetType()
			_ = oldType
			var iifeArg_1 pkg_rhino.
				JSType
		iifeSwitch_1:
			switch {
			default:

				if oldType == nil || oldType.IsJavaNil_() {
					iifeArg_1 = resultType
					break iifeSwitch_1
				}
				iifeArg_1 = oldType.GetLeastSupertype(resultType)
				break iifeSwitch_1
			}
			var_.SetType(iifeArg_1)
		} else if isTypelessConstDecl {
			var_.SetType(resultType)
		}
	case int32(26) /* Token.GETPROP */ :
		var lfs *LinkedFlowScope
		_ = lfs
		if (!(j2gort.InstanceOfBind[*LinkedFlowScope](scope, &lfs)) || !(lfs.HasGlobalFunctionScope())) && target.IsQualifiedName() {
			var qualifiedName string = target.GetQualifiedName()
			_ = qualifiedName
			var declaredSlotType bool = false
			_ = declaredSlotType
			var rawObjType pkg_rhino.JSType = target.GetFirstChild().GetJSType()
			_ = rawObjType
			if rawObjType != nil && !rawObjType.IsJavaNil_() {
				var objType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(rawObjType.RestrictByNotNullOrUndefined())
				_ = objType
				if objType != nil && !objType.IsJavaNil_() {
					var propName string = target.GetString()
					_ = propName
					declaredSlotType = objType.GetObjectTypeBase().IsPropertyTypeDeclared_1_string(propName)
				}
			}
			var safeLeftType pkg_rhino.JSType
			_ = safeLeftType
			if targetType == nil || targetType.IsJavaNil_() {
				safeLeftType = pkg_rhino.Upcast_JSType(this_.unknownType)
			} else {
				safeLeftType = targetType
			}
			scope = j2gort.AsIface[pkg_jscomptype.FlowScope](scope.InferQualifiedSlot(target, qualifiedName, safeLeftType, resultType, declaredSlotType))
		}
		if updateNode != 0 {
			updateNode.SetJSType(resultType)
		}
		this_.priv_TypeInference_ensurePropertyDefined_5515(target, resultType, scope)
	default:
	}
	return scope
}

// Java method: TypeInference.ensurePropertyDefined(com.google.javascript.rhino.Node getprop, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_ensurePropertyDefined_5515(getprop pkg_rhino.Node, rightType pkg_rhino.JSType, scope pkg_jscomptype.FlowScope) {
	_ = getprop
	_ = rightType
	_ = scope
	var propName string = getprop.GetString()
	_ = propName
	var obj pkg_rhino.Node = getprop.GetFirstChild()
	_ = obj
	var nodeType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(obj)
	_ = nodeType
	var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(nodeType.RestrictByNotNullOrUndefined())
	_ = objectType
	var propCreationInConstructor bool = (obj.IsThis() && this_.priv_TypeInference_getJSType_5571(this_.containerScope.GetRootNode_0()).GetJSTypeBase().IsConstructor())
	_ = propCreationInConstructor
	if objectType == nil || objectType.IsJavaNil_() {
		this_.registry.RegisterPropertyOnType(propName, nodeType)
	} else {
		if nodeType.IsStruct() && !(objectType.GetObjectTypeBase().HasProperty_1_string(propName)) {
			var staticPropCreation bool = false
			_ = staticPropCreation
			var maybeAssignStm pkg_rhino.Node = getprop.GetGrandparent()
			_ = maybeAssignStm
			if this_.containerScope.IsGlobal() && NodeUtil_IsPrototypePropertyDeclaration(maybeAssignStm) {
				var propCreationFilename string = maybeAssignStm.GetSourceFileName()
				_ = propCreationFilename
				var ctor pkg_rhino.Node = objectType.GetOwnerFunction().GetSource()
				_ = ctor
				if (ctor != 0) && (j2gort.CheckNotNullString(ctor.GetSourceFileName()) == propCreationFilename) {
					staticPropCreation = true
				}
			}
			if !(propCreationInConstructor) && !(staticPropCreation) {
				return
			}
		}
		if this_.priv_TypeInference_ensurePropertyDeclaredHelper_5517(getprop, objectType, scope) {
			return
		}
		if !(objectType.GetObjectTypeBase().IsPropertyTypeDeclared_1_string(propName)) {
			if objectType.GetObjectTypeBase().HasProperty_1_string(propName) || !(objectType.IsInstanceType()) {
				if "prototype" == propName {
					this_.priv_TypeInference_defineDeclaredProperty_5516(objectType, propName, rightType, getprop)
				} else {
					objectType.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(propName, rightType, getprop)
				}
			} else if propCreationInConstructor {
				objectType.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(propName, rightType, getprop)
			} else {
				this_.registry.RegisterPropertyOnType(propName, pkg_rhino.Upcast_JSType(objectType))
			}
		}
	}
}

// Java method: TypeInference.defineDeclaredProperty(com.google.javascript.rhino.jstype.ObjectType objectType, java.lang.String propName, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.Node getprop)
func (this_ *TypeInference) priv_TypeInference_defineDeclaredProperty_5516(objectType pkg_rhino.ObjectType, propName string, rightType pkg_rhino.JSType, getprop pkg_rhino.Node) bool {
	_ = objectType
	_ = propName
	_ = rightType
	_ = getprop
	if (j2gort.CheckNotNullString(propName) == "prototype") && !(getprop.GetParent().IsExprResult()) {
		var functionType *pkg_rhino.FunctionType = objectType.ToMaybeFunctionType()
		_ = functionType
		if ((functionType != nil) && (functionType.GetSource() != 0)) && functionType.GetSource().IsClass() {
			this_.compiler.Report_1(JSError_Make_3(getprop.GetParent(), StaticGet_TypeInference_REASSIGN_CLASS_PROTOTYPE(), []string{}))
			return true
		}
	}
	return objectType.GetObjectTypeBase().DefineDeclaredProperty_3_string_JSType_Node(propName, rightType, getprop)
}

// Java method: TypeInference.ensurePropertyDeclaredHelper(com.google.javascript.rhino.Node getprop, com.google.javascript.rhino.jstype.ObjectType objectType, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_ensurePropertyDeclaredHelper_5517(getprop pkg_rhino.Node, objectType pkg_rhino.ObjectType, scope pkg_jscomptype.FlowScope) bool {
	_ = getprop
	_ = objectType
	_ = scope
	var rootName string = TypeInference_GetGetPropRootName(getprop)
	_ = rootName
	if rootName != j2gort.NullString {
		var declScope *TypedScope = j2gort.CastPtr[TypedScope](j2gort.AsIface[pkg_rhino.StaticTypedScope](scope.GetDeclarationScope()))
		_ = declScope
		if declScope.MayHaveQualifiedVarForRoot(rootName) && declScope.MayHaveQualifiedVarForNode(getprop) {
			var propName string = getprop.GetString()
			_ = propName
			var qName string = getprop.GetQualifiedName()
			_ = qName
			if qName != j2gort.NullString {
				var var_ *TypedVar = declScope.LookupTypedVar(qName)
				_ = var_
				if (var_ != nil) && !(var_.IsTypeInferred()) {
					if (j2gort.CheckNotNullString(propName) == "prototype") || (!(objectType.GetObjectTypeBase().HasOwnProperty_1_string(propName)) && (!(objectType.IsInstanceType()) || (var_.IsExtern() && !(objectType.GetObjectTypeBase().IsNativeObjectType())))) {
						return this_.priv_TypeInference_defineDeclaredProperty_5516(objectType, propName, var_.GetType(), getprop)
					}
				}
			}
		}
	}
	return false
}

// Java method: TypeInference.traverseDeclaration(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseDeclaration_5518(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	{
		var declarationChild pkg_rhino.Node = n.GetFirstChild()
		_ = declarationChild
		for ; declarationChild != 0; declarationChild = declarationChild.GetNext() {
			scope = this_.priv_TypeInference_traverseDeclarationChild_5519(declarationChild, scope)
		}
	}
	return scope
}

// Java method: TypeInference.traverseDeclarationChild(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseDeclarationChild_5519(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	if n.IsName() {
		return this_.priv_TypeInference_traverseName_5522(n, scope)
	}
	if !(n.IsDestructuringLhs()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	scope = this_.priv_TypeInference_traverse_5500(n.GetSecondChild(), scope)
	return this_.priv_TypeInference_traverseDestructuringPattern_5520(n.GetFirstChild(), scope, this_.priv_TypeInference_getJSType_5571(n.GetSecondChild()), StaticGet_TypeInference_AssignmentType_DECLARATION())
}

// Java method: TypeInference.traverseDestructuringPattern(com.google.javascript.rhino.Node pattern, com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.jstype.JSType patternType, com.google.javascript.jscomp.TypeInference.AssignmentType assignmentType)
func (this_ *TypeInference) priv_TypeInference_traverseDestructuringPattern_5520(pattern pkg_rhino.Node, scope pkg_jscomptype.FlowScope, patternType pkg_rhino.JSType, assignmentType TypeInference_AssignmentType) pkg_jscomptype.FlowScope {
	_ = pattern
	_ = scope
	_ = patternType
	_ = assignmentType
	return this_.priv_TypeInference_traverseDestructuringPatternHelper_5521(pattern, scope, patternType, NewLambda_TypeInference_TypeDeclaringCallback(func(flowScope pkg_jscomptype.FlowScope, targetNode pkg_rhino.Node, targetType pkg_rhino.JSType) pkg_jscomptype.FlowScope {
		_ = flowScope
		_ = targetNode
		_ = targetType
		if targetType != nil && !targetType.IsJavaNil_() {
			targetType = targetType
		} else {
			targetType = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
		return this_.priv_TypeInference_updateScopeForAssignment_5513(flowScope, targetNode, targetType, assignmentType)
	}))
}

// Java method: TypeInference.traverseDestructuringPatternHelper(com.google.javascript.rhino.Node pattern, com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.jstype.JSType patternType, com.google.javascript.jscomp.TypeInference.TypeDeclaringCallback declarer)
func (this_ *TypeInference) priv_TypeInference_traverseDestructuringPatternHelper_5521(pattern pkg_rhino.Node, scope pkg_jscomptype.FlowScope, patternType pkg_rhino.JSType, declarer TypeInference_TypeDeclaringCallback) pkg_jscomptype.FlowScope {
	_ = pattern
	_ = scope
	_ = patternType
	_ = declarer
	if !(pattern.IsDestructuringPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(pattern))
	}
	j2gort.CheckNotNullIface(patternType)
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_JSType_Node(this_.registry, patternType, pattern)))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			target := j2gort.CastPtr[DestructuredTarget](rawElem_9)
			_ = target
			if target.HasComputedProperty() {
				scope = this_.priv_TypeInference_traverse_5500(target.GetComputedProperty().GetFirstChild(), scope)
			}
			var targetNode pkg_rhino.Node = target.GetNode()
			_ = targetNode
			if targetNode.IsDestructuringPattern() {
				if target.HasDefaultValue() {
					this_.priv_TypeInference_traverse_5500(target.GetDefaultValue(), scope)
				}
				var targetType pkg_rhino.JSType = target.InferType()
				_ = targetType
				if targetType != nil && !targetType.IsJavaNil_() {
					targetType = targetType
				} else {
					targetType = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
				}
				scope = this_.priv_TypeInference_traverseDestructuringPatternHelper_5521(targetNode, scope, targetType, declarer)
			} else {
				scope = this_.priv_TypeInference_traverse_5500(targetNode, scope)
				if target.HasDefaultValue() {
					this_.priv_TypeInference_traverse_5500(target.GetDefaultValue(), scope)
				}
				scope = declarer.DeclareTypeInScope(scope, targetNode, target.InferType())
			}
		}
	}
	pattern.SetJSType(patternType)
	return scope
}

// Java method: TypeInference.traverseName(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseName_5522(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var varName string = n.GetString()
	_ = varName
	var value pkg_rhino.Node = n.GetFirstChild()
	_ = value
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if value != 0 {
		scope = this_.priv_TypeInference_traverse_5500(value, scope)
		return this_.priv_TypeInference_updateScopeForAssignment_5513(scope, n, this_.priv_TypeInference_getJSType_5571(value), StaticGet_TypeInference_AssignmentType_DECLARATION())
	}
	if NodeUtil_IsNameDeclaration(n.GetParent()) && TypeInference_IsPossibleMixinApplication(n, pkg_rhino.Node(0)) {
		TypeInference_AddMissingInterfaceProperties(type_)
	}
	if n.GetParent().IsLet() {
		var resultType pkg_rhino.JSType
		_ = resultType
		if type_ != nil && !type_.IsJavaNil_() {
			resultType = type_
		} else {
			resultType = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
		}
		scope = this_.priv_TypeInference_updateScopeForAssignment_5513(scope, n, resultType, StaticGet_TypeInference_AssignmentType_DECLARATION())
		type_ = resultType
	} else {
		var var_ pkg_rhino.StaticTypedSlot = j2gort.AsIface[pkg_rhino.StaticTypedSlot](scope.GetSlot(varName))
		_ = var_
		if var_ != nil && !var_.IsJavaNil_() {
			var isInferred bool = var_.IsTypeInferred()
			_ = isInferred
			var unflowable bool = (isInferred && this_.priv_TypeInference_isUnflowable_5570(TypeInference_GetDeclaredVar(scope, varName)))
			_ = unflowable
			var maybeOuterVar *TypedVar
			_ = maybeOuterVar
			if isInferred && this_.containerScope.IsLocal() {
				maybeOuterVar = j2gort.CastPtr[TypedScope](this_.containerScope.GetParent()).LookupTypedVar(varName)
			} else {
				maybeOuterVar = nil
			}
			var nonLocalInferredSlot bool = (javalang.ObjectsEquals(j2gort.CheckNotNull(var_), j2gort.PtrToAny(maybeOuterVar)) && !(maybeOuterVar.IsMarkedAssignedExactlyOnce()))
			_ = nonLocalInferredSlot
			if !(unflowable) && !(nonLocalInferredSlot) {
				type_ = var_.GetType()
				if type_ == nil || type_.IsJavaNil_() {
					type_ = pkg_rhino.Upcast_JSType(this_.unknownType)
				}
			}
		}
	}
	n.SetJSType(type_)
	return scope
}

// Java method: TypeInference.traverseImportSpec(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node spec)
func (this_ *TypeInference) priv_TypeInference_traverseImportSpec_5523(scope pkg_jscomptype.FlowScope, spec pkg_rhino.Node) pkg_jscomptype.FlowScope {
	_ = scope
	_ = spec
	var exportedName pkg_rhino.Node = spec.GetFirstChild()
	_ = exportedName
	var localName pkg_rhino.Node = spec.GetSecondChild()
	_ = localName
	scope = this_.priv_TypeInference_traverse_5500(localName, scope)
	exportedName.SetJSType(localName.GetJSType())
	return scope
}

// Java static method: TypeInference.getBigIntPresence(com.google.javascript.rhino.jstype.JSType type)
func TypeInference_GetBigIntPresence(type_ pkg_rhino.JSType) TypeInference_BigIntPresence {
	clinit_TypeInference()
	_ = type_
	if type_.GetJSTypeBase().IsOnlyBigInt() {
		return StaticGet_TypeInference_BigIntPresence_ALL_BIGINT()
	}
	var typeAsEnumElement *pkg_rhino.EnumElementType = type_.ToMaybeEnumElementType()
	_ = typeAsEnumElement
	if typeAsEnumElement != nil {
		return TypeInference_GetBigIntPresence(typeAsEnumElement.GetPrimitiveType())
	}
	var typeAsUnion *pkg_rhino.UnionType = type_.ToMaybeUnionType()
	_ = typeAsUnion
	if typeAsUnion != nil {
		var containsBigInt bool = false
		_ = containsBigInt
		var containsNumber bool = false
		_ = containsNumber
		var containsOther bool = false
		_ = containsOther
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = typeAsUnion.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var alternate pkg_rhino.JSType = typeAsUnion.GetAlternateAt(i)
				_ = alternate
				if TypeInference_GetBigIntPresence(alternate) != StaticGet_TypeInference_BigIntPresence_NO_BIGINT() {
					containsBigInt = true
				} else if alternate.GetJSTypeBase().IsNumber() && !(alternate.GetJSTypeBase().IsUnknownType()) {
					containsNumber = true
				} else {
					containsOther = true
				}
			}
		}
		if containsBigInt {
			if containsOther {
				return StaticGet_TypeInference_BigIntPresence_BIGINT_OR_OTHER()
			} else if containsNumber {
				return StaticGet_TypeInference_BigIntPresence_BIGINT_OR_NUMBER()
			} else {
				return StaticGet_TypeInference_BigIntPresence_ALL_BIGINT()
			}
		}
	}
	return StaticGet_TypeInference_BigIntPresence_NO_BIGINT()
}

// Java method: TypeInference.traverseUnaryPlus(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseUnaryPlus_5524(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	var bigintPresenceInOperand TypeInference_BigIntPresence = TypeInference_GetBigIntPresence(this_.priv_TypeInference_getJSType_5571(n.GetFirstChild()))
	_ = bigintPresenceInOperand
	if bigintPresenceInOperand != StaticGet_TypeInference_BigIntPresence_NO_BIGINT() {
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE()))
	} else {
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))
	}
	return scope
}

// Java method: TypeInference.traverseBigIntCompatibleUnaryOperator(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseBigIntCompatibleUnaryOperator_5525(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	switch (TypeInference_GetBigIntPresence(this_.priv_TypeInference_getJSType_5571(n.GetFirstChild()))).Ordinal() {
	case int32(1) /* BigIntPresence.ALL_BIGINT */ :
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE()))
	case int32(0) /* BigIntPresence.NO_BIGINT */ :
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))
	case int32(2) /* BigIntPresence.BIGINT_OR_NUMBER */, int32(3) /* BigIntPresence.BIGINT_OR_OTHER */ :
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER()))
	}
	return scope
}

// Java method: TypeInference.traverseBigIntCompatibleBinaryOperator(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseBigIntCompatibleBinaryOperator_5526(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	var leftBigIntPresence TypeInference_BigIntPresence = TypeInference_GetBigIntPresence(this_.priv_TypeInference_getJSType_5571(n.GetFirstChild()))
	_ = leftBigIntPresence
	var rightBigIntPresence TypeInference_BigIntPresence = TypeInference_GetBigIntPresence(this_.priv_TypeInference_getJSType_5571(n.GetLastChild()))
	_ = rightBigIntPresence
	if leftBigIntPresence != rightBigIntPresence {
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE()))
	} else {
		switch (leftBigIntPresence).Ordinal() {
		case int32(0) /* BigIntPresence.NO_BIGINT */ :
			n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))
		case int32(1) /* BigIntPresence.ALL_BIGINT */ :
			n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE()))
		case int32(2) /* BigIntPresence.BIGINT_OR_NUMBER */ :
			n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER()))
		case int32(3) /* BigIntPresence.BIGINT_OR_OTHER */ :
			n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE()))
		}
	}
	return scope
}

// Java method: TypeInference.traverseClass(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseClass_5527(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverse_5500(n.GetSecondChild(), scope)
	var classMembers pkg_rhino.Node = NodeUtil_GetClassMembers(n)
	_ = classMembers
	{
		var member pkg_rhino.Node = classMembers.GetFirstChild()
		_ = member
		for ; member != 0; member = member.GetNext() {
			if member.IsComputedProp() || member.IsComputedFieldDef() {
				scope = this_.priv_TypeInference_traverse_5500(member.GetFirstChild(), scope)
			}
		}
	}
	{
		var member_2 pkg_rhino.Node = classMembers.GetFirstChild()
		_ = member_2
		for ; member_2 != 0; member_2 = member_2.GetNext() {
			scope = this_.priv_TypeInference_traverseClassMemberRhs_5528(member_2, scope)
		}
	}
	return scope
}

// Java method: TypeInference.traverseClassMemberRhs(com.google.javascript.rhino.Node member, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseClassMemberRhs_5528(member pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = member
	_ = scope
	switch member.TokenSwitchOrd_() {
	case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		var rhs pkg_rhino.Node = TypeInference_GetRhsOfField(member)
		_ = rhs
		if rhs != 0 {
			var computedFieldDefTypedScope *TypedScope = this_.scopeCreator.CreateScope_1(member)
			_ = computedFieldDefTypedScope
			var computedFieldDefFlowScope pkg_jscomptype.FlowScope = scope.WithSyntacticScope(pkg_rhino.ToIface_StaticTypedScope(computedFieldDefTypedScope))
			_ = computedFieldDefFlowScope
			var rhsScope pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverse_5500(rhs, computedFieldDefFlowScope).WithSyntacticScope(j2gort.AsIface[pkg_rhino.StaticTypedScope](scope.GetDeclarationScope()))
			_ = rhsScope
			if member.IsStaticMember() {
				return rhsScope
			}
		}
		return scope
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(90) /* Token.BLOCK */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		return scope
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	panic("unreachable switch expression")
}

// Java static method: TypeInference.getRhsOfField(com.google.javascript.rhino.Node fieldNode)
func TypeInference_GetRhsOfField(fieldNode pkg_rhino.Node) pkg_rhino.Node {
	clinit_TypeInference()
	_ = fieldNode
	switch fieldNode.TokenSwitchOrd_() {
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if fieldNode.HasOneChild() {
			return fieldNode.GetFirstChild()
		}
		return pkg_rhino.Node(0)
	case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		if fieldNode.HasTwoChildren() {
			return fieldNode.GetSecondChild()
		}
		return pkg_rhino.Node(0)
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeInference.traverseArrayLiteral(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseArrayLiteral_5529(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	n.SetJSType(pkg_rhino.ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), []pkg_rhino.JSType{this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())})))
	return scope
}

// Java method: TypeInference.traverseObjectLiteral(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseObjectLiteral_5530(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	j2gort.CheckNotNullIface(type_)
	var spreadOperatorSeen bool = false
	_ = spreadOperatorSeen
	{
		var name pkg_rhino.Node = n.GetFirstChild()
		_ = name
		for ; name != 0; name = name.GetNext() {
			{
				ba_10 := spreadOperatorSeen
				bb_11 := name.IsSpread()
				spreadOperatorSeen = ba_10 || bb_11
			}
			scope = this_.priv_TypeInference_traverseChildren_5549(name, scope)
		}
	}
	var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(type_)
	_ = objectType
	if ((objectType == nil || objectType.IsJavaNil_()) || n.GetBooleanProp(pkg_rhino.StaticGet_Node_REFLECTED_OBJECT())) || objectType.GetObjectTypeBase().IsEnumType() {
		return scope
	}
	var qObjName string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	_ = qObjName
	{
		var key pkg_rhino.Node = n.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if key.IsComputedProp() {
				continue
			}
			if key.IsSpread() {
				n.SetJSType(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()))
				break
			}
			var memberName string = NodeUtil_GetObjectOrClassLitKeyName(key)
			_ = memberName
			if memberName != j2gort.NullString {
				var rawValueType pkg_rhino.JSType = key.GetFirstChild().GetJSType()
				_ = rawValueType
				var valueType pkg_rhino.JSType = TypeCheck_GetObjectLitKeyTypeFromValueType(key, rawValueType)
				_ = valueType
				if valueType == nil || valueType.IsJavaNil_() {
					valueType = pkg_rhino.Upcast_JSType(this_.unknownType)
				}
				if !(spreadOperatorSeen) || !(n.GetJSType().GetJSTypeBase().IsNativeObjectType()) {
					objectType.GetObjectTypeBase().DefineInferredProperty_3_string_JSType_Node(memberName, valueType, key)
				}
				if (qObjName != j2gort.NullString) && key.IsStringKey() {
					var qKeyName string = (j2gort.StrValString(qObjName) + "." + j2gort.StrValString(memberName))
					_ = qKeyName
					var var_ *TypedVar = TypeInference_GetDeclaredVar(scope, qKeyName)
					_ = var_
					var oldType pkg_rhino.JSType
					_ = oldType
					if var_ == nil {
						oldType = nil
					} else {
						oldType = var_.GetType()
					}
					if (var_ != nil) && var_.IsTypeInferred() {
						var iifeArg_2 pkg_rhino.
							JSType
					iifeSwitch_2:
						switch {
						default:

							if oldType == nil || oldType.IsJavaNil_() {
								iifeArg_2 = valueType
								break iifeSwitch_2
							}
							iifeArg_2 = oldType.GetLeastSupertype(oldType)
							break iifeSwitch_2
						}
						var_.SetType(iifeArg_2)
					}
					scope = j2gort.AsIface[pkg_jscomptype.FlowScope](scope.InferQualifiedSlot(key, qKeyName, func() pkg_rhino.JSType {
						if oldType == nil || oldType.IsJavaNil_() {
							return pkg_rhino.Upcast_JSType(this_.unknownType)
						}
						return oldType
					}(), valueType, false))
				}
			} else {
				n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
			}
		}
	}
	return scope
}

// Java method: TypeInference.traverseAdd(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseAdd_5531(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var right pkg_rhino.Node = left.GetNext()
	_ = right
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	var leftType pkg_rhino.JSType = left.GetJSType()
	_ = leftType
	var rightType pkg_rhino.JSType = right.GetJSType()
	_ = rightType
	var type_ pkg_rhino.JSType = pkg_rhino.Upcast_JSType(this_.unknownType)
	_ = type_
	if (leftType != nil && !leftType.IsJavaNil_()) && (rightType != nil && !rightType.IsJavaNil_()) {
		var leftIsUnknown bool = leftType.GetJSTypeBase().IsUnknownType()
		_ = leftIsUnknown
		var rightIsUnknown bool = rightType.GetJSTypeBase().IsUnknownType()
		_ = rightIsUnknown
		if (!(leftIsUnknown) && leftType.GetJSTypeBase().IsString()) || (!(rightIsUnknown) && rightType.GetJSTypeBase().IsString()) {
			type_ = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
		} else if (TypeInference_GetBigIntPresence(leftType) != StaticGet_TypeInference_BigIntPresence_NO_BIGINT()) || (TypeInference_GetBigIntPresence(rightType) != StaticGet_TypeInference_BigIntPresence_NO_BIGINT()) {
			return this_.priv_TypeInference_traverseBigIntCompatibleBinaryOperator_5526(n, scope)
		} else if leftIsUnknown || rightIsUnknown {
			type_ = pkg_rhino.Upcast_JSType(this_.unknownType)
		} else if this_.priv_TypeInference_isAddedAsNumber_5532(leftType) && this_.priv_TypeInference_isAddedAsNumber_5532(rightType) {
			type_ = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
		} else {
			type_ = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING())
		}
	}
	n.SetJSType(type_)
	if n.IsAssignAdd() {
		scope = this_.priv_TypeInference_updateScopeForAssignment_5514(scope, left, type_, pkg_rhino.Node(0), StaticGet_TypeInference_AssignmentType_ASSIGN())
	}
	return scope
}

// Java method: TypeInference.isAddedAsNumber(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeInference) priv_TypeInference_isAddedAsNumber_5532(type_ pkg_rhino.JSType) bool {
	_ = type_
	return type_.GetJSTypeBase().IsSubtypeOf_1(this_.numberAdditionSupertype)
}

// Java method: TypeInference.traverseHook(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseHook_5533(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var condition pkg_rhino.Node = n.GetFirstChild()
	_ = condition
	var trueNode pkg_rhino.Node = condition.GetNext()
	_ = trueNode
	var falseNode pkg_rhino.Node = n.GetLastChild()
	_ = falseNode
	scope = this_.priv_TypeInference_traverse_5500(condition, scope)
	var trueScope pkg_jscomptype.FlowScope = this_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(condition, scope, pkg_rhino.StaticGet_Outcome_TRUE())
	_ = trueScope
	var falseScope pkg_jscomptype.FlowScope = this_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(condition, scope, pkg_rhino.StaticGet_Outcome_FALSE())
	_ = falseScope
	this_.priv_TypeInference_traverse_5500(trueNode, trueScope)
	this_.priv_TypeInference_traverse_5500(falseNode, falseScope)
	var trueType pkg_rhino.JSType = trueNode.GetJSType()
	_ = trueType
	var falseType pkg_rhino.JSType = falseNode.GetJSType()
	_ = falseType
	if (trueType != nil && !trueType.IsJavaNil_()) && (falseType != nil && !falseType.IsJavaNil_()) {
		n.SetJSType(trueType.GetLeastSupertype(falseType))
	} else {
		n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
	}
	return scope
}

// Java method: TypeInference.traverseNullishCoalesce(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseNullishCoalesce_5534(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	if !(n.IsNullishCoalesce() || n.IsAssignNullishCoalesce()) {
		j2gort.CheckArgument(false)
	}
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var right pkg_rhino.Node = n.GetLastChild()
	_ = right
	scope = this_.priv_TypeInference_traverse_5500(left, scope)
	var rightScope pkg_jscomptype.FlowScope = this_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(left, scope, pkg_rhino.StaticGet_Outcome_NULLISH())
	_ = rightScope
	var scopeAfterTraverseRight pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverse_5500(right, rightScope)
	_ = scopeAfterTraverseRight
	var leftType pkg_rhino.JSType = left.GetJSType()
	_ = leftType
	var rightType pkg_rhino.JSType = right.GetJSType()
	_ = rightType
	var type_ pkg_rhino.JSType = pkg_rhino.Upcast_JSType(this_.unknownType)
	_ = type_
	if leftType != nil && !leftType.IsJavaNil_() {
		if !(leftType.IsNullable()) && !(leftType.IsVoidable()) {
			type_ = leftType
		} else if rightType != nil && !rightType.IsJavaNil_() {
			type_ = this_.registry.CreateUnionType_2(leftType.RestrictByNotNullOrUndefined(), rightType)
			scope = this_.Join(scope, scopeAfterTraverseRight)
			if n.IsAssignNullishCoalesce() {
				scope = this_.priv_TypeInference_updateScopeForAssignment_5513(scope, left, type_, StaticGet_TypeInference_AssignmentType_ASSIGN())
			}
		}
	}
	n.SetJSType(type_)
	return scope
}

// Java method: TypeInference.setCallNodeTypeAfterChildrenTraversed(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scopeAfterChildren)
func (this_ *TypeInference) priv_TypeInference_setCallNodeTypeAfterChildrenTraversed_5535(n pkg_rhino.Node, scopeAfterChildren pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scopeAfterChildren
	if (n.IsCall() && NodeUtil_IsExpressionResultUsed(n)) && ModuleImportResolver_IsGoogModuleDependencyCall(n) {
		var calleeNode pkg_rhino.Node = n.GetFirstChild()
		_ = calleeNode
		if calleeNode.MatchesQualifiedName_1_string("goog.requireDynamic") {
			n.SetJSType(Promises_WrapInIThenable(this_.registry, this_.priv_TypeInference_getGoogModuleDependencyCallResultType_5536(n)))
		} else {
			n.SetJSType(this_.priv_TypeInference_getGoogModuleDependencyCallResultType_5536(n))
		}
		return scopeAfterChildren
	}
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var functionType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(left).RestrictByNotNullOrUndefined()
	_ = functionType
	if left.IsSuper() {
		return this_.priv_TypeInference_traverseInstantiation_5547(n, functionType, scopeAfterChildren)
	} else if functionType.GetJSTypeBase().IsFunctionType() {
		var fnType *pkg_rhino.FunctionType = functionType.ToMaybeFunctionType()
		_ = fnType
		n.SetJSType(fnType.GetReturnType())
		this_.priv_TypeInference_backwardsInferenceFromCallSite_5539(n, fnType, scopeAfterChildren)
	} else if functionType.GetJSTypeBase().EqualsJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())) {
		n.SetJSType(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE()))
	} else if !(j2gort.IfaceEquals[pkg_rhino.JSType](left.GetJSType(), nil)) && left.GetJSType().GetJSTypeBase().IsUnknownType() {
		n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
	}
	return scopeAfterChildren
}

// Java method: TypeInference.getGoogModuleDependencyCallResultType(com.google.javascript.rhino.Node callNode)
func (this_ *TypeInference) priv_TypeInference_getGoogModuleDependencyCallResultType_5536(callNode pkg_rhino.Node) pkg_rhino.JSType {
	_ = callNode
	var moduleId string = callNode.GetSecondChild().GetString()
	_ = moduleId
	var module *Module = this_.compiler.GetModuleMap().GetClosureModule(moduleId)
	_ = module
	if module == nil {
		return pkg_rhino.Upcast_JSType(this_.unknownType)
	}
	var name ScopedName = this_.moduleImportResolver.GetClosureNamespaceTypeFromCall(callNode)
	_ = name
	if name != nil && !name.IsJavaNil_() {
		var otherModuleScope *TypedScope = j2gort.CastPtr[TypedScope](this_.scopeCreator.GetNodeToScopeMapper().Apply(j2gort.BoxNode(name.GetScopedNameBase().GetScopeRoot())))
		_ = otherModuleScope
		var otherVar *TypedVar
		_ = otherVar
		if otherModuleScope != nil {
			otherVar = j2gort.CastPtr[TypedVar](otherModuleScope.GetSlot(name.GetScopedNameBase().GetName()))
		} else {
			otherVar = nil
		}
		if otherVar != nil {
			if !(j2gort.IfaceEquals[pkg_rhino.JSType](otherVar.GetType(), nil)) {
				return otherVar.GetType()
			} else {
				return pkg_rhino.Upcast_JSType(this_.unknownType)
			}
		}
		if (otherModuleScope != nil) && (module.Metadata().ModuleType() == StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE()) {
			return otherModuleScope.GetTypeThroughNamespace(moduleId)
		}
	}
	return pkg_rhino.Upcast_JSType(this_.unknownType)
}

// Java method: TypeInference.tightenTypesAfterAssertions(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node callNode)
func (this_ *TypeInference) priv_TypeInference_tightenTypesAfterAssertions_5537(scope pkg_jscomptype.FlowScope, callNode pkg_rhino.Node) pkg_jscomptype.FlowScope {
	_ = scope
	_ = callNode
	var left pkg_rhino.Node = callNode.GetFirstChild()
	_ = left
	var firstParam pkg_rhino.Node = left.GetNext()
	_ = firstParam
	if firstParam == 0 {
		return scope
	}
	var assertionFunctionSpec CodingConvention_AssertionFunctionSpec = this_.assertionFunctionLookup.LookupByCallee(left)
	_ = assertionFunctionSpec
	if assertionFunctionSpec == nil || assertionFunctionSpec.IsJavaNil_() {
		return scope
	}
	var assertedNode pkg_rhino.Node = assertionFunctionSpec.GetCodingConvention_AssertionFunctionSpecBase().GetAssertedArg(firstParam)
	_ = assertedNode
	if assertedNode == 0 {
		return scope
	}
	switch (assertionFunctionSpec.GetAssertionKind()).Ordinal() {
	case int32(0) /* AssertionKind.TRUTHY */ :
		scope = this_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(assertedNode, scope, pkg_rhino.StaticGet_Outcome_TRUE())
		var truthyType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(assertedNode).RestrictByNotNullOrUndefined()
		_ = truthyType
		callNode.SetJSType(truthyType)
	case int32(1) /* AssertionKind.MATCHES_RETURN_TYPE */ :
		var callType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(left.GetJSType())
		_ = callType
		var assertedType pkg_rhino.JSType
		_ = assertedType
		if callType != nil {
			assertedType = callType.GetReturnType()
		} else {
			assertedType = pkg_rhino.Upcast_JSType(this_.unknownType)
		}
		var type_ pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(assertedNode)
		_ = type_
		var narrowed pkg_rhino.JSType = nil
		_ = narrowed
		if assertedType.GetJSTypeBase().IsUnknownType() || type_.GetJSTypeBase().IsUnknownType() {
			narrowed = assertedType
		} else {
			narrowed = type_.GetJSTypeBase().GetGreatestSubtype(assertedType)
		}
		callNode.SetJSType(narrowed)
		if assertedNode.IsQualifiedName() && type_.GetJSTypeBase().DiffersFrom(narrowed) {
			scope = this_.priv_TypeInference_narrowScope_5538(scope, assertedNode, narrowed)
		}
	}
	return scope
}

// Java method: TypeInference.narrowScope(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType narrowed)
func (this_ *TypeInference) priv_TypeInference_narrowScope_5538(scope pkg_jscomptype.FlowScope, node pkg_rhino.Node, narrowed pkg_rhino.JSType) pkg_jscomptype.FlowScope {
	_ = scope
	_ = node
	_ = narrowed
	if node.IsThis() {
		return scope
	}
	if node.IsGetProp() {
		var lfs *LinkedFlowScope
		_ = lfs
		var ok_12 bool
		lfs, ok_12 = j2gort.TryCast[*LinkedFlowScope](scope)
		if ok_12 && (lfs.HasGlobalFunctionScope()) {
			return scope
		}
		return j2gort.AsIface[pkg_jscomptype.FlowScope](scope.InferQualifiedSlot(node, node.GetQualifiedName(), this_.priv_TypeInference_getJSType_5571(node), narrowed, false))
	}
	return this_.priv_TypeInference_redeclareSimpleVar_5569(scope, node, narrowed)
}

// Java method: TypeInference.backwardsInferenceFromCallSite(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType fnType, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_backwardsInferenceFromCallSite_5539(n pkg_rhino.Node, fnType *pkg_rhino.FunctionType, scope pkg_jscomptype.FlowScope) {
	_ = n
	_ = fnType
	_ = scope
	var updatedFnType bool = this_.priv_TypeInference_inferTemplatedTypesForCall_5545(n, fnType, scope)
	_ = updatedFnType
	if updatedFnType {
		fnType = n.GetFirstChild().GetJSType().ToMaybeFunctionType()
	}
	this_.priv_TypeInference_updateTypeOfArguments_5541(n, fnType)
	this_.priv_TypeInference_updateBind_5540(n)
}

// Java method: TypeInference.updateBind(com.google.javascript.rhino.Node n)
func (this_ *TypeInference) priv_TypeInference_updateBind_5540(n pkg_rhino.Node) {
	_ = n
	var bind *CodingConvention_Bind = this_.compiler.GetCodingConvention().DescribeFunctionBind(n, true)
	_ = bind
	if bind == nil {
		return
	}
	var target pkg_rhino.Node = bind.target
	_ = target
	var callTargetFn *pkg_rhino.FunctionType = this_.priv_TypeInference_getJSType_5571(target).RestrictByNotNullOrUndefined().ToMaybeFunctionType()
	_ = callTargetFn
	if callTargetFn == nil {
		return
	}
	if (bind.thisValue != 0) && target.IsFunction() {
		var thisType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(bind.thisValue)
		_ = thisType
		if (!(j2gort.IfaceEquals[pkg_rhino.ObjectType](thisType.GetJSTypeBase().ToObjectType(), nil)) && !(thisType.GetJSTypeBase().IsUnknownType())) && callTargetFn.GetTypeOfThis().GetJSTypeBase().IsUnknownType() {
			var fb_ pkg_rhino.FunctionType_Builder
			callTargetFn.InitBuilder_(&fb_)
			(&fb_).WithTypeOfThis(pkg_rhino.Upcast_JSType(thisType.GetJSTypeBase().ToObjectType()))
			(&fb_).WithSourceNode(pkg_rhino.Node(0))
			callTargetFn = (&fb_).BuildAndResolve()
			target.SetJSType(pkg_rhino.ToIface_JSType(callTargetFn))
		}
	}
	var returnType pkg_rhino.JSType = pkg_rhino.ToIface_JSType(callTargetFn.GetBindReturnType((bind.GetBoundParameterCount() + int32(1))))
	_ = returnType
	var bindType *pkg_rhino.FunctionType = this_.priv_TypeInference_getJSType_5571(n.GetFirstChild()).ToMaybeFunctionType()
	_ = bindType
	if (bindType != nil) && (n.GetFirstChild() != target) {
		var fb_1_ pkg_rhino.FunctionType_Builder
		bindType.InitBuilder_(&fb_1_)
		(&fb_1_).WithReturnType_1(returnType)
		(&fb_1_).WithSourceNode(pkg_rhino.Node(0))
		bindType = (&fb_1_).BuildAndResolve()
		n.GetFirstChild().SetJSType(pkg_rhino.ToIface_JSType(bindType))
	}
	n.SetJSType(returnType)
}

// Java method: TypeInference.updateTypeOfArguments(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType fnType)
func (this_ *TypeInference) priv_TypeInference_updateTypeOfArguments_5541(n pkg_rhino.Node, fnType *pkg_rhino.FunctionType) {
	_ = n
	_ = fnType
	if !(NodeUtil_IsInvocation(n)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var parameters []pkg_rhino.FunctionType_Parameter = fnType.GetParameterArray()
	_ = parameters
	var paramLen int32 = int32(len(parameters))
	_ = paramLen
	var paramIdx int32 = int32(0)
	_ = paramIdx
	var isTaggedTemplate bool = n.IsTaggedTemplateLit()
	_ = isTaggedTemplate
	if isTaggedTemplate {
		if paramLen == int32(0) {
			return
		}
		paramIdx = int32(1)
	}
	var argCursor pkg_rhino.Node = NodeUtil_FirstInvocationArgCursor(n, isTaggedTemplate)
	_ = argCursor
	for (paramIdx < paramLen) && (argCursor != 0) {
		var iArgument pkg_rhino.Node = NodeUtil_InvocationArgFromCursor(argCursor, isTaggedTemplate)
		_ = iArgument
		argCursor = NodeUtil_NextInvocationArgCursor(argCursor, isTaggedTemplate)
		var iArgumentType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(iArgument)
		_ = iArgumentType
		var iParameter pkg_rhino.FunctionType_Parameter = parameters[func() int32 { old_13 := paramIdx; nw_14 := old_13 + 1; paramIdx = nw_14; return old_13 }()]
		_ = iParameter
		var iParameterType pkg_rhino.JSType
		_ = iParameterType
		if !(j2gort.IfaceEquals[pkg_rhino.JSType](iParameter.GetJSType(), nil)) {
			iParameterType = iParameter.GetJSType()
		} else {
			iParameterType = pkg_rhino.Upcast_JSType(this_.unknownType)
		}
		TypeInference_InferPropertyTypesToMatchConstraint(iArgumentType, iParameterType)
		var restrictedParameter *pkg_rhino.FunctionType = nil
		_ = restrictedParameter
		if iParameterType.GetJSTypeBase().IsUnionType() {
			var union *pkg_rhino.UnionType = iParameterType.ToMaybeUnionType()
			_ = union
			{
				var i int32 = int32(0)
				_ = i
				var altLen int32 = union.GetAlternatesCount()
				_ = altLen
				for ; i < altLen; i = i + 1 {
					var alternative pkg_rhino.JSType = union.GetAlternateAt(i)
					_ = alternative
					if alternative.GetJSTypeBase().IsFunctionType() {
						restrictedParameter = alternative.ToMaybeFunctionType()
						break
					}
				}
			}
		} else {
			restrictedParameter = iParameterType.ToMaybeFunctionType()
		}
		if ((restrictedParameter != nil) && iArgument.IsFunction()) && iArgumentType.GetJSTypeBase().IsFunctionType() {
			var argFnType *pkg_rhino.FunctionType = iArgumentType.ToMaybeFunctionType()
			_ = argFnType
			var argJsdoc *pkg_rhino.JSDocInfo = iArgument.GetJSDocInfo()
			_ = argJsdoc
			var declared bool = (((argJsdoc != nil) && argJsdoc.ContainsDeclaration()) || NodeUtil_FunctionHasInlineJsdocs(iArgument))
			_ = declared
			iArgument.SetJSType(pkg_rhino.ToIface_JSType(this_.priv_TypeInference_matchFunction_5542(restrictedParameter, argFnType, declared)))
		}
	}
}

// Java method: TypeInference.matchFunction(com.google.javascript.rhino.jstype.FunctionType expectedType, com.google.javascript.rhino.jstype.FunctionType currentType, boolean declared)
func (this_ *TypeInference) priv_TypeInference_matchFunction_5542(expectedType *pkg_rhino.FunctionType, currentType *pkg_rhino.FunctionType, declared bool) *pkg_rhino.FunctionType {
	_ = expectedType
	_ = currentType
	_ = declared
	if declared {
		if currentType.GetTypeOfThis().GetJSTypeBase().IsUnknownType() && !(expectedType.GetTypeOfThis().GetJSTypeBase().IsUnknownType()) {
			var fb_ pkg_rhino.FunctionType_Builder
			currentType.InitBuilder_(&fb_)
			(&fb_).WithTypeOfThis(expectedType.GetTypeOfThis())
			return (&fb_).BuildAndResolve()
		}
	} else {
		if currentType.GetMaxArity() <= expectedType.GetMaxArity() {
			return expectedType
		}
	}
	return currentType
}

// Java method: TypeInference.buildTypeVariables(java.util.Map<com.google.javascript.rhino.jstype.TemplateType,com.google.javascript.rhino.jstype.JSType> inferredTypes)
func (this_ *TypeInference) priv_TypeInference_buildTypeVariables_5543(inferredTypes j2gort.Ext_java_util_Map) j2gort.Ext_java_util_Map {
	_ = inferredTypes
	var typeVars j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = typeVars
	{
		cur_15 := j2gort.NewMapEntryCursor(inferredTypes)
		for cur_15.NextEntry() {
			if !(j2gort.CastPtr[pkg_rhino.TemplateType](cur_15.GetKey()).IsTypeTransformation()) {
				j2gort.MapPutStr(typeVars, j2gort.CastPtr[pkg_rhino.TemplateType](cur_15.GetKey()).GetReferenceName(), j2gort.AsIface[pkg_rhino.JSType](cur_15.GetValue()))
			}
		}
	}
	return typeVars
}

// Java method: TypeInference.evaluateTypeTransformations(com.google.javascript.rhino.jstype.TemplateTypeMap templateMap, java.util.Map<com.google.javascript.rhino.jstype.TemplateType,com.google.javascript.rhino.jstype.JSType> inferredTypes, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_evaluateTypeTransformations_5544(templateMap *pkg_rhino.TemplateTypeMap, inferredTypes j2gort.Ext_java_util_Map, scope pkg_jscomptype.FlowScope) j2gort.Ext_java_util_Map {
	_ = templateMap
	_ = inferredTypes
	_ = scope
	var typeVars j2gort.Ext_java_util_Map = nil
	_ = typeVars
	var result j2gort.Ext_java_util_Map = nil
	_ = result
	var ttlObj *TypeTransformation = nil
	_ = ttlObj
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = templateMap.Size()
		_ = n
		for ; i < n; i = i + 1 {
			var type_ *pkg_rhino.TemplateType = templateMap.GetTemplateKeyAt(i)
			_ = type_
			if type_.IsTypeTransformation() {
				if ttlObj == nil {
					ttlObj = NewTypeTransformation(this_.compiler, j2gort.AsIface[pkg_rhino.StaticTypedScope](scope.GetDeclarationScope()))
					typeVars = this_.priv_TypeInference_buildTypeVariables_5543(inferredTypes)
					result = j2gort.NewExt_java_util_LinkedHashMap()
				}
				var transformedType pkg_rhino.JSType = ttlObj.Eval_2(type_.GetTypeTransformation(), j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(typeVars))
				_ = transformedType
				result.Put_2(j2gort.PtrToAny(type_), transformedType)
				j2gort.MapPutStr(typeVars, type_.GetReferenceName(), transformedType)
			}
		}
	}
	return result
}

// Java method: TypeInference.inferTemplatedTypesForCall(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType fnType, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_inferTemplatedTypesForCall_5545(n pkg_rhino.Node, fnType *pkg_rhino.FunctionType, scope pkg_jscomptype.FlowScope) bool {
	_ = n
	_ = fnType
	_ = scope
	var fnTemplateMap *pkg_rhino.TemplateTypeMap = fnType.GetTemplateTypeMap()
	_ = fnTemplateMap
	var keyCount int32 = fnTemplateMap.Size()
	_ = keyCount
	if keyCount == int32(0) {
		return false
	}
	var val_16 InvocationTemplateTypeMatcher
	InitInvocationTemplateTypeMatcher(&val_16, this_.registry, fnType, scope.GetTypeOfThis(), n)
	var matcher *InvocationTemplateTypeMatcher = &val_16
	_ = matcher
	matcher.RunMatch()
	var inferenceMap *pkg_rhino.TemplateTypeMap = nil
	_ = inferenceMap
	if (keyCount == int32(1)) && !(fnTemplateMap.GetTemplateKeyAt(int32(0)).IsTypeTransformation()) {
		var k0 *pkg_rhino.TemplateType = fnTemplateMap.GetTemplateKeyAt(int32(0))
		_ = k0
		var v0 pkg_rhino.JSType = matcher.LookupMatched(k0)
		_ = v0
		if v0 == nil || v0.IsJavaNil_() {
			v0 = pkg_rhino.Upcast_JSType(this_.unknownType)
		} else if !(v0.GetJSTypeBase().IsSubtypeOf_1(k0.GetBound())) {
			v0 = k0.GetBound()
		}
		inferenceMap = pkg_rhino.TemplateTypeMap_OfSingleBinding(k0, v0)
	} else {
		var inferred j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap_1_int32((keyCount * int32(2)))
		_ = inferred
		var hasTypeTransformation bool = false
		_ = hasTypeTransformation
		{
			var i int32 = int32(0)
			_ = i
			for ; i < keyCount; i = i + 1 {
				var key *pkg_rhino.TemplateType = fnTemplateMap.GetTemplateKeyAt(i)
				_ = key
				if key.IsTypeTransformation() {
					hasTypeTransformation = true
				}
				var matched pkg_rhino.JSType = matcher.LookupMatched(key)
				_ = matched
				if matched == nil || matched.IsJavaNil_() {
					matched = pkg_rhino.Upcast_JSType(this_.unknownType)
				} else if !(matched.GetJSTypeBase().IsSubtypeOf_1(key.GetBound())) {
					matched = key.GetBound()
				}
				inferred.Put_2(j2gort.PtrToAny(key), matched)
			}
		}
		if hasTypeTransformation {
			var typeTransformations j2gort.Ext_java_util_Map = this_.priv_TypeInference_evaluateTypeTransformations_5544(fnTemplateMap, inferred, scope)
			_ = typeTransformations
			if typeTransformations != nil {
				inferred.PutAll(typeTransformations)
			}
		}
		inferenceMap = pkg_rhino.TemplateTypeMap_FromInferenceMap(this_.registry, inferred)
	}
	var replacer *pkg_rhino.TemplateTypeReplacer = this_.registry.AcquireTemplateTypeReplacer(inferenceMap, true, true, true)
	_ = replacer
	var replacementFnType *pkg_rhino.FunctionType = nil
	_ = replacementFnType
	var madeReplacement bool = false
	_ = madeReplacement
	func() {
		defer func() {
			this_.registry.ReleaseTemplateTypeReplacer(replacer)
		}()
		replacementFnType = replacer.Replace(pkg_rhino.ToIface_JSType(fnType)).ToMaybeFunctionType()
		madeReplacement = replacer.HasMadeReplacement()
	}()
	j2gort.CheckNotNullPtr(replacementFnType)
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	callTarget.SetJSType(pkg_rhino.ToIface_JSType(replacementFnType))
	n.SetJSType(replacementFnType.GetReturnType())
	return madeReplacement
}

// Java method: TypeInference.traverseNew(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseNew_5546(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	var constructor pkg_rhino.Node = n.GetFirstChild()
	_ = constructor
	var constructorType pkg_rhino.JSType = constructor.GetJSType()
	_ = constructorType
	return this_.priv_TypeInference_traverseInstantiation_5547(n, constructorType, scope)
}

// Java method: TypeInference.traverseInstantiation(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType ctorType, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseInstantiation_5547(n pkg_rhino.Node, ctorType pkg_rhino.JSType, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = ctorType
	_ = scope
	if (ctorType == nil || ctorType.IsJavaNil_()) || ctorType.GetJSTypeBase().IsUnknownType() {
		n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
		return scope
	}
	ctorType = ctorType.RestrictByNotNullOrUndefined()
	var ctorFnType *pkg_rhino.FunctionType = ctorType.ToMaybeFunctionType()
	_ = ctorFnType
	var functionType *pkg_rhino.FunctionType
	_ = functionType
	ok_17 := (ctorFnType == nil)
	if ok_17 {
		functionType, ok_17 = j2gort.TryCast[*pkg_rhino.FunctionType](ctorType)
	}
	if ok_17 {
		ctorFnType = functionType
	}
	if (ctorFnType == nil) || !(ctorFnType.IsConstructor()) {
		n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
		return scope
	}
	this_.priv_TypeInference_backwardsInferenceFromCallSite_5539(n, ctorFnType, scope)
	var instantiatedType pkg_rhino.ObjectType = ctorFnType.GetInstanceType()
	_ = instantiatedType
	if j2gort.IfaceEquals[pkg_rhino.JSType](ctorFnType, this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE())) {
	} else if ctorFnType.HasAnyTemplateTypes() {
		if instantiatedType.GetObjectTypeBase().IsTemplatizedType() {
			instantiatedType = instantiatedType.ToMaybeTemplatizedType().GetRawType()
		}
		var inferredTypes j2gort.Ext_com_google_common_collect_ImmutableMap = NewInvocationTemplateTypeMatcher(this_.registry, ctorFnType, scope.GetTypeOfThis(), n).Match()
		_ = inferredTypes
		instantiatedType = this_.registry.CreateTemplatizedType_2_ObjectType_Ext_java_util_Map(instantiatedType, j2gort.AsIface[j2gort.Ext_java_util_Map](inferredTypes)).ToMaybeObjectType()
	}
	n.SetJSType(func() pkg_rhino.JSType {
		if instantiatedType != nil && !instantiatedType.IsJavaNil_() {
			return pkg_rhino.Upcast_JSType(instantiatedType)
		}
		return pkg_rhino.Upcast_JSType(this_.unknownType)
	}())
	return scope
}

// Java method: TypeInference.traverseAnd(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseAnd_5548(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) *TypeInference_BooleanOutcomePair {
	_ = n
	_ = scope
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var leftOutcome *TypeInference_BooleanOutcomePair = this_.priv_TypeInference_traverseWithinShortCircuitingBinOp_5565(left, scope)
	_ = leftOutcome
	return this_.priv_TypeInference_traverseShortCircuitingBinOp_5564(n, left, leftOutcome)
}

// Java method: TypeInference.traverseChildren(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseChildren_5549(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	{
		var el pkg_rhino.Node = n.GetFirstChild()
		_ = el
		for ; el != 0; el = el.GetNext() {
			scope = this_.priv_TypeInference_traverse_5500(el, scope)
		}
	}
	return scope
}

// Java method: TypeInference.traverseGetElem(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope originalScope)
func (this_ *TypeInference) priv_TypeInference_traverseGetElem_5550(n pkg_rhino.Node, originalScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = originalScope
	var executedScope pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverseChildren_5549(n, originalScope)
	_ = executedScope
	return this_.priv_TypeInference_setGetElemNodeTypeAfterChildrenTraversed_5551(n, executedScope)
}

// Java method: TypeInference.setGetElemNodeTypeAfterChildrenTraversed(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scopeAfterChildren)
func (this_ *TypeInference) priv_TypeInference_setGetElemNodeTypeAfterChildrenTraversed_5551(n pkg_rhino.Node, scopeAfterChildren pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scopeAfterChildren
	if !(n.IsGetElem() || n.IsOptChainGetElem()) {
		j2gort.CheckArgument(false)
	}
	this_.priv_TypeInference_inferGetElemType_5552(n)
	scopeAfterChildren = this_.priv_TypeInference_tightenTypeAfterDereference_5560(n.GetFirstChild(), scopeAfterChildren)
	return scopeAfterChildren
}

// Java method: TypeInference.inferGetElemType(com.google.javascript.rhino.Node n)
func (this_ *TypeInference) priv_TypeInference_inferGetElemType_5552(n pkg_rhino.Node) {
	_ = n
	var objType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(n.GetFirstChild()).RestrictByNotNullOrUndefined()
	_ = objType
	var indexKey pkg_rhino.Node = n.GetLastChild()
	_ = indexKey
	var indexType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(indexKey)
	_ = indexType
	var inferredType pkg_rhino.JSType = nil
	_ = inferredType
	base_indexType_ := indexType.GetJSTypeBase()
	if base_indexType_.IsKnownSymbolValueType() {
		inferredType = this_.priv_TypeInference_dereferenceKnownSymbolProp_5553(objType, indexType.ToMaybeKnownSymbolType())
	} else if base_indexType_.IsSymbolValueType() {
		inferredType = pkg_rhino.Upcast_JSType(this_.unknownType)
	} else {
		inferredType = this_.priv_TypeInference_dereferenceIndexSignature_5554(objType)
	}
	n.SetJSType(func() pkg_rhino.JSType {
		if inferredType != nil && !inferredType.IsJavaNil_() {
			return inferredType
		}
		return pkg_rhino.Upcast_JSType(this_.unknownType)
	}())
}

// Java method: TypeInference.dereferenceKnownSymbolProp(com.google.javascript.rhino.jstype.JSType obj, com.google.javascript.rhino.jstype.KnownSymbolType indexSymbol)
func (this_ *TypeInference) priv_TypeInference_dereferenceKnownSymbolProp_5553(obj pkg_rhino.JSType, indexSymbol *pkg_rhino.KnownSymbolType) pkg_rhino.JSType {
	_ = obj
	_ = indexSymbol
	var key pkg_rhino.Property_Key = pkg_rhino.NewProperty_SymbolKey(indexSymbol)
	_ = key
	base_obj_ := obj.GetJSTypeBase()
	if !(base_obj_.IsUnionType()) {
		var objType pkg_rhino.ObjectType = base_obj_.ToMaybeObjectType()
		_ = objType
		if objType != nil && !objType.IsJavaNil_() {
			return objType.GetObjectTypeBase().FindPropertyType_1_Property_Key(key)
		} else {
			return nil
		}
	}
	var union *pkg_rhino.UnionType = obj.ToMaybeUnionType()
	_ = union
	var altCount int32 = union.GetAlternatesCount()
	_ = altCount
	var argumentTypes *pkg_rhino.UnionType_Builder = this_.registry.AcquireUnionTypeBuilder()
	_ = argumentTypes
	defer func() {
		this_.registry.ReleaseUnionTypeBuilder(argumentTypes)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < altCount; i = i + 1 {
				var optionObj pkg_rhino.ObjectType = union.GetAlternateAt(i).GetJSTypeBase().ToMaybeObjectType()
				_ = optionObj
				if optionObj == nil || optionObj.IsJavaNil_() {
					return nil
				}
				var propertyType pkg_rhino.JSType = optionObj.GetObjectTypeBase().FindPropertyType_1_Property_Key(key)
				_ = propertyType
				if propertyType == nil || propertyType.IsJavaNil_() {
					return nil
				}
				argumentTypes.AddAlternate(propertyType)
			}
		}
		return argumentTypes.Build()
	}
	panic("unreachable")
}

// Java method: TypeInference.dereferenceIndexSignature(com.google.javascript.rhino.jstype.JSType obj)
func (this_ *TypeInference) priv_TypeInference_dereferenceIndexSignature_5554(obj pkg_rhino.JSType) pkg_rhino.JSType {
	_ = obj
	if !(obj.GetJSTypeBase().IsUnionType()) {
		var typeMap *pkg_rhino.TemplateTypeMap = obj.GetTemplateTypeMap()
		_ = typeMap
		if !(typeMap.HasTemplateType(this_.registry.GetObjectElementKey())) {
			return nil
		}
		return typeMap.GetResolvedTemplateType(this_.registry.GetObjectElementKey())
	}
	var union *pkg_rhino.UnionType = obj.ToMaybeUnionType()
	_ = union
	var altCount int32 = union.GetAlternatesCount()
	_ = altCount
	var argumentTypes *pkg_rhino.UnionType_Builder = this_.registry.AcquireUnionTypeBuilder()
	_ = argumentTypes
	defer func() {
		this_.registry.ReleaseUnionTypeBuilder(argumentTypes)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < altCount; i = i + 1 {
				var typeMap_2 *pkg_rhino.TemplateTypeMap = union.GetAlternateAt(i).GetTemplateTypeMap()
				_ = typeMap_2
				if !(typeMap_2.HasTemplateType(this_.registry.GetObjectElementKey())) {
					return nil
				}
				argumentTypes.AddAlternate(typeMap_2.GetResolvedTemplateType(this_.registry.GetObjectElementKey()))
			}
		}
		return argumentTypes.Build()
	}
	panic("unreachable")
}

// Java method: TypeInference.traverseGetProp(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseGetProp_5555(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(n, scope)
	return this_.priv_TypeInference_setGetPropNodeTypeAfterChildrenTraversed_5556(n, scope)
}

// Java method: TypeInference.setGetPropNodeTypeAfterChildrenTraversed(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scopeAfterChildren)
func (this_ *TypeInference) priv_TypeInference_setGetPropNodeTypeAfterChildrenTraversed_5556(n pkg_rhino.Node, scopeAfterChildren pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scopeAfterChildren
	if !(n.IsGetProp() || n.IsOptChainGetProp()) {
		j2gort.CheckArgument(false)
	}
	var objNode pkg_rhino.Node = n.GetFirstChild()
	_ = objNode
	n.SetJSType(this_.priv_TypeInference_getPropertyType_5561(objNode.GetJSType(), n.GetString(), n, scopeAfterChildren))
	return this_.priv_TypeInference_tightenTypeAfterDereference_5560(n.GetFirstChild(), scopeAfterChildren)
}

// Java method: TypeInference.setOptChainNodeTypeAfterChildrenTraversed(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scopeAfterChildren)
func (this_ *TypeInference) priv_TypeInference_setOptChainNodeTypeAfterChildrenTraversed_5557(n pkg_rhino.Node, scopeAfterChildren pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scopeAfterChildren
	switch n.TokenSwitchOrd_() {
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return this_.priv_TypeInference_setGetPropNodeTypeAfterChildrenTraversed_5556(n, scopeAfterChildren)
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		return this_.priv_TypeInference_setGetElemNodeTypeAfterChildrenTraversed_5551(n, scopeAfterChildren)
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		return this_.priv_TypeInference_setCallNodeTypeAfterChildrenTraversed_5535(n, scopeAfterChildren)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Illegal token inside finishTraversingOptChain")))
	}
	panic("unreachable switch expression")
}

// Java method: TypeInference.traverseOptChain(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseOptChain_5558(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	if !(NodeUtil_IsOptChainNode(n)) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_IsEndOfOptChainSegment(n) {
		var startOfChain pkg_rhino.Node = NodeUtil_GetStartOfOptChainSegment(n)
		_ = startOfChain
		var val_18 TypeInference_OptChainInfo
		InitTypeInference_OptChainInfo(&val_18, n, startOfChain)
		var optChainInfo *TypeInference_OptChainInfo = &val_18
		_ = optChainInfo
		this_.optChainArrayDeque.AddFirst(optChainInfo)
	}
	var lhsScope pkg_jscomptype.FlowScope = this_.priv_TypeInference_traverse_5500(n.GetFirstChild(), scope)
	_ = lhsScope
	if n.IsOptionalChainStart() {
		{
			r_19 := this_.optChainArrayDeque.PeekFirst()
			r_19.unconditionalScope = lhsScope
		}
	}
	var aboutToExecuteScope pkg_jscomptype.FlowScope = lhsScope
	_ = aboutToExecuteScope
	var nextChild pkg_rhino.Node = n.GetSecondChild()
	_ = nextChild
	for nextChild != 0 {
		aboutToExecuteScope = this_.priv_TypeInference_traverse_5500(nextChild, aboutToExecuteScope)
		nextChild = nextChild.GetNext()
	}
	var executedScope pkg_jscomptype.FlowScope = this_.priv_TypeInference_setOptChainNodeTypeAfterChildrenTraversed_5557(n, aboutToExecuteScope)
	_ = executedScope
	if j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil) {
		n.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
	}
	if NodeUtil_IsEndOfOptChainSegment(n) {
		var startOfChain_2 pkg_rhino.Node = NodeUtil_GetStartOfOptChainSegment(n)
		_ = startOfChain_2
		var currentChain *TypeInference_OptChainInfo = this_.optChainArrayDeque.RemoveFirst()
		_ = currentChain
		if !(currentChain.endOfChain == n) {
			j2gort.CheckState(false)
		}
		if !(currentChain.startOfChain == startOfChain_2) {
			j2gort.CheckState(false)
		}
		var startNode pkg_rhino.Node = j2gort.CheckNotNullNode(startOfChain_2.GetFirstChild())
		_ = startNode
		return this_.priv_TypeInference_updateTypeWhenEndOfOptChain_5559(n, startNode, currentChain, executedScope)
	} else {
		return executedScope
	}
}

// Java method: TypeInference.updateTypeWhenEndOfOptChain(com.google.javascript.rhino.Node optChain, com.google.javascript.rhino.Node startNode, com.google.javascript.jscomp.TypeInference.OptChainInfo currentChain, com.google.javascript.jscomp.type.FlowScope executedScope)
func (this_ *TypeInference) priv_TypeInference_updateTypeWhenEndOfOptChain_5559(optChain pkg_rhino.Node, startNode pkg_rhino.Node, currentChain *TypeInference_OptChainInfo, executedScope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = optChain
	_ = startNode
	_ = currentChain
	_ = executedScope
	var startType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(startNode)
	_ = startType
	base_startType_ := startType.GetJSTypeBase()
	if base_startType_.IsUnknownType() {
		optChain.SetJSType(pkg_rhino.Upcast_JSType(this_.unknownType))
		return this_.Join(executedScope, currentChain.unconditionalScope)
	} else if base_startType_.IsNullType() || base_startType_.IsVoidType() {
		optChain.SetJSType(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()))
		return currentChain.unconditionalScope
	} else if !(startType.IsNullable()) {
		return executedScope
	} else {
		optChain.SetJSType(this_.registry.CreateUnionType_2(this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()), optChain.GetJSType()))
		var unexecutedScope pkg_jscomptype.FlowScope = this_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(startNode, currentChain.unconditionalScope, pkg_rhino.StaticGet_Outcome_NULLISH())
		_ = unexecutedScope
		return this_.Join(executedScope, unexecutedScope)
	}
}

// Java static method: TypeInference.inferPropertyTypesToMatchConstraint(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSType constraint)
func TypeInference_InferPropertyTypesToMatchConstraint(type_ pkg_rhino.JSType, constraint pkg_rhino.JSType) {
	clinit_TypeInference()
	_ = type_
	_ = constraint
	if (type_ == nil || type_.IsJavaNil_()) || (constraint == nil || constraint.IsJavaNil_()) {
		return
	}
	type_.MatchConstraint(constraint)
}

// Java method: TypeInference.tightenTypeAfterDereference(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_tightenTypeAfterDereference_5560(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	if n.IsQualifiedName() {
		var type_ pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(n)
		_ = type_
		var narrowed pkg_rhino.JSType = type_.RestrictByNotNullOrUndefined()
		_ = narrowed
		if !(type_.GetJSTypeBase().EqualsJSType(narrowed)) {
			scope = this_.priv_TypeInference_narrowScope_5538(scope, n, narrowed)
		}
	}
	return scope
}

// Java static method: TypeInference.getGetPropRootName(com.google.javascript.rhino.Node n)
func TypeInference_GetGetPropRootName(n pkg_rhino.Node) string {
	clinit_TypeInference()
	_ = n
	if !(n.IsGetProp()) {
		return j2gort.NullString
	}
	var cur pkg_rhino.Node = n.GetFirstChild()
	_ = cur
	for cur.IsGetProp() {
		cur = cur.GetFirstChild()
	}
	switch cur.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var name string = cur.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "" {
			return j2gort.NullString
		}
		var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
		_ = dot
		{

			if dot < int32(0) {
				return name
			}
			return j2gort.JStrSub(name, int32(0), dot)
		}

	case int32(37) /* Token.THIS */ :
		return "this"
	case int32(110) /* Token.SUPER */ :
		return "super"
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java method: TypeInference.getPropertyType(com.google.javascript.rhino.jstype.JSType objType, java.lang.String propName, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_getPropertyType_5561(objType pkg_rhino.JSType, propName string, n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_rhino.JSType {
	_ = objType
	_ = propName
	_ = n
	_ = scope
	var propertyType pkg_rhino.JSType = nil
	_ = propertyType
	var isLocallyInferred bool = false
	_ = isLocallyInferred
	var rootName string = TypeInference_GetGetPropRootName(n)
	_ = rootName
	var qualifiedName string = j2gort.NullString
	_ = qualifiedName
	var qualifiedNameComputed bool = false
	_ = qualifiedNameComputed
	var var_ pkg_rhino.StaticTypedSlot = nil
	_ = var_
	if rootName != j2gort.NullString {
		var lfs *LinkedFlowScope
		_ = lfs
		var ok_21 bool
		lfs, ok_21 = j2gort.TryCast[*LinkedFlowScope](scope)
		ok_20 := !ok_21
		if ok_20 || (lfs.MayHaveQualifiedSlotForRoot(rootName) && lfs.MayHaveQualifiedSlotForNode(n)) {
			qualifiedName = n.GetQualifiedName()
			qualifiedNameComputed = true
			if qualifiedName != j2gort.NullString {
				var_ = j2gort.AsIface[pkg_rhino.StaticTypedSlot](scope.GetSlot(qualifiedName))
			}
		}
	}
	if var_ != nil && !var_.IsJavaNil_() {
		var varType pkg_rhino.JSType = var_.GetType()
		_ = varType
		if varType != nil && !varType.IsJavaNil_() {
			var isDeclared bool = !(var_.IsTypeInferred())
			_ = isDeclared
			isLocallyInferred = (!(j2gort.IfaceEquals[pkg_rhino.StaticTypedSlot](var_, TypeInference_GetDeclaredVar(scope, qualifiedName))))
			if isDeclared || isLocallyInferred {
				propertyType = varType
			}
		}
	}
	if (propertyType == nil || propertyType.IsJavaNil_()) && (objType != nil && !objType.IsJavaNil_()) {
		var foundType pkg_rhino.JSType = objType.GetJSTypeBase().FindPropertyType_1_string(propName)
		_ = foundType
		if foundType != nil && !foundType.IsJavaNil_() {
			propertyType = foundType
		}
	}
	if ((propertyType == nil || propertyType.IsJavaNil_()) || propertyType.GetJSTypeBase().IsUnknownType()) && (rootName != j2gort.NullString) {
		var declScope pkg_rhino.StaticTypedScope = j2gort.AsIface[pkg_rhino.StaticTypedScope](scope.GetDeclarationScope())
		_ = declScope
		if qualifiedNameComputed || this_.registry.MayHaveQualifiedTypeForGetProp(pkg_rhino.Upcast_StaticScope(declScope), rootName, n) {
			if !(qualifiedNameComputed) {
				qualifiedName = n.GetQualifiedName()
			}
			if qualifiedName != j2gort.NullString {
				var regType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(this_.registry.GetType_2(pkg_rhino.Upcast_StaticScope(declScope), qualifiedName))
				_ = regType
				if regType != nil && !regType.IsJavaNil_() {
					propertyType = pkg_rhino.ToIface_JSType(regType.GetConstructor())
				}
			}
		}
	}
	if propertyType == nil || propertyType.IsJavaNil_() {
		return pkg_rhino.Upcast_JSType(this_.unknownType)
	} else if propertyType.GetJSTypeBase().EqualsJSType(pkg_rhino.Upcast_JSType(this_.unknownType)) && isLocallyInferred {
		return this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())
	} else {
		return propertyType
	}
}

// Java method: TypeInference.traverseOr(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseOr_5562(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) *TypeInference_BooleanOutcomePair {
	_ = n
	_ = scope
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var leftOutcome *TypeInference_BooleanOutcomePair = this_.priv_TypeInference_traverseWithinShortCircuitingBinOp_5565(left, scope)
	_ = leftOutcome
	return this_.priv_TypeInference_traverseShortCircuitingBinOp_5564(n, left, leftOutcome)
}

// Java method: TypeInference.traverseShortCircuitingBinOpAssignment(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseShortCircuitingBinOpAssignment_5563(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = n
	_ = scope
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var nIsAnd bool = n.IsAssignAnd()
	_ = nIsAnd
	var leftOutcome *TypeInference_BooleanOutcomePair = this_.priv_TypeInference_traverseWithinShortCircuitingBinOp_5565(left, scope)
	_ = leftOutcome
	var outcome *TypeInference_BooleanOutcomePair = this_.priv_TypeInference_traverseShortCircuitingBinOp_5564(n, left, leftOutcome)
	_ = outcome
	var outcomeJoinedFlowScope pkg_jscomptype.FlowScope = outcome.GetJoinedFlowScope()
	_ = outcomeJoinedFlowScope
	if leftOutcome.toBooleanOutcomes == pkg_rhino.BooleanLiteralSet_Get(!(nIsAnd)) {
		return outcomeJoinedFlowScope
	}
	return this_.priv_TypeInference_updateScopeForAssignment_5513(outcomeJoinedFlowScope, n.GetFirstChild(), n.GetJSType(), StaticGet_TypeInference_AssignmentType_ASSIGN())
}

// Java method: TypeInference.traverseShortCircuitingBinOp(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node left, com.google.javascript.jscomp.TypeInference.BooleanOutcomePair leftOutcome)
func (this_ *TypeInference) priv_TypeInference_traverseShortCircuitingBinOp_5564(n pkg_rhino.Node, left pkg_rhino.Node, leftOutcome *TypeInference_BooleanOutcomePair) *TypeInference_BooleanOutcomePair {
	_ = n
	_ = left
	_ = leftOutcome
	if !(((n.IsAnd() || n.IsOr()) || n.IsAssignAnd()) || n.IsAssignOr()) {
		j2gort.CheckArgument(false)
	}
	var nIsAnd bool = (n.IsAnd() || n.IsAssignAnd())
	_ = nIsAnd
	var right pkg_rhino.Node = n.GetLastChild()
	_ = right
	var leftType pkg_rhino.JSType = left.GetJSType()
	_ = leftType
	var rightScope pkg_jscomptype.FlowScope = this_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(left, leftOutcome.GetOutcomeFlowScope(left.GetToken(), nIsAnd), pkg_rhino.Outcome_ForBoolean(nIsAnd))
	_ = rightScope
	var rightOutcome *TypeInference_BooleanOutcomePair = this_.priv_TypeInference_traverseWithinShortCircuitingBinOp_5565(right, rightScope)
	_ = rightOutcome
	var rightType pkg_rhino.JSType = right.GetJSType()
	_ = rightType
	var type_ pkg_rhino.JSType = nil
	_ = type_
	var outcome *TypeInference_BooleanOutcomePair = nil
	_ = outcome
	if (leftType != nil && !leftType.IsJavaNil_()) && (rightType != nil && !rightType.IsJavaNil_()) {
		leftType = leftType.GetRestrictedTypeGivenOutcome(pkg_rhino.Outcome_ForBoolean(!(nIsAnd)))
		if leftOutcome.toBooleanOutcomes == pkg_rhino.BooleanLiteralSet_Get(!(nIsAnd)) {
			type_ = leftType
			outcome = leftOutcome
		} else {
			type_ = leftType.GetLeastSupertype(rightType)
			outcome = NewTypeInference_BooleanOutcomePair(this_, TypeInference_JoinBooleanOutcomes(nIsAnd, leftOutcome.toBooleanOutcomes, rightOutcome.toBooleanOutcomes), TypeInference_JoinBooleanOutcomes(nIsAnd, leftOutcome.booleanValues, rightOutcome.booleanValues), leftOutcome.GetJoinedFlowScope(), rightOutcome.GetJoinedFlowScope())
		}
		if (outcome.booleanValues == pkg_rhino.StaticGet_BooleanLiteralSet_EMPTY()) && this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE()).GetJSTypeBase().IsSubtypeOf_1(type_) {
			if type_.GetJSTypeBase().IsUnionType() {
				type_ = type_.ToMaybeUnionType().GetRestrictedUnion(this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE()))
			}
		}
	} else {
		type_ = pkg_rhino.Upcast_JSType(this_.unknownType)
		outcome = NewTypeInference_BooleanOutcomePair(this_, pkg_rhino.StaticGet_BooleanLiteralSet_BOTH(), pkg_rhino.StaticGet_BooleanLiteralSet_BOTH(), leftOutcome.GetJoinedFlowScope(), rightOutcome.GetJoinedFlowScope())
	}
	n.SetJSType(type_)
	return outcome
}

// Java method: TypeInference.traverseWithinShortCircuitingBinOp(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseWithinShortCircuitingBinOp_5565(n pkg_rhino.Node, scope pkg_jscomptype.FlowScope) *TypeInference_BooleanOutcomePair {
	_ = n
	_ = scope
	switch n.TokenSwitchOrd_() {
	case int32(69) /* Token.AND */ :
		return this_.priv_TypeInference_traverseAnd_5548(n, scope)
	case int32(68) /* Token.OR */ :
		return this_.priv_TypeInference_traverseOr_5562(n, scope)
	default:
		scope = this_.priv_TypeInference_traverse_5500(n, scope)
		return this_.priv_TypeInference_newBooleanOutcomePair_5568(n.GetJSType(), scope)
	}
	panic("unreachable switch expression")
}

// Java method: TypeInference.traverseAwait(com.google.javascript.rhino.Node await, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseAwait_5566(await pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = await
	_ = scope
	scope = this_.priv_TypeInference_traverseChildren_5549(await, scope)
	var expr pkg_rhino.Node = await.GetFirstChild()
	_ = expr
	var exprType pkg_rhino.JSType = this_.priv_TypeInference_getJSType_5571(expr)
	_ = exprType
	await.SetJSType(Promises_GetResolvedType(this_.registry, exprType))
	return scope
}

// Java method: TypeInference.traverseDynamicImport(com.google.javascript.rhino.Node dynamicImport, com.google.javascript.jscomp.type.FlowScope scope)
func (this_ *TypeInference) priv_TypeInference_traverseDynamicImport_5567(dynamicImport pkg_rhino.Node, scope pkg_jscomptype.FlowScope) pkg_jscomptype.FlowScope {
	_ = dynamicImport
	_ = scope
	var templateType pkg_rhino.JSType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = templateType
	var moduleMap *ModuleMap = this_.compiler.GetModuleMap()
	_ = moduleMap
	var importSpecifier pkg_rhino.Node = dynamicImport.GetFirstChild()
	_ = importSpecifier
	var input *CompilerInput = this_.compiler.GetInput(NodeUtil_GetInputId(dynamicImport))
	_ = input
	if (importSpecifier.IsStringLit() && (moduleMap != nil)) && (input != nil) {
		var targetPath *ModuleLoader_ModulePath = input.GetPath().ResolveJsModule(importSpecifier.GetString(), importSpecifier.GetSourceFileName(), importSpecifier.GetLineno(), importSpecifier.GetCharno())
		_ = targetPath
		var targetModule *Module
		_ = targetModule
		if targetPath != nil {
			targetModule = moduleMap.GetModule_1_ModuleLoader_ModulePath(targetPath)
		} else {
			targetModule = nil
		}
		if targetModule != nil {
			var scriptNode pkg_rhino.Node = targetModule.Metadata().RootNode()
			_ = scriptNode
			if scriptNode.HasOneChild() && scriptNode.GetFirstChild().IsModuleBody() {
				var exportNamespaceType pkg_rhino.JSType
				_ = exportNamespaceType
				if scriptNode != 0 {
					exportNamespaceType = scriptNode.GetOnlyChild().GetJSType()
				} else {
					exportNamespaceType = nil
				}
				if exportNamespaceType != nil && !exportNamespaceType.IsJavaNil_() {
					templateType = exportNamespaceType
				}
			} else {
				var moduleName pkg_rhino.Node = NodeUtil_FindPreorder(scriptNode, j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_22 uint32) bool {
					node := pkg_rhino.Node(rawNode_22)
					_ = node
					return node.MatchesQualifiedName_1_string(targetPath.ToModuleName())
				}), j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
				_ = moduleName
				if (moduleName != 0) && !(j2gort.IfaceEquals[pkg_rhino.JSType](moduleName.GetJSType(), nil)) {
					templateType = moduleName.GetJSType()
				}
			}
		}
	}
	dynamicImport.SetJSType(pkg_rhino.ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), []pkg_rhino.JSType{templateType})))
	return this_.priv_TypeInference_traverseChildren_5549(dynamicImport, scope)
}

// Java static method: TypeInference.joinBooleanOutcomes(boolean isAnd, com.google.javascript.rhino.jstype.BooleanLiteralSet left, com.google.javascript.rhino.jstype.BooleanLiteralSet right)
func TypeInference_JoinBooleanOutcomes(isAnd bool, left pkg_rhino.BooleanLiteralSet, right pkg_rhino.BooleanLiteralSet) pkg_rhino.BooleanLiteralSet {
	clinit_TypeInference()
	_ = isAnd
	_ = left
	_ = right
	return right.Union(left.Intersection(pkg_rhino.BooleanLiteralSet_Get(!(isAnd))))
}

// Java method: TypeInference.newBooleanOutcomePair(com.google.javascript.rhino.jstype.JSType jsType, com.google.javascript.jscomp.type.FlowScope flowScope)
func (this_ *TypeInference) priv_TypeInference_newBooleanOutcomePair_5568(jsType pkg_rhino.JSType, flowScope pkg_jscomptype.FlowScope) *TypeInference_BooleanOutcomePair {
	_ = jsType
	_ = flowScope
	if jsType == nil || jsType.IsJavaNil_() {
		return NewTypeInference_BooleanOutcomePair(this_, pkg_rhino.StaticGet_BooleanLiteralSet_BOTH(), pkg_rhino.StaticGet_BooleanLiteralSet_BOTH(), flowScope, flowScope)
	}
	return NewTypeInference_BooleanOutcomePair(this_, jsType.GetPossibleToBooleanOutcomes(), func() pkg_rhino.BooleanLiteralSet {
		if this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE()).GetJSTypeBase().IsSubtypeOf_1(jsType) {
			return pkg_rhino.StaticGet_BooleanLiteralSet_BOTH()
		}
		return pkg_rhino.StaticGet_BooleanLiteralSet_EMPTY()
	}(), flowScope, flowScope)
}

// Java method: TypeInference.redeclareSimpleVar(com.google.javascript.jscomp.type.FlowScope scope, com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.jstype.JSType varType)
func (this_ *TypeInference) priv_TypeInference_redeclareSimpleVar_5569(scope pkg_jscomptype.FlowScope, nameNode pkg_rhino.Node, varType pkg_rhino.JSType) pkg_jscomptype.FlowScope {
	_ = scope
	_ = nameNode
	_ = varType
	if !(nameNode.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(nameNode))
	}
	var varName string = nameNode.GetString()
	_ = varName
	if varType == nil || varType.IsJavaNil_() {
		varType = this_.priv_TypeInference_getNativeType_5572(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	if this_.priv_TypeInference_isUnflowable_5570(TypeInference_GetDeclaredVar(scope, varName)) {
		return scope
	}
	return j2gort.AsIface[pkg_jscomptype.FlowScope](scope.InferSlotType(varName, varType))
}

// Java method: TypeInference.isUnflowable(com.google.javascript.jscomp.TypedVar v)
func (this_ *TypeInference) priv_TypeInference_isUnflowable_5570(v *TypedVar) bool {
	_ = v
	return ((((v != nil) && v.IsLocal()) && v.IsMarkedEscaped()) && (j2gort.CastPtr[TypedScope](j2gort.CastPtr[TypedScope](v.GetScope()).GetClosestContainerScope()) == this_.containerScope))
}

// Java method: TypeInference.getJSType(com.google.javascript.rhino.Node n)
func (this_ *TypeInference) priv_TypeInference_getJSType_5571(n pkg_rhino.Node) pkg_rhino.JSType {
	_ = n
	var jsType pkg_rhino.JSType = n.GetJSType()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return pkg_rhino.Upcast_JSType(this_.unknownType)
	} else {
		return jsType
	}
}

// Java method: TypeInference.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *TypeInference) priv_TypeInference_getNativeType_5572(typeId pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = typeId
	return this_.registry.GetNativeType(typeId)
}

// Java static method: TypeInference.getDeclaredVar(com.google.javascript.jscomp.type.FlowScope scope, java.lang.String name)
func TypeInference_GetDeclaredVar(scope pkg_jscomptype.FlowScope, name string) *TypedVar {
	clinit_TypeInference()
	_ = scope
	_ = name
	return (j2gort.CastPtr[TypedScope](j2gort.AsIface[pkg_rhino.StaticTypedScope](scope.GetDeclarationScope()))).LookupTypedVar(name)
}

// Java interface: com.google.javascript.jscomp.TypeInference.TypeDeclaringCallback
type TypeInference_TypeDeclaringCallback interface {
	Implements_jscomp_TypeInference_TypeDeclaringCallback()
	IsJavaNil_() bool
	DeclareTypeInScope(pkg_jscomptype.FlowScope, pkg_rhino.Node, pkg_rhino.JSType) pkg_jscomptype.FlowScope
}

func ToIface_TypeInference_TypeDeclaringCallback[T interface {
	TypeInference_TypeDeclaringCallback
	comparable
}](v T) TypeInference_TypeDeclaringCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_TypeInference_TypeDeclaringCallback[T TypeInference_TypeDeclaringCallback](v T) TypeInference_TypeDeclaringCallback {
	var vi TypeInference_TypeDeclaringCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference.TypeDeclaringCallback", "com.google.javascript.jscomp.TypeInference$TypeDeclaringCallback", func(v_ any) bool {
		return j2gort.InstanceOf[TypeInference_TypeDeclaringCallback](v_)
	})
}

type lambda_TypeInference_TypeDeclaringCallback struct {
	fn func(pkg_jscomptype.FlowScope, pkg_rhino.Node, pkg_rhino.JSType) pkg_jscomptype.FlowScope
}

func (l *lambda_TypeInference_TypeDeclaringCallback) IsJavaNil_() bool { return l == nil }
func NewLambda_TypeInference_TypeDeclaringCallback(fn func(pkg_jscomptype.FlowScope, pkg_rhino.Node, pkg_rhino.JSType) pkg_jscomptype.FlowScope) TypeInference_TypeDeclaringCallback {
	if fn == nil {
		return nil
	}
	return &lambda_TypeInference_TypeDeclaringCallback{fn: fn}
}

func (*lambda_TypeInference_TypeDeclaringCallback) Implements_jscomp_TypeInference_TypeDeclaringCallback() {
}
func (l *lambda_TypeInference_TypeDeclaringCallback) DeclareTypeInScope(p0_ pkg_jscomptype.FlowScope, p1_ pkg_rhino.Node, p2_ pkg_rhino.JSType) pkg_jscomptype.FlowScope {
	return l.fn(p0_, p1_, p2_)
}

// Java class: com.google.javascript.jscomp.TypeInference.ReusableFlowBrancher
type TypeInference_ReusableFlowBrancher struct {
	outer_             *TypeInference
	output             pkg_jscomptype.FlowScope          `gson:"output"`
	conditionFlowScope pkg_jscomptype.FlowScope          `gson:"conditionFlowScope"`
	conditionOutcomes  *TypeInference_BooleanOutcomePair `gson:"conditionOutcomes"`
	source             pkg_rhino.Node                    `gson:"source"`
	condition          pkg_rhino.Node                    `gson:"condition"`
	idHash_            atomic.Int32
}

func (this_ *TypeInference_ReusableFlowBrancher) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeInference_ReusableFlowBrancher) IsJavaNil_() bool { return this_ == nil }
func (*TypeInference_ReusableFlowBrancher) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInference$ReusableFlowBrancher", "ReusableFlowBrancher"
}
func (this_ *TypeInference_ReusableFlowBrancher) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeInference_ReusableFlowBrancher) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference.ReusableFlowBrancher", "com.google.javascript.jscomp.TypeInference$ReusableFlowBrancher", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInference_ReusableFlowBrancher](v_)
	})
}
func (this_ *TypeInference_ReusableFlowBrancher) SetOuter_(o *TypeInference)                 { this_.outer_ = o }
func (this_ *TypeInference_ReusableFlowBrancher) GetOuter_() *TypeInference                  { return this_.outer_ }
func (*TypeInference_ReusableFlowBrancher) Implements_jscomp_DataFlowAnalysis_FlowBrancher() {}

func (this_ *TypeInference_ReusableFlowBrancher) initFields_TypeInference_ReusableFlowBrancher() {
	_ = this_

}

// Java constructor: ReusableFlowBrancher()
func NewTypeInference_ReusableFlowBrancher(outer_ *TypeInference) *TypeInference_ReusableFlowBrancher {
	this_ := &TypeInference_ReusableFlowBrancher{}
	this_.outer_ = outer_
	InitTypeInference_ReusableFlowBrancher(this_)
	return this_
}

func InitTypeInference_ReusableFlowBrancher(this_ *TypeInference_ReusableFlowBrancher) {
	this_.initFields_TypeInference_ReusableFlowBrancher()
}

// Java method: ReusableFlowBrancher.reset(com.google.javascript.rhino.Node source, com.google.javascript.jscomp.type.FlowScope output)
func (this_ *TypeInference_ReusableFlowBrancher) Reset(source pkg_rhino.Node, output pkg_jscomptype.FlowScope) DataFlowAnalysis_FlowBrancher[pkg_jscomptype.FlowScope] {
	_ = source
	_ = output
	this_.source = source
	this_.output = output
	this_.condition = pkg_rhino.Node(0)
	this_.conditionFlowScope = nil
	this_.conditionOutcomes = nil
	return ToIface_DataFlowAnalysis_FlowBrancher[pkg_jscomptype.FlowScope](this_)
}

// Java method: ReusableFlowBrancher.clear()
func (this_ *TypeInference_ReusableFlowBrancher) Clear() {
	this_.source = pkg_rhino.Node(0)
	this_.output = nil
	this_.condition = pkg_rhino.Node(0)
	this_.conditionFlowScope = nil
	this_.conditionOutcomes = nil
}

// Java method: ReusableFlowBrancher.branchFlow(com.google.javascript.jscomp.ControlFlowGraph.Branch branch)
func (this_ *TypeInference_ReusableFlowBrancher) BranchFlow(branch ControlFlowGraph_Branch) pkg_jscomptype.FlowScope {
	_ = branch
	switch (branch).Ordinal() {
	case int32(0) /* Branch.ON_TRUE */ :
		if NodeUtil_IsEnhancedFor(this_.source) {
			return this_.outer_.priv_TypeInference_initializeEnhancedForScope_5501(this_.source, this_.output)
		}
		fallthrough
	case int32(1) /* Branch.ON_FALSE */ :
		if this_.condition == 0 {
			if this_.source.IsCase() {
				this_.condition = this_.source
				this_.conditionFlowScope = this_.outer_.priv_TypeInference_traverse_5500(this_.condition.GetFirstChild(), this_.output)
			} else {
				this_.condition = NodeUtil_GetConditionExpression(this_.source)
				if this_.condition == 0 {
					return this_.output
				}
			}
		}
		if this_.condition.IsAnd() || this_.condition.IsOr() {
			if this_.conditionOutcomes == nil {
				if this_.condition.IsAnd() {
					this_.conditionOutcomes = this_.outer_.priv_TypeInference_traverseAnd_5548(this_.condition, this_.output)
				} else {
					this_.conditionOutcomes = this_.outer_.priv_TypeInference_traverseOr_5562(this_.condition, this_.output)
				}
			}
			return this_.outer_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(this_.condition, this_.conditionOutcomes.GetOutcomeFlowScope(this_.condition.GetToken(), (branch == StaticGet_ControlFlowGraph_Branch_ON_TRUE())), pkg_rhino.Outcome_ForBoolean(func() bool {
				re_23 := branch
				ae_24 := StaticGet_ControlFlowGraph_Branch_ON_TRUE()
				if re_23 == 0 {
					j2gort.ThrowNPE()
				}
				return re_23 == ae_24
			}()))
		}
		if this_.conditionFlowScope == nil || this_.conditionFlowScope.IsJavaNil_() {
			this_.conditionFlowScope = this_.outer_.priv_TypeInference_traverse_5500(this_.condition, this_.output)
		}
		return this_.outer_.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(this_.condition, this_.conditionFlowScope, pkg_rhino.Outcome_ForBoolean(func() bool {
			re_25 := branch
			ae_26 := StaticGet_ControlFlowGraph_Branch_ON_TRUE()
			if re_25 == 0 {
				j2gort.ThrowNPE()
			}
			return re_25 == ae_26
		}()))
	default:
		return this_.output
	}
	panic("unreachable")
}

// Java enum: com.google.javascript.jscomp.TypeInference.AssignmentType
type TypeInference_AssignmentType uint8

const (
	enum_TypeInference_AssignmentType_DECLARATION TypeInference_AssignmentType = 1
	enum_TypeInference_AssignmentType_ASSIGN      TypeInference_AssignmentType = 2
)

var enumBaseTable_TypeInference_AssignmentType = [...]j2gort.EnumBase{
	{Name_: "DECLARATION", Ordinal_: 0},
	{Name_: "ASSIGN", Ordinal_: 1},
}
var enumValues_TypeInference_AssignmentType = []TypeInference_AssignmentType{1, 2}

func StaticGet_TypeInference_AssignmentType_DECLARATION() TypeInference_AssignmentType {
	return enum_TypeInference_AssignmentType_DECLARATION
}
func StaticGet_TypeInference_AssignmentType_ASSIGN() TypeInference_AssignmentType {
	return enum_TypeInference_AssignmentType_ASSIGN
}

func TypeInference_AssignmentType_Values() []TypeInference_AssignmentType {
	return j2gort.CloneSlice(enumValues_TypeInference_AssignmentType)
}

func TypeInference_AssignmentType_ValueOf(name string) TypeInference_AssignmentType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypeInference_AssignmentType {
		if enumBaseTable_TypeInference_AssignmentType[i_].Name_ == name {
			return TypeInference_AssignmentType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.TypeInference.AssignmentType." + name))
}

func (this_ TypeInference_AssignmentType) IsJavaNil_() bool                              { return this_ == 0 }
func (TypeInference_AssignmentType) HasIdentityEquals_()                                 {}
func (TypeInference_AssignmentType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (TypeInference_AssignmentType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypeInference_AssignmentType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (TypeInference_AssignmentType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInference$AssignmentType", "AssignmentType"
}
func (this_ TypeInference_AssignmentType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypeInference_AssignmentType[this_-1]
}
func (this_ TypeInference_AssignmentType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypeInference_AssignmentType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypeInference_AssignmentType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeInference_AssignmentType[this_-1].Name_
}
func (this_ TypeInference_AssignmentType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypeInference_AssignmentType[this_-1].Name_
}
func (this_ TypeInference_AssignmentType) Equals(other_ any) bool {
	o_, ok_ := other_.(TypeInference_AssignmentType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypeInference_AssignmentType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypeInference_AssignmentType[this_-1].HashCode()
}
func (this_ TypeInference_AssignmentType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeInference_AssignmentType[this_-1].CompareTo(other_)
}
func (this_ TypeInference_AssignmentType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypeInference_AssignmentType[this_-1].Name_)
}
func init() {
	vals_ := []any{TypeInference_AssignmentType(1), TypeInference_AssignmentType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeInference.AssignmentType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeInference$AssignmentType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference.AssignmentType", "com.google.javascript.jscomp.TypeInference$AssignmentType", func(v_ any) bool {
		return j2gort.InstanceOf[TypeInference_AssignmentType](v_)
	})
}

// Java enum: com.google.javascript.jscomp.TypeInference.BigIntPresence
type TypeInference_BigIntPresence uint8

const (
	enum_TypeInference_BigIntPresence_NO_BIGINT        TypeInference_BigIntPresence = 1
	enum_TypeInference_BigIntPresence_ALL_BIGINT       TypeInference_BigIntPresence = 2
	enum_TypeInference_BigIntPresence_BIGINT_OR_NUMBER TypeInference_BigIntPresence = 3
	enum_TypeInference_BigIntPresence_BIGINT_OR_OTHER  TypeInference_BigIntPresence = 4
)

var enumBaseTable_TypeInference_BigIntPresence = [...]j2gort.EnumBase{
	{Name_: "NO_BIGINT", Ordinal_: 0},
	{Name_: "ALL_BIGINT", Ordinal_: 1},
	{Name_: "BIGINT_OR_NUMBER", Ordinal_: 2},
	{Name_: "BIGINT_OR_OTHER", Ordinal_: 3},
}
var enumValues_TypeInference_BigIntPresence = []TypeInference_BigIntPresence{1, 2, 3, 4}

func StaticGet_TypeInference_BigIntPresence_NO_BIGINT() TypeInference_BigIntPresence {
	return enum_TypeInference_BigIntPresence_NO_BIGINT
}
func StaticGet_TypeInference_BigIntPresence_ALL_BIGINT() TypeInference_BigIntPresence {
	return enum_TypeInference_BigIntPresence_ALL_BIGINT
}
func StaticGet_TypeInference_BigIntPresence_BIGINT_OR_NUMBER() TypeInference_BigIntPresence {
	return enum_TypeInference_BigIntPresence_BIGINT_OR_NUMBER
}
func StaticGet_TypeInference_BigIntPresence_BIGINT_OR_OTHER() TypeInference_BigIntPresence {
	return enum_TypeInference_BigIntPresence_BIGINT_OR_OTHER
}

func TypeInference_BigIntPresence_Values() []TypeInference_BigIntPresence {
	return j2gort.CloneSlice(enumValues_TypeInference_BigIntPresence)
}

func TypeInference_BigIntPresence_ValueOf(name string) TypeInference_BigIntPresence {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TypeInference_BigIntPresence {
		if enumBaseTable_TypeInference_BigIntPresence[i_].Name_ == name {
			return TypeInference_BigIntPresence(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.TypeInference.BigIntPresence." + name))
}

func (this_ TypeInference_BigIntPresence) IsJavaNil_() bool                              { return this_ == 0 }
func (TypeInference_BigIntPresence) HasIdentityEquals_()                                 {}
func (TypeInference_BigIntPresence) Implements_j2gort_Ext_java_lang_Enum()               {}
func (TypeInference_BigIntPresence) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TypeInference_BigIntPresence) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (TypeInference_BigIntPresence) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInference$BigIntPresence", "BigIntPresence"
}
func (this_ TypeInference_BigIntPresence) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TypeInference_BigIntPresence[this_-1]
}
func (this_ TypeInference_BigIntPresence) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TypeInference_BigIntPresence) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TypeInference_BigIntPresence) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeInference_BigIntPresence[this_-1].Name_
}
func (this_ TypeInference_BigIntPresence) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TypeInference_BigIntPresence[this_-1].Name_
}
func (this_ TypeInference_BigIntPresence) Equals(other_ any) bool {
	o_, ok_ := other_.(TypeInference_BigIntPresence)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TypeInference_BigIntPresence) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TypeInference_BigIntPresence[this_-1].HashCode()
}
func (this_ TypeInference_BigIntPresence) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TypeInference_BigIntPresence[this_-1].CompareTo(other_)
}
func (this_ TypeInference_BigIntPresence) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TypeInference_BigIntPresence[this_-1].Name_)
}
func init() {
	vals_ := []any{TypeInference_BigIntPresence(1), TypeInference_BigIntPresence(2), TypeInference_BigIntPresence(3), TypeInference_BigIntPresence(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeInference.BigIntPresence", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.TypeInference$BigIntPresence", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference.BigIntPresence", "com.google.javascript.jscomp.TypeInference$BigIntPresence", func(v_ any) bool {
		return j2gort.InstanceOf[TypeInference_BigIntPresence](v_)
	})
}

// Java class: com.google.javascript.jscomp.TypeInference.OptChainInfo
type TypeInference_OptChainInfo struct {
	unconditionalScope pkg_jscomptype.FlowScope `gson:"unconditionalScope"`
	endOfChain         pkg_rhino.Node           `gson:"endOfChain"`
	startOfChain       pkg_rhino.Node           `gson:"startOfChain"`
	idHash_            atomic.Int32
}

func (this_ *TypeInference_OptChainInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeInference_OptChainInfo) IsJavaNil_() bool { return this_ == nil }
func (*TypeInference_OptChainInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInference$OptChainInfo", "OptChainInfo"
}
func (this_ *TypeInference_OptChainInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeInference_OptChainInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference.OptChainInfo", "com.google.javascript.jscomp.TypeInference$OptChainInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInference_OptChainInfo](v_)
	})
}

func (this_ *TypeInference_OptChainInfo) initFields_TypeInference_OptChainInfo() {
	_ = this_
}

// Java constructor: OptChainInfo(com.google.javascript.rhino.Node endOfChain, com.google.javascript.rhino.Node startOfChain)
func NewTypeInference_OptChainInfo(endOfChain pkg_rhino.Node, startOfChain pkg_rhino.Node) *TypeInference_OptChainInfo {
	this_ := &TypeInference_OptChainInfo{}
	InitTypeInference_OptChainInfo(this_, endOfChain, startOfChain)
	return this_
}

func InitTypeInference_OptChainInfo(this_ *TypeInference_OptChainInfo, endOfChain pkg_rhino.Node, startOfChain pkg_rhino.Node) {
	_ = endOfChain
	_ = startOfChain
	this_.initFields_TypeInference_OptChainInfo()
	this_.endOfChain = endOfChain
	this_.startOfChain = startOfChain
	this_.unconditionalScope = nil
}

// Java class: com.google.javascript.jscomp.TypeInference.BooleanOutcomePair
type TypeInference_BooleanOutcomePair struct {
	outer_            *TypeInference
	leftScope         pkg_jscomptype.FlowScope `gson:"leftScope"`
	rightScope        pkg_jscomptype.FlowScope `gson:"rightScope"`
	joinedScope       pkg_jscomptype.FlowScope `gson:"joinedScope"`
	idHash_           atomic.Int32
	toBooleanOutcomes pkg_rhino.BooleanLiteralSet `gson:"toBooleanOutcomes"`
	booleanValues     pkg_rhino.BooleanLiteralSet `gson:"booleanValues"`
}

func (this_ *TypeInference_BooleanOutcomePair) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeInference_BooleanOutcomePair) IsJavaNil_() bool { return this_ == nil }
func (*TypeInference_BooleanOutcomePair) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInference$BooleanOutcomePair", "BooleanOutcomePair"
}
func (this_ *TypeInference_BooleanOutcomePair) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeInference_BooleanOutcomePair) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInference.BooleanOutcomePair", "com.google.javascript.jscomp.TypeInference$BooleanOutcomePair", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInference_BooleanOutcomePair](v_)
	})
}
func (this_ *TypeInference_BooleanOutcomePair) SetOuter_(o *TypeInference) { this_.outer_ = o }
func (this_ *TypeInference_BooleanOutcomePair) GetOuter_() *TypeInference  { return this_.outer_ }

func (this_ *TypeInference_BooleanOutcomePair) initFields_TypeInference_BooleanOutcomePair() {
	_ = this_

}

// Java constructor: BooleanOutcomePair(com.google.javascript.rhino.jstype.BooleanLiteralSet toBooleanOutcomes, com.google.javascript.rhino.jstype.BooleanLiteralSet booleanValues, com.google.javascript.jscomp.type.FlowScope leftScope, com.google.javascript.jscomp.type.FlowScope rightScope)
func NewTypeInference_BooleanOutcomePair(outer_ *TypeInference, toBooleanOutcomes pkg_rhino.BooleanLiteralSet, booleanValues pkg_rhino.BooleanLiteralSet, leftScope pkg_jscomptype.FlowScope, rightScope pkg_jscomptype.FlowScope) *TypeInference_BooleanOutcomePair {
	this_ := &TypeInference_BooleanOutcomePair{}
	this_.outer_ = outer_
	InitTypeInference_BooleanOutcomePair(this_, toBooleanOutcomes, booleanValues, leftScope, rightScope)
	return this_
}

func InitTypeInference_BooleanOutcomePair(this_ *TypeInference_BooleanOutcomePair, toBooleanOutcomes pkg_rhino.BooleanLiteralSet, booleanValues pkg_rhino.BooleanLiteralSet, leftScope pkg_jscomptype.FlowScope, rightScope pkg_jscomptype.FlowScope) {
	_ = toBooleanOutcomes
	_ = booleanValues
	_ = leftScope
	_ = rightScope
	this_.initFields_TypeInference_BooleanOutcomePair()
	this_.toBooleanOutcomes = toBooleanOutcomes
	this_.booleanValues = booleanValues
	this_.leftScope = leftScope
	this_.rightScope = rightScope
}

// Java method: BooleanOutcomePair.getJoinedFlowScope()
func (this_ *TypeInference_BooleanOutcomePair) GetJoinedFlowScope() pkg_jscomptype.FlowScope {
	if this_.joinedScope == nil || this_.joinedScope.IsJavaNil_() {
		if this_.leftScope == this_.rightScope || ((this_.leftScope == nil || this_.leftScope.IsJavaNil_()) && (this_.rightScope == nil || this_.rightScope.IsJavaNil_())) {
			this_.joinedScope = this_.rightScope
		} else {
			this_.joinedScope = this_.outer_.Join(this_.leftScope, this_.rightScope)
		}
	}
	return this_.joinedScope
}

// Java method: BooleanOutcomePair.getOutcomeFlowScope(com.google.javascript.rhino.Token nodeType, boolean outcome)
func (this_ *TypeInference_BooleanOutcomePair) GetOutcomeFlowScope(nodeType pkg_rhino.Token, outcome bool) pkg_jscomptype.FlowScope {
	_ = nodeType
	_ = outcome
	if ((nodeType == pkg_rhino.StaticGet_Token_AND()) && outcome) || ((nodeType == pkg_rhino.StaticGet_Token_OR()) && !(outcome)) {
		return this_.rightScope
	} else {
		return this_.GetJoinedFlowScope()
	}
}
