// Port of com/google/javascript/jscomp/serialization/StringPool.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// StringPool aggregates strings into a StringPool.
//
// The zeroth offset in the string pool is always the empty string. This is
// validated inside TypedAstDeserializer.
//
// This implies default/unset/0-valued uuint32 StringPool pointers in protos
// are equivalent to the empty string.
type StringPool struct {
	maxLength int // max length of any individual string in the pool
	pool      *rhino.RhinoStringPool_LazyInternedStringList
}

// Java: StringPool.fromProto(StringPoolProto)
func StringPool_FromProto(proto *StringPoolProto) *StringPool {
	maxLength := int(proto.GetMaxLength())
	decoder := Wtf8_Decoder_(maxLength)

	stringsList := proto.GetStringsList()
	ids := make([]uint32, len(stringsList)+1)
	ids[0] = rhino.RhinoStringPool_InternID("")
	for i := 0; i < len(stringsList); i++ {
		ids[i+1] = decoder.DecodeInternID(stringsList[i])
	}

	return &StringPool{
		maxLength: maxLength,
		pool:      rhino.NewRhinoStringPool_LazyInternedStringListFromIDs(ids),
	}
}

// Java: StringPool.empty()
func StringPool_Empty() *StringPool {
	return stringPool_EMPTY
}

func newStringPool(maxLength int, pool []string) *StringPool {
	p := &StringPool{
		maxLength: maxLength,
		pool:      rhino.NewRhinoStringPool_LazyInternedStringList(pool),
	}
	jutil.CheckState(pool[0] == "")
	return p
}

// Java: StringPool.get(int)
func (p *StringPool) Get(offset int) string {
	return p.pool.Get(offset)
}

// Java: StringPool.getInternedStrings()
func (p *StringPool) GetInternedStrings() *rhino.RhinoStringPool_LazyInternedStringList {
	return p.pool
}

// Java: StringPool.toProto()
func (p *StringPool) ToProto() *StringPoolProto {
	return p.ToProtoWithThreads(1)
}

func (p *StringPool) ToProtoWithThreads(numThreads int) *StringPoolProto {
	n := p.pool.Len()
	var strings [][]byte
	spSize := 0
	if n > 1 {
		count := n - 1
		strings = make([][]byte, count)
		workers := 1
		if numThreads > 1 && count >= 2048 {
			workers = min(numThreads, (count+1023)/1024)
		}
		if workers <= 1 {
			spSize = encodeStringPoolShard(p.pool, strings, 1, n)
		} else {
			shardSizes := make([]int, workers)
			var wg sync.WaitGroup
			wg.Add(workers)
			for w := 0; w < workers; w++ {
				start := 1 + (w*count)/workers
				end := 1 + ((w+1)*count)/workers
				workerIdx := w
				go func() {
					defer wg.Done()
					shardSizes[workerIdx] = encodeStringPoolShard(p.pool, strings, start, end)
				}()
			}
			wg.Wait()
			for _, sz := range shardSizes {
				spSize += sz
			}
		}
	}
	if p.maxLength != 0 {
		spSize += 1 + protoVarintLen(uint64(int64(p.maxLength)))
	}
	return &StringPoolProto{
		strings_:     strings,
		maxLength_:   int32(p.maxLength),
		frozen:       true,
		memoizedSize: spSize,
	}
}

func encodeStringPoolShard(pool *rhino.RhinoStringPool_LazyInternedStringList, out [][]byte, start, end int) int {
	const slabSize = 256 * 1024
	var slab []byte
	spSize := 0
	for i := start; i < end; i++ {
		s := pool.Get(i)
		var b []byte
		if jutil.IsASCII(s) {
			slen := len(s)
			if slen == 0 {
				b = nil
			} else if slen <= 4096 {
				if len(slab) < slen {
					slab = make([]byte, slabSize)
				}
				b = slab[:slen:slen]
				copy(b, s)
				slab = slab[slen:]
			} else {
				b = []byte(s)
			}
		} else {
			b = Wtf8_EncodeToWtf8(s)
		}
		out[i-1] = b
		spSize += 1 + protoVarintLen(uint64(len(b))) + len(b)
	}
	return spSize
}

// StringPool_Builder_ is StringPool.builder().
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name StringPool_Builder for the nested Builder type.
func StringPool_Builder_() *StringPool_Builder {
	return newStringPool_Builder()
}

// StringPool_Builder is StringPool.Builder.
type StringPool_Builder struct {
	maxLength int
	indexByID map[uint32]int
	ids       []uint32
}

func newStringPool_Builder() *StringPool_Builder {
	emptyID := rhino.RhinoStringPool_InternID("")
	return &StringPool_Builder{
		indexByID: map[uint32]int{emptyID: 0},
		ids:       []uint32{emptyID},
	}
}

// PutID inserts the string with the given 32-bit RhinoStringPool ID into the
// string pool if not present and returns its index.
func (b *StringPool_Builder) PutID(id uint32) int {
	if id == 0 {
		return 0
	}
	if idx, ok := b.indexByID[id]; ok {
		return idx
	}
	s := rhino.RhinoStringPool_StringFromID(id)
	if len(s) > b.maxLength {
		if n := jutil.StringLength(s); n > b.maxLength {
			b.maxLength = n
		}
	}
	idx := len(b.ids)
	b.ids = append(b.ids, id)
	b.indexByID[id] = idx
	return idx
}

// Lookup returns the pool index of s if present.
func (b *StringPool_Builder) Lookup(s string) (int, bool) {
	if s == "" {
		return 0, true
	}
	id := rhino.RhinoStringPool_LookupID(s)
	if id == 0 {
		return 0, false
	}
	idx, ok := b.indexByID[id]
	return idx, ok
}

// Put inserts the given string into the string pool if not present and
// returns its index.
//
// Java: StringPool.Builder.put(String)
func (b *StringPool_Builder) Put(s string) int {
	return b.PutID(rhino.RhinoStringPool_InternID(s))
}

// Java: StringPool.Builder.putAnd(String)
func (b *StringPool_Builder) PutAnd(s string) *StringPool_Builder {
	b.Put(s)
	return b
}

// Java: StringPool.Builder.build()
func (b *StringPool_Builder) Build() *StringPool {
	return &StringPool{
		maxLength: b.maxLength,
		pool:      rhino.NewRhinoStringPool_LazyInternedStringListFromIDs(b.ids),
	}
}

var stringPool_EMPTY = StringPool_Builder_().Build()
