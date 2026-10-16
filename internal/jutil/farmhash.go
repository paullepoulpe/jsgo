package jutil

import (
	"encoding/binary"
	"math/bits"
)

// Go-only: language-adapter — Guava Hashing.farmHashFingerprint64() (the
// FarmHash Fingerprint64 / farmhashna algorithm, which Guava guarantees to
// be stable across releases), plus the buffering Hasher that feeds it.

const (
	farmK0 uint64 = 0xc3a5c85c97cb3127
	farmK1 uint64 = 0xb492b66fbe98f273
	farmK2 uint64 = 0x9ae16a3b2f90404f
)

// FarmHashFingerprint64Hasher is the Hasher returned by
// Hashing.farmHashFingerprint64().newHasher(): it buffers the input bytes
// and hashes them on Hash().
type FarmHashFingerprint64Hasher struct {
	buf []byte
}

// Hashing_FarmHashFingerprint64NewHasher is
// Hashing.farmHashFingerprint64().newHasher().
func Hashing_FarmHashFingerprint64NewHasher() *FarmHashFingerprint64Hasher {
	return &FarmHashFingerprint64Hasher{}
}

// PutLong is Hasher.putLong(long): eight little-endian bytes.
func (h *FarmHashFingerprint64Hasher) PutLong(v int64) *FarmHashFingerprint64Hasher {
	h.buf = binary.LittleEndian.AppendUint64(h.buf, uint64(v))
	return h
}

// PutBytes is Hasher.putBytes(byte[]).
func (h *FarmHashFingerprint64Hasher) PutBytes(b []byte) *FarmHashFingerprint64Hasher {
	h.buf = append(h.buf, b...)
	return h
}

// HashAsLong is hasher.hash().asLong().
func (h *FarmHashFingerprint64Hasher) HashAsLong() int64 {
	return int64(FarmHashFingerprint64(h.buf))
}

// FarmHashFingerprint64 is FarmHashFingerprint64.fingerprint(bytes, 0, len).
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
