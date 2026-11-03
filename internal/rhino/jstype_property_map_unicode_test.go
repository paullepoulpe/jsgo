package rhino

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestPropertyMapKeepsByIDWithUnicodeAndHandlesNonCanonicalSurrogates(t *testing.T) {
	reg := NewJSTypeRegistry(nil)
	numType := reg.GetNativeType(JSTypeNative_NUMBER_TYPE)

	pm := NewPropertyMap()
	const numASCII = 200
	for i := range numASCII {
		name := fmt.Sprintf("globalProp_%04d", i)
		pm.PutProperty(name, NewProperty(name, numType, false, NilNode))
	}
	// Insert canonical non-ASCII BMP, CJK, and astral emoji keys.
	unicodeKeys := []string{"π", "café", "日本語", "😀", "🚀_rocket"}
	for _, k := range unicodeKeys {
		pm.PutProperty(k, NewProperty(k, numType, false, NilNode))
	}

	if pm.ext == nil || pm.ext.byID == nil {
		t.Fatalf("expected PropertyMap.ext.byID to stay populated for canonical Unicode keys")
	}
	if pm.hasNonCanonical {
		t.Fatalf("expected hasNonCanonical=false when all keys are valid UTF-8")
	}

	for i := range numASCII {
		name := fmt.Sprintf("globalProp_%04d", i)
		if pm.GetOwnProperty(name) == nil {
			t.Fatalf("missing ASCII property %q", name)
		}
	}
	for _, k := range unicodeKeys {
		if pm.GetOwnProperty(k) == nil {
			t.Fatalf("missing Unicode property %q", k)
		}
	}
	if pm.GetOwnProperty("never_interned_ascii_miss_xyz") != nil {
		t.Fatalf("unexpected hit on missing ASCII key")
	}
	if pm.GetOwnProperty("never_interned_unicode_miss_Ω") != nil {
		t.Fatalf("unexpected hit on missing Unicode key")
	}

	// Non-canonical WTF-8 surrogate pair ("😀" constructed as two concatenated lone surrogates)
	// must still match and replace canonical "😀" via UTF-16 code-unit equality.
	nonCanonSmile := jutil.StringFromChar(0xD83D) + jutil.StringFromChar(0xDE00)
	if got := pm.GetOwnProperty(nonCanonSmile); got == nil {
		t.Fatalf("expected non-canonical surrogate pair lookup to find canonical emoji property")
	}
	replacement := NewProperty(nonCanonSmile, numType, true, NilNode)
	pm.PutProperty(nonCanonSmile, replacement)
	if !pm.hasNonCanonical {
		t.Fatalf("expected hasNonCanonical=true after non-canonical WTF-8 key")
	}
	if got := pm.GetOwnProperty("😀"); got != replacement {
		t.Fatalf("expected canonical lookup after non-canonical replacement to return replacement")
	}
	if got := pm.GetOwnProperty(nonCanonSmile); got != replacement {
		t.Fatalf("expected non-canonical lookup to return replacement")
	}
}

func BenchmarkPropertyMapUnicodeAndGlobalMix(b *testing.B) {
	reg := NewJSTypeRegistry(nil)
	numType := reg.GetNativeType(JSTypeNative_NUMBER_TYPE)

	b.Run("GlobalMix10kPlusUnicode", func(b *testing.B) {
		pm := NewPropertyMap()
		keys := make([]rhinoStringID, 10000)
		for i := range keys {
			name := fmt.Sprintf("externGlobalSymbol_%05d", i)
			keys[i] = rhinoStringPoolIntern(name)
			pm.PutProperty(name, NewProperty(name, numType, false, NilNode))
		}
		pm.PutProperty("π", NewProperty("π", numType, false, NilNode))
		missID := rhinoStringPoolIntern("missingGlobalSymbol")

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			for i := 0; i < 100; i++ {
				if pm.GetOwnPropertyByID(keys[i*97%len(keys)]) == nil {
					b.Fatal("expected hit")
				}
				if pm.GetOwnPropertyByID(missID) != nil {
					b.Fatal("expected miss")
				}
			}
		}
	})

	b.Run("Build1800UnicodeKeys", func(b *testing.B) {
		names := make([]string, 1800)
		props := make([]*Property, 1800)
		for i := range names {
			names[i] = fmt.Sprintf("i18n_日本語_😀_%04d", i)
			props[i] = NewProperty(names[i], numType, false, NilNode)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			pm := NewPropertyMap()
			for i, name := range names {
				pm.PutProperty(name, props[i])
			}
		}
	})
}
