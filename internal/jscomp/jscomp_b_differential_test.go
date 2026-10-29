package jscomp

// Go-only: differential test — replays testdata/jscomp_b.tsv, written by
// tools/oracle/java/JscompBVectors.java against the pinned
// closure-compiler.jar, and compares every result with Java's. Each row
// carries its own inputs (no shared random generator).

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"iter"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const jscompBSafe = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ._-[]()'\"<>#@?+*$%&^~{}"

// jscompBEnc is the oracle's u(): "N" for null, else "u" + the UTF-16 units
// with every unit outside jscompBSafe written as \XXXX.
func jscompBEnc(s *string) string {
	if s == nil {
		return "N"
	}
	var sb strings.Builder
	sb.WriteByte('u')
	for _, c := range jutil.StringToUTF16(*s) {
		if c < 0x80 && strings.IndexByte(jscompBSafe, byte(c)) >= 0 {
			sb.WriteByte(byte(c))
		} else {
			fmt.Fprintf(&sb, "\\%04x", c)
		}
	}
	return sb.String()
}

func jscompBEncS(s string) string { return jscompBEnc(&s) }

// jscompBDec inverts jscompBEnc.
func jscompBDec(s string) *string {
	if s == "N" {
		return nil
	}
	if !strings.HasPrefix(s, "u") {
		panic("bad encoded string: " + s)
	}
	var units []uint16
	for i := 1; i < len(s); i++ {
		if s[i] == '\\' {
			v, err := strconv.ParseUint(s[i+1:i+5], 16, 16)
			if err != nil {
				panic(err)
			}
			units = append(units, uint16(v))
			i += 4
		} else {
			units = append(units, uint16(s[i]))
		}
	}
	r := jutil.StringFromUTF16(units)
	return &r
}

func jscompBDecS(s string) string { return *jscompBDec(s) }

// jscompBCatch runs body; a thrown Java exception becomes "!" + u(toString()).
func jscompBCatch(body func() string) (result string) {
	jutil.Try(func() { result = body() }, func(e jutil.Throwable) {
		result = "!" + jscompBEncS(e.Error())
	})
	return result
}

func jscompBNames(s string) *jutil.LinkedHashSet[string] {
	set := jutil.NewLinkedHashSet[string]()
	if s != "" {
		set.AddAll(strings.Split(s, ",")...)
	}
	return set
}

func jscompBChars(s string) ASCIICharSet {
	var set ASCIICharSet
	for _, ch := range jutil.StringToUTF16(s) {
		set.Add(ch)
	}
	return set
}

func jscompBSeed(s string) *JavaRandom {
	seed, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return NewJavaRandomWithSeed(seed)
}

func jscompBRunNameOps(ops string) string {
	var names []string
	var g NameGenerator
	r := jscompBCatch(func() string {
		for _, op := range strings.Split(ops, ";") {
			f := strings.Split(op, ":")
			switch f[0] {
			case "D0":
				g = NewDefaultNameGenerator()
			case "D3":
				g = NewDefaultNameGeneratorWithReservedChars(jscompBNames(f[1]), f[2], jscompBChars(f[3]))
			case "D4":
				g = NewDefaultNameGeneratorWithFirstAndNonFirstChars(jscompBNames(f[1]), f[2], jscompBChars(f[3]), jscompBChars(f[4]))
			case "R1":
				g = NewRandomNameGenerator(jscompBSeed(f[1]))
			case "R3":
				g = NewRandomNameGeneratorWithReserved(jscompBNames(f[1]), f[2], jscompBChars(f[3]), jscompBSeed(f[4]))
			case "R4":
				g = NewRandomNameGeneratorWithFirstAndNonFirstChars(jscompBNames(f[1]), f[2], jscompBChars(f[3]), jscompBChars(f[4]), jscompBSeed(f[5]))
			case "r3":
				g.Reset(jscompBNames(f[1]), f[2], jscompBChars(f[3]))
			case "r4":
				g.ResetWithFirstAndNonFirstChars(jscompBNames(f[1]), f[2], jscompBChars(f[3]), jscompBChars(f[4]))
			case "c":
				g = g.Clone(jscompBNames(f[1]), f[2], jscompBChars(f[3]))
			case "f":
				g.(*DefaultNameGenerator).Favors(jscompBDecS(f[1]))
			case "g":
				n, _ := strconv.Atoi(f[1])
				for i := 0; i < n; i++ {
					names = append(names, g.GenerateNextName())
				}
			default:
				panic("bad op " + op)
			}
		}
		return ""
	})
	if r != "" {
		names = append(names, r)
	}
	return strings.Join(names, ",")
}

type jscompBGraph struct {
	nodes []string
	edges [][2]string
}

func jscompBParseGraph(nodes, edges string) jscompBGraph {
	g := jscompBGraph{nodes: strings.Split(nodes, ",")}
	if edges != "" {
		for _, e := range strings.Split(edges, ",") {
			a, b, _ := strings.Cut(e, ">")
			g.edges = append(g.edges, [2]string{a, b})
		}
	}
	return g
}

func jscompBAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func (rg jscompBGraph) intGraph() *graph.LinkedDirectedGraph[int, any] {
	g := graph.LinkedDirectedGraph_Create[int, any]()
	for _, n := range rg.nodes {
		g.CreateNode(jscompBAtoi(n))
	}
	for _, e := range rg.edges {
		g.Connect(jscompBAtoi(e[0]), nil, jscompBAtoi(e[1]))
	}
	return g
}

func jscompBColor(directed, tie string, rg jscompBGraph) string {
	return jscompBCatch(func() string {
		var g graph.Graph[string, any]
		if directed == "1" {
			dg := graph.LinkedDirectedGraph_Create[string, any]()
			for _, n := range rg.nodes {
				dg.CreateNode("n" + n)
			}
			for _, e := range rg.edges {
				dg.Connect("n"+e[0], nil, "n"+e[1])
			}
			g = dg
		} else {
			ug := graph.LinkedUndirectedGraph_Create[string, any]()
			for _, n := range rg.nodes {
				ug.CreateNode("n" + n)
			}
			for _, e := range rg.edges {
				ug.Connect("n"+e[0], nil, "n"+e[1])
			}
			g = ug
		}
		var coloring graph.GraphColoring[string, any]
		switch tie {
		case "A":
			coloring = graph.NewGraphColoring_GreedyGraphColoringWithTieBreaker[string, any](g, jutil.StringCompareTo)
		case "D":
			coloring = graph.NewGraphColoring_GreedyGraphColoringWithTieBreaker[string, any](g, func(a, b string) int { return jutil.StringCompareTo(b, a) })
		default:
			coloring = graph.NewGraphColoring_GreedyGraphColoring[string, any](g)
		}
		count := coloring.Color()
		supers := make([]string, len(rg.nodes))
		for i, n := range rg.nodes {
			supers[i] = coloring.GetPartitionSuperNode("n" + n)
		}
		return strconv.Itoa(count) + ";" + strings.Join(supers, ",")
	})
}

func jscompBJoin[T any](seq iter.Seq[T]) string {
	var parts []string
	for x := range seq {
		parts = append(parts, jutil.ToJavaString(x))
	}
	return strings.Join(parts, ",")
}

func jscompBDominators(rg jscompBGraph, entryCol string) string {
	return jscompBCatch(func() string {
		entry := jscompBAtoi(entryCol)
		tree := graph.DominatorTree_Compute[int, any](rg.intGraph(), entry)
		var idoms, sizes, all []string
		for i := 0; i <= len(rg.nodes); i++ {
			if idom, ok := tree.GetImmediateDominator(i); ok {
				idoms = append(idoms, strconv.Itoa(idom))
			} else {
				idoms = append(idoms, "-")
			}
			sizes = append(sizes, strconv.Itoa(tree.GetDominatedSubtreeSize(i)))
		}
		for k, v := range tree.GetAllSubtreeSizes().All() {
			all = append(all, strconv.Itoa(k)+"="+strconv.Itoa(v))
		}
		return strings.Join(idoms, ",") + ";" + strings.Join(sizes, ",") + ";" + strings.Join(all, ",")
	})
}

func jscompBLca(rg jscompBGraph, rootsCol string) string {
	return jscompBCatch(func() string {
		g := rg.intGraph()
		roots := jutil.NewLinkedHashSet[int]()
		for _, r := range strings.Split(rootsCol, ",") {
			roots.Add(jscompBAtoi(r))
		}
		var finder graph.LowestCommonAncestorFinder[int, any] = graph.NewLowestCommonAncestorFinder[int, any](g)
		found := jscompBJoin(finder.FindAll(roots).All())
		again := jscompBJoin(finder.FindAll(roots).All())
		return found + ";" + again
	})
}

func jscompBEntries(m interface {
	All() iter.Seq2[string, string]
}) string {
	var es []string
	for k, v := range m.All() {
		es = append(es, jscompBEncS(k)+"="+jscompBEncS(v))
	}
	return strings.Join(es, ",")
}

func jscompBParseEntries(col string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	if col != "" {
		for _, e := range strings.Split(col, ",") {
			k, v, _ := strings.Cut(e, "=")
			m.Put(jscompBDecS(k), jscompBDecS(v))
		}
	}
	return m
}

func jscompBHex(b []byte) string { return "x" + hex.EncodeToString(b) }

func jscompBEncLookup(s string, ok bool) string {
	if !ok {
		return "N"
	}
	return jscompBEncS(s)
}

func jscompBVariableMap(col string) string {
	m := jscompBParseEntries(col)
	return jscompBCatch(func() string {
		vm := NewVariableMap(m)
		bytes := vm.ToBytes()
		back := jscompBCatch(func() string {
			parsed := VariableMap_FromBytes(bytes)
			return jscompBEntries(parsed.ToMap()) + ";" + jscompBEntries(parsed.GetOriginalNameToNewNameMap()) +
				";" + jscompBEntries(parsed.GetNewNameToOriginalNameMap())
		})
		var lookups []string
		for k, v := range m.All() {
			lookups = append(lookups, jscompBEncLookup(vm.LookupNewName(k))+"/"+jscompBEncLookup(vm.LookupSourceName(v)))
		}
		return jscompBHex(bytes) + ";" + strings.Join(lookups, ",") + ";" + back
	})
}

func jscompBVariableMapBytes(col string) string {
	b, err := hex.DecodeString(strings.TrimPrefix(col, "x"))
	if err != nil {
		panic(err)
	}
	return jscompBCatch(func() string {
		return jscompBEntries(VariableMap_FromBytes(b).ToMap())
	})
}

func jscompBMessageId(projectCol, meaningCol, partsCol string) string {
	return jscompBCatch(func() string {
		var parts []JsMessage_Part
		if partsCol != "" {
			for _, p := range strings.Split(partsCol, ",") {
				s := jscompBDecS(p[1:])
				switch p[0] {
				case 's':
					parts = append(parts, JsMessage_StringPart_Create(s))
				case 'j':
					parts = append(parts, JsMessage_PlaceholderReference_CreateForJsName(s))
				case 'c':
					parts = append(parts, JsMessage_PlaceholderReference_CreateForCanonicalName(s))
				}
			}
		}
		return NewGoogleJsMessageIdGenerator(jscompBDec(projectCol)).GenerateId(jscompBDec(meaningCol), parts)
	})
}

func jscompBAddNumber(bitsCol string) string {
	bits, err := strconv.ParseUint(bitsCol, 16, 64)
	if err != nil {
		panic(err)
	}
	return jscompBCatch(func() string {
		var c CodeConsumer
		c.AddNumber(math.Float64frombits(bits), rhino.NilNode)
		return jscompBEncS(c.GetCode())
	})
}

func TestJscompBDifferential(t *testing.T) {
	f, err := os.Open("testdata/jscomp_b.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)

	xid := NewXid()
	total := map[string]int{}
	failed := map[string]int{}
	line := 0
	for sc.Scan() {
		line++
		col := strings.Split(sc.Text(), "\t")
		kind := col[0]
		want := col[len(col)-1]
		var got string
		switch kind {
		case "ng":
			got = jscompBRunNameOps(col[1])
		case "color":
			got = jscompBColor(col[1], col[2], jscompBParseGraph(col[3], col[4]))
		case "uf":
			continue
		case "dom":
			got = jscompBDominators(jscompBParseGraph(col[1], col[2]), col[3])
		case "lca":
			got = jscompBLca(jscompBParseGraph(col[1], col[2]), col[3])
		case "vmap":
			got = jscompBVariableMap(col[1])
		case "vbytes":
			got = jscompBVariableMapBytes(col[1])
		case "msgid":
			got = jscompBMessageId(col[1], col[2], col[3])
		case "hash":
			got = strconv.FormatInt(JsMessage_Hash_Hash64(jscompBDec(col[1])), 10)
		case "ph":
			s := jscompBDecS(col[1])
			got = strings.Join([]string{
				strconv.FormatBool(JsMessage_IsLowerCamelCaseWithNumericSuffixes(s)),
				strconv.FormatBool(JsMessage_IsCanonicalPlaceholderNameFormat(s)),
				jscompBCatch(func() string { return jscompBEncS(JsMessage_ToLowerCamelCaseWithNumericSuffixes(s)) }),
				jscompBCatch(func() string {
					return jscompBEncS(JsMessage_PlaceholderReference_CreateForJsName(s).GetCanonicalPlaceholderName())
				}),
				jscompBCatch(func() string {
					return jscompBEncS(JsMessage_PlaceholderReference_CreateForCanonicalName(s).GetJsPlaceholderName())
				}),
			}, "\t")
			want = strings.Join(col[2:], "\t")
		case "xid":
			s := jscompBDecS(col[1])
			got = xid.Get(s) + "\t" + strconv.FormatInt(int64(xid.GetAsInt(s)), 10)
			want = col[2] + "\t" + col[3]
		case "path":
			a, b := jscompBDecS(col[2]), jscompBDecS(col[3])
			got = jscompBCatch(func() string {
				switch col[1] {
				case "C":
					return jscompBEncS(PathUtil_CollapseDots(a))
				case "A":
					return jscompBEncS(PathUtil_MakeAbsoluteWithRoot(a, b))
				}
				return jscompBEncS(PathUtil_MakeRelative(a, b))
			})
		case "esc":
			got = jscompBEncS(SourceCodeEscapers_JavascriptEscaper().Escape(jscompBDecS(col[1])))
		case "num":
			got = jscompBAddNumber(col[1])
		default:
			t.Fatalf("line %d: unknown kind %q", line, kind)
		}
		total[kind]++
		if got != want {
			failed[kind]++
			if failed[kind] <= 5 {
				t.Errorf("line %d (%s): input %q\n got: %s\nwant: %s", line, kind, col[1:len(col)-1], got, want)
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	kinds := make([]string, 0, len(total))
	for k := range total {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		t.Logf("%-7s %4d cases, %d diffs", k, total[k], failed[k])
	}
	if len(total) != 13 {
		t.Errorf("expected 13 kinds, got %d", len(total))
	}
}
