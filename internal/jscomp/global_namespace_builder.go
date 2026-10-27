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

func globalNamespace_topNameOfQName(n rhino.Node) string {
	if n.IsNil() {
		return ""
	}
	left := n
	for left.IsGetProp() {
		left = left.GetFirstChild()
	}
	switch left.GetToken() {
	case rhino.Token_NAME:
		return GlobalNamespace_GetTopVarName(left.GetString())
	case rhino.Token_THIS:
		return "this"
	case rhino.Token_SUPER:
		return "super"
	default:
		return ""
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
		if cur.GetString() == "prototype" {
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
			receiverNode := objLitParent.GetSecondChild()
			if globalNamespace_nodeHasPrototype(receiverNode) {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
			}
			topName := globalNamespace_topNameOfQName(receiverNode)
			if topName == "" {
				return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
			}
			jutil.CheckStateMsg1(objLit == receiverNode.GetNext(), "%s", objLitParent)
			root := b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
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
		if !owner.IsStaticMember() || owner.GetString() == "prototype" {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		nameNode := NodeUtil_GetNameNode(ownerParent.GetParent())
		if globalNamespace_nodeHasPrototype(nameNode) {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		topName := globalNamespace_topNameOfQName(nameNode)
		if topName == "" {
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
		root := b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
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
		switch ownerParent.GetToken() {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
		}
	}
	if globalNamespace_nodeHasPrototype(owner) {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	topName := globalNamespace_topNameOfQName(owner)
	if topName == "" {
		return b.putObjLitCache(objLit, "", rhino.NilNode, scope, gnObjLitInvalid)
	}
	root := b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
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

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopName(topName string, s *Scope) rhino.Node {
	return b.outer.getRootNodeOfTopName(topName, s)
}

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopNameWithTraversal(topName string, s *Scope, t *NodeTraversal) rhino.Node {
	return b.getRootNodeOfTopNameByID(0, topName, s, t)
}

func (b *GlobalNamespace_BuildGlobalNamespace) getRootNodeOfTopNameByID(topID uint32, topName string, s *Scope, t *NodeTraversal) rhino.Node {
	if s != nil || t == nil {
		return b.getRootNodeOfTopName(topName, s)
	}
	g := b.outer
	if b.curModuleRoot.IsNil() && len(t.scopes) > 0 {
		globalScope := t.scopes[0].untypedScope
		if globalScope == nil {
			globalScope = t.getScopeSlow(0)
		}
		if globalScope.GetVar(topName) == nil {
			if !g.hasProvidedNames {
				return rhino.NilNode
			}
			if topID == 0 {
				topID = uint32(strpool.Lookup(topName))
			}
			if !g.isProvidedRootID(topID) {
				return rhino.NilNode
			}
		}
		if len(t.scopes) == 1 {
			return g.globalRoot
		}
	}
	v := t.LookupVar(topName)
	if v == nil {
		if g.hasProvidedNames {
			if topID == 0 {
				topID = uint32(strpool.Lookup(topName))
			}
			if g.isProvidedRootID(topID) {
				return g.globalRoot
			}
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
			(n.IsName() && n.GetString() == "CLOSURE_DEFINES") {
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
			if state != gnObjLitLValue {
				return
			}
			key := n.GetString()
			if key == "prototype" || !rhino.TokenStream_IsJSIdentifier(key) {
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
			topName := globalNamespace_topNameOfQName(nameNode)
			if topName == "" {
				return
			}
			root = b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
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
			root = b.getRootNodeOfTopNameWithTraversal(GlobalNamespace_GetTopVarName(nameStr), scope, t)
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
		if parent.IsObjectLit() {
			prefix, r, s, state := b.resolveObjectLit(parent, scope, t)
			if state == gnObjLitInvalid {
				return
			}
			key := n.GetString()
			if key == "prototype" {
				return
			}
			if state == gnObjLitLValue {
				if !rhino.TokenStream_IsJSIdentifier(key) {
					return
				}
				nameStr = prefix + "." + key
				type_ = b.GetValueType(n.GetOnlyChild())
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
			if !patternParent.IsAssign() && !patternParent.IsDestructuringLhs() {
				return
			}
			rhs := patternParent.GetSecondChild()
			topName := globalNamespace_topNameOfQName(rhs)
			if topName == "" {
				return
			}
			root = b.getRootNodeOfTopNameWithTraversal(topName, scope, t)
			if !b.IsTopLevelScopeRoot(root) {
				return
			}
			qname, ok := rhs.GetQualifiedNameString()
			if !ok {
				return
			}
			nameStr = qname + "." + n.GetString()
			type_ = b.GetValueType(n.GetFirstChild())
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
		pTok := parent.GetToken()
		if pTok == rhino.Token_GETPROP {
			// This name is nested in a getprop. Return and only create a Ref for the outermost
			// getprop in the chain.
			return
		}
		if !g.forProcessDefines && (pTok == rhino.Token_CALL || pTok == rhino.Token_OPTCHAIN_CALL) &&
			n.IsFirstChildOf(parent) && b.IsObjectHasOwnPropertyCall(parent) {
			qname, _ := n.GetFirstChild().GetQualifiedNameString()
			globalName := b.GetOrCreateName(qname, b.curMetadata)
			globalName.SetBooleanProperty(GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY)
		}
		if n.GetToken() == rhino.Token_OPTCHAIN_GETPROP {
			return
		}
		// TODO(b/127505242): CAST parents may indicate a set.
		// This may be a variable get or set.
		switch pTok {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			isSet = true
			rvalue := n.GetFirstChild()
			if rvalue.IsNil() {
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			} else {
				type_ = b.GetValueType(rvalue)
			}
		case rhino.Token_ASSIGN:
			if parent.GetFirstChild() == n {
				isSet = true
				type_ = b.GetValueType(n.GetNext())
			}
		case rhino.Token_FUNCTION:
			grandparent := parent.GetParent()
			if grandparent.IsNil() || NodeUtil_IsFunctionExpression(parent) {
				return
			}
			isSet = true
			type_ = GlobalNamespace_NameProp_FUNCTION
		case rhino.Token_CATCH, rhino.Token_INC, rhino.Token_DEC:
			isSet = true
			type_ = GlobalNamespace_NameProp_OTHER_OBJECT
		case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			if parent.GetFirstChild() == n {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		case rhino.Token_CLASS:
			// The first child is the class name, and the second child is the superclass name.
			if parent.GetFirstChild() == n {
				isSet = true
				type_ = GlobalNamespace_NameProp_CLASS
			}
		case rhino.Token_STRING_KEY,
			rhino.Token_ARRAY_PATTERN,
			rhino.Token_DEFAULT_VALUE,
			rhino.Token_COMPUTED_PROP,
			rhino.Token_ITER_REST,
			rhino.Token_OBJECT_REST:
			// This may be a set.
			// TODO(b/120303257): this should extend to qnames too, but doing
			// so causes invalid output. Covered in CollapsePropertiesTest
			if n.IsName() && NodeUtil_IsLhsByDestructuring(n) {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD, rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
			// isSet = false, type = OTHER.
		default:
			if NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n {
				isSet = true
				type_ = GlobalNamespace_NameProp_OTHER_OBJECT
			}
		}
		if g.forProcessDefines && !isSet && (!g.processDefinesCheck || (t != nil && !t.currentHoistScopeRoot.IsNil() && b.curModuleRoot.IsNil())) {
			return
		}
		if n.GetToken() == rhino.Token_NAME {
			topID := n.GetStringID()
			if topID <= 1 || pTok == rhino.Token_PARAM_LIST {
				return
			}
			if t != nil && b.curModuleRoot.IsNil() && !t.currentHoistScopeRoot.IsNil() {
				switch pTok {
				case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CATCH:
					return
				case rhino.Token_CLASS:
					if parent.GetFirstChild() == n {
						return
					}
				}
			}
			topName := strpool.String(strpool.ID(topID))
			if strings.IndexByte(topName, '.') < 0 {
				root = b.getRootNodeOfTopNameByID(topID, topName, scope, t)
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
				root = b.getRootNodeOfTopNameWithTraversal(GlobalNamespace_GetTopVarName(nameStr), scope, t)
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
				root = b.getRootNodeOfTopNameByID(topID, topName, scope, t)
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
				root = b.getRootNodeOfTopNameWithTraversal(GlobalNamespace_GetTopVarName(topName), scope, t)
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
		if parent.IsExprResult() &&
			GlobalNamespace_GOOG_PROVIDE.Matches(n.GetFirstChild()) &&
			n.GetSecondChild().IsStringLit() {
			// goog.provide goes through a different code path than regular sets because it can
			// create multiple names, e.g. `goog.provide('a.b.c');` creates the global names
			// a, a.b, and a.b.c. Other sets only create a single global name.
			b.CreateNamesFromProvide(n.GetSecondChild().GetString())
			return
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
		return b.GetValueType(n.GetLastChild())
	case rhino.Token_HOOK:
		// The same line of reasoning used for the OR case applies here.
		second := n.GetSecondChild()
		t := b.GetValueType(second)
		if t != GlobalNamespace_NameProp_OTHER_OBJECT {
			return t
		}
		third := second.GetNext()
		return b.GetValueType(third)
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
		if n == parent.GetFirstChild() {
			// It is a call target
			type_ = GlobalNamespace_Ref_Type_CALL_GET
		} else if b.IsClassDefiningCall(parent) {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
		} else {
			type_ = GlobalNamespace_Ref_Type_ALIASING_GET
		}
	case rhino.Token_NEW:
		if n == parent.GetFirstChild() {
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
		if n != parent.GetFirstChild() {
			// This node is y or z in (x?y:z) or (x,y). We only know that an alias is
			// not getting created for this name if the result is assigned to
			// the same name (e.g. var a = a ? a : {}).
			type_ = b.determineRefTypeForGetImpl(parent, referenceNode, name, nameObj)
		} else {
			type_ = GlobalNamespace_Ref_Type_DIRECT_GET
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
		for lhs.IsCast() {
			// Case: `/** @type {!Foo} */ (x) = ...`; or multiple casts like `(cast(cast(x)) =`
			lhs = lhs.GetOnlyChild()
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
	jutil.CheckArgumentMsg1(callNode.IsCall() || callNode.IsOptChainCall(), "%s", callNode)
	if !callNode.HasTwoChildren() {
		return false
	}
	callee := callNode.GetFirstChild()
	if !callee.IsGetProp() && !callee.IsOptChainGetProp() {
		return false
	}
	receiver := callee.GetFirstChild()
	propName := callee.GetString()
	return ("hasOwnProperty" == propName || "propertyIsEnumerable" == propName) &&
		receiver.IsQualifiedName()
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
