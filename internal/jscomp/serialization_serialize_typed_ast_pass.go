// Port of com/google/javascript/jscomp/serialization/SerializeTypedAstPass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/serialization/SerializationOptions.java @ closure-compiler v20261021
// (with the AutoBuilder-generated serialization/AutoBuilder_SerializationOptions_Builder.java)

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
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SerializeTypedAstPass is a compiler pass intended to serialize the types in
// the AST.
//
// Serialization means that the AST is converted to a proto representation
// and output to a file or stream. Deserialization means the reverse process.
//
// Go-only: language-adapter — java.util.function.Consumer<TypedAst> is a
// func(*TypedAst); java.io.OutputStream is an io.Writer and
// java.nio.file.Path a file name.
//
// Java: public final class SerializeTypedAstPass implements CompilerPass
type SerializeTypedAstPass struct {
	compiler             AbstractCompiler
	consumer             func(*TypedAst)
	serializationOptions *SerializationOptions
}

// Java: SerializeTypedAstPass(AbstractCompiler, Consumer<TypedAst>, SerializationOptions)
func NewSerializeTypedAstPass(
	compiler AbstractCompiler,
	astConsumer func(*TypedAst),
	serializationOptions *SerializationOptions) *SerializeTypedAstPass {
	p := &SerializeTypedAstPass{}
	p.compiler = compiler
	p.consumer = astConsumer
	p.serializationOptions = serializationOptions
	return p
}

// SerializeTypedAstPass_CreateFromOutputStream serializes a TypedAst to the
// given output stream.
//
// Unlike SerializeTypedAstPass_CreateFromPath, this method does not
// automatically gzip the TypedAST. The "out" parameter may or may not already
// be a GZIPOutputStream.
//
// Java: SerializeTypedAstPass.createFromOutputStream(AbstractCompiler, OutputStream, SerializationOptions)
func SerializeTypedAstPass_CreateFromOutputStream(
	c AbstractCompiler, out io.Writer, serializationOptions *SerializationOptions) *SerializeTypedAstPass {
	toOutputStream := func(ast *TypedAst) {
		// TypedAst.List.newBuilder().addTypedAsts(ast).build().writeTo(out)
		if err := writeTypedAstListItemTo(ast, out); err != nil {
			panic(jutil.NewIllegalArgumentExceptionWithMessageAndCause(
				"Cannot write to stream", jutil.NewIOException(err.Error())))
		}
	}
	return NewSerializeTypedAstPass(c, toOutputStream, serializationOptions)
}

// SerializeTypedAstPass_CreateFromPath serializes a gzipped TypedAst to the
// specified outputPath.
//
// Go-only: language-adapter — jutil's GZIPOutputStream compresses with
// Go's compress/flate, so the file's compressed bytes can differ from
// Java's; it decompresses to the same TypedAst.List bytes (PORTING.md §7).
//
// Java: SerializeTypedAstPass.createFromPath(AbstractCompiler, Path, SerializationOptions)
func SerializeTypedAstPass_CreateFromPath(
	compiler AbstractCompiler, outputPath *jutil.Path, serializationOptions *SerializationOptions) *SerializeTypedAstPass {
	toPath := func(ast *TypedAst) {
		jutil.Try(func() {
			// try (OutputStream out = new FastGzipOutputStream(Files.newOutputStream(outputPath)))
			file, err := os.Create(outputPath.String())
			if err != nil {
				panic(jutil.NewIOException(err.Error()))
			}
			numThreads := 1
			if opts := compiler.GetOptions(); opts != nil {
				numThreads = opts.GetNumParallelThreads()
			}
			out := newSaveStateGzipOutputStream(file, numThreads)
			defer out.Close()
			if err := writeTypedAstListItemTo(ast, out); err != nil {
				panic(jutil.NewIOException(err.Error()))
			}
		}, func(e jutil.IIOException) {
			panic(jutil.NewIllegalArgumentExceptionWithMessageAndCause(
				"Cannot create TypedAst output file", e))
		})
	}
	return NewSerializeTypedAstPass(compiler, toPath, serializationOptions)
}

// Java: SerializeTypedAstPass.process(Node, Node)
func (p *SerializeTypedAstPass) Process(externs rhino.Node, root rhino.Node) {
	NewRemoveCastNodes(p.compiler).Process(externs, root)
	serializer := NewTypedAstSerializer(p.compiler, p.serializationOptions)
	ast := serializer.SerializeRoots(externs, root)
	p.consumer(ast)
}

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

func serializationBoolHashCode(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

// Java: SerializationOptions.hashCode() (record: 31 * h + hash(component))
func (o *SerializationOptions) HashCode() int32 {
	lh := int32(1)
	for _, lib := range o.runtimeLibraries {
		lh = 31*lh + jutil.StringHashCode(lib)
	}
	var h int32
	h = 31*h + serializationBoolHashCode(o.includeDebugInfo)
	h = 31*h + serializationBoolHashCode(o.runValidation)
	h = 31*h + lh
	return h
}

// Java: SerializationOptions.toString() (record)
func (o *SerializationOptions) String() string {
	return "SerializationOptions[includeDebugInfo=" + strconv.FormatBool(o.includeDebugInfo) +
		", runValidation=" + strconv.FormatBool(o.runValidation) +
		", runtimeLibraries=[" + strings.Join(o.runtimeLibraries, ", ") + "]]"
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
