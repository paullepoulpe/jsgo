package rhino_test

// Go-only: language-adapter — small helpers for the ported jstype tests
// (Java literals that Go spells differently).

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
)

// ptr returns a pointer to v: a non-null Java String argument where the Go
// API takes *string.
func ptr[T any](v T) *T { return &v }

// mapOf is ImmutableMap.of(k1, v1, k2, v2, ...) for name to type maps.
func mapOf(keysAndValues ...any) *jutil.LinkedHashMap[string, JSType] {
	m := jutil.NewLinkedHashMap[string, JSType]()
	for i := 0; i < len(keysAndValues); i += 2 {
		m.Put(keysAndValues[i].(string), keysAndValues[i+1].(JSType))
	}
	return m
}

// templateMapOf is ImmutableMap.of(k1, v1, ...) for template type bindings.
func templateMapOf(keysAndValues ...any) *jutil.LinkedHashMap[*TemplateType, JSType] {
	m := jutil.NewLinkedHashMap[*TemplateType, JSType]()
	for i := 0; i < len(keysAndValues); i += 2 {
		m.Put(keysAndValues[i].(*TemplateType), keysAndValues[i+1].(JSType))
	}
	return m
}
