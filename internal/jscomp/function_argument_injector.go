// Port of com/google/javascript/jscomp/FunctionArgumentInjector.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// FunctionArgumentInjector is a nifty set of functions to deal with the
// issues of replacing function parameters with a set of call argument
// expressions.
type FunctionArgumentInjector struct {
	astAnalyzer *AstAnalyzer
}

// A string to use to represent "this".  Anything that is not a valid
// identifier can be used, so we use "this".
const FunctionArgumentInjector_THIS_MARKER = "this"

var FunctionArgumentInjector_THIS_MARKER_NODE = rhino.IR_Name(FunctionArgumentInjector_THIS_MARKER)

const FunctionArgumentInjector_REST_MARKER = "rest param"

// Java: FunctionArgumentInjector(AstAnalyzer)
func NewFunctionArgumentInjector(astAnalyzer *AstAnalyzer) *FunctionArgumentInjector {
	return &FunctionArgumentInjector{astAnalyzer: astAnalyzer}
}

// Inject, with the map provided, replaces the names with expression trees.
//
// node is the root node of the tree within which to perform the
// substitutions; parent is the parent root node; replacements is the map of
// names to template node trees with which to replace the name Nodes.
// Returns the root node or its replacement.
//
// Java: FunctionArgumentInjector.inject(AbstractCompiler, Node, Node, Map<String, Node>)
func (f *FunctionArgumentInjector) Inject(compiler AbstractCompiler, node *rhino.Node, parent *rhino.Node, replacements *jutil.LinkedHashMap[string, *rhino.Node]) *rhino.Node {
	return f.InjectWithReplaceThis(compiler, node, parent, replacements, true /* replaceThis= */)
}

// Java: FunctionArgumentInjector.inject(AbstractCompiler, Node, Node, Map<String, Node>, boolean)
func (f *FunctionArgumentInjector) InjectWithReplaceThis(
	compiler AbstractCompiler,
	node *rhino.Node,
	parent *rhino.Node,
	replacements *jutil.LinkedHashMap[string, *rhino.Node],
	replaceThis bool) *rhino.Node {
	if node.IsName() {
		replacementTemplate := replacements.Get(node.GetString())
		if replacementTemplate != nil {
			// This should not be replacing declared names.
			jutil.CheckStateMsg1(!(parent.IsFunction() || parent.IsVar() || parent.IsCatch()), "%s", parent)
			// The name may need to be replaced more than once,
			// so we need to clone the node.
			replacement := replacementTemplate.CloneTree()
			node.ReplaceWith(replacement)
			return replacement
		}
	} else if replaceThis && node.IsThis() {
		replacementTemplate := replacements.Get(FunctionArgumentInjector_THIS_MARKER)
		jutil.CheckNotNull(replacementTemplate)
		if !replacementTemplate.IsThis() {
			// The name may need to be replaced more than once,
			// so we need to clone the node.
			replacement := replacementTemplate.CloneTree()
			node.ReplaceWith(replacement)

			// Remove the value.  This isn't required but it ensures that we won't
			// inject side-effects multiple times as it will trigger the null
			// check above if we do.
			if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate) {
				replacements.Remove(FunctionArgumentInjector_THIS_MARKER)
			}

			return replacement
		}
	} else if (node.IsFunction() && !node.IsArrowFunction()) ||
		node.IsMemberFieldDef() ||
		NodeUtil_IsClassStaticBlock(node) {
		// Once we enter another non-arrow function, class field initializer, or static block,
		// the "this" value changes. Don't try to replace it within an inner scope.
		replaceThis = false
	}

	if node.IsComputedFieldDef() {
		key := node.GetFirstChild()
		if key != nil {
			key = f.InjectWithReplaceThis(compiler, key, node, replacements, replaceThis)
			value := key.GetNext()
			if value != nil {
				value = f.InjectWithReplaceThis(compiler, value, node, replacements /* replaceThis= */, false)
			}
		}
		return node
	}

	for c := node.GetFirstChild(); c != nil; c = c.GetNext() {
		// We have to reassign c in case it was replaced, because the removed c's
		// getNext() would no longer be correct.
		c = f.InjectWithReplaceThis(compiler, c, node, replacements, replaceThis)
	}

	return node
}

// GetFunctionCallParameterMap gets a mapping for function parameter names to
// call arguments.
//
// Java: FunctionArgumentInjector.getFunctionCallParameterMap(Node, Node, Supplier<String>)
func (f *FunctionArgumentInjector) GetFunctionCallParameterMap(
	fnNode *rhino.Node, callNode *rhino.Node, safeNameIdSupplier func() string) *jutil.LinkedHashMap[string, *FunctionArgumentInjector_ParamArgPair] {
	jutil.CheckNotNull(fnNode)
	// Create an parameterName -> expression map
	argMap := jutil.NewImmutableMapBuilder[string, *FunctionArgumentInjector_ParamArgPair]()

	// CALL NODE: [ NAME, ARG1, ARG2, ... ]
	cArg := callNode.GetSecondChild()
	if cArg != nil && NodeUtil_IsFunctionObjectCall(callNode) {
		if !fnNode.IsArrowFunction() {
			argMap.Put(FunctionArgumentInjector_THIS_MARKER, NewFunctionArgumentInjector_ParamArgPair(FunctionArgumentInjector_THIS_MARKER_NODE, cArg))
		} else {
			argMap.Put(FunctionArgumentInjector_THIS_MARKER, NewFunctionArgumentInjector_ParamArgPair(FunctionArgumentInjector_THIS_MARKER_NODE, rhino.IR_ThisNode()))
		}
		cArg = cArg.GetNext()
	} else {
		// 'apply' isn't supported yet.
		jutil.CheckStateMsg1(!NodeUtil_IsFunctionObjectApply(callNode), "%s", callNode)
		argMap.Put(
			FunctionArgumentInjector_THIS_MARKER, NewFunctionArgumentInjector_ParamArgPair(FunctionArgumentInjector_THIS_MARKER_NODE, NodeUtil_NewUndefinedNode(callNode)))
	}

	for fnParam := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild(); fnParam != nil; fnParam = fnParam.GetNext() {
		if cArg != nil {
			if fnParam.IsRest() {
				jutil.CheckStateMsg1(fnParam.GetOnlyChild().IsName(), "%s", fnParam.GetOnlyChild())
				array := rhino.IR_Arraylit()
				array.SrcrefTreeIfMissing(cArg)
				for cArg != nil {
					array.AddChildToBack(cArg.CloneTree())
					cArg = cArg.GetNext()
				}
				return argMap.
					Put(
						fnParam.GetOnlyChild().GetString(),
						NewFunctionArgumentInjector_ParamArgPair(fnParam.GetOnlyChild(), array)).
					BuildOrThrow()
			} else {
				jutil.CheckStateMsg1(fnParam.IsName(), "%s", fnParam)
				argMap.Put(fnParam.GetString(), NewFunctionArgumentInjector_ParamArgPair(fnParam, cArg))
			}
			cArg = cArg.GetNext()
		} else { // cArg == null
			if fnParam.IsRest() {
				jutil.CheckStateMsg1(fnParam.GetOnlyChild().IsName(), "%s", fnParam)
				// No arguments for REST parameters
				array := rhino.IR_Arraylit().Srcref(fnParam)
				argMap.Put(
					fnParam.GetOnlyChild().GetString(), NewFunctionArgumentInjector_ParamArgPair(fnParam.GetOnlyChild(), array))
			} else {
				jutil.CheckStateMsg1(fnParam.IsName(), "%s", fnParam)
				srcLocation := callNode
				argMap.Put(
					fnParam.GetString(),
					NewFunctionArgumentInjector_ParamArgPair(fnParam, NodeUtil_NewUndefinedNode(srcLocation)))
			}
		}
	}

	// Add temp names for arguments that don't have named parameters in the
	// called function.
	for cArg != nil {
		uniquePlaceholder := FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier)
		argMap.Put(uniquePlaceholder, NewFunctionArgumentInjector_ParamArgPair(rhino.IR_Name(uniquePlaceholder), cArg))
		cArg = cArg.GetNext()
	}

	return argMap.BuildOrThrow()
}

// FunctionArgumentInjector_GetUniqueAnonymousParameterName: parameter names
// will be name unique when at a later time.
//
// Java: FunctionArgumentInjector.getUniqueAnonymousParameterName(Supplier<String>)
func FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier func() string) string {
	return "JSCompiler_inline_anon_param_" + safeNameIdSupplier()
}

// FindModifiedParameters retrieves a set of names that can not be safely
// substituted in place.
//
// Example:
//
//	function(a) {
//	  a = 0;
//	}
//
// Inlining this without taking precautions would cause the call site value
// to be modified (bad).
//
// Java: FunctionArgumentInjector.findModifiedParameters(Node)
func (f *FunctionArgumentInjector) FindModifiedParameters(fnNode *rhino.Node) *jutil.LinkedHashSet[string] {
	names := FunctionArgumentInjector_GetFunctionParameterSet(fnNode)
	if names.IsEmpty() {
		return names
	}
	return FunctionArgumentInjector_FindModifiedParametersIn(fnNode.GetLastChild(), names, false)
}

// FunctionArgumentInjector_FindModifiedParametersIn checks for uses of the
// named value that imply a pass-by-value parameter is expected. This is used
// to prevent cases like:
//
// function (x) { x=2; return x; }
//
// We don't want "undefined" to be substituted for "x", and get undefined=2
//
// n is the node in question; names is the set of names to check;
// inInnerClosure is whether the inspection is occurring inside an inner
// function or non-static class field initializer.
//
// Java: FunctionArgumentInjector.findModifiedParameters(Node, ImmutableSet<String>, boolean)
func FunctionArgumentInjector_FindModifiedParametersIn(
	n *rhino.Node, names *jutil.LinkedHashSet[string], inInnerClosure bool) *jutil.LinkedHashSet[string] {
	unsafe := jutil.NewLinkedHashSet[string]()
	functionArgumentInjector_collectModifiedParameters(n, names, inInnerClosure, unsafe)
	return unsafe
}

func functionArgumentInjector_collectModifiedParameters(
	n *rhino.Node, names *jutil.LinkedHashSet[string], inInnerClosure bool, unsafe *jutil.LinkedHashSet[string]) {
	if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		// A function parameter can not be replaced with a direct inlined value
		// if it is referred to by an inner function or a non-static class field
		// initializer. The inner function or field initializer can out live the
		// call we are replacing, so it must capture a unique name.  This approach
		// does not work within loop bodies so those are forbidden elsewhere.
		inInnerClosure = true
	}

	if n.IsName() {
		if names.Contains(n.GetString()) && (inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n)) {
			unsafe.Add(n.GetString())
		}
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		functionArgumentInjector_collectModifiedParameters(c, names, inInnerClosure, unsafe)
	}
}

// FunctionArgumentInjector_CanNameValueChange is similar to NodeUtil.isLValue
// except that object properties and array member modification aren't
// important ("o" in "o.a = 2" is still "o" after assignment, where in as "o
// = x", "o" is now "x").
//
// This also looks for the redefinition of a name. function (x) {var x;}
//
// n is the NAME node in question.
//
// Java: FunctionArgumentInjector.canNameValueChange(Node)
func FunctionArgumentInjector_CanNameValueChange(n *rhino.Node) bool {
	if !NodeUtil_IsLValue(n) {
		return false
	}
	if NodeUtil_IsDeclarationLValue(n) {
		declaring := NodeUtil_GetDeclaringParent(n)
		return !declaring.IsConst() && !declaring.IsLet()
	}
	return true
}

// GatherCallArgumentsNeedingTemps updates the set of parameter names that
// correspond to arguments from the call site that require aliases.
//
// If an argument requires a temp (e.g. it is side-effectful, affects mutable
// state or due to other reasons), we need to evaluate it into a temporary at
// the callsite before inlining the function block into that callsite.
// However, that argument's early evaluation could also change the values of
// previous arguments. Hence, in this function we decide to hoist all those
// previous args and evaluate them before this side-effectful argument, even
// if they are not actually affected by that side-effectful argument. This
// simplifies the implementation and doesn't hurt performance.
//
// This method populates the given namesNeedingTemps set with all the
// parameter names that require hoisting with a temp.
//
// fnNode is the FUNCTION node to be inlined; argMap is the argument list for
// the call to fnNode; modifiedParameters is the set of parameters known to
// be modified, which automatically need temps.
//
// Java: FunctionArgumentInjector.gatherCallArgumentsNeedingTemps(AbstractCompiler, Node, ImmutableMap<String, ParamArgPair>, ImmutableSet<String>, CodingConvention)
func (f *FunctionArgumentInjector) GatherCallArgumentsNeedingTemps(
	compiler AbstractCompiler,
	fnNode *rhino.Node,
	argMap *jutil.LinkedHashMap[string, *FunctionArgumentInjector_ParamArgPair],
	modifiedParameters *jutil.LinkedHashSet[string],
	convention CodingConvention) *jutil.LinkedHashSet[string] {
	jutil.CheckArgumentMsg1(fnNode.IsFunction(), "%s", fnNode)
	if argMap.IsEmpty() {
		// No arguments to check, we are done.
		return modifiedParameters
	}
	namesNeedingTemps := jutil.NewLinkedHashSet(modifiedParameters.ToSlice()...)
	block := fnNode.GetLastChild()

	/*
	 * This field holds the parameter name corresponding to the last side-effectful argument. Args
	 * corresponding to all parameters before and including this parameter name would get hoisted
	 * using temps.
	 */
	requiresTempsUpToThisParameterName := ""

	argCount := argMap.Size()
	// We limit the "trivial" bodies to those where there is a single expression or
	// return, the expression is
	isTrivialBody := (!block.HasChildren() ||
		(block.HasOneChild() && !f.BodyMayHaveConditionalCode(block.GetLastChild())))
	hasMinimalParameters := (NodeUtil_IsUndefined(argMap.Get(FunctionArgumentInjector_THIS_MARKER).Arg()) ||
		(fnNode.IsArrowFunction() && argMap.Get(FunctionArgumentInjector_THIS_MARKER).Arg().IsThis())) &&
		argCount <= 2 // this + one parameter

	// Get the list of parameters that may need temporaries due to side-effects lazily.
	var namesAfterSideEffects *jutil.LinkedHashSet[string]
	var parametersWithNamesReferencedBefore *jutil.LinkedHashSet[string]
	ensureParametersThatMayNeedTemps := func() {
		if namesAfterSideEffects != nil {
			return
		}
		paramSet := jutil.NewLinkedHashSet[string]()
		for k := range argMap.Keys() {
			paramSet.Add(k)
		}
		parametersThatMayNeedTemps := f.FindParametersReferencedAfterSideEffect(paramSet, block)
		namesAfterSideEffects = parametersThatMayNeedTemps.ParametersReferencedAfterSideEffect()
		parametersWithNamesReferencedBefore = parametersThatMayNeedTemps.ParametersWithNamesReferencedBeforeParameter()
	}

	astAnalyzer := compiler.GetAstAnalyzer()

	// Check for arguments that are evaluated more than once.
	for parameterName, entryValue := range argMap.All() {
		cArg := entryValue.Arg()
		if namesNeedingTemps.Contains(parameterName) {
			requiresTempsUpToThisParameterName = parameterName
			continue
		}
		// Stores whether this arg need to get hoisted using a temporary
		requiresTemporary := false
		references := NodeUtil_GetNameReferenceCount(block, parameterName)

		argSideEffects := astAnalyzer.MayHaveSideEffects(cArg)
		if !argSideEffects && references == 0 {
			requiresTemporary = false
		} else if isTrivialBody &&
			hasMinimalParameters &&
			references == 1 &&
			// The below line is checking: Can this be affected by side-effects in the function body?
			!(NodeUtil_CanBeSideEffected(cArg) && func() bool {
				ensureParametersThatMayNeedTemps()
				return namesAfterSideEffects.Contains(parameterName)
			}()) &&
			// The below line is checking: Can the function body be affected by the side effects of
			// the argument? If the argument has side effects and the function body contains names
			// that are referenced before the argument, then we need to inline using temporaries. If
			// not, then we can enter this condition's else-if block.
			(!argSideEffects || !func() bool {
				ensureParametersThatMayNeedTemps()
				return parametersWithNamesReferencedBefore.Contains(parameterName)
			}()) {
			// For functions with a trivial body, and where the parameter evaluation order
			// can't change, and there aren't any side-effect before the parameter, we can
			// avoid creating a temporary.
			//
			// This is done to help inline common trivial functions
			// TODO: b/407603216 - return `true` when return expressions don't have
			// references that can be side-effected
			requiresTemporary = false
		} else if astAnalyzer.MayEffectMutableState(cArg) && references > 0 {
			// Note: Mutable arguments should be assigned to temps, as the
			// may be within in a loop:
			//   function x(a) {
			//     for(var i=0; i<0; i++) {
			//       foo(a);
			//     }
			//   x( [] );
			// The parameter in the call to foo should not become "[]".
			requiresTemporary = true
		} else if argSideEffects {
			// even if there are no references, we still need to evaluate the
			// expression if it has side-effects.
			requiresTemporary = true
		} else if NodeUtil_CanBeSideEffected(cArg) && func() bool {
			ensureParametersThatMayNeedTemps()
			return namesAfterSideEffects.Contains(parameterName)
		}() {
			requiresTemporary = true
		} else if references > 1 {
			// Safe is a misnomer, this is a check for "large".
			switch cArg.GetToken() {
			case rhino.Token_NAME:
				name := cArg.GetString()
				// Don't worry about whether this is global or local, just check if it is
				// "exported" in either case.
				requiresTemporary =
					(convention.IsExportedWithLocal(name, true) || convention.IsExportedWithLocal(name, false))
			case rhino.Token_THIS:
				requiresTemporary = false
			case rhino.Token_STRINGLIT:
				requiresTemporary = (jutil.StringLength(cArg.GetString()) >= 2)
			default:
				requiresTemporary = !NodeUtil_IsImmutableValue(cArg)
			}
		}

		if requiresTemporary {
			requiresTempsUpToThisParameterName = parameterName
		}
	}

	if requiresTempsUpToThisParameterName != "" {
		// mark all names upto requiresTempsUptoParameterName as namesNeedingTemps
		for parameterName := range argMap.Keys() {
			if parameterName == FunctionArgumentInjector_THIS_MARKER &&
				(NodeUtil_IsUndefined(argMap.Get(FunctionArgumentInjector_THIS_MARKER).Arg()) ||
					(fnNode.IsArrowFunction() && argMap.Get(FunctionArgumentInjector_THIS_MARKER).Arg().IsThis())) {
				/* When there is no explicit this arg passed into the call, the argMap contains an entry
				 * <"this", undefined Node>. See the `getFunctionCallParameterMap` method.
				 *
				 *  We do not want to add an unnecessary temp for "this", when there wasn't an explicit
				 * "this" passed in at the callsite. That is, we want to transform a call-site like:
				 *  {@code `foo(a,b)`}
				 * to
				 * {@code `let a$inline$0 = a; foo(a$inline$0 ,b)`}
				 * and not
				 * {@code `let this$inline$0 = this; let a$inline$0=a; foo(a$inline$0, b)`}.
				 * Hence we skip generating a temporary for an implicit "this" arg
				 */
				continue
			}
			// TODO: b/298828688 Use `NodeUtil.isImmutableValue` to detect some simple immutable cases
			// and skip generating temps for those as well.

			if parameterName == requiresTempsUpToThisParameterName {
				namesNeedingTemps.Add(parameterName)
				break
			} else {
				namesNeedingTemps.Add(parameterName)
			}
		}
	}
	return namesNeedingTemps
}

// BodyMayHaveConditionalCode: we consider a return or expression trivial if
// it doesn't contain a conditional expression, a function, or a non-static
// class field initializer.
//
// Java: FunctionArgumentInjector.bodyMayHaveConditionalCode(Node)
func (f *FunctionArgumentInjector) BodyMayHaveConditionalCode(n *rhino.Node) bool {
	if !n.IsReturn() && !n.IsExprResult() {
		return true
	}
	return f.MayHaveConditionalCode(n)
}

// MayHaveConditionalCode: we consider an expression trivial if it doesn't
// contain a conditional expression, a function, or a non-static class field
// initializer.
//
// Java: FunctionArgumentInjector.mayHaveConditionalCode(Node)
func (f *FunctionArgumentInjector) MayHaveConditionalCode(n *rhino.Node) bool {
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_FUNCTION,
			rhino.Token_AND,
			rhino.Token_OR,
			rhino.Token_HOOK,
			rhino.Token_COALESCE,
			rhino.Token_OPTCHAIN_CALL,
			rhino.Token_OPTCHAIN_GETELEM,
			rhino.Token_OPTCHAIN_GETPROP:
			return true
		default:
			if NodeUtil_IsNonStaticClassFieldInitializer(c) {
				return true
			}
		}
		if f.MayHaveConditionalCode(c) {
			return true
		}
	}
	return false
}

// FindParametersReferencedAfterSideEffect bootstraps a traversal to look for
// parameters referenced after a non-local side-effect, and parameters with
// names referenced before the parameter is referenced.
//
// NOTE: This assumes no-inner functions.
//
// parameters is the set of parameter names; root is the function code
// block. Returns two sets: Set #1 is the subset of parameters referenced
// after the first seen non-local side-effect. Set #2 is the subset of
// parameters with names referenced before the parameter is referenced.
//
// Java: FunctionArgumentInjector.findParametersReferencedAfterSideEffect(ImmutableSet<String>, Node)
func (f *FunctionArgumentInjector) FindParametersReferencedAfterSideEffect(
	parameters *jutil.LinkedHashSet[string], root *rhino.Node) *FunctionArgumentInjector_SetContainer {

	// TODO(johnlenz): Consider using scope for this.
	locals := jutil.NewLinkedHashSet[string]()
	for p := range parameters.All() {
		locals.Add(p)
	}
	FunctionArgumentInjector_GatherLocalNames(root, locals)

	collector := NewFunctionArgumentInjector_ReferencedAfterSideEffect(
		f, parameters, locals)
	NodeUtil_VisitPostOrderWithPred(root, collector, collector.Apply)
	return NewFunctionArgumentInjector_SetContainer(
		collector.GetParametersReferencedAfterSideEffect(),
		collector.GetParametersWithNamesReferencedBeforeParameter())
}

// FunctionArgumentInjector_SetContainer is the record SetContainer.
type FunctionArgumentInjector_SetContainer struct {
	parametersReferencedAfterSideEffect          *jutil.LinkedHashSet[string]
	parametersWithNamesReferencedBeforeParameter *jutil.LinkedHashSet[string]
}

// Java: FunctionArgumentInjector.SetContainer(ImmutableSet<String>, ImmutableSet<String>)
func NewFunctionArgumentInjector_SetContainer(
	parametersReferencedAfterSideEffect *jutil.LinkedHashSet[string],
	parametersWithNamesReferencedBeforeParameter *jutil.LinkedHashSet[string]) *FunctionArgumentInjector_SetContainer {
	return &FunctionArgumentInjector_SetContainer{
		parametersReferencedAfterSideEffect:          parametersReferencedAfterSideEffect,
		parametersWithNamesReferencedBeforeParameter: parametersWithNamesReferencedBeforeParameter,
	}
}

// Java: FunctionArgumentInjector.SetContainer.parametersReferencedAfterSideEffect()
func (s *FunctionArgumentInjector_SetContainer) ParametersReferencedAfterSideEffect() *jutil.LinkedHashSet[string] {
	return s.parametersReferencedAfterSideEffect
}

// Java: FunctionArgumentInjector.SetContainer.parametersWithNamesReferencedBeforeParameter()
func (s *FunctionArgumentInjector_SetContainer) ParametersWithNamesReferencedBeforeParameter() *jutil.LinkedHashSet[string] {
	return s.parametersWithNamesReferencedBeforeParameter
}

// FunctionArgumentInjector_ReferencedAfterSideEffect collects parameter
// names referenced after a non-local side-effect.
//
// Assumptions:
//
//   - We assume parameters are not modified in the function body (that is
//     checked separately).
//   - There are no inner functions (also checked separately).
//
// As we are trying to replace parameters with there passed in values we are
// interested in anything that may affect those value. So, ignoring changes
// to local variables, we look for things that may affect anything outside
// the local-state. Once such a side-effect is seen any following reference
// to the function parameters are collected. These will need to be assigned
// to temporaries to prevent changes to their value as would have happened
// during the function call.
//
// To properly handle loop structures all references to the function
// parameters are recorded and the decision to keep or throw away those
// references is deferred until exiting the loop structure.
type FunctionArgumentInjector_ReferencedAfterSideEffect struct {
	outer                *FunctionArgumentInjector
	parameters           *jutil.LinkedHashSet[string]
	locals               *jutil.LinkedHashSet[string]
	sideEffectSeen       bool
	parametersReferenced *jutil.LinkedHashSet[string]

	nameNodeHasBeenSeen                          bool
	parametersWithNamesReferencedBeforeParameter *jutil.LinkedHashSet[string]

	loopsEntered int
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect(ImmutableSet<String>, ImmutableSet<String>)
func NewFunctionArgumentInjector_ReferencedAfterSideEffect(
	outer *FunctionArgumentInjector, parameters *jutil.LinkedHashSet[string], locals *jutil.LinkedHashSet[string]) *FunctionArgumentInjector_ReferencedAfterSideEffect {
	r := &FunctionArgumentInjector_ReferencedAfterSideEffect{
		outer:                outer,
		parametersReferenced: jutil.NewLinkedHashSet[string](),
		parametersWithNamesReferencedBeforeParameter: jutil.NewLinkedHashSet[string](),
	}
	r.parameters = parameters
	r.locals = locals
	return r
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.getParametersReferencedAfterSideEffect()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) GetParametersReferencedAfterSideEffect() *jutil.LinkedHashSet[string] {
	return r.parametersReferenced
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.getParametersWithNamesReferencedBeforeParameter()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) GetParametersWithNamesReferencedBeforeParameter() *jutil.LinkedHashSet[string] {
	return r.parametersWithNamesReferencedBeforeParameter
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.apply(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) Apply(node *rhino.Node) bool {
	// Keep track of any loop structures entered.
	if NodeUtil_IsLoopStructure(node) {
		r.loopsEntered++
	}

	// If we have found all the parameters, don't bother looking
	// at the children.
	return !(r.sideEffectSeen && r.parameters.Size() == r.parametersReferenced.Size())
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.inLoop()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) InLoop() bool {
	return r.loopsEntered != 0
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.visit(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) Visit(n *rhino.Node) {
	if n.IsName() {
		name := n.GetString()
		if r.parameters.Contains(name) && r.nameNodeHasBeenSeen {
			// We have seen a name node before this parameter. If this parameter has side effects, it
			// is not safe to inline.
			r.parametersWithNamesReferencedBeforeParameter.Add(name)
		}
		r.nameNodeHasBeenSeen = true
	}

	// If we are exiting a loop.
	if NodeUtil_IsLoopStructure(n) {
		r.loopsEntered--
		if !r.InLoop() && !r.sideEffectSeen {
			// Now that the loops has been fully traversed and
			// no side-effects have been seen, throw away
			// the references seen in them.
			r.parametersReferenced.Clear()
		}
	}

	if !r.sideEffectSeen {
		// Look for side-effects.
		if r.HasNonLocalSideEffect(n) {
			r.sideEffectSeen = true
		}
	}

	// If traversing the nodes of a loop save any references
	// that are seen.
	if r.InLoop() || r.sideEffectSeen {
		// Record references to parameters.
		if n.IsName() {
			name := n.GetString()
			if r.parameters.Contains(name) {
				r.parametersReferenced.Add(name)
			}
		} else if n.IsThis() {
			r.parametersReferenced.Add(FunctionArgumentInjector_THIS_MARKER)
		}
	}
}

// HasNonLocalSideEffect returns whether the node may have non-local
// side-effects.
//
// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.hasNonLocalSideEffect(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) HasNonLocalSideEffect(n *rhino.Node) bool {
	sideEffect := false
	type_ := n.GetToken()
	// Note: Only care about changes to non-local names, specifically
	// ignore VAR declaration assignments.
	if NodeUtil_IsAssignmentOp(n) || type_ == rhino.Token_INC || type_ == rhino.Token_DEC {
		lhs := n.GetFirstChild()
		// Ignore changes to local names.
		if !r.IsLocalName(lhs) {
			sideEffect = true
		}
	} else if type_ == rhino.Token_CALL {
		sideEffect = r.outer.astAnalyzer.FunctionCallHasSideEffects(n)
	} else if type_ == rhino.Token_NEW {
		sideEffect = r.outer.astAnalyzer.ConstructorCallHasSideEffects(n)
	} else if type_ == rhino.Token_DELPROP {
		sideEffect = true
	}

	return sideEffect
}

// IsLocalName returns whether node is a reference to locally declared name.
//
// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.isLocalName(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) IsLocalName(node *rhino.Node) bool {
	if node.IsName() {
		name := node.GetString()
		return r.locals.Contains(name)
	}
	return false
}

// FunctionArgumentInjector_GatherLocalNames gathers any names declared in
// the local scope.
//
// Java: FunctionArgumentInjector.gatherLocalNames(Node, Set<String>)
func FunctionArgumentInjector_GatherLocalNames(n *rhino.Node, names *jutil.LinkedHashSet[string]) {
	if n.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(n) {
			names.Add(n.GetFirstChild().GetString())
		}
		// Don't traverse into inner function scopes;
		return
	} else if n.IsName() {
		switch n.GetParent().GetToken() {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH:
			names.Add(n.GetString())
		default:
		}
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		FunctionArgumentInjector_GatherLocalNames(c, names)
	}
}

// FunctionArgumentInjector_GetFunctionParameterSet gets a set of function
// parameter names.
//
// Java: FunctionArgumentInjector.getFunctionParameterSet(Node)
func FunctionArgumentInjector_GetFunctionParameterSet(fnNode *rhino.Node) *jutil.LinkedHashSet[string] {
	builder := jutil.NewLinkedHashSet[string]()
	for n := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild(); n != nil; n = n.GetNext() {
		if n.IsRest() {
			builder.Add(FunctionArgumentInjector_REST_MARKER)
		} else if n.IsDefaultValue() || n.IsObjectPattern() || n.IsArrayPattern() {
			panic(jutil.NewIllegalStateExceptionMsg("Not supported: " + n.String()))
		} else {
			builder.Add(n.GetString())
		}
	}
	return builder
}

// FunctionArgumentInjector_ParamArgPair stores a pair of function formal
// parameter + argument value at a specific call site.
//
// paramNode is the original formal parameter name, such as x in function
// foo(x) { ... } or, when injecting a value for this,
// FunctionArgumentInjector#THIS_MARKER_NODE; arg is the argument passed as
// `paramNode` at some call site, such as 500 in f(500);.
type FunctionArgumentInjector_ParamArgPair struct {
	paramNode *rhino.Node
	arg       *rhino.Node
}

// Java: FunctionArgumentInjector.ParamArgPair(Node, Node)
func NewFunctionArgumentInjector_ParamArgPair(paramNode *rhino.Node, arg *rhino.Node) *FunctionArgumentInjector_ParamArgPair {
	jutil.CheckNotNull(paramNode)
	jutil.CheckNotNull(arg)
	return &FunctionArgumentInjector_ParamArgPair{paramNode: paramNode, arg: arg}
}

// Java: FunctionArgumentInjector.ParamArgPair.paramNode()
func (p *FunctionArgumentInjector_ParamArgPair) ParamNode() *rhino.Node { return p.paramNode }

// Java: FunctionArgumentInjector.ParamArgPair.arg()
func (p *FunctionArgumentInjector_ParamArgPair) Arg() *rhino.Node { return p.arg }
