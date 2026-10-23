// Port of com/google/javascript/rhino/jstype/NamedType.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// NamedType is a named reference to some other type. This provides a
// convenient mechanism for implementing forward references to types; a
// NamedType can be used as a placeholder until its reference is resolved. It
// is also useful for representing type names in JsDoc type annotations,
// some of which may never be resolved (as they may refer to types in host
// systems not yet supported by JSCompiler, such as the JVM.)
//
// An important distinction: NamedType is a type name reference, whereas
// ObjectType is a named type object, such as an Enum name. The Enum itself
// is typically used only in a dot operator to name one of its constants, or
// in a declaration, where its name will appear in a NamedType.
//
// A NamedType is not currently a full-fledged typedef, because it cannot
// resolve to any JavaScript type. It can only resolve to a named
// JSTypeRegistry type, or to FunctionType or EnumType.
//
// If full typedefs are to be supported, then each method on each type class
// needs to be reviewed to make sure that everything works correctly through
// typedefs. Alternatively, we would need to walk through the parse tree and
// unroll each reference to a NamedType to its resolved type before applying
// the rest of the analysis.
//
// TODO(user): Revisit all of this logic.
//
// The existing typing logic is hacky. Unresolved types should get processed
// in a more consistent way, but with the Rhino merge coming, there will be
// much that has to be changed.
type NamedType struct {
	ProxyObjectTypeBase

	reference       string
	sourceName      *string
	resolutionScope StaticTypedScope

	// Validates the type resolution.
	validator func(JSType) bool

	// The following instance properties (`propertyContinuations`,
	// `templateTypes`, and `restrictByNull`) are used to indicate that some
	// type operation should be applied to the type after resolution. This is
	// necessary because type operations are not well-defined when applied to
	// unresolved NamedTypes.
	//
	// TODO(lharker): Generalize this pattern instead of storing these
	// arbitrary fields.

	// Property-defining continuations.
	propertyContinuations []*NamedType_PropertyContinuation

	// Template types defined on a named, not yet resolved type, or null if
	// none. These are ignored during resolution, for backwards compatibility
	// with existing usage. This field is not used for JSCompiler's type
	// checking; it is only needed by Clutz.
	templateTypes []JSType

	lineno         int32
	charno         int32
	resolutionKind NamedType_ResolutionKind

	// Applies the "!" operator to the resolved type, which removes null and
	// undefined
	restrictByNull bool
}

// Java: NamedType.TYPE_CLASS
const NamedType_TYPE_CLASS = JSTypeClass_NAMED

// Java: NamedType.nominalHashCode(ObjectType)
func NamedType_NominalHashCode(type_ ObjectType) int32 {
	jutil.CheckState(type_.HasReferenceName())
	name := jutil.CheckNotNull(type_.GetReferenceName())
	return jutil.StringHashCode(*name)
}

// Java: NamedType(Builder)
func NewNamedType(builder *NamedType_Builder) *NamedType {
	t := builder.registry.allocNamedType()
	t.InitProxyObjectTypeBase(t, builder.registry, builder.referencedType)
	jutil.CheckState(builder.hasReferenceName)
	jutil.CheckState(builder.hasResolutionKind)
	jutil.CheckNotNull(builder.templateTypes)
	if builder.resolutionKind == NamedType_ResolutionKind_TYPEOF {
		jutil.CheckState(jutil.StringStartsWith(builder.referenceName, "typeof "))
	}
	// TODO(lharker): enforce that the scope is not null

	t.restrictByNull = builder.restrictByNull
	t.resolutionScope = builder.scope
	t.reference = builder.referenceName
	t.sourceName = builder.sourceName
	t.lineno = int32(builder.lineno)
	t.charno = int32(builder.charno)
	t.templateTypes = builder.templateTypes
	t.resolutionKind = builder.resolutionKind

	t.registry.GetResolver().ResolveIfClosed(t, NamedType_TYPE_CLASS)
	return t
}

// Java: NamedType.getTypeClass()
func (t *NamedType) GetTypeClass() JSTypeClass {
	return NamedType_TYPE_CLASS
}

// GetBangType returns a new non-null version of this type.
//
// Java: NamedType.getBangType()
func (t *NamedType) GetBangType() JSType {
	if t.restrictByNull {
		return t
	} else if t.IsResolved() {
		// Already resolved, just restrict.
		// TODO(b/146173738): just return
		// getReferencedType().restrictByNotNullOrUndefined() after fixing how
		// conformance checks handle unresolved types.
		if t.IsNoResolvedType() || t.IsUnknownType() {
			return t
		}
		return t.GetReferencedType().RestrictByNotNullOrUndefined()
	}
	return t.ToBuilder().SetRestrictByNull(true).Build()
}

// Java: NamedType.getTemplateTypes()
func (t *NamedType) GetTemplateTypes() []JSType {
	return t.templateTypes
}

// Java: NamedType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *NamedType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool {
	if !t.IsResolved() {
		// If this is an unresolved object type, we need to save all its
		// properties and define them when it is resolved.
		if t.propertyContinuations == nil {
			t.propertyContinuations = []*NamedType_PropertyContinuation{}
		}
		t.propertyContinuations = append(t.propertyContinuations,
			NewNamedType_PropertyContinuation(propertyName, type_, inferred, propertyNode))
		return true
	} else {
		return t.ProxyObjectTypeBase.DefinePropertyByKey(propertyName, type_, inferred, propertyNode)
	}
}

// Java: NamedType.finishPropertyContinuations()
func (t *NamedType) FinishPropertyContinuations() {
	referencedObjType := t.GetReferencedObjTypeInternal()
	if referencedObjType != nil &&
		!referencedObjType.IsUnknownType() &&
		t.propertyContinuations != nil {
		for _, c := range t.propertyContinuations {
			c.Commit(t)
		}
	}
	t.propertyContinuations = nil
}

// GetReferencedType returns the type to which this refers (which is unknown
// if unresolved).
//
// Java: NamedType.getReferencedType()
func (t *NamedType) GetReferencedType() JSType {
	return t.GetReferencedTypeInternal()
}

// Java: NamedType.getReferenceName()
func (t *NamedType) GetReferenceName() *string {
	return &t.reference
}

// Java: NamedType.appendTo(TypeStringBuilder)
func (t *NamedType) AppendTo(sb *TypeStringBuilder) {
	type_ := t.GetReferencedType()
	if !t.IsResolved() || type_.IsNoResolvedType() {
		sb.Append(*t.GetReferenceName())
	} else {
		sb.AppendType(type_)
	}
}

// Java: JSType.isNamedType()
func (t *NamedType) IsNamedType() bool {
	return true
}

// Java: NamedType.toMaybeNamedType()
func (t *NamedType) ToMaybeNamedType() *NamedType {
	return t
}

// Java: NamedType.isNominalType()
func (t *NamedType) IsNominalType() bool {
	if t.IsResolved() {
		return t.ProxyObjectTypeBase.IsNominalType()
	}
	return true
}

// Java: NamedType.recursionUnsafeHashCode()
func (t *NamedType) RecursionUnsafeHashCode() int32 {
	// Recall that equality on `NamedType` uses only the name until successful
	// resolution, then delegates to the resolved type.
	if t.IsSuccessfullyResolved() {
		return t.ProxyObjectTypeBase.RecursionUnsafeHashCode()
	}
	return NamedType_NominalHashCode(t)
}

// ResolveInternal resolves the referenced type within the enclosing scope.
//
// Java: NamedType.resolveInternal(ErrorReporter)
func (t *NamedType) ResolveInternal(reporter ErrorReporter) JSType {
	var resolvedTypeArgs []JSType
	if len(t.templateTypes) != 0 {
		resolvedTypeArgs = make([]JSType, len(t.templateTypes))
		for i, x := range t.templateTypes {
			resolvedTypeArgs[i] = x.Resolve(reporter)
		}
	}

	if t.resolutionKind == NamedType_ResolutionKind_NONE {
		// In some cases (e.g. typeof(ns) when the actual type is just a literal
		// object), a NamedType is created solely for the purpose of naming an
		// already-known type. When that happens, there's nothing to look up, so
		// just resolve the referenced type.
		t.resolutionScope = nil
		return t.ProxyObjectTypeBase.ResolveInternal(reporter)
	}
	jutil.CheckStateMsg(
		t.GetReferencedType().IsUnknownType(),
		"NamedTypes given a referenced type pre-resolution should have ResolutionKind.NONE")

	if t.resolutionScope == nil {
		return t
	}

	// TODO(user): Investigate whether it is really necessary to keep two
	// different mechanisms for resolving named types, and if so, which order
	// makes more sense.
	unused := t.ResolveTypeof(reporter) || t.ResolveViaRegistry(reporter)
	_ = unused

	t.ProxyObjectTypeBase.ResolveInternal(reporter)
	if t.DetectInheritanceCycle() {
		t.HandleTypeCycle(reporter)
	}
	t.FinishPropertyContinuations()
	t.resolutionScope = nil

	result := t.GetReferencedType()
	if t.IsSuccessfullyResolved() {
		resultAsObject := result.ToMaybeObjectType()

		if resultAsObject == nil {
			// For non-object types there is no need to handle template parameters
			// or interface aliases, so we can just return now.
			return result
		}

		if len(resolvedTypeArgs) == 0 || !resultAsObject.IsRawTypeOfTemplatizedType() {
			// No template parameters need to be resolved.
			return result
		}

		numKeys := result.GetTemplateParamCount()
		// TODO(b/287880204): report an error if there are too many type
		// arguments
		if numKeys < len(resolvedTypeArgs) {
			resolvedTypeArgs = jsTypeSubList(resolvedTypeArgs, 0, numKeys)
		}

		result = t.registry.CreateTemplatizedTypeOfList(resultAsObject, resolvedTypeArgs)
		t.SetReferencedType(result.Resolve(reporter))
	}

	return result
}

// ResolveViaRegistry resolves a named type by looking it up in the
// registry. Returns true if we resolved successfully.
//
// Java: NamedType.resolveViaRegistry(ErrorReporter)
func (t *NamedType) ResolveViaRegistry(reporter ErrorReporter) bool {
	type_ := t.registry.GetType(t.resolutionScope, t.reference)
	if type_ == nil {
		t.HandleUnresolvedType(reporter)
		return false
	}
	t.SetReferencedAndResolvedType(type_, reporter)
	return true
}

// Java: NamedType.resolveTypeof(ErrorReporter)
func (t *NamedType) ResolveTypeof(reporter ErrorReporter) bool {
	if t.resolutionKind != NamedType_ResolutionKind_TYPEOF {
		return false
	}

	// TODO(sdh): require var to be const?
	scopeName := jutil.StringSubstringFrom(t.reference, jutil.StringLength("typeof "))
	type_ := t.resolutionScope.LookupQualifiedName(QualifiedName_Of(scopeName))
	if type_ == nil || type_.IsUnknownType() {
		if t.registry.IsForwardDeclaredType(scopeName) {
			// Preserve the "typeof" as a `NoResolvedType`. This is depended on by
			// Clutz so it can generate `typeof ImportedType` instead of `any` when
			// `ImportedType` is not defined in the files it can see.
			t.SetReferencedType(NewNoResolvedType(t.registry, *t.GetReferenceName(), t.GetTemplateTypes()))
			if t.validator != nil {
				unused := t.validator(t.GetReferencedType())
				_ = unused
			}
		} else {
			t.Warning(reporter, "Missing type for `typeof` value. The value must be declared and const.")
			t.SetReferencedAndResolvedType(t.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE), reporter)
		}
	} else {
		if type_.IsLiteralObject() {
			// Create an extra layer of wrapping so that the "typeof" name is
			// preserved for namespaces. This is depended on by Clutz to prevent
			// infinite loops in self-referential typeof types.
			objlit := type_
			type_ = NamedType_Builder_(t.registry, *t.GetReferenceName()).
				SetResolutionKind(NamedType_ResolutionKind_NONE).
				SetReferencedType(objlit).
				Build()
		}
		t.SetReferencedAndResolvedType(type_, reporter)
	}

	return true
}

// Java: NamedType.setReferencedAndResolvedType(JSType, ErrorReporter)
func (t *NamedType) SetReferencedAndResolvedType(type_ JSType, reporter ErrorReporter) {
	if t.restrictByNull {
		type_ = type_.RestrictByNotNullOrUndefined()
	}
	if t.validator != nil {
		unused := t.validator(type_)
		_ = unused
	}
	t.SetReferencedType(type_)
	t.CheckEnumElementCycle(reporter)
	t.CheckProtoCycle(reporter)
}

// Java: NamedType.handleTypeCycle(ErrorReporter)
func (t *NamedType) HandleTypeCycle(reporter ErrorReporter) {
	t.SetReferencedType(t.registry.GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE))
	t.Warning(reporter, "Cycle detected in inheritance chain of type "+t.reference)
}

// Java: NamedType.checkEnumElementCycle(ErrorReporter)
func (t *NamedType) CheckEnumElementCycle(reporter ErrorReporter) {
	referencedType := t.GetReferencedType()
	if enumElementType, ok := referencedType.(*EnumElementType); ok &&
		JSType(t) == enumElementType.GetPrimitiveType() {
		t.HandleTypeCycle(reporter)
	}
}

// Java: NamedType.checkProtoCycle(ErrorReporter)
func (t *NamedType) CheckProtoCycle(reporter ErrorReporter) {
	referencedType := t.GetReferencedType()
	if referencedType == JSType(t) {
		t.HandleTypeCycle(reporter)
	}
}

// HandleUnresolvedType warns about this type being unresolved iff it's not a
// forward-declared type name.
//
// Java: NamedType.handleUnresolvedType(ErrorReporter)
func (t *NamedType) HandleUnresolvedType(reporter ErrorReporter) {
	isForwardDeclared := t.registry.IsForwardDeclaredType(t.reference)
	if !isForwardDeclared {
		msg := "Bad type annotation. Unknown type " + t.reference
		// Look for a local variable that shadows a global namespace to give a
		// clearer message.
		root := t.reference
		if jutil.StringIndexOf(t.reference, ".") >= 0 {
			root = jutil.StringSubstring(t.reference, 0, jutil.StringIndexOf(t.reference, "."))
		}
		if t.LocalVariableShadowsGlobalNamespace(root) {
			msg += "\nIt's possible that a local variable called '" + root +
				"' is shadowing the intended global namespace."
		} else if t.resolutionScope.GetSlot(root) != nil {
			msg += "\nIt's possible that '" + t.reference + "' refers to a value, not a type."
		}
		t.Warning(reporter, msg)
	} else {
		t.SetReferencedType(NewNoResolvedType(t.registry, *t.GetReferenceName(), t.GetTemplateTypes()))
		if t.validator != nil {
			unused := t.validator(t.GetReferencedType())
			_ = unused
		}
	}
}

// LocalVariableShadowsGlobalNamespace checks for an obscure but very
// confusing error condition where a local variable shadows a global
// namespace.
//
// Java: NamedType.localVariableShadowsGlobalNamespace(String)
func (t *NamedType) LocalVariableShadowsGlobalNamespace(root string) bool {
	rootVar := t.resolutionScope.GetSlot(root)
	if rootVar != nil {
		jutil.CheckNotNullMsg1(rootVar.GetScope(), "%s", rootVar)
		parent := rootVar.GetScope().GetParentScope()
		if parent != nil {
			globalVar := parent.GetSlot(root)
			return globalVar != nil
		}
	}
	return false
}

// Java: NamedType.setValidator(Predicate<JSType>)
func (t *NamedType) SetValidator(validator func(JSType) bool) bool {
	// If the type is already resolved, we can validate it now. If the type has
	// not been resolved yet, we need to wait till its resolved before we can
	// validate it.
	if t.IsResolved() {
		return t.ProxyObjectTypeBase.SetValidator(validator)
	} else {
		t.validator = validator
		return true
	}
}

// Java: NamedType.warning(ErrorReporter, String)
func (t *NamedType) Warning(reporter ErrorReporter, message string) {
	reporter.Warning(message, jsTypeDerefString(t.sourceName), int(t.lineno), int(t.charno))
}

// NamedType_PropertyContinuation stores enough information to define a
// property at a later time.
type NamedType_PropertyContinuation struct {
	propertyName Property_Key
	type_        JSType
	inferred     bool
	propertyNode *Node
}

// Java: NamedType.PropertyContinuation(Property.Key, JSType, boolean, Node)
func NewNamedType_PropertyContinuation(
	propertyName Property_Key, type_ JSType, inferred bool, propertyNode *Node) *NamedType_PropertyContinuation {
	return &NamedType_PropertyContinuation{
		propertyName: propertyName,
		type_:        type_,
		inferred:     inferred,
		propertyNode: propertyNode,
	}
}

// Java: NamedType.PropertyContinuation.commit(ObjectType)
func (c *NamedType_PropertyContinuation) Commit(target ObjectType) {
	target.DefinePropertyByKey(c.propertyName, c.type_, c.inferred, c.propertyNode)
}

// Java: NamedType.isObject()
func (t *NamedType) IsObject() bool {
	if t.IsEnumElementType() {
		return t.ToMaybeEnumElementType().IsObject()
	}
	return t.ProxyObjectTypeBase.IsObject()
}

// Java: JSType.isObjectType()
func (t *NamedType) IsObjectType() bool {
	return t.IsObject()
}

// Java: NamedType.visit(Visitor<T>)
func (t *NamedType) Visit(visitor Visitor) any {
	return visitor.CaseNamedType(t)
}

// Java: NamedType.builder(JSTypeRegistry, String)
func NamedType_Builder_(registry *JSTypeRegistry, reference string) *NamedType_Builder {
	return NewNamedType_Builder(registry, reference)
}

// NamedType_ResolutionKind is NamedType.ResolutionKind.
type NamedType_ResolutionKind int32

const (
	NamedType_ResolutionKind_NONE NamedType_ResolutionKind = iota
	NamedType_ResolutionKind_TYPE_NAME
	NamedType_ResolutionKind_TYPEOF
)

var namedType_ResolutionKind_names = [...]string{"NONE", "TYPE_NAME", "TYPEOF"}

// NamedType_ResolutionKind_values is ResolutionKind.values().
func NamedType_ResolutionKind_values() []NamedType_ResolutionKind {
	return []NamedType_ResolutionKind{
		NamedType_ResolutionKind_NONE, NamedType_ResolutionKind_TYPE_NAME, NamedType_ResolutionKind_TYPEOF,
	}
}

// String is Enum.name().
func (k NamedType_ResolutionKind) String() string { return namedType_ResolutionKind_names[k] }

// Java: NamedType.toBuilder()
func (t *NamedType) ToBuilder() *NamedType_Builder {
	jutil.CheckStateMsg(!t.IsResolved(), "Only call toBuilder on unresolved NamedTypes")
	return NewNamedType_Builder(t.registry, t.reference).
		SetScope(t.resolutionScope).
		SetResolutionKind(t.resolutionKind).
		SetErrorReportingLocation(t.sourceName, int(t.lineno), int(t.charno)).
		SetTemplateTypes(t.templateTypes).
		SetReferencedType(t.GetReferencedType()).
		SetRestrictByNull(t.restrictByNull)
}

// NamedType_Builder is NamedType.Builder.
type NamedType_Builder struct {
	registry          *JSTypeRegistry
	referenceName     string
	scope             StaticTypedScope
	sourceName        *string
	lineno            int
	charno            int
	referencedType    JSType
	templateTypes     []JSType
	resolutionKind    NamedType_ResolutionKind
	hasResolutionKind bool
	hasReferenceName  bool
	restrictByNull    bool
}

// Java: NamedType.Builder(JSTypeRegistry, String)
func NewNamedType_Builder(registry *JSTypeRegistry, referenceName string) *NamedType_Builder {
	return &NamedType_Builder{
		registry:         registry,
		referenceName:    referenceName,
		hasReferenceName: true,
		referencedType:   registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE),
		templateTypes:    emptyJSTypeSlice,
	}
}

// Java: NamedType.Builder.setScope(StaticTypedScope)
func (b *NamedType_Builder) SetScope(scope StaticTypedScope) *NamedType_Builder {
	b.scope = scope
	return b
}

// Java: NamedType.Builder.setResolutionKind(ResolutionKind)
func (b *NamedType_Builder) SetResolutionKind(resolutionKind NamedType_ResolutionKind) *NamedType_Builder {
	b.resolutionKind = resolutionKind
	b.hasResolutionKind = true
	return b
}

// Java: NamedType.Builder.setErrorReportingLocation(String, int, int)
func (b *NamedType_Builder) SetErrorReportingLocation(sourceName *string, lineno int, charno int) *NamedType_Builder {
	b.sourceName = sourceName
	b.lineno = lineno
	b.charno = charno
	return b
}

// Java: NamedType.Builder.setErrorReportingLocationFrom(Node)
func (b *NamedType_Builder) SetErrorReportingLocationFrom(source *Node) *NamedType_Builder {
	b.sourceName = source.GetSourceFileName()
	b.lineno = source.GetLineno()
	b.charno = source.GetCharno()
	return b
}

// Java: NamedType.Builder.setTemplateTypes(ImmutableList<JSType>)
func (b *NamedType_Builder) SetTemplateTypes(templateTypes []JSType) *NamedType_Builder {
	b.templateTypes = templateTypes
	return b
}

// Java: NamedType.Builder.setReferencedType(JSType)
func (b *NamedType_Builder) SetReferencedType(referencedType JSType) *NamedType_Builder {
	b.referencedType = referencedType
	return b
}

// Java: NamedType.Builder.setRestrictByNull(boolean)
func (b *NamedType_Builder) SetRestrictByNull(restrictByNull bool) *NamedType_Builder {
	b.restrictByNull = restrictByNull
	return b
}

// Java: NamedType.Builder.build()
func (b *NamedType_Builder) Build() *NamedType {
	return NewNamedType(b)
}
