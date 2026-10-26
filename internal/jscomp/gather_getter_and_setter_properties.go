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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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
	g := &GatherGetterAndSetterProperties{}
	g.compiler = compiler
	return g
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
		accessorSummary = AccessorSummary_Create(GatherGetterAndSetterProperties_Gather(compiler, externs.GetParent()))
	} else {
		accessorSummary = AccessorSummary_CreateAssumingAlwaysGetterAndSetter()
	}
	compiler.SetAccessorSummary(accessorSummary)
}

// Java: GatherGetterAndSetterProperties.gather(AbstractCompiler, Node)
func GatherGetterAndSetterProperties_Gather(compiler AbstractCompiler, root rhino.Node) *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind] {
	if out, ok := gatherGetterAndSetterPropertiesParallel(compiler, root); ok {
		return out
	}
	gatherCallback := NewGatherGetterAndSetterProperties_GatherCallback()
	NodeTraversal_Traverse(compiler, root, gatherCallback)
	return gatherCallback.properties
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

func gatherGetterAndSetterPropertiesParallel(
	compiler AbstractCompiler,
	root rhino.Node,
) (*jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind], bool) {
	var scripts []rhino.Node
	if !collectGatherScripts(root, &scripts) {
		return nil, false
	}
	numWorkers := parallelThreadCount(compiler, len(scripts))
	if numWorkers <= 1 {
		return nil, false
	}

	prewarmCompilerForParallelReads(compiler)
	perScript := make([]*jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind], len(scripts))
	runParallelWork(compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		runItem := func(idx int, rc *recordingCompiler) {
			cb := NewGatherGetterAndSetterProperties_GatherCallback()
			NodeTraversal_Traverse(rc, scripts[idx], cb)
			if !cb.properties.IsEmpty() {
				perScript[idx] = cb.properties
			}
		}
		return runItem, nil
	}, false)

	merged := NewGatherGetterAndSetterProperties_GatherCallback()
	for _, props := range perScript {
		if props == nil {
			continue
		}
		for prop, kind := range props.All() {
			merged.Record(prop, kind)
		}
	}
	return merged.properties, true
}

// GatherGetterAndSetterProperties_GatherCallback is the traversal that
// records the accessors.
type GatherGetterAndSetterProperties_GatherCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	properties *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind]
}

func (*GatherGetterAndSetterProperties_GatherCallback) noVisitLeaves() {}

// Java: GatherCallback()
func NewGatherGetterAndSetterProperties_GatherCallback() *GatherGetterAndSetterProperties_GatherCallback {
	return &GatherGetterAndSetterProperties_GatherCallback{
		properties: jutil.NewLinkedHashMap[string, AccessorSummary_PropertyAccessKind](),
	}
}

// Java: GatherCallback.record(String, PropertyAccessKind)
func (c *GatherGetterAndSetterProperties_GatherCallback) Record(property string, kind AccessorSummary_PropertyAccessKind) {
	// properties.merge(property, kind, PropertyAccessKind::unionWith)
	if old, ok := c.properties.Lookup(property); ok {
		c.properties.Put(property, old.UnionWith(kind))
	} else {
		c.properties.Put(property, kind)
	}
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

	name := getterDef.GetString()
	c.Record(name, AccessorSummary_PropertyAccessKind_GETTER_ONLY)
}

// Java: GatherCallback.recordSetterDef(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) RecordSetterDef(setterDef rhino.Node) {
	jutil.CheckState(setterDef.IsSetterDef())

	name := setterDef.GetString()
	c.Record(name, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
}

// VisitDescriptor looks for getters and setters passed to
// Object.defineProperty or Object.defineProperties.
//
//	Object.defineProperty(obj, 'propertyName', { /* descriptor */ });
//	Object.defineProperties(obj, { 'propertyName': { /* descriptor */ } });
//
// Java: GatherCallback.visitDescriptor(String, Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) VisitDescriptor(propertyName string, descriptor rhino.Node) {
	for key := descriptor.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if key.IsStringKey() || key.IsMemberFunctionDef() {
			if "get" == key.GetString() {
				c.Record(propertyName, AccessorSummary_PropertyAccessKind_GETTER_ONLY)
			} else if "set" == key.GetString() {
				c.Record(propertyName, AccessorSummary_PropertyAccessKind_SETTER_ONLY)
			}
		}
	}
}

// Java: GatherCallback.visitDefineProperty(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) VisitDefineProperty(definePropertyCall rhino.Node) {
	propertyNameNode := definePropertyCall.GetChildAtIndex(2)
	descriptor := definePropertyCall.GetChildAtIndex(3)

	if !propertyNameNode.IsStringLit() || !descriptor.IsObjectLit() {
		return
	}

	propertyName := propertyNameNode.GetString()
	c.VisitDescriptor(propertyName, descriptor)
}

// Java: GatherCallback.visitDefineProperties(Node)
func (c *GatherGetterAndSetterProperties_GatherCallback) VisitDefineProperties(definePropertiesCall rhino.Node) {
	props := definePropertiesCall.GetChildAtIndex(2)

	if !props.IsObjectLit() {
		return
	}

	for prop := props.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
		if prop.IsStringKey() && prop.HasOneChild() && prop.GetFirstChild().IsObjectLit() {
			propertyName := prop.GetString()
			descriptor := prop.GetFirstChild()

			c.VisitDescriptor(propertyName, descriptor)
		}
	}
}
