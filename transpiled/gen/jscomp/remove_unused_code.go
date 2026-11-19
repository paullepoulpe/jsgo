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
// Port of com/google/javascript/jscomp/RemoveUnusedCode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_resources "github.com/paullepoulpe/jsgo/transpiled/gen/resources"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_remove_unused_code_1 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_2 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_3 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_4 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_5 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_6 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_7 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_8 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_9 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_10 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_11 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_12 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_13 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_14 j2gort.Ext_java_util_function_Supplier
var lambdaSingleton_remove_unused_code_15 j2gort.Ext_java_util_function_Supplier

func init() {
	lambdaSingleton_remove_unused_code_1 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return ""
	})
	lambdaSingleton_remove_unused_code_2 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "exported class"
	})
	lambdaSingleton_remove_unused_code_3 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "catch variable"
	})
	lambdaSingleton_remove_unused_code_4 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "for-of or for-in loop variable"
	})
	lambdaSingleton_remove_unused_code_5 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "for-of or for-in loop variable"
	})
	lambdaSingleton_remove_unused_code_6 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "for-loop variable initialization has side-effects"
	})
	lambdaSingleton_remove_unused_code_7 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "exported class"
	})
	lambdaSingleton_remove_unused_code_8 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "base class expression has side-effects"
	})
	lambdaSingleton_remove_unused_code_9 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "class body has side-effects"
	})
	lambdaSingleton_remove_unused_code_10 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "parameter in function using arguments"
	})
	lambdaSingleton_remove_unused_code_11 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "not removing globals"
	})
	lambdaSingleton_remove_unused_code_12 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "not removing locals"
	})
	lambdaSingleton_remove_unused_code_13 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return ""
	})
	lambdaSingleton_remove_unused_code_14 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "<pattern>"
	})
	lambdaSingleton_remove_unused_code_15 = j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return "hasNonLocalOrNonLiteralValue && requiresLocalLiteralValueForRemoval"
	})
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode
type RemoveUnusedCode struct {
	compiler                                                    AbstractCompiler                                        `gson:"compiler"`
	astAnalyzer                                                 *AstAnalyzer                                            `gson:"astAnalyzer"`
	codingConvention                                            CodingConvention                                        `gson:"codingConvention"`
	worklist                                                    j2gort.Deque[RemoveUnusedCode_Continuation]             `gson:"worklist"`
	varInfoKeys                                                 []*Var                                                  `gson:"varInfoKeys"`
	varInfoVals                                                 []*RemoveUnusedCode_VarInfo                             `gson:"varInfoVals"`
	varInfoTable                                                []int32                                                 `gson:"varInfoTable"`
	pinnedPropertyNames                                         j2gort.PlainSet[string]                                 `gson:"pinnedPropertyNames"`
	removablesForPropertyNames                                  j2gort.OrderedMap[string, []RemoveUnusedCode_Removable] `gson:"removablesForPropertyNames"`
	freeBuilder                                                 *RemoveUnusedCode_RemovableBuilder                      `gson:"freeBuilder"`
	canonicalUnremovableVarInfo                                 *RemoveUnusedCode_VarInfo                               `gson:"canonicalUnremovableVarInfo"`
	allFunctionParamScopes                                      []*Scope                                                `gson:"allFunctionParamScopes"`
	polyfills                                                   j2gort.Ext_com_google_common_collect_Multimap           `gson:"polyfills"`
	guardedUsages                                               j2gort.OrderedSet[pkg_rhino.Node]                       `gson:"guardedUsages"`
	polyfillsFromTable                                          *PolyfillUsageFinder_Polyfills                          `gson:"polyfillsFromTable"`
	polyfillUsageFinder                                         *PolyfillUsageFinder                                    `gson:"polyfillUsageFinder"`
	scopeCreator                                                *SyntacticScopeCreator                                  `gson:"scopeCreator"`
	removalLog                                                  pkg_diagnostic.LogFile                                  `gson:"removalLog"`
	unremovableLog                                              pkg_diagnostic.LogFile                                  `gson:"unremovableLog"`
	varInfoSize                                                 int32                                                   `gson:"varInfoSize"`
	idHash_                                                     atomic.Int32
	removeLocalVars                                             bool `gson:"removeLocalVars"`
	removeGlobals                                               bool `gson:"removeGlobals"`
	preserveFunctionExpressionNames                             bool `gson:"preserveFunctionExpressionNames"`
	removeUnusedPrototypeProperties                             bool `gson:"removeUnusedPrototypeProperties"`
	removeUnusedThisProperties                                  bool `gson:"removeUnusedThisProperties"`
	removeUnusedObjectDefinePropertiesDefinitions               bool `gson:"removeUnusedObjectDefinePropertiesDefinitions"`
	removeUnusedPolyfills                                       bool `gson:"removeUnusedPolyfills"`
	assumeGettersArePure                                        bool `gson:"assumeGettersArePure"`
	disambiguateProperties                                      bool `gson:"disambiguateProperties"`
	slab_RemoveUnusedCode_RemovableBuilder                      j2gort.PooledSlab[RemoveUnusedCode_RemovableBuilder]
	slab_RemoveUnusedCode_GlobalPolyfillInfo                    j2gort.PooledSlab[RemoveUnusedCode_GlobalPolyfillInfo]
	slab_RemoveUnusedCode_PrototypePropertyPolyfillInfo         j2gort.PooledSlab[RemoveUnusedCode_PrototypePropertyPolyfillInfo]
	slab_RemoveUnusedCode_StaticPropertyPolyfillInfo            j2gort.PooledSlab[RemoveUnusedCode_StaticPropertyPolyfillInfo]
	slab_RemoveUnusedCode_IndirectAssign                        j2gort.PooledSlab[RemoveUnusedCode_IndirectAssign]
	slab_RemoveUnusedCode_Polyfill                              j2gort.PooledSlab[RemoveUnusedCode_Polyfill]
	slab_RemoveUnusedCode_ClassDeclaration                      j2gort.PooledSlab[RemoveUnusedCode_ClassDeclaration]
	slab_RemoveUnusedCode_NamedClassExpression                  j2gort.PooledSlab[RemoveUnusedCode_NamedClassExpression]
	slab_RemoveUnusedCode_ClassOrPrototypeNamedProperty         j2gort.PooledSlab[RemoveUnusedCode_ClassOrPrototypeNamedProperty]
	slab_RemoveUnusedCode_ObjectDefinePropertiesDefinition      j2gort.PooledSlab[RemoveUnusedCode_ObjectDefinePropertiesDefinition]
	slab_RemoveUnusedCode_FunctionDeclaration                   j2gort.PooledSlab[RemoveUnusedCode_FunctionDeclaration]
	slab_RemoveUnusedCode_NameDeclarationStatement              j2gort.PooledSlab[RemoveUnusedCode_NameDeclarationStatement]
	slab_RemoveUnusedCode_Assign                                j2gort.PooledSlab[RemoveUnusedCode_Assign]
	slab_RemoveUnusedCode_ClassSetupCall                        j2gort.PooledSlab[RemoveUnusedCode_ClassSetupCall]
	slab_RemoveUnusedCode_VanillaForNameDeclaration             j2gort.PooledSlab[RemoveUnusedCode_VanillaForNameDeclaration]
	slab_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign j2gort.PooledSlab[RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign]
	slab_RemoveUnusedCode_IncOrDecOp                            j2gort.PooledSlab[RemoveUnusedCode_IncOrDecOp]
	slab_RemoveUnusedCode_UnusedReadReference                   j2gort.PooledSlab[RemoveUnusedCode_UnusedReadReference]
	slab_RemoveUnusedCode_InstanceofName                        j2gort.PooledSlab[RemoveUnusedCode_InstanceofName]
	slab_RemoveUnusedCode_WeakUsageCall                         j2gort.PooledSlab[RemoveUnusedCode_WeakUsageCall]
	slab_RemoveUnusedCode_VarInfo                               j2gort.PooledSlab[RemoveUnusedCode_VarInfo]
}

func (this_ *RemoveUnusedCode) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnusedCode) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode", "RemoveUnusedCode"
}
func (this_ *RemoveUnusedCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnusedCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode", "com.google.javascript.jscomp.RemoveUnusedCode", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode](v_)
	})
}
func (*RemoveUnusedCode) Implements_jscomp_CompilerPass() {}

var static_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_RemoveUnusedCode()
	return static_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES
}
func StaticSet_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_RemoveUnusedCode()
	static_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = v
}

var static_RemoveUnusedCode_UNREMOVABLE_REF_REASON j2gort.Ext_java_util_function_Supplier

func StaticGet_RemoveUnusedCode_UNREMOVABLE_REF_REASON() j2gort.Ext_java_util_function_Supplier {
	clinit_RemoveUnusedCode()
	return static_RemoveUnusedCode_UNREMOVABLE_REF_REASON
}
func StaticSet_RemoveUnusedCode_UNREMOVABLE_REF_REASON(v j2gort.Ext_java_util_function_Supplier) {
	clinit_RemoveUnusedCode()
	static_RemoveUnusedCode_UNREMOVABLE_REF_REASON = v
}
func StaticGet_RemoveUnusedCode_DOT_PROTOTYPE() string { return ".prototype" }

var clinitDone_RemoveUnusedCode uint32
var clinitState_RemoveUnusedCode j2gort.ClinitState

func clinit_RemoveUnusedCode() {
	if atomic.LoadUint32(&clinitDone_RemoveUnusedCode) != 0 {
		return
	}
	clinitSlow_RemoveUnusedCode()
}

func clinitSlow_RemoveUnusedCode() {
	if !j2gort.EnterClinit(&clinitState_RemoveUnusedCode, &clinitDone_RemoveUnusedCode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RemoveUnusedCode, &clinitDone_RemoveUnusedCode)
	static_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_5("length", "toString", "valueOf", "constructor", "prototype")
	static_RemoveUnusedCode_UNREMOVABLE_REF_REASON = lambdaSingleton_remove_unused_code_1
}

func (this_ *RemoveUnusedCode) initFields_RemoveUnusedCode() {
	_ = this_
	this_.worklist = j2gort.Deque[RemoveUnusedCode_Continuation]{}
	this_.varInfoKeys = j2gort.AcquirePooledSlice[*Var](256)[:256]
	this_.varInfoVals = j2gort.AcquirePooledSlice[*RemoveUnusedCode_VarInfo](256)[:256]
	this_.varInfoTable = j2gort.AcquirePooledSlice[int32](512)[:512]

	this_.pinnedPropertyNames = j2gort.NewPlainSetCap[string](int32(16))
	this_.removablesForPropertyNames = j2gort.NewOrderedMap[string, []RemoveUnusedCode_Removable]()
	this_.allFunctionParamScopes = []*Scope{}
	this_.polyfills = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_HashMultimap_Create_0())
	this_.guardedUsages = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: RemoveUnusedCode(com.google.javascript.jscomp.RemoveUnusedCode.Builder builder)
func NewRemoveUnusedCode(builder *RemoveUnusedCode_Builder) *RemoveUnusedCode {
	this_ := &RemoveUnusedCode{}
	InitRemoveUnusedCode(this_, builder)
	return this_
}

func InitRemoveUnusedCode(this_ *RemoveUnusedCode, builder *RemoveUnusedCode_Builder) {
	_ = builder
	this_.initFields_RemoveUnusedCode()
	this_.compiler = builder.compiler
	this_.astAnalyzer = this_.compiler.GetAstAnalyzer()
	this_.codingConvention = builder.compiler.GetCodingConvention()
	this_.scopeCreator = NewSyntacticScopeCreator_1(builder.compiler)
	this_.removeLocalVars = builder.removeLocalVars
	this_.removeGlobals = builder.removeGlobals
	this_.preserveFunctionExpressionNames = builder.preserveFunctionExpressionNames
	this_.removeUnusedPrototypeProperties = builder.removeUnusedPrototypeProperties
	this_.removeUnusedThisProperties = builder.removeUnusedThisProperties
	this_.removeUnusedObjectDefinePropertiesDefinitions = builder.removeUnusedObjectDefinePropertiesDefinitions
	this_.removeUnusedPolyfills = builder.removeUnusedPolyfills
	this_.disambiguateProperties = builder.disambiguateProperties
	this_.polyfillsFromTable = PolyfillUsageFinder_Polyfills_FromTable(pkg_resources.ResourceLoader_LoadTextResource(j2gort.ClassLiteral("com.google.javascript.jscomp.RemoveUnusedCode", "RemoveUnusedCode"), "js/polyfills.txt"))
	this_.polyfillUsageFinder = NewPolyfillUsageFinder(this_.compiler, this_.polyfillsFromTable, this_.disambiguateProperties)
	this_.assumeGettersArePure = builder.assumeGettersArePure
	this_.canonicalUnremovableVarInfo = NewRemoveUnusedCode_VarInfo(this_)
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_RemoveUnusedCode_IMPLICITLY_USED_PROPERTIES()))
		for cur_1.HasNext() {
			prop := cur_1.NextStr()
			_ = prop
			this_.pinnedPropertyNames.AddVoid(prop)
		}
	}
}

// Java method: RemoveUnusedCode.acquireBuilder()
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_acquireBuilder_4663() *RemoveUnusedCode_RemovableBuilder {
	var b *RemoveUnusedCode_RemovableBuilder = this_.freeBuilder
	_ = b
	if b != nil {
		this_.freeBuilder = nil
		return b
	}
	return allocPooled_RemoveUnusedCode_RemovableBuilder(this_)
}

// Java method: RemoveUnusedCode.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RemoveUnusedCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.ReleaseSlabs_()
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetExternProperties()))
		for cur_2.HasNext() {
			prop := cur_2.NextStr()
			_ = prop
			this_.pinnedPropertyNames.AddVoid(prop)
		}
	}
	func() {
		defer func() {
			this_.removalLog = nil
			this_.unremovableLog = nil
		}()
		var removalLogFile pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "removals.log", []string{})
		_ = removalLogFile
		defer j2gort.CloseResource(removalLogFile)
		var keepLogFile pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "unremovable.log", []string{})
		_ = keepLogFile
		defer j2gort.CloseResource(keepLogFile)
		this_.removalLog = removalLogFile
		this_.unremovableLog = keepLogFile
		this_.priv_RemoveUnusedCode_traverseAndRemoveUnusedReferences_4664(root)
	}()
}

// Java method: RemoveUnusedCode.traverseAndRemoveUnusedReferences(com.google.javascript.rhino.Node root)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseAndRemoveUnusedReferences_4664(root pkg_rhino.Node) {
	_ = root
	var scope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(root.GetParent(), nil)
	_ = scope
	if !(scope.HasSlot("JSCompiler_renameProperty")) {
		scope.Declare("JSCompiler_renameProperty", pkg_rhino.Node(0), nil)
	}
	this_.polyfillUsageFinder.TraverseOnlyGuarded(root, func() j2gort.Ext_java_util_function_Consumer {
		bound_3 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_3.priv_RemoveUnusedCode_storePolyfill_4665(j2gort.CastPtr[PolyfillUsageFinder_PolyfillUsage](mrp0_))
		})
	}())
	this_.worklist.AddLast(NewRemoveUnusedCode_Continuation(this_, root, scope))
	for !(this_.worklist.IsEmpty()) {
		var continuation RemoveUnusedCode_Continuation = this_.worklist.RemoveFirst()
		_ = continuation
		continuation.Apply()
	}
	this_.priv_RemoveUnusedCode_removeUnreferencedVarsAndPolyfills_4710()
	this_.priv_RemoveUnusedCode_removeIndependentlyRemovableProperties_4666()
	for _, fparamScope := range this_.allFunctionParamScopes {
		_ = fparamScope
		this_.priv_RemoveUnusedCode_removeUnreferencedFunctionArgs_4700(fparamScope)
	}

}

// Java method: RemoveUnusedCode.storePolyfill(com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage polyfillUsage)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_storePolyfill_4665(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	_ = polyfillUsage
	this_.guardedUsages.Add(polyfillUsage.Node())
}

// Java method: RemoveUnusedCode.removeIndependentlyRemovableProperties()
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_removeIndependentlyRemovableProperties_4666() {
	{
		cur_6 := this_.removablesForPropertyNames.IterEntries()
		for cur_6.NextEntry() {
			var propName string = cur_6.Key()
			_ = propName
			if this_.removalLog.IsLogging() {
				this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForProperty(propName)))
			}
			var list []RemoveUnusedCode_Removable = cur_6.Value()
			_ = list
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = int32(len(list))
				_ = n
				for ; i < n; i = i + 1 {
					list[i].GetRemoveUnusedCode_RemovableBase().Remove(this_.compiler)
				}
			}
		}
	}
}

// Java method: RemoveUnusedCode.traverseNode(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseNode_4667(n pkg_rhino.Node, scope *Scope) {
	_ = n
	_ = scope
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var type_ int32 = n.TokenSwitchOrd_()
	switch type_ {
	case int32(86) /* Token.CATCH */ :
		this_.priv_RemoveUnusedCode_traverseCatch_4682(n, scope)
	case int32(73) /* Token.FUNCTION */ :
		{
			var varInfo *RemoveUnusedCode_VarInfo = nil
			_ = varInfo
			if NodeUtil_IsFunctionDeclaration(n) {
				varInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(n.GetFirstChild(), scope)
				var functionDeclaration *RemoveUnusedCode_FunctionDeclaration = this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, n, scope)).BuildFunctionDeclaration(n)
				_ = functionDeclaration
				varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(functionDeclaration))
				if parent.IsExport() {
					varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_2)
				}
			} else {
				this_.priv_RemoveUnusedCode_traverseFunction_4698(n, scope)
			}
		}
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_RemoveUnusedCode_traverseAssign_4687(n, scope)
	case int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */ :
		this_.priv_RemoveUnusedCode_traverseCompoundAssign_4671(n, scope)
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		this_.priv_RemoveUnusedCode_traverseIncrementOrDecrementOp_4670(n, scope)
	case int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */ :
		this_.priv_RemoveUnusedCode_traverseCall_4673(n, scope)
	case int32(91) /* Token.SWITCH_BODY */, int32(90) /* Token.BLOCK */ :
		this_.priv_RemoveUnusedCode_traverseChildren_4693(n, func() *Scope {
			if NodeUtil_CreatesBlockScope(n) {
				return this_.scopeCreator.CreateScope_2_Node_Scope(n, scope)
			}
			return scope
		}())
	case int32(123) /* Token.MODULE_BODY */ :
		this_.priv_RemoveUnusedCode_traverseChildren_4693(n, this_.scopeCreator.CreateScope_2_Node_Scope(n, scope))
	case int32(105) /* Token.CLASS */ :
		this_.priv_RemoveUnusedCode_traverseClass_4694(n, scope)
	case int32(106) /* Token.CLASS_MEMBERS */ :
		this_.priv_RemoveUnusedCode_traverseClassMembers_4697(n, scope)
	case int32(102) /* Token.ARRAY_PATTERN */, int32(49) /* Token.PARAM_LIST */ :
		this_.priv_RemoveUnusedCode_traverseIndirectAssignmentList_4691(n, scope)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_RemoveUnusedCode_traverseObjectPattern_4690(n, scope)
	case int32(47) /* Token.OBJECTLIT */ :
		this_.priv_RemoveUnusedCode_traverseObjectLiteral_4677(n, scope)
	case int32(80) /* Token.FOR */ :
		this_.priv_RemoveUnusedCode_traverseVanillaFor_4684(n, scope)
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_RemoveUnusedCode_traverseEnhancedFor_4683(n, scope)
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		if !(NodeUtil_IsStatement(n)) {
			j2gort.CheckState(false)
		}
		this_.priv_RemoveUnusedCode_traverseDeclarationStatement_4686(n, scope)
	case int32(45) /* Token.INSTANCEOF */ :
		this_.priv_RemoveUnusedCode_traverseInstanceof_4668(n, scope)
	case int32(32) /* Token.NAME */ :
		if !(!(n.HasChildren())) {
			j2gort.CheckState(false)
		}
		if !(parent.IsParamList()) {
			if !(!(NodeUtil_IsNameDeclaration(parent))) {
				j2gort.CheckState(false)
			}
			if !(!((parent.IsFunction() || parent.IsClass()) && (parent.GetFirstChild() == n))) {
				j2gort.CheckState(false)
			}
			this_.priv_RemoveUnusedCode_traverseNameNode_4672(n, scope).SetIsExplicitlyNotRemovable(func() j2gort.Ext_java_util_function_Supplier {
				if this_.unremovableLog.IsLogging() {
					return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
						return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("reference found: %s", []any{j2gort.BoxString(n.GetLocation())}))
					})
				}
				return StaticGet_RemoveUnusedCode_UNREMOVABLE_REF_REASON()
			}())
		}
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_RemoveUnusedCode_traverseNormalOrOptChainGetProp_4669(n, scope)
	default:
		this_.priv_RemoveUnusedCode_traverseChildren_4693(n, scope)
	}
}

// Java method: RemoveUnusedCode.traverseInstanceof(com.google.javascript.rhino.Node instanceofNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseInstanceof_4668(instanceofNode pkg_rhino.Node, scope *Scope) {
	_ = instanceofNode
	_ = scope
	if !(instanceofNode.IsInstanceOf()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(instanceofNode))
	}
	var lhs pkg_rhino.Node = instanceofNode.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = lhs.GetNext()
	_ = rhs
	this_.priv_RemoveUnusedCode_traverseNode_4667(lhs, scope)
	if rhs.IsName() {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(rhs, scope)
		_ = varInfo
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
		_ = builder
		varInfo.AddRemovable(builder.BuildInstanceofName(instanceofNode))
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(rhs, scope)
	}
}

// Java method: RemoveUnusedCode.traverseNormalOrOptChainGetProp(com.google.javascript.rhino.Node getProp, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseNormalOrOptChainGetProp_4669(getProp pkg_rhino.Node, scope *Scope) {
	_ = getProp
	_ = scope
	if !(NodeUtil_IsNormalOrOptChainGetProp(getProp)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(getProp))
	}
	var objectNode pkg_rhino.Node = getProp.GetFirstChild()
	_ = objectNode
	var propertyName string = getProp.GetString()
	_ = propertyName
	if j2gort.MultimapContainsKeyStr(this_.polyfills, propertyName) {
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetStr(this_.polyfills, propertyName))))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				info := j2gort.AsIface[RemoveUnusedCode_PolyfillInfo](rawElem_8)
				_ = info
				if info.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
					info.GetRemoveUnusedCode_PolyfillInfoBase().ConsiderPossibleReference(getProp)
				}
			}
		}
	}
	if NodeUtil_IsExpressionResultUsed(getProp) || this_.priv_RemoveUnusedCode_considerForAccessorSideEffects_4703(getProp, StaticGet_AccessorSummary_PropertyAccessKind_GETTER_ONLY()) {
		this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propertyName)
		this_.priv_RemoveUnusedCode_traverseNode_4667(objectNode, scope)
	} else if objectNode.IsThis() {
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsThisDotPropertyReference(true)
		_ = builder
		this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder.BuildUnusedReadReference(getProp, getProp)))
	} else if RemoveUnusedCode_IsDotPrototype(objectNode) {
		var builder_2 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsPrototypeDotPropertyReference(true)
		_ = builder_2
		var objExpression pkg_rhino.Node = objectNode.GetFirstChild()
		_ = objExpression
		if objExpression.IsName() {
			var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(objExpression, scope)
			_ = varInfo
			varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder_2.BuildUnusedReadReference(getProp, getProp)))
		} else {
			if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(objExpression) {
				this_.priv_RemoveUnusedCode_traverseNode_4667(objExpression, scope)
			} else {
				builder_2.AddContinuation(NewRemoveUnusedCode_Continuation(this_, objExpression, scope))
			}
			this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder_2.BuildUnusedReadReference(getProp, getProp)))
		}
	} else {
		this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propertyName)
		this_.priv_RemoveUnusedCode_traverseNode_4667(objectNode, scope)
	}
}

// Java method: RemoveUnusedCode.traverseIncrementOrDecrementOp(com.google.javascript.rhino.Node incOrDecOp, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseIncrementOrDecrementOp_4670(incOrDecOp pkg_rhino.Node, scope *Scope) {
	_ = incOrDecOp
	_ = scope
	if !(incOrDecOp.IsInc() || incOrDecOp.IsDec()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(incOrDecOp))
	}
	var arg pkg_rhino.Node = incOrDecOp.GetOnlyChild()
	_ = arg
	if NodeUtil_IsExpressionResultUsed(incOrDecOp) {
		this_.priv_RemoveUnusedCode_traverseNode_4667(arg, scope)
	} else if arg.IsGetProp() {
		var getPropObj pkg_rhino.Node = arg.GetFirstChild()
		_ = getPropObj
		if this_.priv_RemoveUnusedCode_considerForAccessorSideEffects_4703(arg, StaticGet_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER()) {
			this_.priv_RemoveUnusedCode_traverseNode_4667(getPropObj, scope)
		} else if getPropObj.IsThis() {
			var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsThisDotPropertyReference(true)
			_ = builder
			this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder.BuildIncOrDepOp(incOrDecOp, arg, pkg_rhino.Node(0))))
		} else if RemoveUnusedCode_IsDotPrototype(getPropObj) {
			var exprObj pkg_rhino.Node = getPropObj.GetFirstChild()
			_ = exprObj
			var builder_2 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsPrototypeDotPropertyReference(true)
			_ = builder_2
			if exprObj.IsName() {
				var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(exprObj, scope)
				_ = varInfo
				varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder_2.BuildIncOrDepOp(incOrDecOp, arg, pkg_rhino.Node(0))))
			} else {
				var toPreserve pkg_rhino.Node = pkg_rhino.Node(0)
				_ = toPreserve
				if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(exprObj) {
					toPreserve = exprObj
					this_.priv_RemoveUnusedCode_traverseNode_4667(exprObj, scope)
				} else {
					builder_2.AddContinuation(NewRemoveUnusedCode_Continuation(this_, exprObj, scope))
				}
				this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder_2.BuildIncOrDepOp(incOrDecOp, arg, toPreserve)))
			}
		} else {
			this_.priv_RemoveUnusedCode_traverseNode_4667(arg, scope)
		}
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(arg, scope)
	}
}

// Java method: RemoveUnusedCode.traverseCompoundAssign(com.google.javascript.rhino.Node compoundAssignNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseCompoundAssign_4671(compoundAssignNode pkg_rhino.Node, scope *Scope) {
	_ = compoundAssignNode
	_ = scope
	var targetNode pkg_rhino.Node = compoundAssignNode.GetFirstChild()
	_ = targetNode
	var valueNode pkg_rhino.Node = compoundAssignNode.GetLastChild()
	_ = valueNode
	if targetNode.IsGetProp() {
		if this_.priv_RemoveUnusedCode_considerForAccessorSideEffects_4703(targetNode, StaticGet_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER()) {
			this_.priv_RemoveUnusedCode_traverseNode_4667(targetNode.GetFirstChild(), scope)
			this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
		} else if targetNode.GetFirstChild().IsThis() && !(NodeUtil_IsExpressionResultUsed(compoundAssignNode)) {
			var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsThisDotPropertyReference(true)
			_ = builder
			this_.priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode, builder, scope)
			this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder.BuildNamedPropertyAssign_2(compoundAssignNode, targetNode)))
		} else {
			this_.priv_RemoveUnusedCode_traverseNode_4667(targetNode, scope)
			this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
		}
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(targetNode, scope)
		this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
	}
}

// Java method: RemoveUnusedCode.traverseNameNode(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseNameNode_4672(n pkg_rhino.Node, scope *Scope) *RemoveUnusedCode_VarInfo {
	_ = n
	_ = scope
	if j2gort.MultimapContainsKeyStr(this_.polyfills, n.GetString()) {
		{
			cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetStr(this_.polyfills, n.GetString()))))
			for cur_9.HasNext() {
				rawElem_10 := cur_9.Next()
				info := j2gort.AsIface[RemoveUnusedCode_PolyfillInfo](rawElem_10)
				_ = info
				if info.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
					info.GetRemoveUnusedCode_PolyfillInfoBase().ConsiderPossibleReference(n)
				}
			}
		}
	}
	return this_.priv_RemoveUnusedCode_traverseVar_4707(this_.priv_RemoveUnusedCode_getVarForNameNode_4676(n, scope))
}

// Java method: RemoveUnusedCode.traverseCall(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseCall_4673(callNode pkg_rhino.Node, scope *Scope) {
	_ = callNode
	_ = scope
	var callee pkg_rhino.Node = callNode.GetFirstChild()
	_ = callee
	if this_.codingConvention.IsPropertyRenameFunction(callee) {
		var propertyNameNode pkg_rhino.Node = callee.GetNext()
		_ = propertyNameNode
		if (propertyNameNode != 0) && propertyNameNode.IsStringLit() {
			this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propertyNameNode.GetString())
		}
		this_.priv_RemoveUnusedCode_traverseChildren_4693(callNode, scope)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(callNode) {
		this_.priv_RemoveUnusedCode_traverseObjectDefinePropertiesCall_4674(callNode, scope)
	} else if this_.removeUnusedPolyfills && RemoveUnusedCode_IsJscompPolyfill(callee) {
		var firstArg pkg_rhino.Node = callee.GetNext()
		_ = firstArg
		var polyfillName string = firstArg.GetString()
		_ = polyfillName
		if javalang.StringEndsWith(j2gort.CheckNotNullString(callee.GetString()), "polyfillTypedArrayMethod") {
			polyfillName = ("TypedArray.prototype." + j2gort.StrValString(polyfillName))
		}
		var info RemoveUnusedCode_PolyfillInfo = this_.priv_RemoveUnusedCode_createPolyfillInfo_4712(callNode, scope, polyfillName)
		_ = info
		j2gort.MultimapPutStr(this_.polyfills, info.GetRemoveUnusedCode_PolyfillInfoBase().key, info)
		this_.priv_RemoveUnusedCode_traverseNode_4667(callNode.GetFirstChild(), scope)
	} else if (NodeUtil_IsGoogWeakUsageCall(callNode) && callNode.HasTwoChildren()) && callNode.GetSecondChild().IsName() {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(callNode.GetSecondChild(), scope)
		_ = varInfo
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
		_ = builder
		varInfo.AddRemovable(builder.BuildWeakUsageCall(callNode))
		this_.priv_RemoveUnusedCode_traverseNode_4667(callNode.GetFirstChild(), scope)
	} else {
		var parent pkg_rhino.Node = callNode.GetParent()
		_ = parent
		var classVarName string = j2gort.NullString
		_ = classVarName
		var classDefiningCall bool = false
		_ = classDefiningCall
		if parent.IsExprResult() || (parent.IsComma() && (parent.GetFirstChild() == callNode)) {
			var subclassRelationship *CodingConvention_SubclassRelationship = this_.codingConvention.GetClassesDefinedByCall(callNode)
			_ = subclassRelationship
			if subclassRelationship != nil {
				classVarName = subclassRelationship.subclassName
				classDefiningCall = true
			} else {
				classVarName = this_.codingConvention.GetSingletonGetterClassName(callNode)
			}
		}
		var classVar *Var = nil
		_ = classVar
		if (classVarName != j2gort.NullString) && NodeUtil_IsValidSimpleName(classVarName) {
			classVar = j2gort.CheckNotNullPtr(scope.LookupVar(classVarName))
		}
		if (classVar == nil) || !(classVar.IsGlobal()) {
			this_.priv_RemoveUnusedCode_traverseChildren_4693(callNode, scope)
		} else {
			var builder_2 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
			_ = builder_2
			{
				var child pkg_rhino.Node = callNode.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					builder_2.AddContinuation(NewRemoveUnusedCode_Continuation(this_, child, scope))
				}
			}
			this_.priv_RemoveUnusedCode_traverseVar_4707(classVar).AddRemovable(ToIface_RemoveUnusedCode_Removable(builder_2.BuildClassSetupCall_2(callNode, classDefiningCall)))
		}
	}
}

// Java static method: RemoveUnusedCode.isJscompPolyfill(com.google.javascript.rhino.Node n)
func RemoveUnusedCode_IsJscompPolyfill(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var name string = n.GetString()
		_ = name
		return ((((j2gort.CheckNotNullString(name) == "$jscomp$polyfill") || (j2gort.CheckNotNullString(name) == "$jscomp$patch")) || (j2gort.CheckNotNullString(name) == "$jscomp$polyfillTypedArrayMethod")) && n.GetNext().IsStringLit())
	case int32(26) /* Token.GETPROP */ :
		var propertyName string = n.GetString()
		_ = propertyName
		return ((((((j2gort.CheckNotNullString(propertyName) == "polyfill") || (j2gort.CheckNotNullString(propertyName) == "patch")) || (j2gort.CheckNotNullString(propertyName) == "polyfillTypedArrayMethod")) && n.GetFirstChild().IsName()) && (j2gort.CheckNotNullString(n.GetFirstChild().GetString()) == "$jscomp")) && n.GetNext().IsStringLit())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: RemoveUnusedCode.traverseObjectDefinePropertiesCall(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseObjectDefinePropertiesCall_4674(callNode pkg_rhino.Node, scope *Scope) {
	_ = callNode
	_ = scope
	var callee pkg_rhino.Node = callNode.GetFirstChild()
	_ = callee
	var targetObject pkg_rhino.Node = callNode.GetSecondChild()
	_ = targetObject
	var propertyDefinitions pkg_rhino.Node = targetObject.GetNext()
	_ = propertyDefinitions
	if (targetObject.IsName() || this_.priv_RemoveUnusedCode_isNameDotPrototype_4689(targetObject)) && !(NodeUtil_IsExpressionResultUsed(callNode)) {
		var nameNode pkg_rhino.Node
		_ = nameNode
		if targetObject.IsName() {
			nameNode = targetObject
		} else {
			nameNode = targetObject.GetFirstChild()
		}
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(nameNode, scope)
		_ = varInfo
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
		_ = builder
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(this_, callee, scope))
		if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(propertyDefinitions) {
			this_.priv_RemoveUnusedCode_traverseNode_4667(propertyDefinitions, scope)
		} else {
			builder.AddContinuation(NewRemoveUnusedCode_Continuation(this_, propertyDefinitions, scope))
		}
		varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder.BuildClassSetupCall_1(callNode)))
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(callee, scope)
		this_.priv_RemoveUnusedCode_traverseNode_4667(targetObject, scope)
		this_.priv_RemoveUnusedCode_traverseNode_4667(propertyDefinitions, scope)
	}
}

// Java method: RemoveUnusedCode.traverseObjectDefinePropertiesLiteral(com.google.javascript.rhino.Node propertyDefinitions, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseObjectDefinePropertiesLiteral_4675(propertyDefinitions pkg_rhino.Node, scope *Scope) {
	_ = propertyDefinitions
	_ = scope
	{
		var property pkg_rhino.Node = propertyDefinitions.GetFirstChild()
		_ = property
		for ; property != 0; property = property.GetNext() {
			if property.IsQuotedStringKey() {
				this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(property.GetString())
				this_.priv_RemoveUnusedCode_traverseNode_4667(property.GetOnlyChild(), scope)
			} else if property.IsStringKey() {
				var definition pkg_rhino.Node = property.GetOnlyChild()
				_ = definition
				if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(definition) {
					this_.priv_RemoveUnusedCode_traverseNode_4667(definition, scope)
				} else {
					this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, definition, scope)).BuildObjectDefinePropertiesDefinition(property)))
				}
			} else {
				this_.priv_RemoveUnusedCode_traverseNode_4667(property, scope)
			}
		}
	}
}

// Java method: RemoveUnusedCode.getVarForNameNode(com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_getVarForNameNode_4676(nameNode pkg_rhino.Node, scope *Scope) *Var {
	_ = nameNode
	_ = scope
	return j2gort.CheckNotNullPtr(scope.LookupVar(nameNode.GetString()))
}

// Java method: RemoveUnusedCode.traverseObjectLiteral(com.google.javascript.rhino.Node objectLiteral, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseObjectLiteral_4677(objectLiteral pkg_rhino.Node, scope *Scope) {
	_ = objectLiteral
	_ = scope
	if !(objectLiteral.IsObjectLit()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(objectLiteral))
	}
	if this_.priv_RemoveUnusedCode_isAssignmentToPrototype_4681(objectLiteral.GetParent()) {
		this_.priv_RemoveUnusedCode_traversePrototypeLiteral_4680(objectLiteral, scope)
	} else if this_.priv_RemoveUnusedCode_isObjectDefinePropertiesSecondArgument_4678(objectLiteral) {
		this_.priv_RemoveUnusedCode_traverseObjectDefinePropertiesLiteral_4675(objectLiteral, scope)
	} else {
		this_.priv_RemoveUnusedCode_traverseNonPrototypeObjectLiteral_4679(objectLiteral, scope)
	}
}

// Java method: RemoveUnusedCode.isObjectDefinePropertiesSecondArgument(com.google.javascript.rhino.Node n)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_isObjectDefinePropertiesSecondArgument_4678(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (NodeUtil_IsObjectDefinePropertiesDefinition(parent) && (parent.GetLastChild() == n))
}

// Java method: RemoveUnusedCode.traverseNonPrototypeObjectLiteral(com.google.javascript.rhino.Node objectLiteral, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseNonPrototypeObjectLiteral_4679(objectLiteral pkg_rhino.Node, scope *Scope) {
	_ = objectLiteral
	_ = scope
	{
		var propertyNode pkg_rhino.Node = objectLiteral.GetFirstChild()
		_ = propertyNode
		for ; propertyNode != 0; propertyNode = propertyNode.GetNext() {
			if propertyNode.IsStringKey() {
				this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propertyNode.GetString())
				this_.priv_RemoveUnusedCode_traverseNode_4667(propertyNode.GetFirstChild(), scope)
			} else {
				this_.priv_RemoveUnusedCode_traverseNode_4667(propertyNode, scope)
			}
		}
	}
}

// Java method: RemoveUnusedCode.traversePrototypeLiteral(com.google.javascript.rhino.Node objectLiteral, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traversePrototypeLiteral_4680(objectLiteral pkg_rhino.Node, scope *Scope) {
	_ = objectLiteral
	_ = scope
	{
		var propertyNode pkg_rhino.Node = objectLiteral.GetFirstChild()
		_ = propertyNode
		for ; propertyNode != 0; propertyNode = propertyNode.GetNext() {
			if propertyNode.IsComputedProp() || propertyNode.IsQuotedStringKey() {
				this_.priv_RemoveUnusedCode_traverseChildren_4693(propertyNode, scope)
			} else {
				var valueNode pkg_rhino.Node = propertyNode.GetOnlyChild()
				_ = valueNode
				if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(valueNode) {
					this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
				} else {
					this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, valueNode, scope)).BuildClassOrPrototypeNamedProperty(propertyNode)))
				}
			}
		}
	}
}

// Java method: RemoveUnusedCode.isAssignmentToPrototype(com.google.javascript.rhino.Node n)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_isAssignmentToPrototype_4681(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsAssign() && RemoveUnusedCode_IsDotPrototype(n.GetFirstChild()))
}

// Java static method: RemoveUnusedCode.isDotPrototype(com.google.javascript.rhino.Node n)
func RemoveUnusedCode_IsDotPrototype(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsNormalOrOptChainGetProp(n) && (j2gort.CheckNotNullString(n.GetString()) == "prototype"))
}

// Java method: RemoveUnusedCode.traverseCatch(com.google.javascript.rhino.Node catchNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseCatch_4682(catchNode pkg_rhino.Node, scope *Scope) {
	_ = catchNode
	_ = scope
	var exceptionNameNode pkg_rhino.Node = catchNode.GetFirstChild()
	_ = exceptionNameNode
	var block pkg_rhino.Node = exceptionNameNode.GetNext()
	_ = block
	if exceptionNameNode.IsName() {
		var exceptionVarInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(exceptionNameNode, scope)
		_ = exceptionVarInfo
		exceptionVarInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_3)
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(exceptionNameNode, scope)
	}
	this_.priv_RemoveUnusedCode_traverseNode_4667(block, scope)
}

// Java method: RemoveUnusedCode.traverseEnhancedFor(com.google.javascript.rhino.Node enhancedFor, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseEnhancedFor_4683(enhancedFor pkg_rhino.Node, scope *Scope) {
	_ = enhancedFor
	_ = scope
	var forScope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(enhancedFor, scope)
	_ = forScope
	var iterationTarget pkg_rhino.Node = enhancedFor.GetFirstChild()
	_ = iterationTarget
	var collection pkg_rhino.Node = iterationTarget.GetNext()
	_ = collection
	var body pkg_rhino.Node = collection.GetNext()
	_ = body
	if iterationTarget.IsName() {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(iterationTarget, forScope)
		_ = varInfo
		varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_4)
	} else if NodeUtil_IsNameDeclaration(iterationTarget) {
		var declNode pkg_rhino.Node = iterationTarget.GetOnlyChild()
		_ = declNode
		if declNode.IsDestructuringLhs() {
			this_.priv_RemoveUnusedCode_traverseNode_4667(declNode, forScope)
		} else {
			if !(declNode.IsName()) {
				j2gort.CheckState(false)
			}
			if !(!(declNode.HasChildren())) {
				j2gort.CheckState(false)
			}
			var varInfo_2 *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(declNode, forScope)
			_ = varInfo_2
			varInfo_2.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_5)
		}
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(iterationTarget, forScope)
	}
	this_.priv_RemoveUnusedCode_traverseNode_4667(collection, forScope)
	this_.priv_RemoveUnusedCode_traverseNode_4667(body, forScope)
}

// Java method: RemoveUnusedCode.traverseVanillaFor(com.google.javascript.rhino.Node forNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseVanillaFor_4684(forNode pkg_rhino.Node, scope *Scope) {
	_ = forNode
	_ = scope
	var forScope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(forNode, scope)
	_ = forScope
	var initialization pkg_rhino.Node = forNode.GetFirstChild()
	_ = initialization
	var condition pkg_rhino.Node = initialization.GetNext()
	_ = condition
	var update pkg_rhino.Node = condition.GetNext()
	_ = update
	var block pkg_rhino.Node = update.GetNext()
	_ = block
	if NodeUtil_IsNameDeclaration(initialization) {
		this_.priv_RemoveUnusedCode_traverseVanillaForNameDeclarations_4685(initialization, forScope)
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(initialization, forScope)
	}
	this_.priv_RemoveUnusedCode_traverseNode_4667(condition, forScope)
	this_.priv_RemoveUnusedCode_traverseNode_4667(update, forScope)
	this_.priv_RemoveUnusedCode_traverseNode_4667(block, forScope)
}

// Java method: RemoveUnusedCode.traverseVanillaForNameDeclarations(com.google.javascript.rhino.Node nameDeclaration, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseVanillaForNameDeclarations_4685(nameDeclaration pkg_rhino.Node, scope *Scope) {
	_ = nameDeclaration
	_ = scope
	{
		var child pkg_rhino.Node = nameDeclaration.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if !(child.IsName()) {
				this_.priv_RemoveUnusedCode_traverseNode_4667(child, scope)
			} else {
				var nameNode pkg_rhino.Node = child
				_ = nameNode
				var valueNode pkg_rhino.Node = child.GetFirstChild()
				_ = valueNode
				var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(nameNode, scope)
				_ = varInfo
				if valueNode == 0 {
					varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(this_.priv_RemoveUnusedCode_acquireBuilder_4663().BuildVanillaForNameDeclaration(nameNode)))
				} else if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(valueNode) {
					varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_6)
					this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
				} else {
					var vanillaForNameDeclaration *RemoveUnusedCode_VanillaForNameDeclaration = this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, valueNode, scope)).BuildVanillaForNameDeclaration(nameNode)
					_ = vanillaForNameDeclaration
					varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(vanillaForNameDeclaration))
				}
			}
		}
	}
}

// Java method: RemoveUnusedCode.traverseDeclarationStatement(com.google.javascript.rhino.Node declarationStatement, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseDeclarationStatement_4686(declarationStatement pkg_rhino.Node, scope *Scope) {
	_ = declarationStatement
	_ = scope
	var nameNode pkg_rhino.Node = declarationStatement.GetOnlyChild()
	_ = nameNode
	if !(nameNode.IsName()) {
		this_.priv_RemoveUnusedCode_traverseNode_4667(nameNode, scope)
	} else {
		var valueNode pkg_rhino.Node = nameNode.GetFirstChild()
		_ = valueNode
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(nameNode, scope)
		_ = varInfo
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
		_ = builder
		if valueNode == 0 {
			varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder.BuildNameDeclarationStatement(declarationStatement)))
		} else {
			if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(valueNode) {
				this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
			} else {
				builder.AddContinuation(NewRemoveUnusedCode_Continuation(this_, valueNode, scope))
			}
			var removable *RemoveUnusedCode_NameDeclarationStatement = builder.BuildNameDeclarationStatement(declarationStatement)
			_ = removable
			varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(removable))
		}
	}
}

// Java method: RemoveUnusedCode.traverseAssign(com.google.javascript.rhino.Node assignNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseAssign_4687(assignNode pkg_rhino.Node, scope *Scope) {
	_ = assignNode
	_ = scope
	if !(NodeUtil_IsAssignmentOp(assignNode)) {
		j2gort.CheckState(false)
	}
	var lhs pkg_rhino.Node = assignNode.GetFirstChild()
	_ = lhs
	var valueNode pkg_rhino.Node = assignNode.GetLastChild()
	_ = valueNode
	if lhs.IsName() {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(lhs, scope)
		_ = varInfo
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
		_ = builder
		this_.priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode, builder, scope)
		varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder.BuildVariableAssign(assignNode, varInfo)))
	} else if lhs.IsGetElem() {
		var getElemObj pkg_rhino.Node = lhs.GetFirstChild()
		_ = getElemObj
		var getElemKey pkg_rhino.Node = lhs.GetLastChild()
		_ = getElemKey
		var varNameNode pkg_rhino.Node
		_ = varNameNode
		if getElemObj.IsName() {
			varNameNode = getElemObj
		} else if this_.priv_RemoveUnusedCode_isNameDotPrototype_4689(getElemObj) {
			varNameNode = getElemObj.GetFirstChild()
		} else {
			varNameNode = pkg_rhino.Node(0)
		}
		if varNameNode != 0 {
			var varInfo_2 *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(varNameNode, scope)
			_ = varInfo_2
			var builder_2 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
			_ = builder_2
			if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(getElemKey) {
				this_.priv_RemoveUnusedCode_traverseNode_4667(getElemKey, scope)
			} else {
				builder_2.AddContinuation(NewRemoveUnusedCode_Continuation(this_, getElemKey, scope))
			}
			this_.priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode, builder_2, scope)
			varInfo_2.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder_2.BuildComputedPropertyAssign(assignNode, getElemKey, varInfo_2)))
		} else {
			this_.priv_RemoveUnusedCode_traverseNode_4667(getElemObj, scope)
			this_.priv_RemoveUnusedCode_traverseNode_4667(getElemKey, scope)
			this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
		}
	} else if lhs.IsGetProp() {
		var getPropLhs pkg_rhino.Node = lhs.GetFirstChild()
		_ = getPropLhs
		var isDotPrototypeLhs bool = RemoveUnusedCode_IsDotPrototype(getPropLhs)
		_ = isDotPrototypeLhs
		var isPrototypeMethodDef bool = (isDotPrototypeLhs && valueNode.IsFunction())
		_ = isPrototypeMethodDef
		if !(isPrototypeMethodDef) && this_.priv_RemoveUnusedCode_considerForAccessorSideEffects_4703(lhs, StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY()) {
			this_.priv_RemoveUnusedCode_traverseNode_4667(getPropLhs, scope)
			this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
		} else if getPropLhs.IsName() {
			var varInfo_3 *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(getPropLhs, scope)
			_ = varInfo_3
			var builder_3 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
			_ = builder_3
			this_.priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode, builder_3, scope)
			varInfo_3.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder_3.BuildNamedPropertyAssign_3(assignNode, lhs, varInfo_3)))
		} else if isDotPrototypeLhs {
			var objExpression pkg_rhino.Node = getPropLhs.GetFirstChild()
			_ = objExpression
			var builder_4 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsPrototypeDotPropertyReference(true)
			_ = builder_4
			this_.priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode, builder_4, scope)
			if objExpression.IsName() {
				var varInfo_4 *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(getPropLhs.GetFirstChild(), scope)
				_ = varInfo_4
				varInfo_4.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder_4.BuildNamedPropertyAssign_3(assignNode, lhs, varInfo_4)))
			} else {
				if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(objExpression) {
					this_.priv_RemoveUnusedCode_traverseNode_4667(objExpression, scope)
				} else {
					builder_4.AddContinuation(NewRemoveUnusedCode_Continuation(this_, objExpression, scope))
				}
				this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder_4.BuildAnonymousPrototypeNamedPropertyAssign(assignNode, lhs.GetString())))
			}
		} else if getPropLhs.IsThis() {
			var builder_5 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().SetIsThisDotPropertyReference(true)
			_ = builder_5
			this_.priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode, builder_5, scope)
			this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder_5.BuildNamedPropertyAssign_2(assignNode, lhs)))
		} else {
			this_.priv_RemoveUnusedCode_traverseNode_4667(lhs, scope)
			this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
		}
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(lhs, scope)
		this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
	}
}

// Java method: RemoveUnusedCode.traverseRemovableAssignValue(com.google.javascript.rhino.Node valueNode, com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseRemovableAssignValue_4688(valueNode pkg_rhino.Node, builder *RemoveUnusedCode_RemovableBuilder, scope *Scope) {
	_ = valueNode
	_ = builder
	_ = scope
	if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(valueNode) || NodeUtil_IsExpressionResultUsed(valueNode.GetParent()) {
		this_.priv_RemoveUnusedCode_traverseNode_4667(valueNode, scope)
	} else {
		builder.AddContinuation(NewRemoveUnusedCode_Continuation(this_, valueNode, scope))
	}
}

// Java method: RemoveUnusedCode.isNameDotPrototype(com.google.javascript.rhino.Node n)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_isNameDotPrototype_4689(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsGetProp() && n.GetFirstChild().IsName()) && (j2gort.CheckNotNullString(n.GetString()) == "prototype"))
}

// Java method: RemoveUnusedCode.traverseObjectPattern(com.google.javascript.rhino.Node pattern, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseObjectPattern_4690(pattern pkg_rhino.Node, scope *Scope) {
	_ = pattern
	_ = scope
	if !(pattern.IsObjectPattern()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(pattern))
	}
	{
		var elem pkg_rhino.Node = pattern.GetFirstChild()
		_ = elem
		for ; elem != 0; elem = elem.GetNext() {
			switch elem.TokenSwitchOrd_() {
			case int32(129) /* Token.COMPUTED_PROP */ :
				this_.priv_RemoveUnusedCode_traverseIndirectAssignment_4692(elem, elem.GetSecondChild(), scope)
			case int32(100) /* Token.STRING_KEY */ :
				if !(elem.IsQuotedStringKey()) {
					this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(elem.GetString())
				}
				this_.priv_RemoveUnusedCode_traverseIndirectAssignment_4692(elem, elem.GetOnlyChild(), scope)
			case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
				this_.priv_RemoveUnusedCode_traverseIndirectAssignment_4692(elem, elem.GetOnlyChild(), scope)
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected child of " + j2gort.StrVal(pattern.GetToken()) + ": " + j2gort.StrValString(elem.ToStringTree()))))))
			}
		}
	}
}

// Java method: RemoveUnusedCode.traverseIndirectAssignmentList(com.google.javascript.rhino.Node list, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseIndirectAssignmentList_4691(list pkg_rhino.Node, scope *Scope) {
	_ = list
	_ = scope
	if !(list.IsArrayPattern() || list.IsParamList()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(list))
	}
	{
		var elem pkg_rhino.Node = list.GetFirstChild()
		_ = elem
		for ; elem != 0; elem = elem.GetNext() {
			switch elem.TokenSwitchOrd_() {
			case int32(88) /* Token.EMPTY */ :
			case int32(102) /* Token.ARRAY_PATTERN */, int32(134) /* Token.DEFAULT_VALUE */, int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(32) /* Token.NAME */, int32(103) /* Token.OBJECT_PATTERN */ :
				this_.priv_RemoveUnusedCode_traverseIndirectAssignment_4692(elem, elem, scope)
			case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
				this_.priv_RemoveUnusedCode_traverseIndirectAssignment_4692(elem, elem.GetOnlyChild(), scope)
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected child of " + j2gort.StrVal(list.GetToken()) + ": " + j2gort.StrValString(elem.ToStringTree()))))))
			}
		}
	}
}

// Java method: RemoveUnusedCode.traverseIndirectAssignment(com.google.javascript.rhino.Node root, com.google.javascript.rhino.Node target, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseIndirectAssignment_4692(root pkg_rhino.Node, target pkg_rhino.Node, scope *Scope) {
	_ = root
	_ = target
	_ = scope
	var rootParent pkg_rhino.Node = root.GetParent()
	_ = rootParent
	if !(rootParent.IsDestructuringPattern() || rootParent.IsParamList()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(rootParent))
	}
	if target.IsDefaultValue() {
		target = target.GetFirstChild()
	}
	if target.IsGetProp() {
		this_.priv_RemoveUnusedCode_considerForAccessorSideEffects_4703(target, StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY())
	}
	if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(root) {
		this_.priv_RemoveUnusedCode_traverseNode_4667(root, scope)
	} else if target.IsName() {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(target, scope)
		_ = varInfo
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, root, scope))
		_ = builder
		varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder.BuildIndirectAssign(root, target)))
	} else if this_.priv_RemoveUnusedCode_isNameDotPrototype_4689(target) || RemoveUnusedCode_IsThisDotProperty(target) {
		var builder_2 *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, root, scope))
		_ = builder_2
		this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(builder_2.BuildIndirectAssign(root, target)))
	} else {
		this_.priv_RemoveUnusedCode_traverseNode_4667(root, scope)
	}
}

// Java method: RemoveUnusedCode.traverseChildren(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseChildren_4693(n pkg_rhino.Node, scope *Scope) {
	_ = n
	_ = scope
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.priv_RemoveUnusedCode_traverseNode_4667(c, scope)
		}
	}
}

// Java method: RemoveUnusedCode.traverseClass(com.google.javascript.rhino.Node classNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseClass_4694(classNode pkg_rhino.Node, scope *Scope) {
	_ = classNode
	_ = scope
	if !(classNode.IsClass()) {
		j2gort.CheckArgument(false)
	}
	if NodeUtil_IsClassDeclaration(classNode) {
		this_.priv_RemoveUnusedCode_traverseClassDeclaration_4695(classNode, scope)
	} else {
		this_.priv_RemoveUnusedCode_traverseClassExpression_4696(classNode, scope)
	}
}

// Java method: RemoveUnusedCode.traverseClassDeclaration(com.google.javascript.rhino.Node classNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseClassDeclaration_4695(classNode pkg_rhino.Node, scope *Scope) {
	_ = classNode
	_ = scope
	if !(classNode.IsClass()) {
		j2gort.CheckArgument(false)
	}
	var classNameNode pkg_rhino.Node = classNode.GetFirstChild()
	_ = classNameNode
	var baseClassExpression pkg_rhino.Node = classNameNode.GetNext()
	_ = baseClassExpression
	var classBodyNode pkg_rhino.Node = baseClassExpression.GetNext()
	_ = classBodyNode
	var classScope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(classNode, scope)
	_ = classScope
	var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(classNameNode, scope)
	_ = varInfo
	if classNode.GetParent().IsExport() {
		varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_7)
		this_.priv_RemoveUnusedCode_traverseNode_4667(baseClassExpression, scope)
		this_.priv_RemoveUnusedCode_traverseChildren_4693(classBodyNode, classScope)
	} else if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(baseClassExpression) {
		varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_8)
		this_.priv_RemoveUnusedCode_traverseNode_4667(baseClassExpression, scope)
		this_.priv_RemoveUnusedCode_traverseClassMembers_4697(classBodyNode, classScope)
	} else if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(classBodyNode) {
		varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_9)
		this_.priv_RemoveUnusedCode_traverseNode_4667(baseClassExpression, scope)
		this_.priv_RemoveUnusedCode_traverseClassMembers_4697(classBodyNode, classScope)
	} else {
		var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, baseClassExpression, classScope)).AddContinuation(NewRemoveUnusedCode_Continuation(this_, classBodyNode, classScope))
		_ = builder
		varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(builder.BuildClassDeclaration(classNode)))
	}
}

// Java method: RemoveUnusedCode.traverseClassExpression(com.google.javascript.rhino.Node classNode, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseClassExpression_4696(classNode pkg_rhino.Node, scope *Scope) {
	_ = classNode
	_ = scope
	if !(classNode.IsClass()) {
		j2gort.CheckArgument(false)
	}
	var classNameNode pkg_rhino.Node = classNode.GetFirstChild()
	_ = classNameNode
	var baseClassExpression pkg_rhino.Node = classNameNode.GetNext()
	_ = baseClassExpression
	var classBodyNode pkg_rhino.Node = baseClassExpression.GetNext()
	_ = classBodyNode
	var classScope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(classNode, scope)
	_ = classScope
	if classNameNode.IsName() {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(classNameNode, classScope)
		_ = varInfo
		varInfo.SetHasNonLocalOrNonLiteralValue()
		varInfo.AddRemovable(ToIface_RemoveUnusedCode_Removable(this_.priv_RemoveUnusedCode_acquireBuilder_4663().BuildNamedClassExpression(classNode)))
	}
	this_.priv_RemoveUnusedCode_traverseNode_4667(baseClassExpression, scope)
	this_.priv_RemoveUnusedCode_traverseClassMembers_4697(classBodyNode, classScope)
}

// Java method: RemoveUnusedCode.traverseClassMembers(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Scope scope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseClassMembers_4697(node pkg_rhino.Node, scope *Scope) {
	_ = node
	_ = scope
	if !(node.IsClassMembers()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(node))
	}
	if !(this_.removeUnusedPrototypeProperties) {
		this_.priv_RemoveUnusedCode_traverseChildren_4693(node, scope)
		return
	}
	{
		var member pkg_rhino.Node = node.GetFirstChild()
		_ = member
		for ; member != 0; member = member.GetNext() {
			switch member.TokenSwitchOrd_() {
			case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
				this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, member, scope)).BuildClassOrPrototypeNamedProperty(member)))
			case int32(108) /* Token.MEMBER_FIELD_DEF */ :
				if !(member.HasChildren()) || !(this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(member.GetFirstChild())) {
					this_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(ToIface_RemoveUnusedCode_Removable(this_.priv_RemoveUnusedCode_acquireBuilder_4663().AddContinuation(NewRemoveUnusedCode_Continuation(this_, member, scope)).BuildClassOrPrototypeNamedProperty(member)))
				}
				if member.HasChildren() {
					this_.priv_RemoveUnusedCode_traverseChildren_4693(member, scope)
				}
			case int32(129) /* Token.COMPUTED_PROP */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
				this_.priv_RemoveUnusedCode_traverseChildren_4693(member, scope)
			case int32(90) /* Token.BLOCK */ :
				this_.priv_RemoveUnusedCode_traverseChildren_4693(member, this_.scopeCreator.CreateScope_2_Node_Scope(member, scope))
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected child of CLASS_MEMBERS: "+j2gort.StrValString(member.ToStringTree()))))
			}
		}
	}
}

// Java method: RemoveUnusedCode.traverseFunction(com.google.javascript.rhino.Node function, com.google.javascript.jscomp.Scope parentScope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseFunction_4698(function pkg_rhino.Node, parentScope *Scope) {
	_ = function
	_ = parentScope
	if !(function.HasXChildren(int32(3))) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(function))
	}
	if !(function.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(function))
	}
	var paramlist pkg_rhino.Node = NodeUtil_GetFunctionParameters(function)
	_ = paramlist
	var body pkg_rhino.Node = function.GetLastChild()
	_ = body
	if !((body.GetNext() == 0) && body.IsBlock()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(body))
	}
	var fparamScope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(function, parentScope)
	_ = fparamScope
	var fbodyScope *Scope = this_.scopeCreator.CreateScope_2_Node_Scope(body, fparamScope)
	_ = fbodyScope
	var nameNode pkg_rhino.Node = function.GetFirstChild()
	_ = nameNode
	if !(j2gort.CheckNotNullString(nameNode.GetString()) == "") {
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(nameNode, fparamScope)
		_ = varInfo
		if NodeUtil_IsExpressionResultUsed(function) {
			varInfo.SetHasNonLocalOrNonLiteralValue()
		}
	}
	this_.priv_RemoveUnusedCode_traverseNode_4667(paramlist, fparamScope)
	this_.priv_RemoveUnusedCode_traverseChildren_4693(body, fbodyScope)
	this_.allFunctionParamScopes = append(this_.allFunctionParamScopes, fparamScope)
}

// Java method: RemoveUnusedCode.canRemoveParameters(com.google.javascript.rhino.Node parameterList)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_canRemoveParameters_4699(parameterList pkg_rhino.Node) bool {
	_ = parameterList
	if !(parameterList.IsParamList()) {
		j2gort.CheckState(false)
	}
	var function pkg_rhino.Node = parameterList.GetParent()
	_ = function
	return (this_.removeGlobals && !(NodeUtil_IsGetOrSetKey(function.GetParent())))
}

// Java method: RemoveUnusedCode.removeUnreferencedFunctionArgs(com.google.javascript.jscomp.Scope fparamScope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_removeUnreferencedFunctionArgs_4700(fparamScope *Scope) {
	_ = fparamScope
	if !(this_.removeGlobals) {
		return
	}
	var function pkg_rhino.Node = fparamScope.GetRootNode_0()
	_ = function
	if !(function.IsFunction()) {
		j2gort.CheckState(false)
	}
	if NodeUtil_IsGetOrSetKey(function.GetParent()) {
		return
	}
	var argList pkg_rhino.Node = NodeUtil_GetFunctionParameters(function)
	_ = argList
	this_.priv_RemoveUnusedCode_maybeRemoveUnusedTrailingParameters_4706(argList, fparamScope)
	this_.priv_RemoveUnusedCode_markUnusedParameters_4705(argList, fparamScope)
}

// Java method: RemoveUnusedCode.markPropertyNameAsPinned(java.lang.String propertyName)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propertyName string) {
	_ = propertyName
	if this_.pinnedPropertyNames.Add(propertyName) {
		var removed []RemoveUnusedCode_Removable = this_.removablesForPropertyNames.Remove(propertyName)
		_ = removed
		if removed != nil {
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = int32(len(removed))
				_ = n
				for ; i < n; i = i + 1 {
					removed[i].GetRemoveUnusedCode_RemovableBase().ApplyContinuations()
				}
			}
		}
	}
}

// Java method: RemoveUnusedCode.considerForIndependentRemoval(com.google.javascript.jscomp.RemoveUnusedCode.Removable removable)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_considerForIndependentRemoval_4702(removable RemoveUnusedCode_Removable) {
	_ = removable
	base_removable_ := removable.GetRemoveUnusedCode_RemovableBase()
	if removable.IsNamedProperty() {
		var propertyName string = removable.GetPropertyName()
		_ = propertyName
		if this_.pinnedPropertyNames.Contains(propertyName) || this_.codingConvention.IsExported_2(propertyName, false) {
			base_removable_.ApplyContinuations()
		} else if this_.priv_RemoveUnusedCode_isIndependentlyRemovable_4704(removable) {
			var list []RemoveUnusedCode_Removable = this_.removablesForPropertyNames.Get(propertyName)
			_ = list
			if list == nil {
				list = make([]RemoveUnusedCode_Removable, 0, int32(4))

			}
			list = append(list, removable)
			this_.removablesForPropertyNames.
				Put(propertyName, list)
		} else {
			base_removable_.ApplyContinuations()
			this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propertyName)
		}
	} else {
		base_removable_.ApplyContinuations()
	}
}

// Java method: RemoveUnusedCode.considerForAccessorSideEffects(com.google.javascript.rhino.Node getprop, com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind usage)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_considerForAccessorSideEffects_4703(getprop pkg_rhino.Node, usage *AccessorSummary_PropertyAccessKind) bool {
	_ = getprop
	_ = usage
	if !(NodeUtil_IsNormalOrOptChainGetProp(getprop)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(getprop))
	}
	var propName string = getprop.GetString()
	_ = propName
	var recorded *AccessorSummary_PropertyAccessKind = this_.compiler.GetAccessorSummary().GetKind(propName)
	_ = recorded
	if ((recorded.HasGetter() && usage.HasGetter()) && !(this_.assumeGettersArePure)) || (recorded.HasSetter() && usage.HasSetter()) {
		this_.priv_RemoveUnusedCode_markPropertyNameAsPinned_4701(propName)
		return true
	}
	return false
}

// Java method: RemoveUnusedCode.isIndependentlyRemovable(com.google.javascript.jscomp.RemoveUnusedCode.Removable removable)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_isIndependentlyRemovable_4704(removable RemoveUnusedCode_Removable) bool {
	_ = removable
	if removable.IsPrototypeProperty() {
		return this_.removeUnusedPrototypeProperties
	} else if removable.IsObjectDefinePropertiesDefinition() {
		return this_.removeUnusedObjectDefinePropertiesDefinitions
	} else if removable.IsThisDotPropertyReference() {
		return this_.removeUnusedThisProperties
	} else if removable.IsStaticProperty() {
		return this_.removeUnusedThisProperties
	} else {
		return false
	}
}

// Java method: RemoveUnusedCode.markUnusedParameters(com.google.javascript.rhino.Node paramList, com.google.javascript.jscomp.Scope fparamScope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_markUnusedParameters_4705(paramList pkg_rhino.Node, fparamScope *Scope) {
	_ = paramList
	_ = fparamScope
	if !(paramList.IsParamList()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(paramList))
	}
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			if param.IsUnusedParameter() {
				continue
			}
			var paramNameNode pkg_rhino.Node = RemoveUnusedCode_NameOfParam(param)
			_ = paramNameNode
			if paramNameNode == 0 {
				continue
			}
			var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_traverseNameNode_4672(paramNameNode, fparamScope)
			_ = varInfo
			if varInfo.IsRemovable() {
				param.SetUnusedParameter(true)
				this_.compiler.ReportChangeToEnclosingScope(paramList)
				if this_.removalLog.IsLogging() {
					this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg(paramNameNode, paramList)))
				}
			}
		}
	}
}

// Java method: RemoveUnusedCode.maybeRemoveUnusedTrailingParameters(com.google.javascript.rhino.Node argList, com.google.javascript.jscomp.Scope fparamScope)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_maybeRemoveUnusedTrailingParameters_4706(argList pkg_rhino.Node, fparamScope *Scope) {
	_ = argList
	_ = fparamScope
	if !(argList.IsParamList()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(argList))
	}
	var lastArg pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastArg
	for (func() pkg_rhino.Node {
		var asgn_11 pkg_rhino.Node = argList.GetLastChild()
		lastArg = asgn_11
		return asgn_11
	}()) != 0 {
		var argNode pkg_rhino.Node = lastArg
		_ = argNode
		if lastArg.IsDefaultValue() {
			argNode = lastArg.GetFirstChild()
			if this_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(lastArg.GetLastChild()) {
				break
			}
		}
		if argNode.IsRest() {
			argNode = argNode.GetFirstChild()
		}
		if argNode.IsDestructuringPattern() {
			if argNode.HasChildren() {
				break
			} else {
				NodeUtil_DeleteNode(lastArg, this_.compiler)
				if this_.removalLog.IsLogging() {
					this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg(argList)))
				}
				continue
			}
		}
		var varInfo *RemoveUnusedCode_VarInfo = this_.priv_RemoveUnusedCode_getVarInfo_4708(this_.priv_RemoveUnusedCode_getVarForNameNode_4676(argNode, fparamScope))
		_ = varInfo
		if varInfo.IsRemovable() {
			NodeUtil_DeleteNode(lastArg, this_.compiler)
			if this_.removalLog.IsLogging() {
				this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForNamedArg(argNode, argList)))
			}
		} else {
			break
		}
	}
}

// Java method: RemoveUnusedCode.traverseVar(com.google.javascript.jscomp.Var var)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_traverseVar_4707(var_ *Var) *RemoveUnusedCode_VarInfo {
	_ = var_
	j2gort.CheckNotNullPtr(var_)
	if this_.removeLocalVars && var_.IsArguments() {
		var functionScope *Scope = j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](var_.GetScope()).GetClosestHoistScope())
		_ = functionScope
		var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(functionScope.GetRootNode_0())
		_ = paramList
		{
			var param pkg_rhino.Node = paramList.GetFirstChild()
			_ = param
			for ; param != 0; param = param.GetNext() {
				var lValue pkg_rhino.Node = RemoveUnusedCode_NameOfParam(param)
				_ = lValue
				if lValue == 0 {
					continue
				}
				this_.priv_RemoveUnusedCode_getVarInfo_4708(this_.priv_RemoveUnusedCode_getVarForNameNode_4676(lValue, functionScope)).SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_10)
			}
		}
		return this_.canonicalUnremovableVarInfo
	} else {
		return this_.priv_RemoveUnusedCode_getVarInfo_4708(var_)
	}
}

// Java static method: RemoveUnusedCode.nameOfParam(com.google.javascript.rhino.Node param)
func RemoveUnusedCode_NameOfParam(param pkg_rhino.Node) pkg_rhino.Node {
	_ = param
	switch param.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return param
	case int32(134) /* Token.DEFAULT_VALUE */ :
		return RemoveUnusedCode_NameOfParam(param.GetFirstChild())
	case int32(125) /* Token.ITER_REST */ :
		return RemoveUnusedCode_NameOfParam(param.GetOnlyChild())
	case int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
		return pkg_rhino.Node(0)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected child of PARAM_LIST: "+j2gort.StrValString(param.ToStringTree()))))
	}
	panic("unreachable switch expression")
}

// Java method: RemoveUnusedCode.getVarInfo(com.google.javascript.jscomp.Var var)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_getVarInfo_4708(var_ *Var) *RemoveUnusedCode_VarInfo {
	_ = var_
	j2gort.CheckNotNullPtr(var_)
	var isGlobal bool = var_.IsGlobal()
	_ = isGlobal
	if var_.IsExtern() {
		if this_.unremovableLog.IsLogging() {
			this_.unremovableLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s: extern", []any{j2gort.BoxString(var_.GetName())}))
			}))
		}
		return this_.canonicalUnremovableVarInfo
	} else if this_.codingConvention.IsExported_2(var_.GetName(), !(isGlobal)) {
		if this_.unremovableLog.IsLogging() {
			this_.unremovableLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s: exported by convention", []any{j2gort.BoxString(var_.GetName())}))
			}))
		}
		return this_.canonicalUnremovableVarInfo
	} else if var_.IsArguments() {
		return this_.canonicalUnremovableVarInfo
	} else {
		var h int32 = var_.HashCode()
		_ = h
		var table []int32 = this_.varInfoTable
		_ = table
		var mask int32 = (int32(len(table)) - int32(1))
		_ = mask
		var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
		_ = slot
		for {
			var entry int32 = table[slot]
			_ = entry
			if entry == int32(0) {
				break
			}
			var idx int32 = (entry - int32(1))
			_ = idx
			var k *Var = this_.varInfoKeys[idx]
			_ = k
			if (k == var_) || ((k.GetScopedNameBase().scopeRoot == var_.GetScopedNameBase().scopeRoot) && func() bool {
				sr_12 := k.GetScopedNameBase().name
				se_13 := var_.GetScopedNameBase().name
				return j2gort.CheckNotNullString(sr_12) == se_13
			}()) {
				return this_.varInfoVals[idx]
			}
			slot = ((slot + int32(1)) & mask)
		}
		val_14 := this_.slab_RemoveUnusedCode_VarInfo.Alloc()

		val_14.outer_ = this_
		InitRemoveUnusedCode_VarInfo_1(val_14, var_.GetName())
		var varInfo *RemoveUnusedCode_VarInfo = val_14
		_ = varInfo
		if var_.GetParentNode().IsParamList() {
			varInfo.SetHasNonLocalOrNonLiteralValue()
		}
		if !(this_.removeGlobals) && isGlobal {
			varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_11)
		} else if !(this_.removeLocalVars) && !(isGlobal) {
			varInfo.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_12)
		}
		var newIdx int32 = this_.varInfoSize
		_ = newIdx
		if newIdx == int32(len(this_.varInfoKeys)) {
			var newCap int32 = (newIdx * int32(2))
			_ = newCap
			this_.varInfoKeys = func() []*Var {
				arrSrc_15 := this_.varInfoKeys
				arrLen_16 := newCap
				j2gort.CheckNotNullSlice(arrSrc_15)
				if arrLen_16 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_17 := j2gort.AcquirePooledSlice[*Var](int(arrLen_16))[:arrLen_16]
				copy(arrOut_17, arrSrc_15)
				j2gort.ReleasePooledSlice(arrSrc_15)
				return arrOut_17
			}()
			this_.varInfoVals = func() []*RemoveUnusedCode_VarInfo {
				arrSrc_18 := this_.varInfoVals
				arrLen_19 := newCap
				j2gort.CheckNotNullSlice(arrSrc_18)
				if arrLen_19 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_20 := j2gort.AcquirePooledSlice[*RemoveUnusedCode_VarInfo](int(arrLen_19))[:arrLen_19]
				copy(arrOut_20, arrSrc_18)
				j2gort.ReleasePooledSlice(arrSrc_18)
				return arrOut_20
			}()
		}
		{
			a_21 := this_.varInfoKeys
			i_22 := newIdx
			a_21[i_22] = var_
		}
		{
			a_23 := this_.varInfoVals
			i_24 := newIdx
			a_23[i_24] = varInfo
		}
		this_.varInfoSize = (newIdx + int32(1))
		{
			a_25 := table
			i_26 := slot
			a_25[i_26] = (newIdx + int32(1))
		}
		if (this_.varInfoSize * int32(4)) >= (int32(len(table)) * int32(3)) {
			this_.priv_RemoveUnusedCode_rehashVarInfoTable_4709((int32(len(table)) * int32(2)))
		}
		return varInfo
	}
}

// Java method: RemoveUnusedCode.rehashVarInfoTable(int newTableLen)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_rehashVarInfoTable_4709(newTableLen int32) {
	_ = newTableLen
	var newTable []int32 = j2gort.AcquirePooledSlice[int32](int(newTableLen))[:newTableLen]
	_ = newTable
	var mask int32 = (newTableLen - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.varInfoSize
		_ = n
		for ; i < n; i = i + 1 {
			var h int32 = this_.varInfoKeys[i].HashCode()
			_ = h
			var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
			_ = slot
			for newTable[slot] != int32(0) {
				slot = ((slot + int32(1)) & mask)
			}
			{
				a_27 := newTable
				i_28 := slot
				a_27[i_28] = (i + int32(1))
			}
		}
	}
	j2gort.ReleasePooledSlice(this_.varInfoTable)
	this_.varInfoTable = newTable
}

// Java method: RemoveUnusedCode.removeUnreferencedVarsAndPolyfills()
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_removeUnreferencedVarsAndPolyfills_4710() {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.varInfoSize
		_ = n
		for ; i < n; i = i + 1 {
			var var_ *Var = this_.varInfoKeys[i]
			_ = var_
			var varInfo *RemoveUnusedCode_VarInfo = this_.varInfoVals[i]
			_ = varInfo
			if !(varInfo.IsRemovable()) {
				continue
			}
			if this_.removalLog.IsLogging() {
				this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForVar(var_)))
			}
			varInfo.RemoveAllRemovables()
			var nameNode pkg_rhino.Node = var_.GetNameNode()
			_ = nameNode
			var toRemove pkg_rhino.Node = nameNode.GetParent()
			_ = toRemove
			if (toRemove == 0) || RemoveUnusedCode_AlreadyRemoved(toRemove) {
			} else if NodeUtil_IsFunctionExpression(toRemove) {
				if !(this_.preserveFunctionExpressionNames) {
					var fnNameNode pkg_rhino.Node = toRemove.GetFirstChild()
					_ = fnNameNode
					this_.compiler.ReportChangeToEnclosingScope(fnNameNode)
					fnNameNode.SetString("")
				}
			} else {
				if !(toRemove.IsParamList() || (toRemove.GetParent().IsParamList() && (toRemove.IsDefaultValue() || toRemove.IsRest()))) {
					j2gort.CheckStateFmt(false, "unremoved code: %s", j2gort.BoxNode(toRemove))
				}
			}
		}
	}
	var iter j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.polyfills.Values()).Iterator())
	_ = iter
	for iter.HasNext() {
		var polyfill RemoveUnusedCode_PolyfillInfo = j2gort.AsIface[RemoveUnusedCode_PolyfillInfo](iter.Next())
		_ = polyfill
		if polyfill.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable {
			if this_.removalLog.IsLogging() {
				this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForPolyfill(polyfill)))
			}
			polyfill.GetRemoveUnusedCode_PolyfillInfoBase().removable.Remove(this_.compiler)
			iter.Remove()
		}
	}
}

// Java static method: RemoveUnusedCode.isThisDotProperty(com.google.javascript.rhino.Node n)
func RemoveUnusedCode_IsThisDotProperty(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsNormalOrOptChainGetProp(n) && n.GetFirstChild().IsThis())
}

// Java static method: RemoveUnusedCode.isDotPrototypeDotProperty(com.google.javascript.rhino.Node n)
func RemoveUnusedCode_IsDotPrototypeDotProperty(n pkg_rhino.Node) bool {
	_ = n
	return (NodeUtil_IsNormalOrOptChainGetProp(n) && RemoveUnusedCode_IsDotPrototype(n.GetFirstChild()))
}

// Java static method: RemoveUnusedCode.maybeUnwrapQnameOrDefaultValueNode(com.google.javascript.rhino.Node targetNode, com.google.javascript.rhino.Node valueNode)
func RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(targetNode pkg_rhino.Node, valueNode pkg_rhino.Node) pkg_rhino.Node {
	_ = targetNode
	_ = valueNode
	if valueNode.IsOr() && targetNode.IsQualifiedName() {
		var lhsOfOr pkg_rhino.Node = j2gort.CheckNotNullNode(valueNode.GetFirstChild())
		_ = lhsOfOr
		if lhsOfOr.MatchesQualifiedName_1_Node(targetNode) {
			return valueNode.GetLastChild()
		}
	}
	return valueNode
}

// Java method: RemoveUnusedCode.mayHaveSideEffects(com.google.javascript.rhino.Node node)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_mayHaveSideEffects_4711(node pkg_rhino.Node) bool {
	_ = node
	var jsDocInfo *pkg_rhino.JSDocInfo
	_ = jsDocInfo
	if node.IsCall() {
		jsDocInfo = node.GetJSDocInfo()
	} else {
		jsDocInfo = nil
	}
	if (jsDocInfo != nil) && jsDocInfo.IsPureOrBreakMyCode() {
		return false
	}
	return this_.astAnalyzer.MayHaveSideEffects(node)
}

// Java static method: RemoveUnusedCode.alreadyRemoved(com.google.javascript.rhino.Node n)
func RemoveUnusedCode_AlreadyRemoved(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return true
	}
	if parent.IsRoot() {
		return false
	}
	return RemoveUnusedCode_AlreadyRemoved(parent)
}

// Java method: RemoveUnusedCode.createPolyfillInfo(com.google.javascript.rhino.Node call, com.google.javascript.jscomp.Scope scope, java.lang.String name)
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_createPolyfillInfo_4712(call pkg_rhino.Node, scope *Scope, name string) RemoveUnusedCode_PolyfillInfo {
	_ = call
	_ = scope
	_ = name
	if !(call.GetParent().IsExprResult()) {
		j2gort.CheckState(false)
	}
	var builder *RemoveUnusedCode_RemovableBuilder = this_.priv_RemoveUnusedCode_acquireBuilder_4663()
	_ = builder
	{
		var n pkg_rhino.Node = call.GetFirstChild().GetNext()
		_ = n
		for ; n != 0; n = n.GetNext() {
			builder.AddContinuation(NewRemoveUnusedCode_Continuation(this_, n, scope))
		}
	}
	var removable *RemoveUnusedCode_Polyfill = builder.BuildPolyfill(call.GetParent())
	_ = removable
	var lastDot int32 = j2gort.JStrLastIndexOf(name, ".")
	_ = lastDot
	if lastDot < int32(0) {
		return ToIface_RemoveUnusedCode_PolyfillInfo(allocPooled_RemoveUnusedCode_GlobalPolyfillInfo(this_, removable, name))
	}
	var owner string = j2gort.JStrSub(name, int32(0), lastDot)
	_ = owner
	var prop string = j2gort.JStrSubFrom(name, (lastDot + int32(1)))
	_ = prop
	if javalang.StringEndsWith(j2gort.CheckNotNullString(owner), ".prototype") {
		owner = j2gort.JStrSub(owner, int32(0), (j2gort.JStrLen(owner) - j2gort.JStrLen(".prototype")))
		return ToIface_RemoveUnusedCode_PolyfillInfo(allocPooled_RemoveUnusedCode_PrototypePropertyPolyfillInfo(this_, removable, prop, owner))
	}
	return ToIface_RemoveUnusedCode_PolyfillInfo(allocPooled_RemoveUnusedCode_StaticPropertyPolyfillInfo(this_, removable, prop, owner))
}

// Java method: RemoveUnusedCode.removeExpressionCompletely(com.google.javascript.rhino.Node expression)
func (this_ *RemoveUnusedCode) RemoveExpressionCompletely(expression pkg_rhino.Node) {
	_ = expression
	if !(!(NodeUtil_IsExpressionResultUsed(expression))) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expression))
	}
	var parent pkg_rhino.Node = expression.GetParent()
	_ = parent
	if parent.IsExprResult() {
		NodeUtil_DeleteNode(parent, this_.compiler)
	} else if parent.IsComma() {
		var otherChild pkg_rhino.Node = expression.GetNext()
		_ = otherChild
		if otherChild == 0 {
			otherChild = expression.GetPrevious()
		}
		this_.ReplaceNodeWith(parent, otherChild.Detach())
	} else {
		this_.ReplaceNodeWith(expression, pkg_rhino.IR_Number(float64(int32(0))).Srcref(expression))
	}
}

// Java method: RemoveUnusedCode.replaceNodeWith(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node replacement)
func (this_ *RemoveUnusedCode) ReplaceNodeWith(n pkg_rhino.Node, replacement pkg_rhino.Node) {
	_ = n
	_ = replacement
	this_.compiler.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(replacement)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, this_.compiler)
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.Builder
type RemoveUnusedCode_Builder struct {
	compiler                                      AbstractCompiler `gson:"compiler"`
	idHash_                                       atomic.Int32
	removeLocalVars                               bool `gson:"removeLocalVars"`
	removeGlobals                                 bool `gson:"removeGlobals"`
	preserveFunctionExpressionNames               bool `gson:"preserveFunctionExpressionNames"`
	removeUnusedPrototypeProperties               bool `gson:"removeUnusedPrototypeProperties"`
	removeUnusedThisProperties                    bool `gson:"removeUnusedThisProperties"`
	removeUnusedObjectDefinePropertiesDefinitions bool `gson:"removeUnusedObjectDefinePropertiesDefinitions"`
	removeUnusedPolyfills                         bool `gson:"removeUnusedPolyfills"`
	assumeGettersArePure                          bool `gson:"assumeGettersArePure"`
	disambiguateProperties                        bool `gson:"disambiguateProperties"`
}

func (this_ *RemoveUnusedCode_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnusedCode_Builder) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$Builder", "Builder"
}
func (this_ *RemoveUnusedCode_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnusedCode_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.Builder", "com.google.javascript.jscomp.RemoveUnusedCode$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_Builder](v_)
	})
}

func (this_ *RemoveUnusedCode_Builder) initFields_RemoveUnusedCode_Builder() {
	_ = this_

}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRemoveUnusedCode_Builder(compiler AbstractCompiler) *RemoveUnusedCode_Builder {
	this_ := &RemoveUnusedCode_Builder{}
	InitRemoveUnusedCode_Builder(this_, compiler)
	return this_
}

func InitRemoveUnusedCode_Builder(this_ *RemoveUnusedCode_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RemoveUnusedCode_Builder()
	this_.compiler = compiler
}

// Java method: Builder.removeLocalVars(boolean value)
func (this_ *RemoveUnusedCode_Builder) RemoveLocalVars(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.removeLocalVars = value
	return this_
}

// Java method: Builder.removeGlobals(boolean value)
func (this_ *RemoveUnusedCode_Builder) RemoveGlobals(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.removeGlobals = value
	return this_
}

// Java method: Builder.preserveFunctionExpressionNames(boolean value)
func (this_ *RemoveUnusedCode_Builder) PreserveFunctionExpressionNames(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.preserveFunctionExpressionNames = value
	return this_
}

// Java method: Builder.removeUnusedPrototypeProperties(boolean value)
func (this_ *RemoveUnusedCode_Builder) RemoveUnusedPrototypeProperties(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.removeUnusedPrototypeProperties = value
	return this_
}

// Java method: Builder.removeUnusedThisProperties(boolean value)
func (this_ *RemoveUnusedCode_Builder) RemoveUnusedThisProperties(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.removeUnusedThisProperties = value
	return this_
}

// Java method: Builder.removeUnusedObjectDefinePropertiesDefinitions(boolean value)
func (this_ *RemoveUnusedCode_Builder) RemoveUnusedObjectDefinePropertiesDefinitions(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.removeUnusedObjectDefinePropertiesDefinitions = value
	return this_
}

// Java method: Builder.removeUnusedPolyfills(boolean value)
func (this_ *RemoveUnusedCode_Builder) RemoveUnusedPolyfills(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.removeUnusedPolyfills = value
	return this_
}

// Java method: Builder.assumeGettersArePure(boolean value)
func (this_ *RemoveUnusedCode_Builder) AssumeGettersArePure(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.assumeGettersArePure = value
	return this_
}

// Java method: Builder.disambiguateProperties(boolean value)
func (this_ *RemoveUnusedCode_Builder) DisambiguateProperties(value bool) *RemoveUnusedCode_Builder {
	_ = value
	this_.disambiguateProperties = value
	return this_
}

// Java method: Builder.build()
func (this_ *RemoveUnusedCode_Builder) Build() *RemoveUnusedCode {
	return NewRemoveUnusedCode(this_)
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.RemovalLogRecord
type RemoveUnusedCode_RemovalLogRecord struct {
	kind                 string                                 `gson:"kind"`
	nameSupplier         j2gort.Ext_java_util_function_Supplier `gson:"nameSupplier"`
	functionNameSupplier j2gort.Ext_java_util_function_Supplier `gson:"functionNameSupplier"`
	idHash_              atomic.Int32
}

func (this_ *RemoveUnusedCode_RemovalLogRecord) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnusedCode_RemovalLogRecord) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_RemovalLogRecord) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$RemovalLogRecord", "RemovalLogRecord"
}
func (this_ *RemoveUnusedCode_RemovalLogRecord) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnusedCode_RemovalLogRecord) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.RemovalLogRecord", "com.google.javascript.jscomp.RemoveUnusedCode$RemovalLogRecord", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_RemovalLogRecord](v_)
	})
}
func (*RemoveUnusedCode_RemovalLogRecord) Implements_j2gort_Ext_java_util_function_Supplier() {}

func (this_ *RemoveUnusedCode_RemovalLogRecord) initFields_RemoveUnusedCode_RemovalLogRecord() {
	_ = this_
	this_.kind = j2gort.NullString
}

// Java constructor: RemovalLogRecord(java.lang.String kind, java.util.function.Supplier<java.lang.String> nameSupplier, java.util.function.Supplier<java.lang.String> functionNameSupplier)
func NewRemoveUnusedCode_RemovalLogRecord_3(kind string, nameSupplier j2gort.Ext_java_util_function_Supplier, functionNameSupplier j2gort.Ext_java_util_function_Supplier) *RemoveUnusedCode_RemovalLogRecord {
	this_ := &RemoveUnusedCode_RemovalLogRecord{}
	InitRemoveUnusedCode_RemovalLogRecord_3(this_, kind, nameSupplier, functionNameSupplier)
	return this_
}

func InitRemoveUnusedCode_RemovalLogRecord_3(this_ *RemoveUnusedCode_RemovalLogRecord, kind string, nameSupplier j2gort.Ext_java_util_function_Supplier, functionNameSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = kind
	_ = nameSupplier
	_ = functionNameSupplier
	this_.initFields_RemoveUnusedCode_RemovalLogRecord()
	this_.kind = j2gort.CheckNotNullString(kind)
	this_.nameSupplier = j2gort.CheckNotNull(nameSupplier)
	this_.functionNameSupplier = j2gort.CheckNotNull(functionNameSupplier)
}

// Java constructor: RemovalLogRecord(java.lang.String kind, java.util.function.Supplier<java.lang.String> nameSupplier)
func NewRemoveUnusedCode_RemovalLogRecord_2(kind string, nameSupplier j2gort.Ext_java_util_function_Supplier) *RemoveUnusedCode_RemovalLogRecord {
	this_ := &RemoveUnusedCode_RemovalLogRecord{}
	InitRemoveUnusedCode_RemovalLogRecord_2(this_, kind, nameSupplier)
	return this_
}

func InitRemoveUnusedCode_RemovalLogRecord_2(this_ *RemoveUnusedCode_RemovalLogRecord, kind string, nameSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = kind
	_ = nameSupplier
	InitRemoveUnusedCode_RemovalLogRecord_3(this_, kind, nameSupplier, lambdaSingleton_remove_unused_code_13)
}

// Java method: RemovalLogRecord.get()
func (this_ *RemoveUnusedCode_RemovalLogRecord) Get() any {
	return j2gort.BoxString(j2gort.Ext_java_lang_String_Join_2_any_slice_any("\t", []any{j2gort.BoxString(this_.kind), j2gort.BoxString(j2gort.CastString(this_.nameSupplier.Get())), j2gort.BoxString(j2gort.CastString(this_.functionNameSupplier.Get()))}))
}

// Java static method: RemovalLogRecord.forProperty(java.lang.String propName)
func RemoveUnusedCode_RemovalLogRecord_ForProperty(propName string) *RemoveUnusedCode_RemovalLogRecord {
	_ = propName
	return NewRemoveUnusedCode_RemovalLogRecord_2("prop", j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString(propName)
	}))
}

// Java static method: RemovalLogRecord.forVar(com.google.javascript.jscomp.Var var)
func RemoveUnusedCode_RemovalLogRecord_ForVar(var_ *Var) *RemoveUnusedCode_RemovalLogRecord {
	_ = var_
	return NewRemoveUnusedCode_RemovalLogRecord_2("var", func() j2gort.Ext_java_util_function_Supplier {
		bound_29 := var_
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxString(bound_29.GetName()) })
	}())
}

// Java static method: RemovalLogRecord.forPolyfill(com.google.javascript.jscomp.RemoveUnusedCode.PolyfillInfo polyfillInfo)
func RemoveUnusedCode_RemovalLogRecord_ForPolyfill(polyfillInfo RemoveUnusedCode_PolyfillInfo) *RemoveUnusedCode_RemovalLogRecord {
	_ = polyfillInfo
	return NewRemoveUnusedCode_RemovalLogRecord_2("poly", func() j2gort.Ext_java_util_function_Supplier {
		bound_30 := polyfillInfo
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxString(bound_30.GetName()) })
	}())
}

// Java static method: RemovalLogRecord.forNamedArg(com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node argList)
func RemoveUnusedCode_RemovalLogRecord_ForNamedArg(nameNode pkg_rhino.Node, argList pkg_rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	_ = nameNode
	_ = argList
	return NewRemoveUnusedCode_RemovalLogRecord_3("arg", func() j2gort.Ext_java_util_function_Supplier {
		bound_31 := nameNode
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxString(bound_31.GetString()) })
	}(), RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// Java static method: RemovalLogRecord.forDestructuringArg(com.google.javascript.rhino.Node argList)
func RemoveUnusedCode_RemovalLogRecord_ForDestructuringArg(argList pkg_rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	_ = argList
	return NewRemoveUnusedCode_RemovalLogRecord_3("arg", lambdaSingleton_remove_unused_code_14, RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// Java static method: RemovalLogRecord.forMarkingNamedArg(com.google.javascript.rhino.Node nameNode, com.google.javascript.rhino.Node argList)
func RemoveUnusedCode_RemovalLogRecord_ForMarkingNamedArg(nameNode pkg_rhino.Node, argList pkg_rhino.Node) *RemoveUnusedCode_RemovalLogRecord {
	_ = nameNode
	_ = argList
	return NewRemoveUnusedCode_RemovalLogRecord_3("argmark", func() j2gort.Ext_java_util_function_Supplier {
		bound_32 := nameNode
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxString(bound_32.GetString()) })
	}(), RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList))
}

// Java static method: RemovalLogRecord.getLoggableFunctionNameSupplier(com.google.javascript.rhino.Node argList)
func RemoveUnusedCode_RemovalLogRecord_GetLoggableFunctionNameSupplier(argList pkg_rhino.Node) j2gort.Ext_java_util_function_Supplier {
	_ = argList
	return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		var functionName string = NodeUtil_GetNearestFunctionName(j2gort.CheckNotNullNode(argList).GetParent())
		_ = functionName
		if functionName == j2gort.NullString {
			functionName = "<anonymous>"
		}
		return j2gort.BoxString(functionName)
	})
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.Continuation
type RemoveUnusedCode_Continuation struct {
	outer_ *RemoveUnusedCode
	scope  *Scope         `gson:"scope"`
	node   pkg_rhino.Node `gson:"node"`
}

func (this_ RemoveUnusedCode_Continuation) BindSelf_(s any) {
	_ = s
}
func (this_ RemoveUnusedCode_Continuation) IsJavaNil_() bool {
	return this_ == (RemoveUnusedCode_Continuation{})
}
func (RemoveUnusedCode_Continuation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$Continuation", "Continuation"
}
func (RemoveUnusedCode_Continuation) HasIdentityEquals_() {}
func (this_ RemoveUnusedCode_Continuation) IdentityHash_() int32 {
	if this_ == (RemoveUnusedCode_Continuation{}) {
		return 0
	}
	var h_ int32 = 1
	{
		hp_33 := this_.outer_
		h_ = 31 * h_
		if hp_33 != nil {
			h_ += j2gort.IdentityHash(hp_33)
		}
	}
	h_ = 31*h_ + int32(this_.node)
	{
		hp_34 := this_.scope
		h_ = 31 * h_
		if hp_34 != nil {
			h_ += j2gort.IdentityHash(hp_34)
		}
	}
	return h_
}
func (this_ RemoveUnusedCode_Continuation) FastValHash_() (int32, bool, bool) {
	if this_.IsJavaNil_() {
		return 0, false, true
	}
	return this_.IdentityHash_(), true, true
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.Continuation", "com.google.javascript.jscomp.RemoveUnusedCode$Continuation", func(v_ any) bool {
		return j2gort.InstanceOf[RemoveUnusedCode_Continuation](v_)
	})
}
func (this_ *RemoveUnusedCode_Continuation) SetOuter_(o *RemoveUnusedCode) { this_.outer_ = o }
func (this_ *RemoveUnusedCode_Continuation) GetOuter_() *RemoveUnusedCode  { return this_.outer_ }

func (this_ *RemoveUnusedCode_Continuation) initFields_RemoveUnusedCode_Continuation() {
	_ = this_
}

// Java constructor: Continuation(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Scope scope)
func NewRemoveUnusedCode_Continuation(outer_ *RemoveUnusedCode, node pkg_rhino.Node, scope *Scope) RemoveUnusedCode_Continuation {
	this_ := RemoveUnusedCode_Continuation{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_Continuation(&this_, node, scope)
	return this_
}

func InitRemoveUnusedCode_Continuation(this_ *RemoveUnusedCode_Continuation, node pkg_rhino.Node, scope *Scope) {
	_ = node
	_ = scope
	this_.initFields_RemoveUnusedCode_Continuation()
	this_.node = node
	this_.scope = scope
}

// Java method: Continuation.apply()
func (this_ RemoveUnusedCode_Continuation) Apply() {
	if this_ == (RemoveUnusedCode_Continuation{}) {
		j2gort.ThrowNPE()
	}
	if this_.node.IsFunction() {
		this_.outer_.priv_RemoveUnusedCode_traverseFunction_4698(this_.node, this_.scope)
	} else {
		this_.outer_.priv_RemoveUnusedCode_traverseNode_4667(this_.node, this_.scope)
	}
}

// Java polymorphic class interface: com.google.javascript.jscomp.RemoveUnusedCode.Removable
type RemoveUnusedCode_Removable interface {
	Implements_jscomp_RemoveUnusedCode_Removable()
	IsJavaNil_() bool
	GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase
	GetPropertyName() string
	RemoveInternal(AbstractCompiler)
	Remove(AbstractCompiler)
	ApplyContinuations()
	IsVariableAssignment() bool
	IsNamedProperty() bool
	IsNamedPropertyAssignment() bool
	IsAssignedValueLocal() bool
	GetLocalAssignedValue() pkg_rhino.Node
	IsPrototypeAssignment() bool
	IsPrototypeDotPropertyReference() bool
	IsClassOrPrototypeNamedProperty() bool
	IsPrototypeProperty() bool
	IsThisDotPropertyReference() bool
	IsObjectDefinePropertiesDefinition() bool
	IsStaticProperty() bool
	PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool
}

func ToIface_RemoveUnusedCode_Removable[T interface {
	RemoveUnusedCode_Removable
	comparable
}](v T) RemoveUnusedCode_Removable {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_RemoveUnusedCode_Removable[T RemoveUnusedCode_Removable](v T) RemoveUnusedCode_Removable {
	var vi RemoveUnusedCode_Removable = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.Removable
type RemoveUnusedCode_RemovableBase struct {
	self_                           RemoveUnusedCode_Removable
	outer_                          *RemoveUnusedCode
	continuation0                   RemoveUnusedCode_Continuation   `gson:"continuation0"`
	continuation1                   RemoveUnusedCode_Continuation   `gson:"continuation1"`
	restContinuations               []RemoveUnusedCode_Continuation `gson:"restContinuations"`
	propertyName                    string                          `gson:"propertyName"`
	targetNode                      pkg_rhino.Node                  `gson:"targetNode"`
	idHash_                         atomic.Int32
	isPrototypeDotPropertyReference bool `gson:"isPrototypeDotPropertyReference"`
	isThisDotPropertyReference      bool `gson:"isThisDotPropertyReference"`
	continuationsAreApplied         bool `gson:"continuationsAreApplied"`
	isRemoved                       bool `gson:"isRemoved"`
}

func (this_ *RemoveUnusedCode_RemovableBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(RemoveUnusedCode_Removable)
}
func (this_ *RemoveUnusedCode_RemovableBase) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_RemovableBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$Removable", "Removable"
}
func (this_ *RemoveUnusedCode_RemovableBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.Removable", "com.google.javascript.jscomp.RemoveUnusedCode$Removable", func(v_ any) bool {
		return j2gort.InstanceOf[RemoveUnusedCode_Removable](v_)
	})
}
func (b *RemoveUnusedCode_RemovableBase) GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase {
	return b
}
func (this_ *RemoveUnusedCode_RemovableBase) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_RemovableBase) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_RemovableBase) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_RemovableBase) initFields_RemoveUnusedCode_Removable() {
	_ = this_
	this_.propertyName = j2gort.NullString

}

func InitRemoveUnusedCode_Removable(this_ *RemoveUnusedCode_RemovableBase, targetNode pkg_rhino.Node, builder *RemoveUnusedCode_RemovableBuilder) {
	_ = targetNode
	_ = builder
	this_.initFields_RemoveUnusedCode_Removable()
	this_.continuation0 = builder.continuation0
	this_.continuation1 = builder.continuation1
	this_.restContinuations = builder.restContinuations
	this_.propertyName = builder.propertyName
	this_.isPrototypeDotPropertyReference = builder.isPrototypeDotPropertyReference
	this_.isThisDotPropertyReference = builder.isThisDotPropertyReference
	this_.targetNode = targetNode
}

// Java method: Removable.getPropertyName()
func (this_ *RemoveUnusedCode_RemovableBase) GetPropertyName() string {
	return j2gort.CheckNotNullString(this_.propertyName)
}

// Java method: Removable.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_RemovableBase) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	panic("abstract or native method invoked")
}

// Java method: Removable.remove(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_RemovableBase) Remove(compiler AbstractCompiler) {
	_ = compiler
	if !(this_.isRemoved) {
		this_.isRemoved = true
		this_.self_.RemoveInternal(compiler)
	}
}

// Java method: Removable.applyContinuations()
func (this_ *RemoveUnusedCode_RemovableBase) ApplyContinuations() {
	if !(this_.continuationsAreApplied) {
		this_.continuationsAreApplied = true
		if this_.continuation0 != (RemoveUnusedCode_Continuation{}) {
			this_.outer_.worklist.AddLast(this_.continuation0)
			this_.continuation0 = (RemoveUnusedCode_Continuation{})
			if this_.continuation1 != (RemoveUnusedCode_Continuation{}) {
				this_.outer_.worklist.AddLast(this_.continuation1)
				this_.continuation1 = (RemoveUnusedCode_Continuation{})
				if this_.restContinuations != nil {
					{
						var i int32 = int32(0)
						_ = i
						var n int32 = int32(len(this_.restContinuations))
						_ = n
						for ; i < n; i = i + 1 {
							this_.outer_.worklist.AddLast(this_.restContinuations[i])
						}
					}
					this_.restContinuations = nil
				}
			}
		}
	}
}

// Java method: Removable.isVariableAssignment()
func (this_ *RemoveUnusedCode_RemovableBase) IsVariableAssignment() bool {
	return false
}

// Java method: Removable.isNamedProperty()
func (this_ *RemoveUnusedCode_RemovableBase) IsNamedProperty() bool {
	return (this_.propertyName != j2gort.NullString)
}

// Java method: Removable.isNamedPropertyAssignment()
func (this_ *RemoveUnusedCode_RemovableBase) IsNamedPropertyAssignment() bool {
	return false
}

// Java method: Removable.isAssignedValueLocal()
func (this_ *RemoveUnusedCode_RemovableBase) IsAssignedValueLocal() bool {
	return false
}

// Java method: Removable.getLocalAssignedValue()
func (this_ *RemoveUnusedCode_RemovableBase) GetLocalAssignedValue() pkg_rhino.Node {
	return pkg_rhino.Node(0)
}

// Java method: Removable.isPrototypeAssignment()
func (this_ *RemoveUnusedCode_RemovableBase) IsPrototypeAssignment() bool {
	return (this_.self_.IsNamedPropertyAssignment() && (j2gort.CheckNotNullString(this_.propertyName) == "prototype"))
}

// Java method: Removable.isPrototypeDotPropertyReference()
func (this_ *RemoveUnusedCode_RemovableBase) IsPrototypeDotPropertyReference() bool {
	return this_.isPrototypeDotPropertyReference
}

// Java method: Removable.isClassOrPrototypeNamedProperty()
func (this_ *RemoveUnusedCode_RemovableBase) IsClassOrPrototypeNamedProperty() bool {
	return false
}

// Java method: Removable.isPrototypeProperty()
func (this_ *RemoveUnusedCode_RemovableBase) IsPrototypeProperty() bool {
	return (this_.IsPrototypeDotPropertyReference() || this_.self_.IsClassOrPrototypeNamedProperty())
}

// Java method: Removable.isThisDotPropertyReference()
func (this_ *RemoveUnusedCode_RemovableBase) IsThisDotPropertyReference() bool {
	return this_.isThisDotPropertyReference
}

// Java method: Removable.isObjectDefinePropertiesDefinition()
func (this_ *RemoveUnusedCode_RemovableBase) IsObjectDefinePropertiesDefinition() bool {
	return false
}

// Java method: Removable.isStaticProperty()
func (this_ *RemoveUnusedCode_RemovableBase) IsStaticProperty() bool {
	return false
}

// Java method: Removable.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (this_ *RemoveUnusedCode_RemovableBase) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	return false
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder
type RemoveUnusedCode_RemovableBuilder struct {
	outer_                          *RemoveUnusedCode
	continuation0                   RemoveUnusedCode_Continuation   `gson:"continuation0"`
	continuation1                   RemoveUnusedCode_Continuation   `gson:"continuation1"`
	restContinuations               []RemoveUnusedCode_Continuation `gson:"restContinuations"`
	propertyName                    string                          `gson:"propertyName"`
	idHash_                         atomic.Int32
	isPrototypeDotPropertyReference bool `gson:"isPrototypeDotPropertyReference"`
	isThisDotPropertyReference      bool `gson:"isThisDotPropertyReference"`
}

func (this_ *RemoveUnusedCode_RemovableBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnusedCode_RemovableBuilder) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_RemovableBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$RemovableBuilder", "RemovableBuilder"
}
func (this_ *RemoveUnusedCode_RemovableBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnusedCode_RemovableBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder", "com.google.javascript.jscomp.RemoveUnusedCode$RemovableBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_RemovableBuilder](v_)
	})
}
func (this_ *RemoveUnusedCode_RemovableBuilder) SetOuter_(o *RemoveUnusedCode) { this_.outer_ = o }
func (this_ *RemoveUnusedCode_RemovableBuilder) GetOuter_() *RemoveUnusedCode  { return this_.outer_ }

func (this_ *RemoveUnusedCode_RemovableBuilder) initFields_RemoveUnusedCode_RemovableBuilder() {
	_ = this_
	this_.propertyName = j2gort.NullString

}

// Java constructor: RemovableBuilder()
func NewRemoveUnusedCode_RemovableBuilder(outer_ *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	this_ := &RemoveUnusedCode_RemovableBuilder{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_RemovableBuilder(this_)
	return this_
}
func allocPooled_RemoveUnusedCode_RemovableBuilder(outer_ *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	this_ := outer_.slab_RemoveUnusedCode_RemovableBuilder.Alloc()
	this_.outer_ = outer_
	InitRemoveUnusedCode_RemovableBuilder(this_)
	return this_
}

func InitRemoveUnusedCode_RemovableBuilder(this_ *RemoveUnusedCode_RemovableBuilder) {
	this_.initFields_RemoveUnusedCode_RemovableBuilder()
}

// Java method: RemovableBuilder.finish(T removable)
func (this_ *RemoveUnusedCode_RemovableBuilder) priv_RemovableBuilder_finish_4743(removable RemoveUnusedCode_Removable) RemoveUnusedCode_Removable {
	_ = removable
	this_.continuation0 = (RemoveUnusedCode_Continuation{})
	this_.continuation1 = (RemoveUnusedCode_Continuation{})
	this_.restContinuations = nil
	this_.propertyName = j2gort.NullString
	this_.isPrototypeDotPropertyReference = false
	this_.isThisDotPropertyReference = false
	this_.outer_.freeBuilder = this_
	return removable
}

// Java method: RemovableBuilder.addContinuation(com.google.javascript.jscomp.RemoveUnusedCode.Continuation continuation)
func (this_ *RemoveUnusedCode_RemovableBuilder) AddContinuation(continuation RemoveUnusedCode_Continuation) *RemoveUnusedCode_RemovableBuilder {
	_ = continuation
	if this_.continuation0 == (RemoveUnusedCode_Continuation{}) {
		this_.continuation0 = continuation
	} else if this_.continuation1 == (RemoveUnusedCode_Continuation{}) {
		this_.continuation1 = continuation
	} else {
		if this_.restContinuations == nil {
			this_.restContinuations = make([]RemoveUnusedCode_Continuation, 0, int32(4))
		}
		this_.restContinuations = append(this_.restContinuations, continuation)
	}
	return this_
}

// Java method: RemovableBuilder.setIsPrototypeDotPropertyReference(boolean value)
func (this_ *RemoveUnusedCode_RemovableBuilder) SetIsPrototypeDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	_ = value
	this_.isPrototypeDotPropertyReference = value
	return this_
}

// Java method: RemovableBuilder.setIsThisDotPropertyReference(boolean value)
func (this_ *RemoveUnusedCode_RemovableBuilder) SetIsThisDotPropertyReference(value bool) *RemoveUnusedCode_RemovableBuilder {
	_ = value
	this_.isThisDotPropertyReference = value
	return this_
}

// Java method: RemovableBuilder.buildIndirectAssign(com.google.javascript.rhino.Node root, com.google.javascript.rhino.Node targetNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildIndirectAssign(root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	_ = root
	_ = targetNode
	return j2gort.CastPtr[RemoveUnusedCode_IndirectAssign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_IndirectAssign(this_.outer_, this_, root, targetNode))))
}

// Java method: RemovableBuilder.buildPolyfill(com.google.javascript.rhino.Node polyfillNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildPolyfill(polyfillNode pkg_rhino.Node) *RemoveUnusedCode_Polyfill {
	_ = polyfillNode
	return j2gort.CastPtr[RemoveUnusedCode_Polyfill](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_Polyfill(this_.outer_, this_, polyfillNode))))
}

// Java method: RemovableBuilder.buildClassDeclaration(com.google.javascript.rhino.Node classNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildClassDeclaration(classNode pkg_rhino.Node) *RemoveUnusedCode_ClassDeclaration {
	_ = classNode
	return j2gort.CastPtr[RemoveUnusedCode_ClassDeclaration](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_ClassDeclaration(this_.outer_, this_, classNode))))
}

// Java method: RemovableBuilder.buildNamedClassExpression(com.google.javascript.rhino.Node classNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildNamedClassExpression(classNode pkg_rhino.Node) *RemoveUnusedCode_NamedClassExpression {
	_ = classNode
	return j2gort.CastPtr[RemoveUnusedCode_NamedClassExpression](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_NamedClassExpression(this_.outer_, this_, classNode))))
}

// Java method: RemovableBuilder.buildClassOrPrototypeNamedProperty(com.google.javascript.rhino.Node propertyNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildClassOrPrototypeNamedProperty(propertyNode pkg_rhino.Node) *RemoveUnusedCode_ClassOrPrototypeNamedProperty {
	_ = propertyNode
	if !(((propertyNode.IsMemberFunctionDef() || propertyNode.IsMemberFieldDef()) || NodeUtil_IsGetOrSetKey(propertyNode)) || (propertyNode.IsStringKey() && !(propertyNode.IsQuotedStringKey()))) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(propertyNode))
	}
	this_.propertyName = propertyNode.GetString()
	return j2gort.CastPtr[RemoveUnusedCode_ClassOrPrototypeNamedProperty](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_ClassOrPrototypeNamedProperty(this_.outer_, this_, propertyNode))))
}

// Java method: RemovableBuilder.buildObjectDefinePropertiesDefinition(com.google.javascript.rhino.Node propertyNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildObjectDefinePropertiesDefinition(propertyNode pkg_rhino.Node) *RemoveUnusedCode_ObjectDefinePropertiesDefinition {
	_ = propertyNode
	this_.propertyName = propertyNode.GetString()
	return j2gort.CastPtr[RemoveUnusedCode_ObjectDefinePropertiesDefinition](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_ObjectDefinePropertiesDefinition(this_.outer_, this_, propertyNode))))
}

// Java method: RemovableBuilder.buildFunctionDeclaration(com.google.javascript.rhino.Node functionNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildFunctionDeclaration(functionNode pkg_rhino.Node) *RemoveUnusedCode_FunctionDeclaration {
	_ = functionNode
	return j2gort.CastPtr[RemoveUnusedCode_FunctionDeclaration](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_FunctionDeclaration(this_.outer_, this_, functionNode))))
}

// Java method: RemovableBuilder.buildNameDeclarationStatement(com.google.javascript.rhino.Node declarationStatement)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildNameDeclarationStatement(declarationStatement pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	_ = declarationStatement
	return j2gort.CastPtr[RemoveUnusedCode_NameDeclarationStatement](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_NameDeclarationStatement(this_.outer_, this_, declarationStatement))))
}

// Java method: RemovableBuilder.buildNamedPropertyAssign(com.google.javascript.rhino.Node assignNode, com.google.javascript.rhino.Node propertyNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssign_2(assignNode pkg_rhino.Node, propertyNode pkg_rhino.Node) *RemoveUnusedCode_Assign {
	_ = assignNode
	_ = propertyNode
	return this_.BuildNamedPropertyAssign_3(assignNode, propertyNode, nil)
}

// Java method: RemovableBuilder.buildNamedPropertyAssign(com.google.javascript.rhino.Node assignNode, com.google.javascript.rhino.Node propertyNode, com.google.javascript.jscomp.RemoveUnusedCode.@org.jspecify.annotations.Nullable VarInfo varInfo)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildNamedPropertyAssign_3(assignNode pkg_rhino.Node, propertyNode pkg_rhino.Node, varInfo *RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	_ = assignNode
	_ = propertyNode
	_ = varInfo
	this_.propertyName = propertyNode.GetString()
	return j2gort.CastPtr[RemoveUnusedCode_Assign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_Assign(this_.outer_, this_, assignNode, StaticGet_RemoveUnusedCode_Kind_NAMED_PROPERTY(), propertyNode, varInfo))))
}

// Java method: RemovableBuilder.buildComputedPropertyAssign(com.google.javascript.rhino.Node assignNode, com.google.javascript.rhino.Node propertyNode, com.google.javascript.jscomp.RemoveUnusedCode.VarInfo varInfo)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildComputedPropertyAssign(assignNode pkg_rhino.Node, propertyNode pkg_rhino.Node, varInfo *RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	_ = assignNode
	_ = propertyNode
	_ = varInfo
	return j2gort.CastPtr[RemoveUnusedCode_Assign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_Assign(this_.outer_, this_, assignNode, StaticGet_RemoveUnusedCode_Kind_COMPUTED_PROPERTY(), propertyNode, varInfo))))
}

// Java method: RemovableBuilder.buildVariableAssign(com.google.javascript.rhino.Node assignNode, com.google.javascript.jscomp.RemoveUnusedCode.VarInfo varInfo)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildVariableAssign(assignNode pkg_rhino.Node, varInfo *RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	_ = assignNode
	_ = varInfo
	return j2gort.CastPtr[RemoveUnusedCode_Assign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_Assign(this_.outer_, this_, assignNode, StaticGet_RemoveUnusedCode_Kind_VARIABLE(), pkg_rhino.Node(0), varInfo))))
}

// Java method: RemovableBuilder.buildClassSetupCall(com.google.javascript.rhino.Node callNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCall_1(callNode pkg_rhino.Node) *RemoveUnusedCode_ClassSetupCall {
	_ = callNode
	return this_.BuildClassSetupCall_2(callNode, false)
}

// Java method: RemovableBuilder.buildClassSetupCall(com.google.javascript.rhino.Node callNode, boolean classDefiningCall)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildClassSetupCall_2(callNode pkg_rhino.Node, classDefiningCall bool) *RemoveUnusedCode_ClassSetupCall {
	_ = callNode
	_ = classDefiningCall
	return j2gort.CastPtr[RemoveUnusedCode_ClassSetupCall](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_ClassSetupCall(this_.outer_, this_, callNode, classDefiningCall))))
}

// Java method: RemovableBuilder.buildVanillaForNameDeclaration(com.google.javascript.rhino.Node nameNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildVanillaForNameDeclaration(nameNode pkg_rhino.Node) *RemoveUnusedCode_VanillaForNameDeclaration {
	_ = nameNode
	return j2gort.CastPtr[RemoveUnusedCode_VanillaForNameDeclaration](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_VanillaForNameDeclaration(this_.outer_, this_, nameNode))))
}

// Java method: RemovableBuilder.buildAnonymousPrototypeNamedPropertyAssign(com.google.javascript.rhino.Node assignNode, java.lang.String propertyName)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildAnonymousPrototypeNamedPropertyAssign(assignNode pkg_rhino.Node, propertyName string) *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign {
	_ = assignNode
	_ = propertyName
	this_.propertyName = propertyName
	return j2gort.CastPtr[RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(this_.outer_, this_, assignNode))))
}

// Java method: RemovableBuilder.buildIncOrDepOp(com.google.javascript.rhino.Node incOrDecOp, com.google.javascript.rhino.Node propertyNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node toPreseve)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildIncOrDepOp(incOrDecOp pkg_rhino.Node, propertyNode pkg_rhino.Node, toPreseve pkg_rhino.Node) *RemoveUnusedCode_IncOrDecOp {
	_ = incOrDecOp
	_ = propertyNode
	_ = toPreseve
	this_.propertyName = propertyNode.GetString()
	return j2gort.CastPtr[RemoveUnusedCode_IncOrDecOp](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_IncOrDecOp(this_.outer_, this_, incOrDecOp, toPreseve))))
}

// Java method: RemovableBuilder.buildUnusedReadReference(com.google.javascript.rhino.Node referenceNode, com.google.javascript.rhino.Node propertyNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildUnusedReadReference(referenceNode pkg_rhino.Node, propertyNode pkg_rhino.Node) *RemoveUnusedCode_UnusedReadReference {
	_ = referenceNode
	_ = propertyNode
	this_.propertyName = propertyNode.GetString()
	return j2gort.CastPtr[RemoveUnusedCode_UnusedReadReference](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_UnusedReadReference(this_.outer_, this_, referenceNode))))
}

// Java method: RemovableBuilder.buildInstanceofName(com.google.javascript.rhino.Node instanceofNode)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildInstanceofName(instanceofNode pkg_rhino.Node) RemoveUnusedCode_Removable {
	_ = instanceofNode
	return ToIface_RemoveUnusedCode_Removable(j2gort.CastPtr[RemoveUnusedCode_InstanceofName](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_InstanceofName(this_.outer_, this_, instanceofNode)))))
}

// Java method: RemovableBuilder.buildWeakUsageCall(com.google.javascript.rhino.Node weakUsageCall)
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildWeakUsageCall(weakUsageCall pkg_rhino.Node) RemoveUnusedCode_Removable {
	_ = weakUsageCall
	return ToIface_RemoveUnusedCode_Removable(j2gort.CastPtr[RemoveUnusedCode_WeakUsageCall](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(allocPooled_RemoveUnusedCode_WeakUsageCall(this_.outer_, this_, weakUsageCall)))))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.UnusedReadReference
type RemoveUnusedCode_UnusedReadReference struct {
	RemoveUnusedCode_RemovableBase
	outer_        *RemoveUnusedCode
	referenceNode pkg_rhino.Node `gson:"referenceNode"`
}

func (this_ *RemoveUnusedCode_UnusedReadReference) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_UnusedReadReference) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_UnusedReadReference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$UnusedReadReference", "UnusedReadReference"
}
func (this_ *RemoveUnusedCode_UnusedReadReference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_UnusedReadReference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.UnusedReadReference", "com.google.javascript.jscomp.RemoveUnusedCode$UnusedReadReference", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_UnusedReadReference](v_)
	})
}
func (this_ *RemoveUnusedCode_UnusedReadReference) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_UnusedReadReference) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_UnusedReadReference) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_UnusedReadReference) initFields_RemoveUnusedCode_UnusedReadReference() {
	_ = this_
}

// Java constructor: UnusedReadReference(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node referenceNode)
func NewRemoveUnusedCode_UnusedReadReference(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, referenceNode pkg_rhino.Node) *RemoveUnusedCode_UnusedReadReference {
	this_ := &RemoveUnusedCode_UnusedReadReference{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_UnusedReadReference(this_, builder, referenceNode)
	return this_
}
func allocPooled_RemoveUnusedCode_UnusedReadReference(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, referenceNode pkg_rhino.Node) *RemoveUnusedCode_UnusedReadReference {
	this_ := outer_.slab_RemoveUnusedCode_UnusedReadReference.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_UnusedReadReference(this_, builder, referenceNode)
	return this_
}

func InitRemoveUnusedCode_UnusedReadReference(this_ *RemoveUnusedCode_UnusedReadReference, builder *RemoveUnusedCode_RemovableBuilder, referenceNode pkg_rhino.Node) {
	_ = builder
	_ = referenceNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_UnusedReadReference()
	if !(RemoveUnusedCode_IsThisDotProperty(referenceNode) || RemoveUnusedCode_IsDotPrototypeDotProperty(referenceNode)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(referenceNode))
	}
	this_.referenceNode = referenceNode
}

// Java method: UnusedReadReference.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_UnusedReadReference) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if !(RemoveUnusedCode_AlreadyRemoved(this_.referenceNode)) {
		if RemoveUnusedCode_IsThisDotProperty(this_.referenceNode) {
			this_.outer_.RemoveExpressionCompletely(this_.referenceNode)
		} else {
			if !(RemoveUnusedCode_IsDotPrototypeDotProperty(this_.referenceNode)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(this_.referenceNode))
			}
			var objExpression pkg_rhino.Node = this_.referenceNode.GetFirstFirstChild()
			_ = objExpression
			if this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(objExpression) {
				this_.outer_.ReplaceNodeWith(this_.referenceNode, objExpression.Detach())
			} else {
				this_.outer_.RemoveExpressionCompletely(this_.referenceNode)
			}
		}
	}
}

// Java method: UnusedReadReference.toString()
func (this_ *RemoveUnusedCode_UnusedReadReference) String() string {
	return ("UnusedReadReference:" + j2gort.StrVal(this_.referenceNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.InstanceofName
type RemoveUnusedCode_InstanceofName struct {
	RemoveUnusedCode_RemovableBase
	outer_         *RemoveUnusedCode
	instanceofNode pkg_rhino.Node `gson:"instanceofNode"`
}

func (this_ *RemoveUnusedCode_InstanceofName) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_InstanceofName) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_InstanceofName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$InstanceofName", "InstanceofName"
}
func (this_ *RemoveUnusedCode_InstanceofName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_InstanceofName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.InstanceofName", "com.google.javascript.jscomp.RemoveUnusedCode$InstanceofName", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_InstanceofName](v_)
	})
}
func (this_ *RemoveUnusedCode_InstanceofName) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_InstanceofName) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_InstanceofName) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_InstanceofName) initFields_RemoveUnusedCode_InstanceofName() {
	_ = this_
}

// Java constructor: InstanceofName(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node instanceofNode)
func NewRemoveUnusedCode_InstanceofName(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, instanceofNode pkg_rhino.Node) *RemoveUnusedCode_InstanceofName {
	this_ := &RemoveUnusedCode_InstanceofName{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_InstanceofName(this_, builder, instanceofNode)
	return this_
}
func allocPooled_RemoveUnusedCode_InstanceofName(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, instanceofNode pkg_rhino.Node) *RemoveUnusedCode_InstanceofName {
	this_ := outer_.slab_RemoveUnusedCode_InstanceofName.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_InstanceofName(this_, builder, instanceofNode)
	return this_
}

func InitRemoveUnusedCode_InstanceofName(this_ *RemoveUnusedCode_InstanceofName, builder *RemoveUnusedCode_RemovableBuilder, instanceofNode pkg_rhino.Node) {
	_ = builder
	_ = instanceofNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_InstanceofName()
	if !(instanceofNode.IsInstanceOf()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(instanceofNode))
	}
	this_.instanceofNode = instanceofNode
}

// Java method: InstanceofName.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_InstanceofName) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if !(RemoveUnusedCode_AlreadyRemoved(this_.instanceofNode)) {
		var lhs pkg_rhino.Node = this_.instanceofNode.GetFirstChild()
		_ = lhs
		var falseNode pkg_rhino.Node = pkg_rhino.IR_FalseNode().Srcref(this_.instanceofNode)
		_ = falseNode
		if this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(lhs) {
			this_.outer_.ReplaceNodeWith(this_.instanceofNode, pkg_rhino.IR_Comma(lhs.Detach(), falseNode).Srcref(this_.instanceofNode))
		} else {
			this_.outer_.ReplaceNodeWith(this_.instanceofNode, falseNode)
		}
	}
}

// Java method: InstanceofName.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (this_ *RemoveUnusedCode_InstanceofName) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	return true
}

// Java method: InstanceofName.toString()
func (this_ *RemoveUnusedCode_InstanceofName) String() string {
	return ("InstanceofName:" + j2gort.StrVal(this_.instanceofNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.WeakUsageCall
type RemoveUnusedCode_WeakUsageCall struct {
	RemoveUnusedCode_RemovableBase
	outer_       *RemoveUnusedCode
	weakCallNode pkg_rhino.Node `gson:"weakCallNode"`
}

func (this_ *RemoveUnusedCode_WeakUsageCall) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_WeakUsageCall) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_WeakUsageCall) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$WeakUsageCall", "WeakUsageCall"
}
func (this_ *RemoveUnusedCode_WeakUsageCall) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_WeakUsageCall) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.WeakUsageCall", "com.google.javascript.jscomp.RemoveUnusedCode$WeakUsageCall", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_WeakUsageCall](v_)
	})
}
func (this_ *RemoveUnusedCode_WeakUsageCall) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_WeakUsageCall) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_WeakUsageCall) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_WeakUsageCall) initFields_RemoveUnusedCode_WeakUsageCall() {
	_ = this_
}

// Java constructor: WeakUsageCall(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node weakCallNode)
func NewRemoveUnusedCode_WeakUsageCall(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, weakCallNode pkg_rhino.Node) *RemoveUnusedCode_WeakUsageCall {
	this_ := &RemoveUnusedCode_WeakUsageCall{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_WeakUsageCall(this_, builder, weakCallNode)
	return this_
}
func allocPooled_RemoveUnusedCode_WeakUsageCall(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, weakCallNode pkg_rhino.Node) *RemoveUnusedCode_WeakUsageCall {
	this_ := outer_.slab_RemoveUnusedCode_WeakUsageCall.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_WeakUsageCall(this_, builder, weakCallNode)
	return this_
}

func InitRemoveUnusedCode_WeakUsageCall(this_ *RemoveUnusedCode_WeakUsageCall, builder *RemoveUnusedCode_RemovableBuilder, weakCallNode pkg_rhino.Node) {
	_ = builder
	_ = weakCallNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_WeakUsageCall()
	if !(weakCallNode.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(weakCallNode))
	}
	this_.weakCallNode = weakCallNode
}

// Java method: WeakUsageCall.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_WeakUsageCall) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if !(RemoveUnusedCode_AlreadyRemoved(this_.weakCallNode)) {
		var undefinedNode pkg_rhino.Node = NodeUtil_NewUndefinedNode(this_.weakCallNode)
		_ = undefinedNode
		this_.outer_.ReplaceNodeWith(this_.weakCallNode, undefinedNode)
	}
}

// Java method: WeakUsageCall.toString()
func (this_ *RemoveUnusedCode_WeakUsageCall) String() string {
	return ("weakUsageCall:" + j2gort.StrVal(this_.weakCallNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.IncOrDecOp
type RemoveUnusedCode_IncOrDecOp struct {
	RemoveUnusedCode_RemovableBase
	outer_       *RemoveUnusedCode
	incOrDecNode pkg_rhino.Node `gson:"incOrDecNode"`
	toPreserve   pkg_rhino.Node `gson:"toPreserve"`
}

func (this_ *RemoveUnusedCode_IncOrDecOp) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_IncOrDecOp) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_IncOrDecOp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$IncOrDecOp", "IncOrDecOp"
}
func (this_ *RemoveUnusedCode_IncOrDecOp) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_IncOrDecOp) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.IncOrDecOp", "com.google.javascript.jscomp.RemoveUnusedCode$IncOrDecOp", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_IncOrDecOp](v_)
	})
}
func (this_ *RemoveUnusedCode_IncOrDecOp) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_IncOrDecOp) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_IncOrDecOp) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_IncOrDecOp) initFields_RemoveUnusedCode_IncOrDecOp() {
	_ = this_
}

// Java constructor: IncOrDecOp(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node incOrDecNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node toPreserve)
func NewRemoveUnusedCode_IncOrDecOp(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, incOrDecNode pkg_rhino.Node, toPreserve pkg_rhino.Node) *RemoveUnusedCode_IncOrDecOp {
	this_ := &RemoveUnusedCode_IncOrDecOp{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_IncOrDecOp(this_, builder, incOrDecNode, toPreserve)
	return this_
}
func allocPooled_RemoveUnusedCode_IncOrDecOp(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, incOrDecNode pkg_rhino.Node, toPreserve pkg_rhino.Node) *RemoveUnusedCode_IncOrDecOp {
	this_ := outer_.slab_RemoveUnusedCode_IncOrDecOp.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_IncOrDecOp(this_, builder, incOrDecNode, toPreserve)
	return this_
}

func InitRemoveUnusedCode_IncOrDecOp(this_ *RemoveUnusedCode_IncOrDecOp, builder *RemoveUnusedCode_RemovableBuilder, incOrDecNode pkg_rhino.Node, toPreserve pkg_rhino.Node) {
	_ = builder
	_ = incOrDecNode
	_ = toPreserve
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, incOrDecNode.GetOnlyChild(), builder)
	this_.initFields_RemoveUnusedCode_IncOrDecOp()
	if !(incOrDecNode.IsInc() || incOrDecNode.IsDec()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(incOrDecNode))
	}
	var arg pkg_rhino.Node = incOrDecNode.GetOnlyChild()
	_ = arg
	if !(RemoveUnusedCode_IsThisDotProperty(arg) || RemoveUnusedCode_IsDotPrototypeDotProperty(arg)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(arg))
	}
	this_.incOrDecNode = incOrDecNode
	this_.toPreserve = toPreserve
}

// Java method: IncOrDecOp.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_IncOrDecOp) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if RemoveUnusedCode_AlreadyRemoved(this_.incOrDecNode) {
		return
	}
	var arg pkg_rhino.Node = this_.incOrDecNode.GetOnlyChild()
	_ = arg
	if !(arg.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(arg))
	}
	if this_.toPreserve == 0 {
		this_.outer_.RemoveExpressionCompletely(this_.incOrDecNode)
	} else {
		this_.outer_.ReplaceNodeWith(this_.incOrDecNode, this_.toPreserve.Detach())
	}
}

// Java method: IncOrDecOp.toString()
func (this_ *RemoveUnusedCode_IncOrDecOp) String() string {
	return ("IncOrDecOp:" + j2gort.StrVal(this_.incOrDecNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.IndirectAssign
type RemoveUnusedCode_IndirectAssign struct {
	RemoveUnusedCode_RemovableBase
	outer_ *RemoveUnusedCode
	root   pkg_rhino.Node `gson:"root"`
}

func (this_ *RemoveUnusedCode_IndirectAssign) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_IndirectAssign) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_IndirectAssign) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$IndirectAssign", "IndirectAssign"
}
func (this_ *RemoveUnusedCode_IndirectAssign) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_IndirectAssign) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.IndirectAssign", "com.google.javascript.jscomp.RemoveUnusedCode$IndirectAssign", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_IndirectAssign](v_)
	})
}
func (this_ *RemoveUnusedCode_IndirectAssign) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_IndirectAssign) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_IndirectAssign) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_IndirectAssign) initFields_RemoveUnusedCode_IndirectAssign() {
	_ = this_
}

// Java constructor: IndirectAssign(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node root, com.google.javascript.rhino.Node targetNode)
func NewRemoveUnusedCode_IndirectAssign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	this_ := &RemoveUnusedCode_IndirectAssign{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_IndirectAssign(this_, builder, root, targetNode)
	return this_
}
func allocPooled_RemoveUnusedCode_IndirectAssign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	this_ := outer_.slab_RemoveUnusedCode_IndirectAssign.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_IndirectAssign(this_, builder, root, targetNode)
	return this_
}

func InitRemoveUnusedCode_IndirectAssign(this_ *RemoveUnusedCode_IndirectAssign, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) {
	_ = builder
	_ = root
	_ = targetNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, targetNode, builder)
	this_.initFields_RemoveUnusedCode_IndirectAssign()
	var rootParent pkg_rhino.Node = root.GetParent()
	_ = rootParent
	if !(rootParent.IsDestructuringPattern() || rootParent.IsParamList()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(rootParent))
	}
	if !(targetNode.IsName() || targetNode.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
	}
	this_.root = root
}

// Java method: IndirectAssign.isVariableAssignment()
func (this_ *RemoveUnusedCode_IndirectAssign) IsVariableAssignment() bool {
	return this_.GetRemoveUnusedCode_RemovableBase().targetNode.IsName()
}

// Java method: IndirectAssign.isThisDotPropertyReference()
func (this_ *RemoveUnusedCode_IndirectAssign) IsThisDotPropertyReference() bool {
	return RemoveUnusedCode_IsThisDotProperty(this_.GetRemoveUnusedCode_RemovableBase().targetNode)
}

// Java method: IndirectAssign.isNamedProperty()
func (this_ *RemoveUnusedCode_IndirectAssign) IsNamedProperty() bool {
	return this_.GetRemoveUnusedCode_RemovableBase().targetNode.IsGetProp()
}

// Java method: IndirectAssign.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (this_ *RemoveUnusedCode_IndirectAssign) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	if this_.GetRemoveUnusedCode_RemovableBase().targetNode.IsGetProp() {
		var getPropLhs pkg_rhino.Node = this_.GetRemoveUnusedCode_RemovableBase().targetNode.GetFirstChild()
		_ = getPropLhs
		return (getPropLhs.IsName() || this_.outer_.priv_RemoveUnusedCode_isNameDotPrototype_4689(getPropLhs))
	} else {
		return false
	}
}

// Java method: IndirectAssign.isNamedPropertyAssignment()
func (this_ *RemoveUnusedCode_IndirectAssign) IsNamedPropertyAssignment() bool {
	return this_.GetRemoveUnusedCode_RemovableBase().targetNode.IsGetProp()
}

// Java method: IndirectAssign.getPropertyName()
func (this_ *RemoveUnusedCode_IndirectAssign) GetPropertyName() string {
	if !(this_.GetRemoveUnusedCode_RemovableBase().targetNode.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(this_.GetRemoveUnusedCode_RemovableBase().targetNode))
	}
	return this_.GetRemoveUnusedCode_RemovableBase().targetNode.GetString()
}

// Java method: IndirectAssign.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_IndirectAssign) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if !(RemoveUnusedCode_AlreadyRemoved(this_.GetRemoveUnusedCode_RemovableBase().targetNode)) {
		this_.priv_IndirectAssign_removeRoot_4771()
	}
}

// Java method: IndirectAssign.removeRoot()
func (this_ *RemoveUnusedCode_IndirectAssign) priv_IndirectAssign_removeRoot_4771() {
	var rootParent pkg_rhino.Node = this_.root.GetParent()
	_ = rootParent
	switch rootParent.TokenSwitchOrd_() {
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.outer_.ReplaceNodeWith(this_.root, pkg_rhino.IR_Empty().Srcref(this_.root))
		{
			var maybeEmpty pkg_rhino.Node = rootParent.GetLastChild()
			_ = maybeEmpty
			for ; (maybeEmpty != 0) && maybeEmpty.IsEmpty(); maybeEmpty = rootParent.GetLastChild() {
				maybeEmpty.Detach()
			}
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(rootParent)
	case int32(49) /* Token.PARAM_LIST */ :
		if !(this_.root.IsDefaultValue()) {
			return
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(rootParent)
		var name pkg_rhino.Node = this_.root.GetFirstChild()
		_ = name
		if !(name.IsName()) {
			j2gort.CheckState(false)
		}
		if ((this_.root == rootParent.GetLastChild()) && this_.outer_.removeGlobals) && this_.outer_.priv_RemoveUnusedCode_canRemoveParameters_4699(rootParent) {
			this_.root.Detach()
		} else {
			this_.root.ReplaceWith(name.Detach())
		}
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.root, this_.outer_.compiler)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		NodeUtil_DeleteNode(this_.root, this_.outer_.compiler)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent of indirect assignment: "+j2gort.StrValString(rootParent.ToStringTree()))))
	}
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.Polyfill
type RemoveUnusedCode_Polyfill struct {
	RemoveUnusedCode_RemovableBase
	outer_       *RemoveUnusedCode
	polyfillNode pkg_rhino.Node `gson:"polyfillNode"`
	isPatch      bool           `gson:"isPatch"`
}

func (this_ *RemoveUnusedCode_Polyfill) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_Polyfill) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_Polyfill) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$Polyfill", "Polyfill"
}
func (this_ *RemoveUnusedCode_Polyfill) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_Polyfill) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.Polyfill", "com.google.javascript.jscomp.RemoveUnusedCode$Polyfill", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_Polyfill](v_)
	})
}
func (this_ *RemoveUnusedCode_Polyfill) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_Polyfill) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_Polyfill) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_Polyfill) initFields_RemoveUnusedCode_Polyfill() {
	_ = this_
}

// Java constructor: Polyfill(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node polyfillNode)
func NewRemoveUnusedCode_Polyfill(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, polyfillNode pkg_rhino.Node) *RemoveUnusedCode_Polyfill {
	this_ := &RemoveUnusedCode_Polyfill{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_Polyfill(this_, builder, polyfillNode)
	return this_
}
func allocPooled_RemoveUnusedCode_Polyfill(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, polyfillNode pkg_rhino.Node) *RemoveUnusedCode_Polyfill {
	this_ := outer_.slab_RemoveUnusedCode_Polyfill.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_Polyfill(this_, builder, polyfillNode)
	return this_
}

func InitRemoveUnusedCode_Polyfill(this_ *RemoveUnusedCode_Polyfill, builder *RemoveUnusedCode_RemovableBuilder, polyfillNode pkg_rhino.Node) {
	_ = builder
	_ = polyfillNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_Polyfill()
	this_.polyfillNode = polyfillNode
	this_.isPatch = (j2gort.JStrIndexOf(polyfillNode.GetFirstFirstChild().GetString(), "patch") >= 0)
}

// Java method: Polyfill.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_Polyfill) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	NodeUtil_DeleteNode(this_.polyfillNode, compiler)
}

// Java method: Polyfill.toString()
func (this_ *RemoveUnusedCode_Polyfill) String() string {
	return (j2gort.StrValString(func() string {
		if this_.isPatch {
			return "Patch:"
		}
		return "Polyfill:"
	}()) + j2gort.StrVal(this_.polyfillNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.ClassDeclaration
type RemoveUnusedCode_ClassDeclaration struct {
	RemoveUnusedCode_RemovableBase
	outer_               *RemoveUnusedCode
	classDeclarationNode pkg_rhino.Node `gson:"classDeclarationNode"`
}

func (this_ *RemoveUnusedCode_ClassDeclaration) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_ClassDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_ClassDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$ClassDeclaration", "ClassDeclaration"
}
func (this_ *RemoveUnusedCode_ClassDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_ClassDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.ClassDeclaration", "com.google.javascript.jscomp.RemoveUnusedCode$ClassDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_ClassDeclaration](v_)
	})
}
func (this_ *RemoveUnusedCode_ClassDeclaration) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_ClassDeclaration) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_ClassDeclaration) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_ClassDeclaration) initFields_RemoveUnusedCode_ClassDeclaration() {
	_ = this_
}

// Java constructor: ClassDeclaration(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node classDeclarationNode)
func NewRemoveUnusedCode_ClassDeclaration(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classDeclarationNode pkg_rhino.Node) *RemoveUnusedCode_ClassDeclaration {
	this_ := &RemoveUnusedCode_ClassDeclaration{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ClassDeclaration(this_, builder, classDeclarationNode)
	return this_
}
func allocPooled_RemoveUnusedCode_ClassDeclaration(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classDeclarationNode pkg_rhino.Node) *RemoveUnusedCode_ClassDeclaration {
	this_ := outer_.slab_RemoveUnusedCode_ClassDeclaration.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ClassDeclaration(this_, builder, classDeclarationNode)
	return this_
}

func InitRemoveUnusedCode_ClassDeclaration(this_ *RemoveUnusedCode_ClassDeclaration, builder *RemoveUnusedCode_RemovableBuilder, classDeclarationNode pkg_rhino.Node) {
	_ = builder
	_ = classDeclarationNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, classDeclarationNode.GetFirstChild(), builder)
	this_.initFields_RemoveUnusedCode_ClassDeclaration()
	this_.classDeclarationNode = classDeclarationNode
}

// Java method: ClassDeclaration.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_ClassDeclaration) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	NodeUtil_DeleteNode(this_.classDeclarationNode, compiler)
}

// Java method: ClassDeclaration.isVariableAssignment()
func (this_ *RemoveUnusedCode_ClassDeclaration) IsVariableAssignment() bool {
	return true
}

// Java method: ClassDeclaration.isAssignedValueLocal()
func (this_ *RemoveUnusedCode_ClassDeclaration) IsAssignedValueLocal() bool {
	return true
}

// Java method: ClassDeclaration.getLocalAssignedValue()
func (this_ *RemoveUnusedCode_ClassDeclaration) GetLocalAssignedValue() pkg_rhino.Node {
	return this_.classDeclarationNode
}

// Java method: ClassDeclaration.toString()
func (this_ *RemoveUnusedCode_ClassDeclaration) String() string {
	return ("ClassDeclaration:" + j2gort.StrVal(this_.classDeclarationNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.NamedClassExpression
type RemoveUnusedCode_NamedClassExpression struct {
	RemoveUnusedCode_RemovableBase
	outer_    *RemoveUnusedCode
	classNode pkg_rhino.Node `gson:"classNode"`
}

func (this_ *RemoveUnusedCode_NamedClassExpression) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_NamedClassExpression) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_NamedClassExpression) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$NamedClassExpression", "NamedClassExpression"
}
func (this_ *RemoveUnusedCode_NamedClassExpression) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_NamedClassExpression) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.NamedClassExpression", "com.google.javascript.jscomp.RemoveUnusedCode$NamedClassExpression", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_NamedClassExpression](v_)
	})
}
func (this_ *RemoveUnusedCode_NamedClassExpression) SetOuter_(o *RemoveUnusedCode) { this_.outer_ = o }
func (this_ *RemoveUnusedCode_NamedClassExpression) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_NamedClassExpression) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_NamedClassExpression) initFields_RemoveUnusedCode_NamedClassExpression() {
	_ = this_
}

// Java constructor: NamedClassExpression(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node classNode)
func NewRemoveUnusedCode_NamedClassExpression(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classNode pkg_rhino.Node) *RemoveUnusedCode_NamedClassExpression {
	this_ := &RemoveUnusedCode_NamedClassExpression{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_NamedClassExpression(this_, builder, classNode)
	return this_
}
func allocPooled_RemoveUnusedCode_NamedClassExpression(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, classNode pkg_rhino.Node) *RemoveUnusedCode_NamedClassExpression {
	this_ := outer_.slab_RemoveUnusedCode_NamedClassExpression.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_NamedClassExpression(this_, builder, classNode)
	return this_
}

func InitRemoveUnusedCode_NamedClassExpression(this_ *RemoveUnusedCode_NamedClassExpression, builder *RemoveUnusedCode_RemovableBuilder, classNode pkg_rhino.Node) {
	_ = builder
	_ = classNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, classNode.GetFirstChild(), builder)
	this_.initFields_RemoveUnusedCode_NamedClassExpression()
	this_.classNode = classNode
}

// Java method: NamedClassExpression.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_NamedClassExpression) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if !(RemoveUnusedCode_AlreadyRemoved(this_.classNode)) {
		var nameNode pkg_rhino.Node = this_.classNode.GetFirstChild()
		_ = nameNode
		if !(nameNode.IsEmpty()) {
			nameNode.ReplaceWith(pkg_rhino.IR_Empty().Srcref(nameNode))
			compiler.ReportChangeToEnclosingScope(this_.classNode)
		}
	}
}

// Java method: NamedClassExpression.toString()
func (this_ *RemoveUnusedCode_NamedClassExpression) String() string {
	return ("NamedClassExpression:" + j2gort.StrVal(this_.classNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.ClassOrPrototypeNamedProperty
type RemoveUnusedCode_ClassOrPrototypeNamedProperty struct {
	RemoveUnusedCode_RemovableBase
	outer_       *RemoveUnusedCode
	propertyNode pkg_rhino.Node `gson:"propertyNode"`
}

func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_ClassOrPrototypeNamedProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$ClassOrPrototypeNamedProperty", "ClassOrPrototypeNamedProperty"
}
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_ClassOrPrototypeNamedProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.ClassOrPrototypeNamedProperty", "com.google.javascript.jscomp.RemoveUnusedCode$ClassOrPrototypeNamedProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_ClassOrPrototypeNamedProperty](v_)
	})
}
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_ClassOrPrototypeNamedProperty) Implements_jscomp_RemoveUnusedCode_Removable() {
}

func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) initFields_RemoveUnusedCode_ClassOrPrototypeNamedProperty() {
	_ = this_
}

// Java constructor: ClassOrPrototypeNamedProperty(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node propertyNode)
func NewRemoveUnusedCode_ClassOrPrototypeNamedProperty(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode pkg_rhino.Node) *RemoveUnusedCode_ClassOrPrototypeNamedProperty {
	this_ := &RemoveUnusedCode_ClassOrPrototypeNamedProperty{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ClassOrPrototypeNamedProperty(this_, builder, propertyNode)
	return this_
}
func allocPooled_RemoveUnusedCode_ClassOrPrototypeNamedProperty(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode pkg_rhino.Node) *RemoveUnusedCode_ClassOrPrototypeNamedProperty {
	this_ := outer_.slab_RemoveUnusedCode_ClassOrPrototypeNamedProperty.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ClassOrPrototypeNamedProperty(this_, builder, propertyNode)
	return this_
}

func InitRemoveUnusedCode_ClassOrPrototypeNamedProperty(this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty, builder *RemoveUnusedCode_RemovableBuilder, propertyNode pkg_rhino.Node) {
	_ = builder
	_ = propertyNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_ClassOrPrototypeNamedProperty()
	this_.propertyNode = propertyNode
}

// Java method: ClassOrPrototypeNamedProperty.isStaticProperty()
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) IsStaticProperty() bool {
	return this_.propertyNode.IsStaticMember()
}

// Java method: ClassOrPrototypeNamedProperty.isClassOrPrototypeNamedProperty()
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) IsClassOrPrototypeNamedProperty() bool {
	return !(this_.IsStaticProperty())
}

// Java method: ClassOrPrototypeNamedProperty.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	NodeUtil_DeleteNode(this_.propertyNode, compiler)
}

// Java method: ClassOrPrototypeNamedProperty.toString()
func (this_ *RemoveUnusedCode_ClassOrPrototypeNamedProperty) String() string {
	return ("ClassOrPrototypeNamedProperty:" + j2gort.StrVal(this_.propertyNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.ObjectDefinePropertiesDefinition
type RemoveUnusedCode_ObjectDefinePropertiesDefinition struct {
	RemoveUnusedCode_RemovableBase
	outer_       *RemoveUnusedCode
	propertyNode pkg_rhino.Node `gson:"propertyNode"`
}

func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) IsJavaNil_() bool {
	return this_ == nil
}
func (*RemoveUnusedCode_ObjectDefinePropertiesDefinition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$ObjectDefinePropertiesDefinition", "ObjectDefinePropertiesDefinition"
}
func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_ObjectDefinePropertiesDefinition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.ObjectDefinePropertiesDefinition", "com.google.javascript.jscomp.RemoveUnusedCode$ObjectDefinePropertiesDefinition", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_ObjectDefinePropertiesDefinition](v_)
	})
}
func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_ObjectDefinePropertiesDefinition) Implements_jscomp_RemoveUnusedCode_Removable() {
}

func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) initFields_RemoveUnusedCode_ObjectDefinePropertiesDefinition() {
	_ = this_
}

// Java constructor: ObjectDefinePropertiesDefinition(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node propertyNode)
func NewRemoveUnusedCode_ObjectDefinePropertiesDefinition(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode pkg_rhino.Node) *RemoveUnusedCode_ObjectDefinePropertiesDefinition {
	this_ := &RemoveUnusedCode_ObjectDefinePropertiesDefinition{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ObjectDefinePropertiesDefinition(this_, builder, propertyNode)
	return this_
}
func allocPooled_RemoveUnusedCode_ObjectDefinePropertiesDefinition(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, propertyNode pkg_rhino.Node) *RemoveUnusedCode_ObjectDefinePropertiesDefinition {
	this_ := outer_.slab_RemoveUnusedCode_ObjectDefinePropertiesDefinition.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ObjectDefinePropertiesDefinition(this_, builder, propertyNode)
	return this_
}

func InitRemoveUnusedCode_ObjectDefinePropertiesDefinition(this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition, builder *RemoveUnusedCode_RemovableBuilder, propertyNode pkg_rhino.Node) {
	_ = builder
	_ = propertyNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_ObjectDefinePropertiesDefinition()
	this_.propertyNode = propertyNode
}

// Java method: ObjectDefinePropertiesDefinition.isObjectDefinePropertiesDefinition()
func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) IsObjectDefinePropertiesDefinition() bool {
	return true
}

// Java method: ObjectDefinePropertiesDefinition.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_ObjectDefinePropertiesDefinition) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	NodeUtil_DeleteNode(this_.propertyNode, compiler)
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.FunctionDeclaration
type RemoveUnusedCode_FunctionDeclaration struct {
	RemoveUnusedCode_RemovableBase
	outer_                  *RemoveUnusedCode
	functionDeclarationNode pkg_rhino.Node `gson:"functionDeclarationNode"`
}

func (this_ *RemoveUnusedCode_FunctionDeclaration) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_FunctionDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_FunctionDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$FunctionDeclaration", "FunctionDeclaration"
}
func (this_ *RemoveUnusedCode_FunctionDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_FunctionDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.FunctionDeclaration", "com.google.javascript.jscomp.RemoveUnusedCode$FunctionDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_FunctionDeclaration](v_)
	})
}
func (this_ *RemoveUnusedCode_FunctionDeclaration) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_FunctionDeclaration) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_FunctionDeclaration) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_FunctionDeclaration) initFields_RemoveUnusedCode_FunctionDeclaration() {
	_ = this_
}

// Java constructor: FunctionDeclaration(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node functionDeclarationNode)
func NewRemoveUnusedCode_FunctionDeclaration(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, functionDeclarationNode pkg_rhino.Node) *RemoveUnusedCode_FunctionDeclaration {
	this_ := &RemoveUnusedCode_FunctionDeclaration{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_FunctionDeclaration(this_, builder, functionDeclarationNode)
	return this_
}
func allocPooled_RemoveUnusedCode_FunctionDeclaration(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, functionDeclarationNode pkg_rhino.Node) *RemoveUnusedCode_FunctionDeclaration {
	this_ := outer_.slab_RemoveUnusedCode_FunctionDeclaration.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_FunctionDeclaration(this_, builder, functionDeclarationNode)
	return this_
}

func InitRemoveUnusedCode_FunctionDeclaration(this_ *RemoveUnusedCode_FunctionDeclaration, builder *RemoveUnusedCode_RemovableBuilder, functionDeclarationNode pkg_rhino.Node) {
	_ = builder
	_ = functionDeclarationNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, functionDeclarationNode.GetFirstChild(), builder)
	this_.initFields_RemoveUnusedCode_FunctionDeclaration()
	this_.functionDeclarationNode = functionDeclarationNode
}

// Java method: FunctionDeclaration.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_FunctionDeclaration) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	NodeUtil_DeleteNode(this_.functionDeclarationNode, compiler)
}

// Java method: FunctionDeclaration.isVariableAssignment()
func (this_ *RemoveUnusedCode_FunctionDeclaration) IsVariableAssignment() bool {
	return true
}

// Java method: FunctionDeclaration.isAssignedValueLocal()
func (this_ *RemoveUnusedCode_FunctionDeclaration) IsAssignedValueLocal() bool {
	return true
}

// Java method: FunctionDeclaration.getLocalAssignedValue()
func (this_ *RemoveUnusedCode_FunctionDeclaration) GetLocalAssignedValue() pkg_rhino.Node {
	return this_.functionDeclarationNode
}

// Java method: FunctionDeclaration.toString()
func (this_ *RemoveUnusedCode_FunctionDeclaration) String() string {
	return ("FunctionDeclaration:" + j2gort.StrVal(this_.functionDeclarationNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.NameDeclarationStatement
type RemoveUnusedCode_NameDeclarationStatement struct {
	RemoveUnusedCode_RemovableBase
	outer_               *RemoveUnusedCode
	declarationStatement pkg_rhino.Node `gson:"declarationStatement"`
}

func (this_ *RemoveUnusedCode_NameDeclarationStatement) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_NameDeclarationStatement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$NameDeclarationStatement", "NameDeclarationStatement"
}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_NameDeclarationStatement) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.NameDeclarationStatement", "com.google.javascript.jscomp.RemoveUnusedCode$NameDeclarationStatement", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_NameDeclarationStatement](v_)
	})
}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_NameDeclarationStatement) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_NameDeclarationStatement) initFields_RemoveUnusedCode_NameDeclarationStatement() {
	_ = this_
}

// Java constructor: NameDeclarationStatement(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node declarationStatement)
func NewRemoveUnusedCode_NameDeclarationStatement(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	this_ := &RemoveUnusedCode_NameDeclarationStatement{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_NameDeclarationStatement(this_, builder, declarationStatement)
	return this_
}
func allocPooled_RemoveUnusedCode_NameDeclarationStatement(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	this_ := outer_.slab_RemoveUnusedCode_NameDeclarationStatement.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_NameDeclarationStatement(this_, builder, declarationStatement)
	return this_
}

func InitRemoveUnusedCode_NameDeclarationStatement(this_ *RemoveUnusedCode_NameDeclarationStatement, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement pkg_rhino.Node) {
	_ = builder
	_ = declarationStatement
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, declarationStatement.GetOnlyChild(), builder)
	this_.initFields_RemoveUnusedCode_NameDeclarationStatement()
	if !(NodeUtil_IsNameDeclaration(declarationStatement)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(declarationStatement))
	}
	this_.declarationStatement = declarationStatement
}

// Java method: NameDeclarationStatement.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_NameDeclarationStatement) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	var nameNode pkg_rhino.Node = this_.declarationStatement.GetOnlyChild()
	_ = nameNode
	var valueNode pkg_rhino.Node = nameNode.GetFirstChild()
	_ = valueNode
	if (valueNode != 0) && this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(valueNode) {
		compiler.ReportChangeToEnclosingScope(this_.declarationStatement)
		valueNode.Detach()
		this_.declarationStatement.ReplaceWith(pkg_rhino.IR_ExprResult(valueNode).Srcref(valueNode))
	} else {
		NodeUtil_DeleteNode(this_.declarationStatement, compiler)
	}
}

// Java method: NameDeclarationStatement.isVariableAssignment()
func (this_ *RemoveUnusedCode_NameDeclarationStatement) IsVariableAssignment() bool {
	return true
}

// Java method: NameDeclarationStatement.isAssignedValueLocal()
func (this_ *RemoveUnusedCode_NameDeclarationStatement) IsAssignedValueLocal() bool {
	var nameNode pkg_rhino.Node = this_.declarationStatement.GetOnlyChild()
	_ = nameNode
	var initialValueNode pkg_rhino.Node = nameNode.GetFirstChild()
	_ = initialValueNode
	if initialValueNode == 0 {
		return true
	}
	var valueNode pkg_rhino.Node = RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(nameNode, initialValueNode)
	_ = valueNode
	return NodeUtil_EvaluatesToLocalValue(valueNode)
}

// Java method: NameDeclarationStatement.getLocalAssignedValue()
func (this_ *RemoveUnusedCode_NameDeclarationStatement) GetLocalAssignedValue() pkg_rhino.Node {
	var nameNode pkg_rhino.Node = this_.declarationStatement.GetOnlyChild()
	_ = nameNode
	var initialValueNode pkg_rhino.Node = nameNode.GetFirstChild()
	_ = initialValueNode
	if initialValueNode == 0 {
		return pkg_rhino.Node(0)
	}
	var valueNode pkg_rhino.Node = RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(nameNode, initialValueNode)
	_ = valueNode
	if NodeUtil_EvaluatesToLocalValue(valueNode) {
		return valueNode
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: NameDeclarationStatement.toString()
func (this_ *RemoveUnusedCode_NameDeclarationStatement) String() string {
	return ("NameDeclStmt:" + j2gort.StrVal(this_.declarationStatement))
}

// Java enum: com.google.javascript.jscomp.RemoveUnusedCode.Kind
type RemoveUnusedCode_Kind uint8

const (
	enum_RemoveUnusedCode_Kind_VARIABLE          RemoveUnusedCode_Kind = 1
	enum_RemoveUnusedCode_Kind_NAMED_PROPERTY    RemoveUnusedCode_Kind = 2
	enum_RemoveUnusedCode_Kind_COMPUTED_PROPERTY RemoveUnusedCode_Kind = 3
)

var enumBaseTable_RemoveUnusedCode_Kind = [...]j2gort.EnumBase{
	{Name_: "VARIABLE", Ordinal_: 0},
	{Name_: "NAMED_PROPERTY", Ordinal_: 1},
	{Name_: "COMPUTED_PROPERTY", Ordinal_: 2},
}
var enumValues_RemoveUnusedCode_Kind = []RemoveUnusedCode_Kind{1, 2, 3}

func StaticGet_RemoveUnusedCode_Kind_VARIABLE() RemoveUnusedCode_Kind {
	return enum_RemoveUnusedCode_Kind_VARIABLE
}
func StaticGet_RemoveUnusedCode_Kind_NAMED_PROPERTY() RemoveUnusedCode_Kind {
	return enum_RemoveUnusedCode_Kind_NAMED_PROPERTY
}
func StaticGet_RemoveUnusedCode_Kind_COMPUTED_PROPERTY() RemoveUnusedCode_Kind {
	return enum_RemoveUnusedCode_Kind_COMPUTED_PROPERTY
}

func RemoveUnusedCode_Kind_Values() []RemoveUnusedCode_Kind {
	return j2gort.CloneSlice(enumValues_RemoveUnusedCode_Kind)
}

func RemoveUnusedCode_Kind_ValueOf(name string) RemoveUnusedCode_Kind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_RemoveUnusedCode_Kind {
		if enumBaseTable_RemoveUnusedCode_Kind[i_].Name_ == name {
			return RemoveUnusedCode_Kind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.RemoveUnusedCode.Kind." + name))
}

func (this_ RemoveUnusedCode_Kind) IsJavaNil_() bool                              { return this_ == 0 }
func (RemoveUnusedCode_Kind) HasIdentityEquals_()                                 {}
func (RemoveUnusedCode_Kind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (RemoveUnusedCode_Kind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ RemoveUnusedCode_Kind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (RemoveUnusedCode_Kind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$Kind", "Kind"
}
func (this_ RemoveUnusedCode_Kind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_RemoveUnusedCode_Kind[this_-1]
}
func (this_ RemoveUnusedCode_Kind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ RemoveUnusedCode_Kind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ RemoveUnusedCode_Kind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RemoveUnusedCode_Kind[this_-1].Name_
}
func (this_ RemoveUnusedCode_Kind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_RemoveUnusedCode_Kind[this_-1].Name_
}
func (this_ RemoveUnusedCode_Kind) Equals(other_ any) bool {
	o_, ok_ := other_.(RemoveUnusedCode_Kind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ RemoveUnusedCode_Kind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_RemoveUnusedCode_Kind[this_-1].HashCode()
}
func (this_ RemoveUnusedCode_Kind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RemoveUnusedCode_Kind[this_-1].CompareTo(other_)
}
func (this_ RemoveUnusedCode_Kind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_RemoveUnusedCode_Kind[this_-1].Name_)
}
func init() {
	vals_ := []any{RemoveUnusedCode_Kind(1), RemoveUnusedCode_Kind(2), RemoveUnusedCode_Kind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.RemoveUnusedCode.Kind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.RemoveUnusedCode$Kind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.Kind", "com.google.javascript.jscomp.RemoveUnusedCode$Kind", func(v_ any) bool {
		return j2gort.InstanceOf[RemoveUnusedCode_Kind](v_)
	})
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.Assign
type RemoveUnusedCode_Assign struct {
	RemoveUnusedCode_RemovableBase
	outer_     *RemoveUnusedCode
	varInfo    *RemoveUnusedCode_VarInfo `gson:"varInfo"`
	assignNode pkg_rhino.Node            `gson:"assignNode"`
	kind       RemoveUnusedCode_Kind     `gson:"kind"`
}

func (this_ *RemoveUnusedCode_Assign) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_Assign) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_Assign) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$Assign", "Assign"
}
func (this_ *RemoveUnusedCode_Assign) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_Assign) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.Assign", "com.google.javascript.jscomp.RemoveUnusedCode$Assign", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_Assign](v_)
	})
}
func (this_ *RemoveUnusedCode_Assign) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_Assign) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_Assign) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_Assign) initFields_RemoveUnusedCode_Assign() {
	_ = this_
}

// Java constructor: Assign(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node assignNode, com.google.javascript.jscomp.RemoveUnusedCode.Kind kind, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node propertyNode, com.google.javascript.jscomp.RemoveUnusedCode.@org.jspecify.annotations.Nullable VarInfo varInfo)
func NewRemoveUnusedCode_Assign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node, kind RemoveUnusedCode_Kind, propertyNode pkg_rhino.Node, varInfo *RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	this_ := &RemoveUnusedCode_Assign{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_Assign(this_, builder, assignNode, kind, propertyNode, varInfo)
	return this_
}
func allocPooled_RemoveUnusedCode_Assign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node, kind RemoveUnusedCode_Kind, propertyNode pkg_rhino.Node, varInfo *RemoveUnusedCode_VarInfo) *RemoveUnusedCode_Assign {
	this_ := outer_.slab_RemoveUnusedCode_Assign.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_Assign(this_, builder, assignNode, kind, propertyNode, varInfo)
	return this_
}

func InitRemoveUnusedCode_Assign(this_ *RemoveUnusedCode_Assign, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node, kind RemoveUnusedCode_Kind, propertyNode pkg_rhino.Node, varInfo *RemoveUnusedCode_VarInfo) {
	_ = builder
	_ = assignNode
	_ = kind
	_ = propertyNode
	_ = varInfo
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, assignNode.GetFirstChild(), builder)
	this_.initFields_RemoveUnusedCode_Assign()
	if !(NodeUtil_IsAssignmentOp(assignNode)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(assignNode))
	}
	if kind == StaticGet_RemoveUnusedCode_Kind_VARIABLE() {
		if !(propertyNode == 0) {
			j2gort.CheckArgumentFmt(false, "got property node for simple variable assignment: %s", j2gort.BoxNode(propertyNode))
		}
		if !(varInfo != nil) {
			j2gort.CheckArgumentFmt(false, "missing VarInfo for variable assignment: %s", j2gort.BoxNode(propertyNode))
		}
	} else {
		if !(propertyNode != 0) {
			j2gort.CheckArgumentStr(false, "missing property node")
		}
		if kind == StaticGet_RemoveUnusedCode_Kind_NAMED_PROPERTY() {
			if !(propertyNode.IsGetProp()) {
				j2gort.CheckArgumentFmt(false, "property name is not a GETPROP: %s", j2gort.BoxNode(propertyNode))
			}
		}
	}
	this_.assignNode = assignNode
	this_.kind = kind
	this_.varInfo = varInfo
}

// Java method: Assign.isVariableAssignment()
func (this_ *RemoveUnusedCode_Assign) IsVariableAssignment() bool {
	return (this_.kind == StaticGet_RemoveUnusedCode_Kind_VARIABLE())
}

// Java method: Assign.isAssignedValueLocal()
func (this_ *RemoveUnusedCode_Assign) IsAssignedValueLocal() bool {
	return (this_.GetLocalAssignedValue() != 0)
}

// Java method: Assign.getLocalAssignedValue()
func (this_ *RemoveUnusedCode_Assign) GetLocalAssignedValue() pkg_rhino.Node {
	if NodeUtil_IsExpressionResultUsed(this_.assignNode) {
		return pkg_rhino.Node(0)
	} else {
		var valueNode pkg_rhino.Node = RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(this_.assignNode.GetFirstChild(), this_.assignNode.GetLastChild())
		_ = valueNode
		if NodeUtil_EvaluatesToLocalValue(valueNode) {
			return valueNode
		} else {
			return pkg_rhino.Node(0)
		}
	}
}

// Java method: Assign.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (this_ *RemoveUnusedCode_Assign) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	return (this_.IsNamedPropertyAssignment() || this_.IsComputedPropertyAssignment())
}

// Java method: Assign.isNamedPropertyAssignment()
func (this_ *RemoveUnusedCode_Assign) IsNamedPropertyAssignment() bool {
	return (this_.kind == StaticGet_RemoveUnusedCode_Kind_NAMED_PROPERTY())
}

// Java method: Assign.isComputedPropertyAssignment()
func (this_ *RemoveUnusedCode_Assign) IsComputedPropertyAssignment() bool {
	return (this_.kind == StaticGet_RemoveUnusedCode_Kind_COMPUTED_PROPERTY())
}

// Java method: Assign.isStaticProperty()
func (this_ *RemoveUnusedCode_Assign) IsStaticProperty() bool {
	if ((this_.kind == StaticGet_RemoveUnusedCode_Kind_NAMED_PROPERTY()) && (this_.varInfo != nil)) && this_.varInfo.HasFunctionOrClassLiteralValue() {
		return this_.GetRemoveUnusedCode_RemovableBase().targetNode.GetFirstChild().IsName()
	} else {
		return false
	}
}

// Java method: Assign.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_Assign) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if RemoveUnusedCode_AlreadyRemoved(this_.assignNode) {
		return
	}
	var parent pkg_rhino.Node = this_.assignNode.GetParent()
	_ = parent
	compiler.ReportChangeToEnclosingScope(parent)
	var lhs pkg_rhino.Node = this_.assignNode.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = this_.assignNode.GetSecondChild()
	_ = rhs
	var mustPreserveRhs bool = (this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(rhs) || NodeUtil_IsExpressionResultUsed(this_.assignNode))
	_ = mustPreserveRhs
	var mustPreserveGetElmExpr bool = (lhs.IsGetElem() && this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(lhs.GetLastChild()))
	_ = mustPreserveGetElmExpr
	if mustPreserveRhs && mustPreserveGetElmExpr {
		var replacement pkg_rhino.Node = pkg_rhino.IR_Comma(lhs.GetLastChild().Detach(), rhs.Detach()).Srcref(this_.assignNode)
		_ = replacement
		this_.outer_.ReplaceNodeWith(this_.assignNode, replacement)
	} else if mustPreserveGetElmExpr {
		this_.outer_.ReplaceNodeWith(this_.assignNode, lhs.GetLastChild().Detach())
	} else if mustPreserveRhs {
		this_.outer_.ReplaceNodeWith(this_.assignNode, rhs.Detach())
	} else {
		this_.outer_.RemoveExpressionCompletely(this_.assignNode)
	}
}

// Java method: Assign.toString()
func (this_ *RemoveUnusedCode_Assign) String() string {
	return ("Assign:" + j2gort.StrVal(this_.assignNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign
type RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign struct {
	RemoveUnusedCode_RemovableBase
	outer_     *RemoveUnusedCode
	assignNode pkg_rhino.Node `gson:"assignNode"`
}

func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) IsJavaNil_() bool {
	return this_ == nil
}
func (*RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$AnonymousPrototypeNamedPropertyAssign", "AnonymousPrototypeNamedPropertyAssign"
}
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.AnonymousPrototypeNamedPropertyAssign", "com.google.javascript.jscomp.RemoveUnusedCode$AnonymousPrototypeNamedPropertyAssign", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign](v_)
	})
}
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) Implements_jscomp_RemoveUnusedCode_Removable() {
}

func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) initFields_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign() {
	_ = this_
}

// Java constructor: AnonymousPrototypeNamedPropertyAssign(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node assignNode)
func NewRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node) *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign {
	this_ := &RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(this_, builder, assignNode)
	return this_
}
func allocPooled_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node) *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign {
	this_ := outer_.slab_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(this_, builder, assignNode)
	return this_
}

func InitRemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign(this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node) {
	_ = builder
	_ = assignNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, assignNode.GetFirstChild(), builder)
	this_.initFields_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign()
	j2gort.CheckNotNullString(builder.propertyName)
	if !(assignNode.IsAssign()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(assignNode))
	}
	this_.assignNode = assignNode
}

// Java method: AnonymousPrototypeNamedPropertyAssign.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	if RemoveUnusedCode_AlreadyRemoved(this_.assignNode) {
		return
	}
	var parent pkg_rhino.Node = this_.assignNode.GetParent()
	_ = parent
	compiler.ReportChangeToEnclosingScope(parent)
	var lhs pkg_rhino.Node = this_.assignNode.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = this_.assignNode.GetLastChild()
	_ = rhs
	if !(lhs.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
	}
	var objDotPrototype pkg_rhino.Node = lhs.GetFirstChild()
	_ = objDotPrototype
	if !(objDotPrototype.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(objDotPrototype))
	}
	var objExpression pkg_rhino.Node = objDotPrototype.GetFirstChild()
	_ = objExpression
	if !(j2gort.CheckNotNullString(objDotPrototype.GetString()) == "prototype") {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(objDotPrototype))
	}
	var mustPreserveRhs bool = (this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(rhs) || NodeUtil_IsExpressionResultUsed(this_.assignNode))
	_ = mustPreserveRhs
	var mustPreserveObjExpression bool = this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(objExpression)
	_ = mustPreserveObjExpression
	if mustPreserveRhs && mustPreserveObjExpression {
		var replacement pkg_rhino.Node = pkg_rhino.IR_Comma(objExpression.Detach(), rhs.Detach()).Srcref(this_.assignNode)
		_ = replacement
		this_.outer_.ReplaceNodeWith(this_.assignNode, replacement)
	} else if mustPreserveObjExpression {
		this_.outer_.ReplaceNodeWith(this_.assignNode, objExpression.Detach())
	} else if mustPreserveRhs {
		this_.outer_.ReplaceNodeWith(this_.assignNode, rhs.Detach())
	} else {
		this_.outer_.RemoveExpressionCompletely(this_.assignNode)
	}
}

// Java method: AnonymousPrototypeNamedPropertyAssign.isPrototypeProperty()
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) IsPrototypeProperty() bool {
	return true
}

// Java method: AnonymousPrototypeNamedPropertyAssign.toString()
func (this_ *RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign) String() string {
	return ("AnonymousPrototypeNamedPropertyAssign:" + j2gort.StrVal(this_.assignNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.ClassSetupCall
type RemoveUnusedCode_ClassSetupCall struct {
	RemoveUnusedCode_RemovableBase
	outer_            *RemoveUnusedCode
	callNode          pkg_rhino.Node `gson:"callNode"`
	classDefiningCall bool           `gson:"classDefiningCall"`
}

func (this_ *RemoveUnusedCode_ClassSetupCall) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_ClassSetupCall) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_ClassSetupCall) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$ClassSetupCall", "ClassSetupCall"
}
func (this_ *RemoveUnusedCode_ClassSetupCall) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_ClassSetupCall) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.ClassSetupCall", "com.google.javascript.jscomp.RemoveUnusedCode$ClassSetupCall", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_ClassSetupCall](v_)
	})
}
func (this_ *RemoveUnusedCode_ClassSetupCall) SetOuter_(o *RemoveUnusedCode)           { this_.outer_ = o }
func (this_ *RemoveUnusedCode_ClassSetupCall) GetOuter_() *RemoveUnusedCode            { return this_.outer_ }
func (*RemoveUnusedCode_ClassSetupCall) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_ClassSetupCall) initFields_RemoveUnusedCode_ClassSetupCall() {
	_ = this_
}

// Java constructor: ClassSetupCall(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node callNode, boolean classDefiningCall)
func NewRemoveUnusedCode_ClassSetupCall(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, callNode pkg_rhino.Node, classDefiningCall bool) *RemoveUnusedCode_ClassSetupCall {
	this_ := &RemoveUnusedCode_ClassSetupCall{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ClassSetupCall(this_, builder, callNode, classDefiningCall)
	return this_
}
func allocPooled_RemoveUnusedCode_ClassSetupCall(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, callNode pkg_rhino.Node, classDefiningCall bool) *RemoveUnusedCode_ClassSetupCall {
	this_ := outer_.slab_RemoveUnusedCode_ClassSetupCall.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_ClassSetupCall(this_, builder, callNode, classDefiningCall)
	return this_
}

func InitRemoveUnusedCode_ClassSetupCall(this_ *RemoveUnusedCode_ClassSetupCall, builder *RemoveUnusedCode_RemovableBuilder, callNode pkg_rhino.Node, classDefiningCall bool) {
	_ = builder
	_ = callNode
	_ = classDefiningCall
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, pkg_rhino.Node(0), builder)
	this_.initFields_RemoveUnusedCode_ClassSetupCall()
	this_.callNode = callNode
	this_.classDefiningCall = classDefiningCall
}

// Java method: ClassSetupCall.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_ClassSetupCall) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	var parent pkg_rhino.Node = this_.callNode.GetParent()
	_ = parent
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	this_.callNode.RemoveFirstChild()
	{
		var arg pkg_rhino.Node = this_.callNode.GetLastChild()
		_ = arg
		for ; arg != 0; arg = this_.callNode.GetLastChild() {
			arg.Detach()
			if !(this_.classDefiningCall) && this_.outer_.priv_RemoveUnusedCode_mayHaveSideEffects_4711(arg) {
				if replacement == 0 {
					replacement = arg
				} else {
					replacement = pkg_rhino.IR_Comma(arg, replacement).Srcref(this_.callNode)
				}
			} else {
				NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(arg, compiler)
			}
		}
	}
	if replacement != 0 {
		this_.outer_.ReplaceNodeWith(this_.callNode, replacement)
	} else if parent.IsExprResult() {
		NodeUtil_DeleteNode(parent, compiler)
	} else {
		if (parent.IsComma() || parent.IsAnd()) || parent.IsOr() {
			if parent.GetFirstChild() == this_.callNode {
				var rhs pkg_rhino.Node = j2gort.CheckNotNullNode(this_.callNode.GetNext())
				_ = rhs
				compiler.ReportChangeToEnclosingScope(parent)
				parent.ReplaceWith(rhs.Detach())
			} else {
				var lhs pkg_rhino.Node = parent.GetFirstChild()
				_ = lhs
				compiler.ReportChangeToEnclosingScope(parent)
				parent.ReplaceWith(lhs.Detach())
			}
		} else {
			compiler.ReportChangeToEnclosingScope(parent)
			this_.callNode.ReplaceWith(pkg_rhino.IR_Number(float64(int32(0))))
		}
	}
}

// Java method: ClassSetupCall.preventsRemovalOfVariableWithNonLocalValueOrPrototype()
func (this_ *RemoveUnusedCode_ClassSetupCall) PreventsRemovalOfVariableWithNonLocalValueOrPrototype() bool {
	return true
}

// Java method: ClassSetupCall.toString()
func (this_ *RemoveUnusedCode_ClassSetupCall) String() string {
	return ("ClassSetupCall:" + j2gort.StrVal(this_.callNode))
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.VarInfo
type RemoveUnusedCode_VarInfo struct {
	outer_                              *RemoveUnusedCode
	varName                             string                       `gson:"varName"`
	rem0                                RemoveUnusedCode_Removable   `gson:"rem0"`
	rem1                                RemoveUnusedCode_Removable   `gson:"rem1"`
	restRemovables                      []RemoveUnusedCode_Removable `gson:"restRemovables"`
	idHash_                             atomic.Int32
	canonicalUnremovable                bool `gson:"canonicalUnremovable"`
	isEntirelyRemovable                 bool `gson:"isEntirelyRemovable"`
	hasNonLocalOrNonLiteralValue        bool `gson:"hasNonLocalOrNonLiteralValue"`
	hasFunctionOrClassLiteralValue      bool `gson:"hasFunctionOrClassLiteralValue"`
	requiresLocalLiteralValueForRemoval bool `gson:"requiresLocalLiteralValueForRemoval"`
}

func (this_ *RemoveUnusedCode_VarInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveUnusedCode_VarInfo) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_VarInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$VarInfo", "VarInfo"
}
func (this_ *RemoveUnusedCode_VarInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveUnusedCode_VarInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.VarInfo", "com.google.javascript.jscomp.RemoveUnusedCode$VarInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_VarInfo](v_)
	})
}
func (this_ *RemoveUnusedCode_VarInfo) SetOuter_(o *RemoveUnusedCode) { this_.outer_ = o }
func (this_ *RemoveUnusedCode_VarInfo) GetOuter_() *RemoveUnusedCode  { return this_.outer_ }

func (this_ *RemoveUnusedCode_VarInfo) initFields_RemoveUnusedCode_VarInfo() {
	_ = this_
	this_.varName = j2gort.NullString

}

// Java constructor: VarInfo()
func NewRemoveUnusedCode_VarInfo(outer_ *RemoveUnusedCode) *RemoveUnusedCode_VarInfo {
	this_ := &RemoveUnusedCode_VarInfo{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_VarInfo(this_)
	return this_
}

func InitRemoveUnusedCode_VarInfo(this_ *RemoveUnusedCode_VarInfo) {
	this_.initFields_RemoveUnusedCode_VarInfo()
	this_.varName = ""
	this_.canonicalUnremovable = true
	this_.isEntirelyRemovable = false
}

// Java constructor: VarInfo(java.lang.String varName)
func NewRemoveUnusedCode_VarInfo_1(outer_ *RemoveUnusedCode, varName string) *RemoveUnusedCode_VarInfo {
	this_ := &RemoveUnusedCode_VarInfo{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_VarInfo_1(this_, varName)
	return this_
}

func InitRemoveUnusedCode_VarInfo_1(this_ *RemoveUnusedCode_VarInfo, varName string) {
	_ = varName
	this_.initFields_RemoveUnusedCode_VarInfo()
	this_.varName = varName
	this_.canonicalUnremovable = false
	this_.isEntirelyRemovable = true
}

// Java method: VarInfo.addRemovable(com.google.javascript.jscomp.RemoveUnusedCode.Removable removable)
func (this_ *RemoveUnusedCode_VarInfo) AddRemovable(removable RemoveUnusedCode_Removable) {
	_ = removable
	if this_.canonicalUnremovable {
		this_.outer_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(removable)
		return
	}
	if removable.IsVariableAssignment() {
		if removable.IsAssignedValueLocal() {
			var localValue pkg_rhino.Node = removable.GetLocalAssignedValue()
			_ = localValue
			if (localValue != 0) && (localValue.IsFunction() || localValue.IsClass()) {
				this_.hasFunctionOrClassLiteralValue = true
			}
		} else {
			this_.hasNonLocalOrNonLiteralValue = true
		}
	} else if removable.GetRemoveUnusedCode_RemovableBase().IsPrototypeAssignment() && !(removable.IsAssignedValueLocal()) {
		this_.hasNonLocalOrNonLiteralValue = true
	}
	if removable.PreventsRemovalOfVariableWithNonLocalValueOrPrototype() {
		this_.requiresLocalLiteralValueForRemoval = true
	}
	if this_.hasNonLocalOrNonLiteralValue && this_.requiresLocalLiteralValueForRemoval {
		this_.SetIsExplicitlyNotRemovable(lambdaSingleton_remove_unused_code_15)
	}
	if this_.isEntirelyRemovable {
		if this_.rem0 == nil || this_.rem0.IsJavaNil_() {
			this_.rem0 = removable
		} else if this_.rem1 == nil || this_.rem1.IsJavaNil_() {
			this_.rem1 = removable
		} else {
			if this_.restRemovables == nil {
				this_.restRemovables = make([]RemoveUnusedCode_Removable, 0, int32(4))
			}
			this_.restRemovables = append(this_.restRemovables, removable)
		}
	} else {
		this_.outer_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(removable)
	}
}

// Java method: VarInfo.isRemovable()
func (this_ *RemoveUnusedCode_VarInfo) IsRemovable() bool {
	return this_.isEntirelyRemovable
}

// Java method: VarInfo.setIsExplicitlyNotRemovable(java.util.function.Supplier<java.lang.String> reasonSupplier)
func (this_ *RemoveUnusedCode_VarInfo) SetIsExplicitlyNotRemovable(reasonSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = reasonSupplier
	if this_.isEntirelyRemovable {
		this_.isEntirelyRemovable = false
		if this_.outer_.unremovableLog.IsLogging() {
			this_.outer_.unremovableLog.Log_1_string(j2gort.Ext_java_lang_String_Format_2("%s: %s", []any{j2gort.BoxString(this_.varName), j2gort.BoxString(j2gort.CastString(reasonSupplier.Get()))}))
		}
		if this_.rem0 != nil && !this_.rem0.IsJavaNil_() {
			var r0 RemoveUnusedCode_Removable = this_.rem0
			_ = r0
			var r1 RemoveUnusedCode_Removable = this_.rem1
			_ = r1
			var rest []RemoveUnusedCode_Removable = this_.restRemovables
			_ = rest
			this_.rem0 = nil
			this_.rem1 = nil
			this_.restRemovables = nil
			this_.outer_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(r0)
			if r1 != nil && !r1.IsJavaNil_() {
				this_.outer_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(r1)
				if rest != nil {
					{
						var i int32 = int32(0)
						_ = i
						var n int32 = int32(len(rest))
						_ = n
						for ; i < n; i = i + 1 {
							this_.outer_.priv_RemoveUnusedCode_considerForIndependentRemoval_4702(rest[i])
						}
					}
				}
			}
		}
	}
}

// Java method: VarInfo.setHasNonLocalOrNonLiteralValue()
func (this_ *RemoveUnusedCode_VarInfo) SetHasNonLocalOrNonLiteralValue() {
	if !(this_.canonicalUnremovable) {
		this_.hasNonLocalOrNonLiteralValue = true
	}
}

// Java method: VarInfo.hasFunctionOrClassLiteralValue()
func (this_ *RemoveUnusedCode_VarInfo) HasFunctionOrClassLiteralValue() bool {
	return this_.hasFunctionOrClassLiteralValue
}

// Java method: VarInfo.removeAllRemovables()
func (this_ *RemoveUnusedCode_VarInfo) RemoveAllRemovables() {
	if this_.canonicalUnremovable {
		return
	}
	if !(this_.isEntirelyRemovable) {
		j2gort.CheckState(false)
	}
	if this_.rem0 != nil && !this_.rem0.IsJavaNil_() {
		this_.rem0.GetRemoveUnusedCode_RemovableBase().Remove(this_.outer_.compiler)
		this_.rem0 = nil
		if this_.rem1 != nil && !this_.rem1.IsJavaNil_() {
			this_.rem1.GetRemoveUnusedCode_RemovableBase().Remove(this_.outer_.compiler)
			this_.rem1 = nil
			if this_.restRemovables != nil {
				{
					var i int32 = int32(0)
					_ = i
					var n int32 = int32(len(this_.restRemovables))
					_ = n
					for ; i < n; i = i + 1 {
						this_.restRemovables[i].GetRemoveUnusedCode_RemovableBase().Remove(this_.outer_.compiler)
					}
				}
				this_.restRemovables = nil
			}
		}
	}
}

// Java polymorphic class interface: com.google.javascript.jscomp.RemoveUnusedCode.PolyfillInfo
type RemoveUnusedCode_PolyfillInfo interface {
	Implements_jscomp_RemoveUnusedCode_PolyfillInfo()
	IsJavaNil_() bool
	GetRemoveUnusedCode_PolyfillInfoBase() *RemoveUnusedCode_PolyfillInfoBase
	ConsiderPossibleReference(pkg_rhino.Node)
	GetName() string
	ConsiderPossibleReferenceInternal(pkg_rhino.Node)
}

func ToIface_RemoveUnusedCode_PolyfillInfo[T interface {
	RemoveUnusedCode_PolyfillInfo
	comparable
}](v T) RemoveUnusedCode_PolyfillInfo {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_RemoveUnusedCode_PolyfillInfo[T RemoveUnusedCode_PolyfillInfo](v T) RemoveUnusedCode_PolyfillInfo {
	var vi RemoveUnusedCode_PolyfillInfo = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.PolyfillInfo
type RemoveUnusedCode_PolyfillInfoBase struct {
	self_       RemoveUnusedCode_PolyfillInfo
	outer_      *RemoveUnusedCode
	removable   *RemoveUnusedCode_Polyfill `gson:"removable"`
	key         string                     `gson:"key"`
	idHash_     atomic.Int32
	isRemovable bool `gson:"isRemovable"`
}

func (this_ *RemoveUnusedCode_PolyfillInfoBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(RemoveUnusedCode_PolyfillInfo)
}
func (this_ *RemoveUnusedCode_PolyfillInfoBase) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_PolyfillInfoBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$PolyfillInfo", "PolyfillInfo"
}
func (this_ *RemoveUnusedCode_PolyfillInfoBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.PolyfillInfo", "com.google.javascript.jscomp.RemoveUnusedCode$PolyfillInfo", func(v_ any) bool {
		return j2gort.InstanceOf[RemoveUnusedCode_PolyfillInfo](v_)
	})
}
func (b *RemoveUnusedCode_PolyfillInfoBase) GetRemoveUnusedCode_PolyfillInfoBase() *RemoveUnusedCode_PolyfillInfoBase {
	return b
}
func (this_ *RemoveUnusedCode_PolyfillInfoBase) SetOuter_(o *RemoveUnusedCode)              { this_.outer_ = o }
func (this_ *RemoveUnusedCode_PolyfillInfoBase) GetOuter_() *RemoveUnusedCode               { return this_.outer_ }
func (*RemoveUnusedCode_PolyfillInfoBase) Implements_jscomp_RemoveUnusedCode_PolyfillInfo() {}

func (this_ *RemoveUnusedCode_PolyfillInfoBase) initFields_RemoveUnusedCode_PolyfillInfo() {
	_ = this_
	this_.key = j2gort.NullString
	this_.isRemovable = true
}

func InitRemoveUnusedCode_PolyfillInfo(this_ *RemoveUnusedCode_PolyfillInfoBase, removable *RemoveUnusedCode_Polyfill, key string) {
	_ = removable
	_ = key
	this_.initFields_RemoveUnusedCode_PolyfillInfo()
	this_.removable = removable
	this_.key = key
}

// Java method: PolyfillInfo.considerPossibleReference(com.google.javascript.rhino.Node n)
func (this_ *RemoveUnusedCode_PolyfillInfoBase) ConsiderPossibleReference(n pkg_rhino.Node) {
	_ = n
	if this_.isRemovable && (this_.removable.isPatch || !(this_.outer_.guardedUsages.Contains(n))) {
		this_.self_.ConsiderPossibleReferenceInternal(n)
		if !(this_.isRemovable) {
			this_.removable.ApplyContinuations()
		}
	}
}

// Java method: PolyfillInfo.getName()
func (this_ *RemoveUnusedCode_PolyfillInfoBase) GetName() string {
	return this_.key
}

// Java method: PolyfillInfo.considerPossibleReferenceInternal(com.google.javascript.rhino.Node n)
func (this_ *RemoveUnusedCode_PolyfillInfoBase) ConsiderPossibleReferenceInternal(n pkg_rhino.Node) {
	_ = n
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.GlobalPolyfillInfo
type RemoveUnusedCode_GlobalPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	outer_ *RemoveUnusedCode
}

func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_PolyfillInfoBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_GlobalPolyfillInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$GlobalPolyfillInfo", "GlobalPolyfillInfo"
}
func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_PolyfillInfoBase.IdentityHash_()
}
func (*RemoveUnusedCode_GlobalPolyfillInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.GlobalPolyfillInfo", "com.google.javascript.jscomp.RemoveUnusedCode$GlobalPolyfillInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_GlobalPolyfillInfo](v_)
	})
}
func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) SetOuter_(o *RemoveUnusedCode)              { this_.outer_ = o }
func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) GetOuter_() *RemoveUnusedCode               { return this_.outer_ }
func (*RemoveUnusedCode_GlobalPolyfillInfo) Implements_jscomp_RemoveUnusedCode_PolyfillInfo() {}

func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) initFields_RemoveUnusedCode_GlobalPolyfillInfo() {
	_ = this_
}

// Java constructor: GlobalPolyfillInfo(com.google.javascript.jscomp.RemoveUnusedCode.Polyfill removable, java.lang.String name)
func NewRemoveUnusedCode_GlobalPolyfillInfo(outer_ *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, name string) *RemoveUnusedCode_GlobalPolyfillInfo {
	this_ := &RemoveUnusedCode_GlobalPolyfillInfo{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_GlobalPolyfillInfo(this_, removable, name)
	return this_
}
func allocPooled_RemoveUnusedCode_GlobalPolyfillInfo(outer_ *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, name string) *RemoveUnusedCode_GlobalPolyfillInfo {
	this_ := outer_.slab_RemoveUnusedCode_GlobalPolyfillInfo.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_GlobalPolyfillInfo(this_, removable, name)
	return this_
}

func InitRemoveUnusedCode_GlobalPolyfillInfo(this_ *RemoveUnusedCode_GlobalPolyfillInfo, removable *RemoveUnusedCode_Polyfill, name string) {
	_ = removable
	_ = name
	this_.RemoveUnusedCode_PolyfillInfoBase.outer_ = this_.outer_
	InitRemoveUnusedCode_PolyfillInfo(&this_.RemoveUnusedCode_PolyfillInfoBase, removable, name)
	this_.initFields_RemoveUnusedCode_GlobalPolyfillInfo()
}

// Java method: GlobalPolyfillInfo.considerPossibleReferenceInternal(com.google.javascript.rhino.Node possiblyReferencingNode)
func (this_ *RemoveUnusedCode_GlobalPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode pkg_rhino.Node) {
	_ = possiblyReferencingNode
	if possiblyReferencingNode.IsName() {
		this_.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable = false
	} else if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		this_.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable = false
	}
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.StaticPropertyPolyfillInfo
type RemoveUnusedCode_StaticPropertyPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	outer_            *RemoveUnusedCode
	polyfillOwnerName string `gson:"polyfillOwnerName"`
}

func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_PolyfillInfoBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_StaticPropertyPolyfillInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$StaticPropertyPolyfillInfo", "StaticPropertyPolyfillInfo"
}
func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_PolyfillInfoBase.IdentityHash_()
}
func (*RemoveUnusedCode_StaticPropertyPolyfillInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.StaticPropertyPolyfillInfo", "com.google.javascript.jscomp.RemoveUnusedCode$StaticPropertyPolyfillInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_StaticPropertyPolyfillInfo](v_)
	})
}
func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_StaticPropertyPolyfillInfo) Implements_jscomp_RemoveUnusedCode_PolyfillInfo() {
}

func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) initFields_RemoveUnusedCode_StaticPropertyPolyfillInfo() {
	_ = this_
	this_.polyfillOwnerName = j2gort.NullString
}

// Java constructor: StaticPropertyPolyfillInfo(com.google.javascript.jscomp.RemoveUnusedCode.Polyfill removable, java.lang.String key, java.lang.String ownerName)
func NewRemoveUnusedCode_StaticPropertyPolyfillInfo(outer_ *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string, ownerName string) *RemoveUnusedCode_StaticPropertyPolyfillInfo {
	this_ := &RemoveUnusedCode_StaticPropertyPolyfillInfo{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_StaticPropertyPolyfillInfo(this_, removable, key, ownerName)
	return this_
}
func allocPooled_RemoveUnusedCode_StaticPropertyPolyfillInfo(outer_ *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string, ownerName string) *RemoveUnusedCode_StaticPropertyPolyfillInfo {
	this_ := outer_.slab_RemoveUnusedCode_StaticPropertyPolyfillInfo.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_StaticPropertyPolyfillInfo(this_, removable, key, ownerName)
	return this_
}

func InitRemoveUnusedCode_StaticPropertyPolyfillInfo(this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo, removable *RemoveUnusedCode_Polyfill, key string, ownerName string) {
	_ = removable
	_ = key
	_ = ownerName
	this_.RemoveUnusedCode_PolyfillInfoBase.outer_ = this_.outer_
	InitRemoveUnusedCode_PolyfillInfo(&this_.RemoveUnusedCode_PolyfillInfoBase, removable, key)
	this_.initFields_RemoveUnusedCode_StaticPropertyPolyfillInfo()
	this_.polyfillOwnerName = j2gort.CheckNotNullString(ownerName)
}

// Java method: StaticPropertyPolyfillInfo.getName()
func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) GetName() string {
	return (j2gort.StrValString(this_.polyfillOwnerName) + "." + j2gort.StrValString(this_.GetRemoveUnusedCode_PolyfillInfoBase().key))
}

// Java method: StaticPropertyPolyfillInfo.considerPossibleReferenceInternal(com.google.javascript.rhino.Node possiblyReferencingNode)
func (this_ *RemoveUnusedCode_StaticPropertyPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode pkg_rhino.Node) {
	_ = possiblyReferencingNode
	if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		this_.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable = false
	}
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.PrototypePropertyPolyfillInfo
type RemoveUnusedCode_PrototypePropertyPolyfillInfo struct {
	RemoveUnusedCode_PolyfillInfoBase
	outer_            *RemoveUnusedCode
	polyfillOwnerName string `gson:"polyfillOwnerName"`
}

func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_PolyfillInfoBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_PrototypePropertyPolyfillInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$PrototypePropertyPolyfillInfo", "PrototypePropertyPolyfillInfo"
}
func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_PolyfillInfoBase.IdentityHash_()
}
func (*RemoveUnusedCode_PrototypePropertyPolyfillInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.PrototypePropertyPolyfillInfo", "com.google.javascript.jscomp.RemoveUnusedCode$PrototypePropertyPolyfillInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_PrototypePropertyPolyfillInfo](v_)
	})
}
func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_PrototypePropertyPolyfillInfo) Implements_jscomp_RemoveUnusedCode_PolyfillInfo() {
}

func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) initFields_RemoveUnusedCode_PrototypePropertyPolyfillInfo() {
	_ = this_
	this_.polyfillOwnerName = j2gort.NullString
}

// Java constructor: PrototypePropertyPolyfillInfo(com.google.javascript.jscomp.RemoveUnusedCode.Polyfill removable, java.lang.String key, java.lang.String polyfillOwnerName)
func NewRemoveUnusedCode_PrototypePropertyPolyfillInfo(outer_ *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string, polyfillOwnerName string) *RemoveUnusedCode_PrototypePropertyPolyfillInfo {
	this_ := &RemoveUnusedCode_PrototypePropertyPolyfillInfo{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_PrototypePropertyPolyfillInfo(this_, removable, key, polyfillOwnerName)
	return this_
}
func allocPooled_RemoveUnusedCode_PrototypePropertyPolyfillInfo(outer_ *RemoveUnusedCode, removable *RemoveUnusedCode_Polyfill, key string, polyfillOwnerName string) *RemoveUnusedCode_PrototypePropertyPolyfillInfo {
	this_ := outer_.slab_RemoveUnusedCode_PrototypePropertyPolyfillInfo.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_PrototypePropertyPolyfillInfo(this_, removable, key, polyfillOwnerName)
	return this_
}

func InitRemoveUnusedCode_PrototypePropertyPolyfillInfo(this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo, removable *RemoveUnusedCode_Polyfill, key string, polyfillOwnerName string) {
	_ = removable
	_ = key
	_ = polyfillOwnerName
	this_.RemoveUnusedCode_PolyfillInfoBase.outer_ = this_.outer_
	InitRemoveUnusedCode_PolyfillInfo(&this_.RemoveUnusedCode_PolyfillInfoBase, removable, key)
	this_.initFields_RemoveUnusedCode_PrototypePropertyPolyfillInfo()
	this_.polyfillOwnerName = j2gort.CheckNotNullString(polyfillOwnerName)
}

// Java method: PrototypePropertyPolyfillInfo.getName()
func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) GetName() string {
	return (j2gort.StrValString(this_.polyfillOwnerName) + ".prototype." + j2gort.StrValString(this_.GetRemoveUnusedCode_PolyfillInfoBase().key))
}

// Java method: PrototypePropertyPolyfillInfo.considerPossibleReferenceInternal(com.google.javascript.rhino.Node possiblyReferencingNode)
func (this_ *RemoveUnusedCode_PrototypePropertyPolyfillInfo) ConsiderPossibleReferenceInternal(possiblyReferencingNode pkg_rhino.Node) {
	_ = possiblyReferencingNode
	if NodeUtil_IsNormalOrOptChainGetProp(possiblyReferencingNode) {
		if this_.outer_.disambiguateProperties && this_.outer_.polyfillUsageFinder.IsDisjointFromIterator_2(possiblyReferencingNode, this_.polyfillOwnerName) {
			return
		}
		this_.GetRemoveUnusedCode_PolyfillInfoBase().isRemovable = false
	}
}

// Java class: com.google.javascript.jscomp.RemoveUnusedCode.VanillaForNameDeclaration
type RemoveUnusedCode_VanillaForNameDeclaration struct {
	RemoveUnusedCode_RemovableBase
	outer_   *RemoveUnusedCode
	nameNode pkg_rhino.Node `gson:"nameNode"`
}

func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_VanillaForNameDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveUnusedCode$VanillaForNameDeclaration", "VanillaForNameDeclaration"
}
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RemoveUnusedCode_RemovableBase.IdentityHash_()
}
func (*RemoveUnusedCode_VanillaForNameDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveUnusedCode.VanillaForNameDeclaration", "com.google.javascript.jscomp.RemoveUnusedCode$VanillaForNameDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveUnusedCode_VanillaForNameDeclaration](v_)
	})
}
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) SetOuter_(o *RemoveUnusedCode) {
	this_.outer_ = o
}
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) GetOuter_() *RemoveUnusedCode {
	return this_.outer_
}
func (*RemoveUnusedCode_VanillaForNameDeclaration) Implements_jscomp_RemoveUnusedCode_Removable() {}

func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) initFields_RemoveUnusedCode_VanillaForNameDeclaration() {
	_ = this_
}

// Java constructor: VanillaForNameDeclaration(com.google.javascript.jscomp.RemoveUnusedCode.RemovableBuilder builder, com.google.javascript.rhino.Node nameNode)
func NewRemoveUnusedCode_VanillaForNameDeclaration(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, nameNode pkg_rhino.Node) *RemoveUnusedCode_VanillaForNameDeclaration {
	this_ := &RemoveUnusedCode_VanillaForNameDeclaration{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_VanillaForNameDeclaration(this_, builder, nameNode)
	return this_
}
func allocPooled_RemoveUnusedCode_VanillaForNameDeclaration(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, nameNode pkg_rhino.Node) *RemoveUnusedCode_VanillaForNameDeclaration {
	this_ := outer_.slab_RemoveUnusedCode_VanillaForNameDeclaration.Alloc()
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_VanillaForNameDeclaration(this_, builder, nameNode)
	return this_
}

func InitRemoveUnusedCode_VanillaForNameDeclaration(this_ *RemoveUnusedCode_VanillaForNameDeclaration, builder *RemoveUnusedCode_RemovableBuilder, nameNode pkg_rhino.Node) {
	_ = builder
	_ = nameNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, nameNode, builder)
	this_.initFields_RemoveUnusedCode_VanillaForNameDeclaration()
	this_.nameNode = nameNode
}

// Java method: VanillaForNameDeclaration.removeInternal(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) RemoveInternal(compiler AbstractCompiler) {
	_ = compiler
	var declaration pkg_rhino.Node = j2gort.CheckNotNullNode(this_.nameNode.GetParent())
	_ = declaration
	compiler.ReportChangeToEnclosingScope(declaration)
	if (this_.nameNode.GetPrevious() == 0) && (this_.nameNode.GetNext() == 0) {
		declaration.ReplaceWith(pkg_rhino.IR_Empty().Srcref(declaration))
	} else {
		this_.nameNode.Detach()
	}
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.nameNode, compiler)
}

// Java method: VanillaForNameDeclaration.isVariableAssignment()
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) IsVariableAssignment() bool {
	return true
}

// Java method: VanillaForNameDeclaration.isAssignedValueLocal()
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) IsAssignedValueLocal() bool {
	var initialValueNode pkg_rhino.Node = this_.nameNode.GetFirstChild()
	_ = initialValueNode
	if initialValueNode == 0 {
		return true
	}
	var valueNode pkg_rhino.Node = RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(this_.nameNode, initialValueNode)
	_ = valueNode
	return NodeUtil_EvaluatesToLocalValue(valueNode)
}

// Java method: VanillaForNameDeclaration.getLocalAssignedValue()
func (this_ *RemoveUnusedCode_VanillaForNameDeclaration) GetLocalAssignedValue() pkg_rhino.Node {
	var initialValueNode pkg_rhino.Node = this_.nameNode.GetFirstChild()
	_ = initialValueNode
	if initialValueNode == 0 {
		return pkg_rhino.Node(0)
	}
	var valueNode pkg_rhino.Node = RemoveUnusedCode_MaybeUnwrapQnameOrDefaultValueNode(this_.nameNode, initialValueNode)
	_ = valueNode
	if NodeUtil_EvaluatesToLocalValue(valueNode) {
		return valueNode
	} else {
		return pkg_rhino.Node(0)
	}
}
func (this_ *RemoveUnusedCode) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_RemoveUnusedCode_RemovableBuilder.Release()
	this_.slab_RemoveUnusedCode_GlobalPolyfillInfo.Release()
	this_.slab_RemoveUnusedCode_PrototypePropertyPolyfillInfo.Release()
	this_.slab_RemoveUnusedCode_StaticPropertyPolyfillInfo.Release()
	this_.slab_RemoveUnusedCode_IndirectAssign.Release()
	this_.slab_RemoveUnusedCode_Polyfill.Release()
	this_.slab_RemoveUnusedCode_ClassDeclaration.Release()
	this_.slab_RemoveUnusedCode_NamedClassExpression.Release()
	this_.slab_RemoveUnusedCode_ClassOrPrototypeNamedProperty.Release()
	this_.slab_RemoveUnusedCode_ObjectDefinePropertiesDefinition.Release()
	this_.slab_RemoveUnusedCode_FunctionDeclaration.Release()
	this_.slab_RemoveUnusedCode_NameDeclarationStatement.Release()
	this_.slab_RemoveUnusedCode_Assign.Release()
	this_.slab_RemoveUnusedCode_ClassSetupCall.Release()
	this_.slab_RemoveUnusedCode_VanillaForNameDeclaration.Release()
	this_.slab_RemoveUnusedCode_AnonymousPrototypeNamedPropertyAssign.Release()
	this_.slab_RemoveUnusedCode_IncOrDecOp.Release()
	this_.slab_RemoveUnusedCode_UnusedReadReference.Release()
	this_.slab_RemoveUnusedCode_InstanceofName.Release()
	this_.slab_RemoveUnusedCode_WeakUsageCall.Release()
	j2gort.ReleasePooledSlice(this_.varInfoVals)
	this_.varInfoVals = nil
	this_.slab_RemoveUnusedCode_VarInfo.Release()
	this_.allFunctionParamScopes = nil
	j2gort.ReleasePooledSlice(this_.varInfoKeys)
	this_.varInfoKeys = nil
	j2gort.ReleasePooledSlice(this_.varInfoTable)
	this_.varInfoTable = nil
	this_.varInfoSize = 0
	if this_.scopeCreator != nil {
		this_.scopeCreator.ReleaseSlabs_()
		this_.scopeCreator = nil
	}
}
