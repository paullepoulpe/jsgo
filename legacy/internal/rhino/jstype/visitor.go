// Port of com/google/javascript/rhino/jstype/Visitor.java
// @ closure-compiler v20260415
//
// Java uses a generic type parameter T on the Visitor interface.
// Go interfaces cannot carry generic methods, so each case method
// returns `any`; callers type-assert the concrete result.
// WithDefaultCase is ported as a concrete struct that embeds
// DefaultCaseFunc; implementations embed *WithDefaultCase and
// override only the cases they need.
//
// NOTE: The JSType.visit / JSType.accept wiring is intentionally
// absent here (additive-only phase). Wire it when the visitor is
// actively consumed.

package jstype

// Visitor is a visitor for JSType values.
//
// During visitation an instance may calculate a result value.
// Port of Visitor<T> — T is erased to any.
type Visitor interface {
	// CaseNoType is the bottom type's case.
	CaseNoType(typ JSType) any

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(typ *EnumElementType) any

	// CaseAllType is the all type's case.
	CaseAllType() any

	// CaseBooleanType is the boolean value type's case.
	CaseBooleanType() any

	// CaseNoObjectType is the bottom object type's case.
	CaseNoObjectType() any

	// CaseFunctionType is the function type's case.
	CaseFunctionType(typ *FunctionType) any

	// CaseObjectType is the object type's case.
	CaseObjectType(typ *ObjectType) any

	// CaseUnknownType is the unknown type's case.
	CaseUnknownType() any

	// CaseNullType is the null type's case.
	CaseNullType() any

	// CaseNamedType is the named type's case.
	// TODO: replace JSType with *NamedType once NamedType is ported.
	CaseNamedType(typ JSType) any

	// CaseProxyObjectType is the proxy type's case.
	// TODO: replace JSType with *ProxyObjectType once ProxyObjectType is ported.
	CaseProxyObjectType(typ JSType) any

	// CaseNumberType is the number value type's case.
	CaseNumberType() any

	// CaseBigIntType is the BigInt value type's case.
	CaseBigIntType() any

	// CaseStringType is the string value type's case.
	CaseStringType() any

	// CaseSymbolType is the symbol value type's case.
	// TODO: replace JSType with *SymbolType once SymbolType is ported.
	CaseSymbolType(typ JSType) any

	// CaseVoidType is the void type's case.
	CaseVoidType() any

	// CaseUnionType is the union type's case.
	CaseUnionType(typ *UnionType) any

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(typ *TemplatizedType) any

	// CaseTemplateType is the template type's case.
	CaseTemplateType(typ *TemplateType) any
}

// DefaultCaseFunc is the fallback called by WithDefaultCase for any
// case not overridden by the concrete subtype. Receives nil when the
// originating case carries no JSType argument (e.g. CaseAllType).
// Port of Visitor.WithDefaultCase.caseDefault.
type DefaultCaseFunc func(typ JSType) any

// WithDefaultCase is an embeddable base for Visitor implementations
// that want a uniform default behaviour. Embed it, supply a
// DefaultCase func, and override only the cases you need.
// Port of Visitor.WithDefaultCase<T>.
type WithDefaultCase struct {
	DefaultCase DefaultCaseFunc
}

func (w *WithDefaultCase) CaseNoType(typ JSType) any          { return w.DefaultCase(typ) }
func (w *WithDefaultCase) CaseEnumElementType(typ *EnumElementType) any {
	return w.DefaultCase(typ)
}
func (w *WithDefaultCase) CaseAllType() any                   { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseBooleanType() any               { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseNoObjectType() any              { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseFunctionType(typ *FunctionType) any {
	return w.DefaultCase(typ)
}
func (w *WithDefaultCase) CaseObjectType(typ *ObjectType) any { return w.DefaultCase(typ) }
func (w *WithDefaultCase) CaseUnknownType() any               { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseNullType() any                  { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseNamedType(typ JSType) any       { return w.DefaultCase(typ) }
func (w *WithDefaultCase) CaseProxyObjectType(typ JSType) any { return w.DefaultCase(typ) }
func (w *WithDefaultCase) CaseNumberType() any                { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseBigIntType() any                { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseStringType() any                { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseSymbolType(typ JSType) any      { return w.DefaultCase(typ) }
func (w *WithDefaultCase) CaseVoidType() any                  { return w.DefaultCase(nil) }
func (w *WithDefaultCase) CaseUnionType(typ *UnionType) any   { return w.DefaultCase(typ) }
func (w *WithDefaultCase) CaseTemplatizedType(typ *TemplatizedType) any {
	return w.DefaultCase(typ)
}
func (w *WithDefaultCase) CaseTemplateType(typ *TemplateType) any { return w.DefaultCase(typ) }
