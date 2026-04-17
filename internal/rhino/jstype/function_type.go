// Port of com/google/javascript/rhino/jstype/FunctionType.java @ closure-compiler v20260415
//
// FunctionType in Java is a rich node: it carries parameters, return
// type, this-type, supers, template info, constructor/interface flag,
// and more. Tier-1 only exercises the ordered parameter list + return
// type, so that's all we record here; the Java class's extra fields
// will be added as later tiers require them.

package jstype

import "strings"

type FunctionType struct {
	// Ordered positional parameters. Unknown params (e.g. JSDoc left a
	// name unannotated) carry the Unknown JSType so call-site checking
	// stays conservative.
	Params []JSType
	Return JSType // Unknown when @return was absent / unrecognised.
}

func NewFunctionType(params []JSType, ret JSType) *FunctionType {
	if ret == nil {
		ret = Unknown
	}
	return &FunctionType{Params: params, Return: ret}
}

func (f *FunctionType) TypeName() string {
	var names []string
	for _, p := range f.Params {
		names = append(names, p.TypeName())
	}
	return "function(" + strings.Join(names, ", ") + "): " + f.Return.TypeName()
}

// IsSubtypeOf is deliberately permissive for Tier-1: two FunctionType
// values are equivalent only when every param and the return line up
// exactly, and otherwise we accept AllType/Unknown. A faithful variance
// implementation lands with the Tier-2 subtype work.
func (f *FunctionType) IsSubtypeOf(o JSType) bool {
	if g, ok := o.(*FunctionType); ok {
		if len(f.Params) != len(g.Params) {
			return false
		}
		for i := range f.Params {
			if !IsSubtype(f.Params[i], g.Params[i]) {
				return false
			}
		}
		return IsSubtype(f.Return, g.Return)
	}
	return isAll(o) || isUnknown(o)
}
