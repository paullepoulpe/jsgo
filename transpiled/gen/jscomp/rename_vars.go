/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RenameVars.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_rename_vars_1 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_rename_vars_2 j2gort.Ext_java_util_Comparator

func init() {
	lambdaSingleton_rename_vars_1 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_1 any) int32 {
		arg := j2gort.CastPtr[RenameVars_Assignment](lp_1)
		_ = arg
		return arg.orderOfOccurrence
	})
	lambdaSingleton_rename_vars_2 = j2gort.NewLambda_Ext_java_util_Comparator(func(mrp0_ any, mrp1_ any) int32 {
		return RenameVars_FrequencyComparator(j2gort.CastPtr[RenameVars_Assignment](mrp0_), j2gort.CastPtr[RenameVars_Assignment](mrp1_))
	})
}

// Java class: com.google.javascript.jscomp.RenameVars
type RenameVars struct {
	compiler                       AbstractCompiler                                  `gson:"compiler"`
	globalNameNodes                []pkg_rhino.Node                                  `gson:"globalNameNodes"`
	localNameNodes                 []pkg_rhino.Node                                  `gson:"localNameNodes"`
	originalNameByNode             j2gort.OrderedMap[pkg_rhino.Node, string]         `gson:"originalNameByNode"`
	pseudoNameMap                  *j2gort.OrderedMap[pkg_rhino.Node, string]        `gson:"pseudoNameMap"`
	externNames                    j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"externNames"`
	reservedNames                  j2gort.Ext_java_util_Set                          `gson:"reservedNames"`
	renameMap                      j2gort.OrderedMap[string, string]                 `gson:"renameMap"`
	prevUsedRenameMap              *VariableMap                                      `gson:"prevUsedRenameMap"`
	prefix                         string                                            `gson:"prefix"`
	localBleedingFunctions         j2gort.CustomOrderedSet[*Var]                     `gson:"localBleedingFunctions"`
	localBleedingFunctionsPerScope j2gort.Ext_com_google_common_collect_ListMultimap `gson:"localBleedingFunctionsPerScope"`
	assignments                    j2gort.OrderedMap[string, *RenameVars_Assignment] `gson:"assignments"`
	reservedCharacters             j2gort.Ext_java_util_Set                          `gson:"reservedCharacters"`
	nameGenerator                  NameGenerator                                     `gson:"nameGenerator"`
	assignmentCount                int32                                             `gson:"assignmentCount"`
	idHash_                        atomic.Int32
	localRenamingOnly              bool `gson:"localRenamingOnly"`
	preferStableNames              bool `gson:"preferStableNames"`
}

func (this_ *RenameVars) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameVars) IsJavaNil_() bool { return this_ == nil }
func (*RenameVars) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameVars", "RenameVars"
}
func (this_ *RenameVars) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameVars) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameVars", "com.google.javascript.jscomp.RenameVars", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameVars](v_)
	})
}
func (*RenameVars) Implements_jscomp_CompilerPass() {}

func StaticGet_RenameVars_MAX_LOCALS_IN_SCOPE_TO_TEMP_RENAME() int32 { return int32(1000) }
func StaticGet_RenameVars_LOCAL_VAR_PREFIX() string                  { return "L " }

var static_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_RenameVars()
	return static_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR
}
func StaticSet_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_RenameVars()
	static_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR = v
}

var clinitDone_RenameVars uint32
var clinitState_RenameVars j2gort.ClinitState

func clinit_RenameVars() {
	if atomic.LoadUint32(&clinitDone_RenameVars) != 0 {
		return
	}
	clinitSlow_RenameVars()
}

func clinitSlow_RenameVars() {
	if !j2gort.EnterClinit(&clinitState_RenameVars, &clinitDone_RenameVars) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RenameVars, &clinitDone_RenameVars)
	static_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR = j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_rename_vars_1)
}

func (this_ *RenameVars) initFields_RenameVars() {
	_ = this_
	this_.globalNameNodes = []pkg_rhino.Node{}
	this_.localNameNodes = []pkg_rhino.Node{}
	this_.originalNameByNode = j2gort.NewOrderedMap[pkg_rhino.Node, string]()
	this_.renameMap = j2gort.NewOrderedMap[string, string]()
	this_.prefix = j2gort.NullString

	this_.localBleedingFunctions = j2gort.NewCustomOrderedSet[*Var]()
	this_.localBleedingFunctionsPerScope = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
	this_.assignments = j2gort.NewOrderedMap[string, *RenameVars_Assignment](

	// Java constructor: RenameVars(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String prefix, boolean localRenamingOnly, boolean generatePseudoNames, boolean preferStableNames, com.google.javascript.jscomp.VariableMap prevUsedRenameMap, java.util.Set<java.lang.Character> reservedCharacters, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> reservedNames, com.google.javascript.jscomp.NameGenerator nameGenerator)
	)
}

func NewRenameVars(compiler AbstractCompiler, prefix string, localRenamingOnly bool, generatePseudoNames bool, preferStableNames bool, prevUsedRenameMap *VariableMap, reservedCharacters j2gort.Ext_java_util_Set, reservedNames j2gort.Ext_java_util_Set, nameGenerator NameGenerator) *RenameVars {
	clinit_RenameVars()
	this_ := &RenameVars{}
	InitRenameVars(this_, compiler, prefix, localRenamingOnly, generatePseudoNames, preferStableNames, prevUsedRenameMap, reservedCharacters, reservedNames, nameGenerator)
	return this_
}

func InitRenameVars(this_ *RenameVars, compiler AbstractCompiler, prefix string, localRenamingOnly bool, generatePseudoNames bool, preferStableNames bool, prevUsedRenameMap *VariableMap, reservedCharacters j2gort.Ext_java_util_Set, reservedNames j2gort.Ext_java_util_Set, nameGenerator NameGenerator) {
	_ = compiler
	_ = prefix
	_ = localRenamingOnly
	_ = generatePseudoNames
	_ = preferStableNames
	_ = prevUsedRenameMap
	_ = reservedCharacters
	_ = reservedNames
	_ = nameGenerator
	this_.initFields_RenameVars()
	this_.compiler = compiler
	this_.prefix = j2gort.Ext_com_google_common_base_Strings_NullToEmpty(prefix)
	this_.localRenamingOnly = localRenamingOnly
	if generatePseudoNames {
		this_.pseudoNameMap = &j2gort.OrderedMap[pkg_rhino.Node, string]{}
	} else {
		this_.pseudoNameMap = nil
	}
	this_.prevUsedRenameMap = prevUsedRenameMap
	this_.reservedCharacters = reservedCharacters
	this_.preferStableNames = preferStableNames
	if reservedNames == nil {
		this_.reservedNames = j2gort.NewExt_java_util_LinkedHashSet()
	} else {
		this_.reservedNames = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](reservedNames))
	}
	this_.nameGenerator = nameGenerator
}

// Java static method: RenameVars.frequencyComparator(com.google.javascript.jscomp.RenameVars.Assignment a1, com.google.javascript.jscomp.RenameVars.Assignment a2)
func RenameVars_FrequencyComparator(a1 *RenameVars_Assignment, a2 *RenameVars_Assignment) int32 {
	clinit_RenameVars()
	_ = a1
	_ = a2
	if a1.count != a2.count {
		return (a2.count - a1.count)
	}
	return StaticGet_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR().Compare(j2gort.PtrToAny(a1), j2gort.PtrToAny(a2))
}

// Java method: RenameVars.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RenameVars) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.externNames = NodeUtil_CollectExternVariableNames(this_.compiler, externs)
	this_.originalNameByNode.Clear()
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		this_.localBleedingFunctions.Clear()
		this_.localBleedingFunctionsPerScope.Clear()
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewRenameVars_ProcessVars(this_))
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(root)
	}

	this_.reservedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.externNames))
	var varsByFrequency j2gort.Ext_java_util_SortedSet = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(lambdaSingleton_rename_vars_2)
	_ = varsByFrequency
	varsByFrequency.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.assignments.LiveValues()))
	if this_.prevUsedRenameMap != nil {
		this_.priv_RenameVars_reusePreviouslyUsedVariableMap_4803(varsByFrequency)
	}
	this_.priv_RenameVars_assignNames_4804(varsByFrequency)
	for _, n := range this_.globalNameNodes {
		_ = n
		this_.priv_RenameVars_setNameAndReport_4799(n, this_.priv_RenameVars_getNewGlobalName_4800(n))
	}

	for _, n_2 := range this_.localNameNodes {
		_ = n_2
		this_.priv_RenameVars_setNameAndReport_4799(n_2, this_.priv_RenameVars_getNewLocalName_4801(n_2))
	}

}

// Java method: RenameVars.setNameAndReport(com.google.javascript.rhino.Node n, java.lang.@org.jspecify.annotations.Nullable String newName)
func (this_ *RenameVars) priv_RenameVars_setNameAndReport_4799(n pkg_rhino.Node, newName string) {
	_ = n
	_ = newName
	if (newName != j2gort.NullString) && !(func() bool { sr_4 := newName; se_5 := n.GetString(); return j2gort.CheckNotNullString(sr_4) == se_5 }()) {
		n.SetString(newName)
		{
			var iifeCond_1 bool
			{

				sr_6 := newName
				se_7 := j2gort.OrderedMapGetStr(this_.originalNameByNode, n)
				iifeCond_1 = j2gort.CheckNotNullString(sr_6) == se_7
			}
			if !(iifeCond_1) {
				this_.compiler.ReportChangeToEnclosingScope(n)
				var parent pkg_rhino.Node = n.GetParent()
				_ = parent
				if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
					this_.compiler.ReportChangeToEnclosingScope(parent)
				}
			}
		}
	}
}

// Java method: RenameVars.getNewGlobalName(com.google.javascript.rhino.Node n)
func (this_ *RenameVars) priv_RenameVars_getNewGlobalName_4800(n pkg_rhino.Node) string {
	_ = n
	var oldName string = n.GetString()
	_ = oldName
	var a *RenameVars_Assignment = this_.assignments.Get(oldName)
	_ = a
	if (a.newName != j2gort.NullString) && !(j2gort.CheckNotNullString(a.newName) == oldName) {
		if this_.pseudoNameMap != nil {
			return j2gort.OrderedMapPtrGetStr(this_.pseudoNameMap, n)
		}
		return a.newName
	} else {
		return j2gort.NullString
	}
}

// Java method: RenameVars.getNewLocalName(com.google.javascript.rhino.Node n)
func (this_ *RenameVars) priv_RenameVars_getNewLocalName_4801(n pkg_rhino.Node) string {
	_ = n
	var oldTempName string = n.GetString()
	_ = oldTempName
	var a *RenameVars_Assignment = this_.assignments.Get(oldTempName)
	_ = a
	if !(j2gort.CheckNotNullString(a.newName) == oldTempName) {
		if this_.pseudoNameMap != nil {
			return j2gort.OrderedMapPtrGetStr(this_.pseudoNameMap, n)
		}
		return a.newName
	}
	return j2gort.NullString
}

// Java method: RenameVars.recordPseudoName(com.google.javascript.rhino.Node n)
func (this_ *RenameVars) priv_RenameVars_recordPseudoName_4802(n pkg_rhino.Node) {
	_ = n
	this_.pseudoNameMap.Put(n, "$"+j2gort.StrValString(n.GetString())+"$$")
}

// Java method: RenameVars.reusePreviouslyUsedVariableMap(java.util.SortedSet<com.google.javascript.jscomp.RenameVars.Assignment> varsToRename)
func (this_ *RenameVars) priv_RenameVars_reusePreviouslyUsedVariableMap_4803(varsToRename j2gort.Ext_java_util_SortedSet) {
	_ = varsToRename
	j2gort.CheckNotNull(this_.prevUsedRenameMap.GetNewNameToOriginalNameMap())
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(varsToRename))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			a := j2gort.CastPtr[RenameVars_Assignment](rawElem_9)
			_ = a
			var prevNewName string = this_.prevUsedRenameMap.LookupNewName(a.oldName)
			_ = prevNewName
			if (prevNewName == j2gort.NullString) || j2gort.SetContainsStr(this_.reservedNames, prevNewName) {
				continue
			}
			if a.isLocal || (!(j2gort.SetContainsStr(this_.externNames, a.oldName)) && javalang.StringStartsWith(j2gort.CheckNotNullString(prevNewName), j2gort.CheckNotNullString(this_.prefix))) {
				j2gort.SetAddStr(this_.reservedNames, prevNewName)
				this_.priv_RenameVars_finalizeNameAssignment_4805(a, prevNewName)
			}
		}
	}
}

// Java method: RenameVars.assignNames(java.util.SortedSet<com.google.javascript.jscomp.RenameVars.Assignment> varsToRename)
func (this_ *RenameVars) priv_RenameVars_assignNames_4804(varsToRename j2gort.Ext_java_util_SortedSet) {
	_ = varsToRename
	this_.nameGenerator.Reset_3(this_.reservedNames, this_.prefix, this_.reservedCharacters)
	var globalNameGenerator NameGenerator = this_.nameGenerator
	_ = globalNameGenerator
	var localNameGenerator NameGenerator
	_ = localNameGenerator
	if j2gort.CheckNotNullString(this_.prefix) == "" {
		localNameGenerator = globalNameGenerator
	} else {
		localNameGenerator = this_.nameGenerator.Clone(this_.reservedNames, "", this_.reservedCharacters)
	}
	var pendingAssignments []*RenameVars_Assignment = []*RenameVars_Assignment{}
	_ = pendingAssignments
	var generatedNamesForAssignments []string = []string{}
	_ = generatedNamesForAssignments
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(varsToRename))
		for cur_10.HasNext() {
			rawElem_11 := cur_10.Next()
			a := j2gort.CastPtr[RenameVars_Assignment](rawElem_11)
			_ = a
			if a.newName != j2gort.NullString {
				continue
			}
			if j2gort.SetContainsStr(this_.externNames, a.oldName) {
				continue
			}
			var newName string = j2gort.NullString
			_ = newName
			if a.isLocal {
				newName = localNameGenerator.GenerateNextName()
				this_.priv_RenameVars_finalizeNameAssignment_4805(a, newName)
			} else {
				newName = globalNameGenerator.GenerateNextName()
				pendingAssignments = append(pendingAssignments, a)
				generatedNamesForAssignments = append(generatedNamesForAssignments, newName)
			}
			j2gort.SetAddStr(this_.reservedNames, newName)
		}
	}
	var numPendingAssignments int32 = int32(len(generatedNamesForAssignments))
	_ = numPendingAssignments
	{
		var i int32 = int32(0)
		_ = i
		for i < numPendingAssignments {
			var varsByOrderOfOccurrence j2gort.Ext_java_util_SortedSet = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(StaticGet_RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR())
			_ = varsByOrderOfOccurrence
			var len_ int32 = j2gort.JStrLen(generatedNamesForAssignments[i])
			_ = len_
			{
				var j int32 = i
				_ = j
				for ; (j < numPendingAssignments) && (j2gort.JStrLen(generatedNamesForAssignments[j]) == len_); j = j + 1 {
					varsByOrderOfOccurrence.Add_1(j2gort.PtrToAny(pendingAssignments[j]))
				}
			}
			{
				cur_12 := j2gort.NewIterCursor(varsByOrderOfOccurrence)
				for cur_12.HasNext() {
					rawElem_13 := cur_12.Next()
					a_2 := j2gort.CastPtr[RenameVars_Assignment](rawElem_13)
					_ = a_2
					this_.priv_RenameVars_finalizeNameAssignment_4805(a_2, generatedNamesForAssignments[i])
					i = i + 1
				}
			}
		}
	}
}

// Java method: RenameVars.finalizeNameAssignment(com.google.javascript.jscomp.RenameVars.Assignment a, java.lang.String newName)
func (this_ *RenameVars) priv_RenameVars_finalizeNameAssignment_4805(a *RenameVars_Assignment, newName string) {
	_ = a
	_ = newName
	a.SetNewName(newName)
	this_.renameMap.Put(a.oldName, newName)
}

// Java method: RenameVars.getVariableMap()
func (this_ *RenameVars) GetVariableMap() *VariableMap {
	return NewVariableMap(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.renameMap.ToImmutableMap()))
}

// Java method: RenameVars.okToRenameVar(java.lang.String name, boolean isLocal)
func (this_ *RenameVars) priv_RenameVars_okToRenameVar_4807(name string, isLocal bool) bool {
	_ = name
	_ = isLocal
	return !(this_.compiler.GetCodingConvention().IsExported_2(name, isLocal))
}

// Java method: RenameVars.getLocalVarIndex(com.google.javascript.jscomp.Var v)
func (this_ *RenameVars) priv_RenameVars_getLocalVarIndex_4808(v *Var) int32 {
	_ = v
	var num int32 = v.GetIndex()
	_ = num
	var s *Scope = j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](v.GetScope()).GetParent())
	_ = s
	if s == nil {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Var is not local")))
	}
	var isBleedingIntoScope bool = ((j2gort.CastPtr[Scope](s.GetParent()) != nil) && this_.localBleedingFunctions.Contains(v))
	_ = isBleedingIntoScope
	for j2gort.CastPtr[Scope](s.GetParent()) != nil {
		if isBleedingIntoScope {
			num = (num + (j2gort.AsIface[j2gort.Ext_java_util_List](this_.localBleedingFunctionsPerScope.Get(j2gort.PtrToAny(s))).IndexOf_1(j2gort.PtrToAny(v)) + int32(1)))
			isBleedingIntoScope = false
		} else {
			num = (num + j2gort.AsIface[j2gort.Ext_java_util_List](this_.localBleedingFunctionsPerScope.Get(j2gort.PtrToAny(s))).Size())
		}
		if this_.priv_RenameVars_shouldTemporarilyRenameLocalsInScope_4809(s) {
			num = (num + s.GetVarCount())
		}
		s = j2gort.CastPtr[Scope](s.GetParent())
	}
	return num
}

// Java method: RenameVars.shouldTemporarilyRenameLocalsInScope(com.google.javascript.jscomp.Scope s)
func (this_ *RenameVars) priv_RenameVars_shouldTemporarilyRenameLocalsInScope_4809(s *Scope) bool {
	_ = s
	return (!(this_.preferStableNames) || (s.GetVarCount() <= int32(1000)))
}

// Java class: com.google.javascript.jscomp.RenameVars.Assignment
type RenameVars_Assignment struct {
	outer_            *RenameVars
	oldName           string `gson:"oldName"`
	newName           string `gson:"newName"`
	orderOfOccurrence int32  `gson:"orderOfOccurrence"`
	count             int32  `gson:"count"`
	idHash_           atomic.Int32
	isLocal           bool `gson:"isLocal"`
}

func (this_ *RenameVars_Assignment) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameVars_Assignment) IsJavaNil_() bool { return this_ == nil }
func (*RenameVars_Assignment) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameVars$Assignment", "Assignment"
}
func (this_ *RenameVars_Assignment) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameVars_Assignment) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameVars.Assignment", "com.google.javascript.jscomp.RenameVars$Assignment", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameVars_Assignment](v_)
	})
}
func (this_ *RenameVars_Assignment) SetOuter_(o *RenameVars) { this_.outer_ = o }
func (this_ *RenameVars_Assignment) GetOuter_() *RenameVars  { return this_.outer_ }

func (this_ *RenameVars_Assignment) initFields_RenameVars_Assignment() {
	_ = this_
	this_.oldName = j2gort.NullString
	this_.newName = j2gort.NullString
}

// Java constructor: Assignment(java.lang.String name)
func NewRenameVars_Assignment(outer_ *RenameVars, name string) *RenameVars_Assignment {
	this_ := &RenameVars_Assignment{}
	this_.outer_ = outer_
	InitRenameVars_Assignment(this_, name)
	return this_
}

func InitRenameVars_Assignment(this_ *RenameVars_Assignment, name string) {
	_ = name
	this_.initFields_RenameVars_Assignment()
	this_.isLocal = javalang.StringStartsWith(j2gort.CheckNotNullString(name), "L ")
	this_.oldName = name
	this_.newName = j2gort.NullString
	this_.count = int32(0)
	this_.orderOfOccurrence = func() int32 {
		old_14 := this_.outer_.assignmentCount
		nw_15 := old_14 + 1
		this_.outer_.assignmentCount = nw_15
		return old_14
	}()
}

// Java method: Assignment.setNewName(java.lang.String newName)
func (this_ *RenameVars_Assignment) SetNewName(newName string) {
	_ = newName
	if !(this_.newName == j2gort.NullString) {
		j2gort.CheckState(false)
	}
	this_.newName = newName
}

// Java class: com.google.javascript.jscomp.RenameVars.ProcessVars
type RenameVars_ProcessVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RenameVars
}

func (this_ *RenameVars_ProcessVars) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameVars_ProcessVars) IsJavaNil_() bool { return this_ == nil }
func (*RenameVars_ProcessVars) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameVars$ProcessVars", "ProcessVars"
}
func (this_ *RenameVars_ProcessVars) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RenameVars_ProcessVars) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameVars.ProcessVars", "com.google.javascript.jscomp.RenameVars$ProcessVars", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameVars_ProcessVars](v_)
	})
}
func (this_ *RenameVars_ProcessVars) SetOuter_(o *RenameVars)                              { this_.outer_ = o }
func (this_ *RenameVars_ProcessVars) GetOuter_() *RenameVars                               { return this_.outer_ }
func (*RenameVars_ProcessVars) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*RenameVars_ProcessVars) Implements_jscomp_NodeTraversal_ScopedCallback()            {}
func (*RenameVars_ProcessVars) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *RenameVars_ProcessVars) initFields_RenameVars_ProcessVars() {
	_ = this_
}

// Java constructor: ProcessVars()
func NewRenameVars_ProcessVars(outer_ *RenameVars) *RenameVars_ProcessVars {
	this_ := &RenameVars_ProcessVars{}
	this_.outer_ = outer_
	InitRenameVars_ProcessVars(this_)
	return this_
}

func InitRenameVars_ProcessVars(this_ *RenameVars_ProcessVars) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RenameVars_ProcessVars()
}

// Java method: ProcessVars.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RenameVars_ProcessVars) EnterScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalHoistScope() || !(this_.outer_.priv_RenameVars_shouldTemporarilyRenameLocalsInScope_4809(t.GetScope())) {
		return
	}
	var scope *Scope = t.GetScope()
	_ = scope
	{
		cur_16 := (scope).iterVars_()
		for cur_16.HasNext() {

			current := cur_16.Next()

			_ = current
			if current.IsBleedingFunction() {
				this_.outer_.localBleedingFunctions.Add(current)
				this_.outer_.localBleedingFunctionsPerScope.Put(j2gort.PtrToAny(j2gort.CastPtr[Scope](scope.GetParent())), j2gort.PtrToAny(current))
			}
		}
	}
}

// Java method: ProcessVars.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RenameVars_ProcessVars) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java method: ProcessVars.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RenameVars_ProcessVars) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsName() || n.IsImportStar()) {
		return
	}
	var name string = n.GetString()
	_ = name
	if j2gort.CheckNotNullString(name) == "" {
		return
	}
	if (parent.IsImportSpec() && parent.HasTwoChildren()) && (parent.GetFirstChild() == n) {
		return
	}
	var var_ *Var = t.GetScope().LookupVar(name)
	_ = var_
	var local bool = (((var_ != nil) && var_.IsLocal()) && (j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](var_.GetScope()).GetParent()).IsLocal() || !(var_.IsBleedingFunction())))
	_ = local
	if (var_ != nil) && var_.IsArguments() {
		j2gort.SetAddStr(this_.outer_.reservedNames, name)
		return
	}
	if !(local) && this_.outer_.localRenamingOnly {
		j2gort.SetAddStr(this_.outer_.reservedNames, name)
		return
	}
	if !(this_.outer_.priv_RenameVars_okToRenameVar_4807(name, local)) {
		if local {
			var newName string = MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)
			_ = newName
			if !(j2gort.CheckNotNullString(newName) == name) {
				n.SetString(newName)
			}
		}
		return
	}
	if this_.outer_.pseudoNameMap != nil {
		this_.outer_.priv_RenameVars_recordPseudoName_4802(n)
	}
	if local && this_.outer_.priv_RenameVars_shouldTemporarilyRenameLocalsInScope_4809(j2gort.CastPtr[Scope](var_.GetScope())) {
		var tempName string = ("L " + j2gort.StrValInt32(this_.outer_.priv_RenameVars_getLocalVarIndex_4808(var_)))
		_ = tempName
		this_.IncCount(tempName)
		this_.outer_.localNameNodes = append(this_.outer_.localNameNodes, n)
		this_.outer_.originalNameByNode.Put(n, n.GetString())
		n.SetString(tempName)
	} else if var_ != nil {
		this_.IncCount(name)
		this_.outer_.globalNameNodes = append(this_.outer_.globalNameNodes, n)
	}
}

// Java method: ProcessVars.incCount(java.lang.String name)
func (this_ *RenameVars_ProcessVars) IncCount(name string) {
	_ = name
	var s *RenameVars_Assignment = j2gort.OrderedMapComputeIfAbsent(&this_.outer_.assignments, name, nil, func(mrp0_ string) *RenameVars_Assignment {
		return NewRenameVars_Assignment(this_.outer_, mrp0_)
	})
	_ = s
	r_18 := s
	r_18.count = r_18.count + 1
}
func (*RenameVars_ProcessVars) AlwaysShouldTraverse_() {}
