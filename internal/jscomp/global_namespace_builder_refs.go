// Port of com/google/javascript/jscomp/GlobalNamespace.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// GetValueType gets the type of a value or simple expression.
//
// n is an r-value in an assignment or variable declaration (not null).
//
// Java: GlobalNamespace.BuildGlobalNamespace.getValueType(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) GetValueType(n rhino.Node) GlobalNamespace_NameProp {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		return GlobalNamespace_NameProp_CLASS
	case rhino.Token_OBJECTLIT:
		return GlobalNamespace_NameProp_OBJECTLIT
	case rhino.Token_FUNCTION:
		return GlobalNamespace_NameProp_FUNCTION
	case rhino.Token_OR:
		// Recurse on the second value. If the first value were an object
		// literal or function, then the OR would be meaningless and the
		// second value would be dead code. Assume that if the second value
		// is an object literal or function, then the first value will also
		// evaluate to one when it doesn't evaluate to false.
		return b.GetValueType(ast.MustBinary(n).Right().Node())
	case rhino.Token_HOOK:
		// The same line of reasoning used for the OR case applies here.
		hook := ast.MustHook(n)
		t := b.GetValueType(hook.Then().Node())
		if t != GlobalNamespace_NameProp_OTHER_OBJECT {
			return t
		}
		return b.GetValueType(hook.Else().Node())
	default:
	}
	return GlobalNamespace_NameProp_OTHER_OBJECT
}

// HandleSetFromGlobal updates our representation of the global namespace to
// reflect an assignment to a global name in any scope where variables are
// hoisted to the global scope (i.e. the global scope in an ES5 sense).
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d"); type the type of the value that the
// name is being assigned.
//
// Java: GlobalNamespace.BuildGlobalNamespace.handleSetFromGlobal(JSChunk, Scope, Node, String, NameProp, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) HandleSetFromGlobal(
	chunk *JSChunk, scope *Scope, n rhino.Node, name string, type_ GlobalNamespace_NameProp, metadata *ModuleMetadataMap_ModuleMetadata) {
	if b.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) || b.skipsName(name) {
		return
	}
	b.handleResolvedSetFromGlobal(chunk, scope, n, b.GetOrCreateName(name, metadata), type_)
}

func (b *GlobalNamespace_BuildGlobalNamespace) handleResolvedSetFromGlobal(
	chunk *JSChunk, scope *Scope, n rhino.Node, nameObj *GlobalNamespace_Name, type_ GlobalNamespace_NameProp) {
	if !nameObj.IsGetOrSetDefinition() {
		// Don't change the type of a getter or setter. This is because given:
		//   var a = {set b(item) {}}; a.b = class {};
		// `a.b = class {};` does not change the runtime value of a.b, and we do not want to change
		// the 'type' of a.b to Type.CLASS.
		// TODO(lharker): for non-setter cases, do we really want to just treat the last set of
		// a name as canonical? e.g. what if a name is first set to a class, then an object literal?
		nameObj.SetNameType(type_)
	}
	if n.GetBooleanProp(rhino.Node_MODULE_EXPORT) {
		nameObj.SetBooleanProperty(GlobalNamespace_NameProp_IS_MODULE_PROP)
	}

	if b.IsNestedAssign(n.GetParent()) {
		// This assignment is both a set and a get that creates an alias.
		refType := GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL
		b.AddOrConfirmRef(nameObj, n, refType, scope, chunk)
	} else {
		b.AddOrConfirmRef(nameObj, n, GlobalNamespace_Ref_Type_SET_FROM_GLOBAL, scope, chunk)
		nameObj.SetDeclaredTypeKind(b.GetDeclaredTypeKind(n))
	}
}

// GetDeclaredTypeKind determines whether a set operation is a constructor or
// enumeration or interface declaration. The set operation may either be an
// assignment to a name, a variable declaration, or an object literal key
// mapping.
//
// n is the node that represents the name being set.
//
// Java: GlobalNamespace.BuildGlobalNamespace.getDeclaredTypeKind(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) GetDeclaredTypeKind(n rhino.Node) GlobalNamespace_NameProp {
	valueNode := NodeUtil_GetRValueOfLValue(n)
	var kind GlobalNamespace_NameProp
	if valueNode.IsNil() {
		kind = GlobalNamespace_NameProp_NOT_A_TYPE
	} else if valueNode.IsClass() {
		// Always treat classes as having a declared type. (Transpiled classes are annotated
		// @constructor)
		kind = GlobalNamespace_NameProp_CONSTRUCTOR_TYPE
	} else {
		info := NodeUtil_GetBestJSDocInfo(n)
		// Heed the annotations only if they're sensibly used.
		if info == nil {
			kind = GlobalNamespace_NameProp_NOT_A_TYPE
		} else if info.IsConstructor() && valueNode.IsFunction() {
			kind = GlobalNamespace_NameProp_CONSTRUCTOR_TYPE
		} else if info.IsInterface() && valueNode.IsFunction() {
			kind = GlobalNamespace_NameProp_INTERFACE_TYPE
		} else if info.HasEnumParameterType() && valueNode.IsObjectLit() {
			kind = GlobalNamespace_NameProp_ENUM_TYPE
		} else {
			kind = GlobalNamespace_NameProp_NOT_A_TYPE
		}
	}
	return kind
}

// HandleSetFromLocal updates our representation of the global namespace to
// reflect an assignment to a global name in a local scope.
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d").
//
// Java: GlobalNamespace.BuildGlobalNamespace.handleSetFromLocal(JSChunk, Scope, Node, String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) HandleSetFromLocal(
	chunk *JSChunk, scope *Scope, n rhino.Node, name string, metadata *ModuleMetadataMap_ModuleMetadata) {
	if b.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) || b.skipsName(name) {
		return
	}
	b.handleResolvedSetFromLocal(chunk, scope, n, b.GetOrCreateName(name, metadata))
}

func (b *GlobalNamespace_BuildGlobalNamespace) handleResolvedSetFromLocal(
	chunk *JSChunk, scope *Scope, n rhino.Node, nameObj *GlobalNamespace_Name) {
	if n.GetBooleanProp(rhino.Node_MODULE_EXPORT) {
		nameObj.SetBooleanProperty(GlobalNamespace_NameProp_IS_MODULE_PROP)
	}

	if b.IsNestedAssign(n.GetParent()) {
		// This assignment is both a set and a get that creates an alias.
		b.AddOrConfirmRef(nameObj, n, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL, scope, chunk)
	} else {
		b.AddOrConfirmRef(nameObj, n, GlobalNamespace_Ref_Type_SET_FROM_LOCAL, scope, chunk)
	}
}

// HandleGet updates our representation of the global namespace to reflect a
// read of a global name.
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d").
//
// Java: GlobalNamespace.BuildGlobalNamespace.handleGet(JSChunk, Scope, Node, String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) HandleGet(chunk *JSChunk, scope *Scope, n rhino.Node, name string, metadata *ModuleMetadataMap_ModuleMetadata) {
	if b.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) || b.skipsName(name) {
		return
	}
	type_ := b.DetermineRefTypeForGet(n, n, name)

	b.AddOrConfirmRef(b.GetOrCreateName(name, metadata), n, type_, scope, chunk)
}

func (b *GlobalNamespace_BuildGlobalNamespace) handleResolvedGet(chunk *JSChunk, scope *Scope, n rhino.Node, nameObj *GlobalNamespace_Name) {
	type_ := b.determineRefTypeForGetNode(n, n, nameObj)
	b.AddOrConfirmRef(nameObj, n, type_, scope, chunk)
}

// DetermineRefTypeForGet determines the Ref.Type for referenceNode by
// inspecting parents and recusively ascending as necessary, where n
// represents.
//
// n is the node currently being visited; referenceNode the node hosting the
// name; name the global name (e.g. "a" or "a.b.c.d").
//
// Java: GlobalNamespace.BuildGlobalNamespace.determineRefTypeForGet(Node, Node, String)
func (b *GlobalNamespace_BuildGlobalNamespace) DetermineRefTypeForGet(n rhino.Node, referenceNode rhino.Node, name string) GlobalNamespace_Ref_Type {
	return b.determineRefTypeForGetImpl(n, referenceNode, name, nil)
}

func (b *GlobalNamespace_BuildGlobalNamespace) determineRefTypeForGetNode(n rhino.Node, referenceNode rhino.Node, nameObj *GlobalNamespace_Name) GlobalNamespace_Ref_Type {
	return b.determineRefTypeForGetImpl(n, referenceNode, "", nameObj)
}

func (b *GlobalNamespace_BuildGlobalNamespace) determineRefTypeForGetImpl(
	n rhino.Node, referenceNode rhino.Node, name string, nameObj *GlobalNamespace_Name) GlobalNamespace_Ref_Type {
	var type_ GlobalNamespace_Ref_Type
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_EXPR_RESULT,
		rhino.Token_IF,
		rhino.Token_WHILE,
		rhino.Token_FOR,
		rhino.Token_INSTANCEOF,
		rhino.Token_TYPEOF,
		rhino.Token_VOID,
		rhino.Token_NOT,
		rhino.Token_BITNOT,
		rhino.Token_POS,
		rhino.Token_NEG,
		rhino.Token_SHEQ,
		rhino.Token_EQ,
		rhino.Token_SHNE,
		rhino.Token_NE,
		rhino.Token_LT,
		rhino.Token_LE,
		rhino.Token_GT,
		rhino.Token_GE,
		rhino.Token_ADD,
		rhino.Token_SUB,
		rhino.Token_MUL,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_EXPONENT,
		rhino.Token_BITAND,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_LSH,
		rhino.Token_RSH,
		rhino.Token_URSH:
		type_ = GlobalNamespace_Ref_Type_DIRECT_GET
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if r := ast.RoleIn(n, parent); r == ast.RoleCallCallee || r == ast.RoleOptChainCallCallee {
			// It is a call target
			type_ = GlobalNamespace_Ref_Type_CALL_GET
		} else if b.IsClassDefiningCall(parent) {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_NEW:
		if ast.RoleIn(n, parent) == ast.RoleNewCallee {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_CAST, rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		// This node is x or y in (x||y), (x&&y), or (x??y). We only know that an
		// alias is not getting created for this name if the result is used
		// in a boolean context or assigned to the same name
		// (e.g. var a = a || {}).
		type_ = b.determineRefTypeForGetImpl(parent, referenceNode, name, nameObj)
	case rhino.Token_NAME:
		// Only LET, CONST, VAR declarations have NAME nodes
		// with children.
		// Of particular interest is "var n = n || {}"
		sameName := false
		if n != referenceNode {
			if nameObj != nil {
				sameName = nameObj.parentID == 0 && nameObj.baseNameID == parent.GetStringID()
			} else {
				sameName = name == parent.GetString()
			}
		}
		if sameName {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_COMMA, rhino.Token_HOOK:
		switch ast.RoleIn(n, parent) {
		case ast.RoleBinaryLeft, ast.RoleHookCond:
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		default:
			// This node is y or z in (x?y:z) or (x,y). We only know that an alias is
			// not getting created for this name if the result is assigned to
			// the same name (e.g. var a = a ? a : {}).
			type_ = b.determineRefTypeForGetImpl(parent, referenceNode, name, nameObj)
		}
	case rhino.Token_DELPROP:
		type_ = GlobalNamespace_Ref_Type_DELETE_PROP
	case rhino.Token_CLASS:
		// This node is the superclass in an extends clause.
		type_ = GlobalNamespace_Ref_Type_SUBCLASSING_GET
	case rhino.Token_DESTRUCTURING_LHS, rhino.Token_ASSIGN:
		lhs := n.GetPrevious()
		if lhs.IsNil() {
			// TODO(b/127505242): CAST confused the "is this a get or set?"
			// logic and "handleGet" should not have been called.
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
			break
		}
		for cast, ok := ast.AsCast(lhs); ok; cast, ok = ast.AsCast(lhs) {
			// Case: `/** @type {!Foo} */ (x) = ...`; or multiple casts like `(cast(cast(x)) =`
			lhs = cast.Expr().Node()
		}

		// This is a recursive ascent check if this an assignment
		// to itself.  This handles cases like: "a.b = a.b || {}"
		if n != referenceNode {
			if nameObj != nil {
				if nameObj.MatchesNode(lhs) {
					return GlobalNamespace_Ref_Type_DIRECT_GET
				}
			} else if lhs.MatchesQualifiedName(name) {
				return GlobalNamespace_Ref_Type_DIRECT_GET
			}
		}

		switch lhs.GetToken() {
		case rhino.Token_NAME,
			rhino.Token_GETPROP,
			rhino.Token_GETELEM,
			// The rhs of an assign or a name declaration is escaped if it's assigned to a name
			// directly ...
			rhino.Token_ARRAY_PATTERN,
			rhino.Token_OBJECT_PATTERN:
			// ... or referenced through numeric/object keys.
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected previous sibling of " + n.GetToken().String() + ": " + jutil.ToJavaString(n.GetPrevious())))
		}
	case rhino.Token_OBJECT_PATTERN, // Handle STRING_KEYS in object patterns.
		rhino.Token_ITER_SPREAD,
		rhino.Token_OBJECT_SPREAD,
		rhino.Token_RETURN,
		rhino.Token_THROW:
		// NOTE: There are likely more cases where we should be returning
		// DIRECT_GET.
		type_ = GlobalNamespace_Ref_Type_ALIASING_GET
	default:
		// NOTE: There are likely more cases where we should be returning
		// DIRECT_GET.
		type_ = GlobalNamespace_Ref_Type_ALIASING_GET
	}
	return type_
}

// AddOrConfirmRef: if there is already a Ref for the given name & node,
// confirm it matches what we would create. Otherwise add a new one.
//
// Java: GlobalNamespace.BuildGlobalNamespace.addOrConfirmRef(Name, Node, Ref.Type, Scope, JSChunk)
func (b *GlobalNamespace_BuildGlobalNamespace) AddOrConfirmRef(
	nameObj *GlobalNamespace_Name, node rhino.Node, refType GlobalNamespace_Ref_Type, scope *Scope, chunk *JSChunk) {
	var existingRef *GlobalNamespace_Ref
	if b.confirmExisting {
		existingRef = nameObj.GetRefForNode(node)
	}
	if existingRef == nil {
		refID, ref := nameObj.createNewRefWithID(chunk, scope, node, refType)
		nameObj.putNewRefID(node, refID)
		nameObj.updateStateForAddedRef(refID, ref)
	} else {
		// module and scope are dependent on Node, so not much point in checking them
		existingRefType := existingRef.type_
		jutil.CheckStateMsg2(
			existingRefType == refType,
			"existing ref type: %s expected: %s",
			existingRefType,
			refType)
	}
}

// Java: GlobalNamespace.BuildGlobalNamespace.isClassDefiningCall(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsClassDefiningCall(callNode rhino.Node) bool {
	convention := b.outer.compiler.GetCodingConvention()
	// Look for goog.inherits and J2CL mixin calls
	classes := convention.GetClassesDefinedByCall(callNode)
	if classes != nil {
		return true
	}

	// Look for calls to goog.addSingletonGetter calls.
	className := convention.GetSingletonGetterClassName(callNode)
	return className != nil
}

// IsObjectHasOwnPropertyCall detects calls of the form
// a.b.hasOwnProperty(c); that prevent property collapsing on a.b
//
// Java: GlobalNamespace.BuildGlobalNamespace.isObjectHasOwnPropertyCall(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsObjectHasOwnPropertyCall(callNode rhino.Node) bool {
	call := ast.MustNormalOrOptChainCall(callNode)
	if _, ok := call.OnlyArg(); !ok {
		return false
	}
	callee, ok := call.Callee().AsNormalOrOptChainGetProp()
	if !ok {
		return false
	}
	propID := callee.Node().GetStringID()
	return (propID == gnHasOwnPropertyStrID || propID == gnPropertyIsEnumerableStrID) &&
		callee.Object().Node().IsQualifiedName()
}

// MaybeHandlePrototypePrefix updates our representation of the global
// namespace to reflect a read of a global name's longest prefix before the
// "prototype" property if the name includes the "prototype" property. Does
// nothing otherwise.
//
// scope is the current scope; n the node currently being visited; name the
// global name (e.g. "a" or "a.b.c.d"). Returns whether the name was handled.
//
// Java: GlobalNamespace.BuildGlobalNamespace.maybeHandlePrototypePrefix(JSChunk, Scope, Node, String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) MaybeHandlePrototypePrefix(
	chunk *JSChunk, scope *Scope, n rhino.Node, name string, metadata *ModuleMetadataMap_ModuleMetadata) bool {
	// We use a string-based approach instead of inspecting the parse tree
	// to avoid complexities with object literals, possibly nested, beneath
	// assignments.

	var numLevelsToRemove int
	var prefix string
	if strings.HasSuffix(name, ".prototype") {
		numLevelsToRemove = 1
		prefix = name[:len(name)-10]
	} else {
		i := strings.Index(name, ".prototype.")
		if i == -1 {
			return false
		}
		prefix = name[:i]
		numLevelsToRemove = 2
		for j := i + 11; j < len(name); j++ {
			if name[j] == '.' {
				numLevelsToRemove++
			}
		}
	}

	if NodeUtil_MayBeObjectLitKey(n) {
		// Object literal keys have no prefix that's referenced directly per
		// key, so we're done.
		return true
	}

	if b.skipsName(prefix) {
		return true
	}
	for i := 0; i < numLevelsToRemove; i++ {
		n = n.GetFirstChild()
	}

	b.AddOrConfirmRef(b.GetOrCreateName(prefix, metadata), n, GlobalNamespace_Ref_Type_PROTOTYPE_GET, scope, chunk)
	return true
}

// IsNestedAssign determines whether an assignment is nested (i.e. whether
// its return value is used).
//
// parent is the parent of the current traversal node (not null). Returns
// whether it appears that the return value of the assignment is used.
//
// Java: GlobalNamespace.BuildGlobalNamespace.isNestedAssign(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsNestedAssign(parent rhino.Node) bool {
	return parent.IsAssign() && !parent.GetParent().IsExprResult()
}

func (b *GlobalNamespace_BuildGlobalNamespace) getOrCreateRootByID(baseNameID uint32, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	g := b.outer
	if metadata == nil {
		if id := g.lookupChildID(0, baseNameID); id != 0 {
			return g.nameByID(id)
		}
		node := g.allocName(baseNameID, 0, g.sourceKind)
		node.propertyBitSet |= gnNamePropIndexedGlobal
		g.indexedGlobalCount++
		g.globalNameIDs = append(g.globalNameIDs, node.selfID)
		g.putChildID(0, baseNameID, node.selfID)
		return node
	}
	if g.nameMapByModule != nil {
		if id := g.nameMapByModule[globalNamespaceModuleKey{module: metadata, nameID: baseNameID}]; id != 0 {
			return g.nameByID(id)
		}
	} else {
		g.nameMapByModule = make(map[globalNamespaceModuleKey]uint32)
	}
	node := g.allocName(baseNameID, 0, g.sourceKind)
	g.nameMapByModule[globalNamespaceModuleKey{module: metadata, nameID: baseNameID}] = node.selfID
	return node
}

func (b *GlobalNamespace_BuildGlobalNamespace) getOrCreateChildByID(parent *GlobalNamespace_Name, baseNameID uint32, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	g := b.outer
	if id := g.lookupChildID(parent.selfID, baseNameID); id != 0 {
		return g.nameByID(id)
	}
	node := parent.addPropertyWithBaseNameID(baseNameID, g.sourceKind)
	if metadata == nil {
		node.propertyBitSet |= gnNamePropIndexedGlobal
		g.indexedGlobalCount++
	}
	g.putChildID(parent.selfID, baseNameID, node.selfID)
	return node
}

func (b *GlobalNamespace_BuildGlobalNamespace) resolveChainByID(topID uint32, revSegs []uint32, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	cur := b.getOrCreateRootByID(topID, metadata)
	for i := len(revSegs) - 1; i >= 0; i-- {
		cur = b.getOrCreateChildByID(cur, revSegs[i], metadata)
	}
	return cur
}

// GetOrCreateName gets a Name instance for a global name. Creates it if
// necessary, as well as instances for any of its prefixes that are not yet
// defined.
//
// name is a global name (e.g. "a", "a.b.c.d"). Returns the Name instance
// for name.
//
// Java: GlobalNamespace.BuildGlobalNamespace.getOrCreateName(String, ModuleMetadata)
func (b *GlobalNamespace_BuildGlobalNamespace) GetOrCreateName(name string, metadata *ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	if b.lastNode != nil && name == b.lastName && metadata == b.lastMeta {
		return b.lastNode
	}
	if b.lastParentNode != nil && name == b.lastParentName && metadata == b.lastParentMeta {
		return b.lastParentNode
	}
	var node *GlobalNamespace_Name
	i := strings.LastIndexByte(name, '.')
	if i >= 0 {
		parentName := name[:i]
		baseSeg := name[i+1:]
		parent := b.GetOrCreateName(parentName, metadata)
		b.lastParentName = parentName
		b.lastParentMeta = metadata
		b.lastParentNode = parent
		node = b.getOrCreateChildByID(parent, uint32(strpool.Intern(baseSeg)), metadata)
	} else {
		node = b.getOrCreateRootByID(uint32(strpool.Intern(name)), metadata)
	}
	b.lastName = name
	b.lastMeta = metadata
	b.lastNode = node
	return node
}

// -------------------------------------------------------------------------

// Java: GlobalNamespace.createNameForTesting(String)
func (g *GlobalNamespace) CreateNameForTesting(name string) *GlobalNamespace_Name {
	return NewGlobalNamespace_Name(g, name, nil, GlobalNamespace_SourceKind_CODE)
}
