// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20260415
//
// The Java Node is a large class with many properties (for every stage of
// the pipeline). This port tracks only what Tier-0 needs: token, children,
// a string payload (for NAME / STRINGLIT), a numeric payload (for NUMBER),
// attached JSDocInfo, and source position (1-indexed line, 0-indexed col
// to match the Java getLineno / getCharno convention).

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
