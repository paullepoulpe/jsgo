// Port of com/google/javascript/jscomp/PassListBuilder.java
// @ closure-compiler v20260415

package jscomp

import "fmt"

// PassListBuilder accumulates PassFactory instances through a central
// mechanism that enforces each factory's condition before adding it.
// Mirrors Java's PassListBuilder final class.
type PassListBuilder struct {
	options *CompilerOptions
	passes  []*PassFactory
}

// NewPassListBuilder returns a PassListBuilder backed by the given options.
// Mirrors Java's PassListBuilder(CompilerOptions) constructor.
func NewPassListBuilder(options *CompilerOptions) *PassListBuilder {
	return &PassListBuilder{options: options}
}

// Build returns the accumulated list of pass factories as a slice.
// Mirrors Java's PassListBuilder.build() → ImmutableList.
func (b *PassListBuilder) Build() []*PassFactory {
	out := make([]*PassFactory, len(b.passes))
	copy(out, b.passes)
	return out
}

// AddAll appends all passes from another PassListBuilder.
// Mirrors Java's PassListBuilder.addAll(PassListBuilder).
func (b *PassListBuilder) AddAll(other *PassListBuilder) {
	b.passes = append(b.passes, other.Build()...)
}

// MaybeAdd adds the factory only if its condition evaluates to true.
// Mirrors Java's PassListBuilder.maybeAdd.
func (b *PassListBuilder) MaybeAdd(factory *PassFactory) {
	if factory.Condition == nil || factory.Condition(b.options) {
		b.passes = append(b.passes, factory)
	}
}

// Add unconditionally appends the factory.
func (b *PassListBuilder) Add(factory *PassFactory) {
	b.passes = append(b.passes, factory)
}

// AddBefore inserts the factory before the factory with the given name.
// Mirrors Java's PassListBuilder.addBefore.
func (b *PassListBuilder) AddBefore(factory *PassFactory, passName string) {
	idx := b.findIndexByName(passName)
	b.passes = append(b.passes, nil)
	copy(b.passes[idx+1:], b.passes[idx:])
	b.passes[idx] = factory
}

// AddAfter inserts the factory after the factory with the given name.
// Mirrors Java's PassListBuilder.addAfter.
func (b *PassListBuilder) AddAfter(factory *PassFactory, passName string) {
	idx := b.findIndexByName(passName)
	b.passes = append(b.passes, nil)
	copy(b.passes[idx+2:], b.passes[idx+1:])
	b.passes[idx+1] = factory
}

// AddAllUpTo appends passes from other up to and including passName.
// Deprecated in Java; preserved for API parity.
func (b *PassListBuilder) AddAllUpTo(other *PassListBuilder, passName string) {
	otherPasses := other.Build()
	idx := -1
	for i, p := range otherPasses {
		if p.Name == passName {
			idx = i
			break
		}
	}
	b.passes = append(b.passes, otherPasses[:idx+1]...)
}

// FindByName returns the factory with the given name. Panics if not found.
// Mirrors Java's PassListBuilder.findByName.
func (b *PassListBuilder) FindByName(name string) *PassFactory {
	for _, p := range b.passes {
		if p.Name == name {
			return p
		}
	}
	panic(fmt.Sprintf("No factory named %q in the factory list", name))
}

// RemoveByName removes the first factory with the given name.
// Deprecated in Java; preserved for API parity.
func (b *PassListBuilder) RemoveByName(name string) {
	for i, p := range b.passes {
		if p.Name == name {
			b.passes = append(b.passes[:i], b.passes[i+1:]...)
			return
		}
	}
}

// Contains reports whether the given factory is in the list.
// Mirrors Java's PassListBuilder.contains.
func (b *PassListBuilder) Contains(factory *PassFactory) bool {
	for _, p := range b.passes {
		if p == factory {
			return true
		}
	}
	return false
}

// AssertAllOneTimePasses panics if any pass is marked as a loop pass.
// Mirrors Java's PassListBuilder.assertAllOneTimePasses.
func (b *PassListBuilder) AssertAllOneTimePasses() {
	for _, p := range b.passes {
		if p.RunInFixedPointLoop {
			panic("expected one-time pass but got loop pass: " + p.Name)
		}
	}
}

// AssertAllLoopablePasses panics if any pass is not marked as a loop pass.
// Mirrors Java's PassListBuilder.assertAllLoopablePasses.
func (b *PassListBuilder) AssertAllLoopablePasses() {
	for _, p := range b.passes {
		if !p.RunInFixedPointLoop {
			panic("expected loop pass but got one-time pass: " + p.Name)
		}
	}
}

// AssertPassOrder panics if pass1 is present and ordered after pass2.
// Mirrors Java's PassListBuilder.assertPassOrder.
func (b *PassListBuilder) AssertPassOrder(pass1, pass2 *PassFactory, msg string) {
	idx1 := -1
	idx2 := -1
	for i, p := range b.passes {
		if p == pass1 {
			idx1 = i
		}
		if p == pass2 {
			idx2 = i
		}
	}
	if idx1 != -1 && idx2 != -1 && idx1 >= idx2 {
		panic(msg)
	}
}

// findIndexByName returns the index of the factory with the given name.
// Panics if not found.
func (b *PassListBuilder) findIndexByName(name string) int {
	for i, p := range b.passes {
		if p.Name == name {
			return i
		}
	}
	panic(fmt.Sprintf("No factory named %q in the factory list", name))
}
