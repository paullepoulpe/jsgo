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

// -------------------------------------------------------------------------

type gnObjLitState uint8

const (
	gnObjLitInvalid gnObjLitState = iota
	gnObjLitLValue
	gnObjLitDefineProperties
)

type globalNamespace_objLitCacheEntry struct {
	prefix string
	root   rhino.Node
	scope  *Scope
	state  gnObjLitState
}

// GlobalNamespace_BuildGlobalNamespace builds a tree representation of the
// global namespace. Omits prototypes.
type GlobalNamespace_BuildGlobalNamespace struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *GlobalNamespace

	curModuleRoot rhino.Node
	curMetadata   *ModuleMetadataMap_ModuleMetadata

	lastName       string
	lastMeta       *ModuleMetadataMap_ModuleMetadata
	lastNode       *GlobalNamespace_Name
	lastParentName string
	lastParentMeta *ModuleMetadataMap_ModuleMetadata
	lastParentNode *GlobalNamespace_Name

	objLitCache     map[rhino.Node]globalNamespace_objLitCacheEntry
	confirmExisting bool
}

// Java: GlobalNamespace.BuildGlobalNamespace()
func NewGlobalNamespace_BuildGlobalNamespace(outer *GlobalNamespace) *GlobalNamespace_BuildGlobalNamespace {
	return &GlobalNamespace_BuildGlobalNamespace{outer: outer, curModuleRoot: rhino.NilNode, curMetadata: nil, confirmExisting: true}
}

func globalNamespace_topRootOfQName(n rhino.Node) (uint32, bool) {
	if n.IsNil() {
		return 0, false
	}
	left := n
	for left.IsGetProp() {
		left = left.GetFirstChild()
	}
	switch left.GetToken() {
	case rhino.Token_NAME:
		id := left.GetStringID()
		if id <= 1 {
			return 0, false
		}
		s := strpool.String(strpool.ID(id))
		if i := strings.IndexByte(s, '.'); i >= 0 {
			if i == 0 {
				return 0, false
			}
			return uint32(strpool.Lookup(s[:i])), true
		}
		return id, true
	case rhino.Token_THIS:
		return gnThisStrID, true
	case rhino.Token_SUPER:
		return gnSuperStrID, true
	default:
		return 0, false
	}
}

func (b *GlobalNamespace_BuildGlobalNamespace) putObjLitCache(objLit rhino.Node, prefix string, root rhino.Node, scope *Scope, state gnObjLitState) (string, rhino.Node, *Scope, gnObjLitState) {
	if b.objLitCache == nil {
		b.objLitCache = make(map[rhino.Node]globalNamespace_objLitCacheEntry, 16)
	}
	b.objLitCache[objLit] = globalNamespace_objLitCacheEntry{
		prefix: prefix,
		root:   root,
		scope:  scope,
		state:  state,
	}
	return prefix, root, scope, state
}

func globalNamespace_hasPrototype(s string) bool {
	return strings.HasSuffix(s, ".prototype") || strings.Contains(s, ".prototype.")
}

func globalNamespace_nodeHasPrototype(n rhino.Node) bool {
	for cur := n; !cur.IsNil() && cur.IsGetProp(); cur = cur.GetFirstChild() {
		if cur.GetStringID() == gnPrototypeStrID {
			return true
		}
	}
	return false
}

func (b *GlobalNamespace_BuildGlobalNamespace) resolveObjectLit(objLit rhino.Node, scope *Scope, t *NodeTraversal) (string, rhino.Node, *Scope, gnObjLitState) {
	if e, ok := b.objLitCache[objLit]; ok {
		return e.prefix, e.root, e.scope, e.state
	}

	g := b.outer
	owner := NodeUtil_GetBestLValue(objLit)
	if owner.IsNil() {
		objLitParent := objLit.GetParent()
		if NodeUtil_IsObjectDefinePropertiesDefinition(objLitParent) {
			receiverNode := ast.MustCall(objLitParent).MustFirstArg().Node()
			if globalNamespace_nodeHasPrototype(receiverNode) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
			}
			topID, ok := globalNamespace_topRootOfQName(receiverNode)
			if !ok {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
			}
			jutil.CheckStateMsg1(objLit == receiverNode.GetNext(), "%s", objLitParent)
			root := b.getRootNodeOfTopNameByID(topID, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
			}
			qname, ok := receiverNode.GetQualifiedNameString()
			if !ok || globalNamespace_hasPrototype(qname) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
			}
			if scope == nil && !g.forProcessDefines {
				scope = t.GetScope()
			}
			return b.putObjLitCache(objLit, qname, root, scope, gnObjLitDefineProperties)
		}
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}

	ownerParent := owner.GetParent()
	if ownerParent.IsNil() {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	if ownerParent.IsClassMembers() && !owner.IsComputedProp() {
		if !owner.IsStaticMember() || owner.GetStringID() == gnPrototypeStrID {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		nameNode := NodeUtil_GetNameNode(ownerParent.GetParent())
		if globalNamespace_nodeHasPrototype(nameNode) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		topID, ok := globalNamespace_topRootOfQName(nameNode)
		if !ok {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		root := b.getRootNodeOfTopNameByID(topID, scope, t)
		if !b.IsTopLevelScopeRoot(root) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		className, ok := nameNode.GetQualifiedNameString()
		if !ok {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		ownerName := className + "." + owner.GetString()
		if globalNamespace_hasPrototype(ownerName) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		if scope == nil && !g.forProcessDefines {
			scope = t.GetScope()
		}
		return b.putObjLitCache(objLit, ownerName, root, scope, gnObjLitLValue)
	}
	if NodeUtil_MayBeObjectLitKey(owner) {
		if !ownerParent.IsObjectLit() {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		parentPrefix, root, resolvedScope, parentState := b.resolveObjectLit(ownerParent, scope, t)
		if parentState != gnObjLitLValue {
			return b.putObjLitCache(objLit, "", rhino.NilNode, resolvedScope, gnObjLitInvalid)
		}
		key := NodeUtil_GetObjectOrClassLitKeyName(owner)
		if key == "prototype" || !rhino.TokenStream_IsJSIdentifier(key) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, resolvedScope, gnObjLitInvalid)
		}
		return b.putObjLitCache(objLit, parentPrefix+"."+key, root, resolvedScope, gnObjLitLValue)
	}

	if owner.IsName() && t != nil && b.curModuleRoot.IsNil() && !t.currentHoistScopeRoot.IsNil() {
		if NodeUtil_IsNameDeclaration(ownerParent) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
	}
	if globalNamespace_nodeHasPrototype(owner) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	topID, ok := globalNamespace_topRootOfQName(owner)
	if !ok {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	root := b.getRootNodeOfTopNameByID(topID, scope, t)
	if !b.IsTopLevelScopeRoot(root) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	ownerName, ok := owner.GetQualifiedNameString()
	if !ok || globalNamespace_hasPrototype(ownerName) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	if scope == nil && !g.forProcessDefines {
		scope = t.GetScope()
	}
	return b.putObjLitCache(objLit, ownerName, root, scope, gnObjLitLValue)
}

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopNameByID(topID uint32, s *Scope, t *NodeTraversal) rhino.Node {
	if topID <= 1 {
		return rhino.NilNode
	}
	g := b.outer
	if s != nil || t == nil {
		v := s.GetVarByID(topID)
		if v == nil {
			if g.hasProvidedNames && g.isProvidedRootID(topID) {
				return g.globalRoot
			}
			return rhino.NilNode
		}
		if v.IsLocal() {
			return v.GetScopeRoot()
		}
		return g.globalRoot
	}
	if b.curModuleRoot.IsNil() && len(t.scopes) > 0 {
		globalScope := t.scopes[0].untypedScope
		if globalScope == nil {
			globalScope = t.getScopeSlow(0)
		}
		if globalScope.GetVarByID(topID) == nil {
			if !g.hasProvidedNames || !g.isProvidedRootID(topID) {
				return rhino.NilNode
			}
		}
		if len(t.scopes) == 1 {
			return g.globalRoot
		}
	}
	v := t.LookupVarByID(topID)
	if v == nil {
		if g.hasProvidedNames && g.isProvidedRootID(topID) {
			return g.globalRoot
		}
		return rhino.NilNode
	}
	if v.IsLocal() {
		return v.GetScopeRoot()
	}
	return g.globalRoot
}

// ShouldTraverse collects the references in pre-order.
//
// Java: GlobalNamespace.BuildGlobalNamespace.shouldTraverse(NodeTraversal, Node, Node)
func (b *GlobalNamespace_BuildGlobalNamespace) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	g := b.outer
	if n.IsScript() && g.shouldTraverseScript != nil && !g.shouldTraverseScript(n) {
		return false
	}
	if g.HasExternsRoot() && n.IsScript() {
		// When checking type-summary files, we want to consider them like normal code
		// for some things (like alias inlining) but like externs for other things.
		g.sourceKind = GlobalNamespace_SourceKind_FromScriptNode(n)
	} else if n == g.root {
		g.sourceKind = GlobalNamespace_SourceKind_CODE
	}
	if g.forProcessDefines && g.processDefinesCheck && g.sourceKind == GlobalNamespace_SourceKind_CODE {
		if jsdoc := n.GetJSDocInfo(); (jsdoc != nil && jsdoc.IsDefine()) ||
			(n.IsCall() && n.GetFirstChild().MatchesQualifiedNameNode(processDefines_GOOG_DEFINE)) ||
			(n.IsName() && n.GetStringID() == gnClosureDefinesStrID) {
			if !t.currentScript.IsNil() {
				if g.defineCandidateScripts == nil {
					g.defineCandidateScripts = make(map[rhino.Node]struct{}, 8)
				}
				g.defineCandidateScripts[t.currentScript] = struct{}{}
			} else {
				g.hasNonScriptDefineCand = true
			}
		}
	}
	if n.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		b.SetupModuleMetadata(n)
	} else if n.IsScript() || NodeUtil_IsBundledGoogModuleCall(n) {
		b.curModuleRoot = rhino.NilNode
		b.curMetadata = nil
	}

	b.collect(nil, nil, t, n)

	return true
}

// SetupModuleMetadata initializes the ModuleMetadata for a goog;.module or
// ES module.
//
// moduleRoot is either a MODULE_BODY or a goog.loadModule BLOCK.
//
// Java: GlobalNamespace.BuildGlobalNamespace.setupModuleMetadata(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) SetupModuleMetadata(moduleRoot rhino.Node) {
	moduleMap := b.outer.compiler.GetModuleMap()
	if moduleMap == nil {
		return
	}
	b.curModuleRoot = moduleRoot

	b.curMetadata = jutil.CheckNotNull(
		ModuleImportResolver_GetModuleFromScopeRoot(moduleMap, b.outer.compiler, moduleRoot).
			Metadata())
	if b.curMetadata.IsGoogModule() {
		b.GetOrCreateName("exports", b.curMetadata)
	}
}

// Java: GlobalNamespace.BuildGlobalNamespace.collect(Scope, JSChunk, Node)
func (b *GlobalNamespace_BuildGlobalNamespace) Collect(scope *Scope, chunk *JSChunk, n rhino.Node) {
	b.collect(scope, chunk, nil, n)
}

func (b *GlobalNamespace_BuildGlobalNamespace) collect(scope *Scope, chunk *JSChunk, t *NodeTraversal, n rhino.Node) {
	g := b.outer
	parent := n.GetParent()

	var nameStr string
	var resolvedName *GlobalNamespace_Name
	var root rhino.Node
	hasName := false
	isSet := false
	type_ := GlobalNamespace_NameProp_OTHER_OBJECT

	switch n.GetToken() {
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF:
		if parent.IsObjectLit() {
			prefix, r, s, state := b.resolveObjectLit(parent, scope, t)
			if state != gnObjLitLValue || n.GetStringID() == gnPrototypeStrID {
				return
			}
			key := n.GetString()
			if !rhino.TokenStream_IsJSIdentifier(key) {
				return
			}
			nameStr = prefix + "." + key
			root = r
			scope = s
			hasName = true
		} else if parent.IsClassMembers() {
			if !n.IsStaticMember() || n.IsComputedProp() {
				return // Within a class, only static members define global names.
			}
			nameNode := NodeUtil_GetNameNode(parent.GetParent())
			topID, ok := globalNamespace_topRootOfQName(nameNode)
			if !ok {
				return
			}
			root = b.getRootNodeOfTopNameByID(topID, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
			className, ok := nameNode.GetQualifiedNameString()
			if !ok {
				return
			}
			nameStr = className + "." + n.GetString()
			hasName = true
		} else {
			nameStr, hasName = nodeUtil_getBestLValueNameString(n)
			if !hasName {
				return
			}
			root = b.getRootNodeOfTopNameByID(uint32(strpool.Lookup(GlobalNamespace_GetTopVarName(nameStr))), scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
		}
		isSet = true

		if n.IsMemberFunctionDef() {
			type_ = GlobalNamespace_NameProp_FUNCTION
		} else if n.IsMemberFieldDef() {
			type_ = GlobalNamespace_NameProp_OTHER_OBJECT
		} else {
			type_ = GlobalNamespace_NameProp_GET_SET
		}
	case rhino.Token_STRING_KEY:
		sk := ast.MustStringKey(n)
		if parent.IsObjectLit() {
			prefix, r, s, state := b.resolveObjectLit(parent, scope, t)
			if state == gnObjLitInvalid || n.GetStringID() == gnPrototypeStrID {
				return
			}
			key := sk.Key()
			if state == gnObjLitLValue {
				if !rhino.TokenStream_IsJSIdentifier(key) {
					return
				}
				nameStr = prefix + "." + key
				type_ = b.GetValueType(sk.Value().Node())
			} else {
				nameStr = prefix + "." + key
				type_ = GlobalNamespace_NameProp_GET_SET
			}
			root = r
			scope = s
			hasName = true
			isSet = true
		} else if parent.IsObjectPattern() {
			if g.forProcessDefines && (!g.processDefinesCheck || (t != nil && !t.currentHoistScopeRoot.IsNil() && b.curModuleRoot.IsNil())) {
				return
			}
			patternParent := parent.GetParent()
			var rhs rhino.Node
			if assign, ok := ast.AsAssign(patternParent); ok {
				rhs = assign.Value().Node()
			} else if dlhs, ok := ast.AsDestructuringLhs(patternParent); ok {
				if init, hasInit := dlhs.Init(); hasInit {
					rhs = init.Node()
				}
			} else {
				return
			}
			topID, ok := globalNamespace_topRootOfQName(rhs)
			if !ok {
				return
			}
			root = b.getRootNodeOfTopNameByID(topID, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
			qname, ok := rhs.GetQualifiedNameString()
			if !ok {
				return
			}
			nameStr = qname + "." + sk.Key()
			type_ = b.GetValueType(sk.Value().Node())
			hasName = true
			// not a set
		} else {
			// else not a reference we should record
			return
		}
	case rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		// "a.b?.c" is not a reference to the global name "a.b.c" for the
		// purposes of GlobalNamespace, but we include OPTCHAIN_GETPROP to
		// detect optional chain hasOwnProperty guards.
		if parent.IsGetProp() {
			// This name is nested in a getprop. Return and only create a Ref for the outermost
			// getprop in the chain.
			return
		}
		role := ast.RoleIn(n, parent)
		if !g.forProcessDefines && (role == ast.RoleCallCallee || role == ast.RoleOptChainCallCallee) &&
			b.IsObjectHasOwnPropertyCall(parent) {
			qname, _ := n.GetFirstChild().GetQualifiedNameString()
			globalName := b.GetOrCreateName(qname, b.curMetadata)
			globalName.SetBooleanProperty(GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY)
		}
		if n.GetToken() == rhino.Token_OPTCHAIN_GETPROP {
			return
		}
		// TODO(b/127505242): CAST parents may indicate a set.
		// This may be a variable get or set.
		switch role {
		case ast.RoleNameDeclarationDeclarators:
			isSet = true
			if init, hasInit := ast.MustName(n).Init(); hasInit {
				type_ = b.GetValueType(init.Node())
			} else {
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		case ast.RoleAssignTarget:
			isSet = true
			type_ = b.GetValueType(ast.MustAssign(parent).Value().Node())
		case ast.RoleFunctionName, ast.RoleFunctionBody:
			grandparent := parent.GetParent()
			if grandparent.IsNil() || NodeUtil_IsFunctionExpression(parent) {
				return
			}
			isSet = true
			type_ = GlobalNamespace_NameProp_FUNCTION
		case ast.RoleCatchParam, ast.RoleIncDecTarget,
			ast.RoleForInTarget, ast.RoleForOfTarget, ast.RoleForAwaitOfTarget,
			ast.RoleAssignOpTarget:
			isSet = true
			type_ = GlobalNamespace_NameProp_OTHER_OBJECT
		case ast.RoleClassName:
			// The first child is the class name, and the second child is the superclass name.
			isSet = true
			type_ = GlobalNamespace_NameProp_CLASS
		case ast.RoleStringKeyValue,
			ast.RoleArrayPatternElements,
			ast.RoleDefaultValueTarget,
			ast.RoleDefaultValueValue,
			ast.RoleComputedPropKey,
			ast.RoleComputedPropValue,
			ast.RoleIterRestTarget,
			ast.RoleObjectRestTarget:
			// This may be a set.
			// TODO(b/120303257): this should extend to qnames too, but doing
			// so causes invalid output. Covered in CollapsePropertiesTest
			if n.IsName() && NodeUtil_IsLhsByDestructuring(n) {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		}
		if g.forProcessDefines && !isSet && (!g.processDefinesCheck || (t != nil && !t.currentHoistScopeRoot.IsNil() && b.curModuleRoot.IsNil())) {
			return
		}
		if n.GetToken() == rhino.Token_NAME {
			topID := n.GetStringID()
			if topID <= 1 || role == ast.RoleParamListParams {
				return
			}
			if t != nil && b.curModuleRoot.IsNil() && !t.currentHoistScopeRoot.IsNil() {
				switch role {
				case ast.RoleNameDeclarationDeclarators, ast.RoleFunctionName, ast.RoleCatchParam, ast.RoleClassName:
					return
				}
			}
			topName := strpool.String(strpool.ID(topID))
			if strings.IndexByte(topName, '.') < 0 {
				root = b.getRootNodeOfTopNameByID(topID, scope, t)
				if !b.IsTopLevelScopeRoot(root) || b.skipsLeafID(topID) {
					return
				}
				var nameMetadata *ModuleMetadataMap_ModuleMetadata
				if root != g.globalRoot {
					nameMetadata = b.curMetadata
				}
				resolvedName = b.getOrCreateRootByID(topID, nameMetadata)
			} else {
				nameStr = topName
				root = b.getRootNodeOfTopNameByID(uint32(strpool.Lookup(GlobalNamespace_GetTopVarName(nameStr))), scope, t)
				if !b.IsTopLevelScopeRoot(root) {
					return
				}
				hasName = true
			}
		} else {
			var segBuf [32]uint32
			var segExtra []uint32
			segLen := 0
			protoCut := 0
			hasDotInSeg := false
			var protoTarget rhino.Node
			left := n
			for left.IsGetProp() {
				segID := left.GetStringID()
				if segID == gnPrototypeStrID {
					protoTarget = left.GetFirstChild()
					protoCut = segLen + 1
				} else if !hasDotInSeg && strings.IndexByte(strpool.String(strpool.ID(segID)), '.') >= 0 {
					hasDotInSeg = true
				}
				if segLen < len(segBuf) {
					segBuf[segLen] = segID
				} else {
					if segExtra == nil {
						segExtra = make([]uint32, 0, 64)
						segExtra = append(segExtra, segBuf[:]...)
					}
					segExtra = append(segExtra, segID)
				}
				segLen++
				left = left.GetFirstChild()
			}
			if !left.IsName() {
				return
			}
			topID := left.GetStringID()
			if topID <= 1 {
				return
			}
			topName := strpool.String(strpool.ID(topID))
			if !hasDotInSeg && strings.IndexByte(topName, '.') >= 0 {
				hasDotInSeg = true
			}
			if !hasDotInSeg {
				root = b.getRootNodeOfTopNameByID(topID, scope, t)
				if !b.IsTopLevelScopeRoot(root) {
					return
				}
				segs := segBuf[:segLen]
				if segExtra != nil {
					segs = segExtra
				}
				var nameMetadata *ModuleMetadataMap_ModuleMetadata
				if root != g.globalRoot {
					nameMetadata = b.curMetadata
				}
				if protoTarget.IsNil() {
					if b.skipsLeafID(segs[0]) {
						return
					}
					resolvedName = b.resolveChainByID(topID, segs, nameMetadata)
				} else {
					leafID := topID
					if protoCut < len(segs) {
						leafID = segs[protoCut]
					}
					if b.skipsLeafID(leafID) {
						return
					}
					if !g.forProcessDefines || scope != nil || t == nil {
						if scope == nil {
							scope = t.GetScope()
						}
						if chunk == nil && t != nil {
							chunk = t.GetChunk()
						}
					}
					b.AddOrConfirmRef(b.resolveChainByID(topID, segs[protoCut:], nameMetadata), protoTarget, GlobalNamespace_Ref_Type_PROTOTYPE_GET, scope, chunk)
					return
				}
			} else {
				root = b.getRootNodeOfTopNameByID(uint32(strpool.Lookup(GlobalNamespace_GetTopVarName(topName))), scope, t)
				if !b.IsTopLevelScopeRoot(root) {
					return
				}
				if protoTarget.IsNil() {
					var ok bool
					nameStr, ok = n.GetQualifiedNameString()
					if !ok {
						return
					}
					hasName = true
				} else {
					prefix, ok := protoTarget.GetQualifiedNameString()
					if !ok || b.skipsName(prefix) {
						return
					}
					var nameMetadata *ModuleMetadataMap_ModuleMetadata
					if root != g.globalRoot {
						nameMetadata = b.curMetadata
					}
					if !g.forProcessDefines || scope != nil || t == nil {
						if scope == nil {
							scope = t.GetScope()
						}
						if chunk == nil && t != nil {
							chunk = t.GetChunk()
						}
					}
					b.AddOrConfirmRef(b.GetOrCreateName(prefix, nameMetadata), protoTarget, GlobalNamespace_Ref_Type_PROTOTYPE_GET, scope, chunk)
					return
				}
			}
		}
	case rhino.Token_CALL:
		if parent.IsExprResult() {
			call := ast.MustCall(n)
			if GlobalNamespace_GOOG_PROVIDE.Matches(call.Callee().Node()) {
				if arg, ok := call.FirstArg(); ok {
					if str, ok := arg.AsStringLit(); ok {
						// goog.provide goes through a different code path than regular sets because it can
						// create multiple names, e.g. `goog.provide('a.b.c');` creates the global names
						// a, a.b, and a.b.c. Other sets only create a single global name.
						b.CreateNamesFromProvide(str.Value())
						return
					}
				}
			}
		}
		return
	default:
		return
	}

	if resolvedName != nil {
		if g.forProcessDefines && scope == nil && t != nil {
			if isSet {
				hoistIdx, hoistRoot := t.closestScopeHoistRoot()
				if hoistIdx == 0 || (root != g.globalRoot && hoistRoot == b.curModuleRoot) {
					b.handleResolvedSetFromGlobal(nil, nil, n, resolvedName, type_)
				} else {
					b.handleResolvedSetFromLocal(nil, nil, n, resolvedName)
				}
			} else {
				b.handleResolvedGet(nil, nil, n, resolvedName)
			}
			return
		}
		if scope == nil {
			scope = t.GetScope()
		}
		if chunk == nil && t != nil {
			chunk = t.GetChunk()
		}
		if isSet {
			hoistScope := scope.GetClosestHoistScope()
			if hoistScope.IsGlobal() ||
				(root != g.globalRoot && hoistScope.GetRootNode() == b.curModuleRoot) {
				b.handleResolvedSetFromGlobal(chunk, scope, n, resolvedName, type_)
			} else {
				b.handleResolvedSetFromLocal(chunk, scope, n, resolvedName)
			}
		} else {
			b.handleResolvedGet(chunk, scope, n, resolvedName)
		}
		return
	}

	var nameMetadata *ModuleMetadataMap_ModuleMetadata
	if root != g.globalRoot {
		nameMetadata = b.curMetadata
	}
	if g.forProcessDefines && scope == nil && t != nil {
		if isSet {
			hoistIdx, hoistRoot := t.closestScopeHoistRoot()
			if hoistIdx == 0 || (root != g.globalRoot && hoistRoot == b.curModuleRoot) {
				b.HandleSetFromGlobal(nil, nil, n, nameStr, type_, nameMetadata)
			} else {
				b.HandleSetFromLocal(nil, nil, n, nameStr, nameMetadata)
			}
		} else {
			b.HandleGet(nil, nil, n, nameStr, nameMetadata)
		}
		return
	}
	if scope == nil {
		scope = t.GetScope()
	}
	if chunk == nil && t != nil {
		chunk = t.GetChunk()
	}
	if isSet {
		// Use the closest hoist scope to select handleSetFromGlobal or handleSetFromLocal
		// because they use the term 'global' in an ES5, pre-block-scoping sense.
		hoistScope := scope.GetClosestHoistScope()
		// Consider a set to be 'global' if it is in the hoist scope in which the name is defined.
		// For example, a global name set in a module scope is a 'local' set, but a module-level
		// name set in a module scope is a 'global' set.
		if hoistScope.IsGlobal() ||
			(root != g.globalRoot && hoistScope.GetRootNode() == b.curModuleRoot) {
			b.HandleSetFromGlobal(chunk, scope, n, nameStr, type_, nameMetadata)
		} else {
			b.HandleSetFromLocal(chunk, scope, n, nameStr, nameMetadata)
		}
	} else {
		b.HandleGet(chunk, scope, n, nameStr, nameMetadata)
	}
}

// CreateNamesFromProvide declares all subnamespaces from
// `goog.provide('some.long.namespace')` globally.
//
// Java: GlobalNamespace.BuildGlobalNamespace.createNamesFromProvide(String)
func (b *GlobalNamespace_BuildGlobalNamespace) CreateNamesFromProvide(namespace string) {
	b.outer.hasProvidedNames = true
	var name *GlobalNamespace_Name
	dot := 0

	for dot >= 0 {
		if dot+1 < len(namespace) {
			if idx := strings.IndexByte(namespace[dot+1:], '.'); idx >= 0 {
				dot = dot + 1 + idx
			} else {
				dot = -1
			}
		} else {
			dot = -1
		}
		var subNamespace string
		if dot < 0 {
			subNamespace = namespace
		} else {
			subNamespace = namespace[:dot]
		}
		jutil.CheckState(subNamespace != "")
		// The top-level name resolves the roots of qualified names
		// (getRootNodeOfTopName), so it is kept even for ProcessDefines.
		if name != nil && b.skipsName(subNamespace) {
			continue
		}
		name = b.GetOrCreateName(subNamespace, nil)
		name.SetBooleanProperty(GlobalNamespace_NameProp_IS_PROVIDED)
	}

	if b.skipsName(namespace) {
		return
	}
	newName := b.GetOrCreateName(namespace, nil)
	newName.SetBooleanProperty(GlobalNamespace_NameProp_IS_PROVIDED)
	clear(b.objLitCache)
}

// IsTopLevelScopeRoot reports whether the given name root represents a
// global or module-level name.
//
// This method will return false for functions and blocks and true for
// module bodies and the globalRoot. The one exception is if the function or
// block is from a goog.loadModule argument, as those functions/blocks are
// treated as module roots.
//
// Java: GlobalNamespace.BuildGlobalNamespace.isTopLevelScopeRoot(Node)
func (b *GlobalNamespace_BuildGlobalNamespace) IsTopLevelScopeRoot(root rhino.Node) bool {
	if root.IsNil() {
		return false
	} else if root == b.outer.globalRoot {
		return true
	} else if root == b.curModuleRoot {
		return true
	}
	// Given
	//   goog.loadModule(function(exports) {
	// pretend that assignments to `exports` or `exports.x = ...` are scoped to the function body,
	// although `exports` is really in the enclosing function parameter scope.
	return !b.curModuleRoot.IsNil() && b.curModuleRoot.IsBlock() && root == b.curModuleRoot.GetParent()
}
