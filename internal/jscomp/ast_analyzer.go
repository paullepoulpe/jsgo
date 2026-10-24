// Port of com/google/javascript/jscomp/AstAnalyzer.java @ closure-compiler v20261015
// Port of com/google/javascript/jscomp/AutoBuilder_AstAnalyzer_Options_Builder.java @ closure-compiler v20261015

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AstAnalyzer holds the logic for answering questions about portions of the
// AST.
//
// What kind of methods should go here?
//
// Methods that answer questions about some portion of the AST and that may
// require global information about the compilation, generally taking at
// least one Node as an argument. For example:
//
//   - Does a node have side effects?
//   - Can we statically determine the value of a node?
//
// What kind of logic should not go here?
//
// Really simple logic that requires no global information, like finding the
// parameter list node of a function, should be in NodeUtil. Logic that
// creates new Nodes or modifies the AST should go in AstFactory.
type AstAnalyzer struct {
	typeRegistry                 *rhino.JSTypeRegistry
	accessorSummary              *AccessorSummary
	useTypesForLocalOptimization bool
	assumeGettersArePure         bool
	assumeKnownBuiltinsArePure   bool
	hasRegexpGlobalReferences    bool
}

// The set of builtin constructors that don't have side effects.
//
// TODO(bradfordcsmith): If all of these are annotated sideefectfree, can we
// drop this list?
var AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS = jutil.NewLinkedHashSet(
	"Array", "Date", "Error", "Object", "RegExp", "XMLHttpRequest")

// A list of built-in object creation or primitive type cast functions that
// can also be called as constructors but lack side-effects.
// TODO(johnlenz): consider adding an extern annotation for this.
var AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS = jutil.NewLinkedHashSet(
	"Object", "Array", "String", "Number", "BigInt", "Boolean", "RegExp", "Error", "Symbol")

var AstAnalyzer_REGEXP_METHODS = jutil.NewLinkedHashSet("test", "exec")

var AstAnalyzer_STRING_REGEXP_METHODS = jutil.NewLinkedHashSet("match", "replace", "search", "split")

// Java: AstAnalyzer(AstAnalyzer.Options, JSTypeRegistry, AccessorSummary)
func NewAstAnalyzer(
	options *AstAnalyzer_Options,
	typeRegistry *rhino.JSTypeRegistry,
	accessorSummary *AccessorSummary) *AstAnalyzer {
	jutil.CheckArgumentMsg(
		options.assumeGettersArePure || accessorSummary != nil,
		"accessorSummary must be provided if assumeGettersArePure is false")
	a := &AstAnalyzer{}
	a.typeRegistry = typeRegistry
	a.accessorSummary = accessorSummary
	a.useTypesForLocalOptimization = options.useTypesForLocalOptimization
	a.assumeGettersArePure = options.assumeGettersArePure
	a.hasRegexpGlobalReferences = options.hasRegexpGlobalReferences
	a.assumeKnownBuiltinsArePure = options.assumeKnownBuiltinsArePure
	return a
}

// AstAnalyzer_Options is the record AstAnalyzer.Options.
type AstAnalyzer_Options struct {
	useTypesForLocalOptimization bool
	assumeGettersArePure         bool
	hasRegexpGlobalReferences    bool
	assumeKnownBuiltinsArePure   bool
}

// Java: AstAnalyzer.Options(boolean, boolean, boolean, boolean)
func NewAstAnalyzer_Options(
	useTypesForLocalOptimization bool,
	assumeGettersArePure bool,
	hasRegexpGlobalReferences bool,
	assumeKnownBuiltinsArePure bool) *AstAnalyzer_Options {
	return &AstAnalyzer_Options{
		useTypesForLocalOptimization: useTypesForLocalOptimization,
		assumeGettersArePure:         assumeGettersArePure,
		hasRegexpGlobalReferences:    hasRegexpGlobalReferences,
		assumeKnownBuiltinsArePure:   assumeKnownBuiltinsArePure,
	}
}

// Java: AstAnalyzer.Options.useTypesForLocalOptimization()
func (o *AstAnalyzer_Options) UseTypesForLocalOptimization() bool {
	return o.useTypesForLocalOptimization
}

// Java: AstAnalyzer.Options.assumeGettersArePure()
func (o *AstAnalyzer_Options) AssumeGettersArePure() bool { return o.assumeGettersArePure }

// Java: AstAnalyzer.Options.hasRegexpGlobalReferences()
func (o *AstAnalyzer_Options) HasRegexpGlobalReferences() bool {
	return o.hasRegexpGlobalReferences
}

// Java: AstAnalyzer.Options.assumeKnownBuiltinsArePure()
func (o *AstAnalyzer_Options) AssumeKnownBuiltinsArePure() bool {
	return o.assumeKnownBuiltinsArePure
}

// Equals is the record's componentwise equals.
//
// Java: AstAnalyzer.Options.equals(Object)
func (o *AstAnalyzer_Options) Equals(other any) bool {
	that, ok := other.(*AstAnalyzer_Options)
	return ok && that != nil &&
		o.useTypesForLocalOptimization == that.useTypesForLocalOptimization &&
		o.assumeGettersArePure == that.assumeGettersArePure &&
		o.hasRegexpGlobalReferences == that.hasRegexpGlobalReferences &&
		o.assumeKnownBuiltinsArePure == that.assumeKnownBuiltinsArePure
}

// HashCode is the record's hashCode (31 * h + Boolean.hashCode(c) per
// component, as java.lang.runtime.ObjectMethods computes it).
//
// Java: AstAnalyzer.Options.hashCode()
func (o *AstAnalyzer_Options) HashCode() int32 {
	h := int32(0)
	for _, c := range [...]bool{
		o.useTypesForLocalOptimization,
		o.assumeGettersArePure,
		o.hasRegexpGlobalReferences,
		o.assumeKnownBuiltinsArePure,
	} {
		bh := int32(1237)
		if c {
			bh = 1231
		}
		h = 31*h + bh
	}
	return h
}

// String is the record's toString.
//
// Java: AstAnalyzer.Options.toString()
func (o *AstAnalyzer_Options) String() string {
	return "Options[useTypesForLocalOptimization=" + jutil.ToJavaString(o.useTypesForLocalOptimization) +
		", assumeGettersArePure=" + jutil.ToJavaString(o.assumeGettersArePure) +
		", hasRegexpGlobalReferences=" + jutil.ToJavaString(o.hasRegexpGlobalReferences) +
		", assumeKnownBuiltinsArePure=" + jutil.ToJavaString(o.assumeKnownBuiltinsArePure) + "]"
}

// AstAnalyzer_Options_Builder is the @AutoBuilder interface
// AstAnalyzer.Options.Builder.
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_AstAnalyzer_Options_Builder (its only implementation) is
// folded into this struct.
type AstAnalyzer_Options_Builder struct {
	useTypesForLocalOptimization bool
	assumeGettersArePure         bool
	hasRegexpGlobalReferences    bool
	assumeKnownBuiltinsArePure   bool
	set_0                        int8
}

// Java: AstAnalyzer.Options.builder()
func AstAnalyzer_Options_Builder_() *AstAnalyzer_Options_Builder {
	return NewAutoBuilder_AstAnalyzer_Options_Builder()
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder()
func NewAutoBuilder_AstAnalyzer_Options_Builder() *AstAnalyzer_Options_Builder {
	return &AstAnalyzer_Options_Builder{}
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder(AstAnalyzer.Options)
func NewAutoBuilder_AstAnalyzer_Options_BuilderFromSource(
	source *AstAnalyzer_Options) *AstAnalyzer_Options_Builder {
	b := &AstAnalyzer_Options_Builder{}
	b.useTypesForLocalOptimization = source.UseTypesForLocalOptimization()
	b.assumeGettersArePure = source.AssumeGettersArePure()
	b.hasRegexpGlobalReferences = source.HasRegexpGlobalReferences()
	b.assumeKnownBuiltinsArePure = source.AssumeKnownBuiltinsArePure()
	b.set_0 = int8(0xf)
	return b
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder.setUseTypesForLocalOptimization(boolean)
func (b *AstAnalyzer_Options_Builder) SetUseTypesForLocalOptimization(
	useTypesForLocalOptimization bool) *AstAnalyzer_Options_Builder {
	b.useTypesForLocalOptimization = useTypesForLocalOptimization
	b.set_0 |= int8(1)
	return b
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder.setAssumeGettersArePure(boolean)
func (b *AstAnalyzer_Options_Builder) SetAssumeGettersArePure(
	assumeGettersArePure bool) *AstAnalyzer_Options_Builder {
	b.assumeGettersArePure = assumeGettersArePure
	b.set_0 |= int8(2)
	return b
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder.setHasRegexpGlobalReferences(boolean)
func (b *AstAnalyzer_Options_Builder) SetHasRegexpGlobalReferences(
	hasRegexpGlobalReferences bool) *AstAnalyzer_Options_Builder {
	b.hasRegexpGlobalReferences = hasRegexpGlobalReferences
	b.set_0 |= int8(4)
	return b
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder.setAssumeKnownBuiltinsArePure(boolean)
func (b *AstAnalyzer_Options_Builder) SetAssumeKnownBuiltinsArePure(
	assumeKnownBuiltinsArePure bool) *AstAnalyzer_Options_Builder {
	b.assumeKnownBuiltinsArePure = assumeKnownBuiltinsArePure
	b.set_0 |= int8(8)
	return b
}

// Java: AutoBuilder_AstAnalyzer_Options_Builder.build()
func (b *AstAnalyzer_Options_Builder) Build() *AstAnalyzer_Options {
	if b.set_0 != 0xf {
		var missing strings.Builder
		if (b.set_0 & 1) == 0 {
			missing.WriteString(" useTypesForLocalOptimization")
		}
		if (b.set_0 & 2) == 0 {
			missing.WriteString(" assumeGettersArePure")
		}
		if (b.set_0 & 4) == 0 {
			missing.WriteString(" hasRegexpGlobalReferences")
		}
		if (b.set_0 & 8) == 0 {
			missing.WriteString(" assumeKnownBuiltinsArePure")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return NewAstAnalyzer_Options(
		b.useTypesForLocalOptimization,
		b.assumeGettersArePure,
		b.hasRegexpGlobalReferences,
		b.assumeKnownBuiltinsArePure)
}

// MayEffectMutableState returns true if the node may create new mutable
// state, or change existing state.
//
// See http://www.xkcd.com/326/ (XKCD Cartoon).
//
// Java: AstAnalyzer.mayEffectMutableState(Node)
func (a *AstAnalyzer) MayEffectMutableState(n rhino.Node) bool {
	return a.CheckForStateChangeHelper(n /* checkForNewObjects= */, true)
}

// MayHaveSideEffects returns true if the node which may have side effects
// when executed. This version default to the "safe" assumptions when the
// compiler object is not provided (RegExp have side-effects, etc).
//
// Java: AstAnalyzer.mayHaveSideEffects(Node)
func (a *AstAnalyzer) MayHaveSideEffects(n rhino.Node) bool {
	return a.CheckForStateChangeHelper(n /* checkForNewObjects= */, false)
}

// FunctionCallHasSideEffects returns true if this function call may have
// side effects.
//
// This method is guaranteed to return true all calls that have
// side-effects, but may also return true for calls that have none.
//
// callNode is the function call node.
//
// Java: AstAnalyzer.functionCallHasSideEffects(Node)
func (a *AstAnalyzer) FunctionCallHasSideEffects(callNode rhino.Node) bool {
	jutil.CheckStateMsg1(
		callNode.IsCall() || callNode.IsTaggedTemplateLit() || callNode.IsOptChainCall(), "%s", callNode)

	if callNode.IsNoSideEffectsCall() {
		return false
	}

	if callNode.IsOnlyModifiesArgumentsCall() && NodeUtil_AllArgsUnescapedLocal(callNode) {
		return false
	}

	callee := callNode.GetFirstChild()

	// Built-in functions with no side effects.
	if callee.IsName() {
		name := callee.GetString()
		if a.assumeKnownBuiltinsArePure && AstAnalyzer_BUILTIN_FUNCTIONS_WITHOUT_SIDEEFFECTS.Contains(name) {
			return false
		}
	} else if callee.IsGetProp() || callee.IsOptChainGetProp() {
		method := callee.GetString()
		if a.assumeKnownBuiltinsArePure {
			var pureMethodCall bool
			if method == "valueOf" {
				pureMethodCall = callNode.HasOneChild()
			} else {
				pureMethodCall = method == "toString" &&
					(callNode.HasOneChild() || callNode.HasTwoChildren())
			}
			if pureMethodCall {
				return false
			}
		}

		if callNode.IsOnlyModifiesThisCall() &&
			NodeUtil_EvaluatesToLocalValue(callee.GetFirstChild()) {
			return false
		}

		// Many common Math functions have no side-effects.
		// TODO(nicksantos): This is a terrible terrible hack, until
		// I create a definitionProvider that understands namespacing.
		if a.assumeKnownBuiltinsArePure &&
			callee.GetFirstChild().IsName() &&
			callee.IsQualifiedName() &&
			callee.GetFirstChild().GetString() == "Math" {
			switch method {
			case "abs",
				"acos",
				"acosh",
				"asin",
				"asinh",
				"atan",
				"atanh",
				"atan2",
				"cbrt",
				"ceil",
				"cos",
				"cosh",
				"exp",
				"expm1",
				"floor",
				"hypot",
				"log",
				"log10",
				"log1p",
				"log2",
				"max",
				"min",
				"pow",
				"round",
				"sign",
				"sin",
				"sinh",
				"sqrt",
				"tan",
				"tanh",
				"trunc":
				return false
			case "random":
				// no parameters
				return !callNode.HasOneChild()
			default:
				// Unknown Math.* function, so fall out of this switch statement.
			}
		}

		if !a.hasRegexpGlobalReferences && a.assumeKnownBuiltinsArePure {
			if callee.GetFirstChild().IsRegExp() && AstAnalyzer_REGEXP_METHODS.Contains(method) {
				return false
			} else if a.IsTypedAsString(callee.GetFirstChild()) {
				// Unlike regexs, string methods don't need to be hosted on a string literal
				// to avoid leaking mutating global state changes, it is just necessary that
				// the regex object can't be referenced.
				param := callee.GetNext()
				if !param.IsNil() {
					if param.IsStringLit() {
						if AstAnalyzer_STRING_REGEXP_METHODS.Contains(method) {
							return false
						}
					} else if param.IsRegExp() {
						if "replace" == method {
							// Assume anything but a string constant has side-effects
							return !param.GetNext().IsStringLit()
						} else if AstAnalyzer_STRING_REGEXP_METHODS.Contains(method) {
							return false
						}
					}
				}
			}
		}
	}

	return true
}

// Java: AstAnalyzer.isTypedAsString(Node)
func (a *AstAnalyzer) IsTypedAsString(n rhino.Node) bool {
	if n.IsStringLit() {
		return true
	}

	if a.useTypesForLocalOptimization {
		color := n.GetColor()
		if color != nil {
			return color.Equals(colors.StandardColors_STRING)
		}
		type_ := n.GetJSType()
		if type_ != nil {
			nativeStringType := a.typeRegistry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
			if type_.Equals(nativeStringType) {
				return true
			}
		}
	}

	return false
}

// CheckForStateChangeHelper returns true if some node in n's subtree
// changes application state. If checkForNewObjects is true, we assume that
// newly created mutable objects (like object literals) change state.
// Otherwise, we assume that they have no side effects.
//
// Java: AstAnalyzer.checkForStateChangeHelper(Node, boolean)
func (a *AstAnalyzer) CheckForStateChangeHelper(n rhino.Node, checkForNewObjects bool) bool {
	parent := n.GetParent()
	// Rather than id which ops may have side effects, id the ones
	// that we know to be safe
	switch n.GetToken() {
	case
		// Throw is a side-effect by definition.
		rhino.Token_THROW,
		// Context switches can conceal side-effects.
		rhino.Token_YIELD,
		rhino.Token_AWAIT,
		rhino.Token_FOR_AWAIT_OF,
		// Enhanced for loops are almost always side-effectful; it's not worth checking them
		// further. Particularly, they represent a kind of assignment op.
		rhino.Token_FOR_OF,
		rhino.Token_FOR_IN,
		// Variable declarations are side-effects.
		rhino.Token_VAR,
		rhino.Token_LET,
		rhino.Token_CONST,
		rhino.Token_EXPORT,
		// import() expressions have side effects
		rhino.Token_DYNAMIC_IMPORT:
		return true

	// The super keyword is a noop on its own.
	case rhino.Token_SUPER:
		return false

	case rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT, rhino.Token_REGEXP:
		return checkForNewObjects || a.CheckForChangeStateInChildren(n, checkForNewObjects)
	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		// Object-rest and object-spread may trigger a getter.
		return !a.assumeGettersArePure || a.CheckForChangeStateInChildren(n, checkForNewObjects)
	case rhino.Token_ITER_REST, rhino.Token_ITER_SPREAD:
		return NodeUtil_IteratesImpureIterable(n) ||
			a.CheckForChangeStateInChildren(n, checkForNewObjects)

	case rhino.Token_NAME:
		// NAMEs have children if the left side of a var/let/const
		// TODO(b/129564961): Consider EXPORT declarations.
		return n.HasChildren()
	case rhino.Token_FUNCTION:
		// Function expressions don't have side-effects, but function
		// declarations change the namespace. Either way, we don't need to
		// check the children, since they aren't executed at declaration time.
		return checkForNewObjects || NodeUtil_IsFunctionDeclaration(n)

	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
		// simply defining a member function, getter, or setter has no side effects
		return false

	case rhino.Token_COMPUTED_PROP:
		switch parent.GetToken() {
		case rhino.Token_CLASS_MEMBERS:
			// This is a computed method.
			return a.CheckForStateChangeHelper(n.GetFirstChild(), checkForNewObjects)

		case rhino.Token_OBJECT_PATTERN:
			// Due to language syntax, only the last child can be an OBJECT_REST.
			// `({ ['thisKey']: target, ...rest} = something())`
			// The presence of `thisKey` affects what properties get put into `rest`.
			return parent.GetLastChild().IsObjectRest() ||
				a.CheckForChangeStateInChildren(n, checkForNewObjects)

		case rhino.Token_OBJECTLIT:
			// Assume that COMPUTED_PROP keys in OBJECTLIT never trigger getters.
			return a.CheckForChangeStateInChildren(n, checkForNewObjects)
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Illegal COMPUTED_PROP parent " + parent.String()))
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		return n.IsStaticMember() &&
			n.HasChildren() &&
			a.CheckForStateChangeHelper(n.GetFirstChild(), checkForNewObjects)
	case rhino.Token_COMPUTED_FIELD_DEF:
		return a.CheckForStateChangeHelper(n.GetFirstChild(), checkForNewObjects) ||
			(n.IsStaticMember() &&
				!n.GetSecondChild().IsNil() &&
				a.CheckForStateChangeHelper(n.GetSecondChild(), checkForNewObjects))
	case rhino.Token_CLASS:
		return checkForNewObjects ||
			NodeUtil_IsClassDeclaration(n) ||
			// Check the extends clause for side effects.
			a.CheckForStateChangeHelper(n.GetSecondChild(), checkForNewObjects) ||
			// Check for class members that are computed properties with side effects.
			a.CheckForStateChangeHelper(n.GetLastChild(), checkForNewObjects)

	case rhino.Token_NEW:
		return checkForNewObjects ||
			a.ConstructorCallHasSideEffects(n) ||
			a.CheckForChangeStateInChildren(n, checkForNewObjects)
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL, rhino.Token_TAGGED_TEMPLATELIT:
		// calls to functions that have no side effects have the no
		// side effect property set.
		return a.FunctionCallHasSideEffects(n) || a.CheckForChangeStateInChildren(n, checkForNewObjects)

	case rhino.Token_CAST,
		rhino.Token_AND,
		rhino.Token_BLOCK,
		rhino.Token_ROOT,
		rhino.Token_EXPR_RESULT,
		rhino.Token_HOOK,
		rhino.Token_IF,
		rhino.Token_CLASS_MEMBERS,
		rhino.Token_PARAM_LIST,
		// Any context that supports DEFAULT_VALUE is already an assignment. The possiblity of a
		// default doesn't itself create a side-effect. Therefore, we prefer to defer the
		// decision.
		rhino.Token_DEFAULT_VALUE,
		rhino.Token_NUMBER,
		rhino.Token_BIGINT,
		rhino.Token_OR,
		rhino.Token_COALESCE,
		rhino.Token_THIS,
		rhino.Token_TRUE,
		rhino.Token_FALSE,
		rhino.Token_NULL,
		rhino.Token_STRINGLIT,
		rhino.Token_SWITCH,
		rhino.Token_TEMPLATELIT_SUB,
		rhino.Token_TRY,
		rhino.Token_EMPTY,
		rhino.Token_TEMPLATELIT,
		rhino.Token_TEMPLATELIT_STRING:
		return a.CheckForChangeStateInChildren(n, checkForNewObjects)

	case rhino.Token_STRING_KEY:
		if parent.IsObjectPattern() {
			// This STRING_KEY names a property being read from.
			// Assumption: GETELEM (via a COMPUTED_PROP) never triggers a getter or setter.
			if a.GetPropertyKind(n.GetString()).HasGetter() {
				return true
			} else if parent.GetLastChild().IsObjectRest() {
				// Due to language syntax, only the last child can be an OBJECT_REST.
				// `({ thisKey: target, ...rest} = something())`
				// The presence of `thisKey` affects what properties get put into `rest`.
				return true
			}
		}
		return a.CheckForChangeStateInChildren(n, checkForNewObjects)

	// Since we can't see what property is accessed we cannot tell whether
	// obj[someProp]/obj?.[someProp] will
	// trigger a getter or setter, and thus could have side effects.
	// We will assume it does not. This introduces some risk of code breakage, but the code
	// size cost of assuming all GETELEM/OPTCHAIN_GETELEM nodes have side effects is
	// completely unacceptable.
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM:
		return a.CheckForChangeStateInChildren(n, checkForNewObjects)

	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		// TODO(b/135640150): Use the parent nodes to determine whether this is a get or set.
		return a.GetPropertyKind(n.GetString()).HasGetterOrSetter() ||
			a.CheckForChangeStateInChildren(n, checkForNewObjects)
	default:
		if NodeUtil_IsSimpleOperator(n) {
			return a.CheckForChangeStateInChildren(n, checkForNewObjects)
		}
		if NodeUtil_IsAssignmentOp(n) {
			return a.CheckAssignmentForChangeState(n, checkForNewObjects)
		}
		return true
	}
}

// Java: AstAnalyzer.checkForChangeStateInChildren(Node, boolean)
func (a *AstAnalyzer) CheckForChangeStateInChildren(n rhino.Node, checkForNewObjects bool) bool {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if a.CheckForStateChangeHelper(c, checkForNewObjects) {
			return true
		}
	}

	return false
}

// Java: AstAnalyzer.checkAssignmentForChangeState(Node, boolean)
func (a *AstAnalyzer) CheckAssignmentForChangeState(assignOp rhino.Node, checkForNewObjects bool) bool {
	assignTarget := assignOp.GetFirstChild()

	if assignTarget.IsName() {
		return true
	}

	assignValue := assignOp.GetLastChild()
	// Assignments will have side effects if
	// a) The RHS has side effects, or
	// b) The LHS has side effects, or
	// c) A name on the LHS will exist beyond the life of this statement.
	if a.CheckForStateChangeHelper(assignTarget, checkForNewObjects) ||
		a.CheckForStateChangeHelper(assignValue, checkForNewObjects) {
		return true
	}

	if !NodeUtil_IsNormalGet(assignTarget) {
		// TODO(johnlenz): remove this code and make this an exception. This
		// is here only for legacy reasons, the AST is not valid but
		// preserve existing behavior.
		return !NodeUtil_IsLiteralValue(assignTarget, true)
	}
	// If the object being assigned to is a local object, don't
	// consider this a side-effect as it can't be referenced
	// elsewhere.  Don't do this recursively as the property might
	// be an alias of another object, unlike a literal below.
	current := assignTarget.GetFirstChild()
	if NodeUtil_EvaluatesToLocalValue(current) {
		return false
	}

	// A literal value as defined by "isLiteralValue" is guaranteed
	// not to be an alias, or any components which are aliases of
	// other objects.
	// If the root object is a literal don't consider this a
	// side-effect.
	for NodeUtil_IsNormalGet(current) {
		current = current.GetFirstChild()
	}

	return !NodeUtil_IsLiteralValue(current, true)
}

// ConstructorCallHasSideEffects answers: do calls to this constructor have
// side effects?
//
// newNode is the constructor call node.
//
// Java: AstAnalyzer.constructorCallHasSideEffects(Node)
func (a *AstAnalyzer) ConstructorCallHasSideEffects(newNode rhino.Node) bool {
	jutil.CheckArgumentMsg1(newNode.IsNew(), "Expected NEW node, got %s", newNode.GetToken())

	if newNode.IsNoSideEffectsCall() {
		return false
	}

	// allArgsUnescapedLocal() is actually confirming that all of the arguments are literals or
	// values created at the point they are passed in to the call and are not saved anywhere in the
	// calling scope.
	// TODO(bradfordcsmith): It would be good to rename allArgsUnescapedLocal() to something
	// that makes this clearer.
	if newNode.IsOnlyModifiesArgumentsCall() && NodeUtil_AllArgsUnescapedLocal(newNode) {
		return false
	}

	nameNode := newNode.GetFirstChild()
	return !nameNode.IsName() ||
		!a.assumeKnownBuiltinsArePure ||
		!AstAnalyzer_CONSTRUCTORS_WITHOUT_SIDE_EFFECTS.Contains(nameNode.GetString())
}

// NodeTypeMayHaveSideEffects returns true if the current node's type
// implies side effects.
//
// This is a non-recursive version of the may have side effects check; used
// to check wherever the current node's type is one of the reasons why a
// subtree has side effects.
//
// Java: AstAnalyzer.nodeTypeMayHaveSideEffects(Node)
func (a *AstAnalyzer) NodeTypeMayHaveSideEffects(n rhino.Node) bool {
	if NodeUtil_IsAssignmentOp(n) {
		return true
	}

	switch n.GetToken() {
	case rhino.Token_DELPROP, rhino.Token_DEC, rhino.Token_INC, rhino.Token_YIELD, rhino.Token_THROW,
		rhino.Token_AWAIT, rhino.Token_DYNAMIC_IMPORT:
		return true
	case rhino.Token_FOR_IN, // assigns to a loop LHS
		rhino.Token_FOR_OF,       // assigns to a loop LHS, runs an iterator
		rhino.Token_FOR_AWAIT_OF: // assigns to a loop LHS, runs an iterator, async operations
		return true
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL, rhino.Token_TAGGED_TEMPLATELIT:
		return a.FunctionCallHasSideEffects(n)
	case rhino.Token_NEW:
		return a.ConstructorCallHasSideEffects(n)

	// A variable definition that assigns a value.
	// TODO(b/129564961): Consider EXPORT declarations.
	case rhino.Token_NAME:
		return n.HasChildren()
	// A destructuring declaration statement or assignment. Technically these might contain no
	// lvalues but that case is rare enough to be ignored.
	case rhino.Token_DESTRUCTURING_LHS:
		return true

	// Object-rest and object-spread may trigger a getter.
	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		return !a.assumeGettersArePure
	case rhino.Token_ITER_REST, rhino.Token_ITER_SPREAD:
		return NodeUtil_IteratesImpureIterable(n)
	case rhino.Token_STRING_KEY:
		if n.GetParent().IsObjectPattern() {
			return a.GetPropertyKind(n.GetString()).HasGetter()
		}
		return false
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		return a.GetPropertyKind(n.GetString()).HasGetterOrSetter()
	default:
		return false
	}
}

// Java: AstAnalyzer.getPropertyKind(String)
func (a *AstAnalyzer) GetPropertyKind(name string) AccessorSummary_PropertyAccessKind {
	if a.assumeGettersArePure {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return a.accessorSummary.GetKind(name)
}
