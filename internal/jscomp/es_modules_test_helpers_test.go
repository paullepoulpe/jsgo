package jscomp

// Go-only: test-helper — small helpers shared by the w5-es-modules test ports.

import "github.com/paullepoulpe/jsgo/internal/jutil"

// esModulesTest_mapOf is ImmutableMap.of(k1, v1, ...) for String keys and
// values.
func esModulesTest_mapOf(kv ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i], kv[i+1])
	}
	return m
}
