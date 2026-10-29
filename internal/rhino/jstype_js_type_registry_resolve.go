// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: JSTypeRegistry.DOT_SPLITTER
func JSTypeRegistry_DOT_SPLITTER_Split(s string) []string {
	return strings.Split(s, ".")
}

// GetArrayElementKey returns the template variable for the element type of
// Arrays.
//
// Java: JSTypeRegistry.getArrayElementKey()
func (r *JSTypeRegistry) GetArrayElementKey() *TemplateType {
	return r.arrayElementTemplateKey
}

// GetReadonlyArrayElementKey returns the template variable for the element
// type of ReadonlyArrays.
//
// Java: JSTypeRegistry.getReadonlyArrayElementKey()
func (r *JSTypeRegistry) GetReadonlyArrayElementKey() *TemplateType {
	return r.readonlyArrayElementTemplateKey
}

// GetReadonlyMapKey returns the template variable for the key type of
// ReadonlyMaps.
//
// Java: JSTypeRegistry.getReadonlyMapKey()
func (r *JSTypeRegistry) GetReadonlyMapKey() *TemplateType {
	return r.readonlyMapKeyTemplateKey
}

// GetReadonlyMapValue returns the template variable for the value type of
// ReadonlyMaps.
//
// Java: JSTypeRegistry.getReadonlyMapValue()
func (r *JSTypeRegistry) GetReadonlyMapValue() *TemplateType {
	return r.readonlyMapValueTemplateKey
}

// GetObjectElementKey returns the template variable corresponding to the
// property value type for Javascript Objects and Arrays.
//
// Java: JSTypeRegistry.getObjectElementKey()
func (r *JSTypeRegistry) GetObjectElementKey() *TemplateType {
	return r.iObjectElementTemplateKey
}

// GetObjectIndexKey returns the template variable corresponding to the
// property key type of the built-in Javascript object.
//
// Java: JSTypeRegistry.getObjectIndexKey()
func (r *JSTypeRegistry) GetObjectIndexKey() *TemplateType {
	jutil.CheckNotNull(r.iObjectIndexTemplateKey)
	return r.iObjectIndexTemplateKey
}

// GetIterableValueTemplate returns the value template variable for the
// Iterable interface.
//
// Java: JSTypeRegistry.getIterableValueTemplate()
func (r *JSTypeRegistry) GetIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iterableValueTemplate)
}

// GetIterableReturnTemplate returns the return template variable for the
// Iterable interface.
//
// Java: JSTypeRegistry.getIterableReturnTemplate()
func (r *JSTypeRegistry) GetIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iterableReturnTemplate)
}

// GetIteratorIterableValueTemplate returns the value template variable for
// the IteratorIterable interface.
//
// Java: JSTypeRegistry.getIteratorIterableValueTemplate()
func (r *JSTypeRegistry) GetIteratorIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorIterableValueTemplate)
}

// GetIteratorIterableReturnTemplate returns the return template variable
// for the IteratorIterable interface.
//
// Java: JSTypeRegistry.getIteratorIterableReturnTemplate()
func (r *JSTypeRegistry) GetIteratorIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorIterableReturnTemplate)
}

// GetIIterableResultValueTemplate returns the value template variable for
// the IIterableResult interface.
//
// Java: JSTypeRegistry.getIIterableResultValueTemplate()
func (r *JSTypeRegistry) GetIIterableResultValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iiterableResultValueTemplate)
}

// GetIteratorLikeValueTemplate returns the value template variable for the
// IteratorLike interface.
//
// Java: JSTypeRegistry.getIteratorLikeValueTemplate()
func (r *JSTypeRegistry) GetIteratorLikeValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorLikeValueTemplate)
}

// GetIteratorLikeReturnTemplate returns the return template variable for the
// IteratorLike interface.
//
// Java: JSTypeRegistry.getIteratorLikeReturnTemplate()
func (r *JSTypeRegistry) GetIteratorLikeReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorLikeReturnTemplate)
}

// GetIteratorValueTemplate returns the value template variable for the
// Iterator interface.
//
// Java: JSTypeRegistry.getIteratorValueTemplate()
func (r *JSTypeRegistry) GetIteratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorValueTemplate)
}

// GetIteratorReturnTemplate returns the return template variable for the
// Iterator interface.
//
// Java: JSTypeRegistry.getIteratorReturnTemplate()
func (r *JSTypeRegistry) GetIteratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorReturnTemplate)
}

// GetGeneratorValueTemplate returns the value template variable for the
// Generator interface.
//
// Java: JSTypeRegistry.getGeneratorValueTemplate()
func (r *JSTypeRegistry) GetGeneratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.generatorValueTemplate)
}

// GetGeneratorReturnTemplate returns the return template variable for the
// Generator interface.
//
// Java: JSTypeRegistry.getGeneratorReturnTemplate()
func (r *JSTypeRegistry) GetGeneratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.generatorReturnTemplate)
}

// GetAsyncIterableValueTemplate returns the template variable for the
// AsyncIterable interface.
//
// Java: JSTypeRegistry.getAsyncIterableValueTemplate()
func (r *JSTypeRegistry) GetAsyncIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIterableValueTemplate)
}

// GetAsyncIterableReturnTemplate returns the template variable for the
// AsyncIterableR interface.
//
// Java: JSTypeRegistry.getAsyncIterableReturnTemplate()
func (r *JSTypeRegistry) GetAsyncIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIterableReturnTemplate)
}

// GetAsyncIteratorValueTemplate returns the template variable for the
// AsyncIterator interface.
//
// Java: JSTypeRegistry.getAsyncIteratorValueTemplate()
func (r *JSTypeRegistry) GetAsyncIteratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIteratorValueTemplate)
}

// GetAsyncIteratorReturnTemplate returns the template variable for the
// AsyncIterator interface.
//
// Java: JSTypeRegistry.getAsyncIteratorReturnTemplate()
func (r *JSTypeRegistry) GetAsyncIteratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIteratorReturnTemplate)
}

// GetIThenableTemplate returns the template variable for the IThenable
// interface.
//
// Java: JSTypeRegistry.getIThenableTemplate()
func (r *JSTypeRegistry) GetIThenableTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iThenableTemplateKey)
}

// MaybeGetTemplateTypesOfBuiltin returns an immutable list of template types
// of the given builtin.
//
// Java: JSTypeRegistry.maybeGetTemplateTypesOfBuiltin(StaticScope, String)
func (r *JSTypeRegistry) MaybeGetTemplateTypesOfBuiltin(scope StaticScope, fnName string) []*TemplateType {
	type_ := r.GetType(scope, fnName)
	var objType ObjectType
	if type_ != nil {
		objType = type_.ToObjectType()
	}
	if objType != nil && objType.IsNativeObjectType() {
		return objType.GetTypeParameters()
	}
	return nil
}

// Java: JSTypeRegistry.getRootElementOfName(String)
func JSTypeRegistry_GetRootElementOfName(name string) string {
	index := strings.IndexByte(name, '.')
	if index != -1 {
		return name[:index]
	}
	return name
}

// JSTypeRegistry_GetLookupScope returns which scope in the provided scope
// chain the provided name is declared in, or else null.
//
// This assumed that the Scope construction is complete. It can not be used
// during scope construction to determine if a name is already defined as a
// shadowed name from a parent scope would be returned.
//
// Java: JSTypeRegistry.getLookupScope(StaticScope, String)
func JSTypeRegistry_GetLookupScope(scope StaticScope, name string) StaticScope {
	if scope != nil && scope.GetParentScope() != nil {
		return scope.GetTopmostScopeOfEventualDeclaration(JSTypeRegistry_GetRootElementOfName(name))
	}
	return scope
}

// Java: JSTypeRegistry.getRootNodeForScope(StaticScope)
func (r *JSTypeRegistry) GetRootNodeForScope(scope StaticScope) Node {
	var root Node
	if scope != nil {
		root = scope.GetRootNode()
	}
	if root.IsNil() || root.IsRoot() || root.IsScript() {
		return r.nameTableGlobalRoot
	}
	return root
}

// Java: JSTypeRegistry.isDeclaredForScope(StaticScope, String)
func (r *JSTypeRegistry) IsDeclaredForScope(scope StaticScope, name string) bool {
	return r.GetTypeInternal(scope, name) != nil
}

// Java: JSTypeRegistry.checkTypeName(String)
func JSTypeRegistry_CheckTypeName(typeName string) {
	if strings.IndexByte(typeName, '<') >= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Type names cannot contain template annotations."))
	}
}

// Java: JSTypeRegistry.getTypeInternal(StaticScope, String)
func (r *JSTypeRegistry) GetTypeInternal(scope StaticScope, name string) JSType {
	JSTypeRegistry_CheckTypeName(name)
	if syntheticTemplateScope, ok := scope.(*JSTypeRegistry_SyntheticTemplateScope); ok {
		type_ := syntheticTemplateScope.GetTemplateType(name)
		if type_ != nil {
			return type_
		}
	}
	declarationScope := JSTypeRegistry_GetLookupScope(scope, name)
	resolvedViaTable := r.GetTypeForScopeInternal(declarationScope, name)
	if resolvedViaTable != nil {
		return resolvedViaTable
	}
	bestScope := scope
	if declarationScope != nil {
		bestScope = declarationScope
	}
	return r.ResolveViaComponents(bestScope, name)
}

// Java: JSTypeRegistry.resolveViaComponents(StaticScope, String)
func (r *JSTypeRegistry) ResolveViaComponents(scope StaticScope, qualifiedName string) JSType {
	resolutionScope, ok := scope.(StaticTypedScope)
	if qualifiedName == "" || !ok {
		return nil
	}
	// Skip closure namespace resolution of types whose root component is
	// defined in a local scope (not a global scope). Those will follow the
	// normal resolution scheme. (For legacy compatibility reasons we don't
	// check for global names that are the same as the module root).
	if !JSTypeRegistry_IsNameDefinedLocally(resolutionScope, JSTypeRegistry_GetRootElementOfName(qualifiedName)) {
		resolvedViaClosureNamespace := r.ResolveViaClosureNamespace(qualifiedName)
		if resolvedViaClosureNamespace != nil {
			return resolvedViaClosureNamespace
		}
	}
	return r.ResolveViaProperties(resolutionScope, qualifiedName)
}

// Java: JSTypeRegistry.isNameDefinedLocally(StaticTypedScope, String)
func JSTypeRegistry_IsNameDefinedLocally(resolutionScope StaticTypedScope, reference string) bool {
	slot := resolutionScope.GetSlot(reference)
	return slot != nil && slot.GetScope() != nil && slot.GetScope().GetParentScope() != nil
}

// ResolveViaClosureNamespace resolves a named type by checking for the
// longest prefix that matches some Closure namespace, if any, then
// attempting to resolve via properties based on the type of the `exports`
// object in that namespace.
//
// Java: JSTypeRegistry.resolveViaClosureNamespace(String)
func (r *JSTypeRegistry) ResolveViaClosureNamespace(reference string) JSType {
	if len(r.closureNamespaces) == 0 {
		return nil
	}
	// Find the `exports` type of the longest prefix match of this namespace,
	// if any. Then resolve it via property.
	prefix := reference
	for {
		namespace := r.closureNamespaces[prefix]
		if namespace != nil {
			if namespace.IsLegacy() {
				// Try to resolve this name via registry or properties.
				return nil
			}
			// Always stop resolution here whether successful or not, instead of
			// continuing with resolution via registry or via properties, to
			// match legacy behavior.
			var components []string
			if len(prefix) < len(reference) {
				components = strings.Split(reference[len(prefix)+1:], ".")
			}
			return r.ResolveViaPropertyGivenSlot(
				namespace.Type(), namespace.DefinitionNode(), components)
		}

		lastDot := strings.LastIndexByte(prefix, '.')
		if lastDot < 0 {
			return nil
		}
		prefix = prefix[:lastDot]
	}
}

// Java: JSTypeRegistry.resolveViaProperties(StaticTypedScope, String)
func (r *JSTypeRegistry) ResolveViaProperties(declarationScope StaticTypedScope, qualifiedName string) JSType {
	jutil.CheckArgument(qualifiedName != "")
	dotIdx := strings.IndexByte(qualifiedName, '.')
	rootName := qualifiedName
	if dotIdx != -1 {
		rootName = qualifiedName[:dotIdx]
	}
	slot, _ := declarationScope.GetOwnSlot(rootName).(StaticTypedSlot)
	if slot == nil {
		return nil
	}
	var componentNames []string
	var componentBuf [16]string
	if dotIdx != -1 {
		rem := qualifiedName[dotIdx+1:]
		dots := strings.Count(rem, ".")
		if dots < len(componentBuf) {
			componentNames = componentBuf[:0]
			start := 0
			for i := 0; i < len(rem); i++ {
				if rem[i] == '.' {
					componentNames = append(componentNames, rem[start:i])
					start = i + 1
				}
			}
			componentNames = append(componentNames, rem[start:])
		} else {
			componentNames = strings.Split(rem, ".")
		}
	}
	return r.ResolveViaPropertyGivenSlot(slot.GetType(), NilNode, componentNames)
}

// ResolveViaPropertyGivenSlot resolves a type using a given StaticTypedSlot
// and list of properties on that type.
//
// slotType: the JSType of the slot, possibly null. definitionNode: if known,
// the Node representing the type definition.
//
// Java: JSTypeRegistry.resolveViaPropertyGivenSlot(JSType, Node, List<String>)
func (r *JSTypeRegistry) ResolveViaPropertyGivenSlot(slotType JSType, definitionNode Node, componentNames []string) JSType {
	if len(componentNames) == 0 {
		typedefType := JSTypeRegistry_ResolveTypeFromNodeIfTypedef(definitionNode)
		if typedefType != nil {
			return typedefType
		}
	}

	// If the first component has a type of 'Unknown', then any type names
	// using it should be regarded as silently 'Unknown' rather than be noisy
	// about it.
	if slotType == nil || slotType.IsAllType() || slotType.IsNoType() {
		return nil
	}

	// resolving component by component
	for i, component := range componentNames {
		parentObj := ObjectType_Cast(slotType)
		if parentObj == nil || component == "" || !parentObj.HasOwnProperty(component) {
			return nil
		}
		if i == len(componentNames)-1 {
			// Look for a typedefTypeProp on the definition node of the last
			// component.
			def := parentObj.GetPropertyDefSite(component)
			typedefType := JSTypeRegistry_ResolveTypeFromNodeIfTypedef(def)
			if typedefType != nil {
				return typedefType
			}
		}
		slotType = parentObj.GetPropertyType(component)
	}

	// Translate "constructor" types to "instance" types.
	if slotType == nil {
		return nil
	} else if slotType.IsFunctionType() && (slotType.IsConstructor() || slotType.IsInterface()) {
		return jsTypeOrNil(slotType.ToMaybeFunctionType().GetInstanceType())
	} else if slotType.IsNoObjectType() {
		return r.GetNativeObjectType(JSTypeNative_NO_OBJECT_TYPE)
	} else if enumType, ok := slotType.(*EnumType); ok {
		return enumType.GetElementsType()
	}
	return nil
}

// JSTypeRegistry_ResolveTypeFromNodeIfTypedef checks the given Node for a
// typedef annotation, resolving to that type if existent.
//
// Java: JSTypeRegistry.resolveTypeFromNodeIfTypedef(Node)
func JSTypeRegistry_ResolveTypeFromNodeIfTypedef(node Node) JSType {
	if node.IsNil() {
		return nil
	}
	return node.GetTypedefTypeProp()
}

// Java: JSTypeRegistry.getTypeForScopeInternal(StaticScope, String)
func (r *JSTypeRegistry) GetTypeForScopeInternal(scope StaticScope, name string) JSType {
	return r.scopedNameTable[scopedTypeKey{root: r.GetRootNodeForScope(scope), name: name}]
}

// Java: JSTypeRegistry.registerGlobalType(JSType)
func (r *JSTypeRegistry) RegisterGlobalType(type_ JSType) {
	r.Register(nil, type_, type_.String())
}

// Java: JSTypeRegistry.registerGlobalType(JSType, String)
func (r *JSTypeRegistry) RegisterGlobalTypeWithName(type_ JSType, name string) {
	r.Register(nil, type_, name)
}

// Java: JSTypeRegistry.register(StaticScope, JSType, String)
func (r *JSTypeRegistry) Register(scope StaticScope, type_ JSType, name string) {
	JSTypeRegistry_CheckTypeName(name)
	r.RegisterForScope(JSTypeRegistry_GetLookupScope(scope, name), type_, name)
}

// Java: JSTypeRegistry.registerForScope(StaticScope, JSType, String)
func (r *JSTypeRegistry) RegisterForScope(scope StaticScope, type_ JSType, name string) {
	r.scopedNameTablePut(r.GetRootNodeForScope(scope), name, type_)
}

// scopedNameTablePut is Table.put(row, column, value), which rejects nulls.
func (r *JSTypeRegistry) scopedNameTablePut(rowKey Node, columnKey string, value JSType) {
	CheckNotNull(rowKey)
	jutil.CheckNotNull(value)
	r.scopedNameTable[scopedTypeKey{root: rowKey, name: columnKey}] = value
}

// DeclareType records declared global type names. This makes resolution
// faster and more robust in the common case.
//
// name: the name of the type to be recorded. type: the actual type being
// associated with the name. Returns true if this name is not already
// defined, false otherwise.
//
// Java: JSTypeRegistry.declareType(StaticScope, String, JSType)
func (r *JSTypeRegistry) DeclareType(scope StaticScope, name string, type_ JSType) bool {
	jutil.CheckState(name != "")
	if r.GetTypeForScopeInternal(JSTypeRegistry_GetLookupScope(scope, name), name) != nil {
		return false
	}

	r.Register(scope, type_, name)
	return true
}

// DeclareTypeForExactScope records declared global type names. This makes
// resolution faster and more robust in the common case.
//
// Java: JSTypeRegistry.declareTypeForExactScope(StaticScope, String, JSType)
func (r *JSTypeRegistry) DeclareTypeForExactScope(scope StaticScope, name string, type_ JSType) bool {
	jutil.CheckState(name != "")
	if r.GetTypeForScopeInternal(scope, name) != nil {
		return false
	}

	r.RegisterForScope(scope, type_, name)
	return true
}

// OverwriteDeclaredType overrides a declared global type name. Throws an
// exception if this type name hasn't been declared yet.
//
// Java: JSTypeRegistry.overwriteDeclaredType(StaticScope, String, JSType)
func (r *JSTypeRegistry) OverwriteDeclaredType(scope StaticScope, name string, type_ JSType) {
	jutil.CheckStateMsg1(r.IsDeclaredForScope(scope, name), "missing name %s", name)
	r.Register(scope, type_, name)
}

// IsForwardDeclaredType reports whether this is a forward-declared type
// name.
//
// Java: JSTypeRegistry.isForwardDeclaredType(String)
func (r *JSTypeRegistry) IsForwardDeclaredType(name string) bool {
	return r.forwardDeclaredTypes != nil && r.forwardDeclaredTypes.Contains(name)
}

// GetReadableTypeName first dereferences the JSType to remove
// null/undefined then returns a human-readable type name
//
// Java: JSTypeRegistry.getReadableTypeName(Node)
func (r *JSTypeRegistry) GetReadableTypeName(n Node) string {
	return r.GetReadableJSTypeName(n, true)
}

// Java: JSTypeRegistry.getReadableTypeNameNoDeref(Node)
func (r *JSTypeRegistry) GetReadableTypeNameNoDeref(n Node) string {
	return r.GetReadableJSTypeName(n, false)
}

// Java: JSTypeRegistry.getSimpleReadableJSTypeName(JSType)
func JSTypeRegistry_GetSimpleReadableJSTypeName(type_ JSType) *string {
	if ((uint32(1) << (type_.GetJSTypeBase().objFlags & objKindMask)) & jsKindNonObjectMask) != 0 {
		s := type_.String()
		return &s
	} else if type_.IsFunctionPrototypeType() {
		s := type_.String()
		return &s
	} else if _, ok := type_.(ObjectType); ok {
		if !JSTypeRegistry_IsReadableObjectType(type_.ToMaybeObjectType()) {
			return nil
		}

		source := type_.ToObjectType().GetConstructor().GetSource()
		if !source.IsNil() {
			jutil.CheckStateMsg1(source.IsFunction() || source.IsClass(), "%s", source)
			readable := source.GetFirstChild().GetOriginalName()
			if readable != nil {
				return readable
			}
		}
		s := type_.String()
		return &s

	} else if unionType, ok := type_.(*UnionType); ok {
		alts := unionType.GetAlternates()
		if len(alts) == 0 {
			u := "null)"
			return &u
		}
		var sb strings.Builder
		for i, alternate := range alts {
			name := JSTypeRegistry_GetSimpleReadableJSTypeName(alternate)
			if name == nil {
				return nil
			}
			if i == 0 {
				sb.WriteByte('(')
			} else {
				sb.WriteByte('|')
			}
			sb.WriteString(*name)
		}
		sb.WriteByte(')')
		u := sb.String()
		return &u
	}
	return nil
}

// GetReadableJSTypeName: given a node, get a human-readable name for the
// type of that node so that will be easy for the programmer to find the
// original declaration.
//
// For example, if SubFoo's property "bar" might have the human-readable name
// "Foo.prototype.bar".
//
// n: the node. dereference: if true, the type of the node will be
// dereferenced to an Object type, if possible. Prefer to call
// #getReadableTypeName(String) or #getReadableTypeNameNoDeref(String)
// instead of passing this as an argument.
//
// Java: JSTypeRegistry.getReadableJSTypeName(Node, boolean)
func (r *JSTypeRegistry) GetReadableJSTypeName(n Node, dereference bool) string {
	type_ := r.GetJSTypeOrUnknown(n)
	if dereference {
		autoboxed := type_.Autobox()
		if !autoboxed.IsNoType() {
			type_ = autoboxed
		}
	}

	name := JSTypeRegistry_GetSimpleReadableJSTypeName(type_)
	if name != nil {
		return *name
	}

	// If we're analyzing a GETPROP, the property may be inherited by the
	// prototype chain. So climb the prototype chain and find out where the
	// property was originally defined.
	if n.IsGetProp() {
		objectType := r.GetJSTypeOrUnknown(n.GetFirstChild()).Dereference()
		if objectType != nil {
			propName := n.GetString()
			objectType = objectType.GetClosestDefiningType(propName)

			// Don't show complex function names or anonymous types. Instead, try to
			// get a human-readable type name.
			if JSTypeRegistry_IsReadableObjectType(objectType) {
				return objectType.String() + "." + propName
			}
		}
	}

	if n.IsQualifiedName() {
		return *n.GetQualifiedName()
	} else if type_.IsFunctionType() {
		// Don't show complex function names.
		return "function"
	}
	return type_.String()
}

// Java: JSTypeRegistry.isReadableObjectType(ObjectType)
func JSTypeRegistry_IsReadableObjectType(type_ ObjectType) bool {
	if type_ == nil {
		return false
	} else if type_.IsInstanceType() {
		return true
	} else if type_.IsFunctionPrototypeType() {
		return true
	} else if type_.IsEnumElementType() {
		// TODO(b/147236174): This case should be deleted as impossible.
		return type_.ToMaybeEnumElementType().GetConstructor() != nil
	}
	return false
}

// Java: JSTypeRegistry.getJSTypeOrUnknown(Node)
func (r *JSTypeRegistry) GetJSTypeOrUnknown(n Node) JSType {
	jsType := n.GetJSType()
	if jsType == nil {
		return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	}
	return jsType
}

// Java: JSTypeRegistry.getGlobalType(String)
func (r *JSTypeRegistry) GetGlobalType(jsTypeName string) JSType {
	return r.GetType(nil, jsTypeName)
}

// GetType looks up a native type by name. Returns the corresponding JSType
// object or null it cannot be found.
//
// Java: JSTypeRegistry.getType(StaticScope, String)
func (r *JSTypeRegistry) GetType(scope StaticScope, jsTypeName string) JSType {
	return r.GetTypeInternal(scope, jsTypeName)
}

func (r *JSTypeRegistry) getTypeResolvedOrNil(scope StaticTypedScope, jsTypeName string) JSType {
	switch jsTypeName {
	case "boolean":
		return r.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
	case "number":
		return r.GetNativeType(JSTypeNative_NUMBER_TYPE)
	case "bigint":
		return r.GetNativeType(JSTypeNative_BIGINT_TYPE)
	case "string":
		return r.GetNativeType(JSTypeNative_STRING_TYPE)
	case "undefined", "void":
		return r.GetNativeType(JSTypeNative_VOID_TYPE)
	}
	// Resolve template type names
	var thisType JSType
	if scope != nil && scope.GetTypeOfThis() != nil {
		thisType = jsTypeOrNil(scope.GetTypeOfThis().ToObjectType())
	}
	if thisType != nil {
		// get the 'last' template type key with this name because if a
		// superclass template key shares the same name as a subclass key, the
		// subclass key is guaranteed to be last.
		key := thisType.GetTemplateTypeMap().GetLastTemplateTypeKeyByName(jsTypeName)
		if key != nil {
			type_ := JSType(key)
			jutil.CheckStateMsg1(type_.IsTemplateType(), "expected:%s", type_)
			return type_
		}
	}

	// TODO(sdh): The use of "getType" here is incorrect. This currently will
	// pick up a type in an outer scope if it will be shadowed by a local
	// type. But creating a unique NamedType object for every name referenced
	// (even if interned) in every scope would be expensive.
	//
	// Instead perhaps a standard (untyped) scope object might be used to pick
	// the right scope during type construction.
	var untypedScope StaticScope
	if scope != nil {
		untypedScope = scope
	}
	return r.GetType(untypedScope, jsTypeName)
}

func (r *JSTypeRegistry) getTypeWithLocationString(
	scope StaticTypedScope, jsTypeName string, sourceName string, lineno int, charno int) JSType {
	if type_ := r.getTypeResolvedOrNil(scope, jsTypeName); type_ != nil {
		return type_
	}
	return r.createNamedTypeByID(scope, jsTypeName, rhinoStringPoolIntern(sourceName), lineno, charno)
}

// GetTypeWithLocation looks up a type by name. To allow for forward
// references to types, an unrecognized string has to be bound to a
// NamedType object that will be resolved later.
//
// Java: JSTypeRegistry.getType(StaticTypedScope, String, String, int, int)
func (r *JSTypeRegistry) GetTypeWithLocation(
	scope StaticTypedScope, jsTypeName string, sourceName *string, lineno int, charno int) JSType {
	if type_ := r.getTypeResolvedOrNil(scope, jsTypeName); type_ != nil {
		return type_
	}
	// TODO(user): Each instance should support named type creation using
	// interning.
	return r.CreateNamedType(scope, jsTypeName, sourceName, lineno, charno)
}

// Java: JSTypeRegistry.evaluateTypeExpressionInGlobalScope(JSTypeExpression)
func (r *JSTypeRegistry) EvaluateTypeExpressionInGlobalScope(expr *JSTypeExpression) JSType {
	return expr.Evaluate(nil, r)
}

// RegisterNonLegacyClosureNamespace registers a goog.module namespace (that
// does not have goog.module.declareLegacyNamespace)
//
// This allows JSTypeRegistry to resolve types that refer to goog.modules by
// namespace. These have unique handling because they exist only in the type
// space and do not have a corresponding value space value.
//
// Java: JSTypeRegistry.registerNonLegacyClosureNamespace(String, Node, JSType)
func (r *JSTypeRegistry) RegisterNonLegacyClosureNamespace(moduleName string, definitionNode Node, type_ JSType) {
	if r.closureNamespaces == nil {
		r.closureNamespaces = make(map[string]*JSTypeRegistry_ClosureNamespace)
	}
	r.closureNamespaces[moduleName] = JSTypeRegistry_ClosureNamespace_Create( /* isLegacy= */ false, definitionNode, type_)
}

// RegisterLegacyClosureNamespace registers a goog.provide or legacy
// goog.module namespace with the type registry
//
// Java: JSTypeRegistry.registerLegacyClosureNamespace(String)
func (r *JSTypeRegistry) RegisterLegacyClosureNamespace(moduleName string) {
	if r.closureNamespaces == nil {
		r.closureNamespaces = make(map[string]*JSTypeRegistry_ClosureNamespace)
	}
	r.closureNamespaces[moduleName] = JSTypeRegistry_ClosureNamespace_Create( /* isLegacy= */ true, NilNode, nil)
}

// JSTypeRegistry_ClosureNamespace stores information about a Closure
// namespace.
//
// Java: abstract class JSTypeRegistry.ClosureNamespace, with its generated
// AutoValue_JSTypeRegistry_ClosureNamespace implementation folded in.
type JSTypeRegistry_ClosureNamespace struct {
	isLegacy       bool
	definitionNode Node
	type_          JSType
}

// IsLegacy returns true if this is a goog.provide'd namespace or a
// goog.module namespace followed by `goog.module.declareLegacyNamespace()
//
// Java: JSTypeRegistry.ClosureNamespace.isLegacy()
func (c *JSTypeRegistry_ClosureNamespace) IsLegacy() bool { return c.isLegacy }

// Java: JSTypeRegistry.ClosureNamespace.definitionNode()
func (c *JSTypeRegistry_ClosureNamespace) DefinitionNode() Node { return c.definitionNode }

// Java: JSTypeRegistry.ClosureNamespace.type()
func (c *JSTypeRegistry_ClosureNamespace) Type() JSType { return c.type_ }

// Java: JSTypeRegistry.ClosureNamespace.create(boolean, Node, JSType)
func JSTypeRegistry_ClosureNamespace_Create(isLegacy bool, definitionNode Node, type_ JSType) *JSTypeRegistry_ClosureNamespace {
	return &JSTypeRegistry_ClosureNamespace{isLegacy: isLegacy, definitionNode: definitionNode, type_: type_}
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.toString()
func (c *JSTypeRegistry_ClosureNamespace) String() string {
	return "ClosureNamespace{" +
		"isLegacy=" + jutil.ToJavaString(c.isLegacy) + ", " +
		"definitionNode=" + jutil.ToJavaString(c.definitionNode) + ", " +
		"type=" + jutil.ToJavaString(c.type_) +
		"}"
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.equals(Object)
func (c *JSTypeRegistry_ClosureNamespace) Equals(o any) bool {
	if any(c) == o {
		return true
	}
	if that, ok := o.(*JSTypeRegistry_ClosureNamespace); ok && that != nil {
		return c.isLegacy == that.isLegacy &&
			c.definitionNode == that.definitionNode &&
			((c.type_ == nil && that.type_ == nil) || (c.type_ != nil && c.type_.Equals(that.type_)))
	}
	return false
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.hashCode()
func (c *JSTypeRegistry_ClosureNamespace) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.isLegacy)
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.definitionNode)
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.type_)
	return h
}
