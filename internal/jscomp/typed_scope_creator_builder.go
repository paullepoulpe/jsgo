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

	// Whether we have seen a source file for the current scope.
	hasSourceName bool

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
	b.hasSourceName = false
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
		b.hasSourceName = NodeUtil_GetSourceFile(n) != nil
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
