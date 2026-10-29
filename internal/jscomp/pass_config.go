// Port of com/google/javascript/jscomp/PassConfig.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/graph"
)

// PassConfig holds pass factories and metadata for native Compiler passes.
type PassConfig interface {
	GetWhitespaceOnlyPasses() *PassListBuilder
	GetTranspileOnlyPasses() *PassListBuilder
	GetChecks() *PassListBuilder
	GetOptimizations() *PassListBuilder
	GetFinalizations() *PassListBuilder
	GetBasePassConfig() PassConfig
	Options() *CompilerOptions
}

// BuildPassGraph builds the Graphviz pass graph for pc.
func BuildPassGraph(pc PassConfig) graph.GraphvizGraph {
	passGraph := graph.LinkedDirectedGraph_CreateWithoutAnnotations[string, string]()
	allPasses := append(pc.GetChecks().Build(), pc.GetOptimizations().Build()...)
	var lastPass, loopStart string
	var hasLast, inLoop bool
	for _, pass := range allPasses {
		baseName := pass.GetName()
		passName := baseName
		for i := 1; passGraph.HasNode(passName); i++ {
			passName = baseName + strconv.Itoa(i)
		}
		passGraph.CreateNode(passName)

		if !inLoop && pass.IsRunInFixedPointLoop() {
			loopStart = passName
			inLoop = true
		} else if inLoop && !pass.IsRunInFixedPointLoop() {
			passGraph.Connect(lastPass, "loop", loopStart)
			loopStart = ""
			inLoop = false
		}

		if hasLast {
			passGraph.Connect(lastPass, "", passName)
		}
		lastPass = passName
		hasLast = true
	}
	return passGraph
}

// PassConfigDelegate proxies PassConfig method calls to an inner delegate.
type PassConfigDelegate struct {
	options  *CompilerOptions
	delegate PassConfig
}

func NewPassConfigDelegate(delegate PassConfig) PassConfigDelegate {
	return PassConfigDelegate{
		options:  delegate.Options(),
		delegate: delegate,
	}
}

func (d *PassConfigDelegate) Options() *CompilerOptions {
	return d.options
}

func (d *PassConfigDelegate) GetBasePassConfig() PassConfig {
	return d.delegate.GetBasePassConfig()
}

func (d *PassConfigDelegate) GetWhitespaceOnlyPasses() *PassListBuilder {
	return d.delegate.GetWhitespaceOnlyPasses()
}

func (d *PassConfigDelegate) GetChecks() *PassListBuilder {
	return d.delegate.GetChecks()
}

func (d *PassConfigDelegate) GetOptimizations() *PassListBuilder {
	return d.delegate.GetOptimizations()
}

func (d *PassConfigDelegate) GetFinalizations() *PassListBuilder {
	return d.delegate.GetFinalizations()
}

func (d *PassConfigDelegate) GetTranspileOnlyPasses() *PassListBuilder {
	return d.delegate.GetTranspileOnlyPasses()
}

// ConformancePassConfig runs only the user-supplied conformance checks and any
// earlier passes required by conformance.
type ConformancePassConfig struct {
	PassConfigDelegate
}

func NewConformancePassConfig(delegate PassConfig) *ConformancePassConfig {
	return &ConformancePassConfig{
		PassConfigDelegate: NewPassConfigDelegate(delegate),
	}
}

func (c *ConformancePassConfig) GetChecks() *PassListBuilder {
	// Return every check up to and including "checkConformance", or an empty list if not found.
	passes := NewPassListBuilder(c.options)
	passes.AddAllUpTo(c.delegate.GetChecks(), PassNames_CHECK_CONFORMANCE)
	return passes
}

func (c *ConformancePassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// CheckExtraRequiresPassConfig runs only the CheckExtraRequires linter rule to
// remove unused goog.require imports.
type CheckExtraRequiresPassConfig struct {
	PassConfigDelegate
}

func NewCheckExtraRequiresPassConfig(options *CompilerOptions) *CheckExtraRequiresPassConfig {
	return &CheckExtraRequiresPassConfig{
		PassConfigDelegate: NewPassConfigDelegate(NewDefaultPassConfig(options)),
	}
}

func (c *CheckExtraRequiresPassConfig) GetChecks() *PassListBuilder {
	passes := NewPassListBuilder(c.options)
	passes.MaybeAdd(makePassFactory("removeUnusedImport", func(compiler AbstractCompiler) CompilerPass {
		return NewCombinedCompilerPassOfList(
			compiler,
			[]NodeTraversal_Callback{NewCheckExtraRequires(compiler, c.options.GetUnusedImportsToRemove())},
		)
	}))
	return passes
}

func (c *CheckExtraRequiresPassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}
