// Port of com/google/javascript/rhino/jstype/JSTypeResolver.java @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// JSTypeResolver is a state machine for resolving all JSType instances.
//
// This object shares the lifecycle of a registry. Consider it as a narrow
// interface into the registry for the purpose of resolving types.
//
// Every time a new type is constructed it is (by cooperation of the JSType
// subclasses} added to this resolver. Depending on the state of the resolver
// at the time of addition, the type will either be eagerly resolved, or
// stored for later resolution.
type JSTypeResolver struct {
	registry *JSTypeRegistry

	// The sequence of instantiated types.
	//
	// This allows verification that every new type is captured by this
	// resolver. In general this stack should never be more than a handful of
	// types.
	captureStack []JSType

	// The sequence of types to resolve when the resolver is closed.
	resolutionQueue []JSType

	state JSTypeResolver_State
}

// JSTypeResolver_Closer is a signal to resolve all types known to and close
// the owning resolver.
//
// This is intended to be used in a try-with-resources statement. This
// approach was selected over accepting a callback because it supports some
// level of compile-time enforcement while allowing that enforcement to be
// suppressed in tests.
//
// Go: `try (Closer c = r.openForDefinition()) { ... }` is
// `c := r.OpenForDefinition(); defer c.Close()`.
type JSTypeResolver_Closer struct {
	outer  *JSTypeResolver
	hasRun bool
}

// Java: JSTypeResolver.Closer.close()
func (c *JSTypeResolver_Closer) Close() {
	jutil.CheckState(!c.hasRun)
	c.hasRun = true
	c.outer.ResolveAll()
}

// JSTypeResolver_State is JSTypeResolver.State.
type JSTypeResolver_State int32

const (
	JSTypeResolver_State_CLOSED JSTypeResolver_State = iota
	JSTypeResolver_State_OPEN
	JSTypeResolver_State_CLOSING
)

// Java: JSTypeResolver.create(JSTypeRegistry)
func JSTypeResolver_Create(registry *JSTypeRegistry) *JSTypeResolver {
	if registry.GetResolver() != nil {
		return registry.GetResolver()
	}

	return NewJSTypeResolver(registry)
}

// Java: JSTypeResolver(JSTypeRegistry)
func NewJSTypeResolver(registry *JSTypeRegistry) *JSTypeResolver {
	return &JSTypeResolver{registry: registry, state: JSTypeResolver_State_CLOSED}
}

// AddUnresolved records that a new type has been instantiated.
//
// This should only be called in the JSType constructor, with the intention
// of capturing all new types.
//
// Java: JSTypeResolver.addUnresolved(JSType)
func (r *JSTypeResolver) AddUnresolved(captured JSType) {
	// This method exists to make it easier to debug the source of unresolved
	// types in the future.
	//
	// In theory we should be able to capture all types using
	// `resolveIfClosed`, but this way we get some guarantees against misuse.
	// Separating capturing and resolving allows us to verify our invariants.
	jutil.CheckState(!captured.IsResolved())
	r.captureStack = append(r.captureStack, captured)
}

// ResolveIfClosed: if captured is finished construction, and the resolver is
// closed, resolve it.
//
// caller is used to cooperatively determine whether this call came from the
// final constructor for type.
//
// This should only be called at the end of JSType constructors which invoke
// super. It should be called from all such constructors. Some constructors
// (e.g. AllType) may choose to delegate that call through
// JSType#eagerlyResolveToSelf.
//
// Java: JSTypeResolver.resolveIfClosed(JSType, JSTypeClass)
func (r *JSTypeResolver) ResolveIfClosed(captured JSType, caller JSTypeClass) {
	if !caller.IsTypeOf(captured) {
		return
	}

	if len(r.captureStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	expected := r.captureStack[len(r.captureStack)-1]
	r.captureStack[len(r.captureStack)-1] = nil
	r.captureStack = r.captureStack[:len(r.captureStack)-1]
	if captured != expected {
		jutil.CheckStateMsg2(false, "Captured %s; Expected %s", captured, expected)
	}

	switch r.state {
	case JSTypeResolver_State_CLOSED, JSTypeResolver_State_CLOSING:
		r.DoResolve(captured)
	case JSTypeResolver_State_OPEN:
		r.resolutionQueue = append(r.resolutionQueue, captured)
	}
}

// OpenForDefinition allows new types to be created without eagerly
// resolving them.
//
// This is required by, and only really useful for, type definition. Types
// being derived from previously defined types should not contain reference
// cycles, and so should be eagerly resolvable.
//
// Java: JSTypeResolver.openForDefinition()
func (r *JSTypeResolver) OpenForDefinition() *JSTypeResolver_Closer {
	jutil.CheckState(r.state == JSTypeResolver_State_CLOSED)
	jutil.CheckState(len(r.captureStack) == 0)

	r.state = JSTypeResolver_State_OPEN
	return &JSTypeResolver_Closer{outer: r}
}

// Java: JSTypeResolver.resolveAll()
func (r *JSTypeResolver) ResolveAll() {
	jutil.CheckState(r.state == JSTypeResolver_State_OPEN)
	jutil.CheckState(len(r.captureStack) == 0)

	r.state = JSTypeResolver_State_CLOSING

	for len(r.resolutionQueue) != 0 {
		first := r.resolutionQueue[0]
		r.resolutionQueue[0] = nil
		r.resolutionQueue = r.resolutionQueue[1:]
		r.DoResolve(first)
	}

	// resolutionQueue scales with the size of the application, so it needs
	// to be GC'd.
	r.resolutionQueue = nil

	r.state = JSTypeResolver_State_CLOSED

	// TODO(sdh): Stop doing this here. It's obviously the wrong place. By
	// default, the global "this" type is just an anonymous object. If the user
	// has defined a Window type, make the Window the implicit prototype of
	// "this".
	globalThis := r.registry.GetNativeType(JSTypeNative_GLOBAL_THIS).(PrototypeObjectType)
	windowType := r.registry.GetGlobalType("Window")
	if globalThis.IsUnknownType() {
		windowObjType := ObjectType_Cast(windowType)
		if windowObjType != nil {
			globalThis.SetImplicitPrototype(windowObjType)
		} else {
			globalThis.SetImplicitPrototype(
				r.registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE))
		}
	}
}

// Java: JSTypeResolver.doResolve(JSType)
func (r *JSTypeResolver) DoResolve(type_ JSType) {
	type_.Resolve(r.registry.GetErrorReporter())
}

// AssertLegalToResolveTypes asserts that it's legal to call
// JSType#resolve.
//
// The intent is to enforce that while a type is being resolved, any new
// types synthesized will be immediately resolved.
//
// Java: JSTypeResolver.assertLegalToResolveTypes()
func (r *JSTypeResolver) AssertLegalToResolveTypes() {
	jutil.CheckStateMsg(
		r.state != JSTypeResolver_State_OPEN, "Types cannot be resolved while the registry is open")
}
