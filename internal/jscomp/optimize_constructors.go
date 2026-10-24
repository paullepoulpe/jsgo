// Port of com/google/javascript/jscomp/OptimizeConstructors.java @ closure-compiler v20261015

/*
 * Copyright 2021 The Closure Compiler Authors.
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

// OptimizeConstructors optimizes class declarations by removing explicit
// constructor declarations if the implicit constructor is sufficient.
//
// # The constructor can be removed, if
//
// (0) the class in question defined using ES class syntax
//
// (2) the constructor is empty save for the call to super
//
// (3) either:
//
// - all constructor arguments are forward to the super constructor (via
// rest or arguments)
//
// - super constructor parameters match the number, order, and default
// values of constructor to be removed and the constructor doesn't look for
// additional parameters via references to 'arguments'
//
// When inspecting super class constructors, we don't need to worry about
// (1) reassignment of subclasses or (2) escapes of the superclass. We only
// need to make sure the superclass isn't reassigned. For that we can simply
// check for direct assignments and whether the name is defined in the
// externs as references. The two other reassignment cases that are
// hard/impossible to detect: "globalThis" properties and assignment through
// eval are only only supported with externs.
//
// Note that an alternative exists where the super class could be ignored, if
// all constructor references to the class were known (including direct
// calls, super calls, implicit constructor calls, etc) to validate the
// number of parameters passed. This would require backing off on various
// escapes and is believed it would be both more complicated and less
// effective.
type OptimizeConstructors struct {
	compiler    AbstractCompiler
	astAnalyzer *AstAnalyzer

	// All constructor definition nodes that are to be removed.
	removableConstructors []rhino.Node
}

// Java: OptimizeConstructors(AbstractCompiler)
func NewOptimizeConstructors(compiler AbstractCompiler) *OptimizeConstructors {
	o := &OptimizeConstructors{}
	o.compiler = jutil.CheckNotNull(compiler)
	o.astAnalyzer = compiler.GetAstAnalyzer()
	return o
}

// Process is CompilerPass.process.
//
// Go-only: language-adapter — the two Java overloads process(Node, Node)
// (CompilerPass) and process(Node, Node, ReferenceMap)
// (CallGraphCompilerPass) cannot share a Go name; the second one is
// ProcessWithReferenceMap, and AsCallGraphCompilerPass adapts it.
//
// Java: OptimizeConstructors.process(Node, Node)
func (o *OptimizeConstructors) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(o.compiler.GetLifeCycleStage() == AbstractCompiler_LifeCycleStage_NORMALIZED)

	OptimizeCalls_Builder_().
		SetCompiler(o.compiler).
		SetConsiderExterns(false).
		AddPass(o.AsCallGraphCompilerPass()).
		Build().
		Process(externs, root)
}

// AsCallGraphCompilerPass returns this pass as an
// OptimizeCalls.CallGraphCompilerPass (what Java does by passing `this`).
//
// Go-only: language-adapter — see Process.
func (o *OptimizeConstructors) AsCallGraphCompilerPass() OptimizeCalls_CallGraphCompilerPass {
	return OptimizeCalls_CallGraphCompilerPassFunc(o.ProcessWithReferenceMap)
}

// Java: OptimizeConstructors.process(Node, Node, ReferenceMap)
func (o *OptimizeConstructors) ProcessWithReferenceMap(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	for key, value := range refMap.GetNameReferences() {
		o.AddConstructorsToBeRemoved(key, value)
	}

	// NOTE: Normally, the ReferenceMap must be kept in a consistent state (removing references
	// as they are removed from the AST) so the next pass can reuse the ReferenceMap. However,
	// currently this pass runs by itself so that work is avoided.
	//

	for _, ref := range o.removableConstructors {
		o.RemoveConstructorMethod(ref)
	}
}

// AddConstructorsToBeRemoved iterates over all the references to a symbol.
// There are several interesting references:
//
// -a class definition
//
// -a subclass definition
//
// -class redefinition
//
// -allowed reference (anything not an assignment)
//
// Java: OptimizeConstructors.addConstructorsToBeRemoved(String, ArrayList<Node>)
func (o *OptimizeConstructors) AddConstructorsToBeRemoved(name string, refs *jutil.ArrayList[rhino.Node]) {
	if !OptimizeCalls_MayBeOptimizableName(o.compiler, name) {
		return
	}

	var candidateClassDefinition rhino.Node
	var candidateConstructor rhino.Node

	// Lazily init as most symbols aren't class, and most classes don't have subclasses
	var subclassConstructors []rhino.Node

	for _, n := range refs.Slice() {
		definition := o.GetClassDefinitionOrFunction(n)
		if !definition.IsNil() {
			if !candidateClassDefinition.IsNil() {
				// As a simplification only allow one definition.
				return
			}

			// Be lazy about constructor analysis as it is expected that most classes
			// don't have subclasses.
			candidateClassDefinition = definition

			if candidateClassDefinition.IsClass() {
				// There are two special cases that we can handle just by looking at the extend clause

				// handle no super class
				//    - if no extends clause
				//    - and the constructor body is empty
				//    - and the constructor parameters are side-effect free (no destructuring or
				// side-effect defaults)
				extendsExpr := candidateClassDefinition.GetSecondChild()
				if extendsExpr.IsEmpty() {
					constructor :=
						NodeUtil_GetEs6ClassConstructorMemberFunctionDef(candidateClassDefinition)
					if !constructor.IsNil() {
						fn := constructor.GetLastChild()
						body := fn.GetLastChild()
						if !body.HasChildren() && o.HasRemovableParameterList(constructor) {
							candidateConstructor = constructor
						}
					}
				} else {
					// TODO: handle "Object" super class
					//    - just a call to super and no-side-effects in the formal parameters, or super call
					// arguments
				}
			}
		} else {
			// If this is a reference to the class in an extend clause, then this is a class
			// definition whose constructor should considered for removal.
			if o.IsClassExtendsExpression(n) {
				subclassLiteral := n.GetParent()
				subclassConstructor :=
					NodeUtil_GetEs6ClassConstructorMemberFunctionDef(subclassLiteral)
				if !subclassConstructor.IsNil() &&
					OptimizeConstructors_ConstructorHasRemovableDefinition(subclassConstructor) {

					subclassConstructors = append(subclassConstructors, subclassConstructor)
				}
			} else {
				// Anything assignment that isn't a class definition covered above, is an
				// invalidating assignment.
				if OptimizeConstructors_IsAssigningReference(n) {
					return
				}
			}
		}
	}

	// There is no known class definition
	if candidateClassDefinition.IsNil() {
		return
	}

	if !candidateConstructor.IsNil() {
		o.removableConstructors = append(o.removableConstructors, candidateConstructor)
	}

	// Nothing that invalidated the superclass definition was found, so now check if the
	// candidate subclasses are equivalent to the superclass.

	if subclassConstructors != nil {
		summary := OptimizeConstructors_ClassConstructorSummary_Build(candidateClassDefinition)
		for _, n := range subclassConstructors {
			if summary.IsEquivalentConstructorDefinition(n) {
				o.removableConstructors = append(o.removableConstructors, n)
			}
		}
	}
}

// OptimizeConstructors_IsAssigningReference: is the node in a position in
// the AST so that it might be assigned a value?
//
// Java: OptimizeConstructors.isAssigningReference(Node)
func OptimizeConstructors_IsAssigningReference(n rhino.Node) bool {
	parent := n.GetParent()
	gparent := parent.GetParent()
	switch parent.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		return n.HasChildren() // value assigned
	case rhino.Token_STRING_KEY:
		return gparent.IsObjectPattern()
	case rhino.Token_COMPUTED_PROP:
		return parent.GetLastChild() == n && gparent.IsObjectPattern()
	case rhino.Token_ARRAY_PATTERN,
		rhino.Token_DEFAULT_VALUE, // object or array or function parameter
		rhino.Token_PARAM_LIST,
		rhino.Token_OBJECT_REST,
		rhino.Token_ITER_REST,
		rhino.Token_INC,
		rhino.Token_DEC:
		return true
	case rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_CATCH:
		return parent.GetFirstChild() == n
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_EXPONENT:
		// If n is the target of the assignment (the first child/LHS of the assignment node),
		// then it is an assigning reference.
		return parent.GetFirstChild() == n
	default:
		return false
	}
}

// Java: OptimizeConstructors.isClassExtendsExpression(Node)
func (o *OptimizeConstructors) IsClassExtendsExpression(n rhino.Node) bool {
	parent := n.GetParent()
	return (parent.IsClass() && parent.GetSecondChild() == n)
}

// Java: OptimizeConstructors.getClassDefinitionOrFunction(Node)
func (o *OptimizeConstructors) GetClassDefinitionOrFunction(n rhino.Node) rhino.Node {
	parent := n.GetParent()

	var expr rhino.Node
	if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		expr = parent.GetLastChild()
	} else if n.IsName() && n.HasChildren() {
		expr = n.GetFirstChild()
	} else if parent.IsFunction() && n.IsFirstChildOf(parent) {
		expr = parent
	} else if parent.IsClass() && n.IsFirstChildOf(parent) {
		expr = parent
	} else {
		return rhino.NilNode // Couldn't find a function.
	}

	expr = o.Unwrap(expr)
	if OptimizeConstructors_IsDefinitionClassLiteralOrFunction(expr) {
		return expr
	} else {
		return rhino.NilNode
	}
}

// Java: OptimizeConstructors.unwrap(Node)
func (o *OptimizeConstructors) Unwrap(expr rhino.Node) rhino.Node {
	for expr.IsCast() || expr.IsComma() {
		expr = expr.GetLastChild()
	}
	return expr
}

// Java: OptimizeConstructors.isDefinitionClassLiteralOrFunction(Node)
func OptimizeConstructors_IsDefinitionClassLiteralOrFunction(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		// TODO(b/176208718): ideally this is only return true for normal functions, but it is
		// harmless to include other function types and checking for "normal" function is
		// currently
		// non-trivial.
		return true
	case rhino.Token_CLASS:
		{
			// `class NameNode {`
			// find the constructor
			constructorMemberFunctionDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
			return !constructorMemberFunctionDef.IsNil()
		}
	default:
		return false
	}
}

// OptimizeConstructors_ClassConstructorSummary is an abstraction of a super
// class constructor definition. The goal of this class is to avoid repeated
// analysis of the super class constructor when checking whether a subclass
// constructor is removable.
type OptimizeConstructors_ClassConstructorSummary struct {
	// The number of explicit parameters
	formalParameterCount int
	// "var args" means any number of parameters are allowed
	isVarArgs bool
}

// Java: OptimizeConstructors.ClassConstructorSummary(boolean, int)
func NewOptimizeConstructors_ClassConstructorSummary(isVarArgs bool, formalParameterCount int) *OptimizeConstructors_ClassConstructorSummary {
	return &OptimizeConstructors_ClassConstructorSummary{isVarArgs: isVarArgs, formalParameterCount: formalParameterCount}
}

// Java: OptimizeConstructors.ClassConstructorSummary.build(Node)
func OptimizeConstructors_ClassConstructorSummary_Build(classDefinition rhino.Node) *OptimizeConstructors_ClassConstructorSummary {
	jutil.CheckState(classDefinition.IsClass() || classDefinition.IsFunction())
	var fn rhino.Node
	if classDefinition.IsClass() {
		// NOTE: it would be possible to handle implicit superclass definitions by caching
		// the super class constructor summaries.
		member := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classDefinition)
		fn = member.GetFirstChild()
	} else {
		fn = classDefinition
	}

	argumentsReference := NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn)
	hasVarArgs := argumentsReference || OptimizeConstructors_FunctionHasRest(fn)
	return NewOptimizeConstructors_ClassConstructorSummary(
		hasVarArgs, NodeUtil_GetFunctionParameters(fn).GetChildCount())
}

// Java: OptimizeConstructors.ClassConstructorSummary.isEquivalentConstructorDefinition(Node)
func (s *OptimizeConstructors_ClassConstructorSummary) IsEquivalentConstructorDefinition(constructorMember rhino.Node) bool {
	// `constructorHasRemovableDefinition` has already checked the subclass definition. We
	// know that it has a trivial constructor and that the parameters are passed to the super call
	// in the same order as they are declared in the parameter list, etc.

	// As a result we can simply check the parameters of the constructor and
	// validate they are sufficient for the constructor of the super class

	fn := constructorMember.GetFirstChild()
	paramList := NodeUtil_GetFunctionParameters(fn)

	hasRest := OptimizeConstructors_FunctionHasRest(fn)
	if hasRest {
		return true
	}

	if s.isVarArgs {
		// The count of parameters may matter, and the function may not pass them all on.
		return false
	}

	if s.formalParameterCount != paramList.GetChildCount() {
		return false
	}

	// same count
	return true
}

// Java: OptimizeConstructors.functionHasRest(Node)
func OptimizeConstructors_FunctionHasRest(fn rhino.Node) bool {
	jutil.CheckState(fn.IsFunction())
	params := NodeUtil_GetFunctionParameters(fn)
	lastParam := params.GetLastChild()
	return !lastParam.IsNil() && lastParam.IsRest()
}

// OptimizeConstructors_ConstructorHasRemovableDefinition validates that
// nothing about the constructor definition itself prevents its removal: -
// the body is simply a super call - it passes on exactly what it receives.
// These are necessary but not sufficient for the constructor to be
// removable.
//
// Java: OptimizeConstructors.constructorHasRemovableDefinition(Node)
func OptimizeConstructors_ConstructorHasRemovableDefinition(member rhino.Node) bool {
	fn := member.GetFirstChild()
	superCall := OptimizeConstructors_GetOnlySuperCall(fn)
	if superCall.IsNil() {
		return false
	}

	paramList := NodeUtil_GetFunctionParameters(fn)
	if paramList.GetChildCount() != superCall.GetChildCount()-1 {
		return false
	}

	// TODO(johnlenz): broaden the recognized patterns: default parameters
	param := paramList.GetFirstChild()
	arg := superCall.GetSecondChild()
	for !param.IsNil() {
		if param.IsRest() {
			if !arg.IsSpread() || !param.GetFirstChild().MatchesNameNode(arg.GetFirstChild()) {
				return false
			}
		} else if !param.MatchesNameNode(arg) {
			// not a simple parameter list or matching call args
			return false
		}
		param = param.GetNext()
		arg = arg.GetNext()
	}

	return true
}

// HasRemovableParameterList: true if the parameter list can't cause
// side-effects.
//
// Java: OptimizeConstructors.hasRemovableParameterList(Node)
func (o *OptimizeConstructors) HasRemovableParameterList(member rhino.Node) bool {
	fn := member.GetFirstChild()
	paramList := NodeUtil_GetFunctionParameters(fn)

	for param := paramList.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
		if param.IsName() {
			// a simple name
			continue
		}
		if param.IsRest() && param.GetFirstChild().IsName() {
			// a simple rest expression
			continue
		}

		if param.IsDefaultValue() &&
			param.GetFirstChild().IsName() &&
			!o.astAnalyzer.MayHaveSideEffects(param.GetLastChild()) {
			// a default parameter whose value is determined to be side-effect free
			continue
		}

		// not a simple parameter (destructuring can throw, etc)
		return false
	}

	return true
}

// OptimizeConstructors_GetOnlySuperCall: if the body contains only a call
// to super, return it, otherwise null.
//
// Java: OptimizeConstructors.getOnlySuperCall(Node)
func OptimizeConstructors_GetOnlySuperCall(fn rhino.Node) rhino.Node {
	body := fn.GetLastChild()
	if body.IsBlock() && body.HasOneChild() {
		stmt := body.GetFirstChild()
		if stmt.IsExprResult() {
			call := stmt.GetFirstChild()
			if call.IsCall() && call.GetFirstChild().IsSuper() {
				return call
			}
		}
	}
	return rhino.NilNode
}

// RemoveConstructorMethod removes any candidate constructor if the callers
// are consistent with the definition.
//
// Java: OptimizeConstructors.removeConstructorMethod(Node)
func (o *OptimizeConstructors) RemoveConstructorMethod(member rhino.Node) {
	jutil.CheckState(member.IsMemberFunctionDef())

	o.compiler.ReportFunctionDeleted(member.GetFirstChild())
	o.compiler.ReportChangeToEnclosingScope(member)
	member.Detach()

	// NOTE: As this pass is expected to run by itself, so don't spend the time to update the
	// reference map.
}
