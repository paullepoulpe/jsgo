// Port of com/google/javascript/rhino/Node.java
// @ closure-compiler v20260415 (subset).
//
// Java's Node is a ~3000-line AST carrier with hundreds of methods
// (parent / sibling pointers, property bag, static factories, many
// traversal helpers). This file carries only the fields + accessors
// jsgo's downstream passes currently consult: token kind, children,
// string/number payload, attached JSDocInfo, 1-indexed line / 0-
// indexed col. API names mirror Java (Token, FirstChild, AddChildToBack,
// JSDocInfo/SetJSDocInfo, Lineno/Charno/SetPos) so a later expansion
// slots in without source-level churn. Bits we skip for now:
// parent/next/prev links, prop bag (JSDOC_INFO_PROP …),
// side-effect flags, Token annotations, clone helpers.

package rhino

type Node struct {
	token    Token
	str      string
	num      float64
	children []*Node
	jsdoc    *JSDocInfo
	lineno   int // 1-indexed, matches Node.getLineno()
	charno   int // 0-indexed, matches Node.getCharno()
}

func New(tok Token) *Node { return &Node{token: tok} }

func NewString(tok Token, s string) *Node { return &Node{token: tok, str: s} }

func NewNumber(v float64) *Node { return &Node{token: NUMBER, num: v} }

func (n *Node) Token() Token  { return n.token }
func (n *Node) String_() string { return n.str }
func (n *Node) Number() float64 { return n.num }
func (n *Node) JSDocInfo() *JSDocInfo { return n.jsdoc }
func (n *Node) SetJSDocInfo(info *JSDocInfo) { n.jsdoc = info }
func (n *Node) Lineno() int { return n.lineno }
func (n *Node) Charno() int { return n.charno }
func (n *Node) SetPos(line, col int) { n.lineno, n.charno = line, col }

func (n *Node) AddChildToBack(c *Node) { n.children = append(n.children, c) }
func (n *Node) FirstChild() *Node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[0]
}
func (n *Node) Children() []*Node { return n.children }
