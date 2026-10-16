// Port of com/google/javascript/rhino/jstype/RelationshipVisitor.java @ closure-compiler v20260415

package rhino

// RelationshipVisitor is a vistor for relationships between types: the
// visited type is `this` and the second argument is `that`.
//
// Go-only: language-adapter — RelationshipVisitor<T> with T erased to `any`
// (see Visitor).
type RelationshipVisitor interface {
	// CaseUnknownType is the unknown type's case.
	CaseUnknownType(thisType JSType, thatType JSType) any

	// CaseNoType is the bottom type's case.
	CaseNoType(thatType JSType) any

	// CaseNoObjectType is the bottom Object type's case.
	CaseNoObjectType(thatType JSType) any

	// CaseAllType is the all type's case.
	CaseAllType(thatType JSType) any

	// CaseValueType is the value types' case.
	CaseValueType(thisType ValueType, thatType JSType) any

	// CaseObjectType is the object type's case.
	CaseObjectType(thisType ObjectType, thatType JSType) any

	// CaseFunctionType is the function type's case.
	CaseFunctionType(thisType FunctionType, thatType JSType) any

	// CaseUnionType is the union type's case.
	CaseUnionType(thisType *UnionType, thatType JSType) any

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) any

	// CaseTemplateType is the template type's case.
	CaseTemplateType(thisType *TemplateType, thatType JSType) any

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(typeType *EnumElementType, thatType JSType) any
}
