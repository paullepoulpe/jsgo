package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilMultiCollectionVectors.java
// against the pinned closure-compiler.jar (shaded Guava): seeded operation
// scripts on LinkedHashMultimap, MultimapBuilder.linkedHashKeys()
// .arrayListValues() and HashMultiset, with every result.

func multiCollectionEntries[K, V any](seq func(yield func(K, V) bool)) string {
	var parts []string
	for k, v := range seq {
		parts = append(parts, ToJavaString(k)+"="+ToJavaString(v))
	}
	return CollectionToString(parts)
}

func TestMultiCollectionsMatchGuava(t *testing.T) {
	rows := readTSV(t, "multi_collection.tsv")
	if len(rows) < 900 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		want := strings.Split(r[2], "|")
		var got []string
		ops := strings.Split(r[1], ";")
		switch r[0] {
		case "lhm":
			mm := NewLinkedHashMultimap[string, string]()
			for _, op := range ops {
				a := strings.Split(op, ":")
				switch a[0] {
				case "put":
					got = append(got, strconv.FormatBool(mm.Put(a[1], a[2])))
				case "rm":
					got = append(got, strconv.FormatBool(mm.Remove(a[1], a[2])))
				case "rmAll":
					got = append(got, CollectionToString(mm.RemoveAll(a[1])))
				case "get":
					got = append(got, CollectionToString(mm.Get(a[1])))
				case "ck":
					got = append(got, strconv.FormatBool(mm.ContainsKey(a[1])))
				case "ce":
					got = append(got, strconv.FormatBool(mm.ContainsEntry(a[1], a[2])))
				case "size":
					got = append(got, strconv.Itoa(mm.Size()))
				case "vals":
					var vs []string
					for v := range mm.Values() {
						vs = append(vs, v)
					}
					got = append(got, CollectionToString(vs))
				case "ents":
					got = append(got, multiCollectionEntries(mm.Entries()))
				case "keys":
					var ks []string
					for k := range mm.KeySet() {
						ks = append(ks, k)
					}
					got = append(got, CollectionToString(ks))
				case "str":
					got = append(got, mm.String())
				default:
					t.Fatalf("bad op %q", op)
				}
			}
		case "list":
			mm := ArrayListMultimap_Create[string, string]()
			for _, op := range ops {
				a := strings.Split(op, ":")
				switch a[0] {
				case "put":
					got = append(got, strconv.FormatBool(mm.Put(a[1], a[2])))
				case "rmAll":
					got = append(got, CollectionToString(mm.RemoveAll(a[1])))
				case "get":
					got = append(got, CollectionToString(mm.Get(a[1])))
				case "ck":
					got = append(got, strconv.FormatBool(mm.ContainsKey(a[1])))
				case "size":
					got = append(got, strconv.Itoa(mm.Size()))
				case "ents":
					got = append(got, multiCollectionEntries(mm.Entries()))
				case "keys":
					var ks []string
					for k := range mm.KeySet() {
						ks = append(ks, k)
					}
					got = append(got, CollectionToString(ks))
				case "str":
					got = append(got, mm.String())
				default:
					t.Fatalf("bad op %q", op)
				}
			}
		case "ms":
			continue
		default:
			t.Fatalf("bad row %q", r[0])
		}
		if strings.Join(got, "|") != r[2] {
			for i := range want {
				if i >= len(got) || got[i] != want[i] {
					t.Errorf("%s %s: op %d (%s) = %v, Guava %s", r[0], r[1], i, ops[i], got[min(i, len(got)-1)], want[i])
					break
				}
			}
		}
	}
}
