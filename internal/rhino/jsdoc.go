// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20260415
//
// The Java class is large; this port carries the annotations Tier-0 and
// Tier-1 need: @type, @param, @return. Extended in a later pass to include
// @private/@protected/@public (visibility), @nosideeffects, @abstract,
// @dict, and @struct.

package rhino

// Visibility mirrors JSDocInfo.Visibility in Java. The ordinal order matches
// Java's enum order (PRIVATE=0, PACKAGE=1, PROTECTED=2, PUBLIC=3, INHERITED=4)
// so that comparisons via < / > retain the Java semantics ("0 is most private").
type Visibility int

const (
	VisibilityPrivate   Visibility = iota // PRIVATE
	VisibilityPackage                     // PACKAGE (not yet parsed; reserved for completeness)
	VisibilityProtected                   // PROTECTED
	VisibilityPublic                      // PUBLIC
	VisibilityInherited                   // INHERITED — default when no tag is present
)

// BoundedTemplatePosition records the source position of a `@template {B} T`
// annotation. Used by TypeCheck to emit JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES.
// Port of RhinoErrorReporter.UNSUPPORTED_BOUNDED_GENERIC_TYPES + the
// addParserWarning(JSDOC_TEMPLATE_BOUNDEDGENERICS_USED, lineno, charno) call
// in JsDocInfoParser.java's TEMPLATE case.
type BoundedTemplatePosition struct {
	Line int // 1-indexed absolute file line of the @template annotation
	Col  int // 0-indexed column of the @ character
}

type JSDocInfo struct {
	typeExpr        *JSTypeExpression
	paramNames      []string                     // insertion order (matches Java)
	params          map[string]*JSTypeExpression // name → type
	restParams      map[string]bool              // param names with `{...T}` rest marker
	optionalParams  map[string]bool              // param names with `{T=}` optional marker
	returnExpr      *JSTypeExpression
	isConstructor   bool
	isInterface     bool
	isRecord        bool
	isStatic        bool
	isStruct        bool
	isDict          bool
	isAbstract      bool
	isNoSideEffects bool
	isConst         bool // @const annotation
	visibility      Visibility // default VisibilityInherited
	hasVisibility   bool       // true once any @private/@protected/@public was recorded
	enumExpr        *JSTypeExpression // non-nil when @enum was present
	enumMarked      bool              // true when @enum appeared (even with no type)
	extendsExpr     *JSTypeExpression
	extendsAll      []*JSTypeExpression // every @extends encountered (interfaces can list many)
	implementsAll   []*JSTypeExpression
	typedefExpr     *JSTypeExpression
	templateNames            []string
	ttlTemplateNames         []string // @template T := … =: transformation names
	boundedTemplatePositions []BoundedTemplatePosition
	suppressions    map[string]bool
	hasThisType     bool   // true when @this annotation was present
	isOverride      bool   // true when @override annotation was present
	lendsName       string // non-empty when @lends annotation was present
	isDeprecated    bool
	deprecationReason string
	isFinal         bool
	// isCastHoisted marks a JSDocInfo that was hoisted from a
	// `/** @type {T} */ (expr)` cast expression onto the enclosing
	// declaration node. Go's parser lacks a CAST AST node (Java has one);
	// this flag lets visitVar skip the initializer subtype check (Java
	// uses the CAST node's type as the effective RHS type instead).
	isCastHoisted   bool
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

// Suppressions returns the full set of suppression names. Mirrors
// JSDocInfo.getSuppressions() in Java, used by CheckJSDoc.validateSuppress.
func (i *JSDocInfo) Suppressions() map[string]bool {
	if i == nil {
		return nil
	}
	return i.suppressions
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

// TTLTemplateNames returns @template T := … =: transformation template names.
// Mirrors JSDocInfo.getTypeTransformations() key-set in Java.
func (i *JSDocInfo) TTLTemplateNames() []string {
	if i == nil {
		return nil
	}
	return i.ttlTemplateNames
}
func (i *JSDocInfo) AddTTLTemplateName(n string) {
	i.ttlTemplateNames = append(i.ttlTemplateNames, n)
}

// BoundedTemplatePositions returns all @template {B} T positions recorded
// during JSDoc parsing. Mirrors the addParserWarning(JSDOC_TEMPLATE_BOUNDEDGENERICS_USED)
// call in JsDocInfoParser.java — each entry becomes a JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES
// diagnostic in TypeCheck.
func (i *JSDocInfo) BoundedTemplatePositions() []BoundedTemplatePosition {
	if i == nil {
		return nil
	}
	return i.boundedTemplatePositions
}
func (i *JSDocInfo) AddBoundedTemplatePosition(p BoundedTemplatePosition) {
	i.boundedTemplatePositions = append(i.boundedTemplatePositions, p)
}

// IsInterface / IsRecord mirror JSDocInfo.isInterface() /
// JSDocInfo.usesImplicitMatch() — together they mark a class as a
// structural type.
func (i *JSDocInfo) IsInterface() bool   { return i != nil && i.isInterface }
func (i *JSDocInfo) IsRecord() bool      { return i != nil && i.isRecord }
func (i *JSDocInfo) SetInterface(b bool) { i.isInterface = b }
func (i *JSDocInfo) SetRecord(b bool)    { i.isRecord = b }

// MakesStructs / MakesDicts mirror JSDocInfo.makesStructs() /
// makesDicts() — annotations on a constructor that restrict
// property-access shape. `@struct` forbids `x['prop']` computed
// access; `@dict` forbids `x.prop` dotted access. Propagated onto
// the resulting InstanceType by TypedScopeCreator so TypeCheck can
// emit JSC_ILLEGAL_PROPERTY_ACCESS at receiver sites.
// IsStruct / IsDict are retained as local aliases for convenience.
func (i *JSDocInfo) MakesStructs() bool { return i != nil && i.isStruct }
func (i *JSDocInfo) MakesDicts() bool   { return i != nil && i.isDict }
func (i *JSDocInfo) IsStruct() bool     { return i.MakesStructs() }
func (i *JSDocInfo) IsDict() bool       { return i.MakesDicts() }
func (i *JSDocInfo) SetStruct(b bool)   { i.isStruct = b }
func (i *JSDocInfo) SetDict(b bool)     { i.isDict = b }

// IsAbstract mirrors JSDocInfo.isAbstract() — true when the comment
// carried an `@abstract` annotation.
func (i *JSDocInfo) IsAbstract() bool   { return i != nil && i.isAbstract }
func (i *JSDocInfo) SetAbstract(b bool) { i.isAbstract = b }

// IsNoSideEffects mirrors JSDocInfo.isNoSideEffects() — true when
// the comment carried a `@nosideeffects` annotation.
func (i *JSDocInfo) IsNoSideEffects() bool   { return i != nil && i.isNoSideEffects }
func (i *JSDocInfo) SetNoSideEffects(b bool) { i.isNoSideEffects = b }

// GetVisibility mirrors JSDocInfo.getVisibility() — returns the
// visibility specified by @private, @protected, or @public. Defaults
// to VisibilityInherited when no visibility tag was present.
func (i *JSDocInfo) GetVisibility() Visibility {
	if i == nil || !i.hasVisibility {
		return VisibilityInherited
	}
	return i.visibility
}

// RecordVisibility mirrors JSDocInfo.Builder.recordVisibility() —
// sets the visibility and returns true. Returns false and does not
// overwrite if visibility was already recorded (duplicates ignored,
// matching Java's idempotency guard).
func (i *JSDocInfo) RecordVisibility(v Visibility) bool {
	if i.hasVisibility {
		return false
	}
	i.visibility = v
	i.hasVisibility = true
	return true
}

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
// Returns the FIRST extends tag when the comment has multiple (an
// interface extending multiple supers); callers that need every
// parent use BaseTypes.
func (i *JSDocInfo) BaseType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.extendsExpr
}

// BaseTypes returns every @extends expression on the comment, in
// declaration order. Interfaces permit multi-inheritance so this
// may carry more than one entry; regular constructors carry at
// most one.
func (i *JSDocInfo) BaseTypes() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.extendsAll
}

func (i *JSDocInfo) SetBaseType(t *JSTypeExpression) {
	if i.extendsExpr == nil {
		i.extendsExpr = t
	}
	i.extendsAll = append(i.extendsAll, t)
}

// IsConstructor mirrors JSDocInfo.isConstructor().
func (i *JSDocInfo) IsConstructor() bool { return i != nil && i.isConstructor }

func (i *JSDocInfo) SetConstructor(b bool) { i.isConstructor = b }

// IsConstant mirrors JSDocInfo.isConstant() — true when the @const
// annotation was present. (Java also checks @define / @final / @desc;
// jsgo adds those as needed.)
func (i *JSDocInfo) IsConstant() bool { return i != nil && i.isConst }

// IsConst is an alias used by CheckAccessControls to match Java's
// JSDocInfo.isConstant() accessor name.
func (i *JSDocInfo) IsConst() bool { return i.IsConstant() }

// SetConst records that the @const annotation was present.
func (i *JSDocInfo) SetConst(b bool) { i.isConst = b }

// IsDeprecated / GetDeprecationReason / SetDeprecated mirror
// JSDocInfo.isDeprecated() / getDeprecationReason() and the builder's
// recordDeprecated(reason). Used by CheckAccessControls.
func (i *JSDocInfo) IsDeprecated() bool {
	return i != nil && (i.isDeprecated || i.deprecationReason != "")
}
func (i *JSDocInfo) GetDeprecationReason() string {
	if i == nil {
		return ""
	}
	return i.deprecationReason
}
func (i *JSDocInfo) SetDeprecated(reason string) {
	i.isDeprecated = true
	i.deprecationReason = reason
}

// IsFinal / SetFinal mirror JSDocInfo.isFinal().
func (i *JSDocInfo) IsFinal() bool    { return i != nil && i.isFinal }
func (i *JSDocInfo) SetFinal(b bool)  { i.isFinal = b }

// IsCastHoisted / SetCastHoisted track whether this JSDocInfo was hoisted
// from a `/** @type {T} */ (expr)` type-cast onto its enclosing declaration
// node. Go's parser lacks a CAST AST node (Java has one); this flag lets
// visitVar skip the initializer subtype check when the declared type came
// from a cast (Java uses the CAST node's type as the effective RHS type).
func (i *JSDocInfo) IsCastHoisted() bool   { return i != nil && i.isCastHoisted }
func (i *JSDocInfo) SetCastHoisted(b bool) { i.isCastHoisted = b }

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

// HasThisType mirrors JSDocInfo.hasThisType() — true when @this was present.
func (i *JSDocInfo) HasThisType() bool { return i != nil && i.hasThisType }

// SetHasThisType records that @this appeared in the comment.
func (i *JSDocInfo) SetHasThisType(b bool) { i.hasThisType = b }

// IsOverride mirrors JSDocInfo.isOverride() — true when @override was present.
func (i *JSDocInfo) IsOverride() bool { return i != nil && i.isOverride }

// SetOverride records that @override appeared in the comment.
func (i *JSDocInfo) SetOverride(b bool) { i.isOverride = b }

// HasLendsName mirrors JSDocInfo.hasLendsName() — true when @lends was present.
func (i *JSDocInfo) HasLendsName() bool { return i != nil && i.lendsName != "" }

// LendsName mirrors JSDocInfo.getLendsName() — returns the lends target as a
// JSTypeExpression so callers can inspect the name string.
func (i *JSDocInfo) LendsName() *JSTypeExpression {
	if i == nil || i.lendsName == "" {
		return nil
	}
	return &JSTypeExpression{TypeName: i.lendsName}
}

// SetLendsName records the @lends target name.
func (i *JSDocInfo) SetLendsName(name string) { i.lendsName = name }

// ContainsDeclaration mirrors JSDocInfo.containsDeclaration() — returns
// true when the comment carries any declaration-forming annotation.
func (i *JSDocInfo) ContainsDeclaration() bool {
	if i == nil {
		return false
	}
	return i.typeExpr != nil ||
		i.returnExpr != nil ||
		len(i.params) > 0 ||
		i.isConstructor ||
		i.isInterface ||
		i.isRecord ||
		i.enumMarked ||
		i.typedefExpr != nil ||
		i.extendsExpr != nil ||
		len(i.implementsAll) > 0
}
