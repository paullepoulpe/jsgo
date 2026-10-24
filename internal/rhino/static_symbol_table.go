// Port of com/google/javascript/rhino/StaticSymbolTable.java @ closure-compiler v20261021

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

package rhino

import "iter"

// StaticSymbolTable looks up references by the symbols that they refer to.
type StaticSymbolTable[S StaticSlot, R StaticRef] interface {
	// GetReferences returns the references that point to the given symbol.
	GetReferences(symbol S) iter.Seq[R]

	// GetScope returns the scope for a given symbol.
	GetScope(symbol S) StaticScope

	// GetAllSymbols returns all variables in this symbol table.
	GetAllSymbols() iter.Seq[S]
}
