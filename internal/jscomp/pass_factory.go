// Port of com/google/javascript/jscomp/PassFactory.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoValue_PassFactory.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PassFactory creates JSCompiler passes and holds pass metadata (name,
// preconditions, and whether the pass runs in a fixed-point loop).
type PassFactory struct {
	name                string
	condition           func(options *CompilerOptions) bool
	preconditionCheck   func(options *CompilerOptions) *PassFactory_PreconditionResult
	runInFixedPointLoop bool
	internalFactory     func(compiler AbstractCompiler) CompilerPass
}

func passFactoryAlwaysEnabled(*CompilerOptions) bool { return true }

func passFactoryAlwaysValid(*CompilerOptions) *PassFactory_PreconditionResult {
	return PassFactory_PreconditionResult_SUCCESS
}

func makePassFactory(name string, factory func(AbstractCompiler) CompilerPass) *PassFactory {
	return &PassFactory{
		name:              name,
		condition:         passFactoryAlwaysEnabled,
		preconditionCheck: passFactoryAlwaysValid,
		internalFactory:   factory,
	}
}

func makeLoopPassFactory(name string, factory func(AbstractCompiler) CompilerPass) *PassFactory {
	return &PassFactory{
		name:                name,
		condition:           passFactoryAlwaysEnabled,
		preconditionCheck:   passFactoryAlwaysValid,
		runInFixedPointLoop: true,
		internalFactory:     factory,
	}
}

func makeCheckedPassFactory(
	name string,
	precond func(*CompilerOptions) *PassFactory_PreconditionResult,
	loop bool,
	factory func(AbstractCompiler) CompilerPass,
) *PassFactory {
	return &PassFactory{
		name:                name,
		condition:           passFactoryAlwaysEnabled,
		preconditionCheck:   precond,
		runInFixedPointLoop: loop,
		internalFactory:     factory,
	}
}

// GetName returns the name of the pass as it will appear in logs.
func (f *PassFactory) GetName() string { return f.name }

func (f *PassFactory) GetCondition() func(options *CompilerOptions) bool { return f.condition }

// PassFactory_PreconditionResult is the result of a pass precondition check.
type PassFactory_PreconditionResult struct {
	success bool
	message *string
}

func NewPassFactory_PreconditionResult(success bool, message *string) *PassFactory_PreconditionResult {
	return &PassFactory_PreconditionResult{success: success, message: message}
}

var PassFactory_PreconditionResult_SUCCESS = NewPassFactory_PreconditionResult(true, nil)

func (r *PassFactory_PreconditionResult) Success() bool    { return r.success }
func (r *PassFactory_PreconditionResult) Message() *string { return r.message }

// GetPreconditionCheck runs before this pass runs to validate preconditions.
func (f *PassFactory) GetPreconditionCheck() func(options *CompilerOptions) *PassFactory_PreconditionResult {
	return f.preconditionCheck
}

// IsRunInFixedPointLoop reports whether this factory runs inside a PhaseOptimizer loop.
func (f *PassFactory) IsRunInFixedPointLoop() bool { return f.runInFixedPointLoop }

// GetInternalFactory returns the underlying pass constructor function.
func (f *PassFactory) GetInternalFactory() func(compiler AbstractCompiler) CompilerPass {
	return f.internalFactory
}

func (f *PassFactory) String() string {
	return "PassFactory{" + f.name + "}"
}

func (f *PassFactory) Equals(o any) bool {
	that, ok := o.(*PassFactory)
	if !ok || that == nil {
		return false
	}
	if that == f {
		return true
	}
	return f.name == that.name &&
		f.runInFixedPointLoop == that.runInFixedPointLoop &&
		jutil.ObjectsEquals(f.condition, that.condition) &&
		jutil.ObjectsEquals(f.preconditionCheck, that.preconditionCheck) &&
		jutil.ObjectsEquals(f.internalFactory, that.internalFactory)
}

func (f *PassFactory) ToBuilder() *PassFactory_Builder {
	return &PassFactory_Builder{
		name:                f.name,
		condition:           f.condition,
		preconditionCheck:   f.preconditionCheck,
		runInFixedPointLoop: f.runInFixedPointLoop,
		internalFactory:     f.internalFactory,
	}
}

// PassFactory_Builder builds a PassFactory.
type PassFactory_Builder struct {
	name                string
	condition           func(options *CompilerOptions) bool
	preconditionCheck   func(options *CompilerOptions) *PassFactory_PreconditionResult
	runInFixedPointLoop bool
	internalFactory     func(compiler AbstractCompiler) CompilerPass
}

func PassFactory_Builder_() *PassFactory_Builder {
	return &PassFactory_Builder{
		condition:         passFactoryAlwaysEnabled,
		preconditionCheck: passFactoryAlwaysValid,
	}
}

func (b *PassFactory_Builder) SetName(name string) *PassFactory_Builder {
	b.name = name
	return b
}

func (b *PassFactory_Builder) SetCondition(condition func(options *CompilerOptions) bool) *PassFactory_Builder {
	if condition == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null condition"))
	}
	b.condition = condition
	return b
}

func (b *PassFactory_Builder) SetPreconditionCheck(preconditionCheck func(options *CompilerOptions) *PassFactory_PreconditionResult) *PassFactory_Builder {
	if preconditionCheck == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null preconditionCheck"))
	}
	b.preconditionCheck = preconditionCheck
	return b
}

func (b *PassFactory_Builder) SetRunInFixedPointLoop(runInFixedPointLoop bool) *PassFactory_Builder {
	b.runInFixedPointLoop = runInFixedPointLoop
	return b
}

func (b *PassFactory_Builder) SetInternalFactory(internalFactory func(compiler AbstractCompiler) CompilerPass) *PassFactory_Builder {
	if internalFactory == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null internalFactory"))
	}
	b.internalFactory = internalFactory
	return b
}

func (b *PassFactory_Builder) Build() *PassFactory {
	if b.condition == nil || b.preconditionCheck == nil || b.internalFactory == nil || b.name == "" {
		b.failBuild()
	}
	return &PassFactory{
		name:                b.name,
		condition:           b.condition,
		preconditionCheck:   b.preconditionCheck,
		runInFixedPointLoop: b.runInFixedPointLoop,
		internalFactory:     b.internalFactory,
	}
}

func (b *PassFactory_Builder) failBuild() {
	if b.condition == nil || b.preconditionCheck == nil || b.internalFactory == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required PassFactory properties"))
	}
	jutil.CheckState(b.name != "")
}

func (f *PassFactory) ValidatePreconditions(options *CompilerOptions) {
	canRun := f.preconditionCheck(options)
	if !canRun.success {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.StringFormat("Precondition for pass %s failed: %s", f.name, canRun.message)))
	}
}

// PassFactory_CreateEmptyPass creates a no-op pass that can only run once.
// Used to break up loops and mark pipeline boundaries.
func PassFactory_CreateEmptyPass(name string) *PassFactory {
	if name == "" {
		jutil.CheckState(false)
	}
	return makePassFactory(name, passFactoryEmptyInternal)
}

type passFactory_emptyPass struct{}

var passFactory_emptyPassInstance = &passFactory_emptyPass{}

func passFactoryEmptyInternal(AbstractCompiler) CompilerPass { return passFactory_emptyPassInstance }

func (*passFactory_emptyPass) Process(externs rhino.Node, root rhino.Node) {}

// Create creates a new compiler pass to be run.
func (f *PassFactory) Create(compiler AbstractCompiler) CompilerPass {
	return f.internalFactory(compiler)
}
