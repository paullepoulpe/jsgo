// Port of com/google/javascript/jscomp/AccessControlUtils.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AccessControlUtils holds helper functions for computing the visibility of
// names and properties in JavaScript source code.
//
// See CheckAccessControls.
//
// Java: AccessControlUtils() is private (non-instantiable); the class has
// only static members.
//
// Go-only: language-adapter — a nullable Visibility is a
// *rhino.JSDocInfo_Visibility (nil is Java null), and ImmutableMap.get
// returning null is LinkedHashMap.Lookup reporting false.
type AccessControlUtils struct{}

// AccessControlUtils_GetEffectiveNameVisibility returns the effective
// visibility of the given name. This can differ from the name's declared
// visibility if the file's @fileoverview JsDoc specifies a default
// visibility.
//
// name is the name node to compute effective visibility for; var_ the name
// to compute effective visibility for; fileVisibilityMap a map of
// @fileoverview visibility annotations, used to compute the name's default
// visibility.
//
// Java: AccessControlUtils.getEffectiveNameVisibility(Node, Var, ImmutableMap<StaticSourceFile, Visibility>)
func AccessControlUtils_GetEffectiveNameVisibility(
	name rhino.Node, var_ *Var,
	fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]) rhino.JSDocInfo_Visibility {
	jsDocInfo := var_.GetJSDocInfo()
	// Go-only: language-adapter — JSDocInfo.getVisibility() never returns
	// null (it defaults to INHERITED), so Java's
	// `jsDocInfo.getVisibility() == null` operand is always false.
	var raw rhino.JSDocInfo_Visibility
	if jsDocInfo == nil {
		raw = rhino.JSDocInfo_Visibility_INHERITED
	} else {
		raw = jsDocInfo.GetVisibility()
	}
	if raw != rhino.JSDocInfo_Visibility_INHERITED {
		return raw
	}
	defaultVisibilityForFile, hasDefaultVisibilityForFile := fileVisibilityMap.Lookup(var_.GetSourceFile())
	type_ := name.GetJSType()
	createdFromGoogProvide := type_ != nil && type_.IsLiteralObject()
	// Ignore @fileoverview visibility when computing the effective visibility
	// for names created by goog.provide.
	//
	// ProcessClosurePrimitives rewrites goog.provide()s as object literal
	// declarations, but the exact form depends on the ordering of the
	// input files. If goog.provide('a.b') occurs in the inputs before
	// goog.provide('a'), it is rewritten like
	//
	// var a={};a.b={};
	//
	// If the file containing goog.provide('a.b') also declares a @fileoverview
	// visibility, it must not apply to a, as this would make every a.* namespace
	// effectively package-private.
	if createdFromGoogProvide || !hasDefaultVisibilityForFile {
		return raw
	}
	return defaultVisibilityForFile
}

// AccessControlUtils_GetEffectivePropertyVisibility returns the effective
// visibility of the given property. This can differ from the property's
// declared visibility if the property is inherited from a superclass, or if
// the file's @fileoverview JsDoc specifies a default visibility.
//
// property is the property to compute effective visibility for;
// referenceType the JavaScript type of the property; fileVisibilityMap a map
// of @fileoverview visibility annotations, used to compute the property's
// default visibility.
//
// Java: AccessControlUtils.getEffectivePropertyVisibility(Node, ObjectType, ImmutableMap<StaticSourceFile, Visibility>)
func AccessControlUtils_GetEffectivePropertyVisibility(
	property rhino.Node,
	referenceType rhino.ObjectType,
	fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]) rhino.JSDocInfo_Visibility {
	propertyName := property.GetString()
	definingSource := AccessControlUtils_GetDefiningSource(property, referenceType, propertyName)
	var fileOverviewVisibility *rhino.JSDocInfo_Visibility
	if v, ok := fileVisibilityMap.Lookup(definingSource); ok {
		fileOverviewVisibility = &v
	}
	parent := property.GetParent()
	isOverride := parent.GetJSDocInfo() != nil &&
		parent.IsAssign() &&
		parent.GetFirstChild() == property
	objectType := AccessControlUtils_GetObjectType(
		referenceType, isOverride, propertyName)
	if isOverride {
		overridden := AccessControlUtils_GetOverriddenPropertyVisibility(
			objectType, propertyName)
		return AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(
			overridden, fileOverviewVisibility, propertyName)
	} else {
		return AccessControlUtils_GetEffectiveVisibilityForNonOverriddenProperty(
			property, objectType, fileOverviewVisibility)
	}
}

// AccessControlUtils_GetDefiningSource returns the source file in which the
// given property is defined, or null if it is not known.
//
// Java: AccessControlUtils.getDefiningSource(Node, ObjectType, String)
func AccessControlUtils_GetDefiningSource(
	node rhino.Node, referenceType rhino.ObjectType, propertyName string) rhino.StaticSourceFile {
	if referenceType != nil {
		propDefNode := referenceType.GetPropertyDefSite(propertyName)
		if !propDefNode.IsNil() {
			return propDefNode.GetStaticSourceFile()
		}
	}
	return node.GetStaticSourceFile()
}

// AccessControlUtils_GetObjectType returns the lowest property defined on a
// class with visibility information.
//
// Java: AccessControlUtils.getObjectType(ObjectType, boolean, String)
func AccessControlUtils_GetObjectType(
	referenceType rhino.ObjectType, isOverride bool, propertyName string) rhino.ObjectType {
	if referenceType == nil {
		return nil
	}

	// Find the lowest property defined on a class with visibility information.
	var current rhino.ObjectType
	if isOverride {
		current = referenceType.GetImplicitPrototype()
	} else {
		current = referenceType
	}
	for ; current != nil; current = current.GetImplicitPrototype() {
		docInfo := current.GetOwnPropertyJSDocInfo(propertyName)
		if docInfo != nil && docInfo.GetVisibility() != rhino.JSDocInfo_Visibility_INHERITED {
			return current
		}
	}
	return nil
}

// AccessControlUtils_GetOverriddenPropertyVisibility returns the original
// visibility of an overridden property.
//
// Java: AccessControlUtils.getOverriddenPropertyVisibility(ObjectType, String)
func AccessControlUtils_GetOverriddenPropertyVisibility(objectType rhino.ObjectType, propertyName string) rhino.JSDocInfo_Visibility {
	if objectType != nil {
		return objectType.GetOwnPropertyJSDocInfo(propertyName).GetVisibility()
	}
	return rhino.JSDocInfo_Visibility_INHERITED
}

// AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty returns the
// effective visibility of the given overridden property. An overridden
// property inherits the visibility of the property it overrides.
//
// Java: AccessControlUtils.getEffectiveVisibilityForOverriddenProperty(Visibility, Visibility, String)
func AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(
	visibility rhino.JSDocInfo_Visibility, fileOverviewVisibility *rhino.JSDocInfo_Visibility, propertyName string) rhino.JSDocInfo_Visibility {
	if fileOverviewVisibility != nil &&
		visibility == rhino.JSDocInfo_Visibility_INHERITED {
		return *fileOverviewVisibility
	}
	return visibility
}

// AccessControlUtils_GetEffectiveVisibilityForNonOverriddenProperty returns
// the effective visibility of the given non-overridden property.
// Non-overridden properties without an explicit visibility annotation
// receive the default visibility declared in the file's @fileoverview block,
// if one exists.
//
// Java: AccessControlUtils.getEffectiveVisibilityForNonOverriddenProperty(Node, ObjectType, Visibility)
func AccessControlUtils_GetEffectiveVisibilityForNonOverriddenProperty(
	getprop rhino.Node, objectType rhino.ObjectType, fileOverviewVisibility *rhino.JSDocInfo_Visibility) rhino.JSDocInfo_Visibility {
	propertyName := getprop.GetString()
	raw := rhino.JSDocInfo_Visibility_INHERITED
	if objectType != nil {
		raw = objectType.GetOwnPropertyJSDocInfo(propertyName).GetVisibility()
	}
	type_ := getprop.GetJSType()
	createdFromGoogProvide := type_ != nil && type_.IsLiteralObject()
	// Ignore @fileoverview visibility when computing the effective visibility
	// for properties created by goog.provide.
	//
	// ProcessClosurePrimitives rewrites goog.provide()s as object literal
	// declarations, but the exact form depends on the ordering of the
	// input files. If goog.provide('a.b.c') occurs in the inputs before
	// goog.provide('a'), it is rewritten like
	//
	// var a={};a.b={}a.b.c={};
	//
	// If the file containing goog.provide('a.b.c') also declares
	// a @fileoverview visibility, it must not apply to b, as this would make
	// every a.b.* namespace effectively package-private.
	if raw != rhino.JSDocInfo_Visibility_INHERITED ||
		fileOverviewVisibility == nil ||
		createdFromGoogProvide {
		return raw
	}
	return *fileOverviewVisibility
}
