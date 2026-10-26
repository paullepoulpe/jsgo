// Port of com/google/javascript/jscomp/serialization/SerializationOptions.java @ closure-compiler v20261021
// (with the AutoBuilder-generated serialization/AutoBuilder_SerializationOptions_Builder.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SerializationOptions holds configuration options for serialization time.
type SerializationOptions struct {
	includeDebugInfo bool
	runValidation    bool
	runtimeLibraries []string
}

// Java: SerializationOptions(boolean, boolean, ImmutableList<String>)
func NewSerializationOptions(includeDebugInfo bool, runValidation bool, runtimeLibraries []string) *SerializationOptions {
	if runtimeLibraries == nil {
		panic(jutil.NewNullPointerExceptionMsg("runtimeLibraries"))
	}
	return &SerializationOptions{includeDebugInfo: includeDebugInfo, runValidation: runValidation, runtimeLibraries: runtimeLibraries}
}

// Java: SerializationOptions.includeDebugInfo()
func (o *SerializationOptions) IncludeDebugInfo() bool { return o.includeDebugInfo }

// Java: SerializationOptions.runValidation()
func (o *SerializationOptions) RunValidation() bool { return o.runValidation }

// Java: SerializationOptions.runtimeLibraries()
func (o *SerializationOptions) RuntimeLibraries() []string { return o.runtimeLibraries }

// SerializationOptions_Builder_ is SerializationOptions.builder().
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name SerializationOptions_Builder for the nested Builder type.
func SerializationOptions_Builder_() *SerializationOptions_Builder {
	return NewSerializationOptions_Builder().
		SetRunValidation(false).
		SetIncludeDebugInfo(false).
		SetRuntimeLibraries([]string{})
}

// Java: SerializationOptions.equals(Object) (record)
func (o *SerializationOptions) Equals(obj any) bool {
	that, ok := obj.(*SerializationOptions)
	if !ok || that == nil {
		return false
	}
	return o.includeDebugInfo == that.includeDebugInfo && o.runValidation == that.runValidation &&
		slices.Equal(o.runtimeLibraries, that.runtimeLibraries)
}

// Java: SerializationOptions.hashCode() (record: 31 * h + hash(component))
func (o *SerializationOptions) HashCode() int32 {
	lh := int32(1)
	for _, lib := range o.runtimeLibraries {
		lh = 31*lh + jutil.StringHashCode(lib)
	}
	var h int32
	h = 31*h + jutil.ObjectHashCode(o.includeDebugInfo)
	h = 31*h + jutil.ObjectHashCode(o.runValidation)
	h = 31*h + lh
	return h
}

// Java: SerializationOptions.toString() (record)
func (o *SerializationOptions) String() string {
	return "SerializationOptions[includeDebugInfo=" + jutil.ToJavaString(o.includeDebugInfo) +
		", runValidation=" + jutil.ToJavaString(o.runValidation) +
		", runtimeLibraries=" + jutil.CollectionToString(o.runtimeLibraries) + "]"
}

// SerializationOptions_Builder is SerializationOptions.Builder (the
// AutoBuilder-generated AutoBuilder_SerializationOptions_Builder).
type SerializationOptions_Builder struct {
	includeDebugInfo bool
	runValidation    bool
	runtimeLibraries []string
	set0             byte
}

// Java: AutoBuilder_SerializationOptions_Builder()
func NewSerializationOptions_Builder() *SerializationOptions_Builder {
	return &SerializationOptions_Builder{}
}

// Java: AutoBuilder_SerializationOptions_Builder(SerializationOptions)
func NewSerializationOptions_BuilderFrom(source *SerializationOptions) *SerializationOptions_Builder {
	return &SerializationOptions_Builder{
		includeDebugInfo: source.IncludeDebugInfo(),
		runValidation:    source.RunValidation(),
		runtimeLibraries: source.RuntimeLibraries(),
		set0:             3,
	}
}

// Java: SerializationOptions.Builder.setIncludeDebugInfo(boolean)
func (b *SerializationOptions_Builder) SetIncludeDebugInfo(includeDebugInfo bool) *SerializationOptions_Builder {
	b.includeDebugInfo = includeDebugInfo
	b.set0 |= 1
	return b
}

// Java: SerializationOptions.Builder.setRunValidation(boolean)
func (b *SerializationOptions_Builder) SetRunValidation(runValidation bool) *SerializationOptions_Builder {
	b.runValidation = runValidation
	b.set0 |= 2
	return b
}

// Java: SerializationOptions.Builder.setRuntimeLibraries(ImmutableList<String>)
func (b *SerializationOptions_Builder) SetRuntimeLibraries(runtimeLibraries []string) *SerializationOptions_Builder {
	if runtimeLibraries == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null runtimeLibraries"))
	}
	b.runtimeLibraries = runtimeLibraries
	return b
}

// Java: SerializationOptions.Builder.build()
func (b *SerializationOptions_Builder) Build() *SerializationOptions {
	if b.set0 != 3 || b.runtimeLibraries == nil {
		var missing strings.Builder
		if b.set0&1 == 0 {
			missing.WriteString(" includeDebugInfo")
		}
		if b.set0&2 == 0 {
			missing.WriteString(" runValidation")
		}
		if b.runtimeLibraries == nil {
			missing.WriteString(" runtimeLibraries")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return NewSerializationOptions(b.includeDebugInfo, b.runValidation, b.runtimeLibraries)
}
