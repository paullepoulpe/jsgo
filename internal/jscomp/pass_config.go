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
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PassConfig holds pass factories and meta-data for native Compiler passes.
//
// Go-only: language-adapter — an abstract class with subclasses: this
// interface declares every overridable method, PassConfigBase holds the
// fields and the base implementations. Subclasses embed PassConfigBase, set
// self, then call InitPassConfigBase.
type PassConfig interface {
	// GetWhitespaceOnlyPasses gets additional checking passes that are run
	// always, even in "whitespace only" mode. For very specific cases where
	// processing is required even in a mode which is intended not to have any
	// processing - specifically introduced to support goog.module() usage.
	//
	// Java: PassConfig.getWhitespaceOnlyPasses()
	GetWhitespaceOnlyPasses() *PassListBuilder

	// GetTranspileOnlyPasses gets the transpilation passes
	//
	// Java: PassConfig.getTranspileOnlyPasses()
	GetTranspileOnlyPasses() *PassListBuilder

	// GetChecks gets the checking passes to run.
	//
	// Checking passes revolve around emitting warnings and errors. They also
	// may include pre-processor passes needed to do error analysis more
	// effectively.
	//
	// Clients that only want to analyze code (like IDEs) and not emit code
	// will only run checks and not optimizations.
	//
	// Java: PassConfig.getChecks()
	GetChecks() *PassListBuilder

	// GetOptimizations gets the optimization passes to run.
	//
	// Optimization passes revolve around producing smaller and faster code.
	// They should always run after checking passes.
	//
	// Java: PassConfig.getOptimizations()
	GetOptimizations() *PassListBuilder

	// GetFinalizations gets the finalization passes to run.
	//
	// Finalization passes include the injection of locale-specific code and
	// converting the AST to its final form for output.
	//
	// Java: PassConfig.getFinalizations()
	GetFinalizations() *PassListBuilder

	// GetPassGraph gets a graph of the passes run. For debugging.
	//
	// Java: PassConfig.getPassGraph()
	GetPassGraph() graph.GraphvizGraph

	// GetBasePassConfig finds the first pass provider that does not have a
	// delegate.
	//
	// Java: PassConfig.getBasePassConfig()
	GetBasePassConfig() PassConfig

	GetPassConfigBase() *PassConfigBase
}

// PassConfigBase holds PassConfig's fields and base implementations.
type PassConfigBase struct {
	self PassConfig

	// Used by the subclasses.
	options *CompilerOptions
}

// InitPassConfigBase is the PassConfig(CompilerOptions) constructor.
//
// Java: PassConfig(CompilerOptions)
func (b *PassConfigBase) InitPassConfigBase(options *CompilerOptions) {
	b.options = options
}

// GetPassConfigBase gives access to the base fields through the interface.
func (b *PassConfigBase) GetPassConfigBase() *PassConfigBase { return b }

// Java: PassConfig.getWhitespaceOnlyPasses()
func (b *PassConfigBase) GetWhitespaceOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(b.options)
}

// Java: PassConfig.getTranspileOnlyPasses()
func (b *PassConfigBase) GetTranspileOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(b.options)
}

// Java: PassConfig.getPassGraph()
func (b *PassConfigBase) GetPassGraph() graph.GraphvizGraph {
	passGraph := graph.LinkedDirectedGraph_CreateWithoutAnnotations[string, string]()
	allPasses := append(b.self.GetChecks().Build(), b.self.GetOptimizations().Build()...)
	var lastPass *string = nil
	var loopStart *string = nil
	for _, pass := range allPasses {
		passName := pass.GetName()
		i := 1
		for passGraph.HasNode(passName) {
			passName = pass.GetName() + jutil.ToJavaString(i)
			i++
		}
		passGraph.CreateNode(passName)

		if loopStart == nil && pass.IsRunInFixedPointLoop() {
			loopStart = &passName
		} else if loopStart != nil && !pass.IsRunInFixedPointLoop() {
			passGraph.Connect(*lastPass, "loop", *loopStart)
			loopStart = nil
		}

		if lastPass != nil {
			passGraph.Connect(*lastPass, "", passName)
		}
		lastPass = &passName
	}
	return passGraph
}

// Options returns the CompilerOptions bound to this PassConfig.
func (b *PassConfigBase) Options() *CompilerOptions {
	return b.options
}

// Java: PassConfig.getBasePassConfig()
func (b *PassConfigBase) GetBasePassConfig() PassConfig {
	current := b.self
	for {
		delegate, ok := current.(PassConfig_PassConfigDelegate)
		if !ok {
			break
		}
		current = delegate.GetPassConfig_PassConfigDelegateBase().delegate
	}
	return current
}

// PassConfig_PassConfigDelegate is an implementation of PassConfig that
// just proxies all its method calls into an inner class.
//
// Go-only: language-adapter — a class with subclasses: this interface plus
// PassConfig_PassConfigDelegateBase, which subclasses embed.
type PassConfig_PassConfigDelegate interface {
	PassConfig

	GetPassConfig_PassConfigDelegateBase() *PassConfig_PassConfigDelegateBase
}

// PassConfig_PassConfigDelegateBase holds PassConfigDelegate's fields and
// base implementations.
type PassConfig_PassConfigDelegateBase struct {
	PassConfigBase

	delegate PassConfig
}

// NewPassConfig_PassConfigDelegate allocates a PassConfigDelegate.
//
// Java: PassConfig.PassConfigDelegate(PassConfig)
func NewPassConfig_PassConfigDelegate(delegate PassConfig) *PassConfig_PassConfigDelegateBase {
	d := &PassConfig_PassConfigDelegateBase{}
	d.InitPassConfigDelegate(d, delegate)
	return d
}

// InitPassConfigDelegate initializes a subclass embedding PassConfig_PassConfigDelegateBase.
func (d *PassConfig_PassConfigDelegateBase) InitPassConfigDelegate(self PassConfig, delegate PassConfig) {
	d.self = self
	d.InitPassConfig_PassConfigDelegateBase(delegate)
}

// InitPassConfig_PassConfigDelegateBase is the PassConfigDelegate
// constructor body (super(delegate.options) included).
//
// Java: PassConfig.PassConfigDelegate(PassConfig)
func (d *PassConfig_PassConfigDelegateBase) InitPassConfig_PassConfigDelegateBase(delegate PassConfig) {
	d.InitPassConfigBase(delegate.GetPassConfigBase().options)
	d.delegate = delegate
}

// GetPassConfig_PassConfigDelegateBase gives access to the delegate's
// fields through the interface.
func (d *PassConfig_PassConfigDelegateBase) GetPassConfig_PassConfigDelegateBase() *PassConfig_PassConfigDelegateBase {
	return d
}

// Java: PassConfig.PassConfigDelegate.getWhitespaceOnlyPasses()
func (d *PassConfig_PassConfigDelegateBase) GetWhitespaceOnlyPasses() *PassListBuilder {
	return d.delegate.GetWhitespaceOnlyPasses()
}

// Java: PassConfig.PassConfigDelegate.getChecks()
func (d *PassConfig_PassConfigDelegateBase) GetChecks() *PassListBuilder {
	return d.delegate.GetChecks()
}

// Java: PassConfig.PassConfigDelegate.getOptimizations()
func (d *PassConfig_PassConfigDelegateBase) GetOptimizations() *PassListBuilder {
	return d.delegate.GetOptimizations()
}

// Java: PassConfig.PassConfigDelegate.getFinalizations()
func (d *PassConfig_PassConfigDelegateBase) GetFinalizations() *PassListBuilder {
	return d.delegate.GetFinalizations()
}

// Java: PassConfig.PassConfigDelegate.getTranspileOnlyPasses()
func (d *PassConfig_PassConfigDelegateBase) GetTranspileOnlyPasses() *PassListBuilder {
	return d.delegate.GetTranspileOnlyPasses()
}
