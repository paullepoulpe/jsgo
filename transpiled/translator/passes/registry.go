// Canonical post-pass registry and CLI `-passes=...` flag parser.

package passes

import (
	"fmt"
	"strings"
)

// AllPasses returns all registered post-passes in canonical execution order:
//  1. native_slice_lowering   (lowers non-escaping ArrayList locals and unexported fields to native []T slices)
//  2. native_set_lowering     (lowers non-escaping HashSet/LinkedHashSet locals, fields, and Set-returning funcs to j2gort.OrderedSet[T])
//  3. native_map_lowering     (lowers non-escaping HashMap/LinkedHashMap/IdentityHashMap locals and fields to j2gort.OrderedMap[K, V] / CustomOrderedMap[K, V])
//  4. immutable_list_lowering (lowers ImmutableList fields/locals and local ImmutableList.Builder variables to native []T slices)
//  5. native_deque_lowering   (lowers non-escaping ArrayDeque/Deque/Queue locals and fields to ring-buffer j2gort.Deque[T] / *j2gort.Deque[T])
//  6. string_builder_lowering (lowers non-escaping append-only StringBuilder locals to stack j2gort.LocalStringBuilder)
//  7. slab_alloc_lowering     (lowers package-owned AST tree and graph node/edge allocations onto j2gort.Slab[T] fields on Owner)
//  8. node_traversal_lowering (pools jscomp.NodeTraversal and embeds [12]-deep scope buffers and reusable SyntacticScopeCreator)
//  9. dead_cast_cleanup       (strips redundant T(x) and AsIface[T](x) self-conversions)
//
// 10. abstract_base_devirt    (devirtualizes non-overridden abstract-class interface calls through *<I>Base)
// 11. j2gort_peephole         (eliminates mechanical j2gort Cast/AsIface/Box/StrVal/CheckNotNull*/StringEquals/IfaceEquals wrappers and SliceIterator loops)
func AllPasses() []Pass {
	return []Pass{
		NewNativeSliceLoweringPass(),
		NewNativeSetLoweringPass(),
		NewNativeMapLoweringPass(),
		NewImmutableListLoweringPass(),
		NewNativeDequeLoweringPass(),
		NewStringBuilderLoweringPass(),
		NewSlabAllocLoweringPass(),
		NewNodeTraversalLoweringPass(),
		NewDeadCastCleanupPass(),
		NewAbstractBaseDevirtPass(),
		NewJ2gortPeepholePass(),
	}
}

// PassNames returns the canonical names of all available passes.
func PassNames() []string {
	all := AllPasses()
	names := make([]string, len(all))
	for i, p := range all {
		names[i] = p.Name()
	}
	return names
}

// ParsePassesFlag parses the `-passes=...` CLI flag value (`all`, `none`, or comma-separated names).
func ParsePassesFlag(spec string) ([]Pass, error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" || trimmed == "none" {
		return nil, nil
	}
	if trimmed == "all" {
		return AllPasses(), nil
	}

	registry := make(map[string]Pass)
	for _, p := range AllPasses() {
		registry[p.Name()] = p
	}

	var selected []Pass
	seen := make(map[string]bool)
	for _, raw := range strings.Split(trimmed, ",") {
		name := strings.TrimSpace(raw)
		if name == "" || name == "none" {
			continue
		}
		if name == "all" {
			for _, p := range AllPasses() {
				if !seen[p.Name()] {
					seen[p.Name()] = true
					selected = append(selected, p)
				}
			}
			continue
		}
		p, ok := registry[name]
		if !ok {
			return nil, fmt.Errorf("unknown pass %q (valid passes: %s, all, none)", name, strings.Join(PassNames(), ", "))
		}
		if !seen[name] {
			seen[name] = true
			selected = append(selected, p)
		}
	}
	return selected, nil
}
