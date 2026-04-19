// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20260415
//
// The Java class is large; this port carries the annotations Tier-0 and
// Tier-1 need: @type, @param, @return. Further tags are added as later
// tiers require them.

package rhino

type JSDocInfo struct {
	typeExpr      *JSTypeExpression
	paramNames    []string                     // insertion order (matches Java)
	params        map[string]*JSTypeExpression // name → type
	restParams    map[string]bool              // param names with `{...T}` rest marker
	optionalParams map[string]bool             // param names with `{T=}` optional marker
	returnExpr    *JSTypeExpression
	isConstructor bool
	isInterface   bool
	isRecord      bool
	isStatic      bool
	isStruct      bool
	isDict        bool
	enumExpr      *JSTypeExpression // non-nil when @enum was present
	enumMarked    bool              // true when @enum appeared (even with no type)
	extendsExpr   *JSTypeExpression
	implementsAll []*JSTypeExpression
	typedefExpr   *JSTypeExpression
	templateNames []string
	suppressions  map[string]bool
}

// AddSuppression records a `@suppress {name}` entry. Mirrors
// JSDocInfo.addSuppression in Java.
func (i *JSDocInfo) AddSuppression(name string) {
	if i.suppressions == nil {
		i.suppressions = map[string]bool{}
	}
	i.suppressions[name] = true
}

// HasSuppression reports whether the named suppression was attached
// to this comment. Used by the per-pass diagnostic gates (e.g.
// VariableReferenceCheck honours `@suppress {duplicate}`).
func (i *JSDocInfo) HasSuppression(name string) bool {
	return i != nil && i.suppressions != nil && i.suppressions[name]
}

// IsStatic marks a class member as belonging to the constructor's
// static slot rather than its instance slot. Mirrors Java's
// MemberFunctionDef.isStaticMember and the analogous field flag.
func (i *JSDocInfo) IsStatic() bool   { return i != nil && i.isStatic }
func (i *JSDocInfo) SetStatic(b bool) { i.isStatic = b }

// TypedefType mirrors JSDocInfo.getTypedefType() — returns the aliased
// type expression captured from `/** @typedef {X} */`.
func (i *JSDocInfo) TypedefType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.typedefExpr
}
func (i *JSDocInfo) SetTypedefType(t *JSTypeExpression) { i.typedefExpr = t }

// TemplateNames mirrors JSDocInfo.getTemplateTypeNames() — the in-order
// list of names captured from `@template T` tags.
func (i *JSDocInfo) TemplateNames() []string {
	if i == nil {
		return nil
	}
	return i.templateNames
}
func (i *JSDocInfo) AddTemplateName(n string) {
	i.templateNames = append(i.templateNames, n)
}

// IsInterface / IsRecord mirror JSDocInfo.isInterface() /
// JSDocInfo.usesImplicitMatch() — together they mark a class as a
// structural type.
func (i *JSDocInfo) IsInterface() bool   { return i != nil && i.isInterface }
func (i *JSDocInfo) IsRecord() bool      { return i != nil && i.isRecord }
func (i *JSDocInfo) SetInterface(b bool) { i.isInterface = b }
func (i *JSDocInfo) SetRecord(b bool)    { i.isRecord = b }

// IsStruct / IsDict mirror JSDocInfo.makesStructs() / makesDicts() —
// annotations on a constructor that restrict property-access shape.
// `@struct` forbids `x['prop']` computed access; `@dict` forbids
// `x.prop` dotted access. Propagated onto the resulting InstanceType
// by TypedScopeCreator so TypeCheck can emit
// JSC_ILLEGAL_PROPERTY_ACCESS at receiver sites.
func (i *JSDocInfo) IsStruct() bool   { return i != nil && i.isStruct }
func (i *JSDocInfo) IsDict() bool     { return i != nil && i.isDict }
func (i *JSDocInfo) SetStruct(b bool) { i.isStruct = b }
func (i *JSDocInfo) SetDict(b bool)   { i.isDict = b }

// IsEnum mirrors JSDocInfo.hasEnumParameterType() — true whenever
// the comment carried an `@enum` tag (with or without an explicit
// type argument).
func (i *JSDocInfo) IsEnum() bool { return i != nil && i.enumMarked }

// EnumParameterType mirrors JSDocInfo.getEnumParameterType() —
// returns the bracketed type argument from `@enum {T}`, or nil when
// `@enum` was used without one (Java defaults to Number in that
// case; jsgo matches by materialising Number in the scope creator).
func (i *JSDocInfo) EnumParameterType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.enumExpr
}
func (i *JSDocInfo) SetEnumParameterType(t *JSTypeExpression) {
	i.enumExpr = t
	i.enumMarked = true
}
func (i *JSDocInfo) MarkEnum() { i.enumMarked = true }

// Implementations mirrors JSDocInfo.getImplementedInterfaces().
func (i *JSDocInfo) Implementations() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.implementsAll
}
func (i *JSDocInfo) AddImplementation(t *JSTypeExpression) {
	i.implementsAll = append(i.implementsAll, t)
}

// BaseType mirrors JSDocInfo.getBaseType() — the type expression
// captured from `@extends {Name}` on a constructor or interface.
func (i *JSDocInfo) BaseType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.extendsExpr
}

func (i *JSDocInfo) SetBaseType(t *JSTypeExpression) { i.extendsExpr = t }

// IsConstructor mirrors JSDocInfo.isConstructor().
func (i *JSDocInfo) IsConstructor() bool { return i != nil && i.isConstructor }

func (i *JSDocInfo) SetConstructor(b bool) { i.isConstructor = b }

func (i *JSDocInfo) Type() *JSTypeExpression     { return i.typeExpr }
func (i *JSDocInfo) SetType(t *JSTypeExpression) { i.typeExpr = t }
func (i *JSDocInfo) HasType() bool               { return i != nil && i.typeExpr != nil }

// ParameterNames matches JSDocInfo.getParameterNames() — iteration
// order is the order @param tags appeared in the comment.
func (i *JSDocInfo) ParameterNames() []string {
	if i == nil {
		return nil
	}
	return i.paramNames
}

// ParameterType matches JSDocInfo.getParameterType(String).
func (i *JSDocInfo) ParameterType(name string) *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.params[name]
}

func (i *JSDocInfo) AddParam(name string, t *JSTypeExpression) {
	if i.params == nil {
		i.params = map[string]*JSTypeExpression{}
	}
	if _, seen := i.params[name]; !seen {
		i.paramNames = append(i.paramNames, name)
	}
	i.params[name] = t
}

// IsRestParam reports whether the named @param carried the rest
// marker (`@param {...T} x`). Mirrors JSDocInfo.hasRestParameter()
// + isRestParameter(name) in Java.
func (i *JSDocInfo) IsRestParam(name string) bool {
	return i != nil && i.restParams[name]
}

// MarkRestParam records that the named @param carries the `{...T}`
// rest marker. Called from JsDocInfoParser when the annotation
// has IterRest at the top of the brace.
func (i *JSDocInfo) MarkRestParam(name string) {
	if i.restParams == nil {
		i.restParams = map[string]bool{}
	}
	i.restParams[name] = true
}

// IsOptionalParam reports whether the named @param carried the
// `{T=}` optional marker — callers may omit this argument without
// a WRONG_ARGUMENT_COUNT diagnostic. Mirrors
// JSDocInfo.isOptionalParameter(name) in Java.
func (i *JSDocInfo) IsOptionalParam(name string) bool {
	return i != nil && i.optionalParams[name]
}

// MarkOptionalParam records that the named @param carries `{T=}`.
func (i *JSDocInfo) MarkOptionalParam(name string) {
	if i.optionalParams == nil {
		i.optionalParams = map[string]bool{}
	}
	i.optionalParams[name] = true
}

// ReturnType matches JSDocInfo.getReturnType().
func (i *JSDocInfo) ReturnType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.returnExpr
}

func (i *JSDocInfo) SetReturnType(t *JSTypeExpression) { i.returnExpr = t }
