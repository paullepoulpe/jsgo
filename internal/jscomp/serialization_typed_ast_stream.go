package jscomp

import "io"

// Streaming wire encoder for TypedAst.List items, avoiding intermediate
// packed-slice and message-buffer allocations during SaveState / TypedAST emission.

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
	if tp.cachedSize > 0 {
		return tp.cachedSize - 1
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
	tp.cachedSize = tpSize + 1
	return tpSize
}

func fastObjectTypeProtoSerializedSize(m *ObjectTypeProto) int {
	if m.cachedSize > 0 {
		return m.cachedSize - 1
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
	m.cachedSize = size + 1
	return size
}

func fastUnionTypeProtoSerializedSize(m *UnionTypeProto) int {
	if m.cachedSize > 0 {
		return m.cachedSize - 1
	}
	size := 0
	if len(m.unionMember_) > 0 {
		p := protoPackedInt32PayloadSize(m.unionMember_)
		size += 1 + protoVarintLen(uint64(p)) + p
	}
	if len(m.unknownFields) > 0 {
		size += len(protoCanonicalUnknownFields(m.unknownFields))
	}
	m.cachedSize = size + 1
	return size
}

func fastTypeProtoSerializedSize(m *TypeProto) int {
	if m.cachedSize > 0 {
		return m.cachedSize - 1
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
	m.cachedSize = size + 1
	return size
}

func fastSubtypingEdgeSerializedSize(m *SubtypingEdge) int {
	if m.cachedSize > 0 {
		return m.cachedSize - 1
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
	m.cachedSize = size + 1
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
		if err := flush(); err != nil {
			return err
		}
	}
	for _, x := range m.externAst_ {
		buf = append(buf, 26)
		buf = protoAppendVarint(buf, uint64(x.serializedSize()))
		buf = x.appendTo(buf)
		if len(buf) >= 64*1024 {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	for _, x := range m.codeAst_ {
		buf = append(buf, 34)
		buf = protoAppendVarint(buf, uint64(x.serializedSize()))
		buf = x.appendTo(buf)
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
	return flush()
}
