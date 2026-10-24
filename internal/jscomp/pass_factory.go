// Port of com/google/javascript/jscomp/PassFactory.java @ closure-compiler v20261015
// Port of com/google/javascript/jscomp/AutoValue_PassFactory.java @ closure-compiler v20261015

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

// PassFactory is a factory for creating JSCompiler passes based on the
// Options injected.
//
// Contains all meta-data about compiler passes (like whether it can be run
// multiple times, a human-readable name for logging, etc.).
//
// Go-only: language-adapter — an @AutoValue: the abstract class and the
// generated AutoValue_PassFactory (fields, equals, hashCode, toString) are
// one struct (PORTING.md §4). The java.util.function.Function fields are Go
// funcs; Java compares the lambdas by identity, which jutil.ObjectsEquals
// models for funcs by code pointer.
type PassFactory struct {
	name                string
	condition           func(options *CompilerOptions) bool
	preconditionCheck   func(options *CompilerOptions) *PassFactory_PreconditionResult
	runInFixedPointLoop bool
	internalFactory     func(compiler AbstractCompiler) CompilerPass
}

// Java: AutoValue_PassFactory(String, Function, Function, boolean, Function)
func newAutoValue_PassFactory(
	name string,
	condition func(options *CompilerOptions) bool,
	preconditionCheck func(options *CompilerOptions) *PassFactory_PreconditionResult,
	runInFixedPointLoop bool,
	internalFactory func(compiler AbstractCompiler) CompilerPass) *PassFactory {
	return &PassFactory{
		name:                name,
		condition:           condition,
		preconditionCheck:   preconditionCheck,
		runInFixedPointLoop: runInFixedPointLoop,
		internalFactory:     internalFactory,
	}
}

// GetName returns the name of the pass as it will appear in logs.
//
// Java: PassFactory.getName()
func (f *PassFactory) GetName() string {
	return f.name
}

// Java: PassFactory.getCondition()
func (f *PassFactory) GetCondition() func(options *CompilerOptions) bool {
	return f.condition
}

// PassFactory_PreconditionResult is the result of a precondition check.
//
// Go-only: language-adapter — a Java record: a struct with its accessors
// and value equality.
type PassFactory_PreconditionResult struct {
	success bool
	message *string
}

// Java: PassFactory.PreconditionResult(boolean, String)
func NewPassFactory_PreconditionResult(success bool, message *string) *PassFactory_PreconditionResult {
	return &PassFactory_PreconditionResult{success: success, message: message}
}

// PassFactory_PreconditionResult_SUCCESS is PreconditionResult.SUCCESS.
var PassFactory_PreconditionResult_SUCCESS = NewPassFactory_PreconditionResult(true, nil)

// Java: PassFactory.PreconditionResult.success()
func (r *PassFactory_PreconditionResult) Success() bool { return r.success }

// Java: PassFactory.PreconditionResult.message()
func (r *PassFactory_PreconditionResult) Message() *string { return r.message }

// Java: PassFactory.PreconditionResult.equals(Object)
func (r *PassFactory_PreconditionResult) Equals(o any) bool {
	that, ok := o.(*PassFactory_PreconditionResult)
	if !ok || that == nil {
		return false
	}
	return r.success == that.success && jutil.ObjectsEquals(r.message, that.message)
}

// Java: PassFactory.PreconditionResult.hashCode()
func (r *PassFactory_PreconditionResult) HashCode() int32 {
	// A record's hashCode combines its components as 31 * h + hash(c).
	var h int32
	if r.success {
		h = 1231
	} else {
		h = 1237
	}
	var m int32
	if r.message != nil {
		m = jutil.StringHashCode(*r.message)
	}
	return 31*h + m
}

// Java: PassFactory.PreconditionResult.toString()
func (r *PassFactory_PreconditionResult) String() string {
	return "PreconditionResult[success=" + jutil.ToJavaString(r.success) +
		", message=" + jutil.ToJavaString(r.message) + "]"
}

// GetPreconditionCheck runs before this pass runs to validate
// preconditions. Throws an exception if any are violated (unlike
// GetCondition.)
//
// Java: PassFactory.getPreconditionCheck()
func (f *PassFactory) GetPreconditionCheck() func(options *CompilerOptions) *PassFactory_PreconditionResult {
	return f.preconditionCheck
}

// IsRunInFixedPointLoop reports whether this factory must or must not
// appear in a PhaseOptimizer loop.
//
// Java: PassFactory.isRunInFixedPointLoop()
func (f *PassFactory) IsRunInFixedPointLoop() bool {
	return f.runInFixedPointLoop
}

// GetInternalFactory is a simple factory function for creating actual pass
// instances.
//
// Users should call Create(AbstractCompiler) rather than use this object
// directly.
//
// Java: PassFactory.getInternalFactory()
func (f *PassFactory) GetInternalFactory() func(compiler AbstractCompiler) CompilerPass {
	return f.internalFactory
}

// Java: AutoValue_PassFactory.toString()
func (f *PassFactory) String() string {
	return "PassFactory{" +
		"name=" + f.name + ", " +
		"condition=" + jutil.ToJavaString(f.condition) + ", " +
		"preconditionCheck=" + jutil.ToJavaString(f.preconditionCheck) + ", " +
		"runInFixedPointLoop=" + jutil.ToJavaString(f.runInFixedPointLoop) + ", " +
		"internalFactory=" + jutil.ToJavaString(f.internalFactory) +
		"}"
}

// Java: AutoValue_PassFactory.equals(Object)
func (f *PassFactory) Equals(o any) bool {
	that, ok := o.(*PassFactory)
	if ok && that == f {
		return true
	}
	if ok && that != nil {
		return f.name == that.GetName() &&
			jutil.ObjectsEquals(f.condition, that.GetCondition()) &&
			jutil.ObjectsEquals(f.preconditionCheck, that.GetPreconditionCheck()) &&
			f.runInFixedPointLoop == that.IsRunInFixedPointLoop() &&
			jutil.ObjectsEquals(f.internalFactory, that.GetInternalFactory())
	}
	return false
}

// Java: AutoValue_PassFactory.hashCode()
func (f *PassFactory) HashCode() int32 {
	h := int32(1)
	h *= 1000003
	h ^= jutil.StringHashCode(f.name)
	h *= 1000003
	h ^= jutil.ObjectHashCode(f.condition)
	h *= 1000003
	h ^= jutil.ObjectHashCode(f.preconditionCheck)
	h *= 1000003
	if f.runInFixedPointLoop {
		h ^= 1231
	} else {
		h ^= 1237
	}
	h *= 1000003
	h ^= jutil.ObjectHashCode(f.internalFactory)
	return h
}

// Java: AutoValue_PassFactory.toBuilder()
func (f *PassFactory) ToBuilder() *PassFactory_Builder {
	return newAutoValue_PassFactory_BuilderFrom(f)
}

// PassFactory_Builder is a builder for a PassFactory.
//
// Go-only: language-adapter — the abstract PassFactory.Builder and the
// generated AutoValue_PassFactory.Builder are one struct.
type PassFactory_Builder struct {
	name                *string
	condition           func(options *CompilerOptions) bool
	preconditionCheck   func(options *CompilerOptions) *PassFactory_PreconditionResult
	runInFixedPointLoop bool
	internalFactory     func(compiler AbstractCompiler) CompilerPass
	set_0               byte
}

// Java: AutoValue_PassFactory.Builder()
func newAutoValue_PassFactory_Builder() *PassFactory_Builder {
	return &PassFactory_Builder{}
}

// Java: AutoValue_PassFactory.Builder(PassFactory)
func newAutoValue_PassFactory_BuilderFrom(source *PassFactory) *PassFactory_Builder {
	name := source.GetName()
	return &PassFactory_Builder{
		name:                &name,
		condition:           source.GetCondition(),
		preconditionCheck:   source.GetPreconditionCheck(),
		runInFixedPointLoop: source.IsRunInFixedPointLoop(),
		internalFactory:     source.GetInternalFactory(),
		set_0:               1,
	}
}

// Java: AutoValue_PassFactory.Builder.setName(String)
func (b *PassFactory_Builder) SetName(name string) *PassFactory_Builder {
	// Go-only: language-adapter — a Go string cannot be null, so the
	// generated `if (name == null) throw new NullPointerException("Null name")`
	// cannot fire.
	b.name = &name
	return b
}

// Java: AutoValue_PassFactory.Builder.setCondition(Function)
func (b *PassFactory_Builder) SetCondition(condition func(options *CompilerOptions) bool) *PassFactory_Builder {
	if condition == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null condition"))
	}
	b.condition = condition
	return b
}

// Java: AutoValue_PassFactory.Builder.setPreconditionCheck(Function)
func (b *PassFactory_Builder) SetPreconditionCheck(preconditionCheck func(options *CompilerOptions) *PassFactory_PreconditionResult) *PassFactory_Builder {
	if preconditionCheck == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null preconditionCheck"))
	}
	b.preconditionCheck = preconditionCheck
	return b
}

// Java: AutoValue_PassFactory.Builder.setRunInFixedPointLoop(boolean)
func (b *PassFactory_Builder) SetRunInFixedPointLoop(runInFixedPointLoop bool) *PassFactory_Builder {
	b.runInFixedPointLoop = runInFixedPointLoop
	b.set_0 |= 1
	return b
}

// Java: AutoValue_PassFactory.Builder.setInternalFactory(Function)
func (b *PassFactory_Builder) SetInternalFactory(internalFactory func(compiler AbstractCompiler) CompilerPass) *PassFactory_Builder {
	if internalFactory == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null internalFactory"))
	}
	b.internalFactory = internalFactory
	return b
}

// Java: AutoValue_PassFactory.Builder.autoBuild()
func (b *PassFactory_Builder) AutoBuild() *PassFactory {
	if b.set_0 != 1 ||
		b.name == nil ||
		b.condition == nil ||
		b.preconditionCheck == nil ||
		b.internalFactory == nil {
		missing := jutil.NewStringBuilder()
		if b.name == nil {
			missing.Append(" name")
		}
		if b.condition == nil {
			missing.Append(" condition")
		}
		if b.preconditionCheck == nil {
			missing.Append(" preconditionCheck")
		}
		if (b.set_0 & 1) == 0 {
			missing.Append(" runInFixedPointLoop")
		}
		if b.internalFactory == nil {
			missing.Append(" internalFactory")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return newAutoValue_PassFactory(
		*b.name,
		b.condition,
		b.preconditionCheck,
		b.runInFixedPointLoop,
		b.internalFactory)
}

// Java: PassFactory.Builder.build()
func (b *PassFactory_Builder) Build() *PassFactory {
	result := b.AutoBuild()
	jutil.CheckState(result.GetName() != "")
	return result
}

// Java: PassFactory.builder()
func PassFactory_Builder_() *PassFactory_Builder {
	return newAutoValue_PassFactory_Builder().
		SetRunInFixedPointLoop(false).
		SetCondition(func(o *CompilerOptions) bool { return true }).
		SetPreconditionCheck(func(o *CompilerOptions) *PassFactory_PreconditionResult {
			return PassFactory_PreconditionResult_SUCCESS
		})
}

// Java: PassFactory.validatePreconditions(CompilerOptions)
func (f *PassFactory) ValidatePreconditions(options *CompilerOptions) {
	canRun := f.GetPreconditionCheck()(options)
	if !canRun.success {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.StringFormat("Precondition for pass %s failed: %s", f.GetName(), canRun.Message())))
	}
}

// PassFactory_CreateEmptyPass creates a no-op pass that can only run once.
// Used to break up loops.
//
// Java: PassFactory.createEmptyPass(String)
func PassFactory_CreateEmptyPass(name string) *PassFactory {
	return PassFactory_Builder_().
		SetName(name).
		SetInternalFactory(func(c AbstractCompiler) CompilerPass { return passFactory_emptyPassInstance }).
		Build()
}

// passFactory_emptyPass is the `(externs, root) -> {}` lambda of
// createEmptyPass.
type passFactory_emptyPass struct{}

// passFactory_emptyPassInstance: the lambda captures nothing, so Java
// evaluates it to a single instance.
var passFactory_emptyPassInstance = &passFactory_emptyPass{}

// Java: PassFactory.createEmptyPass(String) (the CompilerPass lambda)
func (*passFactory_emptyPass) Process(externs rhino.Node, root rhino.Node) {}

// Create creates a new compiler pass to be run.
//
// Java: PassFactory.create(AbstractCompiler)
func (f *PassFactory) Create(compiler AbstractCompiler) CompilerPass {
	return f.GetInternalFactory()(compiler)
}
