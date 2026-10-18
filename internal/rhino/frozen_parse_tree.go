package rhino

import (
	"math/bits"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// FrozenParseTree is an immutable, private parse-template representation. It
// stores detached node payloads, integer tree relations, and JSDoc clone plans.
// Clone allocates fresh contiguous nodes, properties, docs, and extended node
// payloads, then fixes their pointers without interpreting the AST again.
//
// Freeze only private parse templates, never a compiler's pooled live AST.
// The current experiment rejects type/color data, shadows, documentation marker
// positions, and unknown object-valued Node properties; callers retain their
// general clone fallback. Strings, InputIds, feature sets and JSDoc properties
// whose ordinary Clone is identity are shared exactly as the existing clone.
// Mutable type-expression/TTL roots and their containers are always copied.
// No pointer into a returned clone is stored back into this template.
type FrozenParseTree struct {
	nodes     []frozenParseNode
	props     []frozenParseProp
	docs      []frozenParseDoc
	exts      []node_ext
	exprCount int
}

type frozenParseNode struct {
	base                          Node
	parent, previous, next, first int
	propStart, propEnd, ext       int
}
type frozenParseProp struct {
	base Node_PropListItem
	doc  int
}
type frozenParseExpr struct {
	root   int
	source string
}
type frozenParseDocProp struct {
	kind  uint8 // ordinary, type, type-list, type-map, TTL node-map
	value any
	names []string
	exprs []frozenParseExpr
	roots []int
}
type frozenParseDoc struct {
	bits, keys int64
	props      []frozenParseDocProp
}

// Node_FreezeParseTree returns nil when the generic clone must be used.
func Node_FreezeParseTree(root *Node) *FrozenParseTree {
	if root == nil {
		return nil
	}
	f := &FrozenParseTree{}
	if _, ok := f.freezeTree(root); !ok {
		return nil
	}
	return f
}

func (f *FrozenParseTree) freezeTree(n *Node) (int, bool) {
	if n.typeOrColor != nil || n.GetClosureUnawareShadow() != nil {
		return 0, false
	}
	id := len(f.nodes)
	plan := frozenParseNode{base: *n, parent: -1, previous: -1, next: -1, first: -1, ext: -1}
	plan.base.parent, plan.base.previous, plan.base.next, plan.base.first, plan.base.propListHead = nil, nil, nil, nil, nil
	plan.base.ext = nil
	if n.ext != nil {
		plan.ext = len(f.exts)
		f.exts = append(f.exts, *n.ext)
	}
	// CloneNode restores the folded subclass token, including unusual nodes
	// whose token has been modified since construction.
	switch n.kind {
	case node_kindNumber:
		plan.base.token = Token_NUMBER
	case node_kindBigInt:
		plan.base.token = Token_BIGINT
	case node_kindTemplateLiteralSubstring:
		plan.base.token = Token_TEMPLATELIT_STRING
	}
	f.nodes = append(f.nodes, plan)
	var props []frozenParseProp
	if doc := n.GetJSDocInfo(); doc != nil {
		docID, ok := f.freezeDoc(doc)
		if !ok {
			return 0, false
		}
		props = append(props, frozenParseProp{base: Node_PropListItem{propType: int8(Node_Prop_JSDOC_INFO), isObject: true}, doc: docID})
	}
	// This is the property order after SetStaticSourceFile followed by
	// SetJSDocInfo, as performed by the original parse-template bind walk.
	props = append(props, frozenParseProp{base: Node_PropListItem{propType: int8(Node_Prop_SOURCE_FILE), isObject: true}, doc: -1})
	for p := n.propListHead; p != nil; p = p.next {
		kind := Node_Prop(p.propType)
		if kind == Node_Prop_JSDOC_INFO || kind == Node_Prop_SOURCE_FILE {
			continue
		}
		if p.isObject && kind != Node_Prop_INPUT_ID && kind != Node_Prop_FEATURE_SET {
			return 0, false
		}
		copy := *p
		copy.next = nil
		props = append(props, frozenParseProp{base: copy, doc: -1})
	}
	f.nodes[id].propStart = len(f.props)
	f.props = append(f.props, props...)
	f.nodes[id].propEnd = len(f.props)
	first, last := -1, -1
	for child := n.first; child != nil; child = child.next {
		childID, ok := f.freezeTree(child)
		if !ok {
			return 0, false
		}
		f.nodes[childID].parent = id
		if first < 0 {
			first = childID
		} else {
			f.nodes[last].next = childID
			f.nodes[childID].previous = last
		}
		last = childID
	}
	if first >= 0 {
		f.nodes[id].first = first
		f.nodes[first].previous = last
	}
	return id, true
}

func (f *FrozenParseTree) freezeExpr(expr *JSTypeExpression) (frozenParseExpr, bool) {
	if expr == nil {
		return frozenParseExpr{root: -1}, true
	}
	root, ok := f.freezeTree(expr.root)
	if !ok {
		return frozenParseExpr{}, false
	}
	f.exprCount++
	return frozenParseExpr{root: root, source: expr.sourceName}, true
}

func (f *FrozenParseTree) freezeDoc(doc *JSDocInfo) (int, bool) {
	id := len(f.docs)
	f.docs = append(f.docs, frozenParseDoc{})
	plan := frozenParseDoc{bits: doc.propertyBits, keys: doc.propertyKeysBitset}
	remaining := uint64(doc.propertyKeysBitset)
	for i := 0; remaining != 0; i++ {
		bit := bits.TrailingZeros64(remaining)
		remaining &= remaining - 1
		value := doc.GetPropertyValueByIndex(i)
		p := frozenParseDocProp{}
		switch prop := JSDocInfo_Property_values[bit].(type) {
		case *JSDocInfo_PropertyBase:
			// The original property's Clone is identity for immutable metadata.
			p.value = value
		case *JSDocInfo_TypeProperty:
			p.kind = 1
			expr, ok := f.freezeExpr(value.(*JSTypeExpression))
			if !ok {
				return 0, false
			}
			p.exprs = []frozenParseExpr{expr}
		case *JSDocInfo_TypeListProperty:
			p.kind = 2
			for _, value := range value.(*jutil.ArrayList[*JSTypeExpression]).Slice() {
				expr, ok := f.freezeExpr(value)
				if !ok {
					return 0, false
				}
				p.exprs = append(p.exprs, expr)
			}
		case *JSDocInfo_TypeMapProperty:
			p.kind = 3
			for name, value := range value.(*jutil.LinkedHashMap[string, *JSTypeExpression]).All() {
				expr, ok := f.freezeExpr(value)
				if !ok {
					return 0, false
				}
				p.names = append(p.names, name)
				p.exprs = append(p.exprs, expr)
			}
		case *JSDocInfo_NodeMapProperty:
			if prop != JSDocInfo_TYPE_TRANSFORMATIONS {
				return 0, false
			}
			p.kind = 4
			for name, tree := range value.(*jutil.LinkedHashMap[string, *Node]).All() {
				root, ok := f.freezeTree(tree)
				if !ok {
					return 0, false
				}
				p.names = append(p.names, name)
				p.roots = append(p.roots, root)
			}
		default:
			return 0, false
		}
		plan.props = append(plan.props, p)
	}
	f.docs[id] = plan
	return id, true
}

// CloneWithSourceFile returns fresh mutable trees with precisely the source
// rebinding of the existing worker cache, including all JSDoc/TTL nodes.
func (f *FrozenParseTree) CloneWithSourceFile(source StaticSourceFile) *Node {
	nodes := make([]Node, len(f.nodes))
	props := make([]Node_PropListItem, len(f.props))
	docs := make([]JSDocInfo, len(f.docs))
	exts := append([]node_ext(nil), f.exts...)
	exprs := make([]JSTypeExpression, f.exprCount)
	exprNext := 0
	ptr := func(id int) *Node {
		if id < 0 {
			return nil
		}
		return &nodes[id]
	}
	for id, plan := range f.nodes {
		n := &nodes[id]
		*n = plan.base
		n.parent, n.previous, n.next, n.first = ptr(plan.parent), ptr(plan.previous), ptr(plan.next), ptr(plan.first)
		if plan.ext >= 0 {
			n.ext = &exts[plan.ext]
		}
		n.propListHead = &props[plan.propStart]
	}
	cloneExpr := func(plan frozenParseExpr) *JSTypeExpression {
		if plan.root < 0 {
			return nil
		}
		e := &exprs[exprNext]
		exprNext++
		e.root, e.sourceName = &nodes[plan.root], plan.source
		return e
	}
	for id, plan := range f.docs {
		doc := &docs[id]
		doc.propertyBits, doc.propertyKeysBitset = plan.bits, plan.keys
		values := make([]any, len(plan.props))
		for i, p := range plan.props {
			switch p.kind {
			case 0:
				values[i] = p.value
			case 1:
				values[i] = cloneExpr(p.exprs[0])
			case 2:
				list := jutil.NewArrayList[*JSTypeExpression]()
				for _, e := range p.exprs {
					list.Add(cloneExpr(e))
				}
				values[i] = list
			case 3:
				m := jutil.NewLinkedHashMap[string, *JSTypeExpression]()
				for j, e := range p.exprs {
					m.Put(p.names[j], cloneExpr(e))
				}
				values[i] = m
			case 4:
				m := jutil.NewLinkedHashMap[string, *Node]()
				for j, root := range p.roots {
					m.Put(p.names[j], &nodes[root])
				}
				values[i] = m
			}
		}
		if len(values) == 1 {
			doc.propertyValues = values[0]
		} else if len(values) > 1 {
			doc.propertyValues = values
		}
	}
	for _, plan := range f.nodes {
		for id := plan.propStart; id < plan.propEnd; id++ {
			p := &props[id]
			*p = f.props[id].base
			if id+1 < plan.propEnd {
				p.next = &props[id+1]
			}
			if f.props[id].doc >= 0 {
				p.objectValue = &docs[f.props[id].doc]
			} else if Node_Prop(p.propType) == Node_Prop_SOURCE_FILE {
				p.objectValue = source
			}
		}
	}
	return &nodes[0]
}
