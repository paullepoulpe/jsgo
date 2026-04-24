// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20260415
//
// The Java class is large; this port carries the annotations Tier-0 and
// Tier-1 need: @type, @param, @return. Extended in a later pass to include
// @private/@protected/@public (visibility), @nosideeffects, @abstract,
// @dict, and @struct.

package rhino

import "strings"

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

// String implements fmt.Stringer for Visibility, matching Java's
// Visibility.toString() / name() output (lower-case).
func (v Visibility) String() string {
	switch v {
	case VisibilityPrivate:
		return "private"
	case VisibilityPackage:
		return "package"
	case VisibilityProtected:
		return "protected"
	case VisibilityPublic:
		return "public"
	case VisibilityInherited:
		return "inherited"
	default:
		return "unknown"
	}
}

// BitField mirrors the private Bit enum in JSDocInfo.java. Each constant is a
// 1-shifted bitmask matching Java's ordinal layout (ordinal 0 = bit 0, etc.).
// Only the subset consumed by jsgo's type-checker is named here; the rest are
// reserved as zero values so the bit positions stay stable.
type BitField uint64

const (
	BitInlineType          BitField = 1 << 0  // INLINE_TYPE
	BitIncludeDocumentation BitField = 1 << 1  // INCLUDE_DOCUMENTATION
	BitConst               BitField = 1 << 2  // CONST
	BitConstructor         BitField = 1 << 3  // CONSTRUCTOR
	BitDefine              BitField = 1 << 4  // DEFINE
	BitTypeSummary         BitField = 1 << 5  // TYPE_SUMMARY
	BitFinal               BitField = 1 << 6  // FINAL
	BitOverride            BitField = 1 << 7  // OVERRIDE
	BitDeprecated          BitField = 1 << 8  // DEPRECATED
	BitInterface           BitField = 1 << 9  // INTERFACE
	BitExport              BitField = 1 << 10 // EXPORT
	BitEnhancedNamespace   BitField = 1 << 11 // ENHANCED_NAMESPACE
	BitNoInline            BitField = 1 << 12 // NOINLINE
	BitFileoverview        BitField = 1 << 13 // FILEOVERVIEW
	BitImplicitCast        BitField = 1 << 14 // IMPLICITCAST
	BitNoSideEffects       BitField = 1 << 15 // NOSIDEEFFECTS
	BitExterns             BitField = 1 << 16 // EXTERNS
	BitNoCompile           BitField = 1 << 17 // NOCOMPILE
	BitNoDts               BitField = 1 << 18 // NODTS
	BitUnrestricted        BitField = 1 << 19 // UNRESTRICTED
	BitUsedViaDotCtor      BitField = 1 << 20 // USED_VIA_DOT_CONSTRUCTOR
	BitStruct              BitField = 1 << 21 // STRUCT
	BitDict                BitField = 1 << 22 // DICT
	BitNoCollapse          BitField = 1 << 23 // NOCOLLAPSE
	BitRecord              BitField = 1 << 24 // RECORD
	BitAbstract            BitField = 1 << 25 // ABSTRACT
	BitPureOrBreakMyCode   BitField = 1 << 26 // PURE_OR_BREAK_MY_CODE
	BitCollapsibleOrBreak  BitField = 1 << 27 // COLLAPSIBLE_OR_BREAK_MY_CODE
	BitNoCoverage          BitField = 1 << 28 // NOCOVERAGE
	BitRequireInlining     BitField = 1 << 29 // REQUIRE_INLINING
	BitEncourageInlining   BitField = 1 << 30 // ENCOURAGE_INLINING
	BitNgInject            BitField = 1 << 31 // NG_INJECT
	BitWizAction           BitField = 1 << 32 // WIZ_ACTION
	BitPolymerBehavior     BitField = 1 << 33 // POLYMER_BEHAVIOR
	BitPolymer             BitField = 1 << 34 // POLYMER
	BitCustomElement       BitField = 1 << 35 // CUSTOM_ELEMENT
	BitMixinClass               BitField = 1 << 36 // MIXIN_CLASS
	BitMixinFunction            BitField = 1 << 37 // MIXIN_FUNCTION
	BitSassGeneratedCssTs       BitField = 1 << 38 // SASS_GENERATED_CSS_TS
	BitClosureUnawareCode       BitField = 1 << 39 // CLOSURE_UNAWARE_CODE
	BitProvideGoog              BitField = 1 << 40 // PROVIDE_GOOG
	BitProvideAlreadyProvided   BitField = 1 << 41 // PROVIDE_ALREADY_PROVIDED
	BitWizCallback              BitField = 1 << 42 // WIZ_CALLBACK
	BitLogTypeInCompiler        BitField = 1 << 43 // LOG_TYPE_IN_COMPILER
)

// IdGenerator mirrors the IdGenerator enum used with @idGenerator in Java.
type IdGenerator int

const (
	IdGeneratorUnique     IdGenerator = iota // plain @idGenerator
	IdGeneratorConsistent                    // @idGenerator {consistent}
	IdGeneratorStable                        // @idGenerator {stable}
	IdGeneratorXid                           // @idGenerator {xid}
	IdGeneratorMapped                        // @idGenerator {mapped}
)

// MarkerPosition records a source range for a single piece of marker information
// (annotation keyword, param name, description, or type expression).
// Mirrors SourcePosition<T> in JSDocInfo.java.
type MarkerPosition struct {
	StartLine int // 1-indexed
	StartCol  int // 0-indexed
	EndLine   int
	EndCol    int
	Item      string // the string item (annotation name, description, or type text)
}

// Marker mirrors JSDocInfo.Marker — source-position-tagged record emitted once
// per block tag when documentation collection is enabled. Callers that need
// precise error positions (e.g. JSC_UNRECOGNIZED_TYPE_ERROR) walk GetMarkers().
type Marker struct {
	Annotation  *MarkerPosition // the annotation keyword (e.g. "param")
	NameNode    *MarkerPosition // the name node (for @param only)
	Description *MarkerPosition // the trailing description text
	Type        *MarkerPosition // the {type} expression, if any
	TypeHasBrackets bool        // mirrors TypePosition.hasBrackets()
}

// BoundedTemplatePosition records the source position of a `@template {B} T`
// annotation. Used by TypeCheck to emit JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES.
// Port of RhinoErrorReporter.UNSUPPORTED_BOUNDED_GENERIC_TYPES + the
// addParserWarning(JSDOC_TEMPLATE_BOUNDEDGENERICS_USED, lineno, charno) call
// in JsDocInfoParser.java's TEMPLATE case.
type BoundedTemplatePosition struct {
	Line int // 1-indexed absolute file line of the @template annotation
	Col  int // 0-indexed column of the @ character
}

// InvalidParamPosition records the source position and raw name of an
// invalid @param annotation. Port of addParserWarning(JSDOC_INVALID_PARAM_NAME).
type InvalidParamPosition struct {
	Line int
	Col  int
	Name string
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
	// originalVisibility mirrors JSDocInfo.getOriginalVisibility() — the
	// visibility value before any merging or inference that may change visibility.
	originalVisibility    Visibility
	hasOriginalVisibility bool
	enumExpr        *JSTypeExpression // non-nil when @enum was present
	enumMarked      bool              // true when @enum appeared (even with no type)
	extendsExpr     *JSTypeExpression
	extendsAll      []*JSTypeExpression // every @extends encountered (interfaces can list many)
	implementsAll   []*JSTypeExpression
	typedefExpr     *JSTypeExpression
	templateNames            []string
	ttlTemplateNames         []string // @template T := … =: transformation names
	boundedTemplatePositions []BoundedTemplatePosition
	invalidParamPositions    []InvalidParamPosition
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

	// --- extended fields: structural port of JSDocInfo property bag ---

	// throwsAnnotations mirrors the THROWS_ANNOTATIONS property.
	// Each entry is the raw text of one @throws annotation.
	throwsAnnotations []string

	// closurePrimitiveId mirrors the CLOSURE_PRIMITIVE_ID property.
	closurePrimitiveId string

	// modifies mirrors the MODIFIES property — set of strings from @modifies{…}.
	modifies map[string]bool

	// markers mirrors the MARKERS property — source-position-tagged annotation list.
	markers []*Marker

	// bit flags for boolean-ish annotations that don't merit individual fields
	bits BitField

	// idGenerator mirrors the ID_GENERATOR property.
	idGenerator    IdGenerator
	hasIdGenerator bool

	// thisType mirrors the THIS_TYPE property — the type expression from @this.
	thisType *JSTypeExpression

	// description mirrors the DESCRIPTION property — text from @desc.
	description string

	// license mirrors the LICENSE property — text from @license.
	license string

	// blockDescription mirrors the BLOCK_DESCRIPTION property — the leading
	// description block before any tag.
	blockDescription string

	// originalCommentPosition mirrors the ORIGINAL_COMMENT_POSITION property.
	originalCommentPosition int

	// originalCommentString mirrors the SOURCE_COMMENT property.
	originalCommentString string

	// fileoverviewDescription mirrors the FILEOVERVIEW_DESCRIPTION property.
	fileoverviewDescription string

	// authors mirrors the AUTHORS property.
	authors []string

	// sees mirrors the SEES property (@see references).
	sees []string

	// returnDescription mirrors the RETURN_DESCRIPTION property.
	returnDescription string

	// paramDescriptions mirrors PARAMETER_DESCRIPTIONS property.
	paramDescriptions map[string]string

	// meaning mirrors the MEANING property — context for i18n translation.
	meaning string

	// alternateMessageId mirrors the ALTERNATE_MESSAGE_ID property.
	alternateMessageId string
}

// AddSuppression records a `@suppress {name}` entry. Mirrors
// JSDocInfo.addSuppression in Java.
// Go-only: language-adapter — AddSuppression ports JSDocInfo.Builder.recordSuppression(String). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) AddSuppression(name string) {
	if i.suppressions == nil {
		i.suppressions = map[string]bool{}
	}
	i.suppressions[name] = true
}

// HasSuppression reports whether the named suppression was attached
// to this comment. Used by the per-pass diagnostic gates (e.g.
// VariableReferenceCheck honours `@suppress {duplicate}`).
// Go-only: language-adapter — HasSuppression ports getSuppressions().contains(name). Go provides a direct predicate; Java callers use getSuppressions().
func (i *JSDocInfo) HasSuppression(name string) bool {
	return i != nil && i.suppressions != nil && i.suppressions[name]
}

// Suppressions returns the full set of suppression names. Mirrors
// JSDocInfo.getSuppressions() in Java, used by CheckJSDoc.validateSuppress.
// Go-only: language-adapter — Suppressions ports JSDocInfo.getSuppressions(). Go exposes the map directly; Java returns an ImmutableSet.
func (i *JSDocInfo) Suppressions() map[string]bool {
	if i == nil {
		return nil
	}
	return i.suppressions
}

// IsStatic marks a class member as belonging to the constructor's
// static slot rather than its instance slot. Mirrors Java's
// MemberFunctionDef.isStaticMember and the analogous field flag.
// Go-only: language-adapter — IsStatic is a Go-specific field on JSDocInfo; Java uses MemberFunctionDef.isStaticMember() on the AST node.
func (i *JSDocInfo) IsStatic() bool   { return i != nil && i.isStatic }
// Go-only: language-adapter — SetStatic is a Go-specific setter for the static member flag; Java uses MemberFunctionDef.isStaticMember() on the AST node.
func (i *JSDocInfo) SetStatic(b bool) { i.isStatic = b }

// TypedefType mirrors JSDocInfo.getTypedefType() — returns the aliased
// type expression captured from `/** @typedef {X} */`.
func (i *JSDocInfo) TypedefType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.typedefExpr
}
// Go-only: language-adapter — SetTypedefType ports JSDocInfo.Builder.recordTypedef(JSTypeExpression).
func (i *JSDocInfo) SetTypedefType(t *JSTypeExpression) { i.typedefExpr = t }

// TemplateNames mirrors JSDocInfo.getTemplateTypeNames() — the in-order
// list of names captured from `@template T` tags.
// Go-only: language-adapter — TemplateNames ports JSDocInfo.getTemplateTypeNames(). Go uses a List; Java returns ImmutableList.
func (i *JSDocInfo) TemplateNames() []string {
	if i == nil {
		return nil
	}
	return i.templateNames
}
// Go-only: language-adapter — AddTemplateName ports JSDocInfo.Builder.recordTemplateTypeName(String). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) AddTemplateName(n string) {
	i.templateNames = append(i.templateNames, n)
}

// TTLTemplateNames returns @template T := … =: transformation template names.
// Mirrors JSDocInfo.getTypeTransformations() key-set in Java.
// Go-only: language-adapter — TTLTemplateNames ports getTypeTransformations() key-set for @template T := … =: transformation names.
func (i *JSDocInfo) TTLTemplateNames() []string {
	if i == nil {
		return nil
	}
	return i.ttlTemplateNames
}
// Go-only: language-adapter — AddTTLTemplateName ports JSDocInfo.Builder.recordTypeTransformation(String, Node). Go stores only the name; Java stores the transformation AST node too.
func (i *JSDocInfo) AddTTLTemplateName(n string) {
	i.ttlTemplateNames = append(i.ttlTemplateNames, n)
}

// BoundedTemplatePositions returns all @template {B} T positions recorded
// during JSDoc parsing. Mirrors the addParserWarning(JSDOC_TEMPLATE_BOUNDEDGENERICS_USED)
// call in JsDocInfoParser.java — each entry becomes a JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES
// diagnostic in TypeCheck.
// Go-only: language-adapter — BoundedTemplatePositions returns the Go-specific list of bounded @template source positions; Java fires JSC warnings directly in JsDocInfoParser.
func (i *JSDocInfo) BoundedTemplatePositions() []BoundedTemplatePosition {
	if i == nil {
		return nil
	}
	return i.boundedTemplatePositions
}
// Go-only: language-adapter — AddBoundedTemplatePosition tracks @template {B} T source positions; Java records the warning directly in JsDocInfoParser.
func (i *JSDocInfo) AddBoundedTemplatePosition(p BoundedTemplatePosition) {
	i.boundedTemplatePositions = append(i.boundedTemplatePositions, p)
}

// Go-only: language-adapter — InvalidParamPositions returns the Go-specific list of invalid @param name positions; Java fires JSC warnings immediately in JsDocInfoParser.
func (i *JSDocInfo) InvalidParamPositions() []InvalidParamPosition {
	if i == nil {
		return nil
	}
	return i.invalidParamPositions
}
// Go-only: language-adapter — AddInvalidParamPosition records invalid @param name positions for JSC_INVALID_PARAM_NAME; Java calls addParserWarning directly.
func (i *JSDocInfo) AddInvalidParamPosition(p InvalidParamPosition) {
	i.invalidParamPositions = append(i.invalidParamPositions, p)
}
// Go-only: language-adapter — UpdateInvalidParamPosition updates a previously recorded invalid @param position; Java fires the warning at parse time.
func (i *JSDocInfo) UpdateInvalidParamPosition(idx, line, col int) {
	if idx >= 0 && idx < len(i.invalidParamPositions) {
		i.invalidParamPositions[idx].Line = line
		i.invalidParamPositions[idx].Col = col
	}
}

// IsInterface / IsRecord mirror JSDocInfo.isInterface() /
// JSDocInfo.usesImplicitMatch() — together they mark a class as a
// structural type.
func (i *JSDocInfo) IsInterface() bool   { return i != nil && i.isInterface }
// Go-only: language-adapter — IsRecord ports isRecordRecorded(). Go uses Go-style naming; Java uses the recorded-predicate pattern.
func (i *JSDocInfo) IsRecord() bool      { return i != nil && i.isRecord }
// Go-only: language-adapter — SetInterface ports JSDocInfo.Builder.recordInterface(). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) SetInterface(b bool) { i.isInterface = b }
// Go-only: language-adapter — SetRecord ports JSDocInfo.Builder.recordInterface() with isRecord=true.
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
// Go-only: language-adapter — SetStruct ports JSDocInfo.Builder.recordStruct(). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) SetStruct(b bool)   { i.isStruct = b }
// Go-only: language-adapter — SetDict ports JSDocInfo.Builder.recordDict(). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) SetDict(b bool)     { i.isDict = b }

// IsAbstract mirrors JSDocInfo.isAbstract() — true when the comment
// carried an `@abstract` annotation.
func (i *JSDocInfo) IsAbstract() bool   { return i != nil && i.isAbstract }
// Go-only: language-adapter — SetAbstract ports JSDocInfo.Builder.recordAbstract(). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) SetAbstract(b bool) { i.isAbstract = b }

// IsNoSideEffects mirrors JSDocInfo.isNoSideEffects() — true when
// the comment carried a `@nosideeffects` annotation.
func (i *JSDocInfo) IsNoSideEffects() bool   { return i != nil && i.isNoSideEffects }
// Go-only: language-adapter — SetNoSideEffects ports JSDocInfo.Builder.recordNoSideEffects().
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
// Go-only: language-adapter — IsEnum ports isEnumParamRecorded(). Go uses isEnum-style naming; Java uses the recorded-predicate pattern.
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
// Go-only: language-adapter — SetEnumParameterType ports JSDocInfo.Builder.recordEnumParameterType(JSTypeExpression).
func (i *JSDocInfo) SetEnumParameterType(t *JSTypeExpression) {
	i.enumExpr = t
	i.enumMarked = true
}
// Go-only: language-adapter — MarkEnum ports the no-arg form of recordEnumParameterType() — marks @enum without capturing a type expression.
func (i *JSDocInfo) MarkEnum() { i.enumMarked = true }

// Implementations mirrors JSDocInfo.getImplementedInterfaces().
// Go-only: language-adapter — Implementations ports JSDocInfo.getImplementedInterfaces(). Go uses a list; Java returns an ImmutableList.
func (i *JSDocInfo) Implementations() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.implementsAll
}
// Go-only: language-adapter — AddImplementation ports recordImplementedInterface(JSTypeExpression). Go uses direct mutation; Java uses Builder.
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
// Go-only: language-adapter — BaseTypes exposes all @extends expressions in interface context; Java's getBaseType() returns only the first. Go splits into BaseType() and BaseTypes().
func (i *JSDocInfo) BaseTypes() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.extendsAll
}

// Go-only: language-adapter — SetBaseType ports JSDocInfo.Builder.recordBaseType(JSTypeExpression). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) SetBaseType(t *JSTypeExpression) {
	if i.extendsExpr == nil {
		i.extendsExpr = t
	}
	i.extendsAll = append(i.extendsAll, t)
}

// IsConstructor mirrors JSDocInfo.isConstructor().
func (i *JSDocInfo) IsConstructor() bool { return i != nil && i.isConstructor }

// Go-only: language-adapter — SetConstructor ports JSDocInfo.Builder.recordConstructor(). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) SetConstructor(b bool) { i.isConstructor = b }

// IsConstant mirrors JSDocInfo.isConstant() — true when the @const
// annotation was present. (Java also checks @define / @final / @desc;
// jsgo adds those as needed.)
func (i *JSDocInfo) IsConstant() bool { return i != nil && i.isConst }

// IsConst is an alias used by CheckAccessControls to match Java's
// JSDocInfo.isConstant() accessor name.
func (i *JSDocInfo) IsConst() bool { return i.IsConstant() }

// SetConst records that the @const annotation was present.
// Go-only: language-adapter — SetConst ports JSDocInfo.Builder.recordConstancy(). Go uses a boolean setter; Java uses a Constancy enum.
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
// Go-only: language-adapter — SetDeprecated ports JSDocInfo.Builder.recordDeprecated() + recordDeprecationReason(String).
func (i *JSDocInfo) SetDeprecated(reason string) {
	i.isDeprecated = true
	i.deprecationReason = reason
}

// IsFinal / SetFinal mirror JSDocInfo.isFinal().
func (i *JSDocInfo) IsFinal() bool    { return i != nil && i.isFinal }
// Go-only: language-adapter — SetFinal ports JSDocInfo.Builder.recordFinality(). Go uses a boolean setter; Java has no args.
func (i *JSDocInfo) SetFinal(b bool)  { i.isFinal = b }

// IsCastHoisted / SetCastHoisted track whether this JSDocInfo was hoisted
// from a `/** @type {T} */ (expr)` type-cast onto its enclosing declaration
// node. Go's parser lacks a CAST AST node (Java has one); this flag lets
// visitVar skip the initializer subtype check when the declared type came
// from a cast (Java uses the CAST node's type as the effective RHS type).
// Go-only: language-adapter — IsCastHoisted is Go-specific: Java represents cast annotations as CAST AST nodes; Go hoists the JSDocInfo to the enclosing declaration.
func (i *JSDocInfo) IsCastHoisted() bool   { return i != nil && i.isCastHoisted }
// Go-only: language-adapter — SetCastHoisted is Go-specific: Java uses a CAST AST node; Go hoists cast JSDoc to the enclosing declaration.
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

// Go-only: language-adapter — AddParam ports JSDocInfo.Builder.recordParameter(String, JSTypeExpression). Go uses direct mutation; Java uses Builder.
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
// Go-only: language-adapter — IsRestParam tracks {...T} rest-param annotations; Java tracks these inside JsDocInfoParser, not on JSDocInfo.
func (i *JSDocInfo) IsRestParam(name string) bool {
	return i != nil && i.restParams[name]
}

// MarkRestParam records that the named @param carries the `{...T}`
// rest marker. Called from JsDocInfoParser when the annotation
// has IterRest at the top of the brace.
// Go-only: language-adapter — MarkRestParam marks a param as {...T} rest; Java tracks this inside JsDocInfoParser, not on JSDocInfo.
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
// Go-only: language-adapter — IsOptionalParam tracks {T=} optional-param annotations; Java tracks these inside JsDocInfoParser, not on JSDocInfo.
func (i *JSDocInfo) IsOptionalParam(name string) bool {
	return i != nil && i.optionalParams[name]
}

// MarkOptionalParam records that the named @param carries `{T=}`.
// Go-only: language-adapter — MarkOptionalParam marks a param as {T=} optional; Java tracks this inside JsDocInfoParser, not on JSDocInfo.
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

// Go-only: language-adapter — SetReturnType ports JSDocInfo.Builder.recordReturnType(JSTypeExpression).
func (i *JSDocInfo) SetReturnType(t *JSTypeExpression) { i.returnExpr = t }

// HasThisType mirrors JSDocInfo.hasThisType() — true when @this was present.
func (i *JSDocInfo) HasThisType() bool { return i != nil && i.hasThisType }

// SetHasThisType records that @this appeared in the comment.
// Go-only: language-adapter — SetHasThisType ports the presence-flag portion of recordThisType(JSTypeExpression).
func (i *JSDocInfo) SetHasThisType(b bool) { i.hasThisType = b }

// IsOverride mirrors JSDocInfo.isOverride() — true when @override was present.
func (i *JSDocInfo) IsOverride() bool { return i != nil && i.isOverride }

// SetOverride records that @override appeared in the comment.
// Go-only: language-adapter — SetOverride ports JSDocInfo.Builder.recordOverride(). Go uses direct mutation; Java uses Builder.
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
// Go-only: language-adapter — SetLendsName ports JSDocInfo.Builder.recordLends(String).
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

// --- originalVisibility accessors ---

// GetOriginalVisibility mirrors JSDocInfo.getOriginalVisibility() — returns the
// visibility annotation before any propagation or inference overwrites it.
// Returns VisibilityInherited when none was explicitly set.
// Go-only: language-adapter — GetOriginalVisibility ports JSDocInfo.getOriginalVisibility() — returns the visibility before any inference-driven merging.
func (i *JSDocInfo) GetOriginalVisibility() Visibility {
	if i == nil || !i.hasOriginalVisibility {
		return VisibilityInherited
	}
	return i.originalVisibility
}

// SetOriginalVisibility records the original visibility. Idempotent: does not
// overwrite an already-recorded value (mirrors Java's idempotency guard).
// Go-only: language-adapter — SetOriginalVisibility stores the pre-inference visibility from the JSDocInfo.originalVisibility field; Go exposes as a setter because there is no Builder.
func (i *JSDocInfo) SetOriginalVisibility(v Visibility) bool {
	if i.hasOriginalVisibility {
		return false
	}
	i.originalVisibility = v
	i.hasOriginalVisibility = true
	return true
}

// --- @throws accessors ---

// GetThrowsAnnotations mirrors JSDocInfo.getThrowsAnnotations() — returns the
// list of raw text strings from each @throws annotation.
func (i *JSDocInfo) GetThrowsAnnotations() []string {
	if i == nil {
		return nil
	}
	return i.throwsAnnotations
}

// AddThrowsAnnotation appends a single @throws annotation text. The first
// element may be an empty string (sentinel) when documentation is not collected
// but the annotation's existence must be recorded for PureFunctionIdentifier.
// Go-only: language-adapter — AddThrowsAnnotation ports JSDocInfo.Builder.recordThrowsAnnotation(String). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) AddThrowsAnnotation(annotation string) {
	i.throwsAnnotations = append(i.throwsAnnotations, annotation)
}

// HasThrows reports whether at least one @throws annotation is present.
// Go-only: language-adapter — HasThrows is a presence check for @throws annotations; Java callers use getThrowsAnnotations() and test for empty.
func (i *JSDocInfo) HasThrows() bool {
	return i != nil && len(i.throwsAnnotations) > 0
}

// --- @closurePrimitive accessors ---

// GetClosurePrimitiveId mirrors JSDocInfo.getClosurePrimitiveId().
func (i *JSDocInfo) GetClosurePrimitiveId() string {
	if i == nil {
		return ""
	}
	return i.closurePrimitiveId
}

// HasClosurePrimitiveId mirrors JSDocInfo.hasClosurePrimitiveId().
func (i *JSDocInfo) HasClosurePrimitiveId() bool {
	return i != nil && i.closurePrimitiveId != ""
}

// SetClosurePrimitiveId records the @closurePrimitive id string.
// Returns false if an id was already recorded (idempotency guard).
// Go-only: language-adapter — SetClosurePrimitiveId ports JSDocInfo.Builder.recordClosurePrimitiveId(String).
func (i *JSDocInfo) SetClosurePrimitiveId(id string) bool {
	if i.closurePrimitiveId != "" {
		return false
	}
	i.closurePrimitiveId = id
	return true
}

// --- @modifies accessors ---

// GetModifies mirrors JSDocInfo.getModifies() — returns the set of
// modification targets declared with @modifies{this} / @modifies{arguments}.
func (i *JSDocInfo) GetModifies() map[string]bool {
	if i == nil {
		return nil
	}
	return i.modifies
}

// ModifiesThis mirrors JSDocInfo.modifiesThis().
func (i *JSDocInfo) ModifiesThis() bool {
	return i != nil && i.modifies["this"]
}

// HasSideEffectsArgumentsAnnotation mirrors hasSideEffectsArgumentsAnnotation().
func (i *JSDocInfo) HasSideEffectsArgumentsAnnotation() bool {
	if i == nil {
		return false
	}
	m := i.modifies
	return len(m) > 1 || (len(m) == 1 && !m["this"])
}

// SetModifies records the @modifies set. Returns false if already set.
// Go-only: language-adapter — SetModifies ports JSDocInfo.Builder.recordModifies(Set<String>).
func (i *JSDocInfo) SetModifies(targets map[string]bool) bool {
	if i.modifies != nil {
		return false
	}
	i.modifies = targets
	return true
}

// --- BitField helpers ---

// CheckBit reports whether the given bit is set.
func (i *JSDocInfo) CheckBit(b BitField) bool {
	return i != nil && (i.bits&b) != 0
}

// SetBit sets or clears the given bit.
func (i *JSDocInfo) SetBit(b BitField, val bool) {
	if val {
		i.bits |= b
	} else {
		i.bits &^= b
	}
}

// --- boolean tag accessors backed by BitField ---

// IsNgInject mirrors JSDocInfo.isNgInject().
func (i *JSDocInfo) IsNgInject() bool { return i.CheckBit(BitNgInject) }

// IsWizaction mirrors JSDocInfo.isWizaction().
func (i *JSDocInfo) IsWizaction() bool { return i.CheckBit(BitWizAction) }

// IsPolymerBehavior mirrors JSDocInfo.isPolymerBehavior().
func (i *JSDocInfo) IsPolymerBehavior() bool { return i.CheckBit(BitPolymerBehavior) }

// IsPolymer mirrors JSDocInfo.isPolymer().
func (i *JSDocInfo) IsPolymer() bool { return i.CheckBit(BitPolymer) }

// IsCustomElement mirrors JSDocInfo.isCustomElement().
func (i *JSDocInfo) IsCustomElement() bool { return i.CheckBit(BitCustomElement) }

// IsMixinClass mirrors JSDocInfo.isMixinClass().
func (i *JSDocInfo) IsMixinClass() bool { return i.CheckBit(BitMixinClass) }

// IsMixinFunction mirrors JSDocInfo.isMixinFunction().
func (i *JSDocInfo) IsMixinFunction() bool { return i.CheckBit(BitMixinFunction) }

// IsExport mirrors JSDocInfo.isExport().
func (i *JSDocInfo) IsExport() bool { return i.CheckBit(BitExport) }

// IsNoCollapse mirrors JSDocInfo.isNoCollapse().
func (i *JSDocInfo) IsNoCollapse() bool { return i.CheckBit(BitNoCollapse) }

// IsNoInline mirrors JSDocInfo.isNoInline().
func (i *JSDocInfo) IsNoInline() bool { return i.CheckBit(BitNoInline) }

// IsImplicitCast mirrors JSDocInfo.isImplicitCast().
func (i *JSDocInfo) IsImplicitCast() bool { return i.CheckBit(BitImplicitCast) }

// IsDefine mirrors JSDocInfo.isDefine().
func (i *JSDocInfo) IsDefine() bool { return i.CheckBit(BitDefine) }

// IsUsedViaDotConstructor mirrors JSDocInfo.isUsedViaDotConstructor().
func (i *JSDocInfo) IsUsedViaDotConstructor() bool { return i.CheckBit(BitUsedViaDotCtor) }

// --- @idGenerator accessors ---

// IsAnyIdGenerator mirrors JSDocInfo.isAnyIdGenerator().
func (i *JSDocInfo) IsAnyIdGenerator() bool { return i != nil && i.hasIdGenerator }

// IsIdGenerator mirrors JSDocInfo.isIdGenerator() — plain @idGenerator.
func (i *JSDocInfo) IsIdGenerator() bool {
	return i != nil && i.hasIdGenerator && i.idGenerator == IdGeneratorUnique
}

// IsConsistentIdGenerator mirrors JSDocInfo.isConsistentIdGenerator().
func (i *JSDocInfo) IsConsistentIdGenerator() bool {
	return i != nil && i.hasIdGenerator && i.idGenerator == IdGeneratorConsistent
}

// IsStableIdGenerator mirrors JSDocInfo.isStableIdGenerator().
func (i *JSDocInfo) IsStableIdGenerator() bool {
	return i != nil && i.hasIdGenerator && i.idGenerator == IdGeneratorStable
}

// IsXidGenerator mirrors JSDocInfo.isXidGenerator().
func (i *JSDocInfo) IsXidGenerator() bool {
	return i != nil && i.hasIdGenerator && i.idGenerator == IdGeneratorXid
}

// IsMappedIdGenerator mirrors JSDocInfo.isMappedIdGenerator().
func (i *JSDocInfo) IsMappedIdGenerator() bool {
	return i != nil && i.hasIdGenerator && i.idGenerator == IdGeneratorMapped
}

// SetIdGenerator records the @idGenerator kind. Returns false if already set.
// Go-only: language-adapter — SetIdGenerator ports the recordIdGenerator / recordConsistentIdGenerator / recordStableIdGenerator family.
func (i *JSDocInfo) SetIdGenerator(kind IdGenerator) bool {
	if i.hasIdGenerator {
		return false
	}
	i.idGenerator = kind
	i.hasIdGenerator = true
	return true
}

// --- Marker accessors ---

// GetMarkers mirrors JSDocInfo.getMarkers() — returns all source-position-tagged
// annotation markers collected during parsing. Only populated when documentation
// collection is enabled.
func (i *JSDocInfo) GetMarkers() []*Marker {
	if i == nil {
		return nil
	}
	return i.markers
}

// AddMarker appends a new Marker and returns a pointer to it so the caller
// can fill in position fields incrementally (mirrors Builder.addMarker()).
func (i *JSDocInfo) AddMarker() *Marker {
	m := &Marker{}
	i.markers = append(i.markers, m)
	return m
}

// --- clone helpers ---

// CloneAndReset returns a shallow clone of the JSDocInfo with the type fields
// reset (typeExpr, returnExpr, params). Mirrors toBuilder / cloneAndReplaceTypeNames
// for the narrow subset jsgo needs: only the structural (non-type) annotations
// are preserved. Used by passes that rewrite type annotations.
// Go-only: language-adapter — CloneAndReset ports JSDocInfo.cloneAndReset() — used by JsDocInfoParser to clone state while resetting mutable fields.
func (i *JSDocInfo) CloneAndReset() *JSDocInfo {
	if i == nil {
		return nil
	}
	cp := *i
	cp.typeExpr = nil
	cp.returnExpr = nil
	cp.params = nil
	cp.paramNames = nil
	cp.restParams = nil
	cp.optionalParams = nil
	return &cp
}

// CloneOn returns a copy of this JSDocInfo attached to the caller's intent to
// move it to a different AST node (structural copy, no deep-clone of maps).
// Mirrors the pattern in JSDocInfo.toBuilder() where props are cloned shallowly.
// Go-only: language-adapter — CloneOn is a Go-specific clone helper for transferring JSDocInfo to a different node; Java assigns the JSDocInfo reference directly.
func (i *JSDocInfo) CloneOn() *JSDocInfo {
	if i == nil {
		return nil
	}
	cp := *i
	return &cp
}

// --- additional accessor ports from JSDocInfo.java ---

// HasConstAnnotation mirrors JSDocInfo.hasConstAnnotation() — true when
// the @const annotation is explicitly present (distinct from isConstant()
// which also checks @define and @final).
func (i *JSDocInfo) HasConstAnnotation() bool {
	return i != nil && i.isConst
}

// IsConstructorOrInterface mirrors JSDocInfo.isConstructorOrInterface() —
// true when any of @constructor, @interface, or @record is present.
func (i *JSDocInfo) IsConstructorOrInterface() bool {
	if i == nil {
		return false
	}
	return i.isConstructor || i.isInterface || i.isRecord
}

// ContainsDeclarationExcludingTypelessConst mirrors
// JSDocInfo.containsDeclarationExcludingTypelessConst() — returns true when
// the comment carries a declaration-forming annotation other than a bare @const.
func (i *JSDocInfo) ContainsDeclarationExcludingTypelessConst() bool {
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
		len(i.implementsAll) > 0 ||
		i.hasThisType ||
		i.hasVisibility ||
		i.isOverride ||
		i.isNoSideEffects ||
		i.CheckBit(BitDefine)
}

// ContainsTypeDeclaration mirrors JSDocInfo.containsTypeDeclaration() — true
// when any of @type, @return, @enum, @typedef, @this, or @param is present.
func (i *JSDocInfo) ContainsTypeDeclaration() bool {
	if i == nil {
		return false
	}
	return i.typeExpr != nil ||
		i.returnExpr != nil ||
		i.enumMarked ||
		i.typedefExpr != nil ||
		i.hasThisType ||
		len(i.params) > 0 ||
		i.extendsExpr != nil
}

// ContainsTypeDefinition mirrors JSDocInfo.containsTypeDefinition() — true
// when the comment defines a new type (@constructor, @interface, @enum, @typedef).
func (i *JSDocInfo) ContainsTypeDefinition() bool {
	if i == nil {
		return false
	}
	return i.isConstructor || i.isInterface || i.enumMarked || i.typedefExpr != nil
}

// HasEnumParameterType mirrors JSDocInfo.hasEnumParameterType() — true when
// @enum was present with an explicit type argument.
func (i *JSDocInfo) HasEnumParameterType() bool {
	return i != nil && i.enumExpr != nil
}

// HasTypedefType mirrors JSDocInfo.hasTypedefType() — true when @typedef is present.
func (i *JSDocInfo) HasTypedefType() bool {
	return i != nil && i.typedefExpr != nil
}

// HasReturnType mirrors JSDocInfo.hasReturnType() — true when @return is present.
func (i *JSDocInfo) HasReturnType() bool {
	return i != nil && i.returnExpr != nil
}

// HasTypeInformation mirrors JSDocInfo.hasTypeInformation() — true when any
// of @type, @typedef, @enum, or @return annotations is present.
func (i *JSDocInfo) HasTypeInformation() bool {
	if i == nil {
		return false
	}
	return i.typeExpr != nil || i.typedefExpr != nil || i.enumMarked || i.returnExpr != nil
}

// HasParameter mirrors JSDocInfo.hasParameter(String) — true when the named
// parameter appears in the comment's @param list.
func (i *JSDocInfo) HasParameter(name string) bool {
	if i == nil {
		return false
	}
	_, ok := i.params[name]
	return ok
}

// HasParameterType mirrors JSDocInfo.hasParameterType(String) — true when the
// named parameter has an attached type expression.
func (i *JSDocInfo) HasParameterType(name string) bool {
	if i == nil {
		return false
	}
	return i.params[name] != nil
}

// GetParameterCount mirrors JSDocInfo.getParameterCount() — returns the number
// of @param annotations recorded.
func (i *JSDocInfo) GetParameterCount() int {
	if i == nil {
		return 0
	}
	return len(i.params)
}

// GetParameterNameAt mirrors JSDocInfo.getParameterNameAt(int) — returns the
// nth parameter name in declaration order, or "" if out of range.
func (i *JSDocInfo) GetParameterNameAt(index int) string {
	if i == nil || index < 0 || index >= len(i.paramNames) {
		return ""
	}
	return i.paramNames[index]
}

// HasBaseType mirrors JSDocInfo.hasBaseType() — true when @extends is present.
func (i *JSDocInfo) HasBaseType() bool {
	return i != nil && i.extendsExpr != nil
}

// GetImplementedInterfaceCount mirrors JSDocInfo.getImplementedInterfaceCount().
func (i *JSDocInfo) GetImplementedInterfaceCount() int {
	if i == nil {
		return 0
	}
	return len(i.implementsAll)
}

// GetExtendedInterfaces mirrors JSDocInfo.getExtendedInterfaces() — returns
// all @extends expressions (interfaces can list multiple supers).
func (i *JSDocInfo) GetExtendedInterfaces() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.extendsAll
}

// GetExtendedInterfacesCount mirrors JSDocInfo.getExtendedInterfacesCount().
func (i *JSDocInfo) GetExtendedInterfacesCount() int {
	if i == nil {
		return 0
	}
	return len(i.extendsAll)
}

// HasModifies mirrors JSDocInfo.hasModifies() — true when @modifies is present.
func (i *JSDocInfo) HasModifies() bool {
	return i != nil && i.modifies != nil
}

// HasFileOverview mirrors JSDocInfo.hasFileOverview() — true when @fileoverview
// bit is set.
func (i *JSDocInfo) HasFileOverview() bool {
	return i.CheckBit(BitFileoverview)
}

// IsExterns mirrors JSDocInfo.isExterns().
func (i *JSDocInfo) IsExterns() bool { return i.CheckBit(BitExterns) }

// IsNoCompile mirrors JSDocInfo.isNoCompile().
func (i *JSDocInfo) IsNoCompile() bool { return i.CheckBit(BitNoCompile) }

// IsNoCoverage mirrors JSDocInfo.isNoCoverage().
func (i *JSDocInfo) IsNoCoverage() bool { return i.CheckBit(BitNoCoverage) }

// IsTypeSummary mirrors JSDocInfo.isTypeSummary().
func (i *JSDocInfo) IsTypeSummary() bool { return i.CheckBit(BitTypeSummary) }

// IsPureOrBreakMyCode mirrors JSDocInfo.isPureOrBreakMyCode().
func (i *JSDocInfo) IsPureOrBreakMyCode() bool { return i.CheckBit(BitPureOrBreakMyCode) }

// IsCollapsibleOrBreakMyCode mirrors JSDocInfo.isCollapsibleOrBreakMyCode().
func (i *JSDocInfo) IsCollapsibleOrBreakMyCode() bool {
	return i.CheckBit(BitCollapsibleOrBreak)
}

// MakesUnrestricted mirrors JSDocInfo.makesUnrestricted() — true when the
// @unrestricted annotation was present.
func (i *JSDocInfo) MakesUnrestricted() bool { return i.CheckBit(BitUnrestricted) }

// IsInlineType mirrors JSDocInfo.isInlineType() — true when the type
// annotation was inlined (e.g. `/** @type {T} */ var x`).
func (i *JSDocInfo) IsInlineType() bool { return i.CheckBit(BitInlineType) }

// IsNgInjectRecorded mirrors JSDocInfo.Builder.isNgInjectRecorded().
func (i *JSDocInfo) IsNgInjectRecorded() bool { return i.IsNgInject() }

// IsWizactionRecorded mirrors JSDocInfo.Builder.isWizactionRecorded().
func (i *JSDocInfo) IsWizactionRecorded() bool { return i.IsWizaction() }

// IsRequireInlining mirrors JSDocInfo.isRequireInlining().
func (i *JSDocInfo) IsRequireInlining() bool { return i.CheckBit(BitRequireInlining) }

// IsEncourageInlining mirrors JSDocInfo.isEncourageInlining().
func (i *JSDocInfo) IsEncourageInlining() bool { return i.CheckBit(BitEncourageInlining) }

// IsNoDts mirrors JSDocInfo.isNoDts().
func (i *JSDocInfo) IsNoDts() bool { return i.CheckBit(BitNoDts) }

// GetTemplateTypes mirrors JSDocInfo.getTemplateTypes() — returns the mapping
// of template names to their bound type expressions (nil for unbounded). The
// Go port stores names as a flat slice; we return a map with nil values.
func (i *JSDocInfo) GetTemplateTypes() map[string]*JSTypeExpression {
	if i == nil || len(i.templateNames) == 0 {
		return nil
	}
	m := make(map[string]*JSTypeExpression, len(i.templateNames))
	for _, name := range i.templateNames {
		m[name] = nil // bound is not tracked at this tier
	}
	return m
}

// GetTypeTransformations mirrors JSDocInfo.getTypeTransformations() — returns
// the map of TTL template names to their transformation ASTs. The Go port
// stores names only (no AST), so we return a nil-valued map.
func (i *JSDocInfo) GetTypeTransformations() map[string]interface{} {
	if i == nil || len(i.ttlTemplateNames) == 0 {
		return nil
	}
	m := make(map[string]interface{}, len(i.ttlTemplateNames))
	for _, name := range i.ttlTemplateNames {
		m[name] = nil
	}
	return m
}

// GetThrowsAnnotations mirrors JSDocInfo.getThrowsAnnotations() — alias of
// the existing GetThrowsAnnotations accessor.
func (i *JSDocInfo) GetThrows() []string {
	return i.GetThrowsAnnotations()
}

// OverwriteVisibility mirrors JSDocInfo.Builder.overwriteVisibility() —
// forcibly sets the visibility even if already recorded. Used by passes that
// propagate fileoverview visibility onto individual declarations.
func (i *JSDocInfo) OverwriteVisibility(v Visibility) {
	i.visibility = v
	i.hasVisibility = true
}

// IsAtSignCodePresent mirrors JSDocInfo.isAtSignCodePresent() — checks
// whether the raw comment string contains "@code". The Go port does not
// store the original comment string, so this always returns false.
// Mirrors: JSDocInfo.isAtSignCodePresent() which calls getOriginalCommentString().
func (i *JSDocInfo) IsAtSignCodePresent() bool {
	// The Go port stores the original comment string when documentation is enabled.
	// Check if the stored comment contains "@code".
	if i == nil || i.originalCommentString == "" {
		return false
	}
	return strings.Contains(i.originalCommentString, "@code")
}

// IsDeprecatedBit mirrors JSDocInfo.isDeprecated() using the Bit.DEPRECATED
// path. The Go port uses a field, so delegate.
// Go-only: language-adapter — IsDeprecatedBit ports isDeprecationReasonRecorded(). Go uses a method; Java uses the bit-field helper.
func (i *JSDocInfo) IsDeprecatedBit() bool {
	return i.IsDeprecated()
}

// IsExport mirrors JSDocInfo.isExport().
// Go-only: language-adapter — IsExportBit ports isExportRecorded(). Go uses a method; Java uses the bit-field helper.
func (i *JSDocInfo) IsExportBit() bool {
	return i.IsExport()
}

// ProvideGoog mirrors JSDocInfo.isProvideGoog() — true when @provideGoog
// is present (only in base.js).
func (i *JSDocInfo) IsProvideGoog() bool { return i.CheckBit(BitProvideGoog) }

// IsProvideAlreadyProvided mirrors JSDocInfo.isProvideAlreadyProvided().
func (i *JSDocInfo) IsProvideAlreadyProvided() bool {
	return i.CheckBit(BitProvideAlreadyProvided)
}

// HasThrowsAnnotations mirrors JSDocInfo.getThrowsAnnotations().size() > 0.
// Alias of HasThrows() for API compatibility.
// Go-only: language-adapter — HasThrowsAnnotations is a presence-check alias for HasThrows; Java callers use getThrowsAnnotations().
func (i *JSDocInfo) HasThrowsAnnotations() bool {
	return i.HasThrows()
}

// --- @this type accessors ---

// GetThisType mirrors JSDocInfo.getThisType() — returns the type expression
// from @this, or nil if not present.
func (i *JSDocInfo) GetThisType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.thisType
}

// SetThisType records the type expression for @this. Called by JsDocInfoParser
// when it processes a @this annotation with a type argument.
// Go-only: language-adapter — SetThisType ports JSDocInfo.Builder.recordThisType(JSTypeExpression).
func (i *JSDocInfo) SetThisType(t *JSTypeExpression) {
	i.thisType = t
	i.hasThisType = t != nil
}

// --- @desc / @license / block-description accessors ---

// GetDescription mirrors JSDocInfo.getDescription() — returns text from @desc.
func (i *JSDocInfo) GetDescription() string {
	if i == nil {
		return ""
	}
	return i.description
}

// SetDescription records the @desc annotation text.
func (i *JSDocInfo) SetDescription(desc string) { i.description = desc }

// GetLicense mirrors JSDocInfo.getLicense() — returns text from @license.
func (i *JSDocInfo) GetLicense() string {
	if i == nil {
		return ""
	}
	return i.license
}

// SetLicense records the @license annotation text.
// Go-only: language-adapter — SetLicense ports JSDocInfo.Builder.recordLicense(String).
func (i *JSDocInfo) SetLicense(license string) { i.license = license }

// GetBlockDescription mirrors JSDocInfo.getBlockDescription() — returns the
// leading description block that precedes any @-tags.
func (i *JSDocInfo) GetBlockDescription() string {
	if i == nil {
		return ""
	}
	return i.blockDescription
}

// SetBlockDescription records the leading block description.
// Go-only: language-adapter — SetBlockDescription ports JSDocInfo.Builder.recordBlockDescription(String).
func (i *JSDocInfo) SetBlockDescription(desc string) { i.blockDescription = desc }

// GetFileOverview mirrors JSDocInfo.getFileOverview() — returns text from
// @fileoverview.
func (i *JSDocInfo) GetFileOverview() string {
	if i == nil {
		return ""
	}
	return i.fileoverviewDescription
}

// SetFileOverview records the @fileoverview description text.
// Go-only: language-adapter — SetFileOverview ports JSDocInfo.Builder.recordFileOverview(String).
func (i *JSDocInfo) SetFileOverview(desc string) { i.fileoverviewDescription = desc }

// --- original comment accessors ---

// GetOriginalCommentString mirrors JSDocInfo.getOriginalCommentString().
// Returns the raw comment text if documentation collection is enabled.
func (i *JSDocInfo) GetOriginalCommentString() string {
	if i == nil {
		return ""
	}
	return i.originalCommentString
}

// SetOriginalCommentString records the raw comment text.
// Go-only: language-adapter — SetOriginalCommentString ports JSDocInfo.Builder.recordOriginalCommentString(String).
func (i *JSDocInfo) SetOriginalCommentString(comment string) {
	i.originalCommentString = comment
}

// GetOriginalCommentPosition mirrors JSDocInfo.getOriginalCommentPosition().
func (i *JSDocInfo) GetOriginalCommentPosition() int {
	if i == nil {
		return 0
	}
	return i.originalCommentPosition
}

// SetOriginalCommentPosition records the comment position.
// Go-only: language-adapter — SetOriginalCommentPosition ports JSDocInfo.Builder.recordOriginalCommentPosition(int).
func (i *JSDocInfo) SetOriginalCommentPosition(pos int) {
	i.originalCommentPosition = pos
}

// --- @author / @see accessors ---

// GetAuthors mirrors JSDocInfo.getAuthors() — returns the list of @author values.
func (i *JSDocInfo) GetAuthors() []string {
	if i == nil {
		return nil
	}
	return i.authors
}

// AddAuthor appends an @author value.
// Go-only: language-adapter — AddAuthor ports JSDocInfo.Builder.recordAuthor(String). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) AddAuthor(author string) {
	i.authors = append(i.authors, author)
}

// GetReferences mirrors JSDocInfo.getReferences() — returns @see values.
func (i *JSDocInfo) GetReferences() []string {
	if i == nil {
		return nil
	}
	return i.sees
}

// AddReference appends a @see reference.
// Go-only: language-adapter — AddReference ports JSDocInfo.Builder.recordReference(String). Go uses direct mutation; Java uses Builder.
func (i *JSDocInfo) AddReference(ref string) {
	i.sees = append(i.sees, ref)
}

// --- @return description ---

// GetReturnDescription mirrors JSDocInfo.getReturnDescription() — returns
// the text that follows the @return type expression.
func (i *JSDocInfo) GetReturnDescription() string {
	if i == nil {
		return ""
	}
	return i.returnDescription
}

// SetReturnDescription records the @return description text.
// Go-only: language-adapter — SetReturnDescription ports JSDocInfo.Builder.recordReturnDescription(String).
func (i *JSDocInfo) SetReturnDescription(desc string) { i.returnDescription = desc }

// --- @param description ---

// HasDescriptionForParameter mirrors JSDocInfo.hasDescriptionForParameter(String).
func (i *JSDocInfo) HasDescriptionForParameter(name string) bool {
	if i == nil {
		return false
	}
	_, ok := i.paramDescriptions[name]
	return ok
}

// GetDescriptionForParameter mirrors JSDocInfo.getDescriptionForParameter(String).
func (i *JSDocInfo) GetDescriptionForParameter(name string) string {
	if i == nil {
		return ""
	}
	return i.paramDescriptions[name]
}

// SetDescriptionForParameter records the description for a named parameter.
// Go-only: language-adapter — SetDescriptionForParameter ports JSDocInfo.Builder.recordParameterDescription(String, String).
func (i *JSDocInfo) SetDescriptionForParameter(name, desc string) {
	if i.paramDescriptions == nil {
		i.paramDescriptions = map[string]string{}
	}
	i.paramDescriptions[name] = desc
}

// --- @meaning / @alternateMessageId ---

// GetMeaning mirrors JSDocInfo.getMeaning() — returns the @meaning text.
func (i *JSDocInfo) GetMeaning() string {
	if i == nil {
		return ""
	}
	return i.meaning
}

// SetMeaning records the @meaning annotation text.
// Go-only: language-adapter — SetMeaning ports JSDocInfo.Builder.recordMeaning(String).
func (i *JSDocInfo) SetMeaning(meaning string) { i.meaning = meaning }

// GetAlternateMessageId mirrors JSDocInfo.getAlternateMessageId().
func (i *JSDocInfo) GetAlternateMessageId() string {
	if i == nil {
		return ""
	}
	return i.alternateMessageId
}

// SetAlternateMessageId records the @alternateMessageId text.
// Go-only: language-adapter — SetAlternateMessageId ports JSDocInfo.Builder.recordAlternateMessageId(String).
func (i *JSDocInfo) SetAlternateMessageId(id string) { i.alternateMessageId = id }

// IsDocumentationIncluded mirrors JSDocInfo.isDocumentationIncluded() — true
// when parseDocumentation was called on the builder (documentation collection
// is enabled). The Go port tracks this via a bit.
func (i *JSDocInfo) IsDocumentationIncluded() bool {
	return i.CheckBit(BitIncludeDocumentation)
}

// ParseDocumentation mirrors Builder.parseDocumentation() — marks that
// documentation collection is enabled.
func (i *JSDocInfo) ParseDocumentation() {
	i.SetBit(BitIncludeDocumentation, true)
}

// IsInlineTypeRecorded mirrors JSDocInfo.isInlineType() for the builder
// check path. Returns true when the @type was inlined.
// Go-only: language-adapter — IsInlineTypeRecorded ports JSDocInfo.isInlineTypeRecorded(). Go forwards to the bit-field CheckBit helper.
func (i *JSDocInfo) IsInlineTypeRecorded() bool {
	return i.IsInlineType()
}

// RecordInlineType records that the type annotation was inlined.
func (i *JSDocInfo) RecordInlineType() {
	i.SetBit(BitInlineType, true)
}

// --- additional mirrors from JSDocInfo.java ---

// IsWizcallback mirrors JSDocInfo.isWizcallback() — true when @wizcallback.
func (i *JSDocInfo) IsWizcallback() bool { return i.CheckBit(BitWizCallback) }

// IsSassGeneratedCssTs mirrors JSDocInfo.isSassGeneratedCssTs().
func (i *JSDocInfo) IsSassGeneratedCssTs() bool { return i.CheckBit(BitSassGeneratedCssTs) }

// IsClosureUnawareCode mirrors JSDocInfo.isClosureUnawareCode() — true when
// @closureUnawareCode is present.
func (i *JSDocInfo) IsClosureUnawareCode() bool { return i.CheckBit(BitClosureUnawareCode) }

// GetLogTypeInCompiler mirrors JSDocInfo.getLogTypeInCompiler().
func (i *JSDocInfo) GetLogTypeInCompiler() bool { return i.CheckBit(BitLogTypeInCompiler) }

// GetTsTypes mirrors JSDocInfo.getTsTypes() — returns the list of @tsType texts.
// Not tracked in the Go port; returns empty slice.
func (i *JSDocInfo) GetTsTypes() []string {
	if i == nil {
		return nil
	}
	return nil // not tracked at this tier
}

// GetSuppressionsAndTheirDescription mirrors
// JSDocInfo.getSuppressionsAndTheirDescription() — in the Go port, descriptions
// are not stored; returns a map from each suppression to "".
func (i *JSDocInfo) GetSuppressionsAndTheirDescription() map[string]string {
	if i == nil || i.suppressions == nil {
		return nil
	}
	m := make(map[string]string, len(i.suppressions))
	for name := range i.suppressions {
		m[name] = ""
	}
	return m
}

// HasEnhance mirrors JSDocInfo.hasEnhance() — true when @enhance is present.
// Not yet tracked in the Go port.
func (i *JSDocInfo) HasEnhance() bool { return i.CheckBit(BitEnhancedNamespace) }

// GetEnhance mirrors JSDocInfo.getEnhance() — returns the enhanced namespace.
// Not yet tracked at this tier.
func (i *JSDocInfo) GetEnhance() string {
	if i == nil {
		return ""
	}
	return "" // not tracked at this tier
}

// HasMods mirrors JSDocInfo.hasMods() — true when @mods is present.
// The mods namespace is not tracked in the Go port at this tier.
func (i *JSDocInfo) HasMods() bool {
	// @mods namespace storage not yet implemented in this port.
	return i != nil && i.CheckBit(0) // always false; placeholder for future port
}

// GetMods mirrors JSDocInfo.getMods() — returns the mods namespace.
// The mods namespace is not tracked in the Go port at this tier.
func (i *JSDocInfo) GetMods() string {
	if i == nil {
		return ""
	}
	return "" // not tracked at this tier
}

// GetTypeExpressions mirrors JSDocInfo.getTypeExpressions() — returns all
// JSTypeExpression values recorded on this info.
func (i *JSDocInfo) GetTypeExpressions() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	var result []*JSTypeExpression
	if i.typeExpr != nil {
		result = append(result, i.typeExpr)
	}
	if i.returnExpr != nil {
		result = append(result, i.returnExpr)
	}
	if i.enumExpr != nil {
		result = append(result, i.enumExpr)
	}
	if i.typedefExpr != nil {
		result = append(result, i.typedefExpr)
	}
	if i.thisType != nil {
		result = append(result, i.thisType)
	}
	for _, p := range i.params {
		if p != nil {
			result = append(result, p)
		}
	}
	if i.extendsExpr != nil {
		result = append(result, i.extendsExpr)
	}
	result = append(result, i.implementsAll...)
	return result
}

// ContainsFunctionDeclaration mirrors JSDocInfo.containsFunctionDeclaration() —
// returns true when the comment contains a function-like declaration.
func (i *JSDocInfo) ContainsFunctionDeclaration() bool {
	if i == nil {
		return false
	}
	return i.isConstructor ||
		i.returnExpr != nil ||
		i.hasThisType ||
		len(i.params) > 0 ||
		i.isNoSideEffects
}

// IsUsedViaDotConstructor mirrors JSDocInfo.isUsedViaDotConstructor().
// Alias of the existing accessor.
// Go-only: language-adapter — IsUsedViaDotConstructorBit ports isUsedViaDotConstructorRecorded(). Go uses a method; Java uses the bit-field helper.
func (i *JSDocInfo) IsUsedViaDotConstructorBit() bool {
	return i.IsUsedViaDotConstructor()
}

// ToStringVerbose mirrors JSDocInfo.toStringVerbose() — returns a debug
// representation of the JSDocInfo.
func (i *JSDocInfo) ToStringVerbose() string {
	if i == nil {
		return "JSDocInfo{nil}"
	}
	return "JSDocInfo{...}"
}
