package main

// Where two builds of the same program differ, per function: the location half
// of the behavioural oracle's coverage rule (PLAN.md §8,
// docs/behavioural-oracle.md, "The coverage rule").
//
// VerifyAstAlphaIsomorphism stops at the first difference. DiffBuilds walks the
// two ASTs the same way but aligns them instead of failing: where the children
// of two matched nodes differ in number or shape, they are aligned by a
// longest common subsequence of name-insensitive subtree hashes, and the
// children left over between two aligned ones are paired by token where they
// can be (and compared further down) or reported as inserted or deleted. Every
// node of build B (fast) is then either matched to a node of build A (compat)
// with the same token, flags and literal values, or inside a reported change.
//
// Names are compared up to a renaming, as RenameVars and RenameProperties make
// them, after the walk, so that one changed function cannot shift the
// renaming seen by the others:
//   - A declared variable of A corresponds to one of B when their declarations
//     are matched nodes; otherwise to the variable of B its references meet
//     most often. Both must be declared in matched scopes. An undeclared name
//     (a global of the externs, such as window) corresponds only to itself.
//   - With props "renamed" (ADVANCED), an unquoted property name of A
//     corresponds to the name of B it meets most often, one to one; a name in
//     the fixed set (the externs' names, which renaming never produces or
//     changes) only to itself. With props "exact" (SIMPLE and WHITESPACE do
//     not rename properties) every property name only to itself.
//   - A label corresponds to the label of the matched LABEL statement.
// A reference whose names do not correspond is reported as a change.
//
// So a function of B with no change inside it (outside its nested functions)
// is the matched function of A up to a consistent renaming: still
// alpha-isomorphic. That covers renamed properties too, under ADVANCED's own
// assumption that the program does not reflect on renamed names.
//
// A change is a node of B (kind "node": the same token with different flags or
// literal; "insert": no counterpart in A), a node of A without counterpart in B
// ("delete"), or a NAME or property name whose names do not correspond ("name",
// "prop", "label"). Offsets are UTF-16 offsets into the build's text (compiled
// output is ASCII, so they are byte offsets too), half-open; the unit of a
// change is the innermost function around it, or the top-level statement.

import (
	"fmt"
	"hash/fnv"
	"math"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DiffChange is one place where the builds differ.
type DiffChange struct {
	Kind  string `json:"kind"`
	What  string `json:"what"`
	A     []int  `json:"a,omitempty"`
	B     []int  `json:"b,omitempty"`
	BAt   int    `json:"b_at"`
	UnitA int    `json:"unit_a"`
	UnitB int    `json:"unit_b"`
}

// DiffUnit is a function, or a top-level statement outside every function.
type DiffUnit struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Line  int    `json:"line"`
}

// DiffResult is the output of DiffBuilds.
type DiffResult struct {
	Identical  bool         `json:"identical"`
	Isomorphic bool         `json:"isomorphic"`
	FunctionsA int          `json:"functions_a"`
	FunctionsB int          `json:"functions_b"`
	ChangedA   int          `json:"changed_functions_a"`
	ChangedB   int          `json:"changed_functions_b"`
	Matched    int          `json:"matched_nodes"`
	NodesA     int          `json:"nodes_a"`
	NodesB     int          `json:"nodes_b"`
	Exact      bool         `json:"exact_positions"`
	Changes    []DiffChange `json:"changes"`
	UnitsA     []DiffUnit   `json:"units_a"`
	UnitsB     []DiffUnit   `json:"units_b"`
}

// DiffOptions selects how names are compared.
type DiffOptions struct {
	RenamedProps bool            // ADVANCED: properties may be renamed differently
	Fixed        map[string]bool // property names that keep their spelling
}

type build struct {
	text  string
	root  rhino.Node
	ins   []int // offsets in the parsed text of the newlines breakLongLines added
	exact bool
	shape map[rhino.Node]uint64
	ends  map[rhino.Node][2]int
	vars  map[rhino.Node]*jscomp.Var
	nodes int
	fns   int
	opts  *DiffOptions
}

// DiffBuilds compares build a (compat) with build b (fast).
func DiffBuilds(textA, textB string, opts DiffOptions) (*DiffResult, error) {
	res := &DiffResult{Identical: textA == textB, Changes: []DiffChange{}, UnitsA: []DiffUnit{}, UnitsB: []DiffUnit{}}
	A, err := loadBuild("compat.js", textA, &opts)
	if err != nil {
		return nil, fmt.Errorf("compat: %w", err)
	}
	B, err := loadBuild("fast.js", textB, &opts)
	if err != nil {
		return nil, fmt.Errorf("fast: %w", err)
	}
	res.FunctionsA, res.FunctionsB = A.fns, B.fns
	res.NodesA, res.NodesB = A.nodes, B.nodes
	res.Exact = A.exact && B.exact
	d := &differ{A: A, B: B, opts: &opts, ab: map[rhino.Node]rhino.Node{}, ba: map[rhino.Node]rhino.Node{}}
	d.match(A.root, B.root)
	d.resolveNames()
	res.Matched = len(d.ab)
	d.report(res)
	res.Isomorphic = len(res.Changes) == 0
	return res, nil
}

func loadBuild(name, text string, opts *DiffOptions) (*build, error) {
	b := &build{text: text, opts: opts, exact: true, ends: map[rhino.Node][2]int{}}
	root, comp, err := parseForIsomorphism(name, text)
	if err != nil {
		return nil, err
	}
	// Node columns stop at 4095 (rhino.Node_MAX_COLUMN_NUMBER), and compiled
	// output has longer lines: parse a copy with newlines after some commas
	// and semicolons, and keep it only if its AST is the same.
	if broken, ins := breakLongLines(text); len(ins) > 0 {
		b1 := &build{opts: opts}
		full := b1.hash(root, true)
		root2, comp2, err2 := parseForIsomorphism(name, broken)
		b2 := &build{opts: opts}
		if err2 == nil && b2.hash(root2, true) == full {
			root, comp = root2, comp2
			b.ins = ins
		} else {
			b.exact = false
		}
	}
	b.root = root
	b.vars = map[rhino.Node]*jscomp.Var{}
	b.resolve(root, nil, jscomp.NewSyntacticScopeCreator(comp))
	b.shape = map[rhino.Node]uint64{}
	b.hash(root, false)
	return b, nil
}

// breakLongLines returns text with a newline after a comma or semicolon (in
// code, not in a string, template, regular expression or comment) wherever a
// line would otherwise exceed 2000 characters, and the offsets of the added
// newlines in the returned text. A newline after a comma or semicolon never
// changes how code parses; the caller checks the AST anyway, since telling a
// regular expression from a division here is a guess.
func breakLongLines(text string) (string, []int) {
	const limit = 2000
	if !hasLongLine(text, limit) {
		return text, nil
	}
	var out strings.Builder
	var ins []int
	col := 0
	var braces []int // template nesting: brace depth at each ${
	depth := 0
	prev := byte(0) // last significant character in code
	prevWord := ""
	i := 0
	emit := func(s string) {
		out.WriteString(s)
		if k := strings.LastIndexByte(s, '\n'); k >= 0 {
			col = len(s) - k - 1
		} else {
			col += len(s)
		}
	}
	scanTemplate := func() { // at the character after ` or }
		for i < len(text) {
			c := text[i]
			if c == '\\' && i+1 < len(text) {
				emit(text[i : i+2])
				i += 2
				continue
			}
			if c == '`' {
				emit("`")
				i++
				prev = '`'
				return
			}
			if c == '$' && i+1 < len(text) && text[i+1] == '{' {
				emit("${")
				i += 2
				braces = append(braces, depth)
				depth++
				prev = '{'
				return
			}
			emit(text[i : i+1])
			i++
		}
	}
	for i < len(text) {
		c := text[i]
		switch {
		case c == '"' || c == '\'':
			j := i + 1
			for j < len(text) && text[j] != c && text[j] != '\n' {
				if text[j] == '\\' {
					j++
				}
				j++
			}
			if j < len(text) {
				j++
			}
			emit(text[i:j])
			i = j
			prev, prevWord = c, ""
		case c == '`':
			emit("`")
			i++
			scanTemplate()
			prevWord = ""
		case c == '/' && i+1 < len(text) && text[i+1] == '*':
			j := strings.Index(text[i+2:], "*/")
			if j < 0 {
				j = len(text)
			} else {
				j += i + 4
			}
			emit(text[i:j])
			i = j
		case c == '/' && i+1 < len(text) && text[i+1] == '/':
			j := strings.IndexByte(text[i:], '\n')
			if j < 0 {
				j = len(text)
			} else {
				j += i
			}
			emit(text[i:j])
			i = j
		case c == '/' && regexAllowed(prev, prevWord):
			j := i + 1
			inClass := false
			for j < len(text) && text[j] != '\n' {
				if text[j] == '\\' {
					j += 2
					continue
				}
				if text[j] == '[' {
					inClass = true
				} else if text[j] == ']' {
					inClass = false
				} else if text[j] == '/' && !inClass {
					break
				}
				j++
			}
			if j < len(text) && text[j] == '/' {
				j++
			}
			for j < len(text) && isIdentChar(text[j]) {
				j++
			}
			if j > len(text) {
				j = len(text)
			}
			emit(text[i:j])
			i = j
			prev, prevWord = ')', "" // a regular expression is an operand
		case c == '{':
			depth++
			emit("{")
			i++
			prev, prevWord = c, ""
		case c == '}':
			if len(braces) > 0 && braces[len(braces)-1] == depth-1 {
				braces = braces[:len(braces)-1]
				depth--
				emit("}")
				i++
				scanTemplate()
				prevWord = ""
				continue
			}
			depth--
			emit("}")
			i++
			prev, prevWord = c, ""
		case isIdentChar(c):
			j := i
			for j < len(text) && isIdentChar(text[j]) {
				j++
			}
			emit(text[i:j])
			prev, prevWord = text[j-1], text[i:j]
			i = j
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			emit(text[i : i+1])
			i++
		default:
			emit(text[i : i+1])
			i++
			prev, prevWord = c, ""
			if (c == ',' || c == ';') && col >= limit {
				ins = append(ins, out.Len())
				emit("\n")
			}
		}
	}
	return out.String(), ins
}

func hasLongLine(text string, limit int) bool {
	start := 0
	for start < len(text) {
		k := strings.IndexByte(text[start:], '\n')
		if k < 0 {
			return len(text)-start > limit
		}
		if k > limit {
			return true
		}
		start += k + 1
	}
	return false
}

func isIdentChar(c byte) bool {
	return c == '_' || c == '$' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}

var regexKeywords = map[string]bool{"return": true, "typeof": true, "instanceof": true, "in": true, "of": true,
	"new": true, "delete": true, "void": true, "throw": true, "case": true, "do": true, "else": true,
	"yield": true, "await": true}

func regexAllowed(prev byte, word string) bool {
	if word != "" {
		return regexKeywords[word]
	}
	switch prev {
	case ')', ']', '}', '"', '\'', '`':
		return false
	}
	return !isIdentChar(prev) || prev == 0
}

// pos returns the [start, end) offsets of the subtree n in the original text.
// A node's own position is not always its start (a GETPROP's is that of the
// property name), so the subtree's extent is taken over its children too.
func (b *build) pos(n rhino.Node) (int, int) {
	s, e := b.extent(n)
	return b.orig(s), b.orig(e)
}

func (b *build) extent(n rhino.Node) (int, int) {
	if x, ok := b.ends[n]; ok {
		return x[0], x[1]
	}
	start := n.GetSourceOffset()
	end := start + n.GetLength()
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		cs, ce := b.extent(c)
		if cs >= 0 && (start < 0 || cs < start) {
			start = cs
		}
		if ce > end {
			end = ce
		}
	}
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	b.ends[n] = [2]int{start, end}
	return start, end
}

// posShallow returns the offsets of n alone: the name of a NAME (whose own
// length covers its initializer) or property node, else the node's own range.
func (b *build) posShallow(n rhino.Node) (int, int) {
	start := n.GetSourceOffset()
	if start < 0 {
		return b.pos(n)
	}
	l := n.GetLength()
	if t := n.GetToken(); t == rhino.Token_NAME || (isPropToken(t) && !n.IsQuotedStringKey()) {
		l = len(n.GetString())
	}
	return b.orig(start), b.orig(start + l)
}

func (b *build) orig(off int) int {
	k := sort.SearchInts(b.ins, off)
	return off - k
}

func mix(h uint64, v uint64) uint64 {
	h ^= v + 0x9e3779b97f4a7c15 + (h << 6) + (h >> 2)
	return h
}

func strHash(s string) uint64 {
	f := fnv.New64a()
	f.Write([]byte(s))
	return f.Sum64()
}

func flagBits(n rhino.Node) uint64 {
	var f uint64
	if n.IsArrowFunction() {
		f |= 1
	}
	if n.IsAsyncFunction() {
		f |= 2
	}
	if n.IsGeneratorFunction() {
		f |= 4
	}
	if n.IsStaticMember() {
		f |= 8
	}
	if n.IsQuotedStringKey() {
		f |= 16
	}
	switch n.GetToken() {
	case rhino.Token_YIELD:
		if n.IsYieldAll() {
			f |= 32
		}
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_CALL:
		if n.IsOptionalChainStart() {
			f |= 64
		}
	case rhino.Token_SCRIPT, rhino.Token_BLOCK:
		if n.IsUseStrict() {
			f |= 128
		}
	}
	return f
}

func isPropToken(t rhino.Token) bool {
	switch t {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FIELD_DEF:
		return true
	}
	return false
}

// literal returns the literal value of n compared exactly (with ok), or ok=false.
func literal(n rhino.Node) (string, bool) {
	switch n.GetToken() {
	case rhino.Token_NUMBER:
		d := n.GetDouble()
		if d != d {
			return "NaN", true
		}
		return fmt.Sprintf("%x", math.Float64bits(d)), true
	case rhino.Token_BIGINT:
		return n.GetBigInt().String(), true
	case rhino.Token_STRINGLIT:
		return n.GetString(), true
	case rhino.Token_TEMPLATELIT_STRING:
		return n.GetRawString(), true
	}
	if isPropToken(n.GetToken()) && n.IsQuotedStringKey() {
		return n.GetString(), true
	}
	return "", false
}

// propFixed reports whether property node n's name must be equal in both builds.
func (b *build) propFixed(name string) bool {
	return !b.opts.RenamedProps || b.opts.Fixed[name]
}

// hash returns a hash of the subtree. With names, every name is included (the
// AST check of breakLongLines); without, only the names that cannot be
// renamed (undeclared names, fixed property names), so that the hash can
// align builds whose renaming differs.
func (b *build) hash(n rhino.Node, names bool) uint64 {
	t := n.GetToken()
	h := mix(uint64(t)+1, flagBits(n))
	if lit, ok := literal(n); ok {
		h = mix(h, strHash(lit))
	} else if isPropToken(t) || t == rhino.Token_LABEL_NAME {
		if s := n.GetString(); names || (isPropToken(t) && b.propFixed(s)) {
			h = mix(h, strHash(s))
		}
	} else if t == rhino.Token_NAME && (names || b.vars[n] == nil) {
		h = mix(h, strHash(n.GetString()))
	}
	if !names {
		b.nodes++
		if t == rhino.Token_FUNCTION {
			b.fns++
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		h = mix(h, b.hash(c, names))
	}
	if !names {
		b.shape[n] = h
	}
	return h
}

// resolve records the variable of every NAME node, creating scopes as
// NodeTraversal does (the same walk as astIsomorphismVerifier.walk).
func (b *build) resolve(n rhino.Node, scope *jscomp.Scope, sc *jscomp.SyntacticScopeCreator) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		name := n.GetFirstChild()
		params := name.GetNext()
		body := params.GetNext()
		isExpr := jscomp.NodeUtil_IsFunctionExpression(n)
		if !isExpr {
			b.resolve(name, scope, sc)
		}
		fs := sc.CreateScopeOfScope(n, scope)
		if isExpr {
			b.resolve(name, fs, sc)
		}
		b.resolve(params, fs, sc)
		b.resolve(body, fs, sc)
		return
	case rhino.Token_NAME:
		if s := n.GetString(); s != "" && scope != nil {
			b.vars[n] = scope.GetVar(s)
		}
	}
	if n.IsRoot() || jscomp.NodeUtil_CreatesBlockScope(n) {
		scope = sc.CreateScopeOfScope(n, scope)
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		b.resolve(c, scope, sc)
	}
}

type change struct {
	kind, what string
	a, b       rhino.Node // either may be nil (rhino.NilNode)
	bAt        rhino.Node // for a delete: the matched B node it was next to
	bAtEnd     bool       // bAt's end rather than its start
	shallow    bool       // a name: only the node itself, not its children
}

type differ struct {
	A, B    *build
	opts    *DiffOptions
	ab, ba  map[rhino.Node]rhino.Node
	changes []change
	names   [][2]rhino.Node
	props   [][2]rhino.Node
	labels  [][2]rhino.Node
}

func children(n rhino.Node) []rhino.Node {
	var out []rhino.Node
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		out = append(out, c)
	}
	return out
}

func (d *differ) match(nA, nB rhino.Node) {
	d.ab[nA] = nB
	d.ba[nB] = nA
	t := nA.GetToken()
	if fa, fb := flagBits(nA), flagBits(nB); fa != fb {
		d.changes = append(d.changes, change{kind: "node", what: fmt.Sprintf("%v flags %d vs %d", t, fa, fb), a: nA, b: nB})
		return
	}
	la, oka := literal(nA)
	lb, _ := literal(nB)
	if oka && la != lb {
		d.changes = append(d.changes, change{kind: "node", what: fmt.Sprintf("%v literal %.40q vs %.40q", t, la, lb), a: nA, b: nB})
		return
	}
	switch {
	case t == rhino.Token_NAME:
		d.names = append(d.names, [2]rhino.Node{nA, nB})
	case t == rhino.Token_LABEL_NAME:
		d.labels = append(d.labels, [2]rhino.Node{nA, nB})
	case isPropToken(t) && !nA.IsQuotedStringKey():
		d.props = append(d.props, [2]rhino.Node{nA, nB})
	}
	cA, cB := children(nA), children(nB)
	if len(cA) == len(cB) {
		same := true
		for i := range cA {
			if cA[i].GetToken() != cB[i].GetToken() || d.A.shape[cA[i]] != d.B.shape[cB[i]] {
				same = false
				break
			}
		}
		if same || fixedArity(t) {
			for i := range cA {
				if cA[i].GetToken() == cB[i].GetToken() {
					d.match(cA[i], cB[i])
				} else {
					d.changes = append(d.changes, change{kind: "node", what: fmt.Sprintf("%v vs %v", cA[i].GetToken(), cB[i].GetToken()), a: cA[i], b: cB[i]})
				}
			}
			return
		}
	}
	d.alignChildren(nA, nB, cA, cB)
}

// fixedArity reports whether the children of t are positional, so that two
// lists of the same length are compared position by position.
func fixedArity(t rhino.Token) bool {
	switch t {
	case rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_BLOCK, rhino.Token_MODULE_BODY, rhino.Token_PARAM_LIST,
		rhino.Token_CALL, rhino.Token_NEW, rhino.Token_OPTCHAIN_CALL, rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT,
		rhino.Token_CLASS_MEMBERS, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_TEMPLATELIT,
		rhino.Token_SWITCH_BODY, rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		return false
	}
	return true
}

func (d *differ) alignChildren(nA, nB rhino.Node, cA, cB []rhino.Node) {
	ka := make([]uint64, len(cA))
	kb := make([]uint64, len(cB))
	for i, c := range cA {
		ka[i] = d.A.shape[c]
	}
	for j, c := range cB {
		kb[j] = d.B.shape[c]
	}
	anchors := lcs(ka, kb)
	// Between two anchors, pair the leftovers by token, in order.
	matchedA := make([]int, len(cA))
	for i := range matchedA {
		matchedA[i] = -1
	}
	pi, pj := 0, 0
	gap := func(i2, j2 int) {
		ta := make([]uint64, i2-pi)
		tb := make([]uint64, j2-pj)
		for k := range ta {
			ta[k] = uint64(cA[pi+k].GetToken())
		}
		for k := range tb {
			tb[k] = uint64(cB[pj+k].GetToken())
		}
		for _, p := range lcs(ta, tb) {
			matchedA[pi+p[0]] = pj + p[1]
		}
	}
	for _, p := range anchors {
		gap(p[0], p[1])
		matchedA[p[0]] = p[1]
		pi, pj = p[0]+1, p[1]+1
	}
	gap(len(cA), len(cB))
	matchedB := make([]bool, len(cB))
	for i, j := range matchedA {
		if j >= 0 {
			matchedB[j] = true
			d.match(cA[i], cB[j])
		}
	}
	for j, ok := range matchedB {
		if !ok {
			d.changes = append(d.changes, change{kind: "insert", what: fmt.Sprintf("%v", cB[j].GetToken()), b: cB[j]})
		}
	}
	for i, j := range matchedA {
		if j >= 0 {
			continue
		}
		ch := change{kind: "delete", what: fmt.Sprintf("%v", cA[i].GetToken()), a: cA[i], bAt: nB, bAtEnd: true}
		for k := i + 1; k < len(cA); k++ {
			if matchedA[k] >= 0 {
				ch.bAt, ch.bAtEnd = cB[matchedA[k]], false
				break
			}
		}
		d.changes = append(d.changes, ch)
	}
}

// lcs returns the index pairs of a longest common subsequence of a and b
// (exact for up to 4M cells, otherwise anchored on the keys unique to both).
func lcs(a, b []uint64) [][2]int {
	var out [][2]int
	lo := 0
	for lo < len(a) && lo < len(b) && a[lo] == b[lo] {
		out = append(out, [2]int{lo, lo})
		lo++
	}
	ha, hb := len(a), len(b)
	var tail [][2]int
	for ha > lo && hb > lo && a[ha-1] == b[hb-1] {
		ha--
		hb--
		tail = append(tail, [2]int{ha, hb})
	}
	out = append(out, lcsCore(a[lo:ha], b[lo:hb], lo, lo)...)
	for k := len(tail) - 1; k >= 0; k-- {
		out = append(out, tail[k])
	}
	return out
}

func lcsCore(a, b []uint64, oa, ob int) [][2]int {
	n, m := len(a), len(b)
	if n == 0 || m == 0 {
		return nil
	}
	if n*m <= 4_000_000 {
		// dp[i][j]: LCS length of a[i:], b[j:].
		w := m + 1
		dp := make([]int32, (n+1)*w)
		for i := n - 1; i >= 0; i-- {
			for j := m - 1; j >= 0; j-- {
				if a[i] == b[j] {
					dp[i*w+j] = dp[(i+1)*w+j+1] + 1
				} else if x, y := dp[(i+1)*w+j], dp[i*w+j+1]; x >= y {
					dp[i*w+j] = x
				} else {
					dp[i*w+j] = y
				}
			}
		}
		var out [][2]int
		i, j := 0, 0
		for i < n && j < m {
			if a[i] == b[j] {
				out = append(out, [2]int{oa + i, ob + j})
				i++
				j++
			} else if dp[(i+1)*w+j] >= dp[i*w+j+1] {
				i++
			} else {
				j++
			}
		}
		return out
	}
	// Patience: anchor on the keys that occur once in each, in increasing order.
	ca, cb := map[uint64]int{}, map[uint64]int{}
	pa, pb := map[uint64]int{}, map[uint64]int{}
	for i, k := range a {
		ca[k]++
		pa[k] = i
	}
	for j, k := range b {
		cb[k]++
		pb[k] = j
	}
	var cand [][2]int
	for i, k := range a {
		if ca[k] == 1 && cb[k] == 1 {
			cand = append(cand, [2]int{i, pb[k]})
		}
	}
	anchors := lis(cand)
	if len(anchors) == 0 {
		// Nothing to anchor on: pair equal keys greedily, in order.
		var out [][2]int
		j := 0
		for i := 0; i < n && j < m; i++ {
			for k := j; k < m && k < j+64; k++ {
				if a[i] == b[k] {
					out = append(out, [2]int{oa + i, ob + k})
					j = k + 1
					break
				}
			}
		}
		return out
	}
	var out [][2]int
	pi, pj := 0, 0
	for _, p := range anchors {
		out = append(out, lcs2(a[pi:p[0]], b[pj:p[1]], oa+pi, ob+pj)...)
		out = append(out, [2]int{oa + p[0], ob + p[1]})
		pi, pj = p[0]+1, p[1]+1
	}
	return append(out, lcs2(a[pi:], b[pj:], oa+pi, ob+pj)...)
}

func lcs2(a, b []uint64, oa, ob int) [][2]int {
	var out [][2]int
	for _, p := range lcs(a, b) {
		out = append(out, [2]int{oa + p[0], ob + p[1]})
	}
	return out
}

// lis returns a longest subsequence of pairs (sorted by first) increasing in second.
func lis(p [][2]int) [][2]int {
	if len(p) == 0 {
		return nil
	}
	tails := []int{} // index into p of the smallest tail of each length
	prev := make([]int, len(p))
	for i := range p {
		k := sort.Search(len(tails), func(k int) bool { return p[tails[k]][1] >= p[i][1] })
		if k > 0 {
			prev[i] = tails[k-1]
		} else {
			prev[i] = -1
		}
		if k == len(tails) {
			tails = append(tails, i)
		} else {
			tails[k] = i
		}
	}
	out := make([][2]int, len(tails))
	for i, k := len(tails)-1, tails[len(tails)-1]; i >= 0; i, k = i-1, prev[k] {
		out[i] = p[k]
	}
	return out
}

type varPair struct{ a, b *jscomp.Var }
type strPair struct{ a, b string }

// resolveNames pairs the variables, property names and labels of the matched
// nodes and reports the references whose names do not correspond.
func (d *differ) resolveNames() {
	// Variables: by votes of the matched references, a matched pair of
	// declarations winning ties; both must be declared in matched scopes.
	va := map[*jscomp.Var]*jscomp.Var{}
	vb := map[*jscomp.Var]*jscomp.Var{}
	scopeOK := func(x, y *jscomp.Var) bool {
		return d.ab[x.GetScopeRoot()] == y.GetScopeRoot()
	}
	votes := map[varPair]int{}
	decl := map[varPair]bool{}
	var order []varPair
	for _, p := range d.names {
		x, y := d.A.vars[p[0]], d.B.vars[p[1]]
		if x == nil || y == nil || !scopeOK(x, y) {
			continue
		}
		k := varPair{x, y}
		if votes[k] == 0 {
			order = append(order, k)
		}
		votes[k]++
		if x.GetNameNode() == p[0] && y.GetNameNode() == p[1] {
			decl[k] = true
		}
	}
	sort.SliceStable(order, func(i, j int) bool {
		if votes[order[i]] != votes[order[j]] {
			return votes[order[i]] > votes[order[j]]
		}
		return decl[order[i]] && !decl[order[j]]
	})
	for _, k := range order {
		if _, ok := va[k.a]; ok {
			continue
		}
		if _, ok := vb[k.b]; ok {
			continue
		}
		va[k.a], vb[k.b] = k.b, k.a
	}
	for _, p := range d.names {
		x, y := d.A.vars[p[0]], d.B.vars[p[1]]
		sa, sb := p[0].GetString(), p[1].GetString()
		switch {
		case x == nil && y == nil:
			if sa != sb {
				d.nameChange("name", fmt.Sprintf("global %s vs %s", sa, sb), p)
			}
		case x == nil || y == nil:
			d.nameChange("name", fmt.Sprintf("%s (declared %v) vs %s (declared %v)", sa, x != nil, sb, y != nil), p)
		case va[x] != y:
			d.nameChange("name", fmt.Sprintf("variable %s vs %s", sa, sb), p)
		}
	}
	// Property names.
	pa := map[string]string{}
	pb := map[string]string{}
	pvotes := map[strPair]int{}
	var porder []strPair
	for _, p := range d.props {
		sa, sb := p[0].GetString(), p[1].GetString()
		if d.A.propFixed(sa) || d.B.propFixed(sb) {
			continue
		}
		k := strPair{sa, sb}
		if pvotes[k] == 0 {
			porder = append(porder, k)
		}
		pvotes[k]++
	}
	sort.SliceStable(porder, func(i, j int) bool {
		if pvotes[porder[i]] != pvotes[porder[j]] {
			return pvotes[porder[i]] > pvotes[porder[j]]
		}
		if porder[i].a != porder[j].a {
			return porder[i].a < porder[j].a
		}
		return porder[i].b < porder[j].b
	})
	for _, k := range porder {
		if _, ok := pa[k.a]; ok {
			continue
		}
		if _, ok := pb[k.b]; ok {
			continue
		}
		pa[k.a], pb[k.b] = k.b, k.a
	}
	for _, p := range d.props {
		sa, sb := p[0].GetString(), p[1].GetString()
		fa, fb := d.A.propFixed(sa), d.B.propFixed(sb)
		if fa || fb {
			if sa != sb {
				d.nameChange("prop", fmt.Sprintf("property %s vs %s", sa, sb), p)
			}
			continue
		}
		if pa[sa] != sb {
			d.nameChange("prop", fmt.Sprintf("property %s vs %s", sa, sb), p)
		}
	}
	// Labels: a label reference must name the matched LABEL statement.
	for _, p := range d.labels {
		la, lb := enclosingLabel(p[0]), enclosingLabel(p[1])
		if la.IsNil() || lb.IsNil() || d.ab[la] != lb {
			d.nameChange("label", fmt.Sprintf("label %s vs %s", p[0].GetString(), p[1].GetString()), p)
		}
	}
}

// enclosingLabel returns the LABEL statement a LABEL_NAME refers to (itself
// for the name of a LABEL).
func enclosingLabel(n rhino.Node) rhino.Node {
	name := n.GetString()
	if p := n.GetParent(); !p.IsNil() && p.GetToken() == rhino.Token_LABEL {
		return p
	}
	for p := n.GetParent(); !p.IsNil(); p = p.GetParent() {
		if p.GetToken() == rhino.Token_LABEL && p.GetFirstChild().GetString() == name {
			return p
		}
		if p.GetToken() == rhino.Token_FUNCTION {
			break
		}
	}
	return rhino.NilNode
}

func (d *differ) nameChange(kind, what string, p [2]rhino.Node) {
	d.changes = append(d.changes, change{kind: kind, what: what, a: p[0], b: p[1], shallow: true})
}

// unitOf returns the innermost function around n (n itself if it is one), or
// the top-level statement that holds it.
func unitOf(n rhino.Node) rhino.Node {
	for p := n; !p.IsNil(); p = p.GetParent() {
		if p.GetToken() == rhino.Token_FUNCTION {
			return p
		}
		if pp := p.GetParent(); pp.IsNil() || pp.GetToken() == rhino.Token_SCRIPT {
			return p
		}
	}
	return n
}

func functionName(n rhino.Node) string {
	if n.GetToken() != rhino.Token_FUNCTION {
		return fmt.Sprintf("<%v>", n.GetToken())
	}
	if s := n.GetFirstChild().GetString(); s != "" {
		return s
	}
	p := n.GetParent()
	switch p.GetToken() {
	case rhino.Token_NAME, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		return p.GetString()
	case rhino.Token_ASSIGN:
		if q, ok := p.GetFirstChild().QualifiedName(); ok {
			return q
		}
	}
	return "<anonymous>"
}

func (d *differ) report(res *DiffResult) {
	unitsA := map[rhino.Node]int{}
	unitsB := map[rhino.Node]int{}
	unit := func(b *build, m map[rhino.Node]int, units *[]DiffUnit, n rhino.Node) int {
		if n.IsNil() {
			return -1
		}
		u := unitOf(n)
		if k, ok := m[u]; ok {
			return k
		}
		s, e := b.pos(u)
		kind := "function"
		if u.GetToken() != rhino.Token_FUNCTION {
			kind = "toplevel"
		}
		m[u] = len(*units)
		*units = append(*units, DiffUnit{Kind: kind, Name: functionName(u), Start: s, End: e, Line: u.GetLineno()})
		return m[u]
	}
	span := func(b *build, n rhino.Node, shallow bool) []int {
		if n.IsNil() {
			return nil
		}
		var s, e int
		switch {
		case shallow:
			s, e = b.posShallow(n)
		default:
			s, e = b.pos(n)
		}
		return []int{s, e}
	}
	for _, c := range d.changes {
		dc := DiffChange{Kind: c.kind, What: c.what, A: span(d.A, c.a, c.shallow), B: span(d.B, c.b, c.shallow),
			UnitA: unit(d.A, unitsA, &res.UnitsA, c.a), UnitB: unit(d.B, unitsB, &res.UnitsB, c.b), BAt: -1}
		if c.kind == "delete" && !c.bAt.IsNil() {
			s, e := d.B.pos(c.bAt)
			if c.bAtEnd {
				dc.BAt = e
			} else {
				dc.BAt = s
			}
			dc.UnitB = unit(d.B, unitsB, &res.UnitsB, c.bAt)
			if c.bAtEnd && c.bAt.GetToken() == rhino.Token_FUNCTION {
				dc.UnitB = unit(d.B, unitsB, &res.UnitsB, c.bAt.GetLastChild())
			}
		}
		res.Changes = append(res.Changes, dc)
	}
	for _, u := range res.UnitsA {
		if u.Kind == "function" {
			res.ChangedA++
		}
	}
	for _, u := range res.UnitsB {
		if u.Kind == "function" {
			res.ChangedB++
		}
	}
}
