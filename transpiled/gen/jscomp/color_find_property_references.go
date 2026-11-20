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
// Port of com/google/javascript/jscomp/disambiguate/ColorFindPropertyReferences.java (generated mechanically by j2go).

package jscomp

import (
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences
type ColorFindPropertyReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	propIndex             j2gort.OrderedMap[string, *PropertyClustering]  `gson:"propIndex"`
	colorGraphNodeFactory *ColorGraphNodeFactory                          `gson:"colorGraphNodeFactory"`
	isPropertyReflector   ColorFindPropertyReferences_IsPropertyReflector `gson:"isPropertyReflector"`
}

func (this_ *ColorFindPropertyReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorFindPropertyReferences) IsJavaNil_() bool { return this_ == nil }
func (*ColorFindPropertyReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences", "ColorFindPropertyReferences"
}
func (this_ *ColorFindPropertyReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ColorFindPropertyReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences", "com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorFindPropertyReferences](v_)
	})
}
func (*ColorFindPropertyReferences) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ColorFindPropertyReferences) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ColorFindPropertyReferences) initFields_ColorFindPropertyReferences() {
	_ = this_
	this_.propIndex = j2gort.NewOrderedMap[string, *PropertyClustering]()
}

// Java constructor: ColorFindPropertyReferences(com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory colorGraphNodeFactory, com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences.IsPropertyReflector isPropertyReflector)
func NewColorFindPropertyReferences(colorGraphNodeFactory *ColorGraphNodeFactory, isPropertyReflector ColorFindPropertyReferences_IsPropertyReflector) *ColorFindPropertyReferences {
	this_ := &ColorFindPropertyReferences{}
	InitColorFindPropertyReferences(this_, colorGraphNodeFactory, isPropertyReflector)
	return this_
}

func InitColorFindPropertyReferences(this_ *ColorFindPropertyReferences, colorGraphNodeFactory *ColorGraphNodeFactory, isPropertyReflector ColorFindPropertyReferences_IsPropertyReflector) {
	_ = colorGraphNodeFactory
	_ = isPropertyReflector
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ColorFindPropertyReferences()
	this_.colorGraphNodeFactory = colorGraphNodeFactory
	this_.isPropertyReflector = isPropertyReflector
}

// Java method: ColorFindPropertyReferences.getOrCreateProperty(java.lang.String name)
func (this_ *ColorFindPropertyReferences) GetOrCreateProperty(name string) *PropertyClustering {
	_ = name
	var prop *PropertyClustering = this_.propIndex.Get(name)
	_ = prop
	if prop == nil {
		prop = NewPropertyClustering(name)
		this_.propIndex.Put(name, prop)
	}
	return prop
}

// Java method: ColorFindPropertyReferences.getProperty(java.lang.String name)
func (this_ *ColorFindPropertyReferences) GetProperty(name string) *PropertyClustering {
	_ = name
	return this_.propIndex.Get(name)
}

// Java method: ColorFindPropertyReferences.collectProperties()
func (this_ *ColorFindPropertyReferences) CollectProperties() j2gort.Ext_java_util_ArrayList {
	var list []*PropertyClustering = j2gort.MakeNonNilSlice[*PropertyClustering](this_.propIndex.Size())
	_ = list
	{
		cur_1 := this_.propIndex.IterValues()
		for cur_1.HasNext() {

			prop := cur_1.
				Next()
			_ = prop
			list = append(list, prop)
		}
	}
	return j2gort.WrapTypedArrayList(list)
}

// Java method: ColorFindPropertyReferences.clearProperties()
func (this_ *ColorFindPropertyReferences) ClearProperties() {
	this_.propIndex.Clear()
}

// Java method: ColorFindPropertyReferences.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ColorFindPropertyReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_ColorFindPropertyReferences_registerPropertyUse_6256(t, n, n.GetFirstChild().GetColor())
	case int32(47) /* Token.OBJECTLIT */ :
		this_.priv_ColorFindPropertyReferences_handleObjectLit_6248(t, n)
	case int32(28) /* Token.CALL */ :
		this_.priv_ColorFindPropertyReferences_handleCall_6249(t, n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_ColorFindPropertyReferences_handleClass_6250(t, n)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_ColorFindPropertyReferences_handleObjectPattern_6252(t, n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_ColorFindPropertyReferences_handleFunction_6251(n)
	default:
	}
}

// Java method: ColorFindPropertyReferences.handleObjectLit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handleObjectLit_6248(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if n.GetParent().IsCall() && NodeUtil_IsObjectDefinePropertiesDefinition(n.GetParent()) {
		return
	}
	var owner *pkg_colors.Color = n.GetColor()
	_ = owner
	this_.priv_ColorFindPropertyReferences_traverseObjectlitLike_6255(t, n, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		m := j2gort.CastNode[pkg_rhino.Node](lp_3)
		_ = m
		return j2gort.PtrToAny(owner)
	}))
}

// Java method: ColorFindPropertyReferences.handleCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handleCall_6249(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var target pkg_rhino.Node = call.GetFirstChild()
	_ = target
	if !(target.IsQualifiedName()) {
		return
	}
	if this_.isPropertyReflector.Test(target) {
		this_.priv_ColorFindPropertyReferences_handlePropertyReflectorCall_6253(t, call)
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(call) {
		this_.priv_ColorFindPropertyReferences_handleObjectDefineProperties_6254(t, call)
	}
}

// Java method: ColorFindPropertyReferences.handleClass(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handleClass_6250(t *NodeTraversal, classNode pkg_rhino.Node) {
	_ = t
	_ = classNode
	var classType *pkg_colors.Color = classNode.GetColor()
	_ = classType
	this_.priv_ColorFindPropertyReferences_traverseObjectlitLike_6255(t, NodeUtil_GetClassMembers(classNode), j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		m := j2gort.CastNode[pkg_rhino.Node](lp_4)
		_ = m
		if m.IsStaticMember() {
			return j2gort.PtrToAny(classType)
		} else if m.IsMemberFieldDef() {
			var classInstanceType []*pkg_colors.Color = classType.GetInstanceColorsArray()
			_ = classInstanceType
			if int32(len(classInstanceType)) == int32(0) {
				return j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN())
			} else {
				return j2gort.PtrToAny(pkg_colors.Color_CreateUnion_1_Color(classInstanceType))
			}
		} else {
			if !((m.IsMemberFunctionDef() || m.IsGetterDef()) || m.IsSetterDef()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(m))
			}
			var classPrototypeType []*pkg_colors.Color = classType.GetPrototypesArray()
			_ = classPrototypeType
			if int32(len(classPrototypeType)) == int32(0) {
				return j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN())
			} else {
				return j2gort.PtrToAny(pkg_colors.Color_CreateUnion_1_Color(classPrototypeType))
			}
		}
	}))
	this_.colorGraphNodeFactory.CreateNode(classType)
}

// Java method: ColorFindPropertyReferences.handleFunction(com.google.javascript.rhino.Node fnNode)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handleFunction_6251(fnNode pkg_rhino.Node) {
	_ = fnNode
	var fnType *pkg_colors.Color = fnNode.GetColor()
	_ = fnType
	if (fnType != nil) && fnType.HasInstanceColors() {
		this_.colorGraphNodeFactory.CreateNode(fnType)
	}
}

// Java method: ColorFindPropertyReferences.handleObjectPattern(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node pattern)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handleObjectPattern_6252(t *NodeTraversal, pattern pkg_rhino.Node) {
	_ = t
	_ = pattern
	var owner *pkg_colors.Color = pattern.GetColor()
	_ = owner
	this_.priv_ColorFindPropertyReferences_traverseObjectlitLike_6255(t, pattern, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		m := j2gort.CastNode[pkg_rhino.Node](lp_5)
		_ = m
		return j2gort.PtrToAny(owner)
	}))
}

// Java method: ColorFindPropertyReferences.handlePropertyReflectorCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handlePropertyReflectorCall_6253(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var name pkg_rhino.Node = call.GetSecondChild()
	_ = name
	if (name == 0) || !(name.IsStringLit()) {
		return
	}
	var obj pkg_rhino.Node = name.GetNext()
	_ = obj
	var objColor *pkg_colors.Color
	_ = objColor
	if obj == 0 {
		objColor = nil
	} else {
		objColor = obj.GetColor()
	}
	this_.priv_ColorFindPropertyReferences_registerPropertyUse_6256(t, name, objColor)
}

// Java method: ColorFindPropertyReferences.handleObjectDefineProperties(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_handleObjectDefineProperties_6254(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	var typeObj pkg_rhino.Node = call.GetSecondChild()
	_ = typeObj
	var objectLiteral pkg_rhino.Node = typeObj.GetNext()
	_ = objectLiteral
	if !(objectLiteral.IsObjectLit()) {
		return
	}
	var type_ *pkg_colors.Color = typeObj.GetColor()
	_ = type_
	this_.priv_ColorFindPropertyReferences_traverseObjectlitLike_6255(t, objectLiteral, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_6 any) any {
		m := j2gort.CastNode[pkg_rhino.Node](lp_6)
		_ = m
		return j2gort.PtrToAny(type_)
	}))
}

// Java method: ColorFindPropertyReferences.traverseObjectlitLike(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, java.util.function.Function<com.google.javascript.rhino.Node,com.google.javascript.jscomp.colors.Color> memberOwnerFn)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_traverseObjectlitLike_6255(t *NodeTraversal, n pkg_rhino.Node, memberOwnerFn j2gort.Ext_java_util_function_Function) {
	_ = t
	_ = n
	_ = memberOwnerFn
	if !((n.IsObjectLit() || n.IsObjectPattern()) || n.IsClassMembers()) {
		j2gort.CheckState(false)
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			switch child.TokenSwitchOrd_() {
			case int32(90) /* Token.BLOCK */, int32(129) /* Token.COMPUTED_PROP */, int32(109) /* Token.COMPUTED_FIELD_DEF */, int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */ :
				continue
			case int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
				if child.IsQuotedStringKey() {
					continue
				}
				this_.priv_ColorFindPropertyReferences_registerPropertyUse_6256(t, child, j2gort.CastPtr[pkg_colors.Color](memberOwnerFn.Apply(j2gort.BoxNode(child))))
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected child of " + j2gort.StrVal(n.GetToken()) + ": " + j2gort.StrValString(child.ToStringTree()))))))
			}
		}
	}
}

// Java method: ColorFindPropertyReferences.registerPropertyUse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node site, com.google.javascript.jscomp.colors.Color owner)
func (this_ *ColorFindPropertyReferences) priv_ColorFindPropertyReferences_registerPropertyUse_6256(t *NodeTraversal, site pkg_rhino.Node, owner *pkg_colors.Color) {
	_ = t
	_ = site
	_ = owner
	var prop *PropertyClustering = this_.GetOrCreateProperty(site.GetString())
	_ = prop
	var flatOwner *ColorGraphNode = this_.colorGraphNodeFactory.CreateNode(owner)
	_ = flatOwner
	flatOwner.GetAssociatedProps().Put(prop, StaticGet_ColorGraphNode_PropAssociation_AST())
	prop.GetClusters().Add(flatOwner)
	prop.PutUseSite(site, flatOwner)
	if t.GetCurrentScript().IsFromExterns() || ((owner != nil) && owner.GetPropertiesKeepOriginalName()) {
		prop.RegisterOriginalNameType(flatOwner)
	}
}

// Java interface: com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences.IsPropertyReflector
type ColorFindPropertyReferences_IsPropertyReflector interface {
	Implements_jscomp_ColorFindPropertyReferences_IsPropertyReflector()
	IsJavaNil_() bool
	Test(pkg_rhino.Node) bool
}

func ToIface_ColorFindPropertyReferences_IsPropertyReflector[T interface {
	ColorFindPropertyReferences_IsPropertyReflector
	comparable
}](v T) ColorFindPropertyReferences_IsPropertyReflector {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ColorFindPropertyReferences_IsPropertyReflector[T ColorFindPropertyReferences_IsPropertyReflector](v T) ColorFindPropertyReferences_IsPropertyReflector {
	var vi ColorFindPropertyReferences_IsPropertyReflector = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences.IsPropertyReflector", "com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences$IsPropertyReflector", func(v_ any) bool {
		return j2gort.InstanceOf[ColorFindPropertyReferences_IsPropertyReflector](v_)
	})
}

type lambda_ColorFindPropertyReferences_IsPropertyReflector struct {
	fn func(pkg_rhino.Node) bool
}

func (l *lambda_ColorFindPropertyReferences_IsPropertyReflector) IsJavaNil_() bool { return l == nil }
func NewLambda_ColorFindPropertyReferences_IsPropertyReflector(fn func(pkg_rhino.Node) bool) ColorFindPropertyReferences_IsPropertyReflector {
	if fn == nil {
		return nil
	}
	return &lambda_ColorFindPropertyReferences_IsPropertyReflector{fn: fn}
}

func (*lambda_ColorFindPropertyReferences_IsPropertyReflector) Implements_jscomp_ColorFindPropertyReferences_IsPropertyReflector() {
}
func (l *lambda_ColorFindPropertyReferences_IsPropertyReflector) Test(p0_ pkg_rhino.Node) bool {
	return l.fn(p0_)
}
func (*ColorFindPropertyReferences) AlwaysShouldTraverse_()      {}
func (*ColorFindPropertyReferences) CanReleaseSyntacticScopes_() {}
