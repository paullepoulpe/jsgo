// Port of com/google/javascript/jscomp/ConformancePassConfig.java @ closure-compiler v20260415

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// ConformancePassConfig runs only the user-supplied conformance checks and any
// earlier passes required by conformance.
type ConformancePassConfig struct {
	PassConfig_PassConfigDelegateBase

	delegate PassConfig
}

// Java: ConformancePassConfig(PassConfig)
func NewConformancePassConfig(delegate PassConfig) *ConformancePassConfig {
	c := &ConformancePassConfig{}
	c.self = c
	c.InitPassConfig_PassConfigDelegateBase(delegate)
	c.delegate = delegate
	return c
}

// Java: ConformancePassConfig.getChecks()
func (c *ConformancePassConfig) GetChecks() *PassListBuilder {
	fromDelegate := c.delegate.GetChecks()

	// Return every check up to and including the "checkConformance" check. Return empty list if
	// "checkConformance" not found. This list may include some unnecessary checks that run before
	// conformance. However, there's no other reliable way to find a list of all the passes
	// that conformance depends on.
	passes := NewPassListBuilder(c.options)
	passes.AddAllUpTo(fromDelegate, PassNames_CHECK_CONFORMANCE)
	return passes
}

// Java: ConformancePassConfig.getOptimizations()
func (c *ConformancePassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}
