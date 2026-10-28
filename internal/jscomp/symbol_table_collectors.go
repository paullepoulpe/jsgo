// Port of com/google/javascript/jscomp/SymbolTable.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// SymbolTable_PropertyRefCollector collects property references.
//
// Java: private class SymbolTable.PropertyRefCollector extends AbstractPostOrderCallback implements CompilerPass
type SymbolTable_PropertyRefCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *SymbolTable
}

// Java: SymbolTable.PropertyRefCollector()
func NewSymbolTable_PropertyRefCollector(outer *SymbolTable) *SymbolTable_PropertyRefCollector {
	return &SymbolTable_PropertyRefCollector{outer: outer}
}

// Java: SymbolTable.PropertyRefCollector.process(Node, Node)
func (c *SymbolTable_PropertyRefCollector) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(c.outer.compiler, c, externs, root)
}

// Java: SymbolTable.PropertyRefCollector.maybeDefineReference(Node, String, Symbol)
func (c *SymbolTable_PropertyRefCollector) maybeDefineReference(n rhino.Node, propName string, ownerSymbol *SymbolTable_Symbol) bool {
	if ownerSymbol == nil {
		return false
	}
	// If current symbol is ObjectType - try using ObjectType.getPropertyNode to find a node.
	// That function searches for property in all extended classes and interfaces, which supports
	// multiple extended interfaces.
	owner := ownerSymbol.GetType()
	if owner != nil && owner.IsObjectType() {
		propNode := owner.ToObjectType().GetPropertyNode(propName)
		propSymbol := c.outer.symbols.Get(symbolTableKey{declNode: propNode, name: propName})
		if propSymbol != nil {
			propSymbol.DefineReferenceAt(n)
			return true
		}
	}
	// Use PropertyScope to find property. PropertyScope doesn't support multiple parents so all
	// cases where multiple parents can occur (e.g. class implementing 2+ interfaces) must be
	// handled above.
	if ownerSymbol.GetPropertyScope() != nil {
		prop := ownerSymbol.GetPropertyScope().GetSlot(propName)
		if prop != nil {
			prop.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// tryDefineLexicalQualifiedNameRef tries to find the symbol by its fully
// qualified name.
//
// Java: SymbolTable.PropertyRefCollector.tryDefineLexicalQualifiedNameRef(String, Node)
func (c *SymbolTable_PropertyRefCollector) tryDefineLexicalQualifiedNameRef(name *string, n rhino.Node) bool {
	if name != nil {
		lexicalSym := c.outer.GetEnclosingScope(n).GetQualifiedSlot(*name)
		if lexicalSym != nil {
			lexicalSym.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// tryRemoveLexicalQualifiedNameRef tries to remove a reference by its fully
// qualified name. If the symbol has no references left, remove it completely.
//
// Java: SymbolTable.PropertyRefCollector.tryRemoveLexicalQualifiedNameRef(String, Node)
func (c *SymbolTable_PropertyRefCollector) tryRemoveLexicalQualifiedNameRef(name *string, n rhino.Node) {
	if name != nil {
		lexicalSym := c.outer.GetEnclosingScope(n).GetQualifiedSlot(*name)
		if lexicalSym != nil &&
			lexicalSym.IsLexicalVariable() &&
			lexicalSym.GetDeclaration().GetNode() == n {
			c.outer.removeSymbol(lexicalSym)
		}
	}
}

// Java: SymbolTable.PropertyRefCollector.maybeDefineTypedReference(Node, String, JSType)
func (c *SymbolTable_PropertyRefCollector) maybeDefineTypedReference(n rhino.Node, propName string, owner rhino.JSType) bool {
	t := c.outer
	if owner.IsGlobalThisType() {
		sym := t.globalScope.GetSlot(propName)
		if sym != nil {
			sym.DefineReferenceAt(n)
			return true
		}
	} else if owner.IsNominalConstructorOrInterface() {
		return c.maybeDefineReference(n, propName, t.GetSymbolDeclaredBy(owner.ToMaybeFunctionType()))
	} else if owner.IsEnumType() {
		return c.maybeDefineReference(n, propName, t.GetSymbolDeclaredByEnum(owner.ToMaybeEnumType()))
	} else {
		defined := false
		for _, ctor := range t.GetAllSymbolsForType(owner) {
			if ctor.GetPropertyScope() != nil {
				ctorType := ctor.GetType()
				isConstructor :=
					ctorType != nil && (ctorType.IsConstructor() || ctorType.IsInterface())
				if !isConstructor && c.maybeDefineReference(n, propName, ctor) {
					defined = true
				}
			}
			if c.maybeDefineReference(n, propName, t.GetSymbolForInstancesOf(ctor)) {
				defined = true
			}
		}
		if defined {
			return true
		}
	}
	// Handle cases like accessing properties on a imported module.
	// const foo = goog.require('some.module.foo');
	// foo.doSomething();
	// here n is foo, propName is doSomething and owner is module type.
	if objType := owner.ToMaybeObjectType(); objType != nil {
		propNodeDecl := objType.GetPropertyNode(propName)
		if propNodeDecl.IsNil() {
			return false
		}
		// There is no handy way to find symbol object the given property node. So we do
		// property node => namespace node => namespace symbol => property symbol.
		if gp, ok := ast.AsNormalOrOptChainGetProp(propNodeDecl); ok {
			namespace := gp.Object().Node()
			namespaceSym := t.GetSymbolForName(namespace, namespace.GetQualifiedName())
			if namespaceSym != nil && namespaceSym.GetPropertyScope() != nil {
				propSym := namespaceSym.GetPropertyScope().GetSlot(propName)
				if propSym != nil && propSym.GetDeclarationNode() != n {
					propSym.DefineReferenceAt(n)
					return true
				}
			}
		}

		// Here handle the case of goog.module exports that have shape of:
		//
		// goog.module('foo.bar');
		// function doOne() { ... }
		// exports = {doOne};
		//
		// propNodeDecl is `doOne` node on the last line with exports. From that node and propName
		// `doOne` we need to get symbol for the `function doOne`. To achieve that we get scope of
		// `doOne` node which should be whole module scope and look for symbol with corresponding
		// name. This relies on the assumption that function name in module scope exported name. But
		// if export is renamed: `exports = {doOneRenamed: doOne}`, then this approach will fail.
		propNodeDeclScope := t.GetEnclosingScope(propNodeDecl)
		if propNodeDeclScope == nil || !propNodeDeclScope.IsModuleScope() {
			return false
		}
		symbol := t.GetSymbolForName(propNodeDecl, &propName)
		if symbol != nil && symbol.GetDeclarationNode() != n {
			symbol.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// Java: SymbolTable.PropertyRefCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_PropertyRefCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// There are two ways to define a property reference:
	// 1) As a fully qualified lexical symbol (e.g., x.y)
	// 2) As a property of another object (e.g., x's y)
	// Property definitions should take precedence over lexical
	// definitions. e.g., for "a.b", it's more useful to record
	// this as "property b of the type of a", than as "symbol a.b".

	if gp, ok := ast.AsNormalOrOptChainGetProp(n); ok {
		owner := gp.Object().Node().GetJSType()
		if owner != nil {
			defined := c.maybeDefineTypedReference(n, gp.Property(), owner)

			if defined {
				c.tryRemoveLexicalQualifiedNameRef(n.GetQualifiedName(), n)
				return
			}
		}

		c.tryDefineLexicalQualifiedNameRef(n.GetQualifiedName(), n)
	} else if sk, ok := ast.AsStringKey(n); ok {
		owner := parent.GetJSType()
		if owner != nil {
			defined := c.maybeDefineTypedReference(n, sk.Key(), owner)

			if defined {
				c.tryRemoveLexicalQualifiedNameRef(NodeUtil_GetBestLValueName(n), n)
				return
			}
		}

		c.tryDefineLexicalQualifiedNameRef(NodeUtil_GetBestLValueName(n), n)
	}
}

// SymbolTable_ThisRefCollector collects references to "this".
//
// Java: private class SymbolTable.ThisRefCollector extends NodeTraversal.AbstractScopedCallback implements CompilerPass
type SymbolTable_ThisRefCollector struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer *SymbolTable

	// The 'this' symbols in the current scope chain.
	//
	// If we don't know how to declare 'this' in a scope chain,
	// then null should be on the stack. But this should be a rare
	// occurrence. We should strive to always be able to come up
	// with some symbol for 'this'.
	//
	// This list only has entries for function scopes and the global scope, and doesn't store a
	// separate `this` value for other block scopes.
	thisStack []*SymbolTable_Symbol
}

// Java: SymbolTable.ThisRefCollector()
func NewSymbolTable_ThisRefCollector(outer *SymbolTable) *SymbolTable_ThisRefCollector {
	return &SymbolTable_ThisRefCollector{outer: outer}
}

// Java: SymbolTable.ThisRefCollector.process(Node, Node)
func (c *SymbolTable_ThisRefCollector) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(c.outer.compiler, c, externs, root)
}

// Java: SymbolTable.ThisRefCollector.enterScope(NodeTraversal)
func (c *SymbolTable_ThisRefCollector) EnterScope(t *NodeTraversal) {
	st := c.outer
	var symbol *SymbolTable_Symbol
	if t.InGlobalScope() {
		// Declare the global this at the first input root.
		// This is a bizarre place to put it, but we need some
		// location with a real file path (because all symbols
		// must have a path).
		// Note that root.lastChild.firstChild is the first non-extern input.
		firstInputRoot := t.GetScopeRoot().GetLastChild().GetFirstChild()
		if !firstInputRoot.IsNil() {
			symbol =
				st.addSymbol(
					SymbolTable_GLOBAL_THIS,
					st.registry.GetNativeType(rhino.JSTypeNative_GLOBAL_THIS),
					false, /* declared */
					st.globalScope,
					firstInputRoot)
			symbol.SetDeclaration(NewSymbolTable_Reference(symbol, firstInputRoot))
		}
		c.thisStack = append(c.thisStack, symbol)
	} else if fn, ok := ast.AsFunction(t.GetScopeRoot()); ok {
		// Otherwise, declare a "this" property when possible.
		scopeRoot := fn.Node()
		scope := st.scopes.Get(scopeRoot)
		if fn.Body().Node().HasChildren() {
			scopeSymbol := st.GetSymbolForScope(scope)
			if scopeSymbol != nil {
				propScope := scopeSymbol.GetPropertyScope()
				if propScope != nil {
					// If a function is assigned multiple times, we only want
					// one addressable "this" symbol.
					symbol = propScope.GetOwnSlot("this")
					if symbol == nil {
						rootType := scopeRoot.GetJSType()
						var fnType rhino.FunctionType
						if rootType != nil {
							fnType = rootType.ToMaybeFunctionType()
						}
						var type_ rhino.JSType
						if fnType != nil {
							type_ = fnType.GetTypeOfThis()
						}
						symbol = st.addSymbol("this", type_ /* inferred= */, false, scope, scopeRoot)
					}
				}
			}
		} else {
			symbolTable_logger.Fine("Skipping empty function: " + scopeRoot.String())
		}
		c.thisStack = append(c.thisStack, symbol)
	}
	// Don't add to the `thisStack` for other block scopes.
}

// Java: SymbolTable.ThisRefCollector.exitScope(NodeTraversal)
func (c *SymbolTable_ThisRefCollector) ExitScope(t *NodeTraversal) {
	if t.InGlobalScope() || t.GetScopeRoot().IsFunction() {
		c.thisStack = c.thisStack[:len(c.thisStack)-1]
	}
}

// Java: SymbolTable.ThisRefCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_ThisRefCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsThis() {
		return
	}

	if len(c.thisStack) == 0 {
		// Java: Iterables.getLast of an empty list.
		panic(jutil.NewNoSuchElementException())
	}
	symbol := c.thisStack[len(c.thisStack)-1]
	if symbol != nil {
		ref := symbol.DefineReferenceAt(n)
		if symbol.GetDeclaration() == nil {
			symbol.SetDeclaration(ref)
		}
	}
}

// SymbolTable_JSDocInfoCollector collects references to types in JSDocInfo.
//
// Java: private class SymbolTable.JSDocInfoCollector extends AbstractPostOrderCallback
type SymbolTable_JSDocInfoCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer        *SymbolTable
	typeRegistry *rhino.JSTypeRegistry
}

// Java: SymbolTable.JSDocInfoCollector(JSTypeRegistry)
func NewSymbolTable_JSDocInfoCollector(outer *SymbolTable, registry *rhino.JSTypeRegistry) *SymbolTable_JSDocInfoCollector {
	return &SymbolTable_JSDocInfoCollector{outer: outer, typeRegistry: registry}
}

// Java: SymbolTable.JSDocInfoCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_JSDocInfoCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	info := n.GetJSDocInfo()
	if info == nil {
		return
	}
	c.outer.docInfos = append(c.outer.docInfos, n)

	var scope *SymbolTable_SymbolScope
	for it := info.TypeExprIter(); it.Next(); {
		if scope == nil {
			scope = c.outer.scopes.Get(t.GetScopeRoot())
			if scope == nil {
				scope = c.outer.globalScope
			}
		}
		c.visitTypeNode(info, scope, it.Node())
	}
}

// Java: SymbolTable.JSDocInfoCollector.isNativeSourcelessType(String)
func (c *SymbolTable_JSDocInfoCollector) isNativeSourcelessType(name string) bool {
	switch name {
	case "null", "undefined", "void":
		return true
	default:
		return false
	}
}

// Java: SymbolTable.JSDocInfoCollector.visitTypeNode(ImmutableMap<String, JSTypeExpression>, SymbolScope, Node)
func (c *SymbolTable_JSDocInfoCollector) visitTypeNode(
	info *rhino.JSDocInfo, scope *SymbolTable_SymbolScope, n rhino.Node) {
	if n.IsStringLit() {
		name := n.GetString()
		if !c.isNativeSourcelessType(name) && !info.HasTemplateTypeName(name) {
			symbol := c.lookupPossiblyDottedName(scope, name)
			if symbol != nil {
				ref := n
				typeString := name
				if orig := n.GetOriginalName(); orig != nil {
					typeString = *orig
				}
				// Qualified names in JSDoc types are kept as a single string: "foo.bar.MyType". In order
				// to have good indexing we need to make the SymbolTable reference include only "MyType"
				// instead of "foo.bar.MyType". To do that we clone the type node and change the source
				// info of the clone to include only the last part of the type ("MyType").
				if dot := strings.LastIndexByte(typeString, '.'); dot >= 0 {
					lastPart := typeString[dot+1:]
					lastPartLen := jutil.StringLength(lastPart)
					copy_ := n.CloneNode()
					copy_.SetLinenoCharno(
						copy_.GetLineno(), copy_.GetCharno()+copy_.GetLength()-lastPartLen)
					copy_.SetLength(lastPartLen)
					ref = copy_
				}
				symbol.DefineReferenceAt(ref)
			}
		}
	}

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		c.visitTypeNode(info, scope, child)
	}
}

// Java: SymbolTable.JSDocInfoCollector.lookupPossiblyDottedName(SymbolScope, String)
func (c *SymbolTable_JSDocInfoCollector) lookupPossiblyDottedName(scope *SymbolTable_SymbolScope, dottedName string) *SymbolTable_Symbol {
	// Try the dotted name to start.
	result := scope.GetQualifiedSlot(dottedName)
	if result != nil {
		return result
	}
	// If we can't find this type, it might be a reference to a
	// primitive type (like {string}). Autobox it to check.
	type_ := c.typeRegistry.GetGlobalType(dottedName)
	var autobox rhino.JSType
	if type_ != nil {
		autobox = type_.AutoboxesTo()
	}
	if autobox == nil {
		return c.outer.GetSymbolForTypeHelper(type_, true)
	}
	return c.outer.GetSymbolForTypeHelper(autobox, true)
}

// SymbolTable_VisibilityCollector collects the visibility information for
// each name/property.
//
// Java: private class SymbolTable.VisibilityCollector extends AbstractPostOrderCallback
type SymbolTable_VisibilityCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer             *SymbolTable
	fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]
}

// Java: SymbolTable.VisibilityCollector(ImmutableMap<StaticSourceFile, Visibility>)
func NewSymbolTable_VisibilityCollector(
	outer *SymbolTable, fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]) *SymbolTable_VisibilityCollector {
	return &SymbolTable_VisibilityCollector{outer: outer, fileVisibilityMap: fileVisibilityMap}
}

// Java: SymbolTable.VisibilityCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_VisibilityCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if nameNode, ok := ast.AsName(n); ok {
		c.visitName(t, nameNode)
	} else if gp, ok := ast.AsNormalOrOptChainGetProp(n); ok {
		c.visitProperty(gp, parent)
	}
}

// Java: SymbolTable.VisibilityCollector.visitName(NodeTraversal, Node)
func (c *SymbolTable_VisibilityCollector) visitName(t *NodeTraversal, nameNode ast.Name) {
	n := nameNode.Node()
	ident := nameNode.Ident()
	symbol := c.outer.symbols.Get(symbolTableKey{declNode: n, name: ident})
	if symbol == nil {
		return
	}
	// Visibility already set.
	if symbol.GetVisibility() != nil {
		return
	}
	var_ := t.GetScope().GetVar(ident)
	if var_ == nil {
		return
	}
	// Go-only: AccessControlUtils.getEffectiveNameVisibility never returns
	// null here (see AccessControlUtils), so Java's null check is dropped.
	v := AccessControlUtils_GetEffectiveNameVisibility(n, var_, c.fileVisibilityMap)
	symbol.SetVisibility(v)
}

// Java: SymbolTable.VisibilityCollector.visitProperty(Node, Node)
func (c *SymbolTable_VisibilityCollector) visitProperty(gp ast.NormalOrOptChainGetProp, parent rhino.Node) {
	getprop := gp.Node()
	propertyName := gp.Property()
	symbol := c.outer.symbols.Get(symbolTableKey{declNode: getprop, name: propertyName})
	if symbol == nil {
		return
	}
	// Visibility already set.
	if symbol.GetVisibility() != nil {
		return
	}
	jsType := gp.Object().Node().GetJSType()
	if jsType == nil {
		return
	}
	assign, isAssign := ast.AsAssign(parent)
	isOverride := isAssign && parent.GetJSDocInfo() != nil && assign.Target().Node() == getprop
	if isOverride {
		// Don't bother with AccessControlUtils for overridden properties.
		// AccessControlUtils currently has complicated logic for detecting
		// visibility mismatches for overridden properties that is still
		// too tightly coupled to CheckAccessControls. TODO(brndn): simplify.
		symbol.SetVisibility(rhino.JSDocInfo_Visibility_INHERITED)
	} else {
		referenceType := rhino.ObjectType_Cast(jsType.Dereference())
		// Go-only: AccessControlUtils.getEffectivePropertyVisibility never
		// returns null (see AccessControlUtils), so Java's null check is dropped.
		v :=
			AccessControlUtils_GetEffectivePropertyVisibility(
				getprop, referenceType, c.fileVisibilityMap)
		symbol.SetVisibility(v)
	}
}
