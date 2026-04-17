// Port of com/google/javascript/jscomp/TypeInference.java @ closure-compiler v20260415
//
// The Java pass does forward dataflow analysis assigning a JSType to every
// expression node. Tier-0 only needs inference for the initializer of a
// single var-decl, which is always a literal — so we implement the one
// case that matters.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type TypeInference struct{}

func NewTypeInference() *TypeInference { return &TypeInference{} }

// InferType mirrors the result of TypeInference for expression nodes.
// Returns Unknown for tokens this scaffold doesn't handle yet.
func (TypeInference) InferType(n *rhino.Node) jstype.JSType {
	switch n.Token() {
	case rhino.STRINGLIT:
		return jstype.String
	case rhino.NUMBER:
		return jstype.Number
	}
	return jstype.Unknown
}
