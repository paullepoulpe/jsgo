package cli

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

type JUnitTestCase = testjunit.JUnitTestCase

func RunJUnit4(t *testing.T, newTest func() JUnitTestCase) {
	testjunit.RunJUnit4(t, newTest)
}

type cliJSErrorSubject struct {
	*testjunit.Subject
	actual *jscomp.JSError
}

func AssertError(err *jscomp.JSError) *cliJSErrorSubject {
	return &cliJSErrorSubject{Subject: testjunit.NewSubject(err), actual: err}
}

func (s *cliJSErrorSubject) check(what string) *testjunit.StandardSubjectBuilder {
	msg := fmt.Sprintf("value of: jSError.%s\njSError was: %v", what, s.actual)
	if s.Message() != "" {
		msg = s.Message() + "\n" + msg
	}
	return testjunit.NewStandardSubjectBuilder(msg)
}

func (s *cliJSErrorSubject) HasType(t *jscomp.DiagnosticType) {
	s.check("getType()").That(s.actual.Type()).IsEqualTo(t)
}

func (s *cliJSErrorSubject) HasMessage(msg string) {
	s.check("description").That(s.actual.Description()).IsEqualTo(msg)
}

type cliNodeSubject struct {
	*testjunit.Subject
	actual     rhino.Node
	serializer func(rhino.Node) string
}

func AssertNode(node rhino.Node) *cliNodeSubject {
	return &cliNodeSubject{Subject: testjunit.NewSubject(node), actual: node}
}

func (s *cliNodeSubject) UsingSerializer(serializer func(rhino.Node) string) *cliNodeSubject {
	s.serializer = serializer
	return s
}

var (
	cliEmptyBlock     = regexp.MustCompile(`\{\n[ \t\n\x0B\f\r]*\}`)
	cliTrailingSpaces = regexp.MustCompile(` +(\n|$)`)
)

func (s *cliNodeSubject) serializeNode(node rhino.Node) string {
	if s.serializer != nil {
		out := s.serializer(node)
		out = cliEmptyBlock.ReplaceAllString(out, "{}")
		out = cliTrailingSpaces.ReplaceAllString(out, "$1")
		return out
	}
	if node.IsNil() {
		return "<Java null>"
	}
	return node.ToStringTree()
}

func cliFindFirstMismatch(actual, expected rhino.Node) (rhino.Node, rhino.Node, bool) {
	if !actual.IsEquivalentToWithFlagsAndSideEffect(expected, false, false, false, false) {
		return actual, expected, true
	}
	ac := actual.GetFirstChild()
	ec := expected.GetFirstChild()
	for !ac.IsNil() {
		if ma, me, ok := cliFindFirstMismatch(ac, ec); ok {
			return ma, me, true
		}
		ac = ac.GetNext()
		ec = ec.GetNext()
	}
	as := actual.GetClosureUnawareShadow()
	es := expected.GetClosureUnawareShadow()
	if !as.IsNil() && !es.IsNil() {
		if ma, me, ok := cliFindFirstMismatch(as, es); ok {
			return ma, me, true
		}
	}
	return rhino.NilNode, rhino.NilNode, false
}

func (s *cliNodeSubject) IsEqualTo(expected rhino.Node) {
	if s.actual.IsNil() {
		s.Failf("expected node not to be null")
		return
	}
	if expected.IsNil() {
		s.Failf("expected comparison target not to be null")
		return
	}
	ma, me, mismatch := cliFindFirstMismatch(s.actual, expected)
	if !mismatch {
		return
	}
	expJS := s.serializeNode(expected)
	actJS := s.serializeNode(s.actual)
	if expJS == actJS {
		s.Failf("Node tree inequality (AST diff):\nexpected:\n%s\nactual:\n%s", me.ToStringTree(), ma.ToStringTree())
		return
	}
	s.Failf("Node tree inequality (JS diff):\nexpected:\n%s\nactual:\n%s", expJS, actJS)
}

type TestExternsBuilder struct {
	includeIterableExterns  bool
	includeArrayExterns     bool
	includeArgumentsExterns bool
	includeAlertExterns     bool
	includeClosureExterns   bool
	extraExterns            []string
}

func NewTestExternsBuilder() *TestExternsBuilder {
	return &TestExternsBuilder{extraExterns: []string{}}
}

func (b *TestExternsBuilder) AddIterable() *TestExternsBuilder {
	b.includeIterableExterns = true
	return b
}

func (b *TestExternsBuilder) AddArray() *TestExternsBuilder {
	b.includeArrayExterns = true
	b.AddIterable()
	return b
}

func (b *TestExternsBuilder) AddArguments() *TestExternsBuilder {
	b.includeArgumentsExterns = true
	b.AddArray()
	b.AddIterable()
	return b
}

func (b *TestExternsBuilder) AddAlert() *TestExternsBuilder {
	b.includeAlertExterns = true
	return b
}

func (b *TestExternsBuilder) AddClosureExterns() *TestExternsBuilder {
	b.includeClosureExterns = true
	return b
}

func (b *TestExternsBuilder) AddExtra(lines ...string) *TestExternsBuilder {
	b.extraExterns = append(b.extraExterns, lines...)
	return b
}

func (b *TestExternsBuilder) Build() string {
	sections := []string{
		`/**
 * @fileoverview
 * @externs
 */
`,
	}
	if b.includeIterableExterns {
		sections = append(sections, cliIterableExterns)
	}
	if b.includeArrayExterns {
		sections = append(sections, cliArrayExterns)
	}
	if b.includeArgumentsExterns {
		sections = append(sections, cliArgumentsExterns)
	}
	if b.includeAlertExterns {
		sections = append(sections, cliAlertExterns)
	}
	if b.includeClosureExterns {
		sections = append(sections, cliClosureExterns)
	}
	sections = append(sections, b.extraExterns...)
	return strings.Join(sections, "\n")
}

func (b *TestExternsBuilder) BuildExternsFile(filePath string) *jscomp.SourceFile {
	return jscomp.SourceFile_FromCode(filePath, b.Build())
}

type JSChunkGraphBuilder struct {
	star   bool
	chunks []string
}

func JSChunkGraphBuilder_ForChain() *JSChunkGraphBuilder {
	return &JSChunkGraphBuilder{star: false}
}

func JSChunkGraphBuilder_ForStar() *JSChunkGraphBuilder {
	return &JSChunkGraphBuilder{star: true}
}

func (b *JSChunkGraphBuilder) AddChunks(inputs []string) *JSChunkGraphBuilder {
	b.chunks = append(b.chunks, inputs...)
	return b
}

func (b *JSChunkGraphBuilder) Build() []*jscomp.JSChunk {
	chunks := make([]*jscomp.JSChunk, len(b.chunks))
	for i, code := range b.chunks {
		c := jscomp.NewJSChunk(fmt.Sprintf("m%d", i))
		c.AddSourceFile(jscomp.SourceFile_FromCode(jutil.LenientFormat("i%s.js", i), code))
		chunks[i] = c
	}
	if b.star {
		for i := 1; i < len(chunks); i++ {
			chunks[i].AddDependency(chunks[0])
		}
	} else {
		for i := 1; i < len(chunks); i++ {
			chunks[i].AddDependency(chunks[i-1])
		}
	}
	return chunks
}

const cliClosureExterns = `/** @const */ var goog = {};
goog.module = function(ns) {};
goog.module.declareLegacyNamespace = function() {};
goog.module.preventModuleExportSealing = function() {};
/** @return {?} */
goog.module.get = function(ns) {};
goog.provide = function(ns) {};
/** @return {?} */
goog.require = function(ns) {};
/** @return {?} */
goog.requireType = function(ns) {};
goog.requireDynamic = function(ns) {};
goog.loadModule = function(ns) {};
/** @return {?} */
goog.forwardDeclare = function(ns) {};
goog.setTestOnly = function() {};
goog.scope = function(fn) {};
goog.declareModuleId = function(ns) {};
/**
 * @param {string} name
 * @param {T} defaultValue
 * @return {T}
 * @template T
 */
goog.define = function(name, defaultValue) {};
goog.exportSymbol = function(publicName, symbol) {};
goog.inherits = function(childCtor, parentCtor) {
  childCtor.superClass_ = parentCtor.prototype;
};
goog.getMsg = function(str) {};
/**
 * @param {T} symbol
 * @return {T|undefined}
 * @template T
 * @noinline
 */
goog.weakUsage = function(symbol) {};
`

const cliIterableExterns = `// Symbol is needed for Symbol.iterator
/**
 * @constructor
 * @param {*=} opt_description
 * @return {symbol}
 * @nosideeffects
 */
function Symbol(opt_description) {}

/** @const {!symbol} */ Symbol.iterator;

/**
 * @record
 * @template TYield
 */
function IIterableResult() {};
/** @type {boolean} */
IIterableResult.prototype.done;
/** @type {TYield} */
IIterableResult.prototype.value;

/**
 * @record
 * @template T, TReturn, TNext
 */
function IteratorLike() {}
/**
 * @param {T=} value
 * @return {!IIterableResult<T>}
 */
IteratorLike.prototype.next;
/**
 * @type {((function(T=): !IIterableResult<T>)|undefined)}
 */
IteratorLike.prototype.return;
/**
 * @type {((function(?=): !IIterableResult<T>)|undefined)}
 */
IteratorLike.prototype.throw;

/**
 * @interface
 * @template T, TReturn, TNext
 */
function Iterable() {}

/**
 * @return {!IteratorLike<T, ?, *>}
 * @suppress {externsValidation}
 */
Iterable.prototype[Symbol.iterator] = function() {};

/**
 * @interface
 * @extends {IteratorLike<T, ?, *>}
 * @extends {Iterable<T, ?, *>}
 * @template T, TReturn, TNext
 */
function IteratorIterable() {}

/**
 * @constructor
 * @abstract
 * @implements {IteratorIterable<T, TReturn, TNext>}
 * @template T, TReturn, TNext
 * @return {?}
 */
function Iterator() {}
/**
 * @param {string|!Iterable<T>|!IteratorLike<T>} iterable
 * @return {!Iterator<T>}
 * @template T
 */
Iterator.from = function(iterable) {};
/**
 * @override
 * @param {?=} value
 * @return {!IIterableResult<T>}
 */
Iterator.prototype.next;
/**
 * @type {((function(T=): !IIterableResult<T>)|undefined)}
 */
Iterator.prototype.return;
/**
 * @type {((function(?=): !IIterableResult<T>)|undefined)}
 */
Iterator.prototype.throw;
/**
 * @override
 * @return {!Iterator<T, TReturn, TNext>}
 * @suppress {externsValidation}
 */
Iterator.prototype[Symbol.iterator] = function() {};
/**
 * @template U
 * @param {function(T, number): U} mapperFn
 * @return {!Iterator<U>}
 */
Iterator.prototype.map;
/**
 * @param {function(T, number): *} filtererFn
 * @return {!Iterator<T>}
 */
Iterator.prototype.filter;
/**
 * @param {number} limit
 * @return {!Iterator<T>}
 */
Iterator.prototype.take;
/**
 * @param {number} limit
 * @return {!Iterator<T>}
 */
Iterator.prototype.drop;
/**
 * @template U
 * @param {function(T, number): (!Iterable<U>|!IteratorLike<U>|!Iterator<U>)} mapperFn
 * @return {!Iterator<U>}
 */
Iterator.prototype.flatMap;
/**
 * @template U
 * @param {function(U, T, number): U} reducerFn
 * @param {U=} initialValue
 * @return {U}
 */
Iterator.prototype.reduce;
/**
 * @return {!Array<T>}
 */
Iterator.prototype.toArray;
/**
 * @param {function(T, number): *} callbackFn
 * @return {void}
 */
Iterator.prototype.forEach;
/**
 * @param {function(T, number): *} predicate
 * @return {boolean}
 */
Iterator.prototype.some;
/**
 * @param {function(T, number): *} predicate
 * @return {boolean}
 */
Iterator.prototype.every;
/**
 * @param {function(T, number): *} predicate
 * @return {T|undefined}
 */
Iterator.prototype.find;

/**
 * @interface
 * @extends {IteratorIterable<T, ?, *>}
 * @template T, TReturn, TNext
 */
function Generator() {}
/**
 * @param {?=} opt_value
 * @return {!IIterableResult<T>}
 * @override
 */
Generator.prototype.next = function(opt_value) {};
/**
 * @param {T=} value
 * @return {!IIterableResult<T>}
 */
Generator.prototype.return = function(value) {};
/**
 * @param {?=} exception
 * @return {!IIterableResult<T>}
 */
Generator.prototype.throw = function(exception) {};
`

const cliArrayExterns = `/**
 * @interface
 * @template KEY1, VALUE1
 */
function IObject() {};

/**
 * @record
 * @extends IObject<number, VALUE2>
 * @template VALUE2
 */
function IArrayLike() {};

/** @type {number} */
IArrayLike.prototype.length;

/**
 * @template T
 * @record
 * @extends {IArrayLike<T>}
 * @extends {Iterable<T>}
 */
function ReadonlyArray(var_args) {}
/** @type {number} */ ReadonlyArray.prototype.length;
/**
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @this {!IArrayLike<T>|string}
 * @template T,S
 * @return {undefined}
 */
ReadonlyArray.prototype.forEach;
/**
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T,S
 */
ReadonlyArray.prototype.filter;
/**
 * @param {...*} var_args
 * @return {!Array<?>}
 * @this {*}
 */
ReadonlyArray.prototype.concat;
/**
 * @param {?number=} begin Zero-based index at which to begin extraction.
 * @param {?number=} end Zero-based index at which to end extraction.  slice
 * extracts up to but not including end.
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
ReadonlyArray.prototype.slice;

/**
 * @return {!IteratorIterable<T>}
 */
ReadonlyArray.prototype.values;

/**
 * @param {T} searchElement
 * @param {number=} fromIndex
 * @return {boolean}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
ReadonlyArray.prototype.includes;
/**
 * @template T
 * @constructor
 * @implements {ReadonlyArray<T>}
 * @implements {IArrayLike<T>}
 * @implements {Iterable<T>}
 * @param {...*} var_args
 * @return {!Array<?>}
 */
function Array(var_args) {}
/** @override */
Array.prototype[Symbol.iterator] = function() {};
/** @override @type {number} */ Array.prototype.length;
/**
 * @param {*} arr
 * @return {boolean}
 */
Array.isArray = function(arr) {};

/**
 * @param {string|!IArrayLike<T>|!Iterable<T>} arrayLike
 * @param {function(this:S, (string|T), number): R=} mapFn
 * @param {S=} thisObj
 * @return {!Array<R>}
 * @template T,S,R
 */
Array.from = function(arrayLike, mapFn, thisObj) {}

/**
 * @param {...T} var_args
 * @return {number} The new length of the array.
 * @this {!IArrayLike<T>}
 * @template T
 * @modifies {this}
 */
Array.prototype.push = function(var_args) {};
/**
 * @this {!IArrayLike<T>}
 * @return {T}
 * @template T
 */
Array.prototype.shift = function() {};
/**
 * @override
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @this {!IArrayLike<T>|string}
 * @template T,S
 * @return {undefined}
 */
Array.prototype.forEach = function(callback, opt_thisobj) {};
/**
 * @override
 * @param {?function(this:S, T, number, !Array<T>): ?} callback
 * @param {S=} opt_thisobj
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T,S
 */
Array.prototype.filter = function(callback, opt_thisobj) {};
/**
 * @override
 * @param {...*} var_args
 * @return {!Array<?>}
 * @this {*}
 */
Array.prototype.concat = function(var_args) {};
/**
 * @override
 * @param {?number=} begin Zero-based index at which to begin extraction.
 * @param {?number=} end Zero-based index at which to end extraction.  slice
 * extracts up to but not including end.
 * @return {!Array<T>}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
Array.prototype.slice = function(begin, end) {};

/**
 * @override
 * @return {!IteratorIterable<T>}
 */
Array.prototype.values;

/**
 * @override
 * @param {T} searchElement
 * @param {number=} fromIndex
 * @return {boolean}
 * @this {!IArrayLike<T>|string}
 * @template T
 * @nosideeffects
 */
Array.prototype.includes = function(searchElement, fromIndex) {};
`

const cliArgumentsExterns = `/**
 * @constructor
 * @implements {IArrayLike<?>}
 * @implements {Iterable<?>}
 */
function Arguments() {}

/** @type {number} */
Arguments.prototype.length;
/** @override */
Arguments.prototype[Symbol.iterator] = function() {};

/** @type {!Arguments} */
var arguments;
`

const cliAlertExterns = `/**
 * @param {*} message
 * @return {undefined}
 */
function alert(message) {}
`
