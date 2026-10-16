// Port of com/google/javascript/rhino/jstype/JSTypeResolver.java
// @ closure-compiler v20260415
//
// JSTypeResolver is a state machine for late-resolving forward-referenced
// JSType instances. It shares the lifecycle of a Registry: every time a
// new type is constructed it is added to the resolver; depending on the
// resolver's state at the time of addition, the type is either eagerly
// resolved or stored for later bulk resolution.
//
// The original Java file is:
//
//	com/google/javascript/rhino/jstype/JSTypeResolver.java
//
// Java licence block:
//
// ***** BEGIN LICENSE BLOCK *****
// Version: MPL 1.1/GPL 2.0
//
// The contents of this file are subject to the Mozilla Public License
// Version 1.1 (the "License"); you may not use this file except in
// compliance with the License. You may obtain a copy of the License at
// http://www.mozilla.org/MPL/
//
// Software distributed under the License is distributed on an "AS IS"
// basis, WITHOUT WARRANTY OF ANY KIND, either express or implied. See
// the License for the specific language governing rights and limitations
// under the License.
//
// The Original Code is Rhino code, released May 6, 1999.
//
// The Initial Developer of the Original Code is Netscape Communications
// Corporation. Portions created by the Initial Developer are Copyright
// (C) 1997-1999 the Initial Developer. All Rights Reserved.
//
// Contributor(s): Bob Jervis, Google Inc.
// ***** END LICENSE BLOCK *****

package jstype

// resolverState mirrors the private State enum in JSTypeResolver.java.
// CLOSED  – no open definition block; new types are resolved immediately.
// OPEN    – inside openForDefinition; new types are queued for later.
// CLOSING – resolveAll is draining the queue; new types resolved immediately.
type resolverState int

const (
	resolverStateClosed  resolverState = iota // CLOSED
	resolverStateOpen                          // OPEN
	resolverStateClosing                       // CLOSING
)

// JSTypeResolver is a state machine for resolving all JSType instances.
//
// This object shares the lifecycle of a Registry. Consider it a narrow
// interface into the registry for the purpose of resolving types.
//
// Every time a new type is constructed it is (by cooperation of the JSType
// subclasses) added to this resolver. Depending on the state of the
// resolver at the time of addition, the type will either be eagerly
// resolved, or stored for later resolution.
//
// Ported from com/google/javascript/rhino/jstype/JSTypeResolver.java.
type JSTypeResolver struct {
	registry *Registry

	// captureStack records types as they are instantiated (addUnresolved).
	// resolveIfClosed removes the top entry when the outermost constructor
	// has finished. In general this should never hold more than a handful
	// of entries. Mirrors ArrayDeque<JSType> captureStack in Java.
	captureStack []JSType

	// resolutionQueue holds types that were created while the resolver was
	// OPEN; they are bulk-resolved when Close() is called.
	// Mirrors ArrayDeque<JSType> resolutionQueue in Java.
	resolutionQueue []JSType

	state resolverState
}

// Closer is a signal to resolve all queued types and close the owning
// resolver. It mirrors the inner class JSTypeResolver.Closer in Java,
// which implements AutoCloseable and is meant for use in a
// try-with-resources block. In Go it is used via defer or an explicit
// call to Close().
//
// Ported from the inner class com/google/javascript/rhino/jstype/JSTypeResolver.Closer.
type Closer struct {
	hasRun   bool
	resolver *JSTypeResolver
}

// Close resolves all types queued during the open-for-definition phase
// and returns the resolver to the CLOSED state. Panics if called more
// than once (mirrors Java's checkState(!hasRun)).
func (c *Closer) Close() {
	if c.hasRun {
		panic("jstype.Closer.Close called more than once")
	}
	c.hasRun = true
	c.resolver.resolveAll()
}

// Create returns the resolver already attached to registry, or allocates
// a fresh one. Mirrors JSTypeResolver.create(JSTypeRegistry) in Java.
func Create(registry *Registry) *JSTypeResolver {
	if registry.GetResolver() != nil {
		return registry.GetResolver()
	}
	return &JSTypeResolver{registry: registry}
}

// AddUnresolved records that a new type has been instantiated.
//
// This should only be called once per type construction, with the intent
// of capturing all new types so the resolver can verify its invariants.
// In Java this is called from the JSType base constructor.
//
// TODO(tier-N): when Go JSType gains isResolved(), add the
// checkState(!captured.isResolved()) guard that Java enforces.
//
// Ported from JSTypeResolver.addUnresolved(JSType).
func (r *JSTypeResolver) AddUnresolved(captured JSType) {
	// TODO(tier-N): checkState(!captured.isResolved())
	r.captureStack = append(r.captureStack, captured)
}

// ResolveIfClosed resolves captured if the resolver is not OPEN.
//
// caller is used to cooperatively decide whether this call came from
// the outermost constructor for captured. In Java this is done via
// JSTypeClass.isTypeOf which checks the runtime class of the object.
// In Go the JSTypeClass mechanism doesn't exist yet, so the check is
// approximated: this method resolves whenever the capture stack records
// captured at its top (i.e. we are in the outermost constructor call).
//
// TODO(tier-N): port JSTypeClass and restore the full caller.isTypeOf
// guard that Java uses to pick the right constructor level.
//
// Ported from JSTypeResolver.resolveIfClosed(JSType, JSTypeClass).
func (r *JSTypeResolver) ResolveIfClosed(captured JSType) {
	if len(r.captureStack) == 0 {
		return
	}
	top := r.captureStack[len(r.captureStack)-1]
	if top != captured {
		return
	}
	r.captureStack = r.captureStack[:len(r.captureStack)-1]

	switch r.state {
	case resolverStateClosed, resolverStateClosing:
		r.doResolve(captured)
	case resolverStateOpen:
		r.resolutionQueue = append(r.resolutionQueue, captured)
	}
}

// OpenForDefinition allows new types to be created without eagerly
// resolving them. Returns a Closer whose Close() method resolves all
// queued types and returns to the CLOSED state.
//
// This is required by, and only useful for, type definition: types being
// derived from previously defined types should not contain reference
// cycles and so should be eagerly resolvable.
//
// Mirrors JSTypeResolver.openForDefinition() which is annotated
// @MustBeClosed and intended for try-with-resources.
//
// Ported from JSTypeResolver.openForDefinition().
func (r *JSTypeResolver) OpenForDefinition() *Closer {
	if r.state != resolverStateClosed {
		panic("jstype.JSTypeResolver.OpenForDefinition: resolver is not CLOSED")
	}
	if len(r.captureStack) != 0 {
		panic("jstype.JSTypeResolver.OpenForDefinition: captureStack is not empty")
	}
	r.state = resolverStateOpen
	return &Closer{resolver: r}
}

// resolveAll drains the resolutionQueue and returns the resolver to
// CLOSED. Called by Closer.Close(). Mirrors JSTypeResolver.resolveAll().
func (r *JSTypeResolver) resolveAll() {
	if r.state != resolverStateOpen {
		panic("jstype.JSTypeResolver.resolveAll: resolver is not OPEN")
	}
	if len(r.captureStack) != 0 {
		panic("jstype.JSTypeResolver.resolveAll: captureStack is not empty")
	}

	r.state = resolverStateClosing

	for len(r.resolutionQueue) > 0 {
		t := r.resolutionQueue[0]
		r.resolutionQueue = r.resolutionQueue[1:]
		r.doResolve(t)
	}

	// Reset the queue so it can be GC'd (matches Java comment:
	// "resolutionQueue scales with the size of the application").
	r.resolutionQueue = nil

	r.state = resolverStateClosed

	// TODO(tier-N): Stop doing this here. It's obviously the wrong place.
	// Java's resolveAll sets globalThis's implicit prototype to Window
	// once the Window type is available in the registry:
	//
	//   PrototypeObjectType globalThis =
	//       (PrototypeObjectType) registry.getNativeType(GLOBAL_THIS);
	//   JSType windowType = registry.getGlobalType("Window");
	//   if (globalThis.isUnknownType()) {
	//     ObjectType windowObjType = ObjectType.cast(windowType);
	//     if (windowObjType != null) {
	//       globalThis.setImplicitPrototype(windowObjType);
	//     } else {
	//       globalThis.setImplicitPrototype(
	//           registry.getNativeObjectType(JSTypeNative.OBJECT_TYPE));
	//     }
	//   }
	//
	// This requires JSTypeNative, PrototypeObjectType, and getNativeType /
	// getGlobalType on Registry — none of which are ported yet.
}

// doResolve resolves a single type. Mirrors JSTypeResolver.doResolve(JSType).
//
// TODO(tier-N): call type.Resolve(registry.GetErrorReporter()) once
// JSType gains a Resolve method and Registry exposes GetErrorReporter.
func (r *JSTypeResolver) doResolve(_ JSType) {
	// No-op until JSType.Resolve and ErrorReporter are ported.
}

// AssertLegalToResolveTypes panics when the resolver is OPEN, enforcing
// that JSType.resolve must not be called during a definition block.
// Mirrors JSTypeResolver.assertLegalToResolveTypes().
func (r *JSTypeResolver) AssertLegalToResolveTypes() {
	if r.state == resolverStateOpen {
		panic("jstype: types cannot be resolved while the registry is open")
	}
}
