/*
 * Copyright 2012 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InlineProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.InlineProperties
type InlineProperties struct {
	compiler AbstractCompiler                                          `gson:"compiler"`
	registry *pkg_colors.ColorRegistry                                 `gson:"registry"`
	props    j2gort.OrderedMap[string, *InlineProperties_PropertyInfo] `gson:"props"`
	idHash_  atomic.Int32
}

func (this_ *InlineProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineProperties) IsJavaNil_() bool { return this_ == nil }
func (*InlineProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineProperties", "InlineProperties"
}
func (this_ *InlineProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineProperties", "com.google.javascript.jscomp.InlineProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineProperties](v_)
	})
}
func (*InlineProperties) Implements_jscomp_CompilerPass() {}

var static_InlineProperties_INVALIDATED *InlineProperties_PropertyInfo

func StaticGet_InlineProperties_INVALIDATED() *InlineProperties_PropertyInfo {
	clinit_InlineProperties()
	return static_InlineProperties_INVALIDATED
}
func StaticSet_InlineProperties_INVALIDATED(v *InlineProperties_PropertyInfo) {
	clinit_InlineProperties()
	static_InlineProperties_INVALIDATED = v
}

var clinitDone_InlineProperties uint32
var clinitState_InlineProperties j2gort.ClinitState

func clinit_InlineProperties() {
	if atomic.LoadUint32(&clinitDone_InlineProperties) != 0 {
		return
	}
	clinitSlow_InlineProperties()
}

func clinitSlow_InlineProperties() {
	if !j2gort.EnterClinit(&clinitState_InlineProperties, &clinitDone_InlineProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InlineProperties, &clinitDone_InlineProperties)
	static_InlineProperties_INVALIDATED = NewInlineProperties_PropertyInfo(nil, pkg_rhino.Node(0))
}

func (this_ *InlineProperties) initFields_InlineProperties() {
	_ = this_
	this_.props = j2gort.NewOrderedMap[string, *InlineProperties_PropertyInfo]()
}

// Java constructor: InlineProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewInlineProperties(compiler AbstractCompiler) *InlineProperties {
	clinit_InlineProperties()
	this_ := &InlineProperties{}
	InitInlineProperties(this_, compiler)
	return this_
}

func InitInlineProperties(this_ *InlineProperties, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_InlineProperties()
	this_.compiler = compiler
	this_.registry = compiler.GetColorRegistry()
	this_.priv_InlineProperties_invalidateExternProperties_3358()
}

// Java method: InlineProperties.invalidateExternProperties()
func (this_ *InlineProperties) priv_InlineProperties_invalidateExternProperties_3358() {
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetExternProperties()))
		for cur_1.HasNext() {
			name := cur_1.NextStr()
			_ = name
			this_.props.Put(name, StaticGet_InlineProperties_INVALIDATED())
		}
	}
}

// Java method: InlineProperties.getColor(com.google.javascript.rhino.Node n)
func (this_ *InlineProperties) priv_InlineProperties_getColor_3359(n pkg_rhino.Node) *pkg_colors.Color {
	_ = n
	var color *pkg_colors.Color = n.GetColor()
	_ = color
	if color == nil {
		return pkg_colors.StaticGet_StandardColors_UNKNOWN()
	} else {
		return color
	}
}

// Java method: InlineProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewInlineProperties_GatherCandidates(this_))
	NodeTraversal_Traverse(this_.compiler, root, NewInlineProperties_ReplaceCandidates(this_))
}

// Java static method: InlineProperties.removeNullAndUndefinedIfUnion(com.google.javascript.jscomp.colors.Color original)
func InlineProperties_RemoveNullAndUndefinedIfUnion(original *pkg_colors.Color) *pkg_colors.Color {
	clinit_InlineProperties()
	_ = original
	if original.IsUnion() {
		return original.SubtractNullOrVoid()
	} else {
		return original
	}
}

// Java class: com.google.javascript.jscomp.InlineProperties.PropertyInfo
type InlineProperties_PropertyInfo struct {
	color   *pkg_colors.Color `gson:"color"`
	value   pkg_rhino.Node    `gson:"value"`
	idHash_ atomic.Int32
}

func (this_ *InlineProperties_PropertyInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineProperties_PropertyInfo) IsJavaNil_() bool { return this_ == nil }
func (*InlineProperties_PropertyInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineProperties$PropertyInfo", "PropertyInfo"
}
func (this_ *InlineProperties_PropertyInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineProperties_PropertyInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineProperties.PropertyInfo", "com.google.javascript.jscomp.InlineProperties$PropertyInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineProperties_PropertyInfo](v_)
	})
}

func (this_ *InlineProperties_PropertyInfo) initFields_InlineProperties_PropertyInfo() {
	_ = this_
}

// Java constructor: PropertyInfo(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color color, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node value)
func NewInlineProperties_PropertyInfo(color *pkg_colors.Color, value pkg_rhino.Node) *InlineProperties_PropertyInfo {
	this_ := &InlineProperties_PropertyInfo{}
	InitInlineProperties_PropertyInfo(this_, color, value)
	return this_
}

func InitInlineProperties_PropertyInfo(this_ *InlineProperties_PropertyInfo, color *pkg_colors.Color, value pkg_rhino.Node) {
	_ = color
	_ = value
	this_.initFields_InlineProperties_PropertyInfo()
	this_.color = color
	this_.value = value
}

// Java class: com.google.javascript.jscomp.InlineProperties.GatherCandidates
type InlineProperties_GatherCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *InlineProperties
}

func (this_ *InlineProperties_GatherCandidates) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineProperties_GatherCandidates) IsJavaNil_() bool { return this_ == nil }
func (*InlineProperties_GatherCandidates) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineProperties$GatherCandidates", "GatherCandidates"
}
func (this_ *InlineProperties_GatherCandidates) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InlineProperties_GatherCandidates) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineProperties.GatherCandidates", "com.google.javascript.jscomp.InlineProperties$GatherCandidates", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineProperties_GatherCandidates](v_)
	})
}
func (this_ *InlineProperties_GatherCandidates) SetOuter_(o *InlineProperties) { this_.outer_ = o }
func (this_ *InlineProperties_GatherCandidates) GetOuter_() *InlineProperties  { return this_.outer_ }
func (*InlineProperties_GatherCandidates) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineProperties_GatherCandidates) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InlineProperties_GatherCandidates) initFields_InlineProperties_GatherCandidates() {
	_ = this_
}

// Java constructor: GatherCandidates()
func NewInlineProperties_GatherCandidates(outer_ *InlineProperties) *InlineProperties_GatherCandidates {
	this_ := &InlineProperties_GatherCandidates{}
	this_.outer_ = outer_
	InitInlineProperties_GatherCandidates(this_)
	return this_
}

func InitInlineProperties_GatherCandidates(this_ *InlineProperties_GatherCandidates) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineProperties_GatherCandidates()
}

// Java method: GatherCandidates.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineProperties_GatherCandidates) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var invalidatingPropRef bool = false
	_ = invalidatingPropRef
	var propName string = j2gort.NullString
	_ = propName
	if n.IsGetProp() {
		propName = n.GetString()
		if parent.IsAssign() {
			invalidatingPropRef = !(this_.priv_GatherCandidates_maybeRecordCandidateGetpropDefinition_3361(t, n, parent))
		} else if NodeUtil_IsLValue(n) {
			invalidatingPropRef = true
		} else if parent.IsDelProp() {
			invalidatingPropRef = true
		} else {
			invalidatingPropRef = false
		}
	} else if (((n.IsStringKey() && !(n.GetParent().IsObjectPattern())) || n.IsGetterDef()) || n.IsSetterDef()) || n.IsMemberFunctionDef() {
		propName = n.GetString()
		invalidatingPropRef = true
	} else if n.IsMemberFieldDef() {
		propName = n.GetString()
		if n.HasChildren() {
			invalidatingPropRef = !(this_.priv_GatherCandidates_maybeRecordCandidateClassFieldDefinition_3360(n))
		} else {
			invalidatingPropRef = true
		}
	} else if ((n.IsCall() && this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(n.GetFirstChild())) && n.HasMoreThanOneChild()) && n.GetSecondChild().IsStringLit() {
		propName = n.GetSecondChild().GetString()
		invalidatingPropRef = true
	} else {
		return
	}
	if invalidatingPropRef {
		j2gort.CheckNotNullString(propName)
		this_.priv_GatherCandidates_invalidateProperty_3362(propName)
	}
}

// Java method: GatherCandidates.maybeRecordCandidateClassFieldDefinition(com.google.javascript.rhino.Node n)
func (this_ *InlineProperties_GatherCandidates) priv_GatherCandidates_maybeRecordCandidateClassFieldDefinition_3360(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsMemberFieldDef()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var value pkg_rhino.Node = n.GetFirstChild()
	_ = value
	var propName string = n.GetString()
	_ = propName
	var classNode pkg_rhino.Node = n.GetGrandparent()
	_ = classNode
	var c *pkg_colors.Color = nil
	_ = c
	if n.IsStaticMember() {
		c = this_.outer_.priv_InlineProperties_getColor_3359(classNode)
	} else {
		var possibleInstances []*pkg_colors.Color = this_.outer_.priv_InlineProperties_getColor_3359(classNode).GetInstanceColorsArray()
		_ = possibleInstances
		if int32(len(possibleInstances)) == int32(0) {
			c = pkg_colors.StaticGet_StandardColors_UNKNOWN()
		} else {
			c = pkg_colors.Color_CreateUnion_1_Color(possibleInstances)
		}
	}
	return this_.priv_GatherCandidates_maybeStoreCandidateValue_3363(c, propName, value)
}

// Java method: GatherCandidates.maybeRecordCandidateGetpropDefinition(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineProperties_GatherCandidates) priv_GatherCandidates_maybeRecordCandidateGetpropDefinition_3361(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if !(n.IsGetProp() && parent.IsAssign()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var src pkg_rhino.Node = n.GetFirstChild()
	_ = src
	var propName string = n.GetString()
	_ = propName
	var value pkg_rhino.Node = parent.GetLastChild()
	_ = value
	if src.IsThis() {
		if this_.priv_GatherCandidates_inConstructor_3364(t) {
			return this_.priv_GatherCandidates_maybeStoreCandidateValue_3363(this_.outer_.priv_InlineProperties_getColor_3359(src), propName, value)
		}
		return false
	} else if (t.InGlobalHoistScope() && src.IsGetProp()) && (j2gort.CheckNotNullString(src.GetString()) == "prototype") {
		var instanceType *pkg_colors.Color = this_.outer_.priv_InlineProperties_getColor_3359(src)
		_ = instanceType
		return this_.priv_GatherCandidates_maybeStoreCandidateValue_3363(instanceType, propName, value)
	} else if t.InGlobalHoistScope() {
		var targetType *pkg_colors.Color = this_.outer_.priv_InlineProperties_getColor_3359(src)
		_ = targetType
		if (targetType != nil) && targetType.IsConstructor() {
			return this_.priv_GatherCandidates_maybeStoreCandidateValue_3363(targetType, propName, value)
		}
	}
	return false
}

// Java method: GatherCandidates.invalidateProperty(java.lang.String propName)
func (this_ *InlineProperties_GatherCandidates) priv_GatherCandidates_invalidateProperty_3362(propName string) {
	_ = propName
	this_.outer_.props.Put(propName, StaticGet_InlineProperties_INVALIDATED())
}

// Java method: GatherCandidates.maybeStoreCandidateValue(com.google.javascript.jscomp.colors.Color color, java.lang.String propName, com.google.javascript.rhino.Node value)
func (this_ *InlineProperties_GatherCandidates) priv_GatherCandidates_maybeStoreCandidateValue_3363(color *pkg_colors.Color, propName string, value pkg_rhino.Node) bool {
	_ = color
	_ = propName
	_ = value
	j2gort.CheckNotNullNode(value)
	if ((!(this_.outer_.props.ContainsKey(propName)) && !(color.IsInvalidating())) && NodeUtil_IsImmutableValue(value)) && NodeUtil_IsExecutedExactlyOnce(value) {
		this_.outer_.props.Put(propName, NewInlineProperties_PropertyInfo(color, value))
		return true
	}
	return false
}

// Java method: GatherCandidates.inConstructor(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *InlineProperties_GatherCandidates) priv_GatherCandidates_inConstructor_3364(t *NodeTraversal) bool {
	_ = t
	var root pkg_rhino.Node = t.GetEnclosingFunction()
	_ = root
	if root == 0 {
		return false
	}
	return NodeUtil_IsConstructor(root)
}

// Java class: com.google.javascript.jscomp.InlineProperties.ReplaceCandidates
type InlineProperties_ReplaceCandidates struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                     *InlineProperties
	hasInSupertypesListSeenSet j2gort.CustomOrderedSet[*pkg_colors.Color] `gson:"hasInSupertypesListSeenSet"`
}

func (this_ *InlineProperties_ReplaceCandidates) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineProperties_ReplaceCandidates) IsJavaNil_() bool { return this_ == nil }
func (*InlineProperties_ReplaceCandidates) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineProperties$ReplaceCandidates", "ReplaceCandidates"
}
func (this_ *InlineProperties_ReplaceCandidates) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InlineProperties_ReplaceCandidates) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineProperties.ReplaceCandidates", "com.google.javascript.jscomp.InlineProperties$ReplaceCandidates", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineProperties_ReplaceCandidates](v_)
	})
}
func (this_ *InlineProperties_ReplaceCandidates) SetOuter_(o *InlineProperties) { this_.outer_ = o }
func (this_ *InlineProperties_ReplaceCandidates) GetOuter_() *InlineProperties  { return this_.outer_ }
func (*InlineProperties_ReplaceCandidates) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineProperties_ReplaceCandidates) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InlineProperties_ReplaceCandidates) initFields_InlineProperties_ReplaceCandidates() {
	_ = this_
	this_.hasInSupertypesListSeenSet = j2gort.NewCustomOrderedSet[*pkg_colors.Color]()
}

// Java constructor: ReplaceCandidates()
func NewInlineProperties_ReplaceCandidates(outer_ *InlineProperties) *InlineProperties_ReplaceCandidates {
	this_ := &InlineProperties_ReplaceCandidates{}
	this_.outer_ = outer_
	InitInlineProperties_ReplaceCandidates(this_)
	return this_
}

func InitInlineProperties_ReplaceCandidates(this_ *InlineProperties_ReplaceCandidates) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineProperties_ReplaceCandidates()
}

// Java method: ReplaceCandidates.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineProperties_ReplaceCandidates) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsGetProp() && !(NodeUtil_IsLValue(n)) {
		var target pkg_rhino.Node = n.GetFirstChild()
		_ = target
		var propName string = n.GetString()
		_ = propName
		var info *InlineProperties_PropertyInfo = this_.outer_.props.Get(propName)
		_ = info
		if ((info != nil) && (info != StaticGet_InlineProperties_INVALIDATED())) && this_.priv_ReplaceCandidates_isMatchingType_3365(target, info.color) {
			var replacement pkg_rhino.Node = info.value.CloneTree_0()
			_ = replacement
			if this_.outer_.compiler.GetAstAnalyzer().MayHaveSideEffects(n.GetFirstChild()) {
				replacement = pkg_rhino.IR_Comma(n.RemoveFirstChild(), replacement).Srcref(n)
			}
			n.ReplaceWith(replacement)
			this_.outer_.compiler.ReportChangeToEnclosingScope(replacement)
		}
	}
}

// Java method: ReplaceCandidates.isMatchingType(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.colors.Color src)
func (this_ *InlineProperties_ReplaceCandidates) priv_ReplaceCandidates_isMatchingType_3365(n pkg_rhino.Node, src *pkg_colors.Color) bool {
	_ = n
	_ = src
	src = InlineProperties_RemoveNullAndUndefinedIfUnion(src)
	var dest *pkg_colors.Color = InlineProperties_RemoveNullAndUndefinedIfUnion(this_.outer_.priv_InlineProperties_getColor_3359(n))
	_ = dest
	if dest.IsInvalidating() {
		return false
	}
	if dest.IsUnion() || src.IsUnion() {
		return false
	}
	return this_.priv_ReplaceCandidates_hasInSupertypesList_3366(dest, src)
}

// Java method: ReplaceCandidates.hasInSupertypesList(com.google.javascript.jscomp.colors.Color subCtor, com.google.javascript.jscomp.colors.Color superCtor)
func (this_ *InlineProperties_ReplaceCandidates) priv_ReplaceCandidates_hasInSupertypesList_3366(subCtor *pkg_colors.Color, superCtor *pkg_colors.Color) bool {
	_ = subCtor
	_ = superCtor
	defer func() {
		this_.hasInSupertypesListSeenSet.Remove(subCtor)
	}()
	{
		if !(this_.hasInSupertypesListSeenSet.Add(subCtor)) {
			return false
		}
		if (subCtor == nil) || (superCtor == nil) {
			return false
		}
		if subCtor.Equals(j2gort.PtrToAny(superCtor)) {
			return true
		}
		var supertypes []*pkg_colors.Color = this_.outer_.registry.GetDisambiguationSupertypesArray(subCtor)
		_ = supertypes
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = int32(len(supertypes))
			_ = n
			for ; i < n; i = i + 1 {
				var immediateSupertype *pkg_colors.Color = supertypes[i]
				_ = immediateSupertype
				if !(immediateSupertype.IsUnion()) && this_.priv_ReplaceCandidates_hasInSupertypesList_3366(immediateSupertype, superCtor) {
					return true
				}
			}
		}
		return false
	}
	panic("unreachable")
}
func (*InlineProperties_GatherCandidates) AlwaysShouldTraverse_()       {}
func (*InlineProperties_GatherCandidates) CanReleaseSyntacticScopes_()  {}
func (*InlineProperties_ReplaceCandidates) AlwaysShouldTraverse_()      {}
func (*InlineProperties_ReplaceCandidates) CanReleaseSyntacticScopes_() {}
