/*
 * Copyright (C) 2015 The Guava Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package guava

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/bits"
	"os"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
)

// ---------------------------------------------------------------------------
// com.google.common.io.BaseEncoding
// ---------------------------------------------------------------------------

type BaseEncoding struct {
	urlSafe bool
	raw     bool
}

var (
	base64Std    = &BaseEncoding{urlSafe: false, raw: false}
	base64StdRaw = &BaseEncoding{urlSafe: false, raw: true}
	base64URL    = &BaseEncoding{urlSafe: true, raw: false}
	base64URLRaw = &BaseEncoding{urlSafe: true, raw: true}
)

func BaseEncoding_Base64() *BaseEncoding    { return base64Std }
func BaseEncoding_Base64Url() *BaseEncoding { return base64URL }

func (b *BaseEncoding) OmitPadding() *BaseEncoding {
	if b.urlSafe {
		return base64URLRaw
	}
	return base64StdRaw
}

func (b *BaseEncoding) enc() *base64.Encoding {
	switch {
	case b.urlSafe && b.raw:
		return base64.RawURLEncoding
	case b.urlSafe:
		return base64.URLEncoding
	case b.raw:
		return base64.RawStdEncoding
	default:
		return base64.StdEncoding
	}
}

func (b *BaseEncoding) Encode(bytes []byte) string {
	return b.enc().EncodeToString(bytes)
}

// Decode decodes encoded base64 (tolerating optional/missing padding on standard/url base64 like Guava).
func (b *BaseEncoding) Decode(encoded string) []byte {
	trimmed := strings.TrimRight(encoded, "=")
	var rawEnc *base64.Encoding
	if b.urlSafe {
		rawEnc = base64.RawURLEncoding
	} else {
		rawEnc = base64.RawStdEncoding
	}
	data, err := rawEnc.DecodeString(trimmed)
	if err != nil {
		panic(javalang.NewIllegalArgumentException("Invalid base64 input: " + err.Error()))
	}
	return data
}

func (b *BaseEncoding) DecodeChecked(encoded string) ([]byte, error) {
	trimmed := strings.TrimRight(encoded, "=")
	var rawEnc *base64.Encoding
	if b.urlSafe {
		rawEnc = base64.RawURLEncoding
	} else {
		rawEnc = base64.RawStdEncoding
	}
	data, err := rawEnc.DecodeString(trimmed)
	if err != nil {
		return nil, javaio.NewIOException("Invalid base64 input: " + err.Error())
	}
	return data, nil
}

// ---------------------------------------------------------------------------
// com.google.common.io.ByteStreams, CharStreams, CharSource, Files
// ---------------------------------------------------------------------------

func ByteStreams_ToByteArray(in io.Reader) ([]byte, error) {
	CheckNotNull(in)
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	return data, nil
}

func ByteStreams_Copy(from io.Reader, to io.Writer) (int64, error) {
	CheckNotNull(from)
	CheckNotNull(to)
	n, err := io.Copy(to, from)
	if err != nil {
		return n, javaio.NewIOException(err.Error())
	}
	return n, nil
}

func CharStreams_ToString(r javaio.Reader) string {
	CheckNotNull(r)
	if ds, ok := r.(interface{ DrainString() (string, bool) }); ok {
		if s, ok := ds.DrainString(); ok {
			return s
		}
	}
	if sa, ok := r.(interface{ ReadAllString() string }); ok {
		return sa.ReadAllString()
	}
	var buf [4096]uint16
	var units []uint16
	for {
		n := r.Read(buf[:], 0, len(buf))
		if n < 0 {
			break
		}
		units = append(units, buf[:n]...)
	}
	return javalang.StringFromUTF16(units)
}

func CharStreams_ReadLines(r javaio.Reader) []string {
	CheckNotNull(r)
	br, ok := r.(*javaio.BufferedReader)
	if !ok {
		br = javaio.NewBufferedReader(r)
	}
	var lines []string
	for {
		line, err := br.ReadLine()
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	if lines == nil {
		return []string{}
	}
	return lines
}

type CharSource struct {
	readFn func() (string, error)
}

func CharSource_Wrap(content string) *CharSource {
	return &CharSource{readFn: func() (string, error) { return content, nil }}
}

func (cs *CharSource) Read() (string, error) {
	return cs.readFn()
}

func (cs *CharSource) ReadLines() ([]string, error) {
	s, err := cs.Read()
	if err != nil {
		return nil, err
	}
	return CharStreams_ReadLines(javaio.NewStringReader(s)), nil
}

func Files_AsCharSource(file *javaio.File, charset *javanio.Charset) *CharSource {
	CheckNotNull(file)
	CheckNotNull(charset)
	return &CharSource{
		readFn: func() (string, error) {
			raw, err := os.ReadFile(file.GetPath())
			if err != nil {
				return "", javaio.NewIOException(err.Error())
			}
			return charset.Decode(raw), nil
		},
	}
}

func Files_ToByteArray(file *javaio.File) ([]byte, error) {
	CheckNotNull(file)
	raw, err := os.ReadFile(file.GetPath())
	if err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	return raw, nil
}

func Files_WriteBytes(bytes []byte, to *javaio.File) error {
	CheckNotNull(bytes)
	CheckNotNull(to)
	if err := os.WriteFile(to.GetPath(), bytes, 0o666); err != nil {
		return javaio.NewIOException(err.Error())
	}
	return nil
}

func Files_AsCharSinkWrite(to *javaio.File, charset *javanio.Charset, content string) error {
	CheckNotNull(to)
	CheckNotNull(charset)
	return Files_WriteBytes(charset.Encode(content), to)
}

// ---------------------------------------------------------------------------
// com.google.common.hash.{Hashing, HashFunction, Hasher, HashCode}
// ---------------------------------------------------------------------------

// HashCode models com.google.common.hash.HashCode (64-bit or 128-bit little-endian).
type HashCode struct {
	bits128 bool
	lo      uint64
	hi      uint64
}

func HashCode_FromLong(val int64) HashCode {
	return HashCode{lo: uint64(val)}
}

func (hc HashCode) AsInt() int32 {
	return int32(uint32(hc.lo))
}

func (hc HashCode) AsLong() int64 {
	return int64(hc.lo)
}

func (hc HashCode) AsBytes() []byte {
	if hc.bits128 {
		var b [16]byte
		binary.LittleEndian.PutUint64(b[0:8], hc.lo)
		binary.LittleEndian.PutUint64(b[8:16], hc.hi)
		return b[:]
	}
	var b [8]byte
	binary.LittleEndian.PutUint64(b[0:8], hc.lo)
	return b[:]
}

func (hc HashCode) String() string {
	return hex.EncodeToString(hc.AsBytes())
}

// Hasher models com.google.common.hash.Hasher.
type Hasher interface {
	PutByte(b byte) Hasher
	PutBytes(b []byte) Hasher
	PutInt(v int32) Hasher
	PutLong(v int64) Hasher
	PutChar(c uint16) Hasher
	PutUnencodedChars(s string) Hasher
	PutString(s string, charset *javanio.Charset) Hasher
	Hash() HashCode
}

// HashFunction models com.google.common.hash.HashFunction.
type HashFunction interface {
	NewHasher() Hasher
	HashBytes(b []byte) HashCode
	HashInt(v int32) HashCode
	HashLong(v int64) HashCode
	HashUnencodedChars(s string) HashCode
	HashString(s string, charset *javanio.Charset) HashCode
}

type bufferedHasher struct {
	buf    []byte
	inline [64]byte
	finish func([]byte) HashCode
}

func newBufferedHasher(finish func([]byte) HashCode) *bufferedHasher {
	h := &bufferedHasher{finish: finish}
	h.buf = h.inline[:0]
	return h
}

func (h *bufferedHasher) PutByte(b byte) Hasher {
	h.buf = append(h.buf, b)
	return h
}

func (h *bufferedHasher) PutBytes(b []byte) Hasher {
	h.buf = append(h.buf, b...)
	return h
}

func (h *bufferedHasher) PutInt(v int32) Hasher {
	h.buf = binary.LittleEndian.AppendUint32(h.buf, uint32(v))
	return h
}

func (h *bufferedHasher) PutLong(v int64) Hasher {
	h.buf = binary.LittleEndian.AppendUint64(h.buf, uint64(v))
	return h
}

func (h *bufferedHasher) PutChar(c uint16) Hasher {
	h.buf = binary.LittleEndian.AppendUint16(h.buf, c)
	return h
}

func (h *bufferedHasher) PutUnencodedChars(s string) Hasher {
	units := javalang.StringToUTF16(s)
	for _, u := range units {
		h.buf = binary.LittleEndian.AppendUint16(h.buf, u)
	}
	return h
}

func (h *bufferedHasher) PutString(s string, charset *javanio.Charset) Hasher {
	CheckNotNull(charset)
	h.buf = append(h.buf, charset.Encode(s)...)
	return h
}

func (h *bufferedHasher) Hash() HashCode {
	return h.finish(h.buf)
}

type simpleHashFunction struct {
	hashBytes func([]byte) HashCode
}

func (f simpleHashFunction) NewHasher() Hasher {
	return newBufferedHasher(f.hashBytes)
}

func (f simpleHashFunction) HashBytes(b []byte) HashCode {
	return f.hashBytes(b)
}

func (f simpleHashFunction) HashInt(v int32) HashCode {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	return f.hashBytes(b[:])
}

func (f simpleHashFunction) HashLong(v int64) HashCode {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	return f.hashBytes(b[:])
}

func (f simpleHashFunction) HashUnencodedChars(s string) HashCode {
	return f.NewHasher().PutUnencodedChars(s).Hash()
}

func (f simpleHashFunction) HashString(s string, charset *javanio.Charset) HashCode {
	return f.HashBytes(charset.Encode(s))
}

var (
	farmHash64Func = simpleHashFunction{
		hashBytes: func(b []byte) HashCode {
			return HashCode{lo: FarmHashFingerprint64(b)}
		},
	}
	murmur3_128ZeroFunc = Hashing_Murmur3_128WithSeed(0)
)

// Hashing_FarmHashFingerprint64 returns Guava's Hashing.farmHashFingerprint64() HashFunction.
func Hashing_FarmHashFingerprint64() HashFunction {
	return farmHash64Func
}

// Hashing_Murmur3_128 returns Guava's Hashing.murmur3_128() HashFunction (seed 0).
func Hashing_Murmur3_128() HashFunction {
	return murmur3_128ZeroFunc
}

// Hashing_Murmur3_128WithSeed returns Guava's Hashing.murmur3_128(seed) HashFunction.
func Hashing_Murmur3_128WithSeed(seed int32) HashFunction {
	uSeed := uint32(seed)
	return simpleHashFunction{
		hashBytes: func(b []byte) HashCode {
			h1, h2 := Murmur3_128(b, uSeed)
			return HashCode{bits128: true, lo: h1, hi: h2}
		},
	}
}

// ---------------------------------------------------------------------------
// FarmHash Fingerprint64 (farmhashna) & MurmurHash3 x64_128 (Guava compatible)
// ---------------------------------------------------------------------------

const (
	farmK0 uint64 = 0xc3a5c85c97cb3127
	farmK1 uint64 = 0xb492b66fbe98f273
	farmK2 uint64 = 0x9ae16a3b2f90404f
)

// FarmHashFingerprint64 computes Guava's FarmHashFingerprint64.fingerprint(b).
func FarmHashFingerprint64(b []byte) uint64 {
	n := len(b)
	switch {
	case n <= 16:
		return farmHashLength0to16(b)
	case n <= 32:
		return farmHashLength17to32(b)
	case n <= 64:
		return farmHashLength33To64(b)
	}
	return farmHashLength65Plus(b)
}

func farmLoad64(b []byte, i int) uint64 { return binary.LittleEndian.Uint64(b[i:]) }
func farmLoad32(b []byte, i int) uint64 { return uint64(binary.LittleEndian.Uint32(b[i:])) }
func farmShiftMix(v uint64) uint64      { return v ^ (v >> 47) }
func farmRotr(v uint64, s int) uint64   { return bits.RotateLeft64(v, -s) }

func farmHashLength16(u, v, mul uint64) uint64 {
	a := (u ^ v) * mul
	a ^= a >> 47
	b := (v ^ a) * mul
	b ^= b >> 47
	b *= mul
	return b
}

func farmWeakHashLength32WithSeeds(b []byte, off int, seedA, seedB uint64) (uint64, uint64) {
	part1 := farmLoad64(b, off)
	part2 := farmLoad64(b, off+8)
	part3 := farmLoad64(b, off+16)
	part4 := farmLoad64(b, off+24)
	seedA += part1
	seedB = farmRotr(seedB+seedA+part4, 21)
	c := seedA
	seedA += part2
	seedA += part3
	seedB += farmRotr(seedA, 44)
	return seedA + part4, seedB + c
}

func farmHashLength0to16(b []byte) uint64 {
	n := len(b)
	if n >= 8 {
		mul := farmK2 + uint64(n)*2
		a := farmLoad64(b, 0) + farmK2
		bb := farmLoad64(b, n-8)
		c := farmRotr(bb, 37)*mul + a
		d := (farmRotr(a, 25) + bb) * mul
		return farmHashLength16(c, d, mul)
	}
	if n >= 4 {
		mul := farmK2 + uint64(n)*2
		a := farmLoad32(b, 0)
		return farmHashLength16(uint64(n)+(a<<3), farmLoad32(b, n-4), mul)
	}
	if n > 0 {
		a := b[0]
		bb := b[n>>1]
		c := b[n-1]
		y := uint32(a) + uint32(bb)<<8
		z := uint32(n) + uint32(c)<<2
		return farmShiftMix(uint64(int64(int32(y)))*farmK2^uint64(int64(int32(z)))*farmK0) * farmK2
	}
	return farmK2
}

func farmHashLength17to32(b []byte) uint64 {
	n := len(b)
	mul := farmK2 + uint64(n)*2
	a := farmLoad64(b, 0) * farmK1
	bb := farmLoad64(b, 8)
	c := farmLoad64(b, n-8) * mul
	d := farmLoad64(b, n-16) * farmK2
	return farmHashLength16(farmRotr(a+bb, 43)+farmRotr(c, 30)+d, a+farmRotr(bb+farmK2, 18)+c, mul)
}

func farmHashLength33To64(b []byte) uint64 {
	n := len(b)
	mul := farmK2 + uint64(n)*2
	a := farmLoad64(b, 0) * farmK2
	bb := farmLoad64(b, 8)
	c := farmLoad64(b, n-8) * mul
	d := farmLoad64(b, n-16) * farmK2
	y := farmRotr(a+bb, 43) + farmRotr(c, 30) + d
	z := farmHashLength16(y, a+farmRotr(bb+farmK2, 18)+c, mul)
	e := farmLoad64(b, 16) * mul
	f := farmLoad64(b, 24)
	g := (y + farmLoad64(b, n-32)) * mul
	h := (z + farmLoad64(b, n-24)) * mul
	return farmHashLength16(farmRotr(e+f, 43)+farmRotr(g, 30)+h, e+farmRotr(f+a, 18)+g, mul)
}

func farmHashLength65Plus(b []byte) uint64 {
	n := len(b)
	var seed uint64 = 81
	x := seed
	y := seed*farmK1 + 113
	z := farmShiftMix(y*farmK2+113) * farmK2
	var v0, v1, w0, w1 uint64
	x = x*farmK2 + farmLoad64(b, 0)
	off := 0
	end := ((n - 1) / 64) * 64
	last64 := end + ((n - 1) & 63) - 63
	for {
		x = farmRotr(x+y+v0+farmLoad64(b, off+8), 37) * farmK1
		y = farmRotr(y+v1+farmLoad64(b, off+48), 42) * farmK1
		x ^= w1
		y += v0 + farmLoad64(b, off+40)
		z = farmRotr(z+w0, 33) * farmK1
		v0, v1 = farmWeakHashLength32WithSeeds(b, off, v1*farmK1, x+w0)
		w0, w1 = farmWeakHashLength32WithSeeds(b, off+32, z+w1, y+farmLoad64(b, off+16))
		x, z = z, x
		off += 64
		if off == end {
			break
		}
	}
	mul := farmK1 + ((z & 0xff) << 1)
	off = last64
	w0 += uint64((n - 1) & 63)
	v0 += w0
	w0 += v0
	x = farmRotr(x+y+v0+farmLoad64(b, off+8), 37) * mul
	y = farmRotr(y+v1+farmLoad64(b, off+48), 42) * mul
	x ^= w1 * 9
	y += v0*9 + farmLoad64(b, off+40)
	z = farmRotr(z+w0, 33) * mul
	v0, v1 = farmWeakHashLength32WithSeeds(b, off, v1*mul, x+w0)
	w0, w1 = farmWeakHashLength32WithSeeds(b, off+32, z+w1, y+farmLoad64(b, off+16))
	return farmHashLength16(farmHashLength16(v0, w0, mul)+farmShiftMix(y)*farmK0+x,
		farmHashLength16(v1, w1, mul)+z, mul)
}

const (
	murmur3C1 uint64 = 0x87c37b91114253d5
	murmur3C2 uint64 = 0x4cf5ad432745937f
)

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

// Murmur3_128 computes Guava's Murmur3_128HashFunction (x64 128-bit) returning (h1, h2).
func Murmur3_128(data []byte, seed uint32) (uint64, uint64) {
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
