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
			out := NewFastGzipOutputStream(file)
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

func protoPackedInt32PayloadSize(xs []int32) int {
	n := 0
	for _, x := range xs {
		n += protoVarintLen(uint64(int64(x)))
	}
	return n
}

func appendPackedInt32Field(b []byte, tagByte byte, xs []int32) []byte {
	start := len(b)
	b = append(b, tagByte, 0)
	for _, x := range xs {
		b = protoAppendInt32(b, x)
	}
	n := len(b) - (start + 2)
	if n < 0x80 {
		b[start+1] = byte(n)
		return b
	}
	vlen := protoVarintLen(uint64(n))
	var pad [4]byte
	b = append(b, pad[:vlen-1]...)
	copy(b[start+1+vlen:], b[start+2:start+2+n])
	protoPutVarint(b[start+1:], uint64(n))
	return b
}

func ensureTypePoolMemoizedSizes(tp *TypePool) int {
	if tp.frozen && tp.memoizedSize >= 0 {
		return tp.memoizedSize
	}
	tpSize := 0
	for _, t := range tp.type_ {
		s := fastTypeProtoSerializedSize(t)
		tpSize += 1 + protoVarintLen(uint64(s)) + s
	}
	for _, e := range tp.disambiguationEdges_ {
		s := fastSubtypingEdgeSerializedSize(e)
		tpSize += 1 + protoVarintLen(uint64(s)) + s
	}
	if tp.debugInfo_ != nil {
		s := tp.debugInfo_.serializedSize()
		tpSize += 1 + protoVarintLen(uint64(s)) + s
	}
	if len(tp.unknownFields) > 0 {
		tpSize += len(protoCanonicalUnknownFields(tp.unknownFields))
	}
	if tp.frozen {
		tp.memoizedSize = tpSize
	}
	return tpSize
}

func fastObjectTypeProtoSerializedSize(m *ObjectTypeProto) int {
	if m.frozen && m.memoizedSize >= 0 {
		return m.memoizedSize
	}
	size := 0
	if m.isInvalidating_ {
		size += 2
	}
	if len(m.uuid_) != 0 {
		size += 1 + protoVarintLen(uint64(len(m.uuid_))) + len(m.uuid_)
	}
	if m.propertiesKeepOriginalName_ {
		size += 2
	}
	if len(m.instanceType_) > 0 {
		p := protoPackedInt32PayloadSize(m.instanceType_)
		size += 1 + protoVarintLen(uint64(p)) + p
	}
	if len(m.prototype_) > 0 {
		p := protoPackedInt32PayloadSize(m.prototype_)
		size += 1 + protoVarintLen(uint64(p)) + p
	}
	if m.markedConstructor_ {
		size += 2
	}
	if len(m.ownProperty_) > 0 {
		p := protoPackedInt32PayloadSize(m.ownProperty_)
		size += 1 + protoVarintLen(uint64(p)) + p
	}
	if m.closureAssert_ {
		size += 2
	}
	if m.debugInfo_ != nil {
		s := m.debugInfo_.serializedSize()
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	if len(m.unknownFields) > 0 {
		size += len(protoCanonicalUnknownFields(m.unknownFields))
	}
	if m.frozen {
		m.memoizedSize = size
	}
	return size
}

func fastUnionTypeProtoSerializedSize(m *UnionTypeProto) int {
	if m.frozen && m.memoizedSize >= 0 {
		return m.memoizedSize
	}
	size := 0
	if len(m.unionMember_) > 0 {
		p := protoPackedInt32PayloadSize(m.unionMember_)
		size += 1 + protoVarintLen(uint64(p)) + p
	}
	if len(m.unknownFields) > 0 {
		size += len(protoCanonicalUnknownFields(m.unknownFields))
	}
	if m.frozen {
		m.memoizedSize = size
	}
	return size
}

func fastTypeProtoSerializedSize(m *TypeProto) int {
	if m.frozen && m.memoizedSize >= 0 {
		return m.memoizedSize
	}
	size := 0
	if m.kindCase_ == 2 {
		s := fastObjectTypeProtoSerializedSize(m.object_)
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	if m.kindCase_ == 4 {
		s := fastUnionTypeProtoSerializedSize(m.union_)
		size += 1 + protoVarintLen(uint64(s)) + s
	}
	if len(m.unknownFields) > 0 {
		size += len(protoCanonicalUnknownFields(m.unknownFields))
	}
	if m.frozen {
		m.memoizedSize = size
	}
	return size
}

func fastSubtypingEdgeSerializedSize(m *SubtypingEdge) int {
	if m.frozen && m.memoizedSize >= 0 {
		return m.memoizedSize
	}
	size := 0
	if m.subtype_ != 0 {
		size += 1 + protoVarintLen(uint64(int64(m.subtype_)))
	}
	if m.supertype_ != 0 {
		size += 1 + protoVarintLen(uint64(int64(m.supertype_)))
	}
	if len(m.unknownFields) > 0 {
		size += len(protoCanonicalUnknownFields(m.unknownFields))
	}
	if m.frozen {
		m.memoizedSize = size
	}
	return size
}

func fastAppendObjectTypeProto(b []byte, m *ObjectTypeProto) []byte {
	if m.isInvalidating_ {
		b = append(b, 16, 1)
	}
	if len(m.uuid_) != 0 {
		b = append(b, 26)
		b = protoAppendBytes(b, m.uuid_)
	}
	if m.propertiesKeepOriginalName_ {
		b = append(b, 32, 1)
	}
	if len(m.instanceType_) > 0 {
		b = appendPackedInt32Field(b, 42, m.instanceType_)
	}
	if len(m.prototype_) > 0 {
		b = appendPackedInt32Field(b, 50, m.prototype_)
	}
	if m.markedConstructor_ {
		b = append(b, 56, 1)
	}
	if len(m.ownProperty_) > 0 {
		b = appendPackedInt32Field(b, 66, m.ownProperty_)
	}
	if m.closureAssert_ {
		b = append(b, 72, 1)
	}
	if m.debugInfo_ != nil {
		b = append(b, 122)
		b = protoAppendVarint(b, uint64(m.debugInfo_.serializedSize()))
		b = m.debugInfo_.appendTo(b)
	}
	if len(m.unknownFields) > 0 {
		b = append(b, protoCanonicalUnknownFields(m.unknownFields)...)
	}
	return b
}

func fastAppendTypeProto(b []byte, m *TypeProto) []byte {
	if m.kindCase_ == 2 {
		b = append(b, 18)
		b = protoAppendVarint(b, uint64(fastObjectTypeProtoSerializedSize(m.object_)))
		b = fastAppendObjectTypeProto(b, m.object_)
	}
	if m.kindCase_ == 4 {
		b = append(b, 34)
		b = protoAppendVarint(b, uint64(fastUnionTypeProtoSerializedSize(m.union_)))
		if len(m.union_.unionMember_) > 0 {
			b = appendPackedInt32Field(b, 10, m.union_.unionMember_)
		}
		if len(m.union_.unknownFields) > 0 {
			b = append(b, protoCanonicalUnknownFields(m.union_.unknownFields)...)
		}
	}
	if len(m.unknownFields) > 0 {
		b = append(b, protoCanonicalUnknownFields(m.unknownFields)...)
	}
	return b
}

// writeTypedAstListItemTo streams a single TypedAst wrapped in a TypedAst.List
// item (field 1) to out without per-message packed slice allocations.
func writeTypedAstListItemTo(m *TypedAst, out io.Writer) error {
	if m.typePool_ != nil {
		ensureTypePoolMemoizedSizes(m.typePool_)
	}
	if m.stringPool_ != nil && m.stringPool_.frozen && m.stringPool_.memoizedSize < 0 {
		spSize := 0
		for _, x := range m.stringPool_.strings_ {
			spSize += 1 + protoVarintLen(uint64(len(x))) + len(x)
		}
		if m.stringPool_.maxLength_ != 0 {
			spSize += 1 + protoVarintLen(uint64(int64(m.stringPool_.maxLength_)))
		}
		if len(m.stringPool_.unknownFields) > 0 {
			spSize += len(protoCanonicalUnknownFields(m.stringPool_.unknownFields))
		}
		m.stringPool_.memoizedSize = spSize
	}
	for _, x := range m.externAst_ {
		if x.frozen && x.memoizedSize < 0 {
			if len(x.sourceMappingUrl_) == 0 && len(x.unknownFields) == 0 {
				s := 0
				if len(x.script_) != 0 {
					s += 1 + protoVarintLen(uint64(len(x.script_))) + len(x.script_)
				}
				if x.sourceFile_ != 0 {
					s += 1 + protoVarintLen(uint64(uint32(x.sourceFile_)))
				}
				x.memoizedSize = s
			} else {
				x.serializedSize()
			}
		}
	}
	for _, x := range m.codeAst_ {
		if x.frozen && x.memoizedSize < 0 {
			if len(x.sourceMappingUrl_) == 0 && len(x.unknownFields) == 0 {
				s := 0
				if len(x.script_) != 0 {
					s += 1 + protoVarintLen(uint64(len(x.script_))) + len(x.script_)
				}
				if x.sourceFile_ != 0 {
					s += 1 + protoVarintLen(uint64(uint32(x.sourceFile_)))
				}
				x.memoizedSize = s
			} else {
				x.serializedSize()
			}
		}
	}
	buf := make([]byte, 0, 64*1024)
	flush := func() error {
		if len(buf) == 0 {
			return nil
		}
		_, err := out.Write(buf)
		buf = buf[:0]
		return err
	}
	buf = append(buf, 10)
	buf = protoAppendVarint(buf, uint64(m.serializedSize()))
	if m.typePool_ != nil {
		buf = append(buf, 10)
		buf = protoAppendVarint(buf, uint64(m.typePool_.serializedSize()))
		for _, x := range m.typePool_.type_ {
			buf = append(buf, 10)
			buf = protoAppendVarint(buf, uint64(fastTypeProtoSerializedSize(x)))
			buf = fastAppendTypeProto(buf, x)
			if len(buf) >= 64*1024 {
				if err := flush(); err != nil {
					return err
				}
			}
		}
		for _, x := range m.typePool_.disambiguationEdges_ {
			buf = append(buf, 18)
			buf = protoAppendVarint(buf, uint64(fastSubtypingEdgeSerializedSize(x)))
			if x.subtype_ != 0 {
				buf = append(buf, 8)
				buf = protoAppendInt32(buf, x.subtype_)
			}
			if x.supertype_ != 0 {
				buf = append(buf, 16)
				buf = protoAppendInt32(buf, x.supertype_)
			}
			if len(x.unknownFields) > 0 {
				buf = append(buf, protoCanonicalUnknownFields(x.unknownFields)...)
			}
			if len(buf) >= 64*1024 {
				if err := flush(); err != nil {
					return err
				}
			}
		}
		if m.typePool_.debugInfo_ != nil {
			buf = append(buf, 26)
			buf = protoAppendVarint(buf, uint64(m.typePool_.debugInfo_.serializedSize()))
			buf = m.typePool_.debugInfo_.appendTo(buf)
		}
		if len(m.typePool_.unknownFields) > 0 {
			buf = append(buf, protoCanonicalUnknownFields(m.typePool_.unknownFields)...)
		}
		if err := flush(); err != nil {
			return err
		}
	}
	if m.stringPool_ != nil {
		buf = append(buf, 18)
		buf = protoAppendVarint(buf, uint64(m.stringPool_.serializedSize()))
		for _, x := range m.stringPool_.strings_ {
			buf = append(buf, 10)
			buf = protoAppendBytes(buf, x)
			if len(buf) >= 64*1024 {
				if err := flush(); err != nil {
					return err
				}
			}
		}
		if m.stringPool_.maxLength_ != 0 {
			buf = append(buf, 16)
			buf = protoAppendInt32(buf, m.stringPool_.maxLength_)
		}
		if len(m.stringPool_.unknownFields) > 0 {
			buf = append(buf, protoCanonicalUnknownFields(m.stringPool_.unknownFields)...)
		}
		if err := flush(); err != nil {
			return err
		}
	}
	for _, x := range m.externAst_ {
		buf = append(buf, 26)
		buf = protoAppendVarint(buf, uint64(x.serializedSize()))
		if len(x.sourceMappingUrl_) == 0 && len(x.unknownFields) == 0 {
			if len(x.script_) != 0 {
				buf = append(buf, 10)
				buf = protoAppendBytes(buf, x.script_)
			}
			if x.sourceFile_ != 0 {
				buf = append(buf, 16)
				buf = protoAppendUint32(buf, x.sourceFile_)
			}
		} else {
			buf = x.appendTo(buf)
		}
		if len(buf) >= 64*1024 {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	for _, x := range m.codeAst_ {
		buf = append(buf, 34)
		buf = protoAppendVarint(buf, uint64(x.serializedSize()))
		if len(x.sourceMappingUrl_) == 0 && len(x.unknownFields) == 0 {
			if len(x.script_) != 0 {
				buf = append(buf, 10)
				buf = protoAppendBytes(buf, x.script_)
			}
			if x.sourceFile_ != 0 {
				buf = append(buf, 16)
				buf = protoAppendUint32(buf, x.sourceFile_)
			}
		} else {
			buf = x.appendTo(buf)
		}
		if len(buf) >= 64*1024 {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	if m.sourceFilePool_ != nil {
		buf = append(buf, 42)
		buf = protoAppendVarint(buf, uint64(m.sourceFilePool_.serializedSize()))
		for _, x := range m.sourceFilePool_.sourceFile_ {
			buf = append(buf, 10)
			buf = protoAppendVarint(buf, uint64(x.serializedSize()))
			buf = x.appendTo(buf)
			if len(buf) >= 64*1024 {
				if err := flush(); err != nil {
					return err
				}
			}
		}
		if len(m.sourceFilePool_.unknownFields) > 0 {
			buf = append(buf, protoCanonicalUnknownFields(m.sourceFilePool_.unknownFields)...)
		}
	}
	if m.externsSummary_ != nil {
		buf = append(buf, 50)
		buf = protoAppendVarint(buf, uint64(m.externsSummary_.serializedSize()))
		buf = m.externsSummary_.appendTo(buf)
	}
	for _, x := range m.runtimeLibraryToInject_ {
		buf = append(buf, 58)
		buf = protoAppendString(buf, x)
	}
	if len(m.unknownFields) > 0 {
		buf = append(buf, protoCanonicalUnknownFields(m.unknownFields)...)
	}
	return flush()
}

// Java: SerializeTypedAstPass.process(Node, Node)
func (p *SerializeTypedAstPass) Process(externs rhino.Node, root rhino.Node) {
	NewRemoveCastNodes(p.compiler).Process(externs, root)
	serializer := NewTypedAstSerializer(p.compiler, p.serializationOptions)
	ast := serializer.SerializeRoots(externs, root)
	p.consumer(ast)
}
