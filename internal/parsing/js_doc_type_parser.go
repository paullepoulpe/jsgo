// Port of com/google/javascript/jscomp/parsing/JsDocInfoParser.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parsing

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsDocInfoParser_ExtendedTypeInfo is a tuple for recording extended types
type JsDocInfoParser_ExtendedTypeInfo struct {
	type_  *rhino.JSTypeExpression
	lineno int
	charno int
}

// Java: JsDocInfoParser.ExtendedTypeInfo(JSTypeExpression, int, int)
func NewJsDocInfoParser_ExtendedTypeInfo(type_ *rhino.JSTypeExpression, lineno int, charno int) *JsDocInfoParser_ExtendedTypeInfo {
	return &JsDocInfoParser_ExtendedTypeInfo{type_: type_, lineno: lineno, charno: charno}
}

// Based on ES4 grammar proposed on July 10, 2008.
// http://wiki.ecmascript.org/doku.php?id=spec:spec
// Deliberately written to line up with the actual grammar rules,
// for maximum flexibility.

// TODO(nicksantos): The current implementation tries to maintain backwards
// compatibility with previous versions of the spec whenever we can.
// We should try to gradually withdraw support for these.

// ParseTypeExpressionAnnotation parses
// TypeExpressionAnnotation := TypeExpression | '{' TopLevelTypeExpression '}'
//
// Java: JsDocInfoParser.parseTypeExpressionAnnotation(JsDocToken)
func (p *JsDocInfoParser) ParseTypeExpressionAnnotation(token JsDocToken) rhino.Node {
	if token == JsDocToken_LEFT_CURLY {
		p.SkipEOLs()
		typeNode := p.ParseTopLevelTypeExpression(p.Next())
		if !typeNode.IsNil() {
			p.SkipEOLs()
			if !p.Match(JsDocToken_RIGHT_CURLY) {
				if typeNode.IsStringLit() && "import" == typeNode.GetString() {
					p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_IMPORT)
				} else {
					p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RC)
				}
			} else {
				p.Next()
			}
		}

		return typeNode
	} else {
		// TODO(tbreisacher): Add a SuggestedFix for this warning.
		p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_BRACES)
		return p.ParseTypeExpression(token)
	}
}

// ParseParamTypeExpression parses a ParamTypeExpression:
//
//	ParamTypeExpression :=
//	    OptionalParameterType |
//	    TopLevelTypeExpression |
//	    '...' TopLevelTypeExpression
//
//	OptionalParameterType :=
//	    TopLevelTypeExpression '='
//
// Java: JsDocInfoParser.parseParamTypeExpression(JsDocToken)
func (p *JsDocInfoParser) ParseParamTypeExpression(token JsDocToken) rhino.Node {
	restArg := false
	if token == JsDocToken_ITER_REST {
		token = p.Next()
		if token == JsDocToken_RIGHT_CURLY {
			p.RestoreLookAhead(token)
			// EMPTY represents the UNKNOWN type in the Type AST.
			return p.WrapNode(rhino.Token_ITER_REST, p.NewNode(rhino.Token_EMPTY))
		}
		restArg = true
	}

	typeNode := p.ParseTopLevelTypeExpression(token)
	if !typeNode.IsNil() {
		p.SkipEOLs()
		if restArg {
			typeNode = p.WrapNode(rhino.Token_ITER_REST, typeNode)
		} else if p.Match(JsDocToken_EQUALS) {
			p.Next()
			p.SkipEOLs()
			typeNode = p.WrapNode(rhino.Token_EQUALS, typeNode)
		}
	}

	return typeNode
}

// ParseParamTypeExpressionAnnotation parses
// ParamTypeExpressionAnnotation := '{' ParamTypeExpression '}'
//
// Java: JsDocInfoParser.parseParamTypeExpressionAnnotation(JsDocToken)
func (p *JsDocInfoParser) ParseParamTypeExpressionAnnotation(token JsDocToken) rhino.Node {
	jutil.CheckArgument(token == JsDocToken_LEFT_CURLY)

	p.SkipEOLs()

	typeNode := p.ParseParamTypeExpression(p.Next())
	if !typeNode.IsNil() {
		if !p.Match(JsDocToken_RIGHT_CURLY) {
			p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RC)
		} else {
			p.Next()
		}
	}

	return typeNode
}

// ParseTypeNameAnnotation parses TypeNameAnnotation := TypeName | '{' TypeName '}'
//
// Java: JsDocInfoParser.parseTypeNameAnnotation(JsDocToken)
func (p *JsDocInfoParser) ParseTypeNameAnnotation(token JsDocToken) rhino.Node {
	if token == JsDocToken_LEFT_CURLY {
		p.SkipEOLs()
		typeNode := p.ParseTypeName(p.Next())
		if !typeNode.IsNil() {
			p.SkipEOLs()
			if !p.Match(JsDocToken_RIGHT_CURLY) {
				p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RC)
			} else {
				p.Next()
			}
		}

		return typeNode
	} else {
		return p.ParseTypeName(token)
	}
}

// ParseTopLevelTypeExpression parses
// TopLevelTypeExpression := TypeExpression | TypeUnionList
//
// We made this rule up, for the sake of backwards compatibility.
//
// Java: JsDocInfoParser.parseTopLevelTypeExpression(JsDocToken)
func (p *JsDocInfoParser) ParseTopLevelTypeExpression(token JsDocToken) rhino.Node {
	typeExpr := p.ParseTypeExpression(token)
	if !typeExpr.IsNil() {
		// top-level unions are allowed
		if p.Match(JsDocToken_PIPE) {
			p.Next()
			p.SkipEOLs()
			token = p.Next()
			return p.ParseUnionTypeWithAlternate(token, typeExpr)
		}
	}
	return typeExpr
}

// ParseTypeExpressionList parses
// TypeExpressionList := TopLevelTypeExpression | TopLevelTypeExpression ',' TypeExpressionList
//
// Java: JsDocInfoParser.parseTypeExpressionList(String, JsDocToken)
func (p *JsDocInfoParser) ParseTypeExpressionList(typeName string, token JsDocToken) rhino.Node {
	typeExpr := p.ParseTopLevelTypeExpression(token)
	if typeExpr.IsNil() {
		return rhino.NilNode
	}
	typeList := p.NewNode(rhino.Token_BLOCK)
	numTypeExprs := 1
	typeList.AddChildToBack(typeExpr)
	for p.Match(JsDocToken_COMMA) {
		p.Next()
		p.SkipEOLs()
		typeExpr = p.ParseTopLevelTypeExpression(p.Next())
		if typeExpr.IsNil() {
			return rhino.NilNode
		}
		numTypeExprs++
		typeList.AddChildToBack(typeExpr)
	}
	if typeName == "Object" && numTypeExprs == 1 {
		// Unlike other generic types, Object<V> means Object<?, V>, not Object<V, ?>.
		typeList.AddChildToFront(p.NewNode(rhino.Token_QMARK))
	}
	return typeList
}

// ParseTypeExpression parses
// TypeExpression := BasicTypeExpression | '?' BasicTypeExpression | '!' BasicTypeExpression |
// BasicTypeExpression '?' | BasicTypeExpression '!' | '?'
//
// Java: JsDocInfoParser.parseTypeExpression(JsDocToken)
func (p *JsDocInfoParser) ParseTypeExpression(token JsDocToken) rhino.Node {
	// Save the source position before we consume additional tokens.
	lineno := p.stream.GetLineno()
	charno := p.stream.GetCharno()
	switch token {
	case JsDocToken_QMARK:
		// A QMARK could mean that a type is nullable, or that it's unknown.
		// We use look-ahead 1 to determine whether it's unknown. Otherwise,
		// we assume it means nullable. There are 8 cases:
		// {?} - right curly
		// ? - EOF (possible when the parseTypeString method is given a bare type expression)
		// {?=} - equals
		// {function(?, number)} - comma
		// {function(number, ?)} - right paren
		// {function(): ?|number} - pipe
		// {Array.<?>} - greater than
		// /** ? */ - EOC (inline types)
		// I'm not a big fan of using look-ahead for this, but it makes
		// the type language a lot nicer.
		token = p.Next()
		if token == JsDocToken_COMMA ||
			token == JsDocToken_EQUALS ||
			token == JsDocToken_RIGHT_SQUARE ||
			token == JsDocToken_RIGHT_CURLY ||
			token == JsDocToken_RIGHT_PAREN ||
			token == JsDocToken_PIPE ||
			token == JsDocToken_RIGHT_ANGLE ||
			token == JsDocToken_EOC ||
			token == JsDocToken_EOL ||
			token == JsDocToken_EOF {
			p.RestoreLookAhead(token)
			return p.NewNode(rhino.Token_QMARK)
		}

		return p.WrapNodeAt(rhino.Token_QMARK, p.ParseBasicTypeExpression(token), lineno, charno)
	case JsDocToken_BANG:
		return p.WrapNodeAt(rhino.Token_BANG, p.ParseBasicTypeExpression(p.Next()), lineno, charno)
	default:
		basicTypeExpr := p.ParseBasicTypeExpression(token)
		lineno = p.stream.GetLineno()
		charno = p.stream.GetCharno()
		if !basicTypeExpr.IsNil() {
			if p.Match(JsDocToken_QMARK) {
				p.Next()
				return p.WrapNodeAt(rhino.Token_QMARK, basicTypeExpr, lineno, charno)
			} else if p.Match(JsDocToken_BANG) {
				p.Next()
				return p.WrapNodeAt(rhino.Token_BANG, basicTypeExpr, lineno, charno)
			}
		}

		return basicTypeExpr
	}
}

// ParseContextTypeExpression parses
// ContextTypeExpression := BasicTypeExpression | '?' For expressions on the
// right hand side of a this: or new:
//
// Java: JsDocInfoParser.parseContextTypeExpression(JsDocToken)
func (p *JsDocInfoParser) ParseContextTypeExpression(token JsDocToken) rhino.Node {
	if token == JsDocToken_QMARK {
		return p.NewNode(rhino.Token_QMARK)
	} else {
		return p.ParseBasicTypeExpression(token)
	}
}

// ParseBasicTypeExpression parses
// BasicTypeExpression := '*' | 'null' | 'undefined' | TypeName | FunctionType | UnionType |
// RecordType | TypeofType
//
// Java: JsDocInfoParser.parseBasicTypeExpression(JsDocToken)
func (p *JsDocInfoParser) ParseBasicTypeExpression(token JsDocToken) rhino.Node {
	if token == JsDocToken_STAR {
		return p.NewNode(rhino.Token_STAR)
	} else if token == JsDocToken_LEFT_CURLY {
		p.SkipEOLs()
		return p.ParseRecordType(p.Next())
	} else if token == JsDocToken_LEFT_PAREN {
		p.SkipEOLs()
		return p.ParseUnionType(p.Next())
	} else if token == JsDocToken_STRING {
		string_ := p.stream.GetString()
		switch string_ {
		case "function":
			p.SkipEOLs()
			return p.ParseFunctionType(p.Next())
		case "null", "undefined":
			return p.NewStringNode(string_)
		case "typeof":
			p.SkipEOLs()
			return p.ParseTypeofType(p.Next())
		default:
			return p.ParseTypeName(token)
		}
	}

	p.RestoreLookAhead(token)
	return p.ReportGenericTypeSyntaxWarning()
}

// Java: JsDocInfoParser.parseNameExpression(JsDocToken)
func (p *JsDocInfoParser) ParseNameExpression(token JsDocToken) rhino.Node {
	if token != JsDocToken_STRING {
		p.AddParserWarningAt(rhino.Msg_JSDOC_NAME_SYNTAX, p.stream.GetLineno(), p.stream.GetCharno())
		return rhino.NilNode
	}

	startLineno := p.stream.GetLineno()
	startCharno := p.stream.GetCharno()
	startOffset := p.stream.GetCursor() - jutil.StringLength(p.stream.GetString())

	name := p.stream.GetString()
	endOffset := p.stream.GetCursor()
	if p.Match(JsDocToken_EOL) && strings.HasSuffix(name, ".") {
		var typeName strings.Builder
		typeName.WriteString(name)
		for p.Match(JsDocToken_EOL) && strings.HasSuffix(typeName.String(), ".") {
			p.SkipEOLs()
			if p.Match(JsDocToken_STRING) {
				p.Next()
				endOffset = p.stream.GetCursor()
				typeName.WriteString(p.stream.GetString())
			}
		}
		name = typeName.String()
	}

	// What we're doing by concatenating these tokens is really hacky. We want that to be obvious.
	str := p.NewStringNode(name)
	str.SetLinenoCharno(startLineno, startCharno)
	str.SetLength(endOffset - startOffset)
	return str
}

// ParseTypeName parses a TypeName:
//
//	TypeName := NameExpression | NameExpression TypeApplication
//	TypeApplication := '.'? '<' TypeExpressionList '>'
//
// Java: JsDocInfoParser.parseTypeName(JsDocToken)
func (p *JsDocInfoParser) ParseTypeName(token JsDocToken) rhino.Node {
	typeNameNode := p.ParseNameExpression(token)

	if p.Match(JsDocToken_LEFT_ANGLE) {
		p.Next()
		p.SkipEOLs()
		memberType := p.ParseTypeExpressionList(typeNameNode.GetString(), p.Next())
		if !memberType.IsNil() {
			typeNameNode.AddChildToFront(memberType)

			p.SkipEOLs()
			if !p.Match(JsDocToken_RIGHT_ANGLE) {
				return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_GT)
			}

			p.Next()
		}
	}
	return typeNameNode
}

// ParseTypeofType parses
// TypeofType := 'typeof' NameExpression | 'typeof' '(' NameExpression ')'
//
// Java: JsDocInfoParser.parseTypeofType(JsDocToken)
func (p *JsDocInfoParser) ParseTypeofType(token JsDocToken) rhino.Node {
	if token == JsDocToken_LEFT_CURLY {
		return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_UNNECESSARY_BRACES)
	}

	typeofType := p.NewNode(rhino.Token_TYPEOF)
	name := p.ParseNameExpression(token)
	if name.IsNil() {
		return rhino.NilNode
	}

	p.SkipEOLs()
	typeofType.AddChildToFront(name)
	return typeofType
}

// ParseFunctionType parses a FunctionType:
//
//	FunctionType := 'function' FunctionSignatureType
//	FunctionSignatureType :=
//	   TypeParameters '(' 'this' ':' TypeName, ParametersType ')' ResultType
//
// The Node that is produced has type Token.FUNCTION but does not look like a
// typical function node. If there is a 'this:' or 'new:' type, that type is
// added as a child. Then, if there are parameters, a PARAM_LIST node is added
// as a child. Finally, if there is a return type, it is added as a child. This
// means that the parameters could be the first or second child, and the
// return type could be the first, second, or third child.
//
// Java: JsDocInfoParser.parseFunctionType(JsDocToken)
func (p *JsDocInfoParser) ParseFunctionType(token JsDocToken) rhino.Node {
	// NOTE(nicksantos): We're not implementing generics at the moment, so
	// just throw out TypeParameters.
	if token != JsDocToken_LEFT_PAREN {
		p.RestoreLookAhead(token)
		return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_LP)
	}

	functionType := p.NewNode(rhino.Token_FUNCTION)
	var parameters rhino.Node
	p.SkipEOLs()
	if !p.Match(JsDocToken_RIGHT_PAREN) {
		token = p.Next()

		hasParams := true
		if token == JsDocToken_STRING {
			tokenStr := p.stream.GetString()
			isThis := "this" == tokenStr
			isNew := "new" == tokenStr
			if isThis || isNew {
				if p.Match(JsDocToken_COLON) {
					p.Next()
					p.SkipEOLs()
					contextTypeToken := rhino.Token_NEW
					if isThis {
						contextTypeToken = rhino.Token_THIS
					}
					contextType := p.WrapNode(contextTypeToken, p.ParseContextTypeExpression(p.Next()))
					if contextType.IsNil() {
						return rhino.NilNode
					}

					functionType.AddChildToFront(contextType)
				} else {
					return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_COLON)
				}

				if p.Match(JsDocToken_COMMA) {
					p.Next()
					p.SkipEOLs()
					token = p.Next()
				} else {
					hasParams = false
				}
			}
		}

		if hasParams {
			parameters = p.ParseParametersType(token)
			if parameters.IsNil() {
				return rhino.NilNode
			}
		}
	}

	if !parameters.IsNil() {
		functionType.AddChildToBack(parameters)
	}

	p.SkipEOLs()
	if !p.Match(JsDocToken_RIGHT_PAREN) {
		return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RP)
	}

	p.SkipEOLs()
	p.Next()
	resultType := p.ParseResultType()
	if resultType.IsNil() {
		return rhino.NilNode
	} else {
		functionType.AddChildToBack(resultType)
	}
	return functionType
}

// ParseParametersType parses a ParametersType:
//
//	ParametersType := RestParameterType | NonRestParametersType
//	    | NonRestParametersType ',' RestParameterType
//	RestParameterType := '...' Identifier
//	NonRestParametersType := ParameterType ',' NonRestParametersType
//	    | ParameterType
//	    | OptionalParametersType
//	OptionalParametersType := OptionalParameterType
//	    | OptionalParameterType, OptionalParametersType
//	OptionalParameterType := ParameterType=
//	ParameterType := TypeExpression | Identifier ':' TypeExpression
//
// NOTE(nicksantos): The official ES4 grammar forces optional and rest
// arguments to come after the required arguments. Our parser does not
// enforce this. Instead we allow them anywhere in the function at parse-time,
// and then warn about them during type resolution.
//
// In theory, it might be mathematically nicer to do the order-checking here.
// But in practice, the order-checking for structural functions is exactly
// the same as the order-checking for @param annotations. And the latter
// has to happen during type resolution. Rather than duplicate the
// order-checking in two places, we just do all of it in type resolution.
//
// Java: JsDocInfoParser.parseParametersType(JsDocToken)
func (p *JsDocInfoParser) ParseParametersType(token JsDocToken) rhino.Node {
	paramsType := p.NewNode(rhino.Token_PARAM_LIST)
	isVarArgs := false
	var paramType rhino.Node
	if token != JsDocToken_RIGHT_PAREN {
		for {
			if !paramType.IsNil() {
				// skip past the comma
				p.Next()
				p.SkipEOLs()
				token = p.Next()
			}

			if token == JsDocToken_ITER_REST {
				// In the latest ES4 proposal, there are no type constraints allowed
				// on variable arguments. We support the old syntax for backwards
				// compatibility, but we should gradually tear it out.
				p.SkipEOLs()
				if p.Match(JsDocToken_RIGHT_PAREN) {
					paramType = p.NewNode(rhino.Token_ITER_REST)
				} else {
					p.SkipEOLs()
					paramType = p.WrapNode(rhino.Token_ITER_REST, p.ParseTypeExpression(p.Next()))
					p.SkipEOLs()
				}

				isVarArgs = true
			} else {
				paramType = p.ParseTypeExpression(token)
				if p.Match(JsDocToken_EQUALS) {
					p.SkipEOLs()
					p.Next()
					paramType = p.WrapNode(rhino.Token_EQUALS, paramType)
				}
			}

			if paramType.IsNil() {
				return rhino.NilNode
			}
			paramsType.AddChildToBack(paramType)
			if isVarArgs {
				break
			}
			if !p.Match(JsDocToken_COMMA) {
				break
			}
		}
	}

	if isVarArgs && p.Match(JsDocToken_COMMA) {
		return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_FUNCTION_VARARGS)
	}

	// The right paren will be checked by parseFunctionType

	return paramsType
}

// ParseResultType parses ResultType := <empty> | ':' void | ':' TypeExpression
//
// Java: JsDocInfoParser.parseResultType()
func (p *JsDocInfoParser) ParseResultType() rhino.Node {
	p.SkipEOLs()
	if !p.Match(JsDocToken_COLON) {
		return p.NewNode(rhino.Token_EMPTY)
	}

	p.Next()
	p.SkipEOLs()
	if p.Match(JsDocToken_STRING) && "void" == p.stream.GetString() {
		p.Next()
		return p.NewNode(rhino.Token_VOID)
	} else {
		return p.ParseTypeExpression(p.Next())
	}
}

// ParseUnionType parses
// UnionType := '(' TypeUnionList ')' TypeUnionList := TypeExpression | TypeExpression '|'
// TypeUnionList
//
// We've removed the empty union type.
//
// Java: JsDocInfoParser.parseUnionType(JsDocToken)
func (p *JsDocInfoParser) ParseUnionType(token JsDocToken) rhino.Node {
	return p.ParseUnionTypeWithAlternate(token, rhino.NilNode)
}

// ParseUnionTypeWithAlternate creates a new union type, with an alternate
// that has already been parsed. The alternate may be nil.
//
// Java: JsDocInfoParser.parseUnionTypeWithAlternate(JsDocToken, Node)
func (p *JsDocInfoParser) ParseUnionTypeWithAlternate(token JsDocToken, alternate rhino.Node) rhino.Node {
	union := p.NewNode(rhino.Token_PIPE)
	if !alternate.IsNil() {
		union.AddChildToBack(alternate)
	}

	var expr rhino.Node
	for {
		if !expr.IsNil() {
			p.SkipEOLs()
			token = p.Next()
			jutil.CheckState(token == JsDocToken_PIPE)

			p.SkipEOLs()
			token = p.Next()
		}
		expr = p.ParseTypeExpression(token)
		if expr.IsNil() {
			return rhino.NilNode
		}

		union.AddChildToBack(expr)
		if !p.Match(JsDocToken_PIPE) {
			break
		}
	}

	if alternate.IsNil() {
		p.SkipEOLs()
		if !p.Match(JsDocToken_RIGHT_PAREN) {
			return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RP)
		}
		p.Next()
	}
	if union.HasOneChild() {
		firstChild := union.GetFirstChild()
		firstChild.Detach()
		return firstChild
	}
	return union
}

// ParseRecordType parses RecordType := '{' FieldTypeList '}'
//
// Java: JsDocInfoParser.parseRecordType(JsDocToken)
func (p *JsDocInfoParser) ParseRecordType(token JsDocToken) rhino.Node {
	recordType := p.NewNode(rhino.Token_LC)
	fieldTypeList := p.ParseFieldTypeList(token)

	if fieldTypeList.IsNil() {
		return p.ReportGenericTypeSyntaxWarning()
	}

	p.SkipEOLs()
	if !p.Match(JsDocToken_RIGHT_CURLY) {
		return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_MISSING_RC)
	}

	p.Next()

	recordType.AddChildToBack(fieldTypeList)
	return recordType
}

// ParseFieldTypeList parses FieldTypeList := FieldType | FieldType ',' FieldTypeList
//
// Java: JsDocInfoParser.parseFieldTypeList(JsDocToken)
func (p *JsDocInfoParser) ParseFieldTypeList(token JsDocToken) rhino.Node {
	fieldTypeList := p.NewNode(rhino.Token_LB)

	var small [8]string
	nSmall := 0
	var namesMap map[string]struct{}

	for {
		fieldType := p.ParseFieldType(token)

		if fieldType.IsNil() {
			return rhino.NilNode
		}

		var name string
		if fieldType.IsStringKey() {
			name = fieldType.GetString()
		} else {
			name = fieldType.GetFirstChild().GetString()
		}
		var dup bool
		if namesMap != nil {
			_, dup = namesMap[name]
			if !dup {
				namesMap[name] = struct{}{}
			}
		} else {
			for i := 0; i < nSmall; i++ {
				if small[i] == name {
					dup = true
					break
				}
			}
			if !dup {
				if nSmall < len(small) {
					small[nSmall] = name
					nSmall++
				} else {
					namesMap = make(map[string]struct{}, 16)
					for i := 0; i < nSmall; i++ {
						namesMap[small[i]] = struct{}{}
					}
					namesMap[name] = struct{}{}
				}
			}
		}
		if !dup {
			fieldTypeList.AddChildToBack(fieldType)
		} else {
			p.AddTypeWarningWithArg(rhino.Msg_JSDOC_TYPE_RECORD_DUPLICATE, name)
		}

		p.SkipEOLs()
		if !p.Match(JsDocToken_COMMA) {
			break
		}

		// Move to the comma token.
		p.Next()

		// Move to the token past the comma
		p.SkipEOLs()

		if p.Match(JsDocToken_RIGHT_CURLY) {
			// Allow trailing comma (ie, right curly following the comma)
			break
		}

		token = p.Next()
	}

	return fieldTypeList
}

// ParseFieldType parses FieldType := FieldName | FieldName ':' TypeExpression
//
// Java: JsDocInfoParser.parseFieldType(JsDocToken)
func (p *JsDocInfoParser) ParseFieldType(token JsDocToken) rhino.Node {
	fieldName := p.ParseFieldName(token)

	if fieldName.IsNil() {
		return rhino.NilNode
	}

	p.SkipEOLs()
	if !p.Match(JsDocToken_COLON) {
		return fieldName
	}

	// Move to the colon.
	p.Next()

	// Move to the token after the colon and parse
	// the type expression.
	p.SkipEOLs()
	typeExpression := p.ParseTypeExpression(p.Next())

	if typeExpression.IsNil() {
		return rhino.NilNode
	}

	fieldType := p.NewNode(rhino.Token_COLON)
	fieldType.AddChildToBack(fieldName)
	fieldType.AddChildToBack(typeExpression)
	return fieldType
}

// ParseFieldName parses
// FieldName := NameExpression | StringLiteral | NumberLiteral | ReservedIdentifier
//
// Java: JsDocInfoParser.parseFieldName(JsDocToken)
func (p *JsDocInfoParser) ParseFieldName(token JsDocToken) rhino.Node {
	if token != JsDocToken_STRING {
		return rhino.NilNode
	}
	s := p.stream.GetString()
	n := p.alloc.Alloc().InitStringFromTemplate(rhino.Token_STRING_KEY, s, p.templateNode)
	n.SetLinenoCharno(p.stream.GetLineno(), p.stream.GetCharno())
	n.SetLength(jutil.StringLength(s))
	return n
}

// Java: JsDocInfoParser.wrapNode(Token, Node)
func (p *JsDocInfoParser) WrapNode(type_ rhino.Token, n rhino.Node) rhino.Node {
	if n.IsNil() {
		return rhino.NilNode
	}
	return p.WrapNodeAt(type_, n, n.GetLineno(), n.GetCharno())
}

// Java: JsDocInfoParser.wrapNode(Token, Node, int, int)
func (p *JsDocInfoParser) WrapNodeAt(type_ rhino.Token, n rhino.Node, lineno int, charno int) rhino.Node {
	if n.IsNil() {
		return rhino.NilNode
	}
	return p.alloc.Alloc().InitWithChildFromTemplate(type_, n, p.templateNode).SetLinenoCharno(lineno, charno)
}

// Java: JsDocInfoParser.newNode(Token)
func (p *JsDocInfoParser) NewNode(type_ rhino.Token) rhino.Node {
	return p.alloc.Alloc().InitFromTemplate(type_, p.templateNode).
		SetLinenoCharno(p.stream.GetLineno(), p.stream.GetCharno())
}

// Java: JsDocInfoParser.newStringNode(String)
func (p *JsDocInfoParser) NewStringNode(s string) rhino.Node {
	n := p.alloc.Alloc().InitStringFromTemplate(rhino.Token_STRINGLIT, s, p.templateNode)
	n.SetLinenoCharno(p.stream.GetLineno(), p.stream.GetCharno())
	n.SetLength(jutil.StringLength(s))
	return n
}

// Java: JsDocInfoParser.reportTypeSyntaxWarning(Msg)
func (p *JsDocInfoParser) ReportTypeSyntaxWarning(warning rhino.Msg) rhino.Node {
	p.AddTypeWarningAt(warning, p.stream.GetLineno(), p.stream.GetCharno())
	return rhino.NilNode
}

// Java: JsDocInfoParser.reportGenericTypeSyntaxWarning()
func (p *JsDocInfoParser) ReportGenericTypeSyntaxWarning() rhino.Node {
	return p.ReportTypeSyntaxWarning(rhino.Msg_JSDOC_TYPE_SYNTAX)
}
