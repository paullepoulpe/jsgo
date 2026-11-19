/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ScopedAliases.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ScopedAliases
type ScopedAliases struct {
	compiler                 AbstractCompiler                              `gson:"compiler"`
	preprocessorSymbolTable  *PreprocessorSymbolTable                      `gson:"preprocessorSymbolTable"`
	moduleMetadataMap        *ModuleMetadataMap                            `gson:"moduleMetadataMap"`
	uniqueIdInput            *CompilerInput                                `gson:"uniqueIdInput"`
	scopedAliasNames         j2gort.Ext_com_google_common_collect_Multiset `gson:"scopedAliasNames"`
	closureNamespaces        j2gort.Ext_java_util_Set                      `gson:"closureNamespaces"`
	idHash_                  atomic.Int32
	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling `gson:"invalidModuleGetHandling"`
}

func (this_ *ScopedAliases) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases", "ScopedAliases"
}
func (this_ *ScopedAliases) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ScopedAliases) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases", "com.google.javascript.jscomp.ScopedAliases", func(v_ any) bool {
		return j2gort.InstanceOf[*ScopedAliases](v_)
	})
}
func (*ScopedAliases) Implements_jscomp_CompilerPass() {}

func StaticGet_ScopedAliases_SCOPING_METHOD_NAME() string    { return "goog.scope" }
func StaticGet_ScopedAliases_SCOPED_ALIASES_PREFIX() string  { return "$jscomp$scope$" }
func StaticGet_ScopedAliases_MISSING_ALIASES_PREFIX() string { return "$$jscomp$missingAlias$" }

var static_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE
}
func StaticSet_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE = v
}

var static_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE
}
func StaticSet_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE = v
}

var static_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS
}
func StaticSet_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS = v
}

var static_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS
}
func StaticSet_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS = v
}

var static_ScopedAliases_GOOG_SCOPE_USES_RETURN *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_USES_RETURN() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_USES_RETURN
}
func StaticSet_ScopedAliases_GOOG_SCOPE_USES_RETURN(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_USES_RETURN = v
}

var static_ScopedAliases_GOOG_SCOPE_USES_THROW *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_USES_THROW() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_USES_THROW
}
func StaticSet_ScopedAliases_GOOG_SCOPE_USES_THROW(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_USES_THROW = v
}

var static_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED
}
func StaticSet_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED = v
}

var static_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE
}
func StaticSet_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE = v
}

var static_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL
}
func StaticSet_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL = v
}

var static_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE *DiagnosticType

func StaticGet_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE() *DiagnosticType {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE
}
func StaticSet_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE(v *DiagnosticType) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE = v
}

var static_ScopedAliases_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_ScopedAliases_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_ScopedAliases()
	return static_ScopedAliases_GOOG_MODULE_GET
}
func StaticSet_ScopedAliases_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_ScopedAliases()
	static_ScopedAliases_GOOG_MODULE_GET = v
}

var clinitDone_ScopedAliases uint32
var clinitState_ScopedAliases j2gort.ClinitState

func clinit_ScopedAliases() {
	if atomic.LoadUint32(&clinitDone_ScopedAliases) != 0 {
		return
	}
	clinitSlow_ScopedAliases()
}

func clinitSlow_ScopedAliases() {
	if !j2gort.EnterClinit(&clinitState_ScopedAliases, &clinitDone_ScopedAliases) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ScopedAliases, &clinitDone_ScopedAliases)
	static_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE = DiagnosticType_Error("JSC_GOOG_SCOPE_MUST_BE_ALONE", "The call to goog.scope must be alone in a single statement.")
	static_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE = DiagnosticType_Error("JSC_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE", "The call to goog.scope must be in the global scope.")
	static_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS = DiagnosticType_Error("JSC_GOOG_SCOPE_HAS_BAD_PARAMETERS", "The call to goog.scope must take only a single parameter.  It must be an anonymous function that itself takes no parameters.")
	static_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS = DiagnosticType_Error("JSC_GOOG_SCOPE_REFERENCES_THIS", "The body of a goog.scope function cannot reference 'this'.")
	static_ScopedAliases_GOOG_SCOPE_USES_RETURN = DiagnosticType_Error("JSC_GOOG_SCOPE_USES_RETURN", "The body of a goog.scope function cannot use 'return'.")
	static_ScopedAliases_GOOG_SCOPE_USES_THROW = DiagnosticType_Error("JSC_GOOG_SCOPE_USES_THROW", "The body of a goog.scope function cannot use 'throw'.")
	static_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED = DiagnosticType_Error("JSC_GOOG_SCOPE_ALIAS_REDEFINED", "The alias {0} is assigned a value more than once.")
	static_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE = DiagnosticType_Error("JSC_GOOG_SCOPE_ALIAS_CYCLE", "The aliases {0} has a cycle.")
	static_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL = DiagnosticType_Error("JSC_GOOG_SCOPE_NON_ALIAS_LOCAL", "The local variable {0} is in a goog.scope and is not an alias.")
	static_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE = DiagnosticType_Error("JSC_GOOG_SCOPE_INVALID_VARIABLE", "The variable {0} cannot be declared in this scope")
	static_ScopedAliases_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
}

func (this_ *ScopedAliases) initFields_ScopedAliases() {
	_ = this_
	this_.scopedAliasNames = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset](j2gort.Ext_com_google_common_collect_HashMultiset_Create_0())
}

// Java constructor: ScopedAliases(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable PreprocessorSymbolTable preprocessorSymbolTable, java.util.Set<java.lang.String> closureNamespaces, com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling invalidModuleGetHandling)
func NewScopedAliases(compiler AbstractCompiler, preprocessorSymbolTable *PreprocessorSymbolTable, closureNamespaces j2gort.Ext_java_util_Set, invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling) *ScopedAliases {
	clinit_ScopedAliases()
	this_ := &ScopedAliases{}
	InitScopedAliases(this_, compiler, preprocessorSymbolTable, closureNamespaces, invalidModuleGetHandling)
	return this_
}

func InitScopedAliases(this_ *ScopedAliases, compiler AbstractCompiler, preprocessorSymbolTable *PreprocessorSymbolTable, closureNamespaces j2gort.Ext_java_util_Set, invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling) {
	_ = compiler
	_ = preprocessorSymbolTable
	_ = closureNamespaces
	_ = invalidModuleGetHandling
	this_.initFields_ScopedAliases()
	this_.compiler = compiler
	this_.preprocessorSymbolTable = preprocessorSymbolTable
	this_.moduleMetadataMap = compiler.GetModuleMetadataMap()
	this_.closureNamespaces = closureNamespaces
	this_.invalidModuleGetHandling = invalidModuleGetHandling
}

// Java static method: ScopedAliases.builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func ScopedAliases_S_Builder(compiler AbstractCompiler) *ScopedAliases_Builder {
	clinit_ScopedAliases()
	_ = compiler
	return NewScopedAliases_Builder(compiler)
}

// Java method: ScopedAliases.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ScopedAliases) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var traversal *ScopedAliases_Traversal = NewScopedAliases_Traversal(this_)
	_ = traversal
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer sc_.ReleaseSlabs_()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(traversal)
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(root)
	}

	if !(traversal.HasErrors()) {
		var aliasWorkQueue j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](traversal.priv_Traversal_getAliasUsages_5034()))
		_ = aliasWorkQueue
		for !(aliasWorkQueue.IsEmpty()) {
			var newQueue j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
			_ = newQueue
			{
				cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(aliasWorkQueue))
				for cur_1.HasNext() {
					rawElem_2 := cur_1.Next()
					aliasUsage := j2gort.AsIface[ScopedAliases_AliasUsage](rawElem_2)
					_ = aliasUsage
					if aliasUsage.GetScopedAliases_AliasUsageBase().ReferencesOtherAlias(traversal.deletedAliasVars) {
						newQueue.Add_1(aliasUsage)
					} else {
						aliasUsage.ApplyAlias(this_.compiler)
					}
				}
			}
			if newQueue.Size() == aliasWorkQueue.Size() {
				var cycleVar *Var = j2gort.AsIface[ScopedAliases_AliasUsage](newQueue.Get(int32(0))).GetScopedAliases_AliasUsageBase().aliasVar
				_ = cycleVar
				this_.compiler.Report_1(JSError_Make_3(cycleVar.GetNode(), StaticGet_ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE(), []string{cycleVar.GetName()}))
				break
			} else {
				aliasWorkQueue = newQueue
			}
		}
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(traversal.GetAliasDefinitionsToDelete()))
			for cur_3.HasNext() {
				aliasDefinition := pkg_rhino.Node(cur_3.NextNode())
				_ = aliasDefinition
				this_.compiler.ReportChangeToEnclosingScope(aliasDefinition)
				if NodeUtil_IsNameDeclaration(aliasDefinition.GetParent()) && aliasDefinition.GetParent().HasOneChild() {
					aliasDefinition.GetParent().Detach()
				} else {
					aliasDefinition.Detach()
				}
			}
		}
		{
			cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(traversal.GetScopeCalls()))
			for cur_4.HasNext() {
				scopeCall := pkg_rhino.Node(cur_4.NextNode())
				_ = scopeCall
				var expressionWithScopeCall pkg_rhino.Node = scopeCall.GetParent()
				_ = expressionWithScopeCall
				var scopeClosureBlock pkg_rhino.Node = scopeCall.GetLastChild().GetLastChild()
				_ = scopeClosureBlock
				scopeClosureBlock.Detach()
				expressionWithScopeCall.ReplaceWith(scopeClosureBlock)
				NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(expressionWithScopeCall, this_.compiler)
				this_.compiler.ReportChangeToEnclosingScope(scopeClosureBlock)
				NodeUtil_TryMergeBlock(scopeClosureBlock, false)
			}
		}
	}
}

// Java static method: ScopedAliases.isValidAliasRhs(com.google.javascript.rhino.Node rhs)
func ScopedAliases_IsValidAliasRhs(rhs pkg_rhino.Node) bool {
	clinit_ScopedAliases()
	_ = rhs
	switch rhs.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		return ScopedAliases_IsValidAliasRhs(rhs.GetFirstChild())
	case int32(32) /* Token.NAME */ :
		return true
	case int32(28) /* Token.CALL */ :
		return NodeUtil_IsCallTo_2_Node_QualifiedName(rhs, StaticGet_ScopedAliases_GOOG_MODULE_GET())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: ScopedAliases.isAliasDefinition(com.google.javascript.rhino.Node nameNode)
func ScopedAliases_IsAliasDefinition(nameNode pkg_rhino.Node) bool {
	clinit_ScopedAliases()
	_ = nameNode
	if !(nameNode.HasChildren()) {
		return false
	}
	var rhs pkg_rhino.Node = nameNode.GetLastChild()
	_ = rhs
	return ScopedAliases_IsValidAliasRhs(rhs)
}

// Java static method: ScopedAliases.getAliasedNamespace(com.google.javascript.rhino.Node rhs)
func ScopedAliases_GetAliasedNamespace(rhs pkg_rhino.Node) string {
	clinit_ScopedAliases()
	_ = rhs
	switch rhs.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		return (j2gort.StrValString(ScopedAliases_GetAliasedNamespace(rhs.GetFirstChild())) + "." + j2gort.StrValString(rhs.GetString()))
	case int32(32) /* Token.NAME */ :
		return rhs.GetString()
	case int32(28) /* Token.CALL */ :
		if !(NodeUtil_IsCallTo_2_Node_QualifiedName(rhs, StaticGet_ScopedAliases_GOOG_MODULE_GET())) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(rhs))
		}
		if !(rhs.HasTwoChildren()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(rhs))
		}
		return rhs.GetLastChild().GetString()
	default:
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Invalid alias RHS:"+j2gort.StrVal(rhs))))
	}
	panic("unreachable switch expression")
}

// Java static method: ScopedAliases.isScopedAliases(java.lang.String name)
func ScopedAliases_IsScopedAliases(name string) bool {
	clinit_ScopedAliases()
	_ = name
	return javalang.StringStartsWith(j2gort.CheckNotNullString(name), "$jscomp$scope$")
}

// Java enum: com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling
type ScopedAliases_InvalidModuleGetHandling uint8

const (
	enum_ScopedAliases_InvalidModuleGetHandling_PRESERVE         ScopedAliases_InvalidModuleGetHandling = 1
	enum_ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME ScopedAliases_InvalidModuleGetHandling = 2
)

var enumBaseTable_ScopedAliases_InvalidModuleGetHandling = [...]j2gort.EnumBase{
	{Name_: "PRESERVE", Ordinal_: 0},
	{Name_: "GIVE_UNIQUE_NAME", Ordinal_: 1},
}
var enumValues_ScopedAliases_InvalidModuleGetHandling = []ScopedAliases_InvalidModuleGetHandling{1, 2}

func StaticGet_ScopedAliases_InvalidModuleGetHandling_PRESERVE() ScopedAliases_InvalidModuleGetHandling {
	return enum_ScopedAliases_InvalidModuleGetHandling_PRESERVE
}
func StaticGet_ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME() ScopedAliases_InvalidModuleGetHandling {
	return enum_ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME
}

func ScopedAliases_InvalidModuleGetHandling_Values() []ScopedAliases_InvalidModuleGetHandling {
	return j2gort.CloneSlice(enumValues_ScopedAliases_InvalidModuleGetHandling)
}

func ScopedAliases_InvalidModuleGetHandling_ValueOf(name string) ScopedAliases_InvalidModuleGetHandling {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ScopedAliases_InvalidModuleGetHandling {
		if enumBaseTable_ScopedAliases_InvalidModuleGetHandling[i_].Name_ == name {
			return ScopedAliases_InvalidModuleGetHandling(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling." + name))
}

// Java method: InvalidModuleGetHandling.shouldGiveUniqueName()
func (this_ ScopedAliases_InvalidModuleGetHandling) ShouldGiveUniqueName() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	{

		re_5 := this_
		ae_6 := StaticGet_ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME()
		if re_5 == 0 {
			j2gort.ThrowNPE()
		}
		return re_5 == ae_6
	}

}

func (this_ ScopedAliases_InvalidModuleGetHandling) IsJavaNil_() bool                              { return this_ == 0 }
func (ScopedAliases_InvalidModuleGetHandling) HasIdentityEquals_()                                 {}
func (ScopedAliases_InvalidModuleGetHandling) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ScopedAliases_InvalidModuleGetHandling) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ScopedAliases_InvalidModuleGetHandling) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ScopedAliases_InvalidModuleGetHandling) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$InvalidModuleGetHandling", "InvalidModuleGetHandling"
}
func (this_ ScopedAliases_InvalidModuleGetHandling) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ScopedAliases_InvalidModuleGetHandling[this_-1]
}
func (this_ ScopedAliases_InvalidModuleGetHandling) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ScopedAliases_InvalidModuleGetHandling) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ScopedAliases_InvalidModuleGetHandling) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ScopedAliases_InvalidModuleGetHandling[this_-1].Name_
}
func (this_ ScopedAliases_InvalidModuleGetHandling) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ScopedAliases_InvalidModuleGetHandling[this_-1].Name_
}
func (this_ ScopedAliases_InvalidModuleGetHandling) Equals(other_ any) bool {
	o_, ok_ := other_.(ScopedAliases_InvalidModuleGetHandling)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ScopedAliases_InvalidModuleGetHandling) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ScopedAliases_InvalidModuleGetHandling[this_-1].HashCode()
}
func (this_ ScopedAliases_InvalidModuleGetHandling) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ScopedAliases_InvalidModuleGetHandling[this_-1].CompareTo(other_)
}
func (this_ ScopedAliases_InvalidModuleGetHandling) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ScopedAliases_InvalidModuleGetHandling[this_-1].Name_)
}
func init() {
	vals_ := []any{ScopedAliases_InvalidModuleGetHandling(1), ScopedAliases_InvalidModuleGetHandling(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ScopedAliases$InvalidModuleGetHandling", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling", "com.google.javascript.jscomp.ScopedAliases$InvalidModuleGetHandling", func(v_ any) bool {
		return j2gort.InstanceOf[ScopedAliases_InvalidModuleGetHandling](v_)
	})
}

// Java class: com.google.javascript.jscomp.ScopedAliases.Builder
type ScopedAliases_Builder struct {
	compiler                 AbstractCompiler         `gson:"compiler"`
	preprocessorSymbolTable  *PreprocessorSymbolTable `gson:"preprocessorSymbolTable"`
	moduleMetadataMap        *ModuleMetadataMap       `gson:"moduleMetadataMap"`
	idHash_                  atomic.Int32
	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling `gson:"invalidModuleGetHandling"`
}

func (this_ *ScopedAliases_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases_Builder) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$Builder", "Builder"
}
func (this_ *ScopedAliases_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ScopedAliases_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases.Builder", "com.google.javascript.jscomp.ScopedAliases$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*ScopedAliases_Builder](v_)
	})
}

func (this_ *ScopedAliases_Builder) initFields_ScopedAliases_Builder() {
	_ = this_

	this_.invalidModuleGetHandling = StaticGet_ScopedAliases_InvalidModuleGetHandling_PRESERVE()
}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewScopedAliases_Builder(compiler AbstractCompiler) *ScopedAliases_Builder {
	this_ := &ScopedAliases_Builder{}
	InitScopedAliases_Builder(this_, compiler)
	return this_
}

func InitScopedAliases_Builder(this_ *ScopedAliases_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ScopedAliases_Builder()
	this_.compiler = compiler
}

// Java method: Builder.setPreprocessorSymbolTable(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable PreprocessorSymbolTable preprocessorSymbolTable)
func (this_ *ScopedAliases_Builder) SetPreprocessorSymbolTable(preprocessorSymbolTable *PreprocessorSymbolTable) *ScopedAliases_Builder {
	_ = preprocessorSymbolTable
	this_.preprocessorSymbolTable = preprocessorSymbolTable
	return this_
}

// Java method: Builder.setModuleMetadataMap(com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func (this_ *ScopedAliases_Builder) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) *ScopedAliases_Builder {
	_ = moduleMetadataMap
	this_.moduleMetadataMap = moduleMetadataMap
	return this_
}

// Java method: Builder.setInvalidModuleGetHandling(com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling invalidModuleGetHandling)
func (this_ *ScopedAliases_Builder) SetInvalidModuleGetHandling(invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling) *ScopedAliases_Builder {
	_ = invalidModuleGetHandling
	this_.invalidModuleGetHandling = invalidModuleGetHandling
	return this_
}

// Java method: Builder.build()
func (this_ *ScopedAliases_Builder) Build() *ScopedAliases {
	return NewScopedAliases(this_.compiler, this_.preprocessorSymbolTable, func() j2gort.Ext_java_util_Set {
		if this_.moduleMetadataMap == nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
		}
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.moduleMetadataMap.GetModulesByGoogNamespace().KeySet_0())
	}(), this_.invalidModuleGetHandling)
}

// Java polymorphic class interface: com.google.javascript.jscomp.ScopedAliases.AliasUsage
type ScopedAliases_AliasUsage interface {
	Implements_jscomp_ScopedAliases_AliasUsage()
	IsJavaNil_() bool
	GetScopedAliases_AliasUsageBase() *ScopedAliases_AliasUsageBase
	ReferencesOtherAlias(j2gort.Ext_java_util_Set) bool
	ApplyAlias(AbstractCompiler)
}

func ToIface_ScopedAliases_AliasUsage[T interface {
	ScopedAliases_AliasUsage
	comparable
}](v T) ScopedAliases_AliasUsage {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ScopedAliases_AliasUsage[T ScopedAliases_AliasUsage](v T) ScopedAliases_AliasUsage {
	var vi ScopedAliases_AliasUsage = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.ScopedAliases.AliasUsage
type ScopedAliases_AliasUsageBase struct {
	aliasVar       *Var           `gson:"aliasVar"`
	aliasReference pkg_rhino.Node `gson:"aliasReference"`
	idHash_        atomic.Int32
}

func (this_ *ScopedAliases_AliasUsageBase) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases_AliasUsageBase) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases_AliasUsageBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$AliasUsage", "AliasUsage"
}
func (this_ *ScopedAliases_AliasUsageBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases.AliasUsage", "com.google.javascript.jscomp.ScopedAliases$AliasUsage", func(v_ any) bool {
		return j2gort.InstanceOf[ScopedAliases_AliasUsage](v_)
	})
}
func (b *ScopedAliases_AliasUsageBase) GetScopedAliases_AliasUsageBase() *ScopedAliases_AliasUsageBase {
	return b
}
func (*ScopedAliases_AliasUsageBase) Implements_jscomp_ScopedAliases_AliasUsage() {}

func (this_ *ScopedAliases_AliasUsageBase) initFields_ScopedAliases_AliasUsage() {
	_ = this_
}

func InitScopedAliases_AliasUsage(this_ *ScopedAliases_AliasUsageBase, aliasVar *Var, aliasReference pkg_rhino.Node) {
	_ = aliasVar
	_ = aliasReference
	this_.initFields_ScopedAliases_AliasUsage()
	this_.aliasVar = aliasVar
	this_.aliasReference = aliasReference
}

// Java method: AliasUsage.referencesOtherAlias(java.util.Set<com.google.javascript.jscomp.Var> deletedAliasVars)
func (this_ *ScopedAliases_AliasUsageBase) ReferencesOtherAlias(deletedAliasVars j2gort.Ext_java_util_Set) bool {
	_ = deletedAliasVars
	var aliasDefinition pkg_rhino.Node = this_.aliasVar.GetInitialValue()
	_ = aliasDefinition
	var qname string = ScopedAliases_GetAliasedNamespace(aliasDefinition)
	_ = qname
	var dotIndex int32 = j2gort.JStrIndexOfCp(qname, int32(uint16(46)))
	_ = dotIndex
	var rootName string
	_ = rootName
	if dotIndex == int32(-1) {
		rootName = qname
	} else {
		rootName = j2gort.JStrSub(qname, int32(0), dotIndex)
	}
	var otherAliasVar *Var = j2gort.CastPtr[Var](j2gort.CastPtr[Scope](this_.aliasVar.GetScope()).GetOwnSlot(rootName))
	_ = otherAliasVar
	return ((otherAliasVar != nil) && !(deletedAliasVars.Contains(j2gort.PtrToAny(otherAliasVar))))
}

// Java method: AliasUsage.applyAlias(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *ScopedAliases_AliasUsageBase) ApplyAlias(compiler AbstractCompiler) {
	_ = compiler
	panic("abstract or native method invoked")
}

// Java class: com.google.javascript.jscomp.ScopedAliases.AliasedNode
type ScopedAliases_AliasedNode struct {
	ScopedAliases_AliasUsageBase
}

func (this_ *ScopedAliases_AliasedNode) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases_AliasedNode) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases_AliasedNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$AliasedNode", "AliasedNode"
}
func (this_ *ScopedAliases_AliasedNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ScopedAliases_AliasUsageBase.IdentityHash_()
}
func (*ScopedAliases_AliasedNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases.AliasedNode", "com.google.javascript.jscomp.ScopedAliases$AliasedNode", func(v_ any) bool {
		return j2gort.InstanceOf[*ScopedAliases_AliasedNode](v_)
	})
}
func (*ScopedAliases_AliasedNode) Implements_jscomp_ScopedAliases_AliasUsage() {}

func (this_ *ScopedAliases_AliasedNode) initFields_ScopedAliases_AliasedNode() {
	_ = this_
}

// Java constructor: AliasedNode(com.google.javascript.jscomp.Var aliasVar, com.google.javascript.rhino.Node aliasReference)
func NewScopedAliases_AliasedNode(aliasVar *Var, aliasReference pkg_rhino.Node) *ScopedAliases_AliasedNode {
	this_ := &ScopedAliases_AliasedNode{}
	InitScopedAliases_AliasedNode(this_, aliasVar, aliasReference)
	return this_
}

func InitScopedAliases_AliasedNode(this_ *ScopedAliases_AliasedNode, aliasVar *Var, aliasReference pkg_rhino.Node) {
	_ = aliasVar
	_ = aliasReference
	InitScopedAliases_AliasUsage(&this_.ScopedAliases_AliasUsageBase, aliasVar, aliasReference)
	this_.initFields_ScopedAliases_AliasedNode()
}

// Java method: AliasedNode.applyAlias(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *ScopedAliases_AliasedNode) ApplyAlias(compiler AbstractCompiler) {
	_ = compiler
	var aliasDefinition pkg_rhino.Node = this_.GetScopedAliases_AliasUsageBase().aliasVar.GetInitialValue()
	_ = aliasDefinition
	var replacement pkg_rhino.Node = aliasDefinition.CloneTree_0()
	_ = replacement
	replacement.SrcrefTree(this_.GetScopedAliases_AliasUsageBase().aliasReference)
	if replacement.HasChildren() {
		replacement.GetFirstChild().MakeNonIndexableRecursive()
	}
	this_.GetScopedAliases_AliasUsageBase().aliasReference.ReplaceWith(replacement)
	compiler.ReportChangeToEnclosingScope(replacement)
}

// Java class: com.google.javascript.jscomp.ScopedAliases.AliasedTypeNode
type ScopedAliases_AliasedTypeNode struct {
	ScopedAliases_AliasUsageBase
}

func (this_ *ScopedAliases_AliasedTypeNode) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases_AliasedTypeNode) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases_AliasedTypeNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$AliasedTypeNode", "AliasedTypeNode"
}
func (this_ *ScopedAliases_AliasedTypeNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ScopedAliases_AliasUsageBase.IdentityHash_()
}
func (*ScopedAliases_AliasedTypeNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases.AliasedTypeNode", "com.google.javascript.jscomp.ScopedAliases$AliasedTypeNode", func(v_ any) bool {
		return j2gort.InstanceOf[*ScopedAliases_AliasedTypeNode](v_)
	})
}
func (*ScopedAliases_AliasedTypeNode) Implements_jscomp_ScopedAliases_AliasUsage() {}

func (this_ *ScopedAliases_AliasedTypeNode) initFields_ScopedAliases_AliasedTypeNode() {
	_ = this_
}

// Java constructor: AliasedTypeNode(com.google.javascript.jscomp.Var aliasVar, com.google.javascript.rhino.Node aliasReference)
func NewScopedAliases_AliasedTypeNode(aliasVar *Var, aliasReference pkg_rhino.Node) *ScopedAliases_AliasedTypeNode {
	this_ := &ScopedAliases_AliasedTypeNode{}
	InitScopedAliases_AliasedTypeNode(this_, aliasVar, aliasReference)
	return this_
}

func InitScopedAliases_AliasedTypeNode(this_ *ScopedAliases_AliasedTypeNode, aliasVar *Var, aliasReference pkg_rhino.Node) {
	_ = aliasVar
	_ = aliasReference
	InitScopedAliases_AliasUsage(&this_.ScopedAliases_AliasUsageBase, aliasVar, aliasReference)
	this_.initFields_ScopedAliases_AliasedTypeNode()
}

// Java method: AliasedTypeNode.applyAlias(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *ScopedAliases_AliasedTypeNode) ApplyAlias(compiler AbstractCompiler) {
	_ = compiler
	var aliasDefinition pkg_rhino.Node = this_.GetScopedAliases_AliasUsageBase().aliasVar.GetInitialValue()
	_ = aliasDefinition
	var aliasName string = this_.GetScopedAliases_AliasUsageBase().aliasVar.GetName()
	_ = aliasName
	var typeName string = this_.GetScopedAliases_AliasUsageBase().aliasReference.GetString()
	_ = typeName
	if javalang.StringStartsWith(j2gort.CheckNotNullString(typeName), "$jscomp$scope$") {
		return
	}
	var aliasExpanded string = j2gort.CheckNotNullString(ScopedAliases_GetAliasedNamespace(aliasDefinition))
	_ = aliasExpanded
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(typeName), j2gort.CheckNotNullString(aliasName))) {
		j2gort.CheckStateFmt(false, "%s must start with %s", j2gort.BoxString(typeName), j2gort.BoxString(aliasName))
	}
	var replacement string = (j2gort.StrValString(aliasExpanded) + j2gort.StrValString(j2gort.JStrSubFrom(typeName, j2gort.JStrLen(aliasName))))
	_ = replacement
	this_.GetScopedAliases_AliasUsageBase().aliasReference.SetString(replacement)
}

// Java class: com.google.javascript.jscomp.ScopedAliases.Traversal
type ScopedAliases_Traversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                   *ScopedAliases
	aliasDefinitionsToDelete j2gort.Ext_java_util_List         `gson:"aliasDefinitionsToDelete"`
	scopeCalls               j2gort.Ext_java_util_List         `gson:"scopeCalls"`
	aliasUsages              j2gort.Ext_java_util_List         `gson:"aliasUsages"`
	aliases                  j2gort.OrderedMap[string, *Var]   `gson:"aliases"`
	injectedDecls            j2gort.OrderedSet[pkg_rhino.Node] `gson:"injectedDecls"`
	deletedAliasVars         j2gort.Ext_java_util_Set          `gson:"deletedAliasVars"`
	forbiddenLocals          j2gort.Ext_java_util_Set          `gson:"forbiddenLocals"`
	scopeFunctionBody        pkg_rhino.Node                    `gson:"scopeFunctionBody"`
	hasNamespaceShadows      bool                              `gson:"hasNamespaceShadows"`
	hasErrors                bool                              `gson:"hasErrors"`
}

func (this_ *ScopedAliases_Traversal) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases_Traversal) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases_Traversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$Traversal", "Traversal"
}
func (this_ *ScopedAliases_Traversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ScopedAliases_Traversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopedAliases.Traversal", "com.google.javascript.jscomp.ScopedAliases$Traversal", func(v_ any) bool {
		return j2gort.InstanceOf[*ScopedAliases_Traversal](v_)
	})
}
func (this_ *ScopedAliases_Traversal) SetOuter_(o *ScopedAliases)                           { this_.outer_ = o }
func (this_ *ScopedAliases_Traversal) GetOuter_() *ScopedAliases                            { return this_.outer_ }
func (*ScopedAliases_Traversal) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ScopedAliases_Traversal) Implements_jscomp_NodeTraversal_ScopedCallback()            {}
func (*ScopedAliases_Traversal) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ScopedAliases_Traversal) initFields_ScopedAliases_Traversal() {
	_ = this_
	this_.aliasDefinitionsToDelete = j2gort.NewExt_java_util_ArrayList()
	this_.scopeCalls = j2gort.NewExt_java_util_ArrayList()
	this_.aliasUsages = j2gort.NewExt_java_util_ArrayList()
	this_.aliases = j2gort.NewOrderedMap[string, *Var]()
	this_.injectedDecls = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.deletedAliasVars = j2gort.NewExt_java_util_LinkedHashSet()
	this_.forbiddenLocals = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.ImmSetOfStr1("$jscomp")))

}

// Java constructor: Traversal()
func NewScopedAliases_Traversal(outer_ *ScopedAliases) *ScopedAliases_Traversal {
	this_ := &ScopedAliases_Traversal{}
	this_.outer_ = outer_
	InitScopedAliases_Traversal(this_)
	return this_
}

func InitScopedAliases_Traversal(this_ *ScopedAliases_Traversal) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ScopedAliases_Traversal()
}

// Java method: Traversal.getAliasDefinitionsToDelete()
func (this_ *ScopedAliases_Traversal) GetAliasDefinitionsToDelete() j2gort.Ext_java_util_Collection {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.aliasDefinitionsToDelete)
}

// Java method: Traversal.getAliasUsages()
func (this_ *ScopedAliases_Traversal) priv_Traversal_getAliasUsages_5034() j2gort.Ext_java_util_List {
	return this_.aliasUsages
}

// Java method: Traversal.getScopeCalls()
func (this_ *ScopedAliases_Traversal) GetScopeCalls() j2gort.Ext_java_util_List {
	return this_.scopeCalls
}

// Java method: Traversal.hasErrors()
func (this_ *ScopedAliases_Traversal) HasErrors() bool {
	return this_.hasErrors
}

// Java method: Traversal.inGoogScopeBody()
func (this_ *ScopedAliases_Traversal) priv_Traversal_inGoogScopeBody_5037() bool {
	return (this_.scopeFunctionBody != 0)
}

// Java method: Traversal.isGoogScopeFunctionBody(com.google.javascript.rhino.Node n)
func (this_ *ScopedAliases_Traversal) priv_Traversal_isGoogScopeFunctionBody_5038(n pkg_rhino.Node) bool {
	_ = n
	return (this_.priv_Traversal_inGoogScopeBody_5037() && (n == this_.scopeFunctionBody))
}

// Java method: Traversal.isCallToScopeMethod(com.google.javascript.rhino.Node n)
func (this_ *ScopedAliases_Traversal) priv_Traversal_isCallToScopeMethod_5039(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsCall() && n.GetFirstChild().MatchesQualifiedName_1_string("goog.scope"))
}

// Java method: Traversal.findScopeMethodCall(com.google.javascript.rhino.Node scopeRoot)
func (this_ *ScopedAliases_Traversal) priv_Traversal_findScopeMethodCall_5040(scopeRoot pkg_rhino.Node) pkg_rhino.Node {
	_ = scopeRoot
	var n pkg_rhino.Node = scopeRoot.GetGrandparent()
	_ = n
	if this_.priv_Traversal_isCallToScopeMethod_5039(n) {
		return n
	}
	return pkg_rhino.Node(0)
}

// Java method: Traversal.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ScopedAliases_Traversal) EnterScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalHoistScope() {
		return
	}
	var scopeMethodCall pkg_rhino.Node = this_.priv_Traversal_findScopeMethodCall_5040(t.GetScopeRoot())
	_ = scopeMethodCall
	if scopeMethodCall != 0 {
		this_.priv_Traversal_findAliases_5043(t.GetScope())
		this_.scopeFunctionBody = scopeMethodCall.GetLastChild().GetLastChild()
	}
}

// Java method: Traversal.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ScopedAliases_Traversal) ExitScope(t *NodeTraversal) {
	_ = t
	if this_.priv_Traversal_isGoogScopeFunctionBody_5038(t.GetScopeRoot()) {
		this_.scopeFunctionBody = pkg_rhino.Node(0)
		this_.priv_Traversal_renameNamespaceShadows_5047(t)
		this_.injectedDecls.Clear()
		this_.aliases.Clear()
		this_.forbiddenLocals.Clear()
		this_.hasNamespaceShadows = false
	} else if this_.priv_Traversal_inGoogScopeBody_5037() {
		this_.priv_Traversal_findNamespaceShadows_5046(t)
		this_.priv_Traversal_reportInvalidVariables_5041(t)
	}
}

// Java method: Traversal.reportInvalidVariables(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ScopedAliases_Traversal) priv_Traversal_reportInvalidVariables_5041(t *NodeTraversal) {
	_ = t
	var scopeRoot pkg_rhino.Node = t.GetScopeRoot()
	_ = scopeRoot
	var enclosingFunctionBody pkg_rhino.Node = t.GetEnclosingFunction().GetLastChild()
	_ = enclosingFunctionBody
	if (this_.priv_Traversal_isGoogScopeFunctionBody_5038(enclosingFunctionBody) && scopeRoot.IsBlock()) && !(scopeRoot.GetParent().IsFunction()) {
		{
			cur_7 := (t.GetScope()).iterVars_()
			for cur_7.HasNext() {

				v := cur_7.Next()

				_ = v
				var parent pkg_rhino.Node = v.GetNameNode().GetParent()
				_ = parent
				if NodeUtil_IsFunctionDeclaration(parent) {
					this_.priv_Traversal_report_5042(v.GetNode(), StaticGet_ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE(), []string{v.GetName()})
				}
			}
		}
	}
}

// Java method: Traversal.report(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType error, java.lang.String[] arguments)
func (this_ *ScopedAliases_Traversal) priv_Traversal_report_5042(n pkg_rhino.Node, error_ *DiagnosticType, arguments []string) {
	_ = n
	_ = error_
	_ = arguments
	this_.outer_.compiler.Report_1(JSError_Make_3(n, error_, arguments))
	this_.hasErrors = true
}

// Java method: Traversal.findAliases(com.google.javascript.jscomp.Scope scope)
func (this_ *ScopedAliases_Traversal) priv_Traversal_findAliases_5043(scope *Scope) {
	_ = scope
	{
		cur_9 := (scope).iterVars_()
		for cur_9.HasNext() {

			v := cur_9.Next()

			_ = v
			var n pkg_rhino.Node = v.GetNode()
			_ = n
			var parent pkg_rhino.Node = n.GetParent()
			_ = parent
			var isVar bool = (NodeUtil_IsNameDeclaration(parent) && parent.GetParent().IsBlock())
			_ = isVar
			var isFunctionDecl bool = NodeUtil_IsFunctionDeclaration(parent)
			_ = isFunctionDecl
			if isVar && ScopedAliases_IsAliasDefinition(n) {
				this_.priv_Traversal_recordAlias_5044(v)
			} else if v.IsBleedingFunction() {
			} else if parent.IsParamList() {
			} else if (isVar || isFunctionDecl) || NodeUtil_IsClassDeclaration(parent) {
				var isHoisted bool = NodeUtil_IsHoistedFunctionDeclaration(parent)
				_ = isHoisted
				var grandparent pkg_rhino.Node = parent.GetParent()
				_ = grandparent
				var value pkg_rhino.Node = v.GetInitialValue()
				_ = value
				var varNode pkg_rhino.Node = pkg_rhino.Node(0)
				_ = varNode
				if n.GetJSDocInfo() != nil {
					var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
					_ = builder
					if isFunctionDecl {
						builder.RecordReturnType(n.GetJSDocInfo().GetType())
					} else {
						builder.RecordType(n.GetJSDocInfo().GetType())
					}
					parent.SetJSDocInfo(builder.Build_0())
					n.SetJSDocInfo(nil)
				}
				var varDocInfo *pkg_rhino.JSDocInfo = v.GetJSDocInfo()
				_ = varDocInfo
				var name string = n.GetString()
				_ = name
				j2gort.MultisetAddStr(this_.outer_.scopedAliasNames, name)
				var uniqueId string = this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(this_.outer_.uniqueIdInput)
				_ = uniqueId
				var globalName string = ("$jscomp$scope$" + j2gort.StrValString(uniqueId) + "$" + j2gort.StrValString(name))
				_ = globalName
				if isFunctionDecl || NodeUtil_IsClassDeclaration(parent) {
					var newName pkg_rhino.Node = pkg_rhino.Node(0)
					_ = newName
					if isFunctionDecl {
						newName = pkg_rhino.IR_Name("")
					} else {
						newName = pkg_rhino.IR_Empty()
					}
					newName.Srcref(n)
					n.ReplaceWith(newName)
					this_.outer_.compiler.ReportChangeToEnclosingScope(newName)
					varNode = pkg_rhino.IR_Var_1(n).Srcref(n)
					parent.ReplaceWith(varNode)
				} else {
					if value != 0 {
						value.Detach()
					}
					varNode = parent
				}
				var newDecl pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node(globalName, value).SetJSDocInfo(varDocInfo).SrcrefTreeIfMissing(n)
				_ = newDecl
				newDecl.Srcref(n)
				newDecl.SetOriginalName(name)
				if isHoisted {
					grandparent.AddChildToFront(newDecl)
				} else {
					newDecl.InsertBefore(varNode)
				}
				if value != 0 {
					this_.outer_.compiler.ReportChangeToEnclosingScope(newDecl)
				}
				this_.injectedDecls.Add(newDecl)
				v.GetNameNode().AddChildToFront(NodeUtil_NewName_4_AbstractCompiler_string_Node_string(this_.outer_.compiler, globalName, n, name))
				this_.priv_Traversal_recordAlias_5044(v)
			} else {
				this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL(), []string{n.GetString()})
			}
		}
	}
}

// Java method: Traversal.recordAlias(com.google.javascript.jscomp.Var aliasVar)
func (this_ *ScopedAliases_Traversal) priv_Traversal_recordAlias_5044(aliasVar *Var) {
	_ = aliasVar
	var initialValue pkg_rhino.Node = aliasVar.GetInitialValue()
	_ = initialValue
	this_.aliasDefinitionsToDelete.Add_1(j2gort.BoxNode(aliasVar.GetNameNode()))
	if this_.outer_.invalidModuleGetHandling.ShouldGiveUniqueName() && this_.priv_Traversal_containsInvalidGoogModuleGet_5045(initialValue) {
		var name string = aliasVar.GetName()
		_ = name
		var fileName string = initialValue.GetSourceFileName()
		_ = fileName
		var metadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.outer_.moduleMetadataMap.GetModulesByPath(), fileName))
		_ = metadata
		var suffix string
		_ = suffix
		if metadata != nil && !metadata.IsJavaNil_() {
			suffix = metadata.Path().ToModuleName()
		} else {
			suffix = "unknownScript"
		}
		var newName string = (j2gort.StrValString(name) + "$$jscomp$missingAlias$" + j2gort.StrValString(suffix))
		_ = newName
		initialValue.ReplaceWith(pkg_rhino.IR_Name(newName))
		this_.deletedAliasVars.Add_1(j2gort.PtrToAny(aliasVar))
		this_.aliases.Put(name, aliasVar)
		return
	}
	var name_2 string = aliasVar.GetName()
	_ = name_2
	this_.aliases.Put(name_2, aliasVar)
	var qualifiedName string = ScopedAliases_GetAliasedNamespace(initialValue)
	_ = qualifiedName
	var rootIndex int32 = j2gort.JStrIndexOfCp(qualifiedName, int32(uint16(46)))
	_ = rootIndex
	if rootIndex != int32(-1) {
		var qNameRoot string = j2gort.JStrSub(qualifiedName, int32(0), rootIndex)
		_ = qNameRoot
		if !(this_.aliases.ContainsKey(qNameRoot)) {
			j2gort.SetAddStr(this_.forbiddenLocals, qNameRoot)
		}
	}
}

// Java method: Traversal.containsInvalidGoogModuleGet(com.google.javascript.rhino.Node expression)
func (this_ *ScopedAliases_Traversal) priv_Traversal_containsInvalidGoogModuleGet_5045(expression pkg_rhino.Node) bool {
	_ = expression
	switch expression.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return false
	case int32(26) /* Token.GETPROP */ :
		return this_.priv_Traversal_containsInvalidGoogModuleGet_5045(expression.GetFirstChild())
	case int32(28) /* Token.CALL */ :
		var namespace string = expression.GetSecondChild().GetString()
		_ = namespace
		return !(j2gort.SetContainsStr(this_.outer_.closureNamespaces, namespace))
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unrecognized alias rhs "+j2gort.StrVal(expression))))
	}
	panic("unreachable switch expression")
}

// Java method: Traversal.findNamespaceShadows(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ScopedAliases_Traversal) priv_Traversal_findNamespaceShadows_5046(t *NodeTraversal) {
	_ = t
	if this_.hasNamespaceShadows {
		return
	}
	var scope *Scope = t.GetScope()
	_ = scope
	{
		cur_11 := (scope).iterVars_()
		for cur_11.HasNext() {

			v := cur_11.Next()

			_ = v
			if j2gort.SetContainsStr(this_.forbiddenLocals, v.GetName()) {
				this_.hasNamespaceShadows = true
				return
			}
		}
	}
}

// Java method: Traversal.renameNamespaceShadows(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ScopedAliases_Traversal) priv_Traversal_renameNamespaceShadows_5047(t *NodeTraversal) {
	_ = t
	if !(NodeUtil_IsFunctionBlock(t.GetScopeRoot())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(t.GetScopeRoot()))
	}
	if this_.hasNamespaceShadows {
		var renamer MakeDeclaredNamesUnique_Renamer = NewMakeDeclaredNamesUnique_TargettedRenamer(Upcast_MakeDeclaredNamesUnique_Renamer(NewMakeDeclaredNamesUnique_ContextualRenamer()), this_.forbiddenLocals)
		_ = renamer
		{
			cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.forbiddenLocals))
			for cur_13.HasNext() {
				s := cur_13.NextStr()
				_ = s
				renamer.AddDeclaredName(s, false)
			}
		}
		var uniquifier *MakeDeclaredNamesUnique = MakeDeclaredNamesUnique_S_Builder().WithRenamer(renamer).Build()
		_ = uniquifier
		defer uniquifier.ReleaseSlabs_()

		NodeTraversal_TraverseScopeRoots(this_.outer_.compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(t.GetScopeRoot()))), ToIface_NodeTraversal_Callback(uniquifier), true)
	}
}

// Java method: Traversal.renameBleedingFunctionName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node fnName)
func (this_ *ScopedAliases_Traversal) priv_Traversal_renameBleedingFunctionName_5048(t *NodeTraversal, fnName pkg_rhino.Node) {
	_ = t
	_ = fnName
	var name string = fnName.GetString()
	_ = name
	var suffix string = this_.outer_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
	_ = suffix
	var cb NodeTraversal_Callback = NewScopedAliases__Anon_79(this_, name, fnName, suffix)
	_ = cb
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.outer_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(cb)
		(&ntb_).SetScopeCreator(t.GetScopeCreator())
		(&ntb_).TraverseAtScope(ToIface_AbstractScope(t.GetScope()))
	}
	fnName.SetString((j2gort.StrValString(name) + "$jscomp$scopedAliases$" + j2gort.StrValString(suffix)))
}

// Java method: Traversal.validateScopeCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ScopedAliases_Traversal) priv_Traversal_validateScopeCall_5049(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.outer_.preprocessorSymbolTable != nil {
		this_.outer_.preprocessorSymbolTable.AddReference_1(n.GetFirstChild())
	}
	if !(parent.IsExprResult()) {
		this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE(), []string{})
	}
	if t.GetEnclosingFunction() != 0 {
		this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE(), []string{})
	}
	if !(n.HasTwoChildren()) {
		this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS(), []string{})
	} else {
		var anonymousFnNode pkg_rhino.Node = n.GetSecondChild()
		_ = anonymousFnNode
		if (!(anonymousFnNode.IsFunction()) || (NodeUtil_GetName(anonymousFnNode) != j2gort.NullString)) || NodeUtil_GetFunctionParameters(anonymousFnNode).HasChildren() {
			this_.priv_Traversal_report_5042(anonymousFnNode, StaticGet_ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS(), []string{})
		} else {
			this_.scopeCalls.Add_1(j2gort.BoxNode(n))
		}
	}
}

// Java method: Traversal.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ScopedAliases_Traversal) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.outer_.uniqueIdInput = t.GetInput()
	if this_.priv_Traversal_isCallToScopeMethod_5039(n) {
		this_.priv_Traversal_validateScopeCall_5049(t, n, n.GetParent())
	}
	if !(this_.priv_Traversal_inGoogScopeBody_5037()) {
		return
	}
	var type_ pkg_rhino.Token = n.GetToken()
	_ = type_
	var aliasVar *Var = nil
	_ = aliasVar
	if type_ == pkg_rhino.StaticGet_Token_NAME() {
		var name string = n.GetString()
		_ = name
		var lexicalVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
		_ = lexicalVar
		if (lexicalVar != nil) && lexicalVar.Equals(j2gort.PtrToAny(this_.aliases.Get(name))) {
			aliasVar = lexicalVar
			if this_.outer_.compiler.GetOptions().PreservesDetailedSourceInfo() && n.IsName() {
				n.SetOriginalName(name)
			}
		}
	}
	if this_.priv_Traversal_isGoogScopeFunctionBody_5038(t.GetEnclosingFunction().GetLastChild()) {
		if (aliasVar != nil) && NodeUtil_IsLValue(n) {
			if aliasVar.GetNode() == n {
				return
			} else {
				this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED(), []string{n.GetString()})
			}
		}
		switch (type_).Ordinal() {
		case int32(0) /* Token.RETURN */ :
			this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_USES_RETURN(), []string{})
		case int32(37) /* Token.THIS */ :
			this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_REFERENCES_THIS(), []string{})
		case int32(43) /* Token.THROW */ :
			this_.priv_Traversal_report_5042(n, StaticGet_ScopedAliases_GOOG_SCOPE_USES_THROW(), []string{})
		default:
		}
	}
	if NodeUtil_IsBleedingFunctionName(n) {
		this_.priv_Traversal_renameBleedingFunctionName_5048(t, n)
	}
	if aliasVar != nil {
		this_.aliasUsages.Add_1(NewScopedAliases_AliasedNode(aliasVar, n))
	}
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if (info != nil) && !(this_.injectedDecls.Contains(n)) {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = info.GetTypeNodeCount()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_Traversal_fixTypeNode_5050(info.GetTypeNodeAt(i))
			}
		}
	}
}

// Java method: Traversal.fixTypeNode(com.google.javascript.rhino.Node typeNode)
func (this_ *ScopedAliases_Traversal) priv_Traversal_fixTypeNode_5050(typeNode pkg_rhino.Node) {
	_ = typeNode
	if typeNode.IsStringLit() {
		var name string = typeNode.GetString()
		_ = name
		var endIndex int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
		_ = endIndex
		if endIndex == int32(-1) {
			endIndex = j2gort.JStrLen(name)
		}
		var baseName string = j2gort.JStrSub(name, int32(0), endIndex)
		_ = baseName
		var aliasVar *Var = this_.aliases.Get(baseName)
		_ = aliasVar
		if aliasVar != nil {
			this_.aliasUsages.Add_1(NewScopedAliases_AliasedTypeNode(aliasVar, typeNode))
		}
		if this_.outer_.compiler.GetOptions().PreservesDetailedSourceInfo() {
			typeNode.SetOriginalName(name)
		}
	}
	{
		var child pkg_rhino.Node = typeNode.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_Traversal_fixTypeNode_5050(child)
		}
	}
}

// Java class:
type ScopedAliases__Anon_79 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_       *ScopedAliases_Traversal
	cap_0_name   string
	cap_2_suffix string
	cap_1_fnName pkg_rhino.Node
}

func (this_ *ScopedAliases__Anon_79) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopedAliases__Anon_79) IsJavaNil_() bool { return this_ == nil }
func (*ScopedAliases__Anon_79) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopedAliases$Traversal$1", ""
}
func (this_ *ScopedAliases__Anon_79) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ScopedAliases__Anon_79) HasIdentityEquals_()                                        {}
func (this_ *ScopedAliases__Anon_79) SetOuter_(o *ScopedAliases_Traversal)                 { this_.outer_ = o }
func (this_ *ScopedAliases__Anon_79) GetOuter_() *ScopedAliases_Traversal                  { return this_.outer_ }
func (*ScopedAliases__Anon_79) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ScopedAliases__Anon_79) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ScopedAliases__Anon_79) initFields_ScopedAliases__Anon_79() {
	_ = this_
}

func NewScopedAliases__Anon_79(outer_ *ScopedAliases_Traversal, cap_p_name string, cap_p_fnName pkg_rhino.Node, cap_p_suffix string) *ScopedAliases__Anon_79 {
	this_ := &ScopedAliases__Anon_79{}
	this_.outer_ = outer_
	this_.cap_0_name = cap_p_name
	this_.cap_1_fnName = cap_p_fnName
	this_.cap_2_suffix = cap_p_suffix
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ScopedAliases__Anon_79()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ScopedAliases__Anon_79) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (n.IsName() && func() bool {
		sr_14 := n.GetString()
		se_15 := this_.cap_0_name
		return j2gort.CheckNotNullString(sr_14) == se_15
	}()) && (j2gort.CastPtr[Var](t.GetScope().GetVar(this_.cap_0_name)).GetNode() == this_.cap_1_fnName) {
		n.SetString((j2gort.StrValString(this_.cap_0_name) + "$jscomp$scopedAliases$" + j2gort.StrValString(this_.cap_2_suffix)))
		this_.outer_.outer_.compiler.ReportChangeToEnclosingScope(n)
	}
}
func (*ScopedAliases_Traversal) AlwaysShouldTraverse_()     {}
func (*ScopedAliases__Anon_79) AlwaysShouldTraverse_()      {}
func (*ScopedAliases__Anon_79) CanReleaseSyntacticScopes_() {}
