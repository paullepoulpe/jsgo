package p

import (
	"maps"
	"slices"
)

type byName map[string]int

func Flagged(m map[string]int, n byName) []string {
	var out []string
	for k := range m { // want
		out = append(out, k)
	}
	for k := range n { // want
		out = append(out, k)
	}
	return append(out, slices.Collect(maps.Keys(m))...) // want
}

func Generic[M ~map[string]int](m M) (s int) {
	for _, v := range m { // want
		s += v
	}
	return s
}

func Annotated(m map[string]int) (s int) {
	// order-independent: only sums the values
	for _, v := range m {
		s += v
	}
	for range m { // order-independent: only counts
		s++
	}
	//order-independent:
	for _, v := range m { // want (no reason given)
		s += v
	}
	return s
}

func NotMaps(xs []int, ch chan int) (s int) {
	for _, x := range xs {
		s += x
	}
	for x := range ch {
		s += x
	}
	return s
}
