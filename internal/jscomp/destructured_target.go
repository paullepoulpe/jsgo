// Port of com/google/javascript/jscomp/DestructuredTarget.java @ closure-compiler v20261015

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
// Call InferType to do type inference on the target lazily. This class is
// designed so that a caller can get information about a target, use that
// information to do additional type inference, and finally call InferType
// if desired.
type DestructuredTarget struct {
	registry *rhino.JSTypeRegistry
	// Holds the STRING_KEY or COMPUTED_PROPERTY for a target in an object
	// pattern. Null for targets in array patterns.
	objectPatternKey *rhino.Node
	// The target being assigned to. Can be a destructuring pattern or name,
	// if from a declaration, or an arbitrary lhs expression in an assign.
	node *rhino.Node
	// A supplier to get the type of the pattern containing this target. e.g.
	// for `a` in `const {a} = {a: 3}`, the supplier provides the record type
	// `{a: number}`
	//
	// This is only called by inferType, not when calling createTarget.
	patternTypeSupplier func() rhino.JSType

	// The default value of this target, or null if none. e.g. for
	// `[a = 3] = rhs`, this is `3`.
	defaultValue *rhino.Node

	// Whether this is a rest key
	isRest bool

	// The destructuring pattern containing this target
	pattern *rhino.Node
}

// Java: DestructuredTarget(Node, Node, Node, Node, JSTypeRegistry, Supplier<JSType>, boolean)
func newDestructuredTarget(
	pattern *rhino.Node,
	node *rhino.Node,
	defaultValue *rhino.Node,
	objectPatternKey *rhino.Node,
	registry *rhino.JSTypeRegistry,
	patternTypeSupplier func() rhino.JSType,
	isRest bool) *DestructuredTarget {
	d := &DestructuredTarget{}
	d.pattern = pattern
	d.node = node
	d.objectPatternKey = objectPatternKey
	d.registry = registry
	d.patternTypeSupplier = patternTypeSupplier
	d.isRest = isRest
	d.defaultValue = defaultValue
	return d
}

// GetComputedProperty returns a COMPUTED_PROP node or null
//
// Java: DestructuredTarget.getComputedProperty()
func (d *DestructuredTarget) GetComputedProperty() *rhino.Node {
	if d.HasComputedProperty() {
		return d.objectPatternKey
	}
	return nil
}

// Java: DestructuredTarget.hasComputedProperty()
func (d *DestructuredTarget) HasComputedProperty() bool {
	return d.objectPatternKey != nil && d.objectPatternKey.IsComputedProp()
}

// Java: DestructuredTarget.hasStringKey()
func (d *DestructuredTarget) HasStringKey() bool {
	return d.objectPatternKey != nil && d.objectPatternKey.IsStringKey()
}

// GetStringKey returns a STRING_KEY node or null
//
// Java: DestructuredTarget.getStringKey()
func (d *DestructuredTarget) GetStringKey() *rhino.Node {
	if d.HasStringKey() {
		return d.objectPatternKey
	}
	return nil
}

// Java: DestructuredTarget.getDefaultValue()
func (d *DestructuredTarget) GetDefaultValue() *rhino.Node {
	return d.defaultValue
}

// Java: DestructuredTarget.hasDefaultValue()
func (d *DestructuredTarget) HasDefaultValue() bool {
	return d.defaultValue != nil
}

// Java: DestructuredTarget.getNode()
func (d *DestructuredTarget) GetNode() *rhino.Node {
	return d.node
}

// DestructuredTarget_Builder is DestructuredTarget.Builder.
type DestructuredTarget_Builder struct {
	registry            *rhino.JSTypeRegistry
	patternTypeSupplier func() rhino.JSType
	pattern             *rhino.Node
	node                *rhino.Node
	defaultValue        *rhino.Node
	objectPatternKey    *rhino.Node
	isRest              bool
}

// Java: DestructuredTarget.Builder(JSTypeRegistry, Node, Supplier<JSType>)
func NewDestructuredTarget_Builder(
	registry *rhino.JSTypeRegistry, pattern *rhino.Node, patternTypeSupplier func() rhino.JSType) *DestructuredTarget_Builder {
	b := &DestructuredTarget_Builder{}
	b.defaultValue = nil
	b.objectPatternKey = nil
	b.isRest = false
	b.registry = registry
	b.patternTypeSupplier = patternTypeSupplier
	b.pattern = pattern
	return b
}

// Java: DestructuredTarget.Builder.setNode(Node)
func (b *DestructuredTarget_Builder) SetNode(node *rhino.Node) *DestructuredTarget_Builder {
	b.node = node
	return b
}

// Java: DestructuredTarget.Builder.setDefaultValue(Node)
func (b *DestructuredTarget_Builder) SetDefaultValue(defaultValue *rhino.Node) *DestructuredTarget_Builder {
	b.defaultValue = defaultValue
	return b
}

// Java: DestructuredTarget.Builder.setObjectPatternKey(Node)
func (b *DestructuredTarget_Builder) SetObjectPatternKey(objectPatternKey *rhino.Node) *DestructuredTarget_Builder {
	b.objectPatternKey = objectPatternKey
	return b
}

// Java: DestructuredTarget.Builder.setIsRest(boolean)
func (b *DestructuredTarget_Builder) SetIsRest(isRest bool) *DestructuredTarget_Builder {
	b.isRest = isRest
	return b
}

// Java: DestructuredTarget.Builder.build()
func (b *DestructuredTarget_Builder) Build() *DestructuredTarget {
	jutil.CheckNotNullMsg(b.node, "Must set a node")

	return newDestructuredTarget(
		b.pattern, b.node, b.defaultValue, b.objectPatternKey, b.registry, b.patternTypeSupplier, b.isRest)
}

// Java: DestructuredTarget.createTarget(JSTypeRegistry, JSType, Node)
func DestructuredTarget_CreateTarget(
	registry *rhino.JSTypeRegistry, destructuringPatternType rhino.JSType, destructuringChild *rhino.Node) *DestructuredTarget {
	return DestructuredTarget_CreateTargetWithSupplier(
		registry, func() rhino.JSType { return destructuringPatternType }, destructuringChild)
}

// DestructuredTarget_CreateTargetWithSupplier converts a given child of a
// destructuring pattern (in the AST) to an instance of this class.
//
// NOTE: does not accept EMPTY nodes
//
// Java: DestructuredTarget.createTarget(JSTypeRegistry, Supplier<JSType>, Node)
func DestructuredTarget_CreateTargetWithSupplier(
	registry *rhino.JSTypeRegistry, destructuringPatternType func() rhino.JSType, destructuringChild *rhino.Node) *DestructuredTarget {
	jutil.CheckArgumentMsg1(destructuringChild.GetParent().IsDestructuringPattern(), "%s", destructuringChild)

	builder := NewDestructuredTarget_Builder(registry, destructuringChild.GetParent(), destructuringPatternType)
	var value *rhino.Node
	switch destructuringChild.GetToken() {
	case rhino.Token_STRING_KEY:
		// const {objectLiteralKey: x} = ...
		builder.SetObjectPatternKey(destructuringChild)
		value = destructuringChild.GetFirstChild()
		if value.IsDefaultValue() {
			builder.SetNode(value.GetFirstChild())
			builder.SetDefaultValue(value.GetSecondChild())
		} else {
			builder.SetNode(value)
		}
	case rhino.Token_COMPUTED_PROP:
		// const {['objectLiteralKey']: x} = ...
		builder.SetObjectPatternKey(destructuringChild)
		value = destructuringChild.GetSecondChild()
		if value.IsDefaultValue() {
			builder.SetNode(value.GetFirstChild())
			builder.SetDefaultValue(value.GetSecondChild())
		} else {
			builder.SetNode(value)
		}
	case rhino.Token_OBJECT_PATTERN, // const [{x}] = ...
		rhino.Token_ARRAY_PATTERN, // const [[x]] = ...
		rhino.Token_NAME,          // const [x] = ...
		rhino.Token_GETELEM,       // [obj[3]] = ...
		rhino.Token_GETPROP:       // [this.x] = ...
		builder.SetNode(destructuringChild)
	case rhino.Token_DEFAULT_VALUE:
		// const [x = 3] = ...
		builder.SetNode(destructuringChild.GetFirstChild())
		builder.SetDefaultValue(destructuringChild.GetSecondChild())
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		// const [...x] = ...
		// const {...x} = ...
		builder.SetNode(destructuringChild.GetFirstChild())
		builder.SetIsRest(true)
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Unexpected child of destructuring pattern " + destructuringChild.String()))
	}
	return builder.Build()
}

// DestructuredTarget_CreateAllNonEmptyTargetsInPattern returns all the
// targets directly in the given pattern, except for EMPTY nodes
//
// EMPTY nodes occur in array patterns with elisions, e.g. `[, , a] = []`
//
// Java: DestructuredTarget.createAllNonEmptyTargetsInPattern(JSTypeRegistry, JSType, Node)
func DestructuredTarget_CreateAllNonEmptyTargetsInPattern(
	registry *rhino.JSTypeRegistry, patternType rhino.JSType, pattern *rhino.Node) []*DestructuredTarget {
	return DestructuredTarget_CreateAllNonEmptyTargetsInPatternWithSupplier(
		registry, func() rhino.JSType { return patternType }, pattern)
}

// DestructuredTarget_CreateAllNonEmptyTargetsInPatternWithSupplier returns
// all the targets directly in the given pattern, except for EMPTY nodes
//
// EMPTY nodes occur in array patterns with elisions, e.g. `[, , a] = []`
//
// Java: DestructuredTarget.createAllNonEmptyTargetsInPattern(JSTypeRegistry, Supplier<JSType>, Node)
func DestructuredTarget_CreateAllNonEmptyTargetsInPatternWithSupplier(
	registry *rhino.JSTypeRegistry, patternType func() rhino.JSType, pattern *rhino.Node) []*DestructuredTarget {
	jutil.CheckArgumentMsg1(pattern.IsDestructuringPattern(), "%s", pattern)
	builder := []*DestructuredTarget{}
	for child := pattern.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsEmpty() {
			continue
		}

		builder = append(builder, DestructuredTarget_CreateTargetWithSupplier(registry, patternType, child))
	}
	return builder
}

// InferTypeWithoutUsingDefaultValue infers the type of this target before
// the default value (if any) is applied.
//
// e.g. given `a` in const {a = 3} = {}; returns `undefined`
//
// Java: DestructuredTarget.inferTypeWithoutUsingDefaultValue()
func (d *DestructuredTarget) InferTypeWithoutUsingDefaultValue() rhino.JSType {
	if d.pattern.IsObjectPattern() {
		return d.InferObjectPatternKeyType()
	} else {
		return d.InferArrayPatternTargetType()
	}
}

// InferType infers the type of this target
//
// e.g. given `a` in const {a = 3} = {}; returns `number`
//
// Java: DestructuredTarget.inferType()
func (d *DestructuredTarget) InferType() rhino.JSType {
	inferredType := d.InferTypeWithoutUsingDefaultValue()
	if !inferredType.IsUnknownType() && d.HasDefaultValue() {
		defaultValueType := d.GetDefaultValue().GetJSType()
		if defaultValueType == nil {
			defaultValueType = d.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		// We effectively replace '|undefined" with '|typeOfDefaultValue'
		return d.registry.CreateUnionType(inferredType.RestrictByNotUndefined(), defaultValueType)
	} else {
		return inferredType
	}
}

// Java: DestructuredTarget.inferObjectPatternKeyType()
func (d *DestructuredTarget) InferObjectPatternKeyType() rhino.JSType {
	patternType := d.patternTypeSupplier()
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
		if patternType != nil {
			return patternType.
				GetTemplateTypeMap().
				GetResolvedTemplateType(d.registry.GetObjectElementKey())
		}
		return d.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected key " + d.objectPatternKey.String()))
	}
}

// Java: DestructuredTarget.inferArrayPatternTargetType()
func (d *DestructuredTarget) InferArrayPatternTargetType() rhino.JSType {
	patternType := d.patternTypeSupplier()

	// e.g. get `number` from `!Iterable<number>`
	templateTypeOfIterable := patternType.
		GetTemplateTypeMap().
		GetResolvedTemplateType(d.registry.GetIterableValueTemplate())

	if d.isRest {
		// return `!Array<number>`
		return d.registry.CreateTemplatizedType(
			d.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), templateTypeOfIterable)
	} else {
		return templateTypeOfIterable
	}
}
