// Command astgen generates internal/ast/nodes.gen.go from
// internal/ast/grammar.txt (see that file for the grammar's syntax).
//
// Usage: go run ./tools/astgen [-grammar file] [-out file]
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"slices"
	"strings"
)

// mode says how a child appears in its parent.
type mode byte

const (
	required mode = iota // always there
	optional             // absent when missing (only at the end)
	emptyOr              // an EMPTY node when missing
	list                 // all the remaining children
)

type child struct {
	role     string // as written, without the leading _
	exported bool
	typ      string // a kind, a category, or Node
	mode     mode
	index    int // position among the children
}

type kind struct {
	name     string
	doc      []string
	tokens   []string
	context  bool // only appears in its parent's role: not recognized by token
	children []child
	payload  string // the accessor of the node's string, if any
	line     string // the grammar line, for the doc comment
}

type category struct {
	name    string
	doc     []string
	members []string // as written: kinds and categories
	kinds   []string // all the kinds, in grammar order
	line    string
}

type grammar struct {
	kinds      []*kind
	categories []*category
	kindByName map[string]*kind
	catByName  map[string]*category
}

func main() {
	grammarPath := flag.String("grammar", "internal/ast/grammar.txt", "the grammar")
	outPath := flag.String("out", "internal/ast/nodes.gen.go", "the generated file")
	flag.Parse()
	g, err := parse(*grammarPath)
	if err != nil {
		log.Fatal(err)
	}
	src := generate(g)
	formatted, err := format.Source(src)
	if err != nil {
		os.WriteFile(*outPath, src, 0o644)
		log.Fatalf("formatting the generated code: %v (unformatted code written to %s)", err, *outPath)
	}
	if err := os.WriteFile(*outPath, formatted, 0o644); err != nil {
		log.Fatal(err)
	}
}

// --- parsing ---

func parse(path string) (*grammar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g := &grammar{kindByName: map[string]*kind{}, catByName: map[string]*category{}}
	var doc []string
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
			doc = nil
			continue
		case strings.HasPrefix(line, "##"):
			doc = append(doc, strings.TrimSpace(strings.TrimPrefix(line, "##")))
			continue
		case strings.HasPrefix(line, "#"):
			doc = nil
			continue
		}
		var err error
		if name, rest, ok := strings.Cut(line, "="); ok {
			err = g.addCategory(strings.TrimSpace(name), rest, doc, line)
		} else {
			err = g.addKind(strings.Fields(line), doc)
		}
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %v", path, lineNo, err)
		}
		doc = nil
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return g, g.resolve()
}

func (g *grammar) addCategory(name, rest string, doc []string, line string) error {
	c := &category{name: name, doc: doc, line: strings.Join(strings.Fields(line), " ")}
	for _, m := range strings.Split(rest, "|") {
		c.members = append(c.members, strings.TrimSpace(m))
	}
	if g.has(name) {
		return fmt.Errorf("%s is defined twice", name)
	}
	g.categories = append(g.categories, c)
	g.catByName[name] = c
	return nil
}

func (g *grammar) addKind(fields []string, doc []string) error {
	if len(fields) < 2 {
		return fmt.Errorf("a kind needs a name and tokens")
	}
	k := &kind{name: fields[0], doc: doc}
	tokens := fields[1]
	if strings.HasSuffix(tokens, "!") {
		k.context = true
		tokens = strings.TrimSuffix(tokens, "!")
	}
	k.tokens = strings.Split(tokens, ",")
	k.line = k.name + " = " + strings.Join(k.tokens, " | ")
	if len(fields) > 2 {
		k.line += "  " + strings.Join(fields[2:], " ")
	}
	for _, f := range fields[2:] {
		if strings.HasPrefix(f, "$") {
			if k.payload != "" {
				return fmt.Errorf("%s: two string accessors", k.name)
			}
			k.payload = strings.TrimPrefix(f, "$")
			continue
		}
		role, typ, ok := strings.Cut(f, ":")
		if !ok {
			return fmt.Errorf("%s: child %q is not role:Type", k.name, f)
		}
		c := child{role: role, exported: true, typ: typ, mode: required, index: len(k.children)}
		if strings.HasPrefix(role, "_") {
			c.role, c.exported = role[1:], false
		}
		switch typ[len(typ)-1] {
		case '?':
			c.mode, c.typ = optional, typ[:len(typ)-1]
		case '~':
			c.mode, c.typ = emptyOr, typ[:len(typ)-1]
		case '*':
			c.mode, c.typ = list, typ[:len(typ)-1]
		}
		k.children = append(k.children, c)
	}
	for i, c := range k.children {
		last := i == len(k.children)-1
		if (c.mode == optional || c.mode == list) && !last {
			return fmt.Errorf("%s: %s must be the last child", k.name, c.role)
		}
	}
	if g.has(k.name) {
		return fmt.Errorf("%s is defined twice", k.name)
	}
	g.kinds = append(g.kinds, k)
	g.kindByName[k.name] = k
	return nil
}

func (g *grammar) has(name string) bool {
	return g.kindByName[name] != nil || g.catByName[name] != nil
}

// resolve checks the references and computes each category's kinds.
func (g *grammar) resolve() error {
	for _, k := range g.kinds {
		for _, c := range k.children {
			if c.typ != "Node" && !g.has(c.typ) {
				return fmt.Errorf("%s.%s: unknown type %s", k.name, c.role, c.typ)
			}
		}
	}
	byToken := map[string]string{}
	for _, k := range g.kinds {
		if k.context {
			continue
		}
		for _, t := range k.tokens {
			if other, ok := byToken[t]; ok {
				return fmt.Errorf("token %s is in both %s and %s", t, other, k.name)
			}
			byToken[t] = k.name
		}
	}
	var expand func(name string, seen map[string]bool) ([]string, error)
	expand = func(name string, seen map[string]bool) ([]string, error) {
		if k := g.kindByName[name]; k != nil {
			if k.context {
				return nil, fmt.Errorf("%s is a context kind and cannot be in a category", name)
			}
			return []string{name}, nil
		}
		c := g.catByName[name]
		if c == nil {
			return nil, fmt.Errorf("unknown kind or category %s", name)
		}
		if seen[name] {
			return nil, fmt.Errorf("category %s contains itself", name)
		}
		seen[name] = true
		var out []string
		for _, m := range c.members {
			ks, err := expand(m, seen)
			if err != nil {
				return nil, err
			}
			out = append(out, ks...)
		}
		delete(seen, name)
		return out, nil
	}
	for _, c := range g.categories {
		ks, err := expand(c.name, map[string]bool{})
		if err != nil {
			return err
		}
		set := map[string]bool{}
		for _, k := range ks {
			set[k] = true
		}
		for _, k := range g.kinds { // grammar order, no duplicates
			if set[k.name] {
				c.kinds = append(c.kinds, k.name)
			}
		}
	}
	return nil
}

// isSubset reports whether every kind of category a is in category b.
func (g *grammar) isSubset(a, b *category) bool {
	for _, k := range a.kinds {
		if !slices.Contains(b.kinds, k) {
			return false
		}
	}
	return true
}

// --- generation ---

type gen struct {
	g   *grammar
	buf bytes.Buffer
}

func (w *gen) p(format string, args ...any) {
	fmt.Fprintf(&w.buf, format, args...)
	w.buf.WriteByte('\n')
}

func generate(g *grammar) []byte {
	w := &gen{g: g}
	w.p("// Code generated by tools/astgen from grammar.txt; DO NOT EDIT.")
	w.p("")
	w.p("package ast")
	w.p("")
	w.p("import (")
	w.p("\t\"iter\"")
	w.p("")
	w.p("\t\"github.com/paullepoulpe/jsgo/internal/rhino\"")
	w.p(")")
	w.p("")
	for _, k := range g.kinds {
		w.kind(k)
	}
	for _, c := range g.categories {
		w.category(c)
	}
	w.elementNavigation()
	w.visitor()
	w.tables()
	return w.buf.Bytes()
}

func (w *gen) doc(lines []string) {
	for _, l := range lines {
		w.p("// %s", l)
	}
}

func tokenConst(t string) string { return "rhino.Token_" + t }

// childExpr returns the Go expression for the child at index i of x.node.
func childExpr(k *kind, i int) string {
	last := i == len(k.children)-1
	if last && k.children[i].mode == required && i > 1 {
		return "x.node.GetLastChild()"
	}
	switch i {
	case 0:
		return "x.node.GetFirstChild()"
	case 1:
		return "x.node.GetSecondChild()"
	}
	return fmt.Sprintf("x.node.GetChildAtIndex(%d)", i)
}

func upper(s string) string { return strings.ToUpper(s[:1]) + s[1:] }

// a returns name with its indefinite article.
func a(name string) string {
	if strings.ContainsRune("AEIOU", rune(name[0])) {
		return "an " + name
	}
	return "a " + name
}

func accessorName(c child) string {
	if c.exported {
		return upper(c.role)
	}
	return c.role
}

// singular returns the name of one element of a list role.
func singular(role string) string {
	if role == "children" {
		return "child"
	}
	return strings.TrimSuffix(role, "s")
}

// wrap returns the Go expression that gives the child expression e the
// type typ.
func wrap(typ, e string) string {
	if typ == "Node" {
		return e
	}
	return typ + "{" + e + "}"
}

func goType(typ string) string {
	if typ == "Node" {
		return "rhino.Node"
	}
	return typ
}

func unwrap(typ, v string) string {
	if typ == "Node" {
		return v
	}
	return v + ".node"
}

var goKeywords = map[string]bool{"break": true, "case": true, "chan": true, "const": true, "continue": true, "default": true,
	"defer": true, "else": true, "fallthrough": true, "for": true, "func": true, "go": true, "goto": true, "if": true,
	"import": true, "interface": true, "map": true, "package": true, "range": true, "return": true, "select": true,
	"struct": true, "switch": true, "type": true, "var": true}

func paramName(role string) string {
	if goKeywords[role] {
		return role + "_"
	}
	return role
}

func (w *gen) kind(k *kind) {
	w.doc(k.doc)
	w.p("//")
	w.p("// Grammar: %s", k.line)
	w.p("type %s struct{ node rhino.Node }", k.name)
	w.p("")
	if !k.context {
		w.p("// As%s returns n as %s; ok is false if n is not one (or is the nil node).", k.name, a(k.name))
		w.p("func As%s(n rhino.Node) (x %s, ok bool) {", k.name, k.name)
		if len(k.tokens) == 1 {
			w.p("if n == rhino.NilNode || n.GetToken() != %s { return %s{}, false }", tokenConst(k.tokens[0]), k.name)
			w.p("return %s{n}, true", k.name)
		} else {
			w.p("if n == rhino.NilNode { return %s{}, false }", k.name)
			w.p("switch n.GetToken() {")
			w.p("case %s:", joinTokens(k.tokens))
			w.p("return %s{n}, true", k.name)
			w.p("}")
			w.p("return %s{}, false", k.name)
		}
		w.p("}")
		w.p("")
	}
	w.p("// Node returns x as a rhino.Node.")
	w.p("func (x %s) Node() rhino.Node { return x.node }", k.name)
	w.p("")
	w.p("// Is reports whether x is the node n.")
	w.p("func (x %s) Is(n rhino.Node) bool { return x.node == n }", k.name)
	w.p("")
	if len(k.tokens) > 1 {
		w.p("// Token returns which of %s x is.", strings.Join(k.tokens, ", "))
		w.p("func (x %s) Token() rhino.Token { return x.node.GetToken() }", k.name)
		w.p("")
	}
	if k.payload != "" {
		w.p("// %s returns the node's string.", k.payload)
		w.p("func (x %s) %s() string { return x.node.GetString() }", k.name, k.payload)
		w.p("")
		w.p("// Set%s sets the node's string.", k.payload)
		w.p("func (x %s) Set%s(s string) { x.node.SetString(s) }", k.name, k.payload)
		w.p("")
	}
	for _, c := range k.children {
		w.childAccessors(k, c)
	}
	w.constructor(k)
	// Upcasts to the categories that contain k.
	for _, c := range w.g.categories {
		if slices.Contains(c.kinds, k.name) {
			w.p("// As%s returns x as %s.", c.name, a(c.name))
			w.p("func (x %s) As%s() %s { return %s{x.node} }", k.name, c.name, c.name, c.name)
			w.p("")
		}
	}
}

func (w *gen) childAccessors(k *kind, c child) {
	name := accessorName(c)
	typ := goType(c.typ)
	e := childExpr(k, c.index)
	switch c.mode {
	case required:
		w.p("// %s returns the %s of x.", name, c.role)
		w.p("func (x %s) %s() %s { return %s }", k.name, name, typ, wrap(c.typ, e))
		w.p("")
		if c.exported {
			w.p("// Set%s replaces the %s of x with v (which must not be in a tree).", name, c.role)
			w.p("func (x %s) Set%s(v %s) { %s.ReplaceWith(%s) }", k.name, name, typ, e, unwrap(c.typ, "v"))
			w.p("")
		}
	case optional:
		w.p("// %s returns the %s of x; ok is false if x has none.", name, c.role)
		w.p("func (x %s) %s() (v %s, ok bool) {", k.name, name, typ)
		w.p("c := %s", e)
		w.p("return %s, c != rhino.NilNode", wrap(c.typ, "c"))
		w.p("}")
		w.p("")
		if c.exported {
			w.p("// Set%s sets the %s of x to v (which must not be in a tree).", name, c.role)
			w.p("func (x %s) Set%s(v %s) {", k.name, name, typ)
			w.p("if c := %s; c != rhino.NilNode { c.ReplaceWith(%s) } else { x.node.AddChildToBack(%s) }", e, unwrap(c.typ, "v"), unwrap(c.typ, "v"))
			w.p("}")
			w.p("")
			w.p("// Remove%s detaches the %s of x, if it has one.", name, c.role)
			w.p("func (x %s) Remove%s() {", k.name, name)
			w.p("if c := %s; c != rhino.NilNode { c.Detach() }", e)
			w.p("}")
			w.p("")
		}
	case emptyOr:
		w.p("// %s returns the %s of x; ok is false if x has none (the child is an EMPTY node).", name, c.role)
		w.p("func (x %s) %s() (v %s, ok bool) {", k.name, name, typ)
		w.p("c := %s", e)
		w.p("if c.IsEmpty() { return %s, false }", zero(c.typ))
		w.p("return %s, true", wrap(c.typ, "c"))
		w.p("}")
		w.p("")
		if c.exported {
			w.p("// Set%s replaces the %s of x (or the EMPTY node in its place) with v (which must not be in a tree).", name, c.role)
			w.p("func (x %s) Set%s(v %s) { %s.ReplaceWith(%s) }", k.name, name, typ, e, unwrap(c.typ, "v"))
			w.p("")
		}
	case list:
		w.listAccessors(k.name, c)
	}
}

// indexExpr returns the Go expression for the child at index i of x.node,
// by position only (unlike childExpr, never as the last child).
func indexExpr(i int) string {
	switch i {
	case 0:
		return "x.node.GetFirstChild()"
	case 1:
		return "x.node.GetSecondChild()"
	}
	return fmt.Sprintf("x.node.GetChildAtIndex(%d)", i)
}

// listAccessors writes the accessors of the list c for the type recv (a
// kind, or a category whose kinds all have c at the same index).
func (w *gen) listAccessors(recv string, c child) {
	name := accessorName(c)
	typ := goType(c.typ)
	one := upper(singular(c.role))
	start := "x.node.GetFirstChild()"
	if c.index > 0 {
		start = indexExpr(c.index-1) + ".GetNext()"
	}
	w.p("// %s iterates over the %s of x. Like a for loop over GetNext, it reads the next", name, c.role)
	w.p("// element after the loop body ran: a body that moves the current element elsewhere")
	w.p("// ends the loop early.")
	w.p("func (x %s) %s() iter.Seq[%s] {", recv, name, typ)
	w.p("return func(yield func(%s) bool) {", typ)
	w.p("for c := %s; c != rhino.NilNode; c = c.GetNext() {", start)
	w.p("if !yield(%s) { return }", wrap(c.typ, "c"))
	w.p("}")
	w.p("}")
	w.p("}")
	w.p("")
	w.p("// First%s returns the first of the %s of x; ok is false if there is none.", one, c.role)
	w.p("func (x %s) First%s() (v %s, ok bool) {", recv, one, typ)
	w.p("c := %s", start)
	w.p("return %s, c != rhino.NilNode", wrap(c.typ, "c"))
	w.p("}")
	w.p("")
	w.p("// Last%s returns the last of the %s of x; ok is false if there is none.", one, c.role)
	w.p("func (x %s) Last%s() (v %s, ok bool) {", recv, one, typ)
	if c.index == 0 {
		w.p("c := x.node.GetLastChild()")
	} else {
		w.p("if %s == rhino.NilNode { return %s, false }", start, zero(c.typ))
		w.p("c := x.node.GetLastChild()")
	}
	w.p("return %s, c != rhino.NilNode", wrap(c.typ, "c"))
	w.p("}")
	w.p("")
	w.p("// Has%s reports whether x has any %s.", upper(c.role), c.role)
	w.p("func (x %s) Has%s() bool { return %s != rhino.NilNode }", recv, upper(c.role), start)
	w.p("")
	w.p("// Append%s adds v (which must not be in a tree) after the %s of x.", one, c.role)
	w.p("func (x %s) Append%s(v %s) { x.node.AddChildToBack(%s) }", recv, one, typ, unwrap(c.typ, "v"))
	w.p("")
	if c.index == 0 {
		w.p("// Prepend%s adds v (which must not be in a tree) before the %s of x.", one, c.role)
		w.p("func (x %s) Prepend%s(v %s) { x.node.AddChildToFront(%s) }", recv, one, typ, unwrap(c.typ, "v"))
		w.p("")
	}
}

func zero(typ string) string {
	if typ == "Node" {
		return "rhino.NilNode"
	}
	return typ + "{}"
}

// constructor writes NewKind for the kinds that need no string and have no
// EMPTY placeholders or unexported children: the required children are
// parameters, a list is the variadic last one, and optional children are
// left out (use SetRole).
func (w *gen) constructor(k *kind) {
	if k.context || k.payload != "" {
		return
	}
	var params, adds []string
	if len(k.tokens) > 1 {
		params = append(params, "token rhino.Token")
	}
	for _, c := range k.children {
		if !c.exported {
			return // custom.go has the constructor
		}
		switch c.mode {
		case emptyOr:
			return
		case required:
			params = append(params, paramName(c.role)+" "+goType(c.typ))
			adds = append(adds, fmt.Sprintf("n.AddChildToBack(%s)", unwrap(c.typ, paramName(c.role))))
		case list:
			params = append(params, paramName(c.role)+" ..."+goType(c.typ))
			adds = append(adds, fmt.Sprintf("for _, v := range %s { n.AddChildToBack(%s) }", paramName(c.role), unwrap(c.typ, "v")))
		}
	}
	w.p("// New%s returns a new %s.", k.name, k.name)
	w.p("func New%s(%s) %s {", k.name, strings.Join(params, ", "), k.name)
	if len(k.tokens) > 1 {
		w.p("switch token {")
		w.p("case %s:", joinTokens(k.tokens))
		w.p("default:")
		w.p("panic(\"ast.New%s: \" + token.String() + \" is not %s\")", k.name, a(k.name))
		w.p("}")
		w.p("n := rhino.NewNode(token)")
	} else {
		w.p("n := rhino.NewNode(%s)", tokenConst(k.tokens[0]))
	}
	for _, a := range adds {
		w.p("%s", a)
	}
	w.p("return %s{n}", k.name)
	w.p("}")
	w.p("")
}

func joinTokens(tokens []string) string {
	var out []string
	for _, t := range tokens {
		out = append(out, tokenConst(t))
	}
	return strings.Join(out, ", ")
}

func (w *gen) categoryTokens(c *category) []string {
	var out []string
	for _, kn := range c.kinds {
		out = append(out, w.g.kindByName[kn].tokens...)
	}
	return out
}

func (w *gen) category(c *category) {
	g := w.g
	w.doc(c.doc)
	w.p("//")
	w.p("// Grammar: %s", c.line)
	w.p("type %s struct{ node rhino.Node }", c.name)
	w.p("")
	w.p("// As%s returns n as %s; ok is false if n is not one (or is the nil node).", c.name, a(c.name))
	w.p("func As%s(n rhino.Node) (x %s, ok bool) {", c.name, c.name)
	w.p("if n == rhino.NilNode { return %s{}, false }", c.name)
	w.p("switch n.GetToken() {")
	w.p("case %s:", joinTokens(w.categoryTokens(c)))
	w.p("return %s{n}, true", c.name)
	w.p("}")
	w.p("return %s{}, false", c.name)
	w.p("}")
	w.p("")
	w.p("// Node returns x as a rhino.Node.")
	w.p("func (x %s) Node() rhino.Node { return x.node }", c.name)
	w.p("")
	w.p("// Is reports whether x is the node n.")
	w.p("func (x %s) Is(n rhino.Node) bool { return x.node == n }", c.name)
	w.p("")
	for _, kn := range c.kinds {
		w.p("// As%s returns x as %s; ok is false if it is not one.", kn, a(kn))
		w.p("func (x %s) As%s() (%s, bool) { return As%s(x.node) }", c.name, kn, kn, kn)
		w.p("")
	}
	for _, d := range g.categories {
		if d == c {
			continue
		}
		if g.isSubset(d, c) {
			w.p("// As%s returns x as %s; ok is false if it is not one.", d.name, a(d.name))
			w.p("func (x %s) As%s() (%s, bool) { return As%s(x.node) }", c.name, d.name, d.name, d.name)
			w.p("")
		} else if g.isSubset(c, d) {
			w.p("// As%s returns x as %s.", d.name, a(d.name))
			w.p("func (x %s) As%s() %s { return %s{x.node} }", c.name, d.name, d.name, d.name)
			w.p("")
		}
	}
	w.sharedAccessors(c)
}

// sharedAccessors gives a category the accessors that all its kinds have
// in the same place: at the same index (e.g. the object of GetProp and
// GetElem, the args of Call and New), or last (e.g. the body of Case and
// DefaultCase).
func (w *gen) sharedAccessors(c *category) {
	first := w.g.kindByName[c.kinds[0]]
	for _, ch := range first.children {
		if !ch.exported || ch.mode == emptyOr {
			continue
		}
		sameIndex, allLast := true, positionKey(first, ch.index) == "last"
		shared := true
		for _, kn := range c.kinds[1:] {
			k := w.g.kindByName[kn]
			found := false
			for _, o := range k.children {
				if o.role == ch.role && o.exported && o.typ == ch.typ && o.mode == ch.mode {
					found = true
					sameIndex = sameIndex && o.index == ch.index
					allLast = allLast && positionKey(k, o.index) == "last"
				}
			}
			if !found {
				shared = false
				break
			}
		}
		// Read the child by its index if it has the same one in every kind,
		// else as the last child if it is last in every kind.
		if !shared || !(sameIndex || allLast) {
			continue
		}
		if ch.mode == list { // never "last", so here at the same index
			w.listAccessors(c.name, ch)
			continue
		}
		e := "x.node.GetLastChild()"
		if sameIndex {
			e = indexExpr(ch.index)
		}
		name := accessorName(ch)
		if ch.mode == required {
			w.p("// %s returns the %s of x.", name, ch.role)
			w.p("func (x %s) %s() %s { return %s }", c.name, name, goType(ch.typ), wrap(ch.typ, e))
		} else {
			w.p("// %s returns the %s of x; ok is false if x has none.", name, ch.role)
			w.p("func (x %s) %s() (v %s, ok bool) {", c.name, name, goType(ch.typ))
			w.p("c := %s", e)
			w.p("return %s, c != rhino.NilNode", wrap(ch.typ, "c"))
			w.p("}")
		}
		w.p("")
	}
	if first.payload != "" {
		shared := true
		for _, kn := range c.kinds[1:] {
			if w.g.kindByName[kn].payload != first.payload {
				shared = false
			}
		}
		if shared {
			w.p("// %s returns the node's string.", first.payload)
			w.p("func (x %s) %s() string { return x.node.GetString() }", c.name, first.payload)
			w.p("")
		}
	}
}

// positionKey describes where child i of k is: "last" when k always has
// the same children and i is the last one, else its index.
func positionKey(k *kind, i int) string {
	for _, c := range k.children {
		if c.mode != required && c.mode != emptyOr {
			return fmt.Sprint(i)
		}
	}
	if i == len(k.children)-1 {
		return "last"
	}
	return fmt.Sprint(i)
}

// elementNavigation gives the element types of lists Next and Previous.
func (w *gen) elementNavigation() {
	type use struct{ start0, other bool }
	uses := map[string]*use{}
	var order []string
	for _, k := range w.g.kinds {
		for _, c := range k.children {
			if c.mode != list || c.typ == "Node" {
				continue
			}
			u := uses[c.typ]
			if u == nil {
				u = &use{}
				uses[c.typ] = u
				order = append(order, c.typ)
			}
			if c.index == 0 {
				u.start0 = true
			} else {
				u.other = true
			}
		}
	}
	for _, t := range order {
		check := "As" + t + "(c)"
		if k := w.g.kindByName[t]; k != nil && k.context {
			continue
		}
		w.p("// Next returns the element after x in its list; ok is false if x is the last one.")
		w.p("func (x %s) Next() (%s, bool) { c := x.node.GetNext(); return %s }", t, t, check)
		w.p("")
		if !uses[t].other {
			w.p("// Previous returns the element before x in its list; ok is false if x is the first one.")
			w.p("func (x %s) Previous() (%s, bool) { c := x.node.GetPrevious(); return %s }", t, t, check)
			w.p("")
		}
	}
}

// visitor writes Visitor: a handler per kind, and the dispatch.
func (w *gen) visitor() {
	var kinds []*kind
	for _, k := range w.g.kinds {
		if !k.context {
			kinds = append(kinds, k)
		}
	}
	w.p("// Visitor calls, for each node it is given, the handler of the node's kind,")
	w.p("// if it has one. T is the type of the traversal (jscomp passes")
	w.p("// *NodeTraversal): see jscomp.NodeTraversal_TraverseVisitor.")
	w.p("type Visitor[T any] struct {")
	for _, k := range kinds {
		w.p("%s func(t T, x %s)", k.name, k.name)
	}
	w.p("}")
	w.p("")
	w.p("// Visit calls the handler of n's kind, if v has one.")
	w.p("func (v *Visitor[T]) Visit(t T, n rhino.Node) {")
	w.p("switch n.GetToken() {")
	for _, k := range kinds {
		w.p("case %s:", joinTokens(k.tokens))
		w.p("if v.%s != nil { v.%s(t, %s{n}) }", k.name, k.name, k.name)
	}
	w.p("}")
	w.p("}")
	w.p("")
	w.p("// Tokens returns the tokens of the kinds v has a handler for.")
	w.p("func (v *Visitor[T]) Tokens() TokenSet {")
	w.p("var s TokenSet")
	for _, k := range kinds {
		w.p("if v.%s != nil {", k.name)
		for _, t := range k.tokens {
			w.p("s.Add(%s)", tokenConst(t))
		}
		w.p("}")
	}
	w.p("return s")
	w.p("}")
	w.p("")
}

// tables writes the grammar as data, for Check (check.go).
func (w *gen) tables() {
	modes := map[mode]string{required: "required", optional: "optional", emptyOr: "emptyOr", list: "list"}
	w.p("// grammarKinds is the grammar's node kinds, for Check.")
	w.p("var grammarKinds = []kindSpec{")
	for _, k := range w.g.kinds {
		var ch []string
		for _, c := range k.children {
			ch = append(ch, fmt.Sprintf("{role: %q, typ: %q, mode: %s}", c.role, c.typ, modes[c.mode]))
		}
		w.p("{name: %q, context: %v, tokens: []rhino.Token{%s}, children: []childSpec{%s}},", k.name, k.context, joinTokens(k.tokens), strings.Join(ch, ", "))
	}
	w.p("}")
	w.p("")
	w.p("// grammarCategories maps each category to its kinds, for Check.")
	w.p("var grammarCategories = map[string][]string{")
	for _, c := range w.g.categories {
		var ks []string
		for _, k := range c.kinds {
			ks = append(ks, fmt.Sprintf("%q", k))
		}
		w.p("%q: {%s},", c.name, strings.Join(ks, ", "))
	}
	w.p("}")
}
