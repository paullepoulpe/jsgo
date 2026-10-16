// NOT a port. Placeholder for com/google/javascript/jscomp/parsing/TypeTransformationParser.java,
// which is being ported on branch port/irfactory (it depends on ParserRunner and IRFactory).
//
// DELETE THIS FILE when type_transformation_parser.go lands: the real port
// declares the same three symbols, so the merge fails to compile ("redeclared")
// until this file is removed. JsDocInfoParser.parseTtlAst calls the real API;
// until then a `@template T := ... =:` type transformation panics here.

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeTransformationParser is not ported yet (see the file comment).
type TypeTransformationParser struct{}

// Java: TypeTransformationParser(String, StaticSourceFile, ErrorReporter, int, int)
func NewTypeTransformationParser(
	typeTransformationString string,
	sourceFile rhino.StaticSourceFile,
	errorReporter rhino.ErrorReporter,
	templateLineno int,
	templateCharno int,
) *TypeTransformationParser {
	panic(jutil.NewUnsupportedOperationExceptionMsg(pendingTypeTransformationParserMessage))
}

// Java: TypeTransformationParser.parseTypeTransformation()
func (p *TypeTransformationParser) ParseTypeTransformation() bool {
	panic(jutil.NewUnsupportedOperationExceptionMsg(pendingTypeTransformationParserMessage))
}

// Java: TypeTransformationParser.getTypeTransformationAst()
func (p *TypeTransformationParser) GetTypeTransformationAst() *rhino.Node {
	panic(jutil.NewUnsupportedOperationExceptionMsg(pendingTypeTransformationParserMessage))
}

// pendingTypeTransformationParserMessage is the message of the panics above
// (tests that reach a type transformation skip when they see it).
const pendingTypeTransformationParserMessage = "TypeTransformationParser is not ported yet (port/irfactory)"
