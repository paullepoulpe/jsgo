// Port of com/google/javascript/jscomp/DestructuredTarget.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DestructuredTarget represents a single target inside a destructuring
// pattern, whether another pattern or a lhs expression.
//
// Call InferType to do type inference on the target lazily. This value type is
// designed so that a caller can get information about a target, use that
// information to do additional type inference, and finally call InferType
// if desired without heap-allocating per-target structs, builders, or closures.
type DestructuredTarget struct {
	registry         *rhino.JSTypeRegistry
	eagerType        rhino.JSType
	lazyInferredType func() rhino.JSType
	// The destructuring pattern containing this target.
	pattern rhino.Node
	// Holds the STRING_KEY or COMPUTED_PROPERTY for a target in an object
	// pattern. NilNode for targets in array patterns.
	objectPatternKey rhino.Node
	// The target being assigned to. Can be a destructuring pattern or name,
	// if from a declaration, or an arbitrary lhs expression in an assign.
	node rhino.Node
	// The default value of this target, or NilNode if none. e.g. for
	// `[a = 3] = rhs`, this is `3`.
	defaultValue rhino.Node
	// Whether this is a rest key.
	isRest bool
}

func (d DestructuredTarget) patternType() rhino.JSType {
	if d.lazyInferredType != nil {
		return d.lazyInferredType()
	}
	return d.eagerType
}

// GetComputedProperty returns a COMPUTED_PROP node or NilNode.
//
// Java: DestructuredTarget.getComputedProperty()
func (d DestructuredTarget) GetComputedProperty() rhino.Node {
	if d.HasComputedProperty() {
		return d.objectPatternKey
	}
	return rhino.NilNode
}

// Java: DestructuredTarget.hasComputedProperty()
func (d DestructuredTarget) HasComputedProperty() bool {
	return !d.objectPatternKey.IsNil() && d.objectPatternKey.IsComputedProp()
}

// Java: DestructuredTarget.hasStringKey()
func (d DestructuredTarget) HasStringKey() bool {
	return !d.objectPatternKey.IsNil() && d.objectPatternKey.IsStringKey()
}

// GetStringKey returns a STRING_KEY node or NilNode.
//
// Java: DestructuredTarget.getStringKey()
func (d DestructuredTarget) GetStringKey() rhino.Node {
	if d.HasStringKey() {
		return d.objectPatternKey
	}
	return rhino.NilNode
}

// Java: DestructuredTarget.getDefaultValue()
func (d DestructuredTarget) GetDefaultValue() rhino.Node {
	return d.defaultValue
}

// Java: DestructuredTarget.hasDefaultValue()
func (d DestructuredTarget) HasDefaultValue() bool {
	return !d.defaultValue.IsNil()
}

// Java: DestructuredTarget.getNode()
func (d DestructuredTarget) GetNode() rhino.Node {
	return d.node
}

// Java: DestructuredTarget.createTarget(JSTypeRegistry, JSType, Node)
func DestructuredTarget_CreateTarget(
	registry *rhino.JSTypeRegistry, destructuringPatternType rhino.JSType, destructuringChild rhino.Node) DestructuredTarget {
	parent := destructuringChild.GetParent()
	if !parent.IsDestructuringPattern() {
		failNotDestructuringChild(destructuringChild)
	}
	return createDestructuredTargetInPattern(registry, destructuringPatternType, nil, parent, destructuringChild)
}

// DestructuredTarget_CreateTargetWithSupplier converts a given child of a
// destructuring pattern (in the AST) to an instance of this value type.
//
// NOTE: does not accept EMPTY nodes.
//
// Java: DestructuredTarget.createTarget(JSTypeRegistry, Supplier<JSType>, Node)
func DestructuredTarget_CreateTargetWithSupplier(
	registry *rhino.JSTypeRegistry, destructuringPatternType func() rhino.JSType, destructuringChild rhino.Node) DestructuredTarget {
	parent := destructuringChild.GetParent()
	if !parent.IsDestructuringPattern() {
		failNotDestructuringChild(destructuringChild)
	}
	return createDestructuredTargetInPattern(registry, nil, destructuringPatternType, parent, destructuringChild)
}

func createDestructuredTargetInPattern(
	registry *rhino.JSTypeRegistry,
	eagerType rhino.JSType,
	lazyInferredType func() rhino.JSType,
	pattern rhino.Node,
	destructuringChild rhino.Node,
) DestructuredTarget {
	d := DestructuredTarget{
		registry:         registry,
		eagerType:        eagerType,
		lazyInferredType: lazyInferredType,
		pattern:          pattern,
	}
	switch destructuringChild.GetToken() {
	case rhino.Token_STRING_KEY:
		// const {objectLiteralKey: x} = ...
		d.objectPatternKey = destructuringChild
		value := destructuringChild.GetFirstChild()
		if value.IsDefaultValue() {
			d.node = value.GetFirstChild()
			d.defaultValue = value.GetSecondChild()
		} else {
			d.node = value
		}
	case rhino.Token_COMPUTED_PROP:
		// const {['objectLiteralKey']: x} = ...
		d.objectPatternKey = destructuringChild
		value := destructuringChild.GetSecondChild()
		if value.IsDefaultValue() {
			d.node = value.GetFirstChild()
			d.defaultValue = value.GetSecondChild()
		} else {
			d.node = value
		}
	case rhino.Token_OBJECT_PATTERN, // const [{x}] = ...
		rhino.Token_ARRAY_PATTERN, // const [[x]] = ...
		rhino.Token_NAME,          // const [x] = ...
		rhino.Token_GETELEM,       // [obj[3]] = ...
		rhino.Token_GETPROP:       // [this.x] = ...
		d.node = destructuringChild
	case rhino.Token_DEFAULT_VALUE:
		// const [x = 3] = ...
		d.node = destructuringChild.GetFirstChild()
		d.defaultValue = destructuringChild.GetSecondChild()
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		// const [...x] = ...
		// const {...x} = ...
		d.node = destructuringChild.GetFirstChild()
		d.isRest = true
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Unexpected child of destructuring pattern " + destructuringChild.String()))
	}
	rhino.CheckNotNullMsg(d.node, "Must set a node")
	return d
}

func failNotDestructuringPattern(pattern rhino.Node) {
	jutil.CheckArgumentMsg1(false, "%s", pattern)
}

func failNotDestructuringChild(child rhino.Node) {
	jutil.CheckArgumentMsg1(false, "%s", child)
}

// ForEachDestructuredTarget visits every non-EMPTY target directly in pattern
// without allocating a slice or per-target heap objects.
func ForEachDestructuredTarget(
	registry *rhino.JSTypeRegistry,
	patternType rhino.JSType,
	pattern rhino.Node,
	visit func(DestructuredTarget),
) {
	if !pattern.IsDestructuringPattern() {
		failNotDestructuringPattern(pattern)
	}
	for child := pattern.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if !child.IsEmpty() {
			visit(createDestructuredTargetInPattern(registry, patternType, nil, pattern, child))
		}
	}
}

// DestructuredTarget_CreateAllNonEmptyTargetsInPattern returns all the
// targets directly in the given pattern, except for EMPTY nodes.
//
// EMPTY nodes occur in array patterns with elisions, e.g. `[, , a] = []`
//
// Java: DestructuredTarget.createAllNonEmptyTargetsInPattern(JSTypeRegistry, JSType, Node)
func DestructuredTarget_CreateAllNonEmptyTargetsInPattern(
	registry *rhino.JSTypeRegistry, patternType rhino.JSType, pattern rhino.Node) []DestructuredTarget {
	if !pattern.IsDestructuringPattern() {
		failNotDestructuringPattern(pattern)
	}
	var targets []DestructuredTarget
	for child := pattern.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if !child.IsEmpty() {
			targets = append(targets, createDestructuredTargetInPattern(registry, patternType, nil, pattern, child))
		}
	}
	return targets
}

// DestructuredTarget_CreateAllNonEmptyTargetsInPatternWithSupplier returns
// all the targets directly in the given pattern, except for EMPTY nodes.
//
// Java: DestructuredTarget.createAllNonEmptyTargetsInPattern(JSTypeRegistry, Supplier<JSType>, Node)
func DestructuredTarget_CreateAllNonEmptyTargetsInPatternWithSupplier(
	registry *rhino.JSTypeRegistry, patternType func() rhino.JSType, pattern rhino.Node) []DestructuredTarget {
	if !pattern.IsDestructuringPattern() {
		failNotDestructuringPattern(pattern)
	}
	var targets []DestructuredTarget
	for child := pattern.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if !child.IsEmpty() {
			targets = append(targets, createDestructuredTargetInPattern(registry, nil, patternType, pattern, child))
		}
	}
	return targets
}

// InferTypeWithoutUsingDefaultValue infers the type of this target before
// the default value (if any) is applied.
//
// e.g. given `a` in const {a = 3} = {}; returns `undefined`
//
// Java: DestructuredTarget.inferTypeWithoutUsingDefaultValue()
func (d DestructuredTarget) InferTypeWithoutUsingDefaultValue() rhino.JSType {
	if d.pattern.IsObjectPattern() {
		return d.InferObjectPatternKeyType()
	}
	return d.InferArrayPatternTargetType()
}

// InferType infers the type of this target.
//
// e.g. given `a` in const {a = 3} = {}; returns `number`
//
// Java: DestructuredTarget.inferType()
func (d DestructuredTarget) InferType() rhino.JSType {
	inferredType := d.InferTypeWithoutUsingDefaultValue()
	if !inferredType.IsUnknownType() && d.HasDefaultValue() {
		defaultValueType := d.GetDefaultValue().GetJSType()
		if defaultValueType == nil {
			defaultValueType = d.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		// We effectively replace '|undefined" with '|typeOfDefaultValue'
		return d.registry.CreateUnionType(inferredType.RestrictByNotUndefined(), defaultValueType)
	}
	return inferredType
}

// Java: DestructuredTarget.inferObjectPatternKeyType()
func (d DestructuredTarget) InferObjectPatternKeyType() rhino.JSType {
	patternType := d.patternType()
	if d.isRest {
		// TODO(b/128355893): Do smarter inferrence. There are a lot of potential issues with
		// inference on object-rest, so for now we just give up and say `Object`.
		return d.registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	}

	if patternType == nil || patternType.IsUnknownType() {
		return d.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	switch d.objectPatternKey.GetToken() {
	case rhino.Token_STRING_KEY:
		propertyType := patternType.FindPropertyType(d.objectPatternKey.GetString())
		if propertyType != nil {
			return propertyType
		}
		return d.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	case rhino.Token_COMPUTED_PROP:
		return patternType.
			GetTemplateTypeMap().
			GetResolvedTemplateType(d.registry.GetObjectElementKey())
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected key " + d.objectPatternKey.String()))
	}
}

// Java: DestructuredTarget.inferArrayPatternTargetType()
func (d DestructuredTarget) InferArrayPatternTargetType() rhino.JSType {
	patternType := d.patternType()

	// e.g. get `number` from `!Iterable<number>`
	templateTypeOfIterable := patternType.
		GetTemplateTypeMap().
		GetResolvedTemplateType(d.registry.GetIterableValueTemplate())

	if d.isRest {
		// return `!Array<number>`
		return d.registry.CreateTemplatizedType(
			d.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), templateTypeOfIterable)
	}
	return templateTypeOfIterable
}
