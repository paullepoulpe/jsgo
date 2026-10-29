// Port of com/google/javascript/jscomp/CheckAccessControls.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoBuilder_CheckAccessControls_PropertyReference_Builder.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckPropertyDeprecation checks the given GETPROP node to ensure that
// access restrictions are obeyed.
//
// Java: CheckAccessControls.checkPropertyDeprecation(NodeTraversal, PropertyReference)
func (c *CheckAccessControls) CheckPropertyDeprecation(t *NodeTraversal, propRef *CheckAccessControls_PropertyReference) {
	if !c.ShouldEmitDeprecationWarningForPropRef(t, propRef) {
		return
	}

	// Don't bother checking constructors.
	if propRef.sourceNode.GetParent().IsNew() {
		return
	}

	objectType := CheckAccessControls_CastToObject(CheckAccessControls_Dereference(propRef.receiverType))
	propertyName := propRef.name

	if objectType != nil {
		if deprecationInfo, ok := CheckAccessControls_GetPropertyDeprecationInfo(objectType, propertyName); ok {
			if deprecationInfo != "" {
				c.compiler.Report(
					JSError_MakeWithNode(
						propRef.sourceNode,
						CheckAccessControls_DEPRECATED_PROP_REASON,
						propertyName,
						propRef.GetReadableTypeNameOrDefault(),
						deprecationInfo))
			} else {
				c.compiler.Report(
					JSError_MakeWithNode(
						propRef.sourceNode,
						CheckAccessControls_DEPRECATED_PROP,
						propertyName,
						propRef.GetReadableTypeNameOrDefault()))
			}
		}
	}
}

// Java: CheckAccessControls.checkPropertyOverrideVisibilityIsSame(Visibility, Visibility, Visibility, PropertyReference)
func (c *CheckAccessControls) CheckPropertyOverrideVisibilityIsSame(
	overriding rhino.JSDocInfo_Visibility,
	overridden rhino.JSDocInfo_Visibility,
	fileOverview *rhino.JSDocInfo_Visibility,
	propRef *CheckAccessControls_PropertyReference) {
	if overriding == rhino.JSDocInfo_Visibility_INHERITED &&
		overriding != overridden &&
		fileOverview != nil &&
		*fileOverview != rhino.JSDocInfo_Visibility_INHERITED {
		c.compiler.Report(
			JSError_MakeWithNode(
				propRef.sourceNode,
				CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY,
				propRef.name,
				fileOverview.String()))
	}
}

// Java: CheckAccessControls.getOverridingPropertyVisibility(PropertyReference)
func CheckAccessControls_GetOverridingPropertyVisibility(propRef *CheckAccessControls_PropertyReference) (rhino.JSDocInfo_Visibility, bool) {
	overridingInfo := propRef.GetJSDocInfo()
	if overridingInfo == nil || !overridingInfo.IsOverride() {
		return 0, false
	}
	return overridingInfo.GetVisibility(), true
}

// CheckConstantProperty determines whether the given constant property got
// reassigned.
//
// Java: CheckAccessControls.checkConstantProperty(PropertyReference, IdentifierBehaviour)
func (c *CheckAccessControls) CheckConstantProperty(
	propRef *CheckAccessControls_PropertyReference, identifierBehaviour CheckAccessControls_IdentifierBehaviour) {
	if propRef == nil || identifierBehaviour == CheckAccessControls_IdentifierBehaviour_ES5_CLASS_NAMESPACE {
		return
	}

	sourceNode := propRef.sourceNode
	if !propRef.mutation && !(propRef.declaration && sourceNode.IsFromExterns()) {
		return
	}

	objectType := CheckAccessControls_Dereference(propRef.receiverType)
	propertyName := propRef.name

	constness := c.IsPropertyDeclaredConstant(objectType, propertyName)
	if constness == CheckAccessControls_Constancy_MUTABLE {
		return
	}

	if sourceNode.IsFromExterns() && propRef.declaration {
		// Treat stub declarations in externs as inits, but never warn on them.
		c.RecordConstPropertyInit(propRef, objectType, constness)
		return
	}

	if !propRef.mutation {
		return
	}

	if propRef.IsDeletion() {
		c.compiler.Report(JSError_MakeWithNode(sourceNode, CheckAccessControls_CONST_PROPERTY_DELETED, propertyName))
		return
	}

	// Can't check for constant properties on generic function types.
	// TODO(johnlenz): I'm not 100% certain this is necessary, or if
	// the type is being inspected incorrectly.
	if objectType.IsFunctionType() && !objectType.ToMaybeFunctionType().IsConstructor() {
		return
	}

	if objectType.IsStructuralType() && !propRef.declaration {
		// We don't know the claess this structural type matches, so assume all assignments are bad.
		c.compiler.Report(
			JSError_MakeWithNode(
				sourceNode,
				CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE,
				propertyName,
				"unknown location due to structural typing"))
		return
	}

	if init, ok := c.GetConstPropertyInit(propRef, objectType); ok {
		diagnostic := CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE
		if init.annotation == CheckAccessControls_Constancy_FINAL {
			diagnostic = CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN
		}
		c.compiler.Report(JSError_MakeWithNode(sourceNode, diagnostic, propertyName, init.node.GetLocation()))
	}

	c.RecordConstPropertyInit(propRef, objectType, constness)
}

func isEquivConstOwnerType(t rhino.JSType) bool {
	return t.IsTemplatizedType() || t.IsRecordType()
}

func (c *CheckAccessControls) checkAccessControls_tableGet(rowKey rhino.JSType, columnKey string) (CheckAccessControls_ConstantDeclaration, bool) {
	if rowKey == nil {
		return CheckAccessControls_ConstantDeclaration{}, false
	}
	if isEquivConstOwnerType(rowKey) {
		for i := range c.equivConstPropertyInits {
			e := &c.equivConstPropertyInits[i]
			if e.prop == columnKey && (e.owner == rowKey || rowKey.Equals(e.owner)) {
				return e.decl, true
			}
		}
		return CheckAccessControls_ConstantDeclaration{}, false
	}
	if c.constPropertyInits == nil {
		return CheckAccessControls_ConstantDeclaration{}, false
	}
	init, ok := c.constPropertyInits[checkAccessControlsConstPropKey{owner: rowKey, prop: columnKey}]
	return init, ok
}

func (c *CheckAccessControls) checkAccessControls_tableRowPutIfAbsent(
	rowKey rhino.JSType, columnKey string, value CheckAccessControls_ConstantDeclaration) {
	if isEquivConstOwnerType(rowKey) {
		for i := range c.equivConstPropertyInits {
			e := &c.equivConstPropertyInits[i]
			if e.prop == columnKey && (e.owner == rowKey || rowKey.Equals(e.owner)) {
				return
			}
		}
		c.equivConstPropertyInits = append(c.equivConstPropertyInits, checkAccessControlsEquivConstPropEntry{
			owner: rowKey,
			prop:  columnKey,
			decl:  value,
		})
		return
	}
	if c.constPropertyInits == nil {
		c.constPropertyInits = make(map[checkAccessControlsConstPropKey]CheckAccessControls_ConstantDeclaration)
	}
	key := checkAccessControlsConstPropKey{owner: rowKey, prop: columnKey}
	if _, exists := c.constPropertyInits[key]; !exists {
		c.constPropertyInits[key] = value
	}
}

// Java: CheckAccessControls.getConstPropertyInit(PropertyReference, ObjectType)
func (c *CheckAccessControls) GetConstPropertyInit(
	ref *CheckAccessControls_PropertyReference, type_ rhino.ObjectType) (CheckAccessControls_ConstantDeclaration, bool) {
	name := ref.name
	for type_ != nil {
		if init, ok := c.checkAccessControls_tableGet(type_, name); ok {
			return init, true
		}
		if canonicalInit, ok := c.checkAccessControls_tableGet(CheckAccessControls_GetCanonicalInstance(type_), name); ok {
			return canonicalInit, true
		}

		type_ = type_.GetImplicitPrototype()
	}

	return CheckAccessControls_ConstantDeclaration{}, false
}

// Java: CheckAccessControls.recordConstPropertyInit(PropertyReference, ObjectType, Constancy)
func (c *CheckAccessControls) RecordConstPropertyInit(
	ref *CheckAccessControls_PropertyReference, type_ rhino.ObjectType, annotation CheckAccessControls_Constancy) {
	decl := CheckAccessControls_ConstantDeclaration{node: ref.sourceNode, annotation: annotation}
	c.checkAccessControls_tableRowPutIfAbsent(type_, ref.name, decl)

	// Add the prototype when we're looking at an instance object
	if type_.IsInstanceType() {
		prototype := type_.GetImplicitPrototype()
		if prototype != nil && prototype.HasProperty(ref.name) {
			c.checkAccessControls_tableRowPutIfAbsent(prototype, ref.name, decl)
		}
	}
}

// CheckAccessControls_GetCanonicalInstance returns an object with the same
// nominal type as obj, but without any possible extra properties that exist
// on obj.
//
// Java: CheckAccessControls.getCanonicalInstance(ObjectType)
func CheckAccessControls_GetCanonicalInstance(obj rhino.ObjectType) rhino.ObjectType {
	ctor := obj.GetConstructor()
	if ctor == nil {
		return obj
	}
	return ctor.GetInstanceType()
}

// CheckAccessControls_Dereference dereferences a type, autoboxing it and
// filtering out null.
//
// Java: CheckAccessControls.dereference(JSType)
func CheckAccessControls_Dereference(type_ rhino.JSType) rhino.ObjectType {
	if type_ == nil {
		return nil
	}
	return type_.Dereference()
}

// Java: CheckAccessControls.typeOrUnknown(JSType)
func (c *CheckAccessControls) TypeOrUnknown(type_ rhino.JSType) rhino.JSType {
	if type_ == nil {
		return c.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	return type_
}

// Java: CheckAccessControls.typeOrUnknown(ObjectType)
func (c *CheckAccessControls) TypeOrUnknownObjectType(type_ rhino.ObjectType) rhino.ObjectType {
	return c.TypeOrUnknown(type_).(rhino.ObjectType)
}

// Java: CheckAccessControls.boxedOrUnknown(JSType)
func (c *CheckAccessControls) BoxedOrUnknown(type_ rhino.JSType) rhino.ObjectType {
	return c.TypeOrUnknownObjectType(CheckAccessControls_Dereference(type_))
}

// CheckPropertyVisibility reports an error if the given property is not
// visible in the current context.
//
// This method covers both:
//
//   - accesses to properties during execution
//   - overrides of properties during declaration
//
// TODO(nickreid): Things would probably be a lot simpler, though a bit
// duplicated, if these two concepts were separated. Much of the underlying
// logic could stop checking various inconsistent definitions of "is this an
// override".
//
// Java: CheckAccessControls.checkPropertyVisibility(PropertyReference)
func (c *CheckAccessControls) CheckPropertyVisibility(propRef *CheckAccessControls_PropertyReference) {
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(propRef.sourceNode) {
		// Class ctor *declarations* can never violate visibility restrictions. They are not
		// accesses and we don't consider them overrides.
		//
		// TODO(nickreid): It would be a lot cleaner if we could model this using `PropertyReference`
		// rather than defining a special case here. I think the problem is that the current
		// implementation of this method conflates "override" with "declaration". But that only works
		// because it ignores cases where there's no overridden definition.
		return
	}

	rawReferenceType := c.TypeOrUnknownObjectType(propRef.receiverType).Autobox()
	referenceType := CheckAccessControls_CastToObject(rawReferenceType)

	propertyName := propRef.name

	// Is this a normal property access, or are we trying to override
	// an existing property?
	isOverride := propRef.IsDocumentedDeclaration() || propRef.override

	objectType := AccessControlUtils_GetObjectType(referenceType, isOverride, propertyName)

	var definingSource rhino.StaticSourceFile
	var fileOverviewVisibility *rhino.JSDocInfo_Visibility
	if len(c.defaultVisibilityForFiles) > 0 {
		definingSource = AccessControlUtils_GetDefiningSource(propRef.sourceNode, referenceType, propertyName)
		fileOverviewVisibility = checkAccessControls_visibilityFor(c.defaultVisibilityForFiles, definingSource)
	}

	var defSite rhino.Node
	if objectType != nil {
		defSite = objectType.GetOwnPropertyDefSite(propertyName)
		if defSite.IsNil() && fileOverviewVisibility == nil {
			return
		}
	} else if fileOverviewVisibility == nil {
		return
	}

	var visibility rhino.JSDocInfo_Visibility
	if propRef.override {
		effObjType := objectType
		if !isOverride {
			effObjType = AccessControlUtils_GetObjectType(referenceType, true, propertyName)
		}
		overridden := AccessControlUtils_GetOverriddenPropertyVisibility(effObjType, propertyName)
		visibility = AccessControlUtils_GetEffectiveVisibilityForOverriddenProperty(
			overridden, fileOverviewVisibility, propertyName)
	} else {
		effObjType := objectType
		if isOverride {
			effObjType = AccessControlUtils_GetObjectType(referenceType, false, propertyName)
		}
		visibility = CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
			propRef, effObjType, fileOverviewVisibility)
	}

	if isOverride {
		if overriding, ok := CheckAccessControls_GetOverridingPropertyVisibility(propRef); ok {
			c.CheckPropertyOverrideVisibilityIsSame(
				overriding, visibility, fileOverviewVisibility, propRef)
		}
	}

	reportType := rawReferenceType
	if objectType != nil {
		if defSite.IsNil() {
			// Assume the property is public.
			return
		}
		reportType = objectType
		definingSource = defSite.GetStaticSourceFile()
	} else if fileOverviewVisibility == nil {
		// We can only check visibility references if we know what file
		// it was defined in.
		// Otherwise just assume the property is public.
		return
	}

	referenceSource := propRef.sourceNode.GetStaticSourceFile()

	if isOverride {
		sameInput := referenceSource != nil && referenceSource.GetName() == definingSource.GetName()
		c.CheckPropertyOverrideVisibility(
			propRef, visibility, fileOverviewVisibility, reportType, sameInput)
	} else {
		c.CheckPropertyAccessVisibility(
			propRef, visibility, reportType, referenceSource, definingSource)
	}
}

var checkAccessControls_visibilityPtrs = [...]rhino.JSDocInfo_Visibility{
	rhino.JSDocInfo_Visibility_PRIVATE:   rhino.JSDocInfo_Visibility_PRIVATE,
	rhino.JSDocInfo_Visibility_PACKAGE:   rhino.JSDocInfo_Visibility_PACKAGE,
	rhino.JSDocInfo_Visibility_PROTECTED: rhino.JSDocInfo_Visibility_PROTECTED,
	rhino.JSDocInfo_Visibility_PUBLIC:    rhino.JSDocInfo_Visibility_PUBLIC,
	rhino.JSDocInfo_Visibility_INHERITED: rhino.JSDocInfo_Visibility_INHERITED,
}

func checkAccessControls_visibilityFor(
	m map[rhino.StaticSourceFile]rhino.JSDocInfo_Visibility,
	file rhino.StaticSourceFile) *rhino.JSDocInfo_Visibility {
	if file == nil || len(m) == 0 {
		return nil
	}
	if v, ok := m[file]; ok {
		return &checkAccessControls_visibilityPtrs[v]
	}
	return nil
}

// CheckEs6ConstructorInvocationVisibility reports visibility violations on
// ES6 class constructor invocations.
//
// Precondition: target has an ES6 class JSType.
//
// Java: CheckAccessControls.checkEs6ConstructorInvocationVisibility(Node)
func (c *CheckAccessControls) CheckEs6ConstructorInvocationVisibility(target rhino.Node) {
	ctorType := target.GetJSType().ToMaybeFunctionType()
	prototypeType := ctorType.GetPrototype()

	// We use the class definition site because classes automatically get a implicit constructor,
	// so there may not be a definition node.
	classDefinition := ctorType.GetSource()

	var definingSource rhino.StaticSourceFile
	if !classDefinition.IsNil() {
		definingSource = AccessControlUtils_GetDefiningSource(classDefinition, prototypeType, "constructor")
	}

	// Synthesize a `PropertyReference` for this constructor call as if we're accessing
	// `Foo.prototype.constructor`. This object allows us to reuse the
	// `checkPropertyAccessVisibility` method which actually reports violations.
	fauxCtorRef := CheckAccessControls_PropertyReference{
		sourceNode:           target,
		name:                 "constructor",
		receiverType:         prototypeType,
		readableInstanceType: ctorType.GetInstanceType(),
	}

	annotatedCtorVisibility :=
		// This function defaults to `INHERITED` which isn't what we want here, but it does handle
		// combining inline and `@fileoverview` visibilities.
		CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
			&fauxCtorRef, prototypeType, checkAccessControls_visibilityFor(c.defaultVisibilityForFiles, definingSource))
	effectiveCtorVisibility := annotatedCtorVisibility
	if annotatedCtorVisibility == rhino.JSDocInfo_Visibility_INHERITED {
		effectiveCtorVisibility = rhino.JSDocInfo_Visibility_PUBLIC
	}

	c.CheckPropertyAccessVisibility(
		&fauxCtorRef,
		effectiveCtorVisibility,
		ctorType,
		target.GetStaticSourceFile(),
		definingSource)
}

// Java: CheckAccessControls.checkPropertyOverrideVisibility(PropertyReference, Visibility, Visibility, JSType, boolean)
func (c *CheckAccessControls) CheckPropertyOverrideVisibility(
	propRef *CheckAccessControls_PropertyReference,
	visibility rhino.JSDocInfo_Visibility,
	fileOverviewVisibility *rhino.JSDocInfo_Visibility,
	objectType rhino.JSType,
	sameInput bool) {
	overridingVisibility := rhino.JSDocInfo_Visibility_INHERITED
	if propRef.override {
		overridingVisibility = propRef.GetJSDocInfo().GetVisibility()
	}

	// Check that:
	// (a) the property *can* be overridden,
	// (b) the visibility of the override is the same as (or broader than) the
	//     visibility of the original property,
	// (c) the visibility is explicitly redeclared if the override is in
	//     a file with default visibility in the @fileoverview block.
	if visibility == rhino.JSDocInfo_Visibility_PRIVATE && !sameInput {
		c.compiler.Report(JSError_MakeWithNode(propRef.sourceNode, CheckAccessControls_PRIVATE_OVERRIDE, objectType.String()))
	} else if !CheckAccessControls_CanOverrideVisibility(visibility, overridingVisibility) &&
		fileOverviewVisibility == nil {
		c.compiler.Report(
			JSError_MakeWithNode(
				propRef.sourceNode,
				CheckAccessControls_VISIBILITY_MISMATCH,
				visibility.String(),
				objectType.String(),
				overridingVisibility.String()))
	}
}

// Java: CheckAccessControls.canOverrideVisibility(Visibility, Visibility)
func CheckAccessControls_CanOverrideVisibility(
	superclassVisibility rhino.JSDocInfo_Visibility, subclassVisibility rhino.JSDocInfo_Visibility) bool {
	// This allows INHERITED to override anything, PUBLIC to override anything (except INHERITED),
	// and PROTECTED to override anything (except PUBLIC or INHERITED). PRIVATE was already handled
	// in a previous check, leaving PACKAGE as the lowest visibility.
	return superclassVisibility.Ordinal()-subclassVisibility.Ordinal() <= 0
}

// Java: CheckAccessControls.checkPropertyAccessVisibility(PropertyReference, Visibility, JSType, StaticSourceFile, StaticSourceFile)
func (c *CheckAccessControls) CheckPropertyAccessVisibility(
	propRef *CheckAccessControls_PropertyReference,
	visibility rhino.JSDocInfo_Visibility,
	objectType rhino.JSType,
	referenceSource rhino.StaticSourceFile,
	definingSource rhino.StaticSourceFile) {
	// private access is always allowed in the same file.
	if referenceSource != nil &&
		definingSource != nil &&
		referenceSource.GetName() == definingSource.GetName() {
		return
	}

	ownerType := CheckAccessControls_InstanceTypeFor(objectType)

	switch visibility {
	case rhino.JSDocInfo_Visibility_PACKAGE:
		c.CheckPackagePropertyVisibility(propRef, referenceSource, definingSource)
	case rhino.JSDocInfo_Visibility_PRIVATE:
		c.CheckPrivatePropertyVisibility(propRef, ownerType)
	case rhino.JSDocInfo_Visibility_PROTECTED:
		c.CheckProtectedPropertyVisibility(propRef, ownerType)
	default:
	}
}

// Java: CheckAccessControls.checkPackagePropertyVisibility(PropertyReference, StaticSourceFile, StaticSourceFile)
func (c *CheckAccessControls) CheckPackagePropertyVisibility(
	propRef *CheckAccessControls_PropertyReference,
	referenceSource rhino.StaticSourceFile,
	definingSource rhino.StaticSourceFile) {
	codingConvention := c.compiler.GetCodingConvention()
	refPackage := codingConvention.GetPackageName(referenceSource)
	defPackage := codingConvention.GetPackageName(definingSource)
	if refPackage == nil || defPackage == nil || *refPackage != *defPackage {
		c.compiler.Report(
			JSError_MakeWithNode(
				propRef.sourceNode,
				CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS,
				propRef.name,
				propRef.GetReadableTypeNameOrDefault()))
	}
}

// Java: CheckAccessControls.checkPrivatePropertyVisibility(PropertyReference, ObjectType)
func (c *CheckAccessControls) CheckPrivatePropertyVisibility(
	propRef *CheckAccessControls_PropertyReference, ownerType rhino.ObjectType) {

	// private access is not allowed outside the file from a different
	// enclosing class.
	// TODO(tbreisacher): Should we also include the filename where ownerType is defined?
	var readableTypeName string
	if ownerType == nil || ownerType.Equals(propRef.receiverType) {
		readableTypeName = propRef.GetReadableTypeNameOrDefault()
	} else {
		readableTypeName = ownerType.String()
	}
	c.compiler.Report(
		JSError_MakeWithNode(
			propRef.sourceNode, CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS, propRef.name, readableTypeName))
}

// Java: CheckAccessControls.checkProtectedPropertyVisibility(PropertyReference, ObjectType)
func (c *CheckAccessControls) CheckProtectedPropertyVisibility(
	propRef *CheckAccessControls_PropertyReference, ownerType rhino.ObjectType) {
	// There are 3 types of legal accesses of a protected property:
	// 1) Accesses in the same file
	// 2) Overriding the property in a subclass
	// 3) Accessing the property from inside a subclass
	// The first two have already been checked for.
	if ownerType != nil {
		for i := len(c.currentClassStack) - 1; i >= 0; i-- {
			scopeType := c.currentClassStack[i]
			if scopeType == nil {
				continue
			} else if scopeType.IsSubtypeOf(ownerType) {
				return
			}
		}
	}

	c.compiler.Report(
		JSError_MakeWithNode(
			propRef.sourceNode,
			CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS,
			propRef.name,
			propRef.GetReadableTypeNameOrDefault()))
}

// ShouldEmitDeprecationWarningForPropRef determines whether a deprecation
// warning should be emitted.
//
// Java: CheckAccessControls.shouldEmitDeprecationWarning(NodeTraversal, PropertyReference)
func (c *CheckAccessControls) ShouldEmitDeprecationWarningForPropRef(t *NodeTraversal, propRef *CheckAccessControls_PropertyReference) bool {
	// In the global scope, there are only two kinds of accesses that should
	// be flagged for warnings:
	// 1) Calls of deprecated functions and methods.
	// 2) Instantiations of deprecated classes.
	// For now, we just let everything else by.
	if t.InGlobalScope() {
		if !NodeUtil_IsInvocationTarget(propRef.sourceNode) {
			return false
		}
	}

	// We can always assign to a deprecated property, to keep it up to date.
	if propRef.mutation {
		return false
	}

	// Don't warn if the node is just declaring the property, not reading it.
	if propRef.declaration {
		jsdoc := propRef.GetJSDocInfo()
		if jsdoc != nil && jsdoc.IsDeprecated() {
			return false
		}
	}

	return !c.CanAccessDeprecatedTypes(t)
}

// IsPropertyDeclaredConstant returns if a property is declared constant.
//
// Java: CheckAccessControls.isPropertyDeclaredConstant(ObjectType, String)
func (c *CheckAccessControls) IsPropertyDeclaredConstant(objectType rhino.ObjectType, prop string) CheckAccessControls_Constancy {
	for ; objectType != nil; objectType = objectType.GetImplicitPrototype() {
		docInfo := objectType.GetOwnPropertyJSDocInfo(prop)
		if docInfo == nil {
			continue
		} else if docInfo.IsFinal() {
			return CheckAccessControls_Constancy_FINAL
		} else if docInfo.IsConstant() {
			return CheckAccessControls_Constancy_OTHER_CONSTANT
		}
	}
	return CheckAccessControls_Constancy_MUTABLE
}

// CheckAccessControls_GetPropertyDeprecationInfo returns the deprecation
// reason for the property if it is marked as being deprecated. Returns ("", true)
// if the property is deprecated without a reason, or ("", false) if not deprecated.
//
// Java: CheckAccessControls.getPropertyDeprecationInfo(ObjectType, String)
func CheckAccessControls_GetPropertyDeprecationInfo(type_ rhino.ObjectType, prop string) (string, bool) {
	for ; type_ != nil; type_ = type_.GetImplicitPrototype() {
		if info := type_.GetOwnPropertyJSDocInfo(prop); info != nil {
			if depReason, ok := CheckAccessControls_GetDeprecationReason(info); ok {
				return depReason, true
			}
		}
	}
	return "", false
}

// CheckAccessControls_PropertyReference is a representation of an object
// property reference in JS code.
//
// This is an abstraction to smooth over the various AST structures that can
// act on properties. It is not useful for names (variables) or anonymous JS
// constructs.
//
// This class should only be used within CheckAccessControls. Having
// package-private visibility is a quirk of AutoValue.
//
// sourceNode: the Node that spawned this reference. receiverType: the type
// from which the property is referenced, not necessarily the one that
// declared it. readableTypeName: a lazy source for a human-readable type name
// to use when generating messages. Most users probably want
// GetReadableTypeNameOrDefault().
//
// Java: CheckAccessControls.PropertyReference (a record)
type CheckAccessControls_PropertyReference struct {
	sourceNode           rhino.Node
	name                 string
	receiverType         rhino.ObjectType
	mutation             bool
	declaration          bool
	override             bool
	readableInstanceType rhino.ObjectType
	typeRegistry         *rhino.JSTypeRegistry
	readableTypeNode     rhino.Node
	jsdoc                *rhino.JSDocInfo
	hasJSDoc             bool
}

// Derived properties.

// Java: CheckAccessControls.PropertyReference.getJSType()
func (r *CheckAccessControls_PropertyReference) GetJSType() rhino.JSType {
	return r.sourceNode.GetJSType()
}

// Java: CheckAccessControls.PropertyReference.getJSDocInfo()
func (r *CheckAccessControls_PropertyReference) GetJSDocInfo() *rhino.JSDocInfo {
	if r.hasJSDoc {
		return r.jsdoc
	}
	r.jsdoc = NodeUtil_GetBestJSDocInfo(r.sourceNode)
	r.hasJSDoc = true
	return r.jsdoc
}

// Java: CheckAccessControls.PropertyReference.isDocumentedDeclaration()
func (r *CheckAccessControls_PropertyReference) IsDocumentedDeclaration() bool {
	return r.declaration && (r.GetJSDocInfo() != nil)
}

// Java: CheckAccessControls.PropertyReference.isDeletion()
func (r *CheckAccessControls_PropertyReference) IsDeletion() bool {
	return r.sourceNode.GetParent().IsDelProp()
}

// Java: CheckAccessControls.PropertyReference.getReadableTypeNameOrDefault()
func (r *CheckAccessControls_PropertyReference) GetReadableTypeNameOrDefault() string {
	var preferred string
	if r.readableInstanceType != nil {
		preferred = r.readableInstanceType.String()
	} else if !r.readableTypeNode.IsNil() && r.typeRegistry != nil {
		preferred = r.typeRegistry.GetReadableTypeName(r.readableTypeNode)
	}
	if preferred == "" {
		return r.receiverType.String()
	}
	return preferred
}

// Java: CheckAccessControls.createPropertyReference(Node)
func (c *CheckAccessControls) CreatePropertyReference(sourceNode rhino.Node) *CheckAccessControls_PropertyReference {
	switch sourceNode.GetToken() {
	case rhino.Token_GETPROP:
		parent := sourceNode.GetParent()
		isLValue := NodeUtil_IsLValue(sourceNode)
		firstChild := sourceNode.GetFirstChild()
		ref := &c.scratchPropRef
		*ref = CheckAccessControls_PropertyReference{
			sourceNode:       sourceNode,
			name:             sourceNode.GetString(),
			receiverType:     c.BoxedOrUnknown(firstChild.GetJSType()),
			mutation:         isLValue || parent.IsDelProp(),
			typeRegistry:     c.typeRegistry,
			readableTypeNode: firstChild,
		}
		if isLValue {
			jsdoc := NodeUtil_GetBestJSDocInfo(sourceNode)
			ref.declaration = parent.IsExprResult() || (jsdoc != nil && jsdoc.IsConstant())
			// TODO(b/113704668): This definition is way too loose. It was used to prevent
			// breakages during refactoring and should be tightened.
			ref.override = jsdoc != nil
			ref.jsdoc = jsdoc
			ref.hasJSDoc = true
		} else {
			ref.declaration = parent.IsExprResult()
		}
		return ref

	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF:
		parent := sourceNode.GetParent()
		switch parent.GetToken() {
		case rhino.Token_OBJECTLIT:
			// TODO(b/80580110): Eventually object-literal members should be covered by
			// `PropertyReference`s. However, doing so initially would have caused too many errors
			// in existing code and delayed support for class syntax.
			parentType := parent.GetJSType()
			if !parentType.IsLiteralObject() {
				// Only add a mutation if the object type is actually a literal object (e.g. a
				// global namespace).  OBJECTLIT tokens are often used to fulfill structural types,
				// which is fine for writing constant properties the first time.
				return nil
			}
			ref := &c.scratchPropRef
			*ref = CheckAccessControls_PropertyReference{
				sourceNode:       sourceNode,
				name:             sourceNode.GetString(),
				receiverType:     c.TypeOrUnknownObjectType(rhino.ObjectType_Cast(parentType)),
				mutation:         true,
				declaration:      true,
				typeRegistry:     c.typeRegistry,
				readableTypeNode: parent,
			}
			return ref

		case rhino.Token_OBJECT_PATTERN:
			ref := &c.scratchPropRef
			*ref = CheckAccessControls_PropertyReference{
				sourceNode:       sourceNode,
				name:             sourceNode.GetString(),
				receiverType:     c.TypeOrUnknownObjectType(rhino.ObjectType_Cast(parent.GetJSType())),
				typeRegistry:     c.typeRegistry,
				readableTypeNode: parent,
			}
			return ref

		case rhino.Token_CLASS_MEMBERS:
			jsdoc := NodeUtil_GetBestJSDocInfo(sourceNode)
			receiverType := c.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
			ctorType := parent.GetParent().GetJSType()
			if ctorType != nil && ctorType.IsFunctionType() {
				ctorFunctionType := ctorType.ToMaybeFunctionType()
				if sourceNode.IsStaticMember() {
					receiverType = ctorFunctionType
				} else if sourceNode.IsMemberFieldDef() {
					receiverType = ctorFunctionType.GetInstanceType()
				} else {
					receiverType = ctorFunctionType.GetPrototype()
				}
			}
			ref := &c.scratchPropRef
			// TODO(b/113704668): override = jsdoc != nil is way too loose. It was used to prevent
			// breakages during refactoring and should be tightened.
			*ref = CheckAccessControls_PropertyReference{
				sourceNode:   sourceNode,
				name:         sourceNode.GetString(),
				receiverType: receiverType,
				mutation:     !(sourceNode.IsMemberFieldDef() && !sourceNode.HasChildren()),
				declaration:  true,
				override:     jsdoc != nil,
				jsdoc:        jsdoc,
				hasJSDoc:     true,
			}
			return ref

		default:
			panic(jutil.NewAssertionError("null"))
		}
	default:
		return nil
	}
}

// CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty returns
// the effective visibility of the given non-overridden property.
// Non-overridden properties without an explicit visibility annotation receive
// the default visibility declared in the file's @fileoverview block, if one
// exists.
//
// TODO(b/111789692): The following methods are forked from
// `AccessControlUtils`. Consider consolidating them.
//
// Java: CheckAccessControls.getEffectiveVisibilityForNonOverriddenProperty(PropertyReference, ObjectType, Visibility)
func CheckAccessControls_GetEffectiveVisibilityForNonOverriddenProperty(
	propRef *CheckAccessControls_PropertyReference,
	objectType rhino.ObjectType,
	fileOverviewVisibility *rhino.JSDocInfo_Visibility) rhino.JSDocInfo_Visibility {
	propertyName := propRef.name
	raw := rhino.JSDocInfo_Visibility_INHERITED
	if objectType != nil {
		jsdoc := objectType.GetOwnPropertyJSDocInfo(propertyName)
		if jsdoc != nil {
			raw = jsdoc.GetVisibility()
		}
	}
	type_ := propRef.GetJSType()
	createdFromGoogProvide := (type_ != nil && type_.IsLiteralObject())
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
	if raw != rhino.JSDocInfo_Visibility_INHERITED || fileOverviewVisibility == nil || createdFromGoogProvide {
		return raw
	}
	return *fileOverviewVisibility
}
