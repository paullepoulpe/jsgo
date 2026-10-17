// Port of com/google/javascript/jscomp/Es6ConvertSuperConstructorCalls.java @ closure-compiler v20260415

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6ConvertSuperConstructorCalls converts `super()` calls.
type Es6ConvertSuperConstructorCalls struct {
	jscompInherits  RuntimeJsLibManager_JsLibField
	jscompConstruct RuntimeJsLibManager_JsLibField

	compiler AbstractCompiler
	// Go-only: language-adapter — the Java ArrayDeque used as a stack: the
	// top (Deque.peek/push/pop) is the last element.
	constructorDataStack []*es6ConvertSuperConstructorCalls_ConstructorData
	astFactory           *AstFactory
	// Note: this GlobalNamespace needs to be the namespace we get after the transpilation passes run.
	// meaning compiler.getTranspilationNamespace() does not work.
	globalNamespace          *GlobalNamespace
	transpilationNamespace   rhino.StaticScope
	uniqueIdSupplier         *UniqueIdSupplier
	es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation
}

// Java: Es6ConvertSuperConstructorCalls.TMP_ERROR
const es6ConvertSuperConstructorCalls_TMP_ERROR = "$jscomp$tmp$error"

// Java: Es6ConvertSuperConstructorCalls.SUPER_THIS
const es6ConvertSuperConstructorCalls_SUPER_THIS = "$jscomp$super$this"

// es6ConvertSuperConstructorCalls_ConstructorData stores superCalls for a
// constructor.
type es6ConvertSuperConstructorCalls_ConstructorData struct {
	constructor *rhino.Node
	superCalls  []*rhino.Node
}

// Java: Es6ConvertSuperConstructorCalls.ConstructorData(Node)
func newEs6ConvertSuperConstructorCalls_ConstructorData(constructor *rhino.Node) *es6ConvertSuperConstructorCalls_ConstructorData {
	c := &es6ConvertSuperConstructorCalls_ConstructorData{}
	c.constructor = constructor
	c.superCalls = []*rhino.Node{}
	return c
}

// Java: Es6ConvertSuperConstructorCalls(AbstractCompiler, Es6SubclassTranspilation)
func NewEs6ConvertSuperConstructorCalls(
	compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) *Es6ConvertSuperConstructorCalls {
	e := &Es6ConvertSuperConstructorCalls{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.transpilationNamespace = compiler.GetTranspilationNamespace()
	e.constructorDataStack = nil
	e.uniqueIdSupplier = compiler.GetUniqueIdSupplier()

	runtimeJsLibManager := compiler.GetRuntimeJsLibManager()
	e.jscompInherits = runtimeJsLibManager.GetJsLibField("$jscomp.inherits")
	e.jscompConstruct = runtimeJsLibManager.GetJsLibField("$jscomp.construct")
	e.es6SubclassTranspilation = es6SubclassTranspilation
	return e
}

// Go-only: language-adapter — Deque.peek(): the top, or nil when empty.
func (e *Es6ConvertSuperConstructorCalls) peekConstructorData() *es6ConvertSuperConstructorCalls_ConstructorData {
	if len(e.constructorDataStack) == 0 {
		return nil
	}
	return e.constructorDataStack[len(e.constructorDataStack)-1]
}

// Java: Es6ConvertSuperConstructorCalls.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6ConvertSuperConstructorCalls) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() && !NodeUtil_GetFeatureSetOfScript(n).ContainsFeature(parser.FeatureSet_Feature_SUPER) {
		// If a script contains Feature.SUPER, only then process super constructor calls for it
		return false
	}
	if n.IsFunction() {
		// TODO(bradfordcsmith): Avoid creating data for non-constructor functions.
		e.constructorDataStack = append(e.constructorDataStack, newEs6ConvertSuperConstructorCalls_ConstructorData(n))
	} else if n.IsSuper() {
		// super(args) or super.prop
		jutil.CheckStateMsg1(n.IsFirstChildOf(parent), "%s", parent)
		if parent.IsGetProp() {
			// TODO(bradfordcsmith): `super.prop` should have been removed before this code executes,
			//     but instead is being left untranspiled when there's no `extends` clause, so
			//     we have to report that problem here.
			t.Report(n, TranspilationUtil_CANNOT_CONVERT_YET, "super access with no extends clause")
			return false
		}
		// must be super(args)
		jutil.CheckStateMsg1(parent.IsCall(), "%s", parent)
		constructorData := jutil.CheckNotNull(e.peekConstructorData())
		constructorData.superCalls = append(constructorData.superCalls, parent)
	}
	return true
}

// Java: Es6ConvertSuperConstructorCalls.visit(NodeTraversal, Node, Node)
func (e *Es6ConvertSuperConstructorCalls) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	constructorData := e.peekConstructorData()
	if constructorData != nil && n == constructorData.constructor {
		e.constructorDataStack[len(e.constructorDataStack)-1] = nil
		e.constructorDataStack = e.constructorDataStack[:len(e.constructorDataStack)-1]
		e.VisitSuper(t, constructorData)
	}
}

// Java: Es6ConvertSuperConstructorCalls.visitSuper(NodeTraversal, ConstructorData)
func (e *Es6ConvertSuperConstructorCalls) VisitSuper(t *NodeTraversal, constructorData *es6ConvertSuperConstructorCalls_ConstructorData) {
	// NOTE: When this pass runs:
	// -   ES6 classes have already been rewritten as ES5 functions.
	// -   All subclasses have $jscomp.inherits() calls connecting them to their parent class.
	// -   All instances of `super` that are not super constructor calls have been rewritten.
	constructor := constructorData.constructor
	superCalls := constructorData.superCalls
	if len(superCalls) == 0 {
		return // nothing to do
	}

	if constructor.IsFromExterns() {
		// This class is defined in an externs file, so it's only a stub, not the actual
		// implementation that should be instantiated.
		// A call to super() shouldn't actually exist for a stub and is problematic to transpile,
		// so just drop it.
		for _, superCall := range superCalls {
			enclosingStatement := NodeUtil_GetEnclosingStatement(superCall)
			enclosingScope := enclosingStatement.GetParent()
			enclosingStatement.Detach()
			e.compiler.ReportChangeToEnclosingScope(enclosingScope)
		}
		return
	}
	// Give a unique name to the $jscomp$super$this variables created when rewriting this super to
	// preserve normalization.
	uniqueSuperThisName :=
		es6ConvertSuperConstructorCalls_SUPER_THIS +
			"$" +
			e.uniqueIdSupplier.GetUniqueId(e.compiler.GetInput(NodeUtil_GetInputId(constructor)))
	// Find the `foo.SuperClass` part of `$jscomp.inherits(foo.SubClass, foo.SuperClass)`
	superClassNameNode := e.GetSuperClassQNameNode(constructor)
	superClassQName := *superClassNameNode.GetQualifiedName()
	thisType := e.GetTypeOfThisForConstructor(constructorData.constructor)

	if e.IsNativeObjectClass(t, superClassQName) {
		// There's no need to call Object as a super constructor, so just replace the call with
		// `this`, which is its correct return value.
		// TODO(bradfordcsmith): Although unlikely, super() could have argument expressions with
		//     side-effects.
		e.ReplaceSuperCallsWithThis(superCalls, thisType)
	} else if e.IsKnownNativeClass(t, superClassQName) {
		// Although we're transpiling down to ES5, it's quite possible that the code will end up
		// running in an environment where native classes are ES6 classes.
		// To correctly extend them with the ES5 classes we're generating here, we must use
		// `$jscomp.construct`, which is our wrapper around `Reflect.construct`.
		e.RewriteSuperCallsToJsCompConstructCalls(
			constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
	} else if e.IsNativeErrorClass(t, superClassQName) {
		// TODO(bradfordcsmith): It might be better to use $jscomp.construct() for these instead
		// of our custom-made, Error-specific workaround.
		e.RewriteNativeErrorSuperCalls(t, superCalls, superClassNameNode, thisType)
	} else if e.IsKnownToReturnOnlyUndefined(superClassQName) {
		e.RewriteSuperCalls(superCalls, superClassNameNode, thisType)
	} else {
		e.RewriteSuperAndPreserveReturnedThis(
			constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
	}
}

// Java: Es6ConvertSuperConstructorCalls.replaceSuperCallsWithThis(List<Node>, AstFactory.Type)
func (e *Es6ConvertSuperConstructorCalls) ReplaceSuperCallsWithThis(superCalls []*rhino.Node, thisType AstFactory_Type) {
	for _, superCall := range superCalls {
		thisNode := e.astFactory.CreateThis(thisType).Srcref(superCall)
		superCall.ReplaceWith(thisNode)
		e.compiler.ReportChangeToEnclosingScope(thisNode)
	}
}

// RewriteSuperCallsToJsCompConstructCalls changes calls to `super` to use
// `$jscomp.construct` instead.
//
//	// note that conversion of the ES6 class to ES5 happens before this pass.
//	// we're just cleaning up the super() calls now
//	var Foo = function(arg1, arg2) {
//	  super(arg1);
//	  this.prop = arg2;
//	}
//	// becomes
//	var Foo = function(arg1, arg2) {
//	  // tmp var and return are necessary, because $jscomp.construct() always creates a new
//	  // object to be used as `this`
//	  var $jscomp$super$this;
//	  $jscomp$super$this = $jscomp.construct(SuperClassName, [arg1], this.constructor);
//	  $jscomp$super$this.prop = arg2
//	  return $jscomp$super$this;
//	}
//
// Java: Es6ConvertSuperConstructorCalls.rewriteSuperCallsToJsCompConstructCalls(Node, List<Node>, Node, AstFactory.Type, String)
func (e *Es6ConvertSuperConstructorCalls) RewriteSuperCallsToJsCompConstructCalls(
	constructor *rhino.Node,
	superCalls []*rhino.Node,
	superClassNameNode *rhino.Node,
	thisType AstFactory_Type,
	uniqueSuperThisName string) {
	constructorBody := jutil.CheckNotNull(constructor.GetChildAtIndex(2))
	firstStatement := constructorBody.GetFirstChild()
	// A constructor body with no call to `super()` is a syntax error for a class that has an
	// extends clause. An error should have been reported and we should never reach this point
	// for an empty constructor body.
	jutil.CheckNotNullMsg(firstStatement, "Empty constructor body")
	firstSuperCall := superCalls[0]

	if constructorBody.HasOneChild() &&
		firstStatement.IsExprResult() &&
		firstStatement.HasOneChild() &&
		firstStatement.GetFirstChild() == firstSuperCall {
		jutil.CheckStateMsg1(len(superCalls) == 1, "%s", constructor)
		// Super call is the entire constructor, so just replace it with.
		// `return $jscomp.construct(SuperClassName, [args], this.constructor);`
		firstStatement.ReplaceWith(
			e.astFactory.CreateReturn(
				e.CreateJSCompConstructorCall(superClassNameNode, firstSuperCall, thisType)))
	} else {
		typeOfThis := e.GetTypeOfThisForConstructor(constructor)
		// `this` -> `$jscomp$super$this` throughout the constructor body,
		// except for super() calls.
		e.UpdateThisToSuperThis(typeOfThis, constructorBody, superCalls, uniqueSuperThisName)
		// Start constructor with `var $jscomp$super$this;`, but place it after any hoisted
		// function declarations
		declaration := e.astFactory.
			CreateSingleVarNameDeclaration(uniqueSuperThisName).
			SrcrefTree(constructorBody)
		insertBeforePoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(constructorBody)
		if insertBeforePoint != nil {
			declaration.InsertBefore(insertBeforePoint)
		} else {
			// functionBody only contains hoisted function declarations
			constructorBody.AddChildToBack(declaration)
		}

		// End constructor with `return $jscomp$super$this;`
		constructorBody.AddChildToBack(
			e.astFactory.
				CreateReturn(e.astFactory.CreateName(uniqueSuperThisName, typeOfThis)).
				SrcrefTree(constructorBody))
		// Replace each super() call with `($jscomp$super$this = $jscomp.construct(...))`
		for _, superCall := range superCalls {
			superCall.ReplaceWith(
				e.astFactory.
					CreateAssign(
						e.astFactory.CreateName(uniqueSuperThisName, typeOfThis).Srcref(superCall),
						e.CreateJSCompConstructorCall(superClassNameNode, superCall, thisType)).
					Srcref(superCall))
		}
	}
	e.compiler.ReportChangeToEnclosingScope(constructorBody)
}

// Java: Es6ConvertSuperConstructorCalls.rewriteNativeErrorSuperCalls(NodeTraversal, List<Node>, Node, AstFactory.Type)
func (e *Es6ConvertSuperConstructorCalls) RewriteNativeErrorSuperCalls(
	t *NodeTraversal, superCalls []*rhino.Node, superClassNameNode *rhino.Node, thisType AstFactory_Type) {
	for _, superCall := range superCalls {
		newSuperCall :=
			e.CreateNewSuperCall(superClassNameNode, superCall, thisType, AstFactory_Type_(superCall))
		e.ReplaceNativeErrorSuperCall(superCall, newSuperCall, t.GetInput())
	}
}

// Java: Es6ConvertSuperConstructorCalls.rewriteSuperCalls(List<Node>, Node, AstFactory.Type)
func (e *Es6ConvertSuperConstructorCalls) RewriteSuperCalls(
	superCalls []*rhino.Node, superClassNameNode *rhino.Node, thisType AstFactory_Type) {
	// super() will not change the value of `this`.
	for _, superCall := range superCalls {
		newSuperCall :=
			e.CreateNewSuperCall(
				superClassNameNode, superCall, thisType, AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID))
		superCallParent := superCall.GetParent()
		if superCallParent.HasOneChild() && NodeUtil_IsStatement(superCallParent) {
			// super() is a statement unto itself
			superCall.ReplaceWith(newSuperCall)
		} else {
			// super() is part of an expression, so it must return `this`.
			superCall.ReplaceWith(
				e.astFactory.
					CreateComma(newSuperCall, e.astFactory.CreateThis(thisType)).
					SrcrefTreeIfMissing(superCall))
		}
		e.compiler.ReportChangeToEnclosingScope(superCallParent)
	}
}

// Java: Es6ConvertSuperConstructorCalls.rewriteSuperAndPreserveReturnedThis(Node, List<Node>, Node, AstFactory.Type, String)
func (e *Es6ConvertSuperConstructorCalls) RewriteSuperAndPreserveReturnedThis(
	constructor *rhino.Node,
	superCalls []*rhino.Node,
	superClassNameNode *rhino.Node,
	thisType AstFactory_Type,
	uniqueSuperThisName string) {
	// Either the superclass constructor returns a value, or we cannot find its definition in
	// the sources, so we don't know if it does.
	//
	// 1. We must use the value it returns, if defined, as the 'this' value in the constructor
	//    we're currently transpiling, and we must also return it from this constructor.
	// 2. It may be an ES6 class defined outside of the sources we can see.
	if e.es6SubclassTranspilation == CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT {
		// To safely extend the ES5 classes we're generating here, we must use
		// `$jscomp.construct`, which is our wrapper around `Reflect.construct`.
		e.RewriteSuperCallsToJsCompConstructCalls(
			constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
		return
	}
	// The code below works as long as the class we're extending is an ES5 class, but will
	// break if we're extending an ES6 class (#2), because it calls the superclass constructor
	// without using `new` or `Reflect.construct()`
	// TODO(b/36789413): we should always use `$jscomp.construct`.
	constructorBody := jutil.CheckNotNull(constructor.GetChildAtIndex(2))
	firstStatement := constructorBody.GetFirstChild()
	firstSuperCall := superCalls[0]

	if constructorBody.HasOneChild() &&
		firstStatement.IsExprResult() &&
		firstStatement.HasOneChild() &&
		firstStatement.GetFirstChild() == firstSuperCall {
		jutil.CheckStateMsg1(len(superCalls) == 1, "%s", constructor)
		// Super call is the entire constructor, so just replace it with.
		// `return <newSuperCall> || this;`
		newReturn :=
			e.astFactory.CreateOr(
				e.CreateNewSuperCall(
					superClassNameNode,
					superCalls[0],
					AstFactory_Type_(superCalls[0]),
					AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
				e.astFactory.CreateThis(thisType))
		firstStatement.ReplaceWith(rhino.IR_ReturnNodeWithExpr(newReturn).SrcrefTreeIfMissing(firstStatement))
	} else {
		typeOfThis := e.GetTypeOfThisForConstructor(constructor)
		// Start constructor with `var $jscomp$super$this;`, but place it after any hoisted
		// function declarations
		declaration :=
			rhino.IR_Var(e.astFactory.CreateName(uniqueSuperThisName, typeOfThis)).
				SrcrefTree(constructorBody)
		insertBeforePoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(constructorBody)
		if insertBeforePoint != nil {
			declaration.InsertBefore(insertBeforePoint)
		} else {
			// functionBody only contains hoisted function declarations
			constructorBody.AddChildToBack(declaration)
		}
		// End constructor with `return $jscomp$super$this;`
		constructorBody.AddChildToBack(
			rhino.IR_ReturnNodeWithExpr(e.astFactory.CreateName(uniqueSuperThisName, typeOfThis)).
				SrcrefTree(constructorBody))

		// Replace each `this` -> `$jscomp$super$this` throughout the constructor body,
		// except for super() calls.
		e.UpdateThisToSuperThis(typeOfThis, constructorBody, superCalls, uniqueSuperThisName)

		// Replace each super() call with `($jscomp$super$this = <newSuperCall> || this)`
		for _, superCall := range superCalls {
			newSuperCall :=
				e.CreateNewSuperCall(
					superClassNameNode, superCall, typeOfThis, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN))
			superCall.ReplaceWith(
				e.astFactory.
					CreateAssign(
						e.astFactory.CreateName(uniqueSuperThisName, typeOfThis),
						e.astFactory.CreateOr(newSuperCall, e.astFactory.CreateThis(typeOfThis))).
					SrcrefTreeIfMissing(superCall))
		}
	}
	e.compiler.ReportChangeToEnclosingScope(constructorBody)
}

// Java: Es6ConvertSuperConstructorCalls.isKnownToReturnOnlyUndefined(String)
func (e *Es6ConvertSuperConstructorCalls) IsKnownToReturnOnlyUndefined(functionQName string) bool {
	if e.globalNamespace == nil {
		return false
	}
	// Go-only: language-adapter — GlobalNamespace.getSlot returns a
	// rhino.StaticSlot (nil when absent); the Java return type is Name.
	slot := e.globalNamespace.GetSlot(functionQName)
	if slot == nil {
		return false
	}
	globalName := slot.(*GlobalNamespace_Name)

	// Go-only: language-adapter — Name.getDeclaration returns a
	// rhino.StaticRef (nil when absent); the Java return type is Ref.
	var declarationRef *GlobalNamespace_Ref
	if d := globalName.GetDeclaration(); d != nil {
		declarationRef = d.(*GlobalNamespace_Ref)
	}
	if declarationRef == nil {
		for _, ref := range globalName.GetRefs() {
			if ref.IsSet() {
				declarationRef = ref
			}
		}
	}
	if declarationRef == nil {
		return false
	}

	declaredVarOrProp := declarationRef.GetNode()
	if declaredVarOrProp.IsFromExterns() {
		return false
	}

	declaration := declaredVarOrProp.GetParent()
	var declaredValue *rhino.Node
	if declaration.IsFunction() {
		declaredValue = declaration
	} else if NodeUtil_IsNameDeclaration(declaration) && declaredVarOrProp.IsName() {
		if declaredVarOrProp.HasChildren() {
			declaredValue = jutil.CheckNotNull(declaredVarOrProp.GetFirstChild())
		} else {
			return false // Declaration without an assigned value.
		}
	} else if declaration.IsAssign() && declaration.GetFirstChild() == declaredVarOrProp {
		declaredValue = jutil.CheckNotNull(declaration.GetSecondChild())
	} else if declaration.IsObjectLit() && declaredVarOrProp.HasOneChild() {
		declaredValue = jutil.CheckNotNull(declaredVarOrProp.GetFirstChild())
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected declaration format:\n" + declaration.ToStringTree()))
	}
	return e.IsNodeKnownToOnlyReturnUndefined(declaredValue)
}

// Java: Es6ConvertSuperConstructorCalls.isNodeKnownToOnlyReturnUndefined(Node)
func (e *Es6ConvertSuperConstructorCalls) IsNodeKnownToOnlyReturnUndefined(node *rhino.Node) bool {
	if node.IsFunction() {
		functionBody := jutil.CheckNotNull(node.GetChildAtIndex(2))
		return !newEs6ConvertSuperConstructorCalls_UndefinedReturnValueCheck(e).MayReturnDefinedValue(functionBody)
	} else if node.IsQualifiedName() {
		return e.IsKnownToReturnOnlyUndefined(*node.GetQualifiedName())
	} else if node.IsHook() {
		// cond ? left : right;
		left := node.GetSecondChild()
		right := node.GetLastChild()
		return e.IsNodeKnownToOnlyReturnUndefined(left) && e.IsNodeKnownToOnlyReturnUndefined(right)
	} else {
		// TODO(bradfordcsmith): What cases are these? Can we do better?
		return false
	}
}

// Java: private class Es6ConvertSuperConstructorCalls.UndefinedReturnValueCheck
type es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck struct {
	outer               *Es6ConvertSuperConstructorCalls
	foundNonEmptyReturn bool
}

// Java: Es6ConvertSuperConstructorCalls.UndefinedReturnValueCheck()
func newEs6ConvertSuperConstructorCalls_UndefinedReturnValueCheck(outer *Es6ConvertSuperConstructorCalls) *es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck {
	return &es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck{outer: outer}
}

// Java: Es6ConvertSuperConstructorCalls.UndefinedReturnValueCheck.mayReturnDefinedValue(Node)
func (u *es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck) MayReturnDefinedValue(functionBody *rhino.Node) bool {
	u.foundNonEmptyReturn = false
	checkForDefinedReturnValue := &es6ConvertSuperConstructorCalls_checkForDefinedReturnValue{outer: u}
	NodeTraversal_Traverse(u.outer.compiler, functionBody, checkForDefinedReturnValue)
	return u.foundNonEmptyReturn
}

// es6ConvertSuperConstructorCalls_checkForDefinedReturnValue is the
// anonymous NodeTraversal.AbstractShallowCallback of
// UndefinedReturnValueCheck.mayReturnDefinedValue.
type es6ConvertSuperConstructorCalls_checkForDefinedReturnValue struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer *es6ConvertSuperConstructorCalls_UndefinedReturnValueCheck
}

// Java: Es6ConvertSuperConstructorCalls.UndefinedReturnValueCheck.mayReturnDefinedValue(Node).new AbstractShallowCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *es6ConvertSuperConstructorCalls_checkForDefinedReturnValue) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !c.outer.foundNonEmptyReturn {
		if n.IsReturn() &&
			n.HasChildren() &&
			!n.GetFirstChild().MatchesName("undefined") {
			c.outer.foundNonEmptyReturn = true
		}
	}
}

// CreateNewSuperCall returns a transpiled version of the super constructor
// call.
//
// The children of the passed in `superCall` are all removed from it by this
// method, but the existing call itself is not replaced in the AST yet. The
// returned node is not yet attached to the AST.
//
// Java: Es6ConvertSuperConstructorCalls.createNewSuperCall(Node, Node, AstFactory.Type, AstFactory.Type)
func (e *Es6ConvertSuperConstructorCalls) CreateNewSuperCall(
	superClassQNameNode *rhino.Node,
	superCall *rhino.Node,
	thisType AstFactory_Type,
	resultType AstFactory_Type) *rhino.Node {
	jutil.CheckArgumentMsg1(superClassQNameNode.IsQualifiedName(), "%s", superClassQNameNode)
	jutil.CheckArgumentMsg1(superCall.IsCall(), "%s", superCall)

	callee := superCall.RemoveFirstChild()
	jutil.CheckStateMsg1(callee.IsSuper(), "%s", callee)

	args := []*rhino.Node{}
	hasSpreadArg := false
	for superCall.HasChildren() {
		arg := superCall.RemoveFirstChild()
		hasSpreadArg = hasSpreadArg || arg.IsSpread()
		args = append(args, arg)
	}

	// Node to which args should be appended
	if hasSpreadArg {
		// We want to convert
		//
		// super(x, ...params, y)
		// to
		// Foo.apply(this, [x, ...params, y])
		//
		// because, after transpilation of spread this becomes
		//
		// Foo.apply(this, [x, $jscomp.arrayFromIterable(params), y])
		//
		// If we used `call`, we'd get this nonsense instead
		//
		// Foo.call.apply(Foo, [this, x, $jscomp.arrayFromIterable(params), y])
		superClassDotApply := e.astFactory.
			CreateGetPropWithUnknownType(superClassQNameNode.CloneTree(), "apply").
			SrcrefTree(callee)
		// Create `SuperClass.call(this)`
		newSuperCall := e.astFactory.CreateCall(superClassDotApply, resultType).Srcref(superCall)
		newSuperCall.AddChildToBack(e.astFactory.CreateThis(thisType).Srcref(callee))
		newSuperCall.PutBooleanProp(rhino.Node_FREE_CALL, false) // callee is now a getprop
		// It's very common to just have `super(...arguments)`, because we generate constructors
		// containing that for extending classes that don't have an explicit constructor.
		// For that case it's more efficient to just convert `super(...arguments)` to
		// `SuperClass.apply(this, arguments)` here rather than relying on later optimizations to
		// convert `[...arguments]` to `arguments`.
		if Es6ConvertSuperConstructorCalls_IsSingleSpreadOfArguments(args) {
			newSuperCall.AddChildToBack(args[0].GetOnlyChild().Detach())
		} else {
			newSuperCall.AddChildToBack(e.astFactory.CreateArraylitOfIterable(slices.Values(args)).Srcref(superCall))
		}
		return newSuperCall
	} else {
		// We want to convert
		//
		// super(arg1, arg2)
		// to
		// Foo.call(this, arg1, arg2)
		//
		// Using `call` is shorter than using `apply`.
		superClassDotCall := e.astFactory.
			CreateGetProp(
				superClassQNameNode.CloneTree(), "call", AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)).
			SrcrefTree(callee)
		newSuperCall := e.astFactory.CreateCall(superClassDotCall, resultType).Srcref(superCall)
		newSuperCall.AddChildToBack(e.astFactory.CreateThis(thisType).Srcref(callee))
		newSuperCall.PutBooleanProp(rhino.Node_FREE_CALL, false) // callee is now a getprop
		for _, arg := range args {
			newSuperCall.AddChildToBack(arg)
		}
		return newSuperCall
	}
}

// CreateJSCompConstructorCall returns a transpiled version of the super
// constructor call using `$jscomp.constructor`.
//
// The children of the passed in `superCall` are all removed from it by this
// method, but the existing call itself is not replaced in the AST yet. The
// returned node is not yet attached to the AST.
//
// Java: Es6ConvertSuperConstructorCalls.createJSCompConstructorCall(Node, Node, AstFactory.Type)
func (e *Es6ConvertSuperConstructorCalls) CreateJSCompConstructorCall(
	superClassQNameNode *rhino.Node, superCall *rhino.Node, thisType AstFactory_Type) *rhino.Node {
	jutil.CheckArgumentMsg1(superClassQNameNode.IsQualifiedName(), "%s", superClassQNameNode)
	jutil.CheckArgumentMsg1(superCall.IsCall(), "%s", superCall)

	callee := jutil.CheckNotNullMsg1(superCall.RemoveFirstChild(), "%s", superCall)
	jutil.CheckStateMsg1(callee.IsSuper(), "%s", callee)

	// `$jscomp.construct`
	jscompDotConstruct :=
		e.astFactory.CreateQNameOfField(e.transpilationNamespace, e.jscompConstruct).SrcrefTree(callee)

	superClassQName := superClassQNameNode.CloneTree()

	// extract the arguments from the super() call and create the arguments list to pass to
	// $jscomp.construct()
	superCallArgList := []*rhino.Node{}
	for superCall.HasChildren() {
		superCallArgList = append(superCallArgList, superCall.RemoveFirstChild())
	}
	// It's very common to just have `super(...arguments)`, because we generate constructors
	// containing that for extending classes that don't have an explicit constructor.
	// For that case it's more efficient to just convert `super(...arguments)` to
	// `$jscomp.construct(SuperClass, arguments, this.constructor)` here rather than relying on
	// later optimizations to
	// convert `[...arguments]` to `arguments`.
	var superArgs *rhino.Node
	if Es6ConvertSuperConstructorCalls_IsSingleSpreadOfArguments(superCallArgList) {
		// pull out `arguments` from `...arguments`
		superArgs = superCallArgList[0].GetOnlyChild().Detach()
	} else {
		superArgs = e.astFactory.CreateArraylitOfIterable(slices.Values(superCallArgList)).Srcref(superCall)
	}

	// `this.constructor`
	thisDotConstructor := e.astFactory.
		CreateGetProp(
			e.astFactory.CreateThis(thisType), "constructor", AstFactory_Type_(superClassQNameNode)).
		SrcrefTree(superCall)

	// `super(arg1, arg2)`
	// becomes
	// `$jscomp.construct(SuperClassName, [arg1, arg2], this.constructor)`
	return e.astFactory.
		CreateCall(
			jscompDotConstruct, AstFactory_Type_(superCall), superClassQName, superArgs, thisDotConstructor).
		Srcref(superCall)
}

// Java: Es6ConvertSuperConstructorCalls.isSingleSpreadOfArguments(List<Node>)
func Es6ConvertSuperConstructorCalls_IsSingleSpreadOfArguments(nodeList []*rhino.Node) bool {
	return len(nodeList) == 1 && Es6ConvertSuperConstructorCalls_IsSpreadOfArguments(nodeList[0])
}

// Java: Es6ConvertSuperConstructorCalls.isSpreadOfArguments(Node)
func Es6ConvertSuperConstructorCalls_IsSpreadOfArguments(node *rhino.Node) bool {
	return node.IsSpread() && node.GetOnlyChild().MatchesName("arguments")
}

// Java: Es6ConvertSuperConstructorCalls.replaceNativeErrorSuperCall(Node, Node, CompilerInput)
func (e *Es6ConvertSuperConstructorCalls) ReplaceNativeErrorSuperCall(
	superCall *rhino.Node, newSuperCall *rhino.Node, compilerInput *CompilerInput) {
	// The native error class constructors always return a new object instead of initializing
	// `this`, so a workaround is needed.
	superStatement := NodeUtil_GetEnclosingStatement(superCall)
	body := superStatement.GetParent()
	jutil.CheckStateMsg1(body.IsBlock(), "%s", body)

	thisType := AstFactory_Type_(newSuperCall)

	// Give a unique name to the $jscomp$tmp$error; variables created when rewriting this super to
	// preserve normalization.
	tmpErrorName := es6ConvertSuperConstructorCalls_TMP_ERROR + "$" + e.uniqueIdSupplier.GetUniqueId(compilerInput)

	// var $jscomp$tmp$error;
	getError :=
		rhino.IR_Var(e.astFactory.CreateName(tmpErrorName, thisType)).SrcrefTreeIfMissing(superCall)
	getError.InsertBefore(superStatement)

	// Create an expression to initialize `this` from temporary Error object at the point
	// where super.apply() was called.
	// $jscomp$tmp$error = Error.call(this, ...),
	getTmpError :=
		e.astFactory.CreateAssign(e.astFactory.CreateName(tmpErrorName, thisType), newSuperCall)
	// this.message = $jscomp$tmp$error.message,
	copyMessage :=
		e.astFactory.CreateAssign(
			e.astFactory.CreateGetProp(
				e.astFactory.CreateThis(thisType), "message", AstFactory_Type_OfColor(colors.StandardColors_STRING)),
			e.astFactory.CreateGetProp(
				e.astFactory.CreateName(tmpErrorName, thisType),
				"message",
				AstFactory_Type_OfColor(colors.StandardColors_STRING)))

	// Old versions of IE Don't set stack until the object is thrown, and won't set it then
	// if it already exists on the object.
	// ('stack' in $jscomp$tmp$error) && (this.stack = $jscomp$tmp$error.stack)
	setStack :=
		e.astFactory.CreateAnd(
			e.astFactory.CreateIn(
				e.astFactory.CreateString("stack"), e.astFactory.CreateName(tmpErrorName, thisType)),
			e.astFactory.CreateAssign(
				e.astFactory.CreateGetProp(
					e.astFactory.CreateThis(thisType), "stack", AstFactory_Type_OfColor(colors.StandardColors_STRING)),
				e.astFactory.CreateGetProp(
					e.astFactory.CreateName(tmpErrorName, thisType),
					"stack",
					AstFactory_Type_OfColor(colors.StandardColors_STRING))))
	superErrorExpr := e.astFactory.
		CreateCommas(getTmpError, copyMessage, setStack, e.astFactory.CreateThis(thisType)).
		SrcrefTreeIfMissing(superCall)
	superCall.ReplaceWith(superErrorExpr)
	e.compiler.ReportChangeToEnclosingScope(superErrorExpr)
}

// Java: Es6ConvertSuperConstructorCalls.isNativeObjectClass(NodeTraversal, String)
func (e *Es6ConvertSuperConstructorCalls) IsNativeObjectClass(t *NodeTraversal, className string) bool {
	return className == "Object" && !e.IsDefinedInSources(t, className)
}

// Java: Es6ConvertSuperConstructorCalls.isNativeErrorClass(NodeTraversal, String)
func (e *Es6ConvertSuperConstructorCalls) IsNativeErrorClass(t *NodeTraversal, superClassName string) bool {
	switch superClassName {
	case "AggregateError",
		"Error",
		"EvalError",
		"RangeError",
		"ReferenceError",
		"SyntaxError",
		"TypeError",
		"URIError":
		// All Error classes listed in the ECMAScript spec as of 2016
		return !e.IsDefinedInSources(t, superClassName)
	default:
		return false
	}
}

// IsKnownNativeClass reports whether `className` is the name of a known
// native JS class for which we haven't seen a definition in the source code
// we're compiling. (Note that our own polyfill definitions don't count as a
// definition being present in the source code.)
//
// Java: Es6ConvertSuperConstructorCalls.isKnownNativeClass(NodeTraversal, String)
func (e *Es6ConvertSuperConstructorCalls) IsKnownNativeClass(t *NodeTraversal, className string) bool {
	// This list originally taken from the list of built-in objects at
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference
	// as of 2016-10-22.
	// - Intl.* classes were left out, because it doesn't seem worth the extra effort
	//   of handling the qualified name.
	// - Deprecated and experimental classes were left out.
	switch className {
	case "Array",
		"ArrayBuffer",
		"Boolean",
		"DataView",
		"Date",
		"Float32Array",
		"Function",
		"Generator",
		"GeneratorFunction",
		"Int16Array",
		"Int32Array",
		"Int8Array",
		"InternalError",
		"Map",
		"Number",
		"Object",
		"Promise",
		"Proxy",
		"RegExp",
		"Set",
		"String",
		"Symbol",
		"TypedArray",
		"Uint16Array",
		"Uint32Array",
		"Uint8Array",
		"Uint8ClampedArray",
		"WeakMap",
		"WeakSet":
		return !e.IsDefinedInSources(t, className)
	default:
		return false
	}
}

// IsDefinedInSources reports whether a variable with the given name is
// defined in the source code being compiled.
//
// Please note that the call to t.getScope() is expensive, so we should
// avoid calling this method when possible.
//
// Java: Es6ConvertSuperConstructorCalls.isDefinedInSources(NodeTraversal, String)
func (e *Es6ConvertSuperConstructorCalls) IsDefinedInSources(t *NodeTraversal, varName string) bool {
	objectVar := t.GetScope().GetVar(varName)
	return objectVar != nil && !objectVar.IsExtern()
}

// Java: Es6ConvertSuperConstructorCalls.updateThisToSuperThis(AstFactory.Type, Node, List<Node>, String)
func (e *Es6ConvertSuperConstructorCalls) UpdateThisToSuperThis(
	typeOfThis AstFactory_Type,
	constructorBody *rhino.Node,
	superCalls []*rhino.Node,
	uniqueSuperThisName string) {
	replaceThisWithSuperThis := &es6ConvertSuperConstructorCalls_replaceThisWithSuperThis{
		outer:               e,
		typeOfThis:          typeOfThis,
		superCalls:          superCalls,
		uniqueSuperThisName: uniqueSuperThisName,
	}
	NodeTraversal_Traverse(e.compiler, constructorBody, replaceThisWithSuperThis)
}

// es6ConvertSuperConstructorCalls_replaceThisWithSuperThis is the anonymous
// NodeTraversal.Callback of updateThisToSuperThis.
type es6ConvertSuperConstructorCalls_replaceThisWithSuperThis struct {
	outer               *Es6ConvertSuperConstructorCalls
	typeOfThis          AstFactory_Type
	superCalls          []*rhino.Node
	uniqueSuperThisName string
}

// Java: Es6ConvertSuperConstructorCalls.updateThisToSuperThis(...).new Callback() {...}.shouldTraverse(NodeTraversal, Node, Node)
func (c *es6ConvertSuperConstructorCalls_replaceThisWithSuperThis) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if slices.Contains(c.superCalls, n) {
		return false // Leave `this` intact on super calls.
	} else if n.IsFunction() && !n.IsArrowFunction() {
		// Don't replace `this` in non-arrow function definitions.
		return false
	} else {
		return true
	}
}

// Java: Es6ConvertSuperConstructorCalls.updateThisToSuperThis(...).new Callback() {...}.visit(NodeTraversal, Node, Node)
func (c *es6ConvertSuperConstructorCalls_replaceThisWithSuperThis) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsThis() {
		superThis := c.outer.astFactory.CreateName(c.uniqueSuperThisName, AstFactory_Type_(n)).Srcref(n)
		n.ReplaceWith(superThis)
	} else if n.IsReturn() && !n.HasChildren() {
		// An empty return needs to be changed to return $jscomp$super$this
		n.AddChildToFront(c.outer.astFactory.CreateName(c.uniqueSuperThisName, c.typeOfThis).Srcref(n))
	}
}

// Java: Es6ConvertSuperConstructorCalls.getTypeOfThisForConstructor(Node)
func (e *Es6ConvertSuperConstructorCalls) GetTypeOfThisForConstructor(constructor *rhino.Node) AstFactory_Type {
	jutil.CheckArgumentMsg1(constructor.IsFunction(), "%s", constructor)
	constructorType := constructor.GetColor()
	if constructorType != nil && !constructorType.GetInstanceColors().IsEmpty() {
		return AstFactory_Type_OfColor(colors.Color_CreateUnion(constructorType.GetInstanceColors()))
	}
	return AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)
}

// Java: Es6ConvertSuperConstructorCalls.getSuperClassQNameNode(Node)
func (e *Es6ConvertSuperConstructorCalls) GetSuperClassQNameNode(constructor *rhino.Node) *rhino.Node {
	className := NodeUtil_GetNameNode(constructor).GetQualifiedName()
	constructorStatement := jutil.CheckNotNull(NodeUtil_GetEnclosingStatement(constructor))

	var superClassNameNode *rhino.Node
	for statement := constructorStatement.GetNext(); statement != nil; statement = statement.GetNext() {
		superClassNameNode = e.GetSuperClassNameNodeIfIsInheritsStatement(statement, className)
		if superClassNameNode != nil {
			break
		}
	}

	return jutil.CheckNotNullMsg(superClassNameNode, "$jscomp.inherits() call not found.")
}

// Java: Es6ConvertSuperConstructorCalls.getSuperClassNameNodeIfIsInheritsStatement(Node, String)
func (e *Es6ConvertSuperConstructorCalls) GetSuperClassNameNodeIfIsInheritsStatement(
	statement *rhino.Node, className *string) *rhino.Node {
	// $jscomp.inherits(ChildClass, SuperClass);
	if !statement.IsExprResult() {
		return nil
	}
	callNode := statement.GetFirstChild()
	if !callNode.IsCall() {
		return nil
	}
	jscompDotInherits := callNode.GetFirstChild()
	if !e.jscompInherits.Matches(jscompDotInherits) {
		return nil
	}
	classNameNode := jutil.CheckNotNull(jscompDotInherits.GetNext())
	if classNameNode.MatchesQualifiedName(*className) {
		return jutil.CheckNotNull(classNameNode.GetNext())
	} else {
		return nil
	}
}

// Java: Es6ConvertSuperConstructorCalls.setGlobalNamespace(GlobalNamespace)
func (e *Es6ConvertSuperConstructorCalls) SetGlobalNamespace(globalNamespace *GlobalNamespace) {
	e.globalNamespace = globalNamespace
}
