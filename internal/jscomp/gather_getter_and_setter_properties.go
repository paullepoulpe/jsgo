// Port of com/google/javascript/jscomp/GatherGetterAndSetterProperties.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	gatherPropGetID           = uint32(strpool.Intern("get"))
	gatherPropSetID           = uint32(strpool.Intern("set"))
	gatherCallbackVisitTokens = makeVisitTokenSet(
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_CALL,
	)
)

// GatherGetterAndSetterProperties finds getter and setter properties in the
// AST.
//
// Used to back off certain optimizations, e.g. code removal.
type GatherGetterAndSetterProperties struct {
	compiler AbstractCompiler
}

// Java: GatherGetterAndSetterProperties(AbstractCompiler)
func NewGatherGetterAndSetterProperties(compiler AbstractCompiler) *GatherGetterAndSetterProperties {
	return &GatherGetterAndSetterProperties{compiler: compiler}
}

// Java: GatherGetterAndSetterProperties.process(Node, Node)
func (g *GatherGetterAndSetterProperties) Process(externs rhino.Node, root rhino.Node) {
	GatherGetterAndSetterProperties_Update(g.compiler, externs, root)
}

// GatherGetterAndSetterProperties_Update gathers all getters and setters in
// the AST.
//
// Java: GatherGetterAndSetterProperties.update(AbstractCompiler, Node, Node)
func GatherGetterAndSetterProperties_Update(compiler AbstractCompiler, externs rhino.Node, root rhino.Node) {
	// TODO(nickreid): We probably don't need to re-gather from the externs. They don't change so
	// the first collection should be good forever.
	// For now we traverse both trees every time because there's no reason we have to treat them
	// differently.
	jutil.CheckState(externs.GetParent() == root.GetParent())
	var accessorSummary *AccessorSummary
	if compiler.GetOptions().GetAssumePropertiesAreStaticallyAnalyzable() {
		accessorSummary = gatherAccessorSummary(compiler, externs.GetParent())
	} else {
		accessorSummary = AccessorSummary_CreateAssumingAlwaysGetterAndSetter()
	}
	compiler.SetAccessorSummary(accessorSummary)
}

// Java: GatherGetterAndSetterProperties.gather(AbstractCompiler, Node)
func GatherGetterAndSetterProperties_Gather(compiler AbstractCompiler, root rhino.Node) *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind] {
	return gatherAccessorSummary(compiler, root).GetAccessors()
}

func gatherAccessorSummary(compiler AbstractCompiler, root rhino.Node) *AccessorSummary {
	if out, ok := gatherAccessorSummaryParallel(compiler, root); ok {
		return out
	}
	var gatherCallback GatherGetterAndSetterProperties_GatherCallback
	NodeTraversal_Traverse(compiler, root, &gatherCallback)
	if len(gatherCallback.entries) == 0 {
		return &AccessorSummary{}
	}
	var batches [1][]accessorIDEntry
	batches[0] = gatherCallback.entries
	return newAccessorSummaryFromIDs(batches[:])
}

func collectGatherScripts(n rhino.Node, out *[]rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	if n.IsScript() {
		*out = append(*out, n)
		return true
	}
	if !n.IsRoot() {
		return false
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !collectGatherScripts(c, out) {
			return false
		}
	}
	return true
}

func gatherAccessorSummaryParallel(
	compiler AbstractCompiler,
	root rhino.Node,
) (*AccessorSummary, bool) {
	var scripts []rhino.Node
	if !collectGatherScripts(root, &scripts) {
		return nil, false
	}
	numWorkers := parallelThreadCount(compiler, len(scripts))
	if numWorkers <= 1 {
		return nil, false
	}

	prewarmCompilerForParallelReads(compiler)
	perScript := make([][]accessorIDEntry, len(scripts))
	runParallelWork(compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		runItem := func(idx int, rc *recordingCompiler) {
			var cb GatherGetterAndSetterProperties_GatherCallback
			NodeTraversal_Traverse(rc, scripts[idx], &cb)
			if len(cb.entries) > 0 {
				perScript[idx] = cb.entries
			}
		}
		return runItem, nil
	}, false)

	return newAccessorSummaryFromIDs(perScript), true
}

// GatherGetterAndSetterProperties_GatherCallback is the traversal that
// records the accessors.
type GatherGetterAndSetterProperties_GatherCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	entries []accessorIDEntry
}

var (
	_ nodeTraversal_noVisitLeaves = (*GatherGetterAndSetterProperties_GatherCallback)(nil)
	_ nodeTraversal_visitTokens   = (*GatherGetterAndSetterProperties_GatherCallback)(nil)
)

func (*GatherGetterAndSetterProperties_GatherCallback) noVisitLeaves() {}

func (*GatherGetterAndSetterProperties_GatherCallback) visitTokens() *ast.TokenSet {
	return &gatherCallbackVisitTokens
}

// Java: GatherCallback()
func NewGatherGetterAndSetterProperties_GatherCallback() *GatherGetterAndSetterProperties_GatherCallback {
	return &GatherGetterAndSetterProperties_GatherCallback{}
}

func (c *GatherGetterAndSetterProperties_GatherCallback) recordByID(propertyID uint32, kind AccessorSummary_PropertyAccessKind) {
	c.entries = append(c.entries, accessorIDEntry{id: propertyID, kind: kind})
}

// Java: GatherCallback.record(String, PropertyAccessKind)
func (c *GatherGetterAndSetterProperties_GatherCallback) Record(property string, kind AccessorSummary_PropertyAccessKind) {
	c.recordByID(uint32(strpool.Intern(property)), kind)
}

// Java: GatherCallback.visit(NodeTraversal, Node, Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETTER_DEF:
		c.RecordGetterDef(n)
	case rhino.Token_SETTER_DEF:
		c.RecordSetterDef(n)
	case rhino.Token_CALL:
		if NodeUtil_IsObjectDefinePropertyDefinition(n) {
			c.VisitDefineProperty(n)
		} else if NodeUtil_IsObjectDefinePropertiesDefinition(n) ||
			NodeUtil_IsObjectCreateWithPropertiesDefinition(n) {
			c.VisitDefineProperties(n)
		}
	default:
	}
}

// Java: GatherCallback.recordGetterDef(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) RecordGetterDef(getterDef rhino.Node) {
	jutil.CheckState(getterDef.IsGetterDef())
	c.recordByID(getterDef.GetStringID(), AccessorSummary_PropertyAccessKind_GETTER_ONLY)
}

// Java: GatherCallback.recordSetterDef(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) RecordSetterDef(setterDef rhino.Node) {
	jutil.CheckState(setterDef.IsSetterDef())
	c.recordByID(setterDef.GetStringID(), AccessorSummary_PropertyAccessKind_SETTER_ONLY)
}

func (c *GatherGetterAndSetterProperties_GatherCallback) visitDescriptorByID(propertyID uint32, descriptor rhino.Node) {
	for key := descriptor.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if key.IsStringKey() || key.IsMemberFunctionDef() {
			switch key.GetStringID() {
			case gatherPropGetID:
				c.recordByID(propertyID, AccessorSummary_PropertyAccessKind_GETTER_ONLY)
			case gatherPropSetID:
				c.recordByID(propertyID, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
			}
		}
	}
}

// VisitDescriptor looks for getters and setters passed to
// Object.defineProperty or Object.defineProperties.
//
//	Object.defineProperty(obj, 'propertyName', { /* descriptor */ });
//	Object.defineProperties(obj, { 'propertyName': { /* descriptor */ } });
//
// Java: GatherCallback.visitDescriptor(String, Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) VisitDescriptor(propertyName string, descriptor rhino.Node) {
	c.visitDescriptorByID(uint32(strpool.Intern(propertyName)), descriptor)
}

// Java: GatherCallback.visitDefineProperty(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) VisitDefineProperty(definePropertyCall rhino.Node) {
	propertyNameNode := definePropertyCall.GetChildAtIndex(2)
	descriptor := definePropertyCall.GetChildAtIndex(3)

	if !propertyNameNode.IsStringLit() || !descriptor.IsObjectLit() {
		return
	}

	c.visitDescriptorByID(propertyNameNode.GetStringID(), descriptor)
}

// Java: GatherCallback.visitDefineProperties(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) VisitDefineProperties(definePropertiesCall rhino.Node) {
	props := definePropertiesCall.GetChildAtIndex(2)

	if !props.IsObjectLit() {
		return
	}

	for prop := props.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
		if prop.IsStringKey() && prop.HasOneChild() && prop.GetFirstChild().IsObjectLit() {
			c.visitDescriptorByID(prop.GetStringID(), prop.GetFirstChild())
		}
	}
}
