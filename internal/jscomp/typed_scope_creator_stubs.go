// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AssertDefinitionNode asserts that it's OK to define this node's name. The
// node should have a source name and be of the specified type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.assertDefinitionNode(Node, Token)
func (b *TypedScopeCreator_scopeBuilder) AssertDefinitionNode(n rhino.Node, type_ rhino.Token) {
	jutil.CheckState(b.hasSourceName)
	jutil.CheckStateMsg1(n.GetToken() == type_, "%s", n)
}

// DefineCatch defines a catch parameter.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineCatch(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineCatch(n rhino.Node) {
	b.AssertDefinitionNode(n, rhino.Token_CATCH)
	// Though almost certainly a terrible idea, it is possible to do destructuring in
	// the catch declaration.
	// e.g. `} catch ({message, errno}) {`
	NodeUtil_VisitLhsNodesInNode(
		n,
		func(catchName rhino.Node) {
			type_ := b.GetDeclaredType(catchName.GetJSDocInfo(), catchName, rhino.NilNode, nil)
			b.defineSlot(typedScopeSlotSpec{
				declNode:                catchName,
				varName:                 catchName.GetString(),
				scope:                   b.currentScope,
				type_:                   type_,
				allowLaterTypeInference: type_ == nil,
			})
		})
}

// DefineAssignAsIfVarDeclaration defines an assignment to a name as if it
// were an actual `var` declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineAssignAsIfVarDeclaration(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineAssignAsIfVarDeclaration(assignment rhino.Node) {
	assign := ast.MustAssign(assignment)
	info := assignment.GetJSDocInfo()
	name := assign.Target().Node()
	jutil.CheckArgumentMsg1(name.IsName(), "%s", name)
	b.DefineName(name, assign.Value().Node(), b.currentScope.GetClosestHoistScope(), info)
}

// DefineVars defines a variable declared with `var`, `let`, or `const`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineVars(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineVars(n rhino.Node) {
	jutil.CheckState(b.hasSourceName)
	decl := ast.MustNameDeclaration(n)
	info := n.GetJSDocInfo()
	// `var` declarations are hoisted, but `let` and `const` are not.
	scope := b.currentScope
	if n.IsVar() {
		scope = b.currentHoistScope
	}

	if decl.NumDeclarators() > 1 && info != nil {
		b.outer.Report(JSError_MakeWithNode(n, TypeCheck_MULTIPLE_VAR_DEF))
	}

	for child := range decl.Declarators() {
		b.DefineVarChild(info, child, scope)
	}
	if only, ok := decl.OnlyDeclarator(); ok && b.IsValidTypedefDeclaration(only.Node(), info) {
		b.DeclareTypedefType(only.Node(), info)
	}
}

// DefineVarChild defines a variable declared with `var`, `let`, or `const`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineVarChild(JSDocInfo, Node, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) DefineVarChild(declarationInfo *rhino.JSDocInfo, child ast.Declarator, scope *TypedScope) {
	if name, ok := child.AsName(); ok {
		nameNode := name.Node()
		if declarationInfo == nil {
			declarationInfo = nameNode.GetJSDocInfo()
			// TODO(bradfordcsmith): Report an error if both the declaration node and the name itself
			//     have JSDoc.
		}
		var initNode rhino.Node
		if init, ok := name.Init(); ok {
			initNode = init.Node()
		}
		b.DefineName(nameNode, initNode, scope, declarationInfo)
	} else {
		lhs := child.MustDestructuringLhs()
		pattern := lhs.Pattern().Node()
		var value rhino.Node
		if init, ok := lhs.Init(); ok {
			value = init.Node()
		}

		if ModuleImportResolver_IsGoogModuleDependencyCall(value) {
			// Define destructuring names here, since goog.require destructuring patterns can only
			// have one level and require some special handling.
			defaultImport := b.outer.moduleImportResolver.GetClosureNamespaceTypeFromCall(value)
			if objPat, ok := ast.AsObjectPattern(pattern); ok {
				for prop := range objPat.Props() {
					if sk, ok := prop.AsStringKey(); ok {
						b.DefineModuleImport(sk.Value().Node(), defaultImport, jutil.StringPtr(sk.Key()), scope)
					}
				}
			}
			return
		}

		b.DefineDestructuringPatternInVarDeclaration(
			pattern,
			scope,
			func() typedScopeRValueInfo {
				// Note that value will be null if we are in an enhanced for loop
				//   for (const {x, y} of data) {
				if !value.IsNil() {
					return typedScopeRValueInfo{
						type_:         b.GetDeclaredRValueType( /* lValue= */ rhino.NilNode, value),
						qualifiedName: value.GetQualifiedNameObject(),
					}
				}
				return typedScopeRValueInfo{type_: b.outer.unknownType}
			})
	}
}

// InferTypeForDestructuredTarget returns information about the qualified
// name and type of the target, if it exists.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.inferTypeForDestructuredTarget(DestructuredTarget, Supplier)
func (b *TypedScopeCreator_scopeBuilder) InferTypeForDestructuredTarget(
	target DestructuredTarget, patternTypeSupplier func() typedScopeRValueInfo) typedScopeRValueInfo {
	// Currently we only do type inference for string key nodes in object patterns here, to
	// handle aliasing types. e.g
	//   const {Foo} = ns;
	// TypeInference takes care of the rest.
	// Note that although DestructuredTarget includes logic for inferring types, we don't use
	// it here because we only do some very limited type inference during TypedScopeCreation,
	// and only return a non-null type here if we are accessing a declared property on a known
	// type.
	if !target.HasStringKey() || target.HasDefaultValue() {
		return typedScopeRValueInfo{}
	}
	rvalue := patternTypeSupplier()
	patternType := rvalue.type_
	propertyName := target.GetStringKey().GetString()
	var qname rhino.QualifiedName
	if rvalue.qualifiedName != nil {
		qname = rvalue.qualifiedName.Getprop(propertyName)
	}
	if patternType == nil || patternType.IsUnknownType() {
		return typedScopeRValueInfo{qualifiedName: qname}
	}
	if patternType.HasProperty(propertyName) {
		type_ := patternType.FindPropertyType(propertyName)
		return typedScopeRValueInfo{type_: type_, qualifiedName: qname}
	}
	return typedScopeRValueInfo{qualifiedName: qname}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.defineDestructuringPatternInVarDeclaration(Node, TypedScope, Supplier)
func (b *TypedScopeCreator_scopeBuilder) DefineDestructuringPatternInVarDeclaration(
	pattern rhino.Node, scope *TypedScope, patternTypeSupplier func() typedScopeRValueInfo) {
	ForEachDestructuredTarget(b.outer.typeRegistry, nil, pattern, func(target DestructuredTarget) {
		typeSupplier := func() typedScopeRValueInfo {
			return b.InferTypeForDestructuredTarget(target, patternTypeSupplier)
		}

		if target.GetNode().IsDestructuringPattern() {
			b.DefineDestructuringPatternInVarDeclaration(target.GetNode(), scope, typeSupplier)
		} else {
			name := target.GetNode()
			jutil.CheckStateMsg1(name.IsName(), "This method is only for declaring variables: %s", name)

			// variable's type
			type_ := b.GetDeclaredType(name.GetJSDocInfo(), name /* rValue= */, rhino.NilNode, typeSupplier)
			if type_ == nil {
				// The variable's type will be inferred.
				if name.IsFromExterns() {
					type_ = b.outer.unknownType
				} else {
					type_ = nil
				}
			}
			b.defineSlot(typedScopeSlotSpec{
				declNode:                name,
				varName:                 name.GetString(),
				scope:                   scope,
				type_:                   type_,
				allowLaterTypeInference: type_ == nil,
			})
		}
	})
}

// DefineClassLiteral defines a class literal. Handles any of the following
// cases:
//
//   - Class declarations: `class Foo { ... }`
//   - Class assignments: `foo.Bar = class { ... }`
//   - Bleeding names: `foo.Bar = class Baz { ... }`
//   - Properties: `{foo: class { ... }}`
//   - Callbacks: `foo(class { ... })`
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineClassLiteral(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineClassLiteral(n rhino.Node) {
	b.AssertDefinitionNode(n, rhino.Token_CLASS)

	// Determine the name and JSDocInfo and l-value for the class.
	// Any of these may be null.
	lValue := NodeUtil_GetBestLValue(n)
	info := NodeUtil_GetBestJSDocInfo(n)
	className := NodeUtil_GetBestLValueName(lValue)
	classTypeIdentifier := b.GetBestTypeName(lValue, className)

	// Create the type and assign it on the CLASS node.
	classType := b.CreateClassTypeFromNodes(n, classTypeIdentifier, className, info, lValue)
	b.outer.SetDeferredType(n, classType)

	// Declare this symbol in the current scope iff it's a class
	// declaration. Otherwise, the declaration will happen in other
	// code paths.
	if NodeUtil_IsClassDeclaration(n) {
		jutil.CheckNotNull(className)
		b.defineSlot(typedScopeSlotSpec{
			declNode:                ast.MustClass(n).NameNode(),
			varName:                 *className,
			scope:                   b.currentScope,
			type_:                   classType,
			allowLaterTypeInference: classType == nil,
		})
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.getBestTypeName(Node, String)
func (b *TypedScopeCreator_scopeBuilder) GetBestTypeName(lvalue rhino.Node, syntacticLvalueName *string) *string {
	if syntacticLvalueName == nil {
		return nil
	}
	if replaced, changed := b.replaceGoogModuleExportsName(lvalue, *syntacticLvalueName); changed {
		return &replaced
	}
	return syntacticLvalueName
}

func (b *TypedScopeCreator_scopeBuilder) replaceGoogModuleExportsName(lvalue rhino.Node, syntacticLvalueName string) (string, bool) {
	if b.IsGoogModuleExports(lvalue) {
		return jutil.StringReplace(syntacticLvalueName, "exports", *b.GetModule().ClosureNamespace()), true
	}
	return syntacticLvalueName, false
}

// DefineFunctionLiteral defines a function literal.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineFunctionLiteral(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineFunctionLiteral(n rhino.Node) {
	b.AssertDefinitionNode(n, rhino.Token_FUNCTION)

	// Determine the name and JSDocInfo and l-value for the function.
	// Any of these may be null.
	lValue := NodeUtil_GetBestLValue(n)
	info := NodeUtil_GetBestJSDocInfo(n)
	var functionName *string
	if !lValue.IsNil() && lValue.IsGetProp() {
		if origName := n.GetOriginalName(); origName != nil && lValue.MatchesQualifiedName(*origName) {
			functionName = origName
		}
	}
	if functionName == nil {
		functionName = NodeUtil_GetBestLValueName(lValue)
	}
	functionType := b.CreateFunctionTypeFromNodes(n, functionName, info, lValue)

	// Assigning the function type to the function node
	b.outer.SetDeferredType(n, functionType)

	// Declare this symbol in the current scope iff it's a function
	// declaration. Otherwise, the declaration will happen in other
	// code paths.
	if NodeUtil_IsFunctionDeclaration(n) {
		b.defineSlot(typedScopeSlotSpec{
			declNode:                ast.MustFunction(n).NameNode().Node(),
			varName:                 *functionName,
			scope:                   b.currentScope,
			type_:                   functionType,
			allowLaterTypeInference: functionType == nil,
		})
	}
}

// DefineName defines a variable based on the NAME node passed.
//
// name is the NAME node; value optionally the value assigned to the name
// node; info the JSDocInfo information relating to this name node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineName(Node, Node, TypedScope, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) DefineName(name rhino.Node, value rhino.Node, scope *TypedScope, info *rhino.JSDocInfo) {
	if ModuleImportResolver_IsGoogModuleDependencyCall(value) {
		b.DefineModuleImport(
			name, b.outer.moduleImportResolver.GetClosureNamespaceTypeFromCall(value), nil, scope)
		return
	}
	type_ := b.GetDeclaredType(info, name, value /* declaredRValueTypeSupplier= */, nil)
	if type_ == nil {
		// The variable's type will be inferred.
		if name.IsFromExterns() {
			type_ = b.outer.unknownType
		} else {
			type_ = nil
		}
	}
	b.defineSlot(typedScopeSlotSpec{
		declNode:                name,
		varName:                 name.GetString(),
		scope:                   scope,
		type_:                   type_,
		allowLaterTypeInference: type_ == nil,
	})
}

// DefineModuleImport defines the name of an import.
//
// localNameNode is the name node of the LHS of the import being defined;
// importedModuleObject the root node of the scope in which the type being
// imported was defined, along with the local name of the overall module
// object inside the scope where it was defined, or null if no module with
// that name exists; optionalProperty the property name of the locally
// imported type on the module object, if destructuring-style importing was
// used, or null if this is a namespace import; scopeToDeclareIn the scope in
// which localNameNode is defined.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineModuleImport(Node, ScopedName, String, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) DefineModuleImport(
	localNameNode rhino.Node,
	importedModuleObject ScopedName,
	optionalProperty *string,
	scopeToDeclareIn *TypedScope) {
	if importedModuleObject == nil {
		// We could not find the module defining this import. If the compiler is assuming forward
		// declared types (or if explicitly forward declared), create a NamedType so it resolves to
		// NoResolvedType. Otherwise, fall back to unknownType to avoid spurious type check errors.
		var fallbackType rhino.JSType
		if b.outer.compiler.GetOptions().AssumeForwardDeclaredForMissingTypes() {
			fallbackType = b.outer.typeRegistry.CreateNamedType(
				b.currentScope,
				localNameNode.GetString(),
				localNameNode.GetSourceFileName(),
				localNameNode.GetLineno(),
				localNameNode.GetCharno())
		} else {
			fallbackType = b.outer.unknownType
		}
		b.defineSlot(typedScopeSlotSpec{
			declNode: localNameNode,
			varName:  localNameNode.GetString(),
			scope:    scopeToDeclareIn,
			type_:    fallbackType,
			allowLaterTypeInference: fallbackType.IsUnknownType() &&
				!b.outer.compiler.GetOptions().AssumeForwardDeclaredForMissingTypes(),
		})
		return
	}

	var exportedName ScopedName
	if optionalProperty == nil {
		exportedName = importedModuleObject
	} else {
		// If this is a destucuring-style import, find its type.
		exportedName = ScopedName_Of(
			importedModuleObject.GetName()+"."+*optionalProperty,
			importedModuleObject.GetScopeRoot())
	}
	// Try getting the actual scope. The scope will be null in the following cases:
	//   - Someone has required a module that does not exist at all.
	//   - Someone has requireType'd or forwardDeclare'd a module that exists, but does not have
	//     an associated scope yet.
	var exportScope *TypedScope
	if !exportedName.GetScopeRoot().IsNil() {
		exportScope = b.outer.memoized[exportedName.GetScopeRoot()]
	}

	// The scope is null for modules that were not visited yet.
	if exportScope != nil {
		type_ := rhino.StaticTypedScope_LookupQualifiedNameString(exportScope, exportedName.GetName())

		// The type is null if either the name is inferred or this is an early ref.
		if type_ != nil {
			b.DeclareAliasTypeIfRvalueIsAliasable(
				localNameNode, rhino.QualifiedName_Of(exportedName.GetName()), type_, exportScope)

			b.defineSlot(typedScopeSlotSpec{
				declNode:                localNameNode,
				varName:                 localNameNode.GetString(),
				scope:                   scopeToDeclareIn,
				type_:                   type_,
				allowLaterTypeInference: type_ == nil,
			})
			return
		}
	}
	// Defer defining this name until after we have visited the entire AST.
	b.outer.weakImports = append(b.outer.weakImports,
		NewTypedScopeCreator_WeakModuleImport(b.outer, localNameNode, exportedName, scopeToDeclareIn))
}

// typedScopeSlotSpec describes a typed variable slot to declare in a scope.
type typedScopeSlotSpec struct {
	declNode                rhino.Node
	varName                 string
	scope                   *TypedScope
	type_                   rhino.JSType
	allowLaterTypeInference bool
	forGoogProvidedName     bool
}

func slotVarNameFromDeclNode(declNode rhino.Node) string {
	parent := declNode.GetParent()
	if declNode.IsName() {
		jutil.CheckArgument(
			parent.IsFunction() ||
				parent.IsClass() ||
				NodeUtil_IsNameDeclaration(parent) ||
				parent.IsParamList() ||
				(parent.IsRest() && parent.GetParent().IsParamList()) ||
				parent.IsCatch())
	} else {
		jutil.CheckArgument(
			declNode.IsGetProp() && (parent.IsAssign() || parent.IsExprResult()))
	}
	name, ok := declNode.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerExceptionMsg("variableName not set"))
	}
	return name
}

// defineSlot defines a typed variable slot and performs related global/ctor
// bookkeeping.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.defineSlot()
func (b *TypedScopeCreator_scopeBuilder) defineSlot(d typedScopeSlotSpec) {
	rhino.CheckNotNullMsg(d.declNode, "declarationNode not set")
	jutil.CheckStateMsg1(
		d.declNode.IsName() ||
			d.declNode.IsGetProp() ||
			NodeUtil_MayBeObjectLitKey(d.declNode) ||
			d.declNode.IsModuleBody() ||
			d.declNode.IsExport() ||
			(d.forGoogProvidedName && NodeUtil_IsExprCall(d.declNode)),
		"declaration node must be an lvalue or goog.provide call, found %s",
		d.declNode)
	jutil.CheckStateMsg(d.allowLaterTypeInference || d.type_ != nil, "null type but inference not allowed")
	jutil.CheckState(d.varName != "")
	jutil.CheckNotNull(d.scope)
	variableName := d.varName

	parent := d.declNode.GetParent()

	scopeToDeclareIn := d.scope

	isGlobalVar := d.declNode.IsName() && scopeToDeclareIn.IsGlobal()
	shouldDeclareOnGlobalThis :=
		(isGlobalVar && (parent.IsVar() || parent.IsFunction())) ||
			(d.forGoogProvidedName && !strings.Contains(variableName, "."))

	// TODO(sdh): Remove this special case.  It is required to reproduce the original
	// non-block-scoped behavior, which is depended on in several places including
	// https://github.com/angular/tsickle/issues/761.  But it's more correct to always
	// declare on the owner scope.  Once all the bugs are fixed, this should be removed.
	// We may be able to get by with checking a "declared" function's source for jsdoc.
	if scopeToDeclareIn != b.currentHoistScope &&
		scopeToDeclareIn.IsGlobal() &&
		scopeToDeclareIn.HasOwnSlot(variableName) &&
		!d.forGoogProvidedName {
		scopeToDeclareIn = b.currentHoistScope
	}

	// The input may be null if we are working with a AST snippet. So read
	// the extern info from the node.

	// declared in closest scope?
	input := b.GetCompilerInput()
	if !scopeToDeclareIn.CanDeclare(variableName) {
		oldVar := scopeToDeclareIn.GetVar(variableName)
		b.outer.validator.ExpectUndeclaredVariable(
			"", input, d.declNode, parent, oldVar, variableName, d.type_)
	} else {
		if d.type_ != nil {
			b.outer.SetDeferredType(d.declNode, d.type_)
		}

		b.Declare(
			scopeToDeclareIn,
			variableName,
			d.declNode,
			d.type_,
			input,
			d.allowLaterTypeInference)
	}

	// We need to do some additional work for constructors and interfaces.
	fnType := rhino.JSType_ToMaybeFunctionType(d.type_)
	if fnType != nil &&
		// We don't want to look at empty function types.
		!d.type_.IsEmptyType() {

		// We want to make sure that when we declare a new instance type
		// (with @constructor) that there's actually a ctor for it.
		// This doesn't apply to structural constructors (like
		// function(new:Array). Checking the constructed type against
		// the variable name is a sufficient check for this.
		if fnType.IsConstructor() || fnType.IsInterface() {
			b.FinishConstructorDefinition(
				d.declNode, variableName, fnType, scopeToDeclareIn, input)
		}
	}

	if shouldDeclareOnGlobalThis {
		globalThis := b.outer.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GLOBAL_THIS)
		if d.allowLaterTypeInference {
			inferredType := d.type_
			if inferredType == nil {
				inferredType = b.outer.GetNativeType(rhino.JSTypeNative_NO_TYPE)
			}
			globalThis.DefineInferredProperty(
				variableName,
				inferredType,
				d.declNode)
		} else {
			globalThis.DefineDeclaredProperty(variableName, d.type_, d.declNode)
		}
	}

	if isGlobalVar &&
		"Window" == variableName &&
		d.type_ != nil &&
		d.type_.IsFunctionType() &&
		d.type_.IsConstructor() {
		globalThisCtor :=
			b.outer.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GLOBAL_THIS).GetConstructor()
		globalThisCtor.GetInstanceType().ClearCachedValues()
		b.outer.typeRegistry.ResetImplicitPrototype(
			globalThisCtor.GetPrototype(), d.type_.ToMaybeFunctionType().GetInstanceType())
	}
}

// typedScopeCreator_sourceNameArg passes the AbstractScopeBuilder's
// sourceName (a Java String that may be null) to
// TypeValidator.expectUndeclaredVariable, whose Go port takes it as a Go
// string (and, like the Java, never reads it).
//
// Go-only: language-adapter — Java null strings are nil *string.
func typedScopeCreator_sourceNameArg(sourceName *string) string {
	if sourceName == nil {
		return ""
	}
	return *sourceName
}

// Declare declares a variable with the given name and type on the given
// scope, returning the newly-declared TypedVar. Additionally checks the
// escapedVarNames and assignedVarNames maps (which were populated during the
// FirstOrderFunctionAnalyzer and marks the result as escaped or assigned
// exactly once if appropriate.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declare(TypedScope, String, Node, JSType, CompilerInput, boolean)
func (b *TypedScopeCreator_scopeBuilder) Declare(
	scope *TypedScope, name string, n rhino.Node, type_ rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	var_, wasReserved := scope.declareWithReserved(name, n, type_, input, inferred)
	if !scope.IsGlobal() {
		var flags uint8
		if len(b.outer.varFlags) > 0 {
			scopedName := typedScopeCreator_ScopedNameKey{name: name, scopeRoot: scope.GetRootNode()}
			flags = b.outer.varFlags[scopedName]
		}
		if (flags & typedScopeCreator_VarFlagEscaped) != 0 {
			var_.MarkEscaped()
		}
		if wasReserved && (flags&typedScopeCreator_VarFlagNotSingleAssigned) == 0 {
			var_.MarkAssignedExactlyOnce()
		}
	}
	return var_
}

// Java: TypedScopeCreator.AbstractScopeBuilder.finishConstructorDefinition(Node, String, FunctionType, TypedScope, CompilerInput)
func (b *TypedScopeCreator_scopeBuilder) FinishConstructorDefinition(
	declarationNode rhino.Node,
	variableName string,
	fnType rhino.FunctionType,
	scopeToDeclareIn *TypedScope,
	input *CompilerInput) {
	// Declare var.prototype in the scope chain.
	superClassCtor := fnType.GetSuperClassConstructor()
	prototypeSlot := fnType.GetSlot("prototype")

	prototypeName := variableName + ".prototype"

	// There are some rare cases where the prototype will already
	// be declared. See TypedScopeCreatorTest#testBogusPrototypeInit.
	// Fortunately, other warnings will complain if this happens.
	prototypeVar := scopeToDeclareIn.GetVar(prototypeName)
	if prototypeVar != nil && prototypeVar.Scope() == scopeToDeclareIn {
		scopeToDeclareIn.Undeclare(prototypeVar)
	}

	scopeToDeclareIn.Declare(
		prototypeName,
		declarationNode,
		prototypeSlot.GetType(),
		input,
		// declared iff there's an explicit supertype
		superClassCtor == nil ||
			superClassCtor.GetInstanceType().Equals(b.outer.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)))
}

// IsLValueRootedInGlobalScope checks if the given node is a property of a
// name in the global scope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isLValueRootedInGlobalScope(Node)
func (b *TypedScopeCreator_scopeBuilder) IsLValueRootedInGlobalScope(n rhino.Node) bool {
	return b.GetLValueRootScope(n).IsGlobal()
}

// GetLValueRootScope returns the scope for the name of the given node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getLValueRootScope(Node)
func (b *TypedScopeCreator_scopeBuilder) GetLValueRootScope(n rhino.Node) *TypedScope {
	root := NodeUtil_GetBestLValueRoot(n)
	if !root.IsNil() {
		if root.IsName() {
			nameParent := root.GetParent()
			switch nameParent.GetToken() {
			case rhino.Token_VAR:
				return b.currentHoistScope
			case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION, rhino.Token_PARAM_LIST, rhino.Token_CATCH:
				return b.currentScope
			case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
				// TODO(bradfordcsmith): Handle array destructuring REST
				jutil.CheckStateMsg1(nameParent.GetParent().IsParamList(), "%s", nameParent)
				return b.currentScope
			default:
				if b.IsGoogModuleExports(root) {
					// Ensure that 'exports = class {}' in a goog.module returns the module scope.
					return b.currentScope
				}
				var_ := b.currentScope.GetVarByID(root.GetStringID())
				if var_ != nil {
					return var_.Scope()
				}
			}
		} else if root.IsThis() || root.IsSuper() {
			// We want the enclosing function scope, or the global scope if not in a function.
			return b.currentHoistScope.GetScopeOfThis()
		}
	}
	return b.currentHoistScope.GetGlobalScope()
}
