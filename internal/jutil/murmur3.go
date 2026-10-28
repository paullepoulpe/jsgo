// Derived from com/google/common/hash/Murmur3_128HashFunction.java
// (Guava 33.4.6-jre): murmur3x64_128 and its helpers follow Guava's
// Murmur3_128Hasher (mixK1, mixK2, fmix64, and processRemaining's mixing of
// both halves of a partial block). MurmurHash3 itself was written by Austin
// Appleby, who placed it in the public domain. The Hasher and HashCode API
// around it is emulated from Guava's documentation.
//
// Copyright (C) 2011 The Guava Authors
// Licensed under the Apache License, Version 2.0; see LICENSES/Apache-2.0.txt.

package jutil

import (
	"encoding/binary"
	"math/bits"
	"strings"
)

// Go-only: language-adapter — Guava's Hashing.murmur3_128() (MurmurHash3
// x64 128-bit, seed 0) and its streaming Hasher. A Hasher is a sequence of
// bytes; Guava's streaming implementation and hashing the concatenated
// bytes at once give the same result.

const (
	murmur3C1 uint64 = 0x87c37b91114253d5
	murmur3C2 uint64 = 0x4cf5ad432745937f
)

// Murmur3_128Hasher is the Hasher returned by
// Hashing.murmur3_128().newHasher().
type Murmur3_128Hasher struct {
	seed uint32
	buf  []byte
}

// Hashing_Murmur3_128NewHasher is Hashing.murmur3_128().newHasher().
func Hashing_Murmur3_128NewHasher() *Murmur3_128Hasher { return &Murmur3_128Hasher{} }

// Hashing_Murmur3_128WithSeedNewHasher is
// Hashing.murmur3_128(seed).newHasher().
func Hashing_Murmur3_128WithSeedNewHasher(seed int32) *Murmur3_128Hasher {
	return &Murmur3_128Hasher{seed: uint32(seed)}
}

// PutByte is Hasher.putByte(byte).
func (h *Murmur3_128Hasher) PutByte(b int8) *Murmur3_128Hasher {
	h.buf = append(h.buf, byte(b))
	return h
}

// PutBytes is Hasher.putBytes(byte[]).
func (h *Murmur3_128Hasher) PutBytes(b []byte) *Murmur3_128Hasher {
	h.buf = append(h.buf, b...)
	return h
}

// PutInt is Hasher.putInt(int) (little-endian).
func (h *Murmur3_128Hasher) PutInt(v int32) *Murmur3_128Hasher {
	h.buf = binary.LittleEndian.AppendUint32(h.buf, uint32(v))
	return h
}

// PutLong is Hasher.putLong(long) (little-endian).
func (h *Murmur3_128Hasher) PutLong(v int64) *Murmur3_128Hasher {
	h.buf = binary.LittleEndian.AppendUint64(h.buf, uint64(v))
	return h
}

// PutBoolean is Hasher.putBoolean(boolean).
func (h *Murmur3_128Hasher) PutBoolean(v bool) *Murmur3_128Hasher {
	if v {
		return h.PutByte(1)
	}
	return h.PutByte(0)
}

// PutChar is Hasher.putChar(char) (little-endian).
func (h *Murmur3_128Hasher) PutChar(c uint16) *Murmur3_128Hasher {
	h.buf = binary.LittleEndian.AppendUint16(h.buf, c)
	return h
}

// PutUnencodedChars is Hasher.putUnencodedChars(CharSequence): every UTF-16
// code unit, little-endian.
func (h *Murmur3_128Hasher) PutUnencodedChars(s string) *Murmur3_128Hasher {
	if isASCII(s) {
		for i := 0; i < len(s); i++ {
			h.buf = append(h.buf, s[i], 0)
		}
		return h
	}
	for _, c := range StringToUTF16(s) {
		h.PutChar(c)
	}
	return h
}

// PutStringUTF8 is Hasher.putString(CharSequence, UTF_8).
func (h *Murmur3_128Hasher) PutStringUTF8(s string) *Murmur3_128Hasher {
	h.buf = append(h.buf, StringGetBytesUTF8(s)...)
	return h
}

// Hash is Hasher.hash().
func (h *Murmur3_128Hasher) Hash() *HashCode {
	h1, h2 := murmur3x64_128(h.buf, h.seed)
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b, h1)
	binary.LittleEndian.PutUint64(b[8:], h2)
	return &HashCode{bytes: b}
}

// HashCode is Guava's com.google.common.hash.HashCode.
type HashCode struct {
	bytes []byte
}

// AsInt is HashCode.asInt(): the first four bytes, little-endian.
func (c *HashCode) AsInt() int32 {
	return int32(binary.LittleEndian.Uint32(c.bytes))
}

// AsLong is HashCode.asLong(): the first eight bytes, little-endian.
func (c *HashCode) AsLong() int64 {
	return int64(binary.LittleEndian.Uint64(c.bytes))
}

// AsBytes is HashCode.asBytes().
func (c *HashCode) AsBytes() []byte { return append([]byte(nil), c.bytes...) }

// Bits is HashCode.bits().
func (c *HashCode) Bits() int { return 8 * len(c.bytes) }

// String is HashCode.toString(): lowercase hex of the bytes.
func (c *HashCode) String() string {
	const hexDigits = "0123456789abcdef"
	var sb strings.Builder
	for _, b := range c.bytes {
		sb.WriteByte(hexDigits[b>>4])
		sb.WriteByte(hexDigits[b&0xf])
	}
	return sb.String()
}

func murmur3MixK1(k1 uint64) uint64 {
	k1 *= murmur3C1
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= murmur3C2
	return k1
}

func murmur3MixK2(k2 uint64) uint64 {
	k2 *= murmur3C2
	k2 = bits.RotateLeft64(k2, 33)
	k2 *= murmur3C1
	return k2
}

func murmur3Fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

func murmur3x64_128(data []byte, seed uint32) (uint64, uint64) {
	h1, h2 := uint64(seed), uint64(seed)
	n := len(data)
	i := 0
	for ; i+16 <= n; i += 16 {
		k1 := binary.LittleEndian.Uint64(data[i:])
		k2 := binary.LittleEndian.Uint64(data[i+8:])
		h1 ^= murmur3MixK1(k1)
		h1 = bits.RotateLeft64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729
		h2 ^= murmur3MixK2(k2)
		h2 = bits.RotateLeft64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}
	if rem := n - i; rem > 0 {
		var k1, k2 uint64
		tail := data[i:]
		for j := rem - 1; j >= 8; j-- {
			k2 ^= uint64(tail[j]) << (8 * uint(j-8))
		}
		for j := min(rem, 8) - 1; j >= 0; j-- {
			k1 ^= uint64(tail[j]) << (8 * uint(j))
		}
		h1 ^= murmur3MixK1(k1)
		h2 ^= murmur3MixK2(k2)
	}
	h1 ^= uint64(n)
	h2 ^= uint64(n)
	h1 += h2
	h2 += h1
	h1 = murmur3Fmix64(h1)
	h2 = murmur3Fmix64(h2)
	h1 += h2
	h2 += h1
	return h1, h2
}

// Hashing_Murmur3_128BytesAsInt returns Hashing.murmur3_128().hashBytes(data).asInt()
// without allocating a Hasher or HashCode struct.
func Hashing_Murmur3_128BytesAsInt(data []byte) int32 {
	h1, _ := murmur3x64_128(data, 0)
	return int32(uint32(h1))
}
