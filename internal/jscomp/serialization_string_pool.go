// Port of com/google/javascript/jscomp/serialization/StringPool.java @ closure-compiler v20260415

package jscomp

import (
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
	decoder := Wtf8_Decoder_(int(proto.GetMaxLength()))

	pool := []string{""}
	stringsList := proto.GetStringsList()
	for i := 0; i < len(stringsList); i++ {
		s := stringsList[i]
		pool = append(pool, decoder.Decode(s))
	}

	return newStringPool(int(proto.GetMaxLength()), pool)
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
	builder := StringPoolProto_NewBuilder().SetMaxLength(int32(p.maxLength))
	first := true
	for s := range p.pool.Stream() {
		if first {
			first = false
			continue
		}
		builder.AddStrings(Wtf8_EncodeToWtf8(s))
	}
	return builder.Build()
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
	pool      *jutil.LinkedHashMap[string, int]
}

func newStringPool_Builder() *StringPool_Builder {
	b := &StringPool_Builder{pool: jutil.NewLinkedHashMap[string, int]()}
	b.Put("")
	return b
}

// Put inserts the given string into the string pool if not present and
// returns its index.
//
// Java: StringPool.Builder.put(String)
func (b *StringPool_Builder) Put(s string) int {
	if n := jutil.StringLength(s); n > b.maxLength {
		b.maxLength = n
	}
	return b.pool.ComputeIfAbsent(s, func(string) int { return b.pool.Size() })
}

// Java: StringPool.Builder.putAnd(String)
func (b *StringPool_Builder) PutAnd(s string) *StringPool_Builder {
	b.Put(s)
	return b
}

// Java: StringPool.Builder.build()
func (b *StringPool_Builder) Build() *StringPool {
	return newStringPool(b.maxLength, b.pool.KeySlice())
}

var stringPool_EMPTY = StringPool_Builder_().Build()
