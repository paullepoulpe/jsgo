/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AccessControlUtils.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.AccessControlUtils
type AccessControlUtils struct {
	idHash_ atomic.Int32
}

func (this_ *AccessControlUtils) BindSelf_(s any) {
	_ = s
}
func (this_ *AccessControlUtils) IsJavaNil_() bool { return this_ == nil }
func (*AccessControlUtils) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AccessControlUtils", "AccessControlUtils"
}
func (this_ *AccessControlUtils) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AccessControlUtils) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AccessControlUtils", "com.google.javascript.jscomp.AccessControlUtils", func(v_ any) bool {
		return j2gort.InstanceOf[*AccessControlUtils](v_)
	})
}

func (this_ *AccessControlUtils) initFields_AccessControlUtils() {
	_ = this_
}

// Java constructor: AccessControlUtils()
func NewAccessControlUtils() *AccessControlUtils {
	this_ := &AccessControlUtils{}
	InitAccessControlUtils(this_)
	return this_
}

func InitAccessControlUtils(this_ *AccessControlUtils) {
	this_.initFields_AccessControlUtils()
}

// Java static method: AccessControlUtils.getEffectiveNameVisibility(com.google.javascript.rhino.Node name, com.google.javascript.jscomp.Var var, com.google.common.collect.ImmutableMap<com.google.javascript.rhino.StaticSourceFile,com.google.javascript.rhino.JSDocInfo.Visibility> fileVisibilityMap)
func AccessControlUtils_GetEffectiveNameVisibility(name pkg_rhino.Node, var_ *Var, fileVisibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap) pkg_rhino.JSDocInfo_Visibility {
	_ = name
	_ = var_
	_ = fileVisibilityMap
	var jsDocInfo *pkg_rhino.JSDocInfo = var_.GetJSDocInfo()
	_ = jsDocInfo
	var raw pkg_rhino.JSDocInfo_Visibility
	_ = raw
	if (jsDocInfo == nil) || (jsDocInfo.GetVisibility() == 0) {
		raw = pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()
	} else {
		raw = jsDocInfo.GetVisibility()
	}
	if raw != pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED() {
		return raw
	}
	var defaultVisibilityForFile pkg_rhino.JSDocInfo_Visibility = j2gort.Cast[pkg_rhino.JSDocInfo_Visibility](fileVisibilityMap.Get(var_.GetSourceFile()))
	_ = defaultVisibilityForFile
	var type_ pkg_rhino.JSType = name.GetJSType()
	_ = type_
	var createdFromGoogProvide bool = ((type_ != nil && !type_.IsJavaNil_()) && type_.GetJSTypeBase().IsLiteralObject())
	_ = createdFromGoogProvide
	if createdFromGoogProvide || (defaultVisibilityForFile == 0) {
		return raw
	} else {
		return defaultVisibilityForFile
	}
}

// Java static method: AccessControlUtils.getEffectivePropertyVisibility(com.google.javascript.rhino.Node property, com.google.javascript.rhino.jstype.ObjectType referenceType, com.google.common.collect.ImmutableMap<com.google.javascript.rhino.StaticSourceFile,com.google.javascript.rhino.JSDocInfo.Visibility> fileVisibilityMap)
func AccessControlUtils_GetEffectivePropertyVisibility(property pkg_rhino.Node, referenceType pkg_rhino.ObjectType, fileVisibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap) pkg_rhino.JSDocInfo_Visibility {
	_ = property
	_ = referenceType
	_ = fileVisibilityMap
	var propertyName string = property.GetString()
	_ = propertyName
	var definingSource pkg_rhino.StaticSourceFile = AccessControlUtils_GetDefiningSource(property, referenceType, propertyName)
	_ = definingSource
	var fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility = j2gort.Cast[pkg_rhino.JSDocInfo_Visibility](fileVisibilityMap.Get(definingSource))
	_ = fileOverviewVisibility
	var parent pkg_rhino.Node = property.GetParent()
	_ = parent
	var isOverride bool = (((parent.GetJSDocInfo() != nil) && parent.IsAssign()) && (parent.GetFirstChild() == property))
	_ = isOverride
	var objectType pkg_rhino.ObjectType = AccessControlUtils_GetObjectType(referenceType, isOverride, propertyName)
	_ = objectType
	if isOverride {
		var overridden pkg_rhino.JSDocInfo_Visibility = AccessControlUtils_GetOverriddenPropertyVisibility(objectType, propertyName)
		_ = overridden
		return AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(overridden, fileOverviewVisibility, propertyName)
	} else {
		return AccessControlUtils_GetEffectiveVisibilityForNonOverriddenProperty(property, objectType, fileOverviewVisibility)
	}
}

// Java static method: AccessControlUtils.getDefiningSource(com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType referenceType, java.lang.String propertyName)
func AccessControlUtils_GetDefiningSource(node pkg_rhino.Node, referenceType pkg_rhino.ObjectType, propertyName string) pkg_rhino.StaticSourceFile {
	_ = node
	_ = referenceType
	_ = propertyName
	if referenceType != nil && !referenceType.IsJavaNil_() {
		var propDefNode pkg_rhino.Node = referenceType.GetObjectTypeBase().GetPropertyDefSite_1_string(propertyName)
		_ = propDefNode
		if propDefNode != 0 {
			return propDefNode.GetStaticSourceFile()
		}
	}
	return node.GetStaticSourceFile()
}

// Java static method: AccessControlUtils.getObjectType(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType referenceType, boolean isOverride, java.lang.String propertyName)
func AccessControlUtils_GetObjectType(referenceType pkg_rhino.ObjectType, isOverride bool, propertyName string) pkg_rhino.ObjectType {
	_ = referenceType
	_ = isOverride
	_ = propertyName
	if referenceType == nil || referenceType.IsJavaNil_() {
		return nil
	}
	var current pkg_rhino.ObjectType
	_ = current
	if isOverride {
		current = referenceType.GetImplicitPrototype()
	} else {
		current = referenceType
	}
	for ; current != nil && !current.IsJavaNil_(); current = current.GetImplicitPrototype() {
		var docInfo *pkg_rhino.JSDocInfo = current.GetObjectTypeBase().GetOwnPropertyJSDocInfo_1_string(propertyName)
		_ = docInfo
		if (docInfo != nil) && (docInfo.GetVisibility() != pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) {
			return current
		}
	}
	return nil
}

// Java static method: AccessControlUtils.getOverriddenPropertyVisibility(com.google.javascript.rhino.jstype.ObjectType objectType, java.lang.String propertyName)
func AccessControlUtils_GetOverriddenPropertyVisibility(objectType pkg_rhino.ObjectType, propertyName string) pkg_rhino.JSDocInfo_Visibility {
	_ = objectType
	_ = propertyName
	if objectType != nil && !objectType.IsJavaNil_() {
		return objectType.GetObjectTypeBase().GetOwnPropertyJSDocInfo_1_string(propertyName).GetVisibility()
	} else {
		return pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()
	}
}

// Java static method: AccessControlUtils.getEffectiveVisibilityForOverriddenProperty(com.google.javascript.rhino.JSDocInfo.Visibility visibility, com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable Visibility fileOverviewVisibility, java.lang.String propertyName)
func AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(visibility pkg_rhino.JSDocInfo_Visibility, fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility, propertyName string) pkg_rhino.JSDocInfo_Visibility {
	_ = visibility
	_ = fileOverviewVisibility
	_ = propertyName
	if (fileOverviewVisibility != 0) && (visibility == pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) {
		return fileOverviewVisibility
	} else {
		return visibility
	}
}

// Java static method: AccessControlUtils.getEffectiveVisibilityForNonOverriddenProperty(com.google.javascript.rhino.Node getprop, com.google.javascript.rhino.jstype.ObjectType objectType, com.google.javascript.rhino.JSDocInfo.@org.jspecify.annotations.Nullable Visibility fileOverviewVisibility)
func AccessControlUtils_GetEffectiveVisibilityForNonOverriddenProperty(getprop pkg_rhino.Node, objectType pkg_rhino.ObjectType, fileOverviewVisibility pkg_rhino.JSDocInfo_Visibility) pkg_rhino.JSDocInfo_Visibility {
	_ = getprop
	_ = objectType
	_ = fileOverviewVisibility
	var propertyName string = getprop.GetString()
	_ = propertyName
	var raw pkg_rhino.JSDocInfo_Visibility = pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()
	_ = raw
	if objectType != nil && !objectType.IsJavaNil_() {
		raw = objectType.GetObjectTypeBase().GetOwnPropertyJSDocInfo_1_string(propertyName).GetVisibility()
	}
	var type_ pkg_rhino.JSType = getprop.GetJSType()
	_ = type_
	var createdFromGoogProvide bool = ((type_ != nil && !type_.IsJavaNil_()) && type_.GetJSTypeBase().IsLiteralObject())
	_ = createdFromGoogProvide
	if ((raw != pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED()) || (fileOverviewVisibility == 0)) || createdFromGoogProvide {
		return raw
	} else {
		return fileOverviewVisibility
	}
}
