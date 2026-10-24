// Port of com/google/javascript/rhino/StaticScope.java @ closure-compiler v20261021

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

// StaticScope must be implemented by any object that defines variables for
// the purposes of static analysis. It is distinguished from the Scriptable
// class that Rhino normally uses to represent a run-time scope.
type StaticScope interface {
	// GetRootNode returns the root node associated with this scope. May be
	// null.
	GetRootNode() Node

	// GetParentScope returns the scope enclosing this one or null if none.
	GetParentScope() StaticScope

	// GetSlot returns any defined slot within this scope for this name. This
	// call continues searching through parent scopes if a slot with this name
	// is not found in the current scope. Returns null if no definition
	// exists.
	GetSlot(name string) StaticSlot

	// GetOwnSlot is like getSlot but does not recurse into parent scopes.
	GetOwnSlot(name string) StaticSlot

	// GetTopmostScopeOfEventualDeclaration returns the topmost slot containing
	// this name, or null if no slots do. May return true in cases where
	// getSlot returns null, though. Do not rely on this method if you need an
	// actual slot.
	//
	// This method is intended for use while scopes are still being built,
	// hence the name 'eventual' declaration. Once scope building is complete,
	// the scope returned from this method must be equivalent to
	// "getSlot(name).getScope()" or null.
	//
	// A Java default method: implementations that do not override it call
	// StaticScope_GetTopmostScopeOfEventualDeclaration.
	GetTopmostScopeOfEventualDeclaration(name string) StaticScope
}

// StaticScope_GetTopmostScopeOfEventualDeclaration is the default body of
// StaticScope.getTopmostScopeOfEventualDeclaration.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func StaticScope_GetTopmostScopeOfEventualDeclaration(s StaticScope, name string) StaticScope {
	slot := s.GetOwnSlot(name)
	if slot != nil {
		return slot.GetScope()
	}
	if s.GetParentScope() != nil {
		return s.GetParentScope().GetTopmostScopeOfEventualDeclaration(name)
	}
	return nil
}
