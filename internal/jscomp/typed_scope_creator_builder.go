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

type scopeBuilderKind uint8

const (
	scopeBuilderNormal scopeBuilderKind = iota
	scopeBuilderFunction
	scopeBuilderClass
)

// TypedScopeCreator_scopeBuilder holds the unified state and traversal methods
// for normal, function, and class scope creation.
type TypedScopeCreator_scopeBuilder struct {
	kind scopeBuilderKind

	// Go-only: language-adapter — the enclosing TypedScopeCreator instance of
	// this inner class.
	outer *TypedScopeCreator

	// The scope that we're building.
	currentScope *TypedScope

	// The current hoist scope.
	currentHoistScope *TypedScope

	// The current source file that we're in.
	sourceName *string

	// The InputId of the current node.
	inputId       *rhino.InputId
	input         *CompilerInput
	inputResolved bool

	// Some actions need to be deferred, such as analyzing object literals with lends annotations,
	// or resolving type-less stubs. These actions are added to this map, keyed by the node that
	// should be waited for before running.
	deferredByNode map[rhino.Node][]*typedScopeDeferredAction
	deferredOrder  []*typedScopeDeferredAction
}

type typedScopeDeferredAction struct {
	fn   func()
	done bool
}

// Java: TypedScopeCreator.AbstractScopeBuilder(TypedScope)
func (b *TypedScopeCreator_scopeBuilder) init(kind scopeBuilderKind, outer *TypedScopeCreator, scope *TypedScope) {
	b.kind = kind
	b.outer = outer
	b.currentScope = scope
	b.currentHoistScope = scope.GetClosestHoistScope()
	b.sourceName = nil
	b.inputId = nil
	b.input = nil
	b.inputResolved = false
	b.deferredByNode = nil
	b.deferredOrder = nil
}

func (b *TypedScopeCreator_scopeBuilder) reset() {
	*b = TypedScopeCreator_scopeBuilder{}
}

func (b *TypedScopeCreator_scopeBuilder) getOwnerLValueNameString(ownerNode rhino.Node) (string, bool) {
	if ownerNode.IsNil() {
		return "", false
	}
	return nodeUtil_getBestLValueNameString(ownerNode)
}

func (b *TypedScopeCreator_scopeBuilder) deferAction(key rhino.Node, fn func()) {
	if b.deferredByNode == nil {
		b.deferredByNode = make(map[rhino.Node][]*typedScopeDeferredAction)
	}
	act := &typedScopeDeferredAction{fn: fn}
	b.deferredByNode[key] = append(b.deferredByNode[key], act)
	b.deferredOrder = append(b.deferredOrder, act)
}

func (b *TypedScopeCreator_scopeBuilder) setInputId(id *rhino.InputId) {
	if b.inputId != id {
		b.inputId = id
		b.input = nil
		b.inputResolved = false
	}
}

// GetCompilerInput returns the current compiler input.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getCompilerInput()
func (b *TypedScopeCreator_scopeBuilder) GetCompilerInput() *CompilerInput {
	if !b.inputResolved {
		b.input = b.outer.compiler.GetInput(b.inputId)
		b.inputResolved = true
	}
	return b.input
}

// Java: TypedScopeCreator.AbstractScopeBuilder.getModule()
func (b *TypedScopeCreator_scopeBuilder) GetModule() *Module {
	return b.currentScope.GetModule()
}

// Build traverses the scope root and builds it.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.build()
func (b *TypedScopeCreator_scopeBuilder) Build() {
	b.InitializeModuleScope(b.currentScope.GetRootNode())

	NodeTraversal_Builder_().
		SetCompiler(b.outer.compiler).
		SetCallback(b).
		SetScopeCreator(ScopeCreator_ASSERT_NO_SCOPES_CREATED).
		TraverseAtScope(b.currentScope)

	b.FinishDeclaringGoogModule()
}

// InitializeModuleScope builds the beginning of a module-scope. This can be
// an ES module or a goog.module.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.initializeModuleScope(Node)
func (b *TypedScopeCreator_scopeBuilder) InitializeModuleScope(moduleBody rhino.Node) {
	if b.GetModule() == nil {
		return
	}
	moduleType := b.GetModule().Metadata().ModuleType()
	switch moduleType {
	case ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE, ModuleMetadataMap_ModuleType_GOOG_MODULE:
		b.DeclareExportsInGoogModuleScope(b.GetModule(), moduleBody)
		b.outer.MarkGoogModuleExportsAsConst(moduleBody)
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		unresolvedImports := b.outer.moduleImportResolver.DeclareEsModuleImports(
			b.GetModule(),
			b.currentScope,
			b.outer.compiler.GetInput(NodeUtil_GetInputId(moduleBody)))
		for key, value := range unresolvedImports.All() {
			b.outer.weakImports = append(b.outer.weakImports,
				NewTypedScopeCreator_WeakModuleImport(b.outer, key, value, b.currentScope))
		}
	default:
		// This should have been a build-breaking DUPLICATE_NAMESPACE_AND_MODULE error in
		// GatherModuleMetadata
		panic(jutil.NewIllegalStateExceptionMsg(
			jutil.StringFormat(
				"Unexpected module type %s in module %s", moduleType, b.GetModule())))
	}
}

// DeclareExportsInGoogModuleScope ensures that the name `exports` is
// declared in goog.module scope.
//
// If a goog.module explicitly assigns to exports, we want to treat that
// assignment inside the scope as if it were a declaration:
// `const exports = ...`. This method only handles cases where we want to
// treat exports as implicitly declared.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareExportsInGoogModuleScope(Module, Node)
func (b *TypedScopeCreator_scopeBuilder) DeclareExportsInGoogModuleScope(googModule *Module, moduleBody rhino.Node) {
	for statement := moduleBody.GetFirstChild(); !statement.IsNil(); statement = statement.GetNext() {
		exprStmt, ok := ast.AsExprStmt(statement)
		if !ok {
			continue
		}
		if assign, ok := exprStmt.Expr().AsAssign(); ok && assign.Target().Node().MatchesName("exports") {
			return // found a direct assignment `exports = [...]`
		}
	}
	root := googModule.Metadata().RootNode()
	// Synthesize an object literal namespace 'exports'
	b.defineSlot(typedScopeSlotSpec{
		declNode:                root.GetFirstChild(),
		varName:                 "exports",
		scope:                   b.currentScope,
		type_:                   b.outer.typeRegistry.CreateAnonymousObjectType(nil),
		allowLaterTypeInference: false,
	})
}

// FinishDeclaringGoogModule declares goog.module.declareLegacyNamespace()
// names in the global scope and annotates the AST with type information
// about module exports.
//
// For example, given
//
//	goog.module('a.b.Foo');
//	goog.module.declareLegacyNamespace();
//
//	exports = class Foo {};
//
// This function is responsible for declaring a global name `a.b.Foo` and
// giving it the type from `exports = class Foo {};`
//
// Java: TypedScopeCreator.AbstractScopeBuilder.finishDeclaringGoogModule()
func (b *TypedScopeCreator_scopeBuilder) FinishDeclaringGoogModule() {
	if b.GetModule() == nil || !b.GetModule().Metadata().IsGoogModule() {
		return
	}
	exportsVar := jutil.CheckNotNull(b.currentScope.GetSlot("exports")).(*TypedVar)

	if b.GetModule().Metadata().IsLegacyGoogModule() {
		b.outer.typeRegistry.RegisterLegacyClosureNamespace(*b.GetModule().ClosureNamespace())
		moduleNamespace := rhino.QualifiedName_Of(*b.GetModule().ClosureNamespace())
		b.defineSlot(typedScopeSlotSpec{
			declNode:                exportsVar.GetNameNode(),
			varName:                 moduleNamespace.Join(),
			scope:                   b.currentScope.GetGlobalScope(),
			type_:                   exportsVar.GetType(),
			allowLaterTypeInference: exportsVar.IsTypeInferred(),
			forGoogProvidedName:     true,
		})

		if !moduleNamespace.IsSimple() && !exportsVar.IsTypeInferred() {
			parentType :=
				b.currentScope.GetGlobalScope().LookupQualifiedName(moduleNamespace.GetOwner())
			if parentType != nil && parentType.ToMaybeObjectType() != nil {
				// Declare the namespace on the parent name, propagating any JSDoc annotations
				// on "exports = Foo".
				typeDoc := NodeUtil_GetBestJSDocInfo(exportsVar.GetNameNode())
				b.DeclarePropertyIfNamespaceType(
					parentType.ToMaybeObjectType(),
					exportsVar.GetNameNode(),
					rhino.NewProperty_StringKey(moduleNamespace.GetComponent()),
					exportsVar.GetType(),
					exportsVar.GetNameNode(),
					typeDoc)
			}
		}
		b.DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
			b.GetModule().ClosureNamespace(),
			exportsVar.GetNameNode(), // Pretend that 'exports = '... is the lvalue node.
			rhino.QualifiedName_Of("exports"),
			exportsVar.GetType(),
			b.currentScope,
			b.currentScope.GetGlobalScope())
	} else {
		b.outer.typeRegistry.RegisterNonLegacyClosureNamespace(
			*b.GetModule().ClosureNamespace(), exportsVar.GetNameNode(), exportsVar.GetType())
	}
	// Store the type of the namespace on the AST for the convenience of later passes that want
	// to access it.
	rootNode := b.currentScope.GetRootNode()
	if rootNode.IsModuleBody() {
		rootNode.SetJSType(exportsVar.GetType())
	} else {
		// For goog.loadModule, give the `exports` parameter the correct type.
		jutil.CheckStateMsg1(rootNode.IsBlock(), "%s", rootNode)
		fn := ast.MustFunction(rootNode.GetParent())
		fn.Params().MustOnlyParam().Node().SetJSType(exportsVar.GetType())
	}
}

// ShouldTraverse is final in Java.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.shouldTraverse(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	b.setInputId(t.GetInputId())

	if n.IsFunction() || n.IsScript() || (parent.IsNil() && b.inputId != nil) {
		jutil.CheckNotNull(b.inputId)
		b.sourceName = NodeUtil_GetSourceName(n)
	}
	if parent.IsNil() || b.InCurrentScope(t) {
		b.visitPreorder(t, n, parent)
		return true
	} else if n.IsModuleBody() {
		// Visit modules pre-order. While this doesn't exactly match execution semantics, it
		// does match how the compiler rewrites modules into the global scope.
		b.outer.CreateScope(n, b.currentScope)
	} else if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		functionScope := b.outer.CreateScope(parent, b.currentScope).(*TypedScope)
		b.outer.CreateScope(n, functionScope)
	}
	return false
}

// Java: TypedScopeCreator.AbstractScopeBuilder.inCurrentScope(NodeTraversal)
func (b *TypedScopeCreator_scopeBuilder) InCurrentScope(t *NodeTraversal) bool {
	return t.GetScopeRoot() == b.currentScope.GetRootNode()
}

// Visit is final in Java.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.visit(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_scopeBuilder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	isScopeRoot := parent.IsNil()
	b.setInputId(t.GetInputId())
	if !isScopeRoot || !n.IsClass() {
		// Do not type literals or visit class scope roots
		b.AttachLiteralTypes(n)
		b.visitPostorder(t, n, parent)
	}
	if b.deferredByNode != nil {
		if !isScopeRoot {
			if actions, ok := b.deferredByNode[n]; ok {
				delete(b.deferredByNode, n)
				for _, act := range actions {
					act.done = true
				}
				for _, act := range actions {
					act.fn()
				}
			}
		} else {
			// Run *all* remaining deferred actions, in case any were missed.
			for i := 0; i < len(b.deferredOrder); i++ {
				act := b.deferredOrder[i]
				if !act.done {
					act.done = true
					act.fn()
				}
			}
		}
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.attachLiteralTypes(Node)
func (b *TypedScopeCreator_scopeBuilder) AttachLiteralTypes(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NULL:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_NULL_TYPE))
	case rhino.Token_VOID:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	case rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT_STRING:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	case rhino.Token_NUMBER:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	case rhino.Token_BIGINT:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
	case rhino.Token_TRUE, rhino.Token_FALSE:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	case rhino.Token_REGEXP:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_REGEXP_TYPE))
	case rhino.Token_OBJECTLIT:
		info := n.GetJSDocInfo()
		if info != nil && info.HasLendsName() {
			// Defer analyzing object literals with a @lends annotation until we
			// reach the root of the statement they're defined in.
			//
			// This ensures that if there are any @lends annotations on the object
			// literals, the type on the @lends annotation resolves correctly.
			//
			// For more information, see http://blickly.github.io/closure-compiler-issues/#314
			b.deferAction(NodeUtil_GetEnclosingStatement(n), func() { b.DefineObjectLiteral(n) })
		} else {
			b.DefineObjectLiteral(n)
		}
	case rhino.Token_CLASS:
		// NOTE(sdh): We can't handle function nodes here because they need special behavior to
		// deal with hoisting.  But since classes aren't hoisted, and may need to be handled in
		// such places as default method initializers (i.e. in a FunctionScope) or class extends
		// clauses (technically part of the ClassScope, but visited instead by the NormalScope),
		// they can be handled consistently in all scopes.
		b.DefineClassLiteral(n)
	// NOTE(johnlenz): If we ever support Array tuples,
	// we will need to handle them here as we do object literals
	// above.
	case rhino.Token_ARRAYLIT:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE))
	default:
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.defineObjectLiteral(Node)
func (b *TypedScopeCreator_scopeBuilder) DefineObjectLiteral(objectLit rhino.Node) {
	// Handle the @lends annotation.
	var type_ rhino.JSType
	info := objectLit.GetJSDocInfo()
	if info != nil && info.HasLendsName() {
		lendsName := info.GetLendsName().GetRoot().GetString()
		lendsVar := b.currentScope.GetVar(lendsName)
		if lendsVar == nil {
			b.outer.Report(JSError_MakeWithNode(objectLit, TypedScopeCreator_UNKNOWN_LENDS, lendsName))
		} else {
			type_ = lendsVar.GetType()
			if type_ == nil {
				type_ = b.outer.unknownType
			}
			if !type_.IsSubtypeOf(b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)) {
				b.outer.Report(JSError_MakeWithNode(objectLit, TypedScopeCreator_LENDS_ON_NON_OBJECT, lendsName, type_.String()))
				type_ = nil
			} else {
				objectLit.SetJSType(type_)
			}
		}
	}

	info = NodeUtil_GetBestJSDocInfo(objectLit)
	createEnumType := info != nil && info.HasEnumParameterType()
	if createEnumType {
		lValue := NodeUtil_GetBestLValue(objectLit)
		lValueName := NodeUtil_GetBestLValueName(lValue)
		type_ = b.CreateEnumTypeFromNodes(objectLit, lValueName, lValue, info)
	}

	if type_ == nil {
		type_ = b.outer.typeRegistry.CreateAnonymousObjectType(info)
	}

	b.outer.SetDeferredType(objectLit, type_)

	// If this is an enum, the properties were already taken care of above.
	b.ProcessObjectLitProperties(
		objectLit, rhino.ObjectType_Cast(objectLit.GetJSType()), !createEnumType)
}

// ProcessObjectLitProperties processes an object literal and all the types
// on it.
//
// objLit is the OBJECTLIT node; objLitType the type of the OBJECTLIT node,
// which might be a named type, because of the lends annotation;
// declareOnOwner, if true, declares properties on the objLitType as well.
// If false, the caller should take care of this.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.processObjectLitProperties(Node, ObjectType, boolean)
func (b *TypedScopeCreator_scopeBuilder) ProcessObjectLitProperties(objLit rhino.Node, objLitType rhino.ObjectType, declareOnOwner bool) {
	lit := ast.MustObjectLit(objLit)
	if !lit.HasMembers() {
		return
	}
	objOwner := NodeUtil_GetBestLValue(objLit)
	objOwnerName, hasObjOwnerName := nodeUtil_getBestLValueNameString(objOwner)
	var objRootScope *TypedScope
	for member := range lit.Members() {
		keyNode := member.Node()
		if keyNode.IsSpread() {
			// Don't try defining computed or spread properties on an object. Note that for spread
			// type inference will try to determine the properties and types. We cannot do it here as
			// we don't have all the type information of the spread object.
			continue
		}
		var memberName rhino.Property_Key
		var value rhino.Node
		var qualifiedName string
		var hasQualifiedName bool
		if cp, ok := member.AsComputedProp(); ok {
			symbol := TypedScopeCreator_ExtractKnownSymbolKey(b.currentScope, cp.Key().Node())
			if symbol == nil {
				continue
			}
			memberName = rhino.NewProperty_SymbolKey(symbol)
			if v, ok := cp.Value(); ok {
				value = v.Node()
			}
		} else {
			value = keyNode.GetFirstChild()
			keyName := NodeUtil_GetObjectOrClassLitKeyName(keyNode)
			memberName = rhino.NewProperty_StringKey(keyName)
			if hasObjOwnerName && rhino.TokenStream_IsJSIdentifier(keyName) {
				qualifiedName = objOwnerName + "." + keyName
				hasQualifiedName = true
			}
		}
		info := keyNode.GetJSDocInfo()
		valueType := b.GetDeclaredType(info, keyNode, value, nil)
		var keyType rhino.JSType
		if objLitType.IsEnumType() {
			keyType = objLitType.ToMaybeEnumType().GetElementsType()
		} else {
			keyType = TypeCheck_GetObjectLitKeyTypeFromValueType(keyNode, valueType)
		}

		// Try to declare this property in the current scope if it
		// has an authoritative name.
		if hasQualifiedName {
			var scopeToDeclareIn *TypedScope
			if keyNode.IsStringKey() {
				if objRootScope == nil {
					objRootScope = b.GetLValueRootScope(objOwner)
				}
				scopeToDeclareIn = objRootScope
			} else {
				scopeToDeclareIn = b.currentHoistScope.GetGlobalScope()
			}
			b.defineSlot(typedScopeSlotSpec{
				declNode:                keyNode,
				varName:                 qualifiedName,
				scope:                   scopeToDeclareIn,
				type_:                   keyType,
				allowLaterTypeInference: keyType == nil,
			})
		} else if keyType != nil {
			b.outer.SetDeferredType(keyNode, keyType)
		}

		if keyType != nil && objLitType != nil && declareOnOwner {
			// Declare this property on its object literal.
			objLitType.DefineDeclaredPropertyByKey(memberName, keyType, keyNode)
		}
	}
}

// GetDeclaredTypeInAnnotation returns the type specified in a JSDoc
// annotation near a GETPROP, NAME, object literal member, or class field.
//
// Extracts type information from the @type tag.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredTypeInAnnotation(Node, JSDocInfo)
func (b *TypedScopeCreator_scopeBuilder) GetDeclaredTypeInAnnotation(node rhino.Node, info *rhino.JSDocInfo) rhino.JSType {
	jutil.CheckArgumentMsg1(info.HasType(), "%s", info)
	ownerTypeKeys := b.FindOwnerTypeKeys(node)

	var templateScope rhino.StaticTypedScope
	if len(ownerTypeKeys) != 0 {
		templateScope = b.outer.typeRegistry.CreateScopeWithTemplates(b.currentScope, ownerTypeKeys)
	} else {
		templateScope = b.currentScope
	}
	return info.GetType().Evaluate(templateScope, b.outer.typeRegistry)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.findOwnerTypeKeys(Node)
func (b *TypedScopeCreator_scopeBuilder) FindOwnerTypeKeys(node rhino.Node) []*rhino.TemplateType {
	if !node.IsNil() && !node.IsStaticMember() && node.GetParent().IsClassMembers() {
		classType := node.GetGrandparent().GetJSType()
		if classType != nil {
			return classType.ToMaybeFunctionType().GetTemplateTypeMap().GetTemplateKeys()
		}
		return emptyTemplateTypeSlice
	}
	ownerNode := NodeUtil_GetBestLValueOwner(node)
	if ownerName, ok := b.getOwnerLValueNameString(ownerNode); ok {
		ownerVar := b.currentScope.GetVar(ownerName)
		if ownerVar != nil {
			ownerType := b.GetPrototypeOwnerType(rhino.ObjectType_Cast(ownerVar.GetType()))
			if ownerType != nil {
				return ownerType.GetTemplateTypeMap().GetTemplateKeys()
			}
		}
	}
	return emptyTemplateTypeSlice
}

// AssertDefinitionNode asserts that it's OK to define this node's name. The
// node should have a source name and be of the specified type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.assertDefinitionNode(Node, Token)
func (b *TypedScopeCreator_scopeBuilder) AssertDefinitionNode(n rhino.Node, type_ rhino.Token) {
	jutil.CheckState(b.sourceName != nil)
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
	jutil.CheckState(b.sourceName != nil)
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
	if b.IsGoogModuleExports(lvalue) {
		return jutil.StringPtr(jutil.StringReplace(*syntacticLvalueName, "exports", *b.GetModule().ClosureNamespace()))
	}
	return syntacticLvalueName
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
			typedScopeCreator_sourceNameArg(b.sourceName), input, d.declNode, parent, oldVar, variableName, d.type_)
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
	if prototypeVar != nil && prototypeVar.GetScope() == rhino.StaticScope(scopeToDeclareIn) {
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
				var_ := b.currentScope.GetVar(root.GetString())
				if var_ != nil {
					return var_.GetScope().(*TypedScope)
				}
			}
		} else if root.IsThis() || root.IsSuper() {
			// We want the enclosing function scope, or the global scope if not in a function.
			return b.currentHoistScope.GetScopeOfThis()
		}
	}
	return b.currentHoistScope.GetGlobalScope()
}
