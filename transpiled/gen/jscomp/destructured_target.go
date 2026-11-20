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
// Port of com/google/javascript/jscomp/DestructuredTarget.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.DestructuredTarget
type DestructuredTarget struct {
	registry            *pkg_rhino.JSTypeRegistry                  `gson:"registry"`
	patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier `gson:"patternTypeSupplier"`
	objectPatternKey    pkg_rhino.Node                             `gson:"objectPatternKey"`
	node                pkg_rhino.Node                             `gson:"node"`
	defaultValue        pkg_rhino.Node                             `gson:"defaultValue"`
	pattern             pkg_rhino.Node                             `gson:"pattern"`
	idHash_             atomic.Int32
	isRest              bool `gson:"isRest"`
}

func (this_ *DestructuredTarget) BindSelf_(s any) {
	_ = s
}
func (this_ *DestructuredTarget) IsJavaNil_() bool { return this_ == nil }
func (*DestructuredTarget) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DestructuredTarget", "DestructuredTarget"
}
func (this_ *DestructuredTarget) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DestructuredTarget) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DestructuredTarget", "com.google.javascript.jscomp.DestructuredTarget", func(v_ any) bool {
		return j2gort.InstanceOf[*DestructuredTarget](v_)
	})
}

func (this_ *DestructuredTarget) initFields_DestructuredTarget() {
	_ = this_
}

// Java constructor: DestructuredTarget(com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.Node node, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node defaultValue, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node objectPatternKey, com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.common.base.Supplier<com.google.javascript.rhino.jstype.JSType> patternTypeSupplier, boolean isRest)
func NewDestructuredTarget(pattern pkg_rhino.Node, node pkg_rhino.Node, defaultValue pkg_rhino.Node, objectPatternKey pkg_rhino.Node, registry *pkg_rhino.JSTypeRegistry, patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier, isRest bool) *DestructuredTarget {
	this_ := &DestructuredTarget{}
	InitDestructuredTarget(this_, pattern, node, defaultValue, objectPatternKey, registry, patternTypeSupplier, isRest)
	return this_
}

func InitDestructuredTarget(this_ *DestructuredTarget, pattern pkg_rhino.Node, node pkg_rhino.Node, defaultValue pkg_rhino.Node, objectPatternKey pkg_rhino.Node, registry *pkg_rhino.JSTypeRegistry, patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier, isRest bool) {
	_ = pattern
	_ = node
	_ = defaultValue
	_ = objectPatternKey
	_ = registry
	_ = patternTypeSupplier
	_ = isRest
	this_.initFields_DestructuredTarget()
	this_.pattern = pattern
	this_.node = node
	this_.objectPatternKey = objectPatternKey
	this_.registry = registry
	this_.patternTypeSupplier = patternTypeSupplier
	this_.isRest = isRest
	this_.defaultValue = defaultValue
}

// Java method: DestructuredTarget.getComputedProperty()
func (this_ *DestructuredTarget) GetComputedProperty() pkg_rhino.Node {
	if this_.HasComputedProperty() {
		return this_.objectPatternKey
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: DestructuredTarget.hasComputedProperty()
func (this_ *DestructuredTarget) HasComputedProperty() bool {
	return ((this_.objectPatternKey != 0) && this_.objectPatternKey.IsComputedProp())
}

// Java method: DestructuredTarget.hasStringKey()
func (this_ *DestructuredTarget) HasStringKey() bool {
	return ((this_.objectPatternKey != 0) && this_.objectPatternKey.IsStringKey())
}

// Java method: DestructuredTarget.getStringKey()
func (this_ *DestructuredTarget) GetStringKey() pkg_rhino.Node {
	if this_.HasStringKey() {
		return this_.objectPatternKey
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: DestructuredTarget.getDefaultValue()
func (this_ *DestructuredTarget) GetDefaultValue() pkg_rhino.Node {
	return this_.defaultValue
}

// Java method: DestructuredTarget.hasDefaultValue()
func (this_ *DestructuredTarget) HasDefaultValue() bool {
	return (this_.defaultValue != 0)
}

// Java method: DestructuredTarget.getNode()
func (this_ *DestructuredTarget) GetNode() pkg_rhino.Node {
	return this_.node
}

// Java static method: DestructuredTarget.createTarget(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType destructuringPatternType, com.google.javascript.rhino.Node destructuringChild)
func DestructuredTarget_CreateTarget_3_JSTypeRegistry_JSType_Node(registry *pkg_rhino.JSTypeRegistry, destructuringPatternType pkg_rhino.JSType, destructuringChild pkg_rhino.Node) *DestructuredTarget {
	_ = registry
	_ = destructuringPatternType
	_ = destructuringChild
	return DestructuredTarget_CreateTarget_3_JSTypeRegistry_Ext_com_google_common_base_Supplier_Node(registry, j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return destructuringPatternType
	}), destructuringChild)
}

// Java static method: DestructuredTarget.createTarget(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.common.base.Supplier<com.google.javascript.rhino.jstype.JSType> destructuringPatternType, com.google.javascript.rhino.Node destructuringChild)
func DestructuredTarget_CreateTarget_3_JSTypeRegistry_Ext_com_google_common_base_Supplier_Node(registry *pkg_rhino.JSTypeRegistry, destructuringPatternType j2gort.Ext_com_google_common_base_Supplier, destructuringChild pkg_rhino.Node) *DestructuredTarget {
	_ = registry
	_ = destructuringPatternType
	_ = destructuringChild
	if !(destructuringChild.GetParent().IsDestructuringPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(destructuringChild))
	}
	var val_1 DestructuredTarget_Builder
	InitDestructuredTarget_Builder(&val_1, registry, destructuringChild.GetParent(), destructuringPatternType)
	var builder *DestructuredTarget_Builder = &val_1
	_ = builder
	var value pkg_rhino.Node = pkg_rhino.Node(0)
	_ = value
	switch destructuringChild.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */ :
		builder.SetObjectPatternKey(destructuringChild)
		value = destructuringChild.GetFirstChild()
		if value.IsDefaultValue() {
			builder.SetNode(value.GetFirstChild())
			builder.SetDefaultValue(value.GetSecondChild())
		} else {
			builder.SetNode(value)
		}
	case int32(129) /* Token.COMPUTED_PROP */ :
		builder.SetObjectPatternKey(destructuringChild)
		value = destructuringChild.GetSecondChild()
		if value.IsDefaultValue() {
			builder.SetNode(value.GetFirstChild())
			builder.SetDefaultValue(value.GetSecondChild())
		} else {
			builder.SetNode(value)
		}
	case int32(103) /* Token.OBJECT_PATTERN */, int32(102) /* Token.ARRAY_PATTERN */, int32(32) /* Token.NAME */, int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
		builder.SetNode(destructuringChild)
	case int32(134) /* Token.DEFAULT_VALUE */ :
		builder.SetNode(destructuringChild.GetFirstChild())
		builder.SetDefaultValue(destructuringChild.GetSecondChild())
	case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
		builder.SetNode(destructuringChild.GetFirstChild())
		builder.SetIsRest(true)
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected child of destructuring pattern "+j2gort.StrVal(destructuringChild))))
	}
	return builder.Build()
}

// Java static method: DestructuredTarget.createAllNonEmptyTargetsInPattern(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType patternType, com.google.javascript.rhino.Node pattern)
func DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_JSType_Node(registry *pkg_rhino.JSTypeRegistry, patternType pkg_rhino.JSType, pattern pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = registry
	_ = patternType
	_ = pattern
	return DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_Ext_com_google_common_base_Supplier_Node(registry, j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return patternType
	}), pattern)
}

// Java static method: DestructuredTarget.createAllNonEmptyTargetsInPattern(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.common.base.Supplier<com.google.javascript.rhino.jstype.JSType> patternType, com.google.javascript.rhino.Node pattern)
func DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_Ext_com_google_common_base_Supplier_Node(registry *pkg_rhino.JSTypeRegistry, patternType j2gort.Ext_com_google_common_base_Supplier, pattern pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = registry
	_ = patternType
	_ = pattern
	if !(pattern.IsDestructuringPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(pattern))
	}
	var builder []*DestructuredTarget = []*DestructuredTarget{}
	_ = builder
	{
		var child pkg_rhino.Node = pattern.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsEmpty() {
				continue
			}
			builder = append(builder, j2gort.CheckNotNullElem(DestructuredTarget_CreateTarget_3_JSTypeRegistry_Ext_com_google_common_base_Supplier_Node(registry, patternType, child)))
		}
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java method: DestructuredTarget.inferTypeWithoutUsingDefaultValue()
func (this_ *DestructuredTarget) InferTypeWithoutUsingDefaultValue() pkg_rhino.JSType {
	if this_.pattern.IsObjectPattern() {
		return this_.priv_DestructuredTarget_inferObjectPatternKeyType_2428()
	} else {
		return this_.priv_DestructuredTarget_inferArrayPatternTargetType_2429()
	}
}

// Java method: DestructuredTarget.inferType()
func (this_ *DestructuredTarget) InferType() pkg_rhino.JSType {
	var inferredType pkg_rhino.JSType = this_.InferTypeWithoutUsingDefaultValue()
	_ = inferredType
	if !(inferredType.GetJSTypeBase().IsUnknownType()) && this_.HasDefaultValue() {
		var defaultValueType pkg_rhino.JSType = this_.GetDefaultValue().GetJSType()
		_ = defaultValueType
		if defaultValueType == nil || defaultValueType.IsJavaNil_() {
			defaultValueType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
		return this_.registry.CreateUnionType_2(inferredType.RestrictByNotUndefined(), defaultValueType)
	} else {
		return inferredType
	}
}

// Java method: DestructuredTarget.inferObjectPatternKeyType()
func (this_ *DestructuredTarget) priv_DestructuredTarget_inferObjectPatternKeyType_2428() pkg_rhino.JSType {
	var patternType pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](this_.patternTypeSupplier.Get())
	_ = patternType
	if this_.isRest {
		return this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
	}
	if (patternType == nil || patternType.IsJavaNil_()) || patternType.GetJSTypeBase().IsUnknownType() {
		return this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	switch this_.objectPatternKey.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */ :
		var propertyType pkg_rhino.JSType = patternType.GetJSTypeBase().FindPropertyType_1_string(this_.objectPatternKey.GetString())
		_ = propertyType
		{

			if propertyType != nil && !propertyType.IsJavaNil_() {
				return propertyType
			}
			return this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}

	case int32(129) /* Token.COMPUTED_PROP */ :
		{

			if patternType != nil && !patternType.IsJavaNil_() {
				return patternType.GetTemplateTypeMap().GetResolvedTemplateType(this_.registry.GetObjectElementKey())
			}
			return this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}

	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected key "+j2gort.StrVal(this_.objectPatternKey))))
	}
	panic("unreachable switch expression")
}

// Java method: DestructuredTarget.inferArrayPatternTargetType()
func (this_ *DestructuredTarget) priv_DestructuredTarget_inferArrayPatternTargetType_2429() pkg_rhino.JSType {
	var patternType pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](this_.patternTypeSupplier.Get())
	_ = patternType
	var templateTypeOfIterable pkg_rhino.JSType = patternType.GetTemplateTypeMap().GetResolvedTemplateType(this_.registry.GetIterableValueTemplate())
	_ = templateTypeOfIterable
	if this_.isRest {
		return pkg_rhino.ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_JSType(this_.registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), []pkg_rhino.JSType{templateTypeOfIterable}))
	} else {
		return templateTypeOfIterable
	}
}

// Java class: com.google.javascript.jscomp.DestructuredTarget.Builder
type DestructuredTarget_Builder struct {
	registry            *pkg_rhino.JSTypeRegistry                  `gson:"registry"`
	patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier `gson:"patternTypeSupplier"`
	pattern             pkg_rhino.Node                             `gson:"pattern"`
	node                pkg_rhino.Node                             `gson:"node"`
	defaultValue        pkg_rhino.Node                             `gson:"defaultValue"`
	objectPatternKey    pkg_rhino.Node                             `gson:"objectPatternKey"`
	idHash_             atomic.Int32
	isRest              bool `gson:"isRest"`
}

func (this_ *DestructuredTarget_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *DestructuredTarget_Builder) IsJavaNil_() bool { return this_ == nil }
func (*DestructuredTarget_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DestructuredTarget$Builder", "Builder"
}
func (this_ *DestructuredTarget_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DestructuredTarget_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DestructuredTarget.Builder", "com.google.javascript.jscomp.DestructuredTarget$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*DestructuredTarget_Builder](v_)
	})
}

func (this_ *DestructuredTarget_Builder) initFields_DestructuredTarget_Builder() {
	_ = this_

}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.Node pattern, com.google.common.base.Supplier<com.google.javascript.rhino.jstype.JSType> patternTypeSupplier)
func NewDestructuredTarget_Builder(registry *pkg_rhino.JSTypeRegistry, pattern pkg_rhino.Node, patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier) *DestructuredTarget_Builder {
	this_ := &DestructuredTarget_Builder{}
	InitDestructuredTarget_Builder(this_, registry, pattern, patternTypeSupplier)
	return this_
}

func InitDestructuredTarget_Builder(this_ *DestructuredTarget_Builder, registry *pkg_rhino.JSTypeRegistry, pattern pkg_rhino.Node, patternTypeSupplier j2gort.Ext_com_google_common_base_Supplier) {
	_ = registry
	_ = pattern
	_ = patternTypeSupplier
	this_.initFields_DestructuredTarget_Builder()
	this_.registry = registry
	this_.patternTypeSupplier = patternTypeSupplier
	this_.pattern = pattern
}

// Java method: Builder.setNode(com.google.javascript.rhino.Node node)
func (this_ *DestructuredTarget_Builder) SetNode(node pkg_rhino.Node) *DestructuredTarget_Builder {
	_ = node
	this_.node = node
	return this_
}

// Java method: Builder.setDefaultValue(com.google.javascript.rhino.Node defaultValue)
func (this_ *DestructuredTarget_Builder) SetDefaultValue(defaultValue pkg_rhino.Node) *DestructuredTarget_Builder {
	_ = defaultValue
	this_.defaultValue = defaultValue
	return this_
}

// Java method: Builder.setObjectPatternKey(com.google.javascript.rhino.Node objectPatternKey)
func (this_ *DestructuredTarget_Builder) SetObjectPatternKey(objectPatternKey pkg_rhino.Node) *DestructuredTarget_Builder {
	_ = objectPatternKey
	this_.objectPatternKey = objectPatternKey
	return this_
}

// Java method: Builder.setIsRest(boolean isRest)
func (this_ *DestructuredTarget_Builder) SetIsRest(isRest bool) *DestructuredTarget_Builder {
	_ = isRest
	this_.isRest = isRest
	return this_
}

// Java method: Builder.build()
func (this_ *DestructuredTarget_Builder) Build() *DestructuredTarget {
	j2gort.CheckNotNullNode(this_.node)
	return NewDestructuredTarget(this_.pattern, this_.node, this_.defaultValue, this_.objectPatternKey, this_.registry, this_.patternTypeSupplier, this_.isRest)
}
