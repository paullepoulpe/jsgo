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

// DeclareAliasTypeIfRvalueIsAliasable: for a const alias, like
// `const alias = other.name`, this may declare `alias` as a type name,
// depending on what other.name is defined to be.
//
// This method recognizes three kinds of type aliases: @typedefs,
// @constructor/@interface types, and @enums.
//
// Given any of those three types, this method redeclares the aliasing name
// in the typeRegistry. For @typedefs and global @enums, this method also
// marks the qualified name referring to the type as non-nullable by
// default.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(Node, QualifiedName, JSType, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) DeclareAliasTypeIfRvalueIsAliasable(
	lValue rhino.Node,
	rValue rhino.QualifiedName,
	rValueType rhino.JSType,
	rValueLookupScope *TypedScope) {
	if rValue == nil {
		return
	}
	b.DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
		lValue.GetQualifiedName(), lValue, rValue, rValueType, rValueLookupScope, b.currentScope)
}

// DeclareAliasTypeIfRvalueIsAliasableWithLValueName: for a const alias, like
// `const alias = other.name`, this may declare `alias` as a type name,
// depending on what other.name is defined to be.
//
// NOTE: in most cases, call the version with fewer arguments. This version
// only exists to handle goog.declareLegacyNamespace, which is strange
// compared to normal type aliasing because 1) there's no GETPROP node
// representing the lvalue and 2) the type is declared in the global scope,
// not the current module-local scope.
//
// lValueName is the fully qualified lValue name, if any. If null, all this
// method will do is propagate the @typedef Node annotation to
// actualLvalueNode. aliasDeclarationScope is the scope in which to declare
// the alias name. In most cases, this should just be the currentScope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(String, Node, QualifiedName, JSType, TypedScope, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
	lValueName *string,
	actualLvalueNode rhino.Node,
	rValue rhino.QualifiedName,
	rValueType rhino.JSType,
	rValueLookupScope *TypedScope,
	aliasDeclarationScope *TypedScope) {
	// NOTE: this allows some strange patterns such allowing instance properties
	// to be aliases of constructors, and then creating a local alias of that to be
	// used as a type name.  Consider restricting this.

	if rValue == nil {
		return
	}

	// Look for a @typedef annotation on the definition node
	definitionNode := b.GetDefinitionNode(rValue, rValueLookupScope)
	if !definitionNode.IsNil() {
		typedefType := definitionNode.GetTypedefTypeProp()
		if typedefType != nil {
			// Propagate typedef type to typedef aliases.
			actualLvalueNode.SetTypedefTypeProp(typedefType)
			if lValueName != nil {
				b.outer.typeRegistry.IdentifyNonNullableName(aliasDeclarationScope, *lValueName)
				b.outer.typeRegistry.DeclareType(aliasDeclarationScope, *lValueName, typedefType)
			}
			return
		}
	}

	if lValueName == nil {
		return
	}

	// Check if the provided rValueType indicates that we should declare this type
	// Note that we only look for enums and constructors/interfaces here: this step cannot work
	// for @typedefs. The 'type' of the TypedVar representing a @typedef'd name is the None type,
	// not the @typedef'd type.
	if rValueType != nil &&
		rValueType.IsFunctionType() &&
		rValueType.ToMaybeFunctionType().HasInstanceType() {
		// Look for @constructor/@interface by checking if the RHS has an instance type
		functionType := rValueType.ToMaybeFunctionType()
		b.outer.typeRegistry.DeclareType(aliasDeclarationScope, *lValueName, functionType.GetInstanceType())
		return
	}

	if rValueType != nil && rValueType.IsEnumType() {
		// Look for cases where the rValue is an Enum namespace
		b.outer.typeRegistry.DeclareType(
			aliasDeclarationScope, *lValueName, rValueType.ToMaybeEnumType().GetElementsType())
		b.outer.typeRegistry.IdentifyNonNullableName(aliasDeclarationScope, *lValueName)
	}
}

// IsGoogModuleExports reports whether this lvalue is either `exports`,
// `exports.x`, or a string key in `exports = {x}`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isGoogModuleExports(Node)
func (b *TypedScopeCreator_scopeBuilder) IsGoogModuleExports(lValue rhino.Node) bool {
	if b.GetModule() == nil || lValue.IsNil() {
		return false
	}
	if _, ok := b.outer.undeclaredNamesForClosure[lValue]; ok {
		// includes "exports" and "exports.x"
		return true
	}
	// ASSIGN
	//   NAME "exports"
	//   OBJECT_LIT
	//     STRING_KEY "x"
	//     [value]
	if !lValue.IsStringKey() || !lValue.GetParent().IsObjectLit() {
		return false
	}
	assign, ok := ast.AsAssign(lValue.GetGrandparent())
	if !ok {
		return false
	}
	target := assign.Target().Node()
	if !target.MatchesName("exports") {
		return false
	}
	_, ok = b.outer.undeclaredNamesForClosure[target]
	return ok
}

// GetDefinitionNode returns the AST node associated with the definition, if
// any.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDefinitionNode(QualifiedName, TypedScope)
func (b *TypedScopeCreator_scopeBuilder) GetDefinitionNode(qname rhino.QualifiedName, scope *TypedScope) rhino.Node {
	if qname.IsSimple() {
		var_ := scope.GetVar(qname.GetComponent())
		if var_ != nil {
			return var_.GetNameNode()
		}
		return rhino.NilNode
	}
	parent := rhino.ObjectType_Cast(scope.LookupQualifiedName(qname.GetOwner()))
	if parent != nil {
		return parent.GetPropertyDefSite(qname.GetComponent())
	}
	return rhino.NilNode
}

// GetDeclaredRValueType checks for common idioms of a typed R-value
// assigned to a const L-value.
//
// Normally, we would only want this sort of propagation to happen under
// type inference. But we want a declared const to be nameable in a type
// annotation, so we need to figure out the type before we try to resolve
// the annotation.
//
// lValue is the lvalue node if this is a simple assignment, null for
// destructuring.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredRValueType(Node, Node)
func (b *TypedScopeCreator_scopeBuilder) GetDeclaredRValueType(lValue rhino.Node, rValue rhino.Node) rhino.JSType {
	// If rValue has a type-cast, we use the type in the type-cast.
	rValueInfo := rValue.GetJSDocInfo()
	if rValue.IsCast() && rValueInfo != nil && rValueInfo.HasType() {
		return rValueInfo.GetType().Evaluate(b.currentScope, b.outer.typeRegistry)
	}

	// Check if the type has already been computed during scope-creation.
	// This is mostly useful for literals like BOOLEAN, NUMBER, STRING, and
	// OBJECT_LITERAL
	type_ := rValue.GetJSType()
	if type_ != nil && !type_.IsUnknownType() {
		return type_
	}

	// If rValue is a name, try looking it up in the current scope.
	if qname, ok := rValue.GetQualifiedNameString(); ok {
		return rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, qname)
	}

	// Check for simple invariant operations, such as "!x" or "+x" or "''+x"
	if NodeUtil_IsBooleanResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	}

	if NodeUtil_IsNumericResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	}

	if NodeUtil_IsBigIntResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	}

	if NodeUtil_IsStringResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	}

	if nw, ok := ast.AsNew(rValue); ok {
		if qname, ok := nw.Callee().Node().GetQualifiedNameString(); ok {
			targetType := rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, qname)
			if targetType != nil {
				fnType := targetType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
				if fnType != nil && fnType.HasInstanceType() {
					return fnType.GetInstanceType()
				}
			}
		}
	}

	if call, ok := ast.AsCall(rValue); ok && call.Callee().Node().MatchesName("Symbol") {
		// Type calls like `Symbol('foo')`
		symbolVar := b.currentScope.GetVarByID(typedScopeSymbolID)
		if symbolVar != nil && symbolVar.Scope().IsGlobal() {
			return b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
		}
	}

	// Check for a very specific JS idiom:
	// var x = x || TYPE;
	// This is used by Closure's base namespace for esoteric
	// reasons, so we only really care about that case.
	if rValue.IsOr() {
		orExpr := ast.MustBinary(rValue)
		if leftName, ok := orExpr.Left().AsName(); ok {
			if lvalName, ok := ast.AsName(lValue); ok && leftName.Ident() == lvalName.Ident() {
				type_ = orExpr.Right().Node().GetJSType()
				if type_ != nil && !type_.IsUnknownType() {
					return type_
				}
			}
		}
	}

	return nil
}

// CheckForClassDefiningCalls looks for class-defining calls. Because JS has
// no 'native' syntax for defining classes, this is often very
// coding-convention dependent and business-logic heavy.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.checkForClassDefiningCalls(Node)
func (b *TypedScopeCreator_scopeBuilder) CheckForClassDefiningCalls(n rhino.Node) {
	relationship := b.outer.codingConvention.GetClassesDefinedByCall(n)
	if relationship != nil {
		superClass :=
			TypeValidator_GetInstanceOfCtor(
				rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, relationship.superclassName))
		subClass :=
			TypeValidator_GetInstanceOfCtor(
				rhino.StaticTypedScope_LookupQualifiedNameString(b.currentScope, relationship.subclassName))
		if superClass != nil && subClass != nil {
			// superCtor and subCtor might be structural constructors
			// (like {function(new:Object)}) so we need to resolve them back
			// to the original ctor objects.
			superCtor := superClass.GetConstructor()
			subCtor := subClass.GetConstructor()
			if superCtor != nil && subCtor != nil {
				b.outer.codingConvention.ApplySubclassRelationship(
					rhino.NewNominalTypeBuilder(superCtor, superClass),
					rhino.NewNominalTypeBuilder(subCtor, subClass),
					relationship.type_)
			}
		}
	}

	singletonGetterClassName := b.outer.codingConvention.GetSingletonGetterClassName(n)
	if singletonGetterClassName != nil {
		objectType :=
			rhino.ObjectType_Cast(b.outer.typeRegistry.GetType(b.currentScope, *singletonGetterClassName))
		if objectType != nil {
			functionType := objectType.GetConstructor()

			if functionType != nil {
				getterType := b.outer.typeRegistry.CreateFunctionType(objectType)
				b.outer.codingConvention.ApplySingletonGetter(
					rhino.NewNominalTypeBuilder(functionType, objectType), getterType)
			}
		}
	}

	objectLiteralCast := b.outer.codingConvention.GetObjectLiteralCast(n)
	if objectLiteralCast != nil {
		if objectLiteralCast.diagnosticType == nil {
			type_ :=
				rhino.ObjectType_Cast(b.outer.typeRegistry.GetType(b.currentScope, *objectLiteralCast.typeName))
			if type_ != nil && type_.GetConstructor() != nil {
				b.outer.SetDeferredType(objectLiteralCast.objectNode, type_)
				objectLiteralCast.objectNode.PutBooleanProp(rhino.Node_REFLECTED_OBJECT, true)
			} else {
				b.outer.Report(JSError_MakeWithNode(n, TypedScopeCreator_CONSTRUCTOR_EXPECTED))
			}
		} else {
			b.outer.Report(JSError_MakeWithNode(n, objectLiteralCast.diagnosticType))
		}
	}
}

// MaybeDeclareGetElem declares a computed property on its owner type.
//
// info is the doc info for this property; n a top-level GETELEM node (it
// should not be contained inside another GETPROP); rhsValue the node that n
// is being initialized to, or null if this is a stub declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.maybeDeclareGetElem(JSDocInfo, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) MaybeDeclareGetElem(info *rhino.JSDocInfo, n rhino.Node, rhsValue rhino.Node) {
	ge := ast.MustGetElem(n)
	ownerNode := ge.Object().Node()
	key := ge.Index().Node()
	keyType := TypedScopeCreator_ExtractKnownSymbolKey(b.currentScope, key)
	if keyType == nil {
		return
	}
	valueType := b.GetDeclaredType(info, n, rhsValue, nil)
	if valueType == nil {
		return
	}
	ownerName, ok := ownerNode.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	ownerType := b.GetObjectSlot(ownerName)
	if ownerType == nil {
		return
	}
	b.DeclarePropertyIfNamespaceType(
		ownerType,
		ownerNode,
		rhino.NewProperty_SymbolKey(keyType),
		valueType,
		n,
		/* jsdocInfo= */ nil)
}

// MaybeDeclareQualifiedName declares the symbol for a qualified name in the
// current scope.
//
// info is the doc info for this property; n a top-level GETPROP node (it
// should not be contained inside another GETPROP); parent the parent of n;
// rhsValue the node that n is being initialized to, or null if this is a
// stub declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.maybeDeclareQualifiedName(NodeTraversal, JSDocInfo, Node, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) MaybeDeclareQualifiedName(
	t *NodeTraversal, info *rhino.JSDocInfo, n rhino.Node, parent rhino.Node, rhsValue rhino.Node) {
	isTypedef := b.IsValidTypedefDeclaration(n, info)
	if isTypedef {
		b.DeclareTypedefType(n, info)
	}

	gp := ast.MustGetProp(n)
	ownerNode := gp.Object().Node()
	var qName string
	var hasQName bool
	if !rhsValue.IsNil() && rhsValue.IsFunction() {
		if origID := rhsValue.GetOriginalNameID(); origID != 0 {
			origStr := rhino.RhinoStringPool_StringFromID(origID)
			if n.MatchesQualifiedName(origStr) {
				qName = origStr
				hasQName = true
			}
		}
	}
	if !hasQName {
		qName, hasQName = n.GetQualifiedNameString()
	}
	propName := gp.Property()
	jutil.CheckArgument(hasQName)
	ownerNameStr := qName[:len(qName)-len(propName)-1]

	// Precedence of type information on GETPROPs:
	// 1) @type annotation / @enum annotation
	// 2) ASSIGN to FUNCTION literal
	// 3) @param/@return annotation (with no function literal)
	// 4) ASSIGN to something marked @const
	// 5) ASSIGN to anything else
	//
	// 1, 3, and 4 are declarations, 5 is inferred, and 2 is a declaration iff
	// the function has JsDoc or has not been declared before.
	//
	// FUNCTION literals are special because TypedScopeCreator is very smart
	// about getting as much type information as possible for them.

	// Determining type for #1 + #2 + #3 + #4
	valueType := b.GetDeclaredType(info, n, rhsValue, nil)
	if valueType == nil && !rhsValue.IsNil() {
		// Determining type for #5
		valueType = rhsValue.GetJSType()
	}

	// Function prototypes are special.
	// It's a common JS idiom to do:
	// F.prototype = { ... };
	// So if F does not have an explicitly declared super type,
	// allow F.prototype to be redefined arbitrarily.
	if "prototype" == propName {
		qVar := b.currentScope.GetVar(qName)
		if qVar != nil {
			// If the programmer has declared that F inherits from Super,
			// and they assign F.prototype to an object literal,
			// then they are responsible for making sure that the object literal's
			// implicit prototype is set up appropriately. We just obey
			// the @extends tag.
			qVarType := rhino.ObjectType_Cast(qVar.GetType())
			if qVarType != nil && !rhsValue.IsNil() && rhsValue.IsObjectLit() {
				b.outer.typeRegistry.ResetImplicitPrototype(
					rhsValue.GetJSType(), qVarType.GetImplicitPrototype())
			} else if !qVar.IsTypeInferred() {
				// If the programmer has declared that F inherits from Super,
				// and they assign F.prototype to some arbitrary expression,
				// there's not much we can do. We just ignore the expression,
				// and hope they've annotated their code in a way to tell us
				// what props are going to be on that prototype.
				return
			}

			qVar.Scope().Undeclare(qVar)
		}
	}

	if valueType == nil {
		if parent.IsExprResult() {
			// t is mutable so make sure to capture the current state before the lambda.
			isExtern := t.GetInput() != nil && t.GetInput().IsExtern()
			b.deferAction(b.currentScope.GetRootNode(), func() { b.ResolveStubDeclaration(n, isExtern, ownerNameStr) })
		}

		return
	}

	inferred := b.IsQualifiedNameInferred(qName, hasQName, n, info, rhsValue, valueType)
	if !inferred {
		ownerType := b.GetObjectSlot(ownerNameStr)
		if ownerType != nil {
			b.DeclarePropertyIfNamespaceType(
				ownerType, ownerNode, rhino.NewProperty_StringKey(propName), valueType, n /* jsdocInfo= */, nil)
		}

		// this is a memory optimization: we don't need to declare .prototype props in the scope.
		// NOTE: in theory we could possibly extend this to more kinds of declarations and only
		// declare simple (non-qualified) names in the scope, but that seems to cause a lot more
		// issues with older code.
		declareInScope := ownerType == nil || !ownerType.IsFunctionPrototypeType()

		// If the property is already declared, the error will be caught when we try to declare it
		// in the current scope.
		if declareInScope {
			b.defineSlot(typedScopeSlotSpec{
				declNode:                n,
				varName:                 qName,
				scope:                   b.GetLValueRootScope(n),
				type_:                   valueType,
				allowLaterTypeInference: inferred,
			})
		}
	}
}

// IsQualifiedNameInferred determines whether a qualified name is inferred.
// NOTE(nicksantos): Determining whether a property is declared or not is
// really really obnoxious.
//
// The problem is that there are two (equally valid) coding styles:
//
// (function() { /* The authoritative definition of goog.bar. / goog.bar =
// function() {}; })();
//
// function f() { goog.bar(); /* Reset goog.bar to a no-op. / goog.bar =
// function() {}; }
//
// In a dynamic language with first-class functions, it's very difficult to
// know which one the user intended without looking at lots of contextual
// information (the second example demonstrates a small case of this, but
// there are some really pathological cases as well).
//
// The current algorithm checks if either the declaration has JsDoc type
// information, or @const with a known type, or a function literal with a
// name we haven't seen before.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isQualifiedNameInferred(String, Node, JSDocInfo, Node, JSType)
func (b *TypedScopeCreator_scopeBuilder) IsQualifiedNameInferred(
	qName string,
	hasQName bool,
	n rhino.Node,
	info *rhino.JSDocInfo,
	rhsValue rhino.Node,
	valueType rhino.JSType) bool {
	// Prototypes of constructors and interfaces are always declared.
	if hasQName && strings.HasSuffix(qName, ".prototype") {
		className := qName[:len(qName)-len(".prototype")]
		slot := b.currentScope.GetVar(className)
		var classType rhino.JSType
		if slot != nil {
			classType = slot.GetType()
		}
		if classType != nil && (classType.IsConstructor() || classType.IsInterface()) {
			return false
		}
	}
	// treat "foo = bar = <VAL>" the same as "foo = <VAL>"
	rhsValue = b.UnwrapIfAssign(rhsValue)

	// If the jsdoc or RHS specifies a concrete type, it's not inferred.
	if (info != nil &&
		(info.HasType() ||
			info.HasEnumParameterType() ||
			b.IsValidTypedefDeclaration(n, info) ||
			FunctionTypeBuilder_IsFunctionTypeDeclaration(info) ||
			(!rhsValue.IsNil() && rhsValue.IsFunction()))) ||
		b.IsTypedConstantDeclaration(info, n, valueType) {
		return false
	}

	// At this point, we're pretty sure it's inferred, since there's neither
	// useful jsdoc info, nor a useful const or doc'd function RHS.  But
	// there's still one case where it may still not be: if the RHS is a
	// class or function that is not
	//   (1) a scoped qualified name (i.e. this.b.c or super.b.c),
	//   (2) already declared in a scope,
	//   (3) assigned in a conditional block, or
	//   (4) escaped to a closure,
	// then we treat it as if it is declared, rather than inferred.
	// Stubs and other values are always considered inferred at this point.
	if rhsValue.IsNil() || (!rhsValue.IsFunction() && !rhsValue.IsClass()) {
		return true
	}

	// "Scoped" qualified names (e.g. this.b.c or super.d) are inferred.
	if !n.IsUnscopedQualifiedName() {
		return true
	}

	// If this qname is already declared then treat this definition as inferred.
	ownerScope := b.GetLValueRootScope(n)
	if ownerScope != nil && ownerScope.HasOwnSlot(qName) {
		return true
	}

	// Check if this is in a conditional block.
	// Functions assigned in conditional blocks are inferred.
	if b.HasControlStructureAncestor(n.GetParent()) {
		return true
	}

	// Check if this is assigned in an inner scope.
	// Functions assigned in inner scopes are inferred.
	if ownerScope != nil && len(b.outer.varFlags) > 0 {
		if (b.outer.varFlags[typedScopeCreator_ScopedNameKey{name: qName, scopeRoot: ownerScope.GetRootNode()}] & typedScopeCreator_VarFlagEscaped) != 0 {
			return true
		}
	}

	return false
}

// Java: TypedScopeCreator.AbstractScopeBuilder.unwrapIfAssign(Node)
func (b *TypedScopeCreator_scopeBuilder) UnwrapIfAssign(maybeAssign rhino.Node) rhino.Node {
	for {
		assign, ok := ast.AsAssign(maybeAssign)
		if !ok {
			return maybeAssign
		}
		maybeAssign = assign.Value().Node()
	}
}

// DeclareProvidedNs: given a `goog.provide()` or legacy `goog.module()` call
// and implicit ProvidedName, declares the name in the global scope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareProvidedNs(Node, ProvidedName)
func (b *TypedScopeCreator_scopeBuilder) DeclareProvidedNs(provideCall rhino.Node, providedName *ProcessClosureProvidesAndRequires_ProvidedName) {
	// Redefine this name if we haven't already added a provide definition.
	// Note: in some cases, this will cause a redefinition error.
	anonymousObjectType := b.outer.typeRegistry.CreateAnonymousObjectType(nil)
	b.defineSlot(typedScopeSlotSpec{
		declNode:                provideCall,
		varName:                 providedName.GetNamespace(),
		scope:                   b.currentScope.GetGlobalScope(),
		type_:                   anonymousObjectType,
		allowLaterTypeInference: false,
		forGoogProvidedName:     true,
	})

	namespace := rhino.QualifiedName_Of(providedName.GetNamespace())
	if !namespace.IsSimple() {
		ownerType := b.currentScope.LookupQualifiedName(namespace.GetOwner())
		if ownerType != nil && ownerType.IsObjectType() {
			ownerType.
				ToMaybeObjectType().
				DefineDeclaredProperty(namespace.GetComponent(), anonymousObjectType, provideCall)
		}
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.isTypedConstantDeclaration(JSDocInfo, Node, JSType)
func (b *TypedScopeCreator_scopeBuilder) IsTypedConstantDeclaration(info *rhino.JSDocInfo, n rhino.Node, valueType rhino.JSType) bool {
	return (NodeUtil_IsConstantDeclaration(info, n) || b.IsGoogModuleExports(n)) &&
		valueType != nil
}

// Java: TypedScopeCreator.AbstractScopeBuilder.hasControlStructureAncestor(Node)
func (b *TypedScopeCreator_scopeBuilder) HasControlStructureAncestor(n rhino.Node) bool {
	for !(n.IsScript() || n.IsFunction()) {
		if NodeUtil_IsControlStructure(n) {
			return true
		}
		n = n.GetParent()
	}
	return false
}

// GetObjectSlot finds the ObjectType associated with the given slot.
//
// slotName is the name of the slot to find the type in. Returns an object
// type, or null if this slot does not contain an object.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getObjectSlot(String)
func (b *TypedScopeCreator_scopeBuilder) GetObjectSlot(slotName string) rhino.ObjectType {
	ownerVar := b.currentScope.GetVar(slotName)
	if ownerVar != nil {
		ownerVarType := ownerVar.GetType()
		if ownerVarType == nil {
			return rhino.ObjectType_Cast(nil)
		}
		return rhino.ObjectType_Cast(ownerVarType.RestrictByNotNullOrUndefined())
	}
	return nil
}

// GetInheritedInterfacePropertyType: when a class has a stub for a
// property, and the property exists on a super interface, use that type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getInheritedInterfacePropertyType(ObjectType, String)
func (b *TypedScopeCreator_scopeBuilder) GetInheritedInterfacePropertyType(obj rhino.ObjectType, propName string) rhino.JSType {
	if obj != nil && obj.IsFunctionPrototypeType() {
		f := obj.GetOwnerFunction()
		for _, i := range f.GetImplementedInterfaces() {
			if i.HasProperty(propName) {
				return i.GetPropertyType(propName)
			}
		}
	}
	return nil
}

// ResolveStubDeclaration resolves any type-less stub declarations to
// unknown types if we could not find types for them during traversal. This
// method is only called as a deferred action after the root node is visted.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.resolveStubDeclaration(Node, boolean, String)
func (b *TypedScopeCreator_scopeBuilder) ResolveStubDeclaration(n rhino.Node, isExtern bool, ownerName string) {
	qName, ok := n.GetQualifiedNameString()
	propName := n.GetString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}

	// TODO(b/111216910): should this be getLValueRoot(n).hasOwnSlot(qName)?
	if b.currentScope.HasOwnSlot(qName) {
		return
	}

	// If we see a stub property, make sure to register this property
	// in the type registry.
	ownerType := b.GetObjectSlot(ownerName)
	inheritedType := b.GetInheritedInterfacePropertyType(ownerType, propName)
	var stubType rhino.JSType
	if inheritedType == nil {
		stubType = b.outer.unknownType
	} else {
		stubType = inheritedType
	}
	b.defineSlot(typedScopeSlotSpec{
		declNode:                n,
		varName:                 slotVarNameFromDeclNode(n),
		scope:                   b.GetLValueRootScope(n),
		type_:                   stubType,
		allowLaterTypeInference: true,
	})

	if ownerType != nil && (isExtern || ownerType.IsFunctionPrototypeType()) {
		// If this is a stub for a prototype, just declare it
		// as an unknown type. These are seen often in externs.
		ownerType.DefineInferredProperty(propName, stubType, n)
	} else {
		var registerOn rhino.JSType
		if ownerType == nil {
			registerOn = stubType
		} else {
			registerOn = ownerType
		}
		b.outer.typeRegistry.RegisterPropertyOnType(propName, registerOn)
	}
}

// IsValidTypedefDeclaration returns whether this is a valid declaration of a
// @typedef.
//
// candidate is a qualified name node; info JSDoc comments.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isValidTypedefDeclaration(Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) IsValidTypedefDeclaration(candidate rhino.Node, info *rhino.JSDocInfo) bool {
	if info == nil || !info.HasTypedefType() {
		return false
	}
	// `isUnscopedQualifiedName` excludes `this` and `super` properties.
	return candidate.IsUnscopedQualifiedName() && !NodeUtil_IsPrototypeProperty(candidate)
}

// DeclareTypedefType declares a typedef'd name in the JSTypeRegistry.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareTypedefType(Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) DeclareTypedefType(candidate rhino.Node, info *rhino.JSDocInfo) {
	typedef, ok := candidate.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}

	// TODO(nicksantos|user): This is a terrible, terrible hack
	// to bail out on recursive typedefs. We'll eventually need
	// to handle these properly.
	b.outer.typeRegistry.DeclareType(b.currentScope, typedef, b.outer.unknownType)

	realType := info.GetTypedefType().Evaluate(b.currentScope, b.outer.typeRegistry)
	if realType == nil {
		b.outer.Report(JSError_MakeWithNode(candidate, TypedScopeCreator_MALFORMED_TYPEDEF, typedef))
	} else {
		candidate.SetTypedefTypeProp(realType)
	}

	b.outer.typeRegistry.OverwriteDeclaredType(b.currentScope, typedef, realType)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.declarePropertyIfNamespaceType(ObjectType, Node, Property.Key, JSType, Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) DeclarePropertyIfNamespaceType(
	ownerType rhino.ObjectType,
	ownerNode rhino.Node,
	propName rhino.Property_Key,
	valueType rhino.JSType,
	declarationNode rhino.Node,
	jsdocInfo *rhino.JSDocInfo) {
	// Only declare this as an official property if it has not been
	// declared yet.
	if ownerType.HasOwnPropertyByKey(propName) && !ownerType.IsPropertyTypeInferredByKey(propName) {
		return
	}
	// Define the property if any of the following are true:
	//   (1) it's a non-native extern type. Native types are excluded here because we don't
	//       want externs of the form "/** @type {!Object} */ var api = {}; api.foo;" to
	//       cause a property "foo" to be declared on Object.
	//   (2) it's a non-instance type. This primarily covers static properties on
	//       constructors (which are FunctionTypes, not InstanceTypes).
	//   (3) it's an assignment to 'this', which covers instance properties assigned in
	//       constructors or other methods.
	isNonNativeExtern :=
		b.GetCompilerInput() != nil &&
			b.GetCompilerInput().IsExtern() &&
			!ownerType.IsNativeObjectType()
	if isNonNativeExtern || !ownerType.IsInstanceType() || ownerNode.IsThis() {
		// If the property is undeclared or inferred, declare it now.
		ownerType.DefineDeclaredPropertyByKey(propName, valueType, declarationNode)
		ownerType.SetPropertyJSDocInfoByKey(propName, jsdocInfo)
	}
} // end AbstractScopeBuilder
