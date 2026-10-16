// Port of com/google/javascript/jscomp/parsing/FeatureCollector.java @ closure-compiler v20260415

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FeatureCollector is a class that deserializes an AstNode-tree representing
// a SCRIPT into a Node-tree.
//
// This process depends on other information from the TypedAST format, but
// the output it limited to only a single SCRIPT. The other deserialized
// content must be provided beforehand.
//
// Java: final class FeatureCollector
type FeatureCollector struct {
	scriptFeatures *parser.FeatureSet
}

// NewFeatureCollector is the implicit default constructor.
//
// Java: FeatureCollector()
func NewFeatureCollector() *FeatureCollector {
	c := &FeatureCollector{}
	c.scriptFeatures = parser.FeatureSet_BARE_MINIMUM
	return c
}

// Java: FeatureCollector.allFeatures()
func (c *FeatureCollector) AllFeatures() *parser.FeatureSet {
	return c.scriptFeatures
}

// Java: FeatureCollector.visitScript(Node)
func (c *FeatureCollector) VisitScript(script *rhino.Node) {
	emptyContext := jutil.NewArrayList[FeatureCollector_FeatureContext]()
	emptyContext.Add(FeatureCollector_FeatureContext_NONE)
	c.VisitRecursive(emptyContext, script)
}

// Java: FeatureCollector.visitRecursive(List<FeatureContext>, Node)
func (c *FeatureCollector) VisitRecursive(contextStack *jutil.ArrayList[FeatureCollector_FeatureContext], node *rhino.Node) {
	newContext := c.VisitSingleNode(contextStack, node)
	contextStack.Add(newContext)
	for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
		c.VisitRecursive(contextStack, child)
	}
	contextStack.RemoveAt(contextStack.Size() - 1)
}

// VisitSingleNode records features for a single node.
//
// contextStack: Ancestor node contexts. Last is the parent of node.
// node: The node to visit.
//
// Java: FeatureCollector.visitSingleNode(List<FeatureContext>, Node)
func (c *FeatureCollector) VisitSingleNode(contextStack *jutil.ArrayList[FeatureCollector_FeatureContext], node *rhino.Node) FeatureCollector_FeatureContext {
	jutil.CheckStateMsg(
		contextStack.Size() >= 1,
		"The provided feature context stack was empty. It must have at least one element.")

	c.RecordScriptFeatures(contextStack, node)
	return FeatureCollector_ContextFor(contextStack, node)
}

// Java: FeatureCollector.recordScriptFeatures(List<FeatureContext>, Node)
func (c *FeatureCollector) RecordScriptFeatures(contextStack *jutil.ArrayList[FeatureCollector_FeatureContext], node *rhino.Node) {
	parentContext := contextStack.GetLast()
	switch node.GetToken() {
	case rhino.Token_FUNCTION:
		if node.IsAsyncGeneratorFunction() {
			c.AddScriptFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS)
		}
		if node.IsArrowFunction() {
			c.AddScriptFeature(parser.FeatureSet_Feature_ARROW_FUNCTIONS)
		}
		if node.IsAsyncFunction() {
			c.AddScriptFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
		}
		if node.IsGeneratorFunction() {
			c.AddScriptFeature(parser.FeatureSet_Feature_GENERATORS)
		}

		if parentContext == FeatureCollector_FeatureContext_BLOCK_SCOPE {
			c.scriptFeatures = c.scriptFeatures.With(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)
		}
	case rhino.Token_PARAM_LIST, rhino.Token_CALL, rhino.Token_NEW:
	case rhino.Token_STRING_KEY:
		if node.IsShorthandProperty() && parentContext == FeatureCollector_FeatureContext_OBJECT_LITERAL {
			c.AddScriptFeature(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
		}
	case rhino.Token_DEFAULT_VALUE:
		if parentContext == FeatureCollector_FeatureContext_PARAM_LIST {
			c.AddScriptFeature(parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
		}
	case rhino.Token_GETTER_DEF:
		c.AddScriptFeature(parser.FeatureSet_Feature_GETTER)
		if parentContext == FeatureCollector_FeatureContext_CLASS_MEMBERS {
			c.AddScriptFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
		}
	case rhino.Token_REGEXP:
		c.AddScriptFeature(parser.FeatureSet_Feature_REGEXP_SYNTAX)
	case rhino.Token_SETTER_DEF:
		c.AddScriptFeature(parser.FeatureSet_Feature_SETTER)
		if parentContext == FeatureCollector_FeatureContext_CLASS_MEMBERS {
			c.AddScriptFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
		}
	case rhino.Token_BLOCK:
		if parentContext == FeatureCollector_FeatureContext_CLASS_MEMBERS {
			c.AddScriptFeature(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK)
		}
	case rhino.Token_EMPTY:
		if parentContext == FeatureCollector_FeatureContext_CATCH {
			c.AddScriptFeature(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
		}
	case rhino.Token_ITER_REST:
		if parentContext == FeatureCollector_FeatureContext_PARAM_LIST {
			c.AddScriptFeature(parser.FeatureSet_Feature_REST_PARAMETERS)
		} else {
			c.AddScriptFeature(parser.FeatureSet_Feature_ARRAY_PATTERN_REST)
		}
	case rhino.Token_ITER_SPREAD:
		c.AddScriptFeature(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)
	case rhino.Token_OBJECT_REST:
		c.AddScriptFeature(parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
	case rhino.Token_OBJECT_SPREAD:
		c.AddScriptFeature(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
	case rhino.Token_BIGINT:
		c.AddScriptFeature(parser.FeatureSet_Feature_BIGINT)
	case rhino.Token_EXPONENT, rhino.Token_ASSIGN_EXPONENT:
		c.AddScriptFeature(parser.FeatureSet_Feature_EXPONENT_OP)
	case rhino.Token_TAGGED_TEMPLATELIT, rhino.Token_TEMPLATELIT:
		c.AddScriptFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	case rhino.Token_NEW_TARGET:
		c.AddScriptFeature(parser.FeatureSet_Feature_NEW_TARGET)
	case rhino.Token_COMPUTED_PROP:
		c.AddScriptFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
		isGetter := node.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER)
		isSetter := node.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER)
		isClassMember := parentContext == FeatureCollector_FeatureContext_CLASS_MEMBERS
		if isGetter {
			c.AddScriptFeature(parser.FeatureSet_Feature_GETTER)
		} else if isSetter {
			c.AddScriptFeature(parser.FeatureSet_Feature_SETTER)
		}
		if (isGetter || isSetter) && isClassMember {
			c.AddScriptFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
		}
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM:
		c.AddScriptFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	case rhino.Token_COALESCE:
		c.AddScriptFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP)
	case rhino.Token_DYNAMIC_IMPORT:
		c.AddScriptFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND:
		c.AddScriptFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	case rhino.Token_ASSIGN_COALESCE:
		c.AddScriptFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP)
		c.AddScriptFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	case rhino.Token_FOR_OF:
		c.AddScriptFeature(parser.FeatureSet_Feature_FOR_OF)
	case rhino.Token_FOR_AWAIT_OF:
		c.AddScriptFeature(parser.FeatureSet_Feature_FOR_AWAIT_OF)
		c.MaybeAddTopLevelAwait(contextStack)
	case rhino.Token_IMPORT, rhino.Token_EXPORT:
		c.AddScriptFeature(parser.FeatureSet_Feature_MODULES)
	case rhino.Token_CONST:
		c.AddScriptFeature(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	case rhino.Token_LET:
		c.AddScriptFeature(parser.FeatureSet_Feature_LET_DECLARATIONS)
	case rhino.Token_CLASS:
		c.AddScriptFeature(parser.FeatureSet_Feature_CLASSES)
	case rhino.Token_MEMBER_FUNCTION_DEF:
		c.AddScriptFeature(parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		c.AddScriptFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	case rhino.Token_SUPER:
		c.AddScriptFeature(parser.FeatureSet_Feature_SUPER)
	case rhino.Token_ARRAY_PATTERN:
		c.AddScriptFeature(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	case rhino.Token_OBJECT_PATTERN:
		c.AddScriptFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	// We need to track this separately from FUNCTION in the case of
	// top-level await.
	case rhino.Token_AWAIT:
		c.AddScriptFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
		c.MaybeAddTopLevelAwait(contextStack)
	default:
	}
}

// Java: FeatureCollector.maybeAddTopLevelAwait(List<FeatureContext>)
func (c *FeatureCollector) MaybeAddTopLevelAwait(contextStack *jutil.ArrayList[FeatureCollector_FeatureContext]) {
	if !contextStack.Contains(FeatureCollector_FeatureContext_FUNCTION) {
		c.AddScriptFeature(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	}
}

// Java: FeatureCollector.addScriptFeature(Feature)
func (c *FeatureCollector) AddScriptFeature(x parser.FeatureSet_Feature) {
	c.scriptFeatures = c.scriptFeatures.With(x)
}

// FeatureCollector_FeatureContext is the context of a node while
// deserializing, specifically for the purpose of tracking FeatureSet.Feature.
//
// These contexts are pushed onto a stack as the AST is traversed. The intent
// is to have a way to model the ancestry of a node being visited before the
// AST is fully built, as parent pointers may not have been instantiated yet.
//
// Java: enum FeatureCollector.FeatureContext
type FeatureCollector_FeatureContext int32

const (
	FeatureCollector_FeatureContext_PARAM_LIST FeatureCollector_FeatureContext = iota
	FeatureCollector_FeatureContext_CLASS_MEMBERS
	FeatureCollector_FeatureContext_CLASS
	FeatureCollector_FeatureContext_CATCH
	// the top of a block scope, e.g. within an if/while/for loop block, or a
	// plain `{ }` block
	FeatureCollector_FeatureContext_BLOCK_SCOPE
	FeatureCollector_FeatureContext_FUNCTION
	// includes objects like `return {x, y};` but /not/ object destructuring
	// like `const {x} = o;`
	FeatureCollector_FeatureContext_OBJECT_LITERAL
	FeatureCollector_FeatureContext_NONE
)

var featureCollector_FeatureContext_names = [...]string{
	"PARAM_LIST", "CLASS_MEMBERS", "CLASS", "CATCH", "BLOCK_SCOPE", "FUNCTION", "OBJECT_LITERAL", "NONE",
}

// Java: FeatureCollector.FeatureContext.values()
func FeatureCollector_FeatureContext_values() []FeatureCollector_FeatureContext {
	return []FeatureCollector_FeatureContext{
		FeatureCollector_FeatureContext_PARAM_LIST,
		FeatureCollector_FeatureContext_CLASS_MEMBERS,
		FeatureCollector_FeatureContext_CLASS,
		FeatureCollector_FeatureContext_CATCH,
		FeatureCollector_FeatureContext_BLOCK_SCOPE,
		FeatureCollector_FeatureContext_FUNCTION,
		FeatureCollector_FeatureContext_OBJECT_LITERAL,
		FeatureCollector_FeatureContext_NONE,
	}
}

// String is Enum.name().
func (c FeatureCollector_FeatureContext) String() string {
	return featureCollector_FeatureContext_names[c]
}

// Java: FeatureCollector.contextFor(List<FeatureContext>, Node)
func FeatureCollector_ContextFor(contextStack *jutil.ArrayList[FeatureCollector_FeatureContext], node *rhino.Node) FeatureCollector_FeatureContext {
	switch node.GetToken() {
	case rhino.Token_PARAM_LIST:
		return FeatureCollector_FeatureContext_PARAM_LIST
	case rhino.Token_CLASS_MEMBERS:
		return FeatureCollector_FeatureContext_CLASS_MEMBERS
	case rhino.Token_CLASS:
		return FeatureCollector_FeatureContext_CLASS
	case rhino.Token_CATCH:
		return FeatureCollector_FeatureContext_CATCH
	case rhino.Token_BLOCK:
		// a function body is not a block scope - BLOCK is just overloaded. all
		// other references to BLOCK are block scopes.
		if contextStack.GetLast() == FeatureCollector_FeatureContext_FUNCTION {
			return FeatureCollector_FeatureContext_NONE
		}
		return FeatureCollector_FeatureContext_BLOCK_SCOPE
	case rhino.Token_FUNCTION:
		return FeatureCollector_FeatureContext_FUNCTION
	case rhino.Token_OBJECTLIT:
		return FeatureCollector_FeatureContext_OBJECT_LITERAL
	default:
		return FeatureCollector_FeatureContext_NONE
	}
}
