// Port of com/google/javascript/jscomp/CollectFileOverviewVisibility.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// CollectFileOverviewVisibility is the compiler pass that collects
// visibility annotations in @fileoverview blocks. Used by
// CheckAccessControls.
type CollectFileOverviewVisibility struct {
	builder *jutil.ImmutableMapBuilder[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]
}

// Java: CollectFileOverviewVisibility(AbstractCompiler)
func NewCollectFileOverviewVisibility(compiler AbstractCompiler) *CollectFileOverviewVisibility {
	return &CollectFileOverviewVisibility{
		builder: jutil.NewImmutableMapBuilder[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility](),
	}
}

// Java: CollectFileOverviewVisibility.process(Node, Node)
func (c *CollectFileOverviewVisibility) Process(externs rhino.Node, root rhino.Node) {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		c.Visit(script)
	}
}

// Java: CollectFileOverviewVisibility.visit(Node)
func (c *CollectFileOverviewVisibility) Visit(scriptNode rhino.Node) {
	jsDocInfo := scriptNode.GetJSDocInfo()
	if jsDocInfo == nil {
		return
	}
	// Go-only: language-adapter — JSDocInfo.getVisibility() never returns
	// null in Java either (it defaults to INHERITED), so Java's
	// `if (v == null) return;` never fires and has no Go counterpart.
	v := jsDocInfo.GetVisibility()
	c.builder.Put(scriptNode.GetStaticSourceFile(), v)
}

// Java: CollectFileOverviewVisibility.getFileOverviewVisibilityMap()
func (c *CollectFileOverviewVisibility) GetFileOverviewVisibilityMap() *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility] {
	return c.builder.BuildOrThrow()
}
