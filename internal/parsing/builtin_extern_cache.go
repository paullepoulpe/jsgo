package parsing

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: performance experiment — a bounded immutable parse-template cache
// for a persistent worker. It is disabled for the normal command-line binary.
var builtinExternCacheEnabled atomic.Bool
var builtinExternCache struct {
	sync.Mutex
	groups []builtinExternCacheGroup
}

type builtinExternCacheGroup struct {
	config *Config
	files  map[builtinExternCacheKey]*ParserRunner_ParseResult
}

type builtinExternCacheKey struct{ name, source string }

func ParserRunner_EnableBuiltinExternCache(enabled bool) {
	builtinExternCacheEnabled.Store(enabled)
	if !enabled {
		builtinExternCache.Lock()
		builtinExternCache.groups = nil
		builtinExternCache.Unlock()
	}
}

func builtinExternCacheEligible(file rhino.StaticSourceFile, config *Config) bool {
	return builtinExternCacheEnabled.Load() && file.IsExtern() &&
		!file.IsClosureUnawareCode() && !file.IsTypeScriptSource() &&
		strings.HasPrefix(file.GetName(), "externs.zip//") &&
		config.JsDocParsingMode() == Config_JsDocParsing_TYPES_ONLY
}

func parseBuiltinExternCached(file rhino.StaticSourceFile, source string, config *Config, reporter rhino.ErrorReporter) *ParserRunner_ParseResult {
	key := builtinExternCacheKey{file.GetName(), source}
	builtinExternCache.Lock()
	for _, group := range builtinExternCache.groups {
		if group.config.Equals(config) {
			if template := group.files[key]; template != nil {
				builtinExternCache.Unlock()
				return cloneBuiltinExternResult(template, file, false)
			}
			break
		}
	}
	builtinExternCache.Unlock()
	record := &builtinExternCacheReporter{target: reporter}
	result := parserRunnerParseUncached(file, source, config, record)
	// Diagnostics must still be emitted on every invocation. Inputs with
	// any parser/JSDoc diagnostic use the original path instead of caching.
	if record.reported || result.Ast == nil || len(result.Comments) != 0 ||
		len(result.TopLevelStatementRanges) != 0 || result.SourceMapURL != nil {
		return result
	}
	template := cloneBuiltinExternResult(result, file, true)
	builtinExternCache.Lock()
	defer builtinExternCache.Unlock()
	for _, group := range builtinExternCache.groups {
		if group.config.Equals(config) {
			if len(group.files) < 512 {
				group.files[key] = template
			}
			return result
		}
	}
	if len(builtinExternCache.groups) == 4 {
		builtinExternCache.groups = builtinExternCache.groups[1:]
	}
	builtinExternCache.groups = append(builtinExternCache.groups, builtinExternCacheGroup{
		config: config, files: map[builtinExternCacheKey]*ParserRunner_ParseResult{key: template},
	})
	return result
}

func cloneBuiltinExternResult(template *ParserRunner_ParseResult, file rhino.StaticSourceFile, private bool) *ParserRunner_ParseResult {
	result := *template
	cloneTree := func(n *rhino.Node) *rhino.Node { return n.CloneTreeWithTypeExprs(false) }
	if private {
		slab := &jutil.Slab[rhino.Node]{}
		cloneTree = func(n *rhino.Node) *rhino.Node { return rhino.Node_CloneTreeInSlab(n, slab) }
	}
	result.Ast = cloneTree(template.Ast)
	var bind func(*rhino.Node)
	bind = func(n *rhino.Node) {
		n.SetStaticSourceFile(file)
		if doc := n.GetJSDocInfo(); doc != nil {
			// Java's ordinary JSDoc clone intentionally drops INLINE_TYPE;
			// a parse-template copy must preserve that parser fact too.
			builder := doc.ToBuilderWithTransform(func(expr *rhino.JSTypeExpression) *rhino.JSTypeExpression {
				return rhino.NewJSTypeExpression(cloneTree(expr.GetRoot()), expr.GetSourceName())
			})
			if transforms := doc.GetTypeTransformations(); !transforms.IsEmpty() {
				copied := jutil.NewLinkedHashMap[string, *rhino.Node]()
				for name, tree := range transforms.All() {
					copied.Put(name, cloneTree(tree))
				}
				builder.SetProp(rhino.JSDocInfo_TYPE_TRANSFORMATIONS, copied)
			}
			if doc.IsInlineType() {
				builder.RecordInlineType()
			}
			doc = builder.Build()
			n.SetJSDocInfo(doc)
			for _, typeNode := range doc.GetTypeNodes() {
				bind(typeNode)
			}
			for tree := range doc.GetTypeTransformations().Values() {
				bind(tree)
			}
		}
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			bind(child)
		}
	}
	bind(result.Ast)
	return &result
}

type builtinExternCacheReporter struct {
	target   rhino.ErrorReporter
	reported bool
}

func (r *builtinExternCacheReporter) Warning(message, source string, line, offset int) {
	r.reported = true
	r.target.Warning(message, source, line, offset)
}
func (r *builtinExternCacheReporter) Error(message, source string, line, offset int) {
	r.reported = true
	r.target.Error(message, source, line, offset)
}
