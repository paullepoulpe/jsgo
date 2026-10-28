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
	GetPassGraph() graph.GraphvizGraph
	GetBasePassConfig() PassConfig
	GetPassConfigBase() *PassConfigBase
}

// PassConfigBase holds PassConfig's fields and base implementations.
type PassConfigBase struct {
	self    PassConfig
	options *CompilerOptions
}

func (b *PassConfigBase) InitPassConfigBase(options *CompilerOptions) {
	b.options = options
}

func (b *PassConfigBase) GetPassConfigBase() *PassConfigBase { return b }

func (b *PassConfigBase) Options() *CompilerOptions { return b.options }

func (b *PassConfigBase) GetWhitespaceOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(b.options)
}

func (b *PassConfigBase) GetTranspileOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(b.options)
}

func (b *PassConfigBase) GetPassGraph() graph.GraphvizGraph {
	passGraph := graph.LinkedDirectedGraph_CreateWithoutAnnotations[string, string]()
	allPasses := append(b.self.GetChecks().Build(), b.self.GetOptimizations().Build()...)
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

func (b *PassConfigBase) GetBasePassConfig() PassConfig {
	current := b.self
	for {
		delegate, ok := current.(PassConfig_PassConfigDelegate)
		if !ok {
			return current
		}
		current = delegate.GetPassConfig_PassConfigDelegateBase().delegate
	}
}

// PassConfig_PassConfigDelegate proxies PassConfig method calls to an inner delegate.
type PassConfig_PassConfigDelegate interface {
	PassConfig
	GetPassConfig_PassConfigDelegateBase() *PassConfig_PassConfigDelegateBase
}

// PassConfig_PassConfigDelegateBase holds PassConfigDelegate's fields and base implementations.
type PassConfig_PassConfigDelegateBase struct {
	PassConfigBase
	delegate PassConfig
}

func NewPassConfig_PassConfigDelegate(delegate PassConfig) *PassConfig_PassConfigDelegateBase {
	d := &PassConfig_PassConfigDelegateBase{}
	d.InitPassConfigDelegate(d, delegate)
	return d
}

func (d *PassConfig_PassConfigDelegateBase) InitPassConfigDelegate(self PassConfig, delegate PassConfig) {
	d.self = self
	d.InitPassConfig_PassConfigDelegateBase(delegate)
}

func (d *PassConfig_PassConfigDelegateBase) InitPassConfig_PassConfigDelegateBase(delegate PassConfig) {
	d.InitPassConfigBase(delegate.GetPassConfigBase().options)
	d.delegate = delegate
}

func (d *PassConfig_PassConfigDelegateBase) GetPassConfig_PassConfigDelegateBase() *PassConfig_PassConfigDelegateBase {
	return d
}

func (d *PassConfig_PassConfigDelegateBase) GetWhitespaceOnlyPasses() *PassListBuilder {
	return d.delegate.GetWhitespaceOnlyPasses()
}

func (d *PassConfig_PassConfigDelegateBase) GetChecks() *PassListBuilder {
	return d.delegate.GetChecks()
}

func (d *PassConfig_PassConfigDelegateBase) GetOptimizations() *PassListBuilder {
	return d.delegate.GetOptimizations()
}

func (d *PassConfig_PassConfigDelegateBase) GetFinalizations() *PassListBuilder {
	return d.delegate.GetFinalizations()
}

func (d *PassConfig_PassConfigDelegateBase) GetTranspileOnlyPasses() *PassListBuilder {
	return d.delegate.GetTranspileOnlyPasses()
}

// ConformancePassConfig runs only the user-supplied conformance checks and any
// earlier passes required by conformance.
type ConformancePassConfig struct {
	PassConfig_PassConfigDelegateBase
}

func NewConformancePassConfig(delegate PassConfig) *ConformancePassConfig {
	c := &ConformancePassConfig{}
	c.InitPassConfigDelegate(c, delegate)
	return c
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
	PassConfig_PassConfigDelegateBase
}

func NewCheckExtraRequiresPassConfig(options *CompilerOptions) *CheckExtraRequiresPassConfig {
	c := &CheckExtraRequiresPassConfig{}
	c.InitPassConfigDelegate(c, NewDefaultPassConfig(options))
	return c
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
