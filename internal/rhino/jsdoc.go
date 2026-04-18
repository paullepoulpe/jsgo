// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20260415
//
// The Java class is large; this port carries the annotations Tier-0 and
// Tier-1 need: @type, @param, @return. Further tags are added as later
// tiers require them.

package rhino

type JSDocInfo struct {
	typeExpr      *JSTypeExpression
	paramNames    []string                     // insertion order (matches Java)
	params        map[string]*JSTypeExpression // name → type
	returnExpr    *JSTypeExpression
	isConstructor bool
	isInterface   bool
	isRecord      bool
	extendsExpr   *JSTypeExpression
	implementsAll []*JSTypeExpression
}

// IsInterface / IsRecord mirror JSDocInfo.isInterface() /
// JSDocInfo.usesImplicitMatch() — together they mark a class as a
// structural type.
func (i *JSDocInfo) IsInterface() bool   { return i != nil && i.isInterface }
func (i *JSDocInfo) IsRecord() bool      { return i != nil && i.isRecord }
func (i *JSDocInfo) SetInterface(b bool) { i.isInterface = b }
func (i *JSDocInfo) SetRecord(b bool)    { i.isRecord = b }

// Implementations mirrors JSDocInfo.getImplementedInterfaces().
func (i *JSDocInfo) Implementations() []*JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.implementsAll
}
func (i *JSDocInfo) AddImplementation(t *JSTypeExpression) {
	i.implementsAll = append(i.implementsAll, t)
}

// BaseType mirrors JSDocInfo.getBaseType() — the type expression
// captured from `@extends {Name}` on a constructor or interface.
func (i *JSDocInfo) BaseType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.extendsExpr
}

func (i *JSDocInfo) SetBaseType(t *JSTypeExpression) { i.extendsExpr = t }

// IsConstructor mirrors JSDocInfo.isConstructor().
func (i *JSDocInfo) IsConstructor() bool { return i != nil && i.isConstructor }

func (i *JSDocInfo) SetConstructor(b bool) { i.isConstructor = b }

func (i *JSDocInfo) Type() *JSTypeExpression     { return i.typeExpr }
func (i *JSDocInfo) SetType(t *JSTypeExpression) { i.typeExpr = t }
func (i *JSDocInfo) HasType() bool               { return i != nil && i.typeExpr != nil }

// ParameterNames matches JSDocInfo.getParameterNames() — iteration
// order is the order @param tags appeared in the comment.
func (i *JSDocInfo) ParameterNames() []string {
	if i == nil {
		return nil
	}
	return i.paramNames
}

// ParameterType matches JSDocInfo.getParameterType(String).
func (i *JSDocInfo) ParameterType(name string) *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.params[name]
}

func (i *JSDocInfo) AddParam(name string, t *JSTypeExpression) {
	if i.params == nil {
		i.params = map[string]*JSTypeExpression{}
	}
	if _, seen := i.params[name]; !seen {
		i.paramNames = append(i.paramNames, name)
	}
	i.params[name] = t
}

// ReturnType matches JSDocInfo.getReturnType().
func (i *JSDocInfo) ReturnType() *JSTypeExpression {
	if i == nil {
		return nil
	}
	return i.returnExpr
}

func (i *JSDocInfo) SetReturnType(t *JSTypeExpression) { i.returnExpr = t }
