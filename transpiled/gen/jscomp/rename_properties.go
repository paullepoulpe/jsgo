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
// Port of com/google/javascript/jscomp/RenameProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_rename_properties_1 j2gort.Ext_java_util_Comparator
var lambdaSingleton_rename_properties_2 j2gort.Ext_java_util_function_Predicate

func init() {
	lambdaSingleton_rename_properties_1 = j2gort.NewLambda_Ext_java_util_Comparator(func(lp_1 any, lp_2 any) int32 {
		p1 := j2gort.CastPtr[RenameProperties_Property](lp_1)
		_ = p1
		p2 := j2gort.CastPtr[RenameProperties_Property](lp_2)
		_ = p2
		if p1.numOccurrences != p2.numOccurrences {
			return (p2.numOccurrences - p1.numOccurrences)
		}
		return j2gort.JStrCompareTo(p1.oldName, p2.oldName)
	})
	lambdaSingleton_rename_properties_2 = j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_3 uint32) bool {
		n := pkg_rhino.Node(rawNode_3)
		_ = n
		return true
	})
}

// Java class: com.google.javascript.jscomp.RenameProperties
type RenameProperties struct {
	compiler                        AbstractCompiler                                      `gson:"compiler"`
	prevUsedPropertyMap             *VariableMap                                          `gson:"prevUsedPropertyMap"`
	toRemove                        []pkg_rhino.Node                                      `gson:"toRemove"`
	stringNodesToRename             []pkg_rhino.Node                                      `gson:"stringNodesToRename"`
	callNodeToParentMap             j2gort.OrderedMap[pkg_rhino.Node, pkg_rhino.Node]     `gson:"callNodeToParentMap"`
	reservedFirstCharacters         j2gort.Ext_java_util_Set                              `gson:"reservedFirstCharacters"`
	reservedNonFirstCharacters      j2gort.Ext_java_util_Set                              `gson:"reservedNonFirstCharacters"`
	propertyMap                     j2gort.OrderedMap[string, *RenameProperties_Property] `gson:"propertyMap"`
	externedNames                   j2gort.Ext_java_util_Set                              `gson:"externedNames"`
	quotedNames                     j2gort.Ext_java_util_Set                              `gson:"quotedNames"`
	nameGenerator                   NameGenerator                                         `gson:"nameGenerator"`
	propertyRenameEligibilityFilter j2gort.Ext_java_util_function_Predicate               `gson:"propertyRenameEligibilityFilter"`
	filteredOutNames                j2gort.OrderedSet[string]                             `gson:"filteredOutNames"`
	idHash_                         atomic.Int32
	generatePseudoNames             bool `gson:"generatePseudoNames"`
}

func (this_ *RenameProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameProperties) IsJavaNil_() bool { return this_ == nil }
func (*RenameProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameProperties", "RenameProperties"
}
func (this_ *RenameProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameProperties", "com.google.javascript.jscomp.RenameProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameProperties](v_)
	})
}
func (*RenameProperties) Implements_jscomp_CompilerPass() {}

var static_RenameProperties_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_RenameProperties_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_RenameProperties()
	return static_RenameProperties_DOT_SPLITTER
}
func StaticSet_RenameProperties_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_RenameProperties()
	static_RenameProperties_DOT_SPLITTER = v
}

var static_RenameProperties_FREQUENCY_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_RenameProperties_FREQUENCY_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_RenameProperties()
	return static_RenameProperties_FREQUENCY_COMPARATOR
}
func StaticSet_RenameProperties_FREQUENCY_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_RenameProperties()
	static_RenameProperties_FREQUENCY_COMPARATOR = v
}

var static_RenameProperties_BAD_CALL *DiagnosticType

func StaticGet_RenameProperties_BAD_CALL() *DiagnosticType {
	clinit_RenameProperties()
	return static_RenameProperties_BAD_CALL
}
func StaticSet_RenameProperties_BAD_CALL(v *DiagnosticType) {
	clinit_RenameProperties()
	static_RenameProperties_BAD_CALL = v
}

var static_RenameProperties_BAD_ARG *DiagnosticType

func StaticGet_RenameProperties_BAD_ARG() *DiagnosticType {
	clinit_RenameProperties()
	return static_RenameProperties_BAD_ARG
}
func StaticSet_RenameProperties_BAD_ARG(v *DiagnosticType) {
	clinit_RenameProperties()
	static_RenameProperties_BAD_ARG = v
}

var clinitDone_RenameProperties uint32
var clinitState_RenameProperties j2gort.ClinitState

func clinit_RenameProperties() {
	if atomic.LoadUint32(&clinitDone_RenameProperties) != 0 {
		return
	}
	clinitSlow_RenameProperties()
}

func clinitSlow_RenameProperties() {
	if !j2gort.EnterClinit(&clinitState_RenameProperties, &clinitDone_RenameProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RenameProperties, &clinitDone_RenameProperties)
	static_RenameProperties_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46))
	static_RenameProperties_FREQUENCY_COMPARATOR = lambdaSingleton_rename_properties_1
	static_RenameProperties_BAD_CALL = DiagnosticType_Error("JSC_BAD_RENAME_PROPERTY_FUNCTION_NAME_CALL", "Bad {0} call - the first argument must be a string literal")
	static_RenameProperties_BAD_ARG = DiagnosticType_Error("JSC_BAD_RENAME_PROPERTY_FUNCTION_NAME_ARG", "Bad {0} argument - ''{1}'' is not a valid JavaScript identifier")
}

func (this_ *RenameProperties) initFields_RenameProperties() {
	_ = this_
	this_.toRemove = []pkg_rhino.Node{}
	this_.stringNodesToRename = []pkg_rhino.Node{}
	this_.callNodeToParentMap = j2gort.NewOrderedMap[pkg_rhino.Node, pkg_rhino.Node]()
	this_.propertyMap = j2gort.NewOrderedMap[string, *RenameProperties_Property]()
	this_.externedNames = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList([]any{"prototype"})))
	this_.quotedNames = j2gort.NewExt_java_util_LinkedHashSet()
	this_.filteredOutNames = j2gort.NewOrderedSet[string]()
}

// Java constructor: RenameProperties(com.google.javascript.jscomp.AbstractCompiler compiler, boolean generatePseudoNames, com.google.javascript.jscomp.VariableMap prevUsedPropertyMap, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters, com.google.javascript.jscomp.NameGenerator nameGenerator)
func NewRenameProperties_6(compiler AbstractCompiler, generatePseudoNames bool, prevUsedPropertyMap *VariableMap, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, nameGenerator NameGenerator) *RenameProperties {
	clinit_RenameProperties()
	this_ := &RenameProperties{}
	InitRenameProperties_6(this_, compiler, generatePseudoNames, prevUsedPropertyMap, reservedFirstCharacters, reservedNonFirstCharacters, nameGenerator)
	return this_
}

func InitRenameProperties_6(this_ *RenameProperties, compiler AbstractCompiler, generatePseudoNames bool, prevUsedPropertyMap *VariableMap, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, nameGenerator NameGenerator) {
	_ = compiler
	_ = generatePseudoNames
	_ = prevUsedPropertyMap
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	_ = nameGenerator
	InitRenameProperties_7(this_, compiler, generatePseudoNames, prevUsedPropertyMap, reservedFirstCharacters, reservedNonFirstCharacters, nameGenerator, lambdaSingleton_rename_properties_2)
}

// Java constructor: RenameProperties(com.google.javascript.jscomp.AbstractCompiler compiler, boolean generatePseudoNames, com.google.javascript.jscomp.VariableMap prevUsedPropertyMap, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters, com.google.javascript.jscomp.NameGenerator nameGenerator, java.util.function.Predicate<com.google.javascript.rhino.Node> propertyRenameEligibilityFilter)
func NewRenameProperties_7(compiler AbstractCompiler, generatePseudoNames bool, prevUsedPropertyMap *VariableMap, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, nameGenerator NameGenerator, propertyRenameEligibilityFilter j2gort.Ext_java_util_function_Predicate) *RenameProperties {
	clinit_RenameProperties()
	this_ := &RenameProperties{}
	InitRenameProperties_7(this_, compiler, generatePseudoNames, prevUsedPropertyMap, reservedFirstCharacters, reservedNonFirstCharacters, nameGenerator, propertyRenameEligibilityFilter)
	return this_
}

func InitRenameProperties_7(this_ *RenameProperties, compiler AbstractCompiler, generatePseudoNames bool, prevUsedPropertyMap *VariableMap, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, nameGenerator NameGenerator, propertyRenameEligibilityFilter j2gort.Ext_java_util_function_Predicate) {
	_ = compiler
	_ = generatePseudoNames
	_ = prevUsedPropertyMap
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	_ = nameGenerator
	_ = propertyRenameEligibilityFilter
	this_.initFields_RenameProperties()
	this_.compiler = compiler
	this_.generatePseudoNames = generatePseudoNames
	this_.prevUsedPropertyMap = prevUsedPropertyMap
	this_.reservedFirstCharacters = reservedFirstCharacters
	this_.reservedNonFirstCharacters = reservedNonFirstCharacters
	this_.nameGenerator = nameGenerator
	this_.propertyRenameEligibilityFilter = propertyRenameEligibilityFilter
	this_.externedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetExternProperties()))
}

// Java method: RenameProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RenameProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	NodeTraversal_Traverse(this_.compiler, root, NewRenameProperties_ProcessProperties(this_))
	var reservedNames j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_NewLinkedHashSetWithExpectedSize((this_.externedNames.Size() + this_.quotedNames.Size())))
	_ = reservedNames
	reservedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.externedNames))
	reservedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.quotedNames))
	var propsByFreq j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(StaticGet_RenameProperties_FREQUENCY_COMPARATOR())
	_ = propsByFreq
	propsByFreq.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.propertyMap.LiveValues()))
	if this_.prevUsedPropertyMap != nil {
		this_.priv_RenameProperties_reusePropertyNames_4824(reservedNames, propsByFreq)
	}
	this_.priv_RenameProperties_generateNames_4825(propsByFreq, reservedNames)
	for _, n := range this_.stringNodesToRename {
		_ = n
		var oldName string = n.GetString()
		_ = oldName
		if this_.filteredOutNames.Contains(oldName) {
			continue
		}
		var p *RenameProperties_Property = this_.propertyMap.Get(oldName)
		_ = p
		if (p != nil) && (p.newName != j2gort.NullString) {
			{
				var iifeCond_1 bool
				{

					sr_5 := oldName
					se_6 := p.oldName
					iifeCond_1 = j2gort.CheckNotNullString(sr_5) == se_6
				}
				if !(iifeCond_1) {
					j2gort.CheckState(false)
				}
			}
			n.SetString(p.newName)
			if !(j2gort.CheckNotNullString(p.newName) == oldName) {
				this_.compiler.ReportChangeToEnclosingScope(n)
			}
		}
	}

	{
		cur_7 := this_.callNodeToParentMap.IterEntries()
		for cur_7.NextEntry() {
			var parent pkg_rhino.Node = cur_7.Value()
			_ = parent
			var firstArg pkg_rhino.Node = cur_7.Key().GetSecondChild()
			_ = firstArg
			var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
			_ = sb
			{
				cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.SplitterSplitStr(StaticGet_RenameProperties_DOT_SPLITTER(), firstArg.GetString())))
				for cur_8.HasNext() {
					oldName_2 := cur_8.NextStr()
					_ = oldName_2
					var p_2 *RenameProperties_Property = this_.propertyMap.Get(oldName_2)
					_ = p_2
					var replacement string = j2gort.NullString
					_ = replacement
					if (p_2 != nil) && (p_2.newName != j2gort.NullString) {
						{
							var iifeCond_1 bool
							{

								sr_9 := oldName_2
								se_10 := p_2.oldName
								iifeCond_1 = j2gort.CheckNotNullString(sr_9) == se_10
							}
							if !(iifeCond_1) {
								j2gort.CheckState(false)
							}
						}
						if this_.filteredOutNames.Contains(oldName_2) {
							replacement = oldName_2
						} else {
							replacement = p_2.newName
						}
					} else {
						replacement = oldName_2
					}
					if sb.Length() > int32(0) {
						sb.Append_1_uint16(uint16(46))
					}
					sb.Append_1_string(replacement)
				}
			}
			cur_7.Key().ReplaceWith(pkg_rhino.IR_String(sb.String()))
			this_.compiler.ReportChangeToEnclosingScope(parent)
		}
	}
	for _, n_2 := range this_.toRemove {
		_ = n_2
		var parent_2 pkg_rhino.Node = n_2.GetParent()
		_ = parent_2
		this_.compiler.ReportChangeToEnclosingScope(n_2)
		n_2.Detach()
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n_2, this_.compiler)
		if !(parent_2.HasChildren()) && !(parent_2.IsScript()) {
			parent_2.Detach()
		}
	}

	this_.compiler.SetLifeCycleStage(StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED())
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java method: RenameProperties.reusePropertyNames(java.util.Set<java.lang.String> reservedNames, java.util.Collection<com.google.javascript.jscomp.RenameProperties.Property> allProps)
func (this_ *RenameProperties) priv_RenameProperties_reusePropertyNames_4824(reservedNames j2gort.Ext_java_util_Set, allProps j2gort.Ext_java_util_Collection) {
	_ = reservedNames
	_ = allProps
	{
		cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(allProps))
		for cur_12.HasNext() {
			rawElem_13 := cur_12.Next()
			prop := j2gort.CastPtr[RenameProperties_Property](rawElem_13)
			_ = prop
			var prevName string = this_.prevUsedPropertyMap.LookupNewName(prop.oldName)
			_ = prevName
			if !(this_.generatePseudoNames) && (prevName != j2gort.NullString) {
				if j2gort.SetContainsStr(reservedNames, prevName) {
					continue
				}
				{
					r_14 := prop
					r_14.newName = prevName
				}
				j2gort.SetAddStr(reservedNames, prevName)
			}
		}
	}
}

// Java method: RenameProperties.generateNames(java.util.Set<com.google.javascript.jscomp.RenameProperties.Property> props, java.util.Set<java.lang.String> reservedNames)
func (this_ *RenameProperties) priv_RenameProperties_generateNames_4825(props j2gort.Ext_java_util_Set, reservedNames j2gort.Ext_java_util_Set) {
	_ = props
	_ = reservedNames
	this_.nameGenerator.Reset_4(reservedNames, "", this_.reservedFirstCharacters, this_.reservedNonFirstCharacters)
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(props))
		for cur_15.HasNext() {
			rawElem_16 := cur_15.Next()
			p := j2gort.CastPtr[RenameProperties_Property](rawElem_16)
			_ = p
			if this_.generatePseudoNames {
				{
					r_17 := p
					r_17.newName = ("$" + j2gort.StrValString(p.oldName) + "$")
				}
			} else {
				if p.newName == j2gort.NullString {
					{
						r_18 := p
						r_18.newName = this_.nameGenerator.GenerateNextName()
					}
				}
			}
			j2gort.SetAddStr(reservedNames, p.newName)
		}
	}
}

// Java method: RenameProperties.getPropertyMap()
func (this_ *RenameProperties) GetPropertyMap() *VariableMap {
	var map_ j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = map_
	{
		cur_19 := this_.propertyMap.IterValues()
		for cur_19.HasNext() {

			p := cur_19.
				Next()
			_ = p
			if p.newName != j2gort.NullString {
				map_.Put_2(j2gort.BoxString(p.oldName), j2gort.BoxString(p.newName))
			}
		}
	}
	return NewVariableMap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](map_.BuildOrThrow()))
}

// Java class: com.google.javascript.jscomp.RenameProperties.ProcessProperties
type RenameProperties_ProcessProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RenameProperties
}

func (this_ *RenameProperties_ProcessProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameProperties_ProcessProperties) IsJavaNil_() bool { return this_ == nil }
func (*RenameProperties_ProcessProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameProperties$ProcessProperties", "ProcessProperties"
}
func (this_ *RenameProperties_ProcessProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RenameProperties_ProcessProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameProperties.ProcessProperties", "com.google.javascript.jscomp.RenameProperties$ProcessProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameProperties_ProcessProperties](v_)
	})
}
func (this_ *RenameProperties_ProcessProperties) SetOuter_(o *RenameProperties) { this_.outer_ = o }
func (this_ *RenameProperties_ProcessProperties) GetOuter_() *RenameProperties  { return this_.outer_ }
func (*RenameProperties_ProcessProperties) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RenameProperties_ProcessProperties) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RenameProperties_ProcessProperties) initFields_RenameProperties_ProcessProperties() {
	_ = this_
}

// Java constructor: ProcessProperties()
func NewRenameProperties_ProcessProperties(outer_ *RenameProperties) *RenameProperties_ProcessProperties {
	this_ := &RenameProperties_ProcessProperties{}
	this_.outer_ = outer_
	InitRenameProperties_ProcessProperties(this_)
	return this_
}

func InitRenameProperties_ProcessProperties(this_ *RenameProperties_ProcessProperties) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RenameProperties_ProcessProperties()
}

// Java method: ProcessProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RenameProperties_ProcessProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
		this_.priv_ProcessProperties_maybeMarkCandidate_4827(n)
	case int32(27) /* Token.GETELEM */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		var child pkg_rhino.Node = n.GetLastChild()
		_ = child
		if (child != 0) && child.IsStringLit() {
			j2gort.SetAddStr(this_.outer_.quotedNames, child.GetString())
		}
	case int32(28) /* Token.CALL */ :
		var fnName pkg_rhino.Node = n.GetFirstChild()
		_ = fnName
		if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(fnName) {
			this_.outer_.callNodeToParentMap.Put(n, parent)
			this_.priv_ProcessProperties_countCallCandidates_4828(t, n)
			var stringArgument pkg_rhino.Node = n.GetSecondChild()
			_ = stringArgument
			if !(j2gort.PredTestNode(this_.outer_.propertyRenameEligibilityFilter, uint32(n))) {
				{
					cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.SplitterSplitStr(StaticGet_RenameProperties_DOT_SPLITTER(), stringArgument.GetString())))
					for cur_21.HasNext() {
						name := cur_21.NextStr()
						_ = name
						this_.outer_.filteredOutNames.Add(name)
					}
				}
			}
		}
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		if !(!(n.IsQuotedStringKey())) {
			j2gort.CheckState(false)
		}
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
			j2gort.SetAddStr(this_.outer_.externedNames, n.GetString())
		} else {
			this_.priv_ProcessProperties_maybeMarkCandidate_4827(n)
		}
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(100) /* Token.STRING_KEY */ :
		if n.IsQuotedStringKey() {
			j2gort.SetAddStr(this_.outer_.quotedNames, n.GetString())
		} else {
			this_.priv_ProcessProperties_maybeMarkCandidate_4827(n)
		}
	case int32(73) /* Token.FUNCTION */ :
		if NodeUtil_IsFunctionDeclaration(n) {
			var name_2 string = n.GetFirstChild().GetString()
			_ = name_2
			if "JSCompiler_renameProperty" == name_2 {
				this_.outer_.toRemove = append(this_.outer_.toRemove, n)
			}
		} else if parent.IsName() && func() bool {
			sr_22 := "JSCompiler_renameProperty"
			se_23 := parent.GetString()
			return j2gort.CheckNotNullString(sr_22) == se_23
		}() {
			var varNode pkg_rhino.Node = parent.GetParent()
			_ = varNode
			if varNode.IsVar() {
				this_.outer_.toRemove = append(this_.outer_.toRemove, parent)
			}
		} else if ((NodeUtil_IsFunctionExpression(n) && parent.IsAssign()) && parent.GetFirstChild().IsGetProp()) && this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(parent.GetFirstChild()) {
			var exprResult pkg_rhino.Node = parent.GetParent()
			_ = exprResult
			if (exprResult.IsExprResult() && NodeUtil_IsStatementBlock(exprResult.GetParent())) && exprResult.GetFirstChild().IsAssign() {
				this_.outer_.toRemove = append(this_.outer_.toRemove, exprResult)
			}
		}
	default:
	}
}

// Java method: ProcessProperties.maybeMarkCandidate(com.google.javascript.rhino.Node n)
func (this_ *RenameProperties_ProcessProperties) priv_ProcessProperties_maybeMarkCandidate_4827(n pkg_rhino.Node) {
	_ = n
	var name string = n.GetString()
	_ = name
	if !(j2gort.SetContainsStr(this_.outer_.externedNames, name)) {
		if !(j2gort.PredTestNode(this_.outer_.propertyRenameEligibilityFilter, uint32(n))) {
			this_.outer_.filteredOutNames.Add(n.GetString())
		}
		this_.outer_.stringNodesToRename = append(this_.outer_.stringNodesToRename, n)
		this_.priv_ProcessProperties_countPropertyOccurrence_4829(name)
	}
}

// Java method: ProcessProperties.countCallCandidates(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode)
func (this_ *RenameProperties_ProcessProperties) priv_ProcessProperties_countCallCandidates_4828(t *NodeTraversal, callNode pkg_rhino.Node) {
	_ = t
	_ = callNode
	var fnName string = callNode.GetFirstChild().GetString()
	_ = fnName
	var firstArg pkg_rhino.Node = callNode.GetSecondChild()
	_ = firstArg
	if !(firstArg.IsStringLit()) {
		t.Report_3(callNode, StaticGet_RenameProperties_BAD_CALL(), []string{fnName})
		return
	}
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.SplitterSplitStr(StaticGet_RenameProperties_DOT_SPLITTER(), firstArg.GetString())))
		for cur_24.HasNext() {
			name := cur_24.NextStr()
			_ = name
			if !(pkg_rhino.TokenStream_IsJSIdentifier(name)) {
				t.Report_3(callNode, StaticGet_RenameProperties_BAD_ARG(), []string{fnName})
				continue
			}
			if !(j2gort.SetContainsStr(this_.outer_.externedNames, name)) {
				this_.priv_ProcessProperties_countPropertyOccurrence_4829(name)
			}
		}
	}
}

// Java method: ProcessProperties.countPropertyOccurrence(java.lang.String name)
func (this_ *RenameProperties_ProcessProperties) priv_ProcessProperties_countPropertyOccurrence_4829(name string) {
	_ = name
	var prop *RenameProperties_Property = this_.outer_.propertyMap.Get(name)
	_ = prop
	if prop == nil {
		prop = NewRenameProperties_Property(name)
		this_.outer_.propertyMap.Put(name, prop)
	}
	r_25 := prop
	r_25.numOccurrences = r_25.numOccurrences + 1
}

// Java class: com.google.javascript.jscomp.RenameProperties.Property
type RenameProperties_Property struct {
	oldName        string `gson:"oldName"`
	newName        string `gson:"newName"`
	numOccurrences int32  `gson:"numOccurrences"`
	idHash_        atomic.Int32
}

func (this_ *RenameProperties_Property) BindSelf_(s any) {
	_ = s
}
func (this_ *RenameProperties_Property) IsJavaNil_() bool { return this_ == nil }
func (*RenameProperties_Property) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RenameProperties$Property", "Property"
}
func (this_ *RenameProperties_Property) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RenameProperties_Property) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RenameProperties.Property", "com.google.javascript.jscomp.RenameProperties$Property", func(v_ any) bool {
		return j2gort.InstanceOf[*RenameProperties_Property](v_)
	})
}

func (this_ *RenameProperties_Property) initFields_RenameProperties_Property() {
	_ = this_
	this_.oldName = j2gort.NullString
	this_.newName = j2gort.NullString
}

// Java constructor: Property(java.lang.String name)
func NewRenameProperties_Property(name string) *RenameProperties_Property {
	this_ := &RenameProperties_Property{}
	InitRenameProperties_Property(this_, name)
	return this_
}

func InitRenameProperties_Property(this_ *RenameProperties_Property, name string) {
	_ = name
	this_.initFields_RenameProperties_Property()
	this_.oldName = name
}
func (*RenameProperties_ProcessProperties) AlwaysShouldTraverse_()      {}
func (*RenameProperties_ProcessProperties) CanReleaseSyntacticScopes_() {}
