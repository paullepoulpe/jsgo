// Port of com/google/javascript/jscomp/PassListBuilder.java @ closure-compiler v20260415

/*
 * Copyright 2022 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// PassListBuilder exists to funnel the list of passes through a central
// mechanism to ensure that the PassFactory can define and enforce conditions
// that much be met in order for the pass they create to be added to the
// compilation process.
type PassListBuilder struct {
	options *CompilerOptions
	passes  []*PassFactory
}

// Java: PassListBuilder(CompilerOptions)
func NewPassListBuilder(options *CompilerOptions) *PassListBuilder {
	b := &PassListBuilder{}
	b.options = options
	return b
}

// Java: PassListBuilder.build()
func (b *PassListBuilder) Build() []*PassFactory {
	return append([]*PassFactory(nil), b.passes...)
}

// Java: PassListBuilder.addAll(PassListBuilder)
func (b *PassListBuilder) AddAll(other *PassListBuilder) {
	// Don't bother to check conditions here as we assume that they were
	// checked when initially added to the other PassListBuilder.
	b.passes = append(b.passes, other.Build()...)
}

// MaybeAdd adds only if the PassFactory condition evaluates to true.
//
// Java: PassListBuilder.maybeAdd(PassFactory)
func (b *PassListBuilder) MaybeAdd(factory *PassFactory) {
	if factory.GetCondition()(b.options) {
		b.passes = append(b.passes, factory)
	}
}

// AddBefore inserts the given pass factory before the factory of the given
// name. Throws if the specified pass is not present
//
// Java: PassListBuilder.addBefore(PassFactory, String)
func (b *PassListBuilder) AddBefore(factory *PassFactory, passName string) {
	b.passes = passListBuilder_insert(b.passes, b.FindIndexByName(passName), factory)
}

// AddAfter inserts the given pass factory before the factory of the given
// name. Throws if the specified pass is not present
//
// Java: PassListBuilder.addAfter(PassFactory, String)
func (b *PassListBuilder) AddAfter(factory *PassFactory, passName string) {
	b.passes = passListBuilder_insert(b.passes, b.FindIndexByName(passName)+1, factory)
}

// passListBuilder_insert is List.add(int, E).
//
// Go-only: language-adapter — java.util.List.add(index, element) on a slice.
func passListBuilder_insert(passes []*PassFactory, index int, factory *PassFactory) []*PassFactory {
	if index < 0 || index > len(passes) {
		panic(jutil.NewIndexOutOfBoundsException(
			"Index: " + jutil.ToJavaString(index) + ", Size: " + jutil.ToJavaString(len(passes))))
	}
	passes = append(passes, nil)
	copy(passes[index+1:], passes[index:])
	passes[index] = factory
	return passes
}

// AddAllUpTo returns the list of pass up to and including the named pass,
// otherwise return an empty set of passes.
//
// Do not add new uses of this method.
//
// Deprecated: replace uses of this method with a more precise definition of
// a pass configuration.
//
// Java: PassListBuilder.addAllUpTo(PassListBuilder, String)
func (b *PassListBuilder) AddAllUpTo(other *PassListBuilder, passName string) {
	otherPasses := other.Build()

	index := -1
	for i := 0; i < len(otherPasses); i++ {
		if passName == otherPasses[i].GetName() {
			index = i
			break
		}
	}
	// Return every check up to and including the "checkConformance" check. Return empty list if
	// "checkConformance" not found. This list may include some unnecessary checks that run before
	// conformance. However, there's no other reliable way to find a list of all the passes
	// that conformance depends on.
	b.passes = append(b.passes, otherPasses[0:index+1]...)
}

// Java: PassListBuilder.findByName(String)
func (b *PassListBuilder) FindByName(name string) *PassFactory {
	for _, pass := range b.passes {
		if pass.GetName() == name {
			return pass
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("No factory named '" + name + "' in the factory list"))
}

// RemoveByName removes a pass with the specified name.
//
// Deprecated: Do not add new uses to the method. This method exists only
// for migration purposes. Removing passes piecemeal is very fragile.
//
// Java: PassListBuilder.removeByName(String)
func (b *PassListBuilder) RemoveByName(name string) {
	for i := 0; i < len(b.passes); i++ {
		if b.passes[i].GetName() == name {
			b.passes = append(b.passes[:i], b.passes[i+1:]...)
			return
		}
	}
}

// FindIndexByName throws an exception if no pass with the given name
// exists.
//
// Java: PassListBuilder.findIndexByName(String)
func (b *PassListBuilder) FindIndexByName(name string) int {
	for i := 0; i < len(b.passes); i++ {
		if b.passes[i].GetName() == name {
			return i
		}
	}

	panic(jutil.NewIllegalArgumentExceptionMsg("No factory named '" + name + "' in the factory list"))
}

// Java: PassListBuilder.contains(PassFactory)
func (b *PassListBuilder) Contains(factory *PassFactory) bool {
	return b.indexOf(factory) != -1
}

// indexOf is List.indexOf(Object), which compares with equals (PassFactory
// is an AutoValue).
//
// Go-only: language-adapter — java.util.List.indexOf on a slice.
func (b *PassListBuilder) indexOf(factory *PassFactory) int {
	for i, pass := range b.passes {
		if factory == nil {
			if pass == nil {
				return i
			}
		} else if factory.Equals(pass) {
			return i
		}
	}
	return -1
}

// AssertAllOneTimePasses verifies that all the passes are one-time passes.
//
// Java: PassListBuilder.assertAllOneTimePasses()
func (b *PassListBuilder) AssertAllOneTimePasses() {
	for _, pass := range b.passes {
		jutil.CheckState(!pass.IsRunInFixedPointLoop())
	}
}

// AssertAllLoopablePasses verifies that all the passes are multi-run passes.
//
// Java: PassListBuilder.assertAllLoopablePasses()
func (b *PassListBuilder) AssertAllLoopablePasses() {
	for _, pass := range b.passes {
		jutil.CheckState(pass.IsRunInFixedPointLoop())
	}
}

// AssertPassOrder asserts that if both PassFactory are present, pass1 is
// ordered before pass2.
//
// Java: PassListBuilder.assertPassOrder(PassFactory, PassFactory, String)
func (b *PassListBuilder) AssertPassOrder(pass1 *PassFactory, pass2 *PassFactory, msg string) {
	pass1Index := b.indexOf(pass1)
	pass2Index := b.indexOf(pass2)
	if pass1Index != -1 && pass2Index != -1 {
		jutil.CheckStateMsg(pass1Index < pass2Index, msg)
	}
}
