// Port of com/google/javascript/rhino/jstype/CanCastToVisitor.java
// @ closure-compiler v20260415
//
// CanCastToVisitor implements RelationshipVisitor to determine whether
// a cast from one JSType to another is permissible.
//
// Java calls thisType.visit(visitor, thatType) to dispatch into the
// visitor; in Go the JSType interface does not carry Visit/accept
// methods yet (additive-only phase). The dispatch is centralised in
// the package-level VisitRelationship function below so callers can
// use the visitor without retrofitting existing concrete types.
//
// Several Java helpers that reference not-yet-ported types
// (NoResolvedType, NamedType, ProxyObjectType, ValueType, SymbolType,
// NoObjectType, NoType) are stubbed with conservative TODO comments.

package jstype

// VisitRelationship dispatches thisType through v with thatType as
// the second argument. This is the Go equivalent of the Java pattern
// thisType.visit(RelationshipVisitor, JSType).
//
// Type cases handled:
//   - *UnionType        → CaseUnionType
//   - *TemplatizedType  → CaseTemplatizedType
//   - *TemplateType     → CaseTemplateType
//   - *FunctionType     → CaseFunctionType
//   - *ObjectType       → CaseObjectType
//   - *EnumElementType  → CaseEnumElementType
//   - unknownType       → CaseUnknownType
//   - allType           → CaseAllType
//   - all other primitives (number/string/boolean/null/void) →
//       CaseValueType (mirrors Java's ValueType superclass for
//       primitive scalar types)
//
// NOTE: NoType, NoObjectType, NamedType, ProxyObjectType and
// ValueType (BigInt, Symbol) are not yet ported; once they are, add
// their cases here. Until then those values fall through to
// CaseValueType, which is the conservative choice.
func VisitRelationship(v RelationshipVisitor, thisType JSType, thatType JSType) any {
	switch t := thisType.(type) {
	case unknownType:
		return v.CaseUnknownType(thisType, thatType)
	case allType:
		return v.CaseAllType(thatType)
	case *UnionType:
		return v.CaseUnionType(t, thatType)
	case *TemplatizedType:
		return v.CaseTemplatizedType(t, thatType)
	case *TemplateType:
		return v.CaseTemplateType(t, thatType)
	case *FunctionType:
		return v.CaseFunctionType(t, thatType)
	case *ObjectType:
		return v.CaseObjectType(t, thatType)
	case *EnumElementType:
		return v.CaseEnumElementType(t, thatType)
	default:
		// Covers primitives (numberType, stringType, booleanType,
		// nullType, voidType) — all are subclasses of Java's ValueType.
		// Also covers not-yet-ported types conservatively.
		return v.CaseValueType(thisType, thatType)
	}
}

// CanCastToVisitor implements the "can cast to" relationship.
// Port of CanCastToVisitor which implements RelationshipVisitor<Boolean>.
type CanCastToVisitor struct{}

// CaseUnknownType returns true — unknown types can cast to anything.
func (c *CanCastToVisitor) CaseUnknownType(thisType JSType, thatType JSType) any {
	return true
}

// CaseNoType returns true — the bottom type can cast to anything.
func (c *CanCastToVisitor) CaseNoType(thatType JSType) any {
	return true
}

// CaseNoObjectType returns true.
// TODO(johnlenz): restrict to objects.
func (c *CanCastToVisitor) CaseNoObjectType(thatType JSType) any {
	return true
}

// CaseAllType returns true — the all type can cast to anything.
func (c *CanCastToVisitor) CaseAllType(thatType JSType) any {
	return true
}

// canCastToUnion reports whether thisType can be cast to any
// alternate of unionType.
func (c *CanCastToVisitor) canCastToUnion(thisType JSType, unionType *UnionType) bool {
	for _, typ := range unionType.Members() {
		if VisitRelationship(c, thisType, typ).(bool) {
			return true
		}
	}
	return false
}

// canCastToFunction reports whether thisType can be cast to
// functionType.
func (c *CanCastToVisitor) canCastToFunction(thisType JSType, functionType *FunctionType) bool {
	if _, ok := thisType.(*FunctionType); ok {
		// TODO(johnlenz): visit function parts.
		return true
	}
	return IsSubtype(thisType, functionType) || IsSubtype(functionType, thisType)
}

// isInterface reports whether type is an interface (structural) object.
func (c *CanCastToVisitor) isInterface(typ JSType) bool {
	objType, ok := typ.(*ObjectType)
	if !ok {
		return false
	}
	// In the Go port, IsStructural covers both @record and @interface.
	// Java checks FunctionType.isInterface() on the constructor; we
	// conservatively treat any structural ObjectType as an interface.
	// TODO: narrow once FunctionType carries an IsInterface flag.
	return objType.IsStructural
}

// castCastToHelper is the common implementation used by the value-type,
// object-type, and function-type cases.
// Port of CanCastToVisitor.castCastToHelper.
func (c *CanCastToVisitor) castCastToHelper(thisType JSType, thatType JSType) bool {
	// TODO: isNoResolvedType not yet ported; conservatively treated as
	// false. Add check once NoResolvedType is ported.

	if isUnknown(thatType) || isAll(thatType) {
		return true
	}

	// TODO(johnlenz): restrict to objects once NoObjectType is ported.
	// For now we have no NoObjectType concrete type, so this check is
	// effectively always false — which is the safe conservative choice.

	// TODO: isNoResolvedType(thatType) not yet ported; treated as false.

	if isNil(thatType) {
		// Corresponds to thatType.isNoType() — not yet ported; conservatively
		// treated as true (bottom type is always castable).
		return true
	}

	if c.isInterface(thisType) || c.isInterface(thatType) {
		// TODO(johnlenz): are there any misuses we can catch?
		return true
	}

	// isRecordType → ObjectType.IsStructural in Go.
	if isRecordType(thisType) || isRecordType(thatType) {
		// TODO(johnlenz): are there any misuses we can catch?
		return true
	}

	if ee, ok := thatType.(*EnumElementType); ok {
		return VisitRelationship(c, thisType, ee.Enum.Underlying).(bool)
	}

	if ut, ok := thatType.(*UnionType); ok {
		return c.canCastToUnion(thisType, ut)
	}

	if ft, ok := thatType.(*FunctionType); ok {
		return c.canCastToFunction(thisType, ft)
	}

	if tt, ok := thatType.(*TemplatizedType); ok {
		// TODO(johnlenz): once the templated type work is finished,
		// restrict the type parameters.
		return VisitRelationship(c, thisType, tt.Base).(bool)
	}

	return IsSubtype(thisType, thatType) || IsSubtype(thatType, thisType)
}

// CaseValueType handles primitive (value) types.
func (c *CanCastToVisitor) CaseValueType(thisType JSType, thatType JSType) any {
	return c.castCastToHelper(thisType, thatType)
}

// CaseObjectType handles object types.
func (c *CanCastToVisitor) CaseObjectType(thisType *ObjectType, thatType JSType) any {
	return c.castCastToHelper(thisType, thatType)
}

// CaseFunctionType handles function types.
func (c *CanCastToVisitor) CaseFunctionType(thisType *FunctionType, thatType JSType) any {
	return c.castCastToHelper(thisType, thatType)
}

// CaseUnionType handles union types. Iterates the alternates, skipping
// null and void members (to prevent null|undefined from casting to any
// nullable type). Handles the special null|undefined union.
func (c *CanCastToVisitor) CaseUnionType(thisType *UnionType, thatType JSType) any {
	visited := false
	for _, typ := range thisType.Members() {
		if isNullType(typ) || isVoidType(typ) {
			// Don't allow if the only match between the types is null or
			// void — otherwise any nullable type would be castable to any
			// other nullable type.
			continue
		}
		visited = true
		if VisitRelationship(c, typ, thatType).(bool) {
			return true
		}
	}

	// Special case: the "null|undefined" union may cast to any type that
	// itself allows either null or undefined.
	if !visited {
		return VisitRelationship(c, Null, thatType).(bool) ||
			VisitRelationship(c, Void, thatType).(bool)
	}

	return false
}

// CaseTemplatizedType handles templatized types.
// TODO(johnlenz): once the templated type work is finished,
// restrict the type parameters.
func (c *CanCastToVisitor) CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) any {
	return VisitRelationship(c, thisType.Base, thatType)
}

// CaseTemplateType handles template (unresolved generic) types.
func (c *CanCastToVisitor) CaseTemplateType(thisType *TemplateType, thatType JSType) any {
	return true
}

// CaseEnumElementType handles enum element types by delegating to the
// primitive underlying type.
func (c *CanCastToVisitor) CaseEnumElementType(thisType *EnumElementType, thatType JSType) any {
	return VisitRelationship(c, thisType.Enum.Underlying, thatType)
}

// isRecordType reports whether t is a structural (@record) object type.
// Mirrors JSType.isRecordType() in Java.
func isRecordType(t JSType) bool {
	if obj, ok := t.(*ObjectType); ok {
		return obj.IsStructural
	}
	return false
}

// isNullType reports whether t is the null primitive.
func isNullType(t JSType) bool {
	_, ok := t.(nullType)
	return ok
}

// isVoidType reports whether t is the void (undefined) primitive.
func isVoidType(t JSType) bool {
	_, ok := t.(voidType)
	return ok
}

// isNil reports whether t is the nil JSType (should not normally
// occur; guards against unintentional nil dereferences).
func isNil(t JSType) bool {
	return t == nil
}
