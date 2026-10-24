// Port of com/google/javascript/rhino/testing/AbstractStaticScope.java @ closure-compiler v20261015

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
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhinotesting

import "github.com/paullepoulpe/jsgo/internal/rhino"

// AbstractStaticScope is a basic implementation of StaticScope for testing.
type AbstractStaticScope interface {
	rhino.StaticTypedScope
}

// AbstractStaticScopeBase holds the methods shared by AbstractStaticScope
// implementations; the implementation provides GetSlot.
type AbstractStaticScopeBase struct {
	self AbstractStaticScope
}

// InitAbstractStaticScopeBase sets the most-derived receiver.
func (s *AbstractStaticScopeBase) InitAbstractStaticScopeBase(self AbstractStaticScope) {
	s.self = self
}

// Java: AbstractStaticScope.getRootNode()
func (s *AbstractStaticScopeBase) GetRootNode() rhino.Node {
	return rhino.Node{}
}

// Java: AbstractStaticScope.getParentScope()
func (s *AbstractStaticScopeBase) GetParentScope() rhino.StaticScope {
	return nil
}

// Java: AbstractStaticScope.getOwnSlot(String)
func (s *AbstractStaticScopeBase) GetOwnSlot(name string) rhino.StaticSlot {
	return s.self.GetSlot(name)
}

// Java: AbstractStaticScope.getTypeOfThis()
func (s *AbstractStaticScopeBase) GetTypeOfThis() rhino.JSType {
	return nil
}

// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (s *AbstractStaticScopeBase) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s.self, name)
}

// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *AbstractStaticScopeBase) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s.self, qname)
}
