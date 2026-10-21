// Port of com/google/javascript/jscomp/serialization/SerializeTypedAstPass.java @ closure-compiler v20261015

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

import (
	"io"
	"os"

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
		if err := ast.writeAsListItemTo(out); err != nil {
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
			out := NewFastGzipOutputStream(file)
			defer out.Close()
			if err := ast.writeAsListItemTo(out); err != nil {
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
func (p *SerializeTypedAstPass) Process(externs *rhino.Node, root *rhino.Node) {
	NewRemoveCastNodes(p.compiler).Process(externs, root)
	serializer := NewTypedAstSerializer(p.compiler, p.serializationOptions)
	ast := serializer.SerializeRoots(externs, root)
	p.consumer(ast)
}
