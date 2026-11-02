package colors

import (
	"fmt"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestColorCreateUnionZeroAllocSmall(t *testing.T) {
	p1 := Color_SingleBuilder().SetId(ColorId_FromUnsigned(101)).Build()
	p2 := Color_SingleBuilder().SetId(ColorId_FromUnsigned(102)).Build()
	i1 := Color_SingleBuilder().SetId(ColorId_FromUnsigned(201)).Build()
	i2 := Color_SingleBuilder().SetId(ColorId_FromUnsigned(202)).Build()

	c1 := Color_SingleBuilder().
		SetId(ColorId_FromUnsigned(1)).
		SetPrototypes([]*Color{p1}).
		SetInstanceColors([]*Color{i1}).
		SetOwnProperties([]string{"a", "b"}).
		SetConstructor(true).
		SetClosureAssert(true).
		Build()
	c2 := Color_SingleBuilder().
		SetId(ColorId_FromUnsigned(2)).
		SetPrototypes([]*Color{p1, p2}).
		SetInstanceColors([]*Color{i2}).
		SetOwnProperties([]string{"b", "c"}).
		SetInvalidating(true).
		Build()
	c3 := Color_SingleBuilder().
		SetId(ColorId_FromUnsigned(3)).
		Build()

	elems := []*Color{c1, c2, c3, c1}

	var slab jutil.Slab[Color]
	var colorArena [64]*Color
	var strArena [64]string
	colorOff := 0
	strOff := 0
	allocColors := func(n int) []*Color {
		s := colorArena[colorOff : colorOff+n : colorOff+n]
		colorOff += n
		return s
	}
	allocStrings := func(n int) []string {
		s := strArena[strOff : strOff+n : strOff+n]
		strOff += n
		return s
	}

	var builder Color_Builder
	builder.SetSlab(&slab)
	// Warm up slab chunk before measuring allocations.
	_ = Color_CreateUnionSliceWithBuilder(elems, &builder, allocColors, allocStrings)

	allocs := testing.AllocsPerRun(100, func() {
		colorOff = 0
		strOff = 0
		u := Color_CreateUnionSliceWithBuilder(elems, &builder, allocColors, allocStrings)
		if len(u.GetUnionElements()) != 3 {
			t.Fatalf("union elements = %d, want 3", len(u.GetUnionElements()))
		}
	})
	if allocs != 0 {
		t.Fatalf("Color_CreateUnionSliceWithBuilder small union allocs = %v, want 0", allocs)
	}
}

func TestColorCreateUnionSingleSliceReuseAndSpill(t *testing.T) {
	proto := Color_SingleBuilder().SetId(ColorId_FromUnsigned(500)).Build()
	inst := Color_SingleBuilder().SetId(ColorId_FromUnsigned(501)).Build()
	cWithMeta := Color_SingleBuilder().
		SetId(ColorId_FromUnsigned(10)).
		SetPrototypes([]*Color{proto}).
		SetInstanceColors([]*Color{inst}).
		SetOwnProperties([]string{"only"}).
		Build()
	cPlain := Color_SingleBuilder().
		SetId(ColorId_FromUnsigned(11)).
		Build()

	uSingle := Color_CreateUnion2(cWithMeta, cPlain)
	if &uSingle.GetPrototypes()[0] != &cWithMeta.GetPrototypes()[0] {
		t.Fatal("expected single prototypes slice to be reused without copying")
	}
	if &uSingle.GetInstanceColors()[0] != &cWithMeta.GetInstanceColors()[0] {
		t.Fatal("expected single instanceColors slice to be reused without copying")
	}
	if &uSingle.GetOwnProperties()[0] != &cWithMeta.GetOwnProperties()[0] {
		t.Fatal("expected single ownProperties slice to be reused without copying")
	}

	// Spill past inline [8] and linear-scan [16] thresholds.
	const count = 24
	members := make([]*Color, count)
	var wantProps []string
	for i := 0; i < count; i++ {
		p := Color_SingleBuilder().SetId(ColorId_FromUnsigned(int64(1000 + i))).Build()
		ic := Color_SingleBuilder().SetId(ColorId_FromUnsigned(int64(2000 + i))).Build()
		prop := fmt.Sprintf("prop_%d", i)
		wantProps = append(wantProps, prop)
		members[i] = Color_SingleBuilder().
			SetId(ColorId_FromUnsigned(int64(3000 + i))).
			SetPrototypes([]*Color{p, proto}).
			SetInstanceColors([]*Color{ic, inst}).
			SetOwnProperties([]string{prop, "shared"}).
			Build()
	}
	wantProps = slices.Insert(wantProps, 1, "shared")

	nested := Color_CreateUnionSlice(members[:5])
	combined := append([]*Color{nested}, members[3:]...)
	uLarge := Color_CreateUnionSlice(combined)
	if got := len(uLarge.GetUnionElements()); got != count {
		t.Fatalf("len(unionElements) = %d, want %d", got, count)
	}
	if got := len(uLarge.GetPrototypes()); got != count+1 {
		t.Fatalf("len(prototypes) = %d, want %d", got, count+1)
	}
	if got := len(uLarge.GetInstanceColors()); got != count+1 {
		t.Fatalf("len(instanceColors) = %d, want %d", got, count+1)
	}
	if !slices.Equal(uLarge.GetOwnProperties(), wantProps) {
		t.Fatalf("ownProperties = %v, want %v", uLarge.GetOwnProperties(), wantProps)
	}

	// Permuted >16-element union with freshly constructed value-equal *Color members.
	reversed := make([]*Color, count)
	for i := range count {
		p := Color_SingleBuilder().SetId(ColorId_FromUnsigned(int64(1000 + (count - 1 - i)))).Build()
		ic := Color_SingleBuilder().SetId(ColorId_FromUnsigned(int64(2000 + (count - 1 - i)))).Build()
		prop := fmt.Sprintf("prop_%d", count-1-i)
		reversed[i] = Color_SingleBuilder().
			SetId(ColorId_FromUnsigned(int64(3000 + (count - 1 - i)))).
			SetPrototypes([]*Color{proto, p}).
			SetInstanceColors([]*Color{inst, ic}).
			SetOwnProperties([]string{"shared", prop}).
			Build()
	}
	uReversed := Color_CreateUnionSlice(reversed)
	if !uLarge.Equals(uReversed) {
		t.Fatal("expected permuted large union with distinct value-equal *Color pointers to compare equal")
	}
}
